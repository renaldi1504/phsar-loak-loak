'use strict';
// Admin: sesi flash sale (tab "Flash sales" di Vouchers & flash sales). Banyak produk per sesi; pembiayaan dipilih per sesi.
addStrings({
  en: {fsUpcoming: '{n} upcoming', fsNew: 'New session', fsEdit: 'Edit session', fsNone: 'No flash sale sessions yet.', fsNameEn: 'Name (English)', fsNameKm: 'Name (Khmer)',
    fsNameZh: 'Name (Chinese)', fsStart: 'Starts', fsEnd: 'Ends', fsFunding: 'Who pays the discount', fsF_seller: 'Seller pays (seller must agree)', fsF_phsar: 'Phsar pays (seller gets the normal price)',
    fsF_split: 'Shared', fsPhsarPct: 'Phsar pays %', fsVouchers: 'Buyers may also use Phsar vouchers on these items', fsMinDisc: 'Minimum discount %', fsPerBuyer: 'Max per buyer (each product)',
    fsOpen: 'Sellers can register products', fsRegUntil: 'Registration closes (empty = when the session starts)', fsPublished: 'Published (visible to buyers)',
    fsSave: 'Save session', fsBack: 'All sessions', fsCancelS: 'Cancel session', fsCancelConfirm: 'Cancel this session? Flash prices stop immediately.',
    fsSt_draft: 'Draft', fsSt_scheduled: 'Scheduled', fsSt_live: 'Live', fsSt_ended: 'Ended', fsSt_cancelled: 'Cancelled',
    fsFundSeller: 'Seller pays', fsFundPhsar: 'Phsar pays', fsFundSplit: 'Phsar {p}% · seller {s}%', fsCounts: '{a} products · {p} to review · {i} waiting for seller',
    fsSold: '{n} sold · ${r} · Phsar paid ${s}', fsItems: 'Products in this session', fsAdd: 'Add products', fsFind: 'Search products (name, SKU or store)',
    fsBulkOff: 'Discount % for selected', fsBulkQty: 'Quota each', fsApply: 'Apply to selected', fsAddN: 'Add {n} products', fsPicked: 'Selected ({n})',
    fsNormal: 'Normal', fsFlashP: 'Flash price', fsOff: 'Off', fsQuota: 'Quota', fsStock: 'Stock', fsSoldH: 'Sold', fsSource: 'From', fsStatus: 'Status', fsPerB: 'Per buyer',
    fsSrc_admin: 'Phsar', fsSrc_seller: 'Seller', fsI_invited: 'Waiting for seller', fsI_pending: 'To review', fsI_approved: 'Approved', fsI_rejected: 'Rejected',
    fsI_declined: 'Declined by seller', fsI_withdrawn: 'Withdrawn', fsI_removed: 'Removed', fsApprove: 'Approve', fsReject: 'Reject', fsRemove: 'Remove', fsSaveRow: 'Save',
    fsRejectWhy: 'Why? (shown to the seller)', fsRemoveConfirm: 'Remove this product from the session?', fsSaved: 'Saved', fsNoItems: 'No products yet. Add products below or open registration for sellers.',
    fsInviteNote: 'Seller pays part of the discount here, so products you add are sent to the seller to accept first.', fsMinNote: 'At least {p}% below the normal price.',
    fsOthers: 'Closed ({n})', fsStockLow: 'max {n}'},
  km: {fsUpcoming: 'នឹងមកដល់ {n}', fsNew: 'វគ្គថ្មី', fsEdit: 'កែវគ្គ', fsNone: 'មិនទាន់មានវគ្គ Flash sale។', fsNameEn: 'ឈ្មោះ (អង់គ្លេស)', fsNameKm: 'ឈ្មោះ (ខ្មែរ)',
    fsNameZh: 'ឈ្មោះ (ចិន)', fsStart: 'ចាប់ផ្តើម', fsEnd: 'បញ្ចប់', fsFunding: 'អ្នកចេញថ្លៃបញ្ចុះ', fsF_seller: 'អ្នកលក់ចេញ (អ្នកលក់ត្រូវយល់ព្រម)', fsF_phsar: 'Phsar ចេញ (អ្នកលក់ទទួលតម្លៃធម្មតា)',
    fsF_split: 'ចែកគ្នា', fsPhsarPct: 'Phsar ចេញ %', fsVouchers: 'អ្នកទិញអាចប្រើប័ណ្ណ Phsar លើទំនិញទាំងនេះផងដែរ', fsMinDisc: 'បញ្ចុះអប្បបរមា %', fsPerBuyer: 'ទិញបានច្រើនបំផុតក្នុងម្នាក់ (ផលិតផលនីមួយៗ)',
    fsOpen: 'អ្នកលក់អាចចុះឈ្មោះផលិតផល', fsRegUntil: 'បិទការចុះឈ្មោះ (ទទេ = ពេលវគ្គចាប់ផ្តើម)', fsPublished: 'បានផ្សាយ (អ្នកទិញមើលឃើញ)',
    fsSave: 'រក្សាទុកវគ្គ', fsBack: 'វគ្គទាំងអស់', fsCancelS: 'បោះបង់វគ្គ', fsCancelConfirm: 'បោះបង់វគ្គនេះ? តម្លៃ flash ឈប់ភ្លាមៗ។',
    fsSt_draft: 'សេចក្តីព្រាង', fsSt_scheduled: 'បានកំណត់ពេល', fsSt_live: 'កំពុងដំណើរការ', fsSt_ended: 'បានបញ្ចប់', fsSt_cancelled: 'បានបោះបង់',
    fsFundSeller: 'អ្នកលក់ចេញ', fsFundPhsar: 'Phsar ចេញ', fsFundSplit: 'Phsar {p}% · អ្នកលក់ {s}%', fsCounts: 'ផលិតផល {a} · ត្រូវពិនិត្យ {p} · រង់ចាំអ្នកលក់ {i}',
    fsSold: 'លក់បាន {n} · ${r} · Phsar ចេញ ${s}', fsItems: 'ផលិតផលក្នុងវគ្គនេះ', fsAdd: 'បន្ថែមផលិតផល', fsFind: 'ស្វែងរកផលិតផល (ឈ្មោះ SKU ឬហាង)',
    fsBulkOff: 'បញ្ចុះ % សម្រាប់ដែលបានជ្រើស', fsBulkQty: 'កូតានីមួយៗ', fsApply: 'អនុវត្តលើដែលបានជ្រើស', fsAddN: 'បន្ថែមផលិតផល {n}', fsPicked: 'បានជ្រើស ({n})',
    fsNormal: 'ធម្មតា', fsFlashP: 'តម្លៃ flash', fsOff: 'បញ្ចុះ', fsQuota: 'កូតា', fsStock: 'ស្តុក', fsSoldH: 'លក់បាន', fsSource: 'ពី', fsStatus: 'ស្ថានភាព', fsPerB: 'ក្នុងម្នាក់',
    fsSrc_admin: 'Phsar', fsSrc_seller: 'អ្នកលក់', fsI_invited: 'រង់ចាំអ្នកលក់', fsI_pending: 'ត្រូវពិនិត្យ', fsI_approved: 'បានអនុម័ត', fsI_rejected: 'បានបដិសេធ',
    fsI_declined: 'អ្នកលក់បដិសេធ', fsI_withdrawn: 'បានដកចេញ', fsI_removed: 'បានយកចេញ', fsApprove: 'អនុម័ត', fsReject: 'បដិសេធ', fsRemove: 'យកចេញ', fsSaveRow: 'រក្សាទុក',
    fsRejectWhy: 'ហេតុអ្វី? (បង្ហាញដល់អ្នកលក់)', fsRemoveConfirm: 'យកផលិតផលនេះចេញពីវគ្គ?', fsSaved: 'បានរក្សាទុក', fsNoItems: 'មិនទាន់មានផលិតផល។ បន្ថែមខាងក្រោម ឬបើកការចុះឈ្មោះសម្រាប់អ្នកលក់។',
    fsInviteNote: 'អ្នកលក់ចេញថ្លៃបញ្ចុះមួយផ្នែកនៅទីនេះ ដូច្នេះផលិតផលដែលអ្នកបន្ថែមត្រូវផ្ញើទៅអ្នកលក់ឱ្យទទួលសិន។', fsMinNote: 'ទាបជាងតម្លៃធម្មតាយ៉ាងតិច {p}%។',
    fsOthers: 'បានបិទ ({n})', fsStockLow: 'អតិបរមា {n}'},
  zh: {fsUpcoming: '即将开始 {n}', fsNew: '新建场次', fsEdit: '编辑场次', fsNone: '暂无限时抢购场次。', fsNameEn: '名称（英文）', fsNameKm: '名称（高棉文）',
    fsNameZh: '名称（中文）', fsStart: '开始', fsEnd: '结束', fsFunding: '折扣由谁承担', fsF_seller: '卖家承担（需卖家同意）', fsF_phsar: 'Phsar 承担（卖家按原价结算）',
    fsF_split: '共同承担', fsPhsarPct: 'Phsar 承担 %', fsVouchers: '买家可同时对这些商品使用 Phsar 优惠券', fsMinDisc: '最低折扣 %', fsPerBuyer: '每位买家限购（每个商品）',
    fsOpen: '卖家可报名商品', fsRegUntil: '报名截止（空 = 场次开始时）', fsPublished: '已发布（买家可见）',
    fsSave: '保存场次', fsBack: '全部场次', fsCancelS: '取消场次', fsCancelConfirm: '取消此场次？抢购价立即停止。',
    fsSt_draft: '草稿', fsSt_scheduled: '已排期', fsSt_live: '进行中', fsSt_ended: '已结束', fsSt_cancelled: '已取消',
    fsFundSeller: '卖家承担', fsFundPhsar: 'Phsar 承担', fsFundSplit: 'Phsar {p}% · 卖家 {s}%', fsCounts: '{a} 个商品 · {p} 待审核 · {i} 等待卖家',
    fsSold: '已售 {n} · ${r} · Phsar 承担 ${s}', fsItems: '本场商品', fsAdd: '添加商品', fsFind: '搜索商品（名称、SKU 或店铺）',
    fsBulkOff: '所选商品折扣 %', fsBulkQty: '每个配额', fsApply: '应用到所选', fsAddN: '添加 {n} 个商品', fsPicked: '已选择（{n}）',
    fsNormal: '原价', fsFlashP: '抢购价', fsOff: '折扣', fsQuota: '配额', fsStock: '库存', fsSoldH: '已售', fsSource: '来源', fsStatus: '状态', fsPerB: '限购',
    fsSrc_admin: 'Phsar', fsSrc_seller: '卖家', fsI_invited: '等待卖家', fsI_pending: '待审核', fsI_approved: '已通过', fsI_rejected: '已拒绝',
    fsI_declined: '卖家已拒绝', fsI_withdrawn: '已撤回', fsI_removed: '已移除', fsApprove: '通过', fsReject: '拒绝', fsRemove: '移除', fsSaveRow: '保存',
    fsRejectWhy: '原因？（卖家可见）', fsRemoveConfirm: '从本场移除此商品？', fsSaved: '已保存', fsNoItems: '暂无商品。在下方添加，或开放卖家报名。',
    fsInviteNote: '本场卖家承担部分折扣，因此您添加的商品需先由卖家接受。', fsMinNote: '至少比原价低 {p}%。',
    fsOthers: '已关闭（{n}）', fsStockLow: '最多 {n}'},
});

