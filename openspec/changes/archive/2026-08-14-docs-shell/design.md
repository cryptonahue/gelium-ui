# Design: docs-shell

## Technical Approach

Wrap `/docs`, `/docs/*` stubs, and `/components/*` in a Scalar-style two-pane chrome composed from Gelium primitives. Keep `/` on the existing `site-header` layout. Single nav builder feeds sidebar + footer sections. No URL moves, no recipe layout merge, no real search/TOC.

Maps proposal approach + spec requirements 1–10.

## Architecture Decisions

| Decision | Options | Choice | Why |
|----------|---------|--------|-----|
| Shell enablement | Always-on layout vs flag | `pageView.DocsNav *docsNavView`; non-nil → shell | `/` must stay unchanged; one `layout.html` |
| Nav source | Extend `navLinks` vs new builder | `docsNavFor(path)` in `docs.go` from `docsSections` + static IA groups | Grouped IA already in `docsSections`; chrome must consume it |
| Sidebar primitive | Full drawer destinations vs list | **List + section labels + divider** (`ui-list`, `ui-list-item-link`, `ui-section-heading` / group label, `ui-divider`) | Drawer requires glyphs; 28 components stay dense without inventing icons |
| Mobile | Modal drawer vs details | **`<details>/<summary>` only** | Spec 0-JS baseline; invokers not required |
| Mobile DOM | One tree vs dual | **Dual render** (desktop `<aside>` + mobile `<details>`) same model | CSS cannot safely share one landmark across sticky aside + disclosure |
| Theme switcher | Duplicate vs relocate | Same `themeSwitcherFor` + `theme-switcher` partial in topbar when shell | No new mechanism; home keeps header slot |
| Patterns/Themes | Deep-link only vs stubs | Thin **`GET /docs/patterns`**, **`GET /docs/themes`** via `renderMarkdown` | Real hrefs, shell stays on; content can grow later |
| Recipes | Shell pages vs outbound | Sidebar links to existing `/recipes/*` only | Spec allows leaving shell; avoids recipe template churn |
| Search/version | Omit vs honest slots | Disabled search control + static version badge (`0.4.0` align asset query) | Spec: placeholders OK if not fake-broken |
| Footer | Keep `navLinks` slice vs rebuild | Rebuild `defaultFooter` sections from `docsNavFor` flat export | One source; avoid drift |

## Data Flow

```
Request path P
    │
    ▼
renderMarkdownStatus / renderErrorPage (layout paths)
    │
    ├─ usesDocsShell(P)? ──yes──► DocsNav = docsNavFor(P)
    │                              ThemeSwitcher = themeSwitcherFor(r, …)
    │                              Breadcrumb default as today
    │
    └─ no (e.g. /) ──────────────► DocsNav = nil; site-header + Nav optional/minimal
    │
    ▼
layout.html
    DocsNav != nil → docs-topbar + mobile details + desktop aside + main
    DocsNav == nil → legacy site-header + main.docs-shell-content (centered)
```

### View model (non-obvious)

```go
type docsNavLink struct {
    Path, Label string
    Current     bool // path == active (exact)
}
type docsNavGroup struct {
    Title string
    Links []docsNavLink
}
type docsNavView struct {
    Groups []docsNavGroup // Getting started, Components/*, Patterns, Recipes, Themes
    // Topbar slots
    Version string // "0.4.0"
    SearchDisabled bool // always true this change
}
```

**IA construction (`docsNavFor`):**

1. **Getting started** — `/docs` (“Documentation”)
2. **Components** — one group per `docsSections` entry (Foundation…Data), links marked Current when `Path == activePath`
3. **Patterns** — `/docs/patterns`
4. **Recipes** — `/recipes/admin-resource`, `/recipes/ops-queue`, `/recipes/public-feed` (outbound)
5. **Themes** — `/docs/themes` (copy explains `?theme=` + switcher; does not replace switcher)

`navLinks()` remains a flat helper for any legacy caller or is reduced to “all component links” derived from `docsSections` so registration order and index cannot drift. Prefer deriving component link lists from `docsSections` only (not re-flattening `componentRoutes` for chrome).

`usesDocsShell(path)`: `path == "/docs" || strings.HasPrefix(path, "/docs/") || strings.HasPrefix(path, "/components/")`.

## Template structure

| Piece | Role |
|-------|------|
| `layout.html` | Branch on `.DocsNav`; keep skip link first; footer + toast after shell |
| `docs-topbar.html` | Brand → `/`, search slot (disabled), version badge, `{{template "theme-switcher"}}` |
| `docs-sidebar.html` | Renders one nav tree: groups → heading + `ul.ui-list` links; `aria-current="page"` on current |
| Mobile block in layout | `<details class="docs-nav-mobile"><summary>Menu</summary>{{template "docs-sidebar" …}}</details>` |
| Desktop | `<aside class="docs-sidebar-desktop"><nav aria-label="Docs">…docs-sidebar…</nav></aside>` |
| Main | `#main-content`; breadcrumb + banner above article inside content column |

**Landmarks:** skip → `#main-content`; docs nav `nav`; content `main`. Theme switcher keeps `aria-label="Visual direction"`.

**Not new public components** — composition CSS under docs chrome layer only.

## CSS

Evolve `base.css` (and small `docs-shell.css` if file size warrants; import from `app.css`):

