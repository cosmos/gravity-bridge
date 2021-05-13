package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
	"github.com/spf13/cobra"
)

func GetQueryCmd() *cobra.Command {
	gravityQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the gravity module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	gravityQueryCmd.AddCommand([]*cobra.Command{
		CmdGetCurrentSignerSetTx(),
		CmdGetSignerSetTx(),
		CmdGetDelegateAddress(),
		CmdGetSignerSetTxSignature(),
		CmdGetPendingSignerSetTx(),
		CmdGetPendingBatchTx(),
		QueryAccepted(),
		QueryApproved(),
	}...)

	return gravityQueryCmd
}

func QueryAccepted() *cobra.Command {
	testingTxCmd := &cobra.Command{
		Use:                        "accepted",
		Short:                      "accepted ETH events",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	testingTxCmd.AddCommand([]*cobra.Command{
		// CmdGetLastAcceptedNonceRequest(storeKey, cdc),
		// CmdGetLastAcceptedNoncesRequest(storeKey, cdc),
		// CmdGetLastAcceptedMultiSigUpdateRequest(storeKey, cdc),
		// CmdGetAllBridgedDenominatorsRequest(storeKey, cdc),
	}...)

	return testingTxCmd
}

// TODO: wtf does this do
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

func CmdGetCurrentSignerSetTx() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "current-signer-set-tx",
		Short: "Query current signer set tx",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			req := &types.CurrentSignerSetTxRequest{}

			res, err := queryClient.CurrentSignerSetTx(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdGetDelegateAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegate-keys [validator]",
		Short: "Get delegate eth and cosmos key for a given validator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			validator, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			req := &types.DelegateKeysByValidatorAddress{
				ValidatorAddress: validator.String(),
			}

			res, err := queryClient.GetDelegateKeyByValidator(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdGetSignerSetTx() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "signer-set-tx [nonce]",
		Short: "Get requested signer set tx with a particular nonce",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			nonce, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			req := &types.SignerSetTxRequest{
				Nonce: nonce,
			}

			res, err := queryClient.SignerSetTx(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdGetSignerSetTxSignature() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "signer-set-tx-signatures [nonce] [bech32 validator address]",
		Short: "Get signer set tx signature with a particular nonce from a particular validator",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			nonce, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			req := &types.SignerSetTxSignatureRequest{
				Nonce:   nonce,
				Address: args[1],
			}

			res, err := queryClient.SignerSetTxSignature(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdGetPendingSignerSetTx() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pending-signer-set-tx [bech32 validator address]",
		Short: "Get the latest signer set tx which has not been signed by a particular validator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			req := &types.LastPendingSignerSetTxByAddrRequest{
				Address: args[0],
			}

			res, err := queryClient.LastPendingSignerSetTxByAddr(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdGetPendingBatchTx() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pending-batch-tx [bech32 validator address]",
		Short: "Get the latest batch tx which has not been signed by a particular validator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			req := &types.LastPendingBatchTxByAddrRequest{
				Address: args[0],
			}

			res, err := queryClient.LastPendingBatchTxByAddr(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
