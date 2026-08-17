# Pagination

Pagination navigates a paginated result set: previous/next plus numbered pages. It is the standalone extraction of the [Data table](/components/data-table) footer nav, reused by feeds and queues. Every page is a real link so no-JS navigation works, the current page is a `span` with `aria-current="page"`, and the boundary links become `aria-disabled` spans instead of dead links.

## Specimen

This page renders the real `pagination` template markup live:

<nav class="ui-pagination" aria-label="Pagination">
  <a class="ui-pagination-page" href="/recipes/ops-queue?page=2">Previous</a>
  <a class="ui-pagination-page" href="/recipes/ops-queue?page=1">1</a>
  <a class="ui-pagination-page" href="/recipes/ops-queue?page=2">2</a>
  <span class="ui-pagination-page ui-pagination-page--current" aria-current="page">3</span>
  <a class="ui-pagination-page" href="/recipes/ops-queue?page=4">4</a>
  <a class="ui-pagination-page" href="/recipes/ops-queue?page=5">5</a>
  <a class="ui-pagination-page" href="/recipes/ops-queue?page=4">Next</a>
</nav>

On the first or last page the boundary control becomes a disabled span instead of a link:

<nav class="ui-pagination" aria-label="Pagination">
  <span class="ui-pagination-page ui-pagination-page--disabled" aria-disabled="true">Previous</span>
  <a class="ui-pagination-page" href="/recipes/ops-queue?page=1">1</a>
  <span class="ui-pagination-page ui-pagination-page--current" aria-current="page">2</span>
  <a class="ui-pagination-page" href="/recipes/ops-queue?page=3">3</a>
  <span class="ui-pagination-page ui-pagination-page--disabled" aria-disabled="true">Next</span>
</nav>

## Guidance

### When to use

Use pagination when a server-rendered list is paginated and the total is more than one page — a data table footer, a feed, a queue, or search results. The server decides the current page, clamps it into `[1, total]`, and renders every destination as a real link.

### When not to use

Do not render pagination for a single page — hide it entirely. Do not use it for a multi-step process: a stepper is a workflow with validation per step, while pagination pages a result set. Do not use it to switch between peer sections of a page; that is [Tabs](/components/tabs). Never render the boundaries as dead links: a disabled span with `aria-disabled="true"` is the contract.

### Usability

- The current page is clamped into `[1, total]`, so the server never renders page 0 or page 99.
- Numbered pages plus Previous/Next follow the same href function as the host — each recipe paginates over its own query vocabulary (`?page=`, `?p=`).
- The nav aligns to the end (`justify-content: flex-end`) to match the data table footer, and pills wrap on narrow screens.
- Pass a custom `aria-label` when the nav needs a more specific name than the default "Pagination".

### Accessibility

- The pager is a `<nav>` with an `aria-label` (default "Pagination"), so it is announced as a navigation landmark.
- Every destination is a real link; the current page is a `span` with `aria-current="page"` and is not focusable.
- Boundary states are `aria-disabled` spans — they are not focusable and never navigate, so there is no dead-link surprise.
- In forced-colors mode links repaint as `LinkText`, the current page as `Highlight`, and disabled boundaries as `GrayText`.

## Anatomy

- **`ui-pagination`** — the `<nav>`: a flex row with `--ui-space-1` gaps, end-aligned.
- **`ui-pagination-page`** — each pill: `min-width: 2rem`, height `2rem`, `--ui-radius-full`, label type `--ui-type-label-sm` in `--ui-pagination-page-color` (`--ui-color-fg-muted`), with a primary hover and the standard focus ring.
- **`ui-pagination-page--current`** — the current page: `--ui-pagination-active-color` (`--ui-color-primary`) at weight 600.
- **`ui-pagination-page--disabled`** — the boundary state: muted ink at `--ui-state-disabled-opacity` with `pointer-events: none`.

The pills use the scoped `--ui-pagination-*` tokens so the primitive works standalone and a theme may override it globally.

## Variants

The partial is one shape — previous/next plus numbered pages — with two positions: **standalone** (its own `<nav>`, as here) and **integrated** (the footer nav of the [Data table](/components/data-table) and the recipes that reuse it). The only configurable surface is the `aria-label` and the href vocabulary.

## Sources

- Registry: `docs/gelium-ui-component-registry.md` (Pagination row, recipe primitive RP; §3 server contract `GET ?page=` with clamping) — `.ui-pagination`, `--ui-pagination-{page-color,active-color}`, `--ui-radius-full`; states current (`aria-current`) and disabled boundary.
- Vocabulary: `docs/gelium-ui-vocabulary.md` §5 — Pagination is result-set navigation; Steps is a process, not a pager.
- Implementation: `lib/templates/pagination.html`, `lib/styles/pagination.css`; view model `paginationView` and `newPaginationView` in `internal/app/pagination.go`.

See also: [Data table](/components/data-table) for the integrated footer, [Ops Queue](/recipes/ops-queue) for a recipe consuming the standalone partial.