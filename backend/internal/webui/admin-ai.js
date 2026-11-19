'use strict';
// Super admin: halaman Integrasi — semua kunci pihak ketiga di satu tempat (tab): Pembayaran · Login & SMS · AI · Peta & lainnya.
// Tab AI: Anthropic (utama) lalu Gemini (cadangan, dipakai kalau Anthropic gagal mis. HTTP 400 saldo habis).
addStrings({
  en: {i_integrations: 'Integrations', intDesc: 'All third-party keys in one place. Only super admins. Keys are encrypted on the server and never shown again — only the last 4 characters.',
    intTabPay: 'Payments', intTabLogin: 'Login & SMS', intTabAI: 'AI & translation', intTabMaps: 'Maps & other',
    aiFlow: 'Order used for translations and reply drafts', aiStep1: 'Anthropic (Claude) — main', aiStep2: 'Gemini — backup, used when Anthropic fails (e.g. HTTP 400 “credit balance too low”, 429, server error)',
    aiStep3: 'No AI key at all → translations wait in the queue; reply drafts use simple rules',
    aiMain: 'Main', aiBackup: 'Backup', aiOn: 'Use this provider', aiModel: 'Model', aiKey: 'API key', aiKeep: 'Saved — ends with {h}. Leave empty to keep.',
    aiKeyNewA: 'console.anthropic.com → API keys (starts with sk-ant-)', aiKeyNewG: 'aistudio.google.com → Get API key',
    aiClear: 'Delete the saved key', aiSave: 'Save', aiTest: 'Test with a short translation', aiSrcAdmin: 'Key from this page', aiSrcEnv: 'Key from server .env',
    aiNoKey: 'No key', aiActive: 'Active', aiOff: 'Turned off', aiModelHintA: 'Default: claude-haiku-5-5', aiModelHintG: 'Default: gemini-3.5-flash (free tier)',
    aiFreeNote: 'Gemini free tier: Google may use the text to improve its products — fine for public product descriptions, not for private data.',
    aiOrderNow: 'Right now: {o}', aiOrderNone: 'Right now: no AI key — translations are waiting.',
    mpTitle: 'Google Maps (pickup points, addresses)', mpKey: 'Browser API key', mpKeyNew: 'Google Cloud → APIs & Services → Credentials (starts with AIza). Restrict it to your domain.',
    mpMapID: 'Map ID', mpMapIDHint: 'Google Cloud → Map Management. DEMO_MAP_ID works for testing.', mpSave: 'Save maps settings', mpActive: 'Maps', mpOff: 'Not set — the simple map is used',
    tgTitle: 'Telegram bot (login & notifications)', tgSet: 'Bot token is set in the server .env ({u}).', tgUnset: 'Not set. Put TELEGRAM_BOT_TOKEN and TELEGRAM_BOT_USERNAME in the server .env.',
    tgWhy: 'The Telegram token stays in .env because the bot connection starts with the server — change it there and restart.'},
  km: {i_integrations: 'ការភ្ជាប់សេវា', intDesc: 'កូនសោសេវាភាគីទីបីទាំងអស់នៅកន្លែងតែមួយ។ សម្រាប់អ្នកគ្រប់គ្រងកំពូលប៉ុណ្ណោះ។ កូនសោត្រូវបានអ៊ិនគ្រីប ហើយមិនបង្ហាញម្តងទៀតទេ — មានតែ 4 តួចុងក្រោយ។',
    intTabPay: 'ការទូទាត់', intTabLogin: 'ការចូល និង SMS', intTabAI: 'AI និងការបកប្រែ', intTabMaps: 'ផែនទី និងផ្សេងៗ',
    aiFlow: 'លំដាប់ដែលប្រើសម្រាប់ការបកប្រែ និងសេចក្តីព្រាងចម្លើយ', aiStep1: 'Anthropic (Claude) — ចម្បង', aiStep2: 'Gemini — បម្រុង ប្រើពេល Anthropic បរាជ័យ (ឧ. HTTP 400 សមតុល្យអស់, 429, កំហុសម៉ាស៊ីនមេ)',
    aiStep3: 'គ្មានកូនសោ AI ទាល់តែសោះ → ការបកប្រែរង់ចាំក្នុងជួរ; សេចក្តីព្រាងចម្លើយប្រើច្បាប់សាមញ្ញ',
    aiMain: 'ចម្បង', aiBackup: 'បម្រុង', aiOn: 'ប្រើអ្នកផ្តល់សេវានេះ', aiModel: 'ម៉ូដែល', aiKey: 'API key', aiKeep: 'បានរក្សាទុក — បញ្ចប់ដោយ {h}។ ទុកទទេដើម្បីរក្សាដដែល។',
    aiKeyNewA: 'console.anthropic.com → API keys (ចាប់ផ្តើមដោយ sk-ant-)', aiKeyNewG: 'aistudio.google.com → Get API key',
    aiClear: 'លុបកូនសោដែលបានរក្សាទុក', aiSave: 'រក្សាទុក', aiTest: 'សាកល្បងដោយបកប្រែខ្លី', aiSrcAdmin: 'កូនសោពីទំព័រនេះ', aiSrcEnv: 'កូនសោពី .env ម៉ាស៊ីនមេ',
    aiNoKey: 'គ្មានកូនសោ', aiActive: 'កំពុងប្រើ', aiOff: 'បានបិទ', aiModelHintA: 'លំនាំដើម៖ claude-haiku-5-5', aiModelHintG: 'លំនាំដើម៖ gemini-3.5-flash (ឥតគិតថ្លៃ)',
    aiFreeNote: 'Gemini ឥតគិតថ្លៃ៖ Google អាចប្រើអត្ថបទដើម្បីកែលម្អផលិតផលរបស់ខ្លួន — មិនអីសម្រាប់ការពិពណ៌នាផលិតផលសាធារណៈ តែកុំប្រើសម្រាប់ទិន្នន័យឯកជន។',
    aiOrderNow: 'ឥឡូវនេះ៖ {o}', aiOrderNone: 'ឥឡូវនេះ៖ គ្មានកូនសោ AI — ការបកប្រែកំពុងរង់ចាំ។',
    mpTitle: 'Google Maps (ទីតាំងទទួល អាសយដ្ឋាន)', mpKey: 'Browser API key', mpKeyNew: 'Google Cloud → APIs & Services → Credentials (ចាប់ផ្តើមដោយ AIza)។ កំណត់ឱ្យប្រើតែដែនរបស់អ្នក។',
    mpMapID: 'Map ID', mpMapIDHint: 'Google Cloud → Map Management។ DEMO_MAP_ID ប្រើបានសម្រាប់សាកល្បង។', mpSave: 'រក្សាទុកការកំណត់ផែនទី', mpActive: 'ផែនទី', mpOff: 'មិនទាន់កំណត់ — ប្រើផែនទីសាមញ្ញ',
    tgTitle: 'Telegram bot (ការចូល និងការជូនដំណឹង)', tgSet: 'Bot token ត្រូវបានកំណត់ក្នុង .env ម៉ាស៊ីនមេ ({u})។', tgUnset: 'មិនទាន់កំណត់។ ដាក់ TELEGRAM_BOT_TOKEN និង TELEGRAM_BOT_USERNAME ក្នុង .env ម៉ាស៊ីនមេ។',
    tgWhy: 'Token Telegram នៅក្នុង .env ព្រោះ bot ភ្ជាប់ពេលម៉ាស៊ីនមេចាប់ផ្តើម — ប្តូរនៅទីនោះ ហើយចាប់ផ្តើមម្តងទៀត។'},
  zh: {i_integrations: '第三方集成', intDesc: '所有第三方密钥集中在一处。仅限超级管理员。密钥在服务器上加密保存，保存后不再显示——只显示最后 4 位。',
    intTabPay: '支付', intTabLogin: '登录与短信', intTabAI: 'AI 与翻译', intTabMaps: '地图及其他',
    aiFlow: '翻译和回复草稿使用的顺序', aiStep1: 'Anthropic (Claude) — 主要', aiStep2: 'Gemini — 备用，Anthropic 失败时使用（如 HTTP 400 余额不足、429、服务器错误）',
    aiStep3: '完全没有 AI 密钥 → 翻译在队列中等待；回复草稿使用简单规则',
    aiMain: '主要', aiBackup: '备用', aiOn: '使用此服务', aiModel: '模型', aiKey: 'API key', aiKeep: '已保存——以 {h} 结尾。留空则保持不变。',
    aiKeyNewA: 'console.anthropic.com → API keys（以 sk-ant- 开头）', aiKeyNewG: 'aistudio.google.com → Get API key',
    aiClear: '删除已保存的密钥', aiSave: '保存', aiTest: '用简短翻译测试', aiSrcAdmin: '来自本页的密钥', aiSrcEnv: '来自服务器 .env 的密钥',
    aiNoKey: '无密钥', aiActive: '使用中', aiOff: '已关闭', aiModelHintA: '默认：claude-haiku-5-5', aiModelHintG: '默认：gemini-3.5-flash（免费额度）',
    aiFreeNote: 'Gemini 免费额度：Google 可能用这些文本改进其产品——公开的商品描述没问题，私人数据不要使用。',
    aiOrderNow: '当前：{o}', aiOrderNone: '当前：没有 AI 密钥——翻译正在等待。',
    mpTitle: 'Google 地图（取件点、地址）', mpKey: '浏览器 API key', mpKeyNew: 'Google Cloud → APIs & Services → Credentials（以 AIza 开头）。请限制为你的域名。',
    mpMapID: 'Map ID', mpMapIDHint: 'Google Cloud → Map Management。测试可用 DEMO_MAP_ID。', mpSave: '保存地图设置', mpActive: '地图', mpOff: '未设置——使用简易地图',
    tgTitle: 'Telegram 机器人（登录与通知）', tgSet: 'Bot token 已在服务器 .env 中设置（{u}）。', tgUnset: '未设置。请在服务器 .env 中填写 TELEGRAM_BOT_TOKEN 和 TELEGRAM_BOT_USERNAME。',
    tgWhy: 'Telegram token 保留在 .env 中，因为机器人连接随服务器启动——请在那里修改并重启。'},
});

