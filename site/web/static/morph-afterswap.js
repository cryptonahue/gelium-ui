// Re-run page-level behavior after an HTMX morph swap.
(function () {
  "use strict";
  document.addEventListener("htmx:before:swap", function (event) {
    if (event.detail && event.detail.ctx && event.detail.ctx.target === document.body) {
      window._geliumPageGeneration = (window._geliumPageGeneration || 0) + 1;
    }
  });
})();
