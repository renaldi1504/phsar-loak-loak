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
// toast(pesan, jenis): notifikasi singkat. jenis = success (default) · error · warning · info. Lihat notify() di bawah.
function toast(msg, type = 'success', opts = {}) { return notify(msg, Object.assign({type}, opts)); }
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
function useSession(s) { P.token = s.access_token; P.user = s.user; P.mfa = s.mfa || null; }
async function refreshSession() {
  try { useSession(await authPost('/refresh')); return true; } catch { P.token = null; return false; }
}
async function logout() { try { await authPost('/logout'); } catch {} P.token = null; P.user = null; P.mfa = null; location.hash = ''; render(); }

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
  toast(bad.length === 1 ? `${bad[0].label}: ${bad[0].msg}` : t('fixN', {n: bad.length}), 'error', {title: t('fixTitle')});
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
    toast(`${fieldName(el)}: ${detail || t('checkField')}`, 'error', {title: t('fixTitle')});
  } else {
    showErrorSummary([{msg}], anchor);
    toast(msg, 'error');
  }
}

// errMsg: pesan error yang dibaca manusia. Kode dikenal → teks terjemahan; error per kolom → "Label kolom: pesan".
function errMsg(e) {
  if (!e) return t('ntError');
  if (t('err_' + e.code) !== 'err_' + e.code) return t('err_' + e.code);
  if (e.field) {
    let el = null;
    try { el = document.querySelector(`[name="${CSS.escape(e.field)}"]`); } catch { /* nama aneh */ }
    return `${(el && fieldName(el)) || e.field.replace(/_/g, ' ')}: ${e.message}`;
  }
  return e.message || t('ntError');
}

