.PHONY: lint
lint: 
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0 run --fast -- $(go list -f '{{.Dir}}/...' -m)

.PHONY: pre-build
pre-build:
	go run cmd/pre-build/pre-build.go

.PHONY: build
build: pre-build
	go build -o ./bin/cli ./cmd/cli/*.go 

.PHONY: test
test:
	
