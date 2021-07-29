module github.com/peggyjv/gravity-bridge/module

go 1.15

require (
	github.com/armon/go-metrics v0.3.8 // indirect
	github.com/confio/ics23/go v0.6.6 // indirect
	github.com/cosmos/cosmos-sdk v0.42.4-0.20210623214207-eb0fc466c99b
	github.com/ethereum/go-ethereum v1.9.23
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.3-0.20201103224600-674baa8c7fc3 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/mitchellh/mapstructure v1.3.3 // indirect
	github.com/otiai10/copy v1.6.0 // indirect
	github.com/pelletier/go-toml v1.8.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.10.0 // indirect
	github.com/prometheus/common v0.23.0 // indirect
	github.com/rakyll/statik v0.1.7
	github.com/regen-network/cosmos-proto v0.3.1
	github.com/rs/zerolog v1.21.0 // indirect
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.8
	github.com/tendermint/tm-db v0.6.4
	golang.org/x/sys v0.0.0-20210510120138-977fb7262007 // indirect
	google.golang.org/genproto v0.0.0-20210114201628-6edceaf6022f
	google.golang.org/grpc v1.37.0
	google.golang.org/protobuf v1.26.0
	gopkg.in/ini.v1 v1.61.0 // indirect
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace github.com/gogo/grpc => google.golang.org/grpc v1.33.2
