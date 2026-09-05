# Gelium UI — UX Patterns

> Recurring user flows of the Gelium UI system, expressed as compositions of real components, state patterns and server contracts.
> Phase E of the system roadmap (`docs/gelium-ui-system-roadmap.md`).
> Base: `docs/gelium-ui-composition-rules.md`, `docs/gelium-ui-vocabulary.md`, `docs/gelium-ui-ux-principles.md`, `docs/handoffs/state-patterns-audit.md`, `docs/handoffs/ux-accessibility-audit.md`.

Each pattern declares: problem, user, context, happy/empty/loading/error/recovery paths, mobile behavior, accessibility, server contract, and when NOT to use it. Each is anchored to real Gelium components and marked:

- **Ready** — the underlying components/contracts exist.
- **Phase G** — the screen recipe is pending; the primitives exist but the pattern is not yet composed as a recipe.

---

## Pattern index

| # | Pattern | Status | Core primitives |
|---|---|---|---|
| 1 | Authentication | Phase G | Text field, Inline alert, Validation summary, Banner, Button, 422 contract |
| 2 | Onboarding | Phase G | Steps (gap), Callout, Text field, Banner, 422 |
| 3 | Resource list | Ready | Data table or List, Empty state, optional Skeleton for an actual waiting region, Pagination, GET params |
| 4 | Search | Ready | GET form, Empty state, Data table/List, optional Skeleton for an actual waiting region |
| 5 | Filters | Ready | Chips, Select, Segmented buttons, GET params |
| 6 | Pagination | Ready | Data table pagination `<nav>` links, GET params |
| 7 | Empty state | Ready | `empty-state.html` |
| 8 | Loading | Ready | Button `aria-busy`, `<progress>`, Skeleton |
| 9 | Error recovery | Ready* | Error state, Inline alert, Validation summary, 422 |
| 10 | Destructive action | Ready | Button (danger), Dialog confirm, POST + 303 |
| 11 | Bulk action | Ready* | Data table checkboxes, Dialog confirm, POST + 303 |
| 12 | Multi-step form | Phase G | Steps (gap), Progress, Validation summary, 422 |
| 13 | Checkout | Phase G | Steps (gap), Text field, 422, Validation summary |
| 14 | Booking | Phase G | Steps (gap), date/datetime input, 422, Success feedback |
| 15 | Notifications | Ready | Toast (`gelium:toast`), Banner, Toast region `aria-live` |
| 16 | Settings | Phase G | Panel, List, Text field, Banner/Inline alert success, 422 |
| 17 | Permissions | Phase G | List + checkboxes, Segmented buttons, Dialog, 422 |
| 18 | Confirmation | Ready | Dialog confirm/cancel, `closedby="any"` |
| 19 | Undo / recovery | Phase G | Toast action (transient), 422 value preservation, POST + 303 |

`*` Ready at the primitive level; the screen recipe lives in Phase G. Transport-error feedback (500/network in HTMX) is still the Phase E gap G5 (`ux-accessibility-audit.md:88`) — it applies to every pattern that uses remote refresh.

---

## 1. Authentication (Phase G)

- **Problem**: prove identity before access; recover from failed attempts without lock-out confusion.
- **User**: any user on a protected surface (public portal, admin, ops).
- **Context**: a full screen — never a dialog; auth is deep and resumable (`composition-rules.md:100-109`).
- **Happy path**: submit credentials → server validates → `POST + 303` redirect to the protected home.
- **Empty path**: first visit — empty form with helper text; no pre-filled credentials.
- **Loading path**: submit button `aria-busy` + spinner (`button.html:4,9`).
- **Error path**: wrong credentials → **422 + `X-Gelium-Validation`** → Inline alert (`role="alert"`) + per-field `aria-invalid` + Validation summary with anchor links.
- **Recovery path**: server preserves the entered value (`text_field.go:62`); focus returns to the failing field (`text_field.go:67`); "Reset password" and "Try again" as real links. Global/session-level errors → Banner.
- **Mobile behavior**: single-column form, full-width primary button; no side-by-side fields.
- **Accessibility**: `role="alert"` on the error summary; never announce credentials via `aria-live`; autofocus to the first error field.
- **Server contract**: `POST` + `422 + X-Gelium-Validation: true` on failure; `POST + 303 SeeOther` on success. Validation never triggers a toast (`toast.go:129-133`).
- **When NOT to use**: a transient "login as" inside a demo, or any surface where auth is not a distinct navigable state.

