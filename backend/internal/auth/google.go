package auth

// Login Google & Sign in with Apple.
//   Web (Google): GET <auth>/google/start?return=… → halaman izin Google → GET <auth>/google/callback?code&state
//        → kode ditukar (dengan Client Secret) jadi id_token → sesi → kembali ke halaman asal. State di cookie
//        HttpOnly (10 menit, sekali pakai), seperti Login Facebook.
//   App  : POST /v1/app/auth/google {id_token} (google_sign_in) · POST /v1/app/auth/apple {id_token, name}
//          (Sign in with Apple native) → tanda tangan diperiksa dengan kunci publik (JWKS) Google / Apple,
//          audience harus Client ID / Bundle ID Phsar.
// Yang disimpan: sub (id unik dari Google/Apple) + nama; email hanya sebagai tampilan identitas.

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"

	"phsar/internal/i18n"
	"phsar/internal/platform/db"
)

var (
	ErrGoogleDisabled = errors.New("login Google belum diaktifkan")
	ErrGoogleFailed   = errors.New("login Google gagal, coba lagi")
	ErrAppleDisabled  = errors.New("Sign in with Apple belum diaktifkan")
	ErrAppleFailed    = errors.New("Sign in with Apple gagal, coba lagi")
)

// GoogleConfig: kredensial Login Google (paket integ).
type GoogleConfig interface {
	Google(ctx context.Context) (clientID, secret string, audiences []string, ok bool)
	GoogleURLs() (auth, token, certs string)
}

// AppleConfig: Bundle ID aplikasi iOS yang diterima (paket integ).
type AppleConfig interface {
	Apple(ctx context.Context) (audiences []string, ok bool)
	AppleKeysURL() string
}

func (s *Service) UseGoogle(g GoogleConfig) { s.google = g }
func (s *Service) UseApple(a AppleConfig)   { s.apple = a }

// GoogleEnabled: tombol "Lanjut dengan Google" ditampilkan.
func (s *Service) GoogleEnabled(ctx context.Context) bool {
	if s.google == nil {
		return false
	}
	_, _, _, ok := s.google.Google(ctx)
	return ok
}

// AppleEnabled: aplikasi iPhone boleh memakai Sign in with Apple (tidak ada tombol di website).
func (s *Service) AppleEnabled(ctx context.Context) bool {
	if s.apple == nil {
		return false
	}
	_, ok := s.apple.Apple(ctx)
	return ok
}

const googleStateCookie = "phsar_gst"

// authBase: alamat publik bagian ini. Seller Center & admin bisa di subdomain sendiri (SELLER_BASE_URL /
// ADMIN_BASE_URL) — redirect URI OAuth harus memakai alamat itu, bukan domain utama.
func (s *Service) authBase(cookiePath string) string {
	switch {
	case cookiePath == "/seller/auth" && s.cfg.SellerBaseURL != "":
		return s.cfg.SellerBaseURL
	case cookiePath == "/admin/auth" && s.cfg.AdminBaseURL != "":
		return s.cfg.AdminBaseURL
	}
	return s.cfg.PublicBaseURL
}

func (s *Service) googleStart(w http.ResponseWriter, r *http.Request, cookiePath string) {
	ret := safeReturn(r.URL.Query().Get("return"), defaultReturn(cookiePath))
	if !s.GoogleEnabled(r.Context()) {
		http.Redirect(w, r, withQuery(ret, "auth_error", "google_disabled"), http.StatusFound)
		return
	}
	clientID, _, _, _ := s.google.Google(r.Context())
	authURL, _, _ := s.google.GoogleURLs()
	state, nonce := randomToken(16), randomToken(16)
	http.SetCookie(w, &http.Cookie{Name: googleStateCookie, Value: state + "|" + nonce + "|" + url.QueryEscape(ret), Path: cookiePath + "/google",
		MaxAge: 600, HttpOnly: true, SameSite: http.SameSiteLaxMode, Secure: s.secureCookies()})
	q := url.Values{"client_id": {clientID}, "redirect_uri": {s.authBase(cookiePath) + cookiePath + "/google/callback"},
		"response_type": {"code"}, "scope": {"openid email profile"}, "state": {state}, "nonce": {nonce},
		"prompt": {"select_account"}}
	if l := string(i18n.From(r.Context())); l != "" {
		q.Set("hl", l)
	}
	http.Redirect(w, r, authURL+"?"+q.Encode(), http.StatusFound)
}

