package shop

// promos: voucher (dibiayai Phsar) & flash sale (pembiayaan dipilih per sesi).
//
// Voucher
//   - Dicek di keranjang (quote, tanpa login: aturan umum) dan di checkout (lengkap: batas per pembeli, pesanan pertama).
//   - Potongan mengurangi total yang dibayar (KHQR). Penghasilan seller tidak berubah: Phsar yang menanggung potongan.
//   - Dibagi ke tiap pesanan sesuai porsi subtotal (orders.discount); sisa pembulatan masuk pesanan terakhir.
//   - Pemakaian dicatat "held" saat checkout, "used" saat lunas, "released" kalau checkout kedaluwarsa/batal
//     (kuota kembali). Total dibayar minimal $0.01 (KHQR tidak bisa 0).
//   - Voucher "pesanan pertama" hanya untuk akun login (Telegram / HP / Facebook), tidak untuk tamu: akun tamu
//     bisa dibuat berulang kali, jadi syarat "pertama" tidak bisa dijamin. Voucher biasa boleh dipakai tamu.
//
// Flash sale (sesi)
//   - Harga khusus untuk barang yang disetujui di sesi yang terbit & sedang berjalan (view flash_live).
//   - Kuota & batas per pembeli: harga flash hanya kalau cukup untuk jumlah diminta (tidak dipecah sebagian).
//   - Pembiayaan per sesi: phsar_share_pct dari potongan dibayar Phsar (orders.flash_subsidy), sisanya seller.

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

// VoucherError: voucher tidak bisa dipakai. Code dipakai halaman untuk pesan dalam bahasa pembeli.
type VoucherError struct {
	Code   string // voucher_not_found · voucher_expired · voucher_not_started · voucher_min_spend · voucher_used_up · voucher_user_limit · voucher_first_order · voucher_login_required
	Amount string // untuk voucher_min_spend: belanja minimal
}

func (e *VoucherError) Error() string { return e.Code }

// VoucherInfo: voucher yang sedang dipakai (dikirim ke halaman keranjang/checkout).
type VoucherInfo struct {
	Code     string `json:"code"`
	Name     string `json:"name,omitempty"`
	Kind     string `json:"kind"`
	Value    string `json:"value"`
	Discount string `json:"discount"`
}

type voucher struct {
	id, code, name, kind string
	value, minSpend      decimal.Decimal
	maxDiscount          *decimal.Decimal
	startsAt             time.Time
	endsAt               *time.Time
	totalLimit           *int
	perUser              int
	firstOnly, active    bool
}

func normCode(c string) string { return strings.ToUpper(strings.TrimSpace(c)) }

type rowQuerier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// loadVoucher: lock=true mengunci baris voucher (di dalam transaksi checkout) supaya kuota tidak terlampaui.
func loadVoucher(ctx context.Context, q rowQuerier, code string, lock bool) (*voucher, error) {
	code = normCode(code)
	if code == "" || len(code) > 24 {
		return nil, &VoucherError{Code: "voucher_not_found"}
	}
	sql := `SELECT id, code, name, kind, value, min_spend, max_discount, starts_at, ends_at, total_limit, per_user_limit, first_order_only, active
		FROM vouchers WHERE upper(code) = $1`
	if lock {
		sql += ` FOR UPDATE`
	}
	v := &voucher{}
	err := q.QueryRow(ctx, sql, code).Scan(&v.id, &v.code, &v.name, &v.kind, &v.value, &v.minSpend, &v.maxDiscount, &v.startsAt, &v.endsAt,
		&v.totalLimit, &v.perUser, &v.firstOnly, &v.active)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &VoucherError{Code: "voucher_not_found"}
	}
	if err != nil {
		return nil, err
	}
	if !v.active {
		return nil, &VoucherError{Code: "voucher_not_found"}
	}
	now := time.Now()
	if now.Before(v.startsAt) {
		return nil, &VoucherError{Code: "voucher_not_started"}
	}
	if v.endsAt != nil && !now.Before(*v.endsAt) {
		return nil, &VoucherError{Code: "voucher_expired"}
	}
	if v.totalLimit != nil {
		var used int
		if err := q.QueryRow(ctx, `SELECT count(*) FROM voucher_redemptions WHERE voucher_id = $1 AND status <> 'released'`, v.id).Scan(&used); err != nil {
			return nil, err
		}
		if used >= *v.totalLimit {
			return nil, &VoucherError{Code: "voucher_used_up"}
		}
	}
	return v, nil
}

