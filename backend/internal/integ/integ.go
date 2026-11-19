// Package integ: layanan luar untuk login yang diatur super admin dari halaman admin —
// provider SMS (PlasGate, untuk kode OTP) dan Login Facebook. Kredensial disimpan terenkripsi.
//
// Prioritas: pengaturan di database; kalau belum pernah disimpan, mode SMS ikut SMS_PROVIDER di .env.
// Store memenuhi auth.SMSGateway dan auth.FacebookConfig, jadi paket auth tidak tahu asal kredensialnya.
package integ

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"phsar/internal/platform/db"
	"phsar/internal/secretbox"
)

const (
	SMSOff      = "off"      // login SMS mati (hanya Telegram / Facebook)
	SMSLog      = "log"      // pengembangan: kode ditulis ke log server
	SMSPlasGate = "plasgate" // PlasGate (Kamboja)

	PlasGateURL = "https://cloudapi.plasgate.com/rest/send"
	GraphURL    = "https://graph.facebook.com/v21.0"
	DialogURL   = "https://www.facebook.com/v21.0/dialog/oauth"
)

// CacheTTL: berapa lama pengaturan disimpan di memori proses (tes memakai 0).
var CacheTTL = 15 * time.Second

// Endpoints: alamat API luar. Hanya diganti oleh tes (server palsu); di produksi selalu alamat resmi.
var Endpoints = struct{ PlasGate, Graph, Dialog string }{PlasGateURL, GraphURL, DialogURL}

var (
	ErrSMSOff  = errors.New("SMS login is not available")
	appIDRe    = regexp.MustCompile(`^[0-9]{5,25}$`)
	senderRe   = regexp.MustCompile(`^[A-Za-z0-9 ._-]{2,11}$|^\+?[0-9]{6,15}$`)
	keyCharsRe = regexp.MustCompile(`^[\x21-\x7e]{8,256}$`)
)

// FieldError: input admin tidak valid.
type FieldError struct{ Field, Msg string }

func (e FieldError) Error() string { return e.Field + ": " + e.Msg }

type Defaults struct {
	Env         string // development · staging · production
	SMSProvider string // SMS_PROVIDER dari .env (log · off)
	BaseURL     string // PUBLIC_BASE_URL, untuk daftar redirect URI Facebook / Google
	SellerURL   string // SELLER_BASE_URL (subdomain Seller Center), kosong = BaseURL
	AdminURL    string // ADMIN_BASE_URL (subdomain admin), kosong = BaseURL
	// Cadangan dari .env kalau belum diisi di admin (halaman Integrasi).
	AnthropicKey, AnthropicModel string
	GeminiKey, GeminiModel       string
	MapsKey, MapsMapID           string
	TelegramSet                  bool   // TELEGRAM_BOT_TOKEN terisi (token tetap di .env karena bot berjalan sejak start)
	TelegramUsername             string // TELEGRAM_BOT_USERNAME
}

type Store struct {
	Pool *pgxpool.Pool
	Box  *secretbox.Box
	Def  Defaults
	Log  *slog.Logger
	HTTP *http.Client
	// Alamat API (diganti di tes).
	PlasGateURL, GraphURL, DialogURL string

	mu  sync.Mutex
	cur *state
	at  time.Time
}

type plasgate struct{ PrivateKey, Secret, Sender string }

type state struct {
	smsMode string
	pg      plasgate
	fbOn    bool
	fbID    string
	fbKey   string
	ai      aiState // kunci AI & Google Maps (ai.go)
	google  googleState
	apple   appleState
}

func (s *Store) client() *http.Client {
	if s.HTTP != nil {
		return s.HTTP
	}
	return &http.Client{Timeout: 15 * time.Second}
}

func (s *Store) url(v, def string) string {
	if v != "" {
		return v
	}
	return def
}

// Invalidate: baca ulang pada permintaan berikutnya.
func (s *Store) Invalidate() {
	s.mu.Lock()
	s.cur = nil
	s.mu.Unlock()
}

