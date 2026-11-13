# embedder — Python (Tahap 1–2)

Layanan internal kecil, hanya bisa dipanggil oleh backend Go (tidak terbuka ke internet):

- `POST /embed/image` → vektor 512 dimensi (CLIP ViT-B/32, ONNX Runtime) untuk cari pakai foto
- `POST /moderate/media` → skor label untuk moderasi otomatis video/foto (hasilnya ke `media_moderation_results`;
  keputusan akhir tetap oleh moderator)
- `GET /health`

Belum diisi di Tahap 0. Akan dibuat saat fitur upload video dan moderasi dikerjakan.
