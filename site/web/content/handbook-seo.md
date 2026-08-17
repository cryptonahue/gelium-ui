# SEO

Gelium UI's SEO contract has one primary answer: the **docs site owns page metadata and indexability**; the library owns semantic markup and media behavior.

## Ownership boundary

- **Library:** semantic HTML, accessible media, and framework-neutral reusable markup.
- **Docs site:** Gelium metadata, canonical URLs, robots, sitemap, Open Graph, Twitter, and JSON-LD.
- **Consumer:** `BASE_URL`/domain, brand, factual content, authors and dates, and social assets.

Consumers must provide their own host and content facts. Do not copy the docs site's canonical URLs or authorship into a product app.

## Indexability rules

Public docs and published blog posts use `index, follow`, a clean canonical, and appear in the sitemap. Recipes, demos, examples, mutation/stateful surfaces use `noindex, nofollow` and are excluded from the sitemap. A sitemap is a filtered view of the same policy, not a second route list.

## Metadata coherence

Use one primary answer per page, descriptive headings, a visible canonical URL, and visible sources for sourced guidance. `og:url`, Twitter URL, canonical, and JSON-LD page identity must agree. Dates are emitted only when the registry or content owner supplies them; never invent a modified date. The docs site does not add frontmatter or media schema here.

## See also

- [AEO](/docs/aeo) — answer-first content and agent discovery.
- [`/llms-ux.txt`](/llms-ux.txt) — compact rules for agents.
