// Package qa: tes integrasi end-to-end untuk seluruh backend Phsar.
//
// Menjalankan server sungguhan (httptest) di atas database PostgreSQL khusus tes, lalu memanggil
// API persis seperti app/Seller Center/admin. Dijalankan lewat `make qa` (lihat QA.md).
//
// Butuh PHSAR_TEST_DATABASE_URL yang nama database-nya berakhiran _test, karena isinya DIHAPUS
// setiap kali tes dijalankan. Tanpa variabel itu, tes ini dilewati.
package qa

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"

	"phsar/internal/auth"
	"phsar/internal/config"
	"phsar/internal/httpx"
	"phsar/internal/media"
	"phsar/internal/platform/db"
	"phsar/internal/storage"
	"phsar/internal/translate"
	"phsar/migrations"
)

// ---- Lingkungan tes bersama (dibuat sekali untuk semua tes) ----

type env struct {
	media  *media.Worker // nil kalau ffmpeg tidak ada
	srv    *httptest.Server
	pool   *pgxpool.Pool
	sms    *captureSMS
	worker *translate.Worker
	cfg    config.Config
}

var (
	shared  *env
	setupMu sync.Mutex
	phoneN  atomic.Int64
)

type captureSMS struct {
	mu   sync.Mutex
	last map[string]string
}

var sixDigits = regexp.MustCompile(`\d{6}`)

func (c *captureSMS) Send(_ context.Context, phone, text string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.last[phone] = sixDigits.FindString(text)
	return nil
}

func (c *captureSMS) code(phone string) string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.last[phone]
}

func getEnv(t *testing.T) *env {
	t.Helper()
	setupMu.Lock()
	defer setupMu.Unlock()
	if shared != nil {
		return shared
	}
	dbURL := os.Getenv("PHSAR_TEST_DATABASE_URL")
	if dbURL == "" {
		t.Skip("PHSAR_TEST_DATABASE_URL belum diisi — jalankan lewat `make qa`")
	}
	u, err := url.Parse(dbURL)
	if err != nil || !strings.HasSuffix(strings.TrimPrefix(u.Path, "/"), "_test") {
		t.Fatalf("demi keamanan, nama database tes harus berakhiran _test (sekarang: %q)", u.Path)
	}
	ctx := context.Background()

	// 1) Kosongkan database, lalu cek migrasi: up → reset (down semua) → up.
	pool, err := db.Open(ctx, dbURL)
	if err != nil && strings.Contains(err.Error(), "3D000") {
		// Database tes belum ada → buat lewat database utama di server yang sama, lalu konek lagi.
		if cerr := createTestDB(ctx, u); cerr != nil {
			t.Fatalf("database tes belum ada dan gagal dibuat: %v", cerr)
		}
		pool, err = db.Open(ctx, dbURL)
	}
	if err != nil {
		t.Fatalf("tidak bisa konek ke database tes: %v", err)
	}
	if _, err := pool.Exec(ctx, `DROP SCHEMA public CASCADE; CREATE SCHEMA public`); err != nil {
		t.Fatalf("gagal mengosongkan database tes: %v", err)
	}
	pool.Close()
	if os.Getenv("QA_VERBOSE") == "" {
		goose.SetLogger(goose.NopLogger())
	}
	for _, step := range []string{"up", "reset", "up"} {
		if err := migrations.Run(ctx, dbURL, step); err != nil {
			t.Fatalf("migrasi %s gagal: %v", step, err)
		}
	}

	pool, err = db.Open(ctx, dbURL)
	if err != nil {
		t.Fatal(err)
	}
	cfg := config.Config{
		Env: "development", PublicBaseURL: "http://qa.local", DatabaseURL: dbURL, DefaultLang: "en", LogLevel: "error",
		JWTSecret: "qa-secret-0123456789-0123456789-0123456789", TelegramMode: "off", SMSProvider: "log",
		StorageDir: t.TempDir(), Translator: "fake", GoogleMapsKey: "qa-maps-key", GoogleMapsMapID: "DEMO_MAP_ID",
	}
	if dir, err := os.MkdirTemp("", "phsar-qa-uploads-"); err == nil {
		cfg.StorageDir = dir // t.TempDir() milik tes pertama saja; folder ini dipakai semua tes
	}
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	if os.Getenv("QA_VERBOSE") != "" {
		log = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}
	authSvc := auth.NewService(pool, cfg, log)
	sms := &captureSMS{last: map[string]string{}}
	authSvc.SetSMSSender(sms)
	router := httpx.NewRouter(cfg, log, pool, authSvc)
	// Setiap balasan 5xx dicatat; TestMain menggagalkan QA kalau ada (500 = bug, bukan validasi).
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &statusRecorder{ResponseWriter: w, status: 200}
		router.ServeHTTP(rec, r)
		if rec.status >= 500 {
			serverErrors.Lock()
			serverErrors.list = append(serverErrors.list, fmt.Sprintf("%d %s %s", rec.status, r.Method, r.URL.Path))
			serverErrors.Unlock()
		}
	}))
	shared = &env{srv: srv, pool: pool, sms: sms, cfg: cfg,
		worker: &translate.Worker{Pool: pool, T: translate.Fake{}, Log: log}}
	if media.Available() {
		shared.media = &media.Worker{Pool: pool, Files: storage.Local{Dir: cfg.StorageDir, BaseURL: cfg.PublicBaseURL + "/media"}, Log: log}
		makeRealFixtures(t)
	}
	return shared
}

