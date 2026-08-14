# Exploration: docs-shell

Scalar-style docs shell for **Gelium UI** (module `loomui`, path `loom-ui`), dogfooding existing Gelium primitives. Landing marketing page is **out of scope**; this change is **docs chrome first**.

## Executive answer

Today “docs shell” is a **flat header + centered prose column**, not a sidebar/topbar product docs frame. Almost all Scalar-like chrome can be **composed** from existing Gelium pieces (Navigation drawer, List, Breadcrumb, Theme switcher, Text field, Footer). First shippable slice: **wrap `/docs` + `/components/*` in a two-pane shell, keep URLs, no search/TOC yet**.

---

## Current State

### Layout chrome today

| Surface | What renders |
|---------|----------------|
| `web/templates/layout.html` | Skip link → `<header class="site-header">` (brand, flat primary nav, theme switcher) → breadcrumb → banner → `<main class="docs-shell">` (prose + optional demos) → footer → toast region |
| `web/styles/base.css` | `.site-header` flex row; `.docs-shell` = **centered content column only** (`width: min(68rem…); margin: 4rem auto`) — **not** a sidebar layout |
| `pageView` (`internal/app/server.go`) | `Nav []navLink`, `ThemeSwitcher`, `Breadcrumb`, `Footer`, `Meta`, content + many demo slots |
| Render choke point | `renderMarkdownStatus`: always sets `Nav = navLinks()`, default footer/breadcrumb, `ThemeSwitcher = themeSwitcherFor`, theme from `?theme=` |

**Primary nav today** (`navLinks()` in `routes.go`):

1. `/docs` → “Docs”
2. Flat list of **28** `componentRoutes()` labels (Button … Tooltip)

No Recipes, Themes, Patterns, or Getting started groups in chrome. Roadmap residual already names this: *“Nav discoverability — Recipes/demos/themes not in navLinks()”* (`docs/gelium-ui-system-roadmap.md`).

### IA / routes today vs target

| Today | Target shell IA (product) |
|-------|---------------------------|
| `/` home (minimal landing markdown + CTA) | Landing **later** (out of scope) |
| `/docs` generated index from `docsSections` | Docs hub / Getting started entry |
| `/components/*` (28 dogfood pages) | Components group |
| `/recipes/*` **standalone** full HTML (own appbar, not `layout`) | Recipes group (link into shell later) |
| `/demo/*` noindex | Optional / secondary |
| Themes via `?theme=` switcher only | Themes section (docs pages later) |
| Public/state patterns: partials, **no** `/components` docs for many | Patterns (policy or pages later) |

**Important:** Grouped IA already exists as data for the **index only**:

- `docsSections` in `internal/app/docs.go` — Foundation / Actions / Input / Feedback / Navigation / Data + hard-coded Demos + Recipes markdown blocks
- Chrome does **not** consume `docsSections`; header dumps a flat 29-link nav

JSON-LD and breadcrumbs already treat `/docs` as the “Components” hub (`componentBreadcrumb`, `componentJSONLD`).

### Theme switcher (already done)

- Catalog: `availableThemes` (Material, Basecoat)
- Middleware: `themeQueryMiddleware` + `?theme=<slug>`
- Partial: `theme-switcher.html` (plain GET links, 0 JS)
- Injected on every layout page via `themeSwitcherFor`
- **Plug into new topbar:** move the same `{{template "theme-switcher"}}` slot from `site-header` into the shell topbar; no new mechanism

### Search / TOC / prev-next

| Capability | Today |
|------------|--------|
| Site docs search | **None**. Recipe/demo local `?q=` only (admin-resource, data-table, whatsapp) |
| TOC (“on this page”) | **None**. `goldmark.New()` default — no auto heading IDs, no TOC extension |
| Prev/next docs | **None**. `pagination.html` is table/list pagination, not doc trail |
| Version in chrome | Asset query `?v=0.4.0` only; no product version control |

### Recipes relationship

Recipes (`recipe-admin-resource.html`, etc.) use **standalone** document shells with recipe-specific appbars and a “Docs” back link. First docs-shell slice should **not** force recipes into `layout.html` (large churn, separate UX). Link them from sidebar only.

---

## Affected Areas

| Path | Why |
|------|-----|
| `web/templates/layout.html` | Restructure chrome: sidebar + topbar + content region |
| `web/styles/base.css` (and/or new `docs-shell.css`) | Replace centered-only `.docs-shell` with sticky two-pane frame |
| `internal/app/routes.go` | Evolve `navLinks` / introduce grouped docs nav model |
| `internal/app/docs.go` | Single source for nav groups (`docsSections` or successor) |
| `internal/app/server.go` | `pageView` fields (Sidebar/DocsNav/ActivePath), render defaults, maybe layout mode |
| `web/templates/navigation-drawer.html` (+ CSS) | Reuse standard + modal variants in real chrome (not only demo) |
| `web/templates/theme-switcher.html` | Relocate into topbar (markup host change) |
| `web/templates/breadcrumb.html` | Stay; likely move under topbar/content header |
| `internal/app/server_test.go` (+ chrome-related tests) | Header/nav/footer/sitemap contracts will need updates |
| `docs/gelium-ui-system-roadmap.md` | Residual “Nav discoverability” closes partially |
| Registries (later apply) | Pattern/composition note if docs-shell is a composition |

