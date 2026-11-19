'use strict';
// Admin: iklan produk bersponsor — tinjau kampanye seller, jeda kampanye bermasalah, pengaturan slot (izin ads.review).
addStrings({
  en: {adHint: 'Sponsored products appear above search and category results. Approve only products that follow the listing rules — the photos, video (if any) and title must match the item.',
    ad_pending: 'To review', ad_running: 'Running', ad_all: 'All', adStore: 'Store', adBid: 'CPC', adBudget: 'Daily budget', adDates: 'Dates', adToday: 'Today', adTotal: 'All time',
    adNote: 'Note to seller (required to reject or pause)', adApprove: 'Approve', adReject: 'Reject', adPause: 'Pause', adResume: 'Resume', adDone: 'Saved', adLive: 'Showing now',
    adNotLive: 'Not showing (dates, budget or balance)', adSettings: 'Ad settings', adEnabled: 'Show sponsored products on Phsar', adSlots: 'Sponsored slots per page',
    adSaveSettings: 'Save settings', adViewProduct: 'View product', adStats: '{i} views · {c} clicks · {s}',
    adSt_pending: 'In review', adSt_approved: 'Approved', adSt_rejected: 'Rejected', adSt_paused: 'Paused', adSt_ended: 'Ended', adBySeller: 'by seller', adByPhsar: 'by Phsar',
    err_ad_state: 'Already handled — refreshing.'},
  km: {adHint: 'ផលិតផលបង់ផ្សាយបង្ហាញនៅលើលទ្ធផលស្វែងរក និងប្រភេទ។ សូមអនុម័តតែផលិតផលដែលគោរពច្បាប់ — វីដេអូ និងចំណងជើងត្រូវតែត្រូវនឹងទំនិញ។',
    ad_pending: 'ត្រូវពិនិត្យ', ad_running: 'កំពុងដំណើរការ', ad_all: 'ទាំងអស់', adStore: 'ហាង', adBid: 'CPC', adBudget: 'ថវិកាប្រចាំថ្ងៃ', adDates: 'កាលបរិច្ឆេទ', adToday: 'ថ្ងៃនេះ', adTotal: 'សរុប',
    adNote: 'កំណត់ចំណាំទៅអ្នកលក់ (ចាំបាច់ពេលបដិសេធ ឬផ្អាក)', adApprove: 'អនុម័ត', adReject: 'បដិសេធ', adPause: 'ផ្អាក', adResume: 'បន្ត', adDone: 'បានរក្សាទុក', adLive: 'កំពុងបង្ហាញ',
    adNotLive: 'មិនបង្ហាញ (កាលបរិច្ឆេទ ថវិកា ឬសមតុល្យ)', adSettings: 'ការកំណត់ការផ្សាយ', adEnabled: 'បង្ហាញផលិតផលបង់ផ្សាយនៅលើ Phsar', adSlots: 'ចំនួនកន្លែងផ្សាយក្នុងមួយទំព័រ',
    adSaveSettings: 'រក្សាទុកការកំណត់', adViewProduct: 'មើលផលិតផល', adStats: 'មើល {i} · ចុច {c} · {s}',
    adSt_pending: 'កំពុងពិនិត្យ', adSt_approved: 'បានអនុម័ត', adSt_rejected: 'បានបដិសេធ', adSt_paused: 'បានផ្អាក', adSt_ended: 'បានបញ្ចប់', adBySeller: 'ដោយអ្នកលក់', adByPhsar: 'ដោយ Phsar',
    err_ad_state: 'បានដោះស្រាយរួចហើយ — កំពុងផ្ទុកឡើងវិញ។'},
  zh: {adHint: '推广商品显示在搜索和分类结果上方。只批准符合上架规则的商品——照片、视频（如有）和标题必须与商品一致。',
    ad_pending: '待审核', ad_running: '投放中', ad_all: '全部', adStore: '店铺', adBid: '单次点击费', adBudget: '每日预算', adDates: '日期', adToday: '今天', adTotal: '累计',
    adNote: '给卖家的说明（拒绝或暂停时必填）', adApprove: '通过', adReject: '拒绝', adPause: '暂停', adResume: '恢复', adDone: '已保存', adLive: '正在展示',
    adNotLive: '未展示（日期、预算或余额）', adSettings: '广告设置', adEnabled: '在 Phsar 显示推广商品', adSlots: '每页推广位数量',
    adSaveSettings: '保存设置', adViewProduct: '查看商品', adStats: '{i} 次展示 · {c} 次点击 · {s}',
    adSt_pending: '审核中', adSt_approved: '已通过', adSt_rejected: '已拒绝', adSt_paused: '已暂停', adSt_ended: '已结束', adBySeller: '卖家暂停', adByPhsar: 'Phsar 暂停',
    err_ad_state: '已处理——正在刷新。'},
});

