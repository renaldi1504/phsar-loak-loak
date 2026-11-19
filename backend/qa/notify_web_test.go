package qa

import (
	"bufio"
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
	"encoding/json"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"phsar/internal/httpx"
	"phsar/internal/notify"
)

// ---------------------------------------------------------------------------
// Notifikasi web: realtime (SSE) + bunyi, admin per izin, pengaturan, Web Push, service worker.
// ---------------------------------------------------------------------------

// sseClient membuka aliran notifikasi dan mengumpulkan event.
type sseClient struct {
	mu     sync.Mutex
	events []map[string]any
	names  []string
	cancel context.CancelFunc
}

func openStream(t *testing.T, e *env, path, token string) *sseClient {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	req, _ := http.NewRequestWithContext(ctx, "GET", e.srv.URL+path, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Lang", "en")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		cancel()
		t.Fatalf("stream: %v", err)
	}
	check(t, res.StatusCode == 200 && strings.HasPrefix(res.Header.Get("Content-Type"), "text/event-stream"), "stream status %d %s", res.StatusCode, res.Header.Get("Content-Type"))
	c := &sseClient{cancel: cancel}
	go func() {
		defer res.Body.Close()
		sc := bufio.NewScanner(res.Body)
		name := ""
		for sc.Scan() {
			line := sc.Text()
			switch {
			case strings.HasPrefix(line, "event: "):
				name = strings.TrimPrefix(line, "event: ")
			case strings.HasPrefix(line, "data: "):
				var v map[string]any
				_ = json.Unmarshal([]byte(strings.TrimPrefix(line, "data: ")), &v)
				c.mu.Lock()
				c.events, c.names = append(c.events, v), append(c.names, name)
				c.mu.Unlock()
			}
		}
	}()
	return c
}

// wait: event bernama `name` yang memenuhi `ok`, maks. 5 detik.
func (c *sseClient) wait(t *testing.T, name string, ok func(map[string]any) bool) map[string]any {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		c.mu.Lock()
		for i, n := range c.names {
			if n == name && (ok == nil || ok(c.events[i])) {
				v := c.events[i]
				c.mu.Unlock()
				return v
			}
		}
		c.mu.Unlock()
		time.Sleep(30 * time.Millisecond)
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	t.Fatalf("event %q tidak datang; yang ada: %v %v", name, c.names, c.events)
	return nil
}

func kindIs(k string) func(map[string]any) bool {
	return func(v map[string]any) bool { return v["kind"] == k }
}

