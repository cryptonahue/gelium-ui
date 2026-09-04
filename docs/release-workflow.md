# Gelium UI — release workflow

This repository publishes the `lib/` workspace as the `gelium-ui` npm package.
The docs application in `site/` consumes that package through the workspace and
serves as its dogfood consumer. A package release must keep the package, docs
consumer, generated projections, and release notes coherent.

## Authority

- `lib/package.json` — canonical npm package name and version.
- `lib/version.go` — canonical static-asset cache-busting version.
- `lib/llms.txt` — packaged agent-facing package metadata.
- `README.md` — repository release heading and monorepo entrypoint.
- `CHANGELOG.md` — release notes; every package version requires a matching
  `## [version]` entry before a release can proceed.
- `site/package.json` and `package-lock.json` — workspace consumer pin and
  resolved dependency graph.

The docs site is a consumer, not a second package authority. `site/web/static/`
and `lib/dist/` are generated; never hand-edit them.

## Required sequence

1. Choose the next SemVer version from the current `lib/package.json`.
2. Add the release entry to `CHANGELOG.md` before publishing. Include Added,
   Changed, Fixed, or Removed sections as applicable.
3. Update `lib/package.json`, `site/package.json`, `lib/version.go`, and the
   packaged metadata (`lib/SKILL.md`, component skills, and `lib/llms.txt`).
4. Update release-facing documentation that contains the current version.
5. Refresh the workspace lockfile with npm; do not hand-edit integrity data.
6. Run `npm run build` to regenerate projections and embedded assets.
7. Run the release gate and the repository verification baseline:

   ```bash
   npm run release:check
   npm run build
   go test ./...
   go vet ./internal/... ./site/... ./lib/...
   bash scripts/ux-detect.sh
   git diff --check
   ```

8. Inspect the package without publishing:

   ```bash
   npm pack --dry-run -w gelium-ui
   ```

9. Commit the verified release preparation. Push separately. Publish only after
   explicit authorization and successful npm authentication.
10. After publishing, verify the exact registry version with:

   ```bash
   npm view gelium-ui version --registry=https://registry.npmjs.org
   ```

## Release gate

`npm run release:check` verifies that the current package version appears in the
root README, changelog, `lib/version.go`, `lib/llms.txt`, `lib/SKILL.md`, the
component implementation skill, the site consumer manifest, the lockfile, and
the generated `site/web/static/llms.txt` projection.

This gate checks metadata consistency. It does not publish, create a Git tag,
or claim that the registry already contains the version.
