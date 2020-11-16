package types

import (
	"encoding/hex"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOutgoingTxBatchCheckpointGold1(t *testing.T) {
	senderAddr, err := sdk.AccAddressFromHex("527FBEE652609AB150F0AEE9D61A2F76CFC4A73E")
	require.NoError(t, err)
	var (
		erc20Addr = NewEthereumAddress("0x22474D350EC2dA53D717E30b96e9a2B7628Ede5b")
	)

	v := NewValset(
		NewUInt64Nonce(1),
		BridgeValidators{{
			EthereumAddress: NewEthereumAddress("0xc783df8a850f42e7F7e57013759C285caa701eB6").Bytes(),
			Power:           6670,
		}},
	)

	src := OutgoingTxBatch{
		Nonce: 1,
		Elements: []*OutgoingTransferTx{
			{
				Id:          0x1,
				Sender:      senderAddr.String(),
				DestAddress: NewEthereumAddress("0x9FC9C2DfBA3b6cF204C37a5F690619772b926e39").Bytes(),
				Amount: &ERC20Token{
					Amount:               sdk.NewInt(0x1),
					Symbol:               "MAX",
					TokenContractAddress: erc20Addr.Bytes(),
				},
				BridgeFee: &ERC20Token{
					Amount:               sdk.NewInt(0x1),
					Symbol:               "MAX",
					TokenContractAddress: erc20Addr.Bytes(),
				},
			},
		},
		TotalFee: &ERC20Token{
			Amount:               sdk.NewInt(0x1),
			Symbol:               "MAX",
			TokenContractAddress: erc20Addr.Bytes(),
		},
		BridgedDenominator: &BridgedDenominator{
			TokenContractAddress: erc20Addr.Bytes(),
			Symbol:               "MAX",
			CosmosVoucherDenom:   "peggy39b512461b",
		},
		Valset:        &v,
		TokenContract: erc20Addr.Bytes(),
	}

	ourHash, err := src.GetCheckpoint()
	require.NoError(t, err)

	// hash from bridge contract
	goldHash := "0x746471abc2232c11039c2160365c4593110dbfbe25ff9a2dcf8b5b7376e9f346"[2:]
	assert.Equal(t, goldHash, hex.EncodeToString(ourHash))
}

// This is the output from the function used to compute the "gold hash" above. It is useful to check for problems.
// The code that produces this output is in /solidity/test/updateValsetAndSubmitBatch.ts
// Be aware that every time that you run the above .ts file, it will use a different tokenContractAddress and thus compute
// a different hash.

// elements in valset checkpoint: {
//   peggyId: '0x666f6f0000000000000000000000000000000000000000000000000000000000',
//   validators: [ '0xc783df8a850f42e7F7e57013759C285caa701eB6' ],
//   valsetMethodName: '0x636865636b706f696e7400000000000000000000000000000000000000000000',
//   valsetNonce: 1,
//   powers: [ 6670 ]
// }
// abiEncodedValset: 0x666f6f0000000000000000000000000000000000000000000000000000000000636865636b706f696e7400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000e00000000000000000000000000000000000000000000000000000000000000001000000000000000000000000c783df8a850f42e7f7e57013759c285caa701eb600000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000001a0e
// valsetCheckpoint: 0xb85e6cb8ee0858a9747229dcf931ca95a9a2ef5b83601592759a65d2c2a30876
// elements in batch cehckpoint: {
//   peggyId: '0x666f6f0000000000000000000000000000000000000000000000000000000000',
//   batchMethodName: '0x76616c736574416e645472616e73616374696f6e426174636800000000000000',
//   valsetCheckpoint: '0xb85e6cb8ee0858a9747229dcf931ca95a9a2ef5b83601592759a65d2c2a30876',
//   txAmounts: [ 1 ],
//   txDestinations: [ '0x9FC9C2DfBA3b6cF204C37a5F690619772b926e39' ],
//   txFees: [ 1 ],
//   batchNonce: 1,
//   tokenContract: '0x22474D350EC2dA53D717E30b96e9a2B7628Ede5b'
// }
// abiEncodedBatch: 0x666f6f000000000000000000000000000000000000000000000000000000000076616c736574416e645472616e73616374696f6e426174636800000000000000b85e6cb8ee0858a9747229dcf931ca95a9a2ef5b83601592759a65d2c2a30876000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001400000000000000000000000000000000000000000000000000000000000000180000000000000000000000000000000000000000000000000000000000000000100000000000000000000000022474d350ec2da53d717e30b96e9a2b7628ede5b0000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000010000000000000000000000009fc9c2dfba3b6cf204c37a5f690619772b926e3900000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001
// batchCheckpoint: 0x746471abc2232c11039c2160365c4593110dbfbe25ff9a2dcf8b5b7376e9f346
