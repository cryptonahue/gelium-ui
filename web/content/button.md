# Button

Button is an open-code component built from native HTML. Use a `<button>` for actions and an `<a>` only for navigation.

## Variants and states

The live preview below uses the same `button` Go template as this documentation layout. Loading controls are native disabled buttons with `aria-busy="true"`; link-shaped controls in disabled or loading state omit `href`, leave the tab order, and expose their inactive state, so neither can activate while work is in progress. In both forms, loading exposes the dynamic accessible name `Loading {Label}` while keeping the visual label hidden from assistive technology.

`IconSVG` is a per-instance inline SVG slot typed as `template.HTML`. Pass only trusted, internal markup—never user input. Icon SVGs are decorative by contract and must include `aria-hidden="true"` and `focusable="false"`; the button `Label` supplies the action text used by the accessible name.

Keyboard focus uses an exterior 3 px `:focus-visible` outline with a 2 px offset. Disabled and loading states do not rely on color alone.
