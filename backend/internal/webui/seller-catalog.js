'use strict';
// Seller Center: alat katalog untuk toko dengan banyak SKU (migrasi 00052).
//   · Daftar produk: tampilan Kartu / Tabel SKU, filter kategori & stok, 24/50/100 per halaman
//   · Tabel SKU: ubah harga, harga coret, stok, ambang peringatan langsung di baris (PATCH /variants/{id})
//   · Aksi massal: sembunyikan, tampilkan lagi, hapus draf, ubah stok, ubah harga %
//   · Ekspor CSV & impor CSV (pratinjau dulu, lalu terapkan)
//   · Profil toko: kode SKU & ambang stok menipis · Formulir produk: pratinjau SKU otomatis
addStrings({
  en: {viewCards: 'Cards', viewTable: 'SKU table', viewAs: 'View', allCats: 'All categories', stockAll: 'Any stock', stockLow: 'Low stock', stockOut: 'Sold out', perPage: '{n} per page',
    exportCsv: 'Export CSV', importCsv: 'Import CSV', colProduct: 'Product', colAlert: 'Alert at', lowChip: 'Low', outChip: 'Sold out', reservedN: '{n} ordered',
    savedShort: 'Saved', sentToReview: 'The price dropped a lot — Phsar checks this product before it is live again.', selectAll: 'Select all on this page', selectProduct: 'Select {p}',
    selectedN: '{n} products selected', bulkArchive: 'Hide', bulkActivate: 'Show again', bulkDelete: 'Delete drafts', bulkStock: 'Stock…', bulkPrice: 'Price %…', clearSel: 'Clear',
    bulkDone: 'Done: {d}', bulkSkipped: 'skipped: {n}', bulkReview: 'checked by Phsar: {n}', confirmBulkDelete: 'Delete the selected draft products? Products that are not drafts are skipped.',
    stockMode: 'How', stockSet: 'Set every variant to', stockAdd: 'Add or remove (e.g. 10 or -5)', stockAmount: 'Amount', pctLabel: 'Change price by (%)', pctHint: '10 = 10% more, -15 = 15% less. Original (crossed-out) prices change too.',
    pctEx: 'Example: {a} → {b}', applyBtn: 'Apply', noSkus: 'No SKUs match.', defaultN: 'store: {n}',
    impTitle: 'Import prices & stock', impHint: 'Use the file from “Export CSV”. Change price, compare_at_price or stock and keep the sku column. Max 2 MB. Nothing changes until you confirm.',
    impChoose: 'Choose CSV file', impSummary: '{c} to change · {u} unchanged · {e} with problems', impReview: '{n} live products will be checked by Phsar because the price drops a lot.',
    impApply: 'Apply {n} changes', impNothing: 'Nothing to change in this file.', impApplied: 'Updated {n} SKUs', impLine: 'Line {n}', impOther: 'Choose another file', impChecking: 'Checking…',
    imp_unknown_sku: 'SKU not found in your store', imp_duplicate_sku: 'SKU appears twice ({d})', imp_bad_price: 'Price is not valid', imp_bad_compare: 'Original price must be higher than the price',
    imp_bad_stock: 'Stock must be a whole number 0–100000', imp_below_reserved: 'Stock can’t be below {d} (already ordered)', imp_banned: 'Product is blocked by Phsar',
    impErr_no_sku_column: 'The file needs a “sku” column.', impErr_no_value_columns: 'Add at least one column: price, compare_at_price or stock.', impErr_empty: 'The file has no rows.',
    impErr_too_large: 'The file is larger than 2 MB.', impErr_bad_csv: 'This doesn’t look like a CSV file.', impErr_not_utf8: 'Save the file as “CSV UTF-8”.', impErr_too_many_rows: 'Too many rows (max 20,000).',
    catSet: 'Products & stock', skuCode: 'SKU code', skuCodeHint: 'New products get SKUs like {ex}. 2–5 letters or numbers. Existing SKUs don’t change.',
    lowStockAt: 'Low stock alert at', lowStockHint: 'We notify you when a variant has this many or fewer left. 0 = only when sold out. You can set a different number per SKU in the SKU table.',
    skuAutoHint: 'Leave SKU empty and we create one, e.g. {ex}.', quickHint: 'Changes go live right away. A very big price drop is checked by Phsar first.'},
  km: {viewCards: 'កាត', viewTable: 'តារាង SKU', viewAs: 'បង្ហាញជា', allCats: 'គ្រប់ប្រភេទ', stockAll: 'ស្តុកទាំងអស់', stockLow: 'ស្តុកជិតអស់', stockOut: 'អស់ស្តុក', perPage: '{n} ក្នុងមួយទំព័រ',
    exportCsv: 'នាំចេញ CSV', importCsv: 'នាំចូល CSV', colProduct: 'ផលិតផល', colAlert: 'ជូនដំណឹងនៅ', lowChip: 'ជិតអស់', outChip: 'អស់ស្តុក', reservedN: 'បានកម្ម៉ង់ {n}',
    savedShort: 'បានរក្សាទុក', sentToReview: 'តម្លៃធ្លាក់ខ្លាំង — Phsar ពិនិត្យផលិតផលនេះ មុនពេលដាក់លក់វិញ។', selectAll: 'ជ្រើសទាំងអស់ក្នុងទំព័រនេះ', selectProduct: 'ជ្រើស {p}',
    selectedN: 'បានជ្រើស {n} ផលិតផល', bulkArchive: 'លាក់', bulkActivate: 'បង្ហាញវិញ', bulkDelete: 'លុបសេចក្តីព្រាង', bulkStock: 'ស្តុក…', bulkPrice: 'តម្លៃ %…', clearSel: 'សម្អាត',
    bulkDone: 'រួចរាល់៖ {d}', bulkSkipped: 'រំលង៖ {n}', bulkReview: 'Phsar ពិនិត្យ៖ {n}', confirmBulkDelete: 'លុបសេចក្តីព្រាងដែលបានជ្រើស? ផលិតផលដែលមិនមែនជាសេចក្តីព្រាងនឹងត្រូវរំលង។',
    stockMode: 'របៀប', stockSet: 'កំណត់គ្រប់ជម្រើសជា', stockAdd: 'បន្ថែម ឬដក (ឧ. 10 ឬ -5)', stockAmount: 'ចំនួន', pctLabel: 'ប្តូរតម្លៃ (%)', pctHint: '10 = ថ្លៃជាង 10%, -15 = ថោកជាង 15%។ តម្លៃដើម (គូសចោល) ក៏ប្តូរដែរ។',
    pctEx: 'ឧទាហរណ៍៖ {a} → {b}', applyBtn: 'អនុវត្ត', noSkus: 'គ្មាន SKU ត្រូវគ្នា។', defaultN: 'ហាង៖ {n}',
    impTitle: 'នាំចូលតម្លៃ និងស្តុក', impHint: 'ប្រើឯកសារពី «នាំចេញ CSV»។ កែ price, compare_at_price ឬ stock ហើយរក្សាជួរ sku។ ធំបំផុត 2 MB។ គ្មានអ្វីប្តូរទេ រហូតដល់អ្នកបញ្ជាក់។',
    impChoose: 'ជ្រើសឯកសារ CSV', impSummary: 'ត្រូវប្តូរ {c} · មិនប្តូរ {u} · មានបញ្ហា {e}', impReview: 'ផលិតផលកំពុងលក់ {n} នឹងត្រូវ Phsar ពិនិត្យ ព្រោះតម្លៃធ្លាក់ខ្លាំង។',
    impApply: 'អនុវត្តការប្តូរ {n}', impNothing: 'គ្មានអ្វីត្រូវប្តូរក្នុងឯកសារនេះទេ។', impApplied: 'បានកែ {n} SKU', impLine: 'ជួរ {n}', impOther: 'ជ្រើសឯកសារផ្សេង', impChecking: 'កំពុងពិនិត្យ…',
    imp_unknown_sku: 'រកមិនឃើញ SKU ក្នុងហាងរបស់អ្នក', imp_duplicate_sku: 'SKU ស្ទួន ({d})', imp_bad_price: 'តម្លៃមិនត្រឹមត្រូវ', imp_bad_compare: 'តម្លៃដើមត្រូវខ្ពស់ជាងតម្លៃលក់',
    imp_bad_stock: 'ស្តុកត្រូវជាចំនួនគត់ 0–100000', imp_below_reserved: 'ស្តុកមិនអាចតិចជាង {d} (បានកម្ម៉ង់រួច)', imp_banned: 'ផលិតផលត្រូវបាន Phsar ហាមឃាត់',
    impErr_no_sku_column: 'ឯកសារត្រូវមានជួរ «sku»។', impErr_no_value_columns: 'បន្ថែមយ៉ាងហោចមួយជួរ៖ price, compare_at_price ឬ stock។', impErr_empty: 'ឯកសារគ្មានជួរទិន្នន័យ។',
    impErr_too_large: 'ឯកសារធំជាង 2 MB។', impErr_bad_csv: 'នេះមិនមែនជាឯកសារ CSV ទេ។', impErr_not_utf8: 'សូមរក្សាទុកជា «CSV UTF-8»។', impErr_too_many_rows: 'ជួរច្រើនពេក (ច្រើនបំផុត 20,000)។',
    catSet: 'ផលិតផល និងស្តុក', skuCode: 'កូដ SKU', skuCodeHint: 'ផលិតផលថ្មីទទួលបាន SKU ដូចជា {ex}។ អក្សរ ឬលេខ 2–5។ SKU ដែលមានស្រាប់មិនប្តូរទេ។',
    lowStockAt: 'ជូនដំណឹងស្តុកជិតអស់នៅ', lowStockHint: 'យើងជូនដំណឹង ពេលជម្រើសណាមួយនៅសល់ប៉ុណ្ណេះ ឬតិចជាង។ 0 = តែពេលអស់។ អ្នកអាចកំណត់ចំនួនផ្សេងសម្រាប់ SKU នីមួយៗក្នុងតារាង SKU។',
    skuAutoHint: 'ទុក SKU ទទេ យើងនឹងបង្កើតឲ្យ ឧ. {ex}។', quickHint: 'ការកែប្រែមានប្រសិទ្ធភាពភ្លាម។ ការបញ្ចុះតម្លៃខ្លាំងពេក Phsar នឹងពិនិត្យមុន។'},
  zh: {viewCards: '卡片', viewTable: 'SKU 表格', viewAs: '显示方式', allCats: '全部类目', stockAll: '全部库存', stockLow: '库存不足', stockOut: '已售罄', perPage: '每页 {n} 个',
    exportCsv: '导出 CSV', importCsv: '导入 CSV', colProduct: '商品', colAlert: '提醒阈值', lowChip: '不足', outChip: '售罄', reservedN: '已下单 {n}',
    savedShort: '已保存', sentToReview: '价格降幅很大——Phsar 会先审核此商品再重新上架。', selectAll: '全选本页', selectProduct: '选择 {p}',
    selectedN: '已选 {n} 个商品', bulkArchive: '下架', bulkActivate: '重新上架', bulkDelete: '删除草稿', bulkStock: '库存…', bulkPrice: '价格 %…', clearSel: '清除',
    bulkDone: '完成：{d}', bulkSkipped: '跳过：{n}', bulkReview: '待 Phsar 审核：{n}', confirmBulkDelete: '删除所选的草稿商品？非草稿商品会被跳过。',
    stockMode: '方式', stockSet: '所有规格设为', stockAdd: '增加或减少（如 10 或 -5）', stockAmount: '数量', pctLabel: '价格调整（%）', pctHint: '10 = 涨价 10%，-15 = 降价 15%。原价（划线价）也会一起调整。',
    pctEx: '示例：{a} → {b}', applyBtn: '应用', noSkus: '没有匹配的 SKU。', defaultN: '店铺：{n}',
    impTitle: '导入价格和库存', impHint: '请使用“导出 CSV”得到的文件。修改 price、compare_at_price 或 stock，保留 sku 列。最大 2 MB。确认前不会做任何修改。',
    impChoose: '选择 CSV 文件', impSummary: '{c} 项将修改 · {u} 项无变化 · {e} 项有问题', impReview: '{n} 个在售商品因降价幅度大，将由 Phsar 审核。',
    impApply: '应用 {n} 项修改', impNothing: '此文件中没有需要修改的内容。', impApplied: '已更新 {n} 个 SKU', impLine: '第 {n} 行', impOther: '选择其他文件', impChecking: '检查中…',
    imp_unknown_sku: '店铺中找不到此 SKU', imp_duplicate_sku: 'SKU 重复（{d}）', imp_bad_price: '价格无效', imp_bad_compare: '原价必须高于售价',
    imp_bad_stock: '库存必须是 0–100000 的整数', imp_below_reserved: '库存不能低于 {d}（已有订单）', imp_banned: '商品已被 Phsar 封禁',
    impErr_no_sku_column: '文件需要包含“sku”列。', impErr_no_value_columns: '请至少包含一列：price、compare_at_price 或 stock。', impErr_empty: '文件中没有数据行。',
    impErr_too_large: '文件大于 2 MB。', impErr_bad_csv: '这不像是 CSV 文件。', impErr_not_utf8: '请将文件保存为“CSV UTF-8”。', impErr_too_many_rows: '行数太多（最多 20,000 行）。',
    catSet: '商品和库存', skuCode: 'SKU 代码', skuCodeHint: '新商品的 SKU 形如 {ex}。2–5 个字母或数字。已有的 SKU 不会改变。',
    lowStockAt: '库存不足提醒阈值', lowStockHint: '某个规格剩余数量小于或等于此数时通知你。0 = 仅在售罄时提醒。可在 SKU 表格中为单个 SKU 设置不同数值。',
    skuAutoHint: 'SKU 留空会自动生成，例如 {ex}。', quickHint: '修改立即生效。降价幅度过大时会先由 Phsar 审核。'},
});

