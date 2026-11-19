// Package seller: onboarding toko — buat toko, KYC, titik penjemputan, kirim untuk diverifikasi.
// Dipakai app (mode seller) dan Seller Center web. Semua route di bawah /stores/{storeID}
// dibatasi ke anggota toko itu.
package seller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"phsar/internal/geo"
	"phsar/internal/platform/db"
	"phsar/internal/storage"
)

var (
	ErrNotMember     = errors.New("kamu bukan anggota toko ini")
	ErrOwnerOnly     = errors.New("hanya pemilik toko yang bisa melakukan ini")
	ErrNotFound      = errors.New("data tidak ditemukan")
	ErrValidation    = errors.New("data tidak valid")
	ErrOutsideArea   = errors.New("lokasi di luar wilayah layanan (saat ini hanya Sihanoukville)")
	ErrNotEditable   = errors.New("toko sedang diperiksa admin, tunggu hasilnya dulu")
	ErrIncomplete    = errors.New("data toko belum lengkap")
	ErrSlugTaken     = errors.New("alamat toko (slug) sudah dipakai")
	ErrTooManyStores = errors.New("satu akun hanya boleh memiliki 1 toko")
)

type validationError struct{ field, msg string }

func (e validationError) Error() string { return e.field + ": " + e.msg }
func (e validationError) Unwrap() error { return ErrValidation }
func invalid(field, msg string) error   { return validationError{field, msg} }

type Service struct {
	pool  *pgxpool.Pool
	files storage.Storage
	now   func() time.Time
	// TranslationOn: false = server tidak punya API key terjemahan (job hanya menunggu di antrean).
	TranslationOn bool
	// TranslationReady (opsional): kunci AI diatur di admin — dicek tiap kali, menggantikan TranslationOn.
	TranslationReady func(context.Context) bool
	// MediaOn: true = ffmpeg tersedia → video dikompres & dibuatkan sampul, foto diperkecil (worker media).
	MediaOn bool
}

func NewService(pool *pgxpool.Pool, files storage.Storage) *Service {
	return &Service{pool: pool, files: files, now: time.Now, TranslationOn: true}
}

// ---- Toko ----------------------------------------------------------------

type CreateStoreInput struct {
	Name      string `json:"name"`
	Slug      string `json:"slug"`       // opsional; dibuat dari nama kalau kosong
	About     string `json:"about"`      // ditulis dalam satu bahasa
	AboutLang string `json:"about_lang"` // en · km · zh; dua bahasa lain diterjemahkan AI
}

// MaxStoresPerOwner: satu akun hanya boleh memiliki satu toko (mencegah spam toko).
const MaxStoresPerOwner = 1

var slugRe = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

func Slugify(s string) string {
	var b strings.Builder
	dash := false
	for _, r := range strings.ToLower(s) {
		switch {
		case r >= 'a' && r <= 'z' || r >= '0' && r <= '9':
			b.WriteRune(r)
			dash = false
		case !dash && b.Len() > 0:
			b.WriteByte('-')
			dash = true
		}
	}
	out := strings.Trim(b.String(), "-")
	if len(out) > 50 {
		out = strings.Trim(out[:50], "-")
	}
	return out
}