---

## 2. Onboarding (Phase G)

- **Problem**: guide a new user through setup steps without overwhelming them.
- **User**: first-time users of the consuming project.
- **Context**: a multi-step flow — page or Steps, never a Dialog (`composition-rules.md:129`).
- **Happy path**: complete step N → validated via 422 → advance → `POST + 303` to the next step.
- **Empty path**: fresh start — Callout with contextual explanation + a clear first field.
- **Loading path**: step submit button `aria-busy`; Progress determinate for position (`progress.html`).
- **Error path**: per-step validation → Inline alert + Validation summary links to the failing field.
- **Recovery path**: back navigation is server-rendered links; values survive re-submit; no data loss.
- **Mobile behavior**: one field/decision per screen; large touch targets.
- **Accessibility**: each step has an `<h1>`; error summary is `role="alert"`; skip-link to content.
- **Server contract**: `POST + 303` per step; `422 + X-Gelium-Validation` per step; **Steps is a Phase D gap** (`vocabulary.md:190-197`).
- **When NOT to use**: flows that fit on one page, or when the user can opt out without penalty.

---

## 3. Resource list (Ready)

- **Problem**: present a set of records for scanning, comparison and per-row action.
- **User**: admin/ops user managing records.
- **Context**: Admin Resource recipe surface (Phase G) — Data table or List.
- **Happy path**: GET with stable params → server renders the table/list → row action or selection → POST + 303.
- **Empty path**: `empty-state.html` row with message + CTA (`data-table.html:68-70`); "Select all" is hidden when empty (`data-table.html:42`).
- **Loading path**: server-rendered first paint does not require a client-like Skeleton; use Skeleton only when a real data region is rendered in a waiting state. A remote refresh may use Progress and transient feedback (`data-table.html:81-91`).
- **Error path**: resource error → `error-state.html` (page) or Inline alert (fragment).
- **Recovery path**: URL carries `?q=&sort=&dir=&page=&selection=`; back/refresh restores state; re-submit re-renders.
- **Mobile behavior**: server-side pagination instead of horizontal scroll (`composition-rules.md:145`); List over Data table when narrow.
- **Accessibility**: `aria-sort` on the active column, `aria-current="page"` on pagination, native checkboxes.
- **Server contract**: `GET ?q=&sort=&dir=&page=&selection=` (`composition-rules.md:170`); `HX-Request` bifurcates fragment vs page.
- **When NOT to use**: ≤5-8 static rows → List (`composition-rules.md:123`).

---

## 4. Search (Ready)

- **Problem**: find an entity by free text across a set.
- **User**: any user on a collection surface.
- **Context**: a GET form above a Data table/List (Search Results recipe, Phase G).
- **Happy path**: submit GET `?q=…` → server-rendered results.
- **Empty path**: `empty-state.html` — "No results for `{query}`. Check the spelling or clear the filter." + CTA (clear filter link).
- **Loading path**: when an enhanced request leaves the results region waiting, Skeleton may preserve its structure; a server-rendered first response does not need a fabricated loading phase. An `aria-live` region announces an actual HTMX swap.
- **Error path**: malformed/vocabulary-invalid query → sanitized, never crashes; server renders normal results or empty.
- **Recovery path**: query persists in the URL (`?q=`); back button returns to results; clear-filter link.
- **Mobile behavior**: search input full-width; `type="search"` for native clear.
- **Accessibility**: visible label ("Filter") and `aria-label` must match (audit G6: the visible label must not be overridden); results region `role="status"`/`aria-live="polite"`.
- **Server contract**: `GET ?q=` (stable param); no client-side filtering (`composition-rules.md:127`).
- **When NOT to use**: tiny static sets where scanning beats searching.

---

## 5. Filters (Ready)

