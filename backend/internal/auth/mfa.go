package auth

// Login admin dua langkah:
//  1. nomor HP (SMS/Telegram) — hanya akun yang terdaftar di admin_users yang mendapat sesi admin;
//  2. kode 6 angka dari aplikasi authenticator (TOTP) atau satu kode cadangan.
// API admin (/admin/api) hanya menerima access token dengan klaim mfa=true.

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

	"phsar/internal/platform/db"
)

var (
	ErrNotAdmin           = errors.New("akun ini bukan admin Phsar")
	ErrMFARequired        = errors.New("masukkan kode authenticator dulu")
	ErrMFAInvalid         = errors.New("kode salah atau sudah dipakai")
	ErrMFALocked          = errors.New("terlalu banyak kode salah, coba lagi 15 menit lagi")
	ErrMFANotEnrolled     = errors.New("MFA belum dipasang")
	ErrMFAAlreadyEnrolled = errors.New("MFA sudah dipasang")
	ErrMFAUnavailable     = errors.New("DATA_ENCRYPTION_KEY belum diisi")
)

const (
	// AdminIdleTimeout: sesi admin ditutup kalau tidak dipakai selama ini (access token 15 menit + jeda).
	AdminIdleTimeout = 30 * time.Minute
	mfaMaxFailures   = 5  // per sesi: lebih dari ini sesi dicabut, harus login ulang
	mfaUserFailures  = 10 // per admin per 15 menit (semua sesi)
	backupCodeCount  = 10
)

// MFAState dikirim bersama sesi admin supaya halaman tahu langkah berikutnya.
type MFAState struct {
	Enrolled bool `json:"enrolled"` // aplikasi authenticator sudah dipasang
	Verified bool `json:"verified"` // sesi ini sudah memasukkan kode
}

