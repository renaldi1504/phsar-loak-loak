package qa

import (
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// Tiket bantuan: pembeli & seller → kotak masuk CS → balasan, catatan internal, status
// ---------------------------------------------------------------------------

func TestSupportTickets(t *testing.T) {
	e := getEnv(t)
	mod := newAdmin(t, e, "moderator")
	cs := newAdmin(t, e, "cs")
	s := newSellable(t, e, mod, "Ticket fan", []string{"5.00"}, []int{20})
	buyer, _ := e.login(t)
	other, _ := e.login(t)
	c := e.call(t, "POST", "/v1/shop/checkouts", checkoutBody("k-ticket-1-x", "5.00", item(s.variants[0], 1)), buyer).want(201)
	e.call(t, "POST", "/dev/pay/"+c.str("id"), nil, buyer).want(200)
	no := c.str("orders.0.order_no")

	var bt, st string
	t.Run("pembeli membuka tiket", func(t *testing.T) {
		body := map[string]any{"category": "payment", "subject": "Paid but pending", "message": "I paid with ABA but the order still shows pending.", "order_no": no}
		bad := func(k string, v any) map[string]any {
			m := map[string]any{}
			for a, b := range body {
				m[a] = b
			}
			m[k] = v
			return m
		}
		e.call(t, "POST", "/v1/shop/tickets", bad("category", "spam"), buyer).wantErr(422, "validation", "category")
		e.call(t, "POST", "/v1/shop/tickets", bad("subject", "a"), buyer).wantErr(422, "validation", "subject")
		e.call(t, "POST", "/v1/shop/tickets", bad("message", "short"), buyer).wantErr(422, "validation", "message")
		e.call(t, "POST", "/v1/shop/tickets", body, other).wantErr(422, "validation", "order_no")
		e.call(t, "POST", "/v1/shop/tickets", bad("x", 1), buyer).wantErr(400, "bad_json")
		r := e.call(t, "POST", "/v1/shop/tickets", body, buyer).want(201)
		r.eq("status", "open").eq("order_no", no).eq("messages.0.author_type", "user").eq("waiting", "true")
		check(t, strings.HasPrefix(r.str("ticket_no"), "TK") && r.str("requester_phone") == "", "nomor tiket & data internal: %s", truncate(string(r.Raw)))
		bt = r.str("id")
		e.call(t, "GET", "/v1/shop/tickets/"+bt, nil, other).want(404)
		e.call(t, "GET", "/v1/shop/tickets/not-a-uuid", nil, buyer).want(404)
		e.call(t, "GET", "/v1/shop/tickets", nil, buyer).want(200).eq("tickets.0.id", bt)
		check(t, e.call(t, "GET", "/v1/shop/tickets", nil, other).want(200).len("tickets") == 0, "pembeli lain tidak boleh melihat tiket")
	})

	t.Run("seller membuka tiket atas nama toko", func(t *testing.T) {
		base := "/v1/seller/stores/" + s.storeID + "/tickets"
		e.call(t, "POST", base, map[string]any{"category": "seller", "subject": "Payout question", "message": "When is the next payout sent?", "order_no": "PH000000-XXXXX"}, s.token).
			wantErr(422, "validation", "order_no")
		r := e.call(t, "POST", base, map[string]any{"category": "seller", "subject": "Payout question", "message": "When is the next payout sent?", "order_no": no}, s.token).want(201)
		st = r.str("id")
		r.eq("source", "seller").eq("order_no", no)
		e.call(t, "GET", base+"/"+st, nil, buyer).want(404)
		e.call(t, "GET", base, nil, s.token).want(200).eq("tickets.0.id", st)
		check(t, e.call(t, "GET", "/v1/shop/tickets/"+st, nil, s.token).Status == 404, "tiket seller tidak muncul di jalur pembeli")
	})

	t.Run("CS: kotak masuk, catatan internal, balasan, status", func(t *testing.T) {
		e.call(t, "GET", "/admin/api/tickets", nil, mod).want(403)
		l := e.call(t, "GET", "/admin/api/tickets?filter=open", nil, cs).want(200)
		check(t, l.len("tickets") >= 2 && l.str("counts.waiting") != "0", "tiket harus tampil: %s", truncate(string(l.Raw)))
		e.call(t, "GET", "/admin/api/tickets?filter=bogus", nil, cs).wantErr(422, "validation", "filter")
		q := e.call(t, "GET", "/admin/api/tickets?filter=all&q="+no, nil, cs).want(200)
		check(t, q.len("tickets") == 2, "cari nomor pesanan harus menemukan 2 tiket: %s", truncate(string(q.Raw)))
		d := e.call(t, "GET", "/admin/api/tickets/"+bt, nil, cs).want(200)
		check(t, d.str("requester_phone") != "", "CS harus melihat HP pembeli")

		e.call(t, "POST", "/admin/api/tickets/"+bt+"/reply", map[string]any{"body": ""}, cs).wantErr(422, "validation", "body")
		e.call(t, "POST", "/admin/api/tickets/"+bt+"/reply", map[string]any{"body": "Checked PayWay: APPROVED, callback missing.", "internal": true}, cs).want(200).
			eq("status", "open").eq("messages.1.internal", "true")
		check(t, e.scalar(t, `SELECT count(*) FROM notifications WHERE kind = 'ticket.reply'`) == "0", "catatan internal tidak boleh memberi notifikasi")
		b := e.call(t, "GET", "/v1/shop/tickets/"+bt, nil, buyer).want(200)
		check(t, b.len("messages") == 1 && !strings.Contains(string(b.Raw), "callback missing"), "catatan internal bocor ke pembeli: %s", truncate(string(b.Raw)))

		e.call(t, "POST", "/admin/api/tickets/"+bt+"/reply", map[string]any{"body": "Your payment is confirmed now, sorry!"}, cs).want(200).
			eq("status", "pending").eq("waiting", "false")
		b = e.call(t, "GET", "/v1/shop/tickets/"+bt, nil, buyer).want(200)
		b.eq("status", "pending").eq("messages.1.author_name", "Phsar Support")
		check(t, e.scalar(t, `SELECT link FROM notifications WHERE kind = 'ticket.reply' ORDER BY id DESC LIMIT 1`) != "", "pembeli harus diberi notifikasi")

		// pembeli membalas → kembali open; menandai selesai; menulis lagi → open
		e.call(t, "POST", "/v1/shop/tickets/"+bt+"/messages", map[string]any{"body": "Thanks, it shows paid now."}, buyer).want(200).eq("status", "open")
		e.call(t, "POST", "/v1/shop/tickets/"+bt+"/solve", nil, buyer).want(200).eq("status", "solved")
		e.call(t, "POST", "/v1/shop/tickets/"+bt+"/messages", map[string]any{"body": "One more question."}, other).want(404)
		e.call(t, "POST", "/v1/shop/tickets/"+bt+"/messages", map[string]any{"body": "One more question."}, buyer).want(200).eq("status", "open")

		// tiket seller: tugaskan ke saya, balas & selesai → notifikasi di lonceng toko
		e.call(t, "POST", "/admin/api/tickets/"+st+"/assign", map[string]any{"to_me": true}, cs).want(200)
		m := e.call(t, "GET", "/admin/api/tickets?filter=mine", nil, cs).want(200)
		check(t, m.len("tickets") == 2 && strings.Contains(string(m.Raw), st), "filter 'mine' (dibalas = otomatis ditugaskan): %s", truncate(string(m.Raw)))
		e.call(t, "POST", "/admin/api/tickets/"+st+"/reply", map[string]any{"body": "Payouts run every Monday.", "status": "solved"}, cs).want(200).eq("status", "solved")
		check(t, e.scalar(t, `SELECT count(*) FROM notifications WHERE kind = 'ticket.reply' AND store_id = $1`, s.storeID) == "1", "seller harus dapat notifikasi di lonceng toko")
		e.call(t, "GET", "/v1/seller/stores/"+s.storeID+"/tickets/"+st, nil, s.token).want(200).eq("status", "solved")
		e.call(t, "POST", "/admin/api/tickets/"+st+"/status", map[string]any{"status": "closed"}, cs).wantErr(422, "validation", "status")
		e.call(t, "POST", "/admin/api/tickets/"+st+"/status", map[string]any{"status": "open"}, cs).want(200).eq("status", "open")
		e.call(t, "POST", "/admin/api/tickets/00000000-0000-0000-0000-000000000000/status", map[string]any{"status": "open"}, cs).want(404)

		dash := e.call(t, "GET", "/admin/api/dashboard", nil, cs).want(200)
		check(t, strings.Contains(string(dash.Raw), `"key":"tickets"`), "antrean tiket harus ada di dashboard")
	})

	t.Run("batas tiket per hari", func(t *testing.T) {
		uid := e.scalar(t, `SELECT requester_id::text FROM tickets WHERE id = $1`, bt)
		for i := 0; i < 10; i++ {
			e.sql(t, `INSERT INTO tickets (ticket_no, source, requester_id, category, subject) VALUES ($1, 'buyer', $2, 'other', 'flood')`, "TKFLOOD-"+itoa(i), uid)
		}
		e.call(t, "POST", "/v1/shop/tickets", map[string]any{"category": "other", "subject": "Another one", "message": "This should be rate limited."}, buyer).wantErr(429, "rate_limited")
	})

	t.Run("halaman bantuan website", func(t *testing.T) {
		r := e.page(t, "/en/support")
		check(t, r.Status == 200 && strings.Contains(string(r.Raw), "/assets/support.js"), "halaman /en/support harus ada")
		check(t, strings.Contains(string(e.page(t, "/km/contact").Raw), "/km/support"), "halaman kontak harus menautkan bantuan")
	})
}
