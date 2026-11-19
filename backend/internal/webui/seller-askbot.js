'use strict';
// Seller Center (askbot): "Jawaban cepat & FAQ" — tanya-jawab toko (opsional per produk) + "Coba tanya bot".
addStrings({
  en: {mi_faq: 'Quick replies & FAQ', mi_faqD: 'Common answers, and a helper that drafts replies to buyers',
    abIntro: 'Write answers to questions buyers ask often. The reply helper uses them — together with your product details — to draft answers for you. Questions for one product also appear on its page.',
    abTryTitle: 'Try the reply helper', abTryHint: 'Type a question like a buyer would, in any language. The helper drafts a reply in the same language. Nothing is sent to anyone.',
    abQuestion: 'Buyer’s question', abQuestionPh: 'e.g. Do you still have size M?', abAbout: 'About which product?', abNoProduct: 'No product (general question)',
    abDraft: 'Draft a reply', abDrafting: 'Writing…', abReply: 'Suggested reply', abCopy: 'Copy', abCopied: 'Copied',
    abCheck: 'Check the reply before you send it. The helper only uses your product details and answers below — if it doesn’t know, it says you will check.',
    abOff: 'The reply helper is turned off by Phsar right now. You can still save your answers.', abRules: 'Simple mode: replies are built from your product details and answers (no AI).',
    abListTitle: 'Your answers', abNone: 'No answers yet. Start with delivery, opening hours or how to measure.', abAdd: 'Add a question',
    abWhole: 'Whole shop', abForProduct: 'For a product', abEdit: 'Edit', abDelete: 'Delete', abConfirmDel: 'Delete this question and answer?',
    abFormNew: 'New question & answer', abFormEdit: 'Edit question & answer', abQ: 'Question', abA: 'Answer', abLang: 'Written in', abFor: 'Applies to',
    abSave: 'Save', abSaved: 'Saved', abDeleted: 'Deleted', abIdeas: 'Ideas:',
    abIdeaShip: 'How long does delivery take?', abIdeaHours: 'What time are you open?', abIdeaSize: 'How do I choose my size?', abIdeaPay: 'How can I pay?', abIdeaOrig: 'Is it original?',
    abCount: '{n} of {max}', err_assist_unavailable: 'The reply helper is turned off right now.', err_rate_limited: 'Too many drafts this hour. Please try again later.'},
  km: {mi_faq: 'ចម្លើយរហ័ស និង FAQ', mi_faqD: 'ចម្លើយទូទៅ និងជំនួយការសរសេរព្រាងចម្លើយឱ្យអ្នកទិញ',
    abIntro: 'សរសេរចម្លើយសម្រាប់សំណួរដែលអ្នកទិញសួរញឹកញាប់។ ជំនួយការប្រើវា រួមជាមួយព័ត៌មានផលិតផល ដើម្បីព្រាងចម្លើយឱ្យអ្នក។ សំណួរសម្រាប់ផលិតផលមួយក៏បង្ហាញនៅទំព័រផលិតផលនោះដែរ។',
    abTryTitle: 'សាកល្បងជំនួយការឆ្លើយ', abTryHint: 'វាយសំណួរដូចអ្នកទិញ ជាភាសាណាក៏បាន។ ជំនួយការព្រាងចម្លើយជាភាសាដដែល។ គ្មានអ្វីត្រូវបានផ្ញើទៅនរណាទេ។',
    abQuestion: 'សំណួររបស់អ្នកទិញ', abQuestionPh: 'ឧ. នៅមានទំហំ M ទេ?', abAbout: 'អំពីផលិតផលណា?', abNoProduct: 'គ្មានផលិតផល (សំណួរទូទៅ)',
    abDraft: 'ព្រាងចម្លើយ', abDrafting: 'កំពុងសរសេរ…', abReply: 'ចម្លើយដែលណែនាំ', abCopy: 'ចម្លង', abCopied: 'បានចម្លង',
    abCheck: 'សូមពិនិត្យចម្លើយមុនពេលផ្ញើ។ ជំនួយការប្រើតែព័ត៌មានផលិតផល និងចម្លើយខាងក្រោមប៉ុណ្ណោះ — បើមិនដឹង វានឹងប្រាប់ថាអ្នកនឹងពិនិត្យ។',
    abOff: 'ជំនួយការឆ្លើយត្រូវបាន Phsar បិទនៅពេលនេះ។ អ្នកនៅតែអាចរក្សាទុកចម្លើយបាន។', abRules: 'របៀបសាមញ្ញ៖ ចម្លើយបង្កើតពីព័ត៌មានផលិតផល និងចម្លើយរបស់អ្នក (គ្មាន AI)។',
    abListTitle: 'ចម្លើយរបស់អ្នក', abNone: 'មិនទាន់មានចម្លើយទេ។ ចាប់ផ្តើមពីការដឹកជញ្ជូន ម៉ោងបើក ឬរបៀបវាស់ទំហំ។', abAdd: 'បន្ថែមសំណួរ',
    abWhole: 'ហាងទាំងមូល', abForProduct: 'សម្រាប់ផលិតផល', abEdit: 'កែ', abDelete: 'លុប', abConfirmDel: 'លុបសំណួរ និងចម្លើយនេះ?',
    abFormNew: 'សំណួរ និងចម្លើយថ្មី', abFormEdit: 'កែសំណួរ និងចម្លើយ', abQ: 'សំណួរ', abA: 'ចម្លើយ', abLang: 'សរសេរជា', abFor: 'អនុវត្តចំពោះ',
    abSave: 'រក្សាទុក', abSaved: 'បានរក្សាទុក', abDeleted: 'បានលុប', abIdeas: 'គំនិត៖',
    abIdeaShip: 'ដឹកជញ្ជូនរយៈពេលប៉ុន្មាន?', abIdeaHours: 'បើកម៉ោងប៉ុន្មាន?', abIdeaSize: 'តើត្រូវជ្រើសទំហំយ៉ាងដូចម្តេច?', abIdeaPay: 'អាចបង់ប្រាក់តាមរបៀបណា?', abIdeaOrig: 'ជារបស់ពិតទេ?',
    abCount: '{n} ក្នុង {max}', err_assist_unavailable: 'ជំនួយការឆ្លើយត្រូវបានបិទនៅពេលនេះ។', err_rate_limited: 'ព្រាងច្រើនពេកក្នុងម៉ោងនេះ។ សូមព្យាយាមម្តងទៀតពេលក្រោយ។'},
  zh: {mi_faq: '快捷回复与常见问题', mi_faqD: '常用回答，以及帮你起草买家回复的助手',
    abIntro: '为买家常问的问题写好答案。回复助手会结合你的商品信息，用这些答案帮你起草回复。针对某个商品的问答也会显示在该商品页面上。',
    abTryTitle: '试试回复助手', abTryHint: '像买家一样输入问题，任何语言都可以。助手会用同一种语言起草回复，不会发送给任何人。',
    abQuestion: '买家的问题', abQuestionPh: '例如：M 码还有吗？', abAbout: '关于哪个商品？', abNoProduct: '不指定商品（一般问题）',
    abDraft: '起草回复', abDrafting: '正在撰写…', abReply: '建议回复', abCopy: '复制', abCopied: '已复制',
    abCheck: '发送前请先检查回复。助手只使用你的商品信息和下方的答案——不确定时会说明你会再确认。',
    abOff: 'Phsar 目前已关闭回复助手。你仍然可以保存答案。', abRules: '简易模式：根据你的商品信息和答案生成回复（不使用 AI）。',
    abListTitle: '你的答案', abNone: '还没有答案。可以先写配送、营业时间或尺码测量。', abAdd: '添加问题',
    abWhole: '全店', abForProduct: '指定商品', abEdit: '编辑', abDelete: '删除', abConfirmDel: '删除这条问答？',
    abFormNew: '新问答', abFormEdit: '编辑问答', abQ: '问题', abA: '回答', abLang: '书写语言', abFor: '适用于',
    abSave: '保存', abSaved: '已保存', abDeleted: '已删除', abIdeas: '参考：',
    abIdeaShip: '配送需要多久？', abIdeaHours: '几点营业？', abIdeaSize: '怎么选尺码？', abIdeaPay: '可以怎么付款？', abIdeaOrig: '是正品吗？',
    abCount: '{n} / {max}', err_assist_unavailable: '回复助手目前已关闭。', err_rate_limited: '本小时起草次数过多，请稍后再试。'},
});

