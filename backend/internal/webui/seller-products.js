'use strict';
// Seller Center: produk — daftar, tambah/edit (varian, video/foto, kondisi), terjemahan AI.
addStrings({
  en: {tabProducts: 'Products', addProduct: 'Add product', editProduct: 'Edit product', noProducts: 'No products here yet.',
    f_all: 'All', f_draft: 'Drafts', f_pending: 'In review', f_live: 'Live', f_changes: 'Needs changes', f_archived: 'Archived',
    ps_draft: 'Draft', ps_review: 'In review', ps_live: 'Live', ps_changes: 'Needs changes', ps_rejected: 'Rejected', ps_archived: 'Archived', ps_banned: 'Blocked', ps_hidden: 'Hidden (store not active)',
    stock: 'Stock', variantsN: '{n} variants', backProducts: '← Products', storeNotActiveHint: 'You can prepare products as drafts now. They can be sent for review once your store is approved.',
    secMedia: 'Video & photos', video: 'Product video', videoHint: 'Required. Up to 3 minutes, max 100 MB. Show the real item — buyers watch this first.',
    uploadVideo: 'Upload video', replace: 'Replace', remove: 'Remove', photos: 'Photos', photosHint: 'At least 1, up to 9. The first photo is the cover.', addPhotos: 'Add photos',
    defects: 'Condition photos', defectsHint: 'Close-ups of scratches, wear or missing parts. Required for used items.', moveFirst: 'Make cover',
    uploading: 'Uploading… {p}%', videoTooLong: 'The video is longer than 3 minutes.', fileTooBig: 'The file is too large.',
    secBasic: 'Product details', productName: 'Product name', namePh: 'e.g. Cotton T-shirt, oversized', category: 'Category', chooseCategory: 'Choose…', description: 'Description',
    descPh: 'Material, size, what’s in the box…', translateHint: 'Write in one language only. We translate the name, description, colours and sizes into the other two languages automatically.',
    secCondition: 'Condition', condition: 'Condition', c_new: 'New', c_like_new: 'Like new', c_used_good: 'Used – good', c_used_fair: 'Used – fair', c_refurbished: 'Refurbished',
    conditionNote: 'Describe the condition honestly', conditionNotePh: 'e.g. small scratch on the back, box missing', usageMonths: 'Used for (months)', newOnlyCategory: 'This category is for new items only.',
    secVariants: 'Price, stock & variants', hasOptions: 'This product comes in different colours, sizes…', price: 'Price (USD)', comparePrice: 'Original price', stockQty: 'Stock', sku: 'SKU',
    optional: 'optional', option: 'Option', optKind: 'Type', k_color: 'Colour', k_size: 'Size', k_material: 'Material', k_capacity: 'Capacity', k_custom: 'Other',
    optName: 'Name shown to buyers', addValue: 'Add', valuePh: 'Type and press Enter', addOption: 'Add another option', removeOption: 'Remove option',
    available: 'Sell', applyAll: 'Apply to all', variant: 'Variant', noVariants: 'Add at least one value to each option.',
    secMore: 'Warranty & returns', warranty: 'Warranty', w_none: 'No warranty', w_store: 'Store warranty', w_official: 'Official warranty', warrantyMonths: 'Months', returnDays: 'Return within (days)',
    saveDraft: 'Save draft', submitReview: 'Send for review', saveResubmit: 'Save', reviewNote: 'Changing the name, description, category, condition, video or photos sends the product for review again. Price and stock changes go live right away.',
    savedDraft: 'Draft saved', sentReview: 'Sent for review', archive: 'Hide from shop', unarchive: 'Show again', confirmDelete: 'Delete this product?',
    modPending: 'Phsar is checking this product (usually within a few hours).', modChanges: 'Phsar asked for changes:', modRejected: 'This product can’t be sold on Phsar:', modLive: 'Live — buyers can see this product.',
    secTranslations: 'Translations', translating: 'Translating…', tr_machine: 'AI translation', tr_needs_review: 'AI — please check', tr_seller: 'Edited by you', tr_reviewer: 'Checked by Phsar', tr_source: 'Your original',
    editTr: 'Edit', useAi: 'Use AI translation', translationsLater: 'Translations appear here after you save.',
    err_store_not_active: 'Your store isn’t active yet. Save as draft for now.',
    compressing: 'Compressing video…', compressed: 'optimised', optimizing: 'Optimising…', trs_pending: 'Translating…', trs_retrying: 'Translation problem', trs_failed: 'Translation failed', trs_off: 'Auto-translate off', trOffShort: 'Not translated', trFailedShort: 'Failed', writeTr: 'Write translation', trOff: 'Automatic translation is not set up on this server (ANTHROPIC_API_KEY is empty). You can write the translations yourself, or ask the Phsar team.', trFailed: 'The automatic translation failed. Try again, or write the translation yourself.', trRetrying: 'The translation service had a problem. We’ll retry automatically in a few minutes.', trRetry: 'Try again now', submitLocked: 'Send for review unlocks after Phsar approves your store.', compareHigher: 'Original price must be higher than the price (it’s shown crossed out). Leave it empty if there’s no discount.', needVideo: 'Add one product video.', needPhoto: 'Add at least one photo.', needDefect: 'Used items need a close-up photo of any scratches or wear.', needValue: 'Add at least one value, e.g. Red or M.', needVariantOn: 'Tick at least one variant to sell.', err_product_banned: 'This product was blocked by Phsar.'},
  km: {tabProducts: 'ផលិតផល', addProduct: 'បន្ថែមផលិតផល', editProduct: 'កែផលិតផល', noProducts: 'មិនទាន់មានផលិតផលទេ។',
    f_all: 'ទាំងអស់', f_draft: 'សេចក្តីព្រាង', f_pending: 'កំពុងពិនិត្យ', f_live: 'កំពុងលក់', f_changes: 'ត្រូវកែ', f_archived: 'បានលាក់',
    ps_draft: 'សេចក្តីព្រាង', ps_review: 'កំពុងពិនិត្យ', ps_live: 'កំពុងលក់', ps_changes: 'ត្រូវកែ', ps_rejected: 'បានបដិសេធ', ps_archived: 'បានលាក់', ps_banned: 'បានហាមឃាត់', ps_hidden: 'លាក់ (ហាងមិនទាន់សកម្ម)',
    stock: 'ស្តុក', variantsN: '{n} ជម្រើស', backProducts: '← ផលិតផល', storeNotActiveHint: 'អ្នកអាចរៀបចំផលិតផលជាសេចក្តីព្រាងសិន។ អាចផ្ញើឱ្យពិនិត្យបាន ពេលហាងត្រូវបានអនុម័ត។',
    secMedia: 'វីដេអូ និងរូបថត', video: 'វីដេអូផលិតផល', videoHint: 'ចាំបាច់។ យូរបំផុត 3 នាទី ទំហំ 100 MB។ បង្ហាញទំនិញពិត — អ្នកទិញមើលវាមុនគេ។',
    uploadVideo: 'ផ្ទុកវីដេអូ', replace: 'ប្តូរ', remove: 'ដកចេញ', photos: 'រូបថត', photosHint: 'យ៉ាងហោចណាស់ 1 ច្រើនបំផុត 9។ រូបទីមួយជារូបគម្រប។', addPhotos: 'បន្ថែមរូបថត',
    defects: 'រូបថតស្ថានភាព', defectsHint: 'ថតជិតកន្លែងឆ្កូត ប្រេះ ឬខ្វះ។ ចាំបាច់សម្រាប់ទំនិញប្រើរួច។', moveFirst: 'ធ្វើជារូបគម្រប',
    uploading: 'កំពុងផ្ទុក… {p}%', videoTooLong: 'វីដេអូលើសពី 3 នាទី។', fileTooBig: 'ឯកសារធំពេក។',
    secBasic: 'ព័ត៌មានផលិតផល', productName: 'ឈ្មោះផលិតផល', namePh: 'ឧ. អាវយឺតកប្បាស', category: 'ប្រភេទ', chooseCategory: 'ជ្រើសរើស…', description: 'ការពិពណ៌នា',
    descPh: 'សាច់ក្រណាត់ ទំហំ អ្វីខ្លះក្នុងប្រអប់…', translateHint: 'សរសេរតែមួយភាសា។ យើងបកប្រែឈ្មោះ ការពិពណ៌នា ពណ៌ និងទំហំ ទៅភាសាពីរទៀតដោយស្វ័យប្រវត្តិ។',
    secCondition: 'ស្ថានភាព', condition: 'ស្ថានភាព', c_new: 'ថ្មី', c_like_new: 'ដូចថ្មី', c_used_good: 'ប្រើរួច – ល្អ', c_used_fair: 'ប្រើរួច – មធ្យម', c_refurbished: 'ជួសជុលឡើងវិញ',
    conditionNote: 'ពិពណ៌នាស្ថានភាពដោយស្មោះត្រង់', conditionNotePh: 'ឧ. ឆ្កូតតិចនៅខាងក្រោយ គ្មានប្រអប់', usageMonths: 'បានប្រើ (ខែ)', newOnlyCategory: 'ប្រភេទនេះសម្រាប់តែទំនិញថ្មី។',
    secVariants: 'តម្លៃ ស្តុក និងជម្រើស', hasOptions: 'ផលិតផលនេះមានពណ៌ ទំហំ… ផ្សេងៗ', price: 'តម្លៃ (USD)', comparePrice: 'តម្លៃដើម', stockQty: 'ស្តុក', sku: 'SKU',
    optional: 'មិនចាំបាច់', option: 'ជម្រើស', optKind: 'ប្រភេទ', k_color: 'ពណ៌', k_size: 'ទំហំ', k_material: 'សាច់', k_capacity: 'ចំណុះ', k_custom: 'ផ្សេងទៀត',
    optName: 'ឈ្មោះបង្ហាញអ្នកទិញ', addValue: 'បន្ថែម', valuePh: 'វាយ ហើយចុច Enter', addOption: 'បន្ថែមជម្រើសមួយទៀត', removeOption: 'ដកជម្រើស',
    available: 'លក់', applyAll: 'ដាក់ទាំងអស់', variant: 'ជម្រើស', noVariants: 'បន្ថែមតម្លៃយ៉ាងហោចមួយក្នុងជម្រើសនីមួយៗ។',
    secMore: 'ធានា និងការប្តូរ', warranty: 'ការធានា', w_none: 'គ្មានការធានា', w_store: 'ធានាដោយហាង', w_official: 'ធានាផ្លូវការ', warrantyMonths: 'ខែ', returnDays: 'ប្តូរក្នុងរយៈ (ថ្ងៃ)',
    saveDraft: 'រក្សាទុកសេចក្តីព្រាង', submitReview: 'ផ្ញើឱ្យពិនិត្យ', saveResubmit: 'រក្សាទុក', reviewNote: 'ការកែឈ្មោះ ការពិពណ៌នា ប្រភេទ ស្ថានភាព វីដេអូ ឬរូបថត នឹងត្រូវពិនិត្យម្តងទៀត។ ការកែតម្លៃ និងស្តុកមានប្រសិទ្ធភាពភ្លាម។',
    savedDraft: 'បានរក្សាទុកសេចក្តីព្រាង', sentReview: 'បានផ្ញើឱ្យពិនិត្យ', archive: 'លាក់ពីហាង', unarchive: 'បង្ហាញវិញ', confirmDelete: 'លុបផលិតផលនេះ?',
    modPending: 'Phsar កំពុងពិនិត្យផលិតផលនេះ (ជាធម្មតាក្នុងពីរបីម៉ោង)។', modChanges: 'Phsar ស្នើឱ្យកែ៖', modRejected: 'ផលិតផលនេះមិនអាចលក់លើ Phsar បានទេ៖', modLive: 'កំពុងលក់ — អ្នកទិញមើលឃើញផលិតផលនេះ។',
    secTranslations: 'ការបកប្រែ', translating: 'កំពុងបកប្រែ…', tr_machine: 'បកប្រែដោយ AI', tr_needs_review: 'AI — សូមពិនិត្យ', tr_seller: 'អ្នកបានកែ', tr_reviewer: 'Phsar បានពិនិត្យ', tr_source: 'អត្ថបទដើមរបស់អ្នក',
    editTr: 'កែ', useAi: 'ប្រើការបកប្រែ AI', translationsLater: 'ការបកប្រែនឹងបង្ហាញនៅទីនេះ បន្ទាប់ពីរក្សាទុក។',
    err_store_not_active: 'ហាងមិនទាន់សកម្ម។ សូមរក្សាទុកជាសេចក្តីព្រាងសិន។',
    compressing: 'កំពុងបង្រួមវីដេអូ…', compressed: 'បានបង្រួម', optimizing: 'កំពុងរៀបចំ…', trs_pending: 'កំពុងបកប្រែ…', trs_retrying: 'បញ្ហាបកប្រែ', trs_failed: 'បកប្រែមិនបាន', trs_off: 'ការបកប្រែស្វ័យប្រវត្តិបិទ', trOffShort: 'មិនទាន់បកប្រែ', trFailedShort: 'បរាជ័យ', writeTr: 'សរសេរការបកប្រែ', trOff: 'ម៉ាស៊ីនមេនេះមិនទាន់កំណត់ការបកប្រែស្វ័យប្រវត្តិ។ អ្នកអាចសរសេរការបកប្រែដោយខ្លួនឯង ឬសួរក្រុម Phsar។', trFailed: 'ការបកប្រែស្វ័យប្រវត្តិបរាជ័យ។ សូមព្យាយាមម្តងទៀត ឬសរសេរដោយខ្លួនឯង។', trRetrying: 'សេវាបកប្រែមានបញ្ហា។ យើងនឹងព្យាយាមម្តងទៀតក្នុងពីរបីនាទី។', trRetry: 'ព្យាយាមឥឡូវនេះ', submitLocked: 'អាចផ្ញើឱ្យពិនិត្យបាន ពេល Phsar អនុម័តហាងរបស់អ្នក។', compareHigher: 'តម្លៃដើមត្រូវខ្ពស់ជាងតម្លៃលក់ (បង្ហាញជាតម្លៃគូសចោល)។ ទុកទទេ បើគ្មានការបញ្ចុះតម្លៃ។', needVideo: 'សូមបន្ថែមវីដេអូផលិតផលមួយ។', needPhoto: 'សូមបន្ថែមរូបថតយ៉ាងហោចណាស់មួយ។', needDefect: 'ទំនិញប្រើរួចត្រូវមានរូបថតជិតកន្លែងឆ្កូត ឬខូច។', needValue: 'សូមបន្ថែមយ៉ាងហោចមួយ ឧ. ក្រហម ឬ M។', needVariantOn: 'សូមធីកជម្រើសយ៉ាងហោចមួយដើម្បីលក់។', err_product_banned: 'ផលិតផលនេះត្រូវបាន Phsar ហាមឃាត់។'},
  zh: {tabProducts: '商品', addProduct: '添加商品', editProduct: '编辑商品', noProducts: '这里还没有商品。',
    f_all: '全部', f_draft: '草稿', f_pending: '审核中', f_live: '在售', f_changes: '需修改', f_archived: '已下架',
    ps_draft: '草稿', ps_review: '审核中', ps_live: '在售', ps_changes: '需修改', ps_rejected: '已拒绝', ps_archived: '已下架', ps_banned: '已封禁', ps_hidden: '未展示（店铺未上线）',
    stock: '库存', variantsN: '{n} 个规格', backProducts: '← 商品', storeNotActiveHint: '现在可以先保存为草稿，店铺审核通过后再提交审核。',
    secMedia: '视频和照片', video: '商品视频', videoHint: '必填。最长 3 分钟，最大 100 MB。请拍摄实物——买家会先看视频。',
    uploadVideo: '上传视频', replace: '更换', remove: '删除', photos: '照片', photosHint: '至少 1 张，最多 9 张。第一张为封面。', addPhotos: '添加照片',
    defects: '成色照片', defectsHint: '近拍划痕、磨损或缺件处。二手商品必填。', moveFirst: '设为封面',
    uploading: '上传中… {p}%', videoTooLong: '视频超过 3 分钟。', fileTooBig: '文件太大。',
    secBasic: '商品信息', productName: '商品名称', namePh: '例如：纯棉宽松T恤', category: '类目', chooseCategory: '请选择…', description: '商品描述',
    descPh: '材质、尺寸、包装内容…', translateHint: '只需用一种语言填写。名称、描述、颜色和尺码会自动翻译成另外两种语言。',
    secCondition: '成色', condition: '成色', c_new: '全新', c_like_new: '几乎全新', c_used_good: '二手 – 良好', c_used_fair: '二手 – 一般', c_refurbished: '翻新',
    conditionNote: '如实描述成色', conditionNotePh: '例如：背面有轻微划痕，无包装盒', usageMonths: '已使用（月）', newOnlyCategory: '该类目只允许全新商品。',
    secVariants: '价格、库存和规格', hasOptions: '该商品有不同颜色、尺码…', price: '价格（美元）', comparePrice: '原价', stockQty: '库存', sku: 'SKU',
    optional: '选填', option: '规格', optKind: '类型', k_color: '颜色', k_size: '尺码', k_material: '材质', k_capacity: '容量', k_custom: '其他',
    optName: '向买家显示的名称', addValue: '添加', valuePh: '输入后按回车', addOption: '再加一个规格', removeOption: '删除规格',
    available: '在售', applyAll: '全部应用', variant: '规格', noVariants: '每个规格至少添加一个值。',
    secMore: '保修和退换', warranty: '保修', w_none: '无保修', w_store: '店铺保修', w_official: '官方保修', warrantyMonths: '月', returnDays: '可退换天数',
    saveDraft: '保存草稿', submitReview: '提交审核', saveResubmit: '保存', reviewNote: '修改名称、描述、类目、成色、视频或照片后需重新审核。修改价格和库存会立即生效。',
    savedDraft: '草稿已保存', sentReview: '已提交审核', archive: '从店铺下架', unarchive: '重新上架', confirmDelete: '删除此商品？',
    modPending: 'Phsar 正在审核此商品（通常几小时内）。', modChanges: 'Phsar 要求修改：', modRejected: '此商品不能在 Phsar 销售：', modLive: '在售——买家可以看到此商品。',
    secTranslations: '翻译', translating: '翻译中…', tr_machine: 'AI 翻译', tr_needs_review: 'AI——请检查', tr_seller: '已由你修改', tr_reviewer: 'Phsar 已核对', tr_source: '你的原文',
    editTr: '编辑', useAi: '使用 AI 翻译', translationsLater: '保存后这里会显示翻译。',
    err_store_not_active: '店铺尚未上线，请先保存为草稿。',
    compressing: '正在压缩视频…', compressed: '已优化', optimizing: '处理中…', trs_pending: '翻译中…', trs_retrying: '翻译出现问题', trs_failed: '翻译失败', trs_off: '自动翻译未开启', trOffShort: '未翻译', trFailedShort: '失败', writeTr: '自己填写翻译', trOff: '此服务器未设置自动翻译（ANTHROPIC_API_KEY 为空）。你可以自己填写翻译，或联系 Phsar 团队。', trFailed: '自动翻译失败。请重试，或自己填写翻译。', trRetrying: '翻译服务出现问题，几分钟后会自动重试。', trRetry: '立即重试', submitLocked: 'Phsar 通过店铺审核后才能提交商品审核。', compareHigher: '原价必须高于售价（会以划线价显示）。没有折扣请留空。', needVideo: '请添加一个商品视频。', needPhoto: '请至少添加一张照片。', needDefect: '二手商品需要划痕或磨损处的近拍照片。', needValue: '请至少添加一个值，例如 红色 或 M。', needVariantOn: '请至少勾选一个要销售的规格。', err_product_banned: '此商品已被 Phsar 封禁。'},
});

