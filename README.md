# Phsar · loak-loak

Marketplace video untuk Kamboja. Launch: **Sihanoukville**, ongkir gratis (kurir Phsar menjemput di titik
penjemputan toko), komisi 8% rata, bahasa **EN (default) · ខ្មែរ · 中文**, bayar KHQR / ABA (tanpa COD).

## Isi repo

| Folder | Isi |
| --- | --- |
| `backend/` | Satu binary Go: API app & seller (`/v1/...`), admin (`/admin`), webhook, website customer (SSR, Tahap 2), worker |
| `backend/migrations/` | Skema PostgreSQL (goose), dari ERD Phsar |
| `app/` | Satu app Flutter: pembeli + mode seller (+ mode kurir) |
| `embedder/` | Layanan Python internal: embedding foto untuk cari-dengan-foto |
| `deploy/` | Online di server: docker compose (Caddy HTTPS, Postgres, app, worker, embedder, backup) + `phsar.sh` — lihat `deploy/README.md` |
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
  Akun harus dijadikan admin dulu (`phsar admin grant`), lalu login 2 langkah (lihat **Login admin + MFA**).

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

### Login admin + MFA

Semua admin login 2 langkah di `/admin`: **nomor HP** (Telegram/SMS) → **kode 6 angka** dari aplikasi authenticator
(Google Authenticator, Microsoft Authenticator, 1Password…).

1. Orangnya login sekali di `/admin` (muncul "bukan admin", tapi akunnya jadi tercatat).
2. Jadikan admin: `./backend/bin/phsar admin grant 012345678 super_admin` (atau dari menu **Admins** oleh super_admin).
3. Login lagi → scan QR → masukkan kode → **simpan 10 kode cadangan** (sekali pakai, untuk kalau HP hilang).

- API admin (`/admin/api/...`) hanya menerima sesi admin yang sudah lolos MFA; token app/Seller Center ditolak.
- Sesi admin: maksimal 8 jam, ditutup otomatis kalau 30 menit tidak dipakai. 5 kode salah → sesi dicabut; 10 salah dalam 15 menit → dikunci 15 menit.
- Secret authenticator disimpan terenkripsi dengan `DATA_ENCRYPTION_KEY`; kode yang sama tidak bisa dipakai dua kali.
- HP admin hilang: super_admin buka **Admins** → **Reset MFA**, atau `./backend/bin/phsar admin reset-mfa 012345678`.
- Menu **Admins** (super_admin): tambah admin, ganti peran, nonaktifkan, reset MFA, buat kode cadangan baru.

Endpoint admin:

`GET /admin/api/stores` (antrean) · `GET /admin/api/stores/{id}` · `POST …/kyc` · `POST /admin/api/pickup-points/{id}/decision` ·
`POST …/stores/{id}/decision` — body `{"decision":"approve"}` atau `{"decision":"reject","note":"…"}`.
Toko aktif kalau KYC terverifikasi + minimal satu titik penjemputan terverifikasi.

File upload disimpan di `backend/data/uploads` (tidak ikut git). File `public/…` disajikan di `/media/public/…`;
dokumen KYC (`private/…`) hanya bisa dibuka admin lewat `/admin/api/files/…` dan setiap akses tercatat di audit log.

## Website pembeli

Buka http://localhost:8080/ → otomatis ke /km/, /en/ atau /zh/ (dari cookie, lalu bahasa browser).

| URL | Isi |
|---|---|
| `/{lang}/` | beranda: kategori, video, produk baru |
| `/{lang}/c/{kategori}` | kategori (urut terbaru / harga), 24 per halaman |
| `/{lang}/search?q=` | pencarian lintas bahasa (nama Khmer ketemu dari halaman Inggris) |
| `/{lang}/p/{slug}` | detail produk: video, foto, varian, harga + riel, kondisi, toko |
| `/{lang}/s/{toko}` | halaman toko |
| `/sitemap.xml`, `/robots.txt` | untuk Google |

