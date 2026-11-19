package reports

// extras: strike dari sistem (tanpa admin), mis. pesanan batal karena telat diterima/dikemas.
// Aturannya (nyala/mati, poin, masa berlaku, batas harian) ada di paket extras; di sini hanya pencatatannya,
// dengan akibat yang sama seperti strike admin: notifikasi seller dan penangguhan otomatis di 3 poin aktif.
// Bisa dicabut admin lewat Revoke seperti strike biasa.

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"phsar/internal/notify"
)

type SystemStrike struct {
	StoreID string
	Reason  string // seller_timeout · pack_timeout
	Note    string // dilihat seller (5–1000 karakter)
	Points  int    // 1–3
	Days    int    // masa berlaku
	OrderID string
	OrderNo string
	Kind    string // jenis notifikasi seller
}

// AddSystemStrike: di dalam tx pemanggil. Satu strike per pesanan (kedua kalinya → "", nil).
func AddSystemStrike(ctx context.Context, tx pgx.Tx, in SystemStrike) (string, error) {
	if in.Points < 1 || in.Points > 3 || in.Days < 1 {
		return "", FieldError{"points", "1–3"}
	}
	var storeName, status string
	err := tx.QueryRow(ctx, `SELECT name, status::text FROM stores WHERE id::text = $1 AND deleted_at IS NULL FOR UPDATE`, in.StoreID).Scan(&storeName, &status)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", err
	}
	var id string
	err = tx.QueryRow(ctx, `INSERT INTO store_strikes (store_id, reason, note, points, auto, order_id, expires_at)
		VALUES ($1, $2, $3, $4, true, nullif($5, '')::uuid, now() + make_interval(days => $6))
		ON CONFLICT (order_id) WHERE order_id IS NOT NULL DO NOTHING RETURNING id`,
		in.StoreID, in.Reason, in.Note, in.Points, in.OrderID, in.Days).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	active, err := activePoints(ctx, tx, in.StoreID)
	if err != nil {
		return "", err
	}
	link := "#/store/" + in.StoreID
	kind := in.Kind
	if kind == "" {
		kind = notify.StoreStrike
	}
	if err := notify.Enqueue(ctx, tx, in.StoreID, kind, notify.Params{Store: storeName, Note: in.Note, Order: in.OrderNo, Point: fmt.Sprint(min(active, SuspendAt))}, link); err != nil {
		return "", err
	}
	if active >= SuspendAt && status == "active" {
		if _, err := tx.Exec(ctx, `UPDATE stores SET status = 'suspended', review_note = $2 WHERE id = $1`, in.StoreID,
			fmt.Sprintf("Suspended automatically after %d active strikes.", SuspendAt)); err != nil {
			return "", err
		}
		if _, err := tx.Exec(ctx, `UPDATE store_strikes SET caused_suspension = true WHERE id = $1`, id); err != nil {
			return "", err
		}
		if err := notify.Enqueue(ctx, tx, in.StoreID, notify.StoreSuspended, notify.Params{Store: storeName}, link); err != nil {
			return "", err
		}
	}
	return id, nil
}