func (s *Store) load(ctx context.Context) *state {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cur != nil && time.Since(s.at) < CacheTTL {
		return s.cur
	}
	st, err := s.read(ctx)
	if err != nil {
		if s.cur != nil {
			s.Log.Warn("pengaturan login/SMS tidak bisa dibaca, memakai yang terakhir", "err", err)
			return s.cur
		}
		s.Log.Error("pengaturan login/SMS tidak bisa dibaca", "err", err)
		fb := &state{smsMode: s.envSMS()}
		fb.ai.keys = AIKeys{AnthropicKey: s.Def.AnthropicKey, AnthropicModel: s.Def.AnthropicModel, GeminiKey: s.Def.GeminiKey, GeminiModel: s.Def.GeminiModel}
		fb.ai.mapsKey, fb.ai.mapID = s.Def.MapsKey, s.Def.MapsMapID
		return fb
	}
	s.cur, s.at = st, time.Now()
	return st
}

func (s *Store) envSMS() string {
	if s.Def.SMSProvider == SMSLog && s.Def.Env == "production" {
		return SMSOff
	}
	if s.Def.SMSProvider == SMSLog {
		return SMSLog
	}
	return SMSOff
}

type querier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type row struct {
	enabled  bool
	config   map[string]string
	secret   []byte // terenkripsi
	hint     string
	verified *time.Time
	last     *TestResult
	exists   bool
}

func readRow(ctx context.Context, q querier, code string) (row, error) {
	var r row
	var cfg, last []byte
	var hint *string
	err := q.QueryRow(ctx, `SELECT enabled, config, secret_enc, secret_hint, verified_at, last_test FROM integrations WHERE code = $1`, code).
		Scan(&r.enabled, &cfg, &r.secret, &hint, &r.verified, &last)
	if errors.Is(err, pgx.ErrNoRows) {
		return row{config: map[string]string{}}, nil
	}
	if err != nil {
		return r, err
	}
	r.exists = true
	r.config = map[string]string{}
	_ = json.Unmarshal(cfg, &r.config)
	if hint != nil {
		r.hint = *hint
	}
	if len(last) > 0 {
		var t TestResult
		if json.Unmarshal(last, &t) == nil {
			r.last = &t
		}
	}
	return r, nil
}

func aad(code string) []byte { return []byte("integration:" + code) }

func (s *Store) open(r row, code string) ([]byte, error) {
	if r.secret == nil {
		return nil, nil
	}
	b, err := s.Box.Open(r.secret, aad(code))
	if err != nil {
		return nil, fmt.Errorf("rahasia %s tidak bisa dibuka (DATA_ENCRYPTION_KEY berubah?): %w", code, err)
	}
	return b, nil
}

func (s *Store) smsSetting(ctx context.Context, q querier) (string, bool, error) {
	var raw []byte
	err := q.QueryRow(ctx, `SELECT value FROM platform_settings WHERE key = 'sms'`).Scan(&raw)
	if errors.Is(err, pgx.ErrNoRows) {
		return s.envSMS(), false, nil
	}
	if err != nil {
		return "", false, err
	}
	var v struct {
		Provider string `json:"provider"`
	}
	_ = json.Unmarshal(raw, &v)
	return v.Provider, true, nil
}

func (s *Store) read(ctx context.Context) (*state, error) {
	st := &state{}
	mode, _, err := s.smsSetting(ctx, s.Pool)
	if err != nil {
		return nil, err
	}
	st.smsMode = mode
	if mode == SMSPlasGate {
		r, err := readRow(ctx, s.Pool, "plasgate")
		if err != nil {
			return nil, err
		}
		raw, err := s.open(r, "plasgate")
		if err != nil {
			return nil, err
		}
		var sec struct{ PrivateKey, Secret string }
		_ = json.Unmarshal(raw, &sec)
		st.pg = plasgate{PrivateKey: sec.PrivateKey, Secret: sec.Secret, Sender: r.config["sender"]}
		if st.pg.PrivateKey == "" || st.pg.Secret == "" || st.pg.Sender == "" {
			s.Log.Error("SMS PlasGate dipilih tapi kredensial belum lengkap — login SMS dimatikan")
			st.smsMode = SMSOff
		}
	}
	if mode == SMSLog && s.Def.Env == "production" {
		st.smsMode = SMSOff
	}
	if mode != SMSPlasGate && mode != SMSLog {
		st.smsMode = SMSOff
	}
	r, err := readRow(ctx, s.Pool, "facebook")
	if err != nil {
		return nil, err
	}
	if r.enabled && r.config["app_id"] != "" && r.secret != nil {
		key, err := s.open(r, "facebook")
		if err != nil {
			return nil, err
		}
		st.fbOn, st.fbID, st.fbKey = true, r.config["app_id"], string(key)
	}
	ai, err := s.readAI(ctx, s.Pool)
	if err != nil {
		return nil, err
	}
	st.ai = ai
	if st.google, st.apple, err = s.readOAuth(ctx, s.Pool); err != nil {
		return nil, err
	}
	return st, nil
}

