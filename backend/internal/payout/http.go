package payout

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"phsar/internal/auth"
	"phsar/internal/paging"
)

type Handler struct {
	S   *Service
	Log *slog.Logger
}

func uid(r *http.Request) string {
	c, _ := auth.UserFrom(r.Context())
	return c.Subject
}

// SellerRoutes: /v1/seller/stores/{storeID}/payouts (khusus pemilik toko; dicek router seller).
func (h *Handler) SellerRoutes() http.Handler {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		s, err := h.S.StoreSummary(r.Context(), chi.URLParam(r, "storeID"))
		h.reply(w, r, s, err)
	})
	r.Get("/ledger", func(w http.ResponseWriter, r *http.Request) { // ?page=&size=
		l, m, err := h.S.Ledger(r.Context(), chi.URLParam(r, "storeID"), paging.FromRequest(r, 20, 100))
		h.reply(w, r, map[string]any{"ledger": l, "page": m}, err)
	})
	r.Put("/account", func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			Bank          string `json:"bank"`
			AccountName   string `json:"account_name"`
			AccountNumber string `json:"account_number"`
		}
		if !decode(w, r, &in) {
			return
		}
		a, err := h.S.SaveAccount(r.Context(), chi.URLParam(r, "storeID"), uid(r), in.Bank, in.AccountName, in.AccountNumber)
		h.reply(w, r, a, err)
	})
	r.Get("/withdraw", h.withdrawState) // extras: tarik saldo manual
	r.Post("/withdraw", h.withdraw)     // extras:
	return r
}

// Admin (izin payout.approve dicek router admin).
func (h *Handler) AdminGetSettings(w http.ResponseWriter, r *http.Request) {
	st, err := h.S.Settings(r.Context())
	h.reply(w, r, map[string]any{"settings": st, "next_run": st.NextRun(timeNow())}, err)
}

func (h *Handler) AdminSaveSettings(w http.ResponseWriter, r *http.Request) {
	var in Settings
	if !decode(w, r, &in) {
		return
	}
	st, err := h.S.SaveSettings(r.Context(), uid(r), in)
	h.reply(w, r, map[string]any{"settings": st, "next_run": st.NextRun(timeNow())}, err)
}

func (h *Handler) AdminList(w http.ResponseWriter, r *http.Request) {
	l, err := h.S.AdminList(r.Context(), uid(r), r.URL.Query().Get("status"))
	h.reply(w, r, map[string]any{"payouts": l}, err)
}

func (h *Handler) AdminGenerate(w http.ResponseWriter, r *http.Request) {
	if _, err := h.S.CompleteDue(r.Context()); err != nil {
		h.reply(w, r, nil, err)
		return
	}
	res, err := h.S.Generate(r.Context(), uid(r))
	h.reply(w, r, res, err)
}

func (h *Handler) AdminPaid(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Reference string `json:"reference"`
		Note      string `json:"note"`
	}
	if !decode(w, r, &in) {
		return
	}
	err := h.S.MarkPaid(r.Context(), uid(r), chi.URLParam(r, "id"), in.Reference, in.Note)
	h.reply(w, r, map[string]string{"result": "ok"}, err)
}

func (h *Handler) AdminCancel(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Note string `json:"note"`
	}
	if !decode(w, r, &in) {
		return
	}
	err := h.S.Cancel(r.Context(), uid(r), chi.URLParam(r, "id"), in.Note)
	h.reply(w, r, map[string]string{"result": "ok"}, err)
}

func (h *Handler) reply(w http.ResponseWriter, r *http.Request, v any, err error) {
	var fe fieldError
	switch {
	case err == nil:
		writeJSON(w, http.StatusOK, v)
	case errors.As(err, &fe):
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"error": map[string]string{"code": "validation", "field": fe.field, "message": fe.msg}})
	case errors.Is(err, ErrNotFound):
		writeErr(w, http.StatusNotFound, "not_found", "not found")
	case errors.Is(err, ErrBadStatus):
		writeErr(w, http.StatusConflict, "bad_status", "This payout was already processed. Refresh.")
	default:
		h.Log.Error("payout", "path", r.URL.Path, "err", err)
		writeErr(w, http.StatusInternalServerError, "internal", "something went wrong")
	}
}

func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	dec := json.NewDecoder(io.LimitReader(r.Body, 32<<10))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		writeErr(w, http.StatusBadRequest, "bad_json", err.Error())
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, code, msg string) {
	writeJSON(w, status, map[string]any{"error": map[string]string{"code": code, "message": msg}})
}
