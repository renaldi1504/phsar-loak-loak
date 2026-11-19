#!/usr/bin/env bash
# Perintah harian Phsar di server. Jalankan dari folder deploy/:  ./phsar.sh <perintah>
set -euo pipefail
cd "$(dirname "$0")"
DC="docker compose"

say()  { printf '%s\n' "$*"; }
die()  { printf '❌ %s\n' "$*" >&2; exit 1; }
envv() { grep -E "^$1=" .env 2>/dev/null | tail -1 | cut -d= -f2- | sed -e 's/^["'\'']//' -e 's/["'\'']$//'; }

usage() {
  cat <<'EOF'
Pemakaian: ./phsar.sh <perintah>

  secrets              isi otomatis nilai acak yang masih kosong di .env (JWT, password DB, kunci webhook, kunci enkripsi)
  check                periksa .env + DNS sebelum online
  up                   build lalu jalankan semuanya (juga dipakai setelah mengubah .env)
  update               ambil versi terbaru dari git, build, jalankan ulang (migrasi otomatis)
  status               daftar layanan + kesehatan
  logs [layanan]       lihat log (app, worker, caddy, db, embedder, backup); Ctrl+C untuk keluar
  admin-grant <hp> <peran>   jadikan admin, mis.  ./phsar.sh admin-grant 012345678 super_admin
  reset-mfa <hp>       reset aplikasi authenticator admin yang kehilangan HP
  backup               backup sekarang (database + file upload)
  backups              daftar file backup
  restore <file.dump>  pulihkan database dari backup (menimpa data sekarang!)
  psql                 buka konsol database
  translate-test       uji kunci Anthropic
  otp [nomor HP]       (staging, SMS_PROVIDER=log) tampilkan kode login terakhir dari log
  media-backfill       kompres video & foto lama
  stop                 hentikan semua layanan
EOF
}

need_env() { [ -f .env ] || die "deploy/.env belum ada. Jalankan:  cp env.production.example .env  lalu  ./phsar.sh secrets"; }

cmd_secrets() {
  [ -f .env ] || cp env.production.example .env
  chmod 600 .env
  fill() { # fill KEY VALUE — hanya kalau masih kosong
    local k=$1 v=$2
    if grep -qE "^$k=.+" .env; then say "  $k sudah terisi — tidak diubah"; return; fi
    if grep -qE "^$k=" .env; then sed -i "s|^$k=.*|$k=$v|" .env; else printf '%s=%s\n' "$k" "$v" >> .env; fi
    say "  $k diisi"
  }
  fill POSTGRES_PASSWORD "$(openssl rand -hex 24)"
  fill JWT_SECRET "$(openssl rand -hex 32)"
  fill TELEGRAM_WEBHOOK_KEY "$(openssl rand -hex 16)"
  fill DATA_ENCRYPTION_KEY "$(openssl rand -base64 32)"
  say "✅ Selesai. Salin isi .env ke password manager (terutama DATA_ENCRYPTION_KEY dan POSTGRES_PASSWORD)."
}

