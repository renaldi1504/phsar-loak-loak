// Package paygw: payment gateway yang diatur super admin dari halaman admin.
//
// Sumber pengaturan (urutan prioritas):
//  1. database — platform_settings['payment'] (mode, lama QR berlaku) + payment_gateways (kredensial terenkripsi);
//  2. .env server (PAYMENT_PROVIDER, PAYWAY_*) — dipakai kalau super admin belum pernah menyimpan.
//
// Store memenuhi payment.Provider, jadi layanan belanja tidak perlu tahu dari mana kredensialnya.
// Pengaturan di-cache sebentar (cacheTTL) supaya tidak membaca database di setiap permintaan;
// perubahan dari admin langsung berlaku di proses yang sama dan paling lambat cacheTTL di proses lain (worker).
package paygw

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"phsar/internal/payment"
	"phsar/internal/platform/db"
	"phsar/internal/secretbox"
)

const (
	ModeOff    = "off"    // pra-peluncuran: checkout ditutup
	ModePayWay = "payway" // ABA PayWay (KHQR)
	ModeFake   = "fake"   // hanya pengembangan: tombol "simulasi bayar"

	SandboxURL    = "https://checkout-sandbox.payway.com.kh"
	ProductionURL = "https://checkout.payway.com.kh"

	cacheTTL = 15 * time.Second
)

var (
	ErrDisabled = errors.New("payment gateway nonaktif")
	merchantRe  = regexp.MustCompile(`^[A-Za-z0-9._-]{1,30}$`)
)

type querier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// FieldError: input admin tidak valid (dikirim balik dengan nama kolomnya).
type FieldError struct{ Field, Msg string }

func (e FieldError) Error() string { return e.Field + ": " + e.Msg }

// Defaults: nilai dari .env server.
type Defaults struct {
	Env               string // development · staging · production
	Mode              string // PAYMENT_PROVIDER
	PayWayBaseURL     string
	PayWayMerchantID  string
	PayWayAPIKey      string
	LifetimeMin       int
	PublicCallbackURL string // https://<domain>/webhooks/payway (ditampilkan ke admin untuk diberikan ke ABA)
}

type Store struct {
	Pool *pgxpool.Pool
	Box  *secretbox.Box
	Def  Defaults
	Fake *payment.Fake // dibagi satu proses (pembayaran palsu disimpan di memori)
	Log  *slog.Logger

	mu  sync.Mutex
	cur *active
	at  time.Time
}

type active struct {
	mode     string
	source   string // admin · env
	lifetime int
	prov     payment.Provider // nil kalau off
}

// ---- payment.Provider ----

func (s *Store) Name() string { return s.load(context.Background()).mode }

func (s *Store) CreateQR(ctx context.Context, r payment.QRRequest) (payment.QR, error) {
	a := s.load(ctx)
	if a.prov == nil {
		return payment.QR{}, ErrDisabled
	}
	return a.prov.CreateQR(ctx, r)
}

func (s *Store) Check(ctx context.Context, tranID string) (payment.Status, error) {
	a := s.load(ctx)
	if a.prov == nil {
		return payment.Status{}, ErrDisabled
	}
	return a.prov.Check(ctx, tranID)
}

func (s *Store) VerifyCallback(h http.Header, body []byte) bool {
	a := s.load(context.Background())
	return a.prov != nil && a.prov.VerifyCallback(h, body)
}

// Open: pembeli boleh checkout (bukan mode off).
func (s *Store) Open(ctx context.Context) bool { return s.load(ctx).mode != ModeOff }

// Test: mode pembayaran palsu (tombol "simulasi bayar").
func (s *Store) Test(ctx context.Context) bool { return s.load(ctx).mode == ModeFake }

// Lifetime: berapa menit QR berlaku.
func (s *Store) Lifetime(ctx context.Context) int { return s.load(ctx).lifetime }

// FakeProvider: penyedia palsu kalau sedang mode fake (untuk endpoint /dev/pay).
func (s *Store) FakeProvider(ctx context.Context) (*payment.Fake, bool) {
	if s.load(ctx).mode != ModeFake {
		return nil, false
	}
	return s.Fake, s.Fake != nil
}

// Invalidate: baca ulang dari database pada permintaan berikutnya.
func (s *Store) Invalidate() {
	s.mu.Lock()
	s.cur = nil
	s.mu.Unlock()
}

func (s *Store) load(ctx context.Context) *active {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cur != nil && time.Since(s.at) < cacheTTL {
		return s.cur
	}
	a, err := s.read(ctx)
	if err != nil {
		if s.cur != nil { // database sesaat bermasalah: pakai pengaturan terakhir
			s.Log.Warn("pengaturan pembayaran tidak bisa dibaca, memakai yang terakhir", "err", err)
			return s.cur
		}
		s.Log.Error("pengaturan pembayaran tidak bisa dibaca, checkout ditutup sementara", "err", err)
		return &active{mode: ModeOff, source: "env", lifetime: s.defLifetime()}
	}
	s.cur, s.at = a, time.Now()
	return a
}

