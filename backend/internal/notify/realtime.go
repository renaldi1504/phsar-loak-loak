package notify

// Realtime: setiap notifikasi baru / "sudah dibaca" dikabarkan Postgres lewat pg_notify('phsar_notify') saat COMMIT
// (trigger di migrasi 00050). Hub mendengarkan satu koneksi LISTEN per proses server dan meneruskan ke aliran SSE
// milik user itu. Karena lewat Postgres, notifikasi yang dibuat proses worker terpisah tetap sampai.

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Event: isi pg_notify.
type Event struct {
	User     string `json:"u"`
	ID       int64  `json:"id"`
	Audience string `json:"a"`
	Ev       string `json:"ev"` // new · read
}

// Hub: pendengar LISTEN bersama untuk semua aliran SSE di proses ini.
type Hub struct {
	Pool *pgxpool.Pool
	Log  *slog.Logger

	once sync.Once
	mu   sync.Mutex
	subs map[string]map[chan Event]struct{}
}

// NewHub: hub belum mendengarkan sampai ada pelanggan pertama.
func NewHub(pool *pgxpool.Pool, log *slog.Logger) *Hub {
	return &Hub{Pool: pool, Log: log, subs: map[string]map[chan Event]struct{}{}}
}

// Subscribe: kanal event untuk satu user. Panggil cancel saat aliran selesai.
func (h *Hub) Subscribe(userID string) (<-chan Event, func()) {
	h.once.Do(func() { go h.run(context.Background()) })
	ch := make(chan Event, 32)
	h.mu.Lock()
	if h.subs[userID] == nil {
		h.subs[userID] = map[chan Event]struct{}{}
	}
	h.subs[userID][ch] = struct{}{}
	h.mu.Unlock()
	return ch, func() {
		h.mu.Lock()
		delete(h.subs[userID], ch)
		if len(h.subs[userID]) == 0 {
			delete(h.subs, userID)
		}
		h.mu.Unlock()
	}
}

func (h *Hub) publish(ev Event) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for ch := range h.subs[ev.User] {
		select {
		case ch <- ev:
		default: // pelanggan lambat: lewati (klien memuat ulang lonceng saat tersambung lagi)
		}
	}
}

// run: LISTEN terus-menerus; tersambung ulang otomatis kalau koneksi putus.
func (h *Hub) run(ctx context.Context) {
	wait := time.Second
	for ctx.Err() == nil {
		start := time.Now()
		err := h.listen(ctx)
		if time.Since(start) > time.Minute { // sempat tersambung lama → mulai lagi dari jeda pendek
			wait = time.Second
		}
		if ctx.Err() != nil {
			return
		}
		if h.Log != nil {
			h.Log.Warn("realtime notifikasi terputus, mencoba lagi", "err", err, "dalam", wait)
		}
		time.Sleep(wait)
		if wait < 30*time.Second {
			wait *= 2
		}
	}
}

// listen memakai koneksi sendiri (bukan dari pool): koneksi LISTEN dipegang selamanya, jadi tidak mengurangi
// jatah pool dan tidak menahan pool.Close() saat server berhenti.
func (h *Hub) listen(ctx context.Context) error {
	conn, err := pgx.ConnectConfig(ctx, h.Pool.Config().ConnConfig.Copy())
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())
	if _, err := conn.Exec(ctx, `LISTEN phsar_notify`); err != nil {
		return err
	}
	for {
		n, err := conn.WaitForNotification(ctx)
		if err != nil {
			return err
		}
		var ev Event
		if json.Unmarshal([]byte(n.Payload), &ev) == nil && ev.User != "" {
			h.publish(ev)
		}
	}
}
