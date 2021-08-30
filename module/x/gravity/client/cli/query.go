package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
	"github.com/ethereum/go-ethereum/common"
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
		CmdBatchTx(),
		CmdBatchTxConfirmations(),
		CmdBatchTxFees(),
		CmdBatchTxs(),
		CmdContractCallTx(),
		CmdContractCallTxConfirmations(),
		CmdContractCallTxs(),
		CmdDenomToERC20Params(),
		CmdERC20ToDenom(),
		CmdLastSubmittedEthereumEvent(),
		CmdLatestSignerSetTx(),
		CmdParams(),
		CmdSignerSetTx(),
		CmdSignerSetTxConfirmations(),
		CmdSignerSetTxs(),
		CmdUnsignedBatchTxs(),
		CmdUnsignedContractCallTxs(),
		CmdUnsignedSignerSetTxs(),
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
		Args:  cobra.ExactArgs(1),
		Short: "query an individual signer set transaction by its nonce",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			nonce, err := parseNonce(args[0])
			if err != nil {
				return err
			}

			res, err := queryClient.SignerSetTx(cmd.Context(), &types.SignerSetTxRequest{SignerSetNonce: nonce})
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
		Args:  cobra.ExactArgs(2),
		Short: "query an outgoing batch by its contract address and nonce",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			contractAddress, err := parseContractAddress(args[0])
			if err != nil {
				return nil
			}

			nonce, err := parseNonce(args[1])
			if err != nil {
				return err
			}

			res, err := queryClient.BatchTx(cmd.Context(), &types.BatchTxRequest{
				TokenContract: contractAddress,
				BatchNonce:    nonce,
			})

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
		Short: "query an outgoing contract call by scope and nonce",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			// TODO: validate this scope somehow
			invalidationScope := []byte(args[0])

			invalidationNonce, err := parseNonce(args[1])
			if err != nil {
				return err
			}

			res, err := queryClient.ContractCallTx(cmd.Context(), &types.ContractCallTxRequest{
				InvalidationScope: invalidationScope,
				InvalidationNonce: invalidationNonce,
			})

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
		Use:   "signer-set-txs",
		Args:  cobra.NoArgs,
		Short: "query all the signer set transactions from the chain",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.SignerSetTxs(cmd.Context(), &types.SignerSetTxsRequest{Pagination: pageReq})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "signer-set-txs")
	return cmd
}

func CmdBatchTxs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batch-txs",
		Args:  cobra.NoArgs,
		Short: "query all the batch transactions from the chain",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.BatchTxs(cmd.Context(), &types.BatchTxsRequest{Pagination: pageReq})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "batch-txs")
	return cmd
}

func CmdContractCallTxs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contract-call-txs",
		Args:  cobra.NoArgs,
		Short: "query all contract call transactions from the chain",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.ContractCallTxs(cmd.Context(), &types.ContractCallTxsRequest{Pagination: pageReq})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "contract-call-txs")
	return cmd
}

