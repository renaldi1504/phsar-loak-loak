// Package reconciliation: mencocokkan ekspor transaksi ABA PayWay dengan pembayaran & refund yang tercatat di Phsar.
//
// Hasil per baris file:
//
//	ok                  dibayar di PayWay dan lunas di Phsar dengan jumlah sama
//	amount_mismatch     jumlah berbeda
//	not_paid_in_phsar   dibayar di PayWay tapi checkout di Phsar belum lunas (callback tidak terproses?)
//	unknown             ID transaksi tidak ada di Phsar
//	duplicate           ID transaksi muncul lebih dari sekali di file
//	ignored             transaksi gagal/ditolak/menunggu — bukan uang masuk
//	refund_ok / refund_mismatch / refund_unrecorded   baris refund dibandingkan refund yang dicatat admin
//
// Ditambah dari sisi Phsar: missing_in_statement = lunas di Phsar (dalam periode file) tapi tidak ada di file.
package reconciliation

import (
	"context"
	"encoding/json"
	"errors"
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
	ErrNotFound = errors.New("tidak ditemukan")
	ErrNotOpen  = errors.New("baris ini bukan selisih yang perlu diselesaikan")
)

type FieldError struct{ Field, Msg string }

func (e FieldError) Error() string { return e.Field + ": " + e.Msg }

// Exceptions: hasil yang perlu diperiksa keuangan.
var Exceptions = []string{"amount_mismatch", "not_paid_in_phsar", "unknown", "duplicate", "missing_in_statement", "refund_mismatch", "refund_unrecorded"}

func isException(r string) bool {
	for _, x := range Exceptions {
		if x == r {
			return true
		}
	}
	return false
}

type Service struct{ Pool *pgxpool.Pool }

var failWords = []string{"declin", "fail", "cancel", "error", "pending", "void", "expire", "reject", "timeout", "unpaid"}

// success: status kosong / APPROVED / SUCCESS / REFUNDED … = uang benar-benar berpindah.
func success(status string) bool {
	s := strings.ToLower(status)
	for _, w := range failWords {
		if strings.Contains(s, w) {
			return false
		}
	}
	return true
}

func parseDate(s string) (*time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	t, err := time.ParseInLocation("2006-01-02", s, Zone)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// Import menyimpan file sebagai batch baru lalu langsung mencocokkan.
func (s *Service) Import(ctx context.Context, actorID, fileName string, data []byte, from, to string) (string, error) {
	fileName = cut(strings.TrimSpace(fileName), 200)
	if !strings.HasSuffix(strings.ToLower(fileName), ".csv") && !strings.HasSuffix(strings.ToLower(fileName), ".txt") {
		return "", FieldError{"file", "upload the PayWay export as a .csv file"}
	}
	lines, err := Parse(data)
	if err != nil {
		return "", FieldError{"file", err.Error()}
	}
	if len(lines) == 0 {
		return "", FieldError{"file", "no transactions found in the file"}
	}
	pf, err1 := parseDate(from)
	pt, err2 := parseDate(to)
	if err1 != nil {
		return "", FieldError{"period_from", "use YYYY-MM-DD"}
	}
	if err2 != nil {
		return "", FieldError{"period_to", "use YYYY-MM-DD"}
	}
	collected, refunded, fees := decimal.Zero, decimal.Zero, decimal.Zero
	var minAt, maxAt *time.Time
	for _, l := range lines {
		if l.At != nil {
			if minAt == nil || l.At.Before(*minAt) {
				minAt = l.At
			}
			if maxAt == nil || l.At.After(*maxAt) {
				maxAt = l.At
			}
		}
		if l.Fee != nil {
			fees = fees.Add(*l.Fee)
		}
		if !success(l.Status) {
			continue
		}
		if l.Kind == "refund" {
			refunded = refunded.Add(l.Amount)
		} else {
			collected = collected.Add(l.Amount)
		}
	}
	if pf == nil {
		pf = minAt
	}
	if pt == nil {
		pt = maxAt
	}
	if pf != nil && pt != nil && pt.Before(*pf) {
		return "", FieldError{"period_to", "must be after the start date"}
	}
	day := func(t *time.Time) any {
		if t == nil {
			return nil
		}
		return t.In(Zone).Format("2006-01-02")
	}
	var id string
	err = db.WithActor(ctx, s.Pool, actorID, "admin", "", func(tx pgx.Tx) error {
		if err := tx.QueryRow(ctx, `INSERT INTO recon_batches (file_name, period_from, period_to, row_count, collected, refunded, fees, uploaded_by)
			VALUES ($1, $2::date, $3::date, $4, $5, $6, $7, $8) RETURNING id`, fileName, day(pf), day(pt), len(lines), collected, refunded, fees, actorID).Scan(&id); err != nil {
			return err
		}
		rows := make([][]any, 0, len(lines))
		for _, l := range lines {
			raw, _ := json.Marshal(l.Raw)
			var amt, fee any
			if l.HasAmount {
				amt = l.Amount
			}
			if l.Fee != nil {
				fee = *l.Fee
			}
			var at any
			if l.At != nil {
				at = *l.At
			}
			rows = append(rows, []any{id, "statement", l.LineNo, l.Kind, l.TranID, nilIfEmpty(l.Reference), nilIfEmpty(l.Status), amt,
				nilIfEmpty(l.Currency), fee, at, string(raw)})
		}
		_, err := tx.CopyFrom(ctx, pgx.Identifier{"recon_rows"}, []string{"batch_id", "source", "line_no", "kind", "tran_id", "reference",
			"status", "amount", "currency", "fee", "occurred_at", "raw"}, pgx.CopyFromRows(rows))
		if err != nil {
			return err
		}
		return match(ctx, tx, id)
	})
	return id, err
}

func nilIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}

