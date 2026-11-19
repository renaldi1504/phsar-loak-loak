'use strict';
// Admin: Categories & brands (izin catalog.manage) — kategori, atribut kategori, merek, kata terlarang untuk moderasi.
addStrings({
  en: {ctHint: 'What buyers browse and sellers choose. Changes show on the website within a minute.', ctTabCats: 'Categories', ctTabBrands: 'Brands', ctTabWords: 'Prohibited words',
    ctAdd: 'Add category', ctEdit: 'Edit category', ctId: 'ID (permanent, e.g. pet-supplies)', ctParent: 'Main category', ctParentNone: '— none (main category) —', ctIcon: 'Icon',
    ctUsed: 'Used / second-hand items allowed', ctUsedNo: 'New only', ctName: 'Name', ctSlug: 'Web address', ctSeo: 'Search engine (optional)', ctMetaTitle: 'Page title (max 70)',
    ctMetaDesc: 'Description (max 160)', ctSave: 'Save', ctCancel: 'Cancel', ctSaved: 'Saved', ctProducts: '{n} products · {l} live', ctHidden: 'Hidden', ctHide: 'Hide', ctShow: 'Show',
    ctDelete: 'Delete', ctUp: 'Move up', ctDown: 'Move down', ctConfirmHide: 'Hide “{c}”? Buyers will not see it and sellers cannot choose it. Existing products stay.',
    ctConfirmDelete: 'Delete “{c}” permanently?', ctNoDelete: 'Has products — hide it instead', ctAttrs: 'Attributes sellers fill in', ctAttrsHint: 'Extra details for this category, for example RAM or material.',
    ctCode: 'Code (e.g. ram_gb)', ctLabel: 'Label', ctType: 'Type', ctType_text: 'Text', ctType_number: 'Number', ctType_bool: 'Yes / no', ctType_select: 'Choice list',
    ctChoices: 'Choices (comma separated)', ctRequired: 'Required', ctFilterable: 'Filter on website', ctAddAttr: 'Add attribute', ctNoAttrs: 'No attributes yet.',
    ctBrand: 'Brand', ctBrandName: 'Brand name', ctBrandSlug: 'Web address (auto if empty)', ctVerified: 'Verified brand', ctActive: 'Active', ctInactive: 'Inactive',
    ctBrandCats: 'Shown in categories (none = all)', ctAddBrand: 'Add brand', ctEditBrand: 'Edit brand', ctLogo: 'Upload logo', ctLogoHint: 'JPG, PNG or WebP, max 2 MB.',
    ctNoBrands: 'No brands yet.', ctConfirmDelBrand: 'Delete brand “{c}”?', ctWord: 'Word', ctLang: 'Language', ctGroup: 'Group (e.g. weapons)', ctAction: 'When found',
    ctAction_review: 'Send to moderator', ctAction_block: 'Block', ctAddWord: 'Add word', ctWordsHint: 'Checked in product names, descriptions and condition notes when a seller sends a product. Block = the seller can’t send it. Send to moderator = it waits for a check instead of going live automatically.',
    ctNoWords: 'No words yet.', ctRemove: 'Remove', ctAll: 'all',
    ctRestricted: 'Restricted — every product is checked by a moderator before it goes live', ctRestrictedShort: 'Restricted',
    ctTabAuto: 'Auto-check', arHint: 'Most products go live immediately. These rules decide which ones wait for a moderator first. Products that went live automatically appear under Products → Auto-published.',
    arAuto: 'Publish clean products automatically', arAutoOff: 'Off: every product waits for a moderator.', arProbation: 'New store: first products checked by a moderator',
    arDiscount: 'Crossed-out price more than … % above the price', arDrop: 'Price dropped more than … % since approval', arReports: 'Hide a product after serious reports from … different buyers',
    arSave: 'Save rules', pfTitle: 'Minimum believable price', pfHint: 'If a product name contains the words and its cheapest option is below this price, a moderator checks it first (stops fake iPhones for $50). Names with “case”, “charger”, “strap”… are skipped.',
    pfKeyword: 'Words in the name', pfMin: 'Minimum price (USD)', pfNote: 'Note', pfAdd: 'Add', pfNone: 'No price rules yet.', pfEdit: 'Edit',},
  km: {ctHint: 'អ្វីដែលអ្នកទិញរកមើល និងអ្នកលក់ជ្រើសរើស។ ការផ្លាស់ប្តូរបង្ហាញលើគេហទំព័រក្នុងមួយនាទី។', ctTabCats: 'ប្រភេទ', ctTabBrands: 'ម៉ាក', ctTabWords: 'ពាក្យហាមឃាត់',
    ctAdd: 'បន្ថែមប្រភេទ', ctEdit: 'កែប្រភេទ', ctId: 'ID (អចិន្ត្រៃយ៍ ឧ. pet-supplies)', ctParent: 'ប្រភេទមេ', ctParentNone: '— គ្មាន (ប្រភេទមេ) —', ctIcon: 'រូបតំណាង',
    ctUsed: 'អនុញ្ញាតទំនិញប្រើរួច', ctUsedNo: 'ថ្មីប៉ុណ្ណោះ', ctName: 'ឈ្មោះ', ctSlug: 'អាសយដ្ឋានវេប', ctSeo: 'ម៉ាស៊ីនស្វែងរក (ស្រេចចិត្ត)', ctMetaTitle: 'ចំណងជើងទំព័រ (អតិបរមា 70)',
    ctMetaDesc: 'ការពិពណ៌នា (អតិបរមា 160)', ctSave: 'រក្សាទុក', ctCancel: 'បោះបង់', ctSaved: 'បានរក្សាទុក', ctProducts: 'ផលិតផល {n} · កំពុងលក់ {l}', ctHidden: 'បានលាក់', ctHide: 'លាក់', ctShow: 'បង្ហាញ',
    ctDelete: 'លុប', ctUp: 'ផ្លាស់ឡើង', ctDown: 'ផ្លាស់ចុះ', ctConfirmHide: 'លាក់ «{c}»? អ្នកទិញនឹងមិនឃើញ ហើយអ្នកលក់មិនអាចជ្រើសបាន។ ផលិតផលដែលមានស្រាប់នៅដដែល។',
    ctConfirmDelete: 'លុប «{c}» ជាអចិន្ត្រៃយ៍?', ctNoDelete: 'មានផលិតផល — សូមលាក់ជំនួសវិញ', ctAttrs: 'លក្ខណៈដែលអ្នកលក់បំពេញ', ctAttrsHint: 'ព័ត៌មានបន្ថែមសម្រាប់ប្រភេទនេះ ឧទាហរណ៍ RAM ឬសម្ភារៈ។',
    ctCode: 'កូដ (ឧ. ram_gb)', ctLabel: 'ស្លាក', ctType: 'ប្រភេទទិន្នន័យ', ctType_text: 'អត្ថបទ', ctType_number: 'លេខ', ctType_bool: 'បាទ/ទេ', ctType_select: 'បញ្ជីជម្រើស',
    ctChoices: 'ជម្រើស (បំបែកដោយក្បៀស)', ctRequired: 'ចាំបាច់', ctFilterable: 'តម្រងលើគេហទំព័រ', ctAddAttr: 'បន្ថែមលក្ខណៈ', ctNoAttrs: 'មិនទាន់មានលក្ខណៈ។',
    ctBrand: 'ម៉ាក', ctBrandName: 'ឈ្មោះម៉ាក', ctBrandSlug: 'អាសយដ្ឋានវេប (ស្វ័យប្រវត្តិ បើទទេ)', ctVerified: 'ម៉ាកបានផ្ទៀងផ្ទាត់', ctActive: 'សកម្ម', ctInactive: 'អសកម្ម',
    ctBrandCats: 'បង្ហាញក្នុងប្រភេទ (គ្មាន = ទាំងអស់)', ctAddBrand: 'បន្ថែមម៉ាក', ctEditBrand: 'កែម៉ាក', ctLogo: 'ផ្ទុកឡូហ្គោ', ctLogoHint: 'JPG, PNG ឬ WebP អតិបរមា 2 MB។',
    ctNoBrands: 'មិនទាន់មានម៉ាក។', ctConfirmDelBrand: 'លុបម៉ាក «{c}»?', ctWord: 'ពាក្យ', ctLang: 'ភាសា', ctGroup: 'ក្រុម (ឧ. weapons)', ctAction: 'ពេលរកឃើញ',
    ctAction_review: 'បញ្ជូនទៅអ្នកត្រួតពិនិត្យ', ctAction_block: 'ហាមឃាត់', ctAddWord: 'បន្ថែមពាក្យ', ctWordsHint: 'ពិនិត្យក្នុងឈ្មោះ ការពិពណ៌នា និងកំណត់ចំណាំស្ថានភាព ពេលអ្នកលក់ផ្ញើផលិតផល។ ហាមឃាត់ = មិនអាចផ្ញើបាន។ បញ្ជូនទៅអ្នកត្រួតពិនិត្យ = រង់ចាំការពិនិត្យ ជំនួសការបង្ហាញស្វ័យប្រវត្តិ។',
    ctNoWords: 'មិនទាន់មានពាក្យ។', ctRemove: 'ដកចេញ', ctAll: 'ទាំងអស់',
    ctRestricted: 'ដាក់កម្រិត — ផលិតផលទាំងអស់ត្រូវអ្នកត្រួតពិនិត្យមើលមុនបង្ហាញ', ctRestrictedShort: 'ដាក់កម្រិត',
    ctTabAuto: 'ពិនិត្យស្វ័យប្រវត្តិ', arHint: 'ផលិតផលភាគច្រើនបង្ហាញភ្លាមៗ។ ច្បាប់ទាំងនេះកំណត់ថាផលិតផលណាត្រូវរង់ចាំអ្នកត្រួតពិនិត្យសិន។ ផលិតផលដែលបង្ហាញស្វ័យប្រវត្តិ មាននៅ ផលិតផល → បង្ហាញស្វ័យប្រវត្តិ។',
    arAuto: 'បង្ហាញផលិតផលស្អាតដោយស្វ័យប្រវត្តិ', arAutoOff: 'បិទ៖ ផលិតផលទាំងអស់រង់ចាំអ្នកត្រួតពិនិត្យ។', arProbation: 'ហាងថ្មី៖ ផលិតផលដំបូងៗត្រូវអ្នកត្រួតពិនិត្យមើល',
    arDiscount: 'តម្លៃគូសចោលខ្ពស់ជាងតម្លៃលក់ច្រើនជាង … %', arDrop: 'តម្លៃធ្លាក់ច្រើនជាង … % ចាប់តាំងពីអនុម័ត', arReports: 'លាក់ផលិតផលបន្ទាប់ពីរបាយការណ៍ធ្ងន់ធ្ងរពីអ្នកទិញ … នាក់ផ្សេងគ្នា',
    arSave: 'រក្សាទុកច្បាប់', pfTitle: 'តម្លៃទាបបំផុតដែលគួរឲ្យជឿ', pfHint: 'បើឈ្មោះផលិតផលមានពាក្យទាំងនេះ ហើយជម្រើសថោកបំផុតទាបជាងតម្លៃនេះ អ្នកត្រួតពិនិត្យនឹងមើលសិន។ ឈ្មោះមាន “case” “charger” “strap”… ត្រូវរំលង។',
    pfKeyword: 'ពាក្យក្នុងឈ្មោះ', pfMin: 'តម្លៃអប្បបរមា (USD)', pfNote: 'កំណត់ចំណាំ', pfAdd: 'បន្ថែម', pfNone: 'មិនទាន់មានច្បាប់តម្លៃ។', pfEdit: 'កែ',},
  zh: {ctHint: '买家浏览、卖家选择的分类。修改会在一分钟内显示在网站上。', ctTabCats: '分类', ctTabBrands: '品牌', ctTabWords: '违禁词',
    ctAdd: '添加分类', ctEdit: '编辑分类', ctId: 'ID（永久，例如 pet-supplies）', ctParent: '上级分类', ctParentNone: '— 无（一级分类）—', ctIcon: '图标',
    ctUsed: '允许二手商品', ctUsedNo: '仅限新品', ctName: '名称', ctSlug: '网址', ctSeo: '搜索引擎（可选）', ctMetaTitle: '页面标题（最多 70）',
    ctMetaDesc: '描述（最多 160）', ctSave: '保存', ctCancel: '取消', ctSaved: '已保存', ctProducts: '{n} 个商品 · {l} 个在售', ctHidden: '已隐藏', ctHide: '隐藏', ctShow: '显示',
    ctDelete: '删除', ctUp: '上移', ctDown: '下移', ctConfirmHide: '隐藏“{c}”？买家将看不到，卖家也无法选择。现有商品保留。',
    ctConfirmDelete: '永久删除“{c}”？', ctNoDelete: '已有商品——请改为隐藏', ctAttrs: '卖家填写的属性', ctAttrsHint: '此分类的附加信息，例如内存或材质。',
    ctCode: '代码（例如 ram_gb）', ctLabel: '标签', ctType: '类型', ctType_text: '文字', ctType_number: '数字', ctType_bool: '是 / 否', ctType_select: '选项列表',
    ctChoices: '选项（用逗号分隔）', ctRequired: '必填', ctFilterable: '网站筛选', ctAddAttr: '添加属性', ctNoAttrs: '暂无属性。',
    ctBrand: '品牌', ctBrandName: '品牌名称', ctBrandSlug: '网址（留空自动生成）', ctVerified: '认证品牌', ctActive: '启用', ctInactive: '停用',
    ctBrandCats: '显示在分类（不选 = 全部）', ctAddBrand: '添加品牌', ctEditBrand: '编辑品牌', ctLogo: '上传标志', ctLogoHint: 'JPG、PNG 或 WebP，最大 2 MB。',
    ctNoBrands: '暂无品牌。', ctConfirmDelBrand: '删除品牌“{c}”？', ctWord: '词语', ctLang: '语言', ctGroup: '分组（例如 weapons）', ctAction: '发现时',
    ctAction_review: '转交审核员', ctAction_block: '拦截', ctAddWord: '添加词语', ctWordsHint: '卖家提交商品时检查名称、描述和成色说明。拦截 = 无法提交；转交审核员 = 不自动上架，等待审核。',
    ctNoWords: '暂无词语。', ctRemove: '移除', ctAll: '全部',
    ctRestricted: '受限——所有商品上架前须经审核员检查', ctRestrictedShort: '受限',
    ctTabAuto: '自动检查', arHint: '大多数商品会立即上架。这些规则决定哪些商品需要先由审核员检查。自动上架的商品会出现在 商品 → 自动上架。',
    arAuto: '自动上架无问题的商品', arAutoOff: '关闭：所有商品都需等待审核。', arProbation: '新店铺：前几件商品由审核员检查',
    arDiscount: '划线价比售价高出 … % 以上', arDrop: '价格较通过审核时下降 … % 以上', arReports: '收到 … 位不同买家的严重举报后隐藏商品',
    arSave: '保存规则', pfTitle: '最低可信价格', pfHint: '如果商品名称包含这些词，且最低规格价格低于此价格，将先由审核员检查（防止 50 美元的假 iPhone）。名称含“case”“charger”“strap”等配件词的会跳过。',
    pfKeyword: '名称中的词', pfMin: '最低价格（美元）', pfNote: '备注', pfAdd: '添加', pfNone: '暂无价格规则。', pfEdit: '编辑',},
});

