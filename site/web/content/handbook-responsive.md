# Design for screen sizes, not devices

Responsive design in Gelium starts from **viewports and content**, not “this is an iPhone” or “this is a tablet.” The same HTML and CSS should reflow cleanly from a narrow pane to a wide monitor. This page states the product stance; control choice details live in [Choose the right control](/docs/choose-the-right-control) and [Forms](/docs/forms).

## Viewports, not device names

- Design and test by **viewport width** (and height/`dvh` when it matters), not by model lists.
- Breakpoints are **layout steps** (`40rem`, `48rem`, …), not `phone` / `tablet` / `desktop` labels in product CSS.
- A foldable phone, a split-screen, and a narrow desktop window can share the same width: layout must not assume “mobile = touch + narrow” as one inseparable bundle.

## Content that reflows

- Prefer **reflow**: stacking columns, rows with `flex-wrap`, tables that scroll **inside** their container.
- Long-form reading uses the docs chrome measure: **`.prose` with `max-width: 65ch`**. Do not stretch paragraphs across the full monitor.
- App shells (for example recipes) honor **`--ui-container-max`** as a sensible column ceiling, not an excuse to force “desktop forever” width.

## Mobile-first, enhance from desktop

- Write the **default stack** for the narrow viewport: one column, stacked actions, headers that do not fight for one row.
- **Enhance from desktop** with `min-width` media queries (or equivalent layers): more columns, toolbars in a row, side-by-side panels.
- Avoid desktop-only CSS patched with opaque mobile exceptions. Recipes (admin-resource, ops-queue, public-feed) **stack headers and contain tables** at narrow widths — copy that pattern.
- Layout utilities such as `.ui-row-from-desktop` encode the same idea: column by default, row from the desktop step.

## Containment: do not mask overflow

- **`overflow-x: hidden` must not mask** a broken layout: do not put it on the document or screen body to “hide” overflow. It clips focus rings, shadows, and real content.
- In flex/grid, set **`min-width: 0`** (or `min-inline-size: 0`) on columns and children that must shrink. Without it, the child’s min-content pushes the viewport.
- Dense tables: wrap `<table>` in a local horizontal scroll container (for example `.ui-data-table-scroll`), not a clip on `html`.
- Measure real **min-content width** (DevTools, captures at ~360–400px and at the breakpoint). If something forces ~780px, the bug is layout, not the user.

## Touch targets and forms

- Interactive controls honor **`--ui-touch-target`** (usable hit-area floor; Material may raise it). Do not invent ad-hoc hit areas per screen.
- Forms: visible labels, native `type`/`inputmode`, do not block paste — see [Forms](/docs/forms).
- Pick the control for the job, not for a “mobile look”: [Choose the right control](/docs/choose-the-right-control).

## Step breakpoints

| Focus | Do | Avoid |
|---|---|---|
| Naming | Widths in `rem` or other layout units | `@media (device-width: 375px)` or “iPad only” |
| Count | Few steps where **layout shape** changes | One breakpoint per lab gadget |
| Order | Narrow base → wide enhancement | Wide base → endless `max-width` patches |
| Testing | Continuous viewports and resizable windows | Only fixed device-chrome emulators |

## Gelium tokens and pieces

| Piece | Role |
|---|---|
| `--ui-touch-target` | Touch / hit-area floor on buttons and icon-buttons |
| `--ui-container-max` | Sensible max width for app shells / recipes |
| `.ui-container` | Width-capped, padded layout column utility |
| `.ui-row-from-desktop` | Stack by default; row from the desktop step |
| `.prose` + **65ch** | Reading measure in docs and long text |
| Recipes + containment | Header stack, `min-width: 0`, local table scroll |

Themes and token vocabulary: [Themes](/docs/themes), [Tokens](/docs/tokens). Accessibility (contrast, focus, motion): [Accessibility](/docs/accessibility).

## Performance honesty

A layout that does not overflow does not replace a payload stance. JS is progressive enhancement; token and theme CSS is large **on purpose**. Measure with the same rule in [Performance](/docs/performance) and positioning in [Why Gelium](/docs/compare).

## What this page is not

- Not a catalog of device frames or a notch simulator.
- Not a license for global `overflow-x: hidden` as a “responsive fix.”
- Not a substitute for component demos or recipes — that is where real stacking shows up.

## Quick checklist

1. Does the layout make sense at ~360px width without page-level horizontal scroll?
2. Do breakpoints name **layout changes**, not hardware brands?
3. Is there `min-width: 0` where flex/grid must shrink?
4. Do tables and pre blocks scroll **inside** their box?
5. Do controls meet `--ui-touch-target` and the [Forms](/docs/forms) contract?
6. Does CSS start from a narrow stack and enhance at larger widths?
