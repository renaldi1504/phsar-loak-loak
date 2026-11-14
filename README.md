# Phsar · loak-loak

Marketplace video untuk Kamboja. Launch: **Sihanoukville**, ongkir gratis (kurir Phsar menjemput di titik
penjemputan toko), komisi 8% rata, bahasa **EN (default) · ខ្មែរ · 中文**, bayar KHQR / ABA (tanpa COD).

## Isi repo

| Folder | Isi |
| --- | --- |
| `backend/` | Satu binary Go: API app & seller (`/v1/...`), admin (`/admin`), webhook, website customer (SSR, Tahap 2), worker |
| `backend/migrations/` | Skema PostgreSQL (goose), dari ERD Phsar |
| `app/` | Satu app Flutter: pembeli + mode seller (+ mode kurir) |
| `embedder/` | Layanan Python internal: embedding foto & moderasi otomatis (menyusul) |
| `api/openapi.yaml` | Kontrak API → client Dart dibuat dari file ini |

## Menjalankan di komputer sendiri

Butuh: Docker Desktop, Go 1.23+, (nanti) Flutter 3.24+.

```bash
cp .env.example .env
docker compose up -d db          # Postgres 16 + pgvector

cd backend
go mod tidy                      # pertama kali: unduh dependency & buat go.sum
make test                        # tes pembulatan uang & cek wilayah
make migrate-up                  # buat semua tabel + data awal launch
make run                         # http://localhost:8080
```

Coba:

```bash
curl localhost:8080/readyz
curl "localhost:8080/v1/app/meta?lang=km"
curl "localhost:8080/v1/app/service-areas/check?lat=10.627&lng=103.522"   # di Sihanoukville → served: true
curl "localhost:8080/v1/app/categories?lang=zh"
```

## Mencoba login (tanpa app)

Isi `.env`: `TELEGRAM_BOT_TOKEN`, `TELEGRAM_BOT_USERNAME` (tanpa @), `TELEGRAM_MODE=polling`, `JWT_SECRET`.
Jalankan `make migrate-up run`, lalu buka **http://localhost:8080/dev/login**:

- **Continue with Telegram** → Telegram terbuka → tekan **Start** di bot → halaman otomatis menampilkan user & token.
- **SMS** → kode OTP muncul di log Terminal server (`SMS_PROVIDER=log`), karena SMS sungguhan belum dipasang.

Alur Telegram: app `POST /v1/app/auth/telegram/start` → buka `deep_link` → bot menerima `/start login_<nonce>`
→ app `POST /telegram/poll` sampai `status: ok`. Di server pakai `TELEGRAM_MODE=webhook` +
`setWebhook` ke `https://<domain>/webhooks/telegram/<TELEGRAM_WEBHOOK_KEY>`.

## Seller Center web & admin

- **http://localhost:8080/seller**: login (Telegram/SMS), buat toko, profil + logo, KYC, titik penjemputan dengan peta, kirim untuk diverifikasi. Tersedia dalam EN/ខ្មែរ/中文.
- **http://localhost:8080/seller** → tab **Products**: tambah/edit produk (video + foto, kondisi baru/bekas + foto kondisi,
  warna × ukuran dengan harga/stok per varian), simpan draf atau kirim ke moderasi, lihat & perbaiki terjemahan AI.
- **Lonceng notifikasi** (toko/KYC/titik/produk disetujui atau ditolak, staf baru) — juga dikirim lewat bot Telegram ke seller yang login dengan Telegram.
- **Tim**: owner mengundang staf lewat link (sekali pakai, 7 hari); staf hanya mengelola produk. **Edit stok & harga cepat** dari daftar produk. **Banner toko**.
- **http://localhost:8080/admin/**: antrean verifikasi toko, buka dokumen KYC, setujui/tolak KYC, titik penjemputan, dan toko.
  Menu **Products**: antrean moderasi produk (setujui · minta perbaikan · tolak), cek otomatis kata terlarang, perbaiki terjemahan.
  Akun harus dijadikan admin dulu (`phsar admin grant`).

Sesi web memakai cookie HttpOnly (refresh token) + access token di memori halaman.

## Onboarding seller (API, Tahap 1)

Semua butuh header `Authorization: Bearer <access_token>` (dari login di atas).

