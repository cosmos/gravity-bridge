package keeper

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/althea-net/peggy/module/x/peggy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func TestQueryValsetConfirm(t *testing.T) {
	var (
		nonce                                          = types.NewUInt64Nonce(1)
		myValidatorCosmosAddr, _                       = sdk.AccAddressFromBech32("cosmos1ees2tqhhhm9ahlhceh2zdguww9lqn2ckukn86l")
		myValidatorEthereumAddr  types.EthereumAddress = createEthAddress(50)
	)
	k, ctx, _ := CreateTestEnv(t)
	k.SetValsetConfirm(ctx, types.MsgValsetConfirm{
		Nonce:      nonce.Uint64(),
		Validator:  myValidatorCosmosAddr.String(),
		EthAddress: myValidatorEthereumAddr.String(),
		Signature:  "alksdjhflkasjdfoiasjdfiasjdfoiasdj",
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
			expResp:  []byte(`{"type":"peggy/MsgValsetConfirm", "value":{"eth_address":"0x3232323232323232323232323232323232323232", "nonce": "1", "validator": "cosmos1ees2tqhhhm9ahlhceh2zdguww9lqn2ckukn86l",  "signature": "alksdjhflkasjdfoiasjdfiasjdfoiasdj"}}`),
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
			got, err := queryValsetConfirm(ctx, []string{spec.srcNonce, spec.srcAddr}, k)
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
	k, ctx, _ := CreateTestEnv(t)

	addrs := []string{
		"cosmos1u508cfnsk2nhakv80vdtq3nf558ngyvldkfjj9",
		"cosmos1krtcsrxhadj54px0vy6j33pjuzcd3jj8kmsazv",
		"cosmos1u94xef3cp9thkcpxecuvhtpwnmg8mhlja8hzkd",
	}
	// seed confirmations
	for i := 0; i < 3; i++ {
		addr, _ := sdk.AccAddressFromBech32(addrs[i])
		msg := types.MsgValsetConfirm{}
		msg.EthAddress = createEthAddress(i + 1).String()
		msg.Nonce = uint64(1)
		msg.Validator = addr.String()
		msg.Signature = fmt.Sprintf("signature %d", i+1)
		k.SetValsetConfirm(ctx, msg)
	}

	specs := map[string]struct {
		srcNonce string
		expErr   bool
		expResp  []byte
	}{
		"all good": {
			srcNonce: "1",
			expResp: []byte(`[
      {"eth_address":"0x0101010101010101010101010101010101010101", "nonce": "1", "validator": "cosmos1u508cfnsk2nhakv80vdtq3nf558ngyvldkfjj9", "signature": "signature 1"},
      {"eth_address":"0x0202020202020202020202020202020202020202", "nonce": "1", "validator": "cosmos1krtcsrxhadj54px0vy6j33pjuzcd3jj8kmsazv", "signature": "signature 2"},
      {"eth_address":"0x0303030303030303030303030303030303030303", "nonce": "1", "validator": "cosmos1u94xef3cp9thkcpxecuvhtpwnmg8mhlja8hzkd", "signature": "signature 3"}
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
			got, err := queryAllValsetConfirms(ctx, spec.srcNonce, k)
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
	k, ctx, _ := CreateTestEnv(t)
	// seed with requests
	for i := 0; i < 6; i++ {
		var validators []sdk.ValAddress
		for j := 0; j <= i; j++ {
			// add an validator each block
			valAddr := bytes.Repeat([]byte{byte(j)}, sdk.AddrLen)
			k.SetEthAddress(ctx, valAddr, createEthAddress(j+1))
			validators = append(validators, valAddr)
		}
		k.StakingKeeper = NewStakingKeeperMock(validators...)
		ctx = ctx.WithBlockHeight(int64(100 + i))
		k.SetValsetRequest(ctx)
	}

	specs := map[string]struct {
		expResp []byte
	}{
		"limit at 5": {
			expResp: []byte(`[
{
  "nonce": "105",
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
			got, err := lastValsetRequests(ctx, k)
			require.NoError(t, err)
			assert.JSONEq(t, string(spec.expResp), string(got), string(got))
		})
	}
}

// TODO: check that it doesn't accidently return a valset that HAS been signed
// Right now it is basically just testing that any valset comes back
func TestPendingValsetRequests(t *testing.T) {
	k, ctx, _ := CreateTestEnv(t)
	// seed with requests
	for i := 0; i < 6; i++ {
		var validators []sdk.ValAddress
		for j := 0; j <= i; j++ {
			// add an validator each block
			valAddr := bytes.Repeat([]byte{byte(j)}, sdk.AddrLen)
			k.SetEthAddress(ctx, valAddr, createEthAddress(j+1))
			validators = append(validators, valAddr)
		}
		k.StakingKeeper = NewStakingKeeperMock(validators...)
		ctx = ctx.WithBlockHeight(int64(100 + i))
		k.SetValsetRequest(ctx)
	}

	specs := map[string]struct {
		expResp []byte
	}{
		"find valset": {
			expResp: []byte(`{
      "type": "peggy/Valset",
      "value": {
        "nonce": "105",
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
			got, err := lastPendingValsetRequest(ctx, valAddr.String(), k)
			require.NoError(t, err)
			assert.JSONEq(t, string(spec.expResp), string(got), string(got))
		})
	}
}

// TODO: check that it actually returns the valset that has NOT been signed, not just any valset
func TestLastPendingBatchRequest(t *testing.T) {
	k, ctx, keepers := CreateTestEnv(t)

	// seed with valset requests and eth addresses to make validators
	// that we will later use to lookup batches to be signed
	for i := 0; i < 6; i++ {
		var validators []sdk.ValAddress
		for j := 0; j <= i; j++ {
			// add an validator each block
			valAddr := bytes.Repeat([]byte{byte(j)}, sdk.AddrLen)
			k.SetEthAddress(ctx, valAddr, createEthAddress(j+1))
			validators = append(validators, valAddr)
		}
		k.StakingKeeper = NewStakingKeeperMock(validators...)
		ctx = ctx.WithBlockHeight(int64(100 + i))
		k.SetValsetRequest(ctx)
	}

	createTestBatch(t, k, ctx, keepers)

	specs := map[string]struct {
		expResp []byte
	}{
		"find batch": {
			expResp: []byte(`{
			"type": "peggy/OutgoingTxBatch",
			"value": {
			"nonce": "1",
			"elements": [
				{
					"id": "2",
					"sender": "cosmos1qyqszqgpqyqszqgpqyqszqgpqyqszqgpjnp7du",
					"dest_address": "0x320915BD0F1bad11cBf06e85D5199DBcAC4E9934",
					"erc20_token": {
						"amount": "101",
						"symbol": "myETHToken",
						"token_contract_address": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
					},
					"erc20_fee": {
						"amount": "3",
						"symbol": "myETHToken",
						"token_contract_address": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
					}
				},
				{
					"id": "1",
					"sender": "cosmos1qyqszqgpqyqszqgpqyqszqgpqyqszqgpjnp7du",
					"dest_address": "0x320915BD0F1bad11cBf06e85D5199DBcAC4E9934",
					"erc20_token": {
						"amount": "100",
						"symbol": "myETHToken",
						"token_contract_address": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
					},
					"erc20_fee": {
						"amount": "2",
						"symbol": "myETHToken",
						"token_contract_address": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
					}
				}
			],
			"erc20_fee": {
				"amount": "5",
				"symbol": "myETHToken",
				"token_contract_address": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
			},
			"bridged_denominator": {
				"token_contract_address": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B",
				"symbol": "myETHToken",
				"cosmos_voucher_denom": "peggyf005bf9aac"
			},
			"valset": {
				"nonce": "105",
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
			got, err := lastPendingBatchRequest(ctx, valAddr.String(), k)
			require.NoError(t, err)
			assert.JSONEq(t, string(spec.expResp), string(got), string(got))
		})
	}
}

func createTestBatch(t *testing.T, k Keeper, ctx sdk.Context, keepers TestKeepers) {
	var (
		mySender            = bytes.Repeat([]byte{1}, sdk.AddrLen)
		myReceiver          = types.NewEthereumAddress("0x320915BD0F1bad11cBf06e85D5199DBcAC4E9934")
		myTokenContractAddr = types.NewEthereumAddress("0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B")
		myETHToken          = "myETHToken"
		voucherDenom        = types.NewVoucherDenom(myTokenContractAddr, myETHToken)
		now                 = time.Now().UTC()
	)
	// mint some voucher first
	allVouchers := sdk.Coins{sdk.NewInt64Coin(string(voucherDenom), 99999)}
	err := keepers.BankKeeper.MintCoins(ctx, types.ModuleName, allVouchers)
	require.NoError(t, err)

	// set senders balance
	keepers.AccountKeeper.NewAccountWithAddress(ctx, mySender)
	err = keepers.BankKeeper.SetBalances(ctx, mySender, allVouchers)
	require.NoError(t, err)

	// store counterpart
	k.StoreCounterpartDenominator(ctx, myTokenContractAddr, myETHToken)

	_ = types.NewBridgedDenominator(myTokenContractAddr, myETHToken)

	// add some TX to the pool
	for i, v := range []int64{2, 3, 2, 1} {
		amount := sdk.NewInt64Coin(string(voucherDenom), int64(i+100))
		fee := sdk.NewInt64Coin(string(voucherDenom), v)
		_, err := k.AddToOutgoingPool(ctx, mySender, myReceiver, amount, fee)
		require.NoError(t, err)
	}
	// when
	ctx = ctx.WithBlockTime(now)

	// tx batch size is 2, so that some of them stay behind
	_, err = k.BuildOutgoingTXBatch(ctx, voucherDenom, 2)
	require.NoError(t, err)
}

// TODO: Query more than one batch confirm
func TestQueryAllBatchConfirms(t *testing.T) {
	k, ctx, _ := CreateTestEnv(t)

	var (
		tokenContract    = types.NewEthereumAddress("0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B")
		validatorAddr, _ = sdk.AccAddressFromBech32("cosmos1mgamdcs9dah0vn0gqupl05up7pedg2mvupe6hh")
	)

	k.SetBatchConfirm(ctx, &types.MsgConfirmBatch{
		Nonce:         1,
		TokenContract: tokenContract.String(),
		EthSigner:     types.NewEthereumAddress("0xf35e2cc8e6523d683ed44870f5b7cc785051a77d").String(),
		Validator:     validatorAddr.String(),
		Signature:     "signature",
	})

	batchConfirms, err := queryAllBatchConfirms(ctx, "1", tokenContract.String(), k)
	require.NoError(t, err)

	expectedJSON := []byte(`[{"eth_signer":"0xF35e2cC8E6523d683eD44870f5B7cC785051a77D", "nonce":"1", "signature":"signature", "token_contract":"0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B", "validator":"cosmos1mgamdcs9dah0vn0gqupl05up7pedg2mvupe6hh"}]`)

	assert.JSONEq(t, string(expectedJSON), string(batchConfirms), "json is equal")
}

// TODO: test that it gets the correct batch, not just any batch.
// Check with multiple nonces and tokenContracts
func TestQueryBatch(t *testing.T) {
	k, ctx, keepers := CreateTestEnv(t)

	var (
		tokenContract = types.NewEthereumAddress("0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B")
	)

	createTestBatch(t, k, ctx, keepers)

	batch, err := queryBatch(ctx, "1", tokenContract.String(), k)
	require.NoError(t, err)

	expectedJSON := []byte(`{
		"type": "peggy/OutgoingTxBatch",
		"value": {
		  "bridged_denominator": {
			"cosmos_voucher_denom": "peggyf005bf9aac",
			"symbol": "myETHToken",
			"token_contract_address": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
		  },
		  "elements": [
			{
			  "erc20_fee": {
				"amount": "3",
				"symbol": "myETHToken",
				"token_contract_address": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
			  },
			  "dest_address": "0x320915BD0F1bad11cBf06e85D5199DBcAC4E9934",
			  "erc20_token": {
				"amount": "101",
				"symbol": "myETHToken",
				"token_contract_address": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
			  },
			  "sender": "cosmos1qyqszqgpqyqszqgpqyqszqgpqyqszqgpjnp7du",
			  "id": "2"
			},
			{
			  "erc20_fee": {
				"amount": "2",
				"symbol": "myETHToken",
				"token_contract_address": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
			  },
			  "dest_address": "0x320915BD0F1bad11cBf06e85D5199DBcAC4E9934",
			  "erc20_token": {
				"amount": "100",
				"symbol": "myETHToken",
				"token_contract_address": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
			  },
			  "sender": "cosmos1qyqszqgpqyqszqgpqyqszqgpqyqszqgpjnp7du",
			  "id": "1"
			}
		  ],
		  "nonce": "1",
		  "token_contract": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B",
		  "erc20_fee": {
			"amount": "5",
			"symbol": "myETHToken",
			"token_contract_address": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
		  },
		  "valset": { 
			  "nonce": "1234567" 
		  }
		}
	  }
	  `)

	// TODO: this test is failing on the empty representation of valset members
	assert.JSONEq(t, string(expectedJSON), string(batch), "json is equal")
}

func TestLastBatchesRequest(t *testing.T) {
	k, ctx, keepers := CreateTestEnv(t)

	createTestBatch(t, k, ctx, keepers)
	createTestBatch(t, k, ctx, keepers)

	lastBatches, err := lastBatchesRequest(ctx, k)
	require.NoError(t, err)

	expectedJSON := []byte(`[
		{
		  "bridged_denominator": {
			"cosmos_voucher_denom": "peggyf005bf9aac",
			"symbol": "myETHToken",
			"token_contract_address": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
		  },
		  "elements": [
			{
			  "erc20_fee": {
				"amount": "3",
				"symbol": "myETHToken",
				"token_contract_address": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
			  },
			  "dest_address": "0x320915BD0F1bad11cBf06e85D5199DBcAC4E9934",
			  "erc20_token": {
				"amount": "101",
				"symbol": "myETHToken",
				"token_contract_address": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
			  },
			  "sender": "cosmos1qyqszqgpqyqszqgpqyqszqgpqyqszqgpjnp7du",
			  "id": "6"
			},
			{
			  "erc20_fee": {
				"amount": "2",
				"symbol": "myETHToken",
				"token_contract_address": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
			  },
			  "dest_address": "0x320915BD0F1bad11cBf06e85D5199DBcAC4E9934",
			  "erc20_token": {
				"amount": "102",
				"symbol": "myETHToken",
				"token_contract_address": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
			  },
			  "sender": "cosmos1qyqszqgpqyqszqgpqyqszqgpqyqszqgpjnp7du",
			  "id": "3"
			}
		  ],
		  "nonce": "2",
		  "token_contract": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B",
		  "erc20_fee": {
			"amount": "5",
			"symbol": "myETHToken",
			"token_contract_address": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
		  },
		  "valset": { "nonce": "1234567" }
		},
		{
		  "bridged_denominator": {
			"cosmos_voucher_denom": "peggyf005bf9aac",
			"symbol": "myETHToken",
			"token_contract_address": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
		  },
		  "elements": [
			{
			  "erc20_fee": {
				"amount": "3",
				"symbol": "myETHToken",
				"token_contract_address": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
			  },
			  "dest_address": "0x320915BD0F1bad11cBf06e85D5199DBcAC4E9934",
			  "erc20_token": {
				"amount": "101",
				"symbol": "myETHToken",
				"token_contract_address": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
			  },
			  "sender": "cosmos1qyqszqgpqyqszqgpqyqszqgpqyqszqgpjnp7du",
			  "id": "2"
			},
			{
			  "erc20_fee": {
				"amount": "2",
				"symbol": "myETHToken",
				"token_contract_address": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
			  },
			  "dest_address": "0x320915BD0F1bad11cBf06e85D5199DBcAC4E9934",
			  "erc20_token": {
				"amount": "100",
				"symbol": "myETHToken",
				"token_contract_address": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
			  },
			  "sender": "cosmos1qyqszqgpqyqszqgpqyqszqgpqyqszqgpjnp7du",
			  "id": "1"
			}
		  ],
		  "nonce": "1",
		  "token_contract": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B",
		  "erc20_fee": {
			"amount": "5",
			"symbol": "myETHToken",
			"token_contract_address": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
		  },
		  "valset": { "nonce": "1234567" }
		}
	  ]
	  `)

	assert.JSONEq(t, string(expectedJSON), string(lastBatches), "json is equal")
}

func createEthAddress(i int) types.EthereumAddress {
	return types.EthereumAddress(gethCommon.BytesToAddress(bytes.Repeat([]byte{byte(i)}, 20)))
}

func createFakeEthSignature(n int) []byte {
	return bytes.Repeat([]byte{byte(n)}, 64)
}

func createValAddress() sdk.ValAddress {
	return sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address())
}
