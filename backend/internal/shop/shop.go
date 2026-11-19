// Package shop: keranjang, checkout, pembayaran, dan pesanan dari sisi pembeli.
//
// Keranjang disimpan di browser/app (tanpa login). Server menghitung ulang semuanya: harga, stok,
// toko aktif — angka dari klien tidak pernah dipercaya. Satu checkout dibayar sekali (KHQR) lalu
// dipecah menjadi satu pesanan per toko.
package shop

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"phsar/internal/geo"
	"phsar/internal/hours"
	"phsar/internal/jobs"
	"phsar/internal/money"
	"phsar/internal/notify"
	"phsar/internal/payment"
	"phsar/internal/platform/db"
	"phsar/internal/storage"
)

var (
	ErrValidation   = errors.New("validation")
	ErrNotFound     = errors.New("not found")
	ErrCartChanged  = errors.New("cart changed")
	ErrNotServed    = errors.New("address outside service area")
	ErrPayment      = errors.New("payment provider unavailable")
	ErrNotCancelled = errors.New("checkout can no longer be cancelled")
	ErrClosed       = errors.New("checkout closed (pre-launch)")
)

type validationError struct{ field, msg string }

func (e validationError) Error() string { return e.field + ": " + e.msg }
func (e validationError) Unwrap() error { return ErrValidation }
func invalid(field, msg string) error   { return validationError{field, msg} }

const (
	maxLines = 50
	maxQty   = 99
)

// Produk yang boleh dibeli: sama dengan yang tampil di website + varian aktif.
const buyable = `p.deleted_at IS NULL AND p.status = 'active' AND p.moderation_status = 'approved'
	AND st.status = 'active' AND st.deleted_at IS NULL AND v.deleted_at IS NULL AND v.status = 'active'`

type Service struct {
	Pool        *pgxpool.Pool
	Pay         payment.Provider
	Files       storage.Storage
	BaseURL     string
	LifetimeMin int
	Log         *slog.Logger
	Closed      bool // pra-peluncuran: checkout ditolak (kalau Pay tidak punya pengaturan sendiri)
}

// gate: penyedia pembayaran yang diatur admin (paygw) — buka/tutup checkout & lama QR berlaku.
type gate interface {
	Open(ctx context.Context) bool
	Lifetime(ctx context.Context) int
}

func (s *Service) open(ctx context.Context) bool {
	if g, ok := s.Pay.(gate); ok {
		return g.Open(ctx)
	}
	return !s.Closed
}

func (s *Service) lifetime(ctx context.Context) int {
	if g, ok := s.Pay.(gate); ok {
		return g.Lifetime(ctx)
	}
	return s.LifetimeMin
}

// ---------------------------------------------------------------------------
// Keranjang
// ---------------------------------------------------------------------------

type CartItem struct {
	VariantID string `json:"variant_id"`
	Qty       int    `json:"qty"`
}

type Line struct {
	VariantID  string `json:"variant_id"`
	ProductID  string `json:"product_id"`
	Slug       string `json:"slug"`
	Name       string `json:"name"`
	Options    string `json:"options"`
	SKU        string `json:"-"`
	CategoryID string `json:"-"`
	ImageURL   string `json:"image_url"`
	Price      string `json:"price"`
	Compare    string `json:"compare,omitempty"`
	Qty        int    `json:"qty"`       // jumlah yang bisa dibeli (bisa dikurangi dari permintaan)
	Requested  int    `json:"requested"` // jumlah yang diminta
	Available  int    `json:"available"`
	LineTotal  string `json:"line_total"`
	// Problem: "" · unavailable (produk/toko tidak aktif) · out_of_stock · qty_reduced
	Problem string `json:"problem,omitempty"`
	// promos: harga flash sale sedang berlaku (Compare = harga normal).
	FlashEndsAt   *time.Time      `json:"flash_ends_at,omitempty"`
	FlashLeft     *int            `json:"flash_left,omitempty"`
	FlashPerBuyer int             `json:"flash_per_buyer,omitempty"`
	FlashItemID   string          `json:"-"`
	flashGap      decimal.Decimal // harga normal − harga flash (per barang)
	flashPct      int             // bagian potongan yang ditanggung Phsar (%)
	flashNoVouch  bool            // sesi tidak boleh digabung voucher

	price decimal.Decimal
}

type StoreGroup struct {
	StoreID   string `json:"store_id"`
	StoreName string `json:"store_name"`
	StoreSlug string `json:"store_slug"`
	Lines     []Line `json:"lines"`
	Subtotal  string `json:"subtotal"`
	// shipping: ongkir toko ini (paket shipping)
	ShippingFee string `json:"shipping_fee"`

	subtotal decimal.Decimal
	shipping decimal.Decimal
	count    int
}

type Quote struct {
	Stores      []StoreGroup `json:"stores"`
	Subtotal    string       `json:"subtotal"`
	ShippingFee string       `json:"shipping_fee"`
	Total       string       `json:"total"`
	ItemCount   int          `json:"item_count"`
	HasProblems bool         `json:"has_problems"`
	// shipping: true = ongkir masih perkiraan (alamat belum dipilih)
	ShippingEstimate bool `json:"shipping_estimate,omitempty"`

	shipping decimal.Decimal
	// promos: voucher (Total sudah dikurangi Discount)
	Discount        string       `json:"discount"`
	Voucher         *VoucherInfo `json:"voucher,omitempty"`
	VoucherError    string       `json:"voucher_error,omitempty"`
	VoucherMinSpend string       `json:"voucher_min_spend,omitempty"`

	total    decimal.Decimal // termasuk ongkir, sudah dikurangi diskon
	discount decimal.Decimal
	voucher  *voucher
}

func normalizeItems(items []CartItem) ([]CartItem, error) {
	if len(items) == 0 {
		return nil, invalid("items", "cart is empty")
	}
	if len(items) > maxLines {
		return nil, invalid("items", fmt.Sprintf("at most %d different items per order", maxLines))
	}
	merged := map[string]int{}
	order := []string{}
	for i, it := range items {
		id := strings.TrimSpace(it.VariantID)
		if len(id) != 36 {
			return nil, invalid(fmt.Sprintf("items[%d].variant_id", i), "invalid item")
		}
		if it.Qty < 1 || it.Qty > maxQty {
			return nil, invalid(fmt.Sprintf("items[%d].qty", i), fmt.Sprintf("quantity must be 1–%d", maxQty))
		}
		if _, ok := merged[id]; !ok {
			order = append(order, id)
		}
		merged[id] += it.Qty
		if merged[id] > maxQty {
			merged[id] = maxQty
		}
	}
	out := make([]CartItem, len(order))
	for i, id := range order {
		out[i] = CartItem{id, merged[id]}
	}
	return out, nil
}

