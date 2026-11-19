'use strict';
// Admin keuangan: jadwal payout, buat payout, transfer manual lalu catat referensi.
addStrings({
  en: {poNextShort: 'Next automatic run', poHoldSub: 'days after delivery', poMinSub: 'smaller balances wait', poBreakdown: 'Per-store payouts', poOrdersCol: 'Orders', poAction: 'Transfer, then record', poFoot: 'Full account numbers are shown only for payouts waiting to be transferred · viewing them is written to the audit log', navPayouts: 'Payouts', poTitle: 'Seller payouts', poHint: 'Transfer each payout from the Phsar bank account to the seller, then record the transfer reference. The seller is notified.',
    poSettings: 'Schedule', poFreq: 'Frequency', poDaily: 'Every day', poWeekly: 'Every week', poWeekday: 'Day', poHour: 'Hour (Phnom Penh)', poHold: 'Hold after delivery (days)',
    poMin: 'Minimum payout ($)', poSave: 'Save schedule', poNext: 'Next automatic run: {d}', poGenerate: 'Create payouts now', poCreated: 'Created {n} payouts · ${t}',
    po_pending: 'To transfer', po_paid: 'Paid', po_cancelled: 'Cancelled', poAccount: 'Account', poRef: 'Bank transfer reference', poNote: 'Note (optional)', poMarkPaid: 'Mark paid',
    poCancel: 'Cancel', poCancelNote: 'Reason for cancelling', poDone: 'Saved', poSkipped: 'Not paid this time', sk_no_account: 'no bank account', sk_account_hold: 'account changed < 48 h',
    sk_below_min: 'below minimum', sk_pending_exists: 'previous payout not transferred yet', poOrders: '{n} orders', copy: 'Copy', copied: 'Copied',
    wd0: 'Sunday', wd1: 'Monday', wd2: 'Tuesday', wd3: 'Wednesday', wd4: 'Thursday', wd5: 'Friday', wd6: 'Saturday', bank_aba: 'ABA Bank', bank_acleda: 'ACLEDA', bank_wing: 'Wing', bank_bakong: 'Bakong', bank_other: 'Other'},
  km: {poNextShort: 'ដំណើរការស្វ័យប្រវត្តិបន្ទាប់', poHoldSub: 'ថ្ងៃបន្ទាប់ពីដឹកដល់', poMinSub: 'សមតុល្យតិចជាងនេះរង់ចាំ', poBreakdown: 'ការទូទាត់តាមហាង', poOrdersCol: 'ការបញ្ជាទិញ', poAction: 'ផ្ទេរ រួចកត់ត្រា', poFoot: 'លេខគណនីពេញបង្ហាញតែសម្រាប់ការទូទាត់ដែលរង់ចាំផ្ទេរ · ការមើលត្រូវកត់ក្នុងកំណត់ហេតុសវនកម្ម', navPayouts: 'ទូទាត់អ្នកលក់', poTitle: 'ការទូទាត់ឱ្យអ្នកលក់', poHint: 'ផ្ទេរប្រាក់ពីគណនី Phsar ទៅអ្នកលក់ រួចកត់លេខយោងការផ្ទេរ។ អ្នកលក់នឹងទទួលដំណឹង។',
    poSettings: 'កាលវិភាគ', poFreq: 'ភាពញឹកញាប់', poDaily: 'រៀងរាល់ថ្ងៃ', poWeekly: 'រៀងរាល់សប្តាហ៍', poWeekday: 'ថ្ងៃ', poHour: 'ម៉ោង (ភ្នំពេញ)', poHold: 'រង់ចាំក្រោយដឹកដល់ (ថ្ងៃ)',
    poMin: 'ទូទាត់អប្បបរមា ($)', poSave: 'រក្សាទុកកាលវិភាគ', poNext: 'ដំណើរការស្វ័យប្រវត្តិបន្ទាប់៖ {d}', poGenerate: 'បង្កើតការទូទាត់ឥឡូវ', poCreated: 'បានបង្កើត {n} · ${t}',
    po_pending: 'ត្រូវផ្ទេរ', po_paid: 'បានទូទាត់', po_cancelled: 'បានបោះបង់', poAccount: 'គណនី', poRef: 'លេខយោងការផ្ទេរ', poNote: 'កំណត់ចំណាំ (មិនចាំបាច់)', poMarkPaid: 'សម្គាល់ថាបានទូទាត់',
    poCancel: 'បោះបង់', poCancelNote: 'មូលហេតុបោះបង់', poDone: 'បានរក្សាទុក', poSkipped: 'មិនបានទូទាត់លើកនេះ', sk_no_account: 'គ្មានគណនីធនាគារ', sk_account_hold: 'ទើបប្តូរគណនី < 48 ម៉ោង',
    sk_below_min: 'តិចជាងអប្បបរមា', sk_pending_exists: 'ការទូទាត់មុនមិនទាន់ផ្ទេរ', poOrders: 'ការបញ្ជាទិញ {n}', copy: 'ចម្លង', copied: 'បានចម្លង',
    wd0: 'អាទិត្យ', wd1: 'ចន្ទ', wd2: 'អង្គារ', wd3: 'ពុធ', wd4: 'ព្រហស្បតិ៍', wd5: 'សុក្រ', wd6: 'សៅរ៍', bank_aba: 'ABA Bank', bank_acleda: 'ACLEDA', bank_wing: 'Wing', bank_bakong: 'Bakong', bank_other: 'ផ្សេងៗ'},
  zh: {poNextShort: '下次自动结算', poHoldSub: '送达后天数', poMinSub: '低于此金额继续累计', poBreakdown: '按店铺结算', poOrdersCol: '订单', poAction: '转账后记录', poFoot: '完整账号仅在待转账时显示 · 查看会记录到审计日志', navPayouts: '卖家结算', poTitle: '卖家结算', poHint: '从 Phsar 银行账户向卖家转账，然后记录转账参考号。系统会通知卖家。',
    poSettings: '结算计划', poFreq: '频率', poDaily: '每天', poWeekly: '每周', poWeekday: '星期', poHour: '时间（金边）', poHold: '送达后冻结天数',
    poMin: '最低结算金额 ($)', poSave: '保存计划', poNext: '下次自动结算：{d}', poGenerate: '立即生成结算', poCreated: '已生成 {n} 笔 · ${t}',
    po_pending: '待转账', po_paid: '已转账', po_cancelled: '已取消', poAccount: '账户', poRef: '银行转账参考号', poNote: '备注（选填）', poMarkPaid: '标记为已转账',
    poCancel: '取消', poCancelNote: '取消原因', poDone: '已保存', poSkipped: '本次未结算', sk_no_account: '未填写银行账户', sk_account_hold: '48 小时内更换过账户',
    sk_below_min: '低于最低金额', sk_pending_exists: '上一笔尚未转账', poOrders: '{n} 个订单', copy: '复制', copied: '已复制',
    wd0: '周日', wd1: '周一', wd2: '周二', wd3: '周三', wd4: '周四', wd5: '周五', wd6: '周六', bank_aba: 'ABA Bank', bank_acleda: 'ACLEDA', bank_wing: 'Wing', bank_bakong: 'Bakong', bank_other: '其他'},
});

