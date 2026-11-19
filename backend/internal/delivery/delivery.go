// Package delivery: kurir Phsar (staf) menjemput paket di toko dan mengantar ke pembeli.
//
// Tugas dibuat saat seller menandai pesanan "siap dijemput". Kurir mengambil tugas sendiri dari
// daftar "Tersedia" (admin juga bisa menugaskan). Di toko kurir memindai barcode label paket —
// harus sama dengan nomor pesanan — supaya tidak salah bawa paket. Di tujuan kurir menandai
// "Terantar" (boleh dengan foto bukti) atau "Gagal" dengan alasan.
package delivery

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"phsar/internal/auth"
	"phsar/internal/notify"
	"phsar/internal/platform/db"
	"phsar/internal/storage"
)

var (
	ErrNotCourier   = errors.New("not a courier")
	ErrNotFound     = errors.New("not found")
	ErrTaken        = errors.New("task already taken")
	ErrBadStatus    = errors.New("task status does not allow this")
	ErrWrongPackage = errors.New("scanned code does not match this order")
	ErrValidation   = errors.New("validation")
	ErrTooMany      = errors.New("too many active tasks")
	ErrMaxAttempts  = errors.New("no delivery attempts left") // courierreturn: harus dikembalikan ke toko
)

type fieldError struct{ field, msg string }

func (e fieldError) Error() string { return e.field + ": " + e.msg }
func (e fieldError) Unwrap() error { return ErrValidation }

// MaxActive: kurir paling banyak memegang sekian tugas yang belum terantar.
const MaxActive = 8

// courierreturn: MaxAttempts — batas percobaan antar (jemput = percobaan 1, tiap "coba lagi" +1).
// Setelah itu paket harus dikembalikan ke toko.
const MaxAttempts = 3

type Service struct {
	Pool     *pgxpool.Pool
	Files    storage.Storage
	Uploader storage.Uploader
	Log      *slog.Logger
}

// ---------------------------------------------------------------------------
// Kurir
// ---------------------------------------------------------------------------

type Courier struct {
	UserID      string     `json:"user_id"`
	Name        string     `json:"name"`
	Phone       string     `json:"phone"`
	Vehicle     string     `json:"vehicle"`
	PlateNo     string     `json:"plate_no,omitempty"`
	Active      bool       `json:"active"`
	Online      bool       `json:"online"`
	LastSeenAt  *time.Time `json:"last_seen_at,omitempty"`
	ActiveTasks int        `json:"active_tasks"`
	DoneToday   int        `json:"done_today"`
}

const courierCols = `c.user_id, c.display_name, coalesce(u.phone, ''), c.vehicle, coalesce(c.plate_no, ''), c.active, c.is_online, c.last_seen_at,
	(SELECT count(*) FROM deliveries d WHERE d.courier_id = c.user_id AND d.status IN ('assigned','picked_up','returning')),
	(SELECT count(*) FROM deliveries d WHERE d.courier_id = c.user_id AND d.status = 'delivered' AND d.delivered_at >= date_trunc('day', now() AT TIME ZONE 'Asia/Phnom_Penh') AT TIME ZONE 'Asia/Phnom_Penh')`

func scanCourier(row pgx.Row) (*Courier, error) {
	var c Courier
	err := row.Scan(&c.UserID, &c.Name, &c.Phone, &c.Vehicle, &c.PlateNo, &c.Active, &c.Online, &c.LastSeenAt, &c.ActiveTasks, &c.DoneToday)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotCourier
	}
	return &c, err
}

// Me: profil kurir yang login (ErrNotCourier kalau bukan kurir aktif).
func (s *Service) Me(ctx context.Context, userID string) (*Courier, error) {
	c, err := scanCourier(s.Pool.QueryRow(ctx, `SELECT `+courierCols+` FROM couriers c JOIN users u ON u.id = c.user_id WHERE c.user_id = $1 AND c.active`, userID))
	if err == nil {
		_, _ = s.Pool.Exec(ctx, `UPDATE couriers SET last_seen_at = now() WHERE user_id = $1`, userID)
	}
	return c, err
}

