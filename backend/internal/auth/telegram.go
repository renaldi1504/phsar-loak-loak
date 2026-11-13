package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"phsar/internal/platform/db"
)

const telegramLoginTTL = 10 * time.Minute

var ErrTelegramDisabled = errors.New("login Telegram belum diaktifkan di server")

// ---- Klien Bot API -------------------------------------------------------

type TelegramClient struct {
	base string
	http *http.Client
}

func NewTelegramClient(token string) *TelegramClient {
	return &TelegramClient{base: "https://api.telegram.org/bot" + token, http: &http.Client{Timeout: 40 * time.Second}}
}

type tgUser struct {
	ID           int64  `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

type TelegramUpdate struct {
	UpdateID int64 `json:"update_id"`
	Message  *struct {
		Chat struct {
			ID   int64  `json:"id"`
			Type string `json:"type"`
		} `json:"chat"`
		From *tgUser `json:"from"`
		Text string  `json:"text"`
	} `json:"message"`
}

func (c *TelegramClient) call(ctx context.Context, method string, in, out any) error {
	body, _ := json.Marshal(in)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.base+"/"+method, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var env struct {
		OK          bool            `json:"ok"`
		Description string          `json:"description"`
		Result      json.RawMessage `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		return err
	}
	if !env.OK {
		return fmt.Errorf("telegram %s: %s", method, env.Description)
	}
	if out != nil {
		return json.Unmarshal(env.Result, out)
	}
	return nil
}

func (c *TelegramClient) getUpdates(ctx context.Context, offset int64) ([]TelegramUpdate, error) {
	var ups []TelegramUpdate
	err := c.call(ctx, "getUpdates", map[string]any{"offset": offset, "timeout": 30, "allowed_updates": []string{"message"}}, &ups)
	return ups, err
}

func (c *TelegramClient) sendMessage(ctx context.Context, chatID int64, text string) error {
	return c.call(ctx, "sendMessage", map[string]any{"chat_id": chatID, "text": text}, nil)
}

// ---- Alur login ----------------------------------------------------------

type TelegramStart struct {
	RequestID string    `json:"request_id"`
	PollToken string    `json:"poll_token"` // rahasia milik app; dipakai saat polling
	DeepLink  string    `json:"deep_link"`  // buka di Telegram: https://t.me/<bot>?start=login_<nonce>
	AppLink   string    `json:"app_link"`   // tg://resolve?... langsung membuka app Telegram
	ExpiresAt time.Time `json:"expires_at"`
}

