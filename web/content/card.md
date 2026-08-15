# Card

Card is a Material container that groups related content into one visual unit. Use a card when related information or actions belong together on a surface people scan — a preview, a profile, a summary — and the group should read as one unit. The card itself stays a plain container: the semantic root is always a native element, chosen by what the card does.

## Anatomy

The `.ui-card` base supplies layout and shape: a flex column, the theme corner radius `--ui-card-radius` (12 px), and the typescale body font. It never carries meaning by itself — content, headings, and actions inside the card do.

## Variants

Pick one variant class on the same element:

- `ui-card-elevated` — surface container with shadow level 1 (`--ui-shadow-1`), for cards that float above the page
- `ui-card-filled` — a raised surface with no shadow, for high-emphasis groups
- `ui-card-outlined` — a surface with a 1 px outline (`--ui-card-outline-color`), for the quietest grouping

Each variant uses its own theme token family: `--ui-card-container-elevated`, `--ui-card-container-filled`, and `--ui-card-container-outlined`.

## Semantic roots

The element you write matters more than the visual:

- `<article class="ui-card ...">` — a static, self-contained group
- `<a class="ui-card ..." href="...">` — the whole card navigates; the anchor carries the accessible name from its text
- `<button class="ui-card ..." type="button">` — the whole card performs an action

Interactive cards inherit `:focus-visible`, so keyboard users always see where focus is.

## Accessibility

A card is a grouping container, never the sole carrier of meaning. Interactive cards must keep a real anchor or button root so they are focusable and operable without JavaScript, and their visible text provides the accessible name. In forced-colors mode the card keeps a `CanvasText` boundary so the grouping survives when color is removed.