| Langkah | Endpoint |
| --- | --- |
| Upload file (logo, foto titik, dokumen KYC) | `POST /v1/seller/uploads` multipart: `kind` (`store_logo` · `store_banner` · `pickup_photo` · `kyc_document`) lalu `file` |
| Buat toko | `POST /v1/seller/stores` `{name, about, about_lang}` |
| Toko saya / detail + checklist | `GET /v1/seller/stores` · `GET /v1/seller/stores/{id}` |
| Ubah toko | `PATCH /v1/seller/stores/{id}` `{logo_key, banner_key, social_links, about, about_lang}` |
| KYC | `PUT /v1/seller/stores/{id}/kyc` `{seller_type, legal_name, id_document_key, tin?}` |
| Titik penjemputan | `POST/GET /v1/seller/stores/{id}/pickup-points` · `PUT/DELETE …/{pointID}` (pin wajib di Sihanoukville) |
| Kirim untuk diverifikasi | `POST /v1/seller/stores/{id}/submit` |

Admin (sementara memakai token login biasa; login admin + MFA menyusul):

```bash
make -C backend build && ./backend/bin/phsar admin grant 012345678 moderator   # jadikan admin (user harus sudah pernah login)
```

`GET /admin/api/stores` (antrean) · `GET /admin/api/stores/{id}` · `POST …/kyc` · `POST /admin/api/pickup-points/{id}/decision` ·
`POST …/stores/{id}/decision` — body `{"decision":"approve"}` atau `{"decision":"reject","note":"…"}`.
Toko aktif kalau KYC terverifikasi + minimal satu titik penjemputan terverifikasi.

File upload disimpan di `backend/data/uploads` (tidak ikut git). File `public/…` disajikan di `/media/public/…`;
dokumen KYC (`private/…`) hanya bisa dibuka admin lewat `/admin/api/files/…` dan setiap akses tercatat di audit log.

## QA (cek semua kode)

```bash
make db && make qa
```

Format, vet, build, unit test, **tes integrasi semua alur API** di database `phsar_test` (migrasi up/reset/up,
setiap error 500 = gagal), sintaks JS, dan cek API key/token yang tidak sengaja ikut git. Wajib lulus sebelum push;
juga jalan otomatis di GitHub Actions. Checklist manual (peta, Telegram, video, HP, bahasa Khmer): [QA.md](QA.md).

## Google Maps (titik penjemputan)

Isi `GOOGLE_MAPS_API_KEY` di `.env` → form titik penjemputan memakai Google Maps: seller mengetik nama tempat,
pin + alamat terisi otomatis, pin bisa digeser, tombol "Gunakan lokasi saya". Tanpa key: peta OpenStreetMap.

Di Google Cloud Console (project yang sama dengan key):
1. Aktifkan **Maps JavaScript API**, **Places API (New)**, dan **Geocoding API**.
2. Credentials → key → *Application restrictions*: **Websites**, isi `http://localhost:8080/*` (nanti tambah domain Phsar).
   *API restrictions*: hanya tiga API di atas.
3. `GOOGLE_MAPS_MAP_ID=DEMO_MAP_ID` cukup untuk pengembangan; untuk production buat Map ID sendiri (Map Management).

Key ini memang terlihat di browser (disajikan lewat `/ui/config.js`), jadi pembatasan domain wajib.

## Kompres video & thumbnail (ffmpeg)

```bash
brew install ffmpeg      # sekali saja
make run                 # log: "worker media berjalan (kompres video + thumbnail)"
make media-backfill      # (sekali) proses video/foto yang di-upload sebelum ffmpeg terpasang
```

- Video dari HP → MP4 H.264, sisi terpanjang ≤ 1280 px, ±10× lebih kecil, bisa diputar sebelum selesai diunduh (faststart),
  plus gambar sampul (poster). Video > 3 menit atau file rusak ditolak saat upload.
- Foto → JPEG ≤ 1600 px, rotasi HP diperbaiki, data EXIF termasuk lokasi GPS dibuang.
- Selama diproses, file asli yang tampil (badge "Compressing video…"). File asli tetap disimpan.
- Tanpa ffmpeg semuanya tetap jalan dengan file asli.

## Produk & terjemahan AI

