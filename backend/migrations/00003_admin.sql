-- +goose Up
-- Admin Phsar: peran, hak akses, dan staf. Admin login di host terpisah + MFA.
CREATE TABLE admin_roles (
  code      varchar(32) PRIMARY KEY,               -- super_admin · moderator · cs · finance · ads_reviewer
  name_i18n jsonb NOT NULL
);

CREATE TABLE role_permissions (
  role_code  varchar(32) NOT NULL REFERENCES admin_roles(code) ON DELETE CASCADE,
  permission varchar(64) NOT NULL,                 -- product.moderate · seller.verify · payout.approve · …
  PRIMARY KEY (role_code, permission)
);

CREATE TABLE admin_users (
  user_id        uuid PRIMARY KEY REFERENCES users(id),
  role_code      varchar(32) NOT NULL REFERENCES admin_roles(code),
  mfa_enabled    boolean NOT NULL DEFAULT false,
  mfa_secret_enc bytea,                            -- TOTP secret, dienkripsi
  active         boolean NOT NULL DEFAULT true,
  created_at     timestamptz NOT NULL DEFAULT now(),
  updated_at     timestamptz NOT NULL DEFAULT now(),
  deleted_at     timestamptz,
  -- Peran finance dan super_admin wajib MFA.
  CONSTRAINT admin_mfa_required CHECK (role_code NOT IN ('finance','super_admin') OR mfa_enabled OR NOT active)
);
CREATE TRIGGER admin_users_updated BEFORE UPDATE ON admin_users FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER admin_users_audit AFTER INSERT OR UPDATE OR DELETE ON admin_users FOR EACH ROW EXECUTE FUNCTION audit_row_change();
CREATE TRIGGER role_permissions_audit AFTER INSERT OR UPDATE OR DELETE ON role_permissions FOR EACH ROW EXECUTE FUNCTION audit_row_change();

-- +goose Down
DROP TABLE admin_users;
DROP TABLE role_permissions;
DROP TABLE admin_roles;
