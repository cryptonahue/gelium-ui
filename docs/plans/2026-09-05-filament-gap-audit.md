# Filament-inspired capability audit

This audit compares Gelium UI's current roadmap with the public Filament documentation at <https://filamentphp.com/docs>. The goal is capability coverage, not framework parity: Filament's Laravel, Livewire, Eloquent, queue, and package APIs are not Gelium contracts.

## Executive answer

Gelium has the core server-rendered resource foundation: list/table, search, filters, pagination, selection, bulk action, detail, forms, dashboard metric contract, feedback states, and recipe extensibility.

The meaningful gaps are documented below. The next product-facing slices should be selected only when a consumer has a real use case.

## Coverage matrix

| Filament capability area | Gelium status | Decision |
|---|---|---|
| Resource listing, tables, columns, sorting, pagination | Covered | Existing Data Table and Admin Resource recipes |
| Search and filters | Covered for current recipes | GET contracts exist; extend only for a real filter requirement |
| Record view/detail and infolists | Covered at recipe level | Do not extract a generic infolist without a second stable consumer |
| Create/edit/delete/confirmation | Covered at recipe/demo level | Consumer owns authorization and mutation rules |
| Notifications and feedback | Covered as presentation | Delivery, persistence, and channels remain consumer-owned |
| Dashboard stats | Contract documented | Real metrics require a consumer-owned source and period model |
| Chart widgets | Deferred | Add only with accessible table/text alternative and real data contract |
| Custom pages/widgets | Partial | Recipe composition is documented; widget layout contract is not yet a product need |
| Global search | Demo slice | `/recipes/search?q=` searches Admin Resource + Ops Queue; production authorization/indexing remain consumer-owned |
| Export | Contract + demo slice | `docs/gelium-ui-export-recipe.md`; Ops Queue has a bounded CSV demo; production async/export policy still waits for a real consumer dataset |
| Import | Contract + demo slice | `docs/gelium-ui-import-recipe.md`; Admin Resource has a bounded synchronous CSV demo; production async/import policy remains consumer-owned |
| Relationship managers / nested resources | Contract + demo slice | `docs/gelium-ui-relationships-recipe.md`; Admin Resource demonstrates Project → Tasks; production persistence and policy remain consumer-owned |
| Multi-tenancy | Boundary only | Add explicit tenant-scope and URL/data-leakage guidance; never implement tenancy in Gelium |
| Authentication / MFA | Intentionally out of scope | Consumer owns identity, sessions, MFA, recovery, and security policy |
| Authorization | Boundary documented | Consumer checks permissions at render and mutation boundaries, including per-record and bulk actions |
| Plugin system | Intentionally deferred | Recipe/component registry first; no runtime plugin mechanism without a real consumer |
| Testing | Partial | Existing contract tests; add a recipe-level verification template |
| Deployment / production optimization | Partial | Release gate exists; operational deployment remains consumer-owned |
| UI component primitives | Broadly covered | Add components only through the registry and implementation contract |

## Gaps worth planning

### 1. Export contract

Filament treats export as more than a button: it includes column selection, formatting, relationship data, row limits, chunking, queued completion, authorization, and CSV formula-injection protection.

Gelium should eventually define a framework-neutral export recipe with:

- explicit authorized dataset and column allowlist;
- synchronous versus asynchronous threshold;
- stable download or job-status URL;
- row and file-size limits;
- CSV/XLSX format policy;
- formula-injection handling;
- sensitive-field exclusion;
- completion, failure, and retry states.

### 2. Import contract

Import is a workflow, not a file input. The contract needs upload validation, column mapping, per-row validation, sensitive-field handling, relationship resolution, limits, progress, retry, and partial-failure semantics.

The likely Gelium shape is a server-rendered upload and mapping flow backed by a consumer-owned job, with an Ops Queue or job-status recipe when processing is asynchronous.

### 3. Global search

A global search recipe should define:

- authorized resource and field scope;
- query normalization and minimum length;
- result grouping and stable URLs;
- no-result and error states;
- tenant isolation;
- keyboard and no-JS navigation;
- rate and result limits.

Do not add a client-only command palette as a substitute for a server contract.

### 4. Relationships and nested resources

Filament exposes relationship managers and nested resources. Gelium now demonstrates one explicit Project → Tasks consumer slice without turning it into a generic relationship runtime. The framework-neutral contract defines:

- which relation is declared by the consumer;
- how relation data is loaded and authorized;
- nested URL structure and breadcrumbs;
- empty, unavailable, and partial-failure states;
- mutation boundaries for related records;
- protection against leaking unauthorized related records.

### 5. Multi-tenancy boundary

Gelium should explicitly document that every consumer-owned query, route, action, export, import, search, metric, and relation must be scoped to the active tenant before rendering or mutation. Gelium must not provide tenant resolution, identity, or policy storage.

### 6. Recipe-level testing

Component contract tests are not enough for a screen recipe. Add a reusable verification template covering:

- ordinary GET and invalid GET parameters;
- empty and error responses;
- authorized and denied actions;
- POST + 303 mutation flow;
- 422 field preservation and validation summary;
- no-JS rendering;
- HX-Request enhancement where supported;
- accessible names, landmarks, roles, and focus relationships;
- tenant/resource scope where the consumer supplies it.

## Deliberately out of scope

These Filament capabilities should not be ported into Gelium's central contract:

- Laravel/Eloquent/Livewire integration;
- authentication, MFA, sessions, password recovery, or user models;
- authorization policy implementation;
- audit persistence;
- notification channels and persistence infrastructure;
- queue workers, cache, filesystem, mail, or database adapters;
- runtime plugin discovery and package lifecycle.

Gelium may document integration boundaries and render the resulting state, but the consuming application owns those systems.

## Recommended order

1. Add explicit multi-tenant boundary guidance and the recipe-level testing template.
2. Choose **one** real consumer problem: export, import, or relationships (global search now has an Ops Queue implementation).
3. Write or refine its framework-neutral contract before implementing UI.
4. Implement one narrow recipe with no-JS and consumer-owned authorization.
5. Re-run the capability audit after that consumer has evidence.

## Sources

- Filament documentation index: <https://filamentphp.com/docs>
- Filament full documentation feed: <https://filamentphp.com/docs/llms-full.txt>
- Gelium component registry: `docs/gelium-ui-component-registry.md`
- Gelium pattern registry: `docs/gelium-ui-pattern-registry.md`
- Gelium screen recipes: `docs/gelium-ui-screen-recipes.md`
- Gelium extensibility contract: `docs/gelium-ui-extensibility.md`