func (s *Store) defLifetime() int {
	if s.Def.LifetimeMin >= 3 && s.Def.LifetimeMin <= 120 {
		return s.Def.LifetimeMin
	}
	return 15
}

type general struct {
	Mode        string `json:"mode"`
	LifetimeMin int    `json:"lifetime_min"`
}

type payway struct {
	Environment string
	MerchantID  string
	BaseURL     string
	APIKey      string
	KeyHint     string
	Source      string // admin · env · none
	VerifiedAt  *time.Time
	LastTest    *TestResult
}

func (s *Store) readGeneral(ctx context.Context, q querier) (general, bool, error) {
	var g general
	var raw []byte
	err := q.QueryRow(ctx, `SELECT value FROM platform_settings WHERE key = 'payment'`).Scan(&raw)
	if errors.Is(err, pgx.ErrNoRows) {
		return general{Mode: s.Def.Mode, LifetimeMin: s.defLifetime()}, false, nil
	}
	if err != nil {
		return g, false, err
	}
	if err := json.Unmarshal(raw, &g); err != nil {
		return g, false, err
	}
	if g.LifetimeMin < 3 || g.LifetimeMin > 120 {
		g.LifetimeMin = s.defLifetime()
	}
	return g, true, nil
}

// readPayWay: kredensial dari database kalau ada, kalau tidak dari .env. withKey=false → API key tidak didekripsi.
func (s *Store) readPayWay(ctx context.Context, q querier, withKey bool) (payway, error) {
	var env string
	var cfgRaw, sealed, lastRaw []byte
	var hint *string
	var verified *time.Time
	err := q.QueryRow(ctx, `SELECT environment, config, secret_enc, secret_hint, verified_at, last_test FROM payment_gateways WHERE code = 'payway'`).
		Scan(&env, &cfgRaw, &sealed, &hint, &verified, &lastRaw)
	if errors.Is(err, pgx.ErrNoRows) {
		p := payway{Environment: "sandbox", MerchantID: s.Def.PayWayMerchantID, BaseURL: s.Def.PayWayBaseURL, Source: "none"}
		if s.Def.PayWayMerchantID != "" && s.Def.PayWayAPIKey != "" {
			p.Source, p.KeyHint = "env", hintOf(s.Def.PayWayAPIKey)
			if withKey {
				p.APIKey = s.Def.PayWayAPIKey
			}
			if !strings.Contains(p.BaseURL, "sandbox") {
				p.Environment = "production"
			}
		}
		return p, nil
	}
	if err != nil {
		return payway{}, err
	}
	var c struct {
		MerchantID string `json:"merchant_id"`
		BaseURL    string `json:"base_url"`
	}
	_ = json.Unmarshal(cfgRaw, &c)
	p := payway{Environment: env, MerchantID: c.MerchantID, BaseURL: c.BaseURL, Source: "admin", VerifiedAt: verified}
	if p.BaseURL == "" {
		p.BaseURL = map[string]string{"sandbox": SandboxURL, "production": ProductionURL}[env]
	}
	if hint != nil {
		p.KeyHint = *hint
	}
	if len(lastRaw) > 0 {
		var t TestResult
		if json.Unmarshal(lastRaw, &t) == nil {
			p.LastTest = &t
		}
	}
	if withKey && sealed != nil {
		k, err := s.Box.Open(sealed, aad("payway"))
		if err != nil {
			return payway{}, fmt.Errorf("API key PayWay tidak bisa dibuka (DATA_ENCRYPTION_KEY berubah?): %w", err)
		}
		p.APIKey = string(k)
	}
	return p, nil
}

func (s *Store) read(ctx context.Context) (*active, error) {
	g, fromDB, err := s.readGeneral(ctx, s.Pool)
	if err != nil {
		return nil, err
	}
	a := &active{mode: g.Mode, source: map[bool]string{true: "admin", false: "env"}[fromDB], lifetime: g.LifetimeMin}
	switch a.mode {
	case ModePayWay:
		p, err := s.readPayWay(ctx, s.Pool, true)
		if err != nil {
			return nil, err
		}
		if p.MerchantID == "" || p.APIKey == "" {
			s.Log.Error("mode PayWay tapi kredensial belum lengkap — checkout ditutup")
			a.mode = ModeOff
			break
		}
		a.prov = &payment.PayWay{BaseURL: strings.TrimRight(p.BaseURL, "/"), MerchantID: p.MerchantID, APIKey: p.APIKey}
	case ModeFake:
		if s.Def.Env == "production" || s.Fake == nil {
			a.mode = ModeOff
			break
		}
		a.prov = s.Fake
	default:
		a.mode = ModeOff
	}
	return a, nil
}

