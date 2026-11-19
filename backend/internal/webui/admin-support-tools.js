'use strict';
// Admin CS: cari pelanggan & pesanan, detail pesanan lengkap, catatan CS, tangguhkan akun, log pembayaran PayWay.
addStrings({
  en: {csSearchPh: 'Phone, name, @telegram, order no., return no. or PayWay ID', csSearch: 'Search', csHint: 'Find a customer by anything they tell you. Opening a customer is recorded in the audit log.',
    csOrdersTab: 'Orders', csCustomersTab: 'Customers', csNoResult: 'Nothing found.', csGuest: 'Guest', csSeller: 'Seller', csCourier: 'Courier', csAdmin: 'Admin',
    csSuspended: 'Suspended', csActive: 'Active', csOrders: 'Orders', csSpent: 'Spent', csOpenReturns: 'Open returns', csJoined: 'Joined', csLastSeen: 'Last seen',
    csLogins: 'Sign-in methods', csSessions: '{n} active sessions', csStores: 'Stores', csAddresses: 'Delivery addresses', csNotes: 'CS notes', csAddNote: 'Add note',
    csNotePh: 'What happened? Only admins can see this.', csSuspend: 'Suspend account', csReactivate: 'Reactivate', csSuspendTitle: 'Suspend this account?',
    csSuspendHint: 'The customer is signed out everywhere and cannot sign in until reactivated. Paid orders continue.', csReason: 'Reason (required)',
    csCancelled: 'Cancelled paid', csReturns: 'Returns', csUnpaid: 'Unpaid checkouts', csNoNotes: 'No notes yet.', csKind_suspended: 'Suspended', csKind_reactivated: 'Reactivated', csKind_note: 'Note',
    csBack: 'Back', csAllStatuses: 'All statuses', odBuyer: 'Buyer', odStore: 'Store', odDeliverTo: 'Deliver to', odPayment: 'Payment', odItems: 'Items', odMoney: 'Money',
    odTimeline: 'Timeline', odDelivery: 'Delivery', odPayLog: 'PayWay log', odLedger: 'Seller balance entries', odReturns: 'Returns', odCommission: 'Commission ({r}%)',
    odSeller: 'Seller receives', odTran: 'PayWay transaction', odPaidAmt: 'Paid', odRef: 'APV / reference', odPayer: 'Payer account', odSplit: 'Part of a checkout with {n} orders (one KHQR payment).',
    odRefund: 'Refund', odCourier: 'Courier', odProof: 'Delivery photo', odNoDelivery: 'Not handed to a courier yet.', odOpenMap: 'Open map', odCustomer: 'Open customer',
    odCopy: 'Copy', odCopied: 'Copied', odPayload: 'Details', odVerified: 'signature OK', odUnverified: 'NOT verified',
    pyHint: 'Everything Phsar sent to or received from ABA PayWay. Problems = callbacks with a bad signature, failed/declined payments, or money received after the order expired.',
    pyAll: 'All', pyCallbacks: 'Callbacks', pyProblems: 'Problems', pyEvents24: 'Events (24 h)', pyPaid24: 'Paid checkouts (24 h)', pyProblems7: 'Problems (7 days)',
    pySearchPh: 'PayWay ID or order no.', odSubtotal: 'Subtotal', odShipping: 'Delivery fee', odPaidAt: 'Paid at', pyOrders: 'Orders', pyCheckout: 'Checkout'},
  km: {csSearchPh: 'ទូរស័ព្ទ ឈ្មោះ @telegram លេខបញ្ជាទិញ លេខប្ដូរ ឬលេខ PayWay', csSearch: 'ស្វែងរក', csHint: 'ស្វែងរកអតិថិជនតាមអ្វីដែលពួកគេប្រាប់អ្នក។ ការបើកមើលអតិថិជនត្រូវកត់ក្នុងកំណត់ហេតុសវនកម្ម។',
    csOrdersTab: 'ការបញ្ជាទិញ', csCustomersTab: 'អតិថិជន', csNoResult: 'រកមិនឃើញ។', csGuest: 'ភ្ញៀវ', csSeller: 'អ្នកលក់', csCourier: 'អ្នកដឹក', csAdmin: 'អ្នកគ្រប់គ្រង',
    csSuspended: 'បានផ្អាក', csActive: 'សកម្ម', csOrders: 'ការបញ្ជាទិញ', csSpent: 'ចំណាយ', csOpenReturns: 'ការប្ដូរកំពុងដំណើរការ', csJoined: 'ចូលរួម', csLastSeen: 'ឃើញចុងក្រោយ',
    csLogins: 'វិធីចូល', csSessions: 'វគ្គសកម្ម {n}', csStores: 'ហាង', csAddresses: 'អាសយដ្ឋានដឹកជញ្ជូន', csNotes: 'កំណត់ចំណាំ CS', csAddNote: 'បន្ថែមកំណត់ចំណាំ',
    csNotePh: 'មានអ្វីកើតឡើង? មានតែអ្នកគ្រប់គ្រងទេដែលមើលឃើញ។', csSuspend: 'ផ្អាកគណនី', csReactivate: 'បើកឡើងវិញ', csSuspendTitle: 'ផ្អាកគណនីនេះ?',
    csSuspendHint: 'អតិថិជនត្រូវចាកចេញពីគ្រប់ឧបករណ៍ ហើយមិនអាចចូលបានរហូតដល់បើកឡើងវិញ។ ការបញ្ជាទិញដែលបានបង់នៅតែដំណើរការ។', csReason: 'មូលហេតុ (ចាំបាច់)',
    csCancelled: 'បានបោះបង់ក្រោយបង់', csReturns: 'ការប្ដូរ', csUnpaid: 'មិនទាន់បង់', csNoNotes: 'មិនទាន់មានកំណត់ចំណាំ។', csKind_suspended: 'បានផ្អាក', csKind_reactivated: 'បានបើកឡើងវិញ', csKind_note: 'កំណត់ចំណាំ',
    csBack: 'ត្រឡប់', csAllStatuses: 'គ្រប់ស្ថានភាព', odBuyer: 'អ្នកទិញ', odStore: 'ហាង', odDeliverTo: 'ដឹកទៅ', odPayment: 'ការបង់ប្រាក់', odItems: 'ទំនិញ', odMoney: 'ប្រាក់',
    odTimeline: 'ប្រវត្តិ', odDelivery: 'ការដឹកជញ្ជូន', odPayLog: 'កំណត់ហេតុ PayWay', odLedger: 'សមតុល្យអ្នកលក់', odReturns: 'ការប្ដូរ', odCommission: 'កម្រៃ ({r}%)',
    odSeller: 'អ្នកលក់ទទួល', odTran: 'ប្រតិបត្តិការ PayWay', odPaidAmt: 'បានបង់', odRef: 'APV / លេខយោង', odPayer: 'គណនីអ្នកបង់', odSplit: 'ជាផ្នែកនៃការទូទាត់ដែលមាន {n} ការបញ្ជាទិញ (KHQR តែមួយ)។',
    odRefund: 'សងប្រាក់', odCourier: 'អ្នកដឹក', odProof: 'រូបថតដឹកដល់', odNoDelivery: 'មិនទាន់ប្រគល់ឱ្យអ្នកដឹក។', odOpenMap: 'បើកផែនទី', odCustomer: 'បើកអតិថិជន',
    odCopy: 'ចម្លង', odCopied: 'បានចម្លង', odPayload: 'ព័ត៌មានលម្អិត', odVerified: 'ហត្ថលេខាត្រឹមត្រូវ', odUnverified: 'មិនបានផ្ទៀងផ្ទាត់',
    pyHint: 'អ្វីៗទាំងអស់ដែល Phsar ផ្ញើទៅ ឬទទួលពី ABA PayWay។ បញ្ហា = callback ហត្ថលេខាខុស ការបង់ប្រាក់បរាជ័យ ឬទទួលប្រាក់ក្រោយការបញ្ជាទិញផុតកំណត់។',
    pyAll: 'ទាំងអស់', pyCallbacks: 'Callback', pyProblems: 'បញ្ហា', pyEvents24: 'ព្រឹត្តិការណ៍ (24 ម៉ោង)', pyPaid24: 'បានបង់ (24 ម៉ោង)', pyProblems7: 'បញ្ហា (7 ថ្ងៃ)',
    pySearchPh: 'លេខ PayWay ឬលេខបញ្ជាទិញ', odSubtotal: 'សរុបរង', odShipping: 'ថ្លៃដឹក', odPaidAt: 'បង់នៅ', pyOrders: 'ការបញ្ជាទិញ', pyCheckout: 'ការទូទាត់'},
  zh: {csSearchPh: '手机号、姓名、@telegram、订单号、退货号或 PayWay 交易号', csSearch: '搜索', csHint: '根据客户提供的任何信息查找。查看客户资料会记录在审计日志中。',
    csOrdersTab: '订单', csCustomersTab: '客户', csNoResult: '没有找到。', csGuest: '访客', csSeller: '卖家', csCourier: '快递员', csAdmin: '管理员',
    csSuspended: '已停用', csActive: '正常', csOrders: '订单', csSpent: '消费', csOpenReturns: '进行中的退货', csJoined: '注册', csLastSeen: '最近活动',
    csLogins: '登录方式', csSessions: '{n} 个活跃会话', csStores: '店铺', csAddresses: '收货地址', csNotes: '客服备注', csAddNote: '添加备注',
    csNotePh: '发生了什么？仅管理员可见。', csSuspend: '停用账户', csReactivate: '重新启用', csSuspendTitle: '停用此账户？',
    csSuspendHint: '客户将在所有设备上退出，重新启用前无法登录。已付款订单继续处理。', csReason: '原因（必填）',
    csCancelled: '付款后取消', csReturns: '退货', csUnpaid: '未付款结账', csNoNotes: '暂无备注。', csKind_suspended: '已停用', csKind_reactivated: '已重新启用', csKind_note: '备注',
    csBack: '返回', csAllStatuses: '全部状态', odBuyer: '买家', odStore: '店铺', odDeliverTo: '收货地址', odPayment: '付款', odItems: '商品', odMoney: '金额',
    odTimeline: '时间线', odDelivery: '配送', odPayLog: 'PayWay 日志', odLedger: '卖家余额记录', odReturns: '退货', odCommission: '佣金（{r}%）',
    odSeller: '卖家收入', odTran: 'PayWay 交易', odPaidAmt: '已付', odRef: 'APV / 参考号', odPayer: '付款账户', odSplit: '属于包含 {n} 个订单的一次结账（同一个 KHQR）。',
    odRefund: '退款', odCourier: '快递员', odProof: '送达照片', odNoDelivery: '尚未交给快递员。', odOpenMap: '打开地图', odCustomer: '打开客户',
    odCopy: '复制', odCopied: '已复制', odPayload: '详情', odVerified: '签名正确', odUnverified: '未验证',
    pyHint: 'Phsar 与 ABA PayWay 之间的所有往来记录。问题 = 签名错误的回调、失败/拒绝的付款，或订单过期后才收到的款项。',
    pyAll: '全部', pyCallbacks: '回调', pyProblems: '问题', pyEvents24: '事件（24 小时）', pyPaid24: '已付款结账（24 小时）', pyProblems7: '问题（7 天）',
    pySearchPh: 'PayWay 交易号或订单号', odSubtotal: '小计', odShipping: '运费', odPaidAt: '付款时间', pyOrders: '订单', pyCheckout: '结账'},
});

