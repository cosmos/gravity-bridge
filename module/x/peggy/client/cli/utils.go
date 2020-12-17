package cli

import (
	"io/ioutil"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/althea-net/peggy/module/x/peggy/types"
)

func ParsePeggyBootstrapProposalWithDeposit(cdc codec.JSONMarshaler, proposalFile string) (types.PeggyBootstrapProposalWithDeposit, error) {
	proposal := types.PeggyBootstrapProposalWithDeposit{}

	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}

func ParsePeggyUpgradeProposalWithDeposit(cdc codec.JSONMarshaler, proposalFile string) (types.PeggyUpgradeProposalWithDeposit, error) {
	proposal := types.PeggyUpgradeProposalWithDeposit{}

	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
