# Icon

Icon is the Gelium-only primitive for trusted inline SVG glyphs. All icons are resolved server-side as trusted `template.HTML` markup — icon markup must never be built from user input.

## The `.ui-icon` utility

The shared `.ui-icon` class gives every icon consistent 24 px sizing, a fixed flex size, and a `currentColor` fill so the glyph matches surrounding text automatically.

## Accessibility contract

Decorative icons must be `aria-hidden="true"` and `focusable="false"` and must sit next to visible text or another accessible name. An icon that carries meaning is named by visible text beside it and never by an `aria-label` alone. The examples below follow that contract.