func TestWebNotifications(t *testing.T) {
	e := getEnv(t)
	tok, sid := newActiveStore(t, e)
	uid := e.scalar(t, `SELECT user_id::text FROM store_members WHERE store_id = $1`, sid)

	t.Run("service worker dilayani di seller, admin, website", func(t *testing.T) {
		for _, p := range []string{"/seller/sw.js", "/admin/sw.js", "/sw.js"} {
			r := e.call(t, "GET", p, nil, "").want(200)
			check(t, strings.Contains(string(r.Raw), "showNotification"), "%s bukan service worker", p)
		}
		e.call(t, "GET", "/ui/icon-192.png", nil, "").want(200)
	})

	t.Run("realtime: pesanan baru → bunyi order 3×", func(t *testing.T) {
		s := openStream(t, e, "/v1/seller/notifications/stream", tok)
		defer s.cancel()
		s.wait(t, "hello", nil)
		e.sql(t, `INSERT INTO notifications (user_id, store_id, kind, params, link, audience)
			VALUES ($1, $2, 'order.new', '{"order":"PH-QA-1"}', '#/store/x/orders', 'seller')`, uid, sid)
		ev := s.wait(t, "notification", kindIs("order.new"))
		check(t, ev["sound"] == "order" && ev["repeat"] == float64(3), "bunyi pesanan baru harus order ×3: %v", ev)
		check(t, strings.Contains(ev["title"].(string), "order") || ev["title"] != "", "judul kosong: %v", ev)
		check(t, ev["unread"].(float64) >= 1, "unread: %v", ev)
		e.sql(t, `INSERT INTO notifications (user_id, store_id, kind, params, audience) VALUES ($1, $2, 'chat.new', '{"person":"Sokha","note":"Hi"}', 'seller')`, uid, sid)
		check(t, s.wait(t, "notification", kindIs("chat.new"))["sound"] == "chat", "chat harus bunyi chat")
		e.call(t, "POST", "/v1/seller/notifications/read", map[string]any{"all": true}, tok).want(204)
		s.wait(t, "read", func(v map[string]any) bool { return v["unread"] == float64(0) })
	})

	t.Run("pengaturan: suara grup mati, jam tenang, validasi", func(t *testing.T) {
		p := e.call(t, "GET", "/v1/seller/notifications/prefs", nil, tok).want(200)
		check(t, p.len("groups") == 6 && p.str("prefs.groups.orders.sound") == "true" && p.str("vapid_public_key") != "", "prefs awal: %s", p.Raw)
		e.call(t, "PUT", "/v1/seller/notifications/prefs", map[string]any{"quiet_from": "23:00"}, tok).wantErr(422, "validation")
		e.call(t, "PUT", "/v1/seller/notifications/prefs", map[string]any{"groups": map[string]any{"orders": map[string]bool{"sound": false, "push": true, "telegram": false}}, "volume": 55}, tok).
			want(200).eq("prefs.groups.orders.sound", "false").eq("prefs.groups.chat.sound", "true").eq("prefs.volume", "55")
		s := openStream(t, e, "/v1/seller/notifications/stream", tok)
		defer s.cancel()
		s.wait(t, "hello", nil)
		e.sql(t, `INSERT INTO notifications (user_id, store_id, kind, params, audience) VALUES ($1, $2, 'order.new', '{}', 'seller')`, uid, sid)
		ev := s.wait(t, "notification", kindIs("order.new"))
		check(t, ev["sound"] == "none" && ev["volume"] == float64(55), "suara pesanan dimatikan → none: %v", ev)
		// jam tenang sepanjang hari (00:00–23:59) → semua tanpa bunyi; notifikasi uji tetap berbunyi
		e.call(t, "PUT", "/v1/seller/notifications/prefs", map[string]any{"quiet_from": "00:00", "quiet_to": "23:59", "volume": 80}, tok).want(200)
		e.sql(t, `INSERT INTO notifications (user_id, store_id, kind, params, audience) VALUES ($1, $2, 'chat.new', '{}', 'seller')`, uid, sid)
		check(t, s.wait(t, "notification", kindIs("chat.new"))["sound"] == "none", "jam tenang → tanpa bunyi")
		e.call(t, "POST", "/v1/seller/notifications/test", nil, tok).want(204)
		check(t, s.wait(t, "notification", kindIs("system.test"))["sound"] == "ding", "notifikasi uji tetap berbunyi")
		e.call(t, "POST", "/v1/seller/notifications/test", nil, tok).wantErr(429, "rate_limited")
		e.call(t, "PUT", "/v1/seller/notifications/prefs", map[string]any{"volume": 80}, tok).want(200)
	})

	t.Run("admin: dibagi per izin peran", func(t *testing.T) {
		mod := newAdmin(t, e, "moderator")
		fin := newAdmin(t, e, "finance")
		cs := newAdmin(t, e, "cs")
		before := func(tok string) int {
			return e.call(t, "GET", "/admin/api/notifications", nil, tok).want(200).len("notifications")
		}
		nFin, nCS := before(fin), before(cs)
		s := openStream(t, e, "/admin/api/notifications/stream", mod)
		defer s.cancel()
		s.wait(t, "hello", nil)
		_, sid2, _ := submittedStore(t, e, mod) // toko dikirim untuk diverifikasi → izin seller.verify
		ev := s.wait(t, "notification", kindIs("admin.store_review"))
		check(t, ev["sound"] == "alert" && ev["link"] == "#/store/"+sid2 && strings.Contains(ev["body"].(string), "Notif Shop"), "notifikasi admin: %v", ev)
		check(t, before(fin) == nFin && before(cs) == nCS, "finance & cs tidak boleh menerima antrean toko")
		// payout menunggu → finance (payout.approve), bukan moderator
		e.sql(t, `INSERT INTO payouts (store_id, amount, status, bank, account_name, account_last4, order_count) VALUES ($1, 12.50, 'pending', 'ABA', 'QA', '1234', 1)`, sid)
		r := e.call(t, "GET", "/admin/api/notifications", nil, fin).want(200)
		r.eq("notifications.0.kind", "admin.payout_pending").eq("notifications.0.link", "#/payouts")
		check(t, strings.Contains(r.str("notifications.0.body"), "$12.50"), "isi payout: %s", r.Raw)
		check(t, !strings.Contains(string(e.call(t, "GET", "/admin/api/notifications", nil, mod).want(200).Raw), "admin.payout_pending"), "moderator tidak boleh menerima payout")
		// admin tidak dikirimi Telegram secara default
		e.call(t, "GET", "/admin/api/notifications/prefs", nil, fin).want(200).eq("prefs.groups.payouts.telegram", "false").eq("audience", "admin")
		// lonceng seller tidak tercampur notifikasi admin
		check(t, !strings.Contains(string(e.call(t, "GET", "/v1/seller/notifications", nil, tok).want(200).Raw), "admin."), "lonceng seller tercampur")
	})

	t.Run("pembeli: lonceng & aliran sendiri", func(t *testing.T) {
		buyer, phone := e.login(t)
		me := e.scalar(t, `SELECT id::text FROM users WHERE phone = '+855' || substr($1, 2)`, phone)
		e.sql(t, `INSERT INTO notifications (user_id, kind, params, link, audience) VALUES ($1, 'buyer.order_on_the_way', '{"order":"PH-QA-9","store":"Toko"}', '/en/orders/PH-QA-9', 'buyer')`, me)
		r := e.call(t, "GET", "/v1/shop/notifications", nil, buyer).want(200)
		r.eq("notifications.0.kind", "buyer.order_on_the_way").eq("unread", "1")
		e.call(t, "GET", "/v1/shop/notifications/prefs", nil, buyer).want(200).eq("audience", "buyer")
	})

	t.Run("Web Push: terenkripsi RFC 8291 + VAPID, 410 menghapus langganan", func(t *testing.T) {
		var mu sync.Mutex
		var got [][]byte
		var authz []string
		status := 201
		push := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			b, _ := io.ReadAll(r.Body)
			mu.Lock()
			got, authz = append(got, b), append(authz, r.Header.Get("Authorization"))
			st := status
			mu.Unlock()
			check(t, r.Header.Get("Content-Encoding") == "aes128gcm" && r.Header.Get("TTL") != "", "header push: %v", r.Header)
			w.WriteHeader(st)
		}))
		defer push.Close()
		ua, _ := ecdh.P256().GenerateKey(rand.Reader)
		secret := make([]byte, 16)
		_, _ = rand.Read(secret)
		enc := base64.RawURLEncoding
		e.call(t, "POST", "/v1/seller/notifications/push", map[string]any{"endpoint": "http://not-https.example/x", "keys": map[string]string{"p256dh": "x", "auth": "y"}}, tok).wantErr(400, "validation")
		e.call(t, "POST", "/v1/seller/notifications/push", map[string]any{"endpoint": push.URL + "/sub/1",
			"keys": map[string]string{"p256dh": enc.EncodeToString(ua.PublicKey().Bytes()), "auth": enc.EncodeToString(secret)}}, tok).want(204)
		e.call(t, "GET", "/v1/seller/notifications/prefs", nil, tok).want(200).eq("push_devices", "1")

		e.sql(t, `UPDATE notifications SET push_status = 'skipped' WHERE push_status = 'pending'`) // antrean tes lain
		time.Sleep(11 * time.Second)                                                               // batas 1 notifikasi uji / 10 detik
		e.call(t, "POST", "/v1/seller/notifications/test", nil, tok).want(204)
		pusher := httpx.NewPusher(e.cfg, e.pool)
		pusher.Client = push.Client()
		w := &notify.Worker{Pool: e.pool, Push: pusher, BaseURL: "https://phsar.test", Log: e.worker.Log}
		for i := 0; i < 10; i++ {
			if n, err := w.RunOnce(context.Background()); err != nil || n == 0 {
				break
			}
		}
		mu.Lock()
		check(t, len(got) == 1, "push terkirim %d kali", len(got))
		body, header := got[0], authz[0]
		mu.Unlock()
		plain := decryptPush(t, body, ua, secret)
		var p map[string]any
		check(t, json.Unmarshal(plain, &p) == nil, "isi push bukan JSON: %q", plain)
		check(t, strings.Contains(p["title"].(string), "Test notification") && p["url"] == "https://phsar.test/seller/" && p["test"] == true, "isi push: %v", p)
		key := e.call(t, "GET", "/v1/seller/notifications/prefs", nil, tok).want(200).str("vapid_public_key")
		checkVAPID(t, header, key, push.URL)
		check(t, e.scalar(t, `SELECT push_status FROM notifications WHERE user_id = $1 AND kind = 'system.test' ORDER BY id DESC LIMIT 1`, uid) == "sent", "status push bukan sent")

		// grup dimatikan → push dilewati
		e.call(t, "PUT", "/v1/seller/notifications/prefs", map[string]any{"groups": map[string]any{"chat": map[string]bool{"sound": true, "push": false, "telegram": true}}, "volume": 80}, tok).want(200)
		e.sql(t, `INSERT INTO notifications (user_id, store_id, kind, params, audience) VALUES ($1, $2, 'chat.new', '{}', 'seller')`, uid, sid)
		_, _ = w.RunOnce(context.Background())
		check(t, e.scalar(t, `SELECT push_status FROM notifications WHERE user_id = $1 AND kind = 'chat.new' ORDER BY id DESC LIMIT 1`, uid) == "skipped", "push chat harusnya dilewati")

		// langganan kedaluwarsa (410) → dihapus
		mu.Lock()
		status = 410
		mu.Unlock()
		e.sql(t, `INSERT INTO notifications (user_id, store_id, kind, params, audience) VALUES ($1, $2, 'order.new', '{}', 'seller')`, uid, sid)
		_, _ = w.RunOnce(context.Background())
		e.call(t, "GET", "/v1/seller/notifications/prefs", nil, tok).want(200).eq("push_devices", "0")
	})

	t.Run("Telegram mengikuti pengaturan", func(t *testing.T) {
		e.sql(t, `INSERT INTO auth_identities (user_id, provider, provider_uid, display) VALUES ($1, 'telegram', '555000777', '@qa') ON CONFLICT DO NOTHING`, uid)
		e.call(t, "PUT", "/v1/seller/notifications/prefs", map[string]any{"groups": map[string]any{"orders": map[string]bool{"sound": true, "push": true, "telegram": false}}, "volume": 80}, tok).want(200)
		e.sql(t, `UPDATE notifications SET telegram_status = 'skipped' WHERE telegram_status = 'pending'`)
		e.sql(t, `INSERT INTO notifications (user_id, store_id, kind, params, audience) VALUES ($1, $2, 'order.new', '{}', 'seller')`, uid, sid)
		e.sql(t, `INSERT INTO notifications (user_id, store_id, kind, params, audience) VALUES ($1, $2, 'payout.sent', '{}', 'seller')`, uid, sid)
		bot := &fakeTelegram{}
		w := &notify.Worker{Pool: e.pool, Bot: bot, BaseURL: "https://phsar.test", Log: e.worker.Log}
		_, _ = w.RunOnce(context.Background())
		check(t, len(bot.sent) == 1, "hanya payout yang dikirim ke Telegram: %v", bot.sent)
		check(t, e.scalar(t, `SELECT telegram_error FROM notifications WHERE user_id = $1 AND kind = 'order.new' ORDER BY id DESC LIMIT 1`, uid) == "dimatikan di pengaturan notifikasi", "alasan skip")
	})
}

