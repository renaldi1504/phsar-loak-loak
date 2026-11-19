'use strict';
// courierreturn: kurir — pengantaran gagal → coba lagi (dibatasi) atau bawa paket kembali ke titik jemput toko.
addStrings({
  en: {st_returning: 'Return to store', st_returned: 'Returned to store', crFailProof: 'Photo (optional)', err_max_attempts: 'No delivery attempts left. Return the package to the store.',
    failSaved: 'Saved. Try again later, or bring the package back to the store.',
    crWhatNow: 'What next?', crAttemptsLeft: 'You can try delivering again ({n} attempt(s) left), or bring the package back to the store.',
    crNoAttempts: 'No delivery attempts left. Bring the package back to the store.', crStart: 'Return to store', crConfirmStart: 'Bring this package back to the store? The order will be cancelled and the buyer refunded.',
    crStarted: 'OK. Bring the package back to the store’s pickup point.', crBackTitle: 'Bring back to the store', crBackHint: 'Hand the package to the store at its pickup point.',
    crProof: 'Photo of the handover at the store (optional)', crDone: 'Handed over to the store', crDoneMsg: 'Thank you! The store will confirm it received the package.',
    crReturnedTitle: 'Returned to the store', crReturnedHint: 'Waiting for the store to confirm. Nothing else to do.'},
  km: {st_returning: 'ប្រគល់ទៅហាងវិញ', st_returned: 'បានប្រគល់ទៅហាងវិញ', crFailProof: 'រូបថត (មិនចាំបាច់)', err_max_attempts: 'អស់ចំនួនដងដឹកហើយ។ សូមប្រគល់កញ្ចប់ទៅហាងវិញ។',
    failSaved: 'បានរក្សាទុក។ ព្យាយាមម្តងទៀតពេលក្រោយ ឬនាំកញ្ចប់ត្រឡប់ទៅហាងវិញ។',
    crWhatNow: 'ធ្វើអ្វីបន្ទាប់?', crAttemptsLeft: 'អ្នកអាចព្យាយាមដឹកម្តងទៀត (នៅសល់ {n} ដង) ឬនាំកញ្ចប់ត្រឡប់ទៅហាងវិញ។',
    crNoAttempts: 'អស់ចំនួនដងដឹកហើយ។ សូមនាំកញ្ចប់ត្រឡប់ទៅហាងវិញ។', crStart: 'ប្រគល់ទៅហាងវិញ', crConfirmStart: 'នាំកញ្ចប់នេះត្រឡប់ទៅហាងវិញ? ការបញ្ជាទិញនឹងត្រូវបោះបង់ ហើយអ្នកទិញនឹងទទួលប្រាក់វិញ។',
    crStarted: 'យល់ព្រម។ សូមនាំកញ្ចប់ត្រឡប់ទៅទីតាំងទទួលរបស់ហាង។', crBackTitle: 'នាំត្រឡប់ទៅហាង', crBackHint: 'ប្រគល់កញ្ចប់ឱ្យហាងនៅទីតាំងទទួលរបស់ហាង។',
    crProof: 'រូបថតពេលប្រគល់នៅហាង (មិនចាំបាច់)', crDone: 'បានប្រគល់ឱ្យហាងហើយ', crDoneMsg: 'អរគុណ! ហាងនឹងបញ្ជាក់ថាបានទទួលកញ្ចប់។',
    crReturnedTitle: 'បានប្រគល់ទៅហាងវិញ', crReturnedHint: 'កំពុងរង់ចាំហាងបញ្ជាក់។ មិនចាំបាច់ធ្វើអ្វីទៀតទេ។'},
  zh: {st_returning: '退回店铺', st_returned: '已退回店铺', crFailProof: '照片（选填）', err_max_attempts: '配送次数已用完，请将包裹退回店铺。',
    failSaved: '已保存。可稍后再试，或将包裹退回店铺。',
    crWhatNow: '下一步？', crAttemptsLeft: '你可以再次配送（还剩 {n} 次），或将包裹退回店铺。',
    crNoAttempts: '配送次数已用完，请将包裹退回店铺。', crStart: '退回店铺', crConfirmStart: '将此包裹退回店铺？订单将被取消，买家获得退款。',
    crStarted: '好的，请将包裹送回店铺取件点。', crBackTitle: '送回店铺', crBackHint: '在店铺取件点将包裹交给店铺。',
    crProof: '在店铺交接的照片（选填）', crDone: '已交给店铺', crDoneMsg: '谢谢！店铺会确认收到包裹。',
    crReturnedTitle: '已退回店铺', crReturnedHint: '等待店铺确认，无需其他操作。'},
});

