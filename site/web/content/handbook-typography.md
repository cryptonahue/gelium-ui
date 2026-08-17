# Typography

Gelium typography is a token contract, not a page-specific scale. These server-rendered specimens read the active theme (`theme-material` or `theme-basecoat`) and the active light/dark class. Use the same role for the same job; do not tune a heading with a new font or arbitrary pixel value.

## Scale and measure

The decomposed tokens make size, weight, line-height, letter-spacing, and family inspectable. Body prose stays at a readable `65ch` measure; Gelium adapts the roughly 68ch USWDS and 30em GOV.UK guidance to one stable docs contract.[1][2]

<div class="docs-foundation-showcase" aria-label="Gelium typography specimens">
  <div class="docs-type-sample"><span class="docs-token-note">display-lg · --ui-type-display-lg-size · --ui-type-display-lg-weight · --ui-type-display-lg-line-height</span><p style="margin:0;font:var(--ui-type-display-lg);letter-spacing:var(--ui-type-display-lg-letter-spacing)">Design with intent</p></div>
  <div class="docs-type-sample"><span class="docs-token-note">headline-sm · --ui-type-headline-sm-size · --ui-type-headline-sm-weight · --ui-type-headline-sm-line-height</span><p style="margin:0;font:var(--ui-type-headline-sm);letter-spacing:var(--ui-type-headline-sm-letter-spacing)">A clear hierarchy helps people scan</p></div>
  <div class="docs-type-sample"><span class="docs-token-note">body-md · --ui-type-body-md-size · --ui-type-body-md-weight · --ui-type-body-md-line-height · measure 65ch</span><p style="max-width:65ch;margin:0;font:var(--ui-type-body-md);letter-spacing:var(--ui-type-body-md-letter-spacing)">Gelium keeps body copy calm and legible. Put one idea in each paragraph, keep the line length bounded, and let the theme provide the family.</p></div>
  <div class="docs-type-sample"><span class="docs-token-note">label-md · --ui-type-label-md-size · --ui-type-label-md-weight · --ui-type-label-md-line-height</span><p style="margin:0;font:var(--ui-type-label-md);letter-spacing:var(--ui-type-label-md-letter-spacing)">LABEL · STATUS · SUPPORTING METADATA</p></div>
</div>

## Rules

- `TYPE-ROLE`: choose display, headline, title, body, or label by content job, then use the matching `--ui-type-*` alias.
- `TYPE-MEASURE`: keep reading content near `65ch`; let controls and data tables use their own layout measure.
- `TYPE-RHYTHM`: preserve the token line-height; do not compress body text to fit a card.
- `TYPE-THEME`: theme overrides the decomposed values; consumers keep the alias so Material and Basecoat stay interchangeable.

## Sources and Gelium adaptation

- [W3C WCAG 2.2, 1.4.12 Text spacing](https://www.w3.org/TR/WCAG22/#text-spacing): Gelium leaves room for user text-spacing overrides instead of clipping prose.[3]
- [GOV.UK type scale](https://design-system.service.gov.uk/styles/type-scale/): Gelium adopts stepped, responsive-friendly roles rather than a fluid novelty scale.[1]
- [USWDS typography](https://designsystem.digital.gov/design-tokens/typesetting/): Gelium adapts readable measure and explicit type tokens to CSS custom properties.[2]

See also: [Tokens](/docs/tokens), [Themes](/docs/themes), [Accessibility](/docs/accessibility).
