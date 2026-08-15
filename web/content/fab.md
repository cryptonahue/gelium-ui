# Floating action button (FAB)

FAB (Floating Action Button) represents the most important action on a screen. Use a FAB when exactly one action dominates the view — compose, add, start — and it must stay reachable in the corner at every screen size. It is icon-anchored and uses a native `<button type="button">` for actions and an `<a href>` for navigation, exactly like Button. Extended FABs are wider to accommodate a visible text label.

## Guidance

### When to use

Use a FAB when exactly one action dominates the view — compose, add, start — and it must stay reachable in the corner at every screen size.

### When not to use

Do not use a FAB when several actions share importance — a [Button](/components/button) or [Icon button](/components/icon-button) in the surface reads better. If the action is not the screen's primary task, it does not deserve a FAB.

### Usability

- Pick the emphasis that fits: `primary`, `surface`, or `secondary`; sizes are `small`, `medium`, and `large`.
- Use an extended FAB when the action needs a visible text label.
- Disabled FABs remove their activation path: native `disabled` on buttons, omitted `href` plus `aria-disabled="true"` on links.

### Accessibility

- An icon-only FAB has no visible text, so its name comes from `aria-label` — required and non-empty.
- The extended FAB supplies a visible label; give it an `aria-label` too when the context needs a richer name.
- Keyboard focus uses an exterior `:focus-visible` outline that does not change geometry.
- Focus, hover, and pressed states are never communicated by color alone.

See [Choose the right control](/docs/choose-the-right-control) for the cross-component decision.
## Variants and states

The live preview below uses the same `fab` Go template as this documentation layout. Color variants are `primary`, `surface`, and `secondary`; sizes are `small`, `medium`, and `large`. The `lowered` variant frees the FAB to a lower elevation. Disabled controls remove their activation path: buttons are natively `disabled` and link-shaped FABs omit `href`, leave the tab order, and expose `aria-disabled="true"`.

`IconSVG` is a per-instance inline SVG slot typed as `template.HTML`. Pass only trusted, internal markup — never user input. FAB icons are decorative by contract and must include `aria-hidden="true"` and `focusable="false"`.

## Accessibility

An icon-only FAB has no visible text, so its name comes from `aria-label`, which is **required and non-empty**. The extended FAB supplies a visible label; give it an `aria-label` too when the context needs a richer name. Keyboard focus uses an exterior `:focus-visible` outline that does not change geometry.

Focus, hover, and pressed states are never communicated by color alone; disabled and inactive states do not rely on color only.
