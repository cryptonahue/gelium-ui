# Badge

Badge is a compact marker that shows a dot, a count, or a short label on the error palette. Use a badge when a small piece of information — a count, a status dot, a short label — must sit on or near another element without taking over the row. It is never color-only: a badge mark must always carry or accompany visible text so the information survives without color.

## Anatomy

- **Dot** — `ui-badge`, a `6px` diameter disc (the `--ui-badge-size` token). It is decorative and must be paired with adjacent visible text; use `aria-hidden="true"` on the dot since the text supplies the meaning.
- **Large** — `ui-badge ui-badge-large`, a pill with `min-width: --ui-badge-large-size` (16 px) and the label inside. Put the count or short label as its text content directly.

Both variants use the theme tokens `--ui-badge-container` (error) and `--ui-badge-fg` (on-error), so they stay consistent across light and dark schemes. The capsule shape comes from `--ui-radius-full`.

## When to use it

A dot badge suits a live "indicator" whose meaning lives in the adjacent text. A large badge is better for a numeric count or a one-word state that should be readable at a glance. If the badge is the only signal, it fails — the surrounding text or a `sr-only` label must carry the same information.

## Accessibility

- The dot is decorative: keep it `aria-hidden` and put its meaning in visible text nearby.
- The large badge reads its own text content; it needs no extra label for its own value.
- In forced-colors mode the badge keeps a `CanvasText` boundary so the marker remains visible when color is removed.