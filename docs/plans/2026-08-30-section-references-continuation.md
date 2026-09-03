# Section references — continuation guide

Repo: `/root/.openfang/workspaces/repos/gelium-ui` (branch `main`).
Docs live (Tailscale): `http://100.121.211.121:8788` — may need `ADDR=100.121.211.121:8788 go run ./cmd/gelium` if down.
Chat language with the human: short Spanish. Artifacts (code, docs, this guide): English.

The first ficha is shipped. Follow this file. Do **not** invent a Saaspo clone or reopen the wireframe.

## Product intent

A first-party catalog of **section types**, not a screenshot gallery.

Each ficha must show, on the same page:

1. **Original (cited)** — official live URL + **block map** of section order (or our own crop). Never a Saaspo/Land-book/Mobbin CDN image.
2. **Gelium remake** — visible HTML using Gelium tokens/primitives (`#gelium-remake`). Same jobs as the cited structure, **our copy**, no their brand/logo/assets/motion. This is “how it would look”, not a link to an unrelated recipe.
3. **Keep / adapt** table — kept → primitive, rejected → why.
4. **Ask before copying** — product questions, not markup.

`/recipes/rich-article` is a **different** media fixture. Do not point new fichas at it.

## What exists

| Route | Role |
|---|---|
| `GET /docs/section-references?type=` | Index (SCREEN list). Closed filter: `all`, `article`, `hero`, `pricing`, `auth`, `faq`, `404`. Unknown → `all`. Empty copy when a type has no published row. |
| `GET /docs/section-references/{id}` | Detail. Unknown/unpublished → docs 404. |
| `/docs/section-references/article` | First ficha (`REF-ARTICLE`). Cite Vercel article. On-page remake of **that** structure. |

Files:

- Catalog: `internal/app/section_references.go` (`sectionRefCatalog`, `sectionRefTypes`)
- Index intro: `site/web/content/handbook-section-references.md`
- Article ficha: `site/web/content/handbook-section-references-article.md` ← **copy this shape**
- Tests: `internal/app/docs_section_references_test.go`
- Nav: Core in `internal/app/docs.go` (`handbookSections`), after Page + section architecture
- Mux: `GET /docs/section-references` and `GET /docs/section-references/{id}` in `internal/app/server.go`
- Sitemap: index via `handbookNavLinks()`; details via `publishedSectionRefPaths()`
- Agent pack: `REF-ARTICLE` in **both** `lib/llms-ux.txt` and `site/web/static/llms-ux.txt`
- Also: `lib/llms.txt`, `site/web/static/llms.txt`

Copy the article ficha. Goldmark allows trusted HTML for the remake (`<article id="gelium-remake">`). Inner headings in that HTML must not impersonate handbook H2s if they pollute TOC — match the article file’s pattern.

## Approved wireframe (do not reopen)

- SURFACE Read. Index = SCREEN **list**. Detail = SCREEN **detail**.
- Index: GET filter, list rows (title + meta + chevron/link), no screenshot grid, no card-as-section.
- Detail primary: jump to / show Gelium remake. Secondary: official original.
- Mobile: same regions; original then remake stacked.
- Docs chrome unchanged.

## How to add the next ficha

1. Pick a **type** already in `sectionRefTypes`.
2. Open the **official** URL (Vercel, Linear, GOV.UK…). Saaspo/Land-book = URL index only. Mobbin is not default. Re-verify the live page; do not trust this backlog’s URLs blindly.
3. Write a block map of **section order**. No their logo/CSS/assets in the remake. Copy is ours (structure stolen, prose not).
4. Add `site/web/content/handbook-section-references-<id>.md` with the same H2 set: Original, Gelium remake (visible `#gelium-remake` HTML), Keep / adapt, Ask before copying.
5. Append a `Published: true` row to `sectionRefCatalog` (`ID` = path segment, `ContentPath` = that markdown).
6. Add `REF-<TYPE>` to **both** `lib/llms-ux.txt` and `site/web/static/llms-ux.txt`.
7. Extend `docs_section_references_test.go`: 200, `id="gelium-remake"`, official host, Keep / adapt marker, table/list class used in the remake. Do **not** assert `/recipes/rich-article`.
8. Gates: `gofmt -w` only touched Go files; `go test ./internal/... ./site/... ./lib/...` (never bare `./...`); `bash scripts/ux-detect.sh`.

Do **not** add a new mux route; `{id}` already dispatches.
Do **not** put detail pages in `handbookNavLinks` / `handbookRoutes` (index only). Details join the sitemap through `publishedSectionRefPaths()`.
Sitemap count test: `componentRoutes + handbookNavLinks + 8 + blogPosts + publishedSectionRefPaths`.
Core insert already retargeted `TestPageSectionArchitectureDocsModel` next → section-references, not Journeys. Do not put it back.

## Hard no

- Republishing Saaspo / Land-book / Mobbin images
- Pixel-clone of the cited product
- Style filters (gradients, bento, motion)
- Operate admin/queue “inspiration”
- Publishing auth/pricing/hero **ficha** before you can reconstruct that structure in Gelium primitives on the page (auth **routes** `/recipes/auth/*` are documented in `docs/gelium-ui-screen-recipes.md` but **not** registered in `server.go` — do not claim they exist)
- Mandatory browser crawl in the agent workflow
- Linking an unrelated existing recipe as the remake

## Page backlog (one type per slice)

Official URLs are starting citations. Re-verify live before the block map.

| Order | Type | ID (proposed) | Official cite (verify live) | Remake | Blocker |
|---|---|---|---|---|---|
| 1 | article | `article` | https://vercel.com/blog/the-end-of-credential-sprawl-for-agents | on-page `#gelium-remake` | **Done** |
| 2 | 404 | `404-vercel` | confirm a real Vercel 404 URL | `error-state` composition on the ficha | next recommended |
| 3 | faq | `faq-govuk` or cited pricing FAQ | GOV.UK or one SaaS FAQ page | native `<details>` / accordion already in library | pick one URL |
| 4 | auth | `auth-register` | one official register URL (not Land-book’s JPG) | on-page form: email, password, terms checkbox, one primary, forgot link, 422 story in Keep / adapt | do **not** wait for missing `/recipes/auth` if the ficha itself is the remake |
| 5 | hero | `hero-linear` | https://linear.app | hub: H1 + short context + **one** primary. No video loop, no gradient | Persuade still uses tokens |
| 6 | pricing | `pricing-linear` | https://linear.app/pricing | plans as cards or table + **one** CTA. No bento 3D | one primary |
| 7 | article | `article-govuk` | a GOV.UK guidance article | second Read cite; same remake pattern | optional |

**Do this next:** slice **404** (block map of a real 404 + Gelium `error-state` remake on the ficha). Then auth register as on-page form reconstruction.

Do not start a 1.300-page catalog. One ficha per PR-sized slice.

## Optional Plan hook (workflow, not this slice)

After SURFACE/SCREEN, before wireframe: if a browser exists, open at most two **official** URLs of that SCREEN, filter with Gelium IDs, tell the human brought/rejected. If not, use this catalog. Fail-open. Do not implement that hook unless the human asks.

## Verify live

If serving: `ADDR=100.121.211.121:8788` only (never `0.0.0.0`). Do not steal port 8787 (existing gelium). Restart `go run` after Go/markdown route changes (`npm run build` does not reload the mux).
