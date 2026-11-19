# embedder — cari dengan foto (Python)

Layanan internal kecil, hanya dipanggil oleh backend Go (jangan dibuka ke internet):

- `POST /embed` — body = byte gambar → `{"model": "clip-vit-b32-q8", "vector": [512 angka]}` (CLIP ViT-B/32, ONNX Runtime)
- `GET /health`

Layanan ini tidak menyimpan apa pun. Backend Go yang menyimpan vektor foto produk (`product_image_vectors`)
dan membandingkannya dengan foto pembeli.

## Menjalankan di laptop

```bash
cd embedder
python3 -m venv .venv && source .venv/bin/activate
pip install -r requirements.txt
./download-model.sh            # ±150 MB, sekali saja (disimpan di embedder/models/, tidak ikut git)
python3 server.py              # http://127.0.0.1:8090
```

Lalu di `.env`: `EMBEDDER_URL=http://127.0.0.1:8090` dan jalankan ulang `make dev`.
Worker otomatis memproses 3 foto pertama tiap produk (juga produk lama). Tanpa `EMBEDDER_URL`
tombol "Foto" tetap tampil, tapi pembeli mendapat pesan bahwa fitur belum tersedia.

Env: `EMBEDDER_MODEL` (path model), `EMBEDDER_ADDR` (default `127.0.0.1:8090`), `EMBEDDER_THREADS` (default 2).
RAM yang dipakai ±400 MB; satu foto ±0,1–0,3 detik di CPU.

Rencana berikutnya: `POST /moderate/media` untuk bantu moderasi video/foto otomatis.
