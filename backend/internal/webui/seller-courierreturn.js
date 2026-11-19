'use strict';
// courierreturn: Seller Center — paket yang gagal diantar dikembalikan kurir ke titik jemput toko.
// Tugas di "Hari ini", kartu di atas Pesanan, riwayat di "Riwayat: Dikembalikan".
addStrings({
  en: {tk_crr: 'Confirm {n} returned package(s)', tk_crrW: 'The courier brought back orders that could not be delivered. Check them and confirm, otherwise Phsar treats them as received.', tk_crrB: 'Check',
    c_delivery_failed: 'Delivery failed — returned to you', crrTitle: 'Packages coming back', crrHistory: 'Returned', crrBack: '← Orders',
    crrs_returning: 'Coming back', crrs_awaiting_seller: 'Please confirm', crrs_received: 'Received',
    crrW_returning: 'The courier could not deliver this order and is bringing it back to {p}. Courier: {c}.',
    crrW_awaiting_seller: 'The courier says the package is back at {p}. Check it, then confirm.',
    crrOk: 'Received — put back in stock', crrDamaged: 'Received, but damaged', crrOkQ: 'Confirm you received the package in good condition? The items go back to your stock and the order is cancelled.',
    crrDmgTitle: 'Package damaged?', crrDmgHint: 'The items will not go back to your stock. The order is cancelled and the buyer refunded.', crrDmgNote: 'What is damaged?', crrDmgGo: 'Confirm damaged',
    crrDone: 'Thank you — recorded.', crrConfirmIn: 'Confirm within {t}', crrLate: 'Time is up — Phsar will record it as received', crrEmpty: 'No returned packages.',
    crrR_seller: 'You confirmed', crrR_admin: 'Recorded by Phsar', crrR_system: 'Recorded automatically', crrCond_ok: 'back in stock', crrCond_damaged: 'damaged, not back in stock',
    crrMoney: 'The buyer gets a full refund. You are not charged for the return delivery.'},
  km: {tk_crr: 'បញ្ជាក់កញ្ចប់ដែលត្រឡប់មកវិញ {n}', tk_crrW: 'អ្នកដឹកបាននាំការបញ្ជាទិញដែលដឹកមិនបានត្រឡប់មកវិញ។ សូមពិនិត្យ ហើយបញ្ជាក់ បើមិនដូច្នោះទេ Phsar នឹងចាត់ទុកថាបានទទួល។', tk_crrB: 'ពិនិត្យ',
    c_delivery_failed: 'ដឹកមិនបាន — បានប្រគល់មកអ្នកវិញ', crrTitle: 'កញ្ចប់កំពុងត្រឡប់មកវិញ', crrHistory: 'បានត្រឡប់', crrBack: '← ការបញ្ជាទិញ',
    crrs_returning: 'កំពុងត្រឡប់មក', crrs_awaiting_seller: 'សូមបញ្ជាក់', crrs_received: 'បានទទួល',
    crrW_returning: 'អ្នកដឹកមិនអាចដឹកជូនការបញ្ជាទិញនេះបាន ហើយកំពុងនាំវាត្រឡប់មក {p}។ អ្នកដឹក៖ {c}។',
    crrW_awaiting_seller: 'អ្នកដឹកប្រាប់ថាកញ្ចប់បានត្រឡប់មកដល់ {p} ហើយ។ សូមពិនិត្យ រួចបញ្ជាក់។',
    crrOk: 'បានទទួល — ដាក់ចូលស្តុកវិញ', crrDamaged: 'បានទទួល ប៉ុន្តែខូច', crrOkQ: 'បញ្ជាក់ថាអ្នកបានទទួលកញ្ចប់ក្នុងស្ថានភាពល្អ? ទំនិញនឹងចូលស្តុកវិញ ហើយការបញ្ជាទិញត្រូវបោះបង់។',
    crrDmgTitle: 'កញ្ចប់ខូច?', crrDmgHint: 'ទំនិញនឹងមិនចូលស្តុកវិញទេ។ ការបញ្ជាទិញត្រូវបោះបង់ ហើយអ្នកទិញទទួលប្រាក់វិញ។', crrDmgNote: 'អ្វីខ្លះខូច?', crrDmgGo: 'បញ្ជាក់ថាខូច',
    crrDone: 'អរគុណ — បានកត់ត្រា។', crrConfirmIn: 'បញ្ជាក់ក្នុង {t}', crrLate: 'ផុតពេល — Phsar នឹងកត់ត្រាថាបានទទួល', crrEmpty: 'មិនមានកញ្ចប់ត្រឡប់ទេ។',
    crrR_seller: 'អ្នកបានបញ្ជាក់', crrR_admin: 'Phsar បានកត់ត្រា', crrR_system: 'កត់ត្រាដោយស្វ័យប្រវត្តិ', crrCond_ok: 'ចូលស្តុកវិញ', crrCond_damaged: 'ខូច មិនចូលស្តុកវិញ',
    crrMoney: 'អ្នកទិញទទួលប្រាក់វិញពេញ។ អ្នកមិនត្រូវបង់ថ្លៃដឹកត្រឡប់ទេ។'},
  zh: {tk_crr: '确认 {n} 个退回的包裹', tk_crrW: '快递员送回了无法送达的订单。请检查并确认，否则 Phsar 将视为已收到。', tk_crrB: '查看',
    c_delivery_failed: '配送失败——已退回给您', crrTitle: '正在退回的包裹', crrHistory: '已退回', crrBack: '← 订单',
    crrs_returning: '退回途中', crrs_awaiting_seller: '请确认', crrs_received: '已收到',
    crrW_returning: '快递员无法送达此订单，正在送回 {p}。快递员：{c}。',
    crrW_awaiting_seller: '快递员表示包裹已送回 {p}。请检查后确认。',
    crrOk: '已收到——放回库存', crrDamaged: '已收到，但有损坏', crrOkQ: '确认已收到完好的包裹？商品将放回库存，订单将取消。',
    crrDmgTitle: '包裹损坏？', crrDmgHint: '商品不会放回库存。订单将取消，买家获得退款。', crrDmgNote: '哪里损坏了？', crrDmgGo: '确认损坏',
    crrDone: '谢谢——已记录。', crrConfirmIn: '{t} 内确认', crrLate: '已超时——Phsar 将记录为已收到', crrEmpty: '没有退回的包裹。',
    crrR_seller: '您已确认', crrR_admin: 'Phsar 已记录', crrR_system: '自动记录', crrCond_ok: '已放回库存', crrCond_damaged: '损坏，未放回库存',
    crrMoney: '买家获得全额退款。退回运费不向您收取。'},
});

