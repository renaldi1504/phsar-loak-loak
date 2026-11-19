package qa

import (
	"strings"
	"testing"
)

// Super admin mengatur payment gateway dari admin: kredensial terenkripsi, tidak pernah dikirim balik,
// menyimpan wajib kode authenticator, dan mode "off" langsung menutup checkout.
func TestPaymentGatewaySettings(t *testing.T) {
	e := getEnv(t)
	tok, acct := newAdminFull(t, e, "super_admin")
	// kode authenticator baru untuk setiap simpan (kode yang sama tidak boleh dipakai dua kali)
	code := func() string {
		e.sql(t, `UPDATE admin_users SET mfa_last_step = NULL WHERE user_id = (SELECT id FROM users WHERE phone = '+855' || substr($1, 2))`, acct.phone)
		return totpAt(acct.secret, 0)
	}
	t.Cleanup(func() { // kembali ikut .env (tes lain memakai pembayaran palsu)
		e.call(t, "PUT", "/admin/api/payments/settings", map[string]any{"mode": "", "mfa_code": code()}, tok).want(200)
	})

	t.Run("hanya super admin", func(t *testing.T) {
		e.call(t, "GET", "/admin/api/payments", nil, newAdmin(t, e, "finance")).want(403)
		e.call(t, "GET", "/admin/api/payments", nil, newAdmin(t, e, "moderator")).want(403)
		v := e.call(t, "GET", "/admin/api/payments", nil, tok).want(200)
		v.eq("source", "env").eq("active", "fake").eq("payway.key_set", "false")
	})

	t.Run("simpan wajib kode authenticator", func(t *testing.T) {
		e.call(t, "PUT", "/admin/api/payments/settings", map[string]any{"mode": "off", "lifetime_min": 15}, tok).wantErr(422, "mfa_invalid", "mfa_code")
		e.call(t, "PUT", "/admin/api/payments/settings", map[string]any{"mode": "off", "lifetime_min": 15, "mfa_code": "000000"}, tok).wantErr(422, "mfa_invalid", "mfa_code")
	})

	t.Run("PayWay belum bisa dipilih tanpa kredensial", func(t *testing.T) {
		e.call(t, "PUT", "/admin/api/payments/settings", map[string]any{"mode": "payway", "lifetime_min": 15, "mfa_code": code()}, tok).wantErr(422, "validation", "mode")
		e.call(t, "PUT", "/admin/api/payments/settings", map[string]any{"mode": "off", "lifetime_min": 500, "mfa_code": code()}, tok).wantErr(422, "validation", "lifetime_min")
	})

	t.Run("kredensial PayWay: validasi, terenkripsi, tidak dikirim balik", func(t *testing.T) {
		pw := func(m map[string]any) map[string]any { m["mfa_code"] = code(); return m }
		e.call(t, "PUT", "/admin/api/payments/payway", pw(map[string]any{"environment": "sandbox", "merchant_id": "bad id!", "api_key": "k-0123456789"}), tok).wantErr(422, "validation", "merchant_id")
		e.call(t, "PUT", "/admin/api/payments/payway", pw(map[string]any{"environment": "live", "merchant_id": "ec000123"}), tok).wantErr(422, "validation", "environment")
		for _, bad := range []string{"http://checkout-sandbox.payway.com.kh", "https://evil.example.com", "https://payway.com.kh.evil.com", "https://user@checkout.payway.com.kh"} {
			e.call(t, "PUT", "/admin/api/payments/payway", pw(map[string]any{"environment": "sandbox", "merchant_id": "ec000123", "base_url": bad}), tok).wantErr(422, "validation", "base_url")
		}
		e.call(t, "PUT", "/admin/api/payments/payway", pw(map[string]any{"environment": "sandbox", "merchant_id": "ec000123", "api_key": "short"}), tok).wantErr(422, "validation", "api_key")

		const key = "qa-secret-payway-key-7f3a"
		r := e.call(t, "PUT", "/admin/api/payments/payway", pw(map[string]any{"environment": "sandbox", "merchant_id": "ec000123", "api_key": key}), tok).want(200)
		r.eq("payway.key_set", "true").eq("payway.key_hint", "7f3a").eq("payway.merchant_id", "ec000123").eq("payway.source", "admin")
		check(t, !strings.Contains(string(r.Raw), key), "API key tidak boleh dikirim balik ke browser")
		check(t, e.scalar(t, `SELECT position($1 in encode(secret_enc, 'escape')) FROM payment_gateways WHERE code = 'payway'`, key) == "0", "API key harus terenkripsi di database")

		// ubah merchant tanpa mengisi key → key lama tetap
		r = e.call(t, "PUT", "/admin/api/payments/payway", pw(map[string]any{"environment": "production", "merchant_id": "ec000999"}), tok).want(200)
		r.eq("payway.key_hint", "7f3a").eq("payway.environment", "production").eq("payway.base_url", "https://checkout.payway.com.kh")

		au := e.call(t, "GET", "/admin/api/audit?table=payment_gateways", nil, tok).want(200)
		check(t, strings.Contains(string(au.Raw), `"secret_enc":"•••"`) && !strings.Contains(string(au.Raw), key), "audit log harus menyamarkan rahasia: %s", au.Raw)
	})

	t.Run("mode off menutup checkout; kembali ke .env membukanya lagi", func(t *testing.T) {
		e.call(t, "PUT", "/admin/api/payments/settings", map[string]any{"mode": "payway", "lifetime_min": 20, "mfa_code": code()}, tok).want(200).
			eq("source", "admin").eq("active", "payway").eq("lifetime_min", "20")
		e.call(t, "PUT", "/admin/api/payments/settings", map[string]any{"mode": "off", "lifetime_min": 20, "mfa_code": code()}, tok).want(200).eq("active", "off")
		buyer, _ := e.login(t)
		e.call(t, "POST", "/v1/shop/checkouts", checkoutBody("k-gw-off-1", "1.00"), buyer).wantErr(409, "checkout_closed")
		home := e.call(t, "GET", "/en/", nil, "").want(200)
		check(t, strings.Contains(string(home.Raw), `class="prelaunch"`), "mode off harus menampilkan pita segera dibuka")
		e.call(t, "PUT", "/admin/api/payments/settings", map[string]any{"mode": "", "mfa_code": code()}, tok).want(200).eq("source", "env").eq("active", "fake")
		home = e.call(t, "GET", "/en/", nil, "").want(200)
		check(t, !strings.Contains(string(home.Raw), `class="prelaunch"`), "kembali ke .env (fake) harus membuka checkout")
	})
}
