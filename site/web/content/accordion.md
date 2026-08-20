# Accordion

Accordion is a native disclosure primitive for progressively enhanced documentation and product flows. The HTML emitted by `lib/templates/accordion.html` is the source of truth: it works with JavaScript disabled.

## Anatomy and contract

The reusable template renders one semantic `<section class="ui-accordion">` containing native `<details class="ui-accordion-item">`, `<summary class="ui-accordion-trigger">`, and a panel `<section>`. `Open` controls the server-rendered initial state. `Heading` renders a visible heading and names the root with `aria-labelledby`; when `Heading` is empty, the backwards-compatible `Label` becomes the root's `aria-label`.

Each item should have a stable `Value`. Gelium preserves it as `data-value` and derives deterministic trigger/panel IDs from it. Values are escaped as attributes; arbitrary `Body` text is escaped by `html/template` and is not promoted to unsafe HTML.

`Multiple=true` permits multiple open items and emits no exclusive grouping attribute. With `Multiple=false` and a non-empty `Name`, the same named-details `name` is emitted on each item. Named details exclusivity is a progressive browser feature: it is not universal enforcement and Gelium does not add JavaScript to simulate it. Server-rendered `open` attributes therefore remain authoritative for the initial response.

## Profiles and execution

The component exposes independent, server-validated profiles: `Behavior` is one of
`native`, `basecoat`, `material`, or `baseui`; `Execution` is `native` or `htmx`.
Unknown values safely fall back to native behavior and native execution. The
active visual theme remains the skin and does not select the behavior.

Native execution is the no-JavaScript baseline. HTMX is JavaScript and is allowed
only as an optional server-driven enhancement; this slice intentionally emits no
fake `hx-*` attributes and adds no persistence endpoint. HTMX mode therefore keeps
the same details/summary fallback when HTMX is missing.

The Theme Gallery demonstrates these independent combinations:

- `?theme=alden&behavior=baseui&execution=native`
- `?theme=vercel&behavior=material&execution=native`
- `?theme=linear&behavior=native&execution=htmx`


- The browser owns disclosure, keyboard operation, focus, and search-visible content; no JavaScript bundle is required.
- Native `<summary>` is not given `role="button"` or a redundant `aria-expanded`; native details state must not be duplicated with ARIA.
- Summary is never given a fake `disabled` attribute. Disabled items are explicitly unsupported in this native slice because `<summary>` has no disabled state. A future disabled design must render a noninteractive heading rather than pretending a summary is disabled.
- Focus-visible, reduced-motion, and forced-colors rules are provided by the core stylesheet. The 44px touch target is preserved.
- HTMX is optional progressive enhancement for lazy or persisted server content. It may replace a server-rendered panel, but it is not required for opening or closing items.

## Visual fidelity and source mapping

The Basecoat skin translates the audited Basecoat/Vega native details/summary
anatomy as a flat, full-width disclosure list: compact neutral trigger rows,
chevron affordance, and divider-separated items. Its exact runtime behaviors
are not reproduced: Gelium keeps native disclosure and does not ship Basecoat's
optional enhancement.

The Material skin adapts the MUI Expansion Panel visual language (surface,
elevation, rounded shape, spacing, and rotating expansion icon) to Gelium's
native `<details>` markup. It does not claim Material Web Accordion parity or
MUI's controlled/disabled JavaScript behavior.

The **Gelium Base UI docs-inspired reference preset** follows the visual
direction of the Base UI documentation accordion, not an upstream package
stylesheet: contiguous high-contrast rows with collapsed 1px borders, square
corners, no card shadow or elevation, a compact 44px-minimum trigger, and a `+`
that rotates 45° into an x-like open state. Root gap is reserved for the optional
semantic heading only; disclosure items stack as one bordered list. In dark mode
this resolves through the active foreground token to a light-on-dark treatment;
in light mode the same semantic tokens invert to dark-on-light. Base UI itself
remains headless and unstyled, so this is Gelium-authored CSS rather than
official Base UI CSS. The documentation demo's surrounding heading/content
framing may differ from Base UI's screenshot; Gelium keeps its optional semantic
`Heading` markup and aligns the disclosure anatomy rather than removing content
for visual mimicry.

## Source mapping and boundaries

- Basecoat native details/summary accordion informed the structure only; no Basecoat JavaScript was copied.
- Base UI `Root`/`Item`/`Header`/`Trigger`/`Panel`, plus multiple and controlled/open concepts, informed the design contract only; its documentation accordion informed Gelium's docs-inspired visual preset. No React code, API, runtime, or official Base UI CSS was copied. “Controlled” is server-initial state here, not a client-controlled React state model.
- Material Web has no official Accordion in its current component inventory, so this component makes no Material Accordion parity claim.
- HTMX remains optional. Native open/close behavior must continue to work without it.

`lib/styles/accordion.css` owns the theme-token-driven presentation; theme files may override the `--ui-accordion-*` tokens without changing this markup contract.
