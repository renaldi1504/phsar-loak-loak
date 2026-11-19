'use strict';
// Admin: ongkir per zona, wilayah layanan, biaya kurir, kurs riel (izin commission.manage).
addStrings({
  en: {shDesc: 'Delivery fee is charged once per store in a checkout. 0 = free delivery. Changes apply to new checkouts only.',
    shZones: 'Delivery zones', shZone: 'Zone', shCode: 'Code', shNameEn: 'Name (English)', shNameKm: 'Name (Khmer)', shNameZh: 'Name (Chinese)',
    shFee: 'Fee per store ($)', shFreeOver: 'Free when store subtotal ≥ ($)', shFreeHint: 'Leave empty for no free-delivery threshold.',
    shCourier: 'Paid to courier per delivery ($)', shCourierHint: 'For cost reporting only — buyers never see this.', shSave: 'Save', shAddZone: 'Add zone',
    shAreas: 'Service areas', shAreasDesc: 'Buyers can only check out inside an active area. The area decides which zone fee applies.',
    shActive: 'Active', shInactive: 'Inactive', shBoundary: 'Boundary (GeoJSON Polygon, [lng, lat])', shStores: '{n} stores', shAddArea: 'Add area',
    shFx: 'Riel exchange rate', shFxDesc: 'Used for the “≈ riel” amounts on the website and app. Prices and payments stay in US dollars.',
    shFxNow: 'Current rate', shFxNew: 'New rate (riel per $1)', shFxNote: 'Note (optional, e.g. NBC 16 Nov)', shFxSave: 'Save new rate', shFxHist: 'History',
    shCost: 'Courier cost, last 30 days', shCostSub: '{n} deliveries · buyers paid {f} in delivery fees', shFreeNow: 'Free delivery', shSaved: 'Saved',
    shSrc_manual: 'Manual', shSrc_nbc: 'NBC', shSrc_payway: 'PayWay', err_last_area: 'At least one service area must stay active.', err_exists: 'This code is already used.'},
  km: {shDesc: 'ថ្លៃដឹកគិតម្តងក្នុងមួយហាងក្នុងការបញ្ជាទិញ។ 0 = ដឹកឥតគិតថ្លៃ។ ការផ្លាស់ប្តូរអនុវត្តចំពោះការបញ្ជាទិញថ្មីប៉ុណ្ណោះ។',
    shZones: 'តំបន់ដឹកជញ្ជូន', shZone: 'តំបន់', shCode: 'កូដ', shNameEn: 'ឈ្មោះ (អង់គ្លេស)', shNameKm: 'ឈ្មោះ (ខ្មែរ)', shNameZh: 'ឈ្មោះ (ចិន)',
    shFee: 'ថ្លៃក្នុងមួយហាង ($)', shFreeOver: 'ឥតគិតថ្លៃ ពេលសរុបរងហាង ≥ ($)', shFreeHint: 'ទុកទទេ ប្រសិនបើគ្មានកម្រិតដឹកឥតគិតថ្លៃ។',
    shCourier: 'បង់ឱ្យអ្នកដឹកក្នុងមួយការដឹក ($)', shCourierHint: 'សម្រាប់របាយការណ៍ចំណាយប៉ុណ្ណោះ — អ្នកទិញមិនឃើញទេ។', shSave: 'រក្សាទុក', shAddZone: 'បន្ថែមតំបន់',
    shAreas: 'តំបន់សេវាកម្ម', shAreasDesc: 'អ្នកទិញអាចបញ្ជាទិញបានតែក្នុងតំបន់សកម្ម។ តំបន់កំណត់ថ្លៃដឹកដែលត្រូវអនុវត្ត។',
    shActive: 'សកម្ម', shInactive: 'អសកម្ម', shBoundary: 'ព្រំដែន (GeoJSON Polygon, [lng, lat])', shStores: 'ហាង {n}', shAddArea: 'បន្ថែមតំបន់សេវាកម្ម',
    shFx: 'អត្រាប្តូរប្រាក់រៀល', shFxDesc: 'ប្រើសម្រាប់ចំនួន «≈ រៀល» នៅលើគេហទំព័រ និងកម្មវិធី។ តម្លៃ និងការបង់ប្រាក់នៅតែជាដុល្លារ។',
    shFxNow: 'អត្រាបច្ចុប្បន្ន', shFxNew: 'អត្រាថ្មី (រៀលក្នុង $1)', shFxNote: 'កំណត់ចំណាំ (ស្រេចចិត្ត ឧ. NBC 16 វិច្ឆិកា)', shFxSave: 'រក្សាទុកអត្រាថ្មី', shFxHist: 'ប្រវត្តិ',
    shCost: 'ចំណាយអ្នកដឹក ៣០ ថ្ងៃចុងក្រោយ', shCostSub: 'ការដឹក {n} · អ្នកទិញបានបង់ថ្លៃដឹក {f}', shFreeNow: 'ដឹកឥតគិតថ្លៃ', shSaved: 'បានរក្សាទុក',
    shSrc_manual: 'ដោយដៃ', shSrc_nbc: 'NBC', shSrc_payway: 'PayWay', err_last_area: 'ត្រូវមានតំបន់សេវាកម្មសកម្មយ៉ាងហោចណាស់មួយ។', err_exists: 'កូដនេះត្រូវបានប្រើរួចហើយ។'},
  zh: {shDesc: '每次结账按店铺收取一次运费。0 = 免运费。修改只对新订单生效。',
    shZones: '配送区域', shZone: '区域', shCode: '代码', shNameEn: '名称（英文）', shNameKm: '名称（高棉文）', shNameZh: '名称（中文）',
    shFee: '每店运费（$）', shFreeOver: '单店小计满（$）免运费', shFreeHint: '留空表示没有免运费门槛。',
    shCourier: '每单支付给快递员（$）', shCourierHint: '仅用于成本报表，买家看不到。', shSave: '保存', shAddZone: '添加区域',
    shAreas: '服务范围', shAreasDesc: '买家只能在启用的服务范围内下单。服务范围决定适用哪个区域的运费。',
    shActive: '启用', shInactive: '停用', shBoundary: '边界（GeoJSON Polygon，[经度, 纬度]）', shStores: '{n} 家店铺', shAddArea: '添加服务范围',
    shFx: '瑞尔汇率', shFxDesc: '用于网站和 App 上“≈ 瑞尔”的显示。价格和付款仍以美元计。',
    shFxNow: '当前汇率', shFxNew: '新汇率（1 美元兑瑞尔）', shFxNote: '备注（可选，例如 NBC 11月16日）', shFxSave: '保存新汇率', shFxHist: '历史',
    shCost: '近 30 天快递成本', shCostSub: '{n} 单 · 买家支付运费 {f}', shFreeNow: '免运费', shSaved: '已保存',
    shSrc_manual: '手动', shSrc_nbc: 'NBC', shSrc_payway: 'PayWay', err_last_area: '至少要保留一个启用的服务范围。', err_exists: '此代码已被使用。'},
});

