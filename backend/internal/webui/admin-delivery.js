'use strict';
// Admin: kurir (staf Phsar) & pengantaran — tambah kurir, pantau tugas, tugaskan / pindahkan, lihat foto bukti.
addStrings({
  en: {delDesc: 'Parcels ready at the store, courier assignment and Phsar staff couriers.', kNeedCourier: 'ready at the store', kOnWay: 'On the way', kFailedSub: 'call the buyer or return to store', navDelivery: 'Delivery', delTitle: 'Delivery', couriersTitle: 'Couriers', addCourier: 'Add / update courier', cName: 'Name', cPhone: 'Phone',
    cVehicle: 'Vehicle', cPlate: 'Plate no.', cActive: 'Active', save: 'Save', saved: 'Saved', courierLink: 'Couriers log in at /courier with this phone number.',
    veh_motorbike: 'Motorbike', veh_tuktuk: 'Tuk-tuk', veh_car: 'Car', veh_bicycle: 'Bicycle', online: 'Online', offline: 'Offline', activeTasks: '{n} active', doneToday: '{n} today',
    d_open: 'In progress', d_failed: 'Failed', d_delivered: 'Delivered (24 h)', unassigned: 'Not assigned', assignTo: 'Assign to', unassign: '— back to list —',
    ds_unassigned: 'Waiting for courier', ds_assigned: 'Courier going to store', ds_picked_up: 'On the way', ds_delivered: 'Delivered', ds_failed: 'Failed',
    pickup: 'Pickup', dropoff: 'Deliver to', viewProof: 'Delivery photo', attempts: 'Attempts', inactive: 'Inactive',
    df_buyer_unreachable: 'Buyer not reachable', df_buyer_refused: 'Buyer refused', df_wrong_address: 'Wrong address', df_damaged: 'Package damaged', df_other: 'Other'},
  km: {delDesc: 'កញ្ចប់រួចរាល់នៅហាង ការចាត់អ្នកដឹក និងអ្នកដឹករបស់ Phsar។', kNeedCourier: 'រួចរាល់នៅហាង', kOnWay: 'កំពុងដឹក', kFailedSub: 'ទូរស័ព្ទអ្នកទិញ ឬប្រគល់ទៅហាងវិញ', navDelivery: 'ការដឹកជញ្ជូន', delTitle: 'ការដឹកជញ្ជូន', couriersTitle: 'អ្នកដឹក', addCourier: 'បន្ថែម / កែអ្នកដឹក', cName: 'ឈ្មោះ', cPhone: 'លេខទូរស័ព្ទ',
    cVehicle: 'យានជំនិះ', cPlate: 'ផ្លាកលេខ', cActive: 'សកម្ម', save: 'រក្សាទុក', saved: 'បានរក្សាទុក', courierLink: 'អ្នកដឹកចូលនៅ /courier ដោយប្រើលេខទូរស័ព្ទនេះ។',
    veh_motorbike: 'ម៉ូតូ', veh_tuktuk: 'តុកតុក', veh_car: 'ឡាន', veh_bicycle: 'កង់', online: 'អនឡាញ', offline: 'ក្រៅបណ្តាញ', activeTasks: 'កំពុងដឹក {n}', doneToday: 'ថ្ងៃនេះ {n}',
    d_open: 'កំពុងដំណើរការ', d_failed: 'ដឹកមិនបាន', d_delivered: 'បានដឹកដល់ (24 ម៉ោង)', unassigned: 'មិនទាន់ចាត់', assignTo: 'ចាត់ឱ្យ', unassign: '— ត្រឡប់ទៅបញ្ជី —',
    ds_unassigned: 'រង់ចាំអ្នកដឹក', ds_assigned: 'អ្នកដឹកកំពុងទៅហាង', ds_picked_up: 'កំពុងដឹក', ds_delivered: 'បានដឹកដល់', ds_failed: 'ដឹកមិនបាន',
    pickup: 'យកពី', dropoff: 'ដឹកទៅ', viewProof: 'រូបថតពេលប្រគល់', attempts: 'ចំនួនដង', inactive: 'អសកម្ម',
    df_buyer_unreachable: 'ទាក់ទងអ្នកទិញមិនបាន', df_buyer_refused: 'អ្នកទិញបដិសេធ', df_wrong_address: 'អាសយដ្ឋានខុស', df_damaged: 'កញ្ចប់ខូច', df_other: 'ផ្សេងៗ'},
  zh: {delDesc: '店铺已备好的包裹、分配快递员以及 Phsar 自有快递员。', kNeedCourier: '店铺已备好', kOnWay: '配送途中', kFailedSub: '联系买家或退回店铺', navDelivery: '配送', delTitle: '配送', couriersTitle: '快递员', addCourier: '添加 / 修改快递员', cName: '姓名', cPhone: '手机号',
    cVehicle: '交通工具', cPlate: '车牌号', cActive: '在职', save: '保存', saved: '已保存', courierLink: '快递员用此手机号在 /courier 登录。',
    veh_motorbike: '摩托车', veh_tuktuk: '嘟嘟车', veh_car: '汽车', veh_bicycle: '自行车', online: '在线', offline: '离线', activeTasks: '进行中 {n}', doneToday: '今日 {n}',
    d_open: '进行中', d_failed: '失败', d_delivered: '已送达（24 小时）', unassigned: '未分配', assignTo: '分配给', unassign: '— 退回列表 —',
    ds_unassigned: '等待快递员', ds_assigned: '快递员前往店铺', ds_picked_up: '配送中', ds_delivered: '已送达', ds_failed: '配送失败',
    pickup: '取件', dropoff: '送至', viewProof: '交付照片', attempts: '尝试次数', inactive: '已停用',
    df_buyer_unreachable: '联系不上买家', df_buyer_refused: '买家拒收', df_wrong_address: '地址错误', df_damaged: '包裹损坏', df_other: '其他'},
});

