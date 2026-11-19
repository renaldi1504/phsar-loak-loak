-- +goose Up
-- Pemakaian layanan pihak ketiga (AI & SMS) untuk laporan + pengingat saldo di Admin → Integrasi → Pemakaian & saldo.
-- Biaya adalah PERKIRAAN dari harga yang diisi super admin (platform_settings 'usage_budget'), bukan tagihan resmi.
CREATE TABLE integration_usage (
  id            bigserial PRIMARY KEY,
  provider      varchar(16) NOT NULL CHECK (provider IN ('anthropic', 'gemini', 'plasgate')),
  purpose       varchar(16) NOT NULL CHECK (purpose IN ('translate', 'askbot', 'test', 'sms')),
  ok            boolean NOT NULL,
  error_kind    varchar(16) CHECK (error_kind IN ('billing', 'quota', 'auth', 'other')),
  error_text    varchar(300),
  model         varchar(80),
  input_tokens  integer NOT NULL DEFAULT 0,
  output_tokens integer NOT NULL DEFAULT 0,
  units         integer NOT NULL DEFAULT 0, -- SMS terkirim
  cost_usd      numeric(12,6) NOT NULL DEFAULT 0,
  created_at    timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX integration_usage_provider_idx ON integration_usage (provider, created_at DESC);

-- +goose Down
DROP TABLE integration_usage;
DELETE FROM platform_settings WHERE key = 'usage_budget';
