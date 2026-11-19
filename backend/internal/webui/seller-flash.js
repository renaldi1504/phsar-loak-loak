'use strict';
// Seller Center (pemilik): flash sale — terima undangan Phsar, daftarkan produk ke sesi yang dibuka, pantau status.
addStrings({
  en: {tabFlash: 'Flash sale', sfIntro: 'Flash sales run in sessions. Join with a lower price and a quota: your products appear in the Flash Sale section on the home page.',
    sfInvites: 'Invitations from Phsar', sfOpen: 'Sessions open for registration', sfMine: 'Your flash sale products', sfNone: 'No flash sale sessions are open right now.',
    sfNoItems: 'You have not joined a flash sale yet.', sfWhen: '{a} → {b}', sfRegUntil: 'Register before {d}', sfYouPay: 'You pay {s}% of the discount, Phsar pays {p}%',
    sfPhsarPays: 'Phsar pays the whole discount — you still receive the normal price', sfSellerPays: 'You pay the discount (your earnings use the flash price)',
    sfMin: 'At least {p}% off · max {n} per buyer', sfRegister: 'Register products', sfFind: 'Search your products', sfSubmit: 'Send {n} products for review', sfSent: 'Sent — Phsar will review.',
    sfAccept: 'Accept', sfDecline: 'Decline', sfWithdraw: 'Withdraw', sfWithdrawConfirm: 'Withdraw this product from the flash sale?', sfDeclineConfirm: 'Decline this invitation?',
    sfNormal: 'Normal', sfFlash: 'Flash price', sfQuota: 'Quota', sfSold: 'Sold', sfYouGet: 'You receive about {a} each', sfSession: 'Session',
    sfI_invited: 'Waiting for you', sfI_pending: 'In review', sfI_approved: 'Joined', sfI_rejected: 'Not accepted', sfI_declined: 'Declined', sfI_withdrawn: 'Withdrawn', sfI_removed: 'Removed by Phsar',
    sfStock: 'stock {n}', sfCancel: 'Cancel', err_flash_state: 'This can no longer be changed.'},
  km: {tabFlash: 'Flash sale', sfIntro: 'Flash sale ដំណើរការជាវគ្គ។ ចូលរួមដោយតម្លៃទាប និងកូតា: ផលិតផលរបស់អ្នកនឹងបង្ហាញក្នុងផ្នែក Flash Sale នៅទំព័រដើម។',
    sfInvites: 'ការអញ្ជើញពី Phsar', sfOpen: 'វគ្គបើកចុះឈ្មោះ', sfMine: 'ផលិតផល Flash sale របស់អ្នក', sfNone: 'មិនមានវគ្គ Flash sale បើកនៅពេលនេះទេ។',
    sfNoItems: 'អ្នកមិនទាន់ចូលរួម Flash sale ទេ។', sfWhen: '{a} → {b}', sfRegUntil: 'ចុះឈ្មោះមុន {d}', sfYouPay: 'អ្នកចេញ {s}% នៃការបញ្ចុះ, Phsar ចេញ {p}%',
    sfPhsarPays: 'Phsar ចេញការបញ្ចុះទាំងអស់ — អ្នកនៅតែទទួលតម្លៃធម្មតា', sfSellerPays: 'អ្នកចេញការបញ្ចុះ (ប្រាក់ចំណូលគិតតាមតម្លៃ flash)',
    sfMin: 'បញ្ចុះយ៉ាងតិច {p}% · អតិបរមា {n} ក្នុងម្នាក់', sfRegister: 'ចុះឈ្មោះផលិតផល', sfFind: 'ស្វែងរកផលិតផលរបស់អ្នក', sfSubmit: 'ផ្ញើផលិតផល {n} ដើម្បីពិនិត្យ', sfSent: 'បានផ្ញើ — Phsar នឹងពិនិត្យ។',
    sfAccept: 'ទទួល', sfDecline: 'បដិសេធ', sfWithdraw: 'ដកចេញ', sfWithdrawConfirm: 'ដកផលិតផលនេះចេញពី Flash sale?', sfDeclineConfirm: 'បដិសេធការអញ្ជើញនេះ?',
    sfNormal: 'ធម្មតា', sfFlash: 'តម្លៃ flash', sfQuota: 'កូតា', sfSold: 'លក់បាន', sfYouGet: 'អ្នកទទួលប្រហែល {a} ក្នុងមួយ', sfSession: 'វគ្គ',
    sfI_invited: 'រង់ចាំអ្នក', sfI_pending: 'កំពុងពិនិត្យ', sfI_approved: 'បានចូលរួម', sfI_rejected: 'មិនត្រូវបានទទួល', sfI_declined: 'បានបដិសេធ', sfI_withdrawn: 'បានដកចេញ', sfI_removed: 'Phsar បានយកចេញ',
    sfStock: 'ស្តុក {n}', sfCancel: 'បោះបង់', err_flash_state: 'មិនអាចផ្លាស់ប្តូរបានទៀតទេ។'},
  zh: {tabFlash: '限时抢购', sfIntro: '限时抢购按场次进行。以更低价格和配额参加：您的商品会出现在首页的限时抢购区。',
    sfInvites: 'Phsar 的邀请', sfOpen: '开放报名的场次', sfMine: '您的抢购商品', sfNone: '目前没有开放报名的场次。',
    sfNoItems: '您还没有参加限时抢购。', sfWhen: '{a} → {b}', sfRegUntil: '请在 {d} 前报名', sfYouPay: '您承担折扣的 {s}%，Phsar 承担 {p}%',
    sfPhsarPays: '折扣全部由 Phsar 承担——您仍按原价结算', sfSellerPays: '折扣由您承担（收入按抢购价计算）',
    sfMin: '至少减 {p}% · 每人限购 {n}', sfRegister: '报名商品', sfFind: '搜索您的商品', sfSubmit: '提交 {n} 个商品审核', sfSent: '已提交——Phsar 将审核。',
    sfAccept: '接受', sfDecline: '拒绝', sfWithdraw: '撤回', sfWithdrawConfirm: '从限时抢购撤回此商品？', sfDeclineConfirm: '拒绝此邀请？',
    sfNormal: '原价', sfFlash: '抢购价', sfQuota: '配额', sfSold: '已售', sfYouGet: '每件约收入 {a}', sfSession: '场次',
    sfI_invited: '等待您确认', sfI_pending: '审核中', sfI_approved: '已参加', sfI_rejected: '未通过', sfI_declined: '已拒绝', sfI_withdrawn: '已撤回', sfI_removed: 'Phsar 已移除',
    sfStock: '库存 {n}', sfCancel: '取消', err_flash_state: '已无法更改。'},
});

