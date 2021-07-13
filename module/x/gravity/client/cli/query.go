package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/gravity-bridge/module/x/gravity/types"
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
	gravityQueryCmd.AddCommand(
		CmdParams(),
		CmdSignerSetTx(),
		CmdBatchTx(),
		CmdContractCallTx(),
		CmdSignerSetTxs(),
		CmdBatchTxs(),
		CmdContractCallTxs(),
		CmdSignerSetTxConfirmations(),
		CmdBatchTxConfirmations(),
		CmdContractCallTxConfirmations(),
		CmdUnsignedSignerSetTxs(),
		CmdUnsignedBatchTxs(),
		CmdUnsignedContractCallTxs(),
		CmdLastSubmittedEthereumEvent(),
		CmdBatchTxFees(),
		CmdERC20ToDenom(),
		CmdDenomToERC20(),
		CmdUnbatchedSendToEthereums(),
		CmdDelegateKeysByValidator(),
		CmdDelegateKeysByEthereumSigner(),
		CmdDelegateKeysByOrchestrator(),
		CmdDelegateKeys(),
	)

	return gravityQueryCmd
}

func CmdParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query votes on a proposal",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			req := types.ParamsRequest{}

			res, err := queryClient.Params(cmd.Context(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdSignerSetTx() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "signer-set-tx [nonce]",
		Args:  cobra.MaximumNArgs(1),
		Short: "", // TODO(levi) provide short description
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			var nonce uint64

			if len(args) > 0 {
				if nonce, err = parseNonce(args[0]); err != nil {
					return err
				}
			}

			req := types.SignerSetTxRequest{
				SignerSetNonce: nonce,
			}

			res, err := queryClient.SignerSetTx(cmd.Context(), &req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdBatchTx() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batch-tx [contract-address] [nonce]",
		Args:  cobra.RangeArgs(1, 2),
		Short: "", // TODO(levi) provide short description
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			var ( // args
				contractAddress string
				nonce           uint64
			)

			contractAddress, err = parseContractAddress(args[0])
			if err != nil {
				return nil
			}

			if len(args) == 2 {
				if nonce, err = parseNonce(args[1]); err != nil {
					return err
				}
			}

			req := types.BatchTxRequest{
				TokenContract: contractAddress,
				BatchNonce:    nonce,
			}

			res, err := queryClient.BatchTx(cmd.Context(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdContractCallTx() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contract-call-tx [invalidation-scope] [invalidation-nonce]",
		Args:  cobra.ExactArgs(2),
		Short: "", // TODO(levi) provide short description
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			invalidationScope := []byte(args[0])
			invalidationNonce, err := strconv.ParseUint(args[1], 10, 64)

			req := types.ContractCallTxRequest{
				InvalidationScope: invalidationScope,
				InvalidationNonce: invalidationNonce,
			}

			res, err := queryClient.ContractCallTx(cmd.Context(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdSignerSetTxs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "signer-set-txs (count)",
		Args:  cobra.NoArgs,
		Short: "", // TODO(levi) provide short description
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			req := types.SignerSetTxsRequest{}

			res, err := queryClient.SignerSetTxs(cmd.Context(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdBatchTxs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batch-txs",
		Args:  cobra.NoArgs,
		Short: "", // TODO(levi) provide short description
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			req := types.BatchTxsRequest{}

			res, err := queryClient.BatchTxs(cmd.Context(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdContractCallTxs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contract-call-txs",
		Args:  cobra.NoArgs,
		Short: "", // TODO(levi) provide short description
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			req := types.ContractCallTxsRequest{}

			res, err := queryClient.ContractCallTxs(cmd.Context(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdSignerSetTxConfirmations() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "signer-set-tx-ethereum-signatures [nonce]",
		Args:  cobra.ExactArgs(1),
		Short: "", // TODO(levi) provide short description
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			nonce, err := parseNonce(args[0])
			if err != nil {
				return err
			}

			req := types.SignerSetTxConfirmationsRequest{
				SignerSetNonce: nonce,
			}

			res, err := queryClient.SignerSetTxConfirmations(cmd.Context(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdBatchTxConfirmations() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batch-tx-ethereum-signatures [nonce] [contract-address]",
		Args:  cobra.MinimumNArgs(2),
		Short: "", // TODO(levi) provide short description
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			var ( // args
				nonce           uint64
				contractAddress string
			)

			if nonce, err = parseNonce(args[0]); err != nil {
				return err
			}

			contractAddress, err = parseContractAddress(args[1])
			if err != nil {
				return nil
			}

			req := types.BatchTxConfirmationsRequest{
				BatchNonce:    nonce,
				TokenContract: contractAddress,
			}

			res, err := queryClient.BatchTxConfirmations(cmd.Context(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdContractCallTxConfirmations() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contract-call-tx-ethereum-signatures [invalidation-scope] [invalidation-nonce]",
		Args:  cobra.MinimumNArgs(2),
		Short: "", // TODO(levi) provide short description
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			invalidationScope := []byte(args[0])
			invalidationNonce, err := parseNonce(args[1])
			if err != nil {
				return err
			}

			req := types.ContractCallTxConfirmationsRequest{
				InvalidationNonce: invalidationNonce,
				InvalidationScope: invalidationScope,
			}

			res, err := queryClient.ContractCallTxConfirmations(cmd.Context(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdUnsignedSignerSetTxs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pending-signer-set-tx-ethereum-signatures [validator-or-orchestrator-address]",
		Args:  cobra.ExactArgs(1),
		Short: "", // TODO(levi) provide short description
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			address := args[0]
			req := types.UnsignedSignerSetTxsRequest{
				Address: address,
			}

			res, err := queryClient.UnsignedSignerSetTxs(cmd.Context(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdUnsignedBatchTxs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pending-batch-tx-ethereum-signatures [address]",
		Args:  cobra.ExactArgs(1),
		Short: "", // TODO(levi) provide short description
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			address := args[0]
			req := types.UnsignedBatchTxsRequest{
				Address: address,
			}

			res, err := queryClient.UnsignedBatchTxs(cmd.Context(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdUnsignedContractCallTxs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pending-contract-call-tx-ethereum-signatures [address]",
		Args:  cobra.ExactArgs(1),
		Short: "", // TODO(levi) provide short description
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			address := args[0]
			req := types.UnsignedContractCallTxsRequest{
				Address: address,
			}

			res, err := queryClient.UnsignedContractCallTxs(cmd.Context(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdLastSubmittedEthereumEvent() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "last-submitted-ethereum-event [address]",
		Args:  cobra.ExactArgs(1),
		Short: "", // TODO(levi) provide short description
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			address := args[0]
			req := types.LastSubmittedEthereumEventRequest{
				Address: address, // TODO(levi) what kind of address is this??
			}

			res, err := queryClient.LastSubmittedEthereumEvent(cmd.Context(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdBatchTxFees() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batch-tx-fees",
		Args:  cobra.NoArgs,
		Short: "", // TODO(levi) provide short description
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			req := types.BatchTxFeesRequest{}

			res, err := queryClient.BatchTxFees(cmd.Context(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdERC20ToDenom() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "erc20-to-denom [erc20]",
		Args:  cobra.ExactArgs(1),
		Short: "", // TODO(levi) provide short description
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			var ( // args
				erc20 string // TODO(levi) init and validate from args[0]
			)

			req := types.ERC20ToDenomRequest{
				Erc20: erc20, // TODO(levi) is this an ethereum address??
			}

			res, err := queryClient.ERC20ToDenom(cmd.Context(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdDenomToERC20() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "denom-to-erc20 [denom]",
		Args:  cobra.ExactArgs(1),
		Short: "", // TODO(levi) provide short description
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			var ( // args
				denom string // TODO(levi) init and validate from args[0]
			)

			req := types.DenomToERC20Request{
				Denom: denom, // TODO(levi) do we validate denoms?? if so, how?
			}

			res, err := queryClient.DenomToERC20(cmd.Context(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdUnbatchedSendToEthereums() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unbatched-send-to-ethereums [sender-address]",
		Args:  cobra.MaximumNArgs(1),
		Short: "", // TODO(levi) provide short description
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			var ( // args
				senderAddress string // TODO(levi) init and validate from args[0]
			)

			req := types.UnbatchedSendToEthereumsRequest{
				SenderAddress: senderAddress, // TODO(levi) is this an ethereum address??
			}

			res, err := queryClient.UnbatchedSendToEthereums(cmd.Context(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdDelegateKeysByValidator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegate-keys-by-validator [validator-address]",
		Args:  cobra.ExactArgs(1),
		Short: "", // TODO(levi) provide short description
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			validatorAddress := args[0]

			req := types.DelegateKeysByValidatorRequest{
				ValidatorAddress: validatorAddress,
			}

			res, err := queryClient.DelegateKeysByValidator(cmd.Context(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdDelegateKeysByEthereumSigner() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegate-keys-by-ethereum-signer [ethereum-signer]",
		Args:  cobra.ExactArgs(1),
		Short: "", // TODO(levi) provide short description
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			ethereumSigner := args[0] // TODO(levi) init and validate from args[0]

			req := types.DelegateKeysByEthereumSignerRequest{
				EthereumSigner: ethereumSigner,
			}

			res, err := queryClient.DelegateKeysByEthereumSigner(cmd.Context(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdDelegateKeysByOrchestrator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegate-keys-by-orchestrator [orchestrator-address]",
		Args:  cobra.ExactArgs(1),
		Short: "", // TODO(levi) provide short description
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			orcAddr := args[0]

			req := types.DelegateKeysByOrchestratorRequest{
				OrchestratorAddress: orcAddr,
			}

			res, err := queryClient.DelegateKeysByOrchestrator(cmd.Context(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}


func CmdDelegateKeys() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all-delegate-keys",
		Args:  cobra.NoArgs,
		Short: "", // TODO(levi) provide short description
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			req := types.DelegateKeysRequest{}

			res, err := queryClient.DelegateKeys(cmd.Context(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}


func newContextAndQueryClient(cmd *cobra.Command) (client.Context, types.QueryClient, error) {
	clientCtx, err := client.GetClientQueryContext(cmd)
	if err != nil {
		return clientCtx, nil, err
	}
	return clientCtx, types.NewQueryClient(clientCtx), nil
}

func parseContractAddress(s string) (string, error) {
	if !common.IsHexAddress(s) {
		return "", fmt.Errorf("%s not a valid contract address, please input a valid contract address", s)
	}
	return s, nil
}

func parseCount(s string) (int64, error) {
	count, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("count %s not a valid int, please input a valid count", s)
	}
	return count, nil
}

func parseNonce(s string) (uint64, error) {
	nonce, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("nonce %s not a valid uint, please input a valid nonce", s)
	}
	return nonce, nil
}
