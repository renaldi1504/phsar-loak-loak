package qa

import (
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// 18. Laporan pembeli & strike toko: 3 strike aktif → toko ditangguhkan; cabut → aktif lagi
// ---------------------------------------------------------------------------

func TestReportsAndStrikes(t *testing.T) {
	e := getEnv(t)
	mod := newAdmin(t, e, "moderator")
	cs := newAdmin(t, e, "cs")
	s := newSellable(t, e, mod, "Report watch", []string{"9.00"}, []int{5})
	b1, _ := e.login(t)
	b2, _ := e.login(t)
	rep := func(tok, target, id, reason, note string) resp {
		return e.call(t, "POST", "/v1/reports", map[string]string{"target": target, "id": id, "reason": reason, "note": note}, tok)
	}

	t.Run("pembeli melapor", func(t *testing.T) {
		rep("", "product", s.productID, "counterfeit", "").want(401)
		guest := e.call(t, "POST", "/v1/app/auth/guest", nil, "").want(200).str("access_token")
		rep(guest, "product", s.productID, "counterfeit", "").want(403)
		rep(b1, "user", s.productID, "counterfeit", "").wantErr(422, "validation", "target")
		rep(b1, "product", s.productID, "ugly", "").wantErr(422, "validation", "reason")
		rep(b1, "product", s.productID, "other", "x").wantErr(422, "validation", "note")
		rep(b1, "product", "00000000-0000-0000-0000-000000000000", "scam", "").want(404)
		rep(b1, "product", s.productID, "counterfeit", "Logo looks fake in the video").want(201)
		rep(b1, "product", s.productID, "scam", "").wantErr(409, "already_reported")
		rep(b2, "product", s.productID, "misleading", "").want(201)
		rep(b2, "store", s.storeID, "scam", "Asked me to pay outside Phsar").want(201)
	})

	t.Run("antrean admin dikelompokkan per target", func(t *testing.T) {
		e.call(t, "GET", "/admin/api/reports", nil, cs).want(403)
		r := e.call(t, "GET", "/admin/api/reports", nil, mod).want(200)
		found := false
		for i := 0; i < r.len("groups"); i++ {
			p := "groups." + itoa(i)
			if r.str(p+".target_id") == s.productID {
				found = true
				check(t, r.str(p+".count") == "2" && r.str(p+".product_name") == "Report watch", "grup produk: %s", truncate(string(r.Raw)))
			}
		}
		check(t, found, "grup produk tidak ada")
		d := e.call(t, "GET", "/admin/api/dashboard", nil, mod).want(200)
		check(t, strings.Contains(string(d.Raw), `"key":"reports"`), "antrean dashboard 'reports' tidak ada")
	})

	t.Run("larang produk + strike; 3 strike → toko ditangguhkan", func(t *testing.T) {
		body := map[string]any{"target": "product", "id": s.productID, "action": "ban_product", "note": "no", "also_strike": true}
		e.call(t, "POST", "/admin/api/reports/resolve", body, mod).wantErr(422, "validation", "note")
		body["note"] = "Counterfeit brand logo confirmed from video"
		e.call(t, "POST", "/admin/api/reports/resolve", body, mod).want(200)
		e.call(t, "POST", "/admin/api/reports/resolve", body, mod).wantErr(409, "no_open_reports")
		check(t, e.scalar(t, `SELECT status::text FROM products WHERE id = $1`, s.productID) == "banned", "produk harus banned")
		check(t, e.scalar(t, `SELECT count(*) FROM notifications WHERE store_id = $1 AND kind = 'store.strike'`, s.storeID) != "0", "seller harus diberi tahu strike")
		e.call(t, "POST", "/admin/api/reports/resolve", map[string]any{"target": "store", "id": s.storeID, "action": "ban_product", "note": "Not a product"}, mod).wantErr(422, "validation", "action")
		e.call(t, "POST", "/admin/api/reports/resolve", map[string]any{"target": "store", "id": s.storeID, "action": "strike", "note": "Off-platform payment request", "points": 2}, mod).want(200)
		check(t, e.scalar(t, `SELECT status::text FROM stores WHERE id = $1`, s.storeID) == "suspended", "3 poin strike → toko harus suspended")
		check(t, e.scalar(t, `SELECT count(*) FROM notifications WHERE store_id = $1 AND kind = 'store.suspended'`, s.storeID) != "0", "seller harus diberi tahu penangguhan")
		sk := e.call(t, "GET", "/v1/seller/stores/"+s.storeID+"/strikes", nil, s.token).want(200)
		sk.eq("active_points", "3")
		check(t, sk.str("strikes.0.by") == "", "seller tidak boleh melihat nama admin")
		e.call(t, "GET", "/v1/seller/stores/"+s.storeID+"/strikes", nil, b1).want(404)
	})

	t.Run("cabut strike → toko aktif lagi; strike manual", func(t *testing.T) {
		h := e.call(t, "GET", "/admin/api/reports/stores/"+s.storeID+"/strikes", nil, mod).want(200)
		var big string
		for i := 0; i < h.len("strikes"); i++ {
			if h.str("strikes."+itoa(i)+".points") == "2" {
				big = h.str("strikes." + itoa(i) + ".id")
			}
		}
		e.call(t, "POST", "/admin/api/reports/strikes/"+big+"/revoke", map[string]string{"note": "x"}, mod).wantErr(422, "validation", "note")
		r := e.call(t, "POST", "/admin/api/reports/strikes/"+big+"/revoke", map[string]string{"note": "Seller showed chat proof"}, mod).want(200)
		r.eq("active_points", "1").eq("store.status", "active")
		e.call(t, "POST", "/admin/api/reports/strikes/"+big+"/revoke", map[string]string{"note": "Seller showed chat proof"}, mod).wantErr(409, "already_revoked")
		check(t, e.scalar(t, `SELECT count(*) FROM notifications WHERE store_id = $1 AND kind = 'store.reinstated'`, s.storeID) != "0", "seller harus diberi tahu toko aktif lagi")
		e.call(t, "POST", "/admin/api/reports/strikes", map[string]any{"store_id": s.storeID, "reason": "bad", "note": "Prohibited item"}, mod).wantErr(422, "validation", "reason")
		e.call(t, "POST", "/admin/api/reports/strikes", map[string]any{"store_id": s.storeID, "reason": "prohibited", "note": "Sold a vape pen"}, mod).want(200).eq("active_points", "2")
		// strike kedaluwarsa tidak dihitung
		e.sql(t, `UPDATE store_strikes SET expires_at = now() - interval '1 day' WHERE store_id = $1`, s.storeID)
		e.call(t, "GET", "/admin/api/reports/stores/"+s.storeID+"/strikes", nil, mod).want(200).eq("active_points", "0")
		cl := e.call(t, "GET", "/admin/api/reports?status=closed", nil, mod).want(200)
		check(t, cl.len("groups") >= 1, "laporan yang ditutup harus tampil di Closed")
	})

	t.Run("batas harian", func(t *testing.T) {
		b3, _ := e.login(t)
		for i := 0; i < 10; i++ {
			e.sql(t, `INSERT INTO reports (reporter_id, target_type, store_id, reason, status) SELECT u.id, 'store', $1, 'other', 'dismissed'
				FROM users u JOIN sessions se ON se.user_id = u.id ORDER BY se.created_at DESC LIMIT 1`, s.storeID)
		}
		rep(b3, "store", s.storeID, "scam", "").wantErr(429, "report_limit")
	})
}
