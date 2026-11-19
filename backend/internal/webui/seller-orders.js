'use strict';
// Seller Center: pesanan — terima / tolak, siap dijemput, cetak label paket.
addStrings({
  en: {os_new: 'Accept', os_processing: 'Pack', os_ready: 'Hand over', osN: '{n} orders', ordHistory: 'History:', cantSend: 'I can’t send this order', readyBig: 'Packed — ready for pickup',
    oh_new: 'Step 1: check you have the items, then tap Accept. Accept within {a} or the order cancels automatically.', oh_processing: 'Step 2: pack the items, print the label and stick it on the package, then tap “Packed — ready for pickup”.{p}', oh_ready: 'Step 3: the Phsar courier comes to your pickup point. Hand over the package — nothing else to do. Your money starts after delivery.',
    oe_new: 'No new orders right now.', oe_processing: 'Nothing to pack.', oe_ready: 'No packages waiting for the courier.', moneyDetail: 'How is this calculated?', tabOrders: 'Orders', of_new: 'New', of_processing: 'Preparing', of_ready: 'Ready for pickup', of_shipping: 'With courier', of_done: 'Completed', of_cancelled: 'Cancelled',
    noOrders: 'No orders here.', orderNo: 'Order', buyer: 'Buyer', items: 'Items', subtotal: 'Subtotal', commission: 'Phsar commission ({r}%)', youGet: 'You receive',
    respondIn: 'Accept within {t}', overdue: 'Time is up — this order will be cancelled', accept: 'Accept order', rejectOrder: 'Reject', markReady: 'Ready for pickup',
    printLabel: 'Print label', waitingCourier: 'Waiting for the Phsar courier.', pickupAt: 'Pickup point', choosePoint: 'Pickup point',
    rejectTitle: 'Reject this order?', rejectHint: 'The buyer is refunded in full. Rejecting often lowers your store rating.', reason: 'Reason',
    r_out_of_stock: 'Out of stock', r_cannot_fulfil: 'Cannot fulfil the order', r_price_error: 'Wrong price', r_other: 'Other', rejectNote: 'Explain (required for “Other”)',
    confirmReject: 'Reject order', keep: 'Keep order', accepted: 'Order accepted — prepare the items.', readyDone: 'Marked ready. The courier will come to {p}.',
    rejected: 'Order rejected. The buyer will be refunded.', paidAt: 'Paid', cancelledBy: 'Cancelled', c_seller_rejected: 'Rejected by store', c_seller_timeout: 'Not accepted in time', c_pack_timeout: 'Not packed in time', packIn: 'Ready for pickup within {t}', packLate: 'Pack time is up', ohPack: ' Do it within {p} or the order cancels automatically.',
    refundPending: 'Refund pending', refunded: 'Refunded', sku: 'SKU', labelTitle: 'Phsar package', labelFrom: 'From', labelItems: 'Contents', labelNote: 'Hand this package to the Phsar courier. Do not seal before checking the items.',
    err_bad_status: 'This order was already updated. Refreshing…', err_no_pickup_point: 'Add a pickup point first (Pickup points tab).', newOrdersBadge: '{n} new'},
  km: {os_new: 'ទទួល', os_processing: 'វេចខ្ចប់', os_ready: 'ប្រគល់', osN: '{n} ការបញ្ជាទិញ', ordHistory: 'ប្រវត្តិ៖', cantSend: 'ខ្ញុំមិនអាចផ្ញើការបញ្ជាទិញនេះបាន', readyBig: 'បានវេចខ្ចប់ — រួចរាល់ឱ្យមកយក',
    oh_new: 'ជំហាន 1៖ ពិនិត្យថាមានទំនិញ រួចចុច ទទួល។ ទទួលក្នុង {a} បើមិនដូច្នេះនឹងបោះបង់ដោយស្វ័យប្រវត្តិ។', oh_processing: 'ជំហាន 2៖ វេចខ្ចប់ទំនិញ បោះពុម្ពស្លាក ហើយបិទលើកញ្ចប់ រួចចុច “បានវេចខ្ចប់ — រួចរាល់ឱ្យមកយក”។{p}', oh_ready: 'ជំហាន 3៖ អ្នកដឹក Phsar មកទីតាំងទទួលរបស់អ្នក។ ប្រគល់កញ្ចប់ — មិនចាំបាច់ធ្វើអ្វីទៀតទេ។ លុយចាប់ផ្តើមបន្ទាប់ពីដឹកដល់។',
    oe_new: 'មិនទាន់មានការបញ្ជាទិញថ្មីទេ។', oe_processing: 'គ្មានអ្វីត្រូវវេចខ្ចប់។', oe_ready: 'គ្មានកញ្ចប់រង់ចាំអ្នកដឹកទេ។', moneyDetail: 'គណនាយ៉ាងដូចម្តេច?', tabOrders: 'ការបញ្ជាទិញ', of_new: 'ថ្មី', of_processing: 'កំពុងរៀបចំ', of_ready: 'រួចរាល់ឱ្យមកយក', of_shipping: 'នៅជាមួយអ្នកដឹក', of_done: 'បានបញ្ចប់', of_cancelled: 'បានបោះបង់',
    noOrders: 'គ្មានការបញ្ជាទិញនៅទីនេះទេ។', orderNo: 'ការបញ្ជាទិញ', buyer: 'អ្នកទិញ', items: 'ទំនិញ', subtotal: 'សរុបរង', commission: 'កម្រៃ Phsar ({r}%)', youGet: 'អ្នកទទួលបាន',
    respondIn: 'ទទួលក្នុងរយៈពេល {t}', overdue: 'ផុតពេលហើយ — ការបញ្ជាទិញនេះនឹងត្រូវបោះបង់', accept: 'ទទួលការបញ្ជាទិញ', rejectOrder: 'បដិសេធ', markReady: 'រួចរាល់ឱ្យមកយក',
    printLabel: 'បោះពុម្ពស្លាក', waitingCourier: 'កំពុងរង់ចាំអ្នកដឹក Phsar។', pickupAt: 'ទីតាំងទទួល', choosePoint: 'ទីតាំងទទួល',
    rejectTitle: 'បដិសេធការបញ្ជាទិញនេះ?', rejectHint: 'អ្នកទិញនឹងទទួលប្រាក់វិញទាំងស្រុង។ ការបដិសេធញឹកញាប់ធ្វើឱ្យពិន្ទុហាងធ្លាក់ចុះ។', reason: 'មូលហេតុ',
    r_out_of_stock: 'អស់ស្តុក', r_cannot_fulfil: 'មិនអាចផ្គត់ផ្គង់បាន', r_price_error: 'តម្លៃខុស', r_other: 'ផ្សេងៗ', rejectNote: 'ពន្យល់ (ចាំបាច់សម្រាប់ “ផ្សេងៗ”)',
    confirmReject: 'បដិសេធការបញ្ជាទិញ', keep: 'រក្សាទុក', accepted: 'បានទទួល — សូមរៀបចំទំនិញ។', readyDone: 'បានសម្គាល់ថារួចរាល់។ អ្នកដឹកនឹងមក {p}។',
    rejected: 'បានបដិសេធ។ អ្នកទិញនឹងទទួលប្រាក់វិញ។', paidAt: 'បានបង់', cancelledBy: 'បានបោះបង់', c_seller_rejected: 'ហាងបានបដិសេធ', c_seller_timeout: 'មិនបានទទួលទាន់ពេល', c_pack_timeout: 'មិនបានវេចខ្ចប់ទាន់ពេល', packIn: 'រួចរាល់ឱ្យមកយកក្នុង {t}', packLate: 'ផុតពេលវេចខ្ចប់', ohPack: ' ធ្វើក្នុង {p} បើមិនដូច្នេះការបញ្ជាទិញនឹងបោះបង់ដោយស្វ័យប្រវត្តិ។',
    refundPending: 'រង់ចាំសងប្រាក់', refunded: 'បានសងប្រាក់', sku: 'SKU', labelTitle: 'កញ្ចប់ Phsar', labelFrom: 'ពី', labelItems: 'ទំនិញក្នុងកញ្ចប់', labelNote: 'ប្រគល់កញ្ចប់នេះឱ្យអ្នកដឹក Phsar។ កុំបិទកញ្ចប់មុនពិនិត្យទំនិញ។',
    err_bad_status: 'ការបញ្ជាទិញនេះបានផ្លាស់ប្តូររួចហើយ។ កំពុងផ្ទុកឡើងវិញ…', err_no_pickup_point: 'សូមបន្ថែមទីតាំងទទួលជាមុនសិន។', newOrdersBadge: 'ថ្មី {n}'},
  zh: {os_new: '接单', os_processing: '打包', os_ready: '交件', osN: '{n} 笔', ordHistory: '历史：', cantSend: '我无法发出此订单', readyBig: '已打包 — 待取件',
    oh_new: '第 1 步：确认有货，然后点“接单”。请在 {a} 内接单，否则订单会自动取消。', oh_processing: '第 2 步：打包商品，打印面单并贴在包裹上，然后点“已打包 — 待取件”。{p}', oh_ready: '第 3 步：Phsar 快递员会到你的取件点。交出包裹即可，无需其他操作。送达后开始结算。',
    oe_new: '暂无新订单。', oe_processing: '暂无需要打包的订单。', oe_ready: '暂无等待取件的包裹。', moneyDetail: '如何计算？', tabOrders: '订单', of_new: '新订单', of_processing: '备货中', of_ready: '待取件', of_shipping: '配送中', of_done: '已完成', of_cancelled: '已取消',
    noOrders: '暂无订单。', orderNo: '订单', buyer: '买家', items: '商品', subtotal: '小计', commission: 'Phsar 佣金（{r}%）', youGet: '你将收到',
    respondIn: '请在 {t} 内接单', overdue: '已超时——该订单将被取消', accept: '接单', rejectOrder: '拒绝', markReady: '已备好待取件',
    printLabel: '打印面单', waitingCourier: '等待 Phsar 快递员取件。', pickupAt: '取件点', choosePoint: '取件点',
    rejectTitle: '拒绝此订单？', rejectHint: '买家将获得全额退款。经常拒单会降低店铺评分。', reason: '原因',
    r_out_of_stock: '缺货', r_cannot_fulfil: '无法履约', r_price_error: '价格错误', r_other: '其他', rejectNote: '说明（选择“其他”时必填）',
    confirmReject: '拒绝订单', keep: '保留订单', accepted: '已接单，请备货。', readyDone: '已标记为待取件，快递员将前往 {p}。',
    rejected: '已拒单，买家将获得退款。', paidAt: '付款时间', cancelledBy: '已取消', c_seller_rejected: '店铺拒单', c_seller_timeout: '未及时接单', c_pack_timeout: '未及时打包', packIn: '{t} 内标记待取件', packLate: '打包已超时', ohPack: ' 请在 {p} 内完成，否则订单会自动取消。',
    refundPending: '待退款', refunded: '已退款', sku: 'SKU', labelTitle: 'Phsar 包裹', labelFrom: '发货方', labelItems: '包裹内容', labelNote: '请将包裹交给 Phsar 快递员。核对商品后再封箱。',
    err_bad_status: '该订单已被更新，正在刷新…', err_no_pickup_point: '请先添加取件点。', newOrdersBadge: '{n} 个新订单'},
});

