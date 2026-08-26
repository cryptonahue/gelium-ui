# Page + Section Architecture Protocol Implementation Plan

> **For Hermes:** Implement task-by-task with `subagent-driven-development`; keep the protocol English in every agent-facing and public artifact.

**Goal:** Teach LLM consumers to make and verify design decisions from product intent through page and section composition, rather than stopping at a component catalog or a page-level recipe.

**Architecture:** Add one canonical machine-facing skill, **Page + Section Architecture**, and one human-readable Handbook page. The protocol inserts an explicit Section layer into Gelium’s existing decision chain: `Product intent → Page/Surface → Section → Components → Tokens/Skin`. It reuses the existing Public Feed and Rich Article recipes as worked applications; no new primitive, runtime route, server contract, CSS system, or duplicate “post detail” demo is required.

**Tech stack:** Markdown guidance, Go docs routing/embedded assets, existing recipe templates/tests, Bash UX detector, Go test, Tailwind build.

---

## 1. Verified current-state map

### What already exists

| Surface | Existing source | What it provides | Protocol gap |
|---|---|---|---|
| Installed agent entry point | `lib/AGENTS.md`, `lib/SKILLS.md`, `lib/llms.txt`, `lib/llms-ux.txt`, `lib/skills/01-09` | Foundations, screen selection, product discovery, DoD, and usability checks | No mandatory decision between a selected page and its component list: sections have no named purpose, contract, hierarchy, or review routine. |
| Human workflow | `site/web/content/handbook-agent-workflow.md`, `handbook-screens.md`, `handbook-ui-definition-of-done.md`, `handbook-patterns.md` | WF-BRIEF → WF-SHAPE → WF-BUILD → WF-AUDIT → WF-POLISH; screen type, page hierarchy, DoD, skeletons | The page is still the smallest compositional unit. “Section” appears as layout prose, not as a purpose-bound grammar. |
| Existing composition model | `docs/gelium-ui-composition-rules.md`, `docs/gelium-ui-screen-composition.md`, `docs/gelium-ui-screen-recipes.md` | Screen grammar, 19-field recipe contract, pattern selection, state matrix, components | Recipe fields enumerate a screen but do not require each major page region to declare its job, audience, primary action relationship, revelation, or boundary. |
| Feed worked example | `/recipes/public-feed`; `internal/app/recipe_public_feed.go`; `site/web/templates/recipe-public-feed.html`; `site/web/styles/recipe-public-feed.css`; `internal/app/recipe_public_feed_test.go` | Read/scan feed, tabs, card list, states, POST+303 reactions, HTMX panel swap | Page regions are implemented but are not labeled or assessed as purpose-specific sections. |
| Detail worked example | `/recipes/rich-article`; `internal/app/recipe_rich_article.go`; `site/web/templates/recipe-rich-article.html`; `site/web/styles/recipe-rich-article.css`; `internal/app/recipe_rich_article_test.go` | Existing read/detail analogue with article header, prose, media, related activity, recoverable states, and related navigation | It is named “Rich Article,” not a new `post-detail` route, but is the appropriate existing post-detail/read-detail demonstrator. |
| Public docs delivery | `internal/app/docs.go`, `internal/app/server.go`, `internal/app/handbook_test.go`, `site/web/assets.go` | Handbook nav, routes, search index, previous/next, footer, sitemap, and embedding are derived from `handbookSections` | A new Handbook page must be registered in the same model; adding only a markdown file would leave it undiscoverable. |
| Mechanical verification | `scripts/ux-detect.sh`, `site/web/styles_recipe_mobile_test.go`, `lib/styles_contract_test.go`, recipe-specific Go tests | Contract IDs, docs existence, shell/color checks, mobile containment, compiled CSS, rendered HTML, route/metadata/mutation checks | No check currently requires the new protocol in the installed/served LLM surfaces or proves that both worked applications expose every required section contract. |

### Important constraints discovered

