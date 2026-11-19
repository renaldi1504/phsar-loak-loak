package qa

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"phsar/internal/jobs"
)

// ---------------------------------------------------------------------------
// courierreturn: gagal antar → coba lagi (maks.) / kembalikan ke toko → seller konfirmasi → batal + refund + stok
// ---------------------------------------------------------------------------

func TestCourierReturn(t *testing.T) {
	e := getEnv(t)
	mod := newAdmin(t, e, "moderator")
	cs := newAdmin(t, e, "cs")           // delivery.dispatch
	finance := newAdmin(t, e, "finance") // refund.decide
	sa := newAdmin(t, e, "super_admin")
	s := newSellable(t, e, mod, "Return kettle", []string{"9.00"}, []int{10})
	e.sql(t, `INSERT INTO store_pickup_points (store_id, label, address_line, lat, lng, contact_name, contact_phone, is_default, verified_at)
		VALUES ($1, 'Kettle shop', 'Street 9', 10.62, 103.52, 'Rina', '012777666', true, now())`, s.storeID)
	other := newSellable(t, e, mod, "Other shop item", []string{"5.00"}, []int{5})
	buyer, _ := e.login(t)
	c1, c1phone := e.login(t)
	c2, c2phone := e.login(t)
	e.call(t, "POST", "/admin/api/couriers", map[string]any{"phone": c1phone, "name": "Rotha"}, cs).want(200)
	e.call(t, "POST", "/admin/api/couriers", map[string]any{"phone": c2phone, "name": "Kosal"}, cs).want(200)
	c2id := e.scalar(t, `SELECT id::text FROM users WHERE phone = '+855' || substr($1, 2)`, c2phone)
	stock := func() string {
		return e.scalar(t, `SELECT stock_qty::text FROM product_variants WHERE id = $1`, s.variants[0])
	}
	base := "/v1/seller/stores/" + s.storeID
	n := 0
	// pesanan dibayar → siap → diambil c1 → dijemput → gagal; mengembalikan (nomor pesanan, id tugas)
	failed := func(reason string) (string, string) {
		t.Helper()
		n++
		c := e.call(t, "POST", "/v1/shop/checkouts", checkoutBody("k-cr-"+itoa(n)+"-x", "18.00", item(s.variants[0], 2)), buyer).want(201)
		e.call(t, "POST", "/dev/pay/"+c.str("id"), nil, buyer).want(200)
		no := c.str("orders.0.order_no")
		e.call(t, "POST", base+"/orders/"+no+"/accept", nil, s.token).want(200)
		e.call(t, "POST", base+"/orders/"+no+"/ready", nil, s.token).want(200)
		task := e.scalar(t, `SELECT d.id::text FROM deliveries d JOIN orders o ON o.id = d.order_id WHERE o.order_no = $1`, no)
		e.call(t, "POST", "/v1/courier/tasks/"+task+"/claim", nil, c1).want(200)
		e.call(t, "POST", "/v1/courier/tasks/"+task+"/pickup", map[string]string{"code": no}, c1).want(200)
		e.call(t, "POST", "/v1/courier/tasks/"+task+"/fail", map[string]string{"reason": reason, "note": "test"}, c1).want(200).eq("status", "failed")
		return no, task
	}
	retID := func(task string) string {
		return e.scalar(t, `SELECT id::text FROM courier_returns WHERE delivery_id = $1`, task)
	}

	start := stock()
	no1, task1 := failed("buyer_unreachable")
	t.Run("gagal antar dengan foto bukti; batas percobaan", func(t *testing.T) {
		check(t, stock() == itoa(atoi(start)-2), "stok harus berkurang saat dibayar: %s", stock())
		// foto bukti gagal antar (milik kurir lain ditolak)
		key := e.courierUpload(t, c1, jpegBytes).want(201).str("key")
		info := e.call(t, "GET", "/v1/courier/returns/"+task1, nil, c1).want(200)
		info.eq("status", "failed").eq("attempts", "1").eq("max_attempts", "3").eq("can_retry", "true").eq("can_return", "true")
		e.call(t, "GET", "/v1/courier/returns/"+task1, nil, c2).want(404)
		e.call(t, "GET", "/v1/courier/returns/"+task1, nil, buyer).wantErr(403, "not_courier")
		// coba lagi sampai batas, gagal lagi dengan foto
		e.call(t, "POST", "/v1/courier/tasks/"+task1+"/retry", nil, c1).want(200).eq("attempts", "2")
		bad := e.courierUpload(t, c2, jpegBytes).want(201).str("key")
		e.call(t, "POST", "/v1/courier/tasks/"+task1+"/fail", map[string]string{"reason": "buyer_refused", "proof_key": bad}, c1).wantErr(422, "validation", "proof_key")
		e.call(t, "POST", "/v1/courier/tasks/"+task1+"/fail", map[string]string{"reason": "buyer_refused", "proof_key": key}, c1).want(200)
		check(t, e.scalar(t, `SELECT fail_proof_key FROM deliveries WHERE id = $1`, task1) == key, "foto bukti gagal antar tidak tersimpan")
		e.call(t, "POST", "/v1/courier/tasks/"+task1+"/retry", nil, c1).want(200).eq("attempts", "3")
		e.call(t, "POST", "/v1/courier/tasks/"+task1+"/fail", map[string]string{"reason": "buyer_unreachable"}, c1).want(200)
		e.call(t, "POST", "/v1/courier/tasks/"+task1+"/retry", nil, c1).wantErr(409, "max_attempts")
		e.call(t, "GET", "/v1/courier/returns/"+task1, nil, c1).want(200).eq("can_retry", "false").eq("can_return", "true")
	})

	t.Run("kurir: kembalikan ke toko → sudah di toko", func(t *testing.T) {
		e.call(t, "POST", "/v1/courier/returns/"+task1+"/returned", nil, c1).wantErr(409, "bad_status")
		e.call(t, "POST", "/v1/courier/returns/"+task1+"/start", nil, c2).want(404)
		r := e.call(t, "POST", "/v1/courier/returns/"+task1+"/start", nil, c1).want(200)
		r.eq("status", "returning").eq("return.status", "returning").eq("return.reason", "buyer_unreachable").eq("return.return_fee", "0.00").eq("return.fee_payer", "phsar")
		e.call(t, "POST", "/v1/courier/returns/"+task1+"/start", nil, c1).want(200) // klik dua kali
		e.call(t, "POST", "/v1/courier/tasks/"+task1+"/deliver", nil, c1).wantErr(409, "bad_status")
		mine := e.call(t, "GET", "/v1/courier/tasks?view=mine", nil, c1).want(200)
		check(t, strings.Contains(string(mine.Raw), `"status":"returning"`), "tugas pengembalian harus ada di daftar kurir: %s", mine.Raw)
		check(t, e.scalar(t, `SELECT count(*) FROM notifications WHERE kind = 'order.courier_return_started' AND store_id = $1`, s.storeID) == "1", "seller tidak diberi tahu")
		check(t, e.scalar(t, `SELECT count(*) FROM notifications WHERE kind = 'buyer.delivery_returning' AND link LIKE '%'||$1`, no1) == "1", "pembeli tidak diberi tahu")
		// seller melihat paket yang sedang kembali; toko lain tidak
		l := e.call(t, "GET", base+"/courier-returns", nil, s.token).want(200)
		l.eq("returns.0.order_no", no1).eq("returns.0.status", "returning").eq("returns.0.courier_name", "Rotha").eq("returns.0.items.0.qty", "2")
		check(t, !strings.Contains(string(l.Raw), "Sokha") && !strings.Contains(string(l.Raw), "012 345 678"), "data pembeli tidak boleh terlihat seller: %s", l.Raw)
		e.call(t, "GET", base+"/courier-returns", nil, other.token).wantErr(404, "store_not_found")
		// sampai di toko (foto opsional)
		key := e.courierUpload(t, c1, jpegBytes).want(201).str("key")
		e.call(t, "POST", "/v1/courier/returns/"+task1+"/returned", map[string]string{"proof_key": key}, c1).want(200).eq("status", "returned").eq("return.status", "awaiting_seller")
		check(t, e.scalar(t, `SELECT (confirm_by > now() + interval '23 hours')::text FROM courier_returns WHERE delivery_id = $1`, task1) == "true", "batas konfirmasi seller harus ~24 jam")
		check(t, e.scalar(t, `SELECT count(*) FROM notifications WHERE kind = 'order.courier_return_arrived' AND store_id = $1`, s.storeID) == "1", "seller tidak diminta konfirmasi")
		e.call(t, "GET", base+"/today", nil, s.token).want(200).eq("courier_returns", "1")
		check(t, e.scalar(t, `SELECT count(*) FROM order_events oe JOIN orders o ON o.id = oe.order_id WHERE o.order_no = $1 AND oe.status IN ('return_to_store','returned_to_store')`, no1) == "2",
			"riwayat pesanan tidak mencatat pengembalian")
	})

	t.Run("seller konfirmasi: validasi, izin, batal + refund + stok kembali", func(t *testing.T) {
		id := retID(task1)
		url := base + "/courier-returns/" + id + "/confirm"
		e.call(t, "POST", url, map[string]string{"condition": "lost"}, s.token).wantErr(422, "validation", "condition")
		e.call(t, "POST", url, map[string]string{"condition": "damaged"}, s.token).wantErr(422, "validation", "note")
		e.call(t, "POST", url, map[string]any{"condition": "ok", "extra": 1}, s.token).want(400)
		e.call(t, "POST", "/v1/seller/stores/"+other.storeID+"/courier-returns/"+id+"/confirm", map[string]string{"condition": "ok"}, other.token).want(404)
		before := stock()
		r := e.call(t, "POST", url, map[string]string{"condition": "ok", "note": "Box fine"}, s.token).want(200)
		r.eq("status", "received").eq("resolved_by", "seller").eq("condition", "ok").eq("restocked", "true").eq("order_status", "cancelled")
		e.call(t, "POST", url, map[string]string{"condition": "ok"}, s.token).wantErr(409, "bad_status")
		check(t, stock() == itoa(atoi(before)+2), "stok harus kembali: %s → %s", before, stock())
		o := e.call(t, "GET", "/v1/shop/orders/"+no1, nil, buyer).want(200)
		o.eq("status", "cancelled")
		check(t, strings.Contains(string(o.Raw), `"return_received"`) && strings.Contains(string(o.Raw), `"refund_pending"`), "riwayat pembeli: %s", o.Raw)
		check(t, e.scalar(t, `SELECT refund_status || ' ' || cancel_reason || ' ' || refund_amount::text FROM orders WHERE order_no = $1`, no1) == "pending delivery_failed 18.00",
			"refund pembeli harus menunggu admin")
		check(t, e.scalar(t, `SELECT count(*) FROM notifications WHERE kind = 'buyer.order_cancelled' AND link LIKE '%'||$1`, no1) == "1", "pembeli tidak diberi tahu batal")
		check(t, e.scalar(t, `SELECT status FROM deliveries WHERE id = $1`, task1) == "returned", "status pengantaran")
		check(t, e.scalar(t, `SELECT count(*) FROM audit_logs WHERE entity_table = 'courier_returns' AND entity_id = $1`, id) != "0", "perubahan harus tercatat di audit")
		e.call(t, "GET", base+"/today", nil, s.token).want(200).eq("courier_returns", "0")
		// refund manual yang ada
		rf := e.call(t, "GET", "/admin/api/refunds", nil, finance).want(200)
		check(t, strings.Contains(string(rf.Raw), no1) && strings.Contains(string(rf.Raw), "delivery_failed"), "harus masuk antrean refund: %s", rf.Raw)
		oid := e.scalar(t, `SELECT id::text FROM orders WHERE order_no = $1`, no1)
		e.call(t, "POST", "/admin/api/refunds/"+oid+"/done", map[string]string{"refund_ref": "RF-CR-1"}, finance).want(200)
		h := e.call(t, "GET", base+"/courier-returns?view=all", nil, s.token).want(200)
		h.eq("returns.0.status", "received")
	})

	t.Run("admin: antrean, izin, kembalikan + tugaskan kurir, selesaikan manual tanpa restock", func(t *testing.T) {
		no2, task2 := failed("damaged")
		e.call(t, "GET", "/admin/api/courier-returns", nil, mod).want(403)
		e.call(t, "POST", "/admin/api/courier-returns", map[string]string{"delivery_id": task2}, mod).want(403)
		e.call(t, "POST", "/admin/api/courier-returns", map[string]string{"delivery_id": task2, "courier_id": "nope"}, cs).wantErr(422, "validation", "courier_id")
		r := e.call(t, "POST", "/admin/api/courier-returns", map[string]string{"delivery_id": task2, "courier_id": c2id}, cs).want(200)
		r.eq("status", "returning").eq("courier_name", "Kosal").eq("started_by", "admin")
		e.call(t, "POST", "/admin/api/courier-returns", map[string]string{"delivery_id": task2}, cs).wantErr(409, "bad_status")
		check(t, e.scalar(t, `SELECT count(*) FROM notifications n WHERE n.kind = 'courier.return_task' AND n.user_id = $1`, c2id) == "1", "kurir pengembalian tidak diberi tahu")
		e.call(t, "GET", "/v1/courier/returns/"+task2, nil, c1).want(404) // dipindah ke Kosal
		e.call(t, "GET", "/v1/courier/returns/"+task2, nil, c2).want(200).eq("status", "returning")
		id := retID(task2)
		e.call(t, "POST", "/admin/api/courier-returns/"+id+"/assign", map[string]string{"courier_id": c2id}, mod).want(403)
		c1id := e.scalar(t, `SELECT id::text FROM users WHERE phone = '+855' || substr($1, 2)`, c1phone)
		e.call(t, "POST", "/admin/api/courier-returns/"+id+"/assign", map[string]string{"courier_id": c1id}, cs).want(200).eq("courier_name", "Rotha")
		l := e.call(t, "GET", "/admin/api/courier-returns", nil, cs).want(200)
		l.eq("returns.0.order_no", no2).eq("counts.returning", "1")
		check(t, e.call(t, "GET", "/admin/api/dashboard", nil, cs).want(200).str("queues") != "", "dasbor")
		// foto bukti gagal antar tidak ada → 404; foto tidak bisa dibuka moderator
		req, _ := http.NewRequest(http.MethodGet, e.srv.URL+"/admin/api/courier-returns/"+id+"/proof?which=fail", nil)
		req.Header.Set("Authorization", "Bearer "+cs)
		check(t, e.send(t, req).Status == 404, "tanpa foto harus 404")
		e.call(t, "POST", "/admin/api/courier-returns/"+id+"/resolve", map[string]any{"note": "x"}, cs).wantErr(422, "validation", "restock")
		before := stock()
		e.call(t, "POST", "/admin/api/courier-returns/"+id+"/resolve", map[string]any{"restock": false, "note": "Store says broken"}, cs).want(200).
			eq("status", "received").eq("resolved_by", "admin").eq("restocked", "false").eq("condition", "damaged")
		check(t, stock() == before, "barang rusak tidak boleh kembali ke stok")
		check(t, e.scalar(t, `SELECT count(*) FROM notifications WHERE kind = 'order.courier_return_closed' AND store_id = $1`, s.storeID) == "1", "seller diberi tahu penyelesaian admin")
		e.call(t, "POST", "/admin/api/courier-returns/"+id+"/resolve", map[string]any{"restock": true}, cs).wantErr(409, "bad_status")
		d := e.call(t, "GET", "/admin/api/courier-returns?status=done", nil, cs).want(200)
		check(t, strings.Contains(string(d.Raw), no2), "riwayat selesai: %s", d.Raw)
	})

	t.Run("job: seller tidak konfirmasi → otomatis diterima", func(t *testing.T) {
		no3, task3 := failed("wrong_address")
		e.call(t, "POST", "/v1/courier/returns/"+task3+"/start", nil, c1).want(200)
		e.call(t, "POST", "/v1/courier/returns/"+task3+"/returned", nil, c1).want(200)
		e.sql(t, `UPDATE courier_returns SET confirm_by = now() - interval '1 minute' WHERE delivery_id = $1`, task3)
		before := stock()
		fn, ok := jobs.Runner("courierreturn.auto_receive")
		check(t, ok, "job harus terdaftar untuk Run now")
		k, err := fn(context.Background())
		check(t, err == nil && k == 1, "job: %d %v", k, err)
		check(t, stock() == itoa(atoi(before)+2), "stok harus kembali (alasan bukan rusak)")
		check(t, e.scalar(t, `SELECT resolved_by_type::text || ' ' || status FROM courier_returns WHERE delivery_id = $1`, task3) == "system received", "harus diterima sistem")
		check(t, e.scalar(t, `SELECT status::text FROM orders WHERE order_no = $1`, no3) == "cancelled", "pesanan harus batal")
		k, _ = fn(context.Background())
		check(t, k == 0, "tidak boleh diproses dua kali")
		jl := e.call(t, "GET", "/admin/api/jobs", nil, sa).want(200)
		check(t, jl.str("jobs."+itoa(len(jobs.Defs)-1)+".name") == "courierreturn.auto_receive", "job harus di akhir daftar")
	})
}

func atoi(s string) int {
	n := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		}
	}
	return n
}
