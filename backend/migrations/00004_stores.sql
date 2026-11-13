-- +goose Up
-- Toko, anggota toko, wilayah layanan, dan titik penjemputan.
-- Launch: hanya Sihanoukville, ongkir gratis, kurir Phsar menjemput di titik toko.

CREATE TABLE shipping_zones (
  code          varchar(32) PRIMARY KEY,
  fee_per_store numeric(12,2) NOT NULL DEFAULT 0.00 CHECK (fee_per_store >= 0),  -- 0.00 saat launch
  is_estimate   boolean NOT NULL DEFAULT true,
  updated_at    timestamptz NOT NULL DEFAULT now()
);
CREATE TRIGGER shipping_zones_audit AFTER INSERT OR UPDATE OR DELETE ON shipping_zones FOR EACH ROW EXECUTE FUNCTION audit_row_change();

CREATE TABLE service_areas (
  code               varchar(32) PRIMARY KEY,
  name               jsonb NOT NULL,                -- {"en":…,"km":…,"zh":…}
  boundary           jsonb NOT NULL,                -- GeoJSON Polygon [lng, lat]
  shipping_zone_code varchar(32) NOT NULL REFERENCES shipping_zones(code),
  is_active          boolean NOT NULL DEFAULT false,
  updated_at         timestamptz NOT NULL DEFAULT now()
);
CREATE TRIGGER service_areas_audit AFTER INSERT OR UPDATE OR DELETE ON service_areas FOR EACH ROW EXECUTE FUNCTION audit_row_change();

CREATE TYPE store_status AS ENUM ('onboarding', 'pending_review', 'active', 'suspended', 'closed');

CREATE TABLE stores (
  id                uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  slug              varchar(60) NOT NULL CHECK (slug ~ '^[a-z0-9]+(-[a-z0-9]+)*$'),
  name              varchar(80) NOT NULL,
  logo_url          text,
  banner_url        text,
  service_area_code varchar(32) REFERENCES service_areas(code),
  social_links      jsonb NOT NULL DEFAULT '{}',    -- {"facebook":…,"telegram":…,"tiktok":…}
  status            store_status NOT NULL DEFAULT 'onboarding',
  verified_at       timestamptz,
  verified_by       uuid REFERENCES users(id),
  created_at        timestamptz NOT NULL DEFAULT now(),
  updated_at        timestamptz NOT NULL DEFAULT now(),
  deleted_at        timestamptz,
  deleted_by        uuid REFERENCES users(id)
);
CREATE UNIQUE INDEX stores_slug_active ON stores (slug) WHERE deleted_at IS NULL;
CREATE TRIGGER stores_updated BEFORE UPDATE ON stores FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER stores_audit AFTER INSERT OR UPDATE OR DELETE ON stores FOR EACH ROW EXECUTE FUNCTION audit_row_change();

CREATE TABLE store_members (
  store_id   uuid NOT NULL REFERENCES stores(id),
  user_id    uuid NOT NULL REFERENCES users(id),
  role       varchar(16) NOT NULL CHECK (role IN ('owner','staff')),
  created_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (store_id, user_id)
);
CREATE INDEX store_members_user_idx ON store_members (user_id);  -- "user ini punya toko?" → tampilkan mode seller
CREATE UNIQUE INDEX store_one_owner ON store_members (store_id) WHERE role = 'owner';
CREATE TRIGGER store_members_audit AFTER INSERT OR UPDATE OR DELETE ON store_members FOR EACH ROW EXECUTE FUNCTION audit_row_change();