1. `site/web/static/llms-ux.txt` is the file served at `/llms-ux.txt`; it is not byte-identical to `lib/llms-ux.txt`. Update both explicitly and add a parity-focused test/check for the new protocol IDs instead of assuming a build copies one to the other.
2. `site/web/assets.go` embeds `content/*.md` and `static/*`; the proposed Handbook markdown is automatically embedded, but its route still must be registered.
3. `handbookSections` in `internal/app/docs.go` drives sidebar navigation, docs search, footer, previous/next navigation, and sitemap through `handbookNavLinks()`. Add the protocol once there rather than maintaining separate lists.
4. Existing recipes are intentionally `noindex, nofollow`; preserve that route policy. The public Handbook page is indexable through the existing handbook registry.
5. Keep the scope documentation and verification only. Do not create a component, token, CSS utility, or “section” runtime abstraction; the protocol governs composition of the components already in the registry.

---

## 2. Protocol contract to introduce

### Required decision order

Every new or substantially changed product surface must record, in this order:

1. **Product intent** — audience, user job, product outcome, non-goal, and relevant lifecycle/event consequences from `skills/08-product-reasoning.md`.
2. **Page/Surface** — URL-level surface mode, screen type, primary action, journey/data/feedback IDs, chrome, and server/no-JS contract.
3. **Section** — each major region has a specific purpose and an explicit boundary. A section is not “a card” or “some spacing”; it is a coherent answer to one reader need within the page.
4. **Components** — registered semantic primitives only, selected to fulfill the section contract; no new component because a composition decision is missing.
5. **Tokens/Skin** — existing semantic type, spacing, color, motion, and skin tokens only. Tokens express the chosen relationship; they do not substitute for section purpose.

### Section contract template

The canonical skill and the Handbook page must include the same recognizable template (the skill is normative; the Handbook explains it):

```text
SECTION-ID / name:
Purpose: What user need this region satisfies now.
Audience and moment: Who reads it, and at what point in the page/journey.
Input / evidence: Data, state, or context it receives.
Output / next move: What understanding or action the reader leaves with.
Hierarchy: Entry, primary, supporting, or recovery; why it has that rank.
Action policy: Primary / secondary / none; relation to the page’s single primary action.
Revelation: Always visible / conditional server-rendered state / user-invoked disclosure; why.
Composition: Semantic landmark or element plus registered Gelium components.
State and recovery: Rest, loading, empty, error, success where applicable; no-JS behavior.
Boundary and rhythm: What belongs inside vs outside; relationship-based --ui-space-* gaps.
Accessibility/content: heading/label/landmark, reading order, announcement, concise factual copy.
Verification evidence: rendered-HTML assertions and mechanical checks that prove the contract.
```

### Rules the protocol must make executable

- **Purpose rule:** Every top-level page child is assigned a section purpose. Decorative wrappers and repeated card anatomy are not automatically sections.
- **Audience-and-moment rule:** State which reader and moment a section serves. Do not present authoring controls, recovery, and reader context as undifferentiated content.
- **Hierarchy rule:** Establish entry → primary → supporting → recovery order. The document order must match the intended reading order; a later support section cannot visually impersonate the page’s primary task.
- **Action rule:** The page retains one primary action. A section may own only a supporting/local action unless it is the page’s declared primary-action section.
- **Revelation rule:** Put identity and the task-critical “now” information in the initial page scan. Use server-rendered conditional sections for meaningful state; use native disclosure only for optional depth; never hide a necessary action or recovery only behind hover/JS.
- **Boundary rule:** Group elements that answer one need; separate unrelated needs. Use the existing spacing contract (`space-1/2` own metadata, `space-3` sibling items, `space-4/6` groups, `space-8` section-to-section) and semantic landmarks/headings, not arbitrary visual containers.
- **Component rule:** A component selection follows section purpose. “Card because it looks organized” is invalid; a card is for a repeated self-contained instance, while panels/sections organize page context.
- **State/recovery rule:** A data or action section identifies its applicable rest/loading/empty/error/success behavior and the server/no-JS path. Do not assign page-global state blindly to every section.
- **Critique rule:** Reviewers can trace each rendered page region back through `Section → Page/Surface → Product intent`; a gap is either missing intent, an uncontracted section, an unregistered component, or a token/skin misuse.

