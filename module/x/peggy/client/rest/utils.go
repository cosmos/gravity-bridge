package rest

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

type PeggyBootstrapProposalReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

	Title                string         `json:"title" yaml:"title"`
	Description          string         `json:"description" yaml:"description"`
	PeggyID              string         `json:"peggy_id" yaml:"peggy_id"`
	ProxyContractHash    string         `json:"proxy_contract_hash" yaml:"proxy_contract_hash"`
	ProxyContractAddress string         `json:"proxy_contract_address" yaml:"proxy_contract_address"`
	LogicContractHash    string         `json:"logic_contract_hash" yaml:"logic_contract_hash"`
	LogicContractAddress string         `json:"logic_contract_address" yaml:"logic_contract_address"`
	StartThreshold       uint64         `json:"start_threshold" yaml:"start_threshold"`
	BridgeChainID        uint64         `json:"bridge_chain_id" yaml:"bridge_chain_id"`
	BootstrapValsetNonce uint64         `json:"bootstrap_valset_nonce" yaml:"bootstrap_valset_nonce"`
	Proposer             sdk.AccAddress `json:"proposer" yaml:"proposer"`
	Deposit              sdk.Coins      `json:"deposit" yaml:"deposit"`
}

type PeggyUpgradeProposalReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

	Title                string         `json:"title" yaml:"title"`
	Description          string         `json:"description" yaml:"description"`
	Version              string         `json:"version" yaml:"version"`
	LogicContractHash    string         `json:"logic_contract_hash" yaml:"logic_contract_hash"`
	LogicContractAddress string         `json:"logic_contract_address" yaml:"logic_contract_address"`
	Proposer             sdk.AccAddress `json:"proposer" yaml:"proposer"`
	Deposit              sdk.Coins      `json:"deposit" yaml:"deposit"`
}
