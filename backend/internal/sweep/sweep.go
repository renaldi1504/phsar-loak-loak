// Package sweep: uang Phsar di rekening penampung PayWay dan pemindahan komisi ("sweep") ke rekening perusahaan.
//
//	saldo penampung (seharusnya) = uang masuk (checkout dibayar) − refund ke pembeli − payout seller yang sudah dibayar − sweep
//	ditahan untuk seller         = saldo seller (per toko, minimal 0) + payout yang belum ditransfer
//	ditahan untuk pesanan        = pesanan dibayar yang belum selesai + refund yang belum dikirim + retur disetujui (pesanan selesai)
//	komisi bersih                = komisi pesanan selesai − komisi yang ikut dikembalikan saat retur
//	boleh di-sweep               = min(komisi bersih − komisi retur yang menunggu refund − sweep, saldo penampung − semua yang ditahan)
//
// Hitungan sengaja konservatif: lebih baik uang tertahan sedikit lebih lama daripada rekening penampung kurang
// saat payout atau refund harus dibayar.
package sweep

import (
	"context"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"phsar/internal/platform/db"
)

var (
	ErrNotFound = errors.New("sweep tidak ditemukan")
	ErrVoided   = errors.New("sweep sudah dibatalkan")
)

// FieldError: input tidak valid (422).
type FieldError struct{ Field, Msg string }

func (e FieldError) Error() string { return e.Field + ": " + e.Msg }

type Service struct{ Pool *pgxpool.Pool }

type querier interface {
	QueryRow(context.Context, string, ...any) pgx.Row
	Query(context.Context, string, ...any) (pgx.Rows, error)
}

// Summary: semua angka dalam USD (string 2 desimal).
type Summary struct {
	Collected        string `json:"collected"`          // checkout dibayar
	Refunds          string `json:"refunds"`            // refund yang sudah dikirim ke pembeli
	PayoutsPaid      string `json:"payouts_paid"`       // payout seller yang sudah ditransfer
	Swept            string `json:"swept"`              // sweep yang tidak dibatalkan
	ExpectedBalance  string `json:"expected_balance"`   // saldo rekening penampung seharusnya
	HeldSellers      string `json:"held_sellers"`       // saldo seller + payout belum ditransfer
	PendingPayouts   string `json:"pending_payouts"`    // bagian dari held_sellers
	HeldOrders       string `json:"held_orders"`        // pesanan belum selesai + refund/retur yang belum dibayar
	CommissionGross  string `json:"commission_gross"`   // komisi + ongkir − voucher (pesanan selesai) + biaya iklan
	CommissionBack   string `json:"commission_back"`    // komisi dikembalikan lewat retur yang sudah direfund
	CommissionNet    string `json:"commission_net"`     // gross − back
	CommissionOnHold string `json:"commission_on_hold"` // komisi dari retur disetujui yang menunggu refund
	FreeCash         string `json:"free_cash"`          // saldo − semua yang ditahan
	Available        string `json:"available"`          // boleh di-sweep sekarang
	AsOf             string `json:"as_of"`
}