func (s *Service) CreateStore(ctx context.Context, userID string, in CreateStoreInput) (Store, error) {
	in.Name = strings.TrimSpace(in.Name)
	if l := len([]rune(in.Name)); l < 2 || l > 80 {
		return Store{}, invalid("name", "2–80 karakter")
	}
	if in.AboutLang == "" {
		in.AboutLang = "km"
	}
	if !validLang(in.AboutLang) {
		return Store{}, invalid("about_lang", "en, km, atau zh")
	}
	slug := strings.TrimSpace(in.Slug)
	if slug == "" {
		slug = Slugify(in.Name)
	}
	if slug == "" { // nama dalam huruf Khmer/Mandarin → slug acak, bisa diganti nanti
		slug = "store-" + fmt.Sprint(s.now().UnixNano()%1_000_000)
	}
	if !slugRe.MatchString(slug) || len(slug) < 3 || len(slug) > 60 {
		return Store{}, invalid("slug", "3–60 huruf kecil latin, angka, dan tanda hubung")
	}

	var storeID string
	err := db.WithActor(ctx, s.pool, userID, "seller", "", func(tx pgx.Tx) error {
		var owned int
		if err := tx.QueryRow(ctx, `SELECT count(*) FROM store_members sm JOIN stores st ON st.id = sm.store_id
			WHERE sm.user_id = $1 AND sm.role = 'owner' AND st.deleted_at IS NULL`, userID).Scan(&owned); err != nil {
			return err
		}
		if owned >= MaxStoresPerOwner { // anti spam toko: 1 akun = 1 toko (staf di toko lain tetap boleh)
			return ErrTooManyStores
		}
		var taken bool
		if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM stores WHERE slug = $1 AND deleted_at IS NULL)`, slug).Scan(&taken); err != nil {
			return err
		}
		if taken {
			return ErrSlugTaken
		}
		if err := tx.QueryRow(ctx, `INSERT INTO stores (slug, name, service_area_code) VALUES ($1, $2, 'sihanoukville') RETURNING id`,
			slug, in.Name).Scan(&storeID); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `INSERT INTO store_members (store_id, user_id, role) VALUES ($1, $2, 'owner')`, storeID, userID); err != nil {
			return err
		}
		return upsertAbout(ctx, tx, storeID, slug, in.About, in.AboutLang)
	})
	if err != nil {
		return Store{}, err
	}
	return s.GetStore(ctx, storeID)
}

// upsertAbout menyimpan deskripsi dalam bahasa yang ditulis seller, lalu mengantrekan terjemahan AI
// untuk dua bahasa lainnya (diproses worker; terjemahan yang sudah diedit seller tidak ditimpa).
func upsertAbout(ctx context.Context, tx pgx.Tx, storeID, slug, about, lang string) error {
	if _, err := tx.Exec(ctx, `
		INSERT INTO store_translations (store_id, lang, about, slug, source) VALUES ($1, $2::lang_code, $3, $4, 'seller')
		ON CONFLICT (store_id, lang) DO UPDATE SET about = EXCLUDED.about, source = 'seller'`,
		storeID, lang, about, slug); err != nil {
		return err
	}
	// Bahasa lain sekarang hasil terjemahan dari teks ini (kalau seller ganti bahasa, baris lama ikut diperbarui AI).
	if _, err := tx.Exec(ctx, `UPDATE store_translations SET source = 'machine' WHERE store_id = $1 AND lang <> $2::lang_code AND source = 'seller'`,
		storeID, lang); err != nil {
		return err
	}
	if strings.TrimSpace(about) == "" {
		return nil
	}
	for _, target := range []string{"en", "km", "zh"} {
		if target == lang {
			continue
		}
		if _, err := tx.Exec(ctx, `
			INSERT INTO translation_jobs (entity_table, entity_id, source_lang, target_lang)
			SELECT 'store_translations', $1::text, $2::lang_code, $3::lang_code
			WHERE NOT EXISTS (SELECT 1 FROM translation_jobs WHERE entity_table = 'store_translations' AND entity_id = $1::text
			                  AND target_lang = $3::lang_code AND status = 'queued')`,
			storeID, lang, target); err != nil {
			return err
		}
	}
	return nil
}

type Store struct {
	ID          string            `json:"id"`
	Slug        string            `json:"slug"`
	Name        string            `json:"name"`
	Status      string            `json:"status"` // onboarding · pending_review · active · suspended · closed
	ReviewNote  string            `json:"review_note,omitempty"`
	LogoURL     string            `json:"logo_url,omitempty"`
	BannerURL   string            `json:"banner_url,omitempty"`
	SocialLinks map[string]string `json:"social_links"`
	About       map[string]string `json:"about"`
	Role        string            `json:"role,omitempty"` // owner · staff (untuk user yang meminta)
	Checklist   *Checklist        `json:"checklist,omitempty"`
	Commission  *Commission       `json:"commission,omitempty"` // hanya di GET toko (kalkulator "Anda terima")
	CreatedAt   time.Time         `json:"created_at"`
}

// Checklist: yang masih kurang sebelum toko bisa dikirim untuk diverifikasi.
type Checklist struct {
	Logo           bool   `json:"logo"`
	About          bool   `json:"about"`
	KYC            string `json:"kyc"` // missing · pending · verified · rejected
	PickupPoints   int    `json:"pickup_points"`
	PickupVerified int    `json:"pickup_points_verified"`
	CanSubmit      bool   `json:"can_submit"`
	CanSell        bool   `json:"can_sell"` // aktif + KYC verified + ≥1 titik terverifikasi
}

func (s *Service) GetStore(ctx context.Context, storeID string) (Store, error) {
	var st Store
	var note, logo, banner *string
	var social []byte
	err := s.pool.QueryRow(ctx, `
		SELECT id, slug, name, status::text, review_note, logo_key, banner_key, social_links, created_at
		FROM stores WHERE id = $1 AND deleted_at IS NULL`, storeID).
		Scan(&st.ID, &st.Slug, &st.Name, &st.Status, &note, &logo, &banner, &social, &st.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return st, ErrNotFound
	}
	if err != nil {
		return st, err
	}
	st.ReviewNote = deref(note)
	st.LogoURL = s.files.PublicURL(deref(logo))
	st.BannerURL = s.files.PublicURL(deref(banner))
	_ = json.Unmarshal(social, &st.SocialLinks)
	st.About = map[string]string{}
	rows, err := s.pool.Query(ctx, `SELECT lang::text, about FROM store_translations WHERE store_id = $1`, storeID)
	if err != nil {
		return st, err
	}
	defer rows.Close()
	for rows.Next() {
		var l, a string
		if err := rows.Scan(&l, &a); err != nil {
			return st, err
		}
		st.About[l] = a
	}
	ck, err := s.checklist(ctx, storeID, logo != nil, st.Status, st.About)
	if err != nil {
		return st, err
	}
	st.Checklist = &ck
	return st, nil
}

func (s *Service) checklist(ctx context.Context, storeID string, hasLogo bool, status string, about map[string]string) (Checklist, error) {
	ck := Checklist{Logo: hasLogo, KYC: "missing"}
	for _, a := range about {
		if strings.TrimSpace(a) != "" {
			ck.About = true
		}
	}
	var kyc *string
	if err := s.pool.QueryRow(ctx, `
		SELECT (SELECT verification_status FROM seller_tax_profiles WHERE store_id = $1),
		       (SELECT count(*) FROM store_pickup_points WHERE store_id = $1 AND deleted_at IS NULL),
		       (SELECT count(*) FROM store_pickup_points WHERE store_id = $1 AND deleted_at IS NULL AND verified_at IS NOT NULL)`,
		storeID).Scan(&kyc, &ck.PickupPoints, &ck.PickupVerified); err != nil {
		return ck, err
	}
	if kyc != nil {
		ck.KYC = *kyc
	}
	ck.CanSubmit = status == "onboarding" && ck.Logo && ck.KYC != "missing" && ck.PickupPoints > 0
	ck.CanSell = status == "active" && ck.KYC == "verified" && ck.PickupVerified > 0
	return ck, nil
}

func (s *Service) MyStores(ctx context.Context, userID string) ([]Store, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT st.id, sm.role FROM store_members sm JOIN stores st ON st.id = sm.store_id
		WHERE sm.user_id = $1 AND st.deleted_at IS NULL ORDER BY st.created_at`, userID)
	if err != nil {
		return nil, err
	}
	type pair struct{ id, role string }
	var ids []pair
	for rows.Next() {
		var p pair
		if err := rows.Scan(&p.id, &p.role); err != nil {
			rows.Close()
			return nil, err
		}
		ids = append(ids, p)
	}
	rows.Close()
	out := []Store{}
	for _, p := range ids {
		st, err := s.GetStore(ctx, p.id)
		if err != nil {
			return nil, err
		}
		st.Role = p.role
		out = append(out, st)
	}
	return out, nil
}

