package qa

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"phsar/internal/media"
	"phsar/internal/notify"
	"phsar/internal/translate"
)

// ---------------------------------------------------------------------------
// 1. Login & sesi
// ---------------------------------------------------------------------------

func TestAuthOTPAndSessions(t *testing.T) {
	e := getEnv(t)
	phone := "0971234560"
	r := e.call(t, "POST", "/v1/app/auth/otp/request", map[string]string{"phone": phone}, "").want(200)
	norm := r.str("phone")
	check(t, norm == "+855971234560", "nomor tidak dinormalisasi ke +855: %s", norm)

	t.Run("minta kode lagi kurang dari 60 detik ditolak", func(t *testing.T) {
		e.call(t, "POST", "/v1/app/auth/otp/request", map[string]string{"phone": phone}, "").wantErr(429, "otp_too_soon")
	})
	t.Run("kode salah dijawab cepat (tidak macet) dan dihitung", func(t *testing.T) {
		start := time.Now()
		e.call(t, "POST", "/v1/app/auth/otp/verify", map[string]string{"phone": phone, "code": "000000"}, "").wantErr(401, "otp_invalid")
		check(t, time.Since(start) < 3*time.Second, "verifikasi kode salah butuh %s (macet?)", time.Since(start))
		check(t, e.scalar(t, `SELECT attempts FROM otp_codes WHERE phone = $1 ORDER BY created_at DESC LIMIT 1`, norm) == "1", "percobaan salah tidak tercatat")
	})
	t.Run("nomor tidak valid ditolak", func(t *testing.T) {
		e.call(t, "POST", "/v1/app/auth/otp/request", map[string]string{"phone": "12"}, "").wantErr(400, "phone_invalid")
	})

	ok := e.call(t, "POST", "/v1/app/auth/otp/verify", map[string]string{"phone": phone, "code": e.sms.code(norm)}, "").want(200)
	access, refresh := ok.str("access_token"), ok.str("refresh_token")
	check(t, access != "" && refresh != "", "token kosong")

	t.Run("kode yang sudah dipakai tidak bisa dipakai lagi", func(t *testing.T) {
		e.call(t, "POST", "/v1/app/auth/otp/verify", map[string]string{"phone": phone, "code": e.sms.code(norm)}, "").wantErr(401, "otp_invalid")
	})
	t.Run("/me dengan token", func(t *testing.T) {
		e.call(t, "GET", "/v1/app/me", nil, access).want(200)
		e.call(t, "GET", "/v1/app/me", nil, "").want(401)
		e.call(t, "GET", "/v1/app/me", nil, "bukan.token.valid").want(401)
	})
	t.Run("refresh berputar; refresh lama ditolak", func(t *testing.T) {
		r2 := e.call(t, "POST", "/v1/app/auth/refresh", map[string]string{"refresh_token": refresh}, "").want(200)
		check(t, r2.str("refresh_token") != refresh, "refresh token tidak berganti")
		e.call(t, "POST", "/v1/app/auth/refresh", map[string]string{"refresh_token": refresh}, "").want(401)
		refresh = r2.str("refresh_token")
		access = r2.str("access_token")
	})
	t.Run("logout mematikan sesi", func(t *testing.T) {
		e.call(t, "POST", "/v1/app/auth/logout", map[string]string{"refresh_token": refresh}, "").want(204)
		e.call(t, "GET", "/v1/app/me", nil, access).want(401)
		e.call(t, "POST", "/v1/app/auth/refresh", map[string]string{"refresh_token": refresh}, "").want(401)
	})
	t.Run("login web wajib header X-Phsar (CSRF)", func(t *testing.T) {
		e.call(t, "POST", "/seller/auth/otp/request", map[string]string{"phone": "0971234561"}, "").want(403)
		e.call(t, "POST", "/seller/auth/otp/request", map[string]string{"phone": "0971234561"}, "", "X-Phsar", "1").want(200)
	})
}

// ---------------------------------------------------------------------------
// 2. Onboarding toko → verifikasi admin → aktif
// ---------------------------------------------------------------------------

// pickup mengembalikan body titik penjemputan yang valid di Sihanoukville (Ochheuteal).
func pickup(label string) map[string]any {
	return map[string]any{
		"label": label, "address_line": "Ochheuteal Beach Road", "lat": 10.6105, "lng": 103.5316,
		"contact_name": "Dara", "contact_phone": "012345678", "pickup_hours": map[string][]string{"mon": {"09:00", "17:00"}},
		"note": "Sebelah apotek", "is_default": true, "google_place_id": "ChIJqa-test",
	}
}

// newActiveStore: toko lengkap yang sudah disetujui admin (dipakai tes produk).
func newActiveStore(t *testing.T, e *env) (token, storeID string) {
	t.Helper()
	token, _ = e.login(t)
	st := e.call(t, "POST", "/v1/seller/stores", map[string]any{"name": "Toko " + time.Now().Format("150405.000"), "about": "Toko baju", "about_lang": "en"}, token).want(201)
	storeID = st.str("id")
	e.sql(t, `UPDATE stores SET status = 'active' WHERE id = $1`, storeID)
	return
}

// newAdmin: user baru → peran admin → login di /admin (HP) → pasang MFA. Mengembalikan access token admin yang lolos MFA.
func newAdmin(t *testing.T, e *env, role string) string {
	t.Helper()
	tok, _ := newAdminFull(t, e, role)
	return tok
}

