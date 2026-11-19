'use strict';
// Super admin: Integrasi → Pemakaian & saldo. Laporan pemakaian AI (Anthropic, Gemini) & SMS (PlasGate) + pengingat saldo.
// Penyedia tidak memberi API saldo, jadi saldo = angka dari console penyedia yang diisi admin, dikurangi perkiraan biaya sejak itu.
addStrings({
  en: {intTabUsage: 'Usage & balance', q_integrations: 'Integration alerts (AI / SMS)',
    usDesc: 'Every AI and SMS call is counted here. Providers do not share your balance with us, so enter the balance you see in their console — the server subtracts the estimated cost of each call since then and warns you before it runs out.',
    usToday: 'Today', us7: '7 days', us30: '30 days', usReq: 'requests', usFail: 'failed', usTok: 'tokens in / out', usSms: 'SMS sent', usCost: 'est. cost',
    usBal: 'Estimated balance', usBalNone: 'Balance not entered', usBalOf: 'of {b} entered {t}', usSpent: 'used since then: {c}',
    usLastErr: 'Last error', usLastOK: 'Last success', usNever: 'never', usFallback: 'Gemini handled {n} requests in 7 days (backup when Anthropic failed or was off).',
    usBudget: 'Balance & reminder', usBalNew: 'Current balance in console (USD)', usBalHint: 'Leave empty to keep counting from the last balance. Saving a number starts counting again from now.',
    usWarn: 'Remind me below (USD)', usPin: 'Price input / 1M tokens (USD)', usPout: 'Price output / 1M tokens (USD)', usPsms: 'Price per SMS (USD)',
    usDaily: 'Free daily limit (requests)', usDailyHint: 'From AI Studio → rate limits. 0 = not set. You get a reminder at 80%.',
    usPriceHint: 'Check the provider’s price page; costs here are estimates, the provider’s invoice is the truth.', usOpenBilling: 'Open billing page',
    usSave: 'Save', usClearBal: 'Forget the balance', usDays: 'Last 14 days', usDay: 'Day', usProv: 'Service', usNoData: 'No calls recorded yet.',
    usAlerts: 'Needs attention', usAllOK: 'All services OK — no balance or error warnings.',
    al_billing: '{p}: balance used up — calls are failing with a billing error. Top up in the console, then enter the new balance here.',
    al_auth: '{p}: the API key was rejected. Check the key in AI & translation (or Login & SMS).',
    al_low_balance: '{p}: estimated balance is low (${v}). Top up soon.',
    al_quota: '{p}: free quota / rate limit reached — requests are retried automatically.',
    al_daily_limit: '{p}: {v} of today’s free requests used.',
    al_no_ai: 'No AI key — product translations are waiting. Add an Anthropic or Gemini key.',
    usSmsNote: 'SMS: counted per message part (160 Latin / 70 Khmer characters). Google Maps runs in the buyer’s browser — see its usage in Google Cloud.'},
  km: {intTabUsage: 'ការប្រើប្រាស់ និងសមតុល្យ', q_integrations: 'ការព្រមានសេវា (AI / SMS)',
    usDesc: 'រាល់ការហៅ AI និង SMS ត្រូវបានរាប់នៅទីនេះ។ អ្នកផ្តល់សេវាមិនចែករំលែកសមតុល្យជាមួយយើងទេ ដូច្នេះសូមបញ្ចូលសមតុល្យដែលអ្នកឃើញក្នុង console — ម៉ាស៊ីនមេដកថ្លៃប៉ាន់ស្មាននៃការហៅនីមួយៗ ហើយព្រមានអ្នកមុនពេលអស់។',
    usToday: 'ថ្ងៃនេះ', us7: '7 ថ្ងៃ', us30: '30 ថ្ងៃ', usReq: 'សំណើ', usFail: 'បរាជ័យ', usTok: 'token ចូល / ចេញ', usSms: 'SMS បានផ្ញើ', usCost: 'ថ្លៃប៉ាន់ស្មាន',
    usBal: 'សមតុល្យប៉ាន់ស្មាន', usBalNone: 'មិនទាន់បញ្ចូលសមតុល្យ', usBalOf: 'ពី {b} បញ្ចូល {t}', usSpent: 'បានប្រើតាំងពីពេលនោះ៖ {c}',
    usLastErr: 'កំហុសចុងក្រោយ', usLastOK: 'ជោគជ័យចុងក្រោយ', usNever: 'មិនដែល', usFallback: 'Gemini បានដោះស្រាយ {n} សំណើក្នុង 7 ថ្ងៃ (បម្រុងពេល Anthropic បរាជ័យ ឬបិទ)។',
    usBudget: 'សមតុល្យ និងការរំលឹក', usBalNew: 'សមតុល្យបច្ចុប្បន្នក្នុង console (USD)', usBalHint: 'ទុកទទេដើម្បីបន្តរាប់ពីសមតុល្យចុងក្រោយ។ ការរក្សាទុកលេខនឹងចាប់ផ្តើមរាប់ឡើងវិញពីឥឡូវ។',
    usWarn: 'រំលឹកខ្ញុំពេលក្រោម (USD)', usPin: 'តម្លៃ input / 1M token (USD)', usPout: 'តម្លៃ output / 1M token (USD)', usPsms: 'តម្លៃក្នុងមួយ SMS (USD)',
    usDaily: 'ដែនកំណត់ឥតគិតថ្លៃប្រចាំថ្ងៃ (សំណើ)', usDailyHint: 'ពី AI Studio → rate limits។ 0 = មិនកំណត់។ អ្នកនឹងទទួលការរំលឹកនៅ 80%។',
    usPriceHint: 'ពិនិត្យទំព័រតម្លៃរបស់អ្នកផ្តល់សេវា។ ថ្លៃនៅទីនេះជាការប៉ាន់ស្មាន វិក្កយបត្ររបស់អ្នកផ្តល់សេវាគឺជាការពិត។', usOpenBilling: 'បើកទំព័រវិក្កយបត្រ',
    usSave: 'រក្សាទុក', usClearBal: 'បំភ្លេចសមតុល្យ', usDays: '14 ថ្ងៃចុងក្រោយ', usDay: 'ថ្ងៃ', usProv: 'សេវា', usNoData: 'មិនទាន់មានការហៅត្រូវបានកត់ត្រា។',
    usAlerts: 'ត្រូវការយកចិត្តទុកដាក់', usAllOK: 'សេវាទាំងអស់ដំណើរការល្អ — គ្មានការព្រមានសមតុល្យ ឬកំហុស។',
    al_billing: '{p}៖ សមតុល្យអស់ — ការហៅបរាជ័យដោយកំហុសវិក្កយបត្រ។ បញ្ចូលទឹកប្រាក់ក្នុង console រួចបញ្ចូលសមតុល្យថ្មីនៅទីនេះ។',
    al_auth: '{p}៖ API key ត្រូវបានបដិសេធ។ ពិនិត្យកូនសោនៅ AI និងការបកប្រែ (ឬ ការចូល និង SMS)។',
    al_low_balance: '{p}៖ សមតុល្យប៉ាន់ស្មានទាប (${v})។ សូមបញ្ចូលទឹកប្រាក់ឆាប់ៗ។',
    al_quota: '{p}៖ ឈានដល់កូតាឥតគិតថ្លៃ / rate limit — សំណើនឹងព្យាយាមម្តងទៀតដោយស្វ័យប្រវត្តិ។',
    al_daily_limit: '{p}៖ បានប្រើ {v} នៃសំណើឥតគិតថ្លៃថ្ងៃនេះ។',
    al_no_ai: 'គ្មានកូនសោ AI — ការបកប្រែផលិតផលកំពុងរង់ចាំ។ បន្ថែមកូនសោ Anthropic ឬ Gemini។',
    usSmsNote: 'SMS៖ រាប់តាមផ្នែកសារ (160 តួឡាតាំង / 70 តួខ្មែរ)។ Google Maps ដំណើរការក្នុងកម្មវិធីរុករករបស់អ្នកទិញ — មើលការប្រើប្រាស់នៅ Google Cloud។'},
  zh: {intTabUsage: '用量与余额', q_integrations: '集成提醒（AI / 短信）',
    usDesc: '每次 AI 和短信调用都会在这里计数。服务商不会把余额提供给我们，所以请填写你在其控制台看到的余额——服务器会扣除此后每次调用的预估费用，并在用完之前提醒你。',
    usToday: '今天', us7: '7 天', us30: '30 天', usReq: '次请求', usFail: '次失败', usTok: 'token 输入 / 输出', usSms: '已发短信', usCost: '预估费用',
    usBal: '预估余额', usBalNone: '未填写余额', usBalOf: '共 {b}，填写于 {t}', usSpent: '此后已用：{c}',
    usLastErr: '最近错误', usLastOK: '最近成功', usNever: '从未', usFallback: '7 天内 Gemini 处理了 {n} 次请求（Anthropic 失败或关闭时的备用）。',
    usBudget: '余额与提醒', usBalNew: '控制台当前余额（USD）', usBalHint: '留空则继续按上次余额计算。保存数字后从现在重新开始计算。',
    usWarn: '低于此金额时提醒（USD）', usPin: '输入价格 / 1M token（USD）', usPout: '输出价格 / 1M token（USD）', usPsms: '每条短信价格（USD）',
    usDaily: '每日免费额度（请求数）', usDailyHint: '见 AI Studio → rate limits。0 = 未设置。用到 80% 时提醒。',
    usPriceHint: '请查看服务商的价格页面；这里的费用是估算，以服务商账单为准。', usOpenBilling: '打开账单页面',
    usSave: '保存', usClearBal: '清除余额', usDays: '最近 14 天', usDay: '日期', usProv: '服务', usNoData: '尚无调用记录。',
    usAlerts: '需要处理', usAllOK: '所有服务正常——没有余额或错误提醒。',
    al_billing: '{p}：余额已用完——调用因账单错误失败。请在控制台充值，然后在这里填写新余额。',
    al_auth: '{p}：API key 被拒绝。请在 AI 与翻译（或登录与短信）中检查密钥。',
    al_low_balance: '{p}：预估余额偏低（${v}）。请尽快充值。',
    al_quota: '{p}：已达免费额度 / 速率限制——请求会自动重试。',
    al_daily_limit: '{p}：今天的免费请求已用 {v}。',
    al_no_ai: '没有 AI 密钥——商品翻译正在等待。请添加 Anthropic 或 Gemini 密钥。',
    usSmsNote: '短信：按分段计数（拉丁字符 160 / 高棉文 70 个字符）。Google 地图在买家浏览器中运行——请在 Google Cloud 查看其用量。'},
});

