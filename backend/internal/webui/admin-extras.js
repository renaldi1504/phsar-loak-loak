'use strict';
// Admin (super admin) — extras: strike otomatis untuk pesanan telat (blok di "Order time limits")
// dan tarik saldo manual seller (blok di "Seller payouts"). Label strike otomatis & payout "diminta seller".
addStrings({
  en: {exAsTitle: 'Automatic strikes for late orders', exAsDesc: 'When an order is cancelled because the seller did not accept or pack it in time, the store can get a strike automatically. Off by default. Admins can revoke automatic strikes like any other strike.',
    exAsOn: 'Give automatic strikes', exAsAccept: 'Not accepted in time', exAsPack: 'Not packed in time', exAsNone: 'No strike', exAsPts: '{n} point(s)',
    exAsDays: 'Strike lasts (days)', exAsMax: 'Max automatic strikes per store per day', exAsMaxHint: 'Stops a burst of late orders on one day from suspending a store at once.',
    exAsSave: 'Save automatic strikes', exSaved: 'Saved',
    exMpTitle: 'Manual withdrawals by sellers', exMpDesc: 'Lets store owners tap “Withdraw now” instead of waiting for the schedule. It creates a normal payout to transfer here. Off by default.',
    exMpOn: 'Allow “Withdraw now”', exMpMin: 'Minimum balance to withdraw ($)', exMpMax: 'Max withdrawals per 7 days', exMpSave: 'Save withdrawal rules',
    exReqBySeller: 'requested by seller', exAuto: 'automatic', exOff: 'Off',
    rr_seller_timeout: 'Not accepted in time', rr_pack_timeout: 'Not packed in time'},
  km: {exAsTitle: 'ការព្រមានស្វ័យប្រវត្តិសម្រាប់ការបញ្ជាទិញយឺត', exAsDesc: 'ពេលការបញ្ជាទិញត្រូវបានបោះបង់ ព្រោះអ្នកលក់មិនបានទទួល ឬវេចខ្ចប់ទាន់ពេល ហាងអាចទទួលការព្រមានដោយស្វ័យប្រវត្តិ។ បិទតាមលំនាំដើម។ អ្នកគ្រប់គ្រងអាចដកការព្រមានស្វ័យប្រវត្តិដូចការព្រមានផ្សេងទៀត។',
    exAsOn: 'ផ្តល់ការព្រមានស្វ័យប្រវត្តិ', exAsAccept: 'មិនបានទទួលទាន់ពេល', exAsPack: 'មិនបានវេចខ្ចប់ទាន់ពេល', exAsNone: 'គ្មានការព្រមាន', exAsPts: '{n} ពិន្ទុ',
    exAsDays: 'ការព្រមានមានសុពលភាព (ថ្ងៃ)', exAsMax: 'ការព្រមានស្វ័យប្រវត្តិអតិបរមាក្នុងមួយហាងក្នុងមួយថ្ងៃ', exAsMaxHint: 'ការពារកុំឱ្យការបញ្ជាទិញយឺតច្រើនក្នុងថ្ងៃតែមួយ ធ្វើឱ្យហាងត្រូវផ្អាកភ្លាមៗ។',
    exAsSave: 'រក្សាទុកការព្រមានស្វ័យប្រវត្តិ', exSaved: 'បានរក្សាទុក',
    exMpTitle: 'ការដកប្រាក់ដោយអ្នកលក់', exMpDesc: 'អនុញ្ញាតឱ្យម្ចាស់ហាងចុច “ដកឥឡូវ” ជំនួសការរង់ចាំកាលវិភាគ។ វាបង្កើតការទូទាត់ធម្មតាដែលត្រូវផ្ទេរនៅទីនេះ។ បិទតាមលំនាំដើម។',
    exMpOn: 'អនុញ្ញាត “ដកឥឡូវ”', exMpMin: 'សមតុល្យអប្បបរមាដើម្បីដក ($)', exMpMax: 'ដកបានអតិបរមាក្នុង 7 ថ្ងៃ', exMpSave: 'រក្សាទុកច្បាប់ដកប្រាក់',
    exReqBySeller: 'អ្នកលក់ស្នើសុំ', exAuto: 'ស្វ័យប្រវត្តិ', exOff: 'បិទ',
    rr_seller_timeout: 'មិនបានទទួលទាន់ពេល', rr_pack_timeout: 'មិនបានវេចខ្ចប់ទាន់ពេល'},
  zh: {exAsTitle: '超时订单自动警告', exAsDesc: '订单因卖家未及时接单或打包而被取消时，可自动给店铺一次警告。默认关闭。自动警告可像其他警告一样被管理员撤销。',
    exAsOn: '启用自动警告', exAsAccept: '未及时接单', exAsPack: '未及时打包', exAsNone: '不警告', exAsPts: '{n} 分',
    exAsDays: '警告有效期（天）', exAsMax: '每家店每天最多自动警告次数', exAsMaxHint: '避免同一天多笔超时订单让店铺立刻被暂停。',
    exAsSave: '保存自动警告', exSaved: '已保存',
    exMpTitle: '卖家手动提现', exMpDesc: '店主可点击“立即提现”，无需等待结算日。会生成一笔普通结算，在此转账处理。默认关闭。',
    exMpOn: '允许“立即提现”', exMpMin: '最低提现余额 ($)', exMpMax: '每 7 天最多提现次数', exMpSave: '保存提现规则',
    exReqBySeller: '卖家申请', exAuto: '自动', exOff: '已关闭',
    rr_seller_timeout: '未及时接单', rr_pack_timeout: '未及时打包'},
});