const SO = {filter: 'new', data: null, timer: null, q: '', page: 1, from: '', to: ''};
addStrings({
  en: {oSearchPh: 'Order no., buyer, product, SKU', oSearchHint: 'Searches all your orders, any status.'},
  km: {oSearchPh: 'លេខបញ្ជាទិញ អ្នកទិញ ផលិតផល SKU', oSearchHint: 'ស្វែងរកគ្រប់ការបញ្ជាទិញ គ្រប់ស្ថានភាព។'},
  zh: {oSearchPh: '订单号、买家、商品、SKU', oSearchHint: '在全部订单中搜索（任何状态）。'},
});
const OFILTERS = ['new', 'processing', 'ready', 'shipping', 'done', 'cancelled'];

function omoney(v) { const n = parseFloat(v); return '$' + (n === Math.round(n) ? n : n.toFixed(2)); }
function fmtTime(s) {
  try { return new Date(s).toLocaleString(P.lang === 'zh' ? 'zh-CN' : P.lang === 'km' ? 'km-KH' : 'en-GB', {dateStyle: 'medium', timeStyle: 'short'}); } catch { return s; }
}
function leftText(iso) {
  const ms = new Date(iso) - Date.now();
  if (ms <= 0) return null;
  const m = Math.floor(ms / 60000);
  return m >= 60 ? `${Math.floor(m / 60)}h ${m % 60}m` : `${m}m`;
}

