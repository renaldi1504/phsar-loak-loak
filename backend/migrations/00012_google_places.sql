-- +goose Up
-- Lokasi dipilih lewat Google Places: simpan place_id sebagai rujukan (lat/lng tetap sumber utama).
ALTER TABLE store_pickup_points ADD COLUMN google_place_id varchar(300);
ALTER TABLE addresses ADD COLUMN google_place_id varchar(300);

-- +goose Down
ALTER TABLE addresses DROP COLUMN google_place_id;
ALTER TABLE store_pickup_points DROP COLUMN google_place_id;
