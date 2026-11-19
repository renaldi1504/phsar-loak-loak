'use strict';
// courierreturn: admin — antrean "Dikembalikan ke toko" di Orders & shipping: mulai pengembalian dari pengantaran gagal,
// tugaskan kurir pengembalian, selesaikan manual kalau seller tidak mengonfirmasi.
addStrings({
  en: {d_returns: 'Returned to store', q_courier_returns: 'Parcels going back to stores', crRetDesc: 'Failed deliveries going back to the store. The order is cancelled and refunded when the store confirms (or automatically after 24 h).',
    crReturnBtn: 'Return to store', crReturnTitle: 'Return this package to the store?', crReturnHint: 'The courier brings it back to the pickup point. When the store confirms, the order is cancelled and the refund is queued.',
    crCourier: 'Return courier', crSameCourier: '— same courier —', crStartGo: 'Start return', crStartedMsg: 'Return started',
    crs_returning: 'Courier bringing back', crs_awaiting_seller: 'Waiting for the store', crs_received: 'Received by the store', crLate: 'Store did not confirm in time',
    crk_returning: 'On the way back', crk_awaiting: 'Waiting for store', crk_overdue: 'Past 24 h', crf_open: 'In progress', crf_done: 'Finished (30 days)',
    crReason: 'Reason', crAttempts: 'Attempts', crStartedBy: 'Started', crBy_courier: 'by courier', crBy_admin: 'by admin', crReturnedAt: 'At the store',
    crConfirmBy: 'Store confirms by', crFee: 'Return delivery cost', crFeePhsar: 'paid by Phsar (not charged to the seller)', crTotal: 'Refund to buyer',
    crResolve: 'Mark received', crResolveTitle: 'The store has the package?', crResolveHint: 'Use this when the store did not confirm. The order is cancelled and the refund goes to the Refunds queue.',
    crRestock: 'Put the items back in stock (uncheck if damaged)', crNote: 'Note (optional)', crResolveGo: 'Mark received', crResolved: 'Recorded as received',
    crHandover: 'Handover photo', crFailPhoto: 'Failed delivery photo', crResolvedBy_seller: 'Confirmed by the store', crResolvedBy_admin: 'Marked by admin', crResolvedBy_system: 'Automatic (no answer in 24 h)',
    crCond_ok: 'Good — back in stock', crCond_damaged: 'Damaged — not back in stock', crItems: 'Items',
    rc_delivery_failed: 'Delivery failed, returned to store', c_delivery_failed: 'Delivery failed, returned to store', err_max_attempts: 'No delivery attempts left.',
    jb_courierreturn_auto_receive: 'Returned parcels: auto-receive', jb_courierreturn_auto_receive_d: 'Treats parcels returned to a store as received when the store does not confirm within 24 hours, then cancels the order and queues the refund.'},
  km: {d_returns: 'បានប្រគល់ទៅហាងវិញ', q_courier_returns: 'កញ្ចប់កំពុងត្រឡប់ទៅហាង', crRetDesc: 'ការដឹកជញ្ជូនបរាជ័យដែលកំពុងត្រឡប់ទៅហាង។ ការបញ្ជាទិញត្រូវបោះបង់ និងសងប្រាក់នៅពេលហាងបញ្ជាក់ (ឬដោយស្វ័យប្រវត្តិក្រោយ 24 ម៉ោង)។',
    crReturnBtn: 'ប្រគល់ទៅហាងវិញ', crReturnTitle: 'ប្រគល់កញ្ចប់នេះទៅហាងវិញ?', crReturnHint: 'អ្នកដឹកនាំវាត្រឡប់ទៅទីតាំងទទួល។ នៅពេលហាងបញ្ជាក់ ការបញ្ជាទិញត្រូវបោះបង់ ហើយការសងប្រាក់ត្រូវបានរៀបចំ។',
    crCourier: 'អ្នកដឹកត្រឡប់', crSameCourier: '— អ្នកដឹកដដែល —', crStartGo: 'ចាប់ផ្តើមប្រគល់វិញ', crStartedMsg: 'បានចាប់ផ្តើមប្រគល់វិញ',
    crs_returning: 'អ្នកដឹកកំពុងនាំត្រឡប់', crs_awaiting_seller: 'រង់ចាំហាង', crs_received: 'ហាងបានទទួល', crLate: 'ហាងមិនបានបញ្ជាក់ទាន់ពេល',
    crk_returning: 'កំពុងត្រឡប់', crk_awaiting: 'រង់ចាំហាង', crk_overdue: 'លើស 24 ម៉ោង', crf_open: 'កំពុងដំណើរការ', crf_done: 'បានបញ្ចប់ (30 ថ្ងៃ)',
    crReason: 'មូលហេតុ', crAttempts: 'ចំនួនដង', crStartedBy: 'បានចាប់ផ្តើម', crBy_courier: 'ដោយអ្នកដឹក', crBy_admin: 'ដោយអ្នកគ្រប់គ្រង', crReturnedAt: 'ដល់ហាង',
    crConfirmBy: 'ហាងបញ្ជាក់មុន', crFee: 'ថ្លៃដឹកត្រឡប់', crFeePhsar: 'Phsar ចេញ (មិនគិតលើអ្នកលក់)', crTotal: 'សងអ្នកទិញ',
    crResolve: 'សម្គាល់ថាបានទទួល', crResolveTitle: 'ហាងមានកញ្ចប់ហើយ?', crResolveHint: 'ប្រើពេលហាងមិនបានបញ្ជាក់។ ការបញ្ជាទិញត្រូវបោះបង់ ហើយការសងប្រាក់ចូលបញ្ជីសងប្រាក់វិញ។',
    crRestock: 'ដាក់ទំនិញចូលស្តុកវិញ (ដោះចេញបើខូច)', crNote: 'កំណត់ចំណាំ (មិនចាំបាច់)', crResolveGo: 'សម្គាល់ថាបានទទួល', crResolved: 'បានកត់ត្រាថាបានទទួល',
    crHandover: 'រូបថតពេលប្រគល់', crFailPhoto: 'រូបថតពេលដឹកមិនបាន', crResolvedBy_seller: 'ហាងបានបញ្ជាក់', crResolvedBy_admin: 'អ្នកគ្រប់គ្រងបានសម្គាល់', crResolvedBy_system: 'ស្វ័យប្រវត្តិ (មិនឆ្លើយក្នុង 24 ម៉ោង)',
    crCond_ok: 'ល្អ — ចូលស្តុកវិញ', crCond_damaged: 'ខូច — មិនចូលស្តុកវិញ', crItems: 'ទំនិញ',
    rc_delivery_failed: 'ដឹកមិនបាន ប្រគល់ទៅហាងវិញ', c_delivery_failed: 'ដឹកមិនបាន ប្រគល់ទៅហាងវិញ', err_max_attempts: 'អស់ចំនួនដងដឹកហើយ។',
    jb_courierreturn_auto_receive: 'កញ្ចប់ត្រឡប់៖ ទទួលស្វ័យប្រវត្តិ', jb_courierreturn_auto_receive_d: 'ចាត់ទុកកញ្ចប់ដែលបានប្រគល់ទៅហាងថាបានទទួល នៅពេលហាងមិនបញ្ជាក់ក្នុង 24 ម៉ោង រួចបោះបង់ការបញ្ជាទិញ និងរៀបចំសងប្រាក់។'},
  zh: {d_returns: '已退回店铺', q_courier_returns: '正在退回店铺的包裹', crRetDesc: '配送失败、正在退回店铺的包裹。店铺确认后（或 24 小时后自动）订单取消并退款。',
    crReturnBtn: '退回店铺', crReturnTitle: '将此包裹退回店铺？', crReturnHint: '快递员将其送回取件点。店铺确认后，订单取消并进入退款队列。',
    crCourier: '退回快递员', crSameCourier: '— 同一快递员 —', crStartGo: '开始退回', crStartedMsg: '已开始退回',
    crs_returning: '快递员正在送回', crs_awaiting_seller: '等待店铺确认', crs_received: '店铺已收到', crLate: '店铺未及时确认',
    crk_returning: '退回途中', crk_awaiting: '等待店铺', crk_overdue: '超过 24 小时', crf_open: '进行中', crf_done: '已完成（30 天）',
    crReason: '原因', crAttempts: '尝试次数', crStartedBy: '开始', crBy_courier: '由快递员', crBy_admin: '由管理员', crReturnedAt: '到达店铺',
    crConfirmBy: '店铺确认截止', crFee: '退回运费', crFeePhsar: '由 Phsar 承担（不向卖家收取）', crTotal: '退还买家',
    crResolve: '标记为已收到', crResolveTitle: '店铺已拿到包裹？', crResolveHint: '店铺未确认时使用。订单将取消，退款进入退款队列。',
    crRestock: '商品放回库存（损坏请取消勾选）', crNote: '备注（选填）', crResolveGo: '标记为已收到', crResolved: '已记录为已收到',
    crHandover: '交接照片', crFailPhoto: '配送失败照片', crResolvedBy_seller: '店铺已确认', crResolvedBy_admin: '管理员标记', crResolvedBy_system: '自动（24 小时未回复）',
    crCond_ok: '完好——已回库存', crCond_damaged: '损坏——未回库存', crItems: '商品',
    rc_delivery_failed: '配送失败，已退回店铺', c_delivery_failed: '配送失败，已退回店铺', err_max_attempts: '配送次数已用完。',
    jb_courierreturn_auto_receive: '退回包裹：自动签收', jb_courierreturn_auto_receive_d: '店铺 24 小时内未确认时，将退回店铺的包裹视为已收到，然后取消订单并安排退款。'},
});

