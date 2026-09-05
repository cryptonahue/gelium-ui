# Admin recipes

Gelium UI provides server-rendered admin recipes for applications that own their data, identity, authorization, and domain rules. Start with a runnable recipe, then adapt its contract to the consumer instead of copying framework-specific APIs.

## Runnable recipes

- [Admin Resource](/recipes/admin-resource) — list, search, filters, pagination, selection, forms, confirmation, and mutations.
- [Ops Queue](/recipes/ops-queue) — operational queue with search, filters, bounded CSV export, POST + 303 transitions, and recovery states.
- [Public/Social Feed](/recipes/public-feed) — server-rendered feed with pagination, reactions, loading, and empty states.

## Contract guides

- [Application integration](/docs/application-integration) — consumer-owned authentication, authorization, audit, and tenant isolation.
- [Dashboard metrics](/docs/dashboard-metrics) — source, period, unit, freshness, delta, and metric visibility.
- [Feedback recipe](/docs/feedback-recipe) — Toast, Banner, Inline alert, validation summary, Error state, and Empty state.
- [Recipe testing](/docs/recipe-testing) — HTTP, authorization, tenant scope, no-JS, HTMX, accessibility, and responsive checks.
- [Export recipe](/docs/export-recipe) — allowlists, limits, sync/async boundaries, downloads, and security.
- [Import recipe](/docs/import-recipe) — upload, mapping, validation, partial failure, and job states.
- [Relationships](/docs/relationships-recipe) — nested URLs, independent authorization, tenant scope, and related-record states.
- [Extensibility](/docs/extensibility) — recipe-first extension order and component/theme boundaries.

## Ownership rule

The consumer application owns authentication, authorization, tenant resolution, data queries, domain mutations, audit persistence, job execution, and business outcomes. Gelium renders the supplied view model and provides HTML, state, feedback, and accessibility contracts.

Never treat a hidden action, disabled control, tenant parameter, or client-side state as a security boundary. Re-check authorization and tenant scope at the mutation, job, export, import, download, and related-record boundaries.

## Server-first baseline

- Use GET for reads and navigable state.
- Use POST + 303 for mutations when navigation follows success.
- Use 422 with preserved values and field associations for validation.
- Keep primary navigation, forms, recovery, and status pages usable without JavaScript.
- Add HTMX only as an enhancement that preserves URL, authorization, and mutation semantics.

## Choosing the next extension

Use an existing recipe first. Add a new component only when the required anatomy or state cannot be composed from registered components. Do not add runtime plugins, framework adapters, database assumptions, or authentication systems to the central Gelium contract without a real consumer and an explicit architecture decision.

## Related

- [Patterns](/docs/patterns)
- [Server contracts](/docs/server-contracts)
- [Accessibility](/docs/accessibility)
