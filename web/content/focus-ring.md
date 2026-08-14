# Focus ring

Focus ring is the Gelium-wide keyboard focus contract. Every focusable element keeps a visible indicator thanks to a broad `:focus-visible` rule, while `.ui-focus-ring` is the explicit utility for custom focusable elements that are not native controls.

## The shared contract

The indicator uses the theme tokens `--ui-focus-thickness` (3 px) and `--ui-focus-offset` (2 px). A 3 px outline with a 2 px offset never changes layout geometry, so focus never causes content to shift.

## How to use it

Native links, buttons, and inputs get `:focus-visible` automatically. A custom focusable element — for example a `tabindex` container — adds the `ui-focus-ring` class to opt in to the same indicator. The live examples below all use the shared contract.