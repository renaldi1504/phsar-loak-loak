'use strict';
// Langkah kedua login admin (aplikasi authenticator) + halaman kelola admin.
addStrings({
  en: {mfaSetupTitle: 'Protect your admin account', mfaSetupSub: 'Every Phsar admin signs in with a code from an authenticator app (Google Authenticator, Microsoft Authenticator, 1Password…).',
    mfaStep1: '1. Install an authenticator app on your phone.', mfaStep2: '2. In the app, tap “+” and scan this QR code.', mfaCantScan: 'Can’t scan? Enter this key manually:',
    mfaStep3: '3. Type the 6-digit code shown in the app.', mfaTurnOn: 'Turn on and continue', mfaCode: '6-digit code',
    mfaVerifyTitle: 'Enter your authenticator code', mfaVerifySub: 'Open your authenticator app and type the 6-digit code for “Phsar Admin”.',
    mfaUseBackup: 'Lost your phone? Use a backup code', mfaUseApp: 'Use the authenticator app', mfaBackupCode: 'Backup code (e.g. abcd-efgh)', mfaContinue: 'Continue', mfaBackupHint: 'Type one of the backup codes you saved when you set up MFA. Each code works once.',
    mfaBackupTitle: 'Save your backup codes', mfaBackupSub: 'Each code works once. Use one if you lose your phone. Keep them in a password manager or print them — they will not be shown again.',
    mfaCopy: 'Copy codes', mfaCopied: 'Copied', mfaSaved: 'I saved my backup codes', mfaBackupLeft: 'Backup code used. {n} left — create new ones in “Admins” if you are running low.',
    err_mfa_invalid: 'Wrong code, or this code was already used. Wait for the next code.', err_mfa_locked: 'Too many wrong codes. Try again in 15 minutes.',
    err_mfa_unavailable: 'MFA is not configured on the server (DATA_ENCRYPTION_KEY).', err_session_invalid: 'Your session ended. Please sign in again.',
    err_mfa_required: 'Enter your authenticator code first.', err_user_not_found: 'This number or Telegram account has never signed in to Phsar. Ask them to sign in once at /admin first.',
    err_self: 'You cannot change your own account here.',
    navAdmins: 'Admins', admins: 'Admins', adminsSub: 'Everyone needs an authenticator app. Reset MFA when someone loses their phone.',
    addAdmin: 'Add admin', role: 'Role', mfa: 'MFA', lastLogin: 'Last sign-in', on: 'On', off: 'Not set up', active: 'Active', inactive: 'Inactive', you: 'you',
    resetMfa: 'Reset MFA', confirmReset: 'Tap again to confirm', deactivate: 'Deactivate', activate: 'Activate', save: 'Save', saved: 'Saved', backupLeft: '{n} backup codes left',
    newBackup: 'New backup codes', role_super_admin: 'Super admin', role_moderator: 'Moderator', role_cs: 'Customer service', role_finance: 'Finance', role_ads_reviewer: 'Ads reviewer', never: 'never'},
  km: {mfaSetupTitle: 'ការពារគណនីអ្នកគ្រប់គ្រង', mfaSetupSub: 'អ្នកគ្រប់គ្រង Phsar ទាំងអស់ចូលដោយប្រើលេខកូដពីកម្មវិធី authenticator (Google Authenticator, Microsoft Authenticator…)។',
    mfaStep1: '1. ដំឡើងកម្មវិធី authenticator លើទូរស័ព្ទ។', mfaStep2: '2. ក្នុងកម្មវិធី ចុច “+” ហើយស្កេន QR នេះ។', mfaCantScan: 'ស្កេនមិនបាន? បញ្ចូលកូដនេះដោយដៃ៖',
    mfaStep3: '3. វាយលេខកូដ 6 ខ្ទង់ដែលបង្ហាញក្នុងកម្មវិធី។', mfaTurnOn: 'បើក ហើយបន្ត', mfaCode: 'លេខកូដ 6 ខ្ទង់',
    mfaVerifyTitle: 'បញ្ចូលលេខកូដ authenticator', mfaVerifySub: 'បើកកម្មវិធី authenticator ហើយវាយលេខកូដ 6 ខ្ទង់សម្រាប់ “Phsar Admin”។',
    mfaUseBackup: 'បាត់ទូរស័ព្ទ? ប្រើលេខកូដបម្រុង', mfaUseApp: 'ប្រើកម្មវិធី authenticator', mfaBackupCode: 'លេខកូដបម្រុង (ឧ. abcd-efgh)', mfaContinue: 'បន្ត', mfaBackupHint: 'វាយលេខកូដបម្រុងមួយដែលអ្នកបានរក្សាទុកពេលរៀបចំ MFA។ លេខកូដនីមួយៗប្រើបានតែម្តង។',
    mfaBackupTitle: 'រក្សាទុកលេខកូដបម្រុង', mfaBackupSub: 'លេខកូដនីមួយៗប្រើបានតែម្តង។ ប្រើពេលបាត់ទូរស័ព្ទ។ រក្សាទុកក្នុង password manager ឬបោះពុម្ព — វានឹងមិនបង្ហាញម្តងទៀតទេ។',
    mfaCopy: 'ចម្លងលេខកូដ', mfaCopied: 'បានចម្លង', mfaSaved: 'ខ្ញុំបានរក្សាទុកលេខកូដបម្រុងហើយ', mfaBackupLeft: 'បានប្រើលេខកូដបម្រុង។ នៅសល់ {n}។',
    err_mfa_invalid: 'លេខកូដខុស ឬបានប្រើរួចហើយ។ សូមរង់ចាំលេខកូដបន្ទាប់។', err_mfa_locked: 'លេខកូដខុសច្រើនពេក។ សូមព្យាយាមម្តងទៀតក្នុង 15 នាទី។',
    err_mfa_unavailable: 'ម៉ាស៊ីនមេមិនទាន់កំណត់ MFA (DATA_ENCRYPTION_KEY)។', err_session_invalid: 'វគ្គបានបញ្ចប់។ សូមចូលម្តងទៀត។',
    err_mfa_required: 'សូមបញ្ចូលលេខកូដ authenticator ជាមុនសិន។', err_user_not_found: 'លេខ ឬគណនី Telegram នេះមិនដែលចូល Phsar ទេ។ សូមឱ្យគាត់ចូល /admin ម្តងសិន។',
    err_self: 'មិនអាចកែគណនីខ្លួនឯងនៅទីនេះទេ។',
    navAdmins: 'អ្នកគ្រប់គ្រង', admins: 'អ្នកគ្រប់គ្រង', adminsSub: 'គ្រប់គ្នាត្រូវមានកម្មវិធី authenticator។ កំណត់ MFA ឡើងវិញពេលនរណាម្នាក់បាត់ទូរស័ព្ទ។',
    addAdmin: 'បន្ថែមអ្នកគ្រប់គ្រង', role: 'តួនាទី', mfa: 'MFA', lastLogin: 'ចូលចុងក្រោយ', on: 'បើក', off: 'មិនទាន់រៀបចំ', active: 'សកម្ម', inactive: 'អសកម្ម', you: 'អ្នក',
    resetMfa: 'កំណត់ MFA ឡើងវិញ', confirmReset: 'ចុចម្តងទៀតដើម្បីបញ្ជាក់', deactivate: 'បិទ', activate: 'បើក', save: 'រក្សាទុក', saved: 'បានរក្សាទុក', backupLeft: 'លេខកូដបម្រុងនៅសល់ {n}',
    newBackup: 'លេខកូដបម្រុងថ្មី', role_super_admin: 'អ្នកគ្រប់គ្រងកំពូល', role_moderator: 'អ្នកត្រួតពិនិត្យ', role_cs: 'សេវាអតិថិជន', role_finance: 'ហិរញ្ញវត្ថុ', role_ads_reviewer: 'ពិនិត្យពាណិជ្ជកម្ម', never: 'មិនដែល'},
  zh: {mfaSetupTitle: '保护你的管理员账号', mfaSetupSub: '所有 Phsar 管理员登录时都需要验证器应用的验证码（Google Authenticator、Microsoft Authenticator、1Password…）。',
    mfaStep1: '1. 在手机上安装验证器应用。', mfaStep2: '2. 在应用中点击“+”并扫描此二维码。', mfaCantScan: '无法扫描？请手动输入此密钥：',
    mfaStep3: '3. 输入应用中显示的 6 位验证码。', mfaTurnOn: '开启并继续', mfaCode: '6 位验证码',
    mfaVerifyTitle: '输入验证器验证码', mfaVerifySub: '打开验证器应用，输入“Phsar Admin”的 6 位验证码。',
    mfaUseBackup: '手机丢了？使用备用码', mfaUseApp: '使用验证器应用', mfaBackupCode: '备用码（例如 abcd-efgh）', mfaContinue: '继续', mfaBackupHint: '输入设置 MFA 时保存的任一备用码，每个只能使用一次。',
    mfaBackupTitle: '保存备用码', mfaBackupSub: '每个备用码只能用一次，手机丢失时使用。请保存在密码管理器或打印出来——之后不会再显示。',
    mfaCopy: '复制备用码', mfaCopied: '已复制', mfaSaved: '我已保存备用码', mfaBackupLeft: '已使用备用码，还剩 {n} 个。',
    err_mfa_invalid: '验证码错误或已使用，请等待下一个验证码。', err_mfa_locked: '错误次数过多，请 15 分钟后再试。',
    err_mfa_unavailable: '服务器未配置 MFA（DATA_ENCRYPTION_KEY）。', err_session_invalid: '会话已结束，请重新登录。',
    err_mfa_required: '请先输入验证器验证码。', err_user_not_found: '该号码或 Telegram 账号从未登录过 Phsar，请先让对方在 /admin 登录一次。',
    err_self: '不能在这里修改自己的账号。',
    navAdmins: '管理员', admins: '管理员', adminsSub: '每个人都需要验证器应用。有人丢失手机时请重置 MFA。',
    addAdmin: '添加管理员', role: '角色', mfa: 'MFA', lastLogin: '最近登录', on: '已开启', off: '未设置', active: '启用', inactive: '停用', you: '你',
    resetMfa: '重置 MFA', confirmReset: '再点一次确认', deactivate: '停用', activate: '启用', save: '保存', saved: '已保存', backupLeft: '剩余 {n} 个备用码',
    newBackup: '生成新备用码', role_super_admin: '超级管理员', role_moderator: '审核员', role_cs: '客服', role_finance: '财务', role_ads_reviewer: '广告审核', never: '从未'},
});

