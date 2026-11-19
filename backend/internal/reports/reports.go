// Package reports: laporan pembeli tentang produk/toko dan strike toko.
//
//	pembeli melapor (produk atau toko) ──► antrean admin (dikelompokkan per target)
//	   admin: abaikan · larang produk (status banned) · beri strike ke toko
//	strike aktif 90 hari; total poin aktif ≥ 3 → toko otomatis "suspended" + seller diberi tahu.
//	strike bisa dicabut; kalau penangguhan disebabkan strike dan poin aktif turun di bawah 3 → toko aktif lagi.
package reports

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"unicode/utf8"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"phsar/internal/modrules"
	"phsar/internal/notify"
	"phsar/internal/platform/db"
)

const (
	SuspendAt   = 3  // poin strike aktif yang membuat toko ditangguhkan
	DailyLimit  = 10 // laporan per pembeli per 24 jam
	maxNoteRune = 1000
)

var (
	ErrNotFound = errors.New("tidak ditemukan")
	ErrLimit    = errors.New("terlalu banyak laporan, coba lagi besok")
	ErrOpen     = errors.New("kamu sudah melaporkan ini, tim Phsar sedang memeriksa")
	ErrNothing  = errors.New("tidak ada laporan terbuka untuk target ini")
	ErrRevoked  = errors.New("strike sudah dicabut")
)

// FieldError: input tidak valid (422).
type FieldError struct{ Field, Msg string }

func (e FieldError) Error() string { return e.Field + ": " + e.Msg }

var Reasons = []string{"counterfeit", "prohibited", "misleading", "scam", "offensive", "other"}

func validReason(r string) bool {
	for _, x := range Reasons {
		if x == r {
			return true
		}
	}
	return false
}

type Service struct {
	Pool *pgxpool.Pool
	Log  *slog.Logger
}

// ---------------------------------------------------------------------------
// Pembeli
// ---------------------------------------------------------------------------

type ReportInput struct {
	Target string `json:"target"` // product · store
	ID     string `json:"id"`
	Reason string `json:"reason"`
	Note   string `json:"note"`
}

func (s *Service) Create(ctx context.Context, userID string, in ReportInput) (string, error) {
	in.Note = strings.TrimSpace(in.Note)
	if in.Target != "product" && in.Target != "store" {
		return "", FieldError{"target", "product or store"}
	}
	if !validReason(in.Reason) {
		return "", FieldError{"reason", "choose a reason"}
	}
	if utf8.RuneCountInString(in.Note) > maxNoteRune {
		return "", FieldError{"note", "too long"}
	}
	if in.Reason == "other" && utf8.RuneCountInString(in.Note) < 5 {
		return "", FieldError{"note", "tell us what is wrong"}
	}
	var id string
	err := db.WithActor(ctx, s.Pool, userID, "buyer", "", func(tx pgx.Tx) error {
		var n int
		if err := tx.QueryRow(ctx, `SELECT count(*) FROM reports WHERE reporter_id = $1 AND created_at > now() - interval '24 hours'`, userID).Scan(&n); err != nil {
			return err
		}
		if n >= DailyLimit {
			return ErrLimit
		}
		var productID *string
		var storeID string
		var err error
		if in.Target == "product" {
			productID = &in.ID
			err = tx.QueryRow(ctx, `SELECT p.store_id FROM products p WHERE p.id::text = $1 AND p.deleted_at IS NULL`, in.ID).Scan(&storeID)
		} else {
			err = tx.QueryRow(ctx, `SELECT id FROM stores WHERE id::text = $1 AND deleted_at IS NULL`, in.ID).Scan(&storeID)
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		err = tx.QueryRow(ctx, `INSERT INTO reports (reporter_id, target_type, product_id, store_id, reason, note)
			VALUES ($1, $2, $3, $4, $5, nullif($6, '')) ON CONFLICT DO NOTHING RETURNING id`,
			userID, in.Target, productID, storeID, in.Reason, in.Note).Scan(&id)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrOpen
		}
		if err != nil || productID == nil || !seriousReason[in.Reason] {
			return err
		}
		return holdIfReported(ctx, tx, *productID)
	})
	return id, err
}

