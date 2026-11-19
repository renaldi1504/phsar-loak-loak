package qa

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"

	"phsar/internal/auth"
	"phsar/internal/httpx"
	"phsar/internal/integ"
)

// Server palsu PlasGate + Facebook (Graph API & dialog).
type fakeLoginAPIs struct {
	mu       sync.Mutex
	sms      []map[string]string // body + private_key + secret
	srv      *httptest.Server
	appID    string
	appKey   string
	fbUserID string
}

func newFakeLoginAPIs(t *testing.T) *fakeLoginAPIs {
	f := &fakeLoginAPIs{appID: "1234567890", appKey: "abcdef0123456789abcdef0123456789", fbUserID: "fb-777"}
	f.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		switch r.URL.Path {
		case "/rest/send":
			var b map[string]string
			_ = json.NewDecoder(r.Body).Decode(&b)
			if q.Get("private_key") != "pk-qa-123456" || r.Header.Get("X-Secret") != "sec-qa-654321" {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"error":"unauthorized"}`))
				return
			}
			f.mu.Lock()
			f.sms = append(f.sms, b)
			f.mu.Unlock()
			_, _ = w.Write([]byte(`{"queue_id":"q-1"}`))
		case "/oauth/access_token":
			if q.Get("client_id") != f.appID || q.Get("client_secret") != f.appKey {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(`{"error":{"message":"Invalid client_secret"}}`))
				return
			}
			if q.Get("grant_type") == "client_credentials" {
				_, _ = w.Write([]byte(`{"access_token":"app-token"}`))
				return
			}
			if q.Get("code") != "good-code" {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(`{"error":{"message":"bad code"}}`))
				return
			}
			_, _ = w.Write([]byte(`{"access_token":"user-token"}`))
		case "/debug_token":
			app := f.appID
			if q.Get("input_token") == "other-app-token" {
				app = "999"
			}
			_, _ = w.Write([]byte(`{"data":{"app_id":"` + app + `","is_valid":true}}`))
		case "/me":
			if q.Get("appsecret_proof") == "" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			_, _ = w.Write([]byte(`{"id":"` + f.fbUserID + `","name":"Dara FB"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(f.srv.Close)
	integ.Endpoints.PlasGate, integ.Endpoints.Graph, integ.Endpoints.Dialog = f.srv.URL+"/rest/send", f.srv.URL, f.srv.URL+"/dialog"
	t.Cleanup(func() {
		integ.Endpoints.PlasGate, integ.Endpoints.Graph, integ.Endpoints.Dialog = integ.PlasGateURL, integ.GraphURL, integ.DialogURL
	})
	return f
}

var noRedirect = &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}

func TestLoginSettingsSMSFacebookGuest(t *testing.T) {
	e := getEnv(t)
	fake := newFakeLoginAPIs(t)
	integ.CacheTTL = 0 // dua server di tes ini membaca pengaturan yang sama
	t.Cleanup(func() { integ.CacheTTL = 15e9 })
	tok, acct := newAdminFull(t, e, "super_admin")
	code := func() string {
		e.sql(t, `UPDATE admin_users SET mfa_last_step = NULL WHERE user_id = (SELECT id FROM users WHERE phone = '+855' || substr($1, 2))`, acct.phone)
		return totpAt(acct.secret, 0)
	}
	// Server kedua tanpa penangkap SMS tes → benar-benar memakai provider dari pengaturan admin.
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	srv := httptest.NewServer(httpx.NewRouter(e.cfg, log, e.pool, auth.NewService(e.pool, e.cfg, log)))
	defer srv.Close()
	live := *e
	live.srv = srv
	t.Cleanup(func() {
		e.call(t, "PUT", "/admin/api/integrations/sms", map[string]any{"mode": "", "mfa_code": code()}, tok).want(200)
		e.call(t, "PUT", "/admin/api/integrations/facebook", map[string]any{"enabled": false, "app_id": fake.appID, "mfa_code": code()}, tok).want(200)
	})

	t.Run("hanya super admin, wajib kode authenticator", func(t *testing.T) {
		e.call(t, "GET", "/admin/api/integrations", nil, newAdmin(t, e, "finance")).want(403)
		v := e.call(t, "GET", "/admin/api/integrations", nil, tok).want(200)
		v.eq("sms.source", "env").eq("sms.active", "log").eq("facebook.active", "false")
		check(t, strings.Contains(string(v.Raw), "/shop/auth/facebook/callback"), "redirect URI Facebook harus ditampilkan")
		e.call(t, "PUT", "/admin/api/integrations/sms", map[string]any{"mode": "off"}, tok).wantErr(422, "mfa_invalid", "mfa_code")
	})

	t.Run("PlasGate: validasi, terenkripsi, kode OTP terkirim lewat PlasGate", func(t *testing.T) {
		e.call(t, "PUT", "/admin/api/integrations/sms", map[string]any{"mode": "plasgate", "mfa_code": code()}, tok).wantErr(422, "validation", "mode")
		e.call(t, "PUT", "/admin/api/integrations/sms", map[string]any{"mode": "off", "private_key": "pk-qa-123456", "mfa_code": code()}, tok).wantErr(422, "validation", "secret")
		e.call(t, "PUT", "/admin/api/integrations/sms", map[string]any{"mode": "off", "sender": "Phsar!!", "mfa_code": code()}, tok).wantErr(422, "validation", "sender")
		r := e.call(t, "PUT", "/admin/api/integrations/sms", map[string]any{"mode": "plasgate", "sender": "Phsar",
			"private_key": "pk-qa-123456", "secret": "sec-qa-654321", "mfa_code": code()}, tok).want(200)
		r.eq("sms.active", "plasgate").eq("sms.key_hint", "4321").eq("sms.source", "admin")
		check(t, !strings.Contains(string(r.Raw), "pk-qa-123456") && !strings.Contains(string(r.Raw), "sec-qa-654321"), "kunci tidak boleh dikirim balik")

		cfgJS := live.call(t, "GET", "/ui/config.js", nil, "").want(200)
		check(t, strings.Contains(string(cfgJS.Raw), `"sms":true`), "config.js: %s", cfgJS.Raw)
		live.call(t, "POST", "/v1/app/auth/otp/request", map[string]string{"phone": "015 222 333"}, "").want(200)
		fake.mu.Lock()
		last := fake.sms[len(fake.sms)-1]
		fake.mu.Unlock()
		check(t, last["to"] == "85515222333" && last["sender"] == "Phsar" && strings.Contains(last["content"], "Phsar"), "SMS PlasGate salah: %v", last)

		e.call(t, "POST", "/admin/api/integrations/sms/test", map[string]string{"phone": "015 222 444"}, tok).want(200).eq("ok", "true")
		e.call(t, "PUT", "/admin/api/integrations/sms", map[string]any{"mode": "off", "mfa_code": code()}, tok).want(200).eq("sms.active", "off")
		live.call(t, "POST", "/v1/app/auth/otp/request", map[string]string{"phone": "015 222 555"}, "").wantErr(503, "sms_disabled")
	})

	t.Run("Facebook web: start → callback → sesi", func(t *testing.T) {
		e.call(t, "PUT", "/admin/api/integrations/facebook", map[string]any{"enabled": true, "app_id": "12x", "mfa_code": code()}, tok).wantErr(422, "validation", "app_id")
		e.call(t, "PUT", "/admin/api/integrations/facebook", map[string]any{"enabled": true, "app_id": fake.appID, "mfa_code": code()}, tok).wantErr(422, "validation", "app_id")
		r := e.call(t, "PUT", "/admin/api/integrations/facebook", map[string]any{"enabled": true, "app_id": fake.appID, "app_secret": fake.appKey, "mfa_code": code()}, tok).want(200)
		r.eq("facebook.active", "true").eq("facebook.key_hint", "6789")
		check(t, !strings.Contains(string(r.Raw), fake.appKey), "App Secret tidak boleh dikirim balik")
		e.call(t, "POST", "/admin/api/integrations/facebook/test", nil, tok).want(200).eq("ok", "true")

		res, err := noRedirect.Get(e.srv.URL + "/shop/auth/facebook/start?return=" + url.QueryEscape("/km/checkout"))
		if err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
		loc, _ := url.Parse(res.Header.Get("Location"))
		check(t, res.StatusCode == 302 && strings.HasPrefix(loc.String(), fake.srv.URL+"/dialog") && loc.Query().Get("client_id") == fake.appID,
			"start harus ke dialog Facebook: %d %s", res.StatusCode, loc)
		check(t, loc.Query().Get("redirect_uri") == "http://qa.local/shop/auth/facebook/callback", "redirect_uri: %s", loc.Query().Get("redirect_uri"))
		var stateCookie *http.Cookie
		for _, c := range res.Cookies() {
			if c.Name == "phsar_fbst" {
				stateCookie = c
			}
		}
		check(t, stateCookie != nil && stateCookie.HttpOnly, "cookie state harus HttpOnly")

		cb := func(state, code string) *http.Response {
			req, _ := http.NewRequest("GET", e.srv.URL+"/shop/auth/facebook/callback?state="+state+"&code="+code, nil)
			req.AddCookie(stateCookie)
			res, err := noRedirect.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			res.Body.Close()
			return res
		}
		bad := cb("wrong", "good-code")
		check(t, strings.Contains(bad.Header.Get("Location"), "auth_error=facebook_expired"), "state salah harus ditolak: %s", bad.Header.Get("Location"))
		ok := cb(loc.Query().Get("state"), "good-code")
		check(t, ok.StatusCode == 302 && ok.Header.Get("Location") == "/km/checkout", "callback harus kembali ke halaman asal: %s", ok.Header.Get("Location"))
		var sess string
		for _, c := range ok.Cookies() {
			if c.Name == "phsar_web" {
				sess = c.Value
			}
		}
		check(t, sess != "", "cookie sesi tidak dipasang")
		me := e.call(t, "POST", "/shop/auth/refresh", nil, "", "X-Phsar", "1", "Cookie", "phsar_web="+sess).want(200)
		me.eq("user.name", "Dara FB")

		// Open redirect ditolak
		res2, _ := noRedirect.Get(e.srv.URL + "/shop/auth/facebook/start?return=" + url.QueryEscape("//evil.example.com"))
		res2.Body.Close()
		check(t, res2.Header.Get("Set-Cookie") != "" && !strings.Contains(res2.Header.Get("Set-Cookie"), "evil"), "return ke situs lain harus diabaikan")
	})

	t.Run("Facebook app: token diperiksa milik App ID Phsar", func(t *testing.T) {
		e.call(t, "POST", "/v1/app/auth/facebook", map[string]string{"access_token": "other-app-token"}, "").wantErr(401, "facebook_failed")
		r := e.call(t, "POST", "/v1/app/auth/facebook", map[string]string{"access_token": "user-token"}, "").want(200)
		r.eq("user.name", "Dara FB")
		check(t, e.scalar(t, `SELECT count(*) FROM auth_identities WHERE provider = 'facebook' AND provider_uid = $1`, fake.fbUserID) == "1", "satu identitas Facebook = satu akun")
	})

	t.Run("checkout sebagai tamu", func(t *testing.T) {
		g := e.call(t, "POST", "/shop/auth/guest", nil, "", csrf...).want(200)
		g.eq("user.is_guest", "true")
		gt := g.str("access_token")
		e.call(t, "GET", "/v1/shop/address", nil, gt).want(200)
		e.call(t, "POST", "/v1/seller/stores", map[string]any{"name": "Guest shop"}, gt).wantErr(403, "guest_not_allowed")
		e.call(t, "POST", "/v1/app/auth/guest", nil, "").want(200).eq("user.is_guest", "true")
		e.call(t, "POST", "/shop/auth/guest", nil, "").wantErr(403, "csrf")
		// batas per IP
		e.sql(t, `INSERT INTO login_attempts (identifier, ip, success) SELECT 'guest', '127.0.0.1', true FROM generate_series(1, 40)`)
		e.call(t, "POST", "/shop/auth/guest", nil, "", csrf...).wantErr(429, "guest_rate_limited")
		e.sql(t, `DELETE FROM login_attempts WHERE identifier = 'guest'`)
		_ = context.Background()
	})
}
