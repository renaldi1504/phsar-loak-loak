'use strict';
// Admin: voucher (dibiayai Phsar) & sesi flash sale (admin-flash.js). Izin commission.manage.
addStrings({
  en: {pmHint: 'Vouchers are paid by Phsar and typed by buyers at checkout. Flash sales run in sessions with many products; who pays the discount is chosen per session.',
    pmVouchers: 'Vouchers', pmFlash: 'Flash sales', pmCost30: 'Voucher cost, 30 days', pmUsed30: '{n} paid orders used a voucher', pmLiveV: 'Live vouchers', pmLiveF: 'Live flash sessions',
    pmNewV: 'New voucher', pmEditV: 'Edit voucher', pmCode: 'Code', pmCodeHint: 'Buyers type this, e.g. WELCOME5', pmName: 'Internal name (optional)', pmKind: 'Type',
    pmPercent: 'Percent off', pmFixed: 'Fixed amount ($)', pmValue: 'Value', pmMaxD: 'Max discount ($, optional)', pmMinS: 'Minimum spend on items ($)',
    pmStart: 'Starts', pmEnd: 'Ends (optional)', pmTotal: 'Total uses (empty = no limit)', pmPerUser: 'Uses per buyer', pmFirst: 'First order only (signed-in buyers, not guests)',
    pmActive: 'Active', pmSave: 'Save', pmCancel: 'Cancel', pmSaved: 'Saved', pmOff: 'Turn off', pmOn: 'Turn on', pmEdit: 'Edit', pmUsage: '{u} used · {h} waiting payment',
    pmCost: 'cost ${c}', pmNoV: 'No vouchers yet.', pmNoF: 'No flash sales yet.', pmNewF: 'New flash sale', pmEditF: 'Edit flash sale', pmFind: 'Find product (name, SKU or store)',
    pmPick: 'Choose', pmPicked: 'Selected', pmNormal: 'normal ${p}', pmSale: 'Flash price ($)', pmCap: 'Quantity limit (empty = until stock runs out)', pmSold: '{s} sold',
    pmSoldCap: '{s} / {c} sold', pmEndNow: 'End now', pmEndConfirm: 'End this flash sale now?', pmStock: 'stock {n}',
    pmSt_live: 'Live', pmSt_scheduled: 'Scheduled', pmSt_ended: 'Ended', pmSt_off: 'Off', pmOffP: '{v}% off', pmOffF: '${v} off', pmMin: 'min ${m}', pmMaxCap: 'max ${m}', pmFirstChip: 'first order'},
  km: {pmHint: 'ប័ណ្ណត្រូវបាន Phsar ចេញថ្លៃ ហើយអ្នកទិញបញ្ចូលនៅពេលបង់ប្រាក់។ Flash sale ដំណើរការជាវគ្គ មានផលិតផលច្រើន; អ្នកចេញថ្លៃបញ្ចុះត្រូវជ្រើសតាមវគ្គ។',
    pmVouchers: 'ប័ណ្ណបញ្ចុះតម្លៃ', pmFlash: 'Flash sale', pmCost30: 'ថ្លៃប័ណ្ណ ៣០ ថ្ងៃ', pmUsed30: 'ការបញ្ជាទិញបានបង់ {n} បានប្រើប័ណ្ណ', pmLiveV: 'ប័ណ្ណកំពុងដំណើរការ', pmLiveF: 'Flash sale កំពុងដំណើរការ',
    pmNewV: 'ប័ណ្ណថ្មី', pmEditV: 'កែប័ណ្ណ', pmCode: 'លេខកូដ', pmCodeHint: 'អ្នកទិញវាយលេខនេះ ឧ. WELCOME5', pmName: 'ឈ្មោះខាងក្នុង (ស្រេចចិត្ត)', pmKind: 'ប្រភេទ',
    pmPercent: 'បញ្ចុះភាគរយ', pmFixed: 'ចំនួនថេរ ($)', pmValue: 'តម្លៃ', pmMaxD: 'បញ្ចុះអតិបរមា ($, ស្រេចចិត្ត)', pmMinS: 'ទិញអប្បបរមា ($)',
    pmStart: 'ចាប់ផ្តើម', pmEnd: 'បញ្ចប់ (ស្រេចចិត្ត)', pmTotal: 'ចំនួនប្រើសរុប (ទទេ = គ្មានកំណត់)', pmPerUser: 'ចំនួនប្រើក្នុងម្នាក់', pmFirst: 'ការបញ្ជាទិញលើកដំបូងប៉ុណ្ណោះ (អ្នកបានចូល មិនមែនភ្ញៀវ)',
    pmActive: 'សកម្ម', pmSave: 'រក្សាទុក', pmCancel: 'បោះបង់', pmSaved: 'បានរក្សាទុក', pmOff: 'បិទ', pmOn: 'បើក', pmEdit: 'កែ', pmUsage: 'បានប្រើ {u} · រង់ចាំបង់ {h}',
    pmCost: 'ថ្លៃ ${c}', pmNoV: 'មិនទាន់មានប័ណ្ណ។', pmNoF: 'មិនទាន់មាន flash sale។', pmNewF: 'Flash sale ថ្មី', pmEditF: 'កែ flash sale', pmFind: 'ស្វែងរកផលិតផល (ឈ្មោះ SKU ឬហាង)',
    pmPick: 'ជ្រើស', pmPicked: 'បានជ្រើស', pmNormal: 'ធម្មតា ${p}', pmSale: 'តម្លៃ flash ($)', pmCap: 'កំណត់ចំនួន (ទទេ = រហូតអស់ស្តុក)', pmSold: 'លក់បាន {s}',
    pmSoldCap: 'លក់បាន {s} / {c}', pmEndNow: 'បញ្ចប់ឥឡូវ', pmEndConfirm: 'បញ្ចប់ flash sale នេះឥឡូវ?', pmStock: 'ស្តុក {n}',
    pmSt_live: 'កំពុងដំណើរការ', pmSt_scheduled: 'បានកំណត់ពេល', pmSt_ended: 'បានបញ្ចប់', pmSt_off: 'បិទ', pmOffP: 'បញ្ចុះ {v}%', pmOffF: 'បញ្ចុះ ${v}', pmMin: 'អប្បបរមា ${m}', pmMaxCap: 'អតិបរមា ${m}', pmFirstChip: 'លើកដំបូង'},
  zh: {pmHint: '优惠券由 Phsar 承担，买家结账时输入。限时抢购按场次进行，每场可有多个商品；折扣由谁承担按场次设置。',
    pmVouchers: '优惠券', pmFlash: '限时抢购', pmCost30: '优惠券成本（30 天）', pmUsed30: '{n} 笔已付款订单使用了优惠券', pmLiveV: '进行中的优惠券', pmLiveF: '进行中的抢购',
    pmNewV: '新建优惠券', pmEditV: '编辑优惠券', pmCode: '优惠码', pmCodeHint: '买家输入的代码，例如 WELCOME5', pmName: '内部名称（可选）', pmKind: '类型',
    pmPercent: '百分比折扣', pmFixed: '固定金额（$）', pmValue: '数值', pmMaxD: '最高折扣（$，可选）', pmMinS: '商品最低消费（$）',
    pmStart: '开始', pmEnd: '结束（可选）', pmTotal: '总使用次数（空 = 不限）', pmPerUser: '每位买家可用次数', pmFirst: '仅限首单（已登录买家，不含游客）',
    pmActive: '启用', pmSave: '保存', pmCancel: '取消', pmSaved: '已保存', pmOff: '停用', pmOn: '启用', pmEdit: '编辑', pmUsage: '已用 {u} · 待付款 {h}',
    pmCost: '成本 ${c}', pmNoV: '暂无优惠券。', pmNoF: '暂无限时抢购。', pmNewF: '新建限时抢购', pmEditF: '编辑限时抢购', pmFind: '查找商品（名称、SKU 或店铺）',
    pmPick: '选择', pmPicked: '已选择', pmNormal: '原价 ${p}', pmSale: '抢购价（$）', pmCap: '数量上限（空 = 直到售完）', pmSold: '已售 {s}',
    pmSoldCap: '已售 {s} / {c}', pmEndNow: '立即结束', pmEndConfirm: '立即结束此限时抢购？', pmStock: '库存 {n}',
    pmSt_live: '进行中', pmSt_scheduled: '已排期', pmSt_ended: '已结束', pmSt_off: '已停用', pmOffP: '减 {v}%', pmOffF: '减 ${v}', pmMin: '满 ${m}', pmMaxCap: '最高 ${m}', pmFirstChip: '首单'},
});