func (s *Service) SetOnline(ctx context.Context, userID string, online bool) (*Courier, error) {
	if _, err := s.Me(ctx, userID); err != nil {
		return nil, err
	}
	if _, err := s.Pool.Exec(ctx, `UPDATE couriers SET is_online = $2, last_seen_at = now() WHERE user_id = $1`, userID, online); err != nil {
		return nil, err
	}
	return s.Me(ctx, userID)
}

// ---------------------------------------------------------------------------
// Tugas
// ---------------------------------------------------------------------------

type Task struct {
	ID          string  `json:"id"`
	Status      string  `json:"status"`
	OrderNo     string  `json:"order_no"`
	OrderStatus string  `json:"order_status"`
	StoreName   string  `json:"store_name"`
	ItemCount   int     `json:"item_count"`
	PickupLabel string  `json:"pickup_label"`
	PickupAddr  string  `json:"pickup_address"`
	PickupLat   float64 `json:"pickup_lat"`
	PickupLng   float64 `json:"pickup_lng"`
	PickupName  string  `json:"pickup_contact"`
	PickupPhone string  `json:"pickup_phone"`
	DropoffLat  float64 `json:"dropoff_lat"`
	DropoffLng  float64 `json:"dropoff_lng"`
	DropoffArea string  `json:"dropoff_area"`
	// Data pembeli hanya untuk kurir yang memegang tugas (dan admin).
	Recipient   string     `json:"recipient,omitempty"`
	Phone       string     `json:"phone,omitempty"`
	AddressLine string     `json:"address_line,omitempty"`
	AddressNote string     `json:"address_note,omitempty"`
	CourierID   string     `json:"courier_id,omitempty"`
	CourierName string     `json:"courier_name,omitempty"`
	Attempts    int        `json:"attempts"`
	FailReason  string     `json:"fail_reason,omitempty"`
	FailNote    string     `json:"fail_note,omitempty"`
	HasProof    bool       `json:"has_proof"`
	ReadyAt     time.Time  `json:"ready_at"`
	AssignedAt  *time.Time `json:"assigned_at,omitempty"`
	PickedUpAt  *time.Time `json:"picked_up_at,omitempty"`
	DeliveredAt *time.Time `json:"delivered_at,omitempty"`
	FailedAt    *time.Time `json:"failed_at,omitempty"`
	Items       []TaskItem `json:"items,omitempty"`
}

type TaskItem struct {
	Name    string `json:"name"`
	Options string `json:"options"`
	Qty     int    `json:"qty"`
}

const taskCols = `d.id, d.status, o.order_no, o.status::text, st.name, o.item_count, d.pickup_label, d.pickup_address, d.pickup_lat::float8, d.pickup_lng::float8,
	d.pickup_contact, d.pickup_phone, d.dropoff_lat::float8, d.dropoff_lng::float8, coalesce(sa.code, ''),
	c.recipient, c.phone, c.address_line, coalesce(c.address_note, ''), coalesce(d.courier_id::text, ''), coalesce(cr.display_name, ''),
	d.attempts, coalesce(d.fail_reason, ''), coalesce(d.fail_note, ''), d.proof_key IS NOT NULL, d.created_at, d.assigned_at, d.picked_up_at, d.delivered_at, d.failed_at`

const taskJoins = ` FROM deliveries d JOIN orders o ON o.id = d.order_id JOIN stores st ON st.id = o.store_id JOIN checkouts c ON c.id = o.checkout_id
	LEFT JOIN service_areas sa ON sa.code = c.service_area_code LEFT JOIN couriers cr ON cr.user_id = d.courier_id `

