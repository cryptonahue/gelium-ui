# Browser support

Gelium UI is a server-rendered library: the baseline is modern HTML, and platform features are applied as progressive enhancements inside `@supports` so older browsers degrade to a usable native control. There is no single minimum version — the meaningful question is which platform APIs each feature relies on. This page consolidates the per-component Baseline notes.

## Platform APIs in use

| API | Baseline status | Where it is used |
|---|---|---|
| Popover API | Baseline 2025 (newly available) | Menu — zero-JS popover positioning |
| CSS anchor positioning | Baseline 2026 | Menu positioning enhancement inside `@supports` |
| Invoker Commands | Baseline 2025 | Dialog/menu light-dismiss triggers |
| `<dialog>` closedby | not Baseline yet | Dialog light-dismiss when available; explicit close control otherwise |

The library never requires these APIs to function: each one is gated by `@supports` or by a documented fallback, so controls keep their native semantics and keyboard behavior on older browsers.

## What always works

- Native HTML semantics: `input`, `button`, `dialog`, `select`, `menu` — no component JavaScript to bundle.
- Server round-trips with HTMX as a progressive enhancement; without JavaScript the same forms submit classically.
- WCAG AA contrast in light and dark, forced-colors borders, reduced-motion, and RTL — enforced by tests, not browser-dependent.

## Per-component notes

See each component page's "Compatibility" section for the exact Baseline claims: [Menu](/components/menu), [Tooltip](/components/tooltip), [Segmented buttons](/components/segmented-button), [Dialog](/components/dialog), [Navigation drawer](/components/navigation-drawer).
