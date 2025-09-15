# Repository Guidelines

## Project Structure & Module Organization
- `core/`: service discovery, naming, exposure modes. (planned)
- `docker/`: Docker API access and event stream. (planned)
- `tailscale/`: certs, MagicDNS checks, Funnel control. (planned)
- `traefik/`: TLS file writer and watchers. (planned)
- `cmd/tailwhale`: Go CLI entrypoint (sync, watch, list).
- `cmd/extension-api`: REST backend for UI. (planned)
- `ui/`: Next.js frontend (scaffolded in this repo).

## Build, Test, and Development Commands
- Node package manager: use `pnpm` (not npm/yarn). Enable via `corepack enable`. Do not commit `package-lock.json` or `yarn.lock`.
- Make targets: `make build`, `make test`, `make dev` — wraps Go/UI tasks and no-ops if a part is missing.
- Go build: `go build ./cmd/tailwhale` — builds the CLI (when available).
- Go run: `go run ./cmd/tailwhale --help` — quick local run (when available).
- Go tests: `go test ./...` — runs unit tests across modules.
- UI dev: `cd ui && pnpm install && pnpm dev` — Next.js dev server.
- UI build: `cd ui && pnpm build` — production build.
- TUI (optional, Bubble Tea): code lives under `cmd/twui` behind tag `tui`.
  - Install dep: `go get github.com/charmbracelet/bubbletea@latest`
  - Build: `go build -tags tui ./cmd/twui` then run `./twui`.

## CLI Usage & Config
- `tailwhale sync`: one-off discover → cert paths → write `traefik/tls.yml`. Flags: `--host`, `--tailnet`, `--tls-path`, `--cert-dir`, `--config <json>`.
- `tailwhale watch`: event-driven (Docker) with ticker fallback; writes `tls.yml` atomically each sync. Same flags as `sync` plus `--interval`.
- `tailwhale list`: `--json` for machine output; `--from-file <containers.json>` for offline dev. See `examples/containers.json`.
- Config file example (JSON): `{ "host": "host1", "tailnet": "tn", "tlsPath": "traefik/tls.yml", "certDir": "/var/lib/tailwhale/certs" }`. Flags override.
  - Sample config at `examples/tailwhale.json`.

## Docker Provider (Build Tags)
- Default build uses a fake provider (no Docker SDK required).
- Real provider behind tag `docker`: `go build -tags docker ./cmd/tailwhale` to enable Docker events-based watch.
 - TUI behind tag `tui`: `go build -tags tui ./cmd/twui` (adds Bubble Tea dep only when building TUI).

## Node Setup (pnpm/Corepack)
- Enable Corepack: `corepack enable`
- Activate pnpm: `corepack prepare pnpm@latest --activate`
- Verify: `pnpm -v`
- First install in `ui/`: `cd ui && pnpm install`
- Fallback (if Corepack unavailable): `npm i -g pnpm` then `pnpm -v`

## Coding Style & Naming Conventions
- Go: formatted by `gofmt`/`goimports`; package names lowercase; exported identifiers use CamelCase; errors with `%w` for wrapping; prefer `context.Context` and structured logs.
- TypeScript/React (ui/): use Biome (`pnpm lint`, `pnpm format`) for lint+format; components PascalCase; files kebab-case; 2-space indent; keep strict `tsconfig`.
- YAML (Traefik): 2-space indent; deterministic key ordering in `traefik/tls.yml`.

## Testing Guidelines
- Go: place tests in `*_test.go`; name `TestXxx`; prefer table-driven tests. Run `go test ./... -race -cover`. Aim for ~80% package-level coverage where practical.
- UI: `cd ui && pnpm typecheck && pnpm lint && pnpm test` (tests stubbed initially). Add screenshots for UI PRs that change visuals.
- Golden tests: see `internal/traefik/tlswriter_golden_test.go` to assert deterministic `tls.yml` output.

## Commit & Pull Request Guidelines
- Commits: Conventional Commits. Examples:
  - `feat(cli): add watch mode for Mode A`
  - `fix(traefik): stable tls.yml sort order`
  - `docs(readme): clarify Funnel mode`
- Types: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`, `revert`, `release`.
- Hooks: Husky `commit-msg` enforces format; `pre-commit` runs `pnpm -C ui format && pnpm -C ui typecheck && pnpm -C ui lint`.
- PRs: must use the template. Required sections: Summary, Linked Issues (use `Closes #ID`), Affected Areas, Test & Checks. Keep descriptions specific and link screenshots when UI changes.

