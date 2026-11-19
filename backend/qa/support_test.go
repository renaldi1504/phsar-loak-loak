package qa

import (
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// 17. Alat CS: cari pelanggan & pesanan, catatan, tangguhkan akun, log pembayaran
// ---------------------------------------------------------------------------

func TestSupportTools(t *testing.T) {
	e := getEnv(t)
	mod := newAdmin(t, e, "moderator")
	cs := newAdmin(t, e, "cs")
	fin := newAdmin(t, e, "finance")
	csID := e.scalar(t, `SELECT user_id::text FROM admin_users a JOIN users u ON u.id = a.user_id WHERE a.role_code = 'cs' ORDER BY a.created_at DESC LIMIT 1`)
	s := newSellable(t, e, mod, "Support lamp", []string{"4.00"}, []int{20})
	buyer, phone := e.login(t)
	c := e.call(t, "POST", "/v1/shop/checkouts", checkoutBody("k-support-1-x", "4.00", item(s.variants[0], 1)), buyer).want(201)
	e.call(t, "POST", "/dev/pay/"+c.str("id"), nil, buyer).want(200)
	no, tran := c.str("orders.0.order_no"), e.scalar(t, `SELECT tran_id FROM checkouts WHERE id = $1`, c.str("id"))
	uid := e.scalar(t, `SELECT buyer_id::text FROM checkouts WHERE id = $1`, c.str("id"))

	t.Run("cari pelanggan", func(t *testing.T) {
		e.call(t, "GET", "/admin/api/customers?q="+phone, nil, fin).want(403)
		for _, q := range []string{phone, "+855" + phone[1:], phone[3:], no, uid} {
			r := e.call(t, "GET", "/admin/api/customers?q="+strings.ReplaceAll(q, "+", "%2B"), nil, cs).want(200)
			check(t, r.len("customers") >= 1 && r.str("customers.0.id") == uid, "cari %q harus menemukan pembeli: %s", q, truncate(string(r.Raw)))
		}
		r := e.call(t, "GET", "/admin/api/customers/"+uid, nil, cs).want(200)
		r.eq("stats.orders", "1").eq("orders.0.order_no", no).eq("status", "active")
		e.call(t, "GET", "/admin/api/customers/00000000-0000-0000-0000-000000000000", nil, cs).want(404)
	})

	t.Run("catatan & tangguhkan akun", func(t *testing.T) {
		e.call(t, "POST", "/admin/api/customers/"+uid+"/notes", map[string]string{"note": "x"}, cs).wantErr(422, "validation", "note")
		e.call(t, "POST", "/admin/api/customers/"+uid+"/notes", map[string]string{"note": "Called about late delivery"}, mod).want(403)
		e.call(t, "POST", "/admin/api/customers/"+uid+"/notes", map[string]string{"note": "Called about late delivery"}, cs).want(200).eq("notes.0.kind", "note")
		e.call(t, "POST", "/admin/api/customers/"+uid+"/status", map[string]string{"status": "suspended", "note": "no"}, cs).wantErr(422, "validation", "note")
		e.call(t, "POST", "/admin/api/customers/"+uid+"/status", map[string]string{"status": "suspended", "note": "Fake payment screenshots"}, cs).want(200).
			eq("status", "suspended").eq("sessions", "0").eq("notes.0.kind", "suspended")
		e.call(t, "GET", "/v1/shop/orders", nil, buyer).want(401)
		e.call(t, "POST", "/admin/api/customers/"+uid+"/status", map[string]string{"status": "suspended", "note": "Again please"}, cs).wantErr(409, "no_change")
		e.call(t, "POST", "/admin/api/customers/"+uid+"/status", map[string]string{"status": "active", "note": "Buyer proved payment"}, cs).want(200).eq("status", "active")
		e.call(t, "POST", "/admin/api/customers/"+csID+"/status", map[string]string{"status": "suspended", "note": "Testing self"}, cs).wantErr(409, "self")
		adminID := e.scalar(t, `SELECT user_id::text FROM admin_users WHERE role_code = 'finance' ORDER BY created_at DESC LIMIT 1`)
		e.call(t, "POST", "/admin/api/customers/"+adminID+"/status", map[string]string{"status": "suspended", "note": "Testing admin"}, cs).wantErr(409, "is_admin")
	})

	t.Run("cari & detail pesanan", func(t *testing.T) {
		e.call(t, "GET", "/admin/api/orders?q="+no, nil, mod).want(403)
		for _, q := range []string{no, no[:10], tran, phone} {
			r := e.call(t, "GET", "/admin/api/orders?q="+q, nil, fin).want(200)
			check(t, r.len("orders") >= 1 && r.str("orders.0.order_no") == no, "cari pesanan %q: %s", q, truncate(string(r.Raw)))
		}
		e.call(t, "GET", "/admin/api/orders?status=bogus", nil, cs).wantErr(422, "validation", "status")
		e.call(t, "GET", "/admin/api/orders?status=paid", nil, cs).want(200)
		d := e.call(t, "GET", "/admin/api/orders/"+strings.ToLower(no), nil, cs).want(200)
		d.eq("status", "paid").eq("payment.tran_id", tran).eq("items.0.qty", "1")
		check(t, d.len("payment_events") >= 1 && d.len("events") >= 1, "detail harus memuat log pembayaran & riwayat: %s", truncate(string(d.Raw)))
		e.call(t, "GET", "/admin/api/orders/PH000000-XXXXX", nil, cs).want(404)
	})

	t.Run("log pembayaran", func(t *testing.T) {
		e.call(t, "GET", "/admin/api/payment-events", nil, mod).want(403)
		r := e.call(t, "GET", "/admin/api/payment-events?q="+tran, nil, fin).want(200)
		check(t, r.len("events") >= 1 && r.str("events.0.tran_id") == tran && r.str("events.0.orders.0") == no, "log harus ketemu: %s", truncate(string(r.Raw)))
		e.call(t, "GET", "/admin/api/payment-events?q="+no, nil, fin).want(200).eq("events.0.tran_id", tran)
		e.sql(t, `INSERT INTO payment_events (checkout_id, tran_id, source, verified, status) VALUES ($1, $2, 'callback', false, 'APPROVED')`, c.str("id"), tran)
		p := e.call(t, "GET", "/admin/api/payment-events?filter=problems&q="+tran, nil, fin).want(200)
		check(t, p.len("events") == 1 && p.str("events.0.verified") == "false", "callback tanpa tanda tangan harus masuk 'problems': %s", truncate(string(p.Raw)))
	})
}
