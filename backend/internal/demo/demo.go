// Package demo: mengisi database dengan toko & produk contoh (phsar seed-demo).
// Toko demo: slug berawalan "demo-", pemilik nomor +85510000000x. Hapus dengan --reset.
package demo

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"phsar/internal/seller"
	"phsar/internal/storage"
)

type Seeder struct {
	Pool     *pgxpool.Pool
	Sellers  *seller.Service // TranslationOn & MediaOn sebaiknya false (teks & media sudah disiapkan di sini)
	Uploader storage.Uploader
	Files    storage.Storage
	Log      *slog.Logger
}

type Result struct {
	Stores, Products, Skipped int
	Videos                    bool
}

var langs = []string{"en", "km", "zh"}

// Seed membuat semua toko demo yang belum ada.
func (s *Seeder) Seed(ctx context.Context) (Result, error) {
	var res Result
	_, err := exec.LookPath("ffmpeg")
	res.Videos = err == nil
	rng := rand.New(rand.NewSource(42))
	for _, st := range stores {
		var exists bool
		if err := s.Pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM stores WHERE slug = $1 AND deleted_at IS NULL)`, st.Slug).Scan(&exists); err != nil {
			return res, err
		}
		if exists {
			res.Skipped++
			continue
		}
		n, err := s.seedStore(ctx, st, res.Videos, rng)
		if err != nil {
			return res, fmt.Errorf("toko %s: %w", st.Name, err)
		}
		res.Stores++
		res.Products += n
		s.Log.Info("toko demo dibuat", "toko", st.Name, "produk", n)
	}
	return res, nil
}

func (s *Seeder) upload(ctx context.Context, owner, kind string, data []byte) (string, error) {
	up, err := s.Uploader.Save(ctx, owner, kind, bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	return up.Key, nil
}

func (s *Seeder) seedStore(ctx context.Context, st dStore, videos bool, rng *rand.Rand) (int, error) {
	phone := "+855" + strings.TrimPrefix(st.Phone, "0")
	var owner string
	err := s.Pool.QueryRow(ctx, `SELECT id FROM users WHERE phone = $1 AND deleted_at IS NULL`, phone).Scan(&owner)
	if errors.Is(err, pgx.ErrNoRows) {
		err = s.Pool.QueryRow(ctx, `INSERT INTO users (name, phone, language) VALUES ($1, $2, $3::lang_code) RETURNING id`, st.Owner, phone, st.Lang).Scan(&owner)
	}
	if err != nil {
		return 0, err
	}
	store, err := s.Sellers.CreateStore(ctx, owner, seller.CreateStoreInput{Name: st.Name, Slug: st.Slug, About: st.About[st.Lang], AboutLang: st.Lang})
	if err != nil {
		return 0, err
	}
	logo, err := s.upload(ctx, owner, "store_logo", storeLogo(st.Colors))
	if err != nil {
		return 0, err
	}
	banner, err := s.upload(ctx, owner, "store_banner", storeBanner(st.Colors))
	if err != nil {
		return 0, err
	}
	if _, err := s.Sellers.UpdateStore(ctx, store.ID, owner, seller.UpdateStoreInput{LogoKey: &logo, BannerKey: &banner}); err != nil {
		return 0, err
	}
	photo, err := s.upload(ctx, owner, "pickup_photo", storeBanner(st.Colors))
	if err != nil {
		return 0, err
	}
	hours := map[string][]string{}
	for _, d := range []string{"mon", "tue", "wed", "thu", "fri", "sat", "sun"} {
		hours[d] = []string{"09:00", "20:00"}
	}
	if _, err := s.Sellers.CreatePickupPoint(ctx, store.ID, owner, seller.PickupInput{Label: st.Pickup.Label, AddressLine: st.Pickup.Address,
		Lat: st.Pickup.Lat, Lng: st.Pickup.Lng, ContactName: st.Owner, ContactPhone: st.Phone, PickupHours: hours, PhotoKey: photo, IsDefault: true}); err != nil {
		return 0, err
	}
	// Disetujui: toko aktif, titik jemput terverifikasi, teks toko 3 bahasa.
	if _, err := s.Pool.Exec(ctx, `UPDATE stores SET status = 'active', verified_at = now() WHERE id = $1`, store.ID); err != nil {
		return 0, err
	}
	if _, err := s.Pool.Exec(ctx, `UPDATE store_pickup_points SET verified_at = now() WHERE store_id = $1`, store.ID); err != nil {
		return 0, err
	}
	for _, l := range langs {
		if _, err := s.Pool.Exec(ctx, `INSERT INTO store_translations (store_id, lang, about, slug, source) VALUES ($1, $2::lang_code, $3, $4, $5)
			ON CONFLICT (store_id, lang) DO UPDATE SET about = EXCLUDED.about`, store.ID, l, st.About[l], st.Slug, map[bool]string{true: "seller", false: "machine"}[l == st.Lang]); err != nil {
			return 0, err
		}
	}
	if _, err := s.Pool.Exec(ctx, `UPDATE translation_jobs SET status = 'done', done_at = now() WHERE entity_table = 'store_translations' AND entity_id = $1 AND status = 'queued'`, store.ID); err != nil {
		return 0, err
	}
	n := 0
	for i, p := range st.Products {
		if err := s.seedProduct(ctx, store.ID, owner, st, p, i, videos, rng); err != nil {
			return n, fmt.Errorf("produk %q: %w", p.Name["en"], err)
		}
		n++
	}
	return n, nil
}

// combos: semua kombinasi nilai opsi.
func combos(opts []dOption) [][]string {
	out := [][]string{{}}
	for _, o := range opts {
		var next [][]string
		for _, c := range out {
			for _, v := range o.Values {
				next = append(next, append(append([]string{}, c...), v))
			}
		}
		out = next
	}
	return out
}

func (s *Seeder) seedProduct(ctx context.Context, storeID, owner string, st dStore, p dProduct, idx int, videos bool, rng *rand.Rand) error {
	main := productPhoto(p.Shape, p.Colors, st.Colors[1], 0)
	media := []seller.MediaInput{}
	tmp, err := os.MkdirTemp("", "phsar-demo-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)
	if videos {
		vid, err := makeVideo(ctx, tmp, main)
		if err != nil {
			return err
		}
		k, err := s.upload(ctx, owner, "product_media", vid)
		if err != nil {
			return err
		}
		media = append(media, seller.MediaInput{Key: k, Kind: "video"})
	}
	for view := 0; view < 2; view++ {
		img := main
		if view == 1 {
			img = productPhoto(p.Shape, p.Colors, st.Colors[0], 1)
		}
		k, err := s.upload(ctx, owner, "product_media", img)
		if err != nil {
			return err
		}
		media = append(media, seller.MediaInput{Key: k, Kind: "photo"})
	}
	if p.Cond == "used_good" || p.Cond == "used_fair" {
		k, err := s.upload(ctx, owner, "product_media", productPhoto(p.Shape, p.Colors, st.Colors[1], 2))
		if err != nil {
			return err
		}
		media = append(media, seller.MediaInput{Key: k, Kind: "defect_photo"})
	}
	in := seller.ProductInput{CategoryID: p.Cat, SourceLang: st.Lang, Name: p.Name[st.Lang], Description: p.Desc[st.Lang], Condition: p.Cond,
		ConditionNote: p.CondNote, WarrantyType: "none", Media: media, Submit: videos}
	if p.WarrantyType != "" {
		in.WarrantyType, in.WarrantyMonths = p.WarrantyType, p.WarrantyMo
	}
	rd := p.ReturnDays
	in.ReturnDays = &rd
	if p.UsageMonths > 0 {
		um := p.UsageMonths
		in.UsageMonths = &um
	}
	for _, o := range p.Options {
		oi := seller.OptionInput{Kind: o.Kind, Name: o.Name}
		for _, v := range o.Values {
			ov := seller.OptionValueInput{Value: v}
			if o.Kind == "color" {
				ov.SwatchHex = swatch[v]
			}
			oi.Values = append(oi.Values, ov)
		}
		in.Options = append(in.Options, oi)
	}
	for _, c := range combos(p.Options) {
		price := p.Price
		for _, v := range c {
			if x, ok := p.Extra[v]; ok {
				price = x
			}
		}
		vi := seller.VariantInput{Options: c, Price: price, StockQty: p.Stock}
		if p.Was != "" {
			pr, _ := decimal.NewFromString(price)
			was, _ := decimal.NewFromString(p.Was)
			if was.GreaterThan(pr) {
				vi.CompareAtPrice = p.Was
			}
		}
		in.Variants = append(in.Variants, vi)
	}
	prod, err := s.Sellers.CreateProduct(ctx, storeID, owner, in)
	if err != nil {
		return err
	}
	// Disetujui moderator; tanggal tayang disebar beberapa hari ke belakang supaya urutan "baru" terlihat wajar.
	ago := rng.Intn(240) + idx
	if _, err := s.Pool.Exec(ctx, `UPDATE products SET status = 'active', moderation_status = 'approved', submitted_at = coalesce(submitted_at, now()),
		published_at = now() - make_interval(hours => $2) WHERE id = $1`, prod.ID, ago); err != nil {
		return err
	}
	// Terjemahan (ditulis manual di data.go) untuk dua bahasa lain + nama opsi/nilai.
	for _, l := range langs {
		if l == st.Lang {
			continue
		}
		if _, err := s.Pool.Exec(ctx, `INSERT INTO product_translations (product_id, lang, name, description, slug, source, review_status, source_lang)
			VALUES ($1, $2::lang_code, $3, $4, $5, 'machine', 'reviewed', $6::lang_code)
			ON CONFLICT (product_id, lang) DO UPDATE SET name = EXCLUDED.name, description = EXCLUDED.description, slug = EXCLUDED.slug`,
			prod.ID, l, p.Name[l], p.Desc[l], seller.ProductSlug(p.Name[l], p.Name["en"], prod.ID), st.Lang); err != nil {
			return err
		}
	}
	if _, err := s.Pool.Exec(ctx, `UPDATE translation_jobs SET status = 'done', done_at = now() WHERE entity_table = 'product_translations' AND entity_id = $1 AND status IN ('queued','running')`, prod.ID); err != nil {
		return err
	}
	for name, t := range optionNameTr {
		if _, err := s.Pool.Exec(ctx, `UPDATE product_options SET name_i18n = name_i18n || jsonb_build_object('km', $3::text, 'zh', $4::text, 'en', $2::text)
			WHERE product_id = $1 AND name_i18n->>$5 = $2`, prod.ID, name, t["km"], t["zh"], st.Lang); err != nil {
			return err
		}
	}
	for v, t := range valueTr {
		if _, err := s.Pool.Exec(ctx, `UPDATE product_option_values ov SET value_i18n = value_i18n || jsonb_build_object('km', $3::text, 'zh', $4::text, 'en', $2::text)
			FROM product_options o WHERE o.id = ov.option_id AND o.product_id = $1 AND ov.value_i18n->>$5 = $2`, prod.ID, v, t["km"], t["zh"], st.Lang); err != nil {
			return err
		}
	}
	// Nilai lain (ukuran, model): sama di semua bahasa.
	if _, err := s.Pool.Exec(ctx, `UPDATE product_option_values ov SET value_i18n = jsonb_build_object('en', ov.value_i18n->>$2, 'km', ov.value_i18n->>$2, 'zh', ov.value_i18n->>$2)
		FROM product_options o WHERE o.id = ov.option_id AND o.product_id = $1 AND NOT (ov.value_i18n ? 'km' AND ov.value_i18n ? 'zh' AND ov.value_i18n ? 'en')`, prod.ID, st.Lang); err != nil {
		return err
	}
	// Video: sampul = foto utama, durasi 6 detik (sudah kecil, tidak perlu dikompres ulang).
	if videos {
		if _, err := s.Pool.Exec(ctx, `UPDATE product_media v SET poster_url = (SELECT coalesce(p.image_url, $2 || p.storage_key) FROM product_media p
			  WHERE p.product_id = v.product_id AND p.kind = 'photo' AND p.deleted_at IS NULL ORDER BY p.sort_order LIMIT 1),
			  duration_sec = 6, width = 720, height = 720, processing = 'ready'
			WHERE v.product_id = $1 AND v.kind = 'video'`, prod.ID, strings.TrimSuffix(s.Files.PublicURL("x"), "x")); err != nil {
			return err
		}
	}
	_, err = s.Pool.Exec(ctx, `UPDATE product_media SET processing = 'ready' WHERE product_id = $1`, prod.ID)
	return err
}

// makeVideo: video 6 detik dari foto (zoom pelan) — butuh ffmpeg.
func makeVideo(ctx context.Context, dir string, photo []byte) ([]byte, error) {
	in, out := filepath.Join(dir, "p.jpg"), filepath.Join(dir, "v.mp4")
	if err := os.WriteFile(in, photo, 0o600); err != nil {
		return nil, err
	}
	cmd := exec.CommandContext(ctx, "ffmpeg", "-hide_banner", "-loglevel", "error", "-y", "-loop", "1", "-i", in, "-t", "6",
		"-vf", "scale=1440:1440,zoompan=z='min(zoom+0.0012,1.18)':x='iw/2-(iw/zoom/2)':y='ih/2-(ih/zoom/2)':d=180:s=720x720:fps=30,format=yuv420p",
		"-c:v", "libx264", "-preset", "veryfast", "-crf", "30", "-movflags", "+faststart", out)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ffmpeg: %v: %s", err, stderr.String())
	}
	return os.ReadFile(out)
}

// Reset: sembunyikan semua toko & produk demo (soft delete). Pesanan lama tetap tersimpan.
func (s *Seeder) Reset(ctx context.Context) (int64, error) {
	if _, err := s.Pool.Exec(ctx, `UPDATE products SET deleted_at = now() WHERE deleted_at IS NULL
		AND store_id IN (SELECT id FROM stores WHERE slug LIKE 'demo-%' AND deleted_at IS NULL)`); err != nil {
		return 0, err
	}
	// Slug dibebaskan (ditambah akhiran) supaya seed berikutnya bisa memakai slug yang sama lagi.
	if _, err := s.Pool.Exec(ctx, `UPDATE store_translations t SET slug = left(t.slug, 40) || '-x' || left(replace(t.store_id::text, '-', ''), 8)
		FROM stores st WHERE st.id = t.store_id AND st.slug LIKE 'demo-%' AND st.deleted_at IS NULL`); err != nil {
		return 0, err
	}
	tag, err := s.Pool.Exec(ctx, `UPDATE stores SET deleted_at = now(), status = 'suspended',
		slug = left(slug, 40) || '-x' || left(replace(id::text, '-', ''), 8) WHERE slug LIKE 'demo-%' AND deleted_at IS NULL`)
	return tag.RowsAffected(), err
}
