package integ

// Login Google & Sign in with Apple — kredensial diatur super admin (Admin → Login & SMS).
//   Google : OAuth Client ID + Client Secret (tipe "Web application") dari Google Cloud Console, ditambah
//            Client ID aplikasi Android/iOS (opsional) supaya id_token dari aplikasi Flutter diterima.
//   Apple  : hanya untuk aplikasi iPhone (Flutter, Sign in with Apple native) — cukup Bundle ID / Services ID
//            sebagai audience id_token. Tombol Apple di website sengaja tidak ada.
// Store memenuhi auth.GoogleConfig dan auth.AppleConfig.

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"phsar/internal/platform/db"
)

const (
	GoogleAuthURL  = "https://accounts.google.com/o/oauth2/v2/auth"
	GoogleTokenURL = "https://oauth2.googleapis.com/token"
	GoogleCertsURL = "https://www.googleapis.com/oauth2/v3/certs"
	AppleKeysURL   = "https://appleid.apple.com/auth/keys"
)

// OAuthEndpoints: alamat Google/Apple (hanya diganti tes).
var OAuthEndpoints = struct{ GoogleAuth, GoogleToken, GoogleCerts, AppleKeys string }{GoogleAuthURL, GoogleTokenURL, GoogleCertsURL, AppleKeysURL}

var (
	googleClientRe = regexp.MustCompile(`^[0-9]{6,20}-[a-z0-9]{10,64}\.apps\.googleusercontent\.com$`)
	bundleRe       = regexp.MustCompile(`^[A-Za-z0-9-]+(\.[A-Za-z0-9-]+){1,8}$`)
)

type googleState struct {
	on       bool
	clientID string
	secret   string
	appIDs   []string // client ID Android / iOS (audience id_token aplikasi)
}

type appleState struct {
	on        bool
	audiences []string // Bundle ID aplikasi iOS (+ Services ID kalau nanti dipakai)
}

