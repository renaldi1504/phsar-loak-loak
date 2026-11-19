package qa

import (
	"strings"
	"testing"
)

// Halaman "Jual di Phsar": 3 bahasa, tarif komisi dari aturan admin, tombol ke Seller Center dengan bahasa yang sama.
func TestSellPage(t *testing.T) {
	e := getEnv(t)
	for lang, want := range map[string]string{"en": "Open your shop on Phsar", "km": "បើកហាងរបស់អ្នកលើ Phsar", "zh": "在 Phsar 免费开店"} {
		r := e.call(t, "GET", "/"+lang+"/sell", nil, "").want(200)
		html := string(r.Raw)
		check(t, strings.Contains(html, want), "%s: judul tidak ada", lang)
		check(t, strings.Contains(html, `href="/seller/?lang=`+lang+`"`), "%s: tombol daftar tidak ke Seller Center", lang)
		check(t, strings.Contains(html, "8%"), "%s: tarif komisi bawaan 8%% tidak tampil", lang)
	}
	// tarif global dari admin dipakai
	e.sql(t, `INSERT INTO commission_rules (scope, rate_percent, valid_from) VALUES ('global', 6.5, now())`)
	defer e.sql(t, `DELETE FROM commission_rules WHERE scope = 'global' AND rate_percent = 6.5`)
	check(t, strings.Contains(string(e.call(t, "GET", "/en/sell", nil, "").want(200).Raw), "Only 6.5% when you sell"), "tarif komisi admin tidak dipakai")
	// footer setiap halaman mengarah ke /sell
	check(t, strings.Contains(string(e.call(t, "GET", "/km/", nil, "").want(200).Raw), `href="/km/sell"`), "link Jual di Phsar di footer")
	check(t, strings.Contains(string(e.call(t, "GET", "/sitemap.xml", nil, "").want(200).Raw), "/en/sell"), "sitemap memuat /sell")
}
