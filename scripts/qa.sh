#!/usr/bin/env bash
# QA otomatis Phsar — jalankan dari folder utama: `make qa`
#
# Memeriksa seluruh kode (termasuk yang ditulis Claude):
#   1. format kode Go (gofmt)            5. tes integrasi semua alur API di database tes
#   2. go vet (bug umum)                 6. sintaks JavaScript halaman web (butuh Node, opsional)
#   3. go build                          7. rahasia (API key / token) yang tidak sengaja ikut git
#   4. unit test
# Hasil akhir: ✅ semua lulus, atau ❌ dengan daftar yang gagal. Exit code ≠ 0 kalau ada yang gagal.
set -uo pipefail
cd "$(dirname "$0")/.."
ROOT=$(pwd)

if [ -t 1 ]; then G=$'\e[32m'; R=$'\e[31m'; Y=$'\e[33m'; B=$'\e[1m'; N=$'\e[0m'; else G= R= Y= B= N=; fi
FAILED=()
step() { printf "\n${B}▶ %s${N}\n" "$1"; }
ok()   { printf "${G}  ✅ %s${N}\n" "$1"; }
bad()  { printf "${R}  ❌ %s${N}\n" "$1"; FAILED+=("$1"); }
warn() { printf "${Y}  ⚠️  %s${N}\n" "$1"; }

# Database tes: default di Postgres Docker (port 5433). Bisa diganti: PHSAR_TEST_DATABASE_URL=… make qa
TEST_DB_URL=${PHSAR_TEST_DATABASE_URL:-postgres://phsar:phsar@localhost:5433/phsar_test?sslmode=disable}

cd "$ROOT/backend"

step "1. Format kode Go (gofmt)"
UNFORMATTED=$(gofmt -l ./cmd ./internal ./migrations ./qa 2>/dev/null)
if [ -z "$UNFORMATTED" ]; then ok "rapi"; else bad "file belum di-gofmt (jalankan: gofmt -w backend):"; echo "$UNFORMATTED" | sed 's/^/     /'; fi

step "2. go vet"
if go vet ./... ; then ok "tidak ada temuan"; else bad "go vet menemukan masalah"; fi

step "3. go build"
if go build -o /dev/null ./cmd/phsar ; then ok "binary bisa dibuat"; else bad "build gagal"; fi

step "4. Unit test"
if go test -count=1 ./internal/... ; then ok "unit test lulus"; else bad "unit test gagal"; fi

step "5. Tes integrasi (login, toko, KYC, titik penjemputan, produk, terjemahan, moderasi, admin)"
# Database tes (phsar_test) dibuat otomatis oleh tesnya kalau belum ada; isinya DIHAPUS setiap run.
if PHSAR_TEST_DATABASE_URL="$TEST_DB_URL" go test -count=1 ${QA_VERBOSE:+-v} ./qa/ ; then ok "semua alur lulus"; else bad "tes integrasi gagal (lihat pesan di atas)"; fi

if ! command -v ffmpeg >/dev/null || ! command -v ffprobe >/dev/null; then
  warn "ffmpeg tidak terpasang — tes kompres video dilewati, dan server memakai video asli (pasang: brew install ffmpeg)"
fi

step "6. Sintaks JavaScript halaman web"
if command -v node >/dev/null; then
  JS_OK=1
  for f in internal/webui/*.js internal/storefront/assets/*.js; do
    if ! out=$(node --check "$f" 2>&1); then echo "$out" | sed 's/^/     /'; JS_OK=0; fi
  done
  # Script halaman dimuat bersama (berbagi ruang nama global): nama ganda (mis. dua `const money`) = halaman mati.
  for page in internal/webui/seller.html internal/webui/admin.html internal/webui/courier.html; do
    tmpd=$(mktemp -d); tmp="$tmpd/page.js"
    for js in $(grep -o 'src="/ui/[a-z-]*\.js"' "$page" | sed 's|src="/ui/||; s|"||' | grep -v '^config.js$'); do
      cat "internal/webui/$js" >> "$tmp"; echo >> "$tmp"
    done
    if ! out=$(node --check "$tmp" 2>&1); then echo "     $(basename "$page"): $(echo "$out" | grep -m1 -E 'Error')"; JS_OK=0; fi
    rm -rf "$tmpd"
  done
  if [ $JS_OK = 1 ]; then ok "$(ls internal/webui/*.js internal/storefront/assets/*.js | wc -l | tr -d ' ') file JS valid (juga digabung per halaman)"; else bad "ada file JS dengan error sintaks / nama ganda"; fi
else
  warn "Node.js tidak terpasang — lewati (pasang: brew install node)"
fi

step "7. Rahasia yang tidak sengaja ikut git"
cd "$ROOT"
PATTERNS='sk-ant-[A-Za-z0-9_-]{20,}|[0-9]{8,10}:AA[A-Za-z0-9_-]{30,}|AIza[0-9A-Za-z_-]{35}|gh[pousr]_[A-Za-z0-9]{30,}|github_pat_[A-Za-z0-9_]{30,}|-----BEGIN [A-Z ]*PRIVATE KEY-----'
if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  LEAKS=$(git ls-files -co --exclude-standard | grep -v '^\.env$' | xargs grep -EIl "$PATTERNS" 2>/dev/null)
  if git ls-files --error-unmatch .env >/dev/null 2>&1; then LEAKS="$LEAKS .env (file .env ikut git!)"; fi
  if [ -z "$LEAKS" ]; then ok "tidak ada API key / token di file yang ikut git"; else bad "kemungkinan rahasia di: $LEAKS"; fi
else
  warn "bukan repo git — lewati"
fi

echo
if [ ${#FAILED[@]} -eq 0 ]; then
  printf "${G}${B}✅ QA LULUS — semua pemeriksaan otomatis berhasil.${N}\n"
  printf "   Sebelum rilis, lanjutkan checklist manual di QA.md.\n"
  exit 0
fi
printf "${R}${B}❌ QA GAGAL (%d):${N}\n" "${#FAILED[@]}"
for f in "${FAILED[@]}"; do printf "${R}   - %s${N}\n" "$f"; done
exit 1
