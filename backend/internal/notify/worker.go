package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"phsar/internal/jobs"
)

// Sender mengirim satu pesan Telegram (bisa dengan tombol link).
type Sender interface {
	Send(ctx context.Context, chatID int64, text, buttonText, buttonURL string) error
}

// ErrPermanent: tidak perlu dicoba lagi (user memblokir bot / chat tidak ada).
var ErrPermanent = errors.New("permanent")

// TelegramBot: klien Bot API minimal untuk sendMessage.
type TelegramBot struct {
	Token string
	Base  string // kosong = https://api.telegram.org (bisa diganti untuk test)
	HTTP  *http.Client
}

func (b *TelegramBot) Send(ctx context.Context, chatID int64, text, buttonText, buttonURL string) error {
	in := map[string]any{"chat_id": chatID, "text": text, "disable_web_page_preview": true}
	// Telegram menolak tombol ke http://localhost — tombol hanya untuk https; selain itu link ditaruh di teks.
	if strings.HasPrefix(buttonURL, "https://") {
		in["reply_markup"] = map[string]any{"inline_keyboard": [][]map[string]string{{{"text": buttonText, "url": buttonURL}}}}
	} else if buttonURL != "" {
		in["text"] = text + "\n\n" + buttonURL
	}
	body, _ := json.Marshal(in)
	base := b.Base
	if base == "" {
		base = "https://api.telegram.org"
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, base+"/bot"+b.Token+"/sendMessage", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	hc := b.HTTP
	if hc == nil {
		hc = &http.Client{Timeout: 20 * time.Second}
	}
	resp, err := hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var env struct {
		OK          bool   `json:"ok"`
		ErrorCode   int    `json:"error_code"`
		Description string `json:"description"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&env)
	if env.OK {
		return nil
	}
	err = fmt.Errorf("telegram %d: %s", env.ErrorCode, env.Description)
	if env.ErrorCode == 403 || (env.ErrorCode == 400 && strings.Contains(strings.ToLower(env.Description), "chat not found")) {
		return fmt.Errorf("%w: %v", ErrPermanent, err)
	}
	return err
}

// Worker mengirim notifikasi yang masih pending lewat Telegram dan Web Push.
type Worker struct {
	Pool    *pgxpool.Pool
	Bot     Sender  // nil = Telegram tidak dikonfigurasi → notifikasi hanya di Seller Center
	Push    *Pusher // nil = Web Push mati
	BaseURL string  // PUBLIC_BASE_URL (website pembeli, kurir)
	// SellerURL / AdminURL: alamat Seller Center & admin kalau di subdomain sendiri. Kosong = BaseURL.
	SellerURL string
	AdminURL  string
	Log       *slog.Logger
	Interval  time.Duration
}

const maxAttempts = 5

// base: alamat aplikasi web milik audience (seller.<domain>, admin.<domain>, atau domain utama).
func (w *Worker) base(aud string) string {
	switch {
	case aud == AudSeller && w.SellerURL != "":
		return w.SellerURL
	case aud == AudAdmin && w.AdminURL != "":
		return w.AdminURL
	}
	return w.BaseURL
}

func (w *Worker) Run(ctx context.Context) {
	if w.Interval == 0 {
		w.Interval = 3 * time.Second
	}
	_, _ = w.Pool.Exec(ctx, `UPDATE notifications SET telegram_status = 'pending' WHERE telegram_status = 'sending'`)
	_, _ = w.Pool.Exec(ctx, `UPDATE notifications SET push_status = 'pending' WHERE push_status = 'sending'`)
	t := time.NewTicker(w.Interval)
	defer t.Stop()
	for {
		if _, err := jobs.Track(ctx, w.Pool, "notify.send", w.RunOnce); err != nil && ctx.Err() == nil { // jobs:
			w.Log.Error("worker notifikasi", "err", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}
	}
}

type item struct {
	id       int64
	user     string
	aud      string
	kind     string
	params   Params
	link     string
	lang     string
	chatID   string
	attempts int
}

// RunOnce mengirim maksimal 20 notifikasi Telegram dan 20 Web Push. Mengembalikan jumlah yang diproses.
func (w *Worker) RunOnce(ctx context.Context) (int, error) {
	// Kebanyakan user tidak login dengan Telegram / belum mengaktifkan push: lewati sekaligus,
	// supaya antrean tidak tersumbat baris yang memang tidak bisa dikirim (mis. notifikasi admin massal).
	if w.Bot != nil {
		if _, err := w.Pool.Exec(ctx, `UPDATE notifications n SET telegram_status = 'skipped', telegram_error = 'user tidak login dengan Telegram'
			WHERE n.telegram_status = 'pending'
			  AND NOT EXISTS (SELECT 1 FROM auth_identities ai WHERE ai.user_id = n.user_id AND ai.provider = 'telegram')`); err != nil {
			return 0, err
		}
	}
	if _, err := w.Pool.Exec(ctx, `UPDATE notifications n SET push_status = 'skipped', push_error = 'tidak ada browser yang berlangganan'
		WHERE n.push_status = 'pending' AND NOT EXISTS (SELECT 1 FROM push_subscriptions ps WHERE ps.user_id = n.user_id AND ps.audience = n.audience)`); err != nil {
		return 0, err
	}
	n, err := w.runTelegram(ctx)
	if err != nil {
		return n, err
	}
	m, err := w.runPush(ctx)
	return n + m, err
}

func (w *Worker) runTelegram(ctx context.Context) (int, error) {
	rows, err := w.Pool.Query(ctx, `
		WITH batch AS (
			SELECT id FROM notifications
			WHERE telegram_status = 'pending' AND (next_attempt_at IS NULL OR next_attempt_at <= now())
			ORDER BY id LIMIT 20 FOR UPDATE SKIP LOCKED)
		UPDATE notifications n SET telegram_status = 'sending', telegram_attempts = n.telegram_attempts + 1
		FROM batch, users u
		WHERE n.id = batch.id AND u.id = n.user_id
		RETURNING n.id, n.user_id::text, n.audience, n.kind, n.params, coalesce(n.link, ''), u.language::text,
		          coalesce((SELECT provider_uid FROM auth_identities ai WHERE ai.user_id = n.user_id AND ai.provider = 'telegram'
		                    ORDER BY last_used_at DESC NULLS LAST LIMIT 1), ''),
		          n.telegram_attempts`)
	if err != nil {
		return 0, err
	}
	items, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (item, error) {
		var it item
		var params []byte
		err := r.Scan(&it.id, &it.user, &it.aud, &it.kind, &params, &it.link, &it.lang, &it.chatID, &it.attempts)
		_ = json.Unmarshal(params, &it.params)
		return it, err
	})
	if err != nil {
		return 0, err
	}
	for _, it := range items {
		chat, perr := strconv.ParseInt(it.chatID, 10, 64)
		switch {
		case w.Bot == nil:
			w.finish(ctx, it.id, "skipped", "telegram tidak dikonfigurasi")
			continue
		case it.chatID == "" || perr != nil:
			w.finish(ctx, it.id, "skipped", "user tidak login dengan Telegram")
			continue
		}
		if prefs, _ := LoadPrefs(ctx, w.Pool, it.user, it.aud); !prefs.TelegramAllowed(it.kind) {
			w.finish(ctx, it.id, "skipped", "dimatikan di pengaturan notifikasi")
			continue
		}
		title, body := Render(it.kind, it.lang, it.params)
		text := title
		if body != "" {
			text += "\n" + body
		}
		url := ""
		if it.link != "" {
			url = LinkURL(w.base(it.aud), it.aud, it.link)
		}
		label := openLabels[it.aud][it.lang]
		if label == "" {
			label = openLabels[it.aud]["en"]
		}
		sctx, cancel := context.WithTimeout(ctx, 25*time.Second)
		err := w.Bot.Send(sctx, chat, text, label, url)
		cancel()
		switch {
		case err == nil:
			_, _ = w.Pool.Exec(ctx, `UPDATE notifications SET telegram_status = 'sent', telegram_sent_at = now(), telegram_error = NULL WHERE id = $1`, it.id)
		case errors.Is(err, ErrPermanent) || it.attempts >= maxAttempts:
			w.Log.Warn("notifikasi Telegram gagal (berhenti)", "id", it.id, "err", err)
			w.finish(ctx, it.id, "failed", err.Error())
		default:
			w.Log.Warn("notifikasi Telegram gagal, dicoba lagi", "id", it.id, "attempt", it.attempts, "err", err)
			_, _ = w.Pool.Exec(ctx, `UPDATE notifications SET telegram_status = 'pending', telegram_error = $2,
				next_attempt_at = now() + make_interval(mins => telegram_attempts * telegram_attempts) WHERE id = $1`, it.id, truncate(err.Error(), 500))
		}
	}
	return len(items), nil
}

func (w *Worker) finish(ctx context.Context, id int64, status, reason string) {
	_, _ = w.Pool.Exec(ctx, `UPDATE notifications SET telegram_status = $2, telegram_error = $3 WHERE id = $1`, id, status, truncate(reason, 500))
}

var openLabels = map[string]map[string]string{
	AudSeller:  openLabel,
	AudAdmin:   {"en": "Open admin", "km": "បើកផ្ទាំងគ្រប់គ្រង", "zh": "打开管理后台"},
	AudBuyer:   {"en": "Open Phsar", "km": "បើក Phsar", "zh": "打开 Phsar"},
	AudCourier: {"en": "Open courier app", "km": "បើកកម្មវិធីអ្នកដឹក", "zh": "打开快递员页面"},
}

// pushPayload: isi pesan yang dibaca service worker (webui/sw.js).
type pushPayload struct {
	ID     int64  `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	URL    string `json:"url"`
	Tag    string `json:"tag"`
	Sound  string `json:"sound"`  // order · alert · chat · ding · none
	Silent bool   `json:"silent"` // jam tenang / suara grup mati
	Test   bool   `json:"test,omitempty"`
}

// runPush: kirim Web Push ke semua browser yang berlangganan untuk audience notifikasi itu.
func (w *Worker) runPush(ctx context.Context) (int, error) {
	rows, err := w.Pool.Query(ctx, `
		WITH batch AS (
			SELECT id FROM notifications
			WHERE push_status = 'pending' AND (push_next_at IS NULL OR push_next_at <= now())
			ORDER BY id LIMIT 20 FOR UPDATE SKIP LOCKED)
		UPDATE notifications n SET push_status = 'sending', push_attempts = n.push_attempts + 1
		FROM batch, users u
		WHERE n.id = batch.id AND u.id = n.user_id
		RETURNING n.id, n.user_id::text, n.audience, n.kind, n.params, coalesce(n.link, ''), u.language::text, n.push_attempts`)
	if err != nil {
		return 0, err
	}
	items, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (item, error) {
		var it item
		var params []byte
		err := r.Scan(&it.id, &it.user, &it.aud, &it.kind, &params, &it.link, &it.lang, &it.attempts)
		_ = json.Unmarshal(params, &it.params)
		return it, err
	})
	if err != nil {
		return 0, err
	}
	for _, it := range items {
		w.pushOne(ctx, it)
	}
	return len(items), nil
}

func (w *Worker) pushOne(ctx context.Context, it item) {
	done := func(status, reason string) {
		_, _ = w.Pool.Exec(ctx, `UPDATE notifications SET push_status = $2, push_error = nullif($3, '') WHERE id = $1`, it.id, status, truncate(reason, 500))
	}
	if w.Push == nil {
		done("skipped", "web push mati")
		return
	}
	prefs, _ := LoadPrefs(ctx, w.Pool, it.user, it.aud)
	if !prefs.PushAllowed(it.kind) {
		done("skipped", "dimatikan di pengaturan notifikasi")
		return
	}
	rows, err := w.Pool.Query(ctx, `SELECT id::text, endpoint, p256dh, auth FROM push_subscriptions WHERE user_id = $1 AND audience = $2`, it.user, it.aud)
	if err != nil {
		done("pending", err.Error())
		return
	}
	type subRow struct {
		id string
		Subscription
	}
	subs, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (subRow, error) {
		var s subRow
		return s, r.Scan(&s.id, &s.Endpoint, &s.P256dh, &s.Auth)
	})
	if err != nil || len(subs) == 0 {
		done("skipped", "tidak ada browser yang berlangganan")
		return
	}
	title, body := Render(it.kind, it.lang, it.params)
	sound, _ := prefs.SoundFor(it.kind, time.Now())
	payload, _ := json.Marshal(pushPayload{ID: it.id, Title: title, Body: body, URL: LinkURL(w.base(it.aud), it.aud, it.link),
		Tag: fmt.Sprintf("%s-%d", it.kind, it.id), Sound: sound, Silent: sound == "none", Test: it.kind == TestKind})
	urgency := "normal"
	if s, _ := Sound(it.kind); s == "order" || s == "alert" {
		urgency = "high"
	}
	sent, retry := 0, false
	var lastErr error
	for _, s := range subs {
		sctx, cancel := context.WithTimeout(ctx, 20*time.Second)
		err := w.Push.Send(sctx, s.Subscription, payload, urgency, 24*time.Hour)
		cancel()
		switch {
		case err == nil:
			sent++
			_, _ = w.Pool.Exec(ctx, `UPDATE push_subscriptions SET last_success_at = now(), failures = 0 WHERE id = $1`, s.id)
		case errors.Is(err, ErrGone):
			_, _ = w.Pool.Exec(ctx, `DELETE FROM push_subscriptions WHERE id = $1`, s.id)
		case errors.Is(err, ErrPermanent):
			lastErr = err
			_, _ = w.Pool.Exec(ctx, `UPDATE push_subscriptions SET failures = failures + 1 WHERE id = $1`, s.id)
			_, _ = w.Pool.Exec(ctx, `DELETE FROM push_subscriptions WHERE id = $1 AND failures >= 5`, s.id)
		default:
			lastErr, retry = err, true
		}
	}
	switch {
	case sent > 0:
		done("sent", "")
	case retry && it.attempts < maxAttempts:
		w.Log.Warn("web push gagal, dicoba lagi", "id", it.id, "attempt", it.attempts, "err", lastErr)
		_, _ = w.Pool.Exec(ctx, `UPDATE notifications SET push_status = 'pending', push_error = $2,
			push_next_at = now() + make_interval(secs => 20 * push_attempts * push_attempts) WHERE id = $1`, it.id, truncate(lastErr.Error(), 500))
	case lastErr != nil:
		done("failed", lastErr.Error())
	default:
		done("skipped", "langganan sudah tidak berlaku")
	}
}
