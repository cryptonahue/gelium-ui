# Feedback

Feedback tells the user **what happened** and **what to do next**. Gelium ships several message primitives; using the wrong one is a common agent mistake (for example, field errors only in a toast).

Rules below are adapted from **GOV.UK Design System**, **USWDS**, **Material 3**, and **Nielsen Norman Group**, then mapped to Gelium components and server contracts.

## Sources (read these)

| Topic | Source |
|---|---|
| Validation errors | [GOV.UK: error summary](https://design-system.service.gov.uk/components/error-summary/), [error message](https://design-system.service.gov.uk/components/error-message/) |
| Service-level notices | [GOV.UK: notification banner](https://design-system.service.gov.uk/components/notification-banner/) |
| Page vs site alerts | [USWDS: alert](https://designsystem.digital.gov/components/alert/), [site alert](https://designsystem.digital.gov/components/site-alert/) |
| Brief low-priority notices | [M3: snackbar](https://m3.material.io/components/snackbar/guidelines) |
| Error severity & timing | [NNG: error-message guidelines](https://www.nngroup.com/articles/error-message-guidelines/) |
| Copy voice | [Content style](/docs/content-style) |
| 422 / toast wire | [Server contracts](/docs/server-contracts) |

Stable rule IDs (`FEED-*`) match [`/llms-ux.txt`](/llms-ux.txt) so agents and humans share one vocabulary.

## Decision matrix

| ID | Situation | Use | Do not use | Gelium |
|---|---|---|---|---|
| **FEED-VAL** | **Field validation** after submit | Error **summary** + **inline** field errors | Toast-only; notification banner for validation | `validation-summary` + field `error` / `aria-invalid`; HTTP **422** + `X-Gelium-Validation` |
| **FEED-INLINE** | **Single field** guidance while typing (rare) | Inline message **after** interaction | Blocking modal; premature errors on empty focus | Prefer submit-time validation ([Forms](/docs/forms)); NNG: avoid grading before the user answered |
| **FEED-OK-TOAST** | **Task succeeded**, user stays in flow | Brief **toast** (optional short action) | Error styling; modal that only says “OK” | `gelium:toast` success; M3 snackbar = low priority, non-blocking |
| **FEED-OK-PAGE** | **Task succeeded** and the next page **is** the result | Success **banner** / **inline-alert** on the result page; toast optional echo | Toast alone with no destination context | After **303** to GET: show durable success in main column |
| **FEED-FAIL** | **Task failed** (server/network), not a field | Persistent **inline-alert** or page **banner** + recovery | Toast that vanishes before it can be read (critical) | `inline-alert` / `banner` error + retry control |
| **FEED-SYS** | **System / service-wide** notice | Site-level alert on layout | Mixing with field validation | Layout `banner`; USWDS **site alert** = every page / system status, not form response |
| **FEED-PAGE** | **Page-level status** (step completed, warning on this view) | Page alert / banner in **main** content | Global chrome for a one-page message | `banner` or `inline-alert` in the main column |
| **FEED-EMPTY** | **Empty collection** | Empty state with reason + next step | Silent blank table | `empty-state` |
| **FEED-LOAD-FAIL** | **Section failed to load** | Error state + retry | Infinite skeleton | `error-state` |
| **FEED-LOAD-LIST** | **List / table loading** | Skeleton rows or list placeholders | Full-page spinner for every partial fetch | `skeleton` scoped to the list region |
| **FEED-LOAD-PAGE** | **First paint of a whole view** | Page-level skeleton **once**; then content | Blocking the app forever | Prefer progressive: shell + skeleton main |
| **FEED-ROW** | **One row / one item** failed (bulk or table action) | Error **in or next to the row** + optional page summary | Toast-only listing every row error | Inline cell/row message; summary if many rows share one cause |
| **FEED-PARTIAL** | **Batch partial success** (some OK, some not) | Summary of outcomes + per-item status | Single green toast “Done” | `banner` or `inline-alert` + list/table status `badge`s |
| **FEED-CONFIRM** | **Consequential choice** (destroy, pay, irrevocable) | Confirm screen or dialog with clear verbs | Toast “Deleted?” with no undo and no confirm | `dialog` confirm route + POST |
| **FEED-FYI** | **Low-priority FYI** after an action | Snackbar/toast-class feedback | Toast for validation lists | `toast` via `gelium:toast` |

### Toast rules (M3 snackbar, adapted)

[M3 snackbars](https://m3.material.io/components/snackbar/guidelines) are **brief** and **low priority**. Gelium `toast` follows that ladder:

| Rule | Do | Don’t |
|---|---|---|
| Purpose | Confirm a completed action or safe FYI | Carry the only copy of field validation |
| Duration | Short; OK if missed for non-critical OK | Rely on toast for money/safety failures |
| Action | At most **one** short action (e.g. “Undo” when truly reversible) | Multi-step workflows inside a toast |
| With redirect | Toast **echo** after 303 is fine; durable success still on the result page when the page *is* the confirmation | Toast-only success when the user lands on a blank list with no context |
| Stacking | One toast at a time in `#gelium-toast-region` | Toast spam per keystroke |

NNG: match **severity** to interruption cost — modals/dialogs for severe blocking issues; transient notices for minimal interaction ([error-message guidelines](https://www.nngroup.com/articles/error-message-guidelines/)).

### Loading rules

| ID | Pattern | Prefer | Avoid |
|---|---|---|---|
| FEED-LOAD-LIST | Collection refetch | Skeleton **inside** the list/table card | Replacing the whole app chrome |
| FEED-LOAD-PAGE | Slow first navigation | One page skeleton, then content | Nested full-page blockers |
| FEED-LOAD-FAIL | Fetch error | `error-state` + retry where the content was | Endless skeleton |

### GOV.UK rules we adopt literally (adapted names)

From [notification banner — when not to use](https://design-system.service.gov.uk/components/notification-banner/):

- **Do not** use a notification-style banner for **validation errors** — use error summary + field errors.
- **Do not** show a notification banner **and** an error summary for the same problem — show the summary.
- Use banners **sparingly** (banner blindness; NNG evidence cited by GOV.UK).

From [error summary — when to use](https://design-system.service.gov.uk/components/error-summary/):

- When there is a validation error, show **both** summary and per-field messages, even if there is only one error.

### USWDS split we adopt

- **Site-wide system status** → site-alert-like treatment (global banner).  
- **Response to user action / page validation** → page alert or form errors — not the global site alert ([site alert: when to consider something else](https://designsystem.digital.gov/components/site-alert/)).

### Material 3 severity ladder (mapped)

M3 distinguishes **dialog** (blocks until addressed) vs **snackbar** (brief, low priority). Gelium:

| M3 idea | Gelium |
|---|---|
| Dialog (high) | Confirm **page** or modal `dialog` |
| Snackbar (low) | `toast` |
| Mid persistence | `inline-alert` / `banner` in page content |

## Copy requirements

All feedback copy follows [Content style](/docs/content-style):

- Errors: **say the fix** (“Enter the project name”), not only the consequence (“Name is required”).
- Toasts: **verb + result** (“Project created”), not “Success”.
- Empty: **what / why / next**.
- Partial batch: **counts + next** (“3 sent, 2 failed — fix the failed rows”).

## Server mapping

| Feedback | HTTP / wire |
|---|---|
| Validation (FEED-VAL) | `422` + `X-Gelium-Validation: true` + summary + fields |
| OK mutate (FEED-OK-*) | `303` redirect to GET result; optional `HX-Trigger: gelium:toast` |
| Transport failure (FEED-FAIL) | Client may toast if non-critical; page still offers retry when critical |
| Partial batch (FEED-PARTIAL) | Prefer **200/303** to a result view that lists per-item outcomes; avoid “all good” toast |

## Anti-patterns

- Validation errors only in a toast (far from fields; easy to miss) — GOV.UK + common NNG-aligned practice.
- Premature inline errors before the user finished the field ([NNG](https://www.nngroup.com/articles/error-message-guidelines/)).
- Error summary without field-level messages (or the reverse) on submit validation.
- Global `banner` for a single field typo.
- Modal dialog for “Saved successfully” with no decision to make.
- **Green toast “Done”** after a partial batch failure.
- **Row action errors** only as a global toast with no row marker.
- **Infinite skeleton** with no path to `error-state`.
- Toast **and** banner **and** summary for the same validation event.

## Checklist (agents)

1. Classify: **validation** | **task OK** | **task fail** | **system** | **empty** | **load** | **row/partial** | **confirm**.
2. Pick the **FEED-*** row — do not invent a new channel.
3. If toast: is it **low priority** and OK if missed? If no → banner/inline/dialog.
4. Write copy per Content style.
5. Wire the matching server contract.
6. Narrow viewports still show the message (no clip via overflow tricks).

## See also

- [Screens](/docs/screens) — which page type you are on  
- [Forms](/docs/forms) — labels and input behavior  
- Agent pack: [`/llms-ux.txt`](/llms-ux.txt)