const MFA = {backup: false, setup: null, codes: null};
const ROLES = ['super_admin', 'moderator', 'cs', 'finance', 'ads_reviewer'];

function mfaFail(e) {
  if (e.code === 'session_invalid' || e.status === 401 && e.code !== 'mfa_invalid') { P.token = null; P.mfa = null; toast(t('err_session_invalid'), 'error'); render(); return; }
  toast(errMsg(e), 'error');
}

// mfaPost: endpoint /admin/auth/mfa/* memakai access token (bukan cookie).
async function mfaPost(path, body) {
  const r = await fetch('/admin/auth/mfa' + path, {method: 'POST', headers: {'Content-Type': 'application/json', 'Authorization': 'Bearer ' + P.token, 'X-Lang': P.lang},
    body: JSON.stringify(body || {})});
  const data = await r.json().catch(() => ({}));
  if (!r.ok) throw Object.assign(new Error(data.error?.message || r.statusText), {code: data.error?.code, status: r.status});
  return data;
}

function qrSVG(text) {
  try { const q = qrcode(0, 'M'); q.addData(text); q.make(); return q.createSvgTag({cellSize: 5, margin: 3, scalable: true}); }
  catch { return ''; }
}

function codeInput(id, ph, backup) {
  return `<input id="${id}" class="code-in${backup ? ' backup' : ''}" ${backup ? 'autocomplete="off" autocapitalize="none" spellcheck="false" maxlength="12"' :
    'inputmode="numeric" autocomplete="one-time-code" maxlength="6" pattern="[0-9]*"'} placeholder="${esc(ph)}" aria-label="${esc(ph)}">`;
}

