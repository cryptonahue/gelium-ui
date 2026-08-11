# Gelium UI — Content Rules

> Copy and content rules for the Gelium UI system: what to say, how to say it, and where.
> Phase E of the system roadmap (`docs/gelium-ui-system-roadmap.md`).
> Base: `docs/gelium-ui-ux-principles.md`, `docs/gelium-ui-ux-patterns.md`, `docs/gelium-ui-vocabulary.md`, `docs/handoffs/state-patterns-audit.md`.

Rules apply to every string rendered by the system and every string a consuming project writes into Gelium components. Each rule: the rule, a good example, a bad example, and where it applies (component/pattern).

---

## 1. Error messages

**Rule.** Every error says **what happened** and **what to do next**. Never a bare "Something went wrong" or a technical-only message. Prefer specific, actionable copy over generic panic.

- **Good**: "We couldn't save your changes. Review the highlighted fields or try again."
- **Bad**: "Something went wrong."
- **Bad**: "Error 500: internal server error."
- **Applies to**: Error state (`error-state.html`), Inline alert (`inline-alert.html`), Banner (`banner.html`), Toast error (`toast.html`), Validation summary (`validation-summary.html`).

**Rule.** Field errors identify the field and the required fix.

- **Good**: "Email is required. Enter the email address you signed up with."
- **Bad**: "Invalid input."
- **Applies to**: Text field error message (`text-field.html:8`, `role="alert"` + `aria-describedby`), Select error, Validation summary list items.

---

## 2. Button labels

**Rule.** Button labels are **action verbs** describing exactly what the button does. Never a bare noun, never "OK"/"Submit" when the action is specific.

- **Good**: "Save", "Delete", "Cancel", "Send", "Retry", "Apply filter".
- **Bad**: "OK", "Submit", "Confirm" (when "Delete" is the action), "Go".
- **Applies to**: Button (`button.html`), Dialog actions (`dialog.html` — "Cancel"/"Confirm" only when "Cancel"/"Confirm" is literally the action).

**Rule.** Destructive buttons state the destructive verb explicitly.

- **Good**: "Delete record".
- **Bad**: "Process".
- **Applies to**: Destructive action pattern, Dialog confirm.

---

## 3. Empty states

**Rule.** An empty state explains why the surface is empty and offers a **real, actionable CTA** (create, clear filter, retry). Never leave a bare "No data" or "0 rows".

- **Good**: "No projects yet. Create your first project to get started." + CTA "Create project".
- **Bad**: "0 rows".
- **Applies to**: Empty state (`empty-state.html`), Data table empty row (`data-table.html:68-70`), Search results, Filters.

**Rule.** After a filtered/search empty, name the active filter and offer "Clear filters".

- **Good**: "No records match `active`. Clear filters to see all records." + CTA "Clear filters".
- **Bad**: "No results."
- **Applies to**: Empty state with active filters (Search, Filters patterns).

---

## 4. Loading states

**Rule.** Loading is always visible and comprehensible: a spinner + label on buttons, determinate `<progress>` on operations, Skeleton for data regions. Never a silent freeze.

- **Good**: Button spinner + sr-only "Loading {Label}" (`button.html:4,9`); `<progress>` for refresh (`data-table.html:81`); Skeleton placeholder (`skeleton.html`).
- **Bad**: A disabled button with no indication of progress.
- **Applies to**: Loading pattern, Button, Progress, Skeleton.

**Rule.** Never tell the user to "wait" as a message; show progress instead.

- **Good**: Determinate progress or skeleton.
- **Bad**: "Please wait…".
- **Applies to**: Loading pattern.

---

## 5. Confirmation copy

**Rule.** A confirmation states **what will happen** and **what can be undone** (or that it cannot be undone). The confirm button repeats the verb.

- **Good**: "Delete project `Acme`? This removes it and its 12 records. This action cannot be undone." + buttons "Cancel" / "Delete".
- **Bad**: "Are you sure?" + buttons "OK"/"Cancel".
- **Applies to**: Confirmation and Destructive action patterns, Dialog (`dialog.html` content + actions).

---

## 6. Permissions and access

**Rule.** Permission copy names the resource, the principal, and the consequence of granting/removing.