func TestStoreOnboardingAndVerification(t *testing.T) {
	e := getEnv(t)
	tok, _ := e.login(t)
	other, _ := e.login(t)

	st := e.call(t, "POST", "/v1/seller/stores", map[string]any{"name": "Grace Beauty", "about": "Kosmetik asli", "about_lang": "en"}, tok).want(201)
	id := st.str("id")
	st.eq("status", "onboarding").eq("checklist.can_submit", "false")

	t.Run("validasi toko", func(t *testing.T) {
		e.call(t, "POST", "/v1/seller/stores", map[string]any{"name": "A", "about": "", "about_lang": "en"}, tok).wantErr(422, "validation", "name")
		e.call(t, "POST", "/v1/seller/stores", map[string]any{"name": "Toko", "about": "", "about_lang": "en", "x": 1}, tok).wantErr(400, "bad_json")
	})
	t.Run("orang lain tidak bisa melihat toko ini (404, bukan 403)", func(t *testing.T) {
		e.call(t, "GET", "/v1/seller/stores/"+id, nil, other).wantErr(404, "store_not_found")
		e.call(t, "GET", "/v1/seller/stores/"+id+"/pickup-points", nil, other).want(404)
	})
	t.Run("deskripsi toko diterjemahkan", func(t *testing.T) {
		e.drainTranslations(t)
		r := e.call(t, "GET", "/v1/seller/stores/"+id, nil, tok).want(200)
		r.eq("about.km", "[km] Kosmetik asli").eq("about.zh", "[zh] Kosmetik asli")
	})

	t.Run("upload: jenis file dicek dari isi, bukan nama", func(t *testing.T) {
		e.upload(t, tok, "store_logo", htmlJPEG).wantErr(415, "upload_type")
		e.upload(t, tok, "unknown_kind", jpegBytes).want(415)
	})
	logo := e.upload(t, tok, "store_logo", pngBytes).want(201)
	e.call(t, "PATCH", "/v1/seller/stores/"+id, map[string]any{"logo_key": logo.str("key")}, tok).want(200)
	t.Run("logo publik bisa dibuka", func(t *testing.T) {
		u := strings.TrimPrefix(logo.str("url"), e.cfg.PublicBaseURL)
		e.call(t, "GET", u, nil, "").want(200)
	})
	t.Run("tidak bisa memakai file milik orang lain", func(t *testing.T) {
		otherStore := e.call(t, "POST", "/v1/seller/stores", map[string]any{"name": "Other shop", "about": "", "about_lang": "en"}, other).want(201)
		e.call(t, "PATCH", "/v1/seller/stores/"+otherStore.str("id"), map[string]any{"logo_key": logo.str("key")}, other).wantErr(422, "upload_not_found")
	})

	doc := e.upload(t, tok, "kyc_document", pdfBytes).want(201)
	check(t, doc.str("url") == "", "dokumen KYC tidak boleh punya URL publik")
	t.Run("dokumen KYC tidak bisa dibuka lewat /media", func(t *testing.T) {
		e.call(t, "GET", "/media/"+doc.str("key"), nil, "").want(404)
	})
	e.call(t, "PUT", "/v1/seller/stores/"+id+"/kyc", map[string]any{"seller_type": "individual_unregistered", "legal_name": "Sok Dara",
		"id_document_key": doc.str("key")}, tok).want(200).eq("kyc.verification_status", "pending")

	var pointID string
	t.Run("titik penjemputan: di luar Sihanoukville ditolak", func(t *testing.T) {
		p := pickup("Phnom Penh")
		p["lat"], p["lng"] = 11.5564, 104.9282
		e.call(t, "POST", "/v1/seller/stores/"+id+"/pickup-points", p, tok).wantErr(422, "outside_service_area")
	})
	t.Run("titik penjemputan: validasi field", func(t *testing.T) {
		p := pickup("")
		e.call(t, "POST", "/v1/seller/stores/"+id+"/pickup-points", p, tok).wantErr(422, "validation", "label")
		p = pickup("Shop")
		p["pickup_hours"] = map[string][]string{"mon": {"17:00", "09:00"}}
		e.call(t, "POST", "/v1/seller/stores/"+id+"/pickup-points", p, tok).wantErr(422, "validation", "pickup_hours")
	})
	t.Run("titik penjemputan: buat + Google place id tersimpan", func(t *testing.T) {
		r := e.call(t, "POST", "/v1/seller/stores/"+id+"/pickup-points", pickup("Shop Ochheuteal"), tok).want(201)
		pointID = r.str("pickup_points.0.id")
		r.eq("pickup_points.0.google_place_id", "ChIJqa-test").eq("pickup_points.0.verified", "false")
	})

	st = e.call(t, "GET", "/v1/seller/stores/"+id, nil, tok).want(200).eq("checklist.can_submit", "true")
	e.call(t, "POST", "/v1/seller/stores/"+id+"/submit", nil, tok).want(200).eq("status", "pending_review")
	t.Run("toko tidak bisa diubah selama direview", func(t *testing.T) {
		e.call(t, "PATCH", "/v1/seller/stores/"+id, map[string]any{"name": "Ganti"}, tok).wantErr(409, "store_under_review")
	})

	admin := newAdmin(t, e, "moderator")
	t.Run("bukan admin ditolak (token app tidak berlaku di API admin)", func(t *testing.T) {
		e.call(t, "GET", "/admin/api/stores", nil, tok).wantErr(401, "token_invalid")
	})
	t.Run("peran tanpa izin ditolak", func(t *testing.T) {
		fin := newAdmin(t, e, "cs")
		e.call(t, "GET", "/admin/api/products", nil, fin).wantErr(403, "forbidden")
	})
	q := e.call(t, "GET", "/admin/api/stores", nil, admin).want(200)
	check(t, strings.Contains(string(q.Raw), id), "toko tidak muncul di antrean admin")
	d := e.call(t, "GET", "/admin/api/stores/"+id, nil, admin).want(200)
	t.Run("admin bisa membuka dokumen KYC (dan tercatat di audit log)", func(t *testing.T) {
		e.call(t, "GET", d.str("id_document_url"), nil, admin).want(200)
		e.call(t, "GET", d.str("id_document_url"), nil, tok).want(401)
		check(t, e.scalar(t, `SELECT count(*) FROM audit_logs WHERE action = 'kyc_document.view'`) != "0", "akses KYC tidak tercatat")
	})
	t.Run("toko belum bisa disetujui sebelum KYC & titik diverifikasi", func(t *testing.T) {
		e.call(t, "POST", "/admin/api/stores/"+id+"/decision", map[string]string{"decision": "approve"}, admin).wantErr(422, "not_ready")
	})
	t.Run("tolak wajib catatan", func(t *testing.T) {
		e.call(t, "POST", "/admin/api/stores/"+id+"/kyc", map[string]string{"decision": "reject"}, admin).want(422)
	})
	e.call(t, "POST", "/admin/api/stores/"+id+"/kyc", map[string]string{"decision": "approve"}, admin).want(200)
	e.call(t, "POST", "/admin/api/pickup-points/"+pointID+"/decision", map[string]string{"decision": "approve"}, admin).want(200)
	e.call(t, "POST", "/admin/api/stores/"+id+"/decision", map[string]string{"decision": "approve"}, admin).want(200)
	st = e.call(t, "GET", "/v1/seller/stores/"+id, nil, tok).want(200)
	st.eq("status", "active").eq("checklist.can_sell", "true")

	t.Run("ubah titik penjemputan: label saja tetap terverifikasi", func(t *testing.T) {
		p := pickup("Shop Ochheuteal (baru)")
		r := e.call(t, "PUT", "/v1/seller/stores/"+id+"/pickup-points/"+pointID, p, tok).want(200)
		r.eq("pickup_points.0.label", "Shop Ochheuteal (baru)").eq("pickup_points.0.verified", "true")
	})
	t.Run("ubah titik penjemputan: lokasi pindah → harus dicek ulang", func(t *testing.T) {
		p := pickup("Shop Ochheuteal (baru)")
		p["lat"], p["lng"], p["address_line"], p["google_place_id"] = 10.6281, 103.5230, "Ekareach Street", "ChIJmoved"
		r := e.call(t, "PUT", "/v1/seller/stores/"+id+"/pickup-points/"+pointID, p, tok).want(200)
		r.eq("pickup_points.0.verified", "false").eq("pickup_points.0.google_place_id", "ChIJmoved").eq("pickup_points.0.lat", "10.6281")
	})
	t.Run("maksimal 5 titik", func(t *testing.T) {
		for i := 0; i < 4; i++ {
			e.call(t, "POST", "/v1/seller/stores/"+id+"/pickup-points", pickup("Extra"), tok).want(201)
		}
		e.call(t, "POST", "/v1/seller/stores/"+id+"/pickup-points", pickup("Extra"), tok).wantErr(422, "validation", "pickup_points")
	})
	t.Run("hapus titik", func(t *testing.T) {
		e.call(t, "DELETE", "/v1/seller/stores/"+id+"/pickup-points/"+pointID, nil, tok).want(200)
	})
	t.Run("perubahan tercatat di audit log dengan pelakunya", func(t *testing.T) {
		n := e.scalar(t, `SELECT count(*) FROM audit_logs WHERE entity_table = 'stores' AND actor_type = 'admin'`)
		check(t, n != "0", "keputusan admin tidak tercatat di audit_logs")
	})
}

