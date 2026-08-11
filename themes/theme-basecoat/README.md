# Gelium UI — Basecoat theme

> **Phase I** — token-only theme on the Loom contract. Same markup, same
> component API, same server contracts, different visual direction.
> Reference audit: `docs/handoffs/basecoat-audit.md`.

## 1. Visual direction

**Basecoat UI v1.0.2, style pack Vega** (the default) —
[basecoatui.com](https://basecoatui.com/) · [github.com/hunvreus/basecoat](https://github.com/hunvreus/basecoat) (MIT).

Basecoat is "all of the shadcn/ui magic, none of the React": a CSS component
library for Tailwind CSS 4 with a shadcn-compatible token set, where the style
pack *is* the visual direction. This theme is a **translation of the Vega style
pack into the `--ui-*` vocabulary** — no `basecoat-css` dependency, no markup
import, no JS. Vega is the neutral grayscale direction: near-black ink on
white, hairline borders, soft shadows, 0.625rem radius base, Geist typography,
and 2.25rem (h-9) control density.

## 2. Token mapping (Vega → `--ui-*`)

| Basecoat / shadcn (Vega) | Gelium `--ui-*` | Note |
|---|---|---|
| `--background` | `--ui-color-canvas` | |
| `--foreground` | `--ui-color-fg` | neutral-950 → neutral-50 in dark |
| `--muted` / `--muted-foreground` | `--ui-color-surface` / `--ui-color-fg-muted` | surface-container reuses muted |
| `--primary` / `--primary-foreground` | `--ui-color-primary` / `--ui-color-primary-fg` | 1:1; near-black ink (light) |
| `--secondary` / `--secondary-foreground` | `--ui-color-secondary` / `--ui-color-secondary-fg` | 1:1 |
| `--destructive` | `--ui-color-danger` / `--ui-color-danger-fg` | light stepped red-700 for AA solid surfaces (see §5) |
| `--border` | `--ui-color-border` | hairline neutral-200 / white-10% |
| `--input` | `--ui-field-border`, `--ui-select-outline` | |
| `--ring` | `--ui-color-focus-ring` | AA step: neutral-500 light / neutral-400 dark |
| `--card` / `--popover` | `--ui-card-container-*`, `--ui-dialog-container`, `--ui-toast-container` | |
| `--radius` (0.625rem) | `--ui-radius-*` | base lands on `--ui-radius-md`; xs=0.375rem (inputs), lg=0.75rem (dialogs/cards) |
| `--font-sans` / `--font-mono` | `--ui-font-sans` / `--ui-font-mono` | Geist Sans/Mono, system fallbacks (not bundled) |
| shadow-xs/sm/md/lg/xl (Tailwind) | `--ui-shadow-1..5` | derived from the utilities Vega hardcodes |
| duration-100/200/300 + ease | `--ui-motion-{short,medium,long}` + `--ui-easing-standard` | Tailwind default ease |
| hover:bg-primary/80, opacity-50 | `--ui-state-{hover,pressed}-opacity` = .20, `--ui-state-disabled-opacity` = .50 | approximates the color-swap hover model |

Out of scope of the mapping (no Phase I consumer): `--popover-*` accents,
`--chart-*`, `--sidebar-*`, icon data-URIs.

## 3. Light / dark

Self-contained in `theme.css`: `.theme-basecoat` (light),
`.theme-basecoat.theme-dark/.dark/[data-theme="dark"]` (explicit dark class) and
`@media (prefers-color-scheme: dark)` (dark media) — a single routine per
scheme, no duplication with drift. Dark colors are the Vega dark set (neutral
900 surfaces, light ink, white-10% borders).

## 4. Variants and states covered

Component CSS is theme-agnostic; the theme supplies the tokens. Covered states
per component (same contract as the Material theme, verified by
`TestThemeMatrixCoversEveryAvailableTheme`):

- **Button**: primary/secondary/outline/text variants; hover, focus-visible,
  pressed, disabled, loading spinner.
- **Text field**: filled/outlined; hover, focus, error, disabled, empty
  (floating label).
- **Card**: elevated/filled/outlined; focus-visible.
- **Badge**: dot/large + error/success/warning/info tones (tone surfaces derive
  from semantic colors, so dark legibility follows the dark palette).
- **Dialog**: modal + page variant; open entrance, backdrop.
- **Toast**: info/success/warning/error icon roles, action, show transition.
- **Data table**: hover/focus/pressed row layers, selected rows, sort, disabled
  pagination, checkbox.

## 5. Divergencias (component × divergence × decision)

| Componente | Divergencia | Decisión |
|---|---|---|
| **Text field** | Basecoat uses a static label on a 2.25rem input; Gelium keeps the floating label (filled/outlined, `:placeholder-shown`) | Keep the floating label with Basecoat aesthetics; `--ui-size-field` steps to 3rem so the label keeps headroom (documented divergence from Basecoat's 2.25rem input) |
| **Toast** | Basecoat Toast needs JS (`toaster` API, auto-dismiss); Gelium demands no-JS end-to-end (`aria-live` + `loom:toast` + inline fallback) | Keep the Gelium server-driven contract untouched; take only the aesthetic: popover surface (light container, hairline border, shadow-lg, rounded-2xl) |
| **Danger solid vs /10 soft** | Basecoat renders `--destructive` as text on a /10 soft fill; Gelium paints `--ui-color-danger` as a solid badge/tone surface | Light `--ui-color-danger` steps red-600 → red-700 (#b91c1c, 5.86:1 with white) to satisfy WCAG AA body text on the solid surface; dark keeps Vega's red-500 with near-black on-color |
| **Focus ring opacity** | Vega ships `--ring` used at 50% opacity; Gelium's focus model is a solid outline | `--ui-color-focus-ring` uses the neutral ring family at an AA non-text step (light neutral-500, dark neutral-400) |
| **Status colors** | Vega tokenizes only destructive; no success/warning/info source | Core semantic palette kept, dark variants tuned for AA contrast (documented in `theme.css`) |
| **State feedback** | Basecoat hovers swap the token color (bg-primary/80, hover:bg-muted); Gelium paints a `color-mix` state layer | `--ui-state-hover/pressed-opacity` raised to .20 to approximate the visible feedback; disabled at .50 per Basecoat's opacity-50 |
| **Badge anatomy** | Basecoat badge is a text pill with data-variant; Gelium badge is a dot/count error-tinted marker (Material) | Dot/count + large label forms kept (Gelium contract); pill-identity variants (outline/ghost/link) are out of Phase I scope |
| **Button/Badge variants** | Basecoat adds ghost/link/destructive and pill variants Gelium lacks | Out of scope: extending variants is a core decision, never theme CSS (contract §7) |
| **FAB / checkbox / radio / switch / slider / progress / select** | No Basecoat equivalent (FAB) or different anatomy | FAB derives from primary/secondary surfaces; form controls use Basecoat's native-control values (size-4 check/radio, 32×18.4 switch, 6px tracks, 2.25rem select) |
| **Dark scrim** | Basecoat backdrop is `bg-black/10` + `backdrop-blur-xs`; Gelium's dialog core renders no blur | Scrim steps to rgb(0 0 0 / .20) light / rgb(0 0 0 / .40) dark so the modal keeps separation without blur |

## 6. Component coverage (Phase I scope)

| Component | Light | Dark | Token families | Notes |
|---|---|---|---|---|
| Button | ✅ | ✅ | `--ui-color-*`, `--ui-radius-full`, `--ui-shadow-*`, `--ui-state-*`, `--ui-focus-*`, `--ui-type-label-lg` | control density 2.25rem |
| Text field | ✅ | ✅ | `--ui-field-*`, `--ui-size-field` | floating label kept (3rem field) |
| Card | ✅ | ✅ | `--ui-card-*`, `--ui-radius-lg`, `--ui-shadow-1` | elevated/outlined = canvas, filled = muted |
| Badge | ✅ | ✅ | `--ui-badge-*` + semantic colors | error pair + tones |
| Dialog | ✅ | ✅ | `--ui-dialog-*`, `--ui-radius-lg` | popover surface |
| Toast | ✅ | ✅ | `--ui-toast-*` incl. `icon-*` | popover surface, no-JS contract |
| Data table | ✅ | ✅ | semantic colors + core size/state/type tokens | anatomy stays scoped in `data-table.css` |

The data table reads the semantic palette (surface container, border, muted
labels, secondary selection tint, primary accent) plus the core size/state/type
scale; its scoped anatomy tokens (row heights, cell padding) stay in the core,
as the contract §3.3 decides for scoped components.

## 7. Verification

```bash
npm run build          # bundles the theme into web/static/app.css
go test ./...          # matrix discovers theme-basecoat by glob, no test edits
go vet ./...
```

The theme entered the bundle in the two Phase H steps: the explicit
`@import "../../themes/theme-basecoat/theme.css"` in `web/styles/app.css` and
the `"theme-basecoat"` entry in `themeClass()` (`internal/app/server.go`).
`TestThemeSelectionIsClassDrivenWithoutJS` still proves selection is a class
swap with zero JS.
