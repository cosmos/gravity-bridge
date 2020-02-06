package nftbridge

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/cosmos/peggy/x/oracle"
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethbridge "github.com/cosmos/peggy/x/ethbridge/types"
	"github.com/cosmos/peggy/x/nftbridge/types"
	"github.com/stretchr/testify/require"
)

const STATUS = "status"

func TestBasicMsgs(t *testing.T) {
	//Setup
	ctx, _, _, _, _, validatorAddresses, handler := CreateTestHandler(t, 0.7, []int64{3, 7})

	valAddress := validatorAddresses[0]

	//Unrecognized type
	res, err := handler(ctx, sdk.NewTestMsg())
	require.Error(t, err)
	require.Nil(t, res)
	require.True(t, strings.Contains(err.Error(), "unrecognized nftbridge message type: "))

	//Normal Creation
	normalCreateNFTMsg := types.CreateTestNFTMsg(t, valAddress, ethbridge.LockText)
	res, err = handler(ctx, normalCreateNFTMsg)
	require.NoError(t, err)
	require.NotNil(t, res)

	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			switch key := string(attribute.Key); key {
			case "module":
				require.Equal(t, value, types.ModuleName)
			case "sender":
				require.Equal(t, value, valAddress.String())
			case "ethereum_sender":
				require.Equal(t, value, types.TestEthereumAddress)
			case "cosmos_receiver":
				require.Equal(t, value, types.TestAddress)
			case "denom":
				require.Equal(t, value, types.TestDenom)
			case "id":
				require.Equal(t, value, types.TestID)
			case STATUS:
				require.Equal(t, value, oracle.StatusTextToString[oracle.PendingStatusText])
			case "claim_type":
				require.Equal(t, value, ethbridge.ClaimTypeToString[ethbridge.LockText])
			default:
				require.Fail(t, fmt.Sprintf("unrecognized event %s", key))
			}
		}
	}

	//Bad Creation
	badCreateMsg := types.CreateTestNFTMsg(t, valAddress, ethbridge.LockText)
	badCreateMsg.Nonce = -1
	err = badCreateMsg.ValidateBasic()
	require.Error(t, err)
}

func TestDuplicateMsgs(t *testing.T) {
	ctx, _, _, _, _, validatorAddresses, handler := CreateTestHandler(t, 0.7, []int64{3, 7})

	valAddress := validatorAddresses[0]

	normalCreateNFTMsg := types.CreateTestNFTMsg(t, valAddress, ethbridge.LockText)
	res, err := handler(ctx, normalCreateNFTMsg)
	require.NoError(t, err)
	require.NotNil(t, res)
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			if string(attribute.Key) == STATUS {
				require.Equal(t, value, oracle.StatusTextToString[oracle.PendingStatusText])
			}
		}
	}

	//Duplicate message from same validator
	res, err = handler(ctx, normalCreateNFTMsg)
	require.Error(t, err)
	require.Nil(t, res)
	require.True(t, strings.Contains(err.Error(), "already processed message from validator for this id"))
}

