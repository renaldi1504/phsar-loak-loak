'use strict';
// Notifikasi web untuk Seller Center & admin: lonceng, realtime (SSE lewat fetch), bunyi (WebAudio),
// notifikasi browser (Web Push + service worker) dan pengaturan per grup. Dimuat setelah common.js.
//   Seller: /v1/seller/notifications   Admin: /admin/api/notifications (dibagi per izin peran di server)
//   Pesanan baru di seller berbunyi 3× (server mengirim sound/repeat/volume sesuai pengaturan user).
addStrings({
  en: {notifications: 'Notifications', noNotifications: 'No notifications yet.', markAllRead: 'Mark all as read',
    tgHint: 'Sign in with Telegram to also get these in Telegram.',
    nfSettings: 'Notification settings', nfSoundOn: 'Turn on notification sound', nfSoundHint: 'Your browser needs one click before it can play sounds.',
    nfThisBrowser: 'This browser', nfPush: 'Browser notifications', nfPushHint: 'Shows notifications even when this tab is closed or in the background.',
    nfPushOn: 'On for this browser', nfPushOff: 'Off', nfPushDenied: 'Blocked in browser settings — allow notifications for this site, then try again.',
    nfPushNo: 'Not supported here. On iPhone/iPad: add this page to the Home Screen first (Share → Add to Home Screen).',
    nfEnable: 'Turn on', nfDisable: 'Turn off', nfTestSound: 'Test sound', nfTestSend: 'Send test notification', nfTestSent: 'Test sent — it should arrive in a moment.',
    nfVolume: 'Volume', nfGroup: 'Type', nfSound: 'Sound', nfPushCol: 'Browser', nfTelegram: 'Telegram', nfTgNotLinked: 'Telegram: sign in with Telegram to use this.',
    nfQuiet: 'Quiet hours (no sound)', nfQuietFrom: 'From', nfQuietTo: 'To', nfQuietHint: 'Notifications still arrive in the bell — just silently. Cambodia time.',
    nfSaved: 'Notification settings saved', nfSave: 'Save', nfOrderRepeat: 'New orders ring 3 times.', nfOpen: 'Open', nfAskPush: 'Get alerts even when this tab is in the background',
    ng_orders: 'Orders', ng_chat: 'Chat', ng_returns: 'Returns', ng_payments: 'Payouts & invoices', ng_store: 'Store & products', ng_promo: 'Promotions & ads',
    ng_stores: 'Store verification', ng_products: 'Product review', ng_refunds: 'Refunds', ng_reports: 'Buyer reports', ng_support: 'Support tickets',
    ng_payouts: 'Payouts', ng_deliveries: 'Deliveries'},
  km: {notifications: 'ការជូនដំណឹង', noNotifications: 'មិនទាន់មានការជូនដំណឹងទេ។', markAllRead: 'សម្គាល់ថាបានអានទាំងអស់',
    tgHint: 'ចូលតាម Telegram ដើម្បីទទួលការជូនដំណឹងក្នុង Telegram ផងដែរ។',
    nfSettings: 'ការកំណត់ការជូនដំណឹង', nfSoundOn: 'បើកសំឡេងជូនដំណឹង', nfSoundHint: 'កម្មវិធីរុករកត្រូវការចុចម្តងសិន ទើបអាចបន្លឺសំឡេងបាន។',
    nfThisBrowser: 'កម្មវិធីរុករកនេះ', nfPush: 'ការជូនដំណឹងក្នុងកម្មវិធីរុករក', nfPushHint: 'បង្ហាញការជូនដំណឹង ទោះបីផ្ទាំងនេះបិទ ឬនៅខាងក្រោយក៏ដោយ។',
    nfPushOn: 'បើកសម្រាប់កម្មវិធីរុករកនេះ', nfPushOff: 'បិទ', nfPushDenied: 'ត្រូវបានរារាំងក្នុងការកំណត់កម្មវិធីរុករក — អនុញ្ញាតការជូនដំណឹងសម្រាប់គេហទំព័រនេះ រួចសាកម្តងទៀត។',
    nfPushNo: 'មិនគាំទ្រនៅទីនេះទេ។ លើ iPhone/iPad៖ បន្ថែមទំព័រនេះទៅអេក្រង់ដើមជាមុន (ចែករំលែក → បន្ថែមទៅអេក្រង់ដើម)។',
    nfEnable: 'បើក', nfDisable: 'បិទ', nfTestSound: 'សាកសំឡេង', nfTestSend: 'ផ្ញើការជូនដំណឹងសាកល្បង', nfTestSent: 'បានផ្ញើការសាកល្បង — វានឹងមកដល់បន្តិចទៀត។',
    nfVolume: 'កម្រិតសំឡេង', nfGroup: 'ប្រភេទ', nfSound: 'សំឡេង', nfPushCol: 'កម្មវិធីរុករក', nfTelegram: 'Telegram', nfTgNotLinked: 'Telegram៖ សូមចូលតាម Telegram ដើម្បីប្រើមុខងារនេះ។',
    nfQuiet: 'ម៉ោងស្ងាត់ (គ្មានសំឡេង)', nfQuietFrom: 'ពី', nfQuietTo: 'ដល់', nfQuietHint: 'ការជូនដំណឹងនៅតែមកដល់ក្នុងកណ្តឹង — គ្រាន់តែគ្មានសំឡេង។ ម៉ោងកម្ពុជា។',
    nfSaved: 'បានរក្សាទុកការកំណត់ការជូនដំណឹង', nfSave: 'រក្សាទុក', nfOrderRepeat: 'ការបញ្ជាទិញថ្មីរោទ៍ 3 ដង។', nfOpen: 'បើក', nfAskPush: 'ទទួលការជូនដំណឹង ទោះផ្ទាំងនេះនៅខាងក្រោយក៏ដោយ',
    ng_orders: 'ការបញ្ជាទិញ', ng_chat: 'ជជែក', ng_returns: 'ការប្រគល់វិញ', ng_payments: 'ការទូទាត់ និងវិក្កយបត្រ', ng_store: 'ហាង និងផលិតផល', ng_promo: 'ប្រូម៉ូសិន និងការផ្សាយ',
    ng_stores: 'ការផ្ទៀងផ្ទាត់ហាង', ng_products: 'ការពិនិត្យផលិតផល', ng_refunds: 'ការសងប្រាក់', ng_reports: 'របាយការណ៍ពីអ្នកទិញ', ng_support: 'សំបុត្រជំនួយ',
    ng_payouts: 'ការទូទាត់ឲ្យអ្នកលក់', ng_deliveries: 'ការដឹកជញ្ជូន'},
  zh: {notifications: '通知', noNotifications: '暂无通知。', markAllRead: '全部标为已读',
    tgHint: '使用 Telegram 登录即可在 Telegram 中收到这些通知。',
    nfSettings: '通知设置', nfSoundOn: '开启通知提示音', nfSoundHint: '浏览器需要先点击一次才能播放声音。',
    nfThisBrowser: '此浏览器', nfPush: '浏览器通知', nfPushHint: '即使此标签页已关闭或在后台，也会显示通知。',
    nfPushOn: '此浏览器已开启', nfPushOff: '已关闭', nfPushDenied: '已在浏览器设置中被阻止——请允许本网站发送通知后重试。',
    nfPushNo: '此处不支持。iPhone/iPad：请先将本页添加到主屏幕（分享 → 添加到主屏幕）。',
    nfEnable: '开启', nfDisable: '关闭', nfTestSound: '试听', nfTestSend: '发送测试通知', nfTestSent: '测试已发送，稍后即会收到。',
    nfVolume: '音量', nfGroup: '类型', nfSound: '声音', nfPushCol: '浏览器', nfTelegram: 'Telegram', nfTgNotLinked: 'Telegram：请使用 Telegram 登录后再开启。',
    nfQuiet: '免打扰时段（静音）', nfQuietFrom: '从', nfQuietTo: '到', nfQuietHint: '通知仍会进入铃铛，只是不响。柬埔寨时间。',
    nfSaved: '通知设置已保存', nfSave: '保存', nfOrderRepeat: '新订单会响 3 次。', nfOpen: '打开', nfAskPush: '标签页在后台时也接收提醒',
    ng_orders: '订单', ng_chat: '聊天', ng_returns: '退货', ng_payments: '结算与发票', ng_store: '店铺与商品', ng_promo: '促销与广告',
    ng_stores: '店铺审核', ng_products: '商品审核', ng_refunds: '退款', ng_reports: '买家举报', ng_support: '客服工单',
    ng_payouts: '结算', ng_deliveries: '配送'},
});

