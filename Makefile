SHELL := /bin/bash

.PHONY: help build test dev ui-dev ui-build go-build go-test

help:
	@echo "Available targets:"
	@echo "  build    - Build Go CLI (if present) and UI"
	@echo "  test     - Run Go tests (if present) and UI tests"
	@echo "  dev      - Start UI dev server (pnpm)"
	@echo "  ui-dev   - Start UI dev server"
	@echo "  ui-build - Build UI for production"
	@echo "  go-build - Build Go CLI if cmd/tailwhale exists"
	@echo "  go-test  - Run Go tests if Go sources exist"

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