const SF = {data: null, reg: null, picked: [], results: []};
const SF_CHIP = {invited: 'warn', pending: 'warn', approved: 'ok', rejected: 'bad', declined: '', withdrawn: '', removed: 'bad'};

function sfName(s) { return (P.lang === 'km' && s.name_km) || (P.lang === 'zh' && s.name_zh) || s.name_en; }
function sfFund(s) {
  return s.phsar_share_pct === 100 ? t('sfPhsarPays') : s.phsar_share_pct === 0 ? t('sfSellerPays') : t('sfYouPay', {s: 100 - s.phsar_share_pct, p: s.phsar_share_pct});
}
// perkiraan penerimaan seller per barang sebelum komisi: harga flash + bagian Phsar dari potongan
function sfGet(it, s) { const n = parseFloat(it.normal_price), f = parseFloat(it.sale_price); return (f + (n - f) * s.phsar_share_pct / 100).toFixed(2); }

async function tabFlash(id) {
  SF.data = await api(`/v1/seller/stores/${id}/flash`);
  drawSellerFlash(id);
}

function drawSellerFlash(id) {
  const d = SF.data, byId = {};
  d.sessions.forEach(x => { byId[x.session.id] = x.session; });
  const invites = d.items.filter(i => i.status === 'invited');
  const mine = d.items.filter(i => i.status !== 'invited');
  const when = s => t('sfWhen', {a: fmtTime(s.starts_at), b: fmtTime(s.ends_at)});
  // sellerpages: satu baris produk yang muat di layar HP (foto · nama · harga normal → flash · kuota · status)
  const itemLine = (i, actions, withSession = true) => {
    const s = byId[i.session_id];
    return `<div class="sp-prod"><img src="${esc(i.image_url || '')}" alt="" onerror="this.style.visibility='hidden'"><div class="px">
      <b>${esc(i.product_name)}${i.options ? ` <span class="muted">${esc(i.options)}</span>` : ''}</b>
      ${withSession && s ? `<span class="muted">${esc(sfName(s))} · ${when(s)}</span>` : ''}
      <span class="sp-price"><s>${omoney(i.normal_price)}</s> → <b>${omoney(i.sale_price)}</b> (-${i.off_pct}%)</span>
      <span class="muted">${t('sfQuota')} ${i.qty_cap} · ${t('sfSold')} ${i.sold}</span>
      ${s ? `<span class="sp-get">${t('sfYouGet', {a: '$' + sfGet(i, s)})}</span>` : ''}${i.note ? `<span class="muted">${esc(i.note)}</span>` : ''}
      <span class="chip ${SF_CHIP[i.status]}">${t('sfI_' + i.status)}</span>${actions}</div></div>`;
  };
  const open = d.sessions.filter(x => x.can_register);
  const groups = {};
  invites.forEach(i => { (groups[i.session_id] = groups[i.session_id] || []).push(i); });
  const step = invites.length || SF.reg ? (SF.picked.length ? 1 : 0) : -1;
  $('#tab').innerHTML = `${spSteps([[t('spFl1'), t('spFl1W')], [t('spFl2'), t('spFl2W')], [t('spFl3'), t('spFl3W')]], step)}
    <p class="ohelp">⚡ ${t('sfIntro')}</p>
    ${invites.length ? `<h2 class="sp-sub">${t('sfInvites')}<span class="count">${invites.length}</span></h2>
      ${Object.keys(groups).sort().map(sid => { const s = byId[sid]; return `<section class="card">
        ${s ? `<div class="sp-ses"><h3>${esc(sfName(s))}</h3><span class="muted">${when(s)}</span><div class="sp-chips"><span class="chip pink">${esc(sfFund(s))}</span></div></div>` : ''}
        <p class="muted" style="margin:0">${t('spFlInvW')}</p>
        ${groups[sid].map(i => itemLine(i, `<div class="sp-acts"><button type="button" class="primary sp-block" data-sfact="accept" data-id="${i.id}">${t('sfAccept')}</button><button type="button" class="textbtn" data-sfact="decline" data-id="${i.id}">${t('sfDecline')}</button></div>`, false)).join('')}</section>`; }).join('')}` : ''}
    <h2 class="sp-sub">${t('sfOpen')}</h2>
    ${open.length ? open.map(x => { const s = x.session, on = SF.reg === s.id; return `<section class="card sp-ses">
      <h3>${esc(sfName(s))}</h3><span class="muted">${when(s)}</span>
      ${x.register_until ? spDue(x.register_until, 7 * 24 * 60, 'spFlRegLeft', 'err_flash_state') : ''}
      <div class="sp-chips"><span class="chip pink">${esc(sfFund(s))}</span><span class="chip">${t('sfMin', {p: s.min_discount_pct, n: s.per_buyer})}</span></div>
      <button type="button" class="${on ? '' : 'primary'} sp-block" data-sfreg="${s.id}">${on ? t('spFlClose') : t('sfRegister')}</button>
      ${on ? sfPickerHTML(s) : ''}</section>`; }).join('') : spEmpty('⚡', t('sfNone'), t('spFlNoneW'))}
    <h2 class="sp-sub">${t('sfMine')}</h2>
    <section class="card">${mine.length ? mine.map(i => {
      const s = byId[i.session_id], canW = (i.status === 'pending' || i.status === 'approved') && s && new Date(s.starts_at) > new Date();
      return itemLine(i, canW ? `<div class="sp-acts" style="align-items:flex-start"><button type="button" class="textbtn" data-sfact="withdraw" data-id="${i.id}">${t('sfWithdraw')}</button></div>` : '');
    }).join('') : `<p class="muted sp-center" style="margin:0">${t('sfNoItems')} ${t('spFlNoMine')}</p>`}</section>`;
  document.querySelectorAll('[data-sfact]').forEach(b => b.onclick = async () => {
    const a = b.dataset.sfact;
    if (a === 'withdraw' && !confirm(t('sfWithdrawConfirm'))) return;
    if (a === 'decline' && !confirm(t('sfDeclineConfirm'))) return;
    try { SF.data = await api(`/v1/seller/stores/${id}/flash/items/${b.dataset.id}/${a}`, {method: 'POST'}); toast(t('saved')); drawSellerFlash(id); loadBell(); }
    catch (e) { toast(errMsg(e), 'error'); if (e.code === 'flash_state') tabFlash(id); }
  });
  document.querySelectorAll('[data-sfreg]').forEach(b => b.onclick = () => { SF.reg = SF.reg === b.dataset.sfreg ? null : b.dataset.sfreg; SF.picked = []; SF.results = []; drawSellerFlash(id); if (SF.reg) sfSearch(id, ''); });
  if (SF.reg) bindSellerPicker(id, byId[SF.reg]);
}

