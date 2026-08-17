# Footer

Footer is the site chrome below every page: brand, secondary navigation, and legal. Use it when a site needs a persistent lower band of navigation and legal context that complements, never replaces, the primary navigation. Sections collapse with native `<details>`/`<summary>` on narrow screens and open into a grid from the desktop breakpoint — zero JavaScript.

## Specimen

This page renders the real `footer` template markup live:

<footer class="ui-footer">
  <p class="ui-footer-brand">Gelium UI</p>
  <nav class="ui-footer-nav" aria-label="Footer">
    <section class="ui-footer-section">
      <details class="ui-footer-details">
        <summary class="ui-footer-heading">Docs</summary>
        <ul class="ui-footer-list">
          <li><a href="/docs">Documentation</a></li>
          <li><a href="/docs/screens">Screens</a></li>
          <li><a href="/blog">Blog</a></li>
        </ul>
      </details>
    </section>
    <section class="ui-footer-section">
      <details class="ui-footer-details">
        <summary class="ui-footer-heading">System</summary>
        <ul class="ui-footer-list">
          <li><a href="/docs/tokens">Tokens</a></li>
          <li><a href="/docs/themes">Themes</a></li>
          <li><a href="/docs/accessibility">Accessibility</a></li>
        </ul>
      </details>
    </section>
  </nav>
  <p class="ui-footer-legal">MIT License · Server-rendered docs with native HTML semantics.</p>
</footer>

## Guidance

### When to use

Use a footer on every site page as the persistent chrome: the brand, the secondary navigation (docs, legal, social), and the legal line. It is the expected place for content that matters globally but does not belong in the primary navigation.

### When not to use

Do not use a footer for primary navigation — the [Navigation bar](/components/navigation-bar) or [Navigation drawer](/components/navigation-drawer) owns the top-level destinations. Do not use it for in-page content footers; a [Divider](/components/divider) or a [Card](/components/card) communicates a content-adjacent boundary better. Do not duplicate the whole primary nav at full depth: the footer is secondary.

### Usability

- Sections stack as a vertical list of native disclosure blocks on narrow screens and expand into a grid (`auto-fit`, `minmax(12rem, 1fr)`) from `48rem`.
- The brand uses `--ui-footer-heading-type` (label-lg); links use `--ui-footer-type` (body-sm) with muted ink that darkens on hover.
- The legal line sits under a top border (`--ui-footer-border`) so it reads as a distinct band from the navigation.

### Accessibility

- The footer is a `<footer>` element containing a `<nav aria-label="Footer">`, so it is announced as a region with a navigation landmark.
- Section headings are native `<summary>` elements: keyboard users open and close them with Enter/Space on narrow screens, and on desktop the disclosure chrome is hidden by CSS with every section forced open.
- In forced-colors mode the footer paints a `CanvasText` top border and its text repaints as `CanvasText`/`LinkText`, so the chrome survives without color.

## Anatomy

- **`ui-footer`** — the chrome band: `--ui-footer-surface`, `--ui-footer-fg` muted ink, `--ui-footer-gap` (`--ui-space-6`) between brand, nav, and legal.
- **`ui-footer-brand`** — the brand line in `--ui-footer-heading-color`.
- **`ui-footer-nav`** — the secondary navigation grid; each `ui-footer-section` holds a `details.ui-footer-details` whose `summary.ui-footer-heading` names the group and whose `ul.ui-footer-list` carries the links.
- **`ui-footer-legal`** — the legal line under the top border.

All colors, type, and spacing come from the scoped `--ui-footer-*` tokens so a theme can retune the chrome without new theme work.

## Variants

The footer is fixed chrome: brand + secondary nav + legal. Its responsive behavior is the only dimension that varies — collapsed `<details>` sections on narrow screens, a forced-open grid from `48rem` — both handled by the stylesheet with zero JavaScript.

## Sources

- Registry: `docs/gelium-ui-component-registry.md` (Footer row, phase P) — `.ui-footer`, `--ui-footer-*`, `--ui-color-*`, `--ui-space-*`; no states; `<details>`/`<summary>` collapsible, grid → stack.
- Vocabulary: `docs/gelium-ui-vocabulary.md` §5 — Footer is the site chrome every Phase G recipe depends on.
- Implementation: `lib/templates/footer.html`, `lib/styles/footer.css`; view model `footerView` and `defaultFooter` in `internal/app/server.go`.

See also: [Screens](/docs/screens) for the chrome contract, [Information architecture](/docs/information-architecture) for nav grouping.