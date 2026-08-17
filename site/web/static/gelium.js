// gelium.js — Gelium UI consumer enhancements (progressive enhancement).
// Ship this file (or bundle it) alongside the component styles:
//   - htmx 422 validation contract (X-Gelium-Validation)
//   - same-document view transitions (htmx 4 config flag, reduced-motion guard)
//   - gelium:toast region (server-rendered toasts + transport errors)
//   - ui-slider fill enhancement (--ui-slider-fill from native range input)
// No framework dependency beyond htmx 4; all features degrade gracefully
// without JS (server-rendered defaults).

document.addEventListener("htmx:before:swap", function (event) {
  var detail = event.detail || {};
  var response = detail.ctx && detail.ctx.response;
  if (
    response &&
    response.status === 422 &&
    response.headers.get("X-Gelium-Validation") === "true"
  ) {
    detail.shouldSwap = true;
    detail.isError = false;
  }
});

// Same-document view transitions (progressive enhancement): htmx 4 runs
// boosted navigations through document.startViewTransition when the config
// flag is set. Same-document VT is NOT auto-disabled under
// prefers-reduced-motion (unlike the cross-document at-rule), so the guard
// below is required (WCAG 2.3.3). No cross-document at-rule: automatic
// cross-document navigation is still flag-gated in Firefox, so VT stays
// same-document only.
if (document.startViewTransition && !window.matchMedia("(prefers-reduced-motion: reduce)").matches) { htmx.config.transitions = true; }

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
  document.addEventListener("htmx:response:error", function () { showToast("error", TRANSPORT_ERROR); });
  document.addEventListener("htmx:error", function () { showToast("error", TRANSPORT_ERROR); });
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

