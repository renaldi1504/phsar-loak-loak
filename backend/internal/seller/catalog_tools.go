package seller

// Alat katalog untuk toko dengan banyak SKU (migrasi 00052):
//   - SKU otomatis rapi: <kode toko>-<nomor urut produk 5 digit>[-<nilai varian>], mis. SKV-00012-RED-M.
//   - Tabel per varian (SKU, varian, harga, stok) dengan filter kategori & stok, ubah cepat per baris.
//   - Aksi massal: arsipkan, aktifkan, hapus draf, ubah stok, ubah harga %.
//   - Ekspor CSV semua SKU dan impor CSV (pratinjau dulu, lalu terapkan) untuk harga/stok.
// Semua perubahan harga/stok lewat editVariants: validasinya sama dengan formulir produk, tercatat di audit
// (db.WithActor) dan inventory_movements, dan produk tayang yang harganya turun jauh / di bawah harga wajar
// dikirim ke antrean moderasi (aturan drop_pct dll. di modrules), sama seperti kalau diubah lewat formulir.

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"

	"phsar/internal/modrules"
	"phsar/internal/money"
	"phsar/internal/paging"
	"phsar/internal/platform/db"
)

var skuPrefixRe = regexp.MustCompile(`^[A-Z0-9]{2,5}$`)

// skuWords: kata latin/angka dari nama (huruf besar). Huruf non-ASCII dibuang.
func skuWords(s string) []string {
	var words []string
	var cur strings.Builder
	for _, r := range strings.ToUpper(s) {
		if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			cur.WriteRune(r)
			continue
		}
		if cur.Len() > 0 {
			words = append(words, cur.String())
			cur.Reset()
		}
	}
	if cur.Len() > 0 {
		words = append(words, cur.String())
	}
	return words
}

// SKUPrefix: kode SKU bawaan dari nama toko — sama dengan fungsi SQL phsar_sku_prefix (migrasi 00052).
//
//	"Beauty" → BEA · "Skin Vibe" → SKV · "Phnom Penh Phone Shop" → PPP · "ហាងស្អាត" → S + 3 karakter id toko
func SKUPrefix(name, storeID string) string {
	w := skuWords(name)
	var p string
	switch len(w) {
	case 0:
	case 1:
		p = cut(w[0], 3)
	case 2:
		p = cut(w[0], 2) + cut(w[1], 1)
	default:
		p = cut(w[0], 1) + cut(w[1], 1) + cut(w[2], 1)
	}
	if len(p) < 2 {
		p = "S" + strings.ToUpper(cut(strings.ReplaceAll(storeID, "-", ""), 3))
	}
	return p
}

func cut(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}
	return s
}

// SKUBase: SKV + 12 → "SKV-00012".
func SKUBase(prefix string, no int) string {
	return fmt.Sprintf("%s-%05d", prefix, no)
}

// skuSuffix: nilai varian → "RED-M". Huruf besar ASCII saja (Khmer/Cina dibuang), maks. 10 karakter per nilai.
func skuSuffix(values []string) string {
	parts := []string{}
	for _, v := range values {
		if p := cut(strings.Join(skuWords(v), ""), 10); p != "" {
			parts = append(parts, p)
		}
	}
	return cut(strings.Join(parts, "-"), 32)
}

