'use strict';
// Admin: kotak masuk tiket bantuan dari pembeli & seller (izin support.handle).
addStrings({
  en: {q_tickets: 'Support tickets waiting for a reply', tkHint: 'Requests from buyers (website/app) and sellers (Seller Center). Oldest waiting first. Replies are sent to the user and notify them; internal notes are only for the team.',
    tf_open: 'Open', tf_pending: 'Waiting for user', tf_mine: 'Mine', tf_unassigned: 'Unassigned', tf_solved: 'Solved', tf_all: 'All', tkSearch: 'Ticket, order, phone, subject…',
    tkFrom: 'From', tkBuyer: 'Buyer', tkSeller: 'Seller', tkOrder: 'Order', tkAssigned: 'Assigned to', tkNobody: 'Nobody', tkAssignMe: 'Assign to me', tkUnassign: 'Unassign',
    tkStatus: 'Status', tkReplyLbl: 'Reply to the user', tkInternal: 'Internal note (not sent to the user)', tkSendReply: 'Send reply', tkSendSolve: 'Reply & mark solved',
    tkAddNote: 'Add note', tkNoteTag: 'Internal note', tkWaitingFor: 'waiting {t}', tkUpdated: 'Saved', tkMsgs: '{n} messages', tkBackList: 'All tickets',
    tkCat_order: 'Order', tkCat_payment: 'Payment', tkCat_return: 'Return', tkCat_account: 'Account', tkCat_seller: 'Seller', tkCat_other: 'Other',
    tkSt_open: 'Open', tkSt_pending: 'Waiting for user', tkSt_solved: 'Solved'},
  km: {q_tickets: 'សំបុត្រជំនួយរង់ចាំការឆ្លើយ', tkHint: 'សំណើពីអ្នកទិញ (គេហទំព័រ/កម្មវិធី) និងអ្នកលក់ (Seller Center)។ រង់ចាំយូរបំផុតនៅខាងលើ។ ការឆ្លើយត្រូវផ្ញើទៅអ្នកប្រើ និងជូនដំណឹង កំណត់ចំណាំខាងក្នុងសម្រាប់ក្រុមប៉ុណ្ណោះ។',
    tf_open: 'បើក', tf_pending: 'រង់ចាំអ្នកប្រើ', tf_mine: 'របស់ខ្ញុំ', tf_unassigned: 'មិនទាន់ចាត់តាំង', tf_solved: 'បានដោះស្រាយ', tf_all: 'ទាំងអស់', tkSearch: 'សំបុត្រ ការបញ្ជាទិញ ទូរស័ព្ទ ចំណងជើង…',
    tkFrom: 'ពី', tkBuyer: 'អ្នកទិញ', tkSeller: 'អ្នកលក់', tkOrder: 'ការបញ្ជាទិញ', tkAssigned: 'ចាត់តាំងឱ្យ', tkNobody: 'គ្មាន', tkAssignMe: 'ចាត់តាំងឱ្យខ្ញុំ', tkUnassign: 'ដកការចាត់តាំង',
    tkStatus: 'ស្ថានភាព', tkReplyLbl: 'ឆ្លើយទៅអ្នកប្រើ', tkInternal: 'កំណត់ចំណាំខាងក្នុង (មិនផ្ញើទៅអ្នកប្រើ)', tkSendReply: 'ផ្ញើការឆ្លើយ', tkSendSolve: 'ឆ្លើយ និងសម្គាល់ថាដោះស្រាយរួច',
    tkAddNote: 'បន្ថែមកំណត់ចំណាំ', tkNoteTag: 'កំណត់ចំណាំខាងក្នុង', tkWaitingFor: 'រង់ចាំ {t}', tkUpdated: 'បានរក្សាទុក', tkMsgs: 'សារ {n}', tkBackList: 'សំបុត្រទាំងអស់',
    tkCat_order: 'ការបញ្ជាទិញ', tkCat_payment: 'ការបង់ប្រាក់', tkCat_return: 'ការប្ដូរ', tkCat_account: 'គណនី', tkCat_seller: 'អ្នកលក់', tkCat_other: 'ផ្សេងៗ',
    tkSt_open: 'បើក', tkSt_pending: 'រង់ចាំអ្នកប្រើ', tkSt_solved: 'បានដោះស្រាយ'},
  zh: {q_tickets: '待回复的支持工单', tkHint: '来自买家（网站/App）和卖家（卖家中心）的请求，等待最久的排在最前。回复会发送给用户并通知对方；内部备注仅团队可见。',
    tf_open: '待处理', tf_pending: '等待用户', tf_mine: '我的', tf_unassigned: '未分配', tf_solved: '已解决', tf_all: '全部', tkSearch: '工单号、订单号、电话、标题…',
    tkFrom: '来自', tkBuyer: '买家', tkSeller: '卖家', tkOrder: '订单', tkAssigned: '负责人', tkNobody: '无', tkAssignMe: '分配给我', tkUnassign: '取消分配',
    tkStatus: '状态', tkReplyLbl: '回复用户', tkInternal: '内部备注（不发送给用户）', tkSendReply: '发送回复', tkSendSolve: '回复并标记已解决',
    tkAddNote: '添加备注', tkNoteTag: '内部备注', tkWaitingFor: '已等待 {t}', tkUpdated: '已保存', tkMsgs: '{n} 条消息', tkBackList: '全部工单',
    tkCat_order: '订单', tkCat_payment: '付款', tkCat_return: '退货', tkCat_account: '账户', tkCat_seller: '卖家', tkCat_other: '其他',
    tkSt_open: '待处理', tkSt_pending: '等待用户', tkSt_solved: '已解决'},
});

