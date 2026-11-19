'use strict';
// Seller Center: strike aktif toko (hanya baca). 3 strike aktif → toko ditangguhkan.
addStrings({
  en: {skTitle: 'Your store has {n} active strike(s) of 3', skHint: 'At 3 active strikes your store is suspended. Strikes expire after 90 days. Contact Phsar support if you think a strike is wrong.',
    skUntil: 'expires {d}', sk_counterfeit: 'Counterfeit', sk_prohibited: 'Prohibited item', sk_misleading: 'Misleading listing', sk_scam: 'Scam', sk_offensive: 'Offensive content', sk_other: 'Other'},
  km: {skTitle: 'ហាងរបស់អ្នកមានការព្រមានសកម្ម {n} ក្នុងចំណោម 3', skHint: 'ពេលមានការព្រមានសកម្ម 3 ហាងនឹងត្រូវផ្អាក។ ការព្រមានផុតកំណត់ក្រោយ 90 ថ្ងៃ។ ទាក់ទងផ្នែកជំនួយ Phsar បើអ្នកគិតថាខុស។',
    skUntil: 'ផុតកំណត់ {d}', sk_counterfeit: 'ក្លែងក្លាយ', sk_prohibited: 'ទំនិញហាមឃាត់', sk_misleading: 'ការផ្សាយបំភាន់', sk_scam: 'បោកប្រាស់', sk_offensive: 'មាតិកាមិនសមរម្យ', sk_other: 'ផ្សេងៗ'},
  zh: {skTitle: '您的店铺有 {n} 次有效警告（满 3 次暂停）', skHint: '累计 3 次有效警告，店铺将被暂停。警告 90 天后失效。如认为有误请联系 Phsar 客服。',
    skUntil: '{d} 失效', sk_counterfeit: '假冒', sk_prohibited: '违禁商品', sk_misleading: '误导性商品信息', sk_scam: '诈骗', sk_offensive: '不当内容', sk_other: '其他'},
});

async function sellerStrikes(storeId) {
  let s;
  try { s = await api(`/v1/seller/stores/${storeId}/strikes`); } catch { return; }
  if (!s.active_points || S.storeId !== storeId) return;
  const host = document.getElementById('storeAlerts');
  if (!host || document.getElementById('skBox')) return;
  const box = document.createElement('div');
  box.id = 'skBox';
  box.className = 'banner ' + (s.active_points >= s.suspend_at ? 'bad' : 'warn');
  box.style.flexDirection = 'column'; box.style.alignItems = 'flex-start';
  box.innerHTML = `<b>${t('skTitle', {n: s.active_points})}</b><span>${t('skHint')}</span><ul style="margin:4px 0 0;padding-left:18px">${s.strikes.filter(k => k.active).map(k =>
    `<li><b>${t('sk_' + k.reason)}</b> — ${esc(typeof exStrikeNote === 'function' ? exStrikeNote(k) : k.note)} <span class="muted">(${t('skUntil', {d: fmtTime(k.expires_at)})})</span></li>`).join('')}</ul>`;
  host.append(box);
}
