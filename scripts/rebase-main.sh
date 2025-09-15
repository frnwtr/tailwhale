#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<EOF
Usage: $0 [--continue|--abort]

Safely rebase the current branch onto origin/main.

Steps
1) Ensure clean working tree
2) Fetch origin
3) Rebase current branch on origin/main

Options
  --continue   Continue an in-progress rebase
  --abort      Abort the current rebase
EOF
}

case "${1:-}" in
  -h|--help) usage; exit 0;;
  --continue) git rebase --continue; exit $?;;
  --abort) git rebase --abort; exit $?;;
esac

# Ensure git
command -v git >/dev/null 2>&1 || { echo "git not found" >&2; exit 1; }

# Ensure inside repo
git rev-parse --is-inside-work-tree >/dev/null 2>&1 || { echo "Not a git repo" >&2; exit 1; }

# Ensure not on main
CURRENT=$(git branch --show-current)
if [ "$CURRENT" = "main" ]; then
  echo "You are on 'main'. Create/checkout a feature branch first." >&2
  exit 1
fi

# Ensure clean
if [ -n "$(git status --porcelain)" ]; then
  echo "Working tree not clean. Commit or stash before rebasing." >&2
  exit 1
fi

git fetch origin
git rebase origin/main
echo "Rebased '$CURRENT' onto origin/main successfully."

