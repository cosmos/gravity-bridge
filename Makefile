#!/usr/bin/make -f

export GO111MODULE = on

all: test_app clean install lint

# The below include contains the tools and runsim targets.
include contrib/devtools/Makefile

build:
	go build ./cmd/ebd
	go build ./cmd/ebcli
	go build ./cmd/ebrelayer

clean:
	rm -f ebd
	rm -f ebcli
	rm -f ebrelayer

install:
	go install ./cmd/ebd
	go install ./cmd/ebcli
	go install ./cmd/ebrelayer

lint:
	@echo "--> Running linter"
	golangci-lint run
	@find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -d -s
	go mod verify

test_app:
	go test ./...

.PHONY: all build clean install test_app lint all