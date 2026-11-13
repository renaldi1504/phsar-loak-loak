package seller

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"phsar/internal/auth"
	"phsar/internal/storage"
)

type Handler struct {
	svc      *Service
	uploader storage.Uploader
	log      *slog.Logger
}

func NewHandler(svc *Service, up storage.Uploader, log *slog.Logger) *Handler {
	return &Handler{svc: svc, uploader: up, log: log}
}

// Routes dipasang di /v1/seller (sudah melewati auth.RequireUser).
func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Post("/uploads", h.upload)
	r.Get("/stores", h.myStores)
	r.Post("/stores", h.createStore)
	r.Route("/stores/{storeID}", func(r chi.Router) {
		r.Use(h.member)
		r.Get("/", h.getStore)
		r.With(ownerOnly).Patch("/", h.updateStore)
		r.With(ownerOnly).Post("/submit", h.submit)
		r.With(ownerOnly).Get("/kyc", h.getKYC)
		r.With(ownerOnly).Put("/kyc", h.putKYC)
		r.Get("/pickup-points", h.listPickup)
		r.With(ownerOnly).Post("/pickup-points", h.createPickup)
		r.With(ownerOnly).Put("/pickup-points/{pointID}", h.updatePickup)
		r.With(ownerOnly).Delete("/pickup-points/{pointID}", h.deletePickup)
	})
	return r
}

type roleKey struct{}

func (h *Handler) member(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, _ := auth.UserFrom(r.Context())
		role, err := h.svc.Role(r.Context(), chi.URLParam(r, "storeID"), c.Subject)
		if err != nil {
			h.fail(w, err)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), roleKey{}, role)))
	})
}

func ownerOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if role, _ := r.Context().Value(roleKey{}).(string); role != "owner" {
			writeErr(w, http.StatusForbidden, "owner_only", ErrOwnerOnly.Error())
			return
		}
		next.ServeHTTP(w, r)
	})
}

func uid(r *http.Request) string {
	c, _ := auth.UserFrom(r.Context())
	return c.Subject
}

// POST /v1/seller/uploads  (multipart: kind, file)
func (h *Handler) upload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 11<<20)
	mr, err := r.MultipartReader()
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad_upload", "kirim multipart/form-data dengan field kind lalu file")
		return
	}
	kind := ""
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			writeErr(w, http.StatusBadRequest, "bad_upload", err.Error())
			return
		}
		switch part.FormName() {
		case "kind":
			b, _ := io.ReadAll(io.LimitReader(part, 64))
			kind = string(b)
		case "file":
			if kind == "" {
				writeErr(w, http.StatusBadRequest, "bad_upload", "field kind harus dikirim sebelum file")
				return
			}
			up, err := h.uploader.Save(r.Context(), uid(r), kind, part)
			if err != nil {
				h.fail(w, err)
				return
			}
			writeJSON(w, http.StatusCreated, up)
			return
		}
	}
	writeErr(w, http.StatusBadRequest, "bad_upload", "field file tidak ada")
}

func (h *Handler) myStores(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.MyStores(r.Context(), uid(r))
	h.respond(w, http.StatusOK, map[string]any{"stores": list}, err)
}

func (h *Handler) createStore(w http.ResponseWriter, r *http.Request) {
	var in CreateStoreInput
	if !decode(w, r, &in) {
		return
	}
	st, err := h.svc.CreateStore(r.Context(), uid(r), in)
	h.respond(w, http.StatusCreated, st, err)
}

func (h *Handler) getStore(w http.ResponseWriter, r *http.Request) {
	st, err := h.svc.GetStore(r.Context(), chi.URLParam(r, "storeID"))
	st.Role, _ = r.Context().Value(roleKey{}).(string)
	h.respond(w, http.StatusOK, st, err)
}

func (h *Handler) updateStore(w http.ResponseWriter, r *http.Request) {
	var in UpdateStoreInput
	if !decode(w, r, &in) {
		return
	}
	st, err := h.svc.UpdateStore(r.Context(), chi.URLParam(r, "storeID"), uid(r), in)
	h.respond(w, http.StatusOK, st, err)
}

