package admin

// Shipping & FX rates: ongkir per zona, wilayah layanan (poligon), biaya kurir, kurs riel (izin commission.manage).

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"phsar/internal/shipping"
)

func (h *Handler) mountShipping(r chi.Router) {
	r.Route("/shipping", func(p chi.Router) {
		p.Use(h.can("commission.manage"))
		p.Get("/", h.shippingView)
		p.Post("/zones", h.shippingSaveZone(true))
		p.Put("/zones/{code}", h.shippingSaveZone(false))
		p.Post("/areas", h.shippingSaveArea(true))
		p.Put("/areas/{code}", h.shippingSaveArea(false))
		p.Post("/fx", h.shippingAddRate)
	})
}

func (h *Handler) shipStore() *shipping.Store { return &shipping.Store{Pool: h.pool} }

func (h *Handler) shippingView(w http.ResponseWriter, r *http.Request) {
	v, err := h.shipStore().View(r.Context())
	if err != nil {
		h.log.Error("ongkir & kurs", "err", err)
		writeErr(w, http.StatusInternalServerError, "internal", "could not load")
		return
	}
	writeJSON(w, http.StatusOK, v)
}

func (h *Handler) shippingSaveZone(create bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var in shipping.ZoneInput
		if !decode(w, r, &in) {
			return
		}
		if !create {
			in.Code = chi.URLParam(r, "code")
		}
		h.shippingReply(w, r, h.shipStore().SaveZone(r.Context(), actor(r), in, create))
	}
}

func (h *Handler) shippingSaveArea(create bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var in shipping.AreaInput
		if !decode(w, r, &in) {
			return
		}
		if !create {
			in.Code = chi.URLParam(r, "code")
		}
		h.shippingReply(w, r, h.shipStore().SaveArea(r.Context(), actor(r), in, create))
	}
}

func (h *Handler) shippingAddRate(w http.ResponseWriter, r *http.Request) {
	var in struct {
		RielPerUSD string `json:"riel_per_usd"`
		Note       string `json:"note"`
	}
	if !decode(w, r, &in) {
		return
	}
	h.shippingReply(w, r, h.shipStore().AddRate(r.Context(), actor(r), in.RielPerUSD, in.Note))
}

func (h *Handler) shippingReply(w http.ResponseWriter, r *http.Request, err error) {
	var fe shipping.FieldError
	switch {
	case err == nil:
		h.shippingView(w, r)
	case errors.As(err, &fe):
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"error": map[string]string{"code": "validation", "field": fe.Field, "message": fe.Msg}})
	case errors.Is(err, shipping.ErrNotFound):
		writeErr(w, http.StatusNotFound, "not_found", err.Error())
	case errors.Is(err, shipping.ErrExists):
		writeJSON(w, http.StatusConflict, map[string]any{"error": map[string]string{"code": "exists", "field": "code", "message": err.Error()}})
	case errors.Is(err, shipping.ErrLastArea):
		writeErr(w, http.StatusConflict, "last_area", err.Error())
	default:
		h.log.Error("ongkir & kurs", "err", err)
		writeErr(w, http.StatusInternalServerError, "internal", "could not save")
	}
}