var serverErrors struct {
	sync.Mutex
	list []string
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) { s.status = code; s.ResponseWriter.WriteHeader(code) }

func TestMain(m *testing.M) {
	code := m.Run()
	serverErrors.Lock()
	if len(serverErrors.list) > 0 {
		fmt.Fprintln(os.Stderr, "\nQA GAGAL: server membalas error 5xx (bug) pada:")
		for _, l := range serverErrors.list {
			fmt.Fprintln(os.Stderr, "  -", l)
		}
		code = 1
	}
	serverErrors.Unlock()
	if shared != nil {
		shared.srv.Close()
		shared.pool.Close()
		_ = os.RemoveAll(shared.cfg.StorageDir)
	}
	os.Exit(code)
}

// createTestDB membuat database <nama>_test dengan konek ke database "phsar" (atau "postgres") di server yang sama.
func createTestDB(ctx context.Context, u *url.URL) error {
	name := strings.TrimPrefix(u.Path, "/")
	if !regexp.MustCompile(`^[a-z0-9_]+_test$`).MatchString(name) {
		return fmt.Errorf("nama database tes tidak valid: %q", name)
	}
	var lastErr error
	for _, adminDB := range []string{"phsar", "postgres"} {
		au := *u
		au.Path = "/" + adminDB
		pool, err := db.Open(ctx, au.String())
		if err != nil {
			lastErr = err
			continue
		}
		_, err = pool.Exec(ctx, `CREATE DATABASE "`+name+`"`)
		pool.Close()
		if err == nil || strings.Contains(err.Error(), "42P04") { // 42P04 = sudah ada
			return nil
		}
		lastErr = err
	}
	return lastErr
}

// ---- Klien HTTP kecil ----

type resp struct {
	t      *testing.T
	Status int
	Body   map[string]any
	Raw    []byte
	Header http.Header
}

func (e *env) call(t *testing.T, method, path string, body any, token string, headers ...string) resp {
	t.Helper()
	var r io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		r = bytes.NewReader(b)
	}
	req, _ := http.NewRequest(method, e.srv.URL+path, r)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	for i := 0; i+1 < len(headers); i += 2 {
		req.Header.Set(headers[i], headers[i+1])
	}
	return e.send(t, req)
}

func (e *env) send(t *testing.T, req *http.Request) resp {
	t.Helper()
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("%s %s: %v", req.Method, req.URL.Path, err)
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	out := resp{t: t, Status: res.StatusCode, Raw: raw, Header: res.Header}
	_ = json.Unmarshal(raw, &out.Body)
	return out
}