func (s *Service) tasks(ctx context.Context, where, order string, limit int, args ...any) ([]Task, error) {
	rows, err := s.Pool.Query(ctx, `SELECT `+taskCols+taskJoins+` WHERE `+where+` ORDER BY `+order+` LIMIT `+fmt.Sprint(limit), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Task{}
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Status, &t.OrderNo, &t.OrderStatus, &t.StoreName, &t.ItemCount, &t.PickupLabel, &t.PickupAddr, &t.PickupLat, &t.PickupLng,
			&t.PickupName, &t.PickupPhone, &t.DropoffLat, &t.DropoffLng, &t.DropoffArea, &t.Recipient, &t.Phone, &t.AddressLine, &t.AddressNote,
			&t.CourierID, &t.CourierName, &t.Attempts, &t.FailReason, &t.FailNote, &t.HasProof, &t.ReadyAt, &t.AssignedAt, &t.PickedUpAt, &t.DeliveredAt, &t.FailedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func hideBuyer(t *Task) {
	t.Recipient, t.Phone, t.AddressLine, t.AddressNote = "", "", "", ""
}

func (s *Service) items(ctx context.Context, taskID string) ([]TaskItem, error) {
	rows, err := s.Pool.Query(ctx, `SELECT oi.name, oi.options, oi.qty FROM order_items oi JOIN deliveries d ON d.order_id = oi.order_id WHERE d.id = $1 ORDER BY oi.name`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []TaskItem{}
	for rows.Next() {
		var it TaskItem
		if err := rows.Scan(&it.Name, &it.Options, &it.Qty); err != nil {
			return nil, err
		}
		out = append(out, it)
	}
	return out, rows.Err()
}

// CourierTasks: view = available (belum diambil) · mine (sedang dipegang) · done (hari ini).
func (s *Service) CourierTasks(ctx context.Context, courierID, view string) ([]Task, error) {
	if _, err := s.Me(ctx, courierID); err != nil {
		return nil, err
	}
	var list []Task
	var err error
	switch view {
	case "mine":
		// courierreturn: + 'returning' (paket dibawa kembali ke toko)
		list, err = s.tasks(ctx, `d.courier_id = $1 AND d.status IN ('assigned','picked_up','failed','returning')`, `d.status = 'picked_up' DESC, d.assigned_at`, 50, courierID)
	case "done":
		list, err = s.tasks(ctx, `d.courier_id = $1 AND d.status = 'delivered' AND d.delivered_at > now() - interval '24 hours'`, `d.delivered_at DESC`, 100, courierID)
	default:
		list, err = s.tasks(ctx, `d.status = 'unassigned' AND o.status = 'ready_for_pickup'`, `d.created_at`, 100)
		for i := range list {
			hideBuyer(&list[i]) // belum diambil: alamat & HP pembeli belum perlu dilihat
		}
	}
	return list, err
}

func (s *Service) CourierTask(ctx context.Context, courierID, taskID string) (*Task, error) {
	if _, err := s.Me(ctx, courierID); err != nil {
		return nil, err
	}
	if len(taskID) != 36 {
		return nil, ErrNotFound
	}
	list, err := s.tasks(ctx, `d.id = $1 AND (d.courier_id = $2 OR d.status = 'unassigned')`, `d.id`, 1, taskID, courierID)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, ErrNotFound
	}
	t := &list[0]
	if t.CourierID != courierID {
		hideBuyer(t)
	}
	t.Items, err = s.items(ctx, t.ID)
	return t, err
}

type locked struct {
	id, orderID, orderNo, status, courierID, storeID, storeName, buyerID, lang string
}

func lockTask(ctx context.Context, tx pgx.Tx, taskID string) (*locked, error) {
	if len(taskID) != 36 {
		return nil, ErrNotFound
	}
	var l locked
	err := tx.QueryRow(ctx, `SELECT d.id, o.id, o.order_no, d.status, coalesce(d.courier_id::text, ''), o.store_id, st.name, o.buyer_id, c.lang::text
		FROM deliveries d JOIN orders o ON o.id = d.order_id JOIN stores st ON st.id = o.store_id JOIN checkouts c ON c.id = o.checkout_id
		WHERE d.id = $1 FOR UPDATE OF d, o`, taskID).Scan(&l.id, &l.orderID, &l.orderNo, &l.status, &l.courierID, &l.storeID, &l.storeName, &l.buyerID, &l.lang)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &l, err
}

func event(ctx context.Context, tx pgx.Tx, orderID, status, actorType, actorID, note string) error {
	_, err := tx.Exec(ctx, `INSERT INTO order_events (order_id, status, actor_type, actor_id, note) VALUES ($1, $2, $3::actor_type, nullif($4, '')::uuid, nullif($5, ''))`,
		orderID, status, actorType, actorID, note)
	return err
}

// Claim: kurir mengambil tugas dari daftar tersedia (yang pertama menang).
func (s *Service) Claim(ctx context.Context, courierID, taskID string) (*Task, error) {
	me, err := s.Me(ctx, courierID)
	if err != nil {
		return nil, err
	}
	if me.ActiveTasks >= MaxActive {
		return nil, ErrTooMany
	}
	err = db.WithActor(ctx, s.Pool, courierID, "courier", "", func(tx pgx.Tx) error {
		l, err := lockTask(ctx, tx, taskID)
		if err != nil {
			return err
		}
		if l.status == "assigned" && l.courierID == courierID {
			return nil
		}
		if l.status != "unassigned" {
			return ErrTaken
		}
		if _, err := tx.Exec(ctx, `UPDATE deliveries SET status = 'assigned', courier_id = $2, assigned_at = now() WHERE id = $1`, l.id, courierID); err != nil {
			return err
		}
		return event(ctx, tx, l.orderID, "courier_assigned", "courier", courierID, me.Name)
	})
	if err != nil {
		return nil, err
	}
	return s.CourierTask(ctx, courierID, taskID)
}

// Release: kurir melepas tugas yang belum dijemput (kembali ke daftar tersedia).
func (s *Service) Release(ctx context.Context, courierID, taskID string) error {
	if _, err := s.Me(ctx, courierID); err != nil {
		return err
	}
	return db.WithActor(ctx, s.Pool, courierID, "courier", "", func(tx pgx.Tx) error {
		l, err := lockTask(ctx, tx, taskID)
		if err != nil {
			return err
		}
		if l.courierID != courierID {
			return ErrNotFound
		}
		if l.status != "assigned" {
			return ErrBadStatus
		}
		if _, err := tx.Exec(ctx, `UPDATE deliveries SET status = 'unassigned', courier_id = NULL, assigned_at = NULL WHERE id = $1`, l.id); err != nil {
			return err
		}
		return event(ctx, tx, l.orderID, "courier_released", "courier", courierID, "")
	})
}

// normalizeCode: hasil scan / ketikan → bentuk nomor pesanan (huruf besar, tanpa spasi).
func normalizeCode(s string) string {
	return strings.ToUpper(strings.Join(strings.Fields(s), ""))
}

// PickUp: kurir di toko memindai label — harus sama dengan nomor pesanan tugas ini.
func (s *Service) PickUp(ctx context.Context, courierID, taskID, code string) (*Task, error) {
	me, err := s.Me(ctx, courierID)
	if err != nil {
		return nil, err
	}
	err = db.WithActor(ctx, s.Pool, courierID, "courier", "", func(tx pgx.Tx) error {
		l, err := lockTask(ctx, tx, taskID)
		if err != nil {
			return err
		}
		if l.courierID != courierID {
			return ErrNotFound
		}
		if l.status == "picked_up" {
			return nil
		}
		if l.status != "assigned" {
			return ErrBadStatus
		}
		if normalizeCode(code) != l.orderNo {
			return ErrWrongPackage
		}
		if _, err := tx.Exec(ctx, `UPDATE deliveries SET status = 'picked_up', picked_up_at = now(), attempts = attempts + 1 WHERE id = $1`, l.id); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `UPDATE orders SET status = 'picked_up' WHERE id = $1 AND status = 'ready_for_pickup'`, l.orderID); err != nil {
			return err
		}
		if err := event(ctx, tx, l.orderID, "picked_up", "courier", courierID, me.Name); err != nil {
			return err
		}
		if err := notify.Enqueue(ctx, tx, l.storeID, notify.OrderPickedUp, notify.Params{Order: l.orderNo, Person: me.Name}, "#/store/"+l.storeID+"/orders"); err != nil {
			return err
		}
		return notify.EnqueueUser(ctx, tx, l.buyerID, notify.BuyerOrderOnTheWay, notify.Params{Order: l.orderNo, Store: l.storeName, Person: me.Name}, "/"+l.lang+"/orders/"+l.orderNo)
	})
	if err != nil {
		return nil, err
	}
	return s.CourierTask(ctx, courierID, taskID)
}

// Deliver: paket sudah diterima pembeli. proofKey = upload "delivery_proof" milik kurir ini (boleh kosong).
func (s *Service) Deliver(ctx context.Context, courierID, taskID, proofKey string) (*Task, error) {
	if _, err := s.Me(ctx, courierID); err != nil {
		return nil, err
	}
	if proofKey != "" {
		if err := storage.CheckOwner(ctx, s.Pool, proofKey, courierID, "delivery_proof"); err != nil {
			return nil, fieldError{"proof_key", "upload the photo again"}
		}
	}
	err := db.WithActor(ctx, s.Pool, courierID, "courier", "", func(tx pgx.Tx) error {
		l, err := lockTask(ctx, tx, taskID)
		if err != nil {
			return err
		}
		if l.courierID != courierID {
			return ErrNotFound
		}
		if l.status == "delivered" {
			return nil
		}
		if l.status != "picked_up" && l.status != "failed" {
			return ErrBadStatus
		}
		if _, err := tx.Exec(ctx, `UPDATE deliveries SET status = 'delivered', delivered_at = now(), proof_key = coalesce(nullif($2, ''), proof_key) WHERE id = $1`, l.id, proofKey); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `UPDATE orders SET status = 'delivered' WHERE id = $1 AND status = 'picked_up'`, l.orderID); err != nil {
			return err
		}
		if err := event(ctx, tx, l.orderID, "delivered", "courier", courierID, ""); err != nil {
			return err
		}
		if err := notify.Enqueue(ctx, tx, l.storeID, notify.OrderDelivered, notify.Params{Order: l.orderNo}, "#/store/"+l.storeID+"/orders"); err != nil {
			return err
		}
		return notify.EnqueueUser(ctx, tx, l.buyerID, notify.BuyerOrderDelivered, notify.Params{Order: l.orderNo, Store: l.storeName}, "/"+l.lang+"/orders/"+l.orderNo)
	})
	if err != nil {
		return nil, err
	}
	return s.CourierTask(ctx, courierID, taskID)
}