CREATE TABLE store_translations (
  store_id         uuid NOT NULL REFERENCES stores(id),
  lang             lang_code NOT NULL,
  about            text NOT NULL DEFAULT '',
  slug             varchar(60) NOT NULL CHECK (slug ~ '^[a-z0-9]+(-[a-z0-9]+)*$'),
  meta_title       varchar(70),
  meta_description varchar(160),
  source           varchar(16) NOT NULL DEFAULT 'seller' CHECK (source IN ('seller','machine','reviewer')),
  updated_at       timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (store_id, lang),
  UNIQUE (lang, slug)
);
CREATE TRIGGER store_translations_updated BEFORE UPDATE ON store_translations FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Titik penjemputan: wajib dan terverifikasi sebelum produk toko tampil.
CREATE TABLE store_pickup_points (
  id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  store_id       uuid NOT NULL REFERENCES stores(id),
  label          varchar(60) NOT NULL,
  address_line   varchar(300) NOT NULL,
  lat            numeric(9,6) NOT NULL CHECK (lat BETWEEN -90 AND 90),
  lng            numeric(9,6) NOT NULL CHECK (lng BETWEEN -180 AND 180),
  contact_name   varchar(80) NOT NULL,
  contact_phone  varchar(30) NOT NULL,
  pickup_hours   jsonb NOT NULL DEFAULT '{}',      -- {"mon":["09:00","17:00"], …}
  note           text,                             -- patokan lokasi
  photo_url      text,                             -- foto depan toko
  is_default     boolean NOT NULL DEFAULT false,
  verified_at    timestamptz,
  verified_by    uuid REFERENCES users(id),
  created_at     timestamptz NOT NULL DEFAULT now(),
  updated_at     timestamptz NOT NULL DEFAULT now(),
  deleted_at     timestamptz
);
CREATE INDEX store_pickup_points_store_idx ON store_pickup_points (store_id) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX store_pickup_one_default ON store_pickup_points (store_id) WHERE is_default AND deleted_at IS NULL;
CREATE TRIGGER store_pickup_points_updated BEFORE UPDATE ON store_pickup_points FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER store_pickup_points_audit AFTER INSERT OR UPDATE OR DELETE ON store_pickup_points FOR EACH ROW EXECUTE FUNCTION audit_row_change();

-- Data pajak & KYC seller.
CREATE TABLE seller_tax_profiles (
  store_id             uuid PRIMARY KEY REFERENCES stores(id),
  seller_type          varchar(32) NOT NULL DEFAULT 'individual_unregistered'
                       CHECK (seller_type IN ('individual_unregistered','individual_registered','company')),
  taxpayer_size        varchar(16) NOT NULL DEFAULT 'none' CHECK (taxpayer_size IN ('none','small','medium','large')),
  tin                  varchar(32),
  vattin               varchar(32),
  vat_registered       boolean NOT NULL DEFAULT false,
  gmv_12m              numeric(14,2) NOT NULL DEFAULT 0,
  gmv_limit            numeric(14,2),
  limit_warned_at      timestamptz,
  limit_reached_at     timestamptz,
  moc_license_no       varchar(64),
  license_expires_at   date,
  legal_name           varchar(120) NOT NULL,
  id_document_key      text,                       -- object key di bucket privat
  verification_status  varchar(16) NOT NULL DEFAULT 'pending' CHECK (verification_status IN ('pending','verified','rejected')),
  rejection_note       text,
  verified_by          uuid REFERENCES users(id),
  verified_at          timestamptz,
  updated_at           timestamptz NOT NULL DEFAULT now()
);
CREATE TRIGGER seller_tax_profiles_updated BEFORE UPDATE ON seller_tax_profiles FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER seller_tax_profiles_audit AFTER INSERT OR UPDATE OR DELETE ON seller_tax_profiles FOR EACH ROW EXECUTE FUNCTION audit_row_change();

CREATE TABLE seller_payout_accounts (
  id                      uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  store_id                uuid NOT NULL REFERENCES stores(id),
  bank                    varchar(16) NOT NULL CHECK (bank IN ('aba','acleda','wing','bakong','other')),
  account_name            varchar(120) NOT NULL,
  account_number_enc      bytea NOT NULL,          -- pgp_sym_encrypt; tidak pernah masuk log/audit
  account_last4           char(4) NOT NULL,
  provider_beneficiary_id varchar(64),
  currency                char(3) NOT NULL DEFAULT 'USD',
  is_default              boolean NOT NULL DEFAULT false,
  verified_at             timestamptz,
  payout_hold_until       timestamptz,             -- ganti rekening = tahan payout 48 jam
  created_at              timestamptz NOT NULL DEFAULT now(),
  deleted_at              timestamptz
);
CREATE INDEX seller_payout_accounts_store_idx ON seller_payout_accounts (store_id) WHERE deleted_at IS NULL;
-- Catatan: sengaja tanpa trigger audit_row_change supaya nomor rekening terenkripsi tidak tersalin ke audit_logs;
-- perubahan rekening dicatat sebagai aksi bisnis "payout_account.change" dari kode Go.

-- +goose Down
DROP TABLE seller_payout_accounts;
DROP TABLE seller_tax_profiles;
DROP TABLE store_pickup_points;
DROP TABLE store_translations;
DROP TABLE store_members;
DROP TABLE stores;
DROP TYPE store_status;
DROP TABLE service_areas;
DROP TABLE shipping_zones;