type querier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

// quote menghitung keranjang. lock=true: varian dikunci (FOR UPDATE) — dipakai di dalam transaksi checkout.
func (s *Service) quote(ctx context.Context, q querier, lang string, items []CartItem, lock bool) (*Quote, error) {
	ids := make([]string, len(items))
	for i, it := range items {
		ids[i] = it.VariantID
	}
	sql := `
		SELECT v.id, p.id, v.sku, v.price, coalesce(v.compare_at_price::text, ''), v.stock_qty - v.reserved_qty,
		       (` + buyable + `) AS ok,
		       st.id, st.name, st.slug, p.category_id,
		       coalesce(t.name, ts.name), coalesce(t.slug, ts.slug),
		       coalesce((SELECT string_agg(coalesce(ov.value_i18n->>$2::text, ov.value_i18n->>p.source_lang::text), ' / ' ORDER BY o.position)
		                 FROM variant_option_values vov JOIN product_option_values ov ON ov.id = vov.option_value_id
		                 JOIN product_options o ON o.id = vov.option_id WHERE vov.variant_id = v.id), ''),
		       coalesce((SELECT coalesce(m.image_url, m.storage_key) FROM product_media m WHERE m.product_id = p.id AND m.deleted_at IS NULL
		                 AND m.kind = 'photo' ORDER BY m.sort_order LIMIT 1), '')
		FROM product_variants v JOIN products p ON p.id = v.product_id JOIN stores st ON st.id = p.store_id
		JOIN product_translations ts ON ts.product_id = p.id AND ts.lang = p.source_lang
		LEFT JOIN product_translations t ON t.product_id = p.id AND t.lang = $2::text::lang_code
		WHERE v.id = ANY($1::uuid[]) ORDER BY v.id`
	if lock {
		sql += ` FOR UPDATE OF v`
	}
	rows, err := q.Query(ctx, sql, ids, lang)
	if err != nil {
		return nil, err
	}
	found := map[string]Line{}
	storeOf := map[string]string{}
	stores := map[string]*StoreGroup{}
	for rows.Next() {
		var l Line
		var ok bool
		var sid, sname, sslug string
		if err := rows.Scan(&l.VariantID, &l.ProductID, &l.SKU, &l.price, &l.Compare, &l.Available, &ok, &sid, &sname, &sslug,
			&l.CategoryID, &l.Name, &l.Slug, &l.Options, &l.ImageURL); err != nil {
			rows.Close()
			return nil, err
		}
		if !ok {
			l.Problem = "unavailable"
		}
		if l.Available < 0 {
			l.Available = 0
		}
		l.ImageURL = s.url(l.ImageURL)
		if stores[sid] == nil {
			stores[sid] = &StoreGroup{StoreID: sid, StoreName: sname, StoreSlug: sslug}
		}
		l.Price = l.price.StringFixed(2)
		found[l.VariantID] = l
		storeOf[l.VariantID] = sid
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	// promos: harga flash sale
	if err := applyFlash(ctx, q, found, items, lock); err != nil {
		return nil, err
	}

	qt := &Quote{Discount: "0.00"}
	order := []string{}
	for _, it := range items {
		l, ok := found[it.VariantID]
		if !ok {
			// Varian sudah dihapus: tidak bisa ditampilkan (tidak tahu tokonya) → diabaikan di sisi server,
			// tapi tetap ditandai supaya keranjang di browser membersihkannya.
			qt.HasProblems = true
			continue
		}
		sid := storeOf[it.VariantID]
		l.Requested = it.Qty
		l.Qty = it.Qty
		switch {
		case l.Problem != "":
			l.Qty = 0
		case l.Available <= 0:
			l.Problem, l.Qty = "out_of_stock", 0
		case it.Qty > l.Available:
			l.Problem, l.Qty = "qty_reduced", l.Available
		}
		if l.Problem != "" {
			qt.HasProblems = true
		}
		lt := l.price.Mul(decimal.NewFromInt(int64(l.Qty)))
		l.LineTotal = lt.StringFixed(2)
		g := stores[sid]
		if g.Lines == nil {
			order = append(order, sid)
		}
		g.Lines = append(g.Lines, l)
		g.subtotal = g.subtotal.Add(lt)
		g.count += l.Qty
		qt.total = qt.total.Add(lt)
		qt.ItemCount += l.Qty
	}
	for _, sid := range order {
		g := stores[sid]
		g.Subtotal = g.subtotal.StringFixed(2)
		qt.Stores = append(qt.Stores, *g)
	}
	if qt.Stores == nil {
		qt.Stores = []StoreGroup{}
	}
	qt.Subtotal = qt.total.StringFixed(2)
	qt.ShippingFee = "0.00" // gratis ongkir di Sihanoukville
	qt.Total = qt.total.StringFixed(2)
	return qt, nil
}

func (s *Service) Quote(ctx context.Context, lang string, items []CartItem) (*Quote, error) {
	if len(items) == 0 {
		return &Quote{Stores: []StoreGroup{}, Subtotal: "0.00", ShippingFee: "0.00", Total: "0.00", Discount: "0.00"}, nil
	}
	items, err := normalizeItems(items)
	if err != nil {
		return nil, err
	}
	qt, err := s.quote(ctx, s.Pool, lang, items, false)
	if err == nil {
		err = s.applyShipping(ctx, s.Pool, qt, "") // shipping: perkiraan dari zona default
	}
	return qt, err
}

func (s *Service) url(keyOrURL string) string {
	if keyOrURL == "" || strings.HasPrefix(keyOrURL, "http") || strings.HasPrefix(keyOrURL, "/") || s.Files == nil {
		return keyOrURL
	}
	return s.Files.PublicURL(keyOrURL)
}

// ---------------------------------------------------------------------------
// Checkout
// ---------------------------------------------------------------------------

type Address struct {
	Recipient   string  `json:"recipient"`
	Phone       string  `json:"phone"`
	AddressLine string  `json:"address_line"`
	Note        string  `json:"note"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	PlaceID     string  `json:"google_place_id"`
}

type CheckoutInput struct {
	IdempotencyKey string     `json:"idempotency_key"`
	Items          []CartItem `json:"items"`
	Address        Address    `json:"address"`
	ExpectedTotal  string     `json:"expected_total"` // total yang dilihat pembeli; beda = harga berubah
	SaveAddress    bool       `json:"save_address"`
	VoucherCode    string     `json:"voucher_code"` // promos
}

// CartChangedError membawa keranjang terbaru supaya halaman bisa menampilkan apa yang berubah.
type CartChangedError struct{ Quote *Quote }

func (e *CartChangedError) Error() string { return "cart changed" }
func (e *CartChangedError) Unwrap() error { return ErrCartChanged }

var phoneRe = func(s string) bool {
	d := 0
	for _, r := range s {
		switch {
		case r >= '0' && r <= '9':
			d++
		case r == '+' || r == ' ' || r == '-':
		default:
			return false
		}
	}
	return d >= 8 && d <= 15
}

func (a *Address) normalize() error {
	a.Recipient = strings.TrimSpace(a.Recipient)
	a.Phone = strings.TrimSpace(a.Phone)
	a.AddressLine = strings.TrimSpace(a.AddressLine)
	a.Note = strings.TrimSpace(a.Note)
	switch {
	case len([]rune(a.Recipient)) < 2 || len([]rune(a.Recipient)) > 80:
		return invalid("address.recipient", "enter the recipient's name")
	case !phoneRe(a.Phone):
		return invalid("address.phone", "enter a valid phone number")
	case len([]rune(a.AddressLine)) < 5 || len([]rune(a.AddressLine)) > 300:
		return invalid("address.address_line", "choose your delivery location on the map")
	case len([]rune(a.Note)) > 300:
		return invalid("address.note", "note is too long")
	case a.Lat == 0 || a.Lng == 0 || a.Lat < -90 || a.Lat > 90 || a.Lng < -180 || a.Lng > 180:
		return invalid("address.address_line", "choose your delivery location on the map")
	}
	if len(a.PlaceID) > 255 {
		a.PlaceID = ""
	}
	return nil
}

// ServiceArea: kode wilayah layanan untuk titik ini, "" kalau di luar.
func (s *Service) ServiceArea(ctx context.Context, lat, lng float64) (string, error) {
	rows, err := s.Pool.Query(ctx, `SELECT code, boundary FROM service_areas WHERE is_active`)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	for rows.Next() {
		var code string
		var raw []byte
		if err := rows.Scan(&code, &raw); err != nil {
			return "", err
		}
		if poly, err := geo.ParsePolygon(raw); err == nil && poly.Contains(lat, lng) {
			return code, nil
		}
	}
	return "", rows.Err()
}

type CheckoutView struct {
	ID          string      `json:"id"`
	Status      string      `json:"status"`
	Provider    string      `json:"provider"`
	Total       string      `json:"total"`
	Subtotal    string      `json:"subtotal"`
	ShippingFee string      `json:"shipping_fee"`
	Discount    string      `json:"discount"` // promos
	QRString    string      `json:"qr_string,omitempty"`
	QRImage     string      `json:"qr_image,omitempty"`
	Deeplink    string      `json:"deeplink,omitempty"`
	ExpiresAt   time.Time   `json:"expires_at"`
	SecondsLeft int         `json:"seconds_left"`
	PaidAt      *time.Time  `json:"paid_at,omitempty"`
	Recipient   string      `json:"recipient"`
	Phone       string      `json:"phone"`
	AddressLine string      `json:"address_line"`
	Orders      []OrderView `json:"orders"`
}

// orderCommission: aturan paling spesifik menang — toko > kategori > global (diatur di admin → Commission rules).
// Aturan toko berlaku untuk seluruh pesanan; tanpa aturan toko, tiap barang memakai tarif kategorinya (atau global).
// Kalau tarif per barang berbeda, commission_rate pesanan = rata-rata tertimbang (komisi / subtotal).
func (s *Service) orderCommission(ctx context.Context, tx pgx.Tx, storeID string, lines []Line, subtotal decimal.Decimal) (rate, comm, seller decimal.Decimal, err error) {
	const active = `valid_from <= now() AND (valid_to IS NULL OR valid_to > now())`
	global := decimal.NewFromInt(8)
	var r decimal.Decimal
	if err = tx.QueryRow(ctx, `SELECT rate_percent FROM commission_rules WHERE scope = 'global' AND `+active+` ORDER BY valid_from DESC LIMIT 1`).Scan(&r); err == nil {
		global = r
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return
	}
	err = tx.QueryRow(ctx, `SELECT rate_percent FROM commission_rules WHERE scope = 'store' AND store_id = $1 AND `+active+` ORDER BY valid_from DESC LIMIT 1`, storeID).Scan(&r)
	if err == nil {
		comm, seller = money.SplitCommission(subtotal, r)
		return r, comm, seller, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return
	}
	cats := map[string]decimal.Decimal{}
	rows, err := tx.Query(ctx, `SELECT DISTINCT ON (category_id) category_id, rate_percent FROM commission_rules
		WHERE scope = 'category' AND `+active+` ORDER BY category_id, valid_from DESC`)
	if err != nil {
		return
	}
	for rows.Next() {
		var c string
		var v decimal.Decimal
		if err = rows.Scan(&c, &v); err != nil {
			rows.Close()
			return
		}
		cats[c] = v
	}
	rows.Close()
	// Sub-kategori tanpa tarif sendiri memakai tarif kategori utamanya.
	parents := map[string]string{}
	prow, err := tx.Query(ctx, `SELECT id, parent_id FROM categories WHERE parent_id IS NOT NULL`)
	if err != nil {
		return
	}
	for prow.Next() {
		var id, parent string
		if err = prow.Scan(&id, &parent); err != nil {
			prow.Close()
			return
		}
		parents[id] = parent
	}
	prow.Close()
	comm = decimal.Zero
	same, first := true, decimal.Decimal{}
	for i, l := range lines {
		lr := global
		if v, ok := cats[l.CategoryID]; ok {
			lr = v
		} else if v, ok := cats[parents[l.CategoryID]]; ok {
			lr = v
		}
		if i == 0 {
			first = lr
		} else if !lr.Equal(first) {
			same = false
		}
		comm = comm.Add(money.Percent(l.price.Mul(decimal.NewFromInt(int64(l.Qty))), lr))
	}
	if len(lines) == 0 {
		first = global
	}
	if same {
		comm, seller = money.SplitCommission(subtotal, first)
		return first, comm, seller, nil
	}
	rate = decimal.Zero
	if subtotal.IsPositive() {
		rate = comm.Mul(decimal.NewFromInt(100)).Div(subtotal).Round(2)
	}
	return rate, comm, subtotal.Sub(comm), nil
}

const codeAlphabet = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ" // tanpa 0/O/1/I supaya mudah dibaca di telepon

func randCode(n int) string {
	b := make([]byte, n)
	for i := range b {
		k, _ := rand.Int(rand.Reader, big.NewInt(int64(len(codeAlphabet))))
		b[i] = codeAlphabet[k.Int64()]
	}
	return string(b)
}

func (s *Service) Checkout(ctx context.Context, buyerID, lang string, in CheckoutInput) (*CheckoutView, error) {
	if !s.open(ctx) {
		return nil, ErrClosed
	}
	in.IdempotencyKey = strings.TrimSpace(in.IdempotencyKey)
	if len(in.IdempotencyKey) < 8 || len(in.IdempotencyKey) > 64 {
		return nil, invalid("idempotency_key", "missing idempotency key")
	}
	// Klik dua kali / jaringan putus lalu kirim ulang → checkout yang sama.
	var existing string
	err := s.Pool.QueryRow(ctx, `SELECT id FROM checkouts WHERE buyer_id = $1 AND idempotency_key = $2`, buyerID, in.IdempotencyKey).Scan(&existing)
	if err == nil {
		return s.ensureQR(ctx, buyerID, existing)
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	items, err := normalizeItems(in.Items)
	if err != nil {
		return nil, err
	}
	if err := in.Address.normalize(); err != nil {
		return nil, err
	}
	area, err := s.ServiceArea(ctx, in.Address.Lat, in.Address.Lng)
	if err != nil {
		return nil, err
	}
	if area == "" {
		return nil, ErrNotServed
	}

	var checkoutID string
	ctx = withFlashBuyer(ctx, buyerID) // promos: batas beli flash per pembeli
	err = db.WithActor(ctx, s.Pool, buyerID, "buyer", "", func(tx pgx.Tx) error {
		qt, err := s.quote(ctx, tx, lang, items, true)
		if err == nil {
			err = s.applyShipping(ctx, tx, qt, area) // shipping: ongkir zona alamat pembeli
		}
		if err != nil {
			return err
		}
		if qt.HasProblems || qt.ItemCount == 0 || len(qt.Stores) == 0 {
			return &CartChangedError{qt}
		}
		// promos: voucher (dicek lengkap & dikunci); expected_total = total setelah potongan
		if err := s.voucherForCheckout(ctx, tx, buyerID, in.VoucherCode, qt); err != nil {
			return err
		}
		if exp, err := decimal.NewFromString(in.ExpectedTotal); err != nil || !exp.Equal(qt.payable()) {
			return &CartChangedError{qt}
		}
		lifetime := s.lifetime(ctx)
		if lifetime < 3 {
			lifetime = 15
		}
		tranID := "P" + randCode(15)
		a := in.Address
		err = tx.QueryRow(ctx, `
			INSERT INTO checkouts (buyer_id, idempotency_key, tran_id, provider, subtotal, shipping_fee, total,
			  recipient, phone, address_line, address_note, lat, lng, google_place_id, service_area_code, lang, expires_at)
			VALUES ($1, $2, $3, $4, $5::numeric - $16::numeric, $16, $5, $6, $7, $8, nullif($9, ''), $10, $11, nullif($12, ''), $13, $14::lang_code,
			        now() + make_interval(mins => $15))
			RETURNING id`, buyerID, in.IdempotencyKey, tranID, s.Pay.Name(), qt.total, a.Recipient, a.Phone, a.AddressLine,
			a.Note, a.Lat, a.Lng, a.PlaceID, area, lang, lifetime, qt.shipping).Scan(&checkoutID)
		if err != nil {
			return err
		}
		// promos: potongan voucher di checkout + pembagian per pesanan
		if err := saveVoucher(ctx, tx, checkoutID, buyerID, qt); err != nil {
			return err
		}
		discounts := qt.allocate()
		day := time.Now().In(phnomPenh).Format("060102")
		for _, g := range qt.Stores {
			// promos: subsidi flash dari Phsar — komisi & penghasilan seller dihitung dari harga flash + subsidi
			grossLines, subsidy := flashGross(g.Lines)
			rate, comm, sellerAmt, err := s.orderCommission(ctx, tx, g.StoreID, grossLines, g.subtotal.Add(subsidy))
			if err != nil {
				return err
			}
			var orderID string
			err = tx.QueryRow(ctx, `
				INSERT INTO orders (order_no, checkout_id, buyer_id, store_id, subtotal, shipping_fee, total, commission_rate,
				  commission_amount, seller_amount, item_count, flash_subsidy)
				VALUES ($1, $2, $3, $4, $5, $10, $5::numeric + $10::numeric, $6, $7, $8, $9, $11) RETURNING id`,
				"PH"+day+"-"+randCode(5), checkoutID, buyerID, g.StoreID, g.subtotal, rate, comm, sellerAmt, g.count, g.shipping, subsidy).Scan(&orderID)
			if err != nil {
				return err
			}
			if err := orderDiscount(ctx, tx, orderID, discounts[g.StoreID]); err != nil { // promos
				return err
			}
			for _, l := range g.Lines {
				var itemID string
				err = tx.QueryRow(ctx, `
					INSERT INTO order_items (order_id, product_id, variant_id, name, options, sku, image_url, unit_price, qty, line_total)
					VALUES ($1, $2, $3, $4, $5, $6, nullif($7, ''), $8::numeric, $9::int, $8::numeric * $9::int) RETURNING id`,
					orderID, l.ProductID, l.VariantID, truncRunes(l.Name, 160), truncRunes(l.Options, 200), l.SKU, l.ImageURL, l.price, l.Qty).Scan(&itemID)
				if err != nil {
					return err
				}
				if err := flashItem(ctx, tx, itemID, l); err != nil { // promos
					return err
				}
				// Tahan stok. Kondisi WHERE menjaga dua pembeli tidak bisa mengambil barang terakhir yang sama.
				tag, err := tx.Exec(ctx, `UPDATE product_variants SET reserved_qty = reserved_qty + $2
					WHERE id = $1 AND stock_qty - reserved_qty >= $2`, l.VariantID, l.Qty)
				if err != nil {
					return err
				}
				if tag.RowsAffected() != 1 {
					return &CartChangedError{qt}
				}
				if _, err := tx.Exec(ctx, `INSERT INTO inventory_movements (variant_id, delta, reason, order_item_id, created_by)
					VALUES ($1, $2, 'reserve', $3, $4)`, l.VariantID, -l.Qty, itemID, buyerID); err != nil {
					return err
				}
			}
		}
		if in.SaveAddress {
			if err := saveAddress(ctx, tx, buyerID, a, area); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		var ce *CartChangedError
		if errors.As(err, &ce) {
			return nil, err
		}
		// Dua request bersamaan dengan kunci idempotensi yang sama: ambil yang menang.
		if strings.Contains(err.Error(), "checkouts_buyer_id_idempotency_key_key") {
			if e2 := s.Pool.QueryRow(ctx, `SELECT id FROM checkouts WHERE buyer_id = $1 AND idempotency_key = $2`, buyerID, in.IdempotencyKey).Scan(&existing); e2 == nil {
				return s.ensureQR(ctx, buyerID, existing)
			}
		}
		return nil, err
	}
	return s.ensureQR(ctx, buyerID, checkoutID)
}

var phnomPenh = func() *time.Location {
	l, err := time.LoadLocation("Asia/Phnom_Penh")
	if err != nil {
		return time.FixedZone("ICT", 7*3600)
	}
	return l
}()

func truncRunes(s string, n int) string {
	if r := []rune(s); len(r) > n {
		return string(r[:n])
	}
	return s
}

func saveAddress(ctx context.Context, tx pgx.Tx, buyerID string, a Address, area string) error {
	var id string
	err := tx.QueryRow(ctx, `SELECT id FROM addresses WHERE user_id = $1 AND is_default AND deleted_at IS NULL`, buyerID).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		_, err = tx.Exec(ctx, `INSERT INTO addresses (user_id, recipient, phone, address_line, service_area_code, lat, lng, note, google_place_id, is_default)
			VALUES ($1, $2, $3, $4, $5, $6, $7, nullif($8, ''), nullif($9, ''), true)`, buyerID, a.Recipient, a.Phone, a.AddressLine, area, a.Lat, a.Lng, a.Note, a.PlaceID)
		return err
	}
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `UPDATE addresses SET recipient = $2, phone = $3, address_line = $4, service_area_code = $5, lat = $6, lng = $7,
		note = nullif($8, ''), google_place_id = nullif($9, '') WHERE id = $1`, id, a.Recipient, a.Phone, a.AddressLine, area, a.Lat, a.Lng, a.Note, a.PlaceID)
	return err
}

// DefaultAddress: alamat terakhir pembeli (untuk mengisi form checkout). nil kalau belum ada.
func (s *Service) DefaultAddress(ctx context.Context, buyerID string) (*Address, error) {
	var a Address
	var note, place *string
	err := s.Pool.QueryRow(ctx, `SELECT recipient, phone, address_line, note, lat::float8, lng::float8, google_place_id
		FROM addresses WHERE user_id = $1 AND deleted_at IS NULL ORDER BY is_default DESC, updated_at DESC LIMIT 1`, buyerID).
		Scan(&a.Recipient, &a.Phone, &a.AddressLine, &note, &a.Lat, &a.Lng, &place)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if note != nil {
		a.Note = *note
	}
	if place != nil {
		a.PlaceID = *place
	}
	return &a, nil
}

// ensureQR: buat QR di penyedia pembayaran kalau belum ada (di luar transaksi: panggilan jaringan).
func (s *Service) ensureQR(ctx context.Context, buyerID, checkoutID string) (*CheckoutView, error) {
	v, err := s.GetCheckout(ctx, buyerID, checkoutID)
	if err != nil {
		return nil, err
	}
	if v.Status != "pending" || v.QRString != "" {
		return v, nil
	}
	var tranID string
	if err := s.Pool.QueryRow(ctx, `SELECT tran_id FROM checkouts WHERE id = $1`, checkoutID).Scan(&tranID); err != nil {
		return nil, err
	}
	items := []payment.Item{}
	for _, o := range v.Orders {
		for _, it := range o.Items {
			if len(items) < 50 {
				items = append(items, payment.Item{Name: truncRunes(it.Name, 60), Quantity: it.Qty, Price: it.UnitPrice})
			}
		}
	}
	if fee, err := decimal.NewFromString(v.ShippingFee); err == nil && fee.IsPositive() { // shipping: ongkir sebagai baris sendiri
		items = append(items, payment.Item{Name: "Delivery", Quantity: 1, Price: fee.StringFixed(2)})
	}
	left := int(time.Until(v.ExpiresAt).Minutes())
	if left < 3 {
		left = 3
	}
	qr, err := s.Pay.CreateQR(ctx, payment.QRRequest{TranID: tranID, Amount: v.Total, FirstName: truncRunes(v.Recipient, 40), Phone: v.Phone,
		Items: items, CallbackURL: s.BaseURL + "/webhooks/payway", LifetimeMin: left})
	if err != nil {
		s.Log.Error("pembuatan QR gagal", "checkout", checkoutID, "err", err)
		s.event(ctx, checkoutID, tranID, "create", false, "ERROR", map[string]any{"error": err.Error()})
		// Pembeli tidak bisa membayar: batalkan sekarang supaya stok kembali.
		_ = s.cancel(ctx, checkoutID, "cancelled", "payment_unavailable")
		return nil, fmt.Errorf("%w: %v", ErrPayment, err)
	}
	if _, err := s.Pool.Exec(ctx, `UPDATE checkouts SET qr_string = $2, qr_image = nullif($3, ''), deeplink = nullif($4, '') WHERE id = $1`,
		checkoutID, qr.String, qr.Image, qr.Deeplink); err != nil {
		return nil, err
	}
	s.event(ctx, checkoutID, tranID, "create", true, "CREATED", map[string]any{"has_image": qr.Image != "", "has_deeplink": qr.Deeplink != ""})
	return s.GetCheckout(ctx, buyerID, checkoutID)
}

func (s *Service) event(ctx context.Context, checkoutID, tranID, source string, verified bool, status string, payload any) {
	var cid any
	if checkoutID != "" {
		cid = checkoutID
	}
	if _, err := s.Pool.Exec(ctx, `INSERT INTO payment_events (checkout_id, tran_id, source, verified, status, payload) VALUES ($1, nullif($2, ''), $3, $4, $5, $6)`,
		cid, tranID, source, verified, status, payload); err != nil {
		s.Log.Error("payment_events", "err", err)
	}
}

// ---------------------------------------------------------------------------
// Lihat checkout & pesanan
// ---------------------------------------------------------------------------

type OrderItemView struct {
	ProductID string `json:"product_id"`
	Name      string `json:"name"`
	Options   string `json:"options"`
	ImageURL  string `json:"image_url"`
	UnitPrice string `json:"unit_price"`
	Qty       int    `json:"qty"`
	LineTotal string `json:"line_total"`
}

type OrderView struct {
	OrderNo      string          `json:"order_no"`
	Status       string          `json:"status"`
	StoreName    string          `json:"store_name"`
	StoreSlug    string          `json:"store_slug"`
	Subtotal     string          `json:"subtotal"`
	ShippingFee  string          `json:"shipping_fee"`
	Discount     string          `json:"discount"` // promos
	Total        string          `json:"total"`
	ItemCount    int             `json:"item_count"`
	CancelReason string          `json:"cancel_reason,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	PaidAt       *time.Time      `json:"paid_at,omitempty"`
	Items        []OrderItemView `json:"items"`
	CheckoutID   string          `json:"checkout_id"`
	Recipient    string          `json:"recipient,omitempty"`
	Phone        string          `json:"phone,omitempty"`
	AddressLine  string          `json:"address_line,omitempty"`
	AddressNote  string          `json:"address_note,omitempty"`
	PaymentType  string          `json:"payment_type,omitempty"`
	PaymentRef   string          `json:"payment_ref,omitempty"`
	RefundStatus string          `json:"refund_status"`
	RefundAmount string          `json:"refund_amount,omitempty"`
	RefundedAt   *time.Time      `json:"refunded_at,omitempty"`
	Events       []OrderEvent    `json:"events,omitempty"`
	id           string
}

func (s *Service) GetCheckout(ctx context.Context, buyerID, id string) (*CheckoutView, error) {
	if len(id) != 36 {
		return nil, ErrNotFound
	}
	v := &CheckoutView{}
	var qrs, qri, dl *string
	err := s.Pool.QueryRow(ctx, `
		SELECT id, status, provider, total::text, subtotal::text, shipping_fee::text, qr_string, qr_image, deeplink, expires_at, paid_at,
		       recipient, phone, address_line, discount::text
		FROM checkouts WHERE id = $1 AND buyer_id = $2`, id, buyerID).
		Scan(&v.ID, &v.Status, &v.Provider, &v.Total, &v.Subtotal, &v.ShippingFee, &qrs, &qri, &dl, &v.ExpiresAt, &v.PaidAt,
			&v.Recipient, &v.Phone, &v.AddressLine, &v.Discount)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if v.Status == "pending" {
		v.QRString, v.QRImage, v.Deeplink = deref(qrs), deref(qri), deref(dl)
		if left := int(time.Until(v.ExpiresAt).Seconds()); left > 0 {
			v.SecondsLeft = left
		}
	}
	v.Orders, err = s.orders(ctx, `o.checkout_id = $1 AND o.buyer_id = $2`, 100, id, buyerID)
	return v, err
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func (s *Service) orders(ctx context.Context, where string, limit int, args ...any) ([]OrderView, error) {
	rows, err := s.Pool.Query(ctx, `
		SELECT o.id, o.order_no, o.status::text, st.name, st.slug, o.subtotal::text, o.shipping_fee::text, o.total::text, o.item_count,
		       coalesce(o.cancel_reason, ''), o.created_at, o.paid_at, o.checkout_id,
		       c.recipient, c.phone, c.address_line, coalesce(c.address_note, ''), coalesce(c.payment_type, ''), coalesce(c.payment_ref, ''),
		       o.refund_status, coalesce(o.refund_amount::text, ''), o.refunded_at, o.discount::text
		FROM orders o JOIN stores st ON st.id = o.store_id JOIN checkouts c ON c.id = o.checkout_id
		WHERE `+where+` ORDER BY o.created_at DESC, o.order_no LIMIT `+fmt.Sprint(limit), args...)
	if err != nil {
		return nil, err
	}
	out := []OrderView{}
	ids := []string{}
	for rows.Next() {
		var o OrderView
		var id string
		if err := rows.Scan(&id, &o.OrderNo, &o.Status, &o.StoreName, &o.StoreSlug, &o.Subtotal, &o.ShippingFee, &o.Total, &o.ItemCount,
			&o.CancelReason, &o.CreatedAt, &o.PaidAt, &o.CheckoutID, &o.Recipient, &o.Phone, &o.AddressLine, &o.AddressNote,
			&o.PaymentType, &o.PaymentRef, &o.RefundStatus, &o.RefundAmount, &o.RefundedAt, &o.Discount); err != nil {
			rows.Close()
			return nil, err
		}
		o.Items = []OrderItemView{}
		o.id = id
		out = append(out, o)
		ids = append(ids, id)
	}
	rows.Close()
	if len(ids) == 0 {
		return out, rows.Err()
	}
	idx := map[string]int{}
	for i, id := range ids {
		idx[id] = i
	}
	rows, err = s.Pool.Query(ctx, `SELECT order_id, product_id, name, options, coalesce(image_url, ''), unit_price::text, qty, line_total::text
		FROM order_items WHERE order_id = ANY($1::uuid[]) ORDER BY name, options`, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var oid string
		var it OrderItemView
		if err := rows.Scan(&oid, &it.ProductID, &it.Name, &it.Options, &it.ImageURL, &it.UnitPrice, &it.Qty, &it.LineTotal); err != nil {
			return nil, err
		}
		out[idx[oid]].Items = append(out[idx[oid]].Items, it)
	}
	return out, rows.Err()
}

// ListOrders: pesanan pembeli (yang belum dibayar & sudah kedaluwarsa tidak ditampilkan supaya tidak membingungkan).
func (s *Service) ListOrders(ctx context.Context, buyerID string) ([]OrderView, error) {
	return s.orders(ctx, `o.buyer_id = $1 AND NOT (o.status = 'cancelled' AND o.cancel_reason IN ('payment_expired','payment_unavailable','buyer_unpaid'))`, 100, buyerID)
}

func (s *Service) GetOrder(ctx context.Context, buyerID, orderNo string) (*OrderView, error) {
	list, err := s.orders(ctx, `o.buyer_id = $1 AND o.order_no = $2`, 1, buyerID, strings.ToUpper(strings.TrimSpace(orderNo)))
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, ErrNotFound
	}
	o := &list[0]
	o.Events, err = s.events(ctx, o.id)
	return o, err
}

// ---------------------------------------------------------------------------
// Pembayaran masuk / kedaluwarsa / batal
// ---------------------------------------------------------------------------

// Refresh: tanya status resmi ke penyedia pembayaran, lalu proses (lunas → pesanan masuk ke seller).
// Dipanggil oleh callback, halaman pembayaran (polling), dan worker.
func (s *Service) Refresh(ctx context.Context, tranID, source string) (string, error) {
	var id, status string
	var total decimal.Decimal
	err := s.Pool.QueryRow(ctx, `SELECT id, status, total FROM checkouts WHERE tran_id = $1`, tranID).Scan(&id, &status, &total)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", err
	}
	if status == "paid" {
		return status, nil
	}
	st, err := s.Pay.Check(ctx, tranID)
	if err != nil {
		return status, err
	}
	if !st.Paid() {
		if st.State != "PENDING" && source != "poll" {
			s.event(ctx, id, tranID, "check", true, st.State, st.Raw)
		}
		return status, nil
	}
	paid, err := decimal.NewFromString(st.Amount)
	if err != nil || !paid.Equal(total) {
		// Nominal tidak cocok: jangan pernah tandai lunas. Admin perlu cek manual.
		s.event(ctx, id, tranID, "check", true, "AMOUNT_MISMATCH", map[string]any{"paid": st.Amount, "expected": total.StringFixed(2), "raw": st.Raw})
		s.Log.Error("nominal pembayaran tidak cocok", "tran_id", tranID, "paid", st.Amount, "expected", total.StringFixed(2))
		return status, nil
	}
	s.event(ctx, id, tranID, "check", true, "APPROVED", map[string]any{"amount": st.Amount, "apv": st.APV, "raw": st.Raw})
	return "paid", s.markPaid(ctx, id, paid, st)
}

