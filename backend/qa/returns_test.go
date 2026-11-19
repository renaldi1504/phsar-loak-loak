package qa

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"testing"

	"phsar/internal/httpx"
	"phsar/internal/returns"
)

func (e *env) buyerUpload(t *testing.T, token, kind string, data []byte) resp {
	t.Helper()
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	_ = mw.WriteField("kind", kind)
	fw, _ := mw.CreateFormFile("file", "photo.png")
	_, _ = fw.Write(data)
	_ = mw.Close()
	req, _ := http.NewRequest(http.MethodPost, e.srv.URL+"/v1/shop/uploads", &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	return e.send(t, req)
}

// ---------------------------------------------------------------------------
// 16. Retur & sengketa: pembeli → seller (48 jam) → Phsar → refund + potong saldo seller
// ---------------------------------------------------------------------------

func TestReturns(t *testing.T) {
	e := getEnv(t)
	ctx := context.Background()
	quiet := slog.New(slog.NewTextHandler(io.Discard, nil))
	mod := newAdmin(t, e, "moderator")
	cs := newAdmin(t, e, "cs")
	fin := newAdmin(t, e, "finance")
	s := newSellable(t, e, mod, "Return shoes", []string{"6.00"}, []int{50})
	e.sql(t, `INSERT INTO store_pickup_points (store_id, label, address_line, lat, lng, contact_name, contact_phone, is_default, verified_at)
		VALUES ($1, 'Shoe shop', 'Street 9', 10.61, 103.53, 'Dara', '012777555', true, now())`, s.storeID)
	courier, cphone := e.login(t)
	e.call(t, "POST", "/admin/api/couriers", map[string]any{"phone": cphone, "name": "Vuth"}, cs).want(200)
	buyer, _ := e.login(t)
	other, _ := e.login(t)
	payouts := httpx.NewPayouts(e.cfg, quiet, e.pool)
	rsvc := &returns.Service{Pool: e.pool, Log: quiet}
	base := "/v1/seller/stores/" + s.storeID
	n := 0
	order := func(deliver bool) string {
		t.Helper()
		n++
		c := e.call(t, "POST", "/v1/shop/checkouts", checkoutBody("k-return-"+itoa(n)+"-x", "6.00", item(s.variants[0], 1)), buyer).want(201)
		e.call(t, "POST", "/dev/pay/"+c.str("id"), nil, buyer).want(200)
		no := c.str("orders.0.order_no")
		if !deliver {
			return no
		}
		e.call(t, "POST", base+"/orders/"+no+"/accept", nil, s.token).want(200)
		e.call(t, "POST", base+"/orders/"+no+"/ready", nil, s.token).want(200)
		task := e.scalar(t, `SELECT d.id::text FROM deliveries d JOIN orders o ON o.id = d.order_id WHERE o.order_no = $1`, no)
		e.call(t, "POST", "/v1/courier/tasks/"+task+"/claim", nil, courier).want(200)
		e.call(t, "POST", "/v1/courier/tasks/"+task+"/pickup", map[string]string{"code": no}, courier).want(200)
		e.call(t, "POST", "/v1/courier/tasks/"+task+"/deliver", nil, courier).want(200)
		return no
	}
	age := func(no string, days int) {
		e.sql(t, `UPDATE deliveries SET delivered_at = now() - make_interval(days => $1) WHERE order_id = (SELECT id FROM orders WHERE order_no = $2)`, days, no)
	}
	ask := func(no string, body map[string]any) resp {
		return e.call(t, "POST", "/v1/shop/orders/"+no+"/return", body, buyer)
	}
	balance := func() string { return e.call(t, "GET", base+"/payouts", nil, s.token).want(200).str("available") }

	t.Run("belum terantar tidak bisa retur", func(t *testing.T) {
		no := order(false)
		e.call(t, "GET", "/v1/shop/orders/"+no+"/return", nil, buyer).want(200).eq("can_request", "false").eq("not_allowed", "not_delivered")
		ask(no, map[string]any{"reason": "damaged", "description": "Box was crushed"}).wantErr(409, "return_not_allowed")
		e.call(t, "GET", "/v1/shop/orders/"+no+"/return", nil, other).want(404)
	})

	var rid, no1 string
	t.Run("rusak: pembeli ajukan dengan foto → seller tolak → Phsar refund sebagian → saldo seller dipotong", func(t *testing.T) {
		no1 = order(true)
		el := e.call(t, "GET", "/v1/shop/orders/"+no1+"/return", nil, buyer).want(200)
		el.eq("can_request", "true").eq("change_mind_ok", "true")
		e.buyerUpload(t, buyer, "product_media", pngBytes).want(400)
		ph := e.buyerUpload(t, buyer, "return_photo", pngBytes).want(201)
		check(t, ph.str("url") == "", "foto retur harus privat")
		mine := e.buyerUpload(t, other, "return_photo", pngBytes).want(201)

		ask(no1, map[string]any{"reason": "broken", "description": "Sole came off"}).wantErr(422, "validation", "reason")
		ask(no1, map[string]any{"reason": "damaged", "description": "bad"}).wantErr(422, "validation", "description")
		ask(no1, map[string]any{"reason": "damaged", "description": "Sole came off on day one", "photo_keys": []string{mine.str("key")}}).wantErr(422, "validation", "photo_keys")
		ask(no1, map[string]any{"reason": "damaged", "description": "Sole came off on day one", "items": []map[string]any{{"order_item_id": el.str("items.0.order_item_id"), "qty": 5}}}).wantErr(422, "validation", "items")
		r := ask(no1, map[string]any{"reason": "damaged", "description": "Sole came off on day one", "photo_keys": []string{ph.str("key")}}).want(201)
		r.eq("status", "requested").eq("requested_amount", "6.00").eq("photo_count", "1").eq("items.0.return_qty", "1")
		rid = r.str("id")
		ask(no1, map[string]any{"reason": "damaged", "description": "Sole came off on day one"}).wantErr(409, "return_state")
		check(t, e.scalar(t, `SELECT count(*) FROM notifications WHERE store_id = $1 AND kind = 'return.requested'`, s.storeID) != "0", "seller harus diberi tahu")

		// selama retur terbuka, pesanan tidak masuk saldo walau masa tahan lewat
		age(no1, 2)
		e.sql(t, `UPDATE platform_settings SET value = jsonb_set(value, '{hold_days}', '1') WHERE key = 'payout'`)
		t.Cleanup(func() {
			e.sql(t, `UPDATE platform_settings SET value = jsonb_set(value, '{hold_days}', '3') WHERE key = 'payout'`)
		})
		if _, err := payouts.CompleteDue(ctx); err != nil {
			t.Fatal(err)
		}
		e.call(t, "GET", "/v1/shop/orders/"+no1, nil, buyer).want(200).eq("status", "delivered")

		// seller melihat foto & menolak
		e.call(t, "GET", base+"/returns", nil, s.token).want(200).eq("returns.0.id", rid)
		e.call(t, "GET", base+"/returns/"+rid+"/photos/0", nil, s.token).want(200)
		e.call(t, "GET", "/v1/shop/returns/"+rid+"/photos/0", nil, other).want(404)
		e.call(t, "GET", "/v1/shop/returns/"+rid+"/photos/0", nil, buyer).want(200)
		e.call(t, "POST", base+"/returns/"+rid+"/respond", map[string]any{"accept": false}, s.token).wantErr(422, "validation", "note")
		e.call(t, "POST", base+"/returns/"+rid+"/respond", map[string]any{"accept": false, "note": "It was fine when packed"}, s.token).want(200).
			eq("status", "disputed").eq("escalated_reason", "seller_rejected")
		e.call(t, "POST", base+"/returns/"+rid+"/respond", map[string]any{"accept": true}, s.token).wantErr(409, "return_state")

		// Phsar: CS memutuskan, keuangan mengembalikan uang
		e.call(t, "GET", "/admin/api/returns?status=disputed", nil, mod).want(403)
		l := e.call(t, "GET", "/admin/api/returns?status=disputed", nil, cs).want(200)
		check(t, l.len("returns") >= 1 && l.str("counts.disputed") != "0", "sengketa harus tampil di admin")
		e.call(t, "GET", "/admin/api/returns/"+rid+"/photos/0", nil, cs).want(200)
		e.call(t, "POST", "/admin/api/returns/"+rid+"/decide", map[string]any{"decision": "refund", "amount": "9.00", "note": "Partial for defect"}, cs).wantErr(422, "validation", "amount")
		e.call(t, "POST", "/admin/api/returns/"+rid+"/decide", map[string]any{"decision": "refund", "amount": "3.00", "note": "ok"}, cs).wantErr(422, "validation", "note")
		e.call(t, "POST", "/admin/api/returns/"+rid+"/decide", map[string]any{"decision": "refund", "amount": "3.00", "note": "Photo shows the sole detached; half refund, buyer keeps item"}, fin).want(403)
		e.call(t, "POST", "/admin/api/returns/"+rid+"/decide", map[string]any{"decision": "refund", "amount": "3.00", "note": "Photo shows the sole detached; half refund, buyer keeps item"}, cs).want(200).
			eq("status", "approved").eq("refund_amount", "3.00").eq("needs_pickup", "false")
		e.call(t, "POST", "/admin/api/returns/"+rid+"/refunded", map[string]any{"reference": "PW-RF-1"}, cs).want(403)
		e.call(t, "POST", "/admin/api/returns/"+rid+"/refunded", map[string]any{"reference": ""}, fin).wantErr(422, "validation", "reference")
		e.call(t, "POST", "/admin/api/returns/"+rid+"/refunded", map[string]any{"reference": "PW-RF-1"}, fin).want(200).eq("status", "refunded").eq("refund_ref", "PW-RF-1")
		e.call(t, "POST", "/admin/api/returns/"+rid+"/refunded", map[string]any{"reference": "PW-RF-1"}, fin).wantErr(409, "return_state")
		// $6, komisi 8% → seller $5.52; refund $3 → potong $2.76 → saldo $2.76
		check(t, balance() == "2.76", "saldo harus 5.52 − 2.76 = 2.76, dapat %s", balance())
		e.call(t, "GET", "/v1/shop/orders/"+no1, nil, buyer).want(200).eq("status", "completed")
		b := e.call(t, "GET", "/v1/shop/orders/"+no1+"/return", nil, buyer).want(200)
		b.eq("return.status", "refunded").eq("not_allowed", "already_requested")
		check(t, b.str("return.refund_ref") == "", "nomor referensi tidak untuk pembeli")
	})

	t.Run("berubah pikiran: seller setuju → refund penuh sebelum masuk saldo → pesanan refunded", func(t *testing.T) {
		before := balance()
		no := order(true)
		r := ask(no, map[string]any{"reason": "changed_mind", "description": "Wrong size for me, unused"}).want(201)
		id := r.str("id")
		e.call(t, "POST", base+"/returns/"+id+"/respond", map[string]any{"accept": true, "note": "OK, bring it back"}, s.token).want(200).
			eq("status", "approved").eq("refund_amount", "6.00").eq("needs_pickup", "true")
		e.call(t, "POST", "/admin/api/returns/"+id+"/refunded", map[string]any{"reference": "PW-RF-2"}, fin).want(200)
		e.call(t, "GET", "/v1/shop/orders/"+no, nil, buyer).want(200).eq("status", "refunded")
		check(t, balance() == before, "refund penuh sebelum masuk saldo tidak mengubah saldo: %s → %s", before, balance())
		check(t, e.scalar(t, `SELECT count(*) FROM seller_ledger l JOIN orders o ON o.id = l.order_id WHERE o.order_no = $1`, no) == "0", "tidak ada catatan saldo")
	})

	t.Run("seller tidak menjawab 48 jam → otomatis ke Phsar; Phsar menolak", func(t *testing.T) {
		no := order(true)
		id := ask(no, map[string]any{"reason": "not_as_described", "description": "Colour is different from the video"}).want(201).str("id")
		e.sql(t, `UPDATE returns SET seller_deadline = now() - interval '1 minute' WHERE id = $1`, id)
		e.call(t, "POST", base+"/returns/"+id+"/respond", map[string]any{"accept": true}, s.token).wantErr(409, "return_deadline")
		if n, err := rsvc.EscalateOverdue(ctx); err != nil || n < 1 {
			t.Fatalf("EscalateOverdue: %d %v", n, err)
		}
		e.call(t, "GET", "/admin/api/returns/"+id, nil, cs).want(200).eq("status", "disputed").eq("escalated_reason", "seller_timeout")
		e.call(t, "POST", "/admin/api/returns/"+id+"/decide", map[string]any{"decision": "reject", "note": "Colour matches the product video"}, cs).want(200).eq("status", "rejected")
		e.call(t, "POST", "/v1/shop/returns/"+id+"/cancel", nil, buyer).wantErr(409, "return_state")
	})

	t.Run("pembeli membatalkan lalu boleh mengajukan lagi; masa retur habis ditolak", func(t *testing.T) {
		no := order(true)
		id := ask(no, map[string]any{"reason": "missing_parts", "description": "Shoe laces were missing"}).want(201).str("id")
		e.call(t, "POST", "/v1/shop/returns/"+id+"/cancel", nil, other).want(404)
		e.call(t, "POST", "/v1/shop/returns/"+id+"/cancel", nil, buyer).want(200).eq("status", "cancelled")
		age(no, 8) // return_days 7: berubah pikiran sudah lewat, barang bermasalah juga (maks(7,3))
		e.call(t, "GET", "/v1/shop/orders/"+no+"/return", nil, buyer).want(200).eq("not_allowed", "window_closed")
		ask(no, map[string]any{"reason": "changed_mind", "description": "Do not need it anymore"}).wantErr(422, "validation", "reason")
		age(no, 5)
		ask(no, map[string]any{"reason": "changed_mind", "description": "Do not need it anymore"}).want(201)
		e.call(t, "GET", "/v1/shop/returns", nil, buyer).want(200)
	})

	t.Run("produk tanpa retur: hanya barang bermasalah, minimal 3 hari", func(t *testing.T) {
		e.sql(t, `UPDATE products SET return_days = 0 WHERE store_id = $1`, s.storeID)
		t.Cleanup(func() { e.sql(t, `UPDATE products SET return_days = 7 WHERE store_id = $1`, s.storeID) })
		no := order(true)
		age(no, 2)
		e.call(t, "GET", "/v1/shop/orders/"+no+"/return", nil, buyer).want(200).eq("can_request", "true").eq("change_mind_ok", "false")
		ask(no, map[string]any{"reason": "changed_mind", "description": "Do not like the colour"}).wantErr(422, "validation", "reason")
		ask(no, map[string]any{"reason": "wrong_item", "description": "Received a different model"}).want(201)
	})
}
