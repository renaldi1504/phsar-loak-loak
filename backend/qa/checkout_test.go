package qa

import (
	"context"
	"io"
	"log/slog"
	"strings"
	"testing"

	"phsar/internal/httpx"
)

// ---------------------------------------------------------------------------
// 10. Keranjang, checkout KHQR, pembayaran, kedaluwarsa
// ---------------------------------------------------------------------------

type sellable struct {
	token, storeID, productID string
	variants                  []string // id varian sesuai urutan dibuat
}

// newSellable: toko aktif + produk disetujui dengan varian (harga, stok) yang diberikan.
func newSellable(t *testing.T, e *env, admin string, name string, prices []string, stocks []int) sellable {
	t.Helper()
	tok, sid := newActiveStore(t, e)
	v := e.upload(t, tok, "product_media", mp4Bytes).want(201)
	f := e.upload(t, tok, "product_media", jpegBytes).want(201)
	vals, vars := []any{}, []any{}
	for i := range prices {
		size := []string{"S", "M", "L", "XL"}[i]
		vals = append(vals, map[string]any{"value": size})
		vars = append(vars, map[string]any{"options": []string{size}, "price": prices[i], "stock_qty": stocks[i]})
	}
	p := e.call(t, "POST", "/v1/seller/stores/"+sid+"/products", map[string]any{
		"category_id": "womens-tops", "source_lang": "en", "name": name, "description": "Checkout test product.", "condition": "new",
		"options":  []any{map[string]any{"kind": "size", "name": "Size", "values": vals}},
		"variants": vars, "submit": true,
		"media": []any{map[string]any{"key": v.str("key"), "kind": "video"}, map[string]any{"key": f.str("key"), "kind": "photo"}},
	}, tok).want(201)
	pid := p.str("id")
	e.call(t, "POST", "/admin/api/products/"+pid+"/decision", map[string]string{"decision": "approve"}, admin).want(200)
	s := sellable{token: tok, storeID: sid, productID: pid}
	byOpt := map[string]string{}
	for i := 0; i < p.len("variants"); i++ {
		byOpt[p.str("variants."+itoa(i)+".options.0")] = p.str("variants." + itoa(i) + ".id")
	}
	for i := range prices {
		s.variants = append(s.variants, byOpt[[]string{"S", "M", "L", "XL"}[i]])
	}
	return s
}

var sihanoukville = map[string]any{"recipient": "Sokha", "phone": "012 345 678", "address_line": "Street 7, Sangkat 3, Sihanoukville",
	"note": "Blue gate", "lat": 10.6105, "lng": 103.5316, "google_place_id": "ChIJqa-buyer"}

func checkoutBody(key, total string, items ...any) map[string]any {
	return map[string]any{"idempotency_key": key, "items": items, "address": sihanoukville, "expected_total": total, "save_address": true}
}

func item(variant string, qty int) map[string]any {
	return map[string]any{"variant_id": variant, "qty": qty}
}

