'use strict';
// Phsar Courier: web HP untuk kurir staf — ambil tugas, scan label di toko, antar ke pembeli.
P.authBase = '/courier/auth';
addStrings({
  en: {appName: 'Courier', logout: 'Sign out', notCourier: 'This phone number is not registered as a Phsar courier. Ask the Phsar admin to add you.',
    online: 'Online', offline: 'Offline', goOnline: 'Go online to get new pickup alerts on Telegram.',
    v_available: 'Available', v_mine: 'My tasks', v_done: 'Done today', noTasks: 'No tasks right now.',
    items: '{n} items', readyAgo: 'Ready {t} ago', take: 'Take task', taken: 'Task is yours. Go to the store.', release: 'Give back',
    st_unassigned: 'Available', st_assigned: 'Go to store', st_picked_up: 'Deliver', st_delivered: 'Delivered', st_failed: 'Delivery failed',
    step1: '1 · Pick up at the store', step2: '2 · Deliver to the buyer', contact: 'Contact', call: 'Call', navigate: 'Navigate', contents: 'Package contents',
    scan: 'Scan package label', typeCode: 'Or type the order number', confirmPickup: 'Confirm pickup', pickedUp: 'Picked up. Deliver to the buyer.',
    scanHint: 'Point the camera at the barcode on the label.', noCamera: 'This phone cannot scan barcodes here. Type the order number from the label.', stopScan: 'Close camera',
    recipient: 'Recipient', note: 'Note', proof: 'Photo of the handover (optional)', delivered: 'Delivered', deliveredDone: 'Delivered ✅ Thank you!',
    failed: 'Could not deliver', failTitle: 'What happened?', f_buyer_unreachable: 'Buyer not reachable', f_buyer_refused: 'Buyer refused',
    f_wrong_address: 'Wrong address', f_damaged: 'Package damaged', f_other: 'Other', failNote: 'Details (required for “Other”)', saveFail: 'Save', cancel: 'Cancel',
    failSaved: 'Saved. Keep the package — Phsar admin will contact you.', retry: 'Try delivering again', back: '← Tasks', uploading: 'Uploading photo…',
    afterPickup: 'The buyer’s address and phone appear after you pick up the package.', order: 'Order', today: 'today', area: 'Area', confirmRelease: 'Give this task back to the list?',
    err_taken: 'Another courier already took this task.', err_wrong_package: 'Wrong package! This label is for a different order.', err_too_many: 'Deliver your current packages first.',
    err_bad_status: 'This task was already updated. Refreshing…', err_not_courier: 'This account is not an active Phsar courier.'},
  km: {appName: 'អ្នកដឹក', logout: 'ចាកចេញ', notCourier: 'លេខទូរស័ព្ទនេះមិនទាន់ចុះឈ្មោះជាអ្នកដឹក Phsar ទេ។ សូមស្នើអ្នកគ្រប់គ្រង Phsar បន្ថែមអ្នក។',
    online: 'អនឡាញ', offline: 'ក្រៅបណ្តាញ', goOnline: 'បើកអនឡាញ ដើម្បីទទួលដំណឹងកញ្ចប់ថ្មីតាម Telegram។',
    v_available: 'អាចទទួលបាន', v_mine: 'ការងាររបស់ខ្ញុំ', v_done: 'បានបញ្ចប់ថ្ងៃនេះ', noTasks: 'មិនទាន់មានការងារទេ។',
    items: '{n} មុខ', readyAgo: 'រួចរាល់ {t} មុន', take: 'ទទួលការងារ', taken: 'ការងារនេះជារបស់អ្នក។ សូមទៅហាង។', release: 'ប្រគល់វិញ',
    st_unassigned: 'អាចទទួលបាន', st_assigned: 'ទៅហាង', st_picked_up: 'ដឹកជូន', st_delivered: 'បានដឹកដល់', st_failed: 'ដឹកមិនបាន',
    step1: '1 · យកនៅហាង', step2: '2 · ដឹកជូនអ្នកទិញ', contact: 'ទំនាក់ទំនង', call: 'ហៅ', navigate: 'នាំផ្លូវ', contents: 'ទំនិញក្នុងកញ្ចប់',
    scan: 'ស្កេនស្លាកកញ្ចប់', typeCode: 'ឬវាយលេខការបញ្ជាទិញ', confirmPickup: 'បញ្ជាក់ការយក', pickedUp: 'បានយកហើយ។ សូមដឹកជូនអ្នកទិញ។',
    scanHint: 'តម្រង់កាមេរ៉ាទៅបាកូដលើស្លាក។', noCamera: 'ទូរស័ព្ទនេះមិនអាចស្កេនបាកូដនៅទីនេះ។ សូមវាយលេខការបញ្ជាទិញពីស្លាក។', stopScan: 'បិទកាមេរ៉ា',
    recipient: 'អ្នកទទួល', note: 'កំណត់ចំណាំ', proof: 'រូបថតពេលប្រគល់ (មិនចាំបាច់)', delivered: 'បានដឹកដល់', deliveredDone: 'បានដឹកដល់ ✅ អរគុណ!',
    failed: 'ដឹកមិនបាន', failTitle: 'មានបញ្ហាអ្វី?', f_buyer_unreachable: 'ទាក់ទងអ្នកទិញមិនបាន', f_buyer_refused: 'អ្នកទិញបដិសេធ',
    f_wrong_address: 'អាសយដ្ឋានខុស', f_damaged: 'កញ្ចប់ខូច', f_other: 'ផ្សេងៗ', failNote: 'ព័ត៌មានលម្អិត (ចាំបាច់សម្រាប់ “ផ្សេងៗ”)', saveFail: 'រក្សាទុក', cancel: 'បោះបង់',
    failSaved: 'បានរក្សាទុក។ សូមរក្សាកញ្ចប់ — អ្នកគ្រប់គ្រង Phsar នឹងទាក់ទងអ្នក។', retry: 'ព្យាយាមដឹកម្តងទៀត', back: '← ការងារ', uploading: 'កំពុងផ្ទុករូបថត…',
    afterPickup: 'អាសយដ្ឋាន និងលេខទូរស័ព្ទអ្នកទិញនឹងបង្ហាញបន្ទាប់ពីអ្នកយកកញ្ចប់។', order: 'ការបញ្ជាទិញ', today: 'ថ្ងៃនេះ', area: 'តំបន់', confirmRelease: 'ប្រគល់ការងារនេះត្រឡប់ទៅបញ្ជីវិញ?',
    err_taken: 'អ្នកដឹកផ្សេងបានទទួលការងារនេះរួចហើយ។', err_wrong_package: 'កញ្ចប់ខុស! ស្លាកនេះសម្រាប់ការបញ្ជាទិញផ្សេង។', err_too_many: 'សូមដឹកកញ្ចប់ដែលមានស្រាប់ជាមុនសិន។',
    err_bad_status: 'ការងារនេះបានផ្លាស់ប្តូររួចហើយ។ កំពុងផ្ទុកឡើងវិញ…', err_not_courier: 'គណនីនេះមិនមែនជាអ្នកដឹក Phsar សកម្មទេ។'},
  zh: {appName: '快递员', logout: '退出', notCourier: '此手机号尚未登记为 Phsar 快递员，请联系 Phsar 管理员添加。',
    online: '在线', offline: '离线', goOnline: '上线后可在 Telegram 收到新取件提醒。',
    v_available: '可接任务', v_mine: '我的任务', v_done: '今日完成', noTasks: '暂无任务。',
    items: '{n} 件', readyAgo: '{t}前备好', take: '接单', taken: '已接单，请前往店铺。', release: '退回任务',
    st_unassigned: '可接', st_assigned: '前往店铺', st_picked_up: '配送中', st_delivered: '已送达', st_failed: '配送失败',
    step1: '1 · 到店取件', step2: '2 · 送给买家', contact: '联系人', call: '拨打', navigate: '导航', contents: '包裹内容',
    scan: '扫描包裹面单', typeCode: '或输入订单号', confirmPickup: '确认取件', pickedUp: '已取件，请送给买家。',
    scanHint: '将摄像头对准面单上的条码。', noCamera: '此手机无法在这里扫码，请输入面单上的订单号。', stopScan: '关闭摄像头',
    recipient: '收件人', note: '备注', proof: '交付照片（选填）', delivered: '已送达', deliveredDone: '已送达 ✅ 辛苦了！',
    failed: '无法送达', failTitle: '发生了什么？', f_buyer_unreachable: '联系不上买家', f_buyer_refused: '买家拒收',
    f_wrong_address: '地址错误', f_damaged: '包裹损坏', f_other: '其他', failNote: '详情（选择“其他”时必填）', saveFail: '保存', cancel: '取消',
    failSaved: '已保存。请保管好包裹，Phsar 管理员会联系你。', retry: '再次配送', back: '← 任务', uploading: '正在上传照片…',
    afterPickup: '取件后将显示买家地址和电话。', order: '订单', today: '今天', area: '区域', confirmRelease: '将此任务退回列表？',
    err_taken: '该任务已被其他快递员接走。', err_wrong_package: '包裹不对！此面单属于其他订单。', err_too_many: '请先送完手上的包裹。',
    err_bad_status: '该任务已更新，正在刷新…', err_not_courier: '此账号不是在职的 Phsar 快递员。'},
});

