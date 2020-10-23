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
		erc20Addr = NewEthereumAddress("0x20Ce94F404343aD2752A2D01b43fa407db9E0D00")
	)

	v := NewValset(
		NewUInt64Nonce(1),
		BridgeValidators{{
			EthereumAddress: NewEthereumAddress("0xc783df8a850f42e7F7e57013759C285caa701eB6"),
			Power:           6670,
		}},
	)

	src := OutgoingTxBatch{
		Nonce: 1,
		Elements: []OutgoingTransferTx{
			{
				ID:          0x1,
				Sender:      senderAddr,
				DestAddress: NewEthereumAddress("0x9FC9C2DfBA3b6cF204C37a5F690619772b926e39"),
				Amount: ERC20Token{
					Amount:               0x1,
					Symbol:               "MAX",
					TokenContractAddress: erc20Addr,
				},
				BridgeFee: ERC20Token{
					Amount:               0x1,
					Symbol:               "MAX",
					TokenContractAddress: erc20Addr,
				},
			},
		},
		TotalFee: ERC20Token{
			Amount:               0x1,
			Symbol:               "MAX",
			TokenContractAddress: erc20Addr,
		},
		BridgedDenominator: BridgedDenominator{
			TokenContractAddress: erc20Addr,
			Symbol:               "MAX",
			CosmosVoucherDenom:   "peggy39b512461b",
		},
		BatchStatus: 1,
		Valset:      v,
	}

	ourHash, err := src.GetCheckpoint()
	require.NoError(t, err)

	// hash from bridge contract
	goldHash := "0xd443c164d8456cd774688c337a11eeb8e6661d45860bd9784557ce56d3e3ea57"[2:]
	assert.Equal(t, goldHash, hex.EncodeToString(ourHash))
}
