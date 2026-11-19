// Package tour menyimpan status "tur cara pakai sudah dilihat" untuk pengguna yang login,
// supaya tur (Seller Center / website pembeli) tidak muncul lagi di perangkat lain.
//
//	GET  /v1/app/tours            → {"done": {"site": "2026-…"}}   (pembeli, bukan tamu)
//	POST /v1/app/tours            {"tour_key":"site"}
//	GET  /v1/seller/tours         → {"done": {"seller_intro": "…"}} (akun seller)
//	POST /v1/seller/tours         {"tour_key":"seller_orders"}
//
// POST idempoten: memanggil dua kali tidak mengubah done_at pertama.
package tour

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"phsar/internal/auth"
	"phsar/internal/platform/db"
)

// Daftar tur yang diizinkan per pintu. Tambah kunci baru di sini (dan di JS) kalau ada tur baru.
var (
	BuyerKeys  = []string{"site"}
	SellerKeys = []string{"seller_intro", "seller_orders", "seller_product"}
)

type Handler struct {
	Pool *pgxpool.Pool
	Log  *slog.Logger
}

// BuyerRoutes dipasang di /v1/app/tours (sudah login, tamu ditolak).
func (h *Handler) BuyerRoutes() http.Handler { return h.routes(BuyerKeys, "buyer") }

// SellerRoutes dipasang di /v1/seller/tours (status per akun, bukan per toko).
func (h *Handler) SellerRoutes() http.Handler { return h.routes(SellerKeys, "seller") }

func (h *Handler) routes(keys []string, actor string) http.Handler {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) { h.list(w, r, keys) })
	r.Post("/", func(w http.ResponseWriter, r *http.Request) { h.done(w, r, keys, actor) })
	return r
}

func uid(r *http.Request) string {
	c, _ := auth.UserFrom(r.Context())
	return c.Subject
}

// GET: tur yang sudah dilihat (hanya kunci milik pintu ini).
func (h *Handler) list(w http.ResponseWriter, r *http.Request, keys []string) {
	rows, err := h.Pool.Query(r.Context(), `SELECT tour_key, done_at FROM user_tours WHERE user_id = $1 AND tour_key = ANY($2)`, uid(r), keys)
	if err != nil {
		h.internal(w, err)
		return
	}
	done := map[string]time.Time{}
	for rows.Next() {
		var k string
		var at time.Time
		if err := rows.Scan(&k, &at); err != nil {
			rows.Close()
			h.internal(w, err)
			return
		}
		done[k] = at
	}
	if err := rows.Err(); err != nil {
		h.internal(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"done": done, "keys": keys})
}

// POST: tandai tur selesai/dilewati.
func (h *Handler) done(w http.ResponseWriter, r *http.Request, keys []string, actor string) {
	var in struct {
		TourKey string `json:"tour_key"`
	}
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4<<10))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&in); err != nil {
		writeErr(w, http.StatusBadRequest, "bad_json", "invalid JSON body")
		return
	}
	key := strings.TrimSpace(in.TourKey)
	if key == "" {
		validation(w, "tour_key", "required")
		return
	}
	if !slices.Contains(keys, key) {
		validation(w, "tour_key", "unknown tour")
		return
	}
	user := uid(r)
	var at time.Time
	err := db.WithActor(r.Context(), h.Pool, user, actor, "", func(tx pgx.Tx) error {
		if _, err := tx.Exec(r.Context(), `INSERT INTO user_tours (user_id, tour_key) VALUES ($1, $2) ON CONFLICT DO NOTHING`, user, key); err != nil {
			return err
		}
		return tx.QueryRow(r.Context(), `SELECT done_at FROM user_tours WHERE user_id = $1 AND tour_key = $2`, user, key).Scan(&at)
	})
	if err != nil {
		h.internal(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"tour_key": key, "done_at": at})
}

func (h *Handler) internal(w http.ResponseWriter, err error) {
	if h.Log != nil {
		h.Log.Error("tour: status tur", "err", err)
	}
	writeErr(w, http.StatusInternalServerError, "internal", "something went wrong")
}

func validation(w http.ResponseWriter, field, msg string) {
	writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"error": map[string]string{"code": "validation", "field": field, "message": msg}})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, code, msg string) {
	writeJSON(w, status, map[string]any{"error": map[string]string{"code": code, "message": msg}})
}
