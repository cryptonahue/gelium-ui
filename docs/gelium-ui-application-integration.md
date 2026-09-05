# Gelium UI application integration

Gelium UI renders server-driven screens. The consuming application owns identity, authorization, domain mutations, and audit persistence. This separation keeps Gelium usable with any language, framework, database, or authentication system.

## Quick path

1. Resolve the current actor from the application's authentication boundary.
2. Ask the application's authorization layer whether the actor may perform the action.
3. Render only actions that are allowed, but enforce the same decision on the request that mutates data.
4. Execute the domain operation in the application layer.
5. Record the outcome through the application's audit service.
6. Return the Gelium response contract: `303` on success, `422` for validation, and `401` or `403` for access failures.

## Responsibility boundary

| Concern | Owner | Gelium UI role |
|---|---|---|
| Authentication and sessions | Consumer application | Render the resulting state |
| Roles, permissions, tenancy, and ownership | Consumer application | Receive an allow/deny decision |
| Record visibility | Consumer application | Render the authorized collection |
| Domain mutation | Consumer application | Provide forms and action affordances |
| Audit storage and retention | Consumer application | Show success or failure feedback |
| HTML semantics, states, and HTTP presentation | Gelium UI | Provide reusable primitives and recipes |

Gelium must not require a particular authentication package, permission vocabulary, ORM, middleware stack, or audit database.

## Authorization contract

The consumer supplies an authorization decision for the current request, action, resource, and optional record:

```text
can(request, action, resource, record) -> allowed
```

The exact API is consumer-defined. The important properties are:

- collection visibility is filtered before rendering;
- row actions are rendered only when the row is actionable;
- direct URLs are protected even when the action is hidden;
- the mutation endpoint checks authorization again;
- authorization is evaluated against current data, not a confirmation page or hidden field.

### HTTP outcomes

| Situation | Response | Meaning |
|---|---:|---|
| No authenticated actor | `401 Unauthorized` | The authentication boundary must establish identity first |
| Authenticated actor lacks permission | `403 Forbidden` | The request is understood but not allowed |
| Resource does not exist or must not be disclosed | `404 Not Found` | The application controls information disclosure |
| Input is invalid | `422 Unprocessable Entity` | Re-render the form with validation feedback |

## Audit contract

The consumer's audit service should record one event for each consequential action:

```text
actor
action
resource IDs
result
partial failures
timestamp
```

Recommended result values are `succeeded`, `partially_succeeded`, and `failed`. The application may add request ID, tenant, reason, source, or correlation data according to its domain and retention policy.

Gelium does not persist audit events and does not infer the actor. A success banner, inline alert, or error state communicates the result to the user; it is not an audit record.

## Bulk actions

Bulk authorization is not equivalent to checking one global permission:

1. Resolve the submitted selection against current records.
2. Re-check authorization for every record immediately before mutation.
3. Apply the domain operation according to the application's transaction policy.
4. Record each affected ID and each rejected or failed ID.
5. Communicate a partial result explicitly; never report complete success when some records were skipped.

The application decides whether a partial operation is acceptable, atomic, retryable, or forbidden by policy.

## Framework integration

A framework integration should be a thin adapter around the consumer's existing authentication, authorization, domain, and audit services. It may map framework-native policies to the contract above, but Gelium's public contract must remain framework-neutral.

## Checklist

- [ ] The consumer owns authentication and authorization.
- [ ] Collection queries exclude records the actor cannot view.
- [ ] Hidden actions are also protected at their endpoints.
- [ ] Mutations re-check authorization using current data.
- [ ] Bulk operations report partial failures explicitly.
- [ ] Audit records contain actor, action, IDs, result, failures, and timestamp.
- [ ] Gelium is not responsible for audit persistence or policy storage.

## Related

- [Server contracts](../lib/skills/05-server-contracts.md)
- [Screen recipes](gelium-ui-screen-recipes.md)
- [Pattern registry](gelium-ui-pattern-registry.md)
