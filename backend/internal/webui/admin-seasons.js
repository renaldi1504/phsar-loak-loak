'use strict';
// Admin (super admin) → Konten: kategori musiman ("Musim ini" di beranda) + sampul kotak kategori.
addStrings({
  en: {snTitle: 'Seasonal picks', snHint: 'Cards under "This season" on the home page, shown first while their dates are active. Leave the end empty to show all year.',
    snAdd: 'Add seasonal pick', snColor: 'Card colour', snc_pink: 'Pink', snc_blue: 'Blue', snc_sand: 'Sand', snc_green: 'Green',
    snSrcCat: 'A category (best selling first)', snSrcSearch: 'Products matching words', snSrcProd: 'Products I pick', snQuery: 'Search words', snQueryPh: 'e.g. party dress, gift',
    snAllYear: 'All year', snPreview: 'Phone preview', cvTitle: 'Category tiles', cvHint: 'Each tile on the home page shows the category icon. Upload your own cover to replace it (square, at least 300 × 300 px).',
    cvAuto: 'Category icon', cvCustom: 'Your cover', cvUse: 'Use icon'},
  km: {snTitle: 'ប្រភេទតាមរដូវ', snHint: 'កាតក្រោម "រដូវនេះ" នៅទំព័រដើម បង្ហាញមុនគេក្នុងរយៈពេលកំណត់។ ទុកថ្ងៃបញ្ចប់ទទេ ដើម្បីបង្ហាញពេញមួយឆ្នាំ។',
    snAdd: 'បន្ថែមប្រភេទតាមរដូវ', snColor: 'ពណ៌កាត', snc_pink: 'ផ្កាឈូក', snc_blue: 'ខៀវ', snc_sand: 'ខ្សាច់', snc_green: 'បៃតង',
    snSrcCat: 'ប្រភេទមួយ (លក់ដាច់មុន)', snSrcSearch: 'ទំនិញដែលត្រូវនឹងពាក្យ', snSrcProd: 'ទំនិញដែលខ្ញុំជ្រើស', snQuery: 'ពាក្យស្វែងរក', snQueryPh: 'ឧ. រ៉ូបជប់លៀង កាដូ',
    snAllYear: 'ពេញមួយឆ្នាំ', snPreview: 'មើលជាមុនលើទូរស័ព្ទ', cvTitle: 'ប្រអប់ប្រភេទ', cvHint: 'ប្រអប់នីមួយៗនៅទំព័រដើម បង្ហាញរូបតំណាងប្រភេទ។ ផ្ទុករូបគម្របផ្ទាល់ខ្លួន ដើម្បីជំនួស (ការ៉េ យ៉ាងតិច 300 × 300 px)។',
    cvAuto: 'រូបតំណាងប្រភេទ', cvCustom: 'រូបគម្របរបស់អ្នក', cvUse: 'ប្រើរូបតំណាង'},
  zh: {snTitle: '季节分类', snHint: '首页"本季精选"下的卡片，在有效期内排在最前。结束时间留空则全年显示。',
    snAdd: '添加季节分类', snColor: '卡片颜色', snc_pink: '粉色', snc_blue: '蓝色', snc_sand: '沙色', snc_green: '绿色',
    snSrcCat: '某个分类（销量优先）', snSrcSearch: '匹配关键词的商品', snSrcProd: '我挑选的商品', snQuery: '搜索关键词', snQueryPh: '例如：派对裙、礼物',
    snAllYear: '全年', snPreview: '手机预览', cvTitle: '分类方块', cvHint: '首页每个方块显示分类图标。可上传自己的封面替换（正方形，至少 300 × 300 像素）。',
    cvAuto: '分类图标', cvCustom: '自定义封面', cvUse: '改用图标'},
});

const SN_COLORS = {pink: '#f50750', blue: '#032ea1', sand: '#ffe9b8', green: '#0b7a44'};

function snCard(x) {
  const fg = x.color === 'sand' ? '#102241' : '#fff';
  const tag = x.ends_at ? new Date(new Date(x.ends_at).getTime() - 1000).toLocaleDateString() : t('snAllYear');
  return `<div class="sn-card" style="background:${SN_COLORS[x.color] || SN_COLORS.pink};color:${fg}"><small>${esc(tag)}</small><b>${esc(ctIn(x.title) || '…')}</b></div>`;
}

