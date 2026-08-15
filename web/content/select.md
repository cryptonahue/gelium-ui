# Select

Select is a dropdown picker built on the native [`<select>`](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/select) element with Material 3 field styling — the browser element stays the focusable, operable control. Use a select when people must choose one option from a known list and the choice submits through a normal `<form>`; the native element keeps the no-JS flow and keyboard behavior for free. No component JavaScript is required.

## Usage

```html
<div class="ui-select ui-select-filled">
  <select id="plan" name="plan">
<option value="" selected disabled>Choose a plan</option>
    <option value="standard">Standard</option>
    <option value="priority">Priority</option>
  </select>
  <label class="ui-select-label" for="plan">Plan</label>
  <span class="ui-select-caret" aria-hidden="true">
    <svg focusable="false" viewBox="0 0 24 24" width="20" height="20" fill="currentColor">
      <path d="m7 10 5 5 5-5z"></path>
    </svg>
  </span>
</div>
```

Wrap with `ui-select-filled` or `ui-select-outlined` to pick the variant.

## Variants

- **Filled** (`ui-select-filled`) — surface-container-highest background, top-only 4px radius, and a bottom active indicator. Focus turns the indicator to the primary color.
- **Outlined** (`ui-select-outlined`) — transparent field with a 1px outline; focus grows the outline to the primary color.

## States

- **Normal** — a floating label sits centered until the field is focused or a value is chosen, then it floats to the top.
- **Populated** — the label floats whenever a real option is selected (`:has(select option:checked:not([value=""]))`), even without focus.
- **Disabled** — add `disabled` to the native element. The field dims its content and stops interaction.
- **Error** — add `aria-invalid="true"` to the select (plus an error message in the form, beyond color). The outline and label turn to the error color.

The label doubles as the placeholder. The first `<option value="" selected disabled>` keeps its visible picker text (for example `Choose a plan`) so the native popup row reads clearly, while the field paints that text `transparent` until a real option is chosen — the resting label is the only prompt in the closed field. The label floats once the field is focused or a real option is selected.

## Design tokens

Field anatomy reuses the `--ui-field-*` palette tokens shared with [Text field](/components/text-field). Select-specific tokens are:

| Token | Role |
| --- | --- |
| `--ui-select-height` | Field height (56px) |
| `--ui-select-radius` | Outlined field corner radius |
| `--ui-select-radius-top` | Filled field top-only corner radius |
| `--ui-select-caret` | Caret icon color |
| `--ui-select-fg` | Selected value text color |
| `--ui-select-label` | Resting floating label color |
| `--ui-select-outline` | Resting outline / active indicator |
| `--ui-select-container-filled` | Filled field background |
| `--ui-select-hover-outline` | Hover outline color |
| `--ui-select-focus` | Focus outline / indicator color |
| `--ui-select-error` | Error outline and label color |
| `--ui-select-disabled-opacity` | Disabled content opacity |
| `--ui-select-list-bg` | Native options popup background |
| `--ui-select-list-fg` | Native options popup text |

## Does select work without JavaScript?

The control is the native `<select>`: keyboard navigation, form submission, and the options popup are browser behavior — zero component JavaScript. The floating label and field surface are pure CSS. In forced-colors mode the select keeps visible `CanvasText` / `GrayText` / `Mark` boundaries and gives the popup `FieldText` contrast.

## Accessibility

- The visible `<label class="ui-select-label" for="...">` is a real label that focuses the select.
- The native element carries name/value semantics for form submission and assistive technology.
- Focus is announced with the shared `--ui-focus-thickness` outline; error state is announced via `aria-invalid` plus visible error text, never color alone.
- A picker's placeholder option is `disabled` so assistive tech and users never confuse it with a real selection.