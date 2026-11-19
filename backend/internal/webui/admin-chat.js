'use strict';
// chat: admin membaca chat pembeli ↔ toko untuk sengketa (izin support.handle, hanya baca).
// Membuka satu percakapan dicatat di audit log (chat.view).
addStrings({
  en: {i_chats: 'Buyer–seller chats', achHint: 'Read-only, for disputes. Search by order number, buyer phone or name, or shop name. Every opened chat is recorded in the audit log.',
    achSearch: 'Order, phone, buyer or shop…', achBuyer: 'Buyer', achShop: 'Shop', achBack: 'All chats', achAudit: 'Opening this chat was recorded in the audit log.',
    achUnread: 'Unread: buyer {b} · shop {s}', achProduct: 'Product', achOrder: 'Order {no}', 'a_chat.view': 'Viewed a buyer–seller chat'},
  km: {i_chats: 'ការជជែកអ្នកទិញ–អ្នកលក់', achHint: 'មើលតែប៉ុណ្ណោះ សម្រាប់វិវាទ។ ស្វែងរកតាមលេខការបញ្ជាទិញ ទូរស័ព្ទ ឬឈ្មោះអ្នកទិញ ឬឈ្មោះហាង។ ការជជែកដែលបានបើកនីមួយៗត្រូវបានកត់ក្នុងកំណត់ហេតុសវនកម្ម។',
    achSearch: 'ការបញ្ជាទិញ ទូរស័ព្ទ អ្នកទិញ ឬហាង…', achBuyer: 'អ្នកទិញ', achShop: 'ហាង', achBack: 'ការជជែកទាំងអស់', achAudit: 'ការបើកការជជែកនេះត្រូវបានកត់ក្នុងកំណត់ហេតុសវនកម្ម។',
    achUnread: 'មិនទាន់អាន៖ អ្នកទិញ {b} · ហាង {s}', achProduct: 'ផលិតផល', achOrder: 'ការបញ្ជាទិញ {no}', 'a_chat.view': 'បានមើលការជជែកអ្នកទិញ–អ្នកលក់'},
  zh: {i_chats: '买家与卖家聊天', achHint: '只读，用于纠纷处理。可按订单号、买家电话或姓名、店铺名搜索。每次打开聊天都会记录在审计日志中。',
    achSearch: '订单、电话、买家或店铺…', achBuyer: '买家', achShop: '店铺', achBack: '全部聊天', achAudit: '打开此聊天已记录在审计日志中。',
    achUnread: '未读：买家 {b} · 店铺 {s}', achProduct: '商品', achOrder: '订单 {no}', 'a_chat.view': '查看买家与卖家聊天'},
});

const ACH = {q: '', photos: {}};

async function renderChats() {
  const m = location.hash.match(/^#\/chats\/([0-9a-f-]{36})/);
  if (m) return renderChatAdmin(m[1]);
  const {conversations} = await api(`/admin/api/chats?q=${encodeURIComponent(ACH.q)}`);
  const rows = conversations.map(c => `<a class="card tk-row" href="#/chats/${c.id}" style="text-decoration:none;color:inherit">
      <div class="spread"><div><b>${esc(c.store_name)}</b> ↔ ${esc(c.buyer_name || c.buyer_phone || '—')}<br>
        <span class="muted small">${t('achBuyer')}: ${esc(c.buyer_phone || '—')} · ${new Date(c.last_message_at).toLocaleString()}</span></div>
        <span class="muted small">${t('achUnread', {b: c.buyer_unread, s: c.seller_unread})}</span></div>
      <p class="muted small" style="margin:0">${esc(c.preview)}</p></a>`).join('');
  $('#app').innerHTML = shell(`${pageHead(t('i_chats'), t('achHint'))}
    <div class="row"><input id="achQ" type="search" placeholder="${t('achSearch')}" value="${esc(ACH.q)}" style="flex:1;min-width:200px"></div>
    ${rows || `<section class="card"><p class="muted">${t('empty')}</p></section>`}`);
  $('#achQ').onkeydown = e => { if (e.key === 'Enter') { ACH.q = e.target.value.trim(); render(); } };
}

function achCard(c) {
  if (!c) return '';
  return c.type === 'order' ? `<div class="muted small">📦 <a href="#/order/${encodeURIComponent(c.order_no)}">${esc(t('achOrder', {no: c.order_no}))}</a> · ${esc(c.status)} · $${esc(c.total)}</div>`
    : `<div class="muted small">🛍️ ${t('achProduct')}: <a href="#/product/${esc(c.product_id)}">${esc(c.name)}</a>${c.price ? ' · $' + esc(c.price) : ''}</div>`;
}

async function renderChatAdmin(id) {
  const c = await api('/admin/api/chats/' + id);
  const msgs = (c.messages || []).map(m => `<div class="tk-amsg ${m.sender === 'seller' ? 'admin' : ''}">${achCard(m.card)}
      ${m.has_photo ? `<img alt="" style="max-width:240px;border-radius:12px" data-achphoto="/admin/api/chats/${c.id}/messages/${m.id}/photo">` : ''}
      ${m.body ? `<div class="tk-abubble">${esc(m.body)}</div>` : ''}
      <span class="muted small">${esc(m.sender === 'seller' ? `${c.store_name}${m.sender_name ? ' (' + m.sender_name + ')' : ''}` : (c.buyer_name || t('achBuyer')))} · ${new Date(m.created_at).toLocaleString()}</span></div>`).join('');
  $('#app').innerHTML = shell(`<a href="#/chats">← ${t('achBack')}</a>
    ${pageHead(`${esc(c.store_name)} ↔ ${esc(c.buyer_name || c.buyer_phone || '—')}`, t('achAudit'))}
    <section class="card"><div class="grid">
      <p><span class="muted">${t('achShop')}</span><br><a href="#/store/${esc(c.store_id)}">${esc(c.store_name)}</a></p>
      <p><span class="muted">${t('achBuyer')}</span><br>${esc(c.buyer_name || '—')} · ${esc(c.buyer_phone || '')}</p>
    </div></section>
    <section class="card"><div class="tk-athread">${msgs}</div></section>`);
  for (const img of document.querySelectorAll('img[data-achphoto]')) {
    const src = img.dataset.achphoto;
    try {
      if (!ACH.photos[src]) {
        const r = await fetch(src, {headers: {Authorization: 'Bearer ' + P.token}});
        if (!r.ok) continue;
        ACH.photos[src] = URL.createObjectURL(await r.blob());
      }
      img.src = ACH.photos[src];
    } catch {}
  }
}