// HTML bagian musim + sampul (dipanggil drawContent).
function ctSeasonsHTML() {
  const d = CNT.data, s = CNT.season;
  const catName = id => (d.categories.find(k => k.id === id) || {}).name || id;
  const prodName = id => { const p = (d.products || []).concat(CNT.pickedInfo || []).find(x => x.id === id); return p ? p.name : id.slice(0, 8); };
  const dates = x => `<div class="grid"><label>${t('ctStart')}<input type="datetime-local" name="starts_at" value="${ctLocal(x.starts_at)}"></label>
    <label>${t('ctEnd')}<input type="datetime-local" name="ends_at" value="${ctLocal(x.ends_at)}"></label></div>`;
  const form = s ? `<form class="ct-form" id="ctSeason" onsubmit="return false">
      ${ctLangFields('title', t('ctTitle'), s.title, 40)}
      <label>${t('snColor')}<span class="sn-colors">${Object.keys(SN_COLORS).map(c => `<label class="sn-color"><input type="radio" name="color" value="${c}" ${(s.color || 'pink') === c ? 'checked' : ''}><span style="background:${SN_COLORS[c]}"></span>${t('snc_' + c)}</label>`).join('')}</span></label>
      <div class="opts" role="radiogroup" aria-label="${t('ctSource')}">
        <label class="opt"><input type="radio" name="source" value="category" ${s.source === 'category' ? 'checked' : ''}><span><b>${t('snSrcCat')}</b></span></label>
        <label class="opt"><input type="radio" name="source" value="search" ${s.source === 'search' ? 'checked' : ''}><span><b>${t('snSrcSearch')}</b></span></label>
        <label class="opt"><input type="radio" name="source" value="products" ${s.source === 'products' ? 'checked' : ''}><span><b>${t('snSrcProd')}</b></span></label></div>
      <label data-src="category" ${s.source === 'category' ? '' : 'hidden'}>${t('ctCategory')}<select name="category_id">${d.categories.map(k => `<option value="${esc(k.id)}" ${s.category_id === k.id ? 'selected' : ''}>${esc(k.name)}</option>`).join('')}</select></label>
      <label data-src="search" ${s.source === 'search' ? '' : 'hidden'}>${t('snQuery')}<input name="query" maxlength="60" value="${esc(s.query || '')}" placeholder="${t('snQueryPh')}"></label>
      <div data-src="products" ${s.source === 'products' ? '' : 'hidden'}>
        <label>${t('ctPick')}<input id="snPickQ" name="product_ids" placeholder="${t('ctPickPh')}" maxlength="80" autocomplete="off"></label>
        <div id="snPickRes" class="ct-picks"></div>
        <p class="small"><b>${t('ctPicked', {n: CNT.spick.length})}</b></p>
        <div class="ct-chips">${CNT.spick.map(id => `<span class="vchip">${esc(prodName(id))}<button type="button" data-sunpick="${id}" aria-label="${t('ctRemove')}">×</button></span>`).join('')}</div></div>
      ${dates(s)}
      <div class="grid"><label>${t('ctSort')}<input type="number" name="sort_order" min="-1000" max="1000" value="${s.sort_order || 0}"></label>
        <label class="inline"><input type="checkbox" name="active" ${s.active !== false ? 'checked' : ''}> ${t('ctActive')}</label></div>
      <div><span class="muted small">${t('snPreview')}</span>${snCard(s)}</div>
      <div class="row"><button class="primary" type="button" id="ctSeasonSave">${t('ctSave')}</button><button type="button" id="ctSeasonCancel">${t('ctCancel')}</button></div></form>` : '';
  const rows = (d.seasons || []).map(x => `<div class="ct-row">${snCard(x)}
      <div class="ct-main"><b>${esc(ctIn(x.title))}</b><span class="muted small">${x.source === 'category' ? `${t('ctCategory')}: ${esc(catName(x.category_id))}` : x.source === 'search' ? `${t('snQuery')}: ${esc(x.query)}` : t('ctPicked', {n: x.product_ids.length})}
        · ${x.starts_at ? new Date(x.starts_at).toLocaleDateString() : '…'} → ${x.ends_at ? new Date(x.ends_at).toLocaleDateString() : t('snAllYear')}</span></div>
      <span class="chip ${CT_CHIP[ctState(x)]}">${t('cts_' + ctState(x))}</span>
      <button type="button" data-sedit="${x.id}">${t('ctEdit')}</button><button type="button" class="danger" data-sdel="${x.id}">${t('ctDelete')}</button></div>`).join('');
  const seasons = `<section class="panel"><div class="ph"><h2>${t('snTitle')}</h2>${s ? '' : `<button type="button" class="primary" id="ctSeasonAdd">${t('snAdd')}</button>`}</div>
    <div class="ct-body"><p class="muted small">${t('snHint')}</p>${form}${rows || (s ? '' : `<p class="muted">${t('ctNone')}</p>`)}</div></section>`;

  const covers = `<section class="panel"><div class="ph"><h2>${t('cvTitle')}</h2></div><div class="ct-body"><p class="muted small">${t('cvHint')}</p>
    <div class="cv-grid">${d.categories.map(k => `<div class="cv-item">
      <div class="cv-img">${k.cover_url ? `<img src="${esc(k.cover_url)}" alt="">` : ''}</div>
      <b>${esc(k.name)}</b><span class="muted small">${t(k.cover_url ? 'cvCustom' : 'cvAuto')}</span>
      <div class="row"><label class="btn sm" style="flex-direction:row">${t('ctChoose')}<input type="file" data-cover="${esc(k.id)}" accept="image/jpeg,image/png,image/webp" class="hidden"></label>
      ${k.cover_url ? `<button type="button" class="sm" data-cover-auto="${esc(k.id)}">${t('cvUse')}</button>` : ''}</div></div>`).join('')}</div></div></section>`;
  return seasons + covers;
}

