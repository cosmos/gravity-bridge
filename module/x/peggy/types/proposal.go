package types

import (
	"fmt"
	"strings"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypePeggyBootstrap = "PeggyBootstrap"
	ProposalTypePeggyUpgrade   = "PeggyUpgrade"
)

var (
	_ govtypes.Content = &PeggyBootstrapProposal{}
	_ govtypes.Content = &PeggyUpgradeProposal{}
)

func NewPeggyBootstrapProposal(
	title string,
	description string,
	peggyID string,
	proxyContractHash string,
	proxyContractAddress string,
	logicContractHash string,
	logicContractAddress string,
	startThreshold uint64,
	bridgeChainID uint64,
	valsetNonce uint64,
) *PeggyBootstrapProposal {
	return &PeggyBootstrapProposal{
		title,
		description,
		peggyID,
		proxyContractHash,
		proxyContractAddress,
		logicContractHash,
		logicContractAddress,
		startThreshold,
		bridgeChainID,
		valsetNonce,
	}
}

// GetTitle returns the title of a community pool spend proposal.
func (pbp *PeggyBootstrapProposal) GetTitle() string { return pbp.Title }

// GetDescription returns the description of a community pool spend proposal.
func (pbp *PeggyBootstrapProposal) GetDescription() string { return pbp.Description }

// GetDescription returns the routing key of a community pool spend proposal.
func (pbp *PeggyBootstrapProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a community pool spend proposal.
func (pbp *PeggyBootstrapProposal) ProposalType() string { return ProposalTypePeggyBootstrap }

// ValidateBasic runs basic stateless validity checks
func (pbp *PeggyBootstrapProposal) ValidateBasic() error {
	// TODO
	return nil
}

// String implements the Stringer interface.
func (pbp PeggyBootstrapProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(
		"Peggy Bootstrap Proposal:\n"+
			"	Title:       			%s\n"+
			"	Description: 			%s\n"+
			"	ProxyContractHash: 		%s\n"+
			"	LogicContractHash: 		%s\n"+
			"	StartThreshold: 		%d\n"+
			"	BridgeChainId: 			%d\n"+
			"	BootstrapValsetNonce: 	%d\n",
		pbp.Title,
		pbp.Description,
		pbp.ProxyContractHash,
		pbp.LogicContractHash,
		pbp.StartThreshold,
		pbp.BridgeChainId,
		pbp.BootstrapValsetNonce,
	))
	return b.String()
}

func NewPeggyUpgradeProposal(
	title string,
	description string,
	version string,
	logicContractHash string,
	logicContractAddress string,
) *PeggyUpgradeProposal {
	return &PeggyUpgradeProposal{
		title,
		description,
		version,
		logicContractHash,
		logicContractAddress,
	}
}

// GetTitle returns the title of a community pool spend proposal.
func (pup *PeggyUpgradeProposal) GetTitle() string { return pup.Title }

// GetDescription returns the description of a community pool spend proposal.
func (pup *PeggyUpgradeProposal) GetDescription() string { return pup.Description }

// GetDescription returns the routing key of a community pool spend proposal.
func (pup *PeggyUpgradeProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a community pool spend proposal.
func (pup *PeggyUpgradeProposal) ProposalType() string { return ProposalTypePeggyUpgrade }

// ValidateBasic runs basic stateless validity checks
func (pup *PeggyUpgradeProposal) ValidateBasic() error {
	// TODO
	return nil
}

// String implements the Stringer interface.
func (pup PeggyUpgradeProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(
		"Peggy Upgrade Proposal:\n"+
			"	Title:       			%s\n"+
			"	Description: 			%s\n"+
			"	Version: 				%s\n"+
			"	LogicContractHash: 		%s\n"+
			"	LogicContractAddress: 	%s\n",
		pup.Title,
		pup.Description,
		pup.Version,
		pup.LogicContractHash,
		pup.LogicContractAddress,
	))
	return b.String()
}
