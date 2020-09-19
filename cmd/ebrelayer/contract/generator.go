package contract

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

const (
	SolcCmdText   = "[SOLC_CMD]"
	DirectoryText = "[DIRECTORY]"
	ContractText  = "[CONTRACT]"
)

var (
	// BaseABIBINGenCmd is the base command for contract compilation to ABI and BIN
	BaseABIBINGenCmd = strings.Join([]string{"solc ",
		fmt.Sprintf("--%s ./testnet-contracts/contracts/%s%s.sol ", SolcCmdText, DirectoryText, ContractText),
		fmt.Sprintf("-o ./cmd/ebrelayer/contract/generated/%s/%s ", SolcCmdText, ContractText),
		"--overwrite ",
		"--allow-paths *,"},
		"")
	// BaseBindingGenCmd is the base command for contract binding generation
	BaseBindingGenCmd = strings.Join([]string{"abigen ",
		fmt.Sprintf("--bin ./cmd/ebrelayer/contract/generated/bin/%s/%s.bin ", ContractText, ContractText),
		fmt.Sprintf("--abi ./cmd/ebrelayer/contract/generated/abi/%s/%s.abi ", ContractText, ContractText),
		fmt.Sprintf("--pkg %s ", ContractText),
		fmt.Sprintf("--type %s ", ContractText),
		fmt.Sprintf("--out ./cmd/ebrelayer/contract/generated/bindings/%s/%s.go", ContractText, ContractText)},
		"")
)

// CompileContracts compiles contracts to BIN and ABI files
func CompileContracts(contracts BridgeContracts) error {
	for _, contract := range contracts {
		// Construct generic BIN/ABI generation cmd with contract's directory path and name
		baseDirectory := ""
		if contract.String() == BridgeBank.String() {
			baseDirectory = contract.String() + "/"
		}
		dirABIBINGenCmd := strings.Replace(BaseABIBINGenCmd, DirectoryText, baseDirectory, -1)
		contractABIBINGenCmd := strings.Replace(dirABIBINGenCmd, ContractText, contract.String(), -1)

		// Segment BIN and ABI generation cmds
		contractBINGenCmd := strings.Replace(contractABIBINGenCmd, SolcCmdText, "bin", -1)
		err := execCmd(contractBINGenCmd)
		if err != nil {
			return err
		}

		contractABIGenCmd := strings.Replace(contractABIBINGenCmd, SolcCmdText, "abi", -1)
		err = execCmd(contractABIGenCmd)
		if err != nil {
			return err
		}
	}
	return nil
}

// GenerateBindings generates bindings for each contract
func GenerateBindings(contracts BridgeContracts) error {
	for _, contract := range contracts {
		genBindingCmd := strings.Replace(BaseBindingGenCmd, ContractText, contract.String(), -1)
		// genBindingCmd = strings.Replace(genBindingCmd, Dir, strings.ToLower(contract.String()), -1)
		err := execCmd(genBindingCmd)
		if err != nil {
			return err
		}
	}
	return nil
}

// execCmd executes a bash cmd
func execCmd(cmd string) error {

	mainCmd := exec.Command("sh", "-c", cmd)

	var out bytes.Buffer
	var stderr bytes.Buffer

	mainCmd.Stdout = &out
	mainCmd.Stderr = &stderr

	err := mainCmd.Run()

	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}

	fmt.Println("Result: Successfully!")

	return nil
}
