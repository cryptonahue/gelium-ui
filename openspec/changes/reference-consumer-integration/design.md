# Design: Reference Consumer Integration

## Boundary

The reference consumer lives outside `lib/` and owns its SQLite connection, schema, repositories, tenant scope, policy checks, and audit writes. Gelium remains a presentation and HTTP contract layer.

## Data model

- `tenants(id, name)`
- `projects(id, tenant_id, name, status, date, owner)`
- `tasks(id, tenant_id, project_id, title, status, due_date, assignee)`
- `activity(id, tenant_id, project_id, type, summary, actor, created_at)`
- `audit_events(id, tenant_id, actor, action, resource, resource_id, created_at)`

Every project, task, activity, and audit lookup requires a tenant scope. Project-child lookups validate both project identity and tenant identity.

## Integration

The reference consumer exposes an application-owned handler/repository composition. Existing recipe handlers remain demo fixtures; the new integration is a separate consumer path to avoid coupling package contracts to SQLite.

Mutations perform policy authorization, repository mutation, and audit write in one consumer-owned transaction. The HTTP response preserves Gelium's POST + 303 and 422 contracts.

## Test strategy

Use a temporary SQLite database per test. Cover schema setup, repository round trips, tenant isolation, denied actions, successful mutations with audit records, and no-JS HTTP navigation. Do not require external services.