### Executable critique workflow

Add `WF-ARCH` as a bounded step between `WF-SHAPE` and `WF-BUILD`, and `WF-SECTION-AUDIT` inside `WF-AUDIT`:

1. Load/confirm `PRODUCT.md` and `DESIGN.md`; run the existing product-reasoning inventory before markup.
2. Name the page/URL, surface mode, screen type, one primary action, and existing `JOURNEY-*`, `DATA-*`, and `FEED-*` IDs.
3. Draw a numbered DOM-outline of major regions only (no CSS or components yet).
4. Complete one section contract per major region; delete, merge, split, or escalate any region without a distinct purpose.
5. Map each section to existing semantic elements/components and relationship-based tokens. Reject new primitives unless the vocabulary/registry process proves the gap.
6. Run a **five-second scan**: title/context/primary action/task-critical status must be recognizable before optional depth.
7. Run **action competition**: count page-primary variants; section-local actions must be visually/semantically subordinate unless they are the page primary action.
8. Run **revelation/recovery**: identify the first visible state, every conditional state, the no-JS path, and a way out of error/empty/partial states.
9. Inspect rendered HTML in narrow and dark/class-routed conditions; then run tests/detectors. Record only failures with a contract ID and rendered evidence.

---

## 3. Smallest coherent artifact set

### New artifacts

1. **`lib/skills/10-page-section-architecture.md`** — canonical, installed, machine-facing protocol. It contains the decision order, section contract template, rules, critique workflow, and both compact worked applications. It is the only new agent skill.
2. **`site/web/content/handbook-page-section-architecture.md`** — human-facing explanation of the same protocol, with short tables/checklists and links to the existing feed and rich-article routes. It must describe `/recipes/rich-article` honestly as the existing post-detail/read-detail analogue, not claim a nonexistent `/recipes/post-detail` route.

### Existing artifacts to extend

| File | Planned change |
|---|---|
| `lib/SKILLS.md` | Add skill 10 and update the ordered workflow so agents read the protocol after product reasoning/surface selection and before component selection. |
| `lib/AGENTS.md` | Change the “01-09” guidance and golden workflow to require the section protocol before markup; retain existing rules as floors. |
| `lib/llms.txt` | Add the Handbook route and state that `/llms-ux.txt` includes the page/section protocol. |
| `lib/llms-ux.txt` | Add concise `ARCH-*` / `SECTION-*` decision IDs, the chain, a compact contract, and `WF-ARCH` / `WF-SECTION-AUDIT`; link readers to the Handbook and skill for details. |
| `site/web/static/llms.txt` | Mirror the new public route and agent-pack discovery statement because this is the file served at `/llms.txt`. |
| `site/web/static/llms-ux.txt` | Mirror the compact protocol IDs because this is the file served at `/llms-ux.txt`; do not rely on `lib/llms-ux.txt` being copied. |
| `site/web/content/handbook-agent-workflow.md` | Insert `WF-ARCH` after Shape and define the bounded section audit within Audit; link to the new Handbook page. |
| `site/web/content/handbook-screens.md` | Add the handoff from page hierarchy to section contracts and a short rule that a “main content” block must be decomposed by purpose before component selection. |
| `site/web/content/handbook-ui-definition-of-done.md` | Add section-level acceptance checks: purpose, hierarchy, action competition, revelation, state/recovery, boundaries/rhythm, and evidence. |
| `site/web/content/handbook-patterns.md` | Extend `SKEL-FORUM` with an explicit link between the feed/topic list and the topic detail/thread; point to the two worked applications. Do not create a new skeleton unless evidence requires one. |
| `internal/app/docs.go` | Add the Core Handbook nav link for `/docs/page-section-architecture` immediately after Screens; add `docsPageSectionArchitecture`; add one curated docs-hub link rather than duplicating the whole Handbook. This automatically feeds search, footer, previous/next, and sitemap via the existing model. |
| `internal/app/server.go` | Register `GET /docs/page-section-architecture` with the new docs handler near other Handbook routes. |
| `internal/app/handbook_test.go` | Add the route to `handbookRoutes` with its exact H1 and a unique protocol marker; assert the nav/sitemap route through the existing table-driven tests. |
| `scripts/ux-detect.sh` | Require the Handbook file, served LLM files, and stable protocol IDs (for example `ARCH-PRODUCT`, `SECTION-CONTRACT`, `WF-ARCH`) so a partial docs-only update fails. Add a targeted check that both worked-application names/paths occur in the canonical guidance. |
| `internal/app/recipe_public_feed_test.go` | Extend rendered-HTML contracts only: assert the feed’s major semantic regions/labels required by the protocol are present and ordered (page context, view selection, feed list, recovery/refresh). Keep existing behavior tests intact. |
| `internal/app/recipe_rich_article_test.go` | Extend rendered-HTML contracts only: assert post-detail section order and semantic anchors (identity/header, body/prose, evidence/media, related activity, recoverable states, related navigation). Preserve current route and noindex/JSON-LD expectations. |

