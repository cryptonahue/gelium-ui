# AGENTS.md — gelium-ui monorepo

This is the repository entrypoint for agents working from the monorepo root.
Gelium UI is a server-rendered Go + HTML/template + Tailwind CSS 4 + HTMX
system, not a client-side framework.

## First route the outcome

Read `lib/SKILL.md` first. It points to the canonical outcome-first routing
contract in `lib/skills/00-agent-routing.md`:

- `direct-exempt` — bounded, already-understood correction;
- `delegated-direct` — broad context or multi-file work without architecture
  change;
- `design-gated` — new screen, flow, or substantial architecture change;
- `escalate` — unresolved product, permission, scope, data, or risk decision;
- `full-sdd` — durable OpenSpec/SDD artifacts materially reduce ambiguity.

Do not apply design-gated ceremony to direct-exempt work. Do not infer the route
from file count alone. When a machine-readable handoff is useful:

```bash
go run ./cmd/gelium-preflight route --route delegated-direct --format json
```

## Repository map

- `lib/` — primary product and published `gelium-ui` npm package. Read `lib/AGENTS.md` and the
  applicable `lib/skills/*.md` before changing package behavior.
- `site/` — docs/dogfood consumer application; it does not define package
  component contracts or package versions.
- `internal/` — Go application, gates, handlers, and contract tests.
- `cmd/` — repository tools and smoke binaries.
- `docs/` — product/system contracts, OpenSpec changes, audits, and diagrams.
- `scripts/` — reproducible build, asset-copy, detector, and packaging scripts.
- `docs/release-workflow.md` — single release/version/changelog workflow.

## Change boundaries

- Preserve unrelated dirty worktree changes. Do not touch `.hermes/`,
  `.worktrees/`, local package archives, credentials, `.env`, or auth stores.
- Do not hand-edit generated assets. Run `npm run build` when source or embedded
  assets change.
- Keep public package source in `lib/`; keep site-only dogfood behavior in
  `site/`.
- Treat `lib/package.json` as the npm version authority. Every package version
  must have a matching `CHANGELOG.md` entry; run `npm run release:check` before
  commit or publish.
- Component implementation uses `lib/skills/14-component-implementation.md`.
  Do not use the removed root-level prompt; it was replaced by this skill.
- Do not commit, push, publish, or deploy unless the user explicitly requests
  that delivery action.

## Verification baseline

Run the narrowest checks first, then the full baseline when the change spans
packages, docs, workflow, or generated assets:

```bash
go test ./...
go vet ./internal/... ./site/... ./lib/...
npm run build
bash scripts/ux-detect.sh
git diff --check
```

Evidence and review do not authorize repository delivery. Commit, push, PR,
publish, and deploy remain separate decisions under ordinary repository policy.
