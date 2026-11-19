package admin

// Promos: voucher & flash sale yang dibiayai Phsar (izin commission.manage).
// Voucher dipakai pembeli di keranjang/checkout; potongan ditanggung Phsar, penghasilan seller tidak berubah.
// Flash sale: harga khusus per varian dalam rentang waktu, boleh dengan kuota barang.

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"

	"phsar/internal/platform/db"
	"phsar/internal/pricehistory" // pricehistory
)

func (h *Handler) mountPromos(r chi.Router) {
	r.Route("/promos", func(p chi.Router) {
		p.Use(h.can("commission.manage"))
		p.Get("/", h.promosView)
		p.Get("/variants", h.promoVariants)
		p.Post("/vouchers", h.voucherSave)
		p.Put("/vouchers/{id}", h.voucherSave)
		p.Post("/vouchers/{id}/active", h.voucherActive)
		h.mountFlash(p) // sesi flash sale (admin/flash.go)
	})
}

type promoFieldErr struct{ field, msg string }

func (e promoFieldErr) Error() string { return e.field + ": " + e.msg }

var errPromoNotFound = errors.New("not found")

func (h *Handler) promoReply(w http.ResponseWriter, r *http.Request, err error) {
	var fe promoFieldErr
	switch {
	case err == nil:
		h.promosView(w, r)
	case errors.As(err, &fe):
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"error": map[string]string{"code": "validation", "field": fe.field, "message": fe.msg}})
	case errors.Is(err, errPromoNotFound), errors.Is(err, pgx.ErrNoRows):
		writeErr(w, http.StatusNotFound, "not_found", "not found")
	default:
		h.log.Error("promos", "err", err)
		writeErr(w, http.StatusInternalServerError, "internal", "could not save")
	}
}

// GET /promos — voucher (dengan pemakaian). Flash sale: GET /promos/flash.
func (h *Handler) promosView(w http.ResponseWriter, r *http.Request) {
	var raw []byte
	err := h.pool.QueryRow(r.Context(), `SELECT json_build_object(
		'vouchers', (SELECT coalesce(json_agg(x ORDER BY x.created_at DESC), '[]') FROM (
			SELECT v.id, v.code, v.name, v.kind, v.value::text AS value, v.max_discount::text AS max_discount, v.min_spend::text AS min_spend,
			  v.starts_at, v.ends_at, v.total_limit, v.per_user_limit, v.first_order_only, v.active, v.created_at,
			  CASE WHEN NOT v.active THEN 'off' WHEN v.starts_at > now() THEN 'scheduled' WHEN v.ends_at IS NOT NULL AND v.ends_at <= now() THEN 'ended' ELSE 'live' END AS state,
			  (SELECT count(*) FROM voucher_redemptions r WHERE r.voucher_id = v.id AND r.status = 'used') AS used,
			  (SELECT count(*) FROM voucher_redemptions r WHERE r.voucher_id = v.id AND r.status = 'held') AS held,
			  (SELECT coalesce(sum(r.discount), 0)::text FROM voucher_redemptions r WHERE r.voucher_id = v.id AND r.status = 'used') AS cost
			FROM vouchers v) x),
		'totals', json_build_object(
			'cost_30d', (SELECT coalesce(sum(discount), 0)::text FROM voucher_redemptions WHERE status = 'used' AND created_at > now() - interval '30 days'),
			'used_30d', (SELECT count(*) FROM voucher_redemptions WHERE status = 'used' AND created_at > now() - interval '30 days'))
		)`).Scan(&raw)
	if err != nil {
		h.log.Error("promos", "err", err)
		writeErr(w, http.StatusInternalServerError, "internal", "could not load")
		return
	}
	writeRawJSON(w, raw)
}