// Role mengembalikan peran user di toko, atau ErrNotMember.
func (s *Service) Role(ctx context.Context, storeID, userID string) (string, error) {
	var role string
	err := s.pool.QueryRow(ctx, `SELECT sm.role FROM store_members sm JOIN stores st ON st.id = sm.store_id
		WHERE sm.store_id = $1 AND sm.user_id = $2 AND st.deleted_at IS NULL`, storeID, userID).Scan(&role)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNotMember
	}
	return role, err
}

type UpdateStoreInput struct {
	Name        *string            `json:"name"`
	LogoKey     *string            `json:"logo_key"`   // dari POST /uploads kind=store_logo
	BannerKey   *string            `json:"banner_key"` // kind=store_banner
	SocialLinks *map[string]string `json:"social_links"`
	About       *string            `json:"about"`
	AboutLang   string             `json:"about_lang"`
}

func (s *Service) UpdateStore(ctx context.Context, storeID, userID string, in UpdateStoreInput) (Store, error) {
	if in.Name != nil {
		n := strings.TrimSpace(*in.Name)
		if l := len([]rune(n)); l < 2 || l > 80 {
			return Store{}, invalid("name", "2–80 karakter")
		}
		in.Name = &n
	}
	if in.LogoKey != nil {
		if err := storage.CheckOwner(ctx, s.pool, *in.LogoKey, userID, "store_logo"); err != nil {
			return Store{}, err
		}
	}
	if in.BannerKey != nil {
		if err := storage.CheckOwner(ctx, s.pool, *in.BannerKey, userID, "store_banner"); err != nil {
			return Store{}, err
		}
	}
	var social []byte
	if in.SocialLinks != nil {
		for k, v := range *in.SocialLinks {
			if k != "facebook" && k != "telegram" && k != "tiktok" && k != "instagram" {
				return Store{}, invalid("social_links", "hanya facebook, telegram, tiktok, instagram")
			}
			if len(v) > 200 {
				return Store{}, invalid("social_links", "maksimal 200 karakter")
			}
		}
		social, _ = json.Marshal(in.SocialLinks)
	}
	if in.About != nil {
		if in.AboutLang == "" || !validLang(in.AboutLang) {
			return Store{}, invalid("about_lang", "wajib en, km, atau zh saat mengubah about")
		}
		if len([]rune(*in.About)) > 2000 {
			return Store{}, invalid("about", "maksimal 2000 karakter")
		}
	}
	err := db.WithActor(ctx, s.pool, userID, "seller", "", func(tx pgx.Tx) error {
		var status, slug string
		if err := tx.QueryRow(ctx, `SELECT status::text, slug FROM stores WHERE id = $1 FOR UPDATE`, storeID).Scan(&status, &slug); err != nil {
			return err
		}
		if status == "pending_review" {
			return ErrNotEditable
		}
		if _, err := tx.Exec(ctx, `UPDATE stores SET
			name = coalesce($2, name), logo_key = coalesce($3, logo_key), banner_key = coalesce($4, banner_key),
			social_links = coalesce($5::jsonb, social_links)
			WHERE id = $1`, storeID, in.Name, emptyToNil(in.LogoKey), emptyToNil(in.BannerKey), nilIfEmptyBytes(social)); err != nil {
			return err
		}
		if in.About != nil {
			return upsertAbout(ctx, tx, storeID, slug, strings.TrimSpace(*in.About), in.AboutLang)
		}
		return nil
	})
	if err != nil {
		return Store{}, err
	}
	return s.GetStore(ctx, storeID)
}

