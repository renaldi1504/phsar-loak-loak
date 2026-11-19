'use strict';
// Service worker Phsar — hanya Web Push (tanpa cache/fetch, jadi tidak mengubah cara halaman dimuat).
// Dilayani di /seller/sw.js, /admin/sw.js dan /sw.js (website pembeli); cakupan = folder masing-masing.
// Kalau halaman aplikasi yang sama sedang terlihat, halaman itu sendiri yang menampilkan + membunyikan notifikasi.
const SCOPE = new URL(self.registration.scope).pathname;
const ICON = '/ui/icon-192.png';

self.addEventListener('install', () => self.skipWaiting());
self.addEventListener('activate', e => e.waitUntil(self.clients.claim()));

self.addEventListener('push', e => {
  let d = {};
  try { d = e.data ? e.data.json() : {}; } catch { d = {title: 'Phsar', body: e.data ? e.data.text() : ''}; }
  e.waitUntil((async () => {
    const wins = await self.clients.matchAll({type: 'window', includeUncontrolled: true});
    const visible = wins.some(c => c.visibilityState === 'visible' && new URL(c.url).pathname.startsWith(SCOPE));
    if (visible && !d.test) return;
    const urgent = d.sound === 'order' || d.sound === 'alert';
    await self.registration.showNotification(d.title || 'Phsar', {
      body: d.body || '',
      tag: d.tag || undefined,
      icon: ICON,
      badge: ICON,
      data: {url: d.url || SCOPE},
      silent: !!d.silent,
      renotify: !!d.tag && !d.silent,
      requireInteraction: d.sound === 'order', // pesanan baru tetap di layar sampai diklik
      vibrate: d.silent ? undefined : (urgent ? [200, 100, 200, 100, 200] : [120]),
    });
  })());
});

self.addEventListener('notificationclick', e => {
  e.notification.close();
  const target = new URL((e.notification.data && e.notification.data.url) || SCOPE, self.location.origin).href;
  e.waitUntil((async () => {
    const wins = await self.clients.matchAll({type: 'window', includeUncontrolled: true});
    const same = wins.find(c => new URL(c.url).pathname.startsWith(SCOPE));
    if (same) {
      await same.focus();
      if (same.url !== target && 'navigate' in same) return same.navigate(target);
      return same;
    }
    return self.clients.openWindow(target);
  })());
});
