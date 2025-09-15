#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<EOF
Usage:
  $(basename "$0") --pr <number> --body <path>
  $(basename "$0") --generate <path>

Options:
  --pr <number>        Pull request number to update
  --body <path>        Markdown file to set as PR body (must contain required sections)
  --generate <path>    Generate a PR body scaffold at <path> from scripts/pr-body-sample.md
  -h, --help           Show help

Notes:
  - After updating the PR body, this script re-runs the "PR Template Check" on the PR branch.
  - Required sections: Summary, Linked Issues (with Closes #ID), Affected Areas, Test & Checks.
EOF
}

PR=""
BODY=""
GEN=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --pr) PR="$2"; shift 2;;
    --body) BODY="$2"; shift 2;;
    --generate) GEN="$2"; shift 2;;
    -h|--help) usage; exit 0;;
    *) echo "Unknown arg: $1" >&2; usage; exit 2;;
  esac
done

if [[ -n "$GEN" ]]; then
  if [[ -e "$GEN" ]]; then
    echo "Refusing to overwrite existing file: $GEN" >&2
    exit 1
  fi
  mkdir -p "$(dirname "$GEN")"
  cp scripts/pr-body-sample.md "$GEN"
  echo "Generated scaffold: $GEN"
  echo "Edit the file, then run: $(basename "$0") --pr <number> --body $GEN"
  exit 0
fi

if [[ -z "$PR" || -z "$BODY" ]]; then
  echo "Missing --pr or --body" >&2
  usage
  exit 2
fi

if ! command -v gh >/dev/null 2>&1; then
  echo "gh CLI is required" >&2
  exit 1
fi

# Update PR body
gh pr edit "$PR" --body-file "$BODY"

# Find branch and rerun PR Template Check
BR=$(gh pr view "$PR" --json headRefName -q .headRefName)
RID=$(gh run list --branch "$BR" --workflow "PR Template Check" --limit 1 --json databaseId -q '.[0].databaseId' || true)
if [[ -n "$RID" ]]; then
  gh run rerun "$RID" >/dev/null || true
fi

echo "Updated PR #$PR and re-ran PR Template Check (if found)."

