# Inline alert

Inline alert is a section- or form-level message that carries a persistent status for one region of the page — a failed save, a permission change, a row that could not be processed. Use an inline alert when a whole section needs a status that survives the interaction, instead of a transient toast or a full-page banner. It renders as a real `div` styled with the `ui-inline-alert` class and needs no component JavaScript.

## Examples

<div class="component-preview">
  <div class="ui-inline-alert ui-inline-alert--error" role="alert">
    <span class="ui-inline-alert-icon" aria-hidden="true">!</span>
    <p class="ui-inline-alert-title">We couldn't save your changes</p>
    <p class="ui-inline-alert-body">Review the highlighted fields or try again.</p>
  </div>
  <div class="ui-inline-alert ui-inline-alert--success" role="status">
    <span class="ui-inline-alert-icon" aria-hidden="true">✓</span>
    <p class="ui-inline-alert-body">Settings saved.</p>
  </div>
</div>

## Guidance

### When to use

Use an inline alert for a status that belongs to one section or form: a failed save, a permission change, a batch row summary. It sits next to the content it describes and persists until the condition is resolved. It is the section-level sibling of the full-page [banner](/components/banner) and the direct context for the [validation summary](/components/validation-summary) when a form fails.

### When not to use

Do not use an inline alert for a message that must disappear on its own — that is a [toast](/components/toast). Do not use it for a whole-site signal like maintenance or consent: that belongs to the layout banner. And do not use it for a single field error: the field itself carries its own message (`aria-invalid` + `aria-describedby`), and the inline alert summarizes the section without replacing field-level errors.

### Usability

- The tone is a closed vocabulary: `ui-inline-alert--{error,warning,success,info}`.
- The `role` derives from the tone: `error` becomes `alert` and everything else becomes `status`.
- Copy says what happened and what to do next — never a bare "Something went wrong" (see [Feedback](/docs/feedback), FEED-FAIL).

### Accessibility

- The icon is decorative: keep it `aria-hidden` and put the meaning in the title or body.
- Never render the alert color-only: the text carries the message and the role announces it.
- In forced-colors mode the alert keeps a `CanvasText` boundary and the error tone paints with `Mark`.

## Anatomy

- **Alert** — `ui-inline-alert`, the flex row with the scoped `--ui-inline-alert-*` tokens (`padding`, `gap`, `radius`, `icon-size`, `bg`, `fg`).
- **Icon** — `ui-inline-alert-icon`, a decorative `aria-hidden` glyph.
- **Title** — `ui-inline-alert-title`, the short label line (`--ui-type-label-lg`).
- **Body** — `ui-inline-alert-body`, the supporting sentence (`--ui-type-body-sm`).

The alert composes with — and never replaces — the field-level error on the Text field and Select contracts.

## States

The inline alert is a static status surface with no hover or focus states of its own. Its lifecycle is server-driven: the server renders it while the condition holds and removes it when the next response resolves the condition. Tones map to `role` (`alert` for error, `status` otherwise), so the announcement matches the severity.

## Checklist

- [ ] One section-level status only — nothing that should expire on its own.
- [ ] Tone class applied and the role follows the tone.
- [ ] Meaning in text, not color: icon `aria-hidden`, clear title and body.
- [ ] Field errors stay on the fields; the alert summarizes the section.
- [ ] Copy says what happened and what to do next.

## Accessibility

The inline alert is announced with the role that matches its tone, so a failed save interrupts (`alert`) while a saved confirmation announces politely (`status`). Interactive content inside it is a real control with its own accessible name. Color is never the sole signal: text, icon and role all reinforce the tone.

## See also

- [Feedback](/docs/feedback) — the decision matrix this component lives in (FEED-FAIL, FEED-ROW, FEED-PARTIAL, FEED-INLINE).
- [Validation summary](/components/validation-summary) — the navigation landmark that lists field errors for a failed form.
- [Banner](/components/banner) — the full-page sibling for site- or page-level signals.
- [Toast](/components/toast) — the transient, auto-dismissing surface for action results.
- [Choose the right control](/docs/choose-the-right-control) — the cross-component decision.