# Acknowledgments

Gelium UI is inspired by, not a copy of, the design systems and component libraries below: we adapted their ideas, patterns, and vocabulary to the Gelium model — server-rendered, zero-JavaScript, token-driven. This page records what we took from each source and how we adapted it.

## What this page is

Every library stands on the shoulders of the systems its authors have read, used, and admired. This page names ours honestly: the sources of inspiration for Gelium UI, what each contributed, and how the idea was changed to fit the Gelium model. Gelium UI ships no code from these projects; it adapts ideas, patterns, and vocabulary into an original semantic token system, server-rendered components, and a contract-tested documentation site.

| System | What we took | How we adapted it | License |
|---|---|---|---|
| [Material Design 3](https://m3.material.io) | Color roles (primary, surface, on-color), state layers, elevation system, type scale discipline | Original semantic `--ui-*` tokens with neutral core defaults; themes map brand palettes onto the same roles; state layers via `color-mix`; elevation as shadow tokens | Apache-2.0 |
| [USWDS](https://designsystem.digital.gov) | Guidance structure, the typographic measure (~68ch), accessibility patterns | Prose column capped at 65ch; every page carries when-to-use guidance; AA contrast is a tested token contract, not a guideline | CC0-1.0 (public domain) |
| [GOV.UK Design System](https://design-system.service.gov.uk) | "When to use / when not to use" rules, plain-language content discipline | Normalized Guidance sections (When to use / When not to use / Usability / Accessibility) on every component page | MIT and OGL |
| [Mozilla Protocol](https://protocol.mozilla.org) | Type scale discipline, design principles | `--ui-type-*` token families drive every text style; the four Design principles are enforced by tests | MPL-2.0 |
| [Base UI](https://base-ui.com) | Headless behavior vocabulary, handbook information architecture, FAQ landing | Handbook concept-before-reference IA; behaviors documented as server contracts instead of JS hooks | MIT |
| [Basecoat UI](https://basecoatui.com) | Open-code component approach, demos-first landing | Dogfooded docs (every page renders the real component it documents); theme layer with a second shipped direction | See project site |
| [Naive UI](https://naiveui.com) | Demos-then-API docs pattern | Component pages lead with a live preview before the API contract | MIT |
| [Name That UI](https://namethatui.com) | Alternate names plus agent prompts | "Alternative names" and "Agent prompt" sections make the docs usable by AI coding agents | See project site |
| [Material Web](https://github.com/material-components/material-web) | Token inventory reference | `--ui-*` token families audited against Material Web's public token families; an original mapping, not a port | See project site |
| [Material Symbols](https://github.com/google/material-design-icons) | Icon glyph vocabulary | Curated set of 21 trusted server-resolved inline SVGs (see the [Icon](/components/icon) page), themeable via `currentColor`, embedded in the binary | Apache-2.0 |
| [Tabler Icons](https://github.com/tabler/tabler-icons) | Outline and filled icon vocabulary | Optional curated catalog alongside Material Symbols; icons remain literal, server-resolved, and extracted only when used | MIT |
| [shadcn/ui and templ](https://ui.shadcn.com) | Distribution presets as a reference | Component source ships as copyable, themeable primitives over semantic tokens, rendered server-side | See project sites |
| [Nielsen Norman Group](https://www.nngroup.com) | Usability and reading research | F-pattern, banner blindness, hierarchy, and recovery guidance are translated into Gelium checks; no NNG code or assets ship | Research source; terms apply |
| [Gentle AI](https://github.com/Gentleman-Programming/gentle-ai) | Agent routing, delegation boundaries, optional SDD, and receipt/review ideas | Gelium uses outcome-first routing and separated delivery authority while keeping its own UX contracts; full RDD is not adopted as a requirement | MIT |
| [Refero](https://styles.refero.design) and visual studies | Visual direction references for Alden, Linear, and Vercel | Token-only skins and documented structural observations; no screenshots, branding, or copied visual assets ship | See project sites |

## License notes

The license column above is the license of the source system, noted for attribution — it does not extend to Gelium UI, which is original work. Where a project's license is not verified here, we point to the project site rather than guess.

Ideas and vocabulary are not copyrightable, and every entry above was changed to fit the Gelium model: semantic tokens instead of hardcoded palettes, server-rendered markup instead of client JavaScript, and contract tests instead of style guides. If you see a pattern here that looks familiar, that is the point — good design systems teach their readers, and we learned from these.