Hanya produk yang **aktif + disetujui moderator** dari **toko aktif** yang tampil. SEO: title/description per bahasa,
canonical, hreflang (km/en/zh-Hans), JSON-LD Product/Offer/Video/Breadcrumb/Store, sitemap dengan hreflang.
Isi `PUBLIC_BASE_URL` dengan domain asli saat online (dipakai untuk canonical & sitemap).

### Cari dengan foto

Tombol **Foto** di kotak pencarian (laptop: di header; HP: di halaman Cari & banner) → pembeli memotret / memilih foto →
muncul produk yang mirip. Butuh layanan `embedder/` (Python + model CLIP), lihat `embedder/README.md`:

```bash
cd embedder && pip install -r requirements.txt && ./download-model.sh && python3 server.py
# .env: EMBEDDER_URL=http://127.0.0.1:8090  → jalankan ulang make dev
```

Foto pembeli tidak disimpan. Maks 10 MB, 12 pencarian per menit per IP.

### Data contoh (demo)

```bash
make migrate-up
make seed-demo         # 8 toko + 39 produk (km/en/zh), langsung tampil di website
make seed-demo-reset   # hapus semua data demo lalu buat ulang
```

Toko demo memakai slug `demo-…` dan pemilik HP `010000001`–`010000008` (login Seller Center seperti biasa, kode OTP muncul di log).
Foto, logo & banner adalah gambar buatan (bukan foto asli); video dibuat dari foto kalau ffmpeg terpasang.
Perintah ini menolak jalan di `APP_ENV=production` (kecuali `--force`) — **jangan** dijalankan di server asli.
## Keranjang & checkout (KHQR)

- Keranjang disimpan di browser; server selalu menghitung ulang harga & stok (`POST /v1/shop/cart/quote`).
- Checkout: login OTP (SMS/Telegram) → alamat di peta (wajib di wilayah Sihanoukville) → **Bayar dengan KHQR**.
- Satu pembayaran, dipecah menjadi **satu pesanan per toko** (komisi 8% dicatat per pesanan).
- Stok **ditahan** saat menekan Bayar; lunas → stok terjual & seller dapat notifikasi; tidak dibayar dalam
  `PAYMENT_LIFETIME_MIN` menit → pesanan batal & stok kembali.
- Pembayaran baru dianggap lunas setelah dicek ke API PayWay (callback saja tidak cukup) dan nominalnya sama.
- Pengembangan: `PAYMENT_PROVIDER=fake` → QR palsu + tombol **Simulasi bayar**. Akun PayWay sungguhan: lihat `GO-LIVE.md`.
- Pembeli: `/{lang}/cart`, `/{lang}/checkout`, `/{lang}/orders` (status, riwayat, info refund).

## Pesanan di Seller Center & refund

- Tab **Pesanan** (pemilik & staf): Baru → Diproses → Siap dijemput → Dengan kurir → Selesai / Batal.
- Seller wajib **menerima dalam 2 jam** setelah dibayar; lewat dari itu pesanan batal otomatis (worker) dan uang dikembalikan.
- Tolak pesanan: pilih alasan (stok habis / tidak bisa dipenuhi / harga salah / lainnya). Stok kembali, kecuali alasan "stok habis".
- **Cetak label paket** (100×150 mm) dengan barcode Code 128 nomor pesanan untuk dipindai kurir.
- Seller tidak melihat nomor HP / alamat pembeli (pengantaran oleh kurir Phsar).
- Pembeli dapat notifikasi Telegram (diterima / dibatalkan / refund terkirim) kalau login dengan Telegram.
- **Refund manual**: admin (izin `refund.decide`, peran finance/super_admin) membuka menu **Refunds** di `/admin`,
  mengembalikan uang di dashboard ABA PayWay (pakai nomor transaksi yang ditampilkan), lalu mencatat nomor referensinya.

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
- **Admin**: semua peran wajib MFA saat login (API admin menolak sesi yang belum lolos MFA).

