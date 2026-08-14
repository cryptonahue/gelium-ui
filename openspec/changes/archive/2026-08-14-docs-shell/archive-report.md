# Archive Report: docs-shell

**Change**: docs-shell  
**Product**: Gelium UI · module `loomui`  
**Archived**: 2026-08-14  
**Mode**: hybrid (Engram + OpenSpec)  
**Branch**: `slice/docs-shell-pr1-nav-model`  
**Status**: archived  

## Gates

| Gate | Result | Notes |
|------|--------|-------|
| Task completion | PASS | 20/20 `[x]` in tasks.md; 0 unchecked |
| Verify verdict | PASS | 10/10 requirements, 19/19 scenarios, 0 CRITICAL |
| Native verify admission | PASS after envelope fix | Removed unknown `vet_command` / `vet_exit_code` / `vet_output_hash` from leading `gentle-ai.verify-result/v1` YAML (vet evidence retained in prose body) |
| Review receipt | N/A | receipt-driven development **off** (global); native status `reviewGate` null; archive dependency `ready` |
| Action context | PASS | `repo-local`; allowed root = loom-ui workspace |
| Destructive delta | None | New domain spec only (ADDED full capability) |

## Specs synced

| Domain | Action | Details |
|--------|--------|---------|
| docs-shell | **Created** | Main spec did not exist; copied delta as full source of truth → `openspec/specs/docs-shell/spec.md` (10 requirements, 19 scenarios) |

No MODIFIED / REMOVED / RENAMED requirements. No destructive merge.

## Archive location

```
openspec/changes/docs-shell/
  → openspec/changes/archive/2026-08-14-docs-shell/
```

### Archive contents

| Artifact | Present |
|----------|---------|
| exploration.md | ✅ |
| proposal.md | ✅ |
| design.md | ✅ |
| tasks.md | ✅ (20/20 complete) |
| specs/docs-shell/spec.md | ✅ |
| verify-report.md | ✅ (PASS) |
| archive-report.md | ✅ (this file) |

Active path `openspec/changes/docs-shell/` removed after move.

## Implementation commits (audit trail; not merged by archive)

| Hash | Slice | Message |
|------|-------|---------|
| `f6ac0d5` | PR1 | feat(docs-shell): add docs nav model and derive navLinks |
| `f85b707` | PR1 | feat(docs-shell): rebuild defaultFooter from docsNavFor |
| `97cb5d4` | PR2 | feat(docs-shell): add two-pane shell chrome, stubs, and CSS |
| `0224886` | PR3 | test(docs-shell): tighten chrome contracts for active path and layout |
| `ecfed36` | PR3 | docs(docs-shell): close nav discoverability residual |

## Engram observation IDs (traceability)

| Artifact | Topic key | Observation ID |
|----------|-----------|----------------|
| explore | `sdd/docs-shell/explore` | 938 |
| proposal | `sdd/docs-shell/proposal` | 939 |
| spec | `sdd/docs-shell/spec` | 950 |
| design | `sdd/docs-shell/design` | 951 |
| tasks | `sdd/docs-shell/tasks` | 952 |
| apply-progress | `sdd/docs-shell/apply-progress` | 954 |
| verify-report | `sdd/docs-shell/verify-report` | 960 |
| archive-report | `sdd/docs-shell/archive-report` | (this save) |

Review topics `sdd/docs-shell/review/*`: none (review mode off).

## Envelope remediation (pre-archive)

Native status initially blocked archive with:

> unknown verify result field `vet_command`

**Action**: stripped non-schema fields `vet_command`, `vet_exit_code`, `vet_output_hash` from the leading verify YAML only. Prose sections still document `go vet ./...` PASS.  
**Validation**: `gentle-ai sdd-verify-validate` → `valid: true`, `verdict: pass`, same `evidence_revision`.  
**Post-fix status**: `nextRecommended: archive`, `dependencies.archive: ready`, `blockedReasons: []`.

Intentional note: envelope field cleanup only — no implementation code change, no 4R review started, no merge/PR.

## Verify summary (from id 960 / archived verify-report)

- Verdict: **PASS**
- Tasks: 20/20  
- Requirements: 10/10 · Scenarios: 19/19  
- Tests: `go test ./... -count=1` exit 0  
- Build: `go build -o /tmp/gelium ./cmd_smoke_main.go` exit 0  
- CRITICAL: none  
- WARNINGs retained: cmd/gelium vs smoke main; home theme switcher optional polish  

## Source of truth updated

- `openspec/specs/docs-shell/spec.md` — two-pane docs shell capability (frame, IA, active route, topbar, theme 0-JS, mobile details, dogfood, landmarks, no-JS, OOS)

## SDD cycle

Planned → implemented (PR1–PR3 on stack branch) → verified PASS → **archived**.

## Human remaining (out of archive scope)

Archive does **not** merge code or open PRs. Remaining for the human/orchestrator:

1. Land stacked PRs for branch `slice/docs-shell-pr1-nav-model` (commits above) toward main — delivery strategy stacked-to-main / auto-chain.  
2. Do **not** treat archive as a substitute for PR review or CI green on the feature branch.  
3. Optional polish from verify WARNINGs (home ThemeSwitcher) is a separate change if desired.  
4. Next SDD change: choose next roadmap residual / product slice (`sdd-new`).

## Rules honored

- Did not delete implementation code  
- Did not merge to main or open PRs  
- Did not start 4R review  
- Did not modify archived artifacts after move except writing this report into the archive folder  