function ordHelp(f) {
  const tm = V2.today?.timing || {accept_minutes: 120, pack_enabled: false};
  return t('oh_' + f, {a: v2Dur(tm.accept_minutes, true), p: tm.pack_enabled ? t('ohPack', {p: v2Dur(tm.pack_minutes, true)}) : ''});
}
const OSTEPS = ['new', 'processing', 'ready'], OHIST = ['shipping', 'done', 'cancelled'];
async function tabOrders(id) {
  clearInterval(SO.timer);
  const searching = !!SO.q, hist = OHIST.includes(SO.filter);
  const dates = searching || hist ? {from: SO.from, to: SO.to} : {};
  SO.data = await api(`/v1/seller/stores/${id}/orders?${qs({filter: SO.filter, q: SO.q, page: SO.page, size: 20, ...dates})}`);
  const pg = SO.data.page;
  if (pg && SO.page > pg.pages) { SO.page = pg.pages; return tabOrders(id); }
  const c = SO.data.counts || {};
  const stepTabs = OSTEPS.map((f, i) => `<button type="button" role="tab" aria-selected="${!searching && SO.filter === f}" data-of="${f}">
      <span class="sn" aria-hidden="true">${i + 1}</span><b>${t('os_' + f)}</b><span class="sc">${c[f] ? (f === 'new' ? `<span class="badge2">${c[f]}</span>` : t('osN', {n: c[f]})) : t('osN', {n: 0})}</span></button>`).join('');
  const histBtns = OHIST.map(f => `<button type="button" aria-selected="${!searching && SO.filter === f}" data-of="${f}">${t('of_' + f)}${c[f] ? ` (${c[f]})` : ''}</button>`).join('');
  const list = SO.data.orders.map(o => orderCard(o)).join('');
  const isStep = !searching && OSTEPS.includes(SO.filter);
  const dateBar = searching || hist ? `<div class="dates">
      <label>${t('dateFrom')}<input type="date" id="oFrom" value="${esc(SO.from)}" max="${esc(SO.to || '')}"></label>
      <label>${t('dateTo')}<input type="date" id="oTo" value="${esc(SO.to)}" min="${esc(SO.from || '')}"></label>
      ${SO.from || SO.to ? `<button type="button" class="small" id="oDates">${t('dateClear')}</button>` : ''}</div>` : '';
  const empty = searching ? t('lsNoMatch', {q: esc(SO.q)}) : isStep ? t('oe_' + SO.filter) : t('noOrders');
  $('#tab').innerHTML = `<div class="listbar">${searchBox('oQ', SO.q, t('oSearchPh'))}</div>
    ${searching ? `<div class="lsresult"><b>${t('lsResults', {q: esc(SO.q)})}</b><span class="muted small">${t('oSearchHint')}</span></div>`
      : `<div class="otabs" role="tablist" aria-label="${t('h_orders')}">${stepTabs}</div>`}
    ${isStep ? `<p class="ohelp">${ordHelp(SO.filter)}</p>` : ''}${typeof closedBanner === 'function' ? closedBanner() : ''}
    ${dateBar}
    ${list || `<section class="card"><p class="muted" style="text-align:center">${empty}</p></section>`}
    ${pagerHTML(pg)}
    ${searching ? '' : `<div class="ohist"><span>${t('ordHistory')}</span>${histBtns}</div>`}`;
  document.querySelectorAll('[data-of]').forEach(b => b.onclick = () => { SO.filter = b.dataset.of; SO.page = 1; tabOrders(id); });
  document.querySelectorAll('[data-oact]').forEach(b => b.onclick = () => orderAction(id, b.dataset.no, b.dataset.oact));
  bindSearch('oQ', v => { if (v === SO.q) return; SO.q = v; SO.page = 1; tabOrders(id); });
  bindPager($('#tab'), p => { SO.page = p; tabOrders(id); });
  const setDates = () => { SO.from = $('#oFrom')?.value || ''; SO.to = $('#oTo')?.value || ''; SO.page = 1; tabOrders(id); };
  if ($('#oFrom')) { $('#oFrom').onchange = setDates; $('#oTo').onchange = setDates; }
  if ($('#oDates')) $('#oDates').onclick = () => { SO.from = SO.to = ''; SO.page = 1; tabOrders(id); };
  if (typeof crAfterOrders === 'function' && !searching) crAfterOrders(id); // courierreturn: paket gagal antar yang kembali ke toko
  // Hitung mundur batas 2 jam + muat ulang berkala supaya pesanan baru muncul.
  let ticks = 0;
  SO.timer = setInterval(() => {
    if (!$('#tab') || S.tab !== 'orders') { clearInterval(SO.timer); return; }
    document.querySelectorAll('[data-left]').forEach(el => {
      const l = leftText(el.dataset.left);
      el.textContent = l ? t('respondIn', {t: l}) : t('overdue');
      el.classList.toggle('bad', !l);
    });
    document.querySelectorAll('[data-pack]').forEach(el => {
      const l = leftText(el.dataset.pack);
      el.textContent = l ? t('packIn', {t: l}) : t('packLate');
      el.classList.toggle('bad', !l);
    });
    // Muat ulang hanya di tab langkah, saat tidak sedang mengetik/memilih (fokus & isian tidak hilang).
    const busy = document.querySelector('dialog[open]') || ['oQ', 'oFrom', 'oTo'].includes(document.activeElement?.id);
    if (++ticks % 6 === 0 && isStep && !busy) tabOrders(id);
  }, 10000);
}

