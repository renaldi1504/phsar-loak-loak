// Package returns: retur & sengketa.
//
//	pembeli ajukan (alasan, keterangan, foto, barang) ──► requested ──seller setuju (48 jam)──► approved
//	                                                       │  └─seller tolak / tidak menjawab──► disputed ──admin──► approved | rejected
//	                                                       └─pembeli batalkan──► cancelled
//	approved ──admin mengembalikan uang lewat PayWay lalu mencatat referensinya──► refunded (saldo seller dipotong bagiannya)
//
// Batas waktu (sesuai Kebijakan Retur): berubah pikiran = return_days produk sejak diterima;
// salah kirim / rusak / tidak sesuai / kurang lengkap = max(return_days, 3 hari), juga untuk produk "tanpa retur".
// Selama retur terbuka, pesanan tidak masuk saldo seller (paket payout melewatinya).
package returns

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"phsar/internal/jobs"
	"phsar/internal/notify"
	"phsar/internal/paging"
	"phsar/internal/platform/db"
	"phsar/internal/storage"
)

const (
	SellerHours   = 48 // waktu seller menjawab
	MinDefectDays = 3  // barang bermasalah: minimal 3 hari sejak diterima
	MaxPhotos     = 6
)

var (
	ErrNotFound   = errors.New("retur tidak ditemukan")
	ErrNotAllowed = errors.New("pesanan ini tidak bisa diretur")
	ErrState      = errors.New("retur sudah diproses")
	ErrDeadline   = errors.New("batas waktu sudah lewat")
)

// FieldError: input tidak valid (422).
type FieldError struct{ Field, Msg string }

func (e FieldError) Error() string { return e.Field + ": " + e.Msg }

var Reasons = []string{"changed_mind", "wrong_item", "damaged", "not_as_described", "missing_parts", "other"}

func validReason(r string) bool {
	for _, x := range Reasons {
		if x == r {
			return true
		}
	}
	return false
}

// defect: alasan "barang bermasalah" — berlaku minimal 3 hari, juga untuk produk tanpa retur.
func defect(reason string) bool { return reason != "changed_mind" }

type Service struct {
	Pool  *pgxpool.Pool
	Files storage.Storage
	Log   *slog.Logger
}

// ---------------------------------------------------------------------------
// Model
// ---------------------------------------------------------------------------

type Item struct {
	OrderItemID string `json:"order_item_id"`
	Name        string `json:"name"`
	Options     string `json:"options"`
	ImageURL    string `json:"image_url,omitempty"`
	UnitPrice   string `json:"unit_price"`
	Qty         int    `json:"qty"`        // dibeli
	ReturnQty   int    `json:"return_qty"` // diretur (0 = tidak ikut)
	ReturnDays  int    `json:"return_days"`
}

type Return struct {
	ID              string     `json:"id"`
	No              string     `json:"return_no"`
	OrderID         string     `json:"order_id,omitempty"`
	OrderNo         string     `json:"order_no"`
	OrderTotal      string     `json:"order_total"`
	StoreID         string     `json:"store_id"`
	StoreName       string     `json:"store_name"`
	BuyerName       string     `json:"buyer_name,omitempty"`  // hanya admin
	BuyerPhone      string     `json:"buyer_phone,omitempty"` // hanya admin (kurir menjemput barang)
	Status          string     `json:"status"`
	Reason          string     `json:"reason"`
	Description     string     `json:"description"`
	PhotoCount      int        `json:"photo_count"`
	Items           []Item     `json:"items"`
	RequestedAmount string     `json:"requested_amount"`
	RefundAmount    *string    `json:"refund_amount"`
	SellerNote      *string    `json:"seller_note"`
	SellerDeadline  time.Time  `json:"seller_deadline"`
	SellerAt        *time.Time `json:"seller_responded_at"`
	Escalated       *string    `json:"escalated_reason"`
	DecisionNote    *string    `json:"decision_note"`
	DecidedAt       *time.Time `json:"decided_at"`
	NeedsPickup     bool       `json:"needs_pickup"`
	RefundRef       *string    `json:"refund_ref,omitempty"` // hanya admin
	RefundedAt      *time.Time `json:"refunded_at"`
	DeliveredAt     *time.Time `json:"delivered_at"`
	CreatedAt       time.Time  `json:"created_at"`
}

