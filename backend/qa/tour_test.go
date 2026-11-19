package qa

import (
	"testing"
)

// ---------------------------------------------------------------------------
// Tur "cara pakai": status sudah dilihat per akun (pembeli /v1/app/tours, seller /v1/seller/tours).
// ---------------------------------------------------------------------------

func TestTourStatus(t *testing.T) {
	e := getEnv(t)
	user, _ := e.login(t)
	other, _ := e.login(t)

	t.Run("harus login", func(t *testing.T) {
		e.call(t, "GET", "/v1/app/tours", nil, "").wantErr(401, "unauthenticated")
		e.call(t, "POST", "/v1/seller/tours", map[string]string{"tour_key": "seller_intro"}, "").wantErr(401, "unauthenticated")
	})

	t.Run("awal: belum ada tur yang dilihat", func(t *testing.T) {
		r := e.call(t, "GET", "/v1/seller/tours", nil, user).want(200)
		check(t, len(r.Body["done"].(map[string]any)) == 0, "done harus kosong: %s", truncate(string(r.Raw)))
		check(t, r.len("keys") == 3, "keys seller = 3: %s", truncate(string(r.Raw)))
	})

	t.Run("tandai selesai: idempoten, done_at tetap", func(t *testing.T) {
		a := e.call(t, "POST", "/v1/seller/tours", map[string]string{"tour_key": "seller_intro"}, user).want(200)
		a.eq("tour_key", "seller_intro")
		b := e.call(t, "POST", "/v1/seller/tours", map[string]string{"tour_key": "seller_intro"}, user).want(200)
		check(t, a.str("done_at") != "" && a.str("done_at") == b.str("done_at"), "done_at berubah: %s vs %s", a.str("done_at"), b.str("done_at"))
		check(t, e.scalar(t, `SELECT count(*) FROM user_tours WHERE tour_key = 'seller_intro'`) == "1", "harus satu baris")
		r := e.call(t, "GET", "/v1/seller/tours", nil, user).want(200)
		check(t, r.str("done.seller_intro") == a.str("done_at"), "GET harus memuat seller_intro: %s", truncate(string(r.Raw)))
	})

	t.Run("status milik akun sendiri", func(t *testing.T) {
		r := e.call(t, "GET", "/v1/seller/tours", nil, other).want(200)
		check(t, len(r.Body["done"].(map[string]any)) == 0, "akun lain tidak boleh melihat status: %s", truncate(string(r.Raw)))
	})

	t.Run("validasi tour_key", func(t *testing.T) {
		e.call(t, "POST", "/v1/seller/tours", map[string]string{"tour_key": ""}, user).wantErr(422, "validation", "tour_key")
		e.call(t, "POST", "/v1/seller/tours", map[string]string{"tour_key": "hack'; drop"}, user).wantErr(422, "validation", "tour_key")
		// kunci pembeli tidak boleh lewat pintu seller, dan sebaliknya
		e.call(t, "POST", "/v1/seller/tours", map[string]string{"tour_key": "site"}, user).wantErr(422, "validation", "tour_key")
		e.call(t, "POST", "/v1/app/tours", map[string]string{"tour_key": "seller_intro"}, user).wantErr(422, "validation", "tour_key")
		e.call(t, "POST", "/v1/app/tours", map[string]any{"tour_key": "site", "extra": 1}, user).wantErr(400, "bad_json")
	})

	t.Run("pembeli: website, pintu terpisah dari seller", func(t *testing.T) {
		e.call(t, "POST", "/v1/app/tours", map[string]string{"tour_key": "site"}, user).want(200).eq("tour_key", "site")
		r := e.call(t, "GET", "/v1/app/tours", nil, user).want(200)
		check(t, r.str("done.site") != "", "done.site kosong: %s", truncate(string(r.Raw)))
		check(t, r.get("done.seller_intro") == nil, "kunci seller tidak boleh muncul di /v1/app: %s", truncate(string(r.Raw)))
	})

	t.Run("tamu ditolak", func(t *testing.T) {
		e.sql(t, `DELETE FROM login_attempts WHERE identifier = 'guest'`)
		g := e.call(t, "POST", "/v1/app/auth/guest", nil, "").want(200)
		e.call(t, "GET", "/v1/app/tours", nil, g.str("access_token")).wantErr(403, "guest_not_allowed")
		e.call(t, "POST", "/v1/seller/tours", map[string]string{"tour_key": "seller_intro"}, g.str("access_token")).wantErr(403, "guest_not_allowed")
	})
}
