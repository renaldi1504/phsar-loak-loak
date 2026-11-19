// Package flash: sesi flash sale.
//
//	admin membuat sesi (nama 3 bahasa, jam mulai/selesai bebas, maks. 14 hari, pembiayaan per sesi)
//	produk masuk lewat dua jalur:
//	  - admin memilih banyak varian sekaligus → langsung "approved" kalau Phsar menanggung 100%,
//	    selain itu "invited" sampai seller menerima (seller tidak pernah menanggung potongan tanpa setuju)
//	  - seller mendaftarkan produknya selama pendaftaran dibuka → "pending" sampai admin menyetujui
//	Aturan barang: potongan minimal sesi, kuota ≤ stok tersedia, satu varian tidak boleh di dua sesi yang waktunya bertabrakan.
//	Harga di keranjang/checkout dihitung paket shop dari view flash_live.
package flash

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

	"phsar/internal/notify"
	"phsar/internal/platform/db"
	"phsar/internal/pricehistory" // pricehistory
)

var (
	ErrNotFound = errors.New("tidak ditemukan")
	ErrState    = errors.New("status tidak mengizinkan perubahan ini")
)

// FieldError: input tidak valid (422).
type FieldError struct{ Field, Msg string }

func (e FieldError) Error() string { return e.Field + ": " + e.Msg }

type Service struct {
	Pool  *pgxpool.Pool
	Files interface{ PublicURL(string) string }
	Log   *slog.Logger
}

// ---------------------------------------------------------------------------
// Sesi
// ---------------------------------------------------------------------------

type SessionInput struct {
	NameEN         string `json:"name_en"`
	NameKM         string `json:"name_km"`
	NameZH         string `json:"name_zh"`
	StartsAt       string `json:"starts_at"` // RFC3339
	EndsAt         string `json:"ends_at"`
	PhsarSharePct  int    `json:"phsar_share_pct"`
	AllowVouchers  bool   `json:"allow_vouchers"`
	MinDiscountPct int    `json:"min_discount_pct"`
	PerBuyer       int    `json:"per_buyer"`
	OpenForSellers bool   `json:"open_for_sellers"`
	RegisterUntil  string `json:"register_until"` // kosong = sampai sesi mulai
	Published      bool   `json:"published"`
}

func trimRunes(s string, n int) string {
	r := []rune(strings.TrimSpace(s))
	if len(r) > n {
		r = r[:n]
	}
	return string(r)
}

func parseTime(field, s string, required bool) (*time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		if required {
			return nil, FieldError{field, "required"}
		}
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil, FieldError{field, "invalid date/time"}
	}
	return &t, nil
}

