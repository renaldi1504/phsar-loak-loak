'use strict';
// Super admin: payment gateway (mode checkout + kredensial ABA PayWay). Kredensial disimpan terenkripsi di server,
// API key tidak pernah dikirim balik — hanya 4 karakter terakhir. Menyimpan wajib kode authenticator.
addStrings({
  en: {i_gateway: 'Payment gateway', gwDesc: 'Only super admins. Credentials are encrypted on the server and are never shown again — only the last 4 characters.',
    gwStatus: 'Checkout', gwOpenPW: 'Open · ABA PayWay', gwOpenFake: 'Open · test payments', gwClosed: 'Closed (pre-launch)',
    gwSource: 'Settings from', gwSrcAdmin: 'Admin (this page)', gwSrcEnv: 'Server .env', gwLife: 'QR valid for', gwMin: '{n} min',
    gwConn: 'PayWay connection', gwVerified: 'Verified {d}', gwNotTested: 'Not tested yet', gwFailed: 'Last test failed',
    gwMode: 'Checkout mode', gwModeOff: 'Closed — pre-launch', gwModeOffSub: 'Website and Seller Center work; buyers cannot check out yet.',
    gwModePW: 'ABA PayWay (KHQR)', gwModePWSub: 'Buyers pay by scanning KHQR with ABA Mobile or any Cambodian bank app.',
    gwModeFake: 'Test payments', gwModeFakeSub: 'Development only: a "simulate payment" button instead of a real QR.',
    gwLifeLabel: 'QR valid for (minutes, 3–120)', gwEnvNote: 'This server currently follows its .env file (PAYMENT_PROVIDER={m}). Saving here takes over.',
    gwFollowEnv: 'Follow server .env again', gwSaveMode: 'Save mode', gwPW: 'ABA PayWay', gwEnv: 'Environment', gwSandbox: 'Sandbox (testing)',
    gwProduction: 'Production (real money)', gwMerchant: 'Merchant ID', gwKey: 'API key', gwKeyKeep: 'Saved key ends with {h} — leave empty to keep it',
    gwKeyNew: 'Paste the API key from ABA', gwAdvanced: 'Advanced: API address', gwBase: 'API address (empty = default for the environment)',
    gwCallback: 'Callback URL — give this to ABA', gwSavePW: 'Save credentials', gwTest: 'Test connection', gwTesting: 'Testing…',
    gwCode: 'Authenticator code', gwCodeHint: 'Required to save — the 6-digit code from your authenticator app.',
    gwSaved: 'Saved', gwSrcPWEnv: 'Credentials currently come from the server .env file. Saving here moves them into the encrypted admin settings.',
    gwProdWarn: 'Production: real buyer money goes to this merchant account.', gwInactive: 'PayWay is selected but the credentials are incomplete — checkout stays closed.',
    err_mfa_invalid: 'Wrong or already-used authenticator code.', err_mfa_locked: 'Too many wrong codes. Try again in 15 minutes.'},
  km: {i_gateway: 'ច្រកទូទាត់', gwDesc: 'សម្រាប់អ្នកគ្រប់គ្រងកំពូលប៉ុណ្ណោះ។ ព័ត៌មានសម្ងាត់ត្រូវបានអ៊ិនគ្រីបនៅលើម៉ាស៊ីនមេ ហើយមិនបង្ហាញម្តងទៀតទេ — មានតែ 4 តួចុងក្រោយ។',
    gwStatus: 'ការទូទាត់', gwOpenPW: 'បើក · ABA PayWay', gwOpenFake: 'បើក · ការទូទាត់សាកល្បង', gwClosed: 'បិទ (មុនបើកដំណើរការ)',
    gwSource: 'ការកំណត់ពី', gwSrcAdmin: 'អ្នកគ្រប់គ្រង (ទំព័រនេះ)', gwSrcEnv: 'ឯកសារ .env ម៉ាស៊ីនមេ', gwLife: 'QR មានសុពលភាព', gwMin: '{n} នាទី',
    gwConn: 'ការតភ្ជាប់ PayWay', gwVerified: 'បានផ្ទៀងផ្ទាត់ {d}', gwNotTested: 'មិនទាន់សាកល្បង', gwFailed: 'ការសាកល្បងចុងក្រោយបរាជ័យ',
    gwMode: 'របៀបទូទាត់', gwModeOff: 'បិទ — មុនបើកដំណើរការ', gwModeOffSub: 'គេហទំព័រ និង Seller Center ដំណើរការ ប៉ុន្តែអ្នកទិញមិនទាន់អាចទូទាត់បានទេ។',
    gwModePW: 'ABA PayWay (KHQR)', gwModePWSub: 'អ្នកទិញបង់ដោយស្កេន KHQR ជាមួយ ABA Mobile ឬកម្មវិធីធនាគារណាមួយ។',
    gwModeFake: 'ការទូទាត់សាកល្បង', gwModeFakeSub: 'សម្រាប់ការអភិវឌ្ឍប៉ុណ្ណោះ៖ ប៊ូតុង «ក្លែងការបង់ប្រាក់» ជំនួស QR ពិត។',
    gwLifeLabel: 'QR មានសុពលភាព (នាទី, 3–120)', gwEnvNote: 'ម៉ាស៊ីនមេនេះកំពុងប្រើឯកសារ .env (PAYMENT_PROVIDER={m})។ ការរក្សាទុកនៅទីនេះនឹងជំនួស។',
    gwFollowEnv: 'ត្រឡប់ទៅប្រើ .env វិញ', gwSaveMode: 'រក្សាទុករបៀប', gwPW: 'ABA PayWay', gwEnv: 'បរិស្ថាន', gwSandbox: 'Sandbox (សាកល្បង)',
    gwProduction: 'Production (ប្រាក់ពិត)', gwMerchant: 'Merchant ID', gwKey: 'API key', gwKeyKeep: 'កូនសោដែលបានរក្សាទុកបញ្ចប់ដោយ {h} — ទុកទទេដើម្បីរក្សាទុកដដែល',
    gwKeyNew: 'បិទភ្ជាប់ API key ពី ABA', gwAdvanced: 'កម្រិតខ្ពស់៖ អាសយដ្ឋាន API', gwBase: 'អាសយដ្ឋាន API (ទទេ = លំនាំដើមតាមបរិស្ថាន)',
    gwCallback: 'Callback URL — ផ្តល់ឱ្យ ABA', gwSavePW: 'រក្សាទុកព័ត៌មាន', gwTest: 'សាកល្បងការតភ្ជាប់', gwTesting: 'កំពុងសាកល្បង…',
    gwCode: 'លេខកូដ Authenticator', gwCodeHint: 'ត្រូវការដើម្បីរក្សាទុក — លេខកូដ 6 ខ្ទង់ពីកម្មវិធី authenticator។',
    gwSaved: 'បានរក្សាទុក', gwSrcPWEnv: 'ព័ត៌មានបច្ចុប្បន្នមកពីឯកសារ .env។ ការរក្សាទុកនៅទីនេះនឹងផ្លាស់ទីទៅការកំណត់ដែលបានអ៊ិនគ្រីប។',
    gwProdWarn: 'Production៖ ប្រាក់ពិតរបស់អ្នកទិញចូលគណនី merchant នេះ។', gwInactive: 'បានជ្រើស PayWay ប៉ុន្តែព័ត៌មានមិនទាន់គ្រប់ — ការទូទាត់នៅតែបិទ។',
    err_mfa_invalid: 'លេខកូដ authenticator ខុស ឬបានប្រើរួច។', err_mfa_locked: 'លេខកូដខុសច្រើនដងពេក។ សូមព្យាយាមម្តងទៀតក្នុង 15 នាទី។'},
  zh: {i_gateway: '支付网关', gwDesc: '仅限超级管理员。凭据在服务器上加密保存，保存后不再显示——只显示最后 4 位。',
    gwStatus: '结账', gwOpenPW: '开放 · ABA PayWay', gwOpenFake: '开放 · 测试支付', gwClosed: '关闭（预上线）',
    gwSource: '设置来源', gwSrcAdmin: '管理后台（本页）', gwSrcEnv: '服务器 .env', gwLife: '二维码有效期', gwMin: '{n} 分钟',
    gwConn: 'PayWay 连接', gwVerified: '已验证 {d}', gwNotTested: '尚未测试', gwFailed: '上次测试失败',
    gwMode: '结账模式', gwModeOff: '关闭——预上线', gwModeOffSub: '网站和卖家中心正常使用，买家暂时不能结账。',
    gwModePW: 'ABA PayWay（KHQR）', gwModePWSub: '买家用 ABA Mobile 或任何柬埔寨银行 App 扫描 KHQR 付款。',
    gwModeFake: '测试支付', gwModeFakeSub: '仅限开发：显示“模拟付款”按钮，而不是真实二维码。',
    gwLifeLabel: '二维码有效期（分钟，3–120）', gwEnvNote: '此服务器目前使用 .env 文件（PAYMENT_PROVIDER={m}）。在此保存后将以本页为准。',
    gwFollowEnv: '恢复使用服务器 .env', gwSaveMode: '保存模式', gwPW: 'ABA PayWay', gwEnv: '环境', gwSandbox: 'Sandbox（测试）',
    gwProduction: 'Production（真实资金）', gwMerchant: 'Merchant ID', gwKey: 'API key', gwKeyKeep: '已保存的密钥以 {h} 结尾——留空则保持不变',
    gwKeyNew: '粘贴 ABA 提供的 API key', gwAdvanced: '高级：API 地址', gwBase: 'API 地址（留空 = 按环境默认）',
    gwCallback: '回调 URL——提供给 ABA', gwSavePW: '保存凭据', gwTest: '测试连接', gwTesting: '正在测试…',
    gwCode: '验证器验证码', gwCodeHint: '保存时需要——验证器 App 中的 6 位验证码。',
    gwSaved: '已保存', gwSrcPWEnv: '凭据目前来自服务器 .env 文件。在此保存后将转入加密的后台设置。',
    gwProdWarn: 'Production：买家的真实资金将进入此商户账户。', gwInactive: '已选择 PayWay，但凭据不完整——结账仍然关闭。',
    err_mfa_invalid: '验证码错误或已被使用。', err_mfa_locked: '错误次数过多，请 15 分钟后再试。'},
});