const LANGS = {en: 'English', km: 'ខ្មែរ', zh: '中文'};
const money = v => v == null ? '—' : '$' + Number(v).toFixed(2);
const PF = {filter: 'all', cats: {}, E: null};

function productState(p, storeActive) {
  if (p.status === 'draft') return ['', 'ps_draft'];
  if (p.status === 'banned') return ['bad', 'ps_banned'];
  if (p.status === 'archived') return ['', 'ps_archived'];
  if (p.moderation_status === 'pending_review') return ['warn', 'ps_review'];
  if (p.moderation_status === 'changes_requested') return ['bad', 'ps_changes'];
  if (p.moderation_status === 'rejected') return ['bad', 'ps_rejected'];
  return storeActive ? ['ok', 'ps_live'] : ['warn', 'ps_hidden'];
}

async function categories() {
  if (!PF.cats[P.lang]) PF.cats[P.lang] = (await api('/v1/app/categories')).categories;
  return PF.cats[P.lang];
}

// ---- Tab daftar produk (di halaman toko) ----
async function tabProducts(id) {
  const [{products}, cats] = await Promise.all([api(`/v1/seller/stores/${id}/products?filter=${PF.filter}`), categories()]);
  const catName = c => (cats.find(x => x.id === c) || {}).name || c;
  const active = S.store.status === 'active';
  const cards = products.map(p => {
    const [cls, key] = productState(p, active);
    const price = p.min_price === p.max_price ? money(p.min_price) : `${money(p.min_price)} – ${money(p.max_price)}`;
    return `<div class="card pcard"><a class="plink" href="#/store/${id}/product/${p.id}">
      <div class="pcover">${p.cover_url ? `<img src="${esc(p.cover_url)}" alt="" loading="lazy">` : ''}${p.has_video ? '<span class="vbadge" aria-hidden="true">▶</span>' : ''}</div>
      <div class="spread"><span class="chip ${cls}">${t(key)}</span><span class="muted">${esc(catName(p.category_id))}</span></div>
      <b class="pname">${esc(p.name)}</b>
      <div class="spread"><b>${price}</b><span class="muted">${t('stock')} ${p.total_stock}${p.variants > 1 ? ' · ' + t('variantsN', {n: p.variants}) : ''}</span></div>
      ${p.moderation_note && ['changes_requested', 'rejected'].includes(p.moderation_status) ? `<p class="muted" style="color:var(--bad)">${esc(p.moderation_note)}</p>` : ''}
    </a>${p.status !== 'banned' ? `<button type="button" class="small" onclick="openStockDialog('${id}','${p.id}')">${t('stockPrice')}</button>` : ''}</div>`;
  }).join('');
  $('#tab').innerHTML = `
    ${active ? '' : `<div class="banner info">${t('storeNotActiveHint')}</div>`}
    <div class="spread"><div class="tabs">${['all', 'draft', 'pending', 'live', 'changes', 'archived'].map(f =>
      `<button type="button" aria-selected="${PF.filter === f}" onclick="PF.filter='${f}';render()">${t('f_' + f)}</button>`).join('')}</div>
      <a class="btn primary" href="#/store/${id}/product/new">+ ${t('addProduct')}</a></div>
    ${products.length ? `<div class="pgrid">${cards}</div>` : `<p class="muted">${t('noProducts')}</p>`}`;
}

