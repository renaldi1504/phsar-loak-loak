package qa

import (
	"regexp"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// Ongkir per zona & kurs riel (admin → Shipping & FX rates)
// ---------------------------------------------------------------------------

var rielRe = regexp.MustCompile(`riel:\s*4100\b`)

func TestShippingAndFX(t *testing.T) {
	e := getEnv(t)
	mod := newAdmin(t, e, "moderator")
	cs := newAdmin(t, e, "cs")
	fin := newAdmin(t, e, "finance")
	s := newSellable(t, e, mod, "Shipping mug", []string{"6.00"}, []int{30})
	buyer, _ := e.login(t)
	zone := map[string]any{"name": map[string]string{"en": "Sihanoukville", "km": "ក្រុងព្រះសីហនុ", "zh": "西哈努克城"}, "fee_per_store": "0", "free_over": "", "courier_fee": "0"}
	t.Cleanup(func() {
		e.sql(t, `UPDATE shipping_zones SET fee_per_store = 0, free_over = NULL, courier_fee = 0 WHERE code = 'sihanoukville'`)
		e.sql(t, `UPDATE service_areas SET is_active = true WHERE code = 'sihanoukville'`)
		e.sql(t, `DELETE FROM service_areas WHERE code = 'qa-otres'`)
		e.sql(t, `DELETE FROM exchange_rates WHERE created_by IS NOT NULL`)
	})
	quote := func(qty int) resp {
		return e.call(t, "POST", "/v1/shop/cart/quote", map[string]any{"items": []any{item(s.variants[0], qty)}}, "").want(200)
	}

	t.Run("default gratis ongkir (perilaku lama)", func(t *testing.T) {
		quote(1).eq("shipping_fee", "0.00").eq("total", "6.00").eq("stores.0.shipping_fee", "0.00")
		check(t, !strings.Contains(string(quote(1).Raw), "shipping_estimate"), "gratis: bukan perkiraan")
	})

	t.Run("izin & validasi", func(t *testing.T) {
		e.call(t, "GET", "/admin/api/shipping", nil, cs).want(403)
		v := e.call(t, "GET", "/admin/api/shipping", nil, fin).want(200)
		check(t, v.len("zones") >= 1 && v.len("areas") >= 1 && v.len("rates") >= 1, "data awal harus tampil: %s", truncate(string(v.Raw)))
		bad := map[string]any{}
		for k, val := range zone {
			bad[k] = val
		}
		bad["fee_per_store"] = "60"
		e.call(t, "PUT", "/admin/api/shipping/zones/sihanoukville", bad, fin).wantErr(422, "validation", "fee_per_store")
		bad["fee_per_store"], bad["free_over"] = "1", "abc"
		e.call(t, "PUT", "/admin/api/shipping/zones/sihanoukville", bad, fin).wantErr(422, "validation", "free_over")
		e.call(t, "PUT", "/admin/api/shipping/zones/nope", zone, fin).want(404)
		e.call(t, "POST", "/admin/api/shipping/fx", map[string]string{"riel_per_usd": "100"}, fin).wantErr(422, "validation", "riel_per_usd")
		area := map[string]any{"name": map[string]string{"en": "Sihanoukville"}, "zone": "sihanoukville", "active": false,
			"boundary": map[string]any{"type": "Polygon", "coordinates": [][][2]float64{{{103.44, 10.55}, {103.62, 10.55}, {103.62, 10.72}, {103.44, 10.72}, {103.44, 10.55}}}}}
		e.call(t, "PUT", "/admin/api/shipping/areas/sihanoukville", area, fin).wantErr(409, "last_area")
		area["code"], area["active"], area["boundary"] = "qa-otres", true, map[string]any{"type": "Point"}
		e.call(t, "POST", "/admin/api/shipping/areas", area, fin).wantErr(422, "validation", "boundary")
	})

	t.Run("ongkir per toko + gratis di atas batas", func(t *testing.T) {
		z := map[string]any{}
		for k, val := range zone {
			z[k] = val
		}
		z["fee_per_store"], z["free_over"], z["courier_fee"] = "1.50", "20", "1.00"
		e.call(t, "PUT", "/admin/api/shipping/zones/sihanoukville", z, fin).want(200).eq("zones.0.fee_per_store", "1.50").eq("zones.0.free_over", "20.00")
		q := quote(1)
		q.eq("subtotal", "6.00").eq("shipping_fee", "1.50").eq("total", "7.50").eq("stores.0.shipping_fee", "1.50").eq("shipping_estimate", "true")
		quote(4).eq("subtotal", "24.00").eq("shipping_fee", "0.00").eq("total", "24.00")

		e.call(t, "POST", "/v1/shop/checkouts", checkoutBody("k-ship-0001", "6.00", item(s.variants[0], 1)), buyer).wantErr(409, "cart_changed")
		c := e.call(t, "POST", "/v1/shop/checkouts", checkoutBody("k-ship-0002", "7.50", item(s.variants[0], 1)), buyer).want(201)
		c.eq("total", "7.50").eq("subtotal", "6.00").eq("shipping_fee", "1.50")
		check(t, e.scalar(t, `SELECT shipping_fee::text || '/' || total::text || '/' || seller_amount::text FROM orders WHERE checkout_id = $1`, c.str("id")) == "1.50/7.50/5.52",
			"pesanan: ongkir 1.50, total 7.50, seller tetap 5.52 (ongkir untuk Phsar)")
		e.call(t, "GET", "/v1/shop/checkouts/"+c.str("id"), nil, buyer).want(200).eq("orders.0.shipping_fee", "1.50")

		home := string(e.page(t, "/en/").Raw)
		check(t, strings.Contains(home, "Fast delivery in Sihanoukville"), "banner tidak boleh lagi menjanjikan gratis ongkir")
	})

	t.Run("kurs riel baru dipakai di website", func(t *testing.T) {
		e.call(t, "POST", "/admin/api/shipping/fx", map[string]string{"riel_per_usd": "4,100", "note": "NBC"}, fin).want(200).eq("riel_per_usd", "4100.0000").eq("rates.0.note", "NBC")
		check(t, rielRe.MatchString(string(e.page(t, "/en/").Raw)), "website harus memakai kurs terbaru")
		e.call(t, "POST", "/admin/api/shipping/fx", map[string]string{"riel_per_usd": "4050"}, fin).want(200).eq("riel_per_usd", "4050.0000")
	})
}
