'use strict';
// tour: tur "cara pakai" Seller Center v2 (mesin tur di tour.js).
//   seller_intro   — "Kenalan dengan Seller Center": otomatis sekali saat toko aktif pertama kali membuka "Hari ini".
//   seller_orders  — pertama kali membuka Pesanan.
//   seller_product — pertama kali membuka wizard Tambah produk.
// Status "sudah dilihat" disimpan di server (GET/POST /v1/seller/tours) supaya tidak muncul lagi di perangkat lain.
// Kode lama hanya memanggil tourSeller(id) setelah render toko dan tourWizard() setelah editor produk (baris `// tour:`).
addStrings({
  en: {tour_next: 'Next', tour_back: 'Back', tour_skip: 'Skip', tour_done: 'Got it', tour_of: '{a}/{b}',
    tour_i0T: 'Welcome! These are your sections', tour_i0B: '{s}. Everything you need is one tap away.',
    tour_i1T: 'Your to-do list', tour_i1B: 'Things to do now appear here, most urgent first. Tap the button on a task to do it.',
    tour_i2T: 'Orders in 3 steps', tour_i2B: '1 {a} · 2 {b} · 3 {c}. Then a Phsar courier picks the parcel up from you.',
    tour_i3T: 'Add a product', tour_i3B: 'Open Products and tap “+ Add product”. A short video sells best.',
    tour_i4T: 'Your money', tour_i4B: 'See your balance, add your bank account and follow your payouts.',
    tour_i5T: 'More', tour_i5B: 'Shop profile, opening hours, team and Help. You can watch this tour again from here.',
    tour_o0B: 'New paid orders wait here. Accept them quickly — orders not accepted in time cancel automatically.',
    tour_o1B: 'Pack accepted orders, print the label and stick it on the parcel.',
    tour_o2B: 'Packed orders wait here for the Phsar courier to pick them up.',
    tour_oStep: 'Step {n} · {s}', tour_o3T: 'Order history', tour_o3B: 'Orders on the way, finished and cancelled orders are kept here.',
    tour_p0T: '4 short steps', tour_p0B: '{s}. You can go back to any step.',
    tour_p1T: 'A tip for every step', tour_p1B: 'Each step shows a short tip here. Follow it and buyers trust your product more.',
    tour_p2T: 'Save anytime', tour_p2B: 'Not finished? Save a draft and continue later.',
    tour_p3T: 'Go to the next step', tour_p3B: 'When this step is done, go to the next one. Phsar checks your product automatically — most go live straight away.',
    tour_again: 'Watch the tour again', tour_againD: 'A quick guide to Seller Center',
    tour_helpH: 'Guided tours', tour_helpP: 'Short guides that show how Seller Center works.',
    tour_rIntro: 'Seller Center', tour_rOrders: 'Orders', tour_rProduct: 'Add product'},
  km: {tour_next: 'បន្ទាប់', tour_back: 'ថយក្រោយ', tour_skip: 'រំលង', tour_done: 'យល់ហើយ', tour_of: '{a}/{b}',
    tour_i0T: 'សូមស្វាគមន៍! នេះជាផ្នែករបស់អ្នក', tour_i0B: '{s}។ អ្វីៗដែលអ្នកត្រូវការ គ្រាន់តែចុចម្តង។',
    tour_i1T: 'បញ្ជីកិច្ចការរបស់អ្នក', tour_i1B: 'កិច្ចការដែលត្រូវធ្វើឥឡូវនេះ បង្ហាញនៅទីនេះ ដោយកិច្ចការបន្ទាន់បំផុតនៅខាងលើ។ ចុចប៊ូតុងលើកិច្ចការ ដើម្បីធ្វើវា។',
    tour_i2T: 'ការបញ្ជាទិញ ៣ ជំហាន', tour_i2B: '១ {a} · ២ {b} · ៣ {c}។ បន្ទាប់មក អ្នកដឹកជញ្ជូន Phsar មកយកកញ្ចប់ពីអ្នក។',
    tour_i3T: 'បន្ថែមផលិតផល', tour_i3B: 'បើក «ផលិតផល» ហើយចុច «+ បន្ថែមផលិតផល»។ វីដេអូខ្លីលក់ដាច់ជាងគេ។',
    tour_i4T: 'លុយរបស់អ្នក', tour_i4B: 'មើលសមតុល្យ បញ្ចូលគណនីធនាគារ និងតាមដានការទូទាត់របស់អ្នក។',
    tour_i5T: 'ច្រើនទៀត', tour_i5B: 'ព័ត៌មានហាង ម៉ោងបើក ក្រុម និងជំនួយ។ អ្នកអាចមើលការណែនាំនេះម្តងទៀតពីទីនេះ។',
    tour_o0B: 'ការបញ្ជាទិញថ្មីដែលបានបង់ប្រាក់ រង់ចាំនៅទីនេះ។ សូមទទួលឱ្យបានឆាប់ — ការបញ្ជាទិញដែលមិនបានទទួលទាន់ពេល នឹងត្រូវបោះបង់ដោយស្វ័យប្រវត្តិ។',
    tour_o1B: 'វេចខ្ចប់ការបញ្ជាទិញដែលបានទទួល បោះពុម្ពស្លាក ហើយបិទលើកញ្ចប់។',
    tour_o2B: 'ការបញ្ជាទិញដែលវេចខ្ចប់រួច រង់ចាំអ្នកដឹកជញ្ជូន Phsar មកយកនៅទីនេះ។',
    tour_oStep: 'ជំហាន {n} · {s}', tour_o3T: 'ប្រវត្តិការបញ្ជាទិញ', tour_o3B: 'ការបញ្ជាទិញកំពុងដឹក បានបញ្ចប់ និងបានបោះបង់ ត្រូវរក្សាទុកនៅទីនេះ។',
    tour_p0T: '៤ ជំហានខ្លីៗ', tour_p0B: '{s}។ អ្នកអាចត្រឡប់ទៅជំហានណាមួយក៏បាន។',
    tour_p1T: 'គន្លឹះសម្រាប់ជំហាននីមួយៗ', tour_p1B: 'ជំហាននីមួយៗ បង្ហាញគន្លឹះខ្លីមួយនៅទីនេះ។ ធ្វើតាមវា អ្នកទិញនឹងទុកចិត្តផលិតផលរបស់អ្នកកាន់តែច្រើន។',
    tour_p2T: 'រក្សាទុកពេលណាក៏បាន', tour_p2B: 'មិនទាន់រួច? រក្សាទុកជាព្រាង ហើយបន្តពេលក្រោយ។',
    tour_p3T: 'ទៅជំហានបន្ទាប់', tour_p3B: 'ពេលជំហាននេះរួច សូមទៅជំហានបន្ទាប់។ Phsar ពិនិត្យផលិតផលរបស់អ្នកដោយស្វ័យប្រវត្តិ — ភាគច្រើនដាក់លក់ភ្លាមៗ។',
    tour_again: 'មើលការណែនាំម្តងទៀត', tour_againD: 'ការណែនាំខ្លីអំពី Seller Center',
    tour_helpH: 'ការណែនាំជាជំហាន', tour_helpP: 'ការណែនាំខ្លីៗ ដែលបង្ហាញពីរបៀបប្រើ Seller Center។',
    tour_rIntro: 'Seller Center', tour_rOrders: 'ការបញ្ជាទិញ', tour_rProduct: 'បន្ថែមផលិតផល'},
  zh: {tour_next: '下一步', tour_back: '上一步', tour_skip: '跳过', tour_done: '知道了', tour_of: '{a}/{b}',
    tour_i0T: '欢迎！这些是你的主要板块', tour_i0B: '{s}。你需要的一切，点一下就能找到。',
    tour_i1T: '待办事项', tour_i1B: '现在需要处理的事情会显示在这里，最紧急的排在最前面。点击任务上的按钮即可处理。',
    tour_i2T: '订单三步走', tour_i2B: '1 {a} · 2 {b} · 3 {c}。之后 Phsar 快递员会上门取件。',
    tour_i3T: '添加商品', tour_i3B: '打开「商品」，点击「+ 添加商品」。短视频最能带动销量。',
    tour_i4T: '我的收入', tour_i4B: '查看余额、添加银行账户，并跟踪你的结算。',
    tour_i5T: '更多', tour_i5B: '店铺资料、营业时间、团队和帮助。你也可以在这里重新观看本导览。',
    tour_o0B: '已付款的新订单在这里等待。请尽快接单——超时未接的订单会自动取消。',
    tour_o1B: '把已接的订单打包，打印面单并贴在包裹上。',
    tour_o2B: '打包好的订单在这里等待 Phsar 快递员上门取件。',
    tour_oStep: '第 {n} 步 · {s}', tour_o3T: '订单记录', tour_o3B: '配送中、已完成和已取消的订单都保存在这里。',
    tour_p0T: '四个简单步骤', tour_p0B: '{s}。可以随时返回任一步骤。',
    tour_p1T: '每一步都有提示', tour_p1B: '每一步都会在这里显示简短提示。照着做，买家会更信任你的商品。',
    tour_p2T: '随时保存', tour_p2B: '还没填完？先保存草稿，稍后继续。',
    tour_p3T: '进入下一步', tour_p3B: '完成本步骤后进入下一步。Phsar 会自动检查商品——大多数会直接上架。',
    tour_again: '重新观看导览', tour_againD: 'Seller Center 快速指南',
    tour_helpH: '使用导览', tour_helpP: '简短指南，带你了解 Seller Center 的用法。',
    tour_rIntro: 'Seller Center', tour_rOrders: '订单', tour_rProduct: '添加商品'},
});

