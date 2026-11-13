# app — Flutter (pembeli + mode seller + mode kurir)

Satu app untuk semua: pembeli belanja, seller pindah ke mode seller lewat tombol di header
(muncul kalau akun punya toko), dan nanti mode kurir untuk pengiriman manual di Sihanoukville.

## Membuat project Flutter pertama kali

Folder ini baru berisi `pubspec.yaml` dan `lib/main.dart`. Folder platform (android/, ios/)
dibuat oleh Flutter sendiri:

```bash
cd app
flutter create --org com.phsar --project-name phsar --platforms android,ios .
flutter pub get
flutter run
```

`flutter create .` tidak menimpa `lib/main.dart` dan `pubspec.yaml` yang sudah ada.

## Rencana struktur (Tahap 2)

```
lib/
  main.dart
  core/        api client (dibuat dari ../api/openapi.yaml), auth, i18n en/km/zh, tema
  features/    feed, search, product, cart, checkout, orders, chat, profile, seller/, courier/
```

Tema mengikuti website Phsar: pink `#F50750`, navy `#102241`, font DM Sans + Noto Sans Khmer.