// Rematch: cocokkan ulang (mis. setelah callback yang tertunda diproses).
func (s *Service) Rematch(ctx context.Context, actorID, batchID string) error {
	return db.WithActor(ctx, s.Pool, actorID, "admin", "", func(tx pgx.Tx) error {
		var ok bool
		if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM recon_batches WHERE id = $1)`, batchID).Scan(&ok); err != nil {
			return err
		}
		if !ok {
			return ErrNotFound
		}
		return match(ctx, tx, batchID)
	})
}

type stmtRow struct {
	id                      int64
	kind, tran, ref, status string
	amount                  *decimal.Decimal
	oldResult               string
	result                  string
	checkoutID              *string
	expected                *decimal.Decimal
}

func match(ctx context.Context, tx pgx.Tx, batchID string) error {
	if _, err := tx.Exec(ctx, `SELECT 1 FROM recon_batches WHERE id = $1 FOR UPDATE`, batchID); err != nil {
		return err
	}
	// penyelesaian baris sisi Phsar disimpan dulu supaya tidak hilang saat dibuat ulang
	type resolution struct {
		by   *string
		at   *time.Time
		note *string
	}
	saved := map[string]resolution{}
	rs, err := tx.Query(ctx, `SELECT tran_id, resolved_by::text, resolved_at, resolution FROM recon_rows WHERE batch_id = $1 AND source = 'phsar' AND resolved_at IS NOT NULL`, batchID)
	if err != nil {
		return err
	}
	for rs.Next() {
		var t string
		var r resolution
		if err := rs.Scan(&t, &r.by, &r.at, &r.note); err != nil {
			rs.Close()
			return err
		}
		saved[t] = r
	}
	rs.Close()
	if _, err := tx.Exec(ctx, `DELETE FROM recon_rows WHERE batch_id = $1 AND source = 'phsar'`, batchID); err != nil {
		return err
	}

	rows, err := tx.Query(ctx, `SELECT id, kind, coalesce(tran_id, ''), coalesce(reference, ''), coalesce(status, ''), amount, result
		FROM recon_rows WHERE batch_id = $1 AND source = 'statement' ORDER BY line_no, id`, batchID)
	if err != nil {
		return err
	}
	var list []*stmtRow
	trans := []string{}
	for rows.Next() {
		r := &stmtRow{}
		if err := rows.Scan(&r.id, &r.kind, &r.tran, &r.ref, &r.status, &r.amount, &r.oldResult); err != nil {
			rows.Close()
			return err
		}
		list = append(list, r)
		trans = append(trans, r.tran)
	}
	rows.Close()

	// checkout menurut ID transaksi
	type co struct {
		id, status string
		expected   decimal.Decimal
	}
	cos := map[string]co{}
	cr, err := tx.Query(ctx, `SELECT tran_id, id::text, status, coalesce(paid_amount, total) FROM checkouts WHERE tran_id = ANY($1)`, trans)
	if err != nil {
		return err
	}
	for cr.Next() {
		var t string
		var c co
		if err := cr.Scan(&t, &c.id, &c.status, &c.expected); err != nil {
			cr.Close()
			return err
		}
		cos[t] = c
	}
	cr.Close()

	// refund yang dicatat admin: refund pesanan dibatalkan + retur yang dikembalikan
	type refund struct {
		ref, tran string
		amount    decimal.Decimal
	}
	var refunds []refund
	rr, err := tx.Query(ctx, `SELECT coalesce(o.refund_ref, ''), c.tran_id, coalesce(o.refund_amount, o.total) FROM orders o JOIN checkouts c ON c.id = o.checkout_id
			WHERE o.refund_status = 'refunded' AND (c.tran_id = ANY($1) OR upper(o.refund_ref) = ANY($2))
		UNION ALL SELECT coalesce(rt.refund_ref, ''), c.tran_id, rt.refund_amount FROM returns rt JOIN orders o ON o.id = rt.order_id JOIN checkouts c ON c.id = o.checkout_id
			WHERE rt.status = 'refunded' AND (c.tran_id = ANY($1) OR upper(rt.refund_ref) = ANY($2))`, trans, upperRefs(list))
	if err != nil {
		return err
	}
	for rr.Next() {
		var f refund
		if err := rr.Scan(&f.ref, &f.tran, &f.amount); err != nil {
			rr.Close()
			return err
		}
		refunds = append(refunds, f)
	}
	rr.Close()

	seen := map[string]bool{}
	paidTrans := []string{}
	for _, r := range list {
		r.result, r.checkoutID, r.expected = "", nil, nil
		if c, ok := cos[r.tran]; ok {
			id := c.id
			r.checkoutID = &id
		}
		if !success(r.status) {
			r.result = "ignored"
			continue
		}
		if r.kind == "refund" {
			var hit []refund
			for _, f := range refunds {
				if r.ref != "" && strings.EqualFold(f.ref, r.ref) {
					hit = []refund{f}
					break
				}
			}
			if hit == nil {
				for _, f := range refunds {
					if f.tran == r.tran {
						hit = append(hit, f)
					}
				}
			}
			if len(hit) == 0 {
				r.result = "refund_unrecorded"
				continue
			}
			sum := decimal.Zero
			for _, f := range hit {
				sum = sum.Add(f.amount)
			}
			r.expected = &sum
			if r.amount != nil && !r.amount.Equal(sum) {
				r.result = "refund_mismatch"
			} else {
				r.result = "refund_ok"
			}
			continue
		}
		if seen[r.tran] {
			r.result = "duplicate"
			continue
		}
		seen[r.tran] = true
		paidTrans = append(paidTrans, r.tran)
		c, ok := cos[r.tran]
		if !ok {
			r.result = "unknown"
			continue
		}
		exp := c.expected
		r.expected = &exp
		switch {
		case c.status != "paid":
			r.result = "not_paid_in_phsar"
		case r.amount != nil && !r.amount.Equal(c.expected):
			r.result = "amount_mismatch"
		default:
			r.result = "ok"
		}
	}
	if len(list) > 0 {
		ids, results, checkouts, expected, keep := make([]int64, len(list)), make([]string, len(list)), make([]*string, len(list)), make([]*string, len(list)), make([]bool, len(list))
		for i, r := range list {
			ids[i], results[i], checkouts[i], keep[i] = r.id, r.result, r.checkoutID, r.result == r.oldResult
			if r.expected != nil {
				v := r.expected.StringFixed(2)
				expected[i] = &v
			}
		}
		// penyelesaian lama hanya dipertahankan kalau hasilnya tidak berubah
		if _, err := tx.Exec(ctx, `UPDATE recon_rows r SET result = u.result, checkout_id = u.checkout_id::uuid, expected = u.expected::numeric,
				resolved_by = CASE WHEN u.keep THEN r.resolved_by END, resolved_at = CASE WHEN u.keep THEN r.resolved_at END,
				resolution = CASE WHEN u.keep THEN r.resolution END
			FROM unnest($1::bigint[], $2::text[], $3::text[], $4::text[], $5::bool[]) AS u(id, result, checkout_id, expected, keep)
			WHERE r.id = u.id`, ids, results, checkouts, expected, keep); err != nil {
			return err
		}
	}

	// lunas di Phsar dalam periode file tapi tidak ada di file (pembayaran palsu/dev tidak dihitung)
	mr, err := tx.Query(ctx, `SELECT c.id::text, c.tran_id, coalesce(c.paid_amount, c.total), c.paid_at FROM checkouts c, recon_batches b
		WHERE b.id = $1 AND b.period_from IS NOT NULL AND b.period_to IS NOT NULL AND c.status = 'paid' AND c.provider <> 'fake'
		  AND (c.paid_at AT TIME ZONE 'Asia/Phnom_Penh')::date BETWEEN b.period_from AND b.period_to AND NOT (c.tran_id = ANY($2))
		ORDER BY c.paid_at LIMIT 5000`, batchID, paidTrans)
	if err != nil {
		return err
	}
	type missing struct {
		id, tran string
		amount   decimal.Decimal
		at       time.Time
	}
	var miss []missing
	for mr.Next() {
		var m missing
		if err := mr.Scan(&m.id, &m.tran, &m.amount, &m.at); err != nil {
			mr.Close()
			return err
		}
		miss = append(miss, m)
	}
	mr.Close()
	for _, m := range miss {
		r := saved[m.tran]
		if _, err := tx.Exec(ctx, `INSERT INTO recon_rows (batch_id, source, kind, tran_id, status, occurred_at, result, checkout_id, expected,
			resolved_by, resolved_at, resolution) VALUES ($1, 'phsar', 'payment', $2, 'PAID IN PHSAR', $3, 'missing_in_statement', $4, $5, $6, $7, $8)`,
			batchID, m.tran, m.at, m.id, m.amount, r.by, r.at, r.note); err != nil {
			return err
		}
	}
	_, err = tx.Exec(ctx, `UPDATE recon_batches SET matched_at = now() WHERE id = $1`, batchID)
	return err
}

func upperRefs(list []*stmtRow) []string {
	out := []string{}
	for _, r := range list {
		if r.ref != "" {
			out = append(out, strings.ToUpper(r.ref))
		}
	}
	return out
}

// ---------------------------------------------------------------------------
// Tampilan (JSON disusun di Postgres)
// ---------------------------------------------------------------------------

const exceptionList = `('amount_mismatch','not_paid_in_phsar','unknown','duplicate','missing_in_statement','refund_mismatch','refund_unrecorded')`

// List: semua batch terbaru.
func (s *Service) List(ctx context.Context) ([]byte, error) {
	var raw []byte
	err := s.Pool.QueryRow(ctx, `SELECT coalesce(json_agg(x ORDER BY x.created_at DESC), '[]') FROM (
		SELECT b.id, b.file_name, b.period_from, b.period_to, b.row_count, b.collected::text AS collected, b.refunded::text AS refunded,
		  b.fees::text AS fees, b.created_at, b.matched_at, coalesce(nullif(u.name, ''), u.phone) AS uploaded_by,
		  (SELECT count(*) FROM recon_rows r WHERE r.batch_id = b.id AND r.result IN `+exceptionList+` AND r.resolved_at IS NULL) AS open_exceptions
		FROM recon_batches b JOIN users u ON u.id = b.uploaded_by ORDER BY b.created_at DESC LIMIT 100) x`).Scan(&raw)
	return raw, err
}

var validFilter = map[string]bool{"open": true, "exceptions": true, "resolved": true, "all": true}

// Get: satu batch + ringkasan + baris (filter: open · exceptions · resolved · all · <hasil>).
func (s *Service) Get(ctx context.Context, batchID, filter string) ([]byte, error) {
	if filter == "" {
		filter = "open"
	}
	where := ""
	switch {
	case filter == "open":
		where = "r.result IN " + exceptionList + " AND r.resolved_at IS NULL"
	case filter == "exceptions":
		where = "r.result IN " + exceptionList
	case filter == "resolved":
		where = "r.resolved_at IS NOT NULL"
	case filter == "all":
		where = "TRUE"
	case isException(filter) || filter == "ok" || filter == "ignored" || filter == "refund_ok":
		where = "r.result = '" + filter + "'"
	default:
		return nil, FieldError{"filter", "unknown filter"}
	}
	var raw []byte
	err := s.Pool.QueryRow(ctx, `SELECT json_build_object(
		'batch', json_build_object('id', b.id, 'file_name', b.file_name, 'period_from', b.period_from, 'period_to', b.period_to,
			'row_count', b.row_count, 'collected', b.collected::text, 'refunded', b.refunded::text, 'fees', b.fees::text,
			'net', (b.collected - b.refunded - b.fees)::text, 'created_at', b.created_at, 'matched_at', b.matched_at,
			'uploaded_by', coalesce(nullif(u.name, ''), u.phone)),
		'counts', (SELECT coalesce(json_object_agg(result, n), '{}') FROM (SELECT result, count(*) AS n FROM recon_rows WHERE batch_id = b.id GROUP BY result) c),
		'open', (SELECT count(*) FROM recon_rows r WHERE r.batch_id = b.id AND r.result IN `+exceptionList+` AND r.resolved_at IS NULL),
		'phsar_paid', (SELECT coalesce(sum(r.expected), 0)::text FROM recon_rows r WHERE r.batch_id = b.id AND r.result IN ('ok','amount_mismatch','missing_in_statement')),
		'rows', (SELECT coalesce(json_agg(x ORDER BY x.source DESC, x.line_no NULLS LAST, x.id), '[]') FROM (
			SELECT r.id, r.source, r.line_no, r.kind, r.tran_id, r.reference, r.status, r.amount::text AS amount, r.currency, r.fee::text AS fee,
			  r.occurred_at, r.result, r.expected::text AS expected, r.checkout_id, c.status AS checkout_status,
			  (SELECT coalesce(array_agg(o.order_no ORDER BY o.order_no), '{}') FROM orders o WHERE o.checkout_id = r.checkout_id) AS orders,
			  r.resolved_at, r.resolution, coalesce(nullif(ru.name, ''), ru.phone) AS resolved_by
			FROM recon_rows r LEFT JOIN checkouts c ON c.id = r.checkout_id LEFT JOIN users ru ON ru.id = r.resolved_by
			WHERE r.batch_id = b.id AND `+where+` ORDER BY r.id LIMIT 2000) x)
		) FROM recon_batches b JOIN users u ON u.id = b.uploaded_by WHERE b.id = $1`, batchID).Scan(&raw)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return raw, err
}

// Resolve: keuangan menjelaskan selisih (mis. "callback diproses manual", "salah input di PayWay").
func (s *Service) Resolve(ctx context.Context, actorID string, rowID int64, note string) (string, error) {
	note = strings.TrimSpace(note)
	if n := utf8.RuneCountInString(note); n < 3 || n > 1000 {
		return "", FieldError{"note", "explain how this difference was settled (3–1000 characters)"}
	}
	var batchID string
	err := db.WithActor(ctx, s.Pool, actorID, "admin", "", func(tx pgx.Tx) error {
		var result string
		err := tx.QueryRow(ctx, `SELECT batch_id::text, result FROM recon_rows WHERE id = $1 FOR UPDATE`, rowID).Scan(&batchID, &result)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		if !isException(result) {
			return ErrNotOpen
		}
		if _, err := tx.Exec(ctx, `UPDATE recon_rows SET resolved_by = $2, resolved_at = now(), resolution = $3 WHERE id = $1`, rowID, actorID, note); err != nil {
			return err
		}
		nd, _ := json.Marshal(map[string]string{"result": result, "resolution": note})
		_, err = tx.Exec(ctx, `INSERT INTO audit_logs (actor_type, actor_user_id, action, entity_table, entity_id, new_data)
			VALUES ('admin', $1, 'reconciliation.resolve', 'recon_rows', $2, $3)`, actorID, strconv.FormatInt(rowID, 10), nd)
		return err
	})
	return batchID, err
}

// Delete: hapus batch yang salah unggah (barisnya ikut terhapus).
func (s *Service) Delete(ctx context.Context, actorID, batchID string) error {
	return db.WithActor(ctx, s.Pool, actorID, "admin", "", func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, `DELETE FROM recon_batches WHERE id = $1`, batchID)
		if err == nil && tag.RowsAffected() == 0 {
			return ErrNotFound
		}
		return err
	})
}