const TOUR_KEYS = ['seller_intro', 'seller_orders', 'seller_product'];
const tourSel = {uid: null, done: null, loading: null, force: null, timer: 0};
const tourUid = () => (P.user && P.user.id) || 'me';

// Status dari server (sekali per akun per muat halaman) + cadangan lokal kalau POST sempat gagal.
async function tourSellerDone() {
  if (tourSel.uid === tourUid() && tourSel.done) return tourSel.done;
  if (tourSel.uid === tourUid() && tourSel.loading) return tourSel.loading;
  tourSel.uid = tourUid();
  tourSel.loading = (async () => {
    const done = {};
    try { Object.assign(done, (await api('/v1/seller/tours')).done || {}); } catch (e) { return null; } // jaringan gagal: jangan ganggu seller
    for (const k of TOUR_KEYS) {
      let pending = false;
      try { pending = localStorage.getItem('phsar-tour-' + tourSel.uid + '-' + k) === '1'; } catch (e) { /* mode privat */ }
      if (pending && !done[k]) { done[k] = true; tourSellerSave(k); }
    }
    tourSel.done = done;
    return done;
  })();
  const d = await tourSel.loading;
  tourSel.loading = null;
  return d;
}

async function tourSellerSave(key) {
  const lk = 'phsar-tour-' + tourUid() + '-' + key;
  if (tourSel.done) tourSel.done[key] = true;
  try {
    await api('/v1/seller/tours', {method: 'POST', body: {tour_key: key}});
    try { localStorage.removeItem(lk); } catch (e) { /* abaikan */ }
  } catch (e) {
    try { localStorage.setItem(lk, '1'); } catch (e2) { /* abaikan */ }
  }
}