const PM = {tab: 'vouchers', data: null, edit: null, picked: null};
const PM_CHIP = {live: 'ok', scheduled: 'warn', ended: '', off: 'bad'};

function pmLocal(d) { // ISO → nilai input datetime-local (waktu lokal)
  if (!d) return '';
  const x = new Date(d); x.setMinutes(x.getMinutes() - x.getTimezoneOffset());
  return x.toISOString().slice(0, 16);
}
function pmISO(v) { return v ? new Date(v).toISOString() : null; }
function pmWhen(d) { return d ? new Date(d).toLocaleString([], {dateStyle: 'medium', timeStyle: 'short'}) : '∞'; }

async function renderPromos() {
  const [d, f] = await Promise.all([api('/admin/api/promos'), api('/admin/api/promos/flash')]);
  PM.data = d; PM.data.sessions = f.sessions;
  drawPromos();
}

function drawPromos() {
  const d = PM.data, live = x => x.state === 'live';
  const tabs = [['vouchers', 'pmVouchers', d.vouchers.filter(live).length], ['flash', 'pmFlash', d.sessions.filter(x => x.status === 'live').length]]
    .map(([k, l, n]) => `<button type="button" aria-selected="${PM.tab === k}" data-pmtab="${k}">${t(l)}${n ? ` <span class="count">${n}</span>` : ''}</button>`).join('');
  $('#app').innerHTML = shell(`${pageHead(t('i_promos'), t('pmHint'))}
    <section class="kpis">
      <div class="kpi"><span>${t('pmCost30')}</span><b>$${esc(d.totals.cost_30d)}</b><span>${t('pmUsed30', {n: d.totals.used_30d})}</span></div>
      <div class="kpi"><span>${t('pmLiveV')}</span><b>${d.vouchers.filter(live).length}</b><span>&nbsp;</span></div>
      <div class="kpi"><span>${t('pmLiveF')}</span><b>${d.sessions.filter(x => x.status === 'live').length}</b><span>${t('fsUpcoming', {n: d.sessions.filter(x => x.status === 'scheduled').length})}</span></div>
    </section>
    <div class="tabs count">${tabs}</div>
    <div id="pmBody"></div>`);
  document.querySelectorAll('[data-pmtab]').forEach(b => b.onclick = () => { PM.tab = b.dataset.pmtab; PM.edit = null; PM.picked = null; drawPromos(); });
  if (PM.tab === 'vouchers') drawVouchers(); else drawFlashSessions();
}

