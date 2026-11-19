'use strict';
// Admin: retur & sengketa. CS memutuskan (refund penuh/sebagian atau tolak); keuangan mengembalikan uang lewat PayWay lalu mencatat referensinya.
addStrings({
  en: {rtHint: 'Sellers have 48 hours to answer. Rejected or unanswered requests land here. Decide fairly from the photos, the product video and both explanations — your note is shown to the buyer and the seller.',
    rs_disputed: 'Disputes', rs_requested: 'Waiting for seller', rs_approved: 'To refund', rs_closed: 'Closed', rtOrder: 'Order', rtStore: 'Store', rtBuyer: 'Buyer', rtReason: 'Reason',
    rtAsked: 'Asked', rtDelivered: 'Delivered', rtSellerDeadline: 'Seller deadline', rtBuyerSays: 'Buyer', rtSellerSays: 'Seller', rtDecision: 'Phsar decision',
    rtRefundAmt: 'Refund amount ($)', rtOrderTotal: 'order total ${t}', rtNeedsPickup: 'Courier collects the item from the buyer', rtNoteLbl: 'Note to buyer and seller (required)',
    rtRefundBtn: 'Approve refund', rtRejectBtn: 'Reject return', rtConfirmReject: 'Reject this return? The buyer and seller are notified.', rtDecided: 'Decision saved',
    rtPayHint: 'Refund {a} in the ABA PayWay dashboard (find the original payment of order {o}), then record the refund reference.', rtRef: 'PayWay refund reference', rtMarkRefunded: 'Mark refunded',
    rtRefunded: 'Recorded as refunded. The seller balance was adjusted.', rtNoAccess: 'Your role can view this but not act on it.', rtPickupChip: 'Pickup needed',
    rtReason_changed_mind: 'Changed mind', rtReason_wrong_item: 'Wrong item', rtReason_damaged: 'Damaged', rtReason_not_as_described: 'Not as shown', rtReason_missing_parts: 'Missing parts', rtReason_other: 'Other',
    rtSt_requested: 'Waiting for seller', rtSt_disputed: 'Dispute', rtSt_approved: 'Approved — to refund', rtSt_rejected: 'Rejected', rtSt_refunded: 'Refunded', rtSt_cancelled: 'Cancelled by buyer',
    rtEsc_seller_timeout: 'Seller did not answer (48 h)', rtEsc_seller_rejected: 'Seller rejected', err_return_state: 'Already handled — refreshing.'},
  km: {rtHint: 'អ្នកលក់មានពេល 48 ម៉ោងដើម្បីឆ្លើយ។ សំណើដែលត្រូវបានបដិសេធ ឬមិនបានឆ្លើយ មកដល់ទីនេះ។ សូមសម្រេចដោយយុត្តិធម៌ពីរូបថត វីដេអូផលិតផល និងការពន្យល់ទាំងសងខាង — កំណត់ចំណាំរបស់អ្នកបង្ហាញដល់អ្នកទិញ និងអ្នកលក់។',
    rs_disputed: 'វិវាទ', rs_requested: 'រង់ចាំអ្នកលក់', rs_approved: 'ត្រូវសងប្រាក់', rs_closed: 'បានបិទ', rtOrder: 'ការបញ្ជាទិញ', rtStore: 'ហាង', rtBuyer: 'អ្នកទិញ', rtReason: 'មូលហេតុ',
    rtAsked: 'ស្នើ', rtDelivered: 'បានដឹកដល់', rtSellerDeadline: 'ផុតកំណត់អ្នកលក់', rtBuyerSays: 'អ្នកទិញ', rtSellerSays: 'អ្នកលក់', rtDecision: 'ការសម្រេចរបស់ Phsar',
    rtRefundAmt: 'ចំនួនសងប្រាក់ ($)', rtOrderTotal: 'សរុបការបញ្ជាទិញ ${t}', rtNeedsPickup: 'អ្នកដឹកយកទំនិញពីអ្នកទិញ', rtNoteLbl: 'កំណត់ចំណាំទៅអ្នកទិញ និងអ្នកលក់ (ចាំបាច់)',
    rtRefundBtn: 'យល់ព្រមសងប្រាក់', rtRejectBtn: 'បដិសេធការប្ដូរ', rtConfirmReject: 'បដិសេធការប្ដូរនេះ? អ្នកទិញ និងអ្នកលក់នឹងទទួលដំណឹង។', rtDecided: 'បានរក្សាទុកការសម្រេច',
    rtPayHint: 'សងប្រាក់ {a} ក្នុងផ្ទាំង ABA PayWay (រកការបង់ប្រាក់ដើមនៃការបញ្ជាទិញ {o}) រួចកត់លេខយោង។', rtRef: 'លេខយោងការសង PayWay', rtMarkRefunded: 'សម្គាល់ថាបានសង',
    rtRefunded: 'បានកត់ត្រាថាបានសង។ សមតុល្យអ្នកលក់ត្រូវបានកែតម្រូវ។', rtNoAccess: 'តួនាទីរបស់អ្នកអាចមើលបាន ប៉ុន្តែមិនអាចធ្វើសកម្មភាពបានទេ។', rtPickupChip: 'ត្រូវមកយក',
    rtReason_changed_mind: 'ប្ដូរចិត្ត', rtReason_wrong_item: 'ទំនិញខុស', rtReason_damaged: 'ខូច', rtReason_not_as_described: 'មិនដូចបង្ហាញ', rtReason_missing_parts: 'ខ្វះគ្រឿង', rtReason_other: 'ផ្សេងៗ',
    rtSt_requested: 'រង់ចាំអ្នកលក់', rtSt_disputed: 'វិវាទ', rtSt_approved: 'បានយល់ព្រម — ត្រូវសង', rtSt_rejected: 'បានបដិសេធ', rtSt_refunded: 'បានសង', rtSt_cancelled: 'អ្នកទិញបានបោះបង់',
    rtEsc_seller_timeout: 'អ្នកលក់មិនបានឆ្លើយ (48 ម៉ោង)', rtEsc_seller_rejected: 'អ្នកលក់បានបដិសេធ', err_return_state: 'បានដោះស្រាយរួចហើយ — កំពុងផ្ទុកឡើងវិញ។'},
  zh: {rtHint: '卖家有 48 小时回复。被拒绝或未回复的申请会到这里。请根据照片、商品视频和双方说明公正裁决——您的备注会显示给买家和卖家。',
    rs_disputed: '纠纷', rs_requested: '等待卖家', rs_approved: '待退款', rs_closed: '已关闭', rtOrder: '订单', rtStore: '店铺', rtBuyer: '买家', rtReason: '原因',
    rtAsked: '申请', rtDelivered: '送达', rtSellerDeadline: '卖家截止', rtBuyerSays: '买家', rtSellerSays: '卖家', rtDecision: 'Phsar 决定',
    rtRefundAmt: '退款金额（$）', rtOrderTotal: '订单总额 ${t}', rtNeedsPickup: '快递员从买家处取回商品', rtNoteLbl: '给买家和卖家的说明（必填）',
    rtRefundBtn: '同意退款', rtRejectBtn: '拒绝退货', rtConfirmReject: '拒绝此退货？买家和卖家都会收到通知。', rtDecided: '决定已保存',
    rtPayHint: '在 ABA PayWay 后台退款 {a}（查找订单 {o} 的原始付款），然后记录退款参考号。', rtRef: 'PayWay 退款参考号', rtMarkRefunded: '标记为已退款',
    rtRefunded: '已记录退款，卖家余额已调整。', rtNoAccess: '您的角色可以查看，但不能操作。', rtPickupChip: '需取件',
    rtReason_changed_mind: '不想要了', rtReason_wrong_item: '发错商品', rtReason_damaged: '损坏', rtReason_not_as_described: '与展示不符', rtReason_missing_parts: '缺少配件', rtReason_other: '其他',
    rtSt_requested: '等待卖家', rtSt_disputed: '纠纷', rtSt_approved: '已同意——待退款', rtSt_rejected: '已拒绝', rtSt_refunded: '已退款', rtSt_cancelled: '买家已取消',
    rtEsc_seller_timeout: '卖家未回复（48 小时）', rtEsc_seller_rejected: '卖家已拒绝', err_return_state: '已处理——正在刷新。'},
});

