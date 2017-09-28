package etgate

const (
    errConflictingChain = fmt.Errorf("Conflicting chain")

    ETGateCodeConflictingChain = abci.CodeType(1001)
)

func ErrConflictingChain(hash common.Hash) error {
    return errors.WithMessage(hash.Hex(), errConflictingChain, ETGateCodeConflictingChain)
}
