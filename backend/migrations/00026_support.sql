-- +goose Up
-- Tiket bantuan dari pembeli (website/app) dan seller (Seller Center) ke tim CS Phsar.
--   open (menunggu CS) ──CS membalas──► pending (menunggu pengguna) ──pengguna membalas──► open
--   └──CS / pengguna menyelesaikan──► solved ──pengguna membalas lagi──► open
CREATE TABLE tickets (
  id                   uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  ticket_no            varchar(20) NOT NULL UNIQUE,
  source               varchar(8) NOT NULL CHECK (source IN ('buyer','seller')),
  requester_id         uuid NOT NULL REFERENCES users(id),
  store_id             uuid REFERENCES stores(id),            -- tiket seller: atas nama toko
  category             varchar(12) NOT NULL CHECK (category IN ('order','payment','return','account','seller','other')),
  subject              varchar(120) NOT NULL CHECK (length(subject) >= 3),
  order_id             uuid REFERENCES orders(id),
  status               varchar(8) NOT NULL DEFAULT 'open' CHECK (status IN ('open','pending','solved')),
  assigned_to          uuid REFERENCES users(id),
  last_user_message_at timestamptz NOT NULL DEFAULT now(),
  last_admin_reply_at  timestamptz,
  solved_at            timestamptz,
  created_at           timestamptz NOT NULL DEFAULT now(),
  updated_at           timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT tickets_seller_store CHECK ((source = 'seller') = (store_id IS NOT NULL))
);
CREATE INDEX tickets_requester_idx ON tickets (requester_id, created_at DESC);
CREATE INDEX tickets_store_idx ON tickets (store_id, created_at DESC) WHERE store_id IS NOT NULL;
CREATE INDEX tickets_status_idx ON tickets (status, last_user_message_at);
CREATE TRIGGER tickets_updated BEFORE UPDATE ON tickets FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER tickets_audit AFTER INSERT OR UPDATE OR DELETE ON tickets FOR EACH ROW EXECUTE FUNCTION audit_row_change();

-- Pesan dalam tiket. internal = catatan CS, tidak pernah dikirim ke pengguna.
CREATE TABLE ticket_messages (
  id          bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  ticket_id   uuid NOT NULL REFERENCES tickets(id),
  author_id   uuid NOT NULL REFERENCES users(id),
  author_type varchar(8) NOT NULL CHECK (author_type IN ('user','admin')),
  internal    boolean NOT NULL DEFAULT false,
  body        text NOT NULL CHECK (length(body) BETWEEN 1 AND 4000),
  created_at  timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT ticket_messages_internal_admin CHECK (NOT internal OR author_type = 'admin')
);
CREATE INDEX ticket_messages_ticket_idx ON ticket_messages (ticket_id, id);
CREATE TRIGGER ticket_messages_audit AFTER INSERT OR UPDATE OR DELETE ON ticket_messages FOR EACH ROW EXECUTE FUNCTION audit_row_change();

-- +goose Down
DROP TABLE ticket_messages;
DROP TABLE tickets;
