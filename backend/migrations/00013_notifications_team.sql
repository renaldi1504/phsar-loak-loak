-- +goose Up
-- Notifikasi untuk seller (di Seller Center + dikirim lewat Telegram) dan undangan staf toko.

-- Satu baris per penerima. Teks dirender saat ditampilkan/dikirim sesuai bahasa user (kind + params).
-- Kolom telegram_* = outbox: worker mengirim, mencoba ulang kalau gagal.
CREATE TABLE notifications (
  id                bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  user_id           uuid NOT NULL REFERENCES users(id),
  store_id          uuid REFERENCES stores(id),
  kind              varchar(40) NOT NULL,              -- store.approved · product.changes_requested · …
  params            jsonb NOT NULL DEFAULT '{}',       -- {"store":"…","product":"…","note":"…"}
  link              varchar(300),                      -- hash Seller Center, mis. #/store/<id>/product/<id>
  read_at           timestamptz,
  telegram_status   varchar(12) NOT NULL DEFAULT 'pending'
                    CHECK (telegram_status IN ('pending','sending','sent','failed','skipped')),
  telegram_attempts integer NOT NULL DEFAULT 0,
  telegram_error    text,
  telegram_sent_at  timestamptz,
  next_attempt_at   timestamptz,
  created_at        timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX notifications_user_idx ON notifications (user_id, created_at DESC);
CREATE INDEX notifications_unread_idx ON notifications (user_id) WHERE read_at IS NULL;
CREATE INDEX notifications_outbox_idx ON notifications (created_at) WHERE telegram_status = 'pending';

-- Undangan staf: owner membuat link, calon staf login lalu menerima. Token disimpan sebagai hash.
CREATE TABLE store_invites (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  store_id    uuid NOT NULL REFERENCES stores(id),
  token_hash  char(64) NOT NULL UNIQUE,
  role        varchar(16) NOT NULL DEFAULT 'staff' CHECK (role = 'staff'),
  created_by  uuid NOT NULL REFERENCES users(id),
  expires_at  timestamptz NOT NULL,
  accepted_by uuid REFERENCES users(id),
  accepted_at timestamptz,
  revoked_at  timestamptz,
  created_at  timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX store_invites_store_idx ON store_invites (store_id, created_at DESC);
CREATE TRIGGER store_invites_audit AFTER INSERT OR UPDATE OR DELETE ON store_invites FOR EACH ROW EXECUTE FUNCTION audit_row_change();

-- +goose Down
DROP TABLE store_invites;
DROP TABLE notifications;
