# Gelium UX Composition — Agent Contract

> **Artifact ID**: `gelium://ux-composition/agent-contract/v0.1` · **Version**: 0.1
> The prompt and contract an agent uses to compose Gelium screens (and the
> machine-readable artifacts that describe them) under the protocol.
> Builds on `docs/gelium-ui-agent-prompts.md` (the four existing workflows) and
> `AI-COMPONENT-IMPLEMENTER-PROMPT.md` (the operational component prompt). It
> does not replace them; it scopes the *composition* work to the protocol.

---

## 1. Scope of this contract

This contract governs **composition authoring**: producing or reviewing a screen
recipe, a flow, and their machine-readable contracts. It does NOT govern
component implementation (see `AI-COMPONENT-IMPLEMENTER-PROMPT.md`) or the
component/recipe/theme/documentary workflows (`gelium-ui-agent-prompts.md`).

## 2. Mandatory constraints (the agent MUST NOT)

1. **No new components.** Compose only components that exist in
   `docs/gelium-ui-component-registry.md` §2. A component not in the registry
   is a blocker, not an excuse to invent.
2. **No new patterns.** Reference only patterns in
   `docs/gelium-ui-pattern-registry.md` (D1–D8, E1–E19, F1–F14).
3. **No new wire contracts.** Use exactly: `422 + X-Loom-Validation: true`,
   `HX-Trigger {"loom:toast":…}`, GET stable params, `POST + 303 SeeOther`.
   No parallel APIs (P5).
4. **No AI / MCP / database / semantic search / new UI implementation.** This
   protocol is documentation + machine-readable contracts. Do not build
   features or components, and do not modify Go/CSS/templates/tests.
5. **No WhatsApp demo code changes.** The WhatsApp chat is evidence only.
6. **No client-side list state, no drag-and-drop, no JS for what a GET form
   solves** (anti-rules `composition-rules.md` §5).
7. **Do not edit shared files** (`agent-prompts.md` §6) and do not modify
   `docs/handoffs/*`, `AI-COMPONENT-IMPLEMENTER-PROMPT.md`,
   `COMPONENT-ROADMAP.md`, `MATERIAL-WEB-PROGRESS.md`.

## 3. Inputs the agent MUST read first

| Purpose | Document |
|---|---|
| Rules and anti-rules | `docs/gelium-ui-composition-rules.md` |
| Recipe template (19 fields) | `docs/gelium-ui-screen-composition.md` §3 |
| Implemented recipes | `docs/gelium-ui-screen-recipes.md` |
| Components | `docs/gelium-ui-component-registry.md` |
| Patterns | `docs/gelium-ui-pattern-registry.md` |
| Protocol model | [docs/gelium/ux-composition-model.md](../gelium/ux-composition-model.md) |
| Decision rules (R-*) | [docs/gelium/ux-decision-rules.md](../gelium/ux-decision-rules.md) |
| Schemas | [docs/gelium/schemas/screen-recipe.schema.json](../gelium/schemas/screen-recipe.schema.json), [docs/gelium/schemas/flow.schema.json](../gelium/schemas/flow.schema.json) |
| Validation | [docs/gelium/composition-validation.md](../gelium/composition-validation.md) |
| Accessibility | `docs/gelium-ui-accessibility-contract.md` |
| Content rules | `docs/gelium-ui-content-rules.md` |
| SEO/GEO | `docs/gelium-ui-seo-contract.md`, `docs/gelium-ui-geo-contract.md` |

## 4. The prompt (copy-paste)

```text
You are composing a Gelium screen (or flow) under the Gelium UX Composition
Protocol v0.1. Work in this repository only, documentation/machine-readable
artifacts only.

READ FIRST, in order:
1. docs/gelium-ui-composition-rules.md       (selection rules §4, anti-rules §5,
                                             state matrix §8, server rules §9)
2. docs/gelium-ui-screen-composition.md §3   (the 19-field recipe template)
3. docs/gelium-ui-screen-recipes.md          (3 implemented recipes as reference)
4. docs/gelium-ui-component-registry.md      (existing components — zero new)
5. docs/gelium-ui-pattern-registry.md        (existing patterns D/E/F)
6. docs/gelium/ux-composition-model.md       (Intent → … → Component)
7. docs/gelium/ux-decision-rules.md          (the R-* rule catalog)

PRODUCE (in this order):
1. A one-line INTENT and the FLOW it belongs to (flow.schema.json shape).
2. The SCREEN(s): one URL per screen; one primary task per surface.
3. The 19-field RECIPE for each screen.
4. A machine-readable contract instance conforming to
   schemas/screen-recipe.schema.json (and flow.schema.json for the flow).
5. A decision_rules array referencing only real R-* IDs from ux-decision-rules.md
   that the screen actually exercises.

CONSTRAINTS (non-negotiable):
- Compose only registry components and patterns. No new component, pattern,
  token or wire contract.
- No-JS end-to-end: the primary flow completes with JavaScript disabled;
  HTMX only enhances and never changes the mutation contract.
- Persistent-contextual feedback never uses loom:toast; validation is never a
  toast (422 + Inline alert + Validation summary).
- Every data surface declares empty, loading and error states with real state
  patterns (R-STATE-01).
- One primary action per surface; destructive is never primary (R-ACTION-01/02).
- Exactly one h1 per page; state never color-only; reduced motion and forced
  colors by design.
- Use only the canonical server contracts (422/X-Loom-Validation, loom:toast,
  GET stable params, POST + 303).
- English artifacts; UI strings server-rendered and localizable.

VALIDATE before reporting:
- jq empty on every JSON file you wrote.
- Validate every contract against its declared schema (see
  composition-validation.md §A/H).
- Confirm every decision-rule ID you reference exists in ux-decision-rules.md.
- Do NOT modify Go/CSS/templates/tests or any shared file.

REPORT: files written, the intent/flow trace, the decision rules exercised,
and the validation commands run with their output.
```

## 5. Review contract (for the orchestrator)

When reviewing protocol artifacts, check in this order:

1. **Consistency**: no contradiction with `composition-rules.md`,
   `screen-recipes.md`, the registries or the schemas.
2. **Rule integrity**: referenced `R-*` IDs exist; applicable rules are not
   silently skipped.
3. **Machine validity**: schemas and instances validate (see
   `composition-validation.md`).
4. **Scope**: docs/machine-readable only; nothing implemented; no new
   contracts or components.
5. **Coherence**: cross-links resolve; `flow_ref` ↔ `screens` links match.

## 6. Deliverable acceptance

A composition is accepted when it passes the gates of
`composition-validation.md` §I (schema conformance, zero new primitives,
complete no-JS branch, required states, real decision-rule references) and
satisfies the review contract above. Leave the changes uncommitted for the
orchestrator to inspect and commit.