// ---------------------------------------------------------------------------
// 3. Produk, varian, media, terjemahan, moderasi
// ---------------------------------------------------------------------------

func TestProducts(t *testing.T) {
	e := getEnv(t)
	tok, sid := newActiveStore(t, e)
	base := "/v1/seller/stores/" + sid + "/products"

	jpg := e.upload(t, tok, "product_media", jpegBytes).want(201)
	png := e.upload(t, tok, "product_media", pngBytes).want(201)
	mp4 := e.upload(t, tok, "product_media", mp4Bytes).want(201).eq("content_type", "video/mp4")
	mov := e.upload(t, tok, "product_media", movBytes).want(201).eq("content_type", "video/quicktime")
	e.upload(t, tok, "product_media", htmlJPEG).wantErr(415, "upload_type")

	product := func(over map[string]any) map[string]any {
		p := map[string]any{
			"category_id": "womens-tops", "source_lang": "en", "name": "Cotton T-shirt basic", "description": "Soft cotton tee.",
			"condition": "new",
			"options": []any{
				map[string]any{"kind": "color", "name": "Color", "values": []any{map[string]any{"value": "Red", "swatch_hex": "#ff0000"}, map[string]any{"value": "Black"}}},
				map[string]any{"kind": "size", "name": "Size", "values": []any{map[string]any{"value": "M"}, map[string]any{"value": "L"}}},
			},
			"variants": []any{
				map[string]any{"options": []string{"Red", "M"}, "price": "5.50", "stock_qty": 10},
				map[string]any{"options": []string{"Red", "L"}, "price": "5.50", "stock_qty": 4},
				map[string]any{"options": []string{"Black", "M"}, "price": "6", "compare_at_price": "8", "stock_qty": 0, "sku": "TEE-BLK-M"},
			},
			"media": []any{map[string]any{"key": mp4.str("key"), "kind": "video"}, map[string]any{"key": jpg.str("key"), "kind": "photo"}},
		}
		for k, v := range over {
			p[k] = v
		}
		return p
	}

	t.Run("validasi", func(t *testing.T) {
		cases := []struct {
			name, field string
			over        map[string]any
		}{
			{"nama terlalu pendek", "name", map[string]any{"name": "ab"}},
			{"bahasa tidak dikenal", "source_lang", map[string]any{"source_lang": "id"}},
			{"kategori kosong", "category_id", map[string]any{"category_id": ""}},
			{"kategori tidak dikenal", "category_id", map[string]any{"category_id": "weapons"}},
			{"barang bekas tanpa catatan", "condition_note", map[string]any{"condition": "used_good"}},
			{"kosmetik bekas", "condition", map[string]any{"category_id": "skincare", "condition": "used_good", "condition_note": "x"}},
			{"nilai varian tidak ada di opsi", "variants[0].options", map[string]any{"variants": []any{map[string]any{"options": []string{"Blue", "M"}, "price": "5", "stock_qty": 1}}}},
			{"harga nol", "variants[0].price", map[string]any{"variants": []any{map[string]any{"options": []string{"Red", "M"}, "price": "0", "stock_qty": 1}}}},
			{"harga coret lebih rendah", "variants[0].compare_at_price", map[string]any{"variants": []any{map[string]any{"options": []string{"Red", "M"}, "price": "5", "compare_at_price": "4", "stock_qty": 1}}}},
			{"stok negatif", "variants[0].stock_qty", map[string]any{"variants": []any{map[string]any{"options": []string{"Red", "M"}, "price": "5", "stock_qty": -1}}}},
			{"varian dobel", "variants[1].options", map[string]any{"variants": []any{map[string]any{"options": []string{"Red", "M"}, "price": "5", "stock_qty": 1}, map[string]any{"options": []string{"red", "m"}, "price": "5", "stock_qty": 1}}}},
			{"foto dikirim sebagai video", "media[0]", map[string]any{"media": []any{map[string]any{"key": jpg.str("key"), "kind": "video"}}}},
			{"dua video", "media", map[string]any{"media": []any{map[string]any{"key": mp4.str("key"), "kind": "video"}, map[string]any{"key": mov.str("key"), "kind": "video"}}}},
			{"kirim tanpa video", "media", map[string]any{"media": []any{map[string]any{"key": jpg.str("key"), "kind": "photo"}}, "submit": true}},
			{"kirim tanpa foto", "media", map[string]any{"media": []any{map[string]any{"key": mp4.str("key"), "kind": "video"}}, "submit": true}},
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				e.call(t, "POST", base, product(c.over), tok).wantErr(422, "validation", c.field)
			})
		}
	})

	p := e.call(t, "POST", base, product(nil), tok).want(201)
	pid := p.str("id")
	p.eq("status", "draft").eq("min_price", "5.50").eq("max_price", "6.00").eq("total_stock", "14").eq("translation_pending", "true")
	check(t, regexp.MustCompile(`^[A-Z0-9]{2,5}-[0-9]{5}-RED-M$`).MatchString(p.str("variants.0.sku")), "SKU otomatis tidak dibuat: %s", p.str("variants.0.sku"))

	t.Run("terjemahan: nama, deskripsi, opsi & nilai varian", func(t *testing.T) {
		e.drainTranslations(t)
		r := e.call(t, "GET", base+"/"+pid, nil, tok).want(200)
		r.eq("translations.km.name", "[km] Cotton T-shirt basic").eq("translations.zh.description", "[zh] Soft cotton tee.")
		r.eq("options.0.name.km", "[km] Color").eq("options.0.values.0.value.zh", "[zh] Red").eq("translation_pending", "false")
		check(t, strings.HasSuffix(r.str("translations.km.slug"), strings.ReplaceAll(pid, "-", "")[:8]), "slug tidak unik per produk")
	})

	t.Run("anggota toko lain tidak bisa melihat atau memakai upload", func(t *testing.T) {
		tok2, sid2 := newActiveStore(t, e)
		e.call(t, "GET", base+"/"+pid, nil, tok2).want(404)
		e.call(t, "GET", "/v1/seller/stores/"+sid2+"/products/"+pid, nil, tok2).want(404)
		e.call(t, "POST", "/v1/seller/stores/"+sid2+"/products", product(nil), tok2).wantErr(422, "upload_not_found")
	})

	t.Run("ubah varian: ID lama tetap, yang dihapus hilang, stok tercatat", func(t *testing.T) {
		ids := map[string]string{}
		for i := 0; i < p.len("variants"); i++ {
			ids[p.str("variants."+itoa(i)+".options.0")+"/"+p.str("variants."+itoa(i)+".options.1")] = p.str("variants." + itoa(i) + ".id")
		}
		r := e.call(t, "PUT", base+"/"+pid, product(map[string]any{"variants": []any{
			map[string]any{"options": []string{"Red", "M"}, "price": "5.50", "stock_qty": 12, "sku": p.str("variants.0.sku")},
			map[string]any{"options": []string{"Black", "L"}, "price": "6.50", "stock_qty": 3},
			map[string]any{"options": []string{"Black", "M"}, "price": "6", "stock_qty": 2, "sku": "TEE-BLK-M"},
		}}), tok).want(200)
		got := map[string]string{}
		for i := 0; i < r.len("variants"); i++ {
			got[r.str("variants."+itoa(i)+".options.0")+"/"+r.str("variants."+itoa(i)+".options.1")] = r.str("variants." + itoa(i) + ".id")
		}
		check(t, got["Red/M"] == ids["Red/M"] && got["Black/M"] == ids["Black/M"], "ID varian berubah: %v vs %v", got, ids)
		_, stillThere := got["Red/L"]
		check(t, !stillThere, "varian yang dihapus masih ada")
		r.eq("translation_pending", "false") // opsi sama → tidak perlu terjemahan ulang
		check(t, e.scalar(t, `SELECT count(*) FROM inventory_movements m JOIN product_variants v ON v.id = m.variant_id WHERE v.product_id = $1`, pid) >= "4",
			"pergerakan stok tidak tercatat")
	})

	t.Run("SKU dipakai produk lain ditolak", func(t *testing.T) {
		e.call(t, "POST", base, product(map[string]any{"name": "Another tee", "variants": []any{
			map[string]any{"options": []string{"Red", "M"}, "price": "5", "stock_qty": 1, "sku": "TEE-BLK-M"}}}), tok).wantErr(422, "validation", "variants")
	})

	var flagged resp
	t.Run("kirim ke moderasi + cek kata terlarang", func(t *testing.T) {
		flagged = e.call(t, "PUT", base+"/"+pid, product(map[string]any{"description": "Soft tee, not a replica.", "submit": true}), tok).want(200)
		flagged.eq("status", "active").eq("moderation_status", "pending_review").eq("visible", "false")
		check(t, strings.Contains(string(flagged.Raw), `"keyword":"replica"`), "kata 'replica' tidak ditandai: %s", flagged.str("auto_flags"))
	})
	t.Run("'gun' tidak cocok dengan 'burgundy'", func(t *testing.T) {
		r := e.call(t, "POST", base, product(map[string]any{"name": "Burgundy shirt", "description": "Colour burgundy", "submit": true,
			"variants": []any{map[string]any{"options": []string{"Red", "L"}, "price": "5", "stock_qty": 1}}}), tok).want(201)
		check(t, !strings.Contains(string(r.Raw), `"keyword":"gun"`), "salah tandai 'gun' di 'burgundy'")
	})

	t.Run("terjemahan yang diedit seller tidak ditimpa AI", func(t *testing.T) {
		e.drainTranslations(t)
		e.call(t, "PUT", base+"/"+pid+"/translations/km", map[string]string{"name": "អាវយឺត", "description": "ក"}, tok).want(200).eq("translations.km.seller_edited", "true")
		e.call(t, "PUT", base+"/"+pid+"/translations/en", map[string]string{"name": "x"}, tok).wantErr(422, "validation", "lang")
		e.call(t, "PUT", base+"/"+pid, product(map[string]any{"name": "Cotton T-shirt premium", "submit": true}), tok).want(200)
		e.drainTranslations(t)
		r := e.call(t, "GET", base+"/"+pid, nil, tok).want(200)
		r.eq("translations.km.name", "អាវយឺត").eq("translations.zh.name", "[zh] Cotton T-shirt premium")
		e.call(t, "DELETE", base+"/"+pid+"/translations/km", nil, tok).want(200)
		e.drainTranslations(t)
		e.call(t, "GET", base+"/"+pid, nil, tok).want(200).eq("translations.km.name", "[km] Cotton T-shirt premium")
	})

	admin := newAdmin(t, e, "moderator")
	t.Run("moderasi: minta perbaikan → seller kirim ulang → setujui", func(t *testing.T) {
		q := e.call(t, "GET", "/admin/api/products", nil, admin).want(200)
		check(t, strings.Contains(string(q.Raw), pid), "produk tidak ada di antrean")
		e.call(t, "POST", "/admin/api/products/"+pid+"/decision", map[string]string{"decision": "reject", "note": ""}, admin).want(422)
		e.call(t, "POST", "/admin/api/products/"+pid+"/decision", map[string]string{"decision": "changes_requested", "reason_code": "unclear_condition", "note": "Tambah tabel ukuran"}, admin).
			want(200).eq("product.moderation_status", "changes_requested")
		l := e.call(t, "GET", base+"?filter=changes", nil, tok).want(200)
		check(t, strings.Contains(string(l.Raw), "Tambah tabel ukuran"), "seller tidak melihat catatan moderator")
		e.call(t, "PUT", base+"/"+pid, product(map[string]any{"name": "Cotton T-shirt premium", "description": "Size: M 100cm, L 106cm"}), tok).
			want(200).eq("moderation_status", "pending_review")
		e.call(t, "PUT", "/admin/api/products/"+pid+"/translations/zh", map[string]string{"name": "纯棉T恤", "description": "尺码表"}, admin).
			want(200).eq("product.translations.zh.source", "reviewer")
		d := e.call(t, "POST", "/admin/api/products/"+pid+"/decision", map[string]string{"decision": "approve"}, admin).want(200)
		d.eq("product.moderation_status", "approved").eq("product.visible", "true")
		check(t, d.len("history") == 2, "riwayat moderasi harus 2, dapat %d", d.len("history"))
	})

	t.Run("ubah harga/stok setelah disetujui: tetap tayang", func(t *testing.T) {
		cur := e.call(t, "GET", base+"/"+pid, nil, tok).want(200)
		r := e.call(t, "PATCH", base+"/"+pid+"/variants", map[string]any{"variants": []any{
			map[string]any{"id": cur.str("variants.0.id"), "price": "4.99", "stock_qty": 20}}}, tok).want(200)
		r.eq("moderation_status", "approved").eq("min_price", "4.99")
		r = e.call(t, "PUT", base+"/"+pid, product(map[string]any{"name": "Cotton T-shirt premium", "description": "Size: M 100cm, L 106cm",
			"variants": []any{map[string]any{"options": []string{"Red", "M"}, "price": "7", "stock_qty": 3}}}), tok).want(200)
		r.eq("moderation_status", "approved")
	})
	t.Run("ubah foto setelah disetujui: dicek ulang", func(t *testing.T) {
		r := e.call(t, "PUT", base+"/"+pid, product(map[string]any{"name": "Cotton T-shirt premium", "description": "Size: M 100cm, L 106cm",
			"media": []any{map[string]any{"key": mp4.str("key"), "kind": "video"}, map[string]any{"key": png.str("key"), "kind": "photo"}}}), tok).want(200)
		r.eq("moderation_status", "pending_review")
	})
	t.Run("arsip & tampilkan lagi", func(t *testing.T) {
		e.call(t, "POST", base+"/"+pid+"/status", map[string]string{"status": "archived"}, tok).want(200).eq("status", "archived").eq("visible", "false")
		e.call(t, "POST", base+"/"+pid+"/status", map[string]string{"status": "active"}, tok).want(200).eq("status", "active")
	})
	t.Run("barang bekas wajib foto kondisi saat dikirim", func(t *testing.T) {
		used := product(map[string]any{"name": "Used jacket", "condition": "used_good", "condition_note": "Lecet kecil di resleting", "submit": true,
			"options": []any{}, "variants": []any{map[string]any{"options": []string{}, "price": "12", "stock_qty": 1}},
			"media": []any{map[string]any{"key": mov.str("key"), "kind": "video"}, map[string]any{"key": png.str("key"), "kind": "photo"}}})
		e.call(t, "POST", base, used, tok).wantErr(422, "validation", "media")
		def := e.upload(t, tok, "product_media", jpegBytes).want(201)
		used["media"] = append(used["media"].([]any), map[string]any{"key": def.str("key"), "kind": "defect_photo"})
		e.call(t, "POST", base, used, tok).want(201).eq("moderation_status", "pending_review")
	})
	t.Run("toko belum aktif: draf boleh, kirim tidak", func(t *testing.T) {
		tok3, _ := e.login(t)
		st := e.call(t, "POST", "/v1/seller/stores", map[string]any{"name": "Toko baru", "about": "", "about_lang": "km"}, tok3).want(201)
		v := e.upload(t, tok3, "product_media", mp4Bytes).want(201)
		f := e.upload(t, tok3, "product_media", jpegBytes).want(201)
		body := product(map[string]any{"media": []any{map[string]any{"key": v.str("key"), "kind": "video"}, map[string]any{"key": f.str("key"), "kind": "photo"}}})
		e.call(t, "POST", "/v1/seller/stores/"+st.str("id")+"/products", body, tok3).want(201)
		body["submit"] = true
		e.call(t, "POST", "/v1/seller/stores/"+st.str("id")+"/products", body, tok3).wantErr(409, "store_not_active")
	})
	t.Run("hapus", func(t *testing.T) {
		e.call(t, "DELETE", base+"/"+pid, nil, tok).want(204)
		e.call(t, "GET", base+"/"+pid, nil, tok).want(404)
	})
}