try { PF.view = localStorage.getItem('phsar-pview') === 'table' ? 'table' : 'cards'; } catch { PF.view = 'cards'; }
PF.sel = new Set();
const PAGE_SIZES = [24, 50, 100];

// ---- Pratinjau SKU otomatis (formulir produk) — aturan sama dengan server (seller.autoSKUs) ----
function skuWordsJS(s) { return (String(s).toUpperCase().match(/[A-Z0-9]+/g) || []); }
function skuSuffixJS(values) {
  return values.map(v => skuWordsJS(v).join('').slice(0, 10)).filter(Boolean).join('-').slice(0, 32);
}
function skuPreviews(f, allCombos) {
  const base = PF.E?.product?.sku_base || S.store?.sku_next || '';
  const out = {single: '', rows: {}, example: ''};
  if (!base) return out;
  if (!f.hasOptions) { out.single = base; return out; }
  const used = new Set(Object.values(f.rows).filter(r => r.on && r.sku).map(r => String(r.sku).trim().toLowerCase()));
  const on = allCombos.filter(c => (f.rows[rowKey(c)] || {on: true}).on);
  for (const c of on) {
    const suf = skuSuffixJS(c), cand = suf ? `${base}-${suf}` : base;
    let s = suf || on.length === 1 ? cand : '';
    for (let n = 1; !s || used.has(s.toLowerCase()); n++) s = `${cand}-${n}`;
    used.add(s.toLowerCase());
    out.rows[rowKey(c)] = s;
    if (!out.example) out.example = s;
  }
  return out;
}

