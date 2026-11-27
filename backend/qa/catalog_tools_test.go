package qa

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"phsar/internal/notify"
	"phsar/internal/seller"
)

// ---------------------------------------------------------------------------
// Alat katalog untuk toko dengan banyak SKU (migrasi 00052):
// A. SKU otomatis rapi · B. tabel varian, filter, ubah cepat, aksi massal · C. CSV · D. peringatan stok menipis
// ---------------------------------------------------------------------------

// catalogStore: toko aktif dengan nama tertentu (kode SKU diturunkan dari nama).
func catalogStore(t *testing.T, e *env, name string) (token, storeID string) {
	t.Helper()
	token, _ = e.login(t)
	slug := fmt.Sprintf("qa-catalog-%d", catalogSlugN.Add(1)) // slug unik: tes lain memakai nama toko yang sama
	st := e.call(t, "POST", "/v1/seller/stores", map[string]any{"name": name, "slug": slug, "about": "Catalog tools", "about_lang": "en"}, token).want(201)
	storeID = st.str("id")
	e.sql(t, `UPDATE stores SET status = 'active' WHERE id = $1`, storeID)
	quietStore(t, e, storeID)
	return token, storeID
}

var catalogSlugN atomic.Int64

// quietStore: tes lain menguras antrean terjemahan bersama (maks. 50 langkah) — toko & produk tes ini tidak perlu
// diterjemahkan, jadi tugasnya dilewati supaya antrean tidak membengkak.
func quietStore(t *testing.T, e *env, storeID string) {
	t.Helper()
	e.sql(t, `UPDATE translation_jobs SET status = 'skipped', done_at = now() WHERE status = 'queued'
		AND ((entity_table = 'store_translations' AND entity_id = $1::text)
		  OR (entity_table = 'product_translations' AND entity_id IN (SELECT id::text FROM products WHERE store_id = $1::uuid)))`, storeID)
}

type catalogKit struct {
	e         *env
	tok, sid  string
	base      string
	jpg, mp4  string
	moderator string
}

func newCatalogKit(t *testing.T, e *env, name string) *catalogKit {
	t.Helper()
	tok, sid := catalogStore(t, e, name)
	k := &catalogKit{e: e, tok: tok, sid: sid, base: "/v1/seller/stores/" + sid}
	k.jpg = e.upload(t, tok, "product_media", jpegBytes).want(201).str("key")
	k.mp4 = e.upload(t, tok, "product_media", mp4Bytes).want(201).str("key")
	return k
}

// product: body produk. opts = nilai opsi tunggal per varian ("" = tanpa opsi); skus kosong = otomatis.
func (k *catalogKit) product(name, category string, values []string, prices []string, stocks []int, skus []string, submit bool) map[string]any {
	body := map[string]any{"category_id": category, "source_lang": "en", "name": name, "description": "Catalog tools test product.",
		"condition": "new", "submit": submit,
		"media": []any{map[string]any{"key": k.mp4, "kind": "video"}, map[string]any{"key": k.jpg, "kind": "photo"}}}
	vars := []any{}
	if len(values) > 0 {
		vals := []any{}
		for _, v := range values {
			vals = append(vals, map[string]any{"value": v})
		}
		body["options"] = []any{map[string]any{"kind": "custom", "name": "Type", "values": vals}}
	}
	for i := range prices {
		v := map[string]any{"price": prices[i], "stock_qty": stocks[i]}
		if len(values) > 0 {
			v["options"] = []string{values[i]}
		}
		if i < len(skus) && skus[i] != "" {
			v["sku"] = skus[i]
		}
		vars = append(vars, v)
	}
	body["variants"] = vars
	return body
}

func (k *catalogKit) create(t *testing.T, body map[string]any) resp {
	t.Helper()
	r := k.e.call(t, "POST", k.base+"/products", body, k.tok).want(201)
	quietStore(t, k.e, k.sid)
	return r
}

// live: produk tayang (dikirim lalu disetujui moderator).
func (k *catalogKit) live(t *testing.T, name string, prices []string, stocks []int) resp {
	t.Helper()
	if k.moderator == "" {
		k.moderator = newAdmin(t, k.e, "moderator")
	}
	values := []string{}
	for i := range prices {
		values = append(values, []string{"Alpha", "Beta", "Gamma", "Delta"}[i])
	}
	p := k.create(t, k.product(name, "womens-tops", values, prices, stocks, nil, true))
	k.e.call(t, "POST", "/admin/api/products/"+p.str("id")+"/decision", map[string]string{"decision": "approve"}, k.moderator).want(200)
	return k.e.call(t, "GET", k.base+"/products/"+p.str("id"), nil, k.tok).want(200).eq("moderation_status", "approved")
}

