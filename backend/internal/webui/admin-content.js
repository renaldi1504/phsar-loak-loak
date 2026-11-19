'use strict';
// Admin (super admin): konten & SEO website pembeli — pita pengumuman, banner beranda, koleksi pilihan, meta SEO, gambar OG.
addStrings({
  en: {ctDesc: 'What buyers see on the website. Changes appear within 30 seconds.', ctAnn: 'Announcement bar', ctAnnHint: 'A thin bar above every page — delivery news, holidays, a promo.',
    ctOn: 'Show the bar', ctColor: 'Colour', ctc_pink: 'Pink', ctc_navy: 'Navy', ctc_green: 'Green', ctc_amber: 'Amber', ctLink: 'Link (optional)', ctLinkPh: '/km/c/fashion or https://…',
    ctText: 'Text', ctStart: 'Start (optional)', ctEnd: 'End (optional)', ctSave: 'Save', ctSaved: 'Saved', ctBanners: 'Home banners', ctBannersHint: 'Wide images in the home carousel, after the search slide. Best size 1600 × 600 px, JPG/PNG/WebP up to 5 MB.',
    ctAddBanner: 'Add banner', ctEdit: 'Edit', ctDelete: 'Delete', ctConfirmDel: 'Delete this item?', ctImage: 'Image', ctChoose: 'Choose image', ctUploading: 'Uploading…',
    ctTitle: 'Title', ctSub: 'Subtitle', ctSort: 'Order', ctActive: 'Active', ctCancel: 'Cancel', ctNone: 'Nothing yet.',
    cts_live: 'Live', cts_scheduled: 'Scheduled', cts_ended: 'Ended', cts_off: 'Off', ctCols: 'Featured collections', ctColsHint: 'Rows of products on the home page, above the videos.',
    ctAddCol: 'Add collection', ctSource: 'Products from', ctSrcCat: 'A category (newest first)', ctSrcProd: 'Products I pick', ctCategory: 'Category', ctPick: 'Find products',
    ctPickPh: 'Product or store name', ctPicked: 'Picked ({n})', ctMax: 'Show up to', ctItems: '{n} products', ctSEO: 'SEO — search results & sharing',
    ctSEOHint: 'Leave empty to use the automatic text. Google shows about 60 characters of the title and 160 of the description.', ctPageHome: 'Home page',
    ctMetaTitle: 'Page title', ctMetaDesc: 'Description', ctOG: 'Default sharing image', ctOGHint: 'Shown when a page is shared on Facebook or Telegram (1200 × 630 px). Product pages use the product photo.',
    ctRemove: 'Remove', ctCustom: 'custom', ctAuto: 'automatic', ctNoProducts: 'No products found.'},
  km: {ctDesc: 'អ្វីដែលអ្នកទិញឃើញនៅលើគេហទំព័រ។ ការផ្លាស់ប្ដូរបង្ហាញក្នុង 30 វិនាទី។', ctAnn: 'របារសេចក្តីជូនដំណឹង', ctAnnHint: 'របារតូចមួយខាងលើគ្រប់ទំព័រ — ព័ត៌មានដឹកជញ្ជូន ថ្ងៃឈប់សម្រាក ឬប្រូម៉ូសិន។',
    ctOn: 'បង្ហាញរបារ', ctColor: 'ពណ៌', ctc_pink: 'ផ្កាឈូក', ctc_navy: 'ខៀវចាស់', ctc_green: 'បៃតង', ctc_amber: 'លឿង', ctLink: 'តំណ (ស្រេចចិត្ត)', ctLinkPh: '/km/c/fashion ឬ https://…',
    ctText: 'អត្ថបទ', ctStart: 'ចាប់ផ្ដើម (ស្រេចចិត្ត)', ctEnd: 'បញ្ចប់ (ស្រេចចិត្ត)', ctSave: 'រក្សាទុក', ctSaved: 'បានរក្សាទុក', ctBanners: 'បដាទំព័រដើម', ctBannersHint: 'រូបភាពធំក្នុងផ្ទាំងរំកិលទំព័រដើម បន្ទាប់ពីផ្ទាំងស្វែងរក។ ទំហំល្អបំផុត 1600 × 600 px, JPG/PNG/WebP រហូតដល់ 5 MB។',
    ctAddBanner: 'បន្ថែមបដា', ctEdit: 'កែ', ctDelete: 'លុប', ctConfirmDel: 'លុបធាតុនេះ?', ctImage: 'រូបភាព', ctChoose: 'ជ្រើសរូបភាព', ctUploading: 'កំពុងផ្ទុកឡើង…',
    ctTitle: 'ចំណងជើង', ctSub: 'ចំណងជើងរង', ctSort: 'លំដាប់', ctActive: 'សកម្ម', ctCancel: 'បោះបង់', ctNone: 'មិនទាន់មានទេ។',
    cts_live: 'កំពុងបង្ហាញ', cts_scheduled: 'កំណត់ពេល', cts_ended: 'បានបញ្ចប់', cts_off: 'បិទ', ctCols: 'បណ្តុំពិសេស', ctColsHint: 'ជួរផលិតផលនៅទំព័រដើម ខាងលើវីដេអូ។',
    ctAddCol: 'បន្ថែមបណ្តុំ', ctSource: 'ផលិតផលពី', ctSrcCat: 'ប្រភេទមួយ (ថ្មីបំផុតមុន)', ctSrcProd: 'ផលិតផលដែលខ្ញុំជ្រើស', ctCategory: 'ប្រភេទ', ctPick: 'ស្វែងរកផលិតផល',
    ctPickPh: 'ឈ្មោះផលិតផល ឬហាង', ctPicked: 'បានជ្រើស ({n})', ctMax: 'បង្ហាញរហូតដល់', ctItems: 'ផលិតផល {n}', ctSEO: 'SEO — លទ្ធផលស្វែងរក និងការចែករំលែក',
    ctSEOHint: 'ទុកទទេដើម្បីប្រើអត្ថបទស្វ័យប្រវត្តិ។ Google បង្ហាញចំណងជើងប្រហែល 60 តួ និងការពិពណ៌នា 160 តួ។', ctPageHome: 'ទំព័រដើម',
    ctMetaTitle: 'ចំណងជើងទំព័រ', ctMetaDesc: 'ការពិពណ៌នា', ctOG: 'រូបភាពចែករំលែកលំនាំដើម', ctOGHint: 'បង្ហាញពេលចែករំលែកទំព័រលើ Facebook ឬ Telegram (1200 × 630 px)។ ទំព័រផលិតផលប្រើរូបផលិតផល។',
    ctRemove: 'ដកចេញ', ctCustom: 'ផ្ទាល់ខ្លួន', ctAuto: 'ស្វ័យប្រវត្តិ', ctNoProducts: 'រកមិនឃើញផលិតផល។'},
  zh: {ctDesc: '买家在网站上看到的内容。修改会在 30 秒内生效。', ctAnn: '公告栏', ctAnnHint: '每个页面顶部的细条——配送通知、假期、促销。',
    ctOn: '显示公告栏', ctColor: '颜色', ctc_pink: '粉色', ctc_navy: '深蓝', ctc_green: '绿色', ctc_amber: '琥珀色', ctLink: '链接（可选）', ctLinkPh: '/km/c/fashion 或 https://…',
    ctText: '文字', ctStart: '开始（可选）', ctEnd: '结束（可选）', ctSave: '保存', ctSaved: '已保存', ctBanners: '首页横幅', ctBannersHint: '首页轮播中的宽图，位于搜索幻灯片之后。最佳尺寸 1600 × 600 像素，JPG/PNG/WebP，最大 5 MB。',
    ctAddBanner: '添加横幅', ctEdit: '编辑', ctDelete: '删除', ctConfirmDel: '删除此项？', ctImage: '图片', ctChoose: '选择图片', ctUploading: '正在上传…',
    ctTitle: '标题', ctSub: '副标题', ctSort: '顺序', ctActive: '启用', ctCancel: '取消', ctNone: '暂无。',
    cts_live: '展示中', cts_scheduled: '已排期', cts_ended: '已结束', cts_off: '关闭', ctCols: '精选合集', ctColsHint: '首页上的商品行，位于视频上方。',
    ctAddCol: '添加合集', ctSource: '商品来源', ctSrcCat: '某个类目（最新优先）', ctSrcProd: '我挑选的商品', ctCategory: '类目', ctPick: '查找商品',
    ctPickPh: '商品或店铺名称', ctPicked: '已选（{n}）', ctMax: '最多显示', ctItems: '{n} 件商品', ctSEO: 'SEO——搜索结果与分享',
    ctSEOHint: '留空则使用自动文字。Google 大约显示标题 60 个字符、描述 160 个字符。', ctPageHome: '首页',
    ctMetaTitle: '页面标题', ctMetaDesc: '描述', ctOG: '默认分享图片', ctOGHint: '在 Facebook 或 Telegram 分享页面时显示（1200 × 630 像素）。商品页使用商品照片。',
    ctRemove: '移除', ctCustom: '自定义', ctAuto: '自动', ctNoProducts: '未找到商品。'},
});

