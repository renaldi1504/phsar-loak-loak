package qa

import (
	"context"
	"errors"
	"testing"
	"time"

	"phsar/internal/jobs"
)

// ---------------------------------------------------------------------------
// Pekerjaan latar: pencatatan run, jeda, jalankan sekarang (super admin)
// ---------------------------------------------------------------------------

func TestScheduledJobs(t *testing.T) {
	e := getEnv(t)
	ctx := context.Background()
	sa := newAdmin(t, e, "super_admin")
	cs := newAdmin(t, e, "cs")
	oldTTL, oldBeat := jobs.CacheTTL, jobs.BeatEvery
	jobs.CacheTTL, jobs.BeatEvery = 0, time.Hour
	t.Cleanup(func() {
		jobs.CacheTTL, jobs.BeatEvery = oldTTL, oldBeat
		e.sql(t, `DELETE FROM platform_settings WHERE key = 'jobs'`)
	})
	count := func(job string) string { return e.scalar(t, `SELECT count(*) FROM job_runs WHERE job = $1`, job) }

	t.Run("Track: putaran kosong hanya detak, yang berisi & error dicatat", func(t *testing.T) {
		job := "qa.track"
		n, err := jobs.Track(ctx, e.pool, job, func(context.Context) (int, error) { return 0, nil })
		check(t, n == 0 && err == nil, "Track: %d %v", n, err)
		check(t, count(job) == "0", "putaran kosong tidak boleh masuk job_runs")
		check(t, e.scalar(t, `SELECT (last_tick_at IS NOT NULL)::text FROM job_state WHERE job = $1`, job) == "true", "detak harus tersimpan")
		_, _ = jobs.Track(ctx, e.pool, job, func(context.Context) (int, error) { return 3, nil })
		check(t, count(job) == "1" && e.scalar(t, `SELECT rows_affected::text || status || trigger FROM job_runs WHERE job = $1`, job) == "3successauto", "run berisi harus dicatat")
		_, err = jobs.Track(ctx, e.pool, job, func(context.Context) (int, error) { return 0, errors.New("boom") })
		check(t, err != nil && count(job) == "2", "error harus dicatat")
		check(t, e.scalar(t, `SELECT last_error FROM job_state WHERE job = $1`, job) == "boom", "error terakhir di job_state")
		n, _ = jobs.Drain(ctx, e.pool, job, func() func(context.Context) (bool, error) {
			k := 0
			return func(context.Context) (bool, error) { k++; return k <= 4, nil }
		}())
		check(t, n == 4, "Drain harus memproses 4 item: %d", n)
	})

	t.Run("daftar job: hanya super admin", func(t *testing.T) {
		e.call(t, "GET", "/admin/api/jobs", nil, cs).want(403)
		r := e.call(t, "GET", "/admin/api/jobs", nil, sa).want(200)
		check(t, r.len("jobs") == len(jobs.Defs), "semua job harus tampil: %d", r.len("jobs"))
		r.eq("jobs.0.name", "checkout.expire").eq("jobs.0.runnable", "true").eq("jobs.4.name", "payout.schedule").eq("jobs.4.runnable", "false")
	})

	t.Run("jeda: worker melewati job", func(t *testing.T) {
		e.call(t, "POST", "/admin/api/jobs/checkout.expire/pause", map[string]bool{"paused": true}, cs).want(403)
		e.call(t, "POST", "/admin/api/jobs/nope/pause", map[string]bool{"paused": true}, sa).want(404)
		e.call(t, "POST", "/admin/api/jobs/checkout.expire/pause", map[string]bool{"paused": true}, sa).want(200).eq("jobs.0.paused", "true")
		ran := false
		_, _ = jobs.Track(ctx, e.pool, "checkout.expire", func(context.Context) (int, error) { ran = true; return 1, nil })
		check(t, !ran, "job yang dijeda tidak boleh jalan")
		e.call(t, "POST", "/admin/api/jobs/checkout.expire/pause", map[string]bool{"paused": false}, sa).want(200).eq("jobs.0.paused", "false")
		_, _ = jobs.Track(ctx, e.pool, "checkout.expire", func(context.Context) (int, error) { ran = true; return 0, nil })
		check(t, ran, "setelah dilanjutkan job harus jalan lagi")
	})

	t.Run("jalankan sekarang: dicatat manual dengan nama admin", func(t *testing.T) {
		e.call(t, "POST", "/admin/api/jobs/checkout.expire/run", nil, cs).want(403)
		e.call(t, "POST", "/admin/api/jobs/payout.schedule/run", nil, sa).wantErr(409, "not_runnable")
		e.call(t, "POST", "/admin/api/jobs/nope/run", nil, sa).want(404)
		before := e.scalar(t, `SELECT count(*) FROM job_runs WHERE job = 'returns.escalate' AND trigger = 'manual'`)
		e.call(t, "POST", "/admin/api/jobs/returns.escalate/run", nil, sa).want(200).eq("ok", "true")
		check(t, e.scalar(t, `SELECT count(*) FROM job_runs WHERE job = 'returns.escalate' AND trigger = 'manual' AND triggered_by IS NOT NULL`) != before, "run manual harus tercatat")
		h := e.call(t, "GET", "/admin/api/jobs/returns.escalate/runs", nil, sa).want(200)
		h.eq("runs.0.trigger", "manual")
		check(t, h.str("runs.0.triggered_by") != "", "nama admin harus tampil")
	})
}