const CT = {data: null, tab: 'cats', edit: null, brandEdit: null, floorEdit: null};
const CT_LANGS = ['en', 'km', 'zh'];
const CT_LNAME = {en: 'English', km: 'ខ្មែរ', zh: '中文'};

async function renderCatalog() {
  CT.data = await api('/admin/api/catalog');
  drawCatalog();
}

function ctName(c) { return c.names[P.lang] || c.names.en || c.id; }

function drawCatalog() {
  const tabs = ['cats', 'brands', 'words', 'auto'].map(k => `<button type="button" role="tab" aria-selected="${CT.tab === k}" data-cttab="${k}">${t('ctTab' + {cats: 'Cats', brands: 'Brands', words: 'Words', auto: 'Auto'}[k])}</button>`).join('');
  const body = {cats: ctCategoriesHTML, brands: ctBrandsHTML, words: ctWordsHTML, auto: ctAutoHTML}[CT.tab]();
  $('#app').innerHTML = shell(`${pageHead(t('i_catalog'), t('ctHint'))}<div class="tabs" role="tablist">${tabs}</div><div id="ctBody">${body}</div>`);
  document.querySelectorAll('[data-cttab]').forEach(b => b.onclick = () => { CT.tab = b.dataset.cttab; CT.edit = null; CT.brandEdit = null; drawCatalog(); });
  ({cats: bindCategories, brands: bindBrands, words: bindWords, auto: bindAuto})[CT.tab]();
}

