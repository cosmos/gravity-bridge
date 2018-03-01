package witness

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
)

type CodeType = sdk.CodeType

const (
    CodeInvalidWitnessMsg  CodeType = 101
    CodeAlreadyCredited CodeType = 201
    CodeWitnessReplay   CodeType = 202
    
)

func codeToDefaultMsg(code CodeType) string {
    switch code {
    case CodeWitnessReplay:
        return "Witness tx replayed"
    default:
        return sdk.CodeToDefaultMsg(code)
    }
}

func ErrInvalidWitnessMsg() sdk.Error {
    return newError(CodeInvalidWitnessMsg, "")
}

func ErrAlreadyCredited() sdk.Error {
    return newError(CodeAlreadyCredited, "")
}

func ErrWitnessReplay() sdk.Error {
    return newError(CodeWitnessReplay, "")
}

func msgOrDefaultMsg(msg string, code CodeType) string {
    if msg != "" {
        return msg
    } else {
        return codeToDefaultMsg(code)
    }
}

func newError(code CodeType, msg string) sdk.Error {
    msg = msgOrDefaultMsg(msg, code)
    return sdk.NewError(code, msg)
}
