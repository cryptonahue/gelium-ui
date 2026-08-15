# Divider

Divider is a thin line that groups or separates related content. Use a divider when adjacent content needs a visual separation that is not a container boundary — between list rows, toolbar groups, or stacked sections. It is a native `hr` element styled with the `ui-divider` class, so it needs no component JavaScript at all.

## Guidance

### When to use

Use a divider when adjacent content needs a visual separation that is not a container boundary — between list rows, toolbar groups, or stacked sections.

### When not to use

Do not divide everything: spacing alone is often enough, and over-dividing adds noise. If the groups are truly distinct surfaces, a [Card](/components/card) or a container boundary communicates the separation better than a line.

### Usability

- A plain `<hr class="ui-divider">` spans the full width of its container.
- Use the inset modifiers to shorten the line: `ui-divider-inset`, `ui-divider-inset-start`, `ui-divider-inset-end`.
- The inset distance is exactly `1rem` with logical properties, so it mirrors correctly in right-to-left layouts.

### Accessibility

- A divider is decorative by default: it does not announce anything to assistive technology.
- When the line represents a real thematic or section break, add `role="separator"` manually.
- In forced-colors mode the divider keeps a `CanvasText` paint so the grouping stays visible without color.

See [Choose the right control](/docs/choose-the-right-control) for the cross-component decision.
## Anatomy

The visible line uses the theme tokens:

- `--ui-divider-color` — the divider paint, mapped to `--ui-color-border` so it follows light and dark schemes automatically
- `--ui-divider-thickness` — 1 px by default

The divider is `display: block`, has no border or margin, and is clipped to its content box so shortened variants do not fade the paint.

## Variants

A plain `<hr class="ui-divider">` spans the full width of its container.

Use the inset modifiers to shorten the line from one or both edges:

- `ui-divider-inset` — 1rem from both edges
- `ui-divider-inset-start` — 1rem from the inline-start edge
- `ui-divider-inset-end` — 1rem from the inline-end edge

The inset distance is exactly `1rem` and uses logical properties, so it mirrors correctly in right-to-left layouts.

## Accessibility

A divider is decorative by default: it does not announce anything to assistive technology. When the line represents a real thematic or section break, add `role="separator"` manually so screen readers treat it as a separator between groups. The divider itself carries a tall color-contrast-friendly paint in forced-colors mode (`CanvasText`), keeping the grouping visible when color is removed.