## Tahapan

0. ✅ Fondasi: repo, migrasi Tahap 0–1, health, meta, cek wilayah, kategori, CI
1. Rekrut seller: ✅ login Telegram (deep link bot) + SMS OTP + sesi per perangkat · ✅ API onboarding toko (toko, KYC, titik penjemputan, verifikasi admin) · ✅ Seller Center web (`/seller`) + halaman verifikasi admin (`/admin`) · ✅ tambah produk + varian + video/foto + terjemahan AI + moderasi produk ·
   ✅ kompres video & thumbnail · ✅ login admin + MFA · ✅ admin sesuai template (dashboard, audit log, matriks peran) · ✅ paket deploy + mode pra-peluncuran · ⏳ Facebook/Apple login
2. Pembeli melihat: website customer (SSR + SEO) + app Flutter (feed, cari, produk)
3. Pembeli membeli: keranjang, checkout Sihanoukville, KHQR/ABA PayWay, ledger
4. Pengiriman manual & setelah beli: siap dijemput, dispatch kurir, mode kurir, lacak, retur, ulasan
5. Chat, notifikasi, payout manual, kepatuhan app store, legal → beta → submit store

## Online di server

Semua langkah ada di [`deploy/README.md`](deploy/README.md): siapkan Linode → `./phsar.sh secrets` → isi `deploy/.env` →
`./phsar.sh check` → `./phsar.sh up`. Tahap awal memakai **mode pra-peluncuran** (`PAYMENT_PROVIDER=off`): seller sudah bisa
daftar dan upload produk, pembeli melihat pita "segera dibuka" dan belum bisa checkout.

## Ke GitHub

```bash
git remote add origin git@github.com:<username>/phsar.git
git push -u origin main
```

## Kurir Phsar (staf)

- Admin menambah kurir di `/admin` → **Delivery** (nomor HP + nama + kendaraan). Butuh izin `delivery.dispatch` (peran cs / super_admin).
- Kurir membuka **`/courier`** di HP, login dengan nomor itu, lalu menyalakan **Online** (dapat notifikasi Telegram saat ada paket siap).
- Alur: **Tersedia** → *Ambil tugas* → di toko **scan barcode label** (Chrome Android; iPhone: ketik nomor pesanan) → status pesanan
  "Dengan kurir" → di tujuan *Terantar* (boleh dengan foto bukti, disimpan privat) atau *Gagal antar* dengan alasan.
- Alamat & HP pembeli baru terlihat setelah kurir mengambil tugas. Label paket yang salah ditolak ("paket salah").
- Admin: pantau tugas berjalan / gagal / terantar, tugaskan atau pindahkan ke kurir lain, buka foto bukti.
- Pembeli & seller mendapat notifikasi saat paket dijemput dan saat terantar.

## Saldo & payout seller

- Pesanan terantar ditahan **N hari** (default 3, untuk komplain) → status *completed* → `seller_amount` masuk **saldo** seller
  (buku besar `seller_ledger`, hanya ditambah, tidak pernah diubah).
- **Jadwal** diatur admin di `/admin` → **Payouts**: harian atau mingguan (hari & jam, waktu Phnom Penh), masa tahan, minimal payout.
  Worker membuat payout *pending* per toko sesuai jadwal; admin juga bisa menekan **Create payouts now**.
- Admin mentransfer **manual** ke rekening seller (nomor lengkap hanya tampil untuk peran finance/super_admin dan dicatat di audit),
  lalu mengisi nomor referensi transfer → *paid* → seller dapat notifikasi. Batal → saldo kembali.
- Seller (pemilik) melihat saldo, jadwal, riwayat, dan mengisi rekening di tab **Payouts** Seller Center.
  Nomor rekening disimpan terenkripsi (AES-GCM, `DATA_ENCRYPTION_KEY`). Ganti rekening → payout ditahan 48 jam.
