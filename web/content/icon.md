# Icon

Icon is the Gelium-only primitive for trusted inline SVG glyphs. Use icons when a control or a message needs a small, themeable graphic that ships with the server render. All icons are resolved server-side as trusted `template.HTML` markup — icon markup must never be built from user input.

## Guidance

### When to use

Use icons when a control or a message needs a small, themeable graphic that ships with the server render.

### When not to use

Do not use an icon where text alone is clearer, and never make an icon the only signal: decorative icons must sit next to visible text, and meaningful icons must be named by visible text. Never build icon markup from user input — glyphs are trusted, server-resolved markup.

### Usability

- The shared `.ui-icon` class gives every icon consistent 24 px sizing and a `currentColor` fill so the glyph matches surrounding text automatically.
- Icons are resolved server-side as trusted `template.HTML` markup; pass only internal glyphs.
- Decorative icons must include `aria-hidden="true"` and `focusable="false"`.

### Accessibility

- Decorative icons must be `aria-hidden="true"` and `focusable="false"` and must sit next to visible text or another accessible name.
- An icon that carries meaning is named by visible text beside it, never by an `aria-label` alone.
- The `currentColor` fill keeps icons at the surrounding text's contrast in every theme, including forced-colors mode.

See [Choose the right control](/docs/choose-the-right-control) for the cross-component decision.
## The `.ui-icon` utility

The shared `.ui-icon` class gives every icon consistent 24 px sizing, a fixed flex size, and a `currentColor` fill so the glyph matches surrounding text automatically.

## Accessibility contract

Decorative icons must be `aria-hidden="true"` and `focusable="false"` and must sit next to visible text or another accessible name. An icon that carries meaning is named by visible text beside it and never by an `aria-label` alone. The examples below follow that contract.