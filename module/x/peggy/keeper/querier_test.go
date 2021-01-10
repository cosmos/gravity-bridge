package keeper

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/althea-net/peggy/module/x/peggy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryValsetConfirm(t *testing.T) {
	var (
		nonce                                       = uint64(1)
		myValidatorCosmosAddr, _                    = sdk.AccAddressFromBech32("cosmos1ees2tqhhhm9ahlhceh2zdguww9lqn2ckukn86l")
		myValidatorEthereumAddr  gethcommon.Address = gethcommon.BytesToAddress(bytes.Repeat([]byte{byte(50)}, 20))
	)
	input := CreateTestEnv(t)
	ctx := input.Context
	input.PeggyKeeper.SetValsetConfirm(ctx, types.MsgValsetConfirm{
		Nonce:        nonce,
		Orchestrator: myValidatorCosmosAddr.String(),
		EthAddress:   myValidatorEthereumAddr.String(),
		Signature:    "alksdjhflkasjdfoiasjdfiasjdfoiasdj",
	})

	specs := map[string]struct {
		srcNonce string
		srcAddr  string
		expErr   bool
		expResp  []byte
	}{
		"all good": {
			srcNonce: "1",
			srcAddr:  myValidatorCosmosAddr.String(),
			expResp:  []byte(`{"type":"peggy/MsgValsetConfirm", "value":{"eth_address":"0x3232323232323232323232323232323232323232", "nonce": "1", "orchestrator": "cosmos1ees2tqhhhm9ahlhceh2zdguww9lqn2ckukn86l",  "signature": "alksdjhflkasjdfoiasjdfiasjdfoiasdj"}}`),
		},
		"unknown nonce": {
			srcNonce: "999999",
			srcAddr:  myValidatorCosmosAddr.String(),
		},
		"invalid address": {
			srcNonce: "1",
			srcAddr:  "not a valid addr",
			expErr:   true,
		},
		"invalid nonce": {
			srcNonce: "not a valid nonce",
			srcAddr:  myValidatorCosmosAddr.String(),
			expErr:   true,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			got, err := queryValsetConfirm(ctx, []string{spec.srcNonce, spec.srcAddr}, input.PeggyKeeper)
			if spec.expErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if spec.expResp == nil {
				assert.Nil(t, got)
				return
			}
			assert.JSONEq(t, string(spec.expResp), string(got))
		})
	}
}

func TestAllValsetConfirmsBynonce(t *testing.T) {
	input := CreateTestEnv(t)
	ctx := input.Context

	addrs := []string{
		"cosmos1u508cfnsk2nhakv80vdtq3nf558ngyvldkfjj9",
		"cosmos1krtcsrxhadj54px0vy6j33pjuzcd3jj8kmsazv",
		"cosmos1u94xef3cp9thkcpxecuvhtpwnmg8mhlja8hzkd",
	}
	// seed confirmations
	for i := 0; i < 3; i++ {
		addr, _ := sdk.AccAddressFromBech32(addrs[i])
		msg := types.MsgValsetConfirm{}
		msg.EthAddress = gethcommon.BytesToAddress(bytes.Repeat([]byte{byte(i + 1)}, 20)).String()
		msg.Nonce = uint64(1)
		msg.Orchestrator = addr.String()
		msg.Signature = fmt.Sprintf("signature %d", i+1)
		input.PeggyKeeper.SetValsetConfirm(ctx, msg)
	}

	specs := map[string]struct {
		srcNonce string
		expErr   bool
		expResp  []byte
	}{
		"all good": {
			srcNonce: "1",
			expResp: []byte(`[
      {"eth_address":"0x0202020202020202020202020202020202020202", "nonce": "1", "orchestrator": "cosmos1krtcsrxhadj54px0vy6j33pjuzcd3jj8kmsazv", "signature": "signature 2"},
	  {"eth_address":"0x0303030303030303030303030303030303030303", "nonce": "1", "orchestrator": "cosmos1u94xef3cp9thkcpxecuvhtpwnmg8mhlja8hzkd", "signature": "signature 3"},
	  {"eth_address":"0x0101010101010101010101010101010101010101", "nonce": "1", "orchestrator": "cosmos1u508cfnsk2nhakv80vdtq3nf558ngyvldkfjj9", "signature": "signature 1"}
]`),
		},
		"unknown nonce": {
			srcNonce: "999999",
			expResp:  nil,
		},
		"invalid nonce": {
			srcNonce: "not a valid nonce",
			expErr:   true,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			got, err := queryAllValsetConfirms(ctx, spec.srcNonce, input.PeggyKeeper)
			if spec.expErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if spec.expResp == nil {
				assert.Nil(t, got)
				return
			}
			assert.JSONEq(t, string(spec.expResp), string(got))
		})
	}
}

