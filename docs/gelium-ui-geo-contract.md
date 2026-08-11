# Gelium UI — GEO Contract

> Contract of the system for Generative Engine Optimization (GEO): how Gelium UI documentation is made findable, understandable, summarizable, entity-relatable and citable by generative engines. Part of Phase E (`docs/gelium-ui-system-roadmap.md`), companion to `docs/gelium-ui-seo-contract.md`.
>
> Baseline evidence: `docs/handoffs/seo-geo-audit.md` (research, read-only). This document turns the audit's GEO findings into a binding contract for content, markup and metadata. It does not modify code; the server-driven metadata implementation is the Phase E engineering work referenced in the roadmap.

## Goal

Make every Gelium UI documentation page citable by generative engines, on top of an already solid SEO foundation (server-rendered HTML, clean URLs, factual content). The contract defines how each GEO dimension is implemented in the Gelium UI system — content rules, visible markup, and machine-readable metadata — so a generative engine can reliably:

1. **find** a page (crawlable, stable URL),
2. **understand** it (explicit summary, definitions, descriptive headings),
3. **summarize** it (self-contained extractable answer),
4. **relate it to entities** (a single unambiguous entity: Gelium UI, version, license, authorship),
5. **cite it** (visible provenance, machine-readable metadata).

GEO builds on SEO. Without crawlable, stable, well-metadataed pages there is nothing for a generative engine to understand.

## Source of truth

| Aspect | Source | Evidence |
|---|---|---|
| Entity name | `Gelium UI` (canonical) | Roadmap naming note (`gelium-ui-system-roadmap.md:7`); audit gap 1 |
| Version | `0.4.0` | `package.json:3`; asset query strings `layout.html:7-9` |
| License | MIT | `README.md:128`, `LICENSE` |
| Authorship | Project maintainer | Decision below (§6) |
| Content | `web/content/*.md`, embedded and dogfooded | `server.go:262-272`, `docs.go:86` |
| URLs | `componentRoutes()` registry | `routes.go:16-47` |
| Markdown pipeline | `goldmark.New()` — default, no frontmatter today | `server.go:63` |

## Contract dimensions

### 1. Answer-first content

- **Rule**: every component page opens with the answer. The first paragraph after the `#` heading is a 1–2 sentence, self-contained statement of what the component is and when to use it.
- **Gelium UI today**: partially met — `button.md:3` and `data-table.md:3` already lead with a factual definition.
- **Requirement**: the intro must not assume the heading; it must be readable in isolation as a direct answer to "What is X?".
- **Verification**: first non-empty paragraph after the H1 contains the component name and the word "Gelium UI" or a stable system term.

### 2. Explicit summaries

- **Rule**: the first paragraph is the citable summary. It is the single source for `meta description` and for `og:description` derivations.
- **Requirements**:
  - no back-references ("this", "it") that resolve only against the H1;
  - full entity name on first use in the summary ("Button is a Gelium UI component built from native HTML…");
  - one summary per page; nothing later in the page may contradict it.
- **Verification**: the summary text is extractable verbatim and still makes sense outside the page.

### 3. Descriptive headings

- **Rule**: headings work as questions/answers, not only as labels. Section headings are short questions or statement-answers that a generative engine can pair with the following block.
- **Gelium UI today**: headings are descriptive but non-interrogative ("Anatomy", "States", "When to use it" — `data-table.md:5,34,46`).
- **Requirement**: convert high-value sections to question form — `What is a data table?`, `When should I use a data table?`, `How does sorting work?`, `How does the no-JS flow work?`. Keep them short and factual; one question per heading; no chained interrogatives.
- **Exemption**: structural sections (Accessibility, Design tokens, Keyboard) keep noun headings; they are reference sections, not FAQ candidates.

### 4. Definitions

