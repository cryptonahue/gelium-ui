# Integrated feedback recipe

Gelium UI keeps feedback server-rendered and contextual. The consuming application owns the operation and business outcome; Gelium provides the presentation contract.

## Decision matrix

| Situation | Canonical feedback | Server contract | Persistence |
|---|---|---|---|
| Transient result of an action | Toast | Consumer may emit the agreed server-driven event; provide a non-JS fallback | Temporary |
| Successful POST that must survive navigation | Banner or inline alert | POST + 303, then render the success state on the destination | Persistent |
| Error attached to a section or field | Inline alert | Re-render the affected region; field errors preserve values and associations | Persistent |
| Form validation with multiple errors | 422 + validation summary | 422 response includes field errors and links to affected controls | Persistent until corrected |
| Resource cannot be found or operation fails outside a form | Error state | Render the appropriate 4xx/5xx state with recovery navigation | Persistent |
| Authorized collection has no records | Empty state | Render the collection response with guidance and optional CTA | Persistent |

## Selection rules

- Use **Toast** only for a non-critical, transient result after an action.
- Use **Banner** for a page- or site-level message that requires awareness or action.
- Use **inline alert** for an error or warning tied to a form or section.
- Use **validation summary** when two or more form fields need attention; keep field-level messages as well.
- Use **error state** for a failed request or unavailable resource when there is no recoverable field context.
- Use **empty state** when the authorized query succeeded but returned no records. Do not turn absence of data into an error.
- Never use a Toast as the only channel for validation, authorization denial, a critical failure, or a success message that must survive navigation.

## Server-rendered baseline

Every recipe must remain understandable with JavaScript disabled:

1. The server returns the resulting page or region with the feedback state.
2. A POST that changes state uses POST + 303 where navigation follows the mutation.
3. A 422 response preserves submitted values and exposes validation errors.
4. A Toast enhancement may be added for supported server-driven interactions, but the server-rendered fallback remains authoritative.
5. Recovery actions are real links or forms with visible labels; never rely on auto-dismiss or client-only navigation.

## Accessibility baseline

- Put persistent feedback in the document order near the affected content.
- Use `role="alert"` for errors that need immediate announcement and `role="status"` for non-error updates.
- Give validation summaries a heading and links to each affected control.
- Keep `aria-invalid` and `aria-describedby` relationships intact on invalid fields.
- Do not communicate severity or success through color alone.
- Toasts must have a live region, a manual dismiss action, and a pauseable timeout when enhancement is available.

## Consumer boundary

The consumer application decides authorization, operation outcome, message content, persistence, and whether a server-driven enhancement is appropriate. Gelium does not authenticate users, authorize actions, persist audit history, or infer business success from an HTTP status alone.

## Recipe checklist

- [ ] The feedback type matches the situation in the decision matrix.
- [ ] The server response is usable without JavaScript.
- [ ] Mutation flows use POST + 303 when navigation follows success.
- [ ] Validation uses 422 and preserves values and field associations.
- [ ] Error and empty states include a clear recovery or next action where appropriate.
- [ ] Accessibility roles and relationships are present.
- [ ] Authorization and business meaning remain consumer-owned.

## Related

- [Pattern vocabulary](gelium-ui-vocabulary.md)
- [Screen recipes](gelium-ui-screen-recipes.md)
- [Application integration](gelium-ui-application-integration.md)