func (s *Service) googleCallback(w http.ResponseWriter, r *http.Request, cookiePath, audience string, set func(http.ResponseWriter, TokenPair)) {
	def := defaultReturn(cookiePath)
	c, err := r.Cookie(googleStateCookie)
	http.SetCookie(w, &http.Cookie{Name: googleStateCookie, Value: "", Path: cookiePath + "/google", MaxAge: -1, HttpOnly: true, SameSite: http.SameSiteLaxMode})
	if err != nil {
		http.Redirect(w, r, withQuery(def, "auth_error", "google_expired"), http.StatusFound)
		return
	}
	parts := strings.SplitN(c.Value, "|", 3)
	if len(parts) != 3 {
		http.Redirect(w, r, withQuery(def, "auth_error", "google_expired"), http.StatusFound)
		return
	}
	state, nonce := parts[0], parts[1]
	ret, _ := url.QueryUnescape(parts[2])
	ret = safeReturn(ret, def)
	if got := r.URL.Query().Get("state"); state == "" || subtle.ConstantTimeCompare([]byte(state), []byte(got)) != 1 {
		http.Redirect(w, r, withQuery(ret, "auth_error", "google_expired"), http.StatusFound)
		return
	}
	if r.URL.Query().Get("error") != "" || r.URL.Query().Get("code") == "" { // pengguna menekan Batal
		http.Redirect(w, r, withQuery(ret, "auth_error", "google_cancelled"), http.StatusFound)
		return
	}
	claims, err := s.googleExchange(r.Context(), r.URL.Query().Get("code"), s.authBase(cookiePath)+cookiePath+"/google/callback")
	if err == nil && subtle.ConstantTimeCompare([]byte(claims.Nonce), []byte(nonce)) != 1 {
		err = errors.New("nonce tidak cocok")
	}
	if err != nil {
		s.log.Warn("google: login gagal", "err", err)
		http.Redirect(w, r, withQuery(ret, "auth_error", "google_failed"), http.StatusFound)
		return
	}
	pair, err := s.oidcLogin(r.Context(), "google", claims, "", audience, string(i18n.From(r.Context())), client(r, "web"))
	if err != nil {
		code := "google_failed"
		switch {
		case errors.Is(err, ErrNotAdmin):
			code = "not_admin"
		case errors.Is(err, ErrUserSuspended):
			code = "user_suspended"
		}
		s.log.Warn("google: login gagal", "err", err)
		http.Redirect(w, r, withQuery(ret, "auth_error", code), http.StatusFound)
		return
	}
	set(w, pair)
	http.Redirect(w, r, ret, http.StatusFound)
}

// googleExchange: tukar kode → id_token (diterima langsung dari Google lewat TLS, tetap diperiksa penuh).
func (s *Service) googleExchange(ctx context.Context, code, redirect string) (idClaims, error) {
	clientID, secret, auds, _ := s.google.Google(ctx)
	_, tokenURL, certs := s.google.GoogleURLs()
	form := url.Values{"client_id": {clientID}, "client_secret": {secret}, "code": {code}, "grant_type": {"authorization_code"}, "redirect_uri": {redirect}}
	ctx2, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx2, http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return idClaims{}, errors.New("tidak bisa menghubungi Google")
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(res.Body, 64<<10))
	var out struct {
		IDToken string `json:"id_token"`
		Error   string `json:"error"`
	}
	_ = json.Unmarshal(raw, &out)
	if res.StatusCode != http.StatusOK || out.IDToken == "" {
		return idClaims{}, fmt.Errorf("token endpoint HTTP %d %s", res.StatusCode, out.Error)
	}
	return verifyIDToken(ctx, out.IDToken, certs, googleIssuers, auds, time.Now())
}