const C = {view: 'mine', me: null, timer: null, stream: null};

function shell(inner) {
  const sw = C.me ? `<button type="button" class="online-toggle ${C.me.online ? 'on' : ''}" id="onlineBtn" aria-pressed="${C.me.online}">
    <span class="dot"></span>${t(C.me.online ? 'online' : 'offline')}</button>` : '';
  return `<header class="top"><a class="brand" href="#/">${brandSVG}<span><b>Phsar</b><small>${esc(t('appName')).toUpperCase()}</small></span></a>
    <div class="top-actions">${sw}${langSwitch()}</div></header><main class="courier">${inner}
    ${P.token ? `<p class="muted center-text">${C.me ? esc(C.me.name) + ' · ' + esc(C.me.phone) + ' · ' : ''}<button type="button" class="link" onclick="logout()">${t('logout')}</button></p>` : ''}</main>`;
}

function ago(iso) {
  const m = Math.max(0, Math.round((Date.now() - new Date(iso)) / 60000));
  return m < 60 ? m + 'm' : Math.floor(m / 60) + 'h ' + (m % 60) + 'm';
}
const mapsDir = (lat, lng) => `https://www.google.com/maps/dir/?api=1&destination=${lat},${lng}&travelmode=driving`;
const tel = p => `tel:${String(p).replace(/[^0-9+]/g, '')}`;