const CNT = {data: null, banner: null, col: null, pick: [], seoOpen: null, season: null, spick: []};
const CNT_LANGS = [['en', 'English'], ['km', 'ខ្មែរ'], ['zh', '中文']];

async function renderContent() {
  CNT.data = await api('/admin/api/content');
  drawContent();
}

// ISO ↔ <input type="datetime-local"> (waktu lokal browser)
function ctLocal(iso) {
  if (!iso) return '';
  const d = new Date(iso), z = n => String(n).padStart(2, '0');
  return `${d.getFullYear()}-${z(d.getMonth() + 1)}-${z(d.getDate())}T${z(d.getHours())}:${z(d.getMinutes())}`;
}
function ctISO(v) { return v ? new Date(v).toISOString() : null; }
function ctState(x) {
  const now = Date.now();
  if (!x.active) return 'off';
  if (x.ends_at && new Date(x.ends_at) <= now) return 'ended';
  if (x.starts_at && new Date(x.starts_at) > now) return 'scheduled';
  return 'live';
}
const CT_CHIP = {live: 'ok', scheduled: 'warn', ended: '', off: ''};
function ctTx(o, l) { return (o && o[l]) || ''; }
function ctIn(o) { return ctTx(o, P.lang) || ctTx(o, 'en') || ctTx(o, 'km') || ctTx(o, 'zh'); }