// handleGoogleToken: aplikasi (google_sign_in) mengirim id_token.
func (s *Service) handleGoogleToken(w http.ResponseWriter, r *http.Request, audience string) {
	var in struct {
		IDToken string `json:"id_token"`
	}
	if !decode(w, r, &in) {
		return
	}
	if !s.GoogleEnabled(r.Context()) {
		s.fail(w, ErrGoogleDisabled)
		return
	}
	_, _, auds, _ := s.google.Google(r.Context())
	_, _, certs := s.google.GoogleURLs()
	claims, err := verifyIDToken(r.Context(), in.IDToken, certs, googleIssuers, auds, time.Now())
	if err != nil {
		s.log.Warn("google: id_token aplikasi ditolak", "err", err)
		s.fail(w, ErrGoogleFailed)
		return
	}
	pair, err := s.oidcLogin(r.Context(), "google", claims, "", audience, string(i18n.From(r.Context())), client(r, r.Header.Get("X-Device-Name")))
	if err != nil {
		s.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, pair)
}

// handleAppleToken: aplikasi iPhone (Sign in with Apple) mengirim id_token. Nama hanya dikirim Apple saat login
// pertama (dari aplikasi, bukan di dalam token), jadi aplikasi meneruskannya sebagai "name".
func (s *Service) handleAppleToken(w http.ResponseWriter, r *http.Request, audience string) {
	var in struct {
		IDToken string `json:"id_token"`
		Name    string `json:"name"`
	}
	if !decode(w, r, &in) {
		return
	}
	if !s.AppleEnabled(r.Context()) {
		s.fail(w, ErrAppleDisabled)
		return
	}
	auds, _ := s.apple.Apple(r.Context())
	claims, err := verifyIDToken(r.Context(), in.IDToken, s.apple.AppleKeysURL(), appleIssuers, auds, time.Now())
	if err != nil {
		s.log.Warn("apple: id_token aplikasi ditolak", "err", err)
		s.fail(w, ErrAppleFailed)
		return
	}
	pair, err := s.oidcLogin(r.Context(), "apple", claims, in.Name, audience, string(i18n.From(r.Context())), client(r, r.Header.Get("X-Device-Name")))
	if err != nil {
		s.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, pair)
}

func (s *Service) oidcLogin(ctx context.Context, provider string, c idClaims, name, audience, lang string, cl ClientInfo) (TokenPair, error) {
	if name == "" {
		name = c.Name
	}
	display := c.Email
	if display == "" {
		display = truncName(name)
	}
	var pair TokenPair
	notAdmin := false
	err := db.WithActor(ctx, s.pool, "", "system", "", func(tx pgx.Tx) error {
		userID, created, err := findOrCreateUser(ctx, tx, identity{provider: provider, uid: c.Sub, name: truncName(name), display: display, lang: langFrom(lang)})
		if err != nil {
			return err
		}
		if ok, err := adminAllowed(ctx, tx, userID, audience); err != nil || !ok {
			notAdmin = err == nil
			return err
		}
		pair, err = s.issueSession(ctx, tx, userID, audience, created, cl)
		return err
	})
	s.recordAttempt(ctx, provider+":"+c.Sub, cl, err == nil && !notAdmin, map[bool]string{true: "not_admin"}[notAdmin])
	if err == nil && notAdmin {
		return TokenPair{}, ErrNotAdmin
	}
	return pair, err
}

// ---- id_token (JWT RS256) ----

var (
	googleIssuers = []string{"https://accounts.google.com", "accounts.google.com"}
	appleIssuers  = []string{"https://appleid.apple.com"}
)

