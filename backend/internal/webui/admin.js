'use strict';
P.authBase = '/admin/auth';
addStrings({
  en: {appName: 'Admin', logout: 'Sign out', queue: 'Store verification', notAdmin: 'This account is not a Phsar admin.',
    f_pending_review: 'Waiting for review', f_onboarding: 'Setting up', f_active: 'Active', store: 'Store', owner: 'Owner', kyc: 'KYC', pickups: 'Pickup points', submitted: 'Submitted',
    empty: 'Nothing here.', back: '← Queue', about: 'About', social: 'Social links', identity: 'Identity (KYC)', sellerType: 'Seller type', legalName: 'Legal name',
    viewDoc: 'Open ID document', tin: 'TIN', moc: 'MoC licence', status: 'Status', approve: 'Approve', reject: 'Reject', notePh: 'Reason (shown to the seller)',
    needNote: 'Write a reason before rejecting.', storeDecision: 'Store decision', storeHint: 'Approve after KYC and at least one pickup point are verified.',
    verified: 'Verified', notVerified: 'Not verified', contact: 'Contact', hours: 'Hours', openMap: 'Open in Google Maps', done: 'Saved',
    st_onboarding: 'Setting up', st_pending_review: 'Under review', st_active: 'Active', st_suspended: 'Suspended',
    t_individual_unregistered: 'Individual (not tax-registered)', t_individual_registered: 'Individual (tax-registered)', t_company: 'Company', kyc_missing: 'Not submitted', kyc_pending: 'Pending', kyc_verified: 'Verified', kyc_rejected: 'Rejected', err_not_ready: 'Verify KYC and a pickup point first.'},
  km: {appName: 'អ្នកគ្រប់គ្រង', logout: 'ចាកចេញ', queue: 'ផ្ទៀងផ្ទាត់ហាង', notAdmin: 'គណនីនេះមិនមែនជាអ្នកគ្រប់គ្រង Phsar ទេ។',
    f_pending_review: 'រង់ចាំពិនិត្យ', f_onboarding: 'កំពុងរៀបចំ', f_active: 'សកម្ម', store: 'ហាង', owner: 'ម្ចាស់', kyc: 'KYC', pickups: 'ទីតាំងទទួល', submitted: 'បានដាក់',
    empty: 'គ្មានទិន្នន័យ។', back: '← បញ្ជីរង់ចាំ', about: 'អំពី', social: 'បណ្តាញសង្គម', identity: 'អត្តសញ្ញាណ (KYC)', sellerType: 'ប្រភេទអ្នកលក់', legalName: 'ឈ្មោះស្របច្បាប់',
    viewDoc: 'បើកឯកសារអត្តសញ្ញាណ', tin: 'TIN', moc: 'អាជ្ញាប័ណ្ណ MoC', status: 'ស្ថានភាព', approve: 'អនុម័ត', reject: 'បដិសេធ', notePh: 'មូលហេតុ (អ្នកលក់មើលឃើញ)',
    needNote: 'សូមសរសេរមូលហេតុមុនបដិសេធ។', storeDecision: 'សម្រេចលើហាង', storeHint: 'អនុម័តបន្ទាប់ពីផ្ទៀងផ្ទាត់ KYC និងទីតាំងទទួលយ៉ាងហោចណាស់មួយ។',
    verified: 'បានផ្ទៀងផ្ទាត់', notVerified: 'មិនទាន់ផ្ទៀងផ្ទាត់', contact: 'ទំនាក់ទំនង', hours: 'ម៉ោង', openMap: 'បើកក្នុង Google Maps', done: 'បានរក្សាទុក',
    st_onboarding: 'កំពុងរៀបចំ', st_pending_review: 'កំពុងពិនិត្យ', st_active: 'សកម្ម', st_suspended: 'ផ្អាក',
    t_individual_unregistered: 'បុគ្គល (មិនទាន់ចុះបញ្ជីពន្ធ)', t_individual_registered: 'បុគ្គល (បានចុះបញ្ជីពន្ធ)', t_company: 'ក្រុមហ៊ុន', kyc_missing: 'មិនទាន់ដាក់', kyc_pending: 'រង់ចាំ', kyc_verified: 'បានផ្ទៀងផ្ទាត់', kyc_rejected: 'បានបដិសេធ', err_not_ready: 'សូមផ្ទៀងផ្ទាត់ KYC និងទីតាំងទទួលជាមុនសិន។'},
  zh: {appName: '管理后台', logout: '退出', queue: '店铺审核', notAdmin: '此账号不是 Phsar 管理员。',
    f_pending_review: '待审核', f_onboarding: '设置中', f_active: '已上线', store: '店铺', owner: '店主', kyc: 'KYC', pickups: '取件点', submitted: '提交时间',
    empty: '暂无数据。', back: '← 审核队列', about: '介绍', social: '社交链接', identity: '身份验证 (KYC)', sellerType: '卖家类型', legalName: '证件姓名',
    viewDoc: '打开身份证件', tin: '税号', moc: '商业部许可证', status: '状态', approve: '通过', reject: '拒绝', notePh: '原因（卖家可见）',
    needNote: '拒绝前请填写原因。', storeDecision: '店铺审核结果', storeHint: 'KYC 与至少一个取件点核实后再通过。',
    verified: '已核实', notVerified: '未核实', contact: '联系人', hours: '时间', openMap: '在 Google 地图中打开', done: '已保存',
    st_onboarding: '设置中', st_pending_review: '审核中', st_active: '已上线', st_suspended: '已暂停',
    t_individual_unregistered: '个人（未税务登记）', t_individual_registered: '个人（已税务登记）', t_company: '公司', kyc_missing: '未提交', kyc_pending: '待审核', kyc_verified: '已验证', kyc_rejected: '已拒绝', err_not_ready: '请先核实 KYC 和取件点。'},
});

