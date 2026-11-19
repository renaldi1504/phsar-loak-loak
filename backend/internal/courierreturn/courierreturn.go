// Package courierreturn: pengantaran gagal → paket dikembalikan ke toko.
//
//	failed ──kurir/admin: kembalikan──► returning ──kurir: sudah di toko──► returned (menunggu seller)
//	                                         │                                   │
//	                                         └──────── seller / admin / otomatis (AutoReceiveAfter) ──► diterima
//	diterima → pesanan batal (cancel_reason delivery_failed) + refund pending (refund manual admin yang ada)
//	         + stok kembali, kecuali seller menandai barang rusak.
//
// Status pengantaran disimpan di deliveries (returning / returned) supaya daftar kurir & dasbor admin ikut;
// rincian pengembalian (alasan, waktu, kondisi barang, ongkir pengembalian) di courier_returns.
package courierreturn

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"phsar/internal/delivery"
	"phsar/internal/jobs"
	"phsar/internal/notify"
	"phsar/internal/platform/db"
	"phsar/internal/storage"
)

var (
	ErrNotCourier = errors.New("not a courier")
	ErrNotFound   = errors.New("not found")
	ErrBadStatus  = errors.New("status does not allow this")
	ErrValidation = errors.New("validation")
)

type fieldError struct{ field, msg string }

func (e fieldError) Error() string { return e.field + ": " + e.msg }
func (e fieldError) Unwrap() error { return ErrValidation }

// AutoReceiveAfter: seller tidak mengonfirmasi dalam selang ini setelah kurir menyerahkan → dianggap diterima.
var AutoReceiveAfter = 24 * time.Hour

type Service struct {
	Pool *pgxpool.Pool
	Log  *slog.Logger
}

// Return: satu pengembalian (tampilan admin & seller; data pembeli tidak pernah ikut).
type Return struct {
	ID           string     `json:"id"`
	DeliveryID   string     `json:"delivery_id"`
	OrderNo      string     `json:"order_no"`
	OrderStatus  string     `json:"order_status"`
	StoreID      string     `json:"store_id"`
	StoreName    string     `json:"store_name"`
	Status       string     `json:"status"`
	Reason       string     `json:"reason"`
	Note         string     `json:"note,omitempty"`
	Attempts     int        `json:"attempts"`
	ReturnFee    string     `json:"return_fee"`
	FeePayer     string     `json:"fee_payer"`
	PickupLabel  string     `json:"pickup_label"`
	PickupAddr   string     `json:"pickup_address"`
	PickupName   string     `json:"pickup_contact"`
	PickupPhone  string     `json:"pickup_phone"`
	CourierID    string     `json:"courier_id"`
	CourierName  string     `json:"courier_name"`
	CourierPhone string     `json:"courier_phone,omitempty"`
	StartedBy    string     `json:"started_by"` // courier · admin
	StartedAt    time.Time  `json:"started_at"`
	ReturnedAt   *time.Time `json:"returned_at,omitempty"`
	HasProof     bool       `json:"has_proof"`      // foto serah terima di toko
	HasFailProof bool       `json:"has_fail_proof"` // foto bukti gagal antar
	ConfirmBy    *time.Time `json:"confirm_by,omitempty"`
	ResolvedBy   string     `json:"resolved_by,omitempty"` // seller · admin · system
	ReceivedAt   *time.Time `json:"received_at,omitempty"`
	Condition    string     `json:"condition,omitempty"`
	Restocked    *bool      `json:"restocked,omitempty"`
	SellerNote   string     `json:"seller_note,omitempty"`
	Total        string     `json:"order_total"`
	Items        []Item     `json:"items"`
}

type Item struct {
	Name    string `json:"name"`
	Options string `json:"options"`
	Qty     int    `json:"qty"`
}

const retCols = `r.id, r.delivery_id, o.order_no, o.status::text, o.store_id, st.name, r.status, coalesce(r.reason, ''), coalesce(r.note, ''), r.attempts,
	r.return_fee::text, r.fee_payer, d.pickup_label, d.pickup_address, d.pickup_contact, d.pickup_phone, coalesce(d.courier_id::text, ''), coalesce(cr.display_name, ''),
	coalesce(cu.phone, ''), r.started_by_type::text, r.started_at, r.returned_at, r.returned_proof_key IS NOT NULL, d.fail_proof_key IS NOT NULL, r.confirm_by, coalesce(r.resolved_by_type::text, ''),
	r.received_at, coalesce(r.item_condition, ''), r.restocked, coalesce(r.seller_note, ''), o.total::text`