// GET /promos/variants?q=&at= — pilih varian untuk flash sale (nama produk, SKU, atau toko).
// pricehistory: low30 = harga normal terendah 30 hari sebelum `at` (jam mulai sesi; kosong = sekarang).
func (h *Handler) promoVariants(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	at, perr := time.Parse(time.RFC3339, r.URL.Query().Get("at")) // pricehistory
	if perr != nil {
		at = time.Now()
	}
	rules, rerr := pricehistory.LoadRules(r.Context(), h.pool) // pricehistory
	if rerr != nil {
		h.log.Error("promos variants", "err", rerr)
		writeErr(w, http.StatusInternalServerError, "internal", "search failed")
		return
	}
	if len([]rune(q)) < 2 {
		writeRawJSON(w, []byte(`{"variants":[]}`))
		return
	}
	var raw []byte
	err := h.pool.QueryRow(r.Context(), `SELECT coalesce(json_agg(x), '[]') FROM (
		SELECT v.id, v.sku, v.price::text AS price, greatest(v.stock_qty - v.reserved_qty, 0) AS stock, coalesce(t.name, '') AS product, st.name AS store,
		  variant_low_price(v.id, $2)::text AS low30,
		  coalesce((SELECT string_agg(coalesce(ov.value_i18n->>'en', ov.value_i18n->>p.source_lang::text), ' / ' ORDER BY o.position)
		    FROM variant_option_values vov JOIN product_option_values ov ON ov.id = vov.option_value_id JOIN product_options o ON o.id = vov.option_id
		    WHERE vov.variant_id = v.id), '') AS options
		FROM product_variants v JOIN products p ON p.id = v.product_id JOIN stores st ON st.id = p.store_id
		LEFT JOIN product_translations t ON t.product_id = p.id AND t.lang = p.source_lang
		WHERE v.deleted_at IS NULL AND v.status = 'active' AND p.deleted_at IS NULL AND p.status = 'active' AND p.moderation_status = 'approved'
		  AND st.status = 'active' AND (t.name ILIKE '%' || $1 || '%' OR v.sku ILIKE $1 || '%' OR st.name ILIKE '%' || $1 || '%')
		ORDER BY t.name, v.sku LIMIT 80) x`, q, at).Scan(&raw)
	if err != nil {
		h.log.Error("promos variants", "err", err)
		writeErr(w, http.StatusInternalServerError, "internal", "search failed")
		return
	}
	writeRawJSON(w, []byte(`{"variants":`+string(raw)+`,"low30_rule":`+strconv.FormatBool(rules.Low30Enabled)+`}`))
}

var voucherCodeRe = regexp.MustCompile(`^[A-Z0-9_-]{3,24}$`)

type voucherIn struct {
	Code           string     `json:"code"`
	Name           string     `json:"name"`
	Kind           string     `json:"kind"`
	Value          string     `json:"value"`
	MaxDiscount    string     `json:"max_discount"`
	MinSpend       string     `json:"min_spend"`
	StartsAt       *time.Time `json:"starts_at"`
	EndsAt         *time.Time `json:"ends_at"`
	TotalLimit     *int       `json:"total_limit"`
	PerUserLimit   int        `json:"per_user_limit"`
	FirstOrderOnly bool       `json:"first_order_only"`
	Active         *bool      `json:"active"`
}

func money2(field, s string, required bool, max float64) (*decimal.Decimal, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		if required {
			return nil, promoFieldErr{field, "required"}
		}
		return nil, nil
	}
	d, err := decimal.NewFromString(s)
	if err != nil || d.Exponent() < -2 || d.IsNegative() || d.GreaterThan(decimal.NewFromFloat(max)) {
		return nil, promoFieldErr{field, "enter an amount like 5 or 2.50"}
	}
	return &d, nil
}

// POST /promos/vouchers · PUT /promos/vouchers/{id}
func (h *Handler) voucherSave(w http.ResponseWriter, r *http.Request) {
	var in voucherIn
	if !decode(w, r, &in) {
		return
	}
	id := chi.URLParam(r, "id")
	h.promoReply(w, r, h.saveVoucher(r, id, in))
}