const AIX = {data: null};

// Tab bersama untuk semua halaman Integrasi (dipanggil juga oleh halaman Pembayaran dan Login & SMS).
function intTabs(active) {
  const tabs = [['gateway', '#/gateway', 'intTabPay'], ['logins', '#/logins', 'intTabLogin'], ['ai', '#/integrations/ai', 'intTabAI'], ['maps', '#/integrations/maps', 'intTabMaps'], ['usage', '#/integrations/usage', 'intTabUsage']];
  return `<nav class="tabs" aria-label="${t('i_integrations')}" style="margin-bottom:4px">${tabs.map(([k, h, l]) =>
    `<a class="tab ${k === active ? 'on' : ''}" href="${h}" ${k === active ? 'aria-current="page"' : ''}>${t(l)}</a>`).join('')}</nav>`;
}

async function renderIntegrations() {
  if (location.hash.startsWith('#/integrations/usage')) return renderUsage();
  AIX.data = await api('/admin/api/integrations');
  if (location.hash.startsWith('#/integrations/maps')) drawMapsTab(); else drawAITab();
}

function aixCode(id) {
  return `<label>${t('gwCode')}<input id="${id}" name="mfa_code" inputmode="numeric" autocomplete="one-time-code" maxlength="6" placeholder="123456" style="max-width:180px"></label>`;
}
function aixResult(x) {
  return x ? `<div class="banner ${x.ok ? 'ok' : 'bad'}">${esc(x.message)} <span class="muted small">· ${new Date(x.at).toLocaleString()}</span></div>` : '';
}
function aixFail(e, form) {
  if (e.field) { const el = $(`#${form} [name="${e.field}"]`); if (el) { el.classList.add('invalid'); el.focus(); } }
  toast(errMsg(e), 'error');
}