const FS = {open: null, form: null, data: null, picked: [], results: []};
const FS_CHIP = {draft: '', scheduled: 'warn', live: 'ok', ended: '', cancelled: 'bad'};
const FS_ICHIP = {invited: 'warn', pending: 'warn', approved: 'ok', rejected: 'bad', declined: 'bad', withdrawn: '', removed: ''};

function fsName(s) { return (P.lang === 'km' && s.name_km) || (P.lang === 'zh' && s.name_zh) || s.name_en; }
function fsFund(s) { return s.phsar_share_pct === 0 ? t('fsFundSeller') : s.phsar_share_pct === 100 ? t('fsFundPhsar') : t('fsFundSplit', {p: s.phsar_share_pct, s: 100 - s.phsar_share_pct}); }
function fsMoney(v) { const n = parseFloat(v); return '$' + (isNaN(n) ? v : n.toFixed(2)); }

async function drawFlashSessions() {
  if (FS.open) return drawFlashDetail(FS.open);
  const list = PM.data.sessions || [];
  const cards = list.map(s => `<section class="card fs-card" data-fsopen="${s.id}" tabindex="0" role="button">
      <div class="spread"><div><h2>${esc(fsName(s))}</h2><span class="muted">${pmWhen(s.starts_at)} → ${pmWhen(s.ends_at)}</span></div><span class="chip ${FS_CHIP[s.status]}">${t('fsSt_' + s.status)}</span></div>
      <div class="row" style="gap:6px;flex-wrap:wrap"><span class="chip pink">${esc(fsFund(s))}</span>${s.open_for_sellers ? `<span class="chip">${t('fsOpen')}</span>` : ''}${s.allow_vouchers ? '' : '<span class="chip">no voucher</span>'}</div>
      <p class="muted small" style="margin:0">${t('fsCounts', {a: s.counts.approved, p: s.counts.pending, i: s.counts.invited})} · ${t('fsSold', {n: s.counts.sold, r: parseFloat(s.counts.revenue).toFixed(2), s: parseFloat(s.counts.subsidy).toFixed(2)})}</p></section>`).join('');
  $('#pmBody').innerHTML = (FS.form ? fsFormHTML(FS.form) : `<div class="row"><button type="button" class="primary" id="fsNew">${t('fsNew')}</button></div>`) +
    (cards || `<section class="card"><p class="muted">${t('fsNone')}</p></section>`) + phRulesHTML(); // pricehistory: sakelar aturan (super admin)
  phBindRules(); // pricehistory
  if ($('#fsNew')) $('#fsNew').onclick = () => { FS.form = {}; drawFlashSessions(); };
  if (FS.form) bindFsForm(() => { FS.form = null; drawFlashSessions(); });
  document.querySelectorAll('[data-fsopen]').forEach(c => { c.onclick = () => { FS.open = c.dataset.fsopen; FS.picked = []; drawFlashSessions(); }; c.onkeydown = e => { if (e.key === 'Enter') c.click(); }; });
}

