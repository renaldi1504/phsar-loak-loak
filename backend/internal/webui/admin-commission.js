'use strict';
// Keuangan / super admin: aturan komisi — default, per kategori, per toko (dengan tanggal akhir, mis. promo 0% seller awal).
addStrings({
  en: {cmDesc: 'The most specific rule wins: store → category → default. Changes apply to new orders only; existing orders keep their rate.',
    cmDefault: 'Default commission', cmCats: 'Category rules', cmStores: 'Store rules', cmNow: 'Current', cmNewRate: 'New rate (%)', cmSave: 'Save',
    cmHistory: 'Default rate history', cmSince: 'since {d}', cmUntil: 'until {d}', cmBy: 'by {n}', cmUseDefault: 'Default ({r}%)', cmOverride: 'Rule',
    cmRemove: 'Remove', cmEnd: 'End now', cmSearch: 'Find a store by name', cmEndDate: 'End date (optional)', cmAdd: 'Add store rule',
    cmNoStores: 'No store rules — every store pays the category or default rate.', cmStoreCol: 'Store', cmRateCol: 'Rate', cmFromCol: 'From', cmToCol: 'Until',
    cmCatCol: 'Category', cmPick: 'Choose a store first', cmSaved: 'Saved', cmConfirmEnd: 'End this rule now? The store/category goes back to the rate above it.',
    cmPromoHint: 'Tip: give new sellers 0% for their first 3 months with an end date.', cmOpen: 'no end date', cmMax: 'Max 50%, two decimals.'},
  km: {cmDesc: 'ច្បាប់ជាក់លាក់បំផុតឈ្នះ៖ ហាង → ប្រភេទ → លំនាំដើម។ ការផ្លាស់ប្តូរអនុវត្តចំពោះការបញ្ជាទិញថ្មីប៉ុណ្ណោះ។',
    cmDefault: 'កម្រៃជើងសារលំនាំដើម', cmCats: 'ច្បាប់តាមប្រភេទ', cmStores: 'ច្បាប់តាមហាង', cmNow: 'បច្ចុប្បន្ន', cmNewRate: 'អត្រាថ្មី (%)', cmSave: 'រក្សាទុក',
    cmHistory: 'ប្រវត្តិអត្រាលំនាំដើម', cmSince: 'ចាប់ពី {d}', cmUntil: 'ដល់ {d}', cmBy: 'ដោយ {n}', cmUseDefault: 'លំនាំដើម ({r}%)', cmOverride: 'ច្បាប់',
    cmRemove: 'ដកចេញ', cmEnd: 'បញ្ចប់ឥឡូវ', cmSearch: 'ស្វែងរកហាងតាមឈ្មោះ', cmEndDate: 'ថ្ងៃបញ្ចប់ (ស្រេចចិត្ត)', cmAdd: 'បន្ថែមច្បាប់ហាង',
    cmNoStores: 'គ្មានច្បាប់ហាង — ហាងទាំងអស់បង់អត្រាប្រភេទ ឬលំនាំដើម។', cmStoreCol: 'ហាង', cmRateCol: 'អត្រា', cmFromCol: 'ចាប់ពី', cmToCol: 'ដល់',
    cmCatCol: 'ប្រភេទ', cmPick: 'សូមជ្រើសហាងជាមុន', cmSaved: 'បានរក្សាទុក', cmConfirmEnd: 'បញ្ចប់ច្បាប់នេះឥឡូវ?',
    cmPromoHint: 'គន្លឹះ៖ ផ្តល់ 0% ដល់អ្នកលក់ថ្មីរយៈពេល 3 ខែដំបូង ជាមួយថ្ងៃបញ្ចប់។', cmOpen: 'គ្មានថ្ងៃបញ្ចប់', cmMax: 'អតិបរមា 50%, ទសភាគពីរខ្ទង់។'},
  zh: {cmDesc: '最具体的规则优先：店铺 → 类目 → 默认。修改只对新订单生效，已有订单保留原佣金率。',
    cmDefault: '默认佣金', cmCats: '类目规则', cmStores: '店铺规则', cmNow: '当前', cmNewRate: '新佣金率 (%)', cmSave: '保存',
    cmHistory: '默认佣金率历史', cmSince: '自 {d}', cmUntil: '至 {d}', cmBy: '操作人 {n}', cmUseDefault: '默认 ({r}%)', cmOverride: '规则',
    cmRemove: '移除', cmEnd: '立即结束', cmSearch: '按名称查找店铺', cmEndDate: '结束日期（可选）', cmAdd: '添加店铺规则',
    cmNoStores: '暂无店铺规则——所有店铺按类目或默认佣金率。', cmStoreCol: '店铺', cmRateCol: '佣金率', cmFromCol: '开始', cmToCol: '结束',
    cmCatCol: '类目', cmPick: '请先选择店铺', cmSaved: '已保存', cmConfirmEnd: '立即结束此规则？店铺/类目将回到上一级佣金率。',
    cmPromoHint: '提示：可给新卖家前 3 个月 0% 佣金并设置结束日期。', cmOpen: '无结束日期', cmMax: '最高 50%，两位小数。'},
});

