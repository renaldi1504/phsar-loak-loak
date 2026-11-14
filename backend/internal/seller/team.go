package seller

// Staf toko: owner membuat link undangan (berlaku 7 hari, sekali pakai), calon staf login lalu menerimanya.
// Staf hanya bisa mengelola produk; profil toko, KYC, titik penjemputan, dan tim hanya untuk owner.

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	"phsar/internal/notify"
	"phsar/internal/platform/db"
)

const (
	maxStaff      = 10
	inviteTTL     = 7 * 24 * time.Hour
	maxOpenInvite = 10
)

var (
	ErrInviteInvalid = errors.New("link undangan tidak berlaku (sudah dipakai, dibatalkan, atau kedaluwarsa)")
	ErrTooManyStaff  = errors.New("maksimal 10 staf per toko")
	ErrCannotRemove  = errors.New("pemilik toko tidak bisa dihapus")
)

type Member struct {
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"` // disamarkan: +85597***4560
	Telegram  string    `json:"telegram,omitempty"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type Invite struct {
	ID        string    `json:"id"`
	Token     string    `json:"token,omitempty"` // hanya dikembalikan sekali, saat dibuat
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type Team struct {
	Members []Member `json:"members"`
	Invites []Invite `json:"invites"`
}

func maskPhone(p string) string {
	if len(p) < 8 {
		return p
	}
	return p[:6] + "***" + p[len(p)-4:]
}

func hashToken(t string) string {
	h := sha256.Sum256([]byte("store-invite:" + t))
	return hex.EncodeToString(h[:])
}

func (s *Service) GetTeam(ctx context.Context, storeID string) (Team, error) {
	team := Team{Members: []Member{}, Invites: []Invite{}}
	rows, err := s.pool.Query(ctx, `
		SELECT u.id, u.name, coalesce(u.phone, ''),
		       coalesce((SELECT display FROM auth_identities ai WHERE ai.user_id = u.id AND ai.provider = 'telegram' LIMIT 1), ''),
		       sm.role, sm.created_at
		FROM store_members sm JOIN users u ON u.id = sm.user_id
		WHERE sm.store_id = $1 ORDER BY sm.role = 'owner' DESC, sm.created_at`, storeID)
	if err != nil {
		return team, err
	}
	for rows.Next() {
		var m Member
		if err := rows.Scan(&m.UserID, &m.Name, &m.Phone, &m.Telegram, &m.Role, &m.CreatedAt); err != nil {
			rows.Close()
			return team, err
		}
		m.Phone = maskPhone(m.Phone)
		team.Members = append(team.Members, m)
	}
	rows.Close()
	rows, err = s.pool.Query(ctx, `SELECT id, expires_at, created_at FROM store_invites
		WHERE store_id = $1 AND accepted_at IS NULL AND revoked_at IS NULL AND expires_at > now() ORDER BY created_at DESC`, storeID)
	if err != nil {
		return team, err
	}
	defer rows.Close()
	for rows.Next() {
		var i Invite
		if err := rows.Scan(&i.ID, &i.ExpiresAt, &i.CreatedAt); err != nil {
			return team, err
		}
		team.Invites = append(team.Invites, i)
	}
	return team, rows.Err()
}

// CreateInvite: token acak 32 byte; yang disimpan hanya hash-nya.
func (s *Service) CreateInvite(ctx context.Context, storeID, userID string) (Invite, error) {
	var staff, open int
	if err := s.pool.QueryRow(ctx, `SELECT
		(SELECT count(*) FROM store_members WHERE store_id = $1 AND role = 'staff'),
		(SELECT count(*) FROM store_invites WHERE store_id = $1 AND accepted_at IS NULL AND revoked_at IS NULL AND expires_at > now())`,
		storeID).Scan(&staff, &open); err != nil {
		return Invite{}, err
	}
	if staff >= maxStaff {
		return Invite{}, ErrTooManyStaff
	}
	if open >= maxOpenInvite {
		return Invite{}, invalid("invites", "terlalu banyak link aktif; batalkan yang lama dulu")
	}
	b := make([]byte, 24)
	_, _ = rand.Read(b)
	inv := Invite{Token: hex.EncodeToString(b)}
	err := db.WithActor(ctx, s.pool, userID, "seller", "", func(tx pgx.Tx) error {
		return tx.QueryRow(ctx, `INSERT INTO store_invites (store_id, token_hash, created_by, expires_at)
			VALUES ($1, $2, $3, now() + make_interval(secs => $4)) RETURNING id, expires_at, created_at`,
			storeID, hashToken(inv.Token), userID, inviteTTL.Seconds()).Scan(&inv.ID, &inv.ExpiresAt, &inv.CreatedAt)
	})
	return inv, err
}

func (s *Service) RevokeInvite(ctx context.Context, storeID, inviteID, userID string) error {
	return db.WithActor(ctx, s.pool, userID, "seller", "", func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, `UPDATE store_invites SET revoked_at = now() WHERE id = $1 AND store_id = $2 AND accepted_at IS NULL AND revoked_at IS NULL`,
			inviteID, storeID)
		if err == nil && tag.RowsAffected() == 0 {
			return ErrNotFound
		}
		return err
	})
}

// RemoveMember: owner menghapus staf, atau staf keluar sendiri.
func (s *Service) RemoveMember(ctx context.Context, storeID, targetUserID, actorID, actorRole string) error {
	if actorRole != "owner" && targetUserID != actorID {
		return ErrOwnerOnly
	}
	return db.WithActor(ctx, s.pool, actorID, "seller", "", func(tx pgx.Tx) error {
		var role string
		err := tx.QueryRow(ctx, `SELECT role FROM store_members WHERE store_id = $1 AND user_id = $2`, storeID, targetUserID).Scan(&role)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		if role == "owner" {
			return ErrCannotRemove
		}
		_, err = tx.Exec(ctx, `DELETE FROM store_members WHERE store_id = $1 AND user_id = $2`, storeID, targetUserID)
		return err
	})
}

type InviteInfo struct {
	StoreID   string `json:"store_id"`
	StoreName string `json:"store_name"`
	LogoURL   string `json:"logo_url,omitempty"`
	Member    bool   `json:"already_member"`
}

func (s *Service) findInvite(ctx context.Context, q interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}, token string, lock bool) (inviteID, storeID string, err error) {
	sql := `SELECT i.id, i.store_id FROM store_invites i JOIN stores st ON st.id = i.store_id AND st.deleted_at IS NULL
		WHERE i.token_hash = $1 AND i.accepted_at IS NULL AND i.revoked_at IS NULL AND i.expires_at > now()`
	if lock {
		sql += " FOR UPDATE OF i"
	}
	err = q.QueryRow(ctx, sql, hashToken(token)).Scan(&inviteID, &storeID)
	if errors.Is(err, pgx.ErrNoRows) {
		err = ErrInviteInvalid
	}
	return
}

// InviteDetails: untuk halaman "Kamu diundang ke toko X".
func (s *Service) InviteDetails(ctx context.Context, token, userID string) (InviteInfo, error) {
	_, storeID, err := s.findInvite(ctx, s.pool, token, false)
	if err != nil {
		return InviteInfo{}, err
	}
	var info InviteInfo
	var logo *string
	if err := s.pool.QueryRow(ctx, `SELECT id, name, logo_key, EXISTS (SELECT 1 FROM store_members WHERE store_id = $1 AND user_id = $2)
		FROM stores WHERE id = $1`, storeID, userID).Scan(&info.StoreID, &info.StoreName, &logo, &info.Member); err != nil {
		return info, err
	}
	info.LogoURL = s.files.PublicURL(deref(logo))
	return info, nil
}

// AcceptInvite: user yang sedang login menjadi staf. Link langsung tidak berlaku lagi.
func (s *Service) AcceptInvite(ctx context.Context, token, userID string) (string, error) {
	var storeID string
	err := db.WithActor(ctx, s.pool, userID, "seller", "", func(tx pgx.Tx) error {
		inviteID, sid, err := s.findInvite(ctx, tx, token, true)
		if err != nil {
			return err
		}
		storeID = sid
		var role string
		err = tx.QueryRow(ctx, `SELECT role FROM store_members WHERE store_id = $1 AND user_id = $2`, sid, userID).Scan(&role)
		if err == nil {
			return nil // sudah anggota (mis. owner membuka link sendiri): link tidak dipakai
		}
		if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
		var staff int
		if err := tx.QueryRow(ctx, `SELECT count(*) FROM store_members WHERE store_id = $1 AND role = 'staff'`, sid).Scan(&staff); err != nil {
			return err
		}
		if staff >= maxStaff {
			return ErrTooManyStaff
		}
		if _, err := tx.Exec(ctx, `INSERT INTO store_members (store_id, user_id, role) VALUES ($1, $2, 'staff')`, sid, userID); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `UPDATE store_invites SET accepted_by = $2, accepted_at = now() WHERE id = $1`, inviteID, userID); err != nil {
			return err
		}
		var person, storeName string
		var name, phone string
		if err := tx.QueryRow(ctx, `SELECT u.name, coalesce(u.phone, ''), st.name FROM users u, stores st WHERE u.id = $1 AND st.id = $2`,
			userID, sid).Scan(&name, &phone, &storeName); err != nil {
			return err
		}
		person = name
		if person == "" {
			person = maskPhone(phone)
		}
		return notify.Enqueue(ctx, tx, sid, notify.StaffJoined, notify.Params{Person: person, Store: storeName}, "#/store/"+sid)
	})
	return storeID, err
}
