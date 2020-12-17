package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/althea-net/peggy/module/x/peggy/client/cli"
	"github.com/althea-net/peggy/module/x/peggy/client/rest"
)

var ProposalPeggyBootstrapHandler = govclient.NewProposalHandler(cli.GetCmdSubmitPeggyBootstrapProposal, rest.ProposalPeggyBootstrapRESTHandler)
var ProposalPeggyUpgradeHandler = govclient.NewProposalHandler(cli.GetCmdSubmitPeggyUpgradeProposal, rest.ProposalPeggyUpgradeRESTHandler)
