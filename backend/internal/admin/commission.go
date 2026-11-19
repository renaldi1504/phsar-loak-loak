package admin

// Aturan komisi: default (global), per kategori, per toko (mis. 0% tiga bulan pertama untuk seller awal).
// Mengubah tarif = aturan lama ditutup (valid_to = sekarang) + aturan baru dibuat → riwayat lengkap & tercatat di audit log.
// Berlaku untuk pesanan BARU; pesanan lama menyimpan tarifnya sendiri.

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"

	"phsar/internal/platform/db"
)

type commissionRule struct {
	ID         string     `json:"id"`
	Scope      string     `json:"scope"`
	CategoryID string     `json:"category_id,omitempty"`
	StoreID    string     `json:"store_id,omitempty"`
	StoreName  string     `json:"store_name,omitempty"`
	Rate       string     `json:"rate_percent"`
	ValidFrom  time.Time  `json:"valid_from"`
	ValidTo    *time.Time `json:"valid_to"`
	By         string     `json:"created_by"`
}

const ruleActive = `r.valid_from <= now() AND (r.valid_to IS NULL OR r.valid_to > now())`

// GET /admin/api/commission
func (h *Handler) commissionView(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	scan := func(q string, args ...any) ([]commissionRule, error) {
		rows, err := h.pool.Query(ctx, q, args...)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		out := []commissionRule{}
		for rows.Next() {
			var c commissionRule
			if err := rows.Scan(&c.ID, &c.Scope, &c.CategoryID, &c.StoreID, &c.StoreName, &c.Rate, &c.ValidFrom, &c.ValidTo, &c.By); err != nil {
				return nil, err
			}
			out = append(out, c)
		}
		return out, rows.Err()
	}
	const cols = `SELECT r.id::text, r.scope, coalesce(r.category_id, ''), coalesce(r.store_id::text, ''), coalesce(s.name, ''),
		r.rate_percent::text, r.valid_from, r.valid_to, coalesce(nullif(u.name, ''), u.phone, '')
		FROM commission_rules r LEFT JOIN stores s ON s.id = r.store_id LEFT JOIN users u ON u.id = r.created_by `
	active, err := scan(cols + `WHERE ` + ruleActive + ` ORDER BY r.scope, s.name, r.category_id`)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "internal", err.Error())
		return
	}
	history, err := scan(cols + `WHERE r.scope = 'global' ORDER BY r.valid_from DESC LIMIT 10`)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "internal", err.Error())
		return
	}
	type cat struct {
		ID    string            `json:"id"`
		Names map[string]string `json:"names"`
	}
	cats := []cat{}
	rows, err := h.pool.Query(ctx, `SELECT c.id, jsonb_object_agg(t.lang, t.name) FROM categories c JOIN category_translations t ON t.category_id = c.id
		WHERE c.parent_id IS NULL GROUP BY c.id, c.sort_order ORDER BY c.sort_order`)
	if err == nil {
		for rows.Next() {
			var c cat
			if rows.Scan(&c.ID, &c.Names) == nil {
				cats = append(cats, c)
			}
		}
		rows.Close()
	}
	writeJSON(w, http.StatusOK, map[string]any{"active": active, "global_history": history, "categories": cats})
}

// GET /admin/api/commission/stores?q= — cari toko untuk aturan khusus.
func (h *Handler) commissionStores(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	rows, err := h.pool.Query(r.Context(), `SELECT id::text, name, slug, status::text FROM stores
		WHERE deleted_at IS NULL AND ($1 = '' OR name ILIKE '%' || $1 || '%' OR slug ILIKE '%' || $1 || '%') ORDER BY name LIMIT 20`, q)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "internal", err.Error())
		return
	}
	defer rows.Close()
	type st struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Slug   string `json:"slug"`
		Status string `json:"status"`
	}
	out := []st{}
	for rows.Next() {
		var s st
		if rows.Scan(&s.ID, &s.Name, &s.Slug, &s.Status) == nil {
			out = append(out, s)
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{"stores": out})
}

