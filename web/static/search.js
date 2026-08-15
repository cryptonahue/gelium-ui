(function () {
  "use strict";

  function initSearch(root) {
    var scope = root || document;
    var script = scope.querySelector ? scope.querySelector("#docs-search-index") : null;
    var input = scope.querySelector ? scope.querySelector("#docs-search") : null;
    var results = scope.querySelector ? scope.querySelector("#docs-search-results") : null;
    if (!script || !input || !results || input.getAttribute("data-gelium-search-initialized")) return;
    input.setAttribute("data-gelium-search-initialized", "true");
    var index;
    try { index = JSON.parse(script.textContent); } catch (err) { return; }
    if (!index || !index.length) return;
    var maxHits = 12;
    function hide() { results.hidden = true; input.setAttribute("aria-expanded", "false"); }
    function render(term) {
      results.innerHTML = "";
      var hits = 0;
      for (var i = 0; i < index.length && hits < maxHits; i++) {
        var entry = index[i];
        if (entry.title.toLowerCase().indexOf(term) === -1 && entry.group.toLowerCase().indexOf(term) === -1) continue;
        var li = document.createElement("li"), link = document.createElement("a");
        link.href = entry.href;
        var title = document.createElement("span"); title.className = "docs-search-result-title"; title.textContent = entry.title;
        var group = document.createElement("span"); group.className = "docs-search-result-group"; group.textContent = entry.group;
        link.appendChild(title); link.appendChild(group); li.appendChild(link); results.appendChild(li); hits++;
      }
      results.hidden = hits === 0; input.setAttribute("aria-expanded", hits > 0 ? "true" : "false");
    }
    input.addEventListener("input", function () { var term = input.value.trim().toLowerCase(); if (term) render(term); else hide(); });
    input.addEventListener("keydown", function (event) {
      if (event.key !== "Enter" || results.hidden) return;
      var first = results.querySelector("a"); if (!first) return;
      event.preventDefault(); window.location.assign(first.href);
    });
    input.addEventListener("blur", function () { window.setTimeout(hide, 150); });
    results.addEventListener("mousedown", function (event) { event.preventDefault(); });
  }
  document.addEventListener("DOMContentLoaded", function () { initSearch(document); });
  document.addEventListener("htmx:after:swap", function (event) { initSearch(event.target || document); });
  document.addEventListener("htmx:before:history:restore", function () { initSearch(document); });
})();
