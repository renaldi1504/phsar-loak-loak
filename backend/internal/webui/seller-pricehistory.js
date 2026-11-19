'use strict';
// pricehistory (Seller Center): tampilkan harga terendah 30 hari di pemilih produk flash sale, supaya seller tidak menebak harga.
addStrings({
  en: {phLow30: 'lowest 30-day price ${a}', phMustBelow: 'flash price must be below ${a}'},
  km: {phLow30: 'តម្លៃទាបបំផុត 30 ថ្ងៃ ${a}', phMustBelow: 'តម្លៃ flash ត្រូវទាបជាង ${a}'},
  zh: {phLow30: '近 30 天最低价 ${a}', phMustBelow: '抢购价须低于 ${a}'},
});

const PHS = {rule: true};

// harga flash yang disarankan: potongan minimal sesi, dan (kalau aturan nyala) di bawah harga terendah 30 hari
function phsSuggest(price, off, low30) {
  let p = Math.floor(parseFloat(price) * (100 - off)) / 100;
  const low = parseFloat(low30);
  if (PHS.rule && !isNaN(low) && p >= low) p = Math.round(low * 100 - 1) / 100;
  return Math.max(p, 0.01).toFixed(2);
}

function phsLowTag(low30) {
  if (!low30) return '';
  const a = parseFloat(low30).toFixed(2);
  return `<div class="muted small ph-low">${t(PHS.rule ? 'phMustBelow' : 'phLow30', {a})}</div>`;
}