const summarySQL = `
WITH
  collected AS (SELECT coalesce(sum(coalesce(paid_amount, total)), 0) v FROM checkouts WHERE paid_at IS NOT NULL),
  -- refund pesanan batal/retur penuh; retur yang dicatat di pesanan (nomor referensi sama) dihitung lewat tabel returns
  refund_orders AS (SELECT coalesce(sum(coalesce(o.refund_amount, o.total)), 0) v FROM orders o
    WHERE o.refund_status = 'refunded' AND NOT EXISTS (SELECT 1 FROM returns r WHERE r.order_id = o.id AND r.status = 'refunded'
      AND r.refund_ref IS NOT DISTINCT FROM o.refund_ref)),
  refund_returns AS (SELECT coalesce(sum(refund_amount), 0) v FROM returns WHERE status = 'refunded'),
  payouts_paid AS (SELECT coalesce(sum(amount), 0) v FROM payouts WHERE status = 'paid'),
  payouts_pending AS (SELECT coalesce(sum(amount), 0) v FROM payouts WHERE status = 'pending'),
  swept AS (SELECT coalesce(sum(amount), 0) v FROM commission_sweeps WHERE voided_at IS NULL),
  seller_bal AS (SELECT coalesce(sum(greatest(b, 0)), 0) v FROM (SELECT sum(amount) b FROM seller_ledger GROUP BY store_id) x),
  open_orders AS (SELECT coalesce(sum(total), 0) v FROM orders WHERE status IN ('paid','accepted','ready_for_pickup','picked_up','delivered')),
  refund_due AS (SELECT coalesce(sum(coalesce(refund_amount, total)), 0) v FROM orders WHERE refund_status = 'pending'),
  return_due AS (SELECT coalesce(sum(r.refund_amount), 0) v,
      coalesce(sum(round(r.refund_amount * o.commission_amount / o.total, 2)), 0) comm
    FROM returns r JOIN orders o ON o.id = r.order_id WHERE r.status = 'approved' AND o.status = 'completed'),
  -- pendapatan Phsar dari pesanan selesai: komisi + ongkir − voucher yang ditanggung Phsar; ditambah biaya iklan yang sudah dipotong dari saldo seller
  comm_gross AS (SELECT (SELECT coalesce(sum(commission_amount + shipping_fee - discount - flash_subsidy), 0) FROM orders WHERE status = 'completed')
    + (SELECT coalesce(-sum(amount), 0) FROM seller_ledger WHERE kind = 'ad_spend') v),
  comm_back AS (SELECT coalesce(sum(r.refund_amount + coalesce(l.amount, 0)), 0) v FROM returns r
    JOIN orders o ON o.id = r.order_id AND o.status = 'completed'
    LEFT JOIN seller_ledger l ON l.return_id = r.id AND l.kind = 'return_refund'
    WHERE r.status = 'refunded')
SELECT collected.v, refund_orders.v + refund_returns.v, payouts_paid.v, payouts_pending.v, swept.v, seller_bal.v,
  open_orders.v + refund_due.v + return_due.v, return_due.comm, comm_gross.v, greatest(comm_back.v, 0)
FROM collected, refund_orders, refund_returns, payouts_paid, payouts_pending, swept, seller_bal, open_orders, refund_due, return_due, comm_gross, comm_back`

func (s *Service) summary(ctx context.Context, q querier) (*Summary, decimal.Decimal, error) {
	var collected, refunds, paid, pending, swept, sellers, orders, commHold, gross, back decimal.Decimal
	if err := q.QueryRow(ctx, summarySQL).Scan(&collected, &refunds, &paid, &pending, &swept, &sellers, &orders, &commHold, &gross, &back); err != nil {
		return nil, decimal.Zero, err
	}
	heldSellers := sellers.Add(pending)
	expected := collected.Sub(refunds).Sub(paid).Sub(swept)
	free := expected.Sub(heldSellers).Sub(orders)
	net := gross.Sub(back)
	avail := decimal.Min(net.Sub(commHold).Sub(swept), free)
	if avail.IsNegative() {
		avail = decimal.Zero
	}
	f := func(d decimal.Decimal) string { return d.StringFixed(2) }
	return &Summary{Collected: f(collected), Refunds: f(refunds), PayoutsPaid: f(paid), Swept: f(swept), ExpectedBalance: f(expected),
		HeldSellers: f(heldSellers), PendingPayouts: f(pending), HeldOrders: f(orders), CommissionGross: f(gross), CommissionBack: f(back),
		CommissionNet: f(net), CommissionOnHold: f(commHold), FreeCash: f(free), Available: f(avail), AsOf: time.Now().UTC().Format(time.RFC3339)}, avail, nil
}

func (s *Service) Summary(ctx context.Context) (*Summary, error) {
	sum, _, err := s.summary(ctx, s.Pool)
	return sum, err
}

type Sweep struct {
	ID        string     `json:"id"`
	Amount    string     `json:"amount"`
	ToAccount string     `json:"to_account"`
	Reference string     `json:"reference"`
	SweptOn   string     `json:"swept_on"`
	Note      *string    `json:"note"`
	By        string     `json:"created_by"`
	CreatedAt time.Time  `json:"created_at"`
	VoidedAt  *time.Time `json:"voided_at"`
	VoidedBy  *string    `json:"voided_by"`
	VoidNote  *string    `json:"void_note"`
}

func (s *Service) List(ctx context.Context, limit int) ([]Sweep, error) {
	return s.list(ctx, "TRUE", limit)
}

