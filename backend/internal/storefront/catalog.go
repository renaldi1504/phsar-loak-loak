package storefront

// Katalog publik (hanya baca). Satu aturan untuk semua halaman: yang tampil hanya produk
// tayang — status active, disetujui moderator, tidak dihapus, dari toko yang aktif.

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"phsar/internal/storage"
)

const visible = `p.deleted_at IS NULL AND p.status IN ('active', 'sold_out') AND p.moderation_status = 'approved'
	AND st.status = 'active' AND st.deleted_at IS NULL`

var errNotFound = errors.New("not found")

type Catalog struct {
	Pool  *pgxpool.Pool
	Files storage.Storage
}

type Card struct {
	ID         string
	Slug       string
	Name       string
	MinPrice   string
	MaxPrice   string
	Compare    string // harga coret tertinggi (kosong = tidak ada diskon)
	CoverURL   string
	PosterURL  string
	HasVideo   bool
	Duration   int
	StoreName  string
	StoreSlug  string
	Condition  string
	Stock      int
	CategoryID string
	VariantID  string     // diisi kalau produk hanya punya 1 varian aktif → tombol keranjang di kartu langsung menambah
	FlashEnds  *time.Time // promos: flash sale berlaku sampai
}

// Kolom kartu: nama/slug dalam bahasa lang, kalau belum diterjemahkan pakai bahasa asli seller.
const cardCols = `p.id,
	coalesce(t.slug, ts.slug), coalesce(t.name, ts.name),
	p.min_price::text, p.max_price::text,
	coalesce((SELECT max(v.compare_at_price)::text FROM product_variants v WHERE v.product_id = p.id AND v.deleted_at IS NULL AND v.status = 'active'), ''),
	coalesce((SELECT coalesce(m.image_url, m.storage_key) FROM product_media m WHERE m.product_id = p.id AND m.deleted_at IS NULL AND m.kind = 'photo' ORDER BY m.sort_order LIMIT 1), ''),
	coalesce(vm.poster_url, ''), vm.id IS NOT NULL, coalesce(vm.duration_sec, 0),
	st.name, st.slug, p.condition::text,
	coalesce((SELECT sum(v.stock_qty - v.reserved_qty) FROM product_variants v WHERE v.product_id = p.id AND v.deleted_at IS NULL AND v.status = 'active'), 0),
	p.category_id,
	coalesce((SELECT CASE WHEN count(*) = 1 THEN min(v.id::text) END FROM product_variants v WHERE v.product_id = p.id AND v.deleted_at IS NULL AND v.status = 'active'), '')`

const cardJoins = `FROM products p
	JOIN stores st ON st.id = p.store_id
	JOIN product_translations ts ON ts.product_id = p.id AND ts.lang = p.source_lang
	LEFT JOIN product_translations t ON t.product_id = p.id AND t.lang = $1::lang_code
	LEFT JOIN LATERAL (SELECT m.id, m.poster_url, m.duration_sec FROM product_media m
		WHERE m.product_id = p.id AND m.deleted_at IS NULL AND m.kind = 'video' ORDER BY m.sort_order LIMIT 1) vm ON true`

func (c *Catalog) url(keyOrURL string) string {
	if keyOrURL == "" || strings.HasPrefix(keyOrURL, "http") || strings.HasPrefix(keyOrURL, "/") {
		return keyOrURL
	}
	return c.Files.PublicURL(keyOrURL)
}

