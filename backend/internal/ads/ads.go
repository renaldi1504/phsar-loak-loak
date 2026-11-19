// Package ads: iklan produk bersponsor (Phase 2 MVP).
//
//	seller (pemilik toko) buat kampanye untuk satu produk tayang: CPC $0.05–$2, anggaran harian ≥ $1, tanggal mulai/selesai
//	→ admin tinjau (pending → approved / rejected) → tayang di hasil cari & halaman kategori, di atas hasil biasa
//	→ klik lewat /ad/{id}: sekali tagih per pengunjung per kampanye per hari (bot diabaikan), dipotong dari saldo seller
//	  (seller_ledger kind ad_spend). Berhenti tayang kalau anggaran hari ini habis atau saldo seller tidak cukup.
package ads

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"phsar/internal/notify"
	"phsar/internal/platform/db"
)

var (
	ErrNotFound = errors.New("kampanye tidak ditemukan")
	ErrState    = errors.New("status kampanye tidak memungkinkan aksi ini")
	MinBid      = decimal.RequireFromString("0.05")
	MaxBid      = decimal.RequireFromString("2.00")
	MinBudget   = decimal.RequireFromString("1.00")
)

// FieldError: input tidak valid (422).
type FieldError struct{ Field, Msg string }

func (e FieldError) Error() string { return e.Field + ": " + e.Msg }

type Service struct {
	Pool *pgxpool.Pool
	Log  *slog.Logger
}

var phnomPenh = func() *time.Location {
	if l, err := time.LoadLocation("Asia/Phnom_Penh"); err == nil {
		return l
	}
	return time.FixedZone("ICT", 7*3600)
}()

// today: tanggal hari ini di Phnom Penh (anggaran harian & statistik memakai tanggal ini).
func today() string { return time.Now().In(phnomPenh).Format("2006-01-02") }

// ---------------------------------------------------------------------------
// Pengaturan global
// ---------------------------------------------------------------------------

type Settings struct {
	Enabled bool `json:"enabled"`
	Slots   int  `json:"slots"` // slot bersponsor per halaman (1–4)
}

func (s *Service) Settings(ctx context.Context) (Settings, error) {
	st := Settings{Enabled: true, Slots: 2}
	var raw []byte
	err := s.Pool.QueryRow(ctx, `SELECT value FROM platform_settings WHERE key = 'ads'`).Scan(&raw)
	if errors.Is(err, pgx.ErrNoRows) {
		return st, nil
	}
	if err != nil {
		return st, err
	}
	_ = json.Unmarshal(raw, &st)
	if st.Slots < 1 || st.Slots > 4 {
		st.Slots = 2
	}
	return st, nil
}

func (s *Service) SaveSettings(ctx context.Context, actorID string, in Settings) (Settings, error) {
	if in.Slots < 1 || in.Slots > 4 {
		return in, FieldError{"slots", "1–4"}
	}
	v, _ := json.Marshal(in)
	err := db.WithActor(ctx, s.Pool, actorID, "admin", "", func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `INSERT INTO platform_settings (key, value, updated_by, updated_at) VALUES ('ads', $1, $2, now())
			ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_by = EXCLUDED.updated_by, updated_at = now()`, v, actorID)
		return err
	})
	return in, err
}

// ---------------------------------------------------------------------------
// Model
// ---------------------------------------------------------------------------

type Stats struct {
	Impressions int    `json:"impressions"`
	Clicks      int    `json:"clicks"`
	Spend       string `json:"spend"`
}

