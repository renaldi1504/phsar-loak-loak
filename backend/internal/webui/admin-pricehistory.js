'use strict';
// pricehistory (admin): harga terendah 30 hari di pemilih produk flash sale + sakelar aturan (super admin, izin flash.rules → "*").
addStrings({
  en: {phLow30: 'lowest 30-day price ${a}', phMustBelow: 'must be below ${a}', phRuleTitle: 'Flash price rule',
    phRuleOn: 'Flash price must be below the lowest normal price of the 30 days before the session starts',
    phRuleHint: 'Stops fake discounts (raising the price just before a sale). Recommended: on.', phRuleSaved: 'Rule saved'},
  km: {phLow30: 'តម្លៃទាបបំផុត 30 ថ្ងៃ ${a}', phMustBelow: 'ត្រូវទាបជាង ${a}', phRuleTitle: 'ច្បាប់តម្លៃ flash',
    phRuleOn: 'តម្លៃ flash ត្រូវទាបជាងតម្លៃធម្មតាទាបបំផុតក្នុង 30 ថ្ងៃមុនវគ្គចាប់ផ្តើម',
    phRuleHint: 'ការពារការបញ្ចុះតម្លៃក្លែងក្លាយ (ដំឡើងតម្លៃមុនការលក់)។ ណែនាំ៖ បើក។', phRuleSaved: 'បានរក្សាទុកច្បាប់'},
  zh: {phLow30: '近 30 天最低价 ${a}', phMustBelow: '须低于 ${a}', phRuleTitle: '抢购价规则',
    phRuleOn: '抢购价必须低于场次开始前 30 天内的最低原价',
    phRuleHint: '防止虚假折扣（促销前先涨价）。建议开启。', phRuleSaved: '规则已保存'},
});

const PH = {rule: true, rules: null};

// harga flash yang disarankan: potongan minimal sesi, dan (kalau aturan nyala) di bawah harga terendah 30 hari
function phSuggest(price, off, low30) {
  let p = Math.floor(parseFloat(price) * (100 - off)) / 100;
  const low = parseFloat(low30);
  if (PH.rule && !isNaN(low) && p >= low) p = Math.round(low * 100 - 1) / 100;
  return Math.max(p, 0.01).toFixed(2);
}

function phLowTag(low30) {
  if (!low30) return '';
  const a = parseFloat(low30).toFixed(2);
  return `<div class="muted small ph-low">${t('phLow30', {a})}${PH.rule ? ` · ${t('phMustBelow', {a})}` : ''}</div>`;
}

// kartu sakelar aturan di daftar sesi flash (hanya super admin)
function phRulesHTML() {
  if (!hasPerm('flash.rules')) return '';
  const on = PH.rules ? PH.rules.low30_enabled : true;
  return `<section class="card" id="phRules"><h2>${t('phRuleTitle')}</h2>
    <label class="inline"><input type="checkbox" id="phRuleOn" ${on ? 'checked' : ''} ${PH.rules ? '' : 'disabled'}> ${t('phRuleOn')}</label>
    <p class="muted small" style="margin:0">${t('phRuleHint')}</p></section>`;
}

async function phBindRules() {
  if (!hasPerm('flash.rules') || !$('#phRules')) return;
  if (!PH.rules) {
    try { PH.rules = await api('/admin/api/flash-rules'); } catch (e) { toast(errMsg(e), 'error'); return; }
    const box = $('#phRules'); if (!box) return;
    box.outerHTML = phRulesHTML();
  }
  const cb = $('#phRuleOn'); if (!cb) return;
  cb.onchange = async () => {
    try { PH.rules = await api('/admin/api/flash-rules', {method: 'PUT', body: {low30_enabled: cb.checked}}); PH.rule = PH.rules.low30_enabled; toast(t('phRuleSaved')); }
    catch (e) { cb.checked = !cb.checked; toast(errMsg(e), 'error'); }
  };
}
