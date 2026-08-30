# Proposal: agent-gates-and-preflight

## Intent

Turn Gelium UI's distributed guidance into a proportional, evidence-backed workflow. Small understood changes stay direct; a new screen, flow, or substantial redesign moves through explicit product, architecture, approval, build, and audit phases. The workflow blocks an **incomplete phase or evidence claim**, not a person's ability to type locally, and it never claims tooling can prove cognitive reading or own delivery.

The proposal is based on a substantial redesign of the authenticated DeepFilter feed using `gelium-ui@0.6.0`. That trial preserved routes, auth, privacy, pagination, POST contracts, no-JS behavior, and theme classes, while exposing preventable gaps in reading attestations, gate state, detector scope, visual evidence, media metadata, and documentation versioning.

## Problem statement

Gelium documents sound principles, but the current protocol leaves the agent to connect several files and record decisions manually. It lacks a proportional route for small versus architectural work, separates neither product intent from technical feasibility nor prebuild criteria from rendered evidence, and gives no canonical task-local evidence shape.

- `AGENTS.md` points to `llms-ux.txt`, but required reading and the reason each skill applies are not attestable.
- `llms-ux.txt` names workflow gates without a consistent input, output, owner, evidence, status, blocking, and recovery model.
- Skills 11 and 12 correctly discuss criteria and wireframe approval but do not distinguish an intent wireframe from the buildable packet a human approves.
- `ux-detect.sh` runs after implementation and cannot attribute a shared logout form, a bounded custom shell, or an approved exception honestly.
- Rendered contract tests can pass while wide/narrow, light/dark, realistic-content, no-JS, and state evidence remain unavailable.
- `MEDIA-IMAGE` can pressure agents to fabricate dimensions when the consumer model only exposes an external URL.
- Version and wire-contract references can drift. The trial found `gelium-ui@0.6.0` while an LLM-facing reference said `0.5.3`; repository configuration still carries historical `loom:*` wording while current runtime surfaces use `gelium:*` / `X-Gelium-*`.

## Goals

- Route work proportionally as `direct-exempt`, `design-gated`, `escalate`, or optional `full-sdd`.
- Separate **Orient → Plan → Architect → Approve → Build → Audit → Release**.
- Give design-gated work a reusable ledger with clear scope, ownership, evidence, and status semantics.
- Require a buildable, human-approved architecture packet before markup for gated work, without slowing ordinary small fixes.
- Separate prebuild planning evidence from postbuild rendered evidence.
- Keep detector findings, shared ownership, bounded exceptions, unknown media metadata, and version/wire drift visible and honest.
- Preserve the distinction between evidence, human decision, ordinary repository delivery policy, and advisory cross-session memory.

## Non-goals

- Do not require client-side JavaScript for any primary flow.
- Do not replace consumer product decisions with Gelium defaults.
- Do not turn every existing consumer edit into SDD, a ledger, or a wireframe ceremony.
- Do not accept fake avatars, fabricated dimensions, invented content, or decorative markup as compliance evidence.
- Do not rename `loom:*` or `X-Loom-*` contracts automatically; reconcile canonical authority first and make any breaking rename separately.
- Do not redesign DeepFilter or alter its data, routes, auth, privacy, or mutation behavior.

## Proposed changes

### 1. Add proportional routing and the Orient gate

Update `AGENTS.md`, `llms-ux.txt`, `SKILLS.md`, and the applicable skills so the agent first selects one route:

| Route | Use |
|---|---|
| `direct-exempt` | A bounded, understood fix with no architecture change. |
| `design-gated` | A new screen, new flow, or substantial redesign. |
| `escalate` | Product intent, scope, risk, or exception cannot converge safely. |
| `full-sdd` | Durable OpenSpec proposal/spec/design/tasks materially reduce cross-cutting ambiguity. |

`design-gated` work begins with **Orient**: inspect product/design artifacts or record their absence, read the relevant Gelium pack, resolve vocabulary/registry options, and identify hard route/server/no-JS contracts. A reading entry is an attestation with a path and status; tools MUST validate shape and referenced files where possible but MUST NOT claim to prove cognitive reading.

### 2. Add distinct Plan and Architecture gates

**Plan** turns Orient facts and the user request into job, audience, SURFACE/screen, primary action, states, non-goals, and an intent wireframe. It MUST NOT invent components, data, identities, media metadata, or approval.

