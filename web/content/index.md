# Gelium UI

Gelium UI is a themeable, open-code component library for server-rendered applications. Use it when you need native HTML semantics, zero component JavaScript, and a Material 3 design system — built with Tailwind CSS v4 and HTMX.

Gelium UI is not an SPA framework and not a React component library. Every component is a server-rendered HTML partial with a documented server contract, styled by CSS custom properties (`--ui-*` tokens). There is no component JavaScript to bundle: menus, dialogs, tooltips, and selection run on declarative platform features (`:checked`, the Popover API, Invoker Commands, Interest Invokers) or plain server round-trips. HTMX is an optional progressive enhancement, never a requirement.

The code is open — this documentation site is itself the first Gelium application, and every partial it renders is copyable. Start from the [GitHub repository](https://github.com/cryptonahue/gelium-ui), embed the handler, and pick the visual direction you want.

## The foundation

Every Gelium component starts from the platform, not from a framework:

- **Native semantics first** — real `<button>`, `<dialog>`, `<input type="radio">`, `<table>`, `<nav>`; ARIA only where the platform has no equivalent.
- **No component JavaScript** — selection, menus, tooltips, and dialogs run on declarative platform features (`:checked`, Popover API, Invoker Commands, Interest Invokers) or plain server round-trips.
- **Server-rendered by design** — HTMX as an enhancement, never a requirement. The docs you are reading are themselves the first Gelium application.
- **Accessible by default** — light/dark, forced-colors, reduced-motion, keyboard focus, and RTL are part of every component contract, tested in the build.

## Quick start

Four steps from zero to a themed page:

1. **Embed the handler.** Gelium UI is a Go module (`geliumui`). `app.New()` returns an `http.Handler` you mount on any path:

   ```go
   package main

   import (
       "net/http"

       "geliumui/internal/app"
   )

   func main() {
       http.ListenAndServe(":8787", app.New())
   }
   ```

   Set `BASE_URL` to your production origin before deploying — canonical URLs, Open Graph tags, breadcrumb JSON-LD, and the sitemap all derive from it. Without it, absolute links point at the default `https://gelium-ui.example`:

   ```sh
   BASE_URL=https://your-domain.example go run .
   ```

2. **Serve the stylesheet bundle.** The compiled single-file bundle (`/static/app.css`) carries every theme; link it once and switch directions at runtime:

   ```html
   <link rel="stylesheet" href="/static/app.css">
   ```

3. **Pick a theme.** Material is the default. Append `?theme=basecoat` to any URL to preview the Basecoat direction, use the Theme control in the docs topbar, or set the class directly on the document root (`<html class="theme-basecoat">`). Dark mode is an explicit class route too: `?scheme=dark` or `<html class="theme-basecoat theme-dark">`.

4. **Copy the partials you need.** Open any component page, copy its markup and its server contract, and paste it into your own templates. This is the open-code model: nothing is generated, nothing is hidden behind a runtime.

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

Material is the default direction, built on **Material 3 (M3)**, Google's design system. Gelium UI takes from M3 the *role-based* approach to design, not a copy of its components:

- **Color roles, not fixed palettes** — `primary`, `secondary`, `tertiary`, `surface`, `error`, and their foregrounds are semantic roles; the theme supplies the values. This is why the same component contract works under a tonal M3 palette and under Basecoat's near-black ink.
- **Type scale** — M3's display/headline/title/body/label hierarchy maps directly onto Gelium's `--ui-type-*` steps.
- **States** — hover, pressed, focus, disabled, and loading follow M3's state-layer model, expressed through `--ui-state-*` opacity tokens.
- **Elevation and shape** — tonal surfaces, rounded corners, and elevation levels (`--ui-shadow-0..5`, `--ui-radius-*`) come from M3's system.

Gelium implements M3 as **token values in `themes/theme-material/theme.css`** — components never hardcode a Material look. M3 defines the *visual direction*; the layer model keeps it swappable, which is exactly why Basecoat can live in the same bundle. (Compare this with Base UI below: M3 contributes aesthetics and roles; Base UI contributes headless *behavior* vocabulary — two different influences, both documented.)

### Basecoat

Basecoat is the Phase I alternative direction — a translation of the Basecoat UI "Vega" style pack (a shadcn-compatible, Tailwind 4 CSS system) into the Gelium `--ui-*` vocabulary: near-black ink on white, hairline borders, soft shadows, a 0.625rem base radius, Geist typography, and 2.25rem control density. It ships in the same bundle as Material and is selected with `?theme=basecoat` or the `theme-basecoat` class. Both themes support the explicit `theme-dark` class route.

## Base UI: vocabulary, never runtime

[Base UI](https://base-ui.com) is a React headless-primitive library that Gelium UI uses as a **vocabulary reference only** — focus management, component states, and keyboard-navigation patterns. It is **never a runtime dependency**: Base UI is React, and Gelium UI ships no React and no component JavaScript. When Gelium documents a "headless primitive" behavior, it implements that behavior natively with platform features (`:checked`, the Popover API, Invoker Commands) or server round-trips — the vocabulary is shared, the implementation is native.

## Get started

Read the [documentation](/docs) or jump straight into a component: [Button](/components/button), [Text field](/components/text-field), [Dialog](/components/dialog), or [Toast](/components/toast).

See the full library in action with the [WhatsApp manager demo](/demo/whatsapp).
