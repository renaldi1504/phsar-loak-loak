package media

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"phsar/internal/jobs"
	"phsar/internal/storage"
)

// Worker memproses antrean media_processing_jobs satu per satu (transcode cukup berat untuk CPU).
type Worker struct {
	Pool     *pgxpool.Pool
	Files    storage.Storage
	Log      *slog.Logger
	Interval time.Duration
}

const maxAttempts = 3

var errPermanent = errors.New("permanent")

func (w *Worker) Run(ctx context.Context) {
	if !Available() {
		w.Log.Warn("ffmpeg tidak terpasang: video & foto dipakai apa adanya (Mac: brew install ffmpeg)")
		return
	}
	if w.Interval == 0 {
		w.Interval = 3 * time.Second
	}
	_, _ = w.Pool.Exec(ctx, `UPDATE media_processing_jobs SET status = 'queued' WHERE status = 'running'`)
	w.Log.Info("worker media berjalan (kompres video + thumbnail)")
	t := time.NewTicker(w.Interval)
	defer t.Stop()
	for {
		if _, err := jobs.Drain(ctx, w.Pool, "media.queue", w.RunOnce); err != nil && ctx.Err() == nil { // jobs: kosongkan antrean
			w.Log.Error("worker media", "err", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}
	}
}

type job struct {
	id, mediaID, kind, key string
	attempts               int
}

// RunOnce memproses satu job. false = antrean kosong.
func (w *Worker) RunOnce(ctx context.Context) (bool, error) {
	var j job
	err := w.Pool.QueryRow(ctx, `
		WITH next AS (
			SELECT id FROM media_processing_jobs
			WHERE status = 'queued' AND step = 'transcode' AND (run_after IS NULL OR run_after <= now())
			ORDER BY created_at LIMIT 1 FOR UPDATE SKIP LOCKED)
		UPDATE media_processing_jobs mj SET status = 'running', attempts = mj.attempts + 1, started_at = now()
		FROM next, product_media pm
		WHERE mj.id = next.id AND pm.id = mj.media_id
		RETURNING mj.id, mj.media_id, pm.kind::text, pm.storage_key, mj.attempts`).Scan(&j.id, &j.mediaID, &j.kind, &j.key, &j.attempts)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	start := time.Now()
	perr := w.process(ctx, j)
	switch {
	case perr == nil:
		w.Log.Info("media selesai diproses", "media", j.mediaID, "kind", j.kind, "detik", time.Since(start).Round(100*time.Millisecond).Seconds())
		_, err = w.Pool.Exec(ctx, `UPDATE media_processing_jobs SET status = 'done', error = NULL, finished_at = now() WHERE id = $1`, j.id)
	case errors.Is(perr, errPermanent) || errors.Is(perr, ErrCorrupt) || j.attempts >= maxAttempts:
		w.Log.Warn("media gagal diproses (berhenti; file asli tetap dipakai)", "media", j.mediaID, "err", perr)
		_, err = w.Pool.Exec(ctx, `UPDATE media_processing_jobs SET status = 'failed', error = $2, finished_at = now() WHERE id = $1`, j.id, trunc(perr.Error()))
		// Produk tetap bisa tampil dengan file asli.
		_, _ = w.Pool.Exec(ctx, `UPDATE product_media SET processing = 'ready' WHERE id = $1 AND processing = 'processing'`, j.mediaID)
	default:
		w.Log.Warn("media gagal diproses, dicoba lagi", "media", j.mediaID, "err", perr)
		_, err = w.Pool.Exec(ctx, `UPDATE media_processing_jobs SET status = 'queued', error = $2,
			run_after = now() + make_interval(mins => attempts * 2) WHERE id = $1`, j.id, trunc(perr.Error()))
	}
	return true, err
}

func (w *Worker) process(ctx context.Context, j job) error {
	dir, err := os.MkdirTemp("", "phsar-media-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	in := filepath.Join(dir, "in"+filepath.Ext(j.key))
	origBytes, err := w.download(ctx, j.key, in)
	if err != nil {
		return fmt.Errorf("%w: file asli tidak ada: %v", errPermanent, err)
	}
	info, err := Probe(ctx, in)
	if err != nil {
		return fmt.Errorf("%w: %v", errPermanent, err)
	}
	base := strings.TrimSuffix(j.key, filepath.Ext(j.key))
	tctx, cancel := context.WithTimeout(ctx, 15*time.Minute)
	defer cancel()

	if j.kind == "video" {
		if info.Duration > MaxVideoSeconds+1 {
			return fmt.Errorf("%w: %v", errPermanent, ErrTooLong)
		}
		outV, outP := filepath.Join(dir, "v.mp4"), filepath.Join(dir, "p.jpg")
		if err := TranscodeVideo(tctx, in, outV); err != nil {
			return err
		}
		if err := Poster(tctx, outV, outP, info.Duration); err != nil {
			return err
		}
		vi, err := Probe(ctx, outV)
		if err != nil {
			return err
		}
		vKey, pKey := base+"-720.mp4", base+"-poster.jpg"
		vBytes, err := w.upload(ctx, outV, vKey)
		if err != nil {
			return err
		}
		if _, err := w.upload(ctx, outP, pKey); err != nil {
			return err
		}
		// Kalau hasil kompres malah lebih besar (video asli sudah kecil), pakai yang asli.
		videoURL := w.Files.PublicURL(vKey)
		if vBytes >= origBytes {
			videoURL, vBytes = w.Files.PublicURL(j.key), origBytes
		}
		dur := int(math.Ceil(info.Duration))
		if dur < 1 {
			dur = 1
		}
		if dur > MaxVideoSeconds {
			dur = MaxVideoSeconds
		}
		_, err = w.Pool.Exec(ctx, `UPDATE product_media SET video_url = $2, poster_url = $3, duration_sec = $4, has_audio = $5,
			width = $6, height = $7, original_bytes = $8, processed_bytes = $9, processing = 'ready'
			WHERE id = $1 AND storage_key = $10`,
			j.mediaID, videoURL, w.Files.PublicURL(pKey), dur, info.HasAudio, vi.Width, vi.Height, origBytes, vBytes, j.key)
		return err
	}

	out := filepath.Join(dir, "photo.jpg")
	if err := ResizePhoto(tctx, in, out); err != nil {
		return err
	}
	pi, err := Probe(ctx, out)
	if err != nil {
		return err
	}
	pKey := base + "-1600.jpg"
	pBytes, err := w.upload(ctx, out, pKey)
	if err != nil {
		return err
	}
	_, err = w.Pool.Exec(ctx, `UPDATE product_media SET image_url = $2, width = $3, height = $4, original_bytes = $5, processed_bytes = $6,
		processing = 'ready' WHERE id = $1 AND storage_key = $7`,
		j.mediaID, w.Files.PublicURL(pKey), pi.Width, pi.Height, origBytes, pBytes, j.key)
	return err
}

func (w *Worker) download(ctx context.Context, key, dst string) (int64, error) {
	r, err := w.Files.Open(ctx, key)
	if err != nil {
		return 0, err
	}
	defer r.Close()
	f, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return io.Copy(f, r)
}

func (w *Worker) upload(ctx context.Context, src, key string) (int64, error) {
	f, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return w.Files.Put(ctx, key, f)
}

// CheckVideo untuk storage.Uploader: tolak video > 3 menit atau yang tidak bisa dibaca, langsung saat upload.
func CheckVideo(files storage.Local) func(ctx context.Context, key string) error {
	return func(ctx context.Context, key string) error {
		p, err := files.Path(key)
		if err != nil {
			return err
		}
		pctx, cancel := context.WithTimeout(ctx, 20*time.Second)
		defer cancel()
		info, err := Probe(pctx, p)
		if err != nil {
			return storage.ErrType
		}
		if info.Duration > MaxVideoSeconds+1 {
			return ErrTooLong
		}
		return nil
	}
}

// Backfill mengantrekan media lama yang belum pernah diproses (yang sudah pernah gagal permanen dilewati).
func Backfill(ctx context.Context, pool *pgxpool.Pool) (int64, error) {
	tag, err := pool.Exec(ctx, `
		INSERT INTO media_processing_jobs (media_id, step)
		SELECT pm.id, 'transcode' FROM product_media pm
		WHERE pm.deleted_at IS NULL
		  AND ((pm.kind = 'video' AND pm.poster_url IS NULL) OR (pm.kind <> 'video' AND pm.width IS NULL))
		  AND NOT EXISTS (SELECT 1 FROM media_processing_jobs j WHERE j.media_id = pm.id AND j.status IN ('queued','running','failed'))`)
	return tag.RowsAffected(), err
}

func trunc(s string) string {
	if len(s) > 1000 {
		return s[:1000]
	}
	return s
}
