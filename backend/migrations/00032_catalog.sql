-- +goose Up
-- Kategori & merek dikelola admin (menu Categories & brands).
-- Kategori "disembunyikan" = categories.deleted_at diisi: tidak tampil di website/app dan tidak bisa dipilih seller,
-- produk lama tetap ada. Hapus permanen hanya untuk kategori tanpa produk.

ALTER TABLE brands
  ADD COLUMN verified   boolean NOT NULL DEFAULT false,   -- merek resmi / distributor terverifikasi
  ADD COLUMN active     boolean NOT NULL DEFAULT true,
  ADD COLUMN updated_at timestamptz NOT NULL DEFAULT now();
CREATE TRIGGER brands_updated BEFORE UPDATE ON brands FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Merek ditawarkan di kategori mana (kosong = semua kategori).
CREATE TABLE brand_categories (
  brand_id    integer NOT NULL REFERENCES brands(id) ON DELETE CASCADE,
  category_id varchar(60) NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
  PRIMARY KEY (brand_id, category_id)
);
CREATE INDEX brand_categories_category_idx ON brand_categories (category_id);

ALTER TABLE category_attributes ADD COLUMN sort_order integer NOT NULL DEFAULT 0;
ALTER TABLE prohibited_keywords ADD COLUMN created_at timestamptz NOT NULL DEFAULT now();

-- Jejak audit untuk semua tabel yang kini bisa diubah admin.
CREATE OR REPLACE TRIGGER brands_audit AFTER INSERT OR UPDATE OR DELETE ON brands FOR EACH ROW EXECUTE FUNCTION audit_row_change();
CREATE OR REPLACE TRIGGER category_translations_audit AFTER INSERT OR UPDATE OR DELETE ON category_translations FOR EACH ROW EXECUTE FUNCTION audit_row_change();
CREATE OR REPLACE TRIGGER category_attributes_audit AFTER INSERT OR UPDATE OR DELETE ON category_attributes FOR EACH ROW EXECUTE FUNCTION audit_row_change();
CREATE OR REPLACE TRIGGER prohibited_keywords_audit AFTER INSERT OR UPDATE OR DELETE ON prohibited_keywords FOR EACH ROW EXECUTE FUNCTION audit_row_change();

-- +goose Down
DROP TRIGGER IF EXISTS prohibited_keywords_audit ON prohibited_keywords;
DROP TRIGGER IF EXISTS category_attributes_audit ON category_attributes;
DROP TRIGGER IF EXISTS category_translations_audit ON category_translations;
DROP TRIGGER IF EXISTS brands_audit ON brands;
ALTER TABLE prohibited_keywords DROP COLUMN created_at;
ALTER TABLE category_attributes DROP COLUMN sort_order;
DROP TABLE brand_categories;
DROP TRIGGER brands_updated ON brands;
ALTER TABLE brands DROP COLUMN updated_at, DROP COLUMN active, DROP COLUMN verified;
