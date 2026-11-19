'use strict';
// Admin (super admin): pekerjaan latar — status, riwayat run, jalankan sekarang, jeda/lanjutkan.
addStrings({
  en: {jbHint: 'Background workers that keep orders, payouts and notifications moving. Empty rounds are not logged — only the last heartbeat.',
    jbJob: 'Job', jbSchedule: 'Schedule', jbLast: 'Last activity', jbStatus: 'Status', jb24h: 'Last 24 h', jbRun: 'Run now', jbPause: 'Pause', jbResume: 'Resume',
    jbHistory: 'History', jbEvery: 'every {n}', jbSec: '{n} s', jbMin: '{n} min', jbPayoutSched: 'per payout settings',
    jbOk: 'Running', jbPaused: 'Paused', jbFailing: 'Failing', jbStale: 'No heartbeat', jbNever: 'Never ran',
    jbRuns: '{n} runs · {i} items', jbFailedN: '{n} failed', jbAvg: 'avg {ms} ms', jbItems: '{n} items', jbDone: 'Done — {n} items processed', jbErr: 'Finished with an error: {e}',
    jbNoRunner: 'Not available on this server', jbManualOnly: 'Run from the Payouts page', jbConfirmPause: 'Pause “{j}”? The worker skips it until you resume.',
    jbStaleHint: 'If this stays, the worker process may be stopped (WORKER_IN_SERVE / phsar worker).', jbNoRuns: 'No logged runs yet.', jbManual: 'manual', jbAuto: 'auto',
    jbCols: 'Started|Status|Items|Duration|Trigger|Error',
    jb_checkout_expire: 'Expire unpaid checkouts', jb_checkout_expire_d: 'Cancels KHQR checkouts that were not paid in time and returns the stock.',
    jb_order_accept_timeout: 'Seller accept limit', jb_order_accept_timeout_d: 'Cancels paid orders a seller did not accept in time and queues the refund.', jb_order_pack_timeout: 'Seller pack limit', jb_order_pack_timeout_d: 'Cancels accepted orders not marked ready for pickup in time and queues the refund.',
    jb_payout_complete: 'Complete delivered orders', jb_payout_complete_d: 'After the hold period, delivered orders become completed and credit the seller balance.',
    jb_payout_schedule: 'Scheduled payouts', jb_payout_schedule_d: 'Creates seller payouts on the day and hour set in Seller payouts.',
    jb_returns_escalate: 'Escalate unanswered returns', jb_returns_escalate_d: 'Sends return requests to Phsar review when the seller does not answer within 48 hours.',
    jb_notify_send: 'Send Telegram notifications', jb_notify_send_d: 'Delivers queued notifications to sellers, buyers and couriers on Telegram.',
    jb_translate_queue: 'AI translations', jb_translate_queue_d: 'Translates product and store texts into Khmer, English and Chinese.',
    jb_media_queue: 'Video & photo processing', jb_media_queue_d: 'Compresses uploaded videos and makes thumbnails.',
    jb_photo_index: 'Photo search index', jb_photo_index_d: 'Builds image vectors so buyers can search with a photo.'},
  km: {jbHint: 'ដំណើរការផ្ទៃខាងក្រោយដែលធ្វើឱ្យការបញ្ជាទិញ ការទូទាត់ និងការជូនដំណឹងដំណើរការ។ ជុំទទេមិនត្រូវកត់ត្រាទេ — មានតែសញ្ញាចុងក្រោយ។',
    jbJob: 'ការងារ', jbSchedule: 'កាលវិភាគ', jbLast: 'សកម្មភាពចុងក្រោយ', jbStatus: 'ស្ថានភាព', jb24h: '24 ម៉ោងចុងក្រោយ', jbRun: 'ដំណើរការឥឡូវ', jbPause: 'ផ្អាក', jbResume: 'បន្ត',
    jbHistory: 'ប្រវត្តិ', jbEvery: 'រៀងរាល់ {n}', jbSec: '{n} វិនាទី', jbMin: '{n} នាទី', jbPayoutSched: 'តាមការកំណត់ការទូទាត់',
    jbOk: 'កំពុងដំណើរការ', jbPaused: 'បានផ្អាក', jbFailing: 'បរាជ័យ', jbStale: 'គ្មានសញ្ញា', jbNever: 'មិនធ្លាប់ដំណើរការ',
    jbRuns: '{n} ដង · {i} ធាតុ', jbFailedN: 'បរាជ័យ {n}', jbAvg: 'មធ្យម {ms} ms', jbItems: '{n} ធាតុ', jbDone: 'រួចរាល់ — បានដំណើរការ {n} ធាតុ', jbErr: 'បញ្ចប់ដោយមានកំហុស៖ {e}',
    jbNoRunner: 'មិនមាននៅលើម៉ាស៊ីនមេនេះ', jbManualOnly: 'ដំណើរការពីទំព័រទូទាត់', jbConfirmPause: 'ផ្អាក «{j}»? ដំណើរការនឹងរំលងវារហូតដល់អ្នកបន្ត។',
    jbStaleHint: 'បើនៅតែបែបនេះ ដំណើរការ worker ប្រហែលជាឈប់ (WORKER_IN_SERVE / phsar worker)។', jbNoRuns: 'មិនទាន់មានការកត់ត្រា។', jbManual: 'ដោយដៃ', jbAuto: 'ស្វ័យប្រវត្តិ',
    jbCols: 'ចាប់ផ្តើម|ស្ថានភាព|ធាតុ|រយៈពេល|ប្រភព|កំហុស',
    jb_checkout_expire: 'បោះបង់ការទូទាត់មិនបានបង់', jb_checkout_expire_d: 'បោះបង់ checkout KHQR ដែលមិនបានបង់ទាន់ពេល ហើយដាក់ស្តុកវិញ។',
    jb_order_accept_timeout: 'ពេលកំណត់ទទួលរបស់អ្នកលក់', jb_order_accept_timeout_d: 'បោះបង់ការបញ្ជាទិញដែលអ្នកលក់មិនទទួលទាន់ពេល ហើយរៀបចំសងប្រាក់។', jb_order_pack_timeout: 'ពេលកំណត់វេចខ្ចប់របស់អ្នកលក់', jb_order_pack_timeout_d: 'បោះបង់ការបញ្ជាទិញដែលបានទទួល តែមិនបានសម្គាល់ថារួចរាល់ទាន់ពេល ហើយរៀបចំសងប្រាក់។',
    jb_payout_complete: 'បញ្ចប់ការបញ្ជាទិញដែលបានដឹក', jb_payout_complete_d: 'ក្រោយរយៈពេលរង់ចាំ ការបញ្ជាទិញដែលបានដឹកត្រូវបញ្ចប់ ហើយបញ្ចូលសមតុល្យអ្នកលក់។',
    jb_payout_schedule: 'ការទូទាត់តាមកាលវិភាគ', jb_payout_schedule_d: 'បង្កើតការទូទាត់អ្នកលក់តាមថ្ងៃ និងម៉ោងដែលបានកំណត់។',
    jb_returns_escalate: 'បញ្ជូនការប្ដូរដែលមិនបានឆ្លើយ', jb_returns_escalate_d: 'បញ្ជូនសំណើប្ដូរទៅ Phsar ពិនិត្យ នៅពេលអ្នកលក់មិនឆ្លើយក្នុង 48 ម៉ោង។',
    jb_notify_send: 'ផ្ញើការជូនដំណឹង Telegram', jb_notify_send_d: 'ផ្ញើការជូនដំណឹងទៅអ្នកលក់ អ្នកទិញ និងអ្នកដឹកតាម Telegram។',
    jb_translate_queue: 'ការបកប្រែ AI', jb_translate_queue_d: 'បកប្រែអត្ថបទផលិតផល និងហាងជាខ្មែរ អង់គ្លេស និងចិន។',
    jb_media_queue: 'ដំណើរការវីដេអូ និងរូបថត', jb_media_queue_d: 'បង្រួមវីដេអូដែលបានផ្ទុកឡើង និងបង្កើតរូបតូច។',
    jb_photo_index: 'សន្ទស្សន៍ស្វែងរកតាមរូប', jb_photo_index_d: 'បង្កើតវ៉ិចទ័ររូបភាព ដើម្បីឱ្យអ្នកទិញស្វែងរកដោយរូបថត។'},
  zh: {jbHint: '让订单、结算和通知持续运转的后台任务。空轮次不记录——只保留最后一次心跳。',
    jbJob: '任务', jbSchedule: '频率', jbLast: '最近活动', jbStatus: '状态', jb24h: '最近 24 小时', jbRun: '立即运行', jbPause: '暂停', jbResume: '恢复',
    jbHistory: '历史', jbEvery: '每 {n}', jbSec: '{n} 秒', jbMin: '{n} 分钟', jbPayoutSched: '按结算设置',
    jbOk: '运行中', jbPaused: '已暂停', jbFailing: '失败', jbStale: '无心跳', jbNever: '从未运行',
    jbRuns: '{n} 次 · {i} 项', jbFailedN: '{n} 次失败', jbAvg: '平均 {ms} ms', jbItems: '{n} 项', jbDone: '完成——已处理 {n} 项', jbErr: '运行出错：{e}',
    jbNoRunner: '此服务器不可用', jbManualOnly: '请在结算页面运行', jbConfirmPause: '暂停“{j}”？恢复前后台将跳过该任务。',
    jbStaleHint: '如果持续如此，后台进程可能已停止（WORKER_IN_SERVE / phsar worker）。', jbNoRuns: '暂无运行记录。', jbManual: '手动', jbAuto: '自动',
    jbCols: '开始|状态|项目|耗时|触发|错误',
    jb_checkout_expire: '取消未付款结账', jb_checkout_expire_d: '取消未按时付款的 KHQR 结账并退回库存。',
    jb_order_accept_timeout: '卖家接单时限', jb_order_accept_timeout_d: '取消卖家未及时接单的已付款订单并安排退款。', jb_order_pack_timeout: '卖家打包时限', jb_order_pack_timeout_d: '取消已接单但未及时标记待取件的订单并安排退款。',
    jb_payout_complete: '完成已送达订单', jb_payout_complete_d: '保留期结束后，已送达订单变为已完成并计入卖家余额。',
    jb_payout_schedule: '定时结算', jb_payout_schedule_d: '按卖家结算中设置的日期和时间生成结算。',
    jb_returns_escalate: '升级未回复的退货', jb_returns_escalate_d: '卖家 48 小时内未回复时，将退货申请转交 Phsar 审核。',
    jb_notify_send: '发送 Telegram 通知', jb_notify_send_d: '通过 Telegram 向卖家、买家和快递员发送排队的通知。',
    jb_translate_queue: 'AI 翻译', jb_translate_queue_d: '将商品和店铺文字翻译成高棉语、英语和中文。',
    jb_media_queue: '视频与照片处理', jb_media_queue_d: '压缩上传的视频并生成缩略图。',
    jb_photo_index: '以图搜索索引', jb_photo_index_d: '生成图片向量，让买家可以用照片搜索。'},
});

