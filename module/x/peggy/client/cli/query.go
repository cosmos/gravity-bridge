package cli

import (
	"errors"
	"fmt"

	"github.com/althea-net/peggy/module/x/peggy/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	peggyQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the peggy module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	peggyQueryCmd.AddCommand(flags.GetCommands(
		CmdGetCurrentValset(storeKey, cdc),
		CmdGetValsetRequest(storeKey, cdc),
		CmdGetValsetConfirm(storeKey, cdc),
		CmdGetPendingValsetRequest(storeKey, cdc),
		CmdGetPendingOutgoingTXBatchRequest(storeKey, cdc),
		CmdGetAllOutgoingTXBatchRequest(storeKey, cdc),
		QueryOracle(storeKey, cdc),
	)...)

	return peggyQueryCmd
}

func QueryOracle(storeKey string, cdc *codec.Codec) *cobra.Command {
	testingTxCmd := &cobra.Command{
		Use:                        "observed",
		Short:                      "observed ETH events",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	testingTxCmd.AddCommand(flags.PostCommands(
		CmdGetLastObservedNonceRequest(storeKey, cdc),
		CmdGetLastObservedNoncesRequest(storeKey, cdc),
	)...)

	return testingTxCmd
}

func CmdGetCurrentValset(storeKey string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "current-valset",
		Short: "Query current valset",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/currentValset", storeKey), nil)
			if err != nil {
				return err
			}
			if len(res) == 0 {
				return errors.New("empty response")
			}

			var out types.Valset
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func CmdGetValsetRequest(storeKey string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "valset-request [nonce]",
		Short: "Get requested valset with a particular nonce",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			nonce := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/valsetRequest/%s", storeKey, nonce), nil)
			if err != nil {
				return err
			}
			if len(res) == 0 {
				return fmt.Errorf("no valset request found for nonce %s", nonce)
			}

			var out types.Valset
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func CmdGetValsetConfirm(storeKey string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "valset-confirm [nonce] [bech32 validator address]",
		Short: "Get valset confirmation with a particular nonce from a particular validator",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/valsetConfirm/%s/%s", storeKey, args[0], args[1]), nil)
			if err != nil {
				return err
			}
			if len(res) == 0 {
				return fmt.Errorf("no valset confirmation found for nonce %s and address %s", args[0], args[1])
			}

			var out types.MsgValsetConfirm
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func CmdGetPendingValsetRequest(storeKey string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "pending-valset-request [bech32 validator address]",
		Short: "Get the latest valset request which has not been signed by a particular validator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/lastPendingValsetRequest/%s", storeKey, args[0]), nil)
			if err != nil {
				return err
			}
			if len(res) == 0 {
				fmt.Println("No pending valset request")
				return nil
			}

			var out types.Valset
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func CmdGetLastObservedNonceRequest(storeKey string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "nonce [claim type]",
		Short: fmt.Sprintf("Get the last nonce that was observed for a claim type of %s", types.AllClaimTypes),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/lastObservedNonce/%s", storeKey, args[0]), nil)
			if err != nil {
				return err
			}
			if len(res) == 0 {
				fmt.Println("No observed nonce, yet")
				return nil
			}

			var out types.Nonce
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func CmdGetLastObservedNoncesRequest(storeKey string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "nonces",
		Short: "Get last observed nonces for all claim types",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/lastObservedNonces", storeKey), nil)
			if err != nil {
				return err
			}
			if len(res) == 0 {
				fmt.Println("No observed nonces, yet")
				return nil
			}

			var out map[string]types.Nonce
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func CmdGetPendingOutgoingTXBatchRequest(storeKey string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "pending-batch-request [bech32 validator address]",
		Short: "Get the latest outgoing TX batch request which has not been signed by a particular validator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/lastPendingBatchRequest/%s", storeKey, args[0]), nil)
			if err != nil {
				return err
			}
			if len(res) == 0 {
				fmt.Println("No pending outgoing batches")
				return nil
			}

			var out types.OutgoingTxBatch
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func CmdGetAllOutgoingTXBatchRequest(storeKey string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "all-batches",
		Short: "Get all batches descending order",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/allBatches", storeKey), nil)
			if err != nil {
				return err
			}
			if len(res) == 0 {
				fmt.Println("No outgoing batches")
				return nil
			}

			var out []types.OutgoingTxBatch
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}
