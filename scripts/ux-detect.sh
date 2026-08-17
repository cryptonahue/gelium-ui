#!/usr/bin/env bash
# ux-detect.sh — deterministic Gelium UX contract / anti-slop checks (no LLM).
# Exit 0 = clean; non-zero = failures. Ethos: contracts + tokens, not aesthetic bans on the design system.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
fail=0
note() { printf '%s\n' "$*"; }
bad() { note "FAIL: $*"; fail=1; }
ok() { note "OK: $*"; }

note "== Gelium UX detectors =="

# Required handbook + agent pack
required=(
  site/web/content/handbook-screens.md
  site/web/content/handbook-feedback.md
  site/web/content/handbook-journeys.md
  site/web/content/handbook-data-display.md
  site/web/content/handbook-agent-workflow.md
  site/web/content/handbook-ui-definition-of-done.md
  site/web/static/llms-ux.txt
  site/web/static/llms.txt
  site/web/content/templates/product.md
  site/web/content/templates/design.md
)
for f in "${required[@]}"; do
  if [[ -f "$f" ]]; then ok "exists $f"; else bad "missing $f"; fi
done

# Contract IDs must remain in agent pack
for id in FEED-VAL JOURNEY-LINEAR DATA-TABLE WF-SHAPE SURFACE Operate SKEL-FORUM; do
  if grep -q "$id" site/web/static/llms-ux.txt; then ok "llms-ux has $id"; else bad "llms-ux missing $id"; fi
done

# Human workflow page must name passes + ethos
for s in "WF-SHAPE" "Surface modes" "Anti-slop" "Ethos"; do
  if grep -q "$s" site/web/content/handbook-agent-workflow.md; then ok "workflow has $s"; else bad "workflow missing $s"; fi
done

# Hard anti-pattern: overflow-x hidden on html/body (masking) in lib + site styles
# Allow in comments only — flag real property use on html/body selectors is complex;
# flag global-looking rules with overflow-x:\s*hidden on the same line as html or body.
if rg -n --glob '*.css' '^(html|body)[^{]*\{[^}]*overflow-x:\s*hidden' lib/styles site/web/styles 2>/dev/null; then
  bad "overflow-x:hidden on html/body rule (masking ban)"
else
  # also catch simple "body { ... overflow-x: hidden" multiline lightly via fixed strings
  if rg -n --glob '*.css' 'overflow-x:\s*hidden' lib/styles site/web/styles 2>/dev/null | rg -i 'html|body' >/dev/null; then
    # warn-level: list but only fail if clearly body/html selector lines nearby — keep soft
    note "WARN: overflow-x:hidden appears near html/body mentions — review manually"
  fi
  ok "no simple html/body overflow-x:hidden rule detected"
fi

# Toast-only validation should stay forbidden in docs
if grep -qi 'toast-only' site/web/content/handbook-feedback.md; then ok "feedback documents toast-only ban"; else bad "feedback lost toast-only guidance"; fi

# llms.txt must point at ux pack
if grep -q 'llms-ux.txt' site/web/static/llms.txt; then ok "llms.txt links llms-ux"; else bad "llms.txt missing llms-ux link"; fi

if [[ "$fail" -ne 0 ]]; then
  note "== RESULT: FAILED =="
  exit 1
fi
note "== RESULT: PASSED =="
exit 0
