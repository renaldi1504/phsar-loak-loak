// Package invoices: faktur komisi bulanan Phsar untuk tiap toko + pengaturan pajak (PPN/VAT).
//
//	Pesanan "completed" di bulan itu (komisi yang dipotong dari seller) − komisi yang dikembalikan saat retur di-refund
//	→ satu faktur per toko per bulan, nomor PREFIX-YYYYMM-####. Isinya disalin saat terbit (snapshot).
//	Void: nomor tetap tercatat, lalu boleh diterbitkan ulang dengan nomor baru (replaces_id).
package invoices

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"phsar/internal/platform/db"
)

var (
	ErrNotFound = errors.New("faktur tidak ditemukan")
	ErrState    = errors.New("faktur sudah dibatalkan")
)

// FieldError: input tidak valid (422).
type FieldError struct{ Field, Msg string }

func (e FieldError) Error() string { return e.Field + ": " + e.Msg }

type Service struct{ Pool *pgxpool.Pool }

var phnomPenh = func() *time.Location {
	if l, err := time.LoadLocation("Asia/Phnom_Penh"); err == nil {
		return l
	}
	return time.FixedZone("ICT", 7*3600)
}()

// ---------------------------------------------------------------------------
// Pengaturan pajak (platform_settings 'tax')
// ---------------------------------------------------------------------------

