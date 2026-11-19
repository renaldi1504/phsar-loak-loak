// Package payout: saldo seller dan pembayaran (payout) terjadwal.
//
//  1. Pesanan terantar → ditahan N hari (hold_days, untuk komplain/retur) → "completed":
//     seller_amount masuk saldo (seller_ledger, kind order_earning).
//  2. Sesuai jadwal admin (harian / mingguan, jam tertentu, waktu Phnom Penh) dibuat payout "pending"
//     per toko untuk seluruh saldo ≥ minimal; saldo langsung dipotong (kind payout).
//  3. Admin mentransfer manual ke rekening seller, lalu mencatat nomor referensi → "paid".
//     Dibatalkan → saldo dikembalikan (kind payout_cancelled).
//
// Nomor rekening disimpan terenkripsi (AES-GCM, terikat ke id toko) dan hanya dibuka untuk admin keuangan.
package payout

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"phsar/internal/jobs"
	"phsar/internal/notify"
	"phsar/internal/paging"
	"phsar/internal/platform/db"
	"phsar/internal/secretbox"
)

var (
	ErrValidation = errors.New("validation")
	ErrNotFound   = errors.New("not found")
	ErrBadStatus  = errors.New("payout status does not allow this")
)

type fieldError struct{ field, msg string }

func (e fieldError) Error() string { return e.field + ": " + e.msg }
func (e fieldError) Unwrap() error { return ErrValidation }

// AccountChangeHold: ganti rekening = payout ditahan 48 jam (melindungi dari pengambilalihan akun).
const AccountChangeHold = 48 * time.Hour

var phnomPenh = func() *time.Location {
	l, err := time.LoadLocation("Asia/Phnom_Penh")
	if err != nil {
		return time.FixedZone("ICT", 7*3600)
	}
	return l
}()

type Service struct {
	Pool *pgxpool.Pool
	Box  *secretbox.Box
	Log  *slog.Logger
}

// ---------------------------------------------------------------------------
// Pengaturan (admin)
// ---------------------------------------------------------------------------

type Settings struct {
	Frequency string `json:"frequency"`  // daily · weekly
	Weekday   int    `json:"weekday"`    // 0 = Minggu … 6 = Sabtu (untuk weekly)
	Hour      int    `json:"hour"`       // jam dibuat (waktu Phnom Penh)
	HoldDays  int    `json:"hold_days"`  // hari setelah terantar sebelum masuk saldo
	MinAmount string `json:"min_amount"` // saldo minimal untuk dibayar
}

func (s *Service) Settings(ctx context.Context) (Settings, error) {
	st := Settings{Frequency: "weekly", Weekday: 1, Hour: 9, HoldDays: 3, MinAmount: "5.00"}
	var raw []byte
	err := s.Pool.QueryRow(ctx, `SELECT value FROM platform_settings WHERE key = 'payout'`).Scan(&raw)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return st, err
	}
	if len(raw) > 0 {
		_ = json.Unmarshal(raw, &st)
	}
	return st, nil
}

func (s *Service) SaveSettings(ctx context.Context, adminID string, st Settings) (Settings, error) {
	switch {
	case st.Frequency != "daily" && st.Frequency != "weekly":
		return st, fieldError{"frequency", "choose daily or weekly"}
	case st.Weekday < 0 || st.Weekday > 6:
		return st, fieldError{"weekday", "choose a day"}
	case st.Hour < 0 || st.Hour > 23:
		return st, fieldError{"hour", "0–23"}
	case st.HoldDays < 0 || st.HoldDays > 30:
		return st, fieldError{"hold_days", "0–30 days"}
	}
	min, err := decimal.NewFromString(strings.TrimSpace(st.MinAmount))
	if err != nil || min.IsNegative() || min.GreaterThan(decimal.NewFromInt(1000)) {
		return st, fieldError{"min_amount", "enter an amount between 0 and 1000"}
	}
	st.MinAmount = min.StringFixed(2)
	raw, _ := json.Marshal(st)
	err = db.WithActor(ctx, s.Pool, adminID, "admin", "", func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `INSERT INTO platform_settings (key, value, updated_by) VALUES ('payout', $1, $2)
			ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_by = EXCLUDED.updated_by, updated_at = now()`, raw, adminID)
		return err
	})
	return st, err
}

