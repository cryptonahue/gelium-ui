# Gelium UI

Gelium UI is a themeable, open-code component library for server-rendered applications. Use it when you need native HTML semantics, zero required component JavaScript, and token themes — built with Tailwind CSS v4 in the monorepo build, shipped to consumers as CSS, themes, templates, and optional JS.

Gelium is not an SPA framework and not a React component library. Components are server-rendered HTML partials with documented server contracts, styled by CSS custom properties (`--ui-*`). Menus, dialogs, tooltips, and selection use platform features (`:checked`, Popover API, Invoker Commands) or server round-trips. HTMX is optional progressive enhancement.

The package on npm is [`gelium-ui`](https://www.npmjs.com/package/gelium-ui). This documentation site is the first dogfood app (Go), not the install path for product UI. Source: [GitHub](https://github.com/cryptonahue/gelium-ui). License: [MIT](https://github.com/cryptonahue/gelium-ui/blob/main/LICENSE).

## The foundation

Every Gelium component starts from the platform, not from a framework:

- **Native semantics first** — real `<button>`, `<dialog>`, `<input type="radio">`, `<table>`, `<nav>`; ARIA only where the platform has no equivalent.
- **No required component JavaScript** — selection, menus, tooltips, and dialogs run on declarative platform features or plain server round-trips.
- **Server-rendered by design** — HTMX as an enhancement, never a requirement.
- **Accessible by default** — light/dark, forced-colors, reduced-motion, keyboard focus, and RTL are part of every component contract, tested in the build.

## Quick start

### 1. Install (consumers)

```bash
npm install gelium-ui
```

### 2. Stylesheet

Drop-in bundle:

```css
@import "gelium-ui/dist/gelium.css";
```

Or compose with your Tailwind build:

```css
@import "tailwindcss";
@import "gelium-ui/themes/theme-material.css";
@import "gelium-ui/themes/theme-basecoat.css";
@import "gelium-ui/styles/index.css";
```

### 3. Theme class

Material is the default direction. Set the class on the document root (no JS required):

```html
<html class="theme-material">
<html class="theme-material theme-dark" data-theme="dark">
<html class="theme-basecoat">
```

On this docs site, preview with `?theme=basecoat` and `?scheme=dark`, or use the topbar controls.

### 4. Partials and optional JS

Copy HTML from `node_modules/gelium-ui/templates/` (and the component pages on this site). Optional:

```html
<script defer src="node_modules/gelium-ui/js/gelium.js"></script>
```

Provides toast region wiring, HTMX 422 swap when `X-Gelium-Validation: true`, and related helpers.

### Docs app only (this repository)

The documentation handler is a Go module used to dogfood the library. Set `BASE_URL` for canonical URLs, Open Graph, and the sitemap when you deploy **this** site:

```sh
BASE_URL=https://your-domain.example go run .
```

Product UIs should depend on **npm `gelium-ui`**, not on embedding `internal/app` as a component library.

## Themes and the layer model

Gelium UI separates *structure* from *visual direction* with six ownership layers:

- **Core tokens** — the structural contract: spacing, type, radius, motion, z-index (`--ui-space-4`, `--ui-type-body-md`). Theme-neutral by design.
- **Themes** — token values only: each theme maps the same `--ui-*` roles onto a visual direction (color, shape, type, motion).
- **Components** — anatomy and behavior of one primitive, consuming tokens, declaring only its own scoped tokens.
- **Patterns** — shared compositions across components (state layers, validation summaries, screen recipes).
- **Recipes** — full screen templates composed from primitives on the canonical server contract.
- **Screens** — complete applications built from recipes (the WhatsApp demo is one).

A theme never changes markup, component anatomy, or server contracts — only token values. That is why switching between Material and Basecoat is a single class change, with no rebuild and no JavaScript.

### Material (default)

Material is the default direction, built on **Material 3 (M3)** roles (color, type, state layers, elevation) expressed as token values in `gelium-ui/themes/theme-material.css` — not a copy of M3 React components.

### Basecoat

Basecoat is the alternative direction in the same package (`gelium-ui/themes/theme-basecoat.css`): near-black ink on white, hairline borders, soft shadows, selected via `theme-basecoat` / `?theme=basecoat`. Both themes support the explicit `theme-dark` class route.

## Base UI: vocabulary, never runtime

[Base UI](https://base-ui.com) is a React headless-primitive library that Gelium UI uses as a **vocabulary reference only** — focus management, component states, and keyboard-navigation patterns. It is **never a runtime dependency**: Base UI is React, and Gelium UI ships no React and no required component JavaScript. When Gelium documents a "headless primitive" behavior, it implements that behavior natively with platform features or server round-trips — the vocabulary is shared, the implementation is native.
