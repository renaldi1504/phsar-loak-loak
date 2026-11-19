package auth

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"phsar/internal/i18n"
)

// CookieRoutes: login untuk halaman web (Seller Center, admin).
// Refresh token disimpan di cookie HttpOnly + SameSite=Strict (tidak bisa dibaca JavaScript);
// access token (15 menit) dikirim di JSON dan hanya disimpan di memori halaman.
func (s *Service) CookieRoutes(audience, cookieName, cookiePath string) http.Handler {
	r := chi.NewRouter()
	set := func(w http.ResponseWriter, p TokenPair) {
		http.SetCookie(w, &http.Cookie{
			Name: cookieName, Value: p.RefreshToken, Path: cookiePath,
			Expires: p.RefreshExpiresAt, HttpOnly: true, SameSite: http.SameSiteStrictMode,
			Secure: strings.HasPrefix(s.cfg.PublicBaseURL, "https://"),
		})
	}
	clear := func(w http.ResponseWriter) {
		http.SetCookie(w, &http.Cookie{Name: cookieName, Value: "", Path: cookiePath, MaxAge: -1, HttpOnly: true, SameSite: http.SameSiteStrictMode})
	}
	public := func(p TokenPair) map[string]any {
		out := map[string]any{"access_token": p.AccessToken, "access_expires_at": p.AccessExpiresAt, "user": p.User}
		if p.MFA != nil {
			out["mfa"] = p.MFA
		}
		return out
	}
	// Cegah CSRF: request dari halaman kita selalu mengirim header ini; form dari situs lain tidak bisa.
	guard := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-Phsar") != "1" {
				writeErr(w, http.StatusForbidden, "csrf", "missing X-Phsar header")
				return
			}
			next(w, r)
		}
	}

	r.Post("/otp/request", guard(s.handleOTPRequest))
	r.Post("/otp/verify", guard(func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			Phone string `json:"phone"`
			Code  string `json:"code"`
		}
		if !decode(w, r, &in) {
			return
		}
		pair, err := s.VerifyOTP(r.Context(), in.Phone, strings.TrimSpace(in.Code), audience, string(i18n.From(r.Context())), client(r, "web"))
		if err != nil {
			s.fail(w, err)
			return
		}
		set(w, pair)
		writeJSON(w, http.StatusOK, public(pair))
	}))
	r.Post("/telegram/start", guard(func(w http.ResponseWriter, r *http.Request) {
		s.handleTelegramStart(w, r, audience)
	}))
	r.Post("/telegram/poll", guard(func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			RequestID string `json:"request_id"`
			PollToken string `json:"poll_token"`
		}
		if !decode(w, r, &in) {
			return
		}
		res, err := s.PollTelegram(r.Context(), in.RequestID, in.PollToken, client(r, "web"))
		if err != nil {
			s.fail(w, err)
			return
		}
		if res.Tokens != nil {
			set(w, *res.Tokens)
			writeJSON(w, http.StatusOK, map[string]any{"status": res.Status, "session": public(*res.Tokens)})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"status": res.Status})
	}))
	r.Post("/refresh", guard(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(cookieName)
		if err != nil || c.Value == "" {
			writeErr(w, http.StatusUnauthorized, "session_invalid", "login required")
			return
		}
		pair, err := s.Refresh(r.Context(), c.Value, client(r, "web"))
		if err != nil {
			clear(w)
			s.fail(w, err)
			return
		}
		set(w, pair)
		writeJSON(w, http.StatusOK, public(pair))
	}))
	if audience == "admin" {
		r.Mount("/mfa", s.MFARoutes())
	}
	// Login Facebook (redirect, bukan fetch — tidak bisa membawa header X-Phsar; dilindungi state acak di cookie).
	r.Get("/facebook/start", func(w http.ResponseWriter, r *http.Request) { s.facebookStart(w, r, cookiePath) })
	r.Get("/facebook/callback", func(w http.ResponseWriter, r *http.Request) { s.facebookCallback(w, r, cookiePath, audience, set) })
	// Checkout sebagai tamu — hanya website pembeli.
	if cookiePath == "/shop/auth" {
		r.Post("/guest", guard(func(w http.ResponseWriter, r *http.Request) {
			pair, err := s.GuestSession(r.Context(), string(i18n.From(r.Context())), client(r, "web"))
			if err != nil {
				s.fail(w, err)
				return
			}
			set(w, pair)
			writeJSON(w, http.StatusOK, public(pair))
		}))
	}
	r.Post("/logout", guard(func(w http.ResponseWriter, r *http.Request) {
		if c, err := r.Cookie(cookieName); err == nil && c.Value != "" {
			_ = s.Logout(r.Context(), c.Value)
		}
		clear(w)
		w.WriteHeader(http.StatusNoContent)
	}))
	return r
}