function fsFormHTML(s) {
  const pct = s.phsar_share_pct ?? 0, mode = pct === 0 ? 'seller' : pct === 100 ? 'phsar' : 'split';
  return `<section class="card" id="fsForm"><h2>${t(s.id ? 'fsEdit' : 'fsNew')}</h2>
    <div class="grid">
      <label>${t('fsNameEn')}<input name="name_en" maxlength="80" value="${esc(s.name_en || '')}" placeholder="12.12 Midnight Sale"></label>
      <label>${t('fsNameKm')}<input name="name_km" maxlength="80" value="${esc(s.name_km || '')}"></label>
      <label>${t('fsNameZh')}<input name="name_zh" maxlength="80" value="${esc(s.name_zh || '')}"></label>
      <label>${t('fsStart')}<input name="starts_at" type="datetime-local" value="${pmLocal(s.starts_at)}"></label>
      <label>${t('fsEnd')}<input name="ends_at" type="datetime-local" value="${pmLocal(s.ends_at)}"></label>
    </div>
    <fieldset class="fs-fund"><legend>${t('fsFunding')}</legend>
      ${['seller', 'phsar', 'split'].map(m => `<label class="inline"><input type="radio" name="fund" value="${m}" ${mode === m ? 'checked' : ''}> ${t('fsF_' + m)}</label>`).join('')}
      <label id="fsPctBox" ${mode === 'split' ? '' : 'hidden'}>${t('fsPhsarPct')}<input name="phsar_share_pct" inputmode="numeric" value="${mode === 'split' ? pct : 50}" style="max-width:120px"></label>
    </fieldset>
    <div class="grid">
      <label>${t('fsMinDisc')}<input name="min_discount_pct" inputmode="numeric" value="${s.min_discount_pct || 10}"></label>
      <label>${t('fsPerBuyer')}<input name="per_buyer" inputmode="numeric" value="${s.per_buyer || 2}"></label>
      <label>${t('fsRegUntil')}<input name="register_until" type="datetime-local" value="${pmLocal(s.register_until)}"></label>
    </div>
    <label class="inline"><input type="checkbox" name="allow_vouchers" ${s.allow_vouchers ?? true ? 'checked' : ''}> ${t('fsVouchers')}</label>
    <label class="inline"><input type="checkbox" name="open_for_sellers" ${s.open_for_sellers ? 'checked' : ''}> ${t('fsOpen')}</label>
    <label class="inline"><input type="checkbox" name="published" ${s.published ? 'checked' : ''}> ${t('fsPublished')}</label>
    <div class="row"><button type="button" class="primary" id="fsSave">${t('fsSave')}</button><button type="button" id="fsCancelForm">${t('pmCancel')}</button></div></section>`;
}

