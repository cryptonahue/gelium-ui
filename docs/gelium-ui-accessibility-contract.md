# Gelium UI — Accessibility Contract

> Contract of the accessibility obligations of the Gelium UI system. Phase E of the system roadmap (`docs/gelium-ui-system-roadmap.md`).
> This document is the single source of truth for how the system satisfies WCAG 2.1 AA and its own accessibility invariants. It complements `docs/gelium-ui-core.md` (§6 HTML-first, §9 Accesibilidad), `docs/gelium-ui-composition-rules.md` and `docs/gelium-ui-vocabulary.md`.
> Evidence baseline: `docs/handoffs/ux-accessibility-audit.md`, `docs/handoffs/state-patterns-audit.md`, `web/templates/*.html`, `web/styles/*.css`, `internal/app/*.go`, `web/static/app.js`.

---

## 0. Non-negotiable invariants

1. **HTML-first, server-first, no-JS end-to-end**: the primary flow of every component and screen works with no JavaScript. ARIA only completes what native semantics cannot express.
2. **Native before ARIA**: `button`, `a`, `dialog`, `progress`, `select`, `table/caption`, native `input` types, `nav`, `aside`, `label`, `fieldset/legend` are preferred over invented widgets. No fake roles: no `tablist`, no fake `menu`, no fake `dialog`.
3. **State is never color-only**: every state ported by color (`:checked`, `aria-sort`, `aria-current`, `aria-invalid`, `aria-busy`, `role`, disabled) is also carried by semantics or text.
4. **Persistent-contextual feedback ≠ transient-action feedback**: nothing persistent is announced through `gelium:toast`; nothing transient occupies a persistent slot. Formalized in Phase D and enforced by tests.
5. **Focus is native**: focus management is delegated to the browser (dialog `showModal`, popover, native controls). Roving focus is only acceptable when the full keyboard contract is already solved natively — and even then links with `aria-current` are preferred.
6. **Reduced motion and forced colors are first-class**: every component that moves or relies on color has a `prefers-reduced-motion` and `forced-colors` treatment.
7. **AA contrast by design**: semantic tokens, not hardcoded values, drive the palette; the theme is responsible for AA ratios in both light and dark.

---

## 1. Contract items

### 1.1 Semantic HTML

**Architecture rule (mandatory)**
Components are delivered as native semantic markup server-rendered. ARIA is additive only when the native element cannot express the required semantics; no `div`/`span` acts as a control; no redundant ARIA on native elements.

**Canonical roots**
| Pattern | Native root | Notes |
|---|---|---|
| Button / CTA | `<button>` or `<a>` | Never `div[role=button]`; icon-only via `aria-label` |
| Link / tab | `<a>` + `aria-current` | Tabs are links, never `tablist` |
| Menu | `<ul>` of native controls | Items are real `button`/`a`/`checkbox`/`radio`; no `role="menu"` |
| Dialog / Select menu | `<dialog>` | Native modal, focus trap, Escape |
| Select | `<select>` | Native options |
| Data table | `<table>` + `<caption>` + `th scope` | Sort/filter/pagination are real links |
| Progress | `<progress>` | Determinate/indeterminate native |
| Segmented control | `fieldset` + radios/checkboxes | Arrows work natively |
| Callout | `<aside>` | Complementary content |
| Banner / inline alert | `<div>` + derived role | See §2 |
| Form | `<form>` + `<label>` + `fieldset/legend` | Server-driven |