// StartTelegram membuat permintaan login dan link ke bot.
func (s *Service) StartTelegram(ctx context.Context, audience string, c ClientInfo) (TelegramStart, error) {
	if s.tg == nil {
		return TelegramStart{}, ErrTelegramDisabled
	}
	nonce := randomToken(16) // 32 hex → "login_" + 32 = 38 karakter (batas Telegram 64)
	poll := randomToken(24)
	exp := s.now().Add(telegramLoginTTL)
	var id string
	if err := s.pool.QueryRow(ctx, `
		INSERT INTO telegram_login_requests (nonce_hash, poll_token_hash, audience, device_name, ip, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		sha256Hex(nonce), sha256Hex(poll), audience, nullIfEmpty(c.DeviceName), c.ipOrNil(), exp).Scan(&id); err != nil {
		return TelegramStart{}, err
	}
	bot := s.cfg.TelegramBotUsername
	return TelegramStart{
		RequestID: id, PollToken: poll, ExpiresAt: exp,
		DeepLink: "https://t.me/" + bot + "?start=login_" + nonce,
		AppLink:  "tg://resolve?domain=" + bot + "&start=login_" + nonce,
	}, nil
}

var botReplies = map[string]map[string]string{
	"ok": {
		"en": "Signed in to Phsar. You can go back to the app now.",
		"km": "បានចូល Phsar ហើយ។ អ្នកអាចត្រឡប់ទៅកម្មវិធីវិញបាន។",
		"zh": "已登录 Phsar，现在可以返回 App。",
	},
	"expired": {
		"en": "This login link has expired. Please tap “Continue with Telegram” in the app again.",
		"km": "តំណចូលនេះផុតកំណត់ហើយ។ សូមចុច “បន្តជាមួយ Telegram” ក្នុងកម្មវិធីម្តងទៀត។",
		"zh": "此登录链接已过期，请在 App 中再次点击“使用 Telegram 继续”。",
	},
	"hello": {
		"en": "Hi! This bot is used to sign in to Phsar. Open the Phsar app and tap “Continue with Telegram”.",
		"km": "សួស្តី! បូតនេះប្រើសម្រាប់ចូល Phsar។ សូមបើកកម្មវិធី Phsar ហើយចុច “បន្តជាមួយ Telegram”។",
		"zh": "你好！此机器人用于登录 Phsar。请打开 Phsar App 并点击“使用 Telegram 继续”。",
	},
}

func reply(key, lang string) string { return botReplies[key][langFrom(lang)] }

// HandleUpdate memproses satu pesan ke bot (dari polling atau webhook).
// Mengembalikan teks balasan (kosong = tidak membalas) supaya mudah dites tanpa jaringan.
func (s *Service) HandleUpdate(ctx context.Context, u TelegramUpdate) (chatID int64, text string, err error) {
	m := u.Message
	if m == nil || m.From == nil || m.From.IsBot || m.Chat.Type != "private" {
		return 0, "", nil
	}
	chatID = m.Chat.ID
	lang := m.From.LanguageCode
	cmd := strings.TrimSpace(m.Text)
	if !strings.HasPrefix(cmd, "/start") {
		return chatID, reply("hello", lang), nil
	}
	param := strings.TrimSpace(strings.TrimPrefix(cmd, "/start"))
	if !strings.HasPrefix(param, "login_") {
		return chatID, reply("hello", lang), nil
	}
	nonce := strings.TrimPrefix(param, "login_")
	tag, err := s.pool.Exec(ctx, `
		UPDATE telegram_login_requests
		SET status = 'confirmed', confirmed_at = now(), telegram_user_id = $2, telegram_username = $3,
		    first_name = $4, last_name = $5, language_code = $6
		WHERE nonce_hash = $1 AND status = 'pending' AND expires_at > now()`,
		sha256Hex(nonce), m.From.ID, nullIfEmpty(m.From.Username), nullIfEmpty(m.From.FirstName),
		nullIfEmpty(m.From.LastName), nullIfEmpty(m.From.LanguageCode))
	if err != nil {
		return chatID, "", err
	}
	if tag.RowsAffected() == 0 {
		return chatID, reply("expired", lang), nil
	}
	return chatID, reply("ok", lang), nil
}

type TelegramPoll struct {
	Status string     `json:"status"` // pending · expired · ok
	Tokens *TokenPair `json:"tokens,omitempty"`
}

// PollTelegram dipanggil app tiap 1–2 detik sampai status bukan pending.
func (s *Service) PollTelegram(ctx context.Context, requestID, pollToken string, c ClientInfo) (TelegramPoll, error) {
	var out TelegramPoll
	err := db.WithActor(ctx, s.pool, "", "system", "", func(tx pgx.Tx) error {
		var status, audience string
		var exp time.Time
		var tgID *int64
		var username, first, last, langCode *string
		err := tx.QueryRow(ctx, `
			SELECT status, audience, expires_at, telegram_user_id, telegram_username, first_name, last_name, language_code
			FROM telegram_login_requests WHERE id = $1 AND poll_token_hash = $2 FOR UPDATE`,
			requestID, sha256Hex(pollToken)).Scan(&status, &audience, &exp, &tgID, &username, &first, &last, &langCode)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrSessionInvalid
		}
		if err != nil {
			return err
		}
		switch {
		case status == "pending" && exp.After(s.now()):
			out.Status = "pending"
			return nil
		case status != "confirmed":
			out.Status = "expired"
			return nil
		}
		name := strings.TrimSpace(deref(first) + " " + deref(last))
		display := ""
		if u := deref(username); u != "" {
			display = "@" + u
		}
		userID, created, err := findOrCreateUser(ctx, tx, identity{
			provider: "telegram", uid: strconv.FormatInt(*tgID, 10), display: display, name: name, lang: langFrom(deref(langCode)),
		})
		if err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `UPDATE telegram_login_requests SET status = 'consumed', consumed_at = now() WHERE id = $1`, requestID); err != nil {
			return err
		}
		pair, err := s.issueSession(ctx, tx, userID, audience, created, c)
		if err != nil {
			return err
		}
		out.Status, out.Tokens = "ok", &pair
		return nil
	})
	return out, err
}

// RunTelegramPolling dipakai di laptop (tanpa domain publik): ambil pesan bot lewat getUpdates.
func (s *Service) RunTelegramPolling(ctx context.Context) {
	if s.tg == nil {
		return
	}
	// Pastikan tidak ada webhook aktif, kalau ada getUpdates akan ditolak Telegram.
	_ = s.tg.call(ctx, "deleteWebhook", map[string]any{}, nil)
	s.log.Info("telegram bot polling aktif", "bot", "@"+s.cfg.TelegramBotUsername)
	for ctx.Err() == nil {
		var offset int64
		_ = s.pool.QueryRow(ctx, `SELECT last_update_id + 1 FROM telegram_bot_state WHERE id = 1`).Scan(&offset)
		ups, err := s.tg.getUpdates(ctx, offset)
		if err != nil {
			if ctx.Err() == nil {
				s.log.Warn("telegram getUpdates gagal", "err", err)
				time.Sleep(5 * time.Second)
			}
			continue
		}
		for _, u := range ups {
			s.processUpdate(ctx, u)
			_, _ = s.pool.Exec(ctx, `UPDATE telegram_bot_state SET last_update_id = $1, updated_at = now() WHERE id = 1`, u.UpdateID)
		}
	}
}

func (s *Service) processUpdate(ctx context.Context, u TelegramUpdate) {
	chatID, text, err := s.HandleUpdate(ctx, u)
	if err != nil {
		s.log.Error("telegram update gagal diproses", "err", err)
		return
	}
	if text != "" && s.tg != nil {
		if err := s.tg.sendMessage(ctx, chatID, text); err != nil {
			s.log.Warn("telegram sendMessage gagal", "err", err)
		}
	}
}

func deref(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func sprintf(f string, a ...any) string { return fmt.Sprintf(f, a...) }
