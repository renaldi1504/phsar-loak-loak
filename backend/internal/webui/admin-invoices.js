'use strict';
// Admin: pengaturan pajak (PPN) + faktur komisi bulanan per toko. Faktur disalin saat terbit; void tetap menyimpan nomornya.
addStrings({
  en: {ivHint: 'One commission invoice per store per month, built from completed orders minus commission returned on refunds. Generate after the month ends.',
    ivTaxTitle: 'Tax settings', ivLegal: 'Company legal name', ivVatTin: 'VAT TIN', ivAddress: 'Address', ivEmail: 'Email', ivPhone: 'Phone', ivRate: 'VAT rate (%)',
    ivInclusive: 'Commission already includes VAT', ivPrefix: 'Invoice number prefix', ivSaveTax: 'Save tax settings', ivSaved: 'Saved',
    ivGenTitle: 'Generate invoices', ivMonth: 'Month', ivGenerate: 'Generate for this month', ivGenDone: '{n} invoices created ({s} already existed)',
    ivList: 'Invoices', ivAll: 'All', iv_issued: 'Issued', iv_void: 'Void', ivSearch: 'Number or store', ivNo: 'Number', ivStore: 'Store', ivPeriod: 'Month',
    ivOrders: 'Orders', ivComm: 'Commission', ivVat: 'VAT', ivTotal: 'Total', ivView: 'View / print', ivVoid: 'Void', ivVoidAsk: 'Why void this invoice? (shown on the invoice)',
    ivVoided: 'Invoice voided. Generate the month again to issue a replacement.', ivEmpty: 'No invoices for this filter.', ivReplaces: 'replaces {n}'},
  km: {ivHint: 'វិក្កយបត្រកម្រៃមួយក្នុងមួយហាងក្នុងមួយខែ ពីការបញ្ជាទិញបានបញ្ចប់ ដកកម្រៃដែលបានសងវិញ។ បង្កើតបន្ទាប់ពីខែបញ្ចប់។',
    ivTaxTitle: 'ការកំណត់ពន្ធ', ivLegal: 'ឈ្មោះក្រុមហ៊ុនផ្លូវការ', ivVatTin: 'លេខ VAT TIN', ivAddress: 'អាសយដ្ឋាន', ivEmail: 'អ៊ីមែល', ivPhone: 'ទូរស័ព្ទ', ivRate: 'អត្រា VAT (%)',
    ivInclusive: 'កម្រៃរួមបញ្ចូល VAT រួចហើយ', ivPrefix: 'បុព្វបទលេខវិក្កយបត្រ', ivSaveTax: 'រក្សាទុកការកំណត់ពន្ធ', ivSaved: 'បានរក្សាទុក',
    ivGenTitle: 'បង្កើតវិក្កយបត្រ', ivMonth: 'ខែ', ivGenerate: 'បង្កើតសម្រាប់ខែនេះ', ivGenDone: 'បានបង្កើត {n} វិក្កយបត្រ (មានរួចហើយ {s})',
    ivList: 'វិក្កយបត្រ', ivAll: 'ទាំងអស់', iv_issued: 'បានចេញ', iv_void: 'បានលុបចោល', ivSearch: 'លេខ ឬហាង', ivNo: 'លេខ', ivStore: 'ហាង', ivPeriod: 'ខែ',
    ivOrders: 'ការបញ្ជាទិញ', ivComm: 'កម្រៃ', ivVat: 'VAT', ivTotal: 'សរុប', ivView: 'មើល / បោះពុម្ព', ivVoid: 'លុបចោល', ivVoidAsk: 'ហេតុអ្វីលុបចោលវិក្កយបត្រនេះ? (បង្ហាញលើវិក្កយបត្រ)',
    ivVoided: 'បានលុបចោល។ បង្កើតខែនោះម្តងទៀតដើម្បីចេញវិក្កយបត្រជំនួស។', ivEmpty: 'គ្មានវិក្កយបត្រទេ។', ivReplaces: 'ជំនួស {n}'},
  zh: {ivHint: '每家店每月一张佣金发票，基于已完成订单减去退款时退还的佣金。请在月份结束后生成。',
    ivTaxTitle: '税务设置', ivLegal: '公司法定名称', ivVatTin: '增值税号 (VAT TIN)', ivAddress: '地址', ivEmail: '邮箱', ivPhone: '电话', ivRate: '增值税率 (%)',
    ivInclusive: '佣金已含增值税', ivPrefix: '发票编号前缀', ivSaveTax: '保存税务设置', ivSaved: '已保存',
    ivGenTitle: '生成发票', ivMonth: '月份', ivGenerate: '生成该月发票', ivGenDone: '已生成 {n} 张发票（{s} 张已存在）',
    ivList: '发票', ivAll: '全部', iv_issued: '已开具', iv_void: '已作废', ivSearch: '编号或店铺', ivNo: '编号', ivStore: '店铺', ivPeriod: '月份',
    ivOrders: '订单', ivComm: '佣金', ivVat: '增值税', ivTotal: '合计', ivView: '查看 / 打印', ivVoid: '作废', ivVoidAsk: '作废原因？（显示在发票上）',
    ivVoided: '发票已作废。再次生成该月即可开具替代发票。', ivEmpty: '没有符合条件的发票。', ivReplaces: '替代 {n}'},
});

