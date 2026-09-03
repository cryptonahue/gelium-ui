# Tasks: agent-gates-and-preflight

## Delivery strategy

This change began as documentation/specification work. The ledger, preflight,
detector, authority, and dogfood units below are now reconciled against the
tracked implementation; remaining unchecked items are the only open work.

| Unit | Goal | Primary verification | Rollback boundary |
|---|---|---|---|
| 1 | Guidance, skills, templates, and public phase language | focused documentation/spec review | Guidance artifacts only |
| 2 | Ledger schema/parser and fixture validation | parser/validator unit tests | New ledger code/templates |
| 3 | `gelium-preflight` prebuild/release modes | CLI/fixture tests | New preflight command |
| 4 | Scoped detector and exception result semantics | detector fixtures + JSON/text output tests | New detector flags/manifest support |
| 5 | Version/wire authority matrix checks | coherence tests with deliberate fixtures | New consistency checks |
| 6 | Media unknown-dimension guidance and detector reporting | media guidance/test contracts | Guidance/detector rule only |
| 7 | Dogfood and advisory rollout | DeepFilter/equivalent consumer record | Consumer fixture and rollout docs |

## Phase 0: OpenSpec revision

- [x] 0.1 Rewrite `proposal.md` around proportional routing and phased gates.
- [x] 0.2 Add `exploration.md` with repository baseline and bounded external lessons.
- [x] 0.3 Add `design.md` with phase, ledger, authority, subagent, and memory boundaries.
- [x] 0.4 Add delta specs for workflow, ledger/preflight, and detector/evidence.
- [x] 0.5 Update `agent-gates-flow.html` for Route → Orient → Plan → Architect → Approve → Build → Audit → Release.
- [x] 0.6 Review terminology and run whitespace validation.

## Phase 1: Guidance and artifacts

- [x] 1.1 RED: Add tests that assert the published agent guidance defines Route, Orient, Plan, Architect, and the prebuild/rendered-audit split.
- [x] 1.2 GREEN: Update `lib/AGENTS.md`, `lib/llms-ux.txt`, `lib/SKILLS.md`, and applicable skills with the proportional workflow.
- [x] 1.3 GREEN: Ship a ledger/packet template from a package-included location, not an unshipped docs-only path.
- [x] 1.4 GREEN: Add public progress names and clarify they are projections, not detailed machine state.
- [x] 1.5 Verify published-package file coverage and focused tests.

## Phase 2: Ledger validation and preflight

- [x] 2.1 RED: Add fixture tests for valid and malformed schema v1 ledgers, missing attestations, unknown statuses, invalid ownership boundaries, and expired exceptions.
- [x] 2.2 GREEN: Implement one versioned parser/validator; do not parse arbitrary YAML with shell regex.
- [x] 2.3 RED: Add prebuild fixtures for direct-exempt, design-gated approved, draft/declined, user-record, and scope mismatch outcomes.
- [x] 2.4 GREEN: Implement `gelium-preflight` prebuild mode with text and machine-readable output.
- [x] 2.5 RED: Add release fixtures requiring audit evidence and authority-matrix results.
- [x] 2.6 GREEN: Implement `release` mode; it reports evidence only and does not control commit/push/deploy.

## Phase 3: Detector, media, and coherence

- [x] 3.1 RED: Add scoped detector fixtures covering an owned finding, shared finding, approved exception, unapproved exception, and malformed manifest.
- [x] 3.2 GREEN: Add scope/exception/result semantics while preserving the existing default detector behavior.
- [x] 3.3 RED: Add media fixtures for known dimensions, unknown dimensions, and decorative/omitted media.
- [x] 3.4 GREEN: Update `MEDIA-IMAGE` and detector result vocabulary without fabricating dimensions.
- [x] 3.5 RED: Define authority-matrix fixtures for package/version and wire-contract drift.
- [x] 3.6 GREEN: Implement scoped coherence checks that report values/paths and never rewrite contracts.

## Phase 4: Dogfood and rollout

- [x] 4.1 Apply the advisory workflow to DeepFilter or an equivalent consumer feed.
- [x] 4.2 Record a buildable packet, rendered audit limitations, scope ownership, and any bounded exceptions.
- [x] 4.3 Review false positives and adoption friction before enabling required mode for future design-gated scopes.
- [x] 4.4 Publish migration guidance; preserve direct-exempt work and legacy exception paths.
- [x] 4.5 Enable required visible-wireframe approval for design-gated work; `continua` without a shown packet is not approval.

## Implementation verification baseline

```text
go test ./internal/... ./site/... ./lib/...
go vet ./internal/... ./site/... ./lib/...
npm run build
git diff --check
```

Run focused parser/CLI/detector fixture tests before the full suite. Treat any rendered-audit limitation or approved exception as visible output, never as a clean pass.