- **Problem**: narrow a set by categories/attributes.
- **User**: admin/ops user on a collection.
- **Context**: a GET form with Chips / Select / Segmented buttons above the set.
- **Happy path**: pick a value → GET with a stable param → filtered server render.
- **Empty path**: no matching records → empty state naming the active filters + "Clear filters" CTA.
- **Loading path**: Skeleton on the list region; active filter state persists in the URL.
- **Error path**: unknown vocabulary → sanitized/ignored server-side; never an error page.
- **Recovery path**: the URL encodes every filter; deep-linkable and shareable.
- **Mobile behavior**: horizontal chip scroll or a collapsible filter panel; never a hover-only filter.
- **Accessibility**: chips/selects are native controls; the active filter is a visible label + link, not color-only.
- **Server contract**: `GET` with stable params; closed vocabularies sanitized (`composition-rules.md:170`).
- **When NOT to use**: filtering a set that is already a URL segment (dedicated page), or when the set has no meaningful attributes.

---

## 6. Pagination (Ready)

- **Problem**: navigate pages of a server-side set.
- **User**: any user on a paginated collection.
- **Context**: `<nav>` below the Data table (`data-table.html:74-78`).
- **Happy path**: click a page link → GET `?page=` → re-render.
- **Empty path**: single page, no controls needed; page 1 with no rows → empty state.
- **Loading path**: HTMX swap replaces the panel; Skeleton optional.
- **Error path**: page out of range → server clamps to a valid range (200), never a broken link.
- **Recovery path**: page + sort + filter + q persist together in the URL.
- **Mobile behavior**: Previous/Next links prominent; page numbers may collapse on narrow surfaces.
- **Accessibility**: `<nav aria-label="Table pages">`, current page `<span aria-current="page">`, real links.
- **Server contract**: `GET ?page=` stable param; `HX-Request` fragment swap.
- **When NOT to use**: ≤1 page of data; or when a "load more" Feed (Phase F) adds real value.

---

## 7. Empty state (Ready)

- **Problem**: a data surface with no records must guide, not look broken.
- **User**: any user on a collection/feed.
- **Context**: `empty-state.html` rendered inside the data region (Data table row `<td colspan>`, List, Feed).
- **Happy path**: a message explains the absence + an optional CTA (create item, clear filter).
- **Empty path**: the state itself — title + body + CTA; compact variant for embedded regions.
- **Loading path**: never flash empty while loading → render Skeleton first, then the real empty.
- **Error path**: an empty result must not be confused with an error; errors use `error-state`, not `empty-state`.
- **Recovery path**: the CTA is a real link (`ui-button` with `href`), e.g. "Clear filter", "Add first record".
- **Mobile behavior**: compact variant fits embedded cards; full-width CTA optional.
- **Accessibility**: `role="status"` (polite); the message is text, never color-only.
- **Server contract**: empty is **server output**, never client state (`state-patterns-audit.md:71`); HTMX swaps the same region.
- **When NOT to use**: an error (use error-state), or a transient moment before data arrives (use Skeleton).

---

## 8. Loading (Ready)

- **Problem**: communicate that work is in progress without blocking.
- **User**: any user triggering an operation or waiting for data.
- **Context**: button/operation (transient) vs data region (persistent placeholder).
- **Happy path**: operation completes → button clears → result/toast.
- **Empty path**: n/a (loading is a state, not a surface).
- **Loading path**: submit button `aria-busy` + spinner + sr-only "Loading {Label}" (`button.html:4,9`); `<progress>` determinate/indeterminate for operations (`progress.html`); Skeleton for data regions.
- **Error path**: loading fails → transport feedback (Phase E G5) or error state; button recovers to enabled.
- **Recovery path**: retry via a real button/link; server re-render replaces the skeleton.
- **Mobile behavior**: spinners and progress must not cause layout shift (reserve space).
- **Accessibility**: `aria-busy` + `role="status"`; never announce purely visual spinners as content.
- **Server contract**: loading is declarative HTML/CSS; no server contract beyond the subsequent request.
- **When NOT to use**: ad-hoc CSS spinners when `<progress>` exists (`composition-rules.md:130`); indefinite spinners for a process that should show determinate progress.

---

## 9. Error recovery (Ready*)

