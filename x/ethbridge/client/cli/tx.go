package cli

import (
	"bufio"
	"regexp"
	"strconv"

	"github.com/trinhtan/peggy/x/ethbridge/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// GetCmdCreateEthBridgeClaim is the CLI command for creating a claim on an ethereum prophecy
//nolint:lll
func GetCmdCreateEthBridgeClaim(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "create-claim [bridge-registry-contract] [nonce] [symbol] [ethereum-sender-address] [cosmos-receiver-address] [validator-address] [amount] [claim-type] --ethereum-chain-id [ethereum-chain-id] --token-contract-address [token-contract-address]",
		Short: "create a claim on an ethereum prophecy",
		Args:  cobra.ExactArgs(8),
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
			if !common.IsHexAddress(tokenContractString) {
				return errors.Errorf("invalid [token-contract-address]: %s", tokenContractString)
			}
			tokenContract := types.NewEthereumAddress(tokenContractString)

			if !common.IsHexAddress(args[0]) {
				return errors.Errorf("invalid [bridge-registry-contract]: %s", args[0])
			}
			bridgeContract := types.NewEthereumAddress(args[0])

			nonce, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}

			symbol := args[2]
			ethereumSender := types.NewEthereumAddress(args[3])
			if !common.IsHexAddress(args[3]) {
				return errors.Errorf("invalid [ethereum-sender-address]: %s", args[0])
			}
			cosmosReceiver, err := sdk.AccAddressFromBech32(args[4])
			if err != nil {
				return err
			}

			validator, err := sdk.ValAddressFromBech32(args[5])
			if err != nil {
				return err
			}

			var digitCheck = regexp.MustCompile(`^[0-9]+$`)
			if !digitCheck.MatchString(args[6]) {
				return types.ErrInvalidAmount
			}
			amount, err := strconv.ParseInt(args[6], 10, 64)
			if err != nil {
				return err
			}
			if amount <= 0 {
				return types.ErrInvalidAmount
			}

			claimType, err := types.StringToClaimType(args[7])
			if err != nil {
				return err
			}

			ethBridgeClaim := types.NewEthBridgeClaim(ethereumChainID, bridgeContract, nonce, symbol, tokenContract,
				ethereumSender, cosmosReceiver, validator, amount, claimType)

			msg := types.NewMsgCreateEthBridgeClaim(ethBridgeClaim)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdBurn is the CLI command for burning some of your eth and triggering an event
//nolint:lll
func GetCmdBurn(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "burn [cosmos-sender-address] [ethereum-receiver-address] [amount] [symbol] --ethereum-chain-id [ethereum-chain-id]",
		Short: "burn cETH or cERC20 on the Cosmos chain",
		Long: `This should be used to burn cETH or cERC20. It will burn your coins on the Cosmos Chain, removing them from your account and deducting them from the supply.
		It will also trigger an event on the Cosmos Chain for relayers to watch so that they can trigger the withdrawal of the original ETH/ERC20 to you from the Ethereum contract!`,
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authtypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			ethereumChainIDString := viper.GetString(types.FlagEthereumChainID)
			ethereumChainID, err := strconv.Atoi(ethereumChainIDString)
			if err != nil {
				return err
			}

			cosmosSender, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			if !common.IsHexAddress(args[1]) {
				return errors.Errorf("invalid [ethereum-receiver-address]: %s", args[1])
			}
			ethereumReceiver := types.NewEthereumAddress(args[1])

			var digitCheck = regexp.MustCompile(`^[0-9]+$`)
			if !digitCheck.MatchString(args[2]) {
				return types.ErrInvalidAmount
			}
			amount, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				return err
			}
			if amount <= 0 {
				return types.ErrInvalidAmount
			}

			symbol := args[3]

			msg := types.NewMsgBurn(ethereumChainID, cosmosSender, ethereumReceiver, amount, symbol)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdLock is the CLI command for locking some of your coins and triggering an event
func GetCmdLock(cdc *codec.Codec) *cobra.Command {
	//nolint:lll
	return &cobra.Command{
		Use:   "lock [cosmos-sender-address] [ethereum-receiver-address] [amount] [symbol] --ethereum-chain-id [ethereum-chain-id]",
		Short: "This should be used to lock Cosmos-originating coins (eg: ATOM). It will lock up your coins in the supply module, removing them from your account. It will also trigger an event on the Cosmos Chain for relayers to watch so that they can trigger the minting of the pegged token on Etherum to you!",
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

			cosmosSender, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			if !common.IsHexAddress(args[1]) {
				return errors.Errorf("invalid [ethereum-receiver-address]: %s", args[1])
			}
			ethereumReceiver := types.NewEthereumAddress(args[1])

			var digitCheck = regexp.MustCompile(`^[0-9]+$`)
			if !digitCheck.MatchString(args[2]) {
				return types.ErrInvalidAmount
			}
			amount, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				return err
			}
			if amount <= 0 {
				return types.ErrInvalidAmount
			}

			symbol := args[3]

			msg := types.NewMsgLock(ethereumChainID, cosmosSender, ethereumReceiver, amount, symbol)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