// ---- Editor produk ----
function emptyProduct() {
  return {source_lang: P.lang, name: '', description: '', category_id: '', condition: 'new', condition_note: '', usage_months: null,
    warranty_type: 'none', warranty_months: 0, return_days: 7, media: [], hasOptions: false,
    options: [{kind: 'color', name: '', values: []}], single: {price: '', compare_at_price: '', stock_qty: 1, sku: ''}, rows: {}};
}

function fromProduct(p) {
  const L = p.source_lang;
  const e = Object.assign(emptyProduct(), {
    id: p.id, source_lang: L, name: p.translations[L]?.name || '', description: p.translations[L]?.description || '', category_id: p.category_id,
    condition: p.condition, condition_note: p.condition_note, usage_months: p.usage_months, warranty_type: p.warranty_type,
    warranty_months: p.warranty_months, return_days: p.return_days,
    media: p.media.map(m => ({key: m.key, url: m.url, kind: m.kind, poster: m.poster_url, processing: m.processing, duration: m.duration_sec})), hasOptions: p.options.length > 0,
  });
  if (p.options.length) {
    e.options = p.options.map(o => ({kind: o.kind, name: o.name[L] || '', values: o.values.map(v => ({value: v.value[L] || '', swatch_hex: v.swatch_hex || ''}))}));
    for (const v of p.variants) e.rows[rowKey(v.options)] = {on: true, price: v.price, compare_at_price: v.compare_at_price || '', stock_qty: v.stock_qty, sku: v.sku};
  } else if (p.variants[0]) {
    const v = p.variants[0];
    e.single = {price: v.price, compare_at_price: v.compare_at_price || '', stock_qty: v.stock_qty, sku: v.sku};
  }
  return e;
}

