# Use bash everywhere (helps on CI/Windows Git-Bash)
SHELL := /usr/bin/env bash

# Detect host OS via Go
GOOS := $(shell go env GOOS)

# Packages we care about
PKGS       := ./cmd/... ./pkg/...
BINARY     := bin/hello

# Fuzz settings (customizable)
FUZZPKG    ?= ./pkg/parser
FUZZTARGET ?= Fuzz
FUZZTIME   ?= 10s

# OS-specific flags
ifeq ($(GOOS),windows)
  # On Windows, avoid cgo issues and skip -race
  CGOEN := 0
  RACE  :=
else
  # On Linux/macOS, enable cgo and race detector
  CGOEN := 1
  RACE  := -race
endif

.PHONY: all build test lint fmt fuzz bench coverage coverage-check clean

# Keep "all" fast & deterministic
all: fmt lint test build

build:
	CGEO_ENABLED=$(CGOEN) go build -o $(BINARY) ./cmd/hello

# NOTE: use explicit package globs (avoid weird dirs)
test:
	CGO_ENABLED=$(CGOEN) go test $(RACE) -coverprofile=coverage.out $(PKGS)

lint:
	golangci-lint run

fmt:
	go fmt ./...

# Fuzz must target a single package (go test -fuzz doesn't support ./...)
fuzz: ## run fuzz tests (~10s)
	go test -run=^$$ -fuzz=$(FUZZTARGET) -fuzztime=$(FUZZTIME) $(FUZZPKG)

bench: ## run benchmarks
	go test -bench=. -benchmem $(PKGS)

coverage:
	go tool cover -func=coverage.out

coverage-check: ## fail if coverage < 70%
	@go tool cover -func=coverage.out | awk '/total:/ {split($$3,a,"%"); if (a[1] < 70) { print "Coverage too low ("$$3")"; exit 1 } else { print "Coverage OK ("$$3")" }}'

clean:
	rm -rf $(BINARY) coverage.out