// Blok terpisah di bawah halaman "Order time limits" (dipanggil dari renderOrderTiming).
async function exAutoStrikesBlock() {
  const anchor = $('#otSave') && $('#otSave').parentElement;
  if (!anchor || document.getElementById('exAs')) return;
  let d;
  try { d = await api('/admin/api/auto-strikes'); } catch { return; }
  const pts = (id, v) => `<select id="${id}">${[0, 1, 2, 3].map(n => `<option value="${n}" ${v === n ? 'selected' : ''}>${n ? t('exAsPts', {n}) : t('exAsNone')}</option>`).join('')}</select>`;
  const sec = document.createElement('section');
  sec.className = 'card'; sec.id = 'exAs'; sec.style.maxWidth = '720px'; sec.style.marginTop = '16px';
  sec.innerHTML = `<h2>${t('exAsTitle')}</h2><p class="muted">${t('exAsDesc')}</p>
    <label class="inline"><input type="checkbox" id="exAsOn" ${d.enabled ? 'checked' : ''}> <b>${t('exAsOn')}</b></label>
    <div class="grid">
      <label>${t('exAsAccept')}${pts('exAsP1', d.seller_timeout_points)}</label>
      <label>${t('exAsPack')}${pts('exAsP2', d.pack_timeout_points)}</label>
      <label>${t('exAsDays')}<input id="exAsDays" type="number" min="1" max="365" value="${d.expiry_days}"></label>
      <label>${t('exAsMax')}<input id="exAsMax" type="number" min="1" max="10" value="${d.daily_max}"><small class="muted">${t('exAsMaxHint')}</small></label>
    </div>
    <div><button type="button" class="primary" id="exAsSave" data-errors-anchor>${t('exAsSave')}</button></div>`;
  anchor.after(sec);
  $('#exAsSave').onclick = async () => {
    try {
      await api('/admin/api/auto-strikes', {method: 'PUT', body: {enabled: $('#exAsOn').checked, seller_timeout_points: +$('#exAsP1').value,
        pack_timeout_points: +$('#exAsP2').value, expiry_days: Math.round(+$('#exAsDays').value), daily_max: Math.round(+$('#exAsMax').value)}});
      toast(t('exSaved'));
    } catch (e) { showError(e, {seller_timeout_points: '#exAsP1', pack_timeout_points: '#exAsP2', expiry_days: '#exAsDays', daily_max: '#exAsMax'}); }
  };
}

// Blok terpisah di halaman "Seller payouts" — hanya super admin (izin payout.rules).
async function exManualPayoutBlock() {
  const anchor = $('#sSave') && $('#sSave').closest('details');
  if (!anchor || document.getElementById('exMp') || !hasPerm('payout.rules')) return;
  let d;
  try { d = await api('/admin/api/manual-payout'); } catch { return; }
  const box = document.createElement('details');
  box.className = 'card'; box.id = 'exMp';
  box.innerHTML = `<summary style="font-weight:700;cursor:pointer">${t('exMpTitle')}${d.enabled ? '' : ` <span class="chip">${t('exOff')}</span>`}</summary>
    <p class="muted">${t('exMpDesc')}</p>
    <label class="inline"><input type="checkbox" id="exMpOn" ${d.enabled ? 'checked' : ''}> <b>${t('exMpOn')}</b></label>
    <div class="grid">
      <label>${t('exMpMin')}<input id="exMpMin" inputmode="decimal" value="${esc(d.min_amount)}"></label>
      <label>${t('exMpMax')}<input id="exMpMax" type="number" min="1" max="14" value="${d.max_per_7_days}"></label>
    </div>
    <div class="row"><button type="button" id="exMpSave">${t('exMpSave')}</button></div>`;
  anchor.after(box);
  $('#exMpSave').onclick = async () => {
    try {
      await api('/admin/api/manual-payout', {method: 'PUT', body: {enabled: $('#exMpOn').checked, min_amount: $('#exMpMin').value, max_per_7_days: Math.round(+$('#exMpMax').value)}});
      toast(t('exSaved'));
    } catch (e) { showError(e, {min_amount: '#exMpMin', max_per_7_days: '#exMpMax'}); }
  };
}
