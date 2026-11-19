'use strict';
// Super admin: provider SMS (kode OTP via PlasGate), Login Google, Login Facebook, Sign in with Apple (aplikasi). Rahasia terenkripsi di server, tidak pernah tampil lagi.
addStrings({
  en: {i_logins: 'Login & SMS', lgDesc: 'Only super admins. Keys are encrypted on the server and never shown again — only the last 4 characters.',
    lgSMS: 'SMS login (OTP codes)', lgSMSActive: 'SMS', lgFB: 'Facebook login', lgSrc: 'Settings from', lgOn: 'On', lgOff: 'Off',
    lgModeOff: 'Off — Telegram / Facebook only', lgModeOffSub: 'The phone-number login is hidden on every sign-in page.',
    lgModePG: 'PlasGate (Cambodia)', lgModePGSub: 'Each login code costs one SMS from your PlasGate balance.',
    lgModeLog: 'Log only (development)', lgModeLogSub: 'Codes are written to the server log — never use on the real server.',
    lgFollowEnv: 'Follow server .env again', lgEnvNote: 'SMS currently follows the server .env file (SMS_PROVIDER={m}). Saving here takes over.',
    lgPK: 'Private key', lgSecret: 'Secret', lgSender: 'Sender ID', lgSenderHint: 'The sender name/number registered with PlasGate (e.g. Phsar).',
    lgKeyKeep: 'Saved — ends with {h}. Leave both empty to keep.', lgKeyNew: 'From PlasGate → SMPP & API → SMS Gateway',
    lgSaveSMS: 'Save SMS settings', lgTestPhone: 'Send a test SMS to', lgTestSMS: 'Send test SMS', lgTestCost: 'Uses one SMS from your PlasGate balance.',
    lgFBEnable: 'Show “Continue with Facebook” on sign-in pages', lgAppID: 'App ID', lgAppSecret: 'App Secret',
    lgAppKeep: 'Saved — ends with {h}. Leave empty to keep.', lgAppNew: 'Facebook Developers → your app → Settings → Basic',
    lgRedirects: 'Valid OAuth Redirect URIs — add all of these in Facebook Login → Settings', lgSaveFB: 'Save Facebook settings', lgTestFB: 'Test App ID & Secret',
    lgGG: 'Google login', lgGGEnable: 'Show “Continue with Google” on sign-in pages', lgGGID: 'Client ID (Web application)', lgGGSecret: 'Client Secret',
    lgGGNew: 'Google Cloud → APIs & Services → Credentials → OAuth client ID', lgGGApps: 'App client IDs (Android / iOS, optional)', lgGGAppsHint: 'For the Flutter app later — one per line.',
    lgGGRedirects: 'Authorized redirect URIs — add all of these in the OAuth client', lgGGOrigins: 'Authorized JavaScript origins', lgSaveGG: 'Save Google settings', lgTestGG: 'Test Client ID & Secret',
    lgAP: 'Sign in with Apple — iPhone app (later)', lgAPNote: 'Only for the iPhone app. There is no Apple button on the website. Turn it on when the Flutter app is ready and an Apple Developer account exists.',
    lgAPEnable: 'Allow Sign in with Apple in the iPhone app', lgAPIDs: 'Bundle ID of the iPhone app', lgSaveAP: 'Save Apple settings',
    lgVerified: 'Verified {d}', lgNotTested: 'Not tested yet', lgFailed: 'Last test failed', lgSaved: 'Saved', lgTesting: 'Testing…'},
  km: {i_logins: 'ការចូល និង SMS', lgDesc: 'សម្រាប់អ្នកគ្រប់គ្រងកំពូលប៉ុណ្ណោះ។ កូនសោត្រូវបានអ៊ិនគ្រីប ហើយមិនបង្ហាញម្តងទៀតទេ — មានតែ 4 តួចុងក្រោយ។',
    lgSMS: 'ការចូលតាម SMS (លេខកូដ OTP)', lgSMSActive: 'SMS', lgFB: 'ការចូលតាម Facebook', lgSrc: 'ការកំណត់ពី', lgOn: 'បើក', lgOff: 'បិទ',
    lgModeOff: 'បិទ — Telegram / Facebook ប៉ុណ្ណោះ', lgModeOffSub: 'ការចូលដោយលេខទូរស័ព្ទត្រូវលាក់នៅគ្រប់ទំព័រចូល។',
    lgModePG: 'PlasGate (កម្ពុជា)', lgModePGSub: 'លេខកូដនីមួយៗចំណាយ SMS មួយពីសមតុល្យ PlasGate។',
    lgModeLog: 'កត់ក្នុង log (អភិវឌ្ឍន៍)', lgModeLogSub: 'លេខកូដត្រូវសរសេរក្នុង log ម៉ាស៊ីនមេ — កុំប្រើលើម៉ាស៊ីនមេពិត។',
    lgFollowEnv: 'ត្រឡប់ទៅប្រើ .env វិញ', lgEnvNote: 'SMS កំពុងប្រើឯកសារ .env (SMS_PROVIDER={m})។ ការរក្សាទុកនៅទីនេះនឹងជំនួស។',
    lgPK: 'Private key', lgSecret: 'Secret', lgSender: 'Sender ID', lgSenderHint: 'ឈ្មោះ/លេខអ្នកផ្ញើដែលបានចុះឈ្មោះជាមួយ PlasGate (ឧ. Phsar)។',
    lgKeyKeep: 'បានរក្សាទុក — បញ្ចប់ដោយ {h}។ ទុកទទេទាំងពីរដើម្បីរក្សាដដែល។', lgKeyNew: 'ពី PlasGate → SMPP & API → SMS Gateway',
    lgSaveSMS: 'រក្សាទុកការកំណត់ SMS', lgTestPhone: 'ផ្ញើ SMS សាកល្បងទៅ', lgTestSMS: 'ផ្ញើ SMS សាកល្បង', lgTestCost: 'ចំណាយ SMS មួយពីសមតុល្យ PlasGate។',
    lgFBEnable: 'បង្ហាញ «បន្តជាមួយ Facebook» នៅទំព័រចូល', lgAppID: 'App ID', lgAppSecret: 'App Secret',
    lgAppKeep: 'បានរក្សាទុក — បញ្ចប់ដោយ {h}។ ទុកទទេដើម្បីរក្សាដដែល។', lgAppNew: 'Facebook Developers → កម្មវិធី → Settings → Basic',
    lgRedirects: 'Valid OAuth Redirect URIs — បន្ថែមទាំងអស់នេះក្នុង Facebook Login → Settings', lgSaveFB: 'រក្សាទុកការកំណត់ Facebook', lgTestFB: 'សាកល្បង App ID និង Secret',
    lgGG: 'ការចូលតាម Google', lgGGEnable: 'បង្ហាញ «បន្តជាមួយ Google» នៅទំព័រចូល', lgGGID: 'Client ID (Web application)', lgGGSecret: 'Client Secret',
    lgGGNew: 'Google Cloud → APIs & Services → Credentials → OAuth client ID', lgGGApps: 'Client ID កម្មវិធី (Android / iOS ស្រេចចិត្ត)', lgGGAppsHint: 'សម្រាប់កម្មវិធី Flutter ពេលក្រោយ — មួយបន្ទាត់មួយ។',
    lgGGRedirects: 'Authorized redirect URIs — បន្ថែមទាំងអស់នេះក្នុង OAuth client', lgGGOrigins: 'Authorized JavaScript origins', lgSaveGG: 'រក្សាទុកការកំណត់ Google', lgTestGG: 'សាកល្បង Client ID និង Secret',
    lgAP: 'Sign in with Apple — កម្មវិធី iPhone (ពេលក្រោយ)', lgAPNote: 'សម្រាប់តែកម្មវិធី iPhone ប៉ុណ្ណោះ។ គ្មានប៊ូតុង Apple នៅលើគេហទំព័រទេ។ បើកពេលកម្មវិធី Flutter រួចរាល់ និងមានគណនី Apple Developer។',
    lgAPEnable: 'អនុញ្ញាត Sign in with Apple ក្នុងកម្មវិធី iPhone', lgAPIDs: 'Bundle ID របស់កម្មវិធី iPhone', lgSaveAP: 'រក្សាទុកការកំណត់ Apple',
    lgVerified: 'បានផ្ទៀងផ្ទាត់ {d}', lgNotTested: 'មិនទាន់សាកល្បង', lgFailed: 'ការសាកល្បងចុងក្រោយបរាជ័យ', lgSaved: 'បានរក្សាទុក', lgTesting: 'កំពុងសាកល្បង…'},
  zh: {i_logins: '登录与短信', lgDesc: '仅限超级管理员。密钥在服务器上加密保存，保存后不再显示——只显示最后 4 位。',
    lgSMS: '短信登录（验证码）', lgSMSActive: '短信', lgFB: 'Facebook 登录', lgSrc: '设置来源', lgOn: '开启', lgOff: '关闭',
    lgModeOff: '关闭——仅 Telegram / Facebook', lgModeOffSub: '所有登录页面都隐藏手机号登录。',
    lgModePG: 'PlasGate（柬埔寨）', lgModePGSub: '每个验证码消耗 PlasGate 余额中的一条短信。',
    lgModeLog: '仅写入日志（开发）', lgModeLogSub: '验证码写入服务器日志——切勿在正式服务器使用。',
    lgFollowEnv: '恢复使用服务器 .env', lgEnvNote: '短信目前使用服务器 .env 文件（SMS_PROVIDER={m}）。在此保存后将以本页为准。',
    lgPK: 'Private key', lgSecret: 'Secret', lgSender: 'Sender ID', lgSenderHint: '在 PlasGate 注册的发送方名称/号码（例如 Phsar）。',
    lgKeyKeep: '已保存——以 {h} 结尾。两项都留空则保持不变。', lgKeyNew: '来自 PlasGate → SMPP & API → SMS Gateway',
    lgSaveSMS: '保存短信设置', lgTestPhone: '发送测试短信到', lgTestSMS: '发送测试短信', lgTestCost: '会消耗 PlasGate 余额中的一条短信。',
    lgFBEnable: '在登录页显示“使用 Facebook 继续”', lgAppID: 'App ID', lgAppSecret: 'App Secret',
    lgAppKeep: '已保存——以 {h} 结尾。留空则保持不变。', lgAppNew: 'Facebook Developers → 应用 → Settings → Basic',
    lgRedirects: 'Valid OAuth Redirect URIs——请在 Facebook Login → Settings 中全部添加', lgSaveFB: '保存 Facebook 设置', lgTestFB: '测试 App ID 和 Secret',
    lgGG: 'Google 登录', lgGGEnable: '在登录页显示“使用 Google 继续”', lgGGID: 'Client ID（Web application）', lgGGSecret: 'Client Secret',
    lgGGNew: 'Google Cloud → APIs & Services → Credentials → OAuth client ID', lgGGApps: 'App Client ID（Android / iOS，可选）', lgGGAppsHint: '供以后的 Flutter App 使用——每行一个。',
    lgGGRedirects: 'Authorized redirect URIs——请在 OAuth client 中全部添加', lgGGOrigins: 'Authorized JavaScript origins', lgSaveGG: '保存 Google 设置', lgTestGG: '测试 Client ID 和 Secret',
    lgAP: 'Sign in with Apple——iPhone App（以后）', lgAPNote: '仅用于 iPhone App，网站上没有 Apple 按钮。等 Flutter App 完成并有 Apple Developer 账号后再开启。',
    lgAPEnable: '允许在 iPhone App 中使用 Sign in with Apple', lgAPIDs: 'iPhone App 的 Bundle ID', lgSaveAP: '保存 Apple 设置',
    lgVerified: '已验证 {d}', lgNotTested: '尚未测试', lgFailed: '上次测试失败', lgSaved: '已保存', lgTesting: '正在测试…'},
});

