'use strict';
// Admin: laporan pembeli (produk/toko) dan strike toko. 3 strike aktif (90 hari) → toko otomatis ditangguhkan.
addStrings({
  en: {rpHint: 'Buyer reports grouped by product or store. Check the product page and video, then dismiss, ban the product, or give the store a strike. 3 active strikes (90 days) suspend the store automatically.',
    rp_open: 'Open', rp_closed: 'Closed (60 days)', rpProduct: 'Product', rpStore: 'Store', rpReports: '{n} reports', rpReport1: '1 report', rpStrikes: 'Active strikes {n}/3',
    rpView: 'View on site', rpNote: 'Decision note (required — shown to the seller for bans and strikes)', rpDismiss: 'Dismiss', rpBan: 'Ban product', rpStrike: 'Give strike',
    rpAlsoStrike: 'Also give the store a strike', rpPoints: 'Points', rpConfirmBan: 'Ban this product? It disappears from Phsar and the seller is notified.',
    rpConfirmStrike: 'Give this store a strike? At 3 active strikes it is suspended.', rpDone: 'Done', rpHistory: 'Strike history', rpNoStrikes: 'No strikes.',
    rpRevoke: 'Revoke', rpRevokeNote: 'Why are you revoking this strike?', rpActive: 'Active', rpExpired: 'Expired', rpRevoked: 'Revoked', rpUntil: 'until {d}',
    rpManual: 'Give a strike without a report', rpReason: 'Reason', rpCausedSusp: 'caused suspension', rpBy: 'by {p}', rpResolution: 'Decision',
    rpRes_dismiss: 'Dismissed', rpRes_ban_product: 'Product banned', rpRes_strike: 'Strike given', rpBanned: 'banned', rpSuspended: 'suspended',
    q_reports: 'Buyer reports to review',
    rr_counterfeit: 'Counterfeit', rr_prohibited: 'Prohibited', rr_misleading: 'Misleading', rr_scam: 'Scam', rr_offensive: 'Offensive', rr_other: 'Other'},
  km: {rpHint: 'របាយការណ៍ពីអ្នកទិញ ដាក់ជាក្រុមតាមផលិតផល ឬហាង។ ពិនិត្យទំព័រផលិតផល និងវីដេអូ រួចមិនអើពើ ហាមផលិតផល ឬផ្តល់ការព្រមានដល់ហាង។ ការព្រមានសកម្ម 3 (90 ថ្ងៃ) ផ្អាកហាងដោយស្វ័យប្រវត្តិ។',
    rp_open: 'បើក', rp_closed: 'បានបិទ (60 ថ្ងៃ)', rpProduct: 'ផលិតផល', rpStore: 'ហាង', rpReports: 'របាយការណ៍ {n}', rpReport1: 'របាយការណ៍ 1', rpStrikes: 'ការព្រមានសកម្ម {n}/3',
    rpView: 'មើលលើគេហទំព័រ', rpNote: 'កំណត់ចំណាំការសម្រេច (ចាំបាច់ — អ្នកលក់មើលឃើញពេលហាម ឬព្រមាន)', rpDismiss: 'មិនអើពើ', rpBan: 'ហាមផលិតផល', rpStrike: 'ផ្តល់ការព្រមាន',
    rpAlsoStrike: 'ផ្តល់ការព្រមានដល់ហាងផងដែរ', rpPoints: 'ពិន្ទុ', rpConfirmBan: 'ហាមផលិតផលនេះ? វានឹងបាត់ពី Phsar ហើយអ្នកលក់នឹងទទួលដំណឹង។',
    rpConfirmStrike: 'ផ្តល់ការព្រមានដល់ហាងនេះ? ពេលមាន 3 ហាងនឹងត្រូវផ្អាក។', rpDone: 'រួចរាល់', rpHistory: 'ប្រវត្តិការព្រមាន', rpNoStrikes: 'គ្មានការព្រមាន។',
    rpRevoke: 'ដកចេញ', rpRevokeNote: 'ហេតុអ្វីបានជាដកការព្រមាននេះ?', rpActive: 'សកម្ម', rpExpired: 'ផុតកំណត់', rpRevoked: 'បានដក', rpUntil: 'រហូតដល់ {d}',
    rpManual: 'ផ្តល់ការព្រមានដោយគ្មានរបាយការណ៍', rpReason: 'មូលហេតុ', rpCausedSusp: 'បណ្តាលឱ្យផ្អាក', rpBy: 'ដោយ {p}', rpResolution: 'ការសម្រេច',
    rpRes_dismiss: 'មិនអើពើ', rpRes_ban_product: 'បានហាមផលិតផល', rpRes_strike: 'បានផ្តល់ការព្រមាន', rpBanned: 'បានហាម', rpSuspended: 'បានផ្អាក',
    q_reports: 'របាយការណ៍អ្នកទិញត្រូវពិនិត្យ',
    rr_counterfeit: 'ក្លែងក្លាយ', rr_prohibited: 'ហាមឃាត់', rr_misleading: 'បំភាន់', rr_scam: 'បោកប្រាស់', rr_offensive: 'មិនសមរម្យ', rr_other: 'ផ្សេងៗ'},
  zh: {rpHint: '买家举报按商品或店铺分组。查看商品页面和视频后，可以忽略、下架商品或给店铺警告。90 天内有效警告满 3 次，店铺自动暂停。',
    rp_open: '待处理', rp_closed: '已关闭（60 天）', rpProduct: '商品', rpStore: '店铺', rpReports: '{n} 条举报', rpReport1: '1 条举报', rpStrikes: '有效警告 {n}/3',
    rpView: '在网站上查看', rpNote: '处理说明（必填——下架或警告时卖家可见）', rpDismiss: '忽略', rpBan: '下架商品', rpStrike: '给予警告',
    rpAlsoStrike: '同时给店铺警告', rpPoints: '分数', rpConfirmBan: '下架此商品？商品将从 Phsar 消失，卖家会收到通知。',
    rpConfirmStrike: '给此店铺警告？累计 3 次有效警告将暂停店铺。', rpDone: '已完成', rpHistory: '警告记录', rpNoStrikes: '暂无警告。',
    rpRevoke: '撤销', rpRevokeNote: '为什么撤销此警告？', rpActive: '有效', rpExpired: '已过期', rpRevoked: '已撤销', rpUntil: '至 {d}',
    rpManual: '无举报直接警告', rpReason: '原因', rpCausedSusp: '导致暂停', rpBy: '由 {p}', rpResolution: '处理',
    rpRes_dismiss: '已忽略', rpRes_ban_product: '商品已下架', rpRes_strike: '已警告', rpBanned: '已下架', rpSuspended: '已暂停',
    q_reports: '待审核的买家举报',
    rr_counterfeit: '假冒', rr_prohibited: '违禁', rr_misleading: '误导', rr_scam: '诈骗', rr_offensive: '不当', rr_other: '其他'},
});

