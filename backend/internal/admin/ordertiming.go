package admin

// Aturan batas waktu pesanan (super admin): berapa lama seller boleh menunda "Terima" dan "Siap dijemput",
// dan apakah hitungannya hanya berjalan saat toko buka. Berlaku untuk pesanan baru.

import (
	"net/http"

	"github.com/jackc/pgx/v5"

	"phsar/internal/hours"
	"phsar/internal/platform/db"
)

// GET /admin/api/order-timing
func (h *Handler) getOrderTiming(w http.ResponseWriter, r *http.Request) {
	t, err := hours.LoadTiming(r.Context(), h.pool)
	if err != nil {
		h.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, t)
}

// PUT /admin/api/order-timing
func (h *Handler) saveOrderTiming(w http.ResponseWriter, r *http.Request) {
	var in hours.Timing
	if !decode(w, r, &in) {
		return
	}
	if err := in.Check(); err != nil {
		h.fail(w, err)
		return
	}
	err := db.WithActor(r.Context(), h.pool, actor(r), "admin", "", func(tx pgx.Tx) error {
		_, err := tx.Exec(r.Context(), `INSERT INTO platform_settings (key, value, updated_by) VALUES ('order_timing', $1, $2)
			ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_by = EXCLUDED.updated_by, updated_at = now()`, in, actor(r))
		return err
	})
	if err != nil {
		h.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, in)
}
