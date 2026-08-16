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

(function () {
  "use strict";
  // applyOptimisticChrome previews the server-driven theme/scheme state on
  // <html> the instant the control changes. The server remains the authority:
  // the boosted form still submits the GET, and htmx:before:swap reconciles
  // class/data-theme from the server response. Without JS the native form
  // submit applies the same state on a full page load.
  function applyOptimisticChrome(form) {
    var root = document.documentElement;
    var scheme = form.querySelector('input[type="checkbox"][name="scheme"]');
    if (scheme) {
      if (scheme.checked) {
        root.classList.add("theme-dark");
        root.setAttribute("data-theme", "dark");
      } else {
        root.classList.remove("theme-dark");
        root.setAttribute("data-theme", "light");
      }
      keepPreservedState("scheme", scheme.checked ? "dark" : "light");
      return;
    }
    var theme = form.querySelector('select[name="theme"]');
    if (theme) {
      var next = null;
      for (var i = 0; i < theme.options.length; i++) {
        var cls = theme.options[i].getAttribute("data-class");
        if (!cls) continue;
        if (i === theme.selectedIndex) next = cls;
        root.classList.remove(cls);
      }
      if (next) root.classList.add(next);
      keepPreservedState("theme", theme.value);
    }
    refreshChromeHrefs();
  }
  // keepPreservedState keeps the OTHER chrome form's hidden preserve input in
  // sync. With hx-swap=none the body never re-renders, so the theme form's
  // hidden scheme (and the scheme form's hidden theme) would otherwise stay
  // stale and the next submission would silently forget the other parameter.
  function keepPreservedState(name, value) {
    var checkbox = document.querySelector('form[data-chrome-form] input[type="checkbox"][name="scheme"]');
    var select = document.querySelector('form[data-chrome-form] select[name="theme"]');
    var target = name === "scheme" ? (select && select.closest("form")) : (checkbox && checkbox.closest("form"));
    var hidden = target && target.querySelector('input[type="hidden"][name="' + name + '"]');
    if (hidden) hidden.value = value;
  }
  // refreshChromeHrefs rewrites every docs-shell chrome href (sidebar,
  // topbar, breadcrumb, prev/next) with the current theme/scheme query.
  // The optimistic toggle submits with hx-swap=none, so the body never
  // re-renders and the server-rendered hrefs from the ORIGINAL load would
  // otherwise stay stale: the next sidebar click would navigate without
  // ?scheme=dark and the server would render light, silently undoing the
  // optimistic toggle. This mirrors server-side chromeHref so navigation
  // preserves direction + light/dark without waiting for a re-render.
  function refreshChromeHrefs() {
    var root = document.documentElement;
    var query = "";
    var params = [];
    var themeSlug = "";
    var scheme = "";
    if (root.classList.contains("theme-dark")) scheme = "dark";
    else if (root.hasAttribute("data-theme") && root.getAttribute("data-theme") === "light") scheme = "light";
    var themeForm = document.querySelector('form[data-chrome-form] select[name="theme"]');
    if (themeForm && themeForm.value) themeSlug = themeForm.value;
    if (scheme) params.push("scheme=" + scheme);
    if (themeSlug) params.push("theme=" + themeSlug);
    if (params.length) query = "?" + params.join("&");
    var links = document.querySelectorAll(
      ".docs-sidebar a[href], .docs-topbar a[href], .ui-breadcrumb a[href], .docs-prev-next a[href]"
    );
    for (var i = 0; i < links.length; i++) {
      var a = links[i];
      var href = a.getAttribute("href");
      if (!href || href.charAt(0) !== "/") continue;
      if (href.indexOf("?") >= 0) href = href.slice(0, href.indexOf("?"));
      a.setAttribute("href", href + query);
    }
  }
  function initChrome(root) {
    var forms = (root || document).querySelectorAll("form[data-chrome-form]");
    for (var i = 0; i < forms.length; i++) {
      var form = forms[i];
      if (form.getAttribute("data-gelium-initialized")) continue;
      form.setAttribute("data-gelium-initialized", "true");
      var submit = form.querySelector('button[type="submit"]'); if (submit) submit.hidden = true;
      form.addEventListener("change", function () {
        applyOptimisticChrome(form);
        // requestSubmit (not submit) fires the submit event so htmx intercepts
        // the boosted GET; form.submit() performs a NATIVE full page load.
        this.requestSubmit();
      });
    }
  }
  function syncDocument(responseBody) {
    if (!responseBody) return;
    var parsed = new DOMParser().parseFromString(responseBody, "text/html");
    var source = parsed.documentElement, target = document.documentElement;
    if (!source) return;
    target.className = source.className;
    if (source.hasAttribute("data-theme")) target.setAttribute("data-theme", source.getAttribute("data-theme"));
    else target.removeAttribute("data-theme");
  }
  var lastScrollY = 0;
  document.addEventListener("htmx:before:swap", function (event) {
    var ctx = event.detail && event.detail.ctx;
    if (ctx && ctx.target === document.body && ctx.text) {
      syncDocument(ctx.text);
    }
  });
  document.addEventListener("htmx:before:request", function (event) {
    var source = event.detail && event.detail.elt;
    if (source && source.closest && source.closest("a[href^='#']")) return;
    lastScrollY = window.scrollY || document.documentElement.scrollTop || 0;
  });
  function afterSwap(event) {
    var detail = event.detail || {};
    if (detail.ctx && detail.ctx.target === document.body) {
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

(function () {
  "use strict";
  // On-this-page scrollspy (progressive enhancement): highlights the rail
  // section currently in view. The rail itself is a plain anchor list that
  // works without JS; this only toggles is-current/aria-current.
  function initOnThisPage() {
    var rail = document.querySelector(".docs-on-this-page");
    if (!rail || rail.getAttribute("data-gelium-spy")) return;
    if (!("IntersectionObserver" in window)) return;
    var links = rail.querySelectorAll('a[href^="#"]');
    if (!links.length) return;
    rail.setAttribute("data-gelium-spy", "true");
    var sections = [];
    links.forEach(function (link) {
      var target = document.getElementById(link.getAttribute("href").slice(1));
      if (target) sections.push({ link: link, target: target });
    });
    var observer = new IntersectionObserver(function (entries) {
      entries.forEach(function (entry) {
        if (!entry.isIntersecting) return;
        sections.forEach(function (s) {
          var active = s.target === entry.target;
          s.link.classList.toggle("is-current", active);
          if (active) s.link.setAttribute("aria-current", "true");
          else s.link.removeAttribute("aria-current");
        });
      });
    }, { rootMargin: "-20% 0px -60% 0px", threshold: 0 });
    sections.forEach(function (s) { observer.observe(s.target); });
  }
  document.addEventListener("DOMContentLoaded", initOnThisPage);
  document.addEventListener("htmx:after:swap", initOnThisPage);
})();
