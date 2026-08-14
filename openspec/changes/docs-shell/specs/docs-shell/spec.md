# docs-shell Specification

## Purpose

Two-pane docs chrome on `/docs` + `/components/*`: grouped IA, topbar slots, 0-JS mobile. Dogfood Gelium; keep URLs. No landing/search/TOC/recipe-merge.

## Requirements

### Requirement: Docs shell frame on docs routes

MUST render sticky sidebar + topbar + main on `/docs` and `/components/*`. MUST NOT change `/`. No path changes/redirects.

#### Scenario: Shell on docs hub
- GIVEN GET `/docs`
- WHEN HTML renders
- THEN sidebar, topbar, main exist

#### Scenario: Shell on component page
- GIVEN GET valid `/components/*`
- WHEN HTML renders
- THEN same shell frame as `/docs`

#### Scenario: Home unchanged
- GIVEN GET `/`
- WHEN HTML renders
- THEN two-pane shell not required; home unchanged

### Requirement: Scalar-full sidebar IA

Sidebar MUST show Getting started/Docs, Components (`docsSections` groups), Patterns, Recipes, Themes. Audience: public docs, embedder bias, clear labels. Patterns/Themes MAY stub/deep-link. Recipes MAY deep-link `/recipes/*` outside shell.

#### Scenario: Required groups visible
- GIVEN shell page
- WHEN sidebar inspected
- THEN all five groups present

#### Scenario: Components use section groups
- GIVEN Components group
- WHEN destinations listed
- THEN use `docsSections` categories, not a flat dump

#### Scenario: Recipes may leave the shell
- GIVEN Recipes link to `/recipes/*`
- WHEN followed
- THEN target MAY leave shell; href MUST be real

### Requirement: Active route indication

MUST mark current-path nav item as current/active.

#### Scenario: Active component link
- GIVEN GET `/components/button`
- WHEN nav renders
- THEN Button item is current; peers are not

#### Scenario: Active docs hub
- GIVEN GET `/docs`
- WHEN nav renders
- THEN Docs hub item current

### Requirement: Topbar chrome slots

Topbar MUST host brand, search slot, version badge, theme switcher. Placeholders OK disabled/coming-soon (Gelium); MUST NOT fake-break.

#### Scenario: Topbar contents present
- GIVEN shell page
- WHEN topbar inspected
- THEN all four slots present

#### Scenario: Honest placeholders
- GIVEN unimplemented search/version
- WHEN slot shown
- THEN disabled or coming-soon; no broken submit

### Requirement: Zero-JS theme switcher in topbar

MUST keep 0-JS `?theme=` links in topbar; allowlisted slugs only.

#### Scenario: Theme links in topbar
- GIVEN shell page path P
- WHEN switcher renders
- THEN each option is GET P with only `?theme=<slug>`

#### Scenario: Theme applies without JS
- GIVEN GET shell path with `?theme=basecoat`
- WHEN HTML renders
- THEN root gets matching theme class; no client script

### Requirement: Mobile details/summary nav

MUST use `<details>`/`<summary>` on narrow viewports; 0 required JS. Modal drawer MUST NOT be baseline-required.

#### Scenario: Details menu without script
- GIVEN shell page, scripts off
- WHEN mobile nav summary opens
- THEN grouped nav as full-page links via details

### Requirement: Dogfood Gelium primitives

MUST compose Gelium primitives/tokens for chrome. MUST NOT invent a parallel docs-only design system.

#### Scenario: Primitive composition
- GIVEN `/docs` chrome
- WHEN inspected
- THEN sidebar/topbar use Gelium ui-* surfaces, not separate docs library

### Requirement: Landmarks and skip link

MUST keep skip link to main; expose nav + main landmarks.

#### Scenario: Skip link preserved
- GIVEN shell page
- WHEN start inspected
- THEN skip link to main exists

#### Scenario: Nav and main landmarks
- GIVEN shell page
- WHEN landmarks inspected
- THEN docs nav is `nav`; content is `main`

### Requirement: No-JS baseline for shell chrome

Shell chrome MUST work with JS off.

#### Scenario: Full chrome without JS
- GIVEN scripts off
- WHEN `/docs` then component nav is opened
- THEN shell and content work via HTTP

### Requirement: Out of scope stays out

MUST NOT: landing redesign, real search, TOC, URL moves/redirects, recipe layout merge, `loom:*` renames.

#### Scenario: No invented search
- GIVEN placeholder search
- WHEN interacted
- THEN no real corpus search

#### Scenario: No path migration
- GIVEN `/docs` and `/components/*`
- WHEN shell ships
- THEN same paths serve same resources; no redirect matrix
