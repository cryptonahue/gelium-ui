# Basecoat / Gelium audit (verified package)

Date: 2026-08-19

## Provenance

- Package: `basecoat-css@1.0.2`, obtained with `npm pack basecoat-css` from the npm registry.
- Repository metadata from `npm view`: `https://github.com/hunvreus/basecoat.git`.
- Archive: `basecoat-css-1.0.2.tgz` in this directory; SHA-256:
  `b493e77a7ee0b945e41398f5aee186ffa692f930751b960f5ac0853e4b54de48`.
- License evidence: `package/LICENSE.md`, MIT License, copyright 2025 Ronan Berder.
- Exact reference files audited: `package/dist/base/base.css`,
  `package/dist/components/accordion.css`, and `package/dist/styles/vega.css`.
- No package dependency, production import, Basecoat class, icon asset, or Basecoat
  runtime was added to Gelium.

## Verified Vega accordion anatomy

The official component and Vega style pack resolve to the following behavior and
visual rules (the component file supplies structure; Vega supplies the visual
values):

| Concern | Evidence in package | Gelium translation |
|---|---|---|
| Root geometry | `.accordion { display:flex; width:100%; flex-direction:column }` | The `html[data-gelium-reference="basecoat"]` / later skin token adapter requests a full-width flex column; no 48rem cap or grid gap. |
| Item separation | Vega: `details:not(:last-child) { border-bottom: 1px }` | Only non-final `.ui-accordion-item` gets a bottom border; no item card border, radius, surface, or shadow. |
| Trigger | `summary`: flex, `items-start`, transparent border, rounded-md, `py-4`, `text-sm font-medium`, `text-left` | 16px block padding, 14px/20px medium sans, start alignment, transparent 1px border and 6px radius. Inline padding is zero because Vega has no `px-*`. `text-align:start` preserves RTL direction. Gelium additionally retains a >=44px native-control floor; Vega does not specify one. |
| Hover | `hover:underline` | Transparent background and underline; previous color-mix fill was not official Vega behavior. |
| Focus | `focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-3` | Ring color plus 3px 50%-alpha ring and border; no invented shadow/surface. |
| Chevron | optional trailing SVG: `text-muted-foreground ms-auto size-4`; open state `rotate-180` | 16px trailing icon, muted foreground, auto inline start margin, 180° open rotation. SVG geometry remains Gelium's existing inline chevron; no Basecoat icon asset copied (Basecoat says icons are not bundled). |
| Panel | `:not(summary) { overflow:hidden; padding-top:0; padding-bottom:1rem; text-sm }` | Zero inline padding, 0/16px block padding, 14px/20px body font. |
| Disabled | component: details or summary `aria-disabled=true` gets `pointer-events:none; opacity:.5` | Same selector and opacity in the Basecoat behavior profile. Optional Basecoat JS additionally prevents click/keyboard toggles. |
| Multiple/open | native `details`; official JS closes siblings unless root has `data-multiple`; `open` controls initial state | Gelium retains native no-JS baseline. Server/native named details behavior remains; Basecoat's single-open and disabled enforcement is documented optional reference behavior, not bundled. |
| Motion/accessibility | Component uses transitions; no RTL-specific visual override; native details semantics | Gelium keeps reduced-motion and forced-colors rules. `text-align:start`, logical margin/padding, and native disclosure preserve RTL. |

## Why the prior page differed

The old translation styled an outer bordered white card, gave every item a border,
used a grid gap, capped the root at 48rem, added horizontal panel/trigger padding,
used a 48px minimum trigger, and applied a color-mix hover background. Those are
not in the official Vega accordion CSS. The resulting page was therefore visibly
more like a card accordion than Basecoat's flat divider list. The icon was also
20px by default instead of Vega's 16px muted trailing icon.

## Direct matches, approximations, and mismatches

### Direct/evidence-backed matches now implemented

- Vega's full-width flex-column root.
- Flat transparent items with bottom dividers only between items.
- `py-4` trigger and panel `pt-0 pb-4` rhythm.
- 14px medium trigger and 14px panel text.
- 6px trigger radius, transparent border, 3px focus ring, ring at 50% opacity.
- Underline hover, muted 16px trailing chevron, and 180° open rotation.
- Disabled opacity/pointer behavior in CSS.
- Existing Gelium semantic `ui-*` markup and no-JS native `<details>/<summary>` baseline.

### Deliberate approximations / translation limits

- Basecoat's Tailwind `ring-ring/50` is represented with `color-mix()` and Gelium's
  focus token; the token is a solid Gelium color, while the official ring is an
  alpha color. This is an equivalent mechanism, not a copied asset.
- `font: var(--ui-accordion-trigger-font)` and panel font encode Vega's `text-sm`
  (14px/20px) because Gelium's semantic typography token is not the Basecoat utility.
- Gelium's inline SVG is retained. Basecoat explicitly does not bundle icons and
  its page examples use inline Lucide SVGs; no unverified SVG/data URI was copied.
- Basecoat's optional JS single-open/disabled enforcement is not included. Native
  HTML remains functional without JS; this is an intentional behavior boundary.

### Remaining non-identity differences

- Gelium's generated template uses `ui-*` classes and a `<section>` panel rather
  than Basecoat's `.accordion` classes; this is required by Gelium's single markup
  contract.
- The theme still maps broad Basecoat semantic colors into Gelium's larger token
  contract. Elevation, state opacity, status colors, and typescale outside this
  accordion are not official Basecoat tokens and remain documented derivations.
- Basecoat's CSS is authored with Tailwind `@apply`; Gelium ships plain CSS and
  custom properties, so computed output is translated rather than wholesale copied.

## Base UI is separate

Base UI's official About and Styling pages state that Base UI components are
**unstyled, do not bundle CSS, and do not prescribe a styling solution**. Its source,
API, state data attributes, and CSS variables are useful behavior/style-hook
references only. There is no official Base UI visual CSS pack to download. Gelium's
`theme-baseui.css` must therefore not claim to be a Base UI visual pack or derive
visual values from Basecoat. No Base UI CSS/runtime/package was downloaded or added
by this audit.

Sources:

- https://basecoatui.com/installation
- https://basecoatui.com/components/accordion/
- https://base-ui.com/react/handbook/styling
- https://base-ui.com/react/overview/about