func (h *Handler) submit(w http.ResponseWriter, r *http.Request) {
	st, err := h.svc.Submit(r.Context(), chi.URLParam(r, "storeID"), uid(r))
	if errors.Is(err, ErrIncomplete) {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{
			"error": map[string]string{"code": "store_incomplete", "message": err.Error()}, "checklist": st.Checklist})
		return
	}
	h.respond(w, http.StatusOK, st, err)
}

func (h *Handler) getKYC(w http.ResponseWriter, r *http.Request) {
	k, err := h.svc.GetKYC(r.Context(), chi.URLParam(r, "storeID"))
	h.respond(w, http.StatusOK, map[string]any{"kyc": k}, err)
}

func (h *Handler) putKYC(w http.ResponseWriter, r *http.Request) {
	var in KYCInput
	if !decode(w, r, &in) {
		return
	}
	k, err := h.svc.PutKYC(r.Context(), chi.URLParam(r, "storeID"), uid(r), in)
	h.respond(w, http.StatusOK, map[string]any{"kyc": k}, err)
}

func (h *Handler) listPickup(w http.ResponseWriter, r *http.Request) {
	l, err := h.svc.ListPickupPoints(r.Context(), chi.URLParam(r, "storeID"))
	h.respond(w, http.StatusOK, map[string]any{"pickup_points": l}, err)
}

func (h *Handler) createPickup(w http.ResponseWriter, r *http.Request) {
	var in PickupInput
	if !decode(w, r, &in) {
		return
	}
	l, err := h.svc.CreatePickupPoint(r.Context(), chi.URLParam(r, "storeID"), uid(r), in)
	h.respond(w, http.StatusCreated, map[string]any{"pickup_points": l}, err)
}

func (h *Handler) updatePickup(w http.ResponseWriter, r *http.Request) {
	var in PickupInput
	if !decode(w, r, &in) {
		return
	}
	l, err := h.svc.UpdatePickupPoint(r.Context(), chi.URLParam(r, "storeID"), chi.URLParam(r, "pointID"), uid(r), in)
	h.respond(w, http.StatusOK, map[string]any{"pickup_points": l}, err)
}

func (h *Handler) deletePickup(w http.ResponseWriter, r *http.Request) {
	l, err := h.svc.DeletePickupPoint(r.Context(), chi.URLParam(r, "storeID"), chi.URLParam(r, "pointID"), uid(r))
	h.respond(w, http.StatusOK, map[string]any{"pickup_points": l}, err)
}

func (h *Handler) respond(w http.ResponseWriter, status int, v any, err error) {
	if err != nil {
		h.fail(w, err)
		return
	}
	writeJSON(w, status, v)
}

func (h *Handler) fail(w http.ResponseWriter, err error) {
	var ve validationError
	switch {
	case errors.As(err, &ve):
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"error": map[string]string{"code": "validation", "field": ve.field, "message": ve.msg}})
	case errors.Is(err, ErrNotMember):
		writeErr(w, http.StatusNotFound, "store_not_found", err.Error()) // jangan bocorkan toko orang lain ada
	case errors.Is(err, ErrNotFound):
		writeErr(w, http.StatusNotFound, "not_found", err.Error())
	case errors.Is(err, ErrOutsideArea):
		writeErr(w, http.StatusUnprocessableEntity, "outside_service_area", err.Error())
	case errors.Is(err, ErrNotEditable):
		writeErr(w, http.StatusConflict, "store_under_review", err.Error())
	case errors.Is(err, ErrSlugTaken):
		writeErr(w, http.StatusConflict, "slug_taken", err.Error())
	case errors.Is(err, ErrTooManyStores):
		writeErr(w, http.StatusConflict, "too_many_stores", err.Error())
	case errors.Is(err, storage.ErrNotOwner):
		writeErr(w, http.StatusUnprocessableEntity, "upload_not_found", err.Error())
	case errors.Is(err, storage.ErrKind), errors.Is(err, storage.ErrType):
		writeErr(w, http.StatusUnsupportedMediaType, "upload_type", err.Error())
	case errors.Is(err, storage.ErrTooLarge):
		writeErr(w, http.StatusRequestEntityTooLarge, "upload_too_large", err.Error())
	default:
		h.log.Error("seller error", "err", err)
		writeErr(w, http.StatusInternalServerError, "internal", "internal error")
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
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, code, msg string) {
	writeJSON(w, status, map[string]any{"error": map[string]string{"code": code, "message": msg}})
}
