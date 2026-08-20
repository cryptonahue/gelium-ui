# Gelium UI — Alden theme

> **Font-pack reference** — the first theme to ship self-hosted webfonts end to
> end (Inter + Source Serif 4, WOFF2, subsetted, metric-adjusted fallbacks) on the
> Gelium contract. Same markup, same component API, same server contracts,
> different visual direction. Reference: Refero style "Alden"
> (`styles.refero.design`). Font pipeline: `docs/gelium-ui-font-contract.md`.

## 1. Visual direction

**"Serene clinic on warm parchment"** — an editorial healthcare register. The
page is 98% achromatic: a near-white canvas warmed by a single cream parchment
surface, deep near-black typography, and exactly **two chromatic accents** used
like punctuation, not decoration:

- **Sky Highlight** (`#97cde5`) — a whisper-soft blue used to recolour the
  chosen word inside a heading (inline annotation).
- **Sage Action** (`#c8dfaa`) — the muted green reserved **only** for the
  primary action / CTA.

Type mixes a geometric sans (**Inter**) for UI with a contemporary serif
(**Source Serif 4**) reserved for the largest editorial headline statements
(display/hero — never UI). Surfaces are flat with very soft card edges (16px),
fully pill buttons, and calm motion. Closer to a premium health magazine than a
SaaS dashboard.

## 2. Fonts (self-hosted, this is the reference)

The reference font (Stk bureau Sans/Serif Book Trial) is commercial, so the
theme uses the **substitutes the reference itself declares**: **Inter** (sans)
and **Source Serif 4** (display serif). Both are SIL OFL.

| Family | Weights | Subset | File (`lib/fonts/`) | Preload |
|---|---|---|---|---|
| Inter | 400 | latin | `theme-alden-inter-400-latin.woff2` | ✅ |
| Inter | 400 | latin-ext | `theme-alden-inter-400-latin-ext.woff2` | ✅ |
| Inter | 500 | latin | `theme-alden-inter-500-latin.woff2` | on demand |
| Inter | 500 | latin-ext | `theme-alden-inter-500-latin-ext.woff2` | on demand |
| Inter | 600 | latin | `theme-alden-inter-600-latin.woff2` | on demand |
| Inter | 600 | latin-ext | `theme-alden-inter-600-latin-ext.woff2` | on demand |
| Source Serif 4 | 400 | latin | `theme-alden-source-serif-4-400-latin.woff2` | on demand |
| Source Serif 4 | 400 | latin-ext | `theme-alden-source-serif-4-400-latin-ext.woff2` | on demand |

Per `gelium-ui-font-contract.md`:

- **Self-hosted, WOFF2, subsetted** to latin + latin-ext (AR/BR/ES coverage).
- **`font-display: swap`** in every `@font-face` (text visible from the start).
- **Metric-adjusted fallbacks** (`--ui-font-sans`/`--ui-font-display` tail with
  `"Alden Sans Fallback"` / `"Alden Display Fallback"` `@font-face` blocks that
  pin the ascent/descent/line-gap to the real font's metrics, so the swap
  causes no layout shift / CLS). Inter 400 size-adjust ≈ 99%.
- **Preload capped**: only the body weight (Inter 400 latin + latin-ext) is
  preloaded in `<head>`; heavier weights and the display serif load on demand
  — they're servable but not preloaded, so preloads don't compete with the LCP
  resource (contract §3.A.4 / §B.7).
- **Licenses**: Inter and Source Serif 4 are SIL OFL — free for commercial
  use. See `docs/gelium-ui-font-contract.md` §E; the filenames ship under
  `lib/fonts/`.

The `@font-face` rules live in `theme.css` alongside the tokens and are bundled
through `lib.Assets` (`//go:embed fonts/*`), served at `/static/fonts/*`.

## 3. Token mapping (Alden palette → `--ui-*`)

| Alden | Gelium `--ui-*` | Note |
|---|---|---|
| Paper White `#ffffff` | `--ui-color-canvas` | dominant canvas & cards |
| Parchment Cream `#f3f1eb` | `--ui-color-surface` / `--ui-color-surface-container` | second, softer layer |
| Ink Black `#28262a` | `--ui-color-fg` | faint plum undertone, softer than #000 |
| Graphite `#4a4a4c` | `--ui-color-fg-muted` | secondary text / captions |
| Sage Action `#c8dfaa` | `--ui-color-primary` / `--ui-color-primary-fg` | the only saturated CTA fill |
| Sky Highlight `#97cde5` | `--ui-color-secondary` | in-line word annotation |
| Fog Border `#dddcdd` | `--ui-color-border` | hairline dividers / badges |
| Silver `#cbcbcb` | `--ui-color-border-strong` | mid dividers |
| Radius system | `--ui-radius-*` | smallUI 4, cards 16, images 24, badges 30, buttons 100 |
| Inter | `--ui-font-sans` + `--ui-type-*-family` | workhorse sans |
| Source Serif 4 | `--ui-font-display` + display steps | hero headline only |

