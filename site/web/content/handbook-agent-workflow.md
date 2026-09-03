# Agent workflow

How humans and LLMs **work with Gelium** without leaving the product ethos: server HTML, token themes, sourced criteria (GOV.UK / USWDS / M3 / NNG), progressive JS only.

This page adapts useful **agent process** ideas (brief → shape → build → audit → polish) — not free-form “make it award-winning” aesthetics that fight Material/Basecoat tokens.

## Start here: route the outcome

The package source of truth is [`skills/00-agent-routing.md`](https://github.com/cryptonahue/gelium-ui/blob/main/lib/skills/00-agent-routing.md).
It chooses `direct-exempt`, `delegated-direct`, `design-gated`, `escalate`, or
`full-sdd` before this downstream UX workflow is loaded. Delegation is per
action, and delivery remains under ordinary repository policy.

## Ethos (non-negotiable)

| Do | Don’t |
|---|---|
| Prefer `ui-*` partials + `--ui-*` tokens | Invent a parallel CSS design system per screen |
| Map `FEED-*` / `DATA-*` / `JOURNEY-*` | Toast-only validation; blank tables |
| One theme class on `<html>` | Random hex, purple gradients “for polish” |
| 0-JS first; enhance | Require a SPA runtime for basic CRUD |
| Cite handbook + sources | Override Forms/Feedback with vibe |

## Surface modes

Choose **one mode per URL/surface** (not per product company). Same product can have a Persuade landing and an Operate admin.

| Mode | Visitor success | Bias | Gelium defaults |
|---|---|---|---|
| **Operate** | Finish a task | Scan, density, consistency, native controls | Admin recipes, data-table, forms, FEED-VAL, compact/cozy |
| **Read** | Understand | Measure 65ch, clear H1, in-page nav | Handbook-style docs, long prose |
| **Persuade** | Decide and act | One primary CTA, honest hierarchy | Hub/start, hero, short proof — still tokens, not stock SaaS gradients |
| **Experience** | Be inside the artifact | Content leads; chrome recedes | Galleries/media; don’t fake it on dense CRUD |

**Operate wins ties** for tools, settings, queues, and tables.

## Workflow passes (bounded)

Do **not** open-ended self-QA loops. Fixed passes:

| Pass | ID | Goal | Exit when |
|---|---|---|---|
| **1. Brief** | `WF-BRIEF` | Audience, job, constraints; optional PRODUCT.md / DESIGN.md | Job sentence + surface mode written |
| **2. Shape** | `WF-SHAPE` | Screen type, journey, data pattern, FEED plan — **no markup yet** | IDs listed (SCREEN / JOURNEY / DATA / FEED) |
| **3. Architecture** | `WF-ARCH` | Major-region purpose, hierarchy, action, revelation, and recovery — **no components yet** | Every major region has a `SECTION-CONTRACT` ([Page + section architecture](/docs/page-section-architecture)) |
| **4. Build** | `WF-BUILD` | Partials, theme, server contracts | Renders; happy path works |
| **5. Audit** | `WF-AUDIT` | Technical + contract checks, including finite `WF-SECTION-AUDIT` (a11y basics, responsive, FEED/DoD) | Detector script + checklist clean or waived with reason |
| **6. Polish** | `WF-POLISH` | Align to theme, spacing tokens, copy (Content style) | One pass only — then stop |
| **Harden** (as needed) | `WF-HARDEN` | Empty/error/i18n overflow, edge cases | Edge paths covered |
| **Onboard** (as needed) | `WF-ONBOARD` | First-run + empty activation | Empty has CTA |

`shape` owns discovery; `architecture` contracts major regions; `build` owns code; `audit` is evidence; `polish` is alignment — not a redesign.

## Anti-slop (Gelium-aware)

Adapted from common AI-frontend failure modes; **scoped so they don’t ban the design system itself**.

| Avoid | Why | Prefer |
|---|---|---|
| Nested cards inside cards for simple forms | Noise, false hierarchy | One surface; `divider` / sections |
| Stock purple→blue hero gradients | Generic AI marketing look | Theme primary/surface tokens |
| Icon-in-rounded-square above every heading | Cliché template | Icon only when it carries meaning |
| Gray text on saturated fills | Contrast failure | Theme fg/on-color tokens; check a11y |
| Bounce/elastic motion everywhere | Dated; motion policy | Short token motion or none ([Motion](/docs/motion)) |
| “Make it bolder” on Operate tables | Hurts scan | Compact density + clear primary |
| New font stack per feature | Breaks system | Theme `--ui-font-*` unless DESIGN.md says otherwise |
| Layout tables for non-tabular content | GOV.UK table guidance | `DATA-LIST` / cards ([Data display](/docs/data-display)) |

**Allowed when the theme/brief says so:** Material type stack, Basecoat direction, product brand fonts recorded in consumer DESIGN.md.

## Consumer context files

For apps **using** gelium-ui (not the monorepo itself):

| File | Purpose |
|---|---|
| `PRODUCT.md` | Who it’s for, jobs, voice, non-goals |
| `DESIGN.md` | Lane (Operate/Persuade/…), theme choice (`theme-material` / `theme-basecoat`), anti-references, density |

Templates: [`/docs/templates/product`](/docs/templates/product) and [`/docs/templates/design`](/docs/templates/design) (also under `site/web/content/templates/`).

Agents: load these before `WF-SHAPE` when present.

## Detectors

Run from repo root:

```bash
bash scripts/ux-detect.sh
```

Deterministic checks for handbook contracts, agent pack IDs, and hard anti-patterns (e.g. masking overflow on `body`/`html` in library/site CSS). Exit non-zero on failure.

## Map to existing docs

| Pass | Deep links |
|---|---|
| Shape | [Screens](/docs/screens), [Page + section architecture](/docs/page-section-architecture), [Journeys](/docs/journeys), [Data display](/docs/data-display), [Patterns](/docs/patterns) |
| Architecture | [Page + section architecture](/docs/page-section-architecture) — `SECTION-CONTRACT` before components or tokens |
| Build | Components sidebar, [Forms](/docs/forms), [Server contracts](/docs/server-contracts) |
| Audit | [Page + section architecture](/docs/page-section-architecture) `WF-SECTION-AUDIT`, [UI definition of done](/docs/ui-definition-of-done), [Feedback](/docs/feedback), [Responsive](/docs/responsive), [Accessibility](/docs/accessibility) |
| Polish | [Themes](/docs/themes), [Tokens](/docs/tokens), [Content style](/docs/content-style), [Density](/docs/density) |

## See also

- Agent pack: [`/llms-ux.txt`](/llms-ux.txt)  
- [UI definition of done](/docs/ui-definition-of-done)
