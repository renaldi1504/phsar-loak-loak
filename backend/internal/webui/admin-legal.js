'use strict';
// Super admin: dokumen legal (3 bahasa, berversi) + identitas & kontak perusahaan yang mengisi {{placeholder}} di dokumen.
addStrings({
  en: {lgCompanyTitle: 'Company & contact', lgCompanyDesc: 'Shown on the Contact page and filled into every document ({{company}}, {{email}} …). Fill these before launch.',
    lgCompany: 'Company / legal name', lgAddress: 'Address', lgEmail: 'Support email', lgPhone: 'Phone', lgTelegram: 'Telegram (@username)', lgFacebook: 'Facebook page (https://…)', lgHours: 'Opening hours',
    lgSaveCompany: 'Save contact details', lgMissing: 'Still empty: {f}. Documents show “—” there until you fill it in.',
    lgDocs: 'Documents', lgCurrent: 'Current version {v} · since {d}', lgView: 'View on site', lgHistory: 'Version history',
    lgNewVersion: 'New version number', lgVersionHint: 'Each publish is kept forever as evidence. Use a new number, e.g. 1.1.',
    lgPreview: 'Preview', lgPublish: 'Publish new version', lgPublished: 'Published version {v}', lgConfirm: 'Publish version {v} of “{d}” in all three languages now?',
    lgFormat: 'Format: ## heading, ### sub-heading, - bullet, 1. numbered, **bold**, [text](https://…). Placeholders: {{company}} {{address}} {{email}} {{phone}} {{telegram}} {{site}} {{base}}.',
    lgChanged: 'edited', lgDraftWarn: 'Have a lawyer review these templates before launch — they are a starting point, not legal advice.',
    legal_buyer_terms: 'Terms of use', legal_privacy_policy: 'Privacy policy', legal_return_policy: 'Returns & refunds', legal_seller_agreement: 'Seller agreement', legal_data_deletion: 'Data deletion'},
  km: {lgCompanyTitle: 'ក្រុមហ៊ុន និងទំនាក់ទំនង', lgCompanyDesc: 'បង្ហាញនៅទំព័រទំនាក់ទំនង ហើយបំពេញក្នុងឯកសារទាំងអស់ ({{company}}, {{email}} …)។ សូមបំពេញមុនបើកដំណើរការ។',
    lgCompany: 'ឈ្មោះក្រុមហ៊ុន', lgAddress: 'អាសយដ្ឋាន', lgEmail: 'អ៊ីមែលជំនួយ', lgPhone: 'ទូរស័ព្ទ', lgTelegram: 'Telegram (@username)', lgFacebook: 'ទំព័រ Facebook (https://…)', lgHours: 'ម៉ោងធ្វើការ',
    lgSaveCompany: 'រក្សាទុកព័ត៌មានទំនាក់ទំនង', lgMissing: 'នៅទទេ៖ {f}។ ឯកសារបង្ហាញ «—» រហូតដល់អ្នកបំពេញ។',
    lgDocs: 'ឯកសារ', lgCurrent: 'កំណែបច្ចុប្បន្ន {v} · តាំងពី {d}', lgView: 'មើលនៅលើគេហទំព័រ', lgHistory: 'ប្រវត្តិកំណែ',
    lgNewVersion: 'លេខកំណែថ្មី', lgVersionHint: 'ការបោះពុម្ពនីមួយៗត្រូវរក្សាទុកជារៀងរហូតជាភស្តុតាង។ ប្រើលេខថ្មី ឧ. 1.1។',
    lgPreview: 'មើលជាមុន', lgPublish: 'បោះពុម្ពកំណែថ្មី', lgPublished: 'បានបោះពុម្ពកំណែ {v}', lgConfirm: 'បោះពុម្ពកំណែ {v} នៃ «{d}» ជាបីភាសាឥឡូវនេះ?',
    lgFormat: 'ទម្រង់៖ ## ចំណងជើង, ### ចំណងជើងរង, - ចំណុច, 1. លេខរៀង, **ដិត**, [អត្ថបទ](https://…)។ Placeholder៖ {{company}} {{address}} {{email}} {{phone}} {{telegram}} {{site}} {{base}}។',
    lgChanged: 'បានកែ', lgDraftWarn: 'សូមឱ្យមេធាវីពិនិត្យគំរូទាំងនេះមុនបើកដំណើរការ — វាជាចំណុចចាប់ផ្តើម មិនមែនជាដំបូន្មានច្បាប់ទេ។',
    legal_buyer_terms: 'លក្ខខណ្ឌប្រើប្រាស់', legal_privacy_policy: 'គោលការណ៍ឯកជនភាព', legal_return_policy: 'ការប្ដូរ និងសងប្រាក់', legal_seller_agreement: 'កិច្ចព្រមព្រៀងអ្នកលក់', legal_data_deletion: 'ការលុបទិន្នន័យ'},
  zh: {lgCompanyTitle: '公司与联系方式', lgCompanyDesc: '显示在“联系我们”页面，并填入所有文件（{{company}}、{{email}} …）。上线前请填写。',
    lgCompany: '公司/法定名称', lgAddress: '地址', lgEmail: '客服邮箱', lgPhone: '电话', lgTelegram: 'Telegram（@用户名）', lgFacebook: 'Facebook 主页（https://…）', lgHours: '营业时间',
    lgSaveCompany: '保存联系方式', lgMissing: '尚未填写：{f}。填写之前文件中显示“—”。',
    lgDocs: '文件', lgCurrent: '当前版本 {v} · 自 {d} 起', lgView: '在网站上查看', lgHistory: '版本历史',
    lgNewVersion: '新版本号', lgVersionHint: '每次发布都会永久保存作为凭证。请使用新编号，例如 1.1。',
    lgPreview: '预览', lgPublish: '发布新版本', lgPublished: '已发布版本 {v}', lgConfirm: '现在以三种语言发布“{d}”的 {v} 版本？',
    lgFormat: '格式：## 标题，### 小标题，- 列表，1. 编号，**加粗**，[文字](https://…)。占位符：{{company}} {{address}} {{email}} {{phone}} {{telegram}} {{site}} {{base}}。',
    lgChanged: '已修改', lgDraftWarn: '上线前请律师审阅这些模板——它们只是起点，并非法律意见。',
    legal_buyer_terms: '使用条款', legal_privacy_policy: '隐私政策', legal_return_policy: '退货与退款', legal_seller_agreement: '卖家协议', legal_data_deletion: '数据删除'},
});