// ---- Kepala daftar: tab status, tampilan, ekspor/impor ----
function catalogHead(id, counts) {
  return `<div class="spread"><div class="tabs">${['all', 'draft', 'pending', 'live', 'changes', 'archived'].map(f =>
      `<button type="button" aria-selected="${PF.filter === f}" data-pf="${f}">${t('f_' + f)}${counts[f] != null ? ` <span class="tcount">${counts[f]}</span>` : ''}</button>`).join('')}</div>
      <a class="btn primary" href="#/store/${id}/product/new">+ ${t('addProduct')}</a></div>
    <div class="cattools"><div class="seg" role="group" aria-label="${t('viewAs')}">
      <button type="button" data-view="cards" aria-pressed="${PF.view === 'cards'}">▦ ${t('viewCards')}</button>
      <button type="button" data-view="table" aria-pressed="${PF.view === 'table'}">☰ ${t('viewTable')}</button></div>
      <span class="grow"></span>
      <button type="button" class="small" id="csvExport">⬇ ${t('exportCsv')}</button>
      <button type="button" class="small" id="csvImport">⬆ ${t('importCsv')}</button></div>`;
}

function catalogFilters(cats) {
  const catOpt = (c, pre) => `<option value="${c.id}" ${PF.category === c.id ? 'selected' : ''}>${pre}${esc(c.name)}</option>`;
  const catList = (cats || []).map(c => catOpt(c, '') + (c.children || []).map(k => catOpt(k, '— ')).join('')).join('');
  return `<div class="catfilters">
    <label class="sr-only" for="pCat">${t('category')}</label><select id="pCat" aria-label="${t('category')}"><option value="">${t('allCats')}</option>${catList}</select>
    <label class="sr-only" for="pStock">${t('stock')}</label><select id="pStock" aria-label="${t('stock')}">${['', 'low', 'out'].map(k =>
      `<option value="${k}" ${PF.stock === k ? 'selected' : ''}>${t(k === 'low' ? 'stockLow' : k === 'out' ? 'stockOut' : 'stockAll')}</option>`).join('')}</select>
    <label class="sr-only" for="pSize">${t('perPage', {n: ''})}</label><select id="pSize" aria-label="${t('perPage', {n: ''})}">${PAGE_SIZES.map(n =>
      `<option value="${n}" ${PF.size === n ? 'selected' : ''}>${t('perPage', {n})}</option>`).join('')}</select></div>`;
}

