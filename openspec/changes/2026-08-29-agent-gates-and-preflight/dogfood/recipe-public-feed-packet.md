# Advisory dogfood packet — recipe public feed

## Identity

```text
Packet version: 1
Change: Advisory workflow dogfood against the existing Gelium recipe public feed.
Date: 2026-08-30
Author: Hermes Agent
Route: design-gated
Scope: Existing feed audit; no product, route, template, CSS, or handler behavior change.
Existing route and contracts: GET /recipes/public-feed; GET query views/pagination; POST reactions with 303; POST refresh; HTMX fragments; class-routed theme; no-JS fallbacks.
Product job / audience: A reader scans recent activity, switches feed views, reacts to a post, and refreshes the demo feed.
Primary action: Read the activity feed; supporting actions are view selection, Like, and Refresh feed.
Non-goals: Redesigning the feed; changing data, auth, privacy, routes, mutation contracts, or existing components.
```

## Orient receipts

```text
Product/design artifacts: This repository recipe is the equivalent consumer. No product change is proposed; the current-chat dogfood record supplies the bounded purpose.
Gelium entrypoint / decision pack / skills: lib/AGENTS.md; lib/llms-ux.txt; lib/SKILLS.md; skills 11 and 12.
Vocabulary and component registry: Existing ui-tabs, ui-card, ui-avatar, ui-badge, ui-button, ui-skeleton, pagination, toast, and progress primitives.
Nearby reusable patterns: internal/app/recipe_public_feed.go; internal/app/recipe_public_feed_test.go; site/web/templates/recipe-public-feed.html; site/web/styles/recipe-public-feed.css.
Hard route, permission, data, and no-JS/server contracts: GET views/pagination, POST+303 reactions, GET/POST refresh behavior, HX fragment response, server-rendered initial paint, class-routed theme.
Open uncertainty: This is a static local audit. It can verify rendered HTML and HTTP state, but not a human visual review at every viewport/theme or real authenticated production data.
```

## Plan — intent wireframe

```text
Screen / SURFACE: Public/server-rendered activity feed.
User job and audience: Readers scan updates and take a lightweight reaction or refresh action.
Major regions in reading order: Brand/header → title/description → feed-view tabs → activity list or empty state/pagination → documented loading placeholder → refresh form and toast region.
Primary and supporting actions: Read is primary; view tabs, Like, and Refresh feed are supporting.
States and recovery: Initial server render; documented Skeleton/loading explanation; empty state; pagination; POST+303 reaction toast; refresh full-page no-JS response or HTMX fragment.
Desktop structural wireframe: Narrow centered reading column, header, tabs, cards, recovery/loading section, then refresh.
Mobile structural wireframe: Same DOM order; wrapping app bar and action rows; responsive content padding.
Constraints and non-goals: Audit only. Do not alter the established recipe or fabricate image metadata.
```

## Architect — buildable wireframe

```text
Route / handler / template: server.go routes GET /recipes/public-feed and POST reaction/refresh to recipe_public_feed.go; template is site/web/templates/recipe-public-feed.html.
Data and permission boundary: In-memory recipe demo store; no auth or privacy boundary is changed by this audit.
Section inventory and SECTION-CONTRACT mapping: Header establishes purpose; tab nav selects server-side view; feed panel carries ordered activity and reaction forms; loading section explains recovery; refresh section invokes server refresh.
Component and token mapping: ui-tabs, ui-card, ui-avatar, ui-badge, ui-button, ui-skeleton, pagination, toast, and progress; recipe layout CSS reads --ui-* tokens.
URL, form, POST+303, 422, and validation contracts: GET query view/page links; reaction POST redirects 303; refresh POST re-renders no-JS page or HTMX fragment. These forms do not collect user-entered fields, so validation-summary is not applicable.
No-JS and accessibility behavior: Native links/forms provide navigation and actions; headings, nav labels, list/article semantics, live toast region, status skeleton, and wrapping layout remain server-first.
Desktop buildable wireframe: Existing rendered candidate exactly as implemented; no requested change.
Mobile buildable wireframe: Existing flex-wrap/min-width:0 and max-width:40rem rules; no requested change.
Material mismatch, exception, or escalation: The generic detector reports custom-shell and form-validation findings for intentional recipe layout and action-only forms. They are retained as bounded, expiring advisory exceptions rather than being silently suppressed or converted into UI changes.
```

## Criteria plan (prebuild)

```text
Hierarchy and DOM order: Header, view selection, feed, recovery/loading, refresh.
Action hierarchy and boundaries: Reading is primary; tabs and mutations are secondary and native/server-first.
Responsive and class-routed theme intent: Use existing theme class on html; verify CSS max-width and narrow padding without media-query dark mode.
States and recovery: List, empty, pagination, reaction success, refresh success, documented Skeleton/loading.
Accessibility and no-JS: Native controls and semantic landmarks; test POST+303 and no-JS response paths.
Preserved contracts: Route/query/POST/HX contracts, toast event/header, component tokens, and no product behavior change.
DESIGN-MEMORY reuse decision: Existing recipe primitives and CSS are audited as-is; no new component or visual pattern is introduced.
```

## Decision

```text
Decision: exception
Approver: CryptoNahue
Date / channel: 2026-08-30 / current Hermes chat
Approved scope or bounded exception: Advisory audit of the existing recipe-public-feed equivalent consumer. No UI implementation is authorized or required.
Reason, risk boundary, owner, and follow-up (if exception): The record uses a design-gated ledger to dogfood the workflow against an existing screen. Static local evidence cannot substitute for a human pixel-level review across all viewport/theme combinations. Owner: maintainers. Follow-up: record only demonstrated detector friction before deciding any required rollout.
```

## Postbuild handoff

The ledger and audit record link the local HTTP/rendered evidence, detector findings, authority-matrix result, focused test, full test, and build results. They do not authorize a commit, publish, deploy, or claim production visual approval.
