#!/usr/bin/env bash
# Persiapan server baru (Linode, Ubuntu 24.04 LTS). Jalankan SEKALI sebagai root:
#   bash setup-server.sh
# Yang dilakukan: update sistem, update keamanan otomatis, jam & zona waktu, swap 2 GB, firewall (22/80/443),
# fail2ban untuk SSH, dan Docker. Tidak menyentuh kunci SSH atau user Anda.
set -euo pipefail
[ "$(id -u)" = 0 ] || { echo "Jalankan sebagai root (sudo bash setup-server.sh)"; exit 1; }
export DEBIAN_FRONTEND=noninteractive

echo "▶ Update sistem"
apt-get update -q && apt-get -yq upgrade
apt-get install -yq ca-certificates curl git ufw fail2ban unattended-upgrades openssl

echo "▶ Update keamanan otomatis"
dpkg-reconfigure -f noninteractive unattended-upgrades

echo "▶ Jam (penting untuk kode authenticator admin) & zona waktu Phnom Penh"
timedatectl set-timezone Asia/Phnom_Penh
timedatectl set-ntp true

echo "▶ Swap 2 GB (cadangan RAM untuk build & embedder)"
if ! swapon --show | grep -q .; then
  fallocate -l 2G /swapfile && chmod 600 /swapfile && mkswap /swapfile && swapon /swapfile
  grep -q '^/swapfile' /etc/fstab || echo '/swapfile none swap sw 0 0' >> /etc/fstab
  sysctl -w vm.swappiness=10 > /dev/null && echo 'vm.swappiness=10' > /etc/sysctl.d/99-phsar.conf
fi

echo "▶ Firewall: hanya SSH, HTTP, HTTPS"
ufw default deny incoming
ufw default allow outgoing
ufw allow OpenSSH
ufw allow 80/tcp
ufw allow 443/tcp
ufw allow 443/udp
ufw --force enable
systemctl enable --now fail2ban

echo "▶ Docker"
if ! command -v docker > /dev/null; then
  curl -fsSL https://get.docker.com | sh
fi
systemctl enable --now docker

echo
echo "✅ Server siap. Langkah berikutnya (lihat deploy/README.md):"
echo "   git clone https://github.com/renaldi1504/phsar-loak-loak.git /opt/phsar && cd /opt/phsar/deploy"
echo "   ./phsar.sh secrets && nano .env && ./phsar.sh check && ./phsar.sh up"
