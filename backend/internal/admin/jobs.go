package admin

// Scheduled jobs (super admin): status worker latar, riwayat run, jalankan sekarang, jeda/lanjutkan.

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

	"phsar/internal/jobs"
	"phsar/internal/platform/db"
)

func (h *Handler) mountJobs(r chi.Router) {
	r.Route("/jobs", func(p chi.Router) {
		p.Use(h.can("*")) // hanya super admin
		p.Get("/", h.jobsList)
		p.Get("/{name}/runs", h.jobRuns)
		p.Post("/{name}/run", h.jobRunNow)
		p.Post("/{name}/pause", h.jobPause)
	})
}

type jobRow struct {
	jobs.Def
	Runnable  bool            `json:"runnable"` // tombol "Run now" tersedia di server ini
	Paused    bool            `json:"paused"`
	State     json.RawMessage `json:"state"`    // detak terakhir (null = belum pernah jalan)
	LastRun   json.RawMessage `json:"last_run"` // run tercatat terakhir
	Runs24h   int             `json:"runs_24h"`
	Failed24h int             `json:"failed_24h"`
	Items24h  int             `json:"items_24h"`
	AvgMs     *int            `json:"avg_ms"`
}

func (h *Handler) jobsList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paused, err := jobs.Paused(ctx, h.pool)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "internal", err.Error())
		return
	}
	out := []jobRow{}
	for _, d := range jobs.Defs {
		_, runnable := jobs.Runner(d.Name)
		row := jobRow{Def: d, Runnable: runnable && d.Safe, Paused: paused[d.Name], State: json.RawMessage("null"), LastRun: json.RawMessage("null")}
		err := h.pool.QueryRow(ctx, `SELECT
			(SELECT row_to_json(s) FROM (SELECT last_tick_at, last_success_at, last_error, last_error_at FROM job_state WHERE job = $1) s),
			(SELECT row_to_json(x) FROM (SELECT id, started_at, finished_at, status, rows_affected AS items, error, duration_ms, trigger
				FROM job_runs WHERE job = $1 ORDER BY started_at DESC LIMIT 1) x),
			count(*), count(*) FILTER (WHERE status = 'failed'), coalesce(sum(rows_affected), 0), avg(duration_ms)::int
			FROM job_runs WHERE job = $1 AND started_at > now() - interval '24 hours'`, d.Name).
			Scan(&row.State, &row.LastRun, &row.Runs24h, &row.Failed24h, &row.Items24h, &row.AvgMs)
		if err != nil {
			h.log.Error("daftar job", "err", err)
			writeErr(w, http.StatusInternalServerError, "internal", "could not load jobs")
			return
		}
		if row.State == nil {
			row.State = json.RawMessage("null")
		}
		if row.LastRun == nil {
			row.LastRun = json.RawMessage("null")
		}
		out = append(out, row)
	}
	writeJSON(w, http.StatusOK, map[string]any{"jobs": out})
}

func (h *Handler) jobDef(w http.ResponseWriter, r *http.Request) (jobs.Def, bool) {
	d, ok := jobs.Find(chi.URLParam(r, "name"))
	if !ok {
		writeErr(w, http.StatusNotFound, "not_found", "unknown job")
	}
	return d, ok
}

func (h *Handler) jobRuns(w http.ResponseWriter, r *http.Request) {
	d, ok := h.jobDef(w, r)
	if !ok {
		return
	}
	var raw []byte
	err := h.pool.QueryRow(r.Context(), `SELECT coalesce(json_agg(x), '[]') FROM (
		SELECT j.id, j.started_at, j.finished_at, j.status, j.rows_affected AS items, j.error, j.duration_ms, j.trigger,
		  coalesce(nullif(u.name, ''), u.phone) AS triggered_by
		FROM job_runs j LEFT JOIN users u ON u.id = j.triggered_by WHERE j.job = $1 ORDER BY j.started_at DESC LIMIT 50) x`, d.Name).Scan(&raw)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "internal", "could not load runs")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"runs":` + string(raw) + `}`))
}

// POST /jobs/{name}/run — jalankan satu putaran sekarang (juga kalau dijeda). Dicatat di job_runs + audit.
func (h *Handler) jobRunNow(w http.ResponseWriter, r *http.Request) {
	d, ok := h.jobDef(w, r)
	if !ok {
		return
	}
	if !d.Safe {
		writeErr(w, http.StatusConflict, "not_runnable", "this job cannot be run manually")
		return
	}
	n, err := jobs.RunNow(r.Context(), h.pool, d.Name, actor(r))
	_, _ = h.pool.Exec(r.Context(), `INSERT INTO audit_logs (actor_type, actor_user_id, action, entity_table, entity_id)
		VALUES ('admin', $1, 'job.run', 'job_runs', $2)`, actor(r), d.Name)
	switch {
	case errors.Is(err, jobs.ErrNoRunner):
		writeErr(w, http.StatusConflict, "not_runnable", "this job is not available on this server (check its settings)")
	case err != nil:
		writeJSON(w, http.StatusOK, map[string]any{"ok": false, "items": n, "error": err.Error()})
	default:
		writeJSON(w, http.StatusOK, map[string]any{"ok": true, "items": n})
	}
}

// POST /jobs/{name}/pause {paused} — worker melewati job yang dijeda (berlaku ≤ 15 detik di semua proses).
func (h *Handler) jobPause(w http.ResponseWriter, r *http.Request) {
	d, ok := h.jobDef(w, r)
	if !ok {
		return
	}
	var in struct {
		Paused bool `json:"paused"`
	}
	if !decode(w, r, &in) {
		return
	}
	err := db.WithActor(r.Context(), h.pool, actor(r), "admin", "", func(tx pgx.Tx) error {
		return jobs.SetPaused(r.Context(), tx, actor(r), d.Name, in.Paused)
	})
	if err != nil {
		h.log.Error("jeda job", "err", err)
		writeErr(w, http.StatusInternalServerError, "internal", "could not save")
		return
	}
	h.jobsList(w, r)
}
