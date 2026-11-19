package notify

// Web Push tanpa dependensi luar:
//   - enkripsi isi: RFC 8291 (aes128gcm, ECDH P-256 + HKDF-SHA256)
//   - identitas server: VAPID, RFC 8292 (JWT ES256)
// Kunci VAPID dibuat sekali lalu disimpan di platform_settings['webpush_vapid'] (kunci privat dienkripsi DATA_KEY
// kalau tersedia). Semua proses (server & worker) memakai kunci yang sama.

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Box: enkripsi kunci privat VAPID (secretbox.Box memenuhi ini). Boleh nil.
type Box interface {
	Seal(plain, aad []byte) ([]byte, error)
	Open(sealed, aad []byte) ([]byte, error)
}

// Pusher mengirim Web Push.
type Pusher struct {
	Pool    *pgxpool.Pool
	Box     Box
	Subject string       // "mailto:…" atau "https://…" (kontak untuk layanan push)
	Client  *http.Client // nil = klien bawaan dengan batas waktu 15 detik

	mu   sync.Mutex
	key  *ecdsa.PrivateKey
	pub  []byte // titik publik tak terkompresi (65 byte)
	jwts map[string]cachedJWT
}

type cachedJWT struct {
	token string
	exp   time.Time
}

// Subscription: satu browser.
type Subscription struct {
	Endpoint string `json:"endpoint"`
	P256dh   string `json:"p256dh"`
	Auth     string `json:"auth"`
}

// ErrGone: langganan sudah tidak berlaku (404/410) — hapus.
var ErrGone = errors.New("langganan push sudah tidak berlaku")

var b64 = base64.RawURLEncoding

func decodeB64(s string) ([]byte, error) {
	s = strings.TrimRight(strings.TrimSpace(s), "=")
	if b, err := b64.DecodeString(s); err == nil {
		return b, nil
	}
	return base64.RawStdEncoding.DecodeString(s)
}

type vapidRow struct {
	Public  string `json:"public"`
	Private string `json:"private"`          // base64url D (32 byte), terenkripsi kalau Sealed
	Sealed  bool   `json:"sealed,omitempty"` // dienkripsi Box
}

// PublicKey: kunci publik VAPID (base64url) untuk PushManager.subscribe di browser.
func (p *Pusher) PublicKey(ctx context.Context) (string, error) {
	if err := p.load(ctx); err != nil {
		return "", err
	}
	return b64.EncodeToString(p.pub), nil
}

func (p *Pusher) load(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.key != nil {
		return nil
	}
	read := func() (vapidRow, bool, error) {
		var raw []byte
		var row vapidRow
		err := p.Pool.QueryRow(ctx, `SELECT value FROM platform_settings WHERE key = 'webpush_vapid'`).Scan(&raw)
		if errors.Is(err, pgx.ErrNoRows) {
			return row, false, nil
		}
		if err != nil {
			return row, false, err
		}
		return row, true, json.Unmarshal(raw, &row)
	}
	row, ok, err := read()
	if err != nil {
		return err
	}
	if !ok {
		k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return err
		}
		d := k.D.FillBytes(make([]byte, 32))
		pubKey, _ := k.PublicKey.ECDH()
		row = vapidRow{Public: b64.EncodeToString(pubKey.Bytes()), Private: b64.EncodeToString(d)}
		if p.Box != nil {
			sealed, err := p.Box.Seal(d, []byte("webpush_vapid"))
			if err != nil {
				return err
			}
			row.Private, row.Sealed = b64.EncodeToString(sealed), true
		}
		val, _ := json.Marshal(row)
		// Dua proses bisa membuat bersamaan: yang pertama menang, lalu baca ulang.
		if _, err := p.Pool.Exec(ctx, `INSERT INTO platform_settings (key, value) VALUES ('webpush_vapid', $1) ON CONFLICT (key) DO NOTHING`, val); err != nil {
			return err
		}
		if row, _, err = read(); err != nil {
			return err
		}
	}
	d, err := decodeB64(row.Private)
	if err != nil {
		return fmt.Errorf("kunci VAPID rusak: %w", err)
	}
	if row.Sealed {
		if p.Box == nil {
			return errors.New("kunci VAPID terenkripsi tetapi DATA_KEY tidak tersedia")
		}
		if d, err = p.Box.Open(d, []byte("webpush_vapid")); err != nil {
			return fmt.Errorf("kunci VAPID tidak bisa dibuka: %w", err)
		}
	}
	priv, err := ecdh.P256().NewPrivateKey(d)
	if err != nil {
		return err
	}
	pub := priv.PublicKey().Bytes()
	k := &ecdsa.PrivateKey{D: new(big.Int).SetBytes(d)}
	k.PublicKey.Curve = elliptic.P256()
	k.PublicKey.X = new(big.Int).SetBytes(pub[1:33])
	k.PublicKey.Y = new(big.Int).SetBytes(pub[33:])
	p.key, p.pub = k, pub
	return nil
}

