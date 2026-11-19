package qa

import (
	"strings"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

// ---------------------------------------------------------------------------
// Faktur komisi bulanan & pengaturan pajak
// ---------------------------------------------------------------------------

func TestInvoices(t *testing.T) {
	e := getEnv(t)
	mod := newAdmin(t, e, "moderator")
	fin := newAdmin(t, e, "finance")
	s := newSellable(t, e, mod, "Invoice bag", []string{"10.00"}, []int{50})
	other := newSellable(t, e, mod, "Other bag", []string{"10.00"}, []int{5})
	buyer, _ := e.login(t)
	t.Cleanup(func() { e.sql(t, `DELETE FROM platform_settings WHERE key = 'tax'`) })

	pp, _ := time.LoadLocation("Asia/Phnom_Penh")
	now := time.Now().In(pp)
	lastMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, pp).AddDate(0, -1, 0)
	period, cur, ym := lastMonth.Format("2006-01"), now.Format("2006-01"), lastMonth.Format("200601")
	mid := lastMonth.AddDate(0, 0, 14)

	// dua pesanan selesai bulan lalu + satu retur yang di-refund bulan lalu
	var orderIDs []string
	for i, k := range []string{"k-invoice-1-x", "k-invoice-2-x"} {
		c := e.call(t, "POST", "/v1/shop/checkouts", checkoutBody(k, "10.00", item(s.variants[0], 1)), buyer).want(201)
		e.call(t, "POST", "/dev/pay/"+c.str("id"), nil, buyer).want(200)
		id := e.scalar(t, `SELECT id::text FROM orders WHERE checkout_id = $1`, c.str("id"))
		e.sql(t, `UPDATE orders SET status = 'completed', completed_at = $2 WHERE id = $1`, id, mid.Add(time.Duration(i)*time.Hour))
		orderIDs = append(orderIDs, id)
	}
	e.sql(t, `INSERT INTO returns (return_no, order_id, buyer_id, store_id, status, reason, description, requested_amount, refund_amount,
		seller_deadline, refund_ref, refunded_at)
		SELECT 'RTINV-' || substr(md5(random()::text), 1, 6), o.id, o.buyer_id, o.store_id, 'refunded', 'damaged', 'Strap broke on arrival', 5, 5, now(), 'PW-1', $2
		FROM orders o WHERE o.id = $1`, orderIDs[0], mid.Add(48*time.Hour))
	c1 := invDec(e.scalar(t, `SELECT commission_amount::text FROM orders WHERE id = $1`, orderIDs[0]))
	c2 := invDec(e.scalar(t, `SELECT commission_amount::text FROM orders WHERE id = $1`, orderIDs[1]))
	back := decimal.NewFromInt(5).Mul(c1).Div(decimal.NewFromInt(10)).Round(2)
	net := c1.Add(c2).Sub(back)
	sub := net.Mul(decimal.NewFromInt(100)).Div(decimal.NewFromInt(110)).Round(2)

	t.Run("pengaturan pajak", func(t *testing.T) {
		e.call(t, "GET", "/admin/api/invoices", nil, mod).want(403)
		e.call(t, "GET", "/admin/api/invoices", nil, fin).want(200).eq("settings.vat_rate", "10").eq("settings.prefix", "PHS").eq("last_month", period)
		base := map[string]any{"legal_name": "Phsar Co., Ltd.", "vat_tin": "K001-901234567", "address": "Street 7, Sihanoukville", "vat_rate": "10", "vat_inclusive": true, "prefix": "PHT"}
		bad := func(k string, v any) map[string]any {
			m := map[string]any{}
			for a, b := range base {
				m[a] = b
			}
			m[k] = v
			return m
		}
		e.call(t, "PUT", "/admin/api/invoices/settings", bad("vat_rate", "50"), fin).wantErr(422, "validation", "vat_rate")
		e.call(t, "PUT", "/admin/api/invoices/settings", bad("prefix", "x"), fin).wantErr(422, "validation", "prefix")
		e.call(t, "PUT", "/admin/api/invoices/settings", bad("legal_name", ""), fin).wantErr(422, "validation", "legal_name")
		e.call(t, "PUT", "/admin/api/invoices/settings", base, mod).want(403)
		e.call(t, "PUT", "/admin/api/invoices/settings", base, fin).want(200).eq("prefix", "PHT")
	})

	var inv1 string
	t.Run("terbitkan bulan lalu: satu faktur per toko, angka benar, aman diulang", func(t *testing.T) {
		e.call(t, "POST", "/admin/api/invoices/generate", map[string]string{"period": cur}, fin).wantErr(422, "validation", "period")
		e.call(t, "POST", "/admin/api/invoices/generate", map[string]string{"period": "2026/10"}, fin).wantErr(422, "validation", "period")
		e.call(t, "POST", "/admin/api/invoices/generate", map[string]string{"period": period}, mod).want(403)
		r := e.call(t, "POST", "/admin/api/invoices/generate", map[string]string{"period": period}, fin).want(200)
		for i := 0; i < r.len("created"); i++ {
			if r.str("created."+itoa(i)+".store_id") == s.storeID {
				inv1 = r.str("created." + itoa(i) + ".id")
			}
		}
		check(t, inv1 != "", "faktur toko harus dibuat: %s", truncate(string(r.Raw)))
		v := e.call(t, "GET", "/admin/api/invoices/"+inv1, nil, fin).want(200)
		v.eq("order_count", "2").eq("commission", net.StringFixed(2)).eq("subtotal", sub.StringFixed(2)).eq("vat_amount", net.Sub(sub).StringFixed(2)).
			eq("total", net.StringFixed(2)).eq("status", "issued").eq("issuer.legal_name", "Phsar Co., Ltd.").eq("period", period)
		check(t, strings.HasPrefix(v.str("number"), "PHT-"+ym+"-"), "nomor faktur: %s", v.str("number"))
		check(t, v.len("lines") == 3 && v.str("lines.2.kind") == "return" && v.str("lines.2.commission") == back.Neg().StringFixed(2), "baris retur: %s", truncate(string(v.Raw)))
		check(t, e.scalar(t, `SELECT count(*) FROM invoices WHERE store_id = $1`, other.storeID) == "0", "toko tanpa transaksi tidak dapat faktur")

		// angka tidak ikut berubah kalau pesanan berubah; terbitkan ulang tidak membuat faktur ganda
		e.sql(t, `UPDATE orders SET completed_at = completed_at + interval '1 minute' WHERE id = $1`, orderIDs[1])
		again := e.call(t, "POST", "/admin/api/invoices/generate", map[string]string{"period": period}, fin).want(200)
		for i := 0; i < again.len("created"); i++ {
			check(t, again.str("created."+itoa(i)+".store_id") != s.storeID, "faktur ganda dibuat")
		}
		check(t, again.str("skipped") != "0", "toko yang sudah punya faktur harus dilewati")
	})

	t.Run("cetak & daftar untuk admin dan pemilik toko", func(t *testing.T) {
		h := e.call(t, "GET", "/admin/api/invoices/"+inv1+"/html", nil, fin).want(200)
		check(t, strings.Contains(string(h.Raw), "Phsar Co., Ltd.") && strings.Contains(string(h.Raw), "Tax invoice"), "HTML faktur")
		e.call(t, "GET", "/admin/api/invoices?status=bogus", nil, fin).wantErr(422, "validation", "status")
		e.call(t, "GET", "/admin/api/invoices?period=13-2026", nil, fin).wantErr(422, "validation", "period")
		l := e.call(t, "GET", "/admin/api/invoices?period="+period+"&q=PHT-", nil, fin).want(200)
		check(t, l.len("invoices") >= 1, "filter daftar")
		base := "/v1/seller/stores/" + s.storeID + "/payouts/invoices"
		sl := e.call(t, "GET", base, nil, s.token).want(200)
		check(t, sl.len("invoices") == 1 && sl.str("invoices.0.id") == inv1, "seller melihat fakturnya: %s", truncate(string(sl.Raw)))
		e.call(t, "GET", base+"/"+inv1+"/html", nil, s.token).want(200)
		e.call(t, "GET", "/v1/seller/stores/"+other.storeID+"/payouts/invoices/"+inv1+"/html", nil, other.token).want(404)
		x := e.call(t, "GET", base, nil, other.token)
		check(t, x.Status == 403 || x.Status == 404, "toko lain tidak boleh melihat: %d", x.Status)
	})

	t.Run("void lalu terbitkan pengganti dengan nomor baru", func(t *testing.T) {
		e.call(t, "POST", "/admin/api/invoices/"+inv1+"/void", map[string]string{"note": "no"}, fin).wantErr(422, "validation", "note")
		e.call(t, "POST", "/admin/api/invoices/"+inv1+"/void", map[string]string{"note": "Wrong legal name"}, mod).want(403)
		old := e.call(t, "POST", "/admin/api/invoices/"+inv1+"/void", map[string]string{"note": "Wrong legal name"}, fin).want(200)
		old.eq("status", "void")
		e.call(t, "POST", "/admin/api/invoices/"+inv1+"/void", map[string]string{"note": "Wrong legal name"}, fin).wantErr(409, "invoice_void")
		e.call(t, "POST", "/admin/api/invoices/00000000-0000-0000-0000-000000000000/void", map[string]string{"note": "Wrong legal name"}, fin).want(404)
		r := e.call(t, "POST", "/admin/api/invoices/generate", map[string]string{"period": period}, fin).want(200)
		var nid string
		for i := 0; i < r.len("created"); i++ {
			if r.str("created."+itoa(i)+".store_id") == s.storeID {
				nid = r.str("created." + itoa(i) + ".id")
			}
		}
		check(t, nid != "" && nid != inv1, "faktur pengganti harus dibuat")
		n := e.call(t, "GET", "/admin/api/invoices/"+nid, nil, fin).want(200)
		n.eq("replaces_number", old.str("number"))
		check(t, n.str("number") != old.str("number"), "nomor baru")
		check(t, strings.Contains(string(e.call(t, "GET", "/admin/api/invoices/"+inv1+"/html", nil, fin).want(200).Raw), "VOID"), "faktur batal bertanda VOID")
	})
}

func invDec(s string) decimal.Decimal { d, _ := decimal.NewFromString(s); return d }
