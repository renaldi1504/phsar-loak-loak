package qa

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"
	"time"
)

func (e *env) reconUpload(t *testing.T, token, name, csv string, fields ...string) resp {
	t.Helper()
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	for i := 0; i+1 < len(fields); i += 2 {
		_ = mw.WriteField(fields[i], fields[i+1])
	}
	fw, _ := mw.CreateFormFile("file", name)
	_, _ = fw.Write([]byte(csv))
	_ = mw.Close()
	req, _ := http.NewRequest(http.MethodPost, e.srv.URL+"/admin/api/reconciliation/batches", &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	return e.send(t, req)
}

// ---------------------------------------------------------------------------
// 18. Rekonsiliasi bank: ekspor PayWay (CSV) ↔ pembayaran & refund di Phsar
// ---------------------------------------------------------------------------

func TestReconciliation(t *testing.T) {
	e := getEnv(t)
	mod := newAdmin(t, e, "moderator")
	cs := newAdmin(t, e, "cs")
	fin := newAdmin(t, e, "finance")
	s := newSellable(t, e, mod, "Recon kettle", []string{"5.00"}, []int{50})
	buyer, _ := e.login(t)
	n := 0
	checkout := func(pay bool) (id, tran string) {
		t.Helper()
		n++
		c := e.call(t, "POST", "/v1/shop/checkouts", checkoutBody("k-recon-"+itoa(n)+"-x", "5.00", item(s.variants[0], 1)), buyer).want(201)
		if pay {
			e.call(t, "POST", "/dev/pay/"+c.str("id"), nil, buyer).want(200)
		}
		return c.str("id"), e.scalar(t, `SELECT tran_id FROM checkouts WHERE id = $1`, c.str("id"))
	}
	_, okTran := checkout(true)         // cocok
	_, wrongTran := checkout(true)      // jumlah berbeda
	pendID, pendTran := checkout(false) // dibayar di PayWay, belum lunas di Phsar
	refID, refTran := checkout(true)    // sudah direfund admin
	missID, missTran := checkout(true)  // lunas di Phsar, tidak ada di file
	e.sql(t, `UPDATE orders SET status = 'cancelled', refund_status = 'refunded', refund_amount = total, refund_ref = 'RF-RECON-1', refunded_at = now()
		WHERE checkout_id = $1`, refID)
	e.sql(t, `UPDATE checkouts SET provider = 'payway' WHERE id = $1`, missID)
	t.Cleanup(func() { e.sql(t, `UPDATE checkouts SET provider = 'fake' WHERE id = $1`, missID) })

	day := time.Now().In(time.FixedZone("ICT", 7*3600)).Format("02/01/2006")
	csv := "\xEF\xBB\xBFABA PayWay - Transaction report\n\n" +
		"Transaction ID;Date;Time;Amount;Currency;Status;APV;Type;Fee\n" +
		okTran + ";" + day + ";10:01;5.00;USD;APPROVED;111111;Payment;0.02\n" +
		wrongTran + ";" + day + ";10:02;4.50;USD;APPROVED;222222;Payment;0.02\n" +
		pendTran + ";" + day + ";10:03;5.00;USD;APPROVED;333333;Payment;0.02\n" +
		"ZZUNKNOWN99;" + day + ";10:04;\"1.234,50\";USD;APPROVED;444444;Payment;\n" +
		okTran + ";" + day + ";10:05;5.00;USD;APPROVED;111111;Payment;0.02\n" +
		"DECL001;" + day + ";10:06;9.00;USD;DECLINED;;Payment;\n" +
		refTran + ";" + day + ";11:00;-5.00;USD;REFUNDED;RF-RECON-1;Refund;\n" +
		okTran + ";" + day + ";11:05;(2.00);USD;APPROVED;RF-NOPE;Refund;\n"

	t.Run("izin & validasi file", func(t *testing.T) {
		e.reconUpload(t, cs, "payway.csv", csv).want(403)
		e.call(t, "GET", "/admin/api/reconciliation", nil, mod).want(403)
		e.reconUpload(t, fin, "payway.xlsx", csv).wantErr(422, "validation", "file")
		e.reconUpload(t, fin, "payway.csv", "name,amount\nfoo,1\n").wantErr(422, "validation", "file")
		e.reconUpload(t, fin, "payway.csv", "tran_id,amount\n\x00\x01").wantErr(422, "validation", "file")
		e.reconUpload(t, fin, "payway.csv", csv, "period_from", "2026-13-01").wantErr(422, "validation", "period_from")
	})

	var batch string
	result := func(r resp, tran, kind string) string {
		for i := 0; i < r.len("rows"); i++ {
			p := "rows." + itoa(i) + "."
			if r.str(p+"tran_id") == tran && (kind == "" || r.str(p+"kind") == kind) {
				return r.str(p + "result")
			}
		}
		return "-"
	}
	t.Run("unggah & cocokkan", func(t *testing.T) {
		r := e.reconUpload(t, fin, "payway.csv", csv).want(200)
		batch = r.str("batch.id")
		r.eq("batch.row_count", "8").eq("batch.collected", "1254.00").eq("batch.refunded", "7.00").eq("batch.fees", "0.08")
		all := e.call(t, "GET", "/admin/api/reconciliation/batches/"+batch+"?filter=all", nil, fin).want(200)
		for tran, want := range map[string]string{okTran: "ok", wrongTran: "amount_mismatch", pendTran: "not_paid_in_phsar", "ZZUNKNOWN99": "unknown",
			"DECL001": "ignored", missTran: "missing_in_statement"} {
			check(t, result(all, tran, "payment") == want, "%s harus %s, dapat %s", tran, want, result(all, tran, "payment"))
		}
		check(t, all.str("counts.duplicate") == "1", "tran ganda harus 'duplicate': %s", truncate(string(all.Raw)))
		check(t, all.str("counts.refund_ok") == "1" && all.str("counts.refund_unrecorded") == "1", "refund: %s", truncate(string(all.Raw)))
		check(t, strings.Contains(string(all.Raw), `"orders":["PH`), "baris harus menautkan nomor pesanan")
		open := e.call(t, "GET", "/admin/api/reconciliation/batches/"+batch, nil, fin).want(200)
		check(t, open.str("open") == itoa(open.len("rows")) && open.len("rows") >= 6, "filter default = selisih terbuka: %s", truncate(string(open.Raw)))
		e.call(t, "GET", "/admin/api/reconciliation/batches/"+batch+"?filter=bogus", nil, fin).wantErr(422, "validation", "filter")
		l := e.call(t, "GET", "/admin/api/reconciliation", nil, fin).want(200)
		check(t, l.str("batches.0.id") == batch, "batch harus tampil di daftar")
	})

	t.Run("selesaikan selisih & cocokkan ulang", func(t *testing.T) {
		all := e.call(t, "GET", "/admin/api/reconciliation/batches/"+batch+"?filter=all", nil, fin).want(200)
		var unknownRow, okRow, missRow string
		for i := 0; i < all.len("rows"); i++ {
			p := "rows." + itoa(i) + "."
			switch all.str(p + "result") {
			case "unknown":
				unknownRow = all.str(p + "id")
			case "ok":
				okRow = all.str(p + "id")
			case "missing_in_statement":
				if all.str(p+"tran_id") == missTran {
					missRow = all.str(p + "id")
				}
			}
		}
		e.call(t, "POST", "/admin/api/reconciliation/rows/"+unknownRow+"/resolve", map[string]string{"note": "x"}, fin).wantErr(422, "validation", "note")
		e.call(t, "POST", "/admin/api/reconciliation/rows/"+okRow+"/resolve", map[string]string{"note": "Nothing to do"}, fin).wantErr(409, "not_exception")
		e.call(t, "POST", "/admin/api/reconciliation/rows/"+unknownRow+"/resolve", map[string]string{"note": "Test payment from PayWay support"}, fin).want(200)
		e.call(t, "POST", "/admin/api/reconciliation/rows/"+missRow+"/resolve", map[string]string{"note": "Settled next day, in tomorrow's file"}, fin).want(200)

		// callback yang tertunda diproses → cocok setelah dicocokkan ulang; penyelesaian sisi Phsar tetap ada
		e.call(t, "POST", "/dev/pay/"+pendID, nil, buyer).want(200)
		r := e.call(t, "POST", "/admin/api/reconciliation/batches/"+batch+"/rematch?filter=all", nil, fin).want(200)
		check(t, result(r, pendTran, "payment") == "ok", "setelah dibayar harus ok: %s", result(r, pendTran, "payment"))
		res := e.call(t, "GET", "/admin/api/reconciliation/batches/"+batch+"?filter=resolved", nil, fin).want(200)
		check(t, res.len("rows") == 2, "dua penyelesaian harus bertahan setelah cocok ulang: %s", truncate(string(res.Raw)))
		check(t, e.scalar(t, `SELECT count(*) FROM audit_logs WHERE action = 'reconciliation.resolve'`) != "0", "penyelesaian harus tercatat di audit")
	})

	t.Run("hapus batch", func(t *testing.T) {
		e.call(t, "DELETE", "/admin/api/reconciliation/batches/"+batch, nil, fin).want(200)
		e.call(t, "GET", "/admin/api/reconciliation/batches/"+batch, nil, fin).want(404)
	})
}
