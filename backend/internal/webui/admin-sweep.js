'use strict';
// Admin keuangan: commission sweep — berapa komisi Phsar yang aman dipindahkan dari rekening penampung PayWay
// ke rekening perusahaan, catat transfernya, dan ringkasan per bulan.
addStrings({
  en: {swHint: 'Commission from completed orders can be moved from the PayWay collection account to the company account. Money for sellers, open orders and refunds always stays.',
    swAvail: 'Available to sweep', swAvailSub: 'safe to transfer now', swBalance: 'Collection account (expected)', swBalanceSub: 'compare with the PayWay settlement balance',
    swHeldSellers: 'Held for sellers', swHeldSellersSub: 'balances + payouts not yet sent ({p})', swHeldOrders: 'Held for orders', swHeldOrdersSub: 'not yet completed, refunds due',
    swCalc: 'How it is calculated', swCollected: 'Paid by buyers', swRefunds: 'Refunded to buyers', swPayouts: 'Paid to sellers', swSwept: 'Already swept',
    swCommGross: 'Commission on completed orders', swCommBack: 'Commission returned on refunds', swCommHold: 'Commission on approved returns (waiting refund)', swCommNet: 'Net commission',
    swFree: 'Cash not owed to anyone', swRecord: 'Record a sweep', swRecordHint: 'Transfer the money in the PayWay / ABA merchant portal first, then record it here.',
    swAmount: 'Amount ($)', swTo: 'To company account', swToPh: 'ABA Phsar Co. ••• 1234', swRef: 'Bank reference', swDate: 'Transfer date', swNote: 'Note (optional)',
    swSave: 'Record sweep', swSaved: 'Sweep recorded', swMax: 'Use maximum', swHistory: 'Sweeps', swNone: 'No sweeps yet.', swVoid: 'Cancel', swVoided: 'Cancelled',
    swVoidAsk: 'Why is this sweep cancelled? (the money counts as still in the collection account)', swVoidDone: 'Sweep cancelled', swMonthly: 'Monthly',
    swMonth: 'Month', swCommission: 'Commission', swExport: 'Export CSV', swBy: 'by {n}', swNothing: 'Nothing to sweep right now.'},
  km: {swHint: 'កម្រៃជើងសារពីការបញ្ជាទិញដែលបានបញ្ចប់ អាចផ្ទេរពីគណនីប្រមូលប្រាក់ PayWay ទៅគណនីក្រុមហ៊ុន។ ប្រាក់សម្រាប់អ្នកលក់ ការបញ្ជាទិញកំពុងដំណើរការ និងការសងប្រាក់ ត្រូវរក្សាទុកជានិច្ច។',
    swAvail: 'អាចផ្ទេរបាន', swAvailSub: 'សុវត្ថិភាពក្នុងការផ្ទេរឥឡូវនេះ', swBalance: 'គណនីប្រមូលប្រាក់ (រំពឹងទុក)', swBalanceSub: 'ប្រៀបធៀបជាមួយសមតុល្យ PayWay',
    swHeldSellers: 'រក្សាទុកសម្រាប់អ្នកលក់', swHeldSellersSub: 'សមតុល្យ + ការទូទាត់មិនទាន់ផ្ញើ ({p})', swHeldOrders: 'រក្សាទុកសម្រាប់ការបញ្ជាទិញ', swHeldOrdersSub: 'មិនទាន់បញ្ចប់ ការសងប្រាក់ត្រូវធ្វើ',
    swCalc: 'របៀបគណនា', swCollected: 'អ្នកទិញបានបង់', swRefunds: 'បានសងអ្នកទិញ', swPayouts: 'បានទូទាត់អ្នកលក់', swSwept: 'បានផ្ទេររួច',
    swCommGross: 'កម្រៃពីការបញ្ជាទិញបានបញ្ចប់', swCommBack: 'កម្រៃដែលបានសងវិញ', swCommHold: 'កម្រៃពីការប្ដូរបានយល់ព្រម (រង់ចាំសង)', swCommNet: 'កម្រៃសុទ្ធ',
    swFree: 'ប្រាក់មិនជំពាក់នរណា', swRecord: 'កត់ត្រាការផ្ទេរ', swRecordHint: 'ផ្ទេរប្រាក់ក្នុងផ្ទាំង PayWay / ABA merchant ជាមុនសិន រួចកត់ត្រានៅទីនេះ។',
    swAmount: 'ចំនួន ($)', swTo: 'ទៅគណនីក្រុមហ៊ុន', swToPh: 'ABA Phsar Co. ••• 1234', swRef: 'លេខយោងធនាគារ', swDate: 'ថ្ងៃផ្ទេរ', swNote: 'កំណត់ចំណាំ (មិនចាំបាច់)',
    swSave: 'កត់ត្រាការផ្ទេរ', swSaved: 'បានកត់ត្រា', swMax: 'ប្រើអតិបរមា', swHistory: 'ការផ្ទេរ', swNone: 'មិនទាន់មានការផ្ទេរ។', swVoid: 'បោះបង់', swVoided: 'បានបោះបង់',
    swVoidAsk: 'ហេតុអ្វីបោះបង់ការផ្ទេរនេះ? (ប្រាក់នឹងរាប់ថានៅក្នុងគណនីប្រមូលប្រាក់)', swVoidDone: 'បានបោះបង់', swMonthly: 'ប្រចាំខែ',
    swMonth: 'ខែ', swCommission: 'កម្រៃ', swExport: 'នាំចេញ CSV', swBy: 'ដោយ {n}', swNothing: 'គ្មានអ្វីត្រូវផ្ទេរពេលនេះទេ។'},
  zh: {swHint: '已完成订单的佣金可以从 PayWay 收款账户转到公司账户。属于卖家、未完成订单和待退款的资金始终保留。',
    swAvail: '可转出', swAvailSub: '现在可安全转出', swBalance: '收款账户（应有余额）', swBalanceSub: '请与 PayWay 结算余额核对',
    swHeldSellers: '为卖家保留', swHeldSellersSub: '余额 + 未转账的结算（{p}）', swHeldOrders: '为订单保留', swHeldOrdersSub: '未完成订单、待退款',
    swCalc: '计算方式', swCollected: '买家已付款', swRefunds: '已退款给买家', swPayouts: '已结算给卖家', swSwept: '已转出',
    swCommGross: '已完成订单佣金', swCommBack: '退款退回的佣金', swCommHold: '已同意退货的佣金（待退款）', swCommNet: '佣金净额',
    swFree: '未欠任何人的资金', swRecord: '记录转出', swRecordHint: '请先在 PayWay / ABA 商户后台完成转账，再在此记录。',
    swAmount: '金额（$）', swTo: '转入公司账户', swToPh: 'ABA Phsar Co. ••• 1234', swRef: '银行参考号', swDate: '转账日期', swNote: '备注（可选）',
    swSave: '记录转出', swSaved: '已记录', swMax: '使用最大值', swHistory: '转出记录', swNone: '暂无转出记录。', swVoid: '作废', swVoided: '已作废',
    swVoidAsk: '为什么作废这笔转出？（该金额将重新计入收款账户）', swVoidDone: '已作废', swMonthly: '按月',
    swMonth: '月份', swCommission: '佣金', swExport: '导出 CSV', swBy: '由 {n}', swNothing: '目前没有可转出的金额。'},
});

