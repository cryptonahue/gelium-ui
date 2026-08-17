# Forms contract

Input rules come **before** component markup. A form is a labeled control group that submits real data. Pick the control with [Choose the right control](/docs/choose-the-right-control), then apply this contract on every field.

## Labels

- Put the **label above** the field. Never use the placeholder as the only label.
- Associate every label with its control: `label[for]` + matching `id`, or wrap the control in the `label`.
- Placeholders may hint at format. They must not replace the visible name of the field.

## `type` and `inputmode`

Pair native `type` with `inputmode` so mobile keyboards match the data. Prefer the native type when the browser already validates it.

| Data | `type` | `inputmode` | Notes |
|---|---|---|---|
| Email | `email` | `email` | Native format check; keep a visible label |
| Phone | `tel` | `tel` | Do not invent a custom dial pad |
| Whole numbers | `text` or `number` | `numeric` | Prefer `numeric` when symbols are not needed |
| Decimals / money | `text` | `decimal` | Avoid `type="number"` for money (locale and spin issues) |
| URL | `url` | `url` | Full URL entry |
| Search | `search` | `search` | Clear affordance on supporting browsers |

Do not invent a custom input when a native control already fits. Use [Text field](/components/text-field), [Checkbox](/components/checkbox), [Radio](/components/radio), [Select](/components/select), [Switch](/components/switch), or [Slider](/components/slider).

## Autocomplete

When the field collects **identity or contact** data, set a standard `autocomplete` token. That lets the browser and password managers fill correctly.

Examples: `name`, `email`, `tel`, `username`, `current-password`, `new-password`, `street-address`, `postal-code`, `country`, `cc-number`, `one-time-code`.

Omit `autocomplete` only when the value is not personal or not reusable (for example a one-off filter that is not an identity field).

## Validation timing

- Validate **after the user interacts** with the field (blur) and **on submit**.
- Do **not** show errors on every keystroke by default. Live checking mid-word is noisy and blocks completion.
- Server validation remains authoritative: HTTP **422** with `aria-invalid` and field messages — see [Server contracts](/docs/server-contracts).
- Error copy follows the **action pattern** (what to do next), not blame — see [Content style](/docs/content-style).

## Disabled vs error

- **Disabled wins over error.** A disabled control must not show `aria-invalid` or an error message.
- Disabled uses the native `disabled` attribute. It leaves the tab order and is not submitted.
- Never disable a field only to “look” invalid. Show the error state instead so the user can fix it.

## Paste, touch, and targets

- **Do not block paste.** Password managers and mobile users rely on paste. Never strip `paste` or force character-by-character entry without a strong security reason documented on the field.
- Touch targets use the existing size tokens (`--ui-size-touch` and related spacing). Do not invent one-off hit areas per form.

## What this page is not

This page is the **cross-field contract**. Anatomy, variants, and dogfood demos live on each component page. Control choice lives on [Choose the right control](/docs/choose-the-right-control). HTTP and HTMX wiring live on [Server contracts](/docs/server-contracts).
