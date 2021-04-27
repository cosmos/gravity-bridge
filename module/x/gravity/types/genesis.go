package types

import (
	"fmt"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// BridgeIDLen defines the length of the random bytes used for signature reuse
// prevention
const BridgeIDLen = 32

// ValidateBasic validates genesis state by looping through the params and
// calling their validation functions
func (g GenesisState) ValidateBasic() error {
	if len(g.BridgeID) != BridgeIDLen {
		return fmt.Errorf("invalid bridge ID bytes length, expected %d, got %d", BridgeIDLen, len(g.BridgeID))
	}

	if err := g.Params.ValidateBasic(); err != nil {
		return sdkerrors.Wrap(err, "params")
	}

	for _, signerSet := range g.SignerSets {
		if signerSet.Height == 0 {
			return fmt.Errorf("signer set height cannot be 0")
		}

		for _, signer := range signerSet.Signers {
			if err := ValidateEthAddress(signer.EthereumAddress); err != nil {
				return err
			}
			if signer.Power <= 0 {
				return fmt.Errorf("signer %s: consensus power cannot be 0 or negative", signer.EthereumAddress)
			}
		}
	}

	for _, batchTx := range g.BatchTxs {
		if batchTx.Block == 0 {
			return fmt.Errorf("batch tx block height cannot be 0")
		}
		if batchTx.Timeout == 0 {
			return fmt.Errorf("batch timeout cannot be 0")
		}
		if err := ValidateEthAddress(batchTx.TokenContract); err != nil {
			return err
		}
		if len(batchTx.Transactions) > int(g.Params.BatchSize) {
			return fmt.Errorf("number of batched txs (%d) > max batch size (%d)", len(batchTx.Transactions), g.Params.BatchSize)
		}
		for _, tx := range batchTx.Transactions {
			if len(tx) == 0 {
				return fmt.Errorf("tx id cannot be empty")
			}
		}
	}

	// TODO: finish
	return nil
}

// DefaultGenesisState returns a genesis state with the default parameters
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}