// decryptPush: sisi browser RFC 8291 (untuk memastikan server mengenkripsi dengan benar).
func decryptPush(t *testing.T, body []byte, ua *ecdh.PrivateKey, secret []byte) []byte {
	t.Helper()
	check(t, len(body) > 86, "body push terlalu pendek")
	salt, idlen := body[:16], int(body[20])
	asPub := body[21 : 21+idlen]
	ct := body[21+idlen:]
	as, err := ecdh.P256().NewPublicKey(asPub)
	check(t, err == nil, "kunci pengirim: %v", err)
	shared, _ := ua.ECDH(as)
	mac := func(k, d []byte) []byte { m := hmac.New(sha256.New, k); m.Write(d); return m.Sum(nil) }
	ikm := mac(mac(secret, shared), append(append(append([]byte("WebPush: info\x00"), ua.PublicKey().Bytes()...), asPub...), 1))
	prk := mac(salt, ikm)
	cek := mac(prk, []byte("Content-Encoding: aes128gcm\x00\x01"))[:16]
	nonce := mac(prk, []byte("Content-Encoding: nonce\x00\x01"))[:12]
	block, _ := aes.NewCipher(cek)
	gcm, _ := cipher.NewGCM(block)
	plain, err := gcm.Open(nil, nonce, ct, nil)
	check(t, err == nil, "dekripsi push gagal: %v", err)
	check(t, len(plain) > 0 && plain[len(plain)-1] == 2, "penanda record terakhir hilang")
	return plain[:len(plain)-1]
}

