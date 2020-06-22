0. There is a governance resolution accepting the peggy code into the full node codebase
    - Once peggy starts running on a validator, it generates an Ethereum keypair.
        - We may get this from config temporarily
    - Peggy publishes a validator's Ethereum address every single block.
        - This is most likely just going to mean putting it in the peggy keeper, but maybe there needs to be some kind of greater tie in with the rest of the validator set information
            - see staking/keeper/keeper.go for an example of getting individual and all validators
    - A peggyId is chosen at this step.
        - This may be hardcoded for now
    - The source code hash of the peggy Ethereum contract is saved here.
        - This may be hardcoded as well
1. There is a governance resolution that says "we are going to start peggy at block x"
    - This is a parameter-changing resolution that the peggy module looks for
2. Right after block x, the peggy module looks at the validator set at block x, and signs over it using its Ethereum keypair.
3. Peggy puts the Ethereum signature from the last step into the consensus state. As part of consensus, the validators check that each of these signatures is valid.
    - The Eth signatures over the Eth addresses of the validator set from the last block are now required in every block going forward.
4. The deployer script hits a full node api, gets the Eth signatures of the valset from the latest block, and deploys the Ethereum contract.
5. The deployer submits the address of the peggy contract that it deployed to Ethereum.
    - We will consider the scenario that many deployers deploy many valid peggy eth contracts.
    - The peggy module checks the Ethereum chain for each submitted address, and makes sure that the peggy contract at that address is using the correct source code, and has the correct validator set.
6. There is a rule in the peggy module that the correct peggy eth contract address with the lowest address in the earliest block that is at least n blocks back is the "official" contract. After n blocks passes, this allows anyone wanting to use peggy to know which ethereum contract to send their money to.