const RP = {status: 'open', open: {}};
const RP_REASONS = ['counterfeit', 'prohibited', 'misleading', 'scam', 'offensive', 'other'];

async function renderReports() {
  const d = await api('/admin/api/reports?status=' + RP.status);
  const open = RP.status === 'open';
  const day = s => s ? new Date(s).toLocaleDateString() : '—';
  const cards = d.groups.map((g, i) => {
    const isP = g.target_type === 'product';
    const reasons = {};
    g.reports.forEach(x => { reasons[x.reason] = (reasons[x.reason] || 0) + 1; });
    const link = isP ? `/en/p/${esc(g.product_slug || '')}` : `/en/s/${esc(g.store_slug)}`;
    const items = g.reports.map(x => `<li><b>${t('rr_' + x.reason)}</b> · ${esc(x.reporter)}${x.reporter_phone ? ' ' + esc(x.reporter_phone) : ''} · <span class="muted">${new Date(x.at).toLocaleString()}</span>
      ${x.note ? `<div style="white-space:pre-wrap">${esc(x.note)}</div>` : ''}${!open && x.resolution ? `<div class="muted small">${t('rpResolution')}: ${t('rpRes_' + x.resolution)} — ${esc(x.resolution_note || '')}</div>` : ''}</li>`).join('');
    const act = open ? `<div class="rt-act">
        <label>${t('rpNote')}<textarea id="rpn_${i}" name="note" maxlength="1000"></textarea></label>
        <div class="row" style="align-items:center">
          ${isP ? `<label class="inline"><input type="checkbox" id="rps_${i}"> ${t('rpAlsoStrike')}</label>` : ''}
          <label class="inline">${t('rpPoints')} <select id="rpp_${i}"><option>1</option><option>2</option><option>3</option></select></label></div>
        <div class="row"><button type="button" data-rp="dismiss" data-i="${i}">${t('rpDismiss')}</button>
          ${isP ? `<button type="button" class="danger" data-rp="ban_product" data-i="${i}">${t('rpBan')}</button>` : ''}
          <button type="button" class="danger" data-rp="strike" data-i="${i}">${t('rpStrike')}</button></div></div>` : '';
    return `<section class="card"><div class="spread"><div><span class="chip ${isP ? 'pink' : ''}">${t(isP ? 'rpProduct' : 'rpStore')}</span>
        <h2 style="margin-top:6px">${esc(isP ? (g.product_name || '—') : g.store_name)}${isP && g.product_status === 'banned' ? ` <span class="chip bad">${t('rpBanned')}</span>` : ''}</h2>
        <span class="muted">${isP ? esc(g.store_name) + ' · ' : ''}${g.count === 1 ? t('rpReport1') : t('rpReports', {n: g.count})} · ${day(g.first_at)}${g.last_at !== g.first_at ? ' – ' + day(g.last_at) : ''}</span></div>
        <div style="text-align:right"><span class="chip ${g.active_strikes >= 3 ? 'bad' : g.active_strikes ? 'warn' : ''}">${t('rpStrikes', {n: g.active_strikes})}</span>
        ${g.store_status === 'suspended' ? `<br><span class="chip bad" style="margin-top:4px">${t('rpSuspended')}</span>` : ''}</div></div>
      <div class="row" style="gap:6px">${Object.entries(reasons).map(([r, n]) => `<span class="chip">${t('rr_' + r)}${n > 1 ? ' × ' + n : ''}</span>`).join('')}</div>
      <div class="rt-says"><ul style="padding-left:18px;margin:0">${items}</ul></div>
      <div class="row"><a href="${link}" target="_blank" rel="noopener">${t('rpView')} ↗</a><button type="button" class="link" data-hist="${g.store_id}">${t('rpHistory')}</button></div>
      <div id="rph_${i}" data-store="${g.store_id}" ${RP.open[g.store_id] ? '' : 'hidden'}></div>
      ${act}</section>`;
  }).join('');
  const tabs = ['open', 'closed'].map(f => `<button type="button" aria-selected="${RP.status === f}" data-rps="${f}">${t('rp_' + f)}${f === 'open' && d.open ? ` <span class="count">${d.open}</span>` : ''}</button>`).join('');
  $('#app').innerHTML = shell(`${pageHead(t('i_reports'), t('rpHint'))}<div class="tabs count">${tabs}</div>
    ${cards || `<section class="card"><p class="muted">${t('empty')}</p></section>`}`);
  document.querySelectorAll('[data-rps]').forEach(b => b.onclick = () => { RP.status = b.dataset.rps; render(); });
  document.querySelectorAll('[data-hist]').forEach(b => b.onclick = () => {
    const id = b.dataset.hist;
    RP.open[id] = !RP.open[id];
    document.querySelectorAll(`[data-store="${id}"]`).forEach(box => { box.hidden = !RP.open[id]; if (RP.open[id]) loadStrikes(id, box); });
  });
  document.querySelectorAll('[data-store]').forEach(box => { if (RP.open[box.dataset.store]) loadStrikes(box.dataset.store, box); });
  document.querySelectorAll('[data-rp]').forEach(b => b.onclick = async () => {
    const i = +b.dataset.i, g = d.groups[i], action = b.dataset.rp;
    const also = action === 'ban_product' && $('#rps_' + i) && $('#rps_' + i).checked;
    if (action === 'ban_product' && !confirm(t('rpConfirmBan'))) return;
    if ((action === 'strike' || also) && !confirm(t('rpConfirmStrike'))) return;
    try {
      await api('/admin/api/reports/resolve', {method: 'POST', body: {target: g.target_type, id: g.target_id, action, note: $('#rpn_' + i).value, also_strike: also, points: +$('#rpp_' + i).value}});
      toast(t('rpDone')); SH.at = 0; await loadBadges(true); render();
    } catch (e) { showError(e, {note: '#rpn_' + i}); }
  });
}

