package types

import (
	"fmt"
	"math/big"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

// TxIDs defines a util type for casting a slice of byte arrays into transaction
// identifiers
type TxIDs []tmbytes.HexBytes

// GetCheckpoint gets the checkpoint signature from the given outgoing tx batch
func (tx BatchTx) GetCheckpoint(bridgeID []byte, transfers []TransferTx) ([]byte, error) {
	if len(transfers) != len(tx.Transactions) {
		return nil, fmt.Errorf(
			"batch transactions length doesn't match with transfer argument (%d ≠ %d)",
			len(transfers), len(tx.Transactions),
		)
	}

	contractABI, err := abi.JSON(strings.NewReader(BatchTxCheckpointABIJSON))
	if err != nil {
		return nil, fmt.Errorf("bad ABI definition in code: %w", err)
	}

	// Create the methodName argument which salts the signature
	methodNameBytes := []uint8("transactionBatch")
	var batchMethodName [32]uint8
	copy(batchMethodName[:], methodNameBytes[:])

	// Run through the elements of the batch and serialize them

	var (
		txAmounts      = make([]*big.Int, len(tx.Transactions))
		txFees         = make([]*big.Int, len(tx.Transactions))
		txDestinations = make([]common.Address, len(tx.Transactions))
	)

	for i := 0; i < len(transfers); i++ {
		txAmounts[i] = transfers[i].Erc20Token.Amount.BigInt()
		txDestinations[i] = common.HexToAddress(transfers[i].EthereumRecipient)
		txFees[i] = transfers[i].Erc20Fee.Amount.BigInt()
	}

	// the methodName needs to be the same as the 'name' above in the checkpointAbiJson
	// but other than that it's a constant that has no impact on the output. This is because
	// it gets encoded as a function name which we must then discard.
	abiEncodedBatch, err := contractABI.Pack(
		"submitBatch",
		bridgeID,                              // bytes32
		batchMethodName,                       // bytes32
		txAmounts,                             // uint256[]
		txDestinations,                        // address[]
		txFees,                                // uint256[]
		big.NewInt(int64(tx.Nonce)),           // uint256
		common.HexToAddress(tx.TokenContract), // address
		big.NewInt(int64(tx.Timeout)),         // uint256
	)

	// this should never happen outside of test since any case that could crash on encoding
	// should be filtered above.
	if err != nil {
		return nil, sdkerrors.Wrap(err, "packing checkpoint")
	}

	// we hash the resulting encoded bytes discarding the first 4 bytes these 4 bytes are the constant
	// method name 'checkpoint'. If you where to replace the checkpoint constant in this code you would
	// then need to adjust how many bytes you truncate off the front to get the output of abi.encode()
	return crypto.Keccak256Hash(abiEncodedBatch[4:]).Bytes(), nil
}

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

// GetCheckpoint gets the checkpoint signature from the given outgoing tx batch
func (tx LogicCallTx) GetCheckpoint(bridgeID []byte, invalidationID tmbytes.HexBytes, invalidationNonce uint64) ([]byte, error) {
	contractABI, err := abi.JSON(strings.NewReader(LogicCallTxABIJSON))
	if err != nil {
		return nil, sdkerrors.Wrap(err, "bad ABI definition in code")
	}

	// Create the methodName argument which salts the signature
	methodNameBytes := []uint8("logicCall")
	var logicCallMethodName [32]uint8
	copy(logicCallMethodName[:], methodNameBytes[:])

	// Run through the elements of the logic call and serialize them
	var (
		transferAmounts        = make([]*big.Int, len(tx.Tokens))
		feeAmounts             = make([]*big.Int, len(tx.Fees))
		transferTokenContracts = make([]common.Address, len(tx.Tokens))
		feeTokenContracts      = make([]common.Address, len(tx.Fees))
	)

	for i, tx := range tx.Tokens {
		transferAmounts[i] = tx.Amount.BigInt()
		transferTokenContracts[i] = common.HexToAddress(tx.Denom)
	}

	for i, tx := range tx.Fees {
		feeAmounts[i] = tx.Amount.BigInt()
		feeTokenContracts[i] = common.HexToAddress(tx.Denom)
	}

	payload := make([]byte, len(tx.Payload))
	copy(payload, tx.Payload)
	var invalidationIDCopy [32]byte
	copy(invalidationIDCopy[:], invalidationID[:])

	// the methodName needs to be the same as the 'name' above in the checkpointAbiJson
	// but other than that it's a constant that has no impact on the output. This is because
	// it gets encoded as a function name which we must then discard.
	abiEncodedCall, err := contractABI.Pack(
		"checkpoint",
		bridgeID,
		logicCallMethodName,
		transferAmounts,
		transferTokenContracts,
		feeAmounts,
		feeTokenContracts,
		common.HexToAddress(tx.LogicContractAddress),
		payload,
		new(big.Int).SetUint64(tx.Timeout),
		invalidationIDCopy,
		new(big.Int).SetUint64(invalidationNonce),
	)

	if err != nil {
		// this should never happen outside of test since any case that could crash on encoding
		// should be filtered above.
		return nil, sdkerrors.Wrap(err, "packing checkpoint")
	}

	hash := crypto.Keccak256Hash(abiEncodedCall[4:])
	return hash.Bytes(), nil
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
