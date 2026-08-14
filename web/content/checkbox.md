# Checkbox

Checkbox is a native `input[type="checkbox"]` enhanced with CSS: the browser default look is removed and Gelium paints the Material box, its outline, and the checked chevron from the `--ui-checkbox-*` tokens. The real control stays the focusable, operable element, so form semantics and assisted input work unchanged.

## Anatomy

- **Input** — the native `ui-checkbox` checkbox, `18px` square (`--ui-checkbox-size`) with a Material radius (`--ui-checkbox-radius`) and a `--ui-checkbox-outline-width` border in `--ui-checkbox-outline`. It keeps `appearance: none` only for styling; keyboard focus, the accessibility tree, and form submission are unaffected.
- **Mark** — `ui-checkbox-mark`, a `span` that sits over the box and paints the chevron in `--ui-checkbox-icon` when the input is checked.
- **Label** — `ui-checkbox-label`, the visible text paired with the box. Put the accessible name here and wrap both in a `label` so clicking the text toggles the input.
- **Error** — set `aria-invalid="true"` on the input to switch the outline to `--ui-checkbox-error`.

## When to use it

Use a checkbox when one option can be toggled independently, or when several options multiply across — never for a single mutually exclusive decision (that is a radio). The label is never color-only: an error is always a visible outline change plus the surrounding form message.

## Accessibility

- Keep the native element: the checkbox keeps its role, name, value, checked state, and keyboard behavior at no cost.
- The input must be `id`-linked to its `label` (or nested inside it) so the accessible name always matches what is on screen.
- Never rely on color alone: checked state carries the chevron, errors carry `aria-invalid`, and disabled state is announced by the platform.
- In forced-colors mode the box keeps a `CanvasText` boundary and the chevron stays on `Canvas`, so the checked state survives without color.