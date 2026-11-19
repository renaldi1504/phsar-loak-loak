package seller

// Jam buka toko: dipakai menghitung batas waktu pesanan kalau super admin menyalakan "hitung hanya saat toko buka".

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"phsar/internal/hours"
)

type hoursView struct {
	Hours         *hours.Week `json:"hours"` // null = memakai jam standar Phsar
	DefaultHours  hours.Week  `json:"default_hours"`
	BusinessHours bool        `json:"business_hours"`
	AcceptMinutes int         `json:"accept_minutes"`
	PackEnabled   bool        `json:"pack_enabled"`
	PackMinutes   int         `json:"pack_minutes"`
	OpenNow       bool        `json:"open_now"`
	OpensAt       *time.Time  `json:"opens_at,omitempty"`
}

func (s *Service) hoursView(r *http.Request, storeID string) (hoursView, error) {
	ctx := r.Context()
	t, err := hours.LoadTiming(ctx, s.pool)
	if err != nil {
		return hoursView{}, err
	}
	w, err := hours.StoreWeek(ctx, s.pool, storeID)
	if err != nil {
		return hoursView{}, err
	}
	v := hoursView{Hours: w, DefaultHours: t.DefaultHours, BusinessHours: t.BusinessHours, AcceptMinutes: t.AcceptMinutes, PackEnabled: t.PackEnabled, PackMinutes: t.PackMinutes}
	v.OpenNow, v.OpensAt = t.Effective(w).OpenAt(time.Now())
	if !t.BusinessHours {
		v.OpenNow, v.OpensAt = true, nil
	}
	return v, nil
}

// GET /v1/seller/stores/{storeID}/hours
func (h *Handler) getHours(w http.ResponseWriter, r *http.Request) {
	v, err := h.svc.hoursView(r, chi.URLParam(r, "storeID"))
	h.respond(w, http.StatusOK, v, err)
}

// PUT /v1/seller/stores/{storeID}/hours {hours: [7]|null} — null = kembali ke jam standar.
func (h *Handler) putHours(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Hours *hours.Week `json:"hours"`
	}
	if !decode(w, r, &in) {
		return
	}
	var raw []byte
	if in.Hours != nil {
		if err := in.Hours.Validate("hours"); err != nil {
			h.fail(w, err)
			return
		}
		raw, _ = json.Marshal(in.Hours)
	}
	id := chi.URLParam(r, "storeID")
	if _, err := h.svc.pool.Exec(r.Context(), `UPDATE stores SET opening_hours = $2 WHERE id = $1`, id, raw); err != nil {
		h.fail(w, err)
		return
	}
	v, err := h.svc.hoursView(r, id)
	h.respond(w, http.StatusOK, v, err)
}
