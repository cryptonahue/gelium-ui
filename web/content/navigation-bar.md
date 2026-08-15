# Navigation bar

> **Labs** — experimental. Material Web ships this component as a Labs feature
> with a stability warning; Gelium mirrors that status. The contract may change
> without a major version bump.

Navigation bar is a Material 3 bottom navigation bar: a fixed-height `80px` container that holds three to five equal destinations, one of which is the active page. Use a navigation bar when an app's primary destinations must stay one tap away at the bottom of a phone-sized viewport. Gelium reimplements it over server-rendered HTML — the root is a real `<nav>`, the destinations are real `<a href>` links (the roadmap's "navegación real con links"), and the active destination is derived server-side from the current page, never from JavaScript. No component JavaScript exists.

## Guidance

### When to use

Use a navigation bar when an app's primary destinations must stay one tap away at the bottom of a phone-sized viewport — three to five equal destinations, one of which is the active page.

### When not to use

Do not use a navigation bar when more destinations exist than fit in a bottom bar — a [Navigation drawer](/components/navigation-drawer) scales better. On desktop-sized layouts, prefer a sidebar or top-level navigation instead of a bottom bar.

### Usability

- Keep the bar to three to five destinations; the active one is derived server-side from the current page.
- The root is a real `<nav>` with real `<a href>` destinations; the active destination carries `aria-current="page"`.
- The `ui-nav-bar--hide-inactive-labels` modifier collapses inactive destinations to icon-only.

### Accessibility

- Keep the native elements: `<nav>` carries the landmark, `<ul>`/`<li>` carry list semantics, and destinations are real anchors.
- The selected destination carries `aria-current="page"`; state is never communicated by color alone.
- The icon slot is `aria-hidden`; the label text supplies the name. In icon-only mode the label stays in the DOM.
- In forced-colors mode the active destination becomes `HighlightText` and the indicator pill is outlined with `Highlight`.

See [Choose the right control](/docs/choose-the-right-control) for the cross-component decision.
## Anatomy

```html
<nav class="ui-nav-bar" aria-label="Primary">
  <ul class="ui-nav-bar-list">
    <li class="ui-nav-bar-item">
      <a class="ui-nav-bar-destination ui-nav-bar-destination--active"
         href="/components/navigation-bar" aria-current="page">
        <span class="ui-nav-bar-icon" aria-hidden="true">
          <span class="ui-nav-bar-indicator"></span>
          [inactive glyph][active glyph]
        </span>
        <span class="ui-nav-bar-label">Navigation bar</span>
      </a>
    </li>
  </ul>
</nav>
```

- **Container** — `ui-nav-bar`, a real `<nav>` `80px` tall on the container
  surface with the Material level-2 elevation.
- **List** — `ui-nav-bar-list`, a real `<ul>` of equal `ui-nav-bar-item` flex
  cells (each at least `48px` square).
- **Destination** — `ui-nav-bar-destination`, a real `<a href>` laid out as a
  flex column. The icon slot (`ui-nav-bar-icon`) is a `64x32` box holding the
  `32px` `ui-nav-bar-indicator` pill behind the `24px` glyph; the `12px`
  `ui-nav-bar-label` sits below it. Material's active/inactive icon pair maps to
  two glyphs in the slot, and only the active glyph is painted when the
  destination is selected.
- **Active indicator** — `ui-nav-bar-indicator`, the `32px` full-radius pill
  (`secondary-container` color) that fades in behind the icon of the selected
  destination.
- **Badge** — an optional badge reuses the existing `.ui-badge` primitive
  (composed, not reinvented) and is anchored to the icon's top-end corner.

## Variants

- **Standard** — every destination shows its icon and label; the active
  destination additionally shows the indicator pill, the active icon/label
  colors, and the prominent `700` label weight.
- **Hide inactive labels** — with the `ui-nav-bar--hide-inactive-labels`
  modifier on the root, inactive destinations collapse to icon-only (`0`
  height, transparent label) while the active destination keeps its label and
  indicator. This mirrors the upstream `hideInactiveLabels` option.

## States

`rest`, `hover`, `focus-visible`, `active`/`pressed`, and `selected` (the active
destination). Hover, focus, and press paint the Material state layer (a
full-bleed `::before` at the shared `--ui-state-*` opacities) and shift the
inactive icon toward on-surface; `:focus-visible` shows the Gelium focus ring.
Because destinations are native links, keyboard activation and tab order come
for free and focus never changes geometry.

## Accessibility

- Keep the native elements: `<nav>` carries the landmark, `<ul>`/`<li>` carry
  list semantics, and destinations are real anchors, so the accessible name is
  the visible label and activation is native.
- The selected destination carries `aria-current="page"`; state is never
  communicated by color alone (the indicator pill is the selected marker, focus
  shows a ring, the label weight changes).
- The icon slot is `aria-hidden`; the label text supplies the name. In
  icon-only mode the label stays in the DOM (only visually collapsed), so
  assistive technology still reads every destination.
- Decorative SVGs are `aria-hidden` and `focusable="false"`.
- In forced-colors mode the container keeps a `CanvasText` boundary, the active
  destination becomes `HighlightText`, and the indicator pill is outlined with
  `Highlight` so selection survives without color.

## No-JS behavior

Navigation is plain HTTP: each destination is a real `<a href>` and clicking it
navigates normally with scripting disabled. The active state is fixed at render
time by the server (the destination whose href matches the request path), so
there is nothing for JavaScript to do.

## Divergences from Material Web

| Area | Material Web | Gelium |
| --- | --- | --- |
| Root semantics | `<div role="tablist">` with roving `tabindex` and JS activation | real `<nav>` with `<a href>` destinations |
| Active state | `activeIndex` property driven by clicks and JS | `aria-current="page"` derived server-side from the current page |
| Arrow/Home/End keyboard | custom keydown handler | deferred (roadmap: requires a demonstrated platform gap); native tab order is sufficient for links |
| Badge | `<md-badge>` inside an `aria-hidden` icon slot | existing `.ui-badge` composed into the same slot |
| Focus ring | inward ring shape `corner-small` | Gelium `:focus-visible` outline (consistent with the rest of Gelium) |

## Compatibility

The component uses only `display:flex`, `position`, `opacity`, `::before`, and
media queries for `prefers-reduced-motion` and `forced-colors` — all baseline
CSS. No custom properties beyond the `--ui-*` token family.

## Trust boundary

The inline SVG glyphs are internal constants embedded in `navigation_bar.go`;
they are `aria-hidden`, unfocusable, and colored by `currentColor`. User input
never reaches them, and the destination hrefs and labels are Go values escaped
by `html/template`.

## Visual checklist

- Bar is `80px` tall and spans its container; destinations divide it evenly.
- The `32px` pill indicator sits behind the `24px` icon on the active
  destination only.
- Active icon/label use `on-secondary-container`/`on-surface`; inactive use
  `on-surface-variant`; the active label weight is `700`.
- Hover/pressed paint the state layer without moving focus or geometry.
- Icon-only (hide inactive labels) keeps the label as the accessible name.
- Badges anchor to the icon's top-end corner and reuse `.ui-badge`.
