# Empty state

Empty state is a server-rendered status message for an empty result — a list with no items, a search with no matches, a feed with nothing to show. Use an empty state when a surface is legitimately empty, so the user gets a reason and a next step instead of a silent blank area. It renders as a real `div` styled with the `ui-empty-state` class and needs no component JavaScript.

## Examples

<div class="component-preview">
  <div class="ui-empty-state" role="status">
    <span class="ui-empty-state-icon" aria-hidden="true">▧</span>
    <p class="ui-empty-state-title">No projects yet</p>
    <p class="ui-empty-state-body">Create your first project to get started.</p>
    <a class="ui-button" href="/projects/new">Create project</a>
  </div>
</div>

<div class="component-preview">
  <div class="ui-empty-state ui-empty-state--compact" role="status">
    <span class="ui-empty-state-icon" aria-hidden="true">▧</span>
    <p class="ui-empty-state-title">No matching rows</p>
    <p class="ui-empty-state-body">No records match <code>status: draft</code>. Clear filters to see all records.</p>
    <a class="ui-button" href="/table?clear=1">Clear filters</a>
  </div>
</div>

## Guidance

### When to use

Use an empty state when a collection, search result or feed is empty: it explains why the surface is empty and offers a real, actionable step. The centered default suits page-level empty feeds and searches; the `ui-empty-state--compact` modifier starts at the inline edge for inline contexts like a table row region.

### When not to use

Do not leave a bare "0 rows" or "no data" with no explanation and no next step (see [Feedback](/docs/feedback), FEED-EMPTY). Do not use an empty state for a loading state — that is a [skeleton](/components/skeleton). And do not use an empty state for a failure: a fetch that came back empty because the request failed is an [error state](/components/error-state), not an empty surface.

### Usability

- The title is short and specific; the body names the reason and the action.
- After a filtered or searched empty, name the active filter and offer a real "Clear filters" control.
- The CTA is a real link (or action) — never an instruction in prose with no control.

### Accessibility

- The root is a live region: `role="status"` announces the message when the empty state appears.
- The icon is decorative: keep it `aria-hidden` and put the meaning in the title and body.
- Under forced colors the icon, title and body repaint `CanvasText` and the CTA stays a link (`LinkText`), so the empty state never depends on color.

## Anatomy

- **Empty state** — `ui-empty-state`, the centered flex column with the scoped `--ui-empty-state-*` tokens (`padding`, `gap`, `icon-size`, `title-color`, `body-color`).
- **Icon** — `ui-empty-state-icon`, a decorative `aria-hidden` glyph.
- **Title** — `ui-empty-state-title`, the short headline (`--ui-type-title-md`).
- **Body** — `ui-empty-state-body`, the explanation and next step (`--ui-type-body-sm`).
- **CTA** — an optional `ui-button` action that is a real control.

## Variants

- **Default** — `ui-empty-state`, centered for page-level feeds, searches and collections.
- **Compact** — `ui-empty-state ui-empty-state--compact`, start-aligned for inline contexts like a table row or a small panel.

## Checklist

- [ ] The surface is truly empty — not loading, not failed.
- [ ] Title + body explain why and what to do next — no bare "0 rows".
- [ ] Filtered empties name the filter and offer a real "Clear filters" control.
- [ ] CTA, when present, is a real control.
- [ ] Icon `aria-hidden`; `role="status"` present; forced-colors output stays readable.

## Accessibility

The empty state announces its message through `role="status"` and keeps all meaning in text: title, body and a real control, never color or icon alone. The compact variant keeps the same announcement while changing only the visual alignment.

## See also

- [Feedback](/docs/feedback) — the decision matrix entry FEED-EMPTY.
- [Error state](/components/error-state) — the failure surface when a request failed instead of returning empty.
- [Skeleton](/components/skeleton) — the loading surface shown before the empty (or full) result arrives.
- [Choose the right control](/docs/choose-the-right-control) — the cross-component decision.