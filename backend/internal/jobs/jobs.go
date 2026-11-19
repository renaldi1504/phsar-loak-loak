// Package jobs: daftar pekerjaan latar (worker) + pencatatan run.
//
//	Track membungkus satu putaran worker: lewati kalau dijeda super admin, catat durasi/jumlah/error.
//	Putaran kosong (0 item, tanpa error) tidak ditulis ke job_runs — cukup detak di job_state (maks. tiap 30 detik),
//	supaya worker yang berputar tiap 3 detik tidak membanjiri database. Run manual selalu dicatat.
//	Riwayat lebih dari 30 hari dihapus otomatis.
package jobs

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Def: satu pekerjaan latar. Teks (3 bahasa) ada di halaman admin, kunci = Name.
type Def struct {
	Name  string `json:"name"`
	Every int    `json:"every_sec"` // jeda putaran worker
	Safe  bool   `json:"safe"`      // boleh dijalankan manual dari admin ("Run now")
}

// Defs: urutan tampil di admin.
var Defs = []Def{
	{"checkout.expire", 20, true},      // checkout tidak dibayar → batal, stok kembali
	{"order.accept_timeout", 20, true}, // seller tidak menerima tepat waktu → batal + refund
	{"order.pack_timeout", 20, true},   // pesanan diterima tapi tidak dikemas tepat waktu → batal + refund
	{"payout.complete", 300, true},     // pesanan lewat masa tahan → saldo seller
	{"payout.schedule", 300, false},    // payout terjadwal (harian/mingguan) — manual lewat halaman Payout
	{"returns.escalate", 300, true},    // retur tidak dijawab seller 48 jam → admin
	{"notify.send", 3, true},           // kirim notifikasi Telegram
	{"translate.queue", 3, true},       // terjemahan AI
	{"media.queue", 3, true},           // kompres video + thumbnail
	{"photo.index", 20, true},          // vektor foto untuk cari-dengan-foto
	// courierreturn: paket kembali ke toko tidak dikonfirmasi seller 24 jam → dianggap diterima
	{"courierreturn.auto_receive", 300, true},
}

func Find(name string) (Def, bool) {
	for _, d := range Defs {
		if d.Name == name {
			return d, true
		}
	}
	return Def{}, false
}

// Func: satu putaran. Mengembalikan jumlah item yang diproses.
type Func func(ctx context.Context) (int, error)

var (
	mu      sync.RWMutex
	runners = map[string]Func{}
)

// SetRunner mendaftarkan fungsi untuk "Run now" di admin (proses API). Aman dipanggil berulang.
func SetRunner(name string, fn Func) {
	mu.Lock()
	runners[name] = fn
	mu.Unlock()
}

func Runner(name string) (Func, bool) {
	mu.RLock()
	defer mu.RUnlock()
	fn, ok := runners[name]
	return fn, ok
}

// ---------------------------------------------------------------------------
// Jeda (platform_settings['jobs'] = {"paused": ["notify.send", …]})
// ---------------------------------------------------------------------------

var pauseCache = struct {
	sync.Mutex
	at     time.Time
	paused map[string]bool
}{}

// CacheTTL: berapa lama status jeda disimpan di memori (tes boleh mengecilkan).
var CacheTTL = 15 * time.Second

func Paused(ctx context.Context, pool *pgxpool.Pool) (map[string]bool, error) {
	pauseCache.Lock()
	defer pauseCache.Unlock()
	if pauseCache.paused != nil && time.Since(pauseCache.at) < CacheTTL {
		return pauseCache.paused, nil
	}
	var raw []byte
	err := pool.QueryRow(ctx, `SELECT value FROM platform_settings WHERE key = 'jobs'`).Scan(&raw)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	var v struct {
		Paused []string `json:"paused"`
	}
	_ = json.Unmarshal(raw, &v)
	m := map[string]bool{}
	for _, n := range v.Paused {
		m[n] = true
	}
	pauseCache.paused, pauseCache.at = m, time.Now()
	return m, nil
}