const NF_ADMIN = location.pathname.startsWith('/admin');
const NF_BASE = NF_ADMIN ? '/admin/api/notifications' : '/v1/seller/notifications';
const NF_SCOPE = NF_ADMIN ? '/admin/' : '/seller/';
// Admin baru boleh memakai notifikasi setelah lolos MFA.
const nfAuthed = () => !!P.token && (!NF_ADMIN || !!(P.mfa && P.mfa.verified));
const NB = {open: false, data: null, unread: 0, prefs: null, pushOn: false, title: document.title};

// ---- Lonceng ----
function bellButton() {
  nfLive.ensure();
  const n = NB.unread;
  return `<div class="bellwrap"><button type="button" class="bell" id="bellBtn" aria-label="${t('notifications')}" aria-expanded="${NB.open}" onclick="toggleBell()">
    <svg width="20" height="20" viewBox="0 0 24 24" aria-hidden="true"><path fill="currentColor" d="M12 22a2.5 2.5 0 0 0 2.4-2h-4.8A2.5 2.5 0 0 0 12 22Zm7-6V11a7 7 0 0 0-5.5-6.8V3a1.5 1.5 0 0 0-3 0v1.2A7 7 0 0 0 5 11v5l-2 2v1h18v-1z"/></svg>
    <span class="badge ${n ? '' : 'hidden'}" id="bellCount">${n > 9 ? '9+' : n}</span></button><div class="bellpanel hidden" id="bellPanel" role="dialog" aria-label="${t('notifications')}"></div></div>`;
}
function nfSetUnread(n) {
  NB.unread = n || 0;
  const c = $('#bellCount');
  if (c) { c.textContent = NB.unread > 9 ? '9+' : NB.unread; c.classList.toggle('hidden', !NB.unread); }
  const base = document.title.replace(/^\(\d+\+?\) /, '');
  document.title = NB.unread ? `(${NB.unread > 99 ? '99+' : NB.unread}) ${base}` : base;
}
async function loadBell() {
  if (!nfAuthed()) return;
  try { NB.data = await api(NF_BASE + '?limit=30'); } catch { return; }
  nfSetUnread(NB.data.unread);
  if (NB.open) drawBell();
}
function drawBell() {
  const panel = $('#bellPanel');
  if (!panel || !NB.data) return;
  panel.classList.toggle('hidden', !NB.open);
  const list = NB.data.notifications;
  const askPush = nfPush.supported() && !NB.pushOn && Notification.permission !== 'denied';
  panel.innerHTML = `<div class="spread"><b>${t('notifications')}</b><span class="nf-acts">${NB.data.unread ? `<button type="button" class="linkbtn" onclick="readAll()">${t('markAllRead')}</button>` : ''}
      <button type="button" class="nf-gear" onclick="nfOpenSettings()" aria-label="${t('nfSettings')}" title="${t('nfSettings')}"><svg width="18" height="18" viewBox="0 0 24 24" aria-hidden="true"><path fill="currentColor" d="M19.4 13a7.5 7.5 0 0 0 0-2l2-1.6-2-3.4-2.4 1a7.6 7.6 0 0 0-1.7-1L15 3.5h-4L10.7 6a7.6 7.6 0 0 0-1.7 1l-2.4-1-2 3.4 2 1.6a7.5 7.5 0 0 0 0 2l-2 1.6 2 3.4 2.4-1a7.6 7.6 0 0 0 1.7 1l.3 2.5h4l.3-2.5a7.6 7.6 0 0 0 1.7-1l2.4 1 2-3.4zM12 15.5a3.5 3.5 0 1 1 0-7 3.5 3.5 0 0 1 0 7z"/></svg></button></span></div>
    ${!nfSound.ready() ? `<button type="button" class="nf-cta" onclick="nfSound.unlock(true)">🔔 ${t('nfSoundOn')}</button>` : ''}
    ${askPush ? `<button type="button" class="nf-cta soft" onclick="nfOpenSettings()">📲 ${t('nfAskPush')}</button>` : ''}
    ${list.length ? list.map(n => `<a class="notif ${n.read ? '' : 'unread'}" href="${esc(n.link || '#/')}" onclick="readOne(${n.id})">
      <b>${esc(n.title)}</b>${n.body ? `<span>${esc(n.body)}</span>` : ''}<small class="muted">${new Date(n.created_at).toLocaleString()}</small></a>`).join('')
      : `<p class="muted">${t('noNotifications')}</p>`}
    ${NF_ADMIN ? '' : `<p class="muted" style="font-size:12px">${t('tgHint')}</p>`}`;
}
function toggleBell() {
  NB.open = !NB.open;
  $('#bellBtn')?.setAttribute('aria-expanded', NB.open);
  nfSound.unlock(); // klik di lonceng = gestur user → suara bisa dibunyikan setelah ini
  if (NB.data) drawBell(); else $('#bellPanel')?.classList.toggle('hidden', !NB.open);
  loadBell();
}
async function readOne(id) {
  NB.open = false;
  try { await api(NF_BASE + '/read', {method: 'POST', body: {ids: [id]}}); } catch {}
  loadBell();
}
async function readAll() {
  try { await api(NF_BASE + '/read', {method: 'POST', body: {all: true}}); } catch {}
  loadBell();
}
document.addEventListener('click', e => {
  if (NB.open && !e.target.closest('.bellwrap')) { NB.open = false; $('#bellPanel')?.classList.add('hidden'); }
});
// Cadangan kalau aliran realtime tidak tersambung (proxy memotong, dll.).
setInterval(() => { if (document.visibilityState === 'visible' && !nfLive.connected && $('#bellBtn')) loadBell(); }, 60000);

