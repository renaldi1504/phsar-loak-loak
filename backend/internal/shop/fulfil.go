package shop

// Seller memproses pesanan & admin mencatat refund manual.
//
//	paid ──terima──► accepted ──siap──► ready_for_pickup ──(kurir, tahap berikut)──► picked_up ► delivered
//	  │                 │
//	  └─tolak / 2 jam tanpa respons──► cancelled + refund_status=pending ──admin──► refunded
//
// Seller tidak melihat nomor HP / alamat pembeli: pengantaran dilakukan kurir Phsar.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"phsar/internal/extras" // extras:
	"phsar/internal/hours"
	"phsar/internal/notify"
	"phsar/internal/paging"
	"phsar/internal/platform/db"
)

// Batas waktu terima & kemas diatur super admin (platform_settings 'order_timing', paket hours);
// hitungan bisa hanya berjalan saat toko buka.

var (
	ErrBadTransition = errors.New("order status does not allow this action")
	ErrNoPickupPoint = errors.New("store has no pickup point")
)

func addEvent(ctx context.Context, tx pgx.Tx, orderID, status, actorType, actorID, note string) error {
	_, err := tx.Exec(ctx, `INSERT INTO order_events (order_id, status, actor_type, actor_id, note) VALUES ($1, $2, $3::actor_type, nullif($4, '')::uuid, nullif($5, ''))`,
		orderID, status, actorType, actorID, truncRunes(note, 300))
	return err
}

// ---------------------------------------------------------------------------
// Tampilan untuk seller
// ---------------------------------------------------------------------------

type SellerOrder struct {
	ID               string       `json:"id"`
	OrderNo          string       `json:"order_no"`
	Status           string       `json:"status"`
	Subtotal         string       `json:"subtotal"`
	CommissionRate   string       `json:"commission_rate"`
	CommissionAmount string       `json:"commission_amount"`
	SellerAmount     string       `json:"seller_amount"`
	ItemCount        int          `json:"item_count"`
	BuyerName        string       `json:"buyer_name"` // nama depan saja
	PaidAt           *time.Time   `json:"paid_at,omitempty"`
	RespondBy        *time.Time   `json:"respond_by,omitempty"`
	ReadyBy          *time.Time   `json:"ready_by,omitempty"` // batas "siap dijemput" (kalau batas kemas aktif)
	AcceptedAt       *time.Time   `json:"accepted_at,omitempty"`
	ReadyAt          *time.Time   `json:"ready_at,omitempty"`
	CancelledAt      *time.Time   `json:"cancelled_at,omitempty"`
	CancelReason     string       `json:"cancel_reason,omitempty"`
	RejectReason     string       `json:"reject_reason,omitempty"`
	SellerNote       string       `json:"seller_note,omitempty"`
	PickupPoint      string       `json:"pickup_point,omitempty"`
	PickupPointID    string       `json:"pickup_point_id,omitempty"`
	Items            []SellerItem `json:"items"`
	Events           []OrderEvent `json:"events,omitempty"`
}

type SellerItem struct {
	ProductID string `json:"product_id"`
	Name      string `json:"name"`
	Options   string `json:"options"`
	SKU       string `json:"sku"`
	ImageURL  string `json:"image_url"`
	UnitPrice string `json:"unit_price"`
	Qty       int    `json:"qty"`
	LineTotal string `json:"line_total"`
}

type OrderEvent struct {
	Status string    `json:"status"`
	Actor  string    `json:"actor"`
	Note   string    `json:"note,omitempty"`
	At     time.Time `json:"at"`
}

// Filter tab di Seller Center.
var sellerFilters = map[string]string{
	"new":        `o.status = 'paid'`,
	"processing": `o.status = 'accepted'`,
	"ready":      `o.status = 'ready_for_pickup'`,
	"shipping":   `o.status = 'picked_up'`,
	"done":       `o.status IN ('delivered','completed')`,
	"cancelled":  `o.status IN ('cancelled','refunded') AND o.paid_at IS NOT NULL`,
}

// Pesanan yang tidak pernah dibayar tidak relevan untuk seller.
const sellerVisible = `o.paid_at IS NOT NULL`

func firstName(s string) string {
	f := strings.Fields(s)
	if len(f) == 0 {
		return ""
	}
	return f[0]
}

// Urutan default tab langkah (terima/kemas/serahkan): tenggat paling dekat dulu.
const sellerDeadlineOrder = `coalesce(CASE WHEN o.status = 'accepted' THEN o.ready_by END, o.respond_by, o.paid_at) ASC, o.order_no`

