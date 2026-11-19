-- +goose Up
-- Voucher & flash sale (dibiayai Phsar).
--   Voucher: potongan dari total yang dibayar pembeli. Penghasilan seller TIDAK berubah (seller_amount = subtotal − komisi);
--   Phsar yang menanggung potongan. Potongan dibagi ke tiap pesanan sesuai porsi subtotal (orders.discount).
--   Flash sale: harga khusus per varian dalam rentang waktu, boleh dibatasi jumlahnya (qty_cap).

CREATE TABLE vouchers (
  id               uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  code             varchar(24) NOT NULL CHECK (code ~ '^[A-Z0-9_-]{3,24}$'),
  name             varchar(80) NOT NULL DEFAULT '',
  kind             varchar(8) NOT NULL CHECK (kind IN ('percent','fixed')),
  value            numeric(12,2) NOT NULL CHECK (value > 0),
  max_discount     numeric(12,2) CHECK (max_discount IS NULL OR max_discount > 0),
  min_spend        numeric(12,2) NOT NULL DEFAULT 0 CHECK (min_spend >= 0),
  starts_at        timestamptz NOT NULL DEFAULT now(),
  ends_at          timestamptz,
  total_limit      integer CHECK (total_limit IS NULL OR total_limit > 0),
  per_user_limit   integer NOT NULL DEFAULT 1 CHECK (per_user_limit > 0),
  first_order_only boolean NOT NULL DEFAULT false,
  active           boolean NOT NULL DEFAULT true,
  created_by       uuid REFERENCES users(id),
  created_at       timestamptz NOT NULL DEFAULT now(),
  updated_at       timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT vouchers_percent CHECK (kind <> 'percent' OR value <= 90),
  CONSTRAINT vouchers_period CHECK (ends_at IS NULL OR ends_at > starts_at)
);
CREATE UNIQUE INDEX vouchers_code_key ON vouchers (upper(code));
CREATE TRIGGER vouchers_updated BEFORE UPDATE ON vouchers FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER vouchers_audit AFTER INSERT OR UPDATE OR DELETE ON vouchers FOR EACH ROW EXECUTE FUNCTION audit_row_change();

-- Pemakaian: held (checkout menunggu bayar) → used (lunas) / released (kedaluwarsa / batal).
CREATE TABLE voucher_redemptions (
  id          bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  voucher_id  uuid NOT NULL REFERENCES vouchers(id),
  checkout_id uuid NOT NULL UNIQUE REFERENCES checkouts(id),
  buyer_id    uuid NOT NULL REFERENCES users(id),
  discount    numeric(12,2) NOT NULL CHECK (discount > 0),
  status      varchar(10) NOT NULL DEFAULT 'held' CHECK (status IN ('held','used','released')),
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX voucher_redemptions_voucher_idx ON voucher_redemptions (voucher_id, buyer_id) WHERE status <> 'released';
CREATE TRIGGER voucher_redemptions_updated BEFORE UPDATE ON voucher_redemptions FOR EACH ROW EXECUTE FUNCTION set_updated_at();

ALTER TABLE checkouts
  ADD COLUMN voucher_id uuid REFERENCES vouchers(id),
  ADD COLUMN discount numeric(12,2) NOT NULL DEFAULT 0 CHECK (discount >= 0);
ALTER TABLE orders ADD COLUMN discount numeric(12,2) NOT NULL DEFAULT 0 CHECK (discount >= 0);
ALTER TABLE orders DROP CONSTRAINT orders_money;
ALTER TABLE orders ADD CONSTRAINT orders_money
  CHECK (total = subtotal + shipping_fee - discount AND seller_amount = subtotal - commission_amount);

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

-- Barang yang dibeli dengan harga flash sale (untuk menghitung sisa kuota).
ALTER TABLE order_items ADD COLUMN flash_sale_id uuid REFERENCES flash_sales(id);
CREATE INDEX order_items_flash_idx ON order_items (flash_sale_id) WHERE flash_sale_id IS NOT NULL;

-- +goose Down
DROP INDEX order_items_flash_idx;
ALTER TABLE order_items DROP COLUMN flash_sale_id;
DROP TABLE flash_sales;
ALTER TABLE orders DROP CONSTRAINT orders_money;
UPDATE orders SET total = total + discount WHERE discount > 0;
ALTER TABLE orders ADD CONSTRAINT orders_money CHECK (total = subtotal + shipping_fee AND seller_amount = subtotal - commission_amount);
ALTER TABLE orders DROP COLUMN discount;
ALTER TABLE checkouts DROP COLUMN discount, DROP COLUMN voucher_id;
DROP TABLE voucher_redemptions;
DROP TABLE vouchers;
