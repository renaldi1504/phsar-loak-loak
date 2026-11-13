-- +goose Up
-- Akun: satu user untuk pembeli, seller, kurir, dan admin. Login Telegram · Facebook · SMS · Apple.
CREATE TYPE auth_provider AS ENUM ('telegram', 'facebook', 'sms', 'apple');

CREATE TABLE users (
  id                      uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name                    varchar(80) NOT NULL DEFAULT '',
  email                   citext,
  phone                   varchar(30) CHECK (phone IS NULL OR phone ~ '^\+[1-9][0-9]{6,14}$'),  -- E.164
  language                lang_code NOT NULL DEFAULT 'en',
  avatar_url              text,
  personalization_enabled boolean NOT NULL DEFAULT true,
  status                  varchar(16) NOT NULL DEFAULT 'active' CHECK (status IN ('active','suspended','deleted')),
  created_at              timestamptz NOT NULL DEFAULT now(),
  updated_at              timestamptz NOT NULL DEFAULT now(),
  deleted_at              timestamptz,
  deleted_by              uuid REFERENCES users(id)
);
CREATE UNIQUE INDEX users_email_active ON users (email) WHERE email IS NOT NULL AND deleted_at IS NULL;
CREATE UNIQUE INDEX users_phone_active ON users (phone) WHERE phone IS NOT NULL AND deleted_at IS NULL;
CREATE TRIGGER users_updated BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER users_audit AFTER INSERT OR UPDATE OR DELETE ON users FOR EACH ROW EXECUTE FUNCTION audit_row_change();

CREATE TABLE auth_identities (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id      uuid NOT NULL REFERENCES users(id),
  provider     auth_provider NOT NULL,
  provider_uid varchar(128) NOT NULL,   -- id Telegram/Facebook/Apple, atau nomor E.164 untuk sms
  display      varchar(128),            -- mis. @username Telegram, untuk ditampilkan
  verified_at  timestamptz,
  last_used_at timestamptz,
  created_at   timestamptz NOT NULL DEFAULT now(),
  UNIQUE (provider, provider_uid)
);
CREATE INDEX auth_identities_user_idx ON auth_identities (user_id);

CREATE TABLE otp_codes (
  id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  phone      varchar(30) NOT NULL,
  code_hash  char(64) NOT NULL,                     -- sha256, kode tidak disimpan mentah
  purpose    varchar(32) NOT NULL CHECK (purpose IN ('login','change_phone','payout_account')),
  channel    varchar(16) NOT NULL CHECK (channel IN ('sms','telegram')),
  attempts   smallint NOT NULL DEFAULT 0 CHECK (attempts <= 5),
  expires_at timestamptz NOT NULL,                  -- 5 menit
  used_at    timestamptz,
  created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX otp_codes_phone_idx ON otp_codes (phone, created_at DESC);

CREATE TABLE sessions (
  id                 uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id            uuid NOT NULL REFERENCES users(id),
  audience           varchar(16) NOT NULL CHECK (audience IN ('app','seller_web','admin')),
  refresh_token_hash char(64) NOT NULL UNIQUE,      -- dirotasi setiap dipakai
  device_name        varchar(120),
  ip                 inet,
  user_agent         text,
  created_at         timestamptz NOT NULL DEFAULT now(),
  last_seen_at       timestamptz NOT NULL DEFAULT now(),
  expires_at         timestamptz NOT NULL,
  revoked_at         timestamptz
);
CREATE INDEX sessions_user_idx ON sessions (user_id) WHERE revoked_at IS NULL;

CREATE TABLE login_attempts (
  id          bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  identifier  varchar(128) NOT NULL,                -- phone / telegram id / email
  ip          inet,
  device_hash varchar(64),
  success     boolean NOT NULL,
  reason      varchar(64),
  occurred_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX login_attempts_ident_idx ON login_attempts (identifier, occurred_at DESC);
CREATE INDEX login_attempts_ip_idx ON login_attempts (ip, occurred_at DESC);

-- Buku alamat pembeli. Titik peta wajib karena kurir Phsar mengantar ke titik ini.
CREATE TABLE addresses (
  id                uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id           uuid NOT NULL REFERENCES users(id),
  label             varchar(40),                    -- Rumah, Kantor
  recipient         varchar(80) NOT NULL,
  phone             varchar(30) NOT NULL,
  address_line      varchar(300) NOT NULL,
  district          varchar(80),
  province          varchar(80),
  service_area_code varchar(32),                    -- diisi setelah dicek terhadap service_areas
  lat               numeric(9,6) NOT NULL CHECK (lat BETWEEN -90 AND 90),
  lng               numeric(9,6) NOT NULL CHECK (lng BETWEEN -180 AND 180),
  note              text,                           -- catatan untuk kurir
  is_default        boolean NOT NULL DEFAULT false,
  created_at        timestamptz NOT NULL DEFAULT now(),
  updated_at        timestamptz NOT NULL DEFAULT now(),
  deleted_at        timestamptz
);
CREATE INDEX addresses_user_idx ON addresses (user_id) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX addresses_one_default ON addresses (user_id) WHERE is_default AND deleted_at IS NULL;
CREATE TRIGGER addresses_updated BEFORE UPDATE ON addresses FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Tour cara pakai (app & Seller Center).
CREATE TABLE user_tour_progress (
  user_id      uuid NOT NULL REFERENCES users(id),
  tour_key     varchar(48) NOT NULL,                -- buyer_first_open · seller_first_product · …
  version      smallint NOT NULL DEFAULT 1,
  surface      varchar(16) NOT NULL CHECK (surface IN ('app','seller_web','seller_app')),
  last_step    smallint NOT NULL DEFAULT 0,
  completed_at timestamptz,
  dismissed_at timestamptz,
  PRIMARY KEY (user_id, tour_key, version)
);

-- +goose Down
DROP TABLE user_tour_progress;
DROP TABLE addresses;
DROP TABLE login_attempts;
DROP TABLE sessions;
DROP TABLE otp_codes;
DROP TABLE auth_identities;
DROP TABLE users;
DROP TYPE auth_provider;
