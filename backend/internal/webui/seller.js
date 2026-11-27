'use strict';
P.authBase = '/seller/auth';
addStrings({
  en: {oneStoreRule: 'One account can own one store. To sell under another name, use a different account.', oneStoreNote: 'One account can own only one store.', err_too_many_stores: 'One account can own only one store.', appName: 'Seller Center', logout: 'Sign out', myStores: 'My stores', createStore: 'Create your store', createSub: 'Fill in the basics. You can edit everything later.',
    storeName: 'Store name', storeAbout: 'About your store', writtenIn: 'Written in', aboutHint: 'Write in one language. We translate it to the other two automatically.',
    create: 'Create store', noStores: 'You don’t have a store yet.', open: 'Open', back: '← My stores',
    st_onboarding: 'Setting up', st_pending_review: 'Under review', st_active: 'Active', st_suspended: 'Suspended', st_closed: 'Closed',
    checklist: 'Before you can sell', ck_logo: 'Store logo', ck_about: 'Store description', ck_kyc: 'Identity check (KYC)', ck_pickup: 'Pickup point in Sihanoukville',
    ck_review: 'Approved by Phsar', submit: 'Send for review', submitHint: 'Phsar checks your identity and pickup point, usually within 1–2 working days.',
    reviewing: 'Phsar is reviewing your store. You can’t edit it until the review is done.', activeMsg: 'Your store is active. Add products in the Products tab.',
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
    searchPlace: 'Search your shop or a nearby place', searchHint: 'Pick a result, then drag the pin to your exact door.', myLocation: 'Use my location', noGeo: 'Could not get your location.',
    verified: 'Verified', notVerified: 'Waiting for check', del: 'Delete', edit: 'Edit', cancel: 'Cancel', default: 'Default', makeDefault: 'Use as default',
    d_mon: 'Mon', d_tue: 'Tue', d_wed: 'Wed', d_thu: 'Thu', d_fri: 'Fri', d_sat: 'Sat', d_sun: 'Sun',
    err_validation: 'Please check the highlighted field.', err_slug_taken: 'That store address is taken. Try another name.', err_outside_service_area: 'That location is outside Sihanoukville.',
    err_store_under_review: 'Your store is under review.', err_upload_type: 'This file type is not allowed.', err_upload_too_large: 'The file is too large.', err_store_incomplete: 'Your store is not complete yet.'},
  km: {oneStoreRule: 'គណនីមួយអាចមានហាងតែមួយ។ ដើម្បីលក់ក្រោមឈ្មោះផ្សេង សូមប្រើគណនីផ្សេង។', oneStoreNote: 'គណនីមួយអាចមានហាងតែមួយប៉ុណ្ណោះ។', err_too_many_stores: 'គណនីមួយអាចមានហាងតែមួយប៉ុណ្ណោះ។', appName: 'មជ្ឈមណ្ឌលអ្នកលក់', logout: 'ចាកចេញ', myStores: 'ហាងរបស់ខ្ញុំ', createStore: 'បង្កើតហាងរបស់អ្នក', createSub: 'បំពេញព័ត៌មានមូលដ្ឋាន។ អាចកែបានពេលក្រោយ។',
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
    searchPlace: 'ស្វែងរកហាង ឬទីកន្លែងក្បែរនោះ', searchHint: 'ជ្រើសលទ្ធផល រួចអូសម្ជុលទៅត្រង់ទ្វារហាង។', myLocation: 'ប្រើទីតាំងខ្ញុំ', noGeo: 'មិនអាចយកទីតាំងរបស់អ្នកបានទេ។',
    verified: 'បានផ្ទៀងផ្ទាត់', notVerified: 'រង់ចាំពិនិត្យ', del: 'លុប', edit: 'កែ', cancel: 'បោះបង់', default: 'លំនាំដើម', makeDefault: 'ប្រើជាលំនាំដើម',
    d_mon: 'ច', d_tue: 'អ', d_wed: 'ព', d_thu: 'ព្រ', d_fri: 'សុ', d_sat: 'ស', d_sun: 'អា',
    err_validation: 'សូមពិនិត្យវាលដែលបានសម្គាល់។', err_slug_taken: 'អាសយដ្ឋានហាងនេះមានគេប្រើហើយ។', err_outside_service_area: 'ទីតាំងនៅក្រៅក្រុងព្រះសីហនុ។',
    err_store_under_review: 'ហាងកំពុងត្រួតពិនិត្យ។', err_upload_type: 'ប្រភេទឯកសារមិនអនុញ្ញាត។', err_upload_too_large: 'ឯកសារធំពេក។', err_store_incomplete: 'ហាងមិនទាន់គ្រប់ព័ត៌មាន។'},
  zh: {oneStoreRule: '一个账号只能拥有一家店铺。如需用其他名称销售，请使用另一个账号。', oneStoreNote: '一个账号只能拥有一家店铺。', err_too_many_stores: '一个账号只能拥有一家店铺。', appName: '卖家中心', logout: '退出', myStores: '我的店铺', createStore: '创建你的店铺', createSub: '先填写基本信息，之后都可以修改。',
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
    searchPlace: '搜索你的店铺或附近地点', searchHint: '选择结果后，把定位拖到店门口。', myLocation: '使用我的位置', noGeo: '无法获取你的位置。',
    verified: '已核实', notVerified: '等待核实', del: '删除', edit: '编辑', cancel: '取消', default: '默认', makeDefault: '设为默认',
    d_mon: '一', d_tue: '二', d_wed: '三', d_thu: '四', d_fri: '五', d_sat: '六', d_sun: '日',
    err_validation: '请检查标出的字段。', err_slug_taken: '该店铺地址已被使用，请换个名称。', err_outside_service_area: '该位置不在西哈努克城。',
    err_store_under_review: '店铺正在审核中。', err_upload_type: '不支持此文件类型。', err_upload_too_large: '文件太大。', err_store_incomplete: '店铺信息尚未完整。'},
});

