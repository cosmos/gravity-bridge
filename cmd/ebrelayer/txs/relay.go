package txs

// ------------------------------------------------------------
//      Relay
//
//      Builds and encodes EthBridgeClaim Msgs with the
//      specified variables, before presenting the unsigned
//      transaction to validators for optional signing.
//      Once signed, the data packets are sent as transactions
//      on the Cosmos Bridge.
// ------------------------------------------------------------

import (
	"log"
	"os"
	"os/exec"
	"strconv"
)

func RelayEvent(claim *WitnessClaim) error {

	// Cast to string
	nonce 				 := strconv.Itoa(claim.Nonce)
	ethereumSender := claim.EthereumSender
	cosmosReceiver := claim.CosmosReceiver.String()
	validator 		 := claim.Validator.String()
	amount 				 := claim.Amount.String()

	// Build the ebcli tx command
	cmd := exec.Command("ebcli tx ethbridge make-claim",
											nonce, ethereumSender, cosmosReceiver, validator, amount)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the cmd
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

	return nil
}
