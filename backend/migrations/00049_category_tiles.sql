-- +goose Up
-- Aman dijalankan ulang (IF NOT EXISTS): versi awal file ini sempat bernomor 00046.
-- Tampilan kategori v2: sampul pilihan admin per kategori + kategori musiman (Tahun Baru, Choul Chnam, Pantai, …).
ALTER TABLE categories ADD COLUMN IF NOT EXISTS cover_key text; -- upload "content_image"; kosong = foto produk terlaris otomatis

CREATE TABLE IF NOT EXISTS season_picks (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  title       jsonb NOT NULL DEFAULT '{}',          -- {"en": "New Year 2027", "km": "…", "zh": "…"}
  color       varchar(8) NOT NULL DEFAULT 'pink' CHECK (color IN ('pink', 'blue', 'sand', 'green')),
  source      varchar(12) NOT NULL CHECK (source IN ('category', 'products', 'search')),
  category_id varchar(60) REFERENCES categories(id),
  product_ids uuid[] NOT NULL DEFAULT '{}',
  query       varchar(60) NOT NULL DEFAULT '',     -- source = search: kata pencarian (nama produk/kategori, semua bahasa)
  starts_at   timestamptz,
  ends_at     timestamptz,                         -- kosong = sepanjang tahun
  sort_order  integer NOT NULL DEFAULT 0,
  active      boolean NOT NULL DEFAULT true,
  created_by  uuid REFERENCES users(id),
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now(),
  CHECK (source <> 'category' OR category_id IS NOT NULL),
  CHECK (source <> 'search' OR query <> ''),
  CHECK (ends_at IS NULL OR starts_at IS NULL OR ends_at > starts_at)
);
DROP TRIGGER IF EXISTS season_picks_updated ON season_picks;
CREATE TRIGGER season_picks_updated BEFORE UPDATE ON season_picks FOR EACH ROW EXECUTE FUNCTION set_updated_at();
DROP TRIGGER IF EXISTS season_picks_audit ON season_picks;
CREATE TRIGGER season_picks_audit AFTER INSERT OR UPDATE OR DELETE ON season_picks FOR EACH ROW EXECUTE FUNCTION audit_row_change();

-- +goose Down
DROP TABLE IF EXISTS season_picks;
ALTER TABLE categories DROP COLUMN IF EXISTS cover_key;
