package witness

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Oracle payload for Eth -> Cosmos
type PayloadLock struct {
	Nonce    uint64
	Coins    sdk.Coins
	DestAddr sdk.Address
}

func (p PayloadLock) Type() string {
	return "witness"
}

func (p PayloadLock) ValidateBasic() sdk.Error {
	if !p.Coins.IsValid() {
		return sdk.ErrInvalidCoins("")
	}
	return nil
}
