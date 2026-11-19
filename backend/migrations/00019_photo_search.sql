-- +goose Up
-- Vektor gambar (CLIP) untuk cari-dengan-foto. Disimpan sebagai real[] dan dicari di memori server
-- (cukup untuk puluhan ribu foto). Nanti bisa dipindah ke pgvector tanpa mengubah API.
CREATE TABLE product_image_vectors (
  media_id    uuid PRIMARY KEY REFERENCES product_media(id),
  product_id  uuid NOT NULL REFERENCES products(id),
  model       text NOT NULL,
  vec         real[],                 -- NULL = gagal diproses (gambar rusak), tidak dicoba ulang
  error       text,
  created_at  timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX product_image_vectors_product ON product_image_vectors (product_id);

-- +goose Down
DROP TABLE product_image_vectors;