const rowKey = vals => vals.map(v => v.toLowerCase()).join('\u001f');
function combos(options) {
  return options.reduce((acc, o) => acc.flatMap(a => o.values.map(v => [...a, v.value])), [[]]);
}

async function renderProductEditor(storeId, productId) {
  const [st, cats] = await Promise.all([api(`/v1/seller/stores/${storeId}`), categories()]);
  S.store = st;
  if (!PF.E || PF.E.storeId !== storeId || PF.E.loadedId !== productId) {
    const p = productId === 'new' ? null : await api(`/v1/seller/stores/${storeId}/products/${productId}`);
    PF.E = {storeId, loadedId: productId, product: p, f: p ? fromProduct(p) : emptyProduct(), editingTr: null};
  }
  drawEditor(cats);
}

function drawEditor(cats) {
  const E = PF.E, f = E.f, p = E.product, storeId = E.storeId;
  const isNew = !p, wasSent = p && p.status !== 'draft';
  const cat = cats.find(c => c.id === f.category_id);
  const media = kind => f.media.map((m, i) => ({...m, i})).filter(m => m.kind === kind);
  const video = media('video')[0], photos = media('photo'), defects = media('defect_photo');
  const thumb = (m, first) => `<figure class="thumb"><img src="${esc(m.url)}" alt="">${m.processing === 'processing' ? `<span class="procbadge">${t('optimizing')}</span>` : ''}
      <figcaption>${first ? '' : `<button type="button" onclick="moveMedia(${m.i})" title="${t('moveFirst')}">★</button>`}<button type="button" onclick="removeMedia(${m.i})" aria-label="${t('remove')}">✕</button></figcaption></figure>`;

  let banner = '';
  if (p) {
    const [cls, key] = productState(p, st_active());
    if (p.status === 'banned') banner = `<div class="banner bad">${t('err_product_banned')}</div>`;
    else if (p.status !== 'draft' && p.moderation_status === 'pending_review') banner = `<div class="banner warn">${t('modPending')}</div>`;
    else if (p.moderation_status === 'changes_requested') banner = `<div class="banner bad"><b>${t('modChanges')}</b> ${esc(p.moderation_note || '')}</div>`;
    else if (p.moderation_status === 'rejected') banner = `<div class="banner bad"><b>${t('modRejected')}</b> ${esc(p.moderation_note || '')}</div>`;
    else if (p.visible) banner = `<div class="banner ok">${t('modLive')}</div>`;
    banner = `<div class="row"><span class="chip ${cls}">${t(key)}</span></div>` + banner;
  }

  const optionBlocks = f.options.map((o, oi) => `<div class="card optcard">
      <div class="grid">
        <label>${t('optKind')}<select data-opt="${oi}" data-field="kind">${['color', 'size', 'material', 'capacity', 'custom'].map(k =>
          `<option value="${k}" ${o.kind === k ? 'selected' : ''}>${t('k_' + k)}</option>`).join('')}</select></label>
        <label>${t('optName')}<input data-opt="${oi}" data-field="name" value="${esc(o.name)}" maxlength="30" placeholder="${t('k_' + o.kind)}"></label>
      </div>
      <div class="row">${o.values.map((v, vi) => `<span class="vchip">${o.kind === 'color' ? `<input type="color" value="${v.swatch_hex || '#cccccc'}" data-opt="${oi}" data-val="${vi}" aria-label="${esc(v.value)}">` : ''}${esc(v.value)}
          <button type="button" onclick="removeValue(${oi},${vi})" aria-label="${t('remove')}">✕</button></span>`).join('')}
        <input class="vinput" data-addval="${oi}" placeholder="${t('valuePh')}" maxlength="40"><button type="button" onclick="addValue(${oi})">${t('addValue')}</button></div>
      ${f.options.length > 1 ? `<div><button type="button" class="danger" onclick="removeOption(${oi})">${t('removeOption')}</button></div>` : ''}
    </div>`).join('');

  const allCombos = f.hasOptions ? combos(f.options) : [];
  const ready = f.hasOptions && f.options.every(o => o.values.length);
  const variantTable = !f.hasOptions ? `<div class="grid">
      <label>${t('price')}<input id="sPrice" inputmode="decimal" value="${esc(f.single.price)}" placeholder="0.00"></label>
      <label>${t('comparePrice')} <span class="muted">(${t('optional')})</span><input id="sCompare" inputmode="decimal" value="${esc(f.single.compare_at_price)}"></label>
      <label>${t('stockQty')}<input id="sStock" type="number" min="0" value="${esc(f.single.stock_qty)}"></label>
      <label>${t('sku')} <span class="muted">(${t('optional')})</span><input id="sSku" value="${esc(f.single.sku)}" maxlength="64"></label></div>`
    : !ready ? `<p class="muted">${t('noVariants')}</p>` : `<div style="overflow-x:auto"><table class="vtable"><thead><tr><th>${t('available')}</th><th>${t('variant')}</th><th>${t('price')}</th><th>${t('comparePrice')}</th><th>${t('stockQty')}</th><th>${t('sku')}</th></tr></thead><tbody>
      <tr class="applyall"><td></td><td><b>${t('applyAll')}</b></td><td><input id="allPrice" inputmode="decimal" placeholder="0.00"></td><td><input id="allCompare" inputmode="decimal"></td><td><input id="allStock" type="number" min="0"></td><td><button type="button" onclick="applyAll()">${t('applyAll')}</button></td></tr>
      ${allCombos.map(c => { const k = rowKey(c), r = f.rows[k] || (f.rows[k] = {on: true, price: '', compare_at_price: '', stock_qty: 1, sku: ''});
        return `<tr data-row="${esc(k)}"><td><input type="checkbox" data-rf="on" ${r.on ? 'checked' : ''} aria-label="${t('available')} ${esc(c.join(' / '))}"></td><td>${esc(c.join(' / '))}</td>
          <td><input data-rf="price" inputmode="decimal" value="${esc(r.price)}" placeholder="0.00"></td><td><input data-rf="compare_at_price" inputmode="decimal" value="${esc(r.compare_at_price)}"></td>
          <td><input data-rf="stock_qty" type="number" min="0" value="${esc(r.stock_qty)}"></td><td><input data-rf="sku" value="${esc(r.sku)}" maxlength="64"></td></tr>`; }).join('')}
    </tbody></table></div>`;

  const others = p ? Object.keys(LANGS).filter(l => l !== p.source_lang) : [];
  const st = p?.translation_state || 'done';
  const trBlock = !p ? `<p class="muted">${t('translationsLater')}</p>` : others.map(l => {
    const tr = p.translations[l];
    const missing = st === 'off' ? t('trOffShort') : st === 'failed' ? t('trFailedShort') : t('translating');
    const tag = !tr ? missing : tr.seller_edited ? t('tr_seller') : tr.source === 'reviewer' ? t('tr_reviewer') : t('tr_' + tr.review_status);
    if (E.editingTr === l) return `<div class="card trcard"><b>${LANGS[l]}</b>
      <label>${t('productName')}<input id="trName" value="${esc(tr?.name || '')}" maxlength="160"></label>
      <label>${t('description')}<textarea id="trDesc">${esc(tr?.description || '')}</textarea></label>
      <div class="row"><button type="button" class="primary" onclick="saveTr('${l}')">${t('save')}</button><button type="button" onclick="PF.E.editingTr=null;drawEditor(PF.cats[P.lang])">${t('cancel')}</button></div></div>`;
    return `<div class="card trcard"><div class="spread"><b>${LANGS[l]}</b><span class="chip ${tr?.review_status === 'needs_review' || (!tr && st !== 'pending') ? 'warn' : ''}">${tag}</span></div>
      ${tr ? `<p><b>${esc(tr.name)}</b></p><p class="muted pre">${esc(tr.description)}</p>` : `<p class="muted">${missing}</p>`}
      <div class="row">${tr || st !== 'pending' ? `<button type="button" onclick="PF.E.editingTr='${l}';drawEditor(PF.cats[P.lang])">${tr ? t('editTr') : t('writeTr')}</button>` : ''}
        ${tr?.seller_edited ? `<button type="button" onclick="resetTr('${l}')">${t('useAi')}</button>` : ''}</div></div>`;
  }).join('');

  $('#app').innerHTML = shell(`
    <a href="#/store/${storeId}">${t('backProducts')}</a>
    <h1>${isNew ? t('addProduct') : t('editProduct')}</h1>
    ${banner}
    ${st_active() ? '' : `<div class="banner info">${t('storeNotActiveHint')}</div>`}
    <section class="card" id="secMedia"><h2>${t('secMedia')}</h2>
      <div><b>${t('video')}</b><p class="muted">${t('videoHint')}</p></div>
      ${video ? `<div class="row"><div class="vwrap"><video class="vprev" src="${esc(video.url)}" ${video.poster ? `poster="${esc(video.poster)}"` : ''} controls preload="metadata"></video>
        ${video.processing === 'processing' ? `<span class="procbadge">${t('compressing')}</span>` : video.duration ? `<span class="procbadge ok">${video.duration}s · ${t('compressed')}</span>` : ''}</div>
          <div class="row"><label class="btn">${t('replace')}<input type="file" class="hidden" accept="video/mp4,video/quicktime,video/webm" data-upload="video"></label>
          <button type="button" class="danger" onclick="removeMedia(${video.i})">${t('remove')}</button></div></div>`
        : `<label class="btn drop" id="videoBox">▶ ${t('uploadVideo')}<input type="file" class="hidden" accept="video/mp4,video/quicktime,video/webm" data-upload="video"></label>`}
      <div><b>${t('photos')}</b><p class="muted">${t('photosHint')}</p></div>
      <div class="thumbs" id="photoBox">${photos.map((m, n) => thumb(m, n === 0)).join('')}
        ${photos.length < 9 ? `<label class="btn drop small">+ ${t('addPhotos')}<input type="file" class="hidden" accept="image/jpeg,image/png,image/webp" multiple data-upload="photo"></label>` : ''}</div>
      ${f.condition !== 'new' ? `<div><b>${t('defects')}</b><p class="muted">${t('defectsHint')}</p></div>
        <div class="thumbs" id="defectBox">${defects.map(m => thumb(m, true)).join('')}
        ${defects.length < 6 ? `<label class="btn drop small">+ ${t('addPhotos')}<input type="file" class="hidden" accept="image/jpeg,image/png,image/webp" multiple data-upload="defect_photo"></label>` : ''}</div>` : ''}
      <p class="muted hidden" id="upState" role="status"></p>
    </section>

    <section class="card"><h2>${t('secBasic')}</h2>
      <div class="grid"><label>${t('writtenIn')}<select id="fLang">${Object.entries(LANGS).map(([c, n]) => `<option value="${c}" ${f.source_lang === c ? 'selected' : ''}>${n}</option>`).join('')}</select></label>
        <label>${t('category')}<select id="fCat"><option value="">${t('chooseCategory')}</option>${cats.map(c => `<option value="${c.id}" ${f.category_id === c.id ? 'selected' : ''}>${esc(c.name)}</option>`).join('')}</select></label></div>
      <label>${t('productName')}<input id="fName" value="${esc(f.name)}" maxlength="160" placeholder="${t('namePh')}"></label>
      <label>${t('description')}<textarea id="fDesc" maxlength="5000" placeholder="${t('descPh')}">${esc(f.description)}</textarea></label>
      <p class="muted">${t('translateHint')}</p>
    </section>

    <section class="card"><h2>${t('secCondition')}</h2>
      <label>${t('condition')}<select id="fCond">${['new', 'like_new', 'used_good', 'used_fair', 'refurbished'].map(c =>
        `<option value="${c}" ${f.condition === c ? 'selected' : ''} ${c !== 'new' && cat && !cat.allows_used ? 'disabled' : ''}>${t('c_' + c)}</option>`).join('')}</select></label>
      ${cat && !cat.allows_used ? `<p class="muted">${t('newOnlyCategory')}</p>` : ''}
      ${f.condition !== 'new' ? `<label>${t('conditionNote')}<textarea id="fNote" maxlength="1000" placeholder="${t('conditionNotePh')}">${esc(f.condition_note)}</textarea></label>
        <label style="max-width:240px">${t('usageMonths')} <span class="muted">(${t('optional')})</span><input id="fUsage" type="number" min="0" max="600" value="${esc(f.usage_months ?? '')}"></label>` : ''}
    </section>

    <section class="card" id="secVariants"><h2>${t('secVariants')}</h2>
      <label class="inline"><input type="checkbox" id="fHasOpt" ${f.hasOptions ? 'checked' : ''}> ${t('hasOptions')}</label>
      ${f.hasOptions ? optionBlocks + (f.options.length < 2 ? `<div><button type="button" onclick="addOption()">+ ${t('addOption')}</button></div>` : '') : ''}
      ${variantTable}
    </section>

    <section class="card"><h2>${t('secMore')}</h2>
      <div class="grid"><label>${t('warranty')}<select id="fWarranty">${['none', 'store', 'official'].map(w => `<option value="${w}" ${f.warranty_type === w ? 'selected' : ''}>${t('w_' + w)}</option>`).join('')}</select></label>
        ${f.warranty_type !== 'none' ? `<label>${t('warrantyMonths')}<input id="fWMonths" type="number" min="1" max="120" value="${esc(f.warranty_months || 12)}"></label>` : ''}
        <label>${t('returnDays')}<input id="fReturn" type="number" min="0" max="30" value="${esc(f.return_days)}"></label></div>
    </section>

    ${wasSent ? `<p class="muted">${t('reviewNote')}</p>` : ''}
    <div class="row sticky-actions" data-errors-anchor>
      ${!wasSent ? `<button type="button" id="btnDraft">${t('saveDraft')}</button>` : ''}
      <button type="button" class="primary" id="btnSubmit" ${p?.status === 'banned' || (!wasSent && !st_active()) ? 'disabled' : ''}>${wasSent ? t('saveResubmit') : t('submitReview')}</button>
      ${!wasSent && !st_active() ? `<span class="muted">${t('submitLocked')}</span>` : ''}
      ${p && p.status !== 'draft' && p.status !== 'banned' ? `<button type="button" onclick="setPStatus('${p.status === 'archived' ? 'active' : 'archived'}')">${p.status === 'archived' ? t('unarchive') : t('archive')}</button>` : ''}
      ${p ? `<button type="button" class="danger" onclick="deleteProduct()">${t('del')}</button>` : ''}
    </div>

    <section class="card"><div class="spread"><h2>${t('secTranslations')}</h2>${p && st !== 'done' ? `<span class="chip ${st === 'pending' ? 'warn' : 'bad'}">${t('trs_' + st)}</span>` : ''}</div>
      ${p ? `<p class="muted">${t('tr_source')}: ${LANGS[p.source_lang]}</p>` : ''}
      ${p && st === 'off' ? `<div class="banner warn">${t('trOff')}</div>` : ''}
      ${p && (st === 'failed' || st === 'retrying') ? `<div class="banner bad"><b>${st === 'failed' ? t('trFailed') : t('trRetrying')}</b>
        ${p.translation_error ? `<br><span class="muted" style="color:inherit">${esc(p.translation_error)}</span>` : ''}
        <div style="margin-top:8px"><button type="button" onclick="retryTr()">${t('trRetry')}</button></div></div>` : ''}
      <div class="grid">${trBlock}</div></section>`);
  bindEditor();
  clearTimeout(PF.poll);
  const mediaBusy = p && p.media.some(m => m.processing === 'processing');
  if (p && (st === 'pending' || st === 'retrying' || mediaBusy)) PF.poll = setTimeout(pollTranslations, st === 'pending' || mediaBusy ? 4000 : 15000);
}

