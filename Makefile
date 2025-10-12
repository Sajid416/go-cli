.PHONY: all build test lint fmt fuzz bench coverage-check

all: build test lint fmt fuzz bench

build:
	go build -o bin/hello ./cmd/hello

test:
	go test ./...

lint:
	golangci-lint run

fmt:
	go fmt ./...

fuzz: ## run fuzz tests (~10s)
	go test -run=^$$ -fuzz=Fuzz -fuzztime=10s ./pkg/parser
	

bench: ## run benchmarks
	go test -bench=. -benchmem ./...

coverage-check: ##fail if coverage <70%
	@go tool cover -func=coverage.out | awk '/total:/ {split($$3,a,"%"); if (a[1]<70) {print "Coverage too low ("$$3")"; exit 1 } else { print "Coverage OK ("$$3")" }}'