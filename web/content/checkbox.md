# Checkbox

Checkbox is a native `input[type="checkbox"]` painted with the Material box, outline, and checked chevron from the `--ui-checkbox-*` tokens. The real control stays the focusable, operable element. Use a checkbox when an option toggles independently or several options multiply across a list. The selection must submit through a real `<form>`. Form semantics and assisted input work unchanged.

## Guidance

### When to use

Use a checkbox when an option toggles independently or several options multiply across a list. The selection must submit through a real `<form>`.

### When not to use

Never use a checkbox for a single mutually exclusive decision — that is a [Radio](/components/radio). When exactly one of many options must be picked and the list is long, a [Select](/components/select) collapses it. For one independent on/off setting that takes effect immediately, prefer a [Switch](/components/switch).

### Usability

- Wrap the input and its text in a `label` so clicking the text toggles the box.
- Give every checkbox a `name` and `value` so the selection submits with the form.
- Errors are a visible outline change (`aria-invalid`) plus the surrounding form message — never color alone.

### Accessibility

- Keep the native element: the checkbox keeps its role, name, value, checked state, and keyboard behavior at no cost.
- The input must be `id`-linked to its `label` (or nested inside it) so the accessible name matches what is on screen.
- Checked state carries the chevron, errors carry `aria-invalid`, and disabled state is announced by the platform.
- In forced-colors mode the box keeps a `CanvasText` boundary and the chevron stays on `Canvas`, so the checked state survives without color.

See [Choose the right control](/docs/choose-the-right-control) for the cross-component decision.
## Anatomy

- **Input** — the native `ui-checkbox` checkbox, `18px` square (`--ui-checkbox-size`) with a Material radius (`--ui-checkbox-radius`) and a `--ui-checkbox-outline-width` border in `--ui-checkbox-outline`. It keeps `appearance: none` only for styling; keyboard focus, the accessibility tree, and form submission are unaffected.
- **Mark** — `ui-checkbox-mark`, a `span` that sits over the box and paints the chevron in `--ui-checkbox-icon` when the input is checked.
- **Label** — `ui-checkbox-label`, the visible text paired with the box. Put the accessible name here and wrap both in a `label` so clicking the text toggles the input.
- **Error** — set `aria-invalid="true"` on the input to switch the outline to `--ui-checkbox-error`.

## Accessibility

- Keep the native element: the checkbox keeps its role, name, value, checked state, and keyboard behavior at no cost.
- The input must be `id`-linked to its `label` (or nested inside it) so the accessible name always matches what is on screen.
- Never rely on color alone: checked state carries the chevron, errors carry `aria-invalid`, and disabled state is announced by the platform.
- In forced-colors mode the box keeps a `CanvasText` boundary and the chevron stays on `Canvas`, so the checked state survives without color.