// ---------------------------------------------------------------------------
// 4. Halaman web, data publik, pengaturan
// ---------------------------------------------------------------------------

func TestWebAndPublic(t *testing.T) {
	e := getEnv(t)
	for _, p := range []string{"/healthz", "/readyz", "/v1/app/meta", "/seller/", "/admin/", "/dev/login",
		"/ui/common.js", "/ui/common.css", "/ui/maps.js", "/ui/seller.js", "/ui/seller-products.js", "/ui/admin.js", "/ui/admin-products.js", "/ui/seller-team.js", "/ui/seller-orders.js", "/ui/admin-refunds.js", "/ui/courier.js", "/ui/admin-delivery.js", "/ui/seller-payouts.js", "/ui/admin-payouts.js"} {
		t.Run(p, func(t *testing.T) { e.call(t, "GET", p, nil, "").want(200) })
	}
	t.Run("config.js berisi kunci Google Maps", func(t *testing.T) {
		r := e.call(t, "GET", "/ui/config.js", nil, "").want(200)
		check(t, strings.Contains(string(r.Raw), `"googleMapsKey":"qa-maps-key"`), "config.js: %s", r.Raw)
	})
	t.Run("kategori 3 bahasa + allows_used", func(t *testing.T) {
		r := e.call(t, "GET", "/v1/app/categories", nil, "", "X-Lang", "km").want(200)
		check(t, r.len("categories") >= 5, "kategori kurang")
		check(t, strings.Contains(string(r.Raw), `"allows_used":false`), "allows_used tidak ada")
	})
	t.Run("cek wilayah layanan", func(t *testing.T) {
		e.call(t, "GET", "/v1/app/service-areas/check?lat=10.6105&lng=103.5316", nil, "").want(200).eq("served", "true")
		e.call(t, "GET", "/v1/app/service-areas/check?lat=11.5564&lng=104.9282", nil, "").want(200).eq("served", "false")
	})
	t.Run("folder media tidak bisa di-list", func(t *testing.T) {
		e.call(t, "GET", "/media/public/", nil, "").want(404)
	})
	t.Run("halaman web memakai header keamanan", func(t *testing.T) {
		r := e.call(t, "GET", "/seller/", nil, "")
		check(t, r.Header.Get("X-Frame-Options") == "DENY", "X-Frame-Options hilang")
	})
	t.Run("route tidak dikenal → 404 JSON", func(t *testing.T) {
		e.call(t, "GET", "/v1/app/tidak-ada", nil, "").want(404)
	})
}

