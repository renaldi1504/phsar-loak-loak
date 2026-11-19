-- +goose Up
-- Halaman legal yang bisa diedit super admin (berversi, tidak pernah ditimpa) + retur & sengketa.

ALTER TABLE legal_documents DROP CONSTRAINT legal_documents_type_check;
ALTER TABLE legal_documents ADD CONSTRAINT legal_documents_type_check
  CHECK (type IN ('buyer_terms','seller_agreement','privacy_policy','ads_terms','return_policy','data_deletion'));
ALTER TABLE legal_documents ADD COLUMN published_by uuid REFERENCES users(id);
CREATE TRIGGER legal_documents_audit AFTER INSERT OR UPDATE OR DELETE ON legal_documents FOR EACH ROW EXECUTE FUNCTION audit_row_change();

-- Retur & sengketa (satu permintaan per pesanan yang sedang berjalan).
--   requested ──seller setuju──► approved ──admin catat transfer──► refunded
--       │  └──seller tolak / 48 jam tanpa jawaban──► disputed ──admin──► approved | rejected
--       └──pembeli batalkan──► cancelled
CREATE TYPE return_status AS ENUM ('requested', 'disputed', 'approved', 'rejected', 'refunded', 'cancelled');
CREATE TABLE returns (
  id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  return_no       varchar(20) NOT NULL UNIQUE,
  order_id        uuid NOT NULL REFERENCES orders(id),
  buyer_id        uuid NOT NULL REFERENCES users(id),
  store_id        uuid NOT NULL REFERENCES stores(id),
  status          return_status NOT NULL DEFAULT 'requested',
  reason          varchar(24) NOT NULL CHECK (reason IN ('changed_mind','wrong_item','damaged','not_as_described','missing_parts','other')),
  description     text NOT NULL CHECK (length(description) BETWEEN 10 AND 2000),
  photo_keys      text[] NOT NULL DEFAULT '{}',
  items           jsonb NOT NULL DEFAULT '[]',   -- [{order_item_id, qty}] barang yang diretur
  requested_amount numeric(12,2) NOT NULL CHECK (requested_amount > 0),
  seller_note     text,
  seller_responded_at timestamptz,
  seller_deadline timestamptz NOT NULL,
  escalated_reason varchar(24) CHECK (escalated_reason IN ('seller_rejected','seller_timeout')),
  decision_note   text,                          -- alasan keputusan admin (dilihat pembeli & seller)
  decided_by      uuid REFERENCES users(id),
  decided_at      timestamptz,
  refund_amount   numeric(12,2) CHECK (refund_amount > 0),
  needs_pickup    boolean NOT NULL DEFAULT false, -- barang perlu dijemput kurir dari pembeli
  refund_ref      varchar(80),
  refunded_by     uuid REFERENCES users(id),
  refunded_at     timestamptz,
  created_at      timestamptz NOT NULL DEFAULT now(),
  updated_at      timestamptz NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX returns_one_open ON returns (order_id) WHERE status IN ('requested','disputed','approved');
CREATE INDEX returns_store_idx ON returns (store_id, created_at DESC);
CREATE INDEX returns_status_idx ON returns (status, created_at);
CREATE TRIGGER returns_updated BEFORE UPDATE ON returns FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER returns_audit AFTER INSERT OR UPDATE OR DELETE ON returns FOR EACH ROW EXECUTE FUNCTION audit_row_change();
CREATE INDEX returns_deadline_idx ON returns (seller_deadline) WHERE status = 'requested';

-- Foto bukti retur dari pembeli (privat) & potongan saldo seller saat retur dikembalikan.
ALTER TABLE uploads DROP CONSTRAINT uploads_kind_check;
ALTER TABLE uploads ADD CONSTRAINT uploads_kind_check
  CHECK (kind IN ('store_logo','store_banner','pickup_photo','kyc_document','product_media','delivery_proof','return_photo'));
ALTER TABLE seller_ledger DROP CONSTRAINT seller_ledger_kind_check;
ALTER TABLE seller_ledger ADD CONSTRAINT seller_ledger_kind_check
  CHECK (kind IN ('order_earning','payout','payout_cancelled','adjustment','return_refund'));
ALTER TABLE seller_ledger ADD COLUMN return_id uuid REFERENCES returns(id);
CREATE UNIQUE INDEX seller_ledger_return_once ON seller_ledger (return_id) WHERE kind = 'return_refund';

-- +goose Down
DROP INDEX seller_ledger_return_once;
ALTER TABLE seller_ledger DROP COLUMN return_id;
DELETE FROM seller_ledger WHERE kind = 'return_refund';
ALTER TABLE seller_ledger DROP CONSTRAINT seller_ledger_kind_check;
ALTER TABLE seller_ledger ADD CONSTRAINT seller_ledger_kind_check CHECK (kind IN ('order_earning','payout','payout_cancelled','adjustment'));
DELETE FROM uploads WHERE kind = 'return_photo';
ALTER TABLE uploads DROP CONSTRAINT uploads_kind_check;
ALTER TABLE uploads ADD CONSTRAINT uploads_kind_check
  CHECK (kind IN ('store_logo','store_banner','pickup_photo','kyc_document','product_media','delivery_proof'));
DROP TABLE returns;
DROP TYPE return_status;
DROP TRIGGER legal_documents_audit ON legal_documents;
ALTER TABLE legal_documents DROP COLUMN published_by;
DELETE FROM legal_documents WHERE type IN ('return_policy','data_deletion');
ALTER TABLE legal_documents DROP CONSTRAINT legal_documents_type_check;
ALTER TABLE legal_documents ADD CONSTRAINT legal_documents_type_check CHECK (type IN ('buyer_terms','seller_agreement','privacy_policy','ads_terms'));
DELETE FROM platform_settings WHERE key = 'legal';