async function loadStrikes(storeId, box) {
  let s;
  try { s = await api(`/admin/api/reports/stores/${storeId}/strikes`); } catch (e) { box.innerHTML = `<p class="muted">${esc(errMsg(e))}</p>`; return; }
  const rows = s.strikes.map(k => {
    const st = k.revoked_at ? 'rpRevoked' : k.active ? 'rpActive' : 'rpExpired';
    return `<li><span class="chip ${k.active ? 'warn' : ''}">${t(st)}</span> <b>${t('rr_' + k.reason)}</b> · ${k.points} · ${new Date(k.created_at).toLocaleDateString()}
      ${k.active ? `<span class="muted small">${t('rpUntil', {d: new Date(k.expires_at).toLocaleDateString()})}</span>` : ''}${k.caused_suspension ? ` <span class="chip bad">${t('rpCausedSusp')}</span>` : ''}${k.auto ? ` <span class="chip">${t('exAuto')}</span>` : ''}
      ${k.by ? `<span class="muted small">${t('rpBy', {p: esc(k.by)})}</span>` : ''}
      <div style="white-space:pre-wrap">${esc(k.note)}</div>${k.revoke_note ? `<div class="muted small">${t('rpRevoked')}: ${esc(k.revoke_note)}</div>` : ''}
      ${k.active ? `<button type="button" class="link" data-revoke="${k.id}">${t('rpRevoke')}</button>` : ''}</li>`;
  }).join('');
  box.innerHTML = `<div class="rt-says"><b>${t('rpHistory')} · ${esc(s.store.name)} · ${t('rpStrikes', {n: s.active_points})}</b>
    ${rows ? `<ul style="padding-left:18px;margin:6px 0 0">${rows}</ul>` : `<p class="muted">${t('rpNoStrikes')}</p>`}
    <details style="margin-top:8px"><summary>${t('rpManual')}</summary><div class="grid" style="margin-top:8px">
      <label>${t('rpReason')}<select id="ms_r_${storeId}">${RP_REASONS.map(r => `<option value="${r}">${t('rr_' + r)}</option>`).join('')}</select></label>
      <label>${t('rpPoints')}<select id="ms_p_${storeId}"><option>1</option><option>2</option><option>3</option></select></label></div>
      <label>${t('rpNote')}<textarea id="ms_n_${storeId}" name="note" maxlength="1000"></textarea></label>
      <button type="button" class="danger" id="ms_go_${storeId}">${t('rpStrike')}</button></details></div>`;
  box.querySelectorAll('[data-revoke]').forEach(b => b.onclick = async () => {
    const note = prompt(t('rpRevokeNote'));
    if (note == null) return;
    try { await api(`/admin/api/reports/strikes/${b.dataset.revoke}/revoke`, {method: 'POST', body: {note}}); toast(t('rpDone')); render(); }
    catch (e) { toast(errMsg(e), 'error'); }
  });
  $('#ms_go_' + storeId).onclick = async () => {
    if (!confirm(t('rpConfirmStrike'))) return;
    try {
      await api('/admin/api/reports/strikes', {method: 'POST', body: {store_id: storeId, reason: $('#ms_r_' + storeId).value, points: +$('#ms_p_' + storeId).value, note: $('#ms_n_' + storeId).value}});
      toast(t('rpDone')); render();
    } catch (e) { showError(e, {note: '#ms_n_' + storeId}); }
  };
}