const selectReturn = `
	SELECT r.id, r.return_no, o.id, o.order_no, o.total::text, r.store_id, st.name, coalesce(u.name, ''), coalesce(u.phone, ''), r.status::text, r.reason,
	  r.description, cardinality(r.photo_keys), r.items, r.requested_amount::text, r.refund_amount::text, r.seller_note, r.seller_deadline,
	  r.seller_responded_at, r.escalated_reason, r.decision_note, r.decided_at, r.needs_pickup, r.refund_ref, r.refunded_at,
	  d.delivered_at, r.created_at
	FROM returns r JOIN orders o ON o.id = r.order_id JOIN stores st ON st.id = r.store_id
	LEFT JOIN users u ON u.id = r.buyer_id LEFT JOIN deliveries d ON d.order_id = o.id`

func (s *Service) query(ctx context.Context, q interface {
	Query(context.Context, string, ...any) (pgx.Rows, error)
}, where string, limit int, args ...any) ([]Return, error) {
	sql := selectReturn + " WHERE " + where
	if limit > 0 { // limit 0: where sudah membawa ORDER BY/LIMIT sendiri
		sql += " ORDER BY r.created_at DESC LIMIT " + fmt.Sprint(limit)
	}
	rows, err := q.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Return{}
	for rows.Next() {
		var r Return
		var items []byte
		if err := rows.Scan(&r.ID, &r.No, &r.OrderID, &r.OrderNo, &r.OrderTotal, &r.StoreID, &r.StoreName, &r.BuyerName, &r.BuyerPhone, &r.Status, &r.Reason,
			&r.Description, &r.PhotoCount, &items, &r.RequestedAmount, &r.RefundAmount, &r.SellerNote, &r.SellerDeadline,
			&r.SellerAt, &r.Escalated, &r.DecisionNote, &r.DecidedAt, &r.NeedsPickup, &r.RefundRef, &r.RefundedAt,
			&r.DeliveredAt, &r.CreatedAt); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(items, &r.Items)
		if r.Items == nil {
			r.Items = []Item{}
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

func (s *Service) one(ctx context.Context, where string, args ...any) (*Return, error) {
	l, err := s.query(ctx, s.Pool, where, 1, args...)
	if err != nil {
		return nil, err
	}
	if len(l) == 0 {
		return nil, ErrNotFound
	}
	return &l[0], nil
}

// forBuyer / forSeller menyembunyikan data yang bukan haknya.
func forBuyer(r *Return) *Return {
	r.BuyerName, r.BuyerPhone, r.RefundRef, r.OrderID = "", "", nil, ""
	return r
}
func forSeller(r *Return) *Return { r.BuyerPhone, r.RefundRef, r.OrderID = "", nil, ""; return r }

// ---------------------------------------------------------------------------
// Pembeli
// ---------------------------------------------------------------------------

// Eligibility: apa yang boleh dilakukan pembeli untuk pesanan ini.
type Eligibility struct {
	CanRequest     bool       `json:"can_request"`
	NotAllowed     string     `json:"not_allowed,omitempty"` // not_delivered · already_requested · window_closed
	ChangeMindOK   bool       `json:"change_mind_ok"`        // masih dalam masa retur "berubah pikiran" untuk minimal satu barang
	DefectDeadline *time.Time `json:"defect_deadline,omitempty"`
	Items          []Item     `json:"items"`
	Return         *Return    `json:"return"` // retur terakhir (kalau ada)
}

type orderInfo struct {
	id, no, status, storeID, storeName, total, buyerLang string
	deliveredAt                                          *time.Time
	items                                                []Item
}

func (s *Service) loadOrder(ctx context.Context, q interface {
	QueryRow(context.Context, string, ...any) pgx.Row
	Query(context.Context, string, ...any) (pgx.Rows, error)
}, buyerID, orderNo string, lock bool) (*orderInfo, error) {
	o := &orderInfo{}
	sfx := ""
	if lock {
		sfx = " FOR UPDATE OF o"
	}
	err := q.QueryRow(ctx, `SELECT o.id, o.order_no, o.status::text, o.store_id, st.name, o.total::text, coalesce(u.language::text, 'km'), d.delivered_at
		FROM orders o JOIN stores st ON st.id = o.store_id JOIN users u ON u.id = o.buyer_id LEFT JOIN deliveries d ON d.order_id = o.id
		WHERE o.buyer_id = $1 AND o.order_no = $2`+sfx, buyerID, strings.ToUpper(strings.TrimSpace(orderNo))).
		Scan(&o.id, &o.no, &o.status, &o.storeID, &o.storeName, &o.total, &o.buyerLang, &o.deliveredAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	rows, err := q.Query(ctx, `SELECT oi.id, oi.name, oi.options, coalesce(oi.image_url, ''), oi.unit_price::text, oi.qty, p.return_days
		FROM order_items oi JOIN products p ON p.id = oi.product_id WHERE oi.order_id = $1 ORDER BY oi.id`, o.id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var it Item
		if err := rows.Scan(&it.OrderItemID, &it.Name, &it.Options, &it.ImageURL, &it.UnitPrice, &it.Qty, &it.ReturnDays); err != nil {
			return nil, err
		}
		o.items = append(o.items, it)
	}
	return o, rows.Err()
}

func deadline(delivered time.Time, days int, reason string) time.Time {
	if defect(reason) && days < MinDefectDays {
		days = MinDefectDays
	}
	return delivered.Add(time.Duration(days) * 24 * time.Hour)
}

func (s *Service) Eligibility(ctx context.Context, buyerID, orderNo string) (*Eligibility, error) {
	o, err := s.loadOrder(ctx, s.Pool, buyerID, orderNo, false)
	if err != nil {
		return nil, err
	}
	el := &Eligibility{Items: o.items}
	if el.Items == nil {
		el.Items = []Item{}
	}
	last, err := s.query(ctx, s.Pool, "r.order_id = $1", 1, o.id)
	if err != nil {
		return nil, err
	}
	if len(last) > 0 {
		el.Return = forBuyer(&last[0])
	}
	now := time.Now()
	switch {
	case (o.status != "delivered" && o.status != "completed") || o.deliveredAt == nil:
		el.NotAllowed = "not_delivered"
	case el.Return != nil && el.Return.Status != "cancelled":
		el.NotAllowed = "already_requested"
	default:
		maxDefect := time.Time{}
		for _, it := range o.items {
			if deadline(*o.deliveredAt, it.ReturnDays, "changed_mind").After(now) && it.ReturnDays > 0 {
				el.ChangeMindOK = true
			}
			if d := deadline(*o.deliveredAt, it.ReturnDays, "damaged"); d.After(maxDefect) {
				maxDefect = d
			}
		}
		if maxDefect.After(now) {
			el.CanRequest, el.DefectDeadline = true, &maxDefect
		} else {
			el.NotAllowed = "window_closed"
		}
	}
	return el, nil
}

type RequestInput struct {
	Reason      string   `json:"reason"`
	Description string   `json:"description"`
	PhotoKeys   []string `json:"photo_keys"`
	Items       []struct {
		OrderItemID string `json:"order_item_id"`
		Qty         int    `json:"qty"`
	} `json:"items"` // kosong = semua barang
}

func (s *Service) Request(ctx context.Context, buyerID, orderNo string, in RequestInput) (*Return, error) {
	in.Description = strings.TrimSpace(in.Description)
	if !validReason(in.Reason) {
		return nil, FieldError{"reason", "choose a reason"}
	}
	if n := len([]rune(in.Description)); n < 10 || n > 2000 {
		return nil, FieldError{"description", "describe the problem (10–2000 characters)"}
	}
	if len(in.PhotoKeys) > MaxPhotos {
		return nil, FieldError{"photo_keys", fmt.Sprintf("at most %d photos", MaxPhotos)}
	}
	for _, k := range in.PhotoKeys {
		if err := storage.CheckOwner(ctx, s.Pool, k, buyerID, "return_photo"); err != nil {
			return nil, FieldError{"photo_keys", "photo not found"}
		}
	}
	var id string
	err := db.WithActor(ctx, s.Pool, buyerID, "buyer", "", func(tx pgx.Tx) error {
		o, err := s.loadOrder(ctx, tx, buyerID, orderNo, true)
		if err != nil {
			return err
		}
		if (o.status != "delivered" && o.status != "completed") || o.deliveredAt == nil {
			return ErrNotAllowed
		}
		var open bool
		if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM returns WHERE order_id = $1 AND status <> 'cancelled')`, o.id).Scan(&open); err != nil {
			return err
		}
		if open {
			return ErrState
		}
		want := map[string]int{}
		for _, x := range in.Items {
			if x.Qty > 0 {
				want[x.OrderItemID] += x.Qty
			}
		}
		now, amount, chosen := time.Now(), decimal.Zero, []map[string]any{}
		for i, it := range o.items {
			q := it.Qty
			if len(want) > 0 {
				q = want[it.OrderItemID]
				delete(want, it.OrderItemID)
			}
			if q <= 0 {
				continue
			}
			if q > it.Qty {
				return FieldError{"items", "quantity is more than you bought"}
			}
			if !defect(in.Reason) && it.ReturnDays == 0 {
				return FieldError{"reason", "this item cannot be returned for change of mind"}
			}
			if !deadline(*o.deliveredAt, it.ReturnDays, in.Reason).After(now) {
				return FieldError{"reason", "the return window for this item has closed"}
			}
			price, _ := decimal.NewFromString(it.UnitPrice)
			amount = amount.Add(price.Mul(decimal.NewFromInt(int64(q))))
			o.items[i].ReturnQty = q
			chosen = append(chosen, map[string]any{"order_item_id": it.OrderItemID, "qty": q})
		}
		if len(want) > 0 {
			return FieldError{"items", "item is not in this order"}
		}
		if len(chosen) == 0 {
			return FieldError{"items", "choose at least one item"}
		}
		total, _ := decimal.NewFromString(o.total)
		if amount.GreaterThan(total) {
			amount = total
		}
		snapshot, _ := json.Marshal(o.items)
		_ = chosen
		photos := in.PhotoKeys
		if photos == nil {
			photos = []string{}
		}
		no := "RT" + time.Now().In(phnomPenh).Format("060102") + "-" + randCode(5)
		if err := tx.QueryRow(ctx, `INSERT INTO returns (return_no, order_id, buyer_id, store_id, reason, description, photo_keys, items,
			requested_amount, seller_deadline) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, now() + make_interval(hours => $10)) RETURNING id`,
			no, o.id, buyerID, o.storeID, in.Reason, in.Description, photos, snapshot, amount, SellerHours).Scan(&id); err != nil {
			return err
		}
		if err := event(ctx, tx, o.id, "return_requested", "buyer", buyerID, no); err != nil {
			return err
		}
		return notify.Enqueue(ctx, tx, o.storeID, notify.ReturnRequested, notify.Params{Order: o.no, Note: "$" + amount.StringFixed(2)}, "#/store/"+o.storeID+"/returns")
	})
	if err != nil {
		return nil, err
	}
	r, err := s.one(ctx, "r.id = $1", id)
	if err != nil {
		return nil, err
	}
	return forBuyer(r), nil
}

func (s *Service) BuyerList(ctx context.Context, buyerID string) ([]Return, error) {
	l, err := s.query(ctx, s.Pool, "r.buyer_id = $1", 100, buyerID)
	for i := range l {
		forBuyer(&l[i])
	}
	return l, err
}

func (s *Service) BuyerCancel(ctx context.Context, buyerID, id string) (*Return, error) {
	err := db.WithActor(ctx, s.Pool, buyerID, "buyer", "", func(tx pgx.Tx) error {
		var orderID, no string
		err := tx.QueryRow(ctx, `UPDATE returns SET status = 'cancelled' WHERE id = $1 AND buyer_id = $2 AND status IN ('requested','disputed')
			RETURNING order_id, return_no`, id, buyerID).Scan(&orderID, &no)
		if errors.Is(err, pgx.ErrNoRows) {
			return s.stateOrMissing(ctx, tx, "id = $1 AND buyer_id = $2", id, buyerID)
		}
		if err != nil {
			return err
		}
		return event(ctx, tx, orderID, "return_cancelled", "buyer", buyerID, no)
	})
	if err != nil {
		return nil, err
	}
	r, err := s.one(ctx, "r.id = $1", id)
	if err != nil {
		return nil, err
	}
	return forBuyer(r), nil
}

// stateOrMissing: UPDATE tidak kena baris → bedakan "tidak ada" (404) dan "status sudah berubah" (409).
func (s *Service) stateOrMissing(ctx context.Context, tx pgx.Tx, where string, args ...any) error {
	var ok bool
	if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM returns WHERE `+where+`)`, args...).Scan(&ok); err != nil {
		return err
	}
	if !ok {
		return ErrNotFound
	}
	return ErrState
}

// Photo membuka foto ke-i (dari 0) milik retur yang boleh dilihat (where membatasi pemilik).
func (s *Service) Photo(ctx context.Context, where string, i int, args ...any) (io.ReadCloser, string, error) {
	var keys []string
	if err := s.Pool.QueryRow(ctx, `SELECT photo_keys FROM returns r WHERE `+where, args...).Scan(&keys); err != nil {
		return nil, "", ErrNotFound
	}
	if i < 0 || i >= len(keys) {
		return nil, "", ErrNotFound
	}
	var ctype string
	if err := s.Pool.QueryRow(ctx, `SELECT content_type FROM uploads WHERE key = $1`, keys[i]).Scan(&ctype); err != nil {
		return nil, "", ErrNotFound
	}
	f, err := s.Files.Open(ctx, keys[i])
	if err != nil {
		return nil, "", ErrNotFound
	}
	return f, ctype, nil
}

// ---------------------------------------------------------------------------
// Seller
// ---------------------------------------------------------------------------

// SellerList: filter = open (default) · all; q = no. retur / no. pesanan / nama pembeli.
// Tanpa pg.On: perilaku lama (maks. 200). total = jumlah semua hasil (untuk halaman).
func (s *Service) SellerList(ctx context.Context, storeID, filter string, pg paging.Params) ([]Return, int, error) {
	where := "r.store_id = $1"
	args := []any{storeID}
	if filter != "all" {
		where += " AND r.status IN ('requested','disputed','approved')"
	}
	if pg.Q != "" {
		args = append(args, paging.Like(pg.Q))
		where += " AND (r.return_no ILIKE $2 OR o.order_no ILIKE $2 OR u.name ILIKE $2)"
	}
	limit, offset := 200, 0
	if pg.On {
		limit, offset = pg.Size, pg.Offset()
	}
	l, err := s.query(ctx, s.Pool, where+" ORDER BY r.created_at DESC LIMIT "+fmt.Sprint(limit)+" OFFSET "+fmt.Sprint(offset), 0, args...)
	if err != nil {
		return nil, 0, err
	}
	for i := range l {
		forSeller(&l[i])
	}
	total := len(l)
	if pg.On {
		if err := s.Pool.QueryRow(ctx, `SELECT count(*) FROM returns r JOIN orders o ON o.id = r.order_id LEFT JOIN users u ON u.id = r.buyer_id
			WHERE `+where, args...).Scan(&total); err != nil {
			return nil, 0, err
		}
	}
	return l, total, nil
}

func (s *Service) SellerGet(ctx context.Context, storeID, id string) (*Return, error) {
	r, err := s.one(ctx, "r.id = $1 AND r.store_id = $2", id, storeID)
	if err != nil {
		return nil, err
	}
	return forSeller(r), nil
}

type SellerAnswer struct {
	Accept   bool   `json:"accept"`
	Note     string `json:"note"`
	KeepItem bool   `json:"keep_item"` // seller setuju refund tanpa barang dikembalikan
}

func (s *Service) SellerRespond(ctx context.Context, storeID, userID, id string, in SellerAnswer) (*Return, error) {
	in.Note = strings.TrimSpace(in.Note)
	if len([]rune(in.Note)) > 1000 {
		return nil, FieldError{"note", "too long"}
	}
	if !in.Accept && len([]rune(in.Note)) < 5 {
		return nil, FieldError{"note", "tell the buyer and Phsar why you reject"}
	}
	err := db.WithActor(ctx, s.Pool, userID, "seller", "", func(tx pgx.Tx) error {
		var orderID, orderNo, buyerID, lang, no, amount string
		var dl time.Time
		err := tx.QueryRow(ctx, `SELECT r.order_id, o.order_no, r.buyer_id, coalesce(u.language::text, 'km'), r.return_no, r.requested_amount::text, r.seller_deadline
			FROM returns r JOIN orders o ON o.id = r.order_id JOIN users u ON u.id = r.buyer_id
			WHERE r.id = $1 AND r.store_id = $2 AND r.status = 'requested' FOR UPDATE OF r`, id, storeID).
			Scan(&orderID, &orderNo, &buyerID, &lang, &no, &amount, &dl)
		if errors.Is(err, pgx.ErrNoRows) {
			return s.stateOrMissing(ctx, tx, "id = $1 AND store_id = $2", id, storeID)
		}
		if err != nil {
			return err
		}
		if time.Now().After(dl) {
			return ErrDeadline
		}
		link := "/" + lang + "/orders/" + orderNo
		if in.Accept {
			if _, err := tx.Exec(ctx, `UPDATE returns SET status = 'approved', seller_note = nullif($2, ''), seller_responded_at = now(),
				refund_amount = requested_amount, needs_pickup = $3 WHERE id = $1`, id, in.Note, !in.KeepItem); err != nil {
				return err
			}
			if err := event(ctx, tx, orderID, "return_approved", "seller", userID, no); err != nil {
				return err
			}
			return notify.EnqueueUser(ctx, tx, buyerID, notify.BuyerReturnApproved, notify.Params{Order: orderNo, Product: "$" + amount, Note: in.Note}, link)
		}
		if _, err := tx.Exec(ctx, `UPDATE returns SET status = 'disputed', escalated_reason = 'seller_rejected', seller_note = $2,
			seller_responded_at = now() WHERE id = $1`, id, in.Note); err != nil {
			return err
		}
		if err := event(ctx, tx, orderID, "return_disputed", "seller", userID, no); err != nil {
			return err
		}
		return notify.EnqueueUser(ctx, tx, buyerID, notify.BuyerReturnReview, notify.Params{Order: orderNo}, link)
	})
	if err != nil {
		return nil, err
	}
	return s.SellerGet(ctx, storeID, id)
}

// EscalateOverdue: seller tidak menjawab dalam 48 jam → Phsar yang memutuskan.
func (s *Service) EscalateOverdue(ctx context.Context) (int, error) {
	n := 0
	err := db.WithActor(ctx, s.Pool, "", "system", "", func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, `UPDATE returns r SET status = 'disputed', escalated_reason = 'seller_timeout'
			FROM orders o, users u WHERE o.id = r.order_id AND u.id = r.buyer_id AND r.status = 'requested' AND r.seller_deadline < now()
			RETURNING r.order_id, o.order_no, r.store_id, r.buyer_id, coalesce(u.language::text, 'km'), r.return_no`)
		if err != nil {
			return err
		}
		type row struct{ order, no, store, buyer, lang, rno string }
		var list []row
		for rows.Next() {
			var r row
			if err := rows.Scan(&r.order, &r.no, &r.store, &r.buyer, &r.lang, &r.rno); err != nil {
				rows.Close()
				return err
			}
			list = append(list, r)
		}
		rows.Close()
		for _, r := range list {
			if err := event(ctx, tx, r.order, "return_disputed", "system", "", r.rno+" seller_timeout"); err != nil {
				return err
			}
			if err := notify.Enqueue(ctx, tx, r.store, notify.ReturnEscalated, notify.Params{Order: r.no}, "#/store/"+r.store+"/returns"); err != nil {
				return err
			}
			if err := notify.EnqueueUser(ctx, tx, r.buyer, notify.BuyerReturnReview, notify.Params{Order: r.no}, "/"+r.lang+"/orders/"+r.no); err != nil {
				return err
			}
			n++
		}
		return nil
	})
	return n, err
}

func (s *Service) Worker(ctx context.Context) {
	t := time.NewTicker(5 * time.Minute)
	defer t.Stop()
	for {
		if n, err := jobs.Track(ctx, s.Pool, "returns.escalate", s.EscalateOverdue); err != nil && ctx.Err() == nil { // jobs:
			s.Log.Error("retur: eskalasi gagal", "err", err)
		} else if n > 0 {
			s.Log.Info("retur diteruskan ke admin (seller tidak menjawab)", "n", n)
		}
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}
	}
}