// vapid: header Authorization untuk origin endpoint (JWT disimpan 1 jam).
func (p *Pusher) vapid(endpoint string) (string, error) {
	u, err := url.Parse(endpoint)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return "", fmt.Errorf("endpoint push tidak valid")
	}
	aud := u.Scheme + "://" + u.Host
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.jwts == nil {
		p.jwts = map[string]cachedJWT{}
	}
	if c, ok := p.jwts[aud]; ok && time.Until(c.exp) > time.Hour {
		return "vapid t=" + c.token + ", k=" + b64.EncodeToString(p.pub), nil
	}
	sub := p.Subject
	if sub == "" {
		sub = "mailto:support@phsar.app"
	}
	exp := time.Now().Add(12 * time.Hour)
	head := b64.EncodeToString([]byte(`{"typ":"JWT","alg":"ES256"}`))
	claims, _ := json.Marshal(map[string]any{"aud": aud, "exp": exp.Unix(), "sub": sub})
	signing := head + "." + b64.EncodeToString(claims)
	sum := sha256.Sum256([]byte(signing))
	r, s, err := ecdsa.Sign(rand.Reader, p.key, sum[:])
	if err != nil {
		return "", err
	}
	sig := append(r.FillBytes(make([]byte, 32)), s.FillBytes(make([]byte, 32))...)
	token := signing + "." + b64.EncodeToString(sig)
	p.jwts[aud] = cachedJWT{token: token, exp: exp}
	return "vapid t=" + token + ", k=" + b64.EncodeToString(p.pub), nil
}

func hmacSHA256(key, data []byte) []byte {
	m := hmac.New(sha256.New, key)
	m.Write(data)
	return m.Sum(nil)
}

// Encrypt: isi pesan untuk satu browser (RFC 8291, satu record aes128gcm).
// salt & asPriv boleh nil (dibuat acak) — diisi hanya oleh tes dengan vektor RFC.
func Encrypt(plain []byte, sub Subscription, salt []byte, asPriv *ecdh.PrivateKey) ([]byte, error) {
	uaPubBytes, err := decodeB64(sub.P256dh)
	if err != nil {
		return nil, fmt.Errorf("p256dh: %w", err)
	}
	authSecret, err := decodeB64(sub.Auth)
	if err != nil || len(authSecret) < 16 {
		return nil, errors.New("auth push tidak valid")
	}
	uaPub, err := ecdh.P256().NewPublicKey(uaPubBytes)
	if err != nil {
		return nil, fmt.Errorf("p256dh: %w", err)
	}
	if asPriv == nil {
		if asPriv, err = ecdh.P256().GenerateKey(rand.Reader); err != nil {
			return nil, err
		}
	}
	if salt == nil {
		salt = make([]byte, 16)
		if _, err := io.ReadFull(rand.Reader, salt); err != nil {
			return nil, err
		}
	}
	shared, err := asPriv.ECDH(uaPub)
	if err != nil {
		return nil, err
	}
	asPub := asPriv.PublicKey().Bytes()
	// IKM = HKDF(auth_secret, ecdh_secret, "WebPush: info" || 0x00 || ua_public || as_public, 32)
	prkKey := hmacSHA256(authSecret, shared)
	keyInfo := append(append(append([]byte("WebPush: info\x00"), uaPubBytes...), asPub...), 0x01)
	ikm := hmacSHA256(prkKey, keyInfo)
	prk := hmacSHA256(salt, ikm)
	cek := hmacSHA256(prk, []byte("Content-Encoding: aes128gcm\x00\x01"))[:16]
	nonce := hmacSHA256(prk, []byte("Content-Encoding: nonce\x00\x01"))[:12]
	block, err := aes.NewCipher(cek)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	record := append(append([]byte{}, plain...), 0x02) // 0x02 = record terakhir, tanpa padding
	body := gcm.Seal(nil, nonce, record, nil)
	var buf bytes.Buffer
	buf.Write(salt)
	_ = binary.Write(&buf, binary.BigEndian, uint32(4096))
	buf.WriteByte(byte(len(asPub)))
	buf.Write(asPub)
	buf.Write(body)
	return buf.Bytes(), nil
}

// Send mengirim satu pesan. urgency: "high" untuk pesanan baru/antrean, selain itu "normal".
func (p *Pusher) Send(ctx context.Context, sub Subscription, payload []byte, urgency string, ttl time.Duration) error {
	if err := p.load(ctx); err != nil {
		return err
	}
	body, err := Encrypt(payload, sub, nil, nil)
	if err != nil {
		return err
	}
	authz, err := p.vapid(sub.Endpoint)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, sub.Endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Content-Encoding", "aes128gcm")
	req.Header.Set("TTL", fmt.Sprint(int(ttl.Seconds())))
	req.Header.Set("Urgency", urgency)
	req.Header.Set("Authorization", authz)
	client := p.Client
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	msg, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
	switch {
	case resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusGone:
		return ErrGone
	case resp.StatusCode >= 400 && resp.StatusCode < 500 && resp.StatusCode != http.StatusTooManyRequests:
		return fmt.Errorf("%w: push %d %s", ErrPermanent, resp.StatusCode, strings.TrimSpace(string(msg)))
	case resp.StatusCode >= 300:
		return fmt.Errorf("push %d %s", resp.StatusCode, strings.TrimSpace(string(msg)))
	}
	return nil
}