func TestCheckout(t *testing.T) {
	e := getEnv(t)
	ctx := context.Background()
	admin := newAdmin(t, e, "moderator")
	a := newSellable(t, e, admin, "Checkout tee A", []string{"5.50", "7.00"}, []int{3, 1})
	b := newSellable(t, e, admin, "Checkout cap B", []string{"4.25"}, []int{10})
	buyer, _ := e.login(t)
	fake := httpx.FakePayments()
	shop := httpx.NewShop(e.cfg, slog.New(slog.NewTextHandler(io.Discard, nil)), e.pool)
	stock := func(variant string) string {
		return e.scalar(t, `SELECT stock_qty || '/' || reserved_qty FROM product_variants WHERE id = $1`, variant)
	}

	t.Run("hitung keranjang: dipecah per toko, harga dari server", func(t *testing.T) {
		q := e.call(t, "POST", "/v1/shop/cart/quote", map[string]any{"items": []any{item(a.variants[0], 2), item(b.variants[0], 1), item(a.variants[1], 1)}}, "").want(200)
		q.eq("total", "22.25").eq("shipping_fee", "0.00").eq("item_count", "4").eq("has_problems", "false")
		check(t, q.len("stores") == 2, "harus 2 toko, dapat %d", q.len("stores"))
		q.eq("stores.0.subtotal", "18.00").eq("stores.1.subtotal", "4.25").eq("stores.0.lines.0.name", "Checkout tee A")
	})
	t.Run("keranjang: stok kurang & barang tidak tayang ditandai", func(t *testing.T) {
		q := e.call(t, "POST", "/v1/shop/cart/quote", map[string]any{"items": []any{item(a.variants[1], 5)}}, "").want(200)
		q.eq("has_problems", "true").eq("stores.0.lines.0.problem", "qty_reduced").eq("stores.0.lines.0.qty", "1").eq("total", "7.00")
		e.call(t, "POST", "/v1/shop/cart/quote", map[string]any{"items": []any{item(a.variants[0], 100)}}, "").wantErr(422, "validation", "items[0].qty")
	})

	t.Run("checkout wajib login", func(t *testing.T) {
		e.call(t, "POST", "/v1/shop/checkouts", checkoutBody("k-nologin-1", "5.50", item(a.variants[0], 1)), "").want(401)
	})
	t.Run("alamat di luar Sihanoukville ditolak", func(t *testing.T) {
		body := checkoutBody("k-pp-000001", "5.50", item(a.variants[0], 1))
		addr := map[string]any{}
		for k, v := range sihanoukville {
			addr[k] = v
		}
		addr["lat"], addr["lng"] = 11.5564, 104.9282 // Phnom Penh
		body["address"] = addr
		e.call(t, "POST", "/v1/shop/checkouts", body, buyer).wantErr(422, "not_served", "address.address_line")
		addr["lat"], addr["phone"] = 10.6105, "abc"
		e.call(t, "POST", "/v1/shop/checkouts", body, buyer).wantErr(422, "validation", "address.phone")
	})
	t.Run("harga berubah sejak dilihat → 409 dengan keranjang terbaru", func(t *testing.T) {
		r := e.call(t, "POST", "/v1/shop/checkouts", checkoutBody("k-price-0001", "5.00", item(a.variants[0], 1)), buyer).wantErr(409, "cart_changed")
		check(t, r.str("quote.total") == "5.50", "quote terbaru tidak dikirim: %s", r.Raw)
	})

	var co resp
	t.Run("checkout: 2 pesanan, stok ditahan, QR dibuat", func(t *testing.T) {
		co = e.call(t, "POST", "/v1/shop/checkouts", checkoutBody("k-main-00001", "22.25", item(a.variants[0], 2), item(a.variants[1], 1), item(b.variants[0], 1)), buyer).want(201)
		co.eq("status", "pending").eq("total", "22.25").eq("provider", "fake")
		check(t, co.str("qr_string") != "" && co.str("seconds_left") > "0", "QR / sisa waktu kosong: %s", co.Raw)
		check(t, co.len("orders") == 2, "pesanan harus 2, dapat %d", co.len("orders"))
		check(t, strings.HasPrefix(co.str("orders.0.order_no"), "PH"), "nomor pesanan: %s", co.str("orders.0.order_no"))
		co.eq("orders.0.status", "pending_payment")
		check(t, stock(a.variants[0]) == "3/2" && stock(a.variants[1]) == "1/1" && stock(b.variants[0]) == "10/1", "stok tidak ditahan: %s %s %s",
			stock(a.variants[0]), stock(a.variants[1]), stock(b.variants[0]))
		comm := e.scalar(t, `SELECT commission_amount || '/' || seller_amount FROM orders WHERE checkout_id = $1 AND store_id = $2`, co.str("id"), a.storeID)
		check(t, comm == "1.44/16.56", "komisi 8%% dari $18.00 salah: %s", comm)
	})
	t.Run("klik dua kali = checkout yang sama", func(t *testing.T) {
		again := e.call(t, "POST", "/v1/shop/checkouts", checkoutBody("k-main-00001", "22.25", item(a.variants[0], 2), item(a.variants[1], 1), item(b.variants[0], 1)), buyer).want(201)
		again.eq("id", co.str("id"))
		check(t, e.scalar(t, `SELECT count(*) FROM orders WHERE checkout_id = $1`, co.str("id")) == "2", "pesanan dobel")
	})
	t.Run("barang terakhir yang sedang ditahan tidak bisa dibeli orang lain", func(t *testing.T) {
		other, _ := e.login(t)
		e.call(t, "POST", "/v1/shop/checkouts", checkoutBody("k-other-0001", "7.00", item(a.variants[1], 1)), other).wantErr(409, "cart_changed")
		e.call(t, "GET", "/v1/shop/checkouts/"+co.str("id"), nil, other).want(404)
	})
	t.Run("callback tanpa tanda tangan tidak membuat lunas", func(t *testing.T) {
		tran := e.scalar(t, `SELECT tran_id FROM checkouts WHERE id = $1`, co.str("id"))
		e.call(t, "POST", "/webhooks/payway", map[string]any{"tran_id": tran, "status": 0, "apv": "999"}, "").want(200)
		e.call(t, "GET", "/v1/shop/checkouts/"+co.str("id"), nil, buyer).want(200).eq("status", "pending")
		check(t, e.scalar(t, `SELECT count(*) FROM payment_events WHERE checkout_id = $1 AND source = 'callback' AND NOT verified`, co.str("id")) == "1",
			"callback tidak dicatat")
	})
	t.Run("dibayar: pesanan masuk ke seller, stok terjual, notifikasi", func(t *testing.T) {
		tran := e.scalar(t, `SELECT tran_id FROM checkouts WHERE id = $1`, co.str("id"))
		fake.Pay(tran, "22.25")
		e.call(t, "POST", "/webhooks/payway", map[string]any{"tran_id": tran, "status": 0}, "").want(200)
		r := e.call(t, "GET", "/v1/shop/checkouts/"+co.str("id"), nil, buyer).want(200)
		r.eq("status", "paid").eq("orders.0.status", "paid").eq("orders.1.status", "paid")
		check(t, r.str("qr_string") == "", "QR masih ditampilkan setelah lunas")
		check(t, stock(a.variants[0]) == "1/0" && stock(a.variants[1]) == "0/0" && stock(b.variants[0]) == "9/0", "stok setelah terjual salah: %s %s %s",
			stock(a.variants[0]), stock(a.variants[1]), stock(b.variants[0]))
		check(t, e.scalar(t, `SELECT count(*) FROM notifications WHERE kind = 'order.new' AND store_id IN ($1, $2)`, a.storeID, b.storeID) == "2",
			"seller tidak dapat notifikasi pesanan baru")
		// callback kedua tidak memproses ulang
		e.call(t, "POST", "/webhooks/payway", map[string]any{"tran_id": tran}, "").want(200)
		check(t, stock(a.variants[0]) == "1/0", "stok dipotong dua kali")
	})
	t.Run("pesanan saya & alamat tersimpan", func(t *testing.T) {
		l := e.call(t, "GET", "/v1/shop/orders", nil, buyer).want(200)
		check(t, l.len("orders") == 2, "pesanan saya harus 2, dapat %d", l.len("orders"))
		no := l.str("orders.0.order_no")
		e.call(t, "GET", "/v1/shop/orders/"+no, nil, buyer).want(200).eq("status", "paid").eq("payment_type", "KHQR")
		other, _ := e.login(t)
		e.call(t, "GET", "/v1/shop/orders/"+no, nil, other).want(404)
		e.call(t, "GET", "/v1/shop/address", nil, buyer).want(200).eq("address.recipient", "Sokha").eq("address.note", "Blue gate")
	})

	t.Run("nominal tidak cocok: tidak pernah lunas", func(t *testing.T) {
		c := e.call(t, "POST", "/v1/shop/checkouts", checkoutBody("k-mismatch-1", "4.25", item(b.variants[0], 1)), buyer).want(201)
		tran := e.scalar(t, `SELECT tran_id FROM checkouts WHERE id = $1`, c.str("id"))
		fake.Pay(tran, "0.01")
		if _, err := shop.Refresh(ctx, tran, "test"); err != nil {
			t.Fatal(err)
		}
		e.call(t, "GET", "/v1/shop/checkouts/"+c.str("id"), nil, buyer).want(200).eq("status", "pending")
		check(t, e.scalar(t, `SELECT count(*) FROM payment_events WHERE checkout_id = $1 AND status = 'AMOUNT_MISMATCH'`, c.str("id")) >= "1", "mismatch tidak dicatat")
		e.call(t, "POST", "/v1/shop/checkouts/"+c.str("id")+"/cancel", nil, buyer).want(200).eq("status", "cancelled")
		check(t, stock(b.variants[0]) == "9/0", "stok tidak dilepas setelah batal: %s", stock(b.variants[0]))
	})

	t.Run("tidak dibayar sampai batas waktu → batal, stok kembali", func(t *testing.T) {
		c := e.call(t, "POST", "/v1/shop/checkouts", checkoutBody("k-expire-01", "8.50", item(b.variants[0], 2)), buyer).want(201)
		check(t, stock(b.variants[0]) == "9/2", "stok tidak ditahan: %s — %s", stock(b.variants[0]),
			e.scalar(t, `SELECT string_agg(m.reason || m.delta || ':' || c.status || ':' || o.status, ', ' ORDER BY m.id) FROM inventory_movements m
				JOIN order_items oi ON oi.id = m.order_item_id JOIN orders o ON o.id = oi.order_id JOIN checkouts c ON c.id = o.checkout_id WHERE m.variant_id = $1`, b.variants[0]))
		e.sql(t, `UPDATE checkouts SET expires_at = now() - interval '5 minutes' WHERE id = $1`, c.str("id"))
		n, err := shop.ExpireDue(ctx)
		if err != nil || n < 1 {
			t.Fatalf("ExpireDue: %d %v", n, err)
		}
		r := e.call(t, "GET", "/v1/shop/checkouts/"+c.str("id"), nil, buyer).want(200)
		r.eq("status", "expired").eq("orders.0.status", "cancelled").eq("orders.0.cancel_reason", "payment_expired")
		check(t, stock(b.variants[0]) == "9/0", "stok tidak kembali: %s", stock(b.variants[0]))
		l := e.call(t, "GET", "/v1/shop/orders", nil, buyer).want(200)
		check(t, !strings.Contains(string(l.Raw), r.str("orders.0.order_no")), "pesanan kedaluwarsa tampil di Pesanan saya")

		t.Run("dibayar terlambat (stok masih ada) → tetap diproses", func(t *testing.T) {
			tran := e.scalar(t, `SELECT tran_id FROM checkouts WHERE id = $1`, c.str("id"))
			fake.Pay(tran, "8.50")
			if _, err := shop.Refresh(ctx, tran, "test"); err != nil {
				t.Fatal(err)
			}
			e.call(t, "GET", "/v1/shop/checkouts/"+c.str("id"), nil, buyer).want(200).eq("status", "paid").eq("orders.0.status", "paid")
			check(t, stock(b.variants[0]) == "7/0", "stok setelah bayar terlambat: %s", stock(b.variants[0]))
		})
	})

	t.Run("barang habis → tampil Sold out, tidak bisa dibeli", func(t *testing.T) {
		c := e.call(t, "POST", "/v1/shop/checkouts", checkoutBody("k-soldout-1", "5.50", item(a.variants[0], 1)), buyer).want(201)
		r := e.call(t, "POST", "/dev/pay/"+c.str("id"), nil, buyer).want(200)
		r.eq("status", "paid")
		slug := e.scalar(t, `SELECT slug FROM product_translations WHERE product_id = $1 AND lang = 'en'`, a.productID)
		e.page(t, "/en/p/"+slug).has(t, "Sold out", "label habis").has(t, "OutOfStock", "JSON-LD habis")
		e.call(t, "POST", "/v1/shop/checkouts", checkoutBody("k-soldout-2", "5.50", item(a.variants[0], 1)), buyer).wantErr(409, "cart_changed")
	})
}

