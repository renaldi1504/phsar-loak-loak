package qa

import (
	"strings"
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// Flash sale berbasis sesi: banyak produk, pembiayaan per sesi, seller mendaftar / menerima undangan.
// ---------------------------------------------------------------------------

func TestFlashSessions(t *testing.T) {
	e := getEnv(t)
	mod := newAdmin(t, e, "moderator")
	fin := newAdmin(t, e, "finance")
	a := newSellable(t, e, mod, "Flash kettle", []string{"10.00"}, []int{20})
	b := newSellable(t, e, mod, "Flash fan", []string{"8.00"}, []int{20})
	c := newSellable(t, e, mod, "Plain lamp", []string{"20.00"}, []int{20})
	buyer, _ := e.login(t)
	buyer2, _ := e.login(t)
	now := time.Now().UTC()
	iso := func(d time.Duration) string { return now.Add(d).Format(time.RFC3339) }
	n := 0
	key := func() string { n++; return "k-flash-" + itoa(n) + "-xyz" }
	quote := func(code string, items ...any) resp {
		return e.call(t, "POST", "/v1/shop/cart/quote", map[string]any{"items": items, "voucher_code": code}, "").want(200)
	}
	session := func(body map[string]any) resp {
		base := map[string]any{"name_en": "Midnight sale", "name_km": "លក់កណ្តាលអធ្រាត្រ", "starts_at": iso(-time.Minute), "ends_at": iso(time.Hour),
			"phsar_share_pct": 100, "allow_vouchers": false, "min_discount_pct": 10, "per_buyer": 2, "published": true}
		for k, v := range body {
			base[k] = v
		}
		return e.call(t, "POST", "/admin/api/promos/flash", base, fin)
	}
	add := func(sid string, items ...map[string]any) resp {
		return e.call(t, "POST", "/admin/api/promos/flash/"+sid+"/items", map[string]any{"items": items}, fin)
	}
	fi := func(v, price string, qty int) map[string]any {
		return map[string]any{"variant_id": v, "sale_price": price, "qty_cap": qty}
	}

	var live string
	t.Run("sesi: izin & validasi", func(t *testing.T) {
		e.call(t, "GET", "/admin/api/promos/flash", nil, mod).want(403)
		session(map[string]any{"name_en": ""}).wantErr(422, "validation", "name_en")
		session(map[string]any{"ends_at": iso(-2 * time.Minute)}).wantErr(422, "validation", "ends_at")
		session(map[string]any{"ends_at": iso(15 * 24 * time.Hour)}).wantErr(422, "validation", "ends_at")
		session(map[string]any{"phsar_share_pct": 120}).wantErr(422, "validation", "phsar_share_pct")
		r := session(nil).want(200)
		r.eq("session.status", "live").eq("session.phsar_share_pct", "100")
		live = r.str("session.id")
	})

	t.Run("admin memilih banyak produk sekaligus (Phsar menanggung → langsung disetujui)", func(t *testing.T) {
		add(live, fi(a.variants[0], "9.50", 3)).wantErr(422, "validation", "items.0.sale_price") // potongan < 10%
		add(live, fi(a.variants[0], "8.00", 3), fi(b.variants[0], "6.00", 50)).wantErr(422, "validation", "items.1.qty_cap")
		r := add(live, fi(a.variants[0], "8.00", 3), fi(b.variants[0], "6.00", 5)).want(200)
		check(t, r.len("items") == 2 && r.str("items.0.status") == "approved" && r.str("items.1.status") == "approved", "dua barang harus disetujui: %s", truncate(string(r.Raw)))
		other := session(map[string]any{"name_en": "Overlap", "starts_at": iso(30 * time.Minute), "ends_at": iso(2 * time.Hour)}).want(200).str("session.id")
		add(other, fi(a.variants[0], "8.00", 2)).wantErr(422, "validation", "items.0.variant_id") // waktu bertabrakan
		e.call(t, "POST", "/admin/api/promos/flash/"+other+"/cancel", nil, fin).want(200).eq("session.status", "cancelled")
	})

	t.Run("harga flash, batas per pembeli, subsidi Phsar masuk penghasilan seller", func(t *testing.T) {
		q := quote("", item(a.variants[0], 1))
		q.eq("total", "8.00").eq("stores.0.lines.0.compare", "10.00").eq("stores.0.lines.0.flash_per_buyer", "2")
		quote("", item(a.variants[0], 3)).eq("total", "30.00") // melebihi batas per pembeli → harga normal
		co := e.call(t, "POST", "/v1/shop/checkouts", checkoutBody(key(), "16.00", item(a.variants[0], 2)), buyer).want(201)
		// $16 dibayar pembeli; Phsar menanggung (10−8)×2 = $4; komisi 8% dari $20 = $1.60; seller $18.40
		row := e.scalar(t, `SELECT subtotal || '|' || flash_subsidy || '|' || commission_amount || '|' || seller_amount FROM orders WHERE checkout_id = $1`, co.str("id"))
		check(t, row == "16.00|4.00|1.60|18.40", "uang pesanan flash: %s", row)
		check(t, e.scalar(t, `SELECT count(*) FROM order_items oi JOIN orders o ON o.id = oi.order_id WHERE o.checkout_id = $1 AND oi.flash_item_id IS NOT NULL AND oi.flash_subsidy = 4`,
			co.str("id")) == "1", "barang harus ditandai flash + subsidi")
		// pembeli yang sama sudah membeli 2 → harga flash tidak berlaku lagi untuknya
		e.call(t, "POST", "/v1/shop/checkouts", checkoutBody(key(), "8.00", item(a.variants[0], 1)), buyer).wantErr(409, "cart_changed")
		e.call(t, "POST", "/v1/shop/checkouts", checkoutBody(key(), "8.00", item(a.variants[0], 1)), buyer2).want(201)
		quote("", item(a.variants[0], 1)).eq("total", "10.00") // kuota 3 habis
	})

	t.Run("sesi tanpa voucher: potongan hanya untuk barang biasa", func(t *testing.T) {
		e.call(t, "POST", "/admin/api/promos/vouchers", map[string]any{"code": "FLASHV", "kind": "percent", "value": "10", "per_user_limit": 5}, fin).want(200)
		quote("FLASHV", item(b.variants[0], 1)).eq("voucher_error", "voucher_flash").eq("total", "6.00")
		quote("FLASHV", item(b.variants[0], 1), item(c.variants[0], 1)).eq("discount", "2.00").eq("total", "24.00")
	})

	t.Run("website: bagian flash sale & halaman sesi", func(t *testing.T) {
		home := string(e.page(t, "/en/").Raw)
		check(t, strings.Contains(home, `id="h-flash"`) && strings.Contains(home, "Midnight sale"), "beranda harus menampilkan flash sale")
		pg := string(e.page(t, "/km/flash-sale").Raw)
		check(t, strings.Contains(pg, "លក់កណ្តាលអធ្រាត្រ") && strings.Contains(pg, "6.00") && strings.Contains(pg, "fbar"), "halaman flash sale")
	})

	var sched string
	t.Run("seller menanggung: undangan admin harus diterima seller; seller mendaftar & admin menyetujui", func(t *testing.T) {
		sched = session(map[string]any{"name_en": "Weekend deals", "starts_at": iso(3 * time.Hour), "ends_at": iso(5 * time.Hour), "phsar_share_pct": 0,
			"allow_vouchers": true, "open_for_sellers": true, "register_until": iso(2 * time.Hour)}).want(200).str("session.id")
		r := add(sched, fi(a.variants[0], "8.50", 4)).want(200)
		r.eq("items.0.status", "invited")
		item0 := r.str("items.0.id")
		check(t, e.scalar(t, `SELECT count(*) FROM notifications WHERE store_id = $1 AND kind = 'flash.invited'`, a.storeID) != "0", "seller harus diberi tahu")
		sv := e.call(t, "GET", "/v1/seller/stores/"+a.storeID+"/flash", nil, a.token).want(200)
		check(t, strings.Contains(string(sv.Raw), item0) && strings.Contains(string(sv.Raw), `"can_register":true`), "seller melihat undangan & sesi terbuka")
		e.call(t, "GET", "/v1/seller/stores/"+a.storeID+"/flash", nil, b.token).want(404)
		e.call(t, "POST", "/v1/seller/stores/"+a.storeID+"/flash/items/"+item0+"/accept", nil, a.token).want(200)
		check(t, e.scalar(t, `SELECT status FROM flash_items WHERE id = $1`, item0) == "approved", "undangan diterima → approved")

		// seller B mendaftar produknya sendiri (bukan produk toko lain)
		base := "/v1/seller/stores/" + b.storeID + "/flash"
		check(t, e.call(t, "GET", base+"/variants", nil, b.token).want(200).len("variants") >= 1, "varian seller")
		e.call(t, "POST", base+"/sessions/"+sched+"/items", map[string]any{"items": []any{fi(a.variants[0], "8.00", 1)}}, b.token).wantErr(422, "validation", "items.0.variant_id")
		e.call(t, "POST", base+"/sessions/"+sched+"/items", map[string]any{"items": []any{fi(b.variants[0], "7.50", 2)}}, b.token).wantErr(422, "validation", "items.0.sale_price")
		e.call(t, "POST", base+"/sessions/"+sched+"/items", map[string]any{"items": []any{fi(b.variants[0], "7.00", 2)}}, b.token).want(200)
		bid := e.scalar(t, `SELECT id::text FROM flash_items WHERE session_id = $1 AND store_id = $2 AND status = 'pending'`, sched, b.storeID)
		e.call(t, "POST", "/admin/api/promos/flash/items/"+bid+"/decide", map[string]any{"action": "reject"}, fin).wantErr(422, "validation", "note")
		e.call(t, "POST", "/admin/api/promos/flash/items/"+bid+"/decide", map[string]any{"action": "approve"}, fin).want(200)
		e.call(t, "POST", base+"/items/"+bid+"/withdraw", nil, b.token).want(200)
		check(t, e.scalar(t, `SELECT status FROM flash_items WHERE id = $1`, bid) == "withdrawn", "seller menarik sebelum sesi mulai")
		quote("", item(b.variants[0], 1)).eq("total", "6.00") // sesi terjadwal belum berlaku; harga sesi yang berjalan
	})

	t.Run("pembiayaan berubah → seller harus setuju lagi; sesi dibatalkan → harga normal", func(t *testing.T) {
		x := session(map[string]any{"name_en": "Phsar day", "starts_at": iso(6 * time.Hour), "ends_at": iso(7 * time.Hour), "phsar_share_pct": 100}).want(200).str("session.id")
		add(x, fi(c.variants[0], "15.00", 2)).want(200).eq("items.0.status", "approved")
		body := map[string]any{"name_en": "Phsar day", "starts_at": iso(6 * time.Hour), "ends_at": iso(7 * time.Hour), "phsar_share_pct": 50, "min_discount_pct": 10, "per_buyer": 2, "published": true}
		e.call(t, "PUT", "/admin/api/promos/flash/"+x, body, fin).want(200).eq("items.0.status", "invited")

		e.call(t, "POST", "/admin/api/promos/flash/"+live+"/cancel", nil, fin).want(200)
		quote("", item(b.variants[0], 1)).eq("total", "8.00")
	})

	t.Run("seller menanggung penuh: tanpa subsidi, penghasilan dari harga flash", func(t *testing.T) {
		e.sql(t, `UPDATE flash_sessions SET starts_at = now() - interval '1 minute', register_until = NULL WHERE id = $1`, sched)
		co := e.call(t, "POST", "/v1/shop/checkouts", checkoutBody(key(), "8.50", item(a.variants[0], 1)), buyer2).want(201)
		row := e.scalar(t, `SELECT flash_subsidy || '|' || commission_amount || '|' || seller_amount FROM orders WHERE checkout_id = $1`, co.str("id"))
		check(t, row == "0.00|0.68|7.82", "seller-funded: %s", row)
		e.call(t, "PUT", "/admin/api/promos/flash/"+sched, map[string]any{"name_en": "Weekend deals", "starts_at": iso(time.Hour), "ends_at": iso(5 * time.Hour), "phsar_share_pct": 0},
			fin).wantErr(422, "validation", "starts_at")
	})
}
