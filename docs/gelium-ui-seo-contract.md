# Gelium UI — SEO Contract

> Contract document for the Gelium UI documentation site (physical repo `loom-ui`).
> Phase E of `docs/gelium-ui-system-roadmap.md`. Base evidence: `docs/handoffs/seo-geo-audit.md`.
>
> This document defines **what the SEO layer must guarantee** and **where each guarantee is implemented**
> (Go handler + layout + template). It is a contract, not a guide: every page that ships a `<head>`
> MUST satisfy these rules. Implementation patterns live in `docs/gelium-ui-seo-patterns.md`.

---

## 0. Non-goals and invariants

### SEO must never depend on the theme

The Markup and metadata emitted into `<head>` are **theme-agnostic**. `ThemeClass` (`layout.html:2`,
`server.go:14-16`, `themeClass` `server.go:23-30`) selects visual direction only. Switching themes
(`theme-material` → `theme-basecoat`, Phase H) must not change a single byte of SEO-relevant markup:
`<title>`, meta tags, canonical, robots, sitemap, JSON-LD and heading structure are identical under
every theme. Rationale: theme is a presentational layer; SEO is a content and server contract.

### Constraints that shape the contract

- Server-rendered Go (`net/http` + `html/template` + `embed`), Markdown rendered with goldmark
  (`server.go:282-301`), no JS framework.
- No CDN (`README.md`), assets embedded (`web/assets.go`), CSS minified.
- Clean, stable URLs (single registry `componentRoutes()`, `routes.go:16-47`).
- Zero mandatory JS: every SEO artifact (canonical, robots, JSON-LD, sitemap) is emitted in the
  HTML response, crawlable without execution.
- Single source of truth for route data is the route registry (`routes.go`) — metadata derives from
  it, never from a parallel hardcoded list.

### Brand

The canonical brand is **Gelium UI** (roadmap naming, `gelium-ui-system-roadmap.md:7`). Legacy
strings "Gelidium UI" (`layout.html:6,13`, `server.go:253`) and "LoomChat" (demos) are unification
work inside Phase E and are assumed resolved for this contract. All metadata uses the unified brand.

---

## 1. Metadata contract per route

Every GET page resolves metadata **in the handler** and emits it **in the layout**. The resolved
value set for a route is:

| Field | Type | Required | Source / resolver | Emission |
|---|---|---|---|---|
| `Title` | `string` | yes | handler (exists, `server.go:66`) | `<title>{{.Title}} · Gelium UI</title>` |
| `Description` | `string` | yes, 150–160 chars | route registry data or first factual sentence of the `.md` | `<meta name="description">` |
| `Canonical` | `string` | yes | `BASE_URL` config + `r.URL.Path` | `<link rel="canonical">` |
| `Robots` | `string` | yes | default `index, follow`; `noindex` for demos/POST pages | `<meta name="robots">` |
| `OGTitle` / `OGDescription` | `string` | yes | derived from `Title` / `Description` | `<meta property="og:title">` etc. |
| `OGType` | `string` | yes | `website` (home, /docs) or `article` (component pages) | `<meta property="og:type">` |
| `OGURL` | `string` | yes | = `Canonical` | `<meta property="og:url">` |
| `OGImage` | `string` | yes (placeholder) | default `https://gelium-ui.example/og.png`; a real static asset path replaces it when shipped | `<meta property="og:image">` |
| `TwitterCard` | `string` | yes | `summary` (default) / `summary_large_image` if `OGImage` | `<meta name="twitter:card">` |
| `JSONLD` | `template.HTML` | yes | per page type (see §12) | `<script type="application/ld+json">` |
| `Lang` | `string` | yes | `en` (fixed, `layout.html:2`); `hreflang` future (§17) | `lang="en"` |
| `ThemeClass` | `string` | yes | existing server-driven resolver | `class="..."` — presentational only |

### Where each field lives

- **Handler layer** (`internal/app/*.go`): each handler sets `Title` (already does) and calls a
  metadata resolver to populate `pageView.Meta`. `renderMarkdownStatus` (`server.go:282`) is the
  single choke point that merges `Meta` into every page before template execution, exactly like it
  already merges `Nav` (`server.go:290`) and `ThemeClass` (`server.go:292`).
