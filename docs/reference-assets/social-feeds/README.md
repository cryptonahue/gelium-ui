# Social feed visual references

Development-only references supplied by the user on 2026-09-03. These captures are evidence for responsive composition and visual comparison; they are not production assets and do not establish backend, brand, permission, or content requirements.

## Provenance

- Origin: `user-supplied-screenshot`
- Source URL: not provided
- Verification: visual-only
- Usage: development comparison only
- Copyright/terms status: unknown; do not publish or bundle without permission
- Exact viewport breakpoints: unknown; image dimensions are recorded, but device/browser chrome may be present

## Reference set

| ID | File | Pixels | Family | Confidence |
|---|---|---:|---|---|
| `REF-SOCIAL-FEED-X-MOBILE` | `x-mobile.png` | 399×850 | X-like social feed | high for visible mobile composition |
| `REF-SOCIAL-FEED-X-MEDIUM` | `x-medium.png` | 830×850 | X-like social feed | high for visible medium composition |
| `REF-SOCIAL-FEED-REDDIT-MOBILE` | `reddit-mobile.png` | 413×857 | Reddit-like community feed | high for visible mobile composition |
| `REF-SOCIAL-FEED-REDDIT-MEDIUM` | `reddit-medium.png` | 848×842 | Reddit-like community feed | high for visible medium composition |

## Image audit

### X-like set

**Mobile — `x-mobile.png`**

Confirmed visually:

- Single-column vertical feed.
- Compact top header with avatar and page title.
- Horizontal topic tabs; the row is clipped or horizontally scrollable.
- Feed post uses avatar + content column.
- Nested quoted-post card with rounded border and media preview.
- Post action row with icon/count controls.
- Fixed bottom navigation with five icon slots.
- Dark theme, thin separators, bright active accent, muted metadata.

Inferred, not proven by this capture:

- Exact breakpoint that activates bottom navigation.
- Whether the feed or whole document owns scrolling.
- Whether the media preview is a real video or a static image with a duration label.
- Exact semantics of unlabeled icons.

**Medium — `x-medium.png`**

Confirmed visually:

- Narrow icon-only left rail.
- Central feed constrained to roughly 600px in the visible capture.
- No visible right rail.
- Topic tabs and composer remain above the feed.
- Post body and quoted card align to the content column rather than the avatar column.
- The quoted media card grows with the available feed width.

Inferred, not proven:

- The left rail is sticky/fixed.
- This is a tablet breakpoint rather than a cropped desktop view.
- The unused right-side space is intentional.

### Reddit-like set

**Mobile — `reddit-mobile.png`**

Confirmed visually:

- Single-column community/feed view.
- Compact top header with menu, identity, messaging, create and notification controls.
- Sort/filter row below the header.
- Post metadata with community, timestamp, context label, join action and overflow action.
- Large rounded media frame using centered portrait media with side letterboxing.
- Video-like overlay controls and duration label.
- Pill-shaped vote/comment/share actions.
- Vertical feed continues below the viewport.

Inferred, not proven:

- Exact behavior of the sort/filter controls.
- Whether the media is video or a captured player state.
- Exact breakpoint and whether the outer gray frame belongs to the product or capture environment.

**Medium — `reddit-medium.png`**

Confirmed visually:

- Global header with menu, search, utility actions and account control.
- Single central feed with no visible sidebars.
- Sort/view controls below the header.
- Large media post with portrait content contained inside a landscape frame.
- Rounded media clipping and overlay playback controls.
- Compact pill action bar below the media.
- Additional feed content continues below the capture.

Inferred, not proven:

- The hidden desktop sidebars are responsive behavior rather than a cropped state.
- Exact search collapse/hide rules at narrower widths.
- Exact media aspect-ratio contract.

## Cross-reference findings

Reusable patterns supported by multiple captures:

- Feed item: metadata → title/body → optional media/quote → actions.
- Two-column post anatomy at wider widths: avatar + content.
- Single-column feed at narrow widths.
- Muted metadata and high-contrast primary content.
- One-pixel dividers instead of heavy card elevation.
- Rounded media wrappers with contained portrait content.
- Tabs or sort controls that need narrow-width overflow handling.
- Action controls that require real touch targets even when the glyph is small.

Do not merge these into one product shell. Treat the X-like and Reddit-like sets as separate reference families with overlapping patterns and different information architecture.

## Gelium adaptation constraints

- Use the existing Gelium shell and registry first; a screenshot cannot authorize a second header, rail, or bottom navigation.
- Use real product data only. Do not copy visible names, handles, posts, metrics, verification marks, logos, or media.
- Treat icons as semantic candidates; resolve them through the chosen Gelium icon catalog rather than copying glyph shapes blindly.
- Preserve server-first URLs, native controls, no-JS completion, 422 recovery, and real empty/loading/error/success states.
- Do not infer exact CSS pixels as tokens. Map spacing, type, colors and radii to existing `--ui-*` tokens.
- Validate wide/narrow viewports, dark/light class-routed themes, keyboard/focus, touch targets, and content clearance around fixed navigation.

## Intended next test

Use one reference family at a time to produce a Gelium social-feed candidate. Compare the candidate against the relevant mobile and medium captures for structure, hierarchy, density, media containment, and responsive transitions. Record matched reference IDs, rejected patterns, and unresolved product decisions in the G5 packet. Do not claim pixel matching or brand fidelity.
