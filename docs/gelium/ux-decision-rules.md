# Gelium UX Composition — Decision Rules

> **Artifact ID**: `gelium://ux-composition/decision-rules/v0.1` · **Version**: 0.1
> **Status**: normative within the protocol. Each rule has a stable ID that
> machine-readable contracts reference (`decision_rules` field in
> `schemas/screen-recipe.schema.json`).
>
> **Source policy**: rules that already exist in the system are cited and
> restated only in binding form (If/Then); the canonical prose lives in
> `docs/gelium-ui-composition-rules.md`, `docs/gelium-ui-ux-patterns.md`,
> `docs/gelium-ui-accessibility-contract.md`, `docs/gelium-ui-content-rules.md`.
> Rules that expand a documented gap (dashboard, tabs, action hierarchy,
> required states, responsive, accessibility, no-JS, errors, destructive) are
> new in this protocol and are marked **NEW**.

---

## How to read

Each rule is `R-<AREA>-<NN>` with a short name, an If/Then statement, a source
and the surfaces it applies to. "Applies to" references the recipe field
(`screen-composition.md` §3) and/or pattern layers. Contracts reference the IDs;
the validation checklist (`composition-validation.md`) verifies the references
are real.

---

## 1. Data surfaces (R-DATA)

| ID | Rule | Source | Applies to |
|---|---|---|---|
| R-DATA-01 | **Table vs list.** If the set is large/remote, homogeneous, and needs columnar comparison with sort/filter/pagination → **Data table**. Else if the set is ≤ 5–8 static rows or discrete units with 1–3 line hierarchy → **List**. | `composition-rules.md` §4.1, §5.1 | UX pattern E3 Resource list; SURFACE/UX_PATTERN/COMPONENTS |
| R-DATA-02 | **List vs queue.** If there is a "next to serve" with per-item operational state and an advance action → **Queue** (List two-line + Badge tone + Button + POST+303). Else it is a presentational index → **List**. | `composition-rules.md` §4.2, §5.2 | vocabulary Queue; Ops Queue recipe |
| R-DATA-03 | **Queue vs board.** If movement is strict sequential (FIFO/pending) → **Queue**. Only if there are multiple parallel lanes with free movement between states → **Board** (columns of Cards + form per item, no drag). | `composition-rules.md` §4.3, §5.2 | vocabulary Board/Queue |
| R-DATA-04 | **Card vs panel.** If it is a repeated instance of a set with its own entry/action → **Card**. If it is a single static layout region grouping page context → **Panel** (`<section>`/`<aside>`, no component of its own). | `composition-rules.md` §4.4 | vocabulary Card/Panel; KPI grids |
| R-DATA-05 | **Feed vs collection.** If the value is novelty and orientation is event/time → **Feed** (Card/List + avatar + Badge, chronological inverse, empty/loading critical). If the value is the total set with filters → **Collection** (Data table/List + filters). | `composition-rules.md` §4.5 | Public Feed recipe; Inbox thread |
| R-DATA-06 | **Timeline vs activity list.** If there is an explicit temporal axis with milestones/dates (process narration) → **Timeline** (`<ol>` + markers + dates). If it is terse recent activity → **Activity list** (List one/two-line, inverse order). | `composition-rules.md` §4.6 | Resource detail, audit trails |
| R-DATA-07 | **Dashboard.** **NEW.** If and only if the primary task is *monitoring* a set of aggregates/KPIs at a glance → **Dashboard**. Each KPI is a **Card**; each card is a real link into the underlying collection (drill-down). A dashboard is never the vehicle for transactional management (that is a table/queue), and never the default home of a workflow surface. Every KPI region declares its own loading/empty/error treatment (state matrix gap for Dashboard must be closed with real state patterns, not ad-hoc); trends are never color-only. | **NEW**; state matrix gap (`composition-rules.md` §8 Dashboard row) | Dashboard recipe (deferred, `roadmap.md`); state patterns D1–D3, D6 |

---

## 2. Overlays and surface choice (R-OVERLAY)

| ID | Rule | Source | Applies to |
|---|---|---|---|
| R-OVERLAY-01 | **Dialog vs page.** If the sub-task is short, focused and reversible with retained context (confirm, quick edit, picker) → **Dialog** (`<dialog>`). If the flow is deep/complete or needs URL + back + deep-link → **Page/Steps**. Never open a long flow in a Dialog. | `composition-rules.md` §4.7, §5.7 | Dialog component; Confirmation E18; forms |
| R-OVERLAY-02 | **Overlay fallback.** **NEW.** Every overlay (Dialog, Select menu, Drawer modal, Menu popover) MUST have a real server-rendered fallback where the native primitive is unavailable (Baseline gates). No Gelium page may leave a control inert. Evidence: gap G1 (`ux-accessibility-audit.md:58`); Dialog page variant (`dialog.html`). | **NEW** (formalizes gap G1) | Dialog, Menu, Navigation drawer |
| R-OVERLAY-03 | **Confirm consequential.** A consequential/destructive action MUST be confirmed in the native Dialog confirm variant (explicit Cancel/Confirm, `aria-labelledby`/`aria-describedby`) before the `POST + 303` commit. Reversible or low-stakes actions do not add confirmation friction. | `ux-patterns.md` E18, E10 | Destructive E10, Confirmation E18, Bulk E11 |

