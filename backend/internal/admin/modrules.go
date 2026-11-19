package admin

// Aturan moderasi otomatis (izin catalog.manage), tampil di halaman Categories:
//   POST/PUT/DELETE /admin/api/catalog/price-floors — harga minimum wajar per kata kunci nama produk
//   PUT /admin/api/catalog/moderation-rules        — sakelar tayang otomatis, masa percobaan, ambang diskon/turun harga, laporan
// Semua membalas katalog lengkap seperti endpoint katalog lain.

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/shopspring/decimal"

	"phsar/internal/modrules"
)

func (h *Handler) catalogSaveFloor(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Keyword  string `json:"keyword"`
		MinPrice string `json:"min_price"`
		Note     string `json:"note"`
	}
	if !decode(w, r, &in) {
		return
	}
	in.Keyword = strings.ToLower(strings.Join(strings.Fields(in.Keyword), " "))
	in.Note = strings.TrimSpace(in.Note)
	if c := utf8.RuneCountInString(in.Keyword); c < 2 || c > 80 {
		h.catalogReply(w, r, catalogFieldError{"keyword", "2–80 characters"})
		return
	}
	price, err := decimal.NewFromString(strings.TrimSpace(in.MinPrice))
	if err != nil || price.Sign() <= 0 || price.GreaterThan(decimal.NewFromInt(1_000_000)) {
		h.catalogReply(w, r, catalogFieldError{"min_price", "enter a price above 0"})
		return
	}
	price = price.Round(2)
	if utf8.RuneCountInString(in.Note) > 200 {
		h.catalogReply(w, r, catalogFieldError{"note", "max 200 characters"})
		return
	}
	fid := chi.URLParam(r, "fid")
	h.catalogReply(w, r, h.catalogTx(r, func(ctx context.Context, tx pgx.Tx) error {
		var tag pgconn.CommandTag
		var err error
		if fid == "" {
			tag, err = tx.Exec(ctx, `INSERT INTO price_floor_rules (keyword, min_price, note) VALUES ($1, $2::numeric, nullif($3, ''))`,
				in.Keyword, price.StringFixed(2), in.Note)
		} else {
			id, convErr := strconv.Atoi(fid)
			if convErr != nil {
				return pgx.ErrNoRows
			}
			tag, err = tx.Exec(ctx, `UPDATE price_floor_rules SET keyword = $2, min_price = $3::numeric, note = nullif($4, '') WHERE id = $1`,
				id, in.Keyword, price.StringFixed(2), in.Note)
		}
		var pe *pgconn.PgError
		if errors.As(err, &pe) && pe.Code == "23505" {
			return catalogFieldError{"keyword", "already in the list"}
		}
		if err == nil && tag.RowsAffected() == 0 {
			return pgx.ErrNoRows
		}
		return err
	}))
}

func (h *Handler) catalogDeleteFloor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "fid"))
	if err != nil {
		writeErr(w, http.StatusNotFound, "not_found", "not found")
		return
	}
	h.catalogReply(w, r, h.catalogTx(r, func(ctx context.Context, tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, `DELETE FROM price_floor_rules WHERE id = $1`, id)
		if err == nil && tag.RowsAffected() == 0 {
			return pgx.ErrNoRows
		}
		return err
	}))
}

func (h *Handler) catalogSaveModRules(w http.ResponseWriter, r *http.Request) {
	in := modrules.Defaults()
	if !decode(w, r, &in) {
		return
	}
	if field, ok := in.Validate(); !ok {
		h.catalogReply(w, r, catalogFieldError{field, "value out of range"})
		return
	}
	h.catalogReply(w, r, h.catalogTx(r, func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `INSERT INTO platform_settings (key, value, updated_by, updated_at) VALUES ($1, $2, $3, now())
			ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_by = EXCLUDED.updated_by, updated_at = now()`,
			modrules.SettingKey, in, actor(r))
		return err
	}))
}
