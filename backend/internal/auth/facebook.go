package auth

// Login Facebook.
//   Web: GET <auth>/facebook/start?return=/km/checkout → dialog Facebook → GET <auth>/facebook/callback?code&state
//        → cookie sesi dipasang → kembali ke halaman asal. State acak disimpan di cookie HttpOnly (10 menit, sekali pakai).
//   App: POST /v1/app/auth/facebook {access_token} dari Facebook SDK → diperiksa dengan debug_token (harus milik App ID kita).
// Profil yang dipakai hanya id + nama (izin public_profile). Kredensial (App ID / App Secret) diatur super admin.

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"phsar/internal/i18n"
	"phsar/internal/platform/db"
)

var (
	ErrFacebookDisabled = errors.New("login Facebook belum diaktifkan")
	ErrFacebookFailed   = errors.New("login Facebook gagal, coba lagi")
)

// FacebookConfig: App ID & App Secret yang diatur super admin (paket integ).
type FacebookConfig interface {
	Facebook(ctx context.Context) (appID, secret string, ok bool)
	FacebookURLs() (dialog, graph string)
}

// UseFacebook memasang sumber pengaturan Login Facebook.
func (s *Service) UseFacebook(f FacebookConfig) { s.fb = f }

// FacebookEnabled: tombol "Lanjut dengan Facebook" ditampilkan.
func (s *Service) FacebookEnabled(ctx context.Context) bool {
	if s.fb == nil {
		return false
	}
	_, _, ok := s.fb.Facebook(ctx)
	return ok
}

const fbStateCookie = "phsar_fbst"

// safeReturn: hanya path di situs sendiri (cegah open redirect ke situs lain).
func safeReturn(p, def string) string {
	if p == "" || !strings.HasPrefix(p, "/") || strings.HasPrefix(p, "//") || strings.HasPrefix(p, "/\\") || strings.ContainsAny(p, "\r\n") {
		return def
	}
	return p
}

// defaultReturn: halaman awal tiap bagian (/shop/auth → /, /seller/auth → /seller/, …).
func defaultReturn(cookiePath string) string {
	if cookiePath == "/shop/auth" {
		return "/"
	}
	return strings.TrimSuffix(cookiePath, "auth")
}

func withQuery(p, k, v string) string {
	sep := "?"
	if strings.Contains(p, "?") {
		sep = "&"
	}
	if i := strings.Index(p, "#"); i >= 0 { // parameter harus sebelum #
		return p[:i] + sep + k + "=" + url.QueryEscape(v) + p[i:]
	}
	return p + sep + k + "=" + url.QueryEscape(v)
}

func (s *Service) secureCookies() bool { return strings.HasPrefix(s.cfg.PublicBaseURL, "https://") }

func (s *Service) facebookStart(w http.ResponseWriter, r *http.Request, cookiePath string) {
	def := defaultReturn(cookiePath)
	ret := safeReturn(r.URL.Query().Get("return"), def)
	appID, _, ok := s.facebookCreds(r.Context())
	if !ok {
		http.Redirect(w, r, withQuery(ret, "auth_error", "facebook_disabled"), http.StatusFound)
		return
	}
	state := randomToken(16)
	http.SetCookie(w, &http.Cookie{Name: fbStateCookie, Value: state + "|" + url.QueryEscape(ret), Path: cookiePath + "/facebook",
		MaxAge: 600, HttpOnly: true, SameSite: http.SameSiteLaxMode, Secure: s.secureCookies()})
	dialog, _ := s.fb.FacebookURLs()
	q := url.Values{"client_id": {appID}, "redirect_uri": {s.authBase(cookiePath) + cookiePath + "/facebook/callback"},
		"state": {state}, "response_type": {"code"}, "scope": {"public_profile"}}
	http.Redirect(w, r, dialog+"?"+q.Encode(), http.StatusFound)
}

func (s *Service) facebookCallback(w http.ResponseWriter, r *http.Request, cookiePath, audience string, set func(http.ResponseWriter, TokenPair)) {
	def := defaultReturn(cookiePath)
	c, err := r.Cookie(fbStateCookie)
	http.SetCookie(w, &http.Cookie{Name: fbStateCookie, Value: "", Path: cookiePath + "/facebook", MaxAge: -1, HttpOnly: true, SameSite: http.SameSiteLaxMode})
	if err != nil {
		http.Redirect(w, r, withQuery(def, "auth_error", "facebook_expired"), http.StatusFound)
		return
	}
	state, retEnc, _ := strings.Cut(c.Value, "|")
	ret, _ := url.QueryUnescape(retEnc)
	ret = safeReturn(ret, def)
	got := r.URL.Query().Get("state")
	if state == "" || subtle.ConstantTimeCompare([]byte(state), []byte(got)) != 1 {
		http.Redirect(w, r, withQuery(ret, "auth_error", "facebook_expired"), http.StatusFound)
		return
	}
	if r.URL.Query().Get("error") != "" || r.URL.Query().Get("code") == "" { // pengguna menekan Batal
		http.Redirect(w, r, withQuery(ret, "auth_error", "facebook_cancelled"), http.StatusFound)
		return
	}
	cl := client(r, "web")
	token, err := s.facebookExchange(r.Context(), r.URL.Query().Get("code"), s.authBase(cookiePath)+cookiePath+"/facebook/callback")
	if err != nil {
		s.log.Warn("facebook: tukar kode gagal", "err", err)
		http.Redirect(w, r, withQuery(ret, "auth_error", "facebook_failed"), http.StatusFound)
		return
	}
	pair, err := s.facebookLogin(r.Context(), token, audience, string(i18n.From(r.Context())), cl)
	if err != nil {
		code := "facebook_failed"
		switch {
		case errors.Is(err, ErrNotAdmin):
			code = "not_admin"
		case errors.Is(err, ErrUserSuspended):
			code = "user_suspended"
		}
		s.log.Warn("facebook: login gagal", "err", err)
		http.Redirect(w, r, withQuery(ret, "auth_error", code), http.StatusFound)
		return
	}
	set(w, pair)
	http.Redirect(w, r, ret, http.StatusFound)
}

