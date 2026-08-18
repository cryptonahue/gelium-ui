# Skill: Server contracts

Gelium renders server-side; state lives in the URL and the response. These are the
wire contracts a consumer must honor.

## Read

- Collection/list state (filter, sort, page, selection) is **GET query params**
  on a stable URL.

## Mutate (create, update, delete, subscribe, like, react)

- **POST** to the resource action.
- On success → **303** redirect to a safe GET (the list or the created item).
- On failure → **422** re-render of the form/action surface with
  `X-Gelium-Validation: true` and a validation summary + inline field errors.
- Add the **`gelium:toast`** HTMX trigger header only for transient,
  auto-dismissable feedback (flash after a redirect).

Example header for a transient toast after redirect:

```
HX-Trigger: {"gelium:toast": {"kind": "success", "text": "Saved"}}
```

## Read-only vs mutating

`<a>` / `GET` only for safe, read-only navigation. Every state-changing action is
a form `POST` (never a `GET` link or a bare `hx-get` that mutates).

## No-JS is the contract

With `js/gelium.js` disabled the same POST+303 and 422 re-render still work —
HTML forms are the mechanism, HTMX only enhances. If a flow needs JS to function,
you have drifted from the contract.

## Server is the authority

Optimistic UI, if any, is enhancement on top of a server-authoritative round trip
— never the source of truth.
