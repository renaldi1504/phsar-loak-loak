package invoices

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"phsar/internal/auth"
)

type Handler struct {
	S   *Service
	Log *slog.Logger
}

func uid(r *http.Request) string {
	c, _ := auth.UserFrom(r.Context())
	return c.Subject
}

// SellerRoutes: /v1/seller/stores/{storeID}/payouts/invoices (pemilik toko).
func (h *Handler) SellerRoutes() http.Handler {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		l, err := h.S.StoreList(r.Context(), chi.URLParam(r, "storeID"))
		h.reply(w, map[string]any{"invoices": l}, err)
	})
	r.Get("/{id}/html", func(w http.ResponseWriter, r *http.Request) {
		v, err := h.S.StoreGet(r.Context(), chi.URLParam(r, "storeID"), chi.URLParam(r, "id"))
		h.html(w, v, err)
	})
	return r
}

// ---- Admin (dipasang paket admin dengan izin commission.manage) ----

func (h *Handler) AdminList(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	l, err := h.S.List(r.Context(), q.Get("period"), q.Get("status"), q.Get("q"))
	if err != nil {
		h.reply(w, nil, err)
		return
	}
	st, err := h.S.Settings(r.Context())
	n := time.Now().In(phnomPenh)
	last := time.Date(n.Year(), n.Month(), 1, 0, 0, 0, 0, time.UTC).AddDate(0, -1, 0).Format("2006-01")
	h.reply(w, map[string]any{"invoices": l, "settings": st, "last_month": last}, err)
}

func (h *Handler) AdminSaveSettings(w http.ResponseWriter, r *http.Request) {
	var in Settings
	if !decode(w, r, &in) {
		return
	}
	st, err := h.S.SaveSettings(r.Context(), uid(r), in)
	h.reply(w, st, err)
}

func (h *Handler) AdminGenerate(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Period string `json:"period"`
	}
	if !decode(w, r, &in) {
		return
	}
	res, err := h.S.Generate(r.Context(), uid(r), in.Period, time.Now())
	h.reply(w, res, err)
}

func (h *Handler) AdminGet(w http.ResponseWriter, r *http.Request) {
	v, err := h.S.Get(r.Context(), chi.URLParam(r, "id"))
	h.reply(w, v, err)
}

func (h *Handler) AdminHTML(w http.ResponseWriter, r *http.Request) {
	v, err := h.S.Get(r.Context(), chi.URLParam(r, "id"))
	h.html(w, v, err)
}

func (h *Handler) AdminVoid(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Note string `json:"note"`
	}
	if !decode(w, r, &in) {
		return
	}
	v, err := h.S.Void(r.Context(), uid(r), chi.URLParam(r, "id"), in.Note)
	h.reply(w, v, err)
}

// ---------------------------------------------------------------------------

func (h *Handler) html(w http.ResponseWriter, v *Invoice, err error) {
	if err != nil {
		h.reply(w, nil, err)
		return
	}
	b, err := HTML(v)
	if err != nil {
		h.reply(w, nil, err)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "private, no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	_, _ = w.Write(b)
}

func (h *Handler) reply(w http.ResponseWriter, v any, err error) {
	var fe FieldError
	switch {
	case err == nil:
		writeJSON(w, http.StatusOK, v)
	case errors.As(err, &fe):
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"error": map[string]string{"code": "validation", "field": fe.Field, "message": fe.Msg}})
	case errors.Is(err, ErrNotFound):
		writeErr(w, http.StatusNotFound, "not_found", err.Error())
	case errors.Is(err, ErrState):
		writeErr(w, http.StatusConflict, "invoice_void", err.Error())
	default:
		if h.Log != nil {
			h.Log.Error("faktur", "err", err)
		}
		writeErr(w, http.StatusInternalServerError, "internal", "something went wrong")
	}
}

func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64<<10))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		writeErr(w, http.StatusBadRequest, "bad_json", "invalid JSON body")
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, code, msg string) {
	writeJSON(w, status, map[string]any{"error": map[string]string{"code": code, "message": msg}})
}
