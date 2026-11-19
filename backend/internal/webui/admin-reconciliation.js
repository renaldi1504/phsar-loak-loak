'use strict';
// Admin: rekonsiliasi bank (izin reconciliation.run) — unggah ekspor transaksi ABA PayWay, cocokkan dengan Phsar, selesaikan selisih.
addStrings({
  en: {rcHint: 'Upload the transaction export from the ABA PayWay merchant portal (CSV). Phsar matches every line with the payments and refunds it recorded and lists the differences.',
    rcFuture: 'Automatic daily download through the PayWay API is planned — for now, export the CSV in PayWay and upload it here.',
    rcUpload: 'Upload PayWay export', rcFile: 'CSV file (max 5 MB)', rcFrom: 'Period from (optional)', rcTo: 'Period to (optional)',
    rcPeriodHint: 'Leave empty to use the first and last date in the file. Payments recorded in Phsar inside the period but missing from the file are flagged.',
    rcSend: 'Upload and match', rcUploading: 'Matching…', rcBatches: 'Uploaded files', rcNone: 'No files uploaded yet.',
    rcColFile: 'File', rcColPeriod: 'Period', rcColRows: 'Lines', rcColCollected: 'Collected', rcColRefunded: 'Refunded', rcColOpen: 'Open differences', rcColBy: 'Uploaded',
    rcAllClear: 'All matched', rcBack: '← All files', rcCollected: 'Collected (PayWay)', rcRefunded: 'Refunded (PayWay)', rcFees: 'Fees', rcNet: 'Net',
    rcPhsarPaid: 'Paid in Phsar', rcOpen: 'Open differences', rcRematch: 'Match again', rcDelete: 'Delete file', rcConfirmDelete: 'Delete this uploaded file and its results?',
    rcf_open: 'Open', rcf_exceptions: 'All differences', rcf_resolved: 'Resolved', rcf_all: 'All lines',
    rcLine: 'Line', rcTran: 'Transaction', rcWhen: 'Date', rcPayway: 'PayWay', rcPhsar: 'Phsar', rcResult: 'Result', rcOrders: 'Orders',
    rcResolveNote: 'How was this settled?', rcResolve: 'Resolve', rcResolved: 'Resolved', rcMatched: 'Matched again', rcUploaded: 'File matched',
    rcRefund: 'Refund', rcFromPhsar: 'Phsar record', rcNoRows: 'Nothing to show for this filter.',
    rr_ok: 'Matched', rr_amount_mismatch: 'Amount differs', rr_not_paid_in_phsar: 'Not paid in Phsar', rr_unknown: 'Not in Phsar', rr_duplicate: 'Duplicate line',
    rr_ignored: 'Failed / declined', rr_missing_in_statement: 'Missing in PayWay file', rr_refund_ok: 'Refund matched', rr_refund_mismatch: 'Refund amount differs',
    rr_refund_unrecorded: 'Refund not recorded in Phsar', rr_pending: 'Not matched yet',
    rx_amount_mismatch: 'PayWay amount is different from what the buyer had to pay in Phsar.',
    rx_not_paid_in_phsar: 'Money arrived at PayWay but the checkout is not paid in Phsar (callback lost or arrived after expiry). Check the order and process it or refund the buyer.',
    rx_unknown: 'This transaction ID does not exist in Phsar — e.g. a test payment or a payment made outside Phsar.',
    rx_duplicate: 'The same transaction appears more than once in the file.',
    rx_missing_in_statement: 'Phsar recorded this payment as paid in this period, but it is not in the PayWay file. Check it in the PayWay portal.',
    rx_refund_mismatch: 'Refund amount differs from the refund recorded in Phsar.', rx_refund_unrecorded: 'PayWay shows a refund that nobody recorded in Phsar (Refunds / Returns).'},
  km: {rcHint: 'ផ្ទុកឯកសារប្រតិបត្តិការពី ABA PayWay (CSV)។ Phsar ផ្ទៀងផ្ទាត់បន្ទាត់នីមួយៗជាមួយការបង់ប្រាក់ និងការសងប្រាក់ដែលបានកត់ត្រា ហើយបង្ហាញភាពខុសគ្នា។',
    rcFuture: 'ការទាញយកដោយស្វ័យប្រវត្តិរាល់ថ្ងៃតាម PayWay API នឹងមាននៅពេលក្រោយ — ឥឡូវនេះ សូមនាំចេញ CSV ពី PayWay ហើយផ្ទុកនៅទីនេះ។',
    rcUpload: 'ផ្ទុកឯកសារ PayWay', rcFile: 'ឯកសារ CSV (អតិបរមា 5 MB)', rcFrom: 'ចាប់ពីថ្ងៃ (ស្រេចចិត្ត)', rcTo: 'ដល់ថ្ងៃ (ស្រេចចិត្ត)',
    rcPeriodHint: 'ទុកទទេដើម្បីប្រើថ្ងៃដំបូង និងចុងក្រោយក្នុងឯកសារ។ ការបង់ប្រាក់ក្នុង Phsar ក្នុងរយៈពេលនេះ ដែលមិនមានក្នុងឯកសារ នឹងត្រូវសម្គាល់។',
    rcSend: 'ផ្ទុក និងផ្ទៀងផ្ទាត់', rcUploading: 'កំពុងផ្ទៀងផ្ទាត់…', rcBatches: 'ឯកសារដែលបានផ្ទុក', rcNone: 'មិនទាន់មានឯកសារទេ។',
    rcColFile: 'ឯកសារ', rcColPeriod: 'រយៈពេល', rcColRows: 'បន្ទាត់', rcColCollected: 'ប្រមូលបាន', rcColRefunded: 'សងវិញ', rcColOpen: 'ភាពខុសគ្នានៅបើក', rcColBy: 'បានផ្ទុក',
    rcAllClear: 'ត្រូវគ្នាទាំងអស់', rcBack: '← ឯកសារទាំងអស់', rcCollected: 'ប្រមូលបាន (PayWay)', rcRefunded: 'សងវិញ (PayWay)', rcFees: 'ថ្លៃសេវា', rcNet: 'សុទ្ធ',
    rcPhsarPaid: 'បានបង់ក្នុង Phsar', rcOpen: 'ភាពខុសគ្នានៅបើក', rcRematch: 'ផ្ទៀងផ្ទាត់ម្តងទៀត', rcDelete: 'លុបឯកសារ', rcConfirmDelete: 'លុបឯកសារនេះ និងលទ្ធផលរបស់វា?',
    rcf_open: 'នៅបើក', rcf_exceptions: 'ភាពខុសគ្នាទាំងអស់', rcf_resolved: 'បានដោះស្រាយ', rcf_all: 'បន្ទាត់ទាំងអស់',
    rcLine: 'បន្ទាត់', rcTran: 'ប្រតិបត្តិការ', rcWhen: 'កាលបរិច្ឆេទ', rcPayway: 'PayWay', rcPhsar: 'Phsar', rcResult: 'លទ្ធផល', rcOrders: 'ការបញ្ជាទិញ',
    rcResolveNote: 'តើបានដោះស្រាយយ៉ាងដូចម្តេច?', rcResolve: 'ដោះស្រាយ', rcResolved: 'បានដោះស្រាយ', rcMatched: 'បានផ្ទៀងផ្ទាត់ម្តងទៀត', rcUploaded: 'បានផ្ទៀងផ្ទាត់ឯកសារ',
    rcRefund: 'សងប្រាក់', rcFromPhsar: 'កំណត់ត្រា Phsar', rcNoRows: 'គ្មានអ្វីសម្រាប់តម្រងនេះទេ។',
    rr_ok: 'ត្រូវគ្នា', rr_amount_mismatch: 'ចំនួនខុសគ្នា', rr_not_paid_in_phsar: 'មិនទាន់បង់ក្នុង Phsar', rr_unknown: 'មិនមានក្នុង Phsar', rr_duplicate: 'បន្ទាត់ស្ទួន',
    rr_ignored: 'បរាជ័យ / បដិសេធ', rr_missing_in_statement: 'មិនមានក្នុងឯកសារ PayWay', rr_refund_ok: 'ការសងត្រូវគ្នា', rr_refund_mismatch: 'ចំនួនសងខុសគ្នា',
    rr_refund_unrecorded: 'ការសងមិនបានកត់ក្នុង Phsar', rr_pending: 'មិនទាន់ផ្ទៀងផ្ទាត់',
    rx_amount_mismatch: 'ចំនួនក្នុង PayWay ខុសពីចំនួនដែលអ្នកទិញត្រូវបង់ក្នុង Phsar។',
    rx_not_paid_in_phsar: 'ប្រាក់បានចូល PayWay ប៉ុន្តែការទូទាត់មិនទាន់បង់ក្នុង Phsar (callback បាត់ ឬមកក្រោយផុតពេល)។ សូមពិនិត្យការបញ្ជាទិញ ហើយដំណើរការ ឬសងប្រាក់អ្នកទិញ។',
    rx_unknown: 'លេខប្រតិបត្តិការនេះមិនមានក្នុង Phsar — ឧ. ការបង់សាកល្បង ឬការបង់នៅក្រៅ Phsar។',
    rx_duplicate: 'ប្រតិបត្តិការដដែលលេចឡើងច្រើនដងក្នុងឯកសារ។',
    rx_missing_in_statement: 'Phsar បានកត់ថាបានបង់ក្នុងរយៈពេលនេះ ប៉ុន្តែមិនមានក្នុងឯកសារ PayWay។ សូមពិនិត្យក្នុង PayWay។',
    rx_refund_mismatch: 'ចំនួនសងខុសពីការសងដែលកត់ក្នុង Phsar។', rx_refund_unrecorded: 'PayWay បង្ហាញការសងប្រាក់ដែលគ្មាននរណាកត់ក្នុង Phsar (សងប្រាក់ / ការប្ដូរ)។'},
  zh: {rcHint: '上传 ABA PayWay 商户后台导出的交易 CSV。Phsar 会把每一行与系统记录的付款和退款核对，并列出差异。',
    rcFuture: '通过 PayWay API 每日自动下载的功能将在之后推出——目前请在 PayWay 导出 CSV 后上传到这里。',
    rcUpload: '上传 PayWay 导出文件', rcFile: 'CSV 文件（最大 5 MB）', rcFrom: '期间开始（可选）', rcTo: '期间结束（可选）',
    rcPeriodHint: '留空则使用文件中的第一天和最后一天。Phsar 在该期间记录为已付款但文件中没有的交易会被标出。',
    rcSend: '上传并核对', rcUploading: '正在核对…', rcBatches: '已上传文件', rcNone: '尚未上传文件。',
    rcColFile: '文件', rcColPeriod: '期间', rcColRows: '行数', rcColCollected: '收款', rcColRefunded: '退款', rcColOpen: '未处理差异', rcColBy: '上传',
    rcAllClear: '全部一致', rcBack: '← 所有文件', rcCollected: '收款（PayWay）', rcRefunded: '退款（PayWay）', rcFees: '手续费', rcNet: '净额',
    rcPhsarPaid: 'Phsar 已付款', rcOpen: '未处理差异', rcRematch: '重新核对', rcDelete: '删除文件', rcConfirmDelete: '删除此上传文件及其结果？',
    rcf_open: '未处理', rcf_exceptions: '全部差异', rcf_resolved: '已处理', rcf_all: '所有行',
    rcLine: '行', rcTran: '交易', rcWhen: '日期', rcPayway: 'PayWay', rcPhsar: 'Phsar', rcResult: '结果', rcOrders: '订单',
    rcResolveNote: '如何处理的？', rcResolve: '标记已处理', rcResolved: '已处理', rcMatched: '已重新核对', rcUploaded: '文件已核对',
    rcRefund: '退款', rcFromPhsar: 'Phsar 记录', rcNoRows: '此筛选下没有内容。',
    rr_ok: '一致', rr_amount_mismatch: '金额不一致', rr_not_paid_in_phsar: 'Phsar 未付款', rr_unknown: 'Phsar 无此交易', rr_duplicate: '重复行',
    rr_ignored: '失败 / 拒绝', rr_missing_in_statement: 'PayWay 文件中缺少', rr_refund_ok: '退款一致', rr_refund_mismatch: '退款金额不一致',
    rr_refund_unrecorded: 'Phsar 未记录此退款', rr_pending: '尚未核对',
    rx_amount_mismatch: 'PayWay 金额与买家在 Phsar 应付金额不同。',
    rx_not_paid_in_phsar: '款项已到 PayWay，但 Phsar 中该结账未付款（回调丢失或超时后到达）。请检查订单并处理或退款。',
    rx_unknown: 'Phsar 中没有这个交易号——例如测试付款或 Phsar 之外的付款。',
    rx_duplicate: '同一笔交易在文件中出现多次。',
    rx_missing_in_statement: 'Phsar 记录该笔在此期间已付款，但 PayWay 文件中没有。请在 PayWay 后台核实。',
    rx_refund_mismatch: '退款金额与 Phsar 记录的退款不同。', rx_refund_unrecorded: 'PayWay 显示了一笔无人在 Phsar 记录的退款（退款 / 退货）。'},
});