const AAD = {filter: 'pending'};
const AAD_CHIP = {pending: 'warn', approved: 'ok', rejected: 'bad', paused: '', ended: ''};

async function renderAds() {
  const {campaigns, counts, settings} = await api('/admin/api/ads?filter=' + AAD.filter);
  const m = v => '$' + Number(v || 0).toFixed(2);
  const stats = s => t('adStats', {i: s.impressions, c: s.clicks, s: m(s.spend)});
  const cards = campaigns.map(c => {
    const act = (a, label, cls = '') => `<button type="button" class="${cls}" data-adact="${a}" data-id="${c.id}">${label}</button>`;
    const actions = c.status === 'pending' ? act('approve', t('adApprove'), 'primary') + act('reject', t('adReject'), 'danger')
      : c.status === 'approved' ? act('pause', t('adPause'), 'danger') : c.status === 'paused' && c.paused_by === 'admin' ? act('resume', t('adResume')) : '';
    return `<section class="card"><div class="spread"><div class="row" style="gap:12px">${c.product_image ? `<img src="${esc(c.product_image)}" alt="" style="width:56px;height:56px;border-radius:12px;object-fit:cover">` : ''}
        <div><h2>${esc(c.product_name)}</h2><span class="muted">${t('adStore')} ${esc(c.store_name)} · <a href="#/product/${c.product_id}">${t('adViewProduct')}</a></span></div></div>
        <div style="text-align:right"><span class="chip ${AAD_CHIP[c.status]}">${t('adSt_' + c.status)}${c.paused_by ? ' · ' + t(c.paused_by === 'admin' ? 'adByPhsar' : 'adBySeller') : ''}</span>
        ${c.status === 'approved' ? `<br><small class="muted">${t(c.live ? 'adLive' : 'adNotLive')}</small>` : ''}</div></div>
      <div class="grid small"><p><span class="muted">${t('adBid')}</span><br><b>${m(c.bid)}</b></p><p><span class="muted">${t('adBudget')}</span><br>${m(c.daily_budget)}</p>
        <p><span class="muted">${t('adDates')}</span><br>${esc(c.starts_on)}${c.ends_on ? ' → ' + esc(c.ends_on) : ''}</p>
        <p><span class="muted">${t('adToday')}</span><br>${stats(c.today)}</p><p><span class="muted">${t('adTotal')}</span><br>${stats(c.total)}</p></div>
      ${c.review_note ? `<p class="muted small">${esc(c.review_note)}</p>` : ''}
      ${actions ? `<label>${t('adNote')}<input id="adn_${c.id}" name="note" maxlength="300"></label><div class="row">${actions}</div>` : ''}</section>`;
  }).join('');
  const tabs = ['pending', 'running', 'all'].map(f => `<button type="button" aria-selected="${AAD.filter === f}" data-adf="${f}">${t('ad_' + f)}${counts[f] ? ` <span class="count">${counts[f]}</span>` : ''}</button>`).join('');
  $('#app').innerHTML = shell(`${pageHead(t('i_ads'), t('adHint'))}
    <section class="card" id="adSet"><h2>${t('adSettings')}</h2><div class="row" style="flex-wrap:wrap;gap:16px;align-items:flex-end">
      <label class="inline"><input type="checkbox" id="adOn" ${settings.enabled ? 'checked' : ''}> ${t('adEnabled')}</label>
      <label>${t('adSlots')}<select id="adSlots" name="slots">${[1, 2, 3, 4].map(n => `<option ${n === settings.slots ? 'selected' : ''}>${n}</option>`).join('')}</select></label>
      <button type="button" id="adSave">${t('adSaveSettings')}</button></div></section>
    <div class="tabs count">${tabs}</div>
    ${cards || `<section class="card"><p class="muted">${t('empty')}</p></section>`}`);
  document.querySelectorAll('[data-adf]').forEach(b => b.onclick = () => { AAD.filter = b.dataset.adf; render(); });
  $('#adSave').onclick = async () => {
    try { await api('/admin/api/ads/settings', {method: 'PUT', body: {enabled: $('#adOn').checked, slots: Number($('#adSlots').value)}}); toast(t('adDone')); }
    catch (e) { showError(e, {slots: '#adSlots'}); }
  };
  document.querySelectorAll('[data-adact]').forEach(b => b.onclick = async () => {
    const id = b.dataset.id;
    try { await api(`/admin/api/ads/${id}/action`, {method: 'POST', body: {action: b.dataset.adact, note: $('#adn_' + id).value.trim()}}); toast(t('adDone')); render(); }
    catch (e) { if (e.code === 'ad_state') { toast(t('err_ad_state'), 'error'); render(); return; } showError(e, {note: '#adn_' + id}); }
  });
}
