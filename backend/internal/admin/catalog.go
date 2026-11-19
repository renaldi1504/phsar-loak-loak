package admin

// Categories & brands (izin catalog.manage): kategori (urutan, ikon, barang bekas boleh/tidak, nama & slug 3 bahasa,
// sembunyikan / hapus), atribut kategori, merek (logo, terverifikasi, aktif, kategori) dan kata terlarang untuk moderasi.
// Website, app dan Seller Center membaca tabel ini langsung (tanpa cache di memori), jadi perubahan langsung berlaku
// (website: paling lambat 60 detik karena Cache-Control).

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"phsar/internal/modrules"
	"phsar/internal/platform/db"
)

var (
	catSlugRe  = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
	attrCodeRe = regexp.MustCompile(`^[a-z][a-z0-9_]{0,39}$`)
	catLangs   = []string{"en", "km", "zh"}
	catIcons   = []string{"headphones", "beauty", "shirt", "star", "cooker", "groceries", "sofa", "bag", "box", "heart", "home", "video",
		"tshirt", "watch", "phone", "laptop", "health", "baby", "ball", "book"}
	kwActions   = map[string]bool{"block": true, "review": true}
	errHasItems = errors.New("masih dipakai")
)

func (h *Handler) mountCatalog(r chi.Router) {
	r.Route("/catalog", func(p chi.Router) {
		p.Use(h.can("catalog.manage"))
		p.Get("/", h.catalogView)
		p.Post("/categories", h.catalogCreateCategory)
		p.Put("/categories/{id}", h.catalogUpdateCategory)
		p.Post("/categories/{id}/move", h.catalogMoveCategory)
		p.Post("/categories/{id}/hide", h.catalogHideCategory)
		p.Delete("/categories/{id}", h.catalogDeleteCategory)
		p.Post("/categories/{id}/attributes", h.catalogSaveAttribute)
		p.Put("/attributes/{aid}", h.catalogSaveAttribute)
		p.Delete("/attributes/{aid}", h.catalogDeleteAttribute)
		p.Post("/brands", h.catalogSaveBrand)
		p.Put("/brands/{bid}", h.catalogSaveBrand)
		p.Post("/brands/{bid}/logo", h.catalogBrandLogo)
		p.Delete("/brands/{bid}", h.catalogDeleteBrand)
		p.Post("/keywords", h.catalogAddKeyword)
		p.Delete("/keywords/{kid}", h.catalogDeleteKeyword)
		p.Post("/price-floors", h.catalogSaveFloor)
		p.Put("/price-floors/{fid}", h.catalogSaveFloor)
		p.Delete("/price-floors/{fid}", h.catalogDeleteFloor)
		p.Put("/moderation-rules", h.catalogSaveModRules)
	})
}

