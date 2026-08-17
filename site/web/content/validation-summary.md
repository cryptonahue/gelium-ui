# Validation summary

Validation summary is a server-rendered, form-level list of validation errors that links to the offending fields — native anchors, no JavaScript. Use a validation summary when a form fails validation, so the user gets one navigation landmark that names each problem and jumps to the field that needs attention. It renders as a real `div` styled with the `ui-validation-summary` class with `role="alert"`.

## Examples

<div class="component-preview">
  <div class="ui-validation-summary" role="alert">
    <h2 class="ui-validation-summary-title">Check the following</h2>
    <ul class="ui-validation-summary-list">
      <li class="ui-validation-summary-item"><a href="#email">Email is required. Enter the address you signed up with.</a></li>
      <li class="ui-validation-summary-item"><a href="#password">Password must be at least 8 characters.</a></li>
      <li class="ui-validation-summary-item"><a href="#region">Select a region.</a></li>
    </ul>
  </div>
</div>

## Guidance

### When to use

Use a validation summary when a form returns validation errors: it is the navigation landmark for the whole form while each field keeps its own `aria-invalid` and `aria-describedby` error message. The summary links to the fields with native anchors, so a keyboard or screen-reader user can jump straight to the first problem (see [Feedback](/docs/feedback), FEED-VAL). Validation errors arrive via the server contract (`HTTP 422` + `X-Gelium-Validation`), not a toast.

### When not to use

Do not use a validation summary for a single field: the field itself carries the message. Do not use a toast for validation — a transient message disappears before the user can act (see [Feedback](/docs/feedback)). And do not redirect or reload the page in a way that loses what the user typed; the server contract preserves the submitted values and returns them with the errors.

### Usability

- The title is short ("Check the following") and the list names each problem tersely.
- Every item is a real link to the field's anchor, so the user can act on each error.
- The summary composes with — and never replaces — the field-level error messages.

### Accessibility

- The root is `role="alert"`, so assistive technology announces the failed submission as a whole.
- The heading level is configurable (`HeadingLevel`) so the summary nests correctly in the page outline.
- Each item is a real link with a visible destination; focus follows the link to the field.

## Anatomy

- **Summary** — `ui-validation-summary`, the flex column with the scoped `--ui-validation-summary-*` tokens (`padding`, `gap`, `radius`, `bg`, `fg`, `title-color`, `item-color`).
- **Title** — `ui-validation-summary-title`, the heading (`h2` by default, level configurable) with `--ui-type-title-md`.
- **List** — `ui-validation-summary-list`, the unstyled list of errors.
- **Items** — `ui-validation-summary-item`, one `<a href>` per offending field carrying the message text.

The summary always renders `role="alert"` because it only ever reports validation errors — there are no tones.

## States

The validation summary has a single state: an active validation failure with `role="alert"`. There are no tones or variants — it fixes to the danger container and announces every time the server returns validation errors.

## Checklist

- [ ] Server returned validation errors (`422` + `X-Gelium-Validation`) — never a toast for validation.
- [ ] Title + linked list of every problem; each item is a real anchor.
- [ ] Each field also carries its own `aria-invalid` + `aria-describedby` message.
- [ ] Submitted values preserved across the validation round trip.
- [ ] Heading level nests in the page outline; `role="alert"` present.
- [ ] Forced-colors output keeps titles and links readable.

## Accessibility

The validation summary announces the failed submission with `role="alert"` and turns every error into a native anchor, so keyboard and screen-reader users can move to each field without re-reading the whole form. It composes with the field-level messages instead of replacing them, and the server contract preserves entered values so a validation round trip never costs the user their input.

## See also

- [Feedback](/docs/feedback) — the decision matrix entry FEED-VAL (validation summary + inline field errors).
- [Inline alert](/components/inline-alert) — the section-level signal for non-validation failures like a failed save.
- [Text field](/components/text-field) — the field-level error message contract this summary composes with.
- [Forms](/docs/forms) — the end-to-end form pattern that pairs them.
- [Choose the right control](/docs/choose-the-right-control) — the cross-component decision.