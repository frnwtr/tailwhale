# TailWhale UI

A minimal Next.js TypeScript app managed with pnpm and Biome.

## Commands
- Install deps: `pnpm install`
- Dev server: `pnpm dev`
- Build: `pnpm build`
- Start: `pnpm start`
- Typecheck: `pnpm typecheck`
- Lint: `pnpm lint`
- Format: `pnpm format`

Enable Corepack globally with `corepack enable` and activate pnpm via `corepack prepare pnpm@latest --activate`.

## Git Hooks
- Husky pre-commit runs typecheck + lint from repo root.
- After `pnpm install`, the `prepare` script auto-installs hooks.