// NextRun: waktu jadwal berikutnya (untuk ditampilkan ke seller & admin).
func (st Settings) NextRun(now time.Time) time.Time {
	l := now.In(phnomPenh)
	at := time.Date(l.Year(), l.Month(), l.Day(), st.Hour, 0, 0, 0, phnomPenh)
	for i := 0; i < 8; i++ {
		c := at.AddDate(0, 0, i)
		if c.After(l) && (st.Frequency == "daily" || int(c.Weekday()) == st.Weekday) {
			return c
		}
	}
	return at.AddDate(0, 0, 7)
}

// ---------------------------------------------------------------------------
// Pesanan selesai → saldo
// ---------------------------------------------------------------------------

// CompleteDue: pesanan terantar yang sudah lewat masa tahan → completed + saldo seller.
func (s *Service) CompleteDue(ctx context.Context) (int, error) {
	st, err := s.Settings(ctx)
	if err != nil {
		return 0, err
	}
	n := 0
	err = db.WithActor(ctx, s.Pool, "", "system", "", func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, `SELECT o.id, o.store_id, o.seller_amount, o.order_no FROM orders o JOIN deliveries d ON d.order_id = o.id
			WHERE o.status = 'delivered' AND d.delivered_at <= now() - make_interval(days => $1)
			  AND NOT EXISTS (SELECT 1 FROM returns r WHERE r.order_id = o.id AND r.status IN ('requested','disputed','approved')) -- retur terbuka: tahan dulu
			ORDER BY d.delivered_at LIMIT 500 FOR UPDATE OF o`, st.HoldDays)
		if err != nil {
			return err
		}
		type row struct {
			id, store, no string
			amt           decimal.Decimal
		}
		var list []row
		for rows.Next() {
			var r row
			if err := rows.Scan(&r.id, &r.store, &r.amt, &r.no); err != nil {
				rows.Close()
				return err
			}
			list = append(list, r)
		}
		rows.Close()
		for _, r := range list {
			if _, err := tx.Exec(ctx, `UPDATE orders SET status = 'completed', completed_at = now() WHERE id = $1`, r.id); err != nil {
				return err
			}
			if _, err := tx.Exec(ctx, `INSERT INTO order_events (order_id, status, actor_type) VALUES ($1, 'completed', 'system')`, r.id); err != nil {
				return err
			}
			if r.amt.IsPositive() {
				if _, err := tx.Exec(ctx, `INSERT INTO seller_ledger (store_id, kind, amount, order_id, note) VALUES ($1, 'order_earning', $2, $3, $4)
					ON CONFLICT DO NOTHING`, r.store, r.amt, r.id, r.no); err != nil {
					return err
				}
			}
			n++
		}
		return nil
	})
	return n, err
}

// ---------------------------------------------------------------------------
// Rekening seller
// ---------------------------------------------------------------------------

