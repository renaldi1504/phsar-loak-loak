package notify

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"phsar/internal/auth"
	"phsar/internal/i18n"
)

type Handler struct {
	Pool *pgxpool.Pool
	Log  *slog.Logger
}

type Item struct {
	ID        int64     `json:"id"`
	Kind      string    `json:"kind"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Link      string    `json:"link,omitempty"`
	StoreID   *string   `json:"store_id,omitempty"`
	Read      bool      `json:"read"`
	CreatedAt time.Time `json:"created_at"`
}

// Routes dipasang di /v1/seller/notifications (sudah login).
func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Get("/", h.list)
	r.Post("/read", h.markRead)
	return r
}

// GET /v1/seller/notifications?limit=30 — teks dalam bahasa header X-Lang.
func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	c, _ := auth.UserFrom(r.Context())
	lang := string(i18n.From(r.Context()))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 30
	}
	rows, err := h.Pool.Query(r.Context(), `SELECT id, kind, params, coalesce(link, ''), store_id::text, read_at IS NOT NULL, created_at
		FROM notifications WHERE user_id = $1 ORDER BY id DESC LIMIT $2`, c.Subject, limit)
	if err != nil {
		h.fail(w, err)
		return
	}
	defer rows.Close()
	out := []Item{}
	for rows.Next() {
		var it Item
		var params []byte
		if err := rows.Scan(&it.ID, &it.Kind, &params, &it.Link, &it.StoreID, &it.Read, &it.CreatedAt); err != nil {
			h.fail(w, err)
			return
		}
		var p Params
		_ = json.Unmarshal(params, &p)
		it.Title, it.Body = Render(it.Kind, lang, p)
		out = append(out, it)
	}
	var unread int
	_ = h.Pool.QueryRow(r.Context(), `SELECT count(*) FROM notifications WHERE user_id = $1 AND read_at IS NULL`, c.Subject).Scan(&unread)
	writeJSON(w, http.StatusOK, map[string]any{"notifications": out, "unread": unread})
}

// POST /v1/seller/notifications/read {"ids":[1,2]} atau {"all":true}
func (h *Handler) markRead(w http.ResponseWriter, r *http.Request) {
	c, _ := auth.UserFrom(r.Context())
	var in struct {
		IDs []int64 `json:"ids"`
		All bool    `json:"all"`
	}
	if err := json.NewDecoder(io.LimitReader(r.Body, 16<<10)).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": map[string]string{"code": "bad_json", "message": err.Error()}})
		return
	}
	var err error
	if in.All {
		_, err = h.Pool.Exec(r.Context(), `UPDATE notifications SET read_at = now() WHERE user_id = $1 AND read_at IS NULL`, c.Subject)
	} else if len(in.IDs) > 0 {
		_, err = h.Pool.Exec(r.Context(), `UPDATE notifications SET read_at = now() WHERE user_id = $1 AND id = ANY($2) AND read_at IS NULL`, c.Subject, in.IDs)
	}
	if err != nil {
		h.fail(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) fail(w http.ResponseWriter, err error) {
	h.Log.Error("notify error", "err", err)
	writeJSON(w, http.StatusInternalServerError, map[string]any{"error": map[string]string{"code": "internal", "message": "internal error"}})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