const AIV = {period: '', status: '', q: ''};

async function renderInvoices() {
  const qs = new URLSearchParams({period: AIV.period, status: AIV.status, q: AIV.q});
  const d = await api('/admin/api/invoices?' + qs);
  const s = d.settings;
  const f = (k, l, extra = '') => `<label>${t(l)}<input name="${k}" value="${esc(s[k] || '')}" ${extra}></label>`;
  const rows = d.invoices.map(v => `<tr><td><b>${esc(v.number)}</b>${v.replaces_number ? `<br><span class="muted small">${esc(t('ivReplaces', {n: v.replaces_number}))}</span>` : ''}</td>
    <td>${esc(v.store_name)}</td><td>${esc(v.period)}</td><td>${v.order_count}</td><td>$${esc(v.commission)}</td><td>$${esc(v.vat_amount)}</td><td><b>$${esc(v.total)}</b></td>
    <td><span class="chip ${v.status === 'issued' ? 'ok' : 'bad'}">${t('iv_' + v.status)}</span>${v.void_note ? `<br><span class="muted small">${esc(v.void_note)}</span>` : ''}</td>
    <td style="white-space:nowrap"><button type="button" data-ivview="${v.id}">${t('ivView')}</button> ${v.status === 'issued' ? `<button type="button" class="danger" data-ivvoid="${v.id}">${t('ivVoid')}</button>` : ''}</td></tr>`).join('');
  $('#app').innerHTML = shell(`${pageHead(t('i_invoices'), t('ivHint'))}
    <section class="panel" id="ivTax"><div class="ph"><h2>${t('ivTaxTitle')}</h2></div>
      <div style="padding:16px 18px;display:flex;flex-direction:column;gap:14px">
        <div style="display:grid;grid-template-columns:repeat(auto-fit,minmax(220px,1fr));gap:14px">
          ${f('legal_name', 'ivLegal', 'maxlength="160"')}${f('vat_tin', 'ivVatTin', 'maxlength="40"')}${f('address', 'ivAddress', 'maxlength="300"')}
          ${f('email', 'ivEmail', 'type="email" maxlength="120"')}${f('phone', 'ivPhone', 'maxlength="40"')}${f('vat_rate', 'ivRate', 'inputmode="decimal" maxlength="6"')}
          ${f('prefix', 'ivPrefix', 'maxlength="8" style="text-transform:uppercase"')}
          <label class="inline" style="align-self:end"><input type="checkbox" name="vat_inclusive" ${s.vat_inclusive ? 'checked' : ''}> ${t('ivInclusive')}</label></div>
        <div><button class="primary" type="button" id="ivSaveTax">${t('ivSaveTax')}</button></div></div></section>
    <section class="panel"><div class="ph"><h2>${t('ivGenTitle')}</h2></div>
      <div style="padding:16px 18px;display:flex;flex-wrap:wrap;gap:12px;align-items:flex-end">
        <label>${t('ivMonth')}<input type="month" id="ivMonth" name="period" value="${esc(d.last_month)}" max="${esc(d.last_month)}"></label>
        <button class="primary" type="button" id="ivGen">${t('ivGenerate')}</button></div></section>
    <section class="panel"><div class="ph"><h2>${t('ivList')}</h2></div>
      <div style="padding:12px 18px;display:flex;flex-wrap:wrap;gap:10px;align-items:flex-end">
        <label>${t('ivPeriod')}<input type="month" id="ivFPeriod" value="${esc(AIV.period)}"></label>
        <label>Status<select id="ivFStatus"><option value="">${t('ivAll')}</option>${['issued', 'void'].map(x => `<option value="${x}" ${AIV.status === x ? 'selected' : ''}>${t('iv_' + x)}</option>`).join('')}</select></label>
        <label style="flex:1;min-width:180px">${t('ivSearch')}<input id="ivFQ" value="${esc(AIV.q)}"></label></div>
      ${rows ? `<div class="tbl-wrap"><table><thead><tr><th>${t('ivNo')}</th><th>${t('ivStore')}</th><th>${t('ivPeriod')}</th><th>${t('ivOrders')}</th><th>${t('ivComm')}</th><th>${t('ivVat')}</th><th>${t('ivTotal')}</th><th>Status</th><th></th></tr></thead><tbody>${rows}</tbody></table></div>`
        : `<p class="muted" style="padding:0 18px 16px">${t('ivEmpty')}</p>`}</section>`);
  $('#ivSaveTax').onclick = async () => {
    document.querySelectorAll('#ivTax .invalid').forEach(el => el.classList.remove('invalid'));
    const body = {};
    document.querySelectorAll('#ivTax input[name]').forEach(i => { body[i.name] = i.type === 'checkbox' ? i.checked : i.value.trim(); });
    try { await api('/admin/api/invoices/settings', {method: 'PUT', body}); toast(t('ivSaved')); renderInvoices(); }
    catch (e) { const el = e.field && $(`#ivTax [name="${e.field}"]`); if (el) { el.classList.add('invalid'); el.focus(); } toast(errMsg(e), 'error'); }
  };
  $('#ivGen').onclick = async () => {
    try {
      const r = await api('/admin/api/invoices/generate', {method: 'POST', body: {period: $('#ivMonth').value}});
      toast(t('ivGenDone', {n: r.created.length, s: r.skipped})); AIV.period = r.period; renderInvoices();
    } catch (e) { toast(errMsg(e), 'error'); }
  };
  const refilter = () => { AIV.period = $('#ivFPeriod').value; AIV.status = $('#ivFStatus').value; AIV.q = $('#ivFQ').value.trim(); renderInvoices(); };
  $('#ivFPeriod').onchange = refilter; $('#ivFStatus').onchange = refilter; $('#ivFQ').onchange = refilter;
  document.querySelectorAll('[data-ivview]').forEach(b => b.onclick = async () => {
    const r = await fetch(`/admin/api/invoices/${b.dataset.ivview}/html`, {headers: {Authorization: 'Bearer ' + P.token}});
    if (!r.ok) { toast(r.statusText, 'error'); return; }
    window.open(URL.createObjectURL(await r.blob()), '_blank', 'noopener');
  });
  document.querySelectorAll('[data-ivvoid]').forEach(b => b.onclick = async () => {
    const note = prompt(t('ivVoidAsk'));
    if (note === null) return;
    try { await api(`/admin/api/invoices/${b.dataset.ivvoid}/void`, {method: 'POST', body: {note}}); toast(t('ivVoided')); renderInvoices(); }
    catch (e) { toast(errMsg(e), 'error'); }
  });
}
