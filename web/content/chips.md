# Chips

Chips help people enter information, make selections, filter content, or trigger actions. They appear dynamically as a group of multiple interactive elements. Loom implements each of the four upstream variants — assist, filter, input, and suggestion — with native server-rendered semantics and no component JavaScript, instead of one generic `div`.

## Types

- **Assist** — a `<a href>` or `<button type="button">` that triggers a discrete action, such as adding an event to a calendar.
- **Filter** — a native `<input type="checkbox">` that toggles a selection; it must work without JavaScript. A selected filter chip fills with the secondary container color.
- **Input** — a chip that holds information a user entered; it has a trailing remove action (an icon button). Removal is completed server-side through a `POST` round-trip with no JavaScript.
- **Suggestion** — a `<button type="button">` (or `<a href>`) that presents a clickable suggestion.

## Anatomy

Every chip is `32px` tall with the `--ui-radius-sm` (Material corner-small, `8px`) container shape and a `label-large` text (`14px` / `20px`). Assist chips may carry a leading icon; input chips carry a trailing remove icon button.

```html
<button type="button" class="ui-chip ui-chip-assist">
  <span class="ui-chip-icon" aria-hidden="true">…</span>
  <span class="ui-chip-label">Add to calendar</span>
</button>
```

### Assist and suggestion

Both are real buttons (or links) so they are focusable, operable, and announced natively. A `disabled` assist or suggestion chip drops to `--ui-state-disabled-opacity` and is not focusable.

### Filter

The filter chip keeps a real checkbox inside its `label`:

```html
<label class="ui-chip ui-chip-filter">
  <input type="checkbox" name="filter" value="docs" checked>
  <span class="ui-chip-selected-icon" aria-hidden="true"></span>
  <span class="ui-chip-label">Docs</span>
</label>
```

The `appearance: none` checkbox paints the full Material chip container. On `:checked` the container fills with the secondary container and the selected icon (a check) appears; on `:indeterminate` the checkbox is the only visual toggle surface. Because it is a native checkbox, Space toggles it without any JavaScript and the state is submitted with the surrounding form.

### Input and removal

Input chips are removable. Loom demonstrates removal as a no-JS server round-trip:

```html
<form method="post" action="/examples/chips/remove">
  <span class="ui-chip ui-chip-input">
    <span class="ui-chip-label">Star Wars</span>
    <button type="submit" class="ui-chip-remove" name="remove" value="star-wars" aria-label="Remove Star Wars">…</button>
  </span>
</form>
```

Clicking the remove icon submits the form; the server re-renders the page with that chip removed and a `role="status"` notice. This keeps the primary flow complete without JavaScript. Removal could equally be enhanced by HTMX for an in-place swap, but a plain POST already completes the action.

## States

The variants cover `rest`, `hover`, `focus-visible`, `active` (pressed), and `disabled`; the filter chip also covers `selected`. Hover and pressed apply an inset state layer of the label color (`--ui-state-hover-opacity` / `--ui-state-pressed-opacity`) exactly like the other Material controls. `focus-visible` uses the shared `--ui-focus-thickness` ring and never changes geometry, so there is no layout shift.

Disabled chips are not focusable and read the disabled opacity; a disabled selected filter chip keeps its container fill with the disabled selected container.

## Accessibility

- Never use a plain `div` for an interactive chip — assist/suggestion are buttons or links, filter is a checkbox, input uses a button for removal.
- Assist chips with only an icon carry their accessible name in `aria-label`; otherwise the visible label is the name.
- The filter checkbox keeps its native name, checked state, and keyboard behavior; the `ui-chip-selected-icon` and other visual overlays are `aria-hidden`.
- The remove button exposes its accessible name via `aria-label`; its icon is decorative and `aria-hidden`.
- Reduced-motion disables transitions. In forced-colors mode chips keep a visible boundary and the focus ring maps to `Highlight`.

## Keyboard

Because assist, suggestion, and input chips are real buttons, Enter and Space activate them by default. The filter chip inherits the native checkbox contract: Space toggles. A chip *set* can add roving arrow-key focus, which is a follow-up beyond this slice.
