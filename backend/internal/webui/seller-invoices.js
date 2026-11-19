'use strict';
// Seller Center (pemilik): faktur komisi bulanan dari Phsar, ditampilkan di bawah tab Payouts.
addStrings({
  en: {sivTitle: 'Commission invoices', sivHint: 'Monthly tax invoices for the commission Phsar deducted. Open one to print or save as PDF.', sivNone: 'No invoices yet — they appear after each month ends.',
    sivOpen: 'Open', sivVoid: 'Void', sivMonth: 'Month', sivComm: 'Commission', sivTotal: 'Total incl. VAT'},
  km: {sivTitle: 'វិក្កយបត្រកម្រៃ', sivHint: 'វិក្កយបត្រអាករប្រចាំខែសម្រាប់កម្រៃដែល Phsar បានកាត់។ បើកដើម្បីបោះពុម្ព ឬរក្សាទុកជា PDF។', sivNone: 'មិនទាន់មានវិក្កយបត្រ — នឹងបង្ហាញបន្ទាប់ពីខែនីមួយៗបញ្ចប់។',
    sivOpen: 'បើក', sivVoid: 'បានលុបចោល', sivMonth: 'ខែ', sivComm: 'កម្រៃ', sivTotal: 'សរុបរួម VAT'},
  zh: {sivTitle: '佣金发票', sivHint: 'Phsar 扣除佣金的月度税务发票。打开即可打印或另存为 PDF。', sivNone: '暂无发票——每月结束后生成。',
    sivOpen: '打开', sivVoid: '已作废', sivMonth: '月份', sivComm: '佣金', sivTotal: '含税合计'},
});

async function sellerInvoices(id) {
  const tab = $('#tab');
  if (!tab) return;
  let list = [];
  try { list = (await api(`/v1/seller/stores/${id}/payouts/invoices`)).invoices; } catch { return; }
  const rows = list.map(v => `<tr><td><b>${esc(v.number)}</b></td><td>${esc(v.period)}</td><td>$${esc(v.commission)}</td><td>$${esc(v.total)}</td>
    <td>${v.status === 'void' ? `<span class="chip bad">${t('sivVoid')}</span>` : `<button type="button" data-siv="${v.id}">${t('sivOpen')}</button>`}</td></tr>`).join('');
  const sec = document.createElement('section');
  sec.className = 'card';
  sec.innerHTML = `<h2>${t('sivTitle')}</h2><p class="muted">${t('sivHint')}</p>${rows ? `<table><tr><th>#</th><th>${t('sivMonth')}</th><th>${t('sivComm')}</th><th>${t('sivTotal')}</th><th></th></tr>${rows}</table>` : `<p class="muted">${t('sivNone')}</p>`}`;
  tab.appendChild(sec);
  sec.querySelectorAll('[data-siv]').forEach(b => b.onclick = async () => {
    const r = await fetch(`/v1/seller/stores/${id}/payouts/invoices/${b.dataset.siv}/html`, {headers: {Authorization: 'Bearer ' + P.token}});
    if (!r.ok) { toast(r.statusText, 'error'); return; }
    window.open(URL.createObjectURL(await r.blob()), '_blank', 'noopener');
  });
}
