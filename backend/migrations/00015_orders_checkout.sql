-- +goose Up
-- Checkout & pesanan. Satu checkout = satu pembayaran KHQR; dipecah menjadi satu pesanan per toko.
--
-- Alur: pembeli menekan "Bayar" → checkout (pending) + pesanan (pending_payment) dibuat, stok ditahan
-- (reserved_qty) → QR KHQR dari ABA PayWay → callback ABA dicek ulang ke API check-transaction →
-- checkout paid, pesanan paid (masuk ke seller), stok terjual. Tidak dibayar sampai expires_at →
-- checkout expired, pesanan dibatalkan, stok dilepas.

CREATE TABLE checkouts (
  id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  buyer_id        uuid NOT NULL REFERENCES users(id),
  idempotency_key varchar(64) NOT NULL,                 -- dari browser: klik dua kali = checkout yang sama
  tran_id         varchar(20) NOT NULL UNIQUE,          -- ID transaksi di PayWay (maks. 20 karakter)
  status          varchar(12) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','paid','expired','cancelled')),
  provider        varchar(12) NOT NULL CHECK (provider IN ('payway','fake')),
  currency        char(3) NOT NULL DEFAULT 'USD' CHECK (currency = 'USD'),
  subtotal        numeric(12,2) NOT NULL CHECK (subtotal > 0),
  shipping_fee    numeric(12,2) NOT NULL DEFAULT 0 CHECK (shipping_fee >= 0),
  total           numeric(12,2) NOT NULL CHECK (total > 0),
  -- alamat pengiriman (salinan; alamat buku alamat boleh berubah setelahnya)
  recipient       varchar(80) NOT NULL,
  phone           varchar(30) NOT NULL,
  address_line    varchar(300) NOT NULL,
  address_note    varchar(300),
  lat             numeric(9,6) NOT NULL CHECK (lat BETWEEN -90 AND 90),
  lng             numeric(9,6) NOT NULL CHECK (lng BETWEEN -180 AND 180),
  google_place_id varchar(255),
  service_area_code varchar(32) NOT NULL REFERENCES service_areas(code),
  lang            lang_code NOT NULL DEFAULT 'en',
  -- pembayaran
  qr_string       text,
  qr_image        text,                                 -- data:image/png;base64,… dari PayWay
  deeplink        text,                                 -- buka aplikasi ABA langsung
  expires_at      timestamptz NOT NULL,
  paid_at         timestamptz,
  paid_amount     numeric(12,2),
  payment_type    varchar(20),                          -- KHQR · ABA Pay
  payment_ref     varchar(64),                          -- kode approval (apv)
  payer_account   varchar(64),
  created_at      timestamptz NOT NULL DEFAULT now(),
  updated_at      timestamptz NOT NULL DEFAULT now(),
  UNIQUE (buyer_id, idempotency_key)
);
CREATE INDEX checkouts_buyer_idx ON checkouts (buyer_id, created_at DESC);
CREATE INDEX checkouts_pending_idx ON checkouts (expires_at) WHERE status = 'pending';
CREATE TRIGGER checkouts_updated BEFORE UPDATE ON checkouts FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER checkouts_audit AFTER INSERT OR UPDATE OR DELETE ON checkouts FOR EACH ROW EXECUTE FUNCTION audit_row_change();

-- Status pesanan. Tahap ini: pending_payment → paid / cancelled. Sisanya dipakai tahap seller & kurir.
CREATE TYPE order_status AS ENUM (
  'pending_payment', 'paid', 'accepted', 'ready_for_pickup', 'picked_up', 'delivered', 'completed',
  'cancelled', 'refunded'
);

