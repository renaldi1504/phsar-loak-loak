// Phsar: keranjang, login, checkout, pembayaran KHQR, pesanan saya. Tanpa library.
// Dijalankan setelah site.js (keranjang & toast) siap: DOMContentLoaded terjadi setelah semua script defer.
document.addEventListener("DOMContentLoaded", function () {
  "use strict";
  var root = document.getElementById("shop");
  if (!root) return;
  var T = window.PHSAR_SHOP || {};
  var P = window.PHSAR || { lang: "en" };
  var Cart = window.PhsarCart;
  var toast = window.PhsarToast || function (m) { alert(m); };
  var L = "/" + P.lang;
  var mode = root.dataset.mode, closed = root.dataset.closed === "true";
  (function () { // kembali dari Facebook dengan ?auth_error=…
    var q = new URLSearchParams(location.search), e = q.get("auth_error");
    if (!e) return;
    q.delete("auth_error");
    history.replaceState(null, "", location.pathname + (q.toString() ? "?" + q : "") + location.hash);
    setTimeout(function () { toast(T["err_" + e] || e); }, 400);
  })();

  // ---------- util ----------
  function t(k, args) {
    var s = T[k] || k;
    for (var a in args || {}) s = s.split("{" + a + "}").join(args[a]);
    return s;
  }
  function glyph(n, mono) { return window.PhsarGlyph ? window.PhsarGlyph(n, mono) : ""; }
  function esc(s) {
    return String(s == null ? "" : s).replace(/[&<>"']/g, function (c) {
      return { "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#39;" }[c];
    });
  }
  function money(v) {
    var n = parseFloat(v);
    if (isNaN(n)) return "$" + v;
    return "$" + (n === Math.round(n) ? String(n) : n.toFixed(2));
  }
  function riel(v) {
    var n = Math.round(parseFloat(v) * (P.riel || 4000) / 100) * 100;
    return n > 0 ? n.toLocaleString("en-US") + " ៛" : "";
  }
  function ss(k, v) {
    try {
      if (v === undefined) return sessionStorage.getItem(k);
      if (v === null) sessionStorage.removeItem(k); else sessionStorage.setItem(k, v);
    } catch (e) { return null; }
  }
  function uuid() {
    if (window.crypto && crypto.randomUUID) return crypto.randomUUID();
    return "k" + Date.now().toString(36) + Math.random().toString(36).slice(2, 12);
  }
  function view(html) { root.innerHTML = html; window.scrollTo(0, 0); }
  function loading() { root.innerHTML = '<div class="loading" aria-busy="true"><span class="spin"></span></div>'; }
  function date(s) {
    try { return new Date(s).toLocaleString(P.lang === "zh" ? "zh-CN" : P.lang === "km" ? "km-KH" : "en-GB", { dateStyle: "medium", timeStyle: "short" }); }
    catch (e) { return s; }
  }

  // ---------- API & login ----------
  var token = null, me = null;
  function errMsg(e) {
    if (!e) return t("error");
    if (e.network) return t("network");
    return T["err_" + e.code] || e.message || t("error");
  }
  function request(method, path, body, headers) {
    var h = { "X-Lang": P.lang };
    if (body !== undefined) h["Content-Type"] = "application/json";
    for (var k in headers || {}) h[k] = headers[k];
    return fetch(path, { method: method, headers: h, body: body === undefined ? undefined : JSON.stringify(body), credentials: "same-origin" })
      .catch(function () { throw { network: true }; })
      .then(function (r) {
        return r.text().then(function (txt) {
          var data = null;
          try { data = txt ? JSON.parse(txt) : null; } catch (e) {}
          if (!r.ok) {
            var err = (data && data.error) || {};
            throw { status: r.status, code: err.code, field: err.field, message: err.message, data: data };
          }
          return data;
        });
      });
  }
  function authPost(path, body) { return request("POST", "/shop/auth" + path, body || {}, { "X-Phsar": "1" }); }
  function useSession(s) { token = s.access_token; me = s.user || me; }
  // Refresh token berotasi: dua refresh bersamaan → yang kedua gagal dan sesi hilang. Jadi satu saja yang berjalan.
  var refreshing = null;
  function refresh() {
    if (refreshing) return refreshing;
    refreshing = authPost("/refresh").then(function (s) { useSession(s); return true; }, function () { token = null; return false; });
    refreshing.then(function () { refreshing = null; });
    return refreshing;
  }
  function api(method, path, body) {
    var go = function () { return request(method, path, body, token ? { Authorization: "Bearer " + token } : {}); };
    return go().catch(function (e) {
      if (e.status !== 401) throw e;
      return refresh().then(function (ok) { if (!ok) throw e; return go(); });
    });
  }
  function ensureLogin() { return token ? Promise.resolve(true) : refresh(); }
  // notify.js (notifikasi realtime + push untuk pembeli) memakai sesi yang sama — tanpa refresh sendiri.
  window.PhsarShop = { token: function () { return token; }, refresh: refresh, ensure: ensureLogin, api: api };

  function loginView() {
    var tg = root.dataset.tg === "true", sms = root.dataset.sms !== "false", fb = root.dataset.fb === "true", gg = root.dataset.gg === "true", guest = mode === "checkout";
    return '<section class="panel login-panel"><h1>' + esc(t("loginTitle")) + '</h1><p class="muted">' + esc(t("loginSub")) + "</p>" +
      (tg ? '<button class="btn tg" id="tgBtn" type="button">' + esc(t("tg")) + '</button><p class="muted small" id="tgWait" hidden>' +
        esc(t("waitingTg")) + "</p>" : "") +
      (gg ? '<button class="btn gg" id="ggBtn" type="button"><svg viewBox="0 0 24 24" width="20" height="20" aria-hidden="true"><path fill="#4285F4" d="M23.5 12.3c0-.8-.1-1.6-.2-2.3H12v4.4h6.5a5.6 5.6 0 0 1-2.4 3.6v3h3.9c2.2-2.1 3.5-5.1 3.5-8.7z"/><path fill="#34A853" d="M12 24c3.2 0 6-1.1 8-2.9l-3.9-3c-1.1.7-2.5 1.2-4.1 1.2-3.1 0-5.8-2.1-6.7-5H1.3v3.1A12 12 0 0 0 12 24z"/><path fill="#FBBC05" d="M5.3 14.3a7.2 7.2 0 0 1 0-4.6V6.6h-4a12 12 0 0 0 0 10.8z"/><path fill="#EA4335" d="M12 4.8c1.8 0 3.3.6 4.6 1.8l3.4-3.4A12 12 0 0 0 1.3 6.6l4 3.1c.9-2.9 3.6-4.9 6.7-4.9z"/></svg>' + esc(t("gg")) + "</button>" : "") +
      (fb ? '<button class="btn fb" id="fbBtn" type="button"><svg viewBox="0 0 24 24" width="20" height="20" aria-hidden="true"><path fill="#fff" d="M24 12a12 12 0 1 0-13.9 11.9v-8.4H7.1V12h3V9.4c0-3 1.8-4.7 4.5-4.7 1.3 0 2.7.2 2.7.2v3h-1.5c-1.5 0-2 .9-2 1.9V12h3.4l-.5 3.5h-2.9v8.4A12 12 0 0 0 24 12z"/></svg>' + esc(t("fb")) + "</button>" : "") +
      ((tg || fb || gg) && sms ? '<div class="or"><span>' + esc(t("or")) + "</span></div>" : "") +
      (sms ? '<label class="field">' + esc(t("phone")) + '<input id="phone" inputmode="tel" autocomplete="tel" placeholder="012 345 678"></label>' +
      '<button class="btn ghost" id="sendBtn" type="button">' + esc(t("sendCode")) + "</button>" +
      '<label class="field">' + esc(t("code")) + '<input id="code" inputmode="numeric" autocomplete="one-time-code" maxlength="6"></label>' +
      '<button class="btn" id="verifyBtn" type="button">' + esc(t("verify")) + "</button>" : "") +
      (guest ? '<div class="or"><span>' + esc(t("or")) + '</span></div><button class="btn ghost" id="guestBtn" type="button">' + esc(t("guest")) +
        '</button><p class="muted small tc">' + esc(t("guestNote")) + "</p>" : "") + "</section>";
  }
  function bindLogin(done) {
    var tgBtn = document.getElementById("tgBtn"), poll = null;
    var fbBtn = document.getElementById("fbBtn"), guestBtn = document.getElementById("guestBtn");
    var ggBtn = document.getElementById("ggBtn");
    if (ggBtn) ggBtn.onclick = function () {
      location.href = "/shop/auth/google/start?return=" + encodeURIComponent(location.pathname + location.search);
    };
    if (fbBtn) fbBtn.onclick = function () {
      location.href = "/shop/auth/facebook/start?return=" + encodeURIComponent(location.pathname + location.search);
    };
    if (guestBtn) guestBtn.onclick = function () {
      guestBtn.disabled = true;
      authPost("/guest").then(function (s) { useSession(s); done(); }, function (e) { guestBtn.disabled = false; toast(errMsg(e)); });
    };
    if (tgBtn) tgBtn.onclick = function () {
      authPost("/telegram/start").then(function (s) {
        window.open(s.deep_link, "_blank", "noopener");
        document.getElementById("tgWait").hidden = false;
        clearInterval(poll);
        poll = setInterval(function () {
          authPost("/telegram/poll", { request_id: s.request_id, poll_token: s.poll_token }).then(function (p) {
            if (p.status === "pending") return;
            clearInterval(poll);
            if (p.status === "ok") { useSession(p.session); done(); } else toast(t("tgExpired"));
          }, function (e) { clearInterval(poll); toast(errMsg(e)); });
        }, 1500);
      }, function (e) { toast(errMsg(e)); });
    };
    if (!document.getElementById("sendBtn")) return;
    document.getElementById("sendBtn").onclick = function () {
      authPost("/otp/request", { phone: document.getElementById("phone").value }).then(function () {
        toast(t("codeSent")); document.getElementById("code").focus();
      }, function (e) { toast(errMsg(e)); });
    };
    var verify = function () {
      authPost("/otp/verify", { phone: document.getElementById("phone").value, code: document.getElementById("code").value })
        .then(function (s) { useSession(s); done(); }, function (e) { toast(errMsg(e)); });
    };
    document.getElementById("verifyBtn").onclick = verify;
    document.getElementById("code").addEventListener("keydown", function (e) { if (e.key === "Enter") verify(); });
  }
  function withLogin(fn) {
    loading();
    ensureLogin().then(function (ok) {
      if (ok) return fn();
      view(loginView());
      bindLogin(function () { loading(); fn(); });
    });
  }

  // ---------- keranjang ----------
  function lineHTML(l, editable) {
    var prob = l.problem ? '<p class="line-problem">' + esc(t("problem_" + l.problem, { n: l.available })) + "</p>" : "";
    var qty = editable
      ? '<div class="stepper" data-v="' + esc(l.variant_id) + '"><button type="button" data-step="-1" aria-label="−">−</button><span>' +
        l.requested + '</span><button type="button" data-step="1" aria-label="+"' + (l.requested >= Math.min(99, l.available) ? " disabled" : "") + ">+</button></div>"
      : '<span class="muted">× ' + l.qty + "</span>";
    return '<div class="line' + (l.problem ? " has-problem" : "") + '">' +
      '<a class="line-img" href="' + L + "/p/" + esc(l.slug) + '">' + (l.image_url ? '<img src="' + esc(l.image_url) + '" alt="">' : "") + "</a>" +
      '<div class="line-main"><a class="line-name" href="' + L + "/p/" + esc(l.slug) + '">' + esc(l.name) + "</a>" +
      (l.options ? '<p class="muted small">' + esc(l.options) + "</p>" : "") + flashHTML(l) + prob +
      '<div class="line-row">' + qty + '<strong class="line-price">' + money(l.line_total !== "0.00" ? l.line_total : l.price) + "</strong></div>" +
      (editable ? '<button type="button" class="link small" data-remove="' + esc(l.variant_id) + '">' + esc(t("remove")) + "</button>" : "") +
      "</div></div>";
  }
  // promos: voucher (kode disimpan di sessionStorage sampai dibayar) & harga flash sale
  function voucherCode() { return ss("phsar_voucher") || ""; }
  function voucherErr(code, amount) { return t("err_" + code, { a: amount || "" }); }
  function voucherHTML(q) {
    if (q.voucher) {
      return '<div class="voucher on"><span>🎟️ ' + esc(t("voucherApplied", { c: q.voucher.code })) + '</span><button type="button" class="link small" id="vRemove">' +
        esc(t("voucherRemove")) + "</button></div>";
    }
    return '<div class="voucher"><label class="sr-only" for="vCode">' + esc(t("voucher")) + '</label><input id="vCode" maxlength="24" autocomplete="off" autocapitalize="characters" placeholder="' +
      esc(t("voucherPh")) + '" value="' + esc(q.voucher_error ? q._code || "" : "") + '"><button type="button" class="btn ghost sm" id="vApply">' + esc(t("voucherApply")) + "</button></div>" +
      (q.voucher_error ? '<p class="voucher-err" role="alert">' + esc(voucherErr(q.voucher_error, q.voucher_min_spend)) + "</p>" : "");
  }
  function bindVoucher(again) {
    var ap = document.getElementById("vApply"), rm = document.getElementById("vRemove"), inp = document.getElementById("vCode");
    if (ap) {
      var go = function () { var c = inp.value.trim().toUpperCase(); if (!c) { inp.focus(); return; } ss("phsar_voucher", c); again(); };
      ap.onclick = go;
      inp.addEventListener("keydown", function (e) { if (e.key === "Enter") { e.preventDefault(); go(); } });
    }
    if (rm) rm.onclick = function () { ss("phsar_voucher", null); again(); };
  }
  // quoteReq: keranjang + voucher. Voucher yang tidak berlaku dilepas (alasannya tetap ditampilkan sekali).
  function quoteReq(items) {
    var code = voucherCode();
    return request("POST", "/v1/shop/cart/quote", { items: items, voucher_code: code }).then(function (q) {
      if (q.voucher_error) { ss("phsar_voucher", null); q._code = code; }
      return q;
    });
  }
  function keepForm() {
    var v = function (n) { var el = root.querySelector('[name="' + n + '"]'); return el ? el.value : ""; };
    if (!root.querySelector('[name="recipient"]')) return;
    state.draft = { recipient: v("recipient"), phone: v("phone"), address_line: v("address_line"), note: v("note"),
      lat: state.lat, lng: state.lng, google_place_id: state.placeId };
  }
  function flashHTML(l) {
    return l.flash_ends_at ? '<p class="line-flash">⚡ ' + esc(t("flashBadge")) + (l.compare ? ' <s class="muted">' + money(l.compare) + "</s>" : "") + "</p>" : "";
  }

  function summaryHTML(q, button, withVoucher) {
    var disc = parseFloat(q.discount || "0") > 0;
    return '<aside class="panel summary"><dl><div><dt>' + esc(t("subtotal")) + "</dt><dd>" + money(q.subtotal) + "</dd></div>" +
      // shipping: ongkir dari server (0 = Gratis); keranjang tanpa alamat = perkiraan
      "<div><dt>" + esc(t("shipping")) + (q.shipping_estimate ? ' <small class="muted">' + esc(t("shipEstimate")) + "</small>" : "") + "</dt>" +
      (parseFloat(q.shipping_fee || 0) > 0 ? "<dd>" + money(q.shipping_fee) + "</dd>" : '<dd class="free">' + esc(t("free")) + "</dd>") + "</div>" +
      (disc ? '<div class="sum-disc"><dt>' + esc(t("discount")) + "</dt><dd>−" + money(q.discount) + "</dd></div>" : "") +
      '<div class="sum-total"><dt>' + esc(t("total")) + "</dt><dd>" + money(q.total) + '<small class="riel">' + riel(q.total) + "</small></dd></div></dl>" +
      (withVoucher ? voucherHTML(q) : "") + (button || "") + "</aside>";
  }
  function syncCart(q) {
    // Cocokkan keranjang browser dengan hasil server: hapus yang tidak ada lagi.
    var known = {};
    q.stores.forEach(function (g) { g.lines.forEach(function (l) { known[l.variant_id] = l; }); });
    Cart.write(Cart.read().filter(function (x) { return known[x.variant_id]; }));
    return known;
  }
  function renderCart() {
    var items = Cart.read();
    if (!items.length) {
      view('<section class="panel center-panel"><div class="big-ico">' + glyph("cart") + '</div><h1>' + esc(t("cartTitle")) + '</h1><p class="muted">' +
        esc(t("cartEmpty")) + '</p><a class="btn" href="' + L + '/">' + esc(t("continueShopping")) + "</a></section>");
      return;
    }
    loading();
    quoteReq(items).then(function (q) {
      var known = syncCart(q);
      if (!Object.keys(known).length) { renderCart(); return; }
      var html = '<h1 class="page-title">' + esc(t("cartTitle")) + ' <small class="muted">(' + q.item_count + ")</small></h1>";
      if (q.has_problems) html += '<div class="notice">' + esc(t("cartChanged")) + ' <button type="button" class="link" id="fixCart">' + esc(t("removeProblems")) + "</button></div>";
      html += '<div class="shop-grid"><div>';
      q.stores.forEach(function (g) {
        html += '<section class="panel store-group"><a class="group-head" href="' + L + "/s/" + esc(g.store_slug) + '">' + esc(g.store_name) + " ›</a>";
        g.lines.forEach(function (l) { html += lineHTML(l, true); });
        html += "</section>";
      });
      html += "</div>" + summaryHTML(q, closed ? closedHTML() : '<a class="btn block' + (q.has_problems || !q.item_count ? " disabled" : "") + '" href="' + L + '/checkout" id="goCheckout">' + esc(t("checkout")) + "</a>" +
        payNoteHTML(), !closed && q.item_count > 0) + "</div>";
      view(html);
      bindVoucher(renderCart);
      root.querySelectorAll("[data-step]").forEach(function (b) {
        b.onclick = function () {
          var v = b.parentNode.dataset.v, d = +b.dataset.step;
          Cart.write(Cart.read().map(function (x) { if (x.variant_id === v) x.qty = Math.max(1, Math.min(99, x.qty + d)); return x; }));
          renderCart();
        };
      });
      root.querySelectorAll("[data-remove]").forEach(function (b) {
        b.onclick = function () { Cart.write(Cart.read().filter(function (x) { return x.variant_id !== b.dataset.remove; })); renderCart(); };
      });
      var fix = document.getElementById("fixCart");
      if (fix) fix.onclick = function () {
        Cart.write(Cart.read().map(function (x) { var l = known[x.variant_id]; if (l) x.qty = l.qty; return x; }).filter(function (x) { return x.qty > 0; }));
        renderCart();
      };
      var go = document.getElementById("goCheckout");
      if (go) go.onclick = function (e) { if (go.classList.contains("disabled")) e.preventDefault(); };
    }, function (e) { view('<p class="empty">' + esc(errMsg(e)) + "</p>"); });
  }

  // ---------- checkout ----------
  var state = { lat: 0, lng: 0, placeId: "", addressTyped: false, quote: null, key: null };

  function localPhone(p) { return p && p.indexOf("+855") === 0 ? "0" + p.slice(4) : (p || ""); }

  function closedHTML() {
    return '<div class="notice soon"><b>' + esc(t("closedTitle")) + "</b><br>" + esc(t("closedText")) + "</div>";
  }

  function renderCheckout() {
    if (closed) { view('<h1 class="page-title">' + esc(t("checkout")) + "</h1>" + closedHTML() + '<p><a class="btn ghost" href="' + L + '/cart">' + esc(t("cartTitle")) + "</a></p>"); return; }
    var pending = ss("phsar_checkout");
    withLogin(function () {
      if (pending) {
        return api("GET", "/v1/shop/checkouts/" + pending).then(function (v) {
          if (v.status === "pending") return renderPayment(v);
          ss("phsar_checkout", null);
          return renderCheckoutForm();
        }, function () { ss("phsar_checkout", null); renderCheckoutForm(); });
      }
      renderCheckoutForm();
    });
  }

  function renderCheckoutForm() {
    var items = Cart.read();
    if (!items.length) { location.href = L + "/cart"; return; }
    Promise.all([quoteReq(items), api("GET", "/v1/shop/address"), me ? Promise.resolve(me) : api("GET", "/v1/app/me")])
      .then(function (r) {
        var q = r[0], a = (r[1] && r[1].address) || null;
        if (state.draft) { a = state.draft; state.draft = null; } // promos: isian tetap saat voucher diganti
        me = r[2] || me;
        syncCart(q);
        if (q.has_problems || !q.item_count) { location.href = L + "/cart"; return; }
        state.quote = q;
        state.key = ss("phsar_ckey") || uuid();
        ss("phsar_ckey", state.key);
        if (a) { state.lat = a.lat; state.lng = a.lng; state.placeId = a.google_place_id || ""; state.addressTyped = true; }
        var name = a ? a.recipient : (me && me.name) || "", phone = a ? a.phone : localPhone(me && me.phone);
        var html = '<h1 class="page-title">' + esc(t("checkout")) + '</h1><div class="shop-grid"><div>' +
          '<section class="panel"><h2>' + esc(t("deliverTo")) + "</h2>" +
          '<div class="two"><label class="field">' + esc(t("recipient")) + '<input name="recipient" autocomplete="name" maxlength="80" value="' + esc(name) + '"></label>' +
          '<label class="field">' + esc(t("recipientPhone")) + '<input name="phone" inputmode="tel" autocomplete="tel" maxlength="30" value="' + esc(phone) + '"></label></div>' +
          '<div class="field" id="placeField"' + (window.Maps && Maps.google() ? "" : " hidden") + '><span>' + esc(t("searchPlace")) + '</span><div id="placeSearch" class="place-box"></div></div>' +
          '<div id="map" class="map" role="application"></div><div class="map-row"><button type="button" class="link small" id="locBtn">' + esc(t("useLocation")) +
          '</button><span class="small pin-state" id="pinState"' + (state.lat ? "" : " hidden") + ">" + esc(t("pinSet")) + "</span></div>" +
          '<p class="muted small">' + esc(t("pinHint")) + "</p>" +
          '<label class="field">' + esc(t("address")) + '<input name="address_line" maxlength="300" placeholder="' + esc(t("addressPh")) + '" value="' + esc(a ? a.address_line : "") + '"></label>' +
          '<label class="field">' + esc(t("note")) + '<input name="note" maxlength="300" placeholder="' + esc(t("notePh")) + '" value="' + esc(a ? a.note : "") + '"></label>' +
          "</section>" + payMethodHTML() + "</div><div>" +
          '<section class="panel"><h2>' + esc(t("yourOrder")) + "</h2>";
        q.stores.forEach(function (g) {
          html += '<p class="group-head small">' + esc(t("fromStore", { store: g.store_name })) + "</p>";
          g.lines.forEach(function (l) { html += lineHTML(l, false); });
        });
        html += "</section>" + summaryHTML(q, '<button type="button" class="btn block" id="payBtn">' + esc(t("payKHQR", { amount: money(q.total).slice(1) })) + "</button>" +
          payNoteHTML(), true) + "</div></div>";
        view(html);
        bindVoucher(function () { keepForm(); renderCheckoutForm(); });
        var addr = root.querySelector('[name="address_line"]');
        addr.addEventListener("input", function () { state.addressTyped = addr.value.trim() !== ""; });
        var picker = null;
        var setPin = function (p) { state.lat = p.lat; state.lng = p.lng; clearErr("address_line"); document.getElementById("pinState").hidden = false; };
        document.getElementById("locBtn").onclick = function () {
          var b = this;
          if (!navigator.geolocation) { toast(t("locationDenied")); return; }
          b.textContent = t("locating");
          navigator.geolocation.getCurrentPosition(function (pos) {
            b.textContent = t("useLocation");
            var p = { lat: +pos.coords.latitude.toFixed(6), lng: +pos.coords.longitude.toFixed(6) };
            if (picker && picker.setPin) picker.setPin(p); else setPin(p);
            setPin(p);
          }, function () { b.textContent = t("useLocation"); toast(t("locationDenied")); }, { enableHighAccuracy: true, timeout: 15000 });
        };
        if (window.Maps) {
          Maps.picker({
            mapEl: document.getElementById("map"), searchEl: document.getElementById("placeSearch"),
            pin: state.lat ? { lat: state.lat, lng: state.lng } : null,
            onPin: setPin,
            onPlace: function (p) {
              state.placeId = p.placeId || "";
              var label = [p.name, p.address].filter(Boolean).join(", ");
              if (label) { addr.value = label.slice(0, 300); state.addressTyped = true; }
            },
            onAddress: function (text, placeId) {
              if (!state.addressTyped || !addr.value) { addr.value = (text || "").slice(0, 300); }
              if (placeId) state.placeId = placeId;
            }
          }).then(function (pk) { picker = pk; }, function (e) { console.warn(e); });
        }
        document.getElementById("payBtn").onclick = pay;
      }, function (e) { view('<p class="empty">' + esc(errMsg(e)) + "</p>"); });
  }

  function fieldEl(name) { return root.querySelector('[name="' + name + '"]'); }
  function clearErr(name) {
    var el = fieldEl(name);
    if (!el) return;
    el.removeAttribute("aria-invalid");
    var m = el.parentNode.querySelector(".field-err");
    if (m) m.remove();
  }
  function markErr(name, msg) {
    var el = fieldEl(name) || fieldEl("address_line");
    el.setAttribute("aria-invalid", "true");
    var m = el.parentNode.querySelector(".field-err");
    if (!m) { m = document.createElement("span"); m.className = "field-err"; el.parentNode.appendChild(m); }
    m.textContent = msg;
    el.focus();
    el.scrollIntoView({ behavior: "smooth", block: "center" });
    el.addEventListener("input", function () { clearErr(name); }, { once: true });
  }

  function pay() {
    var btn = document.getElementById("payBtn");
    var val = function (n) { return fieldEl(n).value.trim(); };
    var req = ["recipient", "phone", "address_line"];
    for (var i = 0; i < req.length; i++) {
      if (!val(req[i])) { markErr(req[i], t("required")); return; }
    }
    if (!state.lat) { markErr("address_line", t("pinMissing")); return; }
    btn.disabled = true;
    btn.textContent = t("placing");
    api("POST", "/v1/shop/checkouts", {
      idempotency_key: state.key, items: Cart.read(), expected_total: state.quote.total, save_address: true, voucher_code: voucherCode(),
      address: { recipient: val("recipient"), phone: val("phone"), address_line: val("address_line"), note: val("note"),
        lat: state.lat, lng: state.lng, google_place_id: state.placeId }
    }).then(function (v) {
      ss("phsar_checkout", v.id);
      ss("phsar_ckey", null);
      renderPayment(v);
    }, function (e) {
      btn.disabled = false;
      btn.textContent = t("payKHQR", { amount: money(state.quote.total).slice(1) });
      if (e.code && e.code.indexOf("voucher_") === 0) { // promos: voucher tidak bisa dipakai → lepas & hitung ulang
        ss("phsar_voucher", null);
        toast(voucherErr(e.code, (e.message.match(/\$([\d.]+)/) || [])[1]));
        keepForm(); renderCheckoutForm();
        return;
      }
      if (e.code === "cart_changed") {
        ss("phsar_ckey", null);
        toast(t("cartChanged"));
        setTimeout(function () { location.href = L + "/cart"; }, 1200);
        return;
      }
      if (e.code === "validation" || e.code === "not_served") {
        ss("phsar_ckey", null);
        state.key = uuid();
        ss("phsar_ckey", state.key);
        markErr((e.field || "address.address_line").replace("address.", ""), errMsg(e));
        return;
      }
      toast(errMsg(e));
    });
  }

  // ---------- metode pembayaran (pedoman ABA PayWay: logo resmi, label "ABA KHQR", keterangan abu-abu) ----------
  function payNoteHTML() {
    return '<p class="pay-note"><img src="/assets/pay/aba-khqr.svg" alt="ABA KHQR" width="24" height="24"><span>' + esc(t("khqrNote")) + "</span></p>";
  }
  function payMethodHTML() {
    return '<section class="panel"><h2>' + esc(t("payMethod")) + '</h2><label class="pay-opt"><input type="radio" name="pay_method" value="abapay_khqr" checked>' +
      '<img src="/assets/pay/aba-khqr.svg" alt="" width="40" height="40"><span><b>ABA KHQR</b><small>' + esc(t("khqrSub")) + "</small></span></label></section>";
  }

  // Nama merchant dari isi KHQR (EMV tag 59) — harus sama dengan nama profil ABA.
  function khqrMerchant(s) {
    var i = 0;
    while (s && i + 4 <= s.length) {
      var id = s.substr(i, 2), n = parseInt(s.substr(i + 2, 2), 10);
      if (isNaN(n)) break;
      if (id === "59") return s.substr(i + 4, n);
      i += 4 + n;
    }
    return "Phsar";
  }

  // ---------- pembayaran ----------
  var timer = null, poller = null;
  function stopTimers() { clearInterval(timer); clearInterval(poller); }

  // QR digambar dari isi KHQR (qrcode-generator /ui/qrcode.js). "Download QR" memakai gambar resmi PayWay kalau ada.
  function drawKHQR(v) {
    var box = document.getElementById("khqrCode"), dl = document.getElementById("dlQR");
    if (window.qrcode && v.qr_string) {
      var q = qrcode(0, "Q");
      q.addData(v.qr_string);
      q.make();
      box.innerHTML = q.createSvgTag({ cellSize: 4, margin: 0, scalable: true });
      if (!v.qr_image) { dl.href = q.createDataURL(8, 16); dl.download = dl.download.replace(/\.png$/, ".gif"); }
    } else if (v.qr_image) {
      box.innerHTML = '<img src="' + esc(v.qr_image) + '" alt="KHQR">';
    }
    if (v.qr_image) dl.href = v.qr_image;
  }

  function renderPayment(v) {
    stopTimers();
    if (v.status === "paid") return renderPaid(v);
    if (v.status !== "pending") return renderClosed(v);
    var test = root.dataset.test === "true";
    var mobile = /Android|iPhone|iPad/i.test(navigator.userAgent);
    var amount = Number(v.total).toFixed(2);
    // Kartu KHQR sesuai standar: kepala merah + logo KHQR, nama merchant, nominal + mata uang, garis putus, QR.
    var card = '<div class="khqr-card"><div class="khqr-card-head"><img src="/assets/pay/khqr-white.svg" alt="KHQR"></div>' +
      '<div class="khqr-card-info"><span>' + esc(khqrMerchant(v.qr_string)) + "</span><b>" + esc(amount) + " <small>USD</small></b></div>" +
      '<div class="khqr-card-qr"><div id="khqrCode" role="img" aria-label="KHQR"></div><span class="khqr-cur" aria-hidden="true">$</span></div></div>';
    view('<section class="panel pay-panel"><div class="pay-body">' +
      '<h1>' + esc(t("payTitle")) + '</h1><p class="muted">' + esc(t("paySub")) + "</p>" +
      (v.deeplink && mobile ? '<a class="btn block aba-btn" href="' + esc(v.deeplink) + '"><img src="/assets/pay/aba-khqr.svg" alt="" width="28" height="28">' + esc(t("payABA")) + "</a>" +
        '<p class="or-line"><span>' + esc(t("orWord")) + "</span></p>" : "") +
      card +
      '<p class="scan-to-pay">' + esc(t("scanToPay")) + '</p><p class="muted small">' + esc(t("orWord")) + "</p>" +
      '<a class="link dl-qr" id="dlQR" href="#" download="Phsar-KHQR-' + esc(amount) + '.png">' + esc(t("downloadQR")) + "</a>" +
      '<p class="muted small">' + esc(t("uploadHint")) + "</p>" +
      '<p class="riel-line muted small">' + riel(v.total) + "</p>" +
      '<p class="countdown" id="countdown"></p><p class="waiting"><span class="spin small"></span> ' + esc(t("waitingPay")) + "</p>" +
      (test ? '<p class="notice">' + esc(t("testMode")) + '</p><button type="button" class="btn block ghost" id="simBtn">' + esc(t("simulate")) + "</button>" : "") +
      '<button type="button" class="link" id="cancelBtn">' + esc(t("cancelPay")) + "</button></div></section>");
    drawKHQR(v);
    var left = v.seconds_left, cd = document.getElementById("countdown");
    var tick = function () {
      var m = Math.floor(left / 60), s = left % 60;
      cd.textContent = t("expiresIn", { t: m + ":" + (s < 10 ? "0" : "") + s });
      if (left <= 0) { stopTimers(); check(); } else left--;
    };
    tick();
    timer = setInterval(tick, 1000);
    var check = function () {
      api("GET", "/v1/shop/checkouts/" + v.id).then(function (n) { if (n.status !== "pending") renderPayment(n); }, function () {});
    };
    poller = setInterval(check, 3000);
    document.addEventListener("visibilitychange", function () { if (!document.hidden) check(); });
    document.getElementById("cancelBtn").onclick = function () {
      if (!confirm(t("confirmCancel"))) return;
      api("POST", "/v1/shop/checkouts/" + v.id + "/cancel").then(renderPayment, function (e) { toast(errMsg(e)); });
    };
    var sim = document.getElementById("simBtn");
    if (sim) sim.onclick = function () {
      sim.disabled = true;
      api("POST", "/dev/pay/" + v.id).then(renderPayment, function (e) { sim.disabled = false; toast(errMsg(e)); });
    };
  }

  function renderPaid(v) {
    stopTimers();
    ss("phsar_checkout", null);
    ss("phsar_voucher", null); // promos
    // Barang yang sudah dibayar keluar dari keranjang.
    var bought = {};
    (v.orders || []).forEach(function (o) { (o.items || []).forEach(function (it) { bought[it.product_id] = true; }); });
    if (state.quote) {
      var ids = {};
      state.quote.stores.forEach(function (g) { g.lines.forEach(function (l) { ids[l.variant_id] = true; }); });
      Cart.write(Cart.read().filter(function (x) { return !ids[x.variant_id]; }));
    } else {
      Cart.write([]);
    }
    var list = (v.orders || []).map(function (o) {
      return '<a class="order-mini" href="' + L + "/orders/" + esc(o.order_no) + '"><strong>' + esc(t("orderNo", { no: o.order_no })) +
        "</strong><span>" + esc(o.store_name) + " · " + money(o.total) + "</span></a>";
    }).join("");
    view('<section class="panel center-panel"><div class="big-ico ok">' + glyph("check", true) + '</div><h1>' + esc(t("paidTitle")) + '</h1><p class="muted">' + esc(t("paidSub")) +
      '</p><div class="order-minis">' + list + "</div>" + (me && me.is_guest ? '<p class="notice">' + esc(t("guestPaid")) + "</p>" : "") +
      '<a class="btn" href="' + L + '/orders">' + esc(t("seeOrders")) + "</a></section>");
  }

  function renderClosed(v) {
    stopTimers();
    ss("phsar_checkout", null);
    var exp = v.status === "expired";
    view('<section class="panel center-panel"><div class="big-ico">' + glyph("box") + '</div><h1>' + esc(t(exp ? "expiredTitle" : "cancelledTitle")) + '</h1><p class="muted">' +
      esc(t(exp ? "expiredSub" : "cancelledSub")) + '</p><a class="btn" href="' + L + '/cart">' + esc(t("backToCart")) + "</a></section>");
  }

  // ---------- pesanan saya ----------
  function statusPill(s) { return '<span class="status st-' + esc(s) + '">' + esc(t("status_" + s)) + "</span>"; }

  function renderOrders() {
    withLogin(function () {
      api("GET", "/v1/shop/orders").then(function (r) {
        var who = me && (me.phone || me.is_guest) ? '<p class="muted small">' + esc(me.is_guest ? t("guestUser") : t("loggedInAs", { p: localPhone(me.phone) })) + ' · <button type="button" class="link" id="logout">' + esc(t("logout")) + "</button>" +
          ' · <a class="link" href="' + L + '/chat">' + esc(t("chatTitle")) + "</a>" + // chat:
          ' · <a class="link" href="' + L + '/support">' + esc(t("supportTitle")) + "</a></p>" : "";
        var html = '<h1 class="page-title">' + esc(t("myOrders")) + "</h1>" + who;
        if (!r.orders.length) {
          html += '<section class="panel center-panel"><p class="muted">' + esc(t("noOrders")) + '</p><a class="btn" href="' + L + '/">' + esc(t("continueShopping")) + "</a></section>";
        }
        r.orders.forEach(function (o) {
          html += '<a class="panel order-card" href="' + L + "/orders/" + esc(o.order_no) + '"><div class="order-top"><strong>' + esc(o.order_no) + "</strong>" + statusPill(o.status) +
            '</div><p class="muted small">' + esc(o.store_name) + " · " + esc(t("placedOn", { d: date(o.created_at) })) + '</p><div class="thumbs-row">' +
            o.items.slice(0, 4).map(function (it) { return it.image_url ? '<img src="' + esc(it.image_url) + '" alt="">' : '<span class="noimg-sm"></span>'; }).join("") +
            '</div><p class="order-total">' + esc(o.item_count === 1 ? t("item1") : t("items", { n: o.item_count })) + " · <strong>" + money(o.total) + "</strong></p></a>";
        });
        view(html);
        var lo = document.getElementById("logout");
        if (lo) lo.onclick = function () {
          // lepas notifikasi browser dulu (butuh token), supaya perangkat bersama tidak menerima kabar user ini lagi
          var off = window.PhsarNotifyOff ? window.PhsarNotifyOff().catch(function () {}) : Promise.resolve();
          off.then(function () { return authPost("/logout"); }).then(function () { token = null; location.href = L + "/"; });
        };
      }, function (e) { view('<p class="empty">' + esc(errMsg(e)) + "</p>"); });
    });
  }

  function renderOrder(no) {
    withLogin(function () {
      api("GET", "/v1/shop/orders/" + encodeURIComponent(no)).then(function (o) {
        var html = '<p class="crumbs"><a href="' + L + '/orders">' + esc(t("myOrders")) + "</a> › " + esc(o.order_no) + "</p>" +
          '<div class="order-top"><h1 class="page-title">' + esc(t("orderNo", { no: o.order_no })) + "</h1>" + statusPill(o.status) + "</div>" +
          '<p class="muted small">' + esc(t("placedOn", { d: date(o.created_at) })) + (o.payment_type ? " · " + esc(t("paidWith", { m: o.payment_type })) : "") + "</p>" +
          (o.refund_status === "pending" ? '<div class="notice">' + esc(t("refundPending", { amount: money(o.refund_amount || o.total).slice(1) })) + "</div>" : "") +
          (o.refund_status === "refunded" ? '<div class="notice ok">' + esc(t("refundDone", { amount: money(o.refund_amount || o.total).slice(1) })) + "</div>" : "") +
          '<div class="shop-grid"><div><div id="retBox"></div><section class="panel"><a class="group-head" href="' + L + "/s/" + esc(o.store_slug) + '">' + esc(o.store_name) + " ›</a>" +
          '<a class="btn ghost sm chat-link" href="' + L + "/chat?order=" + encodeURIComponent(o.order_no) + '">' + esc(t("chatSeller")) + "</a>"; // chat:
        o.items.forEach(function (it) {
          html += '<div class="line"><span class="line-img">' + (it.image_url ? '<img src="' + esc(it.image_url) + '" alt="">' : "") + '</span><div class="line-main"><span class="line-name">' +
            esc(it.name) + "</span>" + (it.options ? '<p class="muted small">' + esc(it.options) + "</p>" : "") +
            '<div class="line-row"><span class="muted">' + money(it.unit_price) + " × " + it.qty + '</span><strong class="line-price">' + money(it.line_total) + "</strong></div></div></div>";
        });
        html += '</section><section class="panel"><h2>' + esc(t("deliverAddress")) + "</h2><p><strong>" + esc(o.recipient) + "</strong> · " + esc(o.phone) + "</p><p>" + esc(o.address_line) + "</p>" +
          (o.address_note ? '<p class="muted">' + esc(o.address_note) + "</p>" : "") + "</section>" +
          ((o.events || []).filter(function (ev) { return T["ev_" + ev.status]; }).length ? '<section class="panel"><h2>' + esc(t("timeline")) + '</h2><ol class="timeline">' + o.events.filter(function (ev) { return T["ev_" + ev.status]; }).map(function (ev) {
            return "<li><b>" + esc(T["ev_" + ev.status] || ev.status) + '</b><span class="muted small">' + esc(date(ev.at)) + "</span></li>";
          }).join("") + "</ol></section>" : "") + "</div>" +
          summaryHTML({ subtotal: o.subtotal, shipping_fee: o.shipping_fee, total: o.total, discount: o.discount }, "") + "</div>";
        view(html);
        if (o.status === "delivered" || o.status === "completed" || o.status === "refunded") loadReturn(o.order_no);
      }, function (e) { view('<p class="empty">' + esc(errMsg(e)) + '</p><p class="tc"><a class="btn" href="' + L + '/orders">' + esc(t("myOrders")) + "</a></p>"); });
    });
  }

  // ---------- retur ----------
  var REASONS = ["damaged", "wrong_item", "not_as_described", "missing_parts", "changed_mind", "other"];
  function authedBlob(path) {
    var go = function () { return fetch(path, { headers: token ? { Authorization: "Bearer " + token } : {}, credentials: "same-origin" }); };
    return go().then(function (r) {
      if (r.status !== 401) return r;
      return refresh().then(function () { return go(); });
    }).then(function (r) { if (!r.ok) throw { status: r.status }; return r.blob(); });
  }
  function loadReturn(no) {
    var box = document.getElementById("retBox");
    if (!box) return;
    api("GET", "/v1/shop/orders/" + encodeURIComponent(no) + "/return").then(function (el) { drawReturn(box, no, el); }, function () {});
  }
  function drawReturn(box, no, el) {
    var r = el.return, html = "";
    if (r && r.status !== "cancelled") {
      var amount = money(r.refund_amount || r.requested_amount);
      var info = t("rtInfo_" + r.status, { d: date(r.status === "refunded" ? r.refunded_at : r.seller_deadline), a: amount });
      if (r.status === "rejected") info = "";
      html = '<section class="panel ret-panel ret-' + esc(r.status) + '"><div class="order-top"><h2>' + esc(t("rtNo", { no: r.return_no })) + '</h2><span class="status rt-' + esc(r.status) + '">' +
        esc(t("rtSt_" + r.status)) + "</span></div>" +
        '<p class="muted small">' + esc(t("rtReason_" + r.reason)) + " · " + esc(t("rtAmount", { a: money(r.requested_amount) })) + " · " + esc(date(r.created_at)) + "</p>" +
        "<p>" + esc(r.description) + "</p>" +
        (r.photo_count ? '<div class="ret-photos" id="retPhotos"></div>' : "") +
        (info ? '<p class="notice' + (r.status === "refunded" || r.status === "approved" ? " ok" : "") + '">' + esc(info) + "</p>" : "") +
        (r.status === "approved" && r.needs_pickup ? '<p class="muted small">' + esc(t("rtPickup")) + "</p>" : "") +
        (r.seller_note ? '<p class="ret-note"><b>' + esc(t("rtSellerSaid")) + ":</b> " + esc(r.seller_note) + "</p>" : "") +
        (r.decision_note ? '<p class="ret-note"><b>' + esc(t("rtPhsarSaid")) + ":</b> " + esc(r.decision_note) + "</p>" : "") +
        (r.status === "requested" || r.status === "disputed" ? '<button type="button" class="btn ghost sm" id="retCancel">' + esc(t("rtCancelReq")) + "</button>" : "") +
        "</section>";
    } else if (el.can_request) {
      html = '<section class="panel ret-panel"><h2>' + esc(t("rtTitle")) + '</h2><p class="muted small">' + esc(t("rtUntil", { d: date(el.defect_deadline) })) + "</p>" +
        (r && r.status === "cancelled" ? '<p class="muted small">' + esc(t("rtInfo_cancelled")) + "</p>" : "") +
        '<div class="row-btns"><button type="button" class="btn ghost" id="retOpen">' + esc(t("rtAsk")) + '</button><a class="link small" href="' + L + '/legal/returns" target="_blank">' + esc(t("rtPolicy")) + "</a></div></section>";
    }
    box.innerHTML = html;
    if (r && r.photo_count) {
      var ph = document.getElementById("retPhotos");
      for (var i = 0; i < r.photo_count; i++) (function (i) {
        authedBlob("/v1/shop/returns/" + r.id + "/photos/" + i).then(function (b) {
          var img = document.createElement("img"); img.src = URL.createObjectURL(b); img.alt = ""; ph.appendChild(img);
        }, function () {});
      })(i);
    }
    var c = document.getElementById("retCancel");
    if (c) c.onclick = function () {
      if (!confirm(t("rtConfirmCancel"))) return;
      c.disabled = true;
      api("POST", "/v1/shop/returns/" + r.id + "/cancel").then(function () { loadReturn(no); }, function (e) { c.disabled = false; toast(errMsg(e)); });
    };
    var op = document.getElementById("retOpen");
    if (op) op.onclick = function () { drawReturnForm(box, no, el); };
  }
  function drawReturnForm(box, no, el) {
    var many = el.items.length > 1 || (el.items[0] && el.items[0].qty > 1), files = [];
    var html = '<section class="panel ret-panel" id="retForm"><h2>' + esc(t("rtAsk")) + "</h2>";
    if (many) {
      html += '<fieldset class="ret-items"><legend>' + esc(t("rtWhich")) + "</legend>" + el.items.map(function (it, i) {
        return '<label class="ret-item"><input type="checkbox" data-item="' + i + '" checked><span class="line-img">' + (it.image_url ? '<img src="' + esc(it.image_url) + '" alt="">' : "") +
          '</span><span class="ret-item-name">' + esc(it.name) + (it.options ? ' <small class="muted">' + esc(it.options) + "</small>" : "") + "</span>" +
          (it.qty > 1 ? '<select data-qty="' + i + '" aria-label="Qty">' + Array.apply(null, { length: it.qty }).map(function (_, k) {
            return '<option value="' + (k + 1) + '"' + (k + 1 === it.qty ? " selected" : "") + ">" + (k + 1) + "</option>"; }).join("") + "</select>" : "") + "</label>";
      }).join("") + "</fieldset>";
    }
    html += '<fieldset class="ret-reasons"><legend>' + esc(t("rtReason")) + "</legend>" + REASONS.map(function (k) {
      var off = k === "changed_mind" && !el.change_mind_ok;
      return '<label class="opt-chip' + (off ? " off" : "") + '"><input type="radio" name="rtReason" value="' + k + '"' + (off ? " disabled" : "") + "><span>" + esc(t("rtReason_" + k)) + "</span></label>";
    }).join("") + "</fieldset>" + (el.change_mind_ok ? "" : '<p class="muted small">' + esc(t("rtNoChangeMind")) + "</p>") +
      '<label class="field">' + esc(t("rtDesc")) + '<textarea id="rtDesc" rows="4" maxlength="2000" placeholder="' + esc(t("rtDescPh")) + '"></textarea></label>' +
      '<div class="field"><span>' + esc(t("rtPhotos")) + '</span><div class="ret-photos" id="rtPrev"></div><label class="btn ghost sm file-btn">' + esc(t("rtAddPhoto")) +
      '<input type="file" id="rtFiles" accept="image/jpeg,image/png,image/webp" multiple hidden></label></div>' +
      '<div class="row-btns"><button type="button" class="btn" id="rtSend">' + esc(t("rtSend")) + '</button><button type="button" class="link" id="rtBack">' + esc(t("back")) + "</button></div></section>";
    box.innerHTML = html;
    box.scrollIntoView({ behavior: "smooth", block: "start" });
    var prev = document.getElementById("rtPrev");
    document.getElementById("rtFiles").onchange = function (ev) {
      Array.prototype.forEach.call(ev.target.files, function (f) {
        if (files.length >= 6) { toast(t("rtTooMany")); return; }
        files.push(f);
        var img = document.createElement("img"); img.src = URL.createObjectURL(f); img.alt = ""; prev.appendChild(img);
      });
      ev.target.value = "";
    };
    document.getElementById("rtBack").onclick = function () { drawReturn(box, no, el); };
    var send = document.getElementById("rtSend");
    send.onclick = function () {
      var reason = (box.querySelector('input[name="rtReason"]:checked') || {}).value, desc = document.getElementById("rtDesc").value.trim();
      if (!reason) { toast(t("rtReason")); return; }
      if (desc.length < 10) { toast(t("rtDesc")); document.getElementById("rtDesc").focus(); return; }
      var items = [];
      if (many) {
        box.querySelectorAll("[data-item]").forEach(function (cb) {
          if (!cb.checked) return;
          var i = +cb.dataset.item, q = box.querySelector('[data-qty="' + i + '"]');
          items.push({ order_item_id: el.items[i].order_item_id, qty: q ? +q.value : el.items[i].qty });
        });
        if (!items.length) { toast(t("rtWhich")); return; }
      }
      send.disabled = true;
      if (files.length) send.textContent = t("rtUploading");
      var keys = [];
      var up = files.reduce(function (p, f) {
        return p.then(function () {
          var fd = new FormData(); fd.append("kind", "return_photo"); fd.append("file", f);
          return fetch("/v1/shop/uploads", { method: "POST", body: fd, headers: { Authorization: "Bearer " + token, "X-Lang": P.lang } }).then(function (r) {
            return r.json().then(function (j) { if (!r.ok) throw (j && j.error) || {}; keys.push(j.key); });
          });
        });
      }, ensureLogin());
      up.then(function () {
        return api("POST", "/v1/shop/orders/" + encodeURIComponent(no) + "/return", { reason: reason, description: desc, photo_keys: keys, items: items });
      }).then(function () { toast(t("rtSent")); loadReturn(no); }, function (e) {
        send.disabled = false; send.textContent = t("rtSend"); toast(errMsg(e));
      });
    };
  }

  // support: halaman tiket bantuan (logika di support.js, memakai login & api dari sini)
  function renderSupport() {
    if (window.PhsarSupport) window.PhsarSupport({ api: api, withLogin: withLogin, view: view, esc: esc, t: t, date: date, toast: toast, errMsg: errMsg, L: L, root: root });
  }
  // chat: chat pembeli ↔ toko (logika di chat.js)
  function renderChat() {
    if (window.PhsarChat) window.PhsarChat({ api: api, withLogin: withLogin, view: view, esc: esc, t: t, date: date, toast: toast, errMsg: errMsg, L: L, root: root, money: money, ensureLogin: ensureLogin, token: function () { return token; } });
  }
  ({ cart: renderCart, checkout: renderCheckout, orders: renderOrders, order: function () { renderOrder(root.dataset.no); }, support: renderSupport, chat: renderChat }[mode] || renderCart)();
});