const AJ = {open: null};
const JB_CHIP = {ok: 'ok', paused: '', failing: 'bad', stale: 'warn', never: ''};

function jbKey(name) { return 'jb_' + name.replace(/\./g, '_'); }
function jbEvery(j) {
  if (j.name === 'payout.schedule') return t('jbPayoutSched');
  return t('jbEvery', {n: j.every_sec >= 60 ? t('jbMin', {n: Math.round(j.every_sec / 60)}) : t('jbSec', {n: j.every_sec})});
}
// Status: dijeda · gagal (error terakhir lebih baru dari sukses) · tanpa detak (> 3× jeda + 2 menit) · berjalan · belum pernah.
function jbStatus(j) {
  if (j.paused) return 'paused';
  const s = j.state;
  if (!s || !s.last_tick_at) return 'never';
  if (s.last_error_at && (!s.last_success_at || new Date(s.last_error_at) >= new Date(s.last_success_at))) return 'failing';
  const limit = (Math.max(j.every_sec, 30) * 3 + 120) * 1000;
  if (Date.now() - new Date(s.last_tick_at) > limit) return 'stale';
  return 'ok';
}

async function renderJobs() {
  const {jobs} = await api('/admin/api/jobs');
  const stat = {ok: 'jbOk', paused: 'jbPaused', failing: 'jbFailing', stale: 'jbStale', never: 'jbNever'};
  const rows = jobs.map(j => {
    const st = jbStatus(j), s = j.state || {}, lr = j.last_run;
    const err = st === 'failing' && s.last_error ? `<div class="small bad-text" style="color:var(--bad);overflow-wrap:anywhere">${esc(s.last_error)}</div>` : '';
    const run = j.name === 'payout.schedule' ? `<span class="muted small">${t('jbManualOnly')}</span>`
      : j.runnable ? `<button type="button" data-jrun="${esc(j.name)}">${t('jbRun')}</button>` : `<span class="muted small">${t('jbNoRunner')}</span>`;
    return `<tr>
      <td><b>${t(jbKey(j.name))}</b><div class="muted small">${t(jbKey(j.name) + '_d')}</div><code class="small muted">${esc(j.name)}</code></td>
      <td class="small">${jbEvery(j)}</td>
      <td class="small">${s.last_tick_at ? `${ago(s.last_tick_at)}` : '—'}${lr ? `<div class="muted">${esc(t('jbItems', {n: lr.items ?? 0}))} · ${ago(lr.started_at)}</div>` : ''}</td>
      <td><span class="chip ${JB_CHIP[st]}">${t(stat[st])}</span>${err}${st === 'stale' ? `<div class="muted small">${t('jbStaleHint')}</div>` : ''}</td>
      <td class="small">${t('jbRuns', {n: j.runs_24h, i: j.items_24h})}${j.failed_24h ? `<div style="color:var(--bad)">${t('jbFailedN', {n: j.failed_24h})}</div>` : ''}${j.avg_ms != null ? `<div class="muted">${t('jbAvg', {ms: j.avg_ms})}</div>` : ''}</td>
      <td><div class="row" style="gap:6px;flex-wrap:wrap">${run}
        <button type="button" data-jpause="${esc(j.name)}" data-p="${j.paused ? '0' : '1'}">${t(j.paused ? 'jbResume' : 'jbPause')}</button>
        <button type="button" data-jhist="${esc(j.name)}" aria-expanded="${AJ.open === j.name}">${t('jbHistory')}</button></div></td></tr>
      ${AJ.open === j.name ? `<tr class="jb-hist"><td colspan="6" id="jbh_${esc(j.name).replace(/\./g, '_')}"><span class="muted">…</span></td></tr>` : ''}`;
  }).join('');
  const failing = jobs.filter(j => jbStatus(j) === 'failing').length, paused = jobs.filter(j => j.paused).length;
  $('#app').innerHTML = shell(`${pageHead(t('i_jobs'), t('jbHint'))}
    <section class="kpis"><div class="kpi"><span>${t('jbOk')}</span><b>${jobs.filter(j => jbStatus(j) === 'ok').length}/${jobs.length}</b></div>
      <div class="kpi"><span>${t('jbFailing')}</span><b class="${failing ? 'warn' : ''}">${failing}</b></div>
      <div class="kpi"><span>${t('jbPaused')}</span><b>${paused}</b></div></section>
    <section class="card"><div class="tbl-wrap"><table><thead><tr><th>${t('jbJob')}</th><th>${t('jbSchedule')}</th><th>${t('jbLast')}</th><th>${t('jbStatus')}</th><th>${t('jb24h')}</th><th></th></tr></thead>
      <tbody>${rows}</tbody></table></div></section>`);
  document.querySelectorAll('[data-jrun]').forEach(b => b.onclick = async () => {
    b.disabled = true;
    try {
      const r = await api(`/admin/api/jobs/${b.dataset.jrun}/run`, {method: 'POST'});
      toast(r.ok ? t('jbDone', {n: r.items}) : t('jbErr', {e: r.error}), r.ok ? 'success' : 'error');
      renderJobs();
    } catch (e) { b.disabled = false; toast(errMsg(e), 'error'); }
  });
  document.querySelectorAll('[data-jpause]').forEach(b => b.onclick = async () => {
    const pause = b.dataset.p === '1';
    if (pause && !confirm(t('jbConfirmPause', {j: t(jbKey(b.dataset.jpause))}))) return;
    try { await api(`/admin/api/jobs/${b.dataset.jpause}/pause`, {method: 'POST', body: {paused: pause}}); renderJobs(); } catch (e) { toast(errMsg(e), 'error'); }
  });
  document.querySelectorAll('[data-jhist]').forEach(b => b.onclick = () => { AJ.open = AJ.open === b.dataset.jhist ? null : b.dataset.jhist; renderJobs(); });
  if (AJ.open) {
    const cell = document.getElementById('jbh_' + AJ.open.replace(/\./g, '_'));
    try {
      const {runs} = await api(`/admin/api/jobs/${AJ.open}/runs`);
      const cols = t('jbCols').split('|');
      if (cell) cell.innerHTML = runs.length ? `<table style="min-width:0;width:100%"><thead><tr>${cols.map(c => `<th>${c}</th>`).join('')}</tr></thead><tbody>${runs.map(r => `<tr>
        <td class="small">${new Date(r.started_at).toLocaleString()}</td><td><span class="chip ${r.status === 'failed' ? 'bad' : 'ok'}">${esc(r.status)}</span></td>
        <td>${r.items ?? 0}</td><td class="small">${r.duration_ms ?? '—'} ms</td>
        <td class="small">${t(r.trigger === 'manual' ? 'jbManual' : 'jbAuto')}${r.triggered_by ? ' · ' + esc(r.triggered_by) : ''}</td>
        <td class="small" style="overflow-wrap:anywhere;max-width:320px">${esc(r.error || '')}</td></tr>`).join('')}</tbody></table>` : `<p class="muted">${t('jbNoRuns')}</p>`;
    } catch (e) { if (cell) cell.textContent = errMsg(e); }
  }
}
