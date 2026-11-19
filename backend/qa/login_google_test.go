package qa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"log/slog"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"phsar/internal/auth"
	"phsar/internal/httpx"
	"phsar/internal/integ"
)

// fakeOIDC: Google (token endpoint + JWKS) dan Apple (JWKS) palsu dengan kunci RSA sendiri.
type fakeOIDC struct {
	srv      *httptest.Server
	key      *rsa.PrivateKey
	clientID string
	secret   string
	mu       sync.Mutex
	nonce    string // nonce yang dimasukkan ke id_token dari token endpoint
}

const qaGoogleClient = "123456789012-qaweb0000000000.apps.googleusercontent.com"
const qaGoogleApp = "123456789012-qaandroid000000.apps.googleusercontent.com"

func newFakeOIDC(t *testing.T) *fakeOIDC {
	k, _ := rsa.GenerateKey(rand.Reader, 2048)
	f := &fakeOIDC{key: k, clientID: qaGoogleClient, secret: "GOCSPX-qa-secret-1234"}
	f.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/certs", "/apple-keys":
			b64 := base64.RawURLEncoding
			_ = json.NewEncoder(w).Encode(map[string]any{"keys": []any{map[string]string{"kty": "RSA", "kid": "qa1", "use": "sig", "alg": "RS256",
				"n": b64.EncodeToString(k.N.Bytes()), "e": b64.EncodeToString(big.NewInt(int64(k.E)).Bytes())}}})
		case "/token":
			_ = r.ParseForm()
			if r.Form.Get("client_id") != f.clientID || r.Form.Get("client_secret") != f.secret {
				w.WriteHeader(401)
				_, _ = io.WriteString(w, `{"error":"invalid_client","error_description":"The OAuth client was not found."}`)
				return
			}
			if r.Form.Get("code") != "good-code" {
				w.WriteHeader(400)
				_, _ = io.WriteString(w, `{"error":"invalid_grant"}`)
				return
			}
			f.mu.Lock()
			n := f.nonce
			f.mu.Unlock()
			_ = json.NewEncoder(w).Encode(map[string]string{"id_token": f.token(t, "accounts.google.com", f.clientID, "g-sub-777", n, time.Hour), "access_token": "x"})
		default:
			w.WriteHeader(404)
		}
	}))
	integ.OAuthEndpoints.GoogleAuth, integ.OAuthEndpoints.GoogleToken = f.srv.URL+"/auth", f.srv.URL+"/token"
	integ.OAuthEndpoints.GoogleCerts, integ.OAuthEndpoints.AppleKeys = f.srv.URL+"/certs", f.srv.URL+"/apple-keys"
	t.Cleanup(func() {
		f.srv.Close()
		integ.OAuthEndpoints.GoogleAuth, integ.OAuthEndpoints.GoogleToken = integ.GoogleAuthURL, integ.GoogleTokenURL
		integ.OAuthEndpoints.GoogleCerts, integ.OAuthEndpoints.AppleKeys = integ.GoogleCertsURL, integ.AppleKeysURL
	})
	return f
}

func (f *fakeOIDC) token(t *testing.T, iss, aud, sub, nonce string, ttl time.Duration) string {
	b64 := base64.RawURLEncoding
	h := b64.EncodeToString([]byte(`{"alg":"RS256","kid":"qa1","typ":"JWT"}`))
	c, _ := json.Marshal(map[string]any{"iss": iss, "aud": aud, "sub": sub, "exp": time.Now().Add(ttl).Unix(), "iat": time.Now().Unix(),
		"nonce": nonce, "email": "dara@example.com", "name": "Dara Google"})
	signing := h + "." + b64.EncodeToString(c)
	sum := sha256.Sum256([]byte(signing))
	sig, err := rsa.SignPKCS1v15(rand.Reader, f.key, crypto.SHA256, sum[:])
	if err != nil {
		t.Fatal(err)
	}
	return signing + "." + b64.EncodeToString(sig)
}