type idClaims struct {
	Iss   string `json:"iss"`
	Sub   string `json:"sub"`
	Aud   any    `json:"aud"` // string atau []string
	Exp   int64  `json:"exp"`
	Iat   int64  `json:"iat"`
	Nonce string `json:"nonce"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (c idClaims) audiences() []string {
	switch v := c.Aud.(type) {
	case string:
		return []string{v}
	case []any:
		var out []string
		for _, x := range v {
			if s, ok := x.(string); ok {
				out = append(out, s)
			}
		}
		return out
	}
	return nil
}

func verifyIDToken(ctx context.Context, token, jwksURL string, issuers, audiences []string, now time.Time) (idClaims, error) {
	var c idClaims
	parts := strings.Split(token, ".")
	if len(token) > 8192 || len(parts) != 3 {
		return c, errors.New("id_token bukan JWT")
	}
	b64 := base64.RawURLEncoding
	hb, err1 := b64.DecodeString(parts[0])
	pb, err2 := b64.DecodeString(parts[1])
	sig, err3 := b64.DecodeString(parts[2])
	if err1 != nil || err2 != nil || err3 != nil {
		return c, errors.New("id_token rusak")
	}
	var head struct {
		Alg, Kid string
	}
	if json.Unmarshal(hb, &head) != nil || head.Alg != "RS256" || head.Kid == "" {
		return c, errors.New("id_token: algoritma tidak didukung")
	}
	key, err := jwks.key(ctx, jwksURL, head.Kid)
	if err != nil {
		return c, err
	}
	sum := sha256.Sum256([]byte(parts[0] + "." + parts[1]))
	if rsa.VerifyPKCS1v15(key, crypto.SHA256, sum[:], sig) != nil {
		return c, errors.New("id_token: tanda tangan tidak sah")
	}
	if json.Unmarshal(pb, &c) != nil {
		return c, errors.New("id_token: isi rusak")
	}
	if !slices.Contains(issuers, c.Iss) {
		return c, fmt.Errorf("id_token: penerbit %q tidak dikenal", c.Iss)
	}
	ok := false
	for _, a := range c.audiences() {
		if a != "" && slices.Contains(audiences, a) {
			ok = true
		}
	}
	if !ok {
		return c, errors.New("id_token: bukan untuk aplikasi Phsar")
	}
	if c.Exp == 0 || now.Unix() > c.Exp+60 {
		return c, errors.New("id_token: kedaluwarsa")
	}
	if c.Iat > now.Unix()+300 {
		return c, errors.New("id_token: waktu terbit di masa depan")
	}
	if c.Sub == "" || len(c.Sub) > 255 {
		return c, errors.New("id_token: tanpa sub")
	}
	return c, nil
}

// jwks: kunci publik Google / Apple, disimpan 1 jam; kid baru (rotasi kunci) → ambil ulang segera (maks. 1×/menit).
var jwks = &jwksCache{sets: map[string]*jwkSet{}}

type jwkSet struct {
	keys    map[string]*rsa.PublicKey
	fetched time.Time
}

type jwksCache struct {
	mu   sync.Mutex
	sets map[string]*jwkSet
}

func (j *jwksCache) key(ctx context.Context, u, kid string) (*rsa.PublicKey, error) {
	j.mu.Lock()
	defer j.mu.Unlock()
	set := j.sets[u]
	if set != nil && time.Since(set.fetched) < time.Hour {
		if k := set.keys[kid]; k != nil {
			return k, nil
		}
		if time.Since(set.fetched) < time.Minute {
			return nil, errors.New("id_token: kunci tidak dikenal")
		}
	}
	fresh, err := fetchJWKS(ctx, u)
	if err != nil {
		if set != nil && set.keys[kid] != nil { // jaringan gagal: pakai yang lama
			return set.keys[kid], nil
		}
		return nil, err
	}
	j.sets[u] = fresh
	if k := fresh.keys[kid]; k != nil {
		return k, nil
	}
	return nil, errors.New("id_token: kunci tidak dikenal")
}

func fetchJWKS(ctx context.Context, u string) (*jwkSet, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.New("tidak bisa mengambil kunci publik login")
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("kunci publik login: HTTP %d", res.StatusCode)
	}
	var doc struct {
		Keys []struct {
			Kty, Kid, N, E, Use, Alg string
		} `json:"keys"`
	}
	if err := json.NewDecoder(io.LimitReader(res.Body, 256<<10)).Decode(&doc); err != nil {
		return nil, err
	}
	set := &jwkSet{keys: map[string]*rsa.PublicKey{}, fetched: time.Now()}
	for _, k := range doc.Keys {
		if k.Kty != "RSA" || k.Kid == "" || (k.Use != "" && k.Use != "sig") {
			continue
		}
		n, err1 := base64.RawURLEncoding.DecodeString(k.N)
		e, err2 := base64.RawURLEncoding.DecodeString(k.E)
		if err1 != nil || err2 != nil || len(e) == 0 || len(e) > 4 {
			continue
		}
		set.keys[k.Kid] = &rsa.PublicKey{N: new(big.Int).SetBytes(n), E: int(new(big.Int).SetBytes(e).Int64())}
	}
	if len(set.keys) == 0 {
		return nil, errors.New("kunci publik login kosong")
	}
	return set, nil
}