**Low touch / out of first slice:** recipe templates, demo whatsapp, theme CSS files, wire contracts, embed API signature (`New() http.Handler` stays).

---

## Gelium primitives inventory (dogfood)

### Reuse as-is (or with thin composition CSS)

| Primitive | Role in Scalar-like shell |
|-----------|---------------------------|
| **Navigation drawer** standard | Desktop sticky sidebar destinations |
| **Navigation drawer** modal (`<dialog>`) | Mobile nav overlay (existing invoker pattern) |
| **List** (nav links) | Alternative denser sidebar if drawer glyphs are heavy |
| **Breadcrumb** | Content header trail (already wired) |
| **Theme switcher** | Topbar visual direction |
| **Text field** | Future search field (type=search pattern exists in demos; text-field is `type="text"` today — may need `type` param or native search input in composition) |
| **Footer** | Site chrome bottom (already defaulted) |
| **Section heading** | In-page section titles; optional group labels |
| **Divider** | Sidebar section separators |
| **Icon button** | Mobile “open nav” trigger if not using Button+command |
| **Skip link** | Keep; retarget still `#main-content` |
| **Pagination** pattern | Inspiration for prev/next links (not drop-in) |

### Gaps that force new work (flag for design/spec)

| Gap | Force new component? | Prefer |
|-----|----------------------|--------|
| **Grouped nav sections** (Getting started / Components / …) | Drawer has **no** section header/subheader API — only flat destinations | **Composition**: `h2`/section-heading or plain label + drawer list per group; small drawer extension only if reused elsewhere |
| **Docs frame layout** (sticky sidebar + topbar + main) | Not a component today; `.docs-shell` is content width | **Composition CSS** in docs chrome layer (tokens only), not a public “Layout” widget unless we want library consumers to copy it |
| **TOC** | No primitive | **Defer** first slice; later goldmark heading IDs + simple `<nav aria-label="On this page">` list |
| **Site search** | No index/search route | **Defer**; later GET `/docs/search?q=` server filter over registry |
| **Version switcher** | None | **Defer** (single version product) |
| **Docs prev/next** | None | **Defer** or tiny link row from ordered registry |
| **Text field `type=search`** | Input is hard-coded `type="text"` | Optional small extension when search ships |

**Dogfood rule:** shell MUST use `.ui-*` primitives and tokens. Ad-hoc `site-header nav a` styles today are already docs-chrome CSS — first slice should **replace** flat header nav with drawer/list primitives rather than invent a third nav skin.

---

## Approaches

### URL strategy

| Option | Description | Pros | Cons | Effort |
|--------|-------------|------|------|--------|
| **A. Keep paths + shell wrap** | Keep `/components/*`, `/docs`; new layout around them | Zero SEO move; sitemap/canonical/JSON-LD stable; least test path churn | Paths less “docs-prefixed” than Scalar | **Low** |
| **B. Move under `/docs/*`** | e.g. `/docs/components/button` + redirects | Cleaner docs tree; matches many design systems | Redirect matrix; canonical/sitemap/JSON-LD/tests/content links all move; embedders bookmarking old paths | **High** |
| **C. Hybrid** | Keep `/components/*`; add `/docs/getting-started`, `/docs/themes` later; optional aliases | Incremental IA without breaking components | Two URL styles forever unless cleaned later | **Medium** |

**Recommendation: A for first slice**, with **C-shaped growth** (new non-component docs pages can live under `/docs/...` without moving components).

### Shell structure options

1. **Compose drawer + layout CSS (recommended)**  
   - Desktop: standard `ui-navigation-drawer` sticky in grid  
   - Mobile: modal drawer + menu trigger (existing)  
   - Topbar: brand + (optional search placeholder) + theme switcher  
   - Content: breadcrumb + article.prose + demos  
   - Pros: max dogfood, reuses tested a11y patterns  
   - Cons: drawer is 360px Material-ish; may need density token tweak for long component lists  
   - Effort: Medium

2. **List-based sidebar only**  
   - Pros: denser, simpler  
   - Cons: less “product” dogfood of Navigation drawer; still need mobile pattern  
   - Effort: Low–Medium

3. **New docs-only chrome components**  
   - Pros: perfect Scalar fit  
   - Cons: violates “compose first”; review budget; dual nav systems  
   - Effort: High — **reject for v1**

### Mobile (no-JS first)

| Pattern | Fit |
|---------|-----|
| **Modal Navigation drawer** (`<dialog>` + button `command`/`commandfor`) | Already implemented; progressive where Invoker Commands exist |
| **No-JS fallback** | Ensure trigger is a real control: page link to `/docs` index **or** `<details>` nav — modal without invokers is a known overlay risk (G1 history). Prefer **details/summary sidebar** as bulletproof no-JS mobile nav, with modal as enhancement **or** always-visible collapsed details |
| **HTMX** | Not required for nav; full page loads are correct |

