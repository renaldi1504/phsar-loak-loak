'use strict';
// Seller Center: saldo, jadwal payout, rekening bank, riwayat (khusus pemilik toko).
addStrings({
  en: {mh_ready: 'Money ready to be paid', mh_auto: 'We send it automatically on {d} to {b}.', mh_noacc: 'Add your bank account so we can send it.', mh_min: 'Sent when it reaches at least ${m}.',
    mf_title: 'How your money moves', mf_1: 'On the way to the buyer', mf_1n: 'Paid, not delivered yet.', mf_2: 'Waiting {d} days', mf_2n: 'After delivery we wait {d} days in case the buyer returns it.', mf_3: 'Ready to be paid', mf_3n: 'Sent on the next payout day.', mf_4: 'Being transferred', mf_4n: 'On the way to your bank.',
    fee_q: 'Why do I get less than the price?', fee_ex: 'Example: buyer pays', fee_cut: 'Phsar fee ({r}%)', fee_get: 'You get', fee_note: 'The fee covers payment, courier handling and promoting your shop. No monthly fee.', moneyHistory: 'History', tabPayouts: 'Payouts', balAvailable: 'Available balance', balOnHold: 'On hold', balUpcoming: 'Coming soon', balPending: 'Being transferred',
    balOnHoldHint: 'Delivered orders wait {d} days for buyer complaints, then move to your balance.', balUpcomingHint: 'Paid orders not delivered yet.',
    nextPayout: 'Next payout: {d}', schedDaily: 'Payouts are made every day.', schedWeekly: 'Payouts are made every {w}.', minPayout: 'Minimum payout: ${m}.',
    bankTitle: 'Bank account for payouts', bankNone: 'Add your bank account so Phsar can pay you.', bank: 'Bank', accName: 'Name on the account', accNumber: 'Account number',
    saveAccount: 'Save account', accSaved: 'Bank account saved.', accHold: 'Account changed — payouts are paused until {d} for your security.', currentAcc: 'Current account',
    changeAcc: 'Change account', bank_aba: 'ABA Bank', bank_acleda: 'ACLEDA', bank_wing: 'Wing', bank_bakong: 'Bakong', bank_other: 'Other',
    historyTitle: 'Payout history', ledgerTitle: 'Balance activity', noHistory: 'No payouts yet.', ps_pending: 'Being transferred', ps_paid: 'Paid', ps_cancelled: 'Cancelled',
    lk_order_earning: 'Order {o} completed', lk_payout: 'Payout', lk_payout_cancelled: 'Payout cancelled — returned to balance', lk_adjustment: 'Adjustment',
    ref: 'Ref.', orders: '{n} orders', wd0: 'Sunday', wd1: 'Monday', wd2: 'Tuesday', wd3: 'Wednesday', wd4: 'Thursday', wd5: 'Friday', wd6: 'Saturday'},
  km: {mh_ready: 'លុយរួចរាល់សម្រាប់ទូទាត់', mh_auto: 'យើងផ្ញើដោយស្វ័យប្រវត្តិនៅ {d} ទៅ {b}។', mh_noacc: 'បន្ថែមគណនីធនាគារ ដើម្បីឱ្យយើងអាចផ្ញើបាន។', mh_min: 'ផ្ញើនៅពេលដល់យ៉ាងតិច ${m}។',
    mf_title: 'លុយរបស់អ្នកដំណើរការយ៉ាងដូចម្តេច', mf_1: 'កំពុងដឹកទៅអ្នកទិញ', mf_1n: 'បានបង់ មិនទាន់ដឹកដល់។', mf_2: 'រង់ចាំ {d} ថ្ងៃ', mf_2n: 'ក្រោយដឹកដល់ យើងរង់ចាំ {d} ថ្ងៃ ក្រែងអ្នកទិញប្តូរវិញ។', mf_3: 'រួចរាល់សម្រាប់ទូទាត់', mf_3n: 'ផ្ញើនៅថ្ងៃទូទាត់បន្ទាប់។', mf_4: 'កំពុងផ្ទេរ', mf_4n: 'កំពុងទៅធនាគាររបស់អ្នក។',
    fee_q: 'ហេតុអ្វីខ្ញុំទទួលបានតិចជាងតម្លៃ?', fee_ex: 'ឧទាហរណ៍៖ អ្នកទិញបង់', fee_cut: 'កម្រៃ Phsar ({r}%)', fee_get: 'អ្នកទទួលបាន', fee_note: 'កម្រៃនេះសម្រាប់ការទូទាត់ ការដឹកជញ្ជូន និងការផ្សព្វផ្សាយហាង។ គ្មានថ្លៃប្រចាំខែ។', moneyHistory: 'ប្រវត្តិ', tabPayouts: 'ការទូទាត់', balAvailable: 'សមតុល្យអាចដកបាន', balOnHold: 'កំពុងរង់ចាំ', balUpcoming: 'នឹងមកដល់', balPending: 'កំពុងផ្ទេរ',
    balOnHoldHint: 'ការបញ្ជាទិញដែលបានដឹកដល់ រង់ចាំ {d} ថ្ងៃសម្រាប់ការតវ៉ា បន្ទាប់មកចូលសមតុល្យ។', balUpcomingHint: 'ការបញ្ជាទិញបានបង់ ប៉ុន្តែមិនទាន់ដឹកដល់។',
    nextPayout: 'ការទូទាត់បន្ទាប់៖ {d}', schedDaily: 'ទូទាត់រៀងរាល់ថ្ងៃ។', schedWeekly: 'ទូទាត់រៀងរាល់ថ្ងៃ{w}។', minPayout: 'ទូទាត់អប្បបរមា ${m}។',
    bankTitle: 'គណនីធនាគារសម្រាប់ទទួលប្រាក់', bankNone: 'បន្ថែមគណនីធនាគារ ដើម្បីឱ្យ Phsar អាចផ្ទេរប្រាក់ជូនអ្នក។', bank: 'ធនាគារ', accName: 'ឈ្មោះលើគណនី', accNumber: 'លេខគណនី',
    saveAccount: 'រក្សាទុកគណនី', accSaved: 'បានរក្សាទុកគណនីធនាគារ។', accHold: 'បានប្តូរគណនី — ផ្អាកការទូទាត់រហូតដល់ {d} ដើម្បីសុវត្ថិភាព។', currentAcc: 'គណនីបច្ចុប្បន្ន',
    changeAcc: 'ប្តូរគណនី', bank_aba: 'ABA Bank', bank_acleda: 'ACLEDA', bank_wing: 'Wing', bank_bakong: 'Bakong', bank_other: 'ផ្សេងៗ',
    historyTitle: 'ប្រវត្តិការទូទាត់', ledgerTitle: 'សកម្មភាពសមតុល្យ', noHistory: 'មិនទាន់មានការទូទាត់។', ps_pending: 'កំពុងផ្ទេរ', ps_paid: 'បានទូទាត់', ps_cancelled: 'បានបោះបង់',
    lk_order_earning: 'ការបញ្ជាទិញ {o} បានបញ្ចប់', lk_payout: 'ការទូទាត់', lk_payout_cancelled: 'បោះបង់ការទូទាត់ — បានត្រឡប់ចូលសមតុល្យ', lk_adjustment: 'ការកែតម្រូវ',
    ref: 'លេខយោង', orders: 'ការបញ្ជាទិញ {n}', wd0: 'អាទិត្យ', wd1: 'ចន្ទ', wd2: 'អង្គារ', wd3: 'ពុធ', wd4: 'ព្រហស្បតិ៍', wd5: 'សុក្រ', wd6: 'សៅរ៍'},
  zh: {mh_ready: '待结算金额', mh_auto: '将于 {d} 自动转到 {b}。', mh_noacc: '请添加银行账户，以便我们打款。', mh_min: '满 ${m} 时打款。',
    mf_title: '钱是怎么到账的', mf_1: '配送中', mf_1n: '已付款，尚未送达。', mf_2: '等待 {d} 天', mf_2n: '送达后等待 {d} 天，以防买家退货。', mf_3: '待结算', mf_3n: '下一个结算日打款。', mf_4: '转账中', mf_4n: '正在转入你的银行账户。',
    fee_q: '为什么比售价少？', fee_ex: '例如：买家支付', fee_cut: 'Phsar 佣金（{r}%）', fee_get: '你收到', fee_note: '佣金用于支付处理、快递服务和店铺推广，无月费。', moneyHistory: '记录', tabPayouts: '结算', balAvailable: '可结算余额', balOnHold: '冻结中', balUpcoming: '即将入账', balPending: '转账中',
    balOnHoldHint: '已送达订单需等待 {d} 天售后期，之后转入余额。', balUpcomingHint: '买家已付款但尚未送达的订单。',
    nextPayout: '下次结算：{d}', schedDaily: '每天结算。', schedWeekly: '每{w}结算。', minPayout: '最低结算金额 ${m}。',
    bankTitle: '收款银行账户', bankNone: '添加银行账户，Phsar 才能向你付款。', bank: '银行', accName: '账户名', accNumber: '账号',
    saveAccount: '保存账户', accSaved: '银行账户已保存。', accHold: '账户已更改——为保障安全，结算暂停至 {d}。', currentAcc: '当前账户',
    changeAcc: '更换账户', bank_aba: 'ABA Bank', bank_acleda: 'ACLEDA', bank_wing: 'Wing', bank_bakong: 'Bakong', bank_other: '其他',
    historyTitle: '结算记录', ledgerTitle: '余额明细', noHistory: '暂无结算。', ps_pending: '转账中', ps_paid: '已到账', ps_cancelled: '已取消',
    lk_order_earning: '订单 {o} 已完成', lk_payout: '结算', lk_payout_cancelled: '结算取消——已退回余额', lk_adjustment: '调整',
    ref: '参考号', orders: '{n} 个订单', wd0: '周日', wd1: '周一', wd2: '周二', wd3: '周三', wd4: '周四', wd5: '周五', wd6: '周六'},
});

