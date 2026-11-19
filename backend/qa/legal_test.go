package qa

import (
	"strings"
	"testing"
)

// Dokumen legal: draf bawaan tampil di 3 bahasa; super admin menerbitkan versi baru (append-only) dan mengisi kontak.
func TestLegalPages(t *testing.T) {
	e := getEnv(t)
	sa := newAdmin(t, e, "super_admin")
	fin := newAdmin(t, e, "finance")
	t.Cleanup(func() {
		e.sql(t, `DELETE FROM platform_settings WHERE key = 'legal'`)
	})

	for _, p := range []string{"/km/legal/terms", "/en/legal/privacy", "/zh/legal/returns", "/en/legal/seller-terms", "/en/legal/data-deletion", "/km/contact"} {
		r := e.page(t, p)
		check(t, r.Status == 200, "%s harus 200, dapat %d", p, r.Status)
	}
	check(t, e.page(t, "/en/legal/nope").Status == 404, "dokumen tak dikenal harus 404")
	home := string(e.page(t, "/en/").Raw)
	check(t, strings.Contains(home, `/en/legal/returns`) && strings.Contains(home, `/en/contact`), "footer harus punya tautan legal")
	check(t, strings.Contains(string(e.page(t, "/sitemap.xml").Raw), "/km/legal/terms"), "sitemap harus memuat halaman legal")

	e.call(t, "GET", "/admin/api/legal", nil, fin).want(403)
	v := e.call(t, "GET", "/admin/api/legal", nil, sa).want(200)
	check(t, v.len("documents") == 5, "harus ada 5 dokumen")

	// kontak: validasi lalu simpan → muncul di halaman Kontak dan di dalam dokumen
	e.call(t, "PUT", "/admin/api/legal/settings", map[string]any{"company": "Phsar Co., Ltd.", "email": "bukan-email"}, sa).wantErr(422, "validation", "email")
	e.call(t, "PUT", "/admin/api/legal/settings", map[string]any{"company": "Phsar Co., Ltd.", "facebook": "javascript:alert(1)"}, sa).wantErr(422, "validation", "facebook")
	e.call(t, "PUT", "/admin/api/legal/settings", map[string]any{"company": "Phsar Co., Ltd.", "address": "Street 7, Sihanoukville", "email": "help@phsar.test", "telegram": "@phsar_help"}, sa).want(200)
	c := string(e.page(t, "/en/contact").Raw)
	check(t, strings.Contains(c, "help@phsar.test") && strings.Contains(c, "https://t.me/phsar_help"), "kontak harus tampil")
	check(t, strings.Contains(string(e.page(t, "/en/legal/privacy").Raw), "Phsar Co., Ltd."), "placeholder {{company}} harus terisi")

	// terbitkan versi baru: semua bahasa wajib, nomor versi unik, isi aman dari XSS
	long := strings.Repeat("Phsar keeps buyers safe. ", 4)
	bodies := map[string]string{"en": "## New returns\n\n" + long + "<script>alert(1)</script> [bad](javascript:alert(1)) [ok](/en/contact)", "km": "## ថ្មី\n\n" + long, "zh": "## 新\n\n" + long}
	e.call(t, "POST", "/admin/api/legal/docs/return_policy", map[string]any{"version": "1.0", "bodies": bodies}, sa).wantErr(422, "validation", "version")
	e.call(t, "POST", "/admin/api/legal/docs/return_policy", map[string]any{"version": "1.1", "bodies": map[string]string{"en": bodies["en"]}}, sa).wantErr(422, "validation", "body_km")
	e.call(t, "POST", "/admin/api/legal/docs/nope", map[string]any{"version": "1.1", "bodies": bodies}, sa).wantErr(422, "validation", "type")
	e.call(t, "POST", "/admin/api/legal/docs/return_policy", map[string]any{"version": "1.1", "bodies": bodies}, fin).want(403)
	e.call(t, "POST", "/admin/api/legal/docs/return_policy", map[string]any{"version": "1.1", "bodies": bodies}, sa).want(200)
	t.Cleanup(func() { e.sql(t, `DELETE FROM legal_documents WHERE type = 'return_policy' AND version = '1.1'`) })

	pg := string(e.page(t, "/en/legal/returns").Raw)
	check(t, strings.Contains(pg, "New returns") && strings.Contains(pg, "Version 1.1"), "versi baru harus tampil")
	check(t, !strings.Contains(pg, "<script>alert") && !strings.Contains(pg, `href="javascript:`) && strings.Contains(pg, `href="/en/contact"`), "isi dokumen harus aman")
	check(t, strings.Contains(string(e.page(t, "/km/legal/returns").Raw), "ថ្មី"), "versi Khmer harus tampil")
	check(t, e.scalar(t, `SELECT count(*) FROM legal_documents WHERE type = 'return_policy'`) == "6", "versi lama tetap tersimpan")

	pv := e.call(t, "POST", "/admin/api/legal/preview", map[string]any{"body": "## Hi\n\nMail {{email}}"}, sa).want(200)
	check(t, strings.Contains(pv.str("html"), "help@phsar.test") && strings.Contains(pv.str("html"), "<h2>"), "pratinjau harus terisi: %s", pv.str("html"))
}