const RC = {filter: 'open'};
const RC_CHIP = {ok: 'ok', refund_ok: 'ok', ignored: '', pending: '', amount_mismatch: 'bad', not_paid_in_phsar: 'bad', unknown: 'warn', duplicate: 'warn',
  missing_in_statement: 'bad', refund_mismatch: 'bad', refund_unrecorded: 'warn'};

async function renderReconciliation() {
  const m = location.hash.match(/^#\/reconciliation\/([0-9a-f-]{36})/);
  if (m) return reconDetail(m[1]);
  const {batches} = await api('/admin/api/reconciliation');
  const day = d => d ? String(d).slice(0, 10) : '—';
  const rows = batches.map(b => `<tr><td><a href="#/reconciliation/${b.id}"><b>${esc(b.file_name)}</b></a></td><td>${day(b.period_from)} → ${day(b.period_to)}</td>
    <td>${b.row_count}</td><td>${amoney(b.collected)}</td><td>${amoney(b.refunded)}</td>
    <td>${b.open_exceptions ? `<span class="chip bad">${b.open_exceptions}</span>` : `<span class="chip ok">${t('rcAllClear')}</span>`}</td>
    <td class="muted small">${esc(b.uploaded_by || '')}<br>${new Date(b.created_at).toLocaleString()}</td></tr>`).join('');
  $('#app').innerHTML = shell(`${pageHead(t('i_reconciliation'), t('rcHint'))}
    <div class="banner info">${t('rcFuture')}</div>
    <section class="card" id="rcForm"><h2>${t('rcUpload')}</h2>
      <div class="grid"><label>${t('rcFile')}<input type="file" id="rcFile" name="file" accept=".csv,text/csv"></label>
        <label>${t('rcFrom')}<input type="date" id="rcFrom" name="period_from"></label><label>${t('rcTo')}<input type="date" id="rcTo" name="period_to"></label></div>
      <p class="muted small">${t('rcPeriodHint')}</p>
      <div><button type="button" class="primary" id="rcSend">${t('rcSend')}</button></div></section>
    <section class="card"><h2>${t('rcBatches')}</h2>${rows ? `<div class="tbl-wrap"><table><thead><tr><th>${t('rcColFile')}</th><th>${t('rcColPeriod')}</th><th>${t('rcColRows')}</th>
      <th>${t('rcColCollected')}</th><th>${t('rcColRefunded')}</th><th>${t('rcColOpen')}</th><th>${t('rcColBy')}</th></tr></thead><tbody>${rows}</tbody></table></div>`
      : `<p class="muted">${t('rcNone')}</p>`}</section>`);
  $('#rcSend').onclick = async () => {
    const f = $('#rcFile').files[0], b = $('#rcSend');
    document.querySelectorAll('#rcForm .invalid').forEach(el => el.classList.remove('invalid'));
    if (!f) { $('#rcFile').classList.add('invalid'); $('#rcFile').focus(); return; }
    const fd = new FormData(); fd.append('file', f); fd.append('period_from', $('#rcFrom').value); fd.append('period_to', $('#rcTo').value);
    b.disabled = true; b.textContent = t('rcUploading');
    try {
      const d = await api('/admin/api/reconciliation/batches', {method: 'POST', body: fd});
      RC.filter = 'open'; toast(t('rcUploaded')); location.hash = '#/reconciliation/' + d.batch.id;
    } catch (e) { b.disabled = false; b.textContent = t('rcSend'); showError(e, {file: '#rcFile', period_from: '#rcFrom', period_to: '#rcTo'}); }
  };
}

async function reconDetail(id, data) {
  const d = data || await api(`/admin/api/reconciliation/batches/${id}?filter=${RC.filter}`);
  const b = d.batch, c = d.counts || {};
  const when = s => s ? new Date(s).toLocaleString() : '—';
  const exc = ['amount_mismatch', 'not_paid_in_phsar', 'unknown', 'duplicate', 'missing_in_statement', 'refund_mismatch', 'refund_unrecorded'];
  const tabs = ['open', 'exceptions', 'resolved', 'all'].map(f => `<button type="button" aria-selected="${RC.filter === f}" data-rcf="${f}">${t('rcf_' + f)}${f === 'open' && d.open ? ` <span class="count">${d.open}</span>` : ''}</button>`)
    .concat(Object.keys(c).filter(k => c[k]).sort((a, z) => exc.includes(z) - exc.includes(a)).map(k => `<button type="button" aria-selected="${RC.filter === k}" data-rcf="${k}">${t('rr_' + k)} <span class="count">${c[k]}</span></button>`)).join('');
  const rows = d.rows.map(r => {
    const open = exc.includes(r.result) && !r.resolved_at;
    const orders = (r.orders || []).map(no => `<a href="#/order/${encodeURIComponent(no)}">${esc(no)}</a>`).join(' ');
    return `<tr class="${open ? 'rc-open' : ''}"><td class="muted">${r.source === 'phsar' ? t('rcFromPhsar') : r.line_no}</td>
      <td><code>${esc(r.tran_id || '')}</code>${r.kind === 'refund' ? ` <span class="chip">${t('rcRefund')}</span>` : ''}${r.reference ? `<br><span class="muted small">${esc(r.reference)}</span>` : ''}</td>
      <td class="small">${when(r.occurred_at)}${r.status ? `<br><span class="muted">${esc(r.status)}</span>` : ''}</td>
      <td>${r.amount != null ? amoney(r.amount) : '—'}${r.fee ? `<br><span class="muted small">${t('rcFees')} ${amoney(r.fee)}</span>` : ''}</td>
      <td>${r.expected != null ? amoney(r.expected) : '—'}</td>
      <td><span class="chip ${RC_CHIP[r.result] || ''}">${t('rr_' + r.result)}</span>${exc.includes(r.result) ? `<div class="muted small rc-why">${t('rx_' + r.result)}</div>` : ''}
        ${r.resolved_at ? `<div class="small"><b>${t('rcResolved')}:</b> ${esc(r.resolution)} <span class="muted">· ${esc(r.resolved_by || '')} · ${when(r.resolved_at)}</span></div>` : ''}
        ${open ? `<div class="row rc-res"><input id="rcn_${r.id}" name="note" maxlength="1000" placeholder="${t('rcResolveNote')}"><button type="button" data-rcres="${r.id}">${t('rcResolve')}</button></div>` : ''}</td>
      <td class="small">${orders || '—'}</td></tr>`;
  }).join('');
  $('#app').innerHTML = shell(`${pageHead(b.file_name, `${String(b.period_from || '—').slice(0, 10)} → ${String(b.period_to || '—').slice(0, 10)} · ${esc(b.uploaded_by || '')} · ${when(b.created_at)}`)}
    <p><a href="#/reconciliation">${t('rcBack')}</a></p>
    <section class="kpis"><div class="kpi"><span>${t('rcCollected')}</span><b>${amoney(b.collected)}</b></div><div class="kpi"><span>${t('rcRefunded')}</span><b>${amoney(b.refunded)}</b></div>
      <div class="kpi"><span>${t('rcFees')}</span><b>${amoney(b.fees)}</b></div><div class="kpi"><span>${t('rcNet')}</span><b>${amoney(b.net)}</b></div>
      <div class="kpi"><span>${t('rcPhsarPaid')}</span><b>${amoney(d.phsar_paid)}</b></div><div class="kpi"><span>${t('rcOpen')}</span><b class="${d.open ? 'warn' : ''}">${d.open}</b></div></section>
    <div class="row" style="justify-content:space-between;flex-wrap:wrap;gap:10px"><div class="tabs count">${tabs}</div>
      <div class="row"><button type="button" id="rcRematch">${t('rcRematch')}</button><button type="button" class="danger" id="rcDelete">${t('rcDelete')}</button></div></div>
    <section class="card">${rows ? `<div class="tbl-wrap"><table class="rc-table"><thead><tr><th>${t('rcLine')}</th><th>${t('rcTran')}</th><th>${t('rcWhen')}</th><th>${t('rcPayway')}</th>
      <th>${t('rcPhsar')}</th><th>${t('rcResult')}</th><th>${t('rcOrders')}</th></tr></thead><tbody>${rows}</tbody></table></div>` : `<p class="muted">${t('rcNoRows')}</p>`}</section>`);
  document.querySelectorAll('[data-rcf]').forEach(x => x.onclick = () => { RC.filter = x.dataset.rcf; reconDetail(id); });
  $('#rcRematch').onclick = async () => {
    try { const n = await api(`/admin/api/reconciliation/batches/${id}/rematch?filter=${RC.filter}`, {method: 'POST'}); toast(t('rcMatched')); reconDetail(id, n); }
    catch (e) { toast(errMsg(e), 'error'); }
  };
  $('#rcDelete').onclick = async () => {
    if (!confirm(t('rcConfirmDelete'))) return;
    try { await api(`/admin/api/reconciliation/batches/${id}`, {method: 'DELETE'}); location.hash = '#/reconciliation'; } catch (e) { toast(errMsg(e), 'error'); }
  };
  document.querySelectorAll('[data-rcres]').forEach(x => x.onclick = async () => {
    const rid = x.dataset.rcres;
    try { const n = await api(`/admin/api/reconciliation/rows/${rid}/resolve?filter=${RC.filter}`, {method: 'POST', body: {note: $('#rcn_' + rid).value}}); toast(t('rcResolved')); reconDetail(id, n); }
    catch (e) { showError(e, {note: '#rcn_' + rid}); }
  });
}