// checkBuyer: aturan per pembeli (hanya di checkout, pembeli sudah login).
func (v *voucher) checkBuyer(ctx context.Context, q rowQuerier, buyerID string) error {
	var guest bool
	var used, paid int
	if err := q.QueryRow(ctx, `SELECT u.is_guest,
		(SELECT count(*) FROM voucher_redemptions r WHERE r.voucher_id = $2 AND r.buyer_id = u.id AND r.status <> 'released'),
		(SELECT count(*) FROM checkouts c WHERE c.buyer_id = u.id AND c.status = 'paid')
		FROM users u WHERE u.id = $1`, buyerID, v.id).Scan(&guest, &used, &paid); err != nil {
		return err
	}
	if v.firstOnly && guest {
		return &VoucherError{Code: "voucher_login_required"}
	}
	if v.firstOnly && paid > 0 {
		return &VoucherError{Code: "voucher_first_order"}
	}
	if used >= v.perUser {
		return &VoucherError{Code: "voucher_user_limit"}
	}
	return nil
}

// voucherDiscount: potongan voucher untuk keranjang ini (barang flash "tanpa voucher" tidak dihitung).
func (qt *Quote) voucherDiscount(v *voucher) (decimal.Decimal, error) {
	base, _ := qt.voucherBase()
	if !base.IsPositive() && dec(qt.Subtotal).IsPositive() {
		return decimal.Zero, &VoucherError{Code: "voucher_flash"}
	}
	return v.discountFor(base, qt.gross())
}

// discountFor: potongan untuk subtotal barang (tanpa ongkir), dibatasi supaya total dibayar ≥ $0.01.
func (v *voucher) discountFor(subtotal, payable decimal.Decimal) (decimal.Decimal, error) {
	if subtotal.LessThan(v.minSpend) {
		return decimal.Zero, &VoucherError{Code: "voucher_min_spend", Amount: v.minSpend.StringFixed(2)}
	}
	d := v.value
	if v.kind == "percent" {
		d = subtotal.Mul(v.value).Div(decimal.NewFromInt(100)).Round(2)
	}
	if v.maxDiscount != nil && d.GreaterThan(*v.maxDiscount) {
		d = *v.maxDiscount
	}
	if d.GreaterThan(subtotal) {
		d = subtotal
	}
	if limit := payable.Sub(decimal.New(1, -2)); d.GreaterThan(limit) {
		d = limit
	}
	if !d.IsPositive() {
		return decimal.Zero, &VoucherError{Code: "voucher_min_spend", Amount: v.minSpend.StringFixed(2)}
	}
	return d, nil
}

func dec(s string) decimal.Decimal {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.Zero
	}
	return d
}

// gross: subtotal + ongkir sebelum potongan (dibaca dari string supaya mengikuti cara ongkir dihitung).
func (qt *Quote) gross() decimal.Decimal { return dec(qt.Subtotal).Add(dec(qt.ShippingFee)) }

// payable: total yang harus dibayar pembeli (setelah potongan voucher).
func (qt *Quote) payable() decimal.Decimal { return qt.gross().Sub(qt.discount) }

// applyVoucher memasang potongan ke quote (Total ikut berubah).
func (qt *Quote) applyVoucher(v *voucher, d decimal.Decimal) {
	qt.discount, qt.voucher = d, v
	qt.Discount = d.StringFixed(2)
	qt.Voucher = &VoucherInfo{Code: v.code, Name: v.name, Kind: v.kind, Value: v.value.StringFixed(2), Discount: qt.Discount}
	qt.Total = qt.payable().StringFixed(2)
}