function orderCard(o) {
  const items = o.items.map(it => `<div class="oitem"><img src="${esc(it.image_url)}" alt="" onerror="this.style.visibility='hidden'">
    <div><b>${esc(it.name)}</b>${it.options ? `<div class="muted">${esc(it.options)}</div>` : ''}<div class="muted small">${t('sku')} ${esc(it.sku)}</div></div>
    <div class="oqty">× ${it.qty}<div class="muted">${omoney(it.line_total)}</div></div></div>`).join('');
  const left = o.status === 'paid' && o.respond_by ? leftText(o.respond_by) : null;
  const statusLine = o.status === 'paid' ? `<span class="chip ${left ? 'warn' : 'bad'}" data-left="${esc(o.respond_by)}">${left ? t('respondIn', {t: left}) : t('overdue')}</span>`
    : o.status === 'cancelled' || o.status === 'refunded' ? `<span class="chip bad">${t(o.cancel_reason ? 'c_' + o.cancel_reason : 'cancelledBy')}</span>`
    : o.status === 'ready_for_pickup' ? `<span class="chip ok">${t('of_ready')}</span>` : o.status === 'accepted' && o.ready_by ? `<span class="chip ${leftText(o.ready_by) ? 'warn' : 'bad'}" data-pack="${esc(o.ready_by)}">${leftText(o.ready_by) ? t('packIn', {t: leftText(o.ready_by)}) : t('packLate')}</span>`
    : o.status === 'accepted' ? `<span class="chip pink">${t('of_processing')}</span>`
    : o.status === 'picked_up' ? `<span class="chip">${t('of_shipping')}</span>` : '';
  const pts = (S.points || []).filter(p => !p.deleted_at);
  const pointSelect = pts.length > 1 ? `<label>${t('choosePoint')}<select id="pp_${esc(o.order_no)}">${pts.map(p => `<option value="${p.id}" ${p.is_default ? 'selected' : ''}>${esc(p.label)}</option>`).join('')}</select></label>` : '';
  const btn = (act, label, cls = '') => `<button type="button" class="${cls}" data-oact="${act}" data-no="${esc(o.order_no)}">${label}</button>`;
  const actions = o.status === 'paid' ? btn('accept', t('accept'), 'primary bigbtn') + `<div style="text-align:center">${btn('reject', t('cantSend'), 'textbtn')}</div>`
    : o.status === 'accepted' ? pointSelect + `<div class="row">${btn('label', '🖨 ' + t('printLabel'))}</div>` + btn('ready', t('readyBig'), 'primary bigbtn') + `<div style="text-align:center">${btn('reject', t('cantSend'), 'textbtn')}</div>`
    : o.status === 'ready_for_pickup' ? `<p class="ohelp">🛵 ${t('waitingCourier')} ${o.pickup_point ? `· ${t('pickupAt')}: <b>${esc(o.pickup_point)}</b>` : ''}</p><div class="row">${btn('label', '🖨 ' + t('printLabel'))}</div>`
    : '';
  return `<section class="card order">
    <div class="spread"><div><h2>${t('orderNo')} ${esc(o.order_no)}</h2><span class="muted">${t('buyer')}: ${esc(o.buyer_name)} · ${fmtTime(o.paid_at)}</span></div>${statusLine}</div>
    <div class="oitems">${items}</div>
    <div class="getline"><span>${t('youGet')}</span><b>${omoney(o.seller_amount)}</b></div>
    <details><summary>${t('moneyDetail')}</summary><dl class="omoney"><div><dt>${t('subtotal')}</dt><dd>${omoney(o.subtotal)}</dd></div>
      <div><dt>${t('commission', {r: parseFloat(o.commission_rate)})}</dt><dd>−${omoney(o.commission_amount)}</dd></div>
      <div class="total"><dt>${t('youGet')}</dt><dd>${omoney(o.seller_amount)}</dd></div></dl></details>
    ${o.seller_note ? `<p class="muted">${esc(o.seller_note)}</p>` : ''}
    ${actions}</section>`;
}