func (s *Service) markPaid(ctx context.Context, checkoutID string, paid decimal.Decimal, st payment.Status) error {
	return db.WithActor(ctx, s.Pool, "", "system", "", func(tx pgx.Tx) error {
		var status string
		if err := tx.QueryRow(ctx, `SELECT status FROM checkouts WHERE id = $1 FOR UPDATE`, checkoutID).Scan(&status); err != nil {
			return err
		}
		if status == "paid" {
			return nil // sudah diproses (callback & polling bisa datang bersamaan)
		}
		ptype := st.PaymentType
		if ptype == "" {
			ptype = "KHQR"
		}
		if _, err := tx.Exec(ctx, `UPDATE checkouts SET status = 'paid', paid_at = now(), paid_amount = $2, payment_type = $3,
			payment_ref = nullif($4, '') WHERE id = $1`, checkoutID, paid, ptype, st.APV); err != nil {
			return err
		}
		if err := redemptionStatus(ctx, tx, checkoutID, "used"); err != nil { // promos
			return err
		}
		rows, err := tx.Query(ctx, `SELECT o.id, o.store_id, o.order_no, o.status::text, o.total::text FROM orders o WHERE o.checkout_id = $1 FOR UPDATE`, checkoutID)
		if err != nil {
			return err
		}
		type ord struct{ id, store, no, status, total string }
		var list []ord
		for rows.Next() {
			var o ord
			if err := rows.Scan(&o.id, &o.store, &o.no, &o.status, &o.total); err != nil {
				rows.Close()
				return err
			}
			list = append(list, o)
		}
		rows.Close()
		for _, o := range list {
			if o.status == "pending_payment" {
				if err := s.sell(ctx, tx, o.id, true); err != nil {
					return err
				}
			} else if o.status == "cancelled" {
				// Dibayar setelah kedaluwarsa (stok sudah dilepas): coba ambil stok lagi.
				if err := s.sell(ctx, tx, o.id, false); err != nil {
					if errors.Is(err, errNoStock) {
						// Barang sudah dibeli orang lain → pesanan tetap batal, uang harus dikembalikan (admin).
						if _, err := tx.Exec(ctx, `UPDATE orders SET cancel_reason = 'paid_after_expiry_refund' WHERE id = $1`, o.id); err != nil {
							return err
						}
						s.Log.Error("dibayar setelah kedaluwarsa & stok habis: perlu refund", "order", o.no)
						continue
					}
					return err
				}
			} else {
				continue
			}
			tm, err := hours.LoadTiming(ctx, tx)
			if err != nil {
				return err
			}
			respondBy, err := hours.StoreDeadline(ctx, tx, tm, o.store, time.Now(), tm.AcceptMinutes)
			if err != nil {
				return err
			}
			if _, err := tx.Exec(ctx, `UPDATE orders SET status = 'paid', paid_at = now(), respond_by = $2,
				cancelled_at = NULL, cancel_reason = NULL WHERE id = $1`, o.id, respondBy); err != nil {
				return err
			}
			if err := addEvent(ctx, tx, o.id, "paid", "system", "", ""); err != nil {
				return err
			}
			if err := notify.Enqueue(ctx, tx, o.store, notify.OrderNew, notify.Params{Order: o.no, Note: "$" + o.total}, "#/store/"+o.store+"/orders"); err != nil {
				return err
			}
		}
		return nil
	})
}

