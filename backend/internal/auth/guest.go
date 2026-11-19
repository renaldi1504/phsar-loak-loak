package auth

// Checkout sebagai tamu: akun sementara tanpa nomor HP / Telegram / Facebook, hanya hidup di browser ini (cookie sesi).
// Tamu bisa belanja & melihat pesanannya sendiri, tapi tidak bisa membuka toko atau memakai API seller/kurir.

import (
	"context"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"

	"phsar/internal/platform/db"
)

var (
	ErrGuestLimit      = errors.New("terlalu banyak akun tamu dari jaringan ini, coba lagi nanti atau login")
	ErrGuestNotAllowed = errors.New("akun tamu tidak bisa memakai fitur ini — silakan login")
)

const guestPerIPHour = 30

// GuestSession membuat akun tamu baru + sesi (audience app).
func (s *Service) GuestSession(ctx context.Context, lang string, c ClientInfo) (TokenPair, error) {
	var n int
	if err := s.pool.QueryRow(ctx, `SELECT count(*) FROM login_attempts WHERE identifier = 'guest' AND ip = $1
		AND occurred_at > now() - interval '1 hour'`, c.ipOrNil()).Scan(&n); err != nil {
		return TokenPair{}, err
	}
	if n >= guestPerIPHour {
		return TokenPair{}, ErrGuestLimit
	}
	var pair TokenPair
	err := db.WithActor(ctx, s.pool, "", "system", "", func(tx pgx.Tx) error {
		var userID string
		if err := tx.QueryRow(ctx, `INSERT INTO users (name, language, is_guest) VALUES ('', $1::lang_code, true) RETURNING id`, langFrom(lang)).Scan(&userID); err != nil {
			return err
		}
		var err error
		pair, err = s.issueSession(ctx, tx, userID, "app", true, c)
		return err
	})
	s.recordAttempt(ctx, "guest", c, err == nil, "")
	return pair, err
}

// NoGuest: pasang setelah RequireUser untuk fitur yang butuh akun sungguhan (seller, kurir).
func NoGuest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if c, ok := UserFrom(r.Context()); ok && c.Guest {
			writeErr(w, http.StatusForbidden, "guest_not_allowed", ErrGuestNotAllowed.Error())
			return
		}
		next.ServeHTTP(w, r)
	})
}
