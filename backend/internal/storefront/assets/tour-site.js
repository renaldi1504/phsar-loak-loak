// tour: tur singkat untuk pengunjung baru website (video, cari, cari dengan foto, keranjang, bahasa).
// Tamu: status di localStorage. Sudah login: juga di server (/v1/app/tours) supaya tidak muncul lagi di perangkat lain.
// Mesin tur (/ui/tour.js + /ui/tour.css) baru dimuat kalau tur memang akan ditampilkan.
(function () {
  "use strict";
  var KEY = "phsar-tour-site", T = window.PHSAR_TOUR;
  function seen() { try { return localStorage.getItem(KEY) === "1"; } catch (e) { return false; } }
  function mark() { try { localStorage.setItem(KEY, "1"); } catch (e) { /* mode privat: tur bisa muncul lagi */ } }
  // Halaman keranjang/checkout/pesanan (shop.js) tidak diganggu tur.
  if (!T || seen() || window.PHSAR_SHOP || !window.fetch || !document.querySelector) return;

  var token = null;
  function json(r) { return r.ok ? r.json() : null; }
  // Sesi login website (cookie) → token akses; tamu / belum login → null.
  function session() {
    return fetch("/shop/auth/refresh", { method: "POST", headers: { "X-Phsar": "1", "Content-Type": "application/json" }, body: "{}", credentials: "same-origin" })
      .then(json).then(function (s) { token = s && s.access_token && !(s.user && s.user.is_guest) ? s.access_token : null; return token; })
      .catch(function () { return null; });
  }
  function serverDone() {
    if (!token) return Promise.resolve(false);
    return fetch("/v1/app/tours", { headers: { Authorization: "Bearer " + token } }).then(json)
      .then(function (d) { return !!(d && d.done && d.done.site); }).catch(function () { return false; });
  }
  function save() {
    mark();
    if (!token) return;
    fetch("/v1/app/tours", { method: "POST", headers: { Authorization: "Bearer " + token, "Content-Type": "application/json" }, body: JSON.stringify({ tour_key: "site" }) }).catch(function () {});
  }
  function load(cb) {
    var css = document.createElement("link");
    css.rel = "stylesheet"; css.href = "/ui/tour.css";
    document.head.appendChild(css);
    var s = document.createElement("script");
    s.src = "/ui/tour.js"; s.onload = cb;
    document.head.appendChild(s);
  }
  function busy() { return document.querySelector("dialog[open], details.lang-menu[open]"); }
  function start() {
    if (seen() || busy() || !window.PhsarTour) return;
    var steps = [
      { sel: "#videos", title: T.videoT, body: T.videoB },
      { sel: [".top-search input[type=search]", ".banner-search input[type=search]", ".big-search input[type=search]", '.tabbar a[href$="/search"]'], title: T.searchT, body: T.searchB },
      { sel: [".top-search .photo-btn", ".banner-search .photo-btn", ".big-search .photo-btn"], title: T.photoT, body: T.photoB },
      { sel: [".top-cart", '.tabbar a[href$="/cart"]'], title: T.cartT, body: T.cartB },
      { sel: ".lang-menu summary", title: T.langT, body: T.langB }
    ];
    window.PhsarTour.start({ steps: steps, labels: { next: T.next, back: T.back, skip: T.skip, done: T.done, of: T.of }, onEnd: save });
  }
  function go() {
    session().then(serverDone).then(function (done) {
      if (done) { mark(); return; }
      setTimeout(function () { load(start); }, 900);
    });
  }
  if (document.readyState === "complete") go(); else window.addEventListener("load", go);
})();