var errNoStock = errors.New("no stock")

// sell: stok terjual. reserved=true: stok sudah ditahan (normal). false: ambil lagi setelah dilepas.
func (s *Service) sell(ctx context.Context, tx pgx.Tx, orderID string, reserved bool) error {
	rows, err := tx.Query(ctx, `SELECT id, variant_id, qty FROM order_items WHERE order_id = $1 ORDER BY variant_id`, orderID)
	if err != nil {
		return err
	}
	type it struct {
		id, variant string
		qty         int
	}
	var items []it
	for rows.Next() {
		var x it
		if err := rows.Scan(&x.id, &x.variant, &x.qty); err != nil {
			rows.Close()
			return err
		}
		items = append(items, x)
	}
	rows.Close()
	for _, x := range items {
		q := `UPDATE product_variants SET stock_qty = stock_qty - $2, reserved_qty = reserved_qty - $2 WHERE id = $1 AND reserved_qty >= $2`
		if !reserved {
			q = `UPDATE product_variants SET stock_qty = stock_qty - $2 WHERE id = $1 AND stock_qty - reserved_qty >= $2`
		}
		tag, err := tx.Exec(ctx, q, x.variant, x.qty)
		if err != nil {
			return err
		}
		if tag.RowsAffected() != 1 {
			return errNoStock
		}
		if _, err := tx.Exec(ctx, `INSERT INTO inventory_movements (variant_id, delta, reason, order_item_id) VALUES ($1, $2, 'sale', $3)`,
			x.variant, -x.qty, x.id); err != nil {
			return err
		}
	}
	// Stok habis: produk tetap "active" (seller cukup menambah stok); website menampilkan label "Habis".
	return nil
}

