# Patterns

Composition patterns for Gelium UI: how to assemble screens for common product shapes. Components stay in the sidebar reference; **recipes** are runnable full screens. This page is the **domain skeleton** layer — enough structure for humans and agents to stop improvising entire apps from a button catalog.

## How to use this page

1. Pick a skeleton that matches the product job.  
2. Expand each screen with [Screens](/docs/screens), [Journeys](/docs/journeys), [Data display](/docs/data-display).  
3. Apply [Feedback](/docs/feedback) `FEED-*` on every collection and form.  
4. Prefer an existing [recipe](/recipes/admin-resource) when it already matches.

## Domain skeletons

### SKEL-FORUM (discussion)

| Screen | Type | Data | Notes |
|---|---|---|---|
| Topic list | list | `DATA-LIST` or `DATA-FEED` | Title, meta (author, count, time), pagination |
| Topic detail / thread | detail + list | posts as `DATA-LIST` | Identity → posts → composer |
| Reply / new topic | form | — | `FEED-VAL` on 422; 303 to thread |
| Empty | — | `FEED-EMPTY` | CTA: create first topic |

**Primary actions:** list → “New topic”; thread → “Reply”.  
**Journey:** `JOURNEY-CRUD`-like between list and detail; not a long wizard.

### SKEL-CATALOG (store browse)

| Screen | Type | Data | Notes |
|---|---|---|---|
| Catalog | list | `DATA-CARDS` or `DATA-LIST` | Filters on GET query |
| Product detail | detail | `DATA-DESC` + actions | Price/status near title |
| Cart (simple) | list | `DATA-LIST` | Qty fields = form controls |
| Checkout | form / linear journey | — | `JOURNEY-LINEAR`; result page after pay |
| Admin products | list | `DATA-TABLE` | [Admin Resource](/recipes/admin-resource) |

**Do not** use a dense data table as the public catalog home.

### SKEL-ADMIN-RESOURCE

Already implemented as recipe **Admin Resource**: table + filters + form + dialog + banner.  
Use for any back-office entity CRUD.

### SKEL-OPS-QUEUE

Recipe **Ops Queue**: status badges, row POST+303, work-list journey (`JOURNEY-QUEUE`).

### SKEL-SETTINGS

| Screen | Type | Notes |
|---|---|---|
| Settings hub | settings / list | Groups of rows |
| Section form | form | One group per URL when large |
| Save | — | `FEED-OK-TOAST` or stay with `FEED-OK-PAGE` |

### SKEL-ONBOARDING

| Screen | Type | Journey |
|---|---|---|
| Start | hub | `JOURNEY-START` |
| Steps | form×N | `JOURNEY-LINEAR` or `JOURNEY-TASKLIST` if resumable |
| Done | result | `FEED-OK-PAGE` + next CTA |

## Runnable recipes

| Recipe | Skeleton affinity |
|---|---|
| [Admin Resource](/recipes/admin-resource) | SKEL-ADMIN-RESOURCE |
| [Ops Queue](/recipes/ops-queue) | SKEL-OPS-QUEUE |
| [Public/Social Feed](/recipes/public-feed) | SKEL-FORUM list/feed slice |

## Anti-patterns

- Inventing a new CSS layout language per domain instead of DATA-* + screen types.
- Forum thread as a single infinitely nested table.
- Checkout as one toast without a result page.

## See also

- [UI definition of done](/docs/ui-definition-of-done) · [`/llms-ux.txt`](/llms-ux.txt)
