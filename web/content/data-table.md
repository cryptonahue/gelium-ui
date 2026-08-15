# Data table

Data table is a server-rendered Material 3 table for displaying dense, sortable, filterable, paginated data. Gelium reimplements it over native HTML table semantics — a real `<table>` with `<thead>`, `<tbody>`, `<th scope="col">`, and `<caption>`. There is no component JavaScript: sort, filter and pagination are real GET requests answered by the server, and row selection uses native checkboxes in a real form. Use a data table when dense, columnar data benefits from server-side sort, filter and pagination.

## Guidance

### When to use

Use a data table for dense, columnar data that benefits from server-side sort, filter and pagination — especially when the dataset lives remotely and no-JS behavior must stay first-class.

### When not to use

For simpler lists of items, prefer the [List](/components/list) component. If the data has no sorting, filtering or pagination needs, plain markup or a list reads better than a table surface.

### Usability

- Keep the native structure: a real `<table>` with `<thead>`, `<th scope="col">`, and a `<caption>` naming the slice.
- Sort, filter and pagination are real GET requests, so the no-JS flow is a normal full-page reload.
- Row selection uses native checkboxes in a real form; the header checkbox submits `selection=all`.

### Accessibility

- Keep the native elements: `<table>`, `<caption>`, `<th scope="col">`, native checkboxes.
- Sort state is never communicated by arrow alone: the active header also exposes `aria-sort` on the `<th>`.
- The caption gives the table an accessible name that also describes the current slice.
- The state layer is decorative; row selection is always reflected in the native checkbox.

See [Choose the right control](/docs/choose-the-right-control) for the cross-component decision.
## Anatomy

The root `ui-data-table` container paints the surface with a `1px` outline and the container radius. Inside it the `table.ui-data-table-table` holds a `caption`, a `thead` of `th.ui-data-table-cell` headers, and a `tbody` of `tr.ui-data-table-row` rows. The pagination footer (`nav.ui-data-table-pagination`) sits under the table inside the same container.

```html
<table class="ui-data-table-table">
  <caption class="ui-data-table-caption">12 rows · page 1 of 3</caption>
  <thead>
    <tr>
      <th scope="col" class="ui-data-table-cell ui-data-table-cell--checkbox">[select all]</th>
      <th scope="col" class="ui-data-table-cell ui-data-table-cell--sortable" aria-sort="ascending">
        <a class="ui-data-table-sort ui-data-table-sort--active" href="?sort=name&dir=desc">Name ↑</a>
      </th>
    </tr>
  </thead>
  <tbody>
    <tr class="ui-data-table-row">
      <td class="ui-data-table-cell ui-data-table-cell--checkbox">[checkbox]</td>
      <td class="ui-data-table-cell ui-data-table-cell--label">Alpha release</td>
    </tr>
  </tbody>
</table>
```

- **Container** — `ui-data-table`, `56px` header row, `52px` body rows, `16px` cell padding, sort icon in the header, checkbox column on the inline-start.
- **Header** — `ui-data-table-cell--sortable` headers are real links (`ui-data-table-sort`) so they activate with a normal GET and show a keyboard focus ring.
- **Rows** — `ui-data-table-row` with the Material hover state layer; selected rows flip the native checkbox and tint the row container.
- **Footer** — `ui-data-table-pagination` with previous/next and numbered page links, plus a `caption` summarizing the slice.

## How do sort, filter and pagination work?

Everything is a real GET request against `/components/data-table`, so the no-JS flow is a normal full-page reload:

- **Sort** — the active header link toggles `?sort=name&dir=desc`; clicking another column sorts ascending. The active `<th>` carries `aria-sort="ascending"` or `"descending"` and shows the direction arrow.
- **Filter** — a real `<form method="get" action="/components/data-table">` submits `?q=...`; the server filters rows by name or status (case-insensitive) and re-renders.
- **Pagination** — page links carry `?page=2` and the current page is marked with `aria-current="page"`.

## Row selection

The selection column uses native checkboxes inside a real `<form method="get">`. Selecting rows and submitting the form sends `?selection=<id>&selection=<id>`; the server marks the matching rows `checked` and re-renders a status notice — the same no-JS multi-select contract as List. The header checkbox submits `selection=all` to select every row on the current page.

## States

- **Sortable header** — rest, hover (label and glyph brighten), `focus-visible` (focus ring).
- **Row** — rest, hover (state layer), selected (`ui-data-table-row--selected` + native `:checked`, the CSS pseudo-class matching a checked native input), and focus that lands on the native checkbox.
- **Disabled** — not part of the table contract; individual controls keep their own disabled semantics.

## Progressive enhancement

The sort, filter and pagination links also carry `hx-get` targeting `#data-table-panel` with `outerHTML` swap, so with HTMX (the progressive-enhancement library that swaps server-rendered fragments) enabled only the table panel re-renders. The no-JS branch is always a real full-page GET — nothing depends on HTMX or JavaScript. The remote refresh demo completes a POST round-trip that re-renders the panel with a determinate `.ui-progress` bar and a `.ui-toast` notification; HTMX swaps just the refresh form and raises `gelium:toast` into the live region.

## Accessibility

- Keep the native elements: `<table>`, `<caption>`, `<th scope="col">`, native checkboxes.
- Sort state is never communicated by arrow alone: the active header also exposes `aria-sort` on the `<th>`.
- The caption gives the table an accessible name that also describes the current slice.
- The state layer is decorative; row selection is always reflected in the native checkbox.
- Decorative sort SVGs are `aria-hidden` and `focusable="false"`.
- In forced-colors mode the container keeps its outline, links keep `LinkText`/`Highlight`, and the selected row switches to `Highlight`/`HighlightText`.

## Design tokens

The `--ui-data-table-*` tokens are declared scoped to the root so the primitive works standalone, and the theme may override them globally.

| Token | Meaning |
| --- | --- |
| `--ui-data-table-container-color` | Container background (surface) |
| `--ui-data-table-outline-color` | Container and row divider outline |
| `--ui-data-table-header-height` | Header row height (`56px`) |
| `--ui-data-table-row-height` | Body row height (`52px`) |
| `--ui-data-table-cell-padding` | Cell horizontal padding (`16px`) |
| `--ui-data-table-header-color` | Header label color (on-surface-variant) |
| `--ui-data-table-header-hover-color` | Header hover color (on-surface) |
| `--ui-data-table-label-color` | Body label color (on-surface) |
| `--ui-data-table-selected-container-color` | Selected row background |
| `--ui-data-table-hover-opacity` | Row hover state layer opacity |
| `--ui-data-table-checkbox-size` | Row checkbox edge (`18px`) |
| `--ui-data-table-sort-icon-size` | Header sort glyph edge (`18px`) |

## Keyboard

Because every control is native, keyboard behavior comes for free: sort headers are links (Tab to focus, Enter to activate), checkboxes toggle with Space, and page links are ordinary anchors. There is no custom roving tabindex or table-navigation JavaScript.
