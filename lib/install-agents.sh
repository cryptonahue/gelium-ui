#!/usr/bin/env bash
# install-agents.sh — copy the gelium-ui agent guidance (AGENTS.md, skills/,
# llms-ux.txt) into the local agent's skill directory so an LLM loading this
# package in any project sees how to apply it.
#
# Targets (first writable wins by default; --target overrides):
#   --hermes | --cursor | --claude | <path>
#   Detectable: ~/.hermes/skills, ~/.cursor/skills, ~/.claude/skills
#   --dry-run  : show what would be copied, copy nothing.
set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SRC="$HERE"
NAME="gelium-ui"
FILES=("AGENTS.md" "SKILLS.md" "llms-ux.txt" "llms.txt" "skills")

target=""
dry=false
while [[ $# -gt 0 ]]; do
  case "$1" in
    --hermes)  target="$HOME/.hermes/skills/$NAME" ;;
    --cursor)  target="$HOME/.cursor/skills/$NAME" ;;
    --claude)  target="$HOME/.claude/skills/$NAME" ;;
    --dry-run) dry=true ;;
    *)
      if [[ "$1" == -* ]]; then
        echo "unknown option: $1" >&2; echo "usage: $0 [--hermes|--cursor|--claude|<path>] [--dry-run]" >&2; exit 2
      fi
      target="${1%/}" ;;
  esac
  shift
done

# Default detection: first existing agent-skill root.
if [[ -z "$target" ]]; then
  for root in "$HOME/.hermes/skills" "$HOME/.cursor/skills" "$HOME/.claude/skills"; do
    if [[ -d "$root" ]]; then target="$root/$NAME"; break; fi
  done
fi
if [[ -z "$target" ]]; then
  echo "no agent skill root detected; pass a target path explicitly" >&2; exit 1
fi

echo "source : $SRC"
echo "target : $target"
if $dry; then
  echo "(dry-run) would copy: ${FILES[*]}"
  exit 0
fi

mkdir -p "$target/skills"
for f in "${FILES[@]}"; do
  if [[ -e "$SRC/$f" ]]; then
    if [[ "$f" == "skills" ]]; then
      cp -R "$SRC/skills/." "$target/skills/"
      echo "copied skills/ -> $target/skills/ ($(ls "$SRC/skills" | wc -l) files)"
    else
      cp "$SRC/$f" "$target/$f"
      echo "copied $f"
    fi
  else
    echo "skip (missing source): $f" >&2
  fi
done

echo
echo "Installed. LLM tools that load $NAME skills will now see Gelium's good-practice guide."