// TODO: Check failure modes
func TestLastValsetRequests(t *testing.T) {
	input := CreateTestEnv(t)
	ctx := input.Context
	// seed with requests
	for i := 0; i < 6; i++ {
		var validators []sdk.ValAddress
		for j := 0; j <= i; j++ {
			// add an validator each block
			valAddr := bytes.Repeat([]byte{byte(j)}, sdk.AddrLen)
			input.PeggyKeeper.SetEthAddress(ctx, valAddr, gethcommon.BytesToAddress(bytes.Repeat([]byte{byte(j + 1)}, 20)).String())
			validators = append(validators, valAddr)
		}
		input.PeggyKeeper.StakingKeeper = NewStakingKeeperMock(validators...)
		ctx = ctx.WithBlockHeight(int64(100 + i))
		input.PeggyKeeper.SetValsetRequest(ctx)
	}

	specs := map[string]struct {
		expResp []byte
	}{
		"limit at 5": {
			expResp: []byte(`[
{
  "nonce": "105",
  "height": "105",
  "members": [
    {
      "power": "715827882",
      "ethereum_address": "0x0101010101010101010101010101010101010101"
    },
    {
      "power": "715827882",
      "ethereum_address": "0x0202020202020202020202020202020202020202"
    },
    {
      "power": "715827882",
      "ethereum_address": "0x0303030303030303030303030303030303030303"
    },
    {
      "power": "715827882",
      "ethereum_address": "0x0404040404040404040404040404040404040404"
    },
    {
      "power": "715827882",
      "ethereum_address": "0x0505050505050505050505050505050505050505"
    },
    {
      "power": "715827882",
      "ethereum_address": "0x0606060606060606060606060606060606060606"
    }
  ]
},
{
  "nonce": "104",
  "height": "104",
  "members": [
    {
      "power": "858993459",
      "ethereum_address": "0x0101010101010101010101010101010101010101"
    },
    {
      "power": "858993459",
      "ethereum_address": "0x0202020202020202020202020202020202020202"
    },
    {
      "power": "858993459",
      "ethereum_address": "0x0303030303030303030303030303030303030303"
    },
    {
      "power": "858993459",
      "ethereum_address": "0x0404040404040404040404040404040404040404"
    },
    {
      "power": "858993459",
      "ethereum_address": "0x0505050505050505050505050505050505050505"
    }
  ]
},
{
  "nonce": "103",
  "height": "103",
  "members": [
    {
      "power": "1073741823",
      "ethereum_address": "0x0101010101010101010101010101010101010101"
    },
    {
      "power": "1073741823",
      "ethereum_address": "0x0202020202020202020202020202020202020202"
    },
    {
      "power": "1073741823",
      "ethereum_address": "0x0303030303030303030303030303030303030303"
    },
    {
      "power": "1073741823",
      "ethereum_address": "0x0404040404040404040404040404040404040404"
    }
  ]
},
{
  "nonce": "102",
  "height": "102",
  "members": [
    {
      "power": "1431655765",
      "ethereum_address": "0x0101010101010101010101010101010101010101"
    },
    {
      "power": "1431655765",
      "ethereum_address": "0x0202020202020202020202020202020202020202"
    },
    {
      "power": "1431655765",
      "ethereum_address": "0x0303030303030303030303030303030303030303"
    }
  ]
},
{
  "nonce": "101",
  "height": "101",
  "members": [
    {
      "power": "2147483647",
      "ethereum_address": "0x0101010101010101010101010101010101010101"
    },
    {
      "power": "2147483647",
      "ethereum_address": "0x0202020202020202020202020202020202020202"
    }
  ]
}
]`),
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			got, err := lastValsetRequests(ctx, input.PeggyKeeper)
			require.NoError(t, err)
			assert.JSONEq(t, string(spec.expResp), string(got), string(got))
		})
	}
}

