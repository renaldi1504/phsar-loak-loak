'use strict';
// Admin (super admin, askbot): bot pembantu seller menjawab pertanyaan pembeli — aktif/nonaktif, batas per toko, pemakaian.
addStrings({
  en: {i_askbot: 'Seller reply helper', abaDesc: 'Drafts replies to buyer questions for sellers, using only product details and the shop’s own answers (FAQ). Sellers always check before sending.',
    abaOn: 'Reply helper is on for all sellers', abaLimit: 'Drafts per shop per hour', abaLimitHint: 'Protects the AI budget. When reached, sellers see “try again later”.',
    abaSave: 'Save', abaSaved: 'Saved', abaEngine: 'Engine', abaAI: 'Claude AI · {m}', abaRules: 'Simple mode (no ANTHROPIC_API_KEY): replies built from product data and FAQ only',
    abaUsage: 'Usage', abaDay: 'Drafts, last 24 h', abaWeek: 'Drafts, last 7 days', abaClaude: 'By AI (7 days)', abaRulesN: 'By simple mode (7 days)',
    abaErrors: 'AI errors → simple mode (7 days)', abaTokens: 'Tokens in / out (7 days)', abaFAQs: 'Saved answers (FAQ)', abaFAQStores: '{n} shops',
    abaTop: 'Most active shops (7 days)', abaNoUse: 'No drafts yet.'},
  km: {i_askbot: 'ជំនួយការឆ្លើយរបស់អ្នកលក់', abaDesc: 'ព្រាងចម្លើយសំណួររបស់អ្នកទិញឱ្យអ្នកលក់ ដោយប្រើតែព័ត៌មានផលិតផល និងចម្លើយផ្ទាល់របស់ហាង (FAQ)។ អ្នកលក់ពិនិត្យជានិច្ចមុនពេលផ្ញើ។',
    abaOn: 'ជំនួយការឆ្លើយបើកសម្រាប់អ្នកលក់ទាំងអស់', abaLimit: 'ចំនួនព្រាងក្នុងមួយហាងក្នុងមួយម៉ោង', abaLimitHint: 'ការពារថវិកា AI។ ពេលដល់កំណត់ អ្នកលក់ឃើញ “សូមព្យាយាមពេលក្រោយ”។',
    abaSave: 'រក្សាទុក', abaSaved: 'បានរក្សាទុក', abaEngine: 'ម៉ាស៊ីន', abaAI: 'Claude AI · {m}', abaRules: 'របៀបសាមញ្ញ (គ្មាន ANTHROPIC_API_KEY)៖ ចម្លើយបង្កើតពីទិន្នន័យផលិតផល និង FAQ ប៉ុណ្ណោះ',
    abaUsage: 'ការប្រើប្រាស់', abaDay: 'ព្រាង ២៤ ម៉ោងចុងក្រោយ', abaWeek: 'ព្រាង ៧ ថ្ងៃចុងក្រោយ', abaClaude: 'ដោយ AI (៧ ថ្ងៃ)', abaRulesN: 'ដោយរបៀបសាមញ្ញ (៧ ថ្ងៃ)',
    abaErrors: 'កំហុស AI → របៀបសាមញ្ញ (៧ ថ្ងៃ)', abaTokens: 'Token ចូល / ចេញ (៧ ថ្ងៃ)', abaFAQs: 'ចម្លើយដែលបានរក្សាទុក (FAQ)', abaFAQStores: '{n} ហាង',
    abaTop: 'ហាងសកម្មបំផុត (៧ ថ្ងៃ)', abaNoUse: 'មិនទាន់មានការព្រាងទេ។'},
  zh: {i_askbot: '卖家回复助手', abaDesc: '为卖家起草买家问题的回复，只使用商品信息和店铺自己的答案（常见问题）。卖家发送前总会先检查。',
    abaOn: '为所有卖家开启回复助手', abaLimit: '每店每小时起草次数', abaLimitHint: '控制 AI 费用。达到上限后，卖家会看到“请稍后再试”。',
    abaSave: '保存', abaSaved: '已保存', abaEngine: '引擎', abaAI: 'Claude AI · {m}', abaRules: '简易模式（未设置 ANTHROPIC_API_KEY）：仅根据商品数据和常见问题生成回复',
    abaUsage: '使用情况', abaDay: '近 24 小时起草', abaWeek: '近 7 天起草', abaClaude: 'AI 生成（7 天）', abaRulesN: '简易模式（7 天）',
    abaErrors: 'AI 出错 → 简易模式（7 天）', abaTokens: 'Token 输入 / 输出（7 天）', abaFAQs: '已保存的答案（常见问题）', abaFAQStores: '{n} 家店铺',
    abaTop: '最活跃的店铺（7 天）', abaNoUse: '还没有起草记录。'},
});

