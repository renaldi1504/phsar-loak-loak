'use strict';
P.authBase = '/seller/auth';
addStrings({
  en: {appName: 'Seller Center', logout: 'Sign out', myStores: 'My stores', createStore: 'Create your store', createSub: 'Fill in the basics. You can edit everything later.',
    storeName: 'Store name', storeAbout: 'About your store', writtenIn: 'Written in', aboutHint: 'Write in one language. We translate it to the other two automatically.',
    create: 'Create store', noStores: 'You don’t have a store yet.', open: 'Open', back: '← My stores',
    st_onboarding: 'Setting up', st_pending_review: 'Under review', st_active: 'Active', st_suspended: 'Suspended', st_closed: 'Closed',
    checklist: 'Before you can sell', ck_logo: 'Store logo', ck_about: 'Store description', ck_kyc: 'Identity check (KYC)', ck_pickup: 'Pickup point in Sihanoukville',
    ck_review: 'Approved by Phsar', submit: 'Send for review', submitHint: 'Phsar checks your identity and pickup point, usually within 1–2 working days.',
    reviewing: 'Phsar is reviewing your store. You can’t edit it until the review is done.', activeMsg: 'Your store is active. Next: add products.',
    needsFix: 'Phsar asked for changes:', tabProfile: 'Profile', tabKyc: 'Identity (KYC)', tabPickup: 'Pickup points',
    logo: 'Logo', uploadLogo: 'Upload logo', social: 'Social links', save: 'Save', saved: 'Saved',
    kycIntro: 'Your documents are only visible to the Phsar verification team.', sellerType: 'Seller type',
    t_individual_unregistered: 'Individual (not tax-registered)', t_individual_registered: 'Individual (tax-registered)', t_company: 'Company',
    legalName: 'Full legal name (as on ID)', idDoc: 'ID card or passport photo', idDocDone: 'Document uploaded', tin: 'TIN (tax ID)', mocNo: 'Ministry of Commerce e-commerce licence no.',
    optional: 'optional', kycStatus: 'Status', kyc_missing: 'Not submitted', kyc_pending: 'Waiting for review', kyc_verified: 'Verified', kyc_rejected: 'Rejected',
    pickupIntro: 'The Phsar courier picks up packages here. Tap the map to place the pin exactly at your door.',
    addPickup: 'Add pickup point', editPickup: 'Edit pickup point', label: 'Name', labelPh: 'e.g. Shop Ochheuteal', address: 'Address', addressPh: 'Street, landmark',
    contactName: 'Contact person', contactPhone: 'Contact phone', hours: 'Pickup hours', note: 'Directions for the courier', photo: 'Photo of the shop front',
    pinInside: 'Inside our delivery area', pinOutside: 'Outside our delivery area (Sihanoukville only for now)', pinNone: 'Tap the map to set the pin',
    verified: 'Verified', notVerified: 'Waiting for check', del: 'Delete', edit: 'Edit', cancel: 'Cancel', default: 'Default', makeDefault: 'Use as default',
    d_mon: 'Mon', d_tue: 'Tue', d_wed: 'Wed', d_thu: 'Thu', d_fri: 'Fri', d_sat: 'Sat', d_sun: 'Sun',
    err_validation: 'Please check the highlighted field.', err_slug_taken: 'That store address is taken. Try another name.', err_outside_service_area: 'That location is outside Sihanoukville.',
    err_store_under_review: 'Your store is under review.', err_upload_type: 'This file type is not allowed.', err_upload_too_large: 'The file is too large.', err_store_incomplete: 'Your store is not complete yet.'},
  km: {appName: 'មជ្ឈមណ្ឌលអ្នកលក់', logout: 'ចាកចេញ', myStores: 'ហាងរបស់ខ្ញុំ', createStore: 'បង្កើតហាងរបស់អ្នក', createSub: 'បំពេញព័ត៌មានមូលដ្ឋាន។ អាចកែបានពេលក្រោយ។',
    storeName: 'ឈ្មោះហាង', storeAbout: 'អំពីហាង', writtenIn: 'សរសេរជា', aboutHint: 'សរសេរតែមួយភាសា។ យើងបកប្រែទៅភាសាពីរទៀតដោយស្វ័យប្រវត្តិ។',
    create: 'បង្កើតហាង', noStores: 'អ្នកមិនទាន់មានហាងទេ។', open: 'បើក', back: '← ហាងរបស់ខ្ញុំ',
    st_onboarding: 'កំពុងរៀបចំ', st_pending_review: 'កំពុងត្រួតពិនិត្យ', st_active: 'សកម្ម', st_suspended: 'ផ្អាក', st_closed: 'បិទ',
    checklist: 'មុនពេលអាចលក់', ck_logo: 'ឡូហ្គោហាង', ck_about: 'ការពិពណ៌នាហាង', ck_kyc: 'ផ្ទៀងផ្ទាត់អត្តសញ្ញាណ (KYC)', ck_pickup: 'ទីតាំងទទួលកញ្ចប់នៅក្រុងព្រះសីហនុ',
    ck_review: 'Phsar បានអនុម័ត', submit: 'ផ្ញើឱ្យត្រួតពិនិត្យ', submitHint: 'Phsar ពិនិត្យអត្តសញ្ញាណ និងទីតាំងទទួលកញ្ចប់ ជាធម្មតាក្នុង 1–2 ថ្ងៃធ្វើការ។',
    reviewing: 'Phsar កំពុងត្រួតពិនិត្យហាងរបស់អ្នក។ មិនអាចកែបានទេរហូតដល់ពិនិត្យរួច។', activeMsg: 'ហាងរបស់អ្នកសកម្មហើយ។ បន្ទាប់៖ បន្ថែមផលិតផល។',
    needsFix: 'Phsar ស្នើឱ្យកែ៖', tabProfile: 'ព័ត៌មានហាង', tabKyc: 'អត្តសញ្ញាណ (KYC)', tabPickup: 'ទីតាំងទទួលកញ្ចប់',
    logo: 'ឡូហ្គោ', uploadLogo: 'ផ្ទុកឡូហ្គោ', social: 'បណ្តាញសង្គម', save: 'រក្សាទុក', saved: 'បានរក្សាទុក',
    kycIntro: 'ឯកសាររបស់អ្នកមើលឃើញតែក្រុមផ្ទៀងផ្ទាត់ Phsar ប៉ុណ្ណោះ។', sellerType: 'ប្រភេទអ្នកលក់',
    t_individual_unregistered: 'បុគ្គល (មិនទាន់ចុះបញ្ជីពន្ធ)', t_individual_registered: 'បុគ្គល (បានចុះបញ្ជីពន្ធ)', t_company: 'ក្រុមហ៊ុន',
    legalName: 'ឈ្មោះពេញ (ដូចអត្តសញ្ញាណប័ណ្ណ)', idDoc: 'រូបថតអត្តសញ្ញាណប័ណ្ណ ឬលិខិតឆ្លងដែន', idDocDone: 'បានផ្ទុកឯកសារ', tin: 'លេខអត្តសញ្ញាណកម្មសារពើពន្ធ (TIN)', mocNo: 'លេខអាជ្ញាប័ណ្ណពាណិជ្ជកម្មអេឡិចត្រូនិក',
    optional: 'មិនចាំបាច់', kycStatus: 'ស្ថានភាព', kyc_missing: 'មិនទាន់ដាក់', kyc_pending: 'រង់ចាំពិនិត្យ', kyc_verified: 'បានផ្ទៀងផ្ទាត់', kyc_rejected: 'បានបដិសេធ',
    pickupIntro: 'អ្នកដឹក Phsar មកទទួលកញ្ចប់នៅទីនេះ។ ចុចលើផែនទីដើម្បីដាក់ម្ជុលត្រង់ទ្វារហាងរបស់អ្នក។',
    addPickup: 'បន្ថែមទីតាំងទទួល', editPickup: 'កែទីតាំងទទួល', label: 'ឈ្មោះ', labelPh: 'ឧ. ហាងអូរឈើទាល', address: 'អាសយដ្ឋាន', addressPh: 'ផ្លូវ សម្គាល់',
    contactName: 'អ្នកទំនាក់ទំនង', contactPhone: 'លេខទូរស័ព្ទ', hours: 'ម៉ោងទទួលកញ្ចប់', note: 'ការណែនាំសម្រាប់អ្នកដឹក', photo: 'រូបថតមុខហាង',
    pinInside: 'នៅក្នុងតំបន់ដឹកជញ្ជូន', pinOutside: 'នៅក្រៅតំបន់ (បច្ចុប្បន្នតែក្រុងព្រះសីហនុ)', pinNone: 'ចុចលើផែនទីដើម្បីដាក់ម្ជុល',
    verified: 'បានផ្ទៀងផ្ទាត់', notVerified: 'រង់ចាំពិនិត្យ', del: 'លុប', edit: 'កែ', cancel: 'បោះបង់', default: 'លំនាំដើម', makeDefault: 'ប្រើជាលំនាំដើម',
    d_mon: 'ច', d_tue: 'អ', d_wed: 'ព', d_thu: 'ព្រ', d_fri: 'សុ', d_sat: 'ស', d_sun: 'អា',
    err_validation: 'សូមពិនិត្យវាលដែលបានសម្គាល់។', err_slug_taken: 'អាសយដ្ឋានហាងនេះមានគេប្រើហើយ។', err_outside_service_area: 'ទីតាំងនៅក្រៅក្រុងព្រះសីហនុ។',
    err_store_under_review: 'ហាងកំពុងត្រួតពិនិត្យ។', err_upload_type: 'ប្រភេទឯកសារមិនអនុញ្ញាត។', err_upload_too_large: 'ឯកសារធំពេក។', err_store_incomplete: 'ហាងមិនទាន់គ្រប់ព័ត៌មាន។'},
  zh: {appName: '卖家中心', logout: '退出', myStores: '我的店铺', createStore: '创建你的店铺', createSub: '先填写基本信息，之后都可以修改。',
    storeName: '店铺名称', storeAbout: '店铺介绍', writtenIn: '书写语言', aboutHint: '只需用一种语言填写，系统会自动翻译成另外两种。',
    create: '创建店铺', noStores: '你还没有店铺。', open: '打开', back: '← 我的店铺',
    st_onboarding: '设置中', st_pending_review: '审核中', st_active: '已上线', st_suspended: '已暂停', st_closed: '已关闭',
    checklist: '开始销售前', ck_logo: '店铺标志', ck_about: '店铺介绍', ck_kyc: '身份验证 (KYC)', ck_pickup: '西哈努克城取件点',
    ck_review: 'Phsar 审核通过', submit: '提交审核', submitHint: 'Phsar 会核实你的身份和取件点，通常 1–2 个工作日。',
    reviewing: 'Phsar 正在审核你的店铺，审核完成前无法修改。', activeMsg: '店铺已上线。下一步：添加商品。',
    needsFix: 'Phsar 要求修改：', tabProfile: '店铺资料', tabKyc: '身份验证', tabPickup: '取件点',
    logo: '标志', uploadLogo: '上传标志', social: '社交链接', save: '保存', saved: '已保存',
    kycIntro: '你的证件只有 Phsar 审核团队可以查看。', sellerType: '卖家类型',
    t_individual_unregistered: '个人（未税务登记）', t_individual_registered: '个人（已税务登记）', t_company: '公司',
    legalName: '证件上的全名', idDoc: '身份证或护照照片', idDocDone: '证件已上传', tin: '税号 (TIN)', mocNo: '商业部电商许可证号',
    optional: '选填', kycStatus: '状态', kyc_missing: '未提交', kyc_pending: '等待审核', kyc_verified: '已验证', kyc_rejected: '已拒绝',
    pickupIntro: 'Phsar 快递员会在这里取件。点击地图，把定位放在店门口。',
    addPickup: '添加取件点', editPickup: '编辑取件点', label: '名称', labelPh: '例如：Ochheuteal 店', address: '地址', addressPh: '街道、地标',
    contactName: '联系人', contactPhone: '联系电话', hours: '取件时间', note: '给快递员的路线说明', photo: '店铺门面照片',
    pinInside: '在配送范围内', pinOutside: '不在配送范围内（目前仅限西哈努克城）', pinNone: '点击地图设置定位',
    verified: '已核实', notVerified: '等待核实', del: '删除', edit: '编辑', cancel: '取消', default: '默认', makeDefault: '设为默认',
    d_mon: '一', d_tue: '二', d_wed: '三', d_thu: '四', d_fri: '五', d_sat: '六', d_sun: '日',
    err_validation: '请检查标出的字段。', err_slug_taken: '该店铺地址已被使用，请换个名称。', err_outside_service_area: '该位置不在西哈努克城。',
    err_store_under_review: '店铺正在审核中。', err_upload_type: '不支持此文件类型。', err_upload_too_large: '文件太大。', err_store_incomplete: '店铺信息尚未完整。'},
});