- **View model** (`pageView`, `server.go:65-104`): gains a `Meta pageMeta` field. `pageMeta`
  carries the fields above; nil-safe so existing tests and partial templates keep working.
- **Registry** (`routes.go:16-47`): `componentRoute` gains optional `Description`, `Robots`,
  `OGType`, `JSONLDType`. The `metaFor(route, req)` resolver builds `pageMeta` from the route +
  `BASE_URL` + the request path. Home, `/docs` and demo pages resolve their own metadata in their
  handlers.
- **Template** (`web/templates/layout.html`): head emits each meta tag behind `{{if .Meta.X}}`
  guards, so a page with zero metadata still renders a valid head (graceful degradation).

### Example — resolver output for `/components/button`

```go
// handler (button.go)
s.renderMarkdownPage(w, pageView{
    Title: "Button",          // existing
    Meta:  metaFor(route, req), // from routes.go registry
}, "content/button.md")
```

```html
<title>Button · Gelium UI</title>
<meta name="description" content="Button is an open-code component built from native HTML. Use a <button> for actions and an <a> only for navigation.">
<link rel="canonical" href="https://gelium-ui.example/components/button">
<meta name="robots" content="index, follow">
<meta property="og:type" content="article">
<meta property="og:title" content="Button · Gelium UI">
<meta property="og:description" content="Button is an open-code component built from native HTML.">
<meta property="og:url" content="https://gelium-ui.example/components/button">
<meta property="og:image" content="https://gelium-ui.example/static/og-cover.png">
<meta name="twitter:card" content="summary_large_image">
<script type="application/ld+json">{...}</script>
```

---

## 2. `<title>` — server-driven per route (exists, kept)

**Architecture rule**: the `<title>` is the single most important on-page signal. It MUST be
server-driven per route, unique, descriptive and ≤ 60 chars before the brand suffix.

**Implementation**: each handler sets `pageView.Title` (`server.go:66`); the layout renders
`<title>{{.Title}} · Gelium UI</title>` (`layout.html:6`). Never derive `<title>` in JavaScript and
never let it depend on the theme.

**Example**:

```html
<title>Data table · Gelium UI</title>
```

---

## 3. `meta description` — per route, 150–160 chars

**Architecture rule**: one description per URL, factual, unique across the site, 150–160 characters
counted on the rendered text (HTML stripped). It is the excerpt search engines and generative
engines quote; a missing or duplicated description makes the engine pick its own text.

**Implementation**: `metaFor` takes the description from the route registry; where a page is backed
by Markdown, the description may be derived at build/handler time from the first factual sentence of
the `.md` (the audit baseline: every `web/content/*.md` starts with a factual intro, e.g.
`button.md:3`). The layout emits `<meta name="description">` in the head. Assertions in
`server_test.go` (pattern `server_test.go:72-82`) verify presence and length range.

**Example**:

```html
<meta name="description" content="A server-side data table with sorting, filtering and real pagination links for Gelium UI.">
```

---

## 4. `robots` — meta per route + robots.txt

**Architecture rule**: every GET page declares its indexing intent explicitly. Public documentation
is `index, follow`; anything that mutates state, is a demo shell, or duplicates canonical content is
`noindex, nofollow`. A `/robots.txt` denies nothing on the public docs site and explicitly
references the sitemap.

**Implementation**:
- Meta: default `index, follow`; the resolver overrides to `noindex, nofollow` for the demo routes
  (`/demo/*`) and any POST companion pages (list at `postOnlyPaths()`, `server.go:164-176`).
- `robots.txt`: a `GET /robots.txt` handler (`server.go:119-133`) serving the crawl policy as a
  static template:
  ```text
  User-agent: *
  Allow: /
  Disallow: /demo/
  Disallow: /examples/
  Disallow: /recipes/
  Sitemap: https://gelium-ui.example/sitemap.xml
  ```
  The demo, example and recipe surfaces are excluded from crawling; the sitemap is advertised.
  Registered alongside `GET /healthz` in `New()`.

**Example**:

```html
<meta name="robots" content="noindex, nofollow">   <!-- /demo/whatsapp/admin -->
<meta name="robots" content="index, follow">        <!-- /components/button -->
```

---

## 5. `sitemap.xml` — generated server-side