const APO = {status: 'pending', last: null};

async function renderPayouts() {
  const [{settings: st, next_run}, {payouts}] = await Promise.all([api('/admin/api/payout-settings'), api('/admin/api/payouts?status=' + APO.status)]);
  const pending = APO.status === 'pending';
  const last = APO.last;
  const total = payouts.reduce((a, p) => a + Number(p.amount || 0), 0);
  const rows = payouts.map(p => `<tr><td><b>${esc(p.store_name)}</b>${p.requested_by_seller ? ' ' + chip('warn', t('exReqBySeller')) : ''}<br><span class="muted small">${new Date(p.created_at).toLocaleString()}</span></td>
    <td class="num">${p.order_count}</td><td class="num"><b>$${esc(p.amount)}</b></td>
    <td>${esc(t('bank_' + p.bank))} · ${esc(p.account_name)}<br>${pending ? `<code>${esc(p.account_number)}</code> <button type="button" class="small" data-copy="${esc(p.account_number)}">${t('copy')}</button>` : '<span class="muted">••' + esc(p.account_last4) + '</span>'}</td>
    <td>${pending ? `<div class="row" style="flex-wrap:nowrap"><input id="ref_${p.id}" placeholder="${t('poRef')}" maxlength="80" style="min-width:150px">
        <button type="button" class="primary" data-paid="${p.id}">${t('poMarkPaid')}</button><button type="button" class="danger" data-cancel="${p.id}" aria-label="${t('poCancel')}">✕</button></div>
        <input id="note_${p.id}" placeholder="${t('poNote')}" maxlength="300" style="margin-top:6px">`
      : `${chip(p.status === 'paid' ? 'ok' : '', t('po_' + p.status))}<br><span class="muted small">${p.reference ? esc(p.reference) + ' · ' : ''}${p.paid_at ? new Date(p.paid_at).toLocaleString() : ''}${p.note ? ' · ' + esc(p.note) : ''}</span>`}</td></tr>`).join('');
  $('#app').innerHTML = shell(`${pageHead(t('i_payouts'), t('poHint'), `<button type="button" class="primary" id="gen">${t('poGenerate')}</button>`)}
    ${last ? `<div class="banner ok"><strong>${t('poCreated', {n: last.created, t: last.total})}</strong>${last.skipped.length ? `<span>${t('poSkipped')}: ${last.skipped.map(s => `${esc(s.store_name)} ($${esc(s.balance)}, ${t('sk_' + s.reason)})`).join(' · ')}</span>` : ''}</div>` : ''}
    <section class="kpis">
      <div class="kpi"><span>${t('k_transfer')}</span><b>${amoney(pending ? total : (SH.dash?.kpis.to_transfer || 0))}</b><span>${t('k_transferSub', {n: pending ? payouts.length : (SH.dash?.kpis.to_transfer_count || 0)})}</span></div>
      <div class="kpi"><span>${t('poNextShort')}</span><b style="font-size:20px">${new Date(next_run).toLocaleDateString([], {weekday: 'short', day: 'numeric', month: 'short'})}</b><span>${new Date(next_run).toLocaleTimeString([], {hour: '2-digit', minute: '2-digit'})} · ${st.frequency === 'daily' ? t('poDaily') : t('poWeekly')}</span></div>
      <div class="kpi"><span>${t('poHold')}</span><b>${st.hold_days}</b><span>${t('poHoldSub')}</span></div>
      <div class="kpi"><span>${t('poMin')}</span><b>$${esc(st.min_amount)}</b><span>${t('poMinSub')}</span></div></section>
    <section class="panel"><div class="ph"><h2>${t('poBreakdown')}</h2><div class="fpills" role="group">${['pending', 'paid', 'cancelled'].map(f => `<button type="button" aria-pressed="${APO.status === f}" onclick="APO.status='${f}';render()">${t('po_' + f)}</button>`).join('')}</div></div>
      <div class="tbl-wrap"><table><thead><tr><th>${t('store')}</th><th class="num">${t('poOrdersCol')}</th><th class="num">${t('rfAmount')}</th><th>${t('poAccount')}</th><th>${pending ? t('poAction') : t('status')}</th></tr></thead>
      <tbody>${rows || `<tr><td colspan="5" class="muted">${t('empty')}</td></tr>`}</tbody></table></div>
      <p class="foot">${t('poFoot')}</p></section>
    <details class="card"><summary style="font-weight:700;cursor:pointer">${t('poSettings')}</summary><div class="grid">
      <label>${t('poFreq')}<select id="sFreq"><option value="daily" ${st.frequency === 'daily' ? 'selected' : ''}>${t('poDaily')}</option><option value="weekly" ${st.frequency === 'weekly' ? 'selected' : ''}>${t('poWeekly')}</option></select></label>
      <label>${t('poWeekday')}<select id="sDay">${[1, 2, 3, 4, 5, 6, 0].map(d => `<option value="${d}" ${st.weekday === d ? 'selected' : ''}>${t('wd' + d)}</option>`).join('')}</select></label>
      <label>${t('poHour')}<input id="sHour" type="number" min="0" max="23" value="${st.hour}"></label>
      <label>${t('poHold')}<input id="sHold" type="number" min="0" max="30" value="${st.hold_days}"></label>
      <label>${t('poMin')}<input id="sMin" inputmode="decimal" value="${esc(st.min_amount)}"></label></div>
      <p class="muted">${t('poNext', {d: new Date(next_run).toLocaleString()})}</p>
      <div class="row"><button type="button" id="sSave">${t('poSave')}</button></div></details>`);
  $('#sFreq').onchange = () => { $('#sDay').disabled = $('#sFreq').value === 'daily'; };
  $('#sFreq').onchange();
  if (typeof exManualPayoutBlock === 'function') exManualPayoutBlock(); // extras: tarik saldo manual (super admin)
  $('#sSave').onclick = async () => {
    try {
      await api('/admin/api/payout-settings', {method: 'PUT', body: {frequency: $('#sFreq').value, weekday: +$('#sDay').value, hour: +$('#sHour').value,
        hold_days: +$('#sHold').value, min_amount: $('#sMin').value}});
      toast(t('poDone')); render();
    } catch (e) { showError(e, {frequency: '#sFreq', weekday: '#sDay', hour: '#sHour', hold_days: '#sHold', min_amount: '#sMin'}); }
  };
  $('#gen').onclick = async () => {
    $('#gen').disabled = true;
    try { APO.last = await api('/admin/api/payouts/generate', {method: 'POST'}); APO.status = 'pending'; render(); }
    catch (e) { toast(errMsg(e), 'error'); $('#gen').disabled = false; }
  };
  document.querySelectorAll('[data-copy]').forEach(b => b.onclick = () => { navigator.clipboard?.writeText(b.dataset.copy); toast(t('copied')); });
  document.querySelectorAll('[data-paid]').forEach(b => b.onclick = async () => {
    const id = b.dataset.paid;
    try { await api(`/admin/api/payouts/${id}/paid`, {method: 'POST', body: {reference: $('#ref_' + id).value, note: $('#note_' + id).value}}); toast(t('poDone')); render(); }
    catch (e) { showError(e, {reference: '#ref_' + id, note: '#note_' + id}); }
  });
  document.querySelectorAll('[data-cancel]').forEach(b => b.onclick = async () => {
    const id = b.dataset.cancel, note = prompt(t('poCancelNote'));
    if (note === null) return;
    try { await api(`/admin/api/payouts/${id}/cancel`, {method: 'POST', body: {note}}); toast(t('poDone')); render(); } catch (e) { toast(errMsg(e), 'error'); }
  });
}
