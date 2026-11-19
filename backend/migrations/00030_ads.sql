-- +goose Up
-- Iklan produk bersponsor (bayar per klik dari saldo seller).
--   seller (pemilik) buat kampanye → admin tinjau → tayang di hasil cari & halaman kategori (maks. N slot)
--   klik lewat /ad/{id} → potong saldo seller (seller_ledger kind ad_spend), sekali per pengunjung per kampanye per hari.

CREATE TABLE ad_campaigns (
  id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  store_id      uuid NOT NULL REFERENCES stores(id),
  product_id    uuid NOT NULL REFERENCES products(id),
  status        varchar(10) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','approved','rejected','paused','ended')),
  paused_by     varchar(6) CHECK (paused_by IN ('seller','admin')),
  bid           numeric(6,2) NOT NULL CHECK (bid BETWEEN 0.05 AND 2.00),        -- biaya per klik (USD)
  daily_budget  numeric(8,2) NOT NULL CHECK (daily_budget >= 1),
  starts_on     date NOT NULL,
  ends_on       date CHECK (ends_on IS NULL OR ends_on >= starts_on),
  review_note   varchar(300),
  reviewed_by   uuid REFERENCES users(id),
  reviewed_at   timestamptz,
  created_by    uuid REFERENCES users(id),
  created_at    timestamptz NOT NULL DEFAULT now(),
  updated_at    timestamptz NOT NULL DEFAULT now()
);
-- satu kampanye berjalan per produk
CREATE UNIQUE INDEX ad_campaigns_one_per_product ON ad_campaigns (product_id) WHERE status IN ('pending','approved','paused');
CREATE INDEX ad_campaigns_store_idx ON ad_campaigns (store_id, created_at DESC);
CREATE INDEX ad_campaigns_live_idx ON ad_campaigns (status, starts_on) WHERE status = 'approved';
CREATE TRIGGER ad_campaigns_updated BEFORE UPDATE ON ad_campaigns FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER ad_campaigns_audit AFTER INSERT OR UPDATE OR DELETE ON ad_campaigns FOR EACH ROW EXECUTE FUNCTION audit_row_change();

-- Statistik harian (hari = tanggal Phnom Penh).
CREATE TABLE ad_stats_daily (
  campaign_id uuid NOT NULL REFERENCES ad_campaigns(id),
  day         date NOT NULL,
  impressions integer NOT NULL DEFAULT 0,
  clicks      integer NOT NULL DEFAULT 0,
  spend       numeric(10,2) NOT NULL DEFAULT 0,
  PRIMARY KEY (campaign_id, day)
);

-- Klik yang ditagih (satu per pengunjung per kampanye per hari). Hanya hash, tanpa IP mentah.
CREATE TABLE ad_clicks (
  id           bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  campaign_id  uuid NOT NULL REFERENCES ad_campaigns(id),
  day          date NOT NULL,
  visitor_hash char(64) NOT NULL,
  ip_hash      char(64) NOT NULL,
  charged      numeric(6,2) NOT NULL,
  created_at   timestamptz NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX ad_clicks_once ON ad_clicks (campaign_id, day, visitor_hash);
CREATE INDEX ad_clicks_ip_idx ON ad_clicks (campaign_id, day, ip_hash);

ALTER TABLE seller_ledger DROP CONSTRAINT seller_ledger_kind_check;
ALTER TABLE seller_ledger ADD CONSTRAINT seller_ledger_kind_check
  CHECK (kind IN ('order_earning','payout','payout_cancelled','adjustment','return_refund','ad_spend'));
ALTER TABLE seller_ledger ADD COLUMN ad_campaign_id uuid REFERENCES ad_campaigns(id);

INSERT INTO platform_settings (key, value) VALUES ('ads', '{"enabled": true, "slots": 2}') ON CONFLICT (key) DO NOTHING;

-- +goose Down
DELETE FROM platform_settings WHERE key = 'ads';
DELETE FROM seller_ledger WHERE kind = 'ad_spend';
ALTER TABLE seller_ledger DROP COLUMN ad_campaign_id;
ALTER TABLE seller_ledger DROP CONSTRAINT seller_ledger_kind_check;
ALTER TABLE seller_ledger ADD CONSTRAINT seller_ledger_kind_check
  CHECK (kind IN ('order_earning','payout','payout_cancelled','adjustment','return_refund'));
DROP TABLE ad_clicks;
DROP TABLE ad_stats_daily;
DROP TABLE ad_campaigns;
