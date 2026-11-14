# QA Phsar

Dua lapis pemeriksaan untuk setiap perubahan, termasuk kode yang ditulis Claude:

1. **Otomatis**: `make qa`. Wajib lulus sebelum `git push`; juga jalan di GitHub (CI) setiap push.
2. **Manual**: checklist di bawah, untuk hal yang harus dilihat atau dicoba manusia (peta, Telegram, video, tampilan HP).

---

## 1. QA otomatis (`make qa`)

```bash
make db      # Postgres Docker harus menyala
make qa
```

| # | Pemeriksaan | Menangkap |
|---|---|---|
| 1 | `gofmt` | kode Go tidak rapi |
| 2 | `go vet` | bug umum Go (format string salah, kode tak terjangkau, dll.) |
| 3 | `go build` | kode tidak bisa dikompilasi |
| 4 | Unit test | perhitungan uang/komisi, normalisasi nomor HP, cek wilayah, klien Claude/Google Translate |
| 5 | **Tes integrasi** (`backend/qa/`) | semua alur API dijalankan sungguhan di database `phsar_test`. Lihat daftar di bawah |
| 6 | `node --check` | error sintaks JavaScript di Seller Center/admin (butuh Node: `brew install node`) |
| 7 | Cek rahasia | API key Anthropic/Google, token Telegram/GitHub, private key yang tidak sengaja masuk git; `.env` ikut git |

Hasil akhir: **✅ QA LULUS** atau **❌ QA GAGAL** dengan daftar yang gagal.

### Isi tes integrasi

Server sungguhan + database tes kosong. Migrasi juga dicek naik → turun semua → naik lagi.
Setiap balasan **500** dari server di endpoint mana pun otomatis membuat QA gagal (500 = bug).

- **Login & sesi**: minta OTP, normalisasi `+855`, jeda 60 detik, kode salah dijawab cepat dan dihitung,
  kode tidak bisa dipakai dua kali, `/me`, refresh berputar dan refresh lama ditolak, logout, CSRF header web.
- **Toko**: validasi, toko orang lain → 404, deskripsi diterjemahkan, upload dicek dari isi file (HTML berkedok JPG ditolak),
  file orang lain tidak bisa dipakai, dokumen KYC tidak publik, titik di luar Sihanoukville ditolak, **ubah titik penjemputan**
  (label saja tetap terverifikasi; lokasi pindah → harus dicek ulang), maks. 5 titik, kirim untuk review, toko terkunci saat direview.
- **Admin**: bukan admin / peran tanpa izin ditolak, antrean, buka dokumen KYC (tercatat di audit log), toko belum bisa disetujui
  sebelum KYC & titik terverifikasi, tolak wajib catatan, toko jadi aktif.
- **Produk**: 15 kasus validasi (nama, kategori, kosmetik bekas, nilai varian, harga 0, harga coret, stok negatif, varian dobel,
  foto sebagai video, dua video, kirim tanpa video/foto), SKU otomatis, cache harga & stok, terjemahan nama/deskripsi/opsi/nilai,
  slug unik, toko lain tidak bisa melihat/memakai upload, ubah varian (ID lama tetap, yang dihapus hilang, pergerakan stok tercatat),
  SKU bentrok, kata terlarang ditandai ("gun" tidak salah cocok dengan "burgundy"), terjemahan editan seller tidak ditimpa AI,
  moderasi (minta perbaikan → kirim ulang → setujui, riwayat), ubah harga/stok tetap tayang, ubah foto → dicek ulang,
  arsip, barang bekas wajib foto kondisi, toko belum aktif hanya boleh draf, hapus.
- **Notifikasi**: keputusan admin → notifikasi untuk anggota toko yang tepat (owner/staf), teks 3 bahasa, terkirim ke Telegram
  (dengan link), gagal sementara → coba lagi dengan jeda, diblokir → berhenti, tanpa Telegram → tetap di lonceng, tandai dibaca.
- **Staf toko**: link undangan (hash di database, sekali pakai, batal, kedaluwarsa), staf hanya bisa produk, owner tidak bisa dihapus,
  hapus staf / staf keluar → akses hilang, notifikasi staf baru ke owner, banner toko.
- **Media** (kalau ffmpeg ada): video >3 menit & file berkedok video ditolak saat upload, video dikompres (≤1280 px) + poster,
  foto diperkecil (≤1600 px), sampul daftar produk memakai foto kecil, ubah urutan tidak memproses ulang, backfill media lama.
