package contract

import (
  "testing"
  "io/ioutil"
  "strings"
  "log"

  "github.com/stretchr/testify/require"
)

// Set up data for parameters and to compare against
func TestLoadABI(t *testing.T) {

	//Get the ABI ready
	rawContractAbi, errorMsg := ioutil.ReadFile("./PeggyABI.json")
	if errorMsg != nil {
		log.Fatal(errorMsg)
	}

	require.True(t, strings.Contains(string(rawContractAbi), "LogLock"))
}