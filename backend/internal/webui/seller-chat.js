'use strict';
// chat: Seller Center — chat dengan pembeli (Lainnya → Pembeli → Chat). Pemilik & staf toko boleh membaca dan membalas.
// Tanpa websocket: daftar dan percakapan yang sedang terbuka diperbarui tiap beberapa detik (hanya saat halaman terlihat).
addStrings({
  en: {mi_chat: 'Chat with buyers', mi_chatD: 'Answer questions about products and orders',
    tk_chat: 'Reply to {n} chat(s)', tk_chatW: 'Buyers are waiting for an answer. Quick replies help you sell more.', tk_chatB: 'Reply',
    chNone: 'No chats yet. When a buyer writes to your shop from a product, the shop page or an order, it shows up here.',
    chBuyer: 'Buyer', chWaiting: 'Waiting for reply', chBack: '← All chats', chPlaceholder: 'Write a reply…', chSend: 'Send', chPhoto: 'Add photo',
    chSuggest: '✨ Suggest reply', chSuggesting: 'Thinking…', chYou: 'You', chProduct: 'Product', chOrder: 'Order {no}', chPhotoAlt: 'Photo from chat',
    chHint: 'Be polite and never ask buyers to pay outside Phsar.', chTooLong: 'At most 2000 characters.', chUploading: 'Uploading photo…',
    err_chat_rate_limited: 'You are sending messages too fast. Please wait a moment.'},
  km: {mi_chat: 'ជជែកជាមួយអ្នកទិញ', mi_chatD: 'ឆ្លើយសំណួរអំពីផលិតផល និងការបញ្ជាទិញ',
    tk_chat: 'ឆ្លើយការជជែក {n}', tk_chatW: 'អ្នកទិញកំពុងរង់ចាំចម្លើយ។ ការឆ្លើយរហ័សជួយឱ្យអ្នកលក់បានច្រើន។', tk_chatB: 'ឆ្លើយ',
    chNone: 'មិនទាន់មានការជជែកទេ។ ពេលអ្នកទិញសរសេរមកហាងពីផលិតផល ទំព័រហាង ឬការបញ្ជាទិញ វានឹងបង្ហាញនៅទីនេះ។',
    chBuyer: 'អ្នកទិញ', chWaiting: 'រង់ចាំការឆ្លើយ', chBack: '← ការជជែកទាំងអស់', chPlaceholder: 'សរសេរចម្លើយ…', chSend: 'ផ្ញើ', chPhoto: 'បន្ថែមរូបថត',
    chSuggest: '✨ ណែនាំចម្លើយ', chSuggesting: 'កំពុងគិត…', chYou: 'អ្នក', chProduct: 'ផលិតផល', chOrder: 'ការបញ្ជាទិញ {no}', chPhotoAlt: 'រូបថតពីការជជែក',
    chHint: 'សូមគួរសម ហើយកុំសុំឱ្យអ្នកទិញបង់ប្រាក់ក្រៅ Phsar។', chTooLong: 'អតិបរមា 2000 តួអក្សរ។', chUploading: 'កំពុងផ្ទុករូបថត…',
    err_chat_rate_limited: 'អ្នកផ្ញើសារលឿនពេក។ សូមរង់ចាំបន្តិច។'},
  zh: {mi_chat: '买家聊天', mi_chatD: '回答关于商品和订单的问题',
    tk_chat: '回复 {n} 个聊天', tk_chatW: '买家正在等待回复，及时回复有助于成交。', tk_chatB: '回复',
    chNone: '暂无聊天。买家从商品、店铺页面或订单给你发消息后，会显示在这里。',
    chBuyer: '买家', chWaiting: '等待回复', chBack: '← 全部聊天', chPlaceholder: '输入回复…', chSend: '发送', chPhoto: '添加图片',
    chSuggest: '✨ 建议回复', chSuggesting: '思考中…', chYou: '你', chProduct: '商品', chOrder: '订单 {no}', chPhotoAlt: '聊天图片',
    chHint: '请礼貌回复，不要让买家在 Phsar 以外付款。', chTooLong: '最多 2000 个字符。', chUploading: '正在上传图片…',
    err_chat_rate_limited: '发送太快了，请稍等片刻。'},
});

const CHAT = {conv: null, last: 0, timer: null, noAssist: false, photos: {}};

function chatBase() { return `/v1/seller/stores/${S.storeId}/chat`; }
function chatStop() { if (CHAT.timer) { clearInterval(CHAT.timer); CHAT.timer = null; } }
// Polling hanya selama elemen chat masih tampil (pindah tab → berhenti sendiri).
function chatPoll(fn, ms, el) {
  chatStop();
  CHAT.timer = setInterval(() => {
    if (!document.getElementById(el)) return chatStop();
    if (document.visibilityState === 'visible') fn().catch(() => {});
  }, ms);
}