const AB = {products: null, data: null};
const AB_LANGS = {en: 'English', km: 'ខ្មែរ', zh: '中文'};

async function abTab(id) {
  const [d, pr] = await Promise.all([api(`/v1/seller/stores/${id}/assist/faqs`), api(`/v1/seller/stores/${id}/products`)]);
  AB.products = pr.products.filter(p => p.status !== 'archived' && p.status !== 'banned');
  AB.data = d;
  const opts = (sel) => `<option value="">${t('abNoProduct')}</option>` + AB.products.map(p => `<option value="${p.id}" ${sel === p.id ? 'selected' : ''}>${esc(p.name || '—')}</option>`).join('');
  const item = f => `<article class="card" style="padding:14px 16px;gap:6px" lang="${f.lang}">
      <div class="spread" style="align-items:flex-start;gap:8px"><b style="overflow-wrap:anywhere">${esc(f.question)}</b>
        <span class="chip">${AB_LANGS[f.lang] || f.lang}</span></div>
      <p style="margin:0;white-space:pre-line;overflow-wrap:anywhere">${esc(f.answer)}</p>
      <div class="row" style="justify-content:space-between;align-items:center"><span class="chip ${f.product_id ? 'pink' : ''}">${f.product_id ? '📦 ' + esc(f.product_name || t('abForProduct')) : '🏪 ' + t('abWhole')}</span>
        <span class="row"><button type="button" class="small" data-abedit="${f.id}">${t('abEdit')}</button><button type="button" class="small danger" data-abdel="${f.id}">${t('abDelete')}</button></span></div></article>`;
  $('#tab').innerHTML = `
    <section class="card"><p style="margin:0">${t('abIntro')}</p></section>
    <section class="card" id="abTry"><h2>💬 ${t('abTryTitle')}</h2>
      ${!d.enabled ? `<div class="banner info">${t('abOff')}</div>` : `<p class="muted">${t('abTryHint')}</p>${d.ai ? '' : `<p class="ohelp">${t('abRules')}</p>`}
      <label>${t('abQuestion')}<textarea id="abQ" maxlength="1000" rows="2" placeholder="${esc(t('abQuestionPh'))}"></textarea></label>
      <label>${t('abAbout')}<select id="abP">${opts('')}</select></label>
      <div data-errors-anchor><button type="button" class="primary" id="abGo">${t('abDraft')}</button></div>
      <div id="abOut"></div>`}
    </section>
    <section class="card"><div class="spread"><h2>${t('abListTitle')}</h2><span class="muted small">${t('abCount', {n: d.faqs.length, max: d.max})}</span></div>
      ${d.faqs.length ? '' : `<p class="muted">${t('abNone')}</p>`}
      <button type="button" class="primary" id="abAddBtn" ${d.faqs.length >= d.max ? 'disabled' : ''}>＋ ${t('abAdd')}</button></section>
    ${d.faqs.map(item).join('')}`;
  if ($('#abGo')) $('#abGo').onclick = () => abSuggest(id);
  $('#abAddBtn').onclick = () => abForm(id, null);
  document.querySelectorAll('[data-abedit]').forEach(b => b.onclick = () => abForm(id, d.faqs.find(f => f.id === b.dataset.abedit)));
  document.querySelectorAll('[data-abdel]').forEach(b => b.onclick = async () => {
    if (!confirm(t('abConfirmDel'))) return;
    try { await api(`/v1/seller/stores/${id}/assist/faqs/${b.dataset.abdel}`, {method: 'DELETE'}); toast(t('abDeleted')); abTab(id); } catch (e) { toast(errMsg(e), 'error'); }
  });
}