function bindCatalogHead(id) {
  const again = () => tabProducts(id);
  document.querySelectorAll('[data-pf]').forEach(b => b.onclick = () => { PF.filter = b.dataset.pf; PF.page = 1; PF.sel.clear(); again(); });
  document.querySelectorAll('[data-view]').forEach(b => b.onclick = () => {
    PF.view = b.dataset.view; PF.page = 1; PF.sel.clear();
    if (PF.view === 'table' && PF.size === 24) PF.size = 50;
    try { localStorage.setItem('phsar-pview', PF.view); } catch {}
    again();
  });
  $('#pCat').onchange = e => { PF.category = e.target.value; PF.page = 1; again(); };
  $('#pStock').onchange = e => { PF.stock = e.target.value; PF.page = 1; again(); };
  $('#pSize').onchange = e => { PF.size = Number(e.target.value); PF.page = 1; again(); };
  $('#csvExport').onclick = () => exportCSV(id);
  $('#csvImport').onclick = () => openImport(id);
}

// ---- Tabel SKU ----
const VSORTS = ['updated', 'name', 'sku', 'price_asc', 'price_desc', 'stock_low'];
async function tabVariantTable(id) {
  if (!PAGE_SIZES.includes(PF.size)) PF.size = 50;
  const sort = VSORTS.includes(PF.sort) ? PF.sort : 'updated';
  const query = qs({filter: PF.filter, q: PF.q, sort, page: PF.page, size: PF.size, category: PF.category, stock: PF.stock});
  const [list, cats] = await Promise.all([api(`/v1/seller/stores/${id}/variants?${query}`), categories()]);
  const {variants, counts = {}, page} = list;
  if (page && PF.page > page.pages) { PF.page = page.pages; return tabVariantTable(id); }
  const active = S.store.status === 'active', def = S.store.low_stock_threshold;
  const rows = variants.map(v => {
    const [cls, key] = productState({status: v.product_status, moderation_status: v.moderation_status}, active);
    const banned = v.product_status === 'banned';
    const stockChip = v.available <= 0 ? `<span class="chip bad">${t('outChip')}</span>` : v.available <= v.threshold ? `<span class="chip warn">${t('lowChip')}</span>` : '';
    const dis = banned ? 'disabled' : '';
    return `<tr data-vid="${v.id}" data-pid="${v.product_id}">
      <td class="c-sel"><input type="checkbox" data-sel="${v.product_id}" ${PF.sel.has(v.product_id) ? 'checked' : ''} aria-label="${esc(t('selectProduct', {p: v.product_name}))}"></td>
      <td class="c-prod"><a href="#/store/${id}/product/${v.product_id}" class="vprod">${v.cover_url ? `<img src="${esc(v.cover_url)}" alt="" loading="lazy">` : '<span class="noimg"></span>'}
        <span><b>${esc(v.product_name || '—')}</b>${v.label ? `<small>${esc(v.label)}</small>` : ''}</span></a>
        <span class="vchips"><span class="chip ${cls}">${t(key)}</span>${stockChip}</span></td>
      <td class="c-sku" data-l="${t('sku')}"><code>${esc(v.sku)}</code></td>
      <td data-l="${t('price')}"><input data-v="price" inputmode="decimal" value="${esc(v.price)}" aria-label="${t('price')} ${esc(v.sku)}" ${dis}></td>
      <td data-l="${t('comparePrice')}"><input data-v="compare_at_price" inputmode="decimal" value="${esc(v.compare_at_price || '')}" aria-label="${t('comparePrice')} ${esc(v.sku)}" ${dis}></td>
      <td data-l="${t('stockQty')}"><input data-v="stock_qty" type="number" min="0" value="${v.stock_qty}" aria-label="${t('stockQty')} ${esc(v.sku)}" ${dis}>${v.reserved_qty ? `<small class="muted">${t('reservedN', {n: v.reserved_qty})}</small>` : ''}</td>
      <td data-l="${t('colAlert')}"><input data-v="low_stock_threshold" type="number" min="0" max="1000" value="${v.low_stock_threshold ?? ''}" placeholder="${def}" title="${t('defaultN', {n: def})}" aria-label="${t('colAlert')} ${esc(v.sku)}" ${dis}></td>
      <td class="c-ok" aria-live="polite"></td></tr>`;
  }).join('');
  const pids = [...new Set(variants.map(v => v.product_id))];
  $('#tab').innerHTML = `
    ${active ? '' : `<div class="banner info">${t('storeNotActiveHint')}</div>`}
    ${catalogHead(id, counts)}
    <div class="listbar">${searchBox('pQ', PF.q, t('pSearchPh'))}
      <label class="sr-only" for="pSort">${t('pSort')}</label>
      <select id="pSort" aria-label="${t('pSort')}">${VSORTS.map(k => `<option value="${k}" ${sort === k ? 'selected' : ''}>${t(k === 'sku' ? 'sku' : 'pSort_' + k)}</option>`).join('')}</select></div>
    ${catalogFilters(cats)}
    <p class="muted small">${t('quickHint')}</p>
    ${variants.length ? `<div class="skuwrap"><table class="skutable"><thead><tr>
      <th class="c-sel"><input type="checkbox" id="selAll" aria-label="${t('selectAll')}" ${pids.length && pids.every(p => PF.sel.has(p)) ? 'checked' : ''}></th>
      <th>${t('colProduct')}</th><th>${t('sku')}</th><th>${t('price')}</th><th>${t('comparePrice')}</th><th>${t('stockQty')}</th><th>${t('colAlert')}</th><th></th></tr></thead>
      <tbody>${rows}</tbody></table></div>` : `<section class="card"><p class="muted" style="text-align:center">${t('noSkus')}</p></section>`}
    ${pagerHTML(page)}
    <div id="bulkBar"></div>`;
  $('#pSort').onchange = e => { PF.sort = e.target.value; PF.page = 1; tabProducts(id); };
  bindSearch('pQ', v => { if (v === PF.q) return; PF.q = v; PF.page = 1; tabProducts(id); });
  bindPager($('#tab'), p => { PF.page = p; tabProducts(id); });
  bindCatalogHead(id);
  const syncSel = () => {
    document.querySelectorAll('[data-sel]').forEach(c => { c.checked = PF.sel.has(c.dataset.sel); });
    if ($('#selAll')) $('#selAll').checked = pids.length > 0 && pids.every(p => PF.sel.has(p));
    drawBulkBar(id);
  };
  document.querySelectorAll('[data-sel]').forEach(c => c.onchange = () => { c.checked ? PF.sel.add(c.dataset.sel) : PF.sel.delete(c.dataset.sel); syncSel(); });
  if ($('#selAll')) $('#selAll').onchange = e => { pids.forEach(p => e.target.checked ? PF.sel.add(p) : PF.sel.delete(p)); syncSel(); };
  document.querySelectorAll('.skutable input[data-v]').forEach(inp => {
    inp.dataset.orig = inp.value;
    inp.onkeydown = e => { if (e.key === 'Enter') { e.preventDefault(); inp.blur(); } if (e.key === 'Escape') { inp.value = inp.dataset.orig; inp.blur(); } };
    inp.onchange = () => saveCell(id, inp);
  });
  drawBulkBar(id);
}

