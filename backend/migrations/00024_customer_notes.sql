-- +goose Up
-- Catatan CS tentang pelanggan (riwayat telepon, keluhan, alasan akun ditangguhkan). Hanya dilihat admin.
CREATE TABLE customer_notes (
  id         bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  user_id    uuid NOT NULL REFERENCES users(id),
  author_id  uuid NOT NULL REFERENCES users(id),
  kind       varchar(12) NOT NULL DEFAULT 'note' CHECK (kind IN ('note','suspended','reactivated')),
  note       text NOT NULL CHECK (length(note) BETWEEN 2 AND 2000),
  created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX customer_notes_user_idx ON customer_notes (user_id, id DESC);
CREATE TRIGGER customer_notes_audit AFTER INSERT OR UPDATE OR DELETE ON customer_notes FOR EACH ROW EXECUTE FUNCTION audit_row_change();

-- Pencarian CS: nomor HP (angka saja) dan nama.
CREATE INDEX IF NOT EXISTS users_phone_digits_idx ON users ((regexp_replace(coalesce(phone, ''), '\D', '', 'g')) text_pattern_ops);
CREATE INDEX IF NOT EXISTS payment_events_tran_idx ON payment_events (tran_id, created_at DESC);

-- +goose Down
DROP INDEX IF EXISTS payment_events_tran_idx;
DROP INDEX IF EXISTS users_phone_digits_idx;
DROP TABLE customer_notes;
