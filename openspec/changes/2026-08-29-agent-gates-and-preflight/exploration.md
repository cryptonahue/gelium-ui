# Exploration: agent-gates-and-preflight

## External workflow reference

The proportional routing, compact public progress states, bounded phase artifacts, fresh-candidate review, and authority limits in this exploration were evaluated against [Gentle-AI](https://github.com/Gentleman-Programming/gentle-ai), especially its [trigger rules](https://raw.githubusercontent.com/Gentleman-Programming/gentle-ai/main/docs/trigger-rules.md), [organic RDD architecture](https://raw.githubusercontent.com/Gentleman-Programming/gentle-ai/main/docs/architecture/organic-rdd.md), and [review authority threat model](https://raw.githubusercontent.com/Gentleman-Programming/gentle-ai/main/docs/review-authority-threat-model.md). They are adapted, not imported as a runtime dependency.

## Outcome

This change turns Gelium UI guidance into a proportional workflow for agents. It does not make every edit a design review and does not claim that local tooling can authenticate an agent's cognition or control ordinary repository delivery.

## Repository baseline

| Surface | Observed baseline | Consequence |
|---|---|---|
| `lib/AGENTS.md` | Points agents to the decision pack and skills, but has no ordered receipt artifact. | Add an explicit Orient route and attestable reading evidence. |
| `lib/llms-ux.txt` | Defines product and workflow guidance, but combines pre-emit judgment with audit guidance. | Split criteria planning from rendered evidence. |
| `lib/skills/11-design-criteria.md` | Requires rendered checks across themes, widths, realistic content, and states. | Keep that evidence postbuild; do not require it before markup exists. |
| `lib/skills/12-wireframe-approval.md` | Correctly requires approval before a new screen, flow, or substantial redesign. | Add a machine-readable packet and phase boundary without applying it to exemptions. |
| `lib/scripts/ux-detect.sh` | Legacy positional detector remains available; `cmd/gelium-ux-detect` adds scoped ownership, exceptions, and machine-readable result semantics. | Keep both paths during migration and preserve raw findings. |
| `lib/llms.txt` / `lib/package.json` / `openspec/config.yaml` | Version and wire vocabulary disagree across release-facing and historical/configuration surfaces. | Define authority matrices before a drift detector attempts to compare files. |

## DeepFilter findings retained by this change

- Direct reading of `llms-ux.txt` was missed before the redesign.
- Wireframe and criteria records were formalized after implementation instead of before it.
- A rendered contract test/build result did not prove wide/narrow, light/dark, realistic-content, and no-JS coverage.
- Shared logout markup needed ownership-aware, visible attribution rather than silent exclusion.
- External image URLs did not contain trustworthy intrinsic dimensions; fake dimensions were correctly avoided.

## Routing model

| Route | Trigger | Required record | Human decision |
|---|---|---|---|
| `direct-exempt` | Small, already-understood copy, token, selector, accessibility, bug, or contract correction with no architecture shift. | Narrow change note only when required by existing repo policy. | No wireframe ceremony. |
| `delegated-direct` | Broad context, read-only research, or multi-file work without a screen/flow architecture shift. | Bounded worker handoff plus parent verification. | No design gate. |
| `design-gated` | New screen, new flow, or substantial redesign. | Ledger + Plan + Architecture packet. | Required before Build. |
| `escalate` | Missing product decision, unbounded exception, unknown risk, or ambiguous scope. | Escalation record. | Required to resume. |
| `full-sdd` | Cross-cutting work where proposal/spec/design/tasks materially reduce ambiguity. | OpenSpec artifacts plus the relevant Gelium records. | Proposal/design acceptance as defined by repo policy. |

`full-sdd` is not a risk score and is not automatic for every UI task. This proposal itself uses it because it spans normative guidance, skills, future scripts, detector behavior, tests, release coherence, and rollout.

## External workflow lessons: adopt, adapt, reject

| Lesson | Decision for Gelium |
|---|---|
| Proportional routing and small work staying direct | Adopt. Use `direct-exempt` rather than forcing a ledger/wireframe for every edit. |
| Compact public progress states | Adopt. Report `Working`, `Needs your decision`, `Checking`, or `Ready`; retain detailed machine states in artifacts. |
| Durable artifacts between phases | Adapt. Use OpenSpec for cross-cutting work and a task-local ledger/packet for consumer surfaces. |
| Fresh-context review after work exists | Adopt. Use a focused auditor for broad work; audit a real candidate, not a narrative. |
| Receipts have bounded authority | Adopt. A ledger attests structured facts and evidence references; it cannot prove cognitive reading or authorize delivery. |
| Immutable review lineages, locks, and delivery machinery | Reject. They solve a broader review-transaction problem and would overcomplicate Gelium UI guidance. |

## Memory boundary

Project memory may be queried during Orient for prior decisions, component-fit findings, approved exceptions, and durable learnings. The agent MUST validate relevant memory against current repository artifacts before relying on it. Memory is never a substitute for required reading, a wireframe approval, audit evidence, or the source of truth for current code.

## Subagent boundary

Use at most two read-only explorers when the task genuinely requires broad context: one for product/UX context and one for technical/Gelium constraints. Pass exact file paths and request bounded findings. One writer owns production edits. A fresh reviewer may inspect the final packet, diff, and evidence. Parallel writers are out of scope for this workflow.