function ctLangFields(prefix, label, obj, max, tag = 'input') {
  return `<div class="grid">${CNT_LANGS.map(([l, name]) => `<label>${label} · ${name}${tag === 'textarea'
    ? `<textarea name="${prefix}_${l}" data-l="${l}" data-g="${prefix}" maxlength="${max}" rows="2">${esc(ctTx(obj, l))}</textarea>`
    : `<input name="${prefix}_${l}" data-l="${l}" data-g="${prefix}" maxlength="${max}" value="${esc(ctTx(obj, l))}">`}</label>`).join('')}</div>`;
}
function ctCollect(form, prefix) {
  const out = {};
  form.querySelectorAll(`[data-g="${prefix}"]`).forEach(el => { const v = el.value.trim(); if (v) out[el.dataset.l] = v; });
  return out;
}
function ctFail(e, form) {
  form.querySelectorAll('.invalid').forEach(el => el.classList.remove('invalid'));
  if (e.field) { const el = form.querySelector(`[name="${e.field}"]`); if (el) { el.classList.add('invalid'); el.focus(); } }
  toast(errMsg(e), 'error');
}
async function ctUpload(file, done) {
  const fd = new FormData(); fd.append('file', file);
  try { done(await api('/admin/api/content/upload', {method: 'POST', body: fd})); } catch (e) { toast(errMsg(e), 'error'); done(null); }
}

