# Gelium UI — GEO Patterns

> Reusable patterns that implement the GEO contract (`docs/gelium-ui-geo-contract.md`). Part of Phase E. Each pattern states its intent, the concrete structure (content or markup), and what to verify. Patterns are reference templates: adapt to the page, never copy without checking the surrounding system vocabulary.
>
> Companion contract: `docs/gelium-ui-seo-contract.md` (discoverability/metadata layer). This file covers the content-and-entity layer: how a page is written, labeled and structured so a generative engine can extract an answer, an entity and a citation.

## 1. Component page pattern

**Intent**: a component documentation page whose first paragraph is a citable answer and whose sections can be paired with questions.

**Structure** (Markdown in `web/content/*.md`):

```markdown
# <Component>

<Answer-first summary>: 1–2 sentences. States what the component is, built on what
native element/contract, and when to use it. Self-contained: names the component
and the system, no back-references. → feeds meta description.

## What is a <Component>?
<Definition of the component and its role in Gelium UI; first-use term definitions.>

## When should I use a <Component>?
<Recommendation, with concrete signals. If an alternative exists, link it
(see Relations pattern).>

## How does <feature> work?
<One question per notable feature (sorting, no-JS flow, states).>

## Anatomy
<Structure/markup; may stay a noun heading — reference section.>

## Accessibility
<Contract-level a11y notes; noun heading, reference section.>

## Design tokens
<`--ui-*` table; noun heading, reference section.>

## Keyboard
<Reference section, noun heading.>
```

**Existing model**: `data-table.md` already follows intent → anatomy → states → when-to-use; only the heading question-form (§3 of the contract) and the forced summary are new.

**Verify**: first paragraph contains `<Component>` + `Gelium UI`; every `## ` after the summary is either a question or a listed reference-section noun heading; no undefined system term before its definition.

## 2. Entity pattern

**Intent**: one unambiguous, consistent system entity across every page and every metadata surface.

**Entity block** (single source, reused verbatim):

```text
name:           Gelium UI
softwareVersion: 0.6.4        (source: lib.AssetsVersion / lib/package.json)
license:        MIT           (source: README.md:128, LICENSE)
publisher/author: Gelium UI maintainer
```

**Rules**:
- Always `Gelium UI`. Never Gelidium, Loom UI, LoomChat, or a mix. This is the audit gap 1 fix and a Phase E blocker.
- Version comes from one source (`lib.AssetsVersion`), never re-typed per page.
- The entity appears in the same three places on every component page: visible (header/footer + provenance line), machine-readable (JSON-LD `SoftwareApplication`/`WebSite`/`Organization`), and in prose (first use in the summary, §1 of the contract).

**Verify**: `rg -i "gelidium|loomchat|loom ui"` over public surfaces returns nothing; visible version equals `softwareVersion` equals `lib.AssetsVersion`.

## 3. Provenance pattern

**Intent**: every page states where its content comes from and when it was last true.

**Visible line** (rendered in the `article`, under the summary — Phase F `Article` slot):

```text
Gelium UI documentation · v0.6.4 · MIT license
Published <datePublished> · Updated <dateModified> · Source: <slug>.md
```

**Metadata fields** per document (source: YAML frontmatter or Go table keyed by slug — one of the two, decided in Phase E):

```yaml
title: Button
description: <first paragraph, verbatim>
version: 0.6.4
published: 2026-08-10
updated: 2026-08-10
author: Gelium UI maintainer
```

**Rules**:
- `dateModified` changes whenever content changes; it is a release-blocking signal if stale.
- The description is *the* summary paragraph, not a re-written marketing sentence.
- Provenance is visible, not metadata-only (contract §15).

## 4. Breadcrumb / relations pattern

**Intent**: declare each page's place in the hierarchy visually and structurally.

**Visual** (every component page, zero JS):

```html
<nav aria-label="Breadcrumb">
  <ol>
    <li><a href="/">Home</a></li>
    <li><a href="/docs">Docs</a></li>
    <li aria-current="page">Button</li>
  </ol>
</nav>
```

**Related links**: close each `When should I use a <Component>?` section with a link to the rejected alternative (existing model: `data-table.md:86` → List).

**JSON-LD** — `BreadcrumbList` mirrors the visual order:

```json
{
  "@context": "https://schema.org",
  "@type": "BreadcrumbList",
  "itemListElement": [
    { "@type": "ListItem", "position": 1, "name": "Home", "item": "https://example.com/" },
    { "@type": "ListItem", "position": 2, "name": "Docs", "item": "https://example.com/docs" },
    { "@type": "ListItem", "position": 3, "name": "Button", "item": "https://example.com/components/button" }
  ]
}
```

**Verify**: visual breadcrumb and `BreadcrumbList` never diverge; the hierarchy matches `docsSections` (`docs.go:9-78`).

## 5. JSON-LD snippets

**Intent**: declarative, server-emitted structured data; zero JS; one source of values. Emitted via `pageView.JSONLD template.HTTP` pre-`</head>` (see contract §14).

**SoftwareApplication** (each component page):

```json
{
  "@context": "https://schema.org",
  "@type": "SoftwareApplication",
  "name": "Gelium UI",
  "applicationCategory": "DeveloperApplication",
  "operatingSystem": "Any",
  "softwareVersion": "0.6.4",
  "license": "https://spdx.org/licenses/MIT.html",
  "url": "https://example.com/components/button",
  "publisher": { "@type": "Organization", "name": "Gelium UI" }
}
```

**BreadcrumbList**: see pattern 4.

**Article / TechArticle** (component page, wraps the dogfooded content):

```json
{
  "@context": "https://schema.org",
  "@type": "TechArticle",
  "headline": "Button",
  "about": "Button component",
  "datePublished": "2026-08-10",
  "dateModified": "2026-08-10",
  "author": { "@type": "Organization", "name": "Gelium UI" },
  "publisher": { "@type": "Organization", "name": "Gelium UI" },
  "mainEntityOfPage": "https://example.com/components/button"
}
```

**WebSite** (home `/`):

```json
{
  "@context": "https://schema.org",
  "@type": "WebSite",
  "name": "Gelium UI",
  "url": "https://example.com/",
  "inLanguage": "en",
  "publisher": { "@type": "Organization", "name": "Gelium UI" }
}
```

**Rules**: `softwareVersion`/`datePublished`/`dateModified` come from the entity/provenance source, never literal in the template; URLs use the configurable base URL (canonical source), never `r.Host`.

## 6. Boundaries pattern

**Intent**: decide, per surface, what is public and citable vs. what must never appear.

| Publish (public, citable) | Never publish |
|---|---|
| Component docs, system contracts | Server internals, environment details, credentials |
| Demos **clearly labeled as demos** | Real user/traffic data, account state |
| Version, license, authorship, dates | Internal roadmap/handoff decisions (unless intentionally published) |
| Mock data in demos | Production data, private identifiers |

**Rules**:
- Demos keep `robots: noindex` (they are stateful surfaces, not canonical content).
- Any content not meant to be a public contract stays in `docs/handoffs/` and private docs, never in `web/content/`.
- The public site must remain renderable with zero app data: all demo rows are synthetic (`demo_whatsapp.go:224-225`).

**Verify**: grep public templates/content for real identifiers, tokens or environment strings; demos carry `noindex`; `web/content/` contains only intended public pages.

## 7. Freshness pattern

**Intent**: docs and released versions cannot drift; staleness is visible and release-blocking.

- **Visible version**: header/footer shows `Gelium UI v0.6.4` on every docs-shell page from `lib.AssetsVersion`; assets keep `?v={{.AssetsVersion}}`.
- **Versioned assets**: keep `?v=<release>` cache-busting (`layout.html:7-9`); it is the SEO asset-cache contract (audit §4), independent of visible version.
- **Content-with-code**: when a component's contract changes, the component's summary + `dateModified` change in the same work unit (`work-unit-commits`) — docs are not updated in a separate sweep.
- **Release checklist**: bump `package.json` version → update visible footer + `softwareVersion` JSON-LD → refresh touched pages' `dateModified` → run the acceptance checks.

**Verify**: visible version == `package.json:3` == JSON-LD `softwareVersion`; no page's `dateModified` predates its last content change.

## Acceptance

- [ ] A component page copy-pasted into a reader: first paragraph reads as the answer; each `## ` is a question or a listed reference section; entity/provenance/breadcrumb are visible.
- [ ] JSON-LD validates (SoftwareApplication + BreadcrumbList + Article/WebSite); values match the entity block and the visible page.
- [ ] No brand residue in public surfaces; demos `noindex`; no private data published.
- [ ] Contract and patterns are referenced from `composition-rules.md` and `screen-recipes.md`.
