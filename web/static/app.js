document.addEventListener("htmx:beforeSwap", function (event) {
  if (
    event.detail.xhr.status === 422 &&
    event.detail.xhr.getResponseHeader("X-Gelium-Validation") === "true"
  ) {
    event.detail.shouldSwap = true;
    event.detail.isError = false;
  }
});

(function () {
  "use strict";

  var REGION_SELECTOR = "#gelium-toast-region";
  var DEFAULT_TIMEOUT = 4000;
  var ERROR_TIMEOUT = 8000;
  var KNOWN_TYPES = { info: true, success: true, warning: true, error: true };

  function region() { return document.querySelector(REGION_SELECTOR); }
  function dismiss(el) {
    if (el.getAttribute("data-gelium-toast-done")) return;
    el.setAttribute("data-gelium-toast-done", "true");
    if (el._timer) clearTimeout(el._timer);
    el.classList.remove("ui-toast-show");
    window.setTimeout(function () { if (el.parentNode) el.parentNode.removeChild(el); }, 300);
  }
  function schedule(el, ms) {
    if (el.getAttribute("data-gelium-toast-done")) return;
    if (el._timer) clearTimeout(el._timer);
    el._timer = window.setTimeout(function () { dismiss(el); }, ms);
  }
  function makeToast(type, message) {
    var normalized = KNOWN_TYPES[type] ? type : "info";
    var el = document.createElement("div");
    el.className = "ui-toast ui-toast-" + normalized;
    el.setAttribute("role", normalized === "error" ? "alert" : "status");
    el._timeoutMs = normalized === "error" ? ERROR_TIMEOUT : DEFAULT_TIMEOUT;
    var msg = document.createElement("span"); msg.className = "ui-toast-message"; msg.textContent = message; el.appendChild(msg);
    var btn = document.createElement("button"); btn.type = "button"; btn.className = "ui-toast-action"; btn.setAttribute("aria-label", "Dismiss notification"); btn.textContent = "Dismiss";
    btn.addEventListener("click", function () { dismiss(el); }); el.appendChild(btn);
    var resume = function () { schedule(el, el._timeoutMs); };
    el.addEventListener("mouseenter", function () { clearTimeout(el._timer); }); el.addEventListener("mouseleave", resume);
    el.addEventListener("focusin", function () { clearTimeout(el._timer); }); el.addEventListener("focusout", resume);
    return el;
  }
  function showToast(type, message) {
    var target = region(); if (!target) return;
    var el = makeToast(type, message); target.appendChild(el);
    window.requestAnimationFrame(function () { el.classList.add("ui-toast-show"); }); schedule(el, el._timeoutMs);
  }
  document.addEventListener("gelium:toast", function (event) { var d = event.detail || {}; if (d.message) showToast(d.type, d.message); });
  var TRANSPORT_ERROR = "We couldn't reach the server. Try again.";
  document.addEventListener("htmx:responseError", function () { showToast("error", TRANSPORT_ERROR); });
  document.addEventListener("htmx:sendError", function () { showToast("error", TRANSPORT_ERROR); });
})();

(function () {
  "use strict";
  function toPercent(input) {
    var min = input.min === "" ? 0 : Number(input.min), max = input.max === "" ? 100 : Number(input.max);
    return ((Number(input.value) - min) / (max - min || 1)) * 100;
  }
  document.addEventListener("input", function (event) {
    var input = event.target; if (!input || input.type !== "range") return;
    var slider = input.closest(".ui-slider[data-ui-slider]"); if (slider) slider.style.setProperty("--ui-slider-fill", toPercent(input) + "%");
  });
})();

(function () {
  "use strict";
  function initChrome(root) {
    var forms = (root || document).querySelectorAll("form[data-chrome-form]");
    for (var i = 0; i < forms.length; i++) {
      var form = forms[i];
      if (form.getAttribute("data-gelium-initialized")) continue;
      form.setAttribute("data-gelium-initialized", "true");
      var submit = form.querySelector('button[type="submit"]'); if (submit) submit.hidden = true;
      form.addEventListener("change", function () { this.submit(); });
    }
  }
  function syncDocument(responseText) {
    if (!responseText) return;
    var parsed = new DOMParser().parseFromString(responseText, "text/html");
    var source = parsed.documentElement, target = document.documentElement;
    if (!source) return;
    target.className = source.className;
    if (source.hasAttribute("data-theme")) target.setAttribute("data-theme", source.getAttribute("data-theme"));
    else target.removeAttribute("data-theme");
  }
  var lastScrollY = 0;
  document.addEventListener("htmx:beforeRequest", function (event) {
    var source = event.detail && event.detail.elt;
    if (source && source.closest && source.closest("a[href^='#']")) return;
    lastScrollY = window.scrollY || document.documentElement.scrollTop || 0;
  });
  document.addEventListener("htmx:beforeSwap", function (event) {
    if (event.detail && event.detail.xhr && event.detail.ctx && event.detail.ctx.target === document.body) {
      event.detail._geliumHtml = event.detail.xhr.responseText || "";
    }
  });
  function afterSwap(event) {
    var detail = event.detail || {};
    if (detail.ctx && detail.ctx.target === document.body) {
      syncDocument(detail._geliumHtml || (detail.xhr && detail.xhr.responseText));
      var main = document.getElementById("main-content");
      if (main) { main.setAttribute("tabindex", "-1"); main.focus({ preventScroll: true }); }
      if (lastScrollY > 0) window.scrollTo(0, lastScrollY);
    }
    initChrome(document);
  }
  document.addEventListener("htmx:after:swap", afterSwap);
  document.addEventListener("htmx:before:history:restore", function () {
    lastScrollY = 0;
  });
  document.addEventListener("DOMContentLoaded", function () { initChrome(document); });
})();