function bindFsForm(done) {
  const s = FS.form, val = n => $(`#fsForm [name="${n}"]`);
  document.querySelectorAll('#fsForm [name="fund"]').forEach(r => r.onchange = () => { $('#fsPctBox').hidden = r.value !== 'split' || !r.checked; });
  $('#fsCancelForm').onclick = done;
  $('#fsSave').onclick = async () => {
    document.querySelectorAll('#fsForm .invalid').forEach(el => el.classList.remove('invalid'));
    const mode = document.querySelector('#fsForm [name="fund"]:checked').value;
    const body = {name_en: val('name_en').value, name_km: val('name_km').value, name_zh: val('name_zh').value,
      starts_at: pmISO(val('starts_at').value) || '', ends_at: pmISO(val('ends_at').value) || '',
      phsar_share_pct: mode === 'seller' ? 0 : mode === 'phsar' ? 100 : parseInt(val('phsar_share_pct').value, 10) || 0,
      allow_vouchers: val('allow_vouchers').checked, min_discount_pct: parseInt(val('min_discount_pct').value, 10) || 0,
      per_buyer: parseInt(val('per_buyer').value, 10) || 0, open_for_sellers: val('open_for_sellers').checked,
      register_until: pmISO(val('register_until').value) || '', published: val('published').checked};
    try {
      const r = await api(s.id ? `/admin/api/promos/flash/${s.id}` : '/admin/api/promos/flash', {method: s.id ? 'PUT' : 'POST', body});
      toast(t('fsSaved')); FS.form = null; FS.open = r.session.id; FS.data = r;
      PM.data.sessions = (await api('/admin/api/promos/flash')).sessions; drawFlashDetail(FS.open, r);
    } catch (e) { pmFail(e, 'fsForm'); }
  };
}