// Badge menu "Orders & shipping" + tautan dari dasbor.
for (const [, items] of NAV) for (const it of items) if (it[0] === 'orders' && !String(it[3] || '').includes('courier_returns')) it[3] = (it[3] ? it[3] + '+' : '') + 'courier_returns';
QROUTE.courier_returns = ['#/delivery', 'cs'];

const CR = {filter: 'open'};

// Tombol di kartu pengantaran gagal (admin-delivery.js).
function crFailedBtn(d) {
  return `<button type="button" class="small" data-crstart="${d.id}" data-crno="${esc(d.order_no)}">↩️ ${t('crReturnBtn')}</button>`;
}

function crBindFailed(couriers) {
  document.querySelectorAll('[data-crstart]').forEach(b => b.onclick = () => crStartDialog(b.dataset.crstart, b.dataset.crno, couriers));
}

function crStartDialog(deliveryId, no, couriers) {
  const d = document.createElement('dialog');
  d.className = 'dlg';
  d.innerHTML = `<form method="dialog"><h2>${t('crReturnTitle')}</h2><p><b>${esc(no)}</b></p><p class="muted">${t('crReturnHint')}</p>
    <label>${t('crCourier')}<select id="crCourierSel"><option value="">${t('crSameCourier')}</option>${couriers.map(c => `<option value="${c.user_id}">${esc(c.name)}${c.online ? ' ●' : ''}</option>`).join('')}</select></label>
    <div class="row"><button type="button" class="primary" id="crGo">${t('crStartGo')}</button><button value="cancel">${t('cancel')}</button></div></form>`;
  document.body.append(d);
  d.showModal();
  d.addEventListener('close', () => d.remove());
  $('#crGo', d).onclick = async () => {
    try {
      await api('/admin/api/courier-returns', {method: 'POST', body: {delivery_id: deliveryId, courier_id: $('#crCourierSel', d).value}});
      d.remove(); toast(t('crStartedMsg')); SH.at = 0; AD.status = 'returns'; render();
    } catch (e) { showError(e, {courier_id: '#crCourierSel'}); }
  };
}

