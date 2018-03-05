package witness 

import (  
    "os"
    "testing"
    "encoding/hex"
    "encoding/json"
    
    sdk "github.com/cosmos/cosmos-sdk/types"
    bam "github.com/cosmos/cosmos-sdk/baseapp"
    "github.com/cosmos/cosmos-sdk/x/bank"
    "github.com/cosmos/cosmos-sdk/x/auth"

    wire "github.com/tendermint/go-wire"
    crypto "github.com/tendermint/go-crypto"
    abci "github.com/tendermint/abci/types"

    dbm "github.com/tendermint/tmlibs/db"
    "github.com/tendermint/tmlibs/log"

    "github.com/stretchr/testify/assert"
)

func newHandler(am sdk.AccountMapper, wmapkey *sdk.KVStoreKey) sdk.Handler {
    wmap := NewWitnessMsgMapper(wmapkey)
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
    }

    app.accountMapper = auth.NewAccountMapperSealed(
            app.capKeyMainStore,
            &auth.BaseAccount{},    
    )

    app.Router().AddRoute("witness", newHandler(app.accountMapper, app.capKeyWmapStore))
    
    app.SetTxDecoder(func(txBytes []byte) (sdk.Tx, sdk.Error) {
        tx := sdk.StdTx{}
        err := app.cdc.UnmarshalBinary(txBytes, &tx)
        if err != nil {
            return nil, sdk.ErrTxParse("").TraceCause(err, "")
        }
        return tx, nil
    })

    app.SetInitChainer(func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
        stateJSON := req.AppStateBytes

        var genesisAccounts []auth.BaseAccount
        err := json.Unmarshal(stateJSON, &genesisAccounts)
        if err != nil {
            panic(err)
        }

        for _, acc := range genesisAccounts {
            app.accountMapper.SetAccount(ctx, &acc)
        }

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
    dest, err := hex.DecodeString("6ab9116baa66282b4dfed248f4cac595a41f4a19")
    if err != nil {
        panic(err)
    }

    ether, err := hex.DecodeString("0000000000000000000000000000000000000000")
    if err != nil {
        panic(err)
    }

    return WitnessMsg {
        Info: LockInfo {
            Destination: dest,
            Amount:      1, //(wei)
            Token:       ether,
            Nonce:       0,
        },
        Signer:      signer,
    }
}

func newTx(priv crypto.PrivKey) sdk.StdTx {
    addr := priv.PubKey().Address()
    msg := newLockMsg(addr)
    return sdk.NewStdTx(msg, []sdk.StdSignature{
        sdk.StdSignature {
            PubKey: priv.PubKey(),
            Signature: priv.Sign(msg.GetSignBytes()),
        },
    })
}

func incSequence(tx sdk.StdTx) sdk.StdTx {
    res := tx
    res.Signatures[0].Sequence++;
    return res
}

func TestLockMsg(t *testing.T) {
    logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "test")
    db := dbm.NewMemDB() 
    app := newTestApp(logger, db)

    priv1 := crypto.GenPrivKeyEd25519()
    priv2 := crypto.GenPrivKeyEd25519() 

    genesisState := []auth.BaseAccount{
        auth.BaseAccount{
            Address: priv1.PubKey().Address(),
            Coins:   nil,
        },
        auth.BaseAccount{
            Address: priv2.PubKey().Address(),
            Coins:   nil,
        },
    }
    stateBytes, err := json.MarshalIndent(genesisState, "", "\t")
    if err != nil {
        panic(err)
    }
    app.InitChain(abci.RequestInitChain{[]abci.Validator{}, stateBytes})

    app.BeginBlock(abci.RequestBeginBlock{})

    tx1 := newTx(priv1)

    res := app.Check(tx1)
    assert.Equal(t, sdk.CodeOK, res.Code, "Should pass: app.Check(tx1)")

    res = app.Deliver(tx1)
    assert.Equal(t, sdk.CodeOK, res.Code, "Should pass: app.Deliver(tx1)")

    tx1 = incSequence(tx1)
    res = app.Deliver(tx1)
    assert.Equal(t, CodeWitnessReplay, res.Code, "Should not pass: app.Deliver(tx1)")

    tx2 := newTx(priv2)
    res = app.Deliver(tx2)
    assert.Equal(t, sdk.CodeOK, res.Code, "Should pass: app.Deliver(tx2)")

    tx2 = incSequence(tx2)
    res = app.Deliver(tx2)
    assert.Equal(t, CodeAlreadyCredited, res.Code, "Should not pass: app.Deliver(tx2)") 
}
