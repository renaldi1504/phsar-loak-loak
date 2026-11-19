# Catatan sebelum online (GO-LIVE)

Satu tempat untuk semua hal yang harus dikerjakan saat Phsar dipindah dari laptop ke server.
Centang `[x]` kalau sudah. Setiap fitur baru yang butuh langkah saat online ditambahkan di sini.

## 0. Tahap 1 — online untuk kumpul seller (pra-peluncuran)
Paket deploy siap di folder `deploy/` — langkah lengkap di **`deploy/README.md`**. Ringkasnya:
- [ ] Linode **Shared 4 GB**, Ubuntu 24.04, region Singapore, **Linode Backups aktif**.
- [ ] Domain dibeli; A record `@` dan `www` → IP server.
- [ ] Bot Telegram **produksi** baru di @BotFather + `/setdomain` ke domain Phsar.
- [ ] Deploy key GitHub (read-only) untuk server.
- [ ] `bash deploy/setup-server.sh` (firewall, swap, jam NTP, Docker).
- [ ] `./phsar.sh secrets` → isi `.env` → `./phsar.sh check` → `./phsar.sh up`.
- [ ] Salin `deploy/.env` ke password manager (terutama `DATA_ENCRYPTION_KEY`, `POSTGRES_PASSWORD`).
- [ ] 2 super admin dibuat (`./phsar.sh admin-grant @usernameTelegram super_admin`) + MFA dipasang.
- [ ] `./phsar.sh backup` lalu uji `./phsar.sh restore <file>` sekali.
- [ ] Mode pra-peluncuran: `PAYMENT_PROVIDER=off` (checkout ditutup, pita "segera dibuka"), `SMS_PROVIDER=off` (login Telegram saja).
- [ ] Bagikan `https://<domain>/seller` ke seller pertama.

## 0a. Server testing — khloak.online (staging)

- [ ] Kode terbaru sudah di-commit & push ke GitHub (server mengambil kode dari sana)
- [ ] Linode Shared 4 GB, Ubuntu 24.04, Singapore, Linode Backups aktif → `bash setup-server.sh`
- [ ] DNS khloak.online: A record `@`, `www`, `seller`, `admin`, `api` → IP server
- [ ] `cp env.staging.example .env` → `./phsar.sh secrets` → isi `PHSAR_DOMAIN=khloak.online`, `ACME_EMAIL`, `PUBLIC_BASE_URL=https://khloak.online` → `./phsar.sh check` → `./phsar.sh up`
- [ ] Admin pertama (`./phsar.sh admin-grant … super_admin`) + MFA
- [ ] Cek: login (kode dari `./phsar.sh otp`), daftar toko, upload produk & video, simulasi bayar, notifikasi berbunyi + notifikasi browser (HTTPS), backup & restore
- [ ] Bot Telegram TESTING terpisah (opsional) → /setdomain khloak.online

## 0b. Legalitas & pengesahan (WAJIB beres sebelum transaksi pertama)
Tahap 0 (kumpul seller, checkout ditutup) boleh jalan duluan. `PAYMENT_PROVIDER` **jangan** dinyalakan sebelum bagian ini selesai.
Denda menjalankan e-commerce tanpa license dari MoC: KHR 10.000.000 (~US$2.500).