## Agent Rules
- Keep docs synced: when commands, flags, config, CI, or paths change, update `README.md`, `AGENTS.md`, and `examples/` in the same PR.
- Use the PR template: fill all required sections; include `Closes #<issue>` so issues auto-close on merge.
- Required checks: `go`, `ui`, `validate` must pass. Auto‑merge is allowed only after checks succeed.
- Preflight before PR: `go test ./...` and, if UI changed, `pnpm -C ui typecheck && pnpm -C ui lint && pnpm -C ui build`.
- One branch per task; avoid committing to `main`.

## PR Template & Auto‑Merge
- Create with body file:
  - `gh pr create --base main --head <branch> --title "..." --body-file ./PR_BODY.md`
- Minimal PR_BODY.md scaffold:
  - `## Summary` one paragraph
  - `## Linked Issues` include `Closes #ID`
  - `## Changes` bullet list
  - `## Affected Areas` scope
  - `## Test & Checks` list of commands/evidence
- Edit and re‑run template check if failed:
  - `gh pr edit <N> --body-file ./PR_BODY.md`
  - `gh run list --workflow "PR Template Check" --branch <branch> --limit 1 | xargs -n1 -I{} gh run rerun {}`
- If checks show “Expected/Waiting”: branch protection must match check‑run names `go`, `ui`, `validate` (not “CI / go”). Update in repo Settings or via API.

Helpers
- `scripts/pr-body-sample.md`: ready-to-use template matching required sections.
- `bash scripts/pr-update.sh --generate ./PR_BODY.md`: create a scaffold file.
- `bash scripts/pr-update.sh --pr <N> --body ./PR_BODY.md`: update PR body and re-run the template check.
- Make targets:
  - `make pr-template` → generates `PR_BODY.md` from the sample.
  - `make pr-update PR=<N>` → updates PR `<N>` using `PR_BODY.md` and re-runs the template check.

## Branching Workflow
- Stay updated: `git fetch origin && git switch main && git pull --ff-only`.
- Create a branch per task: `git switch -c feat/<scope>-<short-desc>` or `fix/<scope>-<short-desc>` (kebab-case).
- Keep in sync: `git fetch origin && git rebase origin/main` (resolve conflicts, re-run tests).
- Push and open PR: `git push -u origin <branch>` then `gh pr create` (or open on GitHub).
- Do not commit directly to `main`.

Helper script
- `scripts/new-branch.sh -t feat -s cli -d "add watch mode" --push --pr`
- Prompts if flags are omitted; always updates from `origin/main` first.

Rebase helper
- `scripts/rebase-main.sh` — fetches and rebases current branch onto `origin/main`.
- Supports `--continue` and `--abort` during conflict resolution.

Release workflow
- Versioning: use tags like `v0.1.0`. Keep `main` releasable.
- `scripts/release.sh <version> [--push] [--gh-release]` to tag, push, and optionally create a GitHub release with notes.

## Security & Configuration Tips
- Never commit secrets; use `.env` and update `.env.example` when adding vars (e.g., Tailscale auth keys).
- Restrict file writes to intended paths (e.g., `traefik/tls.yml`).
- When using Funnel or MagicDNS, verify config with a dry run before enabling `watch`.

## CI
- Combined CI: `.github/workflows/ci.yml` runs two jobs on PRs:
  - `ui`: pnpm install → typecheck → lint → build (skips if no `ui/`).
  - `go`: `go test ./... -race -cover` (skips if no Go code).
- PR template check: `.github/workflows/pr-template-check.yml` ensures PR body uses the template.
- Branch protection: require status checks with names `go`, `ui`, `validate` and “require branches up to date”.

## Code Owners
- File: `.github/CODEOWNERS` assigns reviewers by path. Default owner is `@frnwtr` for the whole repo and key directories.
- If collaborators/teams are added later, update paths to include their handles.
- In branch protection, enable "Require review from Code Owners" to enforce.

## Architecture Overview
TailWhale syncs Docker → Tailscale certs → Traefik TLS. Choose Mode A (host + Traefik), Mode B (sidecar per container), or Mode C (Funnel on Traefik) based on deployment needs.