func aad(code string) []byte { return []byte("payment-gateway:" + code) }

func hintOf(key string) string {
	if len(key) <= 4 {
		return "••••"
	}
	return key[len(key)-4:]
}

// ---- tampilan & perubahan oleh super admin ----

type TestResult struct {
	OK      bool      `json:"ok"`
	Message string    `json:"message"`
	At      time.Time `json:"at"`
}

type PayWayView struct {
	Environment string      `json:"environment"`
	MerchantID  string      `json:"merchant_id"`
	BaseURL     string      `json:"base_url"`
	KeySet      bool        `json:"key_set"`
	KeyHint     string      `json:"key_hint"`
	Source      string      `json:"source"` // admin · env · none
	VerifiedAt  *time.Time  `json:"verified_at"`
	LastTest    *TestResult `json:"last_test"`
}

type View struct {
	Mode        string     `json:"mode"`   // yang tersimpan (off · payway · fake)
	Active      string     `json:"active"` // yang benar-benar berjalan (off kalau kredensial kurang)
	Source      string     `json:"source"` // admin · env
	EnvMode     string     `json:"env_mode"`
	LifetimeMin int        `json:"lifetime_min"`
	Production  bool       `json:"production"`
	CallbackURL string     `json:"callback_url"`
	PayWay      PayWayView `json:"payway"`
}

func (s *Store) View(ctx context.Context) (View, error) {
	g, fromDB, err := s.readGeneral(ctx, s.Pool)
	if err != nil {
		return View{}, err
	}
	p, err := s.readPayWay(ctx, s.Pool, false)
	if err != nil {
		return View{}, err
	}
	s.Invalidate()
	return View{
		Mode: g.Mode, Active: s.load(ctx).mode, Source: map[bool]string{true: "admin", false: "env"}[fromDB], EnvMode: s.Def.Mode,
		LifetimeMin: g.LifetimeMin, Production: s.Def.Env == "production", CallbackURL: s.Def.PublicCallbackURL,
		PayWay: PayWayView{Environment: p.Environment, MerchantID: p.MerchantID, BaseURL: p.BaseURL, KeySet: p.KeyHint != "",
			KeyHint: p.KeyHint, Source: p.Source, VerifiedAt: p.VerifiedAt, LastTest: p.LastTest},
	}, nil
}

// SaveGeneral: mode aktif + lama QR. mode "" = kembali ikut .env server.
func (s *Store) SaveGeneral(ctx context.Context, actorID, mode string, lifetime int) error {
	defer s.Invalidate()
	if mode == "" {
		return db.WithActor(ctx, s.Pool, actorID, "admin", "", func(tx pgx.Tx) error {
			_, err := tx.Exec(ctx, `DELETE FROM platform_settings WHERE key = 'payment'`)
			return err
		})
	}
	switch mode {
	case ModeOff, ModePayWay:
	case ModeFake:
		if s.Def.Env == "production" {
			return FieldError{"mode", "test payments are not allowed in production"}
		}
	default:
		return FieldError{"mode", "unknown mode"}
	}
	if lifetime < 3 || lifetime > 120 {
		return FieldError{"lifetime_min", "must be 3–120 minutes"}
	}
	return db.WithActor(ctx, s.Pool, actorID, "admin", "", func(tx pgx.Tx) error {
		if mode == ModePayWay {
			p, err := s.readPayWay(ctx, tx, false)
			if err != nil {
				return err
			}
			if p.MerchantID == "" || p.KeyHint == "" {
				return FieldError{"mode", "enter the PayWay merchant ID and API key first"}
			}
		}
		v, _ := json.Marshal(general{Mode: mode, LifetimeMin: lifetime})
		_, err := tx.Exec(ctx, `INSERT INTO platform_settings (key, value, updated_by, updated_at) VALUES ('payment', $1, $2, now())
			ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_by = EXCLUDED.updated_by, updated_at = now()`, v, actorID)
		return err
	})
}

type PayWayInput struct {
	Environment string  `json:"environment"`
	MerchantID  string  `json:"merchant_id"`
	BaseURL     string  `json:"base_url"`
	APIKey      *string `json:"api_key"` // nil / kosong = tidak diubah
}