**A. Badan usaha & pajak**
- [ ] Putuskan struktur kepemilikan perusahaan (Co., Ltd., 100% asing atau ada pemegang saham lokal). Pajak perusahaan sama saja; bedanya hanya WHT 14% atas dividen ke pemegang saham **non-residen** pajak. Hindari nominee (nama orang lokal sebagai kedok).
- [ ] Daftar perusahaan lewat CamDX (Online Business Registration): nama Khmer + Inggris, kode kegiatan usaha "online market services", modal minimal ~KHR 4 juta.
- [ ] Alamat kantor terdaftar + surat dari otoritas setempat.
- [ ] Registrasi GDT: TIN, Patent Tax Certificate (diperpanjang tiap 31 Maret), **sertifikat VAT**.
- [ ] Rekening bank perusahaan: operasional, penampung dana seller, cadangan pajak (sesuai desain ledger).
- [ ] Konsultan pajak / akuntan untuk lapor bulanan (VAT, withholding, prepayment TOI, pajak gaji) dan TOI tahunan.
- [ ] Putuskan: komisi 8% **belum** atau **sudah** termasuk VAT 10% → sesuaikan di Admin → Invoices & tax dan Perjanjian seller.
- [ ] Work permit + visa kerja yang sesuai untuk pemilik/staf asing.
- [ ] Model **agen**: omzet Phsar = komisi (+ ongkir kalau dipungut, + iklan nanti), BUKAN total belanja pembeli. Tulis tegas di Perjanjian seller; struk pembeli atas nama seller, Phsar menerbitkan faktur komisi ke seller.
- [ ] Saat daftar di CamDX, masukkan juga kode kegiatan **iklan/periklanan online** supaya nanti tidak perlu amandemen. Fitur Ads di admin tetap **dimatikan** sampai disahkan.

**B. Izin e-commerce**
- [ ] Domain + sertifikat online di MPTC (Kementerian Pos & Telekomunikasi).
- [ ] E-commerce **license** di MoC (portal ecommercelicensing.moc.gov.kh / CamDX, ~US$245, berlaku 3 tahun) — setelah A selesai.
- [ ] Simpan kalender perpanjangan license (3 tahun) & patent tax (tahunan).

**C. Alur dana (risiko izin NBC)**
- [ ] Konsultasi pengacara lokal (DFDL / Tilleke & Gibbins / BNG Legal): apakah model "uang pembeli masuk rekening Phsar → payout ke seller" butuh izin NBC.
- [ ] Tanya ABA PayWay: ada fitur marketplace / split payout / sub-merchant? Kalau ada, pertimbangkan pakai itu.
- [ ] Kalau tetap pakai rekening penampung: kontrak keagenan dengan seller, payout cepat (tanpa saldo mengendap lama), tidak ada fitur mirip dompet/saldo yang bisa dibelanjakan.
- [ ] Merchant PayWay **atas nama perusahaan Phsar**. Kredensial Tuktuk Cafe hanya untuk testing — hapus sebelum live.

**D. Dokumen & data seller**
- [ ] Syarat pembeli, Kebijakan Privasi, Perjanjian seller, Kebijakan retur, Hapus data — direview pengacara (lanjut di bagian 12b). URL publik juga dibutuhkan Apple, Google, Facebook.
- [ ] Perjanjian seller memuat: komisi, jadwal payout, barang terlarang, retur/refund, tanggung jawab barang bekas, informasi minimum perlindungan konsumen (contoh kontrak MoC).
- [ ] Seller perorangan: form KYC minta nomor e-commerce permit (perorangan ~US$50) atau pernyataan termasuk pengecualian; cek batas omzet seller tanpa TIN.
- [ ] Daftar barang terlarang dipasang di moderasi (obat, alkohol, rokok, barang palsu, dll).
- [ ] Kebijakan privasi menyebut: server di luar Kamboja (Linode Singapore) & teks produk dikirim ke API AI untuk terjemahan. Pantau UU Perlindungan Data Pribadi (masih draf).

**E. Akun app store (atas nama perusahaan)**
- [ ] Ajukan D-U-N-S number begitu perusahaan terdaftar (gratis, bisa sampai ~30 hari). Nama & alamat harus sama persis dengan dokumen resmi.
- [ ] Apple Developer Program **Organization** (US$99/tahun, proses 2–4 minggu).
- [ ] Google Play Console **Organization** (US$25; bebas syarat 12 tester × 14 hari).
- [ ] Sebelum submit: Sign in with Apple (karena ada login Telegram/Facebook), hapus akun di dalam app, akun demo pembeli & seller untuk reviewer, target API 36 Android.

