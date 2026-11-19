-- +goose Up
-- Login Google (website + aplikasi) dan Sign in with Apple (aplikasi iPhone). Kredensial di tabel integrations
-- (diatur super admin), identitas di auth_identities.provider = 'google' / 'apple'.
ALTER TYPE auth_provider ADD VALUE IF NOT EXISTS 'google';
ALTER TABLE integrations DROP CONSTRAINT integrations_code_check;
ALTER TABLE integrations ADD CONSTRAINT integrations_code_check
  CHECK (code IN ('plasgate', 'facebook', 'anthropic', 'gemini', 'google_maps', 'google', 'apple'));

-- +goose Down
DELETE FROM integrations WHERE code IN ('google', 'apple');
ALTER TABLE integrations DROP CONSTRAINT integrations_code_check;
ALTER TABLE integrations ADD CONSTRAINT integrations_code_check CHECK (code IN ('plasgate', 'facebook', 'anthropic', 'gemini', 'google_maps'));
-- Nilai enum 'google' tidak bisa dihapus dari auth_provider tanpa membuat ulang tipe; dibiarkan (tidak berbahaya).
