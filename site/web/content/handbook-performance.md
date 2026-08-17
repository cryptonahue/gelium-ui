# Performance stance

Gelium’s size story is intentional product stance, not an apology. **JS is progressive enhancement.** **CSS carries themes, tokens, and components** so the UI stays coherent without a framework runtime.

This page documents what we ship, what we optimize, and how to measure it yourself. For evaluator positioning vs React kits, see [Why Gelium](/docs/compare).

## Stance in one paragraph

A Gelium consumer does **not** need a client framework. Core UI is **server HTML + CSS**. Optional scripts improve toast regions, validation swaps, slider fill, and similar — they are not the substrate. The docs site loads more JS (HTMX, chrome, search) and still stays on the order of **~50 KB** of static JS. Compiled CSS is larger (~**210 KB** on this site) **by design**: tokens, two themes, and the component surface live in CSS.

## Measured sizes (this repo)

Byte counts are **raw file size** (`wc -c`), not gzip. Re-measure after upgrades; treat numbers as orders of magnitude.

### JavaScript (static)

| Asset | ~bytes | Role |
|---|---:|---|
| `gelium.js` | ~4.6 KB (~5 KB) | Library enhancement (toast, 422 helper, slider fill, VT guard). **0 KB required** for core UI. |
| `htmx.min.js` | ~36 KB | Optional progressive-enhancement stack (docs + typical consumers). |
| Docs static JS total | **~50 KB** | HTMX + `gelium.js` + docs chrome (`app.js`) + search + morph helper. |

Typical consumer stack: **`gelium.js` alone (~5 KB)** or **HTMX + gelium (~40 KB)**. The **~50 KB** figure is this docs site’s full static JS surface.

### CSS (compiled)

| Asset | ~bytes | Role |
|---|---:|---|
| Docs `app.css` | **~210 KB** | Built stylesheet for this site: tokens, **two themes**, components, docs chrome. |
| Package `dist/gelium.css` | ~169 KB | Prebuilt library CSS for consumers who skip a Tailwind pipeline. |

**CSS is the biggest asset by design** — not accidental bloat. Shipping unthemed or half-themed UI would “win” on KB and lose on product quality.

### npm package

| Measure | Value |
|---|---|
| `npm pack` tarball (`gelium-ui`) | **~87 KB** package size (gzipped archive; unpacked is larger) |

```bash
cd lib && npm pack --dry-run   # see "package size"
```

Install: [npm gelium-ui](https://www.npmjs.com/package/gelium-ui). Agent brief: [`/llms.txt`](/llms.txt).

## What we optimize

- **No framework runtime required** for core UI (no React/Vue/Svelte substrate on the critical path).
- **Server HTML** as the source of truth: forms, links, native controls work with **JS disabled**.
- **Optional JS** only where progressive enhancement earns its keep.
- **Token + theme CSS** so Material / Basecoat (and future themes) share one markup surface.
- **Honest docs**: same stack we recommend consumers use (dogfood).

## What we do not chase

- **Competing with zero-CSS utility-only pages** or “hello world” KB contests that ship no design system.
- **Shipping unthemed or broken UI** to look smaller on a spreadsheet.
- **Claiming false minification magic** — Tailwind may still include unused utilities until a deliberate purge/review; we do not market phantom KB wins.
- **Matching SPA framework bundles feature-for-feature** while remaining HTML-first — different product shape (see [Why Gelium](/docs/compare)).

## How to measure

From a clean checkout:

```bash
# Docs static JS + CSS (raw bytes)
wc -c site/web/static/*.js site/web/static/app.css

# Library enhancement + prebuilt CSS
wc -c lib/js/gelium.js lib/dist/gelium.css

# Published package tarball size
cd lib && npm pack --dry-run
```

Compare **your** production routes the same way (raw and gzip, same features). Gelium’s claim is structural: **the default path does not require a framework runtime**; CSS size funds tokens and themes.

## Tailwind scope note

The build uses Tailwind CSS 4. Compiled output can still carry **unused utilities** until we deliberately audit and tighten the content/source scan. Prefer a measured unused-utility review over marketing “fully tree-shaken” language we have not proven. Do not equate “ran the bundler” with “every byte is load-bearing.”

## See also

- [Why Gelium (comparison)](/docs/compare) — when to choose Gelium vs Radix / shadcn / Base UI
- [Design principles](/docs/principles)
- [Themes](/docs/themes) — why theme CSS is first-class
- [Tokens](/docs/tokens)
- [Server contracts](/docs/server-contracts)
- Agent brief: [`/llms.txt`](/llms.txt)
- Package: [npm gelium-ui](https://www.npmjs.com/package/gelium-ui)
