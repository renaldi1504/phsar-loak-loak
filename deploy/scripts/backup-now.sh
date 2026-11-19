#!/usr/bin/env bash
# Satu kali backup. Dipanggil oleh backup-loop.sh, atau manual: ./phsar.sh backup
#   /backups/db/phsar-YYYYmmdd-HHMM.dump        (pg_dump format custom; dipulihkan dengan ./phsar.sh restore <file>)
#   /backups/uploads/uploads-YYYYmmdd.tar.gz     (hari Minggu, atau FORCE_UPLOADS=1)
set -euo pipefail
KEEP=${KEEP_DAYS:-14}
ts=$(date +%Y%m%d-%H%M)
mkdir -p /backups/db /backups/uploads

out=/backups/db/phsar-$ts.dump
pg_dump --format=custom --no-owner --file="$out.part"
pg_restore --list "$out.part" > /dev/null   # pastikan file bisa dibaca
mv "$out.part" "$out"
echo "backup: ✅ database → $(basename "$out") ($(du -h "$out" | cut -f1))"
find /backups/db -name 'phsar-*.dump' -mtime +"$KEEP" -delete

if [ "$(date +%u)" = 7 ] || [ "${FORCE_UPLOADS:-0}" = 1 ]; then
  up=/backups/uploads/uploads-$(date +%Y%m%d).tar.gz
  tar -czf "$up.part" -C /uploads .
  mv "$up.part" "$up"
  echo "backup: ✅ file upload → $(basename "$up") ($(du -h "$up" | cut -f1))"
  ls -1t /backups/uploads/uploads-*.tar.gz | tail -n +5 | xargs -r rm -f
fi