const retJoins = ` FROM courier_returns r JOIN deliveries d ON d.id = r.delivery_id JOIN orders o ON o.id = r.order_id JOIN stores st ON st.id = o.store_id
	LEFT JOIN couriers cr ON cr.user_id = d.courier_id LEFT JOIN users cu ON cu.id = d.courier_id `

func (s *Service) list(ctx context.Context, where, order string, args ...any) ([]Return, error) {
	rows, err := s.Pool.Query(ctx, `SELECT `+retCols+retJoins+` WHERE `+where+` ORDER BY `+order+` LIMIT 200`, args...)
	if err != nil {
		return nil, err
	}
	out := []Return{}
	for rows.Next() {
		var x Return
		if err := rows.Scan(&x.ID, &x.DeliveryID, &x.OrderNo, &x.OrderStatus, &x.StoreID, &x.StoreName, &x.Status, &x.Reason, &x.Note, &x.Attempts,
			&x.ReturnFee, &x.FeePayer, &x.PickupLabel, &x.PickupAddr, &x.PickupName, &x.PickupPhone, &x.CourierID, &x.CourierName,
			&x.CourierPhone, &x.StartedBy, &x.StartedAt, &x.ReturnedAt, &x.HasProof, &x.HasFailProof, &x.ConfirmBy, &x.ResolvedBy,
			&x.ReceivedAt, &x.Condition, &x.Restocked, &x.SellerNote, &x.Total); err != nil {
			rows.Close()
			return nil, err
		}
		x.Items = []Item{}
		out = append(out, x)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for i := range out {
		irows, err := s.Pool.Query(ctx, `SELECT oi.name, oi.options, oi.qty FROM order_items oi JOIN courier_returns r ON r.order_id = oi.order_id WHERE r.id = $1 ORDER BY oi.name`, out[i].ID)
		if err != nil {
			return nil, err
		}
		for irows.Next() {
			var it Item
			if err := irows.Scan(&it.Name, &it.Options, &it.Qty); err != nil {
				irows.Close()
				return nil, err
			}
			out[i].Items = append(out[i].Items, it)
		}
		irows.Close()
	}
	return out, nil
}

func (s *Service) get(ctx context.Context, where string, args ...any) (*Return, error) {
	l, err := s.list(ctx, where, "r.id", args...)
	if err != nil {
		return nil, err
	}
	if len(l) == 0 {
		return nil, ErrNotFound
	}
	return &l[0], nil
}

// ---------------------------------------------------------------------------
// Penguncian & langkah bersama
// ---------------------------------------------------------------------------

type lockedDelivery struct {
	id, orderID, orderNo, orderStatus, status, courierID, storeID, storeName, buyerID, lang, failReason, failNote, shippingFee string
	attempts                                                                                                                   int
}

func lockDelivery(ctx context.Context, tx pgx.Tx, deliveryID string) (*lockedDelivery, error) {
	if len(deliveryID) != 36 {
		return nil, ErrNotFound
	}
	var l lockedDelivery
	err := tx.QueryRow(ctx, `SELECT d.id, o.id, o.order_no, o.status::text, d.status, coalesce(d.courier_id::text, ''), o.store_id, st.name, o.buyer_id, c.lang::text,
		  coalesce(d.fail_reason, ''), coalesce(d.fail_note, ''), o.shipping_fee::text, d.attempts
		FROM deliveries d JOIN orders o ON o.id = d.order_id JOIN stores st ON st.id = o.store_id JOIN checkouts c ON c.id = o.checkout_id
		WHERE d.id = $1 FOR UPDATE OF d, o`, deliveryID).Scan(&l.id, &l.orderID, &l.orderNo, &l.orderStatus, &l.status, &l.courierID, &l.storeID, &l.storeName,
		&l.buyerID, &l.lang, &l.failReason, &l.failNote, &l.shippingFee, &l.attempts)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &l, err
}

type lockedReturn struct {
	lockedDelivery
	returnID, returnStatus, reason string
}

// lockReturn: kunci pengembalian + pengantaran + pesanan. storeID kosong = tanpa batas toko (admin / sistem).
func lockReturn(ctx context.Context, tx pgx.Tx, returnID, storeID string) (*lockedReturn, error) {
	if len(returnID) != 36 {
		return nil, ErrNotFound
	}
	var r lockedReturn
	var deliveryID string
	err := tx.QueryRow(ctx, `SELECT r.id, r.status, coalesce(r.reason, ''), r.delivery_id FROM courier_returns r
		WHERE r.id = $1 AND ($2 = '' OR r.store_id::text = $2) FOR UPDATE`, returnID, storeID).Scan(&r.returnID, &r.returnStatus, &r.reason, &deliveryID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	l, err := lockDelivery(ctx, tx, deliveryID)
	if err != nil {
		return nil, err
	}
	r.lockedDelivery = *l
	return &r, nil
}

func event(ctx context.Context, tx pgx.Tx, orderID, status, actorType, actorID, note string) error {
	_, err := tx.Exec(ctx, `INSERT INTO order_events (order_id, status, actor_type, actor_id, note) VALUES ($1, $2, $3::actor_type, nullif($4, '')::uuid, nullif($5, ''))`,
		orderID, status, actorType, actorID, trunc(note, 300))
	return err
}

func trunc(s string, n int) string {
	if r := []rune(s); len(r) > n {
		return string(r[:n])
	}
	return s
}

func activeCourierName(ctx context.Context, tx pgx.Tx, courierID string) (string, error) {
	var name string
	if len(courierID) != 36 {
		return "", fieldError{"courier_id", "choose an active courier"}
	}
	if err := tx.QueryRow(ctx, `SELECT display_name FROM couriers WHERE user_id = $1 AND active`, courierID).Scan(&name); err != nil {
		return "", fieldError{"courier_id", "choose an active courier"}
	}
	return name, nil
}

// start: pengantaran yang gagal mulai dikembalikan ke titik jemput toko (di dalam tx pemanggil).
func start(ctx context.Context, tx pgx.Tx, l *lockedDelivery, actorType, actorID, newCourierID string) error {
	if l.status != "failed" || l.orderStatus != "picked_up" {
		return ErrBadStatus
	}
	courier := l.courierID
	if newCourierID != "" && newCourierID != l.courierID {
		name, err := activeCourierName(ctx, tx, newCourierID)
		if err != nil {
			return err
		}
		courier = newCourierID
		if err := event(ctx, tx, l.orderID, "courier_assigned", actorType, actorID, name); err != nil {
			return err
		}
	}
	if _, err := tx.Exec(ctx, `UPDATE deliveries SET status = 'returning', courier_id = $2 WHERE id = $1`, l.id, courier); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `INSERT INTO courier_returns (delivery_id, order_id, store_id, reason, note, attempts, return_fee, started_by_type, started_by)
		VALUES ($1, $2, $3, nullif($4, ''), nullif($5, ''), $6, $7::numeric, $8::actor_type, nullif($9, '')::uuid)`,
		l.id, l.orderID, l.storeID, l.failReason, l.failNote, l.attempts, l.shippingFee, actorType, actorID); err != nil {
		return err
	}
	if err := event(ctx, tx, l.orderID, "return_to_store", actorType, actorID, l.failReason); err != nil {
		return err
	}
	if err := notify.Enqueue(ctx, tx, l.storeID, notify.CourierReturnStarted, notify.Params{Order: l.orderNo}, "#/store/"+l.storeID+"/orders"); err != nil {
		return err
	}
	if err := notify.EnqueueUser(ctx, tx, l.buyerID, notify.BuyerDeliveryReturning, notify.Params{Order: l.orderNo, Store: l.storeName}, "/"+l.lang+"/orders/"+l.orderNo); err != nil {
		return err
	}
	if courier != l.courierID || actorType == "admin" {
		return notify.EnqueueUser(ctx, tx, courier, notify.CourierReturnTask, notify.Params{Order: l.orderNo, Store: l.storeName}, "/courier/")
	}
	return nil
}

var cancelText = map[string]string{
	"en": "the courier could not deliver it",
	"km": "អ្នកដឹកមិនអាចដឹកជូនបាន",
	"zh": "快递员无法送达",
}

// receive: toko sudah menerima paket kembali → pesanan batal, refund menunggu admin, stok kembali kalau restock.
func receive(ctx context.Context, tx pgx.Tx, r *lockedReturn, actorType, actorID, condition, note string, restock bool) error {
	if r.returnStatus != "returning" && r.returnStatus != "awaiting_seller" {
		return ErrBadStatus
	}
	if r.orderStatus != "picked_up" {
		return ErrBadStatus
	}
	if _, err := tx.Exec(ctx, `UPDATE courier_returns SET status = 'received', resolved_by_type = $2::actor_type, received_by = nullif($3, '')::uuid, received_at = now(),
		item_condition = $4, restocked = $5, seller_note = nullif($6, ''), returned_at = coalesce(returned_at, now()) WHERE id = $1`,
		r.returnID, actorType, actorID, condition, restock, trunc(note, 300)); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `UPDATE deliveries SET status = 'returned' WHERE id = $1`, r.id); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `UPDATE orders SET status = 'cancelled', cancelled_at = now(), cancel_reason = 'delivery_failed',
		refund_status = 'pending', refund_amount = total WHERE id = $1`, r.orderID); err != nil {
		return err
	}
	if restock {
		rows, err := tx.Query(ctx, `SELECT id, variant_id, qty FROM order_items WHERE order_id = $1 ORDER BY variant_id`, r.orderID)
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
			if _, err := tx.Exec(ctx, `UPDATE product_variants SET stock_qty = stock_qty + $2 WHERE id = $1`, x.variant, x.qty); err != nil {
				return err
			}
			if _, err := tx.Exec(ctx, `INSERT INTO inventory_movements (variant_id, delta, reason, order_item_id, created_by) VALUES ($1, $2, 'return', $3, nullif($4, '')::uuid)`,
				x.variant, x.qty, x.id, actorID); err != nil {
				return err
			}
		}
	}
	evNote := condition
	if !restock {
		evNote += " no_restock"
	}
	if note != "" {
		evNote += ": " + note
	}
	if err := event(ctx, tx, r.orderID, "return_received", actorType, actorID, evNote); err != nil {
		return err
	}
	if err := event(ctx, tx, r.orderID, "cancelled", "system", "", "delivery_failed"); err != nil {
		return err
	}
	if err := event(ctx, tx, r.orderID, "refund_pending", "system", "", ""); err != nil {
		return err
	}
	var total string
	if err := tx.QueryRow(ctx, `SELECT total::text FROM orders WHERE id = $1`, r.orderID).Scan(&total); err != nil {
		return err
	}
	txt := cancelText[r.lang]
	if txt == "" {
		txt = cancelText["en"]
	}
	if err := notify.EnqueueUser(ctx, tx, r.buyerID, notify.BuyerOrderCancelled,
		notify.Params{Store: r.storeName, Order: r.orderNo, Note: txt, Product: "$" + total}, "/"+r.lang+"/orders/"+r.orderNo); err != nil {
		return err
	}
	if actorType != "seller" { // seller sendiri yang konfirmasi tidak perlu diberi tahu
		return notify.Enqueue(ctx, tx, r.storeID, notify.CourierReturnClosed, notify.Params{Order: r.orderNo}, "#/store/"+r.storeID+"/orders")
	}
	return nil
}

