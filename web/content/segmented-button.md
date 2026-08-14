# Segmented buttons

> **Labs** — experimental component. Breaking changes may happen outside major-version bumps, mirroring Material Web's own Labs warning.

Segmented buttons display a set of options or actions in a single, mutually exclusive row. Gelium reimplements the Material 3 outlined segmented button set over **server-rendered HTML with native form semantics**: single-select groups use `input[type="radio"]`, multi-select groups use `input[type="checkbox"]`, and non-selection actions use `button[type="button"]`. Selection state is derived from the native `:checked` pseudo-class — there is **no component JavaScript** and the checked values submit through a normal `<form>`.

## Anatomy

- **Set** — `ui-segmented-button-set`, the pill container (40px high, full corner radius, 1px outline). For selection sets the root is a `<fieldset>` with a visually hidden `<legend>` that names the group; for action sets the root is an accessible `role="group"` with an `aria-label` (the roadmap's "fieldset o grupo accesible").
- **Segment** — `ui-segmented-button`, the individual radio/checkbox/button cell. First and last segments cap the pill's outer corners; segments between them carry 1px outline dividers.
- **Input** — the native control inside each segment. A shared `name` makes a radio group single-select and a checkbox group multi-select for free.
- **Graphic** — `ui-segmented-button-graphic`, the leading slot that expands (0 → 18px) to reveal the **checkmark** when the segment is `:checked`.
- **Checkmark** — `ui-segmented-button-checkmark`, the Material check path drawn in `currentColor` with a 2px stroke.
- **Icon** — `ui-segmented-button-icon`, the optional leading glyph. In icon+label segments it is replaced by the checkmark when selected; in icon-only segments it stays visible and the checkmark appears beside it.
- **Label** — `ui-segmented-button-label`, the label-large text that supplies the accessible name (use `aria-label` on the input for icon-only segments).

## Variants and states

- **Single select** — a radio group; one `checked` segment, arrow keys move selection natively.
- **Multi select** — a checkbox group; several `checked` segments, Space toggles each natively.
- **Action set** — `button[type="button"]` segments with no selection state and no checkmark.
- Each segment can be **icon-only**, **label-only**, or **icon + label**.
- States: rest, hover, focus-visible (3px focus ring, no geometry shift), active/pressed (state layer), selected (`:checked`), and disabled (native `disabled`, dimmed content, `not-allowed` cursor).

## Accessibility

- Native controls: radios and checkboxes keep their role, name, checked state, and keyboard behavior at no cost; a radio group is announced as a group with the fieldset legend as its name.
- The `legend` is visually hidden with a standard clip pattern — still in the accessibility tree, never rendered.
- Never color-only: the selected state always carries the checkmark, disabled is announced by the platform, and forced-colors keeps `CanvasText`/`GrayText` boundaries.
- Decorative glyphs are `aria-hidden`; the visible label (or the input's `aria-label`) supplies the accessible name.

## No-JS behavior

The demo wraps the two selection sets in a real `<form method="get" action="/components/segmented-button">`. Submitting navigates to `?transport=drive&formatting=bold&formatting=italic` (whatever is checked) with JavaScript disabled. The action set needs no form at all.

## Compatibility

`<fieldset>`/`<legend>`, native radio/checkbox, `:checked`, and `:focus-visible` are Baseline widely available. `:has()` (used to style the segment from its native control) is Baseline as of late 2023. The checkmark reveal is a CSS transition, so reduced motion disables it; nothing depends on it. No browser is blocked from the no-JS selection flow.

## Trust boundary

The demo markup is static and trusted: `html/template` renders it verbatim, the inline SVGs are trusted internal glyphs (`aria-hidden`, `focusable="false"`), and there is no user-supplied content.

## Divergences from Material Web

- **Selection semantics** — upstream `md-outlined-segmented-button` uses a `<button aria-pressed>` plus component JS in the set to manage selection. Gelium uses native radio/checkbox/button semantics with `:checked`, per the roadmap's "Single/multi select debe preferir radios/checkboxes sin JS". Consequence: radio groups get native arrow-key navigation and form submission instead of per-button tab stops.
- **Geometry** — upstream tokens define the set and segment outer corners as `corner-full` (9999px, a pill); Gelium matches that. The upstream checkmark draw-in animation is a multi-phase keyframe sequence; Gelium reveals it as a stroke-dashoffset transition driven by `:checked`.

## Visual checklist

- [ ] 40px set with full-corner pill and 1px outline
- [ ] First/last segments cap the pill; middle segments square
- [ ] Selected segment: secondary-container fill + checkmark + on-container text
- [ ] Icon+label segment swaps icon for checkmark when selected
- [ ] Icon-only segment keeps the icon and adds the checkmark when selected
- [ ] Hover/focus/pressed state layers on every segment
- [ ] Focus ring without layout shift; disabled dims without color-only
- [ ] Light/dark, narrow/wide, RTL, reduced motion, forced colors
