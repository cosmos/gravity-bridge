PACKAGES=$(shell go list ./... | grep -v '/simulation')

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

# TODO: Update the ldflags with the app, client & server names
ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=Peggy \
	-X github.com/cosmos/cosmos-sdk/version.ServerName=peggyd \
	-X github.com/cosmos/cosmos-sdk/version.ClientName=peggycli \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) 

BUILD_FLAGS := -ldflags '$(ldflags)' -gcflags="all=-N -l"

all: install

install: go.sum
		go install $(BUILD_FLAGS) ./cmd/peggyd
		go install $(BUILD_FLAGS) ./cmd/peggycli

go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		GO111MODULE=on go mod verify

test:
	@go test -mod=readonly $(PACKAGES)

# look into .golangci.yml for enabling / disabling linters
lint:
	@echo "--> Running linter"
	@golangci-lint run
	@go mod verify
