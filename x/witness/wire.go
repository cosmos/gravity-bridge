package witness

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    wire "github.com/tendermint/go-wire"
)

func RegisterWire(cdc *wire.Codec) {
    cdc.RegisterConcrete(WitnessTx{},
        "com.cosmos.peggy.WitnessTx", nil)
}

// used in initBaseAppTxDecoder
func TxDecoder(txBytes []byte) (sdk.Tx, sdk.Error) {
    var protoTx WitnessTx
    if err := proto.Unmarshal(txBytes, &protoTx); err != nil {
        return sdk.Tx{}, err
    }

    var tx sdk.StdTx

    tx.Signatures = []sdk.StdSignature {
        Signature: protoTx.Signature,
        Sequence:  protoTx.Sequence,
    }

    switch innerTx := protoTx.Tx.(type) {
    case WitnessTx_Lock:
        lock := innerTx.Lock
        msg := WitnessMsg {
            Amount:      lock.Value,
            Destination: lock.Dest,
            Token:       lock.Token,
        }
        tx.msg = msg
    /*
    case WitnessTx_Burn:
        burn := innerTx.Burn
        msg := BurnMsg {
            Amount:      burn.Value,
            Destination: burn.Dest,
            Token:       burn.Token,
            Nonce:       burn.Nonce,
        }
        tx.msg = msg
    */
    default: 
        return sdk.Tx{}, errors.New("Not implemented")
    }
    return tx, nil
}
