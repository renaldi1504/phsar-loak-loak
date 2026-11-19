package qa

import (
	"context"
	"io"
	"log/slog"
	"strings"
	"sync"
	"testing"
	"time"

	"phsar/internal/httpx"
)

// ---------------------------------------------------------------------------
// extras A: strike otomatis untuk pesanan batal karena telat (default mati; batas harian; bisa dicabut)
// ---------------------------------------------------------------------------

func TestExtrasAutoStrikes(t *testing.T) {
	e := getEnv(t)
	ctx := context.Background()
	sa := newAdmin(t, e, "super_admin")
	mod := newAdmin(t, e, "moderator")
	s := newSellable(t, e, mod, "Late kettle", []string{"10.00"}, []int{50})
	buyer, _ := e.login(t)
	shop := httpx.NewShop(e.cfg, slog.New(slog.NewTextHandler(io.Discard, nil)), e.pool)
	base := "/v1/seller/stores/" + s.storeID
	rules := func(on bool, p1, p2, days, max int) map[string]any {
		return map[string]any{"enabled": on, "seller_timeout_points": p1, "pack_timeout_points": p2, "expiry_days": days, "daily_max": max}
	}
	t.Cleanup(func() { e.call(t, "PUT", "/admin/api/auto-strikes", rules(false, 1, 1, 30, 1), sa) })
	n := 0
	pay := func() string {
		t.Helper()
		n++
		c := e.call(t, "POST", "/v1/shop/checkouts", checkoutBody("k-exas-"+itoa(n)+"-"+time.Now().Format("150405.000"), "10.00", item(s.variants[0], 1)), buyer).want(201)
		e.call(t, "POST", "/dev/pay/"+c.str("id"), nil, buyer).want(200)
		return e.scalar(t, `SELECT order_no FROM orders WHERE checkout_id = $1`, c.str("id"))
	}
	lateAccept := func() string {
		t.Helper()
		no := pay()
		e.sql(t, `UPDATE orders SET respond_by = now() - interval '1 minute' WHERE order_no = $1`, no)
		if k, err := shop.ExpireUnaccepted(ctx); err != nil || k < 1 {
			t.Fatalf("ExpireUnaccepted: %d %v", k, err)
		}
		return no
	}
	latePack := func() string {
		t.Helper()
		no := pay()
		e.call(t, "POST", base+"/orders/"+no+"/accept", nil, s.token).want(200)
		e.sql(t, `UPDATE orders SET ready_by = now() - interval '1 minute' WHERE order_no = $1`, no)
		if k, err := shop.ExpireUnpacked(ctx); err != nil || k < 1 {
			t.Fatalf("ExpireUnpacked: %d %v", k, err)
		}
		return no
	}
	strikes := func() string { return e.scalar(t, `SELECT count(*) FROM store_strikes WHERE store_id = $1`, s.storeID) }

	t.Run("hanya super admin; default mati; isian divalidasi", func(t *testing.T) {
		e.call(t, "GET", "/admin/api/auto-strikes", nil, mod).want(403)
		e.call(t, "PUT", "/admin/api/auto-strikes", rules(true, 1, 1, 30, 1), mod).want(403)
		e.call(t, "GET", "/admin/api/auto-strikes", nil, sa).want(200).eq("enabled", "false").eq("daily_max", "1")
		e.call(t, "PUT", "/admin/api/auto-strikes", rules(true, 4, 1, 30, 1), sa).wantErr(422, "validation", "seller_timeout_points")
		e.call(t, "PUT", "/admin/api/auto-strikes", rules(true, 1, 1, 0, 1), sa).wantErr(422, "validation", "expiry_days")
		e.call(t, "PUT", "/admin/api/auto-strikes", rules(true, 1, 1, 30, 0), sa).wantErr(422, "validation", "daily_max")
		e.call(t, "PUT", "/admin/api/auto-strikes", map[string]any{"enabled": true, "bogus": 1}, sa).want(400)
	})

	t.Run("mati → pesanan telat tetap batal + refund, tanpa strike", func(t *testing.T) {
		no := lateAccept()
		e.call(t, "GET", base+"/orders/"+no, nil, s.token).want(200).eq("status", "cancelled").eq("cancel_reason", "seller_timeout")
		check(t, strikes() == "0", "fitur mati tidak boleh memberi strike")
	})

	var first string
	t.Run("nyala → strike otomatis dengan poin, masa berlaku, notifikasi; batas harian", func(t *testing.T) {
		e.call(t, "PUT", "/admin/api/auto-strikes", rules(true, 1, 2, 10, 1), sa).want(200).eq("enabled", "true")
		no := lateAccept()
		e.call(t, "GET", base+"/orders/"+no, nil, s.token).want(200).eq("status", "cancelled")
		check(t, strikes() == "1", "harus ada 1 strike otomatis, dapat %s", strikes())
		first = e.scalar(t, `SELECT id::text FROM store_strikes WHERE store_id = $1`, s.storeID)
		check(t, e.scalar(t, `SELECT reason || ':' || points || ':' || auto || ':' || (created_by IS NULL) FROM store_strikes WHERE id = $1`, first) == "seller_timeout:1:true:true",
			"isi strike salah")
		d := e.scalar(t, `SELECT round(extract(epoch FROM expires_at - now()) / 86400)::text FROM store_strikes WHERE id = $1`, first)
		check(t, d == "10", "masa berlaku harus 10 hari, dapat %s", d)
		check(t, e.scalar(t, `SELECT count(*) FROM notifications WHERE store_id = $1 AND kind = 'store.strike_late_accept'`, s.storeID) != "0", "seller tidak diberi tahu")
		sk := e.call(t, "GET", base+"/strikes", nil, s.token).want(200)
		sk.eq("active_points", "1").eq("strikes.0.auto", "true").eq("strikes.0.order_no", no).eq("strikes.0.reason", "seller_timeout")
		// hari yang sama: batas 1 strike otomatis per hari → pesanan telat berikutnya tidak menambah strike
		latePack()
		check(t, strikes() == "1", "batas harian dilanggar, dapat %s", strikes())
	})

	t.Run("hari berikutnya: kemas telat (2 poin) → 3 poin → ditangguhkan; admin mencabut → aktif lagi", func(t *testing.T) {
		e.sql(t, `UPDATE store_strikes SET created_at = now() - interval '1 day' WHERE id = $1`, first)
		no := latePack()
		check(t, strikes() == "2", "harus ada strike kedua, dapat %s", strikes())
		check(t, e.scalar(t, `SELECT points::text FROM store_strikes WHERE store_id = $1 AND order_id = (SELECT id FROM orders WHERE order_no = $2)`, s.storeID, no) == "2",
			"poin kemas telat harus 2")
		check(t, e.scalar(t, `SELECT status::text FROM stores WHERE id = $1`, s.storeID) == "suspended", "3 poin aktif → toko ditangguhkan")
		h := e.call(t, "GET", "/admin/api/reports/stores/"+s.storeID+"/strikes", nil, mod).want(200)
		h.eq("strikes.0.auto", "true").eq("strikes.0.by", "")
		e.call(t, "POST", "/admin/api/reports/strikes/"+h.str("strikes.0.id")+"/revoke", map[string]string{"note": "Courier delay, not the seller"}, mod).want(200).eq("active_points", "1")
		check(t, e.scalar(t, `SELECT status::text FROM stores WHERE id = $1`, s.storeID) == "active", "cabut strike → toko aktif lagi")
	})
}