// handleFacebookToken: app (Facebook SDK) mengirim user access token.
func (s *Service) handleFacebookToken(w http.ResponseWriter, r *http.Request, audience string) {
	var in struct {
		AccessToken string `json:"access_token"`
	}
	if !decode(w, r, &in) {
		return
	}
	if _, _, ok := s.facebookCreds(r.Context()); !ok {
		s.fail(w, ErrFacebookDisabled)
		return
	}
	if err := s.facebookCheckToken(r.Context(), in.AccessToken); err != nil {
		s.log.Warn("facebook: token app ditolak", "err", err)
		s.fail(w, ErrFacebookFailed)
		return
	}
	pair, err := s.facebookLogin(r.Context(), in.AccessToken, audience, string(i18n.From(r.Context())), client(r, r.Header.Get("X-Device-Name")))
	if err != nil {
		s.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, pair)
}

func (s *Service) facebookCreds(ctx context.Context) (string, string, bool) {
	if s.fb == nil {
		return "", "", false
	}
	return s.fb.Facebook(ctx)
}

func (s *Service) graphGet(ctx context.Context, path string, q url.Values, out any) error {
	_, graph := s.fb.FacebookURLs()
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, graph+path+"?"+q.Encode(), nil)
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("graph %s: tidak bisa menghubungi Facebook", path)
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(res.Body, 64<<10))
	if res.StatusCode != http.StatusOK {
		var e struct {
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		_ = json.Unmarshal(raw, &e)
		return fmt.Errorf("graph %s: HTTP %d: %s", path, res.StatusCode, e.Error.Message)
	}
	return json.Unmarshal(raw, out)
}

func (s *Service) facebookExchange(ctx context.Context, code, redirect string) (string, error) {
	appID, secret, _ := s.facebookCreds(ctx)
	var out struct {
		AccessToken string `json:"access_token"`
	}
	err := s.graphGet(ctx, "/oauth/access_token", url.Values{"client_id": {appID}, "client_secret": {secret}, "redirect_uri": {redirect}, "code": {code}}, &out)
	if err == nil && out.AccessToken == "" {
		err = errors.New("tidak ada access_token")
	}
	return out.AccessToken, err
}

// facebookCheckToken: token dari app harus valid dan dibuat untuk App ID Phsar (bukan aplikasi lain).
func (s *Service) facebookCheckToken(ctx context.Context, token string) error {
	if token == "" || len(token) > 1024 {
		return errors.New("token kosong")
	}
	appID, secret, _ := s.facebookCreds(ctx)
	var out struct {
		Data struct {
			AppID   string `json:"app_id"`
			IsValid bool   `json:"is_valid"`
		} `json:"data"`
	}
	if err := s.graphGet(ctx, "/debug_token", url.Values{"input_token": {token}, "access_token": {appID + "|" + secret}}, &out); err != nil {
		return err
	}
	if !out.Data.IsValid || out.Data.AppID != appID {
		return errors.New("token bukan untuk aplikasi ini")
	}
	return nil
}

func (s *Service) facebookLogin(ctx context.Context, token, audience, lang string, c ClientInfo) (TokenPair, error) {
	_, secret, _ := s.facebookCreds(ctx)
	m := hmac.New(sha256.New, []byte(secret))
	m.Write([]byte(token))
	var me struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if err := s.graphGet(ctx, "/me", url.Values{"fields": {"id,name"}, "access_token": {token}, "appsecret_proof": {hex.EncodeToString(m.Sum(nil))}}, &me); err != nil {
		return TokenPair{}, err
	}
	if me.ID == "" {
		return TokenPair{}, ErrFacebookFailed
	}
	var pair TokenPair
	notAdmin := false
	err := db.WithActor(ctx, s.pool, "", "system", "", func(tx pgx.Tx) error {
		userID, created, err := findOrCreateUser(ctx, tx, identity{provider: "facebook", uid: me.ID, name: truncName(me.Name), lang: langFrom(lang)})
		if err != nil {
			return err
		}
		if ok, err := adminAllowed(ctx, tx, userID, audience); err != nil || !ok {
			notAdmin = err == nil
			return err
		}
		pair, err = s.issueSession(ctx, tx, userID, audience, created, c)
		return err
	})
	s.recordAttempt(ctx, "facebook:"+me.ID, c, err == nil && !notAdmin, map[bool]string{true: "not_admin"}[notAdmin])
	if err == nil && notAdmin {
		return TokenPair{}, ErrNotAdmin
	}
	return pair, err
}

func truncName(n string) string {
	r := []rune(strings.TrimSpace(n))
	if len(r) > 80 {
		r = r[:80]
	}
	return string(r)
}