const SHP = {data: null};

async function renderShipping() {
  SHP.data = await api('/admin/api/shipping');
  drawShipping();
}

function drawShipping() {
  const d = SHP.data, nf = n => Number(n).toLocaleString('en-US');
  const zoneForm = (z, isNew) => `<form class="card shz" data-zone="${isNew ? '' : esc(z.code)}" onsubmit="return false">
      <div class="spread"><h3 style="margin:0">${isNew ? t('shAddZone') : esc((z.name && z.name.en) || z.code)}</h3>${isNew ? '' : `<span class="chip ${Number(z.fee_per_store) > 0 ? '' : 'ok'}">${Number(z.fee_per_store) > 0 ? '$' + esc(z.fee_per_store) : t('shFreeNow')}</span>`}</div>
      <div class="grid">
        ${isNew ? `<label>${t('shCode')}<input name="code" maxlength="32" placeholder="kampot"></label>` : ''}
        <label>${t('shNameEn')}<input name="name_en" maxlength="60" value="${esc((z.name || {}).en || '')}"></label>
        <label>${t('shNameKm')}<input name="name_km" maxlength="60" value="${esc((z.name || {}).km || '')}"></label>
        <label>${t('shNameZh')}<input name="name_zh" maxlength="60" value="${esc((z.name || {}).zh || '')}"></label>
        <label>${t('shFee')}<input name="fee_per_store" inputmode="decimal" value="${esc(z.fee_per_store || '0.00')}"></label>
        <label>${t('shFreeOver')}<input name="free_over" inputmode="decimal" value="${esc(z.free_over || '')}" placeholder="—"><small class="muted">${t('shFreeHint')}</small></label>
        <label>${t('shCourier')}<input name="courier_fee" inputmode="decimal" value="${esc(z.courier_fee || '0.00')}"><small class="muted">${t('shCourierHint')}</small></label>
      </div>
      <div><button type="button" class="primary" data-savezone>${isNew ? t('shAddZone') : t('shSave')}</button></div></form>`;
  const zoneOpts = sel => d.zones.map(z => `<option value="${esc(z.code)}" ${z.code === sel ? 'selected' : ''}>${esc((z.name && z.name.en) || z.code)}</option>`).join('');
  const areaForm = (a, isNew) => `<form class="card sha" data-area="${isNew ? '' : esc(a.code)}" onsubmit="return false">
      <div class="spread"><h3 style="margin:0">${isNew ? t('shAddArea') : esc((a.name && a.name.en) || a.code)}</h3>${isNew ? '' : `<span class="chip ${a.active ? 'ok' : ''}">${t(a.active ? 'shActive' : 'shInactive')}</span>`}</div>
      ${isNew ? '' : `<p class="muted small" style="margin:0">${t('shStores', {n: a.stores})} · <code>${esc(a.code)}</code></p>`}
      <div class="grid">
        ${isNew ? `<label>${t('shCode')}<input name="code" maxlength="32" placeholder="otres"></label>` : ''}
        <label>${t('shNameEn')}<input name="name_en" maxlength="60" value="${esc((a.name || {}).en || '')}"></label>
        <label>${t('shNameKm')}<input name="name_km" maxlength="60" value="${esc((a.name || {}).km || '')}"></label>
        <label>${t('shNameZh')}<input name="name_zh" maxlength="60" value="${esc((a.name || {}).zh || '')}"></label>
        <label>${t('shZone')}<select name="zone">${zoneOpts(a.zone)}</select></label>
        <label class="inline" style="align-self:center"><input type="checkbox" name="active" ${a.active || isNew ? 'checked' : ''}> ${t('shActive')}</label>
      </div>
      <label>${t('shBoundary')}<textarea name="boundary" rows="${isNew ? 4 : 3}" spellcheck="false" style="font-family:ui-monospace,Menlo,monospace;font-size:12px">${esc(a.boundary ? JSON.stringify(a.boundary) : '')}</textarea></label>
      <div><button type="button" class="primary" data-savearea>${isNew ? t('shAddArea') : t('shSave')}</button></div></form>`;
  $('#app').innerHTML = shell(`${pageHead(t('i_shipping'), t('shDesc'))}
    <section class="kpis">
      <div class="kpi"><span>${t('shFxNow')}</span><b>${nf(d.riel_per_usd)} ៛</b><span>$1</span></div>
      <div class="kpi"><span>${t('shCost')}</span><b>$${esc(d.courier_30d.cost)}</b><span>${esc(t('shCostSub', {n: d.courier_30d.deliveries, f: '$' + d.courier_30d.fees_collected}))}</span></div>
    </section>
    <section class="panel"><div class="ph"><h2>${t('shZones')}</h2></div><div style="padding:16px 18px;display:flex;flex-direction:column;gap:14px">
      ${d.zones.map(z => zoneForm(z, false)).join('')}
      <details><summary>${t('shAddZone')}</summary>${zoneForm({}, true)}</details></div></section>
    <section class="panel"><div class="ph"><h2>${t('shAreas')}</h2></div><div style="padding:16px 18px;display:flex;flex-direction:column;gap:14px">
      <p class="muted small" style="margin:0">${t('shAreasDesc')}</p>
      ${d.areas.map(a => areaForm(a, false)).join('')}
      <details><summary>${t('shAddArea')}</summary>${areaForm({zone: (d.zones[0] || {}).code}, true)}</details></div></section>
    <section class="panel" id="shFx"><div class="ph"><h2>${t('shFx')}</h2></div><div style="padding:16px 18px;display:flex;flex-direction:column;gap:14px">
      <p class="muted small" style="margin:0">${t('shFxDesc')}</p>
      <div class="row" style="align-items:flex-end;flex-wrap:wrap">
        <label>${t('shFxNew')}<input id="fxRate" name="riel_per_usd" inputmode="decimal" placeholder="4100" style="max-width:160px"></label>
        <label style="flex:1;min-width:200px">${t('shFxNote')}<input id="fxNote" name="note" maxlength="200"></label>
        <button type="button" class="primary" id="fxSave">${t('shFxSave')}</button></div>
      <details><summary>${t('shFxHist')} (${d.rates.length})</summary><div class="tbl-wrap"><table style="min-width:0;width:100%;margin-top:8px"><tbody>
        ${d.rates.map(r => `<tr><td><b>${nf(r.riel_per_usd)} ៛</b></td><td>${esc(r.rate_date)}</td><td>${t('shSrc_' + r.source)}</td><td class="muted">${esc(r.by || '')}${r.note ? ' · ' + esc(r.note) : ''}</td><td class="muted small">${new Date(r.at).toLocaleString()}</td></tr>`).join('')}
      </tbody></table></div></details></div></section>`);

  const val = (f, n) => (f.querySelector(`[name="${n}"]`) || {}).value || '';
  const names = f => ({en: val(f, 'name_en'), km: val(f, 'name_km'), zh: val(f, 'name_zh')});
  const fail = (e, f) => {
    f.querySelectorAll('.invalid').forEach(x => x.classList.remove('invalid'));
    const el = e.field && f.querySelector(`[name="${e.field}"]`);
    if (el) { el.classList.add('invalid'); el.focus(); }
    toast(errMsg(e), 'error');
  };
  document.querySelectorAll('[data-savezone]').forEach(b => b.onclick = async () => {
    const f = b.closest('form'), code = f.dataset.zone;
    const body = {code: code || val(f, 'code').trim(), name: names(f), fee_per_store: val(f, 'fee_per_store'), free_over: val(f, 'free_over'), courier_fee: val(f, 'courier_fee')};
    try { SHP.data = await api('/admin/api/shipping/zones' + (code ? '/' + encodeURIComponent(code) : ''), {method: code ? 'PUT' : 'POST', body}); toast(t('shSaved')); drawShipping(); }
    catch (e) { fail(e, f); }
  });
  document.querySelectorAll('[data-savearea]').forEach(b => b.onclick = async () => {
    const f = b.closest('form'), code = f.dataset.area;
    let boundary = null;
    try { boundary = JSON.parse(val(f, 'boundary')); } catch { fail({field: 'boundary', message: t('shBoundary')}, f); return; }
    const body = {code: code || val(f, 'code').trim(), name: names(f), zone: val(f, 'zone'), active: f.querySelector('[name="active"]').checked, boundary};
    try { SHP.data = await api('/admin/api/shipping/areas' + (code ? '/' + encodeURIComponent(code) : ''), {method: code ? 'PUT' : 'POST', body}); toast(t('shSaved')); drawShipping(); }
    catch (e) { fail(e, f); }
  });
  $('#fxSave').onclick = async () => {
    try { SHP.data = await api('/admin/api/shipping/fx', {method: 'POST', body: {riel_per_usd: $('#fxRate').value, note: $('#fxNote').value}}); toast(t('shSaved')); drawShipping(); }
    catch (e) { fail(e, $('#shFx')); }
  };
}
