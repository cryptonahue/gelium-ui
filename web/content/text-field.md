# Text field

Text field is a single- or multi-line input built on native `<input>` and `<textarea>` elements. Use a text field when people must enter or edit free-form text — names, email addresses, search terms — and the value has to submit through a real `<form>`. Every visible label is explicitly associated through `for` and `id`, while helper and error messages use unique IDs referenced by `aria-describedby`.

## Guidance

### When to use

Use a text field when people must enter or edit free-form text — names, email addresses, search terms — and the value has to submit through a real `<form>`.

### When not to use

Do not use a text field when a constrained choice is faster: one option from a known list belongs in a [Select](/components/select) or [Radio](/components/radio) group, and an approximate magnitude belongs in a [Slider](/components/slider). For short, predefined answers, [Chips](/components/chips) can replace typing.

### Usability

- Single-line inputs use `<input>`; longer answers use the multi-line `<textarea>`.
- Every visible label is explicitly associated through `for` and `id`; helper and error messages use unique IDs referenced by `aria-describedby`.
- Server-side validation returns HTTP 422 with the submitted value preserved and `aria-invalid="true"` on the field.

### Accessibility

- The native `<input>`/`<textarea>` keeps its role, name, value, and keyboard behavior at no cost.
- Every visible label is explicitly associated through `for` and `id`, matching what is on screen.
- Errors are announced via `aria-invalid` plus visible error text — never color alone.
- Disabled keeps the native `disabled` attribute and takes precedence over error.

See [Choose the right control](/docs/choose-the-right-control) for the cross-component decision.
## Alternative names

- Input, text input, form field, entry field

## Agent prompt

Use Text field when people must enter or edit free-form text — names, email addresses, search terms — and the value has to submit through a real `<form>`. It is built on native `<input>` and `<textarea>` elements; every visible label is associated through `for`/`id`, helper and error messages are referenced by `aria-describedby`, and server-side validation returns HTTP 422 with the submitted value preserved, `aria-invalid="true"` and the `X-Gelium-Validation: true` header for the HTMX swap hook. Never build a custom text-entry widget: the native control is the contract, disabled keeps the native `disabled` attribute, and disabled takes precedence over error.

## Variants and states

The dogfooded examples demonstrate outlined and filled variants, normal and helper text, a visibly labelled error that does not rely on color, disabled inputs (filled and outlined), a disabled textarea, and an editable textarea.

Disabled fields keep the native `disabled` attribute: they leave the tab order, cannot receive focus or edits, and are excluded from hover and focus states. Disabled takes explicit precedence over error: a field flagged disabled never advertises `aria-invalid` or an error alert, even if both states are set in the view model. The CSS also resolves a combined `ui-text-field-error ui-text-field-disabled` wrapper to the disabled palettes instead of the error palette.

## How does server-side validation work?

The form below submits a named field to the server with browser validation disabled so the server remains the source of truth.

Whitespace-only input returns **HTTP 422** with the submitted value preserved, `aria-invalid="true"`, and an associated error. Valid input returns HTTP 200 with the preserved value and an accessible success message.

The form works without JavaScript: its native `method` and `action` submit to the same endpoint, which returns a complete documentation page with the updated form while preserving the HTTP 422 or 200 status. HTMX (the progressive-enhancement library that swaps server-rendered fragments) is an optional layer; when it sends `HX-Request: true`, the server returns only the complete form fragment and replaces it using `hx-target="this"` plus `hx-swap="outerHTML"`.

Expected HTMX validation failures return the explicit `X-Gelium-Validation: true` response header. Local `/static/app.js` listens for `htmx:beforeSwap` and only when both that header is `true` and the status is 422 does it set `shouldSwap` to true and `isError` to false. Other HTTP 422 responses remain transport errors and are not swapped by this hook.