// ---------------------------------------------------------------------------
// Admin
// ---------------------------------------------------------------------------

type Counts struct {
	Requested int `json:"requested"`
	Disputed  int `json:"disputed"`
	Approved  int `json:"approved"`
}

func (s *Service) AdminList(ctx context.Context, status string) ([]Return, Counts, error) {
	var c Counts
	if err := s.Pool.QueryRow(ctx, `SELECT count(*) FILTER (WHERE status = 'requested'), count(*) FILTER (WHERE status = 'disputed'),
		count(*) FILTER (WHERE status = 'approved') FROM returns`).Scan(&c.Requested, &c.Disputed, &c.Approved); err != nil {
		return nil, c, err
	}
	where := "r.status IN ('requested','disputed','approved')"
	switch status {
	case "requested", "disputed", "approved":
		where = "r.status = '" + status + "'"
	case "closed":
		where = "r.status IN ('rejected','refunded','cancelled')"
	}
	l, err := s.query(ctx, s.Pool, where, 300)
	return l, c, err
}

func (s *Service) AdminGet(ctx context.Context, id string) (*Return, error) {
	return s.one(ctx, "r.id = $1", id)
}

type Decision struct {
	Decision    string `json:"decision"` // refund · reject
	Amount      string `json:"amount"`   // refund: boleh sebagian
	Note        string `json:"note"`
	NeedsPickup bool   `json:"needs_pickup"`
}

