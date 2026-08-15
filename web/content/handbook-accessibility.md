# Accessibility

Accessibility is part of every component contract, not a retrofit: native semantics first, ARIA only where the platform has no equivalent, and every state tested in the build.

## Native semantics first

Components use the real platform element — `<button>`, `<dialog>`, `<input type="radio">`, `<table>`, `<nav>` — so behavior and assistive-technology support come from the browser. ARIA is added only where the platform has no equivalent. [Button](/components/button) is a real button or a real link; [Dialog](/components/dialog) is a native `<dialog>`.

## Accessible names

Every control has an accessible name. Icon-only controls ([Icon button](/components/icon-button), [FAB](/components/fab)) pair decorative inline SVG (`aria-hidden="true"`, `focusable="false"`) with a real label — visible or `sr-only`.

## States

Disabled and loading are exposed, never color-only. Disabled controls leave the tab order and expose their inactive state; loading controls add `aria-busy="true"` and a dynamic accessible name (`Loading {Label}`). Error fields carry `aria-invalid` and link their message via `aria-describedby`.

## Keyboard and focus

Keyboard focus uses an exterior 3 px `:focus-visible` outline with a 2 px offset ([Focus ring](/components/focus-ring)). Document order is the tab order, and there is no fake roving focus: [Tabs](/components/tabs) are real links with `aria-current`, and navigation components use server-derived active state.

## Live regions and feedback

Transient feedback renders inside the `#loom-toast-region` with `aria-live="polite"`; error toasts use `role="alert"`. Persistent validation errors render inline next to the field with `role="alert"` (error) or `role="status"` (other tones), and the 422 recovery keeps the user in context.

## Reduced motion

`prefers-reduced-motion` is centralized in the core: motion tokens (`--ui-motion-*`) collapse to instant or minimal transitions, with no per-component drift.

## Forced colors

`forced-colors` is centralized too. Components use tokens and `currentColor` for decorative fills, so the Windows high-contrast palette applies cleanly and focus rings stay visible.

## Contrast and target sizes

Text and controls meet WCAG AA contrast in light and dark. Touch targets respect the `--ui-size-*` minimums — [Switch](/components/switch), [Checkbox](/components/checkbox), and [Radio](/components/radio) keep their hit areas. RTL documents mirror automatically via logical properties and `:dir(rtl)` overrides where needed.

## How it is verified

The accessibility contract is enforced by tests: style contract tests assert token-driven states plus reduced-motion and forced-colors coverage; server tests assert `aria-*` attributes, live regions, and 422 recovery. See [Design principles](/docs/principles).
