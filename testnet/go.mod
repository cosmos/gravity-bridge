module github.com/peggyjv/gravity-bridge/testnet

go 1.16

require (
	github.com/BurntSushi/toml v0.4.1
	github.com/cenkalti/backoff/v4 v4.1.1 // indirect
	github.com/containerd/continuity v0.1.0 // indirect
	github.com/cenkalti/backoff/v4 v4.1.1 // indirect
	github.com/cosmos/go-bip39 v1.0.0
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/moby/term v0.0.0-20210610120745-9d4ed1856297 // indirect
	github.com/ory/dockertest/v3 v3.6.5
	github.com/peggyjv/gravity-bridge/module v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.13
	github.com/willf/bitset v1.1.11 // indirect
	gotest.tools/v3 v3.0.3 // indirect
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace github.com/peggyjv/gravity-bridge/module => ../module