**Architect** inspects the real route, handler, templates, CSS, data model, permissions, component registry, and no-JS/server contracts. It turns the intent wireframe into the buildable packet: section contracts, primitive choices, responsive/theme/a11y behavior, recovery, reuse, and trade-offs.

A material incompatibility returns the work to Plan or `escalate`; it is not hidden in CSS polish.

### 3. Make approval conditional and machine-readable

For `design-gated` work, the Architecture packet MUST record a human outcome: `approved`, `exempt`, or bounded `exception`. `draft`, `changes-requested`, and `declined` block Build. A chat approval is valid if it records approved scope, approver, date/channel, and packet version.

The packet retains job/audience, route/contracts, ordered sections, primary action, desktop/mobile wireframes, states/recovery, accessibility/no-JS behavior, reuse, and trade-offs. `direct-exempt` work MUST NOT be delayed by a ceremonial wireframe.

### 4. Define a task-local ledger and public progress projection

A future ledger MUST be selected deterministically by one canonical consumer location or an explicit `--ledger <path>` command input. It MUST use a versioned schema and include route, scope, owned/shared paths, reading attestations, prebuild gates, postbuild audit, exception references, and evidence references.

```yaml
schema_version: 1
route: design-gated
scope:
  routes: [/feed]
  owned_paths: [modules/social/feed.templ]
  shared_paths: [modules/shared/logout.templ]
reading:
  - id: llms-ux
    path: node_modules/gelium-ui/llms-ux.txt
    status: attested
gates:
  plan: { status: pass }
  architecture: { status: pass }
  approval: { status: approved, packet: .gelium/wireframes/feed.md }
  criteria_plan: { status: pass }
  rendered_audit: { status: pending }
```

The shown YAML is conceptual. The implementation MUST choose a real versioned parser or a simpler canonical machine format; it MUST NOT parse unrestricted YAML with Bash regex.

Normal user-facing progress is `Working`, `Needs your decision`, `Checking`, or `Ready`. Detailed status remains in the ledger.

### 5. Split criteria planning from rendered audit

Skill 11 MUST ship a prebuild **criteria plan** with hierarchy, DOM order, primary/supporting actions, boundaries, responsive intent, token/theme intent, accessibility/no-JS constraints, states, preserved contracts, and DESIGN-MEMORY reuse.

After Build, the **rendered audit** marks each applicable criterion `pass`, `fail`, `not-applicable`, `pass-with-escalation`, or `exception`, with evidence or an explicit limitation. It covers wide/narrow, light/dark, realistic content, states, keyboard/focus/touch, forced colors/reduced motion, URL/form/server contracts, and no-JS behavior. A prebuild phase MUST NOT require evidence that can exist only after Build.

### 6. Add phased `gelium-preflight`

Add a documented preflight command with text and machine-readable output. Its precise implementation language and parser are deferred until the schema choice is tested.

| Mode | MUST verify | MUST NOT claim |
|---|---|---|
| `prebuild` | required artifacts/attestations, declared route/scope, criteria plan, Architecture packet, and approval/exemption. | That an agent cannot type locally or that a human cognitively inspected a file. |
| `release` | rendered-audit references, detector output, bounded exceptions, authority-matrix coherence, and applicable test/build evidence. | Authority to commit, push, publish, or deploy. |

Changed paths MUST be checked against the ledger's owned/shared paths. They are evidence of coverage, not a way to infer whether a redesign is substantial.

### 7. Extend `ux-detect.sh` with scoped, honest result semantics

Future detector options MUST support an explicit audited scope, an exception manifest, and machine-readable output while preserving existing default behavior during migration. Every raw finding remains visible.

An exception MUST identify the rule or finding fingerprint, path/selector scope, reason, bounded risk, owner, ledger evidence, and deterministic expiry (`expires_at` or `expires_before_version`). A broad path exclusion is not evidence that a rule passed.

| Result | Meaning |
|---|---|
| `clean-pass` | No unapproved findings remain. |
| `pass-with-exceptions` | Remaining findings have bounded approved exceptions; raw findings and IDs remain visible. |
| `failed` | A required check or finding remains unresolved. |
| `invalid-configuration` | Ledger or exception manifest is malformed, stale, or semantically invalid. |

