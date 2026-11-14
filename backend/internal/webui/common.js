// Dasar bersama untuk Seller Center dan admin: bahasa, panggilan API, login, util.
// Refresh token ada di cookie HttpOnly; access token hanya di memori.
'use strict';
const P = {
  lang: (() => { try { return localStorage.getItem('phsar-lang') || 'en'; } catch { return 'en'; } })(),
  token: null, user: null, authBase: '', T: {},
};

function t(key, vars) {
  const s = (P.T[P.lang] && P.T[P.lang][key]) || (P.T.en && P.T.en[key]) || key;
  return vars ? s.replace(/\{(\w+)\}/g, (_, k) => vars[k] ?? '') : s;
}
function esc(v) {
  return String(v ?? '').replace(/[&<>"']/g, c => ({'&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;'}[c]));
}
function $(sel, root = document) { return root.querySelector(sel); }
function toast(msg) {
  const el = $('#toast'); el.textContent = msg; el.classList.add('show');
  clearTimeout(toast.t); toast.t = setTimeout(() => el.classList.remove('show'), 3200);
}
function setLang(l) {
  P.lang = l; try { localStorage.setItem('phsar-lang', l); } catch {}
  document.documentElement.lang = l; render();
}
function langSwitch() {
  return `<div class="lang" role="group" aria-label="Language">${[['en', 'EN'], ['km', 'ខ្មែរ'], ['zh', '中文']].map(([c, l]) =>
    `<button type="button" aria-pressed="${P.lang === c}" onclick="setLang('${c}')">${l}</button>`).join('')}</div>`;
}
const brandSVG = `<svg width="28" height="32" viewBox="0 0 24 24" aria-hidden="true"><defs><linearGradient id="bm" x1="0" y1="0" x2="1" y2="1"><stop offset="0" stop-color="#ff7ba8"/><stop offset=".56" stop-color="#f50750"/><stop offset="1" stop-color="#c30963"/></linearGradient></defs><path d="M8 8V5a4 4 0 0 1 8 0v3" fill="none" stroke="#f50750" stroke-width="2.4"/><path d="M3 7h18l1 16H2z" fill="url(#bm)"/><path d="m9 10.5 8 4.5-8 4.5z" fill="#fff"/></svg>`;

async function authPost(path, body) {
  const r = await fetch(P.authBase + path, {method: 'POST', credentials: 'same-origin',
    headers: {'Content-Type': 'application/json', 'X-Phsar': '1', 'X-Lang': P.lang}, body: JSON.stringify(body || {})});
  const data = r.status === 204 ? {} : await r.json().catch(() => ({}));
  if (!r.ok) throw Object.assign(new Error(data.error?.message || r.statusText), {code: data.error?.code, data});
  return data;
}
function useSession(s) { P.token = s.access_token; P.user = s.user; }
async function refreshSession() {
  try { useSession(await authPost('/refresh')); return true; } catch { P.token = null; return false; }
}
async function logout() { try { await authPost('/logout'); } catch {} P.token = null; P.user = null; location.hash = ''; render(); }

// api: panggil JSON API dengan access token; perbarui token sekali kalau kedaluwarsa.
async function api(path, opts = {}, retried) {
  const headers = {'X-Lang': P.lang};
  if (P.token) headers.Authorization = 'Bearer ' + P.token;
  let body = opts.body;
  if (body && !(body instanceof FormData)) { headers['Content-Type'] = 'application/json'; body = JSON.stringify(body); }
  const r = await fetch(path, {method: opts.method || 'GET', headers, body});
  if (r.status === 401 && !retried && await refreshSession()) return api(path, opts, true);
  if (r.status === 401) { P.token = null; render(); throw new Error(t('loginAgain')); }
  const data = r.status === 204 ? {} : await r.json().catch(() => ({}));
  if (!r.ok) throw Object.assign(new Error(data.error?.message || r.statusText), {code: data.error?.code, field: data.error?.field, data});
  return data;
}
async function upload(kind, file) {
  const fd = new FormData(); fd.append('kind', kind); fd.append('file', file);
  return api('/v1/seller/uploads', {method: 'POST', body: fd});
}
// ---- Validasi form: tandai field, scroll, fokus ----
const fieldT = {
  en: {pickOne: 'Please choose one.', fixBelow: 'Can’t save yet. Please fix:', fixN: '{n} things need fixing — see the list above the buttons.', required: 'Please fill this in.', checkField: 'Please check this field.', minChars: 'At least {n} characters.', positive: 'Enter an amount above 0.'},
  km: {pickOne: 'សូមជ្រើសរើសមួយ។', fixBelow: 'មិនទាន់អាចរក្សាទុកបានទេ។ សូមកែ៖', fixN: 'មាន {n} ចំណុចត្រូវកែ — មើលបញ្ជីនៅលើប៊ូតុង។', required: 'សូមបំពេញ។', checkField: 'សូមពិនិត្យវាលនេះ។', minChars: 'យ៉ាងហោចណាស់ {n} តួអក្សរ។', positive: 'សូមបញ្ចូលចំនួនលើសពី 0។'},
  zh: {pickOne: '请选择一项。', fixBelow: '暂时无法保存，请修改：', fixN: '有 {n} 处需要修改——请看按钮上方的列表。', required: '请填写此项。', checkField: '请检查此项。', minChars: '至少 {n} 个字符。', positive: '请输入大于 0 的金额。'},
};
for (const l of ['en', 'km', 'zh']) P.T[l] = Object.assign(P.T[l] || {}, fieldT[l]);

function clearFieldErrors(root = document) {
  root.querySelectorAll('.invalid').forEach(el => { el.classList.remove('invalid'); el.removeAttribute('aria-invalid'); });
  root.querySelectorAll('.ferr, #formErrors').forEach(el => el.remove());
}

// fieldName: nama field yang dibaca manusia — dari data-label, label, judul kolom tabel, atau judul di atasnya.
function fieldName(el) {
  if (!el) return '';
  if (el.dataset && el.dataset.label) return el.dataset.label;
  const td = el.closest('td');
  if (td) {
    const tr = td.closest('tr'), i = [...tr.children].indexOf(td);
    const th = tr.closest('table')?.querySelectorAll('thead th')[i];
    const name = [...tr.children].find(c => !c.querySelector('input,select,button') && !['', '—'].includes(c.textContent.trim()))?.textContent.trim() || '';
    return name ? `${name} · ${th?.textContent.trim() || ''}` : (th?.textContent.trim() || '');
  }
  const lab = el.closest('label, .field');
  if (lab) {
    const own = [...lab.childNodes].filter(n => n.nodeType === 3 || n.classList?.contains('flabel')).map(n => n.textContent).join(' ').trim();
    if (own) return own;
  }
  const sec = el.closest('section, .card');
  return sec?.querySelector('h2')?.textContent.trim() || '';
}

// markField: tandai field (selector/elemen), pesan di bawahnya; focus = scroll + kursor ke situ.
function markField(target, msg, {focus = true} = {}) {
  const el = typeof target === 'string' ? document.querySelector(target) : target;
  if (!el) return false;
  el.classList.add('invalid');
  el.setAttribute('aria-invalid', 'true');
  const box = el.closest('label, .field, td') || el;
  const p = document.createElement('p');
  p.className = 'ferr'; p.textContent = msg || t('checkField');
  if (box === el) el.after(p); else box.append(p);
  const clear = () => {
    el.classList.remove('invalid'); el.removeAttribute('aria-invalid'); p.remove();
    const item = document.querySelector(`#formErrors [data-for="${el.dataset.errId}"]`);
    if (item) { item.remove(); if (!document.querySelector('#formErrors li')) document.querySelector('#formErrors')?.remove(); }
  };
  el.addEventListener('input', clear, {once: true});
  el.addEventListener('change', clear, {once: true});
  if (focus) focusField(el);
  return true;
}
function focusField(el) {
  el.scrollIntoView({block: 'center', behavior: 'smooth'});
  const f = el.matches('input,select,textarea,button') ? el : el.querySelector('input,select,textarea,button');
  if (f) setTimeout(() => f.focus({preventScroll: true}), 300);
}

// Kotak merah yang menetap di dekat tombol simpan: daftar semua yang salah, klik → lompat ke field.
// Letaknya sebelum elemen [data-errors-anchor] (baris tombol). items: [{el, label, msg}] atau [{msg}].
let errSeq = 0;
function showErrorSummary(items, anchor) {
  document.querySelector('#formErrors')?.remove();
  const where = (typeof anchor === 'string' ? document.querySelector(anchor) : anchor) || document.querySelector('[data-errors-anchor]');
  if (!where || !items.length) return;
  const box = document.createElement('div');
  box.id = 'formErrors'; box.className = 'banner bad errsum'; box.setAttribute('role', 'alert');
  box.innerHTML = `<b>${esc(t('fixBelow'))}</b><ul></ul>`;
  for (const it of items) {
    const li = document.createElement('li');
    if (it.el) {
      it.el.dataset.errId = it.el.dataset.errId || 'e' + (++errSeq);
      li.dataset.for = it.el.dataset.errId;
      li.innerHTML = `<a href="#" role="button">${it.label ? `<b>${esc(it.label)}</b>: ` : ''}${esc(it.msg)}</a>`;
      li.querySelector('a').onclick = ev => { ev.preventDefault(); focusField(it.el); };
    } else li.textContent = it.msg;
    box.querySelector('ul').append(li);
  }
  if (where.matches('.sticky-actions')) where.prepend(box); else where.before(box);
}

// checkRequired: [[selector|elemen, pesan?, cek?, label?]] — tandai SEMUA yang salah, lompat ke yang pertama,
// tampilkan ringkasan di dekat tombol. true = semua oke.
function checkRequired(rules, anchor) {
  clearFieldErrors();
  const bad = [];
  for (const [sel, msg, test, label] of rules) {
    const el = typeof sel === 'string' ? document.querySelector(sel) : sel;
    if (!el) continue;
    const ok = test ? test(el) : String(el.value ?? '').trim() !== '';
    if (!ok) bad.push({el, msg: msg || t('required'), label: label || fieldName(el)});
  }
  if (!bad.length) return true;
  bad.forEach((b, i) => markField(b.el, b.msg, {focus: i === 0}));
  showErrorSummary(bad, anchor);
  toast(bad.length === 1 ? `${bad[0].label}: ${bad[0].msg}` : t('fixN', {n: bad.length}));
  return false;
}

// showError: error dari server. fields: {field_server: selector} atau fungsi(field) → selector.
// Field yang dikenal ditandai; apa pun yang terjadi, pesan tetap tampil di ringkasan (bukan hanya toast).
function showError(e, fields, anchor) {
  clearFieldErrors();
  const msg = errMsg(e);
  let el = null;
  if (e && fields) {
    const key = e.field || ({outside_service_area: '_location', upload_not_found: 'media', upload_type: 'media'})[e.code];
    if (key) {
      const sel = typeof fields === 'function' ? fields(key) : (fields[key] || fields[key.replace(/\[\d+\].*$/, '')]);
      el = sel ? (typeof sel === 'string' ? document.querySelector(sel) : sel) : null;
    }
  }
  // Pesan validasi dari server berbahasa Indonesia; tampilkan apa adanya sebagai detail di bawah pesan umum.
  const detail = e?.code === 'validation' && e.message ? e.message : '';
  if (el) {
    markField(el, detail || t('checkField'));
    showErrorSummary([{el, label: fieldName(el), msg: detail || t('checkField')}], anchor);
    toast(`${fieldName(el)}: ${detail || t('checkField')}`);
  } else {
    showErrorSummary([{msg: e?.field ? `${e.field}: ${e.message}` : msg}], anchor);
    toast(msg);
  }
}

function errMsg(e) { return t('err_' + e.code) !== 'err_' + e.code ? t('err_' + e.code) : (e.field ? `${e.field}: ${e.message}` : e.message); }

// ---- Login (Telegram / SMS) ----
const loginT = {
  en: {loginTitle: 'Sign in', loginSub: 'Use the same account as the Phsar app.', tg: 'Continue with Telegram', or: 'or',
    phone: 'Phone number', sendCode: 'Send SMS code', code: '6-digit code', verify: 'Sign in', waitingTg: 'Open Telegram and tap Start in the bot…',
    tgExpired: 'The Telegram link expired. Try again.', codeSent: 'Code sent. (Development: the code is in the server log.)', loginAgain: 'Please sign in again.',
    err_otp_invalid: 'Wrong or expired code.', err_otp_too_soon: 'Wait a minute before requesting another code.', err_phone_invalid: 'Invalid phone number.',
    err_telegram_disabled: 'Telegram login is not configured on this server.', err_otp_rate_limited: 'Too many codes requested. Try again later.'},
  km: {loginTitle: 'ចូលគណនី', loginSub: 'ប្រើគណនីដូចកម្មវិធី Phsar។', tg: 'បន្តជាមួយ Telegram', or: 'ឬ',
    phone: 'លេខទូរស័ព្ទ', sendCode: 'ផ្ញើលេខកូដ SMS', code: 'លេខកូដ 6 ខ្ទង់', verify: 'ចូល', waitingTg: 'បើក Telegram ហើយចុច Start ក្នុងបូត…',
    tgExpired: 'តំណ Telegram ផុតកំណត់។ សូមព្យាយាមម្តងទៀត។', codeSent: 'បានផ្ញើលេខកូដ។', loginAgain: 'សូមចូលម្តងទៀត។',
    err_otp_invalid: 'លេខកូដខុស ឬផុតកំណត់។', err_otp_too_soon: 'សូមរង់ចាំមួយនាទីមុនស្នើលេខកូដថ្មី។', err_phone_invalid: 'លេខទូរស័ព្ទមិនត្រឹមត្រូវ។',
    err_telegram_disabled: 'មិនទាន់កំណត់ការចូលតាម Telegram។', err_otp_rate_limited: 'ស្នើលេខកូដច្រើនពេក។ សូមព្យាយាមពេលក្រោយ។'},
  zh: {loginTitle: '登录', loginSub: '使用与 Phsar App 相同的账号。', tg: '使用 Telegram 继续', or: '或',
    phone: '手机号', sendCode: '发送短信验证码', code: '6 位验证码', verify: '登录', waitingTg: '请打开 Telegram 并在机器人中点击 Start…',
    tgExpired: 'Telegram 链接已过期，请重试。', codeSent: '验证码已发送。', loginAgain: '请重新登录。',
    err_otp_invalid: '验证码错误或已过期。', err_otp_too_soon: '请稍等一分钟再获取验证码。', err_phone_invalid: '手机号无效。',
    err_telegram_disabled: '服务器未配置 Telegram 登录。', err_otp_rate_limited: '验证码请求过多，请稍后再试。'},
};
function addStrings(dict) { for (const l of ['en', 'km', 'zh']) P.T[l] = Object.assign(P.T[l] || {}, dict[l]); }
addStrings(loginT);

function loginView(title) {
  return `<section class="card login">
    <div><h1>${esc(title)}</h1><p class="muted">${t('loginSub')}</p></div>
    <button class="tg" id="tgBtn" type="button">${t('tg')}</button>
    <p class="muted hidden" id="tgWait">${t('waitingTg')}</p>
    <div class="or">${t('or')}</div>
    <label>${t('phone')}<input id="phone" inputmode="tel" autocomplete="tel" placeholder="012 345 678"></label>
    <button id="sendBtn" type="button">${t('sendCode')}</button>
    <label>${t('code')}<input id="code" inputmode="numeric" autocomplete="one-time-code" maxlength="6"></label>
    <button class="primary" id="verifyBtn" type="button">${t('verify')}</button>
  </section>`;
}
function bindLogin(onDone) {
  $('#tgBtn').onclick = async () => {
    try {
      const s = await authPost('/telegram/start', {});
      window.open(s.deep_link, '_blank', 'noopener');
      $('#tgWait').classList.remove('hidden');
      clearInterval(bindLogin.poll);
      bindLogin.poll = setInterval(async () => {
        try {
          const p = await authPost('/telegram/poll', {request_id: s.request_id, poll_token: s.poll_token});
          if (p.status === 'pending') return;
          clearInterval(bindLogin.poll);
          if (p.status === 'ok') { useSession(p.session); onDone(); } else { toast(t('tgExpired')); $('#tgWait').classList.add('hidden'); }
        } catch (e) { clearInterval(bindLogin.poll); toast(errMsg(e)); }
      }, 1500);
    } catch (e) { toast(errMsg(e)); }
  };
  $('#sendBtn').onclick = async () => {
    try { await authPost('/otp/request', {phone: $('#phone').value}); toast(t('codeSent')); $('#code').focus(); } catch (e) { toast(errMsg(e)); }
  };
  $('#verifyBtn').onclick = async () => {
    try { useSession(await authPost('/otp/verify', {phone: $('#phone').value, code: $('#code').value})); onDone(); } catch (e) { toast(errMsg(e)); }
  };
}