// cancel: checkout belum dibayar → batal/kedaluwarsa, pesanan batal, stok dilepas.
func (s *Service) cancel(ctx context.Context, checkoutID, newStatus, reason string) error {
	return db.WithActor(ctx, s.Pool, "", "system", "", func(tx pgx.Tx) error {
		var status string
		if err := tx.QueryRow(ctx, `SELECT status FROM checkouts WHERE id = $1 FOR UPDATE`, checkoutID).Scan(&status); err != nil {
			return err
		}
		if status != "pending" {
			return nil
		}
		if _, err := tx.Exec(ctx, `UPDATE checkouts SET status = $2 WHERE id = $1`, checkoutID, newStatus); err != nil {
			return err
		}
		if err := redemptionStatus(ctx, tx, checkoutID, "released"); err != nil { // promos: kuota voucher kembali
			return err
		}
		rows, err := tx.Query(ctx, `SELECT oi.id, oi.variant_id, oi.qty FROM order_items oi JOIN orders o ON o.id = oi.order_id
			WHERE o.checkout_id = $1 AND o.status = 'pending_payment' ORDER BY oi.variant_id`, checkoutID)
		if err != nil {
			return err
		}
		type it struct {
			id, variant string
			qty         int
		}
		var items []it
		for rows.Next() {
			var x it
			if err := rows.Scan(&x.id, &x.variant, &x.qty); err != nil {
				rows.Close()
				return err
			}
			items = append(items, x)
		}
		rows.Close()
		for _, x := range items {
			if _, err := tx.Exec(ctx, `UPDATE product_variants SET reserved_qty = greatest(reserved_qty - $2, 0) WHERE id = $1`, x.variant, x.qty); err != nil {
				return err
			}
			if _, err := tx.Exec(ctx, `INSERT INTO inventory_movements (variant_id, delta, reason, order_item_id) VALUES ($1, $2, 'release', $3)`,
				x.variant, x.qty, x.id); err != nil {
				return err
			}
		}
		_, err = tx.Exec(ctx, `UPDATE orders SET status = 'cancelled', cancelled_at = now(), cancel_reason = $2
			WHERE checkout_id = $1 AND status = 'pending_payment'`, checkoutID, reason)
		return err
	})
}