const USG = {data: null};
const USG_NAMES = {anthropic: 'Anthropic', gemini: 'Gemini', plasgate: 'PlasGate SMS', ai: 'AI'};
const USG_BILLING = {anthropic: 'https://console.anthropic.com/settings/billing', gemini: 'https://aistudio.google.com/usage', plasgate: 'https://cloud.plasgate.com'};
const usgUsd = (v, d = 2) => '$' + Number(v || 0).toLocaleString('en-US', {minimumFractionDigits: d, maximumFractionDigits: d});
const usgSmall = v => Number(v || 0) > 0 && Number(v) < 0.01 ? '<$0.01' : usgUsd(v);
const usgNum = v => Number(v || 0).toLocaleString('en-US');

// Baris peringatan (dipakai juga oleh dasbor).
function usageAlertBanner(a, link) {
  return `<div class="banner ${a.level === 'bad' ? 'bad' : 'warn'}"><strong style="flex:1">${esc(t('al_' + a.code, {p: USG_NAMES[a.provider] || a.provider, v: a.value || ''}))}</strong>${link ? `<a href="#/integrations/usage" style="color:inherit;font-weight:600">${t('open')}</a>` : ''}</div>`;
}

async function renderUsage() {
  USG.data = await api('/admin/api/integrations/usage');
  drawUsage();
}

