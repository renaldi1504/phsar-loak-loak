'use strict';
// Seller Center: retur — setujui atau tolak dalam 48 jam; kalau tidak, Phsar yang memutuskan.
addStrings({
  en: {tabReturns: 'Returns', rf_open: 'Open', rf_all: 'All', noReturns: 'No return requests. 🎉', rtNo: 'Return {no}', rtOrder: 'Order {no}',
    rtAnswerIn: 'Answer within {t}', rtLate: 'Time is up — Phsar will decide', rtAmount: 'Buyer asks', rtRefund: 'Refund',
    rtAccept: 'Accept return', rtReject: 'Reject', rtAcceptTitle: 'Accept this return?', rtAcceptHint: 'The buyer gets {a} back. This amount is taken from the order earnings.',
    rtKeep: 'Buyer may keep the item (no pickup needed)', rtNote: 'Message to the buyer (optional)', rtRejectTitle: 'Reject this return?',
    rtRejectHint: 'Phsar will review the photos and your reason, then decide. Explain clearly what you checked.', rtRejectNote: 'Your reason (required)',
    rtSendAccept: 'Accept and refund', rtSendReject: 'Reject and send to Phsar', rtDone: 'Saved', rtYouSaid: 'You', rtPhsar: 'Phsar decision', rtPickup: 'A Phsar courier collects the item from the buyer and brings it back to you.',
    rtReason_changed_mind: 'Changed mind', rtReason_wrong_item: 'Wrong item', rtReason_damaged: 'Damaged', rtReason_not_as_described: 'Not as shown', rtReason_missing_parts: 'Missing parts', rtReason_other: 'Other',
    rtSt_requested: 'Needs your answer', rtSt_disputed: 'Phsar reviewing', rtSt_approved: 'Approved — refund pending', rtSt_rejected: 'Rejected by Phsar', rtSt_refunded: 'Refunded', rtSt_cancelled: 'Cancelled by buyer',
    rtEsc_seller_timeout: 'Not answered within 48 hours', rtEsc_seller_rejected: 'You rejected', err_return_deadline: 'The 48-hour limit has passed — Phsar will decide.', err_return_state: 'This return was already handled.'},
  km: {tabReturns: 'ការប្ដូរ', rf_open: 'កំពុងដំណើរការ', rf_all: 'ទាំងអស់', noReturns: 'គ្មានសំណើប្ដូរទំនិញទេ។ 🎉', rtNo: 'ការប្ដូរ {no}', rtOrder: 'ការបញ្ជាទិញ {no}',
    rtAnswerIn: 'ឆ្លើយក្នុង {t}', rtLate: 'ផុតពេលហើយ — Phsar នឹងសម្រេច', rtAmount: 'អ្នកទិញស្នើ', rtRefund: 'សងប្រាក់',
    rtAccept: 'យល់ព្រមប្ដូរ', rtReject: 'បដិសេធ', rtAcceptTitle: 'យល់ព្រមការប្ដូរនេះ?', rtAcceptHint: 'អ្នកទិញនឹងទទួលបាន {a} វិញ។ ចំនួននេះត្រូវកាត់ពីប្រាក់ចំណូលនៃការបញ្ជាទិញ។',
    rtKeep: 'អ្នកទិញអាចរក្សាទំនិញ (មិនចាំបាច់មកយក)', rtNote: 'សារទៅអ្នកទិញ (ស្រេចចិត្ត)', rtRejectTitle: 'បដិសេធការប្ដូរនេះ?',
    rtRejectHint: 'Phsar នឹងពិនិត្យរូបថត និងមូលហេតុរបស់អ្នក រួចសម្រេច។ សូមពន្យល់ឱ្យច្បាស់ពីអ្វីដែលអ្នកបានពិនិត្យ។', rtRejectNote: 'មូលហេតុ (ចាំបាច់)',
    rtSendAccept: 'យល់ព្រម និងសងប្រាក់', rtSendReject: 'បដិសេធ និងបញ្ជូនទៅ Phsar', rtDone: 'បានរក្សាទុក', rtYouSaid: 'អ្នក', rtPhsar: 'ការសម្រេចរបស់ Phsar', rtPickup: 'អ្នកដឹក Phsar នឹងយកទំនិញពីអ្នកទិញ ហើយយកមកឱ្យអ្នកវិញ។',
    rtReason_changed_mind: 'ប្ដូរចិត្ត', rtReason_wrong_item: 'ទំនិញខុស', rtReason_damaged: 'ខូច', rtReason_not_as_described: 'មិនដូចបង្ហាញ', rtReason_missing_parts: 'ខ្វះគ្រឿង', rtReason_other: 'ផ្សេងៗ',
    rtSt_requested: 'ត្រូវការចម្លើយរបស់អ្នក', rtSt_disputed: 'Phsar កំពុងពិនិត្យ', rtSt_approved: 'បានយល់ព្រម — រង់ចាំសងប្រាក់', rtSt_rejected: 'Phsar បានបដិសេធ', rtSt_refunded: 'បានសងប្រាក់', rtSt_cancelled: 'អ្នកទិញបានបោះបង់',
    rtEsc_seller_timeout: 'មិនបានឆ្លើយក្នុង 48 ម៉ោង', rtEsc_seller_rejected: 'អ្នកបានបដិសេធ', err_return_deadline: 'ផុតពេល 48 ម៉ោងហើយ — Phsar នឹងសម្រេច។', err_return_state: 'ការប្ដូរនេះត្រូវបានដោះស្រាយរួចហើយ។'},
  zh: {tabReturns: '退货', rf_open: '处理中', rf_all: '全部', noReturns: '暂无退货申请。🎉', rtNo: '退货 {no}', rtOrder: '订单 {no}',
    rtAnswerIn: '请在 {t} 内回复', rtLate: '已超时——将由 Phsar 决定', rtAmount: '买家申请', rtRefund: '退款',
    rtAccept: '同意退货', rtReject: '拒绝', rtAcceptTitle: '同意此退货？', rtAcceptHint: '买家将获得 {a} 退款，该金额从订单收入中扣除。',
    rtKeep: '买家可保留商品（无需取件）', rtNote: '给买家的留言（可选）', rtRejectTitle: '拒绝此退货？',
    rtRejectHint: 'Phsar 将查看照片和您的理由后作出决定。请清楚说明您检查了什么。', rtRejectNote: '拒绝理由（必填）',
    rtSendAccept: '同意并退款', rtSendReject: '拒绝并提交 Phsar', rtDone: '已保存', rtYouSaid: '您', rtPhsar: 'Phsar 决定', rtPickup: 'Phsar 快递员会从买家处取回商品并送还给您。',
    rtReason_changed_mind: '不想要了', rtReason_wrong_item: '发错商品', rtReason_damaged: '损坏', rtReason_not_as_described: '与展示不符', rtReason_missing_parts: '缺少配件', rtReason_other: '其他',
    rtSt_requested: '待您回复', rtSt_disputed: 'Phsar 审核中', rtSt_approved: '已同意——待退款', rtSt_rejected: 'Phsar 已拒绝', rtSt_refunded: '已退款', rtSt_cancelled: '买家已取消',
    rtEsc_seller_timeout: '48 小时内未回复', rtEsc_seller_rejected: '您已拒绝', err_return_deadline: '已超过 48 小时——将由 Phsar 决定。', err_return_state: '此退货已处理。'},
});