// reply: semua endpoint mengembalikan katalog lengkap.
async function ctSave(path, opts, form, fields) {
  try {
    CT.data = await api('/admin/api/catalog' + path, opts);
    toast(t('ctSaved'));
    return true;
  } catch (e) {
    if (fields) showError(e, fields, form); else toast(errMsg(e), 'error');
    return false;
  }
}

// ---------- Kategori ----------
function ctCategoriesHTML() {
  const cats = CT.data.categories, tops = cats.filter(c => !c.parent_id);
  const row = (c, child) => {
    const sibs = cats.filter(x => (x.parent_id || '') === (c.parent_id || ''));
    const i = sibs.indexOf(c);
    return `<tr class="${c.hidden ? 'muted' : ''}"><td>${child ? '<span class="muted">↳</span> ' : ''}<b>${esc(ctName(c))}</b><br><code class="small">${esc(c.id)}</code></td>
      <td class="small">${CT_LANGS.map(l => `${esc(c.names[l] || '—')}`).join('<br>')}</td>
      <td>${esc(c.icon || '—')}</td><td>${c.allows_used ? t('ctUsed') : t('ctUsedNo')}${c.restricted ? `<br><span class="chip warn">${t('ctRestrictedShort')}</span>` : ''}</td>
      <td class="small">${esc(t('ctProducts', {n: c.products, l: c.live_products}))}${c.hidden ? ` <span class="chip warn">${t('ctHidden')}</span>` : ''}</td>
      <td style="white-space:nowrap"><button type="button" class="sm" data-move="up" data-id="${esc(c.id)}" ${i === 0 ? 'disabled' : ''} aria-label="${t('ctUp')}">↑</button>
        <button type="button" class="sm" data-move="down" data-id="${esc(c.id)}" ${i === sibs.length - 1 ? 'disabled' : ''} aria-label="${t('ctDown')}">↓</button>
        <button type="button" class="sm" data-edit="${esc(c.id)}">${t('ctEdit')}</button>
        <button type="button" class="sm" data-hide="${esc(c.id)}">${c.hidden ? t('ctShow') : t('ctHide')}</button>
        <button type="button" class="sm danger" data-del="${esc(c.id)}" ${c.products ? `disabled title="${esc(t('ctNoDelete'))}"` : ''}>${t('ctDelete')}</button></td></tr>`;
  };
  const rows = tops.map(c => row(c, false) + cats.filter(x => x.parent_id === c.id).map(x => row(x, true)).join('')).join('');
  return `${CT.edit ? ctEditorHTML() : `<div class="row"><button type="button" class="primary" id="ctAddBtn">${t('ctAdd')}</button></div>`}
    <section class="panel"><div class="tbl-wrap"><table><thead><tr><th>${t('ctTabCats')}</th><th>${t('ctName')}</th><th>${t('ctIcon')}</th><th></th><th></th><th></th></tr></thead>
    <tbody>${rows}</tbody></table></div></section>`;
}

