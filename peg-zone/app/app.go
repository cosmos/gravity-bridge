package app

import (
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/peggy/peg-zone/x/withdraw"
	"github.com/cosmos/peggy/peg-zone/x/witness"
)

type PeggyApp struct {
	*bam.BaseApp

	capKeyMainStore     *sdk.KVStoreKey
	capKeyWitnessStore  *sdk.KVStoreKey
	capKeyWithdrawStore *sdk.KVStoreKey

	accountMapper    sdk.AccountMapper
	witnessTxMapper  witness.WitnessTxMapper
	withdrawTxMapper withdraw.WithdrawTxMapper
}

func NewPeggy() *PeggyApp {
	app := &PeggyApp{
		BaseApp:             bam.NewBaseApp("Peggy"),
		capKeyMainStore:     sdk.NewKVStoreKey("main"),
		capKeyWitnessStore:  sdk.NewKVStoreKey("witness"),
		capKeyWithdrawStore: sdk.NewKVStoreKey("withdraw"),
	}

	app.accountMapper = auth.NewAccountMapperSealed(
		app.capKeyMainStore, // target store
		&types.AppAccount{}, // prototype
	)

	app.witnessTxMapper = types.NewWitnessTxMapper(app.capKeyWitnessStore)
	app.withdrawTxMapper = types.NewWithdrawTxMapper(app.capKeyWithdrawStore)

	bApp := bam.NewBaseApp(appName)
	mountMultiStore(bApp, app.capKeyMainStore, app.capKeyWithdrawStore, app.capKeyWitnessStore)
	err := bApp.LoadLatestVersion(mainKey)
	if err != nil {
		panic(err)
	}

	// register routes on new application
	accts := types.AccountMapper(mainKey)
	// TODO: Pass WithdrawTxMapper and WitnessTxMapper
	types.RegisterRoutes(bApp.Router(), accts)

	// set up ante and tx parsing
	setAnteHandler(bApp, accts)
	initBaseAppTxDecoder(bApp)

	return &PeggyApp{
		BaseApp: bApp,
		accts:   accts,
	}
}

func mountMultiStore(bApp *baseapp.BaseApp, keys ...*sdk.KVStoreKey) {
	// create substore for every key
	for _, key := range keys {
		bApp.MountStore(key, sdk.StoreTypeIAVL)
	}
}

func setAnteHandler(bApp *baseapp.BaseApp, accts sdk.AccountMapper) {
	// this checks auth, but may take fee is future, check for compatibility
	bApp.SetDefaultAnteHandler(
		auth.NewAnteHandler(accts))
}

func initBaseAppTxDecoder(bApp *baseapp.BaseApp) {
	cdc := types.MakeTxCodec()
	bApp.SetTxDecoder(func(txBytes []byte) (sdk.Tx, sdk.Error) {
		var protoTx types.WitnessTx
		if err := proto.Unmarshal(txBytes, &protoTx); err != nil {
			return sdk.Tx{}, err
		}

		var tx sdk.Tx
		tx.Signatures = []sdk.StdSignature{
			Signature: protoTx.Signature,
			Sequence:  protoTx.Sequence,
		}
		switch innerTx := protoTx.Tx.(type) {
		case WitnessTx_Lock:
			lock := innerTx.Lock
			msg := types.WitnessTx{
				Amount:      lock.Value,
				Destination: lock.Dest,
				Token:       lock.Token,
			}
			tx.Msg = msg
		default:
			return sdk.Tx{}, errors.New("Not implemented")
		}

		// StdTx.Msg is an interface whose concrete
		// types are registered in app/msgs.go.
		//		err := cdc.UnmarshalBinary(txBytes, &tx)
		//		if err != nil {
		//			return nil, sdk.ErrTxParse("").TraceCause(err, "")
		//		}
		return tx, nil
	})
}

/*
func initEndBlocker(bApp *baseapp.BaseApp) {
    bApp.SetEndBlocker(func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {

    })
}*/