---

## 3. Feedback (R-FEEDBACK)

| ID | Rule | Source | Applies to |
|---|---|---|---|
| R-FEEDBACK-01 | **Persistent ≠ transient.** Persistent-contextual feedback NEVER travels via `gelium:toast`; transient action results NEVER occupy a persistent slot (Banner/Inline). | `state-patterns-audit.md:45`, `ux-patterns.md` cross-cutting | Notifications E15; all recipes |
| R-FEEDBACK-02 | **Validation never toast.** Field/form validation → `422 + X-Gelium-Validation` + Inline alert (`role="alert"`) + Validation summary; NEVER a toast. | `composition-rules.md` §5.3, §9.1; `toast.go:129-133` | Form 422; Auth/Editor/Settings |
| R-FEEDBACK-03 | **Channel selection.** Transient result of an action → **Toast** (`HX-Trigger {"gelium:toast":…}`, closed vocabulary info/success/warning/error). Persistent page/site notice requiring action → **Banner**. Persistent contextual/section or field error → **Inline alert**. Ignorable informative note → **Callout**. | `composition-rules.md` §4.8, §5.4 | Notifications E15; state patterns D3–D5, D8 |

---

## 4. Navigation (R-NAV)

| ID | Rule | Source | Applies to |
|---|---|---|---|
| R-NAV-01 | **Tabs.** **NEW.** Use **Tabs** when views compete for the same surface at the same level of the same context. Implement as `<nav>` of real links with `aria-current="page"`; the active view is server-side state. **Never** `role="tablist"`/roving focus without full keyboard. **Never** use Tabs for: cross-context navigation (use Navigation bar/drawer), sequential process steps (use Steps), or hierarchy (use Breadcrumb/nav). | **NEW**; vocabulary Tabs; `composition-rules.md` §5.6 | Public Feed views; Dashboard sections; any same-level view switcher |
| R-NAV-02 | **URL is the state.** Any navigable state is a URL: GET with stable params for list state (sort/filter/page/selection); no client-side list state. No GET for mutations. | `composition-rules.md` §9.6, §5.5 | All list/search/filter patterns |

---

## 5. Action hierarchy (R-ACTION)

| ID | Rule | Source | Applies to |
|---|---|---|---|
| R-ACTION-01 | **One primary action.** **NEW.** Exactly one primary action per surface, encoded `ui-button-primary`, placed where the user completes their primary task (submit of the main form; the main CTA of a landing; the advance button of a queue row). Secondary actions use `ui-button-secondary`/`outline`/`text`. If two candidates compete, split the surface (P7, `ux-principles.md` §8). | **NEW** (formalizes P7) | Every recipe; Hero CTA; primary submit |
| R-ACTION-02 | **Destructive is never primary.** A destructive action is signalled with the danger token and explicit verb copy; it never occupies the primary-action slot. | `ux-principles.md` §8, `content-rules.md` §2 | Destructive E10; forms with delete |

---

## 6. Required states (R-STATE)

| ID | Rule | Source | Applies to |
|---|---|---|---|
| R-STATE-01 | **Required states.** **NEW.** Every data surface declares all applicable states of the state matrix — rest, empty, loading, error (plus selected/success where applicable). Empty and loading are never optional; a GAP in the matrix is documented, never silently filled ad-hoc. | **NEW**; state matrix (`composition-rules.md` §8) | Every data pattern (E3–E9); Dashboard |
| R-STATE-02 | **Empty = message + CTA.** An empty state explains why the surface is empty and offers a real, actionable CTA (create, clear filter, retry). Never "0 rows" or bare "No data". After a filtered empty, name the active filter. | `content-rules.md` §3, `composition-rules.md` §5.9 | Empty state D1; Search E4; Filters E5 |
| R-STATE-03 | **Loading primitives.** Use `.ui-progress` (determinate/indeterminate) for operations and **Skeleton** for data regions; buttons use `aria-busy`. Never ad-hoc spinners; never flash the empty state while loading. | `composition-rules.md` §5.8; `ux-patterns.md` E8 | Loading E8; refresh flows |
| R-STATE-04 | **No color-only.** Every state carried by color is also carried by semantics or visible text (P8). Tones (Badge, Banner, Inline alert, Toast) always pair with a text label. | `accessibility-contract.md` §0.3, `composition-rules.md` §5.10 | Badge tones, status dots, trends |