const LG = {data: null};

async function renderLogins() {
  LG.data = await api('/admin/api/integrations');
  drawLogins();
}

function lgConn(x) {
  return x.verified_at ? t('lgVerified', {d: new Date(x.verified_at).toLocaleDateString()}) : x.last_test && !x.last_test.ok ? t('lgFailed') : t('lgNotTested');
}

function drawLogins() {
  const d = LG.data, s = d.sms, f = d.facebook, g = d.google, ap = d.apple;
  const modes = [['off', 'lgModeOff', 'lgModeOffSub'], ['plasgate', 'lgModePG', 'lgModePGSub']].concat(d.production ? [] : [['log', 'lgModeLog', 'lgModeLogSub']]);
  const radio = ([v, l, sub]) => `<label class="opt"><input type="radio" name="lgMode" value="${v}" ${s.mode === v ? 'checked' : ''}><span><b>${t(l)}</b><small class="muted">${t(sub)}</small></span></label>`;
  const code = id => `<label>${t('gwCode')}<input id="${id}" name="mfa_code" inputmode="numeric" autocomplete="one-time-code" maxlength="6" placeholder="123456" style="max-width:180px"></label>`;
  const result = x => x ? `<div class="banner ${x.ok ? 'ok' : 'bad'}">${esc(x.message)} <span class="muted small">· ${new Date(x.at).toLocaleString()}</span></div>` : '';
  const smsLabel = {off: t('lgOff'), plasgate: 'PlasGate', log: 'Log'}[s.active] || s.active;
  $('#app').innerHTML = shell(`${pageHead(t('i_integrations'), t('intDesc'))}${intTabs('logins')}
    <div class="kpis">
      <div class="kpi"><span>${t('lgSMSActive')}</span><b style="font-size:20px">${esc(smsLabel)}</b></div>
      <div class="kpi"><span>${t('lgSrc')}</span><b style="font-size:20px">${t(s.source === 'admin' ? 'gwSrcAdmin' : 'gwSrcEnv')}</b></div>
      <div class="kpi"><span>${t('lgGG')}</span><b style="font-size:20px">${t(g.active ? 'lgOn' : 'lgOff')}</b></div>
      <div class="kpi"><span>${t('lgFB')}</span><b style="font-size:20px">${t(f.active ? 'lgOn' : 'lgOff')}</b></div>
    </div>
    <section class="panel" id="lgSMSForm"><div class="ph"><h2>${t('lgSMS')}</h2>${s.key_set ? `<span class="chip ${s.verified_at ? 'ok' : 'warn'}">${esc(lgConn(s))}</span>` : ''}</div>
      <div style="padding:16px 18px;display:flex;flex-direction:column;gap:14px">
        ${s.source === 'env' ? `<div class="banner info">${t('lgEnvNote', {m: esc(s.env_mode)})}</div>` : ''}
        <div class="opts" role="radiogroup" aria-label="${t('lgSMS')}">${modes.map(radio).join('')}</div>
        <div style="display:grid;grid-template-columns:repeat(auto-fit,minmax(220px,1fr));gap:14px">
          <label>${t('lgPK')}<input id="lgPK" name="private_key" type="password" autocomplete="new-password" spellcheck="false" placeholder="${esc(s.key_set ? t('lgKeyKeep', {h: '…' + s.key_hint}) : t('lgKeyNew'))}"></label>
          <label>${t('lgSecret')}<input id="lgSecret" name="secret" type="password" autocomplete="new-password" spellcheck="false" placeholder="${esc(s.key_set ? t('lgKeyKeep', {h: '…' + s.key_hint}) : t('lgKeyNew'))}"></label>
          <label>${t('lgSender')}<input id="lgSender" name="sender" maxlength="15" value="${esc(s.sender || '')}" placeholder="Phsar"><small class="muted" style="font-weight:400">${t('lgSenderHint')}</small></label>
        </div>
        <p class="muted small">${t('gwCodeHint')}</p>
        <div style="display:flex;flex-wrap:wrap;gap:10px;align-items:flex-end">${code('lgCode1')}
          <button class="primary" id="lgSaveSMS" type="button">${t('lgSaveSMS')}</button>
          ${s.source === 'admin' ? `<button type="button" id="lgFollowEnv">${t('lgFollowEnv')}</button>` : ''}</div>
        ${s.key_set ? `<div style="display:flex;flex-wrap:wrap;gap:10px;align-items:flex-end;border-top:1px solid var(--line);padding-top:14px">
          <label>${t('lgTestPhone')}<input id="lgTestPhone" name="test_phone" inputmode="tel" placeholder="012 345 678" style="max-width:200px"></label>
          <button type="button" id="lgTestSMS">${t('lgTestSMS')}</button><span class="muted small">${t('lgTestCost')}</span></div>` : ''}
        ${result(s.last_test)}
      </div></section>
    <section class="panel" id="lgGGForm"><div class="ph"><h2>${t('lgGG')}</h2>${g.key_set ? `<span class="chip ${g.verified_at ? 'ok' : 'warn'}">${esc(lgConn(g))}</span>` : ''}</div>
      <div style="padding:16px 18px;display:flex;flex-direction:column;gap:14px">
        <label class="opt" style="flex-direction:row;align-items:center"><input type="checkbox" id="lgGGOn" ${g.enabled ? 'checked' : ''} style="width:18px;height:18px;min-height:0;flex:none;accent-color:var(--pink)"><span><b>${t('lgGGEnable')}</b></span></label>
        <div style="display:grid;grid-template-columns:repeat(auto-fit,minmax(240px,1fr));gap:14px">
          <label>${t('lgGGID')}<input id="lgGGID" name="client_id" spellcheck="false" value="${esc(g.client_id || '')}" placeholder="1234567890-abc….apps.googleusercontent.com"></label>
          <label>${t('lgGGSecret')}<input id="lgGGSecret" name="client_secret" type="password" autocomplete="new-password" spellcheck="false" placeholder="${esc(g.key_set ? t('lgAppKeep', {h: '…' + g.key_hint}) : t('lgGGNew'))}"></label>
        </div>
        <label>${t('lgGGApps')}<textarea id="lgGGApps" name="app_client_ids" rows="2" spellcheck="false">${esc((g.app_client_ids || '').split(',').filter(Boolean).join('\n'))}</textarea><small class="muted" style="font-weight:400">${t('lgGGAppsHint')}</small></label>
        <div><span class="muted small">${t('lgGGRedirects')}</span>
          <ul style="margin:6px 0 0;padding-left:18px;display:flex;flex-direction:column;gap:4px">${g.redirect_uris.map(u => `<li><code style="overflow-wrap:anywhere">${esc(u)}</code></li>`).join('')}</ul>
          <span class="muted small" style="display:block;margin-top:8px">${t('lgGGOrigins')}</span>
          <ul style="margin:6px 0 0;padding-left:18px;display:flex;flex-direction:column;gap:4px">${g.origins.map(u => `<li><code style="overflow-wrap:anywhere">${esc(u)}</code></li>`).join('')}</ul>
          <button type="button" id="lgCopyGG" style="margin-top:8px">${t('copy')}</button></div>
        ${result(g.last_test)}
        <p class="muted small">${t('gwCodeHint')}</p>
        <div style="display:flex;flex-wrap:wrap;gap:10px;align-items:flex-end">${code('lgCode3')}
          <button class="primary" id="lgSaveGG" type="button">${t('lgSaveGG')}</button>
          <button type="button" id="lgTestGG" ${g.key_set ? '' : 'disabled'}>${t('lgTestGG')}</button></div>
      </div></section>
    <section class="panel" id="lgFBForm"><div class="ph"><h2>${t('lgFB')}</h2>${f.key_set ? `<span class="chip ${f.verified_at ? 'ok' : 'warn'}">${esc(lgConn(f))}</span>` : ''}</div>
      <div style="padding:16px 18px;display:flex;flex-direction:column;gap:14px">
        <label class="opt" style="flex-direction:row;align-items:center"><input type="checkbox" id="lgFBOn" ${f.enabled ? 'checked' : ''} style="width:18px;height:18px;min-height:0;flex:none;accent-color:var(--pink)"><span><b>${t('lgFBEnable')}</b></span></label>
        <div style="display:grid;grid-template-columns:repeat(auto-fit,minmax(240px,1fr));gap:14px">
          <label>${t('lgAppID')}<input id="lgAppID" name="app_id" inputmode="numeric" spellcheck="false" value="${esc(f.app_id || '')}"></label>
          <label>${t('lgAppSecret')}<input id="lgAppSecret" name="app_secret" type="password" autocomplete="new-password" spellcheck="false" placeholder="${esc(f.key_set ? t('lgAppKeep', {h: '…' + f.key_hint}) : t('lgAppNew'))}"></label>
        </div>
        <div><span class="muted small">${t('lgRedirects')}</span>
          <ul style="margin:6px 0 0;padding-left:18px;display:flex;flex-direction:column;gap:4px">${f.redirect_uris.map(u => `<li><code style="overflow-wrap:anywhere">${esc(u)}</code></li>`).join('')}</ul>
          <button type="button" id="lgCopy" style="margin-top:8px">${t('copy')}</button></div>
        ${result(f.last_test)}
        <p class="muted small">${t('gwCodeHint')}</p>
        <div style="display:flex;flex-wrap:wrap;gap:10px;align-items:flex-end">${code('lgCode2')}
          <button class="primary" id="lgSaveFB" type="button">${t('lgSaveFB')}</button>
          <button type="button" id="lgTestFB" ${f.key_set ? '' : 'disabled'}>${t('lgTestFB')}</button></div>
      </div></section>
    <details class="panel" id="lgAPForm" ${ap.enabled ? 'open' : ''}><summary class="ph" style="cursor:pointer"><h2>${t('lgAP')}</h2><span class="chip ${ap.active ? 'ok' : ''}">${t(ap.active ? 'lgOn' : 'lgOff')}</span></summary>
      <div style="padding:16px 18px;display:flex;flex-direction:column;gap:14px">
        <p class="muted small">${t('lgAPNote')}</p>
        <label class="opt" style="flex-direction:row;align-items:center"><input type="checkbox" id="lgAPOn" ${ap.enabled ? 'checked' : ''} style="width:18px;height:18px;min-height:0;flex:none;accent-color:var(--pink)"><span><b>${t('lgAPEnable')}</b></span></label>
        <label>${t('lgAPIDs')}<input id="lgAPIDs" name="bundle_ids" spellcheck="false" value="${esc(ap.bundle_ids || '')}" placeholder="com.phsar.app" style="max-width:360px"></label>
        <div style="display:flex;flex-wrap:wrap;gap:10px;align-items:flex-end">${code('lgCode4')}
          <button class="primary" id="lgSaveAP" type="button">${t('lgSaveAP')}</button></div>
      </div></details>`);

  const fail = (e, form) => {
    if (e.field) { const el = $(`#${form} [name="${e.field}"]`); if (el) { el.classList.add('invalid'); el.focus(); } }
    toast(errMsg(e), 'error');
  };
  const clear = form => document.querySelectorAll(`#${form} .invalid`).forEach(el => el.classList.remove('invalid'));
  const saveSMS = async mode => {
    clear('lgSMSForm');
    const body = {mode, sender: $('#lgSender').value.trim(), mfa_code: $('#lgCode1').value};
    if ($('#lgPK').value.trim() || $('#lgSecret').value.trim()) { body.private_key = $('#lgPK').value.trim(); body.secret = $('#lgSecret').value.trim(); }
    try { LG.data = await api('/admin/api/integrations/sms', {method: 'PUT', body}); toast(t('lgSaved')); drawLogins(); } catch (e) { fail(e, 'lgSMSForm'); }
  };
  $('#lgSaveSMS').onclick = () => saveSMS((document.querySelector('input[name="lgMode"]:checked') || {}).value || 'off');
  if ($('#lgFollowEnv')) $('#lgFollowEnv').onclick = () => saveSMS('');
  if ($('#lgTestSMS')) $('#lgTestSMS').onclick = async () => {
    const b = $('#lgTestSMS'); b.disabled = true; b.textContent = t('lgTesting');
    try { await api('/admin/api/integrations/sms/test', {method: 'POST', body: {phone: $('#lgTestPhone').value}}); LG.data = await api('/admin/api/integrations'); drawLogins(); }
    catch (e) { fail(e, 'lgSMSForm'); b.disabled = false; b.textContent = t('lgTestSMS'); }
  };
  $('#lgCopy').onclick = async () => { try { await navigator.clipboard.writeText(f.redirect_uris.join('\n')); toast(t('copied')); } catch {} };
  $('#lgSaveFB').onclick = async () => {
    clear('lgFBForm');
    const body = {enabled: $('#lgFBOn').checked, app_id: $('#lgAppID').value.trim(), mfa_code: $('#lgCode2').value};
    if ($('#lgAppSecret').value.trim()) body.app_secret = $('#lgAppSecret').value.trim();
    try { LG.data = await api('/admin/api/integrations/facebook', {method: 'PUT', body}); toast(t('lgSaved')); drawLogins(); } catch (e) { fail(e, 'lgFBForm'); }
  };
  $('#lgTestFB').onclick = async () => {
    const b = $('#lgTestFB'); b.disabled = true; b.textContent = t('lgTesting');
    try { await api('/admin/api/integrations/facebook/test', {method: 'POST'}); LG.data = await api('/admin/api/integrations'); drawLogins(); }
    catch (e) { toast(errMsg(e), 'error'); b.disabled = false; b.textContent = t('lgTestFB'); }
  };
  $('#lgCopyGG').onclick = async () => { try { await navigator.clipboard.writeText(g.redirect_uris.concat(g.origins).join('\n')); toast(t('copied')); } catch {} };
  $('#lgSaveGG').onclick = async () => {
    clear('lgGGForm');
    const body = {enabled: $('#lgGGOn').checked, client_id: $('#lgGGID').value.trim(), app_client_ids: $('#lgGGApps').value, mfa_code: $('#lgCode3').value};
    if ($('#lgGGSecret').value.trim()) body.client_secret = $('#lgGGSecret').value.trim();
    try { LG.data = await api('/admin/api/integrations/google', {method: 'PUT', body}); toast(t('lgSaved')); drawLogins(); } catch (e) { fail(e, 'lgGGForm'); }
  };
  $('#lgTestGG').onclick = async () => {
    const b = $('#lgTestGG'); b.disabled = true; b.textContent = t('lgTesting');
    try { await api('/admin/api/integrations/google/test', {method: 'POST'}); LG.data = await api('/admin/api/integrations'); drawLogins(); }
    catch (e) { toast(errMsg(e), 'error'); b.disabled = false; b.textContent = t('lgTestGG'); }
  };
  $('#lgSaveAP').onclick = async () => {
    clear('lgAPForm');
    const body = {enabled: $('#lgAPOn').checked, bundle_ids: $('#lgAPIDs').value.trim(), mfa_code: $('#lgCode4').value};
    try { LG.data = await api('/admin/api/integrations/apple', {method: 'PUT', body}); toast(t('lgSaved')); drawLogins(); } catch (e) { fail(e, 'lgAPForm'); }
  };
}
