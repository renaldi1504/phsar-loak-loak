-- +goose Up
-- Pengaturan payment gateway oleh super admin (bukan lagi hanya di .env).
--   platform_settings['payment'] = {"mode": "off" | "payway" | "fake", "lifetime_min": 15}
--     (tidak ada baris = ikut PAYMENT_PROVIDER di .env server)
--   payment_gateways: satu baris per gateway. Rahasia (API key) disimpan terenkripsi (AES-GCM, DATA_ENCRYPTION_KEY)
--   dan tidak pernah dikirim balik ke browser — hanya 4 karakter terakhir.

CREATE TABLE payment_gateways (
  code        varchar(32) PRIMARY KEY CHECK (code IN ('payway')),
  environment varchar(16) NOT NULL DEFAULT 'sandbox' CHECK (environment IN ('sandbox', 'production')),
  config      jsonb NOT NULL DEFAULT '{}',   -- bukan rahasia: merchant_id, base_url
  secret_enc  bytea,                         -- API key terenkripsi
  secret_hint varchar(8),                    -- 4 karakter terakhir, untuk ditampilkan
  verified_at timestamptz,                   -- terakhir lolos "Tes koneksi"
  last_test   jsonb,
  updated_by  uuid REFERENCES users(id),
  updated_at  timestamptz NOT NULL DEFAULT now()
);
CREATE TRIGGER payment_gateways_updated BEFORE UPDATE ON payment_gateways FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER payment_gateways_audit AFTER INSERT OR UPDATE OR DELETE ON payment_gateways FOR EACH ROW EXECUTE FUNCTION audit_row_change();

-- +goose Down
DROP TABLE payment_gateways;
DELETE FROM platform_settings WHERE key = 'payment';