const S = {stores: [], store: null, tab: 'profile', kyc: null, points: [], editing: null, pin: null, map: null, marker: null};
const DAYS = ['mon', 'tue', 'wed', 'thu', 'fri', 'sat', 'sun'];
const statusChip = st => `<span class="chip ${({active: 'ok', pending_review: 'warn', onboarding: 'pink', suspended: 'bad'})[st] || ''}">${t('st_' + st)}</span>`;

function shell(inner) {
  return `<header class="top"><a class="brand" href="#/">${brandSVG}<span><b>Phsar</b><small>${esc(t('appName')).toUpperCase()}</small></span></a>
    <div class="top-actions">${langSwitch()}${P.token ? `<button type="button" onclick="logout()">${t('logout')}</button>` : ''}</div></header>
    <main>${inner}</main>`;
}

async function render() {
  const app = $('#app');
  if (!P.token) {
    app.innerHTML = shell(loginView(t('appName')));
    bindLogin(() => { location.hash = '#/'; render(); });
    return;
  }
  const m = location.hash.match(/^#\/store\/([0-9a-f-]+)/);
  try {
    if (m) await renderStore(m[1]); else await renderStores();
  } catch (e) { toast(errMsg(e)); }
}

// ---- Daftar toko + buat toko ----
async function renderStores() {
  S.stores = (await api('/v1/seller/stores')).stores;
  const list = S.stores.map(s => `<article class="card"><div class="row">
      <img class="logo" src="${esc(s.logo_url || '')}" alt="" onerror="this.style.visibility='hidden'">
      <div style="flex:1"><h2>${esc(s.name)}</h2><p class="muted">phsar.com/s/${esc(s.slug)}</p></div>${statusChip(s.status)}</div>
      <a class="btn" href="#/store/${s.id}">${t('open')}</a></article>`).join('');
  $('#app').innerHTML = shell(`
    <h1>${t('myStores')}</h1>
    ${S.stores.length ? `<div class="grid">${list}</div>` : `<p class="muted">${t('noStores')}</p>`}
    <section class="card" style="max-width:560px">
      <div><h2>${t('createStore')}</h2><p class="muted">${t('createSub')}</p></div>
      <label>${t('storeName')}<input id="cName" maxlength="80"></label>
      <label>${t('storeAbout')}<textarea id="cAbout" maxlength="2000"></textarea></label>
      <label>${t('writtenIn')}<select id="cLang"><option value="km">ខ្មែរ</option><option value="en">English</option><option value="zh">中文</option></select></label>
      <p class="muted">${t('aboutHint')}</p>
      <button class="primary" id="cBtn" type="button">${t('create')}</button>
    </section>`);
  $('#cLang').value = P.lang;
  $('#cBtn').onclick = async () => {
    try {
      const st = await api('/v1/seller/stores', {method: 'POST', body: {name: $('#cName').value, about: $('#cAbout').value, about_lang: $('#cLang').value}});
      location.hash = '#/store/' + st.id;
    } catch (e) { toast(errMsg(e)); }
  };
}

// ---- Halaman satu toko ----
async function renderStore(id) {
  const [st, kyc, pts] = await Promise.all([
    api(`/v1/seller/stores/${id}`),
    api(`/v1/seller/stores/${id}/kyc`).catch(() => ({kyc: null})),
    api(`/v1/seller/stores/${id}/pickup-points`),
  ]);
  S.store = st; S.kyc = kyc.kyc; S.points = pts.pickup_points;
  const c = st.checklist, locked = st.status === 'pending_review';
  const item = (ok, label) => `<div class="check ${ok ? 'done' : ''}"><i>${ok ? '✓' : ''}</i>${label}</div>`;
  const banner = st.status === 'pending_review' ? `<div class="banner warn">${t('reviewing')}</div>`
    : st.status === 'active' ? `<div class="banner ok">${t('activeMsg')}</div>`
    : st.review_note ? `<div class="banner bad"><b>${t('needsFix')}</b> ${esc(st.review_note)}</div>` : '';
  $('#app').innerHTML = shell(`
    <a href="#/">${t('back')}</a>
    <div class="spread"><div class="row"><img class="logo" src="${esc(st.logo_url || '')}" alt="" onerror="this.style.visibility='hidden'">
      <div><h1>${esc(st.name)}</h1><p class="muted">phsar.com/s/${esc(st.slug)}</p></div></div>${statusChip(st.status)}</div>
    ${banner}
    <section class="card">
      <h2>${t('checklist')}</h2>
      ${item(c.logo, t('ck_logo'))}${item(c.about, t('ck_about'))}
      ${item(c.kyc !== 'missing', `${t('ck_kyc')} · <span class="muted">${t('kyc_' + c.kyc)}</span>`)}
      ${item(c.pickup_points > 0, `${t('ck_pickup')} · <span class="muted">${c.pickup_points_verified}/${c.pickup_points} ${t('verified').toLowerCase()}</span>`)}
      ${item(st.status === 'active', t('ck_review'))}
      ${st.status === 'onboarding' ? `<div class="row"><button class="primary" id="submitBtn" type="button" ${c.can_submit ? '' : 'disabled'}>${t('submit')}</button><span class="muted">${t('submitHint')}</span></div>` : ''}
    </section>
    <div class="tabs" role="tablist">${['profile', 'kyc', 'pickup'].map(k =>
      `<button type="button" role="tab" aria-selected="${S.tab === k}" onclick="S.tab='${k}';S.editing=null;render()">${t('tab' + k[0].toUpperCase() + k.slice(1))}</button>`).join('')}</div>
    <div id="tab"></div>`);
  if ($('#submitBtn')) $('#submitBtn').onclick = async () => {
    try { await api(`/v1/seller/stores/${id}/submit`, {method: 'POST'}); render(); } catch (e) { toast(errMsg(e)); }
  };
  ({profile: tabProfile, kyc: tabKyc, pickup: tabPickup})[S.tab](id, locked);
}

function tabProfile(id, locked) {
  const st = S.store, aboutLang = Object.keys(st.about).find(l => st.about[l]) || P.lang;
  $('#tab').innerHTML = `<section class="card">
    <div class="row"><img class="logo" id="logoImg" src="${esc(st.logo_url || '')}" alt="" onerror="this.style.visibility='hidden'">
      <label class="btn" style="flex-direction:row">${t('uploadLogo')}<input type="file" id="logoFile" accept="image/jpeg,image/png,image/webp" class="hidden" ${locked ? 'disabled' : ''}></label></div>
    <label>${t('storeName')}<input id="pName" value="${esc(st.name)}" maxlength="80"></label>
    <label>${t('storeAbout')}<textarea id="pAbout" maxlength="2000">${esc(st.about[aboutLang] || '')}</textarea></label>
    <label>${t('writtenIn')}<select id="pLang"><option value="km">ខ្មែរ</option><option value="en">English</option><option value="zh">中文</option></select></label>
    <p class="muted">${t('aboutHint')}</p>
    <h2>${t('social')}</h2>
    <div class="grid">${['telegram', 'facebook', 'tiktok'].map(k => `<label>${k[0].toUpperCase() + k.slice(1)}<input id="soc_${k}" value="${esc(st.social_links?.[k] || '')}" maxlength="200"></label>`).join('')}</div>
    <div><button class="primary" id="pSave" type="button" ${locked ? 'disabled' : ''}>${t('save')}</button></div></section>`;
  $('#pLang').value = aboutLang;
  $('#logoFile').onchange = async e => {
    try {
      const up = await upload('store_logo', e.target.files[0]);
      await api(`/v1/seller/stores/${id}`, {method: 'PATCH', body: {logo_key: up.key}});
      toast(t('saved')); render();
    } catch (err) { toast(errMsg(err)); }
  };
  $('#pSave').onclick = async () => {
    const social = {};
    for (const k of ['telegram', 'facebook', 'tiktok']) { const v = $('#soc_' + k).value.trim(); if (v) social[k] = v; }
    try {
      await api(`/v1/seller/stores/${id}`, {method: 'PATCH', body: {name: $('#pName').value, about: $('#pAbout').value, about_lang: $('#pLang').value, social_links: social}});
      toast(t('saved')); render();
    } catch (e) { toast(errMsg(e)); }
  };
}

function tabKyc(id) {
  const k = S.kyc || {}, status = k.verification_status || 'missing', done = status === 'verified';
  const chip = {missing: '', pending: 'warn', verified: 'ok', rejected: 'bad'}[status];
  $('#tab').innerHTML = `<section class="card">
    <div class="spread"><p class="muted">${t('kycIntro')}</p><span class="chip ${chip}">${t('kycStatus')}: ${t('kyc_' + status)}</span></div>
    ${k.rejection_note ? `<div class="banner bad">${esc(k.rejection_note)}</div>` : ''}
    <label>${t('sellerType')}<select id="kType" ${done ? 'disabled' : ''}>${['individual_unregistered', 'individual_registered', 'company'].map(v => `<option value="${v}">${t('t_' + v)}</option>`).join('')}</select></label>
    <label>${t('legalName')}<input id="kName" value="${esc(k.legal_name || '')}" maxlength="120" ${done ? 'disabled' : ''}></label>
    <label>${t('idDoc')}<input type="file" id="kDoc" accept="image/jpeg,image/png,image/webp,application/pdf" ${done ? 'disabled' : ''}></label>
    <p class="muted" id="kDocState">${k.has_id_document ? '✓ ' + t('idDocDone') : ''}</p>
    <div class="grid"><label>${t('tin')} <span class="muted">(${t('optional')})</span><input id="kTin" value="${esc(k.tin || '')}" ${done ? 'disabled' : ''}></label>
      <label>${t('mocNo')} <span class="muted">(${t('optional')})</span><input id="kMoc" value="${esc(k.moc_license_no || '')}" ${done ? 'disabled' : ''}></label></div>
    <div><button class="primary" id="kSave" type="button" ${done ? 'disabled' : ''}>${t('save')}</button></div></section>`;
  if (k.seller_type) $('#kType').value = k.seller_type;
  let docKey = null;
  $('#kDoc').onchange = async e => {
    try { docKey = (await upload('kyc_document', e.target.files[0])).key; $('#kDocState').textContent = '✓ ' + t('idDocDone'); } catch (err) { toast(errMsg(err)); }
  };
  $('#kSave').onclick = async () => {
    try {
      if (!docKey && !k.has_id_document) { toast(t('idDoc')); return; }
      await api(`/v1/seller/stores/${id}/kyc`, {method: 'PUT', body: {seller_type: $('#kType').value, legal_name: $('#kName').value,
        id_document_key: docKey || '', tin: $('#kTin').value, moc_license_no: $('#kMoc').value}});
      toast(t('saved')); render();
    } catch (e) { toast(errMsg(e)); }
  };
}

function tabPickup(id, locked) {
  const list = S.points.map(p => `<article class="card"><div class="spread"><h2>${esc(p.label)}</h2>
      <div class="row">${p.is_default ? `<span class="chip pink">${t('default')}</span>` : ''}<span class="chip ${p.verified ? 'ok' : 'warn'}">${p.verified ? t('verified') : t('notVerified')}</span></div></div>
      <p>${esc(p.address_line)}</p><p class="muted">${esc(p.contact_name)} · ${esc(p.contact_phone)}</p>
      <p class="muted">${DAYS.filter(d => p.pickup_hours[d]).map(d => `${t('d_' + d)} ${p.pickup_hours[d].join('–')}`).join(' · ')}</p>
      ${p.rejection_note ? `<div class="banner bad">${esc(p.rejection_note)}</div>` : ''}
      <div class="row"><button type="button" onclick="editPoint('${p.id}')" ${locked ? 'disabled' : ''}>${t('edit')}</button>
      <button type="button" class="danger" onclick="deletePoint('${id}','${p.id}')" ${locked ? 'disabled' : ''}>${t('del')}</button></div></article>`).join('');
  $('#tab').innerHTML = `<p class="muted">${t('pickupIntro')}</p><div class="grid">${list}</div>
    ${locked ? '' : S.editing ? pickupForm() : `<div><button class="primary" type="button" onclick="editPoint('new')">${t('addPickup')}</button></div>`}`;
  if (S.editing) bindPickupForm(id);
}

function editPoint(pid) {
  S.editing = pid === 'new' ? {pickup_hours: {mon: ['09:00', '17:00'], tue: ['09:00', '17:00'], wed: ['09:00', '17:00'], thu: ['09:00', '17:00'], fri: ['09:00', '17:00'], sat: ['09:00', '12:00']}}
    : structuredClone(S.points.find(p => p.id === pid));
  S.pin = S.editing.lat ? [S.editing.lat, S.editing.lng] : null;
  render();
}
async function deletePoint(sid, pid) {
  try { await api(`/v1/seller/stores/${sid}/pickup-points/${pid}`, {method: 'DELETE'}); render(); } catch (e) { toast(errMsg(e)); }
}

function pickupForm() {
  const p = S.editing;
  return `<section class="card"><h2>${p.id ? t('editPickup') : t('addPickup')}</h2>
    <div id="map" class="map"></div><p class="muted" id="pinState">${t('pinNone')}</p>
    <div class="grid"><label>${t('label')}<input id="fLabel" value="${esc(p.label || '')}" placeholder="${t('labelPh')}" maxlength="60"></label>
      <label>${t('address')}<input id="fAddr" value="${esc(p.address_line || '')}" placeholder="${t('addressPh')}" maxlength="300"></label>
      <label>${t('contactName')}<input id="fCName" value="${esc(p.contact_name || '')}"></label>
      <label>${t('contactPhone')}<input id="fCPhone" value="${esc(p.contact_phone || '')}" inputmode="tel"></label></div>
    <fieldset style="border:0;padding:0;margin:0;display:flex;flex-direction:column;gap:8px"><legend style="font-size:13px;font-weight:700;margin-bottom:6px">${t('hours')}</legend>
      <div class="hours">${DAYS.map(d => { const h = p.pickup_hours?.[d]; return `<label style="flex-direction:row;align-items:center"><input type="checkbox" id="h_${d}" ${h ? 'checked' : ''}> ${t('d_' + d)}</label>
        <input type="time" id="hf_${d}" value="${h ? h[0] : '09:00'}" aria-label="${t('d_' + d)} from"><input type="time" id="ht_${d}" value="${h ? h[1] : '17:00'}" aria-label="${t('d_' + d)} to">`; }).join('')}</div></fieldset>
    <label>${t('note')} <textarea id="fNote" maxlength="500">${esc(p.note || '')}</textarea></label>
    <label>${t('photo')} <span class="muted">(${t('optional')})</span><input type="file" id="fPhoto" accept="image/jpeg,image/png,image/webp"></label>
    <label style="flex-direction:row;align-items:center;gap:8px"><input type="checkbox" id="fDefault" style="width:20px;min-height:20px;accent-color:var(--pink)" ${p.is_default ? 'checked' : ''}> ${t('makeDefault')}</label>
    <div class="row"><button class="primary" id="fSave" type="button">${t('save')}</button><button type="button" onclick="S.editing=null;render()">${t('cancel')}</button></div></section>`;
}

async function checkPin() {
  if (!S.pin) return;
  try {
    const r = await api(`/v1/app/service-areas/check?lat=${S.pin[0]}&lng=${S.pin[1]}`);
    $('#pinState').innerHTML = `<span class="chip ${r.served ? 'ok' : 'bad'}">${r.served ? t('pinInside') : t('pinOutside')}</span> <span class="muted">${S.pin[0].toFixed(5)}, ${S.pin[1].toFixed(5)}</span>`;
  } catch (e) { toast(errMsg(e)); }
}

function bindPickupForm(id) {
  const center = S.pin || [10.6274, 103.5223]; // pusat Sihanoukville
  S.map = L.map('map').setView(center, S.pin ? 17 : 13);
  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {maxZoom: 19, attribution: '© OpenStreetMap'}).addTo(S.map);
  S.marker = S.pin ? L.marker(S.pin).addTo(S.map) : null;
  checkPin();
  S.map.on('click', e => {
    S.pin = [Number(e.latlng.lat.toFixed(6)), Number(e.latlng.lng.toFixed(6))];
    if (S.marker) S.marker.setLatLng(S.pin); else S.marker = L.marker(S.pin).addTo(S.map);
    checkPin();
  });
  let photoKey = '';
  $('#fPhoto').onchange = async e => { try { photoKey = (await upload('pickup_photo', e.target.files[0])).key; } catch (err) { toast(errMsg(err)); } };
  $('#fSave').onclick = async () => {
    if (!S.pin) { toast(t('pinNone')); return; }
    const hours = {};
    for (const d of DAYS) if ($('#h_' + d).checked) hours[d] = [$('#hf_' + d).value, $('#ht_' + d).value];
    const body = {label: $('#fLabel').value, address_line: $('#fAddr').value, lat: S.pin[0], lng: S.pin[1], contact_name: $('#fCName').value,
      contact_phone: $('#fCPhone').value, pickup_hours: hours, note: $('#fNote').value, photo_key: photoKey, is_default: $('#fDefault').checked};
    try {
      const pid = S.editing.id;
      await api(`/v1/seller/stores/${id}/pickup-points${pid ? '/' + pid : ''}`, {method: pid ? 'PUT' : 'POST', body});
      S.editing = null; toast(t('saved')); render();
    } catch (e) { toast(errMsg(e)); }
  };
}

window.addEventListener('hashchange', () => { S.editing = null; render(); });
(async () => { document.documentElement.lang = P.lang; await refreshSession(); render(); })();