**Architecture rule**: the sitemap is generated from the route registry, never hand-maintained, so
it cannot drift from the actual library. Only `index, follow` GET pages are listed; canonical URLs
are absolute.

**Implementation**: `GET /sitemap.xml` handler (`server.go:154-199`) iterates `sitemapPaths()`
(`/`, `/docs` + every `componentRoutes()` entry) and emits one `<url><loc>` per page with an
absolute canonical URL and no `lastmod` (static content):

```xml
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url><loc>https://gelium-ui.example/</loc></url>
  <url><loc>https://gelium-ui.example/docs</loc></url>
  <url><loc>https://gelium-ui.example/components/button</loc></url>
  <!-- ...one <url> per registry entry... -->
</urlset>
```

`Content-Type` is `application/xml`; the document is marshaled with `encoding/xml` from the route
registry so it can never drift from the library. Referenced from `robots.txt`. Tested by
`TestSitemapXMLDerivedFromRegistry` (`server_test.go`): every `componentRoutes()` path appears
exactly once and noindex/form surfaces (`/demo/`, `/examples/`, `/recipes/`) are absent.

---

## 6. Open Graph + Twitter Cards

**Architecture rule**: every indexable page ships a minimal OG set — `og:title`, `og:description`,
`og:type`, `og:url`, `og:image` (optional) — plus `twitter:card`. These are derived from the
canonical metadata fields, never written a second time by hand (single source of truth).

**Implementation**: the layout renders the OG block from `pageMeta`; `og:type` is `website` for
home and `/docs`, `article` for component pages. `og:url` equals `Canonical`. Every layout page
resolves a default `OGImage` (`https://gelium-ui.example/og.png`, `server.go:104`) and the layout
emits `og:image` when set plus a `twitter:card` derived from it — `summary_large_image` when an
image is set, `summary` otherwise (`layout.html:13-17`). The `og.png` asset itself is not shipped
yet; the placeholder origin is the documented contract so the markup stays stable.

**Example**:

```html
<meta property="og:type" content="article">
<meta property="og:title" content="Dialog · Gelium UI">
<meta property="og:description" content="A modal dialog built on native semantics with a no-JS fallback.">
<meta property="og:url" content="https://gelium-ui.example/components/dialog">
<meta property="og:image" content="https://gelium-ui.example/og.png">
<meta name="twitter:card" content="summary_large_image">
```

---

## 7. Heading hierarchy — single h1; h2/h3 by content

**Architecture rule**: exactly one `h1` per page, provided by the content (Markdown `# ...`), never
by the layout. Sections use `h2`; sub-sections `h3`. The layout contributes **no** heading: the
brand is an `<a class="brand">` (`layout.html:13`), not a heading.

**Implementation**: goldmark renders `# ` as a single `h1` from every `web/content/*.md`
(audit-verified). The `/docs` index builds its own `# Documentation` (`docs.go:85`). Page-level
components that look like headings (e.g. a future `section-heading` public pattern) must use `h2+`
or a styled `<p>`, reserving `h1` for content. `server_test.go` asserts one `h1` per page.

**Example**: `button.md` → `<h1>Button</h1>` followed by `<h2>Variants and states</h2>`; no other
`h1` anywhere in the rendered page.

---

## 8. Internal links — nav + breadcrumbs + content, descriptive anchor text

**Architecture rule**: the site's link graph is crawlable and label-rich. Three layers compose it:
the primary `<nav>` (`layout.html:14`, from `navLinks()`), breadcrumbs (§11, Phase F), and in-content
links written as Markdown. Anchor text is descriptive ("Data table", not "here").

**Implementation**:
- Nav: real `<a href>` links driven by `navLinks()` (`routes.go:49-58`).
- Content: Markdown links render as normal anchors via goldmark (e.g. `docs.go:91`,
  `data-table.md` cross-links).
- Breadcrumbs: Phase F component (see patterns doc), a real `<nav>` + `<ol>` of links.
- No link relies on JS to be followed; HTMX `hx-*` attributes (e.g. pagination,
  `data-table.html:76`) are progressive enhancements over the same real `href`.

**Example**:

```html
<nav aria-label="Primary"><a href="/docs">Docs</a> <a href="/components/button">Button</a></nav>
<!-- content -->
<p>Compare with <a href="/components/icon-button">Icon button</a> for single-icon actions.</p>
```

