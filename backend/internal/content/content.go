// Package content: konten & SEO website pembeli yang diatur super admin —
// banner beranda, pita pengumuman, koleksi pilihan, meta SEO per halaman, gambar OG bawaan.
//
// Website membaca lewat Snapshot (disimpan di memori 30 detik; disegarkan langsung setelah admin menyimpan).
package content

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"phsar/internal/platform/db"
	"phsar/internal/storage"
)

var Langs = []string{"en", "km", "zh"}

// CacheTTL: lama snapshot disimpan (variabel supaya tes bisa mematikannya).
var CacheTTL = 30 * time.Second

// Text: teks per bahasa {"en": …, "km": …, "zh": …}.
type Text map[string]string

// In: teks bahasa lang, cadangan Inggris lalu bahasa mana pun yang terisi.
func (t Text) In(lang string) string {
	if s := t[lang]; s != "" {
		return s
	}
	if s := t["en"]; s != "" {
		return s
	}
	for _, l := range Langs {
		if t[l] != "" {
			return t[l]
		}
	}
	return ""
}

type Banner struct {
	ID        string     `json:"id"`
	ImageKey  string     `json:"image_key"`
	ImageURL  string     `json:"image_url"`
	Link      string     `json:"link"`
	Title     Text       `json:"title"`
	Subtitle  Text       `json:"subtitle"`
	StartsAt  *time.Time `json:"starts_at"`
	EndsAt    *time.Time `json:"ends_at"`
	SortOrder int        `json:"sort_order"`
	Active    bool       `json:"active"`
}

func live(active bool, from, to *time.Time, now time.Time) bool {
	return active && (from == nil || !now.Before(*from)) && (to == nil || now.Before(*to))
}

func (b Banner) Live(now time.Time) bool { return live(b.Active, b.StartsAt, b.EndsAt, now) }

type Collection struct {
	ID         string   `json:"id"`
	Title      Text     `json:"title"`
	Source     string   `json:"source"` // category · products
	CategoryID string   `json:"category_id"`
	ProductIDs []string `json:"product_ids"`
	MaxItems   int      `json:"max_items"`
	SortOrder  int      `json:"sort_order"`
	Active     bool     `json:"active"`
}

type Announcement struct {
	Active   bool       `json:"active"`
	Text     Text       `json:"text"`
	Link     string     `json:"link"`
	Color    string     `json:"color"` // pink · navy · green · amber
	StartsAt *time.Time `json:"starts_at"`
	EndsAt   *time.Time `json:"ends_at"`
}

func (a Announcement) Live(now time.Time) bool {
	return live(a.Active, a.StartsAt, a.EndsAt, now) && (a.Text["en"] != "" || a.Text["km"] != "" || a.Text["zh"] != "")
}

type SEO struct {
	Page        string `json:"page"` // home · c:<category_id>
	Title       Text   `json:"title"`
	Description Text   `json:"description"`
}

// Snapshot: semua konten (termasuk yang tidak aktif / di luar jadwal; penyaringan saat tampil).
type Snapshot struct {
	Banners      []Banner       `json:"banners"`
	Collections  []Collection   `json:"collections"`
	Seasons      []Season       `json:"seasons"`
	Announcement Announcement   `json:"announcement"`
	SEO          map[string]SEO `json:"-"`
	OGImageKey   string         `json:"og_image_key"`
	OGImageURL   string         `json:"og_image_url"`
}

// LiveBanners / LiveCollections / LiveAnnouncement: yang tampil sekarang.
func (s *Snapshot) LiveBanners(now time.Time) []Banner {
	out := []Banner{}
	for _, b := range s.Banners {
		if b.Live(now) {
			out = append(out, b)
		}
	}
	return out
}

func (s *Snapshot) LiveCollections() []Collection {
	out := []Collection{}
	for _, c := range s.Collections {
		if c.Active {
			out = append(out, c)
		}
	}
	return out
}

func (s *Snapshot) LiveAnnouncement(now time.Time) *Announcement {
	if s.Announcement.Live(now) {
		a := s.Announcement
		return &a
	}
	return nil
}

// Meta: judul & deskripsi pengganti untuk halaman (kosong = pakai bawaan).
func (s *Snapshot) Meta(page, lang string) (title, desc string) {
	m, ok := s.SEO[page]
	if !ok {
		return "", ""
	}
	return strings.TrimSpace(m.Title[lang]), strings.TrimSpace(m.Description[lang]) // per bahasa, tanpa cadangan bahasa lain
}

