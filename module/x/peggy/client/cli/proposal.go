package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/althea-net/peggy/module/x/peggy/types"
)

func GetCmdSubmitPeggyBootstrapProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "peggy-bootstrap [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a peggy bootstrap proposal",
		Example: fmt.Sprintf(
			"%s tx gov submit-proposal peggy-bootstrap <path/to/proposal.json> --from=<key_or_address>",
			version.AppName,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			proposal, err := ParsePeggyBootstrapProposalWithDeposit(clientCtx.JSONMarshaler, args[0])
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoins(proposal.Deposit)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			content := types.NewPeggyBootstrapProposal(
				proposal.Title,
				proposal.Description,
				proposal.PeggyId,
				proposal.ProxyContractHash,
				proposal.ProxyContractAddress,
				proposal.LogicContractHash,
				proposal.LogicContractAddress,
				proposal.StartThreshold,
				proposal.BridgeChainId,
				proposal.BootstrapValsetNonce,
			)

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}

func GetCmdSubmitPeggyUpgradeProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "peggy-upgrade [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a peggy upgrade proposal",
		Example: fmt.Sprintf(
			"%s tx gov submit-proposal peggy-upgrade <path/to/proposal.json> --from=<key_or_address>",
			version.AppName,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			proposal, err := ParsePeggyUpgradeProposalWithDeposit(clientCtx.JSONMarshaler, args[0])
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoins(proposal.Deposit)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			content := types.NewPeggyUpgradeProposal(
				proposal.Title,
				proposal.Description,
				proposal.Version,
				proposal.LogicContractHash,
				proposal.LogicContractAddress,
			)

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}