// SaveSession: id kosong = buat baru. Jam yang sudah lewat tidak bisa dipindah (sesi berjalan/selesai hanya boleh diubah sebagian).
func (s *Service) SaveSession(ctx context.Context, actorID, id string, in SessionInput) (string, error) {
	in.NameEN, in.NameKM, in.NameZH = trimRunes(in.NameEN, 80), trimRunes(in.NameKM, 80), trimRunes(in.NameZH, 80)
	if in.NameEN == "" {
		return "", FieldError{"name_en", "required"}
	}
	start, err := parseTime("starts_at", in.StartsAt, true)
	if err != nil {
		return "", err
	}
	end, err := parseTime("ends_at", in.EndsAt, true)
	if err != nil {
		return "", err
	}
	if !end.After(*start) {
		return "", FieldError{"ends_at", "must be after the start"}
	}
	if end.Sub(*start) > 14*24*time.Hour {
		return "", FieldError{"ends_at", "a session lasts at most 14 days"}
	}
	reg, err := parseTime("register_until", in.RegisterUntil, false)
	if err != nil {
		return "", err
	}
	if reg != nil && reg.After(*start) {
		return "", FieldError{"register_until", "registration must close before the session starts"}
	}
	if in.PhsarSharePct < 0 || in.PhsarSharePct > 100 {
		return "", FieldError{"phsar_share_pct", "0–100"}
	}
	if in.MinDiscountPct == 0 {
		in.MinDiscountPct = 10
	}
	if in.MinDiscountPct < 1 || in.MinDiscountPct > 90 {
		return "", FieldError{"min_discount_pct", "1–90"}
	}
	if in.PerBuyer == 0 {
		in.PerBuyer = 2
	}
	if in.PerBuyer < 1 || in.PerBuyer > 99 {
		return "", FieldError{"per_buyer", "1–99"}
	}
	err = db.WithActor(ctx, s.Pool, actorID, "admin", "", func(tx pgx.Tx) error {
		if id == "" {
			if !end.After(time.Now()) {
				return FieldError{"ends_at", "must be in the future"}
			}
			return tx.QueryRow(ctx, `INSERT INTO flash_sessions (name_en, name_km, name_zh, starts_at, ends_at, phsar_share_pct, allow_vouchers,
				min_discount_pct, per_buyer, open_for_sellers, register_until, published, created_by)
				VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING id`, in.NameEN, in.NameKM, in.NameZH, *start, *end, in.PhsarSharePct,
				in.AllowVouchers, in.MinDiscountPct, in.PerBuyer, in.OpenForSellers, reg, in.Published, actorID).Scan(&id)
		}
		var oldStart, oldEnd time.Time
		var oldPct int
		var cancelled *time.Time
		err := tx.QueryRow(ctx, `SELECT starts_at, ends_at, phsar_share_pct, cancelled_at FROM flash_sessions WHERE id::text = $1 FOR UPDATE`, id).
			Scan(&oldStart, &oldEnd, &oldPct, &cancelled)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		if cancelled != nil || !oldEnd.After(time.Now()) {
			return ErrState
		}
		started := !oldStart.After(time.Now())
		if started {
			// sesi sudah berjalan: jam mulai & pembiayaan dikunci (pembeli & seller sudah memakai aturan ini)
			if !start.Equal(oldStart) {
				return FieldError{"starts_at", "the session has started; the start time cannot change"}
			}
			if in.PhsarSharePct != oldPct {
				return FieldError{"phsar_share_pct", "the session has started; funding cannot change"}
			}
		}
		if _, err := tx.Exec(ctx, `UPDATE flash_sessions SET name_en=$2, name_km=$3, name_zh=$4, starts_at=$5, ends_at=$6, phsar_share_pct=$7,
			allow_vouchers=$8, min_discount_pct=$9, per_buyer=$10, open_for_sellers=$11, register_until=$12, published=$13 WHERE id::text = $1`,
			id, in.NameEN, in.NameKM, in.NameZH, *start, *end, in.PhsarSharePct, in.AllowVouchers, in.MinDiscountPct, in.PerBuyer,
			in.OpenForSellers, reg, in.Published); err != nil {
			return err
		}
		if !start.Equal(oldStart) || !end.Equal(oldEnd) {
			// jam berubah: barang yang kini bertabrakan dengan sesi lain dikeluarkan
			if _, err := tx.Exec(ctx, `UPDATE flash_items i SET status = 'removed', note = 'time overlaps another session', decided_at = now()
				WHERE i.session_id::text = $1 AND i.status IN ('invited','pending','approved') AND EXISTS (
				  SELECT 1 FROM flash_items o JOIN flash_sessions os ON os.id = o.session_id
				  WHERE o.variant_id = i.variant_id AND o.id <> i.id AND o.status IN ('invited','pending','approved') AND os.cancelled_at IS NULL
				    AND os.starts_at < $3 AND os.ends_at > $2)`, id, *start, *end); err != nil {
				return err
			}
		}
		if oldPct > 0 && in.PhsarSharePct < oldPct {
			// Phsar menanggung lebih sedikit dari yang disetujui seller → seller harus setuju lagi
			if _, err := tx.Exec(ctx, `UPDATE flash_items SET status = 'invited', note = 'funding changed', decided_at = now()
				WHERE session_id::text = $1 AND status IN ('approved','pending')`, id); err != nil {
				return err
			}
		} else if in.PhsarSharePct < 100 && oldPct == 100 {
			if _, err := tx.Exec(ctx, `UPDATE flash_items SET status = 'invited', note = 'funding changed', decided_at = now()
				WHERE session_id::text = $1 AND status = 'approved' AND source = 'admin'`, id); err != nil {
				return err
			}
		}
		return nil
	})
	return id, err
}

// CancelSession: sesi dibatalkan (harga flash langsung berhenti). Data tetap disimpan untuk laporan.
func (s *Service) CancelSession(ctx context.Context, actorID, id string) error {
	return db.WithActor(ctx, s.Pool, actorID, "admin", "", func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, `UPDATE flash_sessions SET cancelled_at = now(), published = false
			WHERE id::text = $1 AND cancelled_at IS NULL AND ends_at > now()`, id)
		if err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			return ErrState
		}
		return nil
	})
}