async function saveCell(id, inp) {
  const tr = inp.closest('tr'), f = inp.dataset.v, ok = tr.querySelector('.c-ok');
  clearFieldErrors(tr);
  let val = inp.value.trim();
  const body = {};
  if (f === 'stock_qty') {
    if (val === '' || !/^\d+$/.test(val)) { markField(inp, t('required')); return; }
    body[f] = Number(val);
  } else if (f === 'low_stock_threshold') {
    if (val !== '' && !/^\d+$/.test(val)) { markField(inp, t('required')); return; }
    body[f] = val === '' ? null : Number(val);
  } else if (f === 'price') {
    if (!(Number(val) > 0)) { markField(inp, t('positive')); return; }
    body[f] = val;
  } else {
    const price = Number(tr.querySelector('[data-v=price]').value);
    if (val !== '' && !(Number(val) > price)) { markField(inp, t('compareHigher')); return; }
    body[f] = val;
  }
  ok.textContent = '…';
  try {
    const r = await api(`/v1/seller/stores/${id}/variants/${tr.dataset.vid}`, {method: 'PATCH', body});
    const v = r.variant;
    for (const k of ['price', 'compare_at_price', 'stock_qty', 'low_stock_threshold']) {
      const el = tr.querySelector(`[data-v=${k}]`);
      if (el) el.value = el.dataset.orig = v[k] ?? '';
    }
    ok.innerHTML = `<span class="saved">✓ <span class="sr-only">${t('savedShort')}</span></span>`;
    setTimeout(() => { if (ok.querySelector('.saved')) ok.textContent = ''; }, 2500);
    if (r.review) toast(t('sentToReview'), 'warning');
  } catch (e) {
    ok.textContent = '';
    markField(inp, errMsg(e));
  }
}

