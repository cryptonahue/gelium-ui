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

## What is next (ordered by value, owner reprioritizes freely)

```text
1.  Library/docs split   monorepo npm workspaces: lib/ (publishable CSS+tokens+js+html)
                         vs site/ (docs shell, landing, recipes). Site consumes lib by
                         the same contract as external consumers (dogfooding).
                         Unlocks: npm package, llms.txt, comparison page, Go module rename.
2.  NPM package           publish `gelium-ui` (name verified free) — CSS per component,
                         tokens/themes, consumer js, copyable HTML snippets
3.  Truth sync            README · this roadmap · theme registry · cmd/gelium
4.  DX / discoverability  llms.txt + agent prompts exposed in docs web surface
5.  Comparison page       "why Gelium vs Radix/shadcn/base-ui": server-rendered,
                         ~50KB JS vs 625KB, 0-JS contract, no-go decisions
6.  SEO productization    BASE_URL configurable + real og.png (when there is a domain)
7.  Recipes mobile        admin-resource (646px) / ops-queue (780px) expand the mobile
                         viewport — layout fix, not typography
8.  Theme polish          dark routine already single; scoped ownership of families
9.  Forms contract        document input rules BEFORE components: label above,
                         inputmode/type pairing, autocomplete, validate after interaction
10. Grid/container        -from-desktop modifiers, per-component reflow line (M3 swap
                         style), prose measure ≤75ch check
11. Performance stance    document ~50KB JS / 215KB CSS as a feature; Tailwind scope
                         review (CSS is the biggest asset)
12. Responsive chapter    "design for screen sizes, not devices" in served docs (ES)
13. Optional expansion    third theme · registry JSON runtime · iframe srcdoc demos ·
                         drawer menu upgrade · demo height presets · VT docs chapter ·
                         date-input pattern (3 inputs, inputmode numeric)
```

Items 13 are explicitly optional and wait for a real adoption signal before being
picked up — no date, no promise.

## How to read this page

Checkmarks mean contract-tested, not aspirational: every shipped phase carries tests that pin the behavior. If a phase regresses, the build fails. The "next" list is ordered by value, not by date — the owner reprioritizes freely.