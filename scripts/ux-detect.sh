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
  site/web/content/handbook-page-section-architecture.md
  lib/skills/10-page-section-architecture.md
  lib/llms-ux.txt
  lib/llms.txt
  site/web/static/llms-ux.txt
  site/web/static/llms.txt
  site/web/content/templates/product.md
  site/web/content/templates/design.md
)
for f in "${required[@]}"; do
  if [[ -f "$f" ]]; then ok "exists $f"; else bad "missing $f"; fi
done

# Contract IDs must remain in the package and served agent-facing surfaces.
protocol_surfaces=(
  lib/skills/10-page-section-architecture.md
  lib/llms-ux.txt
  lib/llms.txt
  site/web/static/llms-ux.txt
  site/web/static/llms.txt
)
protocol_ids=(
  ARCH-PRODUCT ARCH-PAGE ARCH-SECTION ARCH-COMPONENTS ARCH-TOKENS
  SECTION-CONTRACT SECTION-HIERARCHY SECTION-ACTION SECTION-REVELATION SECTION-RECOVERY
  WF-ARCH WF-SECTION-AUDIT
)
for f in "${protocol_surfaces[@]}"; do
  for id in "${protocol_ids[@]}"; do
    if grep -Fq "$id" "$f"; then ok "$f has $id"; else bad "$f missing $id"; fi
  done
done

# The two existing recipes are the protocol's worked applications; keep them
# discoverable in every package and served agent-facing surface.
for f in "${protocol_surfaces[@]}"; do
  for recipe in /recipes/public-feed /recipes/rich-article; do
    if grep -Fq "$recipe" "$f"; then ok "$f names $recipe"; else bad "$f missing worked application $recipe"; fi
  done
done

# Content structure grammar
if grep -q 'Content structure' site/web/content/handbook-content-style.md; then ok "content-style has structure grammar"; else bad "content-style missing Content structure"; fi
if grep -q 'Canonical handbook outline' site/web/content/handbook-content-style.md; then ok "content-style has outline"; else bad "content-style missing outline"; fi

# Recipe criteria bridges
for f in site/web/templates/recipe-admin-resource.html site/web/templates/recipe-ops-queue.html site/web/templates/recipe-public-feed.html; do
  if grep -q 'recipe-criteria-bridge' "$f"; then ok "bridge in $f"; else bad "missing criteria bridge in $f"; fi
done

# Landing: only one primary button variant assignment in marketingLanding
primaries=$(grep -c 'Variant: "primary"' internal/app/landing.go || true)
if [[ "$primaries" -eq 1 ]]; then ok "landing.go has exactly one primary Variant"; else bad "landing.go primary Variant count=$primaries want 1"; fi

# Hub: Start here string appears before Deep dive in docsIndex source order
idx_start=$(grep -n 'Start here' internal/app/docs.go | head -1 | cut -d: -f1)
idx_deep=$(grep -n 'Deep dive' internal/app/docs.go | head -1 | cut -d: -f1)
if [[ -n "$idx_start" && -n "$idx_deep" && "$idx_start" -lt "$idx_deep" ]]; then ok "docsIndex Start here before Deep dive"; else bad "docsIndex ordering Start/Deep"; fi

# Hub names Core/System/Meta
if grep -q 'Core' internal/app/docs.go && grep -q 'System' internal/app/docs.go && grep -q 'Meta' internal/app/docs.go; then ok "docsIndex names Core/System/Meta"; else bad "docsIndex missing tier names"; fi

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
