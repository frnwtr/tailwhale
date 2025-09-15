#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<EOF
Usage: $0 <version> [--push] [--gh-release]

Tag a new version and optionally push and create a GitHub release.

Args
  <version>     Semver like v0.1.0 or 0.1.0

Options
  --push        Push tag to origin
  --gh-release  Create a GitHub release via gh (requires gh auth)
  -h, --help    Show this help

Examples
  $0 v0.1.0 --push --gh-release
  $0 0.1.0 --push
EOF
}

case "${1:-}" in
  -h|--help|"") usage; exit 0;;
esac

VERSION="$1"; shift || true
PUSH=false
GHREL=false

while [ $# -gt 0 ]; do
  case "$1" in
    --push) PUSH=true; shift;;
    --gh-release) GHREL=true; shift;;
    *) echo "Unknown option: $1" >&2; usage; exit 2;;
  esac
done

# Normalize: ensure leading 'v'
case "$VERSION" in v*) ;; *) VERSION="v$VERSION";; esac

command -v git >/dev/null 2>&1 || { echo "git not found" >&2; exit 1; }
git rev-parse --is-inside-work-tree >/dev/null 2>&1 || { echo "Not a git repo" >&2; exit 1; }

# Ensure up to date main
git fetch origin
git switch main
git pull --ff-only origin main

# Ensure no existing tag
if git rev-parse "$VERSION" >/dev/null 2>&1; then
  echo "Tag $VERSION already exists" >&2
  exit 1
fi

# Create annotated tag
git tag -a "$VERSION" -m "Release $VERSION"
echo "Created tag $VERSION"

if $PUSH; then
  git push origin "$VERSION"
fi

if $GHREL; then
  if command -v gh >/dev/null 2>&1; then
    gh release create "$VERSION" --generate-notes || {
      echo "Failed to create GitHub release via gh" >&2
      exit 1
    }
  else
    echo "gh CLI not found; skipping GitHub release" >&2
  fi
fi

echo "Done."