### Explicit non-changes

- Do **not** add a `Section` Go type, a section CSS primitive, new tokens, JavaScript, a database/store, or another route for post detail.
- Do **not** alter `site/web/templates/recipe-public-feed.html`, `recipe-rich-article.html`, their handlers, their styles, or source components unless the planned rendered-HTML audit finds that a semantic marker cannot be asserted. The current templates already provide the two applications; the first implementation should document/test their contracts, not redesign them.
- Do **not** edit the historical Spanish architecture documents (`docs/gelium-ui-screen-recipes.md`, `docs/gelium-ui-screen-composition.md`) in the initial slice. The new canonical protocol and all additions are English; a later translation/reconciliation can be separately scoped.
- Do **not** modify `lib/scripts/ux-detect.sh`; it is a consumer-repository detector and cannot reliably inspect Gelium’s own embedded Handbook assets.

---

## 4. Worked applications to author in the protocol

### Application A — Public Feed (`/recipes/public-feed`)

Use the real files `site/web/templates/recipe-public-feed.html` and `internal/app/recipe_public_feed.go`.

| Ordered section | Purpose and judgment | Existing composition / contract | Critique evidence |
|---|---|---|---|
| Page context | Establish that the reader is scanning current activity and where they are. | `recipe-pf-header` with H1 and concise description. | The H1/description precede interaction and list content. |
| View selection | Let the reader change the feed lens without changing the page’s primary job. | `nav.ui-tabs`, real GET links, `aria-current`, `?view=` URL state. | Tab links are before the list and no `role=tablist` JS contract is claimed. |
| Activity feed | Deliver the primary read/scan task. | `ol.recipe-pf-list` of `article.ui-card`; author/time/newness/body/actions. | Reverse chronology, item semantics, avatar/name, and new-state label remain rendered. |
| Local reaction | Let a reader react without competing with discovery as the page purpose. | Per-item native form `POST /recipes/public-feed/{id}/react` → 303 + transient flash. | Like stays local to an item; it is not promoted to a page-level primary action. |
| Collection recovery | Explain empty/loading/error conditions and preserve a way forward. | Shared Empty state, documented Skeleton, 404 error state, real CTA/retry. | Empty output contains a CTA; loading does not flash empty; invalid reaction reaches 404 recovery. |
| Refresh support | Allow a non-primary remote refresh with a no-JS and HTMX-consistent result. | Native POST refresh, inline toast/progress fallback, fragment enhancement. | GET remains 405; no-JS render and HX response are both exercised. |

The worked application must call out a deliberate distinction: the feed card is a repeated entity/event instance, not automatically a page section. The feed list is the section; header/body/footer are card anatomy.

### Application B — Existing Rich Article as post detail (`/recipes/rich-article`)