func itoa(i int) string { return strconv.Itoa(i) }

// ---------------------------------------------------------------------------
// 5. Terjemahan gagal → status jelas → coba lagi
// ---------------------------------------------------------------------------

type failingTranslator struct{}

func (failingTranslator) Translate(context.Context, translate.Request) (translate.Result, error) {
	return translate.Result{}, errors.New("claude: HTTP 400 invalid_request_error: Your credit balance is too low (saldo API habis)")
}

func TestTranslationFailureAndRetry(t *testing.T) {
	e := getEnv(t)
	tok, sid := newActiveStore(t, e)
	e.drainTranslations(t) // kosongkan antrean (deskripsi toko) supaya worker berikutnya memproses produk ini
	base := "/v1/seller/stores/" + sid + "/products"
	v := e.upload(t, tok, "product_media", mp4Bytes).want(201)
	f := e.upload(t, tok, "product_media", jpegBytes).want(201)
	p := e.call(t, "POST", base, map[string]any{"category_id": "kitchenware", "source_lang": "en", "name": "Rice cooker 1.8L", "description": "Non-stick pot",
		"condition": "new", "variants": []any{map[string]any{"options": []string{}, "price": "25", "stock_qty": 3}},
		"media": []any{map[string]any{"key": v.str("key"), "kind": "video"}, map[string]any{"key": f.str("key"), "kind": "photo"}}}, tok).want(201)
	pid := p.str("id")
	p.eq("translation_state", "pending")

	e.worker.T = failingTranslator{}
	defer func() { e.worker.T = translate.Fake{} }()
	if _, err := e.worker.RunOnce(context.Background()); err != nil {
		t.Fatal(err)
	}
	r := e.call(t, "GET", base+"/"+pid, nil, tok).want(200)
	r.eq("translation_state", "retrying")
	check(t, strings.Contains(r.str("translation_error"), "credit balance"), "pesan error tidak sampai ke seller: %q", r.str("translation_error"))

	t.Run("job gagal menunggu jeda (tidak langsung diulang)", func(t *testing.T) {
		did, _ := e.worker.RunOnce(context.Background())
		check(t, !did, "job gagal langsung diulang tanpa jeda")
	})
	t.Run("setelah 5x gagal → failed", func(t *testing.T) {
		e.sql(t, `UPDATE translation_jobs SET attempts = 4, run_after = NULL WHERE entity_id = $1`, pid)
		_, _ = e.worker.RunOnce(context.Background())
		e.call(t, "GET", base+"/"+pid, nil, tok).want(200).eq("translation_state", "failed")
	})
	t.Run("coba lagi sekarang → berhasil", func(t *testing.T) {
		e.worker.T = translate.Fake{}
		e.call(t, "POST", base+"/"+pid+"/translations/retry", nil, tok).want(200).eq("translation_state", "pending")
		e.drainTranslations(t)
		e.call(t, "GET", base+"/"+pid, nil, tok).want(200).eq("translation_state", "done").eq("translations.km.name", "[km] Rice cooker 1.8L")
	})
	t.Run("seller bisa menulis terjemahan sendiri walau AI gagal", func(t *testing.T) {
		e.call(t, "PUT", base+"/"+pid+"/translations/zh", map[string]string{"name": "电饭煲 1.8L", "description": "不粘锅"}, tok).
			want(200).eq("translations.zh.seller_edited", "true")
	})
}

