# Gelium UI — Component Implementation Prompt for AI agents

Reusable prompt for contributing **exactly one Gelium UI component** with an AI
agent (Claude Code, OpenCode, Codex, Hermes, etc.). Copy from "COPYABLE PROMPT",
fill the `{{...}}` placeholders, and hand it to the agent responsible for a single
component.

If you only need fast, safe UI guidance for an existing project, you do not need
this prompt — point the agent at the shipped guidance in the npm package:
`AGENTS.md` (golden rules), `skills/` (7 decision skills), `llms-ux.txt` (full
decision pack). This prompt is for **contributing a new component to Gelium UI
itself**.

## How to use it

1. Pick a component from the canonical registry: `docs/gelium-ui-component-registry.md`.
2. Fill every `{{...}}` assignment placeholder. The final status is chosen by the
   worker after it finishes.
3. Assign a concurrency mode and explicit ownership.
4. Never let two workers write the same physical file.
5. Designate a single integrator for the canonical checkout.
6. The worker must end in exactly one of four allowed statuses
   (`COMPLETE_AWAITING_USER_ACCEPTANCE`, `READY_FOR_INTEGRATION`, `BLOCKED`,
   `ABORTED_ON_DRIFT`) and never declare "done" without evidence.

## Concurrency modes

- `SHARED_HANDOFF`: the worker only creates/edits new or exclusive files. For
  shared files it hands off patches; it does not apply them.
- `ISOLATED_PHYSICAL_WORKSPACE`: works in an authorized physical copy, never the
  canonical checkout. Delivers manifests and patches.
- `EXCLUSIVE_INTEGRATION`: the only mode allowed to modify canonical shared files;
  requires a literal, current reservation for each path.

Today the recommended shape for several AIs is parallel research + isolated
workspaces + serial integration.

---

# COPYABLE PROMPT

````text
# Gelium UI component agent — safe operational contract

Implement exactly ONE Gelium UI component under this contract.

## 0. Parameters

COMPONENT: {{COMPONENT_NAME}}
SLUG: {{COMPONENT_SLUG}}
CATEGORY: {{CORE | GELIUM_ONLY}}
REQUIRED VARIANTS: {{REQUIRED_VARIANTS}}
REQUIRED STATES: {{REQUIRED_STATES}}
REQUIRED BEHAVIORS: {{REQUIRED_BEHAVIORS}}
SERVER/HTMX FLOW, IF ANY: {{SERVER_FLOW_OR_NONE}}
ACCEPTANCE CRITERIA: {{ACCEPTANCE_CRITERIA}}
SUPPORTED BROWSER MATRIX: {{SUPPORTED_BROWSER_MATRIX}}
WEB BASELINE SNAPSHOT DATE: {{WEB_BASELINE_SNAPSHOT_DATE}}