type Campaign struct {
	ID          string     `json:"id"`
	StoreID     string     `json:"store_id"`
	StoreName   string     `json:"store_name"`
	ProductID   string     `json:"product_id"`
	ProductName string     `json:"product_name"`
	ProductURL  string     `json:"product_image,omitempty"`
	Status      string     `json:"status"`
	Live        bool       `json:"live"` // sedang tayang sekarang (disetujui, dalam tanggal, anggaran & saldo cukup)
	PausedBy    *string    `json:"paused_by"`
	Bid         string     `json:"bid"`
	DailyBudget string     `json:"daily_budget"`
	StartsOn    string     `json:"starts_on"`
	EndsOn      *string    `json:"ends_on"`
	ReviewNote  *string    `json:"review_note"`
	ReviewedAt  *time.Time `json:"reviewed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	Today       Stats      `json:"today"`
	Total       Stats      `json:"total"`
}

const selectCampaign = `
	SELECT c.id, c.store_id, st.name, c.product_id, coalesce(t.name, ''),
	  coalesce((SELECT coalesce(m.image_url, '') FROM product_media m WHERE m.product_id = c.product_id AND m.deleted_at IS NULL AND m.kind = 'photo' ORDER BY m.sort_order LIMIT 1), ''),
	  c.status, c.paused_by, c.bid::text, c.daily_budget::text, c.starts_on::text, c.ends_on::text, c.review_note, c.reviewed_at, c.created_at,
	  coalesce(d.impressions, 0), coalesce(d.clicks, 0), coalesce(d.spend, 0)::text,
	  coalesce(tot.impressions, 0), coalesce(tot.clicks, 0), coalesce(tot.spend, 0)::text,
	  (c.status = 'approved' AND c.starts_on <= $1::date AND (c.ends_on IS NULL OR c.ends_on >= $1::date)
	    AND coalesce(d.spend, 0) + c.bid <= c.daily_budget
	    AND (SELECT coalesce(sum(l.amount), 0) FROM seller_ledger l WHERE l.store_id = c.store_id) >= c.bid)
	FROM ad_campaigns c JOIN stores st ON st.id = c.store_id JOIN products p ON p.id = c.product_id
	LEFT JOIN product_translations t ON t.product_id = p.id AND t.lang = p.source_lang
	LEFT JOIN ad_stats_daily d ON d.campaign_id = c.id AND d.day = $1::date
	LEFT JOIN LATERAL (SELECT sum(impressions) AS impressions, sum(clicks) AS clicks, sum(spend) AS spend FROM ad_stats_daily x WHERE x.campaign_id = c.id) tot ON true`

func (s *Service) query(ctx context.Context, where string, limit int, args ...any) ([]Campaign, error) {
	rows, err := s.Pool.Query(ctx, selectCampaign+" WHERE "+where+" ORDER BY c.created_at DESC LIMIT "+fmt.Sprint(limit), append([]any{today()}, args...)...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Campaign{}
	for rows.Next() {
		var c Campaign
		if err := rows.Scan(&c.ID, &c.StoreID, &c.StoreName, &c.ProductID, &c.ProductName, &c.ProductURL, &c.Status, &c.PausedBy, &c.Bid, &c.DailyBudget,
			&c.StartsOn, &c.EndsOn, &c.ReviewNote, &c.ReviewedAt, &c.CreatedAt, &c.Today.Impressions, &c.Today.Clicks, &c.Today.Spend,
			&c.Total.Impressions, &c.Total.Clicks, &c.Total.Spend, &c.Live); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (s *Service) one(ctx context.Context, where string, args ...any) (*Campaign, error) {
	l, err := s.query(ctx, where, 1, args...)
	if err != nil {
		return nil, err
	}
	if len(l) == 0 {
		return nil, ErrNotFound
	}
	return &l[0], nil
}

// ---------------------------------------------------------------------------
// Seller (pemilik toko)
// ---------------------------------------------------------------------------

type Input struct {
	ProductID   string  `json:"product_id"`
	Bid         string  `json:"bid"`
	DailyBudget string  `json:"daily_budget"`
	StartsOn    string  `json:"starts_on"`
	EndsOn      *string `json:"ends_on"`
}

func money(v, field string, min, max decimal.Decimal) (decimal.Decimal, error) {
	d, err := decimal.NewFromString(strings.TrimSpace(strings.TrimPrefix(v, "$")))
	if err != nil || d.Exponent() < -2 || d.LessThan(min) || (!max.IsZero() && d.GreaterThan(max)) {
		if max.IsZero() {
			return d, FieldError{field, "at least $" + min.StringFixed(2)}
		}
		return d, FieldError{field, "$" + min.StringFixed(2) + "–$" + max.StringFixed(2)}
	}
	return d, nil
}

func dates(starts string, ends *string) (string, *string, error) {
	t := today()
	if starts == "" {
		starts = t
	}
	sd, err := time.Parse("2006-01-02", starts)
	if err != nil {
		return "", nil, FieldError{"starts_on", "date YYYY-MM-DD"}
	}
	if starts < t {
		return "", nil, FieldError{"starts_on", "cannot start in the past"}
	}
	if sd.After(time.Now().AddDate(1, 0, 0)) {
		return "", nil, FieldError{"starts_on", "at most one year ahead"}
	}
	if ends != nil && *ends != "" {
		if _, err := time.Parse("2006-01-02", *ends); err != nil || *ends < starts {
			return "", nil, FieldError{"ends_on", "on or after the start date"}
		}
		return starts, ends, nil
	}
	return starts, nil, nil
}

func (s *Service) SellerList(ctx context.Context, storeID string) ([]Campaign, error) {
	return s.query(ctx, "c.store_id = $2", 100, storeID)
}

func (s *Service) Create(ctx context.Context, storeID, userID string, in Input) (*Campaign, error) {
	bid, err := money(in.Bid, "bid", MinBid, MaxBid)
	if err != nil {
		return nil, err
	}
	budget, err := money(in.DailyBudget, "daily_budget", MinBudget, decimal.RequireFromString("1000"))
	if err != nil {
		return nil, err
	}
	if budget.LessThan(bid) {
		return nil, FieldError{"daily_budget", "must be at least one click (the bid)"}
	}
	starts, ends, err := dates(in.StartsOn, in.EndsOn)
	if err != nil {
		return nil, err
	}
	var id string
	err = db.WithActor(ctx, s.Pool, userID, "seller", "", func(tx pgx.Tx) error {
		var ok bool
		// produk harus milik toko ini dan sedang tayang (aktif + disetujui moderasi)
		if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM products WHERE id::text = $1 AND store_id = $2 AND deleted_at IS NULL
			AND status = 'active' AND moderation_status = 'approved')`, in.ProductID, storeID).Scan(&ok); err != nil {
			return err
		}
		if !ok {
			return FieldError{"product_id", "choose an approved, active product from your store"}
		}
		err := tx.QueryRow(ctx, `INSERT INTO ad_campaigns (store_id, product_id, bid, daily_budget, starts_on, ends_on, created_by)
			VALUES ($1, $2, $3, $4, $5::date, $6::date, $7) RETURNING id`, storeID, in.ProductID, bid, budget, starts, ends, userID).Scan(&id)
		if err != nil && strings.Contains(err.Error(), "ad_campaigns_one_per_product") {
			return FieldError{"product_id", "this product already has a campaign"}
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return s.one(ctx, "c.id = $2", id)
}

