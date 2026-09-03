# Gelium UI — SEO Patterns

> Pattern reference for the Gelium UI documentation site (physical repo `loom-ui`).
> Phase E of `docs/gelium-ui-system-roadmap.md`. Companion to `docs/gelium-ui-seo-contract.md`
> (the contract, "what must hold") — this document is the "how to compose" cookbook.
>
> Each pattern states: **when to use it**, a **reference template/markup**, and its **relationship
> to Gelium components**. All patterns are server-driven, zero-JS, and theme-agnostic.

---

## P0 — Server-driven metadata head (the core pattern)

**When to use**: every GET page that renders the layout. This is the composition root: all other
patterns feed values into it.

**Structure**: the handler populates `pageView.Meta` via the resolver `metaFor(route, req)`; the
layout emits the head block behind `{{if}}` guards so a page with partial metadata still renders a
valid head.

**Reference layout head** (`web/templates/layout.html`):

```html
<!doctype html>
<html lang="{{or .Meta.Lang "en"}}" class="{{.ThemeClass}}">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{.Title}} · Gelium UI</title>
  {{if .Meta.Description}}<meta name="description" content="{{.Meta.Description}}">{{end}}
  {{if .Meta.Canonical}}<link rel="canonical" href="{{.Meta.Canonical}}">{{end}}
  {{if .Meta.Robots}}<meta name="robots" content="{{.Meta.Robots}}">{{end}}
  <link rel="stylesheet" href="/static/app.css?v={{.AssetsVersion}}">
  <script defer src="/static/htmx.min.js?v={{.AssetsVersion}}"></script>
  <script defer src="/static/app.js?v={{.AssetsVersion}}"></script>
  {{if .Meta.OGTitle}}<meta property="og:type" content="{{or .Meta.OGType "website"}}">
  <meta property="og:title" content="{{.Meta.OGTitle}}">
  <meta property="og:description" content="{{.Meta.OGDescription}}">
  <meta property="og:url" content="{{.Meta.OGURL}}">
  {{if .Meta.OGImage}}<meta property="og:image" content="{{.Meta.OGImage}}">{{end}}
  <meta name="twitter:card" content="{{or .Meta.TwitterCard "summary"}}">{{end}}
  {{if .Meta.JSONLD}}<script type="application/ld+json">{{.Meta.JSONLD}}</script>{{end}}
</head>
```

**Relationship**: consumes `pageMeta` from the contract (§1); `ThemeClass` stays presentational only.

---

## P1 — Breadcrumb pattern (Phase F)

**When to use**: any page deeper than the root — all `/components/*`, and future `/docs/*`
subpages. Home and `/docs` itself do not need one.

**Reference markup** (Phase F `breadcrumb` component):

```html
<nav aria-label="Breadcrumb">
  <ol>
    {{range .Crumbs}}{{if .Current}}
    <li><span aria-current="page">{{.Label}}</span></li>{{else}}
    <li><a href="{{.Path}}">{{.Label}}</a></li>{{end}}{{end}}
  </ol>
</nav>
```

**Relationship**: the crumb trail is built from `componentRoutes()` (`routes.go:16-47`) — Home →
Docs → Component — so it can never drift from the nav. The same data feeds the `BreadcrumbList`
JSON-LD (P4), keeping visual and structured hierarchy in lockstep. It is a public content pattern,
not a theme component (`gelium-ui-system-roadmap.md:197-200`).

---

## P2 — Heading pattern (single h1)

**When to use**: every page, always. The layout contributes zero headings; the content owns exactly
one `h1` (Markdown `# `), then `h2` sections, then `h3`.

**Reference** (`web/content/button.md:1-5` renders as):

```html
<article class="prose">
  <h1>Button</h1>
  <p>Button is an open-code component...</p>
  <h2>Variants and states</h2>
</article>
```

**Relationship**: the `/docs` index generates its own `# Documentation` (`docs.go:85`); the
`section-heading` public pattern (Phase F) must render `h2`-level headings, never `h1`. Test
assertion: one `h1` per rendered page.

---

## P3 — Link pattern (descriptive anchor text)

**When to use**: every internal link in nav, breadcrumbs, and Markdown content.

**Rules**: anchor text is the destination's label or a descriptive phrase ("Icon button", "Data
table") — never "click here" / "more". Nav labels come from `navLinks()` (`routes.go:49-58`);
content links are written in Markdown and rendered by goldmark.

**Reference**:

```md
Compare with [Icon button](/components/icon-button) for single-icon actions.
```

```html
<p>Compare with <a href="/components/icon-button">Icon button</a> for single-icon actions.</p>
```

**Relationship**: the same label pool drives nav, breadcrumbs and docs index (`docs.go:87-93`), so
anchor text is consistent site-wide. HTMX-enhanced links (`hx-get`) always keep the real `href`.

---

## P4 — JSON-LD snippets per page type

**When to use**: every indexable page, one block per page type (see contract §12).

### WebSite — home

```go
type webSiteLD struct {
    Context  string        `json:"@context"`
    Type     string        `json:"@type"`
    Name     string        `json:"name"`
    URL      string        `json:"url"`
    InLanguage string      `json:"inLanguage"`
    Publisher organizationLD `json:"publisher"`
}
```