func TestMintSuccess(t *testing.T) {
	//Setup
	ctx, _, _, nftKeeper, _, validatorAddresses, handler := CreateTestHandler(t, 0.7, []int64{2, 7, 1})

	valAddressVal1Pow2 := validatorAddresses[0]
	valAddressVal2Pow7 := validatorAddresses[1]
	valAddressVal3Pow1 := validatorAddresses[2]

	//Initial message
	normalCreateNFTMsg := types.CreateTestNFTMsg(t, valAddressVal1Pow2, ethbridge.LockText)
	res, err := handler(ctx, normalCreateNFTMsg)
	require.NoError(t, err)
	require.NotNil(t, res)

	//Message from second validator succeeds and mints new tokens
	normalCreateNFTMsg = types.CreateTestNFTMsg(t, valAddressVal2Pow7, ethbridge.LockText)
	res, err = handler(ctx, normalCreateNFTMsg)
	require.NoError(t, err)
	require.NotNil(t, res)
	receiverAddress, err := sdk.AccAddressFromBech32(types.TestAddress)
	require.NoError(t, err)

	receiverNFT, err := nftKeeper.GetNFT(ctx, types.TestDenom, types.TestID)
	require.NoError(t, err)
	require.True(t, receiverAddress.Equals(receiverNFT.GetOwner()))

	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			if string(attribute.Key) == STATUS {
				require.Equal(t, value, oracle.StatusTextToString[oracle.SuccessStatusText])
			}
		}
	}

	//Additional message from third validator fails and does not mint
	normalCreateNFTMsg = types.CreateTestNFTMsg(t, valAddressVal3Pow1, ethbridge.LockText)
	res, err = handler(ctx, normalCreateNFTMsg)
	require.Error(t, err)
	require.Nil(t, res)
	require.True(t, strings.Contains(err.Error(), "prophecy already finalized"))

	// hasn't changed
	receiverNFT, err = nftKeeper.GetNFT(ctx, types.TestDenom, types.TestID)
	require.NoError(t, err)
	require.True(t, receiverAddress.Equals(receiverNFT.GetOwner()))

}

func TestNoMintFail(t *testing.T) {
	//Setup
	ctx, _, _, nftKeeper, _, validatorAddresses, handler := CreateTestHandler(t, 0.71, []int64{3, 4, 3})

	valAddressVal1Pow3 := validatorAddresses[0]
	valAddressVal2Pow4 := validatorAddresses[1]
	valAddressVal3Pow3 := validatorAddresses[2]

	testTokenContractAddress := ethbridge.NewEthereumAddress(types.TestTokenContractAddress)
	testEthereumAddress := ethbridge.NewEthereumAddress(types.TestEthereumAddress)

	ethClaim1 := types.CreateTestNFTClaim(t, testEthereumAddress, testTokenContractAddress, valAddressVal1Pow3, testEthereumAddress, types.TestDenom, types.TestID, ethbridge.LockText)
	ethMsg1 := NewMsgCreateNFTBridgeClaim(ethClaim1)
	ethClaim2 := types.CreateTestNFTClaim(t, testEthereumAddress, testTokenContractAddress, valAddressVal2Pow4, testEthereumAddress, types.TestDenom, types.TestID, ethbridge.LockText)
	ethMsg2 := NewMsgCreateNFTBridgeClaim(ethClaim2)
	ethClaim3 := types.CreateTestNFTClaim(t, testEthereumAddress, testTokenContractAddress, valAddressVal3Pow3, testEthereumAddress, types.AltTestDenom, types.AltTestID, ethbridge.LockText)
	ethMsg3 := NewMsgCreateNFTBridgeClaim(ethClaim3)

	//Initial message
	res, err := handler(ctx, ethMsg1)
	require.NoError(t, err)
	require.NotNil(t, res)
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			if string(attribute.Key) == STATUS {
				require.Equal(t, value, oracle.StatusTextToString[oracle.PendingStatusText])
			}
		}
	}

	//Different message from second validator succeeds
	res, err = handler(ctx, ethMsg2)
	require.NoError(t, err)
	require.NotNil(t, res)
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			if string(attribute.Key) == STATUS {
				require.Equal(t, value, oracle.StatusTextToString[oracle.PendingStatusText])
			}
		}
	}

	//Different message from third validator succeeds but results in failed prophecy with no minting
	res, err = handler(ctx, ethMsg3)
	require.NoError(t, err)
	require.NotNil(t, res)
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			if string(attribute.Key) == STATUS {
				require.Equal(t, value, oracle.StatusTextToString[oracle.FailedStatusText])
			}
		}
	}
	//TODO: what does this do?
	receiverAddress, err := sdk.AccAddressFromBech32(types.TestAddress)
	require.NoError(t, err)
	receiver1NFT, found := nftKeeper.GetOwnerByDenom(ctx, receiverAddress, types.TestDenom)
	fmt.Println("receiver1NFT", receiver1NFT)
	fmt.Println("found", found)
	// receiver1Coins := bankKeeper.GetCoins(ctx, receiverAddress)
	// require.True(t, receiver1Coins.IsZero())
}

