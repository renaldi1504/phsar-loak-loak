'use strict';
// Seller Center: Bantuan — tiket ke CS Phsar atas nama toko (semua anggota toko bisa melihat & membalas).
addStrings({
  en: {tabHelp: 'Help', hpTitle: 'Phsar support', hpSub: 'Questions about orders, payouts, products or your store? Write to us — we usually reply within one working day.',
    hpNew: 'New request', hpNone: 'No support requests yet.', hpTopic: 'Topic', hpSubject: 'Subject', hpOrder: 'Order number (optional)', hpMessage: 'Message',
    hpSend: 'Send', hpReply: 'Reply', hpSolve: 'Mark as solved', hpBack: 'All requests', hpSent: 'Sent. We will reply here and notify you.', hpYou: 'Your store',
    hpReopen: 'Write again to reopen this request.',
    hpCat_order: 'Order', hpCat_payment: 'Payment / payout', hpCat_return: 'Return', hpCat_account: 'Account', hpCat_seller: 'Store & products', hpCat_other: 'Other',
    hpSt_open: 'Waiting for Phsar', hpSt_pending: 'Phsar replied', hpSt_solved: 'Solved', err_rate_limited: 'Too many requests today. Please try again tomorrow.'},
  km: {tabHelp: 'ជំនួយ', hpTitle: 'ជំនួយ Phsar', hpSub: 'មានសំណួរអំពីការបញ្ជាទិញ ការទូទាត់ ផលិតផល ឬហាង? សរសេរមកយើង — ជាធម្មតាយើងឆ្លើយក្នុងមួយថ្ងៃធ្វើការ។',
    hpNew: 'សំណើថ្មី', hpNone: 'មិនទាន់មានសំណើជំនួយទេ។', hpTopic: 'ប្រធានបទ', hpSubject: 'ចំណងជើង', hpOrder: 'លេខការបញ្ជាទិញ (ស្រេចចិត្ត)', hpMessage: 'សារ',
    hpSend: 'ផ្ញើ', hpReply: 'ឆ្លើយតប', hpSolve: 'សម្គាល់ថាដោះស្រាយរួច', hpBack: 'សំណើទាំងអស់', hpSent: 'បានផ្ញើ។ យើងនឹងឆ្លើយនៅទីនេះ ហើយជូនដំណឹងអ្នក។', hpYou: 'ហាងរបស់អ្នក',
    hpReopen: 'សរសេរម្តងទៀតដើម្បីបើកសំណើនេះឡើងវិញ។',
    hpCat_order: 'ការបញ្ជាទិញ', hpCat_payment: 'ការបង់ប្រាក់ / ការទូទាត់', hpCat_return: 'ការប្ដូរ', hpCat_account: 'គណនី', hpCat_seller: 'ហាង និងផលិតផល', hpCat_other: 'ផ្សេងៗ',
    hpSt_open: 'រង់ចាំ Phsar', hpSt_pending: 'Phsar បានឆ្លើយ', hpSt_solved: 'បានដោះស្រាយ', err_rate_limited: 'សំណើច្រើនពេកថ្ងៃនេះ។ សូមព្យាយាមម្តងទៀតថ្ងៃស្អែក។'},
  zh: {tabHelp: '帮助', hpTitle: 'Phsar 客服', hpSub: '关于订单、结算、商品或店铺有疑问？给我们留言——我们通常在一个工作日内回复。',
    hpNew: '新建请求', hpNone: '暂无支持请求。', hpTopic: '主题', hpSubject: '标题', hpOrder: '订单号（可选）', hpMessage: '内容',
    hpSend: '发送', hpReply: '回复', hpSolve: '标记为已解决', hpBack: '全部请求', hpSent: '已发送，我们会在这里回复并通知您。', hpYou: '您的店铺',
    hpReopen: '再次留言即可重新打开此请求。',
    hpCat_order: '订单', hpCat_payment: '付款 / 结算', hpCat_return: '退货', hpCat_account: '账户', hpCat_seller: '店铺与商品', hpCat_other: '其他',
    hpSt_open: '等待 Phsar 回复', hpSt_pending: 'Phsar 已回复', hpSt_solved: '已解决', err_rate_limited: '今天的请求过多，请明天再试。'},
});

const HP_CATS = ['seller', 'order', 'payment', 'return', 'account', 'other'];
const HP_CHIP = {open: 'warn', pending: 'pink', solved: 'ok'};