async function orderAction(storeId, no, act) {
  const url = `/v1/seller/stores/${storeId}/orders/${encodeURIComponent(no)}`;
  try {
    if (act === 'accept') { await api(url + '/accept', {method: 'POST'}); toast(t('accepted')); }
    else if (act === 'ready') {
      const sel = document.getElementById('pp_' + no);
      const o = await api(url + '/ready', {method: 'POST', body: sel ? {pickup_point_id: sel.value} : {}});
      toast(t('readyDone', {p: o.pickup_point}));
    } else if (act === 'label') { printLabel(SO.data.orders.find(o => o.order_no === no)); return; }
    else if (act === 'reject') { if (!(await rejectDialog(url))) return; toast(t('rejected')); }
    await tabOrders(storeId);
    loadBell();
  } catch (e) {
    toast(errMsg(e), 'error');
    if (e.code === 'bad_status') tabOrders(storeId);
  }
}

function rejectDialog(url) {
  return new Promise(resolve => {
    const d = document.createElement('dialog');
    d.className = 'dlg';
    d.innerHTML = `<form method="dialog"><h2>${t('rejectTitle')}</h2><p class="muted">${t('rejectHint')}</p>
      <label>${t('reason')}<select id="rjReason">${['out_of_stock', 'cannot_fulfil', 'price_error', 'other'].map(r => `<option value="${r}">${t('r_' + r)}</option>`).join('')}</select></label>
      <label>${t('rejectNote')}<textarea id="rjNote" maxlength="300"></textarea></label>
      <div class="row"><button type="button" class="danger" id="rjGo">${t('confirmReject')}</button><button value="cancel">${t('keep')}</button></div></form>`;
    document.body.append(d);
    d.showModal();
    d.addEventListener('close', () => { d.remove(); resolve(false); });
    d.querySelector('#rjGo').onclick = async () => {
      try {
        await api(url + '/reject', {method: 'POST', body: {reason: $('#rjReason', d).value, note: $('#rjNote', d).value}});
        d.remove(); resolve(true);
      } catch (e) { showError(e, {reason: '#rjReason', note: '#rjNote'}); }
    };
  });
}

