# REF-CARD-DETAIL-ENTRY

Use when a list item or card mixes **navigation** (open detail) and **mutation**
(like, save, delete, react) or other nested interactive controls.

## Problem

A whole-card `<a>` wrapping buttons/forms breaks nested interactive HTML, confuses
focus order, and often creates duplicate “View post” CTAs next to a linked title.

## Pattern

1. Choose **one** reading region that opens the canonical detail route:
   title, and optionally context line + summary + safe non-interactive media.
2. Keep mutation controls **outside** that link (sibling forms/buttons).
3. Remove secondary “View / Open / Ver publicación” links when the reading
   region already goes to the same URL.
4. Do not put POST controls inside the reading `<a>`.

## Checklist before wireframe lock

- [ ] Nested interactive elements are not inside the detail link
- [ ] Keyboard focus order: reading link, then each control, in DOM order
- [ ] Touch targets ≥ product minimum on action controls
- [ ] GET for read, POST for mutations (CSRF when required)
- [ ] Mobile order still meta → title → body → media → actions
- [ ] Constant “open detail” affordance is not repeated three ways

## Gelium filter

Registered card/list primitives; real detail route only; no invented paths.
If the product has no detail route, do not fake a link — show plain text.
