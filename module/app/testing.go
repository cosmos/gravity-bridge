package app

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/simapp"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

// Setup initializes a new SimApp. A Nop logger is set in SimApp.
func Setup(isCheckTx bool) *Gravity {
	gravityApp, genesisState := setup(!isCheckTx, 5)
	if !isCheckTx {
		// init chain must be called to stop deliverState from being nil
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		// Initialize the chain
		gravityApp.InitChain(
			abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: simapp.DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	return gravityApp
}

func setup(withGenesis bool, invCheckPeriod uint) (*Gravity, GenesisState) {
	db := dbm.NewMemDB()
	encCdc := MakeEncodingConfig()
	gravityApp := NewGravityApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, DefaultNodeHome, invCheckPeriod, encCdc, simapp.EmptyAppOptions{})
	if withGenesis {
		return gravityApp, NewDefaultGenesisState()
	}
	return gravityApp, GenesisState{}
}
