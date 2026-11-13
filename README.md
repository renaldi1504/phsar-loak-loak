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
1. Rekrut seller: ✅ login Telegram (deep link bot) + SMS OTP + sesi per perangkat · ⏳ Facebook/Apple, Seller Center (daftar toko, KYC, titik penjemputan, tambah produk +
   terjemahan AI), admin (verifikasi seller & titik, moderasi produk, audit)
2. Pembeli melihat: website customer (SSR + SEO) + app Flutter (feed, cari, produk)
3. Pembeli membeli: keranjang, checkout Sihanoukville, KHQR/ABA PayWay, ledger
4. Pengiriman manual & setelah beli: siap dijemput, dispatch kurir, mode kurir, lacak, retur, ulasan
5. Chat, notifikasi, payout manual, kepatuhan app store, legal → beta → submit store

## Ke GitHub

```bash
git remote add origin git@github.com:<username>/phsar.git
git push -u origin main
```