// Laporan yang cukup serius untuk menahan produk sementara.
var seriousReason = map[string]bool{"counterfeit": true, "prohibited": true, "scam": true, "misleading": true}

// holdIfReported: laporan serius terbuka dari N pembeli berbeda (moderation_rules.report_hold) → produk yang tayang
// disembunyikan & masuk antrean moderasi. Admin menyetujui lagi atau melarang; edit seller tidak melepasnya.
func holdIfReported(ctx context.Context, tx pgx.Tx, productID string) error {
	rules, err := modrules.Load(ctx, tx)
	if err != nil {
		return err
	}
	var n int
	if err := tx.QueryRow(ctx, `SELECT count(DISTINCT reporter_id) FROM reports
		WHERE product_id = $1 AND status = 'open' AND reason IN ('counterfeit','prohibited','scam','misleading')`, productID).Scan(&n); err != nil {
		return err
	}
	if n < rules.ReportHold {
		return nil
	}
	_, err = tx.Exec(ctx, `UPDATE products SET moderation_status = 'pending_review', review_lane = 'manual', post_review_pending = false,
		robots_index = false, submitted_at = now(),
		auto_flags = auto_flags || jsonb_build_object('lane', 'manual', 'reports', $2::int,
			'reasons', (SELECT coalesce(jsonb_agg(DISTINCT x), '[]') FROM jsonb_array_elements_text(CASE WHEN jsonb_typeof(auto_flags->'reasons') = 'array' THEN auto_flags->'reasons' ELSE '[]' END
				|| '["reported"]') x))
		WHERE id = $1 AND moderation_status = 'approved' AND deleted_at IS NULL`, productID, n)
	return err
}

// ---------------------------------------------------------------------------
// Admin: antrean
// ---------------------------------------------------------------------------

// Groups: laporan dikelompokkan per target (produk / toko) — status open atau closed.
func (s *Service) Groups(ctx context.Context, status string) (json.RawMessage, error) {
	where := "r.status = 'open'"
	order := "count(*) DESC, min(r.created_at)"
	if status == "closed" {
		where = "r.status <> 'open' AND r.resolved_at > now() - interval '60 days'"
		order = "max(r.resolved_at) DESC"
	}
	var raw []byte
	err := s.Pool.QueryRow(ctx, `SELECT coalesce(json_agg(g), '[]') FROM (
		SELECT r.target_type, coalesce(r.product_id, r.store_id) AS target_id, r.store_id, st.name AS store_name, st.slug AS store_slug,
		  st.status AS store_status, r.product_id,
		  (SELECT coalesce(t.name, '') FROM products p LEFT JOIN product_translations t ON t.product_id = p.id AND t.lang = p.source_lang
		     WHERE p.id = r.product_id) AS product_name,
		  (SELECT p.status FROM products p WHERE p.id = r.product_id) AS product_status,
		  (SELECT t.slug FROM product_translations t WHERE t.product_id = r.product_id ORDER BY (t.lang = 'en') DESC LIMIT 1) AS product_slug,
		  count(*) AS count, min(r.created_at) AS first_at, max(r.created_at) AS last_at,
		  (SELECT coalesce(sum(k.points), 0) FROM store_strikes k WHERE k.store_id = r.store_id AND k.revoked_at IS NULL AND k.expires_at > now()) AS active_strikes,
		  json_agg(json_build_object('id', r.id, 'reason', r.reason, 'note', r.note, 'at', r.created_at, 'status', r.status,
		    'resolution', r.resolution, 'resolution_note', r.resolution_note, 'resolved_at', r.resolved_at,
		    'reporter', coalesce(nullif(u.name, ''), u.phone, 'user'), 'reporter_phone', u.phone) ORDER BY r.created_at) AS reports
		FROM reports r JOIN stores st ON st.id = r.store_id JOIN users u ON u.id = r.reporter_id
		WHERE `+where+`
		GROUP BY r.target_type, coalesce(r.product_id, r.store_id), r.store_id, st.name, st.slug, st.status, r.product_id
		ORDER BY `+order+` LIMIT 200) g`).Scan(&raw)
	return raw, err
}

