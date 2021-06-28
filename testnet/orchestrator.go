package main

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
)

type Orchestrator struct {
	Chain *Chain
	Index uint8

	// Key management
	Mnemonic string
	KeyInfo  keyring.Info
}


func (o *Orchestrator) instanceName() string {
	return fmt.Sprintf("orchestrator%d",  o.Index)
}