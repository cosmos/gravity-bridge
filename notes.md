# Installing Peggy on a Cosmos chain

## Installing Peggy on an unstarted chain

0. The deployed contract address is stored in the gentx's that validators generate
   - Validators generate an Ethereum keypair and submit it as part of their gentx
   - Validators also sign over the gentx validator state with their ethereum keypair
1. Once Gentx signing is complete the Ethereum signatures are collected and used to submit a valset update to the already deployed Ethereum contract
2. The chain starts, the validators all observe the Ethereum contract state and find that it matches their current validator state, proceed with updates as normal

## Installing Peggy on a live Cosmos chain

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
   - At some later step, we will need to check that all validators have their eth address in there
1. There is a governance resolution that says "we are going to start peggy at block x"
   - This is a parameter-changing resolution that the peggy module looks for
1. Right after block x, the peggy module looks at the validator set at block x, and signs over it using its Ethereum keypair.
1. Peggy puts the Ethereum signature from the last step into the consensus state. As part of consensus, the validators check that each of these signatures is valid.
   - The Eth signatures over the Eth addresses of the validator set from the last block are now required in every block going forward.
1. The deployer script hits a full node api, gets the Eth signatures of the valset from the latest block, and deploys the Ethereum contract.
1. The deployer submits the address of the peggy contract that it deployed to Ethereum.
   - We will consider the scenario that many deployers deploy many valid peggy eth contracts.
   - The peggy module checks the Ethereum chain for each submitted address, and makes sure that the peggy contract at that address is using the correct source code, and has the correct validator set.
1. There is a rule in the peggy module that the correct peggy eth contract address with the lowest address in the earliest block that is at least n blocks back is the "official" contract. After n blocks passes, this allows anyone wanting to use peggy to know which ethereum contract to send their money to.

# Creating messages for the Ethereum contract

## Valset Process

- Peggy Daemon on each val submits a "EthAddressTx"
- This adds the Eth addresss to the Eth Adddress store, all validation happens here
- A relayer submits a "ValsetRequest" for the valset of the block in which it is accepted.
- This goes into the ValsetRequestStore, along with the valset that was requested.
- When the peggy daemons see a valset in the ValsetRequestStore, they sign over it with their eth keys, and submit a "ValsetConfirmTx". This goes into the "ValsetConfirmStore".
- Once 66% of the peggy daemons have signatures in the ValsetConfirmStore, for a particular valset, a relayer can submit the valset.

## TX Batch process

- User submits Cosmos TX with requested Eth TX "EthTx"
- This goes into everyone's stores by consensus
<!-- - Relayer chooses a TX batch from the tx's in the store
- Relayer submits "BatchReqTx" to Cosmos, it goes into a BatchReq store by consensus after being validated, all Eth Tx's appearing in the requested batch are removed from the mempool. -->
- --> Peggy module sorts TXs into batches, and puts the batches into the "BatchStore", and all Eth TXs in a batch are removed from the mempool.

- Peggy Daemon on each validator sees all batches in the BatchStore, signs over the batches, sends a "BatchConfirmTx" containing all eth signatures for all the batches.
- The BatchConfirmTx goes into a BatchConfirmStore, now the relayer can relay the batch once there's 66%

- Now the batch is processed by the Eth contract.
- The Peggy Daemons on the validators observe this, and they submit a "BatchSubmittedTx". This TX goes into a Batch SubmittedStore. The cosmos state machine discards the batch permanently once it sees that over 66% of the validators have submitted a BatchSubmittedTx.
  - If there are any older batches in the BatchConfirmStore, they are removed because they can now never be submitted. The Eth Tx's in the old batches are released back into the mempool.

# CosmosSDK / Tendermint considerations

## Signing

In order to update the validator state on the Peggy Ethereum contract we need to perform probably the simplest 'off chain work' possible. Signing the current validator state with an Ethereum key and submitting that as a message
so that a relayer may observe and ferry the signed 'ValSetUpdate' message over to Ethereum.