Use the real files `site/web/templates/recipe-rich-article.html` and `internal/app/recipe_rich_article.go`. Call it **“Post detail / rich article”** in the protocol so consumers recognize the design problem while the source route remains truthful.

| Ordered section | Purpose and judgment | Existing composition / contract | Critique evidence |
|---|---|---|---|
| Context and identity | Identify the post, source/context, author/date, and reading promise before deep content. | Article header with eyebrow, one H1, lead, byline. | One H1; identity metadata stays grouped before prose. |
| Primary reading body | Let the reader understand the article in a readable measure. | `recipe-rich-article-prose`, descriptive H2/H3, paragraphs/lists/quote/code. | Heading order and prose measure are visible in HTML/CSS. |
| Evidence/media | Support comprehension where media adds evidence, with a recovery path rather than decorative interruption. | Local picture, native video/captions/fallback, audio/transcript, explicit safe-embed fallback. | Alt text, dimensions, captions, transcript, and fallback links exist. |
| Supporting reference | Provide structured facts without hijacking the reader’s primary read flow. | Data table section with caption and internal scroll primitive. | It follows the prose/evidence context and is table-semantic. |
| Related activity | Offer adjacent context after the primary reading need is satisfied. | `ui-list` under “Related activity.” | It is a supporting section, not page chrome or a competing hero. |
| Recoverable states | Make absence/loading/failure legible and actionable. | Skeleton, empty state, error alert within a dedicated states section. | State labels and role semantics remain rendered. |
| Related navigation | Give a reader a next path after understanding the post. | `nav[aria-label="Related content"]` with real links. | It occurs after content/recovery and uses navigation semantics. |

The worked application must identify a current limitation instead of hiding it: it is an integration fixture with no mutation/reply action, so the protocol should not invent a comment composer or claim it demonstrates a complete forum thread. It demonstrates **read-detail section architecture**, not the full `SKEL-FORUM` lifecycle.

---

## 5. Implementation tasks

### Task 1: Write the canonical agent protocol

**Objective:** Give installed LLM consumers one normative, scannable source for page-to-section judgment.

**Files:**
- Create: `lib/skills/10-page-section-architecture.md`
- Modify: `lib/SKILLS.md`
- Modify: `lib/AGENTS.md`
- Modify: `lib/llms.txt`
- Modify: `lib/llms-ux.txt`
- Modify: `site/web/static/llms.txt`
- Modify: `site/web/static/llms-ux.txt`

**Steps:**
1. Add the skill with the required chain, section contract template, rules, bounded critique workflow, and both worked applications above.
2. Use stable compact IDs consistently across the skill and packs: `ARCH-PRODUCT`, `ARCH-PAGE`, `ARCH-SECTION`, `ARCH-COMPONENTS`, `ARCH-TOKENS`, `SECTION-CONTRACT`, `SECTION-HIERARCHY`, `SECTION-ACTION`, `SECTION-REVELATION`, `SECTION-RECOVERY`, `WF-ARCH`, and `WF-SECTION-AUDIT`.
3. Update the installed entry points so an LLM cannot reach components/markup from skills 01–09 without seeing the section decision step.
4. Update both served static LLM files in the same change; preserve their current public URLs and concise format.

**Acceptance criteria:**
- The package exposes exactly one new skill at `lib/skills/10-page-section-architecture.md`.
- AGENTS, SKILLS, llms, and both llms-ux surfaces name the same protocol IDs and the new step.
- Both applications name existing routes/files only; neither claims `/recipes/post-detail` exists.
- Guidance says “section purpose first, components second,” and continues to require the pre-existing product/artifact gate.

### Task 2: Publish the human Handbook page and wire it into the docs model

**Objective:** Make the protocol discoverable to human reviewers through the existing docs shell, navigation, search, footer, previous/next, and sitemap machinery.

**Files:**
- Create: `site/web/content/handbook-page-section-architecture.md`
- Modify: `site/web/content/handbook-agent-workflow.md`
- Modify: `site/web/content/handbook-screens.md`
- Modify: `site/web/content/handbook-ui-definition-of-done.md`
- Modify: `site/web/content/handbook-patterns.md`
- Modify: `internal/app/docs.go`
- Modify: `internal/app/server.go`
- Modify: `internal/app/handbook_test.go`