function tourSellerLabels() {
  return {next: t('tour_next'), back: t('tour_back'), skip: t('tour_skip'), done: t('tour_done'), of: t('tour_of')};
}

function tourSellerSteps(key) {
  const nav = k => `.v2nav a[href$="/${k}"]`;
  if (key === 'seller_intro') {
    const owner = S.store && S.store.role === 'owner';
    const secs = (owner ? ['home', 'orders', 'products', 'payouts', 'more'] : ['home', 'orders', 'products', 'more']).map(k => t('nav_' + k)).join(' · ');
    return [
      {sel: '.v2nav', title: t('tour_i0T'), body: t('tour_i0B', {s: secs})},
      {sel: ['#tab .task', '#tab .alldone'], title: t('tour_i1T'), body: t('tour_i1B')},
      {sel: nav('orders'), title: t('tour_i2T'), body: t('tour_i2B', {a: t('os_new'), b: t('os_processing'), c: t('os_ready')})},
      {sel: ['#tab a[href$="/product/new"]', nav('products')], title: t('tour_i3T'), body: t('tour_i3B')},
      {sel: nav('payouts'), title: t('tour_i4T'), body: t('tour_i4B')},
      {sel: nav('more'), title: t('tour_i5T'), body: t('tour_i5B')},
    ];
  }
  if (key === 'seller_orders') {
    return [
      {sel: '.otabs [data-of="new"]', title: t('tour_oStep', {n: 1, s: t('os_new')}), body: t('tour_o0B')},
      {sel: '.otabs [data-of="processing"]', title: t('tour_oStep', {n: 2, s: t('os_processing')}), body: t('tour_o1B')},
      {sel: '.otabs [data-of="ready"]', title: t('tour_oStep', {n: 3, s: t('os_ready')}), body: t('tour_o2B')},
      {sel: '.ohist', title: t('tour_o3T'), body: t('tour_o3B')},
    ];
  }
  return [
    {sel: '.wizbar', title: t('tour_p0T'), body: t('tour_p0B', {s: ['media', 'basic', 'price', 'check'].map(k => t('wz_' + k)).join(' · ')})},
    {sel: '.wizhint', title: t('tour_p1T'), body: t('tour_p1B')},
    {sel: '#btnDraft', title: t('tour_p2T'), body: t('tour_p2B')},
    {sel: '.wizfoot .primary', title: t('tour_p3T'), body: t('tour_p3B')},
  ];
}

