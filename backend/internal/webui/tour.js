'use strict';
// tour: komponen tur "cara pakai" buatan sendiri (tanpa library), dipakai Seller Center dan website pembeli.
//   tourStart({steps: [{sel: '.a' | ['.a', '.b'], title, body}], labels: {next, back, skip, done, of}, onEnd(reason)})
//   - sel kosong = kartu di tengah layar (tanpa sorotan); sel tidak ditemukan / tidak terlihat = langkah dilewati.
//   - Keyboard: Esc = lewati, ← → = kembali / lanjut, Tab berputar di dalam kartu.
//   - Kartu di bawah atau di atas elemen, selalu di dalam layar (HP aman); menghormati prefers-reduced-motion.
//   - Fokus: pindah ke kartu saat mulai, kembali ke elemen semula saat selesai; sisa halaman dibuat inert.
// Semua nama global diawali "tour" (script halaman seller berbagi satu scope).

const tourState = {on: false, steps: [], i: 0, opts: null, root: null, spot: null, card: null, prevFocus: null, inerted: [], raf: 0, mo: null, el: null};

// Elemen pertama yang cocok dan benar-benar terlihat di layar.
function tourFind(sel) {
  if (!sel) return null;
  const list = Array.isArray(sel) ? sel : [sel];
  for (const s of list) {
    let nodes;
    try { nodes = document.querySelectorAll(s); } catch (e) { continue; }
    for (const el of nodes) {
      if (tourState.root && tourState.root.contains(el)) continue;
      const r = el.getBoundingClientRect(), cs = getComputedStyle(el);
      if (r.width > 0 && r.height > 0 && cs.visibility !== 'hidden' && cs.display !== 'none') return el;
    }
  }
  return null;
}

const tourReduced = () => !!(window.matchMedia && matchMedia('(prefers-reduced-motion: reduce)').matches);

function tourStart(opts) {
  if (tourState.on || !opts || !Array.isArray(opts.steps)) return false;
  // Langkah dengan target yang tidak ada dilewati sejak awal, supaya nomor "2/5" jujur.
  const steps = opts.steps.filter(s => !s.sel || tourFind(s.sel));
  if (!steps.length) return false;
  Object.assign(tourState, {on: true, steps, i: 0, opts, prevFocus: document.activeElement, el: null});

  const root = document.createElement('div');
  root.className = 'tour-root' + (tourReduced() ? ' tour-still' : '');
  root.innerHTML = `<div class="tour-block"></div><div class="tour-spot" aria-hidden="true"></div>
    <div class="tour-card" role="dialog" aria-modal="true" aria-labelledby="tourTitle" aria-describedby="tourBody" tabindex="-1">
      <div class="tour-top"><span class="tour-count" id="tourCount"></span><button type="button" class="tour-skip" data-tour="skip"></button></div>
      <h2 id="tourTitle" tabindex="-1"></h2><p id="tourBody"></p>
      <div class="tour-dots" aria-hidden="true"></div>
      <div class="tour-foot"><button type="button" class="tour-back" data-tour="back"></button><button type="button" class="tour-next" data-tour="next"></button></div>
    </div>`;
  document.body.appendChild(root);
  tourState.root = root; tourState.spot = root.querySelector('.tour-spot'); tourState.card = root.querySelector('.tour-card');
  root.addEventListener('click', e => {
    const b = e.target.closest('[data-tour]');
    if (!b) return;
    const a = b.dataset.tour;
    if (a === 'next') tourGo(tourState.i + 1); else if (a === 'back') tourGo(tourState.i - 1); else tourStop('skip');
  });
  // Sisa halaman tidak bisa difokus / dibaca pembaca layar selama tur.
  tourState.inerted = [];
  for (const c of document.body.children) {
    if (c !== root && !c.inert && c.tagName !== 'SCRIPT') { c.inert = true; tourState.inerted.push(c); }
  }
  document.addEventListener('keydown', tourKey, true);
  window.addEventListener('resize', tourQueue);
  window.addEventListener('scroll', tourQueue, true);
  // Halaman bisa menggambar ulang (mis. daftar pesanan dimuat ulang): cari ulang elemen target.
  if (window.MutationObserver) {
    tourState.mo = new MutationObserver(tourQueue);
    tourState.mo.observe(document.body, {childList: true, subtree: true});
  }
  tourGo(0);
  return true;
}

