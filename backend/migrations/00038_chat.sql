-- +goose Up
-- chat: percakapan pembeli ↔ toko. Satu percakapan per (pembeli, toko); pesan boleh membawa kartu produk / pesanan.
-- Penghitung belum dibaca per sisi disimpan di percakapan (pengirim = 0, lawan +1; dinolkan saat dibuka).
-- *_notified_at: pembatas notifikasi (maks. 1 per percakapan per 10 menit per sisi).
CREATE TABLE chat_conversations (
  id                 uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  buyer_id           uuid NOT NULL REFERENCES users(id),
  store_id           uuid NOT NULL REFERENCES stores(id),
  product_id         uuid REFERENCES products(id),     -- kartu produk terakhir (untuk saran balasan)
  order_id           uuid REFERENCES orders(id),       -- kartu pesanan terakhir
  last_message_at    timestamptz,                      -- NULL = belum ada pesan (tidak tampil di daftar)
  last_sender        varchar(6) CHECK (last_sender IN ('buyer','seller')),
  last_preview       varchar(160) NOT NULL DEFAULT '',
  buyer_unread       integer NOT NULL DEFAULT 0 CHECK (buyer_unread >= 0),
  seller_unread      integer NOT NULL DEFAULT 0 CHECK (seller_unread >= 0),
  buyer_notified_at  timestamptz,
  seller_notified_at timestamptz,
  created_at         timestamptz NOT NULL DEFAULT now(),
  updated_at         timestamptz NOT NULL DEFAULT now(),
  UNIQUE (buyer_id, store_id)
);
CREATE INDEX chat_conversations_store_idx ON chat_conversations (store_id, last_message_at DESC) WHERE last_message_at IS NOT NULL;
CREATE INDEX chat_conversations_buyer_idx ON chat_conversations (buyer_id, last_message_at DESC) WHERE last_message_at IS NOT NULL;
CREATE INDEX chat_conversations_waiting_idx ON chat_conversations (store_id) WHERE last_sender = 'buyer';
CREATE TRIGGER chat_conversations_updated BEFORE UPDATE ON chat_conversations FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE chat_messages (
  id              bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  conversation_id uuid NOT NULL REFERENCES chat_conversations(id),
  sender_type     varchar(6) NOT NULL CHECK (sender_type IN ('buyer','seller')),
  sender_id       uuid NOT NULL REFERENCES users(id),
  body            text NOT NULL DEFAULT '' CHECK (length(body) <= 2000),
  photo_key       text,                                -- upload chat_photo (privat)
  product_id      uuid REFERENCES products(id),        -- kartu produk di atas pesan
  order_id        uuid REFERENCES orders(id),          -- kartu pesanan di atas pesan
  created_at      timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT chat_messages_not_empty CHECK (length(body) > 0 OR photo_key IS NOT NULL)
);
CREATE INDEX chat_messages_conv_idx ON chat_messages (conversation_id, id);
CREATE INDEX chat_messages_sender_idx ON chat_messages (sender_id, created_at DESC); -- batas kirim per pengguna
CREATE TRIGGER chat_messages_audit AFTER INSERT OR UPDATE OR DELETE ON chat_messages FOR EACH ROW EXECUTE FUNCTION audit_row_change();

-- Foto chat (privat). Constraint jenis upload ditambah tanpa menulis ulang daftarnya (aman digabung dengan migrasi lain).
-- +goose StatementBegin
DO $$
DECLARE def text;
BEGIN
  SELECT pg_get_constraintdef(oid) INTO def FROM pg_constraint WHERE conname = 'uploads_kind_check' AND conrelid = 'uploads'::regclass;
  IF def IS NOT NULL AND position('chat_photo' IN def) = 0 THEN
    EXECUTE 'ALTER TABLE uploads DROP CONSTRAINT uploads_kind_check';
    EXECUTE 'ALTER TABLE uploads ADD CONSTRAINT uploads_kind_check ' || replace(def, 'ARRAY[', 'ARRAY[''chat_photo''::character varying, ');
  END IF;
END $$;
-- +goose StatementEnd

-- +goose Down
DELETE FROM uploads WHERE kind = 'chat_photo';
-- +goose StatementBegin
DO $$
DECLARE def text;
BEGIN
  SELECT pg_get_constraintdef(oid) INTO def FROM pg_constraint WHERE conname = 'uploads_kind_check' AND conrelid = 'uploads'::regclass;
  IF def IS NOT NULL AND position('chat_photo' IN def) > 0 THEN
    EXECUTE 'ALTER TABLE uploads DROP CONSTRAINT uploads_kind_check';
    EXECUTE 'ALTER TABLE uploads ADD CONSTRAINT uploads_kind_check ' || replace(def, '''chat_photo''::character varying, ', '');
  END IF;
END $$;
-- +goose StatementEnd
DELETE FROM notifications WHERE kind IN ('chat.new', 'chat.reply');
DROP TABLE chat_messages;
DROP TABLE chat_conversations;
