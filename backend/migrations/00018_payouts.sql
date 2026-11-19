-- +goose Up
-- Saldo & payout seller.
--   pesanan terantar ──(masa tahan N hari, untuk komplain/retur)──► completed → saldo seller bertambah (seller_amount)
--   jadwal (harian / mingguan, diatur admin) → payout "pending" per toko (saldo dipotong) → admin transfer manual
--   → catat referensi → "paid". Batal → saldo dikembalikan.

CREATE TABLE platform_settings (
  key        varchar(64) PRIMARY KEY,
  value      jsonb NOT NULL,
  updated_by uuid REFERENCES users(id),
  updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE TRIGGER platform_settings_audit AFTER INSERT OR UPDATE OR DELETE ON platform_settings FOR EACH ROW EXECUTE FUNCTION audit_row_change();
INSERT INTO platform_settings (key, value) VALUES
  ('payout', '{"frequency": "weekly", "weekday": 1, "hour": 9, "hold_days": 3, "min_amount": "5.00"}');

ALTER TABLE orders ADD COLUMN completed_at timestamptz;

CREATE TABLE payouts (
  id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  store_id       uuid NOT NULL REFERENCES stores(id),
  amount         numeric(12,2) NOT NULL CHECK (amount > 0),
  currency       char(3) NOT NULL DEFAULT 'USD',
  status         varchar(12) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','paid','cancelled')),
  account_id     uuid REFERENCES seller_payout_accounts(id),
  bank           varchar(16) NOT NULL,                 -- salinan saat dibuat
  account_name   varchar(120) NOT NULL,
  account_last4  char(4) NOT NULL,
  order_count    integer NOT NULL DEFAULT 0,
  created_by     uuid REFERENCES users(id),            -- NULL = dibuat jadwal otomatis
  paid_at        timestamptz,
  paid_by        uuid REFERENCES users(id),
  reference      varchar(80),                          -- nomor transfer bank
  note           varchar(300),
  cancelled_at   timestamptz,
  created_at     timestamptz NOT NULL DEFAULT now(),
  updated_at     timestamptz NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX payouts_one_pending ON payouts (store_id) WHERE status = 'pending';
CREATE INDEX payouts_store_idx ON payouts (store_id, created_at DESC);
CREATE TRIGGER payouts_updated BEFORE UPDATE ON payouts FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER payouts_audit AFTER INSERT OR UPDATE OR DELETE ON payouts FOR EACH ROW EXECUTE FUNCTION audit_row_change();

-- Buku besar saldo seller: saldo = SUM(amount). Tidak pernah di-UPDATE/DELETE, hanya ditambah.
CREATE TABLE seller_ledger (
  id         bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  store_id   uuid NOT NULL REFERENCES stores(id),
  kind       varchar(20) NOT NULL CHECK (kind IN ('order_earning','payout','payout_cancelled','adjustment')),
  amount     numeric(12,2) NOT NULL CHECK (amount <> 0),
  order_id   uuid REFERENCES orders(id),
  payout_id  uuid REFERENCES payouts(id),
  note       varchar(300),
  created_by uuid REFERENCES users(id),
  created_at timestamptz NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX seller_ledger_order_once ON seller_ledger (order_id) WHERE kind = 'order_earning';
CREATE INDEX seller_ledger_store_idx ON seller_ledger (store_id, id);

-- Jadwal payout yang sudah dijalankan (supaya tidak dobel di hari yang sama).
CREATE TABLE payout_runs (
  id           bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  run_date     date NOT NULL,
  payout_count integer NOT NULL,
  total        numeric(12,2) NOT NULL,
  manual       boolean NOT NULL DEFAULT false,
  created_by   uuid REFERENCES users(id),
  created_at   timestamptz NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX payout_runs_auto_once ON payout_runs (run_date) WHERE NOT manual;

-- +goose Down
DROP TABLE payout_runs;
DROP TABLE seller_ledger;
DROP TABLE payouts;
ALTER TABLE orders DROP COLUMN completed_at;
DROP TABLE platform_settings;