function drawBulkBar(id) {
  const bar = $('#bulkBar');
  if (!bar) return;
  const n = PF.sel.size;
  bar.className = n ? 'bulkbar' : '';
  bar.innerHTML = n ? `<b class="bsel">${t('selectedN', {n})}</b><button type="button" class="linkbtn" id="selClear">${t('clearSel')}</button>
    <button type="button" data-bulk="archive">${t('bulkArchive')}</button><button type="button" data-bulk="activate">${t('bulkActivate')}</button>
    <button type="button" data-bulk="stock">${t('bulkStock')}</button><button type="button" data-bulk="price_pct">${t('bulkPrice')}</button>
    <button type="button" class="danger" data-bulk="delete_drafts">${t('bulkDelete')}</button>` : '';
  if (!n) return;
  $('#selClear').onclick = () => { PF.sel.clear(); tabProducts(id); };
  bar.querySelectorAll('[data-bulk]').forEach(b => b.onclick = () => {
    const a = b.dataset.bulk;
    if (a === 'stock' || a === 'price_pct') return bulkDialog(id, a);
    if (a === 'delete_drafts' && !confirm(t('confirmBulkDelete'))) return;
    runBulk(id, {action: a});
  });
}

async function runBulk(id, body) {
  try {
    const r = await api(`/v1/seller/stores/${id}/products/bulk`, {method: 'POST', body: Object.assign({product_ids: [...PF.sel]}, body)});
    const parts = [t('bulkDone', {d: r.done})];
    if (r.skipped.length) parts.push(t('bulkSkipped', {n: r.skipped.length}));
    if (r.review) parts.push(t('bulkReview', {n: r.review}));
    toast(parts.join(' · '), r.review || r.skipped.length ? 'warning' : 'success');
    PF.sel.clear();
    tabProducts(id);
    return true;
  } catch (e) { toast(errMsg(e), 'error'); return false; }
}

