'use strict';
// Admin (super admin): aturan batas waktu pesanan — terima, kemas, hitung hanya saat toko buka, jam standar.
addStrings({
  en: {otDesc: 'How long sellers have to accept and pack orders before they are cancelled automatically (with a refund). Changes apply to new orders.',
    otAccept: 'Accept a paid order within', otAcceptHint: 'If the seller doesn’t tap Accept in time, the order is cancelled and refunded.',
    otPack: 'Pack time limit', otPackOn: 'Cancel accepted orders that are not packed in time', otPackWithin: 'Mark “ready for pickup” within', otPackHint: 'Counted from the moment the seller accepts.',
    otBiz: 'Count only while the shop is open', otBizHint: 'The clock pauses when the shop is closed. An order paid at night starts counting when the shop opens.',
    otDefault: 'Standard opening hours', otDefaultHint: 'Used for shops that have not set their own hours.', otSave: 'Save rules', otSaved: 'Saved — applies to new orders',
    otMinutes: 'minutes', otHours: 'hours', otPreview: 'Example: paid {from} → accept by {to}', otNeedOne: 'Open at least one day a week.'},
  km: {otDesc: 'រយៈពេលដែលអ្នកលក់មានដើម្បីទទួល និងវេចខ្ចប់ការបញ្ជាទិញ មុនពេលបោះបង់ដោយស្វ័យប្រវត្តិ (សងប្រាក់)។ ការផ្លាស់ប្តូរអនុវត្តចំពោះការបញ្ជាទិញថ្មី។',
    otAccept: 'ទទួលការបញ្ជាទិញដែលបានបង់ក្នុង', otAcceptHint: 'បើអ្នកលក់មិនចុចទទួលទាន់ពេល ការបញ្ជាទិញនឹងត្រូវបោះបង់ និងសងប្រាក់។',
    otPack: 'ពេលកំណត់វេចខ្ចប់', otPackOn: 'បោះបង់ការបញ្ជាទិញដែលបានទទួល តែមិនបានវេចខ្ចប់ទាន់ពេល', otPackWithin: 'សម្គាល់ “រួចរាល់ឱ្យមកយក” ក្នុង', otPackHint: 'រាប់ចាប់ពីពេលអ្នកលក់ទទួល។',
    otBiz: 'រាប់តែពេលហាងបើក', otBizHint: 'នាឡិកាផ្អាកពេលហាងបិទ។ ការបញ្ជាទិញដែលបង់ពេលយប់ ចាប់ផ្តើមរាប់ពេលហាងបើក។',
    otDefault: 'ម៉ោងបើកស្តង់ដារ', otDefaultHint: 'ប្រើសម្រាប់ហាងដែលមិនទាន់កំណត់ម៉ោងផ្ទាល់ខ្លួន។', otSave: 'រក្សាទុកច្បាប់', otSaved: 'បានរក្សាទុក — អនុវត្តចំពោះការបញ្ជាទិញថ្មី',
    otMinutes: 'នាទី', otHours: 'ម៉ោង', otPreview: 'ឧទាហរណ៍៖ បង់ {from} → ទទួលមុន {to}', otNeedOne: 'ត្រូវបើកយ៉ាងហោចណាស់មួយថ្ងៃក្នុងមួយសប្តាហ៍។'},
  zh: {otDesc: '卖家接单和打包的时限，超时订单会自动取消并退款。修改只对新订单生效。',
    otAccept: '已付款订单接单时限', otAcceptHint: '卖家未及时点“接单”，订单将被取消并退款。',
    otPack: '打包时限', otPackOn: '已接单但未及时打包的订单自动取消', otPackWithin: '标记“待取件”的时限', otPackHint: '从卖家接单时开始计算。',
    otBiz: '只在店铺营业时计时', otBizHint: '店铺休息时暂停计时。夜间付款的订单在店铺开门后才开始计时。',
    otDefault: '标准营业时间', otDefaultHint: '用于尚未设置营业时间的店铺。', otSave: '保存规则', otSaved: '已保存——对新订单生效',
    otMinutes: '分钟', otHours: '小时', otPreview: '例如：{from} 付款 → {to} 前接单', otNeedOne: '每周至少营业一天。'},
});

async function renderOrderTiming() {
  const d = await api('/admin/api/order-timing');
  $('#app').innerHTML = shell(`${pageHead(t('i_ordertiming'), t('otDesc'))}
    <section class="card" style="max-width:720px">
      <label>${t('otAccept')}<span class="row"><input id="otAccept" type="number" min="15" max="2880" step="5" value="${d.accept_minutes}" style="max-width:140px"> ${t('otMinutes')}</span>
        <small class="muted">${t('otAcceptHint')}</small></label>
    </section>
    <section class="card" style="max-width:720px"><h2>${t('otPack')}</h2>
      <label class="inline"><input type="checkbox" id="otPackOn" ${d.pack_enabled ? 'checked' : ''}> ${t('otPackOn')}</label>
      <label>${t('otPackWithin')}<span class="row"><input id="otPack" type="number" min="0.5" max="168" step="0.5" value="${d.pack_minutes / 60}" style="max-width:140px"> ${t('otHours')}</span>
        <small class="muted">${t('otPackHint')}</small></label>
    </section>
    <section class="card" style="max-width:720px">
      <label class="inline"><input type="checkbox" id="otBiz" ${d.business_hours ? 'checked' : ''}> <b>${t('otBiz')}</b></label>
      <p class="muted">${t('otBizHint')}</p>
      <h2>${t('otDefault')}</h2><p class="muted">${t('otDefaultHint')}</p>
      ${hoursRows(d.default_hours, 'otHours')}
    </section>
    <div><button type="button" class="primary" id="otSave" data-errors-anchor>${t('otSave')}</button></div>`);
  bindHoursRows('otHours');
  if (typeof exAutoStrikesBlock === 'function') exAutoStrikesBlock(); // extras: blok strike otomatis
  $('#otSave').onclick = async () => {
    const week = readHoursRows('otHours');
    if (!week) { toast(t('otNeedOne'), 'warning'); return; }
    try {
      await api('/admin/api/order-timing', {method: 'PUT', body: {accept_minutes: Math.round(+$('#otAccept').value), pack_enabled: $('#otPackOn').checked,
        pack_minutes: Math.round(+$('#otPack').value * 60), business_hours: $('#otBiz').checked, default_hours: week}});
      toast(t('otSaved'));
      renderOrderTiming();
    } catch (e) { showError(e, {accept_minutes: '#otAccept', pack_minutes: '#otPack'}); }
  };
}
