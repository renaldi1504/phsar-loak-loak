// Phsar storefront: galeri, pilihan varian, teks asli, notifikasi kecil. Tanpa library.
(function () {
  "use strict";
  var P = window.PHSAR || { t: {} };
  var $ = function (s, el) { return (el || document).querySelector(s); };
  var $$ = function (s, el) { return Array.prototype.slice.call((el || document).querySelectorAll(s)); };

  function toast(msg) {
    var t = $("#toast");
    if (!t) return;
    t.textContent = msg;
    t.classList.add("show");
    clearTimeout(toast.timer);
    toast.timer = setTimeout(function () { t.classList.remove("show"); }, 2800);
  }

  function money(v) {
    var n = parseFloat(v);
    if (isNaN(n)) return v;
    return n === Math.round(n) ? String(n) : n.toFixed(2);
  }
  function money2(v) { var n = parseFloat(v); return isNaN(n) ? v : n.toFixed(2); }
  function riel(v) {
    var n = Math.round(parseFloat(v) * (P.riel || 4000) / 100) * 100;
    return n > 0 ? n.toLocaleString("en-US") + " ៛" : "";
  }

  // Urutkan per halaman: dropdown sort langsung kirim.
  $$("select[data-autosubmit]").forEach(function (s) {
    s.addEventListener("change", function () { s.form.submit(); });
  });

  // ---- Keranjang (disimpan di browser; server menghitung ulang harga & stok) ----
  var KEY = "phsar_cart";
  var Cart = {
    read: function () {
      try {
        var v = JSON.parse(localStorage.getItem(KEY) || "[]");
        return Array.isArray(v) ? v.filter(function (x) { return x && typeof x.variant_id === "string" && x.qty > 0; }) : [];
      } catch (e) { return []; }
    },
    write: function (items) {
      try { localStorage.setItem(KEY, JSON.stringify(items.slice(0, 50))); } catch (e) {}
      Cart.badge();
    },
    add: function (variantId, qty) {
      var items = Cart.read(), found = false;
      items.forEach(function (x) { if (x.variant_id === variantId) { x.qty = Math.min(99, x.qty + qty); found = true; } });
      if (!found) items.unshift({ variant_id: variantId, qty: qty });
      Cart.write(items);
    },
    count: function () { return Cart.read().reduce(function (n, x) { return n + x.qty; }, 0); },
    badge: function () {
      var n = Cart.count();
      $$("[data-cart-count]").forEach(function (el) { el.hidden = n === 0; el.textContent = n > 99 ? "99+" : String(n); });
    }
  };
  // Ikon Phsar (sama dengan {{glyph}} di template; defs ada di awal <body>).
  window.PhsarGlyph = function (name, mono) {
    return '<svg class="g' + (mono ? ' mono' : '') + '" viewBox="0 0 24 24" aria-hidden="true"><rect width="24" height="24" fill="' +
      (mono ? 'currentColor' : 'url(#pg)') + '" mask="url(#gm-' + name + ')"/></svg>';
  };
  window.PhsarCart = Cart;
  window.PhsarToast = toast;
  Cart.badge();
  window.addEventListener("storage", function (e) { if (e.key === KEY) Cart.badge(); });

  // Tombol keranjang di kartu (produk dengan satu varian).
  document.addEventListener("click", function (e) {
    var b = e.target.closest("[data-quick-add]");
    if (!b) return;
    e.preventDefault();
    Cart.add(b.getAttribute("data-quick-add"), 1);
    b.classList.add("done");
    setTimeout(function () { b.classList.remove("done"); }, 900);
    toast(P.t.added);
  });

  // Menu bahasa: tutup kalau klik di luar / Esc.
  var lm = $(".lang-menu");
  if (lm) {
    document.addEventListener("click", function (e) { if (lm.open && !lm.contains(e.target)) lm.open = false; });
    document.addEventListener("keydown", function (e) { if (e.key === "Escape") lm.open = false; });
  }

  // Menu kategori desktop: tutup kalau klik di luar / Esc.
  var cm = $(".cat-menu");
  if (cm) {
    document.addEventListener("click", function (e) { if (cm.open && !cm.contains(e.target)) cm.open = false; });
    document.addEventListener("keydown", function (e) { if (e.key === "Escape" && cm.open) { cm.open = false; cm.querySelector("summary").focus(); } });
    catMenu(cm);
  }

  // Tonton: satu video berputar (yang terlihat), geser ke atas untuk berikutnya; tombol suara.
  var feed = $("#watchFeed");
  if (feed && "IntersectionObserver" in window) {
    var soundOn = false;
    var wio = new IntersectionObserver(function (ents) {
      ents.forEach(function (en) {
        var v = en.target;
        if (en.isIntersecting) {
          v.muted = !soundOn; v.preload = "auto";
          var pr = v.play(); if (pr && pr.catch) pr.catch(function () {});
          var nx = v.closest(".watch-item").nextElementSibling;
          if (nx) { var nv = nx.querySelector("video"); if (nv) nv.preload = "auto"; }
        } else { v.pause(); }
      });
    }, { threshold: 0.7 });
    $$("video[data-watch]", feed).forEach(function (v) { wio.observe(v); });
    feed.addEventListener("click", function (e) {
      var b = e.target.closest("[data-sound]");
      if (b) {
        soundOn = !soundOn;
        $$("video[data-watch]", feed).forEach(function (v) { v.muted = !soundOn; });
        $$("[data-sound]", feed).forEach(function (x) { x.classList.toggle("on", soundOn); });
        return;
      }
      var v = e.target.closest("video");
      if (v) { if (v.paused) v.play(); else v.pause(); }
    });
  }
  // Panah geser untuk rail video.
  $$("[data-rail-prev],[data-rail-next]").forEach(function (b) {
    b.addEventListener("click", function () {
      var r = document.getElementById(b.getAttribute("data-rail-prev") || b.getAttribute("data-rail-next"));
      if (!r) return;
      var dir = b.hasAttribute("data-rail-prev") ? -1 : 1;
      r.scrollBy({ left: dir * r.clientWidth * 0.8, behavior: "smooth" });
    });
  });

  // Banner: titik penanda + geser otomatis.
  var track = $("#bannerTrack"), dotsEl = $("#bannerDots");
  if (track && dotsEl) {
    var slides = $$(".slide", track);
    slides.forEach(function () { dotsEl.appendChild(document.createElement("i")); });
    var dots = $$("i", dotsEl);
    var cur = 0;
    var mark = function () {
      cur = Math.round(track.scrollLeft / track.clientWidth);
      dots.forEach(function (d, i) { d.classList.toggle("on", i === cur); });
    };
    track.addEventListener("scroll", function () { clearTimeout(mark.t); mark.t = setTimeout(mark, 60); }, { passive: true });
    mark();
    var reduce = window.matchMedia && matchMedia("(prefers-reduced-motion: reduce)").matches;
    var hold = false;
    ["pointerdown", "focusin", "mouseenter"].forEach(function (ev) { track.addEventListener(ev, function () { hold = true; }); });
    ["mouseleave", "focusout"].forEach(function (ev) { track.addEventListener(ev, function () { hold = false; }); });
    if (!reduce && slides.length > 1) {
      setInterval(function () {
        if (hold || document.hidden) return;
        var next = (cur + 1) % slides.length;
        track.scrollTo({ left: next * track.clientWidth, behavior: "smooth" });
      }, 5000);
    }
  }

  // ---- Cari dengan foto ----
  var sheet = $("#photoSheet"), pin = $("#photoInput");
  if (sheet && pin) {
    var status = $("#photoStatus"), results = $("#photoResults"), prevImg = $("#photoPreview"), pick = $(".photo-pick", sheet);
    var label = $(".photo-choose span", sheet);
    var open = function () { if (!sheet.open) { if (sheet.showModal) sheet.showModal(); else sheet.setAttribute("open", ""); } };
    $$("[data-photo]").forEach(function (b) {
      b.addEventListener("click", function () { open(); pin.click(); });
    });
    $$("[data-close]", sheet).forEach(function (b) { b.addEventListener("click", function () { sheet.close(); }); });
    sheet.addEventListener("click", function (e) { if (e.target === sheet) sheet.close(); });
    var seq = 0, lastURL = null;
    pin.addEventListener("change", function () {
      var f = pin.files && pin.files[0];
      pin.value = "";
      if (!f) return;
      open();
      results.textContent = "";
      if (!/^image\//.test(f.type)) { status.textContent = P.t.photoBad; return; }
      if (f.size > 10 * 1024 * 1024) { status.textContent = P.t.photoTooBig; return; }
      if (lastURL) URL.revokeObjectURL(lastURL);
      lastURL = URL.createObjectURL(f);
      prevImg.src = lastURL; pick.hidden = false;
      $("[data-photo-intro]", sheet).hidden = true;
      label.textContent = label.getAttribute("data-label-b");
      status.innerHTML = '<span class="spin small"></span> ' + P.t.photoWorking;
      var my = ++seq;
      var fd = new FormData();
      fd.append("image", f, f.name || "photo.jpg");
      fetch("/" + P.lang + "/photo-search", { method: "POST", body: fd, headers: { "Accept": "text/html" } })
        .then(function (r) {
          return r.text().then(function (html) { return { ok: r.ok, status: r.status, html: html, n: +(r.headers.get("X-Results") || 0) }; });
        })
        .then(function (r) {
          if (my !== seq) return;
          if (r.status === 503) { status.textContent = P.t.photoOff; return; }
          if (r.status === 429) { status.textContent = P.t.photoBusy; return; }
          if (r.status === 413) { status.textContent = P.t.photoTooBig; return; }
          if (r.status === 415 || r.status === 422) { status.textContent = P.t.photoBad; return; }
          if (!r.ok) { status.textContent = P.t.photoOff; return; }
          status.textContent = r.n ? P.t.photoResults : P.t.photoNone;
          results.innerHTML = r.html;
        })
        .catch(function () { if (my === seq) status.textContent = P.t.photoOff; });
    });
  }

  // ---- Galeri ----
  var stage = $("#stage");
  $$(".thumb").forEach(function (b) {
    b.addEventListener("click", function () {
      $$(".thumb").forEach(function (x) { x.classList.toggle("on", x === b); });
      var old = $("video", stage);
      if (old) old.pause();
      stage.textContent = "";
      var el;
      if (b.dataset.kind === "video") {
        el = document.createElement("video");
        el.controls = true;
        el.playsInline = true;
        el.preload = "metadata";
        if (b.dataset.poster) el.poster = b.dataset.poster;
        el.src = b.dataset.src;
      } else {
        el = document.createElement("img");
        el.alt = "";
        el.src = b.dataset.src;
      }
      stage.appendChild(el);
      if (el.tagName === "VIDEO") el.play().catch(function () {});
    });
  });

  // ---- Teks asli / terjemahan ----
  var tog = $("#toggle-orig");
  if (tog) {
    tog.addEventListener("click", function () {
      var tr = $("#desc-tr"), og = $("#desc-orig");
      var showOrig = og.hidden;
      og.hidden = !showOrig;
      tr.hidden = showOrig;
      tog.textContent = showOrig ? tog.dataset.b : tog.dataset.a;
    });
  }

  // ---- Varian & tombol beli ----
  var vj = $("#variants");
  if (!vj) return;
  var variants = [];
  try { variants = JSON.parse(vj.textContent) || []; } catch (e) { return; }
  var box = $("#options");
  var chosen = function () { return variants.length === 1 && !box ? variants[0] : null; };

  function buy(goCheckout) {
    var v = chosen();
    if (!v) {
      var miss = box ? $$(".opt", box).filter(function (o, k) { return sel[k] === null; })[0] : null;
      var name = miss ? $("legend", miss).firstChild.textContent.replace(/:\s*$/, "") : "";
      toast((P.t.chooseFirst || "").replace("{opt}", name));
      if (box) box.scrollIntoView({ behavior: "smooth", block: "center" });
      return;
    }
    if (v.s <= 0) { toast(P.t.out); return; }
    Cart.add(v.id, 1);
    if (goCheckout) { location.href = "/" + P.lang + "/checkout"; return; }
    toast(P.t.added);
  }
  $$("[data-add]").forEach(function (b) { b.addEventListener("click", function () { buy(false); }); });
  $$("[data-buy]").forEach(function (b) { b.addEventListener("click", function () { buy(true); }); });

  if (!box) return;
  var opts = $$(".opt", box);
  var n = opts.length;
  var sel = [];
  for (var i = 0; i < n; i++) sel.push(null);

  function matches(v, choice) {
    for (var k = 0; k < n; k++) {
      if (choice[k] !== null && (v.i || [])[k] !== choice[k]) return false;
    }
    return true;
  }
  function available(choice) {
    return variants.some(function (v) { return v.s > 0 && matches(v, choice); });
  }

  // Opsi dengan satu nilai langsung dipilih.
  opts.forEach(function (o, k) {
    var vals = $$(".val", o);
    if (vals.length === 1) sel[k] = 0;
  });

  var price = $("#price"), compare = $("#compare"), off = $("#off"), rielEl = $("#riel"), stock = $("#stock"), bb = $("#bb-price");

  function render() {
    $$(".val", box).forEach(function (b) {
      var k = +b.dataset.opt, val = +b.dataset.val;
      var choice = sel.slice();
      choice[k] = val;
      b.classList.toggle("na", !available(choice));
      b.setAttribute("aria-pressed", sel[k] === val ? "true" : "false");
    });
    $$(".opt-val", box).forEach(function (el) {
      var k = +el.dataset.for;
      var b = sel[k] === null ? null : $('.val[data-opt="' + k + '"][data-val="' + sel[k] + '"]', box);
      el.textContent = b ? b.title : "";
    });
    var missing = sel.indexOf(null);
    chosen = function () { return null; };
    stock.className = "stock";
    if (missing >= 0) {
      var name = $("legend", opts[missing]).firstChild.textContent.replace(/:\s*$/, "");
      stock.textContent = P.t.choose + name;
      stock.classList.add("low");
      return;
    }
    var v = variants.filter(function (x) { return matches(x, sel); })[0];
    chosen = function () { return v || null; };
    if (!v) { stock.textContent = P.t.out; stock.classList.add("out"); return; }
    var txt = "$" + money2(v.p);
    price.textContent = txt;
    if (bb) bb.textContent = txt;
    if (rielEl) rielEl.textContent = riel(v.p);
    var bbr = $("#bb-riel");
    if (bbr) bbr.textContent = riel(v.p);
    if (compare) {
      compare.hidden = !v.c;
      compare.textContent = v.c ? "$" + money(v.c) : "";
    }
    if (off) {
      var pc = v.c ? Math.round((parseFloat(v.c) - parseFloat(v.p)) / parseFloat(v.c) * 100) : 0;
      off.hidden = !(pc > 0);
      off.textContent = "-" + pc + "%";
    }
    if (v.s <= 0) { stock.textContent = P.t.out; stock.classList.add("out"); }
    else if (v.s <= 5) { stock.textContent = P.t.left.replace("{n}", v.s); stock.classList.add("low"); }
    else stock.textContent = P.t.inStock;
  }

  box.addEventListener("click", function (e) {
    var b = e.target.closest(".val");
    if (!b) return;
    var k = +b.dataset.opt, val = +b.dataset.val;
    sel[k] = sel[k] === val && $$(".val", opts[k]).length > 1 ? null : val;
    render();
  });
  render();
})();

// promos: hitung mundur flash sale (kartu & halaman produk). Setelah habis, muat ulang supaya harga normal tampil.
(function () {
  "use strict";
  var els = document.querySelectorAll("[data-flash-ends]");
  if (!els.length) return;
  function left(ms) {
    var s = Math.max(0, Math.floor(ms / 1000)), d = Math.floor(s / 86400), h = Math.floor(s % 86400 / 3600), m = Math.floor(s % 3600 / 60);
    var pad = function (n) { return n < 10 ? "0" + n : "" + n; };
    return (d ? d + "d " : "") + pad(h) + ":" + pad(m) + ":" + pad(s % 60);
  }
  var reloaded = false;
  function tick() {
    var now = Date.now();
    els.forEach(function (el) {
      var ms = new Date(el.getAttribute("data-flash-ends")) - now, out = el.querySelector("[data-flash-left]");
      if (ms <= 0) {
        if (out && !reloaded) { reloaded = true; setTimeout(function () { location.reload(); }, 1500); }
        el.hidden = true;
        return;
      }
      if (out) out.textContent = (out.getAttribute("data-tpl") || "{t}").replace("{t}", left(ms));
    });
  }
  tick();
  setInterval(tick, 1000);
})();

// Infinite scroll (kategori, cari, toko, produk terbaru): halaman berikutnya dimuat otomatis saat link
// "Muat lagi" mendekati layar. Setelah AUTO halaman, tombol harus diklik supaya footer tetap terjangkau.
// Tanpa JS / gagal muat: link biasa ke ?page=N tetap berfungsi.
(function () {
  "use strict";
  var AUTO = 4;
  var links = document.querySelectorAll("a[data-more]");
  if (!links.length || !window.fetch) return;
  links.forEach(function (a) {
    var box = a.parentNode, prev = box.previousElementSibling, grid = prev && prev.classList.contains("grid") ? prev : box.parentNode.querySelector(".grid");
    if (!grid) return;
    var loads = 0, busy = false, io = null, label = a.textContent;
    grid.setAttribute("aria-live", "polite");
    function load() {
      if (busy || !a.getAttribute("href")) return;
      busy = true;
      a.setAttribute("aria-busy", "true");
      a.innerHTML = '<span class="spin small"></span> ' + label;
      var u = new URL(a.getAttribute("href"), location.href);
      u.searchParams.set("fragment", "1");
      fetch(u.toString(), { headers: { Accept: "text/html" }, credentials: "same-origin" })
        .then(function (r) { if (!r.ok) throw new Error(r.status); return r.text(); })
        .then(function (html) {
          var tpl = document.createElement("template");
          tpl.innerHTML = html;
          var next = tpl.content.querySelector("a[data-next]");
          if (next) next.remove();
          var first = tpl.content.firstElementChild;
          grid.appendChild(tpl.content);
          loads++;
          if (next) {
            a.setAttribute("href", next.getAttribute("href"));
            a.textContent = label;
            if (loads >= AUTO && io) { io.disconnect(); io = null; } // selanjutnya: klik manual
            else if (io && box.getBoundingClientRect().top < innerHeight + 600) setTimeout(load, 50); // masih dekat layar
          } else {
            if (io) io.disconnect();
            box.remove();
          }
          if (first && document.activeElement === a && first.querySelector("a")) first.querySelector("a").focus({ preventScroll: true });
        })
        .catch(function () { a.textContent = label; if (io) { io.disconnect(); io = null; } })
        .then(function () { busy = false; a.removeAttribute("aria-busy"); });
    }
    a.addEventListener("click", function (e) { e.preventDefault(); load(); });
    if ("IntersectionObserver" in window) {
      io = new IntersectionObserver(function (es) { if (es[0].isIntersecting) load(); }, { rootMargin: "600px 0px" });
      io.observe(box);
    }
  });
})();

// Baris chip kategori bisa digeser: tampilkan chip yang aktif (kategori utama & sub-kategori) tanpa menggeser halaman.
(function () {
  document.querySelectorAll('.chips .chip.on, .watch-chips a.on').forEach(function (c) {
    var row = c.parentElement;
    if (row.scrollWidth <= row.clientWidth) return;
    row.scrollLeft = c.offsetLeft - row.offsetLeft - (row.clientWidth - c.offsetWidth) / 2;
  });
})();

// Menu kategori desktop:
//  - mouse: terbuka saat kursor di atas tombol, tertutup sendiri saat kursor keluar (dengan jeda supaya tidak berkedip);
//  - kategori aktif disimpan (li.on) dan TIDAK pindah saat kursor melewati celah/baris lain menuju panel sub-kategori:
//    selama kursor bergerak ke arah panel (segitiga dari posisi kursor ke tepi kiri panel), pergantian ditunda ("menu aim").
function catMenu(cm) {
  "use strict";
  var list = cm.querySelector(".cat-list"), summary = cm.querySelector("summary");
  if (!list) return;
  var items = Array.prototype.slice.call(list.children);
  cm.classList.add("js");
  function activate(li) {
    if (!li || li.classList.contains("on")) return;
    items.forEach(function (x) { x.classList.toggle("on", x === li); });
  }
  function reset() { activate(items[0]); }
  reset();
  cm.addEventListener("toggle", function () { if (cm.open) reset(); });

  // --- menu aim ---
  var track = [], timer = null, pending = null;
  document.addEventListener("mousemove", function (e) {
    if (!cm.open) return;
    track.push({ x: e.clientX, y: e.clientY });
    if (track.length > 4) track.shift();
  });
  function sub() { var on = list.querySelector("li.on .cat-sub"); return on && on.getBoundingClientRect(); }
  // Masih menuju panel? (titik sekarang di dalam segitiga: titik lama → sudut kiri-atas & kiri-bawah panel)
  function aiming() {
    var r = sub();
    if (!r || track.length < 2) return false;
    var a = track[0], b = track[track.length - 1];
    if (b.x <= a.x) return false; // tidak bergerak ke kanan
    var slope = function (p, q) { return (q.y - p.y) / (q.x - p.x || 0.0001); };
    var up = { x: r.left, y: r.top - 40 }, down = { x: r.left, y: r.bottom + 40 };
    return slope(b, up) < slope(a, up) && slope(b, down) > slope(a, down);
  }
  function attempt() {
    timer = null;
    if (!pending) return;
    if (aiming()) { timer = setTimeout(attempt, 120); return; }
    activate(pending);
    pending = null;
  }
  items.forEach(function (li) {
    li.addEventListener("mouseenter", function () {
      pending = li;
      if (timer) clearTimeout(timer);
      if (aiming()) timer = setTimeout(attempt, 120); else attempt();
    });
    li.addEventListener("focusin", function () { activate(li); });
  });
  list.addEventListener("mouseleave", function () { pending = null; if (timer) { clearTimeout(timer); timer = null; } });

  // --- buka dengan hover (hanya perangkat dengan mouse) ---
  if (!(window.matchMedia && matchMedia("(hover: hover) and (pointer: fine)").matches)) return;
  var openT = null, closeT = null, hoverOpenedAt = 0;
  cm.addEventListener("mouseenter", function () {
    if (closeT) { clearTimeout(closeT); closeT = null; }
    if (!cm.open && !openT) openT = setTimeout(function () { openT = null; cm.open = true; hoverOpenedAt = Date.now(); }, 120);
  });
  cm.addEventListener("mouseleave", function () {
    if (openT) { clearTimeout(openT); openT = null; }
    if (cm.open && !closeT) closeT = setTimeout(function () { closeT = null; if (!cm.contains(document.activeElement) || document.activeElement === summary) cm.open = false; }, 350);
  });
  // Klik tepat setelah terbuka karena hover jangan langsung menutupnya lagi.
  summary.addEventListener("click", function (e) {
    if (openT) { clearTimeout(openT); openT = null; }
    if (cm.open && Date.now() - hoverOpenedAt < 600) e.preventDefault();
  });
}
