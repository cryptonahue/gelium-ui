# Accordion Phase 2 Cascade Implementation Plan

> **For Hermes:** Execute the focused Accordion migration with tests before CSS changes.

**Goal:** Make Accordion visual ownership compose as core → behavior policy → reference preset → product skin → site while retaining native semantics and legacy theme URLs.

**Architecture:** Keep `lib/styles/accordion.css` a neutral token consumer with no behavior-class visual selectors. Add reference- and skin-layer document-root token adapters in the site bundle. Keep `theme-*` sheets as lower-layer legacy adapters; explicit `data-gelium-skin` token adapters win deterministically.

**Tech Stack:** Go tests, embedded CSS assets, Tailwind CSS build, native `<details>/<summary>`.

---

### Task 1: Specify the cascade regression tests

**Files:** `lib/styles_accordion_test.go`, `site/web/styles_theme_mechanism_test.go`, `internal/app/document_selection_test.go`

1. Add tests for declared compiled layer order, selectors in both visual layers, skin-over-reference token ordering, behavior CSS absence of behavior selectors, and exact combo root attributes.
2. Run the narrow Go tests and confirm they fail because the Phase 1 no-op adapter has no reference selectors and behavior CSS owns Basecoat visual branches.

### Task 2: Convert Accordion core to token consumption

**Files:** `lib/styles/accordion.css`, `lib/styles/tokens.css`

1. Preserve `<details>/<summary>` semantics, forced-colors, reduced-motion, and a `max(--ui-touch-target, requested)` 44px floor.
2. Replace every visual behavior branch and behavior-based icon selection with neutral token consumption.

### Task 3: Add scoped reference and skin adapters

**Files:** `site/web/styles/accordion-reference.css`, `site/web/styles/accordion-skin.css`, `site/web/styles/app.css`

1. Import token-only adapters into `gelium.reference` and `gelium.skin` respectively.
2. Add material, audited Basecoat/Vega (with 44px accessibility adaptation), and Gelium Base UI reference selectors.
3. Add material/basecoat/baseui/alden/linear/vercel skin selectors. Ensure `reference=material&skin=basecoat` resolves Basecoat accordion tokens from the later skin layer.

### Task 4: Document and verify

**Files:** `docs/gelium-ui-behavior-contract.md`, `docs/audit/basecoat-official/gelium-translation-audit.md`

1. State exact precedence, legacy migration behavior, Base UI disclaimer, and Basecoat touch-floor adaptation.
2. Run `gofmt`, `npm run build`, `go test ./...`, `go vet ./...`; build a fresh binary and smoke required URLs plus legacy and native/no-hx paths.
