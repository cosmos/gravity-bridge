.PHONY: get_tools get_vendor_deps update_vendor_deps build clean install test

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

install:
	go install ./cmd/ebd
	go install ./cmd/ebcli
	go install ./cmd/ebrelayer

test:
	go test ./...