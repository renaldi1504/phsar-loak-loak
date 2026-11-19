-- +goose Up
-- extras:
--  A) Strike otomatis untuk pesanan yang batal karena telat (seller_timeout / pack_timeout).
--     Diatur super admin di platform_settings['auto_strikes'] (default MATI). Strike sistem: created_by NULL, auto = true,
--     satu strike per pesanan (order_id unik).
--  B) Tarik saldo manual oleh seller → payout pending biasa dengan requested_by = pemilik toko.
--     Diatur super admin di platform_settings['manual_payout'] (default MATI).

ALTER TABLE store_strikes DROP CONSTRAINT store_strikes_reason_check;
ALTER TABLE store_strikes ADD CONSTRAINT store_strikes_reason_check
  CHECK (reason IN ('counterfeit','prohibited','misleading','scam','offensive','other','seller_timeout','pack_timeout'));
ALTER TABLE store_strikes ALTER COLUMN created_by DROP NOT NULL;            -- NULL = sistem (strike otomatis)
ALTER TABLE store_strikes ADD COLUMN auto boolean NOT NULL DEFAULT false;
ALTER TABLE store_strikes ADD COLUMN order_id uuid REFERENCES orders(id);
CREATE UNIQUE INDEX store_strikes_order_once ON store_strikes (order_id) WHERE order_id IS NOT NULL;
CREATE INDEX store_strikes_auto_idx ON store_strikes (store_id, created_at) WHERE auto;

ALTER TABLE payouts ADD COLUMN requested_by uuid REFERENCES users(id);     -- NOT NULL = diminta seller ("Tarik sekarang")
CREATE INDEX payouts_requested_idx ON payouts (store_id, created_at) WHERE requested_by IS NOT NULL;

INSERT INTO platform_settings (key, value) VALUES
  ('auto_strikes', '{"enabled": false, "seller_timeout_points": 1, "pack_timeout_points": 1, "expiry_days": 30, "daily_max": 1}'),
  ('manual_payout', '{"enabled": false, "min_amount": "10.00", "max_per_7_days": 1}')
ON CONFLICT (key) DO NOTHING;

-- +goose Down
DELETE FROM platform_settings WHERE key IN ('auto_strikes', 'manual_payout');
DROP INDEX payouts_requested_idx;
ALTER TABLE payouts DROP COLUMN requested_by;
DELETE FROM store_strikes WHERE auto OR created_by IS NULL;
DROP INDEX store_strikes_auto_idx;
DROP INDEX store_strikes_order_once;
ALTER TABLE store_strikes DROP COLUMN order_id;
ALTER TABLE store_strikes DROP COLUMN auto;
ALTER TABLE store_strikes ALTER COLUMN created_by SET NOT NULL;
ALTER TABLE store_strikes DROP CONSTRAINT store_strikes_reason_check;
ALTER TABLE store_strikes ADD CONSTRAINT store_strikes_reason_check
  CHECK (reason IN ('counterfeit','prohibited','misleading','scam','offensive','other'));
