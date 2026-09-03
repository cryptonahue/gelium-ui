# Roadmap

This page tracks where Gelium UI is and where it is going. The internal system roadmap (`docs/gelium-ui-system-roadmap.md`) is the source of truth; this page mirrors the public status.

## What is done (phases A–J)

Gelium UI shipped its full 10-phase system roadmap: verification and close-out of every phase with contract tests, not promises.

| Phase | Deliverable | Status |
|---|---|---|
| A | Core Go module, server-rendered components, 0-JS contract | ✅ Shipped |
| B | `--ui-*` token vocabulary, typography decomposition | ✅ Shipped |
| C | Theme-agnostic verification (new theme = glob, no test edits) | ✅ Shipped |
| D | Screen state patterns (empty, error, loading, success, toast…) | ✅ Shipped |
| E | SEO/GEO server-driven metadata (dates, JSON-LD, BASE_URL) | ✅ Shipped |
| F | Public content patterns (14 patterns, card slots) | ✅ Shipped |
| G | Screen recipes (Admin Resource, Ops Queue, Public Feed) | ✅ Shipped |
| H | Theme mechanism (class on `<html>`, dark via single class route) | ✅ Shipped |
| I | Basecoat theme (light + dark, same bundle) | ✅ Shipped |
| J | Registries + sync guards (docs cannot drift from code) | ✅ Shipped |

## What is done (v0.5.0+)

| Deliverable | Status |
|---|---|
| HTMX 4 migration (namespaced events, innerMorph, optimistic theme/scheme with server authority) | ✅ Shipped |
| Docs UX: on-this-page rail, prev/next, optimistic dark toggle, themed recipes | ✅ Shipped |
| Material Symbols icons (curated 21 glyphs, embedded server SVGs) | ✅ Shipped |
| Docs coverage: contributing page, browser-support table, MIT named, component inventory via `/docs` | ✅ Shipped |
| Demo-first component pages (title → live demo → body) | ✅ Shipped |
| Release v0.5.0 (tag + changelog + version coherence tests) | ✅ Shipped |
| Chrome href refresh after optimistic toggle (dark mode survives sidebar navigation) | ✅ Shipped |
| Mobile audit — foundations: `--ui-touch-target`, `--ui-container-max`, `prefers-reduced-motion`, hero 320px overflow fix | ✅ Shipped |
| Mobile audit — overflow containment: preview, data-table, chips, whatsapp demo (no masking; min-width:0 + internal scroll) | ✅ Shipped |
| Mobile audit — runtime: same-document view transitions (reduced-motion guard), safe areas, mobile nav 100dvh + aria-label, focus matrix | ✅ Shipped |
| Library/docs split — monorepo npm workspaces: `lib/` publishable package `gelium-ui` (CSS+tokens+themes+templates+consumer JS+dist) + `site/` docs consumer dogfooding by package name; dual go:embed + gates include `./lib/...` | ✅ Shipped |
| npm package `gelium-ui@0.6.2` on the public registry (https://www.npmjs.com/package/gelium-ui) — dist bundle, themes, templates, consumer JS | ✅ Shipped |
| `/llms.txt` agent brief (install, wire contracts, when not to use, component list) | ✅ Shipped |
| Comparison page `/docs/compare` — Why Gelium vs Radix/shadcn/Base UI, ~50KB JS story, no-gos | ✅ Shipped |
| Recipes mobile — admin/ops/feed headers stack + table/list containment under ~40rem (no overflow-x:hidden) | ✅ Shipped |
| Forms contract `/docs/forms` — label above, inputmode/type, autocomplete, validate after interaction | ✅ Shipped |
| Performance stance `/docs/performance` — ~50KB JS / CSS-as-largest-asset as product stance | ✅ Shipped |
| Theme polish — family ownership matrix, lib/themes paths, dark class-only docs | ✅ Shipped |
| Layout utilities — `.ui-container`, `.ui-row-from-desktop`, prose measure ≤75ch (65ch) | ✅ Shipped |
| Responsive chapter `/docs/responsive` — design for screen sizes, not devices | ✅ Shipped |
| Docs hub `/docs` — orientation Start here (no sidebar catalog dump) | ✅ Shipped |
| UX criteria `/docs/screens` + `/docs/feedback` (GOV.UK/USWDS/M3/NNG) + `/llms-ux.txt` agent pack | ✅ Shipped |
| Journeys + data display + patterns skeletons + density/motion + UI DoD checklist | ✅ Shipped |
| Agent workflow (Impeccable-class process adapted) + surface modes + ux-detect + PRODUCT/DESIGN templates | ✅ Shipped |
| Docs sidebar handbook tiers Core / System / Meta (scannable IA) | ✅ Shipped |
| Content structure grammar (H1–H3/list/table) + recipe criteria bridges + harder ux-detect | ✅ Shipped |

## What is next (ordered by value, owner reprioritizes freely)

```text
0.  Icon pack gallery      top icon SVG sets (Tabler, Lucide, Heroicons,
                           Phosphor…) as opt-in `lib/icons/<pack>` packs,
                           each with a searchable Tabler-style gallery on
                           /docs/icon (clipboard copy, themeable currentColor)
1.  SEO productization    BASE_URL configurable + real og.png (when there is a domain)
2.  Optional expansion    third theme · registry JSON runtime · iframe srcdoc demos ·
                          drawer menu upgrade · demo height presets · VT docs chapter ·
                          date-input pattern (3 inputs, inputmode numeric)
3.  Docs localization     translate the internal spec docs (docs/gelium-ui-*.md,
                          docs/handoffs/*) from Spanish to English. Blocked on the
                          docs-spec split first: today the public site (/docs/*) and
                          the package guidance (AGENTS.md, skills/) are English, but the
                          internal specs they reference are Spanish, and contract tests
                          currently pin Spanish anchors ("## 3. Patrones de estado",
                          "Nombres alternativos", "Prompt para agentes"). Sequence:
                          (a) split public vs internal docs (move internal specs to
                          docs/internal/), (b) translate internal specs + their
                          test anchors to English together so the build stays green.
4.  Bun toolchain          swap npm for Bun across the monorepo: `bun install`
                          (fast, replaces package-lock.json with bun.lock),
                          `bun run build` in root/site (tailwindcss CLI + .mjs
                          scripts run natively), and `bun publish`/`bun pack` in
                          scripts/publish-lib.sh. Bun 1.3.11 already on the
                          dev machine. Tradeoffs to confirm before committing:
                          CI/devs still on npm need the lockfile story; keep the
                          Go bake/explore/embed untouched (build of Go binary is
                          independent of JS toolchain). Do the switch on its own
                          PR/worktree, not mixed with feature work.
```

Items 0–4 are ordered by value, not by date, and the owner reprioritizes freely. The
icon pack gallery (item 0) is a product feature: it turns the 7.8k-glyph Material
Symbols source from a 38-icon curated allowlist into discoverable, consistent packs —
not a 10-set visual casserole. Packs stay SVG server-rendered (themeable via
`currentColor`), never icon-fonts. Docs localization (item 3) is maintenance debt, not
a feature: the public product is already fully English, but the internal specs it
references are Spanish, so the translation work is gated behind a docs public/internal
split and done together with the test anchors that currently pin the Spanish.

## How to read this page

Checkmarks mean contract-tested, not aspirational: every shipped phase carries tests that pin the behavior. If a phase regresses, the build fails. The "next" list is ordered by value, not by date — the owner reprioritizes freely.