func checkCondition(condition, note string) (string, error) {
	note = strings.TrimSpace(note)
	if condition != "ok" && condition != "damaged" {
		return "", fieldError{"condition", "choose the condition of the package"}
	}
	if len([]rune(note)) > 300 {
		return "", fieldError{"note", "too long (max 300 characters)"}
	}
	return note, nil
}

// ---------------------------------------------------------------------------
// Kurir
// ---------------------------------------------------------------------------

func (s *Service) courier(ctx context.Context, userID string) (string, error) {
	var name string
	err := s.Pool.QueryRow(ctx, `SELECT display_name FROM couriers WHERE user_id = $1 AND active`, userID).Scan(&name)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNotCourier
	}
	return name, err
}

// CourierInfo: info pengembalian untuk tugas kurir (batas percobaan + pengembalian yang berjalan).
type CourierInfo struct {
	Status      string  `json:"status"`
	Attempts    int     `json:"attempts"`
	MaxAttempts int     `json:"max_attempts"`
	CanRetry    bool    `json:"can_retry"`
	CanReturn   bool    `json:"can_return"`
	Return      *Return `json:"return,omitempty"`
	PickupLat   float64 `json:"pickup_lat"`
	PickupLng   float64 `json:"pickup_lng"`
}

func (s *Service) CourierInfo(ctx context.Context, courierID, deliveryID string) (*CourierInfo, error) {
	if _, err := s.courier(ctx, courierID); err != nil {
		return nil, err
	}
	if len(deliveryID) != 36 {
		return nil, ErrNotFound
	}
	in := CourierInfo{MaxAttempts: delivery.MaxAttempts}
	var orderStatus string
	err := s.Pool.QueryRow(ctx, `SELECT d.status, d.attempts, d.pickup_lat::float8, d.pickup_lng::float8, o.status::text FROM deliveries d JOIN orders o ON o.id = d.order_id
		WHERE d.id = $1 AND d.courier_id = $2`, deliveryID, courierID).Scan(&in.Status, &in.Attempts, &in.PickupLat, &in.PickupLng, &orderStatus)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	in.CanRetry = in.Status == "failed" && in.Attempts < delivery.MaxAttempts
	in.CanReturn = in.Status == "failed" && orderStatus == "picked_up"
	if r, err := s.get(ctx, `r.delivery_id = $1`, deliveryID); err == nil {
		in.Return = r
	} else if !errors.Is(err, ErrNotFound) {
		return nil, err
	}
	return &in, nil
}

