#!/usr/bin/env bash
# install-agents.sh — install the gelium-ui agent skill for a project or agent host.
#
# Project-local (default):
#   bash node_modules/gelium-ui/install-agents.sh
#   bash node_modules/gelium-ui/install-agents.sh --project .
#
# Agent-host targets:
#   --hermes | --cursor | --claude | --codex | --opencode | --gemini |
#   --copilot | --qwen | --global
#   --target <path> | <path>
#   --dry-run
set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SRC="$HERE"
NAME="gelium-ui"
FILES=("AGENTS.md" "SKILL.md" "SKILLS.md" "llms-ux.txt" "llms.txt" "skills")

target=""
dry=false
mode="project-local"

require_value() {
  if [[ $# -lt 2 || -z "$2" ]]; then
    echo "missing value for $1" >&2
    exit 2
  fi
}

while [[ $# -gt 0 ]]; do
  case "$1" in
  --project)
    require_value "$1" "${2:-}"
    project_root="$(cd "$2" && pwd)"
    target="$project_root/.agents/skills/$NAME"
    mode="project-local"
    shift 2
    ;;
  --target)
    require_value "$1" "${2:-}"
    target="${2%/}"
    mode="explicit"
    shift 2
    ;;
  --hermes)
    target="$HOME/.hermes/skills/$NAME"
    mode="global"
    shift
    ;;
  --cursor)
    target="$HOME/.cursor/skills/$NAME"
    mode="global"
    shift
    ;;
  --claude)
    target="$HOME/.claude/skills/$NAME"
    mode="global"
    shift
    ;;
  --codex)
    target="$HOME/.codex/skills/$NAME"
    mode="global"
    shift
    ;;
  --opencode)
    target="$HOME/.config/opencode/skills/$NAME"
    mode="global"
    shift
    ;;
  --gemini)
    target="$HOME/.gemini/skills/$NAME"
    mode="global"
    shift
    ;;
  --copilot)
    target="$HOME/.copilot/skills/$NAME"
    mode="global"
    shift
    ;;
  --qwen)
    target="$HOME/.qwen/skills/$NAME"
    mode="global"
    shift
    ;;
  --global)
    mode="global"
    shift
    ;;
  --dry-run)
    dry=true
    shift
    ;;
  -*)
    echo "unknown option: $1" >&2
    echo "usage: $0 [--project <path>|--target <path>|--hermes|--cursor|--claude|--codex|--opencode|--gemini|--copilot|--qwen|--global|<path>] [--dry-run]" >&2
    exit 2
    ;;
  *)
    if [[ -n "$target" ]]; then
      echo "target already selected; unexpected argument: $1" >&2
      exit 2
    fi
    target="${1%/}"
    mode="explicit"
    shift
    ;;
  esac
done

# Project-local installation is the safe default: it travels with the consumer
# project and does not depend on a particular global agent directory.
if [[ -z "$target" && "$mode" == "project-local" ]]; then
  target="$PWD/.agents/skills/$NAME"
fi

# Explicit global installation keeps the previous host-detection workflow.
if [[ -z "$target" && "$mode" == "global" ]]; then
  for root in \
    "$HOME/.agents/skills" \
    "$HOME/.hermes/skills" \
    "$HOME/.cursor/skills" \
    "$HOME/.claude/skills" \
    "$HOME/.codex/skills" \
    "$HOME/.config/opencode/skills" \
    "$HOME/.gemini/skills" \
    "$HOME/.copilot/skills" \
    "$HOME/.qwen/skills"; do
    if [[ -d "$root" ]]; then
      target="$root/$NAME"
      break
    fi
  done
fi

if [[ -z "$target" ]]; then
  echo "no agent skill root detected; pass --project, --target, or a host option" >&2
  exit 1
fi

echo "source : $SRC"
echo "target : $target"
echo "mode   : $mode"
if $dry; then
  echo "(dry-run) would copy: ${FILES[*]}"
  exit 0
fi

mkdir -p "$target/skills"
for f in "${FILES[@]}"; do
  if [[ -e "$SRC/$f" ]]; then
    if [[ "$f" == "skills" ]]; then
      cp -R "$SRC/skills/." "$target/skills/"
      echo "copied skills/ -> $target/skills/ ($(find "$SRC/skills" -maxdepth 1 -type f | wc -l) files)"
    else
      cp "$SRC/$f" "$target/$f"
      echo "copied $f"
    fi
  else
    echo "skip (missing source): $f" >&2
  fi
done

if [[ ! -f "$target/SKILL.md" || ! -f "$target/skills/00-agent-routing.md" ]]; then
  echo "canonical skill or routing missing after installation: $target" >&2
  exit 1
fi

echo "canonical routing: installed"
if [[ "$mode" == "project-local" ]]; then
  echo "Project-local skill installed. Commit $target if the project shares agent guidance."
else
  echo "Global skill installed."
fi