async function renderMFA() {
  if (MFA.codes) return renderBackupCodes();
  if (!P.mfa.enrolled) {
    if (!MFA.setup) {
      try { MFA.setup = await mfaPost('/setup'); } catch (e) { $('#app').innerHTML = shell(`<div class="banner bad">${esc(errMsg(e))}</div>`); return; }
    }
    const key = MFA.setup.secret.replace(/(.{4})/g, '$1 ').trim();
    $('#app').innerHTML = shell(`<section class="card login mfa">
      <div><h1>${t('mfaSetupTitle')}</h1><p class="muted">${t('mfaSetupSub')}</p></div>
      <p>${t('mfaStep1')}</p><p>${t('mfaStep2')}</p>
      <div class="qr">${qrSVG(MFA.setup.otpauth_url)}</div>
      <p class="muted small">${t('mfaCantScan')}<br><code class="secret">${esc(key)}</code></p>
      <p>${t('mfaStep3')}</p>
      ${codeInput('mfaCode', t('mfaCode'))}
      <button class="primary" id="mfaGo" type="button">${t('mfaTurnOn')}</button>
    </section>`);
    const go = async () => {
      const b = $('#mfaGo'); b.disabled = true;
      try {
        const r = await mfaPost('/enable', {code: $('#mfaCode').value});
        useSession(r.session); MFA.setup = null; MFA.codes = r.backup_codes; render();
      } catch (e) { b.disabled = false; $('#mfaCode').value = ''; $('#mfaCode').focus(); mfaFail(e); }
    };
    bindCode(go);
    return;
  }
  $('#app').innerHTML = shell(`<section class="card login mfa">
    <div><h1>${t('mfaVerifyTitle')}</h1><p class="muted">${MFA.backup ? t('mfaBackupHint') : t('mfaVerifySub')}</p></div>
    ${codeInput('mfaCode', MFA.backup ? t('mfaBackupCode') : t('mfaCode'), MFA.backup)}
    <button class="primary" id="mfaGo" type="button">${t('mfaContinue')}</button>
    <button type="button" class="linkish" id="mfaSwitch">${MFA.backup ? t('mfaUseApp') : t('mfaUseBackup')}</button>
  </section>`);
  $('#mfaSwitch').onclick = () => { MFA.backup = !MFA.backup; render(); };
  const go = async () => {
    const b = $('#mfaGo'); b.disabled = true;
    try {
      const r = await mfaPost('/verify', {code: $('#mfaCode').value});
      useSession(r.session); MFA.backup = false;
      if (r.backup_codes_left !== undefined) toast(t('mfaBackupLeft', {n: r.backup_codes_left}), 'info');
      location.hash = location.hash || '#/'; render();
    } catch (e) { b.disabled = false; $('#mfaCode').value = ''; $('#mfaCode').focus(); mfaFail(e); }
  };
  bindCode(go);
}