func (h *Handler) saveVoucher(r *http.Request, id string, in voucherIn) error {
	in.Code = strings.ToUpper(strings.TrimSpace(in.Code))
	in.Name = strings.TrimSpace(in.Name)
	if !voucherCodeRe.MatchString(in.Code) {
		return promoFieldErr{"code", "3–24 letters, numbers, - or _"}
	}
	if len([]rune(in.Name)) > 80 {
		return promoFieldErr{"name", "too long"}
	}
	if in.Kind != "percent" && in.Kind != "fixed" {
		return promoFieldErr{"kind", "percent or fixed"}
	}
	maxValue := 1000.0
	if in.Kind == "percent" {
		maxValue = 90
	}
	val, err := money2("value", in.Value, true, maxValue)
	if err != nil {
		return err
	}
	if !val.IsPositive() {
		return promoFieldErr{"value", "must be more than 0"}
	}
	maxD, err := money2("max_discount", in.MaxDiscount, false, 1000)
	if err != nil {
		return err
	}
	if maxD != nil && !maxD.IsPositive() {
		maxD = nil
	}
	minS, err := money2("min_spend", in.MinSpend, false, 100000)
	if err != nil {
		return err
	}
	if minS == nil {
		z := decimal.Zero
		minS = &z
	}
	starts := time.Now()
	if in.StartsAt != nil {
		starts = *in.StartsAt
	}
	if in.EndsAt != nil && !in.EndsAt.After(starts) {
		return promoFieldErr{"ends_at", "must be after the start"}
	}
	if in.EndsAt != nil && id == "" && !in.EndsAt.After(time.Now()) {
		return promoFieldErr{"ends_at", "must be in the future"}
	}
	if in.TotalLimit != nil && (*in.TotalLimit < 1 || *in.TotalLimit > 1000000) {
		return promoFieldErr{"total_limit", "1 or more, or empty for no limit"}
	}
	if in.PerUserLimit == 0 {
		in.PerUserLimit = 1
	}
	if in.PerUserLimit < 1 || in.PerUserLimit > 100 {
		return promoFieldErr{"per_user_limit", "1–100"}
	}
	active := true
	if in.Active != nil {
		active = *in.Active
	}
	return db.WithActor(r.Context(), h.pool, actor(r), "admin", "", func(tx pgx.Tx) error {
		var taken bool
		if err := tx.QueryRow(r.Context(), `SELECT EXISTS (SELECT 1 FROM vouchers WHERE upper(code) = $1 AND id::text <> $2)`, in.Code, id).Scan(&taken); err != nil {
			return err
		}
		if taken {
			return promoFieldErr{"code", "this code is already used"}
		}
		if id == "" {
			_, err := tx.Exec(r.Context(), `INSERT INTO vouchers (code, name, kind, value, max_discount, min_spend, starts_at, ends_at, total_limit,
				per_user_limit, first_order_only, active, created_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
				in.Code, in.Name, in.Kind, val, maxD, minS, starts, in.EndsAt, in.TotalLimit, in.PerUserLimit, in.FirstOrderOnly, active, actor(r))
			return err
		}
		if !uuidRe.MatchString(id) {
			return errPromoNotFound
		}
		tag, err := tx.Exec(r.Context(), `UPDATE vouchers SET code = $2, name = $3, kind = $4, value = $5, max_discount = $6, min_spend = $7,
			starts_at = $8, ends_at = $9, total_limit = $10, per_user_limit = $11, first_order_only = $12, active = $13 WHERE id = $1`,
			id, in.Code, in.Name, in.Kind, val, maxD, minS, starts, in.EndsAt, in.TotalLimit, in.PerUserLimit, in.FirstOrderOnly, active)
		if err == nil && tag.RowsAffected() == 0 {
			return errPromoNotFound
		}
		return err
	})
}

// POST /promos/vouchers/{id}/active {active}
func (h *Handler) voucherActive(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Active bool `json:"active"`
	}
	if !decode(w, r, &in) {
		return
	}
	id := chi.URLParam(r, "id")
	if !uuidRe.MatchString(id) {
		h.promoReply(w, r, errPromoNotFound)
		return
	}
	err := db.WithActor(r.Context(), h.pool, actor(r), "admin", "", func(tx pgx.Tx) error {
		tag, err := tx.Exec(r.Context(), `UPDATE vouchers SET active = $2 WHERE id = $1`, id, in.Active)
		if err == nil && tag.RowsAffected() == 0 {
			return errPromoNotFound
		}
		return err
	})
	h.promoReply(w, r, err)
}
