.DEFAULT_GOAL := help
APPNAME := tgBot
OUTPUT_DIR := ./builds
BINARY_ENTRYPOINT := ./telegramBot-cli.go
LINUX := linux
WIN := windows
MAC := darwin
AMD64 := amd64
ARM64 := arm64

help list: ## Show available commands with descriptions
		@echo ""
		@echo "Available Makefile commands:"
		@echo "---------------------------------------------"
#		Searches all Makefiles in top‑level plus included ones
#		Search for a word (letter plus "-" sign) then a ":" and then a "##"
#		Set the separator as ": ##"
#		Print with ANSI escape
		@grep -E '^[a-zA-Z0-9_-]+:.*##' $(MAKEFILE_LIST) | \
			awk 'BEGIN {FS = ":.*## "}; \
			{printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'
		@echo "---------------------------------------------"
		@echo ""

release: ## Build for linux, windows and macOS
	make -s generate-conf
	make -s build-linux-amd64
	make -s build-linux-arm64
	make -s build-win-amd64
	make -s build-macos-amd64
	make -s build-macos-arm64

generate-conf: ## Generate the config template
# 	Get the README.md file
#	Read the toml section
#	Create the file and write the config
	@cat 	./README.md | \
			awk '/^```toml$$/{on=1; next} /^```$$/{on=0} \
			on{print}' > $(OUTPUT_DIR)/config.toml

build-linux-amd64: ## Build for linux-amd64
	GOARCH=$(AMD64) GOOS=$(LINUX) CGO_ENABLED=0 \
	go build -o $(OUTPUT_DIR)/$(APPNAME)_$(LINUX)_$(AMD64) $(BINARY_ENTRYPOINT)
	@sha256sum $(OUTPUT_DIR)/$(APPNAME)_$(LINUX)_$(AMD64) | awk -F ' ' '{print $$1}' > $(OUTPUT_DIR)/$(APPNAME)_$(LINUX)_$(AMD64)_hash.txt

build-linux-arm64: ## Build for linux-arm64
	GOARCH=$(ARM64) GOOS=$(LINUX) CGO_ENABLED=0 \
	go build -o $(OUTPUT_DIR)/$(APPNAME)_$(LINUX)_$(ARM64) $(BINARY_ENTRYPOINT)
	@sha256sum $(OUTPUT_DIR)/$(APPNAME)_$(LINUX)_$(ARM64) | awk -F ' ' '{print $$1}' > $(OUTPUT_DIR)/$(APPNAME)_$(LINUX)_$(ARM64)_hash.txt

build-win-amd64: ## Build for windows-amd64
	GOARCH=$(AMD64) GOOS=$(WIN) CGO_ENABLED=0 \
	go build -o $(OUTPUT_DIR)/$(APPNAME)_$(WIN)_$(AMD64) $(BINARY_ENTRYPOINT)
	@sha256sum $(OUTPUT_DIR)/$(APPNAME)_$(WIN)_$(AMD64) | awk -F ' ' '{print $$1}' > $(OUTPUT_DIR)/$(APPNAME)_$(WIN)_$(AMD64)_hash.txt

build-macos-amd64: ## Build for macos-amd64
	GOARCH=$(AMD64) GOOS=$(MAC) CGO_ENABLED=0 \
	go build -o $(OUTPUT_DIR)/$(APPNAME)_$(MAC)_$(AMD64) $(BINARY_ENTRYPOINT)
	@sha256sum $(OUTPUT_DIR)/$(APPNAME)_$(MAC)_$(AMD64) | awk -F ' ' '{print $$1}' > $(OUTPUT_DIR)/$(APPNAME)_$(MAC)_$(AMD64)_hash.txt

build-macos-arm64: ## Build for macos-arm64
	GOARCH=$(ARM64) GOOS=$(MAC) CGO_ENABLED=0 \
	go build -o $(OUTPUT_DIR)/$(APPNAME)_$(MAC)_$(ARM64) $(BINARY_ENTRYPOINT)
	@sha256sum $(OUTPUT_DIR)/$(APPNAME)_$(MAC)_$(ARM64) | awk -F ' ' '{print $$1}' > $(OUTPUT_DIR)/$(APPNAME)_$(MAC)_$(ARM64)_hash.txt