SHELL := /bin/bash

.PHONY: help build test dev ui-dev ui-build go-build go-test demo pr-template pr-update

help:
	@echo "Available targets:"
	@echo "  build    - Build Go CLI (if present) and UI"
	@echo "  test     - Run Go tests (if present) and UI tests"
	@echo "  dev      - Start UI dev server (pnpm)"
	@echo "  ui-dev   - Start UI dev server"
	@echo "  ui-build - Build UI for production"
	@echo "  go-build - Build Go CLI if cmd/tailwhale exists"
	@echo "  go-test  - Run Go tests if Go sources exist"
	@echo "  demo     - Run CLI demo with examples"
	@echo "  pr-template - Generate PR_BODY.md scaffold from scripts/pr-body-sample.md"
	@echo "  pr-update  - Update PR body with PR_BODY.md (usage: make pr-update PR=123)"

build: go-build ui-build

test: go-test
	@if [ -d ui ]; then \
		cd ui && pnpm test; \
	else \
		echo "[skip] ui/ not found"; \
	fi

dev: ui-dev

ui-dev:
	@if [ -d ui ]; then \
		cd ui && pnpm dev; \
	else \
		echo "ui/ not found. Create ui/ or run 'make help'"; \
	fi

ui-build:
	@if [ -d ui ]; then \
		cd ui && pnpm build; \
	else \
		echo "[skip] ui/ not found"; \
	fi

go-build:
	@if [ -d cmd/tailwhale ]; then \
		go build ./cmd/tailwhale; \
	else \
		echo "[skip] cmd/tailwhale not found"; \
	fi

go-test:
	@if rg -uu --files | rg -q "\.go$$"; then \
		go test ./...; \
	else \
		echo "[skip] no Go files found"; \
	fi

demo:
	@echo "Listing services from examples/containers.json";
	@go run ./cmd/tailwhale list --from-file examples/containers.json || (go build ./cmd/tailwhale && ./tailwhale list --from-file examples/containers.json);
	@echo "Writing tls.yml preview to /tmp/tailwhale_tls.yml using examples/tailwhale.json";
	@go run ./cmd/tailwhale sync --config examples/tailwhale.json --tls-path /tmp/tailwhale_tls.yml || ./tailwhale sync --config examples/tailwhale.json --tls-path /tmp/tailwhale_tls.yml;
	@echo "Wrote: /tmp/tailwhale_tls.yml";

pr-template:
	@bash scripts/pr-update.sh --generate ./PR_BODY.md

pr-update:
	@if [ -z "$(PR)" ]; then echo "Usage: make pr-update PR=123"; exit 2; fi
	@bash scripts/pr-update.sh --pr $(PR) --body ./PR_BODY.md
