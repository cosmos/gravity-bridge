package types

import (
	"fmt"
	"math/big"
	"strings"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// TxIDs defines a util type for casting a slice of byte arrays into transaction
// identifiers
type TxIDs []tmbytes.HexBytes

// GetCheckpoint gets the checkpoint signature from the given outgoing tx batch
func (b BatchTx) GetCheckpoint(bridgeID []byte, transfers []TransferTx) ([]byte, error) {
	if len(transfers) != len(b.Transactions) {
		return nil, fmt.Errorf(
			"batch transactions length doesn't match with transfer argument (%d ≠ %d)",
			len(transfers), len(b.Transactions),
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
		txAmounts      = make([]*big.Int, len(b.Transactions))
		txFees         = make([]*big.Int, len(b.Transactions))
		txDestinations = make([]common.Address, len(b.Transactions))
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
		bridgeID,                             // bytes32
		batchMethodName,                      // bytes32
		txAmounts,                            // uint256[]
		txDestinations,                       // address[]
		txFees,                               // uint256[]
		big.NewInt(int64(b.Nonce)),           // uint256
		common.HexToAddress(b.TokenContract), // address
		big.NewInt(int64(b.Timeout)),         // uint256
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

// GetCheckpoint gets the checkpoint signature from the given outgoing tx batch
func (c LogicCallTx) GetCheckpoint(bridgeID []byte, invalidationID tmbytes.HexBytes, invalidationNonce uint64) ([]byte, error) {
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
		transferAmounts        = make([]*big.Int, len(c.Tokens))
		feeAmounts             = make([]*big.Int, len(c.Fees))
		transferTokenContracts = make([]common.Address, len(c.Tokens))
		feeTokenContracts      = make([]common.Address, len(c.Fees))
	)

	for i, tx := range c.Tokens {
		transferAmounts[i] = tx.Amount.BigInt()
		transferTokenContracts[i] = common.HexToAddress(tx.Denom)
	}
	for i, tx := range c.Fees {
		feeAmounts[i] = tx.Amount.BigInt()
		feeTokenContracts[i] = common.HexToAddress(tx.Denom)
	}
	payload := make([]byte, len(c.Payload))
	copy(payload, c.Payload)
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
		common.HexToAddress(c.LogicContractAddress),
		payload,
		new(big.Int).SetUint64(c.Timeout),
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
