package notify

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"phsar/internal/auth"
	"phsar/internal/i18n"
)

// Handler: lonceng, aliran realtime, pengaturan & langganan push untuk satu audience.
// Dipasang di /v1/seller/notifications (seller), /admin/api/notifications (admin), /v1/shop/notifications (pembeli).
type Handler struct {
	Pool     *pgxpool.Pool
	Log      *slog.Logger
	Audience string  // kosong = seller (kompatibel dengan pemasangan lama)
	Hub      *Hub    // nil = aliran realtime tidak tersedia (503)
	Push     *Pusher // nil = Web Push tidak tersedia
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

func (h *Handler) aud() string {
	if h.Audience == "" {
		return AudSeller
	}
	return h.Audience
}

// Routes: semua butuh login (dipasang di belakang RequireUser).
func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Get("/", h.list)
	r.Post("/read", h.markRead)
	r.Get("/stream", h.stream)
	r.Get("/prefs", h.getPrefs)
	r.Put("/prefs", h.putPrefs)
	r.Post("/push", h.subscribe)
	r.Delete("/push", h.unsubscribe)
	r.Post("/test", h.test)
	return r
}

// GET …/notifications?limit=30 — teks dalam bahasa header X-Lang.
func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	c, _ := auth.UserFrom(r.Context())
	lang := string(i18n.From(r.Context()))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 30
	}
	rows, err := h.Pool.Query(r.Context(), `SELECT id, kind, params, coalesce(link, ''), store_id::text, read_at IS NOT NULL, created_at
		FROM notifications WHERE user_id = $1 AND audience = $2 ORDER BY id DESC LIMIT $3`, c.Subject, h.aud(), limit)
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
	writeJSON(w, http.StatusOK, map[string]any{"notifications": out, "unread": h.unread(r.Context(), c.Subject)})
}

func (h *Handler) unread(ctx context.Context, userID string) int {
	var n int
	_ = h.Pool.QueryRow(ctx, `SELECT count(*) FROM notifications WHERE user_id = $1 AND audience = $2 AND read_at IS NULL`, userID, h.aud()).Scan(&n)
	return n
}

// POST …/notifications/read {"ids":[1,2]} atau {"all":true}
func (h *Handler) markRead(w http.ResponseWriter, r *http.Request) {
	c, _ := auth.UserFrom(r.Context())
	var in struct {
		IDs []int64 `json:"ids"`
		All bool    `json:"all"`
	}
	if err := json.NewDecoder(io.LimitReader(r.Body, 16<<10)).Decode(&in); err != nil {
		badRequest(w, "bad_json", err.Error())
		return
	}
	var err error
	if in.All {
		_, err = h.Pool.Exec(r.Context(), `UPDATE notifications SET read_at = now() WHERE user_id = $1 AND audience = $2 AND read_at IS NULL`, c.Subject, h.aud())
	} else if len(in.IDs) > 0 {
		_, err = h.Pool.Exec(r.Context(), `UPDATE notifications SET read_at = now() WHERE user_id = $1 AND audience = $2 AND id = ANY($3) AND read_at IS NULL`, c.Subject, h.aud(), in.IDs)
	}
	if err != nil {
		h.fail(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// liveItem: event "notification" di aliran SSE.
type liveItem struct {
	Item
	Sound  string `json:"sound"`  // order · alert · chat · ding · none
	Repeat int    `json:"repeat"` // berapa kali bunyi diputar (pesanan baru seller: 3)
	Volume int    `json:"volume"`
	Unread int    `json:"unread"`
}

// GET …/notifications/stream — Server-Sent Events. Klien memakai fetch() supaya bisa mengirim Authorization.
// Event: hello {unread} · notification {…liveItem} · read {unread}. Komentar ": ping" tiap 25 detik.
// Aliran ditutup saat token akses kedaluwarsa; klien memperbarui token lalu tersambung lagi.
func (h *Handler) stream(w http.ResponseWriter, r *http.Request) {
	if h.Hub == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{"error": map[string]string{"code": "unavailable", "message": "realtime off"}})
		return
	}
	c, _ := auth.UserFrom(r.Context())
	lang := string(i18n.From(r.Context()))
	ctx := r.Context()
	if c.ExpiresAt > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithDeadline(ctx, time.Unix(c.ExpiresAt, 0))
		defer cancel()
	}
	rc := http.NewResponseController(w)
	_ = rc.SetWriteDeadline(time.Time{}) // tanpa batas waktu tulis untuk koneksi panjang
	w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-transform")
	w.Header().Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)

	events, unsubscribe := h.Hub.Subscribe(c.Subject)
	defer unsubscribe()
	send := func(event string, v any) bool {
		b, _ := json.Marshal(v)
		if _, err := fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, b); err != nil {
			return false
		}
		return rc.Flush() == nil
	}
	if !send("hello", map[string]int{"unread": h.unread(ctx, c.Subject)}) {
		return
	}
	ping := time.NewTicker(25 * time.Second)
	defer ping.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ping.C:
			if _, err := io.WriteString(w, ": ping\n\n"); err != nil || rc.Flush() != nil {
				return
			}
		case ev := <-events:
			if ev.Audience != h.aud() {
				continue
			}
			if ev.Ev == "read" {
				if !send("read", map[string]int{"unread": h.unread(ctx, c.Subject)}) {
					return
				}
				continue
			}
			it, err := h.load(ctx, c.Subject, ev.ID, lang)
			if err != nil {
				continue
			}
			if !send("notification", it) {
				return
			}
		}
	}
}

