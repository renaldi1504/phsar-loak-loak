-- +goose Up
-- Alat katalog untuk toko dengan banyak SKU:
--   • SKU otomatis rapi: kode toko (stores.sku_prefix) + nomor urut produk per toko (store_sku_counters) + nilai varian,
--     mis. SKV-00012 / SKV-00012-RED-M. Nomor diambil sekali per produk (products.sku_no). SKU lama tidak diubah.
--   • Peringatan stok menipis: ambang per toko (stores.low_stock_threshold, bawaan 3) atau per varian, notifikasi
--     'stock.low' sekali per penurunan (product_variants.low_stock_notified), dibuat trigger supaya pesanan, edit,
--     impor CSV dan retur semuanya tercakup.

ALTER TABLE stores ADD COLUMN IF NOT EXISTS sku_prefix varchar(5);
ALTER TABLE stores ADD COLUMN IF NOT EXISTS low_stock_threshold integer NOT NULL DEFAULT 3;
-- +goose StatementBegin
DO $$ BEGIN
  ALTER TABLE stores ADD CONSTRAINT stores_sku_prefix_check CHECK (sku_prefix ~ '^[A-Z0-9]{2,5}$');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;
-- +goose StatementEnd
-- +goose StatementBegin
DO $$ BEGIN
  ALTER TABLE stores ADD CONSTRAINT stores_low_stock_threshold_check CHECK (low_stock_threshold BETWEEN 0 AND 1000);
EXCEPTION WHEN duplicate_object THEN NULL; END $$;
-- +goose StatementEnd

-- Penghitung terpisah dari tabel stores supaya mengambil nomor tidak membuat jejak audit di stores.
CREATE TABLE IF NOT EXISTS store_sku_counters (
  store_id uuid PRIMARY KEY REFERENCES stores(id),
  last_no  integer NOT NULL DEFAULT 0
);

ALTER TABLE products ADD COLUMN IF NOT EXISTS sku_no integer;
ALTER TABLE product_variants ADD COLUMN IF NOT EXISTS low_stock_threshold integer;
ALTER TABLE product_variants ADD COLUMN IF NOT EXISTS low_stock_notified boolean NOT NULL DEFAULT false;
-- +goose StatementBegin
DO $$ BEGIN
  ALTER TABLE product_variants ADD CONSTRAINT product_variants_low_stock_threshold_check CHECK (low_stock_threshold BETWEEN 0 AND 1000);
EXCEPTION WHEN duplicate_object THEN NULL; END $$;
-- +goose StatementEnd

-- Kode SKU dari nama toko (sama dengan seller.SKUPrefix di Go):
--   1 kata → 3 huruf pertama ("Beauty" → BEA) · 2 kata → 2 huruf + 1 huruf ("Skin Vibe" → SKV) · 3+ kata → inisial 3 kata.
--   Nama tanpa huruf latin (Khmer/Cina) → "S" + 3 karakter id toko.
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION phsar_sku_prefix(name text, store_id uuid) RETURNS text LANGUAGE plpgsql IMMUTABLE AS $$
DECLARE
  w text[];
  p text;
BEGIN
  w := array_remove(regexp_split_to_array(upper(coalesce(name, '')), '[^A-Z0-9]+'), '');
  IF coalesce(array_length(w, 1), 0) = 0 THEN
    p := '';
  ELSIF array_length(w, 1) = 1 THEN
    p := left(w[1], 3);
  ELSIF array_length(w, 1) = 2 THEN
    p := left(w[1], 2) || left(w[2], 1);
  ELSE
    p := left(w[1], 1) || left(w[2], 1) || left(w[3], 1);
  END IF;
  IF length(p) < 2 THEN
    p := 'S' || upper(left(replace(store_id::text, '-', ''), 3));
  END IF;
  RETURN p;
END $$;
-- +goose StatementEnd

UPDATE stores SET sku_prefix = phsar_sku_prefix(name, id) WHERE sku_prefix IS NULL;

-- Toko baru (dari jalur mana pun: Seller Center, admin, demo) langsung punya kode SKU.
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION phsar_store_sku_prefix() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
  IF NEW.sku_prefix IS NULL THEN
    NEW.sku_prefix := phsar_sku_prefix(NEW.name, NEW.id);
  END IF;
  RETURN NEW;