async function tabChat(id) {
  const open = S.chatConv; S.chatConv = null;
  if (open) return chatThread(id, open);
  chatStop();
  await chatList(id);
  chatPoll(() => chatList(id, true), 10000, 'chatList');
}

async function chatList(id, quiet) {
  const r = await api(chatBase());
  const rows = r.conversations.map(c => `<a class="mitem chat-row" href="#/store/${id}/chat/${c.id}" data-conv="${c.id}" style="--mb:${c.waiting ? '#ffe4ec' : '#f1eef3'}">
      <span class="mi" aria-hidden="true">${esc((c.buyer_name || t('chBuyer')).trim().slice(0, 1).toUpperCase())}</span>
      <span class="mx"><b>${esc(c.buyer_name || t('chBuyer'))}${c.unread ? `<span class="mbadge">${c.unread}</span>` : ''}</b>
        <span class="chat-prev">${c.last_sender === 'seller' ? esc(t('chYou')) + ': ' : ''}${esc(c.preview)}</span></span>
      <span class="chat-meta">${c.waiting ? `<span class="chip pink">${t('chWaiting')}</span>` : ''}<small>${fmtTime(c.last_message_at)}</small></span></a>`).join('');
  const html = rows ? `<section class="mgroup"><div class="mlist" id="chatList">${rows}</div></section>`
    : `<section class="card" id="chatList"><p class="muted">${t('chNone')}</p></section>`;
  if (quiet && !$('#chatList')) return;
  $('#tab').innerHTML = html;
  document.querySelectorAll('[data-conv]').forEach(a => a.onclick = e => { e.preventDefault(); chatThread(id, a.dataset.conv); });
}

function chatCard(c) {
  if (!c) return '';
  if (c.type === 'order') return `<div class="chat-card"><span class="cc-ico" aria-hidden="true">📦</span><span><b>${esc(t('chOrder', {no: c.order_no}))}</b><small>${esc(c.status)} · ${pmoney(c.total)}</small></span></div>`;
  return `<div class="chat-card">${c.image_url ? `<img src="${esc(c.image_url)}" alt="">` : '<span class="cc-ico" aria-hidden="true">🛍️</span>'}
    <span><b>${esc(c.name)}</b><small>${t('chProduct')}${c.price ? ' · ' + pmoney(c.price) : ''}</small></span></div>`;
}

function chatMsgHTML(convID, m) {
  const mine = m.sender === 'seller';
  const who = mine ? (m.sender_name || t('chYou')) : (CHAT.conv?.buyer_name || t('chBuyer'));
  return `<div class="hp-msg chat-msg ${mine ? 'mine' : ''}" data-mid="${m.id}">${chatCard(m.card)}
    ${m.has_photo ? `<img class="chat-photo" alt="${t('chPhotoAlt')}" data-photo="${chatBase()}/${convID}/messages/${m.id}/photo">` : ''}
    ${m.body ? `<div class="hp-bubble">${esc(m.body)}</div>` : ''}
    <span class="muted small">${esc(who)} · ${fmtTime(m.created_at)}</span></div>`;
}

// Foto chat privat: dimuat dengan token (bukan URL publik).
async function chatLoadPhotos(root) {
  for (const img of root.querySelectorAll('img[data-photo]')) {
    const src = img.dataset.photo; img.removeAttribute('data-photo');
    try {
      if (!CHAT.photos[src]) {
        const r = await fetch(src, {headers: {Authorization: 'Bearer ' + P.token}});
        if (!r.ok) continue;
        CHAT.photos[src] = URL.createObjectURL(await r.blob());
      }
      img.src = CHAT.photos[src];
      img.onclick = () => window.open(img.src, '_blank', 'noopener');
    } catch {}
  }
}