const SR = {filter: 'open', data: null, timer: null, q: '', page: 1};
addStrings({
  en: {rSearchPh: 'Search return no., order no. or buyer'},
  km: {rSearchPh: 'ស្វែងរកលេខការប្រគល់វិញ លេខការបញ្ជាទិញ ឬអ្នកទិញ'},
  zh: {rSearchPh: '搜索退货单号、订单号或买家'},
});
const RT_CHIP = {requested: 'warn', disputed: 'warn', approved: 'ok', refunded: 'ok', rejected: 'bad', cancelled: ''};

async function tabReturns(id) {
  clearInterval(SR.timer);
  SR.data = await api(`/v1/seller/stores/${id}/returns?${qs({filter: SR.filter, q: SR.q, page: SR.page, size: 20})}`);
  const rs = SR.data.returns, pg = SR.data.page;
  if (pg && SR.page > pg.pages) { SR.page = pg.pages; return tabReturns(id); }
  const openN = pg ? pg.total : rs.length;
  // sellerpages: langkah 1-2-3 + tab "Perlu dijawab / Semua" + hitung mundur 48 jam
  const need = rs.filter(r => r.status === 'requested' && leftText(r.seller_deadline)).length;
  const chips = ['open', 'all'].map(f => `<button type="button" role="tab" aria-selected="${SR.filter === f}" data-rf="${f}">${t('rf_' + f)}${f === 'open' && SR.filter === 'open' && !SR.q && openN ? `<span class="count">${openN}</span>` : ''}</button>`).join('');
  const list = rs.map(returnCard).join('');
  const empty = SR.q ? `<section class="card"><p class="muted" style="text-align:center">${t('lsNoMatch', {q: esc(SR.q)})}</p></section>` : SR.filter === 'open'
    ? spEmpty('🎉', t('spRtNoneT'), t('spRtNoneW'), `<button type="button" data-rf="all">${t('spRtSeeAll')}</button>`)
    : spEmpty('📦', t('noReturns').replace(/\s*🎉$/, ''), t('spRtNoneAllW'), `<button type="button" class="primary" onclick="v2Open('home')">${t('spToday')}</button>`);
  $('#tab').innerHTML = `${spSteps([[t('spRt1'), t('spRt1W')], [t('spRt2'), t('spRt2W')], [t('spRt3'), t('spRt3W')]], need ? 0 : -1)}
    <p class="ohelp">⏱ ${t('spRtHelp')}</p>
    <div class="sp-seg" role="tablist" aria-label="${t('tabReturns')}">${chips}</div>
    <div class="listbar">${searchBox('rQ', SR.q, t('rSearchPh'))}</div>
    ${list || empty}${pagerHTML(pg)}`;
  document.querySelectorAll('[data-rf]').forEach(b => b.onclick = () => { SR.filter = b.dataset.rf; SR.page = 1; tabReturns(id); });
  bindSearch('rQ', v => { if (v === SR.q) return; SR.q = v; SR.page = 1; tabReturns(id); });
  bindPager($('#tab'), p => { SR.page = p; tabReturns(id); });
  document.querySelectorAll('[data-ract]').forEach(b => b.onclick = () => returnDialog(id, SR.data.returns.find(r => r.id === b.dataset.id), b.dataset.ract === 'accept'));
  rs.forEach(r => { for (let i = 0; i < r.photo_count; i++) loadPrivatePhoto(`/v1/seller/stores/${id}/returns/${r.id}/photos/${i}`, $('#rph_' + r.id)); });
  SR.timer = setInterval(() => {
    if (!$('#tab') || S.tab !== 'returns') { clearInterval(SR.timer); return; }
    spDueTick();
  }, 30000);
}

