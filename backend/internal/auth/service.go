// Package auth: login Telegram (deep link bot), SMS OTP, sesi per perangkat
// (access token 15 menit + refresh token yang dirotasi).
package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/netip"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"phsar/internal/config"
	"phsar/internal/platform/db"
)

const (
	AccessTTL = 15 * time.Minute
)

var sessionTTL = map[string]time.Duration{
	"app":        90 * 24 * time.Hour,
	"seller_web": 30 * 24 * time.Hour,
	"admin":      8 * time.Hour,
}

var (
	ErrSessionInvalid = errors.New("sesi tidak valid, silakan login lagi")
	ErrUserSuspended  = errors.New("akun ditangguhkan")
)

type Service struct {
	pool   *pgxpool.Pool
	cfg    config.Config
	log    *slog.Logger
	secret []byte
	sms    SMSSender
	tg     *TelegramClient
	now    func() time.Time
}

func NewService(pool *pgxpool.Pool, cfg config.Config, log *slog.Logger) *Service {
	s := &Service{pool: pool, cfg: cfg, log: log, secret: []byte(cfg.JWTSecret), now: time.Now}
	s.sms = LogSMSSender{log: log}
	if cfg.TelegramBotToken != "" {
		s.tg = NewTelegramClient(cfg.TelegramBotToken)
	}
	return s
}

type TokenPair struct {
	AccessToken      string    `json:"access_token"`
	AccessExpiresAt  time.Time `json:"access_expires_at"`
	RefreshToken     string    `json:"refresh_token"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
	User             UserInfo  `json:"user"`
}

type UserInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone,omitempty"`
	Language string `json:"language"`
	IsNew    bool   `json:"is_new"`
	HasStore bool   `json:"has_store"` // true = tampilkan tombol mode seller
}

type ClientInfo struct {
	DeviceName string
	IP         string
	UserAgent  string
}

func (c ClientInfo) ipOrNil() *netip.Addr {
	a, err := netip.ParseAddr(c.IP)
	if err != nil {
		return nil
	}
	return &a
}

// identity: siapa yang login, dari provider mana.
type identity struct {
	provider string // telegram · sms · facebook · apple
	uid      string
	display  string
	name     string
	phone    string
	lang     string
}

// findOrCreateUser menautkan identitas ke user yang sudah ada, atau membuat user baru.
func findOrCreateUser(ctx context.Context, tx pgx.Tx, id identity) (string, bool, error) {
	var userID, status string
	err := tx.QueryRow(ctx, `
		SELECT u.id, u.status FROM auth_identities ai JOIN users u ON u.id = ai.user_id
		WHERE ai.provider = $1 AND ai.provider_uid = $2 AND u.deleted_at IS NULL`,
		id.provider, id.uid).Scan(&userID, &status)
	if err == nil {
		if status == "suspended" {
			return "", false, ErrUserSuspended
		}
		_, err = tx.Exec(ctx, `UPDATE auth_identities SET last_used_at = now(), display = coalesce($3, display)
			WHERE provider = $1 AND provider_uid = $2`, id.provider, id.uid, nullIfEmpty(id.display))
		return userID, false, err
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return "", false, err
	}

	created := false
	// Nomor HP sudah dipakai user lain (mis. login SMS dulu, sekarang Telegram dengan nomor sama)? Tautkan.
	if id.phone != "" {
		err = tx.QueryRow(ctx, `SELECT id FROM users WHERE phone = $1 AND deleted_at IS NULL`, id.phone).Scan(&userID)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return "", false, err
		}
	}
	if userID == "" {
		lang := id.lang
		if lang == "" {
			lang = "en"
		}
		if err := tx.QueryRow(ctx, `INSERT INTO users (name, phone, language) VALUES ($1, $2, $3::lang_code) RETURNING id`,
			id.name, nullIfEmpty(id.phone), lang).Scan(&userID); err != nil {
			return "", false, err
		}
		created = true
	}
	_, err = tx.Exec(ctx, `INSERT INTO auth_identities (user_id, provider, provider_uid, display, verified_at, last_used_at)
		VALUES ($1, $2::auth_provider, $3, $4, now(), now())`, userID, id.provider, id.uid, nullIfEmpty(id.display))
	return userID, created, err
}

