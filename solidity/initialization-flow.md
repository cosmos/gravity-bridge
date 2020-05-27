0. There is a governance resolution accepting the peggy code into the full node codebase
   - Once peggy starts running on a validator, it generates an Ethereum keypair.
   - Peggy publishes a validator's Ethereum address every single block.
   - A peggyId is chosen at this step.
   - The source code hash of the peggy Ethereum contract is saved here.
1. There is a governance resolution that says "we are going to start peggy at block x"
   - This is a parameter-changing resolution that the peggy module looks for
1. Right after block x, the peggy module looks at the validator set at block x, and signs over it using its Ethereum keypair.
1. Peggy puts the Ethereum signature from the last step into the consensus state. As part of consensus, the validators check that each of these signatures is valid.
   - The Eth signatures over the Eth addresses of the validator set from the last block are now required in every block going forward.
1. Now, the deployer script hits a full node api, gets the Eth signatures of the valset from the latest block, and deploys the Ethereum contract.
   - We will consider the scenario that many deployers deploy many valid peggy eth contracts.
1. The deployer submits the address of the peggy contract that it deployed to Ethereum.
   - The peggy module checks the Ethereum chain for each submitted address, and makes sure that the peggy contract at that address is using the correct source code, and has the correct validator set.
1. There is a rule in the peggy module that the correct peggy eth contract address with the lowest address in the earliest block that is at least n blocks back is the "official" contract. After 50 blocks passes, this implies that anyone wanting to use peggy knows which contract to send their money to.
