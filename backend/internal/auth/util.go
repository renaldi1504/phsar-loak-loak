package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math/big"
	"strings"
	"unicode"
)

var ErrBadPhone = errors.New("nomor telepon tidak valid")

// NormalizePhone mengubah nomor ke format E.164. Nomor lokal Kamboja ("012 345 678")
// menjadi +85512345678. Nomor dengan + dari negara lain diterima apa adanya.
func NormalizePhone(in string) (string, error) {
	plus := strings.HasPrefix(strings.TrimSpace(in), "+")
	var b strings.Builder
	for _, r := range in {
		if unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	d := b.String()
	switch {
	case plus:
	case strings.HasPrefix(d, "855"):
	case strings.HasPrefix(d, "0") && len(d) >= 9 && len(d) <= 10:
		d = "855" + d[1:]
	default:
		return "", ErrBadPhone
	}
	if len(d) < 8 || len(d) > 15 || d[0] == '0' {
		return "", ErrBadPhone
	}
	if strings.HasPrefix(d, "855") && (len(d) < 11 || len(d) > 12) {
		return "", ErrBadPhone
	}
	return "+" + d, nil
}

func randomToken(nBytes int) string {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

func randomDigits(n int) string {
	var b strings.Builder
	for i := 0; i < n; i++ {
		v, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			panic(err)
		}
		b.WriteByte(byte('0' + v.Int64()))
	}
	return b.String()
}

func sha256Hex(parts ...string) string {
	h := sha256.New()
	for i, p := range parts {
		if i > 0 {
			h.Write([]byte{0})
		}
		h.Write([]byte(p))
	}
	return hex.EncodeToString(h.Sum(nil))
}

func langFrom(code string) string {
	switch {
	case strings.HasPrefix(code, "km"):
		return "km"
	case strings.HasPrefix(code, "zh"):
		return "zh"
	}
	return "en"
}
