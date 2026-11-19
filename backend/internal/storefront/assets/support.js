// Phsar: tiket bantuan pembeli (/km/support). Dipanggil shop.js (login, api, teks).
window.PhsarSupport = function (k) {
  "use strict";
  var api = k.api, esc = k.esc, t = k.t, date = k.date, toast = k.toast, errMsg = k.errMsg, L = k.L, view = k.view;
  var CATS = ["order", "payment", "return", "account", "other"];
  function pill(s) { return '<span class="status tk-' + esc(s) + '">' + esc(t("tkSt_" + s)) + "</span>"; }
  function route() {
    var id = location.hash.replace(/^#/, "");
    if (id === "new") return drawNew();
    if (/^[0-9a-f-]{36}$/.test(id)) return drawThread(id);
    drawList();
  }
  function drawList() {
    api("GET", "/v1/shop/tickets").then(function (r) {
      var html = '<div class="order-top"><h1 class="page-title">' + esc(t("supportTitle")) + '</h1><a class="btn sm" href="#new">' + esc(t("tkNew")) + "</a></div>" +
        '<p class="muted">' + esc(t("tkSub")) + "</p>";
      if (!r.tickets.length) html += '<section class="panel center-panel"><p class="muted">' + esc(t("tkNone")) + '</p><a class="btn" href="#new">' + esc(t("tkNew")) + "</a></section>";
      r.tickets.forEach(function (x) {
        html += '<a class="panel order-card tk-card" href="#' + esc(x.id) + '"><div class="order-top"><strong>' + esc(x.subject) + "</strong>" + pill(x.status) + "</div>" +
          '<p class="muted small">' + esc(x.ticket_no) + " · " + esc(t("tkCat_" + x.category)) + (x.order_no ? " · " + esc(x.order_no) : "") + " · " + esc(date(x.updated_at)) + "</p>" +
          '<p class="tk-preview">' + esc(x.preview) + "</p></a>";
      });
      view(html);
    }, function (e) { view('<p class="empty">' + esc(errMsg(e)) + "</p>"); });
  }
  function drawNew() {
    api("GET", "/v1/shop/orders").catch(function () { return { orders: [] }; }).then(function (r) {
      var pre = new URLSearchParams(location.search).get("order") || "";
      var html = '<p class="crumbs"><a href="#">' + esc(t("tkBack")) + "</a></p><h1 class=\"page-title\">" + esc(t("tkNew")) + '</h1><section class="panel" id="tkForm">' +
        '<fieldset class="ret-reasons tk-cats"><legend>' + esc(t("tkTopic")) + "</legend>" + CATS.map(function (c, i) {
          return '<label class="opt-chip"><input type="radio" name="tkCat" value="' + c + '"' + (i === 0 ? " checked" : "") + "><span>" + esc(t("tkCat_" + c)) + "</span></label>";
        }).join("") + "</fieldset>" +
        '<label class="field">' + esc(t("tkOrder")) + '<select id="tkOrder"><option value="">' + esc(t("tkNoOrder")) + "</option>" + (r.orders || []).map(function (o) {
          return '<option value="' + esc(o.order_no) + '"' + (o.order_no === pre ? " selected" : "") + ">" + esc(o.order_no + " · " + o.store_name) + "</option>";
        }).join("") + "</select></label>" +
        '<label class="field">' + esc(t("tkSubject")) + '<input id="tkSubject" maxlength="120" placeholder="' + esc(t("tkSubjectPh")) + '"></label>' +
        '<label class="field">' + esc(t("tkMessage")) + '<textarea id="tkMsg" rows="6" maxlength="4000" placeholder="' + esc(t("tkMessagePh")) + '"></textarea></label>' +
        '<button type="button" class="btn" id="tkSend">' + esc(t("tkSend")) + "</button></section>";
      view(html);
      var send = document.getElementById("tkSend");
      send.onclick = function () {
        var body = { category: (document.querySelector('input[name="tkCat"]:checked') || {}).value, order_no: document.getElementById("tkOrder").value,
          subject: document.getElementById("tkSubject").value, message: document.getElementById("tkMsg").value };
        send.disabled = true;
        api("POST", "/v1/shop/tickets", body).then(function (x) { toast(t("tkSent")); location.hash = x.id; }, function (e) {
          send.disabled = false;
          var f = { subject: "tkSubject", message: "tkMsg", order_no: "tkOrder" }[e.field];
          if (f) { var el = document.getElementById(f); el.classList.add("invalid"); el.focus(); }
          toast(errMsg(e));
        });
      };
    });
  }
  function drawThread(id) {
    api("GET", "/v1/shop/tickets/" + id).then(function (x) {
      var html = '<p class="crumbs"><a href="#">' + esc(t("tkBack")) + "</a> › " + esc(x.ticket_no) + "</p>" +
        '<div class="order-top"><h1 class="page-title">' + esc(x.subject) + "</h1>" + pill(x.status) + "</div>" +
        '<p class="muted small">' + esc(t("tkCat_" + x.category)) + (x.order_no ? ' · <a href="' + L + "/orders/" + esc(x.order_no) + '">' + esc(x.order_no) + "</a>" : "") + "</p>" +
        '<section class="panel tk-thread">' + x.messages.map(function (m) {
          var mine = m.author_type === "user";
          return '<div class="tk-msg ' + (mine ? "mine" : "them") + '"><div class="tk-bubble">' + esc(m.body) + '</div><span class="muted small">' +
            esc(mine ? t("tkYou") : m.author_name) + " · " + esc(date(m.created_at)) + "</span></div>";
        }).join("") + "</section>" +
        (x.status === "solved" ? '<p class="muted small">' + esc(t("tkReopen")) + "</p>" : "") +
        '<section class="panel"><label class="field">' + esc(t("tkReply")) + '<textarea id="tkBody" rows="3" maxlength="4000" placeholder="' + esc(t("tkReplyPh")) + '"></textarea></label>' +
        '<div class="row-btns"><button type="button" class="btn sm" id="tkReplyBtn">' + esc(t("tkSend")) + "</button>" +
        (x.status !== "solved" ? '<button type="button" class="link" id="tkSolveBtn">' + esc(t("tkSolve")) + "</button>" : "") + "</div></section>";
      view(html);
      var b = document.getElementById("tkReplyBtn");
      b.onclick = function () {
        var body = document.getElementById("tkBody").value.trim();
        if (!body) { document.getElementById("tkBody").focus(); return; }
        b.disabled = true;
        api("POST", "/v1/shop/tickets/" + id + "/messages", { body: body }).then(function () { drawThread(id); }, function (e) { b.disabled = false; toast(errMsg(e)); });
      };
      var s = document.getElementById("tkSolveBtn");
      if (s) s.onclick = function () {
        api("POST", "/v1/shop/tickets/" + id + "/solve").then(function () { toast(t("tkSolved")); drawThread(id); }, function (e) { toast(errMsg(e)); });
      };
    }, function (e) { view('<p class="empty">' + esc(errMsg(e)) + '</p><p class="tc"><a class="btn" href="#">' + esc(t("tkBack")) + "</a></p>"); });
  }
  k.withLogin(function () {
    window.addEventListener("hashchange", route);
    route();
  });
};
