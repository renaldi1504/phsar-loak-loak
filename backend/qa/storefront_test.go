package qa

import (
	"net/http"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// 9. Website pembeli: hanya produk tayang yang terlihat, 3 bahasa, SEO, aman dari XSS
// ---------------------------------------------------------------------------

// page: GET tanpa mengikuti redirect (supaya 301/302 bisa dicek).
func (e *env) page(t *testing.T, path string, headers ...string) resp {
	t.Helper()
	req, _ := http.NewRequest(http.MethodGet, e.srv.URL+path, nil)
	for i := 0; i+1 < len(headers); i += 2 {
		req.Header.Set(headers[i], headers[i+1])
	}
	client := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("GET %s: %v", path, err)
	}
	defer res.Body.Close()
	var b strings.Builder
	buf := make([]byte, 32<<10)
	for {
		n, err := res.Body.Read(buf)
		b.Write(buf[:n])
		if err != nil {
			break
		}
	}
	return resp{t: t, Status: res.StatusCode, Raw: []byte(b.String()), Header: res.Header}
}

func (r resp) has(t *testing.T, sub string, msg string) resp {
	t.Helper()
	check(t, strings.Contains(string(r.Raw), sub), "%s: %q tidak ada", msg, sub)
	return r
}

func (r resp) hasNot(t *testing.T, sub string, msg string) resp {
	t.Helper()
	check(t, !strings.Contains(string(r.Raw), sub), "%s: %q seharusnya tidak ada", msg, sub)
	return r
}

