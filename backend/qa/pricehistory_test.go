package qa

import (
	"strings"
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// pricehistory: riwayat harga varian (trigger) + harga terendah 30 hari + aturan harga flash sale.
// ---------------------------------------------------------------------------

func TestPriceHistory(t *testing.T) {
	e := getEnv(t)
	sa := newAdmin(t, e, "super_admin")
	mod := newAdmin(t, e, "moderator")
	fin := newAdmin(t, e, "finance")
	s := newSellable(t, e, mod, "PH rice cooker", []string{"10.00", "20.00"}, []int{20, 20})
	v0, v1 := s.variants[0], s.variants[1]
	t.Cleanup(func() {
		e.sql(t, `UPDATE platform_settings SET value = '{"low30_enabled": true}' WHERE key = 'flash_rules'`)
	})
	rows := func(v string) string {
		return e.scalar(t, `SELECT count(*) FROM variant_price_history WHERE variant_id = $1`, v)
	}
	open := func(v string) string {
		return e.scalar(t, `SELECT price || '|' || coalesce(compare_at_price::text, '-') FROM variant_price_history WHERE variant_id = $1 AND valid_to IS NULL`, v)
	}
	low := func(v, at string) string {
		return e.scalar(t, `SELECT variant_low_price($1::uuid, `+at+`)::text`, v)
	}
	patch := func(p0, c0 string, stock int) {
		e.call(t, "PATCH", "/v1/seller/stores/"+s.storeID+"/products/"+s.productID+"/variants", map[string]any{"variants": []any{
			map[string]any{"id": v0, "price": p0, "compare_at_price": c0, "stock_qty": stock, "status": "active"},
			map[string]any{"id": v1, "price": "20.00", "stock_qty": 20, "status": "active"},
		}}, s.token).want(200)
	}

	t.Run("trigger mencatat setiap perubahan harga", func(t *testing.T) {
		check(t, rows(v0) == "1" && open(v0) == "10.00|-", "varian baru: satu baris terbuka, dapat %s / %s", rows(v0), open(v0))
		patch("12.00", "15.00", 20)
		check(t, rows(v0) == "2" && open(v0) == "12.00|15.00", "PATCH harga → baris baru: %s / %s", rows(v0), open(v0))
		check(t, e.scalar(t, `SELECT count(*) FROM variant_price_history WHERE variant_id = $1 AND valid_to IS NOT NULL AND price = 10`, v0) == "1", "harga lama ditutup")
		patch("12.00", "15.00", 7) // hanya stok berubah
		check(t, rows(v0) == "2", "ubah stok tidak menambah riwayat: %s", rows(v0))
		// diubah dua kali dalam satu transaksi → satu baris (tanpa baris berdurasi nol)
		e.sql(t, `DO $$ BEGIN UPDATE product_variants SET price = 13 WHERE id = '`+v0+`'; UPDATE product_variants SET price = 11.50 WHERE id = '`+v0+`'; END $$`)
		check(t, rows(v0) == "3" && open(v0) == "11.50|15.00", "dua perubahan dalam satu transaksi: %s / %s", rows(v0), open(v0))
		check(t, e.scalar(t, `SELECT count(*) FROM variant_price_history WHERE variant_id = $1 AND valid_to IS NULL`, v0) == "1", "hanya satu harga terbuka")
		patch("10.00", "", 20) // kembali ke semula
	})

	t.Run("harga terendah 30 hari (waktu terkendali)", func(t *testing.T) {
		// v1: $15 (60–40 hari lalu) → $18 (40–10 hari lalu) → $20 (sejak 10 hari lalu)
		e.sql(t, `DELETE FROM variant_price_history WHERE variant_id = $1`, v1)
		e.sql(t, `INSERT INTO variant_price_history (variant_id, price, valid_from, valid_to) VALUES
			($1, 15, now() - interval '60 days', now() - interval '40 days'),
			($1, 18, now() - interval '40 days', now() - interval '10 days'),
			($1, 20, now() - interval '10 days', NULL)`, v1)
		check(t, low(v1, "now()") == "18.00", "30 hari terakhir: $18, dapat %s", low(v1, "now()"))
		check(t, low(v1, "now() - interval '35 days'") == "15.00", "35 hari lalu: $15, dapat %s", low(v1, "now() - interval '35 days'"))
		check(t, low(v1, "now() + interval '25 days'") == "20.00", "25 hari lagi: $18 sudah di luar jendela, dapat %s", low(v1, "now() + interval '25 days'"))
		check(t, low(v1, "now() - interval '90 days'") == "15.00", "sebelum varian ada → harga pertama, dapat %s", low(v1, "now() - interval '90 days'"))
	})

	iso := func(d time.Duration) string { return time.Now().UTC().Add(d).Format(time.RFC3339) }
	session := func(name string, start, end time.Duration, open bool) string {
		return e.call(t, "POST", "/admin/api/promos/flash", map[string]any{"name_en": name, "starts_at": iso(start), "ends_at": iso(end),
			"phsar_share_pct": 100, "min_discount_pct": 10, "per_buyer": 2, "published": true, "open_for_sellers": open}, fin).want(200).str("session.id")
	}
	fi := func(v, price string) map[string]any {
		return map[string]any{"variant_id": v, "sale_price": price, "qty_cap": 2}
	}

	var itemID string
	t.Run("flash: harga harus di bawah harga terendah 30 hari (admin & seller)", func(t *testing.T) {
		sid := session("PH evening", time.Hour, 3*time.Hour, true)
		// potongan 10% dari $20 mengizinkan $18.00, tapi harga terendah 30 hari = $18 → ditolak
		r := e.call(t, "POST", "/admin/api/promos/flash/"+sid+"/items", map[string]any{"items": []any{fi(v0, "8.00"), fi(v1, "18.00")}}, fin).
			wantErr(422, "validation", "items.1.sale_price")
		check(t, strings.Contains(r.str("error.message"), "18.00") && strings.Contains(r.str("error.message"), "17.99"), "pesan menyebut harga terendah: %s", r.str("error.message"))
		km := e.call(t, "POST", "/admin/api/promos/flash/"+sid+"/items", map[string]any{"items": []any{fi(v1, "18.00")}}, fin, "X-Lang", "km").wantErr(422, "validation", "items.0.sale_price")
		check(t, strings.Contains(km.str("error.message"), "30 ថ្ងៃ"), "pesan Khmer: %s", km.str("error.message"))
		zh := e.call(t, "POST", "/admin/api/promos/flash/"+sid+"/items", map[string]any{"items": []any{fi(v1, "18.00")}}, fin, "X-Lang", "zh").wantErr(422, "validation", "items.0.sale_price")
		check(t, strings.Contains(zh.str("error.message"), "30 天"), "pesan Mandarin: %s", zh.str("error.message"))
		a := e.call(t, "POST", "/admin/api/promos/flash/"+sid+"/items", map[string]any{"items": []any{fi(v1, "17.99")}}, fin).want(200)
		itemID = a.str("items.0.id")
		// ubah harga barang: aturan yang sama
		e.call(t, "PUT", "/admin/api/promos/flash/items/"+itemID, map[string]any{"variant_id": v1, "sale_price": "18.00", "qty_cap": 2}, fin).wantErr(422, "validation", "sale_price")

		// pemilih produk menampilkan harga terendah 30 hari (relatif ke jam mulai sesi)
		pv := e.call(t, "GET", "/admin/api/promos/variants?q=PH+rice&at="+iso(time.Hour), nil, fin).want(200)
		check(t, pv.str("low30_rule") == "true" && strings.Contains(string(pv.Raw), `"low30":"18.00"`), "picker admin: %s", truncate(string(pv.Raw)))
		sv := e.call(t, "GET", "/v1/seller/stores/"+s.storeID+"/flash/variants?at="+iso(time.Hour), nil, s.token).want(200)
		check(t, sv.str("low30_rule") == "true" && strings.Contains(string(sv.Raw), `"low30":"18.00"`) && strings.Contains(string(sv.Raw), `"low30":"10.00"`), "picker seller: %s", truncate(string(sv.Raw)))

		// seller mendaftar di sesi lain
		sid2 := session("PH weekend", 4*time.Hour, 6*time.Hour, true)
		base := "/v1/seller/stores/" + s.storeID + "/flash/sessions/" + sid2 + "/items"
		e.call(t, "POST", base, map[string]any{"items": []any{fi(v1, "18.00")}}, s.token).wantErr(422, "validation", "items.0.sale_price")
		e.call(t, "POST", base, map[string]any{"items": []any{fi(v1, "17.50")}}, s.token).want(200)
	})

	t.Run("sakelar aturan: hanya super admin", func(t *testing.T) {
		e.call(t, "GET", "/admin/api/flash-rules", nil, mod).want(403)
		e.call(t, "PUT", "/admin/api/flash-rules", map[string]any{"low30_enabled": false}, fin).want(403)
		e.call(t, "GET", "/admin/api/flash-rules", nil, sa).want(200).eq("low30_enabled", "true")
		e.call(t, "PUT", "/admin/api/flash-rules", map[string]any{}, sa).wantErr(422, "validation", "low30_enabled")
		e.call(t, "PUT", "/admin/api/flash-rules", map[string]any{"low30_enabled": false, "x": 1}, sa).want(400)
		e.call(t, "PUT", "/admin/api/flash-rules", map[string]any{"low30_enabled": false}, sa).want(200).eq("low30_enabled", "false")
		// aturan mati: harga = harga terendah 30 hari boleh (potongan minimal sesi tetap berlaku)
		e.call(t, "PUT", "/admin/api/promos/flash/items/"+itemID, map[string]any{"variant_id": v1, "sale_price": "18.00", "qty_cap": 2}, fin).want(200)
		e.call(t, "PUT", "/admin/api/promos/flash/items/"+itemID, map[string]any{"variant_id": v1, "sale_price": "18.01", "qty_cap": 2}, fin).wantErr(422, "validation", "sale_price")
		check(t, e.call(t, "GET", "/admin/api/promos/variants?q=PH+rice", nil, fin).want(200).str("low30_rule") == "false", "picker tahu aturan mati")
		e.call(t, "PUT", "/admin/api/flash-rules", map[string]any{"low30_enabled": true}, sa).want(200).eq("low30_enabled", "true")
		e.call(t, "PUT", "/admin/api/promos/flash/items/"+itemID, map[string]any{"variant_id": v1, "sale_price": "17.00", "qty_cap": 2}, fin).want(200)
	})

	t.Run("website: harga terendah 30 hari di halaman produk saat flash", func(t *testing.T) {
		slug := e.scalar(t, `SELECT slug FROM product_translations WHERE product_id = $1 AND lang = 'en'`, s.productID)
		before := string(e.page(t, "/en/p/"+slug).Raw)
		check(t, !strings.Contains(before, "ph-low30"), "tanpa flash berjalan: tidak ditampilkan")
		live := session("PH now", -time.Minute, 30*time.Minute, false)
		e.call(t, "POST", "/admin/api/promos/flash/"+live+"/items", map[string]any{"items": []any{fi(v1, "16.00")}}, fin).want(200)
		pg := string(e.page(t, "/en/p/"+slug).Raw)
		check(t, strings.Contains(pg, "Lowest price in the 30 days before this sale: $18.00"), "halaman produk en")
		zh := string(e.page(t, "/zh/p/"+slug).Raw) // slug bahasa mana pun dikenali
		check(t, strings.Contains(zh, "本次抢购前 30 天最低价：$18.00"), "halaman produk zh")
		e.call(t, "POST", "/admin/api/promos/flash/"+live+"/cancel", nil, fin).want(200)
	})
}