**Steps:**
1. Write an English Handbook page using the project’s cognitive-document shape: answer-first lead; quick path; decision chain; section contract table; rules; the two worked applications; critique checklist; sources/see-also.
2. Place `/docs/page-section-architecture` in the Core group after `/docs/screens`, then add the matching handler and GET route.
3. Add one curated link to the docs hub’s Core list, without turning the hub into a full handbook index.
4. Update workflow/Screens/DoD/Patterns by link and concise handoff, not by copying the full protocol into four places.
5. Add table-driven handbook test coverage with a distinct H1 and marker such as `SECTION-CONTRACT`; preserve all existing routes and group order invariants.

**Acceptance criteria:**
- `GET /docs/page-section-architecture` returns 200 under the docs shell with one expected H1 and an identifiable protocol marker.
- The route appears in the Core sidebar, search index, footer, previous/next navigation, and sitemap because it enters `handbookSections`/`handbookNavLinks`.
- Workflow documents `WF-ARCH` between Shape and Build and a finite section audit under Audit.
- DoD contains page-and-section checks in addition to the current page-level checklist.

### Task 3: Make the two existing recipes executable applications of the protocol

**Objective:** Prove that the protocol reviews real Gelium output instead of presenting a purely conceptual template.

**Files:**
- Modify: `internal/app/recipe_public_feed_test.go`
- Modify: `internal/app/recipe_rich_article_test.go`
- Modify only if the tests expose a missing semantic anchor: `site/web/templates/recipe-public-feed.html`, `site/web/templates/recipe-rich-article.html`

**Steps:**
1. Add table-driven section-order assertions to the feed test using existing identifiers/semantics: page header, feed view nav, feed panel/list, item actions, loading/recovery, refresh.
2. Add table-driven section-order assertions to the rich-article test: article header, prose, media/evidence, data reference, related activity, recoverable states, related navigation.
3. Assert semantic containers/labels rather than cosmetic class strings whenever possible. A test should fail when a purpose-bound region disappears, moves before its prerequisite, loses its landmark/heading, or loses recovery semantics.
4. Do not add implementation just to satisfy a brittle test. If existing markup lacks a stable semantic marker needed for the contract, add the smallest semantic attribute/heading/landmark and update the matching template test.

**Acceptance criteria:**
- Feed tests distinguish the feed-list section from repeated card anatomy and continue to pass all existing POST+303, fragment, empty, and metadata assertions.
- Rich-article tests recognize it as the existing post-detail/read-detail application, retain one H1 and valid Article JSON-LD, and prove supporting/recovery sections follow—not replace—the primary reading section.
- No new route, handler, store, component, CSS file, or behavior contract is added.

### Task 4: Add protocol-specific mechanical checks and run the complete validation sequence

**Objective:** Prevent the new protocol from drifting into only one docs/agent surface and verify the docs site still embeds and serves it.

**Files:**
- Modify: `scripts/ux-detect.sh`
- Optionally modify: `internal/app/handbook_test.go` only if a separate cross-surface ID parity test is clearer than extending its existing table-driven test.

**Steps:**
1. Make `scripts/ux-detect.sh` require the new handbook markdown and stable protocol IDs in `site/web/static/llms-ux.txt`.
2. Add explicit checks that the agent-facing guidance names both `/recipes/public-feed` and `/recipes/rich-article`; this prevents the applications from being silently removed while the generic protocol remains.
3. Add a focused Go test or detector check that verifies the package and served LLM surfaces contain the same stable ID set. Do not require whole-file equality because the files already legitimately differ in surrounding content.
4. Run the commands below from repository root; review rendered route HTML where indicated.

**Acceptance criteria:**
- Deleting the Handbook page, removing a protocol ID from served `/llms-ux.txt`, or dropping either worked application makes a deterministic check fail.
- The verification stays feasible without a browser automation dependency: it relies on Go-rendered HTML, static content checks, the existing detector, CSS/build checks, and HTTP route tests.