- **Good**: "Grant `editor` access to `Ana`? She will be able to edit this workspace's projects."
- **Bad**: "Change permissions?"
- **Applies to**: Permissions pattern, Banner/Inline alert on permission change, Dialog confirmation.

---

## 7. Helper text

**Rule.** Helper text prevents errors: it explains requirements, format and consequences *before* submission. It uses `role="status"` when it is static guidance (`text-field.html:8`).

- **Good**: "Use at least 8 characters with a number and a symbol."
- **Bad**: No helper for a field with strict format requirements.
- **Applies to**: Text field helper, Select helper.

---

## 8. Tone

**Rule.** Tone is **clear, direct, neutral, professional** — and never blames the user. Errors describe the situation, not the person.

- **Good**: "We couldn't process your payment. Check the card details or use another card."
- **Bad**: "You entered an invalid card number." (blames) / "Oops, something exploded!" (flippant).
- **Applies to**: All messages (errors, empty, loading, success), all components and patterns.

**Rule.** Success messages confirm the action without flattery or excess.

- **Good**: "Settings saved."
- **Bad**: "Amazing! You did it!!!"
- **Applies to**: Success feedback (Banner/Inline `role="status"`), Toast success.

---

## 9. Language

**Rule.** System technical artifacts and UI strings are **English by default**, unless the consuming project defines another locale. The system is **multi-language by contract**: content is server-rendered data, so a consuming project can localize any string through its own templates/translations without touching components.

- **Good**: `message: "Record deleted."` — the consuming project supplies its localized string.
- **Bad**: Hardcoded Spanish/any locale inside a component template.
- **Applies to**: Every string produced by `internal/app/*` (toast messages, errors, empty states, buttons).

**Rule.** Set `lang` on the rendered document to the actual content language (audit G2: demos in Spanish must not carry `lang="en"`).

- **Good**: `<html lang="es">` for Spanish content.
- **Bad**: `<html lang="en">` with 100% Spanish content.
- **Applies to**: Layout (`layout.html`), demo pages.

---

## 10. Plain language

**Rule.** Use plain, everyday words over jargon, and short titles with an actionable body.

- **Good**: Title "You're offline" + body "Reconnect and try again."
- **Bad**: Title "Network transport error" + body "An unexpected exception occurred while establishing a connection."
- **Applies to**: Error state, Banner, Inline alert, Empty state, Toast.

**Rule.** Keep titles short (5-9 words); put the action in the body or CTA.

- **Good**: "Payment failed" + "Check the card details or use another card."
- **Bad**: "There was a problem processing your payment information with the provider."
- **Applies to**: All message components.

---

## 11. Actionable messages

**Rule.** When an action is available, name it as a real control (link/button), not only in prose.

- **Good**: "No results. Clear filters." where "Clear filters" is a real link.
- **Bad**: "No results — you may want to clear your filters." (instruction with no control).
- **Applies to**: Empty state CTA, Error state retry, Banner CTA (`banner.html`), Validation summary links.

---

## Summary table

| Rule | Good | Bad | Applies to |
|---|---|---|---|
| Errors: what + what to do | "We couldn't save your changes. Review the highlighted fields or try again." | "Something went wrong." | Error state, Inline alert, Banner, Toast |
| Buttons = action verb | "Delete", "Save", "Retry" | "OK", "Submit" | Button, Dialog actions |
| Empty = message + CTA | "No projects yet. Create your first project." | "0 rows" | Empty state, Data table, Search |
| Loading = visible + comprehensible | Spinner + label, progress, skeleton | Silent disabled button | Loading, Button, Progress, Skeleton |
| Confirm = what + undo | "Delete `Acme`? This cannot be undone." | "Are you sure?" | Confirmation, Destructive action, Dialog |
| Tone = neutral, no blame | "We couldn't process your payment." | "You entered an invalid card." | All messages |
| Language = English by default, localized by contract | Server-supplied message | Hardcoded locale in template | `internal/app/*` strings |
| Plain language, short titles | "Payment failed" + action | "Network transport error" + jargon | All message components |
| Actionable = real control | "Clear filters" as link | "you may want to clear your filters" | Empty state, Error state, Banner, Validation summary |

---

**Definition of done (Phase E scope for this doc)**: rules written with good/bad examples and component anchors, referenced by `composition-rules.md` and by the Phase G screen recipes.