// ---- Bunyi (disintesis WebAudio, tanpa file) ----
const nfSound = {
  ctx: null,
  get() {
    if (!this.ctx) { const C = window.AudioContext || window.webkitAudioContext; if (!C) return null; this.ctx = new C(); }
    return this.ctx;
  },
  ready() { return !!this.ctx && this.ctx.state === 'running'; },
  // unlock: dipanggil dari gestur user (klik/tombol). demo = bunyikan "ding" sebagai konfirmasi.
  unlock(demo) {
    const c = this.get();
    if (!c) return;
    const p = c.state === 'running' ? Promise.resolve() : c.resume();
    p.then(() => { $('#nfUnlock')?.remove(); if (demo) this.play('ding', 1, NB.prefs?.volume ?? 80); if (NB.open) drawBell(); }).catch(() => {});
  },
  tone(c, t0, f, dur, vol, type = 'sine') {
    const o = c.createOscillator(), g = c.createGain();
    o.type = type; o.frequency.setValueAtTime(f, t0);
    g.gain.setValueAtTime(0.0001, t0);
    g.gain.exponentialRampToValueAtTime(vol, t0 + 0.012);
    g.gain.exponentialRampToValueAtTime(0.0001, t0 + dur);
    o.connect(g).connect(c.destination);
    o.start(t0); o.stop(t0 + dur + 0.02);
  },
  // tiap bunyi mengembalikan durasinya (detik)
  kinds: {
    order(s, c, t, v) { s.tone(c, t, 988, .14, v); s.tone(c, t + .13, 1319, .14, v); s.tone(c, t + .26, 1760, .32, v); s.tone(c, t + .26, 880, .32, v * .5, 'triangle'); return .6; },
    alert(s, c, t, v) { s.tone(c, t, 740, .2, v, 'triangle'); s.tone(c, t + .24, 554, .3, v, 'triangle'); return .55; },
    chat(s, c, t, v) { s.tone(c, t, 1047, .09, v * .8); s.tone(c, t + .1, 1397, .16, v * .8); return .28; },
    ding(s, c, t, v) { s.tone(c, t, 1319, .5, v * .8); s.tone(c, t, 2637, .25, v * .2); return .5; },
  },
  play(name, repeat = 1, volume = 80) {
    if (!name || name === 'none' || !this.kinds[name]) return false;
    const c = this.get();
    if (!c) return false;
    if (c.state !== 'running') { c.resume().catch(() => {}); nfUnlockBanner(); if (c.state !== 'running') return false; }
    const v = Math.max(0, Math.min(1, volume / 100)) * 0.6;
    if (v <= 0) return false;
    let at = c.currentTime + 0.03;
    for (let i = 0; i < Math.max(1, Math.min(repeat, 5)); i++) at += this.kinds[name](this, c, at, v) + 0.45;
    return true;
  },
};
// Gestur pertama di halaman membuka kunci suara (aturan autoplay browser).
['pointerdown', 'keydown', 'touchend'].forEach(ev => document.addEventListener(ev, () => { if (P.token && !nfSound.ready()) nfSound.unlock(); }, {passive: true}));
function nfUnlockBanner() {
  if ($('#nfUnlock') || !document.body) return;
  const b = document.createElement('button');
  b.type = 'button'; b.id = 'nfUnlock'; b.className = 'nf-unlock';
  b.innerHTML = `🔔 <span><b>${esc(t('nfSoundOn'))}</b><small>${esc(t('nfSoundHint'))}</small></span>`;
  b.onclick = () => nfSound.unlock(true);
  document.body.append(b);
}