// Mulai tur kalau belum pernah dilihat (atau dipaksa dari "Lihat tur lagi"). check() memastikan layar masih sama.
function tourSellerRun(key, check) {
  clearTimeout(tourSel.timer);
  const forced = tourSel.force === key;
  tourSel.timer = setTimeout(async () => {
    if (tourIsOn() || document.querySelector('dialog[open]') || !check()) return;
    if (!forced) {
      const done = await tourSellerDone();
      if (!done || done[key] || tourIsOn() || !check()) return;
    }
    if (tourSel.force === key) tourSel.force = null;
    tourStart({steps: tourSellerSteps(key), labels: tourSellerLabels(), onEnd: () => tourSellerSave(key)});
  }, 350);
}

// Dipanggil setelah halaman toko digambar (seller.js renderStore).
function tourSeller(id) {
  try {
    const st = S.store;
    if (!st || typeof tourStart !== 'function') return;
    if (S.tab === 'more') tourSellerMoreItem(id);
    if (S.tab === 'help') tourSellerHelpCard(id);
    const same = tab => () => S.storeId === id && S.tab === tab && !!$('#tab');
    if (S.tab === 'home' && (st.status === 'active' || tourSel.force === 'seller_intro')) tourSellerRun('seller_intro', same('home'));
    else if (S.tab === 'orders' && (st.status === 'active' || tourSel.force === 'seller_orders')) tourSellerRun('seller_orders', same('orders'));
  } catch (e) { /* tur tidak boleh merusak halaman */ }
}

// Dipanggil setelah wizard produk digambar (seller-products.js renderProductEditor). Hanya untuk produk baru.
function tourWizard() {
  try {
    if (typeof PF === 'undefined' || !PF.E) return;
    const forced = tourSel.force === 'seller_product';
    if (PF.E.product && !forced) return;
    const sid = PF.E.storeId;
    tourSellerRun('seller_product', () => !!PF.E && PF.E.storeId === sid && !!$('.wizbar') && /\/product\//.test(location.hash));
  } catch (e) { /* abaikan */ }
}

// "Lihat tur lagi": buka layar yang benar lalu jalankan tur walaupun sudah pernah dilihat.
function tourReplay(key) {
  if (!TOUR_KEYS.includes(key)) return;
  tourSel.force = key;
  if (key === 'seller_product') {
    const h = '#/store/' + S.storeId + '/product/new';
    if (location.hash === h) render(); else location.hash = h;
    return;
  }
  const tab = key === 'seller_orders' ? 'orders' : 'home';
  if (key === 'seller_orders' && typeof SO !== 'undefined') Object.assign(SO, {filter: 'new', q: '', page: 1});
  v2Open(tab);
}

// Item "Lihat tur lagi" di menu Lainnya (ditambahkan ke grup Bantuan tanpa mengubah seller-home.js).
function tourSellerMoreItem(id) {
  if ($('#tourMoreItem')) return;
  const help = document.querySelector(`#tab .mitem[href$="/${id}/help"]`);
  if (!help) return;
  const a = document.createElement('a');
  a.className = 'mitem'; a.id = 'tourMoreItem'; a.href = '#'; a.style.setProperty('--mb', '#ffe4ec');
  a.innerHTML = `<span class="mi" aria-hidden="true">🧭</span><span class="mx"><b>${esc(t('tour_again'))}</b><span>${esc(t('tour_againD'))}</span></span>`;
  a.onclick = e => { e.preventDefault(); tourReplay('seller_intro'); };
  help.insertAdjacentElement('afterend', a);
}

// Kartu tur di halaman Bantuan.
function tourSellerHelpCard() {
  const tab = $('#tab');
  if (!tab || $('#tourHelpCard') || !$('#hpNewBtn')) return; // hanya di daftar bantuan, bukan di dalam percakapan tiket
  const keys = TOUR_KEYS;
  const lab = {seller_intro: 'tour_rIntro', seller_orders: 'tour_rOrders', seller_product: 'tour_rProduct'};
  const sec = document.createElement('section');
  sec.className = 'card tourhelp'; sec.id = 'tourHelpCard';
  sec.innerHTML = `<div><h2>🧭 ${esc(t('tour_helpH'))}</h2><p class="muted">${esc(t('tour_helpP'))}</p></div>
    <div class="tourbtns">${keys.map(k => `<button type="button" data-tour-replay="${k}">${esc(t(lab[k]))}</button>`).join('')}</div>`;
  sec.querySelectorAll('[data-tour-replay]').forEach(b => { b.onclick = () => tourReplay(b.dataset.tourReplay); });
  tab.appendChild(sec);
}