async function drawFlashDetail(id, preload) {
  const d = preload || await api('/admin/api/promos/flash/' + id);
  FS.data = d;
  const s = d.session, editable = s.status !== 'ended' && s.status !== 'cancelled';
  const active = d.items.filter(i => ['pending', 'invited', 'approved'].includes(i.status)).sort((a, b) => ['pending', 'invited', 'approved'].indexOf(a.status) - ['pending', 'invited', 'approved'].indexOf(b.status));
  const closed = d.items.filter(i => !['pending', 'invited', 'approved'].includes(i.status));
  const row = i => `<tr data-row="${i.id}"><td><div class="fs-prod">${i.image_url ? `<img src="${esc(i.image_url)}" alt="">` : '<span class="fs-noimg"></span>'}<div><b>${esc(i.product_name)}</b>${i.options ? ` <span class="muted">${esc(i.options)}</span>` : ''}<div class="muted small">${esc(i.store_name)}</div></div></div></td>
    <td>${fsMoney(i.normal_price)}</td>
    <td>${editable && i.status !== 'pending' && ['invited', 'approved'].includes(i.status) ? `<input class="fs-in" name="sale_price" value="${esc(i.sale_price)}" inputmode="decimal">` : `<b>${fsMoney(i.sale_price)}</b>`}</td>
    <td><span class="chip pink">-${i.off_pct}%</span></td>
    <td>${editable && ['invited', 'approved', 'pending'].includes(i.status) ? `<input class="fs-in sm" name="qty_cap" value="${i.qty_cap}" inputmode="numeric">` : i.qty_cap}<div class="muted small">${t('fsStockLow', {n: i.stock})}</div></td>
    <td>${i.sold}</td><td>${t('fsSrc_' + i.source)}</td><td><span class="chip ${FS_ICHIP[i.status]}">${t('fsI_' + i.status)}</span>${i.note ? `<div class="muted small">${esc(i.note)}</div>` : ''}</td>
    <td class="fs-act">${!editable ? '' : i.status === 'pending' ? `<button type="button" class="primary" data-fsdec="approve" data-id="${i.id}">${t('fsApprove')}</button><button type="button" class="danger" data-fsdec="reject" data-id="${i.id}">${t('fsReject')}</button>`
      : `<button type="button" data-fssave="${i.id}">${t('fsSaveRow')}</button><button type="button" class="danger" data-fsdec="remove" data-id="${i.id}">${t('fsRemove')}</button>`}</td></tr>`;
  const head = `<thead><tr><th></th><th>${t('fsNormal')}</th><th>${t('fsFlashP')}</th><th>${t('fsOff')}</th><th>${t('fsQuota')}</th><th>${t('fsSoldH')}</th><th>${t('fsSource')}</th><th>${t('fsStatus')}</th><th></th></tr></thead>`;
  $('#pmBody').innerHTML = `<p><a href="#" id="fsBack">← ${t('fsBack')}</a></p>
    ${FS.form ? fsFormHTML(FS.form) : `<section class="card"><div class="spread"><div><h2>${esc(fsName(s))}</h2><span class="muted">${pmWhen(s.starts_at)} → ${pmWhen(s.ends_at)}</span></div><span class="chip ${FS_CHIP[s.status]}">${t('fsSt_' + s.status)}</span></div>
      <div class="row" style="gap:6px;flex-wrap:wrap"><span class="chip pink">${esc(fsFund(s))}</span><span class="chip">${t('fsMinNote', {p: s.min_discount_pct})}</span><span class="chip">${t('fsPerB')}: ${s.per_buyer}</span>
        ${s.open_for_sellers ? `<span class="chip">${t('fsOpen')}</span>` : ''}${s.allow_vouchers ? '' : '<span class="chip">no voucher</span>'}</div>
      <p class="muted small" style="margin:0">${t('fsSold', {n: s.counts.sold, r: parseFloat(s.counts.revenue).toFixed(2), s: parseFloat(s.counts.subsidy).toFixed(2)})}</p>
      ${editable ? `<div class="row"><button type="button" id="fsEditBtn">${t('fsEdit')}</button><button type="button" class="danger" id="fsCancelBtn">${t('fsCancelS')}</button></div>` : ''}</section>`}
    <section class="card"><h2>${t('fsItems')}</h2>
      ${active.length ? `<div class="tbl-wrap"><table class="fs-table">${head}<tbody>${active.map(row).join('')}</tbody></table></div>` : `<p class="muted">${t('fsNoItems')}</p>`}
      ${closed.length ? `<details><summary>${t('fsOthers', {n: closed.length})}</summary><div class="tbl-wrap"><table class="fs-table">${head}<tbody>${closed.map(row).join('')}</tbody></table></div></details>` : ''}</section>
    ${editable ? fsPickerHTML(s) : ''}`;
  $('#fsBack').onclick = e => { e.preventDefault(); FS.open = null; FS.form = null; renderPromos(); };
  if (FS.form) bindFsForm(() => { FS.form = null; drawFlashDetail(id, FS.data); });
  if ($('#fsEditBtn')) $('#fsEditBtn').onclick = () => { FS.form = Object.assign({}, s); drawFlashDetail(id, FS.data); };
  if ($('#fsCancelBtn')) $('#fsCancelBtn').onclick = async () => {
    if (!confirm(t('fsCancelConfirm'))) return;
    try { await api(`/admin/api/promos/flash/${id}/cancel`, {method: 'POST'}); toast(t('fsSaved')); drawFlashDetail(id); } catch (e) { toast(errMsg(e), 'error'); }
  };
  document.querySelectorAll('[data-fsdec]').forEach(b => b.onclick = async () => {
    const action = b.dataset.fsdec;
    let note = '';
    if (action === 'reject') { note = prompt(t('fsRejectWhy')) || ''; if (!note.trim()) return; }
    if (action === 'remove' && !confirm(t('fsRemoveConfirm'))) return;
    try { await api(`/admin/api/promos/flash/items/${b.dataset.id}/decide`, {method: 'POST', body: {action, note}}); toast(t('fsSaved')); drawFlashDetail(id); } catch (e) { toast(errMsg(e), 'error'); }
  });
  document.querySelectorAll('[data-fssave]').forEach(b => b.onclick = async () => {
    const tr = document.querySelector(`[data-row="${b.dataset.fssave}"]`), it = d.items.find(i => i.id === b.dataset.fssave);
    const price = tr.querySelector('[name="sale_price"]'), qty = tr.querySelector('[name="qty_cap"]');
    try {
      await api(`/admin/api/promos/flash/items/${it.id}`, {method: 'PUT', body: {variant_id: it.variant_id, sale_price: price ? price.value.trim() : it.sale_price, qty_cap: parseInt(qty.value, 10) || 0, per_buyer: it.per_buyer}});
      toast(t('fsSaved')); drawFlashDetail(id);
    } catch (e) { tr.querySelectorAll('.fs-in').forEach(x => x.classList.toggle('invalid', e.field === x.name)); toast(errMsg(e), 'error'); }
  });
  if (editable) bindPicker(s);
}