function ctEditorHTML() {
  const isNew = CT.edit === 'new', c = isNew ? {id: '', parent_id: null, icon: '', allows_used: true, restricted: false, names: {}, slugs: {}, meta_titles: {}, meta_descriptions: {}, attributes: []}
    : CT.data.categories.find(x => x.id === CT.edit);
  const tops = CT.data.categories.filter(x => !x.parent_id && x.id !== c.id);
  const langBlock = l => `<fieldset class="ct-lang"><legend>${CT_LNAME[l]}</legend>
    <label>${t('ctName')}<input name="name_${l}" maxlength="80" value="${esc(c.names[l] || '')}"></label>
    <label>${t('ctSlug')}<input name="slug_${l}" maxlength="80" value="${esc(c.slugs[l] || '')}" placeholder="${esc(c.slugs.en || c.id || 'pet-supplies')}"></label>
    <details><summary>${t('ctSeo')}</summary><label>${t('ctMetaTitle')}<input name="meta_title_${l}" maxlength="70" value="${esc(c.meta_titles[l] || '')}"></label>
      <label>${t('ctMetaDesc')}<textarea name="meta_description_${l}" maxlength="160">${esc(c.meta_descriptions[l] || '')}</textarea></label></details></fieldset>`;
  const attrs = isNew ? '' : `<section class="ct-attrs"><h3>${t('ctAttrs')}</h3><p class="muted small">${t('ctAttrsHint')}</p>
    ${c.attributes.length ? `<ul class="ct-attr-list">${c.attributes.map(a => `<li><b>${esc(a.label[P.lang] || a.label.en)}</b> <code>${esc(a.code)}</code>
      <span class="chip">${t('ctType_' + a.data_type)}</span>${a.is_required ? ` <span class="chip pink">${t('ctRequired')}</span>` : ''}${a.is_filterable ? ` <span class="chip">${t('ctFilterable')}</span>` : ''}
      ${a.choices.length ? `<span class="muted small"> · ${esc(a.choices.join(', '))}</span>` : ''}
      <button type="button" class="sm danger" data-delattr="${a.id}" ${a.values ? 'disabled' : ''}>${t('ctRemove')}</button></li>`).join('')}</ul>` : `<p class="muted">${t('ctNoAttrs')}</p>`}
    <div id="ctAttrForm" class="grid">
      <label>${t('ctCode')}<input name="code" maxlength="40"></label>
      ${CT_LANGS.map(l => `<label>${t('ctLabel')} · ${CT_LNAME[l]}<input name="label_${l}" maxlength="60"></label>`).join('')}
      <label>${t('ctType')}<select name="data_type">${['text', 'number', 'bool', 'select'].map(k => `<option value="${k}">${t('ctType_' + k)}</option>`).join('')}</select></label>
      <label>${t('ctChoices')}<input name="choices" placeholder="4, 8, 16"></label>
      <label class="inline"><input type="checkbox" name="is_required"> ${t('ctRequired')}</label>
      <label class="inline"><input type="checkbox" name="is_filterable"> ${t('ctFilterable')}</label>
      <div><button type="button" id="ctAddAttr">${t('ctAddAttr')}</button></div></div></section>`;
  return `<section class="card" id="ctForm"><h2>${isNew ? t('ctAdd') : t('ctEdit') + ' · ' + esc(ctName(c))}</h2>
    <div class="grid">
      ${isNew ? `<label>${t('ctId')}<input name="id" maxlength="60" placeholder="pet-supplies"></label>
        <label>${t('ctParent')}<select name="parent_id"><option value="">${t('ctParentNone')}</option>${tops.map(x => `<option value="${esc(x.id)}">${esc(ctName(x))}</option>`).join('')}</select></label>` : ''}
      <label>${t('ctIcon')}<select name="icon"><option value="">—</option>${CT.data.icons.map(i => `<option ${c.icon === i ? 'selected' : ''}>${i}</option>`).join('')}</select></label>
      <label class="inline" style="align-self:end"><input type="checkbox" name="allows_used" ${c.allows_used ? 'checked' : ''}> ${t('ctUsed')}</label>
      <label class="inline" style="align-self:end"><input type="checkbox" name="restricted" ${c.restricted ? 'checked' : ''}> ${t('ctRestricted')}</label></div>
    <div class="ct-langs">${CT_LANGS.map(langBlock).join('')}</div>
    <div class="row" data-errors-anchor><button type="button" class="primary" id="ctSave">${t('ctSave')}</button><button type="button" id="ctCancel">${t('ctCancel')}</button></div>
    ${attrs}</section>`;
}

