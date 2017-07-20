package etgate

import (
    "bytes"
    "fmt"

    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/rlp"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/trie"   
)


type LogProof struct {
    ReceiptHash common.Hash
    Receipt *types.Receipt
    TxIndex uint
    Index uint
    Proof []rlp.RawValue
}

func (proof *LogProof) Log() (types.Log) {
    return *proof.Receipt.Logs[proof.Index]
}

func (proof *LogProof) IsValid() bool {
    keybuf := new(bytes.Buffer)
    rlp.Encode(keybuf, proof.TxIndex)
    res, err := trie.VerifyProof(proof.ReceiptHash, keybuf.Bytes(), proof.Proof)
    if err != nil {
        fmt.Println("Error in isValid, VerifyProof: ", err)
//        fmt.Printf("log: %+v\n", proof.Log())
//        fmt.Printf("proof: %+v\n", proof)
        return false
    }

    rec, err := rlp.EncodeToBytes(proof.Receipt) 
    if err != nil || !bytes.Equal(rec, res) {
        fmt.Println("Error in isValid, EncodeToBytes: ", err)
        return false
    }

    return true
}
