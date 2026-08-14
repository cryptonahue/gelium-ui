# Gelium UX Composition Protocol — Purpose and Principles

> **Artifact ID**: `gelium://ux-composition/protocol/v0.1` · **Version**: 0.1
> **Status**: draft for review (documentation + machine-readable contracts only).
> Upstream: `docs/gelium-ui-system-roadmap.md`, `docs/gelium-ui-composition-rules.md`,
> `docs/gelium-ui-screen-composition.md`, `docs/gelium-ui-ux-principles.md`,
> `docs/gelium-ui-accessibility-contract.md`, `docs/gelium-ui-content-rules.md`.

---

## 1. Purpose

The Gelium UX Composition Protocol makes the composition of a Gelium screen a
**verifiable process** instead of an implicit one. It defines:

1. **What we compose** — a model that starts one level above the screen: every
   screen exists to serve an *intent* through a *flow* (`ux-composition-model.md`).
2. **How we choose** — explicit decision rules with stable IDs
   (`ux-decision-rules.md`) that reuse the existing selection rules of
   `composition-rules.md` §4 and close its documented gaps.
3. **What we must produce** — the 19-field recipe (already canonical in
   `screen-composition.md` §3) and its machine-readable form
   (`schemas/screen-recipe.schema.json`).
4. **What we may not do** — the anti-rules and non-negotiables consolidated
   from `composition-rules.md` §5, `ux-principles.md` and
   `accessibility-contract.md` §0.
5. **How we prove it** — the validation checklist (`composition-validation.md`).

The protocol **consolidates** the existing 19-field screen recipe composition
and the registries; it does **not duplicate** them. Every rule that already
exists (pattern selection, server-driven contracts, state matrix) is referenced
from its canonical document. New normative content is limited to the explicit
deliverables of v0.1: dashboard, tabs, primary action, required states,
responsive behavior, accessibility, no-JS, errors and destructive actions.

## 2. Non-negotiable principles

The following principles are inherited from the system and are **not** open for
trade-off in any composition. Sources are cited; the protocol only restates the
binding form.

### P1. No-JS end-to-end (server-first)
The primary flow of every screen completes with JavaScript disabled. HTMX only
*enhances*; it never becomes the contract. Evidence: `ux-principles.md`,
`AI-COMPONENT-IMPLEMENTER-PROMPT.md` §12, every recipe's `NO_JS_FLOW`.

### P2. HTML-first, native before ARIA
Native semantics (button, a, form, select, table, dialog, progress, nav) come
before ARIA; no fake roles, no `div` acting as a control, no redundant ARIA.
Evidence: `accessibility-contract.md` §0.2, §1.1.

### P3. The URL is the state
Any navigable state (listing, filter, sort, page, selection, workflow
position) is a URL. GET with stable params for list state; POST + 303 for
mutations; no client-side list state. Evidence: `composition-rules.md` §9.6,
`ux-principles.md` P-recognition.

### P4. Persistent-contextual ≠ transient-action feedback
Persistent feedback (empty, inline alert, banner, callout, error state,
validation summary, persistent success) never travels through `loom:toast`;
transient action results never occupy a persistent slot. Evidence:
`state-patterns-audit.md:45`, `ux-patterns.md` "Cross-cutting rule".

### P5. Canonical wire contracts only
A composition reuses exactly the existing server-driven contracts: `422 +
X-Loom-Validation: true`, `HX-Trigger {"loom:toast":…}`, GET stable params,
`POST + 303 SeeOther`. No new contract is invented by a recipe. Evidence:
`composition-rules.md` §9.

### P6. Recipes are 100% wiring
A screen recipe composes only components, patterns and tokens that already
exist in the registries. A new primitive requires the component-ingress
checklist of `component-registry.md` §5 and the rationale gate of
`composition-rules.md` §11. Evidence: `screen-composition.md` (Admin Resource
decomposition), `screen-recipes.md` (Phase G).

### P7. One primary action per surface
Every surface declares exactly one primary action; destructive is never the
primary action. Evidence: `ux-principles.md` §8, `composition-rules.md` §3.1.

### P8. State is never color-only
Every state carried by color (`:checked`, `aria-sort`, `aria-current`,
`aria-invalid`, `aria-busy`, disabled, tones) is also carried by semantics or
visible text. Evidence: `accessibility-contract.md` §0.3.

### P9. Accessibility is by design
WCAG 2.1 AA, native focus management, exactly one `h1` per page, document order
= visual order, reduced motion and forced colors first-class. Evidence:
`accessibility-contract.md`.

### P10. Tokens `--ui-*`, zero literals
Every visual value in a composition is a public token; recipes never hardcode
color or geometry literals (guard `TestNoColorLiteralsInComponents`). Evidence:
`theme-contract.md`, `agent-prompts.md` §1.3.

### P11. English artifacts, localizable UI
System technical artifacts are English. UI strings are server-rendered data,
localizable by the consuming project without touching components. Evidence:
`content-rules.md` §9.

## 3. Scope of v0.1

- **In scope**: consolidation (19-field recipe, registries), the
  Intent→Flow model, the decision rules of the requested topics, machine-
  readable contracts for Admin Resource / Inbox·Conversation / Public Marketing
  Page, validation checklist, agent contract. See `README.md`.
- **Out of scope**: AI, MCP, database, semantic search, new UI components, new
  wire contracts, implementation code (Go/CSS/templates/tests), WhatsApp demo
  code.

## 4. Artifact chain

```text
ux-composition-model.md        (Intent → User → Task → Flow → Screen → Recipe → Pattern → Component)
        │
ux-decision-rules.md           (normative R-* catalog; references composition-rules.md §4/§5/§8/§9)
        │
schemas/screen-recipe.schema.json   ─┐  machine-readable recipe (19 fields + protocol fields)
schemas/flow.schema.json            ─┴  machine-readable flow (Intent → screens)
        │
contracts/{admin-resource,inbox-conversation,public-marketing-page,inbox-flow}.json
        │
composition-validation.md      (gates: schema conformance, rule compliance, a11y, no-JS, SEO/GEO)
        │
agent-contract.md              (the prompt that applies all of the above)
```

## 5. Definition of done (v0.1)

- Purpose and principles written and consistent with the existing contracts.
- Model documents the eight layers with concrete mappings to real concepts and
  files.
- Decision rules reference, not duplicate, `composition-rules.md`; the gaps
  requested by the deliverables (dashboard, tabs, primary action, required
  states, responsive, accessibility, no-JS, errors, destructive) are explicit.
- Machine-readable files are syntactically valid and validate against their
  declared schemas; links and rule references are coherent.
- No application code, no new components/contracts, no contradictory content.
