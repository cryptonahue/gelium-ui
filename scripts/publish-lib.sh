#!/usr/bin/env bash
# Publish gelium-ui from lib/ after gates. Requires npm auth (npm login or NPM_TOKEN).
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

echo "== gates =="
go test ./internal/... ./site/... ./lib/...
go vet ./internal/... ./site/... ./lib/...
npm run build
git diff --check

echo "== pack =="
cd "$ROOT/lib"
npm pack --dry-run

if ! npm whoami >/dev/null 2>&1; then
  if [[ -n "${NPM_TOKEN:-}" ]]; then
    echo "//registry.npmjs.org/:_authToken=${NPM_TOKEN}" > "$ROOT/lib/.npmrc"
    trap 'rm -f "$ROOT/lib/.npmrc"' EXIT
  else
    echo "Not logged in to npm. Run: npm login"
    echo "Or set NPM_TOKEN and re-run this script."
    exit 1
  fi
fi

echo "== publish =="
# --access public for unscoped package; dry-run first unless PUBLISH_REAL=1
if [[ "${PUBLISH_REAL:-}" == "1" ]]; then
  npm publish --access public
  echo "Published gelium-ui@$(node -p "require('./package.json').version")"
else
  npm publish --access public --dry-run
  echo "Dry-run OK. Re-run with PUBLISH_REAL=1 to publish."
fi