type Update struct {
	Bid         *string `json:"bid"`
	DailyBudget *string `json:"daily_budget"`
	EndsOn      *string `json:"ends_on"`
	Action      string  `json:"action"` // pause · resume · end
}

// SellerUpdate: ubah CPC / anggaran / tanggal selesai, atau jeda / lanjut / akhiri. Perubahan angka tidak perlu ditinjau ulang.
func (s *Service) SellerUpdate(ctx context.Context, storeID, userID, id string, in Update) (*Campaign, error) {
	err := db.WithActor(ctx, s.Pool, userID, "seller", "", func(tx pgx.Tx) error {
		var status string
		var pausedBy *string
		var bid, budget decimal.Decimal
		var starts string
		err := tx.QueryRow(ctx, `SELECT status, paused_by, bid, daily_budget, starts_on::text FROM ad_campaigns WHERE id::text = $1 AND store_id = $2 FOR UPDATE`,
			id, storeID).Scan(&status, &pausedBy, &bid, &budget, &starts)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		if status == "ended" || status == "rejected" {
			return ErrState
		}
		if in.Bid != nil {
			if bid, err = money(*in.Bid, "bid", MinBid, MaxBid); err != nil {
				return err
			}
		}
		if in.DailyBudget != nil {
			if budget, err = money(*in.DailyBudget, "daily_budget", MinBudget, decimal.RequireFromString("1000")); err != nil {
				return err
			}
		}
		if budget.LessThan(bid) {
			return FieldError{"daily_budget", "must be at least one click (the bid)"}
		}
		if in.EndsOn != nil && *in.EndsOn != "" {
			if _, err := time.Parse("2006-01-02", *in.EndsOn); err != nil || *in.EndsOn < starts || *in.EndsOn < today() {
				return FieldError{"ends_on", "today or later, on or after the start date"}
			}
		}
		newStatus, newPaused := status, pausedBy
		seller := "seller"
		switch in.Action {
		case "":
		case "pause":
			if status != "approved" && status != "pending" {
				return ErrState
			}
			if status == "approved" {
				newStatus, newPaused = "paused", &seller
			}
		case "resume":
			if status != "paused" {
				return ErrState
			}
			if pausedBy != nil && *pausedBy == "admin" {
				return ErrState // dijeda Phsar: hanya admin yang boleh melanjutkan
			}
			newStatus, newPaused = "approved", nil
		case "end":
			newStatus = "ended"
		default:
			return FieldError{"action", "pause, resume or end"}
		}
		_, err = tx.Exec(ctx, `UPDATE ad_campaigns SET bid = $2, daily_budget = $3, ends_on = CASE WHEN $4::text IS NULL THEN ends_on WHEN $4 = '' THEN NULL ELSE $4::date END,
			status = $5, paused_by = $6 WHERE id = $1`, id, bid, budget, in.EndsOn, newStatus, newPaused)
		return err
	})
	if err != nil {
		return nil, err
	}
	return s.one(ctx, "c.id::text = $2", id)
}