var failReasons = map[string]bool{"buyer_unreachable": true, "buyer_refused": true, "wrong_address": true, "damaged": true, "other": true}

// Fail: pengantaran gagal. Paket tetap dibawa kurir; admin memutuskan coba lagi atau kembalikan.
func (s *Service) Fail(ctx context.Context, courierID, taskID, reason, note string) (*Task, error) {
	return s.FailWithProof(ctx, courierID, taskID, reason, note, "")
}

// FailWithProof (courierreturn): seperti Fail, dengan foto bukti opsional (upload "delivery_proof" milik kurir ini).
func (s *Service) FailWithProof(ctx context.Context, courierID, taskID, reason, note, proofKey string) (*Task, error) {
	if proofKey != "" {
		if err := storage.CheckOwner(ctx, s.Pool, proofKey, courierID, "delivery_proof"); err != nil {
			return nil, fieldError{"proof_key", "upload the photo again"}
		}
	}
	if !failReasons[reason] {
		return nil, fieldError{"reason", "choose a reason"}
	}
	note = strings.TrimSpace(note)
	if reason == "other" && len([]rune(note)) < 3 {
		return nil, fieldError{"note", "explain what happened"}
	}
	if len([]rune(note)) > 300 {
		note = string([]rune(note)[:300])
	}
	if _, err := s.Me(ctx, courierID); err != nil {
		return nil, err
	}
	err := db.WithActor(ctx, s.Pool, courierID, "courier", "", func(tx pgx.Tx) error {
		l, err := lockTask(ctx, tx, taskID)
		if err != nil {
			return err
		}
		if l.courierID != courierID {
			return ErrNotFound
		}
		if l.status != "picked_up" {
			return ErrBadStatus
		}
		if _, err := tx.Exec(ctx, `UPDATE deliveries SET status = 'failed', failed_at = now(), fail_reason = $2, fail_note = nullif($3, ''),
			fail_proof_key = coalesce(nullif($4, ''), fail_proof_key) WHERE id = $1`, l.id, reason, note, proofKey); err != nil {
			return err
		}
		return event(ctx, tx, l.orderID, "delivery_failed", "courier", courierID, reason)
	})
	if err != nil {
		return nil, err
	}
	return s.CourierTask(ctx, courierID, taskID)
}

