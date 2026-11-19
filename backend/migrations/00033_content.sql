-- +goose Up
-- Konten & SEO website pembeli (diatur super admin): banner beranda, pita pengumuman, koleksi pilihan, meta SEO per halaman.

CREATE TABLE content_banners (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  image_key   text NOT NULL,                       -- upload "content_image" (publik)
  link        varchar(300) NOT NULL DEFAULT '',    -- /km/c/fashion atau https://…; kosong = tanpa tautan
  title       jsonb NOT NULL DEFAULT '{}',         -- {"en": "…", "km": "…", "zh": "…"}
  subtitle    jsonb NOT NULL DEFAULT '{}',
  starts_at   timestamptz,
  ends_at     timestamptz,
  sort_order  integer NOT NULL DEFAULT 0,
  active      boolean NOT NULL DEFAULT true,
  created_by  uuid REFERENCES users(id),
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now(),
  CHECK (ends_at IS NULL OR starts_at IS NULL OR ends_at > starts_at)
);
CREATE TRIGGER content_banners_updated BEFORE UPDATE ON content_banners FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER content_banners_audit AFTER INSERT OR UPDATE OR DELETE ON content_banners FOR EACH ROW EXECUTE FUNCTION audit_row_change();

CREATE TABLE content_collections (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  title       jsonb NOT NULL DEFAULT '{}',
  source      varchar(12) NOT NULL CHECK (source IN ('category','products')),
  category_id varchar(60) REFERENCES categories(id),
  product_ids uuid[] NOT NULL DEFAULT '{}',
  max_items   smallint NOT NULL DEFAULT 8 CHECK (max_items BETWEEN 2 AND 24),
  sort_order  integer NOT NULL DEFAULT 0,
  active      boolean NOT NULL DEFAULT true,
  created_by  uuid REFERENCES users(id),
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now(),
  CHECK (source <> 'category' OR category_id IS NOT NULL)
);
CREATE TRIGGER content_collections_updated BEFORE UPDATE ON content_collections FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER content_collections_audit AFTER INSERT OR UPDATE OR DELETE ON content_collections FOR EACH ROW EXECUTE FUNCTION audit_row_change();

-- Meta SEO: page = 'home' atau 'c:<category_id>'. Kosong = teks bawaan.
CREATE TABLE seo_meta (
  page        varchar(64) PRIMARY KEY,
  title       jsonb NOT NULL DEFAULT '{}',
  description jsonb NOT NULL DEFAULT '{}',
  updated_by  uuid REFERENCES users(id),
  updated_at  timestamptz NOT NULL DEFAULT now()
);
CREATE TRIGGER seo_meta_audit AFTER INSERT OR UPDATE OR DELETE ON seo_meta FOR EACH ROW EXECUTE FUNCTION audit_row_change();

-- Gambar banner / OG (publik).
ALTER TABLE uploads DROP CONSTRAINT uploads_kind_check;
ALTER TABLE uploads ADD CONSTRAINT uploads_kind_check
  CHECK (kind IN ('store_logo','store_banner','pickup_photo','kyc_document','product_media','delivery_proof','return_photo','content_image'));

-- +goose Down
DELETE FROM uploads WHERE kind = 'content_image';
ALTER TABLE uploads DROP CONSTRAINT uploads_kind_check;
ALTER TABLE uploads ADD CONSTRAINT uploads_kind_check
  CHECK (kind IN ('store_logo','store_banner','pickup_photo','kyc_document','product_media','delivery_proof','return_photo'));
DELETE FROM platform_settings WHERE key IN ('content.announcement', 'content.og');
DROP TABLE seo_meta;
DROP TABLE content_collections;
DROP TABLE content_banners;
