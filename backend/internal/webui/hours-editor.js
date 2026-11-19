'use strict';
// Editor jam buka mingguan (dipakai Seller Center: jam buka toko, dan admin: jam standar Phsar).
// Format data: [7] Minggu..Sabtu, null = tutup | {open:"08:00", close:"22:00"}; tutup ≤ buka = lewat tengah malam.
addStrings({
  en: {hrDays: ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'], hrOpen: 'Opens', hrClose: 'Closes', hrIsOpen: 'Open',
    hrClosed: 'Closed', hrCopy: 'Use Monday’s hours for every day', hr24: '24 hours', hrOvernight: 'until the next morning', hrNeedOne: 'Open at least one day a week.'},
  km: {hrDays: ['អាទិត្យ', 'ច័ន្ទ', 'អង្គារ', 'ពុធ', 'ព្រហស្បតិ៍', 'សុក្រ', 'សៅរ៍'], hrOpen: 'បើក', hrClose: 'បិទ', hrIsOpen: 'បើក',
    hrClosed: 'បិទ', hrCopy: 'ប្រើម៉ោងថ្ងៃច័ន្ទសម្រាប់គ្រប់ថ្ងៃ', hr24: '24 ម៉ោង', hrOvernight: 'រហូតដល់ព្រឹកបន្ទាប់', hrNeedOne: 'ត្រូវបើកយ៉ាងហោចណាស់មួយថ្ងៃក្នុងមួយសប្តាហ៍។'},
  zh: {hrDays: ['星期日', '星期一', '星期二', '星期三', '星期四', '星期五', '星期六'], hrOpen: '开门', hrClose: '关门', hrIsOpen: '营业',
    hrClosed: '休息', hrCopy: '每天都使用星期一的时间', hr24: '24 小时', hrOvernight: '至次日早上', hrNeedOne: '每周至少营业一天。'},
});

// Urutan tampil: Senin dulu (lebih umum di Kamboja), Minggu terakhir.
const HR_ORDER = [1, 2, 3, 4, 5, 6, 0];

function hoursRows(week, id) {
  const days = t('hrDays');
  const row = i => {
    const d = week && week[i];
    const night = d && d.close <= d.open && !(d.open === '00:00' && d.close === '24:00');
    return `<div class="hrrow" data-hrday="${i}">
      <label class="inline hrname"><input type="checkbox" data-hr="on" ${d ? 'checked' : ''}> <b>${days[i]}</b></label>
      <span class="hrtimes ${d ? '' : 'hidden'}"><input type="time" data-hr="open" value="${esc(d ? d.open : '08:00')}" aria-label="${days[i]} ${t('hrOpen')}">
        <span aria-hidden="true">–</span><input type="time" data-hr="close" value="${esc(d ? (d.close === '24:00' ? '23:59' : d.close) : '22:00')}" aria-label="${days[i]} ${t('hrClose')}">
        ${night ? `<small class="muted">${t('hrOvernight')}</small>` : ''}</span>
      <span class="muted ${d ? 'hidden' : ''}" data-hr="closedlbl">${t('hrClosed')}</span></div>`;
  };
  return `<div class="hrgrid" id="${id}">${HR_ORDER.map(row).join('')}</div>
    <div><button type="button" class="small" data-hrcopy="${id}">${t('hrCopy')}</button></div>`;
}

function bindHoursRows(id) {
  const box = document.getElementById(id);
  if (!box) return;
  box.querySelectorAll('[data-hr=on]').forEach(cb => cb.onchange = () => {
    const r = cb.closest('.hrrow');
    r.querySelector('.hrtimes').classList.toggle('hidden', !cb.checked);
    r.querySelector('[data-hr=closedlbl]').classList.toggle('hidden', cb.checked);
  });
  const copy = document.querySelector(`[data-hrcopy="${id}"]`);
  if (copy) copy.onclick = () => {
    const mon = box.querySelector('[data-hrday="1"]');
    box.querySelectorAll('.hrrow').forEach(r => {
      if (r === mon) return;
      r.querySelector('[data-hr=on]').checked = mon.querySelector('[data-hr=on]').checked;
      r.querySelector('[data-hr=open]').value = mon.querySelector('[data-hr=open]').value;
      r.querySelector('[data-hr=close]').value = mon.querySelector('[data-hr=close]').value;
      r.querySelector('[data-hr=on]').onchange();
    });
  };
}

// readHoursRows: null kalau semua hari tutup (pemanggil menampilkan hrNeedOne).
function readHoursRows(id) {
  const w = [null, null, null, null, null, null, null];
  let any = false;
  document.querySelectorAll(`#${id} .hrrow`).forEach(r => {
    if (!r.querySelector('[data-hr=on]').checked) return;
    let open = r.querySelector('[data-hr=open]').value || '08:00', close = r.querySelector('[data-hr=close]').value || '22:00';
    if (open === '00:00' && (close === '23:59' || close === '00:00')) close = '24:00'; // 24 jam
    w[+r.dataset.hrday] = {open, close};
    any = true;
  });
  return any ? w : null;
}

// Ringkasan satu baris: "Mon–Sun 08:00–22:00" sederhana per hari.
function hoursSummary(week) {
  if (!week) return '';
  const days = t('hrDays');
  return HR_ORDER.map(i => `${days[i].slice(0, 3)} ${week[i] ? (week[i].open === '00:00' && week[i].close === '24:00' ? t('hr24') : `${week[i].open}–${week[i].close}`) : t('hrClosed')}`).join(' · ');
}
