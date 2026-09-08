GO ?= go
PNPM ?= pnpm
FORGE ?= forge

.PHONY: build build-go build-web build-contracts test dev-web

build: build-go build-web build-contracts

build-go:
	$(GO) build ./cmd/arbex

build-web:
	$(PNPM) --filter @arbex/web build

build-contracts:
	$(FORGE) build

test:
	$(GO) test ./...
	$(PNPM) --filter @arbex/web typecheck
	$(FORGE) test

dev-web:
	$(PNPM) --filter @arbex/web dev