// Beberapa tab terbuka → hanya satu yang membunyikan notifikasi yang sama.
async function nfClaim(id) {
  const key = 'phsar-nf-' + NF_SCOPE + id;
  const take = () => {
    try { if (localStorage.getItem(key)) return false; localStorage.setItem(key, Date.now()); setTimeout(() => { try { localStorage.removeItem(key); } catch {} }, 60000); } catch {}
    return true;
  };
  if (navigator.locks) return navigator.locks.request('phsar-nf', () => take());
  return take();
}

// ---- Realtime (Server-Sent Events lewat fetch supaya bisa mengirim Authorization) ----
const nfLive = {
  running: false, connected: false, ctrl: null, lang: null, retry: 1000,
  ensure() {
    if (nfAuthed() && !this.running) this.loop();
    else if (this.running && this.lang && this.lang !== P.lang) this.ctrl?.abort(); // bahasa berubah → sambung ulang
  },
  tokenExp() {
    try { return JSON.parse(atob(P.token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/'))).exp * 1000; } catch { return 0; }
  },
  async loop() {
    this.running = true;
    while (nfAuthed()) {
      if (this.tokenExp() && this.tokenExp() - Date.now() < 30000 && !(await refreshSession())) break;
      this.ctrl = new AbortController(); this.lang = P.lang;
      const started = Date.now();
      try {
        const r = await fetch(NF_BASE + '/stream', {headers: {Authorization: 'Bearer ' + P.token, 'X-Lang': P.lang, Accept: 'text/event-stream'}, cache: 'no-store', signal: this.ctrl.signal});
        if (r.status === 401) { if (await refreshSession()) continue; break; }
        if (!r.ok || !r.body) throw new Error('stream ' + r.status);
        await this.read(r.body);
      } catch {}
      this.connected = false;
      if (!nfAuthed()) break;
      if (Date.now() - started > 20000) this.retry = 1000; // sempat tersambung lama → sambung lagi segera
      await new Promise(res => setTimeout(res, this.retry));
      this.retry = Math.min(this.retry * 2, 30000);
    }
    this.running = false; this.connected = false;
  },
  async read(body) {
    const rd = body.getReader(), dec = new TextDecoder();
    let buf = '';
    for (;;) {
      const {value, done} = await rd.read();
      if (done) return;
      buf += dec.decode(value, {stream: true});
      let i;
      while ((i = buf.indexOf('\n\n')) >= 0) {
        const chunk = buf.slice(0, i); buf = buf.slice(i + 2);
        let ev = 'message', data = '';
        for (const line of chunk.split('\n')) {
          if (line.startsWith('event: ')) ev = line.slice(7);
          else if (line.startsWith('data: ')) data += line.slice(6);
        }
        if (!data) continue;
        let d; try { d = JSON.parse(data); } catch { continue; }
        this.onEvent(ev, d);
      }
    }
  },
  onEvent(ev, d) {
    if (ev === 'hello') { this.connected = true; nfSetUnread(d.unread); if (NB.open || NB.data) loadBell(); return; }
    if (ev === 'read') { nfSetUnread(d.unread); if (NB.open) loadBell(); else NB.data = null; return; }
    if (ev !== 'notification') return;
    nfSetUnread(d.unread);
    if (NB.data) { NB.data.notifications.unshift(d); NB.data.notifications = NB.data.notifications.slice(0, 30); NB.data.unread = d.unread; }
    if (NB.open) drawBell();
    nfAlert(d);
  },
};

// Kartu di halaman + bunyi (+ notifikasi sistem kalau tab tersembunyi dan push belum aktif di browser ini).
async function nfAlert(d) {
  const mine = await nfClaim(d.id);
  if (mine) nfSound.play(d.sound, d.repeat || 1, d.volume ?? 80);
  if (document.visibilityState === 'visible') {
    const urgent = d.sound === 'order' || d.sound === 'alert';
    const card = notify(d.body || d.title, {type: urgent ? 'warning' : 'info', title: d.body ? d.title : '', timeout: d.sound === 'order' ? 20000 : 8000});
    if (d.link && card) {
      card.classList.add('nf-click');
      card.querySelector('.notice-body').onclick = () => { readOne(d.id); location.href = d.link; card.querySelector('.notice-x').click(); };
    }
  } else if (mine && !NB.pushOn && window.Notification && Notification.permission === 'granted') {
    try {
      const n = new Notification(d.title, {body: d.body || '', tag: 'phsar-' + d.id, icon: '/ui/icon-192.png', silent: true, requireInteraction: d.sound === 'order'});
      n.onclick = () => { window.focus(); if (d.link) location.href = d.link; n.close(); };
    } catch {}
  }
}

// ---- Notifikasi browser (Web Push) ----
const nfPush = {
  supported() { return 'serviceWorker' in navigator && 'PushManager' in window && 'Notification' in window && window.isSecureContext; },
  async sub() {
    if (!this.supported()) return null;
    const reg = await navigator.serviceWorker.getRegistration(NF_SCOPE);
    return reg ? reg.pushManager.getSubscription() : null;
  },
  key(b64) {
    const s = atob(b64.replace(/-/g, '+').replace(/_/g, '/') + '='.repeat((4 - b64.length % 4) % 4));
    return Uint8Array.from(s, c => c.charCodeAt(0));
  },
  async enable(vapid) {
    if (!this.supported()) throw new Error(t('nfPushNo'));
    if (!vapid) throw new Error(t('ntError'));
    const perm = await Notification.requestPermission();
    if (perm !== 'granted') throw new Error(t('nfPushDenied'));
    const reg = await navigator.serviceWorker.register(NF_SCOPE + 'sw.js', {scope: NF_SCOPE});
    await navigator.serviceWorker.ready;
    let sub = await reg.pushManager.getSubscription();
    const want = this.key(vapid);
    if (sub && sub.options && sub.options.applicationServerKey) {
      const have = new Uint8Array(sub.options.applicationServerKey);
      if (have.length !== want.length || have.some((b, i) => b !== want[i])) { await sub.unsubscribe(); sub = null; } // kunci server berganti
    }
    if (!sub) sub = await reg.pushManager.subscribe({userVisibleOnly: true, applicationServerKey: want});
    await this.save(sub);
    NB.pushOn = true;
  },
  async save(sub) {
    const j = sub.toJSON();
    await api(NF_BASE + '/push', {method: 'POST', body: {endpoint: j.endpoint, keys: j.keys}});
  },
  async disable() {
    const sub = await this.sub();
    if (sub) {
      try { await api(NF_BASE + '/push', {method: 'DELETE', body: {endpoint: sub.endpoint}}); } catch {}
      try { await sub.unsubscribe(); } catch {}
    }
    NB.pushOn = false;
  },
  // Saat halaman dibuka: langganan yang sudah ada didaftarkan ulang ke user yang sedang login.
  async sync() {
    try {
      const sub = await this.sub();
      NB.pushOn = !!sub && Notification.permission === 'granted';
      if (NB.pushOn && nfAuthed()) await this.save(sub);
    } catch {}
  },
};
// Keluar → hentikan aliran & lepas push browser ini (komputer bersama tidak menerima notifikasi user sebelumnya).
const nfLogout = logout;
logout = async function () { // eslint-disable-line no-global-assign
  try { if (P.token) await nfPush.disable(); } catch {}
  nfLive.ctrl?.abort();
  NB.data = null; nfSetUnread(0);
  return nfLogout();
};
(function nfBoot() {
  let synced = false;
  const tick = () => {
    if (nfAuthed() && !synced) { synced = true; nfPush.sync(); nfLive.ensure(); }
    if (!nfAuthed()) { synced = false; if (nfLive.running) nfLive.ctrl?.abort(); }
  };
  setInterval(tick, 3000);
  setTimeout(tick, 500);
})();

// ---- Dialog pengaturan ----
async function nfOpenSettings() {
  NB.open = false; $('#bellPanel')?.classList.add('hidden');
  nfSound.unlock();
  let data;
  try { data = await api(NF_BASE + '/prefs'); } catch (e) { return toast(errMsg(e), 'error'); }
  NB.prefs = data.prefs;
  await nfPush.sync();
  let d = $('#nfDlg');
  if (!d) { d = document.createElement('dialog'); d.id = 'nfDlg'; d.className = 'dlg nfdlg'; document.body.append(d); }
  const p = data.prefs, quiet = !!(p.quiet_from && p.quiet_to);
  const pushState = !nfPush.supported() ? 'no' : Notification.permission === 'denied' ? 'denied' : NB.pushOn ? 'on' : 'off';
  const cb = (g, k, dis) => `<td><input type="checkbox" data-g="${g}" data-k="${k}" ${p.groups[g]?.[k] ? 'checked' : ''} ${dis ? 'disabled' : ''} aria-label="${esc(t('ng_' + g))} — ${esc(t({sound: 'nfSound', push: 'nfPushCol', telegram: 'nfTelegram'}[k]))}"></td>`;
  d.innerHTML = `<form method="dialog" onsubmit="return false">
    <h2 tabindex="-1" autofocus>${t('nfSettings')}</h2>
    <section class="nf-sec"><b>${t('nfThisBrowser')}</b>
      <div class="nf-row"><div><span>${t('nfPush')}</span><small class="muted">${t({on: 'nfPushOn', off: 'nfPushHint', denied: 'nfPushDenied', no: 'nfPushNo'}[pushState])}</small></div>
        ${pushState === 'on' ? `<button type="button" onclick="nfTogglePush(false)">${t('nfDisable')}</button>` : pushState === 'off' ? `<button type="button" class="primary" onclick="nfTogglePush(true)">${t('nfEnable')}</button>` : ''}</div>
      <div class="nf-row"><label class="nf-vol">${t('nfVolume')} <input type="range" min="0" max="100" step="5" id="nfVol" value="${p.volume}"></label>
        <button type="button" onclick="nfSound.unlock(); nfSound.play(NF_ADMIN ? 'alert' : 'order', NF_ADMIN ? 1 : 3, +$('#nfVol').value)">▶ ${t('nfTestSound')}</button></div>
      ${NF_ADMIN ? '' : `<small class="muted">🛎️ ${t('nfOrderRepeat')}</small>`}
    </section>
    <div class="nf-tablewrap"><table class="nf-table"><thead><tr><th>${t('nfGroup')}</th><th>${t('nfSound')}</th><th>${t('nfPushCol')}</th>${data.telegram_linked ? `<th>${t('nfTelegram')}</th>` : ''}</tr></thead>
      <tbody>${data.groups.map(g => `<tr><th scope="row">${esc(t('ng_' + g))}</th>${cb(g, 'sound')}${cb(g, 'push')}${data.telegram_linked ? cb(g, 'telegram') : ''}</tr>`).join('')}</tbody></table></div>
    ${data.telegram_linked ? '' : `<small class="muted">${t('nfTgNotLinked')}</small>`}
    <section class="nf-sec"><label class="nf-check"><input type="checkbox" id="nfQuiet" ${quiet ? 'checked' : ''} onchange="$('#nfQuietTimes').classList.toggle('hidden', !this.checked)"> ${t('nfQuiet')}</label>
      <div id="nfQuietTimes" class="nf-row ${quiet ? '' : 'hidden'}"><label>${t('nfQuietFrom')} <input type="time" id="nfQFrom" value="${esc(p.quiet_from || '22:00')}"></label>
        <label>${t('nfQuietTo')} <input type="time" id="nfQTo" value="${esc(p.quiet_to || '07:00')}"></label></div>
      <small class="muted">${t('nfQuietHint')}</small></section>
    <div class="nf-foot"><button type="button" class="linkbtn" onclick="nfSendTest(this)">${t('nfTestSend')}</button><span class="grow"></span>
      <button type="button" onclick="$('#nfDlg').close()">${t('ntClose')}</button><button type="button" class="primary" onclick="nfSaveSettings(this)">${t('nfSave')}</button></div>
  </form>`;
  if (!d.open) d.showModal();
}
async function nfSaveSettings(btn) {
  // kolom yang tidak ditampilkan (mis. Telegram belum terhubung) tetap memakai nilai lama
  const groups = JSON.parse(JSON.stringify(NB.prefs?.groups || {}));
  document.querySelectorAll('#nfDlg input[data-g]').forEach(i => { (groups[i.dataset.g] = groups[i.dataset.g] || {})[i.dataset.k] = i.checked; });
  const quiet = $('#nfQuiet').checked;
  btn.disabled = true;
  try {
    const r = await api(NF_BASE + '/prefs', {method: 'PUT', body: {groups, volume: +$('#nfVol').value,
      quiet_from: quiet ? $('#nfQFrom').value : '', quiet_to: quiet ? $('#nfQTo').value : ''}});
    NB.prefs = r.prefs;
    $('#nfDlg').close(); toast(t('nfSaved'));
  } catch (e) { toast(errMsg(e), 'error'); } finally { btn.disabled = false; }
}
async function nfTogglePush(on) {
  try {
    if (on) { const r = await api(NF_BASE + '/prefs'); await nfPush.enable(r.vapid_public_key); } else await nfPush.disable();
  } catch (e) { toast(e.message || t('ntError'), 'error'); }
  nfOpenSettings();
}
async function nfSendTest(btn) {
  btn.disabled = true;
  try { await api(NF_BASE + '/test', {method: 'POST'}); toast(t('nfTestSent'), 'info'); } catch (e) { toast(errMsg(e), 'error'); }
  setTimeout(() => { btn.disabled = false; }, 10000);
}
