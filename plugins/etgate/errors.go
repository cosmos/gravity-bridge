package etgate

import (
    "fmt"

    abci "github.com/tendermint/abci/types"

    eth "github.com/ethereum/go-ethereum/core/types"
)

var (
    errConflictingChain = fmt.Errorf("Conflicting chain")
    errNoncontinuousHeaderList = fmt.Errorf("Non continuous header list")
    errNotInitialized = fmt.Errorf("Not initialized")

    ETGateCodeConflictingChain = abci.CodeType(1001)
    ETGateCodeNonContinuousHeaderList = abci.CodeType(1002)
)
/*
func ErrConflictingChain(hash common.Hash) error {
    return errors.WithMessage(hash.Hex(), errConflictingChain, ETGateCodeConflictingChain)
}*/

func ErrInvalidLogProof(log eth.Log, err error) error {
    return err
}

