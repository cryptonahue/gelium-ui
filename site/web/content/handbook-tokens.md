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

| Family | Prefix | Example |
|---|---|---|
| Semantic color | `--ui-color-*` | `--ui-color-primary`, `--ui-color-surface` |
| Typography | `--ui-font-*`, `--ui-type-*` | `--ui-font-sans`, `--ui-type-body-md` |
| Spacing | `--ui-space-*` | `--ui-space-4` |
| Radius | `--ui-radius-*` | `--ui-radius-md` |
| Elevation | `--ui-shadow-*` | `--ui-shadow-2` |
| Border | `--ui-border-*` | `--ui-border-width-1` |
| Focus | `--ui-focus-*` | `--ui-focus-thickness` |
| Motion | `--ui-motion-*`, `--ui-easing-*` | `--ui-motion-short`, `--ui-easing-standard` |
| State | `--ui-state-*` | `--ui-state-hover-opacity` |
| Z-index | `--ui-z-*` | `--ui-z-toast` |
| Breakpoints | `--ui-breakpoint-*` | `--ui-breakpoint-md` |
| Density and size | `--ui-density-*`, `--ui-size-*` | `--ui-size-control` |

## Naming conventions

`--ui-<family>-<token>` for transversal families and `--ui-<component>-<role>` for component coverage. Public classes are always `ui-*`; no third-party prefixes enter the contract.

## Rules

1. Everything visual exposed to themes uses `--ui-*` — no fixed values in components for color or control geometry.
2. Every reference resolves: no `var(--ui-*)` may point at an undefined token.
3. No dead tokens: every `--ui-*` has a consumer and an owner.
4. State layers paint with `color-mix()` over the defining `-fg` token, never hardcoded overlays.

Components that consume the shared vocabulary: [Elevation](/components/elevation), [Focus ring](/components/focus-ring), [Icon](/components/icon), and [Divider](/components/divider). Scoped families live on [List](/components/list), [Menu](/components/menu), [Data table](/components/data-table), and the navigation components. See [Themes](/docs/themes) for how themes populate token values.
