#!/usr/bin/env bash
# Publish gelium-ui from lib/ after gates. Requires npm auth (npm login or NPM_TOKEN).
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

write_npmrc() {
  local dir="$1"
  # Granular tokens need the registry host line + always-auth.
  cat > "${dir}/.npmrc" <<EOF
//registry.npmjs.org/:_authToken=${NPM_TOKEN}
registry=https://registry.npmjs.org/
always-auth=true
EOF
  chmod 600 "${dir}/.npmrc"
}

cleanup_npmrc() {
  rm -f "${ROOT}/lib/.npmrc" "${ROOT}/.npmrc"
}

if npm whoami --registry=https://registry.npmjs.org/ >/dev/null 2>&1; then
  echo "== auth: already logged in as $(npm whoami --registry=https://registry.npmjs.org/) =="
elif [[ -n "${NPM_TOKEN:-}" ]]; then
  echo "== auth: using NPM_TOKEN =="
  write_npmrc "${ROOT}/lib"
  write_npmrc "${ROOT}"
  trap cleanup_npmrc EXIT
  if ! npm whoami --registry=https://registry.npmjs.org/ >/dev/null 2>&1; then
    echo "NPM_TOKEN is set but npm whoami still fails."
    echo "Check: token has Read and write on packages; IP allowlist uses 82.197.66.69/32"
    exit 1
  fi
  echo "authenticated as $(npm whoami --registry=https://registry.npmjs.org/)"
else
  echo "Not logged in to npm and NPM_TOKEN is empty."
  echo ""
  echo "In THIS shell run:"
  echo "  export NPM_TOKEN='npm_...your_token...'"
  echo "  [ -n \"\$NPM_TOKEN\" ] && echo token_ok"
  echo "  PUBLISH_REAL=1 bash scripts/publish-lib.sh"
  exit 1
fi

echo "== gates =="
go test ./internal/... ./site/... ./lib/...
go vet ./internal/... ./site/... ./lib/...
npm run build
git diff --check

echo "== pack =="
cd "${ROOT}/lib"
npm pack --dry-run

echo "== publish =="
if [[ "${PUBLISH_REAL:-}" == "1" ]]; then
  npm publish --access public --registry=https://registry.npmjs.org/
  echo "Published gelium-ui@$(node -p "require('./package.json').version")"
else
  npm publish --access public --registry=https://registry.npmjs.org/ --dry-run
  echo "Dry-run OK. Re-run with PUBLISH_REAL=1 to publish."
fi