Shared layout findings may be attributed outside the owned screen surface only when the ledger declares the boundary. They remain in audit and machine-readable output.

### 8. Clarify media contract for unknown dimensions

Update `MEDIA-IMAGE` and media guidance to distinguish known intrinsic dimensions, unknown dimensions, and non-applicable decorative/omitted media. Unknown informative media MUST retain meaningful alt, safe responsive containment, recovery, and an explicit `media-metadata-unknown` or `pass-with-escalation` outcome. It MUST NOT fabricate `width` or `height` values.

### 9. Add authority matrices before coherence checks

Before version and wire-contract checks ship, maintainers MUST define equivalent authority surfaces for monorepo version, published package version, asset version, release-facing documentation, current runtime contracts, and historical/migration references. A checker reports each conflicting path/value and MUST NOT automatically rewrite a wire prefix or release version.

### 10. Define subagent, Engram, and SDD boundaries

For broad discovery, the workflow MAY use at most two read-only explorers: one product/UX explorer and one technical/Gelium explorer. They receive exact paths and return bounded facts, risks, and open questions. One writer owns production changes. A fresh reviewer MAY audit the final packet, diff, and evidence.

Project memory MAY seed Orient with prior decisions, component-fit mismatches, approved exceptions, or lessons. Current repository artifacts always override memory. Memory is never reading proof, approval, audit evidence, task state, or delivery authority.

OpenSpec/SDD is optional for ordinary consumer work and required only when durable proposal/spec/design/tasks materially reduce cross-cutting ambiguity. A screen ledger is operational evidence, not an SDD review transaction.

## Evidence from the DeepFilter trial

The proposal is grounded in these observed outcomes from a real `Read` / `list` / `DATA-FEED` screen:

1. The consumer initially followed the `AGENTS.md` pointer but did not directly read `llms-ux.txt` before implementation. A visible boot stop would have prevented that process gap.
2. The substantial redesign was implemented before the wireframe and criteria records were formalized. A canonical status file would have made the missing gate obvious.
3. The final feed preserved the route, anonymous `303` redirect, private/no-store response, visibility filtering, pagination, likes, saves, CSRF, POST+303, 422, localized chrome, class-routed themes, and no-JS navigation.
4. Rendered contract tests and full package/build checks passed, but live visual viewport evidence was initially unavailable. The workflow needs to report this as a limitation, not imply a complete visual audit.
5. The scoped detector correctly found a shared logout form but could not express ownership and bounded product rationale cleanly. The finding became an explicit attribution rather than being hidden.
6. Feed thumbnails had meaningful alt text and lazy loading, but the consumer model exposed only external URLs. No dimensions were invented; this exposed a real contract gap in `MEDIA-IMAGE`.
7. The package used `gelium-ui@0.6.0` while an LLM-facing version reference still said `0.5.3`, demonstrating that documentation drift is a release-quality problem.

## Acceptance criteria

### Scenario 1 — direct work remains direct

**Given** an already-understood accessibility, copy, token, selector, or contract correction that does not change page architecture
**When** the agent selects `direct-exempt`
**Then** it performs applicable checks without requiring an Architecture wireframe approval.

### Scenario 2 — gated work separates intent from viability

**Given** a task is classified `design-gated`
**When** Orient, Plan, and Architect run
**Then** Plan records the user/product intent and intent wireframe, while Architect records the route/data/component/no-JS validation and buildable wireframe.

### Scenario 3 — a missing human decision blocks Build, not local cognition

**Given** a design-gated Architecture packet has status `draft`, `changes-requested`, or `declined`
**When** `gelium-preflight prebuild` runs
**Then** it reports the packet path, status, and required decision and refuses a passing prebuild result; it does not claim to prevent local file edits.

### Scenario 4 — prebuild does not require postbuild evidence

**Given** a design-gated packet has an approved architecture and criteria plan but no candidate implementation
**When** `gelium-preflight prebuild` runs
**Then** it validates the prebuild artifacts and MUST NOT demand rendered viewport, theme, state, or no-JS evidence.

### Scenario 5 — release requires rendered evidence honestly

**Given** a candidate is ready for release audit
**When** `gelium-preflight release` runs
**Then** it requires rendered audit references, detector outcome, applicable tests/build, and authority-matrix checks; an unobserved viewport/theme/no-JS condition remains a limitation rather than a clean pass.

