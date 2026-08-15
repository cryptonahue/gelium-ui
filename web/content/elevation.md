# Elevation

Elevation is a Gelium-only primitive for visual depth. Use it when an element must sit above the surface it rests on — floating actions, menus, dialogs — without shipping images or per-instance shadows. Apply one of the six utility classes — `.ui-elevation-0` through `.ui-elevation-5` — to lift an element with the theme shadow tokens `--ui-shadow-0` through `--ui-shadow-5`.

## Visual-only by contract

Elevation conveys depth only. It must never be the sole carrier of meaning or state: combine it with text, icons, or semantics when an elevated element communicates something important.

The preview below uses the utility classes directly, so the documentation layout dogfoods the same primitive any feature would use.

## Motion and system preferences

Elevation surfaces transition `box-shadow` on the theme motion token set. `prefers-reduced-motion` disables that transition, and `forced-colors` removes decorative shadows entirely so elevation can never hide borders in high-contrast mode.