-- +goose Up
-- Katalog: kategori, merek, produk induk, varian (warna × ukuran), stok, media.
-- Harga NUMERIC(12,2) USD; pembulatan biasa dilakukan di Go (money.Round2).

CREATE TABLE categories (
  id          varchar(60) PRIMARY KEY CHECK (id ~ '^[a-z0-9]+(-[a-z0-9]+)*$'),
  parent_id   varchar(60) REFERENCES categories(id),
  icon        varchar(40),
  allows_used boolean NOT NULL DEFAULT true,       -- false untuk kosmetik, makanan
  sort_order  integer NOT NULL DEFAULT 0,
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now(),
  deleted_at  timestamptz
);
CREATE TRIGGER categories_updated BEFORE UPDATE ON categories FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER categories_audit AFTER INSERT OR UPDATE OR DELETE ON categories FOR EACH ROW EXECUTE FUNCTION audit_row_change();

CREATE TABLE category_translations (
  category_id      varchar(60) NOT NULL REFERENCES categories(id),
  lang             lang_code NOT NULL,
  name             varchar(80) NOT NULL,
  slug             varchar(80) NOT NULL CHECK (slug ~ '^[a-z0-9]+(-[a-z0-9]+)*$'),
  path             varchar(255) NOT NULL,          -- fashion/men/t-shirts, untuk URL & breadcrumb
  intro            text,
  meta_title       varchar(70),
  meta_description varchar(160),
  PRIMARY KEY (category_id, lang),
  UNIQUE (lang, slug)
);

