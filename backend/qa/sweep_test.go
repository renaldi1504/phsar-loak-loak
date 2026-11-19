package qa

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"phsar/internal/httpx"
)

// ---------------------------------------------------------------------------
// 18. Commission sweep: komisi pesanan selesai boleh dipindahkan ke rekening perusahaan;
//     uang seller, pesanan berjalan dan refund tetap ditahan.
// ---------------------------------------------------------------------------

func TestCommissionSweep(t *testing.T) {
	e := getEnv(t)
	ctx := context.Background()
	mod := newAdmin(t, e, "moderator")
	cs := newAdmin(t, e, "cs")
	fin := newAdmin(t, e, "finance")
	s := newSellable(t, e, mod, "Sweep mug", []string{"6.00"}, []int{50})
	e.sql(t, `INSERT INTO store_pickup_points (store_id, label, address_line, lat, lng, contact_name, contact_phone, is_default, verified_at)
		VALUES ($1, 'Mug shop', 'Street 3', 10.61, 103.53, 'Sok', '012777444', true, now())`, s.storeID)
	courier, cphone := e.login(t)
	e.call(t, "POST", "/admin/api/couriers", map[string]any{"phone": cphone, "name": "Rotha"}, cs).want(200)
	buyer, _ := e.login(t)
	payouts := httpx.NewPayouts(e.cfg, slog.New(slog.NewTextHandler(io.Discard, nil)), e.pool)
	base := "/v1/seller/stores/" + s.storeID
	n := 0
	pay := func() string {
		t.Helper()
		n++
		c := e.call(t, "POST", "/v1/shop/checkouts", checkoutBody("k-sweep-"+itoa(n)+"-x", "6.00", item(s.variants[0], 1)), buyer).want(201)
		e.call(t, "POST", "/dev/pay/"+c.str("id"), nil, buyer).want(200)
		return c.str("orders.0.order_no")
	}
	deliver := func(no string) {
		t.Helper()
		e.call(t, "POST", base+"/orders/"+no+"/accept", nil, s.token).want(200)
		e.call(t, "POST", base+"/orders/"+no+"/ready", nil, s.token).want(200)
		task := e.scalar(t, `SELECT d.id::text FROM deliveries d JOIN orders o ON o.id = d.order_id WHERE o.order_no = $1`, no)
		e.call(t, "POST", "/v1/courier/tasks/"+task+"/claim", nil, courier).want(200)
		e.call(t, "POST", "/v1/courier/tasks/"+task+"/pickup", map[string]string{"code": no}, courier).want(200)
		e.call(t, "POST", "/v1/courier/tasks/"+task+"/deliver", nil, courier).want(200)
	}
	sum := func() resp { return e.call(t, "GET", "/admin/api/sweep", nil, fin).want(200) }
	dec := func(r resp, k string) decimal.Decimal {
		d, err := decimal.NewFromString(r.str("summary." + k))
		if err != nil {
			t.Fatalf("summary.%s = %q", k, r.str("summary."+k))
		}
		return d
	}
	delta := func(a, b resp, k, want string) {
		t.Helper()
		got := dec(b, k).Sub(dec(a, k))
		if !got.Equal(decimal.RequireFromString(want)) {
			t.Errorf("perubahan %s = %s, harusnya %s", k, got.StringFixed(2), want)
		}
	}

	t.Run("izin", func(t *testing.T) {
		e.call(t, "GET", "/admin/api/sweep", nil, mod).want(403)
		e.call(t, "GET", "/admin/api/sweep", nil, cs).want(403)
		e.call(t, "POST", "/admin/api/sweep", map[string]any{"amount": "1"}, cs).want(403)
	})

	t.Run("pesanan dibayar ditahan; setelah selesai komisinya boleh di-sweep", func(t *testing.T) {
		s0 := sum()
		no := pay()
		s1 := sum()
		delta(s0, s1, "collected", "6.00")
		delta(s0, s1, "held_orders", "6.00")
		delta(s0, s1, "commission_net", "0.00")
		delta(s0, s1, "expected_balance", "6.00")

		deliver(no)
		e.sql(t, `UPDATE deliveries SET delivered_at = now() - interval '10 days' WHERE order_id = (SELECT id FROM orders WHERE order_no = $1)`, no)
		if _, err := payouts.CompleteDue(ctx); err != nil {
			t.Fatal(err)
		}
		s2 := sum()
		delta(s1, s2, "held_orders", "-6.00")
		delta(s1, s2, "held_sellers", "5.52") // $6 − komisi 8%
		delta(s1, s2, "commission_gross", "0.48")
		delta(s1, s2, "expected_balance", "0.00")
		check(t, s2.len("monthly") == 12, "ringkasan 12 bulan: %d", s2.len("monthly"))
	})

	t.Run("pesanan ditolak seller: refund ditahan sampai dikirim", func(t *testing.T) {
		s0 := sum()
		no := pay()
		e.call(t, "POST", base+"/orders/"+no+"/reject", map[string]string{"reason": "out_of_stock"}, s.token).want(200)
		s1 := sum()
		delta(s0, s1, "held_orders", "6.00")
		oid := e.scalar(t, `SELECT id::text FROM orders WHERE order_no = $1`, no)
		e.call(t, "POST", "/admin/api/refunds/"+oid+"/done", map[string]string{"refund_ref": "RF-SWEEP-1"}, fin).want(200)
		s2 := sum()
		delta(s1, s2, "held_orders", "-6.00")
		delta(s1, s2, "refunds", "6.00")
		delta(s1, s2, "expected_balance", "-6.00")
		delta(s0, s2, "available", "0.00")
	})

	var sid string
	t.Run("catat sweep: validasi, tidak boleh melebihi yang tersedia", func(t *testing.T) {
		// database tes dipakai bersama tes lain (retur/rekonsiliasi bisa membuat komisi bersih negatif): tambah pesanan selesai seperlunya
		for i := 0; i < 30 && !dec(sum(), "available").IsPositive(); i++ {
			no := pay()
			deliver(no)
			e.sql(t, `UPDATE deliveries SET delivered_at = now() - interval '10 days' WHERE order_id = (SELECT id FROM orders WHERE order_no = $1)`, no)
			if _, err := payouts.CompleteDue(ctx); err != nil {
				t.Fatal(err)
			}
		}
		s0 := sum()
		avail := dec(s0, "available")
		check(t, avail.IsPositive(), "harus ada komisi yang bisa di-sweep: %s", avail)
		today := time.Now().In(time.FixedZone("ICT", 7*3600)).Format("2006-01-02")
		ok := map[string]any{"amount": avail.StringFixed(2), "to_account": "ABA Phsar Co. 1234", "reference": "TRF-0001", "swept_on": today, "note": "November"}
		bad := func(field string, k string, v any) {
			b := map[string]any{}
			for kk, vv := range ok {
				b[kk] = vv
			}
			b[k] = v
			e.call(t, "POST", "/admin/api/sweep", b, fin).wantErr(422, "validation", field)
		}
		bad("amount", "amount", "abc")
		bad("amount", "amount", "0")
		bad("amount", "amount", avail.Add(decimal.RequireFromString("0.01")).StringFixed(2))
		bad("to_account", "to_account", "")
		bad("reference", "reference", "x")
		bad("swept_on", "swept_on", "2099-01-01")
		bad("swept_on", "swept_on", "16/11/2026")
		r := e.call(t, "POST", "/admin/api/sweep", ok, fin).want(201)
		r.eq("summary.available", "0.00").eq("sweeps.0.reference", "TRF-0001").eq("sweeps.0.amount", avail.StringFixed(2))
		sid = r.str("sweeps.0.id")
		delta(s0, r, "swept", avail.StringFixed(2))
		delta(s0, r, "expected_balance", avail.Neg().StringFixed(2))
		e.call(t, "POST", "/admin/api/sweep", map[string]any{"amount": "0.01", "to_account": "ABA Phsar Co. 1234", "reference": "TRF-0002", "swept_on": today}, fin).
			wantErr(422, "validation", "amount")
		check(t, e.scalar(t, `SELECT count(*) FROM audit_logs WHERE entity_table = 'commission_sweeps' AND entity_id = $1`, sid) != "0", "sweep harus tercatat di audit")
	})

	t.Run("batalkan sweep: uang dihitung kembali; tidak bisa dua kali", func(t *testing.T) {
		s0 := sum()
		e.call(t, "POST", "/admin/api/sweep/"+sid+"/void", map[string]string{"note": "x"}, fin).wantErr(422, "validation", "note")
		r := e.call(t, "POST", "/admin/api/sweep/"+sid+"/void", map[string]string{"note": "Typed the wrong amount"}, fin).want(200)
		r.eq("sweeps.0.void_note", "Typed the wrong amount")
		check(t, r.str("sweeps.0.voided_at") != "", "harus tercatat dibatalkan")
		delta(s0, r, "swept", "-"+s0.str("sweeps.0.amount"))
		check(t, dec(r, "available").IsPositive(), "setelah dibatalkan komisi tersedia lagi")
		e.call(t, "POST", "/admin/api/sweep/"+sid+"/void", map[string]string{"note": "Typed the wrong amount"}, fin).wantErr(409, "already_voided")
		e.call(t, "POST", "/admin/api/sweep/00000000-0000-0000-0000-000000000000/void", map[string]string{"note": "Typed the wrong amount"}, fin).want(404)
	})

	t.Run("ekspor CSV", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, e.srv.URL+"/admin/api/sweep/export.csv?what=sweeps", nil)
		req.Header.Set("Authorization", "Bearer "+fin)
		r := e.send(t, req)
		check(t, r.Status == 200 && strings.HasPrefix(r.Header.Get("Content-Type"), "text/csv"), "CSV: %d %s", r.Status, r.Header.Get("Content-Type"))
		body := string(r.Raw)
		check(t, strings.HasPrefix(body, "date,amount_usd,to_account") && strings.Contains(body, "TRF-0001"), "isi CSV: %s", truncate(body))
		req, _ = http.NewRequest(http.MethodGet, e.srv.URL+"/admin/api/sweep/export.csv", nil)
		req.Header.Set("Authorization", "Bearer "+fin)
		check(t, strings.HasPrefix(string(e.send(t, req).Raw), "month,collected_usd"), "CSV bulanan")
	})
}