const st_active = () => S.store?.status === 'active';

// Simpan isi form ke PF.E.f sebelum menggambar ulang.
function readForm() {
  const f = PF.E.f, v = id => $('#' + id)?.value;
  if ($('#fName')) Object.assign(f, {name: v('fName'), description: v('fDesc'), source_lang: v('fLang'), category_id: v('fCat'), condition: v('fCond')});
  if ($('#fNote')) f.condition_note = v('fNote');
  if ($('#fUsage')) f.usage_months = v('fUsage') === '' ? null : Number(v('fUsage'));
  if ($('#fWarranty')) f.warranty_type = v('fWarranty');
  if ($('#fWMonths')) f.warranty_months = Number(v('fWMonths'));
  if ($('#fReturn')) f.return_days = Number(v('fReturn'));
  if ($('#sPrice')) f.single = {price: v('sPrice'), compare_at_price: v('sCompare'), stock_qty: Number(v('sStock')), sku: v('sSku')};
  document.querySelectorAll('[data-opt][data-field]').forEach(el => { f.options[el.dataset.opt][el.dataset.field] = el.value; });
  document.querySelectorAll('[data-opt][data-val]').forEach(el => { const v = f.options[el.dataset.opt].values[el.dataset.val]; if (v.swatch_hex || el.value !== '#cccccc') v.swatch_hex = el.value; });
  document.querySelectorAll('tr[data-row]').forEach(tr => {
    const r = f.rows[tr.dataset.row];
    tr.querySelectorAll('[data-rf]').forEach(el => { r[el.dataset.rf] = el.type === 'checkbox' ? el.checked : el.dataset.rf === 'stock_qty' ? Number(el.value) : el.value; });
  });
}
const redraw = () => { readForm(); drawEditor(PF.cats[P.lang]); };

