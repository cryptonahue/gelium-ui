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
- The tri-state `indeterminate` visual (a dash instead of the chevron) is pure CSS. It paints whenever the platform sets the input's `indeterminate` IDL property. No HTML attribute exists for it — the browser or script sets the property, for example in a partial [Chip](/components/chips) filter selection.

### Accessibility

- Keep the native element: the checkbox keeps its role, name, value, checked state, and keyboard behavior at no cost.
- The input must be `id`-linked to its `label` (or nested inside it) so the accessible name matches what is on screen.
- Checked state carries the chevron, errors carry `aria-invalid`, and disabled state is announced by the platform. For a tri-state control, set `aria-checked="mixed"` alongside the `indeterminate` property so assistive tech announces the partial state.
- In forced-colors mode the box keeps a `CanvasText` boundary and the chevron stays on `Canvas`. Indeterminate mirrors checked: the filled box is `CanvasText` with the dash on `Canvas`.

See [Choose the right control](/docs/choose-the-right-control) for the cross-component decision.
## Anatomy

- **Input** — the native `ui-checkbox` checkbox, `18px` square (`--ui-checkbox-size`) with a Material radius (`--ui-checkbox-radius`) and a `--ui-checkbox-outline-width` border in `--ui-checkbox-outline`. It keeps `appearance: none` only for styling; keyboard focus, the accessibility tree, and form submission are unaffected.
- **Mark** — `ui-checkbox-mark`, a `span` that sits over the box and paints the chevron in `--ui-checkbox-icon` when the input is checked. On `:indeterminate` the mark swaps the chevron for a `10×2` centered dash (`--ui-checkbox-indeterminate-width` × `--ui-checkbox-indeterminate-height`).
- **Label** — `ui-checkbox-label`, the visible text paired with the box. Put the accessible name here and wrap both in a `label` so clicking the text toggles the input.
- **Error** — set `aria-invalid="true"` on the input to switch the outline to `--ui-checkbox-error`.

## States

Checkbox covers `rest`, `hover` (outline grows to `--ui-checkbox-hover-outline`), `focus` (`:focus-visible` ring), `pressed` (`:active` scale), `checked`, `indeterminate`, `disabled`, and `error` (`aria-invalid`). The indeterminate state styles the native box whenever the platform sets the `indeterminate` IDL property, a genuine tri-state. Disabled and indeterminate repaint through the checked-disabled container.

## Accessibility

- Keep the native element: the checkbox keeps its role, name, value, checked state, and keyboard behavior at no cost.
- The input must be `id`-linked to its `label` (or nested inside it) so the accessible name always matches what is on screen.
- Never rely on color alone: checked carries the chevron, indeterminate the dash, errors `aria-invalid`, disabled the platform announcement.
- For tri-state controls the dash is a visual plus `aria-checked="mixed"`; the pair announces and shows the partial state without color.
- In forced-colors mode the box keeps a `CanvasText` boundary and the chevron stays on `Canvas`. Indeterminate mirrors that pair.

## Specimen

The painted box has three visual states. No HTML attribute sets the dash — the browser or script flips the input's `indeterminate` property. Until then the sample renders as unchecked. The `aria-checked="mixed"` attribute keeps assistive tech honest about the tri-state.

<div class="specimen-block">
<p><label class="ui-checkbox"><input type="checkbox"><span class="ui-checkbox-mark"></span><span class="ui-checkbox-label">Unchecked</span></label></p>
</div>

Next, checked.

<div class="specimen-block">
<p><label class="ui-checkbox"><input type="checkbox" checked><span class="ui-checkbox-mark"></span><span class="ui-checkbox-label">Checked</span></label></p>
</div>

Last, indeterminate.

<div class="specimen-block">
<p><label class="ui-checkbox"><input type="checkbox" aria-checked="mixed"><span class="ui-checkbox-mark"></span><span class="ui-checkbox-label">Indeterminate</span></label></p>
</div>