.PHONY: all build test lint fmt

all: build test lint fmt

build:
	go build -o bin/hello ./cmd/hello

test:
	go test ./...

lint:
	golangci-lint run

fmt:
	go fmt ./...
