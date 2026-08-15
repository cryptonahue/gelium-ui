document.addEventListener("htmx:beforeSwap", function (event) {
  if (
    event.detail.xhr.status === 422 &&
    event.detail.xhr.getResponseHeader("X-Loom-Validation") === "true"
  ) {
    event.detail.shouldSwap = true;
    event.detail.isError = false;
  }
});

(function () {
  "use strict";

  // Minimal, framework-free enhancement for the Toast component. It only makes
  // server-driven feedback transient; the no-JS inline flow is complete without it.
  var REGION_SELECTOR = "#loom-toast-region";
  var DEFAULT_TIMEOUT = 4000;
  var ERROR_TIMEOUT = 8000;
  var KNOWN_TYPES = { info: true, success: true, warning: true, error: true };

  function region() {
    return document.querySelector(REGION_SELECTOR);
  }

  function dismiss(el) {
    if (el.getAttribute("data-loom-toast-done")) return;
    el.setAttribute("data-loom-toast-done", "true");
    if (el._timer) clearTimeout(el._timer);
    el.classList.remove("ui-toast-show");
    window.setTimeout(function () {
      if (el.parentNode) el.parentNode.removeChild(el);
    }, 300);
  }

  function schedule(el, ms) {
    if (el.getAttribute("data-loom-toast-done")) return;
    if (el._timer) clearTimeout(el._timer);
    el._timer = window.setTimeout(function () { dismiss(el); }, ms);
  }

  function makeToast(type, message) {
    var normalized = KNOWN_TYPES[type] ? type : "info";
    var el = document.createElement("div");
    el.className = "ui-toast ui-toast-" + normalized;
    el.setAttribute("role", normalized === "error" ? "alert" : "status");
    el._timeoutMs = normalized === "error" ? ERROR_TIMEOUT : DEFAULT_TIMEOUT;

    var msg = document.createElement("span");
    msg.className = "ui-toast-message";
    msg.textContent = message;
    el.appendChild(msg);

    var btn = document.createElement("button");
    btn.type = "button";
    btn.className = "ui-toast-action";
    btn.setAttribute("aria-label", "Dismiss notification");
    btn.textContent = "Dismiss";
    btn.addEventListener("click", function () { dismiss(el); });
    el.appendChild(btn);

    var resume = function () { schedule(el, el._timeoutMs); };
    el.addEventListener("mouseenter", function () { clearTimeout(el._timer); });
    el.addEventListener("mouseleave", resume);
    el.addEventListener("focusin", function () { clearTimeout(el._timer); });
    el.addEventListener("focusout", resume);
    return el;
  }

  function showToast(type, message) {
    if (!region()) return;
    var el = makeToast(type, message);
    region().appendChild(el);
    window.requestAnimationFrame(function () { el.classList.add("ui-toast-show"); });
    schedule(el, el._timeoutMs);
  }

  document.addEventListener("loom:toast", function (event) {
    var detail = event.detail || {};
    if (!detail.message) return;
    showToast(detail.type, detail.message);
  });

  // Transport-level HTMX failures (network down, 5xx) carry no server-driven
  // toast, so surface a transient generic error instead of failing silently.
  // 422-with-header validation never reaches here: the beforeSwap hook above
  // already cleared isError for that case, keeping the "validation is never a
  // toast" contract intact (G5).
  var TRANSPORT_ERROR = "We couldn't reach the server. Try again.";
  document.addEventListener("htmx:responseError", function () {
    showToast("error", TRANSPORT_ERROR);
  });
  document.addEventListener("htmx:sendError", function () {
    showToast("error", TRANSPORT_ERROR);
  });
})();

(function () {
  "use strict";

  // Minimal, framework-free enhancement for the Slider component: it keeps the
  // --ui-slider-fill percentage custom property in sync with the native range
  // while dragging, so the WebKit active-track fill follows the handle. Firefox
  // fills natively through ::-moz-range-progress. The no-JS flow is complete
  // without this: the fill shows the served value and the native input stays
  // fully operable.
  var toPercent = function (input) {
    var min = input.min === "" ? 0 : Number(input.min);
    var max = input.max === "" ? 100 : Number(input.max);
    var span = max - min || 1;
    return ((Number(input.value) - min) / span) * 100;
  };

  document.addEventListener("input", function (event) {
    var input = event.target;
    if (!input || input.type !== "range") return;
    var slider = input.closest(".ui-slider[data-ui-slider]");
    if (!slider) return;
    slider.style.setProperty("--ui-slider-fill", toPercent(input) + "%");
  });
})();

(function () {
  "use strict";

  // Docs chrome enhancement: the theme select and the scheme switch live in
  // 0-JS GET forms with a real submit button. With JavaScript available,
  // changing either control submits the form directly and the button hides;
  // without JavaScript the button still performs the same GET submission.
  // This only submits forms — the server still decides the direction.
  var FORMS = document.querySelectorAll("form[data-chrome-form]");
  for (var i = 0; i < FORMS.length; i++) {
    (function (form) {
      var submit = form.querySelector('button[type="submit"]');
      if (submit) submit.hidden = true;
      form.addEventListener("change", function () {
        form.submit();
      });
    })(FORMS[i]);
  }
})();
