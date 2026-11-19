package qa

import (
	"testing"
)

// ---------------------------------------------------------------------------
// Seller Center v2: layar "Hari ini" (daftar tugas) + tarif komisi di data toko.
// ---------------------------------------------------------------------------

func TestSellerToday(t *testing.T) {
	e := getEnv(t)
	mod := newAdmin(t, e, "moderator")
	s := newSellable(t, e, mod, "Today kettle", []string{"10.00"}, []int{20})
	buyer, _ := e.login(t)
	base := "/v1/seller/stores/" + s.storeID

	t.Run("toko baru: belum ada tugas", func(t *testing.T) {
		r := e.call(t, "GET", base+"/today", nil, s.token).want(200)
		r.eq("new_orders", "0").eq("to_pack", "0").eq("returns", "0").eq("sales_today", "0.00")
	})

	t.Run("pesanan dibayar muncul sebagai tugas 'terima' dengan tenggat", func(t *testing.T) {
		c := e.call(t, "POST", "/v1/shop/checkouts", checkoutBody("k-today-1-xyz", "10.00", item(s.variants[0], 1)), buyer).want(201)
		e.call(t, "POST", "/dev/pay/"+c.str("id"), nil, buyer).want(200)
		r := e.call(t, "GET", base+"/today", nil, s.token).want(200)
		r.eq("new_orders", "1").eq("orders_today", "1").eq("sales_today", "10.00")
		check(t, r.str("new_due") != "", "new_due harus diisi: %s", truncate(string(r.Raw)))
	})

	t.Run("komisi ikut di GET toko; orang lain tidak bisa melihat", func(t *testing.T) {
		r := e.call(t, "GET", base, nil, s.token).want(200)
		check(t, r.str("commission.global") != "", "commission.global kosong: %s", truncate(string(r.Raw)))
		other, _ := e.login(t)
		e.call(t, "GET", base+"/today", nil, other).want(404)
	})
}
