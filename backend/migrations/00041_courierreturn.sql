-- +goose Up
-- courierreturn: pengantaran gagal → "kembalikan ke toko".
--
--   picked_up ──gagal──► failed ──coba lagi (maks. 3 percobaan)──► picked_up
--                          │
--                          └──kembalikan ke toko (kurir / admin)──► returning ──kurir: sudah di toko──► returned
--                                                                                                       │
--   seller konfirmasi diterima (atau admin, atau otomatis setelah 24 jam) ◄──────────────────────────────┘
--     → pesanan batal (cancel_reason = delivery_failed) + refund_status = pending + stok kembali (kecuali barang rusak)
--
-- Ongkir pengembalian hanya dicatat (biaya kurir ditanggung Phsar), belum dibebankan ke seller.

ALTER TABLE deliveries DROP CONSTRAINT deliveries_status_check;
ALTER TABLE deliveries ADD CONSTRAINT deliveries_status_check
  CHECK (status IN ('unassigned','assigned','picked_up','delivered','failed','cancelled','returning','returned'));
ALTER TABLE deliveries ADD COLUMN fail_proof_key text; -- foto bukti gagal antar (privat, opsional)

CREATE TABLE courier_returns (
  id                 uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  delivery_id        uuid NOT NULL UNIQUE REFERENCES deliveries(id),
  order_id           uuid NOT NULL UNIQUE REFERENCES orders(id),
  store_id           uuid NOT NULL REFERENCES stores(id),
  status             varchar(16) NOT NULL DEFAULT 'returning' CHECK (status IN ('returning','awaiting_seller','received')),
  reason             varchar(24),                     -- salinan deliveries.fail_reason saat pengembalian dimulai
  note               varchar(300),
  attempts           integer NOT NULL DEFAULT 0,
  return_fee         numeric(12,2) NOT NULL DEFAULT 0 CHECK (return_fee >= 0), -- biaya kurir pengembalian (dicatat saja)
  fee_payer          varchar(8) NOT NULL DEFAULT 'phsar' CHECK (fee_payer IN ('phsar','seller')),
  started_by_type    actor_type NOT NULL,
  started_by         uuid REFERENCES users(id),
  started_at         timestamptz NOT NULL DEFAULT now(),
  returned_by        uuid REFERENCES users(id),       -- kurir yang menyerahkan ke toko
  returned_at        timestamptz,
  returned_proof_key text,                            -- foto serah terima di toko (privat, opsional)
  confirm_by         timestamptz,                     -- lewat batas ini tanpa konfirmasi seller → dianggap diterima
  resolved_by_type   actor_type,                      -- seller · admin · system
  received_by        uuid REFERENCES users(id),
  received_at        timestamptz,
  item_condition     varchar(8) CHECK (item_condition IN ('ok','damaged')),
  restocked          boolean,
  seller_note        varchar(300),
  created_at         timestamptz NOT NULL DEFAULT now(),
  updated_at         timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX courier_returns_open_idx ON courier_returns (status, started_at) WHERE status IN ('returning','awaiting_seller');
CREATE INDEX courier_returns_store_idx ON courier_returns (store_id, created_at DESC);
CREATE INDEX courier_returns_due_idx ON courier_returns (confirm_by) WHERE status = 'awaiting_seller';
CREATE TRIGGER courier_returns_updated BEFORE UPDATE ON courier_returns FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER courier_returns_audit AFTER INSERT OR UPDATE OR DELETE ON courier_returns FOR EACH ROW EXECUTE FUNCTION audit_row_change();

-- +goose Down
DROP TABLE courier_returns;
ALTER TABLE deliveries DROP COLUMN fail_proof_key;
UPDATE deliveries SET status = 'failed' WHERE status IN ('returning','returned');
ALTER TABLE deliveries DROP CONSTRAINT deliveries_status_check;
ALTER TABLE deliveries ADD CONSTRAINT deliveries_status_check
  CHECK (status IN ('unassigned','assigned','picked_up','delivered','failed','cancelled'));