// ---- Login (Telegram / SMS) ----
const loginT = {
  en: {loginTitle: 'Sign in', loginSub: 'Use the same account as the Phsar app.', tg: 'Continue with Telegram', or: 'or',
    phone: 'Phone number', sendCode: 'Send SMS code', code: '6-digit code', verify: 'Sign in', waitingTg: 'Open Telegram and tap Start in the bot…',
    err_not_admin: 'This account is not a Phsar admin.', tgExpired: 'The Telegram link expired. Try again.', codeSent: 'Code sent. (Development: the code is in the server log.)', loginAgain: 'Please sign in again.',
    err_otp_invalid: 'Wrong or expired code.', err_otp_too_soon: 'Wait a minute before requesting another code.', err_phone_invalid: 'Invalid phone number.',
    err_telegram_disabled: 'Telegram login is not configured on this server.', err_otp_rate_limited: 'Too many codes requested. Try again later.',
    fb: 'Continue with Facebook', err_facebook_cancelled: 'Facebook login was cancelled.', err_facebook_failed: 'Facebook login failed. Please try again.', err_facebook_expired: 'The Facebook login took too long. Please try again.', err_facebook_disabled: 'Facebook login is not available yet.', err_sms_failed: 'The SMS could not be sent right now. Try again or use Telegram.', err_sms_disabled: 'SMS login is not available yet. Please continue with Telegram.', err_user_suspended: 'This account is suspended.'},
  km: {loginTitle: 'ចូលគណនី', loginSub: 'ប្រើគណនីដូចកម្មវិធី Phsar។', tg: 'បន្តជាមួយ Telegram', or: 'ឬ',
    phone: 'លេខទូរស័ព្ទ', sendCode: 'ផ្ញើលេខកូដ SMS', code: 'លេខកូដ 6 ខ្ទង់', verify: 'ចូល', waitingTg: 'បើក Telegram ហើយចុច Start ក្នុងបូត…',
    err_not_admin: 'គណនីនេះមិនមែនជាអ្នកគ្រប់គ្រង Phsar ទេ។', tgExpired: 'តំណ Telegram ផុតកំណត់។ សូមព្យាយាមម្តងទៀត។', codeSent: 'បានផ្ញើលេខកូដ។', loginAgain: 'សូមចូលម្តងទៀត។',
    err_otp_invalid: 'លេខកូដខុស ឬផុតកំណត់។', err_otp_too_soon: 'សូមរង់ចាំមួយនាទីមុនស្នើលេខកូដថ្មី។', err_phone_invalid: 'លេខទូរស័ព្ទមិនត្រឹមត្រូវ។',
    err_telegram_disabled: 'មិនទាន់កំណត់ការចូលតាម Telegram។', err_otp_rate_limited: 'ស្នើលេខកូដច្រើនពេក។ សូមព្យាយាមពេលក្រោយ។',
    fb: 'បន្តជាមួយ Facebook', err_facebook_cancelled: 'បានបោះបង់ការចូលតាម Facebook។', err_facebook_failed: 'ការចូលតាម Facebook បរាជ័យ។ សូមព្យាយាមម្តងទៀត។', err_facebook_expired: 'ការចូលតាម Facebook យូរពេក។ សូមព្យាយាមម្តងទៀត។', err_facebook_disabled: 'មិនទាន់អាចចូលតាម Facebook បានទេ។', err_sms_failed: 'មិនអាចផ្ញើ SMS ឥឡូវនេះទេ។ សូមព្យាយាមម្តងទៀត ឬប្រើ Telegram។', err_sms_disabled: 'មិនទាន់អាចចូលដោយ SMS បានទេ។ សូមបន្តជាមួយ Telegram។', err_user_suspended: 'គណនីនេះត្រូវបានផ្អាក។'},
  zh: {loginTitle: '登录', loginSub: '使用与 Phsar App 相同的账号。', tg: '使用 Telegram 继续', or: '或',
    phone: '手机号', sendCode: '发送短信验证码', code: '6 位验证码', verify: '登录', waitingTg: '请打开 Telegram 并在机器人中点击 Start…',
    err_not_admin: '此账号不是 Phsar 管理员。', tgExpired: 'Telegram 链接已过期，请重试。', codeSent: '验证码已发送。', loginAgain: '请重新登录。',
    err_otp_invalid: '验证码错误或已过期。', err_otp_too_soon: '请稍等一分钟再获取验证码。', err_phone_invalid: '手机号无效。',
    err_telegram_disabled: '服务器未配置 Telegram 登录。', err_otp_rate_limited: '验证码请求过多，请稍后再试。',
    fb: '使用 Facebook 继续', err_facebook_cancelled: '已取消 Facebook 登录。', err_facebook_failed: 'Facebook 登录失败，请重试。', err_facebook_expired: 'Facebook 登录超时，请重试。', err_facebook_disabled: '暂不支持 Facebook 登录。', err_sms_failed: '短信暂时无法发送，请重试或使用 Telegram。', err_sms_disabled: '暂不支持短信登录，请使用 Telegram 继续。', err_user_suspended: '此账号已被停用。'},
};
function addStrings(dict) { for (const l of ['en', 'km', 'zh']) P.T[l] = Object.assign(P.T[l] || {}, dict[l]); }
addStrings(loginT);