// TODO: check that it doesn't accidently return a valset that HAS been signed
// Right now it is basically just testing that any valset comes back
func TestPendingValsetRequests(t *testing.T) {
	input := CreateTestEnv(t)
	ctx := input.Context

	// seed with requests
	for i := 0; i < 6; i++ {
		var validators []sdk.ValAddress
		for j := 0; j <= i; j++ {
			// add an validator each block
			valAddr := bytes.Repeat([]byte{byte(j)}, sdk.AddrLen)
			input.PeggyKeeper.SetEthAddress(ctx, valAddr, gethcommon.BytesToAddress(bytes.Repeat([]byte{byte(j + 1)}, 20)).String())
			validators = append(validators, valAddr)
		}
		input.PeggyKeeper.StakingKeeper = NewStakingKeeperMock(validators...)
		ctx = ctx.WithBlockHeight(int64(100 + i))
		input.PeggyKeeper.SetValsetRequest(ctx)
	}

	specs := map[string]struct {
		expResp []byte
	}{
		"find valset": {
			expResp: []byte(`{
      "type": "peggy/Valset",
      "value": {
        "nonce": "105",
        "height": "105",
        "members": [
          {
            "power": "715827882",
            "ethereum_address": "0x0101010101010101010101010101010101010101"
          },
          {
            "power": "715827882",
            "ethereum_address": "0x0202020202020202020202020202020202020202"
          },
          {
            "power": "715827882",
            "ethereum_address": "0x0303030303030303030303030303030303030303"
          },
          {
            "power": "715827882",
            "ethereum_address": "0x0404040404040404040404040404040404040404"
          },
          {
            "power": "715827882",
            "ethereum_address": "0x0505050505050505050505050505050505050505"
          },
          {
            "power": "715827882",
            "ethereum_address": "0x0606060606060606060606060606060606060606"
          }
        ]
      }
      }`),
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			valAddr := sdk.AccAddress{}
			valAddr = bytes.Repeat([]byte{byte(1)}, sdk.AddrLen)
			got, err := lastPendingValsetRequest(ctx, valAddr.String(), input.PeggyKeeper)
			require.NoError(t, err)
			assert.JSONEq(t, string(spec.expResp), string(got), string(got))
		})
	}
}

// TODO: check that it actually returns the valset that has NOT been signed, not just any valset
func TestLastPendingBatchRequest(t *testing.T) {
	input := CreateTestEnv(t)
	ctx := input.Context

	// seed with valset requests and eth addresses to make validators
	// that we will later use to lookup batches to be signed
	for i := 0; i < 6; i++ {
		var validators []sdk.ValAddress
		for j := 0; j <= i; j++ {
			// add an validator each block
			// TODO: replace with real SDK addresses
			valAddr := bytes.Repeat([]byte{byte(j)}, sdk.AddrLen)
			input.PeggyKeeper.SetEthAddress(ctx, valAddr, gethcommon.BytesToAddress(bytes.Repeat([]byte{byte(j + 1)}, 20)).String())
			validators = append(validators, valAddr)
		}
		input.PeggyKeeper.StakingKeeper = NewStakingKeeperMock(validators...)
		input.PeggyKeeper.SetValsetRequest(ctx)
	}

	createTestBatch(t, input)

	specs := map[string]struct {
		expResp []byte
	}{
		"find batch": {
			expResp: []byte(`{
	"type": "peggy/OutgoingTxBatch",
	"value": {
	"batch_nonce": "1",
	"block": "1234567",
	"transactions": [
		{
		"id": "2",
		"sender": "cosmos1qyqszqgpqyqszqgpqyqszqgpqyqszqgpjnp7du",
		"dest_address": "0x320915BD0F1bad11cBf06e85D5199DBcAC4E9934",
		"erc20_token": {
			"amount": "101",
			"contract": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
		},
		"erc20_fee": {
			"amount": "3",
			"contract": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
		}
		},
		{
		"id": "1",
		"sender": "cosmos1qyqszqgpqyqszqgpqyqszqgpqyqszqgpjnp7du",
		"dest_address": "0x320915BD0F1bad11cBf06e85D5199DBcAC4E9934",
		"erc20_token": {
			"amount": "100",
			"contract": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
		},
		"erc20_fee": {
			"amount": "2",
			"contract": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
		}
		}
	],
	"token_contract": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
	}
}
			`,
			)},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			valAddr := sdk.AccAddress{}
			valAddr = bytes.Repeat([]byte{byte(1)}, sdk.AddrLen)
			got, err := lastPendingBatchRequest(ctx, valAddr.String(), input.PeggyKeeper)
			require.NoError(t, err)
			assert.JSONEq(t, string(spec.expResp), string(got), string(got))
		})
	}
}