type Settings struct {
	LegalName    string `json:"legal_name"`
	VATTIN       string `json:"vat_tin"`
	Address      string `json:"address"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	VATRate      string `json:"vat_rate"`      // persen, mis. "10"
	VATInclusive bool   `json:"vat_inclusive"` // komisi sudah termasuk PPN
	Prefix       string `json:"prefix"`        // awalan nomor faktur
}

func defaults() Settings {
	return Settings{LegalName: "Phsar", Address: "Sihanoukville, Cambodia", VATRate: "10", VATInclusive: true, Prefix: "PHS"}
}

func (s *Service) Settings(ctx context.Context) (Settings, error) {
	st := defaults()
	var raw []byte
	err := s.Pool.QueryRow(ctx, `SELECT value FROM platform_settings WHERE key = 'tax'`).Scan(&raw)
	if errors.Is(err, pgx.ErrNoRows) {
		return st, nil
	}
	if err != nil {
		return st, err
	}
	_ = json.Unmarshal(raw, &st)
	return st, nil
}

var prefixRe = regexp.MustCompile(`^[A-Z0-9]{2,8}$`)

func trim(v string, n int) string {
	r := []rune(strings.TrimSpace(v))
	if len(r) > n {
		r = r[:n]
	}
	return string(r)
}

func (s *Service) SaveSettings(ctx context.Context, actorID string, in Settings) (Settings, error) {
	in = Settings{LegalName: trim(in.LegalName, 160), VATTIN: trim(in.VATTIN, 40), Address: trim(in.Address, 300), Email: trim(in.Email, 120),
		Phone: trim(in.Phone, 40), VATRate: strings.TrimSpace(in.VATRate), VATInclusive: in.VATInclusive, Prefix: strings.ToUpper(strings.TrimSpace(in.Prefix))}
	if len([]rune(in.LegalName)) < 2 {
		return in, FieldError{"legal_name", "required"}
	}
	rate, err := decimal.NewFromString(in.VATRate)
	if err != nil || rate.IsNegative() || rate.GreaterThan(decimal.NewFromInt(30)) || rate.Exponent() < -2 {
		return in, FieldError{"vat_rate", "0–30 (percent)"}
	}
	in.VATRate = rate.String()
	if !prefixRe.MatchString(in.Prefix) {
		return in, FieldError{"prefix", "2–8 letters or digits"}
	}
	if in.Email != "" && !strings.Contains(in.Email, "@") {
		return in, FieldError{"email", "invalid email"}
	}
	v, _ := json.Marshal(in)
	err = db.WithActor(ctx, s.Pool, actorID, "admin", "", func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `INSERT INTO platform_settings (key, value, updated_by, updated_at) VALUES ('tax', $1, $2, now())
			ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_by = EXCLUDED.updated_by, updated_at = now()`, v, actorID)
		return err
	})
	return in, err
}

// ---------------------------------------------------------------------------
// Model
// ---------------------------------------------------------------------------

type Line struct {
	Kind       string    `json:"kind"` // order · return
	Ref        string    `json:"ref"`  // nomor pesanan / retur
	Order      string    `json:"order,omitempty"`
	Date       time.Time `json:"date"`
	Subtotal   string    `json:"subtotal"`   // nilai penjualan (negatif untuk retur)
	Commission string    `json:"commission"` // komisi (negatif = dikembalikan)
}

type Party struct {
	Name       string `json:"name"`
	LegalName  string `json:"legal_name"`
	TIN        string `json:"tin,omitempty"`
	VATTIN     string `json:"vat_tin,omitempty"`
	Address    string `json:"address,omitempty"`
	Email      string `json:"email,omitempty"`
	Phone      string `json:"phone,omitempty"`
	SellerType string `json:"seller_type,omitempty"`
}

type Invoice struct {
	ID           string     `json:"id"`
	Number       string     `json:"number"`
	StoreID      string     `json:"store_id"`
	StoreName    string     `json:"store_name"`
	Period       string     `json:"period"` // YYYY-MM
	Status       string     `json:"status"`
	Seller       Party      `json:"seller"`
	Issuer       Party      `json:"issuer"`
	Lines        []Line     `json:"lines,omitempty"`
	OrderCount   int        `json:"order_count"`
	SalesTotal   string     `json:"sales_total"`
	Commission   string     `json:"commission"`
	VATRate      string     `json:"vat_rate"`
	VATInclusive bool       `json:"vat_inclusive"`
	Subtotal     string     `json:"subtotal"`
	VATAmount    string     `json:"vat_amount"`
	Total        string     `json:"total"`
	ReplacesNo   *string    `json:"replaces_number"`
	IssuedAt     time.Time  `json:"issued_at"`
	VoidedAt     *time.Time `json:"voided_at"`
	VoidNote     *string    `json:"void_note"`
}

const selectInvoice = `SELECT i.id, i.number, i.store_id, st.name, to_char(i.period, 'YYYY-MM'), i.status, i.seller, i.issuer, %s, i.order_count,
	i.sales_total::text, i.commission::text, i.vat_rate::text, i.vat_inclusive, i.subtotal::text, i.vat_amount::text, i.total::text,
	(SELECT p.number FROM invoices p WHERE p.id = i.replaces_id), i.issued_at, i.voided_at, i.void_note
	FROM invoices i JOIN stores st ON st.id = i.store_id`

func (s *Service) query(ctx context.Context, withLines bool, where string, limit int, args ...any) ([]Invoice, error) {
	lines := "'[]'::jsonb"
	if withLines {
		lines = "i.lines"
	}
	rows, err := s.Pool.Query(ctx, fmt.Sprintf(selectInvoice, lines)+" WHERE "+where+" ORDER BY i.period DESC, i.number DESC LIMIT "+fmt.Sprint(limit), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Invoice{}
	for rows.Next() {
		var v Invoice
		var seller, issuer, ls []byte
		if err := rows.Scan(&v.ID, &v.Number, &v.StoreID, &v.StoreName, &v.Period, &v.Status, &seller, &issuer, &ls, &v.OrderCount,
			&v.SalesTotal, &v.Commission, &v.VATRate, &v.VATInclusive, &v.Subtotal, &v.VATAmount, &v.Total, &v.ReplacesNo,
			&v.IssuedAt, &v.VoidedAt, &v.VoidNote); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(seller, &v.Seller)
		_ = json.Unmarshal(issuer, &v.Issuer)
		_ = json.Unmarshal(ls, &v.Lines)
		out = append(out, v)
	}
	return out, rows.Err()
}

func (s *Service) Get(ctx context.Context, id string) (*Invoice, error) {
	return s.one(ctx, "i.id = $1", id)
}

// StoreGet: faktur milik toko ini saja (Seller Center).
func (s *Service) StoreGet(ctx context.Context, storeID, id string) (*Invoice, error) {
	return s.one(ctx, "i.id = $1 AND i.store_id = $2", id, storeID)
}

func (s *Service) one(ctx context.Context, where string, args ...any) (*Invoice, error) {
	if len(args) == 0 || !uuidRe.MatchString(fmt.Sprint(args[0])) {
		return nil, ErrNotFound
	}
	l, err := s.query(ctx, true, where, 1, args...)
	if err != nil {
		return nil, err
	}
	if len(l) == 0 {
		return nil, ErrNotFound
	}
	return &l[0], nil
}

var uuidRe = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

// List untuk admin: period (YYYY-MM), status (issued|void), q (nomor / nama toko).
func (s *Service) List(ctx context.Context, period, status, q string) ([]Invoice, error) {
	where, args := []string{"TRUE"}, []any{}
	if period != "" {
		p, err := parsePeriod(period)
		if err != nil {
			return nil, err
		}
		args = append(args, p)
		where = append(where, fmt.Sprintf("i.period = $%d", len(args)))
	}
	switch status {
	case "":
	case "issued", "void":
		args = append(args, status)
		where = append(where, fmt.Sprintf("i.status = $%d", len(args)))
	default:
		return nil, FieldError{"status", "issued or void"}
	}
	if q = strings.TrimSpace(q); q != "" {
		args = append(args, q)
		where = append(where, fmt.Sprintf("(i.number ILIKE '%%' || $%d || '%%' OR st.name ILIKE '%%' || $%d || '%%')", len(args), len(args)))
	}
	return s.query(ctx, false, strings.Join(where, " AND "), 300, args...)
}

func (s *Service) StoreList(ctx context.Context, storeID string) ([]Invoice, error) {
	return s.query(ctx, false, "i.store_id = $1", 120, storeID)
}

// ---------------------------------------------------------------------------
// Terbitkan
// ---------------------------------------------------------------------------

var periodRe = regexp.MustCompile(`^(\d{4})-(\d{2})$`)

func parsePeriod(p string) (time.Time, error) {
	m := periodRe.FindStringSubmatch(strings.TrimSpace(p))
	if m == nil {
		return time.Time{}, FieldError{"period", "use YYYY-MM"}
	}
	t, err := time.Parse("2006-01", m[0])
	if err != nil || t.Year() < 2024 {
		return time.Time{}, FieldError{"period", "use YYYY-MM"}
	}
	return t, nil
}

type GenerateResult struct {
	Period  string    `json:"period"`
	Created []Invoice `json:"created"`
	Skipped int       `json:"skipped"` // toko yang sudah punya faktur bulan ini
}

// Generate menerbitkan faktur bulan `period` untuk semua toko yang punya transaksi. Aman dijalankan berulang.
// Hanya bulan yang sudah berakhir (waktu Phnom Penh), supaya faktur tidak berubah lagi.
func (s *Service) Generate(ctx context.Context, actorID, period string, now time.Time) (*GenerateResult, error) {
	p, err := parsePeriod(period)
	if err != nil {
		return nil, err
	}
	n := now.In(phnomPenh)
	if !p.Before(time.Date(n.Year(), n.Month(), 1, 0, 0, 0, 0, time.UTC)) {
		return nil, FieldError{"period", "only months that have ended"}
	}
	st, err := s.Settings(ctx)
	if err != nil {
		return nil, err
	}
	rate, _ := decimal.NewFromString(st.VATRate)
	issuer := Party{Name: "Phsar", LegalName: st.LegalName, VATTIN: st.VATTIN, Address: st.Address, Email: st.Email, Phone: st.Phone}
	issuerJSON, _ := json.Marshal(issuer)
	res := &GenerateResult{Period: p.Format("2006-01"), Created: []Invoice{}}
	var ids []string
	err = db.WithActor(ctx, s.Pool, actorID, "admin", "", func(tx pgx.Tx) error {
		// Kunci per bulan: dua admin menekan "Terbitkan" bersamaan tidak membuat faktur ganda.
		if _, err := tx.Exec(ctx, `INSERT INTO invoice_counters (period) VALUES ($1) ON CONFLICT DO NOTHING`, p); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `SELECT 1 FROM invoice_counters WHERE period = $1 FOR UPDATE`, p); err != nil {
			return err
		}
		rng := `>= ($1::date)::timestamp AT TIME ZONE 'Asia/Phnom_Penh' AND %s < (($1::date + interval '1 month')::timestamp AT TIME ZONE 'Asia/Phnom_Penh')`
		inOrder := "o.completed_at " + fmt.Sprintf(rng, "o.completed_at")
		inReturn := "rt.refunded_at " + fmt.Sprintf(rng, "rt.refunded_at")
		rows, err := tx.Query(ctx, `SELECT s.id, EXISTS (SELECT 1 FROM invoices i WHERE i.store_id = s.id AND i.period = $1 AND i.status = 'issued')
			FROM stores s WHERE EXISTS (SELECT 1 FROM orders o WHERE o.store_id = s.id AND `+inOrder+`)
			   OR EXISTS (SELECT 1 FROM returns rt JOIN orders o ON o.id = rt.order_id WHERE rt.store_id = s.id AND rt.status = 'refunded'
			              AND o.completed_at IS NOT NULL AND `+inReturn+`) ORDER BY s.name`, p)
		if err != nil {
			return err
		}
		var stores []string
		for rows.Next() {
			var id string
			var has bool
			if err := rows.Scan(&id, &has); err != nil {
				rows.Close()
				return err
			}
			if has {
				res.Skipped++
			} else {
				stores = append(stores, id)
			}
		}
		rows.Close()
		for _, storeID := range stores {
			lines, err := storeLines(ctx, tx, storeID, p, inOrder, inReturn)
			if err != nil {
				return err
			}
			if len(lines) == 0 {
				continue
			}
			id, err := s.issue(ctx, tx, actorID, storeID, p, lines, st, rate, issuerJSON)
			if err != nil {
				return err
			}
			ids = append(ids, id)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	for _, id := range ids {
		if v, err := s.Get(ctx, id); err == nil {
			v.Lines = nil
			res.Created = append(res.Created, *v)
		}
	}
	return res, nil
}

func storeLines(ctx context.Context, tx pgx.Tx, storeID string, p time.Time, inOrder, inReturn string) ([]Line, error) {
	rows, err := tx.Query(ctx, `
		SELECT 'order', o.order_no, '', o.completed_at, o.subtotal, o.commission_amount FROM orders o
		WHERE o.store_id = $2 AND `+inOrder+`
		UNION ALL
		SELECT 'return', rt.return_no, o.order_no, rt.refunded_at, -rt.refund_amount, -round(rt.refund_amount * o.commission_amount / o.total, 2)
		FROM returns rt JOIN orders o ON o.id = rt.order_id
		WHERE rt.store_id = $2 AND rt.status = 'refunded' AND o.completed_at IS NOT NULL AND `+inReturn+`
		ORDER BY 4, 2`, p, storeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Line
	for rows.Next() {
		var l Line
		var sub, comm decimal.Decimal
		if err := rows.Scan(&l.Kind, &l.Ref, &l.Order, &l.Date, &sub, &comm); err != nil {
			return nil, err
		}
		l.Subtotal, l.Commission = sub.StringFixed(2), comm.StringFixed(2)
		out = append(out, l)
	}
	return out, rows.Err()
}

// amounts menghitung PPN dari komisi bersih.
func amounts(commission, rate decimal.Decimal, inclusive bool) (sub, vat, total decimal.Decimal) {
	hundred := decimal.NewFromInt(100)
	if inclusive {
		total = commission
		sub = commission.Mul(hundred).Div(hundred.Add(rate)).Round(2)
		return sub, total.Sub(sub), total
	}
	sub = commission
	vat = commission.Mul(rate).Div(hundred).Round(2)
	return sub, vat, sub.Add(vat)
}

func (s *Service) issue(ctx context.Context, tx pgx.Tx, actorID, storeID string, p time.Time, lines []Line, st Settings, rate decimal.Decimal, issuerJSON []byte) (string, error) {
	var seller Party
	var legal, tin, vattin, stype, addr, phone *string
	if err := tx.QueryRow(ctx, `SELECT st.name, tp.legal_name, tp.tin, tp.vattin, tp.seller_type,
		(SELECT pp.address_line FROM store_pickup_points pp WHERE pp.store_id = st.id AND pp.deleted_at IS NULL ORDER BY pp.is_default DESC, pp.created_at LIMIT 1),
		(SELECT u.phone FROM store_members sm JOIN users u ON u.id = sm.user_id WHERE sm.store_id = st.id AND sm.role = 'owner' LIMIT 1)
		FROM stores st LEFT JOIN seller_tax_profiles tp ON tp.store_id = st.id WHERE st.id = $1`, storeID).
		Scan(&seller.Name, &legal, &tin, &vattin, &stype, &addr, &phone); err != nil {
		return "", err
	}
	deref := func(p *string) string {
		if p == nil {
			return ""
		}
		return *p
	}
	seller.LegalName, seller.TIN, seller.VATTIN, seller.SellerType, seller.Address, seller.Phone = deref(legal), deref(tin), deref(vattin), deref(stype), deref(addr), deref(phone)
	if seller.LegalName == "" {
		seller.LegalName = seller.Name
	}
	orders, sales, comm := 0, decimal.Zero, decimal.Zero
	for _, l := range lines {
		if l.Kind == "order" {
			orders++
		}
		a, _ := decimal.NewFromString(l.Subtotal)
		c, _ := decimal.NewFromString(l.Commission)
		sales, comm = sales.Add(a), comm.Add(c)
	}
	sub, vat, total := amounts(comm, rate, st.VATInclusive)
	var seq int
	if err := tx.QueryRow(ctx, `UPDATE invoice_counters SET last = last + 1 WHERE period = $1 RETURNING last`, p).Scan(&seq); err != nil {
		return "", err
	}
	number := fmt.Sprintf("%s-%s-%04d", st.Prefix, p.Format("200601"), seq)
	sellerJSON, _ := json.Marshal(seller)
	linesJSON, _ := json.Marshal(lines)
	var id string
	err := tx.QueryRow(ctx, `INSERT INTO invoices (number, store_id, period, seller, issuer, lines, order_count, sales_total, commission, vat_rate,
		vat_inclusive, subtotal, vat_amount, total, replaces_id, issued_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
		  (SELECT v.id FROM invoices v WHERE v.store_id = $2 AND v.period = $3 AND v.status = 'void' ORDER BY v.voided_at DESC LIMIT 1), nullif($15, '')::uuid)
		RETURNING id`, number, storeID, p, sellerJSON, issuerJSON, linesJSON, orders, sales, comm, rate, st.VATInclusive, sub, vat, total, actorID).Scan(&id)
	return id, err
}

// Void membatalkan faktur (nomor tetap tercatat). Setelah itu "Terbitkan" bulan yang sama membuat faktur pengganti.
func (s *Service) Void(ctx context.Context, actorID, id, note string) (*Invoice, error) {
	note = strings.TrimSpace(note)
	if n := len([]rune(note)); n < 5 || n > 300 {
		return nil, FieldError{"note", "explain why (5–300 characters)"}
	}
	if !uuidRe.MatchString(id) {
		return nil, ErrNotFound
	}
	err := db.WithActor(ctx, s.Pool, actorID, "admin", "", func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, `UPDATE invoices SET status = 'void', voided_by = $2, voided_at = now(), void_note = $3 WHERE id = $1 AND status = 'issued'`,
			id, actorID, note)
		if err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			var ok bool
			_ = tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM invoices WHERE id = $1)`, id).Scan(&ok)
			if ok {
				return ErrState
			}
			return ErrNotFound
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return s.Get(ctx, id)
}
