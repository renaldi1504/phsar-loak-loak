package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// JWT HS256 minimal (tanpa library luar). Access token berumur pendek (15 menit);
// sesi panjang dipegang refresh token yang disimpan sebagai hash di tabel sessions.

var (
	ErrTokenInvalid = errors.New("token tidak valid")
	ErrTokenExpired = errors.New("token kedaluwarsa")
)

type Claims struct {
	Subject   string `json:"sub"` // users.id
	Audience  string `json:"aud"` // app · seller_web · admin
	SessionID string `json:"sid"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}

var jwtHeader = base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))

func signJWT(secret []byte, c Claims) (string, error) {
	payload, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	unsigned := jwtHeader + "." + base64.RawURLEncoding.EncodeToString(payload)
	return unsigned + "." + mac(secret, unsigned), nil
}

func verifyJWT(secret []byte, token, audience string, now time.Time) (Claims, error) {
	var c Claims
	parts := strings.Split(token, ".")
	if len(parts) != 3 || parts[0] != jwtHeader {
		return c, ErrTokenInvalid
	}
	if !hmac.Equal([]byte(parts[2]), []byte(mac(secret, parts[0]+"."+parts[1]))) {
		return c, ErrTokenInvalid
	}
	raw, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil || json.Unmarshal(raw, &c) != nil {
		return c, ErrTokenInvalid
	}
	if c.Audience != audience || c.Subject == "" {
		return c, ErrTokenInvalid
	}
	if now.Unix() >= c.ExpiresAt {
		return c, ErrTokenExpired
	}
	return c, nil
}

func mac(secret []byte, s string) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(s))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}
