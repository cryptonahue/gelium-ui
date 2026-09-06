# Gelium UI Neubrutalism theme

Visual direction: loud editorial neubrutalism — warm paper canvas, hard ink borders, zero-radius geometry, offset shadows with no blur, heavy type, and saturated accent surfaces.

References: <https://github.com/neubrutalism/neubrutalism.com> and <https://github.com/ekmas/neobrutalism-components> (MIT). This theme is a token-only translation of their public design-token patterns; it does not copy markup, JavaScript, layout, assets, or page structure.

## Mapping

| Reference token / rule | Gelium token family |
|---|---|
| `--ink` | `--ui-color-border`, `--ui-color-border-strong`, hard shadow color |
| `--bg` | `--ui-color-canvas` |
| `--surface` | `--ui-color-surface`, card/dialog/toast surfaces |
| `--surface-2`, `--bg-warm` | `--ui-color-surface-container` |
| `--yellow` | light-mode secondary surface and warning accent; dark-mode primary |
| `--blue` / `--main` | light-mode AA-safe primary and info accent |
| `--green` | success accent and selected control tone |
| `--pink` / `--red` | badge and danger tone |
| `border: 3px solid var(--ink)` | `--ui-border-width-1`, component border tokens |
| `box-shadow: 5px 5px 0 var(--ink)` | `--ui-shadow-*`, button/card/dialog shadows |
| `border-radius: 0` | `--ui-radius-*`, component radius tokens |

## Light and dark

Light mode follows the warm paper and black-ink model. Its primary is a darker blue because Gelium uses `--ui-color-primary` both as interactive text on a light surface and as a filled-control background; the original bright yellow reaches only 1.44:1 on white. Yellow remains a saturated secondary/warning surface with black foreground. Dark mode follows the “cyber-brutalism” direction: dark surfaces, light structural lines, and yellow primary with dark foreground.

The dark route is class-driven only:

```html
<html class="theme-neubrutalism theme-dark" data-theme="dark">
```

No `prefers-color-scheme` media route is defined by the theme.

## Divergences

- Gelium keeps its HTML, templates, component classes, server contracts, HTMX enhancement hooks, and no-JS behavior unchanged.
- The reference's layout, manifesto page structure, badges, and decorative section composition are not part of this theme. Those would be screen design, not token theming.
- Reference fonts are mapped to system-resolvable stacks. The theme does not add remote Google Fonts or new self-hosted font files in this work unit.
- Fields keep Gelium's floating-label anatomy. Neubrutalist treatment is expressed through border, shadow, radius, color, and type tokens.
- Unlike references that separate `main` surface color from general `foreground`, Gelium's current primary token also colors links and selected text. The light primary is therefore darkened to meet WCAG AA instead of rendering yellow text on white.

## Verification checklist

- `theme-neubrutalism.css` defines light and dark token routes.
- Component token families are present for the existing Gelium matrix.
- The site imports the theme in the reference layer.
- The server catalog exposes the `neubrutalism` slug.
- Build and Go tests must pass after integration.