func splitIDs(s string) []string {
	var out []string
	for _, p := range strings.FieldsFunc(s, func(r rune) bool { return r == ',' || r == ' ' || r == '\n' }) {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func (s *Store) readOAuth(ctx context.Context, q querier) (googleState, appleState, error) {
	var g googleState
	var a appleState
	r, err := readRow(ctx, q, "google")
	if err != nil {
		return g, a, err
	}
	if r.enabled && r.config["client_id"] != "" && r.secret != nil {
		key, err := s.open(r, "google")
		if err != nil {
			return g, a, err
		}
		g = googleState{on: true, clientID: r.config["client_id"], secret: string(key), appIDs: splitIDs(r.config["app_client_ids"])}
	}
	ar, err := readRow(ctx, q, "apple")
	if err != nil {
		return g, a, err
	}
	if ar.enabled && ar.config["bundle_ids"] != "" {
		a = appleState{on: true, audiences: splitIDs(ar.config["bundle_ids"])}
	}
	return g, a, nil
}

// ---- untuk paket auth ----

// Google: Client ID & Secret web, dan semua audience id_token yang diterima (web + aplikasi).
func (s *Store) Google(ctx context.Context) (clientID, secret string, audiences []string, ok bool) {
	g := s.load(ctx).google
	return g.clientID, g.secret, append([]string{g.clientID}, g.appIDs...), g.on
}

// GoogleURLs: halaman izin, token endpoint, kunci publik (JWKS).
func (s *Store) GoogleURLs() (auth, token, certs string) {
	return OAuthEndpoints.GoogleAuth, OAuthEndpoints.GoogleToken, OAuthEndpoints.GoogleCerts
}

// Apple: audience id_token yang diterima (Bundle ID aplikasi iOS).
func (s *Store) Apple(ctx context.Context) (audiences []string, ok bool) {
	a := s.load(ctx).apple
	return a.audiences, a.on
}

func (s *Store) AppleKeysURL() string { return OAuthEndpoints.AppleKeys }

// GoogleRedirectURIs: alamat yang didaftarkan di Google Cloud → Credentials → OAuth client → Authorized redirect URIs.
// Seller Center & admin memakai subdomain sendiri kalau SELLER_BASE_URL / ADMIN_BASE_URL diisi.
func (s *Store) GoogleRedirectURIs() []string {
	base := strings.TrimRight(s.Def.BaseURL, "/")
	seller, admin := base, base
	if s.Def.SellerURL != "" {
		seller = strings.TrimRight(s.Def.SellerURL, "/")
	}
	if s.Def.AdminURL != "" {
		admin = strings.TrimRight(s.Def.AdminURL, "/")
	}
	return []string{base + "/shop/auth/google/callback", seller + "/seller/auth/google/callback",
		base + "/courier/auth/google/callback", admin + "/admin/auth/google/callback"}
}

// GoogleOrigins: Authorized JavaScript origins (tidak wajib untuk alur redirect, tapi Google Console memintanya).
func (s *Store) GoogleOrigins() []string {
	out := []string{strings.TrimRight(s.Def.BaseURL, "/")}
	for _, u := range []string{s.Def.SellerURL, s.Def.AdminURL} {
		if u = strings.TrimRight(u, "/"); u != "" && u != out[0] {
			out = append(out, u)
		}
	}
	return out
}

// ---- tampilan & perubahan oleh super admin ----

type GoogleView struct {
	Enabled      bool        `json:"enabled"`
	Active       bool        `json:"active"`
	ClientID     string      `json:"client_id"`
	AppClientIDs string      `json:"app_client_ids"`
	KeySet       bool        `json:"key_set"`
	KeyHint      string      `json:"key_hint"`
	RedirectURIs []string    `json:"redirect_uris"`
	Origins      []string    `json:"origins"`
	VerifiedAt   *time.Time  `json:"verified_at"`
	LastTest     *TestResult `json:"last_test"`
}

type AppleView struct {
	Enabled   bool   `json:"enabled"`
	Active    bool   `json:"active"`
	BundleIDs string `json:"bundle_ids"`
}

func (s *Store) oauthView(ctx context.Context, st *state) (GoogleView, AppleView, error) {
	var g GoogleView
	var a AppleView
	r, err := readRow(ctx, s.Pool, "google")
	if err != nil {
		return g, a, err
	}
	g = GoogleView{Enabled: r.enabled, Active: st.google.on, ClientID: r.config["client_id"], AppClientIDs: r.config["app_client_ids"],
		KeySet: r.hint != "", KeyHint: r.hint, RedirectURIs: s.GoogleRedirectURIs(), Origins: s.GoogleOrigins(), VerifiedAt: r.verified, LastTest: r.last}
	ar, err := readRow(ctx, s.Pool, "apple")
	if err != nil {
		return g, a, err
	}
	a = AppleView{Enabled: ar.enabled, Active: st.apple.on, BundleIDs: ar.config["bundle_ids"]}
	return g, a, nil
}

type GoogleInput struct {
	Enabled      bool    `json:"enabled"`
	ClientID     string  `json:"client_id"`
	AppClientIDs string  `json:"app_client_ids"`
	ClientSecret *string `json:"client_secret"` // kosong = tidak diubah
}

func (s *Store) SaveGoogle(ctx context.Context, actorID string, in GoogleInput) error {
	defer s.Invalidate()
	in.ClientID = strings.TrimSpace(in.ClientID)
	if in.ClientID != "" && !googleClientRe.MatchString(in.ClientID) {
		return FieldError{"client_id", "Client ID looks like 1234567890-abc….apps.googleusercontent.com"}
	}
	apps := splitIDs(in.AppClientIDs)
	if len(apps) > 6 {
		return FieldError{"app_client_ids", "at most 6 app client IDs"}
	}
	for _, id := range apps {
		if !googleClientRe.MatchString(id) {
			return FieldError{"app_client_ids", "each app Client ID looks like 1234567890-abc….apps.googleusercontent.com"}
		}
	}
	key := strings.TrimSpace(deref(in.ClientSecret))
	if key != "" && !keyCharsRe.MatchString(key) {
		return FieldError{"client_secret", "Client Secret looks wrong"}
	}
	var sealed []byte
	var hint string
	if key != "" {
		var err error
		if sealed, err = s.Box.Seal([]byte(key), aad("google")); err != nil {
			return err
		}
		hint = hintOf(key)
	}
	return db.WithActor(ctx, s.Pool, actorID, "admin", "", func(tx pgx.Tx) error {
		cfg, _ := json.Marshal(map[string]string{"client_id": in.ClientID, "app_client_ids": strings.Join(apps, ",")})
		if _, err := tx.Exec(ctx, `INSERT INTO integrations (code, enabled, config, secret_enc, secret_hint, updated_by)
			VALUES ('google', $1, $2, $3, nullif($4, ''), $5)
			ON CONFLICT (code) DO UPDATE SET enabled = EXCLUDED.enabled, config = EXCLUDED.config,
				secret_enc = coalesce(EXCLUDED.secret_enc, integrations.secret_enc),
				secret_hint = coalesce(EXCLUDED.secret_hint, integrations.secret_hint),
				verified_at = CASE WHEN EXCLUDED.secret_enc IS NOT NULL OR integrations.config <> EXCLUDED.config THEN NULL ELSE integrations.verified_at END,
				updated_by = EXCLUDED.updated_by`, in.Enabled, cfg, sealed, hint, actorID); err != nil {
			return err
		}
		if in.Enabled {
			r, err := readRow(ctx, tx, "google")
			if err != nil {
				return err
			}
			if r.config["client_id"] == "" || r.secret == nil {
				return FieldError{"client_id", "enter the Client ID and Client Secret before turning Google login on"}
			}
		}
		return nil
	})
}

// TestGoogle: tukar kode palsu di token endpoint. Google menjawab "invalid_grant" kalau Client ID & Secret
// cocok (kodenya saja yang salah), dan "invalid_client" / "unauthorized_client" kalau kredensialnya salah.
func (s *Store) TestGoogle(ctx context.Context, actorID string) (TestResult, error) {
	r, err := readRow(ctx, s.Pool, "google")
	if err != nil {
		return TestResult{}, err
	}
	key, err := s.open(r, "google")
	if err != nil {
		return TestResult{}, err
	}
	if r.config["client_id"] == "" || len(key) == 0 {
		return TestResult{}, FieldError{"client_id", "enter the Client ID and Client Secret first"}
	}
	res := TestResult{At: time.Now().UTC()}
	redirect := s.GoogleRedirectURIs()[0]
	form := url.Values{"client_id": {r.config["client_id"]}, "client_secret": {string(key)}, "code": {"phsar-connection-test"},
		"grant_type": {"authorization_code"}, "redirect_uri": {redirect}}
	tctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(tctx, http.MethodPost, OAuthEndpoints.GoogleToken, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := s.client().Do(req)
	if err != nil {
		res.Message = "Cannot reach Google from this server: " + redact(err, string(key)).Error()
	} else {
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
		resp.Body.Close()
		var out struct {
			Error string `json:"error"`
			Desc  string `json:"error_description"`
		}
		_ = json.Unmarshal(raw, &out)
		switch out.Error {
		case "invalid_grant":
			res.OK, res.Message = true, "Google accepted the Client ID and Client Secret."
		case "redirect_uri_mismatch":
			res.Message = "Client ID & Secret are fine, but the redirect URIs are not registered yet. Add all addresses below in Google Cloud Console."
		default:
			res.Message = "Google rejected the Client ID / Client Secret: " + trunc(strings.TrimSpace(out.Error+" "+out.Desc))
		}
	}
	return res, s.saveTest(ctx, actorID, "google", res)
}

type AppleInput struct {
	Enabled   bool   `json:"enabled"`
	BundleIDs string `json:"bundle_ids"`
}

// SaveApple: Bundle ID aplikasi iOS. Tidak ada rahasia — id_token Apple diperiksa dengan kunci publik Apple.
func (s *Store) SaveApple(ctx context.Context, actorID string, in AppleInput) error {
	defer s.Invalidate()
	ids := splitIDs(in.BundleIDs)
	if len(ids) > 4 {
		return FieldError{"bundle_ids", "at most 4 IDs"}
	}
	for _, id := range ids {
		if !bundleRe.MatchString(id) {
			return FieldError{"bundle_ids", "Bundle ID looks like com.phsar.app"}
		}
	}
	if in.Enabled && len(ids) == 0 {
		return FieldError{"bundle_ids", "enter the app Bundle ID before turning Sign in with Apple on"}
	}
	return db.WithActor(ctx, s.Pool, actorID, "admin", "", func(tx pgx.Tx) error {
		cfg, _ := json.Marshal(map[string]string{"bundle_ids": strings.Join(ids, ",")})
		_, err := tx.Exec(ctx, `INSERT INTO integrations (code, enabled, config, updated_by) VALUES ('apple', $1, $2, $3)
			ON CONFLICT (code) DO UPDATE SET enabled = EXCLUDED.enabled, config = EXCLUDED.config, updated_by = EXCLUDED.updated_by`,
			in.Enabled, cfg, actorID)
		return err
	})
}
