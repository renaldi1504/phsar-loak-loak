-- +goose Up
-- Semua kunci pihak ketiga dikumpulkan di halaman Integrasi (super admin): tambah AI (Anthropic utama, Gemini cadangan)
-- dan Google Maps. Rahasia terenkripsi (DATA_ENCRYPTION_KEY); .env tetap jadi cadangan kalau belum diisi di admin.
ALTER TABLE integrations DROP CONSTRAINT integrations_code_check;
ALTER TABLE integrations ADD CONSTRAINT integrations_code_check CHECK (code IN ('plasgate', 'facebook', 'anthropic', 'gemini', 'google_maps'));

-- +goose Down
DELETE FROM integrations WHERE code IN ('anthropic', 'gemini', 'google_maps');
ALTER TABLE integrations DROP CONSTRAINT integrations_code_check;
ALTER TABLE integrations ADD CONSTRAINT integrations_code_check CHECK (code IN ('plasgate', 'facebook'));
