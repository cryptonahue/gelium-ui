# Slider

Slider is a native `input[type="range"]` painted with the Material track, its active fill, and the handle from the `--ui-slider-*` tokens — the real control stays the focusable, operable element. Use a slider when a single value sits on a continuous or discrete ordered range — volume, brightness, a price cap — and approximate precision is acceptable. Form semantics, keyboard input, and assisted input work unchanged.

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

## What states can a slider be in?

The slider covers `unselected` (value at minimum), `populated` (an intermediate value), and `disabled`. Disabled follows the Material contract: the track and handle repaint through `--ui-slider-disabled` at `--ui-slider-disabled-opacity` (`.38`) and interaction stops.

## When to use it

Use a slider to pick a single value from a continuous or a discrete ordered range — volume, brightness, or a price cap. For a small set of discrete labelled options a `select` or radio group is usually clearer. This component covers the **single-value** form; a dual range select is deferred.

## Design tokens

All slider paints are token-driven so states survive light and dark schemes.

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

## Progressive enhancement

The slider is styled `input[type="range"]`; without the stylesheet it degrades to a normal range input. The `data-ui-slider` listener in `app.js` is an optional enhancement that updates `--ui-slider-fill` while dragging on WebKit — the no-JS flow shows the served value and remains fully operable, and Firefox fills natively.

## Accessibility

- Keep the native element: the range input keeps its role, value, min/max, keyboard behavior, and accessible name at no cost.
- Always pair the input with an `aria-label` (or a linked visible label), matching what is on screen.
- Never rely on color alone: the handle position carries the value, and disabled state is announced by the platform.
- The focus ring stays on the native input (`:focus-visible`), so keyboard users always see where the control is.
- In forced-colors mode the track and handle repaint as `CanvasText`, disabled as `GrayText`, and the focus ring uses `Highlight`, so the value survives without color.

## Keyboard

The range input inherits native keyboard behavior: arrow keys step by `step`, Home/End jump to the edges, and Page Up/Page Down step by larger increments — no component JavaScript required.