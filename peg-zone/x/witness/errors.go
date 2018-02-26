package witness

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
)

type CodeType = sdk.CodeType

const (
    CodeAlreadyCredited CodeType = 101
    CodeWitnessReplay   CodeType = 102
)

func codeToDefaultMsg(code CodeType) string {
    switch code {
    case CodeWitnessReplay:
        return "Witness tx replayed"
    default:
        return sdk.CodeToDefaultMsg(code)
    }
}

func ErrAlreadyCredited() sdk.Error {
    return newError(CodeAlreadyCredited, "")
}

func ErrWitnessReplay() sdk.Error {
    return newError(CodeWitnessReplay, "")
}

func newError(code CodeType, msg string) sdk.Error {
    msg = msgOrDefaultMsg(msg, code)
    return sdk.NewError(code, msg)
}
