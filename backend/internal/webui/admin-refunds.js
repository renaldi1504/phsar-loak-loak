'use strict';
// Admin: refund manual — kembalikan uang lewat dashboard ABA PayWay, lalu catat nomor referensinya di sini.
addStrings({
  en: {rfTotalSub: 'to return to buyers', navRefunds: 'Refunds', refundsTitle: 'Refunds', rf_pending: 'To refund', rf_refunded: 'Refunded', rfHint: 'Refund in the ABA PayWay dashboard (search the PayWay transaction ID), then record the refund reference here. The buyer is notified.',
    rfOrder: 'Order', rfStore: 'Store', rfBuyer: 'Buyer', rfAmount: 'Amount', rfReason: 'Reason', rfPayway: 'PayWay transaction', rfWhen: 'Cancelled', rfRef: 'Refund reference', rfNote: 'Note (optional)',
    rfMark: 'Mark refunded', rfDone: 'Recorded as refunded', rfPartOf: 'of a ${t} payment', rc_seller_rejected: 'Rejected by store', rc_seller_timeout: 'Store did not accept in time', rc_pack_timeout: 'Store did not pack in time',
    rc_paid_after_expiry_refund: 'Paid after time-out, item sold out', copy: 'Copy', copied: 'Copied'},
  km: {rfTotalSub: 'ត្រូវសងទៅអ្នកទិញ', navRefunds: 'សងប្រាក់', refundsTitle: 'សងប្រាក់វិញ', rf_pending: 'ត្រូវសង', rf_refunded: 'បានសង', rfHint: 'សងប្រាក់ក្នុងផ្ទាំង ABA PayWay (ស្វែងរកលេខប្រតិបត្តិការ PayWay) រួចកត់លេខយោងនៅទីនេះ។ អ្នកទិញនឹងទទួលដំណឹង។',
    rfOrder: 'ការបញ្ជាទិញ', rfStore: 'ហាង', rfBuyer: 'អ្នកទិញ', rfAmount: 'ចំនួន', rfReason: 'មូលហេតុ', rfPayway: 'ប្រតិបត្តិការ PayWay', rfWhen: 'បានបោះបង់', rfRef: 'លេខយោងការសង', rfNote: 'កំណត់ចំណាំ (មិនចាំបាច់)',
    rfMark: 'សម្គាល់ថាបានសង', rfDone: 'បានកត់ត្រាថាបានសង', rfPartOf: 'នៃការបង់ ${t}', rc_seller_rejected: 'ហាងបានបដិសេធ', rc_seller_timeout: 'ហាងមិនបានទទួលទាន់ពេល', rc_pack_timeout: 'ហាងមិនបានវេចខ្ចប់ទាន់ពេល',
    rc_paid_after_expiry_refund: 'បង់ក្រោយផុតពេល ទំនិញអស់', copy: 'ចម្លង', copied: 'បានចម្លង'},
  zh: {rfTotalSub: '需退还给买家', navRefunds: '退款', refundsTitle: '退款', rf_pending: '待退款', rf_refunded: '已退款', rfHint: '在 ABA PayWay 后台退款（搜索 PayWay 交易号），然后在此记录退款参考号。系统会通知买家。',
    rfOrder: '订单', rfStore: '店铺', rfBuyer: '买家', rfAmount: '金额', rfReason: '原因', rfPayway: 'PayWay 交易', rfWhen: '取消时间', rfRef: '退款参考号', rfNote: '备注（选填）',
    rfMark: '标记为已退款', rfDone: '已记录为已退款', rfPartOf: '（总付款 ${t}）', rc_seller_rejected: '店铺拒单', rc_seller_timeout: '店铺未及时接单', rc_pack_timeout: '店铺未及时打包',
    rc_paid_after_expiry_refund: '超时后付款且商品售罄', copy: '复制', copied: '已复制'},
});

const AR = {status: 'pending'};

async function renderRefunds() {
  const {refunds} = await api('/admin/api/refunds?status=' + AR.status);
  const pending = AR.status === 'pending';
  const cards = refunds.map(r => `<section class="card"><div class="spread"><h2>${esc(r.order_no)}</h2><b style="font-size:20px">$${esc(r.amount)}</b></div>
    <div class="grid"><p><span class="muted">${t('rfStore')}</span><br>${esc(r.store_name)}</p><p><span class="muted">${t('rfBuyer')}</span><br>${esc(r.buyer_name)} · ${esc(r.buyer_phone)}</p>
      <p><span class="muted">${t('rfReason')}</span><br>${esc(t('rc_' + r.cancel_reason))}${r.reject_reason ? ` · ${esc(r.reject_reason)}` : ''}</p>
      <p><span class="muted">${t('rfPayway')}</span><br><code>${esc(r.tran_id)}</code> <button type="button" data-copy="${esc(r.tran_id)}">${t('copy')}</button>
        ${r.checkout_total !== r.amount ? `<br><span class="muted small">${esc(t('rfPartOf').replace('{t}', r.checkout_total))}</span>` : ''}${r.payment_ref ? `<br><span class="muted small">APV ${esc(r.payment_ref)}</span>` : ''}</p>
      <p><span class="muted">${t('rfWhen')}</span><br>${r.cancelled_at ? new Date(r.cancelled_at).toLocaleString() : '—'}</p></div>
    ${pending ? `<div class="row"><input id="ref_${r.order_id}" placeholder="${t('rfRef')}" maxlength="80" style="flex:1;min-width:180px">
      <input id="note_${r.order_id}" placeholder="${t('rfNote')}" maxlength="300" style="flex:1;min-width:180px">
      <button type="button" class="primary" data-refund="${r.order_id}">${t('rfMark')}</button></div>`
      : `<p class="muted">${t('rfRef')}: <b>${esc(r.refund_ref)}</b> · ${r.refunded_at ? new Date(r.refunded_at).toLocaleString() : ''}${r.refund_note ? ' · ' + esc(r.refund_note) : ''}</p>`}</section>`).join('');
  const total = refunds.reduce((a, r) => a + Number(r.amount || 0), 0);
  const oldest = pending && refunds.length ? refunds[0].cancelled_at : null;
  $('#app').innerHTML = shell(`${pageHead(t('i_refunds'), t('rfHint'))}
    ${pending ? `<section class="kpis"><div class="kpi"><span>${t('rf_pending')}</span><b class="${refunds.length ? 'warn' : ''}">${refunds.length}</b><span>${oldest ? t('oldest', {t: ago(oldest)}) : t('nothing')}</span></div>
      <div class="kpi"><span>${t('rfAmount')}</span><b>${amoney(total)}</b><span>${t('rfTotalSub')}</span></div></section>` : ''}
    <div class="tabs">${['pending', 'refunded'].map(f => `<button type="button" aria-selected="${AR.status === f}" onclick="AR.status='${f}';render()">${t('rf_' + f)}</button>`).join('')}</div>
    ${cards || `<section class="card"><p class="muted">${t('empty')}</p></section>`}`);
  document.querySelectorAll('[data-copy]').forEach(b => b.onclick = () => { navigator.clipboard?.writeText(b.dataset.copy); toast(t('copied')); });
  document.querySelectorAll('[data-refund]').forEach(b => b.onclick = async () => {
    const id = b.dataset.refund;
    try {
      await api(`/admin/api/refunds/${id}/done`, {method: 'POST', body: {refund_ref: $('#ref_' + id).value, note: $('#note_' + id).value}});
      toast(t('rfDone')); render();
    } catch (e) { showError(e, {refund_ref: '#ref_' + id}); }
  });
}