// CancelByBuyer: pembeli menekan "Batal" di layar QR. Dicek dulu ke PayWay: kalau ternyata sudah dibayar, jadi lunas.
func (s *Service) CancelByBuyer(ctx context.Context, buyerID, checkoutID string) (*CheckoutView, error) {
	v, err := s.GetCheckout(ctx, buyerID, checkoutID)
	if err != nil {
		return nil, err
	}
	if v.Status != "pending" {
		return v, nil
	}
	var tranID string
	if err := s.Pool.QueryRow(ctx, `SELECT tran_id FROM checkouts WHERE id = $1`, checkoutID).Scan(&tranID); err != nil {
		return nil, err
	}
	st, err := s.Refresh(ctx, tranID, "cancel")
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrPayment, err)
	}
	if st == "paid" {
		return s.GetCheckout(ctx, buyerID, checkoutID)
	}
	if err := s.cancel(ctx, checkoutID, "cancelled", "buyer_unpaid"); err != nil {
		return nil, err
	}
	return s.GetCheckout(ctx, buyerID, checkoutID)
}

// Poll: halaman pembayaran menanyakan status. Ke PayWay paling sering sekali per 4 detik per checkout.
func (s *Service) Poll(ctx context.Context, buyerID, checkoutID string) (*CheckoutView, error) {
	v, err := s.GetCheckout(ctx, buyerID, checkoutID)
	if err != nil || v.Status != "pending" {
		return v, err
	}
	if !pollDue(checkoutID) {
		return v, nil
	}
	var tranID string
	if err := s.Pool.QueryRow(ctx, `SELECT tran_id FROM checkouts WHERE id = $1`, checkoutID).Scan(&tranID); err != nil {
		return nil, err
	}
	if _, err := s.Refresh(ctx, tranID, "poll"); err != nil {
		s.Log.Warn("cek pembayaran gagal", "checkout", checkoutID, "err", err)
		return v, nil
	}
	return s.GetCheckout(ctx, buyerID, checkoutID)
}

