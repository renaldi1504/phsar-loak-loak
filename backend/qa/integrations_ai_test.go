package qa

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"phsar/internal/integ"
	"phsar/internal/translate"
)

// ---------------------------------------------------------------------------
// Halaman Integrasi: kunci AI (Anthropic utama, Gemini cadangan) & Google Maps diatur super admin.
// ---------------------------------------------------------------------------

func TestIntegrationsAIKeys(t *testing.T) {
	e := getEnv(t)
	integ.CacheTTL = 0
	t.Cleanup(func() { integ.CacheTTL = 15e9 })
	tok, acct := newAdminFull(t, e, "super_admin")
	code := func() string {
		e.sql(t, `UPDATE admin_users SET mfa_last_step = NULL WHERE user_id = (SELECT id FROM users WHERE phone = '+855' || substr($1, 2))`, acct.phone)
		return totpAt(acct.secret, 0)
	}
	mod := newAdmin(t, e, "moderator")
	t.Cleanup(func() { e.sql(t, `DELETE FROM integrations WHERE code IN ('anthropic', 'gemini', 'google_maps')`) })

	anthropicKey := "sk-ant-api03-" + strings.Repeat("a", 40) + "QA12"
	geminiKey := "AIzaSy" + strings.Repeat("g", 33)

	t.Run("hanya super admin", func(t *testing.T) {
		e.call(t, "GET", "/admin/api/integrations", nil, mod).want(403)
		e.call(t, "PUT", "/admin/api/integrations/ai", map[string]any{"provider": "gemini", "enabled": true, "key": geminiKey}, mod).want(403)
		r := e.call(t, "GET", "/admin/api/integrations", nil, tok).want(200)
		r.eq("ai.anthropic.model", r.str("ai.anthropic.model")) // ada
	})

	t.Run("validasi + wajib kode MFA", func(t *testing.T) {
		e.call(t, "PUT", "/admin/api/integrations/ai", map[string]any{"provider": "anthropic", "enabled": true, "key": anthropicKey}, tok).wantErr(422, "mfa_invalid", "mfa_code")
		e.call(t, "PUT", "/admin/api/integrations/ai", map[string]any{"provider": "anthropic", "enabled": true, "key": "abc", "mfa_code": code()}, tok).wantErr(422, "validation", "key")
		e.call(t, "PUT", "/admin/api/integrations/ai", map[string]any{"provider": "openai", "enabled": true, "mfa_code": code()}, tok).wantErr(422, "validation", "provider")
		e.call(t, "PUT", "/admin/api/integrations/ai", map[string]any{"provider": "gemini", "enabled": true, "model": "bad model!", "mfa_code": code()}, tok).wantErr(422, "validation", "model")
	})

	t.Run("simpan kunci: hanya 4 karakter terakhir yang tampil, urutan Anthropic → Gemini", func(t *testing.T) {
		e.call(t, "PUT", "/admin/api/integrations/ai", map[string]any{"provider": "anthropic", "enabled": true, "key": anthropicKey, "mfa_code": code()}, tok).want(200)
		r := e.call(t, "PUT", "/admin/api/integrations/ai", map[string]any{"provider": "gemini", "enabled": true, "model": "gemini-3.5-flash", "key": geminiKey, "mfa_code": code()}, tok).want(200)
		r.eq("ai.anthropic.active", "true").eq("ai.anthropic.source", "admin").eq("ai.anthropic.key_hint", "QA12")
		r.eq("ai.gemini.active", "true").eq("ai.order.0", "anthropic").eq("ai.order.1", "gemini")
		check(t, !strings.Contains(string(r.Raw), anthropicKey) && !strings.Contains(string(r.Raw), geminiKey), "kunci penuh bocor ke browser")
		enc := e.scalar(t, `SELECT encode(secret_enc, 'escape') FROM integrations WHERE code = 'gemini'`)
		check(t, !strings.Contains(enc, geminiKey), "kunci harus terenkripsi di database")
	})

	t.Run("tes kunci Gemini (server palsu)", func(t *testing.T) {
		fake := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("x-goog-api-key") != geminiKey {
				w.WriteHeader(400)
				_, _ = w.Write([]byte(`{"error":{"code":400,"message":"API key not valid","status":"INVALID_ARGUMENT"}}`))
				return
			}
			_, _ = w.Write([]byte(`{"candidates":[{"content":{"parts":[{"text":"{\"km\":{\"name\":\"អាវ\"}}"}]}}]}`))
		}))
		defer fake.Close()
		old := translate.GeminiBaseURL
		translate.GeminiBaseURL = fake.URL
		defer func() { translate.GeminiBaseURL = old }()
		e.call(t, "POST", "/admin/api/integrations/ai/test", map[string]any{"provider": "gemini"}, tok).want(200).eq("ok", "true")
		e.call(t, "GET", "/admin/api/integrations", nil, tok).want(200).eq("ai.gemini.last_test.ok", "true")
	})

	t.Run("matikan Anthropic → hanya Gemini; hapus kunci", func(t *testing.T) {
		r := e.call(t, "PUT", "/admin/api/integrations/ai", map[string]any{"provider": "anthropic", "enabled": false, "mfa_code": code()}, tok).want(200)
		r.eq("ai.anthropic.active", "false").eq("ai.order.0", "gemini")
		r = e.call(t, "PUT", "/admin/api/integrations/ai", map[string]any{"provider": "gemini", "enabled": true, "clear": true, "mfa_code": code()}, tok).want(200)
		r.eq("ai.gemini.key_set", "false")
	})

	t.Run("Google Maps: kunci dari admin dipakai config.js, lalu kembali ke .env", func(t *testing.T) {
		mapsKey := "AIzaSyQA" + strings.Repeat("m", 31)
		e.call(t, "PUT", "/admin/api/integrations/maps", map[string]any{"key": "nope", "mfa_code": code()}, tok).wantErr(422, "validation", "key")
		e.call(t, "PUT", "/admin/api/integrations/maps", map[string]any{"key": mapsKey, "map_id": "qa-map", "mfa_code": code()}, tok).want(200).eq("maps.source", "admin")
		js := e.call(t, "GET", "/ui/config.js", nil, "").want(200)
		check(t, strings.Contains(string(js.Raw), mapsKey) && strings.Contains(string(js.Raw), "qa-map"), "config.js: %s", js.Raw)
		e.call(t, "PUT", "/admin/api/integrations/maps", map[string]any{"clear": true, "map_id": "", "mfa_code": code()}, tok).want(200).eq("maps.source", "env")
		js = e.call(t, "GET", "/ui/config.js", nil, "").want(200)
		check(t, strings.Contains(string(js.Raw), `"googleMapsKey":"qa-maps-key"`), "harus kembali ke .env: %s", js.Raw)
	})
}
