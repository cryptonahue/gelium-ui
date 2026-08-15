# Elevation

Elevation is a Gelium-only primitive for visual depth. Use it when an element must sit above the surface it rests on — floating actions, menus, dialogs — without shipping images or per-instance shadows. Apply one of the six utility classes — `.ui-elevation-0` through `.ui-elevation-5` — to lift an element with the theme shadow tokens `--ui-shadow-0` through `--ui-shadow-5`.

## Guidance

### When to use

Use elevation when an element must sit above the surface it rests on — floating actions, menus, dialogs. It ships no images or per-instance shadows. Apply one of the six utility classes (`ui-elevation-0` through `ui-elevation-5`).

### When not to use

Do not elevate everything: depth should be earned by interactive or transient surfaces. Elevation must never be the sole carrier of meaning or state. Keep flat surfaces flat — over-elevation makes dialogs and menus compete instead of stand out.

### Usability

- Apply the utility classes directly on the element that should lift.
- Combine elevation with text, icons, or semantics when the elevated element communicates something important.
- The documentation layout dogfoods the same utility classes any feature would use.

### Accessibility

- Elevation conveys depth only — pair it with visible text or semantics for anything meaningful.
- `prefers-reduced-motion` disables the shadow transition.
- `forced-colors` removes decorative shadows entirely so elevation can never hide borders in high-contrast mode.

See [Choose the right control](/docs/choose-the-right-control) for the cross-component decision.
## Visual-only by contract

Elevation conveys depth only. It must never be the sole carrier of meaning or state: combine it with text, icons, or semantics when an elevated element communicates something important.

The preview below uses the utility classes directly, so the documentation layout dogfoods the same primitive any feature would use.

## Motion and system preferences

Elevation surfaces transition `box-shadow` on the theme motion token set. `prefers-reduced-motion` disables that transition, and `forced-colors` removes decorative shadows entirely so elevation can never hide borders in high-contrast mode.