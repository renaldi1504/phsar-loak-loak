-- +goose Up
-- Faktur komisi bulanan per toko (Phsar → seller). Isi faktur disalin saat diterbitkan (tidak ikut berubah
-- kalau pesanan berubah). Batal (void) tetap menyimpan nomornya; faktur pengganti mendapat nomor baru.

CREATE TABLE invoices (
  id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  number         varchar(40) NOT NULL UNIQUE,              -- PHS-202611-0001
  store_id       uuid NOT NULL REFERENCES stores(id),
  period         date NOT NULL CHECK (extract(day FROM period) = 1),  -- hari pertama bulan
  status         varchar(8) NOT NULL DEFAULT 'issued' CHECK (status IN ('issued','void')),
  seller         jsonb NOT NULL,                            -- salinan: nama toko, nama resmi, TIN, alamat
  issuer         jsonb NOT NULL,                            -- salinan pengaturan pajak Phsar
  lines          jsonb NOT NULL,                            -- [{kind, ref, date, subtotal, commission}]
  order_count    integer NOT NULL,
  sales_total    numeric(12,2) NOT NULL,                    -- jumlah subtotal pesanan (informasi)
  commission     numeric(12,2) NOT NULL,                    -- komisi bersih (setelah komisi yang dikembalikan)
  vat_rate       numeric(5,2) NOT NULL CHECK (vat_rate BETWEEN 0 AND 100),
  vat_inclusive  boolean NOT NULL,
  subtotal       numeric(12,2) NOT NULL,                    -- sebelum PPN
  vat_amount     numeric(12,2) NOT NULL,
  total          numeric(12,2) NOT NULL,
  replaces_id    uuid REFERENCES invoices(id),
  issued_by      uuid REFERENCES users(id),
  issued_at      timestamptz NOT NULL DEFAULT now(),
  voided_by      uuid REFERENCES users(id),
  voided_at      timestamptz,
  void_note      varchar(300),
  CONSTRAINT invoices_money CHECK (total = subtotal + vat_amount)
);
CREATE UNIQUE INDEX invoices_one_issued ON invoices (store_id, period) WHERE status = 'issued';
CREATE INDEX invoices_period_idx ON invoices (period DESC, number);
CREATE INDEX invoices_store_idx ON invoices (store_id, period DESC);
CREATE TRIGGER invoices_audit AFTER INSERT OR UPDATE OR DELETE ON invoices FOR EACH ROW EXECUTE FUNCTION audit_row_change();

-- Nomor urut per bulan (tidak pernah dipakai ulang).
CREATE TABLE invoice_counters (
  period date PRIMARY KEY,
  last   integer NOT NULL DEFAULT 0
);

-- +goose Down
DROP TABLE invoice_counters;
DROP TABLE invoices;
DELETE FROM platform_settings WHERE key = 'tax';