function bindCategories() {
  const add = $('#ctAddBtn');
  if (add) add.onclick = () => { CT.edit = 'new'; drawCatalog(); $('#ctForm input[name=id]').focus(); };
  document.querySelectorAll('[data-edit]').forEach(b => b.onclick = () => { CT.edit = b.dataset.edit; drawCatalog(); $('#ctForm').scrollIntoView({behavior: 'smooth'}); });
  document.querySelectorAll('[data-move]').forEach(b => b.onclick = async () => { if (await ctSave(`/categories/${b.dataset.id}/move`, {method: 'POST', body: {dir: b.dataset.move}})) drawCatalog(); });
  document.querySelectorAll('[data-hide]').forEach(b => b.onclick = async () => {
    const c = CT.data.categories.find(x => x.id === b.dataset.hide);
    if (!c.hidden && !confirm(t('ctConfirmHide', {c: ctName(c)}))) return;
    if (await ctSave(`/categories/${c.id}/hide`, {method: 'POST', body: {hidden: !c.hidden}})) drawCatalog();
  });
  document.querySelectorAll('[data-del]').forEach(b => b.onclick = async () => {
    const c = CT.data.categories.find(x => x.id === b.dataset.del);
    if (!confirm(t('ctConfirmDelete', {c: ctName(c)}))) return;
    if (await ctSave(`/categories/${c.id}`, {method: 'DELETE'})) { if (CT.edit === c.id) CT.edit = null; drawCatalog(); }
  });
  if (!CT.edit) return;
  const f = $('#ctForm'), v = n => (f.querySelector(`[name="${n}"]`) || {}).value || '';
  const fields = {id: '#ctForm [name=id]', parent_id: '#ctForm [name=parent_id]', icon: '#ctForm [name=icon]'};
  CT_LANGS.forEach(l => ['name', 'slug', 'meta_title', 'meta_description'].forEach(k => { fields[k + '_' + l] = `#ctForm [name=${k}_${l}]`; }));
  $('#ctCancel').onclick = () => { CT.edit = null; drawCatalog(); };
  $('#ctSave').onclick = async () => {
    const body = {icon: v('icon'), allows_used: f.querySelector('[name=allows_used]').checked, restricted: f.querySelector('[name=restricted]').checked, names: {}, slugs: {}, meta_titles: {}, meta_descriptions: {}};
    CT_LANGS.forEach(l => { body.names[l] = v('name_' + l); body.slugs[l] = v('slug_' + l); body.meta_titles[l] = v('meta_title_' + l); body.meta_descriptions[l] = v('meta_description_' + l); });
    const isNew = CT.edit === 'new';
    if (isNew) { body.id = v('id').trim(); body.parent_id = v('parent_id'); }
    if (await ctSave(isNew ? '/categories' : `/categories/${CT.edit}`, {method: isNew ? 'POST' : 'PUT', body}, f, fields)) { CT.edit = isNew ? body.id : CT.edit; drawCatalog(); }
  };
  const aa = $('#ctAddAttr');
  if (aa) aa.onclick = async () => {
    const af = $('#ctAttrForm'), av = n => af.querySelector(`[name="${n}"]`);
    const body = {code: av('code').value.trim(), label: {}, data_type: av('data_type').value, choices: av('choices').value.split(',').map(s => s.trim()).filter(Boolean),
      is_required: av('is_required').checked, is_filterable: av('is_filterable').checked};
    CT_LANGS.forEach(l => { body.label[l] = av('label_' + l).value; });
    const af2 = {code: '#ctAttrForm [name=code]', data_type: '#ctAttrForm [name=data_type]', choices: '#ctAttrForm [name=choices]'};
    CT_LANGS.forEach(l => { af2['label_' + l] = `#ctAttrForm [name=label_${l}]`; });
    if (await ctSave(`/categories/${CT.edit}/attributes`, {method: 'POST', body}, af, af2)) drawCatalog();
  };
  document.querySelectorAll('[data-delattr]').forEach(b => b.onclick = async () => { if (await ctSave(`/attributes/${b.dataset.delattr}`, {method: 'DELETE'})) drawCatalog(); });
}

