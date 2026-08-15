# Gelium UI — Wire Contract Reference

> The product and repository identity is **Gelium UI** (Go module `geliumui`, npm
> package `gelium-ui`, CLI/binary `gelium`). The **wire contracts** below use the
> product's `gelium:*` / `X-Gelium-*` prefixes and are canonical.

## Scope

This document is the canonical reference for the server-driven wire contracts of
Gelium UI. The wire prefixes were migrated from their original legacy names to
`gelium:*` / `X-Gelium-*` on **2026-08-15** by owner
decision: the project is brand new (v0.4.0) with no external clients, so the
rename happens now instead of freezing legacy names. There are no legacy
wire references left in code, templates, JS, tests, schemas, or docs.

## Canonical wire contracts

| Contract | Where it lives | Notes |
|---|---|---|
| `X-Gelium-Validation: true` response header | `internal/app/text_field.go`, `select.go`, `newsletter.go`, `recipe_admin_resource.go`, `toast.go`; consumed by `web/static/app.js` (`htmx:beforeSwap` only swaps 422 responses that carry this header) | Canonical validation contract; the served client hook depends on the exact header name |
| `HX-Trigger: {"gelium:toast":{...}}` event | `internal/app/toast.go`, `data_table.go`, `recipe_*.go`; consumed by the `gelium:toast` listener in `web/static/app.js` | HTMX triggers a DOM event named `gelium:toast`; the payload key is the event name |
| `#gelium-toast-region` live region id | `web/templates/toast.html` | Part of the toast client contract documented in `gelium-ui-core.md` and the accessibility contract; tests pin it |
| `data-gelium-toast-dismiss` / `data-gelium-toast-done` attributes | `web/templates/toast.html`, `web/static/app.js` | Internal enhancement-layer contract that ships with the served `app.js`; kept consistent with `gelium:toast` |
| Schema enums `422_X_GELIUM_VALIDATION`, `GELIUM_TOAST` | `docs/gelium/schemas/*.json`, `docs/gelium/schemas/screen-recipe.schema.json` | Machine-readable references to the canonical wire contracts |

## Policy

1. **The gelium wire names are canonical.** Any new contract (headers, events,
   attributes, ids) uses the `gelium:*` / `X-Gelium-*` / `gelium-*` convention.
   Never introduce a legacy loom-branded name.
2. **Branding is Gelium.** Human-facing surfaces, repository documentation,
   prompts, LICENSE, and source identity use "Gelium UI". The physical repo
   path (`loom-ui`) and historical audit documents keep their original names as
   factual/historical references only.

## Historical note (migration of 2026-08-15)

The wire contracts originally carried a legacy loom-branded prefix (event names,
response headers, live-region id, data attributes, and the schema enums) and
were frozen as a legacy-compatibility contract. On 2026-08-15 the owner decided
to migrate them to the product name: the project is brand new with no external
clients ("Hay que cambiarlo igual a eso que dice loom obvio no romper nada.
Porque esto recien nace sino cuando."). The migration was done in atomic
RED→GREEN pairs on branch `wire-prefix-gelium`; all occurrences in Go, JS,
templates, tests, schemas, and docs now use the gelium names, and
`wire_compat_test.go` re-pins the canonical contracts. There are no residual
legacy wire names left anywhere in the tree.

## Tests that pin the canonical contracts

- `internal/app/text_field_test.go`, `select_test.go`, `newsletter_test.go`,
  `recipe_admin_resource_test.go` — assert `X-Gelium-Validation: true` on 422.
- `internal/app/toast_test.go`, `data_table_test.go`, `recipe_*.go` tests —
  assert the `HX-Trigger` payload key `gelium:toast` and the `#gelium-toast-region`.
- `internal/app/wire_compat_test.go` — explicit regression guard that re-checks
  the canonical wire contract so it is never silently reverted.
- `web/styles_contract_test.go` — forbids persistent partials from emitting the
  transient `gelium:toast` signal.
