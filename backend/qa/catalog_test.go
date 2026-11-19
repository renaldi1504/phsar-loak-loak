package qa

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"
)

func (e *env) catalogLogo(t *testing.T, token, brandID string, data []byte) resp {
	t.Helper()
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	fw, _ := mw.CreateFormFile("file", "logo.bin")
	_, _ = fw.Write(data)
	_ = mw.Close()
	req, _ := http.NewRequest(http.MethodPost, e.srv.URL+"/admin/api/catalog/brands/"+brandID+"/logo", &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	return e.send(t, req)
}

// ---------------------------------------------------------------------------
// 18. Kategori & merek dikelola admin (Categories & brands)
// ---------------------------------------------------------------------------

func TestCatalogAdmin(t *testing.T) {
	e := getEnv(t)
	mod := newAdmin(t, e, "moderator") // moderator punya izin catalog.manage
	cs := newAdmin(t, e, "cs")
	names := func(en string) map[string]string {
		return map[string]string{"en": en, "km": en + " km", "zh": en + " zh"}
	}
	cat := func(id string, extra map[string]any) map[string]any {
		m := map[string]any{"id": id, "icon": "box", "allows_used": true, "names": names("Pet supplies"), "slugs": map[string]string{"en": id}}
		for k, v := range extra {
			m[k] = v
		}
		return m
	}
	t.Cleanup(func() {
		e.sql(t, `DELETE FROM category_attributes WHERE category_id IN ('qa-pet-supplies','qa-dog-food')`)
		e.sql(t, `DELETE FROM category_translations WHERE category_id IN ('qa-dog-food','qa-pet-supplies')`)
		e.sql(t, `DELETE FROM categories WHERE id IN ('qa-dog-food','qa-pet-supplies')`)
		e.sql(t, `UPDATE products SET brand_id = NULL WHERE brand_id IN (SELECT id FROM brands WHERE slug LIKE 'qa-%')`)
		e.sql(t, `DELETE FROM brands WHERE slug LIKE 'qa-%'`)
		e.sql(t, `DELETE FROM prohibited_keywords WHERE keyword = 'qa-fake-brand'`)
		e.sql(t, `UPDATE categories SET deleted_at = NULL WHERE id IN ('electronics','beauty','fashion','accessories','kitchen')`)
	})

	t.Run("izin", func(t *testing.T) {
		e.call(t, "GET", "/admin/api/catalog", nil, cs).want(403)
		r := e.call(t, "GET", "/admin/api/catalog", nil, mod).want(200)
		check(t, r.len("categories") >= 5 && r.len("icons") >= 5, "katalog harus berisi kategori & ikon")
	})

	t.Run("buat, ubah, urutkan kategori", func(t *testing.T) {
		e.call(t, "POST", "/admin/api/catalog/categories", cat("Pet Supplies!", nil), mod).wantErr(422, "validation", "id")
		e.call(t, "POST", "/admin/api/catalog/categories", cat("qa-pet-supplies", map[string]any{"names": map[string]string{"en": "Pets"}}), mod).wantErr(422, "validation", "name_km")
		e.call(t, "POST", "/admin/api/catalog/categories", cat("qa-pet-supplies", map[string]any{"icon": "rocket"}), mod).wantErr(422, "validation", "icon")
		e.call(t, "POST", "/admin/api/catalog/categories", cat("qa-pet-supplies", map[string]any{"slugs": map[string]string{"en": "womens-fashion"}}), mod).wantErr(422, "validation", "slug_en")
		e.call(t, "POST", "/admin/api/catalog/categories", cat("fashion", nil), mod).wantErr(422, "validation", "id")
		e.call(t, "POST", "/admin/api/catalog/categories", cat("qa-pet-supplies", nil), mod).want(200)
		check(t, e.scalar(t, `SELECT slug || ' ' || path FROM category_translations WHERE category_id = 'qa-pet-supplies' AND lang = 'km'`) == "qa-pet-supplies qa-pet-supplies",
			"slug bahasa lain mengikuti bahasa Inggris")
		check(t, strings.Contains(string(e.call(t, "GET", "/v1/app/categories", nil, "", "X-Lang", "en").want(200).Raw), `"qa-pet-supplies"`), "kategori baru harus tampil di app")

		// sub-kategori satu tingkat saja
		e.call(t, "POST", "/admin/api/catalog/categories", cat("qa-dog-food", map[string]any{"parent_id": "qa-pet-supplies", "names": names("Dog food")}), mod).want(200)
		check(t, e.scalar(t, `SELECT path FROM category_translations WHERE category_id = 'qa-dog-food' AND lang = 'en'`) == "qa-pet-supplies/qa-dog-food", "path anak")
		e.call(t, "POST", "/admin/api/catalog/categories", cat("qa-cat-food", map[string]any{"parent_id": "qa-dog-food", "names": names("Cat food")}), mod).wantErr(422, "validation", "parent_id")

		// ubah nama, slug & SEO → path anak ikut berubah
		upd := cat("qa-pet-supplies", map[string]any{"names": names("Pets"), "slugs": map[string]string{"en": "qa-pets"}, "allows_used": false,
			"meta_titles": map[string]string{"en": "Pet shop Sihanoukville"}, "meta_descriptions": map[string]string{"en": strings.Repeat("x", 161)}})
		e.call(t, "PUT", "/admin/api/catalog/categories/qa-pet-supplies", upd, mod).wantErr(422, "validation", "meta_description_en")
		upd["meta_descriptions"] = map[string]string{"en": "Food and toys for pets"}
		e.call(t, "PUT", "/admin/api/catalog/categories/qa-pet-supplies", upd, mod).want(200)
		check(t, e.scalar(t, `SELECT path FROM category_translations WHERE category_id = 'qa-dog-food' AND lang = 'en'`) == "qa-pets/qa-dog-food", "path anak harus ikut slug induk")
		check(t, e.scalar(t, `SELECT allows_used::text || meta_title FROM categories c JOIN category_translations t ON t.category_id = c.id AND t.lang = 'en' WHERE c.id = 'qa-pet-supplies'`) == "falsePet shop Sihanoukville", "pengaturan tersimpan")
		check(t, e.call(t, "GET", "/en/c/qa-pets", nil, "").Status == 200, "halaman kategori baru harus bisa dibuka")
		e.call(t, "PUT", "/admin/api/catalog/categories/nope", upd, mod).want(404)

		// urutan: pindah ke atas = bertukar dengan kategori di atasnya
		before := e.scalar(t, `SELECT id FROM categories WHERE parent_id IS NULL ORDER BY sort_order, id OFFSET (SELECT count(*) - 2 FROM categories WHERE parent_id IS NULL) LIMIT 1`)
		e.call(t, "POST", "/admin/api/catalog/categories/qa-pet-supplies/move", map[string]string{"dir": "sideways"}, mod).wantErr(422, "validation", "dir")
		e.call(t, "POST", "/admin/api/catalog/categories/qa-pet-supplies/move", map[string]string{"dir": "up"}, mod).want(200)
		last := e.scalar(t, `SELECT id FROM categories WHERE parent_id IS NULL ORDER BY sort_order DESC, id DESC LIMIT 1`)
		check(t, last == before, "setelah naik, kategori terakhir harus %s, dapat %s", before, last)
	})

	t.Run("sembunyikan & hapus", func(t *testing.T) {
		e.call(t, "POST", "/admin/api/catalog/categories/qa-pet-supplies/hide", map[string]bool{"hidden": true}, mod).want(200)
		check(t, !strings.Contains(string(e.call(t, "GET", "/v1/app/categories", nil, "").Raw), `"qa-pet-supplies"`), "kategori tersembunyi tidak boleh tampil")
		check(t, e.scalar(t, `SELECT count(*) FROM categories WHERE id IN ('qa-pet-supplies','qa-dog-food') AND deleted_at IS NOT NULL`) == "2", "sub-kategori ikut tersembunyi")
		e.call(t, "POST", "/admin/api/catalog/categories/qa-pet-supplies/hide", map[string]bool{"hidden": false}, mod).want(200)
		check(t, strings.Contains(string(e.call(t, "GET", "/v1/app/categories", nil, "").Raw), `"qa-pet-supplies"`), "kategori tampil lagi")

		newSellable(t, e, mod, "Catalog shirt", []string{"5.00"}, []int{3}) // produk di sub-kategori womens-tops
		e.call(t, "DELETE", "/admin/api/catalog/categories/fashion", nil, mod).wantErr(409, "in_use")
		e.call(t, "DELETE", "/admin/api/catalog/categories/qa-pet-supplies", nil, mod).wantErr(409, "in_use") // masih punya sub-kategori
		e.call(t, "DELETE", "/admin/api/catalog/categories/qa-dog-food", nil, mod).want(200)
		e.call(t, "DELETE", "/admin/api/catalog/categories/qa-dog-food", nil, mod).want(404)
	})

	t.Run("atribut kategori", func(t *testing.T) {
		base := "/admin/api/catalog/categories/qa-pet-supplies/attributes"
		attr := map[string]any{"code": "pet_size", "label": names("Pet size"), "data_type": "select", "choices": []string{"S"}, "is_required": true}
		e.call(t, "POST", base, attr, mod).wantErr(422, "validation", "choices")
		attr["choices"] = []string{"Small", "Medium", "Large", "Small"}
		e.call(t, "POST", base, map[string]any{"code": "Pet Size", "label": names("x"), "data_type": "text"}, mod).wantErr(422, "validation", "code")
		e.call(t, "POST", base, map[string]any{"code": "weight", "label": names("Weight"), "data_type": "date"}, mod).wantErr(422, "validation", "data_type")
		e.call(t, "POST", base, attr, mod).want(200)
		e.call(t, "POST", base, attr, mod).wantErr(422, "validation", "code")
		check(t, e.scalar(t, `SELECT choices::text FROM category_attributes WHERE category_id = 'qa-pet-supplies' AND code = 'pet_size'`) == `["Small", "Medium", "Large"]`, "pilihan dibersihkan (tanpa duplikat)")
		aid := e.scalar(t, `SELECT id::text FROM category_attributes WHERE category_id = 'qa-pet-supplies' AND code = 'pet_size'`)
		e.call(t, "PUT", "/admin/api/catalog/attributes/"+aid, map[string]any{"code": "pet_size", "label": names("Size"), "data_type": "text", "is_filterable": true}, mod).want(200)
		e.call(t, "DELETE", "/admin/api/catalog/attributes/"+aid, nil, mod).want(200)
		e.call(t, "DELETE", "/admin/api/catalog/attributes/"+aid, nil, mod).want(404)
	})

	t.Run("merek", func(t *testing.T) {
		e.call(t, "POST", "/admin/api/catalog/brands", map[string]any{"name": ""}, mod).wantErr(422, "validation", "name")
		e.call(t, "POST", "/admin/api/catalog/brands", map[string]any{"name": "QA Acme", "slug": "qa-acme", "category_ids": []string{"nope"}}, mod).wantErr(422, "validation", "category_ids")
		e.call(t, "POST", "/admin/api/catalog/brands", map[string]any{"name": "QA Acme", "slug": "qa-acme", "verified": true, "category_ids": []string{"electronics"}}, mod).want(200)
		e.call(t, "POST", "/admin/api/catalog/brands", map[string]any{"name": "QA Acme 2", "slug": "qa-acme"}, mod).wantErr(422, "validation", "slug")
		bid := e.scalar(t, `SELECT id::text FROM brands WHERE slug = 'qa-acme'`)
		check(t, e.scalar(t, `SELECT verified::text || active::text || (SELECT count(*) FROM brand_categories WHERE brand_id = b.id) FROM brands b WHERE id = $1`, bid) == "truetrue1", "merek tersimpan")

		e.catalogLogo(t, mod, bid, []byte("not an image, just text")).want(415)
		e.catalogLogo(t, cs, bid, pngBytes).want(403)
		e.catalogLogo(t, mod, "99999999", pngBytes).want(404)
		e.catalogLogo(t, mod, bid, pngBytes).want(200)
		check(t, strings.Contains(e.scalar(t, `SELECT coalesce(logo_url, '') FROM brands WHERE id = $1`, bid), "/public/brands/"), "logo tersimpan sebagai file publik")

		e.call(t, "PUT", "/admin/api/catalog/brands/"+bid, map[string]any{"name": "QA Acme", "slug": "qa-acme", "active": false}, mod).want(200)
		check(t, e.scalar(t, `SELECT active::text || (SELECT count(*) FROM brand_categories WHERE brand_id = b.id) FROM brands b WHERE id = $1`, bid) == "false0", "nonaktif & kategori dikosongkan")
		e.sql(t, `UPDATE products SET brand_id = $1 WHERE id = (SELECT id FROM products ORDER BY created_at DESC LIMIT 1)`, bid)
		e.call(t, "DELETE", "/admin/api/catalog/brands/"+bid, nil, mod).wantErr(409, "in_use")
		e.sql(t, `UPDATE products SET brand_id = NULL WHERE brand_id = $1`, bid)
		e.call(t, "DELETE", "/admin/api/catalog/brands/"+bid, nil, mod).want(200)
	})

	t.Run("kata terlarang", func(t *testing.T) {
		e.call(t, "POST", "/admin/api/catalog/keywords", map[string]string{"keyword": "x", "lang": "en", "category": "counterfeit", "action": "review"}, mod).wantErr(422, "validation", "keyword")
		e.call(t, "POST", "/admin/api/catalog/keywords", map[string]string{"keyword": "qa-fake-brand", "lang": "fr", "category": "counterfeit", "action": "review"}, mod).wantErr(422, "validation", "lang")
		e.call(t, "POST", "/admin/api/catalog/keywords", map[string]string{"keyword": "qa-fake-brand", "lang": "en", "category": "counterfeit", "action": "ban"}, mod).wantErr(422, "validation", "action")
		e.call(t, "POST", "/admin/api/catalog/keywords", map[string]string{"keyword": "QA-Fake-Brand", "lang": "en", "category": "counterfeit", "action": "block"}, mod).want(200)
		e.call(t, "POST", "/admin/api/catalog/keywords", map[string]string{"keyword": "qa-fake-brand", "lang": "en", "category": "counterfeit", "action": "review"}, mod).wantErr(422, "validation", "keyword")
		kid := e.scalar(t, `SELECT id::text FROM prohibited_keywords WHERE keyword = 'qa-fake-brand'`)
		e.call(t, "DELETE", "/admin/api/catalog/keywords/"+kid, nil, mod).want(200)
		e.call(t, "DELETE", "/admin/api/catalog/keywords/"+kid, nil, mod).want(404)
	})
	check(t, e.scalar(t, `SELECT count(*) FROM audit_logs WHERE entity_table IN ('categories','brands','prohibited_keywords','category_attributes') AND actor_type = 'admin'`) != "0", "perubahan katalog harus tercatat di audit")
}
