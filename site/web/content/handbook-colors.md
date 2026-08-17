# Colors

Gelium color is semantic: components consume roles such as canvas, foreground, primary, border, focus, and status rather than hard-coded brand values. The swatches below use the live tokens, so changing `?theme=basecoat` or the server-rendered light/dark scheme updates the same markup.

## Roles and states

<div class="docs-foundation-showcase docs-color-grid" aria-label="Gelium semantic color palette">
  <div class="docs-color-swatch" style="--swatch-bg:var(--ui-color-canvas);--swatch-fg:var(--ui-color-fg)"><strong>Canvas / foreground</strong><code>--ui-color-canvas<br>--ui-color-fg</code><span>Reading surface</span></div>
  <div class="docs-color-swatch" style="--swatch-bg:var(--ui-color-surface);--swatch-fg:var(--ui-color-fg)"><strong>Surface / muted</strong><code>--ui-color-surface<br>--ui-color-fg-muted</code><span>Secondary grouping</span></div>
  <div class="docs-color-swatch" style="--swatch-bg:var(--ui-color-primary);--swatch-fg:var(--ui-color-primary-fg)"><strong>Primary</strong><code>--ui-color-primary<br>--ui-color-primary-fg</code><span>Primary action</span></div>
  <div class="docs-color-swatch" style="--swatch-bg:var(--ui-color-success);--swatch-fg:var(--ui-color-success-fg)"><strong>Success</strong><code>--ui-color-success<br>--ui-color-success-fg</code><span>Completed or safe</span></div>
  <div class="docs-color-swatch" style="--swatch-bg:var(--ui-color-warning-container);--swatch-fg:var(--ui-color-warning-fg)"><strong>Warning</strong><code>--ui-color-warning-container<br>--ui-color-warning-fg</code><span>Needs attention</span></div>
  <div class="docs-color-swatch" style="--swatch-bg:var(--ui-color-danger);--swatch-fg:var(--ui-color-danger-fg)"><strong>Danger</strong><code>--ui-color-danger<br>--ui-color-danger-fg</code><span>Destructive or error</span></div>
</div>

Keyboard focus is a role, not a color afterthought: the link below uses the primary pair and gets the theme's focus ring on `:focus-visible`.

<p><a class="docs-focus-example" href="/docs/colors">Focus this link with Tab</a></p>

## Rules

- `COLOR-ROLE`: use semantic foreground/background pairs; never infer text contrast from hue alone.
- `COLOR-CONTRAST`: verify text and non-text contrast for the active theme; status must not be communicated by color alone.
- `COLOR-FOCUS`: preserve `--ui-color-focus-ring`, `--ui-focus-thickness`, and `--ui-focus-offset` for keyboard visibility.
- `COLOR-SCHEME`: theme selection is server-rendered (`?theme=` and `?scheme=dark`), with no client-side palette swap required.

## Light and dark guidance

Material and Basecoat both provide `.theme-dark` variants. Keep markup and semantic roles unchanged; the theme changes token values and `color-scheme`. Test canvas, surface, primary, danger, status, border, and focus in both schemes. Do not assume a light token is safe on a dark surface, and do not use transparent state layers as the only focus indicator.

## Sources and Gelium adaptation

- [W3C WCAG 2.2, 1.4.3 Contrast (Minimum)](https://www.w3.org/TR/WCAG22/#contrast-minimum) and [1.4.11 Non-text Contrast](https://www.w3.org/TR/WCAG22/#non-text-contrast): Gelium keeps semantic pairs and an explicit focus token for both themes.[1][2]
- [USWDS color tokens](https://designsystem.digital.gov/design-tokens/color/): Gelium adapts semantic roles and theme-level overrides to `--ui-color-*`.[3]
- [Material 3 color roles](https://m3.material.io/styles/color/roles): Gelium uses the role-based idea, while retaining its own token names and server-rendered theme switch.[4]

See also: [Themes](/docs/themes), [Accessibility](/docs/accessibility), [Tokens](/docs/tokens).
