# Menu

Menus show a list of choices on a temporary surface anchored to a trigger. Use a menu when an action needs a compact set of options that should not clutter the page. The open/close flow must work without JavaScript. Gelium reimplements the Material 3 menu over **server-rendered HTML with the native Popover API** (the browser's declarative popup mechanism). The menu surface is a `[popover]` element opened and closed declaratively by a real trigger button through `popovertarget` / `popovertargetaction`. CSS anchor positioning positions it when the browser supports it. There is **no component JavaScript** — the trigger, the top-layer open/close, light dismiss and Escape are all platform-native.

## Guidance

### When to use

Use a menu when an action needs a compact set of options that should not clutter the page. The open/close flow must work without JavaScript.

### When not to use

Do not use a menu when the choice is data entry. A form value belongs in a [Radio](/components/radio) group or a [Select](/components/select), not a command menu. For a permanent, always-visible set of choices, prefer [Tabs](/components/tabs) or [Segmented buttons](/components/segmented-button).

### Usability

- Action items are real `<button type="button">` elements; navigation items are real `<a href>` links.
- Selection rows wrap a native checkbox (multi-select) or radio (single-select) in a `<label>` inside a real form.
- The surface opens declaratively via `popovertarget`; light dismiss and Escape are platform-native.

### Accessibility

- The trigger is a real `<button>` with `aria-haspopup="menu"`; `popovertarget` gives the popover its implicit `aria-expanded`/`aria-details` relationship.
- Items are native links/buttons/checkboxes/radios — role, name, keyboard and disabled behavior come from the platform.
- The surface is a `<ul>` so item rows keep list semantics; the divider is `role="separator"`.
- Focus ring uses `:focus-visible`; decorative icons are `aria-hidden`; disabled never relies on color alone.

See [Choose the right control](/docs/choose-the-right-control) for the cross-component decision.
## Anatomy

- **Surface** — `ui-menu`, the elevated container (container color, `--ui-radius-xs` shape, `--ui-shadow-2` elevation, 8px block padding). On the open anatomy demo the surface is a `<ul>` so the items keep native list semantics.
- **Item** — `ui-menu-item`, the 48px-high flex row carrying a leading icon slot, the label and an optional selection control. Items are real interactive elements: a `<button type="button">` for actions or an `<a href>` for navigation.
- **Label** — `ui-menu-item-label`, the body-large text that supplies the accessible name.
- **Leading icon** — `ui-menu-item-icon`, the optional 24px decorative glyph (reuses the Icon contract: trusted inline SVG, `aria-hidden`, `focusable="false"`).
- **Divider** — `ui-menu-divider`, a 1px `role="separator"` break between groups.
- **Selection** — `ui-menu-item--select` rows wrap a native `<input type="checkbox">` or `<input type="radio">` in a `<label>`, exactly like List and Segmented button. The selected state derives from `:checked` (the CSS pseudo-class matching a checked native input), with no JavaScript.

## Variants and states

- **Actions** — `<button type="button">` items.
- **Navigation** — `<a href>` items that navigate to real routes.
- **Selection** — native checkbox (multi-select) and radio (single-select) items inside a real `<form>`.
- Each item can carry an optional leading icon; groups are separated by a divider; any item can be **disabled** (`disabled` on buttons, `:has(input:disabled)` for selection).
- States: rest, hover (state layer), focus-visible (3px focus ring, no geometry shift), active/pressed (state layer), disabled (dimmed), selected (`:checked` for checkbox/radio rows).

## Accessibility

- The trigger is a real `<button>` with `aria-haspopup="menu"`; `popovertarget` gives the popover its implicit `aria-expanded`/`aria-details` relationship for free.
- Items are native links/buttons/checkboxes/radios — role, name, keyboard and disabled behavior come from the platform. The menu intentionally does not force `role="menu"`/`role="menuitem"` over native link/button semantics (the roadmap's "ARIA (role=menu/menuitem vs native links)"). It matches how Tabs uses real links without a fake tablist.
- The surface is a `<ul>` so item rows keep list semantics; the divider is `role="separator"`.
- Focus ring uses `:focus-visible`; decorative icons are `aria-hidden`; disabled never relies on color alone.

## How does the menu open and close without JavaScript?

The whole flow works with JavaScript disabled: the trigger button opens the popover declaratively via `popovertarget`. Navigation items are real links, and action items are real buttons. The selection menu is a real `<form method="get" action="/components/menu">` whose checked values submit through a normal round-trip. Light dismiss (outside click) and Escape close the surface natively. There is no CSS-only imitation: the same markup is a real navigation or form.

## Compatibility

Popover API (`popover`, `popovertarget`, `popovertargetaction`) is Baseline 2025 (newly available). CSS anchor positioning (`anchor-name`/`position-anchor`/`anchor()`) is Baseline 2026 (newly available), applied as a progressive enhancement inside `@supports`. Without it the menu still opens in the top layer with a fallback position. `:has()` and `:focus-visible` are Baseline. No browser is blocked from the no-JS navigation/form flow. A browser without the Popover API simply renders the menu surface statically with its real links and buttons.

## Trust boundary

The demo markup is static and trusted: `html/template` renders it verbatim, and the inline SVGs are trusted internal glyphs (`aria-hidden`, `focusable="false"`). There is no user-supplied content.

## Divergences from Material Web

- **Surface opening** — upstream `md-menu` uses component JavaScript (`show()`/`close()`) to manage a popover-positioned surface. Gelium uses the native `popover` attribute + `popovertarget`/`popovertargetaction` and CSS anchor positioning, with no component JS.
- **Roles** — upstream defaults to `role="menu"` with `role="menuitem"` items and full keyboard/typeahead JS. Gelium keeps native link/button/checkbox/radio semantics (real links navigate, real buttons/forms submit). No roving focus or typeahead JavaScript is needed; Tab/Enter/Space and native radio arrows cover the keyboard.
- **Geometry** — upstream items inherit list-item tokens (56px); Gelium pins the Material menu item height to 48px with a 12px row padding. This follows the menu item token contract.
- **Submenus** — deferred: upstream `md-sub-menu` needs hover/timing JS; Gelium documents nested menus as a follow-up rather than porting the JS controller.

## Visual checklist

- [ ] Elevated surface with container radius and 8px block padding.
- [ ] 48px items with leading icon, label, optional divider.
- [ ] Action, navigation and checkbox/radio selection variants.
- [ ] Hover/focus/pressed state layers; focus ring without layout shift.
- [ ] Disabled items dim without color-only feedback.
- [ ] Selected checkbox/radio rows keyed off `:checked`.
- [ ] Trigger opens the surface declaratively; light dismiss + Escape.
- [ ] Light/dark, narrow/wide, RTL, reduced motion, forced colors.