const CS = {q: '', tab: 'customers', status: ''};
const ODSTAT = ['paid', 'accepted', 'ready_for_pickup', 'picked_up', 'delivered', 'completed', 'cancelled', 'refunded'];
const ODCHIP = {paid: 'warn', accepted: 'pink', ready_for_pickup: 'pink', picked_up: 'pink', delivered: 'ok', completed: 'ok', cancelled: 'bad', refunded: 'bad', pending_payment: ''};
const csDt = d => d ? new Date(d).toLocaleString() : '—';
const csUsd = v => '$' + (v == null ? '0' : v);
function csSt(s) { const k = 'status_' + s; return t(k) !== k ? t(k) : s.replace(/_/g, ' '); }
function csOrdSt(s) { return `<span class="chip ${ODCHIP[s] || ''}">${esc(csSt(s))}</span>`; }

// ---------- Pelanggan & pesanan ----------
async function renderCustomers() {
  const id = (location.hash.match(/^#\/customers\/([0-9a-f-]{36})/) || [])[1];
  if (id) return renderCustomer(id);
  const canOrders = hasPerm('order.view'), canCust = hasPerm('customer.view');
  if (!canCust) CS.tab = 'orders';
  const q = encodeURIComponent(CS.q);
  let body = '';
  if (CS.tab === 'customers') {
    const {customers} = await api('/admin/api/customers?q=' + q);
    body = customers.length ? `<div class="tbl-wrap"><table><thead><tr><th>${t('odBuyer')}</th><th>${t('csOrders')}</th><th>${t('csSpent')}</th><th>${t('csOpenReturns')}</th><th>${t('csJoined')}</th></tr></thead><tbody>${customers.map(c => `
      <tr class="click" data-href="#/customers/${c.id}"><td><b>${esc(c.name || c.phone || '—')}</b>${c.name && c.phone ? `<br><span class="muted small">${esc(c.phone)}</span>` : ''}
        <div class="row" style="gap:4px;margin-top:4px">${c.status === 'suspended' ? `<span class="chip bad">${t('csSuspended')}</span>` : ''}${c.is_guest ? `<span class="chip">${t('csGuest')}</span>` : ''}${c.is_seller ? `<span class="chip pink">${t('csSeller')}</span>` : ''}${c.is_courier ? `<span class="chip">${t('csCourier')}</span>` : ''}${c.is_admin ? `<span class="chip warn">${t('csAdmin')}</span>` : ''}${(c.providers || []).map(p => `<span class="chip">${esc(p)}</span>`).join('')}</div></td>
        <td>${c.orders}</td><td>${csUsd(c.spent)}</td><td>${c.open_returns ? `<span class="chip warn">${c.open_returns}</span>` : '0'}</td><td>${csDt(c.created_at)}</td></tr>`).join('')}</tbody></table></div>`
      : `<p class="muted" style="padding:16px">${t('csNoResult')}</p>`;
  } else {
    const {orders} = await api(`/admin/api/orders?q=${q}&status=${CS.status}`);
    body = orders.length ? `<div class="tbl-wrap"><table><thead><tr><th>${t('rtOrder')}</th><th>${t('odBuyer')}</th><th>${t('odStore')}</th><th>${t('odMoney')}</th><th>Status</th></tr></thead><tbody>${orders.map(o => `
      <tr class="click" data-href="#/order/${esc(o.order_no)}"><td><b>${esc(o.order_no)}</b><br><span class="muted small">${csDt(o.created_at)}</span></td>
        <td>${esc(o.buyer_name || '—')}<br><span class="muted small">${esc(o.buyer_phone || '')}</span></td><td>${esc(o.store_name)}</td><td>${csUsd(o.total)}</td>
        <td>${csOrdSt(o.status)}${o.refund_status === 'pending' ? ` <span class="chip warn">${t('odRefund')}</span>` : ''}${o.return_status ? ` <span class="chip">↩ ${esc(o.return_status)}</span>` : ''}</td></tr>`).join('')}</tbody></table></div>`
      : `<p class="muted" style="padding:16px">${t('csNoResult')}</p>`;
  }
  const tabs = [canCust && 'customers', canOrders && 'orders'].filter(Boolean);
  $('#app').innerHTML = shell(`${pageHead(t('i_customers'), t('csHint'))}
    <form class="row" id="csForm" style="gap:10px"><input id="csQ" type="search" value="${esc(CS.q)}" placeholder="${t('csSearchPh')}" style="flex:1;min-width:220px" autocomplete="off">
      ${CS.tab === 'orders' ? `<select id="csStatus"><option value="">${t('csAllStatuses')}</option>${ODSTAT.map(s => `<option value="${s}" ${CS.status === s ? 'selected' : ''}>${esc(csSt(s))}</option>`).join('')}</select>` : ''}
      <button class="primary" type="submit">${t('csSearch')}</button></form>
    ${tabs.length > 1 ? `<div class="tabs">${tabs.map(k => `<button type="button" aria-selected="${CS.tab === k}" data-cst="${k}">${t(k === 'customers' ? 'csCustomersTab' : 'csOrdersTab')}</button>`).join('')}</div>` : ''}
    <section class="panel">${body}</section>`);
  $('#csForm').onsubmit = e => { e.preventDefault(); CS.q = $('#csQ').value.trim(); if ($('#csStatus')) CS.status = $('#csStatus').value; render(); };
  document.querySelectorAll('[data-cst]').forEach(b => b.onclick = () => { CS.tab = b.dataset.cst; render(); });
  document.querySelectorAll('tr[data-href]').forEach(tr => tr.onclick = () => { location.hash = tr.dataset.href; });
  $('#csQ').focus();
}

async function renderCustomer(id) {
  const c = await api('/admin/api/customers/' + id);
  const canAct = hasPerm('support.handle') && !c.is_admin;
  const s = c.stats;
  const kpi = (l, v, cls = '') => `<div class="kpi"><span>${l}</span><b class="${cls}" style="font-size:22px">${v}</b></div>`;
  $('#app').innerHTML = shell(`<a href="#/customers">← ${t('csBack')}</a>
    ${pageHead(esc(c.name || c.phone || '—'), `${esc(c.phone || '')}${c.email ? ' · ' + esc(c.email) : ''} · ${t('csJoined')} ${csDt(c.created_at)}${c.last_seen_at ? ' · ' + t('csLastSeen') + ' ' + csDt(c.last_seen_at) : ''}`,
      canAct ? `<button type="button" class="${c.status === 'suspended' ? 'primary' : 'danger'}" id="csStatusBtn">${t(c.status === 'suspended' ? 'csReactivate' : 'csSuspend')}</button>` : '')}
    ${c.status === 'suspended' ? `<div class="banner bad"><b>${t('csSuspended')}</b></div>` : ''}
    <section class="kpis">${kpi(t('csOrders'), s.orders)}${kpi(t('csSpent'), csUsd(s.spent))}${kpi(t('csCancelled'), s.cancelled, s.cancelled ? 'warn' : '')}${kpi(t('csReturns'), s.returns)}${kpi(t('csUnpaid'), s.unpaid_checkouts)}</section>
    <div class="cols">
      <div style="display:flex;flex-direction:column;gap:16px;min-width:0">
        <section class="panel"><div class="ph"><h2>${t('csOrders')}</h2></div>${c.orders.length ? `<div class="tbl-wrap"><table style="min-width:0;width:100%"><tbody>${c.orders.map(o => `
          <tr class="click" data-href="#/order/${esc(o.order_no)}"><td><b>${esc(o.order_no)}</b><br><span class="muted small">${csDt(o.created_at)}</span></td><td>${esc(o.store_name)}</td><td>${csUsd(o.total)}</td><td>${csOrdSt(o.status)}</td></tr>`).join('')}</tbody></table></div>` : `<p class="muted" style="padding:16px">${t('csNoResult')}</p>`}</section>
        ${c.returns.length ? `<section class="panel"><div class="ph"><h2>${t('csReturns')}</h2></div><div class="tbl-wrap"><table style="min-width:0;width:100%"><tbody>${c.returns.map(r => `
          <tr class="click" data-href="#/order/${esc(r.order_no)}"><td><b>${esc(r.return_no)}</b><br><span class="muted small">${esc(r.order_no)} · ${csDt(r.created_at)}</span></td><td>${esc(r.reason)}</td><td>${csUsd(r.amount)}</td><td><span class="chip">${esc(r.status)}</span></td></tr>`).join('')}</tbody></table></div></section>` : ''}
      </div>
      <div style="display:flex;flex-direction:column;gap:16px;min-width:0">
        <section class="panel"><div class="ph"><h2>${t('csNotes')}</h2></div><div style="padding:14px 18px;display:flex;flex-direction:column;gap:10px">
          ${hasPerm('support.handle') ? `<textarea id="csNote" maxlength="2000" rows="3" placeholder="${t('csNotePh')}"></textarea><div><button type="button" id="csNoteBtn">${t('csAddNote')}</button></div>` : ''}
          ${c.notes.length ? c.notes.map(n => `<div class="rt-says"><b>${t('csKind_' + n.kind)}</b> <span class="muted small">· ${esc(n.author || '')} · ${csDt(n.created_at)}</span><p>${esc(n.note)}</p></div>`).join('') : `<p class="muted small">${t('csNoNotes')}</p>`}</div></section>
        <section class="panel"><div class="ph"><h2>${t('csLogins')}</h2><span class="muted small">${t('csSessions', {n: c.sessions})}</span></div><div style="padding:14px 18px;display:flex;flex-direction:column;gap:6px">
          ${c.identities.map(i => `<div><span class="chip">${esc(i.provider)}</span> ${esc(i.display || '')} <span class="muted small">· ${csDt(i.last_used_at || i.created_at)}</span></div>`).join('') || '—'}</div></section>
        ${c.stores.length ? `<section class="panel"><div class="ph"><h2>${t('csStores')}</h2></div><div style="padding:14px 18px">${c.stores.map(st => `<div><a href="#/store/${st.id}">${esc(st.name)}</a> <span class="muted small">· ${esc(st.role)} · ${esc(st.status)}</span></div>`).join('')}</div></section>` : ''}
        ${c.addresses.length ? `<section class="panel"><div class="ph"><h2>${t('csAddresses')}</h2></div><div style="padding:14px 18px;display:flex;flex-direction:column;gap:8px">${c.addresses.map(a => `<div><b>${esc(a.recipient)}</b> · ${esc(a.phone)}<br><span class="muted small">${esc(a.address)}</span></div>`).join('')}</div></section>` : ''}
      </div></div>`);
  document.querySelectorAll('tr[data-href]').forEach(tr => tr.onclick = () => { location.hash = tr.dataset.href; });
  if ($('#csNoteBtn')) $('#csNoteBtn').onclick = async () => {
    try { await api(`/admin/api/customers/${id}/notes`, {method: 'POST', body: {note: $('#csNote').value}}); toast(t('saved')); render(); }
    catch (e) { showError(e, {note: '#csNote'}); }
  };
  if ($('#csStatusBtn')) $('#csStatusBtn').onclick = () => {
    const suspend = c.status !== 'suspended';
    const d = document.createElement('dialog');
    d.className = 'dlg';
    d.innerHTML = `<form method="dialog"><h2>${t(suspend ? 'csSuspendTitle' : 'csReactivate')}</h2>${suspend ? `<p class="muted">${t('csSuspendHint')}</p>` : ''}
      <label>${t('csReason')}<textarea id="csWhy" maxlength="2000"></textarea></label>
      <div class="row"><button type="button" class="${suspend ? 'danger' : 'primary'}" id="csGo">${t(suspend ? 'csSuspend' : 'csReactivate')}</button><button value="cancel">${t('csBack')}</button></div></form>`;
    document.body.appendChild(d); d.addEventListener('close', () => d.remove()); d.showModal();
    d.querySelector('#csGo').onclick = async () => {
      try { await api(`/admin/api/customers/${id}/status`, {method: 'POST', body: {status: suspend ? 'suspended' : 'active', note: d.querySelector('#csWhy').value}}); d.close(); toast(t('saved')); render(); }
      catch (e) { if (e.field === 'note') { d.querySelector('#csWhy').classList.add('invalid'); d.querySelector('#csWhy').focus(); } else { d.close(); toast(errMsg(e), 'error'); } }
    };
  };
}

async function renderOrderDetail(no) {
  const o = await api('/admin/api/orders/' + encodeURIComponent(no));
  const p = o.payment, d = o.delivery, a = o.address;
  const dl = rows => `<dl class="omoney">${rows.filter(Boolean).map(([k, v, cls]) => `<div class="${cls || ''}"><dt>${k}</dt><dd>${v}</dd></div>`).join('')}</dl>`;
  const sec = (title, inner) => `<section class="panel"><div class="ph"><h2>${title}</h2></div><div style="padding:14px 18px;display:flex;flex-direction:column;gap:10px">${inner}</div></section>`;
  $('#app').innerHTML = shell(`<a href="#/customers" onclick="CS.tab='orders'">← ${t('csBack')}</a>
    ${pageHead(`${t('rtOrder')} ${esc(o.order_no)}`, csDt(o.created_at), csOrdSt(o.status))}
    ${o.cancel_reason ? `<div class="banner warn">${esc(o.cancel_reason)}${o.reject_reason ? ' · ' + esc(o.reject_reason) : ''}${o.seller_note ? ' · ' + esc(o.seller_note) : ''}</div>` : ''}
    ${o.refund.status !== 'none' ? `<div class="banner ${o.refund.status === 'refunded' ? 'ok' : 'warn'}"><b>${t('odRefund')}: ${esc(o.refund.status)}</b> ${csUsd(o.refund.amount || o.total)}${o.refund.ref ? ' · ' + esc(o.refund.ref) : ''}</div>` : ''}
    <div class="cols"><div style="display:flex;flex-direction:column;gap:16px;min-width:0">
      ${sec(t('odItems'), `<div class="oitems">${o.items.map(it => `<div class="oitem"><img src="${esc(it.image_url || '')}" alt="" onerror="this.style.visibility='hidden'"><div><b>${esc(it.name)}</b>${it.options ? `<div class="muted">${esc(it.options)}</div>` : ''}<div class="muted small">SKU ${esc(it.sku)}</div></div><div class="oqty">× ${it.qty}<div class="muted">${csUsd(it.line_total)}</div></div></div>`).join('')}</div>
        ${dl([[t('odSubtotal'), csUsd(o.subtotal)], +o.shipping_fee ? [t('odShipping'), csUsd(o.shipping_fee)] : null, [t('odCommission', {r: parseFloat(o.commission_rate)}), '−' + csUsd(o.commission_amount)], [t('odSeller'), csUsd(o.seller_amount), 'total']])}`)}
      ${sec(t('odTimeline'), `<ol class="timeline">${o.events.map(e => `<li><b>${esc(e.status)}</b> <span class="muted small">${esc(e.actor_type)}${e.actor ? ' · ' + esc(e.actor) : ''} · ${csDt(e.at)}</span>${e.note ? `<div class="muted small">${esc(e.note)}</div>` : ''}</li>`).join('')}</ol>`)}
      ${sec(t('odPayLog'), o.payment_events.map(pe => `<details><summary><span class="chip ${pe.verified ? 'ok' : pe.source === 'callback' ? 'bad' : ''}">${esc(pe.source)}</span> ${esc(pe.status || '')} <span class="muted small">· ${pe.verified ? t('odVerified') : t('odUnverified')} · ${csDt(pe.at)}</span></summary><pre class="json">${esc(JSON.stringify(pe.payload, null, 2))}</pre></details>`).join('') || '—')}
      ${o.ledger.length ? sec(t('odLedger'), o.ledger.map(l => `<div class="spread"><span>${esc(l.kind)} <span class="muted small">· ${csDt(l.at)}</span></span><b>${csUsd(l.amount)}</b></div>`).join('')) : ''}
    </div><div style="display:flex;flex-direction:column;gap:16px;min-width:0">
      ${sec(t('odBuyer'), `<div><b>${esc(o.buyer.name || '—')}</b> ${o.buyer.is_guest ? `<span class="chip">${t('csGuest')}</span>` : ''}${o.buyer.status === 'suspended' ? ` <span class="chip bad">${t('csSuspended')}</span>` : ''}<br><span class="muted">${esc(o.buyer.phone || '')}</span></div>${hasPerm('customer.view') ? `<a href="#/customers/${o.buyer.id}">${t('odCustomer')} →</a>` : ''}`)}
      ${sec(t('odDeliverTo'), `<div><b>${esc(a.recipient)}</b> · <a href="tel:${esc(a.phone)}">${esc(a.phone)}</a><br>${esc(a.line)}${a.note ? `<br><span class="muted">${esc(a.note)}</span>` : ''}</div>${a.lat ? `<a href="https://www.google.com/maps?q=${a.lat},${a.lng}" target="_blank" rel="noopener">${t('odOpenMap')} ↗</a>` : ''}`)}
      ${sec(t('odStore'), `<div><a href="#/store/${o.store.id}">${esc(o.store.name)}</a><br><span class="muted">${esc(o.store.phone || '')}</span></div>`)}
      ${sec(t('odPayment'), `${dl([[t('odTran'), `<code>${esc(p.tran_id)}</code> <button type="button" data-copy="${esc(p.tran_id)}">${t('odCopy')}</button>`], [t('odPaidAmt'), p.paid_amount ? csUsd(p.paid_amount) : '—'], [t('pyCheckout'), esc(p.checkout_status) + ' · ' + csUsd(p.checkout_total)], p.ref ? [t('odRef'), esc(p.ref)] : null, p.payer ? [t('odPayer'), esc(p.payer)] : null, p.type ? ['Type', esc(p.type)] : null, [t('odPaidAt'), csDt(p.paid_at)]])}
        ${p.orders_in_checkout > 1 ? `<p class="muted small">${t('odSplit', {n: p.orders_in_checkout})}</p>` : ''}`)}
      ${sec(t('odDelivery'), d ? `${dl([['Status', `<span class="chip">${esc(d.status)}</span>`], [t('odCourier'), esc(d.courier || '—') + (d.courier_phone ? ` · <a href="tel:${esc(d.courier_phone)}">${esc(d.courier_phone)}</a>` : '')], ['Pickup', esc(d.pickup || '—')], d.picked_up_at ? ['↑', csDt(d.picked_up_at)] : null, d.delivered_at ? ['✓', csDt(d.delivered_at)] : null, d.fail_reason ? ['✗', esc(d.fail_reason) + (d.fail_note ? ' · ' + esc(d.fail_note) : '')] : null])}
        ${d.has_proof ? `<button type="button" id="odProof">${t('odProof')}</button>` : ''}` : `<p class="muted">${t('odNoDelivery')}</p>`)}
      ${o.returns.length ? sec(t('odReturns'), o.returns.map(r => `<div class="spread"><a href="#/returns">${esc(r.return_no)}</a><span><span class="chip">${esc(r.status)}</span> ${csUsd(r.amount)}</span></div>`).join('')) : ''}
    </div></div>`);
  document.querySelectorAll('[data-copy]').forEach(b => b.onclick = () => { navigator.clipboard?.writeText(b.dataset.copy); toast(t('odCopied')); });
  if ($('#odProof')) $('#odProof').onclick = async () => {
    const r = await fetch(`/admin/api/deliveries/${d.id}/proof`, {headers: {Authorization: 'Bearer ' + P.token}});
    if (!r.ok) { toast(r.statusText, 'error'); return; }
    window.open(URL.createObjectURL(await r.blob()), '_blank', 'noopener');
  };
}

// ---------- Log pembayaran ----------
const PY = {filter: 'all', q: ''};
async function renderPayments() {
  const {events, counts} = await api(`/admin/api/payment-events?filter=${PY.filter}&q=${encodeURIComponent(PY.q)}`);
  const rows = events.map(e => `<tr><td>${csDt(e.at)}</td><td><span class="chip ${e.verified ? 'ok' : e.source === 'callback' ? 'bad' : ''}">${esc(e.source)}</span><br><span class="muted small">${e.verified ? t('odVerified') : t('odUnverified')}</span></td>
    <td><code>${esc(e.tran_id || '')}</code></td><td>${esc(e.status || '—')}</td><td>${e.checkout_status ? esc(e.checkout_status) + ' · ' + csUsd(e.checkout_total) : '—'}</td>
    <td>${(e.orders || []).map(n => `<a href="#/order/${esc(n)}">${esc(n)}</a>`).join('<br>') || '—'}</td>
    <td><details><summary>${t('odPayload')}</summary><pre class="json">${esc(JSON.stringify(e.payload, null, 2))}</pre></details></td></tr>`).join('');
  $('#app').innerHTML = shell(`${pageHead(t('i_payments'), t('pyHint'))}
    <section class="kpis"><div class="kpi"><span>${t('pyEvents24')}</span><b>${counts.today}</b></div><div class="kpi"><span>${t('pyPaid24')}</span><b>${counts.paid_24h}</b></div>
      <div class="kpi"><span>${t('pyProblems7')}</span><b class="${counts.problems ? 'warn' : ''}">${counts.problems}</b></div></section>
    <div class="spread" style="flex-wrap:wrap;gap:10px"><div class="tabs">${['all', 'callback', 'problems'].map(f => `<button type="button" aria-selected="${PY.filter === f}" data-pyf="${f}">${t({all: 'pyAll', callback: 'pyCallbacks', problems: 'pyProblems'}[f])}</button>`).join('')}</div>
      <form id="pyForm" class="row" style="gap:8px"><input id="pyQ" type="search" value="${esc(PY.q)}" placeholder="${t('pySearchPh')}"><button type="submit">${t('csSearch')}</button></form></div>
    <section class="panel">${rows ? `<div class="tbl-wrap"><table><thead><tr><th>Time</th><th>Source</th><th>${t('odTran')}</th><th>Status</th><th>${t('pyCheckout')}</th><th>${t('pyOrders')}</th><th></th></tr></thead><tbody>${rows}</tbody></table></div>` : `<p class="muted" style="padding:16px">${t('csNoResult')}</p>`}</section>`);
  document.querySelectorAll('[data-pyf]').forEach(b => b.onclick = () => { PY.filter = b.dataset.pyf; render(); });
  $('#pyForm').onsubmit = e => { e.preventDefault(); PY.q = $('#pyQ').value.trim(); render(); };
}
