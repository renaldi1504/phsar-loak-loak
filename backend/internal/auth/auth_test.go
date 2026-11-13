package auth

import (
	"testing"
	"time"
)

func TestNormalizePhone(t *testing.T) {
	ok := map[string]string{
		"012 345 678":       "+85512345678",
		"097-123-4567":      "+855971234567",
		"+855 12 345 678":   "+85512345678",
		"85512345678":       "+85512345678",
		"+62 812 3456 7890": "+6281234567890",
	}
	for in, want := range ok {
		got, err := NormalizePhone(in)
		if err != nil || got != want {
			t.Errorf("NormalizePhone(%q) = %q, %v; mau %q", in, got, err, want)
		}
	}
	for _, bad := range []string{"", "12345", "abc", "0123"} {
		if _, err := NormalizePhone(bad); err == nil {
			t.Errorf("NormalizePhone(%q) seharusnya error", bad)
		}
	}
}

func TestJWT(t *testing.T) {
	secret := []byte("0123456789abcdef0123456789abcdef")
	now := time.Unix(1_800_000_000, 0)
	tok, err := signJWT(secret, Claims{Subject: "u1", Audience: "app", SessionID: "s1", IssuedAt: now.Unix(), ExpiresAt: now.Add(AccessTTL).Unix()})
	if err != nil {
		t.Fatal(err)
	}
	if c, err := verifyJWT(secret, tok, "app", now); err != nil || c.Subject != "u1" || c.SessionID != "s1" {
		t.Fatalf("verify gagal: %v %+v", err, c)
	}
	if _, err := verifyJWT(secret, tok, "admin", now); err != ErrTokenInvalid {
		t.Fatal("audience lain harus ditolak")
	}
	if _, err := verifyJWT(secret, tok, "app", now.Add(16*time.Minute)); err != ErrTokenExpired {
		t.Fatal("token kedaluwarsa harus ditolak")
	}
	if _, err := verifyJWT([]byte("secret-lain-0123456789abcdefghijk"), tok, "app", now); err != ErrTokenInvalid {
		t.Fatal("tanda tangan salah harus ditolak")
	}
}