### Scenario 6 — exceptions and shared findings remain visible

**Given** a scoped detector finding has a valid bounded exception or a declared shared ownership boundary
**When** detector output is produced
**Then** it preserves the raw finding, owner/scope attribution, and exception ID; the outcome is `pass-with-exceptions`, not `clean-pass`.

### Scenario 7 — unknown media dimensions are not fabricated

**Given** an informative external image has meaningful alt text but no trustworthy dimensions
**When** the consumer renders it
**Then** the audit records `media-metadata-unknown` or `pass-with-escalation`, preserves responsive/recovery obligations, and does not generate fake `width` or `height` values.

### Scenario 8 — memory is advisory and repository artifacts win

**Given** project memory reports a previous component or contract decision
**When** current repository artifacts disagree
**Then** the agent follows the current artifacts and records the discrepancy when it affects the task.

### Scenario 9 — drift checks compare only authority equivalents

**Given** maintainers define a version or wire-contract authority matrix
**When** a coherence check finds conflicting equivalent surfaces
**Then** it reports each path and value without rewriting a contract; historical references outside the matrix do not create false drift.

## Implementation sequence

1. Add Route, Orient, Plan, Architect, approval, prebuild/rendered-audit separation, public progress states, subagent boundaries, and memory/SDD boundaries to the guidance and skills.
2. Publish a package-included versioned ledger/packet template and tests that assert required guidance surfaces.
3. Select and test a real ledger parser or simpler canonical format; add fixture validation before implementing a CLI.
4. Implement `gelium-preflight prebuild` with text and machine-readable output, scope coverage, and no false claim of cognitive enforcement.
5. Implement `gelium-preflight release` with audit evidence and authority-matrix checks; keep commit/push/deploy under ordinary repository policy.
6. Add detector scope, exception, shared ownership, expiry, and result semantics while retaining default compatibility.
7. Clarify media unknown-dimension guidance and detector reporting.
8. Define authority matrices, then add version/wire coherence fixtures and checks.
9. Dogfood advisory mode on DeepFilter or an equivalent feed; use findings to decide the required-mode rollout.
10. Publish migration/release documentation.

## Rollback plan

The workflow implementation is additive: guidance, templates, scripts, tests, and advisory/required policy wiring can be reverted independently. Revert the relevant unit and remove its CI invocation without removing existing consumer routes, components, or wire contracts. Preserve historical ledger/audit artifacts as evidence but do not let them control ordinary delivery.

## Risks and mitigations

| Risk | Likelihood | Mitigation |
|---|---:|---|
| Preflight becomes ceremonial | Medium | Treat receipts as attestations, validate shape/path/evidence, and retain human decisions explicitly. |
| Ceremony blocks small work | Medium | Route `direct-exempt` work around design-gated artifacts; use `full-sdd` only for material ambiguity. |
| Context/token cost rises through over-delegation | Medium | At most two exact-path read-only explorers, one writer, bounded output, and a fresh reviewer only when scope warrants it. |
| Existing consumers need legitimate exceptions | High | Support bounded scope, owner, evidence, and deterministic expiry; distinguish exception results from clean pass. |
| Parser or CLI complexity grows prematurely | Medium | Freeze semantics first; choose a real parser/format with fixture tests before scripting. |
| Media guidance weakens accessibility | Medium | Keep meaningful alt, responsive containment, recovery, and escalation mandatory when dimensions are unknown. |
| Authority checks create permanent false failures | Medium | Define equivalence matrices before enabling drift checks; retain historical references outside current authority sets. |

## Success criteria

- A reviewer can identify the selected route, current phase, artifact path, scope ownership, decision owner, evidence, and unresolved limitation for a design-gated surface.
- Small, understood changes remain direct and are not blocked by a wireframe ritual.
- Prebuild validates only prebuild facts; rendered audit/release validates postbuild evidence without claiming unseen coverage.
- Detector output distinguishes clean pass, approved visible exceptions, shared attribution, failure, and invalid configuration.
- Unknown image dimensions are represented honestly without fake compliance.
- Version/wire drift checks compare authority equivalents and never rename contracts automatically.
- `go test ./internal/... ./site/... ./lib/...`, `go vet ./internal/... ./site/... ./lib/...`, `npm run build`, and `git diff --check` remain the required future production implementation verification gates.
