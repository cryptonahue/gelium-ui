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

## Docs and DX

- Docs shell: Handbook (concept) before components (reference), search, theme `<select>`, light/dark switch, breadcrumbs, GitHub link
- Choose the right control: the cross-component decision table
- Guidance sections on every component page (when to use / when not to use / usability / accessibility)
- Content style guide: error / toast / empty-state / docs voice, plain English, ≤ 25-word sentences (contract-tested)
- Readability: 65ch measure, `text-wrap: pretty` / `balance`, hyphenation, WCAG AA contrast (contract-tested)
- Acknowledgments: every inspiration source credited (M3, USWDS, GOV.UK, Protocol, Base UI, Basecoat, Naive UI, Name That UI, Material Web, shadcn)

## What is next (post A–J)

```text
1. Truth sync          README · this roadmap · theme registry · cmd/gelium
2. DX / discoverability  demos first-class on component pages
3. SEO productization    BASE_URL configurable + real og.png (when there is a domain)
4. Theme polish          dark routine already single; scoped ownership of families
5. Release               version bump past 0.4.0 + tag
6. Optional expansion    more screen recipes · third theme · registry JSON runtime
```

## How to read this page

Checkmarks mean contract-tested, not aspirational: every shipped phase carries tests that pin the behavior. If a phase regresses, the build fails. The "next" list is ordered by value, not by date — the owner reprioritizes freely.
