package qa

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"strings"
	"testing"
	"time"
)

var csrf = []string{"X-Phsar", "1"}

// totpAt: kode authenticator untuk secret, `off` langkah 30 detik dari sekarang (−1, 0, +1 diterima server).
func totpAt(secret string, off int64) string {
	key, _ := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secret)
	var msg [8]byte
	binary.BigEndian.PutUint64(msg[:], uint64(time.Now().Unix()/30+off))
	h := hmac.New(sha1.New, key)
	h.Write(msg[:])
	sum := h.Sum(nil)
	o := sum[len(sum)-1] & 0x0f
	return fmt.Sprintf("%06d", (binary.BigEndian.Uint32(sum[o:o+4])&0x7fffffff)%1000000)
}

type adminAcct struct {
	phone, secret, token string
	backup               []string
}

// adminStep1: login HP di halaman admin (cookie route) → access token admin (belum lolos MFA).
func (e *env) adminStep1(t *testing.T, phone string) resp {
	t.Helper()
	// tes login berkali-kali dengan nomor yang sama: lewati jeda 60 detik & batas per jam OTP
	e.sql(t, `UPDATE otp_codes SET created_at = created_at - interval '2 hours' WHERE phone = '+855' || substr($1, 2)`, phone)
	r := e.call(t, "POST", "/admin/auth/otp/request", map[string]string{"phone": phone}, "", csrf...)
	r.want(200)
	return e.call(t, "POST", "/admin/auth/otp/verify", map[string]string{"phone": phone, "code": e.sms.code(r.str("phone"))}, "", csrf...)
}

func newAdminFull(t *testing.T, e *env, role string) (string, adminAcct) {
	t.Helper()
	_, phone := e.login(t)
	e.sql(t, `INSERT INTO admin_users (user_id, role_code) SELECT id, $2 FROM users WHERE phone = '+855' || substr($1, 2)`, phone, role)
	s1 := e.adminStep1(t, phone).want(200)
	setup := e.call(t, "POST", "/admin/auth/mfa/setup", nil, s1.str("access_token")).want(200)
	secret := setup.str("secret")
	en := e.call(t, "POST", "/admin/auth/mfa/enable", map[string]string{"code": totpAt(secret, 0)}, s1.str("access_token")).want(200)
	a := adminAcct{phone: phone, secret: secret, token: en.str("session.access_token")}
	for i := 0; i < en.len("backup_codes"); i++ {
		a.backup = append(a.backup, en.str("backup_codes."+itoa(i)))
	}
	return a.token, a
}

// ---------------------------------------------------------------------------
// 15. Login admin + MFA
// ---------------------------------------------------------------------------

