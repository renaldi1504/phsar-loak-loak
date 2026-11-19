package returns

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"phsar/internal/auth"
	"phsar/internal/paging"
	"phsar/internal/storage"
)

type Handler struct {
	S        *Service
	Uploader storage.Uploader
}

func uid(r *http.Request) string {
	c, _ := auth.UserFrom(r.Context())
	return c.Subject
}

// BuyerRoutes dipasang di dalam /v1/shop (sudah login).
func (h *Handler) BuyerRoutes(r chi.Router) {
	r.Get("/orders/{no}/return", func(w http.ResponseWriter, r *http.Request) {
		v, err := h.S.Eligibility(r.Context(), uid(r), chi.URLParam(r, "no"))
		reply(w, v, err)
	})
	r.Post("/orders/{no}/return", func(w http.ResponseWriter, r *http.Request) {
		var in RequestInput
		if !decode(w, r, &in) {
			return
		}
		v, err := h.S.Request(r.Context(), uid(r), chi.URLParam(r, "no"), in)
		replyStatus(w, http.StatusCreated, v, err)
	})
	r.Get("/returns", func(w http.ResponseWriter, r *http.Request) {
		l, err := h.S.BuyerList(r.Context(), uid(r))
		reply(w, map[string]any{"returns": l}, err)
	})
	r.Post("/returns/{id}/cancel", func(w http.ResponseWriter, r *http.Request) {
		v, err := h.S.BuyerCancel(r.Context(), uid(r), chi.URLParam(r, "id"))
		reply(w, v, err)
	})
	r.Get("/returns/{id}/photos/{i}", func(w http.ResponseWriter, r *http.Request) {
		h.photo(w, r, "r.id = $1 AND r.buyer_id = $2", chi.URLParam(r, "id"), uid(r))
	})
	// Foto bukti retur (privat). Hanya jenis return_photo yang boleh lewat jalur pembeli.
	r.Post("/uploads", func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 9<<20)
		f, _, err := r.FormFile("file")
		if err != nil || r.FormValue("kind") != "return_photo" {
			writeErr(w, http.StatusBadRequest, "bad_upload", "send multipart kind=return_photo and file")
			return
		}
		defer f.Close()
		up, err := h.Uploader.Save(r.Context(), uid(r), "return_photo", f)
		switch {
		case errors.Is(err, storage.ErrTooLarge):
			writeErr(w, http.StatusRequestEntityTooLarge, "too_large", err.Error())
		case errors.Is(err, storage.ErrType), errors.Is(err, storage.ErrKind):
			writeErr(w, http.StatusUnsupportedMediaType, "bad_type", err.Error())
		case err != nil:
			h.S.Log.Error("upload foto retur", "err", err)
			writeErr(w, http.StatusInternalServerError, "internal", "upload failed")
		default:
			writeJSON(w, http.StatusCreated, up)
		}
	})
}

// SellerRoutes dipasang di /v1/seller/stores/{storeID}/returns (anggota toko).
func (h *Handler) SellerRoutes() http.Handler {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		pg := paging.FromRequest(r, 20, 100)
		l, total, err := h.S.SellerList(r.Context(), chi.URLParam(r, "storeID"), r.URL.Query().Get("filter"), pg)
		out := map[string]any{"returns": l}
		if pg.On {
			out["page"] = paging.NewMeta(pg, total)
		}
		reply(w, out, err)
	})
	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		v, err := h.S.SellerGet(r.Context(), chi.URLParam(r, "storeID"), chi.URLParam(r, "id"))
		reply(w, v, err)
	})
	r.Post("/{id}/respond", func(w http.ResponseWriter, r *http.Request) {
		var in SellerAnswer
		if !decode(w, r, &in) {
			return
		}
		v, err := h.S.SellerRespond(r.Context(), chi.URLParam(r, "storeID"), uid(r), chi.URLParam(r, "id"), in)
		reply(w, v, err)
	})
	r.Get("/{id}/photos/{i}", func(w http.ResponseWriter, r *http.Request) {
		h.photo(w, r, "r.id = $1 AND r.store_id = $2", chi.URLParam(r, "id"), chi.URLParam(r, "storeID"))
	})
	return r
}

// Admin (perm return.decide) — dipasang paket admin.
func (h *Handler) AdminList(w http.ResponseWriter, r *http.Request) {
	l, c, err := h.S.AdminList(r.Context(), r.URL.Query().Get("status"))
	reply(w, map[string]any{"returns": l, "counts": c}, err)
}

func (h *Handler) AdminGet(w http.ResponseWriter, r *http.Request) {
	v, err := h.S.AdminGet(r.Context(), chi.URLParam(r, "id"))
	reply(w, v, err)
}

func (h *Handler) AdminDecide(w http.ResponseWriter, r *http.Request) {
	var in Decision
	if !decode(w, r, &in) {
		return
	}
	v, err := h.S.Decide(r.Context(), uid(r), chi.URLParam(r, "id"), in)
	reply(w, v, err)
}

func (h *Handler) AdminRefunded(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Reference string `json:"reference"`
		Note      string `json:"note"`
	}
	if !decode(w, r, &in) {
		return
	}
	v, err := h.S.MarkRefunded(r.Context(), uid(r), chi.URLParam(r, "id"), in.Reference, in.Note)
	reply(w, v, err)
}

func (h *Handler) AdminPhoto(w http.ResponseWriter, r *http.Request) {
	h.photo(w, r, "r.id = $1", chi.URLParam(r, "id"))
}

func (h *Handler) photo(w http.ResponseWriter, r *http.Request, where string, args ...any) {
	i, err := strconv.Atoi(chi.URLParam(r, "i"))
	if err != nil {
		writeErr(w, http.StatusNotFound, "not_found", "photo not found")
		return
	}
	f, ctype, err := h.S.Photo(r.Context(), where, i, args...)
	if err != nil {
		writeErr(w, http.StatusNotFound, "not_found", "photo not found")
		return
	}
	defer f.Close()
	w.Header().Set("Content-Type", ctype)
	w.Header().Set("Cache-Control", "private, no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	_, _ = io.Copy(w, f)
}

// ---------------------------------------------------------------------------

func reply(w http.ResponseWriter, v any, err error) { replyStatus(w, http.StatusOK, v, err) }

func replyStatus(w http.ResponseWriter, status int, v any, err error) {
	var fe FieldError
	switch {
	case err == nil:
		writeJSON(w, status, v)
	case errors.As(err, &fe):
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"error": map[string]string{"code": "validation", "field": fe.Field, "message": fe.Msg}})
	case errors.Is(err, ErrNotFound):
		writeErr(w, http.StatusNotFound, "not_found", err.Error())
	case errors.Is(err, ErrNotAllowed):
		writeErr(w, http.StatusConflict, "return_not_allowed", err.Error())
	case errors.Is(err, ErrState):
		writeErr(w, http.StatusConflict, "return_state", err.Error())
	case errors.Is(err, ErrDeadline):
		writeErr(w, http.StatusConflict, "return_deadline", err.Error())
	default:
		writeErr(w, http.StatusInternalServerError, "internal", "something went wrong")
	}
}

func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64<<10))
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
