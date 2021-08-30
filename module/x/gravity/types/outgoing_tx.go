package types

import (
	"math/big"
	"strings"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	_ OutgoingTx = &SignerSetTx{}
	_ OutgoingTx = &BatchTx{}
	_ OutgoingTx = &ContractCallTx{}
)

const (
	_ = iota
	SignerSetTxPrefixByte
	BatchTxPrefixByte
	ContractCallTxPrefixByte
)

///////////////////
// GetStoreIndex //
///////////////////

// TODO: do we need a prefix byte for the different types?
func (sstx *SignerSetTx) GetStoreIndex() []byte {
	return MakeSignerSetTxKey(sstx.Nonce)
}

func (btx *BatchTx) GetStoreIndex() []byte {
	return MakeBatchTxKey(gethcommon.HexToAddress(btx.TokenContract), btx.BatchNonce)
}

func (cctx *ContractCallTx) GetStoreIndex() []byte {
	return MakeContractCallTxKey(cctx.InvalidationScope.Bytes(), cctx.InvalidationNonce)
}

///////////////////
// GetCheckpoint //
///////////////////

func (sstx *SignerSetTx) GetCosmosHeight() uint64 {
	return sstx.Height
}

func (btx *BatchTx) GetCosmosHeight() uint64 {
	return btx.Height
}

func (cctx *ContractCallTx) GetCosmosHeight() uint64 {
	return cctx.Height
}

///////////////////
// GetCheckpoint //
///////////////////

// GetCheckpoint returns the checkpoint
func (sstx SignerSetTx) GetCheckpoint(gravityID string) []byte {

	// the contract argument is not a arbitrary length array but a fixed length 32 byte
	// array, therefore we have to utf8 encode the string (the default in this case) and
	// then copy the variable length encoded data into a fixed length array. This function
	// will panic if gravityId is too long to fit in 32 bytes
	gravityIDFixed, err := strToFixByteArray(gravityID)
	if err != nil {
		panic(err)
	}

	checkpointBytes := []uint8("checkpoint")
	var checkpoint [32]uint8
	copy(checkpoint[:], checkpointBytes[:])

	memberAddresses := make([]gethcommon.Address, len(sstx.Signers))
	convertedPowers := make([]*big.Int, len(sstx.Signers))
	for i, m := range sstx.Signers {
		memberAddresses[i] = gethcommon.HexToAddress(m.EthereumAddress)
		convertedPowers[i] = big.NewInt(int64(m.Power))
	}
	// the word 'checkpoint' needs to be the same as the 'name' above in the checkpointAbiJson
	// but other than that it's a constant that has no impact on the output. This is because
	// it gets encoded as a function name which we must then discard.
	args := []interface{}{
		gravityIDFixed,
		checkpoint,
		big.NewInt(int64(sstx.Nonce)),
		memberAddresses,
		convertedPowers,
	}

	return packCall(SignerSetTxCheckpointABIJSON, "checkpoint", args)
}

// GetCheckpoint gets the checkpoint signature from the given outgoing tx batch
func (btx BatchTx) GetCheckpoint(gravityID string) []byte {

	// the contract argument is not a arbitrary length array but a fixed length 32 byte
	// array, therefore we have to utf8 encode the string (the default in this case) and
	// then copy the variable length encoded data into a fixed length array. This function
	// will panic if gravityId is too long to fit in 32 bytes
	gravityIDFixed, err := strToFixByteArray(gravityID)
	if err != nil {
		panic(err)
	}

	// Create the methodName argument which salts the signature
	methodNameBytes := []uint8("transactionBatch")
	var batchMethodName [32]uint8
	copy(batchMethodName[:], methodNameBytes[:])

	// Run through the elements of the batch and serialize them
	txAmounts := make([]*big.Int, len(btx.Transactions))
	txDestinations := make([]gethcommon.Address, len(btx.Transactions))
	txFees := make([]*big.Int, len(btx.Transactions))
	for i, tx := range btx.Transactions {
		txAmounts[i] = tx.Erc20Token.Amount.BigInt()
		txDestinations[i] = gethcommon.HexToAddress(tx.EthereumRecipient)
		txFees[i] = tx.Erc20Fee.Amount.BigInt()
	}

	// the methodName needs to be the same as the 'name' above in the checkpointAbiJson
	// but other than that it's a constant that has no impact on the output. This is because
	// it gets encoded as a function name which we must then discard.
	args := []interface{}{
		gravityIDFixed,
		batchMethodName,
		txAmounts,
		txDestinations,
		txFees,
		big.NewInt(int64(btx.BatchNonce)),
		gethcommon.HexToAddress(btx.TokenContract),
		big.NewInt(int64(btx.Timeout)),
	}

	return packCall(BatchTxCheckpointABIJSON, "submitBatch", args)
}

// GetCheckpoint gets the checkpoint signature from the given outgoing tx batch
func (cctx ContractCallTx) GetCheckpoint(gravityID string) []byte {

	// Create the methodName argument which salts the signature
	methodNameBytes := []uint8("logicCall")
	var logicCallMethodName [32]uint8
	copy(logicCallMethodName[:], methodNameBytes[:])

	// the contract argument is not a arbitrary length array but a fixed length 32 byte
	// array, therefore we have to utf8 encode the string (the default in this case) and
	// then copy the variable length encoded data into a fixed length array. This function
	// will panic if gravityId is too long to fit in 32 bytes
	gravityIDFixed, err := strToFixByteArray(gravityID)
	if err != nil {
		panic(err)
	}

	// Run through the elements of the logic call and serialize them
	transferAmounts := make([]*big.Int, len(cctx.Tokens))
	transferTokenContracts := make([]gethcommon.Address, len(cctx.Tokens))
	feeAmounts := make([]*big.Int, len(cctx.Fees))
	feeTokenContracts := make([]gethcommon.Address, len(cctx.Fees))
	for i, coin := range cctx.Tokens {
		transferAmounts[i] = coin.Amount.BigInt()
		transferTokenContracts[i] = gethcommon.HexToAddress(coin.Contract)
	}
	for i, coin := range cctx.Fees {
		feeAmounts[i] = coin.Amount.BigInt()
		feeTokenContracts[i] = gethcommon.HexToAddress(coin.Contract)
	}
	payload := make([]byte, len(cctx.Payload))
	copy(payload, cctx.Payload)
	var invalidationId [32]byte
	copy(invalidationId[:], cctx.InvalidationScope[:])

	// the methodName needs to be the same as the 'name' above in the checkpointAbiJson
	// but other than that it's a constant that has no impact on the output. This is because
	// it gets encoded as a function name which we must then discard.
	args := []interface{}{
		gravityIDFixed,
		logicCallMethodName,
		transferAmounts,
		transferTokenContracts,
		feeAmounts,
		feeTokenContracts,
		gethcommon.HexToAddress(cctx.Address),
		payload,
		big.NewInt(int64(cctx.Timeout)),
		invalidationId,
		big.NewInt(int64(cctx.InvalidationNonce)),
	}

	return packCall(ContractCallTxABIJSON, "checkpoint", args)
}

func packCall(abiString, method string, args []interface{}) []byte {
	encodedCall, err := abi.JSON(strings.NewReader(abiString))
	if err != nil {
		panic(sdkerrors.Wrap(err, "bad ABI definition in code"))
	}
	abiEncodedCall, err := encodedCall.Pack(method, args...)
	if err != nil {
		panic(sdkerrors.Wrap(err, "packing checkpoint"))
	}
	return crypto.Keccak256Hash(abiEncodedCall[4:]).Bytes()
}