// CourierStart: kurir memutuskan membawa paket kembali ke toko (dari status gagal).
func (s *Service) CourierStart(ctx context.Context, courierID, deliveryID string) (*CourierInfo, error) {
	if _, err := s.courier(ctx, courierID); err != nil {
		return nil, err
	}
	err := db.WithActor(ctx, s.Pool, courierID, "courier", "", func(tx pgx.Tx) error {
		l, err := lockDelivery(ctx, tx, deliveryID)
		if err != nil {
			return err
		}
		if l.courierID != courierID {
			return ErrNotFound
		}
		if l.status == "returning" {
			return nil // klik dua kali
		}
		return start(ctx, tx, l, "courier", courierID, "")
	})
	if err != nil {
		return nil, err
	}
	return s.CourierInfo(ctx, courierID, deliveryID)
}

// CourierReturned: kurir sudah menyerahkan paket di titik jemput toko (foto opsional). Seller tinggal konfirmasi.
func (s *Service) CourierReturned(ctx context.Context, courierID, deliveryID, proofKey string) (*CourierInfo, error) {
	if _, err := s.courier(ctx, courierID); err != nil {
		return nil, err
	}
	if proofKey != "" {
		if err := storage.CheckOwner(ctx, s.Pool, proofKey, courierID, "delivery_proof"); err != nil {
			return nil, fieldError{"proof_key", "upload the photo again"}
		}
	}
	err := db.WithActor(ctx, s.Pool, courierID, "courier", "", func(tx pgx.Tx) error {
		l, err := lockDelivery(ctx, tx, deliveryID)
		if err != nil {
			return err
		}
		if l.courierID != courierID {
			return ErrNotFound
		}
		if l.status == "returned" {
			return nil
		}
		if l.status != "returning" {
			return ErrBadStatus
		}
		if _, err := tx.Exec(ctx, `UPDATE deliveries SET status = 'returned' WHERE id = $1`, l.id); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `UPDATE courier_returns SET status = 'awaiting_seller', returned_by = $2, returned_at = now(),
			returned_proof_key = nullif($3, ''), confirm_by = now() + make_interval(secs => $4) WHERE delivery_id = $1 AND status = 'returning'`,
			l.id, courierID, proofKey, AutoReceiveAfter.Seconds()); err != nil {
			return err
		}
		if err := event(ctx, tx, l.orderID, "returned_to_store", "courier", courierID, ""); err != nil {
			return err
		}
		return notify.Enqueue(ctx, tx, l.storeID, notify.CourierReturnArrived, notify.Params{Order: l.orderNo}, "#/store/"+l.storeID+"/orders")
	})
	if err != nil {
		return nil, err
	}
	return s.CourierInfo(ctx, courierID, deliveryID)
}