// issueSession membuat baris sessions baru + pasangan token.
func (s *Service) issueSession(ctx context.Context, tx pgx.Tx, userID, audience string, isNew bool, c ClientInfo) (TokenPair, error) {
	now := s.now()
	refresh := randomToken(32)
	refreshExp := now.Add(sessionTTL[audience])
	var sid string
	if err := tx.QueryRow(ctx, `
		INSERT INTO sessions (user_id, audience, refresh_token_hash, device_name, ip, user_agent, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		userID, audience, sha256Hex(refresh), nullIfEmpty(c.DeviceName), c.ipOrNil(), nullIfEmpty(c.UserAgent), refreshExp,
	).Scan(&sid); err != nil {
		return TokenPair{}, err
	}
	return s.tokenPair(ctx, tx, userID, audience, sid, refresh, refreshExp, isNew)
}

func (s *Service) tokenPair(ctx context.Context, q pgx.Tx, userID, audience, sid, refresh string, refreshExp time.Time, isNew bool) (TokenPair, error) {
	now := s.now()
	access, err := signJWT(s.secret, Claims{Subject: userID, Audience: audience, SessionID: sid, IssuedAt: now.Unix(), ExpiresAt: now.Add(AccessTTL).Unix()})
	if err != nil {
		return TokenPair{}, err
	}
	u, err := loadUser(ctx, q, userID)
	if err != nil {
		return TokenPair{}, err
	}
	u.IsNew = isNew
	return TokenPair{AccessToken: access, AccessExpiresAt: now.Add(AccessTTL), RefreshToken: refresh, RefreshExpiresAt: refreshExp, User: u}, nil
}

type queryer interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func loadUser(ctx context.Context, q queryer, userID string) (UserInfo, error) {
	var u UserInfo
	var phone *string
	err := q.QueryRow(ctx, `
		SELECT u.id, u.name, u.phone, u.language::text,
		       EXISTS (SELECT 1 FROM store_members sm JOIN stores st ON st.id = sm.store_id
		               WHERE sm.user_id = u.id AND st.deleted_at IS NULL)
		FROM users u WHERE u.id = $1 AND u.deleted_at IS NULL`, userID).Scan(&u.ID, &u.Name, &phone, &u.Language, &u.HasStore)
	if phone != nil {
		u.Phone = *phone
	}
	return u, err
}

// Refresh menukar refresh token lama dengan pasangan baru (rotasi). Token lama langsung tidak berlaku.
func (s *Service) Refresh(ctx context.Context, refresh string, c ClientInfo) (TokenPair, error) {
	var pair TokenPair
	err := db.WithActor(ctx, s.pool, "", "system", "", func(tx pgx.Tx) error {
		var sid, userID, audience string
		var exp time.Time
		err := tx.QueryRow(ctx, `
			SELECT s.id, s.user_id, s.audience, s.expires_at FROM sessions s JOIN users u ON u.id = s.user_id
			WHERE s.refresh_token_hash = $1 AND s.revoked_at IS NULL AND u.status = 'active' AND u.deleted_at IS NULL
			FOR UPDATE OF s`, sha256Hex(refresh)).Scan(&sid, &userID, &audience, &exp)
		if errors.Is(err, pgx.ErrNoRows) || (err == nil && !exp.After(s.now())) {
			return ErrSessionInvalid
		}
		if err != nil {
			return err
		}
		next := randomToken(32)
		if _, err := tx.Exec(ctx, `UPDATE sessions SET refresh_token_hash = $2, last_seen_at = now(), ip = coalesce($3, ip) WHERE id = $1`,
			sid, sha256Hex(next), c.ipOrNil()); err != nil {
			return err
		}
		pair, err = s.tokenPair(ctx, tx, userID, audience, sid, next, exp, false)
		return err
	})
	return pair, err
}

func (s *Service) Logout(ctx context.Context, refresh string) error {
	_, err := s.pool.Exec(ctx, `UPDATE sessions SET revoked_at = now() WHERE refresh_token_hash = $1 AND revoked_at IS NULL`, sha256Hex(refresh))
	return err
}

// Authenticate memeriksa access token dan memastikan sesinya belum dicabut.
func (s *Service) Authenticate(ctx context.Context, token, audience string) (Claims, error) {
	c, err := verifyJWT(s.secret, token, audience, s.now())
	if err != nil {
		return c, err
	}
	var ok bool
	if err := s.pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM sessions WHERE id = $1 AND revoked_at IS NULL AND expires_at > now())`,
		c.SessionID).Scan(&ok); err != nil {
		return c, err
	}
	if !ok {
		return c, ErrSessionInvalid
	}
	return c, nil
}

func (s *Service) Me(ctx context.Context, userID string) (UserInfo, error) {
	return loadUser(ctx, s.pool, userID)
}

func (s *Service) recordAttempt(ctx context.Context, identifier string, c ClientInfo, success bool, reason string) {
	_, err := s.pool.Exec(ctx, `INSERT INTO login_attempts (identifier, ip, success, reason) VALUES ($1, $2, $3, $4)`,
		identifier, c.ipOrNil(), success, nullIfEmpty(reason))
	if err != nil {
		s.log.Warn("gagal mencatat login_attempts", "err", err)
	}
}

func nullIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