// Submit mengirim toko untuk diverifikasi admin.
func (s *Service) Submit(ctx context.Context, storeID, userID string) (Store, error) {
	st, err := s.GetStore(ctx, storeID)
	if err != nil {
		return st, err
	}
	if !st.Checklist.CanSubmit {
		return st, ErrIncomplete
	}
	err = db.WithActor(ctx, s.pool, userID, "seller", "", func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `UPDATE stores SET status = 'pending_review', submitted_at = now(), review_note = NULL
			WHERE id = $1 AND status = 'onboarding'`, storeID)
		return err
	})
	if err != nil {
		return st, err
	}
	return s.GetStore(ctx, storeID)
}

// ---- KYC -----------------------------------------------------------------

type KYCInput struct {
	SellerType    string `json:"seller_type"` // individual_unregistered · individual_registered · company
	LegalName     string `json:"legal_name"`  // sesuai KTP/paspor/akta
	IDDocumentKey string `json:"id_document_key"`
	TIN           string `json:"tin"`
	MoCLicenseNo  string `json:"moc_license_no"`
}

type KYC struct {
	SellerType         string     `json:"seller_type"`
	LegalName          string     `json:"legal_name"`
	HasIDDocument      bool       `json:"has_id_document"`
	TIN                string     `json:"tin,omitempty"`
	MoCLicenseNo       string     `json:"moc_license_no,omitempty"`
	VerificationStatus string     `json:"verification_status"`
	RejectionNote      string     `json:"rejection_note,omitempty"`
	VerifiedAt         *time.Time `json:"verified_at,omitempty"`
}