func (s *Service) OpenCount(ctx context.Context) (int, error) {
	var n int
	err := s.Pool.QueryRow(ctx, `SELECT count(DISTINCT (target_type, coalesce(product_id, store_id))) FROM reports WHERE status = 'open'`).Scan(&n)
	return n, err
}

type Resolution struct {
	Target       string `json:"target"` // product · store
	ID           string `json:"id"`
	Action       string `json:"action"` // dismiss · ban_product · strike
	Note         string `json:"note"`
	AlsoStrike   bool   `json:"also_strike"`   // ban_product + strike sekaligus
	StrikeReason string `json:"strike_reason"` // default: alasan laporan terbanyak
	Points       int    `json:"points"`        // 1–3, default 1
}

// Resolve menutup semua laporan terbuka untuk satu target sekaligus.
func (s *Service) Resolve(ctx context.Context, adminID string, in Resolution) error {
	in.Note = strings.TrimSpace(in.Note)
	if in.Target != "product" && in.Target != "store" {
		return FieldError{"target", "product or store"}
	}
	if in.Action != "dismiss" && in.Action != "ban_product" && in.Action != "strike" {
		return FieldError{"action", "dismiss, ban_product or strike"}
	}
	if in.Action == "ban_product" && in.Target != "product" {
		return FieldError{"action", "only products can be banned"}
	}
	if n := utf8.RuneCountInString(in.Note); n < 5 || n > maxNoteRune {
		return FieldError{"note", "write a note (5–1000 characters)"}
	}
	strike := in.Action == "strike" || (in.Action == "ban_product" && in.AlsoStrike)
	if in.StrikeReason != "" && !validReason(in.StrikeReason) {
		return FieldError{"strike_reason", "invalid reason"}
	}
	if in.Points == 0 {
		in.Points = 1
	}
	if in.Points < 1 || in.Points > 3 {
		return FieldError{"points", "1–3"}
	}
	return db.WithActor(ctx, s.Pool, adminID, "admin", "", func(tx pgx.Tx) error {
		col := "store_id"
		if in.Target == "product" {
			col = "product_id"
		}
		rows, err := tx.Query(ctx, `SELECT id, store_id, reason FROM reports WHERE status = 'open' AND target_type = $1 AND `+col+`::text = $2 FOR UPDATE`, in.Target, in.ID)
		if err != nil {
			return err
		}
		var storeID string
		count := map[string]int{}
		n := 0
		for rows.Next() {
			var id, st, reason string
			if err := rows.Scan(&id, &st, &reason); err != nil {
				rows.Close()
				return err
			}
			storeID = st
			count[reason]++
			n++
		}
		rows.Close()
		if n == 0 {
			return ErrNothing
		}
		status := "actioned"
		if in.Action == "dismiss" {
			status = "dismissed"
		}
		if _, err := tx.Exec(ctx, `UPDATE reports SET status = $3, resolution = $4, resolution_note = $5, resolved_by = $6, resolved_at = now()
			WHERE status = 'open' AND target_type = $1 AND `+col+`::text = $2`, in.Target, in.ID, status, in.Action, in.Note, adminID); err != nil {
			return err
		}
		if in.Action == "ban_product" {
			if err := banProduct(ctx, tx, adminID, in.ID, in.Note); err != nil {
				return err
			}
		}
		if strike {
			reason := in.StrikeReason
			if reason == "" {
				best := 0
				for _, r := range Reasons { // urutan tetap supaya hasil seri selalu sama
					if count[r] > best {
						reason, best = r, count[r]
					}
				}
			}
			var productID *string
			if in.Target == "product" {
				productID = &in.ID
			}
			if _, err := addStrike(ctx, tx, adminID, storeID, reason, in.Note, in.Points, productID); err != nil {
				return err
			}
		}
		return nil
	})
}