function crTabs() {
  return `<div class="tabs">${['open', 'failed', 'delivered', 'returns'].map(f => `<button type="button" aria-selected="${AD.status === f}" onclick="AD.status='${f}';render()">${t('d_' + f)}</button>`).join('')}</div>`;
}

// Halaman tab "Dikembalikan ke toko" (dipanggil renderDelivery saat AD.status = 'returns').
async function crAdminReturns() {
  const [{couriers}, data] = await Promise.all([api('/admin/api/couriers'), api('/admin/api/courier-returns?status=' + CR.filter)]);
  const active = couriers.filter(c => c.active), k = data.counts;
  const dt = v => v ? new Date(v).toLocaleString() : '—';
  const cards = data.returns.map(r => {
    const late = r.status === 'awaiting_seller' && r.confirm_by && new Date(r.confirm_by) < new Date();
    const chipCls = r.status === 'received' ? 'ok' : late ? 'bad' : r.status === 'awaiting_seller' ? 'pink' : 'warn';
    const items = r.items.map(it => `<li>${it.qty} × ${esc(it.name)}${it.options ? ` <span class="muted">(${esc(it.options)})</span>` : ''}</li>`).join('');
    const assign = r.status === 'returning' ? `<label class="inline">${t('crCourier')}<select data-crassign="${r.id}">${active.map(c => `<option value="${c.user_id}" ${c.user_id === r.courier_id ? 'selected' : ''}>${esc(c.name)}${c.online ? ' ●' : ''}</option>`).join('')}
      ${active.some(c => c.user_id === r.courier_id) ? '' : `<option value="" selected>${esc(r.courier_name)}</option>`}</select></label>` : '';
    const done = r.status === 'received' ? `<div class="banner ${r.restocked ? 'ok' : 'warn'}">${t('crResolvedBy_' + r.resolved_by)} · ${dt(r.received_at)} · ${t('crCond_' + (r.condition || 'ok'))}${r.seller_note ? ' — ' + esc(r.seller_note) : ''}</div>` : '';
    return `<section class="card"><div class="spread"><h2>${esc(r.order_no)}</h2>${chip(chipCls, late ? t('crLate') : t('crs_' + r.status))}</div>
      <div class="grid"><p><span class="muted">${t('pickup')}</span><br><b>${esc(r.store_name)}</b> · ${esc(r.pickup_label)}<br>${esc(r.pickup_address)}<br>${esc(r.pickup_contact)} ${esc(r.pickup_phone)}</p>
        <p><span class="muted">${t('crCourier')}</span><br><b>${esc(r.courier_name)}</b>${r.courier_phone ? ` · <a href="tel:${esc(r.courier_phone)}">${esc(r.courier_phone)}</a>` : ''}</p></div>
      <p class="muted">${t('crItems')}</p><ul class="plain">${items}</ul>
      <p class="muted">${t('crReason')}: <b>${r.reason ? t('df_' + r.reason) : '—'}</b>${r.note ? ' — ' + esc(r.note) : ''} · ${t('crAttempts')}: ${r.attempts}
        · ${t('crStartedBy')} ${t('crBy_' + r.started_by)} ${dt(r.started_at)}${r.returned_at ? ` · ${t('crReturnedAt')} ${dt(r.returned_at)}` : ''}
        ${r.status === 'awaiting_seller' && r.confirm_by ? ` · ${t('crConfirmBy')} ${dt(r.confirm_by)}` : ''}</p>
      <p class="muted">${t('crFee')}: ${amoney(r.return_fee)} (${t('crFeePhsar')}) · ${t('crTotal')}: <b>${amoney(r.order_total)}</b></p>
      ${done}
      <div class="row">${assign}
        ${r.status !== 'received' ? `<button type="button" class="primary small" data-crresolve="${r.id}" data-crno="${esc(r.order_no)}">${t('crResolve')}</button>` : ''}
        ${r.has_proof ? `<button type="button" class="small" data-crproof="${r.id}" data-which="returned">📷 ${t('crHandover')}</button>` : ''}
        ${r.has_fail_proof ? `<button type="button" class="small" data-crproof="${r.id}" data-which="fail">📷 ${t('crFailPhoto')}</button>` : ''}</div></section>`;
  }).join('');
  $('#app').innerHTML = shell(`${pageHead(t('i_orders'), t('crRetDesc'))}
    <section class="kpis"><div class="kpi"><span>${t('crk_returning')}</span><b>${k.returning}</b><span>${t('crs_returning')}</span></div>
      <div class="kpi"><span>${t('crk_awaiting')}</span><b>${k.awaiting_seller}</b><span>${t('crs_awaiting_seller')}</span></div>
      <div class="kpi"><span>${t('crk_overdue')}</span><b class="${k.overdue ? 'warn' : ''}">${k.overdue}</b><span>${t('crLate')}</span></div></section>
    ${crTabs()}
    <div class="tabs">${['open', 'done'].map(f => `<button type="button" aria-selected="${CR.filter === f}" onclick="CR.filter='${f}';render()">${t('crf_' + f)}</button>`).join('')}</div>
    ${cards || `<section class="card"><p class="muted">${t('empty')}</p></section>`}`);
  document.querySelectorAll('[data-crassign]').forEach(s => s.onchange = async () => {
    if (!s.value) return;
    try { await api(`/admin/api/courier-returns/${s.dataset.crassign}/assign`, {method: 'POST', body: {courier_id: s.value}}); toast(t('saved')); render(); }
    catch (e) { toast(errMsg(e), 'error'); render(); }
  });
  document.querySelectorAll('[data-crresolve]').forEach(b => b.onclick = () => crResolveDialog(b.dataset.crresolve, b.dataset.crno));
  document.querySelectorAll('[data-crproof]').forEach(b => b.onclick = async () => {
    const r = await fetch(`/admin/api/courier-returns/${b.dataset.crproof}/proof?which=${b.dataset.which}`, {headers: {Authorization: 'Bearer ' + P.token}});
    if (!r.ok) { toast(r.statusText, 'error'); return; }
    window.open(URL.createObjectURL(await r.blob()), '_blank', 'noopener');
  });
}

function crResolveDialog(id, no) {
  const d = document.createElement('dialog');
  d.className = 'dlg';
  d.innerHTML = `<form method="dialog"><h2>${t('crResolveTitle')}</h2><p><b>${esc(no)}</b></p><p class="muted">${t('crResolveHint')}</p>
    <label class="inline"><input type="checkbox" id="crRestock" checked> ${t('crRestock')}</label>
    <label>${t('crNote')}<textarea id="crNoteIn" maxlength="300"></textarea></label>
    <div class="row"><button type="button" class="primary" id="crGo">${t('crResolveGo')}</button><button value="cancel">${t('cancel')}</button></div></form>`;
  document.body.append(d);
  d.showModal();
  d.addEventListener('close', () => d.remove());
  $('#crGo', d).onclick = async () => {
    try {
      await api(`/admin/api/courier-returns/${id}/resolve`, {method: 'POST', body: {restock: $('#crRestock', d).checked, note: $('#crNoteIn', d).value}});
      d.remove(); toast(t('crResolved')); SH.at = 0; render();
    } catch (e) { showError(e, {note: '#crNoteIn'}); }
  };
}
