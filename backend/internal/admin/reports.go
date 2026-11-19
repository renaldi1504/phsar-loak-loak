package admin

// Reports & strikes: antrean laporan pembeli, larang produk, strike toko (izin report.handle).

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"phsar/internal/reports"
)

func (h *Handler) mountReports(r chi.Router) {
	svc := &reports.Service{Pool: h.pool, Log: h.log}
	r.Route("/reports", func(p chi.Router) {
		p.Use(h.can("report.handle"))
		p.Get("/", func(w http.ResponseWriter, r *http.Request) {
			raw, err := svc.Groups(r.Context(), r.URL.Query().Get("status"))
			if err != nil {
				svc.Fail(w, err)
				return
			}
			n, _ := svc.OpenCount(r.Context())
			reports.WriteRaw(w, []byte(`{"groups":`+string(raw)+`,"open":`+strconv.Itoa(n)+`}`))
		})
		p.Post("/resolve", func(w http.ResponseWriter, r *http.Request) {
			var in reports.Resolution
			if !reports.Decode(w, r, &in) {
				return
			}
			if err := svc.Resolve(r.Context(), actor(r), in); err != nil {
				svc.Fail(w, err)
				return
			}
			reports.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
		})
		p.Get("/stores/{storeID}/strikes", func(w http.ResponseWriter, r *http.Request) {
			raw, err := svc.Strikes(r.Context(), chi.URLParam(r, "storeID"), false)
			if err != nil {
				svc.Fail(w, err)
				return
			}
			reports.WriteRaw(w, raw)
		})
		p.Post("/strikes", func(w http.ResponseWriter, r *http.Request) {
			var in reports.StrikeInput
			if !reports.Decode(w, r, &in) {
				return
			}
			if _, err := svc.AddStrike(r.Context(), actor(r), in); err != nil {
				svc.Fail(w, err)
				return
			}
			raw, _ := svc.Strikes(r.Context(), in.StoreID, false)
			reports.WriteRaw(w, raw)
		})
		p.Post("/strikes/{id}/revoke", func(w http.ResponseWriter, r *http.Request) {
			var in struct {
				Note string `json:"note"`
			}
			if !reports.Decode(w, r, &in) {
				return
			}
			storeID, err := svc.Revoke(r.Context(), actor(r), chi.URLParam(r, "id"), in.Note)
			if err != nil {
				svc.Fail(w, err)
				return
			}
			raw, _ := svc.Strikes(r.Context(), storeID, false)
			reports.WriteRaw(w, raw)
		})
	})
}
