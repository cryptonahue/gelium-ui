# Tasks: docs-shell

## Review Workload Forecast

| Field | Value |
|-------|-------|
| Estimated changed lines | 700–950 authored (CSS+layout+nav+stubs+tests; exclude goldens) |
| 400-line budget risk | High |
| 800-line budget risk | High (single PR) |
| Chained PRs recommended | Yes |
| Suggested split | PR1 nav model → PR2 shell UI/CSS/stubs → PR3 httptest chrome + polish |
| Delivery strategy | auto-chain (session: auto-forecast) |
| Chain strategy | **stacked-to-main** (default; solo maintainer, main ahead of origin) |

Decision needed before apply: No
Chained PRs recommended: Yes
Chain strategy: stacked-to-main
400-line budget risk: High

### Suggested Work Units

| Unit | Goal | Likely PR | Focused test command | Runtime harness | Rollback boundary |
|------|------|-----------|----------------------|-----------------|-------------------|
| 1 | Nav model + `usesDocsShell` + footer source | PR 1 | `go test ./internal/app/ -run 'TestDocsNav\|TestUsesDocsShell\|TestDefaultFooter' -count=1` | N/A — pure model/helpers | `docs.go`, `routes.go`, footer builder in `server.go`, unit tests |
| 2 | Shell layout/CSS/topbar/sidebar + stubs | PR 2 | `go test ./internal/app/ -run 'TestDocsShell\|TestDocsStub\|TestHomeUnchanged' -count=1` | `go run ./cmd/gelium` → `/docs`, `/components/button`, `/` | templates, CSS, `web/static`, stub routes, render choke |
| 3 | Full chrome contracts + sitemap/roadmap | PR 3 | `go test ./internal/app/ -count=1` then `go test ./...` | Spot-check landmarks, theme `?theme=basecoat`, mobile details | `docs_shell_test.go`, `server_test.go`, sitemap, roadmap |

**Req map:** R1 frame · R2 IA · R3 active · R4 topbar · R5 theme · R6 mobile · R7 dogfood · R8 landmarks · R9 no-JS · R10 OOS

---

## Phase 1: Nav model (PR 1) — R2, R3

- [x] 1.1 RED: `docs_shell_test.go` — table `usesDocsShell` (`/docs`, `/docs/*`, `/components/*` true; `/`, `/recipes/*` false)
- [x] 1.2 RED: table `docsNavFor` — five IA blocks; `docsSections` group titles; exact `Current` on `/docs` and `/components/button`
- [x] 1.3 GREEN: types `docsNavLink`/`docsNavGroup`/`docsNavView` + `docsNavFor`/`usesDocsShell` in `docs.go` (Version `0.4.0`, `SearchDisabled`)
- [x] 1.4 GREEN: derive `navLinks()` from `docsSections` in `routes.go` (no dual component list)
- [x] 1.5 GREEN: rebuild `defaultFooter` from `docsNavFor` flat export in `server.go`; unit-assert sections
- [x] 1.6 `gofmt`; `go test ./internal/app/ -run 'TestDocsNav|TestUsesDocsShell|TestDefaultFooter' -count=1`

## Phase 2: Shell chrome (PR 2) — R1, R4–R9

- [x] 2.1 RED: httptest shell on `/docs` + `/components/button` (topbar, sidebar/nav, main); `/` lacks docs-topbar/mobile shell
- [x] 2.2 RED: stubs `/docs/patterns`, `/docs/themes` 200 + shell; search not live submit; skip + `nav` + `main#main-content`
- [x] 2.3 GREEN: `pageView.DocsNav`; set in `renderMarkdownStatus` when `usesDocsShell`; theme switcher only in topbar when shell
- [x] 2.4 GREEN: create `docs-topbar.html` (brand `/`, disabled search, version badge, theme-switcher)
- [x] 2.5 GREEN: create `docs-sidebar.html` (`ui-list*` links, group labels, `ui-divider`, `aria-current="page"`)
- [x] 2.6 GREEN: branch `layout.html` — shell: topbar + dual mobile `<details>`/desktop `<aside>` + main; legacy header when `DocsNav == nil`
- [x] 2.7 GREEN: two-pane CSS in `base.css` and/or `docs-shell.css` + `app.css` import; sticky topbar/sidebar; md ~48rem; home centered utility preserved
- [x] 2.8 GREEN: register stub handlers + thin markdown; sitemap URLs; Recipes outbound `/recipes/*` only
- [x] 2.9 `npm run build`; commit `web/static` if changed; focused tests from 2.1–2.2

## Phase 3: Contracts + polish (PR 3) — R1–R10, residual

- [x] 3.1 RED→GREEN: active `aria-current` peers; IA group strings; theme links `?theme=` only; `?theme=basecoat` root class
- [x] 3.2 Update `server_test.go` layout contracts (`main.docs-shell` / header assumptions)
- [x] 3.3 Sitemap/JSON-LD/footer regressions; paths stable (no redirects)
- [x] 3.4 Close nav-discoverability residual in `docs/gelium-ui-system-roadmap.md`
- [x] 3.5 `go test ./...` · `go vet ./...` · `go build -o /tmp/gelium ./cmd/gelium`
