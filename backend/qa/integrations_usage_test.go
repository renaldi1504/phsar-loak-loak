package qa

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"testing"

	"phsar/internal/httpx"
	"phsar/internal/integ"
)

// ---------------------------------------------------------------------------
// Integrasi → Pemakaian & saldo: catatan pemakaian AI/SMS, perkiraan saldo, pengingat (dasbor + badge).
// ---------------------------------------------------------------------------

func TestIntegrationsUsage(t *testing.T) {
	e := getEnv(t)
	integ.CacheTTL = 0
	t.Cleanup(func() { integ.CacheTTL = 15e9 })
	tok, acct := newAdminFull(t, e, "super_admin")
	code := func() string {
		e.sql(t, `UPDATE admin_users SET mfa_last_step = NULL WHERE user_id = (SELECT id FROM users WHERE phone = '+855' || substr($1, 2))`, acct.phone)
		return totpAt(acct.secret, 0)
	}
	mod := newAdmin(t, e, "moderator")
	clean := func() {
		e.sql(t, `DELETE FROM integration_usage`)
		e.sql(t, `DELETE FROM platform_settings WHERE key = 'usage_budget'`)
		e.sql(t, `DELETE FROM integrations WHERE code IN ('anthropic', 'gemini')`)
	}
	clean()
	t.Cleanup(clean)
	st := httpx.NewIntegrations(e.cfg, slog.Default(), e.pool)
	ctx := context.Background()

	t.Run("hanya super admin", func(t *testing.T) {
		e.call(t, "GET", "/admin/api/integrations/usage", nil, mod).want(403)
		e.call(t, "PUT", "/admin/api/integrations/usage/budget", map[string]any{"provider": "anthropic"}, mod).want(403)
	})

	t.Run("klasifikasi error", func(t *testing.T) {
		for msg, want := range map[string]string{
			"claude: HTTP 400 invalid_request_error: Your credit balance is too low": "billing",
			"gemini: HTTP 429 RESOURCE_EXHAUSTED: quota":                             "quota",
			"claude: HTTP 401 authentication_error: invalid x-api-key":               "auth",
			"gemini: tidak bisa menghubungi server":                                  "other",
		} {
			check(t, integ.ClassifyErr(errors.New(msg)) == want, "%q → %s, mau %s", msg, integ.ClassifyErr(errors.New(msg)), want)
		}
	})

	t.Run("biaya dihitung dari harga; gagal = $0", func(t *testing.T) {
		// harga bawaan Anthropic: $1 / $5 per 1 juta token
		st.RecordUsage(ctx, integ.Usage{Provider: "anthropic", Purpose: "translate", Model: "claude-haiku-5-5", In: 1_000_000, Out: 200_000})
		st.RecordUsage(ctx, integ.Usage{Provider: "anthropic", Purpose: "askbot", Err: errors.New("HTTP 500 overloaded")})
		st.RecordUsage(ctx, integ.Usage{Provider: "gemini", Purpose: "translate", In: 5000, Out: 900})
		r := e.call(t, "GET", "/admin/api/integrations/usage", nil, tok).want(200)
		r.eq("providers.0.provider", "anthropic").eq("providers.0.today.requests", "2").eq("providers.0.today.failed", "1")
		r.eq("providers.0.month.cost_usd", "2").eq("providers.0.last_error.kind", "other")
		r.eq("providers.1.provider", "gemini").eq("providers.1.today.cost_usd", "0").eq("providers.1.fallbacks_7d", "1")
		check(t, r.str("providers.0.remaining_usd") == "" || r.str("providers.0.remaining_usd") == "null", "saldo belum diisi: %s", r.str("providers.0.remaining_usd"))
	})

	t.Run("saldo + pengingat saldo rendah", func(t *testing.T) {
		e.call(t, "PUT", "/admin/api/integrations/usage/budget", map[string]any{"provider": "anthropic", "balance_usd": "3"}, tok).wantErr(422, "mfa_invalid", "mfa_code")
		e.call(t, "PUT", "/admin/api/integrations/usage/budget", map[string]any{"provider": "anthropic", "balance_usd": "abc", "mfa_code": code()}, tok).wantErr(422, "validation", "balance_usd")
		e.call(t, "PUT", "/admin/api/integrations/usage/budget", map[string]any{"provider": "openai", "mfa_code": code()}, tok).wantErr(422, "validation", "provider")
		r := e.call(t, "PUT", "/admin/api/integrations/usage/budget", map[string]any{"provider": "anthropic", "balance_usd": "3", "warn_below_usd": "2",
			"price_in": "1", "price_out": "5", "mfa_code": code()}, tok).want(200)
		r.eq("providers.0.remaining_usd", "3").eq("providers.0.budget.balance_usd", "3")
		// Kunci Anthropic aktif (peringatan hanya untuk layanan yang dipakai).
		e.call(t, "PUT", "/admin/api/integrations/ai", map[string]any{"provider": "anthropic", "enabled": true, "key": "sk-ant-api03-" + strings.Repeat("u", 40) + "US01", "mfa_code": code()}, tok).want(200)
		// Pakai $1.50 sesudah saldo diisi → sisa $1.50 < $2 → peringatan kuning.
		st.RecordUsage(ctx, integ.Usage{Provider: "anthropic", Purpose: "translate", In: 500_000, Out: 200_000})
		r = e.call(t, "GET", "/admin/api/integrations/usage", nil, tok).want(200)
		r.eq("providers.0.remaining_usd", "1.5").eq("alerts.0.provider", "anthropic").eq("alerts.0.code", "low_balance").eq("alerts.0.level", "warn").eq("alerts.0.value", "1.50")
		d := e.call(t, "GET", "/admin/api/dashboard", nil, tok).want(200)
		check(t, strings.Contains(string(d.Raw), `"key":"integrations"`) && strings.Contains(string(d.Raw), `"low_balance"`), "dasbor harus memberi peringatan: %s", d.Raw)
		md := e.call(t, "GET", "/admin/api/dashboard", nil, mod).want(200)
		check(t, !strings.Contains(string(md.Raw), "integration_alerts"), "moderator tidak boleh melihat peringatan integrasi")
	})

	t.Run("saldo habis (HTTP 400 credit balance) → merah, hilang setelah sukses", func(t *testing.T) {
		st.RecordUsage(ctx, integ.Usage{Provider: "anthropic", Purpose: "translate", Err: errors.New("claude: HTTP 400 invalid_request_error: Your credit balance is too low to access the Anthropic API.")})
		r := e.call(t, "GET", "/admin/api/integrations/usage", nil, tok).want(200)
		check(t, strings.Contains(string(r.Raw), `"code":"billing"`) && strings.Contains(string(r.Raw), `"level":"bad"`), "harus ada peringatan saldo habis: %s", r.Raw)
		// Top up: saldo baru $10, lalu panggilan sukses → peringatan billing & saldo rendah hilang.
		e.call(t, "PUT", "/admin/api/integrations/usage/budget", map[string]any{"provider": "anthropic", "balance_usd": "10", "warn_below_usd": "2",
			"price_in": "1", "price_out": "5", "mfa_code": code()}, tok).want(200)
		st.RecordUsage(ctx, integ.Usage{Provider: "anthropic", Purpose: "translate", In: 1000, Out: 100})
		r = e.call(t, "GET", "/admin/api/integrations/usage", nil, tok).want(200)
		check(t, !strings.Contains(string(r.Raw), `"code":"billing"`) && !strings.Contains(string(r.Raw), `"code":"low_balance"`), "peringatan harus hilang: %s", r.Raw)
	})

	t.Run("Gemini: batas harian gratis 80%", func(t *testing.T) {
		e.call(t, "PUT", "/admin/api/integrations/ai", map[string]any{"provider": "gemini", "enabled": true, "key": "AIzaSy" + strings.Repeat("u", 33), "mfa_code": code()}, tok).want(200)
		e.call(t, "PUT", "/admin/api/integrations/usage/budget", map[string]any{"provider": "gemini", "daily_limit": 5, "price_in": "0", "price_out": "0", "mfa_code": code()}, tok).want(200)
		for i := 0; i < 3; i++ {
			st.RecordUsage(ctx, integ.Usage{Provider: "gemini", Purpose: "askbot", In: 10, Out: 10})
		}
		r := e.call(t, "GET", "/admin/api/integrations/usage", nil, tok).want(200)
		check(t, strings.Contains(string(r.Raw), `"code":"daily_limit"`) && strings.Contains(string(r.Raw), `"value":"4/5"`), "peringatan batas harian: %s", r.Raw)
		r.eq("days.0.provider", r.str("days.0.provider")) // tabel harian ada
		check(t, len(r.str("days.0.day")) == 10, "tanggal harian: %s", r.str("days.0.day"))
	})

	t.Run("hapus saldo", func(t *testing.T) {
		r := e.call(t, "PUT", "/admin/api/integrations/usage/budget", map[string]any{"provider": "anthropic", "balance_usd": "", "warn_below_usd": "2", "mfa_code": code()}, tok).want(200)
		r.eq("providers.0.budget.balance_usd", "")
	})
}