function bindEditor() {
  ['fCond', 'fCat', 'fWarranty'].forEach(id => $('#' + id) && ($('#' + id).onchange = redraw));
  $('#fHasOpt').onchange = e => { readForm(); PF.E.f.hasOptions = e.target.checked; drawEditor(PF.cats[P.lang]); };
  document.querySelectorAll('[data-opt][data-field=kind]').forEach(el => el.onchange = redraw);
  document.querySelectorAll('[data-addval]').forEach(el => el.onkeydown = e => { if (e.key === 'Enter') { e.preventDefault(); addValue(Number(el.dataset.addval)); } });
  document.querySelectorAll('[data-upload]').forEach(el => el.onchange = e => uploadMedia(el.dataset.upload, [...e.target.files]));
  if ($('#btnDraft')) $('#btnDraft').onclick = () => saveProduct(false);
  $('#btnSubmit').onclick = () => saveProduct(true);
}

function addValue(oi) {
  readForm();
  const inp = document.querySelector(`[data-addval="${oi}"]`), o = PF.E.f.options[oi];
  for (const raw of inp.value.split(',')) {
    const val = raw.trim();
    if (val && !o.values.some(v => v.value.toLowerCase() === val.toLowerCase())) o.values.push({value: val, swatch_hex: ''});
  }
  drawEditor(PF.cats[P.lang]);
  document.querySelector(`[data-addval="${oi}"]`)?.focus();
}
function removeValue(oi, vi) { readForm(); PF.E.f.options[oi].values.splice(vi, 1); drawEditor(PF.cats[P.lang]); }
function addOption() { readForm(); PF.E.f.options.push({kind: PF.E.f.options[0]?.kind === 'size' ? 'color' : 'size', name: '', values: []}); drawEditor(PF.cats[P.lang]); }
function removeOption(oi) { readForm(); PF.E.f.options.splice(oi, 1); drawEditor(PF.cats[P.lang]); }
function applyAll() {
  readForm();
  const price = $('#allPrice').value, cmp = $('#allCompare').value, stock = $('#allStock').value;
  for (const r of Object.values(PF.E.f.rows)) {
    if (price) r.price = price;
    if (cmp) r.compare_at_price = cmp;
    if (stock !== '') r.stock_qty = Number(stock);
  }
  drawEditor(PF.cats[P.lang]);
}
function removeMedia(i) { readForm(); PF.E.f.media.splice(i, 1); drawEditor(PF.cats[P.lang]); }
function moveMedia(i) {
  readForm();
  const m = PF.E.f.media.splice(i, 1)[0], first = PF.E.f.media.findIndex(x => x.kind === m.kind);
  PF.E.f.media.splice(first < 0 ? PF.E.f.media.length : first, 0, m);
  drawEditor(PF.cats[P.lang]);
}

