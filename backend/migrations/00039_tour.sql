-- +goose Up
-- Tur "cara pakai": status sudah dilihat per pengguna, supaya tur tidak muncul lagi di perangkat lain.
-- tour_key divalidasi di server dari daftar yang diizinkan (internal/tour); di sini hanya bentuknya.
CREATE TABLE user_tours (
  user_id  uuid        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  tour_key varchar(40) NOT NULL CHECK (tour_key ~ '^[a-z][a-z_]{0,39}$'),
  done_at  timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id, tour_key)
);

-- +goose Down
DROP TABLE user_tours;
