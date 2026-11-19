package qa

import (
	"testing"
)

// Aturan komisi dari admin: toko > kategori > default; hanya pesanan baru yang terpengaruh.
func TestCommissionRules(t *testing.T) {
	e := getEnv(t)
	mod := newAdmin(t, e, "moderator")
	fin := newAdmin(t, e, "finance")
	a := newSellable(t, e, mod, "Commission tee", []string{"10.00"}, []int{20})
	buyer, _ := e.login(t)
	order := func(key string) (string, string) {
		c := e.call(t, "POST", "/v1/shop/checkouts", checkoutBody(key, "10.00", item(a.variants[0], 1)), buyer).want(201)
		id := c.str("id")
		return e.scalar(t, `SELECT commission_rate::text FROM orders WHERE checkout_id = $1`, id), e.scalar(t, `SELECT commission_amount::text FROM orders WHERE checkout_id = $1`, id)
	}
	t.Cleanup(func() { // kembalikan default 8% dan hapus aturan tes
		e.sql(t, `UPDATE commission_rules SET valid_to = now() WHERE valid_to IS NULL AND scope <> 'global' AND created_by IS NOT NULL`)
		e.call(t, "POST", "/admin/api/commission/rules", map[string]any{"scope": "global", "rate_percent": "8"}, fin).want(200)
	})

	e.call(t, "GET", "/admin/api/commission", nil, mod).want(403)
	v := e.call(t, "GET", "/admin/api/commission", nil, fin).want(200)
	check(t, v.len("categories") >= 5, "kategori harus tampil")

	rate, comm := order("k-comm-1")
	check(t, rate == "8.00" && comm == "0.80", "default 8%%: %s / %s", rate, comm)

	// kategori fashion 5%
	e.call(t, "POST", "/admin/api/commission/rules", map[string]any{"scope": "category", "category_id": "nope", "rate_percent": "5"}, fin).wantErr(422, "validation", "category_id")
	e.call(t, "POST", "/admin/api/commission/rules", map[string]any{"scope": "category", "category_id": "fashion", "rate_percent": "60"}, fin).wantErr(422, "validation", "rate_percent")
	e.call(t, "POST", "/admin/api/commission/rules", map[string]any{"scope": "category", "category_id": "fashion", "rate_percent": "5"}, fin).want(200)
	rate, comm = order("k-comm-2")
	check(t, rate == "5.00" && comm == "0.50", "kategori 5%%: %s / %s", rate, comm)

	// toko 0% sampai bulan depan (promo seller awal) — menang atas kategori
	e.call(t, "POST", "/admin/api/commission/rules", map[string]any{"scope": "store", "store_id": a.storeID, "rate_percent": "0", "valid_to": "2020-01-01T00:00:00Z"}, fin).wantErr(422, "validation", "valid_to")
	r := e.call(t, "POST", "/admin/api/commission/rules", map[string]any{"scope": "store", "store_id": a.storeID, "rate_percent": "0", "valid_to": "2099-01-01T00:00:00Z"}, fin).want(200)
	rate, comm = order("k-comm-3")
	check(t, rate == "0.00" && comm == "0.00", "toko 0%%: %s / %s", rate, comm)

	// akhiri aturan toko → kembali ke kategori; default berubah 10% tapi kategori tetap 5%
	var storeRule string
	for i := 0; i < r.len("active"); i++ {
		if r.str("active."+itoa(i)+".scope") == "store" && r.str("active."+itoa(i)+".store_id") == a.storeID {
			storeRule = r.str("active." + itoa(i) + ".id")
		}
	}
	e.call(t, "DELETE", "/admin/api/commission/rules/"+storeRule, nil, fin).want(200)
	e.call(t, "DELETE", "/admin/api/commission/rules/"+storeRule, nil, fin).want(404)
	e.call(t, "POST", "/admin/api/commission/rules", map[string]any{"scope": "global", "rate_percent": "10"}, fin).want(200)
	rate, _ = order("k-comm-4")
	check(t, rate == "5.00", "setelah aturan toko berakhir harus kembali ke kategori: %s", rate)
	check(t, e.scalar(t, `SELECT commission_rate::text FROM orders o JOIN checkouts c ON c.id = o.checkout_id WHERE c.idempotency_key = 'k-comm-1'`) == "8.00",
		"pesanan lama tidak boleh berubah")
	check(t, e.scalar(t, `SELECT count(*) FROM commission_rules WHERE scope = 'global' AND (valid_to IS NULL OR valid_to > now())`) == "1", "hanya satu default aktif")
}