// Dipanggil renderTask (courier.js) setelah halaman tugas digambar.
async function crTaskExtra(x, id) {
  if (!['failed', 'returning', 'returned'].includes(x.status)) return;
  let info;
  try { info = await api('/v1/courier/returns/' + id); } catch { return; }
  const head = document.querySelector('main.courier .spread');
  // halaman sudah berganti / digambar ulang oleh render() lain → jangan dobel
  if (!head || $('#crBox') || !head.textContent.includes(x.order_no)) return;
  let html;
  if (x.status === 'failed') {
    const left = Math.max(0, info.max_attempts - info.attempts);
    html = `<section class="card"><h2>${t('crWhatNow')}</h2><p class="muted">${info.can_retry ? t('crAttemptsLeft', {n: left}) : t('crNoAttempts')}</p>
      ${info.can_return ? `<button type="button" class="primary block" id="crStart">↩️ ${t('crStart')}</button>` : ''}</section>`;
    if (!info.can_retry && $('#retryBtn')) $('#retryBtn').remove();
  } else if (x.status === 'returning') {
    html = `<section class="card"><h2>↩️ ${t('crBackTitle')}</h2><p>${t('crBackHint')}</p>
      <p><b>${esc(x.store_name)}</b><br>${esc(x.pickup_label)} · ${esc(x.pickup_address)}</p>
      <p class="muted">${t('contact')}: ${esc(x.pickup_contact)}</p>
      <div class="row"><a class="btn" href="${tel(x.pickup_phone)}">📞 ${t('call')}</a><a class="btn" href="${mapsDir(x.pickup_lat, x.pickup_lng)}" target="_blank" rel="noopener">🧭 ${t('navigate')}</a></div>
      <label>${t('crProof')}<input type="file" id="crProof" accept="image/*" capture="environment"></label>
      <button type="button" class="primary block" id="crDone">✅ ${t('crDone')}</button></section>`;
  } else {
    html = `<section class="card done"><h2>${t('crReturnedTitle')} ✓</h2><p class="muted">${t('crReturnedHint')}</p></section>`;
  }
  head.insertAdjacentHTML('afterend', `<div id="crBox">${html}</div>`);
  if ($('#crStart')) $('#crStart').onclick = async () => {
    if (!confirm(t('crConfirmStart'))) return;
    try { await api(`/v1/courier/returns/${id}/start`, {method: 'POST'}); toast(t('crStarted')); render(); }
    catch (e) { toast(errMsg(e), 'error'); if (e.code === 'bad_status') setTimeout(render, 1200); }
  };
  if ($('#crDone')) $('#crDone').onclick = async () => {
    const btn = $('#crDone');
    btn.disabled = true;
    try {
      let proof = '';
      const f = $('#crProof').files[0];
      if (f) {
        btn.textContent = t('uploading');
        const fd = new FormData(); fd.append('file', await shrink(f));
        proof = (await api('/v1/courier/uploads', {method: 'POST', body: fd})).key;
      }
      await api(`/v1/courier/returns/${id}/returned`, {method: 'POST', body: {proof_key: proof}});
      toast(t('crDoneMsg'));
      location.hash = '#/';
    } catch (e) { toast(errMsg(e), 'error'); btn.disabled = false; btn.textContent = '✅ ' + t('crDone'); }
  };
}