const sessionCols = `s.id, s.name_en, s.name_km, s.name_zh, s.starts_at, s.ends_at, s.phsar_share_pct, s.allow_vouchers, s.min_discount_pct,
	s.per_buyer, s.open_for_sellers, s.register_until, s.published, s.cancelled_at,
	CASE WHEN s.cancelled_at IS NOT NULL THEN 'cancelled' WHEN s.ends_at <= now() THEN 'ended' WHEN NOT s.published THEN 'draft'
	     WHEN s.starts_at <= now() THEN 'live' ELSE 'scheduled' END`

type Session struct {
	ID             string     `json:"id"`
	NameEN         string     `json:"name_en"`
	NameKM         string     `json:"name_km"`
	NameZH         string     `json:"name_zh"`
	StartsAt       time.Time  `json:"starts_at"`
	EndsAt         time.Time  `json:"ends_at"`
	PhsarSharePct  int        `json:"phsar_share_pct"`
	AllowVouchers  bool       `json:"allow_vouchers"`
	MinDiscountPct int        `json:"min_discount_pct"`
	PerBuyer       int        `json:"per_buyer"`
	OpenForSellers bool       `json:"open_for_sellers"`
	RegisterUntil  *time.Time `json:"register_until"`
	Published      bool       `json:"published"`
	CancelledAt    *time.Time `json:"cancelled_at"`
	Status         string     `json:"status"` // draft · scheduled · live · ended · cancelled
	Counts         *Counts    `json:"counts,omitempty"`
}

type Counts struct {
	Approved int    `json:"approved"`
	Pending  int    `json:"pending"`
	Invited  int    `json:"invited"`
	Sold     int    `json:"sold"`
	Revenue  string `json:"revenue"`  // nilai terjual dengan harga flash
	Subsidy  string `json:"subsidy"`  // ditanggung Phsar
	Discount string `json:"discount"` // total potongan (normal − flash)
}

func scanSession(row pgx.Row, withCounts bool) (*Session, error) {
	var x Session
	dest := []any{&x.ID, &x.NameEN, &x.NameKM, &x.NameZH, &x.StartsAt, &x.EndsAt, &x.PhsarSharePct, &x.AllowVouchers, &x.MinDiscountPct,
		&x.PerBuyer, &x.OpenForSellers, &x.RegisterUntil, &x.Published, &x.CancelledAt, &x.Status}
	if withCounts {
		x.Counts = &Counts{}
		dest = append(dest, &x.Counts.Approved, &x.Counts.Pending, &x.Counts.Invited, &x.Counts.Sold, &x.Counts.Revenue, &x.Counts.Subsidy, &x.Counts.Discount)
	}
	return &x, row.Scan(dest...)
}

const countCols = `,
	(SELECT count(*) FROM flash_items i WHERE i.session_id = s.id AND i.status = 'approved')::int,
	(SELECT count(*) FROM flash_items i WHERE i.session_id = s.id AND i.status = 'pending')::int,
	(SELECT count(*) FROM flash_items i WHERE i.session_id = s.id AND i.status = 'invited')::int,
	coalesce((SELECT sum(oi.qty) FROM order_items oi JOIN orders o ON o.id = oi.order_id JOIN flash_items i ON i.id = oi.flash_item_id
	  WHERE i.session_id = s.id AND o.status NOT IN ('cancelled','refunded')), 0)::int,
	coalesce((SELECT sum(oi.line_total) FROM order_items oi JOIN orders o ON o.id = oi.order_id JOIN flash_items i ON i.id = oi.flash_item_id
	  WHERE i.session_id = s.id AND o.status NOT IN ('cancelled','refunded')), 0)::text,
	coalesce((SELECT sum(oi.flash_subsidy) FROM order_items oi JOIN orders o ON o.id = oi.order_id JOIN flash_items i ON i.id = oi.flash_item_id
	  WHERE i.session_id = s.id AND o.status NOT IN ('cancelled','refunded')), 0)::text,
	coalesce((SELECT sum((i.normal_price - oi.unit_price) * oi.qty) FROM order_items oi JOIN orders o ON o.id = oi.order_id JOIN flash_items i ON i.id = oi.flash_item_id
	  WHERE i.session_id = s.id AND o.status NOT IN ('cancelled','refunded')), 0)::text`

