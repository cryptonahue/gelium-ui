# Advisory dogfood audit — recipe public feed

**Candidate:** local Gelium documentation server at `http://127.0.0.1:8788/recipes/public-feed`  
**Scope:** existing `recipe-public-feed` equivalent consumer; no product implementation was changed.  
**Recorded:** 2026-08-30

## Prebuild

```text
$ go run ./cmd/gelium-preflight prebuild --ledger .../recipe-public-feed-ledger.json \
    --changed site/web/templates/recipe-public-feed.html \
    --changed site/web/styles/recipe-public-feed.css \
    --changed internal/app/recipe_public_feed.go --format json
{"status":"pass","issues":[]}
```

## Rendered/server evidence

The local GET candidate returned `HTTP/1.1 200 OK` and a 11,953-byte HTML response. Its rendered HTML contained these contract markers:

- `class="theme-material"` on `html`;
- `recipe-pf-title`, `ui-tabs`, `feed-panel`, `recipe-pf-loading-heading`, `recipe-pf-refresh`, and `gelium-toast-region`;
- native reaction and refresh POST forms.

Server/no-JS interaction checks passed:

```text
reaction-post-303: PASS status=303
refresh-post-nojs: PASS status=200
refresh-post-htmx: PASS status=200
```

The focused route tests passed:

```text
$ go test ./internal/app -run '^TestRecipePublicFeed' -count=1
PASS
```

## Detector evidence

```text
$ go run ./cmd/gelium-ux-detect \
    --owned site/web/templates/recipe-public-feed.html \
    --owned site/web/styles/recipe-public-feed.css \
    --owned internal/app/recipe_public_feed.go \
    --exceptions .../recipe-public-feed-exceptions.json --format json \
    site/web/templates/recipe-public-feed.html \
    site/web/styles/recipe-public-feed.css \
    internal/app/recipe_public_feed.go
{"status":"pass-with-exceptions","findings":[
  {"id":"form-validation","rule":"form-contract","path":"site/web/templates/recipe-public-feed.html","attribution":"owned","exception_id":"recipe-action-forms"},
  {"id":"custom-shell","rule":"shell-contract","path":"site/web/styles/recipe-public-feed.css","attribution":"owned","exception_id":"recipe-layout-shell"}
]}
```

Both findings remain visible and expire on `2026-09-30T00:00:00Z`:

1. `recipe-action-forms`: the recipe uses action-only POST forms, not user-input validation forms.
2. `recipe-layout-shell`: the existing recipe composes documented `ui-*` primitives with scoped `recipe-pf-*` layout CSS.

## Detector calibration finding

The first dogfood scan incorrectly treated the fragment selector `#feed-panel` as a hex color literal. The regression test `TestRunDoesNotTreatIDSelectorsAsColorLiterals` was added before the detector fix. The detector now excludes identifier continuations (`-` and `_`) from its color token match. The repeat scan above has no `color-literal` finding.

This was a real advisory false positive. It was fixed in the detector rather than hidden with an exception.

## Authority evidence

```text
$ go test ./internal/gates -run '^TestCheckAuthorityMatrixReadsCurrentRepositoryMatrix$' -count=1
PASS
```

The matrix compares declared release/version and runtime wire surfaces only; it does not rewrite any contract or scan historical references.

## Limitations and disposition

**Rendered audit status: `pass-with-escalation`.** Local HTTP and rendered-HTML evidence cover the default `theme-material` candidate, server/no-JS mutation paths, and structural markers. A browser-rendered pixel inspection at both wide and narrow viewports, an explicit dark-theme render, keyboard/focus traversal, and production/authenticated content were not observed. The browser harness was unavailable because no supported Chromium browser was running.

This audit is advisory only. It does not authorize a commit, package publication, deployment, or a claim of unconditional visual approval.

## Rollout recommendation

Keep the workflow advisory for the next design-gated consumer scope. Preserve `direct-exempt` work and the legacy detector. Before required mode, validate the two observed heuristic boundaries on another consumer:

- action-only native POST forms without validation fields;
- recipe-scoped layout classes composed around `ui-*` primitives.

No broader UI task, navigation change, or follow-up implementation is created by this record.
