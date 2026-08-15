# Progress

Progress is the native `progress` element decorated with the Material track and its indicator from the `--ui-progress-*` tokens — the real element stays in the document, so its value semantics, `aria-valuenow`/`aria-valuetext` exposure, and determinate/indeterminate behavior work unchanged. Use it when an operation is in flight: determinate when the task has a measurable duration, indeterminate when only activity is known.

## Guidance

### When to use

Use progress when an operation is in flight: determinate when the task has a measurable duration, indeterminate when only activity is known.

### When not to use

Do not use progress for static or instant feedback — a [Toast](/components/toast) or an inline alert communicates results better. If the task can be cancelled or paused, pair progress with an [Icon button](/components/icon-button).

### Usability

- Determinate progress carries a `value`/`max` pair; indeterminate is a `progress` without `value`.
- Keep the native element in the document so its value semantics and `aria-valuenow` exposure work unchanged.
- The browser paints the indeterminate animation; Gelium only re-skins the track.

### Accessibility

- Keep the native element: the `progress` element carries its name, value, and min/max semantics at no cost.
- Always pair the element with an `aria-label` (or a linked visible label) that says what is progressing.
- Never rely on color alone: the indicator length carries the value, and for indeterminate the native motion is the signal.
- In forced-colors mode the track repaints as `Canvas` and the indicator as `CanvasText`.

See [Choose the right control](/docs/choose-the-right-control) for the cross-component decision.
## Anatomy

The component is a single decorated progress element inside a `.ui-progress` wrapper.

```html
<div class="ui-progress">
  <progress id="upload" max="100" value="65" aria-label="Upload progress">65%</progress>
</div>
```

- **`progress`** — the native element, `appearance: none` only for styling. The browser still reflects its `value`/`max` and toggles the indeterminate state.
- **Track** — the element's track paints `--ui-progress-track` with the `--ui-progress-track-height` and `--ui-progress-radius` shape.
- **Indicator** — WebKit draws the value through `::-webkit-progress-value`, Firefox through `::-moz-progress-bar`, both painted `--ui-progress-indicator`.

## States

Progress covers `determinate` (a `value`/`max` pair) and `indeterminate` (a `progress` without `value`). The browser exposes real progress to assistive tech and the platform paints the indeterminate animation; Gelium only re-skins the track so the motion stays the native one.

## Design tokens

All progress paints are token-driven so the component survives light and dark schemes.

| Token | Meaning |
| --- | --- |
| `--ui-progress-track-height` | Track thickness (`4px`) |
| `--ui-progress-radius` | Track corner radius (`--ui-radius-full`) |
| `--ui-progress-track` | Track paint (surface-container-highest) |
| `--ui-progress-indicator` | Indicator paint (primary) |

## Accessibility

- Keep the native element: the `progress` element carries its name, value, and min/max semantics at no cost, and indeterminate state is announced correctly.
- Always pair the element with an `aria-label` (or a linked visible label) that says what is progressing — matching what is on screen.
- Never rely on color alone: the indicator length carries the value, and for indeterminate the native motion is the signal.
- In forced-colors mode the track repaints as `Canvas` and the indicator as `CanvasText`, so progress survives without color.