// Decide: keputusan Phsar untuk retur yang belum selesai (biasanya disputed; boleh juga sebelum seller menjawab).
func (s *Service) Decide(ctx context.Context, adminID, id string, in Decision) (*Return, error) {
	in.Note = strings.TrimSpace(in.Note)
	if in.Decision != "refund" && in.Decision != "reject" {
		return nil, FieldError{"decision", "refund or reject"}
	}
	if n := len([]rune(in.Note)); n < 5 || n > 1000 {
		return nil, FieldError{"note", "explain the decision to the buyer and seller (5–1000 characters)"}
	}
	err := db.WithActor(ctx, s.Pool, adminID, "admin", "", func(tx pgx.Tx) error {
		var orderID, orderNo, storeID, buyerID, lang, no, total string
		err := tx.QueryRow(ctx, `SELECT r.order_id, o.order_no, r.store_id, r.buyer_id, coalesce(u.language::text, 'km'), r.return_no, o.total::text
			FROM returns r JOIN orders o ON o.id = r.order_id JOIN users u ON u.id = r.buyer_id
			WHERE r.id = $1 AND r.status IN ('requested','disputed') FOR UPDATE OF r`, id).Scan(&orderID, &orderNo, &storeID, &buyerID, &lang, &no, &total)
		if errors.Is(err, pgx.ErrNoRows) {
			return s.stateOrMissing(ctx, tx, "id = $1", id)
		}
		if err != nil {
			return err
		}
		link := "/" + lang + "/orders/" + orderNo
		sellerLink := "#/store/" + storeID + "/returns"
		if in.Decision == "reject" {
			if _, err := tx.Exec(ctx, `UPDATE returns SET status = 'rejected', decision_note = $2, decided_by = $3, decided_at = now() WHERE id = $1`,
				id, in.Note, adminID); err != nil {
				return err
			}
			if err := event(ctx, tx, orderID, "return_rejected", "admin", adminID, no); err != nil {
				return err
			}
			if err := notify.EnqueueUser(ctx, tx, buyerID, notify.BuyerReturnRejected, notify.Params{Order: orderNo, Note: in.Note}, link); err != nil {
				return err
			}
			return notify.Enqueue(ctx, tx, storeID, notify.ReturnDecided, notify.Params{Order: orderNo, Note: in.Note}, sellerLink)
		}
		amt, err := decimal.NewFromString(strings.TrimSpace(in.Amount))
		max, _ := decimal.NewFromString(total)
		if err != nil || !amt.IsPositive() || amt.GreaterThan(max) || amt.Exponent() < -2 {
			return FieldError{"amount", "between 0.01 and the order total ($" + total + ")"}
		}
		if _, err := tx.Exec(ctx, `UPDATE returns SET status = 'approved', refund_amount = $2, needs_pickup = $3, decision_note = $4,
			decided_by = $5, decided_at = now() WHERE id = $1`, id, amt, in.NeedsPickup, in.Note, adminID); err != nil {
			return err
		}
		if err := event(ctx, tx, orderID, "return_approved", "admin", adminID, no); err != nil {
			return err
		}
		if err := notify.EnqueueUser(ctx, tx, buyerID, notify.BuyerReturnApproved, notify.Params{Order: orderNo, Product: "$" + amt.StringFixed(2), Note: in.Note}, link); err != nil {
			return err
		}
		return notify.Enqueue(ctx, tx, storeID, notify.ReturnDecided, notify.Params{Order: orderNo, Note: in.Note}, sellerLink)
	})
	if err != nil {
		return nil, err
	}
	return s.AdminGet(ctx, id)
}

