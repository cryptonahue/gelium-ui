# Portable reference catalog

This directory is shipped in the `gelium-ui` npm package so an agent can perform
reference selection **without** cloning the monorepo or reaching the docs site.

## Resolution order

1. Read `catalog.json` and match semantic needs by `kind`, `type`, and `tags`.
2. Read the matching Markdown ficha **before** architecture/wireframe.
3. Optional supplement only: monorepo `docs/reference-assets/` screenshots or the
   live `/docs/section-references` pages — never required when the portable ficha
   is present.
4. Record matched IDs, discovered affordances, product-filter decisions, and
   `no-match` reasons in the plan/G5 packet.

## What ships vs what does not

| Ships in npm | Does not ship |
|---|---|
| `catalog.json` | Third-party screenshots / PNGs |
| Full structural Markdown fichas | Brand assets, logos, copied copy |
| Component patterns (feed, shell, detail-entry) | Unlicensed visual captures |

Screenshots remain development-only under the monorepo's
`docs/reference-assets/` because copyright/terms are unknown. Agents must work
from the **text fichas** in this package; visual captures are an optional
extra when the operator has lawful access.

## Layout

- `section-references/` — portable section fichas (article, 404, auth, faq, hero, pricing)
- `component-references/` — local anatomy and behavior (social-feed, shell, detail-entry)
- `section-references.md` — short index of section fichas

## Boundaries

These are structural references, not copied product skins. Gelium contracts,
real product data, accessibility, server/no-JS behavior, and the consumer's
product/design artifacts remain authoritative. Do not copy brands, logos,
assets, screenshots, motion, text, prices, metrics, or identity data.
