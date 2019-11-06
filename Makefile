#!/usr/bin/make -f

all: test clean install lint

# The below include contains the tools and runsim targets.
include contrib/devtools/Makefile

build:
	go build ./cmd/ebd
	go build ./cmd/ebcli
	go build ./cmd/ebrelayer

build_test_container:
	docker-compose -f ./deploy/test/docker-compose.yml --project-directory . build

start_test_containers:
	docker-compose -f ./deploy/test/docker-compose.yml --project-directory . up

stop_test_containers:
	docker-compose -f ./deploy/test/docker-compose.yml --project-directory . down

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

test:
	go test ./...

.PHONY: all build build_test_container start_test_containers stop_test_containers clean install test lint all
