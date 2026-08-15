# Focus ring

Focus ring is the Gelium-wide keyboard focus contract. Use it when a custom control needs a visible keyboard indicator. Every focusable element keeps one through the broad `:focus-visible` rule. `.ui-focus-ring` is the explicit utility for custom elements that are not native controls.

## Guidance

### When to use

Use the focus ring whenever a custom control needs a visible keyboard indicator. Every focusable element keeps one through the broad `:focus-visible` rule. `.ui-focus-ring` is the explicit utility for custom elements that are not native controls.

### When not to use

Never remove the focus indicator: a control without a visible focus ring fails the keyboard contract. Do not invent a second focus style either — use the shared tokens so focus looks identical across the system.

### Usability

- Native links, buttons, and inputs get `:focus-visible` automatically.
- A custom focusable element — for example a `tabindex` container — adds the `ui-focus-ring` class to opt in to the same indicator.
- The indicator uses the theme tokens `--ui-focus-thickness` (3 px) and `--ui-focus-offset` (2 px).

### Accessibility

- Every focusable element keeps a visible keyboard indicator.
- The 3 px outline with a 2 px offset never changes layout geometry, so focus never causes content to shift.
- Focus rings map to `Highlight` in forced-colors mode.

See [Choose the right control](/docs/choose-the-right-control) for the cross-component decision.
## The shared contract

The indicator uses the theme tokens `--ui-focus-thickness` (3 px) and `--ui-focus-offset` (2 px). A 3 px outline with a 2 px offset never changes layout geometry, so focus never causes content to shift.

## How to use it

Native links, buttons, and inputs get `:focus-visible` automatically. A custom focusable element — for example a `tabindex` container — adds the `ui-focus-ring` class to opt in to the same indicator. The live examples below all use the shared contract.