const AD = {status: 'open'};

async function renderDelivery() {
  if (AD.status === 'returns' && typeof crAdminReturns === 'function') return crAdminReturns(); // courierreturn: tab "Dikembalikan ke toko"
  const [{couriers}, {deliveries}] = await Promise.all([api('/admin/api/couriers'), api('/admin/api/deliveries?status=' + AD.status)]);
  const active = couriers.filter(c => c.active);
  const cRows = couriers.map(c => `<tr><td><b>${esc(c.name)}</b><br><span class="muted">${esc(c.phone)}</span></td>
    <td>${t('veh_' + c.vehicle)}${c.plate_no ? ' · ' + esc(c.plate_no) : ''}</td>
    <td>${!c.active ? chip('bad', t('inactive')) : c.online ? chip('ok', t('online')) : chip('', t('offline'))}</td>
    <td class="muted">${t('activeTasks', {n: c.active_tasks})} · ${t('doneToday', {n: c.done_today})}</td>
    <td><button type="button" class="small" data-edit='${esc(JSON.stringify({phone: c.phone, name: c.name, vehicle: c.vehicle, plate_no: c.plate_no, active: c.active}))}'>✎</button></td></tr>`).join('');
  const opts = (sel, status) => `<option value="">${status === 'unassigned' ? t('unassigned') : t('unassign')}</option>` +
    active.map(c => `<option value="${c.user_id}" ${c.user_id === sel ? 'selected' : ''}>${esc(c.name)}${c.online ? ' ●' : ''}</option>`).join('');
  const dCards = deliveries.map(d => `<section class="card"><div class="spread"><h2>${esc(d.order_no)}</h2>
      ${chip(d.status === 'failed' ? 'bad' : d.status === 'delivered' ? 'ok' : d.status === 'unassigned' ? 'warn' : 'pink', t('ds_' + d.status))}</div>
    <div class="grid"><p><span class="muted">${t('pickup')}</span><br><b>${esc(d.store_name)}</b> · ${esc(d.pickup_label)}<br>${esc(d.pickup_address)}<br>${esc(d.pickup_contact)} ${esc(d.pickup_phone)}</p>
      <p><span class="muted">${t('dropoff')}</span><br><b>${esc(d.recipient)}</b> · ${esc(d.phone)}<br>${esc(d.address_line)}${d.address_note ? `<br><span class="muted">${esc(d.address_note)}</span>` : ''}</p></div>
    ${d.status === 'failed' ? `<div class="banner bad">${t('df_' + d.fail_reason)}${d.fail_note ? ' — ' + esc(d.fail_note) : ''} · ${t('attempts')}: ${d.attempts}</div>` : ''}
    <div class="row">${d.status !== 'delivered' ? `<label class="inline">${t('assignTo')}<select data-assign="${d.id}" data-status="${d.status}">${opts(d.courier_id, d.status)}</select></label>`
      : `<span class="muted">${esc(d.courier_name)} · ${new Date(d.delivered_at).toLocaleString()}</span>`}
      ${d.has_proof ? `<button type="button" class="small" data-proof="${d.id}">📷 ${t('viewProof')}</button>` : ''}
      ${d.status === 'failed' && typeof crFailedBtn === 'function' ? crFailedBtn(d) : ''}</div></section>`).join('');
  const onWay = AD.status === 'open' ? deliveries.filter(d => d.status === 'assigned' || d.status === 'picked_up').length : null;
  $('#app').innerHTML = shell(`${pageHead(t('i_orders'), t('delDesc'))}
    <section class="kpis"><div class="kpi"><span>${t('ds_unassigned')}</span><b class="${SH.badges.deliveries ? 'warn' : ''}">${SH.badges.deliveries || 0}</b><span>${t('kNeedCourier')}</span></div>
      ${onWay != null ? `<div class="kpi"><span>${t('kOnWay')}</span><b>${onWay}</b><span>${t('ds_assigned')} · ${t('ds_picked_up')}</span></div>` : ''}
      <div class="kpi"><span>${t('ds_failed')}</span><b class="${SH.badges.failed ? 'warn' : ''}">${SH.badges.failed || 0}</b><span>${t('kFailedSub')}</span></div>
      <div class="kpi"><span>${t('couriersTitle')}</span><b>${couriers.filter(c => c.active && c.online).length} / ${active.length}</b><span>${t('online')}</span></div></section>
    <div class="tabs">${['open', 'failed', 'delivered', ...(typeof crAdminReturns === 'function' ? ['returns'] : [])].map(f => `<button type="button" aria-selected="${AD.status === f}" onclick="AD.status='${f}';render()">${t('d_' + f)}</button>`).join('')}</div>
    ${dCards || `<section class="card"><p class="muted">${t('empty')}</p></section>`}
    <h2>${t('couriersTitle')}</h2>
    <section class="card"><table class="tbl"><tbody>${cRows || `<tr><td class="muted">${t('empty')}</td></tr>`}</tbody></table></section>
    <section class="card"><h2>${t('addCourier')}</h2><p class="muted">${t('courierLink')}</p>
      <div class="grid"><label>${t('cPhone')}<input id="cPhone" inputmode="tel" placeholder="012 345 678"></label><label>${t('cName')}<input id="cName" maxlength="80"></label>
      <label>${t('cVehicle')}<select id="cVehicle">${['motorbike', 'tuktuk', 'car', 'bicycle'].map(v => `<option value="${v}">${t('veh_' + v)}</option>`).join('')}</select></label>
      <label>${t('cPlate')}<input id="cPlate" maxlength="20"></label></div>
      <label class="inline"><input type="checkbox" id="cActive" checked> ${t('cActive')}</label>
      <div class="row"><button type="button" class="primary" id="cSave">${t('save')}</button></div></section>`);
  document.querySelectorAll('[data-assign]').forEach(s => s.onchange = async () => {
    try { await api(`/admin/api/deliveries/${s.dataset.assign}/assign`, {method: 'POST', body: {courier_id: s.value}}); toast(t('saved')); render(); }
    catch (e) { toast(errMsg(e), 'error'); render(); }
  });
  if (typeof crBindFailed === 'function') crBindFailed(active); // courierreturn
  document.querySelectorAll('[data-proof]').forEach(b => b.onclick = async () => {
    const r = await fetch(`/admin/api/deliveries/${b.dataset.proof}/proof`, {headers: {Authorization: 'Bearer ' + P.token}});
    if (!r.ok) { toast(r.statusText, 'error'); return; }
    window.open(URL.createObjectURL(await r.blob()), '_blank', 'noopener');
  });
  document.querySelectorAll('[data-edit]').forEach(b => b.onclick = () => {
    const c = JSON.parse(b.dataset.edit);
    $('#cPhone').value = c.phone.startsWith('+855') ? '0' + c.phone.slice(4) : c.phone; $('#cName').value = c.name; $('#cVehicle').value = c.vehicle;
    $('#cPlate').value = c.plate_no || ''; $('#cActive').checked = c.active; $('#cPhone').scrollIntoView({behavior: 'smooth', block: 'center'});
  });
  $('#cSave').onclick = async () => {
    try {
      await api('/admin/api/couriers', {method: 'POST', body: {phone: $('#cPhone').value, name: $('#cName').value, vehicle: $('#cVehicle').value, plate_no: $('#cPlate').value, active: $('#cActive').checked}});
      toast(t('saved')); render();
    } catch (e) { showError(e, {phone: '#cPhone', name: '#cName', vehicle: '#cVehicle', plate_no: '#cPlate'}); }
  };
}