function bulkDialog(id, action) {
  const dlg = document.createElement('dialog');
  dlg.className = 'dlg';
  const stock = action === 'stock';
  dlg.innerHTML = `<form method="dialog"><div class="spread"><h2>${t(stock ? 'bulkStock' : 'bulkPrice').replace('…', '')}</h2><button value="cancel" aria-label="${t('cancel')}">✕</button></div>
    <p class="muted">${t('selectedN', {n: PF.sel.size})}</p>
    ${stock ? `<fieldset class="radios"><legend>${t('stockMode')}</legend>
        <label><input type="radio" name="sm" value="set" checked> ${t('stockSet')}</label>
        <label><input type="radio" name="sm" value="add"> ${t('stockAdd')}</label></fieldset>
      <label>${t('stockAmount')}<input id="bAmt" type="number" inputmode="numeric" value="0"></label>`
    : `<label>${t('pctLabel')}<input id="bAmt" inputmode="decimal" value="10"></label><p class="muted small">${t('pctHint')}</p><p class="muted small" id="bEx"></p>`}
    <div class="row" data-errors-anchor><button type="button" class="primary" id="bGo">${t('applyBtn')}</button><button value="cancel">${t('cancel')}</button></div></form>`;
  document.body.append(dlg);
  dlg.addEventListener('close', () => dlg.remove());
  dlg.showModal();
  const amt = dlg.querySelector('#bAmt');
  const ex = () => { const p = Number(amt.value); if (dlg.querySelector('#bEx')) dlg.querySelector('#bEx').textContent = isFinite(p) ? t('pctEx', {a: '$10.00', b: '$' + (10 * (100 + p) / 100).toFixed(2)}) : ''; };
  if (!stock) { amt.oninput = ex; ex(); }
  dlg.querySelector('#bGo').onclick = async () => {
    const v = amt.value.trim();
    const anchor = dlg.querySelector('[data-errors-anchor]');
    if (stock) {
      const mode = dlg.querySelector('input[name=sm]:checked').value;
      if (!checkRequired([[amt, t('required'), el => /^-?\d+$/.test(el.value.trim()) && (mode === 'add' ? Number(el.value) !== 0 : Number(el.value) >= 0)]], anchor)) return;
      if (await runBulk(id, {action: 'stock', stock_mode: mode, stock: Number(v)})) dlg.close();
    } else {
      if (!checkRequired([[amt, t('required'), el => Number(el.value) !== 0 && Number(el.value) > -100 && Number(el.value) <= 500]], anchor)) return;
      if (await runBulk(id, {action: 'price_pct', pct: v})) dlg.close();
    }
  };
}

// ---- CSV ----
async function exportCSV(id) {
  try {
    const r = await fetch(`/v1/seller/stores/${id}/products/export.csv`, {headers: {Authorization: 'Bearer ' + P.token, 'X-Lang': P.lang}});
    if (!r.ok) throw new Error(r.statusText);
    const name = (r.headers.get('Content-Disposition') || '').match(/filename="([^"]+)"/)?.[1] || 'phsar-skus.csv';
    const url = URL.createObjectURL(await r.blob());
    const a = Object.assign(document.createElement('a'), {href: url, download: name});
    document.body.append(a); a.click(); a.remove();
    setTimeout(() => URL.revokeObjectURL(url), 5000);
  } catch (e) { toast(errMsg(e), 'error'); }
}