// autoSKUs: SKU akhir untuk tiap varian. SKU yang diisi seller (termasuk SKU lama) dipakai apa adanya; yang kosong
// diberi SKU otomatis. Nomor produk diambil sekali (products.sku_no) dari penghitung per toko — aman dipakai
// bersamaan karena INSERT … ON CONFLICT DO UPDATE mengunci baris penghitung. Akhiran kosong (nilai varian tanpa
// huruf latin) atau sudah dipakai → ditambah -1, -2, …
func (s *Service) autoSKUs(ctx context.Context, tx pgx.Tx, storeID, productID string, variants []normVariant) ([]string, error) {
	out := make([]string, len(variants))
	used := map[string]bool{}
	need := false
	for i, v := range variants {
		if v.sku != "" {
			out[i] = v.sku
			used[strings.ToLower(v.sku)] = true
		} else {
			need = true
		}
	}
	if !need {
		return out, nil
	}
	var prefix, name string
	var no *int
	if err := tx.QueryRow(ctx, `SELECT coalesce(st.sku_prefix, ''), st.name, p.sku_no FROM products p JOIN stores st ON st.id = p.store_id
		WHERE p.id = $1 FOR UPDATE OF p`, productID).Scan(&prefix, &name, &no); err != nil {
		return nil, err
	}
	if prefix == "" {
		prefix = SKUPrefix(name, storeID)
	}
	if no == nil {
		var n int
		if err := tx.QueryRow(ctx, `INSERT INTO store_sku_counters (store_id, last_no) VALUES ($1, 1)
			ON CONFLICT (store_id) DO UPDATE SET last_no = store_sku_counters.last_no + 1 RETURNING last_no`, storeID).Scan(&n); err != nil {
			return nil, err
		}
		if _, err := tx.Exec(ctx, `UPDATE products SET sku_no = $2 WHERE id = $1`, productID, n); err != nil {
			return nil, err
		}
		no = &n
	}
	base := SKUBase(prefix, *no)
	taken := func(c string) (bool, error) {
		if used[strings.ToLower(c)] {
			return true, nil
		}
		var t bool
		err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM product_variants WHERE store_id = $1 AND lower(sku) = lower($2) AND deleted_at IS NULL)`,
			storeID, c).Scan(&t)
		return t, err
	}
	multi := len(variants) > 1
	for i, v := range variants {
		if out[i] != "" {
			continue
		}
		cand := base
		suf := skuSuffix(v.values)
		if suf != "" {
			cand += "-" + suf
		}
		c := cand
		if suf == "" && multi {
			c = ""
		}
		for n := 1; ; n++ {
			if c != "" {
				t, err := taken(c)
				if err != nil {
					return nil, err
				}
				if !t {
					break
				}
			}
			c = fmt.Sprintf("%s-%d", cand, n)
		}
		used[strings.ToLower(c)] = true
		out[i] = c
	}
	return out, nil
}

// ---- Filter katalog -------------------------------------------------------------

// categoryRe: id kategori (lihat tabel categories).
var categoryRe = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// effThreshold: ambang stok menipis yang berlaku untuk varian alias v.
const effThreshold = `coalesce(v.low_stock_threshold, (SELECT s.low_stock_threshold FROM stores s WHERE s.id = v.store_id))`

// stockCond: kondisi stok untuk satu varian (alias v). out = habis · low = sisa 1..ambang.
func stockCond(stock string) string {
	switch stock {
	case "out":
		return `(v.stock_qty - v.reserved_qty) <= 0`
	case "low":
		return `(v.stock_qty - v.reserved_qty) > 0 AND (v.stock_qty - v.reserved_qty) <= ` + effThreshold
	}
	return ""
}

// catalogWhere: filter kategori (termasuk sub-kategori) dan stok. productLevel = untuk daftar produk (stok dicek
// lewat varian aktifnya), selain itu untuk tabel varian (alias v).
func catalogWhere(q ProductQuery, args []any, productLevel bool) (string, []any) {
	where := ""
	if c := strings.TrimSpace(q.Category); c != "" && categoryRe.MatchString(c) && len(c) <= 60 {
		args = append(args, c)
		n := len(args)
		where += fmt.Sprintf(` AND (p.category_id = $%d OR p.category_id IN (SELECT id FROM categories WHERE parent_id = $%d))`, n, n)
	}
	if cond := stockCond(q.Stock); cond != "" {
		if productLevel {
			where += ` AND EXISTS (SELECT 1 FROM product_variants v WHERE v.product_id = p.id AND v.deleted_at IS NULL AND v.status = 'active' AND ` + cond + `)`
		} else {
			where += ` AND v.status = 'active' AND ` + cond
		}
	}
	return where, args
}

// ---- Tabel varian ---------------------------------------------------------------

// VariantRow: satu baris tabel SKU.
type VariantRow struct {
	ID               string  `json:"id"`
	ProductID        string  `json:"product_id"`
	ProductName      string  `json:"product_name"`
	CoverURL         string  `json:"cover_url"`
	SKU              string  `json:"sku"`
	Label            string  `json:"label"` // nilai opsi, mis. "Merah / M"; kosong untuk produk tanpa varian
	Price            string  `json:"price"`
	CompareAtPrice   *string `json:"compare_at_price"`
	StockQty         int     `json:"stock_qty"`
	ReservedQty      int     `json:"reserved_qty"`
	Available        int     `json:"available"`
	LowStock         *int    `json:"low_stock_threshold"` // ambang khusus varian ini (null = ikut toko)
	Threshold        int     `json:"threshold"`           // ambang yang berlaku
	Status           string  `json:"status"`              // active · inactive (varian)
	ProductStatus    string  `json:"product_status"`
	ModerationStatus string  `json:"moderation_status"`
}

type VariantList struct {
	Variants []VariantRow   `json:"variants"`
	Page     *paging.Meta   `json:"page"`
	Counts   map[string]int `json:"counts"` // jumlah produk per tab (ikut kata kunci, kategori, stok)
}

var variantSorts = map[string]string{
	"updated":    "p.updated_at DESC, p.id, v.is_default DESC, v.sku",
	"name":       "lower(coalesce(t.name, '')), p.id, v.is_default DESC, v.sku",
	"sku":        "lower(v.sku)",
	"price_asc":  "v.price ASC, lower(v.sku)",
	"price_desc": "v.price DESC, lower(v.sku)",
	"stock_low":  "(v.stock_qty - v.reserved_qty) ASC, lower(v.sku)",
}

const variantLabel = `coalesce((SELECT string_agg(ov.value_i18n->>p.source_lang::text, ' / ' ORDER BY o.position)
	FROM variant_option_values vov JOIN product_option_values ov ON ov.id = vov.option_value_id
	JOIN product_options o ON o.id = vov.option_id WHERE vov.variant_id = v.id), '')`

// ListVariants: tabel SKU ringkas (satu baris per varian), dengan halaman.
func (s *Service) ListVariants(ctx context.Context, storeID string, q ProductQuery) (*VariantList, error) {
	tab, ok := productFilters[q.Filter]
	if !ok {
		tab, q.Filter = "true", "all"
	}
	order, ok := variantSorts[q.Sort]
	if !ok {
		order = variantSorts["updated"]
	}
	base := `v.store_id = $1 AND v.deleted_at IS NULL AND p.deleted_at IS NULL`
	args := []any{storeID}
	if q.Page.Q != "" {
		args = append(args, paging.Like(q.Page.Q))
		base += ` AND (v.sku ILIKE $2 OR EXISTS (SELECT 1 FROM product_translations x WHERE x.product_id = p.id AND x.name ILIKE $2))`
	}
	extra, args := catalogWhere(q, args, false)
	base += extra
	if !q.Page.On {
		q.Page.On, q.Page.Page, q.Page.Size = true, 1, 50
	}
	rows, err := s.pool.Query(ctx, `
		SELECT v.id, p.id, coalesce(t.name, ''),
		       coalesce((SELECT coalesce(m.image_url, m.storage_key) FROM product_media m WHERE m.product_id = p.id AND m.deleted_at IS NULL AND m.kind = 'photo'
		                 ORDER BY m.sort_order LIMIT 1), ''),
		       v.sku, `+variantLabel+`, v.price::text, v.compare_at_price::text, v.stock_qty, v.reserved_qty,
		       v.low_stock_threshold, `+effThreshold+`, v.status, p.status::text, p.moderation_status::text
		FROM product_variants v JOIN products p ON p.id = v.product_id
		LEFT JOIN product_translations t ON t.product_id = p.id AND t.lang = p.source_lang
		WHERE `+base+` AND `+tab+`
		ORDER BY `+order+` LIMIT `+fmt.Sprint(q.Page.Size)+` OFFSET `+fmt.Sprint(q.Page.Offset()), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := &VariantList{Variants: []VariantRow{}}
	for rows.Next() {
		var v VariantRow
		var cover string
		if err := rows.Scan(&v.ID, &v.ProductID, &v.ProductName, &cover, &v.SKU, &v.Label, &v.Price, &v.CompareAtPrice,
			&v.StockQty, &v.ReservedQty, &v.LowStock, &v.Threshold, &v.Status, &v.ProductStatus, &v.ModerationStatus); err != nil {
			return nil, err
		}
		v.Available = v.StockQty - v.ReservedQty
		v.CoverURL = cover
		if cover != "" && !strings.HasPrefix(cover, "http") {
			v.CoverURL = s.files.PublicURL(cover)
		}
		out.Variants = append(out.Variants, v)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	rows.Close()
	// Jumlah per tab = jumlah produk (sama dengan tampilan kartu); total halaman = jumlah baris SKU di tab ini.
	cols := make([]string, len(productFilterOrder))
	for i, f := range productFilterOrder {
		cols[i] = `count(DISTINCT p.id) FILTER (WHERE ` + productFilters[f] + `)`
	}
	vals := make([]int, len(productFilterOrder))
	ptrs := make([]any, len(vals)+1)
	for i := range vals {
		ptrs[i] = &vals[i]
	}
	var rowsTotal int
	ptrs[len(vals)] = &rowsTotal
	if err := s.pool.QueryRow(ctx, `SELECT `+strings.Join(cols, ", ")+`, count(*) FILTER (WHERE `+tab+`)
		FROM product_variants v JOIN products p ON p.id = v.product_id WHERE `+base, args...).Scan(ptrs...); err != nil {
		return nil, err
	}
	out.Counts = map[string]int{}
	for i, f := range productFilterOrder {
		out.Counts[f] = vals[i]
	}
	m := paging.NewMeta(q.Page, rowsTotal)
	out.Page = &m
	return out, nil
}

// ---- Ubah harga/stok (satu jalur untuk ubah cepat, aksi massal, dan impor) -------------------------

// variantEdit: nil = tidak diubah. ClearCompare menghapus harga coret.
type variantEdit struct {
	ID           string
	Price        *decimal.Decimal
	Compare      *decimal.Decimal
	ClearCompare bool
	Stock        *int
	StockDelta   *int // untuk aksi massal "tambah/kurangi stok"
	PricePct     *decimal.Decimal
	Status       *string
	Threshold    *int
	SetThreshold bool // true + Threshold nil = ikut ambang toko
	Field        string
}

// editResult: perubahan yang benar-benar terjadi.
type editResult struct {
	Products map[string]bool // produk yang variannya berubah
}

// priceGuard: aturan moderasi harga, dimuat sekali per permintaan.
type priceGuard struct {
	rules  modrules.Rules
	floors []modrules.Floor
}

func (s *Service) loadPriceGuard(ctx context.Context) (*priceGuard, error) {
	rules, err := modrules.Load(ctx, s.pool)
	if err != nil {
		return nil, err
	}
	floors, err := s.priceFloors(ctx)
	if err != nil {
		return nil, err
	}
	return &priceGuard{rules: rules, floors: floors}, nil
}

// editVariants menerapkan perubahan di dalam tx (sudah lewat db.WithActor → tercatat di audit). Validasi sama
// dengan formulir: harga > 0 dan ≤ batas, harga coret > harga, stok 0..batas dan ≥ stok yang sedang dipesan.
func editVariants(ctx context.Context, tx pgx.Tx, storeID, userID string, edits []variantEdit) (*editResult, error) {
	res := &editResult{Products: map[string]bool{}}
	limit := decimal.NewFromInt(maxPrice)
	for _, e := range edits {
		var productID, status, prodStatus string
		var price decimal.Decimal
		var compare *decimal.Decimal
		var stock, reserved int
		var thr *int
		err := tx.QueryRow(ctx, `
			SELECT v.product_id, v.price, v.compare_at_price, v.stock_qty, v.reserved_qty, v.status, v.low_stock_threshold, p.status::text
			FROM product_variants v JOIN products p ON p.id = v.product_id
			WHERE v.id = $1 AND v.store_id = $2 AND v.deleted_at IS NULL AND p.deleted_at IS NULL
			FOR UPDATE OF v`, e.ID, storeID).Scan(&productID, &price, &compare, &stock, &reserved, &status, &thr, &prodStatus)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		if err != nil {
			return nil, err
		}
		if prodStatus == "banned" {
			return nil, ErrProductBanned
		}
		f := e.Field
		newPrice, newCompare, newStock, newStatus, newThr := price, compare, stock, status, thr
		if e.PricePct != nil {
			k := decimal.NewFromInt(100).Add(*e.PricePct).Div(decimal.NewFromInt(100))
			newPrice = money.Round2(price.Mul(k))
			if compare != nil {
				c := money.Round2(compare.Mul(k))
				newCompare = &c
			}
		}
		if e.Price != nil {
			newPrice = money.Round2(*e.Price)
		}
		if e.ClearCompare {
			newCompare = nil
		} else if e.Compare != nil {
			c := money.Round2(*e.Compare)
			newCompare = &c
		}
		if e.Stock != nil {
			newStock = *e.Stock
		}
		if e.StockDelta != nil {
			newStock = stock + *e.StockDelta
			if newStock < reserved {
				newStock = reserved // "kurangi semua 5" tidak boleh membatalkan pesanan yang sudah masuk
			}
		}
		if e.Status != nil {
			newStatus = *e.Status
		}
		if e.SetThreshold {
			newThr = e.Threshold
		}
		if !newPrice.IsPositive() || newPrice.GreaterThan(limit) {
			return nil, invalid(f+"price", fmt.Sprintf("harga USD lebih dari 0, maksimal %d", maxPrice))
		}
		if newCompare != nil && !newCompare.GreaterThan(newPrice) {
			return nil, invalid(f+"compare_at_price", "harus lebih tinggi dari harga jual")
		}
		if newStock < 0 || newStock > maxStock {
			return nil, invalid(f+"stock_qty", fmt.Sprintf("0–%d", maxStock))
		}
		if newStock < reserved {
			return nil, invalid(f+"stock_qty", fmt.Sprintf("minimal %d (sedang dipesan)", reserved))
		}
		if newThr != nil && (*newThr < 0 || *newThr > 1000) {
			return nil, invalid(f+"low_stock_threshold", "0–1000")
		}
		same := newPrice.Equal(price) && newStock == stock && newStatus == status && eqDec(newCompare, compare) && eqInt(newThr, thr)
		if same {
			continue
		}
		var cs *string
		if newCompare != nil {
			c := newCompare.StringFixed(2)
			cs = &c
		}
		if _, err := tx.Exec(ctx, `UPDATE product_variants SET price = $2::numeric, compare_at_price = $3::numeric, stock_qty = $4, status = $5,
			low_stock_threshold = $6 WHERE id = $1`, e.ID, newPrice.StringFixed(2), cs, newStock, newStatus, newThr); err != nil {
			return nil, err
		}
		if d := newStock - stock; d != 0 {
			if _, err := tx.Exec(ctx, `INSERT INTO inventory_movements (variant_id, delta, reason, created_by) VALUES ($1, $2, 'adjustment', $3)`,
				e.ID, d, userID); err != nil {
				return nil, err
			}
		}
		res.Products[productID] = true
	}
	return res, nil
}

func eqDec(a, b *decimal.Decimal) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return a.Equal(*b)
}

func eqInt(a, b *int) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return *a == *b
}

// finishProducts: perbarui cache harga dan cek ulang harga produk yang tayang. Kalau harga baru kena aturan
// (turun lebih dari drop_pct dari harga yang disetujui, di bawah harga wajar, harga coret tidak wajar) →
// produk masuk antrean moderasi (🟡), sama seperti simpan lewat formulir. Mengembalikan produk yang dikirim ke review.
func (s *Service) finishProducts(ctx context.Context, tx pgx.Tx, g *priceGuard, products map[string]bool) ([]string, error) {
	ids := make([]string, 0, len(products))
	for id := range products {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	review := []string{}
	for _, id := range ids {
		if err := refreshPriceCache(ctx, tx, id); err != nil {
			return nil, err
		}
		var status, moderation, name string
		var approvedMin *decimal.Decimal
		if err := tx.QueryRow(ctx, `SELECT p.status::text, p.moderation_status::text, coalesce(t.name, ''), p.approved_min_price
			FROM products p LEFT JOIN product_translations t ON t.product_id = p.id AND t.lang = p.source_lang WHERE p.id = $1`, id).
			Scan(&status, &moderation, &name, &approvedMin); err != nil {
			return nil, err
		}
		if status == "draft" || moderation != "approved" {
			continue // draf: dicek saat dikirim · sudah menunggu admin: admin melihat harga terbaru
		}
		pc := modrules.PriceCheck{Name: name, ApprovedMin: approvedMin}
		rows, err := tx.Query(ctx, `SELECT price, compare_at_price FROM product_variants WHERE product_id = $1 AND deleted_at IS NULL AND status = 'active'`, id)
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var p decimal.Decimal
			var c *decimal.Decimal
			if err := rows.Scan(&p, &c); err != nil {
				rows.Close()
				return nil, err
			}
			pc.Prices = append(pc.Prices, p)
			pc.Compares = append(pc.Compares, c)
		}
		rows.Close()
		hits := modrules.CheckPrice(pc, g.floors, g.rules)
		if len(hits) == 0 {
			continue
		}
		reasons := []string{}
		for _, h := range hits {
			if !containsStr(reasons, h.Reason) {
				reasons = append(reasons, h.Reason)
			}
		}
		flags, _ := json.Marshal(map[string]any{"reasons": reasons, "lane": "manual", "price": hits, "source": "quick_edit"})
		if _, err := tx.Exec(ctx, `UPDATE products SET moderation_status = 'pending_review', submitted_at = now(), moderation_note = NULL,
			auto_flags = $2, review_lane = 'manual', robots_index = false WHERE id = $1`, id, flags); err != nil {
			return nil, err
		}
		review = append(review, id)
	}
	return review, nil
}

func containsStr(xs []string, x string) bool {
	for _, v := range xs {
		if v == x {
			return true
		}
	}
	return false
}

// VariantPatch: ubah cepat satu baris tabel SKU. Field yang tidak dikirim tidak diubah.
type VariantPatch struct {
	Price          *string `json:"price"`
	CompareAtPrice *string `json:"compare_at_price"` // "" = hapus harga coret
	StockQty       *int    `json:"stock_qty"`
	LowStock       optInt  `json:"low_stock_threshold"` // angka, atau null = ikut toko
}

// optInt membedakan "tidak dikirim" dari null (null = hapus ambang khusus varian).
type optInt struct {
	Set bool
	Val *int
}

func (o *optInt) UnmarshalJSON(b []byte) error {
	o.Set = true
	if string(b) == "null" {
		o.Val = nil
		return nil
	}
	var n int
	if err := json.Unmarshal(b, &n); err != nil {
		return invalid("low_stock_threshold", "angka 0–1000 atau kosong")
	}
	o.Val = &n
	return nil
}

// PatchResult: baris yang sudah diperbarui + apakah produknya masuk antrean moderasi karena harga.
type PatchResult struct {
	Variant VariantRow `json:"variant"`
	Review  bool       `json:"review"`
}

func parseMoney(field, s string) (decimal.Decimal, error) {
	d, err := decimal.NewFromString(strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(s), "$")))
	if err != nil {
		return d, invalid(field, "angka, mis. 12.50")
	}
	return d, nil
}

func (s *Service) PatchVariant(ctx context.Context, storeID, variantID, userID string, in VariantPatch) (*PatchResult, error) {
	e := variantEdit{ID: variantID}
	if in.Price != nil {
		d, err := parseMoney("price", *in.Price)
		if err != nil {
			return nil, err
		}
		e.Price = &d
	}
	if in.CompareAtPrice != nil {
		if strings.TrimSpace(*in.CompareAtPrice) == "" {
			e.ClearCompare = true
		} else {
			d, err := parseMoney("compare_at_price", *in.CompareAtPrice)
			if err != nil {
				return nil, err
			}
			e.Compare = &d
		}
	}
	e.Stock = in.StockQty
	if in.LowStock.Set {
		e.SetThreshold, e.Threshold = true, in.LowStock.Val
	}
	g, err := s.loadPriceGuard(ctx)
	if err != nil {
		return nil, err
	}
	var review []string
	err = db.WithActor(ctx, s.pool, userID, "seller", "", func(tx pgx.Tx) error {
		res, err := editVariants(ctx, tx, storeID, userID, []variantEdit{e})
		if err != nil {
			return err
		}
		review, err = s.finishProducts(ctx, tx, g, res.Products)
		return err
	})
	if err != nil {
		return nil, err
	}
	row, err := s.variantRow(ctx, storeID, variantID)
	if err != nil {
		return nil, err
	}
	return &PatchResult{Variant: row, Review: len(review) > 0}, nil
}

func (s *Service) variantRow(ctx context.Context, storeID, variantID string) (VariantRow, error) {
	var v VariantRow
	err := s.pool.QueryRow(ctx, `
		SELECT v.id, p.id, coalesce(t.name, ''), v.sku, `+variantLabel+`, v.price::text, v.compare_at_price::text, v.stock_qty, v.reserved_qty,
		       v.low_stock_threshold, `+effThreshold+`, v.status, p.status::text, p.moderation_status::text
		FROM product_variants v JOIN products p ON p.id = v.product_id
		LEFT JOIN product_translations t ON t.product_id = p.id AND t.lang = p.source_lang
		WHERE v.id = $1 AND v.store_id = $2 AND v.deleted_at IS NULL`, variantID, storeID).
		Scan(&v.ID, &v.ProductID, &v.ProductName, &v.SKU, &v.Label, &v.Price, &v.CompareAtPrice, &v.StockQty, &v.ReservedQty,
			&v.LowStock, &v.Threshold, &v.Status, &v.ProductStatus, &v.ModerationStatus)
	if errors.Is(err, pgx.ErrNoRows) {
		return v, ErrNotFound
	}
	v.Available = v.StockQty - v.ReservedQty
	return v, err
}

// ---- Aksi massal ----------------------------------------------------------------

const maxBulk = 100

type BulkInput struct {
	Action     string   `json:"action"` // archive · activate · delete_drafts · stock · price_pct
	ProductIDs []string `json:"product_ids"`
	StockMode  string   `json:"stock_mode"` // set · add (untuk action=stock)
	Stock      int      `json:"stock"`
	Pct        string   `json:"pct"` // untuk action=price_pct, mis. "-10" atau "5.5"
}

type BulkSkip struct {
	ID     string `json:"id"`
	Code   string `json:"code"` // not_found · banned · not_allowed · validation
	Detail string `json:"detail,omitempty"`
}

type BulkResult struct {
	Done    int        `json:"done"`
	Review  int        `json:"review"` // produk tayang yang masuk antrean moderasi karena harga
	Skipped []BulkSkip `json:"skipped"`
}

// Bulk: aksi untuk banyak produk sekaligus. Tiap produk diproses dalam transaksinya sendiri, jadi satu produk
// yang gagal (mis. stok di bawah jumlah dipesan) tidak membatalkan yang lain — dilaporkan di skipped.
func (s *Service) Bulk(ctx context.Context, storeID, userID string, in BulkInput) (*BulkResult, error) {
	if len(in.ProductIDs) == 0 || len(in.ProductIDs) > maxBulk {
		return nil, invalid("product_ids", fmt.Sprintf("pilih 1–%d produk", maxBulk))
	}
	var pct decimal.Decimal
	switch in.Action {
	case "archive", "activate", "delete_drafts":
	case "stock":
		if in.StockMode != "set" && in.StockMode != "add" {
			return nil, invalid("stock_mode", "set atau add")
		}
		if (in.StockMode == "set" && (in.Stock < 0 || in.Stock > maxStock)) || (in.StockMode == "add" && (in.Stock < -maxStock || in.Stock > maxStock || in.Stock == 0)) {
			return nil, invalid("stock", fmt.Sprintf("0–%d", maxStock))
		}
	case "price_pct":
		var err error
		pct, err = decimal.NewFromString(strings.TrimSpace(in.Pct))
		if err != nil || pct.IsZero() || pct.LessThanOrEqual(decimal.NewFromInt(-100)) || pct.GreaterThan(decimal.NewFromInt(500)) {
			return nil, invalid("pct", "persen antara -99 dan 500, bukan 0")
		}
	default:
		return nil, invalid("action", "archive, activate, delete_drafts, stock, atau price_pct")
	}
	var g *priceGuard
	if in.Action == "price_pct" || in.Action == "stock" {
		var err error
		if g, err = s.loadPriceGuard(ctx); err != nil {
			return nil, err
		}
	}
	res := &BulkResult{Skipped: []BulkSkip{}}
	seen := map[string]bool{}
	for _, id := range in.ProductIDs {
		if seen[id] {
			continue
		}
		seen[id] = true
		skip := func(err error) {
			var ve validationError
			switch {
			case errors.As(err, &ve):
				res.Skipped = append(res.Skipped, BulkSkip{ID: id, Code: "validation", Detail: ve.msg})
			case errors.Is(err, ErrNotFound):
				res.Skipped = append(res.Skipped, BulkSkip{ID: id, Code: "not_found"})
			case errors.Is(err, ErrProductBanned):
				res.Skipped = append(res.Skipped, BulkSkip{ID: id, Code: "banned"})
			default:
				res.Skipped = append(res.Skipped, BulkSkip{ID: id, Code: "error"})
			}
		}
		old, err := s.loadExisting(ctx, storeID, id)
		if err != nil {
			if errors.Is(err, ErrNotFound) || isBadUUID(err) {
				skip(ErrNotFound)
				continue
			}
			return nil, err
		}
		switch in.Action {
		case "archive", "activate":
			want := map[string]string{"archive": "archived", "activate": "active"}[in.Action]
			if old.status == want {
				res.Skipped = append(res.Skipped, BulkSkip{ID: id, Code: "not_allowed", Detail: "already"})
				continue
			}
			if _, err := s.SetProductStatus(ctx, storeID, id, userID, want); err != nil {
				var ve validationError
				if errors.As(err, &ve) {
					res.Skipped = append(res.Skipped, BulkSkip{ID: id, Code: "not_allowed", Detail: old.status})
				} else {
					skip(err)
				}
				continue
			}
		case "delete_drafts":
			if old.status != "draft" {
				res.Skipped = append(res.Skipped, BulkSkip{ID: id, Code: "not_allowed", Detail: old.status})
				continue
			}
			if err := s.DeleteProduct(ctx, storeID, id, userID); err != nil {
				skip(err)
				continue
			}
		case "stock", "price_pct":
			if old.status == "banned" {
				skip(ErrProductBanned)
				continue
			}
			var review []string
			err := db.WithActor(ctx, s.pool, userID, "seller", "", func(tx pgx.Tx) error {
				rows, err := tx.Query(ctx, `SELECT id FROM product_variants WHERE product_id = $1 AND deleted_at IS NULL ORDER BY is_default DESC, sku`, id)
				if err != nil {
					return err
				}
				var edits []variantEdit
				for rows.Next() {
					var vid string
					if err := rows.Scan(&vid); err != nil {
						rows.Close()
						return err
					}
					e := variantEdit{ID: vid}
					switch {
					case in.Action == "price_pct":
						e.PricePct = &pct
					case in.StockMode == "set":
						e.Stock = ptr(in.Stock)
					default:
						e.StockDelta = ptr(in.Stock)
					}
					edits = append(edits, e)
				}
				rows.Close()
				r, err := editVariants(ctx, tx, storeID, userID, edits)
				if err != nil {
					return err
				}
				review, err = s.finishProducts(ctx, tx, g, r.Products)
				return err
			})
			if err != nil {
				skip(err)
				continue
			}
			res.Review += len(review)
		}
		res.Done++
	}
	return res, nil
}

func isBadUUID(err error) bool {
	return err != nil && strings.Contains(err.Error(), "invalid input syntax for type uuid")
}

// ---- CSV ------------------------------------------------------------------------

// MaxImportBytes: batas ukuran file impor (± 20.000 baris).
const MaxImportBytes = 2 << 20

const maxImportRows = 20000

var csvHeader = []string{"sku", "product", "variant", "price", "compare_at_price", "stock", "status"}

// ExportCSV menulis semua SKU toko (UTF-8 dengan BOM supaya Excel membaca huruf Khmer/Cina dengan benar).
// status = status produk (draft · active · sold_out · archived · banned), atau "inactive" kalau variannya nonaktif.
func (s *Service) ExportCSV(ctx context.Context, storeID string, w io.Writer) error {
	rows, err := s.pool.Query(ctx, `
		SELECT v.sku, coalesce(t.name, ''), `+variantLabel+`, v.price::text, coalesce(v.compare_at_price::text, ''), v.stock_qty,
		       CASE WHEN v.status = 'inactive' THEN 'inactive' ELSE p.status::text END
		FROM product_variants v JOIN products p ON p.id = v.product_id
		LEFT JOIN product_translations t ON t.product_id = p.id AND t.lang = p.source_lang
		WHERE v.store_id = $1 AND v.deleted_at IS NULL AND p.deleted_at IS NULL
		ORDER BY lower(coalesce(t.name, '')), p.id, v.is_default DESC, lower(v.sku)`, storeID)
	if err != nil {
		return err
	}
	defer rows.Close()
	if _, err := w.Write([]byte("\xEF\xBB\xBF")); err != nil {
		return err
	}
	cw := csv.NewWriter(w)
	if err := cw.Write(csvHeader); err != nil {
		return err
	}
	for rows.Next() {
		var sku, name, label, price, compare, status string
		var stock int
		if err := rows.Scan(&sku, &name, &label, &price, &compare, &stock, &status); err != nil {
			return err
		}
		if err := cw.Write([]string{csvSafe(sku), csvSafe(name), csvSafe(label), price, compare, strconv.Itoa(stock), status}); err != nil {
			return err
		}
	}
	cw.Flush()
	if err := cw.Error(); err != nil {
		return err
	}
	return rows.Err()
}

// csvSafe: cegah formula injection di Excel (sel yang diawali = + - @ dijalankan sebagai rumus).
func csvSafe(s string) string {
	if s != "" && strings.ContainsRune("=+-@\t\r", rune(s[0])) {
		return "'" + s
	}
	return s
}

// ImportRow: satu baris file impor pada pratinjau.
type ImportRow struct {
	Line    int          `json:"line"`
	SKU     string       `json:"sku"`
	Product string       `json:"product,omitempty"`
	Variant string       `json:"variant,omitempty"`
	Old     *ImportVals  `json:"old,omitempty"`
	New     *ImportVals  `json:"new,omitempty"`
	Error   *ImportError `json:"error,omitempty"`
}

type ImportVals struct {
	Price   string  `json:"price"`
	Compare *string `json:"compare_at_price"`
	Stock   int     `json:"stock"`
}

// ImportError: code untuk teks UI (en/km/zh), detail = pesan teknis.
type ImportError struct {
	Code   string `json:"code"` // unknown_sku · duplicate_sku · bad_price · bad_compare · bad_stock · below_reserved · banned
	Detail string `json:"detail,omitempty"`
}

type ImportResult struct {
	Rows      []ImportRow `json:"rows"` // hanya baris yang berubah atau bermasalah
	Total     int         `json:"total"`
	Changed   int         `json:"changed"`
	Unchanged int         `json:"unchanged"`
	Errors    int         `json:"errors"`
	Applied   bool        `json:"applied"`
	Review    int         `json:"review"` // produk tayang yang akan/telah masuk antrean moderasi karena harga
	Columns   []string    `json:"columns"`
}

// ImportFileError: file tidak bisa dibaca (bukan CSV, tanpa kolom sku, terlalu besar, …).
type ImportFileError struct{ Code, Detail string }

func (e ImportFileError) Error() string { return "import: " + e.Code + " " + e.Detail }

// ImportCSV membaca file lalu (apply=false) mengembalikan pratinjau, atau (apply=true) menerapkan baris yang valid
// dalam satu transaksi. Kolom dikenali dari baris judul: sku wajib, lalu minimal satu dari price, compare_at_price,
// stock. Kolom lain (product, variant, status) diabaikan. Kolom harga coret yang ada tapi kosong = hapus harga coret.
func (s *Service) ImportCSV(ctx context.Context, storeID, userID string, data []byte, apply bool) (*ImportResult, error) {
	if len(data) > MaxImportBytes {
		return nil, ImportFileError{"too_large", fmt.Sprintf("max %d bytes", MaxImportBytes)}
	}
	data = bytes.TrimPrefix(data, []byte("\xEF\xBB\xBF"))
	if !utf8.Valid(data) {
		return nil, ImportFileError{"not_utf8", ""}
	}
	r := csv.NewReader(bytes.NewReader(data))
	r.Comma = sniffComma(data)
	r.FieldsPerRecord = -1
	r.TrimLeadingSpace = true
	records, err := r.ReadAll()
	if err != nil {
		return nil, ImportFileError{"bad_csv", err.Error()}
	}
	if len(records) < 2 {
		return nil, ImportFileError{"empty", ""}
	}
	if len(records)-1 > maxImportRows {
		return nil, ImportFileError{"too_many_rows", fmt.Sprintf("max %d", maxImportRows)}
	}
	col := map[string]int{}
	for i, h := range records[0] {
		k := strings.ToLower(strings.TrimSpace(h))
		k = strings.NewReplacer(" ", "_", "-", "_").Replace(k)
		switch k {
		case "compare", "compare_price", "compare_at", "old_price":
			k = "compare_at_price"
		case "stock_qty", "qty", "quantity":
			k = "stock"
		}
		if _, dup := col[k]; !dup {
			col[k] = i
		}
	}
	skuCol, ok := col["sku"]
	if !ok {
		return nil, ImportFileError{"no_sku_column", ""}
	}
	priceCol, hasPrice := col["price"]
	cmpCol, hasCmp := col["compare_at_price"]
	stockCol, hasStock := col["stock"]
	if !hasPrice && !hasCmp && !hasStock {
		return nil, ImportFileError{"no_value_columns", ""}
	}
	res := &ImportResult{Rows: []ImportRow{}, Columns: []string{"sku"}}
	for _, c := range []struct {
		ok   bool
		name string
	}{{hasPrice, "price"}, {hasCmp, "compare_at_price"}, {hasStock, "stock"}} {
		if c.ok {
			res.Columns = append(res.Columns, c.name)
		}
	}

	// Semua varian toko sekali baca (≤ puluhan ribu baris), dicari dengan SKU huruf kecil.
	type cur struct {
		id, productID, product, label, prodStatus string
		price                                     decimal.Decimal
		compare                                   *decimal.Decimal
		stock, reserved                           int
	}
	vars := map[string]*cur{}
	rows, err := s.pool.Query(ctx, `
		SELECT v.id, p.id, lower(v.sku), coalesce(t.name, ''), `+variantLabel+`, p.status::text, v.price, v.compare_at_price, v.stock_qty, v.reserved_qty
		FROM product_variants v JOIN products p ON p.id = v.product_id
		LEFT JOIN product_translations t ON t.product_id = p.id AND t.lang = p.source_lang
		WHERE v.store_id = $1 AND v.deleted_at IS NULL AND p.deleted_at IS NULL`, storeID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		c := &cur{}
		var sku string
		if err := rows.Scan(&c.id, &c.productID, &sku, &c.product, &c.label, &c.prodStatus, &c.price, &c.compare, &c.stock, &c.reserved); err != nil {
			rows.Close()
			return nil, err
		}
		vars[sku] = c
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	limit := decimal.NewFromInt(maxPrice)
	seen := map[string]int{}
	var edits []variantEdit
	for li, rec := range records[1:] {
		line := li + 2
		get := func(i int) string {
			if i < len(rec) {
				return strings.TrimSpace(rec[i])
			}
			return ""
		}
		sku := strings.TrimPrefix(get(skuCol), "'")
		blank := true
		for _, f := range rec {
			if strings.TrimSpace(f) != "" {
				blank = false
				break
			}
		}
		if blank {
			continue
		}
		res.Total++
		row := ImportRow{Line: line, SKU: sku}
		fail := func(code, detail string) {
			row.Error = &ImportError{Code: code, Detail: detail}
			res.Errors++
			res.Rows = append(res.Rows, row)
		}
		if sku == "" {
			fail("unknown_sku", "")
			continue
		}
		if prev, dup := seen[strings.ToLower(sku)]; dup {
			fail("duplicate_sku", fmt.Sprintf("line %d", prev))
			continue
		}
		seen[strings.ToLower(sku)] = line
		c, ok := vars[strings.ToLower(sku)]
		if !ok {
			fail("unknown_sku", "")
			continue
		}
		row.SKU, row.Product, row.Variant = sku, c.product, c.label
		if c.prodStatus == "banned" {
			fail("banned", "")
			continue
		}
		newPrice, newCompare, newStock := c.price, c.compare, c.stock
		e := variantEdit{ID: c.id, Field: fmt.Sprintf("line %d ", line)}
		if hasPrice && get(priceCol) != "" {
			d, err := decimal.NewFromString(cleanNum(get(priceCol)))
			if err != nil || !d.IsPositive() || d.GreaterThan(limit) {
				fail("bad_price", get(priceCol))
				continue
			}
			newPrice = money.Round2(d)
			e.Price = &newPrice
		}
		if hasCmp {
			if v := get(cmpCol); v == "" {
				newCompare, e.ClearCompare = nil, true
			} else {
				d, err := decimal.NewFromString(cleanNum(v))
				if err != nil {
					fail("bad_compare", v)
					continue
				}
				d = money.Round2(d)
				newCompare, e.Compare = &d, &d
			}
		}
		if newCompare != nil && !newCompare.GreaterThan(newPrice) {
			fail("bad_compare", newCompare.StringFixed(2))
			continue
		}
		if hasStock && get(stockCol) != "" {
			n, err := strconv.Atoi(cleanNum(get(stockCol)))
			if err != nil || n < 0 || n > maxStock {
				fail("bad_stock", get(stockCol))
				continue
			}
			if n < c.reserved {
				fail("below_reserved", strconv.Itoa(c.reserved))
				continue
			}
			newStock = n
			e.Stock = &newStock
		}
		if newPrice.Equal(c.price) && eqDec(newCompare, c.compare) && newStock == c.stock {
			res.Unchanged++
			continue
		}
		row.Old = &ImportVals{Price: c.price.StringFixed(2), Compare: decStr(c.compare), Stock: c.stock}
		row.New = &ImportVals{Price: newPrice.StringFixed(2), Compare: decStr(newCompare), Stock: newStock}
		res.Changed++
		res.Rows = append(res.Rows, row)
		edits = append(edits, e)
	}
	sort.SliceStable(res.Rows, func(i, j int) bool { // bermasalah dulu, lalu urut baris
		if (res.Rows[i].Error != nil) != (res.Rows[j].Error != nil) {
			return res.Rows[i].Error != nil
		}
		return res.Rows[i].Line < res.Rows[j].Line
	})
	if len(edits) == 0 {
		return res, nil
	}
	g, err := s.loadPriceGuard(ctx)
	if err != nil {
		return nil, err
	}
	// Pratinjau juga menjalankan perubahan di transaksi yang lalu dibatalkan — supaya jumlah "masuk review"
	// dan error validasi akhir persis sama dengan saat diterapkan.
	errRollback := errors.New("preview")
	err = db.WithActor(ctx, s.pool, userID, "seller", "", func(tx pgx.Tx) error {
		r, err := editVariants(ctx, tx, storeID, userID, edits)
		if err != nil {
			return err
		}
		review, err := s.finishProducts(ctx, tx, g, r.Products)
		if err != nil {
			return err
		}
		res.Review = len(review)
		if !apply {
			return errRollback
		}
		return nil
	})
	if err != nil && !errors.Is(err, errRollback) {
		return nil, err
	}
	res.Applied = apply
	return res, nil
}

func decStr(d *decimal.Decimal) *string {
	if d == nil {
		return nil
	}
	s := d.StringFixed(2)
	return &s
}

// cleanNum: "$1,250.00" → "1250.00", " 12 " → "12", "12,50" (koma desimal Excel Eropa) → "12.50".
func cleanNum(s string) string {
	s = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(s), "$"))
	if i := strings.IndexByte(s, ','); i >= 0 && strings.Count(s, ",") == 1 && !strings.Contains(s, ".") && len(s)-i-1 <= 2 {
		return strings.Replace(s, ",", ".", 1)
	}
	return strings.ReplaceAll(s, ",", "")
}

// sniffComma: Excel di beberapa bahasa menyimpan CSV dengan titik koma.
func sniffComma(data []byte) rune {
	line := data
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		line = data[:i]
	}
	if bytes.Count(line, []byte(";")) > bytes.Count(line, []byte(",")) {
		return ';'
	}
	if bytes.Count(line, []byte("\t")) > bytes.Count(line, []byte(",")) {
		return '\t'
	}
	return ','
}

// exportName: nama file ekspor, mis. phsar-skv-skus-2026-11-27.csv
func exportName(prefix string, now time.Time) string {
	return fmt.Sprintf("phsar-%s-skus-%s.csv", strings.ToLower(prefix), now.Format("2006-01-02"))
}