// ---------------------------------------------------------------------------
// 6. Notifikasi seller (lonceng + Telegram) dan staf toko
// ---------------------------------------------------------------------------

type fakeTelegram struct {
	sent []string
	fail error
}

func (f *fakeTelegram) Send(_ context.Context, chatID int64, text, _, url string) error {
	if f.fail != nil {
		return f.fail
	}
	f.sent = append(f.sent, strconv.FormatInt(chatID, 10)+"|"+text+"|"+url)
	return nil
}

// submittedStore: toko lengkap yang sudah dikirim untuk review, KYC & titik sudah disetujui admin.
func submittedStore(t *testing.T, e *env, admin string) (tok, storeID, ownerPhone string) {
	t.Helper()
	tok, ownerPhone = e.login(t)
	storeID = e.call(t, "POST", "/v1/seller/stores", map[string]any{"name": "Notif Shop " + ownerPhone[len(ownerPhone)-4:], "about": "Toko", "about_lang": "en"}, tok).want(201).str("id")
	logo := e.upload(t, tok, "store_logo", pngBytes).want(201)
	e.call(t, "PATCH", "/v1/seller/stores/"+storeID, map[string]any{"logo_key": logo.str("key")}, tok).want(200)
	doc := e.upload(t, tok, "kyc_document", pdfBytes).want(201)
	e.call(t, "PUT", "/v1/seller/stores/"+storeID+"/kyc", map[string]any{"seller_type": "individual_unregistered", "legal_name": "Notif Dara", "id_document_key": doc.str("key")}, tok).want(200)
	pt := e.call(t, "POST", "/v1/seller/stores/"+storeID+"/pickup-points", pickup("Gudang"), tok).want(201).str("pickup_points.0.id")
	e.call(t, "POST", "/v1/seller/stores/"+storeID+"/submit", nil, tok).want(200)
	e.call(t, "POST", "/admin/api/stores/"+storeID+"/kyc", map[string]string{"decision": "approve"}, admin).want(200)
	e.call(t, "POST", "/admin/api/pickup-points/"+pt+"/decision", map[string]string{"decision": "approve"}, admin).want(200)
	return
}

func TestNotifications(t *testing.T) {
	e := getEnv(t)
	admin := newAdmin(t, e, "moderator")
	tok, sid, phone := submittedStore(t, e, admin)
	e.call(t, "POST", "/admin/api/stores/"+sid+"/decision", map[string]string{"decision": "approve"}, admin).want(200)

	r := e.call(t, "GET", "/v1/seller/notifications", nil, tok, "X-Lang", "km").want(200)
	check(t, r.str("unread") == "2", "harusnya 2 notifikasi (titik + toko), dapat %s: %s", r.str("unread"), r.Raw)
	r.eq("notifications.0.kind", "store.approved").eq("notifications.0.link", "#/store/"+sid)
	check(t, strings.Contains(r.str("notifications.0.title"), "អនុម័ត") && strings.Contains(r.str("notifications.0.body"), "Notif Shop"),
		"teks Khmer/nama toko salah: %s / %s", r.str("notifications.0.title"), r.str("notifications.0.body"))

	t.Run("user lain tidak melihat notifikasi ini", func(t *testing.T) {
		other, _ := e.login(t)
		e.call(t, "GET", "/v1/seller/notifications", nil, other).want(200).eq("unread", "0")
	})

	t.Run("dikirim ke Telegram (yang login dengan Telegram)", func(t *testing.T) {
		e.sql(t, `INSERT INTO auth_identities (user_id, provider, provider_uid, display)
			SELECT id, 'telegram', '777000111', '@notifshop' FROM users WHERE phone = '+855' || substr($1, 2)`, phone)
		bot := &fakeTelegram{}
		w := &notify.Worker{Pool: e.pool, Bot: bot, BaseURL: "https://phsar.test", Log: e.worker.Log}
		// Kirim semua yang antre (tes lain juga membuat notifikasi).
		total := 0
		for i := 0; i < 50; i++ {
			n, err := w.RunOnce(context.Background())
			check(t, err == nil, "worker: %v", err)
			if n == 0 {
				break
			}
			total += n
		}
		check(t, total >= 2, "worker: n=%d", total)
		joined := strings.Join(bot.sent, "\n")
		check(t, strings.Contains(joined, "777000111|") && strings.Contains(joined, "Notif Shop") && strings.Contains(joined, "https://phsar.test/seller/#/store/"+sid),
			"pesan Telegram salah: %s", joined)
		check(t, e.scalar(t, `SELECT count(*) FROM notifications WHERE store_id = $1 AND telegram_status = 'sent'`, sid) == "2", "status bukan sent")
	})

	t.Run("gagal sementara → dicoba lagi nanti; diblokir → berhenti", func(t *testing.T) {
		pid := newProductForModeration(t, e, tok, sid)
		e.call(t, "POST", "/admin/api/products/"+pid+"/decision", map[string]string{"decision": "changes_requested", "reason_code": "price", "note": "Harga terlalu murah"}, admin).want(200)
		w := &notify.Worker{Pool: e.pool, Bot: &fakeTelegram{fail: errors.New("telegram 502: Bad Gateway")}, Log: e.worker.Log}
		_, _ = w.RunOnce(context.Background())
		e.call(t, "GET", "/v1/seller/notifications", nil, tok).want(200).eq("notifications.0.kind", "product.changes_requested")
		check(t, e.scalar(t, `SELECT telegram_status FROM notifications WHERE store_id = $1 ORDER BY id DESC LIMIT 1`, sid) == "pending", "harusnya pending (retry)")
		n, _ := w.RunOnce(context.Background())
		check(t, n == 0, "dicoba lagi tanpa jeda")
		e.sql(t, `UPDATE notifications SET next_attempt_at = NULL WHERE store_id = $1`, sid)
		w.Bot = &fakeTelegram{fail: fmt.Errorf("%w: telegram 403: bot was blocked by the user", notify.ErrPermanent)}
		_, _ = w.RunOnce(context.Background())
		check(t, e.scalar(t, `SELECT telegram_status FROM notifications WHERE store_id = $1 ORDER BY id DESC LIMIT 1`, sid) == "failed", "harusnya failed")
	})

	t.Run("user tanpa Telegram → skipped (tetap ada di lonceng)", func(t *testing.T) {
		tok2, sid2, _ := submittedStore(t, e, admin)
		w := &notify.Worker{Pool: e.pool, Bot: &fakeTelegram{}, Log: e.worker.Log}
		_, _ = w.RunOnce(context.Background())
		check(t, e.scalar(t, `SELECT count(*) FROM notifications WHERE store_id = $1 AND telegram_status = 'skipped'`, sid2) == "1", "harusnya skipped")
		e.call(t, "GET", "/v1/seller/notifications", nil, tok2).want(200).eq("unread", "1")
	})

	t.Run("tandai sudah dibaca", func(t *testing.T) {
		first := e.call(t, "GET", "/v1/seller/notifications", nil, tok).want(200)
		e.call(t, "POST", "/v1/seller/notifications/read", map[string]any{"ids": []any{first.get("notifications.0.id")}}, tok).want(204)
		left := e.call(t, "GET", "/v1/seller/notifications", nil, tok).want(200)
		left.eq("notifications.0.read", "true")
		e.call(t, "POST", "/v1/seller/notifications/read", map[string]any{"all": true}, tok).want(204)
		e.call(t, "GET", "/v1/seller/notifications", nil, tok).want(200).eq("unread", "0")
	})
}

