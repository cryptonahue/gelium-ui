# Themes

> Live specimens: a real kitchen-sink of every theme (buttons, fields, cards,
> badges, toasts, dialog, form controls) rendered under the active direction.
> → **[Theme gallery](/docs/themes/gallery)**

Themes are visual directions over one markup surface: a theme changes `--ui-*` token values, never HTML, component anatomy, or server contracts.

## The layer model

Gelium UI separates *structure* from *visual direction* with six ownership layers: core tokens, themes, components, patterns, recipes, and screens. A theme sits in the second layer — it maps the same `--ui-*` roles onto a visual direction (color, shape, type, motion). Because markup and contracts never move between layers, switching direction is a single class change with no rebuild and no JavaScript. See [Documentation](/docs) for the full model and [Tokens](/docs/tokens) for the full family ownership matrix.

## Token-family ownership

Themes own **visual direction values**. Core owns structural defaults. Components own anatomy scoped to their root. The split is intentional and enforced by contract tests.

| Layer | Owns | Does not own |
|---|---|---|
| **Core** | Structural defaults: spacing scale, breakpoints, z-index ladder, focus geometry defaults, motion/easing defaults, type *composition* aliases, neutral fallbacks | Product brand color, typeface choice, radius personality |
| **Theme** | Direction values for color, typeface and type scale steps, radius, elevation (shadows), motion timing when the direction needs it, state opacities, and every component surface family the matrix requires | Markup, class names, HTMX contracts, component anatomy (heights, paddings that stay structural) |
| **Component** | Scoped anatomy tokens on the component root (`--ui-list-*`, `--ui-menu-*`, dialog/toast/card structure defaults) | Replacing the public `--ui-*` vocabulary with third-party prefixes |

**Theme-owned families (must be defined for a valid theme):** semantic color (`--ui-color-*`), typography (`--ui-font-*` + decomposed `--ui-type-*-{size,weight,…}`), radius, elevation (`--ui-shadow-*`), and the component surface families in the theme matrix (field, dialog, toast, card, badge, fab, select, switch, slider, …). See [Tokens](/docs/tokens) for the Owner column and the “themes must define” list.

**Core-owned families themes may refine but must not invent a parallel scale for:** spacing (`--ui-space-*`), breakpoints, z-index. Core never redefines a theme’s brand roles under a different name — themes override values under the same `--ui-*` roles.

## Material (default)

Material is the default direction, built on **Material 3 (M3)**: role-based color (`primary`, `secondary`, `tertiary`, `surface`, `error` and their foregrounds), the M3 type scale mapped onto `--ui-type-*` steps, state layers via `--ui-state-*` opacity tokens, and M3 elevation and shape. Gelium implements M3 as token values in `lib/themes/theme-material.css` — components never hardcode a Material look.

## Basecoat

Basecoat is the alternative direction: the Basecoat UI "Vega" style pack translated into the `--ui-*` vocabulary — near-black ink on white, hairline borders, soft shadows, a 0.625rem base radius, Geist typography, and 2.25rem control density. It ships in the same bundle as Material, selected the same way (`lib/themes/theme-basecoat.css`).

## Where theme files live

Themes ship **flat** under the library package — one CSS file per direction, not a nested folder:

| Location | Path |
|---|---|
| Repository | `lib/themes/<name>.css` (e.g. `lib/themes/theme-material.css`) |
| npm package export | `gelium-ui/themes/<name>.css` (e.g. `@import "gelium-ui/themes/theme-material.css"`) |

The site entry imports every theme explicitly (`@import "gelium-ui/themes/….css"`). CSS does not glob: a theme is in the bundle if and only if that import line exists. The old nested layout `themes/<name>/theme.css` is retired — do not document or create it.

## Selecting a theme

- **Query route** — append `?theme=<slug>` to any docs or component URL (`?theme=material`, `?theme=basecoat`). Only allowlisted slugs apply; unknown values keep the default direction.
- **Class route** — set the direction on the document root directly: `<html class="theme-material">` or `<html class="theme-basecoat">`.
- **Docs topbar** — the Theme control in the docs topbar rewrites the current URL's query, and in-shell navigation preserves your selection.

## Dark mode

Dark is an **explicit class route only** — never a media-only path and never `@media (prefers-color-scheme: dark)` as the theme’s dark mechanism:

- Query: `?scheme=dark`
- Class: `<html class="theme-material theme-dark">` (or `.dark` / `[data-theme="dark"]` aliases on the same theme root)

Every theme ships light and dark in the same file; the document root decides. There is a **single dark class route** per theme. Prefer-color-scheme media must not define a second, drifting palette.

## What a theme never changes

Theme switching must not change URLs, landmarks, SEO metadata, markup, or server contracts — only the root class and cascade tokens. [Button](/components/button), [Dialog](/components/dialog), and [Toast](/components/toast) render identically under every direction; the [Tokens](/docs/tokens) handbook page defines what can vary and who owns it.
