.PHONY: all
all: lint generate pre-build build test

.PHONY: lint
lint: 
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0 run --fast -- $(go list -f '{{.Dir}}/...' -m)

generate:
	go generate ./...

.PHONY: pre-build
pre-build: generate
	go run cmd/pre-build/pre-build.go

.PHONY: build
build: pre-build generate
	go build -o ./bin/cli ./cmd/cli/*.go 

.PHONY: test
test:
	go test ./...