```json
{"@context":"https://schema.org","@type":"WebSite","name":"Gelium UI",
 "url":"https://gelium-ui.example/","inLanguage":"en",
 "publisher":{"@type":"Organization","name":"Gelium UI"}}
```

### BreadcrumbList — `/components/*`

```json
{"@context":"https://schema.org","@type":"BreadcrumbList","itemListElement":[
 {"@type":"ListItem","position":1,"name":"Home","item":"https://gelium-ui.example/"},
 {"@type":"ListItem","position":2,"name":"Docs","item":"https://gelium-ui.example/docs"},
 {"@type":"ListItem","position":3,"name":"Button","item":"https://gelium-ui.example/components/button"}]}
```

### SoftwareApplication — `/components/*` (the system itself)

```json
{"@context":"https://schema.org","@type":"SoftwareApplication",
 "name":"Gelium UI","applicationCategory":"DeveloperApplication",
 "softwareVersion":"0.6.1","operatingSystem":"Any (web)","license":"MIT",
 "url":"https://gelium-ui.example/"}
```

### Article / TechArticle — component content (optional wrapper)

```json
{"@context":"https://schema.org","@type":"TechArticle","headline":"Button · Gelium UI",
 "about":"Button component","inLanguage":"en","isPartOf":{"@type":"WebPage","url":"..."}}
```

**Relationship**: built from a typed Go struct and marshaled with `encoding/json` — never string
concatenation — so output always parses. Version and license come from system facts (`package.json`,
`README.md`), not from page content.

---

## P5 — Sitemap generation pattern

**When to use**: a single `GET /sitemap.xml`, generated on each request from the registry.

**Reference**:

```go
func (s *server) sitemap(w http.ResponseWriter, _ *http.Request) {
    var sb strings.Builder
    sb.WriteString(xmlHeader)
    for _, p := range s.indexablePaths() { // "/", "/docs" + componentRoutes() where Robots == index
        fmt.Fprintf(&sb, `<url><loc>%s%s</loc></url>`, baseURL, p)
    }
    w.Header().Set("Content-Type", "application/xml")
    _, _ = w.Write([]byte(sb.String()))
}
```

**Relationship**: `indexablePaths()` is derived from `componentRoutes()` + robots policy — a page
that is `noindex` never appears. Referenced from `robots.txt` (P6).

---

## P6 — Robots pattern

**When to use**: `meta name="robots"` on every page; `GET /robots.txt` on the server.

**Reference**:

```html
<!-- default -->
<meta name="robots" content="index, follow">
<!-- demo / stateful shell -->
<meta name="robots" content="noindex, nofollow">
```

```text
# /robots.txt (handler-served)
User-agent: *
Allow: /
Sitemap: https://gelium-ui.example/sitemap.xml
```

**Relationship**: the resolver computes `Robots` from the route type (demos at `/demo/*`,
`server.go:136-141`, and POST pages at `postOnlyPaths()`, `server.go:164-176`, are noindex). Robots,
canonical and sitemap are computed in one pass so indexability is always consistent.

---

## P7 — Social preview pattern (Open Graph + Twitter)

**When to use**: every indexable page. Values are derived from canonical metadata (`Title`,
`Description`, `Canonical`), never written independently.

**Reference**:

```html
<meta property="og:type" content="website">
<meta property="og:title" content="Documentation · Gelium UI">
<meta property="og:description" content="The Gelium UI component library, dogfooded on every page.">
<meta property="og:url" content="https://gelium-ui.example/docs">
<meta property="og:image" content="https://gelium-ui.example/static/og-cover.png">
<meta name="twitter:card" content="summary_large_image">
```

**Relationship**: a single static `og-cover.png` asset serves all pages (no per-page image
pipeline); `og:type` flips `website` ↔ `article` based on the page type from the registry.

---

## P8 — Pagination pattern (data table, real links)

**When to use**: any server-side list with more rows than a page (the data table demo,
`data_table.go`). Pagination is content navigation, so it must stay index-safe (canonical ignores
query params) and crawlable (real `href`).

**Reference** (already shipped, `web/templates/data-table.html:74-76`):

```html
<nav class="ui-data-table-pagination" aria-label="Table pages">
  {{range .PageLinks}}{{if .Current}}
  <span class="ui-data-table-page ui-data-table-page--current" aria-current="page">{{.Num}}</span>{{else}}
  <a class="ui-data-table-page" href="{{.Href}}" hx-get="{{.Href}}"
     hx-target="#data-table-panel" hx-swap="outerHTML">{{.Num}}</a>{{end}}{{end}}
</nav>
```

**Relationship**: links come from `dataTableHref` (`data_table.go:333-355`) preserving `?q=&sort=&dir=&page=`
stably; the canonical (contract §16) drops all query params so paginated variants never index as
duplicates.

---

## Pattern selection summary

| When you need... | Use |
|---|---|
| The head of any page | P0 (server-driven metadata head) |
| Context/hierarchy on nested pages | P1 (breadcrumb) + P4 (BreadcrumbList) |
| Page structure | P2 (single h1) |
| Navigation and citations | P3 (descriptive links) |
| Rich result eligibility | P4 (JSON-LD) |
| Whole-site discovery | P5 (sitemap) + P6 (robots) |
| Sharing on social networks | P7 (social preview) |
| Long server lists | P8 (pagination) |