function bindCode(go) {
  const inp = $('#mfaCode');
  inp.focus();
  inp.addEventListener('keydown', e => { if (e.key === 'Enter') go(); });
  inp.addEventListener('input', () => { if (!inp.classList.contains('backup') && /^\d{6}$/.test(inp.value)) go(); });
  $('#mfaGo').onclick = go;
}

function renderBackupCodes() {
  $('#app').innerHTML = shell(`<section class="card login mfa">
    <div><h1>${t('mfaBackupTitle')}</h1><p class="muted">${t('mfaBackupSub')}</p></div>
    <ol class="codes">${MFA.codes.map(c => `<li><code>${esc(c)}</code></li>`).join('')}</ol>
    <button type="button" id="copyCodes">${t('mfaCopy')}</button>
    <label class="check"><input type="checkbox" id="savedCodes"> ${t('mfaSaved')}</label>
    <button class="primary" id="codesDone" type="button" disabled>${t('mfaContinue')}</button>
  </section>`);
  $('#copyCodes').onclick = async () => {
    const text = 'Phsar Admin — backup codes\n' + MFA.codes.join('\n');
    try { await navigator.clipboard.writeText(text); toast(t('mfaCopied')); }
    catch { const r = document.createRange(); r.selectNodeContents($('.codes')); getSelection().removeAllRanges(); getSelection().addRange(r); }
  };
  $('#savedCodes').onchange = e => { $('#codesDone').disabled = !e.target.checked; };
  $('#codesDone').onclick = () => { MFA.codes = null; location.hash = location.hash || '#/'; render(); };
}

