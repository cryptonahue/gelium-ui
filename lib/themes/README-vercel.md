# Gelium UI — Vercel theme

Vercel is a light-first, typeset-terminal-on-white-paper translation of the
[Refero Vercel style](https://styles.refero.design/style/f24daf3a-d43f-4dec-85a9-8ac1d5148a03)
and its companion [DESIGN.md](https://github.com/educlopez/design-bites/blob/main/design-mds/vercel.com/DESIGN.md).
It is token-only: Gelium's semantic HTML, server-rendered/0-JS behavior, HTMX
contracts, focus-visible behavior, forced-colors behavior, and component
anatomy remain unchanged.

## Mapping

| Reference direction | Gelium implementation |
|---|---|
| `#fafafa` paper canvas | `--ui-color-canvas` |
| White elevated surfaces | `--ui-color-surface`, card/dialog/toast containers |
| `#f2f2f2` recessed surfaces | `--ui-color-surface-container`, filled controls/cards |
| `#171717` ink | `--ui-color-fg`, black primary button/FAB surfaces |
| `#4d4d4d` secondary ink | `--ui-color-fg-muted` and field labels |
| Vercel blue `#0072f5` | `--ui-color-primary` and interactive controls |
| Focus blue `#005fcc` | `--ui-color-focus-ring` / field focus token |
| `#ebebeb` hairline | `--ui-color-border`, card outline, shadow-as-border |
| Category palette | status tokens, dots, badges and toast icons only; never large status surfaces |
| 4px rhythm, 1200–1280px layout direction | `--ui-space-*` and `--ui-container-max: 80rem` |
| 6px functional corners / 12px panels | anatomy radius tokens; 9999px is reserved for compact pills |
| Geist Sans/Mono direction | Inter/system sans and Geist Mono/system mono fallback stacks |

Display uses the existing Gelium steps with a 64px/1 line and approximately
`-0.06em` tracking. Body remains 16px/1.5 with non-negative tracking. Controls
use a 44px touch-safe geometry even where the reference is more compact.
Primary buttons are black filled; secondary/outline buttons remain the existing
Gelium variants and receive Vercel's hairline/shadow-as-border treatment.

## What Vercel cannot be applied directly

- **Verified Geist font pack unavailable.** The repository ships no verified
  Vercel/Geist font pack for this theme. The CSS names Geist where useful but
  resolves to Inter/system fallbacks; adding webfont loading is separate work.
- **32/40px form heights conflict with the 44px touch floor.** Reference-sized
  compact fields cannot be copied directly without violating Gelium's
  accessibility contract. Fields stay at 3rem and controls at a touch-safe
  2.75rem minimum.
- **The 45px breakpoint list is intentionally not copied.** It overfits the
  reference marketing layout and is not a Gelium token or responsive contract.
  The theme keeps the core responsive behavior and an approximate 80rem maximum.
- **Vercel's brand triangle/wordmark is not a Gelium primitive.** No logo or
  brand-mark markup was introduced; consumers must supply their own approved
  asset outside the component theme contract.
- **Marketing gradients are kept out of core UI.** Decorative gradient heroes,
  glows, and promotional backgrounds do not map to semantic component tokens and
  would reduce contrast and theme portability.
- **The reference's exact bespoke display/mono font metrics cannot be reproduced
  without the missing font pack.** Gelium uses its existing decomposed type
  steps and never redeclares shorthand aliases.
- **A Vercel-specific double focus ring is not applied.** The current core owns
  safe outline and forced-colors behavior; there is no safe themeable focus-shadow
  contract, so this remains a documented phase-2/core capability gap.
- **Reference-specific navigation, avatar stacks, product stages, and status
  layouts are not new primitives.** The gallery uses existing native navigation,
  Badge/status anatomy, live components, and the existing semantic blockquote.

## Accessibility and integration

Status/category hues are dots or compact status anatomy rather than large fills.
Native controls, no-JS navigation, focus-visible indicators, forced-colors
fallbacks, readable body text, and `--ui-touch-target: 44px` remain active.
The theme is imported after the core manifest and allowlisted as `vercel` in the
server theme catalog. Explicit dark mode is available via
`.theme-vercel.theme-dark`, `.theme-vercel.dark`, or `[data-theme="dark"]`; the
light route remains the product default.
