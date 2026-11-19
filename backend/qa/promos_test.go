package qa

import (
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// Voucher: potongan ditanggung Phsar, penghasilan seller tetap. (Flash sale: flash_test.go)
// ---------------------------------------------------------------------------

func TestPromos(t *testing.T) {
	e := getEnv(t)
	mod := newAdmin(t, e, "moderator")
	fin := newAdmin(t, e, "finance")
	a := newSellable(t, e, mod, "Promo shirt", []string{"6.00"}, []int{50})
	b := newSellable(t, e, mod, "Promo socks", []string{"4.00"}, []int{50})
	buyer, _ := e.login(t)
	buyer2, _ := e.login(t)
	buyer3, _ := e.login(t)
	now := time.Now().UTC()
	iso := func(d time.Duration) string { return now.Add(d).Format(time.RFC3339) }
	n := 0
	key := func() string { n++; return "k-promo-" + itoa(n) + "-xyz" }
	voucherBody := func(code string, extra map[string]any) map[string]any {
		m := map[string]any{"code": code, "kind": "percent", "value": "10", "max_discount": "3", "min_spend": "10", "total_limit": 2, "per_user_limit": 1}
		for k, v := range extra {
			m[k] = v
		}
		return m
	}
	quote := func(code string, items ...any) resp {
		return e.call(t, "POST", "/v1/shop/cart/quote", map[string]any{"items": items, "voucher_code": code}, "").want(200)
	}
	withVoucher := func(body map[string]any, code string) map[string]any { body["voucher_code"] = code; return body }

	t.Run("izin & validasi admin", func(t *testing.T) {
		e.call(t, "GET", "/admin/api/promos", nil, mod).want(403)
		e.call(t, "GET", "/admin/api/promos", nil, fin).want(200)
		e.call(t, "POST", "/admin/api/promos/vouchers", voucherBody("a", nil), fin).wantErr(422, "validation", "code")
		e.call(t, "POST", "/admin/api/promos/vouchers", voucherBody("BIG95", map[string]any{"value": "95"}), fin).wantErr(422, "validation", "value")
		e.call(t, "POST", "/admin/api/promos/vouchers", voucherBody("WHEN1", map[string]any{"starts_at": iso(time.Hour), "ends_at": iso(time.Minute)}), fin).wantErr(422, "validation", "ends_at")
		r := e.call(t, "POST", "/admin/api/promos/vouchers", voucherBody("save10", nil), fin).want(200)
		r.eq("vouchers.0.code", "SAVE10").eq("vouchers.0.state", "live")
		e.call(t, "POST", "/admin/api/promos/vouchers", voucherBody("Save10", nil), fin).wantErr(422, "validation", "code")
	})

	var co resp
	t.Run("voucher di keranjang & checkout; seller tetap menerima penuh", func(t *testing.T) {
		q := quote("save10", item(a.variants[0], 2))
		q.eq("subtotal", "12.00").eq("discount", "1.20").eq("total", "10.80").eq("voucher.code", "SAVE10")
		quote("SAVE10", item(a.variants[0], 1)).eq("voucher_error", "voucher_min_spend").eq("voucher_min_spend", "10.00").eq("total", "6.00").eq("discount", "0.00")
		quote("NOPE", item(a.variants[0], 2)).eq("voucher_error", "voucher_not_found")
		// total yang diharapkan harus total setelah potongan
		e.call(t, "POST", "/v1/shop/checkouts", withVoucher(checkoutBody(key(), "12.00", item(a.variants[0], 2)), "SAVE10"), buyer).wantErr(409, "cart_changed")
		co = e.call(t, "POST", "/v1/shop/checkouts", withVoucher(checkoutBody(key(), "10.80", item(a.variants[0], 2)), "SAVE10"), buyer).want(201)
		co.eq("total", "10.80").eq("discount", "1.20").eq("subtotal", "12.00").eq("orders.0.discount", "1.20").eq("orders.0.total", "10.80")
		row := e.scalar(t, `SELECT subtotal || '/' || total || '/' || (seller_amount + commission_amount) FROM orders WHERE checkout_id = $1`, co.str("id"))
		check(t, row == "12.00/10.80/12.00", "subtotal/total/seller+komisi salah: %s", row)
		check(t, e.scalar(t, `SELECT status FROM voucher_redemptions WHERE checkout_id = $1`, co.str("id")) == "held", "pemakaian harus held")
		// batas per pembeli
		e.call(t, "POST", "/v1/shop/checkouts", withVoucher(checkoutBody(key(), "10.80", item(a.variants[0], 2)), "SAVE10"), buyer).wantErr(422, "voucher_user_limit", "voucher_code")
	})

	t.Run("batal → kuota kembali; lunas → used; kuota total habis", func(t *testing.T) {
		e.call(t, "POST", "/v1/shop/checkouts/"+co.str("id")+"/cancel", nil, buyer).want(200)
		check(t, e.scalar(t, `SELECT status FROM voucher_redemptions WHERE checkout_id = $1`, co.str("id")) == "released", "pemakaian harus dilepas")
		c2 := e.call(t, "POST", "/v1/shop/checkouts", withVoucher(checkoutBody(key(), "10.80", item(a.variants[0], 2)), "SAVE10"), buyer).want(201)
		e.call(t, "POST", "/dev/pay/"+c2.str("id"), nil, buyer).want(200).eq("status", "paid")
		check(t, e.scalar(t, `SELECT status FROM voucher_redemptions WHERE checkout_id = $1`, c2.str("id")) == "used", "pemakaian harus used")
		e.call(t, "POST", "/v1/shop/checkouts", withVoucher(checkoutBody(key(), "10.80", item(a.variants[0], 2)), "SAVE10"), buyer2).want(201)
		e.call(t, "POST", "/v1/shop/checkouts", withVoucher(checkoutBody(key(), "10.80", item(a.variants[0], 2)), "SAVE10"), buyer3).wantErr(422, "voucher_used_up")
		quote("SAVE10", item(a.variants[0], 2)).eq("voucher_error", "voucher_used_up")
		v := e.call(t, "GET", "/admin/api/promos", nil, fin).want(200)
		v.eq("vouchers.0.used", "1").eq("vouchers.0.held", "1").eq("vouchers.0.cost", "1.20")
	})

	t.Run("potongan dibagi ke pesanan tiap toko", func(t *testing.T) {
		e.call(t, "POST", "/admin/api/promos/vouchers", map[string]any{"code": "FIVE", "kind": "fixed", "value": "5", "per_user_limit": 3}, fin).want(200)
		c := e.call(t, "POST", "/v1/shop/checkouts", withVoucher(checkoutBody(key(), "11.00", item(a.variants[0], 2), item(b.variants[0], 1)), "five"), buyer3).want(201)
		c.eq("discount", "5.00").eq("total", "11.00")
		got := e.scalar(t, `SELECT string_agg(subtotal || ':' || discount || ':' || total, ',' ORDER BY subtotal DESC) FROM orders WHERE checkout_id = $1`, c.str("id"))
		check(t, got == "12.00:3.75:8.25,4.00:1.25:2.75", "pembagian potongan salah: %s", got)
		// potongan tidak boleh membuat total 0
		qq := quote("FIVE", item(b.variants[0], 1))
		qq.eq("discount", "3.99").eq("total", "0.01")
	})

	t.Run("pesanan pertama: tamu ditolak, pembeli lama ditolak", func(t *testing.T) {
		e.call(t, "POST", "/admin/api/promos/vouchers", map[string]any{"code": "HELLO", "kind": "fixed", "value": "2", "first_order_only": true}, fin).want(200)
		guest := e.call(t, "POST", "/v1/app/auth/guest", nil, "").want(200).str("access_token")
		e.call(t, "POST", "/v1/shop/checkouts", withVoucher(checkoutBody(key(), "4.00", item(a.variants[0], 1)), "HELLO"), guest).wantErr(422, "voucher_login_required")
		e.call(t, "POST", "/v1/shop/checkouts", withVoucher(checkoutBody(key(), "4.00", item(a.variants[0], 1)), "HELLO"), buyer).wantErr(422, "voucher_first_order")
		fresh, _ := e.login(t)
		e.call(t, "POST", "/v1/shop/checkouts", withVoucher(checkoutBody(key(), "4.00", item(a.variants[0], 1)), "HELLO"), fresh).want(201).eq("total", "4.00")
	})

	t.Run("voucher belum mulai / habis / dimatikan", func(t *testing.T) {
		e.call(t, "POST", "/admin/api/promos/vouchers", map[string]any{"code": "LATER", "kind": "fixed", "value": "1", "starts_at": iso(time.Hour)}, fin).want(200)
		quote("LATER", item(a.variants[0], 2)).eq("voucher_error", "voucher_not_started")
		e.sql(t, `INSERT INTO vouchers (code, kind, value, starts_at, ends_at) VALUES ('OLD', 'fixed', 1, now() - interval '2 days', now() - interval '1 day')`)
		quote("OLD", item(a.variants[0], 2)).eq("voucher_error", "voucher_expired")
		id := e.scalar(t, `SELECT id::text FROM vouchers WHERE code = 'LATER'`)
		e.call(t, "POST", "/admin/api/promos/vouchers/"+id+"/active", map[string]any{"active": false}, fin).want(200)
		quote("LATER", item(a.variants[0], 2)).eq("voucher_error", "voucher_not_found")
	})

}