- **Web & publik**: halaman `/seller` `/admin` dan semua file JS, `config.js` (kunci Maps), kategori 3 bahasa,
  cek wilayah, folder media tidak bisa di-list, header keamanan, 404.

Lihat log server saat tes: `QA_VERBOSE=1 make qa`.

### Aturan untuk setiap fitur baru

Setiap fitur atau bug fix **harus menambah tes** di `backend/qa/flows_test.go` (atau unit test di paketnya),
minimal satu skenario sukses dan satu skenario gagal. Bug yang ditemukan manual → tulis tesnya dulu, baru diperbaiki.
Kalau minta Claude membuat fitur, tambahkan kalimat: *"sertakan tes QA-nya dan pastikan `make qa` lulus"*.

---

## 2. Checklist manual

Centang sebelum rilis atau setelah perubahan besar. Coba di laptop **dan** HP (Chrome Android / Safari iPhone),
dan minimal sekali dalam bahasa ខ្មែរ dan 中文.

### Login
- [ ] Login Telegram: tombol membuka bot, tekan Start, halaman otomatis masuk
- [ ] Login SMS: kode masuk (production: SMS sungguhan sampai ke HP Kamboja Smart/Cellcard/Metfone)
- [ ] Refresh halaman → tetap login; Sign out → kembali ke halaman login

### Seller Center
- [ ] Tekan Simpan/Kirim dengan isian kosong → **kotak merah di atas tombol berisi daftar semua yang kurang** (klik → lompat ke isian),
      isian ditandai merah, kursor di isian pertama; kotak hilang setelah semua diperbaiki
- [ ] Toko belum disetujui → tombol "Send for review" nonaktif dengan penjelasan; "Save draft" tetap bisa
- [ ] Buat toko, upload logo (foto dari HP), deskripsi muncul dalam 3 bahasa setelah beberapa detik
- [ ] KYC: upload foto KTP dari kamera HP
- [ ] Titik penjemputan: cari nama toko → pilih hasil → pin & alamat terisi; geser pin; "Gunakan lokasi saya" di HP
- [ ] Titik di luar Sihanoukville → muncul "di luar wilayah"
- [ ] Ubah titik penjemputan → tersimpan
- [ ] Tambah produk: video dari HP (termasuk .mov iPhone) bisa diputar di preview; foto; warna × ukuran; isi semua harga sekaligus
- [ ] Barang bekas: kolom catatan & foto kondisi muncul
- [ ] Video dari HP (termasuk .mov iPhone, video potret) → badge "Compressing video…" lalu "optimised"; hasilnya tidak miring/terbalik
- [ ] Foto dari HP tidak miring setelah diperkecil
- [ ] Setelah kirim: status "In review", terjemahan km/zh muncul, bisa diedit
- [ ] Bahasa Khmer: teks tidak terpotong / tumpang tindih di tombol & tabel

### Notifikasi & tim
- [ ] Admin menyetujui toko/produk → lonceng di Seller Center bertambah; seller yang login dengan Telegram menerima pesan dari bot
- [ ] Klik notifikasi → membuka halaman toko/produk yang benar
- [ ] Owner membuat link undangan → staf membuka di HP, login, "Join store" → hanya tab Produk yang terlihat
- [ ] Owner menghapus staf → staf tidak bisa membuka toko lagi
- [ ] Edit stok & harga dari daftar produk → langsung tersimpan tanpa review

### Admin
- [ ] Antrean toko → buka dokumen KYC → setujui KYC, titik, toko
- [ ] Peta titik penjemputan tampil; "Open in Google Maps" membuka tempat yang benar
- [ ] Antrean produk: video bisa diputar, terjemahan 3 bahasa, tanda kata terlarang, Setujui / Minta perbaikan / Tolak
- [ ] Seller melihat catatan moderator dan bisa kirim ulang

### Terjemahan AI (dengan `ANTHROPIC_API_KEY` asli)
- [ ] Log server: `terjemahan selesai … model=claude-…` (bukan `fake`)
- [ ] Minta orang Khmer & Cina membaca 5 produk contoh: wajar, merek/ukuran tidak ikut diterjemahkan

### Keamanan & data
- [ ] `.env` tidak ikut `git status` / GitHub
- [ ] Key Google Maps dibatasi per domain di Google Cloud
- [ ] Backup database bisa dipulihkan (sebelum production)