// banProduct: produk dilarang dijual (status banned + moderasi rejected) dan seller diberi tahu.
func banProduct(ctx context.Context, tx pgx.Tx, adminID, productID, note string) error {
	var storeID, name, storeName string
	err := tx.QueryRow(ctx, `SELECT p.store_id, coalesce(t.name, ''), st.name FROM products p JOIN stores st ON st.id = p.store_id
		LEFT JOIN product_translations t ON t.product_id = p.id AND t.lang = p.source_lang WHERE p.id::text = $1 FOR UPDATE OF p`, productID).Scan(&storeID, &name, &storeName)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `UPDATE products SET status = 'banned', moderation_status = 'rejected', moderation_note = $2, robots_index = false
		WHERE id = $1`, productID, note); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `INSERT INTO product_moderation_reviews (product_id, decision, reason_code, note, reviewer_id)
		VALUES ($1, 'rejected', 'prohibited', $2, $3)`, productID, note, adminID); err != nil {
		return err
	}
	return notify.Enqueue(ctx, tx, storeID, notify.ProductRejected, notify.Params{Store: storeName, Product: name, Note: note}, "#/store/"+storeID+"/product/"+productID)
}

// ---------------------------------------------------------------------------
// Strike
// ---------------------------------------------------------------------------

func activePoints(ctx context.Context, tx pgx.Tx, storeID string) (int, error) {
	var n int
	err := tx.QueryRow(ctx, `SELECT coalesce(sum(points), 0) FROM store_strikes WHERE store_id = $1 AND revoked_at IS NULL AND expires_at > now()`, storeID).Scan(&n)
	return n, err
}

// addStrike menambah strike; kalau poin aktif ≥ 3 dan toko aktif → ditangguhkan.
func addStrike(ctx context.Context, tx pgx.Tx, adminID, storeID, reason, note string, points int, productID *string) (string, error) {
	var storeName, status string
	err := tx.QueryRow(ctx, `SELECT name, status::text FROM stores WHERE id::text = $1 AND deleted_at IS NULL FOR UPDATE`, storeID).Scan(&storeName, &status)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", err
	}
	var id string
	if err := tx.QueryRow(ctx, `INSERT INTO store_strikes (store_id, reason, note, points, product_id, created_by) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		storeID, reason, note, points, productID, adminID).Scan(&id); err != nil {
		return "", err
	}
	active, err := activePoints(ctx, tx, storeID)
	if err != nil {
		return "", err
	}
	link := "#/store/" + storeID
	if err := notify.Enqueue(ctx, tx, storeID, notify.StoreStrike, notify.Params{Store: storeName, Note: note, Point: fmt.Sprint(min(active, SuspendAt))}, link); err != nil {
		return "", err
	}
	if active >= SuspendAt && status == "active" {
		if _, err := tx.Exec(ctx, `UPDATE stores SET status = 'suspended', review_note = $2 WHERE id = $1`, storeID,
			fmt.Sprintf("Suspended automatically after %d active strikes.", SuspendAt)); err != nil {
			return "", err
		}
		if _, err := tx.Exec(ctx, `UPDATE store_strikes SET caused_suspension = true WHERE id = $1`, id); err != nil {
			return "", err
		}
		if err := notify.Enqueue(ctx, tx, storeID, notify.StoreSuspended, notify.Params{Store: storeName}, link); err != nil {
			return "", err
		}
	}
	return id, nil
}

type StrikeInput struct {
	StoreID string `json:"store_id"`
	Reason  string `json:"reason"`
	Note    string `json:"note"`
	Points  int    `json:"points"`
}

// AddStrike: strike manual (tanpa laporan), mis. dari temuan tim Phsar sendiri.
func (s *Service) AddStrike(ctx context.Context, adminID string, in StrikeInput) (string, error) {
	in.Note = strings.TrimSpace(in.Note)
	if !validReason(in.Reason) {
		return "", FieldError{"reason", "choose a reason"}
	}
	if n := utf8.RuneCountInString(in.Note); n < 5 || n > maxNoteRune {
		return "", FieldError{"note", "write a note for the seller (5–1000 characters)"}
	}
	if in.Points == 0 {
		in.Points = 1
	}
	if in.Points < 1 || in.Points > 3 {
		return "", FieldError{"points", "1–3"}
	}
	var id string
	err := db.WithActor(ctx, s.Pool, adminID, "admin", "", func(tx pgx.Tx) error {
		var err error
		id, err = addStrike(ctx, tx, adminID, in.StoreID, in.Reason, in.Note, in.Points, nil)
		return err
	})
	return id, err
}

// Revoke mencabut strike. Kalau toko ditangguhkan karena strike dan poin aktif < 3 → aktif lagi.
func (s *Service) Revoke(ctx context.Context, adminID, strikeID, note string) (string, error) {
	note = strings.TrimSpace(note)
	if n := utf8.RuneCountInString(note); n < 5 || n > maxNoteRune {
		return "", FieldError{"note", "explain why (5–1000 characters)"}
	}
	var storeID string
	err := db.WithActor(ctx, s.Pool, adminID, "admin", "", func(tx pgx.Tx) error {
		var revoked bool
		err := tx.QueryRow(ctx, `SELECT store_id, revoked_at IS NOT NULL FROM store_strikes WHERE id::text = $1 FOR UPDATE`, strikeID).Scan(&storeID, &revoked)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		if revoked {
			return ErrRevoked
		}
		var storeName, status string
		if err := tx.QueryRow(ctx, `SELECT name, status::text FROM stores WHERE id = $1 FOR UPDATE`, storeID).Scan(&storeName, &status); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `UPDATE store_strikes SET revoked_at = now(), revoked_by = $2, revoke_note = $3 WHERE id = $1`, strikeID, adminID, note); err != nil {
			return err
		}
		active, err := activePoints(ctx, tx, storeID)
		if err != nil {
			return err
		}
		var byStrikes bool
		if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM store_strikes WHERE store_id = $1 AND caused_suspension)`, storeID).Scan(&byStrikes); err != nil {
			return err
		}
		if status == "suspended" && byStrikes && active < SuspendAt {
			if _, err := tx.Exec(ctx, `UPDATE stores SET status = 'active', review_note = NULL WHERE id = $1`, storeID); err != nil {
				return err
			}
			if _, err := tx.Exec(ctx, `UPDATE store_strikes SET caused_suspension = false WHERE store_id = $1`, storeID); err != nil {
				return err
			}
			return notify.Enqueue(ctx, tx, storeID, notify.StoreReinstated, notify.Params{Store: storeName, Note: note}, "#/store/"+storeID)
		}
		return nil
	})
	return storeID, err
}