function fsPickerHTML(s) {
  const rows = FS.picked.map((v, n) => `<tr><td><b>${esc(v.product)}</b>${v.options ? ` <span class="muted">${esc(v.options)}</span>` : ''}<div class="muted small">${esc(v.store)} · ${esc(v.sku)}</div></td>
    <td>${fsMoney(v.price)}${phLowTag(v.low30)}</td><td><input class="fs-in" data-pk="${n}" name="sale_price" value="${esc(v.sale_price)}" inputmode="decimal"></td>
    <td><input class="fs-in sm" data-pk="${n}" name="qty_cap" value="${v.qty_cap}" inputmode="numeric"><div class="muted small">${t('fsStockLow', {n: v.stock})}</div></td>
    <td><button type="button" data-unpick="${n}" aria-label="remove">✕</button></td></tr>`).join('');
  return `<section class="card" id="fsPicker"><h2>${t('fsAdd')}</h2>
    ${s.phsar_share_pct < 100 ? `<div class="banner info">${t('fsInviteNote')}</div>` : ''}
    <label>${t('fsFind')}<input id="fsFind" type="search" autocomplete="off"></label>
    <div id="fsResults" class="fs-results"></div>
    ${FS.picked.length ? `<div class="row fs-bulk"><label>${t('fsBulkOff')}<input id="fsBulkOff" inputmode="numeric" value="${s.min_discount_pct}" style="max-width:110px"></label>
      <label>${t('fsBulkQty')}<input id="fsBulkQty" inputmode="numeric" value="10" style="max-width:110px"></label><button type="button" id="fsBulkApply">${t('fsApply')}</button></div>
    <div class="tbl-wrap"><table class="fs-table"><thead><tr><th>${t('fsPicked', {n: FS.picked.length})}</th><th>${t('fsNormal')}</th><th>${t('fsFlashP')}</th><th>${t('fsQuota')}</th><th></th></tr></thead><tbody>${rows}</tbody></table></div>
    <div class="row"><button type="button" class="primary" id="fsAddBtn">${t('fsAddN', {n: FS.picked.length})}</button></div>` : ''}</section>`;
}

