package auth

import (
	"context"
	"crypto/subtle"
	"errors"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"

	"phsar/internal/platform/db"
)

const (
	otpTTL         = 5 * time.Minute
	otpCooldown    = 60 * time.Second
	otpMaxPerHour  = 5
	otpMaxAttempts = 5
)

var (
	ErrOTPTooSoon   = errors.New("tunggu sebentar sebelum minta kode lagi")
	ErrOTPTooMany   = errors.New("terlalu banyak permintaan kode, coba lagi nanti")
	ErrOTPInvalid   = errors.New("kode salah atau sudah kedaluwarsa")
	ErrOTPExhausted = errors.New("terlalu banyak percobaan, minta kode baru")
)

// SMSSender mengirim SMS. Provider sungguhan (lokal Kamboja / Twilio) ditambahkan nanti.
type SMSSender interface {
	Send(ctx context.Context, phone, text string) error
}

// SetSMSSender mengganti pengirim SMS (dipakai test QA untuk menangkap kode OTP; nanti provider sungguhan).
func (s *Service) SetSMSSender(x SMSSender) { s.sms = x }

// LogSMSSender: hanya untuk pengembangan — kode OTP ditulis ke log server.
type LogSMSSender struct{ log *slog.Logger }

func (l LogSMSSender) Send(_ context.Context, phone, text string) error {
	l.log.Warn("SMS (mode log, tidak benar-benar dikirim)", "to", phone, "text", text)
	return nil
}

var otpText = map[string]string{
	"en": "Phsar code: %s. Valid for 5 minutes. Do not share it with anyone.",
	"km": "លេខកូដ Phsar៖ %s។ មានសុពលភាព 5 នាទី។ កុំប្រាប់អ្នកណាម្នាក់។",
	"zh": "Phsar 验证码：%s，5 分钟内有效。请勿告诉任何人。",
}

// RequestOTP mengirim kode 6 digit ke nomor HP.
func (s *Service) RequestOTP(ctx context.Context, rawPhone, lang string, c ClientInfo) (string, error) {
	phone, err := NormalizePhone(rawPhone)
	if err != nil {
		return "", err
	}
	var last *time.Time
	var lastHour int
	if err := s.pool.QueryRow(ctx, `
		SELECT max(created_at), count(*) FILTER (WHERE created_at > now() - interval '1 hour')
		FROM otp_codes WHERE phone = $1 AND purpose = 'login'`, phone).Scan(&last, &lastHour); err != nil {
		return "", err
	}
	if last != nil && s.now().Sub(*last) < otpCooldown {
		return "", ErrOTPTooSoon
	}
	if lastHour >= otpMaxPerHour {
		s.recordAttempt(ctx, phone, c, false, "otp_rate_limited")
		return "", ErrOTPTooMany
	}
	code := randomDigits(6)
	if _, err := s.pool.Exec(ctx, `
		INSERT INTO otp_codes (phone, code_hash, purpose, channel, expires_at) VALUES ($1, $2, 'login', 'sms', $3)`,
		phone, s.otpHash(phone, code), s.now().Add(otpTTL)); err != nil {
		return "", err
	}
	tmpl, ok := otpText[lang]
	if !ok {
		tmpl = otpText["en"]
	}
	if err := s.sms.Send(ctx, phone, sprintf(tmpl, code)); err != nil {
		return "", err
	}
	return phone, nil
}

// VerifyOTP memeriksa kode lalu membuat sesi. User baru dibuat otomatis.
func (s *Service) VerifyOTP(ctx context.Context, rawPhone, code, audience, lang string, c ClientInfo) (TokenPair, error) {
	phone, err := NormalizePhone(rawPhone)
	if err != nil {
		return TokenPair{}, err
	}
	var pair TokenPair
	var failReason string
	err = db.WithActor(ctx, s.pool, "", "system", "", func(tx pgx.Tx) error {
		var id, hash string
		var attempts int
		err := tx.QueryRow(ctx, `
			SELECT id, code_hash, attempts FROM otp_codes
			WHERE phone = $1 AND purpose = 'login' AND used_at IS NULL AND expires_at > now()
			ORDER BY created_at DESC LIMIT 1 FOR UPDATE`, phone).Scan(&id, &hash, &attempts)
		if errors.Is(err, pgx.ErrNoRows) {
			failReason = "otp_missing"
			return ErrOTPInvalid
		}
		if err != nil {
			return err
		}
		if attempts >= otpMaxAttempts {
			failReason = "otp_exhausted"
			return ErrOTPExhausted
		}
		if subtle.ConstantTimeCompare([]byte(hash), []byte(s.otpHash(phone, code))) != 1 {
			// Catat percobaan salah di transaksi yang sama (baris ini sedang dikunci FOR UPDATE), lalu commit.
			// Dulu memakai s.pool di luar tx → menunggu kunci miliknya sendiri sampai timeout 30 detik.
			failReason = "otp_wrong"
			_, err := tx.Exec(ctx, `UPDATE otp_codes SET attempts = attempts + 1 WHERE id = $1`, id)
			return err
		}
		if _, err := tx.Exec(ctx, `UPDATE otp_codes SET used_at = now() WHERE id = $1`, id); err != nil {
			return err
		}
		userID, created, err := findOrCreateUser(ctx, tx, identity{provider: "sms", uid: phone, phone: phone, lang: lang})
		if err != nil {
			return err
		}
		pair, err = s.issueSession(ctx, tx, userID, audience, created, c)
		return err
	})
	if err == nil && failReason == "otp_wrong" {
		err = ErrOTPInvalid
	}
	s.recordAttempt(ctx, phone, c, err == nil, failReason)
	return pair, err
}

func (s *Service) otpHash(phone, code string) string {
	return sha256Hex("otp", phone, code, string(s.secret))
}
