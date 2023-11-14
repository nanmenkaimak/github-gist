# Tools.
TOOLS = ./tools
TOOLS_BIN = $(TOOLS)/bin

.PHONY: fix-lint
fix-lint: $(TOOLS_BIN)/golangci-lint
	$(TOOLS_BIN)/golangci-lint run --fix

imports: $(TOOLS_BIN)/goimports
	$(TOOLS_BIN)/goimports -local "service" -w ./internal ./cmd

# INSTALL linter
$(TOOLS_BIN)/golangci-lint: export GOBIN = $(shell pwd)/$(TOOLS_BIN)
$(TOOLS_BIN)/golangci-lint:
	mkdir -p $(TOOLS_BIN)
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2


# INSTALL goimports
$(TOOLS_BIN)/goimports: export GOBIN = $(shell pwd)/$(TOOLS_BIN)
$(TOOLS_BIN)/goimports:
	mkdir -p $(TOOLS_BIN)
	go install golang.org/x/tools/cmd/goimports@latest