// MarkRefunded: admin sudah mengembalikan uang lewat PayWay. Saldo seller dipotong bagiannya
// (refund × seller_amount / total). Kalau seluruh pesanan dikembalikan sebelum masuk saldo, pesanan menjadi "refunded".
func (s *Service) MarkRefunded(ctx context.Context, adminID, id, ref, note string) (*Return, error) {
	ref, note = strings.TrimSpace(ref), strings.TrimSpace(note)
	if len(ref) < 3 || len(ref) > 80 {
		return nil, FieldError{"reference", "enter the PayWay refund transaction number"}
	}
	err := db.WithActor(ctx, s.Pool, adminID, "admin", "", func(tx pgx.Tx) error {
		var orderID, orderNo, storeID, buyerID, lang, no, status string
		var refund, total, sellerAmt decimal.Decimal
		err := tx.QueryRow(ctx, `SELECT r.order_id, o.order_no, r.store_id, r.buyer_id, coalesce(u.language::text, 'km'), r.return_no, r.refund_amount,
			o.total, o.seller_amount, o.status::text FROM returns r JOIN orders o ON o.id = r.order_id JOIN users u ON u.id = r.buyer_id
			WHERE r.id = $1 AND r.status = 'approved' FOR UPDATE OF r, o`, id).
			Scan(&orderID, &orderNo, &storeID, &buyerID, &lang, &no, &refund, &total, &sellerAmt, &status)
		if errors.Is(err, pgx.ErrNoRows) {
			return s.stateOrMissing(ctx, tx, "id = $1", id)
		}
		if err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `UPDATE returns SET status = 'refunded', refund_ref = $2, refunded_by = $3, refunded_at = now(),
			decision_note = coalesce(decision_note, nullif($4, '')) WHERE id = $1`, id, ref, adminID, note); err != nil {
			return err
		}
		var earned bool
		if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM seller_ledger WHERE order_id = $1 AND kind = 'order_earning')`, orderID).Scan(&earned); err != nil {
			return err
		}
		share := refund.Mul(sellerAmt).Div(total).Round(2)
		if share.GreaterThan(sellerAmt) {
			share = sellerAmt
		}
		switch {
		case !earned && refund.Equal(total):
			// belum masuk saldo & dikembalikan penuh: seller tidak menerima apa pun.
			if _, err := tx.Exec(ctx, `UPDATE orders SET status = 'refunded', refund_status = 'refunded', refund_amount = $2, refunded_at = now(),
				refund_ref = $3, refunded_by = $4 WHERE id = $1`, orderID, refund, ref, adminID); err != nil {
				return err
			}
			share = decimal.Zero
		case !earned:
			// sebagian: pesanan diselesaikan sekarang (saldo masuk) lalu dipotong bagian retur.
			if _, err := tx.Exec(ctx, `UPDATE orders SET status = 'completed', completed_at = now() WHERE id = $1`, orderID); err != nil {
				return err
			}
			if err := event(ctx, tx, orderID, "completed", "system", "", ""); err != nil {
				return err
			}
			if sellerAmt.IsPositive() {
				if _, err := tx.Exec(ctx, `INSERT INTO seller_ledger (store_id, kind, amount, order_id, note) VALUES ($1, 'order_earning', $2, $3, $4)`,
					storeID, sellerAmt, orderID, orderNo); err != nil {
					return err
				}
			}
		}
		if share.IsPositive() {
			if _, err := tx.Exec(ctx, `INSERT INTO seller_ledger (store_id, kind, amount, order_id, return_id, note, created_by)
				VALUES ($1, 'return_refund', $2, $3, $4, $5, $6)`, storeID, share.Neg(), orderID, id, no+" · "+orderNo, adminID); err != nil {
				return err
			}
		}
		if err := event(ctx, tx, orderID, "return_refunded", "admin", adminID, no); err != nil {
			return err
		}
		if err := notify.EnqueueUser(ctx, tx, buyerID, notify.BuyerOrderRefunded, notify.Params{Order: orderNo, Product: "$" + refund.StringFixed(2)},
			"/"+lang+"/orders/"+orderNo); err != nil {
			return err
		}
		return notify.Enqueue(ctx, tx, storeID, notify.ReturnRefunded, notify.Params{Order: orderNo, Product: "$" + share.StringFixed(2)}, "#/store/"+storeID+"/returns")
	})
	if err != nil {
		return nil, err
	}
	return s.AdminGet(ctx, id)
}

// ---------------------------------------------------------------------------

func event(ctx context.Context, tx pgx.Tx, orderID, status, actorType, actorID, note string) error {
	_, err := tx.Exec(ctx, `INSERT INTO order_events (order_id, status, actor_type, actor_id, note) VALUES ($1, $2, $3, nullif($4, '')::uuid, nullif($5, ''))`,
		orderID, status, actorType, actorID, note)
	return err
}

const codeAlphabet = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"

func randCode(n int) string {
	b := make([]byte, n)
	for i := range b {
		k, _ := rand.Int(rand.Reader, big.NewInt(int64(len(codeAlphabet))))
		b[i] = codeAlphabet[k.Int64()]
	}
	return string(b)
}

var phnomPenh = func() *time.Location {
	if l, err := time.LoadLocation("Asia/Phnom_Penh"); err == nil {
		return l
	}
	return time.FixedZone("ICT", 7*3600)
}()
