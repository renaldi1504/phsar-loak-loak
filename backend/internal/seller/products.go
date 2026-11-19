package seller

// Produk seller: satu produk induk + varian (mis. warna × ukuran), harga & stok per varian,
// video + foto, kondisi baru/bekas. Seller menulis dalam satu bahasa; AI menerjemahkan
// ke dua bahasa lainnya (worker terjemahan). Produk tampil setelah dicek moderator.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/shopspring/decimal"

	"phsar/internal/modrules"
	"phsar/internal/money"
	"phsar/internal/paging"
	"phsar/internal/platform/db"
	"phsar/internal/storage"
)

var (
	ErrStoreNotActive = errors.New("toko belum aktif; produk bisa disimpan sebagai draf dulu dan dikirim setelah toko disetujui")
	ErrProductBanned  = errors.New("produk ini diblokir moderator dan tidak bisa diubah")
)

// Aturan minimal saat produk dikirim ke moderasi.
const (
	// Video tidak wajib (banyak seller hanya punya foto) tapi sangat disarankan:
	// produk tanpa video harus punya foto lebih banyak supaya pembeli tetap yakin.
	minPhotosNoVideo = 3
	maxPhotos        = 9
	maxDefects       = 6
	maxOptions       = 3
	maxValues        = 30
	maxVariants      = 100
	maxPrice         = 100000
	maxStock         = 100000
	maxDescRunes     = 5000
)

var (
	conditions = map[string]bool{"new": true, "like_new": true, "used_good": true, "used_fair": true, "refurbished": true}
	// Barang bekas dengan kondisi ini wajib punya foto cacat/bekas pemakaian.
	needsDefectPhoto = map[string]bool{"used_good": true, "used_fair": true}
	optionKinds      = map[string]bool{"color": true, "size": true, "material": true, "capacity": true, "custom": true}
	skuRe            = regexp.MustCompile(`^[A-Za-z0-9._-]{1,64}$`)
	hexRe            = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)
)

type ProductInput struct {
	CategoryID     string         `json:"category_id"`
	SourceLang     string         `json:"source_lang"` // bahasa yang diketik seller
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	Condition      string         `json:"condition"`
	ConditionNote  string         `json:"condition_note"`
	UsageMonths    *int           `json:"usage_months"`
	WarrantyType   string         `json:"warranty_type"` // none · store · official
	WarrantyMonths int            `json:"warranty_months"`
	ReturnDays     *int           `json:"return_days"`
	Options        []OptionInput  `json:"options"`
	Variants       []VariantInput `json:"variants"`
	Media          []MediaInput   `json:"media"`  // urutan = urutan tampil
	Submit         bool           `json:"submit"` // true = kirim ke moderasi; false = simpan draf
}

type OptionInput struct {
	Kind   string             `json:"kind"` // color · size · material · capacity · custom
	Name   string             `json:"name"` // "Warna", "Size"
	Values []OptionValueInput `json:"values"`
}

type OptionValueInput struct {
	Value     string `json:"value"`
	SwatchHex string `json:"swatch_hex"`
}

type VariantInput struct {
	Options        []string `json:"options"` // nilai per opsi sesuai urutan, mis. ["Merah","M"]
	SKU            string   `json:"sku"`     // kosong = dibuat otomatis
	Price          string   `json:"price"`   // USD, mis. "12.50"
	CompareAtPrice string   `json:"compare_at_price"`
	StockQty       int      `json:"stock_qty"`
	WeightGrams    *int     `json:"weight_grams"`
	Status         string   `json:"status"` // active (default) · inactive
}

type MediaInput struct {
	Key  string `json:"key"`  // dari POST /uploads kind=product_media
	Kind string `json:"kind"` // video · photo · defect_photo
}

