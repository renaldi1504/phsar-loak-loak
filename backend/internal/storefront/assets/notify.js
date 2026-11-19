// Phsar — notifikasi untuk pembeli di website: kabar pesanan realtime (kartu + bunyi) di halaman akun/pesanan,
// dan notifikasi browser (Web Push) supaya kabar tetap sampai saat website ditutup.
// Memakai sesi shop.js (window.PhsarShop) — tidak me-refresh token sendiri, jadi tidak bentrok dengan rotasi token.
document.addEventListener("DOMContentLoaded", function () {
  "use strict";
  var S = window.PhsarShop;
  if (!S || !window.fetch) return;
  var lang = (window.PHSAR && window.PHSAR.lang) || "en";
  var STR = {
    en: { ask: "Get order updates on this device", askSub: "We’ll tell you when your order is confirmed, shipped and arriving — even when Phsar is closed.",
      on: "Turn on", later: "Not now", done: "Order updates are on for this browser.", denied: "Notifications are blocked — allow them for this site in your browser settings.",
      ios: "On iPhone/iPad, add Phsar to your Home Screen first (Share → Add to Home Screen), then turn this on.", off: "Turn off order updates on this device", offDone: "Order updates are off for this browser.", open: "Open" },
    km: { ask: "ទទួលដំណឹងការបញ្ជាទិញលើឧបករណ៍នេះ", askSub: "យើងនឹងប្រាប់អ្នកពេលការបញ្ជាទិញត្រូវបានបញ្ជាក់ ផ្ញើចេញ និងជិតមកដល់ — ទោះបីបិទ Phsar ក៏ដោយ។",
      on: "បើក", later: "ពេលក្រោយ", done: "បានបើកដំណឹងការបញ្ជាទិញសម្រាប់កម្មវិធីរុករកនេះ។", denied: "ការជូនដំណឹងត្រូវបានរារាំង — សូមអនុញ្ញាតសម្រាប់គេហទំព័រនេះក្នុងការកំណត់កម្មវិធីរុករក។",
      ios: "លើ iPhone/iPad សូមបន្ថែម Phsar ទៅអេក្រង់ដើមជាមុន (ចែករំលែក → បន្ថែមទៅអេក្រង់ដើម) រួចបើកមុខងារនេះ។", off: "បិទដំណឹងការបញ្ជាទិញលើឧបករណ៍នេះ", offDone: "បានបិទដំណឹងការបញ្ជាទិញសម្រាប់កម្មវិធីរុករកនេះ។", open: "បើក" },
    zh: { ask: "在此设备接收订单动态", askSub: "订单确认、发货和即将送达时都会通知你——即使关闭了 Phsar。",
      on: "开启", later: "以后再说", done: "此浏览器已开启订单动态。", denied: "通知已被阻止——请在浏览器设置中允许本网站发送通知。",
      ios: "在 iPhone/iPad 上，请先将 Phsar 添加到主屏幕（分享 → 添加到主屏幕），再开启此功能。", off: "关闭此设备的订单动态", offDone: "此浏览器已关闭订单动态。", open: "打开" }
  };
  function t(k) { return (STR[lang] || STR.en)[k] || STR.en[k]; }
  function esc(v) { return String(v == null ? "" : v).replace(/[&<>"']/g, function (c) { return { "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#39;" }[c]; }); }
  var toast = window.PhsarToast || function () {};
  var BASE = "/v1/shop/notifications", SCOPE = "/";
  var mode = (document.getElementById("shop") || {}).dataset ? document.getElementById("shop").dataset.mode : "";

  // ---- bunyi (WebAudio; dibuka oleh klik pertama di halaman) ----
  var ctx = null;
  function audio() { if (!ctx) { var C = window.AudioContext || window.webkitAudioContext; if (C) ctx = new C(); } return ctx; }
  ["pointerdown", "keydown", "touchend"].forEach(function (ev) {
    document.addEventListener(ev, function () { var c = audio(); if (c && c.state !== "running") c.resume().catch(function () {}); }, { passive: true });
  });
  function tone(c, t0, f, dur, vol) {
    var o = c.createOscillator(), g = c.createGain();
    o.frequency.setValueAtTime(f, t0);
    g.gain.setValueAtTime(0.0001, t0); g.gain.exponentialRampToValueAtTime(vol, t0 + 0.012); g.gain.exponentialRampToValueAtTime(0.0001, t0 + dur);
    o.connect(g); g.connect(c.destination); o.start(t0); o.stop(t0 + dur + 0.02);
  }
  function play(name, repeat, volume) {
    var c = ctx;
    if (!c || c.state !== "running" || !name || name === "none") return;
    var v = Math.max(0, Math.min(1, (volume == null ? 80 : volume) / 100)) * 0.5, at = c.currentTime + 0.03;
    for (var i = 0; i < Math.max(1, Math.min(repeat || 1, 3)); i++) {
      if (name === "chat") { tone(c, at, 1047, 0.09, v); tone(c, at + 0.1, 1397, 0.16, v); at += 0.7; }
      else { tone(c, at, 1319, 0.5, v); tone(c, at, 2637, 0.25, v * 0.25); at += 0.9; }
    }
  }

  // ---- kartu notifikasi di halaman ----
  function card(d) {
    var box = document.getElementById("pnotes");
    if (!box) { box = document.createElement("div"); box.id = "pnotes"; box.className = "pnotes"; box.setAttribute("aria-live", "polite"); document.body.appendChild(box); }
    var el = document.createElement(d.link ? "a" : "div");
    el.className = "pnote";
    if (d.link) el.href = d.link;
    el.innerHTML = "<b>" + esc(d.title) + "</b>" + (d.body ? "<span>" + esc(d.body) + "</span>" : "") + '<button type="button" class="pnote-x" aria-label="×">×</button>';
    el.querySelector(".pnote-x").onclick = function (e) { e.preventDefault(); e.stopPropagation(); el.remove(); };
    el.addEventListener("click", function () { if (d.id) S.api("POST", BASE + "/read", { ids: [d.id] }).catch(function () {}); });
    box.insertBefore(el, box.firstChild);
    while (box.children.length > 3) box.removeChild(box.lastChild);
    setTimeout(function () { el.remove(); }, 12000);
  }
  function claim(id) { // beberapa tab → hanya satu yang berbunyi
    var key = "phsar-nfb-" + id;
    try { if (localStorage.getItem(key)) return false; localStorage.setItem(key, "1"); setTimeout(function () { try { localStorage.removeItem(key); } catch (e) {} }, 60000); } catch (e) {}
    return true;
  }

  // ---- realtime ----
  var retry = 1000, running = false, pushOn = false;
  function exp() { try { return JSON.parse(atob(S.token().split(".")[1].replace(/-/g, "+").replace(/_/g, "/"))).exp * 1000; } catch (e) { return 0; } }
  function wait(ms) { return new Promise(function (r) { setTimeout(r, ms); }); }
  function onEvent(ev, d) {
    if (ev !== "notification") return;
    var mine = claim(d.id);
    if (mine) play(d.sound, d.repeat, d.volume);
    if (document.visibilityState === "visible") card(d);
    else if (mine && !pushOn && window.Notification && Notification.permission === "granted") {
      try { var n = new Notification(d.title, { body: d.body || "", tag: "phsar-" + d.id, icon: "/ui/icon-192.png" }); n.onclick = function () { window.focus(); if (d.link) location.href = d.link; }; } catch (e) {}
    }
  }
  function read(body) {
    var rd = body.getReader(), dec = new TextDecoder(), buf = "";
    function pump() {
      return rd.read().then(function (r) {
        if (r.done) return;
        buf += dec.decode(r.value, { stream: true });
        var i;
        while ((i = buf.indexOf("\n\n")) >= 0) {
          var chunk = buf.slice(0, i), ev = "message", data = "";
          buf = buf.slice(i + 2);
          chunk.split("\n").forEach(function (l) { if (l.indexOf("event: ") === 0) ev = l.slice(7); else if (l.indexOf("data: ") === 0) data += l.slice(6); });
          if (data) { try { onEvent(ev, JSON.parse(data)); } catch (e) {} }
        }
        return pump();
      });
    }
    return pump();
  }
  function loop() {
    if (running || !S.token()) return;
    running = true;
    (function next() {
      if (!S.token()) { running = false; return; }
      var pre = exp() && exp() - Date.now() < 30000 ? S.refresh() : Promise.resolve(true);
      var started = Date.now();
      pre.then(function (ok) {
        if (!ok || !S.token()) throw "stop";
        return fetch(BASE + "/stream", { headers: { Authorization: "Bearer " + S.token(), "X-Lang": lang, Accept: "text/event-stream" }, cache: "no-store" });
      }).then(function (r) {
        if (r.status === 401) return S.refresh().then(function (ok) { if (!ok) throw "stop"; retry = 0; });
        if (!r.ok || !r.body) throw "err";
        return read(r.body);
      }).catch(function (e) { if (e === "stop") running = false; }).then(function () {
        if (!running) return;
        if (Date.now() - started > 20000) retry = 1000;
        wait(retry).then(next);
        retry = Math.min(Math.max(retry, 500) * 2, 30000);
      });
    })();
  }

  // ---- Web Push ----
  function supported() { return "serviceWorker" in navigator && "PushManager" in window && "Notification" in window && window.isSecureContext; }
  function isIOS() { return /iPad|iPhone|iPod/.test(navigator.userAgent) || (navigator.platform === "MacIntel" && navigator.maxTouchPoints > 1); }
  function keyBytes(b64) {
    var s = atob(b64.replace(/-/g, "+").replace(/_/g, "/") + "===".slice((b64.length + 3) % 4));
    var a = new Uint8Array(s.length); for (var i = 0; i < s.length; i++) a[i] = s.charCodeAt(i); return a;
  }
  function currentSub() {
    if (!supported()) return Promise.resolve(null);
    return navigator.serviceWorker.getRegistration(SCOPE).then(function (reg) { return reg ? reg.pushManager.getSubscription() : null; });
  }
  function saveSub(sub) { var j = sub.toJSON(); return S.api("POST", BASE + "/push", { endpoint: j.endpoint, keys: j.keys }); }
  function enablePush() {
    return Notification.requestPermission().then(function (perm) {
      if (perm !== "granted") throw new Error(t("denied"));
      return S.api("GET", BASE + "/prefs");
    }).then(function (p) {
      if (!p.vapid_public_key) throw new Error("push off");
      return navigator.serviceWorker.register("/sw.js", { scope: SCOPE }).then(function () { return navigator.serviceWorker.ready; }).then(function (reg) {
        return reg.pushManager.getSubscription().then(function (sub) {
          return sub || reg.pushManager.subscribe({ userVisibleOnly: true, applicationServerKey: keyBytes(p.vapid_public_key) });
        });
      });
    }).then(function (sub) { return saveSub(sub); }).then(function () { pushOn = true; });
  }
  function disablePush() {
    return currentSub().then(function (sub) {
      if (!sub) return;
      return S.api("DELETE", BASE + "/push", { endpoint: sub.endpoint }).catch(function () {}).then(function () { return sub.unsubscribe(); });
    }).then(function () { pushOn = false; });
  }
  window.PhsarNotifyOff = function () { return pushOn ? disablePush() : Promise.resolve(); };
  function dismissed() { try { return localStorage.getItem("phsar-push-ask") === "no"; } catch (e) { return false; } }
  // Ajakan mengaktifkan notifikasi browser — di halaman pesanan saja (saat kabar pesanan paling berguna).
  function offer() {
    if (mode !== "orders" && mode !== "order") return;
    var root = document.getElementById("shop");
    if (!root || document.getElementById("pushAsk")) return;
    var el = document.createElement("section");
    el.id = "pushAsk"; el.className = "panel push-ask";
    if (pushOn) {
      el.className = "push-off";
      el.innerHTML = '<button type="button" class="linkish">' + esc(t("off")) + "</button>";
      el.querySelector("button").onclick = function () { disablePush().then(function () { toast(t("offDone")); el.remove(); }); };
      root.parentNode.insertBefore(el, root.nextSibling);
      return;
    }
    if (dismissed() || (window.Notification && Notification.permission === "denied")) return;
    if (!supported()) {
      if (!isIOS()) return;
      el.innerHTML = '<span class="push-ico" aria-hidden="true">🔔</span><div><b>' + esc(t("ask")) + '</b><p class="muted small">' + esc(t("ios")) + '</p></div><button type="button" class="btn ghost" data-later>' + esc(t("later")) + "</button>";
    } else {
      el.innerHTML = '<span class="push-ico" aria-hidden="true">🔔</span><div><b>' + esc(t("ask")) + '</b><p class="muted small">' + esc(t("askSub")) + '</p></div>' +
        '<div class="push-acts"><button type="button" class="btn" data-on>' + esc(t("on")) + '</button><button type="button" class="btn ghost" data-later>' + esc(t("later")) + "</button></div>";
      el.querySelector("[data-on]").onclick = function (e) {
        e.target.disabled = true;
        enablePush().then(function () { toast(t("done")); el.remove(); }, function (err) { e.target.disabled = false; toast(err.message || String(err)); });
      };
    }
    el.querySelector("[data-later]").onclick = function () { try { localStorage.setItem("phsar-push-ask", "no"); } catch (e) {} el.remove(); };
    var outer = document.createElement("div");
    outer.className = "wrap push-wrap"; outer.appendChild(el);
    root.parentNode.insertBefore(outer, root); // di luar #shop: shop.js menggambar ulang isinya
    el.remove = function () { outer.remove(); };
  }

  // Mulai setelah shop.js punya sesi (pembeli yang belum login tidak memicu apa pun).
  var tries = 0;
  (function boot() {
    if (!S.token()) { if (++tries < 60) setTimeout(boot, 1000); return; }
    currentSub().then(function (sub) {
      pushOn = !!sub && Notification.permission === "granted";
      if (pushOn) saveSub(sub).catch(function () {}); // langganan browser ini milik user yang sedang login
    }).catch(function () {}).then(function () { setTimeout(offer, 600); });
    loop();
  })();
});
