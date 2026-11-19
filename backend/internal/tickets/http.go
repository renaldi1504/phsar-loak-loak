package tickets

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"phsar/internal/auth"
	"phsar/internal/i18n"
)

type Handler struct{ S *Service }

func uid(r *http.Request) string {
	c, _ := auth.UserFrom(r.Context())
	return c.Subject
}

// BuyerRoutes dipasang di dalam /v1/shop (sudah login; tamu juga boleh).
func (h *Handler) BuyerRoutes(r chi.Router) {
	scope := func(r *http.Request) Scope { return Scope{UserID: uid(r), Lang: string(i18n.From(r.Context()))} }
	h.userRoutes(r, "/tickets", scope)
}

// SellerRoutes dipasang di /v1/seller/stores/{storeID}/tickets (anggota toko).
func (h *Handler) SellerRoutes() http.Handler {
	r := chi.NewRouter()
	h.userRoutes(r, "", func(r *http.Request) Scope { return Scope{UserID: uid(r), StoreID: chi.URLParam(r, "storeID")} })
	return r
}

// userRoutes: prefix "/tickets" (di dalam /v1/shop) atau "" (router sendiri yang dipasang di …/tickets).
func (h *Handler) userRoutes(r chi.Router, prefix string, scope func(*http.Request) Scope) {
	root := prefix
	if root == "" {
		root = "/"
	}
	r.Get(root, func(w http.ResponseWriter, r *http.Request) {
		l, err := h.S.List(r.Context(), scope(r))
		reply(w, http.StatusOK, map[string]any{"tickets": l}, err)
	})
	r.Post(root, func(w http.ResponseWriter, r *http.Request) {
		var in CreateInput
		if !decode(w, r, &in) {
			return
		}
		t, err := h.S.Create(r.Context(), scope(r), in)
		reply(w, http.StatusCreated, t, err)
	})
	r.Get(prefix+"/{id}", func(w http.ResponseWriter, r *http.Request) {
		t, err := h.S.Get(r.Context(), scope(r), chi.URLParam(r, "id"))
		reply(w, http.StatusOK, t, err)
	})
	r.Post(prefix+"/{id}/messages", func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			Body string `json:"body"`
		}
		if !decode(w, r, &in) {
			return
		}
		t, err := h.S.Reply(r.Context(), scope(r), chi.URLParam(r, "id"), in.Body)
		reply(w, http.StatusOK, t, err)
	})
	r.Post(prefix+"/{id}/solve", func(w http.ResponseWriter, r *http.Request) {
		t, err := h.S.Solve(r.Context(), scope(r), chi.URLParam(r, "id"))
		reply(w, http.StatusOK, t, err)
	})
}

// Admin — dipasang paket admin (izin support.handle). adminID = admin yang sedang login.
func (h *Handler) AdminRoutes(r chi.Router, adminID func(*http.Request) string) {
	r.Get("/tickets", func(w http.ResponseWriter, r *http.Request) {
		l, c, err := h.S.AdminList(r.Context(), adminID(r), r.URL.Query().Get("filter"), r.URL.Query().Get("q"))
		reply(w, http.StatusOK, map[string]any{"tickets": l, "counts": c}, err)
	})
	r.Get("/tickets/{id}", func(w http.ResponseWriter, r *http.Request) {
		t, err := h.S.AdminGet(r.Context(), chi.URLParam(r, "id"))
		reply(w, http.StatusOK, t, err)
	})
	r.Post("/tickets/{id}/reply", func(w http.ResponseWriter, r *http.Request) {
		var in AdminReply
		if !decode(w, r, &in) {
			return
		}
		t, err := h.S.AdminReplyTo(r.Context(), adminID(r), chi.URLParam(r, "id"), in)
		reply(w, http.StatusOK, t, err)
	})
	r.Post("/tickets/{id}/assign", func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			ToMe bool `json:"to_me"`
		}
		if !decode(w, r, &in) {
			return
		}
		t, err := h.S.Assign(r.Context(), adminID(r), chi.URLParam(r, "id"), in.ToMe)
		reply(w, http.StatusOK, t, err)
	})
	r.Post("/tickets/{id}/status", func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			Status string `json:"status"`
		}
		if !decode(w, r, &in) {
			return
		}
		t, err := h.S.SetStatus(r.Context(), adminID(r), chi.URLParam(r, "id"), in.Status)
		reply(w, http.StatusOK, t, err)
	})
}

func reply(w http.ResponseWriter, status int, v any, err error) {
	var fe FieldError
	switch {
	case err == nil:
		writeJSON(w, status, v)
	case errors.As(err, &fe):
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"error": map[string]string{"code": "validation", "field": fe.Field, "message": fe.Msg}})
	case errors.Is(err, ErrNotFound):
		writeErr(w, http.StatusNotFound, "not_found", err.Error())
	case errors.Is(err, ErrRateLimited):
		writeErr(w, http.StatusTooManyRequests, "rate_limited", err.Error())
	default:
		slog.Error("tiket bantuan", "err", err)
		writeErr(w, http.StatusInternalServerError, "internal", "something went wrong")
	}
}

func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 32<<10))
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
