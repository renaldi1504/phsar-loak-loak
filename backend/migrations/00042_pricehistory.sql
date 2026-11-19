-- +goose Up
-- Riwayat harga varian (harga normal, bukan harga flash/promo — itu ada di flash_items).
--   setiap INSERT / perubahan price atau compare_at_price di product_variants → baris lama ditutup (valid_to), baris baru dibuka.
--   Diisi trigger database, jadi semua jalur (buat/edit produk, PATCH harga-stok, impor, SQL manual) tercatat.
-- Dipakai aturan flash sale: harga flash harus di bawah harga terendah 30 hari sebelum sesi mulai
-- (platform_settings['flash_rules'] = {"low30_enabled": true}, diatur super admin; default nyala).
CREATE TABLE variant_price_history (
  id               bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  variant_id       uuid NOT NULL REFERENCES product_variants(id) ON DELETE CASCADE,
  price            numeric(12,2) NOT NULL,
  compare_at_price numeric(12,2),
  valid_from       timestamptz NOT NULL,
  valid_to         timestamptz, -- NULL = harga sekarang
  CONSTRAINT variant_price_history_range CHECK (valid_to IS NULL OR valid_to >= valid_from)
);
CREATE UNIQUE INDEX variant_price_history_open ON variant_price_history (variant_id) WHERE valid_to IS NULL;
CREATE INDEX variant_price_history_variant ON variant_price_history (variant_id, valid_from);

-- +goose StatementBegin
CREATE FUNCTION record_variant_price() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
  IF TG_OP = 'UPDATE' AND NEW.price IS NOT DISTINCT FROM OLD.price AND NEW.compare_at_price IS NOT DISTINCT FROM OLD.compare_at_price THEN
    RETURN NULL;
  END IF;
  -- diubah lagi dalam transaksi yang sama: perbarui baris yang baru dibuka (tanpa baris berdurasi nol)
  UPDATE variant_price_history SET price = NEW.price, compare_at_price = NEW.compare_at_price
    WHERE variant_id = NEW.id AND valid_to IS NULL AND valid_from = now();
  IF FOUND THEN
    RETURN NULL;
  END IF;
  UPDATE variant_price_history SET valid_to = greatest(now(), valid_from) WHERE variant_id = NEW.id AND valid_to IS NULL;
  INSERT INTO variant_price_history (variant_id, price, compare_at_price, valid_from) VALUES (NEW.id, NEW.price, NEW.compare_at_price, now());
  RETURN NULL;
END $$;
-- +goose StatementEnd
CREATE TRIGGER product_variants_price_history AFTER INSERT OR UPDATE OF price, compare_at_price ON product_variants
  FOR EACH ROW EXECUTE FUNCTION record_variant_price();

-- Harga normal terendah varian dalam p_days hari sebelum p_at (rentang [p_at − p_days, p_at]).
-- Varian yang belum ada pada p_at → harga pertamanya; tanpa riwayat sama sekali → harga sekarang.
-- +goose StatementBegin
CREATE FUNCTION variant_low_price(p_variant uuid, p_at timestamptz, p_days integer DEFAULT 30) RETURNS numeric LANGUAGE sql STABLE AS $$
  SELECT coalesce(
    (SELECT min(h.price) FROM variant_price_history h
      WHERE h.variant_id = p_variant AND h.valid_from <= p_at
        AND (h.valid_to IS NULL OR (h.valid_to > p_at - make_interval(days => p_days) AND h.valid_to > h.valid_from))),
    (SELECT h.price FROM variant_price_history h WHERE h.variant_id = p_variant ORDER BY h.valid_from, h.id LIMIT 1),
    (SELECT v.price FROM product_variants v WHERE v.id = p_variant))
$$;
-- +goose StatementEnd

-- Backfill: harga sekarang berlaku sejak varian dibuat.
INSERT INTO variant_price_history (variant_id, price, compare_at_price, valid_from)
SELECT id, price, compare_at_price, created_at FROM product_variants;

INSERT INTO platform_settings (key, value) VALUES ('flash_rules', '{"low30_enabled": true}') ON CONFLICT (key) DO NOTHING;

-- +goose Down
DELETE FROM platform_settings WHERE key = 'flash_rules';
DROP FUNCTION variant_low_price(uuid, timestamptz, integer);
DROP TRIGGER product_variants_price_history ON product_variants;
DROP FUNCTION record_variant_price();
DROP TABLE variant_price_history;
