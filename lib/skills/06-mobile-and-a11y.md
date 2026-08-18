# Skill: Mobile and accessibility

Base guidelines so a screen is usable on any device and passes WCAG AA.

## Mobile

- Touch targets ≥ 44px (`--ui-touch-target`). Don't depend on hover.
- Contain overflow inside scroll children: child `min-width: 0` + parent
  `overflow-x: auto`. **Never** `overflow-x: hidden` on `body` — that masks
  content loss; scope the scroll to the container.
- Design for screen sizes (breakpoints ~40rem/48rem), not device names.
- Respect safe-area insets for phones with notches.
- Honor `prefers-reduced-motion` — animations off/`animation: none` for
  non-essential motion.

## Keyboard

- Full focus-visible contract: `--ui-color-focus-ring` visible on
  `:focus-visible` (never `outline: none` without a replacement).
- Logical tab order from the document order; no focus traps unless a modal
  (and then Esc closes it).
- Native controls give you free keyboard support — don't reinvent them.

## Semantics & ARIA

- Use native elements first (`<button>`, `<input>`, `<select>`, `<dialog>`,
  `<nav>`, `<article>`, `<section>`): they carry semantics + keyboard for free.
- ARIA only to bridge a genuine gap (e.g. `role="status"` for a live region,
  `aria-label` on an icon-only button with visible text nearby).
- Progress/loading: use `role="status"` / `aria-busy`, prefer semantic
  `<progress>`; name via visible caption where possible.

## Contrast & color

- Text meets WCAG AA (4.5:1; large ≥ 3:1). Use the theme's tokens.
- Never communicate meaning by color alone — pair with icon/text.
- Forced-colors mode: icons stay `currentColor`, borders use
  `CanvasText`/`Highlight` equivalents — do not rely on background fills.

## Media

- Images: meaningful `alt`, intrinsic `width`/`height`, `loading="lazy"` below
  the fold.
- Audio: native `audio[controls]` + typed sources + fallback text + optional
  transcript link.
- Video with meaningful audio: captions track + fallback.
- External embeds: explicit allowlist + consent boundary; otherwise canonical
  fallback.
- Reserve space for loading, provide recovery for error, explain empty.