- **Rule**: define terms at first use, inline, in the same sentence or the immediate next clause. Do not assume reader knowledge of system vocabulary.
- **Required terms** (define on first use anywhere): Gelium UI, server-rendered, open-code component, dogfooding, native semantics, HTMX (enhancement only), design token (`--ui-*`), theme (e.g. theme-material).
- **Format**: `Term — definition.` or `Term is a …`. A page that uses a system term before defining it anywhere fails the contract.
- **Cross-reference**: the canonical vocabulary lives in `docs/gelium-ui-vocabulary.md`; definitions must not contradict it.

### 5. Entities

- **Rule**: one unambiguous system entity — **Gelium UI** — with version, license and authorship. No alternate brand names in public content.
- **Gelium UI today**: brand split — site renders "Gelidium UI" (`layout.html:6,13`, `index.md:1`), repo README says "Loom UI" (`README.md:1`), demos use "LoomChat" (`demo-whatsapp.html:7`, `demo-whatsapp-admin.html:7`, `demo_whatsapp.go:224-225`).
- **Requirement**: unify on `Gelium UI` in all public surfaces (title, brand, demos, content). Expose the entity machine-readably via JSON-LD `WebSite`/`Organization`/`SoftwareApplication` (§14).
- **Entity block** (single source, repeated consistently): name `Gelium UI`, softwareVersion `0.4.0`, license `MIT`, author `Gelium UI maintainer`.

### 6. Authorship

- **Decision**: docs authorship = **project maintainer**, represented as the `Organization`/publisher node in JSON-LD (`name: "Gelium UI"`, `url`, `logo` when available). No individual human name is exposed by default; the docs are maintainer-owned, single-entity output.
- **Rationale**: the site documents a single system maintained by its owner; exposing a personal name adds a second entity with no reader benefit and creates provenance ambiguity if maintainers change.
- **Requirement**: authorship is declared once at site level (home `WebSite` node) and per-page only if per-document authors exist. Today: none — this is a gap to close in Phase E.

### 7. Date / freshness

- **Rule**: every content page carries `datePublished` and `dateModified`, visible in the article and in JSON-LD.
- **Gelium UI today**: zero dates — `web/content/*.md` has no frontmatter and the HTML has no dates (`layout.html:3-10`).
- **Storage**: Markdown has no frontmatter today (goldmark default, `server.go:63`). The contract fixes the *fields*, not the storage; Phase E must choose between (a) minimal YAML frontmatter parsed in the handler, or (b) a Go metadata table keyed by slug (audit §6 proposal). Either way the rendered article and JSON-LD must emit the same values.
- **Verification**: grep confirms no page renders dates today; post-implementation every component page must show an update date distinct from the layout boilerplate.

### 8. Sources / provenance

- **Rule**: every page states where its content comes from, visibly. Provenance = origin + citation + stable reference.
- **Gelium UI today**: origin is implicit (dogfooded content embedded in the binary, `web/assets.go:8`); no visible citation block.
- **Requirements**:
  - a visible provenance line in the article (e.g. "Gelium UI documentation — v0.4.0, MIT. Source: button.md");
  - cross-links count as citations and already exist (`data-table.md:86` links to List);
  - URLs are the stable reference (`routes.go:16-47`) — provenance must reference paths, never query-state (sort/filter/pagination URLs, `data_table.go`, are state, not identity).

### 9. Relations

- **Rule**: pages declare their place in the hierarchy: breadcrumbs + related links.
- **Gelium UI today**: no breadcrumb (`docs/gelium-ui-vocabulary.md:239-243` marks the gap); hierarchy exists only in the generated `/docs` index (`docs.go:9-78`).
- **Requirements**:
  - visual breadcrumb `Home → Docs → <Component>` on every component page (semantic `<nav aria-label="Breadcrumb">` → `<ol>` → `<li>` → `<a>`, zero JS);
  - related links where content already cross-references (`When to use it` sections);
  - `BreadcrumbList` JSON-LD mirroring the visual breadcrumb (§14).

### 10. Stable URLs

