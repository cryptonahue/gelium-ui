# Server contracts

Every Gelium component is a server-rendered HTML partial with a documented HTTP contract: GET with stable query state, POST + 303 for mutations, 422 for validation, and HX-Trigger for transient feedback.

## GET + stable query parameters

List state lives in the URL: `?q=&sort=&dir=&page=&selection=`. The URL is the state — without JavaScript a full reload renders the same page; with HTMX the handler bifurcates on `HX-Request: true` and swaps an `outerHTML` fragment instead. Vocabularies are closed and sanitized, and parameter order is stable. [Data table](/components/data-table) is the canonical consumer.

## POST + 303 See Other

Mutations move state with a plain form POST and a `303 See Other` redirect — no fragments, no client-side orchestration. The destination page re-renders, carrying any persistent success or error surface. The [Admin Resource](/recipes/admin-resource) and [Ops Queue](/recipes/ops-queue) screen recipes build on this contract.

## 422 validation

Server-side field validation answers `422` with `X-Gelium-Validation: true`. The response re-renders the form preserving the submitted value, marking each field `aria-invalid` and associating the message via `aria-describedby`. Without JavaScript the full page re-renders; with HTMX the local `app.js` hook swaps only the 422 fragment (`htmx:before:swap`, using `event.detail.ctx.response`). Validation never fires a toast — errors stay tied to their context. [Text field](/components/text-field) and [Select](/components/select) demonstrate the contract, and the [Dialog](/components/dialog) form flows use it too.

## HX-Trigger toast feedback

Transient feedback is server-driven: the response carries `HX-Trigger: {"gelium:toast":{"type":"success","message":"…"}}`. The vocabulary is closed — `info | success | warning | error` — and `error` renders with `role="alert"` inside the `#gelium-toast-region` live region. Without JavaScript the toast renders as a persistent inline region; with HTMX it auto-dismisses (4s/8s, pausable). [Toast](/components/toast) documents the full pattern.

## Contract rules

- Validation errors are never toasts, and persistent success is never an `HX-Trigger` toast — persistent feedback uses the banner / inline-alert pattern instead.
- The `gelium:*` prefixes and `X-Gelium-*` headers are the canonical wire contract: reference them, never invent new ones.
- Every JavaScript enhancement has a real no-JS fallback; the main flow never breaks without JavaScript.
