# Card

Card is a Material container that groups related content into one visual unit. Use a card when related information or actions belong together on a surface people scan — a preview, a profile, a summary. The group should read as one unit. The card itself stays a plain container: the semantic root is always a native element, chosen by what the card does.

## Guidance

### When to use

Use a card to group related content and actions into one visual unit — a preview, profile, or summary. People scan it as one unit.

### When not to use

Do not use a card for a single static layout region — that is a Panel. Do not wrap a simple click-to-open list of items in cards: a [List](/components/list) of links or a [Data table](/components/data-table) scans better for repeated rows. When the whole card navigates or acts, keep a real `<a>` or `<button>` root.

### Usability

- Pick one variant class on the same element: elevated (floats), filled (high emphasis), or outlined (quietest grouping).
- Choose the semantic root by what the card does: `<article>` for static content, `<a>` when it navigates, `<button>` when it performs an action.
- Keep interactive cards on a real anchor or button so they stay focusable and operable without JavaScript.

### Accessibility

- A card is a grouping container, never the sole carrier of meaning.
- Interactive cards must keep a real anchor or button root; their visible text provides the accessible name.
- In forced-colors mode the card keeps a `CanvasText` boundary so the grouping survives when color is removed.

See [Choose the right control](/docs/choose-the-right-control) for the cross-component decision.
## Alternative names

- Tile, container, panel, media card, info card.

## Agent prompt

Use Card to group related content and actions into one visual unit people scan — previews, profiles, summaries. The `.ui-card` base is a plain container: the semantic root is a native element chosen by what the card does. Use `<article>` for static content, `<a>` when the whole card navigates, `<button>` when it performs an action. Keep interactive cards on a real anchor or button so they stay focusable and operable without JavaScript. Don't use a card for a single static layout region — that is a Panel.

## Anatomy

The `.ui-card` base supplies layout and shape: a flex column, the theme corner radius `--ui-card-radius` (12 px), and the typescale body font. It never carries meaning by itself — content, headings, and actions inside the card do.

## Variants

Pick one variant class on the same element.

- `ui-card-elevated` — surface container with shadow level 1 (`--ui-shadow-1`), for cards that float above the page.
- `ui-card-filled` — a raised surface with no shadow, for high-emphasis groups.
- `ui-card-outlined` — a surface with a 1 px outline (`--ui-card-outline-color`), for the quietest grouping.

Each variant uses its own theme token family: `--ui-card-container-elevated`, `--ui-card-container-filled`, and `--ui-card-container-outlined`.

## Semantic roots

The element you write matters more than the visual.

- `<article class="ui-card ...">` — a static, self-contained group.
- `<a class="ui-card ..." href="...">` — the whole card navigates; the anchor carries the accessible name from its text.
- `<button class="ui-card ..." type="button">` — the whole card performs an action.

Interactive cards inherit `:focus-visible`, so keyboard users always see where focus is.

## Accessibility

A card is a grouping container, never the sole carrier of meaning. Interactive cards must keep a real anchor or button root, so they stay focusable and operable without JavaScript. Their visible text provides the accessible name. In forced-colors mode the card keeps a `CanvasText` boundary so the grouping survives when color is removed.