async function loadPrivatePhoto(url, box) {
  if (!box) return;
  try {
    const r = await fetch(url, {headers: {Authorization: 'Bearer ' + P.token}});
    if (!r.ok) return;
    const src = URL.createObjectURL(await r.blob());
    const a = document.createElement('a'); a.href = src; a.target = '_blank'; a.rel = 'noopener';
    a.innerHTML = '<img alt="">'; a.firstChild.src = src; box.appendChild(a);
  } catch {}
}

function returnCard(r) {
  const items = r.items.filter(it => it.return_qty > 0).map(it => `<div class="oitem"><img src="${esc(it.image_url || '')}" alt="" onerror="this.style.visibility='hidden'">
    <div><b>${esc(it.name)}</b>${it.options ? `<div class="muted">${esc(it.options)}</div>` : ''}</div><div class="oqty">× ${it.return_qty}<div class="muted">${omoney(it.unit_price)}</div></div></div>`).join('');
  const ask = r.status === 'requested', left = ask ? leftText(r.seller_deadline) : null;
  const btn = (act, label, cls) => `<button type="button" class="${cls}" data-ract="${act}" data-id="${r.id}">${label}</button>`;
  const next = ask ? (left ? '' : 'late') : r.status;
  return `<section class="card order">
    ${ask ? spDue(r.seller_deadline, 48 * 60) : ''}
    <div class="spread"><div><h2>${t('rtNo', {no: esc(r.return_no)})}</h2><span class="muted">${t('rtOrder', {no: esc(r.order_no)})} · ${fmtTime(r.created_at)}</span></div>
      ${ask ? '' : `<span class="chip ${RT_CHIP[r.status]}">${t('rtSt_' + r.status)}</span>`}</div>
    <div class="sp-chips"><span class="chip pink">${t('rtReason_' + r.reason)}</span>${r.escalated_reason ? `<span class="chip">${t('rtEsc_' + r.escalated_reason)}</span>` : ''}</div>
    <blockquote class="sp-quote"><small>${t('spRtBuyerSays')}</small><p>${esc(r.description)}</p></blockquote>
    ${r.photo_count ? `<div class="rphotos" id="rph_${r.id}"></div>` : ''}
    <div class="oitems">${items}</div>
    <div class="getline"><span>${r.refund_amount ? t('rtRefund') : t('rtAmount')}</span><b>${omoney(r.refund_amount || r.requested_amount)}</b></div>
    ${r.status === 'approved' && r.needs_pickup ? `<div class="banner info">🛵 ${t('rtPickup')}</div>` : ''}
    ${r.seller_note ? `<p class="muted" style="margin:0"><b>${t('rtYouSaid')}:</b> ${esc(r.seller_note)}</p>` : ''}
    ${r.decision_note ? `<div class="banner ${r.status === 'rejected' ? 'ok' : 'info'}"><b>${t('rtPhsar')}:</b> ${esc(r.decision_note)}</div>` : ''}
    ${next ? `<p class="ohelp">${t('spRtNext_' + next)}</p>` : ''}
    ${ask && left ? `<div class="sp-acts">${btn('accept', t('rtAccept'), 'primary sp-block')}${btn('reject', t('spRtRejectBtn'), 'textbtn danger')}</div>` : ''}</section>`;
}

