# Page + section architecture

A coherent page is an ordered set of purpose-bound sections, not an unordered component catalog. Decide **why each major region exists** before choosing components, CSS, or spacing. This keeps the primary task visible, supporting information subordinate, and recovery actionable in server-rendered, no-JS flows.

## Quick path

1. Name the product intent, URL, surface mode, screen type, and one primary action.
2. Draw the ordered major regions of the page; do not count repeated card anatomy as page sections.
3. Complete a `SECTION-CONTRACT` for every region, then map it to semantic HTML, registered Gelium components, and existing tokens.
4. Inspect rendered HTML and run the bounded section audit before handoff.

## Decision chain

Follow this order. A later decision must not replace an earlier one.

| ID | Decide | Record |
|---|---|---|
| `ARCH-PRODUCT` | Product intent | Audience, user job, outcome, non-goal, and relevant lifecycle consequences |
| `ARCH-PAGE` | Page / surface | URL, surface mode, screen type, one primary action, journey/data/feedback IDs, chrome, and server/no-JS contract |
| `ARCH-SECTION` | Section purpose | One purpose-bound contract for every major region |
| `ARCH-COMPONENTS` | Components | Registered semantic primitives that fulfill each section contract |
| `ARCH-TOKENS` | Tokens / skin | Existing type, spacing, color, motion, and class-routed skin tokens |

```text
Product intent → Page / surface → Section → Components → Tokens / skin
```

Start with the [agent workflow](/docs/agent-workflow), [Screens](/docs/screens), and the product brief/design artifacts when they exist. This protocol extends their existing registry-first, token, and definition-of-done gates; it does not introduce a new runtime, component family, or visual language.

## SECTION-CONTRACT

A section is one coherent answer to one reader need at a moment in the page or journey. It is not a card, decorative wrapper, or repeated item anatomy.

| Field | Decide before markup |
|---|---|
| Identity | `SECTION-ID` / name and the reader need it answers now |
| Audience and moment | Who reads it and when in the page or journey |
| Input and next move | Data/state it receives, plus the understanding or action it leaves the reader with |
| Hierarchy | Entry, primary, supporting, or recovery — and why it has that rank |
| Action policy | Primary, secondary, or none; relation to the page’s one primary action |
| Revelation | Always visible, conditional server-rendered state, or native disclosure — and why |
| Composition | Semantic landmark/element and registered Gelium components |
| State and recovery | Applicable rest, loading, empty, error, success, and no-JS behavior |
| Boundary and rhythm | What belongs inside/outside and relationship-based `--ui-space-*` gaps |
| Accessibility and evidence | Heading/label/landmark, reading order, announcement/copy, rendered-HTML assertions, and mechanical checks |

## Rules

### Purpose and hierarchy

- Give every top-level page child a purpose. Group elements that answer one need; separate a change of user intention with a semantic boundary.
- Order entry → primary → supporting → recovery in the DOM. A later support section cannot impersonate the primary task.
- Keep identity and task-critical “now” information in the initial scan. Use optional disclosure only for optional depth.
- Use relationship spacing, not arbitrary wrappers: `--ui-space-1/2` for an element and metadata, `--ui-space-3` for sibling items, `--ui-space-4/6` for groups, and `--ui-space-8` between sections.

### Actions, components, and revelation

- The page has one primary action. A section owns only supporting or local actions unless it is the declared primary-action section.
- Select components after purpose is known. A repeated self-contained entity may use a card; a card does not replace page context or a section boundary.
- Map sections to semantic HTML and registered primitives. Resolve names in `lib/ui-vocabulary.md`; escalate a new primitive only after the registry process proves a real gap.
- Never hide a necessary action or recovery path behind hover or JavaScript. Conditional state is server-rendered; optional depth uses native disclosure.

### State and recovery

- A data or action section declares its applicable rest, loading, empty, error, and success behavior plus a no-JS path.
- Do not duplicate page-global state in every region. Do preserve a way forward from empty, error, and partial states; persistent feedback is not a toast-only outcome.

## Workflow and critique

Run `WF-ARCH` between `WF-SHAPE` and `WF-BUILD`:

1. Confirm the product/design context and product-reasoning inventory.
2. Name the URL, surface mode, screen type, primary action, and relevant `JOURNEY-*`, `DATA-*`, and `FEED-*` IDs.
3. Draw a numbered DOM outline of major regions only — no CSS or component selection yet.
4. Complete one `SECTION-CONTRACT` per region; merge, split, remove, or escalate a region with no distinct purpose.
5. Map each region to semantic elements, registered components, and existing tokens.

Inside `WF-AUDIT`, run the finite `WF-SECTION-AUDIT`:

1. Five-second scan: title, context, primary action, and task-critical state are recognizable before optional depth.
2. Action competition: count page-primary variants; section-local actions remain subordinate unless declared primary.
3. Revelation and recovery: identify first visible and conditional states, the no-JS path, and a way forward from empty, error, or partial state.
4. Inspect rendered HTML at a narrow viewport and class-routed dark scheme, then run tests/detectors. Record only failures with a protocol ID and rendered evidence.

## Worked application: Public Feed

**Route:** [`/recipes/public-feed`](/recipes/public-feed)

| Ordered section | Existing purpose and contract |
|---|---|
| Page context | `recipe-pf-header`, its H1, and concise description establish scanning context before interaction or list content. |
| View selection | `nav.ui-tabs` uses real GET links, `aria-current`, and `?view=` state to change the feed lens without replacing the read/scan job. |
| Activity feed | `ol.recipe-pf-list` of `article.ui-card` items delivers the primary task in reverse chronology. |
| Local reaction | Native `POST /recipes/public-feed/{id}/react` → 303 reactions remain local and subordinate to discovery. |
| Collection recovery | Empty/loading/error states explain what happened and provide a CTA or retry path; loading does not flash empty. |
| Refresh support | A non-primary native POST refresh preserves no-JS and HTMX-consistent results. |

The feed list is the section. Each feed card is a repeated entity/event instance; its header, body, and footer are card anatomy, not automatic page sections.

## Worked application: Rich Article

**Route:** [`/recipes/rich-article`](/recipes/rich-article)

| Ordered section | Existing purpose and contract |
|---|---|
| Context and identity | The article header, eyebrow, one H1, lead, and byline define the reading promise before prose. |
| Primary reading body | `recipe-rich-article-prose` provides readable structured content with descriptive H2/H3 headings. |
| Evidence and media | Picture, video/captions/fallback, audio/transcript, and safe-embed fallback support comprehension. |
| Supporting reference | The captioned data table adds structured facts without hijacking the reading flow. |
| Related activity | `ui-list` gives adjacent context after the primary reading need. |
| Recoverable states | Skeleton, empty, and error alert states remain legible and actionable. |
| Related navigation | `nav[aria-label="Related content"]` offers a next path after content and recovery. |

This is the existing post-detail/read-detail analogue, not a `/recipes/post-detail` route. It is a read-detail integration fixture, not a complete forum-thread lifecycle.

## Critique checklist

- [ ] Every major region has a distinct purpose and `SECTION-CONTRACT`.
- [ ] DOM order matches entry → primary → supporting → recovery.
- [ ] One page-primary action is visible; local actions are subordinate.
- [ ] Each data/action section has applicable state and a recoverable no-JS path.
- [ ] Semantic boundaries, labels, landmarks, and heading order support reading order.
- [ ] Repeated card anatomy is not mistaken for page architecture.
- [ ] Rendered HTML, narrow viewport, class-routed dark scheme, tests, and detectors provide evidence.

## Sources and see also

| Topic | Source |
|---|---|
| Screen type, hierarchy, and navigation | [Screens](/docs/screens) |
| Bounded implementation passes | [Agent workflow](/docs/agent-workflow) |
| Ship checklist | [UI definition of done](/docs/ui-definition-of-done) |
| Domain skeletons | [Patterns](/docs/patterns) |
| Cited originals vs Gelium remakes | [Section references](/docs/section-references) |
| Product discovery | [GOV.UK patterns](https://design-system.service.gov.uk/patterns/) and [Material guidance](https://m3.material.io/) |

See also: [Journeys](/docs/journeys), [Data display](/docs/data-display), [Feedback](/docs/feedback), [Spacing](/docs/spacing), [Accessibility](/docs/accessibility), and the agent pack [`/llms-ux.txt`](/llms-ux.txt).