// SavePayWay: kredensial PayWay. API key dienkripsi; kalau kredensial berubah, status "terverifikasi" direset.
func (s *Store) SavePayWay(ctx context.Context, actorID string, in PayWayInput) error {
	defer s.Invalidate()
	in.MerchantID = strings.TrimSpace(in.MerchantID)
	in.BaseURL = strings.TrimRight(strings.TrimSpace(in.BaseURL), "/")
	if in.Environment != "sandbox" && in.Environment != "production" {
		return FieldError{"environment", "choose sandbox or production"}
	}
	if !merchantRe.MatchString(in.MerchantID) {
		return FieldError{"merchant_id", "letters, numbers, . _ - (max 30)"}
	}
	if in.BaseURL != "" && !allowedBaseURL(in.BaseURL) {
		return FieldError{"base_url", "must be an https://….payway.com.kh address"}
	}
	var sealed []byte
	var hint string
	if in.APIKey != nil && strings.TrimSpace(*in.APIKey) != "" {
		k := strings.TrimSpace(*in.APIKey)
		if len(k) < 8 || len(k) > 256 || strings.ContainsAny(k, " \t\r\n") {
			return FieldError{"api_key", "looks wrong (8–256 characters, no spaces)"}
		}
		var err error
		if sealed, err = s.Box.Seal([]byte(k), aad("payway")); err != nil {
			return err
		}
		hint = hintOf(k)
	}
	cfg, _ := json.Marshal(map[string]string{"merchant_id": in.MerchantID, "base_url": in.BaseURL})
	return db.WithActor(ctx, s.Pool, actorID, "admin", "", func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `INSERT INTO payment_gateways (code, environment, config, secret_enc, secret_hint, updated_by)
			VALUES ('payway', $1, $2, $3, nullif($4, ''), $5)
			ON CONFLICT (code) DO UPDATE SET
				environment = EXCLUDED.environment,
				config = EXCLUDED.config,
				secret_enc = coalesce(EXCLUDED.secret_enc, payment_gateways.secret_enc),
				secret_hint = coalesce(EXCLUDED.secret_hint, payment_gateways.secret_hint),
				verified_at = CASE WHEN EXCLUDED.secret_enc IS NOT NULL OR payment_gateways.environment <> EXCLUDED.environment
					OR payment_gateways.config <> EXCLUDED.config THEN NULL ELSE payment_gateways.verified_at END,
				updated_by = EXCLUDED.updated_by`, in.Environment, cfg, sealed, hint, actorID)
		return err
	})
}

// allowedBaseURL: hanya server PayWay (API key ikut dikirim ke alamat ini, jadi tidak boleh sembarang).
func allowedBaseURL(raw string) bool {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "https" || u.User != nil || (u.Path != "" && u.Path != "/") || u.RawQuery != "" {
		return false
	}
	h := strings.ToLower(u.Hostname())
	return h == "payway.com.kh" || strings.HasSuffix(h, ".payway.com.kh")
}

// TestPayWay: cek kredensial dengan menanyakan status transaksi yang tidak ada (tidak membuat pembayaran).
// "Tidak ditemukan" = kunci & merchant diterima PayWay.
func (s *Store) TestPayWay(ctx context.Context, actorID string) (TestResult, error) {
	p, err := s.readPayWay(ctx, s.Pool, true)
	if err != nil {
		return TestResult{}, err
	}
	if p.MerchantID == "" || p.APIKey == "" {
		return TestResult{}, FieldError{"api_key", "enter the merchant ID and API key first"}
	}
	pw := &payment.PayWay{BaseURL: strings.TrimRight(p.BaseURL, "/"), MerchantID: p.MerchantID, APIKey: p.APIKey}
	tctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	res := TestResult{At: time.Now().UTC()}
	if _, err := pw.Check(tctx, fmt.Sprintf("PTEST%d", time.Now().Unix())); err != nil {
		res.Message = explain(err.Error())
	} else {
		res.OK, res.Message = true, "PayWay accepted the merchant ID and API key."
	}
	if p.Source == "admin" {
		raw, _ := json.Marshal(res)
		err = db.WithActor(ctx, s.Pool, actorID, "admin", "", func(tx pgx.Tx) error {
			_, err := tx.Exec(ctx, `UPDATE payment_gateways SET last_test = $1, verified_at = CASE WHEN $2 THEN now() ELSE NULL END
				WHERE code = 'payway'`, raw, res.OK)
			return err
		})
	}
	return res, err
}

// explain: pesan PayWay yang lebih mudah dipahami admin (pesan asli tetap disertakan).
func explain(msg string) string {
	switch {
	case strings.Contains(msg, "kode 1:"):
		return "Wrong API key or merchant ID (PayWay code 1). " + msg
	case strings.Contains(msg, "kode 102:"), strings.Contains(msg, "whitelist"):
		return "The server IP/domain is not whitelisted by ABA yet. " + msg
	case strings.Contains(msg, "tidak bisa menghubungi"):
		return "Cannot reach PayWay from this server (network/firewall). " + msg
	}
	return msg
}
