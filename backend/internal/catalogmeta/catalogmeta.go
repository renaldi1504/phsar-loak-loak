// Package catalogmeta menyajikan data referensi yang dibutuhkan app sejak awal:
// bahasa, wilayah layanan, ongkir, kurs tampilan, dan kategori.
//
// Query ditulis langsung dengan pgx di Tahap 0; mulai Tahap 1 query dipindah ke sqlc.
package catalogmeta

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"

	"phsar/internal/geo"
	"phsar/internal/i18n"
)

type Handler struct{ pool *pgxpool.Pool }

func New(pool *pgxpool.Pool) *Handler { return &Handler{pool: pool} }

type serviceArea struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	FeePerStore string `json:"fee_per_store"`
}

// GET /v1/app/meta
func (h *Handler) Meta(w http.ResponseWriter, r *http.Request) {
	lang := i18n.From(r.Context())
	rows, err := h.pool.Query(r.Context(), `
		SELECT sa.code, sa.name, sz.fee_per_store::text
		FROM service_areas sa JOIN shipping_zones sz ON sz.code = sa.shipping_zone_code
		WHERE sa.is_active ORDER BY sa.code`)
	if err != nil {
		writeErr(w, err)
		return
	}
	defer rows.Close()
	areas := []serviceArea{}
	for rows.Next() {
		var a serviceArea
		var name map[string]string
		if err := rows.Scan(&a.Code, &name, &a.FeePerStore); err != nil {
			writeErr(w, err)
			return
		}
		a.Name = i18n.Pick(name, lang)
		areas = append(areas, a)
	}
	var rate string
	err = h.pool.QueryRow(r.Context(),
		`SELECT riel_per_usd::text FROM exchange_rates ORDER BY rate_date DESC, id DESC LIMIT 1`).Scan(&rate)
	if err != nil {
		rate = "4000"
	}
	writeJSON(w, map[string]any{
		"languages":        i18n.Supported,
		"lang":             lang,
		"currency":         "USD",
		"riel_per_usd":     rate,
		"service_areas":    areas,
		"cash_on_delivery": false,
	})
}

// GET /v1/app/service-areas/check?lat=..&lng=..
// Dipakai checkout (alamat pembeli) dan Seller Center (titik penjemputan).
func (h *Handler) CheckPoint(w http.ResponseWriter, r *http.Request) {
	lat, err1 := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	lng, err2 := strconv.ParseFloat(r.URL.Query().Get("lng"), 64)
	if err1 != nil || err2 != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":{"code":"bad_point","message":"lat and lng are required"}}`))
		return
	}
	rows, err := h.pool.Query(r.Context(), `SELECT code, boundary FROM service_areas WHERE is_active`)
	if err != nil {
		writeErr(w, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var code string
		var raw []byte
		if err := rows.Scan(&code, &raw); err != nil {
			writeErr(w, err)
			return
		}
		poly, err := geo.ParsePolygon(raw)
		if err == nil && poly.Contains(lat, lng) {
			writeJSON(w, map[string]any{"served": true, "service_area": code})
			return
		}
	}
	msg := map[i18n.Lang]string{
		i18n.EN: "We currently deliver in Sihanoukville only.",
		i18n.KM: "បច្ចុប្បន្នយើងដឹកជញ្ជូនតែនៅក្រុងព្រះសីហនុប៉ុណ្ណោះ។",
		i18n.ZH: "目前我们只在西哈努克城配送。",
	}[i18n.From(r.Context())]
	writeJSON(w, map[string]any{"served": false, "message": msg})
}

type category struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	Slug     string     `json:"slug"`
	Icon     string     `json:"icon,omitempty"`
	Children []category `json:"children,omitempty"`
}

// GET /v1/app/categories
func (h *Handler) Categories(w http.ResponseWriter, r *http.Request) {
	lang := string(i18n.From(r.Context()))
	rows, err := h.pool.Query(r.Context(), `
		SELECT c.id, c.parent_id, COALESCE(ct.name, cen.name), COALESCE(ct.slug, cen.slug), COALESCE(c.icon, '')
		FROM categories c
		JOIN category_translations cen ON cen.category_id = c.id AND cen.lang = 'en'
		LEFT JOIN category_translations ct ON ct.category_id = c.id AND ct.lang = $1::lang_code
		WHERE c.deleted_at IS NULL
		ORDER BY c.sort_order, c.id`, lang)
	if err != nil {
		writeErr(w, err)
		return
	}
	defer rows.Close()
	byParent := map[string][]category{}
	for rows.Next() {
		var c category
		var parent *string
		if err := rows.Scan(&c.ID, &parent, &c.Name, &c.Slug, &c.Icon); err != nil {
			writeErr(w, err)
			return
		}
		key := ""
		if parent != nil {
			key = *parent
		}
		byParent[key] = append(byParent[key], c)
	}
	var build func(parent string) []category
	build = func(parent string) []category {
		list := byParent[parent]
		for i := range list {
			list[i].Children = build(list[i].ID)
		}
		return list
	}
	out := build("")
	if out == nil {
		out = []category{}
	}
	writeJSON(w, map[string]any{"categories": out})
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(map[string]any{"error": map[string]string{"code": "internal", "message": err.Error()}})
}
