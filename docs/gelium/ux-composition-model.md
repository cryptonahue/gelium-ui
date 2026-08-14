# Gelium UX Composition — The Model

> **Artifact ID**: `gelium://ux-composition/model/v0.1` · **Version**: 0.1
> The composition model of the protocol: the eight layers that turn a product
> intent into concrete Gelium components, and the invariants that govern the
> chain.

---

## 1. The chain

```text
Intent → User → Task → Flow → Screen → Recipe → Pattern → Component
```

Read **top-down** (analysis): an *Intent* is pursued by a *User* performing a
*Task*, delivered through a *Flow* of *Screens*; each *Screen* is specified by a
*Recipe*; each *Recipe* is a composition of *Patterns*; each *Pattern* is a
composition of *Components*.

Read **bottom-up** (composition): *Components* are composed into *Patterns*,
*Patterns* are wired by a *Recipe*, the *Recipe* specifies a *Screen*, the
*Screen* participates in a *Flow*, the *Flow* fulfils a *Task* for a *User* and
realises an *Intent*.

This is the **same dependency direction** already established for the system
(`docs/gelium-ui-dependency-metadata.md` §1: core → components → patterns →
recipes), extended two layers upward (screen → flow, flow → intent).

## 2. Layer definitions

| # | Layer | Definition | Maps to (existing) | Example (Inbox) | Example (Marketing) |
|---|---|---|---|---|---|
| 1 | **Intent** | The product-level goal the user is trying to achieve. One intent drives the design of one or more flows. | New in the protocol (not previously formalized). | "Resolve inbound conversations." | "Convert visitors into users." |
| 2 | **User** | The actor who pursues the intent, described by their primary task — not by the system role. | Screen grammar "usuario" (`composition-rules.md` §2); 19-field `USER` (`screen-composition.md` §3). | Support agent/operator. | Prospective customer/visitor. |
| 3 | **Task** | What the user actually does: one primary task per surface, plus secondary tasks. | 19-field `PRIMARY_TASK` / `SECONDARY_TASKS`. | Read context, reply, triage. | Understand the offer, take the CTA. |
| 4 | **Flow** | An ordered, URL-navigable sequence of screens that fulfils the task. The flow owns the navigation contracts between screens (GET links, POST + 303, 422). | New in the protocol; concretized in `schemas/flow.schema.json`. | Inbox list → Conversation → (reply/close). | Marketing home (single-screen flow in v0.1). |
| 5 | **Screen** | One navigable state (one URL) with a single primary task per surface. The page is the unit of URL; overlays never own flow state. | Screen grammar "superficie"; surface rules (`composition-rules.md` §3); "página = unidad de URL". | `/inbox/{conversation_id}` | `/` (marketing home) |
| 6 | **Recipe** | The 19-field specification that fully defines a screen: surface, user, tasks, pattern, vocabulary, components, states, a11y, content, SEO/GEO, server contract, no-JS, HTMX, responsive, theme, alternatives, rationale. | 19-field template (`screen-composition.md` §3), implemented recipes (`screen-recipes.md`). | `contracts/inbox-conversation.json` | `contracts/public-marketing-page.json` |
| 7 | **Pattern** | A reusable composition of components over a server contract (state D1–D8, UX E1–E19, public F1–F14). | `pattern-registry.md`, `ux-patterns.md`, `public-content-patterns.md`. | Feed (thread) + Form + Error recovery + Notifications. | Hero (F2), Split (F13), Feature Card (F7), Newsletter (F10), CTA Link (F6). |
| 8 | **Component** | A real partial + CSS + optional view-model from the registry. Recipes only wire existing components. | `component-registry.md` §2; dependency metadata §2. | Avatar, List, Text field, Button, Badge, Toast, Empty state. | Hero, Split, Feature Card, Card, Button (link), Newsletter, Footer, Video, Breadcrumb. |

## 3. Invariants

1. **One intent, one or more flows.** An intent may span several flows (e.g.
   "manage resources" → browse flow, create flow, edit flow), but a flow must
   trace back to exactly one intent.
2. **One primary task per surface.** If a screen has two primary tasks, it
   splits into two screens (`composition-rules.md` §3.1). Flows absorb screen
   boundaries; a primary task is never shared across two screens.
3. **A screen is a URL.** Navigable state is a page (deep-linkable, back
   works); overlays are short, reversible sub-tasks that retain context and do
   not own flow state (`composition-rules.md` §3.2, §4.7).
