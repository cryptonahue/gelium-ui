# gelium-ui — Skills for agents

Actionable decision skills. Read them in order for a task; load the ones relevant
to a specific job. The full decision-id pack (SURFACE / SCREEN / WF / DATA / FEED
/ JOURNEY / MEDIA / SKEL) is in `llms-ux.txt`.

| # | Skill | What it answers |
|---|---|---|
| 01 | `01-foundations.md` | tokens, themes by class, 0-JS, layout, mobile guardrails |
| 02 | `02-surface-and-screens.md` | SURFACE mode + screen types → layout blocks |
| 03 | `03-forms-and-controls.md` | which control, form layout, 422 validation |
| 04 | `04-state-and-feedback.md` | states + FEED channel decision matrix |
| 05 | `05-server-contracts.md` | GET / POST+303 / 422 / gelium:toast wire |
| 06 | `06-mobile-and-a11y.md` | touch, keyboard, semantics, contrast, media |
| 07 | `07-dod-and-antislop.md` | Definition of Done checklist + anti-slop |
| 08 | `08-product-reasoning.md` | discovery workflow: find missing product UX before drawing |

## Install these skills into your agent tool

Copy this folder into your agent's skill directory so the LLM loads it in any
project that uses `gelium-ui`:

```bash
bash install-agents.sh          # copies skills/ + AGENTS.md + llms-ux.txt
```

Targets detected: `~/.hermes/skills/gelium-ui`, `~/.cursor/skills/gelium-ui`,
`~/.claude/skills/gelium-ui`. See `install-agents.sh --help`.