func (c *Catalog) cards(ctx context.Context, lang, where, order string, limit int, args ...any) ([]Card, error) {
	q := `SELECT ` + cardCols + ` ` + cardJoins + ` WHERE ` + visible
	if where != "" {
		q += ` AND ` + where
	}
	q += ` ORDER BY ` + order + ` LIMIT ` + itoa(limit)
	rows, err := c.Pool.Query(ctx, q, append([]any{lang}, args...)...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Card{}
	for rows.Next() {
		var k Card
		var min, max *string
		if err := rows.Scan(&k.ID, &k.Slug, &k.Name, &min, &max, &k.Compare, &k.CoverURL, &k.PosterURL, &k.HasVideo, &k.Duration,
			&k.StoreName, &k.StoreSlug, &k.Condition, &k.Stock, &k.CategoryID, &k.VariantID); err != nil {
			return nil, err
		}
		k.MinPrice, k.MaxPrice = deref(min), deref(max)
		k.CoverURL, k.PosterURL = c.url(k.CoverURL), c.url(k.PosterURL)
		if k.CoverURL == "" {
			k.CoverURL = k.PosterURL
		}
		out = append(out, k)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	c.flashCards(ctx, out) // promos
	return out, nil
}

// CardsByIDs: kartu produk tampil untuk id tertentu, urutan sesuai ids (dipakai cari-dengan-foto).
func (c *Catalog) CardsByIDs(ctx context.Context, lang string, ids []string, limit int) ([]Card, error) {
	if len(ids) == 0 {
		return []Card{}, nil
	}
	cards, err := c.cards(ctx, lang, `p.id = ANY($2::uuid[])`, `p.id`, len(ids), ids)
	if err != nil {
		return nil, err
	}
	by := make(map[string]Card, len(cards))
	for _, k := range cards {
		by[k.ID] = k
	}
	out := make([]Card, 0, limit)
	for _, id := range ids {
		if k, ok := by[id]; ok && len(out) < limit {
			out = append(out, k)
		}
	}
	return out, nil
}

// ---- Kategori ----

type Category struct {
	ID, Name, Slug, Icon string
	ParentID             string     // "" = kategori utama
	Count                int        // kategori utama: termasuk produk di sub-kategorinya
	Children             []Category // sub-kategori (hanya di kategori utama)
}

// inCategory: produk di kategori itu atau sub-kategorinya ($n = id kategori; untuk sub-kategori = dirinya sendiri).
func inCategory(n string) string {
	return `p.category_id IN (SELECT k.id FROM categories k WHERE k.deleted_at IS NULL AND (k.id = ` + n + ` OR k.parent_id = ` + n + `))`
}

// Categories: kategori utama beserta sub-kategorinya (2 tingkat), urut sesuai admin.
func (c *Catalog) Categories(ctx context.Context, lang string) ([]Category, error) {
	rows, err := c.Pool.Query(ctx, `
		SELECT c.id, coalesce(c.parent_id, ''), coalesce(t.name, te.name, c.id), coalesce(t.slug, te.slug, c.id), coalesce(c.icon, ''),
		       (SELECT count(*) FROM products p JOIN stores st ON st.id = p.store_id WHERE `+inCategory("c.id")+` AND `+visible+`)
		FROM categories c
		LEFT JOIN categories pc ON pc.id = c.parent_id
		LEFT JOIN category_translations t ON t.category_id = c.id AND t.lang = $1::lang_code
		LEFT JOIN category_translations te ON te.category_id = c.id AND te.lang = 'en'
		WHERE c.deleted_at IS NULL AND (c.parent_id IS NULL OR pc.deleted_at IS NULL)
		ORDER BY c.parent_id NULLS FIRST, c.sort_order, c.id`, lang)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Category{}
	idx := map[string]int{}
	for rows.Next() {
		var k Category
		if err := rows.Scan(&k.ID, &k.ParentID, &k.Name, &k.Slug, &k.Icon, &k.Count); err != nil {
			return nil, err
		}
		if k.ParentID == "" {
			idx[k.ID] = len(out)
			out = append(out, k)
		} else if i, ok := idx[k.ParentID]; ok {
			out[i].Children = append(out[i].Children, k)
		}
	}
	return out, rows.Err()
}

// flatCategories: kategori utama + semua sub-kategori dalam satu daftar (untuk mencari slug/ID).
func flatCategories(cs []Category) []Category {
	out := make([]Category, 0, len(cs)*6)
	for _, c := range cs {
		out = append(out, c)
		out = append(out, c.Children...)
	}
	return out
}

// ---- Daftar ----

func (c *Catalog) VideoFeed(ctx context.Context, lang string, limit int) ([]Card, error) {
	return c.cards(ctx, lang, "vm.id IS NOT NULL", "p.published_at DESC NULLS LAST, p.id", limit)
}

func (c *Catalog) Newest(ctx context.Context, lang string, limit, offset int) ([]Card, error) {
	return c.cards(ctx, lang, "", "p.published_at DESC NULLS LAST, p.id OFFSET "+itoa(offset), limit)
}

var sorts = map[string]string{
	"new":        "p.published_at DESC NULLS LAST, p.id",
	"price_asc":  "p.min_price ASC NULLS LAST, p.id",
	"price_desc": "p.max_price DESC NULLS LAST, p.id",
	"popular":    soldOrder,
}

// soldOrder: terlaris = unit terjual (pesanan dibayar) 30 hari terakhir.
const soldOrder = `(SELECT coalesce(sum(oi.qty), 0) FROM order_items oi JOIN orders o ON o.id = oi.order_id
	WHERE oi.product_id = p.id AND o.paid_at > now() - interval '30 days') DESC, p.published_at DESC NULLS LAST, p.id`

func (c *Catalog) ByCategory(ctx context.Context, lang, categoryID, sort string, limit, offset int) ([]Card, error) {
	order, ok := sorts[sort]
	if !ok {
		order = sorts["new"]
	}
	// Kategori utama ikut menampilkan produk di sub-kategorinya.
	return c.cards(ctx, lang, inCategory("$2"), order+" OFFSET "+itoa(offset), limit, categoryID)
}

// searchWhere: cocokkan nama/deskripsi produk di bahasa apa pun (seller menulis Khmer, pembeli mencari Inggris, dst.),
// nama toko, atau nama kategori. $2 = pola ILIKE. Nanti dilengkapi pencarian vektor (pgvector).
const searchWhere = `(
		EXISTS (SELECT 1 FROM product_translations x WHERE x.product_id = p.id AND (x.name ILIKE $2 OR x.description ILIKE $2))
		OR st.name ILIKE $2
		OR EXISTS (SELECT 1 FROM categories cc JOIN category_translations ct ON ct.category_id IN (cc.id, cc.parent_id)
		           WHERE cc.id = p.category_id AND ct.name ILIKE $2))`

func likePattern(q string) string {
	return "%" + strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(strings.TrimSpace(q)) + "%"
}

// Search: satu halaman hasil cari (limit/offset); nama yang cocok tampil lebih dulu.
func (c *Catalog) Search(ctx context.Context, lang, q string, limit, offset int) ([]Card, error) {
	if strings.TrimSpace(q) == "" {
		return []Card{}, nil
	}
	return c.cards(ctx, lang, searchWhere,
		`(EXISTS (SELECT 1 FROM product_translations x WHERE x.product_id = p.id AND x.name ILIKE $2)) DESC, p.published_at DESC NULLS LAST, p.id OFFSET `+itoa(offset),
		limit, likePattern(q))
}

// SearchCount: jumlah semua hasil cari (untuk "N hasil").
func (c *Catalog) SearchCount(ctx context.Context, q string) (int, error) {
	if strings.TrimSpace(q) == "" {
		return 0, nil
	}
	var n int
	err := c.Pool.QueryRow(ctx, `SELECT count(*) FROM products p JOIN stores st ON st.id = p.store_id
		WHERE $1::text IS NOT NULL AND `+visible+` AND `+searchWhere, "", likePattern(q)).Scan(&n)
	return n, err
}

func (c *Catalog) ByStore(ctx context.Context, lang, storeID string, limit, offset int) ([]Card, error) {
	return c.cards(ctx, lang, "p.store_id = $2", "p.published_at DESC NULLS LAST, p.id OFFSET "+itoa(offset), limit, storeID)
}

func (c *Catalog) Related(ctx context.Context, lang, categoryID, excludeID string, limit int) ([]Card, error) {
	return c.cards(ctx, lang, "p.category_id = $2 AND p.id <> $3", "p.published_at DESC NULLS LAST, p.id", limit, categoryID, excludeID)
}

// ---- Detail produk ----

type Product struct {
	ID, Slug, Name, Description, SourceLang, Lang string
	Translated                                    bool   // teks hasil terjemahan (bukan bahasa asli seller)
	OriginalName, OriginalDescription             string // teks asli seller
	Condition, ConditionNote                      string
	UsageMonths                                   *int
	WarrantyType                                  string
	WarrantyMonths, ReturnDays                    int
	CategoryID, CategoryName, CategorySlug        string
	ParentName, ParentSlug                        string // kategori utama (kosong kalau produk langsung di kategori utama)
	Store                                         StoreCard
	Media                                         []Media
	Options                                       []Option
	Variants                                      []Variant
	MinPrice, MaxPrice, Compare                   string
	Stock                                         int
	Slugs                                         map[string]string // slug per bahasa (untuk hreflang & ganti bahasa)
	PublishedAt                                   *time.Time
	UpdatedAt                                     time.Time
	FlashEnds                                     *time.Time // promos
	Low30                                         string     // pricehistory: harga terendah 30 hari sebelum flash (kosong = tidak flash)
}

type StoreCard struct {
	ID, Name, Slug, LogoURL, BannerURL, About string
	Products                                  int
	Since                                     time.Time
}

type Media struct {
	Kind, URL, PosterURL string
	Width, Height        int
}

type Option struct {
	Kind, Name string
	Values     []OptionValue
}

type OptionValue struct {
	Value, Swatch string
}

type Variant struct {
	ID       string   `json:"id"`
	Options  []string `json:"o"`
	Price    string   `json:"p"`
	Compare  string   `json:"c,omitempty"`
	Stock    int      `json:"s"`
	SKU      string   `json:"-"`
	Position []int    `json:"i"`           // indeks nilai per opsi
	Flash    bool     `json:"f,omitempty"` // promos: harga flash sale
}

// ProductBySlug mencari produk tayang lewat slug di bahasa mana pun. canonical = slug untuk lang.
func (c *Catalog) ProductBySlug(ctx context.Context, lang, slug string) (*Product, error) {
	var id string
	err := c.Pool.QueryRow(ctx, `
		SELECT p.id FROM product_translations x JOIN products p ON p.id = x.product_id JOIN stores st ON st.id = p.store_id
		WHERE x.slug = $1 AND `+visible+` ORDER BY (x.lang = $2::lang_code) DESC LIMIT 1`, slug, lang).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errNotFound
	}
	if err != nil {
		return nil, err
	}
	return c.Product(ctx, lang, id)
}

func (c *Catalog) Product(ctx context.Context, lang, id string) (*Product, error) {
	p := &Product{Lang: lang, Slugs: map[string]string{}}
	var usage *int16
	var wm, rd int16
	var min, max *string
	var storeLogo, storeBanner *string
	err := c.Pool.QueryRow(ctx, `
		SELECT p.id, p.source_lang::text, p.condition::text, coalesce(p.condition_note, ''), p.usage_months,
		       p.warranty_type, p.warranty_months, p.return_days, p.category_id,
		       coalesce(ct.name, cte.name, p.category_id), coalesce(ct.slug, cte.slug, p.category_id),
		       coalesce(pt.name, pte.name, ''), coalesce(pt.slug, pte.slug, ''),
		       p.min_price::text, p.max_price::text, p.published_at, p.updated_at,
		       st.id, st.name, st.slug, st.logo_key, st.banner_key, st.created_at,
		       coalesce((SELECT about FROM store_translations x WHERE x.store_id = st.id AND x.lang = $2::lang_code AND x.about <> ''),
		                (SELECT about FROM store_translations x WHERE x.store_id = st.id AND x.source = 'seller' LIMIT 1), '')
		FROM products p JOIN stores st ON st.id = p.store_id
		LEFT JOIN category_translations ct ON ct.category_id = p.category_id AND ct.lang = $2::lang_code
		LEFT JOIN category_translations cte ON cte.category_id = p.category_id AND cte.lang = 'en'
		LEFT JOIN categories pcat ON pcat.id = (SELECT parent_id FROM categories WHERE id = p.category_id)
		LEFT JOIN category_translations pt ON pt.category_id = pcat.id AND pt.lang = $2::lang_code
		LEFT JOIN category_translations pte ON pte.category_id = pcat.id AND pte.lang = 'en'
		WHERE p.id = $1 AND `+visible, id, lang).
		Scan(&p.ID, &p.SourceLang, &p.Condition, &p.ConditionNote, &usage, &p.WarrantyType, &wm, &rd, &p.CategoryID,
			&p.CategoryName, &p.CategorySlug, &p.ParentName, &p.ParentSlug, &min, &max, &p.PublishedAt, &p.UpdatedAt,
			&p.Store.ID, &p.Store.Name, &p.Store.Slug, &storeLogo, &storeBanner, &p.Store.Since, &p.Store.About)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errNotFound
	}
	if err != nil {
		return nil, err
	}
	p.MinPrice, p.MaxPrice = deref(min), deref(max)
	p.WarrantyMonths, p.ReturnDays = int(wm), int(rd)
	if usage != nil {
		u := int(*usage)
		p.UsageMonths = &u
	}
	p.Store.LogoURL, p.Store.BannerURL = c.url(deref(storeLogo)), c.url(deref(storeBanner))

	// Teks: bahasa yang diminta, kalau belum ada → bahasa asli.
	rows, err := c.Pool.Query(ctx, `SELECT lang::text, name, description, slug FROM product_translations WHERE product_id = $1`, id)
	if err != nil {
		return nil, err
	}
	texts := map[string][3]string{}
	for rows.Next() {
		var l, n, d, s string
		if err := rows.Scan(&l, &n, &d, &s); err != nil {
			rows.Close()
			return nil, err
		}
		texts[l] = [3]string{n, d, s}
		p.Slugs[l] = s
	}
	rows.Close()
	orig := texts[p.SourceLang]
	p.OriginalName, p.OriginalDescription = orig[0], orig[1]
	if tr, ok := texts[lang]; ok {
		p.Name, p.Description, p.Slug = tr[0], tr[1], tr[2]
		p.Translated = lang != p.SourceLang
	} else {
		p.Name, p.Description, p.Slug = orig[0], orig[1], orig[2]
	}
	for _, l := range []string{"en", "km", "zh"} {
		if p.Slugs[l] == "" {
			p.Slugs[l] = orig[2]
		}
	}

	// Media: video dulu, lalu foto, lalu foto kondisi.
	rows, err = c.Pool.Query(ctx, `SELECT kind::text, coalesce(CASE WHEN kind = 'video' THEN video_url ELSE image_url END, storage_key),
		coalesce(poster_url, ''), coalesce(width, 0), coalesce(height, 0)
		FROM product_media WHERE product_id = $1 AND deleted_at IS NULL
		ORDER BY CASE kind WHEN 'video' THEN 0 WHEN 'photo' THEN 1 ELSE 2 END, sort_order`, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var m Media
		if err := rows.Scan(&m.Kind, &m.URL, &m.PosterURL, &m.Width, &m.Height); err != nil {
			rows.Close()
			return nil, err
		}
		m.URL, m.PosterURL = c.url(m.URL), c.url(m.PosterURL)
		p.Media = append(p.Media, m)
	}
	rows.Close()

	// Opsi & nilai (bahasa diminta → asli).
	rows, err = c.Pool.Query(ctx, `
		SELECT o.position, o.kind, coalesce(o.name_i18n->>$2, o.name_i18n->>$3, ''), coalesce(v.value_i18n->>$2, v.value_i18n->>$3, ''),
		       coalesce(v.swatch_hex, ''), coalesce(v.value_i18n->>$3, '')
		FROM product_options o JOIN product_option_values v ON v.option_id = o.id
		WHERE o.product_id = $1 ORDER BY o.position, v.position`, id, lang, p.SourceLang)
	if err != nil {
		return nil, err
	}
	valueIdx := []map[string]int{} // nilai asli → indeks, per opsi
	lastPos := -1
	for rows.Next() {
		var pos int16
		var kind, name, value, swatch, origValue string
		if err := rows.Scan(&pos, &kind, &name, &value, &swatch, &origValue); err != nil {
			rows.Close()
			return nil, err
		}
		if int(pos) != lastPos {
			p.Options = append(p.Options, Option{Kind: kind, Name: name})
			valueIdx = append(valueIdx, map[string]int{})
			lastPos = int(pos)
		}
		o := &p.Options[len(p.Options)-1]
		valueIdx[len(valueIdx)-1][origValue] = len(o.Values)
		o.Values = append(o.Values, OptionValue{Value: value, Swatch: swatch})
	}
	rows.Close()

	rows, err = c.Pool.Query(ctx, `
		SELECT v.id, v.sku, v.price::text, coalesce(v.compare_at_price::text, ''), greatest(v.stock_qty - v.reserved_qty, 0),
		       coalesce(array_agg(ov.value_i18n->>$2 ORDER BY o.position) FILTER (WHERE ov.id IS NOT NULL), '{}')
		FROM product_variants v
		LEFT JOIN variant_option_values vov ON vov.variant_id = v.id
		LEFT JOIN product_option_values ov ON ov.id = vov.option_value_id
		LEFT JOIN product_options o ON o.id = vov.option_id
		WHERE v.product_id = $1 AND v.deleted_at IS NULL AND v.status = 'active'
		GROUP BY v.id ORDER BY v.is_default DESC, v.price, v.sku`, id, p.SourceLang)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var v Variant
		var origVals []string
		if err := rows.Scan(&v.ID, &v.SKU, &v.Price, &v.Compare, &v.Stock, &origVals); err != nil {
			rows.Close()
			return nil, err
		}
		for i, ov := range origVals {
			idx := 0
			if i < len(valueIdx) {
				idx = valueIdx[i][ov]
			}
			v.Position = append(v.Position, idx)
			if i < len(p.Options) && idx < len(p.Options[i].Values) {
				v.Options = append(v.Options, p.Options[i].Values[idx].Value)
			}
		}
		if v.Compare != "" && num(v.Compare) > num(p.Compare) {
			p.Compare = v.Compare
		}
		p.Stock += v.Stock
		p.Variants = append(p.Variants, v)
	}
	rows.Close()
	_ = c.Pool.QueryRow(ctx, `SELECT count(*) FROM products p JOIN stores st ON st.id = p.store_id WHERE p.store_id = $1 AND `+visible,
		p.Store.ID).Scan(&p.Store.Products)
	c.flashProduct(ctx, p) // promos
	return p, nil
}

