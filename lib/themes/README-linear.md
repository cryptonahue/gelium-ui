# Gelium UI — Linear theme

Linear is a dark-first, compact product UI translation of the Refero style
[90ce5883-bb24-4466-93f7-801cd617b0d1](https://styles.refero.design/style/90ce5883-bb24-4466-93f7-801cd617b0d1).
It is token-only: markup, component contracts, HTMX/server behavior, and native
accessibility semantics are unchanged.

## Reference mapping

| Refero role | Gelium token | Linear value (dark) |
|---|---|---|
| Void canvas | `--ui-color-canvas` | `#08090a` |
| Carbon cards/nav | `--ui-color-surface` | `#0f1011` |
| Obsidian elevated | `--ui-color-surface-container` | `#161718` |
| Graphite border | `--ui-color-border` / card outline | `#23252a` |
| Smoke strong border | `--ui-color-border-strong` | `#383b3f` |
| Ash muted text | `--ui-color-fg-muted` | `#62666d` |
| Fog tertiary / labels | field and dialog labels | `#8a8f98` |
| Mist secondary text | `--ui-color-secondary-fg` | `#d0d6e0` |
| Bone foreground | `--ui-color-fg` | `#e5e5e6` |
| Paper headings | component inheritance | `#ffffff` reference accent |
| Acid Lime action | `--ui-color-primary` | `#e4f222` |
| Pulse Green / Coral / Teal | success / danger / info | `#27a644` / `#eb5757` / `#02b8cc` |

Iris Violet (`#6366f1`) and Lavender (`#8b5cf6`) remain available to consuming
surfaces as tag/status accents rather than replacing semantic action colors.

## Type, density, shape

Linear uses system-resolved `Inter, ui-sans-serif, system-ui` and a
`JetBrains Mono`-first mono stack. No font files or preloads are added: the
existing Alden Inter pack is intentionally not presented as Linear-owned fonts.
Display/title tracking is tight, while body and body-sm keep non-negative
tracking and line-height `1.5`; body-lg/body-md remain `1rem` for readability.

The compact register uses a 4px rhythm, 2.25rem controls, 2.75rem fields, 6px
buttons/inputs, 12px cards, and 9999px pills. `--ui-touch-target` stays `44px`
and the small FAB remains `44px`; these floors take precedence over compactness.
Shadows are hairline/inset-first, with restrained outer shadows for overlays.

## Light/dark and accessibility divergences

`.theme-linear` is a usable light fallback for explicit light mode. The intended
route is `.theme-linear.theme-dark`, `.theme-linear.dark`, or
`.theme-linear[data-theme="dark"]`; there is no preference-media route. Strong
focus uses a 2px Acid Lime ring with 2px offset. Status foregrounds are tuned for
solid controls, while the bright accent palette is reserved for controls/tags,
not long body text. Native controls, forced-colors behavior, touch sizes, and
floating-label geometry remain core-owned.

## Coverage

The theme defines the complete semantic palette, state/focus/border/radius,
spacing/size/motion/type decomposition, and all TestThemeMatrix component
families: button, field, dialog, toast, card, badge, checkbox, radio, switch,
slider, progress, FAB, select, divider, plus their themeable anatomy tokens.
It is bundled after the core manifest and allowlisted as `theme-linear` in
`internal/app/server.go`. No component markup or server behavior is changed.
