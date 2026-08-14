# Gelium UI — Wire Compatibility Strategy

> The product and repository identity is **Gelium UI** (Go module `geliumui`, npm
> package `gelium-ui`, CLI/binary `gelium`). The **wire contracts** below keep
> their legacy `loom:*` / `X-Loom-*` names on purpose and are frozen.

## Scope

This document records the deliberate compatibility decision for the Gelium
rename: branding, module, package, and CLI names move to Gelium, while the
server-driven wire contracts that existing consumers depend on are preserved
verbatim. It is the implementation of `openspec/config.yaml`: "Preserve product
brand Gelium UI; do not rename `loom:*` / `X-Loom-*` wire prefixes without an
explicit breaking-change proposal."

## Frozen legacy wire contracts (retained)

Renaming any of these would break existing consumers (the HTMX swap hook, the
client toast listener, or any server integration that already reads them), so
they are retained as the canonical wire contract:

| Contract | Where it lives | Why it is frozen |
|---|---|---|
| `X-Loom-Validation: true` response header | `internal/app/text_field.go`, `select.go`, `newsletter.go`, `recipe_admin_resource.go`, `toast.go`; consumed by `web/static/app.js` (`htmx:beforeSwap` only swaps 422 responses that carry this header) | External consumers and the served client hook depend on the exact header name |
| `HX-Trigger: {"loom:toast":{...}}` event | `internal/app/toast.go`, `data_table.go`, `recipe_*.go`; consumed by the `loom:toast` listener in `web/static/app.js` | HTMX triggers a DOM event named `loom:toast`; the payload key is the event name |
| `#loom-toast-region` live region id | `web/templates/toast.html` | Part of the toast client contract documented in `gelium-ui-core.md` and the accessibility contract; tests pin it |
| `data-loom-toast-dismiss` / `data-loom-toast-done` attributes | `web/templates/toast.html`, `web/static/app.js` | Internal enhancement-layer contract that ships with the served `app.js`; kept consistent with `loom:toast` |
| Schema enums `422_X_LOOM_VALIDATION`, `LOOM_TOAST` | `docs/gelium/schemas/*.json`, `docs/gelium/schemas/screen-recipe.schema.json` | Machine-readable references to the frozen wire contracts |

## Policy

1. **Legacy wire names are canonical.** Do not rename `loom:*` / `X-Loom-*`
   wire prefixes without an explicit breaking-change proposal that migrates
   consumers and the served client hook together.
2. **New public contracts use Gelium names.** Any contract introduced from now
   on (new headers, events, attributes, ids) uses the `gelium:*` / `X-Gelium-*`
   / `gelium-*` convention.
3. **Branding is Gelium.** Human-facing surfaces, repository documentation,
   prompts, LICENSE, and source identity use "Gelium UI". The physical repo
   path (`loom-ui`) and historical audit documents keep their original names as
   factual/historical references only.

## Tests that pin the frozen contracts

- `internal/app/text_field_test.go`, `select_test.go`, `newsletter_test.go`,
  `recipe_admin_resource_test.go` — assert `X-Loom-Validation: true` on 422.
- `internal/app/toast_test.go`, `data_table_test.go`, `recipe_*.go` tests —
  assert the `HX-Trigger` payload key `loom:toast` and the `#loom-toast-region`.
- `internal/app/wire_compat_test.go` — explicit regression guard that re-checks
  the frozen wire contract so the rename is never silently reverted.
- `web/styles_contract_test.go` — forbids persistent partials from emitting the
  transient `loom:toast` signal.