// ---- Halaman kelola admin (super_admin) ----
const ADM = {confirm: null};
async function renderAdmins() {
  const {admins} = await api('/admin/api/admins');
  const roleSel = (id, cur, dis) => `<select data-role="${id}" ${dis ? 'disabled' : ''}>${ROLES.map(r => `<option value="${r}" ${r === cur ? 'selected' : ''}>${t('role_' + r)}</option>`).join('')}</select>`;
  const when = d => d ? new Date(d).toLocaleString() : `<span class="muted">${t('never')}</span>`;
  const rows = admins.map(a => `<tr class="${a.active ? '' : 'muted'}">
    <td><b>${esc(a.name || '—')}</b>${a.self ? ` <span class="chip pink">${t('you')}</span>` : ''}<br><span class="muted">${esc(a.phone)}</span></td>
    <td>${roleSel(a.user_id, a.role, a.self)}</td>
    <td>${a.mfa_enabled ? `<span class="chip ok">${t('on')}</span><br><span class="muted small">${t('backupLeft', {n: a.backup_codes_left})}</span>` : `<span class="chip warn">${t('off')}</span>`}</td>
    <td>${when(a.last_login_at)}</td>
    <td>${a.active ? `<span class="chip ok">${t('active')}</span>` : `<span class="chip">${t('inactive')}</span>`}</td>
    <td class="row" style="justify-content:flex-end">${a.self ? `<button type="button" id="newBackup">${t('newBackup')}</button>` : `
      ${a.mfa_enabled ? `<button type="button" data-reset="${a.user_id}">${ADM.confirm === a.user_id ? t('confirmReset') : t('resetMfa')}</button>` : ''}
      <button type="button" class="${a.active ? 'danger' : ''}" data-active="${a.user_id}" data-to="${!a.active}">${a.active ? t('deactivate') : t('activate')}</button>`}</td></tr>`).join('');
  $('#app').innerHTML = shell(`${pageHead(t('i_admins'), t('adminsSub'))}
    <div class="tabs"><button type="button" aria-selected="true">${t('tabUsers')}</button><button type="button" aria-selected="false" onclick="renderRoles()">${t('tabRoles')}</button></div>
    <section class="card"><h2>${t('addAdmin')}</h2><div class="row addrow">
      <input id="newPhone" placeholder="012 345 678 / @telegram" aria-label="${t('phone')}" style="flex:1;min-width:180px">
      ${roleSel('new', 'moderator', false).replace('data-role="new"', 'id="newRole" aria-label="' + t('role') + '"')}
      <button class="primary" type="button" id="addAdminBtn">${t('addAdmin')}</button></div></section>
    <section class="card" style="overflow-x:auto"><table><thead><tr><th>${t('admins')}</th><th>${t('role')}</th><th>${t('mfa')}</th><th>${t('lastLogin')}</th><th>${t('status')}</th><th></th></tr></thead>
      <tbody>${rows}</tbody></table>
      <ul class="rules"><li>${t('rule1')}</li><li>${t('rule2')}</li><li>${t('rule3')}</li></ul></section>`);
  const act = async (fn) => { try { await fn(); toast(t('saved')); ADM.confirm = null; renderAdmins(); } catch (e) { toast(errMsg(e), 'error'); } };
  $('#addAdminBtn').onclick = () => act(() => api('/admin/api/admins', {method: 'POST', body: {phone: $('#newPhone').value, role: $('#newRole').value}}));
  document.querySelectorAll('[data-role]').forEach(s => s.onchange = () => act(() => api('/admin/api/admins/' + s.dataset.role, {method: 'PATCH', body: {role: s.value}})));
  document.querySelectorAll('[data-active]').forEach(b => b.onclick = () => act(() => api('/admin/api/admins/' + b.dataset.active, {method: 'PATCH', body: {active: b.dataset.to === 'true'}})));
  document.querySelectorAll('[data-reset]').forEach(b => b.onclick = () => {
    if (ADM.confirm !== b.dataset.reset) { ADM.confirm = b.dataset.reset; b.textContent = t('confirmReset'); return; }
    act(() => api('/admin/api/admins/' + b.dataset.reset + '/mfa-reset', {method: 'POST'}));
  });
  const nb = $('#newBackup');
  if (nb) nb.onclick = async () => {
    try { const r = await mfaPost('/backup-codes'); MFA.codes = r.backup_codes; renderBackupCodes(); } catch (e) { mfaFail(e); }
  };
}

const PERMS = ['seller.verify', 'pickup_point.verify', 'product.moderate', 'translation.review', 'report.handle', 'catalog.manage', 'order.view', 'return.decide',
  'support.handle', 'customer.view', 'delivery.dispatch', 'payout.approve', 'refund.decide', 'reconciliation.run', 'commission.manage', 'audit.view', 'ads.review', 'admin.manage', 'payment.configure', 'integration.configure', 'legal.manage', 'content.manage'];
async function renderRoles() {
  const {roles} = await api('/admin/api/roles');
  const has = (r, p) => (roles[r] || []).includes('*') || (roles[r] || []).includes(p);
  $('#app').innerHTML = shell(`${pageHead(t('i_admins'), t('adminsSub'))}
    <div class="tabs"><button type="button" aria-selected="false" onclick="renderAdmins()">${t('tabUsers')}</button><button type="button" aria-selected="true">${t('tabRoles')}</button></div>
    <section class="panel"><div class="tbl-wrap"><table class="matrix"><thead><tr><th>${t('perm')}</th>${ROLES.map(r => `<th>${t('role_' + r)}</th>`).join('')}</tr></thead>
      <tbody>${PERMS.map(p => `<tr><td><code>${p}</code></td>${ROLES.map(r => `<td>${has(r, p) ? '<span class="yes">✓</span>' : '<span class="no">—</span>'}</td>`).join('')}</tr>`).join('')}</tbody></table></div>
      <ul class="rules"><li>${t('rule1')}</li><li>${t('rule2')}</li><li>${t('rule3')}</li></ul></section>`);
}