// Strikes: riwayat strike satu toko (untuk admin; forSeller menyembunyikan siapa yang memberi).
func (s *Service) Strikes(ctx context.Context, storeID string, forSeller bool) (json.RawMessage, error) {
	by := "coalesce(nullif(u.name, ''), u.phone)"
	if forSeller {
		by = "NULL"
	}
	var raw []byte
	err := s.Pool.QueryRow(ctx, `SELECT json_build_object(
		'store', (SELECT json_build_object('id', st.id, 'name', st.name, 'status', st.status) FROM stores st WHERE st.id::text = $1),
		'active_points', (SELECT coalesce(sum(points), 0) FROM store_strikes WHERE store_id::text = $1 AND revoked_at IS NULL AND expires_at > now()),
		'suspend_at', `+fmt.Sprint(SuspendAt)+`,
		'strikes', (SELECT coalesce(json_agg(json_build_object('id', k.id, 'reason', k.reason, 'note', k.note, 'points', k.points,
			'created_at', k.created_at, 'expires_at', k.expires_at, 'revoked_at', k.revoked_at, 'revoke_note', k.revoke_note,
			'active', k.revoked_at IS NULL AND k.expires_at > now(), 'caused_suspension', k.caused_suspension, 'by', `+by+`,
			'auto', k.auto, 'order_no', (SELECT o.order_no FROM orders o WHERE o.id = k.order_id)) -- extras: strike otomatis
			ORDER BY k.created_at DESC), '[]') FROM store_strikes k LEFT JOIN users u ON u.id = k.created_by WHERE k.store_id::text = $1))`, storeID).Scan(&raw)
	return raw, err
}