// GET /catalog — semuanya dalam satu jawaban (data kecil: puluhan kategori/merek).
func (h *Handler) catalogView(w http.ResponseWriter, r *http.Request) {
	var raw []byte
	err := h.pool.QueryRow(r.Context(), `SELECT json_build_object(
		'categories', (SELECT coalesce(json_agg(json_build_object(
			'id', c.id, 'parent_id', c.parent_id, 'icon', coalesce(c.icon, ''), 'allows_used', c.allows_used, 'restricted', c.restricted, 'sort_order', c.sort_order,
			'hidden', c.deleted_at IS NOT NULL,
			'names', (SELECT coalesce(json_object_agg(t.lang, t.name), '{}') FROM category_translations t WHERE t.category_id = c.id),
			'slugs', (SELECT coalesce(json_object_agg(t.lang, t.slug), '{}') FROM category_translations t WHERE t.category_id = c.id),
			'meta_titles', (SELECT coalesce(json_object_agg(t.lang, coalesce(t.meta_title, '')), '{}') FROM category_translations t WHERE t.category_id = c.id),
			'meta_descriptions', (SELECT coalesce(json_object_agg(t.lang, coalesce(t.meta_description, '')), '{}') FROM category_translations t WHERE t.category_id = c.id),
			'products', (SELECT count(*) FROM products p WHERE p.category_id = c.id AND p.deleted_at IS NULL),
			'live_products', (SELECT count(*) FROM products p WHERE p.category_id = c.id AND p.deleted_at IS NULL AND p.status = 'active'),
			'attributes', (SELECT coalesce(json_agg(json_build_object('id', a.id, 'code', a.code, 'label', a.label_i18n, 'data_type', a.data_type,
				'choices', coalesce(a.choices, '[]'), 'is_required', a.is_required, 'is_filterable', a.is_filterable,
				'values', (SELECT count(*) FROM product_attribute_values v WHERE v.attribute_id = a.id)) ORDER BY a.sort_order, a.id), '[]')
				FROM category_attributes a WHERE a.category_id = c.id)
			) ORDER BY c.parent_id NULLS FIRST, c.sort_order, c.id), '[]') FROM categories c),
		'brands', (SELECT coalesce(json_agg(json_build_object('id', b.id, 'slug', b.slug, 'name', b.name, 'logo_url', b.logo_url,
			'verified', b.verified, 'active', b.active,
			'category_ids', (SELECT coalesce(json_agg(bc.category_id ORDER BY bc.category_id), '[]') FROM brand_categories bc WHERE bc.brand_id = b.id),
			'products', (SELECT count(*) FROM products p WHERE p.brand_id = b.id AND p.deleted_at IS NULL)) ORDER BY lower(b.name)), '[]') FROM brands b),
		'keywords', (SELECT coalesce(json_agg(json_build_object('id', k.id, 'keyword', k.keyword, 'lang', k.lang, 'category', k.category,
			'action', k.action) ORDER BY k.category, k.lang, k.keyword), '[]') FROM prohibited_keywords k),
		'price_floors', (SELECT coalesce(json_agg(json_build_object('id', f.id, 'keyword', f.keyword, 'min_price', f.min_price::text,
			'note', coalesce(f.note, '')) ORDER BY lower(f.keyword)), '[]') FROM price_floor_rules f)
		)`).Scan(&raw)
	if err != nil {
		h.log.Error("katalog", "err", err)
		writeErr(w, http.StatusInternalServerError, "internal", "could not load catalog")
		return
	}
	rules, err := modrules.Load(r.Context(), h.pool)
	if err != nil {
		h.log.Error("katalog", "err", err)
		writeErr(w, http.StatusInternalServerError, "internal", "could not load catalog")
		return
	}
	icons, _ := json.Marshal(catIcons)
	rulesJSON, _ := json.Marshal(rules)
	out := append(raw[:len(raw)-1], []byte(`,"icons":`)...)
	out = append(append(out, icons...), []byte(`,"moderation_rules":`)...)
	out = append(append(out, rulesJSON...), '}')
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(out)
}

// catalogReply: setelah perubahan, kirim ulang seluruh katalog; error validasi → 422 dengan field.
func (h *Handler) catalogReply(w http.ResponseWriter, r *http.Request, err error) {
	var fe catalogFieldError
	var pe *pgconn.PgError
	switch {
	case err == nil:
		h.catalogView(w, r)
	case errors.As(err, &fe):
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"error": map[string]string{"code": "validation", "field": fe.Field, "message": fe.Msg}})
	case errors.Is(err, pgx.ErrNoRows):
		writeErr(w, http.StatusNotFound, "not_found", "not found")
	case errors.Is(err, errHasItems):
		writeErr(w, http.StatusConflict, "in_use", "still used by products — hide or deactivate it instead")
	case errors.As(err, &pe) && pe.Code == "23503":
		writeErr(w, http.StatusConflict, "in_use", "still used by products — hide or deactivate it instead")
	default:
		h.log.Error("katalog: simpan", "err", err)
		writeErr(w, http.StatusInternalServerError, "internal", "could not save")
	}
}

