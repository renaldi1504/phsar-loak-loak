// Package extras: dua aturan kecil yang diatur super admin (platform_settings), default MATI.
//
//	auto_strikes  — pesanan batal karena telat (seller_timeout / pack_timeout) → strike otomatis untuk toko
//	                (poin per alasan, masa berlaku, batas strike otomatis per toko per hari Phnom Penh).
//	manual_payout — seller (pemilik) boleh "Tarik sekarang": payout pending yang sama seperti jadwal
//	                (minimal penarikan, maksimal N kali per 7 hari). Logikanya di paket payout (manual.go).
package extras

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"

	"phsar/internal/notify"
	"phsar/internal/reports"
)

// Log: logger paket (strike otomatis yang gagal hanya dicatat, pembatalan pesanan tetap jalan).
var Log = slog.Default()

// FieldError: isian pengaturan tidak valid (422).
type FieldError struct{ Field, Msg string }

func (e FieldError) Error() string { return e.Field + ": " + e.Msg }

type querier interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}

func load(ctx context.Context, q querier, key string, v any) error {
	var raw []byte
	err := q.QueryRow(ctx, `SELECT value FROM platform_settings WHERE key = $1`, key).Scan(&raw)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil
	}
	if err != nil {
		return err
	}
	_ = json.Unmarshal(raw, v) // isian rusak → tetap pakai default
	return nil
}

// Save menyimpan satu kunci pengaturan (dipanggil di dalam tx db.WithActor admin → tercatat di audit).
func Save(ctx context.Context, tx pgx.Tx, key, adminID string, v any) error {
	raw, _ := json.Marshal(v)
	_, err := tx.Exec(ctx, `INSERT INTO platform_settings (key, value, updated_by) VALUES ($1, $2, $3)
		ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_by = EXCLUDED.updated_by, updated_at = now()`, key, raw, adminID)
	return err
}

// ---------------------------------------------------------------------------
// A) Strike otomatis
// ---------------------------------------------------------------------------

const KeyAutoStrikes = "auto_strikes"

type AutoStrikes struct {
	Enabled             bool `json:"enabled"`
	SellerTimeoutPoints int  `json:"seller_timeout_points"` // 0–3 (0 = alasan ini tidak diberi strike)
	PackTimeoutPoints   int  `json:"pack_timeout_points"`   // 0–3
	ExpiryDays          int  `json:"expiry_days"`           // 1–365
	DailyMax            int  `json:"daily_max"`             // strike otomatis maksimal per toko per hari (1–10)
}

func DefaultAutoStrikes() AutoStrikes {
	return AutoStrikes{SellerTimeoutPoints: 1, PackTimeoutPoints: 1, ExpiryDays: 30, DailyMax: 1}
}

func LoadAutoStrikes(ctx context.Context, q querier) (AutoStrikes, error) {
	st := DefaultAutoStrikes()
	err := load(ctx, q, KeyAutoStrikes, &st)
	return st, err
}

func (st AutoStrikes) Check() error {
	switch {
	case st.SellerTimeoutPoints < 0 || st.SellerTimeoutPoints > 3:
		return FieldError{"seller_timeout_points", "0–3"}
	case st.PackTimeoutPoints < 0 || st.PackTimeoutPoints > 3:
		return FieldError{"pack_timeout_points", "0–3"}
	case st.ExpiryDays < 1 || st.ExpiryDays > 365:
		return FieldError{"expiry_days", "1–365 days"}
	case st.DailyMax < 1 || st.DailyMax > 10:
		return FieldError{"daily_max", "1–10"}
	}
	return nil
}

func (st AutoStrikes) points(reason string) int {
	switch reason {
	case "seller_timeout":
		return st.SellerTimeoutPoints
	case "pack_timeout":
		return st.PackTimeoutPoints
	}
	return 0
}