// Retry: kurir mencoba antar lagi tugas yang gagal (mis. pembeli sudah bisa dihubungi).
func (s *Service) Retry(ctx context.Context, courierID, taskID string) (*Task, error) {
	if _, err := s.Me(ctx, courierID); err != nil {
		return nil, err
	}
	err := db.WithActor(ctx, s.Pool, courierID, "courier", "", func(tx pgx.Tx) error {
		l, err := lockTask(ctx, tx, taskID)
		if err != nil {
			return err
		}
		if l.courierID != courierID {
			return ErrNotFound
		}
		if l.status != "failed" {
			return ErrBadStatus
		}
		// courierreturn: batas percobaan → kembalikan ke toko
		var attempts int
		if err := tx.QueryRow(ctx, `SELECT attempts FROM deliveries WHERE id = $1`, l.id).Scan(&attempts); err != nil {
			return err
		}
		if attempts >= MaxAttempts {
			return ErrMaxAttempts
		}
		if _, err := tx.Exec(ctx, `UPDATE deliveries SET status = 'picked_up', attempts = attempts + 1 WHERE id = $1`, l.id); err != nil {
			return err
		}
		return event(ctx, tx, l.orderID, "delivery_retry", "courier", courierID, "")
	})
	if err != nil {
		return nil, err
	}
	return s.CourierTask(ctx, courierID, taskID)
}

