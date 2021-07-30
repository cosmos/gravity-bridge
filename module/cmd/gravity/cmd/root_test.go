package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/cli"
)

type KeyOutput struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Address string `json:"address"`
	PubKey  string `json:"pubkey"`
}

func TestKeyGen(t *testing.T) {
	mnemonic := "weasel lunch attack blossom tone drum unfair worry risk level negative height sight nation inside task oyster client shiver aware neck mansion gun dune"

	// generate key from binary
	keyCmd := keys.AddKeyCommand()
	keyCmd.Flags().String(cli.OutputFlag, "json", "output flag")
	keyCmd.Flags().String(flags.FlagKeyringBackend, keyring.BackendTest, "Select keyring's backend (os|file|kwallet|pass|test|memory)")
	keyCmd.SetArgs([]string{
		"--dry-run=true",
		"--output=json",
		"--recover=true",
		"orch",
	})
	keyCmd.SetIn(strings.NewReader(mnemonic + "\n"))

	buf := bytes.NewBuffer(nil)
	keyCmd.SetOut(buf)
	keyCmd.SetErr(buf)

	err := Execute(keyCmd)
	require.NoError(t, err)

	var key KeyOutput
	output := buf.Bytes()
	t.Log("outputs: ", string(output))
	err = json.Unmarshal(output, &key)
	require.NoError(t, err)

	// generate a memory key directly
	kb, err := keyring.New("testnet", keyring.BackendMemory, "", nil)
	if err != nil {
		return
	}

	keyringAlgos, _ := kb.SupportedAlgorithms()
	algo, err := keyring.NewSigningAlgoFromString(string(hd.Secp256k1Type), keyringAlgos)
	if err != nil {
		return
	}

	account, err := kb.NewAccount("", mnemonic, "", "m/44'/118'/0'/0/0", algo)
	require.NoError(t, err)

	require.Equal(t, account.GetAddress().String(), key.Address)
}
