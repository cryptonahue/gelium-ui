# Skill: Foundations — tokens, themes, layout

Load this first. Everything below is the non-negotiable base for any Gelium UI.

## Theme by class, never hex

Top-level theme lives on the document root (`<html>`), not in component markup:

```html
<html class="theme-material">            <!-- light -->
<html class="theme-material theme-dark" data-theme="dark">  <!-- dark -->
```

`theme-material` (Material 3, default) and `theme-basecoat` ship in the package.
Light/dark is a **class route** — there is no media-only dark toggle. Pick themes
you need at install; both are in `dist/gelium.css`.

## Use `--ui-*` tokens, no literals

Every visual value in shipped content is a token. Never invent a one-off color,
radius, shadow, spacing, or font in a page — map it to a token from the theme.

Key families: `--ui-color-*`, `--ui-type-*`, `--ui-space-*`, `--ui-radius-*`,
`--ui-shadow-*`, `--ui-size-*`, `--ui-touch-target`.

- **Same typographic role for the same job**: one title style per heading level,
  drawn from the `--ui-type-*` scale — never introduce new font families or
  one-off sizes per feature.

## 0-JS default

The primary flow of every screen must work with JavaScript disabled. Native
elements and server round-trips first. JS (`js/gelium.js`) only enhances: 422
validation swap, `gelium:toast` regions, view-transition guard, slider fill.

## Layout utilities

- `.ui-container` — centered constrained content column.
- `.ui-row-from-desktop` — row layout from desktop breakpoint.
- Prose measure ≤ 75ch (use 65ch in docs), `text-wrap: pretty`.

### Spacing contract

All spacing uses `--ui-space-*` tokens (4px base scale). Pick the token from
the **semantic relationship** between elements — never invent a value, never
eyeball it:

| Relationship | Token | Rationale |
|---|---|---|
| Element → its own metadata (title → summary, label → hint) | `--ui-space-1` / `--ui-space-2` | same unit; proximity binds them |
| Item → item inside one list/feed | `--ui-space-3` | siblings of the same group |
| Group → group within a section | `--ui-space-4` / `--ui-space-6` | Gestalt: related groups sit closer than unrelated ones |
| Section → section (page rhythm) | `--ui-space-8` | major page beats |

Rules:

- **Proximity = grouping**: children closer to each other than sibling groups.
  If two blocks read as "glued", you used space-0/1 where the relationship is
  space-3+.
- **One axis per relationship**: pick vertical OR horizontal spacing for a
  pairing, not both.
- Components own their internal padding (`--ui-card-padding`, etc.) — compose
  with margins/gaps *between* components only.
- Never hardcode px/rem in markup or consumer CSS; if no token fits the
  relationship, the scale is wrong — escalate instead of inventing.

Prior art: Material 3 spacing/density rules; the industry-standard 8-point
(4px half-step) grid; classic typographic vertical rhythm.

### Registry-first page shells

Page-level layouts must compose registered components (`ui-container`,
`ui-navigation-bar`, drawer primitives) — never hand-roll a nav header or sticky
page shell. Custom shell CSS is limited to spacing and width. If you need a shell
the registry does not cover, extend the component, not the page.

## Mobile guardrails

- Touch targets ≥ 44px (`--ui-touch-target`).
- Scroll containers use `min-width: 0` + internal overflow, never
  `overflow-x: hidden` on `body` (masking is forbidden — it hides content loss).

## No-JS first check

Is the task's main flow (view, form submit, filter, delete) achievable without
JS? If not, fix the server/HTML path before adding enhancement.