const pchip = (cls, text) => `<span class="chip ${cls}">${esc(text)}</span>`;
function pmoney(v) { const n = parseFloat(v); return (n < 0 ? '−$' : '$') + Math.abs(n).toFixed(2); }
function pdate(s, time = true) {
  try { return new Date(s).toLocaleString(P.lang === 'zh' ? 'zh-CN' : P.lang === 'km' ? 'km-KH' : 'en-GB', time ? {dateStyle: 'medium', timeStyle: 'short'} : {dateStyle: 'full'}); } catch { return s; }
}

async function tabPayouts(id) {
  const s = await api(`/v1/seller/stores/${id}/payouts`);
  const a = s.account;
  const hold = a && a.payout_hold_until && new Date(a.payout_hold_until) > new Date();
  const bal = (label, v, hint, cls = '') => `<div class="stat ${cls}"><span class="muted">${label}</span><b>${pmoney(v)}</b>${hint ? `<small class="muted">${hint}</small>` : ''}</div>`;
  const hist = s.payouts.map(p => `<tr><td>${pdate(p.created_at)}</td><td><b>${pmoney(p.amount)}</b><br><span class="muted small">${t('orders', {n: p.order_count})}</span></td>
    <td>${esc(t('bank_' + p.bank))} ••${esc(p.account_last4)}</td><td>${pchip(p.status === 'paid' ? 'ok' : p.status === 'pending' ? 'warn' : 'bad', t('ps_' + p.status))}
    ${p.reference ? `<br><span class="muted small">${t('ref')} ${esc(p.reference)}</span>` : ''}${p.requested_by_seller ? `<br><span class="muted small">${t('exWdByYou')}</span>` : ''}</td></tr>`).join('');
  const ledger = s.ledger.map(l => `<tr><td class="muted">${pdate(l.at)}</td><td>${esc(t('lk_' + l.kind, {o: l.order_no}))}${l.note && l.kind !== 'order_earning' ? `<br><span class="muted small">${esc(l.note)}</span>` : ''}</td>
    <td style="text-align:right;font-weight:700;color:${parseFloat(l.amount) < 0 ? 'inherit' : 'var(--ok)'}">${parseFloat(l.amount) > 0 ? '+' : ''}${pmoney(l.amount)}</td></tr>`).join('');
  const rate = parseFloat(S.store?.commission?.store || S.store?.commission?.global || 8);
  const ex = 10, cut = Math.round(ex * rate) / 100;
  const fl = (n, title, note, v, bg, fg) => `<div class="fl" style="--fb:${bg};--ff:${fg}"><div class="rail"><i>${n}</i><s></s></div>
    <div class="fx"><div><b>${title}</b><b>${pmoney(v)}</b></div><span>${note}</span></div></div>`;
  $('#tab').innerHTML = `<section class="moneyhero"><span>${t('mh_ready')}</span><b>${pmoney(s.available)}</b>
      <div class="when">${a ? t('mh_auto', {d: pdate(s.next_run, false), b: esc(t('bank_' + a.bank)) + ' ••' + esc(a.account_last4)}) : t('mh_noacc')}
        <br><span style="color:#c3cad8">${s.frequency === 'daily' ? t('schedDaily') : t('schedWeekly', {w: t('wd' + s.weekday)})} ${t('mh_min', {m: s.min_amount})}</span></div><div id="exWd"></div></section>
    ${hold ? `<div class="banner warn">${t('accHold', {d: pdate(a.payout_hold_until)})}</div>` : ''}
    <section class="card flow"><h2>${t('mf_title')}</h2>
      ${fl(1, t('mf_1'), t('mf_1n'), s.upcoming, '#f1eef3', '#5f6678')}
      ${fl(2, t('mf_2', {d: s.hold_days}), t('mf_2n', {d: s.hold_days}), s.on_hold, '#fff0dc', '#8a4300')}
      ${fl(3, t('mf_3'), t('mf_3n'), s.available, '#e3f4ea', '#0b6b3c')}
      ${parseFloat(s.pending_payout) > 0 ? fl(4, t('mf_4'), t('mf_4n'), s.pending_payout, '#e6eaf2', '#102241') : ''}</section>
    <details class="card feebox"><summary>${t('fee_q')}</summary>
      <dl><div><dt>${t('fee_ex')}</dt><dd>${pmoney(ex)}</dd></div><div style="color:var(--bad)"><dt>${t('fee_cut', {r: rate})}</dt><dd>−${pmoney(cut)}</dd></div>
        <div style="border-top:1px solid var(--line);padding-top:6px"><dt>${t('fee_get')}</dt><dd style="color:#0b6b3c">${pmoney(ex - cut)}</dd></div></dl>
      <p class="muted" style="margin-top:8px">${t('fee_note')}</p></details>
    <section class="card"><h2>${t('bankTitle')}</h2>
      ${a ? `<p>${t('currentAcc')}: <b>${esc(t('bank_' + a.bank))}</b> · ${esc(a.account_name)} · ••••${esc(a.account_last4)}</p>` : `<div class="banner warn">${t('bankNone')}</div>`}
      <details ${a ? '' : 'open'}><summary class="link">${a ? t('changeAcc') : t('bankTitle')}</summary>
        <div class="grid"><label>${t('bank')}<select id="pBank">${['aba', 'acleda', 'wing', 'bakong', 'other'].map(b => `<option value="${b}" ${a && a.bank === b ? 'selected' : ''}>${t('bank_' + b)}</option>`).join('')}</select></label>
          <label>${t('accName')}<input id="pName" maxlength="120" value="${esc(a ? a.account_name : '')}"></label>
          <label>${t('accNumber')}<input id="pNumber" inputmode="numeric" autocomplete="off" maxlength="30"></label></div>
        <div class="row"><button type="button" class="primary" id="pSave">${t('saveAccount')}</button></div></details></section>
    ${ledger ? `<section class="card" id="ledgerCard"><h2>${t('ledgerTitle')}</h2><div id="ledgerBody"><table>${ledger}</table></div></section>` : ''}
    <section class="card"><h2>${t('historyTitle')}</h2>${hist ? `<table>${hist}</table>` : `<p class="muted">${t('noHistory')}</p>`}</section>`;
  $('#pSave').onclick = async () => {
    try {
      await api(`/v1/seller/stores/${id}/payouts/account`, {method: 'PUT', body: {bank: $('#pBank').value, account_name: $('#pName').value, account_number: $('#pNumber').value}});
      toast(t('accSaved')); tabPayouts(id);
    } catch (e) { showError(e, {bank: '#pBank', account_name: '#pName', account_number: '#pNumber'}); }
  };
  if (ledger) ledgerPage(id, 1).catch(() => {}); // ganti 50 baris ringkasan dengan versi berhalaman
  if (typeof sellerInvoices === 'function') sellerInvoices(id); // invoices: faktur komisi bulanan
  if (typeof exWithdraw === 'function') exWithdraw(id, a); // extras: tombol "Tarik sekarang"
}

// Riwayat saldo per halaman (ringkasan /payouts hanya membawa 50 baris terbaru).
function ledgerRows(rows) {
  return rows.map(l => `<tr><td class="muted">${pdate(l.at)}</td><td>${esc(t('lk_' + l.kind, {o: l.order_no}))}${l.note && l.kind !== 'order_earning' ? `<br><span class="muted small">${esc(l.note)}</span>` : ''}</td>
    <td style="text-align:right;font-weight:700;color:${parseFloat(l.amount) < 0 ? 'inherit' : 'var(--ok)'}">${parseFloat(l.amount) > 0 ? '+' : ''}${pmoney(l.amount)}</td></tr>`).join('');
}
async function ledgerPage(id, page) {
  const box = $('#ledgerBody');
  if (!box) return;
  const d = await api(`/v1/seller/stores/${id}/payouts/ledger?${qs({page, size: 20})}`);
  box.innerHTML = `<table>${ledgerRows(d.ledger)}</table>${pagerHTML(d.page)}`;
  box.querySelectorAll('.pager [data-page]').forEach(b => b.onclick = () => { ledgerPage(id, +b.dataset.page); $('#ledgerCard').scrollIntoView({block: 'start'}); });
}