**Recommendation:** Desktop permanent drawer; mobile `<details class="…">` labeled “Menu” wrapping the same nav model (0 JS), optional modal later if details UX is insufficient.

### Search / TOC deferral

- **Search:** defer. Placeholder non-submitting field is a product lie — better omit than fake.
- **TOC:** defer until goldmark heading IDs (or manual anchors in markdown) exist; otherwise TOC links are dead.

---

## Recommendation

### Product direction

Ship a **docs chrome composition** that looks Scalar-like:

```text
┌ sidebar (sticky) ┬ topbar (brand · theme) ────────────────┐
│ grouped nav      │ breadcrumb                             │
│  Getting started │ # Title + dogfood preview              │
│  Components *    │ prose                                  │
│  Recipes         │                                        │
│  Themes (links)  │                                        │
└──────────────────┴────────────────────────────────────────┘
```

\* Components expanded via existing `docsSections` categories.

### Technical approach

1. **URL option A** — keep `/docs`, `/components/*`.
2. **Single nav model** derived from `docsSections` + Recipes + Themes entries (fix `navLinks()` flat list as primary chrome).
3. **layout.html** becomes two-pane docs frame; keep one layout for home+docs (home can hide sidebar later or show minimal nav — decide in design; exploration bias: **same shell for `/docs` and `/components/*`**, home may keep simpler header until landing work).
4. **Dogfood:** standard Navigation drawer (or List+section labels) + existing theme switcher + breadcrumb + footer.
5. **No** search, TOC, version switcher, route moves, recipe layout merge, marketing landing.

### Recommended first slice (proposal scope)

**Smallest shippable docs shell**

| In | Out |
|----|-----|
| Grouped sidebar nav data model (from `docsSections` + Recipes + link to theme via switcher or short Themes note) | Route migration / redirects |
| `layout.html` + CSS: sticky sidebar + topbar + main | Site search |
| Active state from `r.URL.Path` | TOC / heading IDs |
| Theme switcher in topbar | Prev/next |
| Mobile: details menu or modal drawer with no-JS path to `/docs` | New themes / recipes content |
| Tests updated for chrome contracts | Registry JSON runtime |
| Optional: slim header on home only if shell hurts landing | Marketing landing |

**Forecast vs 800-line review budget:** one focused PR is plausible if CSS+templates+nav model+tests stay tight; if drawer extensions + home special-case + full test rewrite land together, **chain**: (1) nav model + sidebar markup no visual polish, (2) layout CSS + mobile, (3) test/docs residual. Delivery strategy session default: **auto-forecast**.

---

## Risks

| Risk | Mitigation |
|------|------------|
| **Test churn** — many tests assert layout fragments (`docs-shell`, footer from `navLinks`, header structure) | Update chrome contract tests in same slice; keep SEO/canonical assertions stable (option A) |
| **Shared `layout.html` / `server.go`** — high-contention files | Small PR; avoid recipe merge |
| **Drawer without groups** looks like 40 flat icons | Use section labels + optional no-glyph density variant via composition |
| **Modal drawer no-JS** | Prefer details/summary or `/docs` fallback; don’t rely only on Invoker Commands |
| **SEO** if someone chooses option B later | Stick to A now; document aliases only if needed |
| **Embed API** | `New()` unchanged; chrome is internal to docs handler |
| **800-line budget** | Scope first slice ruthlessly; chain if over |
| **Home vs docs** | Explicit decision: shell on docs routes only vs all layout pages |
| **Footer still lists flat components** | Rebuild footer sections from same grouped model to avoid drift |

---

## Ready for Proposal

**Yes.**

Orchestrator should tell the user:

- Exploration complete for **docs-shell**.
- Recommend **URL keep + shell wrap**, dogfood **Navigation drawer / grouped nav + existing theme switcher**, defer search/TOC/landing.
- Next phase: **sdd-propose** for change `docs-shell` with the first-slice table above.
- Open question for proposal (at most one if interactive): **Apply shell to home (`/`) in v1, or only `/docs` + `/components/*`?** Exploration bias: docs+components only; home waits for landing work.

---

## Evidence index (symbols / files)

- `navLinks`, `componentRoutes` — `internal/app/routes.go`
- `docsSections`, `docsIndex` — `internal/app/docs.go`
- `pageView`, `themeSwitcherFor`, `defaultFooter`, `componentBreadcrumb`, `renderMarkdownStatus`, `New` — `internal/app/server.go`
- Layout chrome — `web/templates/layout.html`, `web/styles/base.css` (`.site-header`, `.docs-shell`)
- Drawer API — `web/templates/navigation-drawer.html`, `internal/app/navigation_drawer.go`
- Theme partial — `web/templates/theme-switcher.html`
- Residual DX — `docs/gelium-ui-system-roadmap.md` (Nav discoverability)
- Recipes outside layout — `web/templates/recipe-*.html`, `renderRecipeTemplate`