func (p *Product) VariantsJSON() string {
	b, _ := json.Marshal(p.Variants)
	return string(b)
}

// ---- Toko ----

func (c *Catalog) StoreBySlug(ctx context.Context, lang, slug string) (*StoreCard, error) {
	s := &StoreCard{}
	var logo, banner *string
	err := c.Pool.QueryRow(ctx, `
		SELECT st.id, st.name, st.slug, st.logo_key, st.banner_key, st.created_at,
		       coalesce((SELECT about FROM store_translations x WHERE x.store_id = st.id AND x.lang = $2::lang_code AND x.about <> ''),
		                (SELECT about FROM store_translations x WHERE x.store_id = st.id AND x.source = 'seller' LIMIT 1), ''),
		       (SELECT count(*) FROM products p WHERE p.store_id = st.id AND p.deleted_at IS NULL AND p.status IN ('active', 'sold_out') AND p.moderation_status = 'approved')
		FROM stores st WHERE st.slug = $1 AND st.status = 'active' AND st.deleted_at IS NULL`, slug, lang).
		Scan(&s.ID, &s.Name, &s.Slug, &logo, &banner, &s.Since, &s.About, &s.Products)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errNotFound
	}
	if err != nil {
		return nil, err
	}
	s.LogoURL, s.BannerURL = c.url(deref(logo)), c.url(deref(banner))
	return s, nil
}

