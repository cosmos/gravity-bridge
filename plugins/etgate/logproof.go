package etgate

import (
    "bytes"
//    "fmt"
    "errors"

    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/rlp"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/trie"   
)


type LogProof struct {
    Receipt []byte // rlp-encoded
    TxIndex uint
//    Index uint // Assuming there is only one log in the receipt
    Proof []rlp.RawValue
    Number uint64
}

func NewLogProof(receipts []*types.Receipt, txIndex uint, /*index uint,*/ number uint64) (LogProof, error) {
    var logReceipt []byte
    keybuf := new(bytes.Buffer)
    trie := new(trie.Trie)

    for i, receipt := range receipts {
        keybuf.Reset()
        rlp.Encode(keybuf, uint(i))

        bytes, err := rlp.EncodeToBytes(receipt)
        if err != nil {
            return LogProof{}, err
        }
        trie.Update(keybuf.Bytes(), bytes)
        
        if txIndex == uint(i) {
            logReceipt = bytes
        }
    }

    if logReceipt == nil {
        return LogProof{}, errors.New("Receipt array does not contain txIndex")
    }

    keybuf.Reset()
    rlp.Encode(keybuf, txIndex)

    return LogProof {
        Receipt: logReceipt,
        TxIndex: txIndex,
//        Index: index,
        Proof: trie.Prove(keybuf.Bytes()),
        Number: number,
    }, nil
}

func (proof LogProof) Log() (types.Log, error) {
    var receipt types.Receipt
    if err := rlp.DecodeBytes(proof.Receipt, &receipt); err != nil {
        return types.Log{}, err
    }/*
    fmt.Printf("len(receipt.Logs): %v, proof.Index: %v\n", len(receipt.Logs), proof.Index)
    if int(proof.Index) >= len(receipt.Logs) {
        return types.Log{}, errors.New("Index out of range while accessing log array")
    }*/
    if len(receipt.Logs) < 1 {
        return types.Log{}, errors.New("No logs found in the receipt")
    }
    return *receipt.Logs[0], nil
}

func (proof LogProof) IsValid(receiptHash common.Hash) bool {
    keybuf := new(bytes.Buffer)
    rlp.Encode(keybuf, proof.TxIndex)
    res, err := trie.VerifyProof(receiptHash, keybuf.Bytes(), proof.Proof)
    if err != nil {
//        fmt.Println("Error in isValid, VerifyProof: ", err)
        return false
    }

    if err != nil || !bytes.Equal(proof.Receipt, res) {
//        fmt.Println("Error in isValid, EncodeToBytes: ", err)
        return false
    }

    return true
}