// allocate: potongan per toko sesuai porsi subtotal.
func (qt *Quote) allocate() map[string]decimal.Decimal {
	out := map[string]decimal.Decimal{}
	if !qt.discount.IsPositive() || len(qt.Stores) == 0 {
		return out
	}
	sub, per := qt.voucherBase()
	left := qt.discount
	for i, g := range qt.Stores {
		var d decimal.Decimal
		if i == len(qt.Stores)-1 {
			d = left
		} else if sub.IsPositive() {
			d = qt.discount.Mul(per[g.StoreID]).Div(sub).RoundDown(2)
		}
		if d.GreaterThan(g.subtotal) {
			d = g.subtotal
		}
		left = left.Sub(d)
		out[g.StoreID] = d
	}
	return out
}

// QuoteWithVoucher: keranjang + voucher (aturan umum saja; batas per pembeli dicek saat checkout).
// Voucher yang tidak berlaku tidak menggagalkan keranjang: alasannya dikirim di voucher_error.
func (s *Service) QuoteWithVoucher(ctx context.Context, lang string, items []CartItem, code string) (*Quote, error) {
	qt, err := s.Quote(ctx, lang, items)
	if err != nil || normCode(code) == "" || qt.ItemCount == 0 {
		return qt, err
	}
	v, err := loadVoucher(ctx, s.Pool, code, false)
	if err == nil {
		var d decimal.Decimal
		if d, err = qt.voucherDiscount(v); err == nil {
			qt.applyVoucher(v, d)
			return qt, nil
		}
	}
	var ve *VoucherError
	if errors.As(err, &ve) {
		qt.VoucherError, qt.VoucherMinSpend = ve.Code, ve.Amount
		return qt, nil
	}
	return nil, err
}

// voucherForCheckout: di dalam transaksi checkout — kunci voucher, cek semua aturan, pasang potongan ke quote.
func (s *Service) voucherForCheckout(ctx context.Context, tx pgx.Tx, buyerID, code string, qt *Quote) error {
	if normCode(code) == "" {
		return nil
	}
	v, err := loadVoucher(ctx, tx, code, true)
	if err != nil {
		return err
	}
	if err := v.checkBuyer(ctx, tx, buyerID); err != nil {
		return err
	}
	d, err := qt.voucherDiscount(v)
	if err != nil {
		return err
	}
	qt.applyVoucher(v, d)
	return nil
}

// saveVoucher: simpan potongan di checkout + catat pemakaian (held).
func saveVoucher(ctx context.Context, tx pgx.Tx, checkoutID, buyerID string, qt *Quote) error {
	if qt.voucher == nil || !qt.discount.IsPositive() {
		return nil
	}
	if _, err := tx.Exec(ctx, `UPDATE checkouts SET voucher_id = $2, discount = $3, total = total - $3 WHERE id = $1`,
		checkoutID, qt.voucher.id, qt.discount); err != nil {
		return err
	}
	_, err := tx.Exec(ctx, `INSERT INTO voucher_redemptions (voucher_id, checkout_id, buyer_id, discount) VALUES ($1, $2, $3, $4)`,
		qt.voucher.id, checkoutID, buyerID, qt.discount)
	return err
}

// orderDiscount: pasang potongan pesanan (total pesanan berkurang, seller_amount tetap).
func orderDiscount(ctx context.Context, tx pgx.Tx, orderID string, d decimal.Decimal) error {
	if !d.IsPositive() {
		return nil
	}
	_, err := tx.Exec(ctx, `UPDATE orders SET discount = $2, total = total - $2 WHERE id = $1`, orderID, d)
	return err
}

// redemptionStatus: held → used (lunas) / released (kedaluwarsa, batal).
func redemptionStatus(ctx context.Context, tx pgx.Tx, checkoutID, status string) error {
	_, err := tx.Exec(ctx, `UPDATE voucher_redemptions SET status = $2 WHERE checkout_id = $1 AND status <> $2`, checkoutID, status)
	return err
}

// ---------------------------------------------------------------------------
// Flash sale (sesi; lihat migrasi 00036 & view flash_live)
// ---------------------------------------------------------------------------