// ---------- Merek ----------
function ctBrandsHTML() {
  const cats = CT.data.categories.filter(c => !c.parent_id), b = CT.brandEdit && CT.brandEdit !== 'new' ? CT.data.brands.find(x => x.id === CT.brandEdit) : null;
  const form = `<section class="card" id="cbForm"><h2>${b ? t('ctEditBrand') + ' · ' + esc(b.name) : t('ctAddBrand')}</h2><div class="grid">
    <label>${t('ctBrandName')}<input name="name" maxlength="80" value="${esc(b ? b.name : '')}"></label>
    <label>${t('ctBrandSlug')}<input name="slug" maxlength="60" value="${esc(b ? b.slug : '')}"></label>
    <label class="inline" style="align-self:end"><input type="checkbox" name="verified" ${b && b.verified ? 'checked' : ''}> ${t('ctVerified')}</label>
    <label class="inline" style="align-self:end"><input type="checkbox" name="active" ${!b || b.active ? 'checked' : ''}> ${t('ctActive')}</label></div>
    <fieldset><legend>${t('ctBrandCats')}</legend><div class="row">${cats.map(c => `<label class="inline"><input type="checkbox" name="cat" value="${esc(c.id)}" ${b && b.category_ids.includes(c.id) ? 'checked' : ''}> ${esc(ctName(c))}</label>`).join('')}</div></fieldset>
    <div class="row" data-errors-anchor><button type="button" class="primary" id="cbSave">${t('ctSave')}</button>${b ? `<button type="button" id="cbCancel">${t('ctCancel')}</button>` : ''}</div></section>`;
  const list = CT.data.brands.map(x => `<tr><td><span class="ct-logo">${x.logo_url ? `<img src="${esc(x.logo_url)}" alt="">` : ''}</span></td>
    <td><b>${esc(x.name)}</b>${x.verified ? ' <span class="chip ok">✓ ' + t('ctVerified') + '</span>' : ''}${x.active ? '' : ' <span class="chip">' + t('ctInactive') + '</span>'}<br><code class="small">${esc(x.slug)}</code></td>
    <td class="small">${x.category_ids.length ? x.category_ids.map(id => esc(ctName(CT.data.categories.find(c => c.id === id) || {names: {}, id}))).join(', ') : t('ctAll')}</td>
    <td class="small">${x.products}</td>
    <td style="white-space:nowrap"><label class="btn sm" style="display:inline-flex;cursor:pointer">${t('ctLogo')}<input type="file" accept="image/jpeg,image/png,image/webp" data-logo="${x.id}" hidden></label>
      <button type="button" class="sm" data-bedit="${x.id}">${t('ctEdit')}</button>
      <button type="button" class="sm danger" data-bdel="${x.id}" ${x.products ? 'disabled' : ''}>${t('ctDelete')}</button></td></tr>`).join('');
  return `${form}<p class="muted small">${t('ctLogoHint')}</p><section class="panel"><div class="tbl-wrap">${list ? `<table><tbody>${list}</tbody></table>` : `<p class="muted" style="padding:16px">${t('ctNoBrands')}</p>`}</div></section>`;
}

