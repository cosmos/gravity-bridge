module github.com/cosmos/gravity-bridge/testnet

go 1.16

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/Microsoft/go-winio v0.5.0 // indirect
	github.com/cenkalti/backoff/v4 v4.1.1 // indirect
	github.com/containerd/continuity v0.1.0 // indirect
	github.com/cosmos/cosmos-sdk v0.43.0
	github.com/cosmos/go-bip39 v1.0.0
	github.com/ethereum/go-ethereum v1.9.25
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/moby/term v0.0.0-20210610120745-9d4ed1856297 // indirect
	github.com/opencontainers/runc v1.0.0-rc95 // indirect
	github.com/ory/dockertest/v3 v3.6.5
	github.com/cosmos/gravity-bridge/module v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.11
	golang.org/x/net v0.0.0-20210610132358-84b48f89b13b // indirect
	golang.org/x/sys v0.0.0-20210611083646-a4fc73990273 // indirect
	gotest.tools/v3 v3.0.3 // indirect
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace github.com/cosmos/gravity-bridge/module => ../module
