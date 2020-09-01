#!/usr/bin/make -f

all: clean test build install lint

# The below include contains the tools and runsim targets.
include contrib/devtools/Makefile

########################################
### Build

build:  go.sum
	@go build -mod=readonly ./...

########################################
### Tools & dependencies

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download
.PHONY: go-mod-cache

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify
	@go mod tidy

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
	go install -mod=readonly ./cmd/ebd
	go install -mod=readonly ./cmd/ebcli
	go install -mod=readonly ./cmd/ebrelayer

# lint:
# 	@echo "--> Running linter"
# 	golangci-lint run
# 	@find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -d -s
# 	go mod verify

test:
	go test ./...

.PHONY: all build go-mod-cache build_test_container start_test_containers stop_test_containers clean install test lint all
