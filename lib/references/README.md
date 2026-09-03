# Portable reference catalog

This directory is shipped in the `gelium-ui` npm package so an agent can perform reference selection without cloning the monorepo or reaching the docs site.

## Resolution order

1. Read `catalog.json` and match semantic needs by `kind`, `type`, and `tags`.
2. Read the matching Markdown ficha before architecture/wireframe.
3. If development-only visual evidence is available in the consumer or supplied by the user, compare it separately; the public package intentionally does not ship unlicensed screenshots.
4. Record matched IDs, discovered affordances, product-filter decisions, and `no-match` reasons in the plan/G5 packet.

## Boundaries

These are structural references, not copied product skins. Gelium contracts, real product data, accessibility, server/no-JS behavior, and the consumer's product/design artifacts remain authoritative. Do not copy brands, logos, assets, screenshots, motion, text, prices, metrics, or identity data.

`section-references/` contains portable structural fichas. `component-references/` contains local behavior and anatomy guidance. `social-feeds/` contains a text-only pattern summary; user-supplied visual captures remain development-only in the monorepo's `docs/reference-assets/`.