END $$;
-- +goose StatementEnd
DROP TRIGGER IF EXISTS stores_sku_prefix ON stores;
CREATE TRIGGER stores_sku_prefix BEFORE INSERT ON stores FOR EACH ROW EXECUTE FUNCTION phsar_store_sku_prefix();

-- Stok menipis: dicek saat stok/stok dipesan/status berubah. Stok tersedia = stock_qty − reserved_qty.
-- Notifikasi hanya untuk produk yang sedang dijual (bukan draf/arsip/terhapus) dan sekali sampai stok naik lagi.
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION phsar_low_stock() RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE
  thr int;
  avail int := NEW.stock_qty - NEW.reserved_qty;
  pname text;
  live boolean;
BEGIN
  SELECT coalesce(NEW.low_stock_threshold, s.low_stock_threshold) INTO thr FROM stores s WHERE s.id = NEW.store_id;
  IF thr IS NULL THEN
    RETURN NEW;
  END IF;
  IF avail > thr THEN
    NEW.low_stock_notified := false;
    RETURN NEW;
  END IF;
  IF NEW.low_stock_notified OR NEW.deleted_at IS NOT NULL OR NEW.status <> 'active' THEN
    RETURN NEW;
  END IF;
  SELECT p.status IN ('active', 'sold_out') AND p.deleted_at IS NULL, coalesce(t.name, '')
    INTO live, pname
    FROM products p LEFT JOIN product_translations t ON t.product_id = p.id AND t.lang = p.source_lang
    WHERE p.id = NEW.product_id;
  IF NOT coalesce(live, false) THEN
    RETURN NEW;
  END IF;
  NEW.low_stock_notified := true;
  INSERT INTO notifications (user_id, store_id, kind, params, link, audience)
  SELECT sm.user_id, NEW.store_id, 'stock.low',
         jsonb_build_object('product', left(pname, 120), 'note', NEW.sku, 'order', greatest(avail, 0)::text),
         '#/store/' || NEW.store_id || '/products', 'seller'
    FROM store_members sm JOIN users u ON u.id = sm.user_id AND u.deleted_at IS NULL
   WHERE sm.store_id = NEW.store_id;
  RETURN NEW;
END $$;
-- +goose StatementEnd

DROP TRIGGER IF EXISTS product_variants_low_stock ON product_variants;
CREATE TRIGGER product_variants_low_stock BEFORE UPDATE OF stock_qty, reserved_qty, status, low_stock_threshold ON product_variants
  FOR EACH ROW EXECUTE FUNCTION phsar_low_stock();

-- Cari cepat berdasarkan SKU (huruf besar/kecil sama) untuk impor CSV.
CREATE INDEX IF NOT EXISTS product_variants_store_sku_lower ON product_variants (store_id, lower(sku)) WHERE deleted_at IS NULL;

-- +goose Down
DROP INDEX IF EXISTS product_variants_store_sku_lower;
DROP TRIGGER IF EXISTS product_variants_low_stock ON product_variants;
DROP FUNCTION IF EXISTS phsar_low_stock();
DROP TRIGGER IF EXISTS stores_sku_prefix ON stores;
DROP FUNCTION IF EXISTS phsar_store_sku_prefix();
DROP FUNCTION IF EXISTS phsar_sku_prefix(text, uuid);
ALTER TABLE product_variants DROP CONSTRAINT IF EXISTS product_variants_low_stock_threshold_check;
ALTER TABLE product_variants DROP COLUMN IF EXISTS low_stock_notified;
ALTER TABLE product_variants DROP COLUMN IF EXISTS low_stock_threshold;
ALTER TABLE products DROP COLUMN IF EXISTS sku_no;
DROP TABLE IF EXISTS store_sku_counters;
ALTER TABLE stores DROP CONSTRAINT IF EXISTS stores_low_stock_threshold_check;
ALTER TABLE stores DROP CONSTRAINT IF EXISTS stores_sku_prefix_check;
ALTER TABLE stores DROP COLUMN IF EXISTS low_stock_threshold;
ALTER TABLE stores DROP COLUMN IF EXISTS sku_prefix;