// ---------------------------------------------------------------------------
// Admin
// ---------------------------------------------------------------------------

type Counts struct {
	Pending int `json:"pending"`
	Running int `json:"running"`
}

func (s *Service) AdminList(ctx context.Context, filter string) ([]Campaign, Counts, error) {
	var c Counts
	if err := s.Pool.QueryRow(ctx, `SELECT count(*) FILTER (WHERE status = 'pending'), count(*) FILTER (WHERE status IN ('approved','paused')) FROM ad_campaigns`).
		Scan(&c.Pending, &c.Running); err != nil {
		return nil, c, err
	}
	where := "c.status = 'pending'"
	switch filter {
	case "running":
		where = "c.status IN ('approved','paused')"
	case "all":
		where = "TRUE"
	}
	l, err := s.query(ctx, where, 300)
	return l, c, err
}

type Review struct {
	Action string `json:"action"` // approve · reject · pause · resume
	Note   string `json:"note"`
}

func (s *Service) AdminAction(ctx context.Context, adminID, id string, in Review) (*Campaign, error) {
	in.Note = strings.TrimSpace(in.Note)
	if len([]rune(in.Note)) > 300 {
		return nil, FieldError{"note", "at most 300 characters"}
	}
	if (in.Action == "reject" || in.Action == "pause") && len([]rune(in.Note)) < 5 {
		return nil, FieldError{"note", "tell the seller why (at least 5 characters)"}
	}
	err := db.WithActor(ctx, s.Pool, adminID, "admin", "", func(tx pgx.Tx) error {
		var status, storeID, product string
		err := tx.QueryRow(ctx, `SELECT c.status, c.store_id, coalesce(t.name, '') FROM ad_campaigns c JOIN products p ON p.id = c.product_id
			LEFT JOIN product_translations t ON t.product_id = p.id AND t.lang = p.source_lang WHERE c.id::text = $1 FOR UPDATE OF c`, id).Scan(&status, &storeID, &product)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		var q, kind string
		switch {
		case in.Action == "approve" && status == "pending":
			q, kind = `UPDATE ad_campaigns SET status = 'approved', review_note = nullif($2, ''), reviewed_by = $3, reviewed_at = now() WHERE id = $1`, notify.AdApproved
		case in.Action == "reject" && status == "pending":
			q, kind = `UPDATE ad_campaigns SET status = 'rejected', review_note = $2, reviewed_by = $3, reviewed_at = now() WHERE id = $1`, notify.AdRejected
		case in.Action == "pause" && status == "approved":
			q, kind = `UPDATE ad_campaigns SET status = 'paused', paused_by = 'admin', review_note = $2, reviewed_by = $3, reviewed_at = now() WHERE id = $1`, notify.AdPaused
		case in.Action == "resume" && status == "paused":
			q = `UPDATE ad_campaigns SET status = 'approved', paused_by = NULL, review_note = nullif($2, ''), reviewed_by = $3, reviewed_at = now() WHERE id = $1`
		case in.Action != "approve" && in.Action != "reject" && in.Action != "pause" && in.Action != "resume":
			return FieldError{"action", "approve, reject, pause or resume"}
		default:
			return ErrState
		}
		if _, err := tx.Exec(ctx, q, id, in.Note, adminID); err != nil {
			return err
		}
		if kind == "" {
			return nil
		}
		return notify.Enqueue(ctx, tx, storeID, kind, notify.Params{Product: product, Note: in.Note}, "#/store/"+storeID+"/ads")
	})
	if err != nil {
		return nil, err
	}
	return s.one(ctx, "c.id::text = $2", id)
}

// ---------------------------------------------------------------------------
// Tayang di website pembeli
// ---------------------------------------------------------------------------

// Slot: produk bersponsor yang dipilih untuk satu halaman.
type Slot struct{ CampaignID, ProductID string }