CREATE TABLE brands (
  id         integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  slug       varchar(60) NOT NULL UNIQUE CHECK (slug ~ '^[a-z0-9]+(-[a-z0-9]+)*$'),
  name       varchar(80) NOT NULL,
  logo_url   text,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE size_charts (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  store_id    uuid NOT NULL REFERENCES stores(id),
  name        varchar(80) NOT NULL,
  size_system varchar(8) NOT NULL CHECK (size_system IN ('letter','EU','US','UK','CM')),
  rows        jsonb NOT NULL DEFAULT '[]',          -- [{"size":"M","chest_cm":100}]
  created_at  timestamptz NOT NULL DEFAULT now(),
  deleted_at  timestamptz
);

CREATE TABLE category_attributes (
  id            integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  category_id   varchar(60) NOT NULL REFERENCES categories(id),
  code          varchar(40) NOT NULL,              -- material · ram_gb · storage_gb
  label_i18n    jsonb NOT NULL,
  data_type     varchar(8) NOT NULL CHECK (data_type IN ('text','number','bool','select')),
  choices       jsonb,
  is_required   boolean NOT NULL DEFAULT false,
  is_filterable boolean NOT NULL DEFAULT false,
  UNIQUE (category_id, code)
);

CREATE TYPE product_condition AS ENUM ('new', 'like_new', 'used_good', 'used_fair', 'refurbished');
CREATE TYPE product_status AS ENUM ('draft', 'active', 'sold_out', 'archived', 'banned');
CREATE TYPE moderation_status AS ENUM ('pending_review', 'approved', 'changes_requested', 'rejected');

CREATE TABLE products (
  id                uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  store_id          uuid NOT NULL REFERENCES stores(id),
  category_id       varchar(60) NOT NULL REFERENCES categories(id),
  brand_id          integer REFERENCES brands(id),
  size_chart_id     uuid REFERENCES size_charts(id),
  source_lang       lang_code NOT NULL,             -- bahasa yang diketik seller; dua lainnya diterjemahkan AI
  condition         product_condition NOT NULL DEFAULT 'new',
  condition_note    text,
  usage_months      smallint CHECK (usage_months >= 0),
  includes          jsonb NOT NULL DEFAULT '[]',
  warranty_type     varchar(16) NOT NULL DEFAULT 'none' CHECK (warranty_type IN ('none','store','official')),
  warranty_months   smallint NOT NULL DEFAULT 0 CHECK (warranty_months >= 0),
  return_days       smallint NOT NULL DEFAULT 7 CHECK (return_days >= 0),
  preorder_days     smallint CHECK (preorder_days > 0),
  min_price         numeric(12,2),                  -- cache dari varian aktif
  max_price         numeric(12,2),
  rating_avg        numeric(2,1),
  rating_count      integer NOT NULL DEFAULT 0,
  robots_index      boolean NOT NULL DEFAULT false,
  status            product_status NOT NULL DEFAULT 'draft',
  moderation_status moderation_status NOT NULL DEFAULT 'pending_review',
  moderation_note   text,
  published_at      timestamptz,
  created_at        timestamptz NOT NULL DEFAULT now(),
  updated_at        timestamptz NOT NULL DEFAULT now(),
  deleted_at        timestamptz,
  deleted_by        uuid REFERENCES users(id),
  -- Barang bukan baru wajib punya catatan kondisi.
  CONSTRAINT used_needs_note CHECK (condition = 'new' OR coalesce(length(trim(condition_note)), 0) > 0)
);
CREATE INDEX products_store_idx ON products (store_id, status) WHERE deleted_at IS NULL;
CREATE INDEX products_feed_idx ON products (published_at DESC)
  WHERE status = 'active' AND moderation_status = 'approved' AND deleted_at IS NULL;
CREATE INDEX products_category_idx ON products (category_id) WHERE deleted_at IS NULL;
CREATE INDEX products_moderation_queue ON products (updated_at) WHERE moderation_status = 'pending_review' AND deleted_at IS NULL;
CREATE TRIGGER products_updated BEFORE UPDATE ON products FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER products_audit AFTER INSERT OR UPDATE OR DELETE ON products FOR EACH ROW EXECUTE FUNCTION audit_row_change();

CREATE TABLE product_translations (
  product_id       uuid NOT NULL REFERENCES products(id),
  lang             lang_code NOT NULL,
  name             varchar(160) NOT NULL,
  description      text NOT NULL DEFAULT '',
  slug             varchar(180) NOT NULL CHECK (slug ~ '^[a-z0-9]+(-[a-z0-9]+)*$'),
  meta_title       varchar(70),
  meta_description varchar(160),
  keywords         text[] NOT NULL DEFAULT '{}',
  source           varchar(16) NOT NULL CHECK (source IN ('seller','machine','reviewer')),
  review_status    varchar(16) NOT NULL DEFAULT 'machine' CHECK (review_status IN ('machine','needs_review','reviewed')),
  source_lang      lang_code,                       -- untuk teks hasil terjemahan
  model            varchar(64),                     -- model AI yang menerjemahkan
  seller_edited    boolean NOT NULL DEFAULT false,  -- true = jangan ditimpa terjemahan ulang
  reviewed_by      uuid REFERENCES users(id),
  reviewed_at      timestamptz,
  updated_at       timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (product_id, lang)
);
CREATE UNIQUE INDEX product_translations_slug ON product_translations (lang, slug);
CREATE INDEX product_translations_review_queue ON product_translations (updated_at) WHERE review_status = 'needs_review';
CREATE TRIGGER product_translations_updated BEFORE UPDATE ON product_translations FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE product_options (
  id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  product_id uuid NOT NULL REFERENCES products(id) ON DELETE CASCADE,
  kind       varchar(16) NOT NULL CHECK (kind IN ('color','size','material','capacity','custom')),
  name_i18n  jsonb NOT NULL,
  position   smallint NOT NULL CHECK (position BETWEEN 1 AND 3),   -- maks. 3 opsi per produk
  UNIQUE (product_id, position)
);

CREATE TABLE product_option_values (
  id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  option_id  uuid NOT NULL REFERENCES product_options(id) ON DELETE CASCADE,
  value_i18n jsonb NOT NULL,
  swatch_hex char(7) CHECK (swatch_hex ~ '^#[0-9a-fA-F]{6}$'),
  position   smallint NOT NULL DEFAULT 0
);
CREATE INDEX product_option_values_option_idx ON product_option_values (option_id);

CREATE TABLE product_variants (
  id               uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  product_id       uuid NOT NULL REFERENCES products(id),
  store_id         uuid NOT NULL REFERENCES stores(id),  -- salinan untuk SKU unik per toko
  sku              varchar(64) NOT NULL,
  barcode          varchar(32),
  price            numeric(12,2) NOT NULL CHECK (price > 0),
  compare_at_price numeric(12,2) CHECK (compare_at_price IS NULL OR compare_at_price > price),
  stock_qty        integer NOT NULL DEFAULT 0 CHECK (stock_qty >= 0),
  reserved_qty     integer NOT NULL DEFAULT 0 CHECK (reserved_qty >= 0 AND reserved_qty <= stock_qty),
  weight_grams     integer CHECK (weight_grams > 0),
  length_mm        integer,
  width_mm         integer,
  height_mm        integer,
  is_default       boolean NOT NULL DEFAULT false,
  status           varchar(16) NOT NULL DEFAULT 'active' CHECK (status IN ('active','inactive')),
  created_at       timestamptz NOT NULL DEFAULT now(),
  updated_at       timestamptz NOT NULL DEFAULT now(),
  deleted_at       timestamptz
);
CREATE UNIQUE INDEX product_variants_sku ON product_variants (store_id, sku) WHERE deleted_at IS NULL;
CREATE INDEX product_variants_product_idx ON product_variants (product_id) WHERE deleted_at IS NULL;
CREATE TRIGGER product_variants_updated BEFORE UPDATE ON product_variants FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER product_variants_audit AFTER INSERT OR UPDATE OR DELETE ON product_variants FOR EACH ROW EXECUTE FUNCTION audit_row_change();

CREATE TABLE variant_option_values (
  variant_id      uuid NOT NULL REFERENCES product_variants(id) ON DELETE CASCADE,
  option_id       uuid NOT NULL REFERENCES product_options(id) ON DELETE CASCADE,
  option_value_id uuid NOT NULL REFERENCES product_option_values(id) ON DELETE CASCADE,
  PRIMARY KEY (variant_id, option_value_id),
  UNIQUE (variant_id, option_id)                   -- satu nilai per opsi
);

CREATE TABLE product_attribute_values (
  product_id   uuid NOT NULL REFERENCES products(id) ON DELETE CASCADE,
  attribute_id integer NOT NULL REFERENCES category_attributes(id),
  value_text   varchar(255),
  value_number numeric,
  PRIMARY KEY (product_id, attribute_id)
);

CREATE TABLE inventory_movements (
  id            bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  variant_id    uuid NOT NULL REFERENCES product_variants(id),
  delta         integer NOT NULL CHECK (delta <> 0),
  reason        varchar(16) NOT NULL CHECK (reason IN ('restock','reserve','release','sale','return','adjustment')),
  order_item_id uuid,                                -- FK ditambahkan saat tabel order dibuat (Tahap 3)
  created_by    uuid REFERENCES users(id),
  created_at    timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX inventory_movements_variant_idx ON inventory_movements (variant_id, created_at DESC);

CREATE TYPE media_kind AS ENUM ('video', 'photo', 'defect_photo');

CREATE TABLE product_media (
  id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  product_id      uuid NOT NULL REFERENCES products(id),
  option_value_id uuid REFERENCES product_option_values(id),  -- foto khusus warna tertentu
  kind            media_kind NOT NULL,
  storage_key     text NOT NULL,                    -- object key file asli (R2)
  video_url       text,                             -- hasil transcode / Cloudflare Stream
  poster_url      text,
  image_url       text,
  duration_sec    smallint CHECK (duration_sec BETWEEN 1 AND 180),
  has_audio       boolean,
  width           integer,
  height          integer,
  alt_text_i18n   jsonb NOT NULL DEFAULT '{}',
  focal_position  varchar(16),
  processing      varchar(16) NOT NULL DEFAULT 'uploaded' CHECK (processing IN ('uploaded','processing','ready','failed','blocked')),
  sort_order      smallint NOT NULL DEFAULT 0,      -- 0 = video utama
  uploaded_at     timestamptz NOT NULL DEFAULT now(),
  deleted_at      timestamptz
);
CREATE INDEX product_media_product_idx ON product_media (product_id, sort_order) WHERE deleted_at IS NULL;

CREATE TABLE media_processing_jobs (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  media_id    uuid NOT NULL REFERENCES product_media(id),
  step        varchar(16) NOT NULL CHECK (step IN ('transcode','thumbnail','moderation','embedding')),
  status      varchar(16) NOT NULL DEFAULT 'queued' CHECK (status IN ('queued','running','done','failed')),
  attempts    integer NOT NULL DEFAULT 0,
  error       text,
  started_at  timestamptz,
  finished_at timestamptz,
  created_at  timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX media_processing_jobs_media_idx ON media_processing_jobs (media_id);

-- +goose Down
DROP TABLE media_processing_jobs;
DROP TABLE product_media;
DROP TYPE media_kind;
DROP TABLE inventory_movements;
DROP TABLE product_attribute_values;
DROP TABLE variant_option_values;
DROP TABLE product_variants;
DROP TABLE product_option_values;
DROP TABLE product_options;
DROP TABLE product_translations;
DROP TABLE products;
DROP TYPE moderation_status;
DROP TYPE product_status;
DROP TYPE product_condition;
DROP TABLE category_attributes;
DROP TABLE size_charts;
DROP TABLE brands;
DROP TABLE category_translations;
DROP TABLE categories;