// RielPerUSD: kurs terbaru (untuk tampilan "≈ 20.000៛").
func (c *Catalog) RielPerUSD(ctx context.Context) float64 {
	var r float64
	if err := c.Pool.QueryRow(ctx, `SELECT riel_per_usd::float8 FROM exchange_rates ORDER BY rate_date DESC, id DESC LIMIT 1`).Scan(&r); err != nil || r <= 0 {
		return 4000
	}
	return r
}

// ---- Sitemap ----

type SitemapEntry struct {
	Slugs   map[string]string
	Updated time.Time
}

func (c *Catalog) SitemapProducts(ctx context.Context, limit int) ([]SitemapEntry, error) {
	rows, err := c.Pool.Query(ctx, `
		SELECT p.updated_at, jsonb_object_agg(x.lang::text, x.slug)
		FROM products p JOIN stores st ON st.id = p.store_id JOIN product_translations x ON x.product_id = p.id
		WHERE `+visible+` GROUP BY p.id, p.updated_at ORDER BY p.updated_at DESC LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []SitemapEntry{}
	for rows.Next() {
		var e SitemapEntry
		var raw []byte
		if err := rows.Scan(&e.Updated, &raw); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(raw, &e.Slugs)
		out = append(out, e)
	}
	return out, rows.Err()
}

func (c *Catalog) SitemapStores(ctx context.Context) ([]string, error) {
	rows, err := c.Pool.Query(ctx, `SELECT slug FROM stores WHERE status = 'active' AND deleted_at IS NULL ORDER BY created_at`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []string{}
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func itoa(i int) string {
	if i < 0 {
		i = 0
	}
	return strconv.Itoa(i)
}

func num(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
