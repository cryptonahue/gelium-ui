# Where Gelium UI is going

The roadmap is a contract, not a wish list. Every shipped phase carries tests that pin its behavior, so checkmarks mean verified, not aspirational. The canonical table lives on the [Roadmap page](/docs/roadmap); this post tells the story behind it.

## What shipped: phases A–J

Gelium UI shipped its full ten-phase roadmap in one release line. Phase A built the core: a Go module, server-rendered components, and the zero-JS contract. Phase B defined the `--ui-*` token vocabulary that every theme shares. Phase C made verification theme-agnostic: adding a theme is a glob, not a test rewrite. Phase D covered screen states — empty, error, loading, success, and toast feedback. Phase E made SEO server-driven: dates, JSON-LD, and a configurable base URL. Phase F shipped public content patterns, and Phase G proved them in three full screen recipes. Phase H delivered the single class-route theme mechanism with dark via one route. Phase I added the Basecoat theme, light and dark in the same bundle. Phase J closed the loop with registries and sync guards, so the docs cannot drift from the code.

## Docs and developer experience

The docs shell landed with a Handbook before the component reference. Search, theme select, light/dark switch, breadcrumbs, and a GitHub link came in the same chrome. Every component page gained guidance sections: when to use, when not to use, usability, and accessibility. The content style guide set one voice: plain English, active verbs, and sentences under 25 words. Readability is contract-tested too: a 65-character measure, text-wrap, hyphenation, and WCAG AA contrast in both themes. Acknowledgments credit every inspiration, from Material 3 to Basecoat.

## What is next

The next list is ordered by value, not by date. Truth sync keeps the README, the roadmap, the theme registry, and the CLI in lockstep. Demos become first-class on component pages, so each one shows a live example. SEO productization waits for a public domain: a real BASE_URL and a real social image. Theme polish continues with scoped ownership of token families. A release bumps the version past 0.4.0 and tags it. Optional expansion adds more screen recipes, a third theme, and a runtime registry.

## How to read the checkmarks

Checkmarks mean contract-tested, not aspirational. Each shipped phase has tests that pin the behavior. If a phase regresses, the build fails. That is the honest bar this project holds itself to.
