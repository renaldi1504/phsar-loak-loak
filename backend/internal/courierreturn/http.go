package courierreturn

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"phsar/internal/auth"
	"phsar/internal/storage"
)

type Handler struct {
	S     *Service
	Files storage.Storage
	Log   *slog.Logger
}

func uid(r *http.Request) string {
	c, _ := auth.UserFrom(r.Context())
	return c.Subject
}

// CourierRoutes: /v1/courier/returns (sudah lewat RequireUser + NoGuest). {id} = id tugas pengantaran.
func (h *Handler) CourierRoutes() http.Handler {
	r := chi.NewRouter()
	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		in, err := h.S.CourierInfo(r.Context(), uid(r), chi.URLParam(r, "id"))
		h.reply(w, r, in, err)
	})
	r.Post("/{id}/start", func(w http.ResponseWriter, r *http.Request) {
		in, err := h.S.CourierStart(r.Context(), uid(r), chi.URLParam(r, "id"))
		h.reply(w, r, in, err)
	})
	r.Post("/{id}/returned", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			ProofKey string `json:"proof_key"`
		}
		if r.ContentLength > 0 && !decode(w, r, &body) {
			return
		}
		in, err := h.S.CourierReturned(r.Context(), uid(r), chi.URLParam(r, "id"), body.ProofKey)
		h.reply(w, r, in, err)
	})
	return r
}

// AdminRoutes: /admin/api/courier-returns (izin delivery.dispatch dicek router admin).
func (h *Handler) AdminRoutes() http.Handler {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		l, c, err := h.S.AdminList(r.Context(), r.URL.Query().Get("status"))
		h.reply(w, r, map[string]any{"returns": l, "counts": c}, err)
	})
	// Mulai pengembalian dari pengantaran yang gagal.
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			DeliveryID string `json:"delivery_id"`
			CourierID  string `json:"courier_id"`
		}
		if !decode(w, r, &in) {
			return
		}
		x, err := h.S.AdminStart(r.Context(), uid(r), in.DeliveryID, in.CourierID)
		h.reply(w, r, x, err)
	})
	r.Post("/{id}/assign", func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			CourierID string `json:"courier_id"`
		}
		if !decode(w, r, &in) {
			return
		}
		x, err := h.S.AdminAssign(r.Context(), uid(r), chi.URLParam(r, "id"), in.CourierID)
		h.reply(w, r, x, err)
	})
	r.Post("/{id}/resolve", func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			Restock *bool  `json:"restock"`
			Note    string `json:"note"`
		}
		if !decode(w, r, &in) {
			return
		}
		if in.Restock == nil {
			h.reply(w, r, nil, fieldError{"restock", "choose whether the items go back to stock"})
			return
		}
		x, err := h.S.AdminResolve(r.Context(), uid(r), chi.URLParam(r, "id"), *in.Restock, in.Note)
		h.reply(w, r, x, err)
	})
	// Foto: ?which=returned (serah terima di toko, default) · fail (bukti gagal antar). Privat, tidak di-cache.
	r.Get("/{id}/proof", func(w http.ResponseWriter, r *http.Request) {
		key, err := h.S.ProofKey(r.Context(), chi.URLParam(r, "id"), r.URL.Query().Get("which"))
		if err != nil {
			h.reply(w, r, nil, err)
			return
		}
		f, err := h.Files.Open(r.Context(), key)
		if err != nil {
			h.reply(w, r, nil, ErrNotFound)
			return
		}
		defer f.Close()
		w.Header().Set("Cache-Control", "private, no-store")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		switch {
		case strings.HasSuffix(key, ".png"):
			w.Header().Set("Content-Type", "image/png")
		case strings.HasSuffix(key, ".webp"):
			w.Header().Set("Content-Type", "image/webp")
		default:
			w.Header().Set("Content-Type", "image/jpeg")
		}
		_, _ = io.Copy(w, f)
	})
	return r
}

// SellerRoutes: /v1/seller/stores/{storeID}/courier-returns (anggota toko: pemilik & staf).
func (h *Handler) SellerRoutes() http.Handler {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		l, err := h.S.SellerList(r.Context(), chi.URLParam(r, "storeID"), r.URL.Query().Get("view"))
		h.reply(w, r, map[string]any{"returns": l}, err)
	})
	r.Post("/{id}/confirm", func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			Condition string `json:"condition"`
			Note      string `json:"note"`
		}
		if !decode(w, r, &in) {
			return
		}
		x, err := h.S.SellerConfirm(r.Context(), chi.URLParam(r, "storeID"), uid(r), chi.URLParam(r, "id"), in.Condition, in.Note)
		h.reply(w, r, x, err)
	})
	return r
}

func (h *Handler) reply(w http.ResponseWriter, r *http.Request, v any, err error) {
	var fe fieldError
	switch {
	case err == nil:
		writeJSON(w, http.StatusOK, v)
	case errors.As(err, &fe):
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"error": map[string]string{"code": "validation", "field": fe.field, "message": fe.msg}})
	case errors.Is(err, ErrNotCourier):
		writeErr(w, http.StatusForbidden, "not_courier", "This account is not an active Phsar courier.")
	case errors.Is(err, ErrNotFound):
		writeErr(w, http.StatusNotFound, "not_found", "Return not found.")
	case errors.Is(err, ErrBadStatus):
		writeErr(w, http.StatusConflict, "bad_status", "This package was already updated. Refresh.")
	default:
		h.Log.Error("courierreturn", "path", r.URL.Path, "err", err)
		writeErr(w, http.StatusInternalServerError, "internal", "something went wrong")
	}
}

func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	dec := json.NewDecoder(io.LimitReader(r.Body, 64<<10))
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