func skus(p resp) []string {
	out := []string{}
	for i := 0; i < p.len("variants"); i++ {
		out = append(out, p.str("variants."+itoa(i)+".sku"))
	}
	return out
}

func hasAll(got []string, want ...string) bool {
	set := map[string]bool{}
	for _, g := range got {
		set[g] = true
	}
	for _, w := range want {
		if !set[w] {
			return false
		}
	}
	return true
}

func (k *catalogKit) importCSV(t *testing.T, csv string, apply bool) resp {
	t.Helper()
	path := k.base + "/products/import"
	if apply {
		path += "?apply=1"
	}
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	fw, _ := mw.CreateFormFile("file", "skus.csv")
	_, _ = fw.Write([]byte(csv))
	_ = mw.Close()
	req, _ := http.NewRequest(http.MethodPost, k.e.srv.URL+path, &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+k.tok)
	return k.e.send(t, req)
}

func TestCatalogSKU(t *testing.T) {
	e := getEnv(t)

	t.Run("kode SKU toko: dari nama, sama di Go dan SQL", func(t *testing.T) {
		for name, want := range map[string]string{"Skin Vibe": "SKV", "Beauty": "BEA", "Phnom Penh Phone Shop": "PPP", "K-Pop 24": "KP2", "A": ""} {
			id := e.scalar(t, `SELECT gen_random_uuid()::text`)
			sql := e.scalar(t, `SELECT phsar_sku_prefix($1, $2::uuid)`, name, id)
			goP := seller.SKUPrefix(name, id)
			check(t, sql == goP, "%q: SQL %s ≠ Go %s", name, sql, goP)
			if want != "" {
				check(t, sql == want, "%q → %s, harusnya %s", name, sql, want)
			}
		}
		id := e.scalar(t, `SELECT gen_random_uuid()::text`)
		km := seller.SKUPrefix("ហាងស្អាត", id)
		check(t, km == "S"+strings.ToUpper(strings.ReplaceAll(id, "-", "")[:3]), "nama Khmer → S + id: %s", km)
	})

	k := newCatalogKit(t, e, "Skin Vibe")
	t.Run("toko baru: kode SKU & pratinjau", func(t *testing.T) {
		e.call(t, "GET", k.base, nil, k.tok).want(200).eq("sku_prefix", "SKV").eq("sku_next", "SKV-00001").eq("low_stock_threshold", "3")
	})

	var p1 resp
	t.Run("SKU otomatis: kode toko + nomor urut + nilai varian; SKU yang diisi tetap", func(t *testing.T) {
		p1 = k.create(t, k.product("Glow serum", "womens-tops", []string{"Red", "Light blue", "Black"}, []string{"5", "5", "6"}, []int{10, 4, 0},
			[]string{"", "", "TEE-BLK-M"}, false))
		check(t, hasAll(skus(p1), "SKV-00001-RED", "SKV-00001-LIGHTBLUE", "TEE-BLK-M"), "SKU: %v", skus(p1))
		p1.eq("sku_base", "SKV-00001")
		e.call(t, "GET", k.base, nil, k.tok).want(200).eq("sku_next", "SKV-00002")
	})
	t.Run("nilai varian tanpa huruf latin → -1, -2; produk tanpa varian → dasar saja", func(t *testing.T) {
		p := k.create(t, k.product("Khmer scarf", "womens-tops", []string{"ក្រហម", "ខៀវ"}, []string{"3", "3"}, []int{1, 1}, nil, false))
		check(t, hasAll(skus(p), "SKV-00002-1", "SKV-00002-2"), "SKU: %v", skus(p))
		p = k.create(t, k.product("Plain mug", "womens-tops", nil, []string{"4"}, []int{2}, nil, false))
		check(t, hasAll(skus(p), "SKV-00003"), "SKU: %v", skus(p))
	})
	t.Run("bentrok dengan SKU yang sudah ada → akhiran -1", func(t *testing.T) {
		k.create(t, k.product("Manual sku", "womens-tops", nil, []string{"4"}, []int{2}, []string{"SKV-00004"}, false))
		p := k.create(t, k.product("Next mug", "womens-tops", nil, []string{"4"}, []int{2}, nil, false))
		check(t, hasAll(skus(p), "SKV-00004-1"), "SKU: %v", skus(p))
	})
	t.Run("simpan ulang: SKU lama tidak berubah, nomor produk tetap; SKU dikosongkan → dibuat lagi", func(t *testing.T) {
		pid := p1.str("id")
		body := k.product("Glow serum", "womens-tops", []string{"Red", "Light blue", "Black", "Green"}, []string{"5", "5", "6", "7"}, []int{10, 4, 0, 1}, nil, false)
		bySKU := map[string]string{}
		for i := 0; i < p1.len("variants"); i++ {
			bySKU[p1.str("variants."+itoa(i)+".options.0")] = p1.str("variants." + itoa(i) + ".sku")
		}
		vars := body["variants"].([]any)
		vars[0].(map[string]any)["sku"] = bySKU["Red"]
		vars[2].(map[string]any)["sku"] = bySKU["Black"]
		// Light blue dikosongkan, Green baru
		r := e.call(t, "PUT", k.base+"/products/"+pid, body, k.tok).want(200)
		quietStore(t, e, k.sid)
		check(t, hasAll(skus(r), "SKV-00001-RED", "SKV-00001-LIGHTBLUE", "TEE-BLK-M", "SKV-00001-GREEN"), "SKU: %v", skus(r))
		r.eq("sku_base", "SKV-00001")
		check(t, e.scalar(t, `SELECT last_no FROM store_sku_counters WHERE store_id = $1`, k.sid) == "4", "nomor tidak diambil ulang")
	})
	t.Run("kode SKU bisa diganti pemilik (untuk SKU baru saja)", func(t *testing.T) {
		e.call(t, "PATCH", k.base, map[string]any{"sku_prefix": "x"}, k.tok).wantErr(422, "validation", "sku_prefix")
		e.call(t, "PATCH", k.base, map[string]any{"sku_prefix": "TOOLONG"}, k.tok).wantErr(422, "validation", "sku_prefix")
		e.call(t, "PATCH", k.base, map[string]any{"sku_prefix": "skn2"}, k.tok).want(200).eq("sku_prefix", "SKN2").eq("sku_next", "SKN2-00005")
		p := k.create(t, k.product("After rename", "womens-tops", nil, []string{"4"}, []int{2}, nil, false))
		check(t, hasAll(skus(p), "SKN2-00005"), "SKU: %v", skus(p))
		check(t, e.scalar(t, `SELECT sku FROM product_variants WHERE product_id = $1 AND sku LIKE 'SKV-00001-RED'`, p1.str("id")) == "SKV-00001-RED", "SKU lama tetap")
		// boleh juga saat menunggu verifikasi (bukan bagian profil yang diverifikasi)
		e.sql(t, `UPDATE stores SET status = 'pending_review' WHERE id = $1`, k.sid)
		e.call(t, "PATCH", k.base, map[string]any{"sku_prefix": "SKV", "low_stock_threshold": 5}, k.tok).want(200).eq("low_stock_threshold", "5")
		e.call(t, "PATCH", k.base, map[string]any{"name": "Skin Vibe 2"}, k.tok).wantErr(409, "store_under_review")
		e.sql(t, `UPDATE stores SET status = 'active' WHERE id = $1`, k.sid)
		e.call(t, "PATCH", k.base, map[string]any{"low_stock_threshold": 3}, k.tok).want(200)
		e.call(t, "PATCH", k.base, map[string]any{"low_stock_threshold": -1}, k.tok).wantErr(422, "validation", "low_stock_threshold")
	})
	t.Run("nomor urut aman dipakai bersamaan", func(t *testing.T) {
		var wg sync.WaitGroup
		var mu sync.Mutex
		statuses := []int{}
		for i := 0; i < 6; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				r := e.call(t, "POST", k.base+"/products", k.product(fmt.Sprintf("Parallel %d", i), "womens-tops", nil, []string{"4"}, []int{1}, nil, false), k.tok)
				mu.Lock()
				statuses = append(statuses, r.Status)
				mu.Unlock()
			}(i)
		}
		wg.Wait()
		quietStore(t, e, k.sid)
		for _, s := range statuses {
			check(t, s == 201, "status %d", s)
		}
		dups := e.scalar(t, `SELECT count(*) - count(DISTINCT sku_no) FROM products WHERE store_id = $1 AND sku_no IS NOT NULL`, k.sid)
		check(t, dups == "0", "nomor produk dobel: %s", dups)
		check(t, e.scalar(t, `SELECT last_no FROM store_sku_counters WHERE store_id = $1`, k.sid) == "11", "penghitung = 11")
	})
	t.Run("toko lama mendapat kode dari migrasi; toko baru dari trigger", func(t *testing.T) {
		check(t, e.scalar(t, `SELECT count(*) FROM stores WHERE sku_prefix IS NULL`) == "0", "semua toko punya kode SKU")
	})
}

