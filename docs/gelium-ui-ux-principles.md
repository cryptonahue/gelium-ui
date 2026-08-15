# Gelium UI — UX Principles

> Principles that govern the experience of the Gelium UI system.
> Phase E of the system roadmap (`docs/gelium-ui-system-roadmap.md`).
> Base: `docs/handoffs/ux-accessibility-audit.md`, `docs/gelium-ui-composition-rules.md`, `docs/gelium-ui-vocabulary.md`, `docs/gelium-ui-core.md` and the real catalog (`internal/app/*`, `web/templates/*`).

Every screen recipe, new component and content decision MUST be justified against these principles. They are the "why" behind the composition rules (§4) and the state patterns of Phase D.

---

## 1. Clarity over ornamentation

**Definition.** The interface communicates structure and meaning before decoration. Every visual element earns its place; nothing exists purely for embellishment.

**How Gelium UI applies it.**
- Themes are single files over the `--ui-*` token contract (`themes/theme-material/theme.css`); components never hardcode Material classes (`roadmap.md:52` → `m3-select-trigger` → `ui-select-trigger`).
- The Material Gelium vocabulary is semantic-first: Card, List, Data table, Badge, Divider, Elevation, Focus ring exist to encode meaning, not style.
- State layers are theme-aware `color-mix()` over `-fg` tokens, never decorative `rgb()` literals (`button.css:17-18`).

**Evidence.** `web/styles/tokens.css` (semantic roles), `COMPONENT-ROADMAP.md`, `demo-whatsapp.css:403` documented as drift (hardcoded fallback defeats clarity).

**Accessibility.** Decoration that carries meaning breaks forced-colors and screen readers; semantic HTML + tokens keeps meaning transport-independent.

---

## 2. Recognition over recall

**Definition.** Users recognize the right option from visible affordances instead of remembering how to reach it. Persistent, visible navigation and state beat memory-based discovery.

**How Gelium UI applies it.**
- Navigation is always visible and native: Tabs (`tabs.html`), Navigation drawer/bar/tab, Breadcrumbs, Menu, List — never hidden behind gestures or JS-only controls.
- The URL is the state (`composition-rules.md:173`): sort, filter, page and selection are recognizable links, not remembered client state.
- Active states are explicit: `aria-current="page"` on tabs, `aria-current="page"` on pagination, `aria-sort` on columns — the user never has to remember what is selected.

**Evidence.** `web/templates/tabs.html:6-18` (`aria-current`), `data-table.html:74-78` (pagination links), `navigation-drawer.html:8-38`.

**Accessibility.** Recognition via text + state attributes (never color-only, `composition-rules.md:132`) is what makes the system work for screen readers and forced-colors.

---

## 3. Visibility of system status

**Definition.** The system always tells the user what is happening: what changed, what failed, what is loading, what was saved.

**How Gelium UI applies it.**
- Transient results → Toast (`gelium:toast` contract, `web/templates/toast.html`, `toast.go:108-127`).
- Persistent, contextual results → Inline alert / Banner / Validation summary (Phase D), never Toast (`composition-rules.md:126`).
- Loading → `aria-busy` on buttons (`button.html:4,9`), native `<progress>` (`progress.html`), Skeleton for data regions (`skeleton.html`).
- Empty and error states are server-rendered output, not silent voids (`empty-state.html`, `error-state.html`).

**Evidence.** `composition-rules.md:119` (toast vs inline alert vs banner vs callout), `ux-accessibility-audit.md:22-24` (persistent vs transient classification).

**Accessibility.** Every status change is announced through `role="status"`/`role="alert"` or `aria-live` — status is never implicit in color alone.

---

## 4. Comprehensible feedback

**Definition.** Feedback states what happened, what it means, and what to do next — in plain language, at the right moment, in the right channel.

**How Gelium UI applies it.**
- Errors say what + what to do: see `docs/gelium-ui-content-rules.md` §Error messages.
- Validation feedback is per-field (`aria-invalid` + `aria-describedby` + persistent `role="alert"`, `text-field.html:5,8`) AND form-level (Validation summary with real anchor links, `validation-summary.html`).
- Persistent feedback is never masked as a transient toast (`toast.go:129-133`).
- Transport errors (network/500 in HTMX) are the Phase E gap G5 (`ux-accessibility-audit.md:88`): until resolved, no silent failure.

**Evidence.** `web/templates/text-field.html:5,8`, `validation-summary.html`, `toast.go:129-133`.

**Accessibility.** `role="alert"` is assertive; `role="status"` is polite. Misusing either interrupts or hides feedback for assistive technology.

---

## 5. Error prevention

**Definition.** The system prevents errors before they happen: constraints, helpful defaults, and pre-submit guidance reduce the need for error messages.

