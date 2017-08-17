package commands

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	tcmd "github.com/tendermint/tendermint/cmd/tendermint/commands"
    basecmd "github.com/tendermint/basecoin/cmd/basecoin/commands"
)

//commands
var (
	InitCmd = &cobra.Command{
		Use:   "init [address]",
		Short: "Initialize a basecoin blockchain",
		RunE:  initCmd,
	}
)

//flags
var (
	chainIDFlag string
)

func init() {
	flags := []basecmd.Flag2Register{
		{&chainIDFlag, "chain-id", "test_chain_id", "Chain ID"},
	}
	basecmd.RegisterFlags(InitCmd, flags)
}

// returns 1 iff it set a file, otherwise 0 (so we can add them)
func setupFile(path, data string, perm os.FileMode) (int, error) {
	_, err := os.Stat(path)
	if !os.IsNotExist(err) { //note, os.IsExist(err) != !os.IsNotExist(err)
		return 0, nil
	}
	err = ioutil.WriteFile(path, []byte(data), perm)
	if err != nil {
		return 0, err
	}
	return 1, nil
}

func initCmd(cmd *cobra.Command, args []string) error {
	// this will ensure that config.toml is there if not yet created, and create dir
	cfg, err := tcmd.ParseConfig()
	if err != nil {
		return err
	}

	if len(args) != 1 {
		return fmt.Errorf("`init` takes one argument, a basecoin account address. Generate one using `basecli keys new mykey`")
	}
	userAddr := args[0]
	// verify this account is correct
	data, err := hex.DecodeString(basecmd.StripHex(userAddr))
	if err != nil {
		return errors.Wrap(err, "Invalid address")
	}
	if len(data) != 20 {
		return errors.New("Address must be 20-bytes in hex")
	}

	// initalize basecoin
    genesisFile := cfg.GenesisFile() 
	privValFile := cfg.PrivValidatorFile()
	keyFile := path.Join(cfg.RootDir, "key.json")

    mod1, err := setupFile(genesisFile, GetGenesisJSON(chainIDFlag, userAddr), 0644)
    if err != nil {
        return err
	}
	mod2, err := setupFile(privValFile, PrivValJSON, 0400)
	if err != nil {
		return err
	}
	mod3, err := setupFile(keyFile, KeyJSON, 0400)
	if err != nil {
		return err
	}

	if (mod1 + mod2 + mod3) > 0 {
		msg := fmt.Sprintf("Initialized %s", cmd.Root().Name())
		logger.Info(msg, "genesis", genesisFile, "priv_validator", privValFile)
	} else {
		logger.Info("Already initialized", "priv_validator", privValFile)
	}

	return nil
}

var PrivValJSON = `{
  "address": "FC2582CC198FA734A751EABDCF8C0977B786D3F9",
  "last_height": 0,
  "last_round": 0,
  "last_signature": null,
  "last_signbytes": "",
  "last_step": 0,
  "priv_key": {
    "type": "secp256k1",
    "data": "13D0B182D12A2C9F0E5510D97E7B8F30A2901041711EBCB67F4552ADE94A67CA"
  },
  "pub_key": {
    "type": "secp256k1",
    "data": "03F21BF92B81FBDC29430C64456F9722E6D84A86F738ABA3CEBDA11539C3BA6841"
  }
}`

// GetGenesisJSON returns a new tendermint genesis with Basecoin app_options
// that grant a large amount of "mycoin" to a single address
// TODO: A better UX for generating genesis files
func GetGenesisJSON(chainID, addr string) string {
	return fmt.Sprintf(`{
  "app_hash": "",
  "chain_id": "%s",
  "genesis_time": "0001-01-01T00:00:00.000Z",
  "validators": [
    {
      "amount": 10,
      "name": "",
      "pub_key": {
        "type": "secp256k1",
        "data": "03F21BF92B81FBDC29430C64456F9722E6D84A86F738ABA3CEBDA11539C3BA6841"
      }
    }
  ],
  "app_options": {
    "accounts": [{
      "address": "%s",
      "coins": [
        {
          "denom": "mycoin",
          "amount": 9007199254740992
        }
      ]
    }]
  }
}`, chainID, addr)
}

// TODO: remove this once not needed for relay
var KeyJSON = `{
  "address": "FC2582CC198FA734A751EABDCF8C0977B786D3F9",
  "priv_key": {
    "type": "secp256k1",
    "data": "13D0B182D12A2C9F0E5510D97E7B8F30A2901041711EBCB67F4552ADE94A67CA"
  },
  "pub_key": { 
    "type": "secp256k1",
    "data": "03F21BF92B81FBDC29430C64456F9722E6D84A86F738ABA3CEBDA11539C3BA6841"
  }
}`