func CmdSignerSetTxConfirmations() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "signer-set-tx-ethereum-signatures [nonce]",
		Args:  cobra.ExactArgs(1),
		Short: "query signer set transaction signatures from the validators identified by nonce",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			nonce, err := parseNonce(args[0])
			if err != nil {
				return err
			}

			res, err := queryClient.SignerSetTxConfirmations(cmd.Context(), &types.SignerSetTxConfirmationsRequest{
				SignerSetNonce: nonce,
			})
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
		Args:  cobra.ExactArgs(2),
		Short: "query signatures for a given batch transaction identified by nonce and contract",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			nonce, err := parseNonce(args[0])
			if err != nil {
				return err
			}

			contractAddress, err := parseContractAddress(args[1])
			if err != nil {
				return nil
			}

			res, err := queryClient.BatchTxConfirmations(cmd.Context(), &types.BatchTxConfirmationsRequest{
				BatchNonce:    nonce,
				TokenContract: contractAddress,
			})

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
		Short: "query signatures for a given contract call transaction identified by invalidation nonce and invalidation scope",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			// TODO: some sort of validation here?
			invalidationScope := []byte(args[0])

			invalidationNonce, err := parseNonce(args[1])
			if err != nil {
				return err
			}

			res, err := queryClient.ContractCallTxConfirmations(cmd.Context(), &types.ContractCallTxConfirmationsRequest{
				InvalidationNonce: invalidationNonce,
				InvalidationScope: invalidationScope,
			})

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
		Use:   "pending-signer-set-tx-ethereum-signatures [validator-or-orchestrator-acc-address]",
		Args:  cobra.ExactArgs(1),
		Short: "query signatures for any pending signer set transactions given a validator or orchestrator address (sdk.AccAddress format)",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			address, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			res, err := queryClient.UnsignedSignerSetTxs(cmd.Context(), &types.UnsignedSignerSetTxsRequest{
				Address: address.String(),
			})

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
		Use:   "pending-batch-tx-ethereum-signatures [validator-or-orchestrator-acc-address]",
		Args:  cobra.ExactArgs(1),
		Short: "query signatures for any pending batch transactions given a validator or orchestrator address (sdk.AccAddress format)",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			address, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			res, err := queryClient.UnsignedBatchTxs(cmd.Context(), &types.UnsignedBatchTxsRequest{
				Address: address.String(),
			})

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
		Use:   "pending-contract-call-tx-ethereum-signatures [validator-or-orchestrator-acc-address]",
		Args:  cobra.ExactArgs(1),
		Short: "query signatures for any pending contract call transactions given a validator or orchestrator address (sdk.AccAddress format)",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			address, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			res, err := queryClient.UnsignedContractCallTxs(cmd.Context(), &types.UnsignedContractCallTxsRequest{
				Address: address.String(),
			})

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdLatestSignerSetTx() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "latest-signer-set-tx",
		Args:  cobra.NoArgs,
		Short: "query for the latest signer set from the chain",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			req := &types.LatestSignerSetTxRequest{}

			res, err := queryClient.LatestSignerSetTx(cmd.Context(), req)
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
		Use:   "last-submitted-ethereum-event [validator-or-orchestrator-acc-address]",
		Args:  cobra.ExactArgs(1),
		Short: "query for the last event nonce that was submitted by a given validator",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			address, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			res, err := queryClient.LastSubmittedEthereumEvent(cmd.Context(), &types.LastSubmittedEthereumEventRequest{
				Address: address.String(),
			})

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// TODO: this looks broken
func CmdBatchTxFees() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batch-tx-fees",
		Args:  cobra.NoArgs,
		Short: "query amount of fees for any unrelayed batches",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			res, err := queryClient.BatchTxFees(cmd.Context(), &types.BatchTxFeesRequest{})
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
		Short: "given an erc20 contract address return the cosmos denom",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			contract, err := parseContractAddress(args[0])
			if err != nil {
				return err
			}

			res, err := queryClient.ERC20ToDenom(cmd.Context(), &types.ERC20ToDenomRequest{
				Erc20: contract,
			})

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdDenomToERC20Params() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "denom-to-erc20-params [denom]",
		Args:  cobra.ExactArgs(1),
		Short: "given a cosmos base denom return the correct erc20 name, symbol and decimals",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			if err := sdk.ValidateDenom(args[0]); err != nil {
				return err
			}

			req := &types.DenomToERC20ParamsRequest{
				Denom: args[0],
			}

			res, err := queryClient.DenomToERC20Params(cmd.Context(), req)
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
		Short: "given a cosmos denom return an erc20 contract address",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			if err := sdk.ValidateDenom(args[0]); err != nil {
				return err
			}

			res, err := queryClient.DenomToERC20(cmd.Context(), &types.DenomToERC20Request{
				Denom: args[0],
			})

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
		Short: "query all unbatched send to ethereum messages",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			sender, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			res, err := queryClient.UnbatchedSendToEthereums(cmd.Context(), &types.UnbatchedSendToEthereumsRequest{
				SenderAddress: sender.String(),
				Pagination:    pageReq,
			})

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "unbatched-send-to-ethereums")
	return cmd
}

func CmdDelegateKeysByValidator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegate-keys-by-validator [validator-address]",
		Args:  cobra.ExactArgs(1),
		Short: "query which public keys/addresses a validator has delegated to",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			validatorAddress, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			res, err := queryClient.DelegateKeysByValidator(cmd.Context(), &types.DelegateKeysByValidatorRequest{
				ValidatorAddress: validatorAddress.String(),
			})

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
		Short: "query the valdiator and orchestartor keys for a given ethsigner",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			if !common.IsHexAddress(args[0]) {
				return fmt.Errorf("address is not an etheruem address")
			}

			res, err := queryClient.DelegateKeysByEthereumSigner(cmd.Context(), &types.DelegateKeysByEthereumSignerRequest{
				EthereumSigner: args[0],
			})
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
		Short: "query the validator and eth signer keys for a given orchestrator address",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			orcAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			res, err := queryClient.DelegateKeysByOrchestrator(cmd.Context(), &types.DelegateKeysByOrchestratorRequest{
				OrchestratorAddress: orcAddr.String(),
			})

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
		Short: "query all delegate keys tracked by the chain",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, queryClient, err := newContextAndQueryClient(cmd)
			if err != nil {
				return err
			}

			res, err := queryClient.DelegateKeys(cmd.Context(), &types.DelegateKeysRequest{})
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
