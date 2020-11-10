package cli

import (
	"errors"
	"fmt"

	"github.com/althea-net/peggy/module/x/peggy/keeper"
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
		CmdGetOutgoingTXBatchByNonceRequest(storeKey, cdc),
		CmdGetAllAttestationsRequest(storeKey, cdc),
		CmdGetAttestationRequest(storeKey, cdc),
		QueryObserved(storeKey, cdc),
		QueryApproved(storeKey, cdc),
	)...)

	return peggyQueryCmd
}

func QueryObserved(storeKey string, cdc *codec.Codec) *cobra.Command {
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
		CmdGetLastObservedMultiSigUpdateRequest(storeKey, cdc),
		CmdGetAllBridgedDenominatorsRequest(storeKey, cdc),
	)...)

	return testingTxCmd
}
func QueryApproved(storeKey string, cdc *codec.Codec) *cobra.Command {
	testingTxCmd := &cobra.Command{
		Use:                        "approved",
		Short:                      "approved cosmos operation",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	testingTxCmd.AddCommand(flags.PostCommands(
		CmdGetLastApprovedNoncesRequest(storeKey, cdc),
		CmdGetLastApprovedMultiSigUpdateRequest(storeKey, cdc),
		CmdGetInflightBatchesRequest(storeKey, cdc),
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
				fmt.Println("Nothing found")
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
		Short: fmt.Sprintf("Get the last nonce that was observed for a claim type of %s", types.ToClaimTypeNames(types.AllOracleClaimTypes...)),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/lastNonce/%s", storeKey, args[0]), nil)
			if err != nil {
				return err
			}
			if len(res) == 0 {
				fmt.Println("Nothing found")
				return nil
			}

			var out types.UInt64Nonce
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
				fmt.Println("Nothing found")
				return nil
			}

			var out map[string]types.UInt64Nonce
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func CmdGetLastObservedMultiSigUpdateRequest(storeKey string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "last-multisig-update",
		Short: "Get last observed multisig update",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/lastObservedMultiSigUpdate", storeKey), nil)
			if err != nil {
				return err
			}
			if len(res) == 0 {
				fmt.Println("Nothing found")
				return nil
			}

			var out keeper.MultiSigUpdateResponse
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}
func CmdGetLastApprovedMultiSigUpdateRequest(storeKey string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "last-multisig-update",
		Short: "Get last approved multisig update",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/lastApprovedMultiSigUpdate", storeKey), nil)
			if err != nil {
				return err
			}
			if len(res) == 0 {
				fmt.Println("Nothing found")
				return nil
			}

			var out keeper.MultiSigUpdateResponse
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
				fmt.Println("Nothing found")
				return nil
			}

			var out types.OutgoingTxBatch
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}
func CmdGetOutgoingTXBatchByNonceRequest(storeKey string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "batch-request [nonce]",
		Short: "Get an outgoing TX batch by nonce",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/batch/%s", storeKey, args[0]), nil)
			if err != nil {
				return err
			}
			if len(res) == 0 {
				fmt.Println("Nothing found")
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

func CmdGetAllAttestationsRequest(storeKey string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "all-attestations [claim type]",
		Short: fmt.Sprintf("Get all attestations by claim type descending order. Claim types: %s", types.ToClaimTypeNames(types.AllOracleClaimTypes...)),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/allAttestations/%s", storeKey, args[0]), nil)
			if err != nil {
				return err
			}
			if len(res) == 0 {
				fmt.Println("Nothing found")
				return nil
			}
			var out []types.Attestation
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}
func CmdGetAttestationRequest(storeKey string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "attestation [claim type] [nonce]",
		Short: fmt.Sprintf("Get attestation by claim type and nonce. Claim types: %s", types.ToClaimTypeNames(types.AllOracleClaimTypes...)),
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			nonce, err := types.UInt64NonceFromString(args[1])
			if err != nil {
				return err
			}
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/attestation/%s/%s", storeKey, args[0], nonce.String()), nil)
			if err != nil {
				return err
			}
			if len(res) == 0 {
				fmt.Println("Nothing found")
				return nil
			}
			var out types.Attestation
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func CmdGetAllBridgedDenominatorsRequest(storeKey string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "all-bridged-denominators",
		Short: "Get all bridged ERC20 denominators on the cosmos side",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/allBridgedDenominators", storeKey), nil)
			if err != nil {
				return err
			}
			if len(res) == 0 {
				fmt.Println("Nothing found")
				return nil
			}
			var out []types.BridgedDenominator
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func CmdGetInflightBatchesRequest(storeKey string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "inflight-batches",
		Short: "Get all batches that have been approved but were not observed, yet",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/inflightBatches", storeKey), nil)
			if err != nil {
				return err
			}
			if len(res) == 0 {
				fmt.Println("Nothing found")
				return nil
			}
			var out []keeper.ApprovedOutgoingTxBatchResponse
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}
func CmdGetLastApprovedNoncesRequest(storeKey string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "nonces",
		Short: "Get last approved nonces for all claim types",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/lastApprovedNonces", storeKey), nil)
			if err != nil {
				return err
			}
			if len(res) == 0 {
				fmt.Println("Nothing found")
				return nil
			}

			var out map[string]types.UInt64Nonce
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}