const GW = {data: null, test: null};

async function renderGateway() {
  GW.data = await api('/admin/api/payments');
  drawGateway();
}

function drawGateway() {
  const d = GW.data, pw = d.payway;
  const status = d.active === 'payway' ? t('gwOpenPW') + ' · ' + (pw.environment === 'production' ? 'Production' : 'Sandbox')
    : d.active === 'fake' ? t('gwOpenFake') : t('gwClosed');
  const last = GW.test || pw.last_test;
  const conn = pw.verified_at ? t('gwVerified', {d: new Date(pw.verified_at).toLocaleDateString()}) : last && !last.ok ? t('gwFailed') : t('gwNotTested');
  const modes = [['off', 'gwModeOff', 'gwModeOffSub'], ['payway', 'gwModePW', 'gwModePWSub']].concat(d.production ? [] : [['fake', 'gwModeFake', 'gwModeFakeSub']]);
  const radio = ([v, l, s]) => `<label class="opt"><input type="radio" name="gwMode" value="${v}" ${d.mode === v ? 'checked' : ''}>
    <span><b>${t(l)}</b><small class="muted">${t(s)}</small></span></label>`;
  const codeField = id => `<label>${t('gwCode')}<input id="${id}" name="mfa_code" inputmode="numeric" autocomplete="one-time-code" maxlength="6" placeholder="123456" style="max-width:180px"></label>`;
  $('#app').innerHTML = shell(`${pageHead(t('i_integrations'), t('intDesc'))}${intTabs('gateway')}<h2 style="font-size:18px">${t('i_gateway')}</h2><p class="muted">${t('gwDesc')}</p>
    <div class="kpis">
      <div class="kpi"><span>${t('gwStatus')}</span><b style="font-size:20px">${esc(status)}</b></div>
      <div class="kpi"><span>${t('gwSource')}</span><b style="font-size:20px">${t(d.source === 'admin' ? 'gwSrcAdmin' : 'gwSrcEnv')}</b></div>
      <div class="kpi"><span>${t('gwLife')}</span><b>${t('gwMin', {n: d.lifetime_min})}</b></div>
      <div class="kpi"><span>${t('gwConn')}</span><b style="font-size:20px" class="${last && !last.ok && !pw.verified_at ? 'warn' : ''}">${esc(conn)}</b></div>
    </div>
    ${d.mode === 'payway' && d.active !== 'payway' ? `<div class="banner bad">${t('gwInactive')}</div>` : ''}
    <section class="panel" id="gwModeForm"><div class="ph"><h2>${t('gwMode')}</h2></div>
      <div style="padding:16px 18px;display:flex;flex-direction:column;gap:14px">
        ${d.source === 'env' ? `<div class="banner info">${t('gwEnvNote', {m: esc(d.env_mode)})}</div>` : ''}
        <div class="opts" role="radiogroup" aria-label="${t('gwMode')}">${modes.map(radio).join('')}</div>
        <label>${t('gwLifeLabel')}<input id="gwLife" name="lifetime_min" type="number" min="3" max="120" value="${d.lifetime_min}" style="max-width:180px"></label>
        <p class="muted small">${t('gwCodeHint')}</p>
        <div class="row" style="display:flex;flex-wrap:wrap;gap:10px;align-items:flex-end">${codeField('gwCode1')}
          <button class="primary" id="gwSaveMode" type="button">${t('gwSaveMode')}</button>
          ${d.source === 'admin' ? `<button type="button" id="gwFollowEnv">${t('gwFollowEnv')}</button>` : ''}</div>
      </div></section>
    <section class="panel" id="gwPWForm"><div class="ph"><h2>${t('gwPW')}</h2>${pw.key_set ? `<span class="chip ${pw.verified_at ? 'ok' : 'warn'}">${esc(conn)}</span>` : ''}</div>
      <div style="padding:16px 18px;display:flex;flex-direction:column;gap:14px">
        ${pw.source === 'env' ? `<div class="banner info">${t('gwSrcPWEnv')}</div>` : ''}
        <div style="display:grid;grid-template-columns:repeat(auto-fit,minmax(240px,1fr));gap:14px">
          <label>${t('gwEnv')}<select id="gwEnv" name="environment">
            <option value="sandbox" ${pw.environment !== 'production' ? 'selected' : ''}>${t('gwSandbox')}</option>
            <option value="production" ${pw.environment === 'production' ? 'selected' : ''}>${t('gwProduction')}</option></select></label>
          <label>${t('gwMerchant')}<input id="gwMerchant" name="merchant_id" autocomplete="off" spellcheck="false" maxlength="30" value="${esc(pw.merchant_id || '')}"></label>
        </div>
        <div class="banner warn" id="gwProdWarn" ${pw.environment === 'production' ? '' : 'hidden'}>${t('gwProdWarn')}</div>
        <label>${t('gwKey')}<input id="gwKey" name="api_key" type="password" autocomplete="new-password" spellcheck="false"
          placeholder="${esc(pw.key_set ? t('gwKeyKeep', {h: '…' + pw.key_hint}) : t('gwKeyNew'))}"></label>
        <details><summary>${t('gwAdvanced')}</summary>
          <label style="margin-top:10px">${t('gwBase')}<input id="gwBase" name="base_url" spellcheck="false" placeholder="https://checkout-sandbox.payway.com.kh"
            value="${esc(pw.source === 'admin' && ![`https://checkout-sandbox.payway.com.kh`, `https://checkout.payway.com.kh`].includes(pw.base_url) ? pw.base_url : '')}"></label></details>
        <div><span class="muted small">${t('gwCallback')}</span><div style="display:flex;gap:8px;align-items:center;flex-wrap:wrap;margin-top:4px">
          <code style="overflow-wrap:anywhere">${esc(d.callback_url)}</code><button type="button" id="gwCopy">${t('copy')}</button></div></div>
        ${last ? `<div class="banner ${last.ok ? 'ok' : 'bad'}">${esc(last.message)} <span class="muted small">· ${new Date(last.at).toLocaleString()}</span></div>` : ''}
        <p class="muted small">${t('gwCodeHint')}</p>
        <div style="display:flex;flex-wrap:wrap;gap:10px;align-items:flex-end">${codeField('gwCode2')}
          <button class="primary" id="gwSavePW" type="button">${t('gwSavePW')}</button>
          <button type="button" id="gwTest" ${pw.key_set ? '' : 'disabled'}>${t('gwTest')}</button></div>
      </div></section>`);

  const fail = (e, form) => {
    if (e.field) { const el = $(`#${form} [name="${e.field}"]`); if (el) { el.classList.add('invalid'); el.focus(); } }
    toast(errMsg(e), 'error');
  };
  const clear = form => document.querySelectorAll(`#${form} .invalid`).forEach(el => el.classList.remove('invalid'));
  $('#gwEnv').onchange = () => { $('#gwProdWarn').hidden = $('#gwEnv').value !== 'production'; };
  $('#gwCopy').onclick = async () => { try { await navigator.clipboard.writeText(d.callback_url); toast(t('copied')); } catch {} };
  $('#gwSaveMode').onclick = async () => {
    clear('gwModeForm');
    const mode = (document.querySelector('input[name="gwMode"]:checked') || {}).value;
    try {
      GW.data = await api('/admin/api/payments/settings', {method: 'PUT', body: {mode, lifetime_min: +$('#gwLife').value, mfa_code: $('#gwCode1').value}});
      toast(t('gwSaved')); drawGateway();
    } catch (e) { fail(e, 'gwModeForm'); }
  };
  if ($('#gwFollowEnv')) $('#gwFollowEnv').onclick = async () => {
    clear('gwModeForm');
    try {
      GW.data = await api('/admin/api/payments/settings', {method: 'PUT', body: {mode: '', lifetime_min: 0, mfa_code: $('#gwCode1').value}});
      toast(t('gwSaved')); drawGateway();
    } catch (e) { fail(e, 'gwModeForm'); }
  };
  $('#gwSavePW').onclick = async () => {
    clear('gwPWForm');
    const body = {environment: $('#gwEnv').value, merchant_id: $('#gwMerchant').value.trim(), base_url: $('#gwBase').value.trim(), mfa_code: $('#gwCode2').value};
    if ($('#gwKey').value.trim()) body.api_key = $('#gwKey').value.trim();
    try {
      GW.data = await api('/admin/api/payments/payway', {method: 'PUT', body});
      GW.test = null; toast(t('gwSaved')); drawGateway();
    } catch (e) { fail(e, 'gwPWForm'); }
  };
  $('#gwTest').onclick = async () => {
    const b = $('#gwTest'); b.disabled = true; b.textContent = t('gwTesting');
    try { GW.test = await api('/admin/api/payments/payway/test', {method: 'POST'}); GW.data = await api('/admin/api/payments'); drawGateway(); }
    catch (e) { toast(errMsg(e), 'error'); b.disabled = false; b.textContent = t('gwTest'); }
  };
}
