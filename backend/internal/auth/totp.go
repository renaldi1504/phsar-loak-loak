package auth

// TOTP (RFC 6238): kode 6 angka yang berganti tiap 30 detik, sama dengan Google Authenticator,
// Microsoft Authenticator, 1Password, dll. Tanpa library luar.

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base32"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
	"time"
)

const (
	totpPeriod = 30
	totpDigits = 6
	totpSkew   = 1 // terima kode 1 langkah sebelum/sesudah (jam HP yang sedikit meleset)
)

var b32 = base32.StdEncoding.WithPadding(base32.NoPadding)

// newTOTPSecret: 20 byte acak (160 bit) dalam base32.
func newTOTPSecret() string {
	b := make([]byte, 20)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return b32.EncodeToString(b)
}

func totpStep(t time.Time) int64 { return t.Unix() / totpPeriod }

func totpCode(secret string, step int64) (string, error) {
	key, err := b32.DecodeString(strings.ToUpper(strings.ReplaceAll(secret, " ", "")))
	if err != nil {
		return "", err
	}
	var msg [8]byte
	binary.BigEndian.PutUint64(msg[:], uint64(step))
	h := hmac.New(sha1.New, key)
	h.Write(msg[:])
	sum := h.Sum(nil)
	off := sum[len(sum)-1] & 0x0f
	v := binary.BigEndian.Uint32(sum[off:off+4]) & 0x7fffffff
	return fmt.Sprintf("%06d", v%1000000), nil
}

// totpMatch mengembalikan langkah yang cocok (untuk mencegah kode dipakai ulang), atau -1.
func totpMatch(secret, code string, now time.Time) int64 {
	code = strings.TrimSpace(strings.ReplaceAll(code, " ", ""))
	if len(code) != totpDigits {
		return -1
	}
	cur := totpStep(now)
	for d := int64(-totpSkew); d <= totpSkew; d++ {
		want, err := totpCode(secret, cur+d)
		if err == nil && hmac.Equal([]byte(want), []byte(code)) {
			return cur + d
		}
	}
	return -1
}

// otpauthURL: isi QR code untuk aplikasi authenticator.
func otpauthURL(secret, account string) string {
	label := url.PathEscape("Phsar Admin:" + account)
	q := url.Values{"secret": {secret}, "issuer": {"Phsar Admin"}, "algorithm": {"SHA1"}, "digits": {"6"}, "period": {"30"}}
	return "otpauth://totp/" + label + "?" + q.Encode()
}

// Kode cadangan: 10 kode sekali pakai, format "abcd-efgh" (tanpa huruf yang mirip: 0/o, 1/l/i).
const backupAlphabet = "abcdefghjkmnpqrstuvwxyz23456789"

func newBackupCodes(n int) []string {
	out := make([]string, n)
	for i := range out {
		b := make([]byte, 8)
		if _, err := rand.Read(b); err != nil {
			panic(err)
		}
		c := make([]byte, 8)
		for j, x := range b {
			c[j] = backupAlphabet[int(x)%len(backupAlphabet)]
		}
		out[i] = string(c[:4]) + "-" + string(c[4:])
	}
	return out
}

func normBackup(code string) string {
	return strings.ToLower(strings.NewReplacer("-", "", " ", "").Replace(strings.TrimSpace(code)))
}

func backupHash(code string) string {
	s := sha256.Sum256([]byte("phsar-backup:" + normBackup(code)))
	return hex.EncodeToString(s[:])
}
