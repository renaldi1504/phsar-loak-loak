-- +goose Up
-- Batas waktu pesanan yang diatur super admin + jam buka toko.
--   dibayar ──(accept_minutes, menit buka)──► harus "Terima" (respond_by), kalau tidak → batal + refund
--   diterima ──(pack_minutes, menit buka)──► harus "Siap dijemput" (ready_by), kalau tidak → batal + refund
-- Kalau business_hours aktif, hitungan berhenti saat toko tutup (jam toko, atau jam standar kalau belum diisi).
ALTER TABLE stores ADD COLUMN opening_hours jsonb; -- [7] Minggu..Sabtu: null = tutup | {"open":"08:00","close":"22:00"}
ALTER TABLE orders ADD COLUMN ready_by timestamptz;
CREATE INDEX orders_ready_by_idx ON orders (ready_by) WHERE status = 'accepted';
INSERT INTO platform_settings (key, value) VALUES ('order_timing', '{"accept_minutes": 120, "pack_enabled": true, "pack_minutes": 720, "business_hours": true,
  "default_hours": [{"open":"08:00","close":"22:00"},{"open":"08:00","close":"22:00"},{"open":"08:00","close":"22:00"},{"open":"08:00","close":"22:00"},
                    {"open":"08:00","close":"22:00"},{"open":"08:00","close":"22:00"},{"open":"08:00","close":"22:00"}]}')
ON CONFLICT (key) DO NOTHING;

-- +goose Down
DELETE FROM platform_settings WHERE key = 'order_timing';
DROP INDEX orders_ready_by_idx;
ALTER TABLE orders DROP COLUMN ready_by;
ALTER TABLE stores DROP COLUMN opening_hours;
