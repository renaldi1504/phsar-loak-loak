package health

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct{ pool *pgxpool.Pool }

func New(pool *pgxpool.Pool) *Handler { return &Handler{pool: pool} }

// Live: proses hidup (dipakai Docker/Caddy).
func (h *Handler) Live(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

// Ready: database bisa dijangkau.
func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	w.Header().Set("Content-Type", "application/json")
	if err := h.pool.Ping(ctx); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "db_unreachable"})
		return
	}
	_, _ = w.Write([]byte(`{"status":"ready"}`))
}
