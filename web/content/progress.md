# Progress

Progress is a native `progress` element decorated with CSS: the browser default chrome is removed and Gelium paints the 4px Material track and its indicator from the `--ui-progress-*` tokens. The real element stays in the document, so its value semantics, `aria-valuenow`/`aria-valuetext` exposure, and determinate/indeterminate behavior work unchanged.

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

## When to use it

Use determinate progress when a task has a measurable duration — file uploads, installs, savings goals. Use indeterminate when only activity is known — boot, connecting, or a first load. If the task can be cancelled or paused, pair it with an icon button.

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