const CM = {data: null, store: null, results: []};
const pct = v => `${Number(v).toFixed(2).replace(/\.00$/, '')}%`;
const dshort = d => new Date(d).toLocaleDateString();

async function renderCommission() {
  CM.data = await api('/admin/api/commission');
  drawCommission();
}

function drawCommission() {
  const d = CM.data, act = d.active;
  const g = act.find(r => r.scope === 'global');
  const gRate = g ? g.rate_percent : '8';
  const catRule = id => act.find(r => r.scope === 'category' && r.category_id === id);
  const stores = act.filter(r => r.scope === 'store');
  const catName = c => (c.names[P.lang] || c.names.en || c.id);
  const cats = d.categories.map(c => {
    const r = catRule(c.id);
    return `<tr><td>${esc(catName(c))}</td><td>${r ? `<b>${pct(r.rate_percent)}</b> <span class="chip pink">${t('cmOverride')}</span>` : `<span class="muted">${t('cmUseDefault', {r: Number(gRate)})}</span>`}</td>
      <td><div style="display:flex;gap:8px;align-items:center;flex-wrap:wrap"><input class="cm-cat" data-cat="${esc(c.id)}" inputmode="decimal" placeholder="%" style="max-width:90px;min-height:36px">
      <button type="button" data-setcat="${esc(c.id)}">${t('cmSave')}</button>${r ? `<button type="button" class="link" data-end="${r.id}">${t('cmRemove')}</button>` : ''}</div></td></tr>`;
  }).join('');
  const storeRows = stores.length ? stores.map(r => `<tr><td><b>${esc(r.store_name)}</b></td><td>${pct(r.rate_percent)}</td><td class="muted">${dshort(r.valid_from)}</td>
      <td class="muted">${r.valid_to ? dshort(r.valid_to) : t('cmOpen')}</td><td><button type="button" data-end="${r.id}">${t('cmEnd')}</button></td></tr>`).join('')
    : `<tr><td colspan="5" class="muted">${t('cmNoStores')}</td></tr>`;
  const hist = d.global_history.map(r => `<li><b>${pct(r.rate_percent)}</b> <span class="muted">${t('cmSince', {d: dshort(r.valid_from)})}${r.valid_to ? ' · ' + t('cmUntil', {d: dshort(r.valid_to)}) : ''}${r.created_by ? ' · ' + t('cmBy', {n: esc(r.created_by)}) : ''}</span></li>`).join('');
  $('#app').innerHTML = shell(`${pageHead(t('i_commission'), t('cmDesc'))}
    <div class="kpis">
      <div class="kpi"><span>${t('cmDefault')}</span><b>${pct(gRate)}</b></div>
      <div class="kpi"><span>${t('cmCats')}</span><b>${act.filter(r => r.scope === 'category').length}</b></div>
      <div class="kpi"><span>${t('cmStores')}</span><b>${stores.length}</b></div>
    </div>
    <div class="cols">
      <section class="panel"><div class="ph"><h2>${t('cmDefault')}</h2><span class="muted">${g ? t('cmSince', {d: dshort(g.valid_from)}) : ''}</span></div>
        <div style="padding:16px 18px;display:flex;flex-direction:column;gap:12px">
          <div style="display:flex;gap:10px;align-items:flex-end;flex-wrap:wrap"><label>${t('cmNewRate')}<input id="cmGlobal" inputmode="decimal" value="${Number(gRate)}" style="max-width:120px"></label>
            <button class="primary" id="cmSaveGlobal" type="button">${t('cmSave')}</button></div>
          <p class="muted small">${t('cmMax')}</p>
          <div><b class="small">${t('cmHistory')}</b><ul style="margin:6px 0 0;padding-left:18px;display:flex;flex-direction:column;gap:4px;font-size:13px">${hist}</ul></div>
        </div></section>
      <section class="panel"><div class="ph"><h2>${t('cmCats')}</h2></div>
        <div class="tbl-wrap"><table style="min-width:0"><thead><tr><th>${t('cmCatCol')}</th><th>${t('cmRateCol')}</th><th></th></tr></thead><tbody>${cats}</tbody></table></div></section>
    </div>
    <section class="panel" id="cmStoreForm"><div class="ph"><h2>${t('cmStores')}</h2><span class="muted">${t('cmPromoHint')}</span></div>
      <div style="padding:16px 18px;display:flex;flex-wrap:wrap;gap:12px;align-items:flex-end;border-bottom:1px solid var(--line)">
        <label style="flex:1 1 240px;position:relative">${t('cmSearch')}<input id="cmQ" name="store_id" autocomplete="off" value="${esc(CM.store ? CM.store.name : '')}">
          <div id="cmResults" class="cm-results"></div></label>
        <label>${t('cmNewRate')}<input id="cmStoreRate" name="rate_percent" inputmode="decimal" placeholder="0" style="max-width:110px"></label>
        <label>${t('cmEndDate')}<input id="cmStoreEnd" name="valid_to" type="date" style="max-width:180px"></label>
        <button class="primary" id="cmAddStore" type="button">${t('cmAdd')}</button>
      </div>
      <div class="tbl-wrap"><table><thead><tr><th>${t('cmStoreCol')}</th><th>${t('cmRateCol')}</th><th>${t('cmFromCol')}</th><th>${t('cmToCol')}</th><th></th></tr></thead><tbody>${storeRows}</tbody></table></div>
    </section>`);

  const post = async (body, form) => {
    try { CM.data = await api('/admin/api/commission/rules', {method: 'POST', body}); toast(t('cmSaved')); CM.store = null; drawCommission(); }
    catch (e) {
      if (form && e.field) { const el = $(`#${form} [name="${e.field}"]`); if (el) { el.classList.add('invalid'); el.focus(); } }
      toast(e.field === 'rate_percent' ? t('cmMax') : errMsg(e), 'error');
    }
  };
  $('#cmSaveGlobal').onclick = () => post({scope: 'global', rate_percent: $('#cmGlobal').value.trim()});
  document.querySelectorAll('[data-setcat]').forEach(b => b.onclick = () => {
    const v = document.querySelector(`.cm-cat[data-cat="${b.dataset.setcat}"]`).value.trim();
    post({scope: 'category', category_id: b.dataset.setcat, rate_percent: v});
  });
  document.querySelectorAll('[data-end]').forEach(b => b.onclick = async () => {
    if (!confirm(t('cmConfirmEnd'))) return;
    try { CM.data = await api('/admin/api/commission/rules/' + b.dataset.end, {method: 'DELETE'}); toast(t('cmSaved')); drawCommission(); } catch (e) { toast(errMsg(e), 'error'); }
  });
  let timer;
  $('#cmQ').oninput = () => {
    CM.store = null; clearTimeout(timer);
    timer = setTimeout(async () => {
      const q = $('#cmQ').value.trim();
      if (!q) { $('#cmResults').innerHTML = ''; return; }
      const r = await api('/admin/api/commission/stores?q=' + encodeURIComponent(q));
      CM.results = r.stores;
      $('#cmResults').innerHTML = r.stores.map((s, i) => `<button type="button" data-pick="${i}">${esc(s.name)} <small class="muted">${esc(s.slug)}</small></button>`).join('');
      document.querySelectorAll('[data-pick]').forEach(b => b.onclick = () => { CM.store = CM.results[+b.dataset.pick]; $('#cmQ').value = CM.store.name; $('#cmResults').innerHTML = ''; });
    }, 250);
  };
  $('#cmAddStore').onclick = () => {
    document.querySelectorAll('#cmStoreForm .invalid').forEach(el => el.classList.remove('invalid'));
    if (!CM.store) { $('#cmQ').classList.add('invalid'); return toast(t('cmPick'), 'warning'); }
    const end = $('#cmStoreEnd').value;
    post({scope: 'store', store_id: CM.store.id, rate_percent: $('#cmStoreRate').value.trim(), valid_to: end ? new Date(end + 'T23:59:59+07:00').toISOString() : null}, 'cmStoreForm');
  };
}