CREATE TABLE orders (
  id                uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  order_no          varchar(20) NOT NULL UNIQUE,         -- PH-241114-7K2Q9 (dibaca pembeli, seller, kurir)
  checkout_id       uuid NOT NULL REFERENCES checkouts(id),
  buyer_id          uuid NOT NULL REFERENCES users(id),
  store_id          uuid NOT NULL REFERENCES stores(id),
  status            order_status NOT NULL DEFAULT 'pending_payment',
  subtotal          numeric(12,2) NOT NULL CHECK (subtotal > 0),
  shipping_fee      numeric(12,2) NOT NULL DEFAULT 0 CHECK (shipping_fee >= 0),
  total             numeric(12,2) NOT NULL CHECK (total > 0),
  commission_rate   numeric(5,2) NOT NULL CHECK (commission_rate BETWEEN 0 AND 100),  -- salinan saat pesan
  commission_amount numeric(12,2) NOT NULL CHECK (commission_amount >= 0),
  seller_amount     numeric(12,2) NOT NULL CHECK (seller_amount >= 0),
  item_count        integer NOT NULL CHECK (item_count > 0),
  cancel_reason     varchar(40),                         -- payment_expired · buyer · seller_rejected · …
  paid_at           timestamptz,
  cancelled_at      timestamptz,
  created_at        timestamptz NOT NULL DEFAULT now(),
  updated_at        timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT orders_money CHECK (total = subtotal + shipping_fee AND seller_amount = subtotal - commission_amount)
);
CREATE INDEX orders_buyer_idx ON orders (buyer_id, created_at DESC);
CREATE INDEX orders_store_idx ON orders (store_id, status, created_at DESC);
CREATE INDEX orders_checkout_idx ON orders (checkout_id);
CREATE TRIGGER orders_updated BEFORE UPDATE ON orders FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER orders_audit AFTER INSERT OR UPDATE OR DELETE ON orders FOR EACH ROW EXECUTE FUNCTION audit_row_change();

-- Barang di pesanan: salinan nama/harga/foto saat dipesan (produk boleh diubah seller setelahnya).
CREATE TABLE order_items (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id    uuid NOT NULL REFERENCES orders(id),
  product_id  uuid NOT NULL REFERENCES products(id),
  variant_id  uuid NOT NULL REFERENCES product_variants(id),
  name        varchar(160) NOT NULL,                    -- dalam bahasa pembeli saat memesan
  options     varchar(200) NOT NULL DEFAULT '',          -- "Merah / M"
  sku         varchar(64) NOT NULL,
  image_url   text,
  unit_price  numeric(12,2) NOT NULL CHECK (unit_price > 0),
  qty         integer NOT NULL CHECK (qty BETWEEN 1 AND 99),
  line_total  numeric(12,2) NOT NULL CHECK (line_total = unit_price * qty)
);
CREATE INDEX order_items_order_idx ON order_items (order_id);
CREATE INDEX order_items_variant_idx ON order_items (variant_id);

ALTER TABLE inventory_movements ADD CONSTRAINT inventory_movements_order_item_fk
  FOREIGN KEY (order_item_id) REFERENCES order_items(id);

-- Semua yang datang dari / ditanyakan ke penyedia pembayaran (bukti untuk sengketa & rekonsiliasi).
CREATE TABLE payment_events (
  id          bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  checkout_id uuid REFERENCES checkouts(id),
  tran_id     varchar(20),
  source      varchar(12) NOT NULL CHECK (source IN ('create','callback','check','fake')),
  verified    boolean NOT NULL DEFAULT false,            -- tanda tangan cocok / hasil dari API resmi
  status      varchar(20),                               -- APPROVED · PENDING · DECLINED · …
  payload     jsonb NOT NULL DEFAULT '{}',
  created_at  timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX payment_events_checkout_idx ON payment_events (checkout_id, created_at);

-- +goose Down
DROP TABLE payment_events;
ALTER TABLE inventory_movements DROP CONSTRAINT inventory_movements_order_item_fk;
DROP TABLE order_items;
DROP TABLE orders;
DROP TYPE order_status;
DROP TABLE checkouts;