// ---- untuk paket auth ----

// SMSEnabled: login dengan kode SMS tersedia.
func (s *Store) SMSEnabled(ctx context.Context) bool { return s.load(ctx).smsMode != SMSOff }

// Send mengirim SMS sesuai provider yang aktif.
func (s *Store) Send(ctx context.Context, phone, text string) error {
	st := s.load(ctx)
	switch st.smsMode {
	case SMSLog:
		s.Log.Warn("SMS (mode log, tidak benar-benar dikirim)", "to", phone, "text", text)
		return nil
	case SMSPlasGate:
		_, err := s.sendPlasGate(ctx, st.pg, phone, text)
		s.RecordUsage(ctx, Usage{Provider: "plasgate", Purpose: "sms", Units: smsParts(text), Err: err})
		return err
	}
	return ErrSMSOff
}

// sendPlasGate: POST /rest/send?private_key=… dengan header X-Secret. Nomor tanpa tanda + (855…).
func (s *Store) sendPlasGate(ctx context.Context, pg plasgate, phone, text string) (string, error) {
	body, _ := json.Marshal(map[string]string{"sender": pg.Sender, "to": strings.TrimPrefix(phone, "+"), "content": text})
	u := s.url(s.PlasGateURL, Endpoints.PlasGate) + "?private_key=" + url.QueryEscape(pg.PrivateKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Secret", pg.Secret)
	res, err := s.client().Do(req)
	if err != nil {
		return "", fmt.Errorf("plasgate: tidak bisa menghubungi server: %w", redact(err, pg.PrivateKey))
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(res.Body, 64<<10))
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("plasgate: HTTP %d: %s", res.StatusCode, trunc(string(raw)))
	}
	var out struct {
		QueueID any `json:"queue_id"`
	}
	_ = json.Unmarshal(raw, &out)
	return fmt.Sprint(out.QueueID), nil
}

// redact: jangan sampai private key (ada di URL) ikut tercatat di log.
func redact(err error, secret string) error {
	if secret == "" {
		return err
	}
	return errors.New(strings.ReplaceAll(err.Error(), url.QueryEscape(secret), "***"))
}

// smsParts: jumlah SMS yang ditagih (160 karakter GSM / 70 Unicode per bagian; pesan panjang 153 / 67).
func smsParts(text string) int {
	unicode := false
	for _, r := range text {
		if r > 127 {
			unicode = true
			break
		}
	}
	n, one, multi := len([]rune(text)), 160, 153
	if unicode {
		one, multi = 70, 67
	}
	if n <= one {
		return 1
	}
	return (n + multi - 1) / multi
}

func trunc(s string) string {
	if len(s) > 300 {
		return s[:300] + "…"
	}
	return s
}

// Facebook: App ID & App Secret kalau Login Facebook aktif.
func (s *Store) Facebook(ctx context.Context) (appID, secret string, ok bool) {
	st := s.load(ctx)
	return st.fbID, st.fbKey, st.fbOn
}

// FacebookURLs: alamat dialog & Graph API (bisa diganti di tes).
func (s *Store) FacebookURLs() (dialog, graph string) {
	return s.url(s.DialogURL, Endpoints.Dialog), s.url(s.GraphURL, Endpoints.Graph)
}

// RedirectURIs: alamat yang harus didaftarkan di Facebook Developers → Facebook Login → Valid OAuth Redirect URIs.
func (s *Store) RedirectURIs() []string {
	out := s.GoogleRedirectURIs() // sama persis, hanya /google/ → /facebook/
	for i, u := range out {
		out[i] = strings.Replace(u, "/auth/google/callback", "/auth/facebook/callback", 1)
	}
	return out
}

// ---- tampilan & perubahan oleh super admin ----

type TestResult struct {
	OK      bool      `json:"ok"`
	Message string    `json:"message"`
	At      time.Time `json:"at"`
}

