# List

List is a continuous, vertical index of text and images. Use a list when people must scan, select, or navigate through a group of related rows — settings, actions, or collections. Gelium reimplements the Material 3 list over semantic HTML — the root is a real `<ul>`, `<ol>`, `<nav>`, or `<menu>`, and each item is a `<li>`. No component JavaScript: navigation items are real `<a href>` links, and selection items are native checkboxes in a real form.

## Guidance

### When to use

Use a list when people must scan, select, or navigate through a group of related rows — settings, actions, or collections.

### When not to use

Do not use a list for dense columnar data that needs sorting, filtering and pagination — a [Data table](/components/data-table) is built for that. For a single choice from many options, a [Select](/components/select) collapses the list; for a compact command menu anchored to a trigger, use a [Menu](/components/menu).

### Usability

- Navigation items wrap a real link (`<li>` → `<a href>`); selection items wrap a native checkbox in a `label`.
- Row heights are 56 px (one-line), 72 px (two-line), and 88 px (three-line) via the modifiers.
- Selection submits with a normal form — no JavaScript required for multi-select.

### Accessibility

- Keep the native elements: the `<ul>`/`<li>` carry list semantics, links are real anchors, checkboxes are native inputs.
- Wrap the selection checkbox in its `label` so the whole row is clickable and the accessible name matches the on-screen text.
- The state layer is decorative; state and selection are never communicated by color alone.
- Focus rings stay on the native link or checkbox (`:focus-visible`).

See [Choose the right control](/docs/choose-the-right-control) for the cross-component decision.
## Anatomy

The root carries the `ui-list` class and paints the container surface. Each `ui-list-item` is a flex row with up to three slots: an optional leading element, the text content, and an optional trailing element.

```html
<ul class="ui-list">
  <li class="ui-list-item">
    <span class="ui-list-item-icon" aria-hidden="true">[icon]</span>
    <span class="ui-list-item-text">
      <span class="ui-list-item-headline">Show your work</span>
      <span class="ui-list-item-supporting">Supporting text</span>
    </span>
    <span class="ui-list-item-icon ui-list-item-icon--end" aria-hidden="true">[icon]</span>
  </li>
</ul>
```

- **Row** — `ui-list-item`, a `56px` (one-line) flex row with `16px` horizontal padding. The `--two-line` (`72px`) and `--three-line` (`88px`) modifiers raise the height.
- **Leading** — `ui-list-item-icon`, a `24px` decorative slot on the inline-start. Icons are `aria-hidden` and unfocusable; visible text supplies the accessible name.
- **Text** — `ui-list-item-text` holds the `ui-list-item-headline` (body-large, on-surface) and optional `ui-list-item-supporting` (body-medium, on-surface-variant). Three-line items clamp the supporting text to two lines.
- **Trailing** — `ui-list-item-icon ui-list-item-icon--end`, a `24px` decorative slot on the inline-end.

## Content types

The roadmap requires a list to distinguish navigation, selection, and static content — Gelium ports each to native semantics.

- **Static** — a plain `<li class="ui-list-item">` with text content.
- **Navigation** — the item wraps a real link: `<li class="ui-list-item"><a class="ui-list-item-link" href="...">`. The anchor fills the row and is the focusable element.
- **Selection** — the item wraps a native checkbox in a `label`: `<li class="ui-list-item ui-list-item--select"><label class="ui-list-item-label"><input type="checkbox" ...>`. The checked values submit with a normal form — no JavaScript required for multi-select.

## States

The list covers `rest`, `hover`, `focus-visible`, `active`/`pressed`, `selected` (the checked checkbox), and `disabled` (on individual items). Hover and active paint the Material state layer (`ui-list-item::before`) at the shared `--ui-state-*` opacities; focus rings land on the native link or checkbox so keyboard users always see the target. Disabled items drop to `38%` opacity and are non-interactive.

## Design tokens

The `--ui-list-*` tokens are declared scoped to the root so the primitive works standalone, and the theme may override them globally.

| Token | Meaning |
| --- | --- |
| `--ui-list-container-color` | List container (surface) |
| `--ui-list-item-leading-space` | Inline-start padding (`16px`) |
| `--ui-list-item-trailing-space` | Inline-end padding (`16px`) |
| `--ui-list-item-icon-size` | Leading/trailing icon edge (`24px`) |
| `--ui-list-item-one-line-height` | One-line row height (`56px`) |
| `--ui-list-item-two-line-height` | Two-line row height (`72px`) |
| `--ui-list-item-three-line-height` | Three-line row height (`88px`) |
| `--ui-list-item-label-color` | Headline color (on-surface) |
| `--ui-list-item-supporting-color` | Supporting text color (on-surface-variant) |
| `--ui-list-item-icon-color` | Leading/trailing icon color |
| `--ui-list-item-hover-opacity` | Hover state layer opacity |
| `--ui-list-item-pressed-opacity` | Pressed state layer opacity |
| `--ui-list-item-focus-opacity` | Focus state layer opacity |
| `--ui-list-item-disabled-opacity` | Disabled row opacity |

## Progressive enhancement

There is no list component JavaScript. Navigation and selection work with scripting disabled: links navigate normally, and the selection form submits the checked checkbox values with a standard GET/POST. Without the stylesheet the markup degrades to a plain `<ul>`/`<ol>` of list items.

## Accessibility

- Keep the native elements: the `<ul>`/`<li>` carry list semantics, links are real anchors, checkboxes are native inputs with their names and values.
- Wrap the selection checkbox in its `label` so the whole row is clickable and the accessible name always matches the on-screen text.
- The state layer is decorative; state and selection are never communicated by color alone (hover changes the layer, selection flips the checkbox, focus shows a ring).
- Focus rings stay on the native link or checkbox (`:focus-visible`), so keyboard users always see where the row focus is.
- Decorative SVGs are `aria-hidden` and `focusable="false"`; the visible text supplies the accessible name.
- In forced-colors mode the container uses `Canvas`, text uses `CanvasText`, links keep `LinkText`, and disabled rows drop to `GrayText`, so content and interactivity survive without color.

## Keyboard

Because the interactive pieces are native, keyboard behavior comes for free: links activate with Enter, checkboxes toggle with Space, and tab order follows document order. There is no custom roving tabindex — the roadmap defers arrow-key list navigation until a platform gap is demonstrated.
