package ads

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"phsar/internal/auth"
)

type Handler struct {
	S *Service
}

func uid(r *http.Request) string {
	c, _ := auth.UserFrom(r.Context())
	return c.Subject
}

// SellerRoutes: /v1/seller/stores/{storeID}/ads (khusus pemilik toko).
func (h *Handler) SellerRoutes() http.Handler {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		l, err := h.S.SellerList(r.Context(), chi.URLParam(r, "storeID"))
		st, _ := h.S.Settings(r.Context())
		reply(w, http.StatusOK, map[string]any{"campaigns": l, "enabled": st.Enabled, "min_bid": MinBid.StringFixed(2), "max_bid": MaxBid.StringFixed(2),
			"min_budget": MinBudget.StringFixed(2)}, err)
	})
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		var in Input
		if !decode(w, r, &in) {
			return
		}
		c, err := h.S.Create(r.Context(), chi.URLParam(r, "storeID"), uid(r), in)
		reply(w, http.StatusCreated, c, err)
	})
	r.Patch("/{id}", func(w http.ResponseWriter, r *http.Request) {
		var in Update
		if !decode(w, r, &in) {
			return
		}
		c, err := h.S.SellerUpdate(r.Context(), chi.URLParam(r, "storeID"), uid(r), chi.URLParam(r, "id"), in)
		reply(w, http.StatusOK, c, err)
	})
	return r
}

// AdminRoutes: /admin/api/ads (izin ads.review dipasang paket admin).
func (h *Handler) AdminRoutes() http.Handler {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		l, c, err := h.S.AdminList(r.Context(), r.URL.Query().Get("filter"))
		st, _ := h.S.Settings(r.Context())
		reply(w, http.StatusOK, map[string]any{"campaigns": l, "counts": c, "settings": st}, err)
	})
	r.Put("/settings", func(w http.ResponseWriter, r *http.Request) {
		var in Settings
		if !decode(w, r, &in) {
			return
		}
		st, err := h.S.SaveSettings(r.Context(), uid(r), in)
		reply(w, http.StatusOK, st, err)
	})
	r.Post("/{id}/action", func(w http.ResponseWriter, r *http.Request) {
		var in Review
		if !decode(w, r, &in) {
			return
		}
		c, err := h.S.AdminAction(r.Context(), uid(r), chi.URLParam(r, "id"), in)
		reply(w, http.StatusOK, c, err)
	})
	return r
}

const visitorCookie = "phsar_vid"

// Click: GET /ad/{id}?l=km — tagih (kalau bukan bot) lalu arahkan ke halaman produk.
func (h *Handler) Click(w http.ResponseWriter, r *http.Request) {
	lang := r.URL.Query().Get("l")
	if lang != "en" && lang != "km" && lang != "zh" {
		lang = "km"
	}
	visitor := ""
	if c, err := r.Cookie(visitorCookie); err == nil && len(c.Value) == 32 {
		visitor = c.Value
	} else {
		b := make([]byte, 16)
		_, _ = rand.Read(b)
		visitor = hex.EncodeToString(b) // pengunjung baru; tanpa cookie (bot/klik curang) dibatasi per IP
		http.SetCookie(w, &http.Cookie{Name: visitorCookie, Value: visitor, Path: "/ad", MaxAge: 365 * 86400,
			HttpOnly: true, SameSite: http.SameSiteLaxMode, Secure: r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https"})
	}
	ua := r.UserAgent()
	slug, err := h.S.Click(r.Context(), chi.URLParam(r, "id"), lang, visitor, auth.ClientFrom(r).IP, ua, !IsBot(ua))
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Robots-Tag", "noindex, nofollow")
	switch {
	case errors.Is(err, ErrNotFound):
		http.Redirect(w, r, "/"+lang+"/", http.StatusFound)
	case err != nil && slug == "":
		h.S.Log.Error("iklan: klik gagal", "err", err)
		http.Redirect(w, r, "/"+lang+"/", http.StatusFound)
	default:
		if err != nil {
			h.S.Log.Error("iklan: tagihan klik gagal", "err", err)
		}
		http.Redirect(w, r, "/"+lang+"/p/"+slug, http.StatusFound)
	}
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
	case errors.Is(err, ErrState):
		writeErr(w, http.StatusConflict, "ad_state", err.Error())
	default:
		writeErr(w, http.StatusInternalServerError, "internal", "something went wrong")
	}
}

func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16<<10))
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