function sfPickerHTML(s) {
  // sellerpages: kartu per produk (bukan tabel) supaya muat di HP; atribut data-sp/name tetap sama.
  return `<div class="sf-picker" id="sfPicker"><p class="muted" style="margin:0">${t('spFlFindHint')}</p>
    <label>${t('sfFind')}<input id="sfFind" type="search" autocomplete="off"></label><div id="sfResults" class="fs-results"></div>
    ${SF.picked.length ? `<h3 class="sp-sub">${t('spFlPicked', {n: SF.picked.length})}</h3><div class="sp-picked">
      ${SF.picked.map((v, n) => `<div class="sp-pick"><div class="ph"><div><b>${esc(v.name)}</b>${v.options ? ` <span class="muted">${esc(v.options)}</span>` : ''}
          <div class="muted small">${t('sfNormal')} ${omoney(v.price)}${phsLowTag(v.low30)} · ${t('sfStock', {n: v.stock})}</div></div>
          <button type="button" data-sfun="${n}" aria-label="remove">✕</button></div>
        <div class="pg"><label>${t('sfFlash')}<input class="fs-in" data-sp="${n}" name="sale_price" value="${esc(v.sale_price)}" inputmode="decimal"></label>
          <label>${t('sfQuota')}<input class="fs-in sm" data-sp="${n}" name="qty_cap" value="${v.qty_cap}" inputmode="numeric"></label></div></div>`).join('')}</div>
      <div class="sp-acts"><button type="button" class="primary sp-block" id="sfSubmit">${t('sfSubmit', {n: SF.picked.length})}</button><button type="button" class="textbtn" id="sfCancel">${t('sfCancel')}</button></div>` : ''}</div>`;
}