// ---- Label paket (dicetak di printer biasa / thermal) ----
// Barcode Code 128 (set B) untuk nomor pesanan: dipindai kurir saat menjemput.
const C128 = ['212222', '222122', '222221', '121223', '121322', '131222', '122213', '122312', '132212', '221213', '221312', '231212', '112232', '122132', '122231', '113222', '123122', '123221', '223211', '221132', '221231', '213212', '223112', '312131', '311222', '321122', '321221', '312212', '322112', '322211', '212123', '212321', '232121', '111323', '131123', '131321', '112313', '132113', '132311', '211313', '231113', '231311', '112133', '112331', '132131', '113123', '113321', '133121', '313121', '211331', '231131', '213113', '213311', '213131', '311123', '311321', '331121', '312113', '312311', '332111', '314111', '221411', '431111', '111224', '111422', '121124', '121421', '141122', '141221', '112214', '112412', '122114', '122411', '142112', '142211', '241211', '221114', '413111', '241112', '134111', '111242', '121142', '121241', '114212', '124112', '124211', '411212', '421112', '421211', '212141', '214121', '412121', '111143', '111341', '131141', '114113', '114311', '411113', '411311', '113141', '114131', '311141', '411131', '211412', '211214', '211232', '2331112'];
function code128SVG(text, h = 70) {
  const codes = [104];
  for (const ch of text) codes.push(ch.charCodeAt(0) - 32);
  let sum = 104;
  for (let i = 1; i < codes.length; i++) sum += codes[i] * i;
  codes.push(sum % 103, 106);
  let x = 10, bars = '';
  for (const c of codes) {
    const p = C128[c];
    for (let i = 0; i < p.length; i++) {
      const w = +p[i] * 2;
      if (i % 2 === 0) bars += `<rect x="${x}" y="0" width="${w}" height="${h}"/>`;
      x += w;
    }
  }
  return `<svg xmlns="http://www.w3.org/2000/svg" width="${x + 10}" height="${h}" viewBox="0 0 ${x + 10} ${h}" shape-rendering="crispEdges" role="img" aria-label="${esc(text)}"><g fill="#000">${bars}</g></svg>`;
}

