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

## Commit & Pull Request Guidelines
- Commits: Conventional Commits style. Examples:
  - `feat(cli): add watch mode for Mode A`
  - `fix(traefik): stable tls.yml sort order`
  - `docs(readme): clarify Funnel mode`
- Types: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`, `revert`, `release`. Format: `type(scope)?: subject`.
- Hook enforcement: Husky `commit-msg` validates Conventional Commits; `pre-commit` runs `pnpm -C ui typecheck && pnpm -C ui lint`. Run `pnpm install` in `ui/` so `prepare` installs hooks.
- PRs: include description, linked issues, affected exposure modes (A/B/C), test evidence (`go test` output or UI screenshots), and upgrade notes if user-facing behavior changes. If UI deps change, commit `pnpm-lock.yaml` updates. Use the template in `.github/pull_request_template.md`.

## Branching Workflow
- Stay updated: `git fetch origin && git switch main && git pull --ff-only`.
- Create a branch per task: `git switch -c feat/<scope>-<short-desc>` or `fix/<scope>-<short-desc>` (kebab-case).
- Keep in sync: `git fetch origin && git rebase origin/main` (resolve conflicts, re-run tests).
- Push and open PR: `git push -u origin <branch>` then `gh pr create` (or open on GitHub).
- Do not commit directly to `main`.

Helper script
- `scripts/new-branch.sh -t feat -s cli -d "add watch mode" --push --pr`
- Prompts if flags are omitted; always updates from `origin/main` first.

## Security & Configuration Tips
- Never commit secrets; use `.env` and update `.env.example` when adding vars (e.g., Tailscale auth keys).
- Restrict file writes to intended paths (e.g., `traefik/tls.yml`).
- When using Funnel or MagicDNS, verify config with a dry run before enabling `watch`.

## CI
- Combined CI: `.github/workflows/ci.yml` runs two jobs on PRs:
  - `ui`: pnpm install → typecheck → lint → build (skips if no `ui/`).
  - `go`: `go test ./... -race -cover` (skips if no Go code).
- Branch protection: in GitHub Settings → Branches, add a rule for your default branch and require status checks `CI / ui` and `CI / go`.

## Code Owners
- File: `.github/CODEOWNERS` assigns reviewers by path. Default owner is `@frnwtr` for the whole repo and key directories.
- If collaborators/teams are added later, update paths to include their handles.
- In branch protection, enable "Require review from Code Owners" to enforce.

## Architecture Overview
TailWhale syncs Docker → Tailscale certs → Traefik TLS. Choose Mode A (host + Traefik), Mode B (sidecar per container), or Mode C (Funnel on Traefik) based on deployment needs.