// AutoStrike dipanggil di dalam tx pembatalan otomatis (shop.expireLate), setelah pesanan dibatalkan.
// Memakai savepoint: kalau strike gagal, hanya strike-nya yang dibatalkan — pembatalan & refund tetap jalan.
// Hasil: true kalau strike dibuat.
func AutoStrike(ctx context.Context, tx pgx.Tx, storeID, orderID, orderNo, reason string) bool {
	st, err := LoadAutoStrikes(ctx, tx)
	if err != nil {
		Log.Error("strike otomatis: pengaturan", "err", err)
		return false
	}
	pts := st.points(reason)
	if !st.Enabled || pts == 0 {
		return false
	}
	sp, err := tx.Begin(ctx) // savepoint
	if err != nil {
		Log.Error("strike otomatis: savepoint", "err", err)
		return false
	}
	defer sp.Rollback(ctx) //nolint:errcheck
	// Kunci baris toko dulu supaya hitungan batas harian aman dari dua pembatalan bersamaan.
	if _, err := sp.Exec(ctx, `SELECT 1 FROM stores WHERE id = $1 FOR UPDATE`, storeID); err != nil {
		Log.Error("strike otomatis: kunci toko", "err", err)
		return false
	}
	var today int
	if err := sp.QueryRow(ctx, `SELECT count(*) FROM store_strikes WHERE store_id = $1 AND auto
		AND created_at >= (date_trunc('day', now() AT TIME ZONE 'Asia/Phnom_Penh') AT TIME ZONE 'Asia/Phnom_Penh')`, storeID).Scan(&today); err != nil {
		Log.Error("strike otomatis: hitung harian", "err", err)
		return false
	}
	if today >= st.DailyMax {
		return false // anti-hukuman beruntun: batas harian tercapai
	}
	kind := notify.StoreStrikeLateAccept
	note := fmt.Sprintf("Order %s was not accepted in time and was cancelled automatically.", orderNo)
	if reason == "pack_timeout" {
		kind = notify.StoreStrikeLatePack
		note = fmt.Sprintf("Order %s was not packed in time and was cancelled automatically.", orderNo)
	}
	id, err := reports.AddSystemStrike(ctx, sp, reports.SystemStrike{StoreID: storeID, Reason: reason, Note: note, Points: pts,
		Days: st.ExpiryDays, OrderID: orderID, OrderNo: orderNo, Kind: kind})
	if err != nil {
		Log.Error("strike otomatis", "order", orderNo, "err", err)
		return false
	}
	if err := sp.Commit(ctx); err != nil {
		Log.Error("strike otomatis: commit savepoint", "err", err)
		return false
	}
	return id != ""
}

// ---------------------------------------------------------------------------
// B) Tarik saldo manual
// ---------------------------------------------------------------------------

const KeyManualPayout = "manual_payout"

type ManualPayout struct {
	Enabled     bool   `json:"enabled"`
	MinAmount   string `json:"min_amount"`     // minimal saldo untuk tarik manual
	MaxPer7Days int    `json:"max_per_7_days"` // 1–14
}

func DefaultManualPayout() ManualPayout {
	return ManualPayout{MinAmount: "10.00", MaxPer7Days: 1}
}

func LoadManualPayout(ctx context.Context, q querier) (ManualPayout, error) {
	st := DefaultManualPayout()
	err := load(ctx, q, KeyManualPayout, &st)
	return st, err
}

// Check memvalidasi dan merapikan min_amount (2 desimal).
func (st *ManualPayout) Check() error {
	min, err := decimal.NewFromString(strings.TrimSpace(st.MinAmount))
	if err != nil || !min.IsPositive() || min.GreaterThan(decimal.NewFromInt(10000)) {
		return FieldError{"min_amount", "enter an amount between 0.01 and 10000"}
	}
	if st.MaxPer7Days < 1 || st.MaxPer7Days > 14 {
		return FieldError{"max_per_7_days", "1–14"}
	}
	st.MinAmount = min.StringFixed(2)
	return nil
}