function printLabel(o) {
  const items = o.items.map(it => `<tr><td>${it.qty} ×</td><td>${esc(it.name)}${it.options ? ` <small>(${esc(it.options)})</small>` : ''}<br><small>${esc(it.sku)}</small></td></tr>`).join('');
  const w = window.open('', '_blank', 'width=480,height=720');
  if (!w) { toast('Pop-up blocked', 'warning'); return; }
  w.document.write(`<!doctype html><html lang="${P.lang}"><head><meta charset="utf-8"><title>${esc(o.order_no)}</title><style>
    @page{size:100mm 150mm;margin:4mm}*{box-sizing:border-box}body{font-family:"DM Sans","Noto Sans Khmer","Noto Sans SC",Arial,sans-serif;margin:0;color:#000}
    .lbl{width:92mm;border:2px solid #000;border-radius:3mm;padding:4mm}.hd{display:flex;justify-content:space-between;align-items:center;border-bottom:2px solid #000;padding-bottom:2mm}
    .hd b{font-size:20pt;letter-spacing:.02em}.no{font-size:22pt;font-weight:800;text-align:center;margin:3mm 0 1mm;letter-spacing:.04em}.bc{text-align:center}.bc svg{max-width:100%;height:18mm}
    h3{font-size:9pt;text-transform:uppercase;margin:3mm 0 1mm;letter-spacing:.06em}p{margin:0;font-size:10pt}table{width:100%;font-size:10pt;border-collapse:collapse}td{vertical-align:top;padding:1mm 0}
    td:first-child{width:12mm;font-weight:700}small{color:#333}.note{font-size:8pt;margin-top:3mm;border-top:1px dashed #000;padding-top:2mm}
    @media screen{body{padding:16px;background:#eee}.lbl{background:#fff;margin:auto}}</style></head><body><div class="lbl">
    <div class="hd"><b>PHSAR</b><span>${esc(t('labelTitle'))}</span></div>
    <div class="no">${esc(o.order_no)}</div><div class="bc">${code128SVG(o.order_no)}</div>
    <h3>${esc(t('labelFrom'))}</h3><p><b>${esc(S.store?.name || '')}</b>${o.pickup_point ? ' · ' + esc(o.pickup_point) : ''}</p>
    <h3>${esc(t('labelItems'))} (${o.item_count})</h3><table>${items}</table>
    <p class="note">${esc(t('labelNote'))}</p></div><script>onload=()=>setTimeout(()=>print(),300)<\/script></body></html>`);
  w.document.close();
}
