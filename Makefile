#!/usr/bin/make -f

export GO111MODULE = on

all: format clean test_app install lint

# The below include contains the tools and runsim targets.
include contrib/devtools/Makefile

DEP := $(shell command -v dep 2> /dev/null)

ldflags = -X github.com/cosmos/sdk-application-tutorial/version.Version=$(VERSION) \
	-X github.com/cosmos/sdk-application-tutorial/version.Commit=$(COMMIT)

build:
	go build ./cmd/ebd
	go build ./cmd/ebcli
	go build ./cmd/ebrelayer

clean:
	rm -f ebd
	rm -f ebcli
	rm -f ebrelayer

install: build
	go install ./cmd/ebd
	go install ./cmd/ebcli
	go install ./cmd/ebrelayer

lint:
	@echo "--> Running linter"
	@golangci-lint run

test_app:
	go test ./...

format: tools
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs goimports -w -local github.com/cosmos/cosmos-sdk

.PHONY: build clean install test format lint