func TestCatalogTable(t *testing.T) {
	e := getEnv(t)
	k := newCatalogKit(t, e, "Table Shop")
	other := e.scalar(t, `SELECT c.id FROM categories c WHERE c.parent_id IS NOT NULL AND c.id <> 'womens-tops' AND c.deleted_at IS NULL
		AND NOT EXISTS (SELECT 1 FROM categories x WHERE x.parent_id = c.id) ORDER BY c.id LIMIT 1`)
	otherParent := e.scalar(t, `SELECT parent_id FROM categories WHERE id = $1`, other)
	for i := 0; i < 30; i++ {
		cat := "womens-tops"
		if i%3 == 0 {
			cat = other
		}
		k.create(t, k.product(fmt.Sprintf("Bulk item %02d", i), cat, []string{"Red", "Blue"}, []string{"5", "6"}, []int{i, 10}, nil, false))
	}
	list := func(q string) resp {
		return e.call(t, "GET", k.base+"/products?page=1&"+strings.ReplaceAll(q, " ", "+"), nil, k.tok).want(200)
	}
	vlist := func(q string) resp {
		return e.call(t, "GET", k.base+"/variants?page=1&"+strings.ReplaceAll(q, " ", "+"), nil, k.tok).want(200)
	}

	t.Run("daftar produk: ukuran halaman 24/50/100, filter kategori & stok", func(t *testing.T) {
		check(t, list("size=24").len("products") == 24, "24 per halaman")
		list("size=50").eq("page.total", "30").eq("page.pages", "1")
		check(t, list("size=100").len("products") == 30, "100 per halaman")
		list("size=50&category="+other).eq("page.total", "10")
		list("size=50&category="+otherParent).eq("page.total", "10") // induk ikut sub-kategori
		list("size=50&category=womens-tops").eq("page.total", "20")
		list("size=50&stock=out").eq("page.total", "1")          // Bulk item 00 Red = 0
		list("size=50&stock=low").eq("page.total", "3")          // stok 1, 2, 3
		list("size=50&category=bad;drop").eq("page.total", "30") // id tidak sah diabaikan
	})
	t.Run("tabel varian: satu baris per SKU, label varian, ambang, filter", func(t *testing.T) {
		r := vlist("size=50&sort=sku")
		r.eq("page.total", "60").eq("page.size", "50")
		check(t, r.len("variants") == 50, "50 baris")
		r = vlist("size=100&q=Bulk item 05")
		r.eq("page.total", "2")
		check(t, r.str("variants.0.label") != "" && r.str("variants.0.threshold") == "3", "label & ambang: %s", string(r.Raw)[:200])
		vlist("size=100&stock=low").eq("page.total", "3")
		vlist("size=100&stock=out&category=womens-tops").eq("page.total", "0")
		vlist("size=100&stock=out").eq("page.total", "1").eq("variants.0.available", "0")
	})

	var row resp
	t.Run("ubah cepat per baris: harga, harga coret, stok, ambang + tercatat", func(t *testing.T) {
		row = vlist("size=1&q=Bulk item 07&sort=sku")
		vid := row.str("variants.0.id")
		r := e.call(t, "PATCH", k.base+"/variants/"+vid, map[string]any{"price": "4.5", "compare_at_price": "9", "stock_qty": 12}, k.tok).want(200)
		r.eq("variant.price", "4.50").eq("variant.compare_at_price", "9.00").eq("variant.stock_qty", "12").eq("review", "false")
		check(t, e.scalar(t, `SELECT count(*) FROM inventory_movements WHERE variant_id = $1 AND reason = 'adjustment'`, vid) == "1", "pergerakan stok")
		check(t, e.scalar(t, `SELECT count(*) FROM audit_logs WHERE entity_table = 'product_variants' AND entity_id = $1 AND actor_type = 'seller'`, vid) != "0", "audit")
		e.call(t, "PATCH", k.base+"/variants/"+vid, map[string]any{"compare_at_price": "4"}, k.tok).wantErr(422, "validation", "compare_at_price")
		e.call(t, "PATCH", k.base+"/variants/"+vid, map[string]any{"price": "0"}, k.tok).wantErr(422, "validation", "price")
		e.call(t, "PATCH", k.base+"/variants/"+vid, map[string]any{"price": "abc"}, k.tok).wantErr(422, "validation", "price")
		e.call(t, "PATCH", k.base+"/variants/"+vid, map[string]any{"stock_qty": -2}, k.tok).wantErr(422, "validation", "stock_qty")
		e.call(t, "PATCH", k.base+"/variants/"+vid, map[string]any{"compare_at_price": ""}, k.tok).want(200).eq("variant.compare_at_price", "")
		e.call(t, "PATCH", k.base+"/variants/"+vid, map[string]any{"low_stock_threshold": 8}, k.tok).want(200).eq("variant.threshold", "8").eq("variant.low_stock_threshold", "8")
		e.call(t, "PATCH", k.base+"/variants/"+vid, map[string]any{"low_stock_threshold": nil}, k.tok).want(200).eq("variant.threshold", "3")
		e.sql(t, `UPDATE product_variants SET reserved_qty = 5 WHERE id = $1`, vid)
		e.call(t, "PATCH", k.base+"/variants/"+vid, map[string]any{"stock_qty": 4}, k.tok).wantErr(422, "validation", "stock_qty")
		e.sql(t, `UPDATE product_variants SET reserved_qty = 0 WHERE id = $1`, vid)
		e.call(t, "PATCH", k.base+"/variants/00000000-0000-0000-0000-000000000000", map[string]any{"stock_qty": 1}, k.tok).want(404)
		e.call(t, "PATCH", k.base+"/variants/not-a-uuid", map[string]any{"stock_qty": 1}, k.tok).want(404)
		// toko lain tidak bisa mengubah varian ini
		o := newCatalogKit(t, e, "Rival Catalog Shop")
		e.call(t, "PATCH", o.base+"/variants/"+vid, map[string]any{"stock_qty": 1}, o.tok).want(404)
	})

	t.Run("produk tayang: harga turun jauh → antrean moderasi (drop_pct), turun sedikit tetap tayang", func(t *testing.T) {
		p := k.live(t, "Live jacket", []string{"20", "22"}, []int{5, 5})
		vid := p.str("variants.0.id")
		e.call(t, "PATCH", k.base+"/variants/"+vid, map[string]any{"price": "18"}, k.tok).want(200).eq("review", "false").eq("variant.moderation_status", "approved")
		r := e.call(t, "PATCH", k.base+"/variants/"+vid, map[string]any{"price": "5"}, k.tok).want(200).eq("review", "true").eq("variant.moderation_status", "pending_review")
		_ = r
		reasons := e.scalar(t, `SELECT auto_flags->'reasons' FROM products WHERE id = $1`, p.str("id"))
		check(t, strings.Contains(reasons, "price_drop"), "alasan price_drop: %s", reasons)
		// jalur lama (PATCH /products/{id}/variants) juga dicek
		q := k.live(t, "Live coat", []string{"30"}, []int{5})
		e.call(t, "PATCH", k.base+"/products/"+q.str("id")+"/variants", map[string]any{"variants": []any{
			map[string]any{"id": q.str("variants.0.id"), "price": "9", "stock_qty": 5, "status": "active"}}}, k.tok).want(200).eq("moderation_status", "pending_review")
	})

	t.Run("aksi massal: arsip, aktifkan, hapus draf, stok, harga %", func(t *testing.T) {
		live1 := k.live(t, "Bulk live A", []string{"10"}, []int{5})
		live2 := k.live(t, "Bulk live B", []string{"10", "12"}, []int{5, 1})
		draft := list("size=1&q=Bulk item 29").str("products.0.id")
		ids := []string{live1.str("id"), live2.str("id"), draft}
		r := e.call(t, "POST", k.base+"/products/bulk", map[string]any{"action": "archive", "product_ids": ids}, k.tok).want(200)
		r.eq("done", "2").eq("skipped.0.id", draft).eq("skipped.0.code", "not_allowed")
		r = e.call(t, "POST", k.base+"/products/bulk", map[string]any{"action": "activate", "product_ids": ids}, k.tok).want(200).eq("done", "2")
		r = e.call(t, "POST", k.base+"/products/bulk", map[string]any{"action": "delete_drafts", "product_ids": ids}, k.tok).want(200)
		r.eq("done", "1")
		check(t, e.scalar(t, `SELECT deleted_at IS NOT NULL FROM products WHERE id = $1`, draft) == "true", "draf terhapus")
		check(t, e.scalar(t, `SELECT deleted_at IS NULL FROM products WHERE id = $1`, live1.str("id")) == "true", "produk tayang tidak terhapus")

		ids = ids[:2]
		e.call(t, "POST", k.base+"/products/bulk", map[string]any{"action": "stock", "stock_mode": "set", "stock": 20, "product_ids": ids}, k.tok).want(200).eq("done", "2")
		check(t, e.scalar(t, `SELECT string_agg(stock_qty::text, ',') FROM product_variants WHERE product_id = $1 AND deleted_at IS NULL`, live2.str("id")) == "20,20", "stok diset 20")
		e.call(t, "POST", k.base+"/products/bulk", map[string]any{"action": "stock", "stock_mode": "add", "stock": -25, "product_ids": ids}, k.tok).want(200)
		check(t, e.scalar(t, `SELECT sum(stock_qty) FROM product_variants WHERE product_id = ANY($1::uuid[]) AND deleted_at IS NULL`, ids) == "0", "dikurangi tidak di bawah 0")
		e.call(t, "POST", k.base+"/products/bulk", map[string]any{"action": "stock", "stock_mode": "add", "stock": 7, "product_ids": ids}, k.tok).want(200)

		e.sql(t, `UPDATE product_variants SET compare_at_price = 15 WHERE product_id = $1`, live1.str("id"))
		r = e.call(t, "POST", k.base+"/products/bulk", map[string]any{"action": "price_pct", "pct": "10", "product_ids": ids}, k.tok).want(200)
		r.eq("done", "2").eq("review", "0")
		check(t, e.scalar(t, `SELECT price || '/' || compare_at_price FROM product_variants WHERE product_id = $1`, live1.str("id")) == "11.00/16.50", "harga & harga coret naik 10%%")
		r = e.call(t, "POST", k.base+"/products/bulk", map[string]any{"action": "price_pct", "pct": "-70", "product_ids": ids}, k.tok).want(200)
		r.eq("done", "2").eq("review", "2")
		check(t, e.scalar(t, `SELECT moderation_status::text FROM products WHERE id = $1`, live1.str("id")) == "pending_review", "turun 70%% → review")

		e.call(t, "POST", k.base+"/products/bulk", map[string]any{"action": "price_pct", "pct": "-100", "product_ids": ids}, k.tok).wantErr(422, "validation", "pct")
		e.call(t, "POST", k.base+"/products/bulk", map[string]any{"action": "explode", "product_ids": ids}, k.tok).wantErr(422, "validation", "action")
		e.call(t, "POST", k.base+"/products/bulk", map[string]any{"action": "archive", "product_ids": []string{}}, k.tok).wantErr(422, "validation", "product_ids")
		r = e.call(t, "POST", k.base+"/products/bulk", map[string]any{"action": "archive", "product_ids": []string{"nope", "00000000-0000-0000-0000-000000000000"}}, k.tok).want(200)
		r.eq("done", "0").eq("skipped.0.code", "not_found").eq("skipped.1.code", "not_found")
	})
}

