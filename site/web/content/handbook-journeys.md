# Journeys

A **journey** is more than one screen: where the user starts, the order of steps, where they land after submit, and how they resume. Gelium stays server-first (GET read, POST+303 mutate). Criteria below adapt **GOV.UK** service patterns and map them to routes and [Feedback](/docs/feedback) IDs.

## Sources

| Topic | Source |
|---|---|
| Starting a service | [GOV.UK: start using a service](https://design-system.service.gov.uk/patterns/start-using-a-service/) |
| Multi-task progress | [GOV.UK: task list](https://design-system.service.gov.uk/components/task-list/) |
| Ordered end-to-end guidance | [GOV.UK: step by step navigation](https://design-system.service.gov.uk/patterns/step-by-step-navigation/) |
| Confirmation after success | GOV.UK confirmation-page practice (reference id, next steps) — see also [Screens](/docs/screens) result type |
| Wire after mutate | [Server contracts](/docs/server-contracts), [Feedback](/docs/feedback) |

## Journey shapes

| ID | Shape | When | When not | Gelium mapping |
|---|---|---|---|---|
| **JOURNEY-START** | Start / hub | User needs context before the first action | Dumping the whole manual on step 1 | Hub screen: short prose + **one** primary CTA ([GOV.UK start](https://design-system.service.gov.uk/patterns/start-using-a-service/)) |
| **JOURNEY-LINEAR** | One path of steps | Strong order; each step is a form or decision | User must freely reorder many independent tasks | Sequence of form/detail URLs; “Continue” primary; back link to previous GET |
| **JOURNEY-TASKLIST** | Task list hub | Long service; evidence users cannot finish in one sitting **and** need to choose order ([GOV.UK task list](https://design-system.service.gov.uk/components/task-list/)) | You can simplify to fewer steps — try that first | List/settings-like hub of tasks + status `badge`; each task → its own route |
| **JOURNEY-BRANCH** | Branch on answer | Next step depends on prior input | Fake branches that always rejoin without reason | POST → 303 to different GET paths by server decision |
| **JOURNEY-CRUD** | List ↔ detail ↔ form | Admin/resource management | Treating every edit as a 12-step wizard | [Admin Resource](/recipes/admin-resource) pattern |
| **JOURNEY-QUEUE** | Work item → action → next | Ops processing | Using a wizard when a queue row action is enough | [Ops Queue](/recipes/ops-queue) |

## After submit: where do they land?

| Outcome | Land on | Feedback ID |
|---|---|---|
| Validation failed | **Same form URL** (422 body) | `FEED-VAL` |
| Created resource | **Detail** or list with new item visible | `FEED-OK-PAGE` and/or `FEED-OK-TOAST` |
| Updated settings | **Same settings GET** or parent section | `FEED-OK-TOAST` or `FEED-OK-PAGE` |
| Deleted | **List** parent (item gone) or explicit result | Prefer list + toast; confirm was `FEED-CONFIRM` |
| Finished multi-step | **Result / confirmation** screen | `FEED-OK-PAGE` — reference id + next links, not toast-only |
| Partial batch | **Result** listing per-item status | `FEED-PARTIAL` |

**Rule:** never end a successful mutate on a blank view with only a vanishing toast (`FEED-OK-PAGE` when the page *is* the confirmation).

## Step UX rules

1. **One job per step URL** — same as one screen type per URL ([Screens](/docs/screens)).
2. **Progress** — show step position or task status when there are 3+ steps (task list statuses or “Step N of M” text).
3. **Back** — always a way to a previous **GET** without resubmitting POST.
4. **Save & resume** — if users leave mid-flow, task list or draft on server; do not rely on client-only state.
5. **Do not use task list** only to show answers — GOV.UK: task list is for control over long services, not a summary of answers.

## Anti-patterns

- Wizard of 8 steps that could be one form with clear sections.
- POST without 303 to a clear next GET.
- Success = toast on `/` with no record of what was created.
- No back path except browser chrome.
- Task list when a single linear form is enough (GOV.UK: simplify first).

## Checklist (agents)

1. Name the journey shape (`JOURNEY-*`).
2. List ordered routes (GET/POST) and the **happy-path landing**.
3. Map each exit to a `FEED-*` id.
4. Define back/cancel targets.
5. Only then generate markup.

## See also

- [Screens](/docs/screens) · [Feedback](/docs/feedback) · [Data display](/docs/data-display) · [`/llms-ux.txt`](/llms-ux.txt)
