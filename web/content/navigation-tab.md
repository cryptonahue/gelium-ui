# Navigation tab

> **Labs** — experimental. Material Web ships this component as a Labs feature
> with a stability warning; Loom mirrors that status. The contract may change
> without a major version bump.

A Material 3 navigation tab: the individual destination of a navigation bar,
reimplemented over server-rendered HTML. The root is a real `<a href>` link —
the roadmap's **"link semántico, no tab falso"** contract. There is no
`role="tab"`, no `role="tablist"`, no roving focus and no component JavaScript:
the active tab is derived server-side from the current page and marked with
`aria-current="page"`.

Because the tab is the same destination contract the delivered navigation bar
uses, the demo composes the existing `.ui-nav-bar` for its in-bar variant and
reuses the existing `.ui-badge` for badge destinations — nothing is reinvented.

## Anatomy

```html
<a class="ui-nav-tab ui-nav-tab--active"
   href="/components/navigation-tab" aria-current="page">
  <span class="ui-nav-tab-icon" aria-hidden="true">
    <span class="ui-nav-tab-indicator"></span>
    [inactive glyph][active glyph]
  </span>
  <span class="ui-nav-tab-label">Navigation tab</span>
</a>
```

- **Tab** — `ui-nav-tab`, a real `<a href>` laid out as a flex column with a
  `48px` minimum hit area (`8px 0 12px` padding). The icon slot
  (`ui-nav-tab-icon`) is a `64x32` box holding the `32px` `ui-nav-tab-indicator`
  pill behind the `24px` glyph; the `12px` `ui-nav-tab-label` sits below it.
- **Active indicator** — `ui-nav-tab-indicator`, the `32px` full-radius pill
  (`secondary-container` color) that fades in behind the icon of the selected
  tab.
- **Glyph swap** — Material's active/inactive icon pair maps to two glyphs in
  the icon slot, and only the active glyph is painted when the tab is selected.
  When both slots carry the same glyph, Loom renders a single copy (there is
  nothing to swap).
- **Badge** — an optional badge reuses the existing `.ui-badge` primitive
  (composed, not reinvented) and is anchored to the icon's top-end corner.
- **In-bar composition** — inside the delivered `.ui-nav-bar`, each destination
  is the same `ui-nav-tab` contract (`.ui-nav-bar-list` / `.ui-nav-bar-item`
  wrap the tabs).

## Variants

- **Standard** — every tab shows its icon and label; the active tab additionally
  shows the indicator pill, the active icon/label colors, and the prominent
  `700` label weight.
- **Hide inactive labels** — with the `ui-nav-tab--hide-inactive-label` modifier,
  inactive tabs collapse to icon-only (`0` height, transparent label) while the
  active tab keeps its label and indicator. This mirrors the upstream
  `hideInactiveLabel` option.
- **Badges** — a dot or count badge reuses `.ui-badge` on the icon slot.
- **Inside the navigation bar** — the tabs composed into the delivered
  `.ui-nav-bar`, proving the shared destination contract.

## States

`rest`, `hover`, `focus-visible`, `active`/`pressed`, and `selected` (the active
tab). Hover, focus, and press paint the Material state layer (a full-bleed
`::before` at the shared `--ui-state-*` opacities) and shift the inactive icon
toward on-surface; `:focus-visible` shows the Loom focus ring. Because tabs are
native links, keyboard activation and tab order come for free and focus never
changes geometry.

## Accessibility

- Keep the native element: each tab is a real anchor, so the accessible name is
  the visible label and activation is native. No fake-tab ARIA.
- The selected tab carries `aria-current="page"`; state is never communicated by
  color alone (the indicator pill is the selected marker, focus shows a ring,
  the label weight changes).
- The icon slot is `aria-hidden`; the label text supplies the name. In icon-only
  mode the label stays in the DOM (only visually collapsed), so assistive
  technology still reads every tab.
- Decorative SVGs are `aria-hidden` and `focusable="false"`.
- In forced-colors mode the active tab becomes `HighlightText`, and the
  indicator pill is outlined with `Highlight` so selection survives without
  color.

## No-JS behavior

Navigation is plain HTTP: each tab is a real `<a href>` and clicking it
navigates normally with scripting disabled. The active state is fixed at render
time by the server (the tab whose href matches the request path), so there is
nothing for JavaScript to do.

## Divergences from Material Web

| Area | Material Web | Loom |
| --- | --- | --- |
| Root semantics | `<button role="tab">` with roving `tabindex` and JS activation | real `<a href>` link (roadmap: "link semántico, no tab falso") |
| Active state | `active` property driven by the bar's `activeIndex` | `aria-current="page"` derived server-side from the current page |
| Arrow/Home/End keyboard | custom keydown handler in the bar | deferred (roadmap: requires a demonstrated platform gap); native tab order is sufficient for links |
| Badge | `<md-badge>` inside an `aria-hidden` icon slot | existing `.ui-badge` composed into the same slot |
| Focus ring | inward ring shape `corner-small` | Loom `:focus-visible` outline (consistent with the rest of Loom) |

## Compatibility

The component uses only `display:flex`, `position`, `opacity`, `::before`, and
media queries for `prefers-reduced-motion` and `forced-colors` — all baseline
CSS. No custom properties beyond the `--ui-*` token family.

## Trust boundary

The inline SVG glyphs are internal constants embedded in `navigation_tab.go`;
they are `aria-hidden`, unfocusable, and colored by `currentColor`. User input
never reaches them, and the tab hrefs and labels are Go values escaped by
`html/template`.

## Visual checklist

- Each tab is a real link with a `48px` minimum hit area; the icon slot is
  `64x32` and holds the `24px` glyph.
- The `32px` pill indicator sits behind the `24px` icon on the active tab only.
- Active icon/label use `on-secondary-container`/`on-surface`; inactive use
  `on-surface-variant`; the active label weight is `700`.
- Hover/pressed paint the state layer without moving focus or geometry.
- Icon-only (hide inactive labels) keeps the label as the accessible name.
- Badges anchor to the icon's top-end corner and reuse `.ui-badge`.
- The in-bar variant composes the delivered `.ui-nav-bar` unchanged.