const LGL = {data: null, type: 'buyer_terms', lang: 'en', drafts: {}};
const LGL_LANG = {en: 'English', km: 'ខ្មែរ', zh: '中文'};

async function renderLegal() {
  LGL.data = await api('/admin/api/legal');
  LGL.drafts = {};
  drawLegal();
}

function nextVersion(v) {
  const m = /^(\d+)\.(\d+)$/.exec(v || '');
  return m ? `${m[1]}.${+m[2] + 1}` : '';
}

function drawLegal() {
  const d = LGL.data, s = d.settings;
  const doc = d.documents.find(x => x.type === LGL.type) || d.documents[0];
  LGL.type = doc.type;
  const cur = doc.latest.en || Object.values(doc.latest)[0] || {};
  const drafts = LGL.drafts[doc.type] || (LGL.drafts[doc.type] = Object.fromEntries(d.langs.map(l => [l, (doc.latest[l] || {}).body || ''])));
  const changed = l => drafts[l] !== ((doc.latest[l] || {}).body || '');
  const missing = [['email', 'lgEmail'], ['phone', 'lgPhone'], ['telegram', 'lgTelegram']].filter(([k]) => !s[k]).map(([, l]) => t(l));
  const field = (k, l, extra = '') => `<label>${t(l)}<input name="${k}" value="${esc(s[k] || '')}" ${extra}></label>`;
  $('#app').innerHTML = shell(`${pageHead(t('i_legal'), t('lgDraftWarn'))}
    <section class="panel" id="lglCompany"><div class="ph"><h2>${t('lgCompanyTitle')}</h2></div>
      <div style="padding:16px 18px;display:flex;flex-direction:column;gap:14px">
        <p class="muted small" style="margin:0">${esc(t('lgCompanyDesc'))}</p>
        ${missing.length ? `<div class="banner info">${esc(t('lgMissing', {f: missing.join(', ')}))}</div>` : ''}
        <div style="display:grid;grid-template-columns:repeat(auto-fit,minmax(220px,1fr));gap:14px">
          ${field('company', 'lgCompany', 'maxlength="120"')}${field('address', 'lgAddress', 'maxlength="200"')}${field('email', 'lgEmail', 'type="email" maxlength="120"')}
          ${field('phone', 'lgPhone', 'maxlength="40" inputmode="tel"')}${field('telegram', 'lgTelegram', 'maxlength="60" placeholder="@phsar_support"')}
          ${field('facebook', 'lgFacebook', 'maxlength="200" placeholder="https://facebook.com/…"')}${field('hours', 'lgHours', 'maxlength="80" placeholder="Mon–Sat 8:00–18:00"')}
        </div>
        <div><button class="primary" id="lglSaveCompany" type="button">${t('lgSaveCompany')}</button></div>
      </div></section>
    <section class="panel" id="lglDoc"><div class="ph"><h2>${t('lgDocs')}</h2></div>
      <div class="tabs" role="tablist" style="padding:14px 18px 0">${d.documents.map(x => `<button type="button" role="tab" data-doc="${x.type}" aria-selected="${x.type === doc.type}">${t('legal_' + x.type)}</button>`).join('')}</div>
      <div style="padding:16px 18px;display:flex;flex-direction:column;gap:14px">
        <div style="display:flex;flex-wrap:wrap;gap:10px;align-items:center;justify-content:space-between">
          <span class="chip ok">${esc(t('lgCurrent', {v: cur.version || '—', d: cur.effective_at ? new Date(cur.effective_at).toLocaleDateString() : '—'}))}</span>
          <a href="/en/legal/${doc.slug}" target="_blank" rel="noopener">${t('lgView')} ↗</a>
        </div>
        <div class="tabs" role="tablist" aria-label="Language">${d.langs.map(l => `<button type="button" role="tab" data-lang="${l}" aria-selected="${l === LGL.lang}">${LGL_LANG[l]}${changed(l) ? ` · <small>${t('lgChanged')}</small>` : ''}</button>`).join('')}</div>
        <textarea id="lglBody" name="body_${LGL.lang}" rows="22" spellcheck="${LGL.lang === 'en'}" lang="${LGL.lang}" style="font-family:ui-monospace,Menlo,monospace;font-size:13px;line-height:1.55;width:100%">${esc(drafts[LGL.lang])}</textarea>
        <p class="muted small" style="margin:0">${esc(t('lgFormat'))}</p>
        <div id="lglPreview" class="legal-preview" hidden></div>
        <div style="display:flex;flex-wrap:wrap;gap:10px;align-items:flex-end">
          <label>${t('lgNewVersion')}<input id="lglVersion" name="version" maxlength="16" value="${esc(nextVersion(cur.version))}" style="max-width:140px"></label>
          <button type="button" id="lglPreviewBtn">${t('lgPreview')}</button>
          <button class="primary" type="button" id="lglPublish">${t('lgPublish')}</button>
        </div>
        <p class="muted small" style="margin:0">${t('lgVersionHint')}</p>
        <details><summary>${t('lgHistory')} (${doc.versions.length})</summary>
          <div class="tbl-wrap"><table style="margin-top:8px;min-width:0;width:100%"><tbody>${doc.versions.map(v => `<tr><td><b>${esc(v.version)}</b></td><td>${new Date(v.effective_at).toLocaleString()}</td><td class="muted">${esc(v.published_by || 'Phsar (default)')}</td></tr>`).join('')}</tbody></table></div>
        </details>
      </div></section>`);

  const keep = () => { drafts[LGL.lang] = $('#lglBody').value; };
  document.querySelectorAll('[data-doc]').forEach(b => b.onclick = () => { keep(); LGL.type = b.dataset.doc; drawLegal(); });
  document.querySelectorAll('[data-lang]').forEach(b => b.onclick = () => { keep(); LGL.lang = b.dataset.lang; drawLegal(); });
  const clear = () => document.querySelectorAll('#app .invalid').forEach(el => el.classList.remove('invalid'));
  const fail = (e, form) => {
    if (e.field) {
      const lang = /^body_(\w+)$/.exec(e.field);
      if (lang && lang[1] !== LGL.lang) { keep(); LGL.lang = lang[1]; drawLegal(); }
      const el = $(`#${form} [name="${e.field}"]`); if (el) { el.classList.add('invalid'); el.focus(); }
    }
    toast(errMsg(e), 'error');
  };
  $('#lglSaveCompany').onclick = async () => {
    clear();
    const body = {};
    document.querySelectorAll('#lglCompany input[name]').forEach(i => { body[i.name] = i.value.trim(); });
    try { const keepDrafts = LGL.drafts; LGL.data = await api('/admin/api/legal/settings', {method: 'PUT', body}); LGL.drafts = keepDrafts; toast(t('lgSaved')); drawLegal(); } catch (e) { fail(e, 'lglCompany'); }
  };
  $('#lglPreviewBtn').onclick = async () => {
    keep();
    try {
      const r = await api('/admin/api/legal/preview', {method: 'POST', body: {body: drafts[LGL.lang]}});
      const p = $('#lglPreview'); p.innerHTML = r.html; p.hidden = false; p.scrollIntoView({behavior: 'smooth', block: 'nearest'});
    } catch (e) { toast(errMsg(e), 'error'); }
  };
  $('#lglPublish').onclick = async () => {
    keep(); clear();
    const version = $('#lglVersion').value.trim();
    if (!confirm(t('lgConfirm', {v: version, d: t('legal_' + doc.type)}))) return;
    try {
      LGL.data = await api(`/admin/api/legal/docs/${doc.type}`, {method: 'POST', body: {version, bodies: drafts}});
      delete LGL.drafts[doc.type];
      toast(t('lgPublished', {v: version})); drawLegal();
    } catch (e) { fail(e, 'lglDoc'); }
  };
}
