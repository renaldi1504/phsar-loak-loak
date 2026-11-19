package delivery

import (
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
	S   *Service
	Log *slog.Logger
}

func uid(r *http.Request) string {
	c, _ := auth.UserFrom(r.Context())
	return c.Subject
}

// CourierRoutes: /v1/courier (sudah lewat RequireUser).
func (h *Handler) CourierRoutes() http.Handler {
	r := chi.NewRouter()
	r.Get("/me", func(w http.ResponseWriter, r *http.Request) {
		c, err := h.S.Me(r.Context(), uid(r))
		h.reply(w, r, c, err)
	})
	r.Post("/me/online", func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			Online bool `json:"online"`
		}
		if !decode(w, r, &in) {
			return
		}
		c, err := h.S.SetOnline(r.Context(), uid(r), in.Online)
		h.reply(w, r, c, err)
	})
	r.Get("/tasks", func(w http.ResponseWriter, r *http.Request) {
		l, err := h.S.CourierTasks(r.Context(), uid(r), r.URL.Query().Get("view"))
		h.reply(w, r, map[string]any{"tasks": l}, err)
	})
	r.Get("/tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		t, err := h.S.CourierTask(r.Context(), uid(r), chi.URLParam(r, "id"))
		h.reply(w, r, t, err)
	})
	r.Post("/tasks/{id}/claim", func(w http.ResponseWriter, r *http.Request) {
		t, err := h.S.Claim(r.Context(), uid(r), chi.URLParam(r, "id"))
		h.reply(w, r, t, err)
	})
	r.Post("/tasks/{id}/release", func(w http.ResponseWriter, r *http.Request) {
		err := h.S.Release(r.Context(), uid(r), chi.URLParam(r, "id"))
		h.reply(w, r, map[string]string{"result": "ok"}, err)
	})
	r.Post("/tasks/{id}/pickup", func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			Code string `json:"code"`
		}
		if !decode(w, r, &in) {
			return
		}
		t, err := h.S.PickUp(r.Context(), uid(r), chi.URLParam(r, "id"), in.Code)
		h.reply(w, r, t, err)
	})
	r.Post("/tasks/{id}/deliver", func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			ProofKey string `json:"proof_key"`
		}
		if r.ContentLength > 0 && !decode(w, r, &in) {
			return
		}
		t, err := h.S.Deliver(r.Context(), uid(r), chi.URLParam(r, "id"), in.ProofKey)
		h.reply(w, r, t, err)
	})
	r.Post("/tasks/{id}/fail", func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			Reason   string `json:"reason"`
			Note     string `json:"note"`
			ProofKey string `json:"proof_key"` // courierreturn: foto bukti opsional
		}
		if !decode(w, r, &in) {
			return
		}
		t, err := h.S.FailWithProof(r.Context(), uid(r), chi.URLParam(r, "id"), in.Reason, in.Note, in.ProofKey)
		h.reply(w, r, t, err)
	})
	r.Post("/tasks/{id}/retry", func(w http.ResponseWriter, r *http.Request) {
		t, err := h.S.Retry(r.Context(), uid(r), chi.URLParam(r, "id"))
		h.reply(w, r, t, err)
	})
	// Foto bukti serah terima: multipart field "file".
	r.Post("/uploads", func(w http.ResponseWriter, r *http.Request) {
		if _, err := h.S.Me(r.Context(), uid(r)); err != nil {
			h.reply(w, r, nil, err)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, 9<<20)
		mr, err := r.MultipartReader()
		if err != nil {
			writeErr(w, http.StatusBadRequest, "bad_upload", "send multipart/form-data with a file field")
			return
		}
		for {
			part, err := mr.NextPart()
			if err != nil {
				writeErr(w, http.StatusBadRequest, "bad_upload", "file field missing")
				return
			}
			if part.FormName() != "file" {
				continue
			}
			up, err := h.S.Uploader.Save(r.Context(), uid(r), "delivery_proof", part)
			if err != nil {
				switch {
				case errors.Is(err, storage.ErrTooLarge):
					writeErr(w, http.StatusRequestEntityTooLarge, "too_large", "photo is too large (max 8 MB)")
				case errors.Is(err, storage.ErrType):
					writeErr(w, http.StatusUnsupportedMediaType, "bad_type", "upload a JPG, PNG or WebP photo")
				default:
					h.reply(w, r, nil, err)
				}
				return
			}
			writeJSON(w, http.StatusCreated, map[string]string{"key": up.Key})
			return
		}
	})
	return r
}

// Admin (izin delivery.dispatch dicek router admin).
func (h *Handler) AdminCouriers(w http.ResponseWriter, r *http.Request) {
	l, err := h.S.ListCouriers(r.Context())
	h.reply(w, r, map[string]any{"couriers": l}, err)
}

func (h *Handler) AdminSaveCourier(w http.ResponseWriter, r *http.Request) {
	var in CourierInput
	if !decode(w, r, &in) {
		return
	}
	c, err := h.S.UpsertCourier(r.Context(), uid(r), in)
	h.reply(w, r, c, err)
}

func (h *Handler) AdminDeliveries(w http.ResponseWriter, r *http.Request) {
	l, err := h.S.AdminDeliveries(r.Context(), r.URL.Query().Get("status"))
	h.reply(w, r, map[string]any{"deliveries": l}, err)
}

func (h *Handler) AdminAssign(w http.ResponseWriter, r *http.Request) {
	var in struct {
		CourierID string `json:"courier_id"`
	}
	if !decode(w, r, &in) {
		return
	}
	t, err := h.S.Assign(r.Context(), uid(r), chi.URLParam(r, "id"), in.CourierID)
	h.reply(w, r, t, err)
}

// AdminProof: foto bukti pengantaran (privat; tidak pernah di-cache).
func (h *Handler) AdminProof(w http.ResponseWriter, r *http.Request) {
	key, err := h.S.ProofKey(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		h.reply(w, r, nil, err)
		return
	}
	f, err := h.S.Files.Open(r.Context(), key)
	if err != nil {
		h.reply(w, r, nil, ErrNotFound)
		return
	}
	defer f.Close()
	w.Header().Set("Cache-Control", "private, no-store")
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if len(key) > 4 && key[len(key)-4:] == ".png" {
		w.Header().Set("Content-Type", "image/png")
	} else if len(key) > 5 && key[len(key)-5:] == ".webp" {
		w.Header().Set("Content-Type", "image/webp")
	}
	_, _ = io.Copy(w, f)
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
		writeErr(w, http.StatusNotFound, "not_found", "Task not found.")
	case errors.Is(err, ErrTaken):
		writeErr(w, http.StatusConflict, "taken", "Another courier already took this task.")
	case errors.Is(err, ErrBadStatus):
		writeErr(w, http.StatusConflict, "bad_status", "This task was already updated. Refresh.")
	case errors.Is(err, ErrWrongPackage):
		writeErr(w, http.StatusUnprocessableEntity, "wrong_package", "This label belongs to a different order.")
	case errors.Is(err, ErrTooMany):
		writeErr(w, http.StatusConflict, "too_many", "Deliver your current packages first.")
	case errors.Is(err, ErrMaxAttempts): // courierreturn
		writeErr(w, http.StatusConflict, "max_attempts", "No delivery attempts left. Return the package to the store.")
	default:
		h.Log.Error("delivery", "path", r.URL.Path, "err", err)
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