func newProductForModeration(t *testing.T, e *env, tok, sid string) string {
	t.Helper()
	v := e.upload(t, tok, "product_media", mp4Bytes).want(201)
	f := e.upload(t, tok, "product_media", jpegBytes).want(201)
	return e.call(t, "POST", "/v1/seller/stores/"+sid+"/products", map[string]any{"category_id": "kitchenware", "source_lang": "en", "name": "Thermos 1L",
		"condition": "new", "variants": []any{map[string]any{"options": []string{}, "price": "9", "stock_qty": 5}}, "submit": true,
		"media": []any{map[string]any{"key": v.str("key"), "kind": "video"}, map[string]any{"key": f.str("key"), "kind": "photo"}}}, tok).want(201).str("id")
}

func TestTeam(t *testing.T) {
	e := getEnv(t)
	owner, sid := newActiveStore(t, e)
	staff, _ := e.login(t)
	stranger, _ := e.login(t)

	inv := e.call(t, "POST", "/v1/seller/stores/"+sid+"/invites", nil, owner).want(201)
	token := inv.str("token")
	check(t, len(token) >= 32, "token terlalu pendek")
	t.Run("token disimpan sebagai hash", func(t *testing.T) {
		check(t, e.scalar(t, `SELECT count(*) FROM store_invites WHERE token_hash = $1`, token) == "0", "token mentah tersimpan di database")
	})
	t.Run("detail undangan", func(t *testing.T) {
		e.call(t, "GET", "/v1/seller/invites/"+token, nil, staff).want(200).eq("store_id", sid).eq("already_member", "false")
		e.call(t, "GET", "/v1/seller/invites/salah", nil, staff).wantErr(410, "invite_invalid")
	})
	e.call(t, "POST", "/v1/seller/invites/"+token+"/accept", nil, staff).want(200).eq("store_id", sid)

	t.Run("link sekali pakai", func(t *testing.T) {
		e.call(t, "POST", "/v1/seller/invites/"+token+"/accept", nil, stranger).wantErr(410, "invite_invalid")
	})
	t.Run("owner dapat notifikasi staf baru", func(t *testing.T) {
		e.call(t, "GET", "/v1/seller/notifications", nil, owner).want(200).eq("notifications.0.kind", "staff.joined")
		e.call(t, "GET", "/v1/seller/notifications", nil, staff).want(200).eq("unread", "0")
	})
	t.Run("staf: boleh kelola produk", func(t *testing.T) {
		e.call(t, "GET", "/v1/seller/stores/"+sid, nil, staff).want(200).eq("role", "staff")
		e.call(t, "GET", "/v1/seller/stores/"+sid+"/products", nil, staff).want(200)
		v := e.upload(t, staff, "product_media", mp4Bytes).want(201)
		f := e.upload(t, staff, "product_media", jpegBytes).want(201)
		e.call(t, "POST", "/v1/seller/stores/"+sid+"/products", map[string]any{"category_id": "kitchenware", "source_lang": "km", "name": "ឆ្នាំងបាយ",
			"condition": "new", "variants": []any{map[string]any{"options": []string{}, "price": "20", "stock_qty": 2}},
			"media": []any{map[string]any{"key": v.str("key"), "kind": "video"}, map[string]any{"key": f.str("key"), "kind": "photo"}}}, staff).want(201)
	})
	t.Run("staf: tidak boleh profil, KYC, titik, tim", func(t *testing.T) {
		e.call(t, "PATCH", "/v1/seller/stores/"+sid, map[string]any{"name": "Ganti"}, staff).wantErr(403, "owner_only")
		e.call(t, "GET", "/v1/seller/stores/"+sid+"/kyc", nil, staff).wantErr(403, "owner_only")
		e.call(t, "POST", "/v1/seller/stores/"+sid+"/pickup-points", pickup("X"), staff).wantErr(403, "owner_only")
		e.call(t, "GET", "/v1/seller/stores/"+sid+"/team", nil, staff).wantErr(403, "owner_only")
		e.call(t, "POST", "/v1/seller/stores/"+sid+"/invites", nil, staff).wantErr(403, "owner_only")
	})
	team := e.call(t, "GET", "/v1/seller/stores/"+sid+"/team", nil, owner).want(200)
	check(t, team.len("members") == 2, "anggota harus 2, dapat %d", team.len("members"))
	team.eq("members.0.role", "owner").eq("members.1.role", "staff")
	check(t, strings.Contains(team.str("members.1.phone"), "***"), "nomor HP staf tidak disamarkan: %s", team.str("members.1.phone"))
	staffID := team.str("members.1.user_id")
	ownerID := team.str("members.0.user_id")

	t.Run("undangan dibatalkan / kedaluwarsa tidak bisa dipakai", func(t *testing.T) {
		i2 := e.call(t, "POST", "/v1/seller/stores/"+sid+"/invites", nil, owner).want(201)
		e.call(t, "DELETE", "/v1/seller/stores/"+sid+"/invites/"+i2.str("id"), nil, owner).want(204)
		e.call(t, "POST", "/v1/seller/invites/"+i2.str("token")+"/accept", nil, stranger).wantErr(410, "invite_invalid")
		i3 := e.call(t, "POST", "/v1/seller/stores/"+sid+"/invites", nil, owner).want(201)
		e.sql(t, `UPDATE store_invites SET expires_at = now() - interval '1 minute' WHERE id = $1`, i3.str("id"))
		e.call(t, "POST", "/v1/seller/invites/"+i3.str("token")+"/accept", nil, stranger).wantErr(410, "invite_invalid")
	})
	t.Run("owner tidak bisa dihapus; staf tidak bisa menghapus orang lain", func(t *testing.T) {
		e.call(t, "DELETE", "/v1/seller/stores/"+sid+"/members/"+ownerID, nil, owner).wantErr(409, "cannot_remove_owner")
		e.call(t, "DELETE", "/v1/seller/stores/"+sid+"/members/"+ownerID, nil, staff).wantErr(403, "owner_only")
	})
	t.Run("owner menghapus staf → akses hilang", func(t *testing.T) {
		e.call(t, "DELETE", "/v1/seller/stores/"+sid+"/members/"+staffID, nil, owner).want(204)
		e.call(t, "GET", "/v1/seller/stores/"+sid, nil, staff).want(404)
	})
	t.Run("staf bisa keluar sendiri", func(t *testing.T) {
		i4 := e.call(t, "POST", "/v1/seller/stores/"+sid+"/invites", nil, owner).want(201)
		e.call(t, "POST", "/v1/seller/invites/"+i4.str("token")+"/accept", nil, stranger).want(200)
		me := e.call(t, "GET", "/v1/seller/stores/"+sid+"/team", nil, owner).want(200)
		id := me.str("members.1.user_id")
		e.call(t, "DELETE", "/v1/seller/stores/"+sid+"/members/"+id, nil, stranger).want(204)
		e.call(t, "GET", "/v1/seller/stores/"+sid, nil, stranger).want(404)
	})
	t.Run("banner toko", func(t *testing.T) {
		b := e.upload(t, owner, "store_banner", jpegBytes).want(201)
		e.call(t, "PATCH", "/v1/seller/stores/"+sid, map[string]any{"banner_key": b.str("key")}, owner).want(200)
		check(t, e.call(t, "GET", "/v1/seller/stores/"+sid, nil, owner).want(200).str("banner_url") != "", "banner_url kosong")
	})
}

