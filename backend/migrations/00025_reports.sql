-- +goose Up
-- Laporan pembeli (produk/toko) dan strike toko. 3 strike aktif (90 hari) → toko otomatis ditangguhkan.

CREATE TABLE reports (
  id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  reporter_id     uuid NOT NULL REFERENCES users(id),
  target_type     varchar(8) NOT NULL CHECK (target_type IN ('product','store')),
  product_id      uuid REFERENCES products(id),
  store_id        uuid NOT NULL REFERENCES stores(id),       -- untuk laporan produk: toko pemilik produk
  reason          varchar(16) NOT NULL CHECK (reason IN ('counterfeit','prohibited','misleading','scam','offensive','other')),
  note            text CHECK (note IS NULL OR length(note) <= 1000),
  status          varchar(10) NOT NULL DEFAULT 'open' CHECK (status IN ('open','dismissed','actioned')),
  resolution      varchar(16) CHECK (resolution IN ('dismiss','ban_product','strike')),
  resolution_note text,
  resolved_by     uuid REFERENCES users(id),
  resolved_at     timestamptz,
  created_at      timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT reports_target CHECK ((target_type = 'product') = (product_id IS NOT NULL))
);
-- satu laporan terbuka per pelapor per target
CREATE UNIQUE INDEX reports_one_open ON reports (reporter_id, target_type, coalesce(product_id, store_id)) WHERE status = 'open';
CREATE INDEX reports_open_idx ON reports (status, created_at);
CREATE INDEX reports_store_idx ON reports (store_id, created_at DESC);
CREATE INDEX reports_reporter_idx ON reports (reporter_id, created_at DESC);
CREATE TRIGGER reports_audit AFTER INSERT OR UPDATE OR DELETE ON reports FOR EACH ROW EXECUTE FUNCTION audit_row_change();

CREATE TABLE store_strikes (
  id                uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  store_id          uuid NOT NULL REFERENCES stores(id),
  reason            varchar(16) NOT NULL CHECK (reason IN ('counterfeit','prohibited','misleading','scam','offensive','other')),
  note              text NOT NULL CHECK (length(note) BETWEEN 5 AND 1000),   -- dilihat seller
  points            smallint NOT NULL DEFAULT 1 CHECK (points BETWEEN 1 AND 3),
  product_id        uuid REFERENCES products(id),
  caused_suspension boolean NOT NULL DEFAULT false,
  created_by        uuid NOT NULL REFERENCES users(id),
  created_at        timestamptz NOT NULL DEFAULT now(),
  expires_at        timestamptz NOT NULL DEFAULT now() + interval '90 days',
  revoked_at        timestamptz,
  revoked_by        uuid REFERENCES users(id),
  revoke_note       text
);
CREATE INDEX store_strikes_store_idx ON store_strikes (store_id, created_at DESC);
CREATE TRIGGER store_strikes_audit AFTER INSERT OR UPDATE OR DELETE ON store_strikes FOR EACH ROW EXECUTE FUNCTION audit_row_change();

-- +goose Down
DROP TABLE store_strikes;
DROP TABLE reports;
