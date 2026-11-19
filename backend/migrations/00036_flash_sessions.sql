-- +goose Up
-- Flash sale berbasis sesi (seperti marketplace besar):
--   admin membuat sesi (jam bebas) → produk masuk dari admin (pilih banyak sekaligus) atau seller mendaftar → admin menyetujui.
--   Pembiayaan dipilih per sesi: phsar_share_pct = bagian potongan yang ditanggung Phsar (0 = seller, 100 = Phsar, 1–99 = patungan).
--   Kalau seller ikut menanggung, produk yang dipilih admin berstatus "invited" sampai seller menerima (seller tidak pernah rugi tanpa setuju).
-- Menggantikan tabel flash_sales (per varian) dari 00031.

DROP INDEX IF EXISTS order_items_flash_idx;
ALTER TABLE order_items DROP COLUMN IF EXISTS flash_sale_id;
DROP TABLE IF EXISTS flash_sales;

CREATE TABLE flash_sessions (
  id               uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name_en          varchar(80) NOT NULL,
  name_km          varchar(80) NOT NULL DEFAULT '',
  name_zh          varchar(80) NOT NULL DEFAULT '',
  starts_at        timestamptz NOT NULL,
  ends_at          timestamptz NOT NULL,
  phsar_share_pct  smallint NOT NULL DEFAULT 0 CHECK (phsar_share_pct BETWEEN 0 AND 100),
  allow_vouchers   boolean NOT NULL DEFAULT true,          -- barang flash boleh ikut potongan voucher Phsar
  min_discount_pct smallint NOT NULL DEFAULT 10 CHECK (min_discount_pct BETWEEN 1 AND 90),
  per_buyer        smallint NOT NULL DEFAULT 2 CHECK (per_buyer BETWEEN 1 AND 99), -- batas beli per pembeli per barang
  open_for_sellers boolean NOT NULL DEFAULT false,         -- seller boleh mendaftarkan produk
  register_until   timestamptz,                            -- batas pendaftaran seller (kosong = sampai sesi mulai)
  published        boolean NOT NULL DEFAULT false,         -- tampil ke pembeli (draf = false)
  cancelled_at     timestamptz,
  created_by       uuid REFERENCES users(id),
  created_at       timestamptz NOT NULL DEFAULT now(),
  updated_at       timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT flash_sessions_period CHECK (ends_at > starts_at AND ends_at <= starts_at + interval '14 days'),
  CONSTRAINT flash_sessions_register CHECK (register_until IS NULL OR register_until <= starts_at)
);
CREATE INDEX flash_sessions_time_idx ON flash_sessions (starts_at, ends_at) WHERE cancelled_at IS NULL;
CREATE TRIGGER flash_sessions_updated BEFORE UPDATE ON flash_sessions FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER flash_sessions_audit AFTER INSERT OR UPDATE OR DELETE ON flash_sessions FOR EACH ROW EXECUTE FUNCTION audit_row_change();

CREATE TABLE flash_items (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  session_id   uuid NOT NULL REFERENCES flash_sessions(id),
  variant_id   uuid NOT NULL REFERENCES product_variants(id),
  product_id   uuid NOT NULL REFERENCES products(id),
  store_id     uuid NOT NULL REFERENCES stores(id),
  normal_price numeric(12,2) NOT NULL CHECK (normal_price > 0),  -- harga normal saat diajukan (acuan diskon & subsidi)
  sale_price   numeric(12,2) NOT NULL CHECK (sale_price > 0),
  qty_cap      integer NOT NULL CHECK (qty_cap BETWEEN 1 AND 100000),
  per_buyer    smallint CHECK (per_buyer IS NULL OR per_buyer BETWEEN 1 AND 99), -- kosong = ikut sesi
  source       varchar(8) NOT NULL CHECK (source IN ('admin','seller')),
  status       varchar(10) NOT NULL CHECK (status IN ('invited','pending','approved','rejected','declined','withdrawn','removed')),
  note         varchar(300),
  submitted_by uuid REFERENCES users(id),
  decided_by   uuid REFERENCES users(id),
  decided_at   timestamptz,
  created_at   timestamptz NOT NULL DEFAULT now(),
  updated_at   timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT flash_items_price CHECK (sale_price < normal_price)
);
CREATE UNIQUE INDEX flash_items_one ON flash_items (session_id, variant_id) WHERE status IN ('invited','pending','approved');
CREATE INDEX flash_items_variant_idx ON flash_items (variant_id) WHERE status = 'approved';
CREATE INDEX flash_items_store_idx ON flash_items (store_id, created_at DESC);
CREATE TRIGGER flash_items_updated BEFORE UPDATE ON flash_items FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER flash_items_audit AFTER INSERT OR UPDATE OR DELETE ON flash_items FOR EACH ROW EXECUTE FUNCTION audit_row_change();

