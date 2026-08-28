# gelium-ui — Skills for agents

Actionable decision skills. Read them in order for a task; load the ones relevant
to a specific job. The full decision-id pack (SURFACE / SCREEN / ARCH / SECTION / WF / DATA /
FEED / JOURNEY / MEDIA / SKEL) is in `llms-ux.txt`.

| # | Skill | What it answers |
|---|---|---|
| 01 | `01-foundations.md` | tokens, themes by class, 0-JS, layout, mobile guardrails |
| 02 | `02-surface-and-screens.md` | SURFACE mode + screen types → layout blocks |
| 03 | `03-forms-and-controls.md` | which control, form layout, 422 validation |
| 04 | `04-state-and-feedback.md` | states + FEED channel decision matrix |
| 05 | `05-server-contracts.md` | GET / POST+303 / 422 / gelium:toast wire |
| 06 | `06-mobile-and-a11y.md` | touch, keyboard, semantics, contrast, media |
| 07 | `07-dod-and-antislop.md` | Definition of Done checklist + anti-slop |
| 08 | `08-product-reasoning.md` | discovery workflow: find missing product-level UX before drawing |
| 09 | `09-usability-checklist.md` | per-screen binary usability checklist |
| 10 | `10-page-section-architecture.md` | purpose-first page and section contracts before components |
| 11 | `11-design-criteria.md` | checkable visual judgment, reuse, and resilient polish |
| 12 | `12-wireframe-approval.md` | conditional approval gate for new screens, flows, and substantial redesigns |

## Required architecture handoff

Keep the existing gates in order: foundations → surface → product reasoning →
page/section → design criteria §1–§3 → conditional wireframe approval for new
screens/flows/substantial redesigns → registered components → tokens/skin →
design criteria critique → usability → DoD. If `PRODUCT.md` or `DESIGN.md` is
absent in a consumer repo, stop and ask before this sequence.

Protocol IDs: `ARCH-PRODUCT`, `ARCH-PAGE`, `ARCH-SECTION`, `ARCH-COMPONENTS`,
`ARCH-TOKENS`, `SECTION-CONTRACT`, `SECTION-HIERARCHY`, `SECTION-ACTION`,
`SECTION-REVELATION`, `SECTION-RECOVERY`, `WF-ARCH`, `WF-SECTION-AUDIT`.

## Description-to-name resolver

**Before creating any new component**, resolve informal UI descriptions ("the
dark see-through layer behind a popup", "the eye icon on the password field")
against `lib/ui-vocabulary.md` — an attributed glossary mapping ~76 UI elements
to canonical names, common aliases, and API symbols. Then check
`docs/gelium-ui-vocabulary.md` §8 (resolved naming conflicts) and compose a
registered component instead of inventing one. See the "Vocabulary resolution"
step in `skills/02-surface-and-screens.md`.

Vocabulary source: namethatui.com (used with attribution; removal on request).

## Install these skills into your agent tool

Copy this folder into your agent's skill directory so the LLM loads it in any
project that uses `gelium-ui`:

```bash
bash install-agents.sh          # copies skills/ + AGENTS.md + llms-ux.txt
```

Targets detected: `~/.hermes/skills/gelium-ui`, `~/.cursor/skills/gelium-ui`,
`~/.claude/skills/gelium-ui`. See `install-agents.sh --help`.