func TestCatalogCSV(t *testing.T) {
	e := getEnv(t)
	k := newCatalogKit(t, e, "Csvstore")
	a := k.create(t, k.product("Csv shirt", "womens-tops", []string{"Red", "Blue", "Green"}, []string{"5", "6", "7"}, []int{10, 4, 1}, []string{"", "", "-HACK"}, false))
	live := k.live(t, "Csv live", []string{"40"}, []int{3})
	red := "CSV-00001-RED"
	liveSKU := live.str("variants.0.sku")
	check(t, hasAll(skus(a), red), "SKU: %v", skus(a))

	t.Run("ekspor: BOM UTF-8, judul kolom, semua SKU, aman dari rumus Excel", func(t *testing.T) {
		r := e.call(t, "GET", k.base+"/products/export.csv", nil, k.tok).want(200)
		check(t, bytes.HasPrefix(r.Raw, []byte("\xEF\xBB\xBF")), "BOM")
		check(t, strings.HasPrefix(r.Header.Get("Content-Type"), "text/csv"), "content-type %s", r.Header.Get("Content-Type"))
		check(t, strings.Contains(r.Header.Get("Content-Disposition"), "phsar-csv-skus-"), "nama file: %s", r.Header.Get("Content-Disposition"))
		body := string(r.Raw[3:])
		check(t, strings.HasPrefix(body, "sku,product,variant,price,compare_at_price,stock,status\n"), "judul: %s", body[:60])
		check(t, strings.Contains(body, red+",Csv shirt,Red,5.00,,10,draft"), "baris: %s", body)
		check(t, strings.Contains(body, "'-HACK"), "rumus dinetralkan")
		check(t, strings.Count(body, "\n") == 5, "4 SKU + judul: %s", body)
	})
	t.Run("impor: pratinjau dulu (tidak mengubah apa pun), error per baris", func(t *testing.T) {
		csv := "\xEF\xBB\xBFsku,product,price,compare_at_price,stock\n" +
			red + ",x,4.50,8,12\n" +
			"csv-00001-blue,x,6.00,,4\n" + // sama persis (huruf kecil tetap cocok) → tidak berubah
			"NOPE-1,x,1,,1\n" +
			red + ",dup,1,,1\n" +
			"'-HACK,x,abc,,1\n" +
			liveSKU + ",x,10,,3\n" +
			"\n"
		r := k.importCSV(t, csv, false).want(200)
		r.eq("total", "6").eq("changed", "2").eq("unchanged", "1").eq("errors", "3").eq("applied", "false").eq("review", "1")
		r.eq("rows.0.error.code", "unknown_sku").eq("rows.0.line", "4")
		r.eq("rows.1.error.code", "duplicate_sku").eq("rows.2.error.code", "bad_price")
		r.eq("rows.3.sku", red).eq("rows.3.old.price", "5.00").eq("rows.3.new.price", "4.50").eq("rows.3.new.compare_at_price", "8.00").eq("rows.3.new.stock", "12")
		check(t, e.scalar(t, `SELECT price || '/' || stock_qty FROM product_variants WHERE sku = $1`, red) == "5.00/10", "pratinjau tidak mengubah data")
		check(t, e.scalar(t, `SELECT moderation_status::text FROM products WHERE id = $1`, live.str("id")) == "approved", "pratinjau tidak mengirim ke review")

		r = k.importCSV(t, csv, true).want(200)
		r.eq("applied", "true").eq("changed", "2").eq("review", "1")
		check(t, e.scalar(t, `SELECT price || '/' || compare_at_price || '/' || stock_qty FROM product_variants WHERE sku = $1`, red) == "4.50/8.00/12", "diterapkan")
		check(t, e.scalar(t, `SELECT moderation_status::text FROM products WHERE id = $1`, live.str("id")) == "pending_review", "harga turun 75%% → review")
		check(t, e.scalar(t, `SELECT count(*) FROM inventory_movements m JOIN product_variants v ON v.id = m.variant_id WHERE v.sku = $1 AND m.reason = 'adjustment'`, red) == "1", "pergerakan stok")
	})
	t.Run("impor: titik koma & koma desimal, kolom harga coret kosong = hapus, stok < dipesan ditolak", func(t *testing.T) {
		e.sql(t, `UPDATE product_variants SET reserved_qty = 2 WHERE sku = 'CSV-00001-BLUE'`)
		r := k.importCSV(t, "SKU;Price;Compare;Qty\n"+red+";4,25;;9\nCSV-00001-BLUE;6;;1\n", true).want(200)
		r.eq("changed", "1").eq("errors", "1").eq("rows.0.error.code", "below_reserved")
		check(t, e.scalar(t, `SELECT price || '/' || coalesce(compare_at_price::text, '-') || '/' || stock_qty FROM product_variants WHERE sku = $1`, red) == "4.25/-/9", "titik koma diterapkan")
		r = k.importCSV(t, "sku,price,compare_at_price\n"+red+",5,4\n", false).want(200)
		r.eq("errors", "1").eq("rows.0.error.code", "bad_compare")
	})
	t.Run("impor: file tidak sah & batas ukuran", func(t *testing.T) {
		k.importCSV(t, "name,price\nx,1\n", false).wantErr(422, "import_no_sku_column")
		k.importCSV(t, "sku,product\nx,y\n", false).wantErr(422, "import_no_value_columns")
		k.importCSV(t, "sku,price\n", false).wantErr(422, "import_empty")
		big := "sku,price\n" + strings.Repeat("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA,1.00\n", seller.MaxImportBytes/38+10)
		r := k.importCSV(t, big, false)
		check(t, r.Status == 413, "file terlalu besar: %d %s", r.Status, truncate(string(r.Raw)))
	})
	t.Run("staf boleh ekspor/impor (aturan produk), orang lain tidak", func(t *testing.T) {
		stranger, _ := e.login(t)
		e.call(t, "GET", k.base+"/products/export.csv", nil, stranger).want(404)
		staff, phone := e.login(t)
		e.sql(t, `INSERT INTO store_members (store_id, user_id, role) SELECT $1, id, 'staff' FROM users WHERE phone = '+855' || substr($2, 2)`, k.sid, phone)
		e.call(t, "GET", k.base+"/products/export.csv", nil, staff).want(200)
		e.call(t, "PATCH", k.base, map[string]any{"sku_prefix": "ABC"}, staff).wantErr(403, "owner_only")
	})
}