type Account struct {
	ID          string     `json:"id"`
	Bank        string     `json:"bank"`
	AccountName string     `json:"account_name"`
	Last4       string     `json:"account_last4"`
	HoldUntil   *time.Time `json:"payout_hold_until,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	Number      string     `json:"account_number,omitempty"` // hanya untuk admin keuangan
}

var banks = map[string]bool{"aba": true, "acleda": true, "wing": true, "bakong": true, "other": true}

func digits(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		} else if r != ' ' && r != '-' {
			return ""
		}
	}
	return b.String()
}

func (s *Service) account(ctx context.Context, q interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}, storeID string) (*Account, []byte, error) {
	var a Account
	var enc []byte
	err := q.QueryRow(ctx, `SELECT id, bank, account_name, account_last4, payout_hold_until, created_at, account_number_enc
		FROM seller_payout_accounts WHERE store_id = $1 AND deleted_at IS NULL ORDER BY is_default DESC, created_at DESC LIMIT 1`, storeID).
		Scan(&a.ID, &a.Bank, &a.AccountName, &a.Last4, &a.HoldUntil, &a.CreatedAt, &enc)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil, nil
	}
	return &a, enc, err
}

// SaveAccount (pemilik toko): rekening baru menggantikan yang lama. Ganti rekening → payout ditahan 48 jam.
func (s *Service) SaveAccount(ctx context.Context, storeID, userID, bank, name, number string) (*Account, error) {
	name = strings.TrimSpace(name)
	num := digits(strings.TrimSpace(number))
	switch {
	case !banks[bank]:
		return nil, fieldError{"bank", "choose a bank"}
	case len([]rune(name)) < 2 || len([]rune(name)) > 120:
		return nil, fieldError{"account_name", "enter the name on the account"}
	case len(num) < 6 || len(num) > 20:
		return nil, fieldError{"account_number", "enter the account number (digits only)"}
	}
	enc, err := s.Box.Seal([]byte(num), []byte(storeID))
	if err != nil {
		return nil, err
	}
	err = db.WithActor(ctx, s.Pool, userID, "seller", "", func(tx pgx.Tx) error {
		old, _, err := s.account(ctx, tx, storeID)
		if err != nil {
			return err
		}
		var hold any
		if old != nil {
			hold = time.Now().Add(AccountChangeHold)
			if _, err := tx.Exec(ctx, `UPDATE seller_payout_accounts SET deleted_at = now(), is_default = false WHERE store_id = $1 AND deleted_at IS NULL`, storeID); err != nil {
				return err
			}
		}
		if _, err := tx.Exec(ctx, `INSERT INTO seller_payout_accounts (store_id, bank, account_name, account_number_enc, account_last4, is_default, payout_hold_until)
			VALUES ($1, $2, $3, $4, $5, true, $6)`, storeID, bank, name, enc, num[len(num)-4:], hold); err != nil {
			return err
		}
		// Tabel rekening sengaja tanpa audit otomatis (nomor terenkripsi tidak boleh tersalin) → catat aksinya saja.
		_, err = tx.Exec(ctx, `INSERT INTO audit_logs (actor_type, actor_user_id, store_id, action, entity_table, new_data)
			VALUES ('seller', $1, $2, 'payout_account.change', 'seller_payout_accounts', $3)`, userID, storeID,
			map[string]string{"bank": bank, "last4": num[len(num)-4:]})
		return err
	})
	if err != nil {
		return nil, err
	}
	a, _, err := s.account(ctx, s.Pool, storeID)
	return a, err
}

// ---------------------------------------------------------------------------
// Ringkasan untuk seller
// ---------------------------------------------------------------------------

type PayoutRow struct {
	ID          string     `json:"id"`
	StoreID     string     `json:"store_id,omitempty"`
	StoreName   string     `json:"store_name,omitempty"`
	OwnerPhone  string     `json:"owner_phone,omitempty"`
	Amount      string     `json:"amount"`
	Status      string     `json:"status"`
	Bank        string     `json:"bank"`
	AccountName string     `json:"account_name"`
	Last4       string     `json:"account_last4"`
	Number      string     `json:"account_number,omitempty"` // admin keuangan
	OrderCount  int        `json:"order_count"`
	Reference   string     `json:"reference,omitempty"`
	Note        string     `json:"note,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	PaidAt      *time.Time `json:"paid_at,omitempty"`
	ByRequest   bool       `json:"requested_by_seller"` // extras: dibuat dari tombol "Tarik sekarang" seller
}

type LedgerRow struct {
	Kind    string    `json:"kind"`
	Amount  string    `json:"amount"`
	OrderNo string    `json:"order_no,omitempty"`
	Note    string    `json:"note,omitempty"`
	At      time.Time `json:"at"`
}

type Summary struct {
	Available     string      `json:"available"`      // saldo siap dibayar
	OnHold        string      `json:"on_hold"`        // terantar, masih masa tahan
	Upcoming      string      `json:"upcoming"`       // sudah dibayar pembeli, belum terantar
	PendingPayout string      `json:"pending_payout"` // sedang ditransfer admin
	HoldDays      int         `json:"hold_days"`
	MinAmount     string      `json:"min_amount"`
	Frequency     string      `json:"frequency"`
	Weekday       int         `json:"weekday"`
	NextRun       time.Time   `json:"next_run"`
	Account       *Account    `json:"account"`
	Payouts       []PayoutRow `json:"payouts"`
	Ledger        []LedgerRow `json:"ledger"`
}

