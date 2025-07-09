# Define LATEST_TAG to get the latest Git tag, or fallback to "v0.0.0" if no tag exists
LATEST_TAG=$(shell git describe --tags $$(git rev-list --tags --max-count=1) 2>/dev/null || echo "v0.0.0")

# Define VERSION to remove the "v" prefix from LATEST_TAG
VERSION=$(shell echo $(LATEST_TAG) | sed 's/^v//')

.PHONY: all
all: lint generate pre-build build test

# Install deps
.PHONY: setup
setup:
	go install github.com/abice/go-enum@latest
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.2.1

.PHONY: lint
lint: 
	golangci-lint run -- $(go list -f '{{.Dir}}/...' -m)

generate:
	go generate ./...

.PHONY: pre-build
pre-build: generate
	go run cmd/pre-build/pre-build.go --write-ranges ranges

.PHONY: build
build: pre-build generate
	go build -ldflags="-X main.version=$(VERSION)" -o ./bin/cli ./cmd/cli/*.go 

.PHONY: test
test:
	go test ./...

.PHONY: get-version
get-version:
	@echo $(VERSION)
