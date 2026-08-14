# Proposal: docs-shell

## Intent

Replace flat header + centered docs column with a Scalar-style two-pane shell. Dogfood Gelium primitives, keep URLs, fix nav discoverability—no landing rewrite.

## Scope

### In Scope

- Sticky sidebar + topbar (brand, search slot, version badge, theme switcher)
- IA: Getting started, Components (`docsSections`), Patterns, Recipes, Themes
- Active path; mobile `<details>/<summary>` (0 JS); breadcrumb in shell
- Compose drawer/list/divider/section-heading/theme switcher/footer
- Patterns/Themes stubs or deep-links; Recipes may link `/recipes/*`
- Non-broken search/version placeholders; httptest chrome tests

### Out of Scope

- `/` changes; URL moves; real search; TOC; prev/next
- Recipe layout merge; new themes; registry JSON; `loom:*` renames

## Capabilities

### New Capabilities

- `docs-shell`: Docs chrome, grouped nav, topbar slots, mobile details, active path

### Modified Capabilities

- None

## Approach

1. Keep `/docs`, `/components/*`; chrome wrap only; optional `/docs/...` stubs
2. Nav from `docsSections` + Patterns/Recipes/Themes; replace `navLinks()`; shared footer
3. Sticky drawer/list + labels; topbar hosts theme switcher; evolve `.docs-shell` to two-pane
4. Shell on docs routes only; `/` unchanged; recipes linked not merged
5. httptest chrome tests; SEO paths stable
6. Auto-forecast 800 lines; chain if needed: nav → CSS/mobile → tests

## Assumptions (binding)

| Topic | Decision |
|-------|----------|
| Audience | Mix public docs; bias embedder |
| URLs | Keep routes; chrome only |
| Sidebar | Scalar full IA |
| Gaps | Stubs/deep-links OK |
| Search/version | Slots OK if not broken |
| Mobile | details/summary only |
| Home | No `/` changes |

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `layout.html`, `base.css` | Modified | Two-pane sticky frame |
| `docs.go`, `routes.go`, `server.go` | Modified | Nav model, active path |
| drawer/list/theme/breadcrumb | Reuse | Chrome composition |
| `internal/app/*_test.go` | Modified | Shell contracts |
| roadmap residual | Modified | Nav discoverability |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Layout test churn | High | Same-slice httptest |
| Dense drawer / footer drift | Med | Labels; one nav model |
| `.docs-shell` CSS shift | Med | Evolve carefully |
| 800-line overrun | Med | Chain PRs |

## Rollback Plan

Revert PR (layout, CSS, `navLinks()` header, tests). No migrations/URL moves; `New()` same.

## Dependencies

Gelium drawer/list/breadcrumb/theme switcher; `docsSections`; theme middleware; CSS build if needed.

## Success Criteria

- [ ] Shell on `/docs` + `/components/*`; `/` unchanged
- [ ] Scalar-full groups; active path; details mobile 0 JS
- [ ] Theme switcher + breadcrumb work; placeholders not broken
- [ ] Dogfood primitives; `go test ./...` pass; SEO stable
- [ ] Recipes linked without layout merge
