// reports: tombol "Laporkan" di halaman produk & toko → modal kecil → POST /v1/reports. Tanpa library.
(function () {
  "use strict";
  var L = (window.PHSAR && window.PHSAR.lang) || "en";
  var T = {
    en: { title: "Report this {t}", product: "product", store: "store", why: "What is wrong?", note: "Details (optional)", noteReq: "Details (required for “Other”)",
      send: "Send report", cancel: "Cancel", thanks: "Thank you. The Phsar team will review it.", signin: "Please sign in to report.", signinBtn: "Sign in",
      already: "You already reported this. We are reviewing it.", limit: "Too many reports today. Please try again tomorrow.", error: "Something went wrong. Please try again.",
      r_counterfeit: "Fake or counterfeit", r_prohibited: "Prohibited item", r_misleading: "Misleading photos or description", r_scam: "Scam or fraud",
      r_offensive: "Offensive content", r_other: "Other" },
    km: { title: "រាយការណ៍{t}នេះ", product: "ផលិតផល", store: "ហាង", why: "មានបញ្ហាអ្វី?", note: "ព័ត៌មានលម្អិត (ស្រេចចិត្ត)", noteReq: "ព័ត៌មានលម្អិត (ចាំបាច់សម្រាប់ «ផ្សេងៗ»)",
      send: "ផ្ញើរបាយការណ៍", cancel: "បោះបង់", thanks: "សូមអរគុណ។ ក្រុម Phsar នឹងពិនិត្យ។", signin: "សូមចូលគណនីដើម្បីរាយការណ៍។", signinBtn: "ចូល",
      already: "អ្នកបានរាយការណ៍រួចហើយ។ យើងកំពុងពិនិត្យ។", limit: "រាយការណ៍ច្រើនពេកថ្ងៃនេះ។ សូមព្យាយាមម្តងទៀតថ្ងៃស្អែក។", error: "មានបញ្ហា។ សូមព្យាយាមម្តងទៀត។",
      r_counterfeit: "ក្លែងក្លាយ", r_prohibited: "ទំនិញហាមឃាត់", r_misleading: "រូបភាព ឬការពិពណ៌នាបំភាន់", r_scam: "ការបោកប្រាស់",
      r_offensive: "មាតិកាមិនសមរម្យ", r_other: "ផ្សេងៗ" },
    zh: { title: "举报此{t}", product: "商品", store: "店铺", why: "有什么问题？", note: "详细说明（可选）", noteReq: "详细说明（选择“其他”时必填）",
      send: "提交举报", cancel: "取消", thanks: "谢谢，Phsar 团队会进行审核。", signin: "请先登录再举报。", signinBtn: "登录",
      already: "您已举报过，我们正在审核。", limit: "今天举报次数过多，请明天再试。", error: "出错了，请重试。",
      r_counterfeit: "假冒伪劣", r_prohibited: "违禁商品", r_misleading: "图片或描述误导", r_scam: "诈骗",
      r_offensive: "不当内容", r_other: "其他" }
  }[L] || {};
  var REASONS = ["counterfeit", "prohibited", "misleading", "scam", "offensive", "other"];
  var css = ".report-link{border:0;background:none;color:var(--muted,#697386);font-size:13px;text-decoration:underline;text-underline-offset:2px;cursor:pointer;padding:6px 0;margin-top:6px}" +
    ".rpt-dlg{border:0;border-radius:20px;padding:0;max-width:440px;width:calc(100% - 32px)}.rpt-dlg::backdrop{background:rgba(16,24,40,.45)}" +
    ".rpt-dlg form{padding:20px;display:flex;flex-direction:column;gap:12px}.rpt-dlg h2{margin:0;font-size:18px}" +
    ".rpt-dlg label.r{display:flex;gap:10px;align-items:center;font-size:15px;padding:4px 0}.rpt-dlg input[type=radio]{width:18px;height:18px;accent-color:var(--pink,#e5195e)}" +
    ".rpt-dlg textarea{font:inherit;font-size:15px;border:1px solid #e6e3ea;border-radius:12px;padding:10px;min-height:80px;resize:vertical}" +
    ".rpt-dlg .row{display:flex;gap:10px;justify-content:flex-end;flex-wrap:wrap}.rpt-dlg .msg{margin:0;font-size:14px}";
  function t(k, a) { var s = T[k] || k; for (var x in a || {}) s = s.split("{" + x + "}").join(a[x]); return s; }
  function esc(s) { return String(s == null ? "" : s).replace(/[&<>"']/g, function (c) { return { "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#39;" }[c]; }); }
  function session() {
    return fetch("/shop/auth/refresh", { method: "POST", headers: { "X-Phsar": "1", "Content-Type": "application/json" }, body: "{}", credentials: "same-origin" })
      .then(function (r) { return r.ok ? r.json() : null; }, function () { return null; });
  }
  function open(kind, id) {
    var d = document.createElement("dialog");
    d.className = "rpt-dlg";
    d.innerHTML = '<form method="dialog"><h2>' + esc(t("title", { t: t(kind) })) + '</h2><p class="msg muted" id="rptMsg" hidden></p>' +
      '<fieldset style="border:0;padding:0;margin:0"><legend class="muted" style="font-size:14px;margin-bottom:6px">' + esc(t("why")) + "</legend>" +
      REASONS.map(function (r) { return '<label class="r"><input type="radio" name="rpt" value="' + r + '"> ' + esc(t("r_" + r)) + "</label>"; }).join("") +
      '</fieldset><label class="muted" style="font-size:14px;display:flex;flex-direction:column;gap:6px"><span id="rptNoteLbl">' + esc(t("note")) +
      '</span><textarea id="rptNote" maxlength="1000"></textarea></label><div class="row"><button type="button" class="btn ghost" value="cancel" id="rptCancel">' +
      esc(t("cancel")) + '</button><button type="button" class="btn" id="rptSend">' + esc(t("send")) + "</button></div></form>";
    document.body.appendChild(d);
    d.addEventListener("close", function () { d.remove(); });
    d.showModal();
    var msg = d.querySelector("#rptMsg"), send = d.querySelector("#rptSend");
    var say = function (k) { msg.textContent = t(k); msg.hidden = false; };
    d.querySelector("#rptCancel").onclick = function () { d.close(); };
    d.querySelectorAll('input[name="rpt"]').forEach(function (i) {
      i.onchange = function () { d.querySelector("#rptNoteLbl").textContent = t(i.value === "other" ? "noteReq" : "note"); };
    });
    send.onclick = function () {
      var r = d.querySelector('input[name="rpt"]:checked'), note = d.querySelector("#rptNote").value.trim();
      if (!r) { say("why"); return; }
      if (r.value === "other" && note.length < 5) { say("noteReq"); d.querySelector("#rptNote").focus(); return; }
      send.disabled = true;
      session().then(function (s) {
        if (!s || !s.access_token || (s.user && s.user.is_guest)) {
          say("signin");
          send.disabled = false; send.textContent = t("signinBtn");
          send.onclick = function () { location.href = "/" + L + "/orders"; };
          return;
        }
        return fetch("/v1/reports", { method: "POST", credentials: "same-origin",
          headers: { "Content-Type": "application/json", Authorization: "Bearer " + s.access_token, "X-Lang": L },
          body: JSON.stringify({ target: kind, id: id, reason: r.value, note: note }) }).then(function (res) {
          if (res.ok) {
            d.querySelector("fieldset").remove(); d.querySelector("label.muted").remove(); send.remove();
            d.querySelector("#rptCancel").textContent = "OK"; say("thanks"); return;
          }
          send.disabled = false;
          say(res.status === 409 ? "already" : res.status === 429 ? "limit" : "error");
        });
      }).catch(function () { send.disabled = false; say("error"); });
    };
  }
  document.addEventListener("DOMContentLoaded", function () {
    var btns = document.querySelectorAll("[data-report]");
    if (!btns.length) return;
    var st = document.createElement("style"); st.textContent = css; document.head.appendChild(st);
    btns.forEach(function (b) {
      b.textContent = b.textContent || t("title", { t: t(b.dataset.report) });
      b.addEventListener("click", function () { open(b.dataset.report, b.dataset.id); });
    });
  });
})();
