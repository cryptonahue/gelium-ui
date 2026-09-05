# Import recipe contract

Import is a controlled data workflow, not only a file-upload control. The consuming application owns parsing, authorization, validation, persistence, jobs, and rollback policy; Gelium renders the workflow and its states.

## Quick path

1. Accept a bounded file through a server-rendered form.
2. Validate file type, size, encoding, and structure before mapping.
3. Let the consumer confirm an explicit source-to-field mapping.
4. Validate rows before mutation and show actionable errors.
5. Process synchronously only when bounded; otherwise expose an authorized job-status page.
6. Report complete, partial, failed, cancelled, and retryable outcomes.

## Contract

An import request should declare:

```text
resource
file
mapping
mode
requested-by
tenant-scope
```

The consumer must define:

- accepted formats, encoding, delimiter, and header rules;
- maximum file size, rows, columns, and processing time;
- required and optional destination fields;
- type casting, normalization, and null semantics;
- create-only versus update-or-create behavior;
- duplicate and conflict policy;
- relationship resolution policy;
- sensitive-field policy;
- transaction and partial-failure behavior;
- retention and cleanup for uploaded files.

Unknown fields, mappings, formats, or modes must be rejected or normalized through closed vocabularies. Never map arbitrary input directly to every model field.

## Workflow states

```text
uploaded -> mapping_required -> validated -> queued -> processing -> completed
uploaded -> rejected
mapping_required -> validation_failed
queued -> cancelled
processing -> partial_failed
processing -> failed -> retryable|terminal
```

The consumer owns state persistence and job execution. A status page must be a real URL and must re-check actor, tenant, job ownership, expiration, and authorization.

## Validation

Validate in layers:

1. **File:** size, format, encoding, delimiter, header, and malware/storage policy.
2. **Mapping:** required destination fields, duplicate mappings, unsupported fields, and type declarations.
3. **Row:** required values, types, lengths, formats, closed vocabularies, and domain rules.
4. **Cross-row:** duplicates, ordering, references, and conflicts.
5. **Mutation:** current authorization, tenant scope, and record state immediately before persistence.

A validation response should identify row, column, field, reason, and recovery action without exposing secrets or unrelated records.

## Partial failures

The consumer must choose and document one policy:

- **Atomic:** any invalid row prevents all writes.
- **Best effort:** valid rows persist and invalid rows are reported.
- **Staged:** rows are validated first and a separate confirmation commits the batch.

Never label a best-effort import as fully successful. Report counts for accepted, rejected, skipped, and failed rows, plus a downloadable or viewable error report when appropriate.

## Security and scope

- Re-check authorization and tenant scope at upload, validation, processing, and status/download boundaries.
- Treat uploaded content as untrusted data; never execute formulas, macros, or embedded code.
- Exclude sensitive destination fields unless explicitly allowlisted.
- Limit rows, bytes, columns, processing time, concurrency, and retry count.
- Avoid storing raw file contents longer than the consumer's retention policy allows.
- Do not trust filenames, MIME types, hidden fields, or client-supplied tenant IDs.
- Prevent unauthorized users from observing row values, counts, job identifiers, or failure details.

## Feedback and states

| Situation | Feedback |
|---|---|
| File rejected | 422 + inline alert or validation summary |
| Mapping required | Server-rendered mapping form with required-field guidance |
| Row validation failed | 422 or validation result with row/field errors |
| Import queued | Persistent status page or banner with job link |
| Import processing | Loading/progress state with readable counts when available |
| Complete | Persistent success with accepted count and recovery/download links |
| Partial failure | Persistent warning with accepted/rejected/failed counts and error report |
| Cancelled | Persistent status with explicit reason and restart action when safe |
| Terminal failure | Error state with retry guidance only when idempotency is defined |
| No rows | Empty state; do not claim success for an empty file |

Toast can supplement a server-driven update, but it must never be the only channel for validation, partial failure, or completion.

## No-JS baseline

Upload, mapping, validation results, confirmation, status, retry, cancellation, and error-report links must work through native forms and links. HTMX may enhance progress polling or region updates, but it must not own job state, authorization, or recovery.

## Accessibility

- Label file, format, mapping, mode, and confirmation controls.
- State accepted formats and limits before upload.
- Associate row/field errors with controls where possible and provide a summary for multiple errors.
- Announce processing, completion, and partial failure persistently.
- Provide a readable alternative to downloadable error reports.
- Keep focus and headings coherent across mapping, validation, and status pages.
- Do not use color or transient notifications as the only status signal.

## Testing checklist

- [ ] File, format, encoding, size, row, and column limits are enforced.
- [ ] Mapping is explicit and rejects unsupported destination fields.
- [ ] Row and cross-row validation report actionable locations.
- [ ] Authorization and tenant scope are re-checked at every boundary.
- [ ] Atomic, best-effort, or staged behavior is explicit.
- [ ] Duplicate, retry, cancellation, and idempotency behavior is defined.
- [ ] Sensitive data is excluded from output and error reports.
- [ ] Empty, validation, queued, processing, complete, partial, and failed states are covered.
- [ ] No-JS upload, mapping, status, and recovery paths work.
- [ ] Accessibility semantics and readable error alternatives are verified.

## Gelium boundary

Gelium does not parse files, persist imports, run jobs, choose domain mappings, authorize records, resolve tenants, or decide transaction policy. It provides components and recipes for the consumer-owned contract.

## Related

- [Filament capability audit](plans/2026-09-05-filament-gap-audit.md)
- [Application integration](gelium-ui-application-integration.md)
- [Recipe testing](gelium-ui-recipe-testing.md)
- [Ops Queue](gelium-ui-screen-recipes.md)
- [Export recipe](gelium-ui-export-recipe.md)
