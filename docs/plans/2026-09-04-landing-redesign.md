# Landing redesign — implementation record

**Route:** `/`
**Route:** `design-gated`
**Brief:** Developers evaluating and installing `gelium-ui`.
**Primary action:** Install the package and reach the first component.
**Approved packet:** A + C1 + C2 + C3 + C4; mobile drawer from `inline-start`.
**Approval:** User approved `ok lo que recomiende gelium` in chat on 2026-09-04.

## Architecture

- Preserve the shared marketing shell on desktop with inline primary links.
- On narrow screens, replace the wrapped inline nav with one native responsive
  disclosure/drawer from `inline-start`; do not add a second body navigation.
- Reuse registered `ui-footer` with real grouped links and native disclosures.
- Keep the main flow server-rendered and usable with JavaScript disabled.
- Make the hero CTA point to the install/first-component section; keep exactly
  one primary button variant on the page.
- Keep the docs consumer and WhatsApp demo as secondary evidence, not the main
  product promise.

## Evidence and references

- `REF-HERO`: promise, short context, one primary CTA, optional supporting proof;
  no copied brand or fake metrics.
- `REF-SHELL`: preserve one shared shell, use real routes, responsive collapse,
  native links, accessible names, no duplicate chrome.
- Existing `lib/templates/footer.html`, `defaultFooter()`, and footer contract
  tests: groups already collapse with `<details>/<summary>` on narrow screens.
- Existing `lib/templates/navigation-drawer.html`: registered drawer anatomy;
  implementation must preserve no-JS and logical inline-start behavior.

## Scope

- Modify landing view data, landing template/styles, marketing shell template/
  styles, and focused contract tests.
- Do not redesign the docs shell.
- Do not add invented metrics, avatars, logos, screenshots, or routes.

## Verification

- Focused landing/shell/footer tests fail before implementation and pass after.
- `npm run build`
- `npm run release:check`
- `go test ./...`
- `go vet ./internal/... ./site/... ./lib/...`
- `bash scripts/ux-detect.sh`
- `git diff --check`
- Render wide/narrow, light/dark, keyboard, no-JS, open/closed drawer, and
  footer disclosure states before release.