// ProofKey: foto serah terima di toko atau foto bukti gagal antar (untuk admin).
func (s *Service) ProofKey(ctx context.Context, returnID, which string) (string, error) {
	if len(returnID) != 36 {
		return "", ErrNotFound
	}
	col := "r.returned_proof_key"
	if which == "fail" {
		col = "d.fail_proof_key"
	}
	var key *string
	err := s.Pool.QueryRow(ctx, `SELECT `+col+` FROM courier_returns r JOIN deliveries d ON d.id = r.delivery_id WHERE r.id = $1`, returnID).Scan(&key)
	if errors.Is(err, pgx.ErrNoRows) || (err == nil && key == nil) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", err
	}
	return *key, nil
}

// ---------------------------------------------------------------------------
// Admin
// ---------------------------------------------------------------------------

type Counts struct {
	Returning      int `json:"returning"`
	AwaitingSeller int `json:"awaiting_seller"`
	Overdue        int `json:"overdue"`
}

// AdminList: status = open (berjalan) · done (diterima, 30 hari terakhir).
func (s *Service) AdminList(ctx context.Context, status string) ([]Return, Counts, error) {
	var c Counts
	if err := s.Pool.QueryRow(ctx, `SELECT count(*) FILTER (WHERE status = 'returning'), count(*) FILTER (WHERE status = 'awaiting_seller'),
		count(*) FILTER (WHERE status = 'awaiting_seller' AND confirm_by < now()) FROM courier_returns WHERE status <> 'received'`).Scan(&c.Returning, &c.AwaitingSeller, &c.Overdue); err != nil {
		return nil, c, err
	}
	var l []Return
	var err error
	if status == "done" {
		l, err = s.list(ctx, `r.status = 'received' AND r.received_at > now() - interval '30 days'`, `r.received_at DESC`)
	} else {
		l, err = s.list(ctx, `r.status IN ('returning','awaiting_seller')`, `r.status = 'awaiting_seller' DESC, r.started_at`)
	}
	return l, c, err
}