func TestCatalogLowStock(t *testing.T) {
	e := getEnv(t)
	k := newCatalogKit(t, e, "Low Stock Shop")
	p := k.live(t, "Low stock tee", []string{"5", "5"}, []int{10, 10})
	v0, v1 := p.str("variants.0.id"), p.str("variants.1.id")
	count := func() string {
		return e.scalar(t, `SELECT count(*) FROM notifications WHERE store_id = $1 AND kind = 'stock.low'`, k.sid)
	}
	patch := func(v string, body map[string]any) {
		e.call(t, "PATCH", k.base+"/variants/"+v, body, k.tok).want(200)
	}

	t.Run("turun ke ambang toko (3) → satu notifikasi, sekali per penurunan", func(t *testing.T) {
		patch(v0, map[string]any{"stock_qty": 4})
		check(t, count() == "0", "di atas ambang: %s", count())
		patch(v0, map[string]any{"stock_qty": 3})
		check(t, count() == "1", "sampai ambang: %s", count())
		patch(v0, map[string]any{"stock_qty": 1})
		check(t, count() == "1", "tidak dikirim lagi: %s", count())
		patch(v0, map[string]any{"stock_qty": 10})
		patch(v0, map[string]any{"stock_qty": 2})
		check(t, count() == "2", "naik lalu turun lagi → dikirim lagi: %s", count())
		n := e.call(t, "GET", "/v1/seller/notifications", nil, k.tok, "X-Lang", "km").want(200)
		n.eq("notifications.0.kind", "stock.low")
		check(t, strings.Contains(n.str("notifications.0.title"), "ស្តុក") && strings.Contains(n.str("notifications.0.body"), "Low stock tee"), "teks Khmer: %s / %s",
			n.str("notifications.0.title"), n.str("notifications.0.body"))
		check(t, strings.Contains(n.str("notifications.0.body"), "2"), "jumlah tersisa: %s", n.str("notifications.0.body"))
	})
	t.Run("pesanan (stok dipesan) juga memicu", func(t *testing.T) {
		e.sql(t, `UPDATE product_variants SET reserved_qty = 8 WHERE id = $1`, v1) // tersedia 2
		check(t, count() == "3", "dari pesanan: %s", count())
		e.sql(t, `UPDATE product_variants SET reserved_qty = 0 WHERE id = $1`, v1)
	})
	t.Run("ambang per varian & ambang toko bisa diubah", func(t *testing.T) {
		patch(v1, map[string]any{"low_stock_threshold": 0})
		patch(v1, map[string]any{"stock_qty": 1})
		check(t, count() == "3", "ambang varian 0 → belum: %s", count())
		patch(v1, map[string]any{"stock_qty": 0})
		check(t, count() == "4", "habis → dikirim: %s", count())
		patch(v1, map[string]any{"stock_qty": 20, "low_stock_threshold": nil})
		e.call(t, "PATCH", k.base, map[string]any{"low_stock_threshold": 15}, k.tok).want(200)
		patch(v1, map[string]any{"stock_qty": 15})
		check(t, count() == "5", "ambang toko 15: %s", count())
	})
	t.Run("draf / varian nonaktif tidak memicu", func(t *testing.T) {
		d := k.create(t, k.product("Draft low", "womens-tops", nil, []string{"4"}, []int{10}, nil, false))
		patch(d.str("variants.0.id"), map[string]any{"stock_qty": 1})
		check(t, count() == "5", "draf: %s", count())
	})
	t.Run("teks en/km/zh, grup store, bunyi ding", func(t *testing.T) {
		for _, l := range []string{"en", "km", "zh"} {
			title, body := notify.Render(notify.StockLow, l, notify.Params{Product: "Tee", Note: "SKV-00001", Order: "2"})
			check(t, title != notify.StockLow && strings.Contains(body, "SKV-00001") && strings.Contains(body, "2"), "%s: %s / %s", l, title, body)
		}
		check(t, notify.GroupOf(notify.StockLow) == "store", "grup %s", notify.GroupOf(notify.StockLow))
		s, _ := notify.Sound(notify.StockLow)
		check(t, s == "ding", "bunyi %s", s)
	})
}
