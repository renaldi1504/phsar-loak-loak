package chat

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"phsar/internal/auth"
	"phsar/internal/i18n"
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

func lang(r *http.Request) string { return string(i18n.From(r.Context())) }

// BuyerRoutes dipasang di dalam /v1/shop (sudah login). Akun tamu tidak boleh chat.
func (h *Handler) BuyerRoutes(r chi.Router) {
	scope := func(r *http.Request) Scope { return Scope{UserID: uid(r), Lang: lang(r)} }
	r.With(auth.NoGuest).Route("/chat", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			var in OpenInput
			if !decode(w, r, &in) {
				return
			}
			c, created, err := h.S.Open(r.Context(), scope(r), in)
			status := http.StatusOK
			if created {
				status = http.StatusCreated
			}
			reply(w, status, c, err)
		})
		h.common(r, scope)
	})
}

// SellerRoutes dipasang di /v1/seller/stores/{storeID}/chat (anggota toko: pemilik & staf).
func (h *Handler) SellerRoutes() http.Handler {
	r := chi.NewRouter()
	h.common(r, func(r *http.Request) Scope {
		return Scope{UserID: uid(r), StoreID: chi.URLParam(r, "storeID"), Lang: lang(r)}
	})
	return r
}

// common: daftar, isi (polling ?after=), kirim, upload foto, lihat foto — sama untuk pembeli & toko.
func (h *Handler) common(r chi.Router, scope func(*http.Request) Scope) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		l, unread, waiting, err := h.S.List(r.Context(), scope(r))
		reply(w, http.StatusOK, map[string]any{"conversations": l, "unread": unread, "waiting": waiting}, err)
	})
	r.Post("/uploads", func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 9<<20)
		f, _, err := r.FormFile("file")
		if err != nil || (r.FormValue("kind") != "" && r.FormValue("kind") != PhotoKind) {
			writeErr(w, http.StatusBadRequest, "bad_upload", "send multipart kind=chat_photo and file")
			return
		}
		defer f.Close()
		up, err := h.Uploader.Save(r.Context(), uid(r), PhotoKind, f)
		switch {
		case errors.Is(err, storage.ErrTooLarge):
			writeErr(w, http.StatusRequestEntityTooLarge, "too_large", err.Error())
		case errors.Is(err, storage.ErrType), errors.Is(err, storage.ErrKind):
			writeErr(w, http.StatusUnsupportedMediaType, "bad_type", err.Error())
		case err != nil:
			h.S.Log.Error("upload foto chat", "err", err)
			writeErr(w, http.StatusInternalServerError, "internal", "upload failed")
		default:
			writeJSON(w, http.StatusCreated, up)
		}
	})
	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		after, _ := strconv.ParseInt(r.URL.Query().Get("after"), 10, 64)
		c, err := h.S.Get(r.Context(), scope(r), chi.URLParam(r, "id"), max(after, 0))
		reply(w, http.StatusOK, c, err)
	})
	r.Post("/{id}/messages", func(w http.ResponseWriter, r *http.Request) {
		var in SendInput
		if !decode(w, r, &in) {
			return
		}
		m, err := h.S.Send(r.Context(), scope(r), chi.URLParam(r, "id"), in)
		reply(w, http.StatusCreated, m, err)
	})
	r.Get("/{id}/messages/{mid}/photo", func(w http.ResponseWriter, r *http.Request) {
		sc := scope(r)
		h.photo(w, r, &sc)
	})
}

// AdminRoutes dipasang paket admin (izin support.handle). adminID = admin yang sedang login.
func (h *Handler) AdminRoutes(r chi.Router, adminID func(*http.Request) string) {
	r.Get("/chats", func(w http.ResponseWriter, r *http.Request) {
		l, err := h.S.AdminList(r.Context(), r.URL.Query().Get("q"))
		reply(w, http.StatusOK, map[string]any{"conversations": l}, err)
	})
	r.Get("/chats/{id}", func(w http.ResponseWriter, r *http.Request) {
		c, err := h.S.AdminGet(r.Context(), adminID(r), chi.URLParam(r, "id"))
		reply(w, http.StatusOK, c, err)
	})
	r.Get("/chats/{id}/messages/{mid}/photo", func(w http.ResponseWriter, r *http.Request) { h.photo(w, r, nil) })
}

func (h *Handler) photo(w http.ResponseWriter, r *http.Request, sc *Scope) {
	mid, err := strconv.ParseInt(chi.URLParam(r, "mid"), 10, 64)
	if err != nil {
		writeErr(w, http.StatusNotFound, "not_found", "photo not found")
		return
	}
	f, ctype, err := h.S.Photo(r.Context(), sc, chi.URLParam(r, "id"), mid)
	if err != nil {
		writeErr(w, http.StatusNotFound, "not_found", "photo not found")
		return
	}
	defer f.Close()
	w.Header().Set("Content-Type", ctype)
	w.Header().Set("Cache-Control", "private, max-age=3600")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	_, _ = io.Copy(w, f)
}

// ---------------------------------------------------------------------------

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
		writeErr(w, http.StatusTooManyRequests, "chat_rate_limited", err.Error())
	default:
		slog.Error("chat", "err", err)
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
