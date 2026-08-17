# Spacing

Gelium spacing is a small, named rhythm shared by components and compositions. The live bars below use `--ui-space-1` through `--ui-space-8`; the labels are generated in HTML so the page remains useful with JavaScript disabled.

## Scale

<div class="docs-foundation-showcase docs-space-scale" aria-label="Gelium spacing scale">
  <div class="docs-space-row"><code>space-1</code><span class="docs-space-bar" style="--space-demo:var(--ui-space-1)"></span><code>.25rem</code></div>
  <div class="docs-space-row"><code>space-2</code><span class="docs-space-bar" style="--space-demo:var(--ui-space-2)"></span><code>.5rem</code></div>
  <div class="docs-space-row"><code>space-3</code><span class="docs-space-bar" style="--space-demo:var(--ui-space-3)"></span><code>.75rem</code></div>
  <div class="docs-space-row"><code>space-4</code><span class="docs-space-bar" style="--space-demo:var(--ui-space-4)"></span><code>1rem</code></div>
  <div class="docs-space-row"><code>space-6</code><span class="docs-space-bar" style="--space-demo:var(--ui-space-6)"></span><code>1.5rem</code></div>
  <div class="docs-space-row"><code>space-8</code><span class="docs-space-bar" style="--space-demo:var(--ui-space-8)"></span><code>2rem</code></div>
</div>

## Composition

This small server-rendered composition shows the distinction between a layout gap (`space-4`), component padding (`space-3`), and text rhythm (`space-2`). It is deliberately plain: spacing should clarify grouping, not become decoration.

<div class="docs-space-composition" aria-label="Spacing composition example">
  <div><strong>Queue</strong><span>12 open items</span><a href="/docs/spacing">Review queue</a></div>
  <div><strong>Review</strong><span>3 need attention</span><a href="/docs/spacing">Open review</a></div>
  <div><strong>Done</strong><span>48 completed</span><a href="/docs/spacing">See history</a></div>
</div>

## Rules

- `SPACE-SCALE`: use `--ui-space-1`, `2`, `3`, `4`, `6`, or `8`; do not create a second numeric scale in feature CSS.
- `SPACE-GROUP`: use smaller gaps inside a component and larger gaps between groups; whitespace communicates hierarchy.
- `SPACE-COLLAPSE`: at narrow widths, stack compositions instead of forcing a minimum width or hiding overflow.
- `SPACE-TARGET`: spacing never shrinks the interactive hit area below `--ui-touch-target`.

## Sources and Gelium adaptation

- [GOV.UK spacing](https://design-system.service.gov.uk/styles/spacing/): Gelium adapts a constrained rhythm and responsive composition to the `--ui-space-*` family.[1]
- [USWDS spacing tokens](https://designsystem.digital.gov/design-tokens/spacing/): Gelium keeps spacing semantic and reusable across components.[2]
- [W3C WCAG 2.2, 1.4.12 Text spacing](https://www.w3.org/TR/WCAG22/#text-spacing): compositions must survive increased user spacing.[3]

See also: [Responsive](/docs/responsive), [Density](/docs/density), [Tokens](/docs/tokens).