- **Rule**: URLs are clean and permanent. Content identity lives in the path.
- **Gelium UI today**: met — `componentRoutes()` is the single registry (`routes.go:16-47`), no hashes, no `.html` suffix, stable query strings only for server state on data-table.
- **Requirements**:
  - never rename or repurpose a path; if a component is renamed, add a redirect and keep the old path;
  - add `<link rel="canonical">` (requires a configurable base URL — new `BASE_URL`/host config, audit §6);
  - document demos paths as non-canonical if indexed later (`robots: noindex` on `/demo/*`).

### 11. Facts

- **Rule**: content is factual, verifiable and unambiguous. Every statement maps to something observable (rendered component, token value, HTTP contract).
- **Gelium UI today**: strong — content is technical and dogfooded; the live preview is the real component (`docs.go:86`).
- **Requirements**:
  - no marketing superlatives or unverifiable claims;
  - token values quoted exactly as declared (`--ui-*` tables like `data-table.md:69-82`);
  - state the server contract in the same terms the code uses (HTTP status, `loom:toast`, GET semantics — roadmap §61);
  - do not state a behavior as universal if it is theme-specific (e.g. Material-specific hexes).

### 12. Boundaries

- **Rule**: clear separation between public and non-public content; the public site never exposes application data, internal state or credentials.
- **Gelium UI today**: demos are self-contained mock data (`demo_whatsapp.go:224-225`); no secrets rendered.
- **Requirements**:
  - publish only intended public contracts: component docs, system contracts, demos clearly labeled as demos;
  - never render server-internal identifiers, environment details, or user data;
  - demos stay synthetic; do not leak real traffic or account state into public pages;
  - internal roadmap/handoff decisions stay out of public pages unless intentionally published as system contracts.

### 13. Freshness

- **Rule**: docs stay in sync with released versions; stale content is updated, never left to drift.
- **Gelium UI today**: version exists only in asset query strings (`?v=0.4.0`, `layout.html:7-9`) — not as visible text, not in metadata.
- **Requirements**:
  - visible version on every page (header/footer, Phase F footer pattern);
  - `softwareVersion` in JSON-LD always matches the release (`package.json:3`);
  - when a component's contract changes, the component's intro summary and dates are updated in the same change (content travels with code — `work-unit-commits`);
  - an outdated page (version mismatch, stale API example) is a release-blocking defect.

### 14. JSON-LD

- **Rule**: structured data is emitted server-side, declaratively, zero JS, in every relevant page.
- **Gelium UI today**: zero — no `application/ld+json` anywhere (audit §3).
- **Types and placement**:

| Type | Page | Content |
|---|---|---|
| `WebSite` | `/` (home) | `name` Gelium UI, `url`, `inLanguage`, `publisher` |
| `Organization` | `/`, `/docs` | publisher node, `url`, `logo` (when available) |
| `SoftwareApplication` | each `/components/*` | `name`, `applicationCategory` DeveloperApplication, `softwareVersion` 0.4.0, `operatingSystem`, `license` MIT |
| `BreadcrumbList` | each `/components/*` | Home → Docs → Component |
| `TechArticle` (optional) | each `/components/*` | `headline`, `datePublished`, `dateModified`, `about`, `author` |

- **Emission**: extend `pageView` (`server.go:65-104`) with `JSONLD template.HTML`, populated in the handler and rendered pre-`</head>` in `layout.html` (audit §6). Values derive from the entity block (§5) and provenance (§7-8) — no duplicated literals in the template.

### 15. Visible citations

- **Rule**: citations are human-visible on the page, not only in metadata. A generative engine should be able to see what it is citing.
- **Gelium UI today**: no visible citation block.
- **Requirements**: on every component page, render at minimum: version, license (linked to `LICENSE`/`README.md`), repository/source reference, breadcrumb. This is the visible counterpart of the JSON-LD entity and the provenance line (§8).

### 16. Machine-readable metadata