async function renderAskbot() {
  const d = await api('/admin/api/askbot'), u = d.usage;
  const stat = (k, v) => `<div class="card" style="padding:14px;gap:4px"><span class="muted small">${t(k)}</span><b style="font-size:22px">${esc(v)}</b></div>`;
  $('#app').innerHTML = shell(`${pageHead(t('i_askbot'), t('abaDesc'))}
    <section class="card" style="max-width:720px">
      <label class="inline"><input type="checkbox" id="abaOn" ${d.settings.enabled ? 'checked' : ''}> <b>${t('abaOn')}</b></label>
      <label>${t('abaLimit')}<input id="abaLimit" type="number" min="1" max="1000" step="1" value="${d.settings.hourly_limit}" style="max-width:140px">
        <small class="muted">${t('abaLimitHint')}</small></label>
      <p class="muted"><b>${t('abaEngine')}:</b> ${d.ai ? esc(t('abaAI', {m: d.model})) : t('abaRules')}</p>
      <div data-errors-anchor><button type="button" class="primary" id="abaSave">${t('abaSave')}</button></div>
    </section>
    <section class="card" style="max-width:720px"><h2>${t('abaUsage')}</h2>
      <div style="display:grid;grid-template-columns:repeat(auto-fit,minmax(150px,1fr));gap:10px">
        ${stat('abaDay', u.day)}${stat('abaWeek', u.week)}${stat('abaClaude', u.claude_week)}${stat('abaRulesN', u.rules_week)}
        ${stat('abaErrors', u.errors_week)}${stat('abaTokens', u.tokens_in_week + ' / ' + u.tokens_out_week)}${stat('abaFAQs', u.faqs + ' · ' + t('abaFAQStores', {n: u.stores_with_faqs}))}
      </div>
      <h2>${t('abaTop')}</h2>
      ${u.top_stores.length ? `<table><tbody>${u.top_stores.map(s => `<tr><td><a href="#/store/${s.store_id}">${esc(s.name)}</a></td><td style="text-align:right">${s.count}</td></tr>`).join('')}</tbody></table>` : `<p class="muted">${t('abaNoUse')}</p>`}
    </section>`);
  $('#abaSave').onclick = async () => {
    try {
      await api('/admin/api/askbot', {method: 'PUT', body: {enabled: $('#abaOn').checked, hourly_limit: Math.round(+$('#abaLimit').value)}});
      toast(t('abaSaved')); renderAskbot();
    } catch (e) { showError(e, {hourly_limit: '#abaLimit'}); }
  };
}

// Menu: tambahkan item ke grup "Pertumbuhan" tanpa mengubah admin-shell.js (hanya super admin, izin "*").
(function abaNav() {
  const g = NAV.find(x => x[0] === 'growth');
  if (g && !g[1].some(it => it[0] === 'askbot')) g[1].push(['askbot', '#/askbot', '*']);
  ROUTE_ITEM['#/askbot'] = 'askbot';
})();