func createTestBatch(t *testing.T, input TestInput) {
	var (
		mySender            = bytes.Repeat([]byte{1}, sdk.AddrLen)
		myReceiver          = "0x320915BD0F1bad11cBf06e85D5199DBcAC4E9934"
		myTokenContractAddr = "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
		now                 = time.Now().UTC()
	)
	// mint some voucher first
	allVouchers := sdk.Coins{types.NewERC20Token(99999, myTokenContractAddr).PeggyCoin()}
	err := input.BankKeeper.MintCoins(input.Context, types.ModuleName, allVouchers)
	require.NoError(t, err)

	// set senders balance
	input.AccountKeeper.NewAccountWithAddress(input.Context, mySender)
	err = input.BankKeeper.SetBalances(input.Context, mySender, allVouchers)
	require.NoError(t, err)

	// add some TX to the pool
	for i, v := range []uint64{2, 3, 2, 1} {
		amount := types.NewERC20Token(uint64(i+100), myTokenContractAddr).PeggyCoin()
		fee := types.NewERC20Token(v, myTokenContractAddr).PeggyCoin()
		_, err := input.PeggyKeeper.AddToOutgoingPool(input.Context, mySender, myReceiver, amount, fee)
		require.NoError(t, err)
	}
	// when
	input.Context = input.Context.WithBlockTime(now)

	// tx batch size is 2, so that some of them stay behind
	_, err = input.PeggyKeeper.BuildOutgoingTXBatch(input.Context, myTokenContractAddr, 2)
	require.NoError(t, err)
}

func TestQueryAllBatchConfirms(t *testing.T) {
	input := CreateTestEnv(t)
	ctx := input.Context

	var (
		tokenContract    = "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
		validatorAddr, _ = sdk.AccAddressFromBech32("cosmos1mgamdcs9dah0vn0gqupl05up7pedg2mvupe6hh")
	)

	input.PeggyKeeper.SetBatchConfirm(ctx, &types.MsgConfirmBatch{
		Nonce:         1,
		TokenContract: tokenContract,
		EthSigner:     "0xf35e2cc8e6523d683ed44870f5b7cc785051a77d",
		Orchestrator:  validatorAddr.String(),
		Signature:     "signature",
	})

	batchConfirms, err := queryAllBatchConfirms(ctx, "1", tokenContract, input.PeggyKeeper)
	require.NoError(t, err)

	expectedJSON := []byte(`[{"eth_signer":"0xf35e2cc8e6523d683ed44870f5b7cc785051a77d", "nonce":"1", "signature":"signature", "token_contract":"0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B", "orchestrator":"cosmos1mgamdcs9dah0vn0gqupl05up7pedg2mvupe6hh"}]`)

	assert.JSONEq(t, string(expectedJSON), string(batchConfirms), "json is equal")
}