// Server memberi tahu lewat /ui/config.js apakah login SMS / Telegram tersedia (production awal: hanya Telegram).
function loginModes() {
  const c = window.PHSAR_CONFIG || {};
  return {tg: c.telegram !== false, sms: c.sms !== false, fb: c.facebook === true};
}
const fbIcon = '<svg viewBox="0 0 24 24" width="20" height="20" aria-hidden="true"><path fill="#fff" d="M24 12a12 12 0 1 0-13.9 11.9v-8.4H7.1V12h3V9.4c0-3 1.8-4.7 4.5-4.7 1.3 0 2.7.2 2.7.2v3h-1.5c-1.5 0-2 .9-2 1.9V12h3.4l-.5 3.5h-2.9v8.4A12 12 0 0 0 24 12z"/></svg>';
// Kembali dari Facebook dengan ?auth_error=… → tampilkan pesan, lalu hapus dari alamat.
(function () {
  const q = new URLSearchParams(location.search), e = q.get('auth_error');
  if (!e) return;
  q.delete('auth_error');
  history.replaceState(null, '', location.pathname + (q.toString() ? '?' + q : '') + location.hash);
  window.addEventListener('load', () => setTimeout(() => toast(t('err_' + e) !== 'err_' + e ? t('err_' + e) : e, 'error'), 300));
})();
function loginView(title) {
  const m = loginModes();
  return `<section class="card login">
    <div><h1>${esc(title)}</h1><p class="muted">${t('loginSub')}</p></div>
    ${m.tg ? `<button class="tg" id="tgBtn" type="button">${t('tg')}</button>
    <p class="muted hidden" id="tgWait">${t('waitingTg')}</p>` : ''}
    ${m.fb ? `<button class="fb" id="fbBtn" type="button">${fbIcon}${t('fb')}</button>` : ''}
    ${(m.tg || m.fb) && m.sms ? `<div class="or">${t('or')}</div>` : ''}
    ${m.sms ? `<label>${t('phone')}<input id="phone" inputmode="tel" autocomplete="tel" placeholder="012 345 678"></label>
    <button id="sendBtn" type="button">${t('sendCode')}</button>
    <label>${t('code')}<input id="code" inputmode="numeric" autocomplete="one-time-code" maxlength="6"></label>
    <button class="primary" id="verifyBtn" type="button">${t('verify')}</button>` : ''}
  </section>`;
}
function bindLogin(onDone) {
  if ($('#fbBtn')) $('#fbBtn').onclick = () => { location.href = P.authBase + '/facebook/start?return=' + encodeURIComponent(location.pathname + location.search + location.hash); };
  if ($('#tgBtn')) $('#tgBtn').onclick = async () => {
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
          if (p.status === 'ok') { useSession(p.session); onDone(); }
          else { toast(t(p.status === 'not_admin' ? 'err_not_admin' : 'tgExpired'), 'error'); $('#tgWait').classList.add('hidden'); }
        } catch (e) { clearInterval(bindLogin.poll); toast(errMsg(e), 'error'); }
      }, 1500);
    } catch (e) { toast(errMsg(e), 'error'); }
  };
  if (!$('#sendBtn')) return;
  $('#sendBtn').onclick = async () => {
    try { await authPost('/otp/request', {phone: $('#phone').value}); toast(t('codeSent')); $('#code').focus(); } catch (e) { toast(errMsg(e), 'error'); }
  };
  $('#verifyBtn').onclick = async () => {
    try { useSession(await authPost('/otp/verify', {phone: $('#phone').value, code: $('#code').value})); onDone(); } catch (e) { toast(errMsg(e), 'error'); }
  };
}

// ---- Daftar bertingkat: kotak cari + halaman (Seller Center & admin) --------------------------
addStrings({
  en: {pgPrev: 'Previous', pgNext: 'Next', pgNav: 'Pages', pgPage: 'Page {p}', pgTotal: '{n} results', pgRange: '{a}–{b} of {n}',
    lsClear: 'Clear search', lsNoMatch: 'Nothing matches “{q}”.', lsResults: 'Results for “{q}”', dateFrom: 'From', dateTo: 'To', dateClear: 'All dates'},
  km: {pgPrev: 'មុន', pgNext: 'បន្ទាប់', pgNav: 'ទំព័រ', pgPage: 'ទំព័រ {p}', pgTotal: 'លទ្ធផល {n}', pgRange: '{a}–{b} នៃ {n}',
    lsClear: 'សម្អាតការស្វែងរក', lsNoMatch: 'គ្មានអ្វីត្រូវនឹង “{q}” ទេ។', lsResults: 'លទ្ធផលសម្រាប់ “{q}”', dateFrom: 'ពី', dateTo: 'ដល់', dateClear: 'គ្រប់កាលបរិច្ឆេទ'},
  zh: {pgPrev: '上一页', pgNext: '下一页', pgNav: '分页', pgPage: '第 {p} 页', pgTotal: '{n} 条结果', pgRange: '第 {a}–{b} 条，共 {n} 条',
    lsClear: '清除搜索', lsNoMatch: '没有与“{q}”匹配的结果。', lsResults: '“{q}”的搜索结果', dateFrom: '从', dateTo: '至', dateClear: '全部日期'},
});

const LS = {focus: null, timers: {}};