type View struct {
	SMS struct {
		Mode       string      `json:"mode"`   // tersimpan
		Active     string      `json:"active"` // yang berjalan
		Source     string      `json:"source"` // admin · env
		EnvMode    string      `json:"env_mode"`
		Sender     string      `json:"sender"`
		KeySet     bool        `json:"key_set"`
		KeyHint    string      `json:"key_hint"`
		VerifiedAt *time.Time  `json:"verified_at"`
		LastTest   *TestResult `json:"last_test"`
	} `json:"sms"`
	Facebook struct {
		Enabled      bool        `json:"enabled"`
		Active       bool        `json:"active"`
		AppID        string      `json:"app_id"`
		KeySet       bool        `json:"key_set"`
		KeyHint      string      `json:"key_hint"`
		RedirectURIs []string    `json:"redirect_uris"`
		VerifiedAt   *time.Time  `json:"verified_at"`
		LastTest     *TestResult `json:"last_test"`
	} `json:"facebook"`
	Google   GoogleView `json:"google"`
	Apple    AppleView  `json:"apple"`
	AI       AIView     `json:"ai"`
	Maps     MapsView   `json:"maps"`
	Telegram struct {
		BotSet   bool   `json:"bot_set"` // diatur di .env (TELEGRAM_BOT_TOKEN), perlu restart
		Username string `json:"username"`
	} `json:"telegram"`
	Production bool `json:"production"`
}

func (s *Store) View(ctx context.Context) (View, error) {
	var v View
	mode, fromDB, err := s.smsSetting(ctx, s.Pool)
	if err != nil {
		return v, err
	}
	pg, err := readRow(ctx, s.Pool, "plasgate")
	if err != nil {
		return v, err
	}
	fb, err := readRow(ctx, s.Pool, "facebook")
	if err != nil {
		return v, err
	}
	s.Invalidate()
	st := s.load(ctx)
	v.SMS.Mode, v.SMS.Active, v.SMS.EnvMode = mode, st.smsMode, s.Def.SMSProvider
	v.SMS.Source = map[bool]string{true: "admin", false: "env"}[fromDB]
	v.SMS.Sender, v.SMS.KeySet, v.SMS.KeyHint, v.SMS.VerifiedAt, v.SMS.LastTest = pg.config["sender"], pg.hint != "", pg.hint, pg.verified, pg.last
	v.Facebook.Enabled, v.Facebook.Active, v.Facebook.AppID = fb.enabled, st.fbOn, fb.config["app_id"]
	v.Facebook.KeySet, v.Facebook.KeyHint, v.Facebook.VerifiedAt, v.Facebook.LastTest = fb.hint != "", fb.hint, fb.verified, fb.last
	v.Facebook.RedirectURIs = s.RedirectURIs()
	v.Production = s.Def.Env == "production"
	if v.AI, v.Maps, err = s.aiView(ctx, st); err != nil {
		return v, err
	}
	if v.Google, v.Apple, err = s.oauthView(ctx, st); err != nil {
		return v, err
	}
	v.Telegram.BotSet, v.Telegram.Username = s.Def.TelegramSet, s.Def.TelegramUsername
	return v, nil
}

func hintOf(k string) string {
	if len(k) <= 4 {
		return "••••"
	}
	return k[len(k)-4:]
}

type SMSInput struct {
	Mode       string  `json:"mode"` // off · plasgate · log · "" (ikut .env)
	Sender     string  `json:"sender"`
	PrivateKey *string `json:"private_key"` // kosong = tidak diubah
	Secret     *string `json:"secret"`
}

