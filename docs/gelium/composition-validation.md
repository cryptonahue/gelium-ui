# Gelium UX Composition — Validation Checklist

> **Artifact ID**: `gelium://ux-composition/validation/v0.1` · **Version**: 0.1
> How to prove that a screen composition (or a whole flow) is compliant with
> the Gelium UX Composition Protocol v0.1. Used by the agent contract
> (`agent-contract.md`) and by the orchestrator when reviewing protocol
> artifacts.
> Related: `composition-rules.md` §11 (rationale gate), `component-registry.md`
> §5 (ingress checklist), `agent-prompts.md` §5 (documentary workflow).

---

## A. Machine-readable contract checks

1. **Syntax.** Every `.schema.json` and every contract instance is valid JSON
   (e.g. `jq empty <file>` or `python3 -m json.tool`).
2. **Schema conformance.** Every contract instance validates against the schema
   declared in its `schema` field (`screen-recipe.schema.json` /
   `flow.schema.json`), e.g. `python3 -m jsonschema -i <instance> <schema>`.
3. **IDs.** `protocol` is `gelium-ux-composition`; `version` matches the schema
   (`0.1`); `slug` is `[a-z0-9-]+`; `$id`s are `gelium://ux-composition/<artifact>/v0.1`.
4. **Decision-rule references.** Every ID in `decision_rules` exists in the
   rule catalog (`ux-decision-rules.md`), regex `^R-[A-Z0-9]+-\d+$`.
5. **Schema references.** `flow_ref` in a screen matches the `$id` (or slug) of
   a real flow contract; a flow's `screens` entries reference real screen
   contracts.

## B. Composition checks (19 fields)

For every screen contract, verify the 19 fields of the recipe template
(`screen-composition.md` §3) are present and coherent:

- [ ] SURFACE, USER, PRIMARY_TASK, SECONDARY_TASKS present.
- [ ] UX_PATTERN references exist in `pattern-registry.md` (D1–D8, E1–E19, F1–F14).
- [ ] COMPONENTS reference only components in `component-registry.md` §2.
- [ ] Zero new primitives: any component not in the registry fails the check
      unless a documented, approved ingress (`component-registry.md` §5).
- [ ] STATES are a subset of the closed set (rest/hover/focus/pressed/selected/
      disabled/empty/loading/error/success) and every data surface declares
      empty + loading + error (R-STATE-01).
- [ ] SERVER_CONTRACT uses only canonical contracts (P5): 422 + `X-Loom-Validation`,
      `loom:toast`, GET stable params, POST + 303. No invented contract.
- [ ] NO_JS_FLOW is a complete, walkable branch (no JS required).
- [ ] HTMX_ENHANCEMENT never changes the mutation contract (R-NOJS-03).
- [ ] ALTERNATIVES_REJECTED names the patterns/anti-rules considered and why
      they were rejected (against `composition-rules.md` §4/§5).
- [ ] RATIONALE answers the six questions of `composition-rules.md` §11.

## C. Decision-rule compliance

Walk each applicable rule and confirm the composition satisfies it:

- Data surfaces: R-DATA-01…07 (table/list, list/queue, queue/board, card/panel,
  feed/collection, timeline/activity, dashboard).
- Overlays: R-OVERLAY-01…03 (dialog vs page, server fallback, confirm).
- Feedback: R-FEEDBACK-01…03 (persistent ≠ transient, validation never toast,
  channel selection).
- Navigation: R-NAV-01 (tabs), R-NAV-02 (URL is state).
- Actions: R-ACTION-01 (one primary), R-ACTION-02 (destructive not primary).
- States: R-STATE-01…04. Responsive: R-RESP-01…02.
- Accessibility: R-A11Y-01…07. No-JS: R-NOJS-01…03.
- Errors: R-ERROR-01…04. Destructive: R-DESTRUCT-01.
- SEO/GEO: R-SEO-01, R-GEO-01. Content: R-CONTENT-01…05.

The screen's `decision_rules` field lists the rules the composition explicitly
claims; the reviewer must confirm at least the applicable set (no applicable
rule left unlisted without rationale).

## D. Accessibility checklist

- [ ] One `h1` per page; headings descend (R-A11Y-03).
- [ ] All controls have accessible names; icon-only controls carry `aria-label`;
      decorative SVG `aria-hidden` (R-A11Y-02).
- [ ] Form controls: native `<label for>`; errors `aria-invalid` +
      `aria-describedby` (R-A11Y-01).
- [ ] Overlays: native `<dialog>` focus trap; server-rendered fallback
      available (R-OVERLAY-02, R-A11Y-04).
- [ ] Tabs: links + `aria-current`, no `tablist` (R-NAV-01).
- [ ] Data table: `aria-sort`, `aria-current="page"` pagination, native
      checkboxes (R-A11Y-05).
- [ ] Landmarks match roles; document order = visual order (R-A11Y-06).
- [ ] Reduced motion + forced colors verified (R-A11Y-07).
- [ ] No color-only state anywhere (R-STATE-04).

## E. No-JS and HTMX checklist

- [ ] Walk the primary flow with JS/HTMX disabled end-to-end (R-NOJS-01).
- [ ] Mutations are `POST + 303`; GET on a POST-only path returns 405 (R-NOJS-02).
- [ ] `HX-Request: true` returns the fragment; without it, the full page
      (R-NOJS-03).
- [ ] Transport error (network/5xx) shows a transient generic toast; nothing
      fails silently (R-ERROR-03).

## F. SEO/GEO checklist (indexable surfaces)

- [ ] Per-route `<title>` + description; clean canonical without query; one `h1`;
      correct `lang` (R-SEO-01).
- [ ] Demo/example/recipe surfaces are `noindex, nofollow` (R-SEO-01).
- [ ] Unique brand entity Gelium UI; factual citable content; stable deep links;
      visible provenance; JSON-LD where applicable (R-GEO-01).

## G. Repository hygiene

- [ ] Documentation-only slice: no Go/CSS/template/test files changed.
- [ ] No new components, patterns, tokens or wire contracts introduced.
- [ ] `git diff --check` clean (no trailing whitespace, no bad diffs).
- [ ] Cross-links resolve: every relative doc link points to an existing file;
      every rule/schema/contract reference is coherent.

## H. Technical verification commands

```bash
# JSON syntax + schema conformance (see §A)
jq empty docs/gelium/schemas/*.json docs/gelium/schemas/contracts/*.json
python3 -m jsonschema -i docs/gelium/schemas/contracts/admin-resource.json docs/gelium/schemas/screen-recipe.schema.json
python3 -m jsonschema -i docs/gelium/schemas/contracts/inbox-conversation.json docs/gelium/schemas/screen-recipe.schema.json
python3 -m jsonschema -i docs/gelium/schemas/contracts/public-marketing-page.json docs/gelium/schemas/screen-recipe.schema.json
python3 -m jsonschema -i docs/gelium/schemas/contracts/inbox-flow.json docs/gelium/schemas/flow.schema.json

# Documentation slice hygiene
git diff --check

# (Not required for a docs-only slice, but run if any code was touched)
# npm run build && go test ./... && go vet ./...
```

## I. Gates (do not pass without)

1. Schema conformance for every instance (A2).
2. Zero new primitives/contracts (B, G).
3. Complete no-JS branch (E1).
4. Required states declared for every data surface (B, R-STATE-01).
5. Decision-rule references are real and the applicable set is covered (A4, C).
