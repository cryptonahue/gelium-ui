# Navigation drawer

> **Labs** — experimental. Material Web ships this component as a Labs feature
> with a stability warning; Gelium mirrors that status. The contract may change
> without a major version bump.

Navigation drawer is a Material 3 navigation drawer with two variants: a **modal** drawer that opens over a scrim and a **standard** (permanent) drawer embedded in the layout. Use a drawer when an app has more destinations than fit in a bottom bar and the navigation must stay server-rendered and no-JS. Gelium reimplements it over server-rendered HTML — the standard variant is a real `<nav>` in the layout, the modal variant is a native `<dialog>`, and every destination is a real `<a href>` whose active state is derived server-side from the current page (the roadmap's "navegación real con links"). No component JavaScript exists.

## Anatomy

```html
<nav class="ui-navigation-drawer ui-navigation-drawer--standard" aria-label="Primary">
  <ul class="ui-navigation-drawer-list">
    <li class="ui-navigation-drawer-item">
      <a class="ui-navigation-drawer-destination ui-navigation-drawer-destination--active"
         href="/components/navigation-drawer" aria-current="page">
        <span class="ui-navigation-drawer-indicator" aria-hidden="true"></span>
        <span class="ui-navigation-drawer-glyph" aria-hidden="true">[24px glyph][optional badge]</span>
        <span class="ui-navigation-drawer-label">Navigation drawer</span>
      </a>
    </li>
  </ul>
</nav>
```

- **Container** — `ui-navigation-drawer`, `360px` wide and full height on the
  Material surface, shaped `corner-large-end` (rounded on the end side, so the
  start edge sits flush against the screen).
- **List** — `ui-navigation-drawer-list`, a real `<ul>` of 56px-tall
  `ui-navigation-drawer-item` rows.
- **Destination** — `ui-navigation-drawer-destination`, a real `<a href>`
  laying out the `24px` icon and the `label-large` label with a `12px` gap.
  Unlike the navigation bar, the drawer keeps a single glyph per destination
  and recolors it when selected (Material does not swap icons in a drawer).
- **Active indicator** — `ui-navigation-drawer-indicator`, the full-width
  `336px × 56px` `secondary-container` pill (`corner-full`) behind the selected
  destination, mirroring the upstream `active-indicator.width/height/shape`
  tokens.
- **Badge** — an optional badge reuses the existing `.ui-badge` primitive
  (composed, not reinvented) and is anchored to the icon's top-end corner.

## Variants

- **Standard** — the permanent drawer: a real `<nav>` embedded in the layout,
  always visible. This is the roadmap's "variante permanente como `<nav>` en
  layout". The active destination carries `aria-current="page"`.
- **Modal** — the roadmap's "variante modal sobre `<dialog>`": a native
  `<dialog class="ui-navigation-drawer ui-navigation-drawer--modal">` with a
  scrim (its `::backdrop`), opened by a trigger button using the native invoker
  command `command="show-modal"` + `commandfor` (the Invoker Commands API —
  declarative `command`/`commandfor` dialog control). Native dialog behavior gives
  the top layer, focus move into the drawer, and Escape-to-close for free.
  `closedby="any"` adds scrim (light) dismiss in supporting browsers only;
  the explicit Escape and focus behavior remain in compatible browsers.

## States

`rest`, `hover`, `focus-visible`, `active`/`pressed`, and `selected` (the active
destination). Hover, focus, and press paint the Material state layer (a
full-bleed `::before` at the shared `--ui-state-*` opacities) and shift the
inactive icon/label toward on-surface; the selected destination uses
`on-secondary-container` for its state layer. `:focus-visible` shows the Gelium
focus ring. Because destinations are native links, keyboard activation and tab
order come for free and focus never changes geometry.

## Accessibility

- Keep the native elements: the standard variant is a real `<nav>` landmark,
  `<ul>`/`<li>` carry list semantics, and destinations are real anchors, so the
  accessible name is the visible label and activation is native. The modal
  variant is a native `<dialog>`: the browser manages the top layer, focus, and
  Escape, and the drawer carries `aria-label`.
- The selected destination carries `aria-current="page"`; state is never
  communicated by color alone (the indicator pill is the selected marker, focus
  shows a ring, the label weight changes).
- The icon slot is `aria-hidden`; the label text supplies the name. Badges live
  inside the aria-hidden slot (decorative or a visual count only).
- Decorative SVGs are `aria-hidden` and `focusable="false"`.
- In forced-colors mode the drawer keeps a `CanvasText` boundary, the active
  destination becomes `HighlightText`, and the indicator pill is outlined with
  `Highlight` so selection survives without color.

## How does the drawer open and close without JavaScript?

Navigation is plain HTTP: each destination is a real `<a href>` and clicking it
navigates normally with scripting disabled; the active state is fixed at render
time by the server. The modal drawer opens through the native invoker command
`command="show-modal"` (declarative, no component JavaScript) exactly like the
shipped Dialog component; in browsers without the Invoker Commands API (the
native `command`/`commandfor` mechanism for declarative dialog control) the
trigger does nothing, so consumers supporting them need a server-rendered
fallback or adapter — the same documented gap as Dialog.

## HTMX

The roadmap reserves HTMX for remote content, not for basic semantics. The
drawer has no remote content today, so it ships with no `hx-*` attributes; the
destinations are plain links.

## Divergences from Material Web

| Area | Material Web | Gelium |
| --- | --- | --- |
| Root semantics (standard) | `<div role="dialog">` inside a host with the `opened` property | real `<nav>` embedded in the layout |
| Root semantics (modal) | custom host + separate scrim `<div>`, JS `opened` toggling | native `<dialog>` + `::backdrop` scrim, declarative invoker trigger |
| Active state | `aria-expanded`/`aria-hidden` on a wrapper | `aria-current="page"` derived server-side from the current page |
| Icon | two `slot` copies (active/inactive), consumer-swapped | single `24px` glyph per destination, recolored by state |
| Active indicator | `active-indicator.width: 336px`, `height: 56px`, `corner-full` | same geometry via the `ui-navigation-drawer-indicator` pill |
| Focus ring | inward focus indicator | Gelium `:focus-visible` outline (consistent with the rest of Gelium) |
| Scrim | deprecated `scrim` tokens (neutral-variant20) | existing `--ui-dialog-scrim` (32% on black) reused from the Dialog contract |
| Modal container | `surface-container-low`, level-1 elevation | `--ui-dialog-container` (carries `surface-container-low`) + `--ui-shadow-1` |
| `pivot` start/end | property on both variants | not ported — no demonstrated consumer need; logical properties anchor start |
| `navigation-drawer-changed` event | fires 250ms after open/close | not ported — server-rendered links don't need a JS event |

## Compatibility

The component uses only baseline CSS (`display:flex`, `position`, `opacity`,
`::before`, logical properties) plus the modern `@starting-style` and
`overlay`/`display` discrete transitions for modal motion — the same
progressive-enhancement pattern as Dialog, where motion may be instant in
browsers without them. The modal trigger relies on the native Invoker Commands
API (`command`/`commandfor`, recent Baseline Low features, same as Dialog); the
permanent variant needs nothing beyond baseline HTML.

## Trust boundary

The inline SVG glyphs are internal constants embedded in `navigation_drawer.go`;
they are `aria-hidden`, unfocusable, and colored by `currentColor`. User input
never reaches them, and the destination hrefs and labels are Go values escaped
by `html/template`. The trigger and dialog id are Go values emitted as plain
attributes.

## Visual checklist

- The drawer is `360px` wide, full height, surface-colored, rounded on the end
  side.
- Destinations are `56px` tall with a `24px` icon, a `12px` gap, and a
  `label-large` label; the icon sits `12px` from the start edge.
- The active destination shows the full-width `336px × 56px` `secondary-container`
  pill; its icon/label use `on-secondary-container` and the label weight is `700`.
- Hover/pressed paint the state layer without moving focus or geometry.
- The modal drawer slides in from the start edge over a dimmed scrim and closes
  with Escape or the native close affordances.
- Badges anchor to the icon's top-end corner and reuse `.ui-badge`.