// catalogFieldError: input tidak valid untuk field tertentu (422).
type catalogFieldError struct{ Field, Msg string }

func (e catalogFieldError) Error() string { return e.Field + ": " + e.Msg }

func (h *Handler) catalogTx(r *http.Request, fn func(ctx context.Context, tx pgx.Tx) error) error {
	ctx := r.Context()
	return db.WithActor(ctx, h.pool, actor(r), "admin", "", func(tx pgx.Tx) error { return fn(ctx, tx) })
}

// ---------------------------------------------------------------------------
// Kategori
// ---------------------------------------------------------------------------

type categoryInput struct {
	ID               string            `json:"id"`
	ParentID         string            `json:"parent_id"`
	Icon             string            `json:"icon"`
	AllowsUsed       bool              `json:"allows_used"`
	Restricted       bool              `json:"restricted"` // produk di kategori ini selalu ditinjau manual
	Names            map[string]string `json:"names"`
	Slugs            map[string]string `json:"slugs"`
	MetaTitles       map[string]string `json:"meta_titles"`
	MetaDescriptions map[string]string `json:"meta_descriptions"`
}

func validIcon(i string) bool {
	for _, x := range catIcons {
		if x == i {
			return true
		}
	}
	return false
}

func (in *categoryInput) clean() error {
	if in.Icon != "" && !validIcon(in.Icon) {
		return catalogFieldError{"icon", "unknown icon"}
	}
	if in.Slugs == nil {
		in.Slugs = map[string]string{}
	}
	if in.Names == nil {
		in.Names = map[string]string{}
	}
	for _, l := range catLangs {
		n := strings.TrimSpace(in.Names[l])
		if c := utf8.RuneCountInString(n); c < 1 || c > 80 {
			return catalogFieldError{"name_" + l, "1–80 characters"}
		}
		in.Names[l] = n
		s := strings.ToLower(strings.TrimSpace(in.Slugs[l]))
		if s == "" { // slug bahasa lain mengikuti bahasa Inggris (URL /km/c/fashion tetap rapi)
			s = strings.ToLower(strings.TrimSpace(in.Slugs["en"]))
		}
		if s == "" {
			s = in.ID
		}
		if !catSlugRe.MatchString(s) || len(s) > 80 {
			return catalogFieldError{"slug_" + l, "lowercase letters, numbers and dashes"}
		}
		in.Slugs[l] = s
		if utf8.RuneCountInString(in.MetaTitles[l]) > 70 {
			return catalogFieldError{"meta_title_" + l, "max 70 characters"}
		}
		if utf8.RuneCountInString(in.MetaDescriptions[l]) > 160 {
			return catalogFieldError{"meta_description_" + l, "max 160 characters"}
		}
	}
	return nil
}