// searchBox: kolom cari. Ketik → onChange(nilai) setelah 350 ms (Enter langsung).
function searchBox(id, value, placeholder) {
  return `<div class="lsearch"><span class="lsearch-ico" aria-hidden="true"></span>
    <input type="search" id="${id}" value="${esc(value || '')}" placeholder="${esc(placeholder)}" aria-label="${esc(placeholder)}" autocomplete="off" enterkeyhint="search" maxlength="80">
    ${value ? `<button type="button" class="lsearch-x" data-ls-clear="${id}" aria-label="${t('lsClear')}">×</button>` : ''}</div>`;
}
function bindSearch(id, onChange) {
  const el = document.getElementById(id);
  if (!el) return;
  const fire = () => { clearTimeout(LS.timers[id]); LS.focus = document.activeElement === el ? id : null; onChange(el.value.trim()); };
  el.addEventListener('input', () => { clearTimeout(LS.timers[id]); LS.timers[id] = setTimeout(fire, 350); });
  el.addEventListener('keydown', e => { if (e.key === 'Enter') { e.preventDefault(); fire(); } else if (e.key === 'Escape' && el.value) { el.value = ''; fire(); } });
  const x = document.querySelector(`[data-ls-clear="${id}"]`);
  if (x) x.onclick = () => { el.value = ''; LS.focus = id; onChange(''); };
  if (LS.focus === id) { // kembalikan fokus & kursor setelah daftar dirender ulang
    LS.focus = null;
    el.focus();
    el.setSelectionRange(el.value.length, el.value.length);
  }
}

// pagerHTML: navigasi halaman dari meta {page, size, total, pages}. Kosong kalau cuma satu halaman.
function pagerHTML(m) {
  if (!m || !m.total) return '';
  const a = (m.page - 1) * m.size + 1, b = Math.min(m.page * m.size, m.total);
  const sum = `<span class="pager-sum">${t('pgRange', {a, b, n: m.total})}</span>`;
  if (m.pages <= 1) return `<nav class="pager single" aria-label="${t('pgNav')}">${sum}</nav>`;
  const nums = [];
  for (let p = 1; p <= m.pages; p++) {
    if (p === 1 || p === m.pages || Math.abs(p - m.page) <= 1) nums.push(p);
    else if (nums[nums.length - 1] !== '…') nums.push('…');
  }
  const btn = (p, label, extra = '') => `<button type="button" data-page="${p}" ${extra}>${label}</button>`;
  return `<nav class="pager" aria-label="${t('pgNav')}">${sum}<div class="pager-btns">
    ${btn(m.page - 1, '‹ ' + t('pgPrev'), `class="pager-step" ${m.page <= 1 ? 'disabled' : ''}`)}
    ${nums.map(p => p === '…' ? '<span class="pager-gap" aria-hidden="true">…</span>'
      : btn(p, p, `aria-label="${t('pgPage', {p})}" ${p === m.page ? 'aria-current="page" class="on"' : ''}`)).join('')}
    ${btn(m.page + 1, t('pgNext') + ' ›', `class="pager-step" ${m.page >= m.pages ? 'disabled' : ''}`)}</div></nav>`;
}
function bindPager(root, onPage) {
  (root || document).querySelectorAll('.pager [data-page]').forEach(b => b.onclick = () => {
    onPage(+b.dataset.page);
    const top = document.getElementById('tab') || document.querySelector('main');
    if (top) top.scrollIntoView({block: 'start', behavior: matchMedia('(prefers-reduced-motion: reduce)').matches ? 'auto' : 'smooth'});
  });
}

// qs: query string dari objek, nilai kosong dilewati.
function qs(o) {
  return Object.entries(o).filter(([, v]) => v !== '' && v != null).map(([k, v]) => encodeURIComponent(k) + '=' + encodeURIComponent(v)).join('&');
}