function pmFail(e, box) {
  if (e.field) { const el = $(`#${box} [name="${e.field}"]`); if (el) { el.classList.add('invalid'); el.focus(); } }
  toast(errMsg(e), 'error');
}

function drawVouchers() {
  const d = PM.data, v = PM.edit;
  const rows = d.vouchers.map(x => `<section class="card"><div class="spread"><div><h2 style="font-family:ui-monospace,Menlo,monospace">${esc(x.code)}</h2>
      ${x.name ? `<span class="muted">${esc(x.name)}</span>` : ''}</div><span class="chip ${PM_CHIP[x.state]}">${t('pmSt_' + x.state)}</span></div>
    <div class="row" style="gap:6px;flex-wrap:wrap"><span class="chip pink">${x.kind === 'percent' ? t('pmOffP', {v: parseFloat(x.value)}) : t('pmOffF', {v: x.value})}</span>
      ${x.max_discount ? `<span class="chip">${t('pmMaxCap', {m: x.max_discount})}</span>` : ''}${parseFloat(x.min_spend) > 0 ? `<span class="chip">${t('pmMin', {m: x.min_spend})}</span>` : ''}
      ${x.first_order_only ? `<span class="chip">${t('pmFirstChip')}</span>` : ''}</div>
    <p class="muted small" style="margin:0">${pmWhen(x.starts_at)} → ${pmWhen(x.ends_at)} · ${t('pmUsage', {u: x.used + (x.total_limit ? ' / ' + x.total_limit : ''), h: x.held})} · ${t('pmCost', {c: x.cost})}</p>
    <div class="row"><button type="button" data-vedit="${x.id}">${t('pmEdit')}</button><button type="button" class="${x.active ? 'danger' : ''}" data-vact="${x.id}" data-on="${!x.active}">${t(x.active ? 'pmOff' : 'pmOn')}</button></div></section>`).join('');
  const f = v || {kind: 'percent', per_user_limit: 1, active: true, min_spend: '0'};
  $('#pmBody').innerHTML = `<section class="card" id="vForm"><h2>${t(v && v.id ? 'pmEditV' : 'pmNewV')}</h2>
    <div class="grid">
      <label>${t('pmCode')}<input name="code" maxlength="24" value="${esc(f.code || '')}" style="text-transform:uppercase"><small class="muted">${t('pmCodeHint')}</small></label>
      <label>${t('pmName')}<input name="name" maxlength="80" value="${esc(f.name || '')}"></label>
      <label>${t('pmKind')}<select name="kind"><option value="percent" ${f.kind === 'percent' ? 'selected' : ''}>${t('pmPercent')}</option><option value="fixed" ${f.kind === 'fixed' ? 'selected' : ''}>${t('pmFixed')}</option></select></label>
      <label>${t('pmValue')}<input name="value" inputmode="decimal" value="${esc(f.value ? String(parseFloat(f.value)) : '')}"></label>
      <label>${t('pmMaxD')}<input name="max_discount" inputmode="decimal" value="${esc(f.max_discount || '')}"></label>
      <label>${t('pmMinS')}<input name="min_spend" inputmode="decimal" value="${esc(f.min_spend || '0')}"></label>
      <label>${t('pmStart')}<input name="starts_at" type="datetime-local" value="${pmLocal(f.starts_at)}"></label>
      <label>${t('pmEnd')}<input name="ends_at" type="datetime-local" value="${pmLocal(f.ends_at)}"></label>
      <label>${t('pmTotal')}<input name="total_limit" inputmode="numeric" value="${f.total_limit || ''}"></label>
      <label>${t('pmPerUser')}<input name="per_user_limit" inputmode="numeric" value="${f.per_user_limit || 1}"></label>
    </div>
    <label class="inline"><input type="checkbox" name="first_order_only" ${f.first_order_only ? 'checked' : ''}> ${t('pmFirst')}</label>
    <label class="inline"><input type="checkbox" name="active" ${f.active ? 'checked' : ''}> ${t('pmActive')}</label>
    <div class="row"><button type="button" class="primary" id="vSave">${t('pmSave')}</button>${v ? `<button type="button" id="vCancel">${t('pmCancel')}</button>` : ''}</div></section>
    ${rows || `<section class="card"><p class="muted">${t('pmNoV')}</p></section>`}`;
  const val = n => $(`#vForm [name="${n}"]`);
  $('#vSave').onclick = async () => {
    document.querySelectorAll('#vForm .invalid').forEach(el => el.classList.remove('invalid'));
    const tl = val('total_limit').value.trim();
    const body = {code: val('code').value, name: val('name').value, kind: val('kind').value, value: val('value').value.trim(), max_discount: val('max_discount').value.trim(),
      min_spend: val('min_spend').value.trim(), starts_at: pmISO(val('starts_at').value), ends_at: pmISO(val('ends_at').value),
      total_limit: tl ? parseInt(tl, 10) : null, per_user_limit: parseInt(val('per_user_limit').value, 10) || 1,
      first_order_only: val('first_order_only').checked, active: val('active').checked};
    try {
      await api(v && v.id ? `/admin/api/promos/vouchers/${v.id}` : '/admin/api/promos/vouchers', {method: v && v.id ? 'PUT' : 'POST', body});
      PM.edit = null; toast(t('pmSaved')); renderPromos();
    } catch (e) { pmFail(e, 'vForm'); }
  };
  if ($('#vCancel')) $('#vCancel').onclick = () => { PM.edit = null; drawVouchers(); };
  document.querySelectorAll('[data-vedit]').forEach(b => b.onclick = () => { PM.edit = d.vouchers.find(x => x.id === b.dataset.vedit); drawVouchers(); $('#vForm').scrollIntoView({behavior: 'smooth'}); });
  document.querySelectorAll('[data-vact]').forEach(b => b.onclick = async () => {
    try { await api(`/admin/api/promos/vouchers/${b.dataset.vact}/active`, {method: 'POST', body: {active: b.dataset.on === 'true'}}); toast(t('pmSaved')); renderPromos(); }
    catch (e) { toast(errMsg(e), 'error'); }
  });
}