// saveTranslations: nama, slug & path (induk/anak) per bahasa. Slug bentrok → 422 slug_<lang>.
func saveTranslations(ctx context.Context, tx pgx.Tx, id, parent string, in categoryInput) error {
	for _, l := range catLangs {
		path := in.Slugs[l]
		if parent != "" {
			var pp string
			if err := tx.QueryRow(ctx, `SELECT path FROM category_translations WHERE category_id = $1 AND lang = $2::lang_code`, parent, l).Scan(&pp); err == nil {
				path = pp + "/" + path
			}
		}
		_, err := tx.Exec(ctx, `INSERT INTO category_translations (category_id, lang, name, slug, path, meta_title, meta_description)
			VALUES ($1, $2::lang_code, $3, $4, $5, nullif($6, ''), nullif($7, ''))
			ON CONFLICT (category_id, lang) DO UPDATE SET name = EXCLUDED.name, slug = EXCLUDED.slug, path = EXCLUDED.path,
			  meta_title = EXCLUDED.meta_title, meta_description = EXCLUDED.meta_description`,
			id, l, in.Names[l], in.Slugs[l], path, strings.TrimSpace(in.MetaTitles[l]), strings.TrimSpace(in.MetaDescriptions[l]))
		var pe *pgconn.PgError
		if errors.As(err, &pe) && pe.Code == "23505" {
			return catalogFieldError{"slug_" + l, "this web address is already used by another category"}
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *Handler) catalogCreateCategory(w http.ResponseWriter, r *http.Request) {
	var in categoryInput
	if !decode(w, r, &in) {
		return
	}
	in.ID = strings.ToLower(strings.TrimSpace(in.ID))
	if !catSlugRe.MatchString(in.ID) || len(in.ID) > 60 {
		writeFieldErr(w, "id")
		return
	}
	if err := in.clean(); err != nil {
		h.catalogReply(w, r, err)
		return
	}
	h.catalogReply(w, r, h.catalogTx(r, func(ctx context.Context, tx pgx.Tx) error {
		var exists bool
		if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM categories WHERE id = $1)`, in.ID).Scan(&exists); err != nil {
			return err
		}
		if exists {
			return catalogFieldError{"id", "this ID is already used"}
		}
		if in.ParentID != "" {
			var top bool
			if err := tx.QueryRow(ctx, `SELECT parent_id IS NULL FROM categories WHERE id = $1`, in.ParentID).Scan(&top); err != nil || !top {
				return catalogFieldError{"parent_id", "choose a main category (one level only)"}
			}
		}
		if _, err := tx.Exec(ctx, `INSERT INTO categories (id, parent_id, icon, allows_used, restricted, sort_order)
			VALUES ($1, nullif($2, ''), nullif($3, ''), $4, $5, (SELECT coalesce(max(sort_order), 0) + 10 FROM categories WHERE parent_id IS NOT DISTINCT FROM nullif($2, '')))`,
			in.ID, in.ParentID, in.Icon, in.AllowsUsed, in.Restricted); err != nil {
			return err
		}
		return saveTranslations(ctx, tx, in.ID, in.ParentID, in)
	}))
}

func (h *Handler) catalogUpdateCategory(w http.ResponseWriter, r *http.Request) {
	var in categoryInput
	if !decode(w, r, &in) {
		return
	}
	in.ID = chi.URLParam(r, "id")
	if err := in.clean(); err != nil {
		h.catalogReply(w, r, err)
		return
	}
	h.catalogReply(w, r, h.catalogTx(r, func(ctx context.Context, tx pgx.Tx) error {
		var parent *string
		if err := tx.QueryRow(ctx, `UPDATE categories SET icon = nullif($2, ''), allows_used = $3, restricted = $4 WHERE id = $1 RETURNING parent_id`,
			in.ID, in.Icon, in.AllowsUsed, in.Restricted).Scan(&parent); err != nil {
			return err
		}
		p := ""
		if parent != nil {
			p = *parent
		}
		if err := saveTranslations(ctx, tx, in.ID, p, in); err != nil {
			return err
		}
		// path anak ikut berubah kalau slug induk berubah
		for _, l := range catLangs {
			if _, err := tx.Exec(ctx, `UPDATE category_translations t SET path = $2 || '/' || t.slug
				FROM categories c WHERE c.id = t.category_id AND c.parent_id = $1 AND t.lang = $3::lang_code`, in.ID, in.Slugs[l], l); err != nil {
				return err
			}
		}
		return nil
	}))
}

// POST /categories/{id}/move {dir: up|down} — tukar urutan dengan kategori sebelah (induk yang sama).
func (h *Handler) catalogMoveCategory(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Dir string `json:"dir"`
	}
	if !decode(w, r, &in) {
		return
	}
	if in.Dir != "up" && in.Dir != "down" {
		writeFieldErr(w, "dir")
		return
	}
	id := chi.URLParam(r, "id")
	h.catalogReply(w, r, h.catalogTx(r, func(ctx context.Context, tx pgx.Tx) error {
		// rapikan dulu: urutan 10, 20, 30 … supaya tidak ada nilai kembar
		if _, err := tx.Exec(ctx, `UPDATE categories c SET sort_order = x.n * 10 FROM (
			SELECT id, row_number() OVER (PARTITION BY parent_id ORDER BY sort_order, id) AS n FROM categories) x
			WHERE x.id = c.id AND c.sort_order <> x.n * 10`); err != nil {
			return err
		}
		var parent *string
		var order int
		if err := tx.QueryRow(ctx, `SELECT parent_id, sort_order FROM categories WHERE id = $1`, id).Scan(&parent, &order); err != nil {
			return err
		}
		cmp, dir := "<", "DESC"
		if in.Dir == "down" {
			cmp, dir = ">", "ASC"
		}
		var other string
		var otherOrder int
		err := tx.QueryRow(ctx, `SELECT id, sort_order FROM categories WHERE parent_id IS NOT DISTINCT FROM $1 AND sort_order `+cmp+` $2
			ORDER BY sort_order `+dir+` LIMIT 1`, parent, order).Scan(&other, &otherOrder)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil // sudah paling atas / bawah
		}
		if err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `UPDATE categories SET sort_order = $2 WHERE id = $1`, id, otherOrder); err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `UPDATE categories SET sort_order = $2 WHERE id = $1`, other, order)
		return err
	}))
}

// POST /categories/{id}/hide {hidden} — sembunyikan dari website/app/Seller Center (produk lama tetap ada).
func (h *Handler) catalogHideCategory(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Hidden bool `json:"hidden"`
	}
	if !decode(w, r, &in) {
		return
	}
	h.catalogReply(w, r, h.catalogTx(r, func(ctx context.Context, tx pgx.Tx) error {
		q := `UPDATE categories SET deleted_at = NULL WHERE id = $1 OR parent_id = $1`
		if in.Hidden {
			q = `UPDATE categories SET deleted_at = coalesce(deleted_at, now()) WHERE id = $1 OR parent_id = $1`
		}
		tag, err := tx.Exec(ctx, q, chi.URLParam(r, "id"))
		if err == nil && tag.RowsAffected() == 0 {
			return pgx.ErrNoRows
		}
		return err
	}))
}

// DELETE /categories/{id} — hanya kategori yang belum pernah punya produk & sub-kategori.
func (h *Handler) catalogDeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.catalogReply(w, r, h.catalogTx(r, func(ctx context.Context, tx pgx.Tx) error {
		var exists, used bool
		if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM categories WHERE id = $1),
			EXISTS (SELECT 1 FROM products WHERE category_id = $1) OR EXISTS (SELECT 1 FROM categories WHERE parent_id = $1)
			OR EXISTS (SELECT 1 FROM commission_rules WHERE category_id = $1)`, id).Scan(&exists, &used); err != nil {
			return err
		}
		if !exists {
			return pgx.ErrNoRows
		}
		if used {
			return errHasItems
		}
		for _, q := range []string{`DELETE FROM category_attributes WHERE category_id = $1`, `DELETE FROM brand_categories WHERE category_id = $1`,
			`DELETE FROM category_translations WHERE category_id = $1`, `DELETE FROM categories WHERE id = $1`} {
			if _, err := tx.Exec(ctx, q, id); err != nil {
				return err
			}
		}
		return nil
	}))
}