function drawAITab() {
  const d = AIX.data.ai, names = {anthropic: 'Anthropic', gemini: 'Gemini'};
  const card = (k, p, role) => {
    const status = !p.enabled ? `<span class="chip">${t('aiOff')}</span>` : p.active ? `<span class="chip ok">${t('aiActive')} · ${t(p.source === 'admin' ? 'aiSrcAdmin' : 'aiSrcEnv')}</span>` : `<span class="chip warn">${t('aiNoKey')}</span>`;
    return `<section class="panel" id="ai_${k}"><div class="ph"><h2>${role === 'main' ? '1' : '2'} · ${names[k]} <span class="chip ${role === 'main' ? 'pink' : ''}">${t(role === 'main' ? 'aiMain' : 'aiBackup')}</span></h2>${status}</div>
      <div style="padding:16px 18px;display:flex;flex-direction:column;gap:14px">
        <label class="opt" style="flex-direction:row;align-items:center"><input type="checkbox" name="enabled" ${p.enabled ? 'checked' : ''} style="width:18px;height:18px;min-height:0;flex:none;accent-color:var(--pink)"><span><b>${t('aiOn')}</b></span></label>
        <div style="display:grid;grid-template-columns:repeat(auto-fit,minmax(240px,1fr));gap:14px">
          <label>${t('aiKey')}<input name="key" type="password" autocomplete="new-password" spellcheck="false" placeholder="${esc(p.key_set ? t('aiKeep', {h: '…' + p.key_hint}) : t(k === 'anthropic' ? 'aiKeyNewA' : 'aiKeyNewG'))}"></label>
          <label>${t('aiModel')}<input name="model" spellcheck="false" value="${esc(p.model || '')}"><small class="muted" style="font-weight:400">${t(k === 'anthropic' ? 'aiModelHintA' : 'aiModelHintG')}</small></label>
        </div>
        ${p.key_set ? `<label class="opt" style="flex-direction:row;align-items:center"><input type="checkbox" name="clear" style="width:18px;height:18px;min-height:0;flex:none"><span>${t('aiClear')}</span></label>` : ''}
        ${k === 'gemini' ? `<p class="muted small">${t('aiFreeNote')}</p>` : ''}
        ${aixResult(p.last_test)}
        <p class="muted small">${t('gwCodeHint')}</p>
        <div style="display:flex;flex-wrap:wrap;gap:10px;align-items:flex-end">${aixCode('aiCode_' + k)}
          <button class="primary" type="button" data-aisave="${k}">${t('aiSave')}</button>
          <button type="button" data-aitest="${k}" ${p.active ? '' : 'disabled'}>${t('aiTest')}</button></div>
      </div></section>`;
  };
  const order = (d.order || []).map(k => names[k]).join(' → ');
  $('#app').innerHTML = shell(`${pageHead(t('i_integrations'), t('intDesc'))}${intTabs('ai')}
    <section class="panel"><div class="ph"><h2>${t('aiFlow')}</h2><span class="chip ${order ? 'ok' : 'warn'}">${order ? esc(t('aiOrderNow', {o: order})) : t('aiOrderNone')}</span></div>
      <ol style="margin:0;padding:14px 18px 16px 38px;display:flex;flex-direction:column;gap:6px"><li>${t('aiStep1')}</li><li>${t('aiStep2')}</li><li class="muted">${t('aiStep3')}</li></ol></section>
    ${card('anthropic', d.anthropic, 'main')}
    ${card('gemini', d.gemini, 'backup')}`);
  document.querySelectorAll('[data-aisave]').forEach(b => b.onclick = async () => {
    const k = b.dataset.aisave, f = $('#ai_' + k);
    f.querySelectorAll('.invalid').forEach(el => el.classList.remove('invalid'));
    const body = {provider: k, enabled: f.querySelector('[name=enabled]').checked, model: f.querySelector('[name=model]').value.trim(),
      clear: !!(f.querySelector('[name=clear]') || {}).checked, mfa_code: $('#aiCode_' + k).value};
    const key = f.querySelector('[name=key]').value.trim();
    if (key) body.key = key;
    try { AIX.data = await api('/admin/api/integrations/ai', {method: 'PUT', body}); toast(t('lgSaved')); drawAITab(); } catch (e) { aixFail(e, 'ai_' + k); }
  });
  document.querySelectorAll('[data-aitest]').forEach(b => b.onclick = async () => {
    b.disabled = true; b.textContent = t('lgTesting');
    try { await api('/admin/api/integrations/ai/test', {method: 'POST', body: {provider: b.dataset.aitest}}); } catch (e) { toast(errMsg(e), 'error'); }
    AIX.data = await api('/admin/api/integrations'); drawAITab();
  });
}

