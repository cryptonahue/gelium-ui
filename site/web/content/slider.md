# Slider

Slider is a native `input[type="range"]` painted with the Material track, its active fill, and the handle from the `--ui-slider-*` tokens. The real control stays the focusable, operable element. Use a slider when a single value sits on a continuous or discrete ordered range. Approximate precision is acceptable — volume, brightness, a price cap. Form semantics, keyboard input, and assisted input work unchanged.

## Guidance

### When to use

Use a slider when a single value sits on a continuous or discrete ordered range. Approximate precision is acceptable — volume, brightness, a price cap.

### When not to use

Do not use a slider when exact values matter — price, dates, account numbers. A [Text field](/components/text-field) is accurate where the user cannot aim precisely. For a small set of discrete labelled options, a [Select](/components/select) or [Radio](/components/radio) group is usually clearer.

### Usability

- Put the visual `--ui-slider-fill` percentage on the `.ui-slider` wrapper and keep the native range input inside it.
- Always pair the input with an `aria-label` (or a linked visible label) matching what is on screen.
- Arrow keys step by `step`, Home/End jump to the edges, and Page Up/Page Down step by larger increments — natively.
- Use `--ticks` for discrete (`step`) ranges; `--value-label` shows the exact value on focus or press.
- For the value label, mirror the served value in `data-value` on the wrapper. It drives the bubble via `attr(data-value)` with no JavaScript.

### Accessibility

- Keep the native element: the range input keeps its role, value, min/max, keyboard behavior, and accessible name at no cost.
- Never rely on color alone: the handle position carries the value, and disabled state is announced by the platform.
- The focus ring stays on the native input (`:focus-visible`), so keyboard users always see where the control is.
- The value label bubble is a visual companion to the native `aria-valuenow`. Announced information never depends on the bubble.
- In forced-colors mode the track and handle repaint as `CanvasText` and the focus ring uses `Highlight`. Ticks become notches cut into the track; the bubble becomes a `CanvasText` chip with `Canvas` ink.

See [Choose the right control](/docs/choose-the-right-control) for the cross-component decision.
## Specimen

The two M3-parity variants render live on this page. Ticks mark the discrete rhythm of a `step` range.

<div class="specimen-block">
<div class="ui-slider ui-slider--ticks" style="--ui-slider-fill: 40%"><input type="range" min="0" max="10" step="1" value="4"></div>
</div>

The bubble appears on focus or press.

<div class="specimen-block">
<div class="ui-slider ui-slider--value-label" data-value="65" style="--ui-slider-fill: 65%"><input type="range" value="65"></div>
</div>

## Anatomy

The component is a single decorated range input. Put the visual `--ui-slider-fill` percentage on the `.ui-slider` wrapper and keep the input inside it.

```html
<div class="ui-slider" data-ui-slider style="--ui-slider-fill: 65%">
  <input type="range" id="volume" min="0" max="100" value="65" aria-label="Volume">
</div>
```

- **Range input** — the native control, `appearance: none` only for styling. The browser still positions the thumb exactly at the current value, so keyboard and pointer input need no component JavaScript.
- **Track** — the input's track paints the Material baseline: the active portion fills from the start to `--ui-slider-fill`, the remainder uses `--ui-slider-inactive`. WebKit draws the fill as a linear-gradient keyed to `--ui-slider-fill`; Firefox fills natively through `::-moz-range-progress`.
- **Handle** — the thumb is a 20px circle in `--ui-slider-handle` with `--ui-slider-handle-elevation`, growing to `--ui-slider-handle-pressed-size` while pressed.
- **`--ui-slider-fill`** — a percentage custom property. Set it inline to the value you render. The `data-ui-slider` enhancement in `app.js` keeps it in sync while dragging; without JavaScript the fill reflects the served value.
- **`ui-slider--ticks`** — the discrete-slider modifier on the wrapper. It paints a repeating tick rhythm over the track (`--ui-slider-tick-*` scoped tokens). Ticks are decorative — the step geometry belongs to the native input — and disappear when the input is disabled (`:has(input:disabled)`).
- **`ui-slider--value-label`** — the value-indicator modifier on the wrapper. Its `::after` reads `data-value` and floats over the handle while the input is `:focus-visible` or `:active`. Its `left` follows `--ui-slider-fill`.

## What states can a slider be in?

The slider covers `unselected` (value at minimum), `populated` (an intermediate value), and `disabled`. Disabled follows the Material contract: the track and handle repaint through `--ui-slider-disabled` at `--ui-slider-disabled-opacity` (`.38`) and interaction stops. The ticks variant hides its marks when disabled. The value label appears on `focus-visible` and `active` — the same moments the handle is operated — and never on rest.

## Design tokens

All slider paints are token-driven so states survive light and dark schemes. The variants read scoped tokens on the wrapper (no theme tokens needed).

| Token | Meaning |
| --- | --- |
| `--ui-slider-track-height` | Track thickness (`4px`) |
| `--ui-slider-track-radius` | Track and handle corner radius (`--ui-radius-full`) |
| `--ui-slider-handle-size` | Handle diameter (`20px`) |
| `--ui-slider-handle-pressed-size` | Handle diameter while pressed (`24px`) |
| `--ui-slider-active` | Active (filled) track paint (primary) |
| `--ui-slider-inactive` | Inactive track paint (surface-container-highest) |
| `--ui-slider-handle` | Handle paint (primary) |
| `--ui-slider-handle-elevation` | Handle shadow |
| `--ui-slider-disabled` | Disabled track and handle paint (on-surface) |
| `--ui-slider-disabled-opacity` | Disabled opacity (`.38`) |
| `--ui-slider-tick-width` | Tick mark thickness (scoped default `2px`) |
| `--ui-slider-tick-interval` | Tick rhythm (scoped default `1.5rem`) |
| `--ui-slider-tick` | Tick paint (scoped default `--ui-slider-inactive`) |
| `--ui-slider-value-bubble-bg` / `--ui-slider-value-bubble-fg` | Bubble chip paints (scoped defaults primary / on-primary) |

## Progressive enhancement

The slider is styled `input[type="range"]`; without the stylesheet it degrades to a normal range input. The `data-ui-slider` listener in `app.js` is an optional enhancement that updates `--ui-slider-fill` while dragging on WebKit. The no-JS flow shows the served value and remains fully operable, and Firefox fills natively. The value label follows the same contract. With the enhancement the bubble tracks the live value; without it, the served `data-value`. The native input's `aria-valuenow` is always the source of truth.

## Accessibility

- Keep the native element: the range input keeps its role, value, min/max, keyboard behavior, and accessible name at no cost.
- Always pair the input with an `aria-label` (or a linked visible label), matching what is on screen.
- Never rely on color alone: the handle position carries the value, and disabled state is announced by the platform.
- The focus ring stays on the native input (`:focus-visible`), so keyboard users always see where the control is.
- The value bubble and tick marks are visual helpers; the native input already exposes the value. The focused handle shows the bubble without color.
- In forced-colors mode the track and handle repaint as `CanvasText`, disabled as `GrayText`, and the focus ring uses `Highlight`. Ticks cut notches into the track and the bubble becomes a `CanvasText` chip, so both variants survive without color.

## Keyboard

The range input inherits keyboard behavior: arrow keys step by `step`, Home/End jump to the edges, and Page Up/Page Down use larger increments. No component JavaScript required.