---

## 6. Verification commands and evidence

Run from `/root/.openfang/workspaces/repos/gelium-ui` after implementation:

```bash
# Focused docs, guidance, and worked-application tests.
go test -count=1 ./internal/app ./site/web ./lib

# Full repository Go validation required by Gelium guidance.
go test -count=1 ./...
go vet ./...

# Rebuild both the docs-site asset and publishable library bundle.
npm run build

# Resolve the workspace package exactly as the docs consumer does.
npm run test:workspace-resolution

# Run monorepo documentation/agent contract detectors.
bash scripts/ux-detect.sh

# Serve the newly built binary for route-level smoke checks.
go build -o /tmp/gelium-ui ./cmd/gelium-ui
/tmp/gelium-ui
```

In a second shell, verify the rendered/served outcomes (replace port only if the server reports another value):

```bash
curl -fsS http://127.0.0.1:8080/docs/page-section-architecture | grep -F 'Page + Section Architecture'
curl -fsS http://127.0.0.1:8080/llms-ux.txt | grep -F 'SECTION-CONTRACT'
curl -fsS http://127.0.0.1:8080/recipes/public-feed | grep -F 'id="feed-panel"'
curl -fsS http://127.0.0.1:8080/recipes/rich-article | grep -F 'Related content'
curl -fsS http://127.0.0.1:8080/sitemap.xml | grep -F '/docs/page-section-architecture'
```

Before reporting completion, inspect `git diff --check` and `git status --short`; confirm the diff contains only the protocol/docs/tests/build outputs expected by the implementation work. This planning task itself must leave only this plan file changed and must not commit or push.

---

## 7. Final acceptance checklist

- [ ] LLM consumers have an installed skill and concise pack IDs that enforce `Product intent → Page/Surface → Section → Components → Tokens/Skin`.
- [ ] The protocol includes an actionable, reusable section contract template—not merely descriptive prose.
- [ ] It defines visual hierarchy, action competition, intended audience/moment, revelation, boundaries/rhythm, state/recovery, and component/token mapping rules.
- [ ] `WF-ARCH` and `WF-SECTION-AUDIT` form a bounded, executable critique workflow with rendered evidence.
- [ ] Public Feed and existing Rich Article are documented as the two real worked applications; no invented post-detail route or fake component is introduced.
- [ ] The human Handbook page is a routed, navigable, searchable, indexable Core document; recipes remain noindex.
- [ ] Existing recipe tests prove the section contracts on rendered HTML without turning presentation details into brittle requirements.
- [ ] `scripts/ux-detect.sh`, focused tests, full tests/vet, build, workspace-resolution test, and route smoke checks pass.
- [ ] No production source/component/example change occurs unless a test demonstrates that the existing semantic markup lacks the minimal anchor required for a contract assertion.
- [ ] No commit or push is made.

## Risks and mitigations

| Risk | Mitigation |
|---|---|
| Guidance is duplicated and drifts between npm, served static packs, and Handbook docs. | Make the new skill normative, use compact IDs elsewhere, update both independently served static files, and assert the stable ID set mechanically. |
| “Section” becomes a new visual wrapper/component taxonomy. | Define it as a purpose-and-boundary contract; explicitly prohibit a new runtime abstraction, token family, or primitive in this scope. |
| The new protocol overrules product discovery or DoD. | Place it after existing product/surface reasoning and before build; it feeds existing DoD/usability checks rather than replacing them. |
| The rich article is misrepresented as a complete social post/thread. | Label it post-detail/read-detail analogue, preserve the real route/name, and document its no-mutation limitation. |
| Recipe tests become tightly coupled to CSS class names. | Prefer semantic tags, labels, IDs, and ordered rendered text; add a minimal stable semantic marker only when necessary. |
| Docs routing is manually duplicated and becomes invisible in one destination. | Register through `handbookSections`; existing route/nav/footer/search/sitemap tests prove the shared model. |
