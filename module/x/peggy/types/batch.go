package types

import (
	"math/big"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
)

type BatchStatus uint8

const (
	BatchStatusUnknown   BatchStatus = 0
	BatchStatusPending   BatchStatus = 1 // initial status
	BatchStatusSubmitted BatchStatus = 2 // in flight to ETH
	BatchStatusProcessed BatchStatus = 3 // observed - end state
	BatchStatusCancelled BatchStatus = 4 // end state
)

func (b BatchStatus) String() string {
	return []string{"unknown", "pending", "submitted", "observed", "processed", "cancelled"}[b]
}

type OutgoingTxBatch struct {
	Nonce              UInt64Nonce          `json:"nonce"`
	Elements           []OutgoingTransferTx `json:"elements"`
	CreatedAt          time.Time            `json:"created_at"`
	TotalFee           ERC20Token           `json:"total_fee"`
	BridgedDenominator BridgedDenominator   `json:"bridged_denominator"`
	BatchStatus        BatchStatus          `json:"batch_status"`
	Valset             Valset               `json:"valset"`
	TokenContract      EthereumAddress      `json:"tokenContract"`
}

func (b *OutgoingTxBatch) Cancel() error {
	if b.BatchStatus != BatchStatusPending {
		return sdkerrors.Wrap(ErrInvalid, "status - batch not pending")
	}
	b.BatchStatus = BatchStatusCancelled
	return nil
}

func (b OutgoingTxBatch) GetCheckpoint() ([]byte, error) {

	// TODO replace hardcoded "foo" here with a getter to retrieve the correct PeggyID from the store
	// this will work for now because 'foo' is the test Peggy ID we are using
	var peggyIDString = "foo"

	// The go-ethereum ABI encoder *only* encodes function calls and then it only encodes
	// function calls for which you provide an ABI json just like you would get out of the
	// solidity compiler with your compiled contract.
	// You are supposed to compile your contract, use abigen to generate an ABI , import
	// this generated go module and then use for that for all testing and development.
	// This abstraction layer is more trouble than it's worth, because we don't want to
	// encode a function call at all, but instead we want to emulate a Solidity encode operation
	// which has no equal available from go-ethereum.
	//
	// In order to work around this absurd series of problems we have to manually write the below
	// 'function specification' that will encode the same arguments into a function call. We can then
	// truncate the first several bytes where the call name is encoded to finally get the equal of the
	const checkpointAbiJSON = `[{
	  "inputs": [
	    {
	      "internalType": "bytes32",
	      "name": "_peggyId",
	      "type": "bytes32"
	    },
	    {
	      "internalType": "bytes32",
	      "name": "_methodName",
	      "type": "bytes32"
		},
		{
		  "internalType": "bytes32",
		  "name": "_checkPoint",
		  "type": "bytes32"
		},
		
	    {
	      "internalType": "uint256[]",
	      "name": "_amounts",
	      "type": "uint256[]"
	    },
	    {
	      "internalType": "address[]",
	      "name": "_destinations",
	      "type": "address[]"
	    },
	    {
	      "internalType": "uint256[]",
	      "name": "_fees",
	      "type": "uint256[]"
	    },

		{
		  "internalType": "uint256",
		  "name": "_batchNonce",
		  "type": "uint256"
		},
		{
		  "internalType": "address",
		  "name": "_tokenContract",
		  "type": "address"
		}
	  ],
	  "name": "updateValsetAndSubmitBatch",
	  "outputs": [
	    {
	      "internalType": "bytes32",
	      "name": "",
	      "type": "bytes32"
	    }
	  ],
	  "stateMutability": "pure",
	  "type": "function"
	}]`
	// Solidity abi.Encode() call.
	// error case here should not occur outside of testing since the above is a constant
	contractAbi, err := abi.JSON(strings.NewReader(checkpointAbiJSON))
	if err != nil {
		return nil, sdkerrors.Wrap(err, "bad ABI definition in code")
	}
	peggyIDBytes := []uint8(peggyIDString)

	// the contract argument is not a arbitrary length array but a fixed length 32 byte
	// array, therefore we have to utf8 encode the string (the default in this case) and
	// then copy the variable length encoded data into a fixed length array. This function
	// will panic if peggyId is too long to fit in 32 bytes
	var peggyID [32]uint8
	copy(peggyID[:], peggyIDBytes[:])

	// Create the methodName argument which salts the signature
	methodNameBytes := []uint8("updateValsetAndSubmitBatch")
	var methodName [32]uint8
	copy(methodName[:], methodNameBytes[:])

	// Run through the elements of the batch and serialize them
	amounts := make([]*big.Int, len(b.Elements))
	destinations := make([]EthereumAddress, len(b.Elements))
	fees := make([]*big.Int, len(b.Elements))
	for i, tx := range b.Elements {
		amounts[i] = big.NewInt(int64(tx.Amount.Amount))
		destinations[i] = tx.DestAddress
		fees[i] = big.NewInt(int64(tx.BridgeFee.Amount))
	}

	batchNonce := big.NewInt(int64(b.Nonce))

	valsetCheckpointBytes := b.Valset.GetCheckpoint()
	var valsetCheckpoint [32]uint8
	copy(valsetCheckpoint[:], valsetCheckpointBytes[:])

	tokenContractBytes := b.TokenContract.Bytes()
	var tokenContract [32]uint8
	copy(tokenContract[:], tokenContractBytes[:])

	// the methodName needs to be the same as the 'name' above in the checkpointAbiJson
	// but other than that it's a constant that has no impact on the output. This is because
	// it gets encoded as a function name which we must then discard.
	bytes, err := contractAbi.Pack("updateValsetAndSubmitBatch",
		peggyID,
		methodName,
		valsetCheckpoint,
		amounts,
		destinations,
		fees,
		batchNonce,
		tokenContract,
	)
	// this should never happen outside of test since any case that could crash on encoding
	// should be filtered above.
	if err != nil {
		return nil, sdkerrors.Wrap(err, "packing checkpoint")
	}

	// we hash the resulting encoded bytes discarding the first 4 bytes these 4 bytes are the constant
	// method name 'checkpoint'. If you where to replace the checkpoint constant in this code you would
	// then need to adjust how many bytes you truncate off the front to get the output of abi.encode()
	hash := crypto.Keccak256Hash(bytes[4:])

	return hash.Bytes(), nil
}

func (b *OutgoingTxBatch) Observed() error {
	if b.BatchStatus != BatchStatusPending && b.BatchStatus != BatchStatusSubmitted {
		return sdkerrors.Wrap(ErrInvalid, "status")
	}
	b.BatchStatus = BatchStatusProcessed
	return nil
}

type OutgoingTransferTx struct {
	ID          uint64          `json:"txid"`
	Sender      sdk.AccAddress  `json:"sender"`
	DestAddress EthereumAddress `json:"dest_address"`
	Amount      ERC20Token      `json:"send"`
	BridgeFee   ERC20Token      `json:"bridge_fee"`
}
