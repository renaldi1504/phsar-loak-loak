package qa

import (
	"context"
	"io"
	"log/slog"
	"strings"
	"testing"
	"time"

	"phsar/internal/httpx"
)

// ---------------------------------------------------------------------------
// 13. Saldo seller & payout manual terjadwal
// ---------------------------------------------------------------------------

func TestPayouts(t *testing.T) {
	e := getEnv(t)
	ctx := context.Background()
	mod := newAdmin(t, e, "moderator")
	cs := newAdmin(t, e, "cs")
	fin := newAdmin(t, e, "finance")
	s := newSellable(t, e, mod, "Payout bag", []string{"6.00"}, []int{50})
	e.sql(t, `INSERT INTO store_pickup_points (store_id, label, address_line, lat, lng, contact_name, contact_phone, is_default, verified_at)
		VALUES ($1, 'Bag shop', 'Street 5', 10.61, 103.53, 'Lina', '012777666', true, now())`, s.storeID)
	courier, cphone := e.login(t)
	e.call(t, "POST", "/admin/api/couriers", map[string]any{"phone": cphone, "name": "Rith"}, cs).want(200)
	buyer, _ := e.login(t)
	svc := httpx.NewPayouts(e.cfg, slog.New(slog.NewTextHandler(io.Discard, nil)), e.pool)
	base := "/v1/seller/stores/" + s.storeID
	n := 0
	deliver := func() string { // pesanan $6 → terantar; seller_amount $5.52
		t.Helper()
		n++
		c := e.call(t, "POST", "/v1/shop/checkouts", checkoutBody("k-payout-"+itoa(n)+"-x", "6.00", item(s.variants[0], 1)), buyer).want(201)
		e.call(t, "POST", "/dev/pay/"+c.str("id"), nil, buyer).want(200)
		no := c.str("orders.0.order_no")
		e.call(t, "POST", base+"/orders/"+no+"/accept", nil, s.token).want(200)
		e.call(t, "POST", base+"/orders/"+no+"/ready", nil, s.token).want(200)
		task := e.scalar(t, `SELECT d.id::text FROM deliveries d JOIN orders o ON o.id = d.order_id WHERE o.order_no = $1`, no)
		e.call(t, "POST", "/v1/courier/tasks/"+task+"/claim", nil, courier).want(200)
		e.call(t, "POST", "/v1/courier/tasks/"+task+"/pickup", map[string]string{"code": no}, courier).want(200)
		e.call(t, "POST", "/v1/courier/tasks/"+task+"/deliver", nil, courier).want(200)
		return no
	}
	age := func(days int) { // pura-pura pesanan terantar beberapa hari lalu
		e.sql(t, `UPDATE deliveries SET delivered_at = now() - make_interval(days => $1) WHERE status = 'delivered'
			AND order_id IN (SELECT id FROM orders WHERE store_id = $2 AND status = 'delivered')`, days, s.storeID)
	}

	no1 := deliver()
	deliver()
	t.Run("saldo: terantar masih ditahan, lalu masuk saldo setelah masa tahan", func(t *testing.T) {
		r := e.call(t, "GET", base+"/payouts", nil, s.token).want(200)
		r.eq("available", "0.00").eq("on_hold", "11.04").eq("hold_days", "3")
		if n, err := svc.CompleteDue(ctx); err != nil || n != 0 {
			t.Fatalf("masih masa tahan, harus 0: %d %v", n, err)
		}
		age(4)
		if n, err := svc.CompleteDue(ctx); err != nil || n != 2 {
			t.Fatalf("CompleteDue: %d %v", n, err)
		}
		r = e.call(t, "GET", base+"/payouts", nil, s.token).want(200)
		r.eq("available", "11.04").eq("on_hold", "0.00").eq("ledger.0.kind", "order_earning")
		e.call(t, "GET", "/v1/shop/orders/"+no1, nil, buyer).want(200).eq("status", "completed")
		if n, _ := svc.CompleteDue(ctx); n != 0 {
			t.Fatal("pesanan diproses dua kali")
		}
	})
	t.Run("tanpa rekening tidak dibayar; rekening disimpan terenkripsi", func(t *testing.T) {
		r := e.call(t, "POST", "/admin/api/payouts/generate", nil, fin).want(200)
		r.eq("created", "0")
		check(t, strings.Contains(string(r.Raw), `"reason":"no_account"`), "alasan no_account tidak ada: %s", r.Raw)
		e.call(t, "PUT", base+"/payouts/account", map[string]string{"bank": "aba", "account_name": "Lina", "account_number": "12ab"}, s.token).wantErr(422, "validation", "account_number")
		a := e.call(t, "PUT", base+"/payouts/account", map[string]string{"bank": "aba", "account_name": "Lina Sok", "account_number": "000 123 456"}, s.token).want(200)
		a.eq("account_last4", "3456")
		check(t, a.str("payout_hold_until") == "", "rekening pertama tidak perlu ditahan")
		check(t, !strings.Contains(string(a.Raw), "000123456"), "nomor rekening lengkap dikirim ke seller")
		check(t, e.scalar(t, `SELECT position('000123456' in encode(account_number_enc, 'escape'))::text FROM seller_payout_accounts WHERE store_id = $1`, s.storeID) == "0",
			"nomor rekening tersimpan tanpa enkripsi")
		staff := e.call(t, "GET", base+"/payouts", nil, buyer)
		check(t, staff.Status == 403 || staff.Status == 404, "bukan anggota toko bisa melihat saldo: %d", staff.Status)
	})
	var pid string
	t.Run("admin membuat payout → transfer manual → tandai dibayar", func(t *testing.T) {
		e.call(t, "POST", "/admin/api/payouts/generate", nil, mod).want(403)
		r := e.call(t, "POST", "/admin/api/payouts/generate", nil, fin).want(200)
		r.eq("created", "1").eq("total", "11.04")
		e.call(t, "GET", base+"/payouts", nil, s.token).want(200).eq("available", "0.00").eq("pending_payout", "11.04")
		again := e.call(t, "POST", "/admin/api/payouts/generate", nil, fin).want(200)
		again.eq("created", "0")
		l := e.call(t, "GET", "/admin/api/payouts?status=pending", nil, fin).want(200)
		l.eq("payouts.0.account_number", "000123456").eq("payouts.0.amount", "11.04").eq("payouts.0.order_count", "2")
		pid = l.str("payouts.0.id")
		check(t, e.scalar(t, `SELECT count(*) FROM audit_logs WHERE action = 'payout_account.view'`) >= "1", "melihat nomor rekening tidak dicatat")
		e.call(t, "POST", "/admin/api/payouts/"+pid+"/paid", map[string]string{"reference": ""}, fin).wantErr(422, "validation", "reference")
		e.call(t, "POST", "/admin/api/payouts/"+pid+"/paid", map[string]string{"reference": "ABA-TRX-998877"}, fin).want(200)
		e.call(t, "POST", "/admin/api/payouts/"+pid+"/paid", map[string]string{"reference": "ABA-TRX-998877"}, fin).wantErr(409, "bad_status")
		h := e.call(t, "GET", base+"/payouts", nil, s.token).want(200)
		h.eq("payouts.0.status", "paid").eq("payouts.0.reference", "ABA-TRX-998877").eq("pending_payout", "0.00")
		check(t, h.str("payouts.0.account_number") == "", "nomor lengkap bocor di riwayat seller")
		check(t, e.scalar(t, `SELECT count(*) FROM notifications WHERE kind = 'payout.sent' AND store_id = $1`, s.storeID) == "1", "seller tidak diberi tahu")
	})
	t.Run("ganti rekening → ditahan 48 jam; batal payout → saldo kembali", func(t *testing.T) {
		deliver()
		age(5)
		if _, err := svc.CompleteDue(ctx); err != nil {
			t.Fatal(err)
		}
		a := e.call(t, "PUT", base+"/payouts/account", map[string]string{"bank": "acleda", "account_name": "Lina Sok", "account_number": "777888999"}, s.token).want(200)
		check(t, a.str("payout_hold_until") != "", "ganti rekening harus ditahan")
		r := e.call(t, "POST", "/admin/api/payouts/generate", nil, fin).want(200)
		r.eq("created", "0")
		check(t, strings.Contains(string(r.Raw), "account_hold"), "alasan account_hold tidak ada: %s", r.Raw)
		e.sql(t, `UPDATE seller_payout_accounts SET payout_hold_until = now() - interval '1 minute' WHERE store_id = $1 AND deleted_at IS NULL`, s.storeID)
		e.call(t, "POST", "/admin/api/payouts/generate", nil, fin).want(200).eq("created", "1").eq("total", "5.52")
		p2 := e.call(t, "GET", "/admin/api/payouts?status=pending", nil, fin).want(200)
		p2.eq("payouts.0.bank", "acleda").eq("payouts.0.account_number", "777888999")
		id := p2.str("payouts.0.id")
		e.call(t, "POST", "/admin/api/payouts/"+id+"/cancel", map[string]string{"note": ""}, fin).wantErr(422, "validation", "note")
		e.call(t, "POST", "/admin/api/payouts/"+id+"/cancel", map[string]string{"note": "Wrong account name"}, fin).want(200)
		e.call(t, "GET", base+"/payouts", nil, s.token).want(200).eq("available", "5.52").eq("pending_payout", "0.00")
	})
	t.Run("pengaturan jadwal oleh admin; jadwal otomatis sekali per hari", func(t *testing.T) {
		e.call(t, "PUT", "/admin/api/payout-settings", map[string]any{"frequency": "monthly", "weekday": 1, "hour": 9, "hold_days": 3, "min_amount": "5"}, fin).wantErr(422, "validation", "frequency")
		e.call(t, "PUT", "/admin/api/payout-settings", map[string]any{"frequency": "daily", "weekday": 0, "hour": 0, "hold_days": 2, "min_amount": "1"}, fin).want(200).
			eq("settings.frequency", "daily").eq("settings.min_amount", "1.00")
		e.call(t, "GET", base+"/payouts", nil, s.token).want(200).eq("frequency", "daily").eq("hold_days", "2")
		e.sql(t, `DELETE FROM payout_runs`)
		ran, err := svc.RunIfDue(ctx, time.Now())
		if err != nil || !ran {
			t.Fatalf("jadwal harian jam 00 harus jalan: %v %v", ran, err)
		}
		e.call(t, "GET", base+"/payouts", nil, s.token).want(200).eq("pending_payout", "5.52")
		ran, _ = svc.RunIfDue(ctx, time.Now())
		check(t, !ran, "jadwal berjalan dua kali di hari yang sama")
		e.call(t, "PUT", "/admin/api/payout-settings", map[string]any{"frequency": "weekly", "weekday": 1, "hour": 9, "hold_days": 3, "min_amount": "5.00"}, fin).want(200)
	})
}
