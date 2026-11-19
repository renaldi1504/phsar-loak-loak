-- +goose Up
-- Kurir (staf Phsar) & pengantaran. Satu pesanan "siap dijemput" = satu tugas pengantaran.
--
--   unassigned ──ambil──► assigned ──scan label di toko──► picked_up ──► delivered
--        ▲                   │                                  │
--        └────lepas / admin──┘                                  └──gagal──► failed (admin: coba lagi / kembalikan)

CREATE TABLE couriers (
  user_id      uuid PRIMARY KEY REFERENCES users(id),
  display_name varchar(80) NOT NULL,
  vehicle      varchar(16) NOT NULL DEFAULT 'motorbike' CHECK (vehicle IN ('motorbike','tuktuk','car','bicycle')),
  plate_no     varchar(20),
  active       boolean NOT NULL DEFAULT true,
  is_online    boolean NOT NULL DEFAULT false,
  last_seen_at timestamptz,
  created_by   uuid REFERENCES users(id),
  created_at   timestamptz NOT NULL DEFAULT now(),
  updated_at   timestamptz NOT NULL DEFAULT now()
);
CREATE TRIGGER couriers_updated BEFORE UPDATE ON couriers FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER couriers_audit AFTER INSERT OR UPDATE OR DELETE ON couriers FOR EACH ROW EXECUTE FUNCTION audit_row_change();

CREATE TABLE deliveries (
  id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id        uuid NOT NULL UNIQUE REFERENCES orders(id),
  courier_id      uuid REFERENCES couriers(user_id),
  status          varchar(12) NOT NULL DEFAULT 'unassigned'
                    CHECK (status IN ('unassigned','assigned','picked_up','delivered','failed','cancelled')),
  -- jemput (salinan titik penjemputan saat siap)
  pickup_label    varchar(80) NOT NULL,
  pickup_address  varchar(300) NOT NULL,
  pickup_lat      numeric(9,6) NOT NULL,
  pickup_lng      numeric(9,6) NOT NULL,
  pickup_contact  varchar(80) NOT NULL,
  pickup_phone    varchar(30) NOT NULL,
  -- antar (dari checkout)
  dropoff_lat     numeric(9,6) NOT NULL,
  dropoff_lng     numeric(9,6) NOT NULL,
  attempts        integer NOT NULL DEFAULT 0,
  fail_reason     varchar(24) CHECK (fail_reason IN ('buyer_unreachable','buyer_refused','wrong_address','damaged','other')),
  fail_note       varchar(300),
  proof_key       text,                               -- foto bukti serah terima (privat)
  assigned_at     timestamptz,
  picked_up_at    timestamptz,
  delivered_at    timestamptz,
  failed_at       timestamptz,
  created_at      timestamptz NOT NULL DEFAULT now(),
  updated_at      timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT deliveries_courier_needed CHECK (status IN ('unassigned','cancelled') OR courier_id IS NOT NULL)
);
CREATE INDEX deliveries_open_idx ON deliveries (status, created_at) WHERE status IN ('unassigned','assigned','picked_up','failed');
CREATE INDEX deliveries_courier_idx ON deliveries (courier_id, status);
CREATE TRIGGER deliveries_updated BEFORE UPDATE ON deliveries FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER deliveries_audit AFTER INSERT OR UPDATE OR DELETE ON deliveries FOR EACH ROW EXECUTE FUNCTION audit_row_change();

-- Pesanan yang sudah "siap dijemput" sebelum fitur ini ada → buat tugasnya.
INSERT INTO deliveries (order_id, pickup_label, pickup_address, pickup_lat, pickup_lng, pickup_contact, pickup_phone, dropoff_lat, dropoff_lng)
SELECT o.id, pp.label, pp.address_line, pp.lat, pp.lng, pp.contact_name, pp.contact_phone, c.lat, c.lng
FROM orders o JOIN store_pickup_points pp ON pp.id = o.pickup_point_id JOIN checkouts c ON c.id = o.checkout_id
WHERE o.status = 'ready_for_pickup';

-- Foto bukti pengantaran (disimpan privat).
ALTER TABLE uploads DROP CONSTRAINT uploads_kind_check;
ALTER TABLE uploads ADD CONSTRAINT uploads_kind_check
  CHECK (kind IN ('store_logo','store_banner','pickup_photo','kyc_document','product_media','delivery_proof'));

-- +goose Down
ALTER TABLE uploads DROP CONSTRAINT uploads_kind_check;
ALTER TABLE uploads ADD CONSTRAINT uploads_kind_check
  CHECK (kind IN ('store_logo','store_banner','pickup_photo','kyc_document','product_media'));
DROP TABLE deliveries;
DROP TABLE couriers;