function drawMapsTab() {
  const m = AIX.data.maps, tg = AIX.data.telegram;
  $('#app').innerHTML = shell(`${pageHead(t('i_integrations'), t('intDesc'))}${intTabs('maps')}
    <section class="panel" id="mpForm"><div class="ph"><h2>${t('mpTitle')}</h2>${m.active ? `<span class="chip ok">${t('aiActive')} · ${t(m.source === 'admin' ? 'aiSrcAdmin' : 'aiSrcEnv')}</span>` : `<span class="chip warn">${t('mpOff')}</span>`}</div>
      <div style="padding:16px 18px;display:flex;flex-direction:column;gap:14px">
        <div style="display:grid;grid-template-columns:repeat(auto-fit,minmax(240px,1fr));gap:14px">
          <label>${t('mpKey')}<input name="key" type="password" autocomplete="new-password" spellcheck="false" placeholder="${esc(m.key_set ? t('aiKeep', {h: '…' + m.key_hint}) : 'AIza…')}"><small class="muted" style="font-weight:400">${t('mpKeyNew')}</small></label>
          <label>${t('mpMapID')}<input name="map_id" spellcheck="false" value="${esc(m.map_id || '')}"><small class="muted" style="font-weight:400">${t('mpMapIDHint')}</small></label>
        </div>
        ${m.key_set ? `<label class="opt" style="flex-direction:row;align-items:center"><input type="checkbox" name="clear" style="width:18px;height:18px;min-height:0;flex:none"><span>${t('aiClear')}</span></label>` : ''}
        <p class="muted small">${t('gwCodeHint')}</p>
        <div style="display:flex;flex-wrap:wrap;gap:10px;align-items:flex-end">${aixCode('mpCode')}<button class="primary" type="button" id="mpSave">${t('mpSave')}</button></div>
      </div></section>
    <section class="panel"><div class="ph"><h2>${t('tgTitle')}</h2><span class="chip ${tg.bot_set ? 'ok' : 'warn'}">${tg.bot_set ? t('aiActive') : t('aiNoKey')}</span></div>
      <div style="padding:16px 18px;display:flex;flex-direction:column;gap:8px">
        <p>${tg.bot_set ? t('tgSet', {u: esc(tg.username ? '@' + tg.username : '—')}) : t('tgUnset')}</p><p class="muted small">${t('tgWhy')}</p></div></section>`);
  $('#mpSave').onclick = async () => {
    const f = $('#mpForm');
    f.querySelectorAll('.invalid').forEach(el => el.classList.remove('invalid'));
    const body = {map_id: f.querySelector('[name=map_id]').value.trim(), clear: !!(f.querySelector('[name=clear]') || {}).checked, mfa_code: $('#mpCode').value};
    const key = f.querySelector('[name=key]').value.trim();
    if (key) body.key = key;
    try { AIX.data = await api('/admin/api/integrations/maps', {method: 'PUT', body}); toast(t('lgSaved')); drawMapsTab(); } catch (e) { aixFail(e, 'mpForm'); }
  };
}
