# Section heading

Section heading is a typographic utility that labels a block of related content. Use it when a section of a page needs a title that is scannable but quieter than the page title. The heading always renders an `h2` — never an `h1` — because the page owns a single `h1`. It is text-only, so it needs no component JavaScript at all.

## Examples

A plain section heading spans the content column.

<p class="ui-section-heading-eyebrow">Eyebrow</p>
<h2 class="ui-section-heading">Features</h2>

A centered heading suits a landing block that should stay symmetric on wide screens.

<p class="ui-section-heading-eyebrow">Eyebrow</p>
<h2 class="ui-section-heading ui-section-heading--centered">Features</h2>

Both specimens above are the live markup the template `section-heading.html` emits: an optional eyebrow paragraph followed by the `h2` itself.

## Guidance

### When to use

Use a section heading to title a block of related content inside a page — a features grid, a pricing section, a form group. It earns its place when the section needs a label that is readable at a glance without competing with the page title.

### When not to use

Do not use a section heading as the page title: the page owns a single `h1`. Do not wrap every paragraph in a heading — a heading must label a real block of related content, not decorate prose. If the block is a distinct, standalone unit people scan independently, a [Card](/components/card) or [Feature card](/components/feature-card) communicates the grouping better than a heading alone.

### Usability

- Pair the optional eyebrow (`ui-section-heading-eyebrow`) with the heading when a categorical label helps: the eyebrow uses the label typescale and muted color.
- The centered variant (`ui-section-heading--centered`) is text alignment only; keep it on the same element as the base class.
- Both variants use the scoped `--ui-section-heading-*` tokens, so they stay consistent across light and dark schemes without new theme work.

### Accessibility

- The heading is a real `h2` in document order, so assistive technology navigates the section structure for free.
- The eyebrow is decorative: it is a `p`, not a heading, and should repeat nothing essential that the heading does not already say.
- Text-only styling means no forced-colors block is needed: the heading inherits the canvas text color when color is removed.

See [Choose the right control](/docs/choose-the-right-control) for the cross-component decision.

## Anatomy

The heading uses the theme tokens.

- `--ui-section-heading-type` — the headline-sm typescale, mapped to `--ui-type-headline-sm`.
- `--ui-section-heading-color` — the heading paint, mapped to `--ui-color-fg`.
- `--ui-section-heading-margin` — the vertical rhythm around the heading, mapped to `--ui-space-*`.

The eyebrow uses the label-lg typescale and the muted foreground (`--ui-color-fg-muted`) with `--ui-space-1` below it.

## Variants

- `ui-section-heading` — the base h2 utility.
- `ui-section-heading--centered` — centers the text within its container.
- `ui-section-heading-eyebrow` — the optional categorical label above the heading.

## Anti-patterns

- Never use a heading level other than `h2` for a section label; `h1` belongs to the page title alone.
- Never style a `p` or `span` to look like a heading: the outline must stay semantic so the document structure survives without styling.

## Checklist

- The section has a real `h2` in document order.
- The page has exactly one `h1`, and it is not this heading.
- The heading labels a real block of related content.
- The eyebrow, when present, adds a categorical label and repeats no essential information.

## Accessibility

The heading participates in the document outline as an `h2`, so screen-reader navigation and the "On this page" rail pick it up automatically. Because the utility is text-only, it keeps the canvas text color in forced-colors mode without a dedicated block.