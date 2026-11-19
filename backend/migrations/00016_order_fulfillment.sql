-- +goose Up
-- Seller memproses pesanan: terima → siapkan → siap dijemput; tolak / tidak merespons 2 jam → batal + refund.
-- Refund sementara manual (admin mengembalikan dari dashboard PayWay lalu menandai di sini).

ALTER TABLE orders
  ADD COLUMN respond_by       timestamptz,                -- paid_at + 2 jam: batas seller menerima
  ADD COLUMN accepted_at      timestamptz,
  ADD COLUMN ready_at         timestamptz,
  ADD COLUMN pickup_point_id  uuid REFERENCES store_pickup_points(id),
  ADD COLUMN reject_reason    varchar(40),                -- out_of_stock · cannot_fulfil · other
  ADD COLUMN seller_note      varchar(300),
  ADD COLUMN handled_by       uuid REFERENCES users(id),  -- anggota toko yang terakhir memproses
  ADD COLUMN refund_status    varchar(12) NOT NULL DEFAULT 'none' CHECK (refund_status IN ('none','pending','refunded')),
  ADD COLUMN refund_amount    numeric(12,2) CHECK (refund_amount IS NULL OR refund_amount >= 0),
  ADD COLUMN refunded_at      timestamptz,
  ADD COLUMN refund_ref       varchar(80),                -- nomor transaksi refund di PayWay
  ADD COLUMN refund_note      varchar(300),
  ADD COLUMN refunded_by      uuid REFERENCES users(id);

CREATE INDEX orders_respond_idx ON orders (respond_by) WHERE status = 'paid';
CREATE INDEX orders_refund_idx ON orders (cancelled_at) WHERE refund_status = 'pending';

-- Pesanan lama yang sudah dibayar: beri batas waktu dari sekarang.
UPDATE orders SET respond_by = now() + interval '2 hours' WHERE status = 'paid' AND respond_by IS NULL;

-- Riwayat status (untuk timeline pembeli/seller dan penyelesaian sengketa).
CREATE TABLE order_events (
  id         bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  order_id   uuid NOT NULL REFERENCES orders(id),
  status     varchar(24) NOT NULL,                 -- status baru, atau refund_pending / refunded
  actor_type actor_type NOT NULL,
  actor_id   uuid REFERENCES users(id),
  note       varchar(300),
  created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX order_events_order_idx ON order_events (order_id, id);

-- +goose Down
DROP TABLE order_events;
DROP INDEX IF EXISTS orders_refund_idx;
DROP INDEX IF EXISTS orders_respond_idx;
ALTER TABLE orders
  DROP COLUMN respond_by, DROP COLUMN accepted_at, DROP COLUMN ready_at, DROP COLUMN pickup_point_id,
  DROP COLUMN reject_reason, DROP COLUMN seller_note, DROP COLUMN handled_by, DROP COLUMN refund_status,
  DROP COLUMN refund_amount, DROP COLUMN refunded_at, DROP COLUMN refund_ref, DROP COLUMN refund_note, DROP COLUMN refunded_by;