function tourGo(i) {
  if (!tourState.on) return;
  const n = tourState.steps.length;
  if (i >= n) return tourStop('done');
  if (i < 0) i = 0;
  // Target hilang setelah tur mulai: lompati ke arah yang dituju.
  const dir = i >= tourState.i ? 1 : -1;
  while (i >= 0 && i < n && tourState.steps[i].sel && !tourFind(tourState.steps[i].sel)) i += dir;
  if (i >= n) return tourStop('done');
  if (i < 0) i = tourState.i;
  tourState.i = i;
  const s = tourState.steps[i], L = tourState.opts.labels || {}, card = tourState.card;
  card.querySelector('#tourCount').textContent = (L.of || '{a}/{b}').replace('{a}', i + 1).replace('{b}', n);
  card.querySelector('#tourTitle').textContent = s.title || '';
  card.querySelector('#tourBody').textContent = s.body || '';
  card.querySelector('.tour-dots').innerHTML = tourState.steps.map((_, k) => `<i class="${k === i ? 'on' : ''}"></i>`).join('');
  const skip = card.querySelector('.tour-skip'), back = card.querySelector('.tour-back'), next = card.querySelector('.tour-next');
  skip.textContent = L.skip || 'Skip';
  skip.hidden = i === n - 1;
  back.textContent = L.back || 'Back';
  back.hidden = i === 0;
  next.textContent = i === n - 1 ? (L.done || 'Done') : (L.next || 'Next');
  tourState.el = tourFind(s.sel);
  if (tourState.el && tourState.el.scrollIntoView) {
    const r = tourState.el.getBoundingClientRect();
    if (r.top < 70 || r.bottom > innerHeight - 70) tourState.el.scrollIntoView({block: 'center', behavior: tourReduced() ? 'auto' : 'smooth'});
  }
  tourPlace();
  // Pembaca layar membaca judul langkah baru; fokus tetap di dalam kartu.
  const h = card.querySelector('#tourTitle');
  try { h.focus({preventScroll: true}); } catch (e) { h.focus(); }
}

function tourQueue() {
  if (!tourState.on || tourState.raf) return;
  tourState.raf = requestAnimationFrame(() => {
    tourState.raf = 0;
    if (!tourState.on) return;
    const s = tourState.steps[tourState.i];
    if (s.sel && (!tourState.el || !document.contains(tourState.el))) tourState.el = tourFind(s.sel);
    tourPlace();
  });
}

// Letakkan sorotan + kartu. Kartu: di bawah elemen kalau muat, kalau tidak di atas, kalau tidak di bawah layar.
function tourPlace() {
  const el = tourState.el, spot = tourState.spot, card = tourState.card;
  const vw = document.documentElement.clientWidth || innerWidth, vh = innerHeight, m = 12, gap = 12;
  const cw = Math.min(380, vw - 2 * m);
  card.style.width = cw + 'px';
  const ch = card.offsetHeight;
  let top, left;
  if (el) {
    const r = el.getBoundingClientRect(), pad = 6;
    const sx = Math.max(4, r.left - pad), sy = Math.max(4, r.top - pad);
    const sw = Math.min(vw - 4, r.right + pad) - sx, sh = Math.min(vh - 4, r.bottom + pad) - sy;
    Object.assign(spot.style, {display: 'block', left: sx + 'px', top: sy + 'px', width: Math.max(0, sw) + 'px', height: Math.max(0, sh) + 'px'});
    tourState.root.classList.remove('tour-center');
    if (sy + sh + gap + ch <= vh - m) top = sy + sh + gap;
    else if (sy - gap - ch >= m) top = sy - gap - ch;
    else top = vh - m - ch; // elemen besar: kartu menempel di bawah layar
    left = r.left + r.width / 2 - cw / 2;
  } else {
    spot.style.display = 'none';
    tourState.root.classList.add('tour-center');
    top = (vh - ch) / 2; left = (vw - cw) / 2;
  }
  card.style.left = Math.round(Math.max(m, Math.min(left, vw - cw - m))) + 'px';
  card.style.top = Math.round(Math.max(m, Math.min(top, vh - ch - m))) + 'px';
}

function tourKey(e) {
  if (!tourState.on) return;
  if (e.key === 'Escape') { e.preventDefault(); e.stopPropagation(); tourStop('skip'); return; }
  const tag = (e.target && e.target.tagName) || '';
  if (e.key === 'ArrowRight' && tag !== 'INPUT') { e.preventDefault(); tourGo(tourState.i + 1); return; }
  if (e.key === 'ArrowLeft' && tag !== 'INPUT') { e.preventDefault(); if (tourState.i > 0) tourGo(tourState.i - 1); return; }
  if (e.key === 'Tab') {
    // Fokus berputar di tombol kartu.
    const f = [...tourState.card.querySelectorAll('button')].filter(b => !b.hidden);
    if (!f.length) return;
    const k = f.indexOf(document.activeElement);
    e.preventDefault();
    const nx = e.shiftKey ? (k <= 0 ? f.length - 1 : k - 1) : (k + 1) % f.length;
    f[nx].focus();
  }
}

function tourStop(reason) {
  if (!tourState.on) return;
  tourState.on = false;
  document.removeEventListener('keydown', tourKey, true);
  window.removeEventListener('resize', tourQueue);
  window.removeEventListener('scroll', tourQueue, true);
  if (tourState.mo) tourState.mo.disconnect();
  if (tourState.raf) cancelAnimationFrame(tourState.raf);
  tourState.raf = 0; tourState.mo = null;
  for (const c of tourState.inerted) c.inert = false;
  tourState.inerted = [];
  if (tourState.root) tourState.root.remove();
  tourState.root = tourState.spot = tourState.card = tourState.el = null;
  const pf = tourState.prevFocus;
  if (pf && document.contains(pf) && pf.focus) { try { pf.focus({preventScroll: true}); } catch (e) { /* abaikan */ } }
  const cb = tourState.opts && tourState.opts.onEnd;
  tourState.opts = null;
  if (typeof cb === 'function') { try { cb(reason || 'skip'); } catch (e) { /* abaikan */ } }
}

const tourIsOn = () => tourState.on;
window.PhsarTour = {start: tourStart, stop: tourStop, isOn: tourIsOn};