func (h *Handler) load(ctx context.Context, userID string, id int64, lang string) (liveItem, error) {
	var it liveItem
	var params []byte
	err := h.Pool.QueryRow(ctx, `SELECT id, kind, params, coalesce(link, ''), store_id::text, read_at IS NOT NULL, created_at
		FROM notifications WHERE id = $1 AND user_id = $2`, id, userID).
		Scan(&it.ID, &it.Kind, &params, &it.Link, &it.StoreID, &it.Read, &it.CreatedAt)
	if err != nil {
		return it, err
	}
	var p Params
	_ = json.Unmarshal(params, &p)
	it.Title, it.Body = Render(it.Kind, lang, p)
	prefs, _ := LoadPrefs(ctx, h.Pool, userID, h.aud())
	it.Sound, it.Repeat = prefs.SoundFor(it.Kind, time.Now())
	it.Volume = prefs.Volume
	it.Unread = h.unread(ctx, userID)
	return it, nil
}

// GET …/notifications/prefs — pengaturan + grup yang tersedia + kunci publik VAPID + jumlah browser berlangganan.
func (h *Handler) getPrefs(w http.ResponseWriter, r *http.Request) {
	c, _ := auth.UserFrom(r.Context())
	prefs, err := LoadPrefs(r.Context(), h.Pool, c.Subject, h.aud())
	if err != nil {
		h.fail(w, err)
		return
	}
	h.replyPrefs(w, r, c.Subject, prefs)
}

func (h *Handler) replyPrefs(w http.ResponseWriter, r *http.Request, userID string, prefs Prefs) {
	out := map[string]any{"prefs": prefs, "groups": Groups[h.aud()], "audience": h.aud()}
	if h.Push != nil {
		if key, err := h.Push.PublicKey(r.Context()); err == nil {
			out["vapid_public_key"] = key
		} else {
			h.Log.Warn("kunci VAPID", "err", err)
		}
	}
	var devices int
	_ = h.Pool.QueryRow(r.Context(), `SELECT count(*) FROM push_subscriptions WHERE user_id = $1 AND audience = $2`, userID, h.aud()).Scan(&devices)
	out["push_devices"] = devices
	var tg bool
	_ = h.Pool.QueryRow(r.Context(), `SELECT EXISTS (SELECT 1 FROM auth_identities WHERE user_id = $1 AND provider = 'telegram')`, userID).Scan(&tg)
	out["telegram_linked"] = tg
	writeJSON(w, http.StatusOK, out)
}

// PUT …/notifications/prefs {groups, quiet_from, quiet_to, volume}
func (h *Handler) putPrefs(w http.ResponseWriter, r *http.Request) {
	c, _ := auth.UserFrom(r.Context())
	var in Prefs
	if err := json.NewDecoder(io.LimitReader(r.Body, 16<<10)).Decode(&in); err != nil {
		badRequest(w, "bad_json", err.Error())
		return
	}
	clean, err := in.Clean(h.aud())
	if err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"error": map[string]string{"code": "validation", "field": "quiet_from", "message": err.Error()}})
		return
	}
	raw, _ := json.Marshal(clean)
	if _, err := h.Pool.Exec(r.Context(), `INSERT INTO notification_prefs (user_id, audience, prefs) VALUES ($1, $2, $3)
		ON CONFLICT (user_id, audience) DO UPDATE SET prefs = EXCLUDED.prefs, updated_at = now()`, c.Subject, h.aud(), raw); err != nil {
		h.fail(w, err)
		return
	}
	h.replyPrefs(w, r, c.Subject, clean)
}

