MISE ?= mise

.DEFAULT_GOAL := help

.PHONY: help setup build build-all test lint vet release-check clean

help:
	@printf "Available targets:\n"
	@printf "  make setup         Install tools from mise.toml\n"
	@printf "  make build         Build local dellkvm binary with GoReleaser\n"
	@printf "  make build-all     Build all configured GoReleaser targets\n"
	@printf "  make test          Run tests\n"
	@printf "  make lint          Run golangci-lint\n"
	@printf "  make vet           Run go vet\n"
	@printf "  make release-check Validate GoReleaser config\n"
	@printf "  make clean         Remove local build artifacts\n"

setup:
	@command -v $(MISE) >/dev/null 2>&1 || { \
		echo "mise not found. Install it first: https://mise.jdx.dev/getting-started.html"; \
		exit 1; \
	}
	$(MISE) trust
	$(MISE) install

build:
	$(MISE) run build

build-all:
	$(MISE) run build:all

test:
	$(MISE) run test

lint:
	$(MISE) run lint

vet:
	$(MISE) run vet

release-check:
	$(MISE) run release:check

clean:
	rm -rf dist dellkvm