// TODO: test that it gets the correct batch, not just any batch.
// Check with multiple nonces and tokenContracts
func TestQueryBatch(t *testing.T) {
	input := CreateTestEnv(t)
	ctx := input.Context

	var (
		tokenContract = "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
	)

	createTestBatch(t, input)

	batch, err := queryBatch(ctx, "1", tokenContract, input.PeggyKeeper)
	require.NoError(t, err)

	expectedJSON := []byte(`{
		"type": "peggy/OutgoingTxBatch",
		"value": {
		  "transactions": [
			{
			  "erc20_fee": {
				"amount": "3",
				"contract": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
			  },
			  "dest_address": "0x320915BD0F1bad11cBf06e85D5199DBcAC4E9934",
			  "erc20_token": {
				"amount": "101",
				"contract": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
			  },
			  "sender": "cosmos1qyqszqgpqyqszqgpqyqszqgpqyqszqgpjnp7du",
			  "id": "2"
			},
			{
			  "erc20_fee": {
				"amount": "2",
				"contract": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
			  },
			  "dest_address": "0x320915BD0F1bad11cBf06e85D5199DBcAC4E9934",
			  "erc20_token": {
				"amount": "100",
				"contract": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
			  },
			  "sender": "cosmos1qyqszqgpqyqszqgpqyqszqgpqyqszqgpjnp7du",
			  "id": "1"
			}
		  ],
		  "batch_nonce": "1",
		  "block": "1234567",
		  "token_contract": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
		}
	  }
	  `)

	// TODO: this test is failing on the empty representation of valset members
	assert.JSONEq(t, string(expectedJSON), string(batch), string(batch))
}

func TestLastBatchesRequest(t *testing.T) {
	input := CreateTestEnv(t)
	ctx := input.Context

	createTestBatch(t, input)
	createTestBatch(t, input)

	lastBatches, err := lastBatchesRequest(ctx, input.PeggyKeeper)
	require.NoError(t, err)

	expectedJSON := []byte(`[
		{
		  "transactions": [
			{
			  "erc20_fee": {
				"amount": "3",
				"contract": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
			  },
			  "dest_address": "0x320915BD0F1bad11cBf06e85D5199DBcAC4E9934",
			  "erc20_token": {
				"amount": "101",
				"contract": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
			  },
			  "sender": "cosmos1qyqszqgpqyqszqgpqyqszqgpqyqszqgpjnp7du",
			  "id": "6"
			},
			{
			  "erc20_fee": {
				"amount": "2",
				"contract": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
			  },
			  "dest_address": "0x320915BD0F1bad11cBf06e85D5199DBcAC4E9934",
			  "erc20_token": {
				"amount": "102",
				"contract": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
			  },
			  "sender": "cosmos1qyqszqgpqyqszqgpqyqszqgpqyqszqgpjnp7du",
			  "id": "3"
			}
		  ],
		  "batch_nonce": "2",
		  "block": "1234567",
		  "token_contract": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
		},
		{
		  "transactions": [
			{
			  "erc20_fee": {
				"amount": "3",
				"contract": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
			  },
			  "dest_address": "0x320915BD0F1bad11cBf06e85D5199DBcAC4E9934",
			  "erc20_token": {
				"amount": "101",
				"contract": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
			  },
			  "sender": "cosmos1qyqszqgpqyqszqgpqyqszqgpqyqszqgpjnp7du",
			  "id": "2"
			},
			{
			  "erc20_fee": {
				"amount": "2",
				"contract": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
			  },
			  "dest_address": "0x320915BD0F1bad11cBf06e85D5199DBcAC4E9934",
			  "erc20_token": {
				"amount": "100",
				"contract": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
			  },
			  "sender": "cosmos1qyqszqgpqyqszqgpqyqszqgpqyqszqgpjnp7du",
			  "id": "1"
			}
		  ],
		  "batch_nonce": "1",
		  "block": "1234567",
		  "token_contract": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
		}
	  ]
	  `)

	assert.JSONEq(t, string(expectedJSON), string(lastBatches), "json is equal")
}

// tests setting and querying eth address and orchestrator addresses
func TestQueryCurrentValset(t *testing.T) {
	var (
		ethAddress                = "0xb462864E395d88d6bc7C5dd5F3F5eb4cc2599255"
		valAddress sdk.ValAddress = bytes.Repeat([]byte{0x2}, sdk.AddrLen)
	)
	input := CreateTestEnv(t)
	input.PeggyKeeper.StakingKeeper = NewStakingKeeperMock(valAddress)
	ctx := input.Context
	input.PeggyKeeper.SetEthAddress(ctx, valAddress, ethAddress)

	currentValset := input.PeggyKeeper.GetCurrentValset(ctx)

	bridgeVal := types.BridgeValidator{EthereumAddress: ethAddress, Power: 4294967295}
	expectedValset := types.Valset{Nonce: 1234567, Height: 1234567, Members: []*types.BridgeValidator{&bridgeVal}}
	assert.Equal(t, &expectedValset, currentValset)
}