// ---------------------------------------------------------------------------
// Admin: kelola kurir & pantau / tugaskan pengantaran
// ---------------------------------------------------------------------------

func (s *Service) ListCouriers(ctx context.Context) ([]Courier, error) {
	rows, err := s.Pool.Query(ctx, `SELECT `+courierCols+` FROM couriers c JOIN users u ON u.id = c.user_id ORDER BY c.active DESC, c.is_online DESC, c.display_name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Courier{}
	for rows.Next() {
		c, err := scanCourier(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *c)
	}
	return out, rows.Err()
}

type CourierInput struct {
	Phone   string `json:"phone"`
	Name    string `json:"name"`
	Vehicle string `json:"vehicle"`
	PlateNo string `json:"plate_no"`
	Active  *bool  `json:"active"`
}

// UpsertCourier: tambah/ubah kurir berdasarkan nomor HP. Kalau nomor belum pernah login, akunnya dibuat
// (kurir tinggal login dengan nomor itu).
func (s *Service) UpsertCourier(ctx context.Context, adminID string, in CourierInput) (*Courier, error) {
	phone, err := auth.NormalizePhone(in.Phone)
	if err != nil {
		return nil, fieldError{"phone", "enter a valid phone number"}
	}
	in.Name = strings.TrimSpace(in.Name)
	if len([]rune(in.Name)) < 2 || len([]rune(in.Name)) > 80 {
		return nil, fieldError{"name", "enter the courier's name"}
	}
	if in.Vehicle == "" {
		in.Vehicle = "motorbike"
	}
	if !map[string]bool{"motorbike": true, "tuktuk": true, "car": true, "bicycle": true}[in.Vehicle] {
		return nil, fieldError{"vehicle", "choose a vehicle"}
	}
	in.PlateNo = strings.TrimSpace(in.PlateNo)
	if len(in.PlateNo) > 20 {
		return nil, fieldError{"plate_no", "too long"}
	}
	active := in.Active == nil || *in.Active
	var userID string
	err = db.WithActor(ctx, s.Pool, adminID, "admin", "", func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx, `SELECT id FROM users WHERE phone = $1 AND deleted_at IS NULL`, phone).Scan(&userID)
		if errors.Is(err, pgx.ErrNoRows) {
			err = tx.QueryRow(ctx, `INSERT INTO users (name, phone, language) VALUES ($1, $2, 'km') RETURNING id`, in.Name, phone).Scan(&userID)
		}
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `INSERT INTO couriers (user_id, display_name, vehicle, plate_no, active, created_by) VALUES ($1, $2, $3, nullif($4, ''), $5, $6)
			ON CONFLICT (user_id) DO UPDATE SET display_name = EXCLUDED.display_name, vehicle = EXCLUDED.vehicle, plate_no = EXCLUDED.plate_no,
			  active = EXCLUDED.active, is_online = couriers.is_online AND EXCLUDED.active`, userID, in.Name, in.Vehicle, in.PlateNo, active, adminID)
		return err
	})
	if err != nil {
		return nil, err
	}
	return scanCourier(s.Pool.QueryRow(ctx, `SELECT `+courierCols+` FROM couriers c JOIN users u ON u.id = c.user_id WHERE c.user_id = $1`, userID))
}

// AdminDeliveries: status = open (belum selesai) · failed · delivered (24 jam terakhir).
func (s *Service) AdminDeliveries(ctx context.Context, status string) ([]Task, error) {
	switch status {
	case "failed":
		return s.tasks(ctx, `d.status = 'failed'`, `d.failed_at`, 200)
	case "delivered":
		return s.tasks(ctx, `d.status = 'delivered' AND d.delivered_at > now() - interval '24 hours'`, `d.delivered_at DESC`, 200)
	default:
		return s.tasks(ctx, `d.status IN ('unassigned','assigned','picked_up')`, `d.status = 'unassigned' DESC, d.created_at`, 200)
	}
}

// Assign: admin menugaskan / memindahkan tugas ke kurir lain (courierID kosong = kembalikan ke daftar tersedia).
// Paket yang sudah dijemput hanya bisa dipindah ke kurir lain (bukan dilepas).
func (s *Service) Assign(ctx context.Context, adminID, taskID, courierID string) (*Task, error) {
	err := db.WithActor(ctx, s.Pool, adminID, "admin", "", func(tx pgx.Tx) error {
		l, err := lockTask(ctx, tx, taskID)
		if err != nil {
			return err
		}
		if courierID == "" {
			if l.status != "assigned" {
				return ErrBadStatus
			}
			if _, err := tx.Exec(ctx, `UPDATE deliveries SET status = 'unassigned', courier_id = NULL, assigned_at = NULL WHERE id = $1`, l.id); err != nil {
				return err
			}
			return event(ctx, tx, l.orderID, "courier_released", "admin", adminID, "")
		}
		var name string
		if err := tx.QueryRow(ctx, `SELECT display_name FROM couriers WHERE user_id = $1 AND active`, courierID).Scan(&name); err != nil {
			return fieldError{"courier_id", "choose an active courier"}
		}
		switch l.status {
		case "unassigned", "assigned":
			_, err = tx.Exec(ctx, `UPDATE deliveries SET status = 'assigned', courier_id = $2, assigned_at = now() WHERE id = $1`, l.id, courierID)
		case "picked_up", "failed":
			_, err = tx.Exec(ctx, `UPDATE deliveries SET courier_id = $2 WHERE id = $1`, l.id, courierID)
		default:
			return ErrBadStatus
		}
		if err != nil {
			return err
		}
		if err := event(ctx, tx, l.orderID, "courier_assigned", "admin", adminID, name); err != nil {
			return err
		}
		return notify.EnqueueUser(ctx, tx, courierID, notify.CourierNewTask, notify.Params{Store: l.storeName, Order: l.orderNo}, "/courier/")
	})
	if err != nil {
		return nil, err
	}
	list, err := s.tasks(ctx, `d.id = $1`, `d.id`, 1, taskID)
	if err != nil || len(list) == 0 {
		return nil, err
	}
	return &list[0], nil
}

// ProofKey: kunci file foto bukti (untuk admin).
func (s *Service) ProofKey(ctx context.Context, taskID string) (string, error) {
	var key *string
	err := s.Pool.QueryRow(ctx, `SELECT proof_key FROM deliveries WHERE id = $1`, taskID).Scan(&key)
	if errors.Is(err, pgx.ErrNoRows) || (err == nil && key == nil) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", err
	}
	return *key, nil
}