---

## 7. Responsive behavior (R-RESP)

| ID | Rule | Source | Applies to |
|---|---|---|---|
| R-RESP-01 | **Fluid first.** Layouts are fluid first: `auto-fit/minmax`, `min()/clamp()` before breakpoints. Large tables resolve **server-side** via pagination, never horizontal scroll of the component. Document order = visual order. | `composition-rules.md` §7 | All recipes; landing sections |
| R-RESP-02 | **Overlays and drawer.** Overlays are fluid (`calc(100vw - n)`) and need no breakpoints; Navigation drawer resolves responsive by variants (modal vs permanent), not media queries. | `composition-rules.md` §7.4–7.5 | Dialog, Menu, Drawer |

---

## 8. Accessibility (R-A11Y)

| ID | Rule | Source | Applies to |
|---|---|---|---|
| R-A11Y-01 | **Native before ARIA.** Native semantics before ARIA; no fake roles; ARIA additive only. (P2) | `accessibility-contract.md` §0.2, §1.1 | Every component/recipe |
| R-A11Y-02 | **Accessible names.** Every control has an accessible name (text by default; `aria-label` for icon-only). Decorative SVG is `aria-hidden` + `focusable=false`. `aria-label` never overrides a visible label. | `accessibility-contract.md` §1.2 | Icon button, FAB, avatar, icons |
| R-A11Y-03 | **One h1 per page.** Exactly one `h1`; headings descend without skipping. State patterns accept a configurable heading level. | `accessibility-contract.md` §1.3 | All pages; error-state, validation-summary |
| R-A11Y-04 | **Native focus.** Focus is native (`dialog` focus trap, popover, native controls); on 422 the failing field receives focus (skipped in the HTMX branch to avoid stealing focus). | `accessibility-contract.md` §1.5; `ux-principles.md` §6 | Forms; overlays; error recovery |
| R-A11Y-05 | **Keyboard.** The full task is operable by keyboard with native contracts (Tab/Enter/Space/Escape/arrows). Roving focus only when the keyboard contract is fully solved natively. | `accessibility-contract.md` §1.5 | Tabs, Segmented, Menu, Data table |
| R-A11Y-06 | **Landmarks and order.** Landmarks (`header`/`nav`/`main`/`aside`/`footer`) match surface roles; document order equals visual order. | `composition-rules.md` §10 | All recipes |
| R-A11Y-07 | **Motion and colors.** Reduced motion and forced colors are first-class for every component that moves or relies on color. | `accessibility-contract.md` §0.6 | Skeleton shimmer, dialog transitions, tones |

---

## 9. No-JS / HTMX (R-NOJS)

| ID | Rule | Source | Applies to |
|---|---|---|---|
| R-NOJS-01 | **No-JS end-to-end.** The primary flow completes with JS disabled; HTMX only enhances. (P1) | `AI-COMPONENT-IMPLEMENTER-PROMPT.md` §12 | Every recipe |
| R-NOJS-02 | **Mutations are POST + 303.** All mutations use `POST + 303 SeeOther`; GET on a POST-only path responds `405` with `Allow: POST` (`postOnlyPaths()`). | `composition-rules.md` §9.4 | All mutating routes |
| R-NOJS-03 | **HTMX never changes the mutation contract.** HTMX enhancement swaps fragments (`hx-get`) for reading/refresh; mutations remain `POST + 303` (no `hx-post`). The refresh fragment may add `HX-Trigger gelium:toast` (transient). | `composition-rules.md` §9; recipes HTMX_ENHANCEMENT | Refresh flows |

---

## 10. Errors (R-ERROR)

| ID | Rule | Source | Applies to |
|---|---|---|---|
| R-ERROR-01 | **Validation recovery.** Field/form validation → `422 + X-Gelium-Validation: true`, per-field `aria-invalid` + `aria-describedby`, Inline alert + Validation summary with real anchor links, submitted value preserved, focus returns to the failing field (not in the HTMX branch). | `composition-rules.md` §9.1; `ux-patterns.md` E9; `text_field.go:62-67` | Forms (E1/E9/E12–E14/E16/E17) |
| R-ERROR-02 | **Resource and global errors.** Missing/invalid resource or server failure → real HTTP status (404/500/503) + `error-state` with retry GET. Persistent global notice → Banner (`role="alert"`, no auto-dismiss). Never a toast for these. | `composition-rules.md` §5.4; `ux-patterns.md` E9; recipes 404 | 404/500 handling; global banners |
| R-ERROR-03 | **Transport errors.** HTMX transport failures (network/5xx) → generic transient toast (never validation copy); no silent failure. **GAP G5** remains an open residual; until resolved, never swallow the error. | `ux-patterns.md` note (G5); `app.js` | Any remote refresh |
| R-ERROR-04 | **Recovery is explicit and lossless.** Errors state what + how to resolve (`content-rules.md` §1); retry re-submits without data loss; `POST + 303` keeps workflows resumable. | `content-rules.md` §1, `ux-patterns.md` E9 | All error paths |