const A = {filter: 'pending_review', map: null};
const chip = (cls, text) => `<span class="chip ${cls}">${esc(text)}</span>`;
const kycCls = s => ({pending: 'warn', verified: 'ok', rejected: 'bad'})[s] || '';

function shell(inner) {
  return `<header class="top"><a class="brand" href="#/">${brandSVG}<span><b>Phsar</b><small>${esc(t('appName')).toUpperCase()}</small></span></a>
    <div class="top-actions">${langSwitch()}${P.token ? `<button type="button" onclick="logout()">${t('logout')}</button>` : ''}</div></header><main>${inner}</main>`;
}

async function render() {
  if (!P.token) {
    $('#app').innerHTML = shell(loginView('Phsar ' + t('appName')));
    bindLogin(() => { location.hash = '#/'; render(); });
    return;
  }
  const m = location.hash.match(/^#\/store\/([0-9a-f-]+)/);
  try { if (m) await renderDetail(m[1]); else await renderQueue(); }
  catch (e) { $('#app').innerHTML = shell(`<div class="banner bad">${e.code === 'not_admin' ? t('notAdmin') : esc(errMsg(e))}</div>`); }
}

async function renderQueue() {
  const {stores} = await api('/admin/api/stores?status=' + A.filter);
  const rows = stores.map(s => `<tr><td><a href="#/store/${s.id}"><b>${esc(s.name)}</b></a><br><span class="muted">${esc(s.slug)}</span></td>
    <td>${esc(s.owner_name || '')} ${esc(s.owner_phone)}</td><td>${chip(kycCls(s.kyc), t('kyc_' + s.kyc))}</td><td>${s.pickup_points}</td>
    <td class="muted">${s.submitted_at ? new Date(s.submitted_at).toLocaleString() : '—'}</td></tr>`).join('');
  $('#app').innerHTML = shell(`<h1>${t('queue')}</h1>
    <div class="tabs">${['pending_review', 'onboarding', 'active'].map(f => `<button type="button" aria-selected="${A.filter === f}" onclick="A.filter='${f}';render()">${t('f_' + f)}</button>`).join('')}</div>
    <section class="card" style="overflow-x:auto">${stores.length ? `<table><thead><tr><th>${t('store')}</th><th>${t('owner')}</th><th>${t('kyc')}</th><th>${t('pickups')}</th><th>${t('submitted')}</th></tr></thead><tbody>${rows}</tbody></table>`
      : `<p class="muted">${t('empty')}</p>`}</section>`);
}

function decisionBox(id) {
  return `<div class="row"><input id="note_${id}" placeholder="${t('notePh')}" style="flex:1;min-width:200px">
    <button type="button" class="primary" data-act="approve" data-id="${id}">${t('approve')}</button>
    <button type="button" class="danger" data-act="reject" data-id="${id}">${t('reject')}</button></div>`;
}

async function renderDetail(id) {
  const d = await api('/admin/api/stores/' + id);
  const st = d.store, k = d.kyc, pending = st.status === 'pending_review';
  const days = ['mon', 'tue', 'wed', 'thu', 'fri', 'sat', 'sun'];
  $('#app').innerHTML = shell(`<a href="#/">${t('back')}</a>
    <div class="spread"><div class="row"><img class="logo" src="${esc(st.logo_url || '')}" alt="" onerror="this.style.visibility='hidden'">
      <div><h1>${esc(st.name)}</h1><p class="muted">${esc(st.slug)}</p></div></div>${chip({active: 'ok', pending_review: 'warn'}[st.status] || 'pink', t('st_' + st.status))}</div>
    <section class="card"><h2>${t('about')}</h2>${Object.entries(st.about).map(([l, a]) => `<p><span class="chip">${l}</span> ${esc(a)}</p>`).join('') || '<p class="muted">—</p>'}
      <p class="muted">${t('social')}: ${Object.entries(st.social_links || {}).map(([s, v]) => `${esc(s)} ${esc(v)}`).join(' · ') || '—'}</p></section>
    <section class="card"><div class="spread"><h2>${t('identity')}</h2>${chip(kycCls(k?.verification_status), t('kyc_' + (k?.verification_status || 'missing')))}</div>
      ${k ? `<div class="grid"><p><span class="muted">${t('sellerType')}</span><br>${esc(t('t_' + k.seller_type))}</p><p><span class="muted">${t('legalName')}</span><br><b>${esc(k.legal_name)}</b></p>
        <p><span class="muted">${t('tin')}</span><br>${esc(k.tin || '—')}</p><p><span class="muted">${t('moc')}</span><br>${esc(k.moc_license_no || '—')}</p></div>
        ${d.id_document_url ? `<div><button type="button" id="docBtn">${t('viewDoc')}</button></div>` : ''}
        ${k.rejection_note ? `<div class="banner bad">${esc(k.rejection_note)}</div>` : ''}
        ${k.verification_status !== 'verified' ? decisionBox('kyc') : ''}` : `<p class="muted">${t('kyc_missing')}</p>`}</section>
    <section class="card"><h2>${t('pickups')}</h2><div id="amap" class="map"></div>
      ${d.pickup_points.map(p => `<div class="card" style="padding:14px"><div class="spread"><b>${esc(p.label)}</b>${chip(p.verified ? 'ok' : 'warn', p.verified ? t('verified') : t('notVerified'))}</div>
        <p>${esc(p.address_line)}</p><p class="muted">${t('contact')}: ${esc(p.contact_name)} · ${esc(p.contact_phone)}</p>
        <p class="muted">${t('hours')}: ${days.filter(x => p.pickup_hours[x]).map(x => `${x} ${p.pickup_hours[x].join('–')}`).join(' · ')}</p>
        ${p.note ? `<p class="muted">${esc(p.note)}</p>` : ''}${p.photo_url ? `<img src="${esc(p.photo_url)}" alt="" style="max-width:240px;border-radius:16px">` : ''}
        <a href="https://www.google.com/maps?q=${p.lat},${p.lng}" target="_blank" rel="noopener">${t('openMap')}</a>
        ${p.rejection_note ? `<div class="banner bad">${esc(p.rejection_note)}</div>` : ''}${p.verified ? '' : decisionBox('pp_' + p.id)}</div>`).join('')}</section>
    ${pending ? `<section class="card"><h2>${t('storeDecision')}</h2><p class="muted">${t('storeHint')}</p>${decisionBox('store')}</section>` : ''}`);

  if (d.pickup_points.length) {
    A.map = L.map('amap');
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {maxZoom: 19, attribution: '© OpenStreetMap'}).addTo(A.map);
    const pts = d.pickup_points.map(p => (L.marker([p.lat, p.lng]).addTo(A.map).bindPopup(esc(p.label)), [p.lat, p.lng]));
    A.map.fitBounds(pts, {maxZoom: 17, padding: [30, 30]});
  } else $('#amap').classList.add('hidden');

  if ($('#docBtn')) $('#docBtn').onclick = async () => {
    // Dokumen privat butuh token → ambil sebagai blob lalu buka di tab baru.
    const r = await fetch(d.id_document_url, {headers: {Authorization: 'Bearer ' + P.token}});
    if (!r.ok) { toast(r.statusText); return; }
    window.open(URL.createObjectURL(await r.blob()), '_blank', 'noopener');
  };
  document.querySelectorAll('[data-act]').forEach(b => b.onclick = async () => {
    const box = b.dataset.id, decision = b.dataset.act, note = $('#note_' + box).value.trim();
    if (decision === 'reject' && !note) { toast(t('needNote')); return; }
    const url = box === 'kyc' ? `/admin/api/stores/${id}/kyc` : box === 'store' ? `/admin/api/stores/${id}/decision`
      : `/admin/api/pickup-points/${box.slice(3)}/decision`;
    try { await api(url, {method: 'POST', body: {decision, note}}); toast(t('done')); render(); } catch (e) { toast(errMsg(e)); }
  });
}

window.addEventListener('hashchange', render);
(async () => { document.documentElement.lang = P.lang; await refreshSession(); render(); })();
