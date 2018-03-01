package witness 

import (  
    "os"
    "testing"
    
    sdk "github.com/cosmos/cosmos-sdk/types"
    bam "github.com/cosmos/cosmos-sdk/baseapp"
    "github.com/cosmos/cosmos-sdk/x/bank"
    "github.com/cosmos/cosmos-sdk/x/auth"

    wire "github.com/tendermint/go-wire"
    crypto "github.com/tendermint/go-crypto"
    abci "github.com/tendermint/abci/types"

    dbm "github.com/tendermint/tmlibs/db"
    "github.com/tendermint/tmlibs/log"
    "github.com/tendermint/tmlibs/common"

    "github.com/stretchr/testify/assert"
)

func newHandler(am sdk.AccountMapper) sdk.Handler {
    wmap := NewWitnessMsgMapper(sdk.NewKVStoreKey("witness"))
    ck := bank.NewCoinKeeper(am)
    return NewHandler(wmap, ck)
}

func makeTxCodec() *wire.Codec {
    cdc := wire.NewCodec()
    crypto.RegisterWire(cdc)
    bank.RegisterWire(cdc)
    RegisterWire(cdc)
    return cdc
}

type TestApp struct {
    *bam.BaseApp

    cdc *wire.Codec

    capKeyMainStore *sdk.KVStoreKey
    capKeyWmapStore *sdk.KVStoreKey

    accountMapper sdk.AccountMapper
}

func newTestApp(logger log.Logger, db dbm.DB) *TestApp {
    app := &TestApp {
        BaseApp:         bam.NewBaseApp("TestApp", logger, db),
        cdc:             makeTxCodec(),
        capKeyMainStore: sdk.NewKVStoreKey("main"),
        capKeyWmapStore: sdk.NewKVStoreKey("witness"),
        accountMapper:   auth.NewAccountMapperSealed(
            sdk.NewKVStoreKey("main"),
            &auth.BaseAccount{},
        ),
    }

    app.Router().AddRoute("witness", newHandler(app.accountMapper))
    
    app.SetTxDecoder(func(txBytes []byte) (sdk.Tx, sdk.Error) {
        tx := sdk.StdTx{}
        err := app.cdc.UnmarshalBinary(txBytes, &tx)
        if err != nil {
            return nil, sdk.ErrTxParse("").TraceCause(err, "")
        }
        return tx, nil
    })

    app.SetInitChainer(func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
        return abci.ResponseInitChain{}
    })

    app.MountStoresIAVL(app.capKeyMainStore, app.capKeyWmapStore)
    app.SetAnteHandler(auth.NewAnteHandler(app.accountMapper))
    err := app.LoadLatestVersion(app.capKeyMainStore)
    if err != nil {
        panic(err)
    }

    return app
}

func newLockMsg(signer crypto.Address) WitnessMsg {
    var dest *common.HexBytes
    dest.UnmarshalJSON([]byte("0x6ab9116baa66282b4dfed248f4cac595a41f4a19"))

    var ether *common.HexBytes
    ether.UnmarshalJSON([]byte("0x0000000000000000000000000000000000000000"))

    return WitnessMsg {
        Info: LockInfo {
            Destination: *dest,
            Amount:      1, //(wei)
            Token:       *ether,
        },
        Signer:      signer,
    }
}

func newTx() sdk.StdTx {
    priv := crypto.GenPrivKeyEd25519()
    addr := priv.PubKey().Address()
    msg := newLockMsg(addr)
    return sdk.NewStdTx(msg, []sdk.StdSignature{
        sdk.StdSignature {
            PubKey: priv.PubKey(),
            Signature: priv.Sign(msg.GetSignBytes()),
        },
    })
}

func TestLockMsg(t *testing.T) {
    logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module, test")
    db := dbm.NewMemDB() 
    app := newTestApp(logger, db)

    app.BeginBlock(abci.RequestBeginBlock{})

    tx1 := newTx()

    assert.Equal(t, sdk.CodeOK, app.Check(tx1), "Should pass: app.Check(tx1)")
    assert.Equal(t, sdk.CodeOK, app.Deliver(tx1), "Should pass: app.Deliver(tx1)")
    assert.Equal(t, CodeWitnessReplay, app.Deliver(tx1), "Should not pass: app.Deliver(tx1)")

    tx2 := newTx()
    assert.Equal(t, sdk.CodeOK, app.Deliver(tx2), "Should pass: app.Deliver(tx2)")
    assert.Equal(t, CodeAlreadyCredited, app.Deliver(tx2), "Should not pass: app.Deliver(tx2)") 
}
