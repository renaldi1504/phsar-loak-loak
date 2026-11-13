-- +goose Up
-- Login Telegram lewat deep link bot (cocok untuk app native, tanpa perlu domain):
-- app minta login → buka t.me/<bot>?start=login_<nonce> → user tekan Start → bot mengonfirmasi
-- → app (polling) menerima sesi.
CREATE TABLE telegram_login_requests (
  id               uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  nonce_hash       char(64) NOT NULL UNIQUE,          -- sha256 dari nonce di link /start
  poll_token_hash  char(64) NOT NULL,                 -- sha256 dari token rahasia milik app
  audience         varchar(16) NOT NULL DEFAULT 'app' CHECK (audience IN ('app','seller_web')),
  status           varchar(16) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','confirmed','consumed','expired')),
  telegram_user_id bigint,
  telegram_username varchar(64),
  first_name       varchar(80),
  last_name        varchar(80),
  language_code    varchar(16),
  device_name      varchar(120),
  ip               inet,
  created_at       timestamptz NOT NULL DEFAULT now(),
  expires_at       timestamptz NOT NULL,              -- 10 menit
  confirmed_at     timestamptz,
  consumed_at      timestamptz
);
CREATE INDEX telegram_login_requests_pending ON telegram_login_requests (expires_at) WHERE status = 'pending';

-- Posisi update terakhir dari Telegram (mode polling untuk pengembangan lokal).
CREATE TABLE telegram_bot_state (
  id             smallint PRIMARY KEY DEFAULT 1 CHECK (id = 1),
  last_update_id bigint NOT NULL DEFAULT 0,
  updated_at     timestamptz NOT NULL DEFAULT now()
);
INSERT INTO telegram_bot_state (id) VALUES (1);

-- +goose Down
DROP TABLE telegram_bot_state;
DROP TABLE telegram_login_requests;
