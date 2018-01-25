.PHONY: install test test_unit test_cli get_vendor_deps tools

COMMIT_HASH := $(shell git rev-parse --short HEAD)
LINKER_FLAGS:="-X github.com/cosmos/cosmos-sdk/client/commands.CommitHash=${COMMIT_HASH}"
NOVENDOR := $(shell glide novendor)

GOTOOLS =	github.com/mitchellh/gox \
			github.com/Masterminds/glide \
			github.com/rigelrozanski/shelldown/cmd/shelldown

install:
	@go install -ldflags $(LINKER_FLAGS) ./cmd/...

test: test_unit test_cli

test_unit:
	@go test $(NOVENDOR)

test_cli:
	./tests/cli/basictx.sh
	./tests/cli/keys.sh
	./tests/cli/init.sh
	./tests/cli/rpc.sh
	./tests/cli/restart.sh


get_vendor_deps: tools
	@glide install

tools:
	@go get $(GOTOOLS)