// AdminStart: admin memutuskan pengembalian untuk pengantaran yang gagal; courierID kosong = kurir yang sama.
func (s *Service) AdminStart(ctx context.Context, adminID, deliveryID, courierID string) (*Return, error) {
	err := db.WithActor(ctx, s.Pool, adminID, "admin", "", func(tx pgx.Tx) error {
		l, err := lockDelivery(ctx, tx, deliveryID)
		if err != nil {
			return err
		}
		return start(ctx, tx, l, "admin", adminID, courierID)
	})
	if err != nil {
		return nil, err
	}
	return s.get(ctx, `r.delivery_id = $1`, deliveryID)
}

// AdminAssign: tugaskan / pindahkan kurir pengembalian (hanya selama paket belum sampai di toko).
func (s *Service) AdminAssign(ctx context.Context, adminID, returnID, courierID string) (*Return, error) {
	err := db.WithActor(ctx, s.Pool, adminID, "admin", "", func(tx pgx.Tx) error {
		r, err := lockReturn(ctx, tx, returnID, "")
		if err != nil {
			return err
		}
		if r.returnStatus != "returning" {
			return ErrBadStatus
		}
		name, err := activeCourierName(ctx, tx, courierID)
		if err != nil {
			return err
		}
		if courierID == r.courierID {
			return nil
		}
		if _, err := tx.Exec(ctx, `UPDATE deliveries SET courier_id = $2 WHERE id = $1`, r.id, courierID); err != nil {
			return err
		}
		if err := event(ctx, tx, r.orderID, "courier_assigned", "admin", adminID, name); err != nil {
			return err
		}
		return notify.EnqueueUser(ctx, tx, courierID, notify.CourierReturnTask, notify.Params{Order: r.orderNo, Store: r.storeName}, "/courier/")
	})
	if err != nil {
		return nil, err
	}
	return s.get(ctx, `r.id = $1`, returnID)
}

// AdminResolve: admin menyelesaikan manual (mis. seller tidak menjawab, toko sudah memegang paket).
func (s *Service) AdminResolve(ctx context.Context, adminID, returnID string, restock bool, note string) (*Return, error) {
	note = strings.TrimSpace(note)
	if len([]rune(note)) > 300 {
		return nil, fieldError{"note", "too long (max 300 characters)"}
	}
	condition := "ok"
	if !restock {
		condition = "damaged"
	}
	err := db.WithActor(ctx, s.Pool, adminID, "admin", "", func(tx pgx.Tx) error {
		r, err := lockReturn(ctx, tx, returnID, "")
		if err != nil {
			return err
		}
		return receive(ctx, tx, r, "admin", adminID, condition, note, restock)
	})
	if err != nil {
		return nil, err
	}
	return s.get(ctx, `r.id = $1`, returnID)
}