function bindBrands() {
  const f = $('#cbForm');
  $('#cbSave').onclick = async () => {
    const id = CT.brandEdit && CT.brandEdit !== 'new' ? CT.brandEdit : null;
    const body = {name: f.querySelector('[name=name]').value, slug: f.querySelector('[name=slug]').value, verified: f.querySelector('[name=verified]').checked,
      active: f.querySelector('[name=active]').checked, category_ids: [...f.querySelectorAll('[name=cat]:checked')].map(c => c.value)};
    if (await ctSave(id ? `/brands/${id}` : '/brands', {method: id ? 'PUT' : 'POST', body}, f, {name: '#cbForm [name=name]', slug: '#cbForm [name=slug]'})) { CT.brandEdit = null; drawCatalog(); }
  };
  const c = $('#cbCancel');
  if (c) c.onclick = () => { CT.brandEdit = null; drawCatalog(); };
  document.querySelectorAll('[data-bedit]').forEach(b => b.onclick = () => { CT.brandEdit = +b.dataset.bedit; drawCatalog(); window.scrollTo(0, 0); });
  document.querySelectorAll('[data-bdel]').forEach(b => b.onclick = async () => {
    const x = CT.data.brands.find(y => y.id === +b.dataset.bdel);
    if (confirm(t('ctConfirmDelBrand', {c: x.name})) && await ctSave(`/brands/${x.id}`, {method: 'DELETE'})) drawCatalog();
  });
  document.querySelectorAll('[data-logo]').forEach(inp => inp.onchange = async () => {
    const file = inp.files[0];
    if (!file) return;
    if (file.size > 2 * 1024 * 1024) { toast(t('ctLogoHint'), 'warning'); return; }
    const fd = new FormData(); fd.append('file', file);
    const r = await fetch(`/admin/api/catalog/brands/${inp.dataset.logo}/logo`, {method: 'POST', body: fd, headers: {Authorization: 'Bearer ' + P.token}});
    const j = await r.json().catch(() => ({}));
    if (!r.ok) { toast(errMsg(j.error || {}), 'error'); return; }
    CT.data = j; toast(t('ctSaved')); drawCatalog();
  });
}

// ---------- Kata terlarang ----------
function ctWordsHTML() {
  const rows = CT.data.keywords.map(k => `<tr><td><b>${esc(k.keyword)}</b></td><td>${esc(k.lang)}</td><td>${esc(k.category)}</td>
    <td><span class="chip ${k.action === 'block' ? 'bad' : 'warn'}">${t('ctAction_' + k.action)}</span></td>
    <td><button type="button" class="sm danger" data-kdel="${k.id}">${t('ctRemove')}</button></td></tr>`).join('');
  return `<section class="card" id="ckForm"><p class="muted small" style="margin:0">${t('ctWordsHint')}</p><div class="grid">
    <label>${t('ctWord')}<input name="keyword" maxlength="80"></label>
    <label>${t('ctLang')}<select name="lang">${CT_LANGS.map(l => `<option value="${l}">${CT_LNAME[l]}</option>`).join('')}</select></label>
    <label>${t('ctGroup')}<input name="category" maxlength="40" list="ckGroups"><datalist id="ckGroups">${[...new Set(CT.data.keywords.map(k => k.category))].map(g => `<option value="${esc(g)}">`).join('')}</datalist></label>
    <label>${t('ctAction')}<select name="action"><option value="review">${t('ctAction_review')}</option><option value="block">${t('ctAction_block')}</option></select></label></div>
    <div class="row" data-errors-anchor><button type="button" class="primary" id="ckAdd">${t('ctAddWord')}</button></div></section>
    <section class="panel"><div class="tbl-wrap">${rows ? `<table><thead><tr><th>${t('ctWord')}</th><th>${t('ctLang')}</th><th>${t('ctGroup')}</th><th>${t('ctAction')}</th><th></th></tr></thead><tbody>${rows}</tbody></table>`
      : `<p class="muted" style="padding:16px">${t('ctNoWords')}</p>`}</div></section>`;
}

function bindWords() {
  const f = $('#ckForm'), v = n => f.querySelector(`[name=${n}]`).value;
  $('#ckAdd').onclick = async () => {
    const fields = {keyword: '#ckForm [name=keyword]', lang: '#ckForm [name=lang]', category: '#ckForm [name=category]', action: '#ckForm [name=action]'};
    if (await ctSave('/keywords', {method: 'POST', body: {keyword: v('keyword'), lang: v('lang'), category: v('category'), action: v('action')}}, f, fields)) drawCatalog();
  };
  document.querySelectorAll('[data-kdel]').forEach(b => b.onclick = async () => { if (await ctSave(`/keywords/${b.dataset.kdel}`, {method: 'DELETE'})) drawCatalog(); });
}