function openImport(id) {
  const dlg = document.createElement('dialog');
  dlg.className = 'stockdlg impdlg';
  document.body.append(dlg);
  dlg.addEventListener('close', () => dlg.remove());
  let file = null;
  const send = async apply => {
    const fd = new FormData();
    fd.append('file', file);
    return api(`/v1/seller/stores/${id}/products/import${apply ? '?apply=1' : ''}`, {method: 'POST', body: fd});
  };
  const head = `<div class="spread"><h2>${t('impTitle')}</h2><button type="button" data-close aria-label="${t('cancel')}">✕</button></div>`;
  const pick = (err) => {
    dlg.innerHTML = `<form method="dialog">${head}<p class="muted">${t('impHint')}</p>
      ${err ? `<div class="banner bad" role="alert">${esc(err)}</div>` : ''}
      <label class="btn primary impfile">⬆ ${t('impChoose')}<input type="file" id="impFile" accept=".csv,text/csv,text/plain" class="hidden"></label>
      <button type="button" class="linkbtn" id="impExport">⬇ ${t('exportCsv')}</button></form>`;
    bindClose();
    dlg.querySelector('#impExport').onclick = () => exportCSV(id);
    dlg.querySelector('#impFile').onchange = async e => {
      file = e.target.files[0];
      if (!file) return;
      if (file.size > 2 * 1024 * 1024) return pick(t('impErr_too_large'));
      dlg.querySelector('.impfile').textContent = t('impChecking');
      try { preview(await send(false)); } catch (er) { pick(er.code?.startsWith('import_') ? t('impErr_' + er.code.slice(7)) : errMsg(er)); }
    };
  };
  const errText = r => r.error ? t('imp_' + r.error.code, {d: esc(r.error.detail || '')}) : '';
  const fmt = (o, k) => o ? (k === 'stock' ? String(o.stock) : o[k] == null ? '—' : '$' + o[k]) : '';
  const cell = (r, k) => { const a = fmt(r.old, k), b = fmt(r.new, k); return a === b ? `<span class="muted">${esc(a)}</span>` : `<s class="muted">${esc(a)}</s> <b>${esc(b)}</b>`; };
  const preview = r => {
    const rows = r.rows.slice(0, 300).map(x => `<tr class="${x.error ? 'imperr' : ''}"><td data-l="${t('impLine', {n: ''}).trim()}">${x.line}</td>
      <td data-l="${t('sku')}"><code>${esc(x.sku)}</code>${x.product ? `<small>${esc(x.product)}${x.variant ? ' · ' + esc(x.variant) : ''}</small>` : ''}</td>
      ${x.error ? `<td colspan="3" class="bad">${errText(x)}</td>` : `<td data-l="${t('price')}">${cell(x, 'price')}</td><td data-l="${t('comparePrice')}">${cell(x, 'compare_at_price')}</td><td data-l="${t('stockQty')}">${cell(x, 'stock')}</td>`}</tr>`).join('');
    dlg.innerHTML = `<form method="dialog">${head}
      <p><b>${esc(file.name)}</b></p>
      <div class="impsum"><span class="chip ok">${t('impSummary', {c: r.changed, u: r.unchanged, e: r.errors})}</span></div>
      ${r.review ? `<div class="banner warn">${t('impReview', {n: r.review})}</div>` : ''}
      ${r.rows.length ? `<div class="skuwrap"><table class="skutable imptable"><thead><tr><th>#</th><th>${t('sku')}</th><th>${t('price')}</th><th>${t('comparePrice')}</th><th>${t('stockQty')}</th></tr></thead><tbody>${rows}</tbody></table></div>`
        : `<p class="muted">${t('impNothing')}</p>`}
      <div class="wizfoot"><button type="button" class="primary" id="impApply" ${r.changed ? '' : 'disabled'}>${t('impApply', {n: r.changed})}</button>
        <span class="grow"></span><button type="button" id="impOther">${t('impOther')}</button></div></form>`;
    bindClose();
    dlg.querySelector('#impOther').onclick = () => pick();
    dlg.querySelector('#impApply').onclick = async () => {
      dlg.querySelector('#impApply').disabled = true;
      try {
        const a = await send(true);
        dlg.close();
        toast(t('impApplied', {n: a.changed}) + (a.review ? ' · ' + t('bulkReview', {n: a.review}) : ''), a.review ? 'warning' : 'success');
        tabProducts(id);
      } catch (er) { dlg.querySelector('#impApply').disabled = false; toast(errMsg(er), 'error'); }
    };
  };
  const bindClose = () => dlg.querySelectorAll('[data-close]').forEach(b => b.onclick = () => dlg.close());
  pick();
  dlg.showModal();
}

// ---- Profil toko: kode SKU & ambang stok (pemilik) ----
function catalogSettingsHTML(st) {
  return `<section class="card" id="catSettings"><h2>${t('catSet')}</h2>
    <div class="grid">
      <label>${t('skuCode')}<input id="skuPrefix" value="${esc(st.sku_prefix || '')}" maxlength="5" autocapitalize="characters" spellcheck="false"></label>
      <label>${t('lowStockAt')}<input id="lowStock" type="number" min="0" max="1000" inputmode="numeric" value="${st.low_stock_threshold ?? 3}"></label></div>
    <p class="muted small" id="skuEx">${t('skuCodeHint', {ex: esc(st.sku_next || '')})}</p>
    <p class="muted small">${t('lowStockHint')}</p>
    <div data-errors-anchor><button type="button" id="catSave">${t('save')}</button></div></section>`;
}

function bindCatalogSettings(id) {
  const inp = $('#skuPrefix');
  if (!inp) return;
  const no = (S.store.sku_next || '').split('-').pop();
  inp.oninput = () => {
    inp.value = inp.value.toUpperCase().replace(/[^A-Z0-9]/g, '');
    $('#skuEx').textContent = t('skuCodeHint', {ex: `${inp.value || '…'}-${no}`});
  };
  $('#catSave').onclick = async () => {
    const anchor = $('#catSettings [data-errors-anchor]');
    if (!checkRequired([[inp, t('skuCodeHint', {ex: 'ABC-00001'}), el => /^[A-Z0-9]{2,5}$/.test(el.value)],
      [$('#lowStock'), t('required'), el => /^\d+$/.test(el.value) && Number(el.value) <= 1000]], anchor)) return;
    try {
      S.store = await api(`/v1/seller/stores/${id}`, {method: 'PATCH', body: {sku_prefix: inp.value, low_stock_threshold: Number($('#lowStock').value)}});
      toast(t('saved'));
      $('#skuEx').textContent = t('skuCodeHint', {ex: esc(S.store.sku_next)});
    } catch (e) { showError(e, {sku_prefix: '#skuPrefix', low_stock_threshold: '#lowStock'}, anchor); }
  };
}