- **Problem**: after a failure, the user understands why and can retry without redoing work.
- **User**: any user submitting or refreshing.
- **Context**: forms (validation), resources (404/500), fragments (HTMX).
- **Happy path**: submit succeeds → success feedback / redirect.
- **Empty path**: n/a.
- **Loading path**: button `aria-busy` during retry.
- **Error path**: validation → `422 + X-Gelium-Validation` + Inline alert + Validation summary; resource → `error-state.html` with status + retry link.
- **Recovery path**: value preserved (`text_field.go:62`), focus returns to the failing field (`text_field.go:67`), summary links jump to fields, retry link re-GETs.
- **Mobile behavior**: summary links are full-width tap targets on narrow surfaces.
- **Accessibility**: `role="alert"` for errors; focus returns to the first error; transport errors announced via `aria-live` (G5 pending).
- **Server contract**: `422 + X-Gelium-Validation: true`; real HTTP status (404/500/503); `POST + 303` for resumable workflows.
- **When NOT to use**: an empty result (use empty state); a transient action result (use Toast).

---

## 10. Destructive action (Ready)

- **Problem**: irreversible or high-impact actions (delete, reset) need explicit confirmation and clear intent.
- **User**: admin/ops user.
- **Context**: destructive Button (danger token) + Dialog confirm (`dialog.html`) + `POST + 303`.
- **Happy path**: user confirms → POST → 303 redirect → success feedback (persistent) or Toast.
- **Empty path**: n/a.
- **Loading path**: confirm button `aria-busy` while the request runs.
- **Error path**: failure → Inline alert/Banner (not Toast); validation is never used here.
- **Recovery path**: Cancel in the Dialog aborts; POST + 303 makes the result a new deep-linkable state; undo (see #19) where feasible.
- **Mobile behavior**: Dialog is fluid (`calc(100vw - n)`); danger button full-width optional.
- **Accessibility**: `<dialog closedby="any">` native focus trap; Cancel carries `autofocus` (`dialog.go:18`); destructive intent is text + button, never color-only; server-rendered dialog fallback for non-Baseline browsers (G1).
- **Server contract**: `POST` + `303 SeeOther` (`composition-rules.md:171`); no GET for mutations.
- **When NOT to use**: reversible actions don't need confirmation; low-risk actions shouldn't add friction.

---

## 11. Bulk action (Ready*)

- **Problem**: act on multiple records at once (archive, assign, delete).
- **User**: admin/ops user on a Data table/List.
- **Context**: native row checkboxes + selection state + one primary action.
- **Happy path**: select rows → submit selection → POST → 303 → success feedback.
- **Empty path**: selection hidden when no rows (`data-table.html:42`); empty state with CTA instead.
- **Loading path**: action button `aria-busy`.
- **Error path**: partial failure → error feedback naming what failed; never a silent partial success.
- **Recovery path**: selection persists in `?selection=`; re-submit after failure.
- **Mobile behavior**: List with checkboxes over Data table on narrow surfaces.
- **Accessibility**: native checkboxes with `aria-label` per row; selected rows derive state from `:checked`, not color.
- **Server contract**: `GET ?selection=` for state; `POST + 303` for the mutation.
- **When NOT to use**: single-record actions; bulk selection spanning >1 page without clear semantics.

---

## 12. Multi-step form (Phase G)

- **Problem**: collect complex input across logical steps without overwhelming the user.
- **User**: users creating/editing complex resources.
- **Context**: Steps pattern (`vocabulary.md:190-197`, Phase D gap) — page-based, never Dialog.
- **Happy path**: validate a step via 422 → advance → complete → POST + 303.
- **Empty path**: first step with helper text; no invented defaults.
- **Loading path**: step submit `aria-busy`; Progress determinate for position.
- **Error path**: per-step 422 + Inline alert + Validation summary.
- **Recovery path**: back/forward server-rendered; values preserved on validation failure.
- **Mobile behavior**: one field-group per step; sticky step submit.
- **Accessibility**: `<h1>` per step; summary `role="alert"`; focus to the first error field.
- **Server contract**: `422 + X-Gelium-Validation` per step; `POST + 303` between steps.
- **When NOT to use**: a flow that fits on one page (`vocabulary.md:197`).

---

## 13. Checkout (Phase G)

- **Problem**: complete a purchase reliably; reassure the user at each stage.
- **User**: consumer on a public surface.
- **Context**: multi-step flow (Steps) with payments handled by the consuming backend.
- **Happy path**: steps validated → confirm → POST + 303 → success page with persistent confirmation.
- **Empty path**: empty cart → empty state + CTA ("Continue browsing").
- **Loading path**: `aria-busy` on payment/confirm; Progress determinate.
- **Error path**: card/validation errors → 422 + Inline alert per step; never a toast for validation.
- **Recovery path**: values preserved; clear "Back" and order-summary consistency.
- **Mobile behavior**: payment fields collapse to a single column; security badges visible but never decorative.
- **Accessibility**: logical tab order through payment fields; errors `role="alert"`.
- **Server contract**: `POST + 303` per step; `422 + X-Gelium-Validation`; success → persistent Success feedback (Banner/inline, never only Toast).
- **When NOT to use**: single-screen purchases (one page is enough).

---

## 14. Booking (Phase G)

- **Problem**: select a date/slot/option and confirm a reservation.
- **User**: consumers on a public surface.
- **Context**: Steps + date/datetime inputs (`vocabulary.md:295-297` — native input candidate, calendar component pending).
- **Happy path**: pick slot → confirm → POST + 303 → booking confirmation with persistent success.
- **Empty path**: no available slots → empty state + alternative CTA (notify me, another date).
- **Loading path**: `aria-busy` on confirm; Progress determinate through steps.
- **Error path**: 422 + Inline alert on the failing field; summary for multi-field steps.
- **Recovery path**: values preserved on failure; resumable by URL/redirect.
- **Mobile behavior**: slot pickers collapse to one column; native `input type="date"` preferred.
- **Accessibility**: native inputs for keyboard/calendar semantics; error summary `role="alert"`.
- **Server contract**: `POST + 303` on success; `422 + X-Gelium-Validation` on failure.
- **When NOT to use**: when a single GET-parameter slot picker (e.g. a queue) suffices; or when there is no time dimension.

---

## 15. Notifications (Ready)

- **Problem**: surface transient results and persistent site-level announcements without mixing them.
- **User**: any user.
- **Context**: transient → Toast (`gelium:toast`); persistent/global → Banner; contextual → Inline alert / Callout.
- **Happy path**: action completes → server emits `gelium:toast` → auto-dismiss 4s/8s, pausable (`app.js:11-77`).
- **Empty path**: no notifications; nothing rendered.
- **Loading path**: n/a.
- **Error path**: global error/session expiry → Banner `role="alert"` (persistent, no auto-dismiss); never a toast for persistent/critical feedback (`composition-rules.md:126`).
- **Recovery path**: Banner dismiss = `POST + 303` (`banner.html`); toast dismiss manual + auto.
- **Mobile behavior**: toast region stacked bottom on narrow surfaces; Banner full-width top.
- **Accessibility**: `#gelium-toast-region` `aria-live="polite"` `aria-atomic="false"` (`toast.html:10`); `role="alert"` for error toasts.
- **Server contract**: `HX-Trigger: {"gelium:toast":{...}}` with closed vocabulary `info|success|warning|error` (`toast.go:13-14,45`); validation never toast (`toast.go:129-133`).
- **When NOT to use**: persistent feedback (Inline alert/Banner instead); validation errors (422 inline instead).

---

## 16. Settings (Phase G)

- **Problem**: let users view and change configuration without losing context or data.
- **User**: admin/user on the account surface.
- **Context**: Panel-based page (Settings recipe, Phase G) — Panels group related fields.
- **Happy path**: edit a field → submit → 422 if invalid → POST + 303 → persistent success banner/inline on the saved section.
- **Empty path**: new/empty config → sensible defaults + helper text.
- **Loading path**: submit `aria-busy`; sections render server-side.
- **Error path**: validation → 422 + Inline alert + Validation summary; never a toast (`composition-rules.md:126`).
- **Recovery path**: values preserved (`text_field.go:62`); success is a non-ephemeral Banner/Inline `role="status"` (Success feedback, `vocabulary.md:140-151`).
- **Mobile behavior**: sections stack; save action sticky.
- **Accessibility**: grouped fields in `<fieldset>`/`<legend>`; success `role="status"` polite.
- **Server contract**: `POST + 303` after save; `422 + X-Gelium-Validation` on validation; persistent success never via `gelium:toast`.
- **When NOT to use**: quick single toggles that belong inline in the current page (Segmented/Checkbox + POST).

---

## 17. Permissions (Phase G)

- **Problem**: show and change who can do what, without dangerous defaults.
- **User**: admin user.
- **Context**: List + native checkboxes/radios/Segmented buttons, grouped by resource (`composition-rules.md:124`).
- **Happy path**: change a permission → POST → 303 → persistent success.
- **Empty path**: no permissions yet → default deny + helper text.
- **Loading path**: submit `aria-busy`.
- **Error path**: 422 + Inline alert on the failing group; confirmation dialog for revoking own access / removing roles.
- **Recovery path**: change is a URL/redirect state; values preserved on failure.
- **Mobile behavior**: one permission per row on narrow surfaces.
- **Accessibility**: native controls (`checkbox.html`, `radio.html`, `segmented-button.html`); groups in `<fieldset>`.
- **Server contract**: `POST + 303`; `422 + X-Gelium-Validation`; destructive grants use Dialog confirm (see #18).
- **When NOT to use**: when permissions come from an external identity provider and are read-only.

---

## 18. Confirmation (Ready)

- **Problem**: confirm an action before committing, especially when it is hard to reverse.
- **User**: any user triggering a consequential action.
- **Context**: `<dialog closedby="any">` with Cancel (`autofocus`) and Confirm (`dialog.html`).
- **Happy path**: confirm → POST + 303 → result state.
- **Empty path**: n/a.
- **Loading path**: confirm button `aria-busy`.
- **Error path**: POST failure → Inline alert/Banner (never toast for the failure itself).
- **Recovery path**: Cancel/light-dismiss/Escape abort; POST + 303 keeps the flow resumable.
- **Mobile behavior**: fluid dialog (`dialog.css:3-5`); actions full-width on narrow surfaces.
- **Accessibility**: `aria-labelledby` headline + `aria-describedby` content; focus trap native; Cancel autofocused; server-rendered fallback for non-Baseline browsers (G1).
- **Server contract**: `POST + 303` for the confirmed action; no GET mutation.
- **When NOT to use**: reversible actions, or low-stakes actions where confirmation adds friction.

---

## 19. Undo / recovery (Phase G)

- **Problem**: let users reverse a mistake or resume after interruption.
- **User**: any user.
- **Context**: transient undo affordance (Toast action) for small mutations; POST + 303 resumability for workflows; 422 value preservation for forms.
- **Happy path**: act → toast with "Undo" action → undo re-POSTs the previous state.
- **Empty path**: n/a.
- **Loading path**: undo button `aria-busy`.
- **Error path**: undo failed → inline/Banner error; re-submit.
- **Recovery path**: server-side undo endpoint (POST + 303); form values preserved on validation (`text_field.go:62`); navigation state in URL.
- **Mobile behavior**: toast with action must be reachable (pause auto-dismiss while hovered/focused, `app.js`).
- **Accessibility**: the undo action is a real button; announcement via `role="status"`.
- **Server contract**: POST + 303 for both the action and the undo; no new contract invented. **Undo is not yet a canonical Gelium contract** — Phase G must define it before recipes rely on it.
- **When NOT to use**: destructive actions with explicit confirmation and no practical reversal window; use re-run/retry instead of undo for non-deterministic operations.

---

## Status vs server contracts

| Contract | Patterns that own it |
|---|---|
| `422 + X-Gelium-Validation: true` | Authentication, Onboarding, Multi-step form, Checkout, Booking, Settings, Permissions, Error recovery |
| `HX-Trigger {"gelium:toast":…}` | Notifications, Undo (transient), action results in Resource list / Bulk / Destructive |
| `GET` stable params (`?q=&sort=&dir=&page=&selection=`) | Resource list, Search, Filters, Pagination |
| `POST + 303 SeeOther` | Authentication, Onboarding, Checkout, Booking, Destructive, Bulk, Confirmation, Settings, Permissions, Undo |
| Persistent success (Banner/Inline `role="status"`) | Checkout, Booking, Settings, Destructive (post-redirect) |

**Cross-cutting rule**: persistent feedback never travels through `gelium:toast`; transient action results never occupy a persistent slot (`state-patterns-audit.md:45`).

---

**Definition of done (Phase E scope for this doc)**: patterns anchored to real components/contracts, readiness marked, cross-referenced from `composition-rules.md` and referenced by the Phase G screen recipes.
