'use strict';
// Seller Center: lonceng notifikasi, banner, tim (staf) + halaman undangan.
addStrings({
  en: {notifications: 'Notifications', noNotifications: 'No notifications yet.', markAllRead: 'Mark all as read', tgHint: 'Sign in with Telegram to also get these in Telegram.',
    tabTeam: 'Team', staffNote: 'You are staff in this store: you can add and edit products. Store settings are managed by the owner.',
    teamIntro: 'Staff can add and edit products. Only you (the owner) can change the store profile, identity, pickup points and team.',
    owner: 'Owner', staff: 'Staff', remove: 'Remove', confirmRemove: 'Remove {name} from this store?', inviteStaff: 'Invite staff',
    inviteHint: 'Create a link and send it to your staff (e.g. in Telegram). It works once and expires in 7 days.', createLink: 'Create invite link',
    copy: 'Copy', copied: 'Copied', openInvites: 'Active invite links', cancelLink: 'Cancel', expires: 'expires {d}', linkOnce: 'Copy this link now — it is shown only once.',
    inviteTitle: 'You’re invited', inviteText: 'Join {store} as staff. You will be able to add and edit its products.', acceptInvite: 'Join store',
    alreadyMember: 'You are already a member of this store.', inviteGone: 'This invite link is no longer valid. Ask the store owner for a new one.',
    leaveStore: 'Leave this store', confirmLeave: 'Leave this store? You will lose access.',
    uploadBanner: 'Upload banner', changeBanner: 'Change banner', bannerHint: 'Store banner — shown at the top of your shop page', bannerSize: 'Wide photo, e.g. 1600 × 600',
    err_invite_invalid: 'This invite link is no longer valid.', err_too_many_staff: 'A store can have up to 10 staff.', err_owner_only: 'Only the store owner can do this.',
    stockPrice: 'Stock & price', saveChanges: 'Save', stockSaved: 'Stock & price updated', quickHint: 'Changes go live right away — no review needed.'},
  km: {notifications: 'ការជូនដំណឹង', noNotifications: 'មិនទាន់មានការជូនដំណឹងទេ។', markAllRead: 'សម្គាល់ថាបានអានទាំងអស់', tgHint: 'ចូលតាម Telegram ដើម្បីទទួលការជូនដំណឹងក្នុង Telegram ផងដែរ។',
    tabTeam: 'ក្រុមការងារ', staffNote: 'អ្នកជាបុគ្គលិកនៃហាងនេះ៖ អាចបន្ថែម និងកែផលិតផល។ ការកំណត់ហាងគ្រប់គ្រងដោយម្ចាស់ហាង។',
    teamIntro: 'បុគ្គលិកអាចបន្ថែម និងកែផលិតផល។ មានតែអ្នក (ម្ចាស់) ដែលអាចកែព័ត៌មានហាង អត្តសញ្ញាណ ទីតាំងទទួល និងក្រុមការងារ។',
    owner: 'ម្ចាស់', staff: 'បុគ្គលិក', remove: 'ដកចេញ', confirmRemove: 'ដក {name} ចេញពីហាងនេះ?', inviteStaff: 'អញ្ជើញបុគ្គលិក',
    inviteHint: 'បង្កើតតំណ ហើយផ្ញើទៅបុគ្គលិក (ឧ. តាម Telegram)។ ប្រើបានតែម្តង និងផុតកំណត់ក្នុង 7 ថ្ងៃ។', createLink: 'បង្កើតតំណអញ្ជើញ',
    copy: 'ចម្លង', copied: 'បានចម្លង', openInvites: 'តំណអញ្ជើញសកម្ម', cancelLink: 'បោះបង់', expires: 'ផុតកំណត់ {d}', linkOnce: 'សូមចម្លងតំណនេះឥឡូវ — វាបង្ហាញតែម្តងប៉ុណ្ណោះ។',
    inviteTitle: 'អ្នកត្រូវបានអញ្ជើញ', inviteText: 'ចូលរួមជាបុគ្គលិករបស់ {store}។ អ្នកនឹងអាចបន្ថែម និងកែផលិតផលរបស់ហាង។', acceptInvite: 'ចូលរួមហាង',
    alreadyMember: 'អ្នកជាសមាជិកនៃហាងនេះរួចហើយ។', inviteGone: 'តំណអញ្ជើញនេះលែងប្រើបានហើយ។ សូមស្នើតំណថ្មីពីម្ចាស់ហាង។',
    leaveStore: 'ចាកចេញពីហាងនេះ', confirmLeave: 'ចាកចេញពីហាងនេះ? អ្នកនឹងបាត់បង់សិទ្ធិចូល។',
    uploadBanner: 'ផ្ទុកបដា', changeBanner: 'ប្តូរបដា', bannerHint: 'បដាហាង — បង្ហាញនៅខាងលើទំព័រហាង', bannerSize: 'រូបថតទទឹង ឧ. 1600 × 600',
    err_invite_invalid: 'តំណអញ្ជើញនេះលែងប្រើបានហើយ។', err_too_many_staff: 'ហាងមួយមានបុគ្គលិកបានច្រើនបំផុត 10 នាក់។', err_owner_only: 'មានតែម្ចាស់ហាងប៉ុណ្ណោះដែលអាចធ្វើបាន។',
    stockPrice: 'ស្តុក និងតម្លៃ', saveChanges: 'រក្សាទុក', stockSaved: 'បានកែស្តុក និងតម្លៃ', quickHint: 'ការកែប្រែមានប្រសិទ្ធភាពភ្លាម — មិនចាំបាច់ពិនិត្យ។'},
  zh: {notifications: '通知', noNotifications: '暂无通知。', markAllRead: '全部标为已读', tgHint: '使用 Telegram 登录即可在 Telegram 中收到这些通知。',
    tabTeam: '团队', staffNote: '你是这家店的店员：可以添加和编辑商品。店铺设置由店主管理。',
    teamIntro: '店员可以添加和编辑商品。只有你（店主）可以修改店铺资料、身份验证、取件点和团队。',
    owner: '店主', staff: '店员', remove: '移除', confirmRemove: '把 {name} 从店铺移除？', inviteStaff: '邀请店员',
    inviteHint: '生成链接并发给店员（例如通过 Telegram）。链接只能使用一次，7 天后失效。', createLink: '生成邀请链接',
    copy: '复制', copied: '已复制', openInvites: '有效的邀请链接', cancelLink: '取消', expires: '{d} 失效', linkOnce: '请现在复制此链接——只显示一次。',
    inviteTitle: '你收到邀请', inviteText: '以店员身份加入 {store}，可以添加和编辑店铺商品。', acceptInvite: '加入店铺',
    alreadyMember: '你已经是这家店的成员。', inviteGone: '此邀请链接已失效，请向店主索取新链接。',
    leaveStore: '退出此店铺', confirmLeave: '退出此店铺？你将失去访问权限。',
    uploadBanner: '上传横幅', changeBanner: '更换横幅', bannerHint: '店铺横幅——显示在店铺页面顶部', bannerSize: '横向图片，例如 1600 × 600',
    err_invite_invalid: '此邀请链接已失效。', err_too_many_staff: '每家店最多 10 名店员。', err_owner_only: '只有店主可以执行此操作。',
    stockPrice: '库存和价格', saveChanges: '保存', stockSaved: '库存和价格已更新', quickHint: '修改立即生效，无需审核。'},
});