// SetPaused menyimpan jeda (dipanggil admin di dalam transaksinya, supaya tercatat di audit).
func SetPaused(ctx context.Context, tx pgx.Tx, actorID, name string, paused bool) error {
	var raw []byte
	err := tx.QueryRow(ctx, `SELECT value FROM platform_settings WHERE key = 'jobs' FOR UPDATE`).Scan(&raw)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	var v struct {
		Paused []string `json:"paused"`
	}
	_ = json.Unmarshal(raw, &v)
	out := []string{}
	for _, n := range v.Paused {
		if n != name {
			out = append(out, n)
		}
	}
	if paused {
		out = append(out, name)
	}
	b, _ := json.Marshal(map[string]any{"paused": out})
	if _, err := tx.Exec(ctx, `INSERT INTO platform_settings (key, value, updated_by, updated_at) VALUES ('jobs', $1, $2, now())
		ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_by = EXCLUDED.updated_by, updated_at = now()`, b, actorID); err != nil {
		return err
	}
	pauseCache.Lock()
	pauseCache.paused = nil // baca ulang di putaran berikutnya
	pauseCache.Unlock()
	return nil
}

// ---------------------------------------------------------------------------
// Pencatatan
// ---------------------------------------------------------------------------

var beat = struct {
	sync.Mutex
	last  map[string]time.Time
	prune time.Time
}{last: map[string]time.Time{}}

// BeatEvery: detak kosong paling sering ditulis sekali per selang ini per job.
var BeatEvery = 30 * time.Second

// Track: putaran otomatis dari worker. Dilewati (0, nil) kalau job dijeda.
func Track(ctx context.Context, pool *pgxpool.Pool, name string, fn Func) (int, error) {
	if p, err := Paused(ctx, pool); err == nil && p[name] {
		return 0, nil
	}
	return run(ctx, pool, name, "auto", "", fn)
}

// Drain: worker antrean yang memproses satu item per panggilan (false = antrean kosong) — satu putaran = kosongkan antrean.
func Drain(ctx context.Context, pool *pgxpool.Pool, name string, once func(ctx context.Context) (bool, error)) (int, error) {
	return Track(ctx, pool, name, func(ctx context.Context) (int, error) {
		n := 0
		for ctx.Err() == nil {
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
	})
}

// RunNow: dijalankan super admin (juga saat dijeda), selalu dicatat.
func RunNow(ctx context.Context, pool *pgxpool.Pool, name, adminID string) (int, error) {
	fn, ok := Runner(name)
	if !ok {
		return 0, ErrNoRunner
	}
	return run(ctx, pool, name, "manual", adminID, fn)
}

var ErrNoRunner = errors.New("job ini tidak bisa dijalankan dari sini")

func run(ctx context.Context, pool *pgxpool.Pool, name, trigger, by string, fn Func) (int, error) {
	start := time.Now()
	n, err := fn(ctx)
	if ctx.Err() != nil { // server sedang berhenti: jangan catat
		return n, err
	}
	ms := int(time.Since(start) / time.Millisecond)
	now := time.Now()
	meaningful := n > 0 || err != nil || trigger == "manual"

	beat.Lock()
	due := now.Sub(beat.last[name]) >= BeatEvery
	if due || meaningful {
		beat.last[name] = now
	}
	pruneDue := now.Sub(beat.prune) >= time.Hour
	if pruneDue {
		beat.prune = now
	}
	beat.Unlock()

	// Pencatatan memakai context sendiri yang pendek: kegagalan mencatat tidak boleh mengganggu worker.
	wctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if meaningful {
		status, msg := "success", ""
		if err != nil {
			status, msg = "failed", truncate(err.Error(), 1000)
		}
		_, _ = pool.Exec(wctx, `INSERT INTO job_runs (job, started_at, finished_at, status, rows_affected, error, duration_ms, trigger, triggered_by)
			VALUES ($1, $2, $3, $4, $5, nullif($6, ''), $7, $8, nullif($9, '')::uuid)`, name, start, now, status, n, msg, ms, trigger, by)
	}
	if due || meaningful {
		if err != nil {
			_, _ = pool.Exec(wctx, `INSERT INTO job_state (job, last_tick_at, last_error, last_error_at, updated_at) VALUES ($1, $2, $3, $2, now())
				ON CONFLICT (job) DO UPDATE SET last_tick_at = EXCLUDED.last_tick_at, last_error = EXCLUDED.last_error,
				last_error_at = EXCLUDED.last_error_at, updated_at = now()`, name, now, truncate(err.Error(), 1000))
		} else {
			_, _ = pool.Exec(wctx, `INSERT INTO job_state (job, last_tick_at, last_success_at, updated_at) VALUES ($1, $2, $2, now())
				ON CONFLICT (job) DO UPDATE SET last_tick_at = EXCLUDED.last_tick_at, last_success_at = EXCLUDED.last_success_at, updated_at = now()`, name, now)
		}
	}
	if pruneDue {
		_, _ = pool.Exec(wctx, `DELETE FROM job_runs WHERE started_at < now() - interval '30 days'`)
	}
	return n, err
}

func truncate(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}
	return s
}
