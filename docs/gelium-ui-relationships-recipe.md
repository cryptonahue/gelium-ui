# Relationships and nested resources recipe

Related records are an application-owned data and authorization concern. Gelium can render a relationship manager or nested resource recipe, but it must never infer joins, fetch undeclared fields, or assume that visibility of a parent grants access to its children.

## Gelium demo consumer

The Admin Resource demo implements deliberately small `Project → Tasks` and read-only `Project → Activity` relationships:

```text
GET  /recipes/admin-resource/{project-id}/tasks?status=
POST /recipes/admin-resource/{project-id}/tasks
GET  /recipes/admin-resource/{project-id}/activity?type=
```

It demonstrates project-scoped reads, closed task/activity filters, native task creation, `422` validation, and `POST + 303` navigation. Activity is read-only. Its in-memory data and allow-all demo authorization are not production persistence or policy implementations.

## Quick path

1. Declare the parent resource and relationship explicitly.
2. Resolve the active actor and tenant before loading either side.
3. Query only authorized related records with bounded fields and pagination.
4. Keep relationship state navigable with stable GET URLs.
5. Re-check authorization and current parent/child state for every mutation.
6. Render empty, unavailable, partial, and error states distinctly.

## Relationship declaration

A consumer-owned relationship definition should specify:

```text
parent-resource
parent-id
relationship-key
child-resource
read-scope
allowed-actions
fields
sort
pagination
tenant-scope
```

`relationship-key` and fields use closed vocabularies. The consumer decides whether the relationship is one-to-one, one-to-many, many-to-many, polymorphic, computed, or unavailable for the current actor.

## URL contract

Nested navigation should use stable URLs, for example:

```text
GET /admin/{parent}/{parent-id}/related/{relationship}
GET /admin/{parent}/{parent-id}/related/{relationship}?q=&sort=&page=
GET /admin/{parent}/{parent-id}/related/{relationship}/new
POST /admin/{parent}/{parent-id}/related/{relationship}
GET /admin/{parent}/{parent-id}/related/{relationship}/{child-id}/edit
POST /admin/{parent}/{parent-id}/related/{relationship}/{child-id}/edit
```

The exact path is consumer-defined. Parent identity, relationship key, filters, pagination, and selection must be validated server-side. Breadcrumbs and back links must preserve the parent context without trusting client-supplied labels.

## Read scope and disclosure

The consumer must authorize independently:

- access to the parent;
- visibility of the relationship itself;
- visibility of each child record;
- fields shown for each child;
- actions available on each child.

If the parent or relationship is unavailable, use the consumer's non-disclosing response, commonly `404`. Do not reveal counts, names, IDs, or relationship existence through disabled controls, empty labels, metadata, or error messages.

## Mutation contract

Relationship mutations use native forms and follow the normal Gelium contracts:

- create/update/delete uses POST;
- successful navigation uses 303;
- validation uses 422 and preserves submitted values;
- authorization and tenant scope are re-checked at the mutation boundary;
- parent and child current state are checked for stale or conflicting operations;
- partial relationship operations report affected, rejected, and failed records explicitly.

A hidden parent ID or relationship key is not authorization. The consumer must derive or validate both from the current route and request scope.

## States

| Situation | Feedback |
|---|---|
| Parent unavailable | Non-disclosing 404/error state with recovery link |
| Relationship unavailable | Non-disclosing 404/error state or consumer-defined access response |
| Authorized relationship with no children | Empty state with an appropriate create or recovery action |
| Related data loading | Server-rendered loading/progress state only for an actual waiting region |
| Related data failure | Error state or contextual alert; never an empty success state |
| Invalid relationship filter | 422 or safe normalization using closed values |
| Child validation failure | 422 with field errors and summary when multiple fields fail |
| Concurrent/stale child state | Persistent warning or error with refresh/retry guidance |
| Successful mutation | Persistent success after POST + 303 |
| Partial bulk mutation | Persistent warning with explicit counts and recovery |

## Tenant isolation

When tenancy applies, the active tenant scope must be present in parent lookup, relationship lookup, filters, pagination, mutations, background jobs, and download URLs. A child belonging to another tenant must not be reachable through a parent route, guessed child ID, relationship count, or search query.

## No-JS baseline

Parent navigation, relationship tabs/links, search, filters, pagination, create/edit forms, confirmation, validation, and recovery must work with JavaScript disabled. HTMX may swap a relationship region, but it must not change authorization, tenant scope, URL meaning, or mutation semantics.

## Accessibility

- Use a heading and landmark that identify the parent and relationship context.
- Keep breadcrumbs and back links keyboard reachable and descriptive.
- Associate child forms and validation messages with their controls.
- Preserve `aria-current` for the active relationship navigation.
- Make empty, unavailable, stale, and partial states readable without color alone.
- Avoid nested interactive controls inside row links or action regions.

## Testing checklist

- [ ] Parent, relationship, child, field, and action visibility are independently authorized.
- [ ] Foreign parent and child identifiers are non-disclosing.
- [ ] Relationship keys, fields, filters, sort, and pagination use closed validation.
- [ ] Tenant scope is re-checked for reads, mutations, jobs, and downloads.
- [ ] Empty, unavailable, failure, stale, validation, and partial states are distinct.
- [ ] POST + 303 and 422 value preservation work for relationship mutations.
- [ ] No-JS navigation and forms work.
- [ ] HTMX enhancement preserves the same URL and authorization contracts.
- [ ] Breadcrumb, heading, focus, labels, and `aria-current` semantics are valid.

## Gelium boundary

Gelium does not define ORM relationships, query joins, eager loading, tenant resolution, authorization policies, cascade behavior, transactions, or referential integrity. It renders the consumer's explicit relationship contract.

## Related

- [Filament capability audit](plans/2026-09-05-filament-gap-audit.md)
- [Application integration](gelium-ui-application-integration.md)
- [Recipe testing](gelium-ui-recipe-testing.md)
- [Extensibility](gelium-ui-extensibility.md)
- [Screen recipes](gelium-ui-screen-recipes.md)