function bindSeasons(save) {
  const d = CNT.data;
  if ($('#ctSeasonAdd')) $('#ctSeasonAdd').onclick = () => {
    CNT.season = {active: true, title: {}, color: 'pink', source: 'category', category_id: (d.categories[0] || {}).id}; CNT.spick = [];
    CNT.banner = null; CNT.col = null; drawContent(); $('#ctSeason').scrollIntoView({block: 'center'});
  };
  document.querySelectorAll('[data-sedit]').forEach(b => b.onclick = () => {
    CNT.season = JSON.parse(JSON.stringify(d.seasons.find(x => x.id === b.dataset.sedit))); CNT.spick = CNT.season.product_ids.slice();
    drawContent(); $('#ctSeason').scrollIntoView({block: 'center'});
  });
  document.querySelectorAll('[data-sdel]').forEach(b => b.onclick = () => { if (confirm(t('ctConfirmDel'))) save('/seasons/' + b.dataset.sdel, 'DELETE'); });
  document.querySelectorAll('[data-cover]').forEach(inp => inp.onchange = e => {
    const f = e.target.files[0]; if (!f) return;
    ctUpload(f, up => { if (up) save('/categories/' + encodeURIComponent(inp.dataset.cover) + '/cover', 'PUT', {image_key: up.key}); });
  });
  document.querySelectorAll('[data-cover-auto]').forEach(b => b.onclick = () => save('/categories/' + encodeURIComponent(b.dataset.coverAuto) + '/cover', 'PUT', {image_key: ''}));
  const f = $('#ctSeason');
  if (!f) return;
  const keep = () => Object.assign(CNT.season, {title: ctCollect(f, 'title'), color: f.querySelector('[name="color"]:checked').value,
    source: f.querySelector('[name="source"]:checked').value, category_id: f.querySelector('[name="category_id"]').value, query: f.querySelector('[name="query"]').value,
    starts_at: ctISO(f.querySelector('[name="starts_at"]').value), ends_at: ctISO(f.querySelector('[name="ends_at"]').value),
    sort_order: Number(f.querySelector('[name="sort_order"]').value) || 0, active: f.querySelector('[name="active"]').checked});
  f.querySelectorAll('[name="source"]').forEach(r => r.onchange = () => {
    const v = f.querySelector('[name="source"]:checked').value;
    f.querySelectorAll('[data-src]').forEach(el => { el.hidden = el.dataset.src !== v; });
  });
  f.querySelectorAll('[data-g="title"], [name="color"]').forEach(el => el.addEventListener('input', () => {
    keep(); f.querySelector('.sn-card').outerHTML = snCard(CNT.season);
  }));
  let timer;
  $('#snPickQ').oninput = e => {
    clearTimeout(timer);
    const q = e.target.value.trim();
    timer = setTimeout(async () => {
      if (q.length < 2) { $('#snPickRes').innerHTML = ''; return; }
      try {
        const {products} = await api('/admin/api/content/products?q=' + encodeURIComponent(q));
        CNT.pickedInfo = (CNT.pickedInfo || []).concat(products);
        $('#snPickRes').innerHTML = products.length ? products.map(p => `<button type="button" class="ct-pick" data-spick="${p.id}" ${CNT.spick.includes(p.id) ? 'disabled' : ''}>
          ${p.image_url ? `<img src="${esc(p.image_url)}" alt="">` : '<span class="ct-noimg"></span>'}<span><b>${esc(p.name)}</b><small class="muted">${esc(p.store)}</small></span></button>`).join('')
          : `<p class="muted small">${t('ctNoProducts')}</p>`;
        document.querySelectorAll('[data-spick]').forEach(b => b.onclick = () => { if (CNT.spick.length < 60) { keep(); CNT.spick.push(b.dataset.spick); drawContent(); } });
      } catch (err) { toast(errMsg(err), 'error'); }
    }, 250);
  };
  document.querySelectorAll('[data-sunpick]').forEach(b => b.onclick = () => { keep(); CNT.spick = CNT.spick.filter(x => x !== b.dataset.sunpick); drawContent(); });
  $('#ctSeasonCancel').onclick = () => { CNT.season = null; drawContent(); };
  $('#ctSeasonSave').onclick = () => {
    keep();
    const x = CNT.season;
    save(x.id ? '/seasons/' + x.id : '/seasons', x.id ? 'PUT' : 'POST', {title: x.title, color: x.color, source: x.source, category_id: x.category_id,
      product_ids: x.source === 'products' ? CNT.spick : [], query: x.query, starts_at: x.starts_at, ends_at: x.ends_at, sort_order: x.sort_order, active: x.active},
    f, () => { CNT.season = null; });
  };
}
