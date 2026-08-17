# Breadcrumb

Breadcrumb is the canonical location trail: Home → section → page. Use it to show where the current page sits in the hierarchy and give one-click way back. It is a native `nav` → `ol` → `li` structure styled with `ui-breadcrumb`, so it needs no component JavaScript. The current crumb is never a link — it is a `span` with `aria-current="page"` — and the separator is pure CSS, so no literal characters leak into markup.

## Specimen

This page renders the real `breadcrumb` template markup live:

<nav aria-label="Breadcrumb">
  <ol class="ui-breadcrumb">
    <li class="ui-breadcrumb-item"><a href="/">Home</a></li>
    <li class="ui-breadcrumb-item"><a href="/docs">Docs</a></li>
    <li class="ui-breadcrumb-item"><span aria-current="page">Breadcrumb</span></li>
  </ol>
</nav>

## Guidance

### When to use

Use a breadcrumb when a page sits more than one level below a section root and the hierarchy is worth showing — a deep docs page, a product detail under a catalog, a settings page under a workspace. The crumb data comes from the same server-side model as the JSON-LD `BreadcrumbList`, so the visible trail and the machine-readable trail cannot drift.

### When not to use

Do not add a breadcrumb on the home page or on single-level sites where the trail adds no information. Do not make the current crumb a link: linking to the page you are on is a dead loop. For switching between peer sections use [Tabs](/components/tabs) instead, and for navigating pages of a result set use [Pagination](/components/pagination) — a breadcrumb is location, not paging.

### Usability

- Keep the trail short: three to five crumbs, with the current page named plainly as the last crumb.
- The separator is `--ui-breadcrumb-separator` painted by CSS, so locale copy never has to ship a literal "›" in markup.
- Items are `display: inline-flex` with `--ui-breadcrumb-gap` and the list wraps (`flex-wrap`), so long trails collapse gracefully on narrow screens.

### Accessibility

- The trail is a `<nav aria-label="Breadcrumb">` wrapping an `<ol>` of `<li>` items, so it is announced as a navigation landmark and read in order.
- The current crumb is a `span` with `aria-current="page"` — screen readers announce it as the current page, and it is not focusable.
- In forced-colors mode links repaint as `LinkText` and the current crumb plus separator repaint as `CanvasText`, so the trail survives without color.

## Anatomy

- **`ui-breadcrumb`** — the `<ol>`: a flex list with no list chrome (`list-style: none`, zero margin/padding) and `--ui-breadcrumb-gap` (`--ui-space-1`) between items.
- **`ui-breadcrumb-item`** — each `<li>`: label type `--ui-breadcrumb-type` (`--ui-type-label-sm`) in muted ink (`--ui-breadcrumb-color`).
- **Separator** — the `::before` of every item after the first, painting `--ui-breadcrumb-separator` in muted ink; links hover to the current-color ink.
- **Current crumb** — the last `<li>` holds a `span` with `aria-current="page"` in `--ui-breadcrumb-current-color` (`--ui-color-fg`); it is never an anchor.

All ink and spacing come from the scoped `--ui-breadcrumb-*` tokens, so a consumer can retune the trail per placement.

## Variants

The trail has no visual variants: `nav` → `ol` → `li` (current as `span[aria-current="page"]`) is the canonical contract, shared by the docs shell breadcrumb and `componentBreadcrumb`. The only dimension that varies is the data — which crumbs, in which order — supplied by the server.

## Sources

- Registry: `docs/gelium-ui-component-registry.md` (Breadcrumb row, phase P) — `.ui-breadcrumb`, `--ui-breadcrumb-*`, `--ui-color-fg[-muted]`; state current (`aria-current="page"`); CSS separator, never literal text.
- Vocabulary: `docs/gelium-ui-vocabulary.md` §5 — Breadcrumb blocks GEO §9/§14; data shares its source with the JSON-LD `BreadcrumbList`.
- Implementation: `lib/templates/breadcrumb.html`, `lib/styles/breadcrumb.css`; crumb data in `internal/app/server.go` (`breadcrumbView`, `componentBreadcrumb`).

See also: [SEO](/docs/seo) for the `BreadcrumbList` metadata, [Information architecture](/docs/information-architecture) for trails inside the docs IA.