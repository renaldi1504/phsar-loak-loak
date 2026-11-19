-- +goose Up
-- Rekonsiliasi bank: ekspor transaksi ABA PayWay (CSV) dicocokkan dengan pembayaran & refund di Phsar.
--   statement → baris dari file PayWay; phsar → pembayaran yang tercatat lunas di Phsar tapi tidak ada di file (dalam periode).

CREATE TABLE recon_batches (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  file_name    varchar(200) NOT NULL,
  period_from  date,
  period_to    date,
  row_count    integer NOT NULL DEFAULT 0,
  collected    numeric(14,2) NOT NULL DEFAULT 0,   -- pembayaran berhasil menurut PayWay
  refunded     numeric(14,2) NOT NULL DEFAULT 0,
  fees         numeric(14,2) NOT NULL DEFAULT 0,   -- kalau file punya kolom biaya
  uploaded_by  uuid NOT NULL REFERENCES users(id),
  matched_at   timestamptz,
  created_at   timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX recon_batches_created_idx ON recon_batches (created_at DESC);
CREATE TRIGGER recon_batches_audit AFTER INSERT OR UPDATE OR DELETE ON recon_batches FOR EACH ROW EXECUTE FUNCTION audit_row_change();

CREATE TABLE recon_rows (
  id           bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  batch_id     uuid NOT NULL REFERENCES recon_batches(id) ON DELETE CASCADE,
  source       varchar(10) NOT NULL DEFAULT 'statement' CHECK (source IN ('statement','phsar')),
  line_no      integer,                         -- nomor baris di file (statement)
  kind         varchar(8) NOT NULL DEFAULT 'payment' CHECK (kind IN ('payment','refund')),
  tran_id      varchar(64),
  reference    varchar(80),                     -- APV / nomor referensi refund
  status       varchar(40),
  amount       numeric(14,2),
  currency     varchar(3),
  fee          numeric(14,2),
  occurred_at  timestamptz,
  raw          jsonb NOT NULL DEFAULT '{}',
  result       varchar(24) NOT NULL DEFAULT 'pending' CHECK (result IN ('pending','ok','amount_mismatch','not_paid_in_phsar','unknown',
                 'duplicate','ignored','missing_in_statement','refund_ok','refund_mismatch','refund_unrecorded')),
  checkout_id  uuid REFERENCES checkouts(id),
  expected     numeric(14,2),                   -- jumlah menurut Phsar
  resolved_by  uuid REFERENCES users(id),
  resolved_at  timestamptz,
  resolution   text
);
CREATE INDEX recon_rows_batch_idx ON recon_rows (batch_id, result, id);
CREATE INDEX recon_rows_tran_idx ON recon_rows (tran_id);
-- (tanpa trigger audit per baris: pencocokan ulang mengubah ribuan baris; penyelesaian dicatat ke audit_logs oleh aplikasi)

-- +goose Down
DROP TABLE recon_rows;
DROP TABLE recon_batches;
