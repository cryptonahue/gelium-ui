# Export recipe contract

Export is a data workflow, not merely a download button. The consuming application owns the dataset, authorization, formatting, and job execution; Gelium renders the server-driven flow and its states.

## Quick path

1. Resolve the authorized dataset and exportable column allowlist.
2. Validate format, filters, selected records, and row/file limits.
3. Choose synchronous download or an asynchronous export job.
4. Re-check authorization and tenant scope when generating and downloading.
5. Communicate completion, failure, limits, and retry through the feedback recipe.

## Contract

An export request should declare:

```text
resource
filters
sort
selection
columns
format
requested-by
tenant-scope
```

The consumer must define:

- which fields are exportable;
- how labels, units, dates, timezones, and nulls are formatted;
- whether relationships may be included;
- maximum rows and output size;
- supported formats and encoding;
- synchronous versus asynchronous threshold;
- retention and expiration of generated files;
- authorization and tenant-scope rules.

Unknown columns, formats, filters, or record IDs must be rejected or normalized using a closed vocabulary. Never export every field by default.

## Synchronous and asynchronous paths

### Synchronous

Use only for bounded output that can complete within the application's request budget. Return the file with a content disposition and a safe, consumer-owned filename. The endpoint must enforce authorization and current scope.

### Asynchronous

Use a server-rendered request form followed by a job-status page or an authorized download link. The consumer owns the queue, job state, worker, storage, retry policy, and cleanup. Suggested states are:

```text
queued -> processing -> completed
queued -> cancelled
processing -> failed -> retryable|terminal
```

A download URL must re-check authorization, tenant scope, expiration, and job ownership. Do not expose storage paths or predictable identifiers.

## Security

- Exclude sensitive or internal fields unless explicitly allowlisted.
- Re-check per-record authorization for the final dataset.
- Scope queries, jobs, and downloads to the active tenant where applicable.
- Neutralize spreadsheet formula injection according to the consumer's security policy.
- Enforce row, byte, time, and concurrency limits.
- Escape or encode values according to the selected format.
- Avoid putting secrets, raw queries, or authorization decisions in URLs.
- Expire generated files and provide an explicit failure state after expiration.

## Feedback and states

| Situation | Feedback |
|---|---|
| Export request accepted | Persistent success or status page with job reference |
| Export processing | Loading/progress state with readable status |
| Export completed | Authorized download link with expiration context |
| No authorized records | Empty state; do not fabricate a zero-row file unless requested by policy |
| Invalid filters/columns | 422 + validation summary or inline alert |
| Limit exceeded | Inline alert with an actionable smaller-scope suggestion |
| Job failure | Error state or persistent banner with retry when safe |
| Expired or unauthorized download | Non-disclosing error state |

A Toast may supplement a completed server-driven action, but it must not be the only indication that a file is ready or a job failed.

## No-JS baseline

The primary export request must work as a native form submission. The status page, refresh/retry action, and download link must be real server-rendered navigation or forms. HTMX may enhance progress refreshes, but it must not own authorization, job state, or the only recovery path.

## Accessibility

- Label format, scope, filters, and column selection controls.
- Explain limits before submission where possible.
- Announce processing and completion with persistent readable status.
- Give download links descriptive names and expiration context.
- Keep validation errors associated with their controls and provide a summary for multiple errors.
- Do not rely on color or a transient notification to communicate completion or failure.

## Testing checklist

- [ ] Exported columns are an explicit allowlist.
- [ ] Filters, sort, selection, format, and limits use closed validation.
- [ ] Unauthorized records and foreign-tenant records are excluded.
- [ ] Authorization is re-checked during generation and download.
- [ ] Synchronous and asynchronous thresholds are documented.
- [ ] Job status, retry, expiration, and cleanup behavior are defined.
- [ ] Sensitive fields and formula injection are covered.
- [ ] Empty, validation, limit, processing, success, and failure states are covered.
- [ ] No-JS request, status, retry, and download paths work.
- [ ] Accessibility semantics and download naming are verified.

## Gelium boundary

Gelium does not query databases, run workers, store files, authorize records, resolve tenants, or decide which fields are sensitive. It provides components and recipes for the consumer-owned contract.

## Related

- [Filament capability audit](plans/2026-09-05-filament-gap-audit.md)
- [Application integration](gelium-ui-application-integration.md)
- [Recipe testing](gelium-ui-recipe-testing.md)
- [Feedback recipe](gelium-ui-feedback-recipe.md)
- [Extensibility](gelium-ui-extensibility.md)
