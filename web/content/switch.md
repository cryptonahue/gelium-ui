# Switch

Switch is a native `input[type="checkbox"]` styled as the Material `52x32` track, its outline, and the sliding handle from the `--ui-switch-*` tokens. The real control stays the focusable, operable element. Use a switch when a setting is either on or off and should read as an immediate toggle — enabling a feature, dark mode, notifications. The state submits through a real `<form>`. Form semantics and assisted input work unchanged.

## Guidance

### When to use

Use a switch when a setting is either on or off and should read as an immediate toggle. The state submits through a real `<form>`.

### When not to use

Do not use a switch for a group of related options that submit together — [Checkboxes](/components/checkbox) are for multi-select. When the choice is one of several mutually exclusive modes, a [Radio](/components/radio) group reads better than a row of switches.

### Usability

- Keep the order `input`, `track`, `handle`, `label` exactly — the CSS uses both the adjacent and general sibling selectors to paint states.
- A switch is for one independent setting that takes effect immediately; it is not a submit-required form control.
- The track and handle spans are decorative and `aria-hidden`; put the accessible name in the label.

### Accessibility

- Keep the native element: the checkbox keeps its role, name, value, checked state, and keyboard behavior at no cost.
- The input must be nested inside its `label` (or `id`-linked) so the accessible name always matches what is on screen.
- Never rely on color alone: checked state carries the sliding handle, and disabled state is announced by the platform.
- In forced-colors mode the checked track and handle repaint as `ButtonText` and the disabled track outline drops to `GrayText`.

See [Choose the right control](/docs/choose-the-right-control) for the cross-component decision.
## Anatomy

The label wraps the native checkbox plus two decorative overlay spans. Keep the order `input`, `track`, `handle`, `label` exactly — the CSS uses both the adjacent (`+`) and general (`~`) sibling selectors to paint states.

```html
<label class="ui-switch">
  <input type="checkbox" id="dark-mode" checked>
  <span class="ui-switch-track" aria-hidden="true"></span>
  <span class="ui-switch-handle" aria-hidden="true"></span>
  <span class="ui-switch-label">Dark mode</span>
</label>
```

- **Input** — the native `ui-switch` checkbox, a `52x32` track (`--ui-switch-width` / `--ui-switch-height`) with a fully round radius (`--ui-switch-radius`) and a `--ui-switch-outline-width` border in `--ui-switch-track-outline`. It keeps `appearance: none` only for styling; keyboard focus, the accessibility tree, and form submission are unaffected.
- **Track** — `ui-switch-track`, a decorative `span` that paints the unselected fill in `--ui-switch-track-unselected` and flips to the primary `--ui-switch-track` when the input is checked.
- **Handle** — `ui-switch-handle`, a circular `span` that starts at `16px`, grows to `24px` when checked, and stretches to `28px` while pressed (`--ui-switch-handle-size`, `--ui-switch-handle-selected-size`, `--ui-switch-handle-pressed-size`). It slides to the opposite edge via `translateX`.
- **Label** — `ui-switch-label`, the visible text paired with the track. Put the accessible name here and wrap both in a `label` so clicking the text toggles the input.

## What states can a switch be in?

The switch covers `checked`, `unchecked`, and `disabled` (both disabled states shown in the preview). Disabled follows the Material 3 split-opacity contract: the track drops to `--ui-switch-disabled-track-opacity` (`.12`) while the unchecked handle drops to `--ui-switch-disabled-handle-opacity` (`.38`). A disabled checked handle keeps full opacity and repaints through `--ui-switch-disabled-handle`.

## Design Tokens

All switch paints are token-driven so states survive light and dark schemes.

| Token | Meaning |
| --- | --- |
| `--ui-switch-width` | Track width (`52px`) |
| `--ui-switch-height` | Track height (`32px`) |
| `--ui-switch-radius` | Track and handle corner radius (`--ui-radius-full`) |
| `--ui-switch-outline-width` | Track outline thickness |
| `--ui-switch-track` | Checked track paint (primary) |
| `--ui-switch-track-unselected` | Unchecked track fill |
| `--ui-switch-track-outline` | Track outline color |
| `--ui-switch-handle` | Unchecked handle paint |
| `--ui-switch-handle-selected` | Checked handle paint |
| `--ui-switch-handle-size` | Handle diameter when unchecked (`16px`) |
| `--ui-switch-handle-selected-size` | Handle diameter when checked (`24px`) |
| `--ui-switch-handle-pressed-size` | Handle diameter while pressed (`28px`) |
| `--ui-switch-disabled-track-opacity` | Disabled track opacity (`.12`) |
| `--ui-switch-disabled-handle-opacity` | Disabled unchecked handle opacity (`.38`) |
| `--ui-switch-disabled-handle` | Disabled handle paint (on-surface) |

## Progressive enhancement

The switch is pure CSS over a native checkbox — there is no component JavaScript. Without the stylesheet (or in a user style) the controls degrade to a normal checkbox group. The `ui-switch-track` and `ui-switch-handle` spans are decorative and `aria-hidden`, so nothing extra is read to assistive technology.

## Accessibility

- Keep the native element: the checkbox keeps its role, name, value, checked state, and keyboard behavior at no cost.
- The input must be nested inside its `label` (or `id`-linked) so the accessible name always matches what is on screen.
- Never rely on color alone: checked state carries the sliding handle, and disabled state is announced by the platform.
- The focus ring stays on the native input (`:focus-visible` on `input`), so keyboard users always see where the toggle is.
- In forced-colors mode the track keeps a `CanvasText` boundary, the checked track and handle repaint as `ButtonText`, and the disabled outline drops to `GrayText`. The on/off state survives without color.

## Keyboard

Because the switch inherits native checkbox keyboard behavior, Space toggles and the arrow keys do nothing special — the checkbox contract applies unchanged.