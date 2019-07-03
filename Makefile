.PHONY: get_tools get_vendor_deps update_vendor_deps install test

DEP := $(shell command -v dep 2> /dev/null)

ldflags = -X github.com/cosmos/sdk-application-tutorial/version.Version=$(VERSION) \
	-X github.com/cosmos/sdk-application-tutorial/version.Commit=$(COMMIT)

get_tools:
ifndef DEP
	@echo "Installing dep"
	go get -u -v github.com/golang/dep/cmd/dep
else
	@echo "Dep is already installed..."
endif

get_vendor_deps:
	@echo "--> Generating vendor directory via dep ensure"
	@rm -rf .vendor-new
	@dep ensure -v -vendor-only

update_vendor_deps:
	@echo "--> Running dep ensure"
	@rm -rf .vendor-new
	@dep ensure -v -update

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
	gotestsum