- **Rule**: every page carries machine-readable metadata beyond JSON-LD: `meta description`, `canonical`, `robots`, Open Graph, Twitter card.
- **Gelium UI today**: zero — head is only charset/viewport/title/styles/scripts (`layout.html:3-10`).
- **Source**: description derived from the answer-first summary (§1-2); canonical from configurable base URL + path; robots default `index,follow`, `noindex` on demos and POST-only surfaces; OG/Twitter derived from title/description/canonical. Implemented server-driven per route (audit §6).
- **Frontmatter decision**: no frontmatter today (goldmark default). Phase E chooses one source of truth per document — YAML frontmatter parsed in the handler, or a Go table keyed by slug. Both must feed §7, §8, §14 and §16 from a single place.

## What GEO is not

- **GEO is not a replacement for SEO.** SEO gets the page crawled, indexed and ranked in traditional search (metadata, canonical, sitemap, robots, performance). GEO operates on top of that foundation to make the page understood and cited by generative engines. Gelium UI keeps both contracts; GEO does not subsume the SEO contract (`docs/gelium-ui-seo-contract.md`).
- **GEO does not guarantee citation.** Citation is decided by each engine's retrieval and generation pipeline. GEO only raises the probability: unambiguous entity, visible provenance, extractable answers, machine-readable metadata. No contract, however complete, can force a generative engine to cite a page.
- **GEO is not keyword stuffing.** Generative engines read for meaning and entity consistency, not keyword density. Repetition of "Gelium UI" or of component names degrades readability and summary quality. The contract optimizes for a single clear statement per page, not for term frequency.
- **GEO does not require llms.txt.** An `llms.txt` file is an optional helper for LLM-facing tooling, not a GEO requirement. **Decision: not required for Gelium UI at this stage.** The site is fully server-rendered with one canonical content URL per component (no duplication, no JS-rendered walls), so the value an `llms.txt` adds is marginal. Re-evaluate when the site grows beyond a single content surface, adds multi-format content (video, changelog archives), or gains enough pages that engine-side crawl is a bottleneck. Document the decision in the roadmap notes; do not ship `llms.txt` as part of Phase E.

## Baseline (from `docs/handoffs/seo-geo-audit.md`)

| Dimension | State today | Notes |
|---|---|---|
| Answer-first content | ◐ | `button.md:3`, `data-table.md:3` lead with a definition; not formalized |
| Explicit summaries | ◐ | Intro is citable in practice; not wired to metadata |
| Descriptive headings | ❌ | Non-interrogative (`data-table.md:5,46`) |
| Definitions | ◐ | Terms used inline; vocabulary not enforced |
| Entities | ❌ | Brand split Gelium/Gelidium/LoomChat; version only in query strings |
| Authorship | ❌ | None |
| Date / freshness | ❌ | No dates anywhere |
| Sources / provenance | ❌ | Origin implicit; no citation block |
| Relations | ❌ | No breadcrumb; hierarchy only in generated `/docs` |
| Stable URLs | ✅ | `routes.go:16-47` |
| Facts | ✅ | Factual, dogfooded content |
| Boundaries | ✅ | Synthetic demo data; no secrets |
| Freshness | ❌ | Version not visible as text/metadata |
| JSON-LD | ❌ | Zero |
| Visible citations | ❌ | None |
| Machine-readable metadata | ❌ | Zero (no description/canonical/OG) |

## Acceptance (Phase E DoD, GEO part)

- [ ] Single entity `Gelium UI` in title, brand, demos and content; no LoomChat/Loom UI/Gelidium residue in public surfaces.
- [ ] Every component page renders: answer-first summary, visible date, version, license, breadcrumb, provenance line.
- [ ] JSON-LD `WebSite` + `SoftwareApplication` + `BreadcrumbList` (+ `TechArticle`) emitted server-side, zero JS, values from one source.
- [ ] `meta description`, `canonical`, `robots`, OG, Twitter present per route.
- [ ] Goldmark content pipeline unchanged or extended only by the chosen frontmatter/metadata source; build `npm run build`, `go test ./...`, `go vet ./...` green.
- [ ] Contract is referenced from `composition-rules.md` and `screen-recipes.md` (roadmap Phase E integration requirement).