async function render() {
  stopScan();
  clearInterval(C.timer);
  if (!P.token) {
    $('#app').innerHTML = shell(loginView('Phsar ' + t('appName')));
    bindLogin(() => { location.hash = '#/'; render(); });
    return;
  }
  try {
    C.me = await api('/v1/courier/me');
  } catch (e) {
    C.me = null;
    $('#app').innerHTML = shell(`<div class="banner bad">${esc(e.code === 'not_courier' ? t('notCourier') : errMsg(e))}</div>`);
    return;
  }
  const m = location.hash.match(/^#\/task\/([0-9a-f-]+)/);
  try {
    if (m) await renderTask(m[1]); else await renderList();
  } catch (e) { $('#app').innerHTML = shell(`<div class="banner bad">${esc(errMsg(e))}</div>`); }
  const ob = $('#onlineBtn');
  if (ob) ob.onclick = async () => {
    try { C.me = await api('/v1/courier/me/online', {method: 'POST', body: {online: !C.me.online}}); render(); } catch (e) { toast(errMsg(e), 'error'); }
  };
}

async function renderList() {
  const [{tasks}, avail] = await Promise.all([api('/v1/courier/tasks?view=' + C.view), C.view === 'available' ? null : api('/v1/courier/tasks?view=available')]);
  const nAvail = C.view === 'available' ? tasks.length : avail.tasks.length;
  const counts = {available: nAvail, mine: C.me.active_tasks, done: C.me.done_today};
  const tabs = ['mine', 'available', 'done'].map(v => `<button type="button" aria-selected="${C.view === v}" data-view="${v}">${t('v_' + v)}${counts[v] ? ` <span class="count">${counts[v]}</span>` : ''}</button>`).join('');
  const cards = tasks.map(x => `<a class="card task st-${x.status}" href="#/task/${x.id}">
      <div class="spread"><b>${esc(x.store_name)}</b><span class="chip ${x.status === 'failed' ? 'bad' : x.status === 'picked_up' ? 'pink' : x.status === 'delivered' ? 'ok' : 'warn'}">${t('st_' + x.status)}</span></div>
      <div class="muted">${esc(x.pickup_label)} · ${esc(x.pickup_address)}</div>
      <div class="task-meta"><span>${t('order')} <b>${esc(x.order_no)}</b></span><span>${t('items', {n: x.item_count})}</span>
        ${x.status === 'unassigned' ? `<span>${t('readyAgo', {t: ago(x.ready_at)})}</span>` : ''}</div>
      ${x.status === 'picked_up' || x.status === 'failed' ? `<div class="muted">→ ${esc(x.recipient)} · ${esc(x.address_line)}</div>` : ''}
      ${x.status === 'unassigned' ? `<button type="button" class="primary block" data-take="${x.id}">${t('take')}</button>` : ''}</a>`).join('');
  $('#app').innerHTML = shell(`${!C.me.online ? `<div class="banner info">${t('goOnline')}</div>` : ''}<div class="tabs">${tabs}</div>
    ${cards || `<section class="card"><p class="muted">${t('noTasks')}</p></section>`}`);
  document.querySelectorAll('[data-view]').forEach(b => b.onclick = () => { C.view = b.dataset.view; render(); });
  document.querySelectorAll('[data-take]').forEach(b => b.onclick = async e => {
    e.preventDefault();
    b.disabled = true;
    try { await api(`/v1/courier/tasks/${b.dataset.take}/claim`, {method: 'POST'}); toast(t('taken'), 'warning'); location.hash = '#/task/' + b.dataset.take; }
    catch (err) { toast(errMsg(err), 'error'); render(); }
  });
  C.timer = setInterval(() => { if (!document.hidden && !location.hash.startsWith('#/task/')) render(); }, 20000);
}

async function renderTask(id) {
  const x = await api('/v1/courier/tasks/' + id);
  const items = x.items.map(it => `<li><b>${it.qty} ×</b> ${esc(it.name)}${it.options ? ` <span class="muted">(${esc(it.options)})</span>` : ''}</li>`).join('');
  const pickupDone = ['picked_up', 'delivered', 'failed', 'returning', 'returned'].includes(x.status); // courierreturn: + returning/returned
  const step1 = `<section class="card ${pickupDone ? 'done' : ''}"><h2>${t('step1')} ${pickupDone ? '✓' : ''}</h2>
    <p><b>${esc(x.store_name)}</b><br>${esc(x.pickup_label)} · ${esc(x.pickup_address)}</p>
    <p class="muted">${t('contact')}: ${esc(x.pickup_contact)}</p>
    <div class="row"><a class="btn" href="${tel(x.pickup_phone)}">📞 ${t('call')}</a><a class="btn" href="${mapsDir(x.pickup_lat, x.pickup_lng)}" target="_blank" rel="noopener">🧭 ${t('navigate')}</a></div>
    <h3>${t('contents')} · ${t('order')} ${esc(x.order_no)}</h3><ul class="plain">${items}</ul>
    ${x.status === 'unassigned' ? `<button type="button" class="primary block" id="takeBtn">${t('take')}</button>` : ''}
    ${x.status === 'assigned' ? `<button type="button" class="primary block" id="scanBtn">📷 ${t('scan')}</button><div id="scanBox"></div>
      <label>${t('typeCode')}<input id="codeIn" autocapitalize="characters" autocomplete="off" placeholder="PH000000-XXXXX"></label>
      <div class="row"><button type="button" id="pickBtn">${t('confirmPickup')}</button><button type="button" class="ghost" id="relBtn">${t('release')}</button></div>` : ''}</section>`;
  const step2 = x.status === 'unassigned' || x.status === 'assigned' ? `<section class="card muted-card"><h2>${t('step2')}</h2><p class="muted">${t('afterPickup')}</p></section>`
    : `<section class="card ${x.status === 'delivered' ? 'done' : ''}"><h2>${t('step2')} ${x.status === 'delivered' ? '✓' : ''}</h2>
      <p>${t('recipient')}: <b>${esc(x.recipient)}</b><br>${esc(x.address_line)}</p>${x.address_note ? `<p class="muted">${t('note')}: ${esc(x.address_note)}</p>` : ''}
      <div class="row"><a class="btn" href="${tel(x.phone)}">📞 ${t('call')}</a><a class="btn" href="${mapsDir(x.dropoff_lat, x.dropoff_lng)}" target="_blank" rel="noopener">🧭 ${t('navigate')}</a></div>
      ${x.status === 'picked_up' ? `<label>${t('proof')}<input type="file" id="proofIn" accept="image/*" capture="environment"></label>
        <button type="button" class="primary block" id="delBtn">✅ ${t('delivered')}</button><button type="button" class="danger block" id="failBtn">${t('failed')}</button>` : ''}
      ${x.status === 'failed' ? `<div class="banner bad">${t('f_' + x.fail_reason)}${x.fail_note ? ' — ' + esc(x.fail_note) : ''}</div>
        <button type="button" class="primary block" id="retryBtn">${t('retry')}</button>` : ''}</section>`;
  $('#app').innerHTML = shell(`<p><a href="#/">${t('back')}</a></p><div class="spread"><h1>${esc(x.order_no)}</h1><span class="chip">${t('st_' + x.status)}</span></div>${step1}${step2}`);
  const post = async (path, body, okMsg) => {
    try { await api(`/v1/courier/tasks/${id}/${path}`, {method: 'POST', body}); if (okMsg) toast(okMsg); render(); }
    catch (e) { toast(errMsg(e), 'error'); if (e.code === 'bad_status' || e.code === 'taken') setTimeout(render, 1200); }
  };
  $('#takeBtn') && ($('#takeBtn').onclick = () => post('claim', undefined, t('taken')));
  $('#relBtn') && ($('#relBtn').onclick = () => { if (confirm(t('confirmRelease'))) post('release').then(() => { location.hash = '#/'; }); });
  $('#pickBtn') && ($('#pickBtn').onclick = () => post('pickup', {code: $('#codeIn').value}, t('pickedUp')));
  $('#scanBtn') && ($('#scanBtn').onclick = () => startScan(code => { $('#codeIn').value = code; post('pickup', {code}, t('pickedUp')); }));
  $('#retryBtn') && ($('#retryBtn').onclick = () => post('retry'));
  $('#failBtn') && ($('#failBtn').onclick = () => failDialog(id));
  if (typeof crTaskExtra === 'function') await crTaskExtra(x, id); // courierreturn: kembalikan ke toko
  $('#delBtn') && ($('#delBtn').onclick = async () => {
    const btn = $('#delBtn');
    btn.disabled = true;
    let proof = '';
    const f = $('#proofIn').files[0];
    try {
      if (f) {
        btn.textContent = t('uploading');
        const fd = new FormData(); fd.append('file', await shrink(f));
        proof = (await api('/v1/courier/uploads', {method: 'POST', body: fd})).key;
      }
      await api(`/v1/courier/tasks/${id}/deliver`, {method: 'POST', body: {proof_key: proof}});
      toast(t('deliveredDone'));
      location.hash = '#/';
    } catch (e) { toast(errMsg(e), 'error'); btn.disabled = false; btn.textContent = '✅ ' + t('delivered'); }
  });
}

// Perkecil foto dari kamera HP (bisa 5–10 MB) sebelum diunggah.
function shrink(file) {
  return new Promise(resolve => {
    const img = new Image();
    img.onload = () => {
      const k = Math.min(1, 1600 / Math.max(img.width, img.height));
      const c = document.createElement('canvas');
      c.width = Math.round(img.width * k); c.height = Math.round(img.height * k);
      c.getContext('2d').drawImage(img, 0, 0, c.width, c.height);
      c.toBlob(b => resolve(b ? new File([b], 'proof.jpg', {type: 'image/jpeg'}) : file), 'image/jpeg', 0.82);
      URL.revokeObjectURL(img.src);
    };
    img.onerror = () => resolve(file);
    img.src = URL.createObjectURL(file);
  });
}

function failDialog(id) {
  const d = document.createElement('dialog');
  d.className = 'dlg';
  d.innerHTML = `<form method="dialog"><h2>${t('failTitle')}</h2>
    <label>${t('failTitle')}<select id="fReason">${['buyer_unreachable', 'buyer_refused', 'wrong_address', 'damaged', 'other'].map(r => `<option value="${r}">${t('f_' + r)}</option>`).join('')}</select></label>
    <label>${t('failNote')}<textarea id="fNote" maxlength="300"></textarea></label>
    <label>${t('crFailProof')}<input type="file" id="fProof" accept="image/*" capture="environment"></label>
    <div class="row"><button type="button" class="danger" id="fGo">${t('saveFail')}</button><button value="cancel">${t('cancel')}</button></div></form>`;
  document.body.append(d);
  d.showModal();
  d.addEventListener('close', () => d.remove());
  $('#fGo', d).onclick = async () => {
    try {
      // courierreturn: foto bukti gagal antar (opsional)
      let proof = '';
      const f = $('#fProof', d).files[0];
      if (f) { const fd = new FormData(); fd.append('file', await shrink(f)); proof = (await api('/v1/courier/uploads', {method: 'POST', body: fd})).key; }
      await api(`/v1/courier/tasks/${id}/fail`, {method: 'POST', body: {reason: $('#fReason', d).value, note: $('#fNote', d).value, proof_key: proof}});
      d.remove(); toast(t('failSaved')); render();
    } catch (e) { showError(e, {reason: '#fReason', note: '#fNote'}); }
  };
}

// ---- Scan barcode label (BarcodeDetector: Chrome Android). iPhone: ketik nomor pesanan. ----
async function startScan(onCode) {
  const box = $('#scanBox');
  let formats = [];
  try { formats = 'BarcodeDetector' in window ? await BarcodeDetector.getSupportedFormats() : []; } catch {}
  if (!formats.includes('code_128') || !navigator.mediaDevices?.getUserMedia) { toast(t('noCamera'), 'warning'); $('#codeIn').focus(); return; }
  let stream;
  try { stream = await navigator.mediaDevices.getUserMedia({video: {facingMode: 'environment'}, audio: false}); }
  catch { toast(t('noCamera'), 'warning'); $('#codeIn').focus(); return; }
  C.stream = stream;
  box.innerHTML = `<div class="scanner"><video playsinline muted></video><div class="scan-line"></div></div><p class="muted">${t('scanHint')}</p>
    <button type="button" class="ghost" id="stopScan">${t('stopScan')}</button>`;
  const video = $('video', box);
  video.srcObject = stream;
  await video.play().catch(() => {});
  $('#stopScan').onclick = stopScan;
  const det = new BarcodeDetector({formats: ['code_128']});
  const loop = async () => {
    if (C.stream !== stream) return;
    try {
      const codes = await det.detect(video);
      const hit = codes.find(c => /^PH\d{6}-[A-Z0-9]{5}$/i.test(c.rawValue.trim()));
      if (hit) { if (navigator.vibrate) navigator.vibrate(80); stopScan(); onCode(hit.rawValue.trim()); return; }
    } catch {}
    requestAnimationFrame(loop);
  };
  loop();
}

function stopScan() {
  if (C.stream) { C.stream.getTracks().forEach(tr => tr.stop()); C.stream = null; }
  const box = $('#scanBox');
  if (box) box.innerHTML = '';
}

window.addEventListener('hashchange', render);
(async () => { document.documentElement.lang = P.lang; await refreshSession(); render(); })();