// eligible: kampanye yang boleh tayang sekarang (disetujui, dalam tanggal, anggaran hari ini & saldo seller cukup untuk satu klik,
// produk & toko tayang dan masih ada stok).
const eligible = `c.status = 'approved' AND c.starts_on <= $1::date AND (c.ends_on IS NULL OR c.ends_on >= $1::date)
	AND coalesce((SELECT d.spend FROM ad_stats_daily d WHERE d.campaign_id = c.id AND d.day = $1::date), 0) + c.bid <= c.daily_budget
	AND (SELECT coalesce(sum(l.amount), 0) FROM seller_ledger l WHERE l.store_id = c.store_id) >= c.bid
	AND p.deleted_at IS NULL AND p.status = 'active' AND p.moderation_status = 'approved' AND st.status = 'active' AND st.deleted_at IS NULL
	AND (SELECT coalesce(sum(v.stock_qty - v.reserved_qty), 0) FROM product_variants v WHERE v.product_id = p.id AND v.deleted_at IS NULL AND v.status = 'active') > 0`

// Pick memilih slot untuk halaman kategori (categoryID) atau hasil cari (query). Tawaran tertinggi menang;
// tawaran sama diacak supaya adil. Kosong kalau iklan dimatikan.
func (s *Service) Pick(ctx context.Context, categoryID, query string) ([]Slot, error) {
	st, err := s.Settings(ctx)
	if err != nil || !st.Enabled {
		return nil, err
	}
	match, args := "p.category_id IN (SELECT id FROM categories WHERE id = $2 OR parent_id = $2)", []any{today(), categoryID}
	if categoryID == "" {
		q := strings.TrimSpace(query)
		if q == "" {
			return nil, nil
		}
		like := "%" + strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(q) + "%"
		match, args = `(EXISTS (SELECT 1 FROM product_translations x WHERE x.product_id = p.id AND x.name ILIKE $2)
			OR EXISTS (SELECT 1 FROM categories cc JOIN category_translations ct ON ct.category_id IN (cc.id, cc.parent_id)
				WHERE cc.id = p.category_id AND ct.name ILIKE $2))`, []any{today(), like}
	}
	rows, err := s.Pool.Query(ctx, `SELECT DISTINCT ON (c.store_id) c.id, c.product_id, c.bid, random() AS r FROM ad_campaigns c
		JOIN products p ON p.id = c.product_id JOIN stores st ON st.id = c.store_id
		WHERE `+eligible+` AND `+match+` ORDER BY c.store_id, c.bid DESC, r`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	type cand struct {
		s   Slot
		bid decimal.Decimal
		r   float64
	}
	var all []cand
	for rows.Next() {
		var c cand
		if err := rows.Scan(&c.s.CampaignID, &c.s.ProductID, &c.bid, &c.r); err != nil {
			return nil, err
		}
		all = append(all, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	// satu slot per toko (DISTINCT ON di atas), lalu urutkan tawaran tertinggi
	for i := 1; i < len(all); i++ {
		for j := i; j > 0 && (all[j].bid.GreaterThan(all[j-1].bid) || (all[j].bid.Equal(all[j-1].bid) && all[j].r < all[j-1].r)); j-- {
			all[j], all[j-1] = all[j-1], all[j]
		}
	}
	out := []Slot{}
	for _, c := range all {
		if len(out) >= st.Slots {
			break
		}
		out = append(out, c.s)
	}
	return out, nil
}

// Impress mencatat tayangan untuk semua slot satu halaman dalam satu pernyataan.
func (s *Service) Impress(ctx context.Context, campaignIDs []string) error {
	if len(campaignIDs) == 0 {
		return nil
	}
	_, err := s.Pool.Exec(ctx, `INSERT INTO ad_stats_daily (campaign_id, day, impressions)
		SELECT id, $1::date, 1 FROM unnest($2::uuid[]) AS id
		ON CONFLICT (campaign_id, day) DO UPDATE SET impressions = ad_stats_daily.impressions + 1`, today(), campaignIDs)
	return err
}

// MaxChargesPerIP: klik bertagih maksimal dari satu IP per kampanye per hari (cegah klik curang tanpa cookie).
const MaxChargesPerIP = 3

func hash(s string) string {
	h := sha256.Sum256([]byte("phsar-ad|" + s))
	return hex.EncodeToString(h[:])
}

// Click mencatat klik dan mengembalikan slug produk (bahasa lang) untuk diarahkan.
// visitor = cookie pengunjung (kosong → IP+UA), ip = alamat klien. charge=false untuk bot.
func (s *Service) Click(ctx context.Context, campaignID, lang, visitor, ip, ua string, charge bool) (slug string, err error) {
	err = s.Pool.QueryRow(ctx, `SELECT coalesce(t.slug, ts.slug) FROM ad_campaigns c JOIN products p ON p.id = c.product_id
		JOIN product_translations ts ON ts.product_id = p.id AND ts.lang = p.source_lang
		LEFT JOIN product_translations t ON t.product_id = p.id AND t.lang = $2::lang_code WHERE c.id::text = $1`, campaignID, lang).Scan(&slug)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNotFound
	}
	if err != nil || !charge {
		return slug, err
	}
	if visitor == "" {
		visitor = "ipua|" + ip + "|" + ua
	}
	day := today()
	err = db.WithActor(ctx, s.Pool, "", "system", "", func(tx pgx.Tx) error {
		var storeID, product string
		var bid decimal.Decimal
		// kunci kampanye supaya klik bersamaan tidak melewati anggaran/saldo
		err := tx.QueryRow(ctx, `SELECT c.store_id, c.bid, coalesce(ts.name, '') FROM ad_campaigns c JOIN products p ON p.id = c.product_id JOIN stores st ON st.id = c.store_id
			LEFT JOIN product_translations ts ON ts.product_id = p.id AND ts.lang = p.source_lang
			WHERE c.id::text = $2 AND `+eligible+` FOR UPDATE OF c`, day, campaignID).Scan(&storeID, &bid, &product)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil // tidak tayang lagi (anggaran/saldo habis, dijeda): arahkan saja tanpa tagihan
		}
		if err != nil {
			return err
		}
		// saldo dikunci per toko (beberapa kampanye satu toko bisa diklik bersamaan)
		if _, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtext('ad-balance:' || $1))`, storeID); err != nil {
			return err
		}
		var bal decimal.Decimal
		if err := tx.QueryRow(ctx, `SELECT coalesce(sum(amount), 0) FROM seller_ledger WHERE store_id = $1`, storeID).Scan(&bal); err != nil {
			return err
		}
		if bal.LessThan(bid) {
			return nil
		}
		ipHash := hash("ip|" + ip)
		var n int
		if err := tx.QueryRow(ctx, `SELECT count(*) FROM ad_clicks WHERE campaign_id = $1 AND day = $2::date AND ip_hash = $3`, campaignID, day, ipHash).Scan(&n); err != nil {
			return err
		}
		if n >= MaxChargesPerIP {
			return nil
		}
		tag, err := tx.Exec(ctx, `INSERT INTO ad_clicks (campaign_id, day, visitor_hash, ip_hash, charged) VALUES ($1, $2::date, $3, $4, $5)
			ON CONFLICT (campaign_id, day, visitor_hash) DO NOTHING`, campaignID, day, hash("v|"+visitor), ipHash, bid)
		if err != nil || tag.RowsAffected() == 0 {
			return err // klik ulang hari ini: tidak ditagih lagi
		}
		if _, err := tx.Exec(ctx, `INSERT INTO seller_ledger (store_id, kind, amount, ad_campaign_id, note) VALUES ($1, 'ad_spend', $2, $3, $4)`,
			storeID, bid.Neg(), campaignID, truncate("Ad click · "+product, 300)); err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `INSERT INTO ad_stats_daily (campaign_id, day, clicks, spend) VALUES ($1, $2::date, 1, $3)
			ON CONFLICT (campaign_id, day) DO UPDATE SET clicks = ad_stats_daily.clicks + 1, spend = ad_stats_daily.spend + EXCLUDED.spend`, campaignID, day, bid)
		return err
	})
	return slug, err
}

func truncate(s string, n int) string {
	r := []rune(s)
	if len(r) > n {
		return string(r[:n])
	}
	return s
}

// IsBot: crawler & pratinjau tautan tidak ditagih (dan tidak dihitung tayangannya).
func IsBot(ua string) bool {
	u := strings.ToLower(ua)
	if u == "" {
		return true
	}
	for _, k := range []string{"bot", "crawl", "spider", "slurp", "preview", "facebookexternalhit", "whatsapp", "curl", "wget", "python", "go-http", "headless", "lighthouse"} {
		if strings.Contains(u, k) {
			return true
		}
	}
	return false
}