const SCR = {history: false};

// "Hari ini": buka Pesanan (kartu paket kembali ada di atas).
function crSellerOpen() { SCR.history = false; v2Open('orders'); }

function crCard(r) {
  const items = r.items.map(it => `<div class="oitem"><div><b>${esc(it.name)}</b>${it.options ? `<div class="muted">${esc(it.options)}</div>` : ''}</div><div class="oqty">× ${it.qty}</div></div>`).join('');
  const open = r.status !== 'received';
  const left = r.status === 'awaiting_seller' && r.confirm_by ? leftText(r.confirm_by) : null;
  const chipHTML = r.status === 'received' ? `<span class="chip ok">${t('crrs_received')}</span>`
    : r.status === 'awaiting_seller' ? `<span class="chip ${left ? 'warn' : 'bad'}">${left ? t('crrConfirmIn', {t: left}) : t('crrLate')}</span>`
    : `<span class="chip">${t('crrs_returning')}</span>`;
  const outcome = r.status === 'received' ? `<p class="ohelp">${t('crrR_' + r.resolved_by)} · ${fmtTime(r.received_at)} · ${t('crrCond_' + (r.condition || 'ok'))}${r.seller_note ? ' — ' + esc(r.seller_note) : ''}</p>` : '';
  return `<section class="card order"><div class="spread"><div><h2>↩️ ${t('orderNo')} ${esc(r.order_no)}</h2><span class="muted">${fmtTime(r.started_at)}</span></div>${chipHTML}</div>
    ${open ? `<p class="ohelp">${t('crrW_' + r.status, {p: esc(r.pickup_label), c: esc(r.courier_name)})}</p>` : ''}
    <div class="oitems">${items}</div>${outcome}
    ${open ? `<p class="muted small">${t('crrMoney')}</p>
      <button type="button" class="primary bigbtn" data-crrok="${r.id}">${t('crrOk')}</button>
      <div style="text-align:center"><button type="button" class="textbtn" data-crrdmg="${r.id}">${t('crrDamaged')}</button></div>` : ''}</section>`;
}