function fsPriceAt(price, off) { return (Math.floor(parseFloat(price) * (100 - off)) / 100).toFixed(2); }

function bindPicker(s) {
  let timer = null;
  const keepInputs = () => document.querySelectorAll('#fsPicker [data-pk]').forEach(x => { FS.picked[+x.dataset.pk][x.name] = x.value; });
  const drawResults = () => {
    const have = new Set(FS.picked.map(v => v.id).concat(FS.data.items.filter(i => ['pending', 'invited', 'approved'].includes(i.status)).map(i => i.variant_id)));
    $('#fsResults').innerHTML = FS.results.map((v, i) => `<label class="fs-res"><input type="checkbox" data-res="${i}" ${have.has(v.id) ? 'checked disabled' : ''}>
      <span><b>${esc(v.product)}</b>${v.options ? ` <span class="muted">${esc(v.options)}</span>` : ''}<span class="muted small"> · ${esc(v.store)} · ${fsMoney(v.price)} · ${t('pmStock', {n: v.stock})}${v.low30 ? ' · ' + t('phLow30', {a: parseFloat(v.low30).toFixed(2)}) : ''}</span></span></label>`).join('');
    document.querySelectorAll('[data-res]').forEach(cb => cb.onchange = () => {
      if (!cb.checked) return;
      keepInputs();
      const v = FS.results[+cb.dataset.res];
      FS.picked.push(Object.assign({}, v, {sale_price: phSuggest(v.price, s.min_discount_pct, v.low30), qty_cap: Math.min(v.stock, 10)}));
      drawFlashDetail(s.id, FS.data).then(() => { $('#fsFind').value = FS.lastQ || ''; drawResults(); });
    });
  };
  $('#fsFind').oninput = () => {
    clearTimeout(timer);
    timer = setTimeout(async () => {
      const q = $('#fsFind').value.trim(); FS.lastQ = q;
      if (q.length < 2) { FS.results = []; $('#fsResults').innerHTML = ''; return; }
      try {
        // pricehistory: at = jam mulai sesi → low30 (harga terendah 30 hari) per varian
        const r = await api('/admin/api/promos/variants?q=' + encodeURIComponent(q) + '&at=' + encodeURIComponent(s.starts_at));
        FS.results = r.variants; PH.rule = r.low30_rule !== false; drawResults();
      } catch (e) { toast(errMsg(e), 'error'); }
    }, 250);
  };
  if (FS.results.length) drawResults();
  document.querySelectorAll('[data-unpick]').forEach(b => b.onclick = () => { keepInputs(); FS.picked.splice(+b.dataset.unpick, 1); drawFlashDetail(s.id, FS.data); });
  if ($('#fsBulkApply')) $('#fsBulkApply').onclick = () => {
    const off = parseInt($('#fsBulkOff').value, 10) || s.min_discount_pct, qty = parseInt($('#fsBulkQty').value, 10) || 1;
    FS.picked.forEach(v => { v.sale_price = phSuggest(v.price, off, v.low30); v.qty_cap = Math.min(qty, v.stock); });
    drawFlashDetail(s.id, FS.data);
  };
  if ($('#fsAddBtn')) $('#fsAddBtn').onclick = async () => {
    keepInputs();
    const items = FS.picked.map(v => ({variant_id: v.id, sale_price: String(v.sale_price).trim(), qty_cap: parseInt(v.qty_cap, 10) || 0}));
    try {
      const r = await api(`/admin/api/promos/flash/${s.id}/items`, {method: 'POST', body: {items}});
      FS.picked = []; FS.results = []; toast(t('fsSaved')); drawFlashDetail(s.id, r);
    } catch (e) {
      const m = /^items\.(\d+)\.(\w+)$/.exec(e.field || '');
      toast(errMsg(e), 'error');
      if (m) { const el = document.querySelector(`#fsPicker [data-pk="${m[1]}"][name="${m[2]}"]`); if (el) { el.classList.add('invalid'); el.focus(); } }
    }
  };
}