type flashBuyerKey struct{}

// withFlashBuyer: checkout memberi tahu pembelinya supaya batas beli per pembeli dihitung dari pembelian sebelumnya.
func withFlashBuyer(ctx context.Context, buyerID string) context.Context {
	return context.WithValue(ctx, flashBuyerKey{}, buyerID)
}

// applyFlash: harga flash untuk varian yang sedang flash sale — hanya kalau sisa kuota dan batas per pembeli
// cukup untuk jumlah yang diminta (tidak dipecah sebagian harga flash, sebagian normal).
// lock=true (checkout): baris flash_items dikunci supaya kuota tidak terlampaui oleh dua pembeli bersamaan.
func applyFlash(ctx context.Context, q querier, found map[string]Line, items []CartItem, lock bool) error {
	if len(found) == 0 {
		return nil
	}
	want := map[string]int{}
	ids := []string{}
	for _, it := range items {
		if _, ok := found[it.VariantID]; ok {
			want[it.VariantID] = it.Qty
			ids = append(ids, it.VariantID)
		}
	}
	if lock {
		r, err := q.Query(ctx, `SELECT id FROM flash_items WHERE variant_id = ANY($1::uuid[]) AND status = 'approved' ORDER BY id FOR UPDATE`, ids)
		if err != nil {
			return err
		}
		r.Close()
	}
	buyer, _ := ctx.Value(flashBuyerKey{}).(string)
	rows, err := q.Query(ctx, `SELECT DISTINCT ON (f.variant_id) f.variant_id, f.id, f.sale_price, f.normal_price, f.ends_at, f.qty_cap - f.sold,
		  f.per_buyer, f.phsar_share_pct, f.allow_vouchers,
		  CASE WHEN $2 = '' THEN 0 ELSE coalesce((SELECT sum(oi.qty) FROM order_items oi JOIN orders o ON o.id = oi.order_id
		    WHERE oi.flash_item_id = f.id AND o.buyer_id::text = $2 AND o.status NOT IN ('cancelled','refunded')), 0) END::int
		FROM flash_live f WHERE f.variant_id = ANY($1::uuid[])
		ORDER BY f.variant_id, f.sale_price`, ids, buyer)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var vid, fid string
		var price, normal decimal.Decimal
		var ends time.Time
		var left, perBuyer, pct, bought int
		var vouchers bool
		if err := rows.Scan(&vid, &fid, &price, &normal, &ends, &left, &perBuyer, &pct, &vouchers, &bought); err != nil {
			return err
		}
		l := found[vid]
		if left < want[vid] || bought+want[vid] > perBuyer {
			continue // kuota habis / melebihi batas per pembeli → harga normal
		}
		if !price.LessThan(l.price) {
			continue
		}
		if l.Compare == "" || dec(l.Compare).LessThan(l.price) {
			l.Compare = l.price.StringFixed(2)
		}
		gap := normal.Sub(price)
		if cur := l.price.Sub(price); cur.LessThan(gap) {
			gap = cur // harga normal turun sejak diajukan: subsidi dihitung dari harga sekarang
		}
		l.price, l.Price = price, price.StringFixed(2)
		l.FlashItemID, l.FlashEndsAt, l.FlashPerBuyer = fid, &ends, perBuyer
		l.flashGap, l.flashPct, l.flashNoVouch = gap, pct, !vouchers
		n := left
		l.FlashLeft = &n
		found[vid] = l
	}
	return rows.Err()
}

// flashSubsidy: bagian potongan baris ini yang ditanggung Phsar.
func flashSubsidy(l Line) decimal.Decimal {
	if l.FlashItemID == "" || l.flashPct <= 0 {
		return decimal.Zero
	}
	return l.flashGap.Mul(decimal.NewFromInt(int64(l.Qty))).Mul(decimal.NewFromInt(int64(l.flashPct))).Div(decimal.NewFromInt(100)).Round(2)
}