// ExpireDue: worker — checkout lewat batas waktu: cek terakhir ke PayWay, lalu batalkan.
func (s *Service) ExpireDue(ctx context.Context) (int, error) {
	rows, err := s.Pool.Query(ctx, `SELECT id, tran_id FROM checkouts WHERE status = 'pending' AND expires_at < now() - interval '30 seconds'
		ORDER BY expires_at LIMIT 50`)
	if err != nil {
		return 0, err
	}
	type c struct{ id, tran string }
	var list []c
	for rows.Next() {
		var x c
		if err := rows.Scan(&x.id, &x.tran); err != nil {
			rows.Close()
			return 0, err
		}
		list = append(list, x)
	}
	rows.Close()
	n := 0
	for _, x := range list {
		st, err := s.Refresh(ctx, x.tran, "expire")
		if err != nil {
			s.Log.Warn("cek pembayaran sebelum kedaluwarsa gagal (dicoba lagi nanti)", "checkout", x.id, "err", err)
			continue
		}
		if st == "paid" {
			continue
		}
		if err := s.cancel(ctx, x.id, "expired", "payment_expired"); err != nil {
			return n, err
		}
		n++
	}
	return n, nil
}

// Worker: kedaluwarsa + cek ulang pembayaran yang tertunda.
func (s *Service) Worker(ctx context.Context) {
	t := time.NewTicker(20 * time.Second)
	defer t.Stop()
	for {
		if n, err := jobs.Track(ctx, s.Pool, "checkout.expire", s.ExpireDue); err != nil && ctx.Err() == nil { // jobs:
			s.Log.Error("worker checkout", "err", err)
		} else if n > 0 {
			s.Log.Info("checkout kedaluwarsa dibatalkan", "jumlah", n)
		}
		if n, err := jobs.Track(ctx, s.Pool, "order.accept_timeout", s.ExpireUnaccepted); err != nil && ctx.Err() == nil { // jobs:
			s.Log.Error("worker pesanan (batas terima seller)", "err", err)
		} else if n > 0 {
			s.Log.Info("pesanan tidak diterima seller tepat waktu → batal + refund", "jumlah", n)
		}
		if n, err := jobs.Track(ctx, s.Pool, "order.pack_timeout", s.ExpireUnpacked); err != nil && ctx.Err() == nil { // jobs:
			s.Log.Error("worker pesanan (batas kemas)", "err", err)
		} else if n > 0 {
			s.Log.Info("pesanan tidak dikemas tepat waktu → batal + refund", "jumlah", n)
		}
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}
	}
}

// pollDue: batas pengecekan ke PayWay per checkout (di memori).
var pollSeen = struct {
	sync.Mutex
	m map[string]time.Time
}{m: map[string]time.Time{}}

func pollDue(id string) bool {
	pollSeen.Lock()
	defer pollSeen.Unlock()
	now := time.Now()
	if t, ok := pollSeen.m[id]; ok && now.Sub(t) < 4*time.Second {
		return false
	}
	if len(pollSeen.m) > 10000 {
		keys := make([]string, 0, len(pollSeen.m))
		for k := range pollSeen.m {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys[:5000] {
			delete(pollSeen.m, k)
		}
	}
	pollSeen.m[id] = now
	return true
}
