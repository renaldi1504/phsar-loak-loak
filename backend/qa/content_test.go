package qa

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"
	"time"

	"phsar/internal/content"
)

func (e *env) adminContentUpload(t *testing.T, token string, data []byte) resp {
	t.Helper()
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	fw, _ := mw.CreateFormFile("file", "banner.png")
	_, _ = fw.Write(data)
	_ = mw.Close()
	req, _ := http.NewRequest(http.MethodPost, e.srv.URL+"/admin/api/content/upload", &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	return e.send(t, req)
}

// ---------------------------------------------------------------------------
// 18. Konten & SEO: banner beranda, pita pengumuman, koleksi pilihan, meta SEO, gambar OG
// ---------------------------------------------------------------------------

func TestContentAndSEO(t *testing.T) {
	e := getEnv(t)
	old := content.CacheTTL
	content.CacheTTL = 0
	t.Cleanup(func() {
		content.CacheTTL = old
		e.sql(t, `DELETE FROM content_banners`)
		e.sql(t, `DELETE FROM content_collections`)
		e.sql(t, `DELETE FROM seo_meta`)
		e.sql(t, `DELETE FROM platform_settings WHERE key IN ('content.announcement', 'content.og')`)
	})
	sa := newAdmin(t, e, "super_admin")
	mod := newAdmin(t, e, "moderator")
	s := newSellable(t, e, mod, "Content silk scarf", []string{"12.00"}, []int{5})
	home := func(lang string) string { return string(e.page(t, "/"+lang+"/").Raw) }

	e.call(t, "GET", "/admin/api/content", nil, mod).want(403)
	e.call(t, "GET", "/admin/api/content", nil, sa).want(200).eq("banners", "[]")

	var img string
	t.Run("banner beranda: validasi, jadwal, tampil", func(t *testing.T) {
		e.adminContentUpload(t, mod, pngBytes).want(403)
		e.adminContentUpload(t, sa, []byte("not an image at all, just text")).want(415)
		up := e.adminContentUpload(t, sa, pngBytes).want(201)
		img = up.str("key")
		check(t, strings.Contains(up.str("url"), "/public/content/"), "gambar banner harus publik: %s", up.str("url"))
		e.call(t, "POST", "/admin/api/content/banners", map[string]any{"title": map[string]string{"en": "x"}}, sa).wantErr(422, "validation", "image_key")
		e.call(t, "POST", "/admin/api/content/banners", map[string]any{"image_key": img, "link": "javascript:alert(1)"}, sa).wantErr(422, "validation", "link")
		e.call(t, "POST", "/admin/api/content/banners", map[string]any{"image_key": img, "link": "//evil.example"}, sa).wantErr(422, "validation", "link")
		e.call(t, "POST", "/admin/api/content/banners", map[string]any{"image_key": img, "title": map[string]string{"xx": "bad"}}, sa).wantErr(422, "validation", "title")
		e.call(t, "POST", "/admin/api/content/banners", map[string]any{"image_key": img, "title": map[string]string{"en": strings.Repeat("a", 81)}}, sa).wantErr(422, "validation", "title_en")
		r := e.call(t, "POST", "/admin/api/content/banners", map[string]any{"image_key": img, "link": "/en/c/fashion", "active": true,
			"title": map[string]string{"en": "Big Rainy Season Sale", "km": "លក់ធំរដូវវស្សា"}, "subtitle": map[string]string{"en": "Up to 30% off"}}, sa).want(200)
		check(t, r.len("banners") == 1, "banner harus tersimpan")
		tomorrow := time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)
		e.call(t, "POST", "/admin/api/content/banners", map[string]any{"image_key": img, "active": true, "starts_at": tomorrow,
			"title": map[string]string{"en": "Later promo"}}, sa).want(200)
		e.call(t, "POST", "/admin/api/content/banners", map[string]any{"image_key": img, "active": false,
			"title": map[string]string{"en": "Switched off promo"}}, sa).want(200)
		h := home("en")
		check(t, strings.Contains(h, "Big Rainy Season Sale") && strings.Contains(h, "Up to 30% off") && strings.Contains(h, `href="/en/c/fashion"`), "banner harus tampil di beranda")
		check(t, !strings.Contains(h, "Later promo") && !strings.Contains(h, "Switched off promo"), "banner terjadwal / mati tidak boleh tampil")
		check(t, strings.Contains(home("km"), "លក់ធំរដូវវស្សា"), "judul Khmer harus tampil di beranda Khmer")
		check(t, strings.Contains(home("zh"), "Big Rainy Season Sale"), "bahasa tanpa teks memakai bahasa Inggris")
		e.call(t, "PUT", "/admin/api/content/banners/not-a-uuid", map[string]any{"image_key": img}, sa).want(404)
	})

	t.Run("pita pengumuman di semua halaman", func(t *testing.T) {
		e.call(t, "PUT", "/admin/api/content/announcement", map[string]any{"active": true, "color": "purple", "text": map[string]string{"en": "Hi"}}, sa).wantErr(422, "validation", "color")
		e.call(t, "PUT", "/admin/api/content/announcement", map[string]any{"active": true, "color": "navy"}, sa).wantErr(422, "validation", "text_en")
		e.call(t, "PUT", "/admin/api/content/announcement", map[string]any{"active": true, "color": "navy", "link": "/en/legal/returns",
			"text": map[string]string{"en": "Free delivery all week <b>now</b>", "km": "ដឹកជញ្ជូនឥតគិតថ្លៃពេញមួយសប្តាហ៍"}}, sa).want(200)
		for _, p := range []string{"/en/", "/en/search?q=scarf", "/en/legal/terms"} {
			pg := string(e.page(t, p).Raw)
			check(t, strings.Contains(pg, "announce-navy") && strings.Contains(pg, "Free delivery all week &lt;b&gt;now&lt;/b&gt;"), "pengumuman (aman) harus ada di %s", p)
		}
		check(t, strings.Contains(home("km"), "ដឹកជញ្ជូនឥតគិតថ្លៃពេញមួយសប្តាហ៍"), "pengumuman Khmer")
		past := time.Now().Add(-time.Hour).UTC().Format(time.RFC3339)
		e.call(t, "PUT", "/admin/api/content/announcement", map[string]any{"active": true, "color": "navy", "ends_at": past,
			"text": map[string]string{"en": "Free delivery all week"}}, sa).want(200)
		check(t, !strings.Contains(home("en"), "announce-navy"), "pengumuman yang sudah berakhir tidak tampil")
	})

	t.Run("koleksi pilihan", func(t *testing.T) {
		e.call(t, "POST", "/admin/api/content/collections", map[string]any{"source": "category", "category_id": "fashion"}, sa).wantErr(422, "validation", "title_en")
		e.call(t, "POST", "/admin/api/content/collections", map[string]any{"title": map[string]string{"en": "X"}, "source": "category", "category_id": "nope"}, sa).wantErr(422, "validation", "category_id")
		e.call(t, "POST", "/admin/api/content/collections", map[string]any{"title": map[string]string{"en": "X"}, "source": "products", "product_ids": []string{"00000000-0000-0000-0000-000000000000"}}, sa).wantErr(422, "validation", "product_ids")
		e.call(t, "POST", "/admin/api/content/collections", map[string]any{"title": map[string]string{"en": "X"}, "source": "products", "product_ids": []string{s.productID}, "max_items": 50}, sa).wantErr(422, "validation", "max_items")
		e.call(t, "POST", "/admin/api/content/collections", map[string]any{"title": map[string]string{"en": "Hand picked for you"}, "source": "products",
			"product_ids": []string{s.productID, s.productID}, "active": true}, sa).want(200).eq("collections.0.product_ids.0", s.productID).eq("products.0.id", s.productID)
		e.call(t, "POST", "/admin/api/content/collections", map[string]any{"title": map[string]string{"en": "Fashion corner"}, "source": "category",
			"category_id": "fashion", "max_items": 4, "active": true}, sa).want(200)
		h := home("en")
		check(t, strings.Contains(h, "Hand picked for you") && strings.Contains(h, "Content silk scarf"), "koleksi pilihan harus tampil")
		check(t, strings.Contains(h, "Fashion corner") && strings.Contains(h, `href="/en/c/fashion"`), "koleksi kategori + lihat semua")
		search := e.call(t, "GET", "/admin/api/content/products?q=silk%20scarf", nil, sa).want(200)
		check(t, search.len("products") >= 1, "pencarian produk untuk koleksi")
	})

	t.Run("SEO & gambar OG", func(t *testing.T) {
		e.call(t, "PUT", "/admin/api/content/seo/nope", map[string]any{"title": map[string]string{"en": "x"}}, sa).wantErr(422, "validation", "page")
		e.call(t, "PUT", "/admin/api/content/seo/c:nope", map[string]any{"title": map[string]string{"en": "x"}}, sa).wantErr(422, "validation", "page")
		e.call(t, "PUT", "/admin/api/content/seo/home", map[string]any{"title": map[string]string{"en": "Phsar — video shopping in Sihanoukville"},
			"description": map[string]string{"en": "Watch real videos and buy from local shops."}}, sa).want(200)
		e.call(t, "PUT", "/admin/api/content/seo/c:fashion", map[string]any{"title": map[string]string{"en": "Clothes & shoes in Sihanoukville"}}, sa).want(200)
		h := home("en")
		check(t, strings.Contains(h, "<title>Phsar — video shopping in Sihanoukville</title>") && strings.Contains(h, `content="Watch real videos and buy from local shops."`), "SEO beranda")
		check(t, !strings.Contains(home("km"), "video shopping in Sihanoukville"), "SEO per bahasa: Khmer tetap teks bawaan")
		check(t, strings.Contains(string(e.page(t, "/en/c/womens-fashion").Raw), "<title>Clothes &amp; shoes in Sihanoukville</title>"), "SEO kategori")
		e.call(t, "PUT", "/admin/api/content/seo/c:fashion", map[string]any{"title": map[string]string{}}, sa).want(200)
		check(t, e.scalar(t, `SELECT count(*) FROM seo_meta`) == "1", "SEO kosong = kembali ke bawaan")
		e.call(t, "PUT", "/admin/api/content/og", map[string]any{"image_key": "public/products/nope.jpg"}, sa).wantErr(422, "validation", "image_key")
		e.call(t, "PUT", "/admin/api/content/og", map[string]any{"image_key": img}, sa).want(200)
		check(t, strings.Contains(home("en"), `og:image" content="`) && strings.Contains(home("en"), img), "gambar OG bawaan")
	})

	t.Run("hapus banner", func(t *testing.T) {
		v := e.call(t, "GET", "/admin/api/content", nil, sa).want(200)
		id := v.str("banners.0.id")
		e.call(t, "DELETE", "/admin/api/content/banners/"+id, nil, sa).want(200)
		e.call(t, "DELETE", "/admin/api/content/banners/"+id, nil, sa).want(404)
		check(t, !strings.Contains(home("en"), "Big Rainy Season Sale"), "banner terhapus tidak tampil")
	})
}
