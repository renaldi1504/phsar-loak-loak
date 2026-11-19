package admin

// pricehistory: aturan flash sale (super admin, izin "flash.rules" → hanya "*"):
// harga flash harus di bawah harga normal terendah 30 hari sebelum sesi mulai. Default nyala.

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/jackc/pgx/v5"

	"phsar/internal/platform/db"
	"phsar/internal/pricehistory"
)

// GET /admin/api/flash-rules
func (h *Handler) phGetRules(w http.ResponseWriter, r *http.Request) {
	rules, err := pricehistory.LoadRules(r.Context(), h.pool)
	if err != nil {
		h.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, rules)
}

// PUT /admin/api/flash-rules {"low30_enabled": bool}
func (h *Handler) phSaveRules(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Low30Enabled *bool `json:"low30_enabled"`
	}
	dec := json.NewDecoder(io.LimitReader(r.Body, 4<<10))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&in); err != nil {
		writeErr(w, http.StatusBadRequest, "bad_json", "invalid JSON body")
		return
	}
	if in.Low30Enabled == nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"error": map[string]string{"code": "validation", "field": "low30_enabled", "message": "required"}})
		return
	}
	rules := pricehistory.Rules{Low30Enabled: *in.Low30Enabled}
	err := db.WithActor(r.Context(), h.pool, actor(r), "admin", "", func(tx pgx.Tx) error {
		_, err := tx.Exec(r.Context(), `INSERT INTO platform_settings (key, value, updated_by, updated_at) VALUES ('flash_rules', $1, $2, now())
			ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_by = EXCLUDED.updated_by, updated_at = now()`, rules, actor(r))
		return err
	})
	if err != nil {
		h.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, rules)
}
