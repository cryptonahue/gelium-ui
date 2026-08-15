# Button

Button is an open-code component built from native HTML. Use a `<button>` for actions and an `<a>` only for navigation.

## Alternative names

- Push button, action button, submit button, call to action (CTA)

## Agent prompt

Use Button when the user performs an action in the page; use an `<a>` only for navigation. The component is a native `<button>` with primary, secondary, outline and text variants; loading renders a native disabled button with `aria-busy="true"` and the dynamic accessible name `Loading {Label}`, and disabled/loading states never rely on color alone. Never replace it with a styled `<div>`: the semantic root must stay a real control with its own focus, and decorative icons must stay `aria-hidden`.

## Variants and states

The live preview below uses the same `button` Go template as this documentation layout. Loading controls are native disabled buttons with `aria-busy="true"`; link-shaped controls in disabled or loading state omit `href`, leave the tab order, and expose their inactive state, so neither can activate while work is in progress. In both forms, loading exposes the dynamic accessible name `Loading {Label}` while keeping the visual label hidden from assistive technology.

`IconSVG` is a per-instance inline SVG slot typed as `template.HTML`. Pass only trusted, internal markup—never user input. Icon SVGs are decorative by contract and must include `aria-hidden="true"` and `focusable="false"`; the button `Label` supplies the action text used by the accessible name.

Keyboard focus uses an exterior 3 px `:focus-visible` outline with a 2 px offset. Disabled and loading states do not rely on color alone.
