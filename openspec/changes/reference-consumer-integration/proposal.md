# Proposal: Reference Consumer Integration

## Problem

Gelium's recipes demonstrate resource workflows, but the demos currently use in-memory stores and permissive policies. They do not show how a consumer owns persistence, tenant scope, authorization, or audit boundaries.

## Outcome

Add a small reference consumer application using SQLite that exercises Gelium recipes without moving consumer concerns into `lib/`.

## Scope

- SQLite-backed projects, tasks, and activity records in the reference consumer.
- Explicit consumer-owned repository, tenant scope, authorization policy, and audit sink boundaries.
- Server-rendered integration for the existing Admin Resource workflows.
- Deterministic tests for persistence, tenant isolation, authorization, and audit events.
- Documentation explaining which contracts belong to Gelium and which belong to the consumer.

## Non-goals

- No ORM abstraction in Gelium.
- No authentication, authorization, tenancy, audit, or job runtime in `lib/`.
- No production deployment or migration framework.
- No replacement of the existing in-memory recipe fixtures unless the reference consumer owns the integration path.
- No automatic schema discovery or generic CRUD generation.

## Constraints

- HTML-first and no-JavaScript by default.
- GET for navigable state; POST + 303 for successful mutations; 422 for validation.
- SQLite is limited to the reference consumer boundary.
- Delivery actions remain separate from implementation.
