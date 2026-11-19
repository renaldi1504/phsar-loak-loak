// Package pricehistory: riwayat harga normal varian + aturan flash sale "harga flash di bawah harga terendah 30 hari".
//
//	product_variants.price berubah ──trigger──► variant_price_history (valid_from / valid_to)
//	variant_low_price(variant, T) = harga normal terendah dalam 30 hari sebelum T (fungsi SQL, migrasi 00042)
//	flash: sale_price < variant_low_price(variant, sesi mulai)  — bisa dimatikan super admin (platform_settings['flash_rules'])
package pricehistory

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"

	"phsar/internal/i18n"
)

// Days: panjang jendela harga terendah.
const Days = 30

type querier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// Rules: platform_settings['flash_rules'].
type Rules struct {
	Low30Enabled bool `json:"low30_enabled"`
}

// LoadRules: tanpa baris = aturan nyala (default aman).
func LoadRules(ctx context.Context, q querier) (Rules, error) {
	r := Rules{Low30Enabled: true}
	var raw []byte
	err := q.QueryRow(ctx, `SELECT value FROM platform_settings WHERE key = 'flash_rules'`).Scan(&raw)
	if errors.Is(err, pgx.ErrNoRows) {
		return r, nil
	}
	if err != nil {
		return r, err
	}
	if err := json.Unmarshal(raw, &r); err != nil {
		return Rules{Low30Enabled: true}, nil
	}
	return r, nil
}

// Low: harga normal terendah varian dalam 30 hari sebelum `at`.
func Low(ctx context.Context, q querier, variantID string, at time.Time) (decimal.Decimal, error) {
	var low decimal.Decimal
	err := q.QueryRow(ctx, `SELECT variant_low_price($1::uuid, $2, $3)`, variantID, at, Days).Scan(&low)
	return low, err
}

// CheckFlashPrice: "" = boleh; selain itu pesan (bahasa request) untuk field sale_price.
// Dipanggil paket flash di dalam transaksi tambah/daftar/ubah barang.
func CheckFlashPrice(ctx context.Context, q querier, variantID string, start time.Time, sale decimal.Decimal, name string) (string, error) {
	rules, err := LoadRules(ctx, q)
	if err != nil || !rules.Low30Enabled {
		return "", err
	}
	low, err := Low(ctx, q, variantID, start)
	if err != nil {
		return "", err
	}
	if sale.LessThan(low) {
		return "", nil
	}
	return Msg(i18n.From(ctx), name, low), nil
}

var msgs = map[string]string{
	"en": "%s: the flash price must be below the lowest price of the last 30 days ($%s) — at most $%s",
	"km": "%s: តម្លៃ flash ត្រូវតែទាបជាងតម្លៃទាបបំផុតក្នុង 30 ថ្ងៃចុងក្រោយ ($%s) — អតិបរមា $%s",
	"zh": "%s：抢购价必须低于近 30 天最低价（$%s）——最高 $%s",
}

// Msg: pesan penolakan terlokalisasi (en / km / zh).
func Msg(l i18n.Lang, name string, low decimal.Decimal) string {
	max := low.Sub(decimal.New(1, -2))
	if max.IsNegative() {
		max = decimal.Zero
	}
	return fmt.Sprintf(i18n.Pick(msgs, l), name, low.StringFixed(2), max.StringFixed(2))
}
