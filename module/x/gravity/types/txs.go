package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (tx BatchTx) Validate() error {
	if tx.Nonce == 0 {
		return fmt.Errorf("nonce cannot be 0")
	}
	if tx.Block == 0 {
		return fmt.Errorf("batch tx block height cannot be 0")
	}
	if tx.Timeout == 0 {
		return fmt.Errorf("batch timeout cannot be 0")
	}
	if err := ValidateEthAddress(tx.TokenContract); err != nil {
		return err
	}
	for _, tx := range tx.Transactions {
		if len(tx) == 0 {
			return fmt.Errorf("tx id cannot be empty")
		}
	}

	return nil
}

func (tx TransferTx) Validate() error {
	if tx.Nonce == 0 {
		return fmt.Errorf("nonce cannot be 0")
	}
	_, err := sdk.AccAddressFromBech32(tx.Sender)
	if err != nil {
		return err
	}
	if err := ValidateEthAddress(tx.EthereumRecipient); err != nil {
		return err
	}
	if !tx.Erc20Token.IsPositive() {
		return fmt.Errorf("erc20 token amount must be positive: %s", tx.Erc20Token)
	}
	if err := tx.Erc20Token.Validate(); err != nil {
		return fmt.Errorf("invalid erc20 token: %s", err)
	}
	if err := tx.Erc20Fee.Validate(); err != nil {
		return fmt.Errorf("invalid erc20 fee: %s", err)
	}

	return nil
}

func (tx LogicCallTx) Validate() error {
	if tx.Nonce == 0 {
		return fmt.Errorf("nonce cannot be 0")
	}
	if err := tx.Tokens.Validate(); err != nil {
		return fmt.Errorf("invalid erc20 tokens: %s", err)
	}
	if err := tx.Fees.Validate(); err != nil {
		return fmt.Errorf("invalid erc20 fees: %s", err)
	}
	if err := ValidateEthAddress(tx.LogicContractAddress); err != nil {
		return err
	}
	if len(tx.Payload) == 0 {
		return fmt.Errorf("payload bytes cannot be empty")
	}
	if tx.Timeout == 0 {
		return fmt.Errorf("tx timeout cannot be 0")
	}

	return nil
}
