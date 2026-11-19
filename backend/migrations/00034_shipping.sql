-- +goose Up
-- Ongkir per zona (dibayar per toko dalam satu checkout) + gratis ongkir di atas nilai tertentu + biaya kurir (laporan)
-- dan kurs riel yang diatur admin (riwayat disimpan; yang terbaru berlaku).
ALTER TABLE shipping_zones
  ADD COLUMN name        jsonb NOT NULL DEFAULT '{}',                                        -- {"en":…,"km":…,"zh":…}
  ADD COLUMN free_over   numeric(12,2) CHECK (free_over IS NULL OR free_over > 0),          -- subtotal toko ≥ ini → gratis
  ADD COLUMN courier_fee numeric(12,2) NOT NULL DEFAULT 0 CHECK (courier_fee >= 0),          -- yang dibayar Phsar ke kurir per antaran
  ADD COLUMN updated_by  uuid REFERENCES users(id);
UPDATE shipping_zones SET name = '{"en":"Sihanoukville","km":"ក្រុងព្រះសីហនុ","zh":"西哈努克城"}' WHERE code = 'sihanoukville' AND name = '{}';
ALTER TABLE shipping_zones ADD CONSTRAINT shipping_zones_fee_max CHECK (fee_per_store <= 50);

-- Kurs: boleh beberapa kali sehari (yang terakhir berlaku: rate_date DESC, id DESC).
ALTER TABLE exchange_rates DROP CONSTRAINT IF EXISTS exchange_rates_rate_date_source_key;
ALTER TABLE exchange_rates ADD COLUMN created_by uuid REFERENCES users(id), ADD COLUMN note varchar(200);
CREATE INDEX IF NOT EXISTS exchange_rates_latest_idx ON exchange_rates (rate_date DESC, id DESC);
CREATE TRIGGER exchange_rates_audit AFTER INSERT OR UPDATE OR DELETE ON exchange_rates FOR EACH ROW EXECUTE FUNCTION audit_row_change();

-- +goose Down
DROP TRIGGER IF EXISTS exchange_rates_audit ON exchange_rates;
DROP INDEX IF EXISTS exchange_rates_latest_idx;
ALTER TABLE exchange_rates DROP COLUMN note, DROP COLUMN created_by;
DELETE FROM exchange_rates a USING exchange_rates b WHERE a.rate_date = b.rate_date AND a.source = b.source AND a.id < b.id;
ALTER TABLE exchange_rates ADD CONSTRAINT exchange_rates_rate_date_source_key UNIQUE (rate_date, source);
ALTER TABLE shipping_zones DROP CONSTRAINT shipping_zones_fee_max;
ALTER TABLE shipping_zones DROP COLUMN updated_by, DROP COLUMN courier_fee, DROP COLUMN free_over, DROP COLUMN name;