async function tabHelp(id) {
  const open = S.helpTicket; S.helpTicket = null;
  if (open) return helpThread(id, open);
  const {tickets} = await api(`/v1/seller/stores/${id}/tickets`);
  // sellerpages: kartu ajakan "Tulis ke Phsar" di atas, lalu permintaan sebelumnya.
  const list = tickets.map(x => `<section class="card order sp-tk" role="link" tabindex="0" data-tk="${x.id}"><div class="spread"><div style="min-width:0"><h2>${esc(x.subject)}</h2>
      <span class="muted">${esc(x.ticket_no)} · ${t('hpCat_' + x.category)}${x.order_no ? ' · ' + esc(x.order_no) : ''} · ${fmtTime(x.updated_at)}</span></div>
      <span class="chip ${HP_CHIP[x.status]}">${t('hpSt_' + x.status)}</span></div><p class="muted">${esc(x.preview)}</p></section>`).join('');
  $('#tab').innerHTML = `<section class="sp-hero"><h2>💬 ${t('hpTitle')}</h2><p>${t('hpSub')}</p>
      <button type="button" class="primary sp-block" id="hpNewBtn">${t('spHpWrite')}</button></section>
    ${list ? `<h2 class="sp-sub">${t('spHpMine')}</h2>${list}` : spEmpty('📭', t('hpNone'), t('spHpNoneW'))}`;
  $('#hpNewBtn').onclick = () => helpNew(id);
  document.querySelectorAll('[data-tk]').forEach(c => { c.onclick = () => helpThread(id, c.dataset.tk); c.onkeydown = e => { if (e.key === 'Enter') helpThread(id, c.dataset.tk); }; });
}

function helpNew(id) {
  $('#tab').innerHTML = `<button type="button" class="sp-backbtn" id="hpBack">← ${t('hpBack')}</button><section class="card" id="hpForm"><h2>${t('hpNew')}</h2>
    <label>${t('hpTopic')}<select id="hpCat">${HP_CATS.map(c => `<option value="${c}">${t('hpCat_' + c)}</option>`).join('')}</select></label>
    <label>${t('hpSubject')}<input id="hpSubject" name="subject" maxlength="120"></label>
    <label>${t('hpOrder')}<input id="hpOrder" name="order_no" maxlength="20" placeholder="PH261114-XXXXX"></label>
    <label>${t('hpMessage')}<textarea id="hpMsg" name="message" maxlength="4000" rows="6"></textarea></label>
    <div data-errors-anchor><button type="button" class="primary sp-block" id="hpSend">${t('hpSend')}</button></div></section>`;
  $('#hpBack').onclick = () => tabHelp(id);
  $('#hpSend').onclick = async () => {
    try {
      const x = await api(`/v1/seller/stores/${id}/tickets`, {method: 'POST', body: {category: $('#hpCat').value, subject: $('#hpSubject').value, order_no: $('#hpOrder').value, message: $('#hpMsg').value}});
      toast(t('hpSent')); helpThread(id, x.id);
    } catch (e) { showError(e, {subject: '#hpSubject', order_no: '#hpOrder', message: '#hpMsg', category: '#hpCat'}); }
  };
}

async function helpThread(id, tid) {
  const x = await api(`/v1/seller/stores/${id}/tickets/${tid}`);
  const msgs = x.messages.map(m => `<div class="hp-msg ${m.author_type === 'user' ? 'mine' : ''}"><div class="hp-bubble">${esc(m.body)}</div>
    <span class="muted small">${m.author_type === 'user' ? esc(m.author_name || t('hpYou')) : esc(m.author_name)} · ${fmtTime(m.created_at)}</span></div>`).join('');
  $('#tab').innerHTML = `<button type="button" class="sp-backbtn" id="hpBack">← ${t('hpBack')}</button>
    <section class="card"><div class="spread"><div><h2>${esc(x.subject)}</h2><span class="muted">${esc(x.ticket_no)} · ${t('hpCat_' + x.category)}${x.order_no ? ' · ' + esc(x.order_no) : ''}</span></div>
      <span class="chip ${HP_CHIP[x.status]}">${t('hpSt_' + x.status)}</span></div>
      <div class="hp-thread">${msgs}</div>
      ${x.status === 'solved' ? `<p class="muted small">${t('hpReopen')}</p>` : ''}
      <label>${t('hpReply')}<textarea id="hpBody" maxlength="4000" rows="3" placeholder="${t('spHpReplyPh')}"></textarea></label>
      <div class="sp-acts"><button type="button" class="primary sp-block" id="hpReply">${t('hpSend')}</button>${x.status !== 'solved' ? `<button type="button" class="textbtn" id="hpSolve">✓ ${t('hpSolve')}</button>` : ''}</div></section>`;
  $('#hpBack').onclick = () => tabHelp(id);
  $('#hpReply').onclick = async () => {
    const body = $('#hpBody').value.trim();
    if (!body) { $('#hpBody').focus(); return; }
    try { await api(`/v1/seller/stores/${id}/tickets/${tid}/messages`, {method: 'POST', body: {body}}); helpThread(id, tid); } catch (e) { toast(errMsg(e), 'error'); }
  };
  if ($('#hpSolve')) $('#hpSolve').onclick = async () => {
    try { await api(`/v1/seller/stores/${id}/tickets/${tid}/solve`, {method: 'POST'}); helpThread(id, tid); } catch (e) { toast(errMsg(e), 'error'); }
  };
}
