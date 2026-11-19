-- +goose Up
-- Login Facebook, provider SMS (PlasGate) yang diatur super admin, dan checkout sebagai tamu.
--   integrations: satu baris per layanan luar. Rahasia (secret PlasGate, app secret Facebook) terenkripsi (DATA_ENCRYPTION_KEY),
--   tidak pernah dikirim balik ke browser — hanya 4 karakter terakhir.
--   platform_settings['sms'] = {"provider": "off" | "plasgate" | "log"}  (tidak ada baris = ikut SMS_PROVIDER di .env)

CREATE TABLE integrations (
  code        varchar(32) PRIMARY KEY CHECK (code IN ('plasgate', 'facebook')),
  enabled     boolean NOT NULL DEFAULT false,
  config      jsonb NOT NULL DEFAULT '{}',   -- bukan rahasia: sender PlasGate, app_id Facebook
  secret_enc  bytea,
  secret_hint varchar(8),
  verified_at timestamptz,
  last_test   jsonb,
  updated_by  uuid REFERENCES users(id),
  updated_at  timestamptz NOT NULL DEFAULT now()
);
CREATE TRIGGER integrations_updated BEFORE UPDATE ON integrations FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER integrations_audit AFTER INSERT OR UPDATE OR DELETE ON integrations FOR EACH ROW EXECUTE FUNCTION audit_row_change();

-- Pembeli tamu: akun tanpa nomor HP / Telegram / Facebook, hanya hidup di satu browser (sesi). Tidak boleh membuka toko.
ALTER TABLE users ADD COLUMN is_guest boolean NOT NULL DEFAULT false;

-- +goose Down
ALTER TABLE users DROP COLUMN is_guest;
DROP TABLE integrations;
DELETE FROM platform_settings WHERE key = 'sms';