function drawUsage() {
  const d = USG.data;
  const alerts = d.alerts.length ? d.alerts.map(a => usageAlertBanner(a, false)).join('') : `<div class="banner ok"><strong>${t('usAllOK')}</strong></div>`;
  const stat = (p, label, sms) => `<div class="kpi" style="border-radius:16px;padding:12px 14px"><span>${label}</span><b style="font-size:20px">${usgSmall(p.cost_usd)}</b>
    <span>${usgNum(p.requests)} ${t('usReq')}${p.failed ? ` · <b style="font-size:13px;color:var(--bad,#c62828)">${usgNum(p.failed)} ${t('usFail')}</b>` : ''}</span>
    <span class="small">${sms ? `${usgNum(p.units)} ${t('usSms')}` : `${usgNum(p.input_tokens)} / ${usgNum(p.output_tokens)} ${t('usTok')}`}</span></div>`;
  const card = pu => {
    const b = pu.budget, sms = pu.provider === 'plasgate', k = pu.provider;
    let bal = `<p class="muted">${t('usBalNone')}</p>`;
    if (pu.remaining_usd != null) {
      const total = Number(b.balance_usd) || 0, left = pu.remaining_usd, pct = total > 0 ? Math.max(0, Math.min(100, left / total * 100)) : 0;
      const low = b.warn_below_usd !== '' && left < Number(b.warn_below_usd);
      bal = `<div style="display:flex;align-items:baseline;gap:10px;flex-wrap:wrap"><b style="font-size:26px;font-variant-numeric:tabular-nums;color:${left <= 0 ? 'var(--bad,#c62828)' : low ? 'var(--warn)' : 'inherit'}">${usgUsd(left)}</b>
        <span class="muted small">${esc(t('usBalOf', {b: usgUsd(total), t: new Date(b.balance_at).toLocaleDateString()}))} · ${esc(t('usSpent', {c: usgUsd(pu.spent_since, 4)}))}</span></div>
        <div role="meter" aria-label="${t('usBal')}" aria-valuemin="0" aria-valuemax="100" aria-valuenow="${Math.round(pct)}" style="height:8px;border-radius:8px;background:var(--line);overflow:hidden;margin-top:6px">
          <div style="height:100%;width:${pct}%;background:${left <= 0 || pct < 15 ? 'var(--bad,#c62828)' : low ? 'var(--warn)' : 'var(--ok,#2e7d32)'}"></div></div>`;
    }
    const le = pu.last_error;
    const field = (name, label, val, hint) => `<label>${label}<input name="${name}" inputmode="decimal" spellcheck="false" value="${esc(val == null ? '' : String(val))}">${hint ? `<small class="muted" style="font-weight:400">${hint}</small>` : ''}</label>`;
    return `<section class="panel" id="us_${k}"><div class="ph"><h2>${USG_NAMES[k]}</h2><a class="btn" href="${USG_BILLING[k]}" target="_blank" rel="noopener noreferrer">${t('usOpenBilling')} ↗</a></div>
      <div style="padding:16px 18px;display:flex;flex-direction:column;gap:14px">
        ${k === 'gemini' && !b.balance_usd ? '' : `<div><div class="muted small" style="margin-bottom:4px">${t('usBal')}</div>${bal}</div>`}
        <div style="display:grid;grid-template-columns:repeat(auto-fit,minmax(96px,1fr));gap:8px">${stat(pu.today, t('usToday'), sms)}${stat(pu.week, t('us7'), sms)}${stat(pu.month, t('us30'), sms)}</div>
        ${k === 'gemini' && pu.fallbacks_7d ? `<p class="muted small">${t('usFallback', {n: pu.fallbacks_7d})}</p>` : ''}
        <p class="small">${t('usLastOK')}: <b>${pu.last_ok ? ago(pu.last_ok) : t('usNever')}</b>${le ? ` · ${t('usLastErr')}: <b>${ago(le.at)}</b> <span class="chip ${le.kind === 'billing' || le.kind === 'auth' ? 'bad' : 'warn'}">${esc(le.kind)}</span><br><span class="muted small" style="word-break:break-word">${esc(le.text)}</span>` : ''}</p>
        <details><summary style="cursor:pointer;font-weight:600">${t('usBudget')}</summary>
          <div style="display:flex;flex-direction:column;gap:12px;margin-top:12px">
            <div style="display:grid;grid-template-columns:repeat(auto-fit,minmax(200px,1fr));gap:12px">
              ${k === 'gemini' ? field('daily_limit', t('usDaily'), b.daily_limit || 0, t('usDailyHint')) : field('balance_usd', t('usBalNew'), '', t('usBalHint'))}
              ${k === 'gemini' ? '' : field('warn_below_usd', t('usWarn'), b.warn_below_usd)}
              ${sms ? field('price_per_sms', t('usPsms'), b.price_per_sms) : field('price_in', t('usPin'), b.price_in) + field('price_out', t('usPout'), b.price_out)}
            </div>
            <p class="muted small">${t('usPriceHint')}</p>
            ${b.balance_usd && k !== 'gemini' ? `<label class="opt" style="flex-direction:row;align-items:center"><input type="checkbox" name="clear_balance" style="width:18px;height:18px;min-height:0;flex:none"><span>${t('usClearBal')}</span></label>` : ''}
            <p class="muted small">${t('gwCodeHint')}</p>
            <div style="display:flex;flex-wrap:wrap;gap:10px;align-items:flex-end">${aixCode('usCode_' + k)}<button class="primary" type="button" data-ussave="${k}">${t('usSave')}</button></div>
          </div></details>
      </div></section>`;
  };
  const rows = d.days.map(r => `<tr><td>${esc(r.day)}</td><td>${USG_NAMES[r.provider] || esc(r.provider)}</td><td style="text-align:right">${usgNum(r.requests)}</td>
    <td style="text-align:right">${r.failed ? `<b style="color:var(--bad,#c62828)">${usgNum(r.failed)}</b>` : '0'}</td><td style="text-align:right">${usgSmall(r.cost_usd)}</td></tr>`).join('');
  $('#app').innerHTML = shell(`${pageHead(t('i_integrations'), t('usDesc'))}${intTabs('usage')}
    <section class="alerts" aria-label="${t('usAlerts')}">${alerts}</section>
    ${d.providers.map(card).join('')}
    <section class="panel"><div class="ph"><h2>${t('usDays')}</h2></div>
      ${rows ? `<div class="tbl-wrap"><table style="min-width:0"><thead><tr><th>${t('usDay')}</th><th>${t('usProv')}</th><th style="text-align:right">${t('usReq')}</th><th style="text-align:right">${t('usFail')}</th><th style="text-align:right">${t('usCost')}</th></tr></thead><tbody>${rows}</tbody></table></div>`
        : `<p class="muted" style="padding:16px 18px">${t('usNoData')}</p>`}
      <p class="muted small" style="padding:0 18px 16px">${t('usSmsNote')}</p></section>`);
  document.querySelectorAll('[data-ussave]').forEach(btn => btn.onclick = async () => {
    const k = btn.dataset.ussave, f = $('#us_' + k), v = n => { const el = f.querySelector(`[name=${n}]`); return el ? el.value.trim() : undefined; };
    f.querySelectorAll('.invalid').forEach(el => el.classList.remove('invalid'));
    const b = USG.data.providers.find(p => p.provider === k).budget;
    const body = {provider: k, warn_below_usd: v('warn_below_usd') ?? b.warn_below_usd, price_in: v('price_in') ?? b.price_in, price_out: v('price_out') ?? b.price_out,
      price_per_sms: v('price_per_sms') ?? b.price_per_sms, daily_limit: Number(v('daily_limit') ?? b.daily_limit) || 0, mfa_code: $('#usCode_' + k).value};
    const bal = v('balance_usd');
    if (bal) body.balance_usd = bal;
    if ((f.querySelector('[name=clear_balance]') || {}).checked) body.balance_usd = '';
    try { USG.data = await api('/admin/api/integrations/usage/budget', {method: 'PUT', body}); toast(t('lgSaved')); drawUsage(); loadBadges(true); } catch (e) { aixFail(e, 'us_' + k); }
  });
}