// SaveSMS menyimpan provider + kredensial PlasGate. Kredensial baru wajib berpasangan (private key + secret).
func (s *Store) SaveSMS(ctx context.Context, actorID string, in SMSInput) error {
	defer s.Invalidate()
	in.Sender = strings.TrimSpace(in.Sender)
	switch in.Mode {
	case "", SMSOff, SMSPlasGate:
	case SMSLog:
		if s.Def.Env == "production" {
			return FieldError{"mode", "log mode is not allowed in production"}
		}
	default:
		return FieldError{"mode", "unknown provider"}
	}
	pk, sec := strings.TrimSpace(deref(in.PrivateKey)), strings.TrimSpace(deref(in.Secret))
	if (pk == "") != (sec == "") {
		return FieldError{map[bool]string{true: "secret", false: "private_key"}[pk != ""], "enter the private key and secret together (same API key pair)"}
	}
	if pk != "" && !keyCharsRe.MatchString(pk) {
		return FieldError{"private_key", "looks wrong (8–256 characters, no spaces)"}
	}
	if sec != "" && !keyCharsRe.MatchString(sec) {
		return FieldError{"secret", "looks wrong (8–256 characters, no spaces)"}
	}
	if in.Sender != "" && !senderRe.MatchString(in.Sender) {
		return FieldError{"sender", "2–11 letters/numbers (registered sender ID) or a phone number"}
	}
	return db.WithActor(ctx, s.Pool, actorID, "admin", "", func(tx pgx.Tx) error {
		touchPG := pk != "" || in.Sender != ""
		if touchPG {
			var sealed []byte
			var hint string
			if pk != "" {
				b, _ := json.Marshal(map[string]string{"PrivateKey": pk, "Secret": sec})
				var err error
				if sealed, err = s.Box.Seal(b, aad("plasgate")); err != nil {
					return err
				}
				hint = hintOf(sec)
			}
			cfg, _ := json.Marshal(map[string]string{"sender": in.Sender})
			if _, err := tx.Exec(ctx, `INSERT INTO integrations (code, enabled, config, secret_enc, secret_hint, updated_by)
				VALUES ('plasgate', true, $1, $2, nullif($3, ''), $4)
				ON CONFLICT (code) DO UPDATE SET config = EXCLUDED.config,
					secret_enc = coalesce(EXCLUDED.secret_enc, integrations.secret_enc),
					secret_hint = coalesce(EXCLUDED.secret_hint, integrations.secret_hint),
					verified_at = CASE WHEN EXCLUDED.secret_enc IS NOT NULL OR integrations.config <> EXCLUDED.config THEN NULL ELSE integrations.verified_at END,
					updated_by = EXCLUDED.updated_by`, cfg, sealed, hint, actorID); err != nil {
				return err
			}
		}
		if in.Mode == "" {
			_, err := tx.Exec(ctx, `DELETE FROM platform_settings WHERE key = 'sms'`)
			return err
		}
		if in.Mode == SMSPlasGate {
			r, err := readRow(ctx, tx, "plasgate")
			if err != nil {
				return err
			}
			if r.secret == nil || r.config["sender"] == "" {
				return FieldError{"mode", "enter the PlasGate private key, secret and sender ID first"}
			}
		}
		v, _ := json.Marshal(map[string]string{"provider": in.Mode})
		_, err := tx.Exec(ctx, `INSERT INTO platform_settings (key, value, updated_by, updated_at) VALUES ('sms', $1, $2, now())
			ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_by = EXCLUDED.updated_by, updated_at = now()`, v, actorID)
		return err
	})
}

func deref(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

// TestSMS mengirim SMS uji ke nomor admin sendiri lewat PlasGate (kredensial tersimpan, walau belum aktif).
func (s *Store) TestSMS(ctx context.Context, actorID, phone string) (TestResult, error) {
	r, err := readRow(ctx, s.Pool, "plasgate")
	if err != nil {
		return TestResult{}, err
	}
	raw, err := s.open(r, "plasgate")
	if err != nil {
		return TestResult{}, err
	}
	var sec struct{ PrivateKey, Secret string }
	_ = json.Unmarshal(raw, &sec)
	if sec.PrivateKey == "" || r.config["sender"] == "" {
		return TestResult{}, FieldError{"private_key", "enter the PlasGate private key, secret and sender ID first"}
	}
	res := TestResult{At: time.Now().UTC()}
	tctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	id, err := s.sendPlasGate(tctx, plasgate{PrivateKey: sec.PrivateKey, Secret: sec.Secret, Sender: r.config["sender"]}, phone, "Phsar: SMS test OK")
	if err != nil {
		res.Message = err.Error()
		if strings.Contains(res.Message, "HTTP 401") {
			res.Message = "Wrong private key or secret (they must come from the same API key pair). " + res.Message
		}
	} else {
		res.OK, res.Message = true, "Test SMS sent to "+phone+" (queue "+id+")."
	}
	return res, s.saveTest(ctx, actorID, "plasgate", res)
}

func (s *Store) saveTest(ctx context.Context, actorID, code string, res TestResult) error {
	raw, _ := json.Marshal(res)
	return db.WithActor(ctx, s.Pool, actorID, "admin", "", func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `UPDATE integrations SET last_test = $2, verified_at = CASE WHEN $3 THEN now() ELSE NULL END WHERE code = $1`, code, raw, res.OK)
		return err
	})
}

