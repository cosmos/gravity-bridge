package commands

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path"
    "path/filepath"

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

    cfgpath := filepath.Join(cfg.Consensus.RootDir, "config.toml")

    cfgfile, err := ioutil.ReadFile(cfgpath)
    if err != nil {
        return err
    }

    cfgfile = append(cfgfile, []byte("\n[consensus]\ntimeout_commit = 10000\ncreate_empty_blocks = false\ncreate_empty_blocks_interval = 1800\n")...) // 30 min

    ioutil.WriteFile(cfgpath, cfgfile, 0644)

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
        fmt.Printf("%s\n", msg)
//		logger.Info(msg, "genesis", genesisFile, "priv_validator", privValFile)
	} else {
        fmt.Printf("Already initialized\n")
//		logger.Info("Already initialized", "priv_validator", privValFile)
	}

	return nil
}

var PrivValJSON = `{
  "address": "5683A1E28A791BADE84A7A1D562241E314D61DAF",
  "last_height": 0,
  "last_round": 0,
  "last_signature": null,
  "last_signbytes": "",
  "last_step": 0,
  "priv_key": {
    "type": "secp256k1",
    "data": "F1C075541D98BC45F9DFA19EA2691E61DE9E347B83AA65935D1EA90BE0369D80"
  },
  "pub_key": {
    "type": "secp256k1",
    "data": "0276954FB8E6E2303CA72BC082D5F5EE89A4C50451F277197C9236E5931D779774"
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
        "data": "0276954FB8E6E2303CA72BC082D5F5EE89A4C50451F277197C9236E5931D779774"
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
  "address": "5683A1E28A791BADE84A7A1D562241E314D61DAF",
  "priv_key": {
    "type": "secp256k1",
    "data": "F1C075541D98BC45F9DFA19EA2691E61DE9E347B83AA65935D1EA90BE0369D80"
  },
  "pub_key": { 
    "type": "secp256k1",
    "data": "0276954FB8E6E2303CA72BC082D5F5EE89A4C50451F277197C9236E5931D779774"
  }
}`