const SWP = {data: null};

async function renderSweep() {
  SWP.data = await api('/admin/api/sweep');
  drawSweep();
}

function drawSweep() {
  const {summary: s, sweeps, monthly} = SWP.data;
  const row = (label, v, sign = '', strong = false) => `<tr${strong ? ' style="font-weight:700"' : ''}><td>${label}</td><td style="text-align:right">${sign}${amoney(v)}</td></tr>`;
  const today = new Date(Date.now() + 7 * 3600e3).toISOString().slice(0, 10);
  const hist = sweeps.map(x => `<tr${x.voided_at ? ' class="muted"' : ''}><td>${esc(x.swept_on)}</td><td style="text-align:right"><b>${x.voided_at ? '<s>' + amoney(x.amount) + '</s>' : amoney(x.amount)}</b></td>
      <td>${esc(x.to_account)}<br><span class="muted small">${esc(x.reference)}${x.note ? ' · ' + esc(x.note) : ''}</span></td>
      <td class="muted small">${esc(t('swBy', {n: x.created_by}))}<br>${new Date(x.created_at).toLocaleString()}</td>
      <td>${x.voided_at ? `<span class="chip">${t('swVoided')}</span><br><span class="muted small">${esc(x.void_note || '')}</span>` : `<button type="button" data-void="${x.id}">${t('swVoid')}</button>`}</td></tr>`).join('');
  const months = monthly.map(m => `<tr><td>${esc(m.month)}</td><td style="text-align:right">${amoney(m.collected)}</td><td style="text-align:right">${amoney(m.refunds)}</td>
      <td style="text-align:right">${amoney(m.payouts)}</td><td style="text-align:right">${amoney(m.commission)}</td><td style="text-align:right">${amoney(m.swept)}</td></tr>`).join('');
  const avail = Number(s.available) > 0;
  $('#app').innerHTML = shell(`${pageHead(t('i_sweep'), t('swHint'))}
    <section class="kpis">
      <div class="kpi"><span>${t('swAvail')}</span><b class="${avail ? 'ok' : ''}">${amoney(s.available)}</b><span>${t('swAvailSub')}</span></div>
      <div class="kpi"><span>${t('swBalance')}</span><b>${amoney(s.expected_balance)}</b><span>${t('swBalanceSub')}</span></div>
      <div class="kpi"><span>${t('swHeldSellers')}</span><b>${amoney(s.held_sellers)}</b><span>${t('swHeldSellersSub', {p: amoney(s.pending_payouts)})}</span></div>
      <div class="kpi"><span>${t('swHeldOrders')}</span><b>${amoney(s.held_orders)}</b><span>${t('swHeldOrdersSub')}</span></div>
    </section>
    <div class="cols">
      <section class="panel" id="swForm"><div class="ph"><h2>${t('swRecord')}</h2></div>
        <div style="padding:16px 18px;display:flex;flex-direction:column;gap:14px">
          <p class="muted small" style="margin:0">${t('swRecordHint')}</p>
          ${avail ? '' : `<div class="banner info">${t('swNothing')}</div>`}
          <div class="grid">
            <label>${t('swAmount')}<input id="swAmount" name="amount" inputmode="decimal" placeholder="0.00"><button type="button" class="link small" id="swMax" style="align-self:flex-start">${t('swMax')} (${amoney(s.available)})</button></label>
            <label>${t('swDate')}<input id="swDate" name="swept_on" type="date" value="${today}" max="${today}"></label>
            <label>${t('swTo')}<input id="swTo" name="to_account" maxlength="120" placeholder="${esc(t('swToPh'))}"></label>
            <label>${t('swRef')}<input id="swRef" name="reference" maxlength="80"></label>
          </div>
          <label>${t('swNote')}<input id="swNote" name="note" maxlength="300"></label>
          <div><button type="button" class="primary" id="swSave" ${avail ? '' : 'disabled'}>${t('swSave')}</button></div>
        </div></section>
      <section class="panel"><div class="ph"><h2>${t('swCalc')}</h2></div>
        <div class="tbl-wrap" style="padding:8px 18px 16px"><table style="min-width:0;width:100%"><tbody>
          ${row(t('swCollected'), s.collected)}${row(t('swRefunds'), s.refunds, '− ')}${row(t('swPayouts'), s.payouts_paid, '− ')}${row(t('swSwept'), s.swept, '− ')}
          ${row(t('swBalance'), s.expected_balance, '', true)}
          ${row(t('swHeldSellers'), s.held_sellers, '− ')}${row(t('swHeldOrders'), s.held_orders, '− ')}${row(t('swFree'), s.free_cash, '', true)}
          <tr><td colspan="2" style="padding-top:14px"></td></tr>
          ${row(t('swCommGross'), s.commission_gross)}${row(t('swCommBack'), s.commission_back, '− ')}${row(t('swCommNet'), s.commission_net, '', true)}
          ${row(t('swCommHold'), s.commission_on_hold, '− ')}${row(t('swSwept'), s.swept, '− ')}
          ${row(t('swAvail'), s.available, '', true)}
        </tbody></table></div></section>
    </div>
    <section class="panel"><div class="ph"><h2>${t('swHistory')}</h2><a class="btn" href="/admin/api/sweep/export.csv?what=sweeps" data-csv>${t('swExport')}</a></div>
      ${hist ? `<div class="tbl-wrap"><table><tbody>${hist}</tbody></table></div>` : `<p class="muted" style="padding:0 18px 16px">${t('swNone')}</p>`}</section>
    <section class="panel"><div class="ph"><h2>${t('swMonthly')}</h2><a class="btn" href="/admin/api/sweep/export.csv?what=monthly" data-csv>${t('swExport')}</a></div>
      <div class="tbl-wrap"><table><thead><tr><th>${t('swMonth')}</th><th style="text-align:right">${t('swCollected')}</th><th style="text-align:right">${t('swRefunds')}</th>
        <th style="text-align:right">${t('swPayouts')}</th><th style="text-align:right">${t('swCommission')}</th><th style="text-align:right">${t('swSwept')}</th></tr></thead><tbody>${months}</tbody></table></div></section>`);

  $('#swMax').onclick = () => { $('#swAmount').value = s.available; };
  $('#swSave').onclick = async () => {
    document.querySelectorAll('#swForm .invalid').forEach(el => el.classList.remove('invalid'));
    const body = {amount: $('#swAmount').value.trim(), swept_on: $('#swDate').value, to_account: $('#swTo').value.trim(), reference: $('#swRef').value.trim(), note: $('#swNote').value.trim()};
    try { SWP.data = await api('/admin/api/sweep', {method: 'POST', body}); toast(t('swSaved')); drawSweep(); }
    catch (e) { showError(e, {amount: '#swAmount', swept_on: '#swDate', to_account: '#swTo', reference: '#swRef', note: '#swNote'}); }
  };
  document.querySelectorAll('[data-void]').forEach(b => b.onclick = async () => {
    const note = prompt(t('swVoidAsk'));
    if (!note) return;
    try { SWP.data = await api(`/admin/api/sweep/${b.dataset.void}/void`, {method: 'POST', body: {note}}); toast(t('swVoidDone')); drawSweep(); }
    catch (e) { toast(errMsg(e), 'error'); }
  });
  // CSV butuh token → unduh sebagai blob.
  document.querySelectorAll('[data-csv]').forEach(a => a.onclick = async ev => {
    ev.preventDefault();
    const r = await fetch(a.getAttribute('href'), {headers: {Authorization: 'Bearer ' + P.token}});
    if (!r.ok) { toast(r.statusText, 'error'); return; }
    const url = URL.createObjectURL(await r.blob()), link = document.createElement('a');
    link.href = url; link.download = (r.headers.get('Content-Disposition') || '').split('filename="')[1]?.replace('"', '') || 'phsar-sweep.csv';
    document.body.appendChild(link); link.click(); link.remove();
  });
}
