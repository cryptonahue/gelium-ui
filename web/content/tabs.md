# Tabs

Tabs organize groups of related content at the same level of hierarchy. Gelium reimplements the Material 3 tab contract as server-rendered navigation links: each tab is a real `<a href>` to its own page or section, and the selected tab is marked by the handler from the current URL — there is no component JavaScript, no `role="tablist"` and no roving focus.

## Anatomy

A tab bar is a `<nav>` wrapping a `<ul>` of `<li>` items, each containing an `<a class="ui-tab">`. The anchor carries the icon (optional), the label (optional), and the active indicator. Keep the order `icon`, `label`, `indicator` exactly.

```html
<nav class="ui-tabs" aria-label="Primary tabs example">
  <ul class="ui-tabs-list">
    <li class="ui-tabs-item">
      <a class="ui-tab ui-tab-primary ui-tab-stacked" href="/components/tabs?tab=photos" aria-current="page">
        <span class="ui-tab-icon" aria-hidden="true">…</span>
        <span class="ui-tab-label">Photos</span>
        <span class="ui-tab-indicator" aria-hidden="true"></span>
      </a>
    </li>
  </ul>
</nav>
```

- **Bar** — `ui-tabs` is a flex column that fills the width and draws the Material 1px divider below the tabs. The `ui-tabs-list` aligns items to the bottom edge so mixed-height tabs sit on the same baseline.
- **Tab** — `ui-tab` is the real link. Its height is `48px` (`--ui-tab-height`), horizontal padding is `16px`, and the label uses the Material `title-small` typescale (`500 .875rem/1.25rem`, tracking `0.00625rem`). The inactive label and icon paint `--ui-color-fg-muted` (on-surface-variant).
- **Icon** — `ui-tab-icon` is a `24px` decorative slot; the visible label (or an `aria-label` on icon-only tabs) supplies the accessible name.
- **Indicator** — `ui-tab-indicator` is the `3px` active indicator in `--ui-color-primary` (primary tabs) or `2px` (secondary tabs). It is hidden unless the tab is selected.
- **Variant** — `ui-tab-primary` and `ui-tab-secondary` set the active colors: primary tabs use `--ui-color-primary` for the selected label/icon, secondary tabs use `--ui-color-fg`.
- **Stacked** — primary tabs with both icon and label add `ui-tab-stacked`: the icon stacks above the label (`column`, gap `2px`) and the tab grows to `64px`, matching the Material `with-icon-and-label-text-container-height`.

## Variants

Material defines two tab variants, and each supports icon-only, label-only, and icon+label content:

- **Primary tabs** — placed at the top of the content pane. With an icon and label the content stacks vertically (`ui-tab-stacked`, `64px`); label-only and icon-only tabs stay `48px`.
- **Secondary tabs** — used within a content area. The icon always renders inline with the label (row layout, `48px`).

Use the same variant for every tab in a bar, and never mix primary and secondary tabs in one bar.

## States

The selected tab is decided by the server, so the bar always shows its state without JavaScript:

- **Rest** — inactive label/icon in on-surface-variant.
- **Hover** — state layer at `--ui-state-hover-opacity` (`.08`).
- **Focus-visible** — state layer at `--ui-state-focus-opacity` (`.10`) plus the standard `--ui-focus-thickness` ring in `--ui-color-focus-ring`. The ring does not change geometry or shift layout.
- **Active/pressed** — state layer at `--ui-state-pressed-opacity` (`.10`).
- **Selected** — `aria-current="page"` turns the label/icon to the active color and shows the indicator.

The state layer color follows the Material split: inactive tabs paint on-surface (`--ui-color-fg`), selected primary tabs paint primary (`--ui-color-primary`), selected secondary tabs keep on-surface.

## Selection is server-driven

Each tab is a real link to its own section. On this page the primary tabs link to `?tab=photos|videos|music` and the secondary tabs to `?sub=travel|hotel|activities`. The handler validates the parameter against a closed vocabulary and marks the matching tab with `aria-current="page"`; unknown values fall back to the first tab. Clicking a tab is a normal GET navigation — the full page is re-rendered with the new selection and the matching panel content. Nothing requires client-side state.

## Accessibility

- **Plain navigation semantics** — the bar is a `<nav aria-label>` with `<ul>`/`<li>`/`<a>`, and the selected tab uses `aria-current="page"`. This is the simplest semantics that satisfies accessibility for the link pattern.
- **No `role="tablist"`** — the full tablist keyboard contract (roving tabindex, arrow-key focus, Home/End) cannot be satisfied without JavaScript, so Gelium does not claim it. Native link keyboard behavior (Tab to move, Enter to activate) is complete and announced correctly.
- **No roving focus** — arrow-key roving focus would require a demonstrated platform gap; the roadmap only allows it when the gap is real, and for real navigation links there is none.
- **Icon-only tabs** must carry an `aria-label`; icon and indicator are decorative and `aria-hidden`.
- In forced-colors mode the indicator repaints as `CanvasText`, selected text becomes `Highlight`, and the focus ring becomes `Highlight`, so selection survives without color.

## Keyboard

Because every tab is a real link, keyboard use is native: `Tab` moves between tabs, `Enter` activates the focused tab, and browser conventions (focus outline, current-page announcement) apply unchanged.

## Progressive enhancement

There is no enhancement layer: the server-rendered link flow is the complete flow. HTMX panel swapping could be added later, but it is not required — the fallback already navigates normally with plain links.
