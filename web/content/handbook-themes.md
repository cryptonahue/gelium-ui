# Themes

Themes are visual directions over one markup surface: a theme changes `--ui-*` token values, never HTML, component anatomy, or server contracts.

## The layer model

Gelium UI separates *structure* from *visual direction* with six ownership layers: core tokens, themes, components, patterns, recipes, and screens. A theme sits in the second layer — it maps the same `--ui-*` roles onto a visual direction (color, shape, type, motion). Because markup and contracts never move between layers, switching direction is a single class change with no rebuild and no JavaScript. See [Documentation](/docs) for the full model.

## Material (default)

Material is the default direction, built on **Material 3 (M3)**: role-based color (`primary`, `secondary`, `tertiary`, `surface`, `error` and their foregrounds), the M3 type scale mapped onto `--ui-type-*` steps, state layers via `--ui-state-*` opacity tokens, and M3 elevation and shape. Gelium implements M3 as token values in `themes/theme-material/theme.css` — components never hardcode a Material look.

## Basecoat

Basecoat is the alternative direction: the Basecoat UI "Vega" style pack translated into the `--ui-*` vocabulary — near-black ink on white, hairline borders, soft shadows, a 0.625rem base radius, Geist typography, and 2.25rem control density. It ships in the same bundle as Material, selected the same way.

## Selecting a theme

- **Query route** — append `?theme=<slug>` to any docs or component URL (`?theme=material`, `?theme=basecoat`). Only allowlisted slugs apply; unknown values keep the default direction.
- **Class route** — set the direction on the document root directly: `<html class="theme-material">` or `<html class="theme-basecoat">`.
- **Docs topbar** — the Theme control in the docs topbar rewrites the current URL's query, and in-shell navigation preserves your selection.

## Dark mode

Dark is an explicit class route, not a media-query guess: `?scheme=dark` or `<html class="theme-material theme-dark">`. Every theme ships light and dark in the same bundle; the document root decides.

## What a theme never changes

Theme switching must not change URLs, landmarks, SEO metadata, markup, or server contracts — only the root class and cascade tokens. [Button](/components/button), [Dialog](/components/dialog), and [Toast](/components/toast) render identically under every direction; the [Tokens](/docs/tokens) handbook page defines what can vary and who owns it.
