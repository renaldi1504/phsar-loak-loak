package admin

// Pengelolaan akun admin (super_admin): tambah staf, ganti peran, nonaktifkan, reset MFA.

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

	"phsar/internal/auth"
	"phsar/internal/platform/db"
)

var uuidRe = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

var adminRoles = map[string]bool{"super_admin": true, "moderator": true, "cs": true, "finance": true, "ads_reviewer": true}

type adminRow struct {
	UserID      string     `json:"user_id"`
	Name        string     `json:"name"`
	Phone       string     `json:"phone"`
	Role        string     `json:"role"`
	Active      bool       `json:"active"`
	MFAEnabled  bool       `json:"mfa_enabled"`
	MFAAt       *time.Time `json:"mfa_enabled_at"`
	LastLoginAt *time.Time `json:"last_login_at"`
	BackupLeft  int        `json:"backup_codes_left"`
	Self        bool       `json:"self"`
}

// GET /admin/api/admins
func (h *Handler) listAdmins(w http.ResponseWriter, r *http.Request) {
	me := actor(r)
	rows, err := h.pool.Query(r.Context(), `SELECT a.user_id::text, coalesce(u.name, ''),
			coalesce(nullif(u.phone, ''), (SELECT ai.display FROM auth_identities ai WHERE ai.user_id = u.id AND ai.provider = 'telegram' AND ai.display IS NOT NULL LIMIT 1), ''), a.role_code, a.active,
			a.mfa_enabled, a.mfa_enabled_at, a.last_login_at,
			(SELECT count(*) FROM admin_mfa_backup_codes b WHERE b.user_id = a.user_id AND b.used_at IS NULL)
		FROM admin_users a JOIN users u ON u.id = a.user_id
		WHERE a.deleted_at IS NULL ORDER BY a.active DESC, a.role_code, u.name`)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "internal", err.Error())
		return
	}
	defer rows.Close()
	out := []adminRow{}
	for rows.Next() {
		var a adminRow
		if err := rows.Scan(&a.UserID, &a.Name, &a.Phone, &a.Role, &a.Active, &a.MFAEnabled, &a.MFAAt, &a.LastLoginAt, &a.BackupLeft); err != nil {
			writeErr(w, http.StatusInternalServerError, "internal", err.Error())
			return
		}
		a.Self = a.UserID == me
		out = append(out, a)
	}
	writeJSON(w, http.StatusOK, map[string]any{"admins": out})
}

// POST /admin/api/admins {phone, role} — user harus sudah pernah login sekali (app, Seller Center, atau halaman admin).
func (h *Handler) addAdmin(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Phone string `json:"phone"`
		Role  string `json:"role"`
	}
	if !decode(w, r, &in) {
		return
	}
	in.Phone = strings.TrimSpace(in.Phone)
	if strings.HasPrefix(in.Phone, "@") {
		if !auth.IsTelegramHandle(in.Phone) {
			writeFieldErr(w, "phone")
			return
		}
	} else if _, err := auth.NormalizePhone(in.Phone); err != nil {
		writeFieldErr(w, "phone")
		return
	}
	if !adminRoles[in.Role] {
		writeFieldErr(w, "role")
		return
	}
	err := db.WithActor(r.Context(), h.pool, actor(r), "admin", "", func(tx pgx.Tx) error {
		uid, _, err := auth.FindUser(r.Context(), tx, in.Phone)
		if err != nil {
			return err
		}
		_, err = tx.Exec(r.Context(), `INSERT INTO admin_users (user_id, role_code) VALUES ($1, $2)
			ON CONFLICT (user_id) DO UPDATE SET role_code = EXCLUDED.role_code, active = true, deleted_at = NULL`, uid, in.Role)
		return err
	})
	if errors.Is(err, pgx.ErrNoRows) {
		writeErr(w, http.StatusNotFound, "user_not_found", "nomor / akun Telegram ini belum pernah login di Phsar")
		return
	}
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "internal", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// PATCH /admin/api/admins/{userID} {role?, active?} — tidak bisa mengubah akun sendiri.
func (h *Handler) updateAdmin(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "userID")
	if !uuidRe.MatchString(id) {
		writeErr(w, http.StatusNotFound, "not_found", "admin tidak ditemukan")
		return
	}
	if id == actor(r) {
		writeErr(w, http.StatusConflict, "self", "tidak bisa mengubah akun sendiri")
		return
	}
	var in struct {
		Role   *string `json:"role"`
		Active *bool   `json:"active"`
	}
	if !decode(w, r, &in) {
		return
	}
	if in.Role != nil && !adminRoles[*in.Role] {
		writeFieldErr(w, "role")
		return
	}
	err := db.WithActor(r.Context(), h.pool, actor(r), "admin", "", func(tx pgx.Tx) error {
		tag, err := tx.Exec(r.Context(), `UPDATE admin_users SET role_code = coalesce($2, role_code), active = coalesce($3, active)
			WHERE user_id = $1 AND deleted_at IS NULL`, id, in.Role, in.Active)
		if err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			return pgx.ErrNoRows
		}
		if in.Active != nil && !*in.Active { // dinonaktifkan: tutup semua sesi admin-nya
			_, err = tx.Exec(r.Context(), `UPDATE sessions SET revoked_at = now() WHERE user_id = $1 AND audience = 'admin' AND revoked_at IS NULL`, id)
		}
		return err
	})
	if errors.Is(err, pgx.ErrNoRows) {
		writeErr(w, http.StatusNotFound, "not_found", "admin tidak ditemukan")
		return
	}
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "internal", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// POST /admin/api/admins/{userID}/mfa-reset — admin yang kehilangan HP memasang MFA lagi saat login berikutnya.
func (h *Handler) resetAdminMFA(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "userID")
	if !uuidRe.MatchString(id) {
		writeErr(w, http.StatusNotFound, "not_found", "admin tidak ditemukan")
		return
	}
	if id == actor(r) {
		writeErr(w, http.StatusConflict, "self", "pakai kode cadangan untuk akun sendiri")
		return
	}
	err := db.WithActor(r.Context(), h.pool, actor(r), "admin", "", func(tx pgx.Tx) error {
		return auth.ResetMFA(r.Context(), tx, actor(r), id)
	})
	if errors.Is(err, pgx.ErrNoRows) {
		writeErr(w, http.StatusNotFound, "not_found", "admin tidak ditemukan")
		return
	}
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "internal", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func writeFieldErr(w http.ResponseWriter, field string) {
	writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"error": map[string]string{"code": "validation", "field": field, "message": field}})
}