// POST /admin/api/commission/rules {scope, category_id?, store_id?, rate_percent, valid_to?}
func (h *Handler) commissionSet(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Scope      string     `json:"scope"`
		CategoryID string     `json:"category_id"`
		StoreID    string     `json:"store_id"`
		Rate       string     `json:"rate_percent"`
		ValidTo    *time.Time `json:"valid_to"`
	}
	if !decode(w, r, &in) {
		return
	}
	rate, err := decimal.NewFromString(strings.TrimSpace(in.Rate))
	if err != nil || rate.IsNegative() || rate.GreaterThan(decimal.NewFromInt(50)) || !rate.Equal(rate.Round(2)) {
		writeFieldErr(w, "rate_percent")
		return
	}
	if in.ValidTo != nil && !in.ValidTo.After(time.Now().Add(time.Hour)) {
		writeFieldErr(w, "valid_to")
		return
	}
	var match string
	var args []any
	switch in.Scope {
	case "global":
		if in.ValidTo != nil {
			writeFieldErr(w, "valid_to") // tarif default tidak punya tanggal akhir
			return
		}
		match = `scope = 'global'`
	case "category":
		match, args = `scope = 'category' AND category_id = $1`, []any{in.CategoryID}
	case "store":
		if !uuidRe.MatchString(in.StoreID) {
			writeFieldErr(w, "store_id")
			return
		}
		match, args = `scope = 'store' AND store_id = $1::uuid`, []any{in.StoreID}
	default:
		writeFieldErr(w, "scope")
		return
	}
	err = db.WithActor(r.Context(), h.pool, actor(r), "admin", "", func(tx pgx.Tx) error {
		if in.Scope == "category" {
			var ok bool
			if err := tx.QueryRow(r.Context(), `SELECT EXISTS (SELECT 1 FROM categories WHERE id = $1)`, in.CategoryID).Scan(&ok); err != nil || !ok {
				return errBadTarget
			}
		}
		if in.Scope == "store" {
			var ok bool
			if err := tx.QueryRow(r.Context(), `SELECT EXISTS (SELECT 1 FROM stores WHERE id = $1 AND deleted_at IS NULL)`, in.StoreID).Scan(&ok); err != nil || !ok {
				return errBadTarget
			}
		}
		if _, err := tx.Exec(r.Context(), `UPDATE commission_rules SET valid_to = now()
			WHERE `+match+` AND valid_from <= now() AND (valid_to IS NULL OR valid_to > now())`, args...); err != nil {
			return err
		}
		var cat, store *string
		if in.Scope == "category" {
			cat = &in.CategoryID
		}
		if in.Scope == "store" {
			store = &in.StoreID
		}
		_, err := tx.Exec(r.Context(), `INSERT INTO commission_rules (scope, category_id, store_id, rate_percent, valid_from, valid_to, created_by)
			VALUES ($1, $2, $3::uuid, $4, now() + interval '1 millisecond', $5, $6)`, in.Scope, cat, store, rate, in.ValidTo, actor(r))
		return err
	})
	if errors.Is(err, errBadTarget) {
		writeFieldErr(w, map[string]string{"category": "category_id", "store": "store_id"}[in.Scope])
		return
	}
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "internal", err.Error())
		return
	}
	h.commissionView(w, r)
}

var errBadTarget = errors.New("target tidak ditemukan")

// DELETE /admin/api/commission/rules/{id} — akhiri aturan kategori/toko (kembali ke tarif di atasnya). Default tidak bisa dihapus.
func (h *Handler) commissionEnd(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !uuidRe.MatchString(id) {
		writeErr(w, http.StatusNotFound, "not_found", "rule not found")
		return
	}
	var n int64
	err := db.WithActor(r.Context(), h.pool, actor(r), "admin", "", func(tx pgx.Tx) error {
		tag, err := tx.Exec(r.Context(), `UPDATE commission_rules r SET valid_to = now() WHERE r.id = $1 AND r.scope <> 'global' AND `+ruleActive, id)
		n = tag.RowsAffected()
		return err
	})
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "internal", err.Error())
		return
	}
	if n == 0 {
		writeErr(w, http.StatusNotFound, "not_found", "rule not found or already ended")
		return
	}
	h.commissionView(w, r)
}
