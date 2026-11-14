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

// Worker mengirim notifikasi yang masih pending lewat Telegram.
type Worker struct {
	Pool     *pgxpool.Pool
	Bot      Sender // nil = Telegram tidak dikonfigurasi → notifikasi hanya di Seller Center
	BaseURL  string // PUBLIC_BASE_URL, untuk link ke Seller Center
	Log      *slog.Logger
	Interval time.Duration
}

const maxAttempts = 5

func (w *Worker) Run(ctx context.Context) {
	if w.Interval == 0 {
		w.Interval = 3 * time.Second
	}
	_, _ = w.Pool.Exec(ctx, `UPDATE notifications SET telegram_status = 'pending' WHERE telegram_status = 'sending'`)
	t := time.NewTicker(w.Interval)
	defer t.Stop()
	for {
		if _, err := w.RunOnce(ctx); err != nil && ctx.Err() == nil {
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
	kind     string
	params   Params
	link     string
	lang     string
	chatID   string
	attempts int
}

// RunOnce mengirim maksimal 20 notifikasi. Mengembalikan jumlah yang diproses.
func (w *Worker) RunOnce(ctx context.Context) (int, error) {
	rows, err := w.Pool.Query(ctx, `
		WITH batch AS (
			SELECT id FROM notifications
			WHERE telegram_status = 'pending' AND (next_attempt_at IS NULL OR next_attempt_at <= now())
			ORDER BY id LIMIT 20 FOR UPDATE SKIP LOCKED)
		UPDATE notifications n SET telegram_status = 'sending', telegram_attempts = n.telegram_attempts + 1
		FROM batch, users u
		WHERE n.id = batch.id AND u.id = n.user_id
		RETURNING n.id, n.kind, n.params, coalesce(n.link, ''), u.language::text,
		          coalesce((SELECT provider_uid FROM auth_identities ai WHERE ai.user_id = n.user_id AND ai.provider = 'telegram'
		                    ORDER BY last_used_at DESC NULLS LAST LIMIT 1), ''),
		          n.telegram_attempts`)
	if err != nil {
		return 0, err
	}
	items, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (item, error) {
		var it item
		var params []byte
		err := r.Scan(&it.id, &it.kind, &params, &it.link, &it.lang, &it.chatID, &it.attempts)
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
		title, body := Render(it.kind, it.lang, it.params)
		text := title
		if body != "" {
			text += "\n" + body
		}
		url := ""
		if it.link != "" {
			url = strings.TrimRight(w.BaseURL, "/") + "/seller/" + it.link
		}
		label := openLabel[it.lang]
		if label == "" {
			label = openLabel["en"]
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