// ---- Lonceng notifikasi ----
const NB = {open: false, data: null};
function bellButton() {
  return `<div class="bellwrap"><button type="button" class="bell" id="bellBtn" aria-label="${t('notifications')}" aria-expanded="${NB.open}" onclick="toggleBell()">
    <svg width="20" height="20" viewBox="0 0 24 24" aria-hidden="true"><path fill="currentColor" d="M12 22a2.5 2.5 0 0 0 2.4-2h-4.8A2.5 2.5 0 0 0 12 22Zm7-6V11a7 7 0 0 0-5.5-6.8V3a1.5 1.5 0 0 0-3 0v1.2A7 7 0 0 0 5 11v5l-2 2v1h18v-1z"/></svg>
    <span class="badge hidden" id="bellCount"></span></button><div class="bellpanel hidden" id="bellPanel" role="dialog" aria-label="${t('notifications')}"></div></div>`;
}
async function loadBell() {
  if (!P.token || !$('#bellBtn')) return;
  try { NB.data = await api('/v1/seller/notifications?limit=30'); } catch { return; }
  const c = $('#bellCount');
  if (c) { c.textContent = NB.data.unread > 9 ? '9+' : NB.data.unread; c.classList.toggle('hidden', !NB.data.unread); }
  if (NB.open) drawBell();
}
function drawBell() {
  const panel = $('#bellPanel');
  if (!panel || !NB.data) return;
  panel.classList.toggle('hidden', !NB.open);
  const list = NB.data.notifications;
  panel.innerHTML = `<div class="spread"><b>${t('notifications')}</b>${NB.data.unread ? `<button type="button" class="linkbtn" onclick="readAll()">${t('markAllRead')}</button>` : ''}</div>
    ${list.length ? list.map(n => `<a class="notif ${n.read ? '' : 'unread'}" href="${esc(n.link || '#/')}" onclick="readOne(${n.id})">
      <b>${esc(n.title)}</b>${n.body ? `<span>${esc(n.body)}</span>` : ''}<small class="muted">${new Date(n.created_at).toLocaleString()}</small></a>`).join('')
      : `<p class="muted">${t('noNotifications')}</p>`}
    <p class="muted" style="font-size:12px">${t('tgHint')}</p>`;
}
function toggleBell() {
  NB.open = !NB.open;
  $('#bellBtn')?.setAttribute('aria-expanded', NB.open);
  if (NB.data) drawBell(); else $('#bellPanel')?.classList.toggle('hidden', !NB.open);
  loadBell();
}
async function readOne(id) {
  NB.open = false;
  try { await api('/v1/seller/notifications/read', {method: 'POST', body: {ids: [id]}}); } catch {}
  loadBell();
}
async function readAll() {
  try { await api('/v1/seller/notifications/read', {method: 'POST', body: {all: true}}); } catch {}
  loadBell();
}
document.addEventListener('click', e => {
  if (NB.open && !e.target.closest('.bellwrap')) { NB.open = false; $('#bellPanel')?.classList.add('hidden'); }
});
setInterval(() => { if (document.visibilityState === 'visible') loadBell(); }, 60000);