function videoDuration(file) {
  return new Promise(res => {
    const v = document.createElement('video');
    v.preload = 'metadata';
    v.onloadedmetadata = () => { URL.revokeObjectURL(v.src); res(v.duration); };
    v.onerror = () => res(0); // format tidak bisa dibaca browser (mis. HEVC) → biar server yang menilai
    v.src = URL.createObjectURL(file);
  });
}

// Upload dengan progress (video bisa besar). Token diperbarui sekali kalau kedaluwarsa.
function xhrUpload(file, onProgress, retried) {
  return new Promise((resolve, reject) => {
    const fd = new FormData(); fd.append('kind', 'product_media'); fd.append('file', file);
    const x = new XMLHttpRequest();
    x.open('POST', '/v1/seller/uploads');
    x.setRequestHeader('Authorization', 'Bearer ' + P.token);
    x.setRequestHeader('X-Lang', P.lang);
    x.upload.onprogress = e => e.lengthComputable && onProgress(Math.round(e.loaded / e.total * 100));
    x.onload = async () => {
      let data = {}; try { data = JSON.parse(x.responseText); } catch {}
      if (x.status === 401 && !retried && await refreshSession()) return resolve(xhrUpload(file, onProgress, true));
      if (x.status >= 200 && x.status < 300) resolve(data);
      else reject(Object.assign(new Error(data.error?.message || x.statusText), {code: data.error?.code}));
    };
    x.onerror = () => reject(new Error('network'));
    x.send(fd);
  });
}

async function uploadMedia(kind, files) {
  readForm();
  const state = $('#upState');
  for (const file of files) {
    if (kind === 'video') {
      if (file.size > 100 * 1024 * 1024) { toast(t('fileTooBig')); continue; }
      if (await videoDuration(file) > 181) { toast(t('videoTooLong')); continue; }
    } else if (file.size > 8 * 1024 * 1024) { toast(t('fileTooBig')); continue; }
    if (kind === 'photo' && PF.E.f.media.filter(m => m.kind === 'photo').length >= 9) break;
    state.classList.remove('hidden');
    try {
      const up = await xhrUpload(file, pct => { state.textContent = t('uploading', {p: pct}); });
      if (kind === 'video') PF.E.f.media = PF.E.f.media.filter(m => m.kind !== 'video');
      const item = {key: up.key, url: up.url, kind};
      if (kind === 'video') PF.E.f.media.unshift(item); else PF.E.f.media.push(item);
    } catch (e) { toast(errMsg(e)); }
  }
  drawEditor(PF.cats[P.lang]);
}

function buildBody(submit) {
  const f = PF.E.f;
  const body = {category_id: f.category_id, source_lang: f.source_lang, name: f.name, description: f.description, condition: f.condition,
    condition_note: f.condition === 'new' ? '' : f.condition_note, usage_months: f.condition === 'new' ? null : f.usage_months,
    warranty_type: f.warranty_type, warranty_months: f.warranty_type === 'none' ? 0 : f.warranty_months, return_days: f.return_days,
    media: f.media.filter(m => f.condition !== 'new' || m.kind !== 'defect_photo').map(m => ({key: m.key, kind: m.kind})), submit,
    options: [], variants: []};
  PF.E.sentRows = [];
  if (f.hasOptions) {
    body.options = f.options.map(o => ({kind: o.kind, name: o.name.trim() || t('k_' + o.kind), values: o.values.map(v => ({value: v.value, swatch_hex: o.kind === 'color' ? (v.swatch_hex || '') : ''}))}));
    for (const c of combos(f.options)) {
      const r = f.rows[rowKey(c)];
      if (r && r.on) PF.E.sentRows.push(rowKey(c));
      if (r && r.on) body.variants.push({options: c, price: String(r.price).trim(), compare_at_price: String(r.compare_at_price || '').trim(), stock_qty: Number(r.stock_qty) || 0, sku: (r.sku || '').trim()});
    }
  } else {
    body.variants = [{options: [], price: String(f.single.price).trim(), compare_at_price: String(f.single.compare_at_price || '').trim(), stock_qty: Number(f.single.stock_qty) || 0, sku: (f.single.sku || '').trim()}];
  }
  return body;
}