function drawContent() {
  const d = CNT.data, a = d.announcement || {text: {}};
  const dates = (x) => `<div class="grid"><label>${t('ctStart')}<input type="datetime-local" name="starts_at" value="${ctLocal(x.starts_at)}"></label>
    <label>${t('ctEnd')}<input type="datetime-local" name="ends_at" value="${ctLocal(x.ends_at)}"></label></div>`;
  const prodName = id => { const p = (d.products || []).concat(CNT.pickedInfo || []).find(x => x.id === id); return p ? p.name : id.slice(0, 8); };

  // --- pengumuman
  const ann = `<section class="panel" id="ctAnn"><div class="ph"><h2>${t('ctAnn')}</h2><span class="chip ${CT_CHIP[ctState(a)]}">${t('cts_' + ctState(a))}</span></div>
    <div class="ct-body"><p class="muted small">${t('ctAnnHint')}</p>
      <div class="grid"><label class="inline"><input type="checkbox" name="active" ${a.active ? 'checked' : ''}> ${t('ctOn')}</label>
        <label>${t('ctColor')}<select name="color">${['pink', 'navy', 'green', 'amber'].map(c => `<option value="${c}" ${a.color === c ? 'selected' : ''}>${t('ctc_' + c)}</option>`).join('')}</select></label>
        <label>${t('ctLink')}<input name="link" maxlength="300" value="${esc(a.link || '')}" placeholder="${t('ctLinkPh')}"></label></div>
      ${ctLangFields('text', t('ctText'), a.text, 160)}${dates(a)}
      <div class="ct-preview announce announce-${esc(a.color || 'pink')}" id="ctAnnPrev">${esc(ctIn(a.text) || '…')}</div>
      <div><button class="primary" type="button" id="ctAnnSave">${t('ctSave')}</button></div></div></section>`;

  // --- banner
  const b = CNT.banner;
  const bannerForm = b ? `<form class="ct-form" id="ctBanner" onsubmit="return false">
      <div class="ct-img">${b.image_url ? `<img src="${esc(b.image_url)}" alt="">` : `<span class="muted">${t('ctImage')}</span>`}</div>
      <label class="btn" style="flex-direction:row;align-self:flex-start">${t('ctChoose')}<input type="file" id="ctBannerFile" name="image_key" accept="image/jpeg,image/png,image/webp" class="hidden"></label>
      <label>${t('ctLink')}<input name="link" maxlength="300" value="${esc(b.link || '')}" placeholder="${t('ctLinkPh')}"></label>
      ${ctLangFields('title', t('ctTitle'), b.title, 80)}${ctLangFields('subtitle', t('ctSub'), b.subtitle, 140)}${dates(b)}
      <div class="grid"><label>${t('ctSort')}<input type="number" name="sort_order" min="-1000" max="1000" value="${b.sort_order || 0}"></label>
        <label class="inline"><input type="checkbox" name="active" ${b.active !== false ? 'checked' : ''}> ${t('ctActive')}</label></div>
      <div class="row"><button class="primary" type="button" id="ctBannerSave">${t('ctSave')}</button><button type="button" id="ctBannerCancel">${t('ctCancel')}</button></div></form>` : '';
  const bannerRows = d.banners.map(x => `<div class="ct-row"><img src="${esc(x.image_url)}" alt="" class="ct-thumb">
      <div class="ct-main"><b>${esc(ctIn(x.title) || '—')}</b><span class="muted small">${esc(x.link || '')}${x.starts_at || x.ends_at ? ` · ${x.starts_at ? new Date(x.starts_at).toLocaleString() : '…'} → ${x.ends_at ? new Date(x.ends_at).toLocaleString() : '…'}` : ''}</span></div>
      <span class="chip ${CT_CHIP[ctState(x)]}">${t('cts_' + ctState(x))}</span>
      <button type="button" data-bedit="${x.id}">${t('ctEdit')}</button><button type="button" class="danger" data-bdel="${x.id}">${t('ctDelete')}</button></div>`).join('');
  const banners = `<section class="panel"><div class="ph"><h2>${t('ctBanners')}</h2>${b ? '' : `<button type="button" class="primary" id="ctBannerAdd">${t('ctAddBanner')}</button>`}</div>
    <div class="ct-body"><p class="muted small">${t('ctBannersHint')}</p>${bannerForm}${bannerRows || (b ? '' : `<p class="muted">${t('ctNone')}</p>`)}</div></section>`;

  // --- koleksi
  const c = CNT.col;
  const colForm = c ? `<form class="ct-form" id="ctCol" onsubmit="return false">
      ${ctLangFields('title', t('ctTitle'), c.title, 60)}
      <div class="opts" role="radiogroup" aria-label="${t('ctSource')}">
        <label class="opt"><input type="radio" name="source" value="category" ${c.source !== 'products' ? 'checked' : ''}><span><b>${t('ctSrcCat')}</b></span></label>
        <label class="opt"><input type="radio" name="source" value="products" ${c.source === 'products' ? 'checked' : ''}><span><b>${t('ctSrcProd')}</b></span></label></div>
      <label id="ctCatWrap" ${c.source === 'products' ? 'hidden' : ''}>${t('ctCategory')}<select name="category_id">${d.categories.map(k => `<option value="${esc(k.id)}" ${c.category_id === k.id ? 'selected' : ''}>${esc(k.name)}</option>`).join('')}</select></label>
      <div id="ctPickWrap" ${c.source === 'products' ? '' : 'hidden'}>
        <label>${t('ctPick')}<input id="ctPickQ" name="product_ids" placeholder="${t('ctPickPh')}" maxlength="80" autocomplete="off"></label>
        <div id="ctPickRes" class="ct-picks"></div>
        <p class="small"><b>${t('ctPicked', {n: CNT.pick.length})}</b></p>
        <div class="ct-chips">${CNT.pick.map(id => `<span class="vchip">${esc(prodName(id))}<button type="button" data-unpick="${id}" aria-label="${t('ctRemove')}">×</button></span>`).join('')}</div></div>
      <div class="grid"><label>${t('ctMax')}<select name="max_items">${[4, 6, 8, 12, 16, 24].map(n => `<option value="${n}" ${(c.max_items || 8) === n ? 'selected' : ''}>${t('ctItems', {n})}</option>`).join('')}</select></label>
        <label>${t('ctSort')}<input type="number" name="sort_order" min="-1000" max="1000" value="${c.sort_order || 0}"></label>
        <label class="inline"><input type="checkbox" name="active" ${c.active !== false ? 'checked' : ''}> ${t('ctActive')}</label></div>
      <div class="row"><button class="primary" type="button" id="ctColSave">${t('ctSave')}</button><button type="button" id="ctColCancel">${t('ctCancel')}</button></div></form>` : '';
  const catName = id => (d.categories.find(k => k.id === id) || {}).name || id;
  const colRows = d.collections.map(x => `<div class="ct-row"><div class="ct-main"><b>${esc(ctIn(x.title))}</b>
      <span class="muted small">${x.source === 'category' ? `${t('ctCategory')}: ${esc(catName(x.category_id))}` : t('ctPicked', {n: x.product_ids.length})} · ${t('ctItems', {n: x.max_items})}</span></div>
      <span class="chip ${x.active ? 'ok' : ''}">${t(x.active ? 'cts_live' : 'cts_off')}</span>
      <button type="button" data-cedit="${x.id}">${t('ctEdit')}</button><button type="button" class="danger" data-cdel="${x.id}">${t('ctDelete')}</button></div>`).join('');
  const cols = `<section class="panel"><div class="ph"><h2>${t('ctCols')}</h2>${c ? '' : `<button type="button" class="primary" id="ctColAdd">${t('ctAddCol')}</button>`}</div>
    <div class="ct-body"><p class="muted small">${t('ctColsHint')}</p>${colForm}${colRows || (c ? '' : `<p class="muted">${t('ctNone')}</p>`)}</div></section>`;

  // --- SEO
  const seoBy = Object.fromEntries((d.seo || []).map(m => [m.page, m]));
  const pages = [['home', t('ctPageHome')]].concat(d.categories.map(k => ['c:' + k.id, k.name]));
  const seoRows = pages.map(([page, label]) => {
    const m = seoBy[page], open = CNT.seoOpen === page;
    return `<div class="ct-seo"><button type="button" class="ct-seo-head" data-seo="${esc(page)}" aria-expanded="${open}"><b>${esc(label)}</b>
        <span class="chip ${m ? 'pink' : ''}">${t(m ? 'ctCustom' : 'ctAuto')}</span>${m ? `<span class="muted small ct-seo-t">${esc(ctIn(m.title))}</span>` : ''}</button>
      ${open ? `<form class="ct-form" id="ctSEO" data-page="${esc(page)}" onsubmit="return false">${ctLangFields('title', t('ctMetaTitle'), m && m.title, 70)}
        ${ctLangFields('description', t('ctMetaDesc'), m && m.description, 170, 'textarea')}
        <div class="row"><button class="primary" type="button" id="ctSEOSave">${t('ctSave')}</button></div></form>` : ''}</div>`;
  }).join('');
  const seo = `<section class="panel"><div class="ph"><h2>${t('ctSEO')}</h2></div><div class="ct-body"><p class="muted small">${t('ctSEOHint')}</p>
      <div class="ct-row"><div class="ct-img ct-og">${d.og_image_url ? `<img src="${esc(d.og_image_url)}" alt="">` : '<span class="muted">1200 × 630</span>'}</div>
        <div class="ct-main"><b>${t('ctOG')}</b><span class="muted small">${t('ctOGHint')}</span>
          <div class="row"><label class="btn" style="flex-direction:row">${t('ctChoose')}<input type="file" id="ctOGFile" accept="image/jpeg,image/png,image/webp" class="hidden"></label>
          ${d.og_image_key ? `<button type="button" id="ctOGDel">${t('ctRemove')}</button>` : ''}</div></div></div>
      ${seoRows}</div></section>`;

  $('#app').innerHTML = shell(`${pageHead(t('i_content'), t('ctDesc'))}${ann}${banners}${ctSeasonsHTML()}${cols}${seo}`);
  bindContent();
}