**F. Setelah launch**
- [ ] Daftarkan merek "Phsar" di Departemen Kekayaan Intelektual MoC.
- [ ] Sambungkan Meta business portfolio ke perusahaan (Facebook Login).

## 1. Server & domain
- [x] Paket Docker: Caddy (HTTPS otomatis Let's Encrypt), Postgres 16 + pgvector, app, worker terpisah, embedder, backup harian (`deploy/`).
- [ ] Domain Phsar dibeli, DNS diarahkan ke server (Linode).
- [ ] Backup harian otomatis jalan (`./phsar.sh backups` ada file baru tiap pagi) + uji restore sekali.
- [x] ffmpeg terpasang di server (sudah ada di `backend/Dockerfile`).
- [x] Migrasi otomatis setiap `./phsar.sh up` / `update` (layanan `migrate`).
- [ ] Buat admin pertama: login sekali di `https://<domain>/admin`, lalu `./phsar.sh admin-grant <nomor HP atau @username Telegram> super_admin`,
      login lagi dan pasang aplikasi authenticator. Simpan 10 kode cadangan di password manager.
- [x] Production: worker terpisah (`WORKER_IN_SERVE=false` + layanan `worker` di docker compose).

## 2. Isi `.env` di server (jangan pernah di-commit / dikirim ke chat)
- [ ] `APP_ENV=production`
- [ ] `PUBLIC_BASE_URL=https://<domain>` — dipakai canonical, sitemap, link di Telegram, URL media.
- [ ] `JWT_SECRET` = string acak baru (≥ 32 karakter), **beda** dengan yang di laptop.
- [ ] `DATABASE_URL` ke database server (pakai `sslmode=require` kalau database terpisah).
- [ ] `DEFAULT_LANG=km` (bahasa default website pembeli).
- [ ] `LOG_LEVEL=info`
- [ ] `STORAGE_DIR` ke disk yang ikut di-backup (sampai pindah ke Cloudflare R2).

## 3. Telegram
- [ ] Bot dari @BotFather: `TELEGRAM_BOT_TOKEN`, `TELEGRAM_BOT_USERNAME`.
- [ ] `TELEGRAM_MODE=webhook`, isi `TELEGRAM_WEBHOOK_KEY` (acak; `./phsar.sh secrets`). Webhook didaftarkan otomatis saat server start
      (cek log: "telegram webhook aktif").
- [ ] Set domain login di @BotFather (`/setdomain`) kalau memakai tombol Telegram Login di web.

## 4. SMS OTP
- [x] Sementara `SMS_PROVIDER=off`: login SMS disembunyikan, semua orang login lewat Telegram (server menolak `log` di production).
- [ ] Akun **PlasGate** (plasgate.com) atas nama Phsar: daftarkan Sender ID (mis. "Phsar"), isi saldo SMS.
- [ ] Admin → **Login & SMS**: pilih PlasGate, isi Private key + Secret (SMPP & API → SMS Gateway) + Sender ID, simpan,
      lalu **Send test SMS** ke HP sendiri. Setiap kode login = 1 SMS dari saldo PlasGate.

## 4b. Login Facebook
- [ ] developers.facebook.com → Create App (tipe *Consumer*) → tambah produk **Facebook Login**.
- [ ] Facebook Login → Settings → **Valid OAuth Redirect URIs**: salin keempat alamat dari admin → Login & SMS.
- [ ] App Settings → Basic: isi Privacy Policy URL + Terms URL (halaman legal Phsar), App Domain = domain Phsar, lalu **Live mode**.
- [ ] Admin → Login & SMS: isi App ID + App Secret, centang "Show Continue with Facebook", simpan, klik **Test App ID & Secret**.
- [ ] Aplikasi Flutter nanti memakai App ID yang sama (endpoint `POST /v1/app/auth/facebook`).

## 4d. Login Google (& Apple untuk aplikasi iPhone)

- [ ] Google Cloud (project yang sama dengan Maps) → APIs & Services → **OAuth consent screen**: nama aplikasi Phsar, logo, email kontak, domain khloak.online / domain produksi → Publish
- [ ] Credentials → **Create OAuth client ID** → tipe *Web application* → tempel semua **Authorized redirect URIs** dan **JavaScript origins** dari Admin → Integrasi → Login & SMS → kartu Google (tombol Copy)
- [ ] Client ID + Client Secret diisi di kartu Google (bukan di .env / chat) → centang tampilkan → Simpan → **Test Client ID & Secret** harus ✅
- [ ] Coba login "Lanjut dengan Google" di Seller Center, website (keranjang), admin
- [ ] Nanti untuk aplikasi Flutter: buat OAuth client Android & iOS → isi di "App client IDs"
- [ ] Sign in with Apple: hanya aplikasi iPhone, aktifkan setelah akun Apple Developer ada (Bundle ID di kartu Apple)

## 4c. Checkout sebagai tamu
- [x] Tombol "Continue as guest" di checkout: tanpa akun; pesanan tersimpan di browser itu. Tamu tidak bisa membuka toko.
- [ ] Batas 30 akun tamu per IP per jam (anti-spam) — tinjau ulang setelah ada data.

## 5. Terjemahan AI
- [ ] AI terjemahan & draf balasan: isi kunci di **Admin → Integrasi → AI & terjemahan** (Anthropic utama + saldo; Gemini cadangan gratis dari aistudio.google.com), lalu tekan "Tes". Kunci Google Maps juga di Admin → Integrasi → Peta. `.env` hanya cadangan.
- [ ] Cek: `make translate-test` (atau `phsar translate-test`).
- [ ] Ulangi terjemahan yang dulu gagal (tombol "Coba lagi" di Seller Center / admin).
- [ ] Opsional cadangan: `GOOGLE_TRANSLATE_API_KEY` + `TRANSLATION_FALLBACK=google`.

## 6. Google Maps / Places
- [ ] Di Google Cloud: tambah `https://<domain>/*` ke *Website restrictions* API key.
- [ ] Buat Map ID sendiri (ganti `DEMO_MAP_ID`) → `GOOGLE_MAPS_MAP_ID`.
- [ ] Pasang batas kuota/budget alert di Google Cloud Billing.

## 7. Media
- [ ] Sekali setelah deploy: `phsar media-backfill` (kompres video & foto lama).
- [ ] Nanti: pindah ke Cloudflare R2 (file + CDN), lalu salin folder `uploads` lama.

## 8. Website pembeli & SEO
- [ ] Buka `https://<domain>/` → harus pindah ke `/km/` (atau bahasa browser).
- [ ] Cek `https://<domain>/robots.txt` dan `https://<domain>/sitemap.xml` berisi domain asli (bukan localhost).
- [ ] Daftarkan domain di Google Search Console, kirim `sitemap.xml`.
- [ ] Uji 1 halaman produk di Rich Results Test (data Product muncul).
- [ ] Pastikan tidak ada data demo di database server (`SELECT count(*) FROM stores WHERE slug LIKE 'demo-%'` = 0).
- [ ] Ganti gambar share `backend/internal/storefront/assets/og.png` dengan desain resmi (1200×630) kalau sudah ada.

## 8b. Cari dengan foto
- [ ] Jalankan `embedder/` di server sebagai proses terpisah (systemd/docker), hanya listen di `127.0.0.1:8090`:
      `pip install -r embedder/requirements.txt && embedder/download-model.sh && python3 embedder/server.py`.
- [ ] `.env` server: `EMBEDDER_URL=http://127.0.0.1:8090`. Cek log: "cari-dengan-foto aktif" dan "indeks foto".
- [ ] RAM ±400 MB untuk embedder — pastikan paket Linode cukup (min. 2 GB total untuk Go + Postgres + embedder).
- [ ] Uji dari HP: tombol Foto → potret barang yang dijual → produknya muncul.
- [ ] Kalau produk sudah puluhan ribu: pindahkan pencarian vektor ke pgvector (sekarang dicari di memori).

## 9. Keamanan & admin
- [x] Login admin + MFA (aplikasi authenticator) sudah dibuat — semua admin wajib MFA.
- [ ] Minimal **2 super_admin** (kalau satu kehilangan HP + kode cadangan, yang lain bisa Reset MFA). Cadangan terakhir: `phsar admin reset-mfa <hp>` di server.
- [ ] Jam server sinkron (NTP aktif) — kode authenticator bergantung pada jam yang tepat.
- [ ] Nonaktifkan admin yang keluar dari tim lewat menu **Admins** (sesinya langsung dicabut).
- [ ] `/dev/login` hanya aktif di `APP_ENV=development` — pastikan production tidak development.
- [ ] Rahasia lama yang pernah terlihat (token GitHub, bot, API key) di-rotate.
- [x] Firewall server: hanya port 22/80/443 (`setup-server.sh`), fail2ban untuk SSH, database tidak terbuka ke internet.
- [x] IP pembeli asli dari Caddy (header palsu dari luar dibuang) — untuk rate limit & audit log.
- [ ] Opsional: Cloudflare di depan server — kalau diaktifkan, atur `trusted_proxies` Caddy dulu (lihat `deploy/README.md`).

## 10. Pembayaran KHQR (ABA PayWay) & pesanan
**Syarat sebelum daftar merchant PayWay atas nama Phsar** (urutannya penting):
- [ ] Perusahaan terdaftar di Ministry of Commerce (biasanya Private Limited) lewat portal online CamDX — ±1–3 minggu.
- [ ] TIN dari General Department of Taxation (GDT) + patent tax tahunan.
- [ ] Lisensi e-commerce dari MoC (marketplace/platform = *license* untuk badan hukum, bukan *permit* perorangan).
- [ ] Rekening bisnis ABA atas nama perusahaan: (1) rekening operasional Phsar, (2) rekening terpisah dana seller.
- [ ] Website online dengan halaman Syarat & Ketentuan, Privasi, Kebijakan Retur/Refund, dan Kontak (biasanya dicek bank).
- [ ] Email ke paywaysales@ababank.com: jelaskan model marketplace (uang pembeli masuk Phsar → payout ke seller),
      minta daftar dokumen, biaya (MDR) KHQR/kartu, akses sandbox, Refund API, dan apakah payout ke seller boleh/bisa lewat API.
- [ ] Tanyakan ke ABA/konsultan: apakah menampung dana seller butuh izin NBC tambahan (payment service) atau cukup akun merchant.
- [ ] Teknis: digitalsupport@ababank.com untuk integrasi; minta production base URL + whitelist IP server + URL callback.

- [ ] Daftar akun merchant ABA PayWay atas nama Phsar (butuh dokumen usaha; bisa beberapa minggu — mulai secepatnya).
      Sementara memakai kredensial Tuktuk Cafe untuk uji coba: sebelum buka untuk umum, ganti ke akun Phsar sendiri
      (uang pembeli & payout seller jangan tercampur dengan rekening cafe; jenis usaha di ABA juga harus marketplace).
- [ ] Minta akun **sandbox** dulu (https://sandbox.payway.com.kh/register-sandbox/), isi `.env` laptop:
      `PAYMENT_PROVIDER=payway`, `PAYWAY_MERCHANT_ID`, `PAYWAY_API_KEY`, `PAYWAY_BASE_URL=https://checkout-sandbox.payway.com.kh`,
      lalu coba 1 pesanan sungguhan dari HP (scan QR sandbox).
- [ ] Uji kunci tanpa membuka website: `cd backend && go run ./cmd/phsar payway-test` → buat QR $0.01 (gambar `payway-test-qr.png`),
      scan & bayar, lalu `go run ./cmd/phsar payway-test <tran_id>` harus `APPROVED`. Di server: `./phsar.sh` → `docker compose run --rm app payway-test`.
- [ ] Kredensial Tuktuk Cafe yang lama (contoh PHP di Downloads/PWsample) memakai **Purchase API** (popup checkout). Phsar memakai
      **QR API** (`generate-qr`) — kalau `payway-test` menjawab kode 23/32, minta ABA mengaktifkan QR API untuk merchant tersebut.
- [ ] Kalau QR sandbox ditolak: cek pesan error di log (`payway generate-qr: kode …`) — urutan hash mengikuti dokumentasi PayWay;
      tanyakan ke tim integrasi ABA bila format berbeda (mis. `qr_image_template` wajib).
- [ ] Minta ke ABA alamat **production** API (`PAYWAY_BASE_URL`) dan kunci production; jangan pakai kunci sandbox di server.
- [ ] Beri tahu ABA URL callback: `https://<domain>/webhooks/payway` (dan whitelist IP server kalau diminta).
- [ ] Kredensial PayWay diisi super admin di **admin → Payment gateway** (merchant ID, API key, sandbox/production) — disimpan terenkripsi,
      tidak pernah tampil lagi, wajib kode authenticator. Klik **Test connection** dari server asli (IP server harus di-whitelist ABA).
- [ ] Buka belanja: di admin → Payment gateway pilih **ABA PayWay**, simpan (pita "segera dibuka" hilang otomatis, tanpa restart).
      `PAYMENT_PROVIDER` di `deploy/.env` hanya cadangan selama admin belum pernah menyimpan.
      Server menolak start kalau `fake` di production.
- [ ] `PAYMENT_LIFETIME_MIN` (default 15 menit) — lama QR berlaku sebelum pesanan batal & stok kembali.
- [ ] Pastikan worker jalan (`phsar worker` atau `WORKER_IN_SERVE=true`): worker yang membatalkan pesanan tidak dibayar.
- [ ] Uji akhir di production dengan nominal kecil ($0.10–$1): bayar → pesanan "Paid" → seller dapat notifikasi Telegram.
- [ ] Refund masih manual: tentukan siapa yang memegang menu **Refunds** di `/admin` (peran finance/super_admin) dan
      cek daftar "To refund" setiap hari. Setelah meeting dengan PayWay: putuskan apakah memakai Refund API otomatis.
- [ ] Batas seller menerima pesanan = 2 jam (`SellerResponseMinutes` di `internal/shop/fulfil.go`). Pesanan malam hari
      bisa batal saat toko tutup — tinjau ulang setelah ada data jam buka toko.
- [ ] Beri tahu seller: aktifkan login Telegram supaya pesanan baru masuk ke Telegram mereka.
- [ ] Komisi default 8% ada di tabel `commission_rules` (bisa diubah per toko).

## 11. Kurir
- [ ] Tambahkan kurir staf di `/admin` → Delivery; minta tiap kurir login di `https://<domain>/courier` dan login Telegram sekali
      (supaya notifikasi paket baru masuk ke Telegram mereka). Simpan halaman `/courier` di layar utama HP.
- [ ] Kurir sebaiknya memakai **Chrome di Android** (bisa scan barcode label). Di iPhone kurir mengetik nomor pesanan.
- [ ] Siapkan printer label (thermal 100×150 mm atau kertas biasa) di toko-toko awal; uji scan 1 label dari HP kurir.
- [ ] Paket **gagal antar**: kurir boleh coba lagi sampai 3 percobaan, lalu "Kembalikan ke toko" (kurir sendiri atau admin di
      Orders & shipping → tab "Returned to store"). Seller mengonfirmasi di Seller Center (Pesanan); tanpa konfirmasi 24 jam
      dianggap diterima (job `courierreturn.auto_receive`). Pesanan lalu batal + refund masuk antrean Refunds. Ongkir
      pengembalian ditanggung Phsar (hanya dicatat). Jelaskan alur ini ke kurir & seller awal.

## 12. Payout seller
- [ ] Buat `DATA_ENCRYPTION_KEY` baru untuk server (`openssl rand -base64 32`), simpan juga di password manager.
      **Jangan pernah diganti/hilang** — nomor rekening seller tidak bisa dibaca lagi tanpa kunci ini.
- [ ] Atur jadwal di `/admin` → Payouts (harian/mingguan, jam, masa tahan, minimal) sesuai kemampuan tim keuangan.
- [ ] Tentukan rekening Phsar yang dipakai untuk mentransfer ke seller (terpisah dari rekening cafe) dan siapa yang memegang peran finance.
- [ ] Peran finance/super_admin memegang data rekening seller — wajib MFA sebelum online (lihat bagian 9).
- [ ] Setelah meeting PayWay: cek apakah payout bisa lewat API (bukan manual).

## 12b. Komisi, dokumen legal, retur & sengketa
- [ ] Admin → **Commission rules**: cek tarif default (8%), tambahkan tarif per kategori / promo toko kalau perlu. Hanya pesanan baru yang terpengaruh.
- [ ] Admin → **Legal documents** (super admin): isi nama perusahaan, alamat, email, telepon, Telegram — tanpa ini dokumen menampilkan "—".
- [ ] Minta pengacara memeriksa 5 dokumen bawaan (syarat pembeli, privasi, retur, perjanjian seller, hapus data) di 3 bahasa, lalu **Publish new version** (versi lama tetap tersimpan sebagai bukti).
- [ ] Coba alur retur sekali: pembeli ajukan dengan foto → seller tolak → CS memutuskan di **Returns & disputes** → keuangan refund di PayWay lalu catat nomornya. Cek saldo seller terpotong.
- Aturan: seller 48 jam menjawab (lewat → otomatis ke Phsar). Selama retur terbuka, pesanan tidak masuk saldo seller.

## 12c. Menu admin lainnya (cek sekali sebelum launch)
- [ ] **Shipping & FX rates:** isi kurs riel terbaru; ongkir default tetap gratis sampai diubah.
- [ ] **Invoices & tax:** isi nama resmi perusahaan, VAT TIN, alamat; cek tarif PPN (default 10%).
- [ ] **Categories & brands:** cek nama kategori 3 bahasa; tambah kata terlarang bila perlu.
- [ ] **Content & SEO:** pasang banner beranda & judul/deskripsi SEO beranda (km/en/zh).
- [ ] **Vouchers & flash sales:** buat voucher pembukaan (mis. WELCOME, dibiayai Phsar). Flash sale: buat sesi (jam bebas), pilih siapa yang menanggung potongan (seller / Phsar / patungan), pilih banyak produk sekaligus atau buka pendaftaran untuk seller. Kalau seller ikut menanggung, produk pilihan admin harus diterima seller dulu di Seller Center → Flash sale.
- [ ] **Ads:** biarkan aktif atau matikan sampai seller cukup banyak.
- [ ] **Bank reconciliation:** tiap minggu unggah CSV transaksi dari portal PayWay.
- [ ] **Commission sweep:** catat setiap transfer komisi dari rekening penampung ke rekening perusahaan.
- [ ] **Scheduled jobs:** pastikan semua worker "Running".
- CS memakai **Customers** (cari pembeli/pesanan), **Support tickets**, **Reports & strikes**, **Payments & webhooks**.

## 13. Terakhir
- [ ] `make qa` lulus di versi yang di-deploy.
- [ ] Checklist manual di `QA.md` dijalankan sekali di server asli (HP Android + iPhone).
- [ ] `git push` versi final.