// ---------------------------------------------------------------------------
// extras B: tarik saldo manual oleh pemilik toko (default mati; minimal; N kali / 7 hari; aman dari klik ganda)
// ---------------------------------------------------------------------------

func TestExtrasManualPayout(t *testing.T) {
	e := getEnv(t)
	sa := newAdmin(t, e, "super_admin")
	fin := newAdmin(t, e, "finance")
	tok, sid := newActiveStore(t, e)
	staff, staffPhone := e.login(t)
	e.sql(t, `INSERT INTO store_members (store_id, user_id, role) SELECT $1, id, 'staff' FROM users WHERE phone = '+855' || substr($2, 2)`, sid, staffPhone)
	base := "/v1/seller/stores/" + sid + "/payouts"
	rules := func(on bool, min string, max int) map[string]any {
		return map[string]any{"enabled": on, "min_amount": min, "max_per_7_days": max}
	}
	t.Cleanup(func() {
		e.call(t, "PUT", "/admin/api/manual-payout", rules(false, "10.00", 1), sa)
		// saldo toko tes dinolkan supaya tes payout lain (yang menghitung semua toko) tidak terpengaruh
		e.sql(t, `UPDATE payouts SET status = 'cancelled', cancelled_at = now() WHERE store_id = $1 AND status = 'pending'`, sid)
		e.sql(t, `INSERT INTO seller_ledger (store_id, kind, amount, note) SELECT $1, 'adjustment', -sum(amount), 'test reset'
			FROM seller_ledger WHERE store_id = $1 HAVING sum(amount) <> 0`, sid)
	})
	credit := func(amount string) {
		e.sql(t, `INSERT INTO seller_ledger (store_id, kind, amount, note) VALUES ($1, 'adjustment', $2::numeric, 'test credit')`, sid, amount)
	}
	reason := func(r resp) string { return r.str("reason") }

	t.Run("pengaturan: hanya super admin; default mati; validasi", func(t *testing.T) {
		e.call(t, "GET", "/admin/api/manual-payout", nil, fin).want(403)
		e.call(t, "PUT", "/admin/api/manual-payout", rules(true, "5", 1), fin).want(403)
		e.call(t, "GET", "/admin/api/manual-payout", nil, sa).want(200).eq("enabled", "false").eq("min_amount", "10.00")
		e.call(t, "PUT", "/admin/api/manual-payout", rules(true, "0", 1), sa).wantErr(422, "validation", "min_amount")
		e.call(t, "PUT", "/admin/api/manual-payout", rules(true, "abc", 1), sa).wantErr(422, "validation", "min_amount")
		e.call(t, "PUT", "/admin/api/manual-payout", rules(true, "10", 0), sa).wantErr(422, "validation", "max_per_7_days")
	})

	t.Run("mati → tidak bisa tarik; tidak mengubah ringkasan saldo", func(t *testing.T) {
		credit("12.00")
		w := e.call(t, "GET", base+"/withdraw", nil, tok).want(200)
		w.eq("enabled", "false").eq("can", "false")
		check(t, reason(w) == "off", "alasan harus off, dapat %s", reason(w))
		r := e.call(t, "POST", base+"/withdraw", nil, tok).wantErr(409, "withdraw_blocked")
		check(t, strings.Contains(string(r.Raw), `"reason":"off"`), "alasan off tidak ada: %s", r.Raw)
		e.call(t, "GET", base, nil, tok).want(200).eq("available", "12.00").eq("pending_payout", "0.00")
	})

	t.Run("izin: hanya pemilik toko", func(t *testing.T) {
		e.call(t, "PUT", "/admin/api/manual-payout", rules(true, "10.00", 1), sa).want(200)
		st := e.call(t, "POST", base+"/withdraw", nil, staff)
		check(t, st.Status == 403 || st.Status == 404, "staf bisa menarik saldo: %d", st.Status)
		other, _ := e.login(t)
		o := e.call(t, "GET", base+"/withdraw", nil, other)
		check(t, o.Status == 403 || o.Status == 404, "orang lain bisa melihat: %d", o.Status)
	})

	t.Run("syarat: rekening, ditahan, minimal", func(t *testing.T) {
		w := e.call(t, "GET", base+"/withdraw", nil, tok).want(200)
		check(t, reason(w) == "no_account", "harus no_account, dapat %s", reason(w))
		e.call(t, "PUT", base+"/account", map[string]string{"bank": "aba", "account_name": "Dara Chan", "account_number": "000111222"}, tok).want(200)
		e.call(t, "PUT", base+"/account", map[string]string{"bank": "wing", "account_name": "Dara Chan", "account_number": "000111333"}, tok).want(200) // ganti → ditahan 48 jam
		w = e.call(t, "GET", base+"/withdraw", nil, tok).want(200)
		check(t, reason(w) == "account_hold" && w.str("hold_until") != "", "harus account_hold, dapat %s", reason(w))
		e.sql(t, `UPDATE seller_payout_accounts SET payout_hold_until = now() - interval '1 minute' WHERE store_id = $1 AND deleted_at IS NULL`, sid)
		e.call(t, "PUT", "/admin/api/manual-payout", rules(true, "20", 1), sa).want(200).eq("min_amount", "20.00")
		w = e.call(t, "GET", base+"/withdraw", nil, tok).want(200)
		check(t, reason(w) == "below_min", "harus below_min, dapat %s", reason(w))
		e.call(t, "POST", base+"/withdraw", nil, tok).wantErr(409, "withdraw_blocked")
		e.call(t, "PUT", "/admin/api/manual-payout", rules(true, "10", 1), sa).want(200)
		e.call(t, "GET", base+"/withdraw", nil, tok).want(200).eq("can", "true").eq("available", "12.00")
	})

	var pid string
	t.Run("tarik sekarang (klik ganda bersamaan) → tepat satu payout pending, label diminta seller", func(t *testing.T) {
		var wg sync.WaitGroup
		codes := make([]int, 2)
		for i := range codes {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				codes[i] = e.call(t, "POST", base+"/withdraw", nil, tok).Status
			}(i)
		}
		wg.Wait()
		check(t, (codes[0] == 200 && codes[1] == 409) || (codes[0] == 409 && codes[1] == 200), "harus satu 200 dan satu 409, dapat %v", codes)
		check(t, e.scalar(t, `SELECT count(*) FROM payouts WHERE store_id = $1`, sid) == "1", "payout dobel")
		e.call(t, "GET", base, nil, tok).want(200).eq("available", "0.00").eq("pending_payout", "12.00").eq("payouts.0.requested_by_seller", "true")
		w := e.call(t, "GET", base+"/withdraw", nil, tok).want(200)
		check(t, reason(w) == "pending_exists", "harus pending_exists, dapat %s", reason(w))
		check(t, e.scalar(t, `SELECT count(*) FROM audit_logs WHERE action = 'payout.request' AND store_id = $1`, sid) == "1", "permintaan tidak diaudit")
		l := e.call(t, "GET", "/admin/api/payouts?status=pending", nil, fin).want(200)
		found := false
		for i := 0; i < l.len("payouts"); i++ {
			if l.str("payouts."+itoa(i)+".store_id") == sid {
				found = true
				pid = l.str("payouts." + itoa(i) + ".id")
				check(t, l.str("payouts."+itoa(i)+".requested_by_seller") == "true", "admin tidak melihat label diminta seller")
				check(t, l.str("payouts."+itoa(i)+".account_number") == "000111333", "rekening salah")
			}
		}
		check(t, found, "payout manual tidak muncul di antrean admin")
		// jadwal admin tidak membuat payout kedua
		g := e.call(t, "POST", "/admin/api/payouts/generate", nil, fin).want(200)
		check(t, e.scalar(t, `SELECT count(*) FROM payouts WHERE store_id = $1`, sid) == "1", "jadwal membuat payout kedua: %s", g.Raw)
	})

	t.Run("admin tandai dibayar seperti biasa; batas N kali per 7 hari", func(t *testing.T) {
		e.call(t, "POST", "/admin/api/payouts/"+pid+"/paid", map[string]string{"reference": "ABA-MANUAL-1"}, fin).want(200)
		credit("15.00")
		w := e.call(t, "GET", base+"/withdraw", nil, tok).want(200)
		check(t, reason(w) == "limit" && w.str("next_at") != "", "harus limit, dapat %s", reason(w))
		w.eq("used_7_days", "1")
		e.call(t, "POST", base+"/withdraw", nil, tok).wantErr(409, "withdraw_blocked")
		e.call(t, "PUT", "/admin/api/manual-payout", rules(true, "10", 2), sa).want(200)
		e.call(t, "POST", base+"/withdraw", nil, tok).want(200).eq("amount", "15.00").eq("withdraw.reason", "pending_exists")
		// payout manual yang dibatalkan admin tidak dihitung ke batas
		p2 := e.scalar(t, `SELECT id::text FROM payouts WHERE store_id = $1 AND status = 'pending'`, sid)
		e.call(t, "POST", "/admin/api/payouts/"+p2+"/cancel", map[string]string{"note": "Wrong account"}, fin).want(200)
		e.call(t, "GET", base+"/withdraw", nil, tok).want(200).eq("used_7_days", "1").eq("can", "true").eq("available", "15.00")
	})
}
