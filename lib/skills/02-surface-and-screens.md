# Skill: Surface mode and screen types

Pick the SURFACE mode and one SCREEN type before writing markup. This decides the
layout blocks, not aesthetics.

## SURFACE modes

| Mode | Bias | When |
| --- | --- | --- |
| Operate | task, scan, table/form | default for admin/tools |
| Read | 65ch prose, structure | articles, docs, essays |
| Persuade | one CTA, hub/start | marketing/landing (still tokens, not stock gradients) |
| Experience | artifact first | distinctive showcase (not fake CRUD) |

Default for admin and tool screens is **Operate**.

## Vocabulary resolution (do this before creating any NEW component)

Before naming or building a new component, resolve the user's informal
description against `lib/ui-vocabulary.md` (grep it for keywords: "the dark
see-through layer behind a popup" → scrim; "the gray text inside the box" →
placeholder). Then check `docs/gelium-ui-vocabulary.md` §8 for already-resolved
naming conflicts.

- Canonical name exists in the registry → compose the registered component.
- Name is an alias → use the canonical name (e.g. snackbar → toast,
  facepile → avatar group).
- Only escalate a genuinely new primitive if neither file resolves it.

## Screen types → blocks

| Type | Blocks |
| --- | --- |
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

0. Brief gate: read consumer `PRODUCT.md`/`DESIGN.md` when present. If either
   is absent or incomplete, run the plain-language brief in `llms-ux.txt`.
   Accept `Unknown`/`To decide`, record assumptions, and stop only when an
   unresolved decision would change page or flow architecture. Do not invent
   product intent, chrome, or visual direction silently.
1. USER JOB → one sentence.
2. SURFACE mode + SCREEN type.
3. ONE primary action.
4. JOURNEY-*if multi-step; DATA-* if collection; FEED-* for feedback.
5. States (empty/error/loading/success).
6. Shell/chrome inventory — ask explicitly if no user decision exists:
   - Does this product have a footer? What carries it (legal / contact / links)?
   - Top bar vs side navigation? (Decide once per product, not per screen.)
   - Density by surface (Operate tables denser than Read prose)?
   If no explicit user decision exists, ASK — this is part of the artifacts
   gate above. Do not invent chrome silently.

## Multi-screen journeys

Source: `handbook-journeys.md` (`docs/journeys` in the consumer repo).

- **Start/hub pattern**: every journey starts at a hub that gives context,
  then offers exactly ONE primary CTA into the flow.
- **Step order**: fixed steps are strongly ordered (each step's only forward
  path is the next step); optional/free steps are freely navigable from the hub.
- **Post-submit landing**:

  | Outcome | Landing |
  | --- | --- |
  | validation fail | same form re-rendered with **422** + errors |
  | success | **303** redirect to detail/result page |

- **Resume rule**: navigable state = URL. If a user can resume mid-journey,
  that state must be reachable by URL (query params or path), never JS-only.

## Data display choice

| Display | Use for |
| --- | --- |
| table | tabular comparison data (same fields across many rows) |
| cards | browsable rich items (varied content per item) |
| list | scannable homogeneous items |

## Anti-slop

Avoid: nested cards for simple forms; stock purple-blue heroes; an icon tile above
every H1; gray on saturated fills; bounce motion; "bolder for its own sake" on
Operate tables; a new font per feature; layout tables for non-data.
