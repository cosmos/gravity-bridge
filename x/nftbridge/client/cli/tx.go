package cli

import (
	"bufio"
	"strconv"

	ethbridge "github.com/cosmos/peggy/x/ethbridge/types"
	"github.com/cosmos/peggy/x/nftbridge/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// GetCmdCreateNFTBridgeClaim is the CLI command for creating a claim on an nft prophecy
func GetCmdCreateNFTBridgeClaim(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "create-nft-claim [ethereum-chain-id] [bridge-contract] [nonce] [symbol] [token-contract] [ethereum-sender-address] [cosmos-receiver-address] [validator-address] [denom] [id] [claim-type]",
		Short: "create a claim on an ethereum prophecy",
		Args:  cobra.ExactArgs(11),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authtypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			ethereumChainID, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			bridgeContract := ethbridge.NewEthereumAddress(args[1])

			nonce, err := strconv.Atoi(args[2])
			if err != nil {
				return err
			}

			symbol := args[3]
			tokenContract := ethbridge.NewEthereumAddress(args[4])
			ethereumSender := ethbridge.NewEthereumAddress(args[5])
			cosmosReceiver, err := sdk.AccAddressFromBech32(args[6])
			if err != nil {
				return err
			}

			validator, err := sdk.ValAddressFromBech32(args[7])
			if err != nil {
				return err
			}

			denom := args[8]
			id := args[9]

			claimType, err := ethbridge.StringToClaimType(args[10])
			if err != nil {
				return err
			}

			nftBridgeClaim := types.NewNFTBridgeClaim(ethereumChainID, bridgeContract, nonce, symbol, tokenContract, ethereumSender, cosmosReceiver, validator, denom, id, claimType)

			msg := types.NewMsgCreateNFTBridgeClaim(nftBridgeClaim)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdBurnNFT is the CLI command for burning some of your eth and triggering an event
func GetCmdBurnNFT(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "burnNFT [cosmos-sender-address] [ethereum-receiver-address] [denom] [id] --ethereum-chain-id [ethereum-chain-id] --token-contract-address [token-contract-address]",
		Short: "burn cNFT on the Cosmos chain",
		Long:  "This should be used to burn cNFT. It will burn your NFT on the Cosmos Chain, removing them from your account and deducting them from the supply. It will also trigger an event on the Cosmos Chain for relayers to watch so that they can trigger the withdrawal of the original NFT to you from the Ethereum contract!",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authtypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			ethereumChainIDString := viper.GetString(types.FlagEthereumChainID)
			ethereumChainID, err := strconv.Atoi(ethereumChainIDString)
			if err != nil {
				return err
			}

			tokenContractString := viper.GetString(types.FlagTokenContractAddr)
			tokenContract := ethbridge.NewEthereumAddress(tokenContractString)

			cosmosSender, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			ethereumReceiver := ethbridge.NewEthereumAddress(args[1])
			denom := args[2]
			id := args[3]

			msg := types.NewMsgBurnNFT(ethereumChainID, tokenContract, cosmosSender, ethereumReceiver, denom, id)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdLockNFT is the CLI command for locking some of your coins and triggering an event
func GetCmdLockNFT(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "lockNFT [cosmos-sender-address] [ethereum-receiver-address] [denom] [id] --ethereum-chain-id [ethereum-chain-id] --token-contract-address [token-contract-address]",
		Short: "This should be used to lock Cosmos-originating NFT. It will lock up your NFT in the bridge module, removing them from your account. It will also trigger an event on the Cosmos Chain for relayers to watch so that they can trigger the minting of the pegged NFT on Etherum to you!",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authtypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			ethereumChainIDString := viper.GetString(types.FlagEthereumChainID)
			ethereumChainID, err := strconv.Atoi(ethereumChainIDString)
			if err != nil {
				return err
			}

			tokenContractString := viper.GetString(types.FlagTokenContractAddr)
			tokenContract := ethbridge.NewEthereumAddress(tokenContractString)

			cosmosSender, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			ethereumReceiver := ethbridge.NewEthereumAddress(args[1])
			denom := args[2]
			id := args[3]

			msg := types.NewMsgLockNFT(ethereumChainID, tokenContract, cosmosSender, ethereumReceiver, denom, id)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
