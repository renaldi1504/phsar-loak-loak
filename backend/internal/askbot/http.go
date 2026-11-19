package askbot

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

	"phsar/internal/auth"
	"phsar/internal/i18n"
	"phsar/internal/platform/db"
)

// SellerRoutes: /v1/seller/stores/{storeID}/assist (member toko: pemilik & staf; dicek paket seller).
//
//	POST   /suggest     {"question", "product_id"?} → {"reply"}   (KONTRAK — dipakai modul chat)
//	GET    /faqs        ?product_id=<uuid>|store   → {"faqs", "enabled", "ai"}
//	POST   /faqs        FAQInput → FAQ
//	PUT    /faqs/{id}   FAQInput → FAQ
//	DELETE /faqs/{id}
func (s *Service) SellerRoutes() http.Handler {
	r := chi.NewRouter()
	store := func(r *http.Request) string { return chi.URLParam(r, "storeID") }
	r.Post("/suggest", func(w http.ResponseWriter, r *http.Request) {
		var in SuggestInput
		if !decode(w, r, &in) {
			return
		}
		res, err := s.Suggest(r.Context(), store(r), uid(r), in)
		reply(w, http.StatusOK, res, err)
	})
	r.Get("/faqs", func(w http.ResponseWriter, r *http.Request) {
		list, err := s.ListFAQs(r.Context(), store(r), string(i18n.From(r.Context())), r.URL.Query().Get("product_id"))
		if err != nil {
			reply(w, 0, nil, err)
			return
		}
		st, err := s.Settings(r.Context())
		reply(w, http.StatusOK, map[string]any{"faqs": list, "enabled": st.Enabled, "ai": s.aiOn(r.Context()), "max": MaxFAQs}, err)
	})
	r.Post("/faqs", func(w http.ResponseWriter, r *http.Request) {
		var in FAQInput
		if !decode(w, r, &in) {
			return
		}
		f, err := s.CreateFAQ(r.Context(), store(r), uid(r), in)
		reply(w, http.StatusCreated, f, err)
	})
	r.Put("/faqs/{id}", func(w http.ResponseWriter, r *http.Request) {
		var in FAQInput
		if !decode(w, r, &in) {
			return
		}
		f, err := s.UpdateFAQ(r.Context(), store(r), uid(r), chi.URLParam(r, "id"), in)
		reply(w, http.StatusOK, f, err)
	})
	r.Delete("/faqs/{id}", func(w http.ResponseWriter, r *http.Request) {
		err := s.DeleteFAQ(r.Context(), store(r), uid(r), chi.URLParam(r, "id"))
		reply(w, http.StatusOK, map[string]bool{"ok": true}, err)
	})
	return r
}

// AdminRoutes: /admin/api/askbot (super admin; izin dicek paket admin).
func (s *Service) AdminRoutes(model string) http.Handler {
	r := chi.NewRouter()
	view := func(w http.ResponseWriter, r *http.Request) {
		st, err := s.Settings(r.Context())
		if err != nil {
			reply(w, 0, nil, err)
			return
		}
		u, err := s.Usage(r.Context())
		m := ""
		if s.Gen != nil {
			m = model
		}
		reply(w, http.StatusOK, map[string]any{"settings": st, "ai": s.aiOn(r.Context()), "model": m, "usage": u}, err)
	}
	r.Get("/", view)
	r.Put("/", func(w http.ResponseWriter, r *http.Request) {
		var in Settings
		if !decode(w, r, &in) {
			return
		}
		if err := in.Check(); err != nil {
			reply(w, 0, nil, err)
			return
		}
		err := db.WithActor(r.Context(), s.Pool, uid(r), "admin", "", func(tx pgx.Tx) error {
			_, err := tx.Exec(r.Context(), `INSERT INTO platform_settings (key, value, updated_by) VALUES ('askbot', $1, $2)
				ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_by = EXCLUDED.updated_by, updated_at = now()`, in, uid(r))
			return err
		})
		if err != nil {
			reply(w, 0, nil, err)
			return
		}
		view(w, r)
	})
	return r
}

func uid(r *http.Request) string {
	c, _ := auth.UserFrom(r.Context())
	return c.Subject
}

func reply(w http.ResponseWriter, status int, v any, err error) {
	var fe FieldError
	switch {
	case err == nil:
		writeJSON(w, status, v)
	case errors.As(err, &fe):
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"error": map[string]string{"code": "validation", "field": fe.Field, "message": fe.Msg}})
	case errors.Is(err, ErrNotFound):
		writeJSON(w, http.StatusNotFound, map[string]any{"error": map[string]string{"code": "not_found", "message": "not found"}})
	case errors.Is(err, ErrDisabled):
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{"error": map[string]string{"code": "assist_unavailable", "message": "the reply assistant is turned off"}})
	case errors.Is(err, ErrRateLimited):
		w.Header().Set("Retry-After", "600")
		writeJSON(w, http.StatusTooManyRequests, map[string]any{"error": map[string]string{"code": "rate_limited", "message": "too many suggestions this hour, try again later"}})
	default:
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": map[string]string{"code": "internal", "message": "something went wrong"}})
	}
}

func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64<<10))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": map[string]string{"code": "bad_json", "message": "invalid JSON body"}})
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
