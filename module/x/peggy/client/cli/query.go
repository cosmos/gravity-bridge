package cli

import (
	"strconv"

	"github.com/althea-net/peggy/module/x/peggy/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
)

func GetQueryCmd() *cobra.Command {
	peggyQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the peggy module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	peggyQueryCmd.AddCommand([]*cobra.Command{
		CmdGetCurrentValset(),
		CmdGetValsetRequest(),
		CmdGetValsetConfirm(),
		CmdGetPendingValsetRequest(),
		CmdGetPendingOutgoingTXBatchRequest(),
		// CmdGetAllOutgoingTXBatchRequest(),
		// CmdGetOutgoingTXBatchByNonceRequest(),
		// CmdGetAllAttestationsRequest(),
		// CmdGetAttestationRequest(),
		QueryObserved(),
		QueryApproved(),
	}...)

	return peggyQueryCmd
}

func QueryObserved() *cobra.Command {
	testingTxCmd := &cobra.Command{
		Use:                        "observed",
		Short:                      "observed ETH events",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	testingTxCmd.AddCommand([]*cobra.Command{
		// CmdGetLastObservedNonceRequest(storeKey, cdc),
		// CmdGetLastObservedNoncesRequest(storeKey, cdc),
		// CmdGetLastObservedMultiSigUpdateRequest(storeKey, cdc),
		// CmdGetAllBridgedDenominatorsRequest(storeKey, cdc),
	}...)

	return testingTxCmd
}
func QueryApproved() *cobra.Command {
	testingTxCmd := &cobra.Command{
		Use:                        "approved",
		Short:                      "approved cosmos operation",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	testingTxCmd.AddCommand([]*cobra.Command{
		// CmdGetLastApprovedNoncesRequest(storeKey, cdc),
		// CmdGetLastApprovedMultiSigUpdateRequest(storeKey, cdc),
		// CmdGetInflightBatchesRequest(storeKey, cdc),
	}...)

	return testingTxCmd
}

func CmdGetCurrentValset() *cobra.Command {
	return &cobra.Command{
		Use:   "current-valset",
		Short: "Query current valset",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryCurrentValsetRequest{}

			res, err := queryClient.CurrentValset(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
}

func CmdGetValsetRequest() *cobra.Command {
	return &cobra.Command{
		Use:   "valset-request [nonce]",
		Short: "Get requested valset with a particular nonce",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			nonce, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			req := &types.QueryValsetRequestRequest{
				Nonce: nonce,
			}

			res, err := queryClient.ValsetRequest(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
}

func CmdGetValsetConfirm() *cobra.Command {
	return &cobra.Command{
		Use:   "valset-confirm [nonce] [bech32 validator address]",
		Short: "Get valset confirmation with a particular nonce from a particular validator",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			nonce, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			req := &types.QueryValsetConfirmRequest{
				Nonce:   nonce,
				Address: args[1],
			}

			res, err := queryClient.ValsetConfirm(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
}

func CmdGetPendingValsetRequest() *cobra.Command {
	return &cobra.Command{
		Use:   "pending-valset-request [bech32 validator address]",
		Short: "Get the latest valset request which has not been signed by a particular validator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryLastPendingValsetRequestByAddrRequest{
				Address: args[0],
			}

			res, err := queryClient.LastPendingValsetRequestByAddr(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
}

func CmdGetPendingOutgoingTXBatchRequest() *cobra.Command {
	return &cobra.Command{
		Use:   "pending-batch-request [bech32 validator address]",
		Short: "Get the latest outgoing TX batch request which has not been signed by a particular validator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryLastPendingBatchRequestByAddrRequest{
				Address: args[0],
			}

			res, err := queryClient.LastPendingBatchRequestByAddr(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
}
