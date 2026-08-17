# Density and shell

**Density** is how much UI fits per viewport. **Shell** is the persistent chrome (nav, top bar, main column). Getting both wrong makes a correct component set feel “not product-grade.”

## Sources

| Topic | Source |
|---|---|
| Touch / hit targets | Material 3 touch guidance; Gelium `--ui-touch-target` ([Responsive](/docs/responsive)) |
| Reading measure | Common readability practice; Gelium docs prose **65ch** (≤75ch) |
| App destinations | [M3: navigation bar](https://m3.material.io/components/navigation-bar/guidelines) |
| Section nav | [USWDS: side navigation](https://designsystem.digital.gov/components/side-navigation/) |
| Container | Gelium `--ui-container-max`, `.ui-container` |

## Density modes

| Mode | When | Traits | Don’t |
|---|---|---|---|
| **Comfortable** | Consumer, reading, marketing-adjacent | Larger gaps, fewer columns, card/list bias | Cram admin tables into marketing spacing |
| **Cozy (default app)** | General tools | Token spacing as shipped (`--ui-space-*`) | Mix random px gaps beside tokens |
| **Compact** | Power admin, high data | Tighter row padding; still **touch-target floor** on controls | Shrink hit areas below `--ui-touch-target` |

Pick **one mode per product surface** (admin vs public), not per component instance.

## Shell layout

| Region | Role | Rules |
|---|---|---|
| **Top bar** | Brand, global actions, user menu | Keep height stable; don’t hide primary nav only in unlabeled icon clusters |
| **Side nav** | Section hierarchy 1–3 levels | USWDS: skip heavy sidenav on tiny sites |
| **Main** | Page job | One H1; primary action in top band ([Screens](/docs/screens)) |
| **Aside** | Secondary only | Don’t put the only CTA exclusively in aside on narrow |
| **Max width** | Readable app column | `--ui-container-max` / `.ui-container` for app shells; prose 65ch for long text |

### from-desktop

Default to **stacked** shell regions; enhance to row/multi-column from the desktop step (`.ui-row-from-desktop`, media `min-width` ~48rem). See [Responsive](/docs/responsive).

## Anti-patterns

- Compact density + 44px targets violated.
- Main content full-bleed paragraphs on ultrawide.
- Side nav + bottom nav + tabs all competing for the same destinations.
- Different density on every page of the same admin.

## Checklist

1. Name density mode for the surface.
2. Draw shell regions and what is persistent.
3. Main column uses container/measure tokens.
4. Narrow: stack; no `overflow-x: hidden` on body.

## See also

- [Screens](/docs/screens) · [Responsive](/docs/responsive) · [Tokens](/docs/tokens) · [`/llms-ux.txt`](/llms-ux.txt)