func (s *Service) sellerOrders(ctx context.Context, where, order string, limit, offset int, args ...any) ([]SellerOrder, error) {
	rows, err := s.Pool.Query(ctx, `
		SELECT o.id, o.order_no, o.status::text, o.subtotal::text, o.commission_rate::text, o.commission_amount::text, o.seller_amount::text,
		       o.item_count, c.recipient, o.paid_at, o.respond_by, o.ready_by, o.accepted_at, o.ready_at, o.cancelled_at, coalesce(o.cancel_reason, ''),
		       coalesce(o.reject_reason, ''), coalesce(o.seller_note, ''), coalesce(pp.label, ''), coalesce(o.pickup_point_id::text, '')
		FROM orders o JOIN checkouts c ON c.id = o.checkout_id LEFT JOIN store_pickup_points pp ON pp.id = o.pickup_point_id
		WHERE `+sellerVisible+` AND `+where+` ORDER BY `+order+` LIMIT `+fmt.Sprint(limit)+` OFFSET `+fmt.Sprint(offset), args...)
	if err != nil {
		return nil, err
	}
	out := []SellerOrder{}
	ids := []string{}
	for rows.Next() {
		var o SellerOrder
		var recipient string
		if err := rows.Scan(&o.ID, &o.OrderNo, &o.Status, &o.Subtotal, &o.CommissionRate, &o.CommissionAmount, &o.SellerAmount, &o.ItemCount,
			&recipient, &o.PaidAt, &o.RespondBy, &o.ReadyBy, &o.AcceptedAt, &o.ReadyAt, &o.CancelledAt, &o.CancelReason, &o.RejectReason, &o.SellerNote,
			&o.PickupPoint, &o.PickupPointID); err != nil {
			rows.Close()
			return nil, err
		}
		o.BuyerName = firstName(recipient)
		o.Items = []SellerItem{}
		out = append(out, o)
		ids = append(ids, o.ID)
	}
	rows.Close()
	if len(ids) == 0 {
		return out, rows.Err()
	}
	idx := map[string]int{}
	for i, id := range ids {
		idx[id] = i
	}
	rows, err = s.Pool.Query(ctx, `SELECT order_id, product_id, name, options, sku, coalesce(image_url, ''), unit_price::text, qty, line_total::text
		FROM order_items WHERE order_id = ANY($1::uuid[]) ORDER BY name, options`, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var oid string
		var it SellerItem
		if err := rows.Scan(&oid, &it.ProductID, &it.Name, &it.Options, &it.SKU, &it.ImageURL, &it.UnitPrice, &it.Qty, &it.LineTotal); err != nil {
			return nil, err
		}
		out[idx[oid]].Items = append(out[idx[oid]].Items, it)
	}
	return out, rows.Err()
}

type SellerOrderList struct {
	Orders []SellerOrder  `json:"orders"`
	Counts map[string]int `json:"counts"`
	Filter string         `json:"filter"`         // "search" saat mencari
	Page   *paging.Meta   `json:"page,omitempty"` // hanya kalau klien meminta halaman
}

// SellerOrderQuery: tab, kata kunci, rentang tanggal bayar [From, To), halaman.
type SellerOrderQuery struct {
	Filter   string
	Page     paging.Params
	From, To *time.Time
}

// Kata kunci pesanan: no. pesanan, nama penerima, nama produk / SKU. Nomor telepon sengaja tidak
// (seller hanya melihat nama depan pembeli). $2 = pola ILIKE.
const sellerOrderSearch = `(o.order_no ILIKE $2 OR c.recipient ILIKE $2
	OR EXISTS (SELECT 1 FROM order_items i WHERE i.order_id = o.id AND (i.name ILIKE $2 OR i.sku ILIKE $2)))`

