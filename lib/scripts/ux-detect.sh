#!/usr/bin/env bash
# ux-detect.sh — deterministic Gelium UX contract / anti-slop checks for CONSUMER repos.
# Ships inside lib/ so npm consumers can run it without the monorepo.
# Exit 0 = clean; non-zero = failures. No LLM, just greps over your markup/styles.
#
# Usage: run from the consumer repo root:
#   bash node_modules/@gelium/ui/scripts/ux-detect.sh [paths...]
# Paths default to the whole repo excluding .git/node_modules/dist.
# Scans .html, .templ, .go (non-test) and .css sources; vendored bundles
# (gelium.css dist, anything under a vendor/ or node_modules/ path) are skipped.
set -uo pipefail
ROOT="$(pwd)"
TARGETS=("$@")
if [[ ${#TARGETS[@]} -eq 0 ]]; then
  mapfile -t TARGETS < <(find . \( -name .git -o -name node_modules -o -name dist \) -prune -o \
    \( -name '*.html' -o -name '*.go' -o -name '*.templ' -o -name '*.css' \) -print)
fi

# Source files actually shipped by the consumer: drop test fixtures and
# vendored/prebuilt bundles from every scan.
mapfile -t SOURCES < <(
  printf '%s\n' "${TARGETS[@]}" | grep -v -E '(_test\.go|\.test\.|\.spec\.)' | grep -v -E '(^|/)(node_modules|vendor|dist)/' | grep -v '/gelium\.css$' || true
)
if [[ ${#SOURCES[@]} -eq 0 ]]; then
  echo "ux-detect: no source files to scan" >&2
  exit 0
fi

fail=0
note() { printf '%s\n' "$*"; }
bad() {
  note "FAIL: $*"
  fail=1
}
ok() { note "OK: $*"; }
# Truncate each match so a minified one-line bundle can't flood stdout.
clip() { cut -c1-200; }

note "== Gelium UX detectors (consumer) =="

# Step-0 gate reminder: artifacts are not required to exist, but if they are absent
# the plain-language brief must establish the job, scope, and unresolved design
# decisions before generating UI. The script cannot verify a conversation; it only reminds.
if [[ ! -f PRODUCT.md || ! -f DESIGN.md ]]; then
  note "NOTE: no PRODUCT.md/DESIGN.md found — confirm the plain-language brief established job, scope, and unresolved design decisions before generating UI."
fi

# F-2: dark mode is class-routed; no media-query-only dark overrides with literals.
# CSS only — Go/HTML/templ files legitimately mention the string in comments/tests.
css_sources=($(printf '%s\n' "${SOURCES[@]}" | grep '\.css$' || true))
markup_sources=($(printf '%s\n' "${SOURCES[@]}" | grep -E '\.(html|templ|go)$' || true))
if [[ ${#css_sources[@]} -gt 0 ]] && rg -n '@media[^{]*prefers-color-scheme:\s*dark' "${css_sources[@]}" 2>/dev/null | clip; then
  bad "prefers-color-scheme: dark override found — dark mode must use the theme-dark class on <html>, not media queries"
else
  ok "no prefers-color-scheme dark overrides"
fi

# Dark surfaces must actually set the class somewhere in shipped markup.
if [[ ${#markup_sources[@]} -gt 0 ]] && rg -l 'theme-dark' "${markup_sources[@]}" 2>/dev/null | rg -q .; then
  ok "theme-dark class referenced"
else
  note "WARN: no 'theme-dark' usage found — if any surface supports dark, route it via <html class=\"... theme-dark\">."
fi

# F-1: registry-first page shells. Any consumer-defined nav/header/page-shell
# class with layout rules is a hand-rolled shell, sticky or not — flag names that
# are NOT registered Gelium components (ui-*).
shell_hits=""
for f in "${css_sources[@]:-}"; do
  [[ -z "$f" || ! -f "$f" ]] && continue
  hits=$(rg -n --no-heading '\.[a-zA-Z][a-zA-Z0-9_-]*(site-header|site-nav|page-header|topbar|appbar|navbar)[a-zA-Z0-9_-]*\s*\{' "$f" 2>/dev/null |
    rg -v '\.ui-' | clip || true)
  [[ -n "$hits" ]] && shell_hits+="$hits"$'\n'
done
if [[ -n "$shell_hits" ]]; then
  bad "custom page-shell classes found — page shells must compose registered components (ui-navigation-bar, ui-container); custom shell CSS limited to spacing/width:"
  note "$shell_hits" | head -10
else
  ok "no hand-rolled page shells detected"
fi

# Hard anti-pattern: overflow-x hidden on html/body (masking ban). Match across
# selector lines too: report any overflow-x:hidden whose file also styles html/body,
# flagged for manual review rather than silently passed.
if [[ ${#css_sources[@]} -gt 0 ]]; then
  hidden=$(rg -n --no-heading 'overflow-x:\s*hidden' "${css_sources[@]}" 2>/dev/null | clip || true)
  if [[ -n "$hidden" ]]; then
    if rg -q -i '(^|[,\s{}])(html|body)\s*[,{]' "${css_sources[@]}" 2>/dev/null && echo "$hidden" | rg -qi 'body|html'; then
      bad "overflow-x:hidden near html/body selectors (masking ban) — review manually:"
      note "$hidden" | head -5
    else
      note "NOTE: overflow-x:hidden found on non-root elements (allowed):"
      note "$hidden" | head -3
    fi
  else
    ok "no overflow-x:hidden"
  fi
fi

# One-off color literals in shipped markup/CSS outside token definitions.
lit=""
if [[ ${#markup_sources[@]} -gt 0 ]]; then
  lit+="$(rg -n --no-heading '#[0-9a-fA-F]{3,8}\b' "${markup_sources[@]}" 2>/dev/null | rg -v 'var\(|--ui-' | clip || true)"$'\n'
fi
if [[ ${#css_sources[@]} -gt 0 ]]; then
  # In CSS, skip lines that DEFINE tokens (--ui-*:) — those are the theme itself.
  lit+="$(rg -n --no-heading '#[0-9a-fA-F]{3,8}\b' "${css_sources[@]}" 2>/dev/null | rg -v -- '--ui-[a-z0-9-]+\s*:' | clip || true)"$'\n'
fi
lit=$(echo "$lit" | sed '/^$/d' | head -5)
if [[ -n "$lit" ]]; then
  bad "one-off color literals in shipped markup/CSS — use --ui-* tokens:"
  note "$lit"
else
  ok "no one-off color literals in markup/CSS"
fi

# Forms need validation summary hook. Forms live in .html AND .templ files;
# templ compiles to Go but the source of truth for review is the template.
form_files=()
while IFS= read -r f; do form_files+=("$f"); done < <(rg -l '<form' $(printf '%s\n' "${SOURCES[@]}" | grep -E '\.(html|templ)$') 2>/dev/null || true)
missing_forms=0
for f in "${form_files[@]:-}"; do
  [[ -z "$f" || ! -f "$f" ]] && continue
  if ! grep -qE 'validation-summary|X-Gelium-Validation|gelium-validation' "$f"; then
    bad "form without validation-summary contract: $f"
    missing_forms=1
  fi
done
[[ $missing_forms -eq 0 ]] && ok "all forms carry a validation-summary hook (${#form_files[@]} checked)"

if [[ "$fail" -ne 0 ]]; then
  note "== RESULT: FAILED =="
  exit 1
fi
note "== RESULT: PASSED =="
exit 0