func TestAdminMFA(t *testing.T) {
	e := getEnv(t)

	t.Run("bukan admin tidak mendapat sesi admin; token app ditolak API admin", func(t *testing.T) {
		app, phone := e.login(t)
		e.adminStep1(t, phone).wantErr(403, "not_admin")
		e.call(t, "GET", "/admin/api/me", nil, app).want(401)
	})

	_, phone := e.login(t)
	e.sql(t, `INSERT INTO admin_users (user_id, role_code) SELECT id, 'super_admin' FROM users WHERE phone = '+855' || substr($1, 2)`, phone)
	var secret, tok string
	var backups []string
	t.Run("pasang MFA: QR/secret → kode pertama → sesi lolos + 10 kode cadangan", func(t *testing.T) {
		s1 := e.adminStep1(t, phone).want(200)
		s1.eq("mfa.enrolled", "false").eq("mfa.verified", "false")
		pre := s1.str("access_token")
		e.call(t, "GET", "/admin/api/me", nil, pre).wantErr(403, "mfa_required")
		e.call(t, "POST", "/admin/auth/mfa/verify", map[string]string{"code": "123456"}, pre).wantErr(409, "mfa_not_enrolled")
		setup := e.call(t, "POST", "/admin/auth/mfa/setup", nil, pre).want(200)
		secret = setup.str("secret")
		check(t, len(secret) == 32 && strings.HasPrefix(setup.str("otpauth_url"), "otpauth://totp/Phsar%20Admin:"), "secret/QR tidak valid: %s", setup.Raw)
		check(t, e.scalar(t, `SELECT position($1 in encode(mfa_pending_enc, 'escape'))::text FROM admin_users a JOIN users u ON u.id = a.user_id WHERE u.phone = '+855' || substr($2, 2)`, secret, phone) == "0",
			"secret MFA tersimpan tanpa enkripsi")
		e.call(t, "POST", "/admin/auth/mfa/enable", map[string]string{"code": "000000"}, pre).wantErr(401, "mfa_invalid")
		en := e.call(t, "POST", "/admin/auth/mfa/enable", map[string]string{"code": totpAt(secret, 0)}, pre).want(200)
		en.eq("session.mfa.verified", "true").eq("session.mfa.enrolled", "true")
		check(t, en.len("backup_codes") == 10, "kode cadangan harus 10")
		for i := 0; i < 10; i++ {
			backups = append(backups, en.str("backup_codes."+itoa(i)))
		}
		tok = en.str("session.access_token")
		e.call(t, "GET", "/admin/api/me", nil, tok).want(200).eq("role", "super_admin").eq("backup_codes_left", "10")
		e.call(t, "POST", "/admin/auth/mfa/setup", nil, tok).wantErr(409, "mfa_enrolled")
	})
	t.Run("login berikutnya: kode authenticator; kode yang sama tidak bisa dipakai ulang", func(t *testing.T) {
		pre := e.adminStep1(t, phone).want(200)
		pre.eq("mfa.enrolled", "true").eq("mfa.verified", "false")
		e.call(t, "POST", "/admin/auth/mfa/verify", map[string]string{"code": totpAt(secret, 0)}, pre.str("access_token")).wantErr(401, "mfa_invalid") // dipakai saat pasang
		v := e.call(t, "POST", "/admin/auth/mfa/verify", map[string]string{"code": totpAt(secret, 1)}, pre.str("access_token")).want(200)
		e.call(t, "GET", "/admin/api/me", nil, v.str("session.access_token")).want(200)
		check(t, e.scalar(t, `SELECT count(*) FROM audit_logs WHERE action = 'admin.login'`) >= "1", "login admin tidak dicatat")
	})
	t.Run("kode cadangan: sekali pakai", func(t *testing.T) {
		pre := e.adminStep1(t, phone).want(200).str("access_token")
		v := e.call(t, "POST", "/admin/auth/mfa/verify", map[string]string{"code": strings.ToUpper(backups[0])}, pre).want(200)
		v.eq("backup_codes_left", "9")
		pre2 := e.adminStep1(t, phone).want(200).str("access_token")
		e.call(t, "POST", "/admin/auth/mfa/verify", map[string]string{"code": backups[0]}, pre2).wantErr(401, "mfa_invalid")
	})
	t.Run("5 kode salah → sesi dicabut", func(t *testing.T) {
		pre := e.adminStep1(t, phone).want(200).str("access_token")
		for i := 0; i < 5; i++ {
			e.call(t, "POST", "/admin/auth/mfa/verify", map[string]string{"code": "111111"}, pre).wantErr(401, "mfa_invalid")
		}
		e.call(t, "POST", "/admin/auth/mfa/verify", map[string]string{"code": backups[1]}, pre).want(401) // sesi sudah dicabut
	})
	t.Run("kelola admin: tambah, ganti peran, reset MFA, nonaktifkan", func(t *testing.T) {
		mod := newAdmin(t, e, "moderator")
		e.call(t, "GET", "/admin/api/admins", nil, mod).want(403)
		_, staffPhone := e.login(t)
		e.call(t, "POST", "/admin/api/admins", map[string]string{"phone": "099 000 111", "role": "cs"}, tok).wantErr(404, "user_not_found")
		e.call(t, "POST", "/admin/api/admins", map[string]string{"phone": staffPhone, "role": "boss"}, tok).wantErr(422, "validation", "role")
		e.call(t, "POST", "/admin/api/admins", map[string]string{"phone": staffPhone, "role": "cs"}, tok).want(200)
		l := e.call(t, "GET", "/admin/api/admins", nil, tok).want(200)
		check(t, strings.Contains(string(l.Raw), `"role":"cs"`) && strings.Contains(string(l.Raw), `"self":true`), "daftar admin salah: %s", l.Raw)
		staffID := e.scalar(t, `SELECT id::text FROM users WHERE phone = '+855' || substr($1, 2)`, staffPhone)
		myID := e.scalar(t, `SELECT id::text FROM users WHERE phone = '+855' || substr($1, 2)`, phone)
		e.call(t, "POST", "/admin/api/admins/"+myID+"/mfa-reset", nil, tok).wantErr(409, "self")
		e.call(t, "PATCH", "/admin/api/admins/"+myID, map[string]any{"active": false}, tok).wantErr(409, "self")

		// staf memasang MFA, lalu HP hilang → super_admin reset → sesi staf dicabut, MFA dipasang ulang
		s1 := e.adminStep1(t, staffPhone).want(200).str("access_token")
		sec := e.call(t, "POST", "/admin/auth/mfa/setup", nil, s1).want(200).str("secret")
		staff := e.call(t, "POST", "/admin/auth/mfa/enable", map[string]string{"code": totpAt(sec, 0)}, s1).want(200).str("session.access_token")
		e.call(t, "GET", "/admin/api/me", nil, staff).want(200).eq("role", "cs")
		e.call(t, "POST", "/admin/api/admins/"+staffID+"/mfa-reset", nil, tok).want(200)
		e.call(t, "GET", "/admin/api/me", nil, staff).want(401)
		e.adminStep1(t, staffPhone).want(200).eq("mfa.enrolled", "false")

		e.call(t, "PATCH", "/admin/api/admins/"+staffID, map[string]any{"active": false}, tok).want(200)
		e.adminStep1(t, staffPhone).wantErr(403, "not_admin")
		e.call(t, "PATCH", "/admin/api/admins/not-a-uuid", map[string]any{"active": true}, tok).want(404)
	})
	t.Run("sesi admin ditutup kalau lama tidak dipakai", func(t *testing.T) {
		v := e.adminStep1(t, phone).want(200)
		cookieOf := func(r resp) string {
			c := r.Header.Get("Set-Cookie")
			check(t, strings.Contains(c, "phsar_admin=") && strings.Contains(c, "HttpOnly"), "cookie sesi admin tidak aman: %s", c)
			return strings.SplitN(strings.TrimPrefix(c, "phsar_admin="), ";", 2)[0]
		}
		val := cookieOf(v)
		r := e.call(t, "POST", "/admin/auth/refresh", nil, "", "X-Phsar", "1", "Cookie", "phsar_admin="+val).want(200)
		val = cookieOf(r)
		e.sql(t, `UPDATE sessions SET last_seen_at = now() - interval '31 minutes' WHERE audience = 'admin' AND revoked_at IS NULL
			AND user_id = (SELECT id FROM users WHERE phone = '+855' || substr($1, 2))`, phone)
		e.call(t, "POST", "/admin/auth/refresh", nil, "", "X-Phsar", "1", "Cookie", "phsar_admin="+val).wantErr(401, "session_invalid")
	})
}