// ---- Tab tim (owner) ----
const TM = {lastLink: null};
const inviteURL = token => `${location.origin}/seller/#/invite/${token}`;
async function tabTeam(id) {
  const team = await api(`/v1/seller/stores/${id}/team`);
  $('#tab').innerHTML = `<p class="muted">${t('teamIntro')}</p>
    <section class="card">
      ${team.members.map(m => { const title = m.name || m.telegram || m.phone || '—'; const sub = [m.phone, m.telegram].filter(x => x && x !== title).join(' · ');
        return `<div class="spread member"><div><b>${esc(title)}</b>${sub ? `<span class="muted">${esc(sub)}</span>` : ''}</div>
        <div class="row"><span class="chip ${m.role === 'owner' ? 'pink' : ''}">${t(m.role)}</span>
          ${m.role === 'staff' ? `<button type="button" class="danger" onclick="removeMember('${id}','${m.user_id}','${esc(m.name || m.phone)}')">${t('remove')}</button>` : ''}</div></div>`; }).join('')}
    </section>
    <section class="card"><h2>${t('inviteStaff')}</h2><p class="muted">${t('inviteHint')}</p>
      ${TM.lastLink ? `<div class="banner ok"><b>${t('linkOnce')}</b><div class="row" style="margin-top:8px">
        <input id="inviteLink" readonly value="${esc(TM.lastLink)}" style="flex:1;min-width:220px"><button type="button" onclick="copyInvite()">${t('copy')}</button></div></div>` : ''}
      <div><button type="button" class="primary" id="newInvite">${t('createLink')}</button></div>
      ${team.invites.length ? `<b>${t('openInvites')}</b>${team.invites.map(i => `<div class="spread"><span class="muted">${t('expires', {d: new Date(i.expires_at).toLocaleDateString()})}</span>
        <button type="button" onclick="revokeInvite('${id}','${i.id}')">${t('cancelLink')}</button></div>`).join('')}` : ''}
    </section>`;
  $('#newInvite').onclick = async () => {
    try { const inv = await api(`/v1/seller/stores/${id}/invites`, {method: 'POST'}); TM.lastLink = inviteURL(inv.token); render(); }
    catch (e) { showError(e); }
  };
}
async function copyInvite() {
  const el = $('#inviteLink');
  try { await navigator.clipboard.writeText(el.value); } catch { el.select(); document.execCommand('copy'); }
  toast(t('copied'));
}
async function revokeInvite(sid, iid) {
  try { await api(`/v1/seller/stores/${sid}/invites/${iid}`, {method: 'DELETE'}); TM.lastLink = null; render(); } catch (e) { toast(errMsg(e)); }
}
async function removeMember(sid, uid, name) {
  if (!confirm(t('confirmRemove', {name}))) return;
  try { await api(`/v1/seller/stores/${sid}/members/${uid}`, {method: 'DELETE'}); render(); } catch (e) { toast(errMsg(e)); }
}
window.addEventListener('hashchange', () => { TM.lastLink = null; });

