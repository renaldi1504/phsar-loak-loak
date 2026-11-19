package auth

import (
	"context"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5"
)

// Querier: *pgxpool.Pool atau pgx.Tx.
type Querier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

var (
	uuidLike  = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	tgName    = regexp.MustCompile(`^@[A-Za-z0-9_]{4,32}$`)
	tgNumeric = regexp.MustCompile(`^tg:[0-9]{3,20}$`)
)

// IsTelegramHandle: "@username" Telegram yang valid.
func IsTelegramHandle(s string) bool { return tgName.MatchString(s) }

// FindUser mencari akun dari nomor HP, @username Telegram, tg:<id Telegram>, atau user_id.
// Akun yang login hanya lewat Telegram tidak punya nomor HP, jadi admin bisa ditambahkan lewat @username.
// Mengembalikan pgx.ErrNoRows kalau tidak ada.
func FindUser(ctx context.Context, q Querier, who string) (userID, normalized string, err error) {
	who = strings.TrimSpace(who)
	var sql string
	switch {
	case uuidLike.MatchString(who):
		sql = `SELECT id FROM users WHERE id = $1::uuid AND deleted_at IS NULL`
	case tgName.MatchString(who):
		sql = `SELECT u.id FROM auth_identities ai JOIN users u ON u.id = ai.user_id
			WHERE ai.provider = 'telegram' AND lower(ai.display) = lower($1) AND u.deleted_at IS NULL
			ORDER BY ai.last_used_at DESC NULLS LAST LIMIT 1`
	case tgNumeric.MatchString(who):
		sql = `SELECT u.id FROM auth_identities ai JOIN users u ON u.id = ai.user_id
			WHERE ai.provider = 'telegram' AND ai.provider_uid = substr($1, 4) AND u.deleted_at IS NULL`
	default:
		phone, perr := NormalizePhone(who)
		if perr != nil {
			return "", who, perr
		}
		who, sql = phone, `SELECT id FROM users WHERE phone = $1 AND deleted_at IS NULL`
	}
	err = q.QueryRow(ctx, sql, who).Scan(&userID)
	return userID, who, err
}