function returnDialog(storeId, r, accept) {
  const d = document.createElement('dialog');
  d.className = 'dlg';
  d.innerHTML = `<form method="dialog"><h2>${t(accept ? 'rtAcceptTitle' : 'rtRejectTitle')}</h2>
    <p class="muted">${accept ? t('rtAcceptHint', {a: omoney(r.requested_amount)}) : t('rtRejectHint')}</p>
    ${accept ? `<label class="inline"><input type="checkbox" id="rtKeep"> ${t('rtKeep')}</label>` : ''}
    <label>${t(accept ? 'rtNote' : 'rtRejectNote')}<textarea id="rtNote" maxlength="1000"></textarea></label>
    <div class="sp-acts"><button type="button" class="${accept ? 'primary' : 'danger'} sp-block" id="rtGo">${t(accept ? 'rtSendAccept' : 'rtSendReject')}</button><button value="cancel" class="textbtn">${t('spBack')}</button></div></form>`;
  document.body.appendChild(d);
  d.addEventListener('close', () => d.remove());
  d.showModal();
  d.querySelector('#rtGo').onclick = async () => {
    const note = d.querySelector('#rtNote').value.trim();
    if (!accept && note.length < 5) { d.querySelector('#rtNote').classList.add('invalid'); d.querySelector('#rtNote').focus(); return; }
    try {
      await api(`/v1/seller/stores/${storeId}/returns/${r.id}/respond`, {method: 'POST', body: {accept, note, keep_item: accept && d.querySelector('#rtKeep').checked}});
      d.close(); toast(t('rtDone')); tabReturns(storeId); loadBell();
    } catch (e) { toast(errMsg(e), 'error'); if (e.code === 'return_state' || e.code === 'return_deadline') { d.close(); tabReturns(storeId); } }
  };
}
