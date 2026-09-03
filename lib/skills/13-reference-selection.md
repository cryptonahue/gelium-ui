# Skill: Semantic reference selection

Use after selecting `SURFACE` and `SCREEN`, before page/section architecture and the wireframe. References are design evidence, not a moodboard or a replacement for Gelium contracts.

## Procedure

1. Extract the semantic needs from the request and the selected screen:
   page sections, local components, user actions, data shape, states, and recovery.
2. Match each need against published section references and the registered component vocabulary/registry.
3. Run a discovery pass over the reference before filtering it. Enumerate visible affordances and composition choices: what opens the detail, which actions deserve a glyph, search/discovery entry points, account/menu treatment, shell regions, responsive transitions, media controls, and recovery paths. Do not omit a candidate because it is not in the current contract; classify it first.
4. For every candidate, record the reference ID/name and why it applies. Then classify it as existing contract, B alert, C improvement, product fork, or rejected/no-match.
5. If a matching reference is available, consult it before the wireframe. If browser access, source verification, or the reference itself is unavailable, fail open and record the reason.
6. Apply the product filter: keep only structure and decisions that fit the product, data, audience, permissions, no-JS contract, and Gelium tokens/primitives.
7. Put the selected references, discovered candidates, and rejected candidates in the plan/G5 packet. They inform A/B/C; they do not silently add scope.

## Semantic audit lenses

Before accepting a reference or existing implementation, run these checks:

- **Identity/account:** decide whether `Me`, avatar, account menu, settings, and logout are one cluster or separate product destinations. Do not duplicate account actions in the page body when the shared shell owns them.
- **Actions:** distinguish navigation from mutation. A repeated action may use a catalog glyph (for example bookmark) only when its accessible name, state, POST/GET contract, and touch target remain clear.
- **Constant state:** challenge labels or badges rendered on every item. If a state is true for every item in the current surface, remove the repeated badge or move the explanation to the section context; do not spend row space proving a constant.
- **Visibility semantics:** show public/private/followers status only where visibility can differ and the viewer is allowed to know it (typically an owner profile or management surface). Use a catalog globe/lock plus visible text when the distinction matters; never imply private content is present in a public-only feed.
- **Detail entry:** choose one clear reading target. Avoid a whole-card link when nested forms or controls exist; make the reading region link to the canonical detail route.

A mismatch with the current contract becomes a visible B, C, or product fork. It is not silently discarded merely because the existing template does it.

- Section references guide page composition; component references guide local behavior, anatomy, and states.
- Gelium's registered primitives, server/no-JS contracts, accessibility rules, and real product data are authoritative.
- Never copy another product's logo, brand, assets, screenshots, motion, visual skin, or unverified commercial claims.
- A reference cannot justify a fake avatar, invented price, fabricated media metadata, or an unregistered route.
- Operate/admin screens do not use inspiration galleries as a substitute for product requirements.

## G5 handoff

The packet should include a compact note like:

```text
References: REF-PRICING + ui-data-table — pricing comparison detected; both passed product filter.
No-match: no published reference for billing-history; use registered list/detail patterns.
Rejected: REF-HERO — no landing promise in this screen.
```

References are mandatory inputs only when semantically matched and available. They are never a blocking prerequisite for continuing the work.