// ---------------------------------------------------------------------------
// Atribut kategori (spesifikasi yang diisi seller per kategori)
// ---------------------------------------------------------------------------

func (h *Handler) catalogSaveAttribute(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Code         string            `json:"code"`
		Label        map[string]string `json:"label"`
		DataType     string            `json:"data_type"`
		Choices      []string          `json:"choices"`
		IsRequired   bool              `json:"is_required"`
		IsFilterable bool              `json:"is_filterable"`
	}
	if !decode(w, r, &in) {
		return
	}
	in.Code = strings.ToLower(strings.TrimSpace(in.Code))
	if !attrCodeRe.MatchString(in.Code) {
		writeFieldErr(w, "code")
		return
	}
	if in.Label == nil {
		in.Label = map[string]string{}
	}
	for _, l := range catLangs {
		v := strings.TrimSpace(in.Label[l])
		if c := utf8.RuneCountInString(v); c < 1 || c > 60 {
			writeFieldErr(w, "label_"+l)
			return
		}
		in.Label[l] = v
	}
	switch in.DataType {
	case "text", "number", "bool":
		in.Choices = nil
	case "select":
		clean := []string{}
		seen := map[string]bool{}
		for _, c := range in.Choices {
			c = strings.TrimSpace(c)
			if c != "" && !seen[c] && utf8.RuneCountInString(c) <= 40 {
				seen[c] = true
				clean = append(clean, c)
			}
		}
		if len(clean) < 2 || len(clean) > 50 {
			writeFieldErr(w, "choices")
			return
		}
		in.Choices = clean
	default:
		writeFieldErr(w, "data_type")
		return
	}
	label, _ := json.Marshal(in.Label)
	var choices any
	if in.Choices != nil {
		b, _ := json.Marshal(in.Choices)
		choices = b
	}
	aid, catID := chi.URLParam(r, "aid"), chi.URLParam(r, "id")
	h.catalogReply(w, r, h.catalogTx(r, func(ctx context.Context, tx pgx.Tx) error {
		var err error
		if aid == "" {
			_, err = tx.Exec(ctx, `INSERT INTO category_attributes (category_id, code, label_i18n, data_type, choices, is_required, is_filterable, sort_order)
				SELECT c.id, $2, $3, $4, $5, $6, $7, (SELECT coalesce(max(sort_order), 0) + 10 FROM category_attributes WHERE category_id = c.id)
				FROM categories c WHERE c.id = $1`, catID, in.Code, label, in.DataType, choices, in.IsRequired, in.IsFilterable)
		} else {
			var tag pgconn.CommandTag
			// tipe data tidak boleh diganti kalau sudah ada nilai dari produk
			tag, err = tx.Exec(ctx, `UPDATE category_attributes a SET code = $2, label_i18n = $3, data_type = $4, choices = $5, is_required = $6, is_filterable = $7
				WHERE a.id = $1 AND (a.data_type = $4 OR NOT EXISTS (SELECT 1 FROM product_attribute_values v WHERE v.attribute_id = a.id))`,
				aid, in.Code, label, in.DataType, choices, in.IsRequired, in.IsFilterable)
			if err == nil && tag.RowsAffected() == 0 {
				var exists bool
				_ = tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM category_attributes WHERE id = $1)`, aid).Scan(&exists)
				if !exists {
					return pgx.ErrNoRows
				}
				return catalogFieldError{"data_type", "products already use this attribute — the type cannot change"}
			}
		}
		var pe *pgconn.PgError
		if errors.As(err, &pe) && pe.Code == "23505" {
			return catalogFieldError{"code", "this code already exists in the category"}
		}
		if errors.As(err, &pe) && pe.Code == "22P02" {
			return pgx.ErrNoRows
		}
		return err
	}))
}

func (h *Handler) catalogDeleteAttribute(w http.ResponseWriter, r *http.Request) {
	aid, err := strconv.Atoi(chi.URLParam(r, "aid"))
	if err != nil {
		writeErr(w, http.StatusNotFound, "not_found", "not found")
		return
	}
	h.catalogReply(w, r, h.catalogTx(r, func(ctx context.Context, tx pgx.Tx) error {
		var used bool
		if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM product_attribute_values WHERE attribute_id = $1)`, aid).Scan(&used); err != nil {
			return err
		}
		if used {
			return errHasItems
		}
		tag, err := tx.Exec(ctx, `DELETE FROM category_attributes WHERE id = $1`, aid)
		if err == nil && tag.RowsAffected() == 0 {
			return pgx.ErrNoRows
		}
		return err
	}))
}

