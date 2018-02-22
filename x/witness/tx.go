package withdraw

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    crypto "github.com/tendermint/go-crypto"
)

type WitnessTx struct {
    amount int64
    destination crypto.Address
    Token crypto.Address
}

var _ sdk.Msg = (*WitnessTx)(nil)

func (wtx WitnessTx) ValidateBasic() sdk.Error {
    return nil
}

func (wtx WitnessTx) Type() string {
    return "WitnessTx"
}

type WitnessData struct {
    Witnesses      []crypto.Address
    Amount         int64
    Destination    crypto.Address
    credited       bool
}

