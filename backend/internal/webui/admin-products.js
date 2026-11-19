'use strict';
// Admin: antrean moderasi produk + detail (video, foto, terjemahan, varian, keputusan).
addStrings({
  en: {pqDesc: 'Products that need a manual check wait here (new stores, restricted categories, sensitive words, unusual prices, buyer reports). Clean products go live automatically — see Auto-published. Used items must have a condition note and condition photos.', pickProduct: 'Choose a product on the left.', navStores: 'Stores', navProducts: 'Products', pqueue: 'Product moderation', pf_pending_review: 'Waiting', pf_changes_requested: 'Changes requested', pf_approved: 'Approved', pf_rejected: 'Rejected', pf_auto: 'Auto-published', autoDesc: 'These went live automatically because they passed every check. Open a few each day — “Looks good” clears it, or request changes / reject.', looksGood: 'Looks good', system: 'System (auto)', laneAuto: 'Went live automatically', laneManual: 'Waiting for a moderator because:', priceHit: 'Price {p} (limit {l})', rs_new_store: 'New store (first products)', rs_restricted_category: 'Restricted category', rs_keyword: 'Sensitive word', rs_price_floor: 'Price too low for this item', rs_big_discount: 'Unrealistic discount', rs_price_drop: 'Price dropped a lot', rs_store_strikes: 'Store has strikes', rs_changes_requested: 'Resent after changes requested', rs_was_rejected: 'Was rejected before', rs_reported: 'Reported by buyers', rs_auto_publish_off: 'Auto-publish is off',
    product: 'Product', category: 'Category', price: 'Price', flags: 'Flags', translatingNow: 'translating…', backProducts: '← Products',
    media: 'Video & photos', defectPhotos: 'Condition photos', photoOnly: 'Photos only — no video. Check that the photos show the real item from several angles.', condition: 'Condition', conditionNote: 'Condition note', usage: 'Used for', months: 'months',
    warranty: 'Warranty', returns: 'Returns', days: 'days', variants: 'Variants', stock: 'Stock', sku: 'SKU', translations: 'Name & description',
    original: 'original', autoChecks: 'Automatic checks', kwHit: 'Keyword “{k}” ({c})', usedItem: 'Used item — check condition photos', noFlags: 'Nothing flagged.',
    decision: 'Decision', reason: 'Reason', approveP: 'Approve', changesP: 'Request changes', rejectP: 'Reject (can’t be sold)', notePhP: 'Message to the seller',
    needReason: 'Choose a reason and write a message.', history: 'History', editTr: 'Edit', saveTr: 'Save', cancel: 'Cancel', storeLink: 'Store',
    r_defect_photo_missing: 'Condition photos missing', r_unclear_condition: 'Condition unclear', r_prohibited: 'Prohibited item', r_counterfeit: 'Counterfeit', r_wrong_category: 'Wrong category', r_price: 'Price problem', r_other: 'Other',
    c_new: 'New', c_like_new: 'Like new', c_used_good: 'Used – good', c_used_fair: 'Used – fair', c_refurbished: 'Refurbished',
    m_approved: 'Approved', m_changes_requested: 'Changes requested', m_rejected: 'Rejected', m_pending_review: 'Waiting',
    tr_machine: 'AI', tr_needs_review: 'AI — check', tr_reviewed: 'Checked', err_store_not_active: 'Approve the store first.'},
  km: {pqDesc: 'ផលិតផលដែលត្រូវពិនិត្យដោយដៃ រង់ចាំនៅទីនេះ (ហាងថ្មី ប្រភេទដាក់កម្រិត ពាក្យងាយរងគ្រោះ តម្លៃមិនធម្មតា របាយការណ៍អ្នកទិញ)។ ផលិតផលស្អាតបង្ហាញស្វ័យប្រវត្តិ — មើល បង្ហាញស្វ័យប្រវត្តិ។ ទំនិញប្រើរួចត្រូវមានកំណត់ចំណាំ និងរូបថតស្ថានភាព។', pickProduct: 'ជ្រើសផលិតផលនៅខាងឆ្វេង។', navStores: 'ហាង', navProducts: 'ផលិតផល', pqueue: 'ពិនិត្យផលិតផល', pf_pending_review: 'រង់ចាំ', pf_changes_requested: 'ស្នើឱ្យកែ', pf_approved: 'បានអនុម័ត', pf_rejected: 'បានបដិសេធ', pf_auto: 'បង្ហាញស្វ័យប្រវត្តិ', autoDesc: 'ផលិតផលទាំងនេះបង្ហាញស្វ័យប្រវត្តិ ព្រោះឆ្លងកាត់ការពិនិត្យទាំងអស់។ សូមបើកមើលខ្លះរៀងរាល់ថ្ងៃ — «ល្អហើយ» ដើម្បីដកចេញ ឬស្នើកែ / បដិសេធ។', looksGood: 'ល្អហើយ', system: 'ប្រព័ន្ធ (ស្វ័យប្រវត្តិ)', laneAuto: 'បានបង្ហាញស្វ័យប្រវត្តិ', laneManual: 'រង់ចាំអ្នកត្រួតពិនិត្យ ព្រោះ៖', priceHit: 'តម្លៃ {p} (កម្រិត {l})', rs_new_store: 'ហាងថ្មី (ផលិតផលដំបូងៗ)', rs_restricted_category: 'ប្រភេទដាក់កម្រិត', rs_keyword: 'ពាក្យងាយរងគ្រោះ', rs_price_floor: 'តម្លៃទាបពេកសម្រាប់ទំនិញនេះ', rs_big_discount: 'បញ្ចុះតម្លៃមិនសមហេតុផល', rs_price_drop: 'តម្លៃធ្លាក់ខ្លាំង', rs_store_strikes: 'ហាងមានការព្រមាន', rs_changes_requested: 'ផ្ញើឡើងវិញក្រោយស្នើកែ', rs_was_rejected: 'ធ្លាប់ត្រូវបដិសេធ', rs_reported: 'អ្នកទិញរាយការណ៍', rs_auto_publish_off: 'ការបង្ហាញស្វ័យប្រវត្តិបិទ',
    product: 'ផលិតផល', category: 'ប្រភេទ', price: 'តម្លៃ', flags: 'សញ្ញាព្រមាន', translatingNow: 'កំពុងបកប្រែ…', backProducts: '← ផលិតផល',
    media: 'វីដេអូ និងរូបថត', defectPhotos: 'រូបថតស្ថានភាព', photoOnly: 'មានតែរូបថត — គ្មានវីដេអូ។ សូមពិនិត្យថារូបបង្ហាញទំនិញពិតពីមុំច្រើន។', condition: 'ស្ថានភាព', conditionNote: 'កំណត់ចំណាំស្ថានភាព', usage: 'បានប្រើ', months: 'ខែ',
    warranty: 'ការធានា', returns: 'ការប្តូរ', days: 'ថ្ងៃ', variants: 'ជម្រើស', stock: 'ស្តុក', sku: 'SKU', translations: 'ឈ្មោះ និងការពិពណ៌នា',
    original: 'ដើម', autoChecks: 'ការពិនិត្យស្វ័យប្រវត្តិ', kwHit: 'ពាក្យ “{k}” ({c})', usedItem: 'ទំនិញប្រើរួច — ពិនិត្យរូបថតស្ថានភាព', noFlags: 'គ្មានសញ្ញាព្រមាន។',
    decision: 'សេចក្តីសម្រេច', reason: 'មូលហេតុ', approveP: 'អនុម័ត', changesP: 'ស្នើឱ្យកែ', rejectP: 'បដិសេធ (មិនអាចលក់)', notePhP: 'សារទៅអ្នកលក់',
    needReason: 'សូមជ្រើសមូលហេតុ និងសរសេរសារ។', history: 'ប្រវត្តិ', editTr: 'កែ', saveTr: 'រក្សាទុក', cancel: 'បោះបង់', storeLink: 'ហាង',
    r_defect_photo_missing: 'ខ្វះរូបថតស្ថានភាព', r_unclear_condition: 'ស្ថានភាពមិនច្បាស់', r_prohibited: 'ទំនិញហាមឃាត់', r_counterfeit: 'ក្លែងក្លាយ', r_wrong_category: 'ប្រភេទខុស', r_price: 'បញ្ហាតម្លៃ', r_other: 'ផ្សេងៗ',
    c_new: 'ថ្មី', c_like_new: 'ដូចថ្មី', c_used_good: 'ប្រើរួច – ល្អ', c_used_fair: 'ប្រើរួច – មធ្យម', c_refurbished: 'ជួសជុលឡើងវិញ',
    m_approved: 'បានអនុម័ត', m_changes_requested: 'ស្នើឱ្យកែ', m_rejected: 'បានបដិសេធ', m_pending_review: 'រង់ចាំ',
    tr_machine: 'AI', tr_needs_review: 'AI — ពិនិត្យ', tr_reviewed: 'បានពិនិត្យ', err_store_not_active: 'សូមអនុម័តហាងជាមុនសិន។'},
  zh: {pqDesc: '需要人工审核的商品在这里等待（新店铺、受限类目、敏感词、价格异常、买家举报）。无问题的商品会自动上架——见“自动上架”。二手商品必须有成色说明和成色照片。', pickProduct: '请在左侧选择商品。', navStores: '店铺', navProducts: '商品', pqueue: '商品审核', pf_pending_review: '待审核', pf_changes_requested: '已要求修改', pf_approved: '已通过', pf_rejected: '已拒绝', pf_auto: '自动上架', autoDesc: '这些商品通过了所有检查，已自动上架。每天抽查几件——点“没问题”即可移出，或要求修改 / 拒绝。', looksGood: '没问题', system: '系统（自动）', laneAuto: '已自动上架', laneManual: '等待审核，原因：', priceHit: '价格 {p}（限制 {l}）', rs_new_store: '新店铺（前几件商品）', rs_restricted_category: '受限类目', rs_keyword: '敏感词', rs_price_floor: '对此商品而言价格过低', rs_big_discount: '折扣不合理', rs_price_drop: '价格大幅下降', rs_store_strikes: '店铺有违规记录', rs_changes_requested: '要求修改后重新提交', rs_was_rejected: '曾被拒绝', rs_reported: '被买家举报', rs_auto_publish_off: '自动上架已关闭',
    product: '商品', category: '类目', price: '价格', flags: '提示', translatingNow: '翻译中…', backProducts: '← 商品',
    media: '视频和照片', defectPhotos: '成色照片', photoOnly: '仅有照片——无视频。请确认照片从多个角度展示了实物。', condition: '成色', conditionNote: '成色说明', usage: '已使用', months: '个月',
    warranty: '保修', returns: '退换', days: '天', variants: '规格', stock: '库存', sku: 'SKU', translations: '名称和描述',
    original: '原文', autoChecks: '自动检查', kwHit: '关键词“{k}”（{c}）', usedItem: '二手商品——请检查成色照片', noFlags: '没有提示。',
    decision: '审核结果', reason: '原因', approveP: '通过', changesP: '要求修改', rejectP: '拒绝（不可销售）', notePhP: '给卖家的说明',
    needReason: '请选择原因并填写说明。', history: '记录', editTr: '编辑', saveTr: '保存', cancel: '取消', storeLink: '店铺',
    r_defect_photo_missing: '缺少成色照片', r_unclear_condition: '成色不清楚', r_prohibited: '违禁商品', r_counterfeit: '假冒商品', r_wrong_category: '类目错误', r_price: '价格问题', r_other: '其他',
    c_new: '全新', c_like_new: '几乎全新', c_used_good: '二手 – 良好', c_used_fair: '二手 – 一般', c_refurbished: '翻新',
    m_approved: '已通过', m_changes_requested: '已要求修改', m_rejected: '已拒绝', m_pending_review: '待审核',
    tr_machine: 'AI', tr_needs_review: 'AI——需检查', tr_reviewed: '已核对', err_store_not_active: '请先通过店铺审核。'},
});