async function sfSearch(id, q) {
  try {
    // pricehistory: at = jam mulai sesi → low30 (harga terendah 30 hari) per varian
    const ses = SF.data.sessions.find(x => x.session.id === SF.reg), at = ses ? '&at=' + encodeURIComponent(ses.session.starts_at) : '';
    const r = await api(`/v1/seller/stores/${id}/flash/variants?q=${encodeURIComponent(q)}${at}`);
    SF.results = r.variants; PHS.rule = r.low30_rule !== false; sfDrawResults(id);
  } catch (e) { toast(errMsg(e), 'error'); }
}

function sfDrawResults(id) {
  const s = SF.data.sessions.find(x => x.session.id === SF.reg).session;
  const taken = new Set(SF.picked.map(v => v.variant_id).concat(SF.data.items.filter(i => i.session_id === s.id && ['invited', 'pending', 'approved'].includes(i.status)).map(i => i.variant_id)));
  const box = $('#sfResults'); if (!box) return;
  box.innerHTML = SF.results.map((v, i) => `<label class="fs-res"><input type="checkbox" data-sfres="${i}" ${taken.has(v.variant_id) ? 'checked disabled' : ''}>
    <span><b>${esc(v.name)}</b>${v.options ? ` <span class="muted">${esc(v.options)}</span>` : ''}<span class="muted small"> · ${omoney(v.price)} · ${t('sfStock', {n: v.stock})}${v.low30 ? ' · ' + t('phLow30', {a: parseFloat(v.low30).toFixed(2)}) : ''}</span></span></label>`).join('');
  document.querySelectorAll('[data-sfres]').forEach(cb => cb.onchange = () => {
    if (!cb.checked) return;
    sfKeep();
    const v = SF.results[+cb.dataset.sfres];
    SF.picked.push(Object.assign({}, v, {sale_price: phsSuggest(v.price, s.min_discount_pct, v.low30), qty_cap: Math.min(v.stock, 10)}));
    const q = $('#sfFind').value; drawSellerFlash(id); $('#sfFind').value = q; sfDrawResults(id);
  });
}

function sfKeep() { document.querySelectorAll('#sfPicker [data-sp]').forEach(x => { SF.picked[+x.dataset.sp][x.name] = x.value; }); }

function bindSellerPicker(id, s) {
  let timer = null;
  $('#sfFind').oninput = () => { clearTimeout(timer); timer = setTimeout(() => sfSearch(id, $('#sfFind').value.trim()), 250); };
  if (SF.results.length) sfDrawResults(id);
  document.querySelectorAll('[data-sfun]').forEach(b => b.onclick = () => { sfKeep(); SF.picked.splice(+b.dataset.sfun, 1); drawSellerFlash(id); });
  if ($('#sfCancel')) $('#sfCancel').onclick = () => { SF.reg = null; SF.picked = []; drawSellerFlash(id); };
  if ($('#sfSubmit')) $('#sfSubmit').onclick = async () => {
    sfKeep();
    const items = SF.picked.map(v => ({variant_id: v.variant_id, sale_price: String(v.sale_price).trim(), qty_cap: parseInt(v.qty_cap, 10) || 0}));
    try {
      SF.data = await api(`/v1/seller/stores/${id}/flash/sessions/${s.id}/items`, {method: 'POST', body: {items}});
      SF.reg = null; SF.picked = []; toast(t('sfSent')); drawSellerFlash(id);
    } catch (e) {
      const m = /^items\.(\d+)\.(\w+)$/.exec(e.field || '');
      toast(errMsg(e), 'error');
      if (m) { const el = document.querySelector(`#sfPicker [data-sp="${m[1]}"][name="${m[2]}"]`); if (el) { el.classList.add('invalid'); el.focus(); } }
    }
  };
}