// ---- Notifikasi: satu template untuk admin, Seller Center, dan kurir ---------------------------
// Kartu di pojok kanan atas (HP: atas, selebar layar). success/info hilang sendiri; error & warning tampil
// lebih lama dan berhenti menghitung saat disentuh/di-hover. Error yang menyangkut kolom yang baru
// ditandai merah juga ditulis di bawah kolom itu.
addStrings({
  en: {ntError: 'Something went wrong', ntErrTitle: 'That didn’t work', ntWarnTitle: 'Check this', ntInfoTitle: 'Note', ntOkTitle: 'Done', ntClose: 'Close', fixTitle: 'Can’t save yet'},
  km: {ntError: 'មានបញ្ហាអ្វីមួយ', ntErrTitle: 'មិនបានសម្រេច', ntWarnTitle: 'សូមពិនិត្យ', ntInfoTitle: 'ចំណាំ', ntOkTitle: 'រួចរាល់', ntClose: 'បិទ', fixTitle: 'មិនទាន់អាចរក្សាទុកបាន'},
  zh: {ntError: '出了点问题', ntErrTitle: '操作未完成', ntWarnTitle: '请检查', ntInfoTitle: '提示', ntOkTitle: '完成', ntClose: '关闭', fixTitle: '暂时无法保存'},
});
const NT_ICON = {
  success: '<path d="M6 12.5l4 4 8-9" fill="none" stroke="currentColor" stroke-width="2.6" stroke-linecap="round" stroke-linejoin="round"/>',
  error: '<path d="M12 6.5v7" stroke="currentColor" stroke-width="2.6" stroke-linecap="round"/><circle cx="12" cy="17.6" r="1.6" fill="currentColor"/>',
  warning: '<path d="M12 6.5v7" stroke="currentColor" stroke-width="2.6" stroke-linecap="round"/><circle cx="12" cy="17.6" r="1.6" fill="currentColor"/>',
  info: '<circle cx="12" cy="6.8" r="1.6" fill="currentColor"/><path d="M12 10.5v7" stroke="currentColor" stroke-width="2.6" stroke-linecap="round"/>',
};
const NT_TIME = {success: 4000, info: 5000, warning: 7000, error: 9000};

function notify(msg, {type = 'success', title = '', timeout} = {}) {
  if (!NT_TIME[type]) type = 'info';
  msg = String(msg ?? '');
  if (type === 'error') ntInlineField(msg);
  let box = document.getElementById('notices');
  if (!box) {
    box = document.createElement('div');
    box.id = 'notices'; box.className = 'notices';
    document.body.append(box);
  }
  // Pesan yang sama berturut-turut: perbarui kartu lama, jangan menumpuk.
  const same = [...box.children].find(n => n.dataset.msg === type + '|' + msg);
  if (same) same.remove();
  const n = document.createElement('div');
  n.className = 'notice notice-' + type;
  n.dataset.msg = type + '|' + msg;
  n.setAttribute('role', type === 'error' ? 'alert' : 'status');
  const head = title || (type === 'success' ? '' : t({error: 'ntErrTitle', warning: 'ntWarnTitle', info: 'ntInfoTitle'}[type]));
  n.innerHTML = `<span class="notice-ico" aria-hidden="true"><svg viewBox="0 0 24 24" width="20" height="20">${NT_ICON[type]}</svg></span>
    <div class="notice-body">${head ? `<b class="notice-title">${esc(head)}</b>` : ''}<span class="notice-msg">${esc(msg)}</span></div>
    <button type="button" class="notice-x" aria-label="${esc(t('ntClose'))}">×</button><i class="notice-bar" aria-hidden="true"></i>`;
  box.prepend(n);
  while (box.children.length > 3) box.lastElementChild.remove();
  const ms = timeout ?? NT_TIME[type];
  const bar = n.querySelector('.notice-bar');
  bar.style.animationDuration = ms + 'ms';
  let left = ms, started = Date.now(), timer = null;
  const close = () => { clearTimeout(timer); n.classList.add('out'); setTimeout(() => n.remove(), 180); };
  const run = () => { started = Date.now(); timer = setTimeout(close, left); n.classList.remove('paused'); };
  const pause = () => { clearTimeout(timer); left -= Date.now() - started; n.classList.add('paused'); };
  n.querySelector('.notice-x').onclick = close;
  n.addEventListener('mouseenter', pause); n.addEventListener('mouseleave', run);
  n.addEventListener('focusin', pause); n.addEventListener('focusout', run);
  run();
  return n;
}

// Kolom yang baru ditandai merah (dan sedang fokus) tanpa pesan di bawahnya → tulis pesannya di situ.
function ntInlineField(msg) {
  const el = document.activeElement;
  if (!el || !el.classList || !el.classList.contains('invalid')) return;
  const box = el.closest('label, .field, td') || el.parentElement;
  if (!box || box.querySelector('.ferr')) return;
  const label = fieldName(el);
  const text = label && msg.startsWith(label + ': ') ? msg.slice(label.length + 2) : msg;
  markField(el, text, {focus: false});
}