**How Gelium UI applies it.**
- Native constraints first: `<input type=…>`, `required`, pattern hints, `<select>`, radio/checkbox/segmented semantics reduce free-form input (`checkbox.html`, `radio.html`, `segmented-button.html`).
- Helper text (`role="status"`, `text-field.html:8`) explains requirements *before* the user errs.
- Closed vocabularies are sanitized server-side (`composition-rules.md:170`) — the server never trusts the client.
- Confirmation dialogs guard high-impact actions before they commit (see §Action hierarchy).

**Evidence.** `web/templates/text-field.html:8`, `segmented-button.html:10,50`, `composition-rules.md:168-171`.

**Accessibility.** Prevention lowers the volume of `role="alert"` interruptions, which is itself an accessibility win.

---

## 6. Error recovery

**Definition.** When an error happens, the user can recover quickly and understand why. Recovery is explicit, low-cost, and loses no data.

**How Gelium UI applies it.**
- **422 + `X-Gelium-Validation`** preserves the submitted value and returns focus to the failing field (`text_field.go:62-67`) — nothing typed is lost.
- Validation summary links (`<a href="#{field}-error">`) let users jump straight to each problem (`validation-summary.html`).
- Page/resource errors render a real `error-state.html` with status code + retry button; POST + 303 keeps workflows resumable.
- Re-submit re-renders; helper text replaces the error (`ux-accessibility-audit.md:89`).

**Evidence.** `text_field.go:64-91`, `error-state.html`, `web/static/app.js:1-9` (HTMX 422 swap), `ux-accessibility-audit.md:84-89`.

**Accessibility.** `autofocus` return-to-field is skipped in the HTMX branch so it does not steal focus (`ux-accessibility-audit.md:79`) — recovery must not itself be a disruption.

---

## 7. Consistency (internal + platform)

**Definition.** The same concept is always represented by the same component, token and word — both inside the product and against platform/native expectations.

**How Gelium UI applies it.**
- **Internal**: one canonical vocabulary (`docs/gelium-ui-vocabulary.md`); a Toast is a Toast everywhere, a Callout is always an ignorable note, validation is always 422. Conflicts resolved in `vocabulary.md:328-339`.
- **Platform**: native controls over custom widgets (native `checkbox`/`radio`/`select`/`progress`/`<dialog>`), so keyboard, focus and semantics match what users expect from the platform.
- **Themes**: Material and future Basecoat map to the same Gelium contract; markup is never forked (`roadmap.md:52,284`).

**Evidence.** `docs/gelium-ui-vocabulary.md` §8 (naming conflicts resolved), `dialog.html:3` (`<dialog closedby="any">`), `tokens.css` (semantic color vocabulary).

**Accessibility.** Platform-native semantics give assistive tech the expected roles for free — the strongest ARIA is the one you do not write.

---

## 8. Action hierarchy (primary / secondary / destructive)

**Definition.** Every surface has a clear primary action; secondary actions are visually quieter; destructive actions are explicit, separated, and confirmed.

**How Gelium UI applies it.**
- Button variants encode hierarchy: `ui-button-primary`, `ui-button-secondary`, `ui-button-outline`, `ui-button-text` (`web/styles/button.css:25-28`).
- Destructive intent is signalled through the danger semantic token (`--ui-color-danger`, `tokens.css:32`) and MUST be confirmed in a Dialog before the `POST + 303` commit (`dialog.html`, `demo_whatsapp.go:559`).
- Only one primary action per surface (`composition-rules.md:35`).

**Evidence.** `web/styles/button.css:25-28`, `dialog.html:6-9`, `tokens.css:32`, `ux-accessibility-audit.md:50` (dialog confirm focus).

**Accessibility.** The danger token is never the sole signal: text + confirmation dialog + focus management carry the meaning (never color-only, `composition-rules.md:132`).

---

## 9. Progressive disclosure

**Definition.** Complexity is revealed on demand: show what is needed now, hide what belongs to later steps or optional detail.

**How Gelium UI applies it.**
- Multi-step flows (Steps, Editor/Booking) validate step-by-step with 422 + Inline alert per step (`vocabulary.md:190-197`).
- Dialog is used for short, reversible sub-tasks so the page does not grow (`composition-rules.md:37`).
- Callout presents optional informational context that can be ignored (`vocabulary.md:131-138`).
- Menus and drawers hide secondary actions until invoked — but never essential actions (see §13 Tooltip).

**Evidence.** `vocabulary.md:190-197` (Steps), `composition-rules.md:108-109` (dialog anti-rule), `menu.html`.

**Accessibility.** Disclosure must not hide content from assistive tech: hidden detail remains reachable via native dialogs/popovers/anchors, never via hover-only or JS-only reveals.

---

## 10. Density by surface