// Sessions: daftar untuk admin (terbaru dulu, sesi yang belum selesai di atas).
func (s *Service) Sessions(ctx context.Context) ([]Session, error) {
	rows, err := s.Pool.Query(ctx, `SELECT `+sessionCols+countCols+` FROM flash_sessions s
		ORDER BY (s.ends_at > now() AND s.cancelled_at IS NULL) DESC, s.starts_at DESC LIMIT 100`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Session{}
	for rows.Next() {
		x, err := scanSession(rows, true)
		if err != nil {
			return nil, err
		}
		out = append(out, *x)
	}
	return out, rows.Err()
}

func (s *Service) Session(ctx context.Context, id string) (*Session, error) {
	x, err := scanSession(s.Pool.QueryRow(ctx, `SELECT `+sessionCols+countCols+` FROM flash_sessions s WHERE s.id::text = $1`, id), true)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return x, err
}

// ---------------------------------------------------------------------------
// Barang
// ---------------------------------------------------------------------------

type Item struct {
	ID          string     `json:"id"`
	SessionID   string     `json:"session_id"`
	VariantID   string     `json:"variant_id"`
	ProductID   string     `json:"product_id"`
	ProductName string     `json:"product_name"`
	Options     string     `json:"options"`
	ImageURL    string     `json:"image_url"`
	StoreID     string     `json:"store_id"`
	StoreName   string     `json:"store_name"`
	NormalPrice string     `json:"normal_price"`
	SalePrice   string     `json:"sale_price"`
	CurrentPrc  string     `json:"current_price"`
	OffPct      int        `json:"off_pct"`
	QtyCap      int        `json:"qty_cap"`
	PerBuyer    *int       `json:"per_buyer"`
	Stock       int        `json:"stock"`
	Sold        int        `json:"sold"`
	Source      string     `json:"source"`
	Status      string     `json:"status"`
	Note        *string    `json:"note"`
	CreatedAt   time.Time  `json:"created_at"`
	DecidedAt   *time.Time `json:"decided_at"`
	Session     *Session   `json:"session,omitempty"` // daftar milik seller
}

const itemSelect = `SELECT i.id, i.session_id, i.variant_id, i.product_id,
	coalesce((SELECT t.name FROM product_translations t WHERE t.product_id = p.id AND t.lang = p.source_lang), ''),
	coalesce((SELECT string_agg(coalesce(ov.value_i18n->>'en', ov.value_i18n->>p.source_lang::text), ' / ' ORDER BY o.position)
	  FROM variant_option_values vov JOIN product_option_values ov ON ov.id = vov.option_value_id JOIN product_options o ON o.id = vov.option_id
	  WHERE vov.variant_id = v.id), ''),
	coalesce((SELECT coalesce(m.image_url, m.storage_key) FROM product_media m WHERE m.product_id = p.id AND m.deleted_at IS NULL AND m.kind = 'photo'
	  ORDER BY m.sort_order LIMIT 1), ''),
	i.store_id, st.name, i.normal_price::text, i.sale_price::text, v.price::text,
	round((1 - i.sale_price / i.normal_price) * 100)::int, i.qty_cap, i.per_buyer, greatest(v.stock_qty - v.reserved_qty, 0),
	coalesce((SELECT sum(oi.qty) FROM order_items oi JOIN orders o ON o.id = oi.order_id WHERE oi.flash_item_id = i.id
	  AND o.status NOT IN ('cancelled','refunded')), 0)::int,
	i.source, i.status, i.note, i.created_at, i.decided_at
	FROM flash_items i JOIN product_variants v ON v.id = i.variant_id JOIN products p ON p.id = i.product_id JOIN stores st ON st.id = i.store_id`

func (s *Service) items(ctx context.Context, where string, args ...any) ([]Item, error) {
	rows, err := s.Pool.Query(ctx, itemSelect+` WHERE `+where+` ORDER BY i.created_at DESC LIMIT 500`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Item{}
	for rows.Next() {
		var x Item
		if err := rows.Scan(&x.ID, &x.SessionID, &x.VariantID, &x.ProductID, &x.ProductName, &x.Options, &x.ImageURL, &x.StoreID, &x.StoreName,
			&x.NormalPrice, &x.SalePrice, &x.CurrentPrc, &x.OffPct, &x.QtyCap, &x.PerBuyer, &x.Stock, &x.Sold, &x.Source, &x.Status, &x.Note,
			&x.CreatedAt, &x.DecidedAt); err != nil {
			return nil, err
		}
		if x.ImageURL != "" && !strings.HasPrefix(x.ImageURL, "http") && !strings.HasPrefix(x.ImageURL, "/") && s.Files != nil {
			x.ImageURL = s.Files.PublicURL(x.ImageURL)
		}
		out = append(out, x)
	}
	return out, rows.Err()
}

// SessionItems: semua barang sesi (admin), yang aktif di atas.
func (s *Service) SessionItems(ctx context.Context, sessionID string) ([]Item, error) {
	return s.items(ctx, `i.session_id::text = $1`, sessionID)
}

type ItemInput struct {
	VariantID string `json:"variant_id"`
	SalePrice string `json:"sale_price"`
	QtyCap    int    `json:"qty_cap"`
	PerBuyer  *int   `json:"per_buyer"`
}

type sessionRow struct {
	id                     string
	start, end             time.Time
	pct, minDisc           int
	openSellers, cancelled bool
	registerUntil          *time.Time
	nameEN                 string
}

func lockSession(ctx context.Context, tx pgx.Tx, id string) (*sessionRow, error) {
	var x sessionRow
	var cancelled *time.Time
	err := tx.QueryRow(ctx, `SELECT id, starts_at, ends_at, phsar_share_pct, min_discount_pct, open_for_sellers, register_until, cancelled_at, name_en
		FROM flash_sessions WHERE id::text = $1 FOR UPDATE`, id).Scan(&x.id, &x.start, &x.end, &x.pct, &x.minDisc, &x.openSellers, &x.registerUntil, &cancelled, &x.nameEN)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	x.cancelled = cancelled != nil
	return &x, err
}

type variantRow struct {
	productID, storeID, name string
	price                    decimal.Decimal
	stock                    int
}

// checkItem: aturan satu barang. storeID != "" = jalur seller (varian harus milik tokonya).
func checkItem(ctx context.Context, tx pgx.Tx, ses *sessionRow, storeID string, n int, in ItemInput) (*variantRow, decimal.Decimal, error) {
	f := func(field, msg string) error { return FieldError{fmt.Sprintf("items.%d.%s", n, field), msg} }
	var v variantRow
	var ok bool
	err := tx.QueryRow(ctx, `SELECT p.id, p.store_id, coalesce((SELECT t.name FROM product_translations t WHERE t.product_id = p.id AND t.lang = p.source_lang), ''),
		v.price, greatest(v.stock_qty - v.reserved_qty, 0),
		v.deleted_at IS NULL AND v.status = 'active' AND p.deleted_at IS NULL AND p.status = 'active' AND p.moderation_status = 'approved' AND st.status = 'active'
		FROM product_variants v JOIN products p ON p.id = v.product_id JOIN stores st ON st.id = p.store_id WHERE v.id::text = $1`, in.VariantID).
		Scan(&v.productID, &v.storeID, &v.name, &v.price, &v.stock, &ok)
	if errors.Is(err, pgx.ErrNoRows) || (err == nil && storeID != "" && v.storeID != storeID) {
		return nil, decimal.Zero, f("variant_id", "product not found")
	}
	if err != nil {
		return nil, decimal.Zero, err
	}
	if !ok {
		return nil, decimal.Zero, f("variant_id", v.name+": product is not live")
	}
	price, err := decimal.NewFromString(strings.TrimSpace(in.SalePrice))
	if err != nil || !price.IsPositive() || price.Exponent() < -2 {
		return nil, decimal.Zero, f("sale_price", v.name+": enter a price like 4.50")
	}
	maxPrice := v.price.Mul(decimal.NewFromInt(int64(100 - ses.minDisc))).Div(decimal.NewFromInt(100)).RoundDown(2)
	if price.GreaterThan(maxPrice) {
		return nil, decimal.Zero, f("sale_price", fmt.Sprintf("%s: at least %d%% off — at most $%s (normal $%s)", v.name, ses.minDisc, maxPrice.StringFixed(2), v.price.StringFixed(2)))
	}
	// pricehistory: harga flash harus di bawah harga normal terendah 30 hari sebelum sesi mulai (bisa dimatikan super admin)
	if msg, err := pricehistory.CheckFlashPrice(ctx, tx, in.VariantID, ses.start, price, v.name); err != nil {
		return nil, decimal.Zero, err
	} else if msg != "" {
		return nil, decimal.Zero, f("sale_price", msg)
	}
	if in.QtyCap < 1 {
		return nil, decimal.Zero, f("qty_cap", v.name+": quota at least 1")
	}
	if in.QtyCap > v.stock {
		return nil, decimal.Zero, f("qty_cap", fmt.Sprintf("%s: only %d in stock", v.name, v.stock))
	}
	if in.PerBuyer != nil && (*in.PerBuyer < 1 || *in.PerBuyer > 99) {
		return nil, decimal.Zero, f("per_buyer", "1–99")
	}
	var clash bool
	if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM flash_items i JOIN flash_sessions s ON s.id = i.session_id
		WHERE i.variant_id::text = $1 AND i.status IN ('invited','pending','approved') AND s.cancelled_at IS NULL
		  AND s.starts_at < $3 AND s.ends_at > $2 AND s.id::text <> $4)`, in.VariantID, ses.start, ses.end, ses.id).Scan(&clash); err != nil {
		return nil, decimal.Zero, err
	}
	if clash {
		return nil, decimal.Zero, f("variant_id", v.name+": already in another flash sale at this time")
	}
	return &v, price, nil
}

// AddItems (admin): banyak varian sekaligus. Semua atau tidak sama sekali (pesan menyebut produk yang bermasalah).
func (s *Service) AddItems(ctx context.Context, actorID, sessionID string, list []ItemInput) (int, error) {
	if len(list) == 0 || len(list) > 200 {
		return 0, FieldError{"items", "choose 1–200 products"}
	}
	added := 0
	err := db.WithActor(ctx, s.Pool, actorID, "admin", "", func(tx pgx.Tx) error {
		ses, err := lockSession(ctx, tx, sessionID)
		if err != nil {
			return err
		}
		if ses.cancelled || !ses.end.After(time.Now()) {
			return ErrState
		}
		status := "invited"
		if ses.pct == 100 {
			status = "approved" // Phsar menanggung semua: seller tetap menerima harga normal
		}
		seen := map[string]bool{}
		invited := map[string]int{}
		for n, in := range list {
			if seen[in.VariantID] {
				return FieldError{fmt.Sprintf("items.%d.variant_id", n), "chosen twice"}
			}
			seen[in.VariantID] = true
			v, price, err := checkItem(ctx, tx, ses, "", n, in)
			if err != nil {
				return err
			}
			// varian yang sudah ada di sesi ini (aktif) → perbarui harga/kuota
			tag, err := tx.Exec(ctx, `UPDATE flash_items SET sale_price = $3, qty_cap = $4, per_buyer = $5, normal_price = $6,
				status = CASE WHEN $7 = 'approved' OR status = 'pending' THEN $7 ELSE 'invited' END, decided_by = $8, decided_at = now()
				WHERE session_id::text = $1 AND variant_id::text = $2 AND status IN ('invited','pending','approved')`,
				sessionID, in.VariantID, price, in.QtyCap, in.PerBuyer, v.price, status, actorID)
			if err != nil {
				return err
			}
			if tag.RowsAffected() == 0 {
				if _, err := tx.Exec(ctx, `INSERT INTO flash_items (session_id, variant_id, product_id, store_id, normal_price, sale_price, qty_cap, per_buyer,
					source, status, submitted_by, decided_by, decided_at)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'admin', $9, $10, $10, now())`,
					sessionID, in.VariantID, v.productID, v.storeID, v.price, price, in.QtyCap, in.PerBuyer, status, actorID); err != nil {
					return err
				}
			}
			if status == "invited" {
				invited[v.storeID]++
			}
			added++
		}
		for store, n := range invited {
			if err := notify.Enqueue(ctx, tx, store, notify.FlashInvited, notify.Params{Product: fmt.Sprint(n), Note: ses.nameEN},
				"#/store/"+store+"/flash"); err != nil {
				return err
			}
		}
		return nil
	})
	return added, err
}

// UpdateItem (admin): kuota & batas per pembeli kapan saja; harga hanya sebelum sesi mulai
// (kalau seller ikut menanggung, harga baru harus disetujui seller lagi).
func (s *Service) UpdateItem(ctx context.Context, actorID, itemID string, in ItemInput) error {
	return db.WithActor(ctx, s.Pool, actorID, "admin", "", func(tx pgx.Tx) error {
		var sessionID, variantID, status, oldPrice string
		if err := tx.QueryRow(ctx, `SELECT session_id::text, variant_id::text, status, sale_price::text FROM flash_items WHERE id::text = $1 FOR UPDATE`, itemID).
			Scan(&sessionID, &variantID, &status, &oldPrice); errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		} else if err != nil {
			return err
		}
		if status != "invited" && status != "pending" && status != "approved" {
			return ErrState
		}
		ses, err := lockSession(ctx, tx, sessionID)
		if err != nil {
			return err
		}
		if ses.cancelled || !ses.end.After(time.Now()) {
			return ErrState
		}
		in.VariantID = variantID
		v, price, err := checkItem(ctx, tx, ses, "", 0, in)
		if err != nil {
			var fe FieldError
			if errors.As(err, &fe) {
				fe.Field = strings.TrimPrefix(fe.Field, "items.0.")
				return fe
			}
			return err
		}
		changed := !price.Equal(decimal.RequireFromString(oldPrice))
		if changed && !ses.start.After(time.Now()) {
			return FieldError{"sale_price", "the session has started; the price cannot change"}
		}
		newStatus := status
		if changed && ses.pct < 100 && status == "approved" {
			newStatus = "invited"
		}
		_, err = tx.Exec(ctx, `UPDATE flash_items SET sale_price = $2, qty_cap = $3, per_buyer = $4, normal_price = $5, status = $6, decided_by = $7, decided_at = now()
			WHERE id::text = $1`, itemID, price, in.QtyCap, in.PerBuyer, v.price, newStatus, actorID)
		return err
	})
}

// Decide (admin): setujui / tolak pendaftaran seller; atau keluarkan barang (remove) kapan saja.
func (s *Service) Decide(ctx context.Context, actorID, itemID, action, note string) error {
	note = trimRunes(note, 300)
	return db.WithActor(ctx, s.Pool, actorID, "admin", "", func(tx pgx.Tx) error {
		var store, status, session, product string
		if err := tx.QueryRow(ctx, `SELECT i.store_id::text, i.status, s.name_en, coalesce((SELECT t.name FROM product_translations t JOIN products p ON p.id = t.product_id
			WHERE t.product_id = i.product_id AND t.lang = p.source_lang), '')
			FROM flash_items i JOIN flash_sessions s ON s.id = i.session_id WHERE i.id::text = $1 FOR UPDATE OF i`, itemID).Scan(&store, &status, &session, &product); errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		} else if err != nil {
			return err
		}
		var next, kind string
		switch {
		case action == "approve" && status == "pending":
			next, kind = "approved", notify.FlashApproved
		case action == "reject" && status == "pending":
			if len([]rune(note)) < 3 {
				return FieldError{"note", "tell the seller why"}
			}
			next, kind = "rejected", notify.FlashRejected
		case action == "remove" && (status == "approved" || status == "invited" || status == "pending"):
			next = "removed"
			if status == "approved" {
				kind = notify.FlashRejected
			}
		default:
			if action != "approve" && action != "reject" && action != "remove" {
				return FieldError{"action", "approve, reject or remove"}
			}
			return ErrState
		}
		if _, err := tx.Exec(ctx, `UPDATE flash_items SET status = $2, note = nullif($3, ''), decided_by = $4, decided_at = now() WHERE id::text = $1`,
			itemID, next, note, actorID); err != nil {
			return err
		}
		if kind != "" {
			return notify.Enqueue(ctx, tx, store, kind, notify.Params{Product: product, Note: strings.TrimSpace(session + " " + note)}, "#/store/"+store+"/flash")
		}
		return nil
	})
}

// ---------------------------------------------------------------------------
// Seller
// ---------------------------------------------------------------------------

// SellerView: sesi yang membuka pendaftaran / mengundang toko ini + semua barang toko di sesi yang belum lewat lama.
func (s *Service) SellerView(ctx context.Context, storeID string) (map[string]any, error) {
	rows, err := s.Pool.Query(ctx, `SELECT `+sessionCols+` FROM flash_sessions s WHERE s.cancelled_at IS NULL AND s.ends_at > now()
		AND (s.open_for_sellers OR EXISTS (SELECT 1 FROM flash_items i WHERE i.session_id = s.id AND i.store_id::text = $1))
		ORDER BY s.starts_at`, storeID)
	if err != nil {
		return nil, err
	}
	sessions := []map[string]any{}
	for rows.Next() {
		x, err := scanSession(rows, false)
		if err != nil {
			rows.Close()
			return nil, err
		}
		until := x.StartsAt
		if x.RegisterUntil != nil {
			until = *x.RegisterUntil
		}
		sessions = append(sessions, map[string]any{"session": x, "can_register": x.OpenForSellers && time.Now().Before(until), "register_until": until})
	}
	rows.Close()
	items, err := s.items(ctx, `i.store_id::text = $1 AND i.created_at > now() - interval '90 days'`, storeID)
	if err != nil {
		return nil, err
	}
	return map[string]any{"sessions": sessions, "items": items}, nil
}

// SellerVariants: varian toko yang bisa didaftarkan (produk tayang, stok > 0).
// pricehistory: low30 = harga normal terendah 30 hari sebelum `at` (jam mulai sesi; nol = sekarang).
func (s *Service) SellerVariants(ctx context.Context, storeID, q string, at time.Time) ([]map[string]any, error) {
	if at.IsZero() {
		at = time.Now()
	}
	var raw []byte
	err := s.Pool.QueryRow(ctx, `SELECT coalesce(json_agg(x), '[]') FROM (
		SELECT v.id AS variant_id, p.id AS product_id, coalesce(t.name, '') AS name, v.sku, v.price::text AS price,
		  greatest(v.stock_qty - v.reserved_qty, 0) AS stock, variant_low_price(v.id, $3)::text AS low30,
		  coalesce((SELECT string_agg(coalesce(ov.value_i18n->>'en', ov.value_i18n->>p.source_lang::text), ' / ' ORDER BY o.position)
		    FROM variant_option_values vov JOIN product_option_values ov ON ov.id = vov.option_value_id JOIN product_options o ON o.id = vov.option_id
		    WHERE vov.variant_id = v.id), '') AS options
		FROM product_variants v JOIN products p ON p.id = v.product_id
		LEFT JOIN product_translations t ON t.product_id = p.id AND t.lang = p.source_lang
		WHERE p.store_id::text = $1 AND v.deleted_at IS NULL AND v.status = 'active' AND p.deleted_at IS NULL AND p.status = 'active'
		  AND p.moderation_status = 'approved' AND v.stock_qty - v.reserved_qty > 0
		  AND ($2 = '' OR t.name ILIKE '%' || $2 || '%' OR v.sku ILIKE '%' || $2 || '%')
		ORDER BY t.name, v.price LIMIT 100) x`, storeID, strings.TrimSpace(q), at).Scan(&raw)
	if err != nil {
		return nil, err
	}
	out := []map[string]any{}
	_ = json.Unmarshal(raw, &out)
	return out, nil
}

// Submit (seller): daftarkan barang ke sesi yang membuka pendaftaran → pending.
func (s *Service) Submit(ctx context.Context, storeID, userID, sessionID string, list []ItemInput) (int, error) {
	if len(list) == 0 || len(list) > 50 {
		return 0, FieldError{"items", "choose 1–50 products"}
	}
	err := db.WithActor(ctx, s.Pool, userID, "seller", "", func(tx pgx.Tx) error {
		ses, err := lockSession(ctx, tx, sessionID)
		if err != nil {
			return err
		}
		until := ses.start
		if ses.registerUntil != nil {
			until = *ses.registerUntil
		}
		if ses.cancelled || !ses.openSellers || !time.Now().Before(until) {
			return ErrState
		}
		for n, in := range list {
			v, price, err := checkItem(ctx, tx, ses, storeID, n, in)
			if err != nil {
				return err
			}
			var exists bool
			if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM flash_items WHERE session_id::text = $1 AND variant_id::text = $2
				AND status IN ('invited','pending','approved'))`, sessionID, in.VariantID).Scan(&exists); err != nil {
				return err
			}
			if exists {
				return FieldError{fmt.Sprintf("items.%d.variant_id", n), v.name + ": already in this session"}
			}
			if _, err := tx.Exec(ctx, `INSERT INTO flash_items (session_id, variant_id, product_id, store_id, normal_price, sale_price, qty_cap, per_buyer,
				source, status, submitted_by) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'seller', 'pending', $9)`,
				sessionID, in.VariantID, v.productID, storeID, v.price, price, in.QtyCap, in.PerBuyer, userID); err != nil {
				return err
			}
		}
		return nil
	})
	return len(list), err
}