const ATK = {filter: 'open', q: ''};
const ATK_CHIP = {open: 'warn', pending: 'pink', solved: 'ok'};

async function renderSupport() {
  const m = location.hash.match(/^#\/support\/([0-9a-f-]{36})/);
  if (m) return renderTicket(m[1]);
  const {tickets, counts} = await api(`/admin/api/tickets?filter=${ATK.filter}&q=${encodeURIComponent(ATK.q)}`);
  const n = {open: counts.open, pending: counts.pending, mine: counts.mine, unassigned: counts.unassigned};
  const tabs = ['open', 'pending', 'mine', 'unassigned', 'solved', 'all'].map(f => `<button type="button" aria-selected="${ATK.filter === f}" data-tf="${f}">${t('tf_' + f)}${n[f] ? ` <span class="count">${n[f]}</span>` : ''}</button>`).join('');
  const rows = tickets.map(x => `<a class="card tk-row" href="#/support/${x.id}" style="text-decoration:none;color:inherit">
      <div class="spread"><div><b>${esc(x.subject)}</b><br><span class="muted small">${esc(x.ticket_no)} · ${t(x.source === 'seller' ? 'tkSeller' : 'tkBuyer')}: ${esc(x.source === 'seller' ? (x.store_name || '') : (x.requester_name || x.requester_phone || '—'))}
        · ${t('tkCat_' + x.category)}${x.order_no ? ' · ' + esc(x.order_no) : ''}</span></div>
        <div style="text-align:right"><span class="chip ${ATK_CHIP[x.status]}">${t('tkSt_' + x.status)}</span>${x.waiting ? `<br><span class="muted small">${t('tkWaitingFor', {t: ago(x.last_user_message_at)})}</span>` : ''}</div></div>
      <p class="muted small" style="margin:0">${esc(x.preview)}</p>
      <span class="muted small">${t('tkAssigned')}: ${esc(x.assigned_name || t('tkNobody'))} · ${t('tkMsgs', {n: x.message_count})}</span></a>`).join('');
  $('#app').innerHTML = shell(`${pageHead(t('i_support'), t('tkHint'))}
    <div class="row" style="gap:10px;flex-wrap:wrap"><div class="tabs count">${tabs}</div>
      <input id="tkQ" type="search" placeholder="${t('tkSearch')}" value="${esc(ATK.q)}" style="flex:1;min-width:200px"></div>
    ${rows || `<section class="card"><p class="muted">${t('empty')}</p></section>`}`);
  document.querySelectorAll('[data-tf]').forEach(b => b.onclick = () => { ATK.filter = b.dataset.tf; render(); });
  $('#tkQ').onkeydown = e => { if (e.key === 'Enter') { ATK.q = e.target.value.trim(); render(); } };
}

async function renderTicket(id) {
  const x = await api('/admin/api/tickets/' + id);
  const who = x.source === 'seller'
    ? `${t('tkSeller')}: <a href="#/store/${esc(x.store_id)}">${esc(x.store_name)}</a> · ${esc(x.requester_name || '')} ${esc(x.requester_phone || '')}`
    : `${t('tkBuyer')}: ${esc(x.requester_name || '—')} · ${esc(x.requester_phone || '')}`;
  const msgs = x.messages.map(m => `<div class="tk-amsg ${m.author_type}${m.internal ? ' internal' : ''}">
      <div class="tk-abubble">${m.internal ? `<b class="small">${t('tkNoteTag')}</b><br>` : ''}${esc(m.body)}</div>
      <span class="muted small">${esc(m.author_name)} · ${new Date(m.created_at).toLocaleString()}</span></div>`).join('');
  const mine = x.assigned_to && A.me && x.assigned_to === A.me.user_id;
  $('#app').innerHTML = shell(`<a href="#/support">← ${t('tkBackList')}</a>
    ${pageHead(x.subject, `${esc(x.ticket_no)} · ${t('tkCat_' + x.category)} · ${new Date(x.created_at).toLocaleString()}`)}
    <section class="card"><div class="grid">
      <p><span class="muted">${t('tkFrom')}</span><br>${who}</p>
      <p><span class="muted">${t('tkOrder')}</span><br>${x.order_no ? `<a href="#/order/${encodeURIComponent(x.order_no)}">${esc(x.order_no)}</a>` : '—'}</p>
      <p><span class="muted">${t('tkAssigned')}</span><br>${esc(x.assigned_name || t('tkNobody'))} <button type="button" id="tkAssign">${t(mine ? 'tkUnassign' : 'tkAssignMe')}</button></p>
      <p><span class="muted">${t('tkStatus')}</span><br><select id="tkStatus">${['open', 'pending', 'solved'].map(s => `<option value="${s}" ${x.status === s ? 'selected' : ''}>${t('tkSt_' + s)}</option>`).join('')}</select></p>
    </div></section>
    <section class="card"><div class="tk-athread">${msgs}</div></section>
    <section class="card" id="tkReplyBox"><label>${t('tkReplyLbl')}<textarea id="tkBody" name="body" rows="4" maxlength="4000"></textarea></label>
      <label class="inline"><input type="checkbox" id="tkInternal"> ${t('tkInternal')}</label>
      <div class="row"><button type="button" class="primary" data-send="pending">${t('tkSendReply')}</button><button type="button" data-send="solved">${t('tkSendSolve')}</button></div></section>`);
  $('#tkInternal').onchange = e => {
    document.querySelector('[data-send="pending"]').textContent = t(e.target.checked ? 'tkAddNote' : 'tkSendReply');
    document.querySelector('[data-send="solved"]').hidden = e.target.checked;
  };
  document.querySelectorAll('[data-send]').forEach(b => b.onclick = async () => {
    const internal = $('#tkInternal').checked, body = $('#tkBody').value.trim();
    if (!body) { $('#tkBody').focus(); return; }
    try {
      await api(`/admin/api/tickets/${id}/reply`, {method: 'POST', body: {body, internal, status: internal ? '' : b.dataset.send}});
      toast(t('tkUpdated')); SH.at = 0; render();
    } catch (e) { showError(e, {body: '#tkBody'}); }
  });
  $('#tkAssign').onclick = async () => { try { await api(`/admin/api/tickets/${id}/assign`, {method: 'POST', body: {to_me: !mine}}); render(); } catch (e) { toast(errMsg(e), 'error'); } };
  $('#tkStatus').onchange = async e => { try { await api(`/admin/api/tickets/${id}/status`, {method: 'POST', body: {status: e.target.value}}); toast(t('tkUpdated')); SH.at = 0; render(); } catch (err) { toast(errMsg(err), 'error'); } };
}