func (s *Service) SellerListOrders(ctx context.Context, storeID string, q SellerOrderQuery) (*SellerOrderList, error) {
	filter := q.Filter
	where, ok := sellerFilters[filter]
	if !ok {
		filter, where = "new", sellerFilters["new"]
	}
	order := sellerDeadlineOrder
	if filter == "shipping" || filter == "done" || filter == "cancelled" { // riwayat: terbaru dulu
		order = "o.paid_at DESC, o.order_no DESC"
	}
	cond := `o.store_id = $1 AND `
	args := []any{storeID}
	if q.Page.Q != "" { // cari di semua status
		filter, where, order = "search", "true", "o.paid_at DESC, o.order_no DESC"
		cond += sellerOrderSearch + ` AND `
		args = append(args, paging.Like(q.Page.Q))
	}
	if q.From != nil {
		args = append(args, *q.From)
		cond += fmt.Sprintf(`o.paid_at >= $%d AND `, len(args))
	}
	if q.To != nil {
		args = append(args, *q.To)
		cond += fmt.Sprintf(`o.paid_at < $%d AND `, len(args))
	}
	limit, offset := 200, 0
	if q.Page.On {
		limit, offset = q.Page.Size, q.Page.Offset()
	}
	list, err := s.sellerOrders(ctx, cond+where, order, limit, offset, args...)
	if err != nil {
		return nil, err
	}
	out := &SellerOrderList{Orders: list, Counts: map[string]int{}, Filter: filter}
	if q.Page.On {
		var total int
		if err := s.Pool.QueryRow(ctx, `SELECT count(*) FROM orders o JOIN checkouts c ON c.id = o.checkout_id
			WHERE `+sellerVisible+` AND `+cond+where, args...).Scan(&total); err != nil {
			return nil, err
		}
		m := paging.NewMeta(q.Page, total)
		out.Page = &m
	}
	for k, w := range sellerFilters {
		if k == "done" || k == "cancelled" {
			continue
		}
		var n int
		if err := s.Pool.QueryRow(ctx, `SELECT count(*) FROM orders o WHERE o.store_id = $1 AND `+sellerVisible+` AND `+w, storeID).Scan(&n); err != nil {
			return nil, err
		}
		out.Counts[k] = n
	}
	return out, nil
}

func (s *Service) SellerGetOrder(ctx context.Context, storeID, orderNo string) (*SellerOrder, error) {
	list, err := s.sellerOrders(ctx, `o.store_id = $1 AND o.order_no = $2`, sellerDeadlineOrder, 1, 0, storeID, strings.ToUpper(strings.TrimSpace(orderNo)))
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, ErrNotFound
	}
	o := &list[0]
	o.Events, err = s.events(ctx, o.ID)
	return o, err
}