---

## 9. Image alt — decorative SVG inline `aria-hidden`; content images with alt

**Architecture rule**: decorative SVGs (icons) are inline, `aria-hidden="true"` and
`focusable="false"` (already the button contract: `button.md:9`); they must never carry `alt`.
Images that communicate content require a meaningful `alt`.

**Implementation**: the icon slot is `template.HTML` internal-only (`IconSVG`, `button.go`), marked
`aria-hidden` + `focusable="false"`. The `og:image` (an asset, not an `<img>`) is metadata, not
page content. Any future content `<img>` in Markdown must supply `alt` via Markdown syntax.

**Example**:

```html
<!-- decorative icon -->
<svg aria-hidden="true" focusable="false" ...><path .../></svg>
```

---

## 10. Pagination — real links + `aria-current` (data table)

**Architecture rule**: pagination is a set of real GET links preserving the current query (`?q=`,
`?sort=`, `?dir=`, `?page=`), with the current page exposed as non-link text carrying
`aria-current="page"`. Page state is in the URL, never in JS memory.

**Implementation**: exists in the data table — `dataTableHref` (`data_table.go:333-355`) builds
stable links; the template renders current page as `<span aria-current="page">` and others as
`<a href>` with `hx-get` enhancement (`data-table.html:74-76`). The `noindex` rule (§4) keeps
pagination-descendant URLs out of the index so infinite indexed duplicates don't occur.

**Example**:

```html
<nav class="ui-data-table-pagination" aria-label="Table pages">
  <a class="ui-data-table-page" href="?page=1">1</a>
  <span class="ui-data-table-page ui-data-table-page--current" aria-current="page">2</span>
  <a class="ui-data-table-page" href="?page=3">3</a>
</nav>
```

---

## 11. Breadcrumbs — Phase F pattern + `BreadcrumbList` contract

**Architecture rule**: every nested page (all `/components/*`, future `/docs/*` subpages) exposes a
visible breadcrumb as a real `<nav>` + `<ol>` of links, and a matching `BreadcrumbList` JSON-LD
block. Breadcrumb is a **public content pattern** (Phase F), never a theme concern.

**Implementation**: a `breadcrumb` template consumes `[]navLink`-like items ending at the current
page (rendered as `<span aria-current="page">`); the handler builds the trail from the route
registry (Home → Docs → Component). The JSON-LD for the page is emitted by the resolver (§12).

**Example**:

```html
<nav aria-label="Breadcrumb"><ol>
  <li><a href="/">Home</a></li>
  <li><a href="/docs">Docs</a></li>
  <li><span aria-current="page">Button</span></li>
</ol></nav>
```

---

## 12. Structured data JSON-LD — WebSite, BreadcrumbList, Article, SoftwareApplication

**Architecture rule**: JSON-LD is emitted server-side per page type, declarative, zero-JS, valid
JSON (must parse with `encoding/json` in tests). Types:

| Page | JSON-LD |
|---|---|
| `/` | `WebSite` (+ `Organization`/`SoftwareSourceCode` for the repo) |
| `/docs` | `CollectionPage`/`WebPage` |
| `/components/*` | `@graph` of `BreadcrumbList` + `TechArticle` |
| demos | none (noindex) |

**Implementation**: `pageView.Meta.JSONLD template.JS` is populated by the resolver and emitted by
the layout before `</head>`. The home page ships the `WebSite` block; every registered component
page ships a single `@graph` document with the `BreadcrumbList` trail (Home > Components > page,
matching the visible breadcrumb) plus a `TechArticle` carrying the page headline and canonical URL
(`componentJSONLD`, `server.go:214-260`). Values are derived from the registry and fixed system
facts (version `0.4.0` from `package.json:3`, license MIT). No JSON is built by string
concatenation in templates; Go `encoding/json` marshals typed structs so escaping and validity are
guaranteed (tested with `json.Valid` in `server_test.go`).

**Example**:

```json
{
  "@context": "https://schema.org",
  "@type": "WebSite",
  "name": "Gelium UI",
  "url": "https://gelium-ui.example/",
  "inLanguage": "en",
  "publisher": { "@type": "Organization", "name": "Gelium UI" }
}
```