// Cek di browser dulu supaya seller langsung diarahkan ke isian yang kurang.
function validateProductForm(submit) {
  const f = PF.E.f, kinds = f.media.map(m => m.kind);
  const price = el => Number(String(el.value).trim()) > 0;
  const rules = [['#fName', t('minChars', {n: 3}), el => el.value.trim().length >= 3], ['#fCat', t('pickOne')]];
  if (f.condition !== 'new') rules.push(['#fNote']);
  // Harga coret (original price) opsional, tapi kalau diisi harus lebih tinggi dari harga jual.
  const compareOk = (cmpEl, priceEl) => { const c = String(cmpEl.value).trim(); return c === '' || Number(c) > Number(priceEl.value); };
  if (!f.hasOptions) {
    rules.push(['#sPrice', t('positive'), price]);
    rules.push(['#sCompare', t('compareHigher'), el => compareOk(el, $('#sPrice'))]);
  }
  else {
    f.options.forEach((o, i) => rules.push([`[data-addval="${i}"]`, t('needValue'), () => o.values.length > 0, `${t('option')} ${i + 1} (${t('k_' + o.kind)})`]));
    rules.push(['#secVariants', t('needVariantOn'), () => !document.querySelector('tr[data-row]') || !!document.querySelector('tr[data-row] [data-rf=on]:checked'), t('secVariants')]);
    document.querySelectorAll('tr[data-row]').forEach(tr => {
      if (!tr.querySelector('[data-rf=on]').checked) return;
      rules.push([tr.querySelector('[data-rf=price]'), t('positive'), price]);
      rules.push([tr.querySelector('[data-rf=compare_at_price]'), t('compareHigher'), el => compareOk(el, tr.querySelector('[data-rf=price]'))]);
    });
  }
  if (submit) {
    rules.unshift(['#videoBox, video.vprev', t('needVideo'), () => kinds.includes('video'), t('video')]);
    rules.splice(1, 0, ['#photoBox', t('needPhoto'), () => kinds.includes('photo'), t('photos')]);
    if (['used_good', 'used_fair'].includes(f.condition)) rules.push(['#defectBox', t('needDefect'), () => kinds.includes('defect_photo'), t('defects')]);
  }
  return checkRequired(rules); // selector atau elemen
}

function productErrorField(field) {
  const simple = {name: '#fName', description: '#fDesc', category_id: '#fCat', source_lang: '#fLang', condition: '#fCond',
    condition_note: '#fNote', usage_months: '#fUsage', warranty_type: '#fWarranty', warranty_months: '#fWMonths', return_days: '#fReturn',
    options: '#secVariants', variants: '#secVariants', media: '#secMedia'};
  if (simple[field]) return simple[field];
  let m = field.match(/^options\[(\d+)\]\.(\w+)/);
  if (m) return m[2] === 'values' ? `[data-addval="${m[1]}"]` : `[data-opt="${m[1]}"][data-field="${m[2]}"]`;
  m = field.match(/^options\[(\d+)\]/);
  if (m) return `[data-addval="${m[1]}"]`;
  m = field.match(/^variants\[(\d+)\]\.(\w+)/);
  if (m) {
    const col = m[2] === 'options' ? 'on' : m[2];
    if (!PF.E.f.hasOptions) return ({price: '#sPrice', compare_at_price: '#sCompare', stock_qty: '#sStock', sku: '#sSku'})[col];
    const key = PF.E.sentRows?.[Number(m[1])];
    return key ? document.querySelector(`tr[data-row="${CSS.escape(key)}"] [data-rf="${col}"]`) : '#secVariants';
  }
  if (/^media\[(\d+)\]/.test(field)) return '#secMedia';
  return null;
}

async function saveProduct(submit) {
  readForm();
  if (!validateProductForm(submit)) return;
  const E = PF.E, url = `/v1/seller/stores/${E.storeId}/products`;
  try {
    const p = E.product ? await api(`${url}/${E.product.id}`, {method: 'PUT', body: buildBody(submit)})
      : await api(url, {method: 'POST', body: buildBody(submit)});
    toast(submit ? t('sentReview') : t('savedDraft'));
    PF.E = {storeId: E.storeId, loadedId: p.id, product: p, f: fromProduct(p), editingTr: null};
    if (location.hash !== `#/store/${E.storeId}/product/${p.id}`) history.replaceState(null, '', `#/store/${E.storeId}/product/${p.id}`);
    drawEditor(PF.cats[P.lang]);
  } catch (e) { showError(e, productErrorField); }
}

async function refreshProduct() {
  const E = PF.E;
  readForm();
  const p = await api(`/v1/seller/stores/${E.storeId}/products/${E.product.id}`);
  E.product = p;
  // Hasil kompres/sampul yang baru selesai → perbarui pratinjau tanpa mengubah isian form.
  for (const m of E.f.media) {
    const fresh = p.media.find(x => x.key === m.key);
    if (fresh) Object.assign(m, {url: fresh.url, poster: fresh.poster_url, processing: fresh.processing, duration: fresh.duration_sec});
  }
  drawEditor(PF.cats[P.lang]);
}
async function pollTranslations() {
  if (!PF.E?.product || !location.hash.includes('/product/' + PF.E.product.id) || PF.E.editingTr) return;
  try { await refreshProduct(); } catch {}
}
async function saveTr(lang) {
  const E = PF.E;
  if (!checkRequired([['#trName']])) return;
  try {
    E.product = await api(`/v1/seller/stores/${E.storeId}/products/${E.product.id}/translations/${lang}`, {method: 'PUT', body: {name: $('#trName').value, description: $('#trDesc').value}});
    E.editingTr = null; toast(t('saved')); drawEditor(PF.cats[P.lang]);
  } catch (e) { showError(e, {name: '#trName', description: '#trDesc'}); }
}
async function retryTr() {
  const E = PF.E;
  try { readForm(); E.product = await api(`/v1/seller/stores/${E.storeId}/products/${E.product.id}/translations/retry`, {method: 'POST'}); drawEditor(PF.cats[P.lang]); }
  catch (e) { toast(errMsg(e)); }
}
async function resetTr(lang) {
  const E = PF.E;
  try { readForm(); E.product = await api(`/v1/seller/stores/${E.storeId}/products/${E.product.id}/translations/${lang}`, {method: 'DELETE'}); drawEditor(PF.cats[P.lang]); }
  catch (e) { toast(errMsg(e)); }
}
async function setPStatus(status) {
  const E = PF.E;
  try { readForm(); E.product = await api(`/v1/seller/stores/${E.storeId}/products/${E.product.id}/status`, {method: 'POST', body: {status}}); drawEditor(PF.cats[P.lang]); }
  catch (e) { toast(errMsg(e)); }
}
async function deleteProduct() {
  const E = PF.E;
  if (!confirm(t('confirmDelete'))) return;
  try { await api(`/v1/seller/stores/${E.storeId}/products/${E.product.id}`, {method: 'DELETE'}); PF.E = null; location.hash = '#/store/' + E.storeId; }
  catch (e) { toast(errMsg(e)); }
}
