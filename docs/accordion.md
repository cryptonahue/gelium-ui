# Accordion behavior profiles

Gelium Accordion is server-rendered native HTML: each item is a `<details>` with a
`<summary>`. Opening, closing, keyboard interaction, focus, and the initial
`open` state work with JavaScript disabled. Profiles change interaction policy,
not the active visual skin.

## Profiles

Visual fidelity is skin-specific and behavior remains independent: Basecoat is
a token-only reproduction of its native accordion presentation, while Material
is an MUI Expansion Panel-inspired adaptation. Material Web has no official
Accordion in its current inventory, and no React or custom runtime is included.

- **native** preserves the view's `Multiple`, `Name`, and every initial `Open`
  value. The demo is multiple-open by default.
- **baseui** preserves multiple-open/headless native composition and initial
  open state. Its **Gelium Base UI docs-inspired reference preset** maps the
  documentation accordion's flat, high-contrast disclosure anatomy into
  token-only CSS; it does not copy React, Base UI runtime behavior, or official
  Base UI CSS.
- **basecoat** defaults to one native named-details group when `Multiple` was
  not explicitly selected. It derives a stable name from the accordion ID and
  keeps only the first initial open item. No Basecoat JavaScript is copied.
- **material** follows the same native exclusive policy: stable named-details
  group, one deterministic initial open item, and no runtime. This is
  Material-inspired/MUI-like, not an official Material Web Accordion.

`MultipleSet` distinguishes an explicit `false` from the Go zero value. Set it
when a caller intentionally selects exclusive or multiple behavior; an explicit
`Multiple: true` remains authoritative for compatibility. Callers may provide a
custom `Name`; otherwise Basecoat and Material derive `<accordion-id>-group`.

Named-details exclusivity is a browser-native progressive capability and is
browser-dependent. Browsers without support fall back to ordinary details
behavior (multiple items can be opened); Gelium deliberately does not add
JavaScript to polyfill it. The server still makes the initial state deterministic
for Basecoat and Material.

Behavior and execution are independently allowlisted. Invalid behavior values
fall back to `native`; invalid execution values fall back to `native`. HTMX, when
selected, remains optional enhancement and does not add fake `hx-*` attributes or
replace native disclosure. Theme/skin query parameters only change CSS tokens and
never alter these policies.

The Base UI documentation demo and Gelium may frame the disclosure list with
different surrounding content. Gelium retains the optional semantic Accordion
`Heading`; visual alignment applies to the root/item/trigger/panel anatomy, not
to deleting heading markup based on a behavior profile.
