package types

import (
	"math/big"
	"strings"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
)

// GetCheckpoint gets the checkpoint signature from the given outgoing tx batch
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
	const abiJSON = `[{
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
	abi, err := abi.JSON(strings.NewReader(abiJSON))
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
	methodNameBytes := []uint8("valsetAndTransactionBatch")
	var batchMethodName [32]uint8
	copy(batchMethodName[:], methodNameBytes[:])

	// Run through the elements of the batch and serialize them
	txAmounts := make([]*big.Int, len(b.Elements))
	txDestinations := make([]EthereumAddress, len(b.Elements))
	txFees := make([]*big.Int, len(b.Elements))
	for i, tx := range b.Elements {
		txAmounts[i] = tx.Amount.Amount.BigInt()
		txDestinations[i] = NewEthereumAddress(string(tx.DestAddress))
		txFees[i] = tx.BridgeFee.Amount.BigInt()
	}

	batchNonce := big.NewInt(int64(b.Nonce))

	valsetCheckpointBytes := (*b.Valset).GetCheckpoint()
	var valsetCheckpoint [32]uint8
	copy(valsetCheckpoint[:], valsetCheckpointBytes[:])

	// tokenContractBytes := b.TokenContract.Bytes()
	// var tokenContract [20]uint8
	// copy(tokenContract[:], tokenContractBytes[:])

	tokenContract := b.TokenContract

	// the methodName needs to be the same as the 'name' above in the checkpointAbiJson
	// but other than that it's a constant that has no impact on the output. This is because
	// it gets encoded as a function name which we must then discard.
	abiEncodedBatch, err := abi.Pack("updateValsetAndSubmitBatch",
		peggyID,
		batchMethodName,
		valsetCheckpoint,
		txAmounts,
		txDestinations,
		txFees,
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
	batchDigest := crypto.Keccak256Hash(abiEncodedBatch[4:])

	// fmt.Printf(`
	//   valsetCheckpoint: 0x%v
	//   elements in batch digest: {
	// 	peggyId: 0x%v,
	// 	batchMethodName: 0x%v,
	// 	valsetCheckpoint: 0x%v,
	// 	txAmounts: %v,
	// 	txDestinations: %v,
	// 	txFees: %v,
	// 	batchNonce: %v,
	// 	tokenContract: %v
	//   }
	//   abiEncodedBatch: 0x%v
	//   batchDigest: 0x%v
	// `,
	// 	// peggyID, validators, valsetMethodName, valsetNonce, powers,
	// 	// abiEncodedValset,
	// 	common.Bytes2Hex(valsetCheckpoint[:]),
	// 	common.Bytes2Hex(peggyID[:]), common.Bytes2Hex(batchMethodName[:]), common.Bytes2Hex(valsetCheckpoint[:]), txAmounts, txDestinations, txFees, batchNonce, tokenContract,
	// 	common.Bytes2Hex(abiEncodedBatch[:]),
	// 	common.Bytes2Hex(batchDigest[:]),
	// )

	return batchDigest.Bytes(), nil
}

// func (b *OutgoingTxBatch) Observed() error {
// 	if b.BatchStatus != BatchStatusPending && b.BatchStatus != BatchStatusSubmitted {
// 		return sdkerrors.Wrap(ErrInvalid, "status")
// 	}
// 	b.BatchStatus = BatchStatusProcessed
// 	return nil
// }

// type OutgoingTransferTx struct {
// 	ID          uint64          `json:"txid"`
// 	Sender      sdk.AccAddress  `json:"sender"`
// 	DestAddress EthereumAddress `json:"dest_address"`
// 	Amount      ERC20Token      `json:"send"`
// 	BridgeFee   ERC20Token      `json:"bridge_fee"`
// }