| Class | Behavior |
|-------|----------|
| `.docs-shell` | **Repurpose** to two-pane frame: CSS grid `sidebar \| main` desktop; single column mobile. Sticky topbar full width. |
| `.docs-shell-main` / content column | Former centered column width (`min(68rem…)`) moves **here**, not full-page margin auto on shell root |
| `.docs-topbar` | Sticky `top:0`; flex; surface + shadow tokens |
| `.docs-sidebar-desktop` | Sticky under topbar; `overflow:auto`; width ~16–18rem; hide `< md` |
| `.docs-nav-mobile` | Show `< md` only; summary uses label type tokens |
| Home / non-shell | Keep `.site-header` + centered content utility (rename old centered rule to `.docs-content` if shell consumes `.docs-shell`) |

**Breakpoint:** ~48rem (`md`). Tokens only (`--ui-color-*`, type, radius, shadow). Theme-agnostic markup.

**Sticky:** topbar always; desktop sidebar `position: sticky; top: <topbar-height>; max-height: calc(100vh - topbar)`.

## Gelium composition (exact)

| Slot | Primitive / class |
|------|-------------------|
| Sidebar links | `ui-list`, `ui-list-item`, `ui-list-item-link` (+ current modifier class if needed, e.g. `aria-current` styled in docs CSS) |
| Group titles | `ui-section-heading` **or** plain `p.docs-nav-group-label` using `--ui-type-label-*` if heading level noise is high |
| Separators | `ui-divider` between top-level IA blocks |
| Theme | `theme-switcher` → `.ui-theme-switcher*` |
| Trail | `breadcrumb` → `.ui-breadcrumb*` |
| Footer | `footer` → `.ui-footer*` |
| Skip | `.ui-skip-link` unchanged |

Drawer (`ui-navigation-drawer*`) **not** required in v1 chrome (density); remains dogfood on its docs page.

## Stubs & routes

| Route | Handler | Shell |
|-------|---------|-------|
| `GET /docs/patterns` | Thin markdown: pointer to pattern registries / Phase F–G names | Yes |
| `GET /docs/themes` | Thin markdown: Material/Basecoat + `?theme=` | Yes |
| Existing `/docs`, `/components/*` | Unchanged handlers | Yes via choke point |
| `/recipes/*` | Unchanged standalone | No (linked only) |
| `/` | Unchanged | No |

Register stubs next to `GET /docs` in `New()`. Sitemap: add stub URLs if public-indexable (yes — same as `/docs`).

## Theme switcher

No logic change. Shell topbar hosts `{{if .ThemeSwitcher}}{{template "theme-switcher" .ThemeSwitcher}}{{end}}`. Legacy header omits switcher when `DocsNav != nil` (avoid duplicate). `themeSwitcherFor` still strips non-theme query params.

## Testing Strategy

Strict TDD (`go test`, httptest). Table-driven where multiple paths share contracts.

| Layer | What | Approach |
|-------|------|----------|
| Unit | `docsNavFor`, `usesDocsShell` | Table: path → current link, group titles present |
| Integration | Shell chrome | httptest `GET /docs`, `/components/button`, `/docs/patterns`, `/docs/themes` |
| Integration | Home unchanged | `GET /` — no `.docs-topbar` / mobile details shell; no two-pane requirement |
| Integration | Active state | `/components/button` → Button `aria-current="page"`; peers absent |
| Integration | IA groups | Sidebar contains Getting started, Components section titles from `docsSections`, Patterns, Recipes, Themes |
| Integration | Theme | Topbar contains `?theme=` links; `?theme=basecoat` sets root class |
| Integration | Mobile markup | `details` + `summary` present on shell pages |
| Integration | A11y | Skip link + `main#main-content` + docs `nav` |
| Integration | Placeholders | Search not a live submitting corpus form |
| Regression | Footer/breadcrumb/JSON-LD/sitemap | Update expectations; paths stable |

Primary files: `internal/app/docs_shell_test.go` (new) + adjust `server_test.go` layout contracts that hard-code `main.docs-shell` / header nav.

## Threat Matrix

N/A — no shell/subprocess, VCS/PR automation, or executable-file classification boundary. HTTP stub routes are normal docs mux registration (allowlisted paths only).

## Migration / Rollout

1. Add model + `docsNavFor` + tests (RED→GREEN).
2. Partials + layout branch + CSS; enable shell only via `usesDocsShell` in `renderMarkdownStatus` (and error pages only if they already use layout on docs paths — optional, default skip).
3. Register stub routes; rebuild footer from nav model.
4. `npm run build` if CSS entry changes; commit `web/static` if required by repo practice.
5. No feature flag; revert PR to rollback.

**Chained PR forecast (800 budget, auto-forecast):** **High** if CSS+layout+nav+stubs+tests land together. Recommended slices for tasks:

1. Nav model + unit tests + footer rebuild  
2. Layout/partials/CSS/mobile + stub routes  
3. httptest chrome contracts + sitemap/docs residual  

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `internal/app/docs.go` | Modify | `docsNavFor`, stub handlers or shared thin md helper; keep `docsSections` |
| `internal/app/routes.go` | Modify | Flat helpers if still needed; avoid dual source of component lists |
| `internal/app/server.go` | Modify | `pageView.DocsNav`; `usesDocsShell` at render choke; register stubs; footer |
| `web/templates/layout.html` | Modify | Shell branch vs legacy header |
| `web/templates/docs-topbar.html` | Create | Topbar slots |
| `web/templates/docs-sidebar.html` | Create | Grouped nav tree |
| `web/styles/base.css` and/or `docs-shell.css` | Modify/Create | Two-pane, sticky, breakpoints |
| `web/styles/app.css` | Modify | Import if new CSS file |
| `internal/app/docs_shell_test.go` | Create | Chrome + IA + active + home negative |
| `internal/app/server_test.go` | Modify | Layout contracts / main class |
| `docs/gelium-ui-system-roadmap.md` | Modify | Close nav discoverability residual (apply phase) |

## Open Questions

None blocking — decided above.
