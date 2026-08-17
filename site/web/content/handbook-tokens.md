# Tokens

The `--ui-*` vocabulary is the contract between components and themes: components consume tokens, themes supply their values, and nothing visual is hardcoded.

## What a token is

A token is a CSS custom property that names a design decision. Gelium UI classifies every `--ui-*` token by ownership:

- **Core tokens** — structural defaults owned by the core (`--ui-space-4`, `--ui-focus-thickness`, `--ui-state-hover-opacity`). Theme-neutral by design.
- **Theme tokens** — visual-direction values owned by a theme (`--ui-color-primary`, `--ui-radius-md`); another theme may override them.
- **Component tokens** — anatomy scoped to one component root (`--ui-dialog-scrim`, `--ui-card-radius`); a theme may override only the degrees of freedom the component contract declares (color, shape, type, motion).
- **Pattern tokens** — values shared by a composition pattern.
- **Internal tokens** — private implementation details, documented as out of the public contract.

## Core families

| Family | Prefix | Example | Owner |
|---|---|---|---|
| Semantic color | `--ui-color-*` | `--ui-color-primary`, `--ui-color-surface` | **theme** (roles in core; values per theme) |
| Typography | `--ui-font-*`, `--ui-type-*` | `--ui-font-sans`, `--ui-type-body-md` | **theme** for faces and scale steps; **core** owns shorthand composition aliases |
| Spacing | `--ui-space-*` | `--ui-space-4` | **core** (themes must not invent a parallel scale) |
| Radius | `--ui-radius-*` | `--ui-radius-md` | **theme** |
| Elevation | `--ui-shadow-*` | `--ui-shadow-2` | **theme** |
| Border | `--ui-border-*` | `--ui-border-width-1` | **core** defaults; theme may refine |
| Focus | `--ui-focus-*` | `--ui-focus-thickness` | **core** defaults; theme may refine ring color via `--ui-color-focus-ring` |
| Motion | `--ui-motion-*`, `--ui-easing-*` | `--ui-motion-short`, `--ui-easing-standard` | **core** defaults; **theme** may override timings |
| State | `--ui-state-*` | `--ui-state-hover-opacity` | **core** defaults; **theme** may override opacities |
| Z-index | `--ui-z-*` | `--ui-z-toast` | **core** (themes never own stacking order) |
| Breakpoints | `--ui-breakpoint-*` | `--ui-breakpoint-md` | **core** |
| Density and size | `--ui-density-*`, `--ui-size-*` | `--ui-size-control` | **core** defaults; **theme** may override control density |

### What a valid theme must define

A theme **MUST** define (light + single dark class route where color-bearing) at least:

- Semantic **color** family (`--ui-color-*` roles and `-fg` pairs used by components)
- **Typography** faces and decomposed type-scale steps the matrix requires
- **Radius** scale used by controls and surfaces
- **Elevation** (`--ui-shadow-*`) when components consume elevation tokens
- **Motion** overrides when the direction differs from core defaults (Material and Basecoat both set them)
- Every **component surface family** listed in the theme matrix (field, dialog, toast, card, badge, checkbox, radio, switch, slider, progress, fab, select, divider, …)

### What core never overrides

Core **never** re-owns or renames a theme’s brand direction under a second public prefix. Core does **not** ship product primary/secondary brand hexes as the long-term source of truth — those live on the theme root (`.theme-material`, `.theme-basecoat`). Core also does **not** override component anatomy that is scoped on the component root for structural layout (list item height, menu padding, data-table cell geometry) unless the component contract promotes a default into core for all themes to share.

Scoped component families (Owner = **component**): List, Menu, Data table, Navigation bar/tab/drawer, Segmented button, Tooltip — anatomy stays in the component CSS; themes may only override declared degrees of freedom (color, shape, type, motion).

## Naming conventions

`--ui-<family>-<token>` for transversal families and `--ui-<component>-<role>` for component coverage. Public classes are always `ui-*`; no third-party prefixes enter the contract.

## Rules

1. Everything visual exposed to themes uses `--ui-*` — no fixed values in components for color or control geometry.
2. Every reference resolves: no `var(--ui-*)` may point at an undefined token.
3. No dead tokens: every `--ui-*` has a consumer and an owner.
4. State layers paint with `color-mix()` over the defining `-fg` token, never hardcoded overlays.

Components that consume the shared vocabulary: [Elevation](/components/elevation), [Focus ring](/components/focus-ring), [Icon](/components/icon), and [Divider](/components/divider). Scoped families live on [List](/components/list), [Menu](/components/menu), [Data table](/components/data-table), and the navigation components. See [Themes](/docs/themes) for how themes populate token values, file paths (`lib/themes/*.css` / `gelium-ui/themes/*`), and the single dark class route.
