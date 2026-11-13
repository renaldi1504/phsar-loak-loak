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