CANONICAL REPOSITORY: {{CANONICAL_REPO_PATH}}
COMPONENT REGISTRY: {{CANONICAL_REPO_PATH}}/docs/gelium-ui-component-registry.md
AGENT GUIDANCE: {{CANONICAL_REPO_PATH}}/lib/AGENTS.md
SKILLS: {{CANONICAL_REPO_PATH}}/lib/skills/*.md
DECISION PACK: {{CANONICAL_REPO_PATH}}/lib/llms-ux.txt

ASSIGNED MODE: {{SHARED_HANDOFF | ISOLATED_PHYSICAL_WORKSPACE | EXCLUSIVE_INTEGRATION}}
AUTHORIZED WORKSPACE: {{AUTHORIZED_WORKSPACE_PATH}}
WORKER ID: {{WORKER_ID}}
RESERVATION ID: {{RESERVATION_ID}}
OWNED PATHS: {{OWNED_PATHS}}
ADDITIONAL FORBIDDEN PATHS: {{FORBIDDEN_PATHS_OR_NONE}}
BASELINE SHA-256: {{EXPECTED_HASHES_OR_NONE}}
PROPOSED NEW ASSET VERSION: {{NEW_ASSET_VERSION_OR_INTEGRATOR_OWNED}}

Do not expand scope to a second component. Record non-essential cross-cutting
needs as follow-ups.

## 1. Mandatory reading — read BEFORE you write

Read in order, do not skip:

1. `lib/AGENTS.md` — the golden rules that apply to every task.
2. `lib/skills/` — every skill; they are short and actionable. Start with
   `01-foundations.md`. Confirm compliance with `07-dod-and-antislop.md` before
   calling anything done.
3. `lib/llms-ux.txt` — the full decision pack (SURFACE / SCREEN / WF / DATA /
   FEED / JOURNEY / MEDIA / SKEL ids).
4. `docs/gelium-ui-component-registry.md` — the row and states for this component.
5. `docs/gelium-ui-vocabulary.md` — the canonical meaning of the pattern(s) used.
6. `README.md` and `lib/README.md` — product and package contract.
7. Existing similar components: `internal/app/<x>.go` handler,
   `lib/templates/<x>.html` partial, `lib/styles/<x>.css`, and the corresponding
   `lib/styles_<x>_test.go`.

Report the mode, ownership, protected paths, and port/process state (see §3)
before producing anything.

## 2. Immutable constraints

**Target stack (fixed):** Go `net/http` + `html/template`; server-rendered
components; public `--ui-*` tokens; Tailwind CSS 4 as a bundler/optimizer only;
HTMX only as progressive enhancement; modern HTML/CSS before JavaScript; zero JS
in the main flow (a JS-disabled user completes the primary path). Reimplement the
upstream contract; do not simplify the design, drop states, or alter the visual
contract just because you switch from a web-component runtime to server-rendered
HTML.

Do NOT use or introduce: React; Lit; Shadow DOM; Astro; `templ`; JSX; Custom
Elements as a runtime requirement; CDN; unnecessary external dependencies; or
component JavaScript without a demonstrated platform gap.

Do not add JavaScript for convenience or to replace native semantics.

Do not init or use Git: no `git init`, status, diff, branches, commits,
worktrees, stash, reset, checkout, fetch, pull or push.

Do not read, print, request or modify: credentials, tokens, cookies, keys;
`.env`; credential stores; Git/npm authentication.

Read-only references: keep any reference checkout (e.g. Material Web upstream)
read-only; never copy its CSS/code without reviewing provenance and license.
Reimplement contracts; do not copy.

Do not invent compatibility, paths, tests, reviews or results. Use `UNKNOWN`
where appropriate.

## 3. Port and process safety

The accepted app may be running on `:8787` and a built binary may exist
(`cmd/gelium`). Absolute worker rules:

- never start anything on `:8787`;
- never stop, restart or signal the `:8787` process;
- never overwrite, rebuild, run or replace the production binary;
- never kill another process;
- never test changes against `:8787`;
- use `:8788` exclusively for your smoke server, and only if it is free.

Inspect read-only before touching anything:

```bash
SS -ltnp | grep -E ':(8787|8788)[[:space:]]' || true
curl -fsS --max-time 3 http://localhost:8787/healthz || true
```

Do not derive destructive actions from those results. For smoke, start
`PORT=8788 go run ./cmd/gelium` with your tool's background-process capability,
keep its session/process ID, and stop only that process when done. Do not use
`&`, `nohup`, or orphan processes.

If `:8788` is occupied by a process you did not start: do not kill it, do not
reuse it, report `BLOCKED_PORT_8788`, and ask the coordinator.

## 4. Concurrency without Git

### 4.1 Use only the assigned mode

`SHARED_HANDOFF` — edit only new or exclusive files listed in `OWNED PATHS`. Do
not edit shared files. For every shared change, deliver an integration manifest
and textual patch. If integration is pending, end as `READY_FOR_INTEGRATION`,
not `COMPLETE`.

`ISOLATED_PHYSICAL_WORKSPACE` — work only in `AUTHORIZED WORKSPACE`. Do not write
to the canonical checkout. Builds, tests and `:8788` happen in the isolated
copy. Do not auto-copy results to canonical. Deliver a full manifest and
reproducible patches.

`EXCLUSIVE_INTEGRATION` — list each shared file before writing it; verify every
one appears literally in the reservation; capture baseline hashes; a component
reservation does not imply global ownership; if a path's ownership is missing,
do not edit it.

### 4.2 Shared by default

Unless explicitly reserved, do not edit concurrently:

```text
internal/app/server.go
internal/app/routes.go
internal/app/docs.go
internal/app/server_test.go
lib/styles/*.css  (regenerate dist, do not hand-edit the bundle)
lib/dist/gelium.css
lib/package.json
lib/package-lock.json
site/web/static/app.css
site/web/static/app.js
site/web/static/htmx.min.js
site/web/assets.go
README.md
package.json
package-lock.json
go.mod
go.sum
cmd/gelium/*.go
```

Router/mux, central page view, navigation, registry, layout, bundles, indexes and
generated files are shared too.

### 4.3 Ownership and drift

Before writing, present:

| Path | New/existing | Exclusive/shared | Owner | Reservation | Action |
|---|---|---|---|---|---|

For each authorized existing file: record SHA-256 and size; recompute the hash
immediately before writing; recompute again before handoff; compare against your
last known state.

```bash
sha256sum path/to/file
wc -c < path/to/file
```

If a file changed without your own write: stop; do not overwrite or auto-merge;
do not revert others' work; keep your proposal as a handoff; report
`ABORTED_ON_DRIFT` with path, expected and observed hash.

### 4.4 Integration manifest

For each shared change you did not apply, deliver:

```text
FILE: path
BASELINE_SHA256: hash
OWNER_REQUIRED: integration-owner
ANCHOR: stable, unique text
PURPOSE: reason
DEPENDENCIES: related contracts
PATCH:
...diff or exact block...
POSTCONDITION:
...expected result...
TESTS THAT PROVE IT:
...tests/commands...
```

`STALE_REBASE_REQUIRED` is not a final worker status: the coordinator applies it
to a previously emitted handoff whose baseline went stale. If the worker detects
drift during its run, its final status is `ABORTED_ON_DRIFT` and it generates no
patch against the changed baseline.

## 5. Required discovery — do not write yet

1. Read `AGENTS.md`, the skills, `llms-ux.txt`, the component registry row, and
   the vocabulary entry for the patterns used.
2. Inspect existing similar Gelium components (handler, partial, CSS, tests).
3. Inspect Go, templates, tests, CSS, theme, docs, build and relevant asset
   versioning.
4. Report mode, ownership, protected paths and port/process state.
5. Distinguish facts, inferences, `UNKNOWN` and decisions needing approval.

Do not write production code during this phase.

## 6. Upstream evidence

Do not implement from memory, screenshots, or generic summaries. For the
upstream component you are porting/contracting, locate literal paths for:
documentation and demos; render/source; shared primitives; CSS; public tokens;
hard-coded geometry; behavior/accessibility tests. Record exact paths, variants
and states, anatomy, dimensions, spacing, type, colors, motion, keyboard/focus,
ARIA, and proposed divergences. Use `UNKNOWN`/`BLOCKED` and ask the coordinator
for a baseline when a required snapshot cannot be verified. Reference is not a
runtime dependency; reimplement contracts, never blindly copy.

## 7. Mandatory platform-first audit

Before proposing JavaScript, build a table:

| Capability | Native HTML/CSS | Current support | Baseline | No-JS | Real gap | Solution |
|---|---|---|---|---|---|---|

Audit as applicable: HTML elements/attributes; forms and navigation; CSS
selectors/properties; Popover API, Invoker Commands, top layer, inertness and
dismissal; keyboard and focus; anchor positioning; reduced motion and forced
colors; RTL and responsive; server/network semantics; real browser probes.
Sources in order: WHATWG/W3C; MDN; Browser Compat Data; Web Features/Baseline;
real probe. Date and evaluate each capability against the provided
`SUPPORTED BROWSER MATRIX` and `WEB BASELINE SNAPSHOT DATE`. If either is
missing, mark the decision `UNKNOWN`/`BLOCKED` rather than silently picking a
matrix, and never use a historical matrix as if it were current.

JavaScript is accepted only if a functional gap remains that cannot be solved
correctly with (1) semantic HTML, (2) declarative CSS, (3) server-rendered
form/navigation, (4) HTMX as enhancement. If still necessary: document the exact
gap; write a RED test first; use a minimal vanilla/framework-free module; keep a
real no-JS flow; never make JS a requirement to read, navigate, submit or
complete the main flow. For browsers without the modern primitive, use a real
server-rendered route/page. Do not simulate modals, tabs, menus or validation
with inaccessible CSS hacks.

## 8. Specification before code

Produce a matrix:

| Feature | Upstream contract | Gelium strategy | Test | Divergence |
|---|---|---|---|---|

Cover: semantic root and anatomy; variants and states; combinations and
precedence; rest/hover/focus-visible/active; disabled/loading/error/empty where
applicable; keyboard and focus lifecycle; labels, names and descriptions; forms,
values, submission and HTTP codes; no-JS and HTMX; light/dark, narrow/wide, RTL;
reduced motion and forced colors; assets and trust boundaries. Any non-trivial
divergence requires approval BEFORE implementation. Do not call an unverified
approximation "parity".

## 9. Strict TDD

Law: NO PRODUCTION CODE WITHOUT A FAILING TEST FIRST.

Use vertical cycles, one per behavior: write a minimal test; run only that test;
observe and log the expected RED; confirm it is not a typo/setup; implement the
minimum; observe GREEN; run the related suite; refactor only in green; repeat.
Do not write all tests and then all implementation. Keep this log:

| Cycle | Test | RED command | Observed failure | GREEN change | GREEN command | Result |
|---|---|---|---|---|---|---|

If a test passes on the first attempt it does not count as RED — the behavior
already existed, or the test is wrong. Never delete or revert others' work.

Base commands (run from the repo root):

```bash
go test ./internal/app -run 'TestName' -count=1 -v
go test ./lib/... -count=1
go test ./... -count=1
go vet ./...
go mod verify
```

## 10. Templates, attributes and trust

Use `html/template` and concrete Go view models. Required: closed vocabularies
for variant/type/command/placement; explicit defaults; real booleans for boolean
attributes; escapable strings for text/values; individually emitted attributes;
escaping/omission/invalid-input tests. Prohibited: `map[string]any` as API;
attrs/raw HTML strings; dynamic `template.HTMLAttr`; caller-controlled arbitrary
classes; user-controlled tag/attr names; `template.HTML` or `template.URL` over
untrusted content. `template.HTML` only for trusted internal markup, with a trust
boundary comment and tests. Decorative trusted SVG must carry
`aria-hidden="true"` and `focusable="false"`; visible text provides the
accessible name.

## 11. CSS, tokens and accessibility

Native elements before ARIA; never use `div`/`span` to replace controls; no
redundant ARIA. All public tokens use `--ui-*`; Material mappings live in the
theme. Do not add token families the slice does not need. Define explicit state
precedence. Focus must not change geometry or cause layout shift; never remove
outlines without an equivalent indicator. Test: light/dark; narrow/wide; RTL
where applicable; reduced motion; forced colors; disabled and combined states;
text/errors not communicated by color alone; contrast before claiming WCAG.

## 12. Complete flow without JavaScript

The agreed main flow must complete with JS/HTMX disabled within the supported
matrix. Where applicable test: a real `href`; form `method`/`action`; a complete
HTML response; preserved values/errors; the correct HTTP status; a return anchor;
an equivalent server-rendered action. HTMX may receive fragments, but a non-HX
branch must exist. Test normal requests and `HX-Request: true` separately. Do not
claim "no JS" if it only renders, does not complete the action, depends solely on
`hx-*`, or imitates semantics with CSS.

## 13. Dogfooded docs

Create the docs page for the component using the real implementation, never
duplicated markup. In `site/web/content/<slug>.md` document: purpose and
anatomy; variants/states; accessibility; no-JS behavior; HTMX if applicable;
compatibility/Baseline; trust boundary; divergences; visual checklist. In
`SHARED_HANDOFF`, route/nav/layout are delivered as an integration manifest.

## 14. Build and assets

Source CSS lives in `lib/styles/<slug>.css` (and is `@import`ed through the
package manifest). The published bundle `lib/dist/gelium.css` is generated — do
not hand-edit it; regenerate with the package build script and verify the output
matches the source. Docs static assets (`site/web/`) are embedded; any embedded
asset change requires a new versioned URL or content hash — `Cache-Control`
alone is not enough. Test that HTML references the new version, the served CSS
contains the component marker, package/docs/layout do not diverge, the build is
reproducible, there is no CDN, and output matches source. Do not build the
production binary.

## 15. Separate reviews

Gate A — Spec review: after green tests, review only against parameters, matrix,
upstream, no-JS and acceptance. Mark each requirement `PASS | FAIL | BLOCKED |
UNKNOWN` with path/test evidence. Fix FAIL with TDD and repeat to PASS.

Gate B — Quality review: only after Spec PASS, review security, escaping, trust,
semantics, accessibility, CSS, maintainability, ownership/drift, scope, deps,
assets, docs and test fragility. Fix Critical/Important with TDD and repeat to
APPROVED. Use independent reviewers when your platform allows; otherwise declare
the limitation and do not invent independent approval.

## 16. Final verification

If authorized:

```bash
npm run build                 # from site/ (regenerates site/web/static assets)
node scripts/build-lib-dist.mjs   # from repo root (regenerates lib/dist)
go test ./... -count=1
go vet ./...
go mod verify
node --check site/web/static/app.js
```

Do not hide warnings. If the build touches a path without ownership, stop and
report a scope violation.

## 17. Smoke on :8788

Only after build/tests and from the authorized workspace: confirm `:8788` is
free; start your own server in the background; wait for `/healthz`; test route,
assets and version; test the non-HX flow; test HX if applicable; test
status/headers; use a real browser and console; stop only your process.
Validate light/dark, narrow/wide, keyboard, focus, disabled/error, reduced
motion, forced colors where feasible, and visual comparison upstream. Never use
`:8787` to test changes.

## 18. User acceptance

Deliver an observable, specific checklist covering at least: anatomy and
variants; rest/hover/focus/active; disabled/loading/error; keyboard/focus;
no-JS; HTMX as enhancement; light/dark; narrow/wide/RTL; reduced
motion/forced colors; console; real docs; versioned assets; confirmation that
`:8787`/the production binary were not touched. Do not integrate another
component or replace the stable app. Wait for explicit acceptance.

## 19. Allowed final statuses

Use exactly one:

- `COMPLETE_AWAITING_USER_ACCEPTANCE`: implementation and integration complete,
  reviews and smoke approved;
- `READY_FOR_INTEGRATION`: handoff ready, shared patches not yet applied;
- `BLOCKED`: missing evidence, decision, ownership or port;
- `ABORTED_ON_DRIFT`: a baseline changed concurrently.

`SHARED_HANDOFF` cannot declare COMPLETE while patches remain.

## 20. Mandatory delivery

### Status
`[ONE_ALLOWED_STATUS_CHOSEN_BY_WORKER]`

### Summary
Component, scope, mode and result.

### App safety
State of `:8787`, the production binary, `:8788`, and your own process.

### Upstream evidence
Snapshot/manifest provided, paths and divergences.

### Platform-first audit
Date, table and JS gap, or "not required".

### Ownership and drift
Paths, reservations, baseline/final hashes and drift.

### TDD log
Each RED/GREEN with real commands and results.

### Files
Separately list created, modified, generated, proposed-not-applied and
protected-verified files.

### Integration manifest
Patches for shared files not edited.

### Reviews
Spec and Quality with findings and resolution.

### Verification
Build, tests, vet, mod verify, smoke, browser, no-JS and HTMX.

### User checklist
Observable cases.

### Issues/blockers
`UNKNOWN`, risks and pending work.

### Confirmation
Explicitly declare: no Git, no credentials, no CDN, no React/Lit/Shadow
DOM/Astro/templ, no JS without a gap, read-only references, `:8787`/production
binary intact, no concurrent overwrite, no invented results.
````

## Minimal assignment template

```text
COMPONENT: Checkbox
SLUG: checkbox
CATEGORY: CORE
REQUIRED VARIANTS: checked, unchecked, indeterminate if approved
REQUIRED STATES: rest, hover, focus, pressed, disabled, error
REQUIRED BEHAVIORS: clickable label, normal form submit, native keyboard
SERVER/HTMX FLOW, IF ANY: normal submit mandatory; HTMX optional
ACCEPTANCE CRITERIA: complete after upstream/platform-first audit
SUPPORTED BROWSER MATRIX: {{COORDINATOR_MUST_DEFINE}}
WEB BASELINE SNAPSHOT DATE: {{COORDINATOR_MUST_DEFINE}}

CANONICAL REPOSITORY: {{COORDINATOR_MUST_DEFINE}}
ASSIGNED MODE: SHARED_HANDOFF
AUTHORIZED WORKSPACE: {{COORDINATOR_MUST_DEFINE}}
WORKER ID: checkbox-worker-01
RESERVATION ID: checkbox-exclusive-files-only
OWNED PATHS: lib/templates/checkbox.html, lib/styles/checkbox.css, site/web/content/checkbox.md
ADDITIONAL FORBIDDEN PATHS: all shared files; hand off patches
BASELINE SHA-256: hashes provided by integrator
PROPOSED NEW ASSET VERSION: integrator-owned
```

## Coordination note

Until the Gelium roadmap Wave P completes, assigning two components to two AIs
does NOT mean allowing two writers on the canonical checkout. Workers prepare
contracts/artifacts in parallel; a single integrator incorporates each lane
serially.
