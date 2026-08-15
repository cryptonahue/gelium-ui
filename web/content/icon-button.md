# Icon button

Icon button is a compact, single-tap action control built from a native HTML `<button>` (for actions) or `<a>` (for navigation). Use an icon button when a toolbar or surface needs a small action whose meaning is carried by a glyph plus an accessible name — edit, delete, refresh — and the tap target must stay touch-friendly. It renders a Material Web–compatible circular touch target with a trusted inline SVG glyph.

## Guidance

### When to use

Use an icon button when a toolbar or surface needs a small action whose meaning is carried by a glyph plus an accessible name — edit, delete, refresh — and the tap target must stay touch-friendly.

### When not to use

Do not use an icon button when the action needs a visible label to be understood — a [Button](/components/button) with text is clearer. If several icon actions pile up on one surface, a [Menu](/components/menu) can hold them without crowding the toolbar.

### Usability

- Four variants match Material Web anatomy: standard, filled, filled-tonal, and outlined.
- Provide an `Href` to render an `<a>` for navigation; otherwise the control is a native `<button>`.
- Toggle icon buttons reflect their state with `aria-pressed` and may expose a state-specific name.

### Accessibility

- An icon button is never icon-only without a name: use the visible label, or supply an `aria-label` — the template refuses to emit an empty one.
- Keyboard focus uses an exterior 3 px `:focus-visible` outline with a 2 px offset.
- Disabled and pressed states do not rely on color alone; forced-colors and reduced-motion themes are covered by the shared build.

See [Choose the right control](/docs/choose-the-right-control) for the cross-component decision.
## Variants

The live preview below uses the same `icon-button` Go template as this documentation layout. Four variants are available, matching Material Web anatomy:

- **standard** — transparent, lowest emphasis.
- **filled** — primary container, high emphasis.
- **filled-tonal** — secondary container, medium emphasis.
- **outlined** — outlined circle, medium emphasis.

## Accessible name (required)

An icon button is never icon-only without a name. Use the visible label on a standard icon button, or supply an `aria-label` for icon-only buttons. The template refuses to emit an empty `aria-label=""`; an unlabelled instance still exposes the glyph only if the caller provides a real name.

## Toggle

Toggle icon buttons switch between selected and unselected states, reflecting the toggle with `aria-pressed`. A selected toggle may render a distinct selected glyph and, when `AriaLabelSelected` is given, expose a state-specific accessible name. Disabled toggle buttons lose the activation path (`disabled` on the native button; an `<a>` in inactive mode omits `href`, leaves the tab order, and carries `aria-disabled`).

## Navigation

Provide an `Href` to render an `<a>` for navigation. For un-navigable (disabled) href controls, the link omits the destination, leaves the tab order, and exposes its inactive state — no activation path remains.

## Trust boundary

`IconSVG` and `SelectedIcon` are per-instance inline SVG slots typed as `template.HTML`. Pass only trusted, internal markup — never user input. Icon glyphs are decorative by contract and must include `aria-hidden="true"` and `focusable="false"`; the accessible name comes from the visible label or `aria-label`.

## Accessibility

Keyboard focus uses an exterior 3 px `:focus-visible` outline with a 2 px offset. Disabled and pressed states do not rely on color alone; forced-colors and reduced-motion themes are covered by the shared build.
