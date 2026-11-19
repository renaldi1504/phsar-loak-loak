// Phsar: chat pembeli ↔ toko (/km/chat). Dipanggil shop.js (login, api, teks).
// /km/chat?store=…&product=…  atau  ?order=PH…  → buka (atau buat) percakapan dengan toko itu, kartu produk/pesanan
// ikut pesan pertama. /km/chat#<id> = satu percakapan. Pesan baru diambil tiap 6 detik selama halaman terlihat.
window.PhsarChat = function (k) {
  "use strict";
  var api = k.api, esc = k.esc, t = k.t, date = k.date, toast = k.toast, errMsg = k.errMsg, L = k.L, view = k.view, money = k.money;
  var timer = null, conv = null, last = 0, pending = null, photos = {};

  function stop() { if (timer) { clearInterval(timer); timer = null; } }
  function poll(fn, ms) {
    stop();
    timer = setInterval(function () { if (document.visibilityState === "visible") fn(); }, ms);
  }
  function failView(e) {
    stop();
    view('<p class="empty">' + esc(errMsg(e)) + '</p><p class="tc"><a class="btn" href="' + L + '/chat">' + esc(t("chatBack")) + "</a></p>");
  }

  function route() {
    var id = location.hash.replace(/^#/, "");
    if (/^[0-9a-f-]{36}$/.test(id)) return drawThread(id);
    var q = new URLSearchParams(location.search);
    if (q.get("store") || q.get("product") || q.get("order")) return openFrom(q);
    drawList();
  }

  // Dari halaman produk / toko / pesanan: cari atau buat percakapan, simpan kartu untuk pesan pertama.
  function openFrom(q) {
    var body = {};
    if (q.get("store")) body.store_id = q.get("store");
    if (q.get("product")) body.product_id = q.get("product");
    if (q.get("order")) body.order_no = q.get("order");
    api("POST", "/v1/shop/chat", body).then(function (c) {
      pending = c.card ? { card: c.card, product_id: body.product_id || "", order_no: body.order_no || "" } : null;
      history.replaceState(null, "", L + "/chat#" + c.id);
      drawThread(c.id);
    }, failView);
  }

  function drawList() {
    stop();
    api("GET", "/v1/shop/chat").then(function (r) {
      var html = '<h1 class="page-title">' + esc(t("chatTitle")) + "</h1>";
      if (!r.conversations.length) html += '<section class="panel center-panel"><p class="muted">' + esc(t("chatNone")) + '</p><a class="btn" href="' + L + '/">' + esc(t("continueShopping")) + "</a></section>";
      r.conversations.forEach(function (c) {
        html += '<a class="panel chat-row" href="#' + esc(c.id) + '"><span class="chat-logo">' + (c.store_logo ? '<img src="' + esc(c.store_logo) + '" alt="">' : esc((c.store_name || "?").slice(0, 1).toUpperCase())) +
          '</span><span class="chat-main"><strong>' + esc(c.store_name) + (c.unread ? ' <span class="chat-badge">' + c.unread + "</span>" : "") + "</strong>" +
          '<span class="muted small chat-prev">' + (c.last_sender === "buyer" ? esc(t("chatYou")) + ": " : "") + esc(c.preview) + "</span></span>" +
          '<span class="muted small">' + esc(date(c.last_message_at)) + "</span></a>";
      });
      view(html);
      poll(function () { api("GET", "/v1/shop/chat").then(function (r2) { if (!conv && sig(r2) !== sig(r)) drawList(); }, function () {}); }, 15000);
    }, failView);
  }
  function sig(r) { return r.conversations.map(function (c) { return c.id + ":" + c.last_message_at + ":" + c.unread; }).join("|"); }

  function cardHTML(c) {
    if (!c) return "";
    if (c.type === "order") {
      return '<a class="chat-card" href="' + L + "/orders/" + encodeURIComponent(c.order_no) + '"><span class="cc-ico" aria-hidden="true">📦</span><span><b>' +
        esc(t("chatOrder", { no: c.order_no })) + "</b><small>" + esc(t("status_" + c.status)) + " · " + money(c.total) + "</small></span></a>";
    }
    return '<a class="chat-card" href="' + L + "/p/" + esc(c.slug) + '">' + (c.image_url ? '<img src="' + esc(c.image_url) + '" alt="">' : '<span class="cc-ico" aria-hidden="true">🛍️</span>') +
      "<span><b>" + esc(c.name) + "</b><small>" + esc(t("chatProduct")) + (c.price ? " · " + money(c.price) : "") + "</small></span></a>";
  }

  function msgHTML(m) {
    var mine = m.sender === "buyer";
    return '<div class="tk-msg chat-msg ' + (mine ? "mine" : "them") + '" data-mid="' + m.id + '">' + cardHTML(m.card) +
      (m.has_photo ? '<img class="chat-photo" alt="' + esc(t("chatPhotoAlt")) + '" data-photo="/v1/shop/chat/' + esc(conv.id) + "/messages/" + m.id + '/photo">' : "") +
      (m.body ? '<div class="tk-bubble">' + esc(m.body) + "</div>" : "") +
      '<span class="muted small">' + esc(mine ? t("chatYou") : conv.store_name) + " · " + esc(date(m.created_at)) + "</span></div>";
  }

  // Foto chat privat: diambil dengan token login (bukan URL publik).
  function loadPhotos(box) {
    Array.prototype.forEach.call(box.querySelectorAll("img[data-photo]"), function (img) {
      var src = img.getAttribute("data-photo");
      img.removeAttribute("data-photo");
      if (photos[src]) { img.src = photos[src]; return; }
      k.ensureLogin().then(function () {
        return fetch(src, { headers: { Authorization: "Bearer " + k.token() }, credentials: "same-origin" });
      }).then(function (r) { return r.ok ? r.blob() : null; }).then(function (b) {
        if (!b) return;
        photos[src] = URL.createObjectURL(b);
        img.src = photos[src];
        img.onclick = function () { window.open(img.src, "_blank", "noopener"); };
      }, function () {});
    });
  }

  function drawThread(id) {
    stop();
    api("GET", "/v1/shop/chat/" + id).then(function (c) {
      c.messages = c.messages || []; // kosong → tidak dikirim (omitempty)
      conv = c;
      last = c.messages.length ? c.messages[c.messages.length - 1].id : 0;
      var html = '<p class="crumbs"><a href="' + L + '/chat" id="chatAll">' + esc(t("chatBack")) + "</a> › " + esc(c.store_name) + "</p>" +
        '<div class="order-top"><h1 class="page-title">' + esc(c.store_name) + '</h1><a class="link" href="' + L + "/s/" + esc(c.store_slug) + '">' + esc(t("chatStore")) + " ›</a></div>" +
        '<section class="panel chat-box"><div class="tk-thread chat-thread" id="chatMsgs" aria-live="polite">' +
        (c.messages.length ? c.messages.map(msgHTML).join("") : '<p class="muted chat-empty">' + esc(t("chatEmpty")) + "</p>") + "</div>" +
        (pending ? '<div class="chat-about"><span class="muted small">' + esc(t("chatAbout")) + "</span>" + cardHTML(pending.card) + "</div>" : "") +
        '<p class="muted small">' + esc(t("chatHint")) + "</p>" +
        '<label class="field"><span class="chat-sr">' + esc(t("chatPh")) + '</span><textarea id="chatBody" rows="3" maxlength="2000" placeholder="' + esc(t("chatPh")) + '"></textarea></label>' +
        '<div class="row-btns"><button type="button" class="btn sm" id="chatSend">' + esc(t("chatSend")) + "</button>" +
        '<label class="btn ghost sm chat-file">📷 ' + esc(t("chatPhoto")) + '<input type="file" id="chatFile" accept="image/jpeg,image/png,image/webp" hidden></label></div></section>';
      view(html);
      document.getElementById("chatAll").onclick = function (e) { e.preventDefault(); conv = null; pending = null; history.pushState(null, "", L + "/chat"); drawList(); };
      var box = document.getElementById("chatMsgs");
      box.scrollTop = box.scrollHeight;
      loadPhotos(box);
      var send = document.getElementById("chatSend"), ta = document.getElementById("chatBody");
      send.onclick = function () {
        var body = ta.value.trim();
        if (!body) { ta.focus(); return; }
        if (body.length > 2000) { toast(t("chatTooLong")); return; }
        send.disabled = true;
        post({ body: body }).then(function () { ta.value = ""; send.disabled = false; }, function (e) { send.disabled = false; toast(errMsg(e)); });
      };
      document.getElementById("chatFile").onchange = function (e) {
        var f = e.target.files[0];
        e.target.value = "";
        if (!f) return;
        toast(t("chatUploading"));
        var fd = new FormData();
        fd.append("kind", "chat_photo");
        fd.append("file", f);
        k.ensureLogin().then(function () {
          return fetch("/v1/shop/chat/uploads", { method: "POST", body: fd, headers: { Authorization: "Bearer " + k.token(), "X-Lang": (window.PHSAR || {}).lang || "en" } });
        }).then(function (r) {
          return r.json().then(function (j) { if (!r.ok) throw (j && j.error) || {}; return j; });
        }).then(function (up) { return post({ body: "", photo_key: up.key }); }).then(null, function (err) { toast(errMsg(err)); });
      };
      poll(pull, 6000);
    }, failView);
  }

  // Kirim: kartu produk/pesanan ikut pesan pertama yang dikirim dari halaman produk/pesanan.
  function post(body) {
    if (pending) { body.product_id = pending.product_id; body.order_no = pending.order_no; }
    return api("POST", "/v1/shop/chat/" + conv.id + "/messages", body).then(function () {
      pending = null;
      var about = document.querySelector(".chat-about");
      if (about) about.parentNode.removeChild(about);
      return pull();
    });
  }

  function pull() {
    if (!conv) return Promise.resolve();
    var id = conv.id;
    return api("GET", "/v1/shop/chat/" + id + "?after=" + last).then(function (c) {
      var box = document.getElementById("chatMsgs");
      if (!box || !conv || conv.id !== id || !c.messages || !c.messages.length) return;
      var empty = box.querySelector(".chat-empty");
      if (empty) empty.parentNode.removeChild(empty);
      var near = box.scrollHeight - box.scrollTop - box.clientHeight < 80;
      c.messages.forEach(function (m) {
        if (!box.querySelector('[data-mid="' + m.id + '"]')) box.insertAdjacentHTML("beforeend", msgHTML(m));
      });
      last = c.messages[c.messages.length - 1].id;
      if (near) box.scrollTop = box.scrollHeight;
      loadPhotos(box);
    }, function () {});
  }

  k.withLogin(function () {
    window.addEventListener("popstate", function () { conv = null; route(); }); // juga terpicu saat #hash berubah
    route();
  });
};
