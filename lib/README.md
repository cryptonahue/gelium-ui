# gelium-ui

Server-rendered UI components for **Tailwind CSS 4** and **HTMX**. HTML-first, 0-JS contract by default; progressive enhancement for toast, validation swap, and slider fill.

## Install

```bash
npm install gelium-ui
```

## CSS (drop-in bundle)

```html
<link rel="stylesheet" href="node_modules/gelium-ui/dist/gelium.css" />
```

Or from a bundler / CSS pipeline:

```css
@import "gelium-ui/dist/gelium.css";
```

The prebuilt bundle includes Tailwind preflight, **theme-material**, **theme-basecoat**, tokens, and component styles.

### Source styles (your own Tailwind entry)

```css
@import "tailwindcss";
@import "gelium-ui/styles/index.css";
@import "gelium-ui/themes/theme-material.css";
@import "gelium-ui/themes/theme-basecoat.css";
```

Pick themes you need; both ship with the package.

### Theme selection

Put the direction class on the document root:

```html
<html class="theme-material">
```

Dark scheme (class route, no media-only dark):

```html
<html class="theme-material theme-dark" data-theme="dark">
```

## Consumer JS (optional)

```html
<script defer src="node_modules/gelium-ui/js/gelium.js"></script>
```

Provides:

- HTMX 422 swap when `X-Gelium-Validation: true`
- `gelium:toast` region helpers
- Same-document view transitions guard (`prefers-reduced-motion`)
- Native range slider fill (`--ui-slider-fill`)

## Templates

Copyable HTML partials live under `templates/` (e.g. `templates/button.html`). Wire them with your server or paste into pages.

## What this package is / is not

| Is | Is not |
|----|--------|
| CSS + HTML open-code components | A React/Vue runtime |
| Token + class themes | A full app framework |
| Optional small JS enhancements | Required hydration |

## Agent guidance

This package teaches LLM tools how to apply Gelium's good practices. When you
write UI with an AI agent, start with the canonical routing layer:

- `SKILL.md` — discoverable wrapper for agent registries such as Gentle AI.
- `skills/00-agent-routing.md` — outcome-first route selection and delegation
  boundary.
- `AGENTS.md` — Gelium entry point: golden rules + how to read the guidance.
- `skills/` — actionable decision skills for foundations, screens, forms,
  states, server contracts, a11y, architecture, criteria, approval, and
  references.
- `llms-ux.txt` — the compact decision-id pack (SURFACE / SCREEN / WF / DATA /
  FEED / JOURNEY / MEDIA / SKEL) for fast agent lookup.

Install into your agent's skill directory so the LLM loads it in any project:

```bash
bash node_modules/gelium-ui/install-agents.sh   # auto-detects hermes/cursor/claude
```

For projects using Gentle AI, refresh its project-local registry after
installation:

```bash
./node_modules/gelium-ui/scripts/agent-start.sh .
```

Then select a task route explicitly when implementation begins:

```bash
go run ./cmd/gelium-preflight route --route delegated-direct --format json
```

The startup hook refreshes discovery only; it does not infer intent, enable
SDD/RDD, or authorize delivery. See `SKILLS.md` for the full skill index.

## Monorepo

This package is `lib/` in [gelium-ui](https://github.com/cryptonahue/gelium-ui). The docs site dogfoods it via npm workspaces.

## License

MIT