**Definition.** Visual density adapts to the surface: public portals favor spacious, reading-friendly layouts; admin surfaces allow denser scanning; ops surfaces favor compact information-dense rows.

**How Gelium UI applies it.**
- Density is the responsibility of the core/theme via `--ui-density-*` tokens (`composition-rules.md:135-140`), never of screen recipes.
- Recipes consume geometry through `--ui-size-*`; literals in recipes are forbidden.
- Until density tokens exist (Phase B), no screen may define per-screen densities (`composition-rules.md:137`).

**Evidence.** `composition-rules.md:135-140`, `tokens.css:119` (density family deliberately absent until consumed).

**Accessibility.** Density must respect target sizes and readability; compact surfaces still require accessible touch targets and AA contrast.

---

## 11. Responsive / mobile-first behavior

**Definition.** Layouts are fluid first and adapt by content, not by device sniffing or breakpoint-worship.

**How Gelium UI applies it.**
- Fluid-first: `auto-fit/minmax`, `min()/clamp()` before breakpoints (`composition-rules.md:141-148`, `card.css:30`, `dialog.css:3-5`).
- Breakpoints only when fluid layout is insufficient; overlays are already fluid (`calc(100vw - n)`).
- Large tables are resolved server-side via pagination, never horizontal scroll of the component (`composition-rules.md:145`).
- Drawer resolves responsive via variants (modal vs permanent), not media queries.

**Evidence.** `composition-rules.md:141-148`, `card.css:30`, `dialog.css:3-5`, `data-table.html:74-78` (pagination).

**Accessibility.** Document order equals visual order (`composition-rules.md:178`); fluid layouts preserve reading order without CSS reordering.

---

## 12. Page vs Dialog

**Definition.** Choose the surface by task depth and reversibility: pages own URL + back + deep-link; dialogs own short, reversible sub-tasks with retained context.

**How Gelium UI applies it.**
- **Use Dialog** for: confirm, quick edit, picker — short and focused (`dialog.html`).
- **Use Page** for: deep or complete flows, anything needing URL/back/deep-link (`composition-rules.md:34-38`).
- **Never** open a long flow (editor/booking multi-step) in a Dialog → Steps/page instead (`composition-rules.md:108-109`).
- The URL is the state: if the state is navigable, it is a page (`composition-rules.md:173`).

**Evidence.** `composition-rules.md:100-109` (§4.7), `dialog.html`, `vocabulary.md:245-255`.

**Accessibility.** `<dialog>` provides native focus trap, Escape and light dismiss (`closedby="any"`); a fallback server-rendered variant is required where Invoker Commands are not Baseline (gap G1, `ux-accessibility-audit.md:58`).

---

## 13. When NOT to use a Tooltip

**Definition.** Tooltips are a supplement, never the source. Essential information — errors, requirements, navigation, security — must never live only in a tooltip.

**How Gelium UI applies it.**
- Tooltip (`tooltip.html`, `role="tooltip"` + `aria-describedby`) is limited to short clarification of a labelled control (`vocabulary.md:271-279`).
- Essential info lives in visible text: helper text (`text-field.html:8`), inline alert, callout — never in a hover-only layer.
- Rich tooltips with interactive content are discouraged; action links inside tooltips are not reachable by mouse (`tooltip.css:43`, gap G10).

**Evidence.** `vocabulary.md:271-279`, `tooltip.html`, `ux-accessibility-audit.md:62,110` (G10).

**Accessibility.** Hover-only content is invisible to touch and unreliable for keyboard; `:focus-within` reveal (`tooltip.css:52-56`) is the acceptable minimum.

---

## Anti-rules derived from principles

| Principle | Anti-rule |
|---|---|
| Status visibility | Validation is never Toast (`toast.go:129-133`) |
| Feedback channels | Persistent/critical feedback never Toast → Inline alert / Banner (`composition-rules.md:126`) |
| Recognition | No `role="tablist"`/roving focus without full keyboard → links + `aria-current` (`composition-rules.md:128`) |
| Error recovery | No invented list state off the URL → GET with stable params (`composition-rules.md:127`) |
| Surface choice | No long flow in Dialog → page/steps (`composition-rules.md:129`) |
| Clarity | No spinner where Progress exists → `.ui-progress` + `aria-busy` (`composition-rules.md:130`) |
| Status visibility | No empty state without message/CTA (`composition-rules.md:131`) |
| Consistency | No color-only state → native control + forced colors (`composition-rules.md:132`) |
| Platform consistency | No JS for what a GET form already solves (`composition-rules.md:133`) |

---

**Definition of done (Phase E scope for this doc)**: principles written, anchored to real components/tokens/contracts, cross-referenced from `composition-rules.md` and referenced by the screen recipes.