// ---------------------------------------------------------------------------

type Store struct {
	Pool  *pgxpool.Pool
	Files storage.Storage

	mu   sync.Mutex
	snap *Snapshot
	at   time.Time
}

func (s *Store) url(key string) string {
	if key == "" || s.Files == nil {
		return ""
	}
	return s.Files.PublicURL(key)
}

// Snapshot dari cache (disegarkan setiap CacheTTL). Gagal baca → snapshot kosong (website tetap jalan).
func (s *Store) Snapshot(ctx context.Context) *Snapshot {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.snap != nil && time.Since(s.at) < CacheTTL {
		return s.snap
	}
	snap, err := s.load(ctx)
	if err != nil {
		if s.snap != nil {
			return s.snap
		}
		return &Snapshot{SEO: map[string]SEO{}}
	}
	s.snap, s.at = snap, time.Now()
	return snap
}

func (s *Store) invalidate() {
	s.mu.Lock()
	s.snap = nil
	s.mu.Unlock()
}

func (s *Store) load(ctx context.Context) (*Snapshot, error) {
	snap := &Snapshot{Banners: []Banner{}, Collections: []Collection{}, Seasons: []Season{}, SEO: map[string]SEO{}, Announcement: Announcement{Text: Text{}, Color: "pink"}}
	rows, err := s.Pool.Query(ctx, `SELECT id, image_key, link, title, subtitle, starts_at, ends_at, sort_order, active
		FROM content_banners ORDER BY sort_order, created_at`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var b Banner
		if err := rows.Scan(&b.ID, &b.ImageKey, &b.Link, &b.Title, &b.Subtitle, &b.StartsAt, &b.EndsAt, &b.SortOrder, &b.Active); err != nil {
			rows.Close()
			return nil, err
		}
		b.ImageURL = s.url(b.ImageKey)
		snap.Banners = append(snap.Banners, b)
	}
	rows.Close()
	rows, err = s.Pool.Query(ctx, `SELECT id, title, source, coalesce(category_id, ''), product_ids::text[], max_items, sort_order, active
		FROM content_collections ORDER BY sort_order, created_at`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var c Collection
		if err := rows.Scan(&c.ID, &c.Title, &c.Source, &c.CategoryID, &c.ProductIDs, &c.MaxItems, &c.SortOrder, &c.Active); err != nil {
			rows.Close()
			return nil, err
		}
		if c.ProductIDs == nil {
			c.ProductIDs = []string{}
		}
		snap.Collections = append(snap.Collections, c)
	}
	rows.Close()
	if err := s.loadSeasons(ctx, snap); err != nil {
		return nil, err
	}
	rows, err = s.Pool.Query(ctx, `SELECT page, title, description FROM seo_meta`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var m SEO
		if err := rows.Scan(&m.Page, &m.Title, &m.Description); err != nil {
			rows.Close()
			return nil, err
		}
		snap.SEO[m.Page] = m
	}
	rows.Close()
	var raw []byte
	if err := s.Pool.QueryRow(ctx, `SELECT value FROM platform_settings WHERE key = 'content.announcement'`).Scan(&raw); err == nil {
		_ = json.Unmarshal(raw, &snap.Announcement)
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	if snap.Announcement.Text == nil {
		snap.Announcement.Text = Text{}
	}
	var og struct {
		ImageKey string `json:"image_key"`
	}
	if err := s.Pool.QueryRow(ctx, `SELECT value FROM platform_settings WHERE key = 'content.og'`).Scan(&raw); err == nil {
		_ = json.Unmarshal(raw, &og)
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	snap.OGImageKey, snap.OGImageURL = og.ImageKey, s.url(og.ImageKey)
	return snap, nil
}

// ---------------------------------------------------------------------------
// Admin: validasi & simpan
// ---------------------------------------------------------------------------

type FieldError struct{ Field, Msg string }

func (e FieldError) Error() string { return e.Field + ": " + e.Msg }

var ErrNotFound = errors.New("tidak ditemukan")

// cleanText: rapikan teks per bahasa, tolak bahasa tak dikenal & teks terlalu panjang.
func cleanText(field string, t Text, max int) (Text, error) {
	out := Text{}
	for k, v := range t {
		ok := false
		for _, l := range Langs {
			ok = ok || l == k
		}
		if !ok {
			return nil, FieldError{field, "unknown language " + k}
		}
		v = strings.TrimSpace(v)
		if utf8.RuneCountInString(v) > max {
			return nil, FieldError{field + "_" + k, "too long"}
		}
		if strings.ContainsAny(v, "\r\n") {
			v = strings.Join(strings.Fields(v), " ")
		}
		if v != "" {
			out[k] = v
		}
	}
	return out, nil
}

// cleanLink: kosong, path situs sendiri (/km/c/fashion), atau https://….
func cleanLink(field, l string) (string, error) {
	l = strings.TrimSpace(l)
	if l == "" {
		return "", nil
	}
	if len(l) > 300 || strings.ContainsAny(l, " \t\r\n\"'<>") {
		return "", FieldError{field, "invalid link"}
	}
	if strings.HasPrefix(l, "/") && !strings.HasPrefix(l, "//") && !strings.HasPrefix(l, "/\\") {
		return l, nil
	}
	if u, err := url.Parse(l); err == nil && u.Scheme == "https" && u.Host != "" {
		return l, nil
	}
	return "", FieldError{field, "use a site path like /km/c/fashion or an https:// link"}
}

func (s *Store) checkImage(ctx context.Context, field, key string) error {
	var ok bool
	if err := s.Pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM uploads WHERE key = $1 AND kind = 'content_image')`, key).Scan(&ok); err != nil {
		return err
	}
	if !ok {
		return FieldError{field, "upload an image first"}
	}
	return nil
}

type BannerInput struct {
	ImageKey  string     `json:"image_key"`
	Link      string     `json:"link"`
	Title     Text       `json:"title"`
	Subtitle  Text       `json:"subtitle"`
	StartsAt  *time.Time `json:"starts_at"`
	EndsAt    *time.Time `json:"ends_at"`
	SortOrder int        `json:"sort_order"`
	Active    bool       `json:"active"`
}

func (s *Store) cleanBanner(ctx context.Context, in BannerInput) (BannerInput, error) {
	var err error
	if in.ImageKey == "" {
		return in, FieldError{"image_key", "upload an image first"}
	}
	if err = s.checkImage(ctx, "image_key", in.ImageKey); err != nil {
		return in, err
	}
	if in.Link, err = cleanLink("link", in.Link); err != nil {
		return in, err
	}
	if in.Title, err = cleanText("title", in.Title, 80); err != nil {
		return in, err
	}
	if in.Subtitle, err = cleanText("subtitle", in.Subtitle, 140); err != nil {
		return in, err
	}
	if in.StartsAt != nil && in.EndsAt != nil && !in.EndsAt.After(*in.StartsAt) {
		return in, FieldError{"ends_at", "must be after the start"}
	}
	if in.SortOrder < -1000 || in.SortOrder > 1000 {
		return in, FieldError{"sort_order", "between -1000 and 1000"}
	}
	return in, nil
}

// SaveBanner: id kosong = baru.
func (s *Store) SaveBanner(ctx context.Context, actorID, id string, in BannerInput) (string, error) {
	in, err := s.cleanBanner(ctx, in)
	if err != nil {
		return "", err
	}
	if id != "" && !uuidRe.MatchString(id) {
		return "", ErrNotFound
	}
	err = db.WithActor(ctx, s.Pool, actorID, "admin", "", func(tx pgx.Tx) error {
		if id == "" {
			var count int
			if err := tx.QueryRow(ctx, `SELECT count(*) FROM content_banners`).Scan(&count); err != nil {
				return err
			}
			if count >= 20 {
				return FieldError{"banners", "at most 20 banners — delete old ones first"}
			}
			return tx.QueryRow(ctx, `INSERT INTO content_banners (image_key, link, title, subtitle, starts_at, ends_at, sort_order, active, created_by)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`, in.ImageKey, in.Link, in.Title, in.Subtitle, in.StartsAt, in.EndsAt,
				in.SortOrder, in.Active, actorID).Scan(&id)
		}
		tag, err := tx.Exec(ctx, `UPDATE content_banners SET image_key = $2, link = $3, title = $4, subtitle = $5, starts_at = $6, ends_at = $7,
			sort_order = $8, active = $9 WHERE id = $1`, id, in.ImageKey, in.Link, in.Title, in.Subtitle, in.StartsAt, in.EndsAt, in.SortOrder, in.Active)
		if err == nil && tag.RowsAffected() == 0 {
			return ErrNotFound
		}
		return err
	})
	s.invalidate()
	return id, err
}

func (s *Store) DeleteBanner(ctx context.Context, actorID, id string) error {
	return s.del(ctx, actorID, "content_banners", id)
}

func (s *Store) DeleteCollection(ctx context.Context, actorID, id string) error {
	return s.del(ctx, actorID, "content_collections", id)
}

func (s *Store) del(ctx context.Context, actorID, table, id string) error {
	err := db.WithActor(ctx, s.Pool, actorID, "admin", "", func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, `DELETE FROM `+table+` WHERE id::text = $1`, id)
		if err == nil && tag.RowsAffected() == 0 {
			return ErrNotFound
		}
		return err
	})
	s.invalidate()
	return err
}

type CollectionInput struct {
	Title      Text     `json:"title"`
	Source     string   `json:"source"`
	CategoryID string   `json:"category_id"`
	ProductIDs []string `json:"product_ids"`
	MaxItems   int      `json:"max_items"`
	SortOrder  int      `json:"sort_order"`
	Active     bool     `json:"active"`
}

var uuidRe = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

func (s *Store) SaveCollection(ctx context.Context, actorID, id string, in CollectionInput) (string, error) {
	var err error
	if in.Title, err = cleanText("title", in.Title, 60); err != nil {
		return "", err
	}
	if in.Title["en"] == "" && in.Title["km"] == "" && in.Title["zh"] == "" {
		return "", FieldError{"title_en", "required"}
	}
	if in.MaxItems == 0 {
		in.MaxItems = 8
	}
	if in.MaxItems < 2 || in.MaxItems > 24 {
		return "", FieldError{"max_items", "between 2 and 24"}
	}
	switch in.Source {
	case "category":
		var ok bool
		if err := s.Pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM categories WHERE id = $1 AND deleted_at IS NULL)`, in.CategoryID).Scan(&ok); err != nil {
			return "", err
		}
		if !ok {
			return "", FieldError{"category_id", "choose a category"}
		}
		in.ProductIDs = []string{}
	case "products":
		in.CategoryID = ""
		seen, ids := map[string]bool{}, []string{}
		for _, p := range in.ProductIDs {
			p = strings.ToLower(strings.TrimSpace(p))
			if !uuidRe.MatchString(p) {
				return "", FieldError{"product_ids", "invalid product"}
			}
			if !seen[p] {
				seen[p] = true
				ids = append(ids, p)
			}
		}
		if len(ids) == 0 || len(ids) > 24 {
			return "", FieldError{"product_ids", "pick 1–24 products"}
		}
		var n int
		if err := s.Pool.QueryRow(ctx, `SELECT count(*) FROM products WHERE id = ANY($1::uuid[]) AND deleted_at IS NULL`, ids).Scan(&n); err != nil {
			return "", err
		}
		if n != len(ids) {
			return "", FieldError{"product_ids", "a product was not found"}
		}
		in.ProductIDs = ids
	default:
		return "", FieldError{"source", "category or products"}
	}
	if id != "" && !uuidRe.MatchString(id) {
		return "", ErrNotFound
	}
	err = db.WithActor(ctx, s.Pool, actorID, "admin", "", func(tx pgx.Tx) error {
		if id == "" {
			var count int
			if err := tx.QueryRow(ctx, `SELECT count(*) FROM content_collections`).Scan(&count); err != nil {
				return err
			}
			if count >= 12 {
				return FieldError{"collections", "at most 12 collections"}
			}
			return tx.QueryRow(ctx, `INSERT INTO content_collections (title, source, category_id, product_ids, max_items, sort_order, active, created_by)
				VALUES ($1, $2, nullif($3, ''), $4::uuid[], $5, $6, $7, $8) RETURNING id`, in.Title, in.Source, in.CategoryID, in.ProductIDs,
				in.MaxItems, in.SortOrder, in.Active, actorID).Scan(&id)
		}
		tag, err := tx.Exec(ctx, `UPDATE content_collections SET title = $2, source = $3, category_id = nullif($4, ''), product_ids = $5::uuid[],
			max_items = $6, sort_order = $7, active = $8 WHERE id = $1`, id, in.Title, in.Source, in.CategoryID, in.ProductIDs, in.MaxItems, in.SortOrder, in.Active)
		if err == nil && tag.RowsAffected() == 0 {
			return ErrNotFound
		}
		return err
	})
	s.invalidate()
	return id, err
}

