#!/usr/bin/env bash
set -euo pipefail

ALLOWED_TYPES="feat fix docs chore ci refactor perf test release"

usage() {
  cat <<EOF
Usage: $0 [-t type] [-s scope] [-d desc] [--push] [--pr]

Creates a new branch from updated origin/main following our conventions.

Options:
  -t, --type   One of: ${ALLOWED_TYPES}
  -s, --scope  Short scope (e.g., cli, ui, traefik). Default: repo
  -d, --desc   Short description (will be slugified to kebab-case)
      --push   Push branch to origin and set upstream
      --pr     Open a PR via GitHub CLI (gh)
  -h, --help   Show help

Examples:
  $0 -t feat -s cli -d "add watch mode" --push --pr
  $0 --type fix --scope ui --desc "button alignment"
EOF
}

slugify() {
  # lower, replace non-alnum with '-', collapse repeats, trim '-'
  local s="$*"
  s="${s,,}"
  s="$(printf '%s' "$s" | sed -E 's/[^a-z0-9]+/-/g; s/-+/-/g; s/^-+|-+$//g')"
  printf '%s' "$s"
}

has() { command -v "$1" >/dev/null 2>&1; }

TYPE=""
SCOPE="repo"
DESC=""
PUSH=false
PR=false

# Parse args (supports short and simple long flags)
while [ $# -gt 0 ]; do
  case "$1" in
    -t|--type) TYPE=${2:-}; shift 2;;
    -s|--scope) SCOPE=${2:-}; shift 2;;
    -d|--desc) DESC=${2:-}; shift 2;;
    --push) PUSH=true; shift;;
    --pr) PR=true; shift;;
    -h|--help) usage; exit 0;;
    *) echo "Unknown arg: $1" >&2; usage; exit 2;;
  esac
done

if ! has git; then echo "git not found" >&2; exit 1; fi

# Validate inputs; prompt if missing
if [ -z "$TYPE" ]; then read -rp "Type (${ALLOWED_TYPES}): " TYPE; fi
if [ -z "$SCOPE" ]; then read -rp "Scope (default repo): " SCOPE; SCOPE=${SCOPE:-repo}; fi
if [ -z "$DESC" ]; then read -rp "Short description: " DESC; fi

# Validate type
case " ${ALLOWED_TYPES} " in
  *" ${TYPE} "*) ;;
  *) echo "Invalid type: ${TYPE}. Allowed: ${ALLOWED_TYPES}" >&2; exit 2;;
esac

# Ensure in a git repo
git rev-parse --is-inside-work-tree >/dev/null 2>&1 || { echo "Not in a git repo" >&2; exit 1; }

# Ensure clean working tree
if [ -n "$(git status --porcelain)" ]; then
  echo "Working tree not clean. Commit or stash changes before branching." >&2
  exit 1
fi

# Update from origin/main
git fetch origin
git switch main
git pull --ff-only --no-rebase origin main

SLUG=$(slugify "$DESC")
SCOPE_SLUG=$(slugify "$SCOPE")
BRANCH="${TYPE}/${SCOPE_SLUG}-${SLUG}"

# Create branch
git switch -c "$BRANCH"
echo "Created branch: $BRANCH"

if $PUSH; then
  git push -u origin "$BRANCH"
fi

if $PR; then
  if has gh; then
    TITLE="${TYPE}(${SCOPE_SLUG}): ${DESC}"
    gh pr create --base main --head "$BRANCH" --title "$TITLE" --body "New branch created via helper script."
  else
    echo "gh CLI not found; skipping PR creation" >&2
  fi
fi

echo "Done."