// SellerAction: accept / decline undangan admin; withdraw pendaftaran/barang sebelum sesi mulai.
func (s *Service) SellerAction(ctx context.Context, storeID, userID, itemID, action string) error {
	return db.WithActor(ctx, s.Pool, userID, "seller", "", func(tx pgx.Tx) error {
		var status string
		var start, end time.Time
		var cancelled *time.Time
		err := tx.QueryRow(ctx, `SELECT i.status, s.starts_at, s.ends_at, s.cancelled_at FROM flash_items i JOIN flash_sessions s ON s.id = i.session_id
			WHERE i.id::text = $1 AND i.store_id::text = $2 FOR UPDATE OF i`, itemID, storeID).Scan(&status, &start, &end, &cancelled)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		if cancelled != nil || !end.After(time.Now()) {
			return ErrState
		}
		var next string
		switch {
		case action == "accept" && status == "invited":
			next = "approved"
		case action == "decline" && status == "invited":
			next = "declined"
		case action == "withdraw" && (status == "pending" || status == "approved") && start.After(time.Now()):
			next = "withdrawn"
		default:
			if action != "accept" && action != "decline" && action != "withdraw" {
				return FieldError{"action", "accept, decline or withdraw"}
			}
			return ErrState
		}
		_, err = tx.Exec(ctx, `UPDATE flash_items SET status = $2, decided_at = now() WHERE id::text = $1`, itemID, next)
		return err
	})
}