var colors = map[string]bool{"pink": true, "navy": true, "green": true, "amber": true}

func (s *Store) SaveAnnouncement(ctx context.Context, actorID string, in Announcement) error {
	var err error
	if in.Text, err = cleanText("text", in.Text, 160); err != nil {
		return err
	}
	if in.Link, err = cleanLink("link", in.Link); err != nil {
		return err
	}
	if in.Color == "" {
		in.Color = "pink"
	}
	if !colors[in.Color] {
		return FieldError{"color", "pink, navy, green or amber"}
	}
	if in.Active && in.Text["en"] == "" && in.Text["km"] == "" && in.Text["zh"] == "" {
		return FieldError{"text_en", "write the announcement first"}
	}
	if in.StartsAt != nil && in.EndsAt != nil && !in.EndsAt.After(*in.StartsAt) {
		return FieldError{"ends_at", "must be after the start"}
	}
	return s.setting(ctx, actorID, "content.announcement", in)
}

func (s *Store) SaveOGImage(ctx context.Context, actorID, key string) error {
	if key != "" {
		if err := s.checkImage(ctx, "image_key", key); err != nil {
			return err
		}
	}
	return s.setting(ctx, actorID, "content.og", map[string]string{"image_key": key})
}

func (s *Store) setting(ctx context.Context, actorID, key string, v any) error {
	raw, _ := json.Marshal(v)
	err := db.WithActor(ctx, s.Pool, actorID, "admin", "", func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `INSERT INTO platform_settings (key, value, updated_by, updated_at) VALUES ($1, $2, $3, now())
			ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_by = EXCLUDED.updated_by, updated_at = now()`, key, raw, actorID)
		return err
	})
	s.invalidate()
	return err
}