const S = {stores: [], store: null, tab: 'profile', kyc: null, points: [], editing: null, pin: null, map: null, marker: null};
const DAYS = ['mon', 'tue', 'wed', 'thu', 'fri', 'sat', 'sun'];
// Daftar / masuk Seller Center: login = daftar (nomor baru otomatis dibuatkan akun) + peta langkah menjadi seller.
addStrings({
  en: {jnLogin: 'Register or sign in', jnLoginSub: 'New here? Choose any sign-in below — your account is created automatically. Already a seller? Sign in the same way.',
    jnTitle: 'Become a seller in 6 steps', jn1: 'Sign up — takes a minute', jn2: 'Create your store', jn3: 'Logo & description', jn4: 'Verify your identity (KYC)', jn5: 'Add your pickup point', jn6: 'Send for review — then sell',
    jnFree: 'Free to open. Commission only when you sell.', jnMore: 'Learn more about selling on Phsar', jnStep: 'Step {a} of {b}', jnWelcome: 'Welcome! Let’s set up your store.'},
  km: {jnLogin: 'ចុះឈ្មោះ ឬចូលគណនី', jnLoginSub: 'ថ្មីនៅទីនេះ? ជ្រើសរើសវិធីចូលណាមួយខាងក្រោម — គណនីនឹងត្រូវបង្កើតដោយស្វ័យប្រវត្តិ។ ជាអ្នកលក់រួចហើយ? ចូលតាមរបៀបដដែល។',
    jnTitle: 'ក្លាយជាអ្នកលក់ក្នុង 6 ជំហាន', jn1: 'ចុះឈ្មោះ — ប្រើពេលមួយនាទី', jn2: 'បង្កើតហាងរបស់អ្នក', jn3: 'ឡូហ្គោ និងការពិពណ៌នា', jn4: 'ផ្ទៀងផ្ទាត់អត្តសញ្ញាណ (KYC)', jn5: 'បន្ថែមទីតាំងទទួលកញ្ចប់', jn6: 'ផ្ញើឲ្យពិនិត្យ — រួចចាប់ផ្តើមលក់',
    jnFree: 'បើកហាងឥតគិតថ្លៃ។ បង់កម្រៃជើងសារតែពេលលក់ដាច់។', jnMore: 'ស្វែងយល់បន្ថែមអំពីការលក់លើ Phsar', jnStep: 'ជំហាន {a} នៃ {b}', jnWelcome: 'សូមស្វាគមន៍! តោះរៀបចំហាងរបស់អ្នក។'},
  zh: {jnLogin: '注册或登录', jnLoginSub: '第一次来？选择下面任一登录方式，账户会自动创建。已经是卖家？用同样的方式登录。',
    jnTitle: '6 步成为卖家', jn1: '注册——一分钟搞定', jn2: '创建店铺', jn3: '标志和介绍', jn4: '身份验证 (KYC)', jn5: '添加取件点', jn6: '提交审核，然后开始销售',
    jnFree: '免费开店，成交才收佣金。', jnMore: '了解更多在 Phsar 开店的信息', jnStep: '第 {a} 步，共 {b} 步', jnWelcome: '欢迎！我们来设置你的店铺。'},
});
// journey: daftar 5 langkah; langkah < active = selesai, = active = sekarang.
function journey(active) {
  return `<ol class="journey">${[1, 2, 3, 4, 5, 6].map(n => `<li class="${n < active ? 'done' : n === active ? 'now' : ''}"${n === active ? ' aria-current="step"' : ''}>
    <span class="n" aria-hidden="true">${n < active ? '✓' : n}</span><span>${t('jn' + n)}</span></li>`).join('')}</ol>`;
}
function sellerLoginView() {
  return `<div class="jn-login">${loginView(t('jnLogin')).replace(`<p class="muted">${t('loginSub')}</p>`, `<p class="muted">${t('jnLoginSub')}</p>`)}
    <aside class="card jn-side"><h2>${t('jnTitle')}</h2>${journey(1)}<p class="muted">${t('jnFree')}</p>
      <a href="/${P.lang}/sell">${t('jnMore')} →</a></aside></div>`;
}