| Langkah | Endpoint (`/v1/seller/stores/{id}/…`) |
| --- | --- |
| Upload video/foto | `POST /v1/seller/uploads` `kind=product_media` (foto ≤ 8 MB; video mp4/mov/webm ≤ 100 MB) |
| Daftar produk | `GET products?filter=all\|draft\|pending\|live\|changes\|archived` |
| Buat / ubah | `POST products` · `PUT products/{pid}` — `{category_id, source_lang, name, description, condition, condition_note, options[], variants[], media[], submit}` |
| Ubah harga/stok saja (tanpa moderasi ulang) | `PATCH products/{pid}/variants` |
| Arsip / tampilkan lagi · hapus | `POST products/{pid}/status {status}` · `DELETE products/{pid}` |
| Perbaiki terjemahan sendiri · kembali ke AI | `PUT products/{pid}/translations/{lang}` · `DELETE …/{lang}` |

Syarat kirim ke moderasi: toko aktif, 1 video, minimal 1 foto, barang bekas (used_good/used_fair) wajib foto kondisi.
Ubah nama/deskripsi/kategori/kondisi/foto/video → dicek ulang; ubah harga/stok → langsung berlaku.

Admin: `GET /admin/api/products?status=pending_review` · `GET /admin/api/products/{id}` ·
`POST …/{id}/decision {decision: approve|changes_requested|reject, reason_code, note}` · `PUT …/{id}/translations/{lang}`.

**Terjemahan**: seller menulis dalam satu bahasa. Worker (ikut jalan di `make run`) menerjemahkan nama, deskripsi,
nama opsi dan nilai varian ke dua bahasa lain lewat Claude (`ANTHROPIC_API_KEY`), cadangan Google
(`GOOGLE_TRANSLATE_API_KEY`). Terjemahan yang diedit seller tidak ditimpa. Tanpa API key, job menunggu di antrean
`translation_jobs`. Untuk mencoba tanpa key: `TRANSLATOR=fake`. Gagal → dicoba lagi (maks. 5×, jeda makin lama).

## Aturan yang sudah dipasang di database

- **Uang**: `NUMERIC(12,2)` USD, dihitung dengan `shopspring/decimal`, pembulatan biasa (12.5960 → 12.60).
  Komisi: seller = total − komisi, jadi jumlahnya selalu pas. Riel hanya tampilan (dibulatkan ke ៛100).
- **Audit log**: setiap perubahan tabel penting tercatat otomatis (siapa, kapan, kolom apa, dari nilai apa ke apa).
  Pelaku dibaca dari `SET LOCAL app.actor_id` → di Go selalu ubah data lewat `db.WithActor(...)`.
  `audit_logs` tidak bisa di-UPDATE/DELETE.
- **Soft delete**: kolom `deleted_at`; email, nomor HP, slug, dan SKU unik hanya untuk data yang aktif.
- **Wilayah layanan**: `service_areas` (GeoJSON). Seed masih **kotak kasar** Sihanoukville — ganti dengan poligon
  sungguhan sebelum launch.
- **Titik penjemputan** wajib per toko; produk baru tampil kalau titik sudah diverifikasi admin.
- **Barang bekas** wajib catatan kondisi. **Produk** baru tampil kalau `moderation_status = approved`.
- **Komisi** `commission_rules`: global 8% saat ini; bisa ditambah per kategori/toko dari admin tanpa ubah kode.
- **Admin** finance & super_admin wajib MFA (dipaksa CHECK constraint).

## Tahapan

0. ✅ Fondasi: repo, migrasi Tahap 0–1, health, meta, cek wilayah, kategori, CI
1. Rekrut seller: ✅ login Telegram (deep link bot) + SMS OTP + sesi per perangkat · ✅ API onboarding toko (toko, KYC, titik penjemputan, verifikasi admin) · ✅ Seller Center web (`/seller`) + halaman verifikasi admin (`/admin`) · ✅ tambah produk + varian + video/foto + terjemahan AI + moderasi produk ·
   ⏳ transcode/kompres video & thumbnail, Facebook/Apple login, login admin + MFA, halaman audit
2. Pembeli melihat: website customer (SSR + SEO) + app Flutter (feed, cari, produk)
3. Pembeli membeli: keranjang, checkout Sihanoukville, KHQR/ABA PayWay, ledger
4. Pengiriman manual & setelah beli: siap dijemput, dispatch kurir, mode kurir, lacak, retur, ulasan
5. Chat, notifikasi, payout manual, kepatuhan app store, legal → beta → submit store

## Ke GitHub

```bash
git remote add origin git@github.com:<username>/phsar.git
git push -u origin main
```