const ARt = {status: 'disputed'};
const ART_CHIP = {requested: 'warn', disputed: 'bad', approved: 'warn', refunded: 'ok', rejected: '', cancelled: ''};

async function renderReturns() {
  const q = new URLSearchParams(location.hash.split('?')[1] || '');
  if (q.get('s')) { ARt.status = q.get('s'); history.replaceState(null, '', '#/returns'); }
  const {returns, counts} = await api('/admin/api/returns?status=' + ARt.status);
  const canDecide = hasPerm('return.decide'), canRefund = hasPerm('refund.decide');
  const when = d => d ? new Date(d).toLocaleString() : '—';
  const cards = returns.map(r => {
    const items = r.items.filter(it => it.return_qty > 0).map(it => `<li>${esc(it.name)}${it.options ? ` <span class="muted">${esc(it.options)}</span>` : ''} × ${it.return_qty} · $${esc(it.unit_price)}</li>`).join('');
    const open = r.status === 'requested' || r.status === 'disputed';
    let act = '';
    if (open && canDecide) act = `<div class="rt-act"><div class="grid">
        <label>${t('rtRefundAmt')}<input id="amt_${r.id}" name="amount" inputmode="decimal" value="${esc(r.requested_amount)}"><small class="muted">${esc(t('rtOrderTotal').replace('{t}', r.order_total))}</small></label>
        <label class="inline" style="align-self:center"><input type="checkbox" id="pk_${r.id}" ${r.reason === 'changed_mind' || r.reason === 'wrong_item' ? 'checked' : ''}> ${t('rtNeedsPickup')}</label></div>
      <label>${t('rtNoteLbl')}<textarea id="dn_${r.id}" name="note" maxlength="1000"></textarea></label>
      <div class="row"><button type="button" class="primary" data-dec="refund" data-id="${r.id}">${t('rtRefundBtn')}</button><button type="button" class="danger" data-dec="reject" data-id="${r.id}">${t('rtRejectBtn')}</button></div></div>`;
    else if (r.status === 'approved' && canRefund) act = `<div class="rt-act"><div class="banner info">${esc(t('rtPayHint', {a: '$' + r.refund_amount, o: r.order_no}))}</div>
      <div class="row"><input id="ref_${r.id}" name="reference" placeholder="${t('rtRef')}" maxlength="80" style="flex:1;min-width:200px">
      <button type="button" class="primary" data-refunded="${r.id}">${t('rtMarkRefunded')}</button></div></div>`;
    else if (open || r.status === 'approved') act = `<p class="muted small">${t('rtNoAccess')}</p>`;
    return `<section class="card" id="rt_${r.id}"><div class="spread"><div><h2>${esc(r.return_no)}</h2><span class="muted">${t('rtOrder')} ${esc(r.order_no)} · ${t('rtStore')} ${esc(r.store_name)} · ${t('rtBuyer')} ${esc(r.buyer_name || '—')}${r.buyer_phone ? ' · ' + esc(r.buyer_phone) : ''}</span></div>
        <div style="text-align:right"><span class="chip ${ART_CHIP[r.status]}">${t('rtSt_' + r.status)}</span><br><b style="font-size:20px">$${esc(r.refund_amount || r.requested_amount)}</b></div></div>
      <div class="row" style="gap:6px"><span class="chip pink">${t('rtReason_' + r.reason)}</span>${r.escalated_reason ? `<span class="chip">${t('rtEsc_' + r.escalated_reason)}</span>` : ''}${r.needs_pickup && r.status !== 'refunded' ? `<span class="chip warn">${t('rtPickupChip')}</span>` : ''}</div>
      <div class="grid small"><p><span class="muted">${t('rtAsked')}</span><br>${when(r.created_at)}</p><p><span class="muted">${t('rtDelivered')}</span><br>${when(r.delivered_at)}</p>
        <p><span class="muted">${t('rtSellerDeadline')}</span><br>${when(r.seller_deadline)}</p></div>
      <div class="rt-says"><b>${t('rtBuyerSays')}</b><p>${esc(r.description)}</p>${r.photo_count ? `<div class="rphotos" id="aph_${r.id}"></div>` : ''}<ul>${items}</ul></div>
      ${r.seller_note ? `<div class="rt-says"><b>${t('rtSellerSays')}</b><p>${esc(r.seller_note)}</p></div>` : ''}
      ${r.decision_note ? `<div class="rt-says"><b>${t('rtDecision')}</b><p>${esc(r.decision_note)}</p></div>` : ''}
      ${r.refund_ref ? `<p class="muted small">${t('rtRef')}: <b>${esc(r.refund_ref)}</b> · ${when(r.refunded_at)}</p>` : ''}
      ${act}</section>`;
  }).join('');
  const tabs = ['disputed', 'requested', 'approved', 'closed'].map(f => `<button type="button" aria-selected="${ARt.status === f}" data-rs="${f}">${t('rs_' + f)}${counts[f] ? ` <span class="count">${counts[f]}</span>` : ''}</button>`).join('');
  $('#app').innerHTML = shell(`${pageHead(t('i_returns'), t('rtHint'))}
    <div class="tabs count">${tabs}</div>
    ${cards || `<section class="card"><p class="muted">${t('empty')}</p></section>`}`);
  document.querySelectorAll('[data-rs]').forEach(b => b.onclick = () => { ARt.status = b.dataset.rs; render(); });
  returns.forEach(r => { for (let i = 0; i < r.photo_count; i++) adminPhoto(`/admin/api/returns/${r.id}/photos/${i}`, $('#aph_' + r.id)); });
  const failed = (e, id) => { if (e.code === 'return_state') { toast(t('err_return_state'), 'error'); render(); return; } showError(e, {amount: '#amt_' + id, note: '#dn_' + id, reference: '#ref_' + id}); };
  document.querySelectorAll('[data-dec]').forEach(b => b.onclick = async () => {
    const id = b.dataset.id, decision = b.dataset.dec;
    if (decision === 'reject' && !confirm(t('rtConfirmReject'))) return;
    try {
      await api(`/admin/api/returns/${id}/decide`, {method: 'POST', body: {decision, amount: $('#amt_' + id).value.trim(), needs_pickup: $('#pk_' + id).checked, note: $('#dn_' + id).value}});
      toast(t('rtDecided')); loadBadges(true); render();
    } catch (e) { failed(e, id); }
  });
  document.querySelectorAll('[data-refunded]').forEach(b => b.onclick = async () => {
    const id = b.dataset.refunded;
    try { await api(`/admin/api/returns/${id}/refunded`, {method: 'POST', body: {reference: $('#ref_' + id).value}}); toast(t('rtRefunded')); loadBadges(true); render(); }
    catch (e) { failed(e, id); }
  });
}

async function adminPhoto(url, box) {
  if (!box) return;
  try {
    const r = await fetch(url, {headers: {Authorization: 'Bearer ' + P.token}});
    if (!r.ok) return;
    const src = URL.createObjectURL(await r.blob());
    const a = document.createElement('a'); a.href = src; a.target = '_blank'; a.rel = 'noopener';
    const img = document.createElement('img'); img.src = src; img.alt = ''; a.appendChild(img); box.appendChild(a);
  } catch {}
}
