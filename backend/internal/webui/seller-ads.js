'use strict';
// Seller Center (pemilik): iklan produk bersponsor — bayar per klik dari saldo Phsar.
addStrings({
  en: {tabAds: 'Ads', adsIntro: 'Show a product at the top of search and category pages. You pay only when a buyer clicks (once per buyer per day), from your Phsar balance.',
    adsOff: 'Ads are paused on Phsar right now.', adsBalance: 'Available balance: {a}', adsNew: 'New campaign', adsProduct: 'Product', adsChoose: 'Choose a live product',
    adsBid: 'Cost per click ($)', adsBidHint: '{min}–{max}. Higher bids show first.', adsBudget: 'Daily budget ($)', adsBudgetHint: 'At least {min}. Ads stop for the day when it is used up.',
    adsStart: 'Start', adsEnd: 'End (optional)', adsCreate: 'Send for review', adsCreated: 'Sent — Phsar reviews ads within 1 working day.', adsNone: 'No campaigns yet.',
    adsToday: 'Today', adsTotal: 'All time', adsImpr: 'Views', adsClicks: 'Clicks', adsSpend: 'Spent', adsCtr: 'CTR', adsLive: 'Showing now', adsNotLive: 'Not showing',
    adsNotLiveWhy: 'Check dates, today’s budget and your balance.', adsPause: 'Pause', adsResume: 'Resume', adsEnd2: 'End', adsSave: 'Save', adsSaved: 'Saved',
    adsConfirmEnd: 'End this campaign? You can create a new one later.', adsNoProducts: 'You need an approved, active product first.',
    adsSt_pending: 'In review', adsSt_approved: 'Approved', adsSt_rejected: 'Not approved', adsSt_paused: 'Paused', adsSt_ended: 'Ended', adsPausedByPhsar: 'Paused by Phsar', adsPhsar: 'Phsar',
    err_ad_state: 'This campaign cannot do that right now.'},
  km: {tabAds: 'ផ្សាយពាណិជ្ជកម្ម', adsIntro: 'បង្ហាញផលិតផលនៅខាងលើទំព័រស្វែងរក និងប្រភេទ។ អ្នកបង់តែពេលអ្នកទិញចុច (ម្តងក្នុងមួយថ្ងៃក្នុងអ្នកទិញម្នាក់) ពីសមតុល្យ Phsar របស់អ្នក។',
    adsOff: 'ការផ្សាយពាណិជ្ជកម្មត្រូវបានផ្អាកនៅលើ Phsar ពេលនេះ។', adsBalance: 'សមតុល្យអាចប្រើបាន៖ {a}', adsNew: 'យុទ្ធនាការថ្មី', adsProduct: 'ផលិតផល', adsChoose: 'ជ្រើសផលិតផលដែលកំពុងលក់',
    adsBid: 'តម្លៃក្នុងមួយការចុច ($)', adsBidHint: '{min}–{max}។ តម្លៃខ្ពស់បង្ហាញមុន។', adsBudget: 'ថវិកាប្រចាំថ្ងៃ ($)', adsBudgetHint: 'យ៉ាងតិច {min}។ ការផ្សាយឈប់សម្រាប់ថ្ងៃនោះពេលអស់ថវិកា។',
    adsStart: 'ចាប់ផ្តើម', adsEnd: 'បញ្ចប់ (ស្រេចចិត្ត)', adsCreate: 'ផ្ញើឱ្យពិនិត្យ', adsCreated: 'បានផ្ញើ — Phsar ពិនិត្យក្នុង 1 ថ្ងៃធ្វើការ។', adsNone: 'មិនទាន់មានយុទ្ធនាការ។',
    adsToday: 'ថ្ងៃនេះ', adsTotal: 'សរុប', adsImpr: 'ការមើល', adsClicks: 'ការចុច', adsSpend: 'ចំណាយ', adsCtr: 'CTR', adsLive: 'កំពុងបង្ហាញ', adsNotLive: 'មិនកំពុងបង្ហាញ',
    adsNotLiveWhy: 'សូមពិនិត្យកាលបរិច្ឆេទ ថវិកាថ្ងៃនេះ និងសមតុល្យ។', adsPause: 'ផ្អាក', adsResume: 'បន្ត', adsEnd2: 'បញ្ចប់', adsSave: 'រក្សាទុក', adsSaved: 'បានរក្សាទុក',
    adsConfirmEnd: 'បញ្ចប់យុទ្ធនាការនេះ? អ្នកអាចបង្កើតថ្មីពេលក្រោយ។', adsNoProducts: 'អ្នកត្រូវមានផលិតផលដែលបានអនុម័ត និងកំពុងលក់ជាមុនសិន។',
    adsSt_pending: 'កំពុងពិនិត្យ', adsSt_approved: 'បានអនុម័ត', adsSt_rejected: 'មិនអនុម័ត', adsSt_paused: 'បានផ្អាក', adsSt_ended: 'បានបញ្ចប់', adsPausedByPhsar: 'Phsar បានផ្អាក', adsPhsar: 'Phsar',
    err_ad_state: 'យុទ្ធនាការនេះមិនអាចធ្វើបែបនេះបានទេពេលនេះ។'},
  zh: {tabAds: '广告', adsIntro: '让商品显示在搜索和分类页面顶部。只有买家点击时才付费（每位买家每天一次），从您的 Phsar 余额扣除。',
    adsOff: 'Phsar 目前暂停了广告。', adsBalance: '可用余额：{a}', adsNew: '新建推广', adsProduct: '商品', adsChoose: '选择在售商品',
    adsBid: '每次点击费用（$）', adsBidHint: '{min}–{max}。出价高的优先显示。', adsBudget: '每日预算（$）', adsBudgetHint: '至少 {min}。当天预算用完后停止展示。',
    adsStart: '开始', adsEnd: '结束（可选）', adsCreate: '提交审核', adsCreated: '已提交——Phsar 会在 1 个工作日内审核。', adsNone: '暂无推广。',
    adsToday: '今天', adsTotal: '累计', adsImpr: '展示', adsClicks: '点击', adsSpend: '花费', adsCtr: '点击率', adsLive: '正在展示', adsNotLive: '未展示',
    adsNotLiveWhy: '请检查日期、今日预算和余额。', adsPause: '暂停', adsResume: '继续', adsEnd2: '结束', adsSave: '保存', adsSaved: '已保存',
    adsConfirmEnd: '结束此推广？之后可以新建。', adsNoProducts: '您需要先有一个已通过审核的在售商品。',
    adsSt_pending: '审核中', adsSt_approved: '已通过', adsSt_rejected: '未通过', adsSt_paused: '已暂停', adsSt_ended: '已结束', adsPausedByPhsar: 'Phsar 已暂停', adsPhsar: 'Phsar',
    err_ad_state: '此推广目前无法执行该操作。'},
});

