# AEO

Answer-engine optimization is a content contract, not a ranking guarantee. Make each page easy to quote without pretending that `llms.txt` controls retrieval or ranking.

## Answer-first pages

- Give each page one primary answer near the H1.
- Use descriptive headings that name the question or task.
- Keep one primary answer per page; move distinct jobs to separate pages.
- Show sources visibly for factual or normative guidance.
- Preserve stable IDs and links so humans and agents can deep-link to the answer.

## Discovery and ownership

`/llms.txt` and `/llms-ux.txt` are discovery aids and compact navigation aids, not ranking guarantees. The library provides semantic markup and media contracts. This docs site owns Gelium metadata, canonical/robots/sitemap, social tags, JSON-LD, and the agent packs. Consumers own `BASE_URL`/domain, brand, content facts, authors, dates, and social assets.

## Indexability boundary

Docs and registered blog posts are indexable only when their page policy says `index, follow`. Recipes, demos, examples, and stateful/mutation surfaces remain `noindex, nofollow`; they are not made indexable for agent discovery.

## See also

- [SEO](/docs/seo) — canonical, social, sitemap, and schema ownership.
- [`/llms-ux.txt`](/llms-ux.txt) — deterministic UX and AEO IDs.
