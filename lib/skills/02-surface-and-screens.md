# Skill: Surface mode and screen types

Pick the SURFACE mode and one SCREEN type before writing markup. This decides the
layout blocks, not aesthetics.

## SURFACE modes

| Mode | Bias | When |
|---|---|---|
| Operate | task, scan, table/form | default for admin/tools |
| Read | 65ch prose, structure | articles, docs, essays |
| Persuade | one CTA, hub/start | marketing/landing (still tokens, not stock gradients) |
| Experience | artifact first | distinctive showcase (not fake CRUD) |

Default for admin and tool screens is **Operate**.

## Screen types → blocks

| Type | Blocks |
|---|---|
| hub | title, short context, **one primary button** |
| list | data-table OR list/card, filters, pagination, empty-state |
| detail | title, status badge, body, secondary actions |
| form | fields, validation-summary, primary + cancel |
| confirm | dialog/confirm route, destructive button, cancel |
| settings | list rows, switch/select |
| queue | list/table, badge, row POST actions |
| result | success banner/alert, next links, optional toast |

## One primary action

Each screen has exactly **one** primary action (highlighted button). Everything
else is secondary/link.

## Workflow before drawing

1. USER JOB → one sentence.
2. SURFACE mode + SCREEN type.
3. ONE primary action.
4. JOURNEY-* if multi-step; DATA-* if collection; FEED-* for feedback.
5. States (empty/error/loading/success).

## Anti-slop

Avoid: nested cards for simple forms; stock purple-blue heroes; an icon tile above
every H1; gray on saturated fills; bounce motion; "bolder for its own sake" on
Operate tables; a new font per feature; layout tables for non-data.
