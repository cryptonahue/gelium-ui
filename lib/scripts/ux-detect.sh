#!/usr/bin/env bash
# ux-detect.sh — deterministic Gelium UX contract / anti-slop checks for CONSUMER repos.
# Ships inside lib/ so npm consumers can run it without the monorepo.
# Exit 0 = clean; non-zero = failures. No LLM, just greps over your markup/styles.
#
# Usage: run from the consumer repo root:
#   bash node_modules/@gelium/ui/scripts/ux-detect.sh [paths...]
# Paths default to the whole repo excluding .git/node_modules/dist.
set -uo pipefail
ROOT="$(pwd)"
TARGETS=("$@")
if [[ ${#TARGETS[@]} -eq 0 ]]; then
  mapfile -t TARGETS < <(find . \( -name .git -o -name node_modules -o -name dist \) -prune -o \
    \( -name '*.html' -o -name '*.go' -o -name '*.templ' -o -name '*.css' \) -print)
fi
fail=0
note() { printf '%s\n' "$*"; }
bad() { note "FAIL: $*"; fail=1; }
ok() { note "OK: $*"; }

note "== Gelium UX detectors (consumer) =="

# Step-0 gate reminder: artifacts are not required to exist, but if they are absent
# you must have asked the user before generating UI. The script cannot verify a
# conversation; it only reminds.
if [[ ! -f PRODUCT.md && ! -f DESIGN.md ]]; then
  note "NOTE: no PRODUCT.md/DESIGN.md found — confirm you asked the user about job, SURFACE mode, and visual direction before generating UI."
fi

# F-2: dark mode is class-routed; no media-query-only dark overrides with literals
if rg -n --glob '*.css' '@media[^{]*prefers-color-scheme:\s*dark' "${TARGETS[@]}" 2>/dev/null; then
  bad "prefers-color-scheme: dark override found — dark mode must use the theme-dark class on <html>, not media queries"
else
  ok "no prefers-color-scheme dark overrides"
fi

# Dark surfaces must actually set the class somewhere in shipped markup
if rg -l 'theme-dark' --glob '*.html' --glob '*.go' "${TARGETS[@]}" 2>/dev/null | rg -q .; then
  ok "theme-dark class referenced"
else
  note "WARN: no 'theme-dark' usage found — if any surface supports dark, route it via <html class=\"... theme-dark\">."
fi

# F-1: registry-first page shells — flag hand-rolled nav headers
# A page shell that defines its own sticky/fixed header/nav CSS instead of composing
# ui-navigation-bar / drawer primitives.
if rg -ln --glob '*.css' '(nav|header)[^{]*\{[^}]*(position:\s*(sticky|fixed))' "${TARGETS[@]}" 2>/dev/null; then
  bad "custom sticky/fixed header or nav CSS found — page shells must compose registered components (ui-navigation-bar, ui-container); custom shell CSS limited to spacing/width"
else
  ok "no hand-rolled sticky/fixed page shells detected"
fi

# Hard anti-pattern: overflow-x hidden on html/body (masking ban)
if rg -n --glob '*.css' 'overflow-x:\s*hidden' "${TARGETS[@]}" 2>/dev/null | rg -i 'html|body' >/dev/null; then
  bad "possible overflow-x:hidden on html/body (masking ban) — review manually"
else
  ok "no overflow-x:hidden near html/body"
fi

# One-off color literals in shipped markup/CSS outside token definitions
lit=$(rg -n --glob '*.html' --glob '*.go' '#[0-9a-fA-F]{3,8}\b' "${TARGETS[@]}" 2>/dev/null | rg -v 'var\(|--ui-' | head -5 || true)
if [[ -n "$lit" ]]; then
  bad "one-off color literals in markup — use --ui-* tokens:"
  note "$lit"
else
  ok "no one-off color literals in markup"
fi

# Forms need validation summary hook
for f in $(rg -l --glob '*.html' '<form' "${TARGETS[@]}" 2>/dev/null); do
  if ! grep -qE 'validation-summary|X-Gelium-Validation|gelium-validation' "$f"; then
    bad "form without validation-summary contract: $f"
  fi
done
[[ $fail -eq 0 ]] && ok "all forms carry a validation-summary hook"

if [[ "$fail" -ne 0 ]]; then
  note "== RESULT: FAILED =="
  exit 1
fi
note "== RESULT: PASSED =="
exit 0
