# Button

Button is an open-code component built from native HTML. Use a `<button>` for actions and an `<a>` only for navigation.

## Guidance

### When to use

Use a button for any action the user performs in the page — submit, save, open, confirm. Use an `<a>` only for navigation to another page or section.

### When not to use

Do not use a button for navigation: a link is the navigational control. Never style a `<div>` as a button — the semantic root must stay a real control. For a compact glyph-only action in a toolbar, prefer an [Icon button](/components/icon-button); when exactly one action dominates the screen, a [FAB](/components/fab) keeps it reachable.

### Usability

- Pick the emphasis that fits: primary for the main action, secondary/outlined for alternatives, text for the quietest level.
- Loading renders a native disabled button with `aria-busy="true"` and the dynamic accessible name `Loading {Label}`.
- Decorative icons are `aria-hidden` and `focusable="false"`; the label supplies the action text.

### Accessibility

- The native `<button>` keeps its focus, activation, and form behavior at no cost.
- Keyboard focus uses an exterior 3 px `:focus-visible` outline with a 2 px offset that never changes geometry.
- Disabled and loading states never rely on color alone. Link-shaped controls in disabled or loading state omit `href`, leave the tab order, and expose their inactive state.

See [Choose the right control](/docs/choose-the-right-control) for the cross-component decision.
## Alternative names

- Push button, action button, submit button, call to action (CTA).

## Agent prompt

Use Button when the user performs an action in the page. Use an `<a>` only for navigation. The component is a native `<button>` with primary, secondary, outline and text variants. Loading renders a native disabled button with `aria-busy="true"` and the dynamic accessible name `Loading {Label}`. Disabled and loading states never rely on color alone. Never replace it with a styled `<div>`: the semantic root must stay a real control with its own focus, and decorative icons must stay `aria-hidden`.

## Variants and states

The live preview below uses the same `button` Go template as this documentation layout. Loading controls are native disabled buttons with `aria-busy="true"`. Link-shaped controls in disabled or loading state omit `href`, leave the tab order, and expose their inactive state. Neither can activate while work is in progress. In both forms, loading exposes the dynamic accessible name `Loading {Label}` while keeping the visual label hidden from assistive technology.

`IconSVG` is a per-instance inline SVG slot typed as `template.HTML`. Pass only trusted, internal markup—never user input. Icon SVGs are decorative by contract and must include `aria-hidden="true"` and `focusable="false"`; the button `Label` supplies the action text used by the accessible name.

Keyboard focus uses an exterior 3 px `:focus-visible` outline with a 2 px offset. Disabled and loading states do not rely on color alone.