function bindContent() {
  const d = CNT.data;
  const save = async (path, method, body, form, after) => {
    try { CNT.data = await api('/admin/api/content' + path, {method, body}); if (after) after(); toast(t('ctSaved')); drawContent(); }
    catch (e) { ctFail(e, form || $('#app')); }
  };
  // pengumuman
  const af = $('#ctAnn');
  af.querySelectorAll('[data-g="text"], [name="color"]').forEach(el => el.oninput = () => {
    const prev = $('#ctAnnPrev'), txt = ctCollect(af, 'text');
    prev.textContent = txt[P.lang] || txt.en || txt.km || txt.zh || '…';
    prev.className = 'ct-preview announce announce-' + af.querySelector('[name="color"]').value;
  });
  $('#ctAnnSave').onclick = () => save('/announcement', 'PUT', {active: af.querySelector('[name="active"]').checked, color: af.querySelector('[name="color"]').value,
    link: af.querySelector('[name="link"]').value, text: ctCollect(af, 'text'),
    starts_at: ctISO(af.querySelector('[name="starts_at"]').value), ends_at: ctISO(af.querySelector('[name="ends_at"]').value)}, af);

  bindSeasons(save); // kategori musiman + sampul kotak kategori

  // banner
  if ($('#ctBannerAdd')) $('#ctBannerAdd').onclick = () => { CNT.banner = {active: true, title: {}, subtitle: {}}; CNT.col = null; drawContent(); $('#ctBanner').scrollIntoView({block: 'center'}); };
  document.querySelectorAll('[data-bedit]').forEach(btn => btn.onclick = () => { CNT.banner = JSON.parse(JSON.stringify(d.banners.find(x => x.id === btn.dataset.bedit))); drawContent(); $('#ctBanner').scrollIntoView({block: 'center'}); });
  document.querySelectorAll('[data-bdel]').forEach(btn => btn.onclick = () => { if (confirm(t('ctConfirmDel'))) save('/banners/' + btn.dataset.bdel, 'DELETE'); });
  const bf = $('#ctBanner');
  if (bf) {
    const keep = () => { // simpan isian supaya tidak hilang saat gambar diganti
      Object.assign(CNT.banner, {link: bf.querySelector('[name="link"]').value, title: ctCollect(bf, 'title'), subtitle: ctCollect(bf, 'subtitle'),
        starts_at: ctISO(bf.querySelector('[name="starts_at"]').value), ends_at: ctISO(bf.querySelector('[name="ends_at"]').value),
        sort_order: Number(bf.querySelector('[name="sort_order"]').value) || 0, active: bf.querySelector('[name="active"]').checked});
    };
    $('#ctBannerFile').onchange = e => {
      const f = e.target.files[0]; if (!f) return;
      keep(); bf.querySelector('.ct-img').innerHTML = `<span class="muted">${t('ctUploading')}</span>`;
      ctUpload(f, up => { if (up) { CNT.banner.image_key = up.key; CNT.banner.image_url = up.url; } drawContent(); });
    };
    $('#ctBannerCancel').onclick = () => { CNT.banner = null; drawContent(); };
    $('#ctBannerSave').onclick = () => {
      keep();
      const x = CNT.banner;
      save(x.id ? '/banners/' + x.id : '/banners', x.id ? 'PUT' : 'POST', {image_key: x.image_key || '', link: x.link, title: x.title, subtitle: x.subtitle,
        starts_at: x.starts_at, ends_at: x.ends_at, sort_order: x.sort_order, active: x.active}, bf, () => { CNT.banner = null; });
    };
  }

  // koleksi
  if ($('#ctColAdd')) $('#ctColAdd').onclick = () => { CNT.col = {active: true, title: {}, source: 'category', max_items: 8}; CNT.pick = []; CNT.banner = null; drawContent(); $('#ctCol').scrollIntoView({block: 'center'}); };
  document.querySelectorAll('[data-cedit]').forEach(btn => btn.onclick = () => {
    CNT.col = JSON.parse(JSON.stringify(d.collections.find(x => x.id === btn.dataset.cedit))); CNT.pick = CNT.col.product_ids.slice(); drawContent(); $('#ctCol').scrollIntoView({block: 'center'});
  });
  document.querySelectorAll('[data-cdel]').forEach(btn => btn.onclick = () => { if (confirm(t('ctConfirmDel'))) save('/collections/' + btn.dataset.cdel, 'DELETE'); });
  const cf = $('#ctCol');
  if (cf) {
    const keep = () => Object.assign(CNT.col, {title: ctCollect(cf, 'title'), source: cf.querySelector('[name="source"]:checked').value,
      category_id: cf.querySelector('[name="category_id"]').value, max_items: Number(cf.querySelector('[name="max_items"]').value),
      sort_order: Number(cf.querySelector('[name="sort_order"]').value) || 0, active: cf.querySelector('[name="active"]').checked});
    cf.querySelectorAll('[name="source"]').forEach(r => r.onchange = () => {
      const prod = cf.querySelector('[name="source"]:checked').value === 'products';
      $('#ctCatWrap').hidden = prod; $('#ctPickWrap').hidden = !prod;
    });
    let timer;
    $('#ctPickQ').oninput = e => {
      clearTimeout(timer);
      const q = e.target.value.trim();
      timer = setTimeout(async () => {
        if (q.length < 2) { $('#ctPickRes').innerHTML = ''; return; }
        try {
          const {products} = await api('/admin/api/content/products?q=' + encodeURIComponent(q));
          CNT.pickedInfo = (CNT.pickedInfo || []).concat(products);
          $('#ctPickRes').innerHTML = products.length ? products.map(p => `<button type="button" class="ct-pick" data-pick="${p.id}" ${CNT.pick.includes(p.id) ? 'disabled' : ''}>
            ${p.image_url ? `<img src="${esc(p.image_url)}" alt="">` : '<span class="ct-noimg"></span>'}<span><b>${esc(p.name)}</b><small class="muted">${esc(p.store)}</small></span></button>`).join('')
            : `<p class="muted small">${t('ctNoProducts')}</p>`;
          document.querySelectorAll('[data-pick]').forEach(b => b.onclick = () => { if (CNT.pick.length < 24) { keep(); CNT.pick.push(b.dataset.pick); drawContent(); } });
        } catch (err) { toast(errMsg(err), 'error'); }
      }, 250);
    };
    document.querySelectorAll('[data-unpick]').forEach(b => b.onclick = () => { keep(); CNT.pick = CNT.pick.filter(x => x !== b.dataset.unpick); drawContent(); });
    $('#ctColCancel').onclick = () => { CNT.col = null; drawContent(); };
    $('#ctColSave').onclick = () => {
      keep();
      const x = CNT.col;
      save(x.id ? '/collections/' + x.id : '/collections', x.id ? 'PUT' : 'POST', {title: x.title, source: x.source, category_id: x.category_id,
        product_ids: x.source === 'products' ? CNT.pick : [], max_items: x.max_items, sort_order: x.sort_order, active: x.active}, cf, () => { CNT.col = null; });
    };
  }

  // SEO
  document.querySelectorAll('[data-seo]').forEach(b => b.onclick = () => { CNT.seoOpen = CNT.seoOpen === b.dataset.seo ? null : b.dataset.seo; drawContent(); });
  const sf = $('#ctSEO');
  if (sf) $('#ctSEOSave').onclick = () => save('/seo/' + encodeURIComponent(sf.dataset.page), 'PUT', {title: ctCollect(sf, 'title'), description: ctCollect(sf, 'description')}, sf);
  $('#ctOGFile').onchange = e => { const f = e.target.files[0]; if (f) ctUpload(f, up => { if (up) save('/og', 'PUT', {image_key: up.key}); }); };
  if ($('#ctOGDel')) $('#ctOGDel').onclick = () => save('/og', 'PUT', {image_key: ''});
}