func (s *Service) GetKYC(ctx context.Context, storeID string) (*KYC, error) {
	var k KYC
	var doc, tin, moc, note *string
	err := s.pool.QueryRow(ctx, `SELECT seller_type, legal_name, id_document_key, tin, moc_license_no, verification_status, rejection_note, verified_at
		FROM seller_tax_profiles WHERE store_id = $1`, storeID).
		Scan(&k.SellerType, &k.LegalName, &doc, &tin, &moc, &k.VerificationStatus, &note, &k.VerifiedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	k.HasIDDocument, k.TIN, k.MoCLicenseNo, k.RejectionNote = doc != nil, deref(tin), deref(moc), deref(note)
	return &k, nil
}

func (s *Service) PutKYC(ctx context.Context, storeID, userID string, in KYCInput) (*KYC, error) {
	switch in.SellerType {
	case "individual_unregistered", "individual_registered", "company":
	default:
		return nil, invalid("seller_type", "individual_unregistered, individual_registered, atau company")
	}
	in.LegalName = strings.TrimSpace(in.LegalName)
	if l := len([]rune(in.LegalName)); l < 3 || l > 120 {
		return nil, invalid("legal_name", "3–120 karakter, sesuai dokumen identitas")
	}
	if in.IDDocumentKey == "" { // tidak upload ulang → pakai dokumen yang sudah ada
		var existing *string
		_ = s.pool.QueryRow(ctx, `SELECT id_document_key FROM seller_tax_profiles WHERE store_id = $1`, storeID).Scan(&existing)
		if existing == nil || *existing == "" {
			return nil, invalid("id_document_key", "upload foto KTP/paspor dulu (kind=kyc_document)")
		}
		in.IDDocumentKey = *existing
	} else if err := storage.CheckOwner(ctx, s.pool, in.IDDocumentKey, userID, "kyc_document"); err != nil {
		return nil, err
	}
	if in.SellerType == "company" && strings.TrimSpace(in.TIN) == "" {
		return nil, invalid("tin", "wajib untuk perusahaan")
	}
	err := db.WithActor(ctx, s.pool, userID, "seller", "", func(tx pgx.Tx) error {
		var status *string
		_ = tx.QueryRow(ctx, `SELECT verification_status FROM seller_tax_profiles WHERE store_id = $1 FOR UPDATE`, storeID).Scan(&status)
		if status != nil && *status == "verified" {
			return invalid("kyc", "sudah diverifikasi; hubungi Phsar untuk mengubah data")
		}
		_, err := tx.Exec(ctx, `
			INSERT INTO seller_tax_profiles (store_id, seller_type, legal_name, id_document_key, tin, moc_license_no, verification_status)
			VALUES ($1, $2, $3, $4, $5, $6, 'pending')
			ON CONFLICT (store_id) DO UPDATE SET seller_type = EXCLUDED.seller_type, legal_name = EXCLUDED.legal_name,
			  id_document_key = EXCLUDED.id_document_key, tin = EXCLUDED.tin, moc_license_no = EXCLUDED.moc_license_no,
			  verification_status = 'pending', rejection_note = NULL`,
			storeID, in.SellerType, in.LegalName, in.IDDocumentKey, nilIfBlank(in.TIN), nilIfBlank(in.MoCLicenseNo))
		return err
	})
	if err != nil {
		return nil, err
	}
	return s.GetKYC(ctx, storeID)
}

// ---- Titik penjemputan ---------------------------------------------------

type PickupInput struct {
	Label        string              `json:"label"`
	AddressLine  string              `json:"address_line"`
	Lat          float64             `json:"lat"`
	Lng          float64             `json:"lng"`
	ContactName  string              `json:"contact_name"`
	ContactPhone string              `json:"contact_phone"`
	PickupHours  map[string][]string `json:"pickup_hours"` // {"mon":["09:00","17:00"], …}
	Note         string              `json:"note"`
	PhotoKey     string              `json:"photo_key"` // kind=pickup_photo
	IsDefault    bool                `json:"is_default"`
	PlaceID      string              `json:"google_place_id"` // dari Google Places (opsional)
}

type PickupPoint struct {
	ID            string              `json:"id"`
	Label         string              `json:"label"`
	AddressLine   string              `json:"address_line"`
	Lat           float64             `json:"lat"`
	Lng           float64             `json:"lng"`
	ContactName   string              `json:"contact_name"`
	ContactPhone  string              `json:"contact_phone"`
	PickupHours   map[string][]string `json:"pickup_hours"`
	Note          string              `json:"note,omitempty"`
	PhotoURL      string              `json:"photo_url,omitempty"`
	IsDefault     bool                `json:"is_default"`
	PlaceID       string              `json:"google_place_id,omitempty"`
	Verified      bool                `json:"verified"`
	RejectionNote string              `json:"rejection_note,omitempty"`
}

var days = map[string]bool{"mon": true, "tue": true, "wed": true, "thu": true, "fri": true, "sat": true, "sun": true}
var hhmm = regexp.MustCompile(`^([01][0-9]|2[0-3]):[0-5][0-9]$`)

func (s *Service) validatePickup(ctx context.Context, userID string, in *PickupInput) error {
	in.Label, in.AddressLine = strings.TrimSpace(in.Label), strings.TrimSpace(in.AddressLine)
	in.ContactName, in.ContactPhone = strings.TrimSpace(in.ContactName), strings.TrimSpace(in.ContactPhone)
	in.PlaceID = strings.TrimSpace(in.PlaceID)
	switch {
	case len(in.PlaceID) > 300:
		return invalid("google_place_id", "terlalu panjang")
	case in.Label == "" || len([]rune(in.Label)) > 60:
		return invalid("label", "1–60 karakter")
	case len([]rune(in.AddressLine)) < 5 || len([]rune(in.AddressLine)) > 300:
		return invalid("address_line", "5–300 karakter")
	case in.ContactName == "":
		return invalid("contact_name", "wajib")
	case len(in.ContactPhone) < 8:
		return invalid("contact_phone", "wajib")
	case len(in.PickupHours) == 0:
		return invalid("pickup_hours", "isi minimal satu hari")
	}
	for d, span := range in.PickupHours {
		if !days[d] || len(span) != 2 || !hhmm.MatchString(span[0]) || !hhmm.MatchString(span[1]) || span[0] >= span[1] {
			return invalid("pickup_hours", `format {"mon":["09:00","17:00"]}`)
		}
	}
	if err := storage.CheckOwner(ctx, s.pool, in.PhotoKey, userID, "pickup_photo"); err != nil {
		return err
	}
	return s.inServiceArea(ctx, in.Lat, in.Lng)
}

func (s *Service) inServiceArea(ctx context.Context, lat, lng float64) error {
	rows, err := s.pool.Query(ctx, `SELECT boundary FROM service_areas WHERE is_active`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var raw []byte
		if err := rows.Scan(&raw); err != nil {
			return err
		}
		if p, err := geo.ParsePolygon(raw); err == nil && p.Contains(lat, lng) {
			return nil
		}
	}
	return ErrOutsideArea
}

func (s *Service) ListPickupPoints(ctx context.Context, storeID string) ([]PickupPoint, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, label, address_line, lat::float8, lng::float8, contact_name, contact_phone, pickup_hours,
		       coalesce(note, ''), coalesce(photo_key, ''), is_default, verified_at IS NOT NULL, coalesce(rejection_note, ''),
		       coalesce(google_place_id, '')
		FROM store_pickup_points WHERE store_id = $1 AND deleted_at IS NULL ORDER BY is_default DESC, created_at`, storeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []PickupPoint{}
	for rows.Next() {
		var p PickupPoint
		var hours []byte
		var photo string
		if err := rows.Scan(&p.ID, &p.Label, &p.AddressLine, &p.Lat, &p.Lng, &p.ContactName, &p.ContactPhone, &hours,
			&p.Note, &photo, &p.IsDefault, &p.Verified, &p.RejectionNote, &p.PlaceID); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(hours, &p.PickupHours)
		p.PhotoURL = s.files.PublicURL(photo)
		out = append(out, p)
	}
	return out, nil
}

func (s *Service) CreatePickupPoint(ctx context.Context, storeID, userID string, in PickupInput) ([]PickupPoint, error) {
	if err := s.validatePickup(ctx, userID, &in); err != nil {
		return nil, err
	}
	hours, _ := json.Marshal(in.PickupHours)
	err := db.WithActor(ctx, s.pool, userID, "seller", "", func(tx pgx.Tx) error {
		var n int
		if err := tx.QueryRow(ctx, `SELECT count(*) FROM store_pickup_points WHERE store_id = $1 AND deleted_at IS NULL`, storeID).Scan(&n); err != nil {
			return err
		}
		if n >= 5 {
			return invalid("pickup_points", "maksimal 5 titik per toko")
		}
		isDefault := in.IsDefault || n == 0
		if isDefault {
			if _, err := tx.Exec(ctx, `UPDATE store_pickup_points SET is_default = false WHERE store_id = $1 AND is_default`, storeID); err != nil {
				return err
			}
		}
		_, err := tx.Exec(ctx, `
			INSERT INTO store_pickup_points (store_id, label, address_line, lat, lng, contact_name, contact_phone, pickup_hours, note, photo_key, is_default, google_place_id)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
			storeID, in.Label, in.AddressLine, in.Lat, in.Lng, in.ContactName, in.ContactPhone, hours, nilIfBlank(in.Note), nilIfBlank(in.PhotoKey), isDefault,
			nilIfBlank(in.PlaceID))
		return err
	})
	if err != nil {
		return nil, err
	}
	return s.ListPickupPoints(ctx, storeID)
}

