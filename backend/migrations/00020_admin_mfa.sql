-- +goose Up
-- Login admin + MFA (TOTP, aplikasi authenticator). Semua peran admin wajib MFA saat login
-- (dicek di backend: API admin hanya menerima sesi yang sudah lolos MFA). Aturan lama
-- "finance/super_admin harus mfa_enabled" diganti karena admin baru baru bisa memasang MFA
-- setelah diberi peran.
ALTER TABLE admin_users DROP CONSTRAINT admin_mfa_required;
ALTER TABLE admin_users
  ADD COLUMN mfa_pending_enc bytea,        -- secret yang sedang dipasang (belum dikonfirmasi kode)
  ADD COLUMN mfa_last_step   bigint,       -- langkah 30 detik terakhir yang dipakai (kode tidak bisa dipakai ulang)
  ADD COLUMN mfa_enabled_at  timestamptz,
  ADD COLUMN last_login_at   timestamptz;

ALTER TABLE sessions
  ADD COLUMN mfa_verified_at timestamptz,  -- sesi admin baru boleh memakai API setelah ini terisi
  ADD COLUMN mfa_failures    int NOT NULL DEFAULT 0;

CREATE TABLE admin_mfa_backup_codes (
  id         bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  user_id    uuid NOT NULL REFERENCES users(id),
  code_hash  char(64) NOT NULL,           -- sha256, kode asli hanya ditampilkan sekali
  used_at    timestamptz,
  created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX admin_mfa_backup_codes_user ON admin_mfa_backup_codes (user_id) WHERE used_at IS NULL;

-- +goose Down
DROP TABLE admin_mfa_backup_codes;
ALTER TABLE sessions DROP COLUMN mfa_failures, DROP COLUMN mfa_verified_at;
ALTER TABLE admin_users DROP COLUMN last_login_at, DROP COLUMN mfa_enabled_at, DROP COLUMN mfa_last_step, DROP COLUMN mfa_pending_enc;
ALTER TABLE admin_users ADD CONSTRAINT admin_mfa_required
  CHECK (role_code NOT IN ('finance','super_admin') OR mfa_enabled OR NOT active) NOT VALID;
