# Feature card

Feature card is a composition of [Card](/components/card) and a call to action, used to highlight one offer, product, or article in a set people scan. Use it when a preview needs a media block, a short title, a supporting line, and one clear action — and all four should read as one unit. It is not a primitive: the wrapper reuses the `.ui-card` surface, so every visual signal (surface, shadow, focus, forced-colors) comes from Card itself.

## Examples

An elevated feature card with media, title, body copy, and a primary action.

<div class="specimen-block">
<article class="ui-card ui-card-elevated ui-feature-card">
  <div class="ui-feature-card-media">
    <img src="https://images.unsplash.com/photo-1494526585095-c41746248156?auto=format&fit=crop&w=800&q=60" alt="A landscape lit by a low sun" width="800" height="450" loading="lazy">
  </div>
  <div class="ui-feature-card-body">
    <h3 class="ui-card-title">Plan a weekend escape</h3>
    <p class="ui-card-body">Three quiet trails, one map, zero planning. Pick a route and go.</p>
    <div class="ui-card-action"><a class="ui-button ui-button-primary" href="/components/feature-card"><span>Explore routes</span></a></div>
  </div>
</article>
</div>

The specimen above is the live markup the template `feature-card.html` emits: an `<article>` carrying both the `ui-card` variant and the `ui-feature-card` wrapper, with the media block and body slots filled.

## Guidance

### When to use

Use a feature card when one item in a scannable set deserves a visual preview plus a single action — a promo, a highlight, a featured article. The media block earns its place when the image communicates the offer faster than the title alone.

### When not to use

Do not use a feature card for a static layout region that never acts or links — that is a plain [Card](/components/card) or a section heading. Do not stack every card in a set with feature treatments: when the items are homogeneous records to compare, a [List](/components/list) or [Data table](/components/data-table) scans better. The horizontal variant is deprecated upstream and is not part of the Gelium contract.

### Usability

- Compose on top of the Card surface: pick the card variant (`ui-card-elevated`, `ui-card-filled`, `ui-card-outlined`) the same way you would for any card.
- Give the media block a real image or video; the 16:9 ratio is structural geometry (literal `aspect-ratio`), never a token.
- Keep the CTA a real link or button in the `ui-card-action` slot so the action stays focusable and operable without JavaScript.

### Accessibility

- The card is a grouping container, never the sole carrier of meaning; the title, body, and action text carry the meaning.
- The media `img` needs descriptive alt text — the feature image is content, not decoration.
- In forced-colors mode the Card surface boundary keeps the grouping visible when color is removed.

See [Choose the right control](/docs/choose-the-right-control) for the cross-component decision.

## Anatomy

- **`.ui-feature-card`** — the wrapper: zero padding and hidden overflow so the media block sits flush and corners stay rounded. It must accompany a `ui-card` variant class on the same element.
- **`.ui-feature-card-media`** — a 16:9 media block (`aspect-ratio: 16 / 9`, literal and not tokenized) with cover-fill for images and video.
- **`.ui-feature-card-body`** — a flex column with `--ui-space-2` gaps and `--ui-space-4` padding, holding the title, body, and action slots.
- **`.ui-card-title` / `.ui-card-body` / `.ui-card-action`** — the Card slots reused unchanged; the template renders the title as `h3` inside the feature card.

## Anti-patterns

- Do not invent a fourth "horizontal" layout: it was deprecated upstream and is not part of the Gelium contract.
- Do not make the whole card one big link and then nest a second link inside it: nested interactive elements are invalid. Choose one action.
- Do not add scoped tokens for the media ratio; aspect-ratio stays literal (same rule as Video and Card media).

## Checklist

- The element carries both a `ui-card` variant and `ui-feature-card`.
- The media slot has a real asset with descriptive alt text.
- The action slot holds one real link or button.
- The card is one of a scannable set, not a lone static region.

## Accessibility

Feature card inherits everything from Card: the surface boundary survives forced-colors mode, `:focus-visible` appears on interactive roots, and the accessible name comes from the visible text. The media image is content and therefore needs descriptive alt text — never an empty `alt` unless the image is purely decorative.
