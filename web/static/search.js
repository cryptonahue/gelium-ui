(function () {
  "use strict";

  // Docs search enhancement: the search input is a real GET form to the /docs
  // hub (0-JS fallback: Enter submits ?q=<term>). When JavaScript is available
  // this filters the nav index emitted by the shell ({title, href, group} in
  // #docs-search-index) into a results list as you type. Without the index the
  // input degrades to the native form submission — no JS required.

  var script = document.getElementById("docs-search-index");
  if (!script) return;

  var INDEX;
  try {
    INDEX = JSON.parse(script.textContent);
  } catch (err) {
    return;
  }

  var input = document.getElementById("docs-search");
  var results = document.getElementById("docs-search-results");
  if (!input || !results || !INDEX || !INDEX.length) return;

  var MAX_HITS = 12;

  function matches(term, entry) {
    return (
      entry.title.toLowerCase().indexOf(term) !== -1 ||
      entry.group.toLowerCase().indexOf(term) !== -1
    );
  }

  function hide() {
    results.hidden = true;
    input.setAttribute("aria-expanded", "false");
  }

  function render(term) {
    results.innerHTML = "";
    var hits = 0;
    for (var i = 0; i < INDEX.length && hits < MAX_HITS; i++) {
      if (!matches(term, INDEX[i])) continue;
      var li = document.createElement("li");
      var link = document.createElement("a");
      link.href = INDEX[i].href;
      var title = document.createElement("span");
      title.className = "docs-search-result-title";
      title.textContent = INDEX[i].title;
      var group = document.createElement("span");
      group.className = "docs-search-result-group";
      group.textContent = INDEX[i].group;
      link.appendChild(title);
      link.appendChild(group);
      li.appendChild(link);
      results.appendChild(li);
      hits++;
    }
    results.hidden = hits === 0;
    input.setAttribute("aria-expanded", hits > 0 ? "true" : "false");
  }

  input.addEventListener("input", function () {
    var term = input.value.trim().toLowerCase();
    if (!term) {
      hide();
      return;
    }
    render(term);
  });

  input.addEventListener("keydown", function (event) {
    if (event.key !== "Enter") return;
    if (results.hidden) return; // fall back to the native GET /docs?q= submission
    var first = results.querySelector("a");
    if (!first) return;
    event.preventDefault();
    window.location.assign(first.href);
  });

  input.addEventListener("blur", function () {
    // Delayed so a mousedown on a result still navigates.
    window.setTimeout(hide, 150);
  });

  // Keep input focus while interacting with the results (mousedown precedes
  // blur); the native link click then navigates normally.
  results.addEventListener("mousedown", function (event) {
    event.preventDefault();
  });
})();