var pageRe = regexp.MustCompile(`^(home|c:[a-z0-9]+(-[a-z0-9]+)*)$`)

// SaveSEO: title/description kosong semua = hapus pengganti (kembali ke teks bawaan).
func (s *Store) SaveSEO(ctx context.Context, actorID, page string, title, desc Text) error {
	if !pageRe.MatchString(page) {
		return FieldError{"page", "unknown page"}
	}
	var err error
	if title, err = cleanText("title", title, 70); err != nil {
		return err
	}
	if desc, err = cleanText("description", desc, 170); err != nil {
		return err
	}
	if strings.HasPrefix(page, "c:") {
		var ok bool
		if err := s.Pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM categories WHERE id = $1 AND deleted_at IS NULL)`, page[2:]).Scan(&ok); err != nil {
			return err
		}
		if !ok {
			return FieldError{"page", "unknown category"}
		}
	}
	err = db.WithActor(ctx, s.Pool, actorID, "admin", "", func(tx pgx.Tx) error {
		if len(title) == 0 && len(desc) == 0 {
			_, err := tx.Exec(ctx, `DELETE FROM seo_meta WHERE page = $1`, page)
			return err
		}
		_, err := tx.Exec(ctx, `INSERT INTO seo_meta (page, title, description, updated_by, updated_at) VALUES ($1, $2, $3, $4, now())
			ON CONFLICT (page) DO UPDATE SET title = EXCLUDED.title, description = EXCLUDED.description, updated_by = EXCLUDED.updated_by, updated_at = now()`,
			page, title, desc, actorID)
		return err
	})
	s.invalidate()
	return err
}

// Fresh: snapshot langsung dari database (untuk admin; tanpa cache).
func (s *Store) Fresh(ctx context.Context) (*Snapshot, error) {
	s.invalidate()
	snap, err := s.load(ctx)
	if err != nil {
		return nil, err
	}
	s.mu.Lock()
	s.snap, s.at = snap, time.Now()
	s.mu.Unlock()
	return snap, nil
}

// SEOList: semua meta SEO (untuk admin).
func (s *Snapshot) SEOList() []SEO {
	out := []SEO{}
	for _, m := range s.SEO {
		out = append(out, m)
	}
	return out
}