func TestBurnEthFail(t *testing.T) {

}

func TestBurnEthSuccess(t *testing.T) {
	ctx, _, _, nftKeeper, _, validatorAddresses, handler := CreateTestHandler(t, 0.5, []int64{5})
	valAddressVal1Pow5 := validatorAddresses[0]

	moduleAccountAddress := sdk.AccAddress(crypto.AddressHash([]byte(ModuleName)))
	// moduleAccountAddress := moduleAccount.GetAddress()

	// Initial message to mint some eth
	// coinsToMint := "7ethereum"
	denom := types.TestDenom
	id := types.TestID

	testTokenContractAddress := ethbridge.NewEthereumAddress(types.TestTokenContractAddress)
	testEthereumAddress := ethbridge.NewEthereumAddress(types.TestEthereumAddress)

	ethClaim1 := types.CreateTestNFTClaim(t, testEthereumAddress, testTokenContractAddress, valAddressVal1Pow5, testEthereumAddress, denom, id, ethbridge.LockText)
	ethMsg1 := NewMsgCreateNFTBridgeClaim(ethClaim1)

	// Initial message succeeds and mints eth
	res, err := handler(ctx, ethMsg1)
	require.NoError(t, err)
	require.NotNil(t, res)
	receiverAddress, err := sdk.AccAddressFromBech32(types.TestAddress)
	require.NoError(t, err)

	receiverCollection, found := nftKeeper.GetOwnerByDenom(ctx, receiverAddress, denom)
	require.True(t, found)
	found = receiverCollection.Exists(id)
	require.True(t, found)

	// coinsToBurn := "3ethereum"
	ethereumReceiver := ethbridge.NewEthereumAddress(types.AltTestEthereumAddress)

	// Second message succeeds, burns eth and fires correct event
	burnMsg := types.CreateTestBurnMsg(t, types.TestAddress, ethereumReceiver, denom, id)
	res, err = handler(ctx, burnMsg)
	require.NoError(t, err)
	require.NotNil(t, res)
	senderAddress := receiverAddress

	receiverCollection, found = nftKeeper.GetOwnerByDenom(ctx, senderAddress, denom)
	require.True(t, found)
	found = receiverCollection.Exists(id)
	require.False(t, found)

	eventEthereumChainID := ""
	eventTokenContract := ""
	eventCosmosSender := ""
	eventEthereumReceiver := ""
	eventDenom := ""
	eventID := ""
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			switch key := string(attribute.Key); key {
			case "sender":
				require.Equal(t, value, senderAddress.String())
			case "recipient":
				require.Equal(t, value, moduleAccountAddress.String())
			case "module":
				require.Equal(t, value, ModuleName)
			case "ethereum_chain_id":
				eventEthereumChainID = value
			case "token_contract_address":
				eventTokenContract = value
			case "cosmos_sender":
				eventCosmosSender = value
			case "ethereum_receiver":
				eventEthereumReceiver = value
			case "denom":
				eventDenom = value
			case "id":
				eventID = value
			default:
				require.Fail(t, fmt.Sprintf("unrecognized event %s", key))
			}
		}
	}
	require.Equal(t, eventEthereumChainID, strconv.Itoa(types.TestEthereumChainID))
	require.Equal(t, eventTokenContract, types.TestTokenContractAddress)
	require.Equal(t, eventCosmosSender, senderAddress.String())
	require.Equal(t, eventEthereumReceiver, ethereumReceiver.String())
	require.Equal(t, eventDenom, denom)
	require.Equal(t, eventID, id)

	// Third message fails, no longer owns NFT
	res, err = handler(ctx, burnMsg)
	require.Error(t, err)
	require.Nil(t, res)
}
