# Progress

Progress is the native `progress` element decorated with the Material track and its indicator from the `--ui-progress-*` tokens. The real element stays in the document, so its value semantics, `aria-valuenow`/`aria-valuetext` exposure, and determinate/indeterminate behavior work unchanged. Use it when an operation is in flight: determinate when the task has a measurable duration, indeterminate when only activity is known.

## Guidance

### When to use

Use progress when an operation is in flight. Use determinate when the task has a measurable duration, indeterminate when only activity is known. Use the circular variant for compact surfaces where a full-width track does not fit — a card action, an inline refresh.

### When not to use

Do not use progress for static or instant feedback — a [Toast](/components/toast) or an inline alert communicates results better. If the task can be cancelled or paused, pair progress with an [Icon button](/components/icon-button).

### Usability

- Determinate progress carries a `value`/`max` pair; indeterminate is a `progress` without `value`.
- Keep the native element in the document so its value semantics and `aria-valuenow` exposure work unchanged.
- The browser paints the indeterminate animation; Gelium only re-skins the track.
- For the circular variant, put the ring inside a `role="status"` wrapper. Name the wrapper with an `aria-label` or a linked visible label; the ring stays `aria-hidden`.
- Under `prefers-reduced-motion` the ring stops as a static arc. The status region still announces, so the signal never depends on motion.

### Accessibility

- Keep the native element: the `progress` element carries its name, value, and min/max semantics at no cost.
- Always pair the element with an `aria-label` (or a linked visible label) that says what is progressing.
- Never rely on color alone: the indicator length carries the value, and for indeterminate the native motion is the signal.
- The circular variant needs an explicit live region. `role="status"` announces the operation without JavaScript; the decorative ring stays `aria-hidden`.
- In forced-colors mode the track repaints as `Canvas` and the indicator as `CanvasText`. The circular ring repaints as `CanvasText` with its active arc in `Highlight` — progress survives without color.

See [Choose the right control](/docs/choose-the-right-control) for the cross-component decision.
## Specimen

The native determinate and indeterminate bars render live on this page. The circular variant is a decorative `ui-progress-circle` ring inside a `role="status"` wrapper. Its live specimen renders below.

<div class="specimen-block">
<p><span class="ui-progress ui-progress--circular" role="status"><span class="ui-progress-circle" aria-hidden="true"></span></span></p>
</div>

## Anatomy

The component is a single decorated progress element inside a `.ui-progress` wrapper. The circular variant is the same wrapper with the `--circular` modifier, holding an `aria-hidden` ring and a `role="status"` region.

```html
<div class="ui-progress">
  <progress id="upload" max="100" value="65" aria-label="Upload progress">65%</progress>
</div>
```

```html
<div class="ui-progress ui-progress--circular" role="status" aria-label="Loading search results">
  <span class="ui-progress-circle" aria-hidden="true"></span>
</div>
```

- **`progress`** — the native element, `appearance: none` only for styling. The browser still reflects its `value`/`max` and toggles the indeterminate state.
- **Track** — the element's track paints `--ui-progress-track` with the `--ui-progress-track-height` and `--ui-progress-radius` shape.
- **Indicator** — WebKit draws the value through `::-webkit-progress-value`, Firefox through `::-moz-progress-bar`, both painted `--ui-progress-indicator`.
- **`ui-progress--circular`** — the circular indeterminate modifier on the wrapper. It is an inline-flex box sized by `--ui-progress-circular-size`, with `role="status"` naming the operation.
- **`ui-progress-circle`** — the ring: a bordered circle in `--ui-progress-track` whose `--ui-progress-indicator` arc sweeps on the shared `ui-spin` keyframes. Decorative (`aria-hidden`), so the accessible name always lives on the wrapper.

## States

Progress covers `determinate` (a `value`/`max` pair) and `indeterminate` (a `progress` without `value`). The browser exposes real progress to assistive tech and the platform paints the indeterminate animation. Gelium only re-skins the track so the motion stays the native one. The circular variant is indeterminate by design — a `role="status"` region around a decorative ring, sized for compact surfaces.

## Design tokens

All progress paints are token-driven so the component survives light and dark schemes. The circular variant reads two scoped tokens on the wrapper (no theme tokens needed).

| Token | Meaning |
| --- | --- |
| `--ui-progress-track-height` | Track thickness (`4px`) |
| `--ui-progress-radius` | Track corner radius (`--ui-radius-full`) |
| `--ui-progress-track` | Track paint (surface-container-highest) |
| `--ui-progress-indicator` | Indicator paint (primary) |
| `--ui-progress-circular-size` | Ring diameter (scoped default `--ui-size-item`) |
| `--ui-progress-circular-stroke` | Ring stroke thickness (scoped default `4px`) |

## Accessibility

- Keep the native element: the `progress` element carries its name, value, and min/max semantics at no cost, and indeterminate state is announced correctly.
- Always pair the element with an `aria-label` (or a linked visible label) that says what is progressing — matching what is on screen. For the circular variant the label lives on the `role="status"` wrapper.
- Never rely on color alone: the indicator length carries the value, and for indeterminate the native motion is the signal.
- Under reduced motion the ring stops as a static arc. The status region continues to announce the operation.
- In forced-colors mode the track repaints as `Canvas` and the indicator as `CanvasText`. The circular ring repaints as `CanvasText` with its active arc in `Highlight`, so progress survives without color.