const SAD = {chip: {pending: 'warn', approved: 'ok', rejected: 'bad', paused: '', ended: ''}, form: false}; // sellerpages: form = formulir kampanye baru dibuka

function adCtr(s) { return s.impressions ? (100 * s.clicks / s.impressions).toFixed(1) + '%' : '—'; }

async function tabAds(id) {
  const [d, prods, pay] = await Promise.all([
    api(`/v1/seller/stores/${id}/ads`),
    api(`/v1/seller/stores/${id}/products`).catch(() => ({products: []})),
    api(`/v1/seller/stores/${id}/payouts`).catch(() => null),
  ]);
  const taken = new Set(d.campaigns.filter(c => ['pending', 'approved', 'paused'].includes(c.status)).map(c => c.product_id));
  const live = prods.products.filter(p => p.status === 'active' && p.moderation_status === 'approved' && !taken.has(p.id));
  const today = new Date().toLocaleDateString('en-CA', {timeZone: 'Asia/Phnom_Penh'});
  // sellerpages: kartu kampanye — angka hari ini sebagai 3 kotak, pengaturan harga/anggaran dilipat.
  const cards = d.campaigns.map(c => {
    const open = ['pending', 'approved', 'paused'].includes(c.status);
    const byPhsar = c.status === 'paused' && c.paused_by === 'admin';
    return `<section class="card order" data-ad="${c.id}">
      <div class="spread"><div class="row" style="gap:12px;flex-wrap:nowrap;min-width:0">${c.product_image ? `<img class="logo" src="${esc(c.product_image)}" alt="">` : ''}
        <div style="min-width:0"><h2>${esc(c.product_name)}</h2><span class="muted">${esc(c.starts_on)}${c.ends_on ? ' → ' + esc(c.ends_on) : ''}</span></div></div>
        <span class="chip ${SAD.chip[c.status]}">${t(byPhsar ? 'adsPausedByPhsar' : 'adsSt_' + c.status)}</span></div>
      ${c.status === 'approved' ? `<p class="ohelp" style="margin:0">${c.live ? '🟢 ' + t('adsLive') : '⚪ ' + t('adsNotLive') + ' — ' + t('adsNotLiveWhy')}</p>` : ''}
      ${c.review_note ? `<div class="banner ${c.status === 'rejected' || byPhsar ? 'bad' : 'info'}"><b>${t('adsPhsar')}:</b> ${esc(c.review_note)}</div>` : ''}
      <div><b style="font-size:13px">${t('adsToday')}</b><div class="sp-stats"><div><span>${t('adsImpr')}</span><b>${c.today.impressions}</b></div>
        <div><span>${t('adsClicks')}</span><b>${c.today.clicks}</b></div><div><span>${t('adsSpend')}</span><b>${omoney(c.today.spend)}</b></div></div></div>
      <p class="muted small" style="margin:0">${t('spAdAll', {v: c.total.impressions, c: c.total.clicks, s: omoney(c.total.spend), r: adCtr(c.total)})}</p>
      ${open ? `<div class="sp-row2">${c.status === 'approved' ? `<button type="button" data-act="pause">${t('adsPause')}</button>` : ''}
          ${c.status === 'paused' && !byPhsar ? `<button type="button" class="primary" data-act="resume">${t('adsResume')}</button>` : ''}
          <button type="button" class="danger" data-act="end">${t('adsEnd2')}</button></div>
        <details class="sp-more"><summary>${t('spAdChange')}</summary><div>
          <label>${t('adsBid')}<input name="bid" inputmode="decimal" value="${esc(c.bid)}"></label>
          <label>${t('adsBudget')}<input name="daily_budget" inputmode="decimal" value="${esc(c.daily_budget)}"></label>
          <label>${t('adsEnd')}<input name="ends_on" type="date" min="${today}" value="${esc(c.ends_on || '')}"></label>
          <button type="button" class="primary sp-block" data-act="save">${t('adsSave')}</button></div></details>` : ''}
    </section>`;
  }).join('');
  const showForm = SAD.form || !d.campaigns.length;
  const form = live.length ? `<section class="card"><h2>${t('adsNew')}</h2><p class="muted" style="margin:0">${t('adsIntro')}</p>
      <div class="grid" id="adForm">
        <label>${t('adsProduct')}<select name="product_id"><option value="">${t('adsChoose')}</option>${live.map(p => `<option value="${p.id}">${esc(p.name)}</option>`).join('')}</select></label>
        <label>${t('adsBid')}<input name="bid" inputmode="decimal" value="0.10"><small class="muted">${t('adsBidHint', {min: '$' + d.min_bid, max: '$' + d.max_bid})}</small></label>
        <label>${t('adsBudget')}<input name="daily_budget" inputmode="decimal" value="2.00"><small class="muted">${t('adsBudgetHint', {min: '$' + d.min_budget})}</small></label>
        <label>${t('adsStart')}<input name="starts_on" type="date" min="${today}" value="${today}"></label>
        <label>${t('adsEnd')}<input name="ends_on" type="date" min="${today}"></label></div>
      <div class="sp-acts" data-errors-anchor><button type="button" class="primary sp-block" id="adCreate">${t('adsCreate')}</button>
        ${d.campaigns.length ? `<button type="button" class="textbtn" id="spAdHide">${t('cancel')}</button>` : ''}</div></section>`
    : spEmpty('📦', t('adsNoProducts'), '', `<a class="btn primary" href="#/store/${id}/product/new">${t('su_prod')}</a>`);
  $('#tab').innerHTML = `${spSteps([[t('spAd1'), t('spAd1W')], [t('spAd2'), t('spAd2W')], [t('spAd3'), t('spAd3W')]], showForm && live.length ? 0 : -1)}
    ${d.enabled ? '' : `<div class="banner warn">${t('adsOff')}</div>`}
    ${pay ? `<div class="sp-balance"><span>${t('spAdBal')}</span><b>${omoney(pay.available)}</b><small>${t('spAdBalW')}</small></div>` : ''}
    ${d.campaigns.length ? `<h2 class="sp-sub">${t('spAdYours')}</h2>${cards}
      ${showForm ? form : live.length ? `<button type="button" class="primary sp-block" id="spAdNew">+ ${t('adsNew')}</button>` : ''}`
    : `${spEmpty('📣', t('spAdNoneT'), t('spAdNoneW'))}${form}`}`;
  if ($('#spAdNew')) $('#spAdNew').onclick = () => { SAD.form = true; tabAds(id); };
  if ($('#spAdHide')) $('#spAdHide').onclick = () => { SAD.form = false; tabAds(id); };
  const val = (box, n) => (box.querySelector(`[name="${n}"]`) || {}).value || '';
  const fields = sel => ({product_id: sel + ' [name=product_id]', bid: sel + ' [name=bid]', daily_budget: sel + ' [name=daily_budget]', starts_on: sel + ' [name=starts_on]', ends_on: sel + ' [name=ends_on]'});
  if ($('#adCreate')) $('#adCreate').onclick = async () => {
    const f = $('#adForm');
    const body = {product_id: val(f, 'product_id'), bid: val(f, 'bid').trim(), daily_budget: val(f, 'daily_budget').trim(), starts_on: val(f, 'starts_on')};
    if (val(f, 'ends_on')) body.ends_on = val(f, 'ends_on');
    try { await api(`/v1/seller/stores/${id}/ads`, {method: 'POST', body}); toast(t('adsCreated')); SAD.form = false; tabAds(id); }
    catch (e) { showError(e, fields('#adForm')); }
  };
  document.querySelectorAll('[data-ad] [data-act]').forEach(b => b.onclick = async () => {
    const box = b.closest('[data-ad]'), cid = box.dataset.ad, act = b.dataset.act;
    if (act === 'end' && !confirm(t('adsConfirmEnd'))) return;
    const body = act === 'save' ? {bid: val(box, 'bid').trim(), daily_budget: val(box, 'daily_budget').trim(), ends_on: val(box, 'ends_on')} : {action: act};
    try { await api(`/v1/seller/stores/${id}/ads/${cid}`, {method: 'PATCH', body}); toast(t('adsSaved')); tabAds(id); }
    catch (e) { const dt = box.querySelector('details'); if (dt && act === 'save') dt.open = true; showError(e, fields(`[data-ad="${cid}"]`)); }
  });
}
