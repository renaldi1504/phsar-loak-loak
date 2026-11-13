package auth

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"phsar/internal/i18n"
)

// Routes untuk /v1/app/auth. audience = "app" (Flutter) atau "seller_web".
func (s *Service) Routes(audience string) http.Handler {
	r := chi.NewRouter()
	r.Post("/otp/request", s.handleOTPRequest)
	r.Post("/otp/verify", func(w http.ResponseWriter, r *http.Request) { s.handleOTPVerify(w, r, audience) })
	r.Post("/telegram/start", func(w http.ResponseWriter, r *http.Request) { s.handleTelegramStart(w, r, audience) })
	r.Post("/telegram/poll", s.handleTelegramPoll)
	r.Post("/refresh", s.handleRefresh)
	r.Post("/logout", s.handleLogout)
	return r
}

// TelegramWebhook dipakai di server (TELEGRAM_MODE=webhook): POST /webhooks/telegram/{key}.
func (s *Service) TelegramWebhook(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	if s.cfg.TelegramWebhookKey == "" || subtle.ConstantTimeCompare([]byte(key), []byte(s.cfg.TelegramWebhookKey)) != 1 {
		http.NotFound(w, r)
		return
	}
	var u TelegramUpdate
	if err := json.NewDecoder(io.LimitReader(r.Body, 1<<20)).Decode(&u); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	s.processUpdate(r.Context(), u)
	w.WriteHeader(http.StatusOK)
}

type ctxKey struct{}

// RequireUser memeriksa header Authorization: Bearer <access_token>.
func (s *Service) RequireUser(audience string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tok := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			if tok == "" || tok == r.Header.Get("Authorization") {
				writeErr(w, http.StatusUnauthorized, "unauthenticated", "login required")
				return
			}
			c, err := s.Authenticate(r.Context(), tok, audience)
			if err != nil {
				code := "token_invalid"
				switch {
				case errors.Is(err, ErrTokenExpired):
					code = "token_expired" // app memanggil /auth/refresh lalu mengulang request
				case errors.Is(err, ErrSessionInvalid):
					code = "session_invalid" // sudah logout / dicabut: kembali ke layar login
				}
				writeErr(w, http.StatusUnauthorized, code, err.Error())
				return
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKey{}, c)))
		})
	}
}

// UserFrom mengambil klaim user yang sedang login dari context.
func UserFrom(ctx context.Context) (Claims, bool) {
	c, ok := ctx.Value(ctxKey{}).(Claims)
	return c, ok
}

func (s *Service) HandleMe(w http.ResponseWriter, r *http.Request) {
	c, _ := UserFrom(r.Context())
	u, err := s.Me(r.Context(), c.Subject)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "internal", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, u)
}

func (s *Service) handleOTPRequest(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Phone string `json:"phone"`
	}
	if !decode(w, r, &in) {
		return
	}
	phone, err := s.RequestOTP(r.Context(), in.Phone, string(i18n.From(r.Context())), client(r, ""))
	if err != nil {
		s.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"phone": phone, "expires_in": int(otpTTL.Seconds()), "resend_after": int(otpCooldown.Seconds())})
}

func (s *Service) handleOTPVerify(w http.ResponseWriter, r *http.Request, audience string) {
	var in struct {
		Phone      string `json:"phone"`
		Code       string `json:"code"`
		DeviceName string `json:"device_name"`
	}
	if !decode(w, r, &in) {
		return
	}
	pair, err := s.VerifyOTP(r.Context(), in.Phone, strings.TrimSpace(in.Code), audience, string(i18n.From(r.Context())), client(r, in.DeviceName))
	if err != nil {
		s.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, pair)
}

func (s *Service) handleTelegramStart(w http.ResponseWriter, r *http.Request, audience string) {
	var in struct {
		DeviceName string `json:"device_name"`
	}
	_ = json.NewDecoder(io.LimitReader(r.Body, 4096)).Decode(&in)
	st, err := s.StartTelegram(r.Context(), audience, client(r, in.DeviceName))
	if err != nil {
		s.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, st)
}

func (s *Service) handleTelegramPoll(w http.ResponseWriter, r *http.Request) {
	var in struct {
		RequestID string `json:"request_id"`
		PollToken string `json:"poll_token"`
	}
	if !decode(w, r, &in) {
		return
	}
	res, err := s.PollTelegram(r.Context(), in.RequestID, in.PollToken, client(r, ""))
	if err != nil {
		s.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (s *Service) handleRefresh(w http.ResponseWriter, r *http.Request) {
	var in struct {
		RefreshToken string `json:"refresh_token"`
	}
	if !decode(w, r, &in) {
		return
	}
	pair, err := s.Refresh(r.Context(), in.RefreshToken, client(r, ""))
	if err != nil {
		s.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, pair)
}

func (s *Service) handleLogout(w http.ResponseWriter, r *http.Request) {
	var in struct {
		RefreshToken string `json:"refresh_token"`
	}
	if !decode(w, r, &in) {
		return
	}
	if err := s.Logout(r.Context(), in.RefreshToken); err != nil {
		s.fail(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// fail memetakan error domain ke kode HTTP + kode error yang diterjemahkan app.
func (s *Service) fail(w http.ResponseWriter, err error) {
	type m struct {
		status int
		code   string
	}
	table := map[error]m{
		ErrBadPhone:         {http.StatusBadRequest, "phone_invalid"},
		ErrOTPTooSoon:       {http.StatusTooManyRequests, "otp_too_soon"},
		ErrOTPTooMany:       {http.StatusTooManyRequests, "otp_rate_limited"},
		ErrOTPInvalid:       {http.StatusUnauthorized, "otp_invalid"},
		ErrOTPExhausted:     {http.StatusUnauthorized, "otp_exhausted"},
		ErrSessionInvalid:   {http.StatusUnauthorized, "session_invalid"},
		ErrUserSuspended:    {http.StatusForbidden, "user_suspended"},
		ErrTelegramDisabled: {http.StatusServiceUnavailable, "telegram_disabled"},
	}
	for e, v := range table {
		if errors.Is(err, e) {
			writeErr(w, v.status, v.code, err.Error())
			return
		}
	}
	s.log.Error("auth error", "err", err)
	writeErr(w, http.StatusInternalServerError, "internal", "internal error")
}

func client(r *http.Request, device string) ClientInfo {
	ip := r.RemoteAddr
	if h, _, err := net.SplitHostPort(ip); err == nil {
		ip = h
	}
	return ClientInfo{DeviceName: device, IP: ip, UserAgent: r.UserAgent()}
}

func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	if err := json.NewDecoder(io.LimitReader(r.Body, 64<<10)).Decode(v); err != nil {
		writeErr(w, http.StatusBadRequest, "bad_json", "invalid JSON body")
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, code, msg string) {
	writeJSON(w, status, map[string]any{"error": map[string]string{"code": code, "message": msg}})
}
