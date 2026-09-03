# REF-SOCIAL-FEED

Portable structural reference for `SCREEN=list` social / activity feeds.
Origin: responsive visual audits (X-like and Reddit-like families). Screenshots stay
development-only outside this package (copyright unknown). This ficha ships the
**confirmed patterns** an agent needs without monorepo or docs site access.

## Families (do not merge into one shell)

| Family | Narrow | Medium | Distinct IA (usually reject for product) |
|---|---|---|---|
| X-like | single-column feed, compact header, topic tabs, bottom nav | icon-only left rail, constrained central column, composer above feed | bottom nav, left rail, composer-on-feed, quote embeds |
| Reddit-like | single-column, sort/filter under header, large contained media | global header+search, central feed, pill actions | sort chrome, join-on-card, vote as primary |

Treat families as separate references with **overlapping item anatomy**.

## Confirmed item anatomy (cross-family)

Reading order for each post:

1. **Metadata** (author/community, relative time, optional context) — muted
2. **Title / body** — high contrast primary reading
3. **Optional media or quoted block** — contained, not full-bleed chrome
4. **Actions** — subordinate; real touch targets even if glyph is small

Wider layouts may use **avatar + content column**. Narrow layouts collapse to one column.
Prefer light dividers over heavy elevation between items.

## Confirmed composition patterns

- Single-column feed at narrow widths.
- Constrained central feed at medium (≈ content measure, not full viewport chrome).
- Muted metadata / high-contrast primary content.
- Rounded media wrappers; portrait media letterboxed inside a landscape frame when needed.
- Tabs or sort rows need narrow-width overflow handling (scroll/clip) — only if product owns those views.
- Fixed bottom navigation and icon rails are **shell choices**, not feed-item requirements.

## States the feed surface must still own

Empty, loading (or sync GET with no fake client spinner), error + recovery, success after local mutation (like/save), pagination or next-page when the product has it.

## Gelium product filter (mandatory)

**Keep / map**

- Item order: meta → title/body → optional media → actions
- Existing consumer shell (header/nav); do not add a second system from the screenshot
- Real product fields only (author/space, date, title, summary, safe media URL)
- Local mutations as subordinate POST controls outside the reading link
- Media via `.ui-card-media` / registered media primitives when a **safe URL** exists — never invent width/height
- Icons only through the Gelium catalog + accessible name; text-first is valid default

**Reject unless product data + routes already own them**

- Second header, left icon rail, fixed bottom nav
- Composer on the feed (publish belongs on its own route if that is the product)
- Quoted nested posts, ads, “suggested/join” chrome on every card
- Vote up/down when the product only has like/save
- Comment counts without a feed query field
- Sort/topic tabs without a real `?view=` / sort contract
- Fake avatars, copied handles, metrics, logos, verification marks
- Constant badges on every row (constant-state lens) — move to section context
- Whole-card link when nested POST controls exist (use REF-CARD-DETAIL-ENTRY)

## G5 packet note (example)

```text
References: REF-SOCIAL-FEED — item order + media containment + single-column narrow; passed filter.
Rejected: bottom nav, left rail, composer-on-feed, vote, fake avatar (no data).
No-match visual PNGs: package text-only; monorepo docs/reference-assets optional supplement.
```

## Companion refs

- `REF-CARD-DETAIL-ENTRY` — one reading target vs nested mutations
- `REF-SHELL` — do not duplicate chrome the shared layout already owns
- Recipe public-feed (docs site) — richer live composition when available; not required if this ficha is present
