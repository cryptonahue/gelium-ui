# Contributing

Gelium UI is open-code: every partial the docs render is copyable, and contributions are welcome. This page describes how the project is developed so a contribution lands cleanly.

## Where things live

- **Repository**: [github.com/cryptonahue/gelium-ui](https://github.com/cryptonahue/gelium-ui)
- **Issues and pull requests**: GitHub, on the repository above.
- **License**: MIT — see the `LICENSE` file at the repository root.

## Development setup

Requirements are Go 1.24+ and Node.js 20+ with npm. From the repository root:

```bash
npm install
go mod download
npm run build
```

`npm run build` (root) chains the site build — compiles `site/web/styles/app.css` with Tailwind CSS 4 into `site/web/static/app.css`, copies `htmx.min.js`, regenerates `internal/app/icons.go` and copies `lib/js/gelium.js` into `site/web/static/` — then builds the library dist (`lib/styles/dist-entry.css` → `lib/dist/gelium.css`). The final files in `site/web/static/` and `lib/dist/` are build artifacts and are embedded in the Go binary.

## Running the checks

Every contribution must pass the same gates the project enforces. NEVER use `go test ./...` — the tree contains pre-split entrypoints the gates exclude on purpose:

```bash
go test ./internal/... ./site/... ./lib/...
go vet ./internal/... ./site/... ./lib/...
npm run build
git diff --check
gofmt -l internal site lib
```

The contract tests in `internal/app/*_test.go`, `site/web/styles_*_test.go` and `lib/*_test.go` are not optional decoration: they pin observable behavior (semantics, aria attributes, tokens, copy, server contracts, the lib/site template boundary and the package manifest). Do not weaken a test to make it pass — if a contract changed legitimately, rescope the test and say so in the pull request.

## Commit conventions

- Conventional commits (`feat:`, `fix:`, `docs:`, `test:`, `style:`, `refactor:`, `release:`).
- No AI attribution or "Co-Authored-By" trailers.

## Pull request workflow

1. Branch from `main` with a descriptive name (for example `feat/on-this-page`).
2. Make focused commits in reviewable work units.
3. Run the gates above before opening the pull request.
4. Describe the change, the contract impact (if any), and the test evidence in the PR body.

## What to work on

The [Roadmap](/docs/roadmap) tracks where the project is going. The [Information architecture](/docs/information-architecture) handbook page documents the rules for adding docs pages, and [Choose the right control](/docs/choose-the-right-control) governs component decisions — read them before adding either.
