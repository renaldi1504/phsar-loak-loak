'use strict';
// Seller Center — extras: tombol "Tarik sekarang" di layar Uang (pemilik toko) dan label strike otomatis.
addStrings({
  en: {exWdNow: 'Withdraw now', exWdConfirm: 'Send {a} to {b} now? Phsar transfers it to your bank, usually within 1 business day.', exWdDone: 'Withdrawal requested — we will transfer {a} soon.',
    exWd_pending_exists: 'A payout is already on its way to your bank.', exWd_no_account: 'Add your bank account first.', exWd_account_hold: 'Bank account just changed — you can withdraw after {d}.',
    exWd_limit: 'You can withdraw again on {d}.', exWd_below_min: 'Minimum ${m} to withdraw now.', exWdByYou: 'Requested by you',
    sk_seller_timeout: 'Order not accepted in time', sk_pack_timeout: 'Order not packed in time',
    exSk_seller_timeout: 'Order {o} was cancelled automatically because it was not accepted in time.', exSk_pack_timeout: 'Order {o} was cancelled automatically because it was not packed in time.'},
  km: {exWdNow: 'ដកឥឡូវ', exWdConfirm: 'ផ្ញើ {a} ទៅ {b} ឥឡូវ? Phsar នឹងផ្ទេរទៅធនាគាររបស់អ្នក ជាធម្មតាក្នុង 1 ថ្ងៃធ្វើការ។', exWdDone: 'បានស្នើដកប្រាក់ — យើងនឹងផ្ទេរ {a} ឆាប់ៗ។',
    exWd_pending_exists: 'មានការទូទាត់មួយកំពុងទៅធនាគាររបស់អ្នករួចហើយ។', exWd_no_account: 'សូមបន្ថែមគណនីធនាគារជាមុនសិន។', exWd_account_hold: 'ទើបប្តូរគណនីធនាគារ — អាចដកបានក្រោយ {d}។',
    exWd_limit: 'អ្នកអាចដកម្តងទៀតនៅ {d}។', exWd_below_min: 'ត្រូវមានយ៉ាងតិច ${m} ដើម្បីដកឥឡូវ។', exWdByYou: 'អ្នកបានស្នើ',
    sk_seller_timeout: 'មិនបានទទួលការបញ្ជាទិញទាន់ពេល', sk_pack_timeout: 'មិនបានវេចខ្ចប់ការបញ្ជាទិញទាន់ពេល',
    exSk_seller_timeout: 'ការបញ្ជាទិញ {o} ត្រូវបានបោះបង់ដោយស្វ័យប្រវត្តិ ព្រោះមិនបានទទួលទាន់ពេល។', exSk_pack_timeout: 'ការបញ្ជាទិញ {o} ត្រូវបានបោះបង់ដោយស្វ័យប្រវត្តិ ព្រោះមិនបានវេចខ្ចប់ទាន់ពេល។'},
  zh: {exWdNow: '立即提现', exWdConfirm: '现在把 {a} 转到 {b}？Phsar 通常在 1 个工作日内转账。', exWdDone: '已申请提现——我们会尽快转账 {a}。',
    exWd_pending_exists: '已有一笔结算正在转入你的银行账户。', exWd_no_account: '请先添加银行账户。', exWd_account_hold: '银行账户刚更换——{d} 后可提现。',
    exWd_limit: '{d} 可再次提现。', exWd_below_min: '满 ${m} 才能立即提现。', exWdByYou: '你申请的提现',
    sk_seller_timeout: '未及时接单', sk_pack_timeout: '未及时打包',
    exSk_seller_timeout: '订单 {o} 因未及时接单已被自动取消。', exSk_pack_timeout: '订单 {o} 因未及时打包已被自动取消。'},
});

// Dipanggil dari tabPayouts: isi #exWd di dalam kartu saldo. Fitur mati → tidak tampil apa-apa.
async function exWithdraw(storeId, acc) {
  const host = document.getElementById('exWd');
  if (!host) return;
  let w;
  try { w = await api(`/v1/seller/stores/${storeId}/payouts/withdraw`); } catch { return; }
  if (!w.enabled || !document.body.contains(host)) return;
  const money = v => '$' + Number(v || 0).toFixed(2);
  if (!w.can) {
    const d = w.reason === 'account_hold' ? w.hold_until : w.reason === 'limit' ? w.next_at : null;
    host.innerHTML = `<p style="margin:8px 0 0;color:#c3cad8;font-size:14px">${esc(t('exWd_' + w.reason, {m: w.min_amount, d: d ? pdate(d) : ''}))}</p>`;
    return;
  }
  host.innerHTML = `<button type="button" class="primary" id="exWdBtn" style="margin-top:10px;min-height:48px;width:100%;border-radius:999px;font-size:16px;font-weight:800">${t('exWdNow')} · ${money(w.available)}</button>`;
  $('#exWdBtn').onclick = async () => {
    const b = acc ? t('bank_' + acc.bank) + ' ••' + acc.account_last4 : '';
    if (!confirm(t('exWdConfirm', {a: money(w.available), b}))) return;
    $('#exWdBtn').disabled = true;
    try {
      const r = await api(`/v1/seller/stores/${storeId}/payouts/withdraw`, {method: 'POST'});
      toast(t('exWdDone', {a: money(r.amount)}));
    } catch (e) {
      const reason = e && e.data && e.data.error && e.data.error.reason;
      toast(reason ? t('exWd_' + reason, {m: w.min_amount, d: ''}) : errMsg(e), 'error');
    }
    tabPayouts(storeId);
  };
}

// Teks strike otomatis dalam bahasa seller (catatan server berbahasa Inggris).
function exStrikeNote(k) {
  return k.auto ? t('exSk_' + k.reason, {o: k.order_no || ''}) : k.note;
}