// ---------------------------------------------------------------------------
// 11. Seller memproses pesanan, batas 2 jam, refund manual
// ---------------------------------------------------------------------------

func TestOrderFulfilment(t *testing.T) {
	e := getEnv(t)
	ctx := context.Background()
	mod := newAdmin(t, e, "moderator")
	s := newSellable(t, e, mod, "Fulfil shirt", []string{"10.00"}, []int{20})
	buyer, _ := e.login(t)
	shop := httpx.NewShop(e.cfg, slog.New(slog.NewTextHandler(io.Discard, nil)), e.pool)
	base := "/v1/seller/stores/" + s.storeID + "/orders"
	stock := func() string {
		return e.scalar(t, `SELECT stock_qty || '/' || reserved_qty FROM product_variants WHERE id = $1`, s.variants[0])
	}
	n := 0
	order := func(qty int) string { // checkout + bayar → nomor pesanan
		t.Helper()
		n++
		total := itoa(10*qty) + ".00"
		c := e.call(t, "POST", "/v1/shop/checkouts", checkoutBody("k-fulfil-"+itoa(n)+"-x", total, item(s.variants[0], qty)), buyer).want(201)
		e.call(t, "POST", "/dev/pay/"+c.str("id"), nil, buyer).want(200).eq("status", "paid")
		return c.str("orders.0.order_no")
	}

	no := order(2)
	t.Run("pesanan baru muncul di tab Baru, tanpa data pribadi pembeli", func(t *testing.T) {
		l := e.call(t, "GET", base+"?filter=new", nil, s.token).want(200)
		l.eq("counts.new", "1").eq("orders.0.order_no", no).eq("orders.0.buyer_name", "Sokha").eq("orders.0.seller_amount", "18.40").eq("orders.0.commission_amount", "1.60")
		check(t, !strings.Contains(string(l.Raw), "012 345 678") && !strings.Contains(string(l.Raw), "Street 7"), "seller melihat HP/alamat pembeli")
		check(t, l.str("orders.0.respond_by") != "", "batas waktu respons kosong")
		other, otherStore := newActiveStore(t, e)
		e.call(t, "GET", "/v1/seller/stores/"+otherStore+"/orders/"+no, nil, other).want(404)
		e.call(t, "GET", base+"/"+no, nil, other).wantErr(404, "store_not_found")
	})
	t.Run("terima → siap dijemput (butuh titik penjemputan)", func(t *testing.T) {
		e.call(t, "POST", base+"/"+no+"/ready", nil, s.token).wantErr(409, "bad_status")
		e.call(t, "POST", base+"/"+no+"/accept", nil, s.token).want(200).eq("status", "accepted")
		e.call(t, "POST", base+"/"+no+"/accept", nil, s.token).want(200).eq("status", "accepted") // klik dua kali
		check(t, e.scalar(t, `SELECT count(*) FROM notifications WHERE kind = 'buyer.order_accepted' AND link LIKE '%'||$1`, no) == "1", "pembeli tidak diberi tahu")
		e.call(t, "POST", base+"/"+no+"/ready", nil, s.token).wantErr(409, "no_pickup_point")
		e.sql(t, `INSERT INTO store_pickup_points (store_id, label, address_line, lat, lng, contact_name, contact_phone, is_default, verified_at)
			VALUES ($1, 'Shop front', 'Street 2', 10.61, 103.53, 'Dara', '012000111', true, now())`, s.storeID)
		r := e.call(t, "POST", base+"/"+no+"/ready", nil, s.token).want(200)
		r.eq("status", "ready_for_pickup").eq("pickup_point", "Shop front")
		check(t, r.len("events") == 3, "riwayat harus paid, accepted, ready: %s", r.Raw)
		e.call(t, "POST", base+"/"+no+"/reject", map[string]string{"reason": "cannot_fulfil"}, s.token).wantErr(409, "bad_status")
		e.call(t, "GET", base+"?filter=ready", nil, s.token).want(200).eq("counts.ready", "1").eq("orders.0.order_no", no)
		e.call(t, "GET", "/v1/shop/orders/"+no, nil, buyer).want(200).eq("status", "ready_for_pickup")
	})

	var rejected, timedOut string
	t.Run("tolak: stok habis → stok tidak dikembalikan; alasan lain → stok kembali", func(t *testing.T) {
		before := stock()
		rejected = order(1)
		e.call(t, "POST", base+"/"+rejected+"/reject", map[string]string{"reason": "other"}, s.token).wantErr(422, "validation", "note")
		e.call(t, "POST", base+"/"+rejected+"/reject", map[string]string{"reason": "out_of_stock"}, s.token).want(200).eq("status", "cancelled")
		check(t, stock() != before, "stok seharusnya tetap berkurang (stok seller salah): %s → %s", before, stock())
		r := e.call(t, "GET", "/v1/shop/orders/"+rejected, nil, buyer).want(200)
		r.eq("status", "cancelled").eq("refund_status", "pending").eq("refund_amount", "10.00")
		check(t, e.scalar(t, `SELECT count(*) FROM notifications WHERE kind = 'buyer.order_cancelled' AND link LIKE '%'||$1`, rejected) == "1", "pembeli tidak diberi tahu pembatalan")

		before = stock()
		other := order(3)
		e.call(t, "POST", base+"/"+other+"/reject", map[string]any{"reason": "cannot_fulfil", "note": "Machine broken"}, s.token).want(200)
		check(t, stock() == before, "stok tidak kembali setelah ditolak: %s → %s", before, stock())
	})
	t.Run("tidak diterima dalam 2 jam → batal otomatis + refund", func(t *testing.T) {
		timedOut = order(1)
		e.sql(t, `UPDATE orders SET respond_by = now() - interval '1 minute' WHERE order_no = $1`, timedOut)
		if n, err := shop.ExpireUnaccepted(ctx); err != nil || n != 1 {
			t.Fatalf("ExpireUnaccepted: %d %v", n, err)
		}
		e.call(t, "GET", base+"/"+timedOut, nil, s.token).want(200).eq("status", "cancelled").eq("cancel_reason", "seller_timeout")
		check(t, e.scalar(t, `SELECT count(*) FROM notifications WHERE kind = 'order.seller_timeout' AND store_id = $1`, s.storeID) >= "1", "seller tidak diberi tahu")
		e.call(t, "POST", base+"/"+timedOut+"/accept", nil, s.token).wantErr(409, "bad_status")
	})
	t.Run("admin: refund manual dengan nomor referensi PayWay", func(t *testing.T) {
		e.call(t, "GET", "/admin/api/refunds", nil, mod).want(403)
		fin := newAdmin(t, e, "finance")
		l := e.call(t, "GET", "/admin/api/refunds", nil, fin).want(200)
		check(t, strings.Contains(string(l.Raw), rejected) && strings.Contains(string(l.Raw), timedOut), "pesanan tidak ada di daftar refund: %s", l.Raw)
		check(t, strings.Contains(string(l.Raw), `"tran_id":"P`), "tran_id PayWay tidak ditampilkan")
		oid := e.scalar(t, `SELECT id::text FROM orders WHERE order_no = $1`, rejected)
		e.call(t, "POST", "/admin/api/refunds/"+oid+"/done", map[string]string{"refund_ref": ""}, fin).wantErr(422, "validation", "refund_ref")
		e.call(t, "POST", "/admin/api/refunds/"+oid+"/done", map[string]string{"refund_ref": "RF-123456", "note": "Refunded via PayWay dashboard"}, fin).want(200)
		e.call(t, "POST", "/admin/api/refunds/"+oid+"/done", map[string]string{"refund_ref": "RF-123456"}, fin).wantErr(409, "bad_status")
		r := e.call(t, "GET", "/v1/shop/orders/"+rejected, nil, buyer).want(200)
		r.eq("status", "refunded").eq("refund_status", "refunded")
		check(t, strings.Contains(string(r.Raw), `"status":"refunded"`) && r.len("events") >= 4, "riwayat refund tidak tercatat: %s", r.Raw)
		check(t, e.scalar(t, `SELECT count(*) FROM notifications WHERE kind = 'buyer.order_refunded' AND link LIKE '%'||$1`, rejected) == "1", "pembeli tidak diberi tahu refund")
		done := e.call(t, "GET", "/admin/api/refunds?status=refunded", nil, fin).want(200)
		check(t, strings.Contains(string(done.Raw), "RF-123456"), "riwayat refund tidak tampil")
	})
	t.Run("notifikasi pembeli tidak muncul di lonceng seller", func(t *testing.T) {
		r := e.call(t, "GET", "/v1/seller/notifications", nil, buyer).want(200)
		check(t, !strings.Contains(string(r.Raw), "buyer."), "notifikasi pembeli bocor ke lonceng seller: %s", r.Raw)
	})
}