`--ui-font-display` is Alden's new optional token: the serif reserved for the
`display-lg` / `display-sm` typescale steps. Every other step stays on the
sans. `--ui-font-display` has a system serif fallback so an unloaded font still
renders.

## 4. Light / dark

Self-contained in `theme-alden.css`: `.theme-alden` (light),
`.theme-alden.theme-dark/.dark/[data-theme="dark"]` (explicit dark class) — a
**single dark class route**, no `@media (prefers-color-scheme: dark)`. Dark
inverts the register: deep warm ink canvas (`#1a191c`), parchment becomes the
elevated surface, and the two chromatic accents survive unchanged as
punctuation. The font families are scheme-independent (fonts are only declared
in light; dark resolves identical — enforced by the type snapshot test).

## 5. Variants and states covered

Component CSS is theme-agnostic; the theme supplies the tokens. Covered states
per component (verified by `TestThemeMatrixCoversEveryAvailableTheme` which
discovers `theme-alden` by glob with **zero test edits**):

- **Button**: primary/secondary/outline/text variants; hover, focus-visible,
  pressed, disabled, loading spinner. Primary = sage pill (Ink text).
- **Text field**: floating label, hairline border on paper; hover, focus,
  error, disabled, empty.
- **Card**: elevated/filled/outlined; focus-visible. Flat faces on parchment.
- **Badge**: dot/large + error/success/warning/info tones (semantic colors).
- **Dialog / Toast**: soft parchment popovers, flat; no-JS server contract.
- **Data table**: semantic warm palette + core size/state/type tokens.
- **Checkbox/radio/switch/slider/progress/select/fab/divider**: full analogue
  families (matrix coverage), colors derived from the warm palette.

## 6. Implementing / verification

```bash
npm run build          # bundles theme-alden + imported fonts into static/app.css
go test ./...          # matrix + font packs discover theme-alden by glob
go vet ./...
```

The theme entered the bundle in the two Phase H steps: the explicit
`@import "gelium-ui/themes/theme-alden.css"` in `site/web/styles/app.css` and
the `theme-alden` entry (with its `Fonts`) in `availableThemes`
(`internal/app/server.go`). Preload `<link>`s are emitted automatically by
`layout.html` (`themePreloadFonts`); the `/static/fonts/*` route serves only
allowlisted theme fonts (closed namespace). `internal/app/fonts_test.go`
proves the real Alden fonts serve with `font/woff2` and that the layout emits
preloads for Alden but none for Material.

## 7. Divergencias (component × divergence × decision)

| Componente | Divergencia | Decisión |
|---|---|---|
| **Display serif in UI** | Alden's reference reserves the serif for the biggest headlines only ("never use it in UI") | `--ui-font-display` maps only the `display-*` steps; all body/label/caption stay on Inter. Enforced by the typescale mapping, not theme CSS. |
| **Primary = sage** | A single saturated fill | Matches the reference; the sage-tinted primary is used sparingly (CTAs, selected controls), everything else stays achromatic. |
| **Flat cards** | Alden uses cream/warm surfaces with no elevation | `--ui-shadow-1..5` are near-flat; card separation is by color (paper vs parchment), not shadow. |
| **Button radius** | Fully pill (100px) | `--ui-radius-full: 100px` for buttons; cards keep 16px. Matches Alden's radius system. |
| **Fonts** | Alden ships Stk bureau fonts (commercial) | Substituted with the reference's own declared replacements (Inter / Source Serif 4), self-hosted + subsetted + metric-adjusted per font contract. |

## 8. Font quality (contract §3) — checklist

- [x] `.woff2` self-hosted, no external CDN
- [x] Subsetted latin + latin-ext; only needed weights (400/500/600 sans, 400 serif)
- [x] `font-display: swap` on every `@font-face`
- [x] Metric-adjusted fallbacks (anti-CLS) for sans + display
- [x] Preload capped to the body weight (2 files)
- [x] License documented (Inter, Source Serif 4 — SIL OFL)
- [x] Body ≥ 16px; tracking not negative at body sizes; AA contrast via palette
