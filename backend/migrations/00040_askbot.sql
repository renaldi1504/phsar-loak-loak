-- +goose Up
-- askbot: bot pembantu seller menjawab pertanyaan pembeli.
--   FAQ toko (opsional per produk) → dipakai bot + tampil di halaman produk ("Tanya jawab").
--   askbot_usage: satu baris per permintaan draf (batas per toko per jam + laporan admin).
--   platform_settings['askbot'] = {"enabled": true, "hourly_limit": 60} (super admin).
CREATE TABLE store_faqs (
  id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  store_id   uuid NOT NULL REFERENCES stores(id),
  product_id uuid REFERENCES products(id) ON DELETE CASCADE, -- null = berlaku untuk seluruh toko
  lang       lang_code NOT NULL,                             -- bahasa yang ditulis seller
  question   varchar(300) NOT NULL CHECK (length(trim(question)) > 0),
  answer     text NOT NULL CHECK (length(trim(answer)) > 0 AND length(answer) <= 2000),
  position   smallint NOT NULL DEFAULT 0,
  created_by uuid REFERENCES users(id),
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX store_faqs_store_idx ON store_faqs (store_id, product_id);
CREATE TRIGGER store_faqs_updated BEFORE UPDATE ON store_faqs FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER store_faqs_audit AFTER INSERT OR UPDATE OR DELETE ON store_faqs FOR EACH ROW EXECUTE FUNCTION audit_row_change();

CREATE TABLE askbot_usage (
  id            bigserial PRIMARY KEY,
  store_id      uuid NOT NULL REFERENCES stores(id),
  user_id       uuid REFERENCES users(id),
  product_id    uuid,
  lang          lang_code NOT NULL,
  mode          varchar(16) NOT NULL CHECK (mode IN ('claude','rules')),
  model         varchar(64),
  input_tokens  integer NOT NULL DEFAULT 0,
  output_tokens integer NOT NULL DEFAULT 0,
  error         text,                                        -- Claude gagal → jatuh ke aturan
  created_at    timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX askbot_usage_store_idx ON askbot_usage (store_id, created_at DESC);
CREATE INDEX askbot_usage_time_idx ON askbot_usage (created_at);

INSERT INTO platform_settings (key, value) VALUES ('askbot', '{"enabled": true, "hourly_limit": 60}') ON CONFLICT (key) DO NOTHING;

-- +goose Down
DELETE FROM platform_settings WHERE key = 'askbot';
DROP TABLE askbot_usage;
DROP TABLE store_faqs;