type Product struct {
	ID                 string                        `json:"id"`
	StoreID            string                        `json:"store_id"`
	CategoryID         string                        `json:"category_id"`
	SourceLang         string                        `json:"source_lang"`
	Condition          string                        `json:"condition"`
	ConditionNote      string                        `json:"condition_note"`
	UsageMonths        *int                          `json:"usage_months"`
	WarrantyType       string                        `json:"warranty_type"`
	WarrantyMonths     int                           `json:"warranty_months"`
	ReturnDays         int                           `json:"return_days"`
	Status             string                        `json:"status"`            // draft · active · sold_out · archived · banned
	ModerationStatus   string                        `json:"moderation_status"` // pending_review · approved · changes_requested · rejected
	ModerationNote     string                        `json:"moderation_note,omitempty"`
	Visible            bool                          `json:"visible"` // tampil ke pembeli
	Translations       map[string]ProductTranslation `json:"translations"`
	TranslationPending bool                          `json:"translation_pending"`
	// TranslationState: done · pending · retrying (gagal, dicoba lagi otomatis) · failed (berhenti) · off (server tanpa API key)
	TranslationState string         `json:"translation_state"`
	TranslationError string         `json:"translation_error,omitempty"`
	Options          []Option       `json:"options"`
	Variants         []Variant      `json:"variants"`
	Media            []Media        `json:"media"`
	MinPrice         *string        `json:"min_price"`
	MaxPrice         *string        `json:"max_price"`
	TotalStock       int            `json:"total_stock"`
	AutoFlags        map[string]any `json:"auto_flags"`
	ReviewLane       *string        `json:"review_lane"`         // auto · manual (jalur moderasi terakhir)
	PostReview       bool           `json:"post_review_pending"` // tayang otomatis, belum dilihat admin
	SubmittedAt      *time.Time     `json:"submitted_at"`
	PublishedAt      *time.Time     `json:"published_at"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

type ProductTranslation struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Slug         string `json:"slug"`
	Source       string `json:"source"`        // seller · machine · reviewer
	ReviewStatus string `json:"review_status"` // machine · needs_review · reviewed
	SellerEdited bool   `json:"seller_edited"`
	Model        string `json:"model,omitempty"`
}

type Option struct {
	ID       string            `json:"id"`
	Kind     string            `json:"kind"`
	Name     map[string]string `json:"name"`
	Position int               `json:"position"`
	Values   []OptionValue     `json:"values"`
}

type OptionValue struct {
	ID        string            `json:"id"`
	Value     map[string]string `json:"value"`
	SwatchHex string            `json:"swatch_hex,omitempty"`
}

type Variant struct {
	ID             string   `json:"id"`
	SKU            string   `json:"sku"`
	Price          string   `json:"price"`
	CompareAtPrice *string  `json:"compare_at_price"`
	StockQty       int      `json:"stock_qty"`
	ReservedQty    int      `json:"reserved_qty"`
	WeightGrams    *int     `json:"weight_grams"`
	Status         string   `json:"status"`
	Options        []string `json:"options"` // nilai dalam bahasa asli seller, urut sesuai opsi
}

type Media struct {
	ID          string `json:"id"`
	Kind        string `json:"kind"`
	Key         string `json:"key"`
	URL         string `json:"url"`          // versi terbaik untuk ditampilkan (hasil kompres, atau file asli)
	OriginalURL string `json:"original_url"` // file asli dari HP
	PosterURL   string `json:"poster_url,omitempty"`
	Processing  string `json:"processing"` // processing · ready · failed
	DurationSec *int   `json:"duration_sec,omitempty"`
	Width       *int   `json:"width,omitempty"`
	Height      *int   `json:"height,omitempty"`
	SortOrder   int    `json:"sort_order"`
}

type ProductSummary struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	CategoryID       string    `json:"category_id"`
	Condition        string    `json:"condition"`
	Status           string    `json:"status"`
	ModerationStatus string    `json:"moderation_status"`
	ModerationNote   string    `json:"moderation_note,omitempty"`
	CoverURL         string    `json:"cover_url"`
	HasVideo         bool      `json:"has_video"`
	MinPrice         *string   `json:"min_price"`
	MaxPrice         *string   `json:"max_price"`
	TotalStock       int       `json:"total_stock"`
	Variants         int       `json:"variants"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// ProductSlug: slug URL per bahasa, mis. "kaos-polos-katun-3f2a9c1d". Nama Khmer/Cina tidak punya
// huruf latin, jadi memakai nama cadangan (biasanya bahasa asli / Inggris) atau "item".
func ProductSlug(name, fallback, productID string) string {
	s := Slugify(name)
	if s == "" {
		s = Slugify(fallback)
	}
	if s == "" {
		s = "item"
	}
	id := strings.ReplaceAll(productID, "-", "")
	if len(id) > 8 {
		id = id[:8]
	}
	return s + "-" + id
}

// ---- Validasi ----------------------------------------------------------------

type normVariant struct {
	key     string // nilai opsi (huruf kecil) digabung, untuk mencocokkan varian lama
	values  []string
	sku     string
	price   decimal.Decimal
	compare *decimal.Decimal
	stock   int
	weight  *int
	status  string
}

func normalizeProduct(in *ProductInput) ([]normVariant, error) {
	in.Name = strings.TrimSpace(in.Name)
	in.Description = strings.TrimSpace(in.Description)
	in.ConditionNote = strings.TrimSpace(in.ConditionNote)
	if !validLang(in.SourceLang) {
		return nil, invalid("source_lang", "en, km, atau zh")
	}
	if l := len([]rune(in.Name)); l < 3 || l > 160 {
		return nil, invalid("name", "3–160 karakter")
	}
	if len([]rune(in.Description)) > maxDescRunes {
		return nil, invalid("description", fmt.Sprintf("maksimal %d karakter", maxDescRunes))
	}
	if in.CategoryID == "" {
		return nil, invalid("category_id", "pilih kategori")
	}
	if in.Condition == "" {
		in.Condition = "new"
	}
	if !conditions[in.Condition] {
		return nil, invalid("condition", "new, like_new, used_good, used_fair, atau refurbished")
	}
	if in.Condition != "new" && in.ConditionNote == "" {
		return nil, invalid("condition_note", "jelaskan kondisi barang (lecet, kekurangan, kelengkapan)")
	}
	if in.Condition == "new" {
		in.UsageMonths = nil
	} else if in.UsageMonths != nil && (*in.UsageMonths < 0 || *in.UsageMonths > 600) {
		return nil, invalid("usage_months", "0–600")
	}
	if in.WarrantyType == "" {
		in.WarrantyType = "none"
	}
	switch in.WarrantyType {
	case "none":
		in.WarrantyMonths = 0
	case "store", "official":
		if in.WarrantyMonths < 1 || in.WarrantyMonths > 120 {
			return nil, invalid("warranty_months", "1–120 bulan")
		}
	default:
		return nil, invalid("warranty_type", "none, store, atau official")
	}
	if in.ReturnDays == nil {
		d := 7
		in.ReturnDays = &d
	}
	if *in.ReturnDays < 0 || *in.ReturnDays > 30 {
		return nil, invalid("return_days", "0–30 hari")
	}

	// Opsi & nilai
	if len(in.Options) > maxOptions {
		return nil, invalid("options", fmt.Sprintf("maksimal %d opsi", maxOptions))
	}
	valueIdx := make([]map[string]bool, len(in.Options))
	kinds := map[string]bool{}
	for i := range in.Options {
		o := &in.Options[i]
		o.Name = strings.TrimSpace(o.Name)
		if !optionKinds[o.Kind] {
			return nil, invalid(fmt.Sprintf("options[%d].kind", i), "color, size, material, capacity, atau custom")
		}
		if o.Kind != "custom" && kinds[o.Kind] {
			return nil, invalid(fmt.Sprintf("options[%d].kind", i), "jenis opsi dobel")
		}
		kinds[o.Kind] = true
		if l := len([]rune(o.Name)); l < 1 || l > 30 {
			return nil, invalid(fmt.Sprintf("options[%d].name", i), "1–30 karakter")
		}
		if len(o.Values) == 0 || len(o.Values) > maxValues {
			return nil, invalid(fmt.Sprintf("options[%d].values", i), fmt.Sprintf("1–%d nilai", maxValues))
		}
		valueIdx[i] = map[string]bool{}
		for j := range o.Values {
			v := &o.Values[j]
			v.Value = strings.TrimSpace(v.Value)
			if l := len([]rune(v.Value)); l < 1 || l > 40 {
				return nil, invalid(fmt.Sprintf("options[%d].values[%d]", i, j), "1–40 karakter")
			}
			k := strings.ToLower(v.Value)
			if valueIdx[i][k] {
				return nil, invalid(fmt.Sprintf("options[%d].values[%d]", i, j), "nilai dobel: "+v.Value)
			}
			valueIdx[i][k] = true
			if v.SwatchHex != "" && !hexRe.MatchString(v.SwatchHex) {
				return nil, invalid(fmt.Sprintf("options[%d].values[%d].swatch_hex", i, j), "format #RRGGBB")
			}
		}
	}

	// Varian
	if len(in.Variants) == 0 {
		return nil, invalid("variants", "minimal satu varian (harga & stok)")
	}
	if len(in.Options) == 0 && len(in.Variants) != 1 {
		return nil, invalid("variants", "tanpa opsi hanya boleh satu varian")
	}
	if len(in.Variants) > maxVariants {
		return nil, invalid("variants", fmt.Sprintf("maksimal %d varian", maxVariants))
	}
	out := make([]normVariant, 0, len(in.Variants))
	seen, skus := map[string]bool{}, map[string]bool{}
	limit := decimal.NewFromInt(maxPrice)
	for i, v := range in.Variants {
		f := func(name string) string { return fmt.Sprintf("variants[%d].%s", i, name) }
		if len(v.Options) != len(in.Options) {
			return nil, invalid(f("options"), "isi satu nilai untuk setiap opsi")
		}
		lower := make([]string, len(v.Options))
		vals := make([]string, len(v.Options))
		for j, val := range v.Options {
			val = strings.TrimSpace(val)
			if !valueIdx[j][strings.ToLower(val)] {
				return nil, invalid(f("options"), "nilai tidak ada di opsi: "+val)
			}
			lower[j], vals[j] = strings.ToLower(val), val
		}
		key := strings.Join(lower, "\x1f")
		if seen[key] {
			return nil, invalid(f("options"), "kombinasi varian dobel")
		}
		seen[key] = true
		price, err := decimal.NewFromString(strings.TrimSpace(v.Price))
		if err != nil || !price.IsPositive() || price.GreaterThan(limit) {
			return nil, invalid(f("price"), "harga USD lebih dari 0, mis. 12.50")
		}
		price = money.Round2(price)
		nv := normVariant{key: key, values: vals, price: price, stock: v.StockQty, weight: v.WeightGrams, status: "active"}
		if s := strings.TrimSpace(v.CompareAtPrice); s != "" {
			c, err := decimal.NewFromString(s)
			if err != nil || !money.Round2(c).GreaterThan(price) || c.GreaterThan(limit) {
				return nil, invalid(f("compare_at_price"), "harga coret harus lebih tinggi dari harga jual")
			}
			c = money.Round2(c)
			nv.compare = &c
		}
		if v.StockQty < 0 || v.StockQty > maxStock {
			return nil, invalid(f("stock_qty"), fmt.Sprintf("0–%d", maxStock))
		}
		if v.WeightGrams != nil && (*v.WeightGrams < 1 || *v.WeightGrams > 100000) {
			return nil, invalid(f("weight_grams"), "1–100000 gram")
		}
		if v.Status == "inactive" {
			nv.status = "inactive"
		}
		nv.sku = strings.TrimSpace(v.SKU)
		if nv.sku != "" {
			if !skuRe.MatchString(nv.sku) {
				return nil, invalid(f("sku"), "huruf, angka, titik, strip; maks. 64")
			}
			if skus[strings.ToLower(nv.sku)] {
				return nil, invalid(f("sku"), "SKU dobel: "+nv.sku)
			}
			skus[strings.ToLower(nv.sku)] = true
		}
		out = append(out, nv)
	}

	// Media
	counts := map[string]int{}
	keys := map[string]bool{}
	for i, m := range in.Media {
		if m.Kind != "video" && m.Kind != "photo" && m.Kind != "defect_photo" {
			return nil, invalid(fmt.Sprintf("media[%d].kind", i), "video, photo, atau defect_photo")
		}
		if m.Key == "" || keys[m.Key] {
			return nil, invalid(fmt.Sprintf("media[%d].key", i), "file kosong atau dobel")
		}
		keys[m.Key] = true
		counts[m.Kind]++
	}
	if counts["video"] > 1 {
		return nil, invalid("media", "maksimal 1 video")
	}
	if counts["photo"] > maxPhotos {
		return nil, invalid("media", fmt.Sprintf("maksimal %d foto", maxPhotos))
	}
	if counts["defect_photo"] > maxDefects {
		return nil, invalid("media", fmt.Sprintf("maksimal %d foto kondisi", maxDefects))
	}
	return out, nil
}

// checkSubmission: syarat tambahan saat dikirim ke moderasi.
func checkSubmission(in ProductInput) error {
	counts := map[string]int{}
	for _, m := range in.Media {
		counts[m.Kind]++
	}
	if counts["photo"] == 0 {
		return invalid("media", "tambahkan minimal satu foto (dipakai sebagai sampul)")
	}
	if counts["video"] == 0 && counts["photo"] < minPhotosNoVideo {
		return invalid("media", fmt.Sprintf("tanpa video, tambahkan minimal %d foto dari sudut berbeda", minPhotosNoVideo))
	}
	if needsDefectPhoto[in.Condition] && counts["defect_photo"] == 0 {
		return invalid("media", "barang bekas wajib punya foto bagian yang lecet/rusak")
	}
	return nil
}

// ---- Simpan -------------------------------------------------------------------

func (s *Service) CreateProduct(ctx context.Context, storeID, userID string, in ProductInput) (Product, error) {
	return s.saveProduct(ctx, storeID, "", userID, in)
}

func (s *Service) UpdateProduct(ctx context.Context, storeID, productID, userID string, in ProductInput) (Product, error) {
	return s.saveProduct(ctx, storeID, productID, userID, in)
}

type existingProduct struct {
	status, moderation, sourceLang, name, desc, category, condition, note string
	media                                                                 []string
	approvedMin                                                           *decimal.Decimal
	cover                                                                 string // storage key foto sampul
	reported                                                              bool   // ditahan karena laporan pembeli
}

func (s *Service) saveProduct(ctx context.Context, storeID, productID, userID string, in ProductInput) (Product, error) {
	variants, err := normalizeProduct(&in)
	if err != nil {
		return Product{}, err
	}
	// Wajib sub-kategori (kategori tanpa anak). Bekas/terbatas mengikuti induk.
	var allowsUsed, restricted, hasChildren bool
	err = s.pool.QueryRow(ctx, `
		SELECT c.allows_used AND coalesce(pc.allows_used, true), c.restricted OR coalesce(pc.restricted, false),
		       EXISTS (SELECT 1 FROM categories k WHERE k.parent_id = c.id AND k.deleted_at IS NULL)
		FROM categories c LEFT JOIN categories pc ON pc.id = c.parent_id
		WHERE c.id = $1 AND c.deleted_at IS NULL AND pc.deleted_at IS NULL`, in.CategoryID).Scan(&allowsUsed, &restricted, &hasChildren)
	if errors.Is(err, pgx.ErrNoRows) {
		return Product{}, invalid("category_id", "kategori tidak dikenal")
	}
	if err != nil {
		return Product{}, err
	}
	if hasChildren {
		return Product{}, invalid("category_id", "pilih sub-kategori")
	}
	if !allowsUsed && in.Condition != "new" {
		return Product{}, invalid("condition", "kategori ini hanya untuk barang baru")
	}

	var old *existingProduct
	if productID != "" {
		old, err = s.loadExisting(ctx, storeID, productID)
		if err != nil {
			return Product{}, err
		}
		if old.status == "banned" {
			return Product{}, ErrProductBanned
		}
		// Produk yang sudah pernah dikirim: setiap simpan = kirim ulang.
		if old.status != "draft" {
			in.Submit = true
		}
	}
	if err := s.checkMedia(ctx, userID, productID, in.Media); err != nil {
		return Product{}, err
	}

	contentChanged := old == nil || old.name != in.Name || old.desc != in.Description || old.category != in.CategoryID ||
		old.condition != in.Condition || old.note != in.ConditionNote || old.sourceLang != in.SourceLang || !sameKeys(old.media, in.Media)
	textChanged := old == nil || old.name != in.Name || old.desc != in.Description || old.sourceLang != in.SourceLang

	status, moderation := "draft", "pending_review"
	var submittedNow, publishNow, postReview bool
	var lane *string
	flags := map[string]any{}
	minPrice := minActivePrice(variants)
	if in.Submit {
		var storeStatus string
		if err := s.pool.QueryRow(ctx, `SELECT status::text FROM stores WHERE id = $1`, storeID).Scan(&storeStatus); err != nil {
			return Product{}, err
		}
		if storeStatus != "active" {
			return Product{}, ErrStoreNotActive
		}
		if err := checkSubmission(in); err != nil {
			return Product{}, err
		}
		status = "active"
		if old != nil && old.status == "archived" {
			status = "archived"
		}
		d, err := s.decideLane(ctx, storeID, old, in, variants, restricted, contentChanged)
		if err != nil {
			return Product{}, err
		}
		wasApproved := old != nil && old.moderation == "approved"
		switch {
		case wasApproved && !contentChanged && len(d.reasons) == 0:
			moderation = "approved" // hanya harga/stok/varian berubah dan harganya wajar
		case len(d.reasons) == 0:
			// 🟢 tayang otomatis. Produk baru / edit besar → masuk daftar "Baru tayang" untuk dicek admin.
			moderation, submittedNow, flags = "approved", true, d.flags
			publishNow = !wasApproved || d.major
			postReview = publishNow
			if publishNow {
				lane = ptr("auto")
			}
		default:
			// 🟡 menunggu admin
			submittedNow, flags, lane = true, d.flags, ptr("manual")
		}
	}

	flagsJSON, _ := json.Marshal(flags)
	err = db.WithActor(ctx, s.pool, userID, "seller", "", func(tx pgx.Tx) error {
		if productID == "" {
			if err := tx.QueryRow(ctx, `
				INSERT INTO products (store_id, category_id, source_lang, condition, condition_note, usage_months,
					warranty_type, warranty_months, return_days, status, moderation_status, submitted_at, auto_flags,
					published_at, robots_index, review_lane, post_review_pending, approved_min_price)
				VALUES ($1, $2, $3::lang_code, $4::product_condition, $5, $6, $7, $8, $9, $10::product_status, $11::moderation_status,
					CASE WHEN $12 THEN now() END, $13,
					CASE WHEN $14 THEN now() END, $11::moderation_status = 'approved', $15, $16, CASE WHEN $14 THEN $17::numeric END)
				RETURNING id`,
				storeID, in.CategoryID, in.SourceLang, in.Condition, nilIfEmpty(in.ConditionNote), in.UsageMonths,
				in.WarrantyType, in.WarrantyMonths, *in.ReturnDays, status, moderation, submittedNow, flagsJSON,
				publishNow, lane, postReview, minPrice).Scan(&productID); err != nil {
				return err
			}
		} else {
			if _, err := tx.Exec(ctx, `
				UPDATE products SET category_id = $2, source_lang = $3::lang_code, condition = $4::product_condition, condition_note = $5,
					usage_months = $6, warranty_type = $7, warranty_months = $8, return_days = $9,
					status = $10::product_status, moderation_status = $11::moderation_status,
					submitted_at = CASE WHEN $12 THEN now() ELSE submitted_at END,
					moderation_note = CASE WHEN $12 THEN NULL ELSE moderation_note END,
					auto_flags = CASE WHEN $12 THEN $13 ELSE auto_flags END,
					published_at = CASE WHEN $14 THEN coalesce(published_at, now()) ELSE published_at END,
					robots_index = ($11::moderation_status = 'approved'),
					review_lane = coalesce($15, review_lane),
					post_review_pending = ($16 OR (post_review_pending AND $11::moderation_status = 'approved')),
					approved_min_price = CASE WHEN $14 THEN $17::numeric ELSE approved_min_price END
				WHERE id = $1`,
				productID, in.CategoryID, in.SourceLang, in.Condition, nilIfEmpty(in.ConditionNote), in.UsageMonths,
				in.WarrantyType, in.WarrantyMonths, *in.ReturnDays, status, moderation, submittedNow, flagsJSON,
				publishNow, lane, postReview, minPrice); err != nil {
				return err
			}
		}
		if publishNow { // keputusan otomatis ikut dicatat di riwayat moderasi (reviewer kosong = sistem)
			if _, err := tx.Exec(ctx, `INSERT INTO product_moderation_reviews (product_id, decision, auto_flags)
				VALUES ($1, 'approved', $2)`, productID, flagsJSON); err != nil {
				return err
			}
		}
		// Teks bahasa asli. Kalau bahasa asli berubah, baris lama ikut jadi "terjemahan" yang akan diperbarui AI.
		if _, err := tx.Exec(ctx, `
			INSERT INTO product_translations (product_id, lang, name, description, slug, source, review_status, seller_edited)
			VALUES ($1, $2::lang_code, $3, $4, $5, 'seller', 'reviewed', false)
			ON CONFLICT (product_id, lang) DO UPDATE SET name = EXCLUDED.name, description = EXCLUDED.description,
				slug = EXCLUDED.slug, source = 'seller', review_status = 'reviewed', seller_edited = false,
				model = NULL, source_lang = NULL`,
			productID, in.SourceLang, in.Name, in.Description, ProductSlug(in.Name, "", productID)); err != nil {
			return err
		}
		if old != nil && old.sourceLang != in.SourceLang {
			if _, err := tx.Exec(ctx, `UPDATE product_translations SET source = 'machine', review_status = 'machine'
				WHERE product_id = $1 AND lang = $2::lang_code AND NOT seller_edited`, productID, old.sourceLang); err != nil {
				return err
			}
		}
		optionsChanged, err := s.writeVariants(ctx, tx, storeID, productID, userID, old, in, variants)
		if err != nil {
			return err
		}
		if err := s.writeMedia(ctx, tx, productID, in.Media); err != nil {
			return err
		}
		if textChanged || optionsChanged {
			priority := 0
			if in.Submit {
				priority = 1 // moderator menunggu terjemahan
			}
			return enqueueProductTranslation(ctx, tx, productID, in.SourceLang, priority)
		}
		return nil
	})
	if err != nil {
		var pe *pgconn.PgError
		if errors.As(err, &pe) && pe.Code == "23505" && strings.Contains(pe.ConstraintName, "sku") {
			return Product{}, invalid("variants", "SKU sudah dipakai produk lain di toko ini")
		}
		return Product{}, err
	}
	return s.GetProduct(ctx, storeID, productID)
}

func sameKeys(old []string, media []MediaInput) bool {
	if len(old) != len(media) {
		return false
	}
	for i := range old {
		if old[i] != media[i].Key {
			return false
		}
	}
	return true
}

func (s *Service) loadExisting(ctx context.Context, storeID, productID string) (*existingProduct, error) {
	e := &existingProduct{}
	err := s.pool.QueryRow(ctx, `
		SELECT p.status::text, p.moderation_status::text, p.source_lang::text, coalesce(t.name, ''), coalesce(t.description, ''),
		       p.category_id, p.condition::text, coalesce(p.condition_note, ''),
		       coalesce((SELECT array_agg(m.storage_key ORDER BY m.sort_order) FROM product_media m
		                 WHERE m.product_id = p.id AND m.deleted_at IS NULL), '{}'),
		       coalesce((SELECT m.storage_key FROM product_media m WHERE m.product_id = p.id AND m.deleted_at IS NULL AND m.kind = 'photo'
		                 ORDER BY m.sort_order LIMIT 1), ''),
		       p.approved_min_price, p.moderation_status = 'pending_review' AND coalesce(p.auto_flags->'reasons' ? 'reported', false)
		FROM products p LEFT JOIN product_translations t ON t.product_id = p.id AND t.lang = p.source_lang
		WHERE p.id = $1 AND p.store_id = $2 AND p.deleted_at IS NULL`, productID, storeID).
		Scan(&e.status, &e.moderation, &e.sourceLang, &e.name, &e.desc, &e.category, &e.condition, &e.note, &e.media,
			&e.cover, &e.approvedMin, &e.reported)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return e, err
}

// checkMedia: file harus milik user ini (atau sudah terpasang di produk ini), dan jenisnya cocok.
func (s *Service) checkMedia(ctx context.Context, userID, productID string, media []MediaInput) error {
	for i, m := range media {
		var ctype string
		err := s.pool.QueryRow(ctx, `
			SELECT u.content_type FROM uploads u
			WHERE u.key = $1 AND u.kind = 'product_media'
			  AND (u.owner_user_id = $2 OR EXISTS (
			       SELECT 1 FROM product_media pm WHERE pm.storage_key = u.key AND pm.product_id::text = $3 AND pm.deleted_at IS NULL))`,
			m.Key, userID, productID).Scan(&ctype)
		if errors.Is(err, pgx.ErrNoRows) {
			return storage.ErrNotOwner
		}
		if err != nil {
			return err
		}
		if m.Kind == "video" && !storage.IsVideo(ctype) || m.Kind != "video" && !storage.IsImage(ctype) {
			return invalid(fmt.Sprintf("media[%d]", i), "jenis file tidak cocok ("+ctype+")")
		}
	}
	return nil
}

// writeVariants mengganti opsi & varian. Varian lama dengan kombinasi yang sama tetap memakai
// ID yang sama (penting untuk keranjang/pesanan); yang hilang di-soft-delete.
func (s *Service) writeVariants(ctx context.Context, tx pgx.Tx, storeID, productID, userID string, old *existingProduct,
	in ProductInput, variants []normVariant) (optionsChanged bool, err error) {
	type oldVariant struct {
		id    string
		stock int
	}
	existing := map[string]oldVariant{}
	oldOptions := ""
	if old != nil {
		rows, err := tx.Query(ctx, `
			SELECT v.id, v.stock_qty,
			       coalesce(string_agg(lower(ov.value_i18n->>$2), E'\x1f' ORDER BY o.position), '')
			FROM product_variants v
			LEFT JOIN variant_option_values vov ON vov.variant_id = v.id
			LEFT JOIN product_option_values ov ON ov.id = vov.option_value_id
			LEFT JOIN product_options o ON o.id = vov.option_id
			WHERE v.product_id = $1 AND v.deleted_at IS NULL
			GROUP BY v.id, v.stock_qty`, productID, old.sourceLang)
		if err != nil {
			return false, err
		}
		for rows.Next() {
			var id, key string
			var stock int
			if err := rows.Scan(&id, &stock, &key); err != nil {
				rows.Close()
				return false, err
			}
			existing[key] = oldVariant{id, stock}
		}
		rows.Close()
		if err := tx.QueryRow(ctx, `
			SELECT coalesce(string_agg(o.kind || ':' || (o.name_i18n->>$2) || '=' || vals, ';' ORDER BY o.position), '')
			FROM product_options o, LATERAL (
				SELECT string_agg((v.value_i18n->>$2) || coalesce(v.swatch_hex, ''), ',' ORDER BY v.position) AS vals
				FROM product_option_values v WHERE v.option_id = o.id) x
			WHERE o.product_id = $1`, productID, old.sourceLang).Scan(&oldOptions); err != nil {
			return false, err
		}
	}
	newParts := []string{}
	for _, o := range in.Options {
		vals := []string{}
		for _, v := range o.Values {
			vals = append(vals, v.Value+v.SwatchHex)
		}
		newParts = append(newParts, o.Kind+":"+o.Name+"="+strings.Join(vals, ","))
	}
	optionsChanged = old == nil || old.sourceLang != in.SourceLang || oldOptions != strings.Join(newParts, ";")

	// Varian lama yang tidak ada lagi → soft delete. SKU sisanya diberi nama sementara supaya tukar SKU antar varian tidak bentrok.
	keep := map[string]bool{}
	for _, v := range variants {
		if ov, ok := existing[v.key]; ok {
			keep[ov.id] = true
		}
	}
	for _, ov := range existing {
		if !keep[ov.id] {
			if _, err := tx.Exec(ctx, `UPDATE product_variants SET deleted_at = now(), status = 'inactive' WHERE id = $1`, ov.id); err != nil {
				return false, err
			}
		}
	}
	if _, err := tx.Exec(ctx, `UPDATE product_variants SET sku = 'tmp-' || id::text WHERE product_id = $1 AND deleted_at IS NULL`, productID); err != nil {
		return false, err
	}

	if optionsChanged {
		if _, err := tx.Exec(ctx, `DELETE FROM product_options WHERE product_id = $1`, productID); err != nil {
			return false, err
		}
	}
	optionIDs := make([]string, len(in.Options))
	valueIDs := make([]map[string]string, len(in.Options))
	for i, o := range in.Options {
		valueIDs[i] = map[string]string{}
		if optionsChanged {
			name, _ := json.Marshal(map[string]string{in.SourceLang: o.Name})
			if err := tx.QueryRow(ctx, `INSERT INTO product_options (product_id, kind, name_i18n, position) VALUES ($1, $2, $3, $4) RETURNING id`,
				productID, o.Kind, name, i+1).Scan(&optionIDs[i]); err != nil {
				return false, err
			}
			for j, v := range o.Values {
				val, _ := json.Marshal(map[string]string{in.SourceLang: v.Value})
				var id string
				if err := tx.QueryRow(ctx, `INSERT INTO product_option_values (option_id, value_i18n, swatch_hex, position) VALUES ($1, $2, $3, $4) RETURNING id`,
					optionIDs[i], val, nilIfEmpty(v.SwatchHex), j).Scan(&id); err != nil {
					return false, err
				}
				valueIDs[i][strings.ToLower(v.Value)] = id
			}
		} else {
			// Opsi sama persis: pakai ID lama (terjemahannya tetap).
			rows, err := tx.Query(ctx, `
				SELECT o.id, v.id, lower(v.value_i18n->>$3) FROM product_options o JOIN product_option_values v ON v.option_id = o.id
				WHERE o.product_id = $1 AND o.position = $2`, productID, i+1, in.SourceLang)
			if err != nil {
				return false, err
			}
			for rows.Next() {
				var oid, vid, val string
				if err := rows.Scan(&oid, &vid, &val); err != nil {
					rows.Close()
					return false, err
				}
				optionIDs[i], valueIDs[i][val] = oid, vid
			}
			rows.Close()
		}
	}

	prefix := strings.ToUpper(strings.ReplaceAll(productID, "-", "")[:8])
	for n, v := range variants {
		sku := v.sku
		if sku == "" {
			sku = fmt.Sprintf("P%s-%d", prefix, n+1)
		}
		var compare *string
		if v.compare != nil {
			c := v.compare.StringFixed(2)
			compare = &c
		}
		var variantID string
		delta, reason := v.stock, "restock"
		if ov, ok := existing[v.key]; ok {
			variantID = ov.id
			delta, reason = v.stock-ov.stock, "adjustment"
			if _, err := tx.Exec(ctx, `
				UPDATE product_variants SET sku = $2, price = $3::numeric, compare_at_price = $4::numeric, stock_qty = $5,
					weight_grams = $6, status = $7, is_default = $8 WHERE id = $1`,
				variantID, sku, v.price.StringFixed(2), compare, v.stock, v.weight, v.status, n == 0); err != nil {
				return false, err
			}
			if _, err := tx.Exec(ctx, `DELETE FROM variant_option_values WHERE variant_id = $1`, variantID); err != nil {
				return false, err
			}
		} else {
			if err := tx.QueryRow(ctx, `
				INSERT INTO product_variants (product_id, store_id, sku, price, compare_at_price, stock_qty, weight_grams, status, is_default)
				VALUES ($1, $2, $3, $4::numeric, $5::numeric, $6, $7, $8, $9) RETURNING id`,
				productID, storeID, sku, v.price.StringFixed(2), compare, v.stock, v.weight, v.status, n == 0).Scan(&variantID); err != nil {
				return false, err
			}
		}
		for i, val := range v.values {
			if _, err := tx.Exec(ctx, `INSERT INTO variant_option_values (variant_id, option_id, option_value_id) VALUES ($1, $2, $3)`,
				variantID, optionIDs[i], valueIDs[i][strings.ToLower(val)]); err != nil {
				return false, err
			}
		}
		if delta != 0 {
			if _, err := tx.Exec(ctx, `INSERT INTO inventory_movements (variant_id, delta, reason, created_by) VALUES ($1, $2, $3, $4)`,
				variantID, delta, reason, userID); err != nil {
				return false, err
			}
		}
	}
	return optionsChanged, refreshPriceCache(ctx, tx, productID)
}

func refreshPriceCache(ctx context.Context, tx pgx.Tx, productID string) error {
	_, err := tx.Exec(ctx, `
		UPDATE products p SET min_price = x.mn, max_price = x.mx
		FROM (SELECT min(price) mn, max(price) mx FROM product_variants
		      WHERE product_id = $1 AND deleted_at IS NULL AND status = 'active') x
		WHERE p.id = $1`, productID)
	return err
}

func (s *Service) writeMedia(ctx context.Context, tx pgx.Tx, productID string, media []MediaInput) error {
	keys := make([]string, len(media))
	for i, m := range media {
		keys[i] = m.Key
	}
	if _, err := tx.Exec(ctx, `UPDATE product_media SET deleted_at = now()
		WHERE product_id = $1 AND deleted_at IS NULL AND NOT (storage_key = ANY($2))`, productID, keys); err != nil {
		return err
	}
	for i, m := range media {
		// Media yang sudah ada: hanya urutan & jenis (foto ↔ foto kondisi) yang berubah; hasil kompres tetap.
		tag, err := tx.Exec(ctx, `UPDATE product_media SET kind = $3::media_kind, sort_order = $4
			WHERE product_id = $1 AND storage_key = $2 AND deleted_at IS NULL`, productID, m.Key, m.Kind, i)
		if err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			// Baru: dengan ffmpeg → antre kompres/perkecil ("processing", file asli tampil sementara);
			// tanpa ffmpeg → file asli langsung dipakai ("ready").
			status := "ready"
			if s.MediaOn {
				status = "processing"
			}
			var mediaID string
			if err := tx.QueryRow(ctx, `INSERT INTO product_media (product_id, kind, storage_key, processing, sort_order)
				VALUES ($1, $2::media_kind, $3, $4, $5) RETURNING id`, productID, m.Kind, m.Key, status, i).Scan(&mediaID); err != nil {
				return err
			}
			if s.MediaOn {
				if _, err := tx.Exec(ctx, `INSERT INTO media_processing_jobs (media_id, step) VALUES ($1, 'transcode')`, mediaID); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func enqueueProductTranslation(ctx context.Context, tx pgx.Tx, productID, source string, priority int) error {
	for _, target := range []string{"en", "km", "zh"} {
		if target == source {
			continue
		}
		if _, err := tx.Exec(ctx, `
			INSERT INTO translation_jobs (entity_table, entity_id, source_lang, target_lang, priority)
			SELECT 'product_translations', $1::text, $2::lang_code, $3::lang_code, $4
			WHERE NOT EXISTS (SELECT 1 FROM product_translations WHERE product_id = $1::uuid AND lang = $3::lang_code AND seller_edited)
			  AND NOT EXISTS (SELECT 1 FROM translation_jobs WHERE entity_table = 'product_translations' AND entity_id = $1::text
			                  AND target_lang = $3::lang_code AND status = 'queued')`,
			productID, source, target, priority); err != nil {
			return err
		}
	}
	return nil
}

// laneDecision: hasil cek otomatis saat produk dikirim.
type laneDecision struct {
	reasons []string       // kosong = tayang otomatis
	major   bool           // produk tayang diubah besar
	flags   map[string]any // disimpan di products.auto_flags untuk moderator
}

// decideLane: 🔴 kata "block" → error; selain itu kumpulkan alasan jalur manual (paket modrules).
func (s *Service) decideLane(ctx context.Context, storeID string, old *existingProduct, in ProductInput, variants []normVariant,
	restricted, contentChanged bool) (laneDecision, error) {
	rules, err := modrules.Load(ctx, s.pool)
	if err != nil {
		return laneDecision{}, err
	}
	blocked, review, err := s.keywordHits(ctx, in)
	if err != nil {
		return laneDecision{}, err
	}
	if len(blocked) > 0 {
		return laneDecision{}, invalid("name", fmt.Sprintf("produk dengan kata “%s” tidak boleh dijual di Phsar", blocked[0]["keyword"]))
	}
	sit := modrules.Situation{Restricted: restricted, ContentChanged: contentChanged, KeywordReview: len(review) > 0}
	if old != nil {
		if old.status != "draft" {
			sit.PrevModeration = old.moderation
		}
		sit.Reported = old.reported
		sit.MajorEdit = old.moderation == "approved" && isMajorEdit(old, in)
	}
	pc := modrules.PriceCheck{Name: in.Name}
	for _, v := range variants {
		if v.status != "active" {
			continue
		}
		pc.Prices = append(pc.Prices, v.price)
		pc.Compares = append(pc.Compares, v.compare)
	}
	if old != nil && old.moderation == "approved" {
		pc.ApprovedMin = old.approvedMin
	}
	floors, err := s.priceFloors(ctx)
	if err != nil {
		return laneDecision{}, err
	}
	sit.PriceHits = modrules.CheckPrice(pc, floors, rules)
	// Fakta toko hanya dibutuhkan kalau jalur penuh dievaluasi.
	if sit.PrevModeration != "approved" || sit.MajorEdit {
		if err := s.pool.QueryRow(ctx, `
			SELECT (SELECT count(DISTINCT r.product_id) FROM product_moderation_reviews r JOIN products p ON p.id = r.product_id
			        WHERE p.store_id = $1 AND r.decision = 'approved' AND r.reviewer_id IS NOT NULL) < $2,
			       (SELECT coalesce(sum(points), 0) FROM store_strikes WHERE store_id = $1 AND revoked_at IS NULL AND expires_at > now())`,
			storeID, rules.Probation).Scan(&sit.ProbationLeft, &sit.StrikePoints); err != nil {
			return laneDecision{}, err
		}
	}
	reasons := modrules.Decide(sit, rules)
	if reasons == nil {
		reasons = []string{}
	}
	flags := map[string]any{"reasons": reasons, "lane": "auto"}
	if len(reasons) > 0 {
		flags["lane"] = "manual"
	}
	if len(review) > 0 {
		flags["keywords"] = review
	}
	if len(sit.PriceHits) > 0 {
		flags["price"] = sit.PriceHits
	}
	if in.Condition != "new" {
		flags["used_item"] = true
	}
	return laneDecision{reasons: reasons, major: sit.MajorEdit, flags: flags}, nil
}

// isMajorEdit: perubahan yang membuat produk tayang dicek ulang penuh — nama, kategori, kondisi, atau foto sampul.
func isMajorEdit(old *existingProduct, in ProductInput) bool {
	if !strings.EqualFold(strings.Join(strings.Fields(old.name), " "), strings.Join(strings.Fields(in.Name), " ")) ||
		old.category != in.CategoryID || old.condition != in.Condition {
		return true
	}
	return old.cover != coverKey(in.Media)
}

// coverKey: foto pertama (sampul).
func coverKey(media []MediaInput) string {
	for _, m := range media {
		if m.Kind == "photo" {
			return m.Key
		}
	}
	return ""
}

// keywordHits: kata terlarang di nama, deskripsi & catatan kondisi — dipisah per aksi.
func (s *Service) keywordHits(ctx context.Context, in ProductInput) (blocked, review []map[string]string, err error) {
	text := strings.ToLower(in.Name + "\n" + in.Description + "\n" + in.ConditionNote)
	rows, err := s.pool.Query(ctx, `SELECT keyword, category, action FROM prohibited_keywords ORDER BY action, keyword`)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var kw, cat, action string
		if err := rows.Scan(&kw, &cat, &action); err != nil {
			return nil, nil, err
		}
		if modrules.ContainsWord(text, strings.ToLower(kw)) {
			hit := map[string]string{"keyword": kw, "category": cat, "action": action}
			if action == "block" {
				blocked = append(blocked, hit)
			} else {
				review = append(review, hit)
			}
		}
	}
	return blocked, review, rows.Err()
}

func (s *Service) priceFloors(ctx context.Context) ([]modrules.Floor, error) {
	rows, err := s.pool.Query(ctx, `SELECT keyword, min_price FROM price_floor_rules`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []modrules.Floor
	for rows.Next() {
		var f modrules.Floor
		if err := rows.Scan(&f.Keyword, &f.MinPrice); err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

// minActivePrice: harga termurah varian aktif (nil kalau tidak ada) — teks untuk parameter numeric.
func minActivePrice(variants []normVariant) *string {
	var min *decimal.Decimal
	for i := range variants {
		v := &variants[i]
		if v.status == "active" && (min == nil || v.price.LessThan(*min)) {
			min = &v.price
		}
	}
	if min == nil {
		return nil
	}
	str := min.StringFixed(2)
	return &str
}

func ptr[T any](v T) *T { return &v }

// ---- Baca ---------------------------------------------------------------------

// GetProduct: storeID kosong = tanpa cek toko (dipakai admin).
func (s *Service) GetProduct(ctx context.Context, storeID, productID string) (Product, error) {
	p := Product{Translations: map[string]ProductTranslation{}, Options: []Option{}, Variants: []Variant{}, Media: []Media{}}
	var note *string
	var flags []byte
	var usage *int16
	var warranty, ret int16
	var storeActive bool
	err := s.pool.QueryRow(ctx, `
		SELECT p.id, p.store_id, p.category_id, p.source_lang::text, p.condition::text, coalesce(p.condition_note, ''), p.usage_months,
		       p.warranty_type, p.warranty_months, p.return_days, p.status::text, p.moderation_status::text, p.moderation_note,
		       p.min_price::text, p.max_price::text, p.auto_flags, p.review_lane, p.post_review_pending,
		       p.submitted_at, p.published_at, p.created_at, p.updated_at,
		       st.status = 'active',
		       EXISTS (SELECT 1 FROM translation_jobs j WHERE j.entity_table = 'product_translations' AND j.entity_id = p.id::text
		               AND j.status IN ('queued','running'))
		FROM products p JOIN stores st ON st.id = p.store_id
		WHERE p.id = $1 AND ($2 = '' OR p.store_id::text = $2) AND p.deleted_at IS NULL`, productID, storeID).
		Scan(&p.ID, &p.StoreID, &p.CategoryID, &p.SourceLang, &p.Condition, &p.ConditionNote, &usage,
			&p.WarrantyType, &warranty, &ret, &p.Status, &p.ModerationStatus, &note,
			&p.MinPrice, &p.MaxPrice, &flags, &p.ReviewLane, &p.PostReview, &p.SubmittedAt, &p.PublishedAt, &p.CreatedAt, &p.UpdatedAt, &storeActive, &p.TranslationPending)
	if errors.Is(err, pgx.ErrNoRows) {
		return p, ErrNotFound
	}
	if err != nil {
		return p, err
	}
	if usage != nil {
		u := int(*usage)
		p.UsageMonths = &u
	}
	p.WarrantyMonths, p.ReturnDays = int(warranty), int(ret)
	p.ModerationNote = deref(note)
	_ = json.Unmarshal(flags, &p.AutoFlags)
	p.Visible = storeActive && p.Status == "active" && p.ModerationStatus == "approved"

	p.TranslationState = "done"
	var jobStatus, jobErr string
	err = s.pool.QueryRow(ctx, `SELECT status, coalesce(error, '') FROM translation_jobs
		WHERE entity_table = 'product_translations' AND entity_id = $1 AND status IN ('queued','running','failed','done')
		ORDER BY id DESC LIMIT 1`, productID).Scan(&jobStatus, &jobErr)
	if err == nil {
		switch {
		case jobStatus == "failed":
			p.TranslationState = "failed"
		case (jobStatus == "queued" || jobStatus == "running") && jobErr != "":
			p.TranslationState = "retrying"
		case jobStatus == "queued" || jobStatus == "running":
			p.TranslationState = "pending"
		}
		if p.TranslationState != "done" {
			p.TranslationError = truncate(jobErr, 240)
		}
	}
	on := s.TranslationOn
	if s.TranslationReady != nil {
		on = s.TranslationReady(ctx)
	}
	if !on && (p.TranslationState == "pending" || p.TranslationState == "retrying") {
		p.TranslationState = "off"
	}

	rows, err := s.pool.Query(ctx, `SELECT lang::text, name, description, slug, source, review_status, seller_edited, coalesce(model, '')
		FROM product_translations WHERE product_id = $1`, productID)
	if err != nil {
		return p, err
	}
	for rows.Next() {
		var l string
		var t ProductTranslation
		if err := rows.Scan(&l, &t.Name, &t.Description, &t.Slug, &t.Source, &t.ReviewStatus, &t.SellerEdited, &t.Model); err != nil {
			rows.Close()
			return p, err
		}
		p.Translations[l] = t
	}
	rows.Close()

	rows, err = s.pool.Query(ctx, `
		SELECT o.id, o.kind, o.name_i18n, o.position, v.id, v.value_i18n, coalesce(v.swatch_hex, '')
		FROM product_options o JOIN product_option_values v ON v.option_id = o.id
		WHERE o.product_id = $1 ORDER BY o.position, v.position`, productID)
	if err != nil {
		return p, err
	}
	for rows.Next() {
		var o Option
		var v OptionValue
		var oname, vval []byte
		var pos int16
		if err := rows.Scan(&o.ID, &o.Kind, &oname, &pos, &v.ID, &vval, &v.SwatchHex); err != nil {
			rows.Close()
			return p, err
		}
		if len(p.Options) == 0 || p.Options[len(p.Options)-1].ID != o.ID {
			_ = json.Unmarshal(oname, &o.Name)
			o.Position = int(pos)
			o.Values = []OptionValue{}
			p.Options = append(p.Options, o)
		}
		_ = json.Unmarshal(vval, &v.Value)
		last := &p.Options[len(p.Options)-1]
		last.Values = append(last.Values, v)
	}
	rows.Close()

	rows, err = s.pool.Query(ctx, `
		SELECT v.id, v.sku, v.price::text, v.compare_at_price::text, v.stock_qty, v.reserved_qty, v.weight_grams, v.status,
		       coalesce(array_agg(ov.value_i18n->>$2 ORDER BY o.position) FILTER (WHERE ov.id IS NOT NULL), '{}')
		FROM product_variants v
		LEFT JOIN variant_option_values vov ON vov.variant_id = v.id
		LEFT JOIN product_option_values ov ON ov.id = vov.option_value_id
		LEFT JOIN product_options o ON o.id = vov.option_id
		WHERE v.product_id = $1 AND v.deleted_at IS NULL
		GROUP BY v.id ORDER BY v.is_default DESC, v.created_at, v.sku`, productID, p.SourceLang)
	if err != nil {
		return p, err
	}
	for rows.Next() {
		var v Variant
		if err := rows.Scan(&v.ID, &v.SKU, &v.Price, &v.CompareAtPrice, &v.StockQty, &v.ReservedQty, &v.WeightGrams, &v.Status, &v.Options); err != nil {
			rows.Close()
			return p, err
		}
		if v.Status == "active" {
			p.TotalStock += v.StockQty
		}
		p.Variants = append(p.Variants, v)
	}
	rows.Close()
	sortVariants(p)

	rows, err = s.pool.Query(ctx, `SELECT id, kind::text, storage_key, processing, sort_order,
		       coalesce(CASE WHEN kind = 'video' THEN video_url ELSE image_url END, ''), coalesce(poster_url, ''),
		       duration_sec::int, width, height
		FROM product_media WHERE product_id = $1 AND deleted_at IS NULL ORDER BY sort_order`, productID)
	if err != nil {
		return p, err
	}
	for rows.Next() {
		var m Media
		var so int16
		var best string
		if err := rows.Scan(&m.ID, &m.Kind, &m.Key, &m.Processing, &so, &best, &m.PosterURL, &m.DurationSec, &m.Width, &m.Height); err != nil {
			rows.Close()
			return p, err
		}
		m.SortOrder = int(so)
		m.OriginalURL = s.files.PublicURL(m.Key)
		m.URL = best
		if m.URL == "" {
			m.URL = m.OriginalURL
		}
		p.Media = append(p.Media, m)
	}
	rows.Close()
	return p, nil
}

// sortVariants mengurutkan varian sesuai urutan nilai opsi (Merah-S, Merah-M, Hitam-S, …).
func sortVariants(p Product) {
	order := make([]map[string]int, len(p.Options))
	for i, o := range p.Options {
		order[i] = map[string]int{}
		for j, v := range o.Values {
			order[i][v.Value[p.SourceLang]] = j
		}
	}
	sort.SliceStable(p.Variants, func(a, b int) bool {
		va, vb := p.Variants[a].Options, p.Variants[b].Options
		for i := range order {
			if i >= len(va) || i >= len(vb) {
				break
			}
			if order[i][va[i]] != order[i][vb[i]] {
				return order[i][va[i]] < order[i][vb[i]]
			}
		}
		return false
	})
}

// productFilters: tab daftar produk seller.
var productFilters = map[string]string{
	"all":      "true",
	"draft":    "p.status = 'draft'",
	"pending":  "p.status <> 'draft' AND p.moderation_status = 'pending_review'",
	"live":     "p.status = 'active' AND p.moderation_status = 'approved'",
	"changes":  "p.status <> 'draft' AND p.moderation_status IN ('changes_requested','rejected')",
	"archived": "p.status IN ('archived','banned')",
}

var productFilterOrder = []string{"all", "draft", "pending", "live", "changes", "archived"}

var productSorts = map[string]string{
	"updated":    "p.updated_at DESC, p.id",
	"name":       "lower(coalesce(t.name, '')), p.id",
	"price_asc":  "p.min_price ASC NULLS LAST, p.id",
	"price_desc": "p.max_price DESC NULLS LAST, p.id",
	"stock_low":  "total_stock ASC, p.updated_at DESC, p.id",
}

// ProductQuery: filter (tab), q (nama di bahasa apa pun / SKU), sort, halaman.
type ProductQuery struct {
	Filter, Sort string
	Page         paging.Params
}

type ProductList struct {
	Products []ProductSummary `json:"products"`
	Counts   map[string]int   `json:"counts,omitempty"` // jumlah per tab (ikut kata kunci)
	Page     *paging.Meta     `json:"page,omitempty"`
}

// productSearch: kondisi kata kunci ($2 = pola ILIKE).
const productSearch = `(EXISTS (SELECT 1 FROM product_translations x WHERE x.product_id = p.id AND x.name ILIKE $2)
	OR EXISTS (SELECT 1 FROM product_variants v WHERE v.product_id = p.id AND v.deleted_at IS NULL AND v.sku ILIKE $2))`

// ListProducts: daftar produk toko. Tanpa q.Page.On: perilaku lama (maks. 500, tanpa jumlah per tab).
func (s *Service) ListProducts(ctx context.Context, storeID string, q ProductQuery) (*ProductList, error) {
	where, ok := productFilters[q.Filter]
	if !ok {
		where = "true"
	}
	order, ok := productSorts[q.Sort]
	if !ok {
		order = productSorts["updated"]
	}
	base := `p.store_id = $1 AND p.deleted_at IS NULL`
	args := []any{storeID}
	if q.Page.Q != "" {
		base += ` AND ` + productSearch
		args = append(args, paging.Like(q.Page.Q))
	}
	limit, offset := 500, 0
	if q.Page.On {
		limit, offset = q.Page.Size, q.Page.Offset()
	}
	rows, err := s.pool.Query(ctx, `
		SELECT p.id, coalesce(t.name, ''), p.category_id, p.condition::text, p.status::text, p.moderation_status::text,
		       coalesce(p.moderation_note, ''),
		       coalesce((SELECT coalesce(m.image_url, m.storage_key) FROM product_media m WHERE m.product_id = p.id AND m.deleted_at IS NULL AND m.kind = 'photo'
		                 ORDER BY m.sort_order LIMIT 1), ''),
		       EXISTS (SELECT 1 FROM product_media m WHERE m.product_id = p.id AND m.deleted_at IS NULL AND m.kind = 'video'),
		       p.min_price::text, p.max_price::text,
		       coalesce((SELECT sum(stock_qty) FROM product_variants v WHERE v.product_id = p.id AND v.deleted_at IS NULL AND v.status = 'active'), 0) AS total_stock,
		       (SELECT count(*) FROM product_variants v WHERE v.product_id = p.id AND v.deleted_at IS NULL),
		       p.updated_at
		FROM products p LEFT JOIN product_translations t ON t.product_id = p.id AND t.lang = p.source_lang
		WHERE `+base+` AND `+where+`
		ORDER BY `+order+` LIMIT `+fmt.Sprint(limit)+` OFFSET `+fmt.Sprint(offset), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := &ProductList{Products: []ProductSummary{}}
	for rows.Next() {
		var p ProductSummary
		var cover string
		if err := rows.Scan(&p.ID, &p.Name, &p.CategoryID, &p.Condition, &p.Status, &p.ModerationStatus, &p.ModerationNote,
			&cover, &p.HasVideo, &p.MinPrice, &p.MaxPrice, &p.TotalStock, &p.Variants, &p.UpdatedAt); err != nil {
			return nil, err
		}
		p.CoverURL = cover
		if !strings.HasPrefix(cover, "http") {
			p.CoverURL = s.files.PublicURL(cover)
		}
		out.Products = append(out.Products, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	rows.Close()
	if !q.Page.On {
		return out, nil
	}
	// Jumlah per tab dalam satu query (ikut kata kunci, supaya seller tahu di tab mana hasilnya).
	cols := make([]string, len(productFilterOrder))
	for i, f := range productFilterOrder {
		cols[i] = `count(*) FILTER (WHERE ` + productFilters[f] + `)`
	}
	vals := make([]int, len(productFilterOrder))
	ptrs := make([]any, len(vals))
	for i := range vals {
		ptrs[i] = &vals[i]
	}
	if err := s.pool.QueryRow(ctx, `SELECT `+strings.Join(cols, ", ")+` FROM products p WHERE `+base, args...).Scan(ptrs...); err != nil {
		return nil, err
	}
	out.Counts = map[string]int{}
	for i, f := range productFilterOrder {
		out.Counts[f] = vals[i]
	}
	filter := q.Filter
	if _, ok := productFilters[filter]; !ok {
		filter = "all"
	}
	m := paging.NewMeta(q.Page, out.Counts[filter])
	out.Page = &m
	return out, nil
}

// ---- Ubah cepat ---------------------------------------------------------------

type StockPriceInput struct {
	Variants []struct {
		ID             string `json:"id"`
		Price          string `json:"price"`
		CompareAtPrice string `json:"compare_at_price"`
		StockQty       int    `json:"stock_qty"`
		Status         string `json:"status"`
	} `json:"variants"`
}

// UpdateStockPrice: ubah harga/stok tanpa dicek ulang moderator.
func (s *Service) UpdateStockPrice(ctx context.Context, storeID, productID, userID string, in StockPriceInput) (Product, error) {
	if _, err := s.loadExisting(ctx, storeID, productID); err != nil {
		return Product{}, err
	}
	limit := decimal.NewFromInt(maxPrice)
	err := db.WithActor(ctx, s.pool, userID, "seller", "", func(tx pgx.Tx) error {
		for i, v := range in.Variants {
			price, err := decimal.NewFromString(strings.TrimSpace(v.Price))
			if err != nil || !price.IsPositive() || price.GreaterThan(limit) {
				return invalid(fmt.Sprintf("variants[%d].price", i), "harga USD lebih dari 0")
			}
			price = money.Round2(price)
			var compare *string
			if c := strings.TrimSpace(v.CompareAtPrice); c != "" {
				cd, err := decimal.NewFromString(c)
				if err != nil || !money.Round2(cd).GreaterThan(price) {
					return invalid(fmt.Sprintf("variants[%d].compare_at_price", i), "harus lebih tinggi dari harga jual")
				}
				cs := money.Round2(cd).StringFixed(2)
				compare = &cs
			}
			if v.StockQty < 0 || v.StockQty > maxStock {
				return invalid(fmt.Sprintf("variants[%d].stock_qty", i), fmt.Sprintf("0–%d", maxStock))
			}
			status := "active"
			if v.Status == "inactive" {
				status = "inactive"
			}
			var oldStock, reserved int
			err = tx.QueryRow(ctx, `SELECT stock_qty, reserved_qty FROM product_variants WHERE id = $1 AND product_id = $2 AND deleted_at IS NULL FOR UPDATE`,
				v.ID, productID).Scan(&oldStock, &reserved)
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrNotFound
			}
			if err != nil {
				return err
			}
			if v.StockQty < reserved {
				return invalid(fmt.Sprintf("variants[%d].stock_qty", i), fmt.Sprintf("minimal %d (sedang dipesan)", reserved))
			}
			if _, err := tx.Exec(ctx, `UPDATE product_variants SET price = $2::numeric, compare_at_price = $3::numeric, stock_qty = $4, status = $5 WHERE id = $1`,
				v.ID, price.StringFixed(2), compare, v.StockQty, status); err != nil {
				return err
			}
			if d := v.StockQty - oldStock; d != 0 {
				if _, err := tx.Exec(ctx, `INSERT INTO inventory_movements (variant_id, delta, reason, created_by) VALUES ($1, $2, 'adjustment', $3)`,
					v.ID, d, userID); err != nil {
					return err
				}
			}
		}
		return refreshPriceCache(ctx, tx, productID)
	})
	if err != nil {
		return Product{}, err
	}
	return s.GetProduct(ctx, storeID, productID)
}

// SetProductStatus: archived (sembunyikan) atau active (tampilkan lagi).
func (s *Service) SetProductStatus(ctx context.Context, storeID, productID, userID, status string) (Product, error) {
	old, err := s.loadExisting(ctx, storeID, productID)
	if err != nil {
		return Product{}, err
	}
	if old.status == "banned" {
		return Product{}, ErrProductBanned
	}
	switch {
	case status == "archived" && old.status != "draft":
	case status == "active" && old.status == "archived":
	default:
		return Product{}, invalid("status", "hanya bisa arsipkan produk yang sudah dikirim, atau aktifkan lagi produk yang diarsipkan")
	}
	err = db.WithActor(ctx, s.pool, userID, "seller", "", func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `UPDATE products SET status = $2::product_status WHERE id = $1`, productID, status)
		return err
	})
	if err != nil {
		return Product{}, err
	}
	return s.GetProduct(ctx, storeID, productID)
}

func (s *Service) DeleteProduct(ctx context.Context, storeID, productID, userID string) error {
	if _, err := s.loadExisting(ctx, storeID, productID); err != nil {
		return err
	}
	return db.WithActor(ctx, s.pool, userID, "seller", "", func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `UPDATE products SET deleted_at = now(), deleted_by = $2, status = 'archived' WHERE id = $1`, productID, userID); err != nil {
			return err
		}
		_, err := tx.Exec(ctx, `UPDATE translation_jobs SET status = 'skipped', done_at = now()
			WHERE entity_table = 'product_translations' AND entity_id = $1 AND status = 'queued'`, productID)
		return err
	})
}

// EditTranslation: seller memperbaiki terjemahan satu bahasa sendiri → tidak ditimpa AI lagi.
func (s *Service) EditTranslation(ctx context.Context, storeID, productID, userID, lang string, name, desc string) (Product, error) {
	old, err := s.loadExisting(ctx, storeID, productID)
	if err != nil {
		return Product{}, err
	}
	if !validLang(lang) || lang == old.sourceLang {
		return Product{}, invalid("lang", "ubah teks bahasa asli lewat form produk")
	}
	name, desc = strings.TrimSpace(name), strings.TrimSpace(desc)
	if l := len([]rune(name)); l < 1 || l > 160 {
		return Product{}, invalid("name", "1–160 karakter")
	}
	if len([]rune(desc)) > maxDescRunes {
		return Product{}, invalid("description", fmt.Sprintf("maksimal %d karakter", maxDescRunes))
	}
	err = db.WithActor(ctx, s.pool, userID, "seller", "", func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `
			INSERT INTO product_translations (product_id, lang, name, description, slug, source, review_status, seller_edited, source_lang)
			VALUES ($1, $2::lang_code, $3, $4, $5, 'seller', 'reviewed', true, $6::lang_code)
			ON CONFLICT (product_id, lang) DO UPDATE SET name = EXCLUDED.name, description = EXCLUDED.description,
				slug = EXCLUDED.slug, source = 'seller', review_status = 'reviewed', seller_edited = true, model = NULL`,
			productID, lang, name, desc, ProductSlug(name, old.name, productID), old.sourceLang)
		return err
	})
	if err != nil {
		return Product{}, err
	}
	return s.GetProduct(ctx, storeID, productID)
}

// ResetTranslation: kembalikan satu bahasa ke terjemahan AI.
func (s *Service) ResetTranslation(ctx context.Context, storeID, productID, userID, lang string) (Product, error) {
	old, err := s.loadExisting(ctx, storeID, productID)
	if err != nil {
		return Product{}, err
	}
	if !validLang(lang) || lang == old.sourceLang {
		return Product{}, invalid("lang", "bukan bahasa terjemahan")
	}
	err = db.WithActor(ctx, s.pool, userID, "seller", "", func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `UPDATE product_translations SET seller_edited = false WHERE product_id = $1 AND lang = $2::lang_code`,
			productID, lang); err != nil {
			return err
		}
		_, err := tx.Exec(ctx, `
			INSERT INTO translation_jobs (entity_table, entity_id, source_lang, target_lang, priority)
			SELECT 'product_translations', $1::text, $2::lang_code, $3::lang_code, 1
			WHERE NOT EXISTS (SELECT 1 FROM translation_jobs WHERE entity_table = 'product_translations' AND entity_id = $1::text
			                  AND target_lang = $3::lang_code AND status = 'queued')`, productID, old.sourceLang, lang)
		return err
	})
	if err != nil {
		return Product{}, err
	}
	return s.GetProduct(ctx, storeID, productID)
}

// RetryTranslation: coba terjemahkan lagi sekarang (job gagal / sedang menunggu jeda dijalankan ulang).
func (s *Service) RetryTranslation(ctx context.Context, storeID, productID, userID string) (Product, error) {
	old, err := s.loadExisting(ctx, storeID, productID)
	if err != nil {
		return Product{}, err
	}
	err = db.WithActor(ctx, s.pool, userID, "seller", "", func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `UPDATE translation_jobs SET status = 'queued', attempts = 0, run_after = NULL, error = NULL
			WHERE entity_table = 'product_translations' AND entity_id = $1 AND status IN ('failed','queued')`, productID); err != nil {
			return err
		}
		return enqueueProductTranslation(ctx, tx, productID, old.sourceLang, 1)
	})
	if err != nil {
		return Product{}, err
	}
	return s.GetProduct(ctx, storeID, productID)
}

func truncate(s string, n int) string {
	if r := []rune(s); len(r) > n {
		return string(r[:n]) + "…"
	}
	return s
}

func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
