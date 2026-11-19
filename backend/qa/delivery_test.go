package qa

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// 12. Kurir Phsar: ambil tugas, scan label, antar, gagal, admin menugaskan
// ---------------------------------------------------------------------------

func (e *env) courierUpload(t *testing.T, token string, data []byte) resp {
	t.Helper()
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	fw, _ := mw.CreateFormFile("file", "proof.jpg")
	_, _ = fw.Write(data)
	_ = mw.Close()
	req, _ := http.NewRequest(http.MethodPost, e.srv.URL+"/v1/courier/uploads", &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	return e.send(t, req)
}

func TestDelivery(t *testing.T) {
	e := getEnv(t)
	mod := newAdmin(t, e, "moderator")
	cs := newAdmin(t, e, "cs") // izin delivery.dispatch
	s := newSellable(t, e, mod, "Delivery lamp", []string{"6.00"}, []int{30})
	e.sql(t, `INSERT INTO store_pickup_points (store_id, label, address_line, lat, lng, contact_name, contact_phone, is_default, verified_at)
		VALUES ($1, 'Lamp shop', 'Street 3', 10.61, 103.53, 'Vanna', '012999888', true, now())`, s.storeID)
	buyer, _ := e.login(t)
	base := "/v1/seller/stores/" + s.storeID + "/orders"
	n := 0
	ready := func() (string, string) { // pesanan dibayar, diterima, siap dijemput → (nomor pesanan, id tugas)
		t.Helper()
		n++
		c := e.call(t, "POST", "/v1/shop/checkouts", checkoutBody("k-deliv-"+itoa(n)+"-x", "6.00", item(s.variants[0], 1)), buyer).want(201)
		e.call(t, "POST", "/dev/pay/"+c.str("id"), nil, buyer).want(200)
		no := c.str("orders.0.order_no")
		e.call(t, "POST", base+"/"+no+"/accept", nil, s.token).want(200)
		e.call(t, "POST", base+"/"+no+"/ready", nil, s.token).want(200)
		return no, e.scalar(t, `SELECT d.id::text FROM deliveries d JOIN orders o ON o.id = d.order_id WHERE o.order_no = $1`, no)
	}

	c1, c1phone := e.login(t)
	c2, c2phone := e.login(t)
	t.Run("hanya kurir terdaftar; admin menambah kurir", func(t *testing.T) {
		e.call(t, "GET", "/v1/courier/me", nil, c1).wantErr(403, "not_courier")
		e.call(t, "POST", "/admin/api/couriers", map[string]any{"phone": c1phone, "name": "Dara"}, mod).want(403)
		e.call(t, "POST", "/admin/api/couriers", map[string]any{"phone": "abc", "name": "Dara"}, cs).wantErr(422, "validation", "phone")
		e.call(t, "POST", "/admin/api/couriers", map[string]any{"phone": c1phone, "name": "Dara", "vehicle": "motorbike", "plate_no": "PP-1234"}, cs).want(200).eq("name", "Dara")
		e.call(t, "POST", "/admin/api/couriers", map[string]any{"phone": c2phone, "name": "Sok"}, cs).want(200)
		// nomor yang belum pernah login → akun dibuat, kurir tinggal login
		e.call(t, "POST", "/admin/api/couriers", map[string]any{"phone": "0971234567", "name": "Baru"}, cs).want(200)
		e.call(t, "GET", "/v1/courier/me", nil, c1).want(200).eq("name", "Dara").eq("online", "false")
		e.call(t, "POST", "/v1/courier/me/online", map[string]bool{"online": true}, c1).want(200).eq("online", "true")
		check(t, e.page(t, "/courier/").Status == 200, "halaman kurir tidak terbuka")
	})

	no1, task1 := ready()
	t.Run("tugas tersedia tanpa data pembeli; kurir online dapat notifikasi", func(t *testing.T) {
		check(t, e.scalar(t, `SELECT count(*) FROM notifications n JOIN users u ON u.id = n.user_id WHERE n.kind = 'courier.new_task' AND u.phone = '+855' || substr($1, 2)`, c1phone) == "1",
			"kurir online tidak dapat notifikasi tugas baru")
		l := e.call(t, "GET", "/v1/courier/tasks?view=available", nil, c2).want(200)
		check(t, strings.Contains(string(l.Raw), no1), "tugas tidak muncul: %s", l.Raw)
		check(t, !strings.Contains(string(l.Raw), "012 345 678") && !strings.Contains(string(l.Raw), "Street 7"), "data pembeli terlihat sebelum tugas diambil")
		e.call(t, "GET", "/v1/courier/tasks?view=available", nil, buyer).wantErr(403, "not_courier")
	})
	t.Run("ambil tugas: yang pertama menang", func(t *testing.T) {
		r := e.call(t, "POST", "/v1/courier/tasks/"+task1+"/claim", nil, c1).want(200)
		r.eq("status", "assigned").eq("pickup_phone", "012999888").eq("recipient", "Sokha")
		check(t, r.len("items") == 1, "isi paket tidak tampil")
		e.call(t, "POST", "/v1/courier/tasks/"+task1+"/claim", nil, c2).wantErr(409, "taken")
		e.call(t, "GET", "/v1/courier/tasks/"+task1, nil, c2).want(404)
	})
	t.Run("scan label: salah paket ditolak, benar → dijemput", func(t *testing.T) {
		e.call(t, "POST", "/v1/courier/tasks/"+task1+"/deliver", nil, c1).wantErr(409, "bad_status")
		e.call(t, "POST", "/v1/courier/tasks/"+task1+"/pickup", map[string]string{"code": "PH000000-AAAAA"}, c1).wantErr(422, "wrong_package")
		e.call(t, "POST", "/v1/courier/tasks/"+task1+"/pickup", map[string]string{"code": " " + strings.ToLower(no1) + " "}, c1).want(200).eq("status", "picked_up")
		e.call(t, "GET", "/v1/shop/orders/"+no1, nil, buyer).want(200).eq("status", "picked_up")
		e.call(t, "GET", base+"?filter=shipping", nil, s.token).want(200).eq("orders.0.order_no", no1)
		check(t, e.scalar(t, `SELECT count(*) FROM notifications WHERE kind = 'buyer.order_on_the_way' AND link LIKE '%'||$1`, no1) == "1", "pembeli tidak diberi tahu")
	})
	t.Run("terantar dengan foto bukti", func(t *testing.T) {
		bad := e.courierUpload(t, c2, jpegBytes).want(201).str("key")
		e.call(t, "POST", "/v1/courier/tasks/"+task1+"/deliver", map[string]string{"proof_key": bad}, c1).wantErr(422, "validation", "proof_key")
		key := e.courierUpload(t, c1, jpegBytes).want(201).str("key")
		check(t, strings.HasPrefix(key, "private/delivery/"), "foto bukti harus privat: %s", key)
		e.call(t, "POST", "/v1/courier/tasks/"+task1+"/deliver", map[string]string{"proof_key": key}, c1).want(200).eq("status", "delivered").eq("has_proof", "true")
		r := e.call(t, "GET", "/v1/shop/orders/"+no1, nil, buyer).want(200)
		r.eq("status", "delivered")
		check(t, strings.Contains(string(r.Raw), `"status":"delivered"`), "riwayat tidak mencatat terantar")
		check(t, e.scalar(t, `SELECT count(*) FROM notifications WHERE kind = 'buyer.order_delivered' AND link LIKE '%'||$1`, no1) == "1", "pembeli tidak diberi tahu terantar")
		e.call(t, "GET", base+"?filter=done", nil, s.token).want(200).eq("orders.0.order_no", no1)
		e.call(t, "GET", "/v1/courier/tasks?view=done", nil, c1).want(200).eq("tasks.0.order_no", no1)
		e.call(t, "GET", "/v1/courier/me", nil, c1).want(200).eq("done_today", "1")
		// foto bukti hanya untuk admin berizin
		req, _ := http.NewRequest(http.MethodGet, e.srv.URL+"/admin/api/deliveries/"+task1+"/proof", nil)
		req.Header.Set("Authorization", "Bearer "+cs)
		pr := e.send(t, req)
		check(t, pr.Status == 200 && len(pr.Raw) == len(jpegBytes), "foto bukti tidak bisa dibuka admin: %d", pr.Status)
		req2, _ := http.NewRequest(http.MethodGet, e.srv.URL+"/admin/api/deliveries/"+task1+"/proof", nil)
		req2.Header.Set("Authorization", "Bearer "+mod)
		check(t, e.send(t, req2).Status == 403, "moderator bisa membuka foto bukti")
	})

	t.Run("gagal antar → coba lagi; admin memindahkan ke kurir lain", func(t *testing.T) {
		_, task2 := ready()
		e.call(t, "POST", "/v1/courier/tasks/"+task2+"/claim", nil, c1).want(200)
		no2 := e.scalar(t, `SELECT o.order_no FROM deliveries d JOIN orders o ON o.id = d.order_id WHERE d.id = $1`, task2)
		e.call(t, "POST", "/v1/courier/tasks/"+task2+"/pickup", map[string]string{"code": no2}, c1).want(200)
		e.call(t, "POST", "/v1/courier/tasks/"+task2+"/fail", map[string]string{"reason": "other"}, c1).wantErr(422, "validation", "note")
		e.call(t, "POST", "/v1/courier/tasks/"+task2+"/fail", map[string]string{"reason": "buyer_unreachable", "note": "Phone off"}, c1).want(200).eq("status", "failed")
		l := e.call(t, "GET", "/admin/api/deliveries?status=failed", nil, cs).want(200)
		check(t, strings.Contains(string(l.Raw), no2) && strings.Contains(string(l.Raw), "Phone off"), "admin tidak melihat yang gagal: %s", l.Raw)
		e.call(t, "POST", "/v1/courier/tasks/"+task2+"/retry", nil, c1).want(200).eq("status", "picked_up").eq("attempts", "2")
		c2id := e.scalar(t, `SELECT id::text FROM users WHERE phone = '+855' || substr($1, 2)`, c2phone)
		e.call(t, "POST", "/admin/api/deliveries/"+task2+"/assign", map[string]string{"courier_id": c2id}, cs).want(200).eq("courier_name", "Sok")
		e.call(t, "GET", "/v1/courier/tasks/"+task2, nil, c1).want(404)
		e.call(t, "GET", "/v1/courier/tasks?view=mine", nil, c2).want(200).eq("tasks.0.order_no", no2)
		e.call(t, "POST", "/admin/api/deliveries/"+task2+"/assign", map[string]string{"courier_id": ""}, cs).wantErr(409, "bad_status") // sudah dijemput: tidak bisa dilepas
		e.call(t, "POST", "/v1/courier/tasks/"+task2+"/deliver", nil, c2).want(200).eq("status", "delivered")
	})
	t.Run("lepas tugas sebelum dijemput; kurir nonaktif ditolak", func(t *testing.T) {
		_, task3 := ready()
		e.call(t, "POST", "/v1/courier/tasks/"+task3+"/claim", nil, c2).want(200)
		e.call(t, "POST", "/v1/courier/tasks/"+task3+"/release", nil, c2).want(200)
		e.call(t, "GET", "/v1/courier/tasks/"+task3, nil, c1).want(200).eq("status", "unassigned")
		e.call(t, "POST", "/admin/api/couriers", map[string]any{"phone": c2phone, "name": "Sok", "active": false}, cs).want(200).eq("active", "false")
		e.call(t, "GET", "/v1/courier/me", nil, c2).wantErr(403, "not_courier")
		e.call(t, "POST", "/v1/courier/tasks/"+task3+"/claim", nil, c2).wantErr(403, "not_courier")
		cl := e.call(t, "GET", "/admin/api/couriers", nil, cs).want(200)
		check(t, cl.len("couriers") >= 3, "daftar kurir: %s", cl.Raw)
	})
}
