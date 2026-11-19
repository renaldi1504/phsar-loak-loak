package qa

import (
	"strings"
	"testing"
	"time"

	"phsar/internal/content"
	"phsar/internal/storefront"
)

// ---------------------------------------------------------------------------
// Tampilan kategori v2: kategori musiman, kotak kategori (sampul), subkategori, Tonton per kategori, menu desktop.
// ---------------------------------------------------------------------------

func TestCategoryTilesAndSeasons(t *testing.T) {
	e := getEnv(t)
	oldC, oldT := content.CacheTTL, storefront.TileTTL
	content.CacheTTL, storefront.TileTTL = 0, 0
	t.Cleanup(func() {
		content.CacheTTL, storefront.TileTTL = oldC, oldT
		e.sql(t, `DELETE FROM season_picks`)
		e.sql(t, `UPDATE categories SET cover_key = NULL`)
	})
	sa := newAdmin(t, e, "super_admin")
	mod := newAdmin(t, e, "moderator")
	s := newSellable(t, e, mod, "Tiles linen shirt", []string{"14.00"}, []int{5})
	page := func(path string) string { return string(e.page(t, path).Raw) }

	t.Run("hanya pengelola konten", func(t *testing.T) {
		e.call(t, "POST", "/admin/api/content/seasons", map[string]any{"title": map[string]string{"en": "x"}}, mod).want(403)
		e.call(t, "PUT", "/admin/api/content/categories/fashion/cover", map[string]any{"image_key": ""}, mod).want(403)
	})

	t.Run("kategori musiman: validasi, jadwal, tampil di beranda & menu", func(t *testing.T) {
		e.call(t, "POST", "/admin/api/content/seasons", map[string]any{"source": "category", "category_id": "fashion"}, sa).wantErr(422, "validation", "title_en")
		e.call(t, "POST", "/admin/api/content/seasons", map[string]any{"title": map[string]string{"en": "X"}, "source": "category", "category_id": "nope"}, sa).wantErr(422, "validation", "category_id")
		e.call(t, "POST", "/admin/api/content/seasons", map[string]any{"title": map[string]string{"en": "X"}, "source": "search", "query": "a"}, sa).wantErr(422, "validation", "query")
		e.call(t, "POST", "/admin/api/content/seasons", map[string]any{"title": map[string]string{"en": "X"}, "source": "category", "category_id": "fashion", "color": "purple"}, sa).wantErr(422, "validation", "color")
		end := time.Date(2099, 1, 2, 0, 0, 0, 0, time.FixedZone("ICT", 7*3600)).UTC().Format(time.RFC3339)
		r := e.call(t, "POST", "/admin/api/content/seasons", map[string]any{"title": map[string]string{"en": "New Year Party", "km": "ឆ្នាំថ្មីជប់លៀង"},
			"source": "category", "category_id": "fashion", "color": "pink", "ends_at": end, "active": true}, sa).want(200)
		check(t, r.len("seasons") == 1, "musim harus tersimpan")
		e.call(t, "POST", "/admin/api/content/seasons", map[string]any{"title": map[string]string{"en": "Linen picks"}, "source": "products",
			"product_ids": []string{s.productID}, "color": "blue", "active": true}, sa).want(200)
		tomorrow := time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)
		e.call(t, "POST", "/admin/api/content/seasons", map[string]any{"title": map[string]string{"en": "Later season"}, "source": "search", "query": "linen",
			"starts_at": tomorrow, "active": true}, sa).want(200)
		h := page("/en/")
		check(t, strings.Contains(h, "New Year Party") && strings.Contains(h, "Until 1 Jan") && strings.Contains(h, "Linen picks") && strings.Contains(h, "All year"), "musim harus tampil: beranda")
		check(t, !strings.Contains(h, "Later season"), "musim terjadwal tidak boleh tampil")
		check(t, strings.Contains(page("/km/"), "ឆ្នាំថ្មីជប់លៀង") && strings.Contains(page("/km/"), "ដល់ 1 មករា"), "musim dalam bahasa Khmer")
		check(t, strings.Contains(page("/zh/"), "至 1月1日"), "tanggal dalam format Mandarin")
		check(t, strings.Contains(page("/en/c/womens-fashion"), `class="cat-menu"`), "menu kategori desktop ada di setiap halaman")
		check(t, strings.Contains(h, `href="/en/c/womens-tops"`), "menu desktop: sub-kategori menuju halamannya sendiri")
	})

	t.Run("halaman musim", func(t *testing.T) {
		snap := e.scalar(t, `SELECT left(id::text, 8) FROM season_picks WHERE title->>'en' = 'Linen picks'`)
		p := page("/en/season/" + snap)
		check(t, strings.Contains(p, "Linen picks") && strings.Contains(p, "Tiles linen shirt"), "halaman musim harus berisi produk pilihan")
		later := e.scalar(t, `SELECT left(id::text, 8) FROM season_picks WHERE title->>'en' = 'Later season'`)
		e.page(t, "/en/season/"+later).want(404)
		e.page(t, "/en/season/zzzzzzzz").want(404)
	})

	t.Run("sampul kotak kategori", func(t *testing.T) {
		e.call(t, "PUT", "/admin/api/content/categories/fashion/cover", map[string]any{"image_key": "not-uploaded"}, sa).wantErr(422, "validation", "image_key")
		e.call(t, "PUT", "/admin/api/content/categories/nope/cover", map[string]any{"image_key": ""}, sa).want(404)
		up := e.adminContentUpload(t, sa, pngBytes).want(201)
		r := e.call(t, "PUT", "/admin/api/content/categories/fashion/cover", map[string]any{"image_key": up.str("key")}, sa).want(200)
		check(t, strings.Contains(string(r.Raw), up.str("key")), "sampul tersimpan: %s", r.Raw)
		h := page("/en/")
		check(t, strings.Contains(h, `class="tiles"`) && strings.Contains(h, up.str("url")), "kotak kategori memakai sampul pilihan admin")
		e.call(t, "PUT", "/admin/api/content/categories/fashion/cover", map[string]any{"image_key": ""}, sa).want(200)
		check(t, !strings.Contains(page("/en/"), up.str("url")), "kembali ke foto otomatis")
	})

	t.Run("subkategori: chip, halaman sendiri, produk induk termasuk anak", func(t *testing.T) {
		// Produk newSellable ada di sub-kategori womens-tops (induk: fashion / womens-fashion).
		p := page("/en/c/womens-fashion")
		check(t, strings.Contains(p, `href="/en/c/womens-tops"`) && strings.Contains(p, "Tiles linen shirt"), "induk menampilkan chip & produk subkategori")
		check(t, strings.Contains(p, "/en/c/womens-fashion/watch"), "tombol Tonton ada")
		sub := page("/en/c/womens-tops")
		check(t, strings.Contains(sub, "Tiles linen shirt") && strings.Contains(sub, `aria-current="page">Tops`), "halaman subkategori + chip aktif")
		check(t, strings.Contains(sub, "/en/c/womens-tops/watch") && strings.Contains(sub, `href="/en/c/womens-fashion">Women`), "Tonton subkategori & breadcrumb induk")
		r := e.page(t, "/en/c/womens-fashion?sub=womens-tops")
		check(t, r.Status == 301, "link lama ?sub= dialihkan (status %d)", r.Status)
		check(t, strings.Contains(page("/en/c/womens-fashion?sort=popular"), `value="popular" selected`), "urutan terlaris")
	})

	t.Run("Tonton per kategori", func(t *testing.T) {
		e.drainMedia(t)
		w := page("/en/c/womens-fashion/watch")
		check(t, strings.Contains(w, `class="watch"`) && strings.Contains(w, "noindex"), "halaman Tonton")
		if e.scalar(t, `SELECT count(*) FROM product_media WHERE product_id = $1 AND kind = 'video' AND video_url IS NOT NULL`, s.productID) != "0" {
			check(t, strings.Contains(w, "Tiles linen shirt") && strings.Contains(w, "data-watch"), "video produk tampil di Tonton")
			check(t, strings.Contains(page("/en/c/womens-tops/watch"), "Tiles linen shirt"), "Tonton per subkategori")
		}
		check(t, strings.Contains(page("/en/c/womens-tops/watch"), `aria-current="page">Tops`), "Tonton subkategori: chip saudara aktif")
		check(t, e.page(t, "/en/c/fashion/watch").Status == 301, "ID/slug lama dialihkan ke slug bahasa ini")
		e.page(t, "/en/c/nope/watch").want(404)
	})
}
