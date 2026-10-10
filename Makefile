.DEFAULT_GOAL := help

APP_NAME := telegram-cli
OUTPUT_DIR := builds
ENTRYPOINT := ./cmd/telegram-cli
VERSION ?= dev
LDFLAGS := -s -w -X github.com/P1zz/telegramBot-cli/cmd.version=$(VERSION)
PLATFORMS := linux/amd64 linux/arm64 windows/amd64 darwin/amd64 darwin/arm64

.PHONY: help test vet fmt-check build release clean generate-config

help: ## Show available commands
	@awk 'BEGIN {FS = ":.*## "}; /^[a-zA-Z0-9_-]+:.*## / {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

test: ## Run unit tests
	go test ./...

vet: ## Run Go static analysis
	go vet ./...

fmt-check: ## Check Go formatting
	@test -z "$$(gofmt -l .)" || (gofmt -l . && exit 1)

build: ## Build for the current platform
	@mkdir -p $(OUTPUT_DIR)
	go build -trimpath -ldflags "$(LDFLAGS)" -o $(OUTPUT_DIR)/$(APP_NAME) $(ENTRYPOINT)

release: clean generate-config ## Build release binaries and SHA-256 checksums
	@set -e; \
	for platform in $(PLATFORMS); do \
		os=$${platform%/*}; arch=$${platform#*/}; extension=""; \
		if [ "$$os" = windows ]; then extension=".exe"; fi; \
		name="$(APP_NAME)_$${os}_$${arch}$${extension}"; \
		echo "Building $$name"; \
		GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o "$(OUTPUT_DIR)/$$name" $(ENTRYPOINT); \
		sha256sum "$(OUTPUT_DIR)/$$name" > "$(OUTPUT_DIR)/$$name.sha256"; \
	done

generate-config: ## Generate example config from the README
	@mkdir -p $(OUTPUT_DIR)
	@awk '/^```toml$$/{on=1; next} /^```$$/{if(on){exit}} on{print}' README.md > $(OUTPUT_DIR)/config.example.toml

clean: ## Remove generated builds
	@rm -rf -- $(OUTPUT_DIR)