```json
{
  "@context": "https://schema.org",
  "@type": "BreadcrumbList",
  "itemListElement": [
    {"@type": "ListItem", "position": 1, "name": "Home", "item": "https://gelium-ui.example/"},
    {"@type": "ListItem", "position": 2, "name": "Docs", "item": "https://gelium-ui.example/docs"},
    {"@type": "ListItem", "position": 3, "name": "Button", "item": "https://gelium-ui.example/components/button"}
  ]
}
```

---

## 13. Performance — server-rendered, no CDN, minified CSS, semantic HTML

**Architecture rule**: HTML arrives from the server already complete (crawlable, indexable,
renderable without JS). Assets are embedded, CSS is minified (`package.json` build lane), scripts
are `defer` and optional (`layout.html:8-9`). The SEO layer must not add weight: metadata and
JSON-LD are a few hundred bytes.

**Implementation**: keep the head order CSS → deferred JS → metadata (`layout.html:3-10`). Budget:
metadata + JSON-LD total ≤ 1 KB per page. The audit's remaining performance items (gzip, immutable
cache for `/static/*` with the existing `?v=` cache-buster, `server.go:228-249`) are delivery
concerns tracked with this contract but not part of the markup contract.

---

## 14. Mobile rendering — responsive, no mandatory JS

**Architecture rule**: every page is fully usable and readable at mobile widths with JS disabled;
viewport is present (`layout.html:5`). Metadata must not assume desktop rendering (no viewport-gated
title swaps, no JS-injected descriptions).

**Implementation**: goldmark output + existing responsive layout handle content; HTMX only enhances.
Nothing in the SEO layer changes per viewport.

---

## 15. Indexability — robots + canonical + no duplicates

**Architecture rule**: a page is indexable iff `robots` is `index` AND it has a canonical. Every
indexable page has exactly one URL identity. This is the gate for the sitemap: only indexable pages
appear there.

**Implementation**: `metaFor` computes `Robots` and `Canonical` together; `sitemap.xml` filters on
`robots == index`. A test asserts: `index, follow` ⇒ canonical present, and noindex pages are absent
from the sitemap.

---

## 16. Duplicate content — canonical; no params in canonical URLs

**Architecture rule**: the canonical URL is the clean path with **no query parameters**. Stateful
variants (`?q=`, `?sort=`, `?page=`, `data_table.go`) are the same document; they must not generate
canonical URLs or sitemap entries. This keeps the index free of duplicate permutations.

**Implementation**: `Canonical = BASE_URL + r.URL.Path` (path only, query dropped). `metaFor`
never reads `r.URL.RawQuery` when building the canonical.

**Example**:

```
Request:  /components/data-table?sort=name&page=2
Canonical: https://gelium-ui.example/components/data-table
```

---

## 17. hreflang — future multilingual

**Architecture rule**: when a second language ships, every localized page emits a self-referencing
`hreflang` block with all language alternatives, and `Lang` becomes per-route data instead of the
fixed `en` (`layout.html:2`).

**Implementation (future)**: `pageMeta` gains `HrefLang []struct{ Lang, URL string }`; the layout
emits one `<link rel="alternate" hreflang="X" href="...">` per entry plus
`<link rel="alternate" hreflang="x-default">`. The contract does not change markup today; it only
forbids shipping markup that makes a future `hreflang` impossible (e.g. hardcoding a wrong `lang` or
building canonical URLs from the request host).

**Example**:

```html
<link rel="alternate" hreflang="en" href="https://gelium-ui.example/components/button">
<link rel="alternate" hreflang="es" href="https://gelium-ui.example/es/components/button">
<link rel="alternate" hreflang="x-default" href="https://gelium-ui.example/components/button">
```

---

## 18. Verification

Phase E DoD (`gelium-ui-system-roadmap.md:189`) requires metadata server-driven and green build/test.
The SEO contract is verified with:

```bash
go test ./...   # asserts: one h1, title per route, description 150-160, canonical present,
                # robots default/noindex, JSON-LD parses, sitemap == registry, one canonical per URL
go vet ./...
```

Acceptance: `curl -s localhost:<port>/components/button | grep -E 'canonical|description|ld\+json'`
returns all three.
