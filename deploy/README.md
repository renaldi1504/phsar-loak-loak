# Deploy Phsar ke server (Linode)

Satu server menjalankan semuanya dengan Docker: website pembeli, Seller Center (`/seller`), admin (`/admin`),
kurir (`/courier`), API untuk aplikasi, worker, cari-dengan-foto, database, HTTPS otomatis, dan backup harian.

```
internet ──► caddy (HTTPS) ──► app ──► db (Postgres 16)
                               worker      ▲
                               embedder    └── backup (setiap malam → deploy/backups/)
```

## Tahap pra-peluncuran (kumpul seller dulu)

Dengan `PAYMENT_PROVIDER=off` (bawaan di `env.production.example`):

- website tampil dengan pita **"Phsar segera dibuka"** dan tombol **Buka toko Anda** ke Seller Center;
- seller bisa login (Telegram), daftar toko, KYC, titik penjemputan, upload produk & video;
- admin memverifikasi toko dan memoderasi produk seperti biasa;
- pembeli bisa melihat-lihat dan mengisi keranjang, tapi **checkout ditutup** (API juga menolak).

Login SMS dimatikan (`SMS_PROVIDER=off`) sampai ada provider SMS Kamboja. Semua orang login lewat Telegram.

Saat siap jualan: isi `PAYMENT_PROVIDER=payway` + kunci PayWay produksi, lalu `./phsar.sh up`.

## Yang dibutuhkan

| | |
|---|---|
| Server | Linode **Shared 4 GB** (2 CPU) Ubuntu 24.04 LTS — 2 GB terlalu sempit untuk Go + Postgres + embedder + build. Region: **Singapore**. Aktifkan **Linode Backups** (cadangan seluruh disk, di luar server). |
| Domain | A record `@`, `www`, `seller`, `admin`, `api` → IP server (atau satu wildcard `*`). |
| Bot Telegram produksi | Bot baru dari @BotFather (bukan @loakloak_dev_bot). `/setdomain` → domain Phsar. |
| Akses GitHub | Repo privat → buat **deploy key** (read-only) di GitHub → Settings → Deploy keys. |

## Langkah pertama kali

**1. Siapkan server** (login SSH sebagai root):

```bash
curl -fsSLO https://raw.githubusercontent.com/renaldi1504/phsar-loak-loak/main/deploy/setup-server.sh  # atau salin manual
bash setup-server.sh
```

**2. Ambil kode** (dengan deploy key):

```bash
ssh-keygen -t ed25519 -f ~/.ssh/phsar_deploy -N ""      # tempel ~/.ssh/phsar_deploy.pub ke GitHub Deploy keys
printf 'Host github.com\n  IdentityFile ~/.ssh/phsar_deploy\n' >> ~/.ssh/config
git clone git@github.com:renaldi1504/phsar-loak-loak.git /opt/phsar
cd /opt/phsar/deploy
```

**3. Isi pengaturan** — nilai acak dibuat otomatis, sisanya diisi sendiri:

```bash
./phsar.sh secrets        # isi POSTGRES_PASSWORD, JWT_SECRET, TELEGRAM_WEBHOOK_KEY, DATA_ENCRYPTION_KEY
nano .env                 # isi PHSAR_DOMAIN, ACME_EMAIL, PUBLIC_BASE_URL, bot Telegram, ANTHROPIC_API_KEY, GOOGLE_MAPS_API_KEY
./phsar.sh check          # semua harus ✅
```

> Salin isi `.env` ke password manager. **`DATA_ENCRYPTION_KEY` tidak boleh hilang atau diganti** — tanpa kunci itu
> nomor rekening seller tidak bisa dibaca lagi.

**4. Jalankan:**

```bash
./phsar.sh up             # build ±5–10 menit pertama kali (unduh model foto ±150 MB)
./phsar.sh logs app       # tunggu "server berjalan" + "telegram webhook aktif"; Ctrl+C
```

Sertifikat HTTPS dibuat otomatis dalam ±1 menit untuk setiap alamat:

