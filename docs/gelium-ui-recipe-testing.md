# Recipe-level testing contract

Component contract tests prove that primitives render correctly. Recipe tests prove that a complete server-rendered screen preserves its data, state, authorization, accessibility, and recovery contracts.

## Quick path

For every recipe, test one representative scenario for each row below. Keep the consumer's data and authorization adapter explicit so the test does not accidentally become a test of framework defaults.

## Required scenarios

| Area | Scenario | Expected evidence |
|---|---|---|
| Initial read | Ordinary `GET` | Correct status, heading, landmarks, data, and links |
| Navigable state | Valid and invalid query parameters | Stable URL; closed values; invalid input is rejected or normalized |
| Empty | Authorized query returns no records | Empty state, explanation, and real recovery/CTA when useful |
| Missing | Record or route is unavailable | Appropriate 404/error state without sensitive disclosure |
| Failure | Data source or operation fails | Error state or contextual alert; no fabricated success or zero |
| Authorization | Allowed and denied actor/action | Allowed affordance works; denied request is protected at endpoint |
| Tenant scope | Cross-tenant record/query/action | Foreign data is absent and direct access is denied/non-disclosing |
| Mutation | Valid POST | Domain operation occurs once, then `303` to destination when appropriate |
| Validation | Invalid POST | `422`, submitted values preserved, field errors associated, summary links work |
| Bulk action | Mixed allowed/denied selection | Per-record re-check and explicit partial result; no false total success |
| No JavaScript | Browser/client enhancement absent | Primary read, navigation, form, validation, and recovery remain usable |
| HTMX enhancement | `HX-Request: true`, when supported | Correct fragment contract without changing business authorization |
| Accessibility | Keyboard and assistive semantics | Labels, headings, roles, focus, `aria-invalid`, and descriptions are coherent |
| Responsive | Narrow and wide viewport | No clipped critical content; actions remain discoverable and usable |
| Theme | Light and dark | Contrast and focus remain valid in both themes |

## Mutation assertions

A recipe test must assert all of the following for a state-changing action:

1. The endpoint re-checks actor, tenant, resource, and current record state.
2. The operation is not performed when authorization fails.
3. A successful operation does not duplicate on refresh.
4. The response communicates the result through the correct persistent or transient feedback pattern.
5. Partial failures name the affected or rejected records according to the consumer policy.

## Data-scope assertions

Use fixtures from at least two tenants when tenancy applies. Assert that:

- list, search, filter, sort, pagination, detail, metrics, export, and import queries are tenant-scoped;
- a foreign identifier cannot retrieve or mutate another tenant's record;
- related records do not cross the active tenant boundary;
- background-job and download URLs re-check scope instead of trusting an earlier request;
- empty and error responses do not reveal foreign-record existence.

## Review evidence

Record the command, fixture scope, response status, and relevant HTML assertions for each scenario. If a scenario cannot run because the consumer adapter is not present, mark it as blocked rather than treating the recipe as verified.

## Checklist

- [ ] Ordinary GET and invalid navigable state are covered.
- [ ] Empty, missing, and dependency-error states are distinct.
- [ ] Authorization is tested at render and endpoint boundaries.
- [ ] Tenant isolation is tested with foreign fixtures when applicable.
- [ ] POST + 303 and 422 value preservation are covered.
- [ ] Bulk partial failures are explicit.
- [ ] No-JS behavior is covered.
- [ ] HTMX behavior is covered when the recipe offers it.
- [ ] Accessibility, responsive, theme, and focus behavior are checked.
- [ ] Evidence records commands and concrete assertions.

## Related

- [Application integration](gelium-ui-application-integration.md)
- [Feedback recipe](gelium-ui-feedback-recipe.md)
- [Screen recipes](gelium-ui-screen-recipes.md)
- [Server contracts](../lib/skills/05-server-contracts.md)