The concept of 'validators signing some known value' is very core to Tendermint and of course any proof of stake system. So when it's presented as a step in any Peggy process no one bats an eye. But we're not coding Peggy as a Tendermint extension, it's a CosmosSDK module.

CosmosSDK modules don't have access to validator private keys or a signing context to work with. In order to get around this we perform the signing in a separate codebase that interacts with the main module. Typically this would be called a 'relayer' but since we're writing a module where the validators specifically must perform actions with their own private keys it may be better to term this as a 'validator external signer' or something along those lines.

The need for an external signer doubles the number of states required to produce a working Peggy CosmosSDK module state machine. For example the ValSetUpdate message generation process requires a trigger message, this goes into a store where the external signers observe it and submit their own signatures. This 'waiting for sigs' state could be eliminated if signing the state update could be processed as part of the trigger message handler itself.

This obviously isn't a show stopper, but if it's easy _and_ maintainable we should consider using ABCI to do this at the Tendermint level.

## Peggy Consensus

Performing our signing at the CosmosSDK level rather than the Tendermint level has other implications. Mainly it changes the nature of the slashing and halting conditions. At the Tendermint level if signing the ValSetUpdate was part of processing the message failing to do so would result in downtime for that validator on that block. On the other hand submitting a ValSetUpdate signature in a CosmosSDK module is just another message, having no consensus impact other than slashing conditions we may add. Since slashing conditions are slow this produces the following potential vulnerability.

If a validator failing to produce ValSetUpdates and the process is implemented in Tendermint they are simply racking up downtime and have no capabilities as a validator. But if the process is implemented at the CosmosSDK level they will continue to operate normally as a validator.

My intuition about vulnerabilities here is that they could only be used to halt the bridge using 1/3rd of the stake. Since that's roughly the same as halting the chain using 1/3rd of the active stake I don't think it's an issue.
Ethereum event feed

- There is a governance parameter calle EthBlockDelay, for example 50 blocks
- Peggy Daemons get the current block number from their Geth, then get the events from the block EthBlockDelay lower than the current one
- They send an EthBlockData message
- These messages go in an EthBlockDataStore, indexed by the block number and the validator that sent them.
- Once there is a version of the block data that matches from at least 66% of the validators, it is considered legit
- Once a block goes over 66% (and all previous blocks are also over 66%), tokens are minted as a result of EthToCosmos transfers in that block.
- (optional downtime slashing???) Something watches for validators that have not submitted matching block data, and triggers a downtime slashing

Alternate experimental ethereum input

- Let's say that accounts have to send a ClaimTokens message to get their tokens that have been transferred over from Eth (just for sake of argument)
- Each validator connects directly to Geth (instead of the peggy daemon handling it)
- When a ClaimTokens message comes in, each validator in the state machine, checks it's own eth block DB (this is seperate from the Cosmos KV stores, and is access only by the peggy module. The peggy daemon fills it up with block data from a different process) to see if the tokens have been transferred to that account.
- Validators with a different opinion on the state of Eth will arrive at different conclusions, and produce different blocks. Validators that disagree will have downtime, according to Tendermint. This allows Tendermint to handle all consensus, and we don't think about it.

Ethereum to Cosmos transfers

- Function in contract takes a destination Cosmos address and transfer amount
- Does the transfer and logs a EthToCosmosTransfer event with the amount and destination

Alternate batch assembly by validators

- Allowing anyone to request batches as we have envisioned previously opens up an attack vector where someone assembles unprofitable batches to stop the bridge.
- Instead, why not just have the validators assemble batches?
- In the cosmos state machine, they look at all transactions, sort them from lowest to highest fee, and chop that list into batches.
- Now relayers can try to submit the batches.
- Batches are submitted to Eth, and transactions build up in the pool, in parallel.
- At some point the validators make new batches sorted in the same way.
- Relayers may choose to submit new batches which invalidate older low fee batches.
