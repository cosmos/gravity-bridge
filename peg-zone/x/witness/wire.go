package witness

import (
	wire "github.com/tendermint/go-wire"
)

func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(PayloadLock{}, "peggy/oracle/Lock", nil)
}