| Alamat | Untuk |
|---|---|
| `https://<domain>/` | Pembeli (website, keranjang), kurir `/courier` |
| `https://seller.<domain>/` | Seller Center (`<domain>/seller` dialihkan ke sini) |
| `https://admin.<domain>/` | Admin (`<domain>/admin` sengaja 404) |
| `https://api.<domain>/v1/…` | API aplikasi Flutter |

Login di tiap subdomain terpisah. Pembagiannya diatur di `Caddyfile`; backend tetap satu.

**5. Admin pertama:**

1. Buka `https://admin.<domain>/`, login dengan Telegram (akan ditolak karena belum admin — itu normal, akunnya jadi terdaftar).
2. Di server: `./phsar.sh admin-grant 0xx xxx xxx super_admin`
3. Login lagi → pindai QR dengan aplikasi authenticator → simpan 10 kode cadangan.
4. Ulangi untuk super admin kedua (supaya bisa saling reset MFA).

## Server testing (staging) — coba dulu sebelum produksi

Sama persis dengan langkah di atas, dengan perbedaan:

- Server terpisah (Linode Shared 4 GB, bisa dihapus setelah selesai). Jangan dicampur dengan server yang sudah menjalankan web lain di port 80/443.
- Domain: `test.<domain>` + `seller.test`, `admin.test`, `api.test` (A record → IP server testing), atau tanpa domain: `<IP-pakai-strip>.sslip.io` (subdomain `admin.<IP>.sslip.io` dll. langsung jalan).
- Pengaturan: `cp env.staging.example .env` **sebelum** `./phsar.sh secrets`. Isinya `APP_ENV=staging`,
  `PAYMENT_PROVIDER=fake` (tombol simulasi bayar), `SMS_PROVIDER=log` (kode login tidak dikirim), `ROBOTS_TAG=noindex`.
- Login tanpa SMS/Telegram: masukkan nomor HP di halaman login, lalu di server `./phsar.sh otp 012345678` untuk melihat kodenya.
- Kunci AI, Maps, PayWay sandbox: isi di Admin → Integrasi (bukan di chat).

## Perintah harian

| Perintah | Untuk apa |
|---|---|
| `./phsar.sh update` | Ambil versi terbaru dari GitHub, build, jalankan ulang. Migrasi database otomatis. |
| `./phsar.sh status` | Layanan mana yang jalan / sehat. |
| `./phsar.sh logs app` | Log server (juga `worker`, `caddy`, `db`, `embedder`, `backup`). |
| `./phsar.sh backup` | Backup sekarang (sebelum perubahan besar). |
| `./phsar.sh backups` | Daftar backup. |
| `./phsar.sh restore backups/db/phsar-….dump` | Pulihkan database (minta konfirmasi). |
| `./phsar.sh admin-grant <hp> <peran>` | Tambah admin dari server (biasanya cukup lewat menu Admins). |
| `./phsar.sh reset-mfa <hp>` | Admin kehilangan HP dan kode cadangan. |
| `./phsar.sh translate-test` | Cek kunci Anthropic. |
| `./phsar.sh otp <hp>` | (staging) Kode login terakhir dari log. |

## Backup

- **Database**: setiap hari jam 03:00 → `deploy/backups/db/` (disimpan 14 hari).
- **File upload** (foto, video, dokumen KYC): setiap Minggu → `deploy/backups/uploads/` (disimpan 4).
- Folder `backups/` ada di disk server yang sama → **aktifkan Linode Backups** supaya ada salinan di luar server.
  Nanti saat pindah ke Cloudflare R2, backup juga bisa dikirim ke sana.
- Uji restore sekali setelah online: `./phsar.sh backup` lalu `./phsar.sh restore <file terbaru>`.

## Kalau pakai Cloudflare di depan server

Set SSL/TLS ke **Full (strict)**. IP pembeli yang tercatat akan menjadi IP Cloudflare sampai Caddy dikonfigurasi
`trusted_proxies` untuk Cloudflare — beri tahu sebelum mengaktifkan proxy oranye.
