# Screens

UI starts with **what kind of screen this is**, not with which component looks pretty. This page defines screen types, hierarchy rules, and a build checklist. Criteria are adapted from public design-system guidance (GOV.UK, USWDS, Material 3) and research (Nielsen Norman Group), then mapped onto Gelium’s server-rendered, 0-JS-first model.

## Sources (read these)

| Topic | Source |
|---|---|
| Start / journey framing | [GOV.UK: start using a service](https://design-system.service.gov.uk/patterns/start-using-a-service/) |
| Long multi-task flows | [GOV.UK: task list](https://design-system.service.gov.uk/components/task-list/) |
| Side / section navigation | [USWDS: side navigation](https://designsystem.digital.gov/components/side-navigation/) |
| Compact app destinations | [M3: navigation bar](https://m3.material.io/components/navigation-bar/guidelines) |
| How people scan pages | [NNG: F-shaped pattern](https://www.nngroup.com/articles/f-shaped-pattern-reading-web-content/) |
| Control choice (forms) | [Choose the right control](/docs/choose-the-right-control), [Forms](/docs/forms) |
| After submit feedback | [Feedback](/docs/feedback) |

Gelium **adapts** these sources: we do not ship React shells; we ship HTML partials, tokens, and server contracts.

## Screen types

Pick **one primary type** per URL. Mixing “settings + giant table + marketing hero” on one route usually means split the page.

| Type | User job | Typical Gelium building blocks | Notes |
|---|---|---|---|
| **Hub / start** | Understand what this area is and begin | `hero` or title + short prose, **one** primary `button`, links out | GOV.UK start pages: enough context to begin, not the whole manual |
| **List** | Find an item among many | `data-table` or `list`/`card`, filters (`text-field`, `select`, `chips`), `pagination`, `empty-state` | Recipe: [Admin Resource](/recipes/admin-resource), [Public feed](/recipes/public-feed) |
| **Detail** | Understand one thing | title, metadata, `divider`, actions, related lists | Lead with identity (name/id), then status, then body |
| **Form** | Provide or edit data | `text-field`, `select`, `checkbox`/`radio`/`switch`, `validation-summary`, primary + secondary actions | See [Forms](/docs/forms); validate on submit / after interaction, not hostile early errors ([NNG error timing](https://www.nngroup.com/articles/error-message-guidelines/)) |
| **Confirm** | Decide on a consequential action | `dialog` page variant or confirm route, destructive `button`, cancel link | Prefer real URL + POST for destructive work (Gelium dialog recipe) |
| **Settings** | Change durable preferences | `list` of rows, `switch`/`select`, save if needed | One setting group per view when possible |
| **Queue / work list** | Process items in order or by status | list/table + status `badge` + row actions | Recipe: [Ops Queue](/recipes/ops-queue) |
| **Result / confirmation** | Know it worked and what next | success `banner` or clear H1, reference id, next links | Do not only flash a toast and leave the user nowhere |

### When not to invent a new type

If the job is “pick a value in a form,” that is still a **form** (maybe a step), not a new screen taxonomy. Prefer fewer types and clearer IA.

## Hierarchy on the page

Adapted from [NNG F-pattern](https://www.nngroup.com/articles/f-shaped-pattern-reading-web-content/) and common service-pattern practice:

1. **Title** — what is this page (one H1).
2. **Context** — one short line or status if needed (not a paragraph wall).
3. **Primary action** — the main thing the user can do next (usually one visible primary `button`).
4. **Main content** — list, form, or detail body.
5. **Secondary actions** — outline/secondary buttons, row actions, “back”.

| Rule | Do | Don’t |
|---|---|---|
| Primary actions | One obvious primary per view | Three filled primary buttons competing |
| Scanning | Put identity and action in the top band | Bury the only CTA under a long table |
| Density | Admin lists may be dense; reading views stay measured (`65ch` prose) | Stretch paragraphs full bleed “because desktop” |
| Empty | `empty-state` with reason + next step | Blank table with no explanation |

## Navigation patterns (product IA)

| Pattern | When to use | When not | Gelium |
|---|---|---|---|
| **Side nav (section hierarchy)** | 1–3 levels inside a product area; many destinations | Tiny site (&lt; ~5 pages) — USWDS: consider no hierarchy | Docs shell sidebar; `navigation-drawer` for app chrome |
| **Top / bar destinations** | Few peer areas the user switches often | Deep trees of nested pages | `navigation-bar`, `navigation-tab` |
| **Tabs** | Alternate **views of the same object** or peer modes | Primary site IA for everything; avoid stacking tab bars | `tabs`, `segmented-button` for local mode |
| **Bottom nav (compact apps)** | ~3–5 top-level app areas on small viewports (M3 nav bar guidance) | More than ~5 peers; deep hierarchy | `navigation-bar` patterns in components |
| **In-page (on this page)** | Long single document sections | Short pages | Docs on-this-page rail |

USWDS: side navigation is for **section sub-navigation**, not a substitute for simplifying a small site ([side navigation](https://designsystem.digital.gov/components/side-navigation/)).  
M3: navigation bars are for switching top-level views on smaller windows and should stay consistent ([navigation bar](https://m3.material.io/components/navigation-bar/guidelines)).

## Build checklist (humans and agents)

Before writing markup:

1. **Screen type** — hub | list | detail | form | confirm | settings | queue | result.
2. **One user job** — finish the sentence: “The user came here to ___.”
3. **One primary action** — name it; demote everything else.
4. **States** — loading, empty, error, success (see [Feedback](/docs/feedback)).
5. **Controls** — map fields/actions via [Forms](/docs/forms) and [Choose the right control](/docs/choose-the-right-control).
6. **Server** — GET for read, POST+redirect for mutates, 422 + summary for validation ([Server contracts](/docs/server-contracts)).
7. **Theme** — `theme-*` on `<html>`; no one-off hex.
8. **Narrow viewport** — stack headers, `min-width: 0`, local table scroll ([Responsive](/docs/responsive)).

## Anti-patterns

- Dumping the full component catalog on a hub page (the nav already does that).
- Using a dialog for simple external navigation.
- Validation only in a toast (see [Feedback](/docs/feedback)).
- Desktop-only CSS patched with `overflow-x: hidden` on `body`.

## See also

- [Feedback](/docs/feedback) — which message component for which situation  
- [Performance](/docs/performance) — payload stance  
- [Why Gelium](/docs/compare) — when not to use this stack  
- Agent pack: [`/llms-ux.txt`](/llms-ux.txt)