// adminAllowed: untuk audience admin, hanya admin aktif yang boleh mendapat sesi.
func adminAllowed(ctx context.Context, q queryer, userID, audience string) (bool, error) {
	if audience != "admin" {
		return true, nil
	}
	var ok bool
	err := q.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM admin_users WHERE user_id = $1 AND active AND deleted_at IS NULL)`, userID).Scan(&ok)
	return ok, err
}

func mfaState(ctx context.Context, q queryer, userID, sid string) (MFAState, error) {
	var st MFAState
	err := q.QueryRow(ctx, `SELECT coalesce((SELECT mfa_enabled FROM admin_users WHERE user_id = $1), false),
		coalesce((SELECT mfa_verified_at IS NOT NULL FROM sessions WHERE id = $2), false)`, userID, sid).Scan(&st.Enrolled, &st.Verified)
	return st, err
}

// RequireMFA: pasang setelah RequireUser("admin").
func (s *Service) RequireMFA(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, ok := UserFrom(r.Context())
		if !ok || c.Audience != "admin" || !c.MFA {
			writeErr(w, http.StatusForbidden, "mfa_required", ErrMFARequired.Error())
			return
		}
		next.ServeHTTP(w, r)
	})
}

// AccessToken: token akses baru untuk sesi yang sama (dipakai setelah MFA lolos; refresh token tidak berubah).
type AccessToken struct {
	AccessToken     string    `json:"access_token"`
	AccessExpiresAt time.Time `json:"access_expires_at"`
	User            UserInfo  `json:"user"`
	MFA             *MFAState `json:"mfa"`
}

func (s *Service) accessFor(ctx context.Context, q queryer, c Claims) (AccessToken, error) {
	st, err := mfaState(ctx, q, c.Subject, c.SessionID)
	if err != nil {
		return AccessToken{}, err
	}
	now := s.now()
	tok, err := signJWT(s.secret, Claims{Subject: c.Subject, Audience: c.Audience, SessionID: c.SessionID, IssuedAt: now.Unix(),
		ExpiresAt: now.Add(AccessTTL).Unix(), MFA: st.Verified})
	if err != nil {
		return AccessToken{}, err
	}
	u, err := loadUser(ctx, q, c.Subject)
	if err != nil {
		return AccessToken{}, err
	}
	return AccessToken{AccessToken: tok, AccessExpiresAt: now.Add(AccessTTL), User: u, MFA: &st}, nil
}

func (s *Service) needBox() error {
	if s.box == nil {
		return ErrMFAUnavailable
	}
	return nil
}

func secretAAD(userID string) []byte { return []byte("admin-mfa:" + userID) }

// MFASetup membuat secret baru (belum aktif sampai dikonfirmasi dengan kode dari aplikasi).
func (s *Service) MFASetup(ctx context.Context, c Claims) (secret, url string, err error) {
	if err := s.needBox(); err != nil {
		return "", "", err
	}
	err = db.WithActor(ctx, s.pool, c.Subject, "admin", "", func(tx pgx.Tx) error {
		var enabled bool
		var phone *string
		if err := tx.QueryRow(ctx, `SELECT a.mfa_enabled, u.phone FROM admin_users a JOIN users u ON u.id = a.user_id
			WHERE a.user_id = $1 AND a.active AND a.deleted_at IS NULL FOR UPDATE OF a`, c.Subject).Scan(&enabled, &phone); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrNotAdmin
			}
			return err
		}
		if enabled {
			return ErrMFAAlreadyEnrolled
		}
		secret = newTOTPSecret()
		sealed, err := s.box.Seal([]byte(secret), secretAAD(c.Subject))
		if err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `UPDATE admin_users SET mfa_pending_enc = $2 WHERE user_id = $1`, c.Subject, sealed); err != nil {
			return err
		}
		account := deref(phone)
		if account == "" {
			account = c.Subject[:8]
		}
		url = otpauthURL(secret, account)
		return nil
	})
	return secret, url, err
}

// MFAEnable: kode pertama dari aplikasi → MFA aktif, sesi ini lolos, kode cadangan dibuat (ditampilkan sekali).
func (s *Service) MFAEnable(ctx context.Context, c Claims, code string, cl ClientInfo) (AccessToken, []string, error) {
	if err := s.needBox(); err != nil {
		return AccessToken{}, nil, err
	}
	if err := s.checkLock(ctx, c); err != nil {
		return AccessToken{}, nil, err
	}
	var out AccessToken
	var codes []string
	wrong := false
	err := db.WithActor(ctx, s.pool, c.Subject, "admin", "", func(tx pgx.Tx) error {
		var enabled bool
		var pending []byte
		if err := tx.QueryRow(ctx, `SELECT mfa_enabled, mfa_pending_enc FROM admin_users
			WHERE user_id = $1 AND active AND deleted_at IS NULL FOR UPDATE`, c.Subject).Scan(&enabled, &pending); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrNotAdmin
			}
			return err
		}
		if enabled {
			return ErrMFAAlreadyEnrolled
		}
		if pending == nil {
			return ErrMFANotEnrolled
		}
		secret, err := s.box.Open(pending, secretAAD(c.Subject))
		if err != nil {
			return err
		}
		step := totpMatch(string(secret), code, s.now())
		if step < 0 {
			wrong = true
			return s.failMFA(ctx, tx, c)
		}
		if _, err := tx.Exec(ctx, `UPDATE admin_users SET mfa_secret_enc = mfa_pending_enc, mfa_pending_enc = NULL, mfa_enabled = true,
			mfa_enabled_at = now(), mfa_last_step = $2, last_login_at = now() WHERE user_id = $1`, c.Subject, step); err != nil {
			return err
		}
		codes, err = replaceBackupCodes(ctx, tx, c.Subject)
		if err != nil {
			return err
		}
		if err := markVerified(ctx, tx, c, cl, "admin.mfa_enabled"); err != nil {
			return err
		}
		out, err = s.accessFor(ctx, tx, c)
		return err
	})
	s.recordAttempt(ctx, "mfa:"+c.Subject, cl, err == nil && !wrong, map[bool]string{true: "mfa_wrong"}[wrong])
	if err == nil && wrong {
		return AccessToken{}, nil, ErrMFAInvalid
	}
	return out, codes, err
}

// MFAVerify: langkah kedua saat login (kode authenticator atau kode cadangan).
func (s *Service) MFAVerify(ctx context.Context, c Claims, code string, cl ClientInfo) (AccessToken, int, error) {
	if err := s.needBox(); err != nil {
		return AccessToken{}, 0, err
	}
	if err := s.checkLock(ctx, c); err != nil {
		return AccessToken{}, 0, err
	}
	var out AccessToken
	backupLeft := -1
	wrong := false
	err := db.WithActor(ctx, s.pool, c.Subject, "admin", "", func(tx pgx.Tx) error {
		var enabled bool
		var sealed []byte
		var last *int64
		if err := tx.QueryRow(ctx, `SELECT mfa_enabled, mfa_secret_enc, mfa_last_step FROM admin_users
			WHERE user_id = $1 AND active AND deleted_at IS NULL FOR UPDATE`, c.Subject).Scan(&enabled, &sealed, &last); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrNotAdmin
			}
			return err
		}
		if !enabled || sealed == nil {
			return ErrMFANotEnrolled
		}
		action := "admin.login"
		clean := strings.TrimSpace(code)
		if len(normBackup(clean)) == 8 {
			// Kode cadangan (berisi huruf).
			tag, err := tx.Exec(ctx, `UPDATE admin_mfa_backup_codes SET used_at = now()
				WHERE id = (SELECT id FROM admin_mfa_backup_codes WHERE user_id = $1 AND code_hash = $2 AND used_at IS NULL LIMIT 1)`,
				c.Subject, backupHash(clean))
			if err != nil {
				return err
			}
			if tag.RowsAffected() == 0 {
				wrong = true
				return s.failMFA(ctx, tx, c)
			}
			if err := tx.QueryRow(ctx, `SELECT count(*) FROM admin_mfa_backup_codes WHERE user_id = $1 AND used_at IS NULL`, c.Subject).Scan(&backupLeft); err != nil {
				return err
			}
			action = "admin.login_backup_code"
		} else {
			secret, err := s.box.Open(sealed, secretAAD(c.Subject))
			if err != nil {
				return err
			}
			step := totpMatch(string(secret), clean, s.now())
			if step < 0 || (last != nil && step <= *last) { // salah, atau kode yang sama dipakai ulang
				wrong = true
				return s.failMFA(ctx, tx, c)
			}
			if _, err := tx.Exec(ctx, `UPDATE admin_users SET mfa_last_step = $2 WHERE user_id = $1`, c.Subject, step); err != nil {
				return err
			}
		}
		if _, err := tx.Exec(ctx, `UPDATE admin_users SET last_login_at = now() WHERE user_id = $1`, c.Subject); err != nil {
			return err
		}
		if err := markVerified(ctx, tx, c, cl, action); err != nil {
			return err
		}
		var err error
		out, err = s.accessFor(ctx, tx, c)
		return err
	})
	s.recordAttempt(ctx, "mfa:"+c.Subject, cl, err == nil && !wrong, map[bool]string{true: "mfa_wrong"}[wrong])
	if err == nil && wrong {
		return AccessToken{}, 0, ErrMFAInvalid
	}
	return out, backupLeft, err
}

// ConfirmTOTP: minta kode authenticator lagi sebelum aksi sensitif (mis. menyimpan kredensial pembayaran).
// Hanya kode dari aplikasi (bukan kode cadangan); kode yang sudah dipakai ditolak; salah berkali-kali → terkunci 15 menit.
func (s *Service) ConfirmTOTP(ctx context.Context, userID, code string, cl ClientInfo) error {
	if s.cfg.AdminMFADevSkip && s.cfg.Env == "development" {
		return nil // mode dev: konfirmasi kode dilewati (lihat config.AdminMFADevSkip)
	}
	if err := s.needBox(); err != nil {
		return err
	}
	if err := s.checkLock(ctx, Claims{Subject: userID}); err != nil {
		return err
	}
	wrong := false
	err := db.WithActor(ctx, s.pool, userID, "admin", "", func(tx pgx.Tx) error {
		var sealed []byte
		var last *int64
		if err := tx.QueryRow(ctx, `SELECT mfa_secret_enc, mfa_last_step FROM admin_users
			WHERE user_id = $1 AND active AND mfa_enabled AND deleted_at IS NULL FOR UPDATE`, userID).Scan(&sealed, &last); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrMFANotEnrolled
			}
			return err
		}
		secret, err := s.box.Open(sealed, secretAAD(userID))
		if err != nil {
			return err
		}
		step := totpMatch(string(secret), strings.TrimSpace(code), s.now())
		if step < 0 || (last != nil && step <= *last) {
			wrong = true
			return nil
		}
		_, err = tx.Exec(ctx, `UPDATE admin_users SET mfa_last_step = $2 WHERE user_id = $1`, userID, step)
		return err
	})
	s.recordAttempt(ctx, "mfa:"+userID, cl, err == nil && !wrong, map[bool]string{true: "mfa_confirm_wrong"}[wrong])
	if err == nil && wrong {
		return ErrMFAInvalid
	}
	return err
}

// failMFA menambah hitungan salah; setelah mfaMaxFailures sesi dicabut (harus login HP lagi).
func (s *Service) failMFA(ctx context.Context, tx pgx.Tx, c Claims) error {
	var n int
	if err := tx.QueryRow(ctx, `UPDATE sessions SET mfa_failures = mfa_failures + 1 WHERE id = $1 RETURNING mfa_failures`, c.SessionID).Scan(&n); err != nil {
		return err
	}
	if n >= mfaMaxFailures {
		_, err := tx.Exec(ctx, `UPDATE sessions SET revoked_at = now() WHERE id = $1`, c.SessionID)
		return err
	}
	return nil
}

func (s *Service) checkLock(ctx context.Context, c Claims) error {
	var n int
	if err := s.pool.QueryRow(ctx, `SELECT count(*) FROM login_attempts WHERE identifier = $1 AND NOT success
		AND occurred_at > now() - interval '15 minutes'`, "mfa:"+c.Subject).Scan(&n); err != nil {
		return err
	}
	if n >= mfaUserFailures {
		return ErrMFALocked
	}
	return nil
}

func markVerified(ctx context.Context, tx pgx.Tx, c Claims, cl ClientInfo, action string) error {
	if _, err := tx.Exec(ctx, `UPDATE sessions SET mfa_verified_at = now(), mfa_failures = 0, last_seen_at = now() WHERE id = $1`, c.SessionID); err != nil {
		return err
	}
	_, err := tx.Exec(ctx, `INSERT INTO audit_logs (actor_type, actor_user_id, action, entity_table, entity_id, ip, user_agent)
		VALUES ('admin', $1, $2, 'sessions', $3, $4, $5)`, c.Subject, action, c.SessionID, cl.ipOrNil(), nullIfEmpty(cl.UserAgent))
	return err
}

func replaceBackupCodes(ctx context.Context, tx pgx.Tx, userID string) ([]string, error) {
	if _, err := tx.Exec(ctx, `UPDATE admin_mfa_backup_codes SET used_at = now() WHERE user_id = $1 AND used_at IS NULL`, userID); err != nil {
		return nil, err
	}
	codes := newBackupCodes(backupCodeCount)
	for _, code := range codes {
		if _, err := tx.Exec(ctx, `INSERT INTO admin_mfa_backup_codes (user_id, code_hash) VALUES ($1, $2)`, userID, backupHash(code)); err != nil {
			return nil, err
		}
	}
	return codes, nil
}

// NewBackupCodes: admin yang sudah lolos MFA membuat ulang kode cadangan (yang lama tidak berlaku).
func (s *Service) NewBackupCodes(ctx context.Context, c Claims) ([]string, error) {
	var codes []string
	err := db.WithActor(ctx, s.pool, c.Subject, "admin", "", func(tx pgx.Tx) error {
		var err error
		codes, err = replaceBackupCodes(ctx, tx, c.Subject)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `INSERT INTO audit_logs (actor_type, actor_user_id, action, entity_table, entity_id)
			VALUES ('admin', $1, 'admin.backup_codes_renewed', 'admin_users', $1)`, c.Subject)
		return err
	})
	return codes, err
}

// ResetMFA: super_admin menghapus MFA admin lain (HP hilang). Semua sesi admin orang itu dicabut.
func ResetMFA(ctx context.Context, tx pgx.Tx, actorID, userID string) error {
	tag, err := tx.Exec(ctx, `UPDATE admin_users SET mfa_enabled = false, mfa_secret_enc = NULL, mfa_pending_enc = NULL,
		mfa_last_step = NULL, mfa_enabled_at = NULL WHERE user_id = $1 AND deleted_at IS NULL`, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	if _, err := tx.Exec(ctx, `UPDATE admin_mfa_backup_codes SET used_at = now() WHERE user_id = $1 AND used_at IS NULL`, userID); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `UPDATE sessions SET revoked_at = now() WHERE user_id = $1 AND audience = 'admin' AND revoked_at IS NULL`, userID); err != nil {
		return err
	}
	actorType := "admin"
	var actor any = actorID
	if actorID == "" {
		actorType, actor = "system", nil
	}
	_, err = tx.Exec(ctx, `INSERT INTO audit_logs (actor_type, actor_user_id, action, entity_table, entity_id)
		VALUES ($1::actor_type, $2, 'admin.mfa_reset', 'admin_users', $3)`, actorType, actor, userID)
	return err
}

// MFARoutes dipasang di /admin/auth/mfa (butuh access token admin; MFA belum perlu lolos).
func (s *Service) MFARoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(s.RequireUser("admin"))
	claims := func(r *http.Request) Claims { c, _ := UserFrom(r.Context()); return c }
	r.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		c := claims(r)
		st, err := mfaState(r.Context(), s.pool, c.Subject, c.SessionID)
		if err != nil {
			s.fail(w, err)
			return
		}
		writeJSON(w, http.StatusOK, st)
	})
	r.Post("/setup", func(w http.ResponseWriter, r *http.Request) {
		secret, url, err := s.MFASetup(r.Context(), claims(r))
		if err != nil {
			s.fail(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"secret": secret, "otpauth_url": url})
	})
	r.Post("/enable", func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			Code string `json:"code"`
		}
		if !decode(w, r, &in) {
			return
		}
		tok, codes, err := s.MFAEnable(r.Context(), claims(r), in.Code, client(r, "web"))
		if err != nil {
			s.fail(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"session": tok, "backup_codes": codes})
	})
	r.Post("/verify", func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			Code string `json:"code"`
		}
		if !decode(w, r, &in) {
			return
		}
		tok, left, err := s.MFAVerify(r.Context(), claims(r), in.Code, client(r, "web"))
		if err != nil {
			s.fail(w, err)
			return
		}
		out := map[string]any{"session": tok}
		if left >= 0 {
			out["backup_codes_left"] = left
		}
		writeJSON(w, http.StatusOK, out)
	})
	r.With(s.RequireMFA).Post("/backup-codes", func(w http.ResponseWriter, r *http.Request) {
		codes, err := s.NewBackupCodes(r.Context(), claims(r))
		if err != nil {
			s.fail(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"backup_codes": codes})
	})
	return r
}
