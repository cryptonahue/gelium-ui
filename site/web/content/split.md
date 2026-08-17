# Split

Split is a two-column editorial composition that pairs media with a body block and stacks into a single column on narrow screens. Use it when a section pairs one visual with a short story — an image or video on one side, a heading, copy, and one action on the other. It is pure layout: the body copy is deliberately not styled editorially, so the consumer keeps the typography contract.

## Examples

A split with media, an eyebrow, a title, copy, and one action.

<div class="specimen-block">
<section class="ui-split">
  <div class="ui-split-media">
    <img src="https://images.unsplash.com/photo-1441974231531-c6227db76b6e?auto=format&fit=crop&w=800&q=60" alt="Sunlight through a forest canopy" width="800" height="600" loading="lazy">
  </div>
  <div class="ui-split-body">
    <p class="ui-split-eyebrow">Field notes</p>
    <h2 class="ui-split-title">A slower way to read</h2>
    <p class="ui-split-copy">Long-form pages read better when the media gets a side of its own. Split gives the image a column and the story the other.</p>
    <div class="ui-split-action"><a class="ui-button ui-button-outline" href="/components/split"><span>Read the notes</span></a></div>
  </div>
</section>
</div>

The media slot is optional: without it the body takes the full row instead of a half column.

<div class="specimen-block">
<section class="ui-split">
  <div class="ui-split-body">
    <p class="ui-split-eyebrow">Announcement</p>
    <h2 class="ui-split-title">Gelium UI 0.5 is out</h2>
    <p class="ui-split-copy">Native semantics, server-driven state, and zero component JavaScript. The release notes walk through every change.</p>
    <div class="ui-split-action"><a class="ui-button ui-button-primary" href="/components/split"><span>Read the release notes</span></a></div>
  </div>
</section>
</div>

Both specimens above are the live markup the template `split.html` emits: an optional media slot followed by the body slots.

## Guidance

### When to use

Use a split when a section pairs one visual with a short narrative — a feature highlight, a field note, an announcement. It earns its place when the media communicates the section faster than the copy alone and the two belong side by side on wide screens.

### When not to use

Do not use a split for a grid of repeated items — that is a set of [Cards](/components/card) or [Feature cards](/components/feature-card). Do not use it when the media is decorative filler: the media slot should carry meaning. If the section is a single static region without media, a plain [Card](/components/card) or section heading communicates the grouping better.

### Usability

- Put one visual in the media slot and keep the body tight: an eyebrow, a title, a few lines, one action.
- Leave the body copy unstyled by default — Split is layout only. Style the copy with your prose type scale (the site prose contract) once per page, not per split.
- The columns stack automatically below `47.99rem`, so narrow screens read media-first, then body.

### Accessibility

- The split is a `<section>`: give it an accessible name from the title inside it, and keep the media's alt text descriptive when the image is content.
- RTL documents mirror the columns automatically — the tracks flow in the document direction, so no `left`/`right` literals exist to flip.
- In forced-colors mode the media keeps a `CanvasText` boundary so the two-column grouping stays visible when color is removed.

See [Choose the right control](/docs/choose-the-right-control) for the cross-component decision.

## Anatomy

- **`.ui-split`** — a two-track grid (`repeat(2, minmax(0, 1fr))`) with `--ui-split-gap` (`--ui-space-6`) and centered alignment.
- **`.ui-split-media`** — the media slot: hidden overflow, the theme small radius, and the border color; images and video render block and fluid.
- **`.ui-split-body`** — the copy slot; as the only child it spans the full row.
- **`.ui-split-eyebrow` / `.ui-split-title` / `.ui-split-copy`** — the body type slots using the label-lg, headline-sm, and body-lg typescales with `--ui-space-*` rhythm.
- **`.ui-split-action`** — the action slot, spaced `--ui-space-4` above the copy.

## Variants

- `ui-split` — the base two-column composition.
- With media — the default pairing of media and body.
- Body only — no media slot: the body spans the full row (`ui-split-body:only-child`).

## Anti-patterns

- Do not place two visuals in the media slot or two CTAs in the action slot: the split pairs exactly one visual with one action.
- Do not apply editorial typography inside the split template; the copy stays unstyled so consumers keep their prose contract.
- Do not add `direction`-specific layout rules: the columns already mirror through document direction.

## Checklist

- The section pairs at most one media item with one body block.
- The body copy is unstyled, ready for the consumer's prose typescale.
- The media image has descriptive alt text.
- The action slot holds one real link or button.

## Accessibility

Split stacks in source order on narrow screens, so assistive technology and reading order never diverge from the visual layout at any viewport. The media boundary and the body text keep `CanvasText` paint in forced-colors mode, and RTL documents mirror the columns without a separate layout override.
