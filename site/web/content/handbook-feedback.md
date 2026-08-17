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

## Decision matrix

| Situation | Use | Do not use | Gelium |
|---|---|---|---|
| **Field validation** after submit | Error **summary** + **inline** field errors | Toast-only; notification banner for validation | `validation-summary` + field `error` / `aria-invalid`; HTTP **422** + `X-Gelium-Validation` |
| **Single field** guidance while typing (rare) | Inline message after interaction | Blocking modal; premature errors on empty focus | Prefer submit-time validation ([Forms](/docs/forms)); NNG: avoid grading before the user answered |
| **Task succeeded** (saved, sent, deleted) | Short confirmation: toast and/or inline success on the result page | Error styling; modal that only says “OK” | `gelium:toast` success; or success `banner` / `inline-alert` when the page itself is the result |
| **Task failed** (server/network), not a field | Persistent inline or page alert + recovery action | Toast that disappears before it can be read (for critical failure) | `inline-alert` / `banner` error; toast only if non-blocking and repeated safe |
| **System / service-wide** notice | Site-level alert pattern | Mixing with field validation on the same focus | `banner` at layout level; USWDS **site alert** = every page / system status, not form response |
| **Page-level status** (you completed a step) | Page alert / banner in main content | Site-wide chrome for a one-page message | `banner` or `inline-alert` in the main column |
| **Empty collection** | Empty state with reason + next step | Silent blank table | `empty-state` |
| **Section failed to load** | Error state + retry | Infinite skeleton | `error-state` |
| **Content loading** | Skeleton / progressive reveal | Full-page blocker for every fetch | `skeleton` |
| **Consequential choice** (destroy, pay, irrevocable) | Confirm screen or dialog with clear verbs | Toast “Deleted?” with no undo and no confirm | `dialog` confirm route + POST |
| **Low-priority FYI** after an action | Snackbar/toast-class feedback (M3 snackbar: brief, non-blocking) | Using toast for validation lists | `toast` via `gelium:toast` |

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

## Server mapping

| Feedback | HTTP / wire |
|---|---|
| Validation | `422` + `X-Gelium-Validation: true` + summary + fields |
| OK mutate | `303` redirect to GET result; optional `HX-Trigger: gelium:toast` |
| Transport failure | Client may toast error; page should still offer retry if critical |

## Anti-patterns

- Validation errors only in a toast (far from fields; easy to miss) — called out across GOV.UK and industry guidance summarizing NNG.
- Premature inline errors before the user finished the field ([NNG](https://www.nngroup.com/articles/error-message-guidelines/)).
- Error summary without field-level messages (or the reverse) on submit validation.
- Global `banner` for a single field typo.
- Modal dialog for “Saved successfully” with no decision to make.

## Checklist (agents)

1. Is this **validation**, **task result**, **system status**, or **empty/load**?
2. Pick the row in the matrix — do not invent a fifth channel.
3. Write copy per Content style.
4. Wire the matching server contract.
5. Check narrow viewports still show the message (no clip via overflow tricks).

## See also

- [Screens](/docs/screens) — which page type you are on  
- [Forms](/docs/forms) — labels and input behavior  
- Agent pack: [`/llms-ux.txt`](/llms-ux.txt)
