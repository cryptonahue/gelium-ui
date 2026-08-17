# Callout

Callout is a contextual, ignorable note rendered inline in content — a tip box that explains context, clarifies a term, or points at a useful alternative. Use a callout when a paragraph needs a visual nudge that is not a state signal. It is a real `<aside>` element styled with the `ui-callout` class and needs no component JavaScript.

## Examples

<aside class="ui-callout ui-callout--tip">
  <span class="ui-callout-icon" aria-hidden="true">i</span>
  <div>
    <h3 class="ui-callout-heading">Tip</h3>
    <p class="ui-callout-body">You can match recipients by email or by username — both work the same way.</p>
  </div>
</aside>

<aside class="ui-callout ui-callout--info">
  <span class="ui-callout-icon" aria-hidden="true">i</span>
  <div>
    <h3 class="ui-callout-heading">Note</h3>
    <p class="ui-callout-body">This view updates every 30 seconds. Changes you make elsewhere appear here automatically.</p>
  </div>
</aside>

## Guidance

### When to use

Use a callout when content itself needs a visible aside: a tip, a clarification, or a cross-reference the reader may act on. It renders inline, keeps its `--ui-callout-*` tokens scoped, and carries no role because it is ignorable — there is no state to announce.

### When not to use

Do not use a callout for an error, a warning or a success: those are signals with a role and belong to the [inline alert](/components/inline-alert) or the [banner](/components/banner). Do not use a callout for a message that must be dismissed: callouts have no dismiss because there is no state to clean up. If the note is essential to the task, say it in the content rather than pinning it to a box.

### Usability

- A plain callout (`ui-callout`) uses the muted accent; `ui-callout--info` and `ui-callout--tip` switch the accent bar to the info or secondary color.
- The accent bar is the visual note mark; the heading and body text carry the meaning, never the color alone.
- An optional `ui-button` link turns the note into an action without making it a state signal.

### Accessibility

- The callout carries no role: it is content, read in place like any paragraph.
- The icon is decorative: keep it `aria-hidden` and let the heading or body say what the note is about.
- In forced-colors mode the callout keeps a `CanvasText` boundary so the note stays visible without color.

## Anatomy

- **Callout** — `ui-callout`, the `<aside>` flex row with the scoped `--ui-callout-*` tokens (`padding`, `gap`, `radius`, `icon-size`, `bg`, `fg`).
- **Accent bar** — the 4px inline-start border that marks the note visually (`--ui-callout-accent`).
- **Icon** — `ui-callout-icon`, a decorative `aria-hidden` glyph.
- **Heading** — `ui-callout-heading`, the optional `h3` note title (`--ui-type-label-lg`).
- **Body** — `ui-callout-body`, the note text (`--ui-type-body-sm`).

## Variants

- **Plain** — `ui-callout`, the default surface-container note with a muted accent.
- **Info** — `ui-callout ui-callout--info`, info-tinted background and accent for clarifications.
- **Tip** — `ui-callout ui-callout--tip`, secondary-tinted background and accent for actionable advice.

The variants never communicate state: error, warning and success live on the inline alert and banner.

## Checklist

- [ ] Content is ignorable — the page makes sense if the callout were removed.
- [ ] No role, no tones, no dismiss: the note is not a state signal.
- [ ] Heading or body carries the meaning; the icon is `aria-hidden`.
- [ ] Accent bar renders; forced-colors keeps the box visible.
- [ ] Optional CTA, when present, is a real link.

## Accessibility

The callout is read in place as ordinary content, so it never interrupts with a live region. Text (heading and body) carries the meaning, the icon is decorative, and the forced-colors block keeps the note visible when color is removed.

## See also

- [Feedback](/docs/feedback) — callouts are not feedback; use the matrix to find the right signal component for real state.
- [Inline alert](/components/inline-alert) — the section-level surface for a real status signal.
- [Choose the right control](/docs/choose-the-right-control) — the cross-component decision.