-- Barang pesanan yang dibeli dengan harga flash + subsidi Phsar untuk baris itu.
ALTER TABLE order_items
  ADD COLUMN flash_item_id uuid REFERENCES flash_items(id),
  ADD COLUMN flash_subsidy numeric(12,2) NOT NULL DEFAULT 0 CHECK (flash_subsidy >= 0);
CREATE INDEX order_items_flash_item_idx ON order_items (flash_item_id) WHERE flash_item_id IS NOT NULL;

-- Subsidi flash yang ditanggung Phsar per pesanan: seller menerima seolah-olah harganya lebih tinggi.
ALTER TABLE orders ADD COLUMN flash_subsidy numeric(12,2) NOT NULL DEFAULT 0 CHECK (flash_subsidy >= 0);
ALTER TABLE orders DROP CONSTRAINT orders_money;
ALTER TABLE orders ADD CONSTRAINT orders_money
  CHECK (total = subtotal + shipping_fee - discount AND seller_amount = subtotal + flash_subsidy - commission_amount);

-- Barang flash yang sedang berlaku untuk pembeli (sesi terbit, tidak batal, sedang berjalan; barang disetujui).
CREATE VIEW flash_live AS
SELECT i.id, i.session_id, i.variant_id, i.product_id, i.store_id, i.normal_price, i.sale_price, i.qty_cap,
       coalesce(i.per_buyer, s.per_buyer) AS per_buyer, s.phsar_share_pct, s.allow_vouchers, s.starts_at, s.ends_at,
       coalesce((SELECT sum(oi.qty) FROM order_items oi JOIN orders o ON o.id = oi.order_id
                 WHERE oi.flash_item_id = i.id AND o.status NOT IN ('cancelled','refunded')), 0)::int AS sold
FROM flash_items i JOIN flash_sessions s ON s.id = i.session_id
WHERE i.status = 'approved' AND s.published AND s.cancelled_at IS NULL AND s.starts_at <= now() AND s.ends_at > now();

-- +goose Down
DROP VIEW flash_live;
ALTER TABLE orders DROP CONSTRAINT orders_money;
ALTER TABLE orders ADD CONSTRAINT orders_money
  CHECK (total = subtotal + shipping_fee - discount AND seller_amount = subtotal - commission_amount);
ALTER TABLE orders DROP COLUMN flash_subsidy;
DROP INDEX order_items_flash_item_idx;
ALTER TABLE order_items DROP COLUMN flash_subsidy, DROP COLUMN flash_item_id;
DROP TABLE flash_items;
DROP TABLE flash_sessions;
CREATE TABLE flash_sales (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  variant_id  uuid NOT NULL REFERENCES product_variants(id),
  sale_price  numeric(12,2) NOT NULL CHECK (sale_price > 0),
  starts_at   timestamptz NOT NULL,
  ends_at     timestamptz NOT NULL,
  qty_cap     integer CHECK (qty_cap IS NULL OR qty_cap > 0),
  active      boolean NOT NULL DEFAULT true,
  created_by  uuid REFERENCES users(id),
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT flash_sales_period CHECK (ends_at > starts_at)
);
CREATE INDEX flash_sales_variant_idx ON flash_sales (variant_id, ends_at) WHERE active;
CREATE TRIGGER flash_sales_updated BEFORE UPDATE ON flash_sales FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER flash_sales_audit AFTER INSERT OR UPDATE OR DELETE ON flash_sales FOR EACH ROW EXECUTE FUNCTION audit_row_change();
ALTER TABLE order_items ADD COLUMN flash_sale_id uuid REFERENCES flash_sales(id);
CREATE INDEX order_items_flash_idx ON order_items (flash_sale_id) WHERE flash_sale_id IS NOT NULL;