4. **Dependency direction is one-way.** Components never depend on patterns,
   patterns never on recipes, recipes never on flows. A recipe composes the
   layers below it and nothing above.
5. **Recipes are 100% wiring.** Pattern and component references must exist in
   the registries. A gap (e.g. `dashboard` states in the state matrix) is a
   documented GAP, never silently re-invented by a recipe.
6. **Contracts are shared, not invented.** Every screen/flow uses the canonical
   server contracts (P5 in the protocol). Flow navigation between screens is
   GET links for reading and POST + 303 for mutations.

## 4. Screen and flow contracts

- A **Screen** is machine-describable by `schemas/screen-recipe.schema.json`
  (the 19 fields + intent/flow/decision-rule references). Screens are the unit
  of the recipe template.
- A **Flow** is machine-describable by `schemas/flow.schema.json`: intent, user,
  task, entry screen, ordered screens, navigation edges with their contracts,
  success state and no-JS/HTMX notes.
- The relationship is **flow → screens**: a flow references screens by `$id`;
  a screen references its flow back via `flow_ref`. Both are validated.

## 5. Worked traces

### 5.1 Inbox / Conversation (flow → screen)

```text
Intent  : "Resolve inbound conversations."
User    : support agent/operator
Task    : read context → reply/triage (primary); filter, refresh (secondary)
Flow    : Inbox list → Conversation        (see contracts/inbox-flow.json)
Screen  : /inbox/{conversation_id}         (contracts/inbox-conversation.json)
Recipe  : 19 fields; pattern = Feed (thread) + Form(reply) + Error recovery + Notifications
Patterns: E9, E15, F? none public; vocabulary Feed + Form
Server  : GET page; POST reply → 303 (+ 422 validation); refresh POST-only + loom:toast
Components: avatar, badge, list, text-field, button, toast, empty-state, banner, error-state
```

The thread composition mirrors the WhatsApp chat pattern already in the system
(`internal/app/demo_whatsapp.go`, `web/templates/demo-whatsapp.html`), which is
evidence, not new implementation.

### 5.2 Public Marketing Page (single-screen flow)

```text
Intent  : "Convert visitors into users."
User    : prospective customer/visitor
Task    : understand the offer and take the primary conversion CTA
Flow    : Marketing home (single screen in v0.1)
Screen  : / (contracts/public-marketing-page.json)
Patterns: F2 Hero, F13 Split, F7 Feature Card, F6 CTA Link, F10 Newsletter, F8 Footer, F12 Section Heading, F1 Article, F3 Breadcrumb, F14 Video
Components: hero, split, feature-card, card, button (link), newsletter, footer, breadcrumb, video, language-switcher
```

Evidence in the real system: the marketing home is implemented by
`internal/app/landing.go` (`marketingLanding()`) composing the same public
patterns, and it is the indexable surface (`index, follow`, clean canonical,
JSON-LD) per `internal/app/server.go` `resolveMeta`/`jsonLDBreadcrumb`.

### 5.3 Admin Resource (existing recipe, machine-readable)

The implemented recipe (`screen-recipes.md` §1) is the reference for the
machine-readable contract `contracts/admin-resource.json`; no behavior changes.

## 6. Consistency with the existing composition flow

The 6-step composition flow of `screen-composition.md` §1 stays intact. The
protocol adds two steps at the front (state the intent; trace the flow) and
attaches machine-readable artifacts at the end:

```text
0. State INTENT and trace the FLOW (model, this doc).
1. Choose the pattern with the screen grammar (composition-rules §2).   [unchanged]
2. Select the pattern from the registry (composition-rules §4).          [unchanged]
3. Decompose into registry components — zero new primitives.             [unchanged]
4. Resolve states with the state patterns.                               [unchanged]
5. Fix the server contract (composition-rules §9).                       [unchanged]
6. Complete the 19 fields AND emit the machine-readable contract
   (screen-recipe.schema.json) + decision-rule references.               [added by protocol]
```

## 7. Definition of done

- The eight layers map 1:1 onto existing concepts and files (no invented
  vocabulary beyond Intent/Flow, which are the requested additions).
- The dependency direction matches `dependency-metadata.md` §1.
- The model is demonstrated by at least one flow instance and three screen
  instances that validate against the schemas.