// ---------------------------------------------------------------------------
// Merek
// ---------------------------------------------------------------------------

var nonSlug = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(s string) string {
	return strings.Trim(nonSlug.ReplaceAllString(strings.ToLower(s), "-"), "-")
}

func (h *Handler) catalogSaveBrand(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name        string   `json:"name"`
		Slug        string   `json:"slug"`
		Verified    bool     `json:"verified"`
		Active      *bool    `json:"active"`
		CategoryIDs []string `json:"category_ids"`
	}
	if !decode(w, r, &in) {
		return
	}
	in.Name = strings.TrimSpace(in.Name)
	if c := utf8.RuneCountInString(in.Name); c < 1 || c > 80 {
		writeFieldErr(w, "name")
		return
	}
	in.Slug = strings.ToLower(strings.TrimSpace(in.Slug))
	if in.Slug == "" {
		in.Slug = slugify(in.Name)
	}
	if !catSlugRe.MatchString(in.Slug) || len(in.Slug) > 60 {
		writeFieldErr(w, "slug")
		return
	}
	active := in.Active == nil || *in.Active
	bid := chi.URLParam(r, "bid")
	h.catalogReply(w, r, h.catalogTx(r, func(ctx context.Context, tx pgx.Tx) error {
		var id int
		var err error
		if bid == "" {
			err = tx.QueryRow(ctx, `INSERT INTO brands (slug, name, verified, active) VALUES ($1, $2, $3, $4) RETURNING id`,
				in.Slug, in.Name, in.Verified, active).Scan(&id)
		} else {
			err = tx.QueryRow(ctx, `UPDATE brands SET slug = $2, name = $3, verified = $4, active = $5 WHERE id = $1::int RETURNING id`,
				bid, in.Slug, in.Name, in.Verified, active).Scan(&id)
		}
		var pe *pgconn.PgError
		if errors.As(err, &pe) && pe.Code == "23505" {
			return catalogFieldError{"slug", "another brand already uses this name"}
		}
		if errors.As(err, &pe) && pe.Code == "22P02" {
			return pgx.ErrNoRows
		}
		if err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `DELETE FROM brand_categories WHERE brand_id = $1`, id); err != nil {
			return err
		}
		for _, c := range in.CategoryIDs {
			tag, err := tx.Exec(ctx, `INSERT INTO brand_categories (brand_id, category_id) SELECT $1, id FROM categories WHERE id = $2 ON CONFLICT DO NOTHING`, id, c)
			if err != nil {
				return err
			}
			if tag.RowsAffected() == 0 {
				var dup bool
				_ = tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM categories WHERE id = $1)`, c).Scan(&dup)
				if !dup {
					return catalogFieldError{"category_ids", "unknown category"}
				}
			}
		}
		return nil
	}))
}

var logoTypes = map[string]string{"image/jpeg": "jpg", "image/png": "png", "image/webp": "webp"}

// POST /brands/{bid}/logo (multipart: file) — logo publik, maks. 2 MB, isi file diperiksa (bukan nama file).
func (h *Handler) catalogBrandLogo(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 3<<20)
	f, _, err := r.FormFile("file")
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad_upload", "send multipart with field file")
		return
	}
	defer f.Close()
	data, err := io.ReadAll(io.LimitReader(f, (2<<20)+1))
	if err != nil || len(data) == 0 {
		writeErr(w, http.StatusBadRequest, "bad_upload", "could not read file")
		return
	}
	if len(data) > 2<<20 {
		writeErr(w, http.StatusRequestEntityTooLarge, "too_large", "logo must be 2 MB or smaller")
		return
	}
	ext, ok := logoTypes[http.DetectContentType(data)]
	if !ok {
		writeErr(w, http.StatusUnsupportedMediaType, "bad_type", "use a JPG, PNG or WebP image")
		return
	}
	bid, err := strconv.Atoi(chi.URLParam(r, "bid"))
	if err != nil {
		writeErr(w, http.StatusNotFound, "not_found", "not found")
		return
	}
	var exists bool
	if err := h.pool.QueryRow(r.Context(), `SELECT EXISTS (SELECT 1 FROM brands WHERE id = $1)`, bid).Scan(&exists); err != nil || !exists {
		writeErr(w, http.StatusNotFound, "not_found", "not found")
		return
	}
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	key := "public/brands/" + strconv.Itoa(bid) + "-" + hex.EncodeToString(b) + "." + ext
	if _, err := h.files.Put(r.Context(), key, bytes.NewReader(data)); err != nil {
		h.log.Error("logo merek", "err", err)
		writeErr(w, http.StatusInternalServerError, "internal", "could not save logo")
		return
	}
	h.catalogReply(w, r, h.catalogTx(r, func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `UPDATE brands SET logo_url = $2 WHERE id = $1`, bid, h.files.PublicURL(key))
		return err
	}))
}

func (h *Handler) catalogDeleteBrand(w http.ResponseWriter, r *http.Request) {
	bid, err := strconv.Atoi(chi.URLParam(r, "bid"))
	if err != nil {
		writeErr(w, http.StatusNotFound, "not_found", "not found")
		return
	}
	h.catalogReply(w, r, h.catalogTx(r, func(ctx context.Context, tx pgx.Tx) error {
		var used bool
		if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM products WHERE brand_id = $1)`, bid).Scan(&used); err != nil {
			return err
		}
		if used {
			return errHasItems
		}
		tag, err := tx.Exec(ctx, `DELETE FROM brands WHERE id = $1`, bid)
		if err == nil && tag.RowsAffected() == 0 {
			return pgx.ErrNoRows
		}
		return err
	}))
}