func (e *env) upload(t *testing.T, token, kind string, data []byte) resp {
	t.Helper()
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	_ = mw.WriteField("kind", kind)
	fw, _ := mw.CreateFormFile("file", "file.bin")
	_, _ = fw.Write(data)
	_ = mw.Close()
	req, _ := http.NewRequest(http.MethodPost, e.srv.URL+"/v1/seller/uploads", &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	return e.send(t, req)
}

// login membuat user baru lewat SMS OTP dan mengembalikan (access token, nomor HP).
func (e *env) login(t *testing.T) (string, string) {
	t.Helper()
	phone := fmt.Sprintf("0970%06d", phoneN.Add(1)) // database tes dikosongkan tiap run, jadi cukup berurutan
	r := e.call(t, "POST", "/v1/app/auth/otp/request", map[string]string{"phone": phone}, "")
	r.want(200)
	r = e.call(t, "POST", "/v1/app/auth/otp/verify", map[string]string{"phone": phone, "code": e.sms.code(r.str("phone"))}, "")
	r.want(200)
	return r.str("access_token"), phone
}

func (e *env) sql(t *testing.T, q string, args ...any) {
	t.Helper()
	if _, err := e.pool.Exec(context.Background(), q, args...); err != nil {
		t.Fatalf("sql %q: %v", q, err)
	}
}

func (e *env) scalar(t *testing.T, q string, args ...any) string {
	t.Helper()
	var v any
	if err := e.pool.QueryRow(context.Background(), q, args...).Scan(&v); err != nil {
		t.Fatalf("sql %q: %v", q, err)
	}
	return fmt.Sprint(v)
}

// drainTranslations menjalankan worker terjemahan sampai antrean kosong.
func (e *env) drainTranslations(t *testing.T) {
	t.Helper()
	for i := 0; i < 50; i++ {
		did, err := e.worker.RunOnce(context.Background())
		if err != nil {
			t.Fatalf("worker terjemahan: %v", err)
		}
		if !did {
			return
		}
	}
	t.Fatal("antrean terjemahan tidak habis-habis")
}

// ---- Pemeriksa ----

func (r resp) want(status int) resp {
	r.t.Helper()
	if r.Status != status {
		r.t.Fatalf("status %d, harusnya %d: %s", r.Status, status, truncate(string(r.Raw)))
	}
	return r
}

// wantErr: status + kode error (+ field, kalau diisi).
func (r resp) wantErr(status int, code string, field ...string) resp {
	r.t.Helper()
	r.want(status)
	if got := r.str("error.code"); got != code {
		r.t.Fatalf("kode error %q, harusnya %q: %s", got, code, truncate(string(r.Raw)))
	}
	if len(field) > 0 && r.str("error.field") != field[0] {
		r.t.Fatalf("field error %q, harusnya %q: %s", r.str("error.field"), field[0], truncate(string(r.Raw)))
	}
	return r
}

// get mengambil nilai dengan path "a.b.0.c".
func (r resp) get(path string) any {
	var cur any = r.Body
	for _, p := range strings.Split(path, ".") {
		switch c := cur.(type) {
		case map[string]any:
			cur = c[p]
		case []any:
			i, err := strconv.Atoi(p)
			if err != nil || i >= len(c) {
				return nil
			}
			cur = c[i]
		default:
			return nil
		}
	}
	return cur
}

func (r resp) str(path string) string {
	v := r.get(path)
	if v == nil {
		return ""
	}
	if f, ok := v.(float64); ok {
		return strconv.FormatFloat(f, 'f', -1, 64)
	}
	return fmt.Sprint(v)
}

func (r resp) len(path string) int {
	switch v := r.get(path).(type) {
	case []any:
		return len(v)
	case map[string]any:
		return len(v)
	}
	return 0
}

func (r resp) eq(path, want string) resp {
	r.t.Helper()
	if got := r.str(path); got != want {
		r.t.Fatalf("%s = %q, harusnya %q", path, got, want)
	}
	return r
}

func check(t *testing.T, ok bool, msg string, args ...any) {
	t.Helper()
	if !ok {
		t.Fatalf(msg, args...)
	}
}

func truncate(s string) string {
	if len(s) > 600 {
		return s[:600] + "…"
	}
	return s
}

// drainMedia menjalankan worker media sampai antrean kosong (kalau ffmpeg ada).
func (e *env) drainMedia(t *testing.T) {
	t.Helper()
	if e.media == nil {
		return
	}
	for i := 0; i < 50; i++ {
		did, err := e.media.RunOnce(context.Background())
		if err != nil {
			t.Fatalf("worker media: %v", err)
		}
		if !did {
			return
		}
	}
}

// makeRealFixtures: dengan ffmpeg, upload video diperiksa sungguhan (durasi), jadi file uji harus video/foto asli.
func makeRealFixtures(t *testing.T) {
	dir, err := os.MkdirTemp("", "phsar-qa-fixtures-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	gen := func(name string, args ...string) []byte {
		p := filepath.Join(dir, name)
		cmd := exec.Command("ffmpeg", append(append([]string{"-hide_banner", "-loglevel", "error", "-y"}, args...), p)...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("buat file uji %s: %v %s", name, err, out)
		}
		b, _ := os.ReadFile(p)
		return b
	}
	src := []string{"-f", "lavfi", "-i", "testsrc2=size=540x960:rate=24", "-f", "lavfi", "-i", "sine=frequency=330", "-t", "2", "-shortest",
		"-c:v", "libx264", "-preset", "ultrafast", "-c:a", "aac"}
	mp4Bytes = gen("v.mp4", src...)
	movBytes = gen("v.mov", append(src, "-f", "mov")...)
	jpegBytes = gen("p.jpg", "-f", "lavfi", "-i", "testsrc2=size=2400x1800", "-frames:v", "1")
	pngBytes = gen("p.png", "-f", "lavfi", "-i", "testsrc2=size=640x480", "-frames:v", "1")
	longVideo = gen("long.mp4", "-f", "lavfi", "-i", "color=c=black:size=64x64:rate=1", "-t", "200", "-c:v", "libx264", "-preset", "ultrafast")
}

// ---- File contoh (isi minimal yang lolos deteksi jenis file; diganti file asli kalau ffmpeg ada) ----

var (
	jpegBytes = append([]byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 'J', 'F', 'I', 'F', 0x00, 0x01}, make([]byte, 600)...)
	pngBytes  = append([]byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1A, '\n'}, make([]byte, 600)...)
	mp4Bytes  = append([]byte{0, 0, 0, 24, 'f', 't', 'y', 'p', 'i', 's', 'o', 'm', 0, 0, 2, 0, 'i', 's', 'o', 'm', 'i', 's', 'o', '2'}, make([]byte, 2000)...)
	movBytes  = append([]byte{0, 0, 0, 20, 'f', 't', 'y', 'p', 'q', 't', ' ', ' ', 0, 0, 2, 0, 'q', 't', ' ', ' '}, make([]byte, 2000)...)
	pdfBytes  = append([]byte("%PDF-1.4\n"), make([]byte, 600)...)
	longVideo []byte
	htmlJPEG  = append([]byte("<html><script>alert(1)</script></html>"), bytes.Repeat([]byte(" "), 600)...)
)
