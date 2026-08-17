# Why Gelium (comparison)

This page is for evaluators: **when Gelium is the right tool, and when it is not.** It is not a dunk on other libraries — Radix, shadcn/ui, Base UI, and friends are excellent at what they optimize for.

## One-line positioning

| Library family | Optimizes for |
|---|---|
| **Gelium UI** | Server-rendered HTML + CSS tokens + optional tiny JS. Works with **JS disabled**. |
| **Radix / Base UI / React Aria** | Accessible **headless** primitives in a **React** (or similar) tree. |
| **shadcn/ui** | Copy-paste **React** components styled with Tailwind — you own the code in your repo. |
| **Naive UI / MUI / etc.** | Full **client component** kits with rich interaction models. |

Gelium is closer to **GOV.UK Frontend / open-code HTML** than to a React design system.

## Payload reality (orders of magnitude)

Numbers are **approximate and stack-dependent**. Measure your own bundle.

| Surface | JS (order of magnitude) | Notes |
|---|---|---|
| **Gelium consumer (enhancement only)** | **~5 KB** (`gelium.js`) | Toast, 422 swap helper, slider fill, VT guard. **0 KB required** for core UI. |
| **Gelium + HTMX 4** | **~40 KB** | `htmx.min.js` (~36 KB) + `gelium.js` (~5 KB). Typical progressive-enhancement stack. |
| **This docs site (all static JS)** | **~50 KB** | HTMX + gelium + docs chrome + search + morph helper. CSS is larger (~210 KB compiled app CSS including themes). |
| **Typical React + Radix/shadcn app shell** | **hundreds of KB** (often cited ~500–650 KB+ gzipped less, raw more) | Framework runtime + primitives + your tree. Not evil — different product shape. |

**Takeaway:** Gelium’s bet is **HTML and CSS first**. JS is progressive enhancement, not the substrate.

## Contract comparison

| Concern | Gelium | Radix / Base UI | shadcn/ui |
|---|---|---|---|
| **Runtime** | Browser + your server (any language) | React (or port) | React + your copy of components |
| **Works without JS** | **Yes** (forms, links, native controls) | No (component tree) | No (component tree) |
| **Styling** | `--ui-*` tokens + themes by **class on `<html>`** | Bring your own | Tailwind classes you own |
| **Install** | `npm install gelium-ui` (CSS/HTML/JS assets) | npm packages | CLI / copy into repo |
| **Server authority** | First-class (422 + `X-Gelium-Validation`, `gelium:toast`) | App-defined | App-defined |
| **Theming** | Ship `theme-material` + `theme-basecoat`; select by class | Tokens / CSS vars in your app | Your Tailwind theme |
| **Copy-paste HTML** | `templates/*.html` in the package | N/A (JSX) | You already own TSX |
| **HTMX** | Designed to pair | Possible but not the model | Possible but not the model |

## When to choose Gelium

- Multi-page or **server-rendered** apps (Go, Rails, Laravel, Django, PHP, plain HTML).
- **HTMX** (or similar) progressive enhancement.
- Legal / gov / internal tools where **no-JS** and predictable HTML matter.
- You want a **published CSS package** (`gelium-ui`) without adopting React.
- Token themes switched by **document class**, not a JS theme provider.

## When NOT to choose Gelium (no-gos)

Be honest with yourself:

1. **You need React/Vue/Svelte components as the API** — use Radix, Base UI, shadcn, Naive, MUI.
2. **You need a rich client-only SPA** with client routers and heavy client state as the core — Gelium will feel like the wrong layer.
3. **You need Shadow DOM / Web Components runtime** — not Gelium’s model (open DOM + CSS).
4. **You expect a Figma-to-production closed kit with every pattern already productized** — Gelium is still growing; recipes and components are open-code, not a closed enterprise suite.
5. **You refuse to own server HTML** — if the team only ships client bundles, pick a client library.

## “50 KB vs 625 KB”

Use it as a **conversation starter**, not a lab measurement of every competitor:

- **~50 KB** ≈ this docs site’s **JS** surface (HTMX + enhancements + chrome), or a lean HTMX+Gelium consumer.
- **~625 KB** ≈ a ballpark many teams see for **framework + headless UI + app chrome** before tree-shaking/gzip storytelling.

Always compare **your** production bundles (raw and gzip), same routes, same features. Gelium’s claim is structural: **the default path does not require a framework runtime.**

## How Gelium relates to shadcn-style workflows

| shadcn idea | Gelium analogue |
|---|---|
| You own the code | Open-code HTML/CSS in `templates/` + `styles/`; fork or copy freely (MIT) |
| Tailwind | Yes (v4); also a prebuilt `dist/gelium.css` if you do not want a pipeline |
| CLI add component | Not yet — install the package and copy/import partials |
| Registry of components | Docs `/components/*` + npm package files |

## Install reminder

```bash
npm install gelium-ui
```

```css
@import "gelium-ui/dist/gelium.css";
```

Details: [npm gelium-ui](https://www.npmjs.com/package/gelium-ui), agent brief [`/llms.txt`](/llms.txt), [Contributing](/docs/contributing).

## See also

- [Design principles](/docs/principles)
- [Server contracts](/docs/server-contracts)
- [Acknowledgments](/docs/acknowledgments) (what we learned from other systems)
- [Roadmap](/docs/roadmap)
