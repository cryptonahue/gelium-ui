# Gelium UI behavior and execution contract

Gelium separates three concerns on interactive components:

- **Skin** is visual appearance: the existing allowlisted `theme-*` class and its
  design tokens. A skin must not change the semantic component contract.
- **Behavior** is interaction policy: Accordion profiles are `native`, `basecoat`,
  `material`, and `baseui`. These profiles are allowlisted typed values and may
  change only native HTML/CSS presentation or progressive native capabilities.
- **Execution** is where an optional enhancement runs: `native` or `htmx`.

## Accordion baseline

`native` execution is the baseline and is fully functional with JavaScript
 disabled. Native `<details>` and `<summary>` own disclosure, keyboard operation,
focus, and initial `open` state. Multiple-open mode uses normal details; exclusive
mode may use the browser's named-details `name` attribute when supported. No
client runtime is used to simulate these rules.

`htmx` execution is an optional server-driven enhancement profile. It is allowed
because HTMX is JavaScript, but it is never a requirement for opening or closing
an Accordion. This first slice exposes the validated execution profile and its
fallback contract without adding fake `hx-*` attributes or a persistence endpoint.
Therefore an HTMX request still renders the same native details/summary markup and
works unchanged when HTMX is absent.

### Behavior policies

The profiles are real native policies, not cosmetic metadata. `native` preserves
the view's `Multiple`, `Name`, and initial `Open` values. `baseui` preserves the
multiple-open/headless native composition. `basecoat` and `material` default to
a native exclusive named-details group when multiple-open was not explicitly
selected; they derive a stable group name from the accordion ID and retain only
the first initially open item. `MultipleSet` distinguishes an explicit false
from Go's zero value, while an explicit `Multiple=true` remains authoritative.

Material is Material-inspired/MUI-like, not an official Material Web Accordion.
Named-details exclusivity is browser-dependent. Browsers without support fall
back to ordinary multiple-open details behavior; Gelium does not ship a
JavaScript polyfill. This boundary is intentional and preserves no-JS behavior.

The server emits only allowlisted `data-behavior` and `data-execution` values plus
corresponding component classes. Invalid query or view values fall back to
`behavior=native` and `execution=native`; arbitrary user attributes are never
copied into the component.

Examples:

- `/docs/themes/gallery?theme=alden&behavior=baseui&execution=native`
- `/docs/themes/gallery?theme=vercel&behavior=material&execution=native`
- `/docs/themes/gallery?theme=linear&behavior=native&execution=htmx`

## Phase 1 independent selection model

The document has four independent, allowlisted inputs. They are resolved on the
server; raw query text is never copied to an HTML class or attribute.

- **Behavior** (`native`, `material`, `basecoat`, `baseui`) is component-level
  interaction policy only. It must not own color, typography, geometry, icons,
  borders, or shadows. Its classes/data attributes remain diagnostics and policy
  hooks, not visual selectors.
- **Reference visual preset** (`auto`/`default`, `none`, `material`, `basecoat`,
  `baseui`) is a visual baseline. `auto` resolves from behavior: native→none,
  material→material, basecoat→basecoat, baseui→baseui.
- **Product skin** (`none`, `material`, `basecoat`, `baseui`, `alden`, `linear`,
  `vercel`) is a product visual overlay and wins over a reference. The current
  Phase 2 slice applies token-only Accordion adapters for every listed skin;
  it does not claim that every component has completed this migration.
- **Scheme** is `system`, `light`, or `dark`; execution remains `native` or
  optional `htmx`.

`<html>` exposes the resolved visual values as `data-gelium-reference`,
`data-gelium-skin`, and `data-gelium-scheme`. Behavior and execution stay on
component roots (for example Accordion), rather than claiming a misleading
site-wide interaction guarantee.

Legacy URLs remain supported: `theme=material` maps to reference=material,
skin=none; `theme=basecoat` maps to reference=basecoat, skin=none; and
`theme=alden` maps to reference=none, skin=alden. Explicit `reference` and
`skin` override their respective legacy-mapped values. `theme-*` classes remain
legacy CSS adapters in the **reference** layer during transition. For an
explicit reference/skin render, the later `data-gelium-skin` adapter overrides
both the selected reference and those legacy class values; legacy classes cannot
silently defeat the explicit skin.

### Phase 2 Accordion cascade

The actual Accordion cascade is ordered CSS layers, not a convention:

```text
`gelium.core` (semantic native markup, neutral token fallbacks, 44px floor)
→ behavior policy (server/native details semantics only)
→ `gelium.reference` (`html[data-gelium-reference]` token preset)
→ `gelium.skin` (`html[data-gelium-skin]` product overlay)
→ `gelium.site` (site composition)
```

The reference and skin adapters contain tokens only; no `.theme-*` or
`.ui-accordion--behavior-*` selector chooses an Accordion appearance. A skin
wins because `gelium.skin` is later than `gelium.reference`, including the exact
`behavior=material&reference=material&skin=basecoat` case. The source and
compiled bundle tests inspect those selectors and an Accordion cascade token.

Basecoat light and dark values are Gelium translations of the audited official
Basecoat/Vega CSS in `docs/audit/basecoat-official/`; no Basecoat package, class,
or runtime is used. Vega's dense summary has no touch minimum, so the adapter is
intentionally adapted: Accordion controls retain Gelium's `>=44px` core floor.

Base UI is headless/unstyled upstream. **Gelium Base UI reference preset** is
Gelium's own visual reference direction, not official Base UI CSS, a visual pack,
or parity claim.
