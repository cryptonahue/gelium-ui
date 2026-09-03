# Design: agent-gates-and-preflight

## Technical approach

The change defines a proportional operating protocol and its machine-readable
artifacts. The ledger, detector, and preflight contracts are implemented in
`internal/gates`, `cmd/gelium-preflight`, and `cmd/gelium-ux-detect`; this design
records their boundaries and does not grant delivery authority.

```text
ROUTE → ORIENT → PLAN → ARCHITECT → APPROVE → BUILD → AUDIT → RELEASE
```

The normal public projection is deliberately compact:

```text
Working | Needs your decision | Checking | Ready
```

Detailed statuses remain in artifacts so a user does not need to reason about parser state, receipt IDs, or internal lifecycle details.

## Phase contracts

| Phase | Inputs | Output | Block / recovery |
|---|---|---|---|
| Route | Request, changed surface, known risk | `direct-exempt`, `delegated-direct`, `design-gated`, `escalate`, or `full-sdd` | Ambiguity → `escalate`. |
| Orient | Product/design artifacts, Gelium pack, vocabulary/registry, prior memory, hard contracts | Constraint map and reading attestation | Missing product intent → human record or artifact exception. |
| Plan | Constraint map, user request, product intent | Brief, states, primary action, Plan wireframe | Unresolved product choice → `Needs your decision`. |
| Architect | Plan wireframe, route/handler/template/data/component inspection | Buildable wireframe, section contracts, component/contract mapping | Component/data mismatch → revise or exception. |
| Approve | Architecture packet | `approved`, `exempt`, or bounded `exception` | `draft`, `changes-requested`, `declined`, and missing decision block Build. |
| Build | Approved packet | Candidate implementation | Material scope change → Route again. |
| Audit | Candidate, rendered output, detector, tests | Rendered evidence and findings | Failure/unbounded mismatch → Build or escalation. |
| Release | Audit evidence, authority matrices, relevant checks | `clean-pass`, `pass-with-exceptions`, or failure | Failure → repair/re-audit; delivery remains repo-policy-owned. |

## Wireframe representations

1. **Plan wireframe:** expresses user job, information hierarchy, action policy, states, recovery, and intended section order without pretending implementation details.
2. **Architecture wireframe:** maps that intent to actual data, contracts, registered primitives, responsive behavior, themes, accessibility, and no-JS fallback. This is the approval packet.

## Ledger model

The JSON v1 ledger is selected deterministically, either by one canonical
consumer path or explicit `--ledger <path>` input. It includes a
`schema_version`; a scope declaration; route classification; owned and shared
paths; reading attestations; prebuild gates; postbuild evidence; and exception
references.

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

The shown YAML is a conceptual shape only. The implemented canonical machine
format is versioned JSON, validated by `internal/gates`; Bash does not parse
unrestricted YAML.

## Preflight boundary

The implemented preflight has three modes:

| Mode | Checks | Does not claim |
|---|---|---|
| `route` | Validate the selected route and return its next action and design-gate requirement. | Infer intent from file count or authorize implementation/delivery. |
| `prebuild` | Artifact presence, attestations, declared scope coverage, plan/architecture/approval state. | That no local markup can be typed or that a human cognitively reviewed a file. |
| `release` | Rendered evidence references, detector results, exceptions, authority-matrix coherence, tests/build results. | Authority to commit, push, publish, or deploy. |

Changed paths validate declared scope coverage. They MUST NOT be used to infer whether a conceptual redesign is substantial.

## Detector and exception design

The detector retains raw findings. A matching exception narrows a finding only when it records finding/rule identity, path/selector scope, reason, risk, owner, ledger evidence, and a deterministic expiry (`expires_at` or `expires_before_version`).

| Result | Meaning |
|---|---|
| `clean-pass` | No unapproved findings remain. |
| `pass-with-exceptions` | Every remaining finding has a bounded approved exception; raw findings remain visible. |
| `failed` | An unapproved finding or failed required check remains. |
| `invalid-configuration` | Ledger/manifest is malformed, stale, or semantically invalid. |

## Subagents, memory, and SDD

- Exact-path, bounded-output explorers protect the parent context; they do not own approval or write production code.
- One writer owns the implementation.
- A fresh reviewer is useful after Build because it examines a candidate and avoids author confirmation bias.
- Memory is advisory cross-session context; repository files and current artifacts override it.
- OpenSpec/SDD owns substantial cross-cutting planning. Screen ledgers are operational evidence, not an SDD transaction system.

## Authority matrices

Before drift checks ship, maintainers must define an authoritative matrix for:

1. monorepo version versus published npm package version versus asset version;
2. release-facing documentation versus historical documentation;
3. canonical server headers/events versus historical/migration references.

A checker compares equivalent authority surfaces only. It reports drift and never renames a wire contract automatically.