---

## 11. Destructive actions (R-DESTRUCT)

| ID | Rule | Source | Applies to |
|---|---|---|---|
| R-DESTRUCT-01 | **Destructive protocol.** A destructive action uses: danger token (`--ui-color-danger`), explicit destructive verb copy (never "OK"), Dialog confirmation stating what happens and irreversibility ("This action cannot be undone."), and `POST + 303` commit. Signal is never color-only (text + confirm + focus). | `ux-principles.md` §8, `content-rules.md` §2/§5, `ux-patterns.md` E10/E18 | Delete/archive/reset anywhere |

---

## 12. SEO / GEO (R-SEO, R-GEO)

| ID | Rule | Source | Applies to |
|---|---|---|---|
| R-SEO-01 | **Indexable metadata.** Every indexable page emits per-route `<title>` + `<meta name="description">`, a clean canonical without query (`siteBaseURL + routePath`), one `h1`, correct `lang`. Demo (`/demo/*`), example (`/examples/*`) and recipe (`/recipes/*`) surfaces are `noindex, nofollow`. | `composition-rules.md` §9.5, `seo-contract.md`, `seo-patterns.md` | All recipes; landing |
| R-GEO-01 | **Entity and provenance.** Content is factual, citable, with stable deep-linkable URLs and visible provenance. The brand entity is unique (**Gelium UI**). Emit `BreadcrumbList`/Organization JSON-LD where applicable (landing, component pages). No unverified universal claims. | `geo-contract.md`, `geo-patterns.md`; `server.go` `jsonLDBreadcrumb` | Landing; public pages |

---

## 13. Content (R-CONTENT)

| ID | Rule | Source | Applies to |
|---|---|---|---|
| R-CONTENT-01 | **Buttons are action verbs.** Labels name the exact action ("Save", "Delete", "Send"); never bare "OK"/"Submit" when the action is specific. | `content-rules.md` §2 | Button, Dialog actions, CTAs |
| R-CONTENT-02 | **Errors say what + how to resolve.** Errors identify the field and the fix; never generic panic, never blame. | `content-rules.md` §1, §8 | Error state, Inline alert, Banner, Validation summary |
| R-CONTENT-03 | **Empty copy.** Empty states explain why + offer a real CTA (naming the active filter where relevant). | `content-rules.md` §3 | Empty state D1 |
| R-CONTENT-04 | **Tone.** Neutral, direct, professional; no flattery, no blame, plain language with short titles (5–9 words). | `content-rules.md` §8, §10 | All message components |
| R-CONTENT-05 | **Language.** UI strings are English by default and localizable by contract (server-rendered data); document `lang` matches the content language. | `content-rules.md` §9 | All `internal/app` strings; layout |

---

## 14. Anti-rule consolidation

The hard anti-rules of `composition-rules.md` §5 are binding in the protocol:

| Anti-rule | Binding form |
|---|---|
| §5.1 | No Data table for ≤ 5–8 static rows → List. |
| §5.2 | No Board for strict FIFO queues → Queue. |
| §5.3 | No validation errors as Toast → 422 + Inline alert. |
| §5.4 | No Toast for persistent/critical feedback → Inline alert/Banner. |
| §5.5 | No list state off the URL → GET stable params. |
| §5.6 | No `role="tablist"`/roving focus without full keyboard → links + `aria-current`. |
| §5.7 | No long flow in Dialog → Page/Steps. |
| §5.8 | No ad-hoc spinner where Progress exists → `.ui-progress` + `aria-busy`. |
| §5.9 | No empty state without message/CTA → real empty + CTA. |
| §5.10 | No color-only state → native control + forced colors. |
| §5.11 | No JS for what a GET form already solves → platform-first. |

## 15. Rule coverage vs the state matrix

The state matrix (`composition-rules.md` §8) marks GAPs for List/Data table/Queue/
Feed/Dashboard empty-loading rows. Per R-STATE-01 those GAPs are closed with the
real state patterns (D1 Empty, D2 Skeleton, D3 Inline alert, D6 Error state)
before a recipe ships — never with ad-hoc markup. The dashboard-specific
obligations live in R-DATA-07.
