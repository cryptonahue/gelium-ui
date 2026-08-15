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

  // Minimal, framework-free enhancement for the Toast component. It only makes
  // server-driven feedback transient; the no-JS inline flow is complete without it.
  var REGION_SELECTOR = "#gelium-toast-region";
  var DEFAULT_TIMEOUT = 4000;
  var ERROR_TIMEOUT = 8000;
  var KNOWN_TYPES = { info: true, success: true, warning: true, error: true };

  function region() {
    return document.querySelector(REGION_SELECTOR);
  }

  function dismiss(el) {
    if (el.getAttribute("data-gelium-toast-done")) return;
    el.setAttribute("data-gelium-toast-done", "true");
    if (el._timer) clearTimeout(el._timer);
    el.classList.remove("ui-toast-show");
    window.setTimeout(function () {
      if (el.parentNode) el.parentNode.removeChild(el);
    }, 300);
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

  document.addEventListener("gelium:toast", function (event) {
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

(function () {
  "use strict";

  // hx-boost navigation enhancement. The body is boosted (hx-boost="true"),
  // so internal links and GET forms navigate via AJAX and htmx swaps only
  // the <body>. Three concerns need explicit handling:

  // 1. Theme/scheme live on <html class> (and data-theme), which htmx does
  //    NOT swap. The boosted response carries the server-decided html; copy
  //    its class + data-theme onto the real <html> so a theme or scheme
  //    change survives the swap. 0-JS fallback is untouched.
  function syncDocumentClass() {
    var html = document.documentElement;
    var response = document.body.getAttribute("data-gelium-boost-html");
    if (!response) return;
    var match = /<html[^>]*class="([^"]*)"/.exec(response);
    if (match && html.className !== match[1]) html.className = match[1];
    var themeMatch = /<html[^>]*data-theme="([^"]*)"/.exec(response);
    if (themeMatch) {
      if (themeMatch[1]) html.setAttribute("data-theme", themeMatch[1]);
      else html.removeAttribute("data-theme");
    }
  }

  // 2. Scroll position: without JS htmx restores to the top of the target;
  //    for docs navigation we prefer to keep the reader's place (e.g. a
  //    sidebar click while scrolled deep). Restore the previous scrollY
  //    after the swap, falling back to htmx defaults for fragment links.
  var lastScrollY = 0;
  document.body.addEventListener("htmx:beforeRequest", function () {
    lastScrollY = window.scrollY || document.documentElement.scrollTop || 0;
  });
  // Keep the raw response html available for the class sync: htmx swaps the
  // <body> but the theme/scheme live on <html>, so we stash the response on
  // the body before the swap and read it in htmx:afterSwap.
  document.body.addEventListener("htmx:beforeSwap", function (event) {
    if (event.detail && event.detail.xhr) {
      document.body.setAttribute("data-gelium-boost-html", event.detail.xhr.responseText || "");
    }
  });
  document.body.addEventListener("htmx:afterSwap", function () {
    syncDocumentClass();
    var h = document.getElementById("main-content");
    if (h) {
      // Move focus into the new content so keyboard/screen-reader users
      // land on the page they navigated to, not the stale link.
      h.setAttribute("tabindex", "-1");
      h.focus({ preventScroll: true });
    }
    if (lastScrollY > 0) {
      window.scrollTo(0, lastScrollY);
    }
    document.body.removeAttribute("data-gelium-boost-html");
  });
})();