async function abSuggest(id) {
  const btn = $('#abGo'), q = $('#abQ').value.trim();
  if (!checkRequired([['#abQ', t('required'), () => q.length > 0]])) return;
  btn.disabled = true; btn.textContent = t('abDrafting');
  try {
    const body = {question: q};
    if ($('#abP').value) body.product_id = $('#abP').value;
    const r = await api(`/v1/seller/stores/${id}/assist/suggest`, {method: 'POST', body});
    $('#abOut').innerHTML = `<div class="nextcard" style="margin-top:12px;display:flex;flex-direction:column;gap:8px">
        <b>${t('abReply')}</b><p id="abReplyText" style="margin:0;white-space:pre-line;overflow-wrap:anywhere">${esc(r.reply)}</p>
        <div><button type="button" id="abCopy">${t('abCopy')}</button></div><p class="muted small" style="margin:0">${t('abCheck')}</p></div>`;
    $('#abCopy').onclick = async () => { try { await navigator.clipboard.writeText(r.reply); toast(t('abCopied')); } catch { toast(r.reply); } };
  } catch (e) { showError(e, {question: '#abQ', product_id: '#abP'}); }
  finally { btn.disabled = false; btn.textContent = t('abDraft'); }
}

function abForm(id, f) {
  const prods = AB.products || [];
  const ideas = ['abIdeaShip', 'abIdeaHours', 'abIdeaSize', 'abIdeaPay', 'abIdeaOrig'];
  $('#tab').innerHTML = `<a href="javascript:void 0" id="abBack">← ${t('mi_faq')}</a>
    <section class="card" id="abFormCard"><h2>${t(f ? 'abFormEdit' : 'abFormNew')}</h2>
      ${f ? '' : `<div class="row" style="flex-wrap:wrap;gap:6px;align-items:center"><span class="muted small">${t('abIdeas')}</span>${ideas.map(k => `<button type="button" class="small" data-abidea="${k}">${t(k)}</button>`).join('')}</div>`}
      <label>${t('abQ')}<input id="abFQ" maxlength="300" value="${esc(f ? f.question : '')}"></label>
      <label>${t('abA')}<textarea id="abFA" maxlength="2000" rows="5">${esc(f ? f.answer : '')}</textarea></label>
      <label>${t('abLang')}<select id="abFL">${Object.entries(AB_LANGS).map(([k, v]) => `<option value="${k}" ${(f ? f.lang : P.lang) === k ? 'selected' : ''}>${v}</option>`).join('')}</select></label>
      <label>${t('abFor')}<select id="abFP"><option value="">🏪 ${t('abWhole')}</option>${prods.map(p => `<option value="${p.id}" ${f && f.product_id === p.id ? 'selected' : ''}>📦 ${esc(p.name || '—')}</option>`).join('')}</select></label>
      <div class="row" data-errors-anchor><button type="button" class="primary" id="abFSave">${t('abSave')}</button><button type="button" id="abFCancel">${t('cancel')}</button></div>
    </section>`;
  const back = () => abTab(id);
  $('#abBack').onclick = back; $('#abFCancel').onclick = back;
  document.querySelectorAll('[data-abidea]').forEach(b => b.onclick = () => { $('#abFQ').value = t(b.dataset.abidea); $('#abFL').value = P.lang; $('#abFA').focus(); });
  $('#abFSave').onclick = async () => {
    if (!checkRequired([['#abFQ', t('required'), el => el.value.trim().length > 0], ['#abFA', t('required'), el => el.value.trim().length > 0]])) return;
    const body = {question: $('#abFQ').value, answer: $('#abFA').value, lang: $('#abFL').value, product_id: $('#abFP').value || null};
    try {
      await api(`/v1/seller/stores/${id}/assist/faqs${f ? '/' + f.id : ''}`, {method: f ? 'PUT' : 'POST', body});
      toast(t('abSaved')); back();
    } catch (e) { showError(e, {question: '#abFQ', answer: '#abFA', lang: '#abFL', product_id: '#abFP'}); }
  };
}