// ---- Halaman undangan: #/invite/<token> ----
async function renderInvite(token) {
  let info;
  try { info = await api('/v1/seller/invites/' + token); }
  catch (e) {
    $('#app').innerHTML = shell(`<section class="card login"><h1>${t('inviteTitle')}</h1><div class="banner bad">${t('inviteGone')}</div><a class="btn" href="#/">${t('myStores')}</a></section>`);
    return;
  }
  $('#app').innerHTML = shell(`<section class="card login">
    <img class="logo" src="${esc(info.logo_url || '')}" alt="" onerror="this.style.visibility='hidden'">
    <h1>${t('inviteTitle')}</h1>
    <p>${esc(t('inviteText', {store: info.store_name}))}</p>
    ${info.already_member ? `<div class="banner info">${t('alreadyMember')}</div><a class="btn primary" href="#/store/${info.store_id}">${t('open')}</a>`
      : `<button type="button" class="primary" id="acceptBtn">${t('acceptInvite')}</button>`}
  </section>`);
  if ($('#acceptBtn')) $('#acceptBtn').onclick = async () => {
    try { const r = await api(`/v1/seller/invites/${token}/accept`, {method: 'POST'}); location.hash = '#/store/' + r.store_id; }
    catch (e) { showError(e); }
  };
}

// ---- Edit stok & harga cepat (dari daftar produk) ----
async function openStockDialog(sid, pid) {
  let p;
  try { p = await api(`/v1/seller/stores/${sid}/products/${pid}`); } catch (e) { toast(errMsg(e)); return; }
  const dlg = document.createElement('dialog');
  dlg.className = 'stockdlg';
  dlg.innerHTML = `<form method="dialog"><div class="spread"><h2>${t('stockPrice')}</h2><button value="cancel" aria-label="${t('cancel')}">✕</button></div>
    <p class="muted">${esc(p.translations[p.source_lang]?.name || '')}</p><p class="muted" style="font-size:13px">${t('quickHint')}</p>
    <div style="overflow-x:auto"><table class="vtable"><thead><tr><th>${t('variant')}</th><th>${t('price')}</th><th>${t('comparePrice')}</th><th>${t('stockQty')}</th><th>${t('available')}</th></tr></thead><tbody>
    ${p.variants.map(v => `<tr data-vid="${v.id}"><td>${esc(v.options.join(' / ') || '—')}</td>
      <td><input data-q="price" inputmode="decimal" value="${esc(v.price)}"></td><td><input data-q="compare_at_price" inputmode="decimal" value="${esc(v.compare_at_price || '')}"></td>
      <td><input data-q="stock_qty" type="number" min="0" value="${v.stock_qty}"></td><td><input data-q="on" type="checkbox" ${v.status === 'active' ? 'checked' : ''}></td></tr>`).join('')}
    </tbody></table></div>
    <div class="row" data-errors-anchor><button type="button" class="primary" id="stockSave">${t('saveChanges')}</button><button value="cancel">${t('cancel')}</button></div></form>`;
  document.body.append(dlg);
  dlg.addEventListener('close', () => dlg.remove());
  dlg.showModal();
  dlg.querySelector('#stockSave').onclick = async () => {
    const rows = [...dlg.querySelectorAll('tr[data-vid]')];
    const rules = [];
    for (const tr of rows) {
      const price = tr.querySelector('[data-q=price]'), cmp = tr.querySelector('[data-q=compare_at_price]');
      rules.push([price, t('positive'), el => Number(el.value) > 0]);
      rules.push([cmp, t('compareHigher'), el => el.value.trim() === '' || Number(el.value) > Number(price.value)]);
      rules.push([tr.querySelector('[data-q=stock_qty]'), t('required'), el => el.value !== '' && Number(el.value) >= 0]);
    }
    if (!checkRequired(rules, dlg.querySelector('[data-errors-anchor]'))) return;
    const body = {variants: rows.map(tr => ({id: tr.dataset.vid, price: tr.querySelector('[data-q=price]').value.trim(),
      compare_at_price: tr.querySelector('[data-q=compare_at_price]').value.trim(), stock_qty: Number(tr.querySelector('[data-q=stock_qty]').value),
      status: tr.querySelector('[data-q=on]').checked ? 'active' : 'inactive'}))};
    try {
      await api(`/v1/seller/stores/${sid}/products/${pid}/variants`, {method: 'PATCH', body});
      dlg.close(); toast(t('stockSaved')); render();
    } catch (e) {
      showError(e, f => { const m = f.match(/^variants\[(\d+)\]\.(\w+)/); return m ? rows[Number(m[1])]?.querySelector(`[data-q=${m[2]}]`) : null; },
        dlg.querySelector('[data-errors-anchor]'));
    }
  };
}