func TestLoginGoogleApple(t *testing.T) {
	e := getEnv(t)
	fake := newFakeOIDC(t)
	integ.CacheTTL = 0
	t.Cleanup(func() { integ.CacheTTL = 15e9 })
	tok, acct := newAdminFull(t, e, "super_admin")
	code := func() string {
		e.sql(t, `UPDATE admin_users SET mfa_last_step = NULL WHERE user_id = (SELECT id FROM users WHERE phone = '+855' || substr($1, 2))`, acct.phone)
		return totpAt(acct.secret, 0)
	}
	t.Cleanup(func() {
		e.sql(t, `DELETE FROM integrations WHERE code IN ('google', 'apple')`)
	})

	t.Run("belum diatur: tombol tidak ada, endpoint menolak", func(t *testing.T) {
		check(t, strings.Contains(string(e.call(t, "GET", "/ui/config.js", nil, "").want(200).Raw), `"google":false`), "config.js harus google:false")
		e.call(t, "POST", "/v1/app/auth/google", map[string]string{"id_token": "x.y.z"}, "").wantErr(403, "google_disabled")
		e.call(t, "POST", "/v1/app/auth/apple", map[string]string{"id_token": "x.y.z"}, "").wantErr(403, "apple_disabled")
		res, _ := noRedirect.Get(e.srv.URL + "/seller/auth/google/start")
		res.Body.Close()
		check(t, strings.Contains(res.Header.Get("Location"), "auth_error=google_disabled"), "start tanpa pengaturan: %s", res.Header.Get("Location"))
	})

	t.Run("admin: simpan, validasi, uji koneksi", func(t *testing.T) {
		e.call(t, "PUT", "/admin/api/integrations/google", map[string]any{"enabled": true, "client_id": "abc", "mfa_code": code()}, tok).wantErr(422, "validation", "client_id")
		e.call(t, "PUT", "/admin/api/integrations/google", map[string]any{"enabled": true, "client_id": qaGoogleClient, "mfa_code": code()}, tok).wantErr(422, "validation", "client_id")
		e.call(t, "PUT", "/admin/api/integrations/google", map[string]any{"enabled": true, "client_id": qaGoogleClient, "client_secret": fake.secret, "app_client_ids": "bad", "mfa_code": code()}, tok).wantErr(422, "validation", "app_client_ids")
		e.call(t, "PUT", "/admin/api/integrations/google", map[string]any{"enabled": true, "client_id": qaGoogleClient, "client_secret": fake.secret}, tok).wantErr(422, "mfa_invalid", "mfa_code")
		r := e.call(t, "PUT", "/admin/api/integrations/google", map[string]any{"enabled": true, "client_id": qaGoogleClient, "client_secret": fake.secret,
			"app_client_ids": qaGoogleApp, "mfa_code": code()}, tok).want(200)
		r.eq("google.active", "true").eq("google.key_hint", "1234").eq("google.redirect_uris.1", "http://qa.local/seller/auth/google/callback")
		check(t, !strings.Contains(string(r.Raw), fake.secret), "Client Secret tidak boleh dikirim balik")
		e.call(t, "POST", "/admin/api/integrations/google/test", nil, tok).want(200).eq("ok", "true")
		check(t, strings.Contains(string(e.call(t, "GET", "/ui/config.js", nil, "").want(200).Raw), `"google":true`), "config.js harus google:true")
		check(t, strings.Contains(string(e.call(t, "GET", "/km/cart", nil, "").want(200).Raw), `data-gg="true"`), "website: tombol Google di login keranjang")
		// rahasia salah → uji gagal dengan pesan jelas
		e.call(t, "PUT", "/admin/api/integrations/google", map[string]any{"enabled": true, "client_id": qaGoogleClient, "client_secret": "GOCSPX-wrong-0000", "app_client_ids": qaGoogleApp, "mfa_code": code()}, tok).want(200)
		e.call(t, "POST", "/admin/api/integrations/google/test", nil, tok).want(200).eq("ok", "false")
		e.call(t, "PUT", "/admin/api/integrations/google", map[string]any{"enabled": true, "client_id": qaGoogleClient, "client_secret": fake.secret, "app_client_ids": qaGoogleApp, "mfa_code": code()}, tok).want(200)
	})

	t.Run("web: start → Google → callback → sesi", func(t *testing.T) {
		res, err := noRedirect.Get(e.srv.URL + "/shop/auth/google/start?return=" + url.QueryEscape("/km/checkout"))
		if err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
		loc, _ := url.Parse(res.Header.Get("Location"))
		q := loc.Query()
		check(t, strings.HasPrefix(loc.String(), fake.srv.URL+"/auth") && q.Get("client_id") == qaGoogleClient && q.Get("scope") == "openid email profile",
			"start harus ke halaman izin Google: %s", loc)
		check(t, q.Get("redirect_uri") == "http://qa.local/shop/auth/google/callback" && q.Get("nonce") != "" && q.Get("state") != "", "redirect_uri/nonce/state: %s", loc)
		var st *http.Cookie
		for _, c := range res.Cookies() {
			if c.Name == "phsar_gst" {
				st = c
			}
		}
		check(t, st != nil && st.HttpOnly, "cookie state harus HttpOnly")
		cb := func(state, code string) *http.Response {
			req, _ := http.NewRequest("GET", e.srv.URL+"/shop/auth/google/callback?state="+state+"&code="+code, nil)
			req.AddCookie(st)
			r, err := noRedirect.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			r.Body.Close()
			return r
		}
		check(t, strings.Contains(cb("wrong", "good-code").Header.Get("Location"), "auth_error=google_expired"), "state salah harus ditolak")
		fake.mu.Lock()
		fake.nonce = "other-nonce"
		fake.mu.Unlock()
		check(t, strings.Contains(cb(q.Get("state"), "good-code").Header.Get("Location"), "auth_error=google_failed"), "nonce salah harus ditolak")
		fake.mu.Lock()
		fake.nonce = q.Get("nonce")
		fake.mu.Unlock()
		ok := cb(q.Get("state"), "good-code")
		check(t, ok.Header.Get("Location") == "/km/checkout", "callback harus kembali ke halaman asal: %s", ok.Header.Get("Location"))
		var sess string
		for _, c := range ok.Cookies() {
			if c.Name == "phsar_web" {
				sess = c.Value
			}
		}
		check(t, sess != "", "cookie sesi tidak dipasang")
		e.call(t, "POST", "/shop/auth/refresh", nil, "", "X-Phsar", "1", "Cookie", "phsar_web="+sess).want(200).eq("user.name", "Dara Google")
		check(t, e.scalar(t, `SELECT display FROM auth_identities WHERE provider = 'google' AND provider_uid = 'g-sub-777'`) == "dara@example.com", "email sebagai tampilan identitas")
	})

	t.Run("Seller Center di subdomain: redirect_uri memakai SELLER_BASE_URL", func(t *testing.T) {
		cfg := e.cfg
		cfg.SellerBaseURL, cfg.AdminBaseURL = "https://seller.qa.local", "https://admin.qa.local"
		log := slog.New(slog.NewTextHandler(io.Discard, nil))
		srv := httptest.NewServer(httpx.NewRouter(cfg, log, e.pool, auth.NewService(e.pool, cfg, log)))
		defer srv.Close()
		for path, want := range map[string]string{"/seller/auth/google/start": "https://seller.qa.local/seller/auth/google/callback",
			"/admin/auth/google/start": "https://admin.qa.local/admin/auth/google/callback"} {
			res, err := noRedirect.Get(srv.URL + path)
			if err != nil {
				t.Fatal(err)
			}
			res.Body.Close()
			loc, _ := url.Parse(res.Header.Get("Location"))
			check(t, loc.Query().Get("redirect_uri") == want, "%s: redirect_uri %s", path, loc.Query().Get("redirect_uri"))
		}
	})

	t.Run("aplikasi: id_token Google diperiksa tanda tangan, audience, kedaluwarsa", func(t *testing.T) {
		e.call(t, "POST", "/v1/app/auth/google", map[string]string{"id_token": fake.token(t, "https://accounts.google.com", "999-other.apps.googleusercontent.com", "g-sub-777", "", time.Hour)}, "").wantErr(401, "google_failed")
		e.call(t, "POST", "/v1/app/auth/google", map[string]string{"id_token": fake.token(t, "https://accounts.google.com", qaGoogleApp, "g-sub-777", "", -time.Hour)}, "").wantErr(401, "google_failed")
		e.call(t, "POST", "/v1/app/auth/google", map[string]string{"id_token": fake.token(t, "https://evil.example", qaGoogleApp, "g-sub-777", "", time.Hour)}, "").wantErr(401, "google_failed")
		good := fake.token(t, "https://accounts.google.com", qaGoogleApp, "g-sub-777", "", time.Hour)
		tampered := good[:len(good)-6] + "AAAAAA"
		e.call(t, "POST", "/v1/app/auth/google", map[string]string{"id_token": tampered}, "").wantErr(401, "google_failed")
		e.call(t, "POST", "/v1/app/auth/google", map[string]string{"id_token": good}, "").want(200).eq("user.name", "Dara Google")
		check(t, e.scalar(t, `SELECT count(*) FROM auth_identities WHERE provider = 'google' AND provider_uid = 'g-sub-777'`) == "1", "satu akun Google = satu identitas (web & app)")
	})

	t.Run("Apple (aplikasi iPhone): Bundle ID, nama dari aplikasi", func(t *testing.T) {
		e.call(t, "PUT", "/admin/api/integrations/apple", map[string]any{"enabled": true, "bundle_ids": "", "mfa_code": code()}, tok).wantErr(422, "validation", "bundle_ids")
		e.call(t, "PUT", "/admin/api/integrations/apple", map[string]any{"enabled": true, "bundle_ids": "com.phsar.app", "mfa_code": code()}, tok).want(200).eq("apple.active", "true")
		check(t, !strings.Contains(string(e.call(t, "GET", "/ui/config.js", nil, "").want(200).Raw), "apple"), "tidak ada tombol Apple di website")
		e.call(t, "POST", "/v1/app/auth/apple", map[string]string{"id_token": fake.token(t, "https://appleid.apple.com", "com.other.app", "a-sub-1", "", time.Hour)}, "").wantErr(401, "apple_failed")
		r := e.call(t, "POST", "/v1/app/auth/apple", map[string]string{"id_token": fake.token(t, "https://appleid.apple.com", "com.phsar.app", "a-sub-1", "", time.Hour), "name": "Sokha Apple"}, "").want(200)
		r.eq("user.name", "Sokha Apple")
		check(t, e.scalar(t, `SELECT count(*) FROM auth_identities WHERE provider = 'apple' AND provider_uid = 'a-sub-1'`) == "1", "identitas Apple tersimpan")
	})
}