// checkVAPID: header "vapid t=<JWT>, k=<kunci>" dengan tanda tangan ES256 yang sah untuk origin endpoint.
func checkVAPID(t *testing.T, header, key, origin string) {
	t.Helper()
	check(t, strings.HasPrefix(header, "vapid t=") && strings.HasSuffix(header, ", k="+key), "header VAPID: %s", header)
	jwt := strings.TrimSuffix(strings.TrimPrefix(header, "vapid t="), ", k="+key)
	parts := strings.Split(jwt, ".")
	check(t, len(parts) == 3, "JWT: %s", jwt)
	enc := base64.RawURLEncoding
	claims, _ := enc.DecodeString(parts[1])
	var c map[string]any
	_ = json.Unmarshal(claims, &c)
	check(t, c["aud"] == origin && c["sub"] != "" && c["exp"].(float64) > float64(time.Now().Unix()), "klaim VAPID: %v", c)
	pub, _ := enc.DecodeString(key)
	sig, _ := enc.DecodeString(parts[2])
	k := ecdsa.PublicKey{Curve: elliptic.P256(), X: new(big.Int).SetBytes(pub[1:33]), Y: new(big.Int).SetBytes(pub[33:])}
	sum := sha256.Sum256([]byte(parts[0] + "." + parts[1]))
	check(t, len(sig) == 64 && ecdsa.Verify(&k, sum[:], new(big.Int).SetBytes(sig[:32]), new(big.Int).SetBytes(sig[32:])), "tanda tangan VAPID tidak sah")
}
