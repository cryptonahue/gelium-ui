#!/usr/bin/env bash
# agent-start.sh — optional startup hook for projects using gelium-ui.
# Refreshes Gentle AI's project-local skill registry without selecting a task
# route, enabling SDD, or changing repository delivery state.
set -euo pipefail

project_root="${1:-.}"
if [[ ! -d "$project_root" ]]; then
  echo "agent-start: project root does not exist: $project_root" >&2
  exit 2
fi

if ! command -v gentle-ai >/dev/null 2>&1; then
  echo "agent-start: gentle-ai not installed; registry refresh skipped"
  exit 0
fi

cd "$project_root"
gentle-ai skill-registry refresh --force
echo "agent-start: skill registry refreshed"
