package flash

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"phsar/internal/auth"
	"phsar/internal/pricehistory" // pricehistory
)

// SellerRoutes: /v1/seller/stores/{storeID}/flash (pemilik toko).
func (s *Service) SellerRoutes() http.Handler {
	r := chi.NewRouter()
	store := func(r *http.Request) string { return chi.URLParam(r, "storeID") }
	uid := func(r *http.Request) string { c, _ := auth.UserFrom(r.Context()); return c.Subject }
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		v, err := s.SellerView(r.Context(), store(r))
		reply(w, v, err)
	})
	r.Get("/variants", func(w http.ResponseWriter, r *http.Request) {
		// pricehistory: ?at=<jam mulai sesi> → low30 per varian; low30_rule = aturan harga terendah 30 hari nyala
		at, _ := time.Parse(time.RFC3339, r.URL.Query().Get("at"))
		v, err := s.SellerVariants(r.Context(), store(r), r.URL.Query().Get("q"), at)
		rules, rerr := pricehistory.LoadRules(r.Context(), s.Pool)
		if err == nil {
			err = rerr
		}
		reply(w, map[string]any{"variants": v, "low30_rule": rules.Low30Enabled}, err)
	})
	r.Post("/sessions/{id}/items", func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			Items []ItemInput `json:"items"`
		}
		if !decode(w, r, &in) {
			return
		}
		if _, err := s.Submit(r.Context(), store(r), uid(r), chi.URLParam(r, "id"), in.Items); err != nil {
			reply(w, nil, err)
			return
		}
		v, err := s.SellerView(r.Context(), store(r))
		reply(w, v, err)
	})
	r.Post("/items/{item}/{action}", func(w http.ResponseWriter, r *http.Request) {
		if err := s.SellerAction(r.Context(), store(r), uid(r), chi.URLParam(r, "item"), chi.URLParam(r, "action")); err != nil {
			reply(w, nil, err)
			return
		}
		v, err := s.SellerView(r.Context(), store(r))
		reply(w, v, err)
	})
	return r
}

func reply(w http.ResponseWriter, v any, err error) {
	var fe FieldError
	switch {
	case err == nil:
		writeJSON(w, http.StatusOK, v)
	case errors.As(err, &fe):
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"error": map[string]string{"code": "validation", "field": fe.Field, "message": fe.Msg}})
	case errors.Is(err, ErrNotFound):
		writeJSON(w, http.StatusNotFound, map[string]any{"error": map[string]string{"code": "not_found", "message": "not found"}})
	case errors.Is(err, ErrState):
		writeJSON(w, http.StatusConflict, map[string]any{"error": map[string]string{"code": "flash_state", "message": "this can no longer be changed"}})
	default:
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": map[string]string{"code": "internal", "message": "something went wrong"}})
	}
}

func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 256<<10))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": map[string]string{"code": "bad_json", "message": "invalid JSON body"}})
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