// ---------------------------------------------------------------------------
// Seller
// ---------------------------------------------------------------------------

// SellerList: view = open (sedang kembali / menunggu konfirmasi) · all (riwayat 90 hari).
func (s *Service) SellerList(ctx context.Context, storeID, view string) ([]Return, error) {
	var l []Return
	var err error
	if view == "all" {
		l, err = s.list(ctx, `r.store_id = $1 AND r.created_at > now() - interval '90 days'`, `r.created_at DESC`, storeID)
	} else {
		l, err = s.list(ctx, `r.store_id = $1 AND r.status IN ('returning','awaiting_seller')`, `r.status = 'awaiting_seller' DESC, r.started_at`, storeID)
	}
	for i := range l {
		l[i].CourierPhone = "" // seller cukup nama kurir
	}
	return l, err
}

// SellerConfirm: seller menerima paket kembali. condition = ok (stok kembali) · damaged (stok tidak kembali).
func (s *Service) SellerConfirm(ctx context.Context, storeID, userID, returnID, condition, note string) (*Return, error) {
	note, err := checkCondition(condition, note)
	if err != nil {
		return nil, err
	}
	if condition == "damaged" && len([]rune(note)) < 3 {
		return nil, fieldError{"note", "describe the damage"}
	}
	err = db.WithActor(ctx, s.Pool, userID, "seller", "", func(tx pgx.Tx) error {
		r, err := lockReturn(ctx, tx, returnID, storeID)
		if err != nil {
			return err
		}
		return receive(ctx, tx, r, "seller", userID, condition, note, condition == "ok")
	})
	if err != nil {
		return nil, err
	}
	r, err := s.get(ctx, `r.id = $1`, returnID)
	if r != nil {
		r.CourierPhone = ""
	}
	return r, err
}

// ---------------------------------------------------------------------------
// Job latar: seller tidak mengonfirmasi → dianggap diterima
// ---------------------------------------------------------------------------

// AutoReceive: pengembalian "menunggu seller" yang lewat confirm_by → diterima otomatis.
// Stok kembali kecuali alasan gagal antar "paket rusak".
func (s *Service) AutoReceive(ctx context.Context) (int, error) {
	rows, err := s.Pool.Query(ctx, `SELECT id FROM courier_returns WHERE status = 'awaiting_seller' AND confirm_by < now() ORDER BY confirm_by LIMIT 50`)
	if err != nil {
		return 0, err
	}
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return 0, err
		}
		ids = append(ids, id)
	}
	rows.Close()
	n := 0
	for _, id := range ids {
		err := db.WithActor(ctx, s.Pool, "", "system", "", func(tx pgx.Tx) error {
			r, err := lockReturn(ctx, tx, id, "")
			if err != nil {
				return err
			}
			if r.returnStatus != "awaiting_seller" {
				return nil
			}
			restock := r.reason != "damaged"
			cond := "ok"
			if !restock {
				cond = "damaged"
			}
			return receive(ctx, tx, r, "system", "", cond, "", restock)
		})
		if errors.Is(err, ErrBadStatus) { // pesanan sudah diubah di tempat lain — biarkan admin yang menangani
			s.Log.Warn("pengembalian kurir: lewati terima otomatis", "id", id)
			continue
		}
		if err != nil {
			return n, err
		}
		n++
	}
	return n, nil
}

// Worker: jalankan AutoReceive tiap 5 menit (dicatat lewat jobs.Track).
func (s *Service) Worker(ctx context.Context) {
	t := time.NewTicker(5 * time.Minute)
	defer t.Stop()
	for {
		if n, err := jobs.Track(ctx, s.Pool, JobName, s.AutoReceive); err != nil && ctx.Err() == nil {
			s.Log.Error("pengembalian kurir: terima otomatis gagal", "err", err)
		} else if n > 0 {
			s.Log.Info("pengembalian kurir diterima otomatis (seller tidak konfirmasi)", "n", n)
		}
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}
	}
}

// JobName: nama job di jobs.Defs (halaman admin Scheduled jobs).
const JobName = "courierreturn.auto_receive"
