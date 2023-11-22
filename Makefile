# Tools.
TOOLS = ./tools
TOOLS_BIN = $(TOOLS)/bin

generate-swagger-gist:
	docker run --rm -it  \
		-u $(shell id -u):$(shell id -g) \
		-e GOPATH=$(shell go env GOPATH):/go \
		-e GOCACHE=/tmp \
		-v $(HOME):$(HOME) \
		-w $(shell pwd) \
		quay.io/goswagger/swagger:0.30.4 \
		generate spec -c ./cmd/gist --scan-models -c ./internal/gist -o ./swagger/OpenAPI/gist.rest.swagger.json

generate-swagger-auth:
	docker run --rm -it  \
		-u $(shell id -u):$(shell id -g) \
		-e GOPATH=$(shell go env GOPATH):/go \
		-e GOCACHE=/tmp \
		-v $(HOME):$(HOME) \
		-w $(shell pwd) \
		quay.io/goswagger/swagger:0.30.4 \
		generate spec -c ./cmd/auth --scan-models -c ./internal/auth -o ./swagger/OpenAPI/auth.rest.swagger.json

generate-swagger-admin:
	docker run --rm -it  \
		-u $(shell id -u):$(shell id -g) \
		-e GOPATH=$(shell go env GOPATH):/go \
		-e GOCACHE=/tmp \
		-v $(HOME):$(HOME) \
		-w $(shell pwd) \
		quay.io/goswagger/swagger:0.30.4 \
		generate spec -c ./cmd/admin --scan-models -c ./internal/admin -o ./swagger/OpenAPI/admin.rest.swagger.json


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