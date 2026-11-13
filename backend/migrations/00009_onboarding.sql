-- +goose Up
-- Onboarding seller: file yang di-upload, catatan review toko.

-- Setiap file yang di-upload dicatat pemiliknya, supaya user tidak bisa memakai file milik orang lain.
CREATE TABLE uploads (
  key           text PRIMARY KEY,                  -- public/stores/… atau private/kyc/…
  owner_user_id uuid NOT NULL REFERENCES users(id),
  kind          varchar(24) NOT NULL CHECK (kind IN ('store_logo','store_banner','pickup_photo','kyc_document','product_media')),
  content_type  varchar(64) NOT NULL,
  bytes         integer NOT NULL CHECK (bytes > 0),
  is_private    boolean NOT NULL,
  created_at    timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX uploads_owner_idx ON uploads (owner_user_id, created_at DESC);

ALTER TABLE stores ADD COLUMN review_note text;         -- alasan ditolak / perlu perbaikan, terlihat oleh seller
ALTER TABLE stores ADD COLUMN submitted_at timestamptz;
ALTER TABLE stores ADD COLUMN logo_key text;
ALTER TABLE stores ADD COLUMN banner_key text;
ALTER TABLE store_pickup_points ADD COLUMN photo_key text;
ALTER TABLE store_pickup_points ADD COLUMN rejection_note text;

-- +goose Down
ALTER TABLE store_pickup_points DROP COLUMN rejection_note;
ALTER TABLE store_pickup_points DROP COLUMN photo_key;
ALTER TABLE stores DROP COLUMN banner_key;
ALTER TABLE stores DROP COLUMN logo_key;
ALTER TABLE stores DROP COLUMN submitted_at;
ALTER TABLE stores DROP COLUMN review_note;
DROP TABLE uploads;