// ---------- Cek otomatis: aturan jalur moderasi + harga minimum ----------
function ctAutoHTML() {
  const r = CT.data.moderation_rules, num = (k, min, max) => `<input type="number" name="${k}" min="${min}" max="${max}" value="${r[k]}" style="max-width:110px">`;
  const fe = CT.floorEdit && CT.data.price_floors.find(f => f.id === CT.floorEdit);
  const rows = CT.data.price_floors.map(f => `<tr><td><b>${esc(f.keyword)}</b></td><td>$${esc(f.min_price)}</td><td class="muted small">${esc(f.note)}</td>
    <td style="white-space:nowrap"><button type="button" class="sm" data-fedit="${f.id}">${t('pfEdit')}</button>
      <button type="button" class="sm danger" data-fdel="${f.id}">${t('ctRemove')}</button></td></tr>`).join('');
  return `<section class="card" id="arForm"><p class="muted small" style="margin:0">${t('arHint')}</p>
    <label class="inline"><input type="checkbox" name="auto_publish" ${r.auto_publish ? 'checked' : ''}> <b>${t('arAuto')}</b></label>
    ${r.auto_publish ? '' : `<div class="banner warn">${t('arAutoOff')}</div>`}
    <div class="grid">
      <label>${t('arProbation')}${num('probation', 0, 50)}</label>
      <label>${t('arDiscount')}${num('discount_pct', 10, 95)}</label>
      <label>${t('arDrop')}${num('drop_pct', 10, 95)}</label>
      <label>${t('arReports')}${num('report_hold', 1, 50)}</label></div>
    <div class="row" data-errors-anchor><button type="button" class="primary" id="arSave">${t('arSave')}</button></div></section>
    <section class="card" id="pfForm"><h2>${t('pfTitle')}</h2><p class="muted small" style="margin:0">${t('pfHint')}</p>
    <div class="grid">
      <label>${t('pfKeyword')}<input name="keyword" maxlength="80" placeholder="iphone 15" value="${esc(fe ? fe.keyword : '')}"></label>
      <label>${t('pfMin')}<input name="min_price" type="number" min="0.01" step="0.01" placeholder="250" value="${esc(fe ? fe.min_price : '')}"></label>
      <label>${t('pfNote')}<input name="note" maxlength="200" value="${esc(fe ? fe.note : '')}"></label></div>
    <div class="row" data-errors-anchor><button type="button" class="primary" id="pfSave">${fe ? t('ctSave') : t('pfAdd')}</button>${fe ? `<button type="button" id="pfCancel">${t('ctCancel')}</button>` : ''}</div></section>
    <section class="panel"><div class="tbl-wrap">${rows ? `<table><thead><tr><th>${t('pfKeyword')}</th><th>${t('pfMin')}</th><th>${t('pfNote')}</th><th></th></tr></thead><tbody>${rows}</tbody></table>`
      : `<p class="muted" style="padding:16px">${t('pfNone')}</p>`}</div></section>`;
}

function bindAuto() {
  const af = $('#arForm'), n = k => Number(af.querySelector(`[name=${k}]`).value);
  af.querySelector('[name=auto_publish]').onchange = e => { CT.data.moderation_rules.auto_publish = e.target.checked; };
  $('#arSave').onclick = async () => {
    const body = {auto_publish: af.querySelector('[name=auto_publish]').checked, probation: n('probation'), discount_pct: n('discount_pct'), drop_pct: n('drop_pct'), report_hold: n('report_hold')};
    const fields = {probation: '#arForm [name=probation]', discount_pct: '#arForm [name=discount_pct]', drop_pct: '#arForm [name=drop_pct]', report_hold: '#arForm [name=report_hold]'};
    if (await ctSave('/moderation-rules', {method: 'PUT', body}, af, fields)) drawCatalog();
  };
  const pf = $('#pfForm'), v = k => pf.querySelector(`[name=${k}]`).value;
  $('#pfSave').onclick = async () => {
    const fields = {keyword: '#pfForm [name=keyword]', min_price: '#pfForm [name=min_price]', note: '#pfForm [name=note]'};
    const path = CT.floorEdit ? `/price-floors/${CT.floorEdit}` : '/price-floors';
    if (await ctSave(path, {method: CT.floorEdit ? 'PUT' : 'POST', body: {keyword: v('keyword'), min_price: v('min_price'), note: v('note')}}, pf, fields)) { CT.floorEdit = null; drawCatalog(); }
  };
  const pc = $('#pfCancel');
  if (pc) pc.onclick = () => { CT.floorEdit = null; drawCatalog(); };
  document.querySelectorAll('[data-fedit]').forEach(b => b.onclick = () => { CT.floorEdit = Number(b.dataset.fedit); drawCatalog(); $('#pfForm [name=min_price]').focus(); });
  document.querySelectorAll('[data-fdel]').forEach(b => b.onclick = async () => { if (await ctSave(`/price-floors/${b.dataset.fdel}`, {method: 'DELETE'})) { CT.floorEdit = null; drawCatalog(); } });
}