**Application in Gelium UI**
- `web/templates/*.html` implement the roots above; `menu.html` uses `ul.ui-menu` with native item controls; `tabs.html` renders `<a>` tabs with `aria-current`; `dialog.html` uses `<dialog closedby="any">`; `select.html` uses native `<select>` for both the plain select and the select-menu demo (the G1 fix replaced the `<dialog>` menu with the component's own native `<select>`).
- `role="link"` on disabled `<a>` is intentional in the disabled-link pattern (`button.html:4`); `role="separator"` on `li.ui-menu-divider` is valid.

**Evidence / tests**
- `web/styles_contract_test.go` (class vocabulary closed, component CSS wired into `app.css`).
- Per-component class-vocabulary tests (e.g. `styles_menu_test.go`, `styles_tabs_test.go`, `styles_dialog_test.go`).
- `internal/app/server_test.go` (`TestLayoutRendersThemeClassOnHTMLElement`, error-state single `h1`).

---

### 1.2 Accessible names

**Architecture rule (mandatory)**
Every control must have an accessible name. Text content is the default source; `aria-label`/`aria-labelledby` are allowed for icon-only controls; decorative SVG must be `aria-hidden="true"` + `focusable="false"` and must never be the name source.

**Application in Gelium UI**
- Icon-only controls (icon button, FAB, dismiss buttons) carry `aria-label` (`banner.html:9` "Dismiss", `data-table.html:42,59`, toast dismiss `toast.html:5`).
- Decorative icons: `.ui-icon`, spinner, caret, error icon all use `aria-hidden="true"` (+ `focusable="false"` where they are SVGs).
- Loading buttons keep a real name: `sr-only` "Loading {Label}" + visible label hidden from AT (`button.html:4,9`).
- CTA links use their visible text; icon + text buttons use text.
- **Rule against overriding visible labels**: `aria-label` must not shadow a visible `<label>`/caption (G6 — the demo overrides were removed; components never reproduce the pattern).

**Evidence / tests**
- `web/styles_icon_button_test.go`, `web/styles_fab_test.go`, `web/styles_icon_test.go` (icon SVG `aria-hidden`/`focusable` contract).
- `web/styles_button_test.go` (spinner `aria-hidden`, `sr-only` loading name).
- `web/styles_banner_test.go` (dismiss `aria-label`).

---

### 1.3 Heading hierarchy

**Architecture rule (mandatory)**
Exactly one `h1` per page. The layout shell must not declare its own `h1`; the `h1` is owned by the page content (server-rendered markdown heading or the error-state component). Headings descend without skipping levels where practical; state patterns accept a configurable heading level (`h{{.HeadingLevel}}`).

**Application in Gelium UI**
- `layout.html` contains no `h1`; the `<article class="prose">` renders the page's markdown with its single `h1`; `error-state.html:3` renders `<h1 class="ui-error-state-title">`.
- `validation-summary.html:2` accepts `HeadingLevel` so it can fit any page without breaking the sequence.
- `callout.html:3` uses an optional `<h3>`; `empty-state.html` uses `<p>` titles (no heading pressure) — a page may combine one `h1` with `h2`/`h3` subsections.

**Evidence / tests**
- `internal/app/server_test.go` (`TestErrorStateMarkupContracts`, `TestUnknownRouteRendersErrorStatePage` assert exactly one `h1`; `TestHomeRendersMarkdownInsideDogfoodedLayout` asserts `<h1>` from content).
- `web/styles_error_state_test.go`.

---

### 1.4 Labels

**Architecture rule (mandatory)**
Form controls are labeled with a native `<label for="…">`. Help text and error text are associated with `aria-describedby`. Validation state is exposed with `aria-invalid="true"` on the control. The error/helper message is a real `<p>` referenced by `aria-describedby`, never a title/tooltip-only.

**Application in Gelium UI**
- `text-field.html:4-5`: `<label for="{{.ID}}">`; on error the control gets `aria-invalid="true"` + `aria-describedby="{{.ID}}-error"`; otherwise helper via `aria-describedby="{{.ID}}-help"`.
- Error message `<p id="{{.ID}}-error" role="alert">`; helper `<p id="{{.ID}}-help">` with optional `role="status"`.
- Data table uses `aria-label` on checkboxes (there is no visible label per checkbox) — acceptable for icon-like controls inside cells.
- The filter input in the demo is named by its visible "Filter" label; the previous `aria-label` override was removed as part of G6, so the visible label is the single name source.

**Evidence / tests**
- `web/styles_text_field_test.go` (`aria-invalid`, `aria-describedby` wiring; forced-colors input contract).
- `internal/app/server_test.go` validation flow; `internal/app/text_field_test.go`.

---

### 1.5 Keyboard

**Architecture rule (mandatory)**
The full task must be operable by keyboard using native contracts: `Tab` for navigation, `Enter`/`Space` to activate, `Escape`/light-dismiss to close overlays, native arrows for radios/selects/range. Roving focus is permitted only when the complete keyboard contract is already resolved natively — in that case prefer real links with `aria-current` over a custom `tablist`/`tabpanel` widget.

**Application in Gelium UI**
- Tabs: real links + `aria-current="page"` (`tabs.html:10`); no arrows, no `tablist` — a documented decision; link keyboard (Tab/Enter) is sufficient and expected.
- Dialog: native `closedby="any"` + Escape + focus trap (`dialog.html:3`).
- Select menu: native `<select>` (the G1 fix) — arrow keys, type-ahead and Escape are native; the docs page demo posts the value server-side (`select.html:77-81`).
- Menu popover: native Popover API — Escape + light dismiss; items are real controls; Tab navigation (no roving focus, because no `role="menu"`).
- Segmented buttons, radio/checkbox/switch/slider/select: native arrow behavior via real inputs in `fieldset`.
- Data table: sort/filter/pagination are real links; checkboxes native.
- Focusable-but-not-activatable patterns (e.g. a `span[tabindex=0]` without an Enter handler) are forbidden — `focus-ring.html` demo of that anti-pattern is documentation of what NOT to build.

**Evidence / tests**
- `web/styles_tabs_test.go` (`aria-current`, focus-visible), `web/styles_menu_test.go`, `web/styles_select_menu_test.go`, `web/styles_dialog_test.go`, `web/styles_segmented_button_test.go`.

---

### 1.6 Focus visible

**Architecture rule (mandatory)**
Every focusable element must show a visible focus indicator. The system uses `:focus-visible` (keyboard-initiated focus only) with the `--ui-focus-*` token family; in forced colors the indicator switches to `Highlight`. No element may lose the indicator; no focus indicator may be removed globally.

**Application in Gelium UI**
- Global `:focus-visible` rules in `focus-ring.css` + `app.css`; per-component rules (`tabs.css` indicator, `toast.css` action ring, `text-field.css` control outline).
- Tokens: `--ui-focus-thickness: 3px`, `--ui-focus-offset: 2px`, `--ui-color-focus-ring` (`tokens.css:63-64`).
- Forced colors: indicator uses `outline-color: Highlight` per component (see §1.15).

**Evidence / tests**
- `web/styles_focus_ring_test.go`, and `:focus-visible` assertions across `styles_*_test.go` (button, text-field, tabs, toast, tooltip, validation-summary, switch, etc.).
- `web/styles_contract_test.go` (focus family presence in compiled CSS).

---

### 1.7 Focus initial (422 recovery)

**Architecture rule (mandatory)**
On server-side validation failure (HTTP 422), focus must be placed back on the first invalid control so the user can recover. This applies to the **no-JS** full-page re-render. In the HTMX branch, `autofocus` is deliberately omitted so the enhanced flow does not steal focus.

**Application in Gelium UI**
- `text_field.go:65-68`: on empty value the field is marked with `Error` and `Autofocus = !isHX`; the re-rendered `<input>`/`<textarea>` carries `autofocus` (`text-field.html:5`).
- `internal/app/server.go` 422 status + `X-Gelium-Validation: true` header in HX branch (`text_field.go:87-89`); `app.js:1-9` swaps only 422 with that header.
- The error message (`role="alert"`) is right below the field, so the recovered focus lands adjacent to the announcement.

**Evidence / tests**
- `internal/app/text_field_test.go` (422 status, autofocus present on re-render, absent in HX branch).
- `internal/app/server_test.go` validation demo route.

---

### 1.8 Focus restoration

**Architecture rule (mandatory)**
When a modal overlay closes, focus must return to the element that opened it. The system relies on native `<dialog>` `showModal` behavior for trap + restoration; any future JS-driven overlay must reproduce trap and restoration or be rejected.

**Application in Gelium UI**
- Dialog: native `showModal` restores focus to the trigger on close; the trigger is the `command="show-modal"`/`commandfor` button (`dialog.go:17`, `button.html:9`).
- Menu popover: native Popover API handles focus in/out.
- **G1 resolved**: the base flow is a server-rendered page — the dialog trigger is a real link to `/components/dialog/confirm` and the select-menu demo is a native `<select>` (commit `86a7f71`); "restoration" for the page variant is the Cancel link back. The `<dialog>` + Invoker Commands remains as an opt-in progressive enhancement only.

**Evidence / tests**
- `web/styles_dialog_test.go`, `web/styles_menu_test.go` (popover contract).
- Gap tracked as G1 (see §3).

---

### 1.9 Escape / cancel dialog behavior

**Architecture rule (mandatory)**
Dialogs must close on `Escape` and on cancel, with both the "close/cancel" command and the top layer behavior natively connected. The dialog must not require JS to close.

**Application in Gelium UI**
- `<dialog closedby="any">` — Escape, cancel, light-dismiss all close natively (`dialog.html:3`).
- Cancel button carries `autofocus` + `command="request-close"` (`dialog.go:18`); Confirm uses `command="close"`.
- Invoker Commands (`command`/`commandfor`) are the no-JS opening mechanism in supporting browsers.
- **G1 resolved**: the base dialog flow is the server-rendered confirm page (`/components/dialog/confirm`) and the select-menu demo is a native `<select>` (commit `86a7f71`); the modal keeps the Escape/cancel contract where it exists.

**Evidence / tests**
- `web/styles_dialog_test.go`, `web/styles_select_menu_test.go` (closedby contract), `web/styles_button_test.go` (command attributes preserved).
- Gap tracked as G1 (see §3).

---

### 1.10 Disabled / unavailable

**Architecture rule (mandatory)**
Unavailable controls must be visibly distinct (never color-only), unreachable or explicitly skipped by assistive tech, and not focusable.

**Application in Gelium UI**
- Real `<button disabled>` for buttons (`button.html:9`); no interaction, natively excluded from tab order and AT.
- Disabled links: the `<a>` variant renders with `role="link" aria-disabled="true" tabindex="-1"` (`button.html:4`) so it is announced but not focusable.
- Disabled opacity token: `--ui-state-disabled-opacity: .38` (`tokens.css:60`); applied with the disabled state — color plus semantics, not color-only.
- Pagination: disabled page links render as `<span aria-disabled="true">` (`data-table.html:75,77`).
- Text fields/selects use the native `disabled` attribute (`text-field.html:5`, `select.html:24`).

**Evidence / tests**
- `web/styles_button_test.go` (disabled link contract), `web/styles_data_table_test.go` (disabled page spans), `web/styles_fab_test.go` (`tabindex="-1"`).

---

### 1.11 Loading

**Architecture rule (mandatory)**
Loading state must be announced: the control exposes `aria-busy="true"` and its accessible name changes to `Loading {Label}`. Operation-level loading uses the native `<progress>` element. Never use an ad-hoc spinner without `aria-busy` and a dynamic name.

**Application in Gelium UI**
- Button/CTA loading: `aria-busy="true"`, spinner `aria-hidden="true"`, `sr-only` "Loading {Label}", visible label kept (`button.html:4,9`).
- Operation loading: `<progress>` determinate/indeterminate (`progress.html:5-23`); the refresh demo reuses `.ui-progress` + toast (`data-table.html:81-87`).
- Skeleton (initial data load): container `role="status"` + `sr-only` "Loading" text + `aria-hidden` placeholder blocks (`skeleton.html:1-3`). Shimmer disabled under reduced motion (`skeleton.css`).
- Loading is a transient-action pattern; never used for persistent feedback (invariant §0.4).

**Evidence / tests**
- `web/styles_button_test.go` (aria-busy, sr-only name), `web/styles_skeleton_test.go` (`TestSkeletonReducedMotionDisablesAnimation`), `web/styles_progress_test.go`.

---

### 1.12 Validation HTTP 422

**Architecture rule (mandatory)**
Server-side validation failures return HTTP 422 with the `X-Gelium-Validation: true` header. The re-render (full page no-JS, or fragment with HTMX) carries: per-field `aria-invalid="true"` + `role="alert"` message + `aria-describedby`, and a form-level **validation summary** with real anchor links to each invalid field. Validation never fires `gelium:toast`.

**Application in Gelium UI**
- Contract header: `X-Gelium-Validation: true` (`text_field.go:88`); `app.js:1-9` instructs HTMX to swap 422-with-header as success.
- Per-field: `text-field.html:5,8` (`aria-invalid` + `role="alert"` + `aria-describedby`); `select.html:84` (`ui-select-error` `role="alert"`).
- Validation summary: `validation-summary.html` — `role="alert"` container, heading level configurable, `<ul>` of `<li><a href="#{field}-error">` anchors that jump natively to the field.
- Value preservation on 422 (`text_field.go:62`); focus recovery via autofocus (§1.7).
- Rule "validation ≠ toast" enforced in `toast.go:129-133` and by contract.

**Evidence / tests**
- `web/styles_validation_summary_test.go`, `web/styles_text_field_test.go`, `internal/app/text_field_test.go`, `internal/app/select_test.go`, `internal/app/toast_test.go` (validation-never-toast).

---

### 1.13 Live regions

**Architecture rule (mandatory)**
Transient feedback is announced via a polite live region; persistent errors bound to context use `role="alert"` (assertive) or `role="status"` (polite) placed near the context. The toast region is the single live region for transient action feedback.

**Application in Gelium UI**
- Toast region: `#gelium-toast-region` with `aria-live="polite" aria-atomic="false" aria-relevant="additions text"` (`toast.html:10`); each toast sets `role="status"` (info/success/warning) or `role="alert"` (error) (`toast.html:2`, `app.js:45`).
- Server-driven contract: `HX-Trigger {"gelium:toast":…}` with closed vocabulary `info|success|warning|error`; no-JS inline toast fallback (`toast.go:161-164`).
- State patterns map: `role="alert"` for error tones, `role="status"` for info/success (see §2); `empty-state` uses `role="status"`; skeleton uses `role="status"`.
- Persistent-contextual feedback is never announced through the toast region (§0.4).

**Evidence / tests**
- `web/styles_toast_test.go` (aria-live, role derivation, reduced-motion, forced-colors), `internal/app/toast_test.go`.
- State-pattern role derivation: `web/styles_banner_test.go`, `web/styles_inline_alert_test.go`, `web/styles_empty_state_test.go`, `web/styles_error_state_test.go`, `web/styles_skeleton_test.go`, `web/styles_validation_summary_test.go`.

---

### 1.14 Reduced motion

**Architecture rule (mandatory)**
Under `prefers-reduced-motion: reduce` every animation/transition that could cause vestibular issues is disabled or reduced to an opacity change. Each moving component owns a reduced-motion block; static primitives must not declare one (and must not animate at all).

**Application in Gelium UI**
- Central block in `app.css:60-77` (button, text-field, dialog, toast, elevation, switch, select, select-menu, slider, progress, fab, list) + per-component blocks (tabs, navigation-*, segmented-button, icon-button, tooltip, menu, chips, data-table, skeleton, toast, text-field).
- Static state primitives (empty-state, callout, error-state, inline-alert, validation-summary, banner) deliberately declare no animation and no reduced-motion block — enforced by tests.
- **G11 resolved**: the checkbox/radio `:active` `scale(.92)` is dropped under reduced motion (`checkbox.css:95`, `radio.css:103`).

**Evidence / tests**
- `web/styles_toast_test.go`, `web/styles_tooltip_test.go`, `web/styles_text_field_test.go`, `web/styles_tabs_test.go`, `web/styles_skeleton_test.go` (must declare reduced-motion block) and `styles_empty_state_test.go` / `styles_callout_test.go` / `styles_error_state_test.go` / `styles_inline_alert_test.go` / `styles_validation_summary_test.go` / `styles_banner_test.go` (must NOT declare it).

---

### 1.15 Forced colors

**Architecture rule (mandatory)**
Under `forced-colors: active`, every component must stay discernible using the system palette (`CanvasText`, `LinkText`, `Highlight`, `Canvas`, `Field`, etc.), must not depend on background color alone, and its focus indicator must switch to `Highlight`.

**Application in Gelium UI**
- Central block `app.css:79-213` + per-component blocks (menu, chips, data-table, navigation-*, segmented-button, icon-button, tooltip, tabs, toast, text-field, dialog, switch, banner, callout, empty-state, error-state, inline-alert, skeleton, validation-summary).
- Focus ring switches to `Highlight` (e.g. `tabs.css`, `tooltip.css`).
- Containers get `forced-color-adjust: auto` + explicit `border: 1px solid CanvasText` where the surface would otherwise vanish.

**Evidence / tests**
- Forced-colors assertions in the vast majority of `web/styles_*_test.go` (banner, button, callout, card, checkbox, chips, data_table, dialog, divider, elevation, empty_state, error_state, fab, icon_button, inline_alert, list, menu, navigation_*, progress, radio, segmented_button, select, select_menu, skeleton, slider, switch, tabs, text_field, toast, tooltip, validation_summary).
- `web/styles_contract_test.go` (central forced-colors wiring in compiled `app.css`).

---

### 1.16 Contrast (AA)

**Architecture rule (mandatory)**
Text and UI components must meet WCAG 2.1 AA contrast (4.5:1 body text, 3:1 large text and non-text contrast) in light and dark. Contrast is the responsibility of **semantic tokens**, not component hardcoded colors; components reference the role (`var(--ui-color-*)`), the theme defines AA-compliant values.

**Application in Gelium UI**
- Roles drive color: `--ui-color-fg` on `--ui-color-canvas`, `--ui-color-danger` for errors, `--ui-color-fg-muted` for secondary text, etc. (`theme-material/theme.css`, `tokens.css`).
- Audit baseline (Phase E evidence): the Material palette satisfies AA by design (e.g. `--ui-color-fg-muted #49454f` on `#fff7ff`; danger ≈ 8:1); demo tones also comply.
- Token drift **G11 resolved**: `--ui-color-error` is a compatibility alias of the canonical `--ui-color-danger` (`tokens.css:38`) and no component hardcodes a fallback; consumers reference `danger` (or the alias, which resolves to it).

**Evidence / tests**
- `internal/app/server_test.go` (`TestMaterialDarkThemeKeepsFilledFieldDistinctFromSurface`, `TestMaterialThemeDefinesTextFieldTypescaleTokens`, `TestMaterialThemeExposesSemanticFoundationContracts`).
- `web/styles_contract_test.go` (semantic token family presence). Visual AA verification is manual per theme (`roadmap.md` DoD).

---

### 1.17 Target sizes

**Architecture rule (mandatory)**
Interactive targets must offer a minimum 24px hit area (WCAG 2.1 AA non-text contrast exception baseline; Gelium targets the AA minimum and prefers 44px where density allows). Sizes come from the `--ui-size-*` token family, never per-component literals.

**Application in Gelium UI**
- `--ui-size-control: 2.5rem` (40px) drives button, icon-button, segmented-button minimum geometry (`tokens.css:120`); `--ui-size-icon: 1.5rem`; `--ui-size-item-*` for list/menu/drawer/data-table row heights.
- `button.css:7` `min-height: var(--ui-size-control)`; `icon-button.css:5-6` width/height `var(--ui-size-control)`.
- Dismiss buttons in banner/toast are `.ui-icon-button` (40px) so the dismiss target meets the minimum.

**Evidence / tests**
- `web/styles_button_test.go`, `web/styles_icon_button_test.go`, `web/styles_fab_test.go` (size token wiring), `web/styles_contract_test.go` (size family presence).
- `web/styles_banner_test.go`, `web/styles_toast_test.go` (dismiss hit area).

---

## 2. State patterns — ARIA contract

Every state pattern is a server-rendered component. "Role" is derived, never overridable per instance; tone and role are coupled.

| Pattern | Role | When | Server contract | Persistence |
|---|---|---|---|---|
| Empty state | `role="status"` | A list/table/search has zero results; gives the user guidance or a CTA | GET params; output of the handler, never client state (`data_table.go` empty row, `colspan`) | Persistent-contextual |
| Skeleton | `role="status"` + `sr-only` text + `aria-hidden` blocks | Initial data load placeholder | Server-rendered static HTML; replaced by the next request | Persistent-contextual (placeholder) |
| Inline alert | `role="alert"` (error tone) / `role="status"` (info/success/warning) | Section/form-level feedback bound to context; per-field error uses the field's own `role="alert"` message | 422 + `X-Gelium-Validation` re-render, or fragment | Persistent-contextual |
| Banner | `role="alert"` (error tone) / `role="status"` (info/success/warning) | Page/site-level notice; never auto-dismissed; dismiss = form `POST` + 303 | POST + 303 redirect (dismiss); server-rendered in layout slot between `</header>` and `<main>` | Persistent-contextual |
| Callout | none (static `aside`) | Optional informational note, ignorable; tone variants must never be color-only | None (static markup) | Persistent-contextual |
| Error state | `role="alert"` | Full-page/resource failure; single `h1`, descriptive body, real retry link, real HTTP status | HTTP status (404/500/503) + server-rendered re-render | Persistent-contextual |
| Validation summary | `role="alert"` | Form-level list of invalid fields; anchors `#field-error` jump natively; heading level configurable | 422 + `X-Gelium-Validation` (no new header) | Persistent-contextual |
| Success (persistent) | `role="status"` | Post-save/page confirmation that must survive navigation; reuses Banner/Inline alert success tone | POST + 303 redirect re-render, or post-submit fragment; NEVER `gelium:toast` | Persistent-contextual |
| Toast | `role="status"` (info/success/warning) / `role="alert"` (error) | Transient action result; auto-dismiss 4s/8s pausable; single live region | `HX-Trigger {"gelium:toast":…}`; no-JS inline fallback | Transient-action |
| Loading (button/operation) | `aria-busy` + dynamic name `Loading {Label}`; `<progress>` for operations | In-progress action | Server-rendered state | Transient-action |

**Rule**: `role="alert"` only for content that needs immediate assertive interruption (errors, required action); everything else polite (`role="status"`). No persistent pattern emits `gelium:toast`; no transient pattern occupies a persistent slot (§0.4).

---

## 3. Known gaps — resolution status

Tracked in `docs/handoffs/ux-accessibility-audit.md` §8 and `docs/gelium-ui-system-roadmap.md` Phase E.

| Gap | Severity | Description | Status |
|---|---|---|---|
| **G1** — Overlays without no-Chromium fallback | Critical | Dialog and Select menu open only via Invoker Commands (`command`/`commandfor`); in Firefox/Safari the trigger is an inert button; `select-menu.css:2-8` promises a native `<select>` fallback that is not rendered | **RESOLVED (commit `86a7f71`)** — the select-menu demo is the component's own native `<select>` (no dead form) and the dialog trigger is a real link to the server-rendered `/components/dialog/confirm` page; the `<dialog>` + Invoker Commands stays only as an opt-in enhancement. Escape/cancel contract kept where the modal exists (§1.9) |
| **G2** — `lang="en"` on Spanish demos | High | `demo-whatsapp.html`, `demo-whatsapp-admin.html` declare `lang="en"` while content is Spanish — SR reads Spanish with an English voice; also `layout.html` hardcodes `lang="en"` | **RESOLVED (commit `9504216`)** — both demos render `lang="es"` (`demoMetaES`); `layout.html` keeps `en` for the English docs site |
| **G3** — Dead webhook form | High | Admin webhook form `POST /demo/whatsapp/admin` has no registered POST handler → 405 with no feedback; several placeholder `href="#"` links and an inert "regenerate" button | **RESOLVED (commit `9504216`)** — `POST /demo/whatsapp/admin/webhook` persists and redirects (POST+303); placeholder `href="#"` demo links and the inert "regenerate" button remain documented demo scaffolding |
| **G4** — Data table without empty state | High | 0 rows rendered as a `0 rows` caption + empty tbody; "Select all" checkbox wrongly `checked` with 0 rows | **RESOLVED (Phase D)** — `data_table.go` renders an `empty-state` row with `colspan`; select-all is omitted when `Total == 0`; covered by `data_table_test.go` |
| **G5** — No transport error feedback in HTMX | High | `app.js` handles only 422-with-header; a 500/network failure during an HTMX request shows nothing (no `hx-on::response-error`, no transport live region, no retry) | **RESOLVED (commit `9504216`)** — `app.js:88-94` surfaces a generic error toast on `htmx:responseError`/`htmx:sendError`; the 422-with-header hook keeps "validation is never a toast" intact |
| G6 | Medium | Accessible-name drift: `aria-label` overrides visible labels (filter, slider, progress) in demos | **RESOLVED** — the `aria-label` overrides were removed; every control is named by its visible `<label>` (`data-table.html:10`, `slider.html`, `progress.html`); covered by tests |
| G7 | Medium | No skip link; `main` landmark missing in the admin demo | **RESOLVED** — `layout.html` ships a skip link as the first focusable element (`#main-content`) with a `:focus` treatment (`base.css`); the admin demo exposes a `<main class="demo-wa-admin">` |
| G8 | Medium | No `aria-expanded` on popover/menu/select triggers | **RESOLVED** — the three menu popover triggers carry `aria-expanded="false"` (`menu.html`); the native `<select>` needs none; the dialog/drawer modal trigger remains UA-managed (no JS to keep the attribute in sync — documented nuance) |
| G9 | Medium | Admin active tabs class-only (`aria-current` missing); emoji action links without `aria-label` | **RESOLVED** — the active admin tab carries `aria-current="location"` and every emoji-only action link has an `aria-label` (`demo-whatsapp-admin.html`) |
| G10 | Low | Invalid `aria-selected` on select-menu buttons; redundant `role="list"`; no-JS toast dismiss handler unbound; tooltip rich action not mouse-clickable | **PARTIALLY RESOLVED** — `aria-selected` is gone (native `<select>`), redundant `role="list"` removed from the demo; the no-JS toast dismiss affordance and the tooltip rich action remain documented demo limitations |
| G11 | Low | Checkbox/radio activation animation under reduced motion; `--ui-color-error` drift | **RESOLVED** — the `:active` `scale(.92)` is dropped under reduced motion (`checkbox.css`, `radio.css`); `--ui-color-error` is a compatibility alias of the canonical `--ui-color-danger` (`tokens.css:38`) with no hardcoded fallbacks |

**Summary**: G1–G9 fully resolved (G1, G2, G3, G5 in Phase E commits `86a7f71`/`9504216`; G6–G9 in the cleanup slice); G10/G11 partially resolved with the remaining demo-scaffold limitations documented in their rows.

---

## 4. Tests that guarantee the contract

### 4.1 Server contract tests (`internal/app/`)
- `server_test.go` — banner slot position and role derivation; inline-alert role derivation; error-state markup with single `h1`; unknown route renders error state; theme class server-driven; semantic foundation contracts; layout h1 ownership; docs index.
- `text_field_test.go` — 422 + `X-Gelium-Validation`, autofocus on no-JS re-render, no autofocus in HX branch, value preservation, `role="status"` success helper.
- `select_test.go` — 422 validation, error `role="alert"`, native select vs select-menu.
- `toast_test.go` — `gelium:toast` contract, closed vocabulary, role derivation, no-JS inline fallback, **validation never toast**.
- `data_table_test.go` — empty state row (`colspan`, message, CTA), select-all absent at 0 rows, HX empty-state fragment, query escaping, sort/pagination links.
- `demo_whatsapp_test.go` — demo screens render (baseline for G2/G3 cleanup).

### 4.2 Style contract tests (`web/styles_*_test.go`)
Per component, three guarantees: primitive CSS maps tokens, contract wired into compiled `app.css` (`@media (forced-colors:active)` presence), and closed class vocabulary matching the template.
- Semantic HTML / roots: `styles_contract_test.go`, `styles_menu_test.go`, `styles_tabs_test.go`, `styles_dialog_test.go`, `styles_data_table_test.go`, `styles_select_menu_test.go`, `styles_progress_test.go`, `styles_segmented_button_test.go`.
- Accessible names: `styles_icon_button_test.go`, `styles_fab_test.go`, `styles_icon_test.go`, `styles_button_test.go`, `styles_banner_test.go`, `styles_toast_test.go`.
- Focus visible: `styles_focus_ring_test.go`, `:focus-visible` assertions in `styles_button/text_field/tabs/toast/tooltip/switch/validation_summary` tests.
- Disabled: `styles_button_test.go`, `styles_fab_test.go`, `styles_data_table_test.go`.
- Loading: `styles_button_test.go`, `styles_skeleton_test.go` (reduced-motion disable), `styles_progress_test.go`.
- Validation: `styles_validation_summary_test.go`, `styles_text_field_test.go`.
- Live regions / state patterns: `styles_banner_test.go`, `styles_inline_alert_test.go`, `styles_empty_state_test.go`, `styles_error_state_test.go`, `styles_skeleton_test.go`, `styles_validation_summary_test.go`, `styles_callout_test.go`, `styles_toast_test.go`.
- Reduced motion: blocks required in moving components (`toast`, `tooltip`, `text_field`, `tabs`, `skeleton`) and forbidden in static ones (`empty_state`, `callout`, `error_state`, `inline_alert`, `validation_summary`, `banner`).
- Forced colors: present across all `styles_*_test.go` (banner, button, callout, card, checkbox, chips, data_table, dialog, divider, elevation, empty_state, error_state, fab, icon_button, inline_alert, list, menu, navigation_*, progress, radio, segmented_button, select, select_menu, skeleton, slider, switch, tabs, text_field, toast, tooltip, validation_summary).
- Target sizes: `styles_button_test.go`, `styles_icon_button_test.go`, `styles_fab_test.go`, `styles_banner_test.go`, `styles_toast_test.go` (size-token wiring).

### 4.3 Verification commands (every change)
```bash
npm run build
go test ./...
go vet ./...
```
Manual smoke per `roadmap.md`: light/dark, narrow/wide, reduced motion, forced colors, keyboard, no-JS, HTMX, empty/loading/error, console without errors.

---

## 5. Sources of authority

`docs/gelium-ui-core.md` (§6 HTML-first, §7 server-driven contracts, §9 Accesibilidad), `docs/gelium-ui-system-roadmap.md` (Phase D/E), `docs/handoffs/ux-accessibility-audit.md`, `docs/handoffs/state-patterns-audit.md`, `docs/gelium-ui-composition-rules.md`, `docs/gelium-ui-vocabulary.md`, `web/templates/*.html`, `web/styles/*.css`, `web/static/app.js`, `internal/app/*.go`, `themes/theme-material/theme.css`.