// ---------------------------------------------------------------------------
// Kata terlarang (cek otomatis saat seller mengirim produk ke moderasi)
// ---------------------------------------------------------------------------

func (h *Handler) catalogAddKeyword(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Keyword  string `json:"keyword"`
		Lang     string `json:"lang"`
		Category string `json:"category"`
		Action   string `json:"action"`
	}
	if !decode(w, r, &in) {
		return
	}
	in.Keyword = strings.ToLower(strings.TrimSpace(in.Keyword))
	in.Category = strings.ToLower(strings.TrimSpace(in.Category))
	if c := utf8.RuneCountInString(in.Keyword); c < 2 || c > 80 {
		writeFieldErr(w, "keyword")
		return
	}
	if in.Lang != "en" && in.Lang != "km" && in.Lang != "zh" {
		writeFieldErr(w, "lang")
		return
	}
	if !regexp.MustCompile(`^[a-z][a-z_]{1,39}$`).MatchString(in.Category) {
		writeFieldErr(w, "category")
		return
	}
	if !kwActions[in.Action] {
		writeFieldErr(w, "action")
		return
	}
	h.catalogReply(w, r, h.catalogTx(r, func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `INSERT INTO prohibited_keywords (keyword, lang, category, action) VALUES ($1, $2::lang_code, $3, $4)`,
			in.Keyword, in.Lang, in.Category, in.Action)
		var pe *pgconn.PgError
		if errors.As(err, &pe) && pe.Code == "23505" {
			return catalogFieldError{"keyword", "already in the list"}
		}
		return err
	}))
}

func (h *Handler) catalogDeleteKeyword(w http.ResponseWriter, r *http.Request) {
	kid, err := strconv.Atoi(chi.URLParam(r, "kid"))
	if err != nil {
		writeErr(w, http.StatusNotFound, "not_found", "not found")
		return
	}
	h.catalogReply(w, r, h.catalogTx(r, func(ctx context.Context, tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, `DELETE FROM prohibited_keywords WHERE id = $1`, kid)
		if err == nil && tag.RowsAffected() == 0 {
			return pgx.ErrNoRows
		}
		return err
	}))
}