const AP = {filter: 'pending_review', editTr: null};
const LN = {en: 'English', km: 'ខ្មែរ', zh: '中文'};
const usd = v => v == null ? '—' : '$' + Number(v).toFixed(2);
const modCls = s => ({approved: 'ok', pending_review: 'warn', changes_requested: 'bad', rejected: 'bad'})[s] || '';

async function productList() {
  const {products} = await api('/admin/api/products?status=' + AP.filter);
  AP.list = products;
  return products;
}
function productQueueHTML(products, cur) {
  return products.map(p => `<a class="qcard" href="#/product/${p.id}" aria-current="${p.id === cur}"><span class="t"><b>${esc(p.name)}</b><small>${p.submitted_at ? ago(p.submitted_at) : ''}</small></span>
    <span class="s">${esc(p.store_name)} · ${usd(p.min_price)}</span>
    <span class="c">${chip(p.condition === 'new' ? '' : 'warn', t('c_' + p.condition))}${p.flags ? chip('bad', '⚑ ' + p.flags) : ''}${(p.reasons || []).slice(0, 2).map(r => chip('warn', t('rs_' + r))).join('')}${p.translating ? chip('warn', t('translatingNow')) : ''}</span></a>`).join('') || `<p class="muted">${t('empty')}</p>`;
}
function productHead() {
  return `${pageHead(t('pqueue'), AP.filter === 'auto' ? t('autoDesc') : t('pqDesc'))}
    <div class="tabs count">${['pending_review', 'auto', 'changes_requested', 'approved', 'rejected'].map(f => `<button type="button" aria-selected="${AP.filter === f}" onclick="AP.filter='${f}';location.hash='#/products';render()">${t('pf_' + f)}${f === 'pending_review' && SH.badges.products ? ` <span>· ${SH.badges.products}</span>` : ''}</button>`).join('')}</div>`;
}