func TestStorefront(t *testing.T) {
	e := getEnv(t)
	tok, sid := newActiveStore(t, e)
	base := "/v1/seller/stores/" + sid + "/products"
	storeSlug := e.scalar(t, `SELECT slug FROM stores WHERE id = $1`, sid)
	admin := newAdmin(t, e, "moderator")

	newProduct := func(lang, name string, submit, approve bool) string {
		t.Helper()
		v := e.upload(t, tok, "product_media", mp4Bytes).want(201)
		f := e.upload(t, tok, "product_media", jpegBytes).want(201)
		body := map[string]any{
			"category_id": "womens-tops", "source_lang": lang, "name": name, "description": name + " — description line.", "condition": "new",
			"options": []any{map[string]any{"kind": "size", "name": "Size", "values": []any{map[string]any{"value": "M"}, map[string]any{"value": "L"}}}},
			"variants": []any{
				map[string]any{"options": []string{"M"}, "price": "5.50", "compare_at_price": "8", "stock_qty": 3},
				map[string]any{"options": []string{"L"}, "price": "7", "stock_qty": 0},
			},
			"media":  []any{map[string]any{"key": v.str("key"), "kind": "video"}, map[string]any{"key": f.str("key"), "kind": "photo"}},
			"submit": submit,
		}
		id := e.call(t, "POST", base, body, tok).want(201).str("id")
		if approve {
			e.call(t, "POST", "/admin/api/products/"+id+"/decision", map[string]string{"decision": "approve"}, admin).want(200)
		}
		return id
	}
	khmer := newProduct("km", "អាវយឺតពណ៌ក្រហម Phsarqa", true, true)
	xss := newProduct("en", `Tee <script>alert(1)</script> xssqa`, true, true)
	draft := newProduct("en", "Draft hidden draftqa", false, false)
	pending := newProduct("en", "Pending hidden pendingqa", true, false)
	e.drainTranslations(t)
	e.drainMedia(t)

	slug := func(id, lang string) string {
		return e.scalar(t, `SELECT slug FROM product_translations WHERE product_id = $1 AND lang = $2`, id, lang)
	}

	t.Run("/ memilih bahasa: cookie, Accept-Language, default", func(t *testing.T) {
		r := e.page(t, "/", "Accept-Language", "zh-CN,zh;q=0.9")
		check(t, r.Status == 302 && r.Header.Get("Location") == "/zh/", "zh: %d %s", r.Status, r.Header.Get("Location"))
		r = e.page(t, "/", "Cookie", "phsar_lang=km", "Accept-Language", "en")
		check(t, r.Header.Get("Location") == "/km/", "cookie km: %s", r.Header.Get("Location"))
		r = e.page(t, "/", "Accept-Language", "fr")
		check(t, r.Header.Get("Location") == "/en/", "default: %s", r.Header.Get("Location"))
	})

	t.Run("beranda: hanya produk tayang, nama berbahaya di-escape", func(t *testing.T) {
		for _, l := range []string{"en", "km", "zh"} {
			r := e.page(t, "/"+l+"/")
			check(t, r.Status == 200, "/%s/ status %d", l, r.Status)
			r.has(t, "Phsarqa", "produk Khmer").has(t, "xssqa", "produk XSS").hasNot(t, "draftqa", "draf").hasNot(t, "pendingqa", "belum disetujui")
			r.hasNot(t, "<script>alert(1)", "XSS").has(t, "&lt;script&gt;", "escape")
			r.has(t, `hreflang="zh-Hans"`, "hreflang").has(t, `"@type":"WebSite"`, "JSON-LD")
		}
		e.page(t, "/km/").has(t, `lang="km"`, "html lang").has(t, "ស្វែងរក", "teks Khmer")
	})

	t.Run("detail produk: 3 bahasa, hreflang, JSON-LD, harga riel", func(t *testing.T) {
		en := slug(khmer, "en")
		r := e.page(t, "/en/p/"+en)
		check(t, r.Status == 200, "status %d: %s", r.Status, truncate(string(r.Raw)))
		r.has(t, `"@type":"Product"`, "JSON-LD Product").has(t, `"priceCurrency":"USD"`, "mata uang").has(t, `"lowPrice":"5.50"`, "harga").
			has(t, `"@type":"BreadcrumbList"`, "breadcrumb").has(t, `rel="canonical" href="http://qa.local/en/p/`+en+`"`, "canonical").
			has(t, `hreflang="km" href="http://qa.local/km/p/`+slug(khmer, "km")+`"`, "hreflang km").
			has(t, "22,000 ៛", "riel (5.50 × 4000)").has(t, "-31%", "diskon").has(t, "Show original", "teks asli").
			has(t, `"VideoObject"`, "video di JSON-LD").has(t, `id="variants"`, "data varian")
		check(t, strings.Contains(r.Header.Get("Cache-Control"), "max-age"), "cache-control: %s", r.Header.Get("Cache-Control"))
		e.page(t, "/zh/p/"+slug(khmer, "zh")).has(t, "[zh]", "terjemahan zh").has(t, `lang="zh-Hans"`, "html lang zh")
		x := e.page(t, "/en/p/"+slug(xss, "en"))
		x.hasNot(t, "<script>alert(1)", "XSS di halaman produk").hasNot(t, `"name":"Tee <script>`, "XSS di JSON-LD")
	})

	t.Run("slug bahasa lain → 301 ke URL resmi", func(t *testing.T) {
		r := e.page(t, "/km/p/"+slug(khmer, "en"))
		check(t, r.Status == 301 && r.Header.Get("Location") == "/km/p/"+slug(khmer, "km"), "redirect: %d %s", r.Status, r.Header.Get("Location"))
	})

	t.Run("draf, belum disetujui, tidak ada → 404 noindex", func(t *testing.T) {
		for _, p := range []string{"/en/p/" + slug(draft, "en"), "/en/p/" + slug(pending, "en"), "/en/p/tidak-ada", "/en/s/tidak-ada", "/en/c/tidak-ada"} {
			r := e.page(t, p)
			check(t, r.Status == 404, "%s: status %d", p, r.Status)
			r.has(t, `content="noindex"`, p+" noindex")
		}
	})

	t.Run("pencarian lintas bahasa", func(t *testing.T) {
		e.page(t, "/en/search?q=អាវយឺត").has(t, "Phsarqa", "cari teks Khmer dari halaman Inggris").has(t, `content="noindex, follow"`, "noindex")
		e.page(t, "/km/search?q=xssqa").has(t, "xssqa", "cari nama Inggris dari halaman Khmer")
		e.page(t, "/en/search?q=pendingqa").hasNot(t, "Pending hidden", "produk belum disetujui").has(t, "Nothing found", "kosong")
		e.page(t, "/en/search?q=%25").hasNot(t, "Phsarqa", "% harus literal")
		e.page(t, "/en/search?q=%3Cscript%3Ex").hasNot(t, "<script>x", "query di-escape").has(t, `value="&lt;script&gt;x"`, "query di kotak cari")
	})

	t.Run("kategori & toko", func(t *testing.T) {
		cat := e.scalar(t, `SELECT slug FROM category_translations WHERE category_id = 'fashion' AND lang = 'en'`)
		r := e.page(t, "/en/c/"+cat)
		check(t, r.Status == 200, "kategori status %d", r.Status)
		r.has(t, "Phsarqa", "produk di kategori").hasNot(t, "draftqa", "draf di kategori")
		e.page(t, "/en/c/"+cat+"?sort=price_asc").has(t, `content="noindex, follow"`, "urutan lain noindex")
		s := e.page(t, "/en/s/"+storeSlug)
		check(t, s.Status == 200, "toko status %d", s.Status)
		s.has(t, "Phsarqa", "produk di toko").has(t, `"@type":"Store"`, "JSON-LD toko")
	})

	t.Run("robots.txt, sitemap.xml, aset", func(t *testing.T) {
		e.page(t, "/robots.txt").has(t, "Sitemap: http://qa.local/sitemap.xml", "robots").has(t, "Disallow: /admin/", "robots admin")
		sm := e.page(t, "/sitemap.xml")
		sm.has(t, "/km/p/"+slug(khmer, "km"), "sitemap produk").has(t, "/en/s/"+storeSlug, "sitemap toko").
			hasNot(t, slug(draft, "en"), "draf di sitemap").hasNot(t, slug(pending, "en"), "pending di sitemap").
			has(t, `hreflang="zh-Hans"`, "hreflang di sitemap")
		for _, a := range []string{"/assets/site.css", "/assets/site.js", "/assets/og.png", "/assets/icon.svg", "/assets/shop.js", "/en/cart"} {
			check(t, e.page(t, a).Status == 200, "%s tidak 200", a)
		}
	})

	t.Run("halaman belanja: noindex, tidak di-cache", func(t *testing.T) {
		for _, p := range []string{"/en/cart", "/km/checkout", "/zh/orders", "/en/orders/PH000000-XXXXX"} {
			r := e.page(t, p)
			check(t, r.Status == 200, "%s status %d", p, r.Status)
			check(t, r.Header.Get("Cache-Control") == "no-store", "%s cache: %s", p, r.Header.Get("Cache-Control"))
			r.has(t, `content="noindex, nofollow"`, p+" noindex").has(t, `id="shop"`, p+" kerangka").has(t, "PHSAR_SHOP", p+" teks")
		}
		e.page(t, "/km/checkout").has(t, "/ui/maps.js", "peta di checkout").has(t, "ចូលគណនីដើម្បីបញ្ជាទិញ", "teks Khmer")
	})

	t.Run("produk diarsip & toko dibekukan → hilang dari website", func(t *testing.T) {
		e.call(t, "POST", base+"/"+xss+"/status", map[string]string{"status": "archived"}, tok).want(200)
		check(t, e.page(t, "/en/p/"+slug(xss, "en")).Status == 404, "produk arsip masih tampil")
		e.page(t, "/en/").hasNot(t, "xssqa", "produk arsip di beranda")
		e.sql(t, `UPDATE stores SET status = 'suspended' WHERE id = $1`, sid)
		check(t, e.page(t, "/en/p/"+slug(khmer, "en")).Status == 404, "produk toko beku masih tampil")
		check(t, e.page(t, "/en/s/"+storeSlug).Status == 404, "toko beku masih tampil")
		e.page(t, "/sitemap.xml").hasNot(t, slug(khmer, "km"), "toko beku di sitemap")
	})
}
