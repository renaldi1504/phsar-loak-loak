package qa

import (
	"context"
	"io"
	"log/slog"
	"net/http/httptest"
	"strings"
	"testing"

	"phsar/internal/auth"
	"phsar/internal/httpx"
)

// Server tahap pra-peluncuran: PAYMENT_PROVIDER=off (checkout ditutup), SMS_PROVIDER=off (login hanya Telegram).
func TestPrelaunchMode(t *testing.T) {
	e := getEnv(t)
	cfg := e.cfg
	cfg.PaymentProvider, cfg.SMSProvider, cfg.TelegramMode = "off", "off", "webhook"
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	srv := httptest.NewServer(httpx.NewRouter(cfg, log, e.pool, auth.NewService(e.pool, cfg, log)))
	defer srv.Close()
	pre := *e
	pre.srv = srv

	t.Run("halaman web tahu SMS mati", func(t *testing.T) {
		r := pre.call(t, "GET", "/ui/config.js", nil, "").want(200)
		check(t, strings.Contains(string(r.Raw), `"sms":false`) && strings.Contains(string(r.Raw), `"telegram":true`), "config.js: %s", r.Raw)
	})
	t.Run("minta kode SMS ditolak dengan pesan jelas", func(t *testing.T) {
		pre.call(t, "POST", "/seller/auth/otp/request", map[string]string{"phone": "012 345 678"}, "", csrf...).wantErr(503, "sms_disabled")
		pre.call(t, "POST", "/v1/app/auth/otp/request", map[string]string{"phone": "012 345 678"}, "").wantErr(503, "sms_disabled")
	})
	t.Run("website: pita segera dibuka, keranjang tanpa tombol checkout", func(t *testing.T) {
		home := pre.call(t, "GET", "/en/", nil, "").want(200)
		check(t, strings.Contains(string(home.Raw), `class="prelaunch"`) && strings.Contains(string(home.Raw), `href="/en/sell"`), "pita pra-peluncuran tidak ada")
		cart := pre.call(t, "GET", "/km/cart", nil, "").want(200)
		check(t, strings.Contains(string(cart.Raw), `data-closed="true"`) && strings.Contains(string(cart.Raw), `data-sms="false"`), "halaman keranjang tidak tahu checkout ditutup")
		normal := e.call(t, "GET", "/en/", nil, "").want(200)
		check(t, !strings.Contains(string(normal.Raw), `class="prelaunch"`), "server normal tidak boleh menampilkan pita")
	})
	t.Run("API checkout ditolak", func(t *testing.T) {
		buyer, _ := e.login(t)
		pre.call(t, "POST", "/v1/shop/checkouts", checkoutBody("k-prelaunch-1", "1.00"), buyer).wantErr(409, "checkout_closed")
		check(t, e.scalar(t, `SELECT count(*) FROM checkouts WHERE idempotency_key = 'k-prelaunch-1'`) == "0", "checkout tidak boleh tersimpan")
	})
}

// Akun yang login hanya lewat Telegram tidak punya nomor HP → admin ditambahkan lewat @username.
func TestAdminByTelegramHandle(t *testing.T) {
	e := getEnv(t)
	tok := newAdmin(t, e, "super_admin")
	_, phone := e.login(t)
	uid := e.scalar(t, `SELECT id::text FROM users WHERE phone = '+855' || substr($1, 2)`, phone)
	for _, q := range []string{`UPDATE users SET phone = NULL WHERE id = $1`,
		`INSERT INTO auth_identities (user_id, provider, provider_uid, display) VALUES ($1, 'telegram', '880011223', '@QA_TgOnly')`} {
		if _, err := e.pool.Exec(context.Background(), q, uid); err != nil {
			t.Fatal(err)
		}
	}
	e.call(t, "POST", "/admin/api/admins", map[string]string{"phone": "@ab", "role": "cs"}, tok).wantErr(422, "validation", "phone")
	e.call(t, "POST", "/admin/api/admins", map[string]string{"phone": "@nobody_here", "role": "cs"}, tok).wantErr(404, "user_not_found")
	e.call(t, "POST", "/admin/api/admins", map[string]string{"phone": "@qa_tgonly", "role": "cs"}, tok).want(200)
	l := e.call(t, "GET", "/admin/api/admins", nil, tok).want(200)
	check(t, strings.Contains(string(l.Raw), `"phone":"@QA_TgOnly"`), "admin Telegram harus tampil dengan @username: %s", l.Raw)
}