async function renderProductQueue() {
  const products = await productList();
  if (products.length && window.innerWidth > 1100) return renderProductDetail(products[0].id, products);
  $('#app').innerHTML = shell(`${productHead()}<div class="split"><div class="qlist">${productQueueHTML(products)}</div>
    <div class="detail"><section class="card"><p class="muted">${products.length ? t('pickProduct') : t('empty')}</p></section></div></div>`);
}

async function renderProductDetail(id, list) {
  const [d] = await Promise.all([api('/admin/api/products/' + id), list ? Promise.resolve(list) : productList()]);
  drawProductDetail(d);
}

function drawProductDetail(d) {
  const p = d.product, st = d.store, src = p.source_lang;
  const video = p.media.find(m => m.kind === 'video'), photos = p.media.filter(m => m.kind === 'photo'), defects = p.media.filter(m => m.kind === 'defect_photo');
  const flags = p.auto_flags || {};
  const flagList = [...(flags.keywords || []).map(k => t('kwHit', {k: k.keyword, c: k.category})),
    ...(flags.price || []).map(h => `${t('rs_' + h.reason)}${h.keyword ? ' · “' + h.keyword + '”' : ''} — ${t('priceHit', {p: usd(h.price), l: usd(h.limit)})}`),
    ...(flags.used_item ? [t('usedItem')] : [])];
  const laneLine = flags.lane === 'manual' ? `<p><b>${t('laneManual')}</b> ${(flags.reasons || []).map(r => chip('warn', t('rs_' + r))).join(' ')}</p>`
    : flags.lane === 'auto' ? `<p>${chip('ok', t('laneAuto'))}</p>` : '';
  const optionLine = p.options.map(o => `<p><b>${esc(o.name[src] || '')}</b>: ${o.values.map(v => `${v.swatch_hex ? `<span class="sw" style="background:${esc(v.swatch_hex)}"></span>` : ''}${esc(v.value[src] || '')}`).join(', ')}
      <span class="muted">· ${['en', 'km', 'zh'].filter(l => l !== src).map(l => `${l}: ${esc(o.name[l] || '…')} (${o.values.map(v => esc(v.value[l] || '…')).join(', ')})`).join(' · ')}</span></p>`).join('');
  const trCards = ['en', 'km', 'zh'].map(l => {
    const tr = p.translations[l];
    const isSrc = l === src;
    if (AP.editTr === l) return `<div class="card trcard"><b>${LN[l]}</b><label>${t('product')}<input id="atrName" value="${esc(tr?.name || '')}" maxlength="160"></label>
      <label>${t('about')}<textarea id="atrDesc">${esc(tr?.description || '')}</textarea></label>
      <div class="row"><button type="button" class="primary" onclick="saveAdminTr('${p.id}','${l}')">${t('saveTr')}</button><button type="button" onclick="AP.editTr=null;drawProductDetail(AP.d)">${t('cancel')}</button></div></div>`;
    return `<div class="card trcard"><div class="spread"><b>${LN[l]}</b>${isSrc ? `<span class="chip pink">${t('original')}</span>` : tr ? `<span class="chip ${tr.review_status === 'needs_review' ? 'warn' : ''}">${t('tr_' + tr.review_status)}${tr.model ? ' · ' + esc(tr.model) : ''}</span>` : `<span class="chip warn">${t('translatingNow')}</span>`}</div>
      ${tr ? `<p><b>${esc(tr.name)}</b></p><p class="muted pre">${esc(tr.description)}</p>` : ''}
      ${!isSrc && tr ? `<div><button type="button" onclick="AP.editTr='${l}';drawProductDetail(AP.d)">${t('editTr')}</button></div>` : ''}</div>`;
  }).join('');
  AP.d = d;
  $('#app').innerHTML = shell(`${productHead()}<div class="split"><div class="qlist">${productQueueHTML(AP.list || [], p.id)}</div><div class="detail">
    <div class="spread"><div><h2 style="font-size:22px">${esc(p.translations[src]?.name || '')}</h2>
      <p class="muted">${t('storeLink')}: <a href="#/store/${st.id}">${esc(st.name)}</a> · ${esc(p.category_id)} · ${t('c_' + p.condition)}</p></div>
      <span class="chip ${modCls(p.moderation_status)}">${t('m_' + p.moderation_status)}</span></div>
    ${p.moderation_note ? `<div class="banner bad">${esc(p.moderation_note)}</div>` : ''}
    <section class="card"><h2>${t('autoChecks')}</h2>${laneLine}${flagList.length ? flagList.map(f => `<div class="banner warn">${esc(f)}</div>`).join('') : `<p class="muted">${t('noFlags')}</p>`}</section>
    <section class="card"><h2>${t('media')}</h2>${video ? '' : `<div class="banner info">${t('photoOnly')}</div>`}
      <div class="row" style="align-items:flex-start">${video ? `<video class="vprev" src="${esc(video.url)}" ${video.poster_url ? `poster="${esc(video.poster_url)}"` : ''} controls preload="metadata"></video>` : ''}
        <div class="thumbs">${photos.map(m => `<a href="${esc(m.url)}" target="_blank" rel="noopener" class="thumb"><img src="${esc(m.url)}" alt=""></a>`).join('')}</div></div>
      ${defects.length ? `<b>${t('defectPhotos')}</b><div class="thumbs">${defects.map(m => `<a href="${esc(m.url)}" target="_blank" rel="noopener" class="thumb"><img src="${esc(m.url)}" alt=""></a>`).join('')}</div>` : ''}
    </section>
    ${p.condition !== 'new' ? `<section class="card"><h2>${t('condition')}: ${t('c_' + p.condition)}</h2><p class="pre">${esc(p.condition_note)}</p>${p.usage_months != null ? `<p class="muted">${t('usage')} ${p.usage_months} ${t('months')}</p>` : ''}</section>` : ''}
    <section class="card"><h2>${t('translations')}</h2><div class="grid">${trCards}</div></section>
    <section class="card"><h2>${t('variants')}</h2>${optionLine}
      <div style="overflow-x:auto"><table><thead><tr><th>${t('variants')}</th><th>${t('price')}</th><th>${t('stock')}</th><th>${t('sku')}</th></tr></thead><tbody>
      ${p.variants.map(v => `<tr><td>${esc(v.options.join(' / ') || '—')}</td><td>${usd(v.price)}${v.compare_at_price ? ` <s class="muted">${usd(v.compare_at_price)}</s>` : ''}</td><td>${v.stock_qty}</td><td class="muted">${esc(v.sku)}</td></tr>`).join('')}
      </tbody></table></div>
      <p class="muted">${t('warranty')}: ${p.warranty_type === 'none' ? '—' : `${esc(p.warranty_type)} ${p.warranty_months} ${t('months')}`} · ${t('returns')}: ${p.return_days} ${t('days')}</p></section>
    ${p.status !== 'draft' ? `<section class="card"><h2>${t('decision')}</h2>
      <div class="grid"><label>${t('reason')}<select id="pReason"><option value="">—</option>${['defect_photo_missing', 'unclear_condition', 'prohibited', 'counterfeit', 'wrong_category', 'price', 'other'].map(r => `<option value="${r}">${t('r_' + r)}</option>`).join('')}</select></label>
        <label>${t('notePhP')}<input id="pNote" maxlength="500"></label></div>
      <div class="row"><button type="button" class="primary" data-pd="approve">${p.moderation_status === 'approved' && p.post_review_pending ? t('looksGood') : t('approveP')}</button>
        <button type="button" data-pd="changes_requested">${t('changesP')}</button>
        <button type="button" class="danger" data-pd="reject">${t('rejectP')}</button></div></section>` : ''}
    ${d.history.length ? `<section class="card"><h2>${t('history')}</h2>${d.history.map(h => `<p><span class="chip ${modCls(h.decision)}">${t('m_' + h.decision)}</span>
      ${h.reason_code ? t('r_' + h.reason_code) + ' · ' : ''}${esc(h.note || '')} <span class="muted">— ${esc(h.reviewer || t('system'))}, ${new Date(h.created_at).toLocaleString()}</span></p>`).join('')}</section>` : ''}
    </div></div>`);
  document.querySelectorAll('[data-pd]').forEach(b => b.onclick = async () => {
    const decision = b.dataset.pd, reason_code = $('#pReason').value, note = $('#pNote').value.trim();
    if (decision !== 'approve' && (!reason_code || !note)) { toast(t('needReason'), 'warning'); return; }
    try { const r = await api(`/admin/api/products/${p.id}/decision`, {method: 'POST', body: {decision, reason_code, note}}); SH.at = 0; await loadBadges(); await productList(); drawProductDetail(r); toast(t('done')); }
    catch (e) { toast(errMsg(e), 'error'); }
  });
}

async function saveAdminTr(id, lang) {
  try {
    const d = await api(`/admin/api/products/${id}/translations/${lang}`, {method: 'PUT', body: {name: $('#atrName').value, description: $('#atrDesc').value}});
    AP.editTr = null; drawProductDetail(d); toast(t('done'));
  } catch (e) { toast(errMsg(e), 'error'); }
}