func (s *Service) list(ctx context.Context, where string, limit int, args ...any) ([]Sweep, error) {
	rows, err := s.Pool.Query(ctx, `SELECT s.id, s.amount::text, s.to_account, s.reference, s.swept_on::text, s.note,
		coalesce(nullif(u.name, ''), u.phone, ''), s.created_at, s.voided_at, coalesce(nullif(v.name, ''), v.phone), s.void_note
		FROM commission_sweeps s JOIN users u ON u.id = s.created_by LEFT JOIN users v ON v.id = s.voided_by
		WHERE `+where+` ORDER BY s.swept_on DESC, s.created_at DESC LIMIT `+strconv.Itoa(limit), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Sweep{}
	for rows.Next() {
		var x Sweep
		if err := rows.Scan(&x.ID, &x.Amount, &x.ToAccount, &x.Reference, &x.SweptOn, &x.Note, &x.By, &x.CreatedAt, &x.VoidedAt, &x.VoidedBy, &x.VoidNote); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}

// Month: ringkasan per bulan (zona waktu Phnom Penh).
type Month struct {
	Month      string `json:"month"` // 2026-11
	Collected  string `json:"collected"`
	Refunds    string `json:"refunds"`
	Payouts    string `json:"payouts"`
	Commission string `json:"commission"` // komisi pesanan yang selesai bulan itu − komisi retur yang direfund bulan itu
	Swept      string `json:"swept"`
}

func (s *Service) Monthly(ctx context.Context, months int) ([]Month, error) {
	rows, err := s.Pool.Query(ctx, `
		WITH m AS (SELECT generate_series(date_trunc('month', now() AT TIME ZONE 'Asia/Phnom_Penh') - make_interval(months => $1 - 1),
			date_trunc('month', now() AT TIME ZONE 'Asia/Phnom_Penh'), interval '1 month') AS start),
		r AS (SELECT start, start + interval '1 month' AS stop FROM m)
		SELECT to_char(r.start, 'YYYY-MM'),
		  (SELECT coalesce(sum(coalesce(paid_amount, total)), 0) FROM checkouts WHERE paid_at IS NOT NULL
		    AND paid_at AT TIME ZONE 'Asia/Phnom_Penh' >= r.start AND paid_at AT TIME ZONE 'Asia/Phnom_Penh' < r.stop)::text,
		  ((SELECT coalesce(sum(coalesce(o.refund_amount, o.total)), 0) FROM orders o WHERE o.refund_status = 'refunded'
		    AND o.refunded_at AT TIME ZONE 'Asia/Phnom_Penh' >= r.start AND o.refunded_at AT TIME ZONE 'Asia/Phnom_Penh' < r.stop
		    AND NOT EXISTS (SELECT 1 FROM returns x WHERE x.order_id = o.id AND x.status = 'refunded' AND x.refund_ref IS NOT DISTINCT FROM o.refund_ref))
		   + (SELECT coalesce(sum(refund_amount), 0) FROM returns WHERE status = 'refunded'
		    AND refunded_at AT TIME ZONE 'Asia/Phnom_Penh' >= r.start AND refunded_at AT TIME ZONE 'Asia/Phnom_Penh' < r.stop))::text,
		  (SELECT coalesce(sum(amount), 0) FROM payouts WHERE status = 'paid'
		    AND paid_at AT TIME ZONE 'Asia/Phnom_Penh' >= r.start AND paid_at AT TIME ZONE 'Asia/Phnom_Penh' < r.stop)::text,
		  ((SELECT coalesce(sum(commission_amount + shipping_fee - discount - flash_subsidy), 0) FROM orders WHERE status = 'completed'
		    AND completed_at AT TIME ZONE 'Asia/Phnom_Penh' >= r.start AND completed_at AT TIME ZONE 'Asia/Phnom_Penh' < r.stop)
		   + (SELECT coalesce(-sum(amount), 0) FROM seller_ledger WHERE kind = 'ad_spend'
		    AND created_at AT TIME ZONE 'Asia/Phnom_Penh' >= r.start AND created_at AT TIME ZONE 'Asia/Phnom_Penh' < r.stop)
		   - (SELECT coalesce(sum(x.refund_amount + coalesce(l.amount, 0)), 0) FROM returns x JOIN orders o ON o.id = x.order_id AND o.status = 'completed'
		    LEFT JOIN seller_ledger l ON l.return_id = x.id AND l.kind = 'return_refund'
		    WHERE x.status = 'refunded' AND x.refunded_at AT TIME ZONE 'Asia/Phnom_Penh' >= r.start AND x.refunded_at AT TIME ZONE 'Asia/Phnom_Penh' < r.stop))::text,
		  (SELECT coalesce(sum(amount), 0) FROM commission_sweeps WHERE voided_at IS NULL AND swept_on >= r.start::date AND swept_on < r.stop::date)::text
		FROM r ORDER BY r.start DESC`, months)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Month{}
	for rows.Next() {
		var m Month
		if err := rows.Scan(&m.Month, &m.Collected, &m.Refunds, &m.Payouts, &m.Commission, &m.Swept); err != nil {
			return nil, err
		}
		for _, p := range []*string{&m.Collected, &m.Refunds, &m.Payouts, &m.Commission, &m.Swept} {
			d, _ := decimal.NewFromString(*p)
			*p = d.StringFixed(2)
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

type Input struct {
	Amount    string `json:"amount"`
	ToAccount string `json:"to_account"`
	Reference string `json:"reference"`
	SweptOn   string `json:"swept_on"` // YYYY-MM-DD
	Note      string `json:"note"`
}

var amountRe = regexp.MustCompile(`^\d{1,9}(\.\d{1,2})?$`)

var phnomPenh = func() *time.Location {
	if l, err := time.LoadLocation("Asia/Phnom_Penh"); err == nil {
		return l
	}
	return time.FixedZone("ICT", 7*3600)
}()

func runes(s string) int { return utf8.RuneCountInString(s) }

// Record mencatat sweep. Dikunci (advisory lock) supaya dua admin tidak bisa melebihi saldo bersamaan.
func (s *Service) Record(ctx context.Context, actorID string, in Input) (*Sweep, error) {
	in.Amount, in.ToAccount, in.Reference, in.Note = strings.TrimSpace(in.Amount), strings.TrimSpace(in.ToAccount), strings.TrimSpace(in.Reference), strings.TrimSpace(in.Note)
	amt, err := decimal.NewFromString(in.Amount)
	if !amountRe.MatchString(in.Amount) || err != nil || !amt.IsPositive() {
		return nil, FieldError{"amount", "enter an amount like 120.50"}
	}
	if n := runes(in.ToAccount); n < 3 || n > 120 {
		return nil, FieldError{"to_account", "which company account received the money (3–120 characters)"}
	}
	if n := runes(in.Reference); n < 3 || n > 80 {
		return nil, FieldError{"reference", "enter the bank transfer reference"}
	}
	if runes(in.Note) > 300 {
		return nil, FieldError{"note", "at most 300 characters"}
	}
	day, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(in.SweptOn), phnomPenh)
	today := time.Now().In(phnomPenh)
	if err != nil || day.After(today) || day.Before(today.AddDate(-1, 0, 0)) {
		return nil, FieldError{"swept_on", "transfer date: today or within the last year"}
	}
	var id string
	err = db.WithActor(ctx, s.Pool, actorID, "admin", "", func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtext('phsar.commission_sweep'))`); err != nil {
			return err
		}
		_, avail, err := s.summary(ctx, tx)
		if err != nil {
			return err
		}
		if amt.GreaterThan(avail) {
			return FieldError{"amount", "more than available to sweep ($" + avail.StringFixed(2) + ")"}
		}
		return tx.QueryRow(ctx, `INSERT INTO commission_sweeps (amount, to_account, reference, swept_on, note, created_by)
			VALUES ($1, $2, $3, $4, nullif($5, ''), $6) RETURNING id`, amt, in.ToAccount, in.Reference, day.Format("2006-01-02"), in.Note, actorID).Scan(&id)
	})
	if err != nil {
		return nil, err
	}
	return s.get(ctx, id)
}

// Void membatalkan sweep yang salah catat (uangnya kembali dihitung ada di rekening penampung).
func (s *Service) Void(ctx context.Context, actorID, id, note string) (*Sweep, error) {
	note = strings.TrimSpace(note)
	if n := runes(note); n < 5 || n > 300 {
		return nil, FieldError{"note", "why is this sweep cancelled (5–300 characters)"}
	}
	err := db.WithActor(ctx, s.Pool, actorID, "admin", "", func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, `UPDATE commission_sweeps SET voided_at = now(), voided_by = $2, void_note = $3 WHERE id::text = $1 AND voided_at IS NULL`, id, actorID, note)
		if err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			var voided bool
			if err := tx.QueryRow(ctx, `SELECT voided_at IS NOT NULL FROM commission_sweeps WHERE id::text = $1`, id).Scan(&voided); err != nil {
				return ErrNotFound
			}
			return ErrVoided
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return s.get(ctx, id)
}

func (s *Service) get(ctx context.Context, id string) (*Sweep, error) {
	l, err := s.list(ctx, "s.id::text = $1", 1, id)
	if err != nil {
		return nil, err
	}
	if len(l) == 0 {
		return nil, ErrNotFound
	}
	return &l[0], nil
}