func (s *Service) balance(ctx context.Context, q interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}, storeID string) (decimal.Decimal, error) {
	var b decimal.Decimal
	err := q.QueryRow(ctx, `SELECT coalesce(sum(amount), 0) FROM seller_ledger WHERE store_id = $1`, storeID).Scan(&b)
	return b, err
}

func (s *Service) StoreSummary(ctx context.Context, storeID string) (*Summary, error) {
	st, err := s.Settings(ctx)
	if err != nil {
		return nil, err
	}
	sum := &Summary{HoldDays: st.HoldDays, MinAmount: st.MinAmount, Frequency: st.Frequency, Weekday: st.Weekday, NextRun: st.NextRun(time.Now())}
	bal, err := s.balance(ctx, s.Pool, storeID)
	if err != nil {
		return nil, err
	}
	sum.Available = bal.StringFixed(2)
	var hold, up, pend decimal.Decimal
	if err := s.Pool.QueryRow(ctx, `SELECT
		  coalesce(sum(seller_amount) FILTER (WHERE status = 'delivered'), 0),
		  coalesce(sum(seller_amount) FILTER (WHERE status IN ('paid','accepted','ready_for_pickup','picked_up')), 0)
		FROM orders WHERE store_id = $1`, storeID).Scan(&hold, &up); err != nil {
		return nil, err
	}
	if err := s.Pool.QueryRow(ctx, `SELECT coalesce(sum(amount), 0) FROM payouts WHERE store_id = $1 AND status = 'pending'`, storeID).Scan(&pend); err != nil {
		return nil, err
	}
	sum.OnHold, sum.Upcoming, sum.PendingPayout = hold.StringFixed(2), up.StringFixed(2), pend.StringFixed(2)
	if sum.Account, _, err = s.account(ctx, s.Pool, storeID); err != nil {
		return nil, err
	}
	if sum.Payouts, err = s.payouts(ctx, `p.store_id = $1`, 30, false, storeID); err != nil {
		return nil, err
	}
	rows, err := s.Pool.Query(ctx, `SELECT l.kind, l.amount::text, coalesce(o.order_no, ''), coalesce(l.note, ''), l.created_at
		FROM seller_ledger l LEFT JOIN orders o ON o.id = l.order_id WHERE l.store_id = $1 ORDER BY l.id DESC LIMIT 50`, storeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	sum.Ledger = []LedgerRow{}
	for rows.Next() {
		var l LedgerRow
		if err := rows.Scan(&l.Kind, &l.Amount, &l.OrderNo, &l.Note, &l.At); err != nil {
			return nil, err
		}
		sum.Ledger = append(sum.Ledger, l)
	}
	return sum, rows.Err()
}

// Ledger: riwayat saldo toko per halaman (terbaru dulu).
func (s *Service) Ledger(ctx context.Context, storeID string, pg paging.Params) ([]LedgerRow, paging.Meta, error) {
	var total int
	if err := s.Pool.QueryRow(ctx, `SELECT count(*) FROM seller_ledger WHERE store_id = $1`, storeID).Scan(&total); err != nil {
		return nil, paging.Meta{}, err
	}
	rows, err := s.Pool.Query(ctx, `SELECT l.kind, l.amount::text, coalesce(o.order_no, ''), coalesce(l.note, ''), l.created_at
		FROM seller_ledger l LEFT JOIN orders o ON o.id = l.order_id WHERE l.store_id = $1 ORDER BY l.id DESC LIMIT $2 OFFSET $3`,
		storeID, pg.Size, pg.Offset())
	if err != nil {
		return nil, paging.Meta{}, err
	}
	defer rows.Close()
	out := []LedgerRow{}
	for rows.Next() {
		var l LedgerRow
		if err := rows.Scan(&l.Kind, &l.Amount, &l.OrderNo, &l.Note, &l.At); err != nil {
			return nil, paging.Meta{}, err
		}
		out = append(out, l)
	}
	return out, paging.NewMeta(pg, total), rows.Err()
}

func (s *Service) payouts(ctx context.Context, where string, limit int, withNumber bool, args ...any) ([]PayoutRow, error) {
	rows, err := s.Pool.Query(ctx, `
		SELECT p.id, p.store_id, st.name, coalesce((SELECT u.phone FROM store_members m JOIN users u ON u.id = m.user_id
		         WHERE m.store_id = p.store_id AND m.role = 'owner' LIMIT 1), ''),
		       p.amount::text, p.status, p.bank, p.account_name, p.account_last4, p.order_count, coalesce(p.reference, ''), coalesce(p.note, ''),
		       p.created_at, p.paid_at, a.account_number_enc, p.requested_by IS NOT NULL
		FROM payouts p JOIN stores st ON st.id = p.store_id LEFT JOIN seller_payout_accounts a ON a.id = p.account_id
		WHERE `+where+` ORDER BY p.created_at DESC LIMIT `+fmt.Sprint(limit), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []PayoutRow{}
	for rows.Next() {
		var p PayoutRow
		var enc []byte
		if err := rows.Scan(&p.ID, &p.StoreID, &p.StoreName, &p.OwnerPhone, &p.Amount, &p.Status, &p.Bank, &p.AccountName, &p.Last4, &p.OrderCount,
			&p.Reference, &p.Note, &p.CreatedAt, &p.PaidAt, &enc, &p.ByRequest); err != nil {
			return nil, err
		}
		if withNumber && len(enc) > 0 {
			if b, err := s.Box.Open(enc, []byte(p.StoreID)); err == nil {
				p.Number = string(b)
			} else {
				s.Log.Error("nomor rekening tidak bisa dibuka (DATA_ENCRYPTION_KEY berubah?)", "payout", p.ID)
				p.Number = "?"
			}
		}
		if !withNumber {
			p.OwnerPhone = ""
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// ---------------------------------------------------------------------------
// Membuat payout (jadwal / tombol admin)
// ---------------------------------------------------------------------------

type RunResult struct {
	Created int            `json:"created"`
	Total   string         `json:"total"`
	Skipped []SkippedStore `json:"skipped"`
}

type SkippedStore struct {
	StoreName string `json:"store_name"`
	Balance   string `json:"balance"`
	Reason    string `json:"reason"` // no_account · account_hold · below_min · pending_exists
}

// Generate: satu payout per toko untuk seluruh saldo. adminID kosong = jadwal otomatis.
func (s *Service) Generate(ctx context.Context, adminID string) (*RunResult, error) {
	st, err := s.Settings(ctx)
	if err != nil {
		return nil, err
	}
	min, _ := decimal.NewFromString(st.MinAmount)
	res := &RunResult{Skipped: []SkippedStore{}}
	total := decimal.Zero
	actorType := "system"
	if adminID != "" {
		actorType = "admin"
	}
	rows, err := s.Pool.Query(ctx, `SELECT l.store_id, st.name, sum(l.amount) FROM seller_ledger l JOIN stores st ON st.id = l.store_id
		GROUP BY l.store_id, st.name HAVING sum(l.amount) > 0 ORDER BY st.name`)
	if err != nil {
		return nil, err
	}
	type cand struct {
		id, name string
		bal      decimal.Decimal
	}
	var list []cand
	for rows.Next() {
		var c cand
		if err := rows.Scan(&c.id, &c.name, &c.bal); err != nil {
			rows.Close()
			return nil, err
		}
		list = append(list, c)
	}
	rows.Close()
	for _, c := range list {
		skip := ""
		err := db.WithActor(ctx, s.Pool, adminID, actorType, "", func(tx pgx.Tx) error {
			// Kunci per toko supaya jadwal & tombol admin bersamaan tidak membuat dua payout.
			if _, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtext('payout:' || $1))`, c.id); err != nil {
				return err
			}
			var pending bool
			if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM payouts WHERE store_id = $1 AND status = 'pending')`, c.id).Scan(&pending); err != nil {
				return err
			}
			if pending {
				skip = "pending_exists"
				return nil
			}
			bal, err := s.balance(ctx, tx, c.id)
			if err != nil {
				return err
			}
			if bal.LessThan(min) || !bal.IsPositive() {
				skip = "below_min"
				return nil
			}
			acc, _, err := s.account(ctx, tx, c.id)
			if err != nil {
				return err
			}
			if acc == nil {
				skip = "no_account"
				return nil
			}
			if acc.HoldUntil != nil && acc.HoldUntil.After(time.Now()) {
				skip = "account_hold"
				return nil
			}
			var orders int
			if err := tx.QueryRow(ctx, `SELECT count(*) FROM seller_ledger l WHERE l.store_id = $1 AND l.kind = 'order_earning'
				AND l.id > coalesce((SELECT max(id) FROM seller_ledger WHERE store_id = $1 AND kind = 'payout'), 0)`, c.id).Scan(&orders); err != nil {
				return err
			}
			var pid string
			if err := tx.QueryRow(ctx, `INSERT INTO payouts (store_id, amount, account_id, bank, account_name, account_last4, order_count, created_by)
				VALUES ($1, $2, $3, $4, $5, $6, $7, nullif($8, '')::uuid) RETURNING id`, c.id, bal, acc.ID, acc.Bank, acc.AccountName, acc.Last4, orders, adminID).Scan(&pid); err != nil {
				return err
			}
			if _, err := tx.Exec(ctx, `INSERT INTO seller_ledger (store_id, kind, amount, payout_id, created_by) VALUES ($1, 'payout', $2, $3, nullif($4, '')::uuid)`,
				c.id, bal.Neg(), pid, adminID); err != nil {
				return err
			}
			res.Created++
			total = total.Add(bal)
			return nil
		})
		if err != nil {
			return nil, err
		}
		if skip != "" {
			res.Skipped = append(res.Skipped, SkippedStore{StoreName: c.name, Balance: c.bal.StringFixed(2), Reason: skip})
		}
	}
	res.Total = total.StringFixed(2)
	today := time.Now().In(phnomPenh).Format("2006-01-02")
	_, err = s.Pool.Exec(ctx, `INSERT INTO payout_runs (run_date, payout_count, total, manual, created_by) VALUES ($1, $2, $3, $4, nullif($5, '')::uuid)
		ON CONFLICT DO NOTHING`, today, res.Created, total, adminID != "", adminID)
	return res, err
}

// RunIfDue (worker): jalankan jadwal kalau hari & jamnya sudah tiba dan hari ini belum dijalankan.
func (s *Service) RunIfDue(ctx context.Context, now time.Time) (bool, error) {
	st, err := s.Settings(ctx)
	if err != nil {
		return false, err
	}
	l := now.In(phnomPenh)
	if l.Hour() < st.Hour || (st.Frequency == "weekly" && int(l.Weekday()) != st.Weekday) {
		return false, nil
	}
	var done bool
	if err := s.Pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM payout_runs WHERE run_date = $1 AND NOT manual)`, l.Format("2006-01-02")).Scan(&done); err != nil {
		return false, err
	}
	if done {
		return false, nil
	}
	res, err := s.Generate(ctx, "")
	if err == nil {
		s.Log.Info("payout terjadwal dibuat", "jumlah", res.Created, "total", res.Total)
	}
	return err == nil, err
}

// runScheduled: RunIfDue sebagai satu putaran job (1 = payout dibuat hari ini).
func (s *Service) runScheduled(ctx context.Context) (int, error) { // jobs:
	ok, err := s.RunIfDue(ctx, time.Now())
	if ok {
		return 1, err
	}
	return 0, err
}

func (s *Service) Worker(ctx context.Context) {
	t := time.NewTicker(5 * time.Minute)
	defer t.Stop()
	for {
		if n, err := jobs.Track(ctx, s.Pool, "payout.complete", s.CompleteDue); err != nil && ctx.Err() == nil { // jobs:
			s.Log.Error("worker payout (pesanan selesai)", "err", err)
		} else if n > 0 {
			s.Log.Info("pesanan selesai → saldo seller", "jumlah", n)
		}
		if _, err := jobs.Track(ctx, s.Pool, "payout.schedule", s.runScheduled); err != nil && ctx.Err() == nil { // jobs:
			s.Log.Error("worker payout (jadwal)", "err", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}
	}
}

// ---------------------------------------------------------------------------
// Admin keuangan
// ---------------------------------------------------------------------------

// AdminList: status pending (dengan nomor rekening lengkap — dicatat di audit) · paid · cancelled.
func (s *Service) AdminList(ctx context.Context, adminID, status string) ([]PayoutRow, error) {
	if status != "paid" && status != "cancelled" {
		status = "pending"
	}
	list, err := s.payouts(ctx, `p.status = $1`, 300, status == "pending", status)
	if err == nil && status == "pending" && len(list) > 0 {
		_, _ = s.Pool.Exec(ctx, `INSERT INTO audit_logs (actor_type, actor_user_id, action, entity_table, entity_id)
			VALUES ('admin', $1, 'payout_account.view', 'payouts', $2)`, adminID, fmt.Sprint(len(list))+" pending")
	}
	return list, err
}

func lockPayout(ctx context.Context, tx pgx.Tx, id string) (status, storeID string, amount decimal.Decimal, err error) {
	if len(id) != 36 {
		return "", "", amount, ErrNotFound
	}
	err = tx.QueryRow(ctx, `SELECT status, store_id, amount FROM payouts WHERE id = $1 FOR UPDATE`, id).Scan(&status, &storeID, &amount)
	if errors.Is(err, pgx.ErrNoRows) {
		err = ErrNotFound
	}
	return
}

// MarkPaid: admin sudah mentransfer → catat nomor referensi bank; seller diberi tahu.
func (s *Service) MarkPaid(ctx context.Context, adminID, id, reference, note string) error {
	reference, note = strings.TrimSpace(reference), strings.TrimSpace(note)
	if len(reference) < 3 || len(reference) > 80 {
		return fieldError{"reference", "enter the bank transfer reference"}
	}
	if len([]rune(note)) > 300 {
		return fieldError{"note", "too long"}
	}
	return db.WithActor(ctx, s.Pool, adminID, "admin", "", func(tx pgx.Tx) error {
		status, storeID, amount, err := lockPayout(ctx, tx, id)
		if err != nil {
			return err
		}
		if status != "pending" {
			return ErrBadStatus
		}
		if _, err := tx.Exec(ctx, `UPDATE payouts SET status = 'paid', paid_at = now(), paid_by = $2, reference = $3, note = nullif($4, '') WHERE id = $1`,
			id, adminID, reference, note); err != nil {
			return err
		}
		// Rekening terbukti benar → tandai terverifikasi.
		if _, err := tx.Exec(ctx, `UPDATE seller_payout_accounts SET verified_at = coalesce(verified_at, now()) WHERE id = (SELECT account_id FROM payouts WHERE id = $1)`, id); err != nil {
			return err
		}
		return notify.Enqueue(ctx, tx, storeID, notify.PayoutSent, notify.Params{Note: "$" + amount.StringFixed(2), Product: reference}, "#/store/"+storeID+"/payouts")
	})
}

// Cancel: payout dibatalkan (mis. rekening salah) → saldo kembali ke seller.
func (s *Service) Cancel(ctx context.Context, adminID, id, note string) error {
	note = strings.TrimSpace(note)
	if len([]rune(note)) < 3 || len([]rune(note)) > 300 {
		return fieldError{"note", "write the reason"}
	}
	return db.WithActor(ctx, s.Pool, adminID, "admin", "", func(tx pgx.Tx) error {
		status, storeID, amount, err := lockPayout(ctx, tx, id)
		if err != nil {
			return err
		}
		if status != "pending" {
			return ErrBadStatus
		}
		if _, err := tx.Exec(ctx, `UPDATE payouts SET status = 'cancelled', cancelled_at = now(), note = $2 WHERE id = $1`, id, note); err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `INSERT INTO seller_ledger (store_id, kind, amount, payout_id, note, created_by) VALUES ($1, 'payout_cancelled', $2, $3, $4, $5)`,
			storeID, amount, id, note, adminID)
		return err
	})
}

var timeNow = time.Now