const statusChip = st => `<span class="chip ${({active: 'ok', pending_review: 'warn', onboarding: 'pink', suspended: 'bad'})[st] || ''}">${t('st_' + st)}</span>`;

function shell(inner) {
  return `<header class="top"><a class="brand" href="#/">${brandSVG}<span><b>Phsar</b><small>${esc(t('appName')).toUpperCase()}</small></span></a>
    <div class="top-actions">${P.token ? bellButton() : ''}${langSwitch()}${P.token ? `<button type="button" onclick="logout()">${t('logout')}</button>` : ''}</div></header>
    <main>${inner}</main>`;
}

async function render() {
  const app = $('#app');
  if (!P.token) {
    app.innerHTML = shell(sellerLoginView());
    // Setelah login, kembali ke halaman yang dibuka (mis. link undangan staf).
    bindLogin(() => { if (!/^#\/invite\//.test(location.hash)) location.hash = '#/'; render(); });
    return;
  }
  const inv = location.hash.match(/^#\/invite\/([0-9a-f]+)/);
  const pm = location.hash.match(/^#\/store\/([0-9a-f-]+)\/product\/(new|[0-9a-f-]+)/);
  const m = location.hash.match(/^#\/store\/([0-9a-f-]+)/);
  const sub = location.hash.match(/^#\/store\/[0-9a-f-]+\/(home|more|orders|products|payouts|returns|help|chat|ads|flash|profile|kyc|pickup|team|hours|faq)(?:\/([0-9a-f-]{36}))?/);
  if (sub && sub[1] === 'help') S.helpTicket = sub[2] || null; // support: buka tiket dari notifikasi
  if (sub && sub[1] === 'chat') S.chatConv = sub[2] || null; // chat: buka percakapan dari notifikasi
  if (sub) { S.tab = sub[1]; S.storeId = m[1]; try { sessionStorage.setItem('phsar-tab-' + m[1], sub[1]); } catch {} history.replaceState(null, '', '#/store/' + m[1]); }
  try {
    if (inv) await renderInvite(inv[1]);
    else if (pm) { if ($('#fName') && PF.E) readForm(); await renderProductEditor(pm[1], pm[2]); }
    else if (m) await renderStore(m[1]); else await renderStores();
  } catch (e) { toast(errMsg(e), 'error'); }
  loadBell();
}

// ---- Daftar toko + buat toko ----
async function renderStores() {
  S.stores = (await api('/v1/seller/stores')).stores;
  const list = S.stores.map(s => `<article class="card"><div class="row">
      <img class="logo" src="${esc(s.logo_url || '')}" alt="" onerror="this.style.visibility='hidden'">
      <div style="flex:1"><h2>${esc(s.name)}</h2><p class="muted">phsar.com/s/${esc(s.slug)}</p></div>${statusChip(s.status)}</div>
      <a class="btn" href="#/store/${s.id}">${t('open')}</a></article>`).join('');
  const owns = S.stores.some(s => s.role === 'owner');
  // 1 akun = 1 toko: kalau hanya punya satu toko (tanpa jadi staf di toko lain), langsung buka tokonya.
  if (S.stores.length === 1 && owns && !S.storesShown) { S.storesShown = true; location.hash = '#/store/' + S.stores[0].id + '/home'; return; }
  $('#app').innerHTML = shell(`
    <h1>${t('myStores')}</h1>
    ${S.stores.length ? `<div class="grid">${list}</div>` : `<p class="muted">${t('jnWelcome')}</p>${owns ? '' : journey(2)}`}
    ${owns ? `<p class="muted">${t('oneStoreRule')}</p>` : `<section class="card" style="max-width:560px">
      <div><span class="chip pink">${t('jnStep', {a: 2, b: 6})}</span><h2>${t('createStore')}</h2><p class="muted">${t('createSub')}</p></div>
      <label>${t('storeName')}<input id="cName" maxlength="80"></label>
      <label>${t('storeAbout')}<textarea id="cAbout" maxlength="2000"></textarea></label>
      <label>${t('writtenIn')}<select id="cLang"><option value="km">ខ្មែរ</option><option value="en">English</option><option value="zh">中文</option></select></label>
      <p class="muted">${t('aboutHint')}</p>
      <p class="muted">${t('oneStoreNote')}</p>
      <button class="primary" id="cBtn" type="button" data-errors-anchor>${t('create')}</button>
    </section>`}`);
  if (owns) return;
  $('#cLang').value = P.lang;
  $('#cBtn').onclick = async () => {
    if (!checkRequired([['#cName', t('minChars', {n: 2}), el => el.value.trim().length >= 2]])) return;
    try {
      const st = await api('/v1/seller/stores', {method: 'POST', body: {name: $('#cName').value, about: $('#cAbout').value, about_lang: $('#cLang').value}});
      location.hash = '#/store/' + st.id;
    } catch (e) { showError(e, {name: '#cName', about: '#cAbout', about_lang: '#cLang'}); }
  };
}

// ---- Halaman satu toko (Seller Center v2: Hari ini · Pesanan · Produk · Uang · Lainnya) ----
const storeTabs = () => ({home: tabHome, more: tabMore, help: tabHelp, flash: tabFlash, ads: tabAds, returns: tabReturns, orders: tabOrders,
  payouts: tabPayouts, products: tabProducts, profile: tabProfile, kyc: tabKyc, pickup: tabPickup, team: tabTeam, hours: tabHours,
  faq: abTab, chat: tabChat}); // askbot: Jawaban cepat & FAQ · chat:
async function renderStore(id) {
  const [st, kyc, pts, today] = await Promise.all([
    api(`/v1/seller/stores/${id}`),
    api(`/v1/seller/stores/${id}/kyc`).catch(() => ({kyc: null})),
    api(`/v1/seller/stores/${id}/pickup-points`),
    api(`/v1/seller/stores/${id}/today`).catch(() => null),
  ]);
  const owner = st.role === 'owner';
  if (S.storeId !== id) {
    S.storeId = id;
    let saved = null;
    try { saved = sessionStorage.getItem('phsar-tab-' + id); } catch {}
    S.tab = saved || 'home';
  }
  const allowed = owner ? Object.keys(storeTabs()) : ['home', 'more', 'orders', 'returns', 'products', 'help', 'faq', 'chat']; // staf: pesanan, retur, produk, FAQ (askbot), chat
  if (!allowed.includes(S.tab)) S.tab = 'home';
  S.store = st; S.kyc = kyc.kyc; S.points = pts.pickup_points; V2.today = today;
  const locked = st.status === 'pending_review';
  const banner = locked && ['profile', 'kyc', 'pickup'].includes(S.tab) ? `<div class="banner warn">${t('reviewing')}</div>`
    : st.status === 'suspended' && S.tab !== 'home' ? `<div class="banner bad">${t('su_suspended')}</div>`
    : st.review_note && st.status === 'onboarding' && ['profile', 'kyc', 'pickup'].includes(S.tab) ? `<div class="banner bad"><b>${t('needsFix')}</b> ${esc(st.review_note)}</div>` : '';
  const setup = inSetup(st, S.tab); // langkah pendaftaran: judul "Langkah X dari 6" + tombol ke langkah berikutnya
  $('#app').innerHTML = storeShell(st, `<div id="storeAlerts"></div>${setup ? setupHead(st, S.tab) : v2Header(S.tab)}${banner}<div id="tab"></div>${setup ? setupNext(st, S.tab) : ''}`);
  if (typeof sellerStrikes === 'function') sellerStrikes(id); // reports:
  await storeTabs()[S.tab](id, locked);
  if (setup) { // langkah ini baru saja selesai → arahkan ke tombol lanjut
    const key = id + ':' + S.tab, wasDone = S.setupDone?.[key], nowDone = !!setupNext(st, S.tab);
    S.setupDone = Object.assign(S.setupDone || {}, {[key]: nowDone});
    if (nowDone && wasDone === false) $('#swNext')?.scrollIntoView({behavior: 'smooth', block: 'center'});
  }
  if (typeof tourSeller === 'function') tourSeller(id); // tour: tur cara pakai (seller-tour.js)
}

function tabProfile(id, locked) {
  const st = S.store, aboutLang = Object.keys(st.about).find(l => st.about[l]) || P.lang;
  // sellerpages: pratinjau toko di atas, lalu nama & deskripsi, lalu sosial (opsional); satu tombol Simpan.
  $('#tab').innerHTML = `<h2 class="sp-sub">${t('spPfLook')}</h2><section class="card sp-look">
    <div class="store-banner">${st.banner_url ? `<img src="${esc(st.banner_url)}" alt="">` : `<span class="muted">${t('bannerHint')}</span>`}</div>
    <div class="sp-who"><img class="logo" id="logoImg" src="${esc(st.logo_url || '')}" alt="" onerror="this.style.visibility='hidden'"><b>${esc(st.name)}</b></div>
    <div class="sp-pbtns"><label class="btn">🖼 ${st.banner_url ? t('changeBanner') : t('uploadBanner')}<input type="file" id="bannerFile" accept="image/jpeg,image/png,image/webp" class="hidden" ${locked ? 'disabled' : ''}></label>
      <label class="btn">🏪 ${st.logo_url ? t('spPfLogo') : t('uploadLogo')}<input type="file" id="logoFile" accept="image/jpeg,image/png,image/webp" class="hidden" ${locked ? 'disabled' : ''}></label></div>
    <p class="muted small">${t('bannerSize')}</p></section>
    <section class="card"><h2>${t('spPfName')}</h2>
    <label>${t('storeName')}<input id="pName" value="${esc(st.name)}" maxlength="80"></label>
    <label>${t('storeAbout')}<textarea id="pAbout" maxlength="2000" rows="4">${esc(st.about[aboutLang] || '')}</textarea></label>
    <label>${t('writtenIn')}<select id="pLang"><option value="km">ខ្មែរ</option><option value="en">English</option><option value="zh">中文</option></select></label>
    <p class="muted" style="margin:0">${t('aboutHint')}</p></section>
    <section class="card"><h2>${t('social')}</h2><p class="muted" style="margin:0">${t('spPfSocialW')}</p>
    ${['telegram', 'facebook', 'tiktok'].map(k => `<label>${k[0].toUpperCase() + k.slice(1)}<input id="soc_${k}" value="${esc(st.social_links?.[k] || '')}" maxlength="200"></label>`).join('')}</section>
    <div data-errors-anchor><button class="primary sp-block" id="pSave" type="button" ${locked ? 'disabled' : ''}>${t('save')}</button></div>
    ${catalogSettingsHTML(st)}`;
  $('#pLang').value = aboutLang;
  bindCatalogSettings(id); // catalog: kode SKU & ambang stok (bisa diubah juga saat toko menunggu verifikasi)
  $('#bannerFile').onchange = async e => {
    const f = e.target.files[0];
    if (!f) return;
    if (f.size > 8 * 1024 * 1024) { toast(t('err_upload_too_large'), 'error'); return; }
    try {
      const up = await upload('store_banner', f);
      await api(`/v1/seller/stores/${id}`, {method: 'PATCH', body: {banner_key: up.key}});
      toast(t('saved')); render();
    } catch (err) { showError(err, {banner_key: '#bannerFile', media: '#bannerFile'}); }
  };
  $('#logoFile').onchange = async e => {
    try {
      const up = await upload('store_logo', e.target.files[0]);
      await api(`/v1/seller/stores/${id}`, {method: 'PATCH', body: {logo_key: up.key}});
      toast(t('saved')); render();
    } catch (err) { toast(errMsg(err), 'error'); }
  };
  $('#pSave').onclick = async () => {
    const social = {};
    for (const k of ['telegram', 'facebook', 'tiktok']) { const v = $('#soc_' + k).value.trim(); if (v) social[k] = v; }
    if (!checkRequired([['#pName', t('minChars', {n: 2}), el => el.value.trim().length >= 2]])) return;
    try {
      await api(`/v1/seller/stores/${id}`, {method: 'PATCH', body: {name: $('#pName').value, about: $('#pAbout').value, about_lang: $('#pLang').value, social_links: social}});
      toast(t('saved')); render();
    } catch (e) { showError(e, f => f.startsWith('social') ? '#soc_telegram' : ({name: '#pName', about: '#pAbout', about_lang: '#pLang', logo_key: '#logoFile'})[f]); }
  };
}

function tabKyc(id) {
  const k = S.kyc || {}, status = k.verification_status || 'missing', done = status === 'verified';
  const chip = {missing: '', pending: 'warn', verified: 'ok', rejected: 'bad'}[status];
  // sellerpages: status besar di atas + 3 langkah bernomor; data pajak (opsional) dilipat.
  const sIcon = {missing: '🪪', pending: '⏳', verified: '✅', rejected: '⚠️'}[status];
  const dis = done ? 'disabled' : '';
  $('#tab').innerHTML = `<div class="sp-status ${chip}" role="status"><span class="si" aria-hidden="true">${sIcon}</span><div>
      <b>${t('kycStatus')}: ${t('spKy_' + status)}</b><p>${t('spKy_' + status + 'W')}</p></div></div>
    ${k.rejection_note ? `<div class="banner bad">${esc(k.rejection_note)}</div>` : ''}
    <section class="card">
    <div class="sp-num"><span class="n" aria-hidden="true">1</span><label>${t('sellerType')}<select id="kType" ${dis}>${['individual_unregistered', 'individual_registered', 'company'].map(v => `<option value="${v}">${t('t_' + v)}</option>`).join('')}</select></label></div>
    <div class="sp-num"><span class="n" aria-hidden="true">2</span><label>${t('legalName')}<input id="kName" value="${esc(k.legal_name || '')}" maxlength="120" ${dis}></label></div>
    <div class="sp-num"><span class="n" aria-hidden="true">3</span><div><label class="sp-file">📷 ${t('idDoc')}<small>${t('spKyFileW')}</small><input type="file" id="kDoc" accept="image/jpeg,image/png,image/webp,application/pdf" ${dis}></label>
      <p class="muted" id="kDocState" style="margin:6px 0 0;color:var(--ok);font-weight:700">${k.has_id_document ? '✓ ' + t('idDocDone') : ''}</p></div></div>
    <details class="sp-more" id="spKyMore" ${k.tin || k.moc_license_no ? 'open' : ''}><summary>${t('spKyMore')}</summary><div>
      <label>${t('tin')} <span class="muted">(${t('optional')})</span><input id="kTin" value="${esc(k.tin || '')}" ${dis}></label>
      <label>${t('mocNo')} <span class="muted">(${t('optional')})</span><input id="kMoc" value="${esc(k.moc_license_no || '')}" ${dis}></label></div></details>
    <div data-errors-anchor><button class="primary sp-block" id="kSave" type="button" ${dis}>${t('save')}</button></div>
    <p class="sp-lock">🔒 ${t('kycIntro')}</p></section>`;
  if (k.seller_type) $('#kType').value = k.seller_type;
  let docKey = null;
  $('#kDoc').onchange = async e => {
    try { docKey = (await upload('kyc_document', e.target.files[0])).key; $('#kDocState').textContent = '✓ ' + t('idDocDone'); } catch (err) { toast(errMsg(err), 'error'); }
  };
  $('#kSave').onclick = async () => {
    if (!checkRequired([['#kName'], ['#kDoc', t('required'), () => !!(docKey || k.has_id_document)]])) return;
    try {
      await api(`/v1/seller/stores/${id}/kyc`, {method: 'PUT', body: {seller_type: $('#kType').value, legal_name: $('#kName').value,
        id_document_key: docKey || '', tin: $('#kTin').value, moc_license_no: $('#kMoc').value}});
      toast(t('saved')); render();
    } catch (e) { if (/^(tin|moc_license_no)$/.test(e.field || '')) $('#spKyMore').open = true; showError(e, {seller_type: '#kType', legal_name: '#kName', id_document_key: '#kDoc', media: '#kDoc', tin: '#kTin', moc_license_no: '#kMoc'}); }
  };
}

function tabPickup(id, locked) {
  // sellerpages: kartu titik jemput yang mudah dibaca + kosong-state dengan ajakan menambah.
  const list = S.points.map(p => `<article class="card sp-pt"><div class="spread"><h2>${esc(p.label)}</h2>
      <div class="sp-chips">${p.is_default ? `<span class="chip pink">${t('default')}</span>` : ''}<span class="chip ${p.verified ? 'ok' : 'warn'}">${p.verified ? t('verified') : t('notVerified')}</span></div></div>
      <p class="sp-ln"><span aria-hidden="true">📍</span><span>${esc(p.address_line)}</span></p><p class="sp-ln"><span aria-hidden="true">👤</span><span>${esc(p.contact_name)} · ${esc(p.contact_phone)}</span></p>
      <p class="sp-ln sp-hh"><span aria-hidden="true">🕘</span><span>${DAYS.filter(d => p.pickup_hours[d]).map(d => `${t('d_' + d)} ${p.pickup_hours[d].join('–')}`).join(' · ')}</span></p>
      ${p.rejection_note ? `<div class="banner bad">${esc(p.rejection_note)}</div>` : ''}
      <div class="sp-row2"><button type="button" onclick="editPoint('${p.id}')" ${locked ? 'disabled' : ''}>✏️ ${t('edit')}</button>
      <button type="button" class="danger" onclick="deletePoint('${id}','${p.id}')" ${locked ? 'disabled' : ''}>${t('del')}</button></div></article>`).join('');
  const add = `<button class="primary sp-block" type="button" onclick="editPoint('new')">+ ${t('addPickup')}</button>`;
  $('#tab').innerHTML = S.points.length || S.editing ? `<p class="ohelp">🛵 ${t('pickupIntro')}</p>${list}
    ${locked ? '' : S.editing ? pickupForm() : add}`
    : spEmpty('📍', t('spPkNoneT'), t('pickupIntro'), locked ? '' : add);
  if (S.editing) bindPickupForm(id);
}

function editPoint(pid) {
  S.editing = pid === 'new' ? {pickup_hours: {mon: ['09:00', '17:00'], tue: ['09:00', '17:00'], wed: ['09:00', '17:00'], thu: ['09:00', '17:00'], fri: ['09:00', '17:00'], sat: ['09:00', '12:00']}}
    : structuredClone(S.points.find(p => p.id === pid));
  S.pin = S.editing.lat ? [S.editing.lat, S.editing.lng] : null;
  render();
}
async function deletePoint(sid, pid) {
  if (!confirm(t('spPkDelQ'))) return; // sellerpages: konfirmasi sebelum hapus
  try { await api(`/v1/seller/stores/${sid}/pickup-points/${pid}`, {method: 'DELETE'}); render(); } catch (e) { toast(errMsg(e), 'error'); }
}

function pickupForm() {
  const p = S.editing;
  const h3 = (n, k) => `<h3><span class="n" aria-hidden="true">${n}</span>${t(k)}</h3>`;
  // sellerpages: formulir dibagi 4 bagian bernomor; jam jemput tidak lagi melebar di HP.
  return `<section class="card sp-form"><h2>${p.id ? t('editPickup') : t('addPickup')}</h2>
    ${h3(1, 'spPkS1')}
    ${Maps.google() ? `<div class="field"><span class="flabel">${t('searchPlace')}</span><div id="placeSearch" class="place-wrap"></div></div><p class="muted">${t('searchHint')}</p>` : ''}
    <div id="map" class="map"></div>
    <div class="spread"><p class="muted" id="pinState">${t('pinNone')}</p><button type="button" id="myLoc">📍 ${t('myLocation')}</button></div>
    ${h3(2, 'spPkS2')}
    <div class="grid"><label>${t('label')}<input id="fLabel" value="${esc(p.label || '')}" placeholder="${t('labelPh')}" maxlength="60"></label>
      <label>${t('address')}<input id="fAddr" value="${esc(p.address_line || '')}" placeholder="${t('addressPh')}" maxlength="300"></label>
      <label>${t('contactName')}<input id="fCName" value="${esc(p.contact_name || '')}"></label>
      <label>${t('contactPhone')}<input id="fCPhone" value="${esc(p.contact_phone || '')}" inputmode="tel"></label></div>
    ${h3(3, 'spPkS3')}
    <fieldset style="border:0;padding:0;margin:0;display:flex;flex-direction:column;gap:8px"><legend class="sr-only" style="position:absolute;left:-9999px">${t('hours')}</legend>
      <div class="hours">${DAYS.map(d => { const h = p.pickup_hours?.[d]; return `<label style="flex-direction:row;align-items:center"><input type="checkbox" id="h_${d}" ${h ? 'checked' : ''}> ${t('d_' + d)}</label>
        <input type="time" id="hf_${d}" value="${h ? h[0] : '09:00'}" aria-label="${t('d_' + d)} from"><input type="time" id="ht_${d}" value="${h ? h[1] : '17:00'}" aria-label="${t('d_' + d)} to">`; }).join('')}</div></fieldset>
    ${h3(4, 'spPkS4')}
    <label>${t('note')} <textarea id="fNote" maxlength="500">${esc(p.note || '')}</textarea></label>
    <label>${t('photo')} <span class="muted">(${t('optional')})</span><input type="file" id="fPhoto" accept="image/jpeg,image/png,image/webp"></label>
    <label style="flex-direction:row;align-items:center;gap:8px;min-height:44px"><input type="checkbox" id="fDefault" style="width:20px;min-height:20px;accent-color:var(--pink)" ${p.is_default ? 'checked' : ''}> ${t('makeDefault')}</label>
    <div class="sp-acts" data-errors-anchor><button class="primary sp-block" id="fSave" type="button">${t('save')}</button><button type="button" class="textbtn" onclick="S.editing=null;render()">${t('cancel')}</button></div></section>`;
}

async function checkPin() {
  if (!S.pin) return;
  try {
    const r = await api(`/v1/app/service-areas/check?lat=${S.pin[0]}&lng=${S.pin[1]}`);
    $('#pinState').innerHTML = `<span class="chip ${r.served ? 'ok' : 'bad'}">${r.served ? t('pinInside') : t('pinOutside')}</span> <span class="muted">${S.pin[0].toFixed(5)}, ${S.pin[1].toFixed(5)}</span>`;
  } catch (e) { toast(errMsg(e), 'error'); }
}

async function bindPickupForm(id) {
  S.placeId = S.editing.google_place_id || '';
  let autoAddr = null; // alamat terakhir yang diisi otomatis; kalau seller mengetik sendiri, tidak ditimpa
  const setAddr = (addr, force) => {
    const el = $('#fAddr');
    if (!el || !addr) return;
    if (force || !el.value.trim() || el.value === autoAddr) { el.value = addr.slice(0, 300); autoAddr = el.value; }
  };
  const pickerP = Maps.picker({
    mapEl: $('#map'), searchEl: $('#placeSearch'),
    pin: S.pin ? {lat: S.pin[0], lng: S.pin[1]} : null,
    onPin: p => { S.pin = [p.lat, p.lng]; checkPin(); },
    onPlace: pl => {
      S.placeId = pl.placeId;
      setAddr(pl.address, true);
      if (!$('#fLabel').value.trim() && pl.name) $('#fLabel').value = pl.name.slice(0, 60);
    },
    onAddress: (addr, placeId) => { setAddr(addr); if (placeId && !S.placeId) S.placeId = placeId; },
  });
  checkPin();
  $('#myLoc').onclick = async () => {
    const picker = await pickerP;
    if (!navigator.geolocation) { toast(t('noGeo'), 'warning'); return; }
    navigator.geolocation.getCurrentPosition(
      pos => { S.placeId = ''; picker.setPin({lat: pos.coords.latitude, lng: pos.coords.longitude}); },
      () => toast(t('noGeo'), 'warning'), {enableHighAccuracy: true, timeout: 10000});
  };
  let photoKey = '';
  $('#fPhoto').onchange = async e => { try { photoKey = (await upload('pickup_photo', e.target.files[0])).key; } catch (err) { toast(errMsg(err), 'error'); } };
  $('#fSave').onclick = async () => {
    if (!checkRequired([
      ['#map', t('pinNone'), () => !!S.pin],
      ['#fLabel'],
      ['#fAddr', t('minChars', {n: 5}), el => el.value.trim().length >= 5],
      ['#fCName'],
      ['#fCPhone', t('minChars', {n: 8}), el => el.value.replace(/\D/g, '').length >= 8],
      ['.hours', t('required'), () => DAYS.some(d => $('#h_' + d).checked)],
    ])) return;
    const hours = {};
    for (const d of DAYS) if ($('#h_' + d).checked) hours[d] = [$('#hf_' + d).value, $('#ht_' + d).value];
    const body = {label: $('#fLabel').value, address_line: $('#fAddr').value, lat: S.pin[0], lng: S.pin[1], contact_name: $('#fCName').value,
      contact_phone: $('#fCPhone').value, pickup_hours: hours, note: $('#fNote').value, photo_key: photoKey, is_default: $('#fDefault').checked,
      google_place_id: S.placeId || ''};
    try {
      const pid = S.editing.id;
      await api(`/v1/seller/stores/${id}/pickup-points${pid ? '/' + pid : ''}`, {method: pid ? 'PUT' : 'POST', body});
      S.editing = null; toast(t('saved')); render();
    } catch (e) {
      showError(e, {label: '#fLabel', address_line: '#fAddr', contact_name: '#fCName', contact_phone: '#fCPhone', pickup_hours: '.hours',
        note: '#fNote', photo_key: '#fPhoto', media: '#fPhoto', google_place_id: '#map', _location: '#map', lat: '#map', lng: '#map', pickup_points: '#fSave'});
    }
  };
}

window.addEventListener('hashchange', () => {
  S.editing = null;
  if (!/\/product\//.test(location.hash)) PF.E = null; // keluar dari editor → form dibuang
  render();
});
(async () => { document.documentElement.lang = P.lang; await refreshSession(); render(); })();