// ---------------------------------------------------------------------------
// 7. Kompres video + thumbnail, foto diperkecil (butuh ffmpeg)
// ---------------------------------------------------------------------------

func TestMediaProcessing(t *testing.T) {
	e := getEnv(t)
	if e.media == nil {
		t.Skip("ffmpeg tidak terpasang (Mac: brew install ffmpeg) — kompres video tidak dites")
	}
	tok, sid := newActiveStore(t, e)
	base := "/v1/seller/stores/" + sid + "/products"

	t.Run("video lebih dari 3 menit ditolak saat upload", func(t *testing.T) {
		e.upload(t, tok, "product_media", longVideo).wantErr(422, "video_too_long")
	})
	t.Run("file berkedok video ditolak", func(t *testing.T) {
		fake := append([]byte{0, 0, 0, 24, 'f', 't', 'y', 'p', 'i', 's', 'o', 'm', 0, 0, 2, 0, 'i', 's', 'o', 'm', 'i', 's', 'o', '2'}, make([]byte, 3000)...)
		e.upload(t, tok, "product_media", fake).wantErr(415, "upload_type")
	})

	v := e.upload(t, tok, "product_media", movBytes).want(201)
	f := e.upload(t, tok, "product_media", jpegBytes).want(201)
	p := e.call(t, "POST", base, map[string]any{"category_id": "kitchenware", "source_lang": "en", "name": "Blender 1.5L", "condition": "new",
		"variants": []any{map[string]any{"options": []string{}, "price": "30", "stock_qty": 4}},
		"media":    []any{map[string]any{"key": v.str("key"), "kind": "video"}, map[string]any{"key": f.str("key"), "kind": "photo"}}}, tok).want(201)
	pid := p.str("id")
	p.eq("media.0.processing", "processing").eq("media.1.processing", "processing")
	check(t, p.str("media.0.url") == p.str("media.0.original_url"), "sebelum diproses, file asli yang tampil")

	e.drainMedia(t)
	r := e.call(t, "GET", base+"/"+pid, nil, tok).want(200)
	r.eq("media.0.processing", "ready").eq("media.1.processing", "ready")
	check(t, strings.HasSuffix(r.str("media.0.url"), "-720.mp4") || r.str("media.0.url") == r.str("media.0.original_url"), "video: %s", r.str("media.0.url"))
	check(t, strings.HasSuffix(r.str("media.0.poster_url"), "-poster.jpg"), "poster: %s", r.str("media.0.poster_url"))
	check(t, r.str("media.0.duration_sec") == "2", "durasi %s", r.str("media.0.duration_sec"))
	check(t, strings.HasSuffix(r.str("media.1.url"), "-1600.jpg"), "foto: %s", r.str("media.1.url"))
	r.eq("media.1.width", "1600").eq("media.1.height", "1200")
	for _, u := range []string{r.str("media.0.url"), r.str("media.0.poster_url"), r.str("media.1.url")} {
		e.call(t, "GET", strings.TrimPrefix(u, e.cfg.PublicBaseURL), nil, "").want(200)
	}
	t.Run("sampul di daftar produk memakai foto yang sudah diperkecil", func(t *testing.T) {
		l := e.call(t, "GET", base, nil, tok).want(200)
		check(t, strings.HasSuffix(l.str("products.0.cover_url"), "-1600.jpg"), "cover: %s", l.str("products.0.cover_url"))
	})
	t.Run("ubah urutan tidak memproses ulang", func(t *testing.T) {
		n := e.scalar(t, `SELECT count(*) FROM media_processing_jobs j JOIN product_media m ON m.id = j.media_id WHERE m.product_id = $1`, pid)
		e.call(t, "PUT", base+"/"+pid, map[string]any{"category_id": "kitchenware", "source_lang": "en", "name": "Blender 1.5L", "condition": "new",
			"variants": []any{map[string]any{"options": []string{}, "price": "30", "stock_qty": 4}},
			"media":    []any{map[string]any{"key": f.str("key"), "kind": "photo"}, map[string]any{"key": v.str("key"), "kind": "video"}}}, tok).want(200)
		check(t, e.scalar(t, `SELECT count(*) FROM media_processing_jobs j JOIN product_media m ON m.id = j.media_id WHERE m.product_id = $1`, pid) == n,
			"media lama diproses ulang")
	})
	t.Run("foto rusak: langsung pakai file asli (tidak macet 'processing')", func(t *testing.T) {
		bad := e.upload(t, tok, "product_media", append([]byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 'J', 'F', 'I', 'F', 0x00, 0x01}, make([]byte, 600)...)).want(201)
		e.call(t, "PUT", base+"/"+pid, map[string]any{"category_id": "kitchenware", "source_lang": "en", "name": "Blender 1.5L", "condition": "new",
			"variants": []any{map[string]any{"options": []string{}, "price": "30", "stock_qty": 4}},
			"media": []any{map[string]any{"key": v.str("key"), "kind": "video"}, map[string]any{"key": f.str("key"), "kind": "photo"},
				map[string]any{"key": bad.str("key"), "kind": "photo"}}}, tok).want(200)
		e.drainMedia(t)
		r := e.call(t, "GET", base+"/"+pid, nil, tok).want(200)
		r.eq("media.2.processing", "ready")
		check(t, r.str("media.2.url") == r.str("media.2.original_url"), "foto rusak harus memakai file asli")
	})
	t.Run("backfill media lama", func(t *testing.T) {
		e.sql(t, `UPDATE product_media SET poster_url = NULL WHERE product_id = $1 AND kind = 'video'`, pid)
		n, err := media.Backfill(context.Background(), e.pool)
		check(t, err == nil && n == 1, "backfill: n=%d err=%v", n, err)
		e.drainMedia(t)
		e.call(t, "GET", base+"/"+pid, nil, tok).want(200)
	})
}