cmd_check() {
  need_env
  local bad=0 d env req
  env=$(envv APP_ENV); env=${env:-production}
  case "$env" in
    production) say "ℹ️  APP_ENV=production";;
    staging) say "ℹ️  APP_ENV=staging — server TESTING (bukan untuk pembeli sungguhan)";;
    *) say "❌ APP_ENV harus production atau staging"; bad=1;;
  esac
  req="PHSAR_DOMAIN ACME_EMAIL PUBLIC_BASE_URL POSTGRES_PASSWORD JWT_SECRET DATA_ENCRYPTION_KEY"
  # Staging dengan SMS_PROVIDER=log boleh tanpa Telegram (kode login dibaca dari ./phsar.sh otp).
  if [ "$env" = production ] || [ "$(envv SMS_PROVIDER)" != log ] || [ -n "$(envv TELEGRAM_BOT_TOKEN)" ]; then
    req="$req TELEGRAM_BOT_TOKEN TELEGRAM_BOT_USERNAME TELEGRAM_WEBHOOK_KEY"
  fi
  for k in $req; do
    if [ -z "$(envv "$k")" ]; then say "❌ $k kosong"; bad=1; else say "✅ $k"; fi
  done
  d=$(envv PHSAR_DOMAIN)
  [ "$(envv PUBLIC_BASE_URL)" = "https://$d" ] || { say "❌ PUBLIC_BASE_URL harus https://$d"; bad=1; }
  [[ "$(envv POSTGRES_PASSWORD)" =~ ^[A-Za-z0-9]+$ ]] || { say "❌ POSTGRES_PASSWORD hanya boleh huruf/angka"; bad=1; }
  if [ "$(envv SMS_PROVIDER)" = "log" ]; then
    [ "$env" = staging ] && say "ℹ️  SMS_PROVIDER=log — kode login TIDAK dikirim, lihat dengan ./phsar.sh otp" || { say "❌ SMS_PROVIDER=log tidak boleh di production (pakai off)"; bad=1; }
  fi
  case "$(envv PAYMENT_PROVIDER)" in
    off) say "ℹ️  PAYMENT_PROVIDER=off — mode pra-peluncuran (checkout ditutup, seller bisa daftar)";;
    fake) [ "$env" = staging ] && say "ℹ️  PAYMENT_PROVIDER=fake — tombol \"simulasi bayar\" (tanpa uang sungguhan)" || { say "❌ PAYMENT_PROVIDER=fake hanya untuk APP_ENV=staging"; bad=1; };;
    payway) for k in PAYWAY_BASE_URL PAYWAY_MERCHANT_ID PAYWAY_API_KEY; do [ -n "$(envv $k)" ] || { say "❌ $k kosong"; bad=1; }; done;;
    *) say "❌ PAYMENT_PROVIDER harus off, payway, atau fake (staging)"; bad=1;;
  esac
  [ -n "$(envv ANTHROPIC_API_KEY)$(envv GEMINI_API_KEY)" ] || say "ℹ️  Kunci AI belum ada di .env — isi nanti di Admin → Integrasi → AI & terjemahan"
  [ -n "$(envv GOOGLE_MAPS_API_KEY)" ] || say "⚠️  GOOGLE_MAPS_API_KEY kosong — peta memakai OpenStreetMap biasa"
  if [ -n "$d" ]; then
    local me dns
    me=$(curl -fsS4 --max-time 5 https://api.ipify.org 2>/dev/null || true)
    dns=$(getent ahostsv4 "$d" 2>/dev/null | awk 'NR==1{print $1}' || true)
    if [ -z "$dns" ]; then say "❌ DNS $d belum mengarah ke mana pun (buat A record ke ${me:-IP server})"; bad=1
    elif [ -n "$me" ] && [ "$dns" != "$me" ]; then say "⚠️  DNS $d → $dns, tapi IP server ini $me (benar kalau pakai Cloudflare proxy)"
    else say "✅ DNS $d → $dns"; fi
    local sub sd
    for sub in seller admin api; do
      sd=$(getent ahostsv4 "$sub.$d" 2>/dev/null | awk 'NR==1{print $1}' || true)
      if [ -z "$sd" ]; then say "⚠️  DNS $sub.$d belum ada — buat A record $sub → ${me:-IP server} (atau wildcard *)"
      else say "✅ DNS $sub.$d → $sd"; fi
    done
  fi
  [ $bad = 0 ] && say "✅ Siap: ./phsar.sh up" || die "Perbaiki dulu yang bertanda ❌ di deploy/.env"
}

cmd_up()      { need_env; mkdir -p backups; $DC up -d --build --remove-orphans; sleep 5; cmd_status; }
cmd_update()  { need_env; git -C .. pull --ff-only; cmd_up; docker image prune -f > /dev/null; }
cmd_status()  { $DC ps; }
cmd_logs()    { $DC logs -f --tail=200 "$@"; }
cmd_stop()    { $DC stop; }
run_app()     { need_env; $DC run --rm --no-deps -T app "$@"; }

cmd_restore() {
  local f=${1:-}
  [ -n "$f" ] || die "Pakai:  ./phsar.sh restore backups/db/phsar-YYYYmmdd-HHMM.dump"
  [ -f "$f" ] || die "File tidak ada: $f"
  case "$(realpath "$f")" in "$(realpath backups)"/*) ;; *) die "File backup harus ada di folder deploy/backups/";; esac
  say "⚠️  Database sekarang akan DITIMPA dengan $f."
  read -r -p "Ketik PULIHKAN untuk lanjut: " ok
  [ "$ok" = "PULIHKAN" ] || die "Dibatalkan."
  $DC stop app worker
  $DC exec -T backup pg_restore --clean --if-exists --no-owner --dbname=phsar "/backups/${f#*backups/}"
  $DC up -d app worker
  say "✅ Database dipulihkan."
}

case "${1:-}" in
  secrets) cmd_secrets;;
  check) cmd_check;;
  up) cmd_up;;
  update) cmd_update;;
  status|ps) cmd_status;;
  logs) shift; cmd_logs "$@";;
  stop) cmd_stop;;
  admin-grant) [ $# -eq 3 ] || die "Pakai: ./phsar.sh admin-grant <nomor HP> <super_admin|moderator|cs|finance|ads_reviewer>"; run_app admin grant "$2" "$3";;
  reset-mfa) [ $# -eq 2 ] || die "Pakai: ./phsar.sh reset-mfa <nomor HP>"; run_app admin reset-mfa "$2";;
  backup) $DC exec -T -e FORCE_UPLOADS=1 backup bash /scripts/backup-now.sh;;
  backups) ls -lh backups/db backups/uploads 2>/dev/null || say "Belum ada backup.";;
  restore) cmd_restore "${2:-}";;
  psql) $DC exec db psql -U phsar phsar;;
  translate-test) run_app translate-test;;
  otp) [ "$(envv SMS_PROVIDER)" = log ] || die "Hanya untuk staging dengan SMS_PROVIDER=log"
       $DC logs --no-color --since 30m app | grep -a "SMS (mode log" | grep -a -- "${2#0}" | tail -5 || say "Belum ada kode dalam 30 menit terakhir.";;
  media-backfill) run_app media-backfill;;
  *) usage;;
esac
