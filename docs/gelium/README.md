# Gelium UX Composition Protocol — v0.1

> The composition protocol of the Gelium UI system: how an *intent* becomes a
> *flow* of *screens*, and how every screen is composed from the existing
> registries, decision rules and machine-readable contracts.
>
> **Status**: v0.1 (draft for review). Documentation and machine-readable
> contract artifacts only. It does NOT implement AI, MCP, database, semantic
> search or new UI components, and it does not modify application code.
>
> **Scope boundary**: this protocol *describes* and *validates* composition; it
> never invents components, patterns, tokens or wire contracts. Everything it
> references already exists in the system (`docs/gelium-ui-*.md`, the component
> registry, the pattern registry and the screen recipes).

---

## What this protocol is

Gelium UI already defines how a screen is composed from registries
(`docs/gelium-ui-screen-composition.md`), the selection rules
(`docs/gelium-ui-composition-rules.md`) and three implemented screen recipes
(`docs/gelium-ui-screen-recipes.md`). The **Gelium UX Composition Protocol
v0.1** consolidates that work into a single, agent-readable, machine-verifiable
contract:

1. A **purpose + principles** statement (`ux-composition-protocol.md`).
2. The **composition model** `Intent → User → Task → Flow → Screen → Recipe →
   Pattern → Component` (`ux-composition-model.md`), which adds the two layers
   that precede a screen (Intent, Flow) on top of the existing grammar.
3. **Explicit UX decision rules** (`ux-decision-rules.md`) covering data
   surfaces, overlays, feedback, tabs, primary action, required states,
   responsive behavior, accessibility, no-JS, errors and destructive actions.
   Existing selection rules are referenced, not duplicated; only the gaps
   (dashboard, tabs, action hierarchy, …) are expanded.
4. **Machine-readable contracts** (`schemas/` + `contracts/`) for three screen
   types — Admin Resource, Inbox / Conversation and Public Marketing Page.
5. A **composition validation checklist** (`composition-validation.md`).
6. An **agent-facing prompt/contract** (`agent-contract.md`).

## Relationship to existing documents

| Existing document | Role | How the protocol relates |
|---|---|---|
| `docs/gelium-ui-composition-rules.md` | Pattern selection rules (§4), anti-rules (§5), state matrix (§8), server-driven rules (§9) | **Authoritative source** for the selection and anti-rules. The decision rules reference it instead of re-stating it. |
| `docs/gelium-ui-screen-composition.md` | 6-step composition flow, real Admin Resource decomposition, 19-field template (§3), recipe × component matrix (§4) | The 19-field template is the **canonical recipe contract**; `schemas/screen-recipe.schema.json` is its machine-readable form. |
| `docs/gelium-ui-screen-recipes.md` | The 3 implemented recipes (Admin Resource, Ops Queue, Public Feed) | `contracts/admin-resource.json` is a machine-readable derivation of recipe §1. |
| `docs/gelium-ui-component-registry.md` | Canonical component inventory | Components referenced by contracts MUST exist here. |
| `docs/gelium-ui-pattern-registry.md` | State (D), UX (E) and public (F) patterns | `ux_pattern` fields reference these IDs (D1–D8, E1–E19, F1–F14). |
| `docs/gelium-ui-agent-prompts.md` | Agent workflows (component, recipe, theme, documentary) | `agent-contract.md` is the protocol-specific prompt built on it. |
| `docs/gelium-ui-system-roadmap.md` | System phases A–J and residuals | The protocol is a post-A–J consolidation slice (see roadmap "Próximo slice", optional expansion). |

## File map

```text
docs/gelium/
├── README.md                        this index + scope + versioning
├── ux-composition-protocol.md       purpose, principles, non-negotiables
├── ux-composition-model.md          Intent → User → Task → Flow → Screen → Recipe → Pattern → Component
├── ux-decision-rules.md             normative decision rules (R-* catalog)
├── composition-validation.md        validation checklist and gates
├── agent-contract.md                agent-facing prompt + contract
└── schemas/
    ├── screen-recipe.schema.json    $id gelium://ux-composition/screen-recipe/v0.1
    ├── flow.schema.json             $id gelium://ux-composition/flow/v0.1
    └── contracts/
        ├── admin-resource.json          screen contract (Admin Resource)
        ├── inbox-conversation.json      screen contract (Inbox / Conversation)
        ├── public-marketing-page.json   screen contract (Public Marketing Page)
        ├── admin-resource-flow.json     flow contract (Admin Resource)
        ├── inbox-flow.json              flow contract (Inbox)
        └── public-marketing-flow.json   flow contract (Marketing Home)
```

## Versioning and schema IDs

- Protocol version: **0.1**. The version lives in the filename, in each
  document header and in the `version` field of every contract instance.
- Schema `$id` values are **opaque identifiers**, not resolvable URLs:
  `gelium://ux-composition/<artifact>/v0.1`. They exist to pin the exact
  contract a document or instance claims to satisfy. A breaking change to the
  protocol bumps the major version and the `$id` (e.g. `/v0.2`).
- Instances declare the schema they conform to via `schema` (exact `$id`).
  Validating an instance against its declared schema is part of the
  validation checklist (`composition-validation.md` §A).
- Decision rules are versioned with the protocol; each rule has a stable
  `R-<AREA>-<NN>` ID that contracts reference (`decision_rules` field).

## Reading order

1. `ux-composition-protocol.md` — purpose and principles (what and why).
2. `ux-composition-model.md` — the eight-layer model (how screens are reached).
3. `ux-decision-rules.md` — the normative rules (what is allowed/required).
4. `schemas/screen-recipe.schema.json` + `schemas/flow.schema.json` — the
   machine-readable shape of a screen and a flow.
5. `contracts/*.json` — worked, validated examples.
6. `composition-validation.md` — how to prove a composition is compliant.
7. `agent-contract.md` — the prompt an agent uses to apply all of the above.

## v0.1 scope — in and out

**In scope**
- Consolidation of the existing 19-field recipe template and the registries
  into one protocol with explicit decision rules and machine-readable
  contracts.
- New normative content ONLY where the system has a documented gap or where an
  explicit rule is required by the deliverables (dashboard selection, tabs,
  primary action, required states, responsive, accessibility, no-JS, errors,
  destructive actions).
- Three screen contracts (Admin Resource, Inbox / Conversation, Public
  Marketing Page) and three flow contracts (one per screen type), all
  referencing real components, patterns and contracts.

**Out of scope (v0.1)**
- AI features, MCP servers, databases, semantic search, new UI components or
  new wire contracts. Contracts only reference components and contracts that
  already exist.
- Implementation work: this slice is documentation + machine-readable
  artifacts only. No Go, CSS, templates or tests are added or modified.
- WhatsApp demo implementation code. The WhatsApp chat pattern is referenced
  only as evidence for the Inbox/Conversation composition (`demo-whatsapp*`).
- `docs/handoffs/*`, `AI-COMPONENT-IMPLEMENTER-PROMPT.md`,
  `COMPONENT-ROADMAP.md`, `MATERIAL-WEB-PROGRESS.md` are not modified.

## Validation performed

See `composition-validation.md`. Every JSON schema and contract instance is
validated syntactically and against its declared schema; every internal link
and rule reference is checked.