// Dipanggil tabOrders (seller-orders.js) setelah daftar pesanan digambar.
async function crAfterOrders(id) {
  const tab = $('#tab');
  if (!tab || S.tab !== 'orders') return;
  let data;
  try { data = await api(`/v1/seller/stores/${id}/courier-returns${SCR.history ? '?view=all' : ''}`); } catch { return; }
  if (S.tab !== 'orders' || $('#tab') !== tab || tab.querySelector('#crrHistBtn, #crrBackBtn')) return; // sudah digambar oleh panggilan lain
  if (SCR.history) {
    tab.innerHTML = `<a class="backlink" href="#" id="crrBackBtn">${t('crrBack')}</a><div class="v2head"><h1>${t('crrHistory')}</h1></div>
      ${data.returns.map(crCard).join('') || `<section class="card"><p class="muted" style="text-align:center">${t('crrEmpty')}</p></section>`}`;
    $('#crrBackBtn').onclick = e => { e.preventDefault(); SCR.history = false; tabOrders(id); };
  } else {
    if (data.returns.length) tab.insertAdjacentHTML('afterbegin', `<h2>${t('crrTitle')}</h2>${data.returns.map(crCard).join('')}`);
    const hist = tab.querySelector('.ohist');
    if (hist) hist.insertAdjacentHTML('beforeend', `<button type="button" id="crrHistBtn">${t('crrHistory')}</button>`);
    if ($('#crrHistBtn')) $('#crrHistBtn').onclick = () => { SCR.history = true; tabOrders(id); };
  }
  tab.querySelectorAll('[data-crrok]').forEach(b => b.onclick = async () => {
    if (!confirm(t('crrOkQ'))) return;
    crConfirm(id, b.dataset.crrok, {condition: 'ok'}).catch(() => {});
  });
  tab.querySelectorAll('[data-crrdmg]').forEach(b => b.onclick = () => crDamagedDialog(id, b.dataset.crrdmg));
}

async function crConfirm(storeId, retId, body) {
  try {
    await api(`/v1/seller/stores/${storeId}/courier-returns/${retId}/confirm`, {method: 'POST', body});
    toast(t('crrDone'));
    V2.today = await api(`/v1/seller/stores/${storeId}/today`).catch(() => V2.today);
    tabOrders(storeId);
    return true;
  } catch (e) {
    toast(errMsg(e), 'error');
    if (e.code === 'bad_status') tabOrders(storeId);
    throw e;
  }
}

function crDamagedDialog(storeId, retId) {
  const d = document.createElement('dialog');
  d.className = 'dlg';
  d.innerHTML = `<form method="dialog"><h2>${t('crrDmgTitle')}</h2><p class="muted">${t('crrDmgHint')}</p>
    <label>${t('crrDmgNote')}<textarea id="crrNote" maxlength="300"></textarea></label>
    <div class="row"><button type="button" class="danger" id="crrGo">${t('crrDmgGo')}</button><button value="cancel">${t('keep')}</button></div></form>`;
  document.body.append(d);
  d.showModal();
  d.addEventListener('close', () => d.remove());
  $('#crrGo', d).onclick = async () => {
    try { await crConfirm(storeId, retId, {condition: 'damaged', note: $('#crrNote', d).value}); d.remove(); }
    catch (e) { if (e.code === 'validation') showError(e, {note: '#crrNote'}); }
  };
}