async function chatThread(id, cid) {
  chatStop();
  const c = await api(`${chatBase()}/${cid}`);
  c.messages = c.messages || []; // kosong → tidak dikirim (omitempty)
  CHAT.conv = c; CHAT.last = c.messages.length ? c.messages[c.messages.length - 1].id : 0;
  $('#tab').innerHTML = `<a class="backlink" href="#/store/${id}/chat" id="chBack">${t('chBack')}</a>
    <section class="card chat-box" id="chatThread">
      <div class="spread"><h2>${esc(c.buyer_name || t('chBuyer'))}</h2></div>
      <div class="hp-thread chat-thread" id="chMsgs" aria-live="polite">${c.messages.map(m => chatMsgHTML(cid, m)).join('')}</div>
      <p class="muted small">${t('chHint')}</p>
      <label class="chat-sr" for="chBody">${t('chPlaceholder')}</label>
      <textarea id="chBody" name="body" rows="3" maxlength="2000" placeholder="${t('chPlaceholder')}"></textarea>
      <div class="row chat-actions">
        <button type="button" class="primary" id="chSend">${t('chSend')}</button>
        <label class="btn chat-file">📷 ${t('chPhoto')}<input type="file" id="chFile" accept="image/jpeg,image/png,image/webp" class="hidden"></label>
        <button type="button" id="chSuggest" class="hidden">${t('chSuggest')}</button>
      </div></section>`;
  $('#chBack').onclick = e => { e.preventDefault(); tabChat(id); };
  const box = $('#chMsgs');
  box.scrollTop = box.scrollHeight;
  chatLoadPhotos(box);
  chatSuggestSetup(id);
  $('#chSend').onclick = () => chatSend(id, cid);
  $('#chBody').onkeydown = e => { if (e.key === 'Enter' && (e.ctrlKey || e.metaKey)) chatSend(id, cid); };
  $('#chFile').onchange = async e => {
    const f = e.target.files[0]; e.target.value = '';
    if (!f) return;
    toast(t('chUploading'), 'info');
    try {
      const fd = new FormData(); fd.append('kind', 'chat_photo'); fd.append('file', f);
      const up = await api(`${chatBase()}/uploads`, {method: 'POST', body: fd});
      await chatPost(id, cid, {body: '', photo_key: up.key});
    } catch (err) { toast(errMsg(err), 'error'); }
  };
  chatPoll(() => chatPull(cid), 6000, 'chatThread');
}

// Tombol "✨ Saran balasan": memanggil asisten toko (endpoint dari fitur lain). Gagal / tidak ada → tombol disembunyikan diam-diam.
function chatSuggestSetup(id) {
  const b = $('#chSuggest');
  if (!b) return;
  const lastBuyer = [...(CHAT.conv?.messages || [])].reverse().find(m => m.sender === 'buyer' && m.body);
  // atribut hidden kalah oleh CSS tombol (display:inline-flex) → pakai kelas .hidden
  b.classList.toggle('hidden', CHAT.noAssist || !lastBuyer || CHAT.conv.messages[CHAT.conv.messages.length - 1].sender !== 'buyer');
  b.onclick = async () => {
    const q = [...CHAT.conv.messages].reverse().find(m => m.sender === 'buyer' && m.body);
    if (!q) return;
    b.disabled = true; b.textContent = t('chSuggesting');
    try {
      const body = {question: q.body};
      if (CHAT.conv.product_id) body.product_id = CHAT.conv.product_id;
      const r = await api(`/v1/seller/stores/${id}/assist/suggest`, {method: 'POST', body});
      if (!r || typeof r.reply !== 'string' || !r.reply.trim()) throw new Error('empty');
      const ta = $('#chBody'); if (ta) { ta.value = r.reply.slice(0, 2000); ta.focus(); }
      b.disabled = false; b.textContent = t('chSuggest');
    } catch {
      CHAT.noAssist = true; b.classList.add('hidden');
    }
  };
}

async function chatSend(id, cid) {
  const ta = $('#chBody'), body = ta.value.trim();
  if (!body) { ta.focus(); return; }
  if (body.length > 2000) { toast(t('chTooLong'), 'warning'); return; }
  $('#chSend').disabled = true;
  try { await chatPost(id, cid, {body}); ta.value = ''; } catch (e) { toast(errMsg(e), 'error'); }
  if ($('#chSend')) $('#chSend').disabled = false;
}

async function chatPost(id, cid, body) {
  await api(`${chatBase()}/${cid}/messages`, {method: 'POST', body});
  await chatPull(cid);
}

// chatPull: ambil pesan baru saja (?after=) lalu tambahkan ke bawah.
async function chatPull(cid) {
  const c = await api(`${chatBase()}/${cid}?after=${CHAT.last}`);
  const box = $('#chMsgs');
  if (!box || !c.messages || !c.messages.length) return;
  const nearBottom = box.scrollHeight - box.scrollTop - box.clientHeight < 80;
  c.messages.forEach(m => {
    if (box.querySelector(`[data-mid="${m.id}"]`)) return;
    box.insertAdjacentHTML('beforeend', chatMsgHTML(cid, m));
    CHAT.conv.messages.push(m);
  });
  CHAT.last = c.messages[c.messages.length - 1].id;
  if (nearBottom) box.scrollTop = box.scrollHeight;
  chatLoadPhotos(box);
  chatSuggestSetup(S.storeId);
}
