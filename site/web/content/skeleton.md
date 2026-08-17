# Skeleton

Skeleton is a server-rendered loading placeholder for a region the server will fill with a later response — an HTMX swap or the next full page GET. Use a skeleton while data loads, so the layout stays stable and the user sees progress instead of a silent freeze. The root is a live region (`role="status"`) whose only accessible content is an `sr-only` label; the visual blocks are `aria-hidden` decoration. No component JavaScript is involved.

## Examples

<div class="component-preview">
  <div class="ui-skeleton" role="status">
    <span class="sr-only">Loading conversations</span>
    <div class="ui-skeleton-blocks" aria-hidden="true">
      <span class="ui-skeleton-block ui-skeleton-block--title"></span>
      <span class="ui-skeleton-block ui-skeleton-block--line"></span>
      <span class="ui-skeleton-block ui-skeleton-block--line"></span>
      <span class="ui-skeleton-block ui-skeleton-block--line ui-skeleton-block--short"></span>
    </div>
  </div>
</div>

<div class="component-preview">
  <div class="ui-skeleton ui-skeleton--avatar" role="status">
    <span class="sr-only">Loading profile</span>
    <div class="ui-skeleton-blocks" aria-hidden="true">
      <span class="ui-skeleton-block ui-skeleton-block--circle"></span>
      <span class="ui-skeleton-block ui-skeleton-block--title"></span>
      <span class="ui-skeleton-block ui-skeleton-block--line"></span>
      <span class="ui-skeleton-block ui-skeleton-block--line ui-skeleton-block--short"></span>
    </div>
  </div>
</div>

## Guidance

### When to use

Use a skeleton for a region that is about to be filled: a list that will arrive via swap, a table awaiting its rows, a profile pane on first paint. Scope it to the loading region — one skeleton for the list card, not nested full-page blockers (see [Feedback](/docs/feedback), FEED-LOAD-LIST and FEED-LOAD-PAGE). A page-level skeleton is rendered once, then replaced by content.

### When not to use

Do not use a skeleton for an operation that fails: a fetch error is a real status and belongs to the [error state](/components/error-state) with a retry. Do not show a skeleton forever — never an infinite placeholder. And do not use a skeleton to load a button: a button carries its own spinner and `aria-busy` state. If the whole page has nothing else to show, prefer progressive rendering: the shell first, then one skeleton for the main region.

### Usability

- A plain `ui-skeleton` stacks title, line and short blocks; `ui-skeleton--avatar` adds a circle in a grid with the same blocks.
- The blocks are bare decorations collapsed to nothing — only the `sr-only` label is announced, so screen readers hear "Loading…" instead of empty boxes.
- The fill uses the Material 3 pulse (opacity 0.5 ↔ 1) disabled under `prefers-reduced-motion`.

### Accessibility

- The root is a live region: `role="status"` announces the `sr-only` label when the skeleton appears.
- The blocks are `aria-hidden` decoration that must carry no announced content.
- Under reduced motion the pulse is disabled; under forced colors the blocks paint `CanvasText` so the placeholder stays visible without color.

## Anatomy

- **Skeleton** — `ui-skeleton`, the root live region with the scoped `--ui-skeleton-*` tokens (`padding`, `gap`, `block-height`, `block-color`, `pulse-duration`).
- **Label** — an `sr-only` span inside the root that carries the accessible "Loading…" text.
- **Blocks** — `ui-skeleton-blocks`, an `aria-hidden` container for the decorative shapes.
- **Shapes** — `ui-skeleton-block` with the modifiers `--title`, `--line`, `--short` (60% width) and, on the avatar variant, `--circle`.

## Variants

- **Plain** — `ui-skeleton`, a stacked set of title and line blocks for lists, tables and panes.
- **Avatar** — `ui-skeleton ui-skeleton--avatar`, a grid that pairs a circle with title, line and short blocks for profiles and feed items.

## Checklist

- [ ] Skeleton scoped to the loading region, not the whole app chrome.
- [ ] One `sr-only` label with a real message; blocks are `aria-hidden`.
- [ ] Replaced by content when the swap lands — never left forever.
- [ ] A failure path exists: the region ends in content or an error state, not an endless skeleton.
- [ ] Pulse disabled under reduced motion; blocks visible under forced colors.

## Accessibility

The skeleton announces a loading state once (live region) and its decorative blocks stay hidden from assistive technology. Because the placeholder is a deception of final content, the label must be honest ("Loading conversations…") and the region must resolve to real content or an error state. The pulse adheres to reduced-motion preferences.

## See also

- [Feedback](/docs/feedback) — loading rules (FEED-LOAD-LIST, FEED-LOAD-PAGE, FEED-LOAD-FAIL).
- [Error state](/components/error-state) — the failure path when the loaded region cannot be delivered.
- [Progress](/components/progress) — the determinate alternative for operations with a known length.
- [Choose the right control](/docs/choose-the-right-control) — the cross-component decision.