-- +goose Up
-- Daftar bertingkat di Seller Center: pesanan urut tanggal bayar, produk urut terakhir diubah.
CREATE INDEX IF NOT EXISTS orders_store_paid_idx ON orders (store_id, paid_at DESC) WHERE paid_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS products_store_updated_idx ON products (store_id, updated_at DESC) WHERE deleted_at IS NULL;

-- +goose Down
DROP INDEX IF EXISTS orders_store_paid_idx;
DROP INDEX IF EXISTS products_store_updated_idx;