// flashGross: baris dengan harga "untuk seller" (harga flash + subsidi Phsar) dan total subsidinya.
func flashGross(lines []Line) ([]Line, decimal.Decimal) {
	out := make([]Line, len(lines))
	total := decimal.Zero
	for i, l := range lines {
		if s := flashSubsidy(l); s.IsPositive() {
			l.price = l.price.Add(s.Div(decimal.NewFromInt(int64(l.Qty))))
			total = total.Add(s)
		}
		out[i] = l
	}
	return out, total
}

// flashItem: tandai barang pesanan yang dibeli dengan harga flash (kuota, batas per pembeli, subsidi).
func flashItem(ctx context.Context, tx pgx.Tx, itemID string, l Line) error {
	if l.FlashItemID == "" {
		return nil
	}
	_, err := tx.Exec(ctx, `UPDATE order_items SET flash_item_id = $2, flash_subsidy = $3 WHERE id = $1`, itemID, l.FlashItemID, flashSubsidy(l))
	return err
}

// voucherBase: subtotal barang yang boleh kena voucher (barang flash dari sesi "tanpa voucher" tidak ikut).
func (qt *Quote) voucherBase() (decimal.Decimal, map[string]decimal.Decimal) {
	total, per := decimal.Zero, map[string]decimal.Decimal{}
	for _, g := range qt.Stores {
		sub := decimal.Zero
		for _, l := range g.Lines {
			if l.Qty <= 0 || l.flashNoVouch {
				continue
			}
			sub = sub.Add(l.price.Mul(decimal.NewFromInt(int64(l.Qty))))
		}
		per[g.StoreID] = sub
		total = total.Add(sub)
	}
	return total, per
}

var voucherMsgs = map[string]map[string]string{
	"voucher_not_found":      {"en": "This voucher code is not valid.", "km": "លេខកូដប័ណ្ណនេះមិនត្រឹមត្រូវទេ។", "zh": "此优惠码无效。"},
	"voucher_expired":        {"en": "This voucher has expired.", "km": "ប័ណ្ណនេះផុតកំណត់ហើយ។", "zh": "此优惠券已过期。"},
	"voucher_not_started":    {"en": "This voucher is not active yet.", "km": "ប័ណ្ណនេះមិនទាន់អាចប្រើបានទេ។", "zh": "此优惠券尚未开始。"},
	"voucher_min_spend":      {"en": "Spend at least ${a} on items to use this voucher.", "km": "ទិញយ៉ាងហោចណាស់ ${a} ដើម្បីប្រើប័ណ្ណនេះ។", "zh": "商品满 ${a} 才能使用此优惠券。"},
	"voucher_used_up":        {"en": "This voucher has been fully used.", "km": "ប័ណ្ណនេះត្រូវបានប្រើអស់ហើយ។", "zh": "此优惠券已被领完。"},
	"voucher_user_limit":     {"en": "You have already used this voucher.", "km": "អ្នកបានប្រើប័ណ្ណនេះរួចហើយ។", "zh": "您已使用过此优惠券。"},
	"voucher_first_order":    {"en": "This voucher is for your first order only.", "km": "ប័ណ្ណនេះសម្រាប់ការបញ្ជាទិញលើកដំបូងប៉ុណ្ណោះ។", "zh": "此优惠券仅限首单使用。"},
	"voucher_flash":          {"en": "Flash sale items in this session cannot be combined with vouchers.", "km": "ទំនិញ Flash sale ក្នុងវគ្គនេះមិនអាចប្រើជាមួយប័ណ្ណបានទេ។", "zh": "本场限时抢购商品不能与优惠券同时使用。"},
	"voucher_login_required": {"en": "Sign in with Telegram, phone or Facebook to use this voucher.", "km": "សូមចូលដោយ Telegram លេខទូរស័ព្ទ ឬ Facebook ដើម្បីប្រើប័ណ្ណនេះ។", "zh": "请使用 Telegram、手机号或 Facebook 登录后使用此优惠券。"},
}

func voucherMsg(e *VoucherError, lang string) string {
	m := voucherMsgs[e.Code]
	s, ok := m[lang]
	if !ok {
		s = m["en"]
	}
	return strings.ReplaceAll(s, "{a}", e.Amount)
}
