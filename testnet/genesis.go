package main

import (
	"encoding/json"
	"time"

	"github.com/cosmos/cosmos-sdk/x/auth/types"
	types2 "github.com/cosmos/cosmos-sdk/x/bank/types"
	gravitytypes "github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

type Auth struct {
	Params   map[string]string   `json:"params"`
	Accounts []types.BaseAccount `json:"accounts"`
}

type DenomUnit struct {
	Denom    string            `json:"denom"`
	Exponent uint              `json:"exponent"`
	Aliases  []json.RawMessage `json:"aliases"`
}

type DenomMetadata struct {
	Description string      `json:"description"`
	Display     string      `json:"display"`
	Base        string      `json:"base"`
	DenomUnits  []DenomUnit `json:"denom_units"`
}

type Bank struct {
	Params   map[string]json.RawMessage `json:"params"`
	Balances []types2.Balance           `json:"balances"`

	Supply        []json.RawMessage `json:"supply"`
	DenomMetadata []DenomMetadata   `json:"denom_metadata"`
}

type GenUtil struct {
	GenTxs []json.RawMessage `json:"gen_txs"`
}

type AppState struct {
	Auth Auth `json:"auth"`
	Bank Bank `json:"bank"`

	GenUtil GenUtil                   `json:"genutil"`
	Gravity gravitytypes.GenesisState `json:"gravity"`

	Capability   map[string]json.RawMessage `json:"capability"`
	Crisis       map[string]json.RawMessage `json:"crisis"`
	Distribution map[string]json.RawMessage `json:"distribution"`
	Evidence     map[string]json.RawMessage `json:"evidence"`
	Gov          map[string]json.RawMessage `json:"gov"`
	IBC          map[string]json.RawMessage `json:"ibc"`
	Mint         map[string]json.RawMessage `json:"mint"`
	Params       map[string]json.RawMessage `json:"params"`
	Slashing     map[string]json.RawMessage `json:"slashing"`
	Staking      map[string]json.RawMessage `json:"staking"`
	Transfer     map[string]json.RawMessage `json:"transfer"`
	Upgrade      map[string]json.RawMessage `json:"upgrade"`
	Vesting      map[string]json.RawMessage `json:"vesting"`
}

type GenesisState struct {
	GenesisTime   time.Time `json:"genesis_time"`
	ChainID       string    `json:"chain_id"`
	InitialHeight uint64    `json:"initial_height"`

	AppHash  string   `json:"app_hash"`
	AppState AppState `json:"app_state"`

	// These will remain as json.RawMessage until we get a chance to flesh them out.
	// They should be explicitly defined.
	ConsensusParams map[string]json.RawMessage `json:"consensus_params"`
}
