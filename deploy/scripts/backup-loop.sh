#!/usr/bin/env bash
# Berjalan terus di container "backup": setiap hari jam BACKUP_HOUR membuat backup, lalu tidur sampai besok.
set -uo pipefail
HOUR=${BACKUP_HOUR:-3}
echo "backup: aktif — dump database setiap hari jam ${HOUR}:00 ($(date +%Z)), file upload setiap Minggu"
while true; do
  now=$(date +%s)
  next=$(date -d "today ${HOUR}:00" +%s)
  [ "$next" -le "$now" ] && next=$(date -d "tomorrow ${HOUR}:00" +%s)
  sleep $((next - now))
  bash /scripts/backup-now.sh || echo "backup: ❌ GAGAL $(date '+%F %T')"
done