// UpdatePickupPoint: mengubah lokasi/alamat membuat titik harus diverifikasi ulang.
func (s *Service) UpdatePickupPoint(ctx context.Context, storeID, pointID, userID string, in PickupInput) ([]PickupPoint, error) {
	if err := s.validatePickup(ctx, userID, &in); err != nil {
		return nil, err
	}
	hours, _ := json.Marshal(in.PickupHours)
	err := db.WithActor(ctx, s.pool, userID, "seller", "", func(tx pgx.Tx) error {
		if in.IsDefault {
			if _, err := tx.Exec(ctx, `UPDATE store_pickup_points SET is_default = false WHERE store_id = $1 AND is_default AND id <> $2`, storeID, pointID); err != nil {
				return err
			}
		}
		tag, err := tx.Exec(ctx, `
			UPDATE store_pickup_points SET label = $3::text, address_line = $4::text, lat = $5::numeric, lng = $6::numeric,
			  contact_name = $7::text, contact_phone = $8::text,
			  pickup_hours = $9, note = $10::text, photo_key = coalesce($11::text, photo_key), is_default = is_default OR $12::bool,
			  google_place_id = $13::text,
			  -- Lokasi/alamat berubah → harus diverifikasi admin lagi.
			  verified_at = CASE WHEN lat = $5::numeric AND lng = $6::numeric AND address_line = $4::text THEN verified_at END,
			  verified_by = CASE WHEN lat = $5::numeric AND lng = $6::numeric AND address_line = $4::text THEN verified_by END
			WHERE id = $2 AND store_id = $1 AND deleted_at IS NULL`,
			storeID, pointID, in.Label, in.AddressLine, in.Lat, in.Lng, in.ContactName, in.ContactPhone, hours,
			nilIfBlank(in.Note), nilIfBlank(in.PhotoKey), in.IsDefault, nilIfBlank(in.PlaceID))
		if err == nil && tag.RowsAffected() == 0 {
			return ErrNotFound
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return s.ListPickupPoints(ctx, storeID)
}

func (s *Service) DeletePickupPoint(ctx context.Context, storeID, pointID, userID string) ([]PickupPoint, error) {
	err := db.WithActor(ctx, s.pool, userID, "seller", "", func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, `UPDATE store_pickup_points SET deleted_at = now(), is_default = false
			WHERE id = $2 AND store_id = $1 AND deleted_at IS NULL`, storeID, pointID)
		if err == nil && tag.RowsAffected() == 0 {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		// Pastikan tetap ada satu titik default.
		_, err = tx.Exec(ctx, `UPDATE store_pickup_points SET is_default = true WHERE id = (
			SELECT id FROM store_pickup_points WHERE store_id = $1 AND deleted_at IS NULL ORDER BY created_at LIMIT 1)
			AND NOT EXISTS (SELECT 1 FROM store_pickup_points WHERE store_id = $1 AND deleted_at IS NULL AND is_default)`, storeID)
		return err
	})
	if err != nil {
		return nil, err
	}
	return s.ListPickupPoints(ctx, storeID)
}

// ---- util ----------------------------------------------------------------

func validLang(l string) bool { return l == "en" || l == "km" || l == "zh" }

func deref(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func nilIfBlank(s string) *string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return &s
}

func emptyToNil(p *string) *string {
	if p == nil || *p == "" {
		return nil
	}
	return p
}

func nilIfEmptyBytes(b []byte) *string {
	if len(b) == 0 {
		return nil
	}
	s := string(b)
	return &s
}