// POST …/notifications/push {endpoint, keys:{p256dh, auth}} — PushSubscription.toJSON() dari browser.
func (h *Handler) subscribe(w http.ResponseWriter, r *http.Request) {
	c, _ := auth.UserFrom(r.Context())
	var in struct {
		Endpoint string `json:"endpoint"`
		Keys     struct {
			P256dh string `json:"p256dh"`
			Auth   string `json:"auth"`
		} `json:"keys"`
	}
	if err := json.NewDecoder(io.LimitReader(r.Body, 8<<10)).Decode(&in); err != nil {
		badRequest(w, "bad_json", err.Error())
		return
	}
	u, err := url.Parse(in.Endpoint)
	if err != nil || u.Scheme != "https" || u.Host == "" || len(in.Endpoint) > 1000 {
		badRequest(w, "validation", "endpoint push tidak valid")
		return
	}
	// cek kunci bisa dipakai (gagal di sini lebih baik daripada gagal diam-diam saat mengirim)
	if _, err := Encrypt([]byte("x"), Subscription{Endpoint: in.Endpoint, P256dh: in.Keys.P256dh, Auth: in.Keys.Auth}, nil, nil); err != nil {
		badRequest(w, "validation", "kunci push tidak valid")
		return
	}
	ua := r.UserAgent()
	if len(ua) > 200 {
		ua = ua[:200]
	}
	if _, err := h.Pool.Exec(r.Context(), `INSERT INTO push_subscriptions (user_id, audience, endpoint, p256dh, auth, user_agent)
		VALUES ($1, $2, $3, $4, $5, nullif($6, ''))
		ON CONFLICT (endpoint) DO UPDATE SET user_id = EXCLUDED.user_id, audience = EXCLUDED.audience, p256dh = EXCLUDED.p256dh,
		  auth = EXCLUDED.auth, user_agent = EXCLUDED.user_agent, failures = 0`,
		c.Subject, h.aud(), in.Endpoint, in.Keys.P256dh, in.Keys.Auth, ua); err != nil {
		h.fail(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DELETE …/notifications/push {endpoint}
func (h *Handler) unsubscribe(w http.ResponseWriter, r *http.Request) {
	c, _ := auth.UserFrom(r.Context())
	var in struct {
		Endpoint string `json:"endpoint"`
	}
	if err := json.NewDecoder(io.LimitReader(r.Body, 8<<10)).Decode(&in); err != nil {
		badRequest(w, "bad_json", err.Error())
		return
	}
	if _, err := h.Pool.Exec(r.Context(), `DELETE FROM push_subscriptions WHERE user_id = $1 AND endpoint = $2`, c.Subject, strings.TrimSpace(in.Endpoint)); err != nil {
		h.fail(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// POST …/notifications/test — notifikasi uji ke diri sendiri (lonceng + realtime + push), maks. 1 per 10 detik.
func (h *Handler) test(w http.ResponseWriter, r *http.Request) {
	c, _ := auth.UserFrom(r.Context())
	var recent bool
	_ = h.Pool.QueryRow(r.Context(), `SELECT EXISTS (SELECT 1 FROM notifications WHERE user_id = $1 AND kind = $2 AND created_at > now() - interval '10 seconds')`,
		c.Subject, TestKind).Scan(&recent)
	if recent {
		writeJSON(w, http.StatusTooManyRequests, map[string]any{"error": map[string]string{"code": "rate_limited", "message": "tunggu sebentar"}})
		return
	}
	if err := EnqueueTest(r.Context(), h.Pool, c.Subject, h.aud()); err != nil {
		h.fail(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) fail(w http.ResponseWriter, err error) {
	h.Log.Error("notify error", "err", err)
	writeJSON(w, http.StatusInternalServerError, map[string]any{"error": map[string]string{"code": "internal", "message": "internal error"}})
}

func badRequest(w http.ResponseWriter, code, msg string) {
	writeJSON(w, http.StatusBadRequest, map[string]any{"error": map[string]string{"code": code, "message": msg}})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
