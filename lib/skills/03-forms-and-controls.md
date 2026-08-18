# Skill: Forms and controls

Server-rendered forms follow a strict contract. Choice of control comes before
layout.

## What control for what

| Need | Use | Not |
|---|---|---|
| One of a small set | radio group (`fieldset`+`legend`) | menu |
| Many options, single | native `<select>` | custom menu |
| On/off, boolean | `<input type="checkbox">` or switch | select |
| A range / a level | slider | text field |
| Free text | text field (matching `type`) | div contenteditable |
| List of actions | menu / button group | a generic list |

## Form layout contract

- **Label above** the field, visible, never placeholder-as-label alone.
- Use correct `type` + `inputmode` (email, url, tel, number, numeric, search).
- `autocomplete` for known fields (name, email, postal-code, current-password…).
- Validate after interaction, not on every keystroke.
- Group radios/checkboxes in `fieldset` with `legend`.

## Validation (server, 0-JS)

- Invalid submit → **422** with header `X-Gelium-Validation: true`.
- Re-render the form with a **validation-summary** + inline field errors
  (no hiding meaning behind color alone).
- `js/gelium.js` swaps the 422 response into the form; without JS the 422 page
  renders the same summary natively.

## POST + 303

Mutating actions (create, update, delete, subscribe, like) are **POST + 303**
redirect back to a safe GET. Persistent success is a banner/result page, not a
toast. Transient, auto-dismissable feedback may use `gelium:toast`.
