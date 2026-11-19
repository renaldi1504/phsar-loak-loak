package qa

import (
	"strings"
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// ads: iklan produk bersponsor — seller buat kampanye → admin tinjau → tayang → klik ditagih dari saldo seller
// ---------------------------------------------------------------------------

const browserUA = "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 Mobile/15E148 Safari/604.1"

func TestAds(t *testing.T) {
	e := getEnv(t)
	mod := newAdmin(t, e, "moderator")
	rev := newAdmin(t, e, "ads_reviewer")
	s := newSellable(t, e, mod, "Zebrastripe ad sneakers", []string{"12.00"}, []int{5})
	other := newSellable(t, e, mod, "Other shop item", []string{"3.00"}, []int{5})
	base := "/v1/seller/stores/" + s.storeID + "/ads"
	slug := e.scalar(t, `SELECT slug FROM product_translations WHERE product_id = $1 AND lang = 'en'`, s.productID)
	catSlug := e.scalar(t, `SELECT slug FROM category_translations WHERE category_id = 'fashion' AND lang = 'en'`)
	searchURL := "/en/search?q=Zebrastripe"
	page := func(path string) string { return string(e.page(t, path, "User-Agent", browserUA).Raw) }
	shown := func(path string) bool { return strings.Contains(page(path), `class="card ad-card"`) }
	balance := func() string {
		return e.call(t, "GET", "/v1/seller/stores/"+s.storeID+"/payouts", nil, s.token).want(200).str("available")
	}
	e.sql(t, `INSERT INTO seller_ledger (store_id, kind, amount, note) VALUES ($1, 'adjustment', 5.00, 'test top-up')`, s.storeID)
	t.Cleanup(func() {
		e.sql(t, `UPDATE platform_settings SET value = '{"enabled": true, "slots": 2}' WHERE key = 'ads'`)
	})

	var id string
	t.Run("seller membuat kampanye", func(t *testing.T) {
		ok := map[string]any{"product_id": s.productID, "bid": "0.10", "daily_budget": "1.00"}
		with := func(k string, v any) map[string]any {
			m := map[string]any{}
			for a, b := range ok {
				m[a] = b
			}
			m[k] = v
			return m
		}
		e.call(t, "POST", base, with("bid", "0.01"), s.token).wantErr(422, "validation", "bid")
		e.call(t, "POST", base, with("bid", "2.50"), s.token).wantErr(422, "validation", "bid")
		e.call(t, "POST", base, with("daily_budget", "0.50"), s.token).wantErr(422, "validation", "daily_budget")
		e.call(t, "POST", base, with("product_id", other.productID), s.token).wantErr(422, "validation", "product_id")
		e.call(t, "POST", base, with("starts_on", "2020-01-01"), s.token).wantErr(422, "validation", "starts_on")
		e.call(t, "POST", base, ok, other.token).want(404) // bukan anggota toko
		r := e.call(t, "POST", base, ok, s.token).want(201)
		r.eq("status", "pending").eq("bid", "0.10").eq("live", "false")
		id = r.str("id")
		e.call(t, "POST", base, ok, s.token).wantErr(422, "validation", "product_id")
		check(t, !shown(searchURL), "kampanye yang belum ditinjau tidak boleh tayang")
	})

	t.Run("admin meninjau", func(t *testing.T) {
		e.call(t, "GET", "/admin/api/ads", nil, mod).want(403)
		l := e.call(t, "GET", "/admin/api/ads", nil, rev).want(200)
		check(t, l.str("counts.pending") != "0" && l.str("campaigns.0.id") == id, "kampanye harus menunggu tinjauan: %s", truncate(string(l.Raw)))
		e.call(t, "POST", "/admin/api/ads/"+id+"/action", map[string]string{"action": "reject"}, rev).wantErr(422, "validation", "note")
		e.call(t, "POST", "/admin/api/ads/"+id+"/action", map[string]string{"action": "approve"}, rev).want(200).eq("status", "approved").eq("live", "true")
		e.call(t, "POST", "/admin/api/ads/"+id+"/action", map[string]string{"action": "approve"}, rev).wantErr(409, "ad_state")
		check(t, e.scalar(t, `SELECT count(*) FROM notifications WHERE store_id = $1 AND kind = 'ad.approved'`, s.storeID) != "0", "seller harus diberi tahu")
	})

	t.Run("tayang di hasil cari & kategori, tanpa duplikat, tayangan dihitung", func(t *testing.T) {
		html := page(searchURL)
		check(t, strings.Contains(html, `class="card ad-card"`) && strings.Contains(html, `/ad/`+id+`?l=en`) && strings.Contains(html, "Sponsored"), "iklan harus tampil di hasil cari")
		check(t, !strings.Contains(html, `href="/en/p/`+slug+`"`), "produk bersponsor tidak boleh muncul dua kali")
		check(t, strings.Contains(page("/km/search?q=Zebrastripe"), "ផ្សាយពាណិជ្ជកម្ម"), "label Khmer")
		check(t, shown("/en/c/"+catSlug), "iklan harus tampil di halaman kategori")
		var imp string
		for i := 0; i < 30; i++ {
			if imp = e.scalar(t, `SELECT coalesce(sum(impressions), 0)::text FROM ad_stats_daily WHERE campaign_id = $1`, id); imp >= "3" {
				break
			}
			time.Sleep(50 * time.Millisecond)
		}
		check(t, imp == "3", "3 tayangan harus tercatat (bot tidak dihitung): %s", imp)
		check(t, strings.Contains(page("/robots.txt"), "Disallow: /ad/"), "robots.txt")
	})

	t.Run("klik: ditagih sekali per pengunjung per hari, bot gratis, maks. per IP", func(t *testing.T) {
		r := e.page(t, "/ad/"+id+"?l=en", "User-Agent", browserUA)
		check(t, r.Status == 302 && r.Header.Get("Location") == "/en/p/"+slug, "klik harus diarahkan ke produk: %d %s", r.Status, r.Header.Get("Location"))
		check(t, strings.Contains(r.Header.Get("Set-Cookie"), "phsar_vid="), "cookie pengunjung harus dipasang")
		vid := strings.SplitN(strings.TrimPrefix(strings.SplitN(r.Header.Get("Set-Cookie"), ";", 2)[0], "phsar_vid="), ";", 2)[0]
		check(t, balance() == "4.90", "saldo harus 5.00 − 0.10: %s", balance())
		e.page(t, "/ad/"+id+"?l=en", "User-Agent", browserUA, "Cookie", "phsar_vid="+vid)
		check(t, balance() == "4.90", "klik ulang pengunjung sama tidak ditagih: %s", balance())
		e.page(t, "/ad/"+id+"?l=en", "User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1)", "Cookie", "phsar_vid="+strings.Repeat("b", 32))
		check(t, balance() == "4.90", "bot tidak ditagih: %s", balance())
		for _, v := range []string{"c", "d", "f"} { // IP sama (127.0.0.1): maks. 3 tagihan per hari
			e.page(t, "/ad/"+id+"?l=en", "User-Agent", browserUA, "Cookie", "phsar_vid="+strings.Repeat(v, 32))
		}
		check(t, balance() == "4.70", "maks. 3 tagihan per IP per hari: %s", balance())
		e.call(t, "GET", base, nil, s.token).want(200).eq("campaigns.0.today.clicks", "3").eq("campaigns.0.today.spend", "0.30")
		check(t, e.scalar(t, `SELECT count(*) FROM seller_ledger WHERE store_id = $1 AND kind = 'ad_spend'`, s.storeID) == "3", "3 catatan ad_spend")
		r = e.page(t, "/ad/00000000-0000-0000-0000-000000000000?l=km", "User-Agent", browserUA)
		check(t, r.Status == 302 && r.Header.Get("Location") == "/km/", "kampanye tak dikenal → beranda")
	})

	t.Run("berhenti tayang: anggaran habis, saldo kurang, jeda, iklan dimatikan", func(t *testing.T) {
		e.sql(t, `UPDATE ad_stats_daily SET spend = 0.95 WHERE campaign_id = $1`, id)
		check(t, !shown(searchURL), "anggaran harian habis → tidak tayang")
		e.page(t, "/ad/"+id+"?l=en", "User-Agent", browserUA, "Cookie", "phsar_vid="+strings.Repeat("e", 32))
		e.sql(t, `UPDATE ad_stats_daily SET spend = 0.30 WHERE campaign_id = $1`, id)
		check(t, shown(searchURL), "anggaran tersedia lagi → tayang")

		e.sql(t, `INSERT INTO seller_ledger (store_id, kind, amount, note) VALUES ($1, 'adjustment', -4.65, 'test')`, s.storeID)
		check(t, !shown(searchURL), "saldo $0.05 < CPC → tidak tayang")
		e.sql(t, `INSERT INTO seller_ledger (store_id, kind, amount, note) VALUES ($1, 'adjustment', 4.65, 'test')`, s.storeID)

		e.call(t, "PATCH", base+"/"+id, map[string]any{"action": "pause"}, s.token).want(200).eq("status", "paused")
		check(t, !shown(searchURL), "dijeda seller → tidak tayang")
		e.call(t, "PATCH", base+"/"+id, map[string]any{"action": "resume", "bid": "0.20"}, s.token).want(200).eq("status", "approved").eq("bid", "0.20")
		e.call(t, "PATCH", base+"/"+id, map[string]any{"daily_budget": "0.10"}, s.token).wantErr(422, "validation", "daily_budget")

		e.call(t, "POST", "/admin/api/ads/"+id+"/action", map[string]string{"action": "pause", "note": "Misleading video"}, rev).want(200).eq("paused_by", "admin")
		e.call(t, "PATCH", base+"/"+id, map[string]any{"action": "resume"}, s.token).wantErr(409, "ad_state")
		e.call(t, "POST", "/admin/api/ads/"+id+"/action", map[string]string{"action": "resume"}, rev).want(200).eq("status", "approved")

		e.call(t, "PUT", "/admin/api/ads/settings", map[string]any{"enabled": true, "slots": 9}, rev).wantErr(422, "validation", "slots")
		e.call(t, "PUT", "/admin/api/ads/settings", map[string]any{"enabled": false, "slots": 2}, rev).want(200)
		check(t, !shown(searchURL), "iklan dimatikan admin → tidak tayang")
		e.call(t, "PUT", "/admin/api/ads/settings", map[string]any{"enabled": true, "slots": 2}, rev).want(200)
		check(t, shown(searchURL), "iklan dinyalakan lagi")

		e.call(t, "PATCH", base+"/"+id, map[string]any{"action": "end"}, s.token).want(200).eq("status", "ended")
		check(t, !shown(searchURL), "kampanye berakhir → tidak tayang")
		e.call(t, "POST", base, map[string]any{"product_id": s.productID, "bid": "0.10", "daily_budget": "1.00"}, s.token).want(201) // boleh buat baru setelah berakhir
	})
}
