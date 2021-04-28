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

		if err := signerSet.Signers.ValidateBasic(); err != nil {
			return err
		}
	}

	for _, batchTx := range g.BatchTxs {
		if len(batchTx.Transactions) > int(g.Params.BatchSize) {
			return fmt.Errorf("number of batched txs (%d) > max batch size (%d)", len(batchTx.Transactions), g.Params.BatchSize)
		}

		if err := batchTx.Validate(); err != nil {
			return err
		}
	}

	for _, tx := range g.LogicCallTxs {
		if err := tx.Validate(); err != nil {
			return err
		}
	}

	for _, tx := range g.TransferTxs {
		if err := tx.Validate(); err != nil {
			return err
		}
	}

	for _, attestation := range g.Attestations {
		if err := attestation.Validate(); err != nil {
			return err
		}
	}

	for _, keyDelegation := range g.DelegateKeys {
		if err := keyDelegation.ValidateBasic(); err != nil {
			return err
		}
	}

	for _, e := range g.Erc20ToDenoms {
		if err := e.Validate(); err != nil {
			return err
		}
	}

	// TODO: validate confirms
	return nil
}

// DefaultGenesisState returns a genesis state with the default parameters
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}
