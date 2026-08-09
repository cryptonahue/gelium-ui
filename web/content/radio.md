# Radio

Radio is a native `input[type="radio"]` enhanced with CSS: the browser default look is removed and Loom paints the Material ring and the inner checked dot from the `--ui-radio-*` tokens. The real control stays the focusable, operable element, so form semantics, grouping, and assisted input work unchanged.

## Anatomy

- **Input** — the native `ui-radio` radio, `20px` in diameter (`--ui-radio-size`) with a fully round radius (`--ui-radio-radius`) and a `--ui-radio-outline-width` border in `--ui-radio-outline`. It keeps `appearance: none` only for styling; keyboard focus, the accessibility tree, and form submission are unaffected.
- **Mark** — `ui-radio-mark`, a `span` that sits over the ring and paints the `10px` inner dot in `--ui-radio-checked` when the input is selected.
- **Label** — `ui-radio-label`, the visible text paired with the ring. Put the accessible name here and wrap both in a `label` so clicking the text selects the input.
- **Error** — set `aria-invalid="true"` on the input to switch the ring to `--ui-field-error`.

## When to use it

Use a radio group when exactly one option from a mutually exclusive set must be chosen — never for independent yes/no options (that is a checkbox). Group radios by sharing the same `name` so the browser enforces the single-select rule. The label is never color-only: an error is always a visible ring change plus the surrounding form message.

## Accessibility

- Keep the native element: the radio keeps its role, name, value, checked state, and keyboard behavior at no cost.
- The input must be nested inside its `label` (or `id`-linked) so the accessible name always matches what is on screen.
- Group a related set with `fieldset` and `legend`: each radio stays a member of the same-named group, and the legend becomes the group's label.
- Never rely on color alone: selected state carries the inner dot, errors carry `aria-invalid`, and disabled state is announced by the platform.
- In forced-colors mode the ring keeps a `CanvasText` boundary, the selected dot stays on `CanvasText`, and the disabled dot drops to `GrayText`, so selection survives without color.