func (s *Service) events(ctx context.Context, orderID string) ([]OrderEvent, error) {
	rows, err := s.Pool.Query(ctx, `SELECT status, actor_type::text, coalesce(note, ''), created_at FROM order_events WHERE order_id = $1 ORDER BY id`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []OrderEvent{}
	for rows.Next() {
		var e OrderEvent
		if err := rows.Scan(&e.Status, &e.Actor, &e.Note, &e.At); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

// ---------------------------------------------------------------------------
// Aksi seller
// ---------------------------------------------------------------------------

type lockedOrder struct {
	id, no, status, storeID, storeName, buyerID, total, lang string
}

func lockOrder(ctx context.Context, tx pgx.Tx, storeID, orderNo string) (*lockedOrder, error) {
	var o lockedOrder
	err := tx.QueryRow(ctx, `SELECT o.id, o.order_no, o.status::text, o.store_id, st.name, o.buyer_id, o.total::text, c.lang::text
		FROM orders o JOIN stores st ON st.id = o.store_id JOIN checkouts c ON c.id = o.checkout_id
		WHERE o.store_id = $1 AND o.order_no = $2 AND o.paid_at IS NOT NULL FOR UPDATE OF o`, storeID, strings.ToUpper(strings.TrimSpace(orderNo))).
		Scan(&o.id, &o.no, &o.status, &o.storeID, &o.storeName, &o.buyerID, &o.total, &o.lang)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &o, err
}

func (s *Service) Accept(ctx context.Context, storeID, userID, orderNo string) (*SellerOrder, error) {
	err := db.WithActor(ctx, s.Pool, userID, "seller", "", func(tx pgx.Tx) error {
		o, err := lockOrder(ctx, tx, storeID, orderNo)
		if err != nil {
			return err
		}
		if o.status == "accepted" {
			return nil // klik dua kali
		}
		if o.status != "paid" {
			return ErrBadTransition
		}
		tm, err := hours.LoadTiming(ctx, tx)
		if err != nil {
			return err
		}
		var readyBy *time.Time
		if tm.PackEnabled {
			d, err := hours.StoreDeadline(ctx, tx, tm, storeID, time.Now(), tm.PackMinutes)
			if err != nil {
				return err
			}
			readyBy = &d
		}
		if _, err := tx.Exec(ctx, `UPDATE orders SET status = 'accepted', accepted_at = now(), handled_by = $2, ready_by = $3 WHERE id = $1`, o.id, userID, readyBy); err != nil {
			return err
		}
		if err := addEvent(ctx, tx, o.id, "accepted", "seller", userID, ""); err != nil {
			return err
		}
		return notify.EnqueueUser(ctx, tx, o.buyerID, notify.BuyerOrderAccepted, notify.Params{Store: o.storeName, Order: o.no}, "/"+o.lang+"/orders/"+o.no)
	})
	if err != nil {
		return nil, err
	}
	return s.SellerGetOrder(ctx, storeID, orderNo)
}

var rejectReasons = map[string]bool{"out_of_stock": true, "cannot_fulfil": true, "price_error": true, "other": true}

var reasonText = map[string]map[string]string{
	"en": {"out_of_stock": "out of stock", "cannot_fulfil": "the store cannot fulfil it", "price_error": "price error", "other": "cancelled by the store", "seller_timeout": "the store did not respond in time"},
	"km": {"out_of_stock": "អស់ស្តុក", "cannot_fulfil": "ហាងមិនអាចផ្គត់ផ្គង់បាន", "price_error": "តម្លៃខុស", "other": "ហាងបានបោះបង់", "seller_timeout": "ហាងមិនបានឆ្លើយតបទាន់ពេល"},
	"zh": {"out_of_stock": "缺货", "cannot_fulfil": "店铺无法履约", "price_error": "价格错误", "other": "店铺取消", "seller_timeout": "店铺未及时响应"},
}

func (s *Service) Reject(ctx context.Context, storeID, userID, orderNo, reason, note string) (*SellerOrder, error) {
	if !rejectReasons[reason] {
		return nil, invalid("reason", "choose a reason")
	}
	note = strings.TrimSpace(note)
	if reason == "other" && len([]rune(note)) < 3 {
		return nil, invalid("note", "explain why")
	}
	err := db.WithActor(ctx, s.Pool, userID, "seller", "", func(tx pgx.Tx) error {
		o, err := lockOrder(ctx, tx, storeID, orderNo)
		if err != nil {
			return err
		}
		if o.status != "paid" && o.status != "accepted" {
			return ErrBadTransition
		}
		return s.cancelPaid(ctx, tx, o, "seller_rejected", reason, note, "seller", userID)
	})
	if err != nil {
		return nil, err
	}
	return s.SellerGetOrder(ctx, storeID, orderNo)
}

// cancelPaid: pesanan yang sudah dibayar dibatalkan → stok kembali, refund menunggu admin, pembeli diberi tahu.
func (s *Service) cancelPaid(ctx context.Context, tx pgx.Tx, o *lockedOrder, cancelReason, rejectReason, note, actorType, actorID string) error {
	if _, err := tx.Exec(ctx, `UPDATE orders SET status = 'cancelled', cancelled_at = now(), cancel_reason = $2, reject_reason = nullif($3, ''),
		seller_note = nullif($4, ''), handled_by = nullif($5, '')::uuid, refund_status = 'pending', refund_amount = total WHERE id = $1`,
		o.id, cancelReason, rejectReason, truncRunes(note, 300), actorID); err != nil {
		return err
	}
	rows, err := tx.Query(ctx, `SELECT id, variant_id, qty FROM order_items WHERE order_id = $1 ORDER BY variant_id`, o.id)
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
	// Barang tidak jadi terjual → stok kembali (kecuali alasan "stok habis": stok seller memang salah, jangan dikembalikan).
	if rejectReason != "out_of_stock" {
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
	reason := rejectReason
	if reason == "" {
		reason = cancelReason
	}
	if err := addEvent(ctx, tx, o.id, "cancelled", actorType, actorID, reason+trimNote(note)); err != nil {
		return err
	}
	if err := addEvent(ctx, tx, o.id, "refund_pending", "system", "", ""); err != nil {
		return err
	}
	txt := reasonText[o.lang][reason]
	if txt == "" {
		txt = reasonText["en"]["other"]
	}
	return notify.EnqueueUser(ctx, tx, o.buyerID, notify.BuyerOrderCancelled,
		notify.Params{Store: o.storeName, Order: o.no, Note: txt, Product: "$" + o.total}, "/"+o.lang+"/orders/"+o.no)
}

func trimNote(n string) string {
	if n == "" {
		return ""
	}
	return ": " + n
}

func (s *Service) Ready(ctx context.Context, storeID, userID, orderNo, pickupPointID string) (*SellerOrder, error) {
	err := db.WithActor(ctx, s.Pool, userID, "seller", "", func(tx pgx.Tx) error {
		o, err := lockOrder(ctx, tx, storeID, orderNo)
		if err != nil {
			return err
		}
		if o.status == "ready_for_pickup" {
			return nil
		}
		if o.status != "accepted" {
			return ErrBadTransition
		}
		// Titik penjemputan: yang dipilih, atau default toko (yang sudah diverifikasi didahulukan).
		var pp string
		err = tx.QueryRow(ctx, `SELECT id FROM store_pickup_points WHERE store_id = $1 AND deleted_at IS NULL AND ($2 = '' OR id::text = $2)
			ORDER BY (verified_at IS NOT NULL) DESC, is_default DESC, created_at LIMIT 1`, storeID, pickupPointID).Scan(&pp)
		if errors.Is(err, pgx.ErrNoRows) {
			if pickupPointID != "" {
				return invalid("pickup_point_id", "pickup point not found")
			}
			return ErrNoPickupPoint
		}
		if err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `UPDATE orders SET status = 'ready_for_pickup', ready_at = now(), pickup_point_id = $2, handled_by = $3 WHERE id = $1`,
			o.id, pp, userID); err != nil {
			return err
		}
		// Tugas untuk kurir Phsar (salinan titik jemput + titik antar).
		if _, err := tx.Exec(ctx, `
			INSERT INTO deliveries (order_id, pickup_label, pickup_address, pickup_lat, pickup_lng, pickup_contact, pickup_phone, dropoff_lat, dropoff_lng)
			SELECT o.id, pp.label, pp.address_line, pp.lat, pp.lng, pp.contact_name, pp.contact_phone, c.lat, c.lng
			FROM orders o JOIN store_pickup_points pp ON pp.id = $2 JOIN checkouts c ON c.id = o.checkout_id WHERE o.id = $1
			ON CONFLICT (order_id) DO NOTHING`, o.id, pp); err != nil {
			return err
		}
		if err := addEvent(ctx, tx, o.id, "ready_for_pickup", "seller", userID, ""); err != nil {
			return err
		}
		// Kabari kurir yang sedang online (Telegram).
		rows, err := tx.Query(ctx, `SELECT user_id FROM couriers WHERE active AND is_online`)
		if err != nil {
			return err
		}
		var couriers []string
		for rows.Next() {
			var id string
			if err := rows.Scan(&id); err != nil {
				rows.Close()
				return err
			}
			couriers = append(couriers, id)
		}
		rows.Close()
		for _, c := range couriers {
			if err := notify.EnqueueUser(ctx, tx, c, notify.CourierNewTask, notify.Params{Store: o.storeName, Order: o.no}, "/courier/"); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return s.SellerGetOrder(ctx, storeID, orderNo)
}

// ExpireUnaccepted (worker): pesanan dibayar yang tidak diterima seller sebelum respond_by → batal + refund.
func (s *Service) ExpireUnaccepted(ctx context.Context) (int, error) {
	return s.expireLate(ctx, "paid", "respond_by", "seller_timeout", notify.OrderSellerTimeout)
}

// ExpireUnpacked (worker): pesanan diterima tapi tidak "siap dijemput" sebelum ready_by → batal + refund.
func (s *Service) ExpireUnpacked(ctx context.Context) (int, error) {
	return s.expireLate(ctx, "accepted", "ready_by", "pack_timeout", notify.OrderPackTimeout)
}

func (s *Service) expireLate(ctx context.Context, status, col, reason, kind string) (int, error) {
	rows, err := s.Pool.Query(ctx, `SELECT store_id, order_no FROM orders WHERE status = $1 AND `+col+` < now() ORDER BY `+col+` LIMIT 50`, status)
	if err != nil {
		return 0, err
	}
	type ref struct{ store, no string }
	var list []ref
	for rows.Next() {
		var r ref
		if err := rows.Scan(&r.store, &r.no); err != nil {
			rows.Close()
			return 0, err
		}
		list = append(list, r)
	}
	rows.Close()
	n := 0
	for _, r := range list {
		err := db.WithActor(ctx, s.Pool, "", "system", "", func(tx pgx.Tx) error {
			o, err := lockOrder(ctx, tx, r.store, r.no)
			if err != nil {
				return err
			}
			if o.status != status {
				return nil
			}
			if err := s.cancelPaid(ctx, tx, o, reason, "", "", "system", ""); err != nil {
				return err
			}
			extras.AutoStrike(ctx, tx, o.storeID, o.id, o.no, reason) // extras: strike otomatis (kalau dinyalakan super admin)
			return notify.Enqueue(ctx, tx, o.storeID, kind, notify.Params{Order: o.no}, "#/store/"+o.storeID+"/orders")
		})
		if err != nil {
			return n, err
		}
		n++
	}
	return n, nil
}

// ---------------------------------------------------------------------------
// Refund manual (admin)
// ---------------------------------------------------------------------------

type RefundRow struct {
	OrderID      string     `json:"order_id"`
	OrderNo      string     `json:"order_no"`
	StoreName    string     `json:"store_name"`
	Amount       string     `json:"amount"`
	Status       string     `json:"refund_status"`
	CancelReason string     `json:"cancel_reason"`
	RejectReason string     `json:"reject_reason,omitempty"`
	CancelledAt  *time.Time `json:"cancelled_at"`
	BuyerName    string     `json:"buyer_name"`
	BuyerPhone   string     `json:"buyer_phone"`
	TranID       string     `json:"tran_id"`     // transaksi PayWay asli (cari di dashboard PayWay)
	PaymentRef   string     `json:"payment_ref"` // APV
	PaidAmount   string     `json:"checkout_total"`
	RefundedAt   *time.Time `json:"refunded_at,omitempty"`
	RefundRef    string     `json:"refund_ref,omitempty"`
	RefundNote   string     `json:"refund_note,omitempty"`
}

func (s *Service) ListRefunds(ctx context.Context, status string) ([]RefundRow, error) {
	if status != "refunded" {
		status = "pending"
	}
	rows, err := s.Pool.Query(ctx, `
		SELECT o.id, o.order_no, st.name, coalesce(o.refund_amount, o.total)::text, o.refund_status, coalesce(o.cancel_reason, ''), coalesce(o.reject_reason, ''),
		       o.cancelled_at, c.recipient, c.phone, c.tran_id, coalesce(c.payment_ref, ''), c.total::text, o.refunded_at, coalesce(o.refund_ref, ''), coalesce(o.refund_note, '')
		FROM orders o JOIN stores st ON st.id = o.store_id JOIN checkouts c ON c.id = o.checkout_id
		WHERE o.refund_status = $1 ORDER BY o.cancelled_at `+map[string]string{"pending": "ASC", "refunded": "DESC"}[status]+` LIMIT 200`, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []RefundRow{}
	for rows.Next() {
		var r RefundRow
		if err := rows.Scan(&r.OrderID, &r.OrderNo, &r.StoreName, &r.Amount, &r.Status, &r.CancelReason, &r.RejectReason, &r.CancelledAt,
			&r.BuyerName, &r.BuyerPhone, &r.TranID, &r.PaymentRef, &r.PaidAmount, &r.RefundedAt, &r.RefundRef, &r.RefundNote); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// MarkRefunded: admin sudah mengembalikan uang dari dashboard PayWay → catat nomor referensinya.
func (s *Service) MarkRefunded(ctx context.Context, adminID, orderID, ref, note string) error {
	ref, note = strings.TrimSpace(ref), strings.TrimSpace(note)
	if len(ref) < 3 || len(ref) > 80 {
		return invalid("refund_ref", "enter the refund reference from PayWay")
	}
	if len(orderID) != 36 {
		return ErrNotFound
	}
	return db.WithActor(ctx, s.Pool, adminID, "admin", "", func(tx pgx.Tx) error {
		var status, no, buyer, lang, amount string
		err := tx.QueryRow(ctx, `SELECT o.refund_status, o.order_no, o.buyer_id, c.lang::text, coalesce(o.refund_amount, o.total)::text
			FROM orders o JOIN checkouts c ON c.id = o.checkout_id WHERE o.id = $1 FOR UPDATE OF o`, orderID).Scan(&status, &no, &buyer, &lang, &amount)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		if status != "pending" {
			return ErrBadTransition
		}
		if _, err := tx.Exec(ctx, `UPDATE orders SET refund_status = 'refunded', refunded_at = now(), refund_ref = $2, refund_note = nullif($3, ''),
			refunded_by = $4, status = 'refunded' WHERE id = $1`, orderID, ref, truncRunes(note, 300), adminID); err != nil {
			return err
		}
		if err := addEvent(ctx, tx, orderID, "refunded", "admin", adminID, ref); err != nil {
			return err
		}
		return notify.EnqueueUser(ctx, tx, buyer, notify.BuyerOrderRefunded, notify.Params{Order: no, Product: "$" + amount}, "/"+lang+"/orders/"+no)
	})
}