type FacebookInput struct {
	Enabled   bool    `json:"enabled"`
	AppID     string  `json:"app_id"`
	AppSecret *string `json:"app_secret"` // kosong = tidak diubah
}

func (s *Store) SaveFacebook(ctx context.Context, actorID string, in FacebookInput) error {
	defer s.Invalidate()
	in.AppID = strings.TrimSpace(in.AppID)
	if in.AppID != "" && !appIDRe.MatchString(in.AppID) {
		return FieldError{"app_id", "App ID is a number from Facebook Developers"}
	}
	key := strings.TrimSpace(deref(in.AppSecret))
	if key != "" && !regexp.MustCompile(`^[0-9a-fA-F]{16,64}$`).MatchString(key) {
		return FieldError{"app_secret", "App Secret looks wrong (letters a–f and numbers)"}
	}
	var sealed []byte
	var hint string
	if key != "" {
		var err error
		if sealed, err = s.Box.Seal([]byte(key), aad("facebook")); err != nil {
			return err
		}
		hint = hintOf(key)
	}
	return db.WithActor(ctx, s.Pool, actorID, "admin", "", func(tx pgx.Tx) error {
		cfg, _ := json.Marshal(map[string]string{"app_id": in.AppID})
		if _, err := tx.Exec(ctx, `INSERT INTO integrations (code, enabled, config, secret_enc, secret_hint, updated_by)
			VALUES ('facebook', $1, $2, $3, nullif($4, ''), $5)
			ON CONFLICT (code) DO UPDATE SET enabled = EXCLUDED.enabled, config = EXCLUDED.config,
				secret_enc = coalesce(EXCLUDED.secret_enc, integrations.secret_enc),
				secret_hint = coalesce(EXCLUDED.secret_hint, integrations.secret_hint),
				verified_at = CASE WHEN EXCLUDED.secret_enc IS NOT NULL OR integrations.config <> EXCLUDED.config THEN NULL ELSE integrations.verified_at END,
				updated_by = EXCLUDED.updated_by`, in.Enabled, cfg, sealed, hint, actorID); err != nil {
			return err
		}
		if in.Enabled {
			r, err := readRow(ctx, tx, "facebook")
			if err != nil {
				return err
			}
			if r.config["app_id"] == "" || r.secret == nil {
				return FieldError{"app_id", "enter the App ID and App Secret before turning Facebook login on"}
			}
		}
		return nil
	})
}

// TestFacebook: minta app access token (client_credentials) — berhasil berarti App ID & Secret cocok.
func (s *Store) TestFacebook(ctx context.Context, actorID string) (TestResult, error) {
	r, err := readRow(ctx, s.Pool, "facebook")
	if err != nil {
		return TestResult{}, err
	}
	key, err := s.open(r, "facebook")
	if err != nil {
		return TestResult{}, err
	}
	if r.config["app_id"] == "" || len(key) == 0 {
		return TestResult{}, FieldError{"app_id", "enter the App ID and App Secret first"}
	}
	res := TestResult{At: time.Now().UTC()}
	_, graph := s.FacebookURLs()
	q := url.Values{"client_id": {r.config["app_id"]}, "client_secret": {string(key)}, "grant_type": {"client_credentials"}}
	tctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(tctx, http.MethodGet, graph+"/oauth/access_token?"+q.Encode(), nil)
	resp, err := s.client().Do(req)
	if err != nil {
		res.Message = "Cannot reach Facebook from this server: " + redact(err, string(key)).Error()
	} else {
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
		resp.Body.Close()
		var out struct {
			AccessToken string `json:"access_token"`
			Error       struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		_ = json.Unmarshal(raw, &out)
		if resp.StatusCode == http.StatusOK && out.AccessToken != "" {
			res.OK, res.Message = true, "Facebook accepted the App ID and App Secret."
		} else {
			res.Message = "Facebook rejected the App ID / App Secret: " + trunc(out.Error.Message)
		}
	}
	return res, s.saveTest(ctx, actorID, "facebook", res)
}
