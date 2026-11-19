package httpx

// Fungsi untuk tombol "Run now" di admin → Scheduled jobs (proses API). Worker latar memakai fungsi yang sama lewat jobs.Track.

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"phsar/internal/config"
	"phsar/internal/jobs"
	"phsar/internal/media"
	"phsar/internal/payout"
	"phsar/internal/returns"
	"phsar/internal/shop"
)

// drainMax: kosongkan antrean, tapi paling banyak max item per klik (supaya request admin tidak menggantung).
func drainMax(once func(context.Context) (bool, error), max int) jobs.Func {
	return func(ctx context.Context) (int, error) {
		n := 0
		for n < max && ctx.Err() == nil {
			did, err := once(ctx)
			if err != nil {
				return n, err
			}
			if !did {
				break
			}
			n++
		}
		return n, nil
	}
}

func registerJobs(cfg config.Config, log *slog.Logger, pool *pgxpool.Pool, shopSvc *shop.Service, paySvc *payout.Service, retSvc *returns.Service) {
	jobs.SetRunner("checkout.expire", shopSvc.ExpireDue)
	jobs.SetRunner("order.accept_timeout", shopSvc.ExpireUnaccepted)
	jobs.SetRunner("order.pack_timeout", shopSvc.ExpireUnpacked)
	jobs.SetRunner("payout.complete", paySvc.CompleteDue)
	jobs.SetRunner("returns.escalate", retSvc.EscalateOverdue)
	jobs.SetRunner("notify.send", NewNotifyWorker(cfg, log, pool).RunOnce)
	if tw := NewTranslationWorker(cfg, log, pool); tw.T != nil {
		jobs.SetRunner("translate.queue", drainMax(tw.RunOnce, 20))
	}
	if media.Available() {
		mw := &media.Worker{Pool: pool, Files: newFiles(cfg), Log: log}
		jobs.SetRunner("media.queue", drainMax(mw.RunOnce, 2)) // transcode berat: 2 item per klik
	}
	if ps := NewPhotoSearch(cfg, log, pool); ps.Enabled() {
		jobs.SetRunner("photo.index", func(ctx context.Context) (int, error) { return ps.IndexPending(ctx, 40) })
	}
}
