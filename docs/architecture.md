# Architecture

The following document describe all relevant actors and pieces of software within the Peggy (v2) suite. Generally it is divided between the Cosmos SDK modules, the EVM based smart contracts, and the Golang relayer client operated by validators that watch both chains and relay information according to events which are heard on either chain.

## EVM Contracts, Actors & Data
- **Claim**
    - A data type that represents a Claim that some event happened on an SDK chain. The Claim is submitted as a transaction from one or more SDK chain validators with reference to a “Prophecy” that includes possible events of “Lock” or “Burn”. Claims could be made about other events but “Lock” and “Burn” are the only two relevant to the transfer of fungible tokens.
- **Prophecy**
    - A data type that represents the first Claim made on the contract about an event coming from an SDK chain. This Claim generates a new Prophecy that has a unique ID used as a reference for subsequent Claims made to support the initial Prophecy.
- **Valset**
    - A contract that keeps track of a set of SDK chain validators and their corresponding token weights.
- **Valset Operator**
    - An actor with the permission to update the validator set within the `Valset` contract.
    - This is a point of centralization that is addressed in Peggy v3.
- **CosmosBridge**
    - A contract which receives new Prophecies from validators, concludes Prophecies after a threshold of Claims have been made against them within the `Oracle` contract and issues wrapped tokens on behalf of the `BridgeBank` or unlocks tokens previously in escrow.
- **CosmosBridge Operator**
    - An actor with the ability to update the `CosmosBridge` references to the `BridgeBank` and `Oracle` contracts.
    - This is a point of centralization that is addressed in Peggy v3.
- **Oracle**
    - A contract which receives Claims from validators against specific Prophecies until a threshold is reached whereupon the Prophecy is concluded within the `CosmosBridge`.
- **BridgeBank**
    - A group of contracts comprised of `BridgeBank`, `CosmosBank` and `EthereumBank` which manages locked EVM assets as well as newly minted SDK based assets.
- **BridgeBank Operator**
    - Actor who is able to deploy new token denominations used to represent SDK chain assets as well as make ETH deposits on the contract directly.
    - This is a point of centralization that is addressed in Peggy v3.
- **BridgeToken**
    - Contract template for standard ERC-20 token that is managed by the `BridgeBank` in order to represent SDK chain assets.
- **BridgeRegistry**
    - Contract which records the values of as well as executes an event containing the contract addresses of `Valset`, `CosmosBridge`, `Oracle` and `BridgeBank`.

## SDK Chain Modules and Messages
- **x/oracle**
    - Manages storage of Prophecies, regardless of their content, and manages subsequent Claims and validators contained within Prophecies. 
- **x/ethbridge**
    - References the oracle module in order to process Claims and subsequent successful and unsuccessful Prophecies by minting, burning and locking Coins.
- **MsgLock**
    - A message that denotes a SDK chain asset should be locked and minted on the EVM chain as a `BridgeToken`.
- **MsgBurn**
    - A message that denotes an EVM chain asset which exists on the SDK chain should be burned and subsequently unlocked on the EVM chain.
- **MsgCreateEthBridgeClaim**
    - A message that creates a Claim on behalf of a validator for some event that occurred on the EVM chain. This message creates a new Prophecy if one does not previously exist and adds a Claim to that Prophecy or another that was previously registered.

## Relayer Events
- **Cosmos Event MsgBurn/MsgLock**
    - When the Cosmos Listener hears this event, it handles it by executing the `NewProphecyClaim()` function within the `CosmosBridge` using the current validator key.
- **Ethereum Event LogLock**
    - When `LogLock` event is heard by the relayer it converts it to a `ProphecyClaim`, and relays a tx to Cosmos.
- **Ethereum Event LogNewProphecyClaim**
    - When `LogNewProphecyClaim` is heard by the Relayer it creates subsequent Claims to push the Prophecy to completion.

## User Flows
### Initialization

In order for the EVM side of Peggy to exist it must first be deployed. There are a total of 5 required contract deployments: `Valset` , `CosmosBridge` , `Oracle`, `BridgeBank` and `BridgeRegistry`. They should be deployed in that order as subsequent deployments make reference to already deployed contracts. The `Valset`, `CosmosBridge` and `BridgeBank` record the address of the account that deployed the contracts as `operator` of each contract. These can be separate addresses or the same.

When moving ERC-20 Tokens on an EVM chain to the Peggy equipped SDK chain it should be noted that prior to the transfer an `approve()` must be called on the ERC-20 to allow the `BridgeBank` contract to move tokens on behalf of the user.

When moving Coins from an SDK chain to an EVM chain, the operator of the `BridgeToken`  must first have initiated a new `BridgeToken` by executing the `createNewBridgeToken()` function.

Environment variables need to be set for both `test-contracts` and the root directory within the corresponding `.env` files. These contain private keys and mnemonic phrases for deploying contracts and operating validator relayers.

### EVM → SDK

After proper initialization has been made the first step to moving an asset from the EVM chain to the SDK chain is to execute the `lock()` function on the `BridgeBank` contract, referencing the cosmos address of the recipient, the amount to transfer and the token address within the EVM setting. If the asset is Ether or whatever native asset exists on this EVM chain, the token address should be the address type equivalent of the number 0. The asset should be transferred to the `BridgeBank` address and the event `lockFunds` should be emitted.

The relayer should be running with settings that include the `chain-id` of the EVM chain and the `BridgeBank` address in order for the relayer to set up a websocket and listen for events. Once `lockFunds` is heard it will convert the event into a `ProphecyClaim` and submit it as a cosmos `MsgCreateEthBridgeClaim` within the `ethbridge` module.

Upon receiving `MsgCreateEthBridgeClaim` the `ethbridge` module will store the Claim information within the `oracle` module keeper where a tally is kept about all references to this Claim. As more validators who were running relayers submit more Claims the Prophecy will reach a threshold of acceptable weight and execute the `Mint` module to create the new tokens on behalf of the designated recipient.

### EVM→ SDK → EVM

Should the user desire to move their EVM based SDK asset back to the EVM setting they will need to sign and execute `MsgBurn` within the SDK chain referencing the EVM based recipient and amount to burn. This Coin amount will be burned within the SDK and an event will be emitted called `MsgBurn` . When the relayer hears this message it will convert the information into a ProphecyClaim and submit it on the `CosmosBridge`. Upon receiving the Claim on the Cosmos Bridge the total weight of all Claims is stored within the `Oracle` contract until a desired threshold is reached and the `BridgeBank` is tasked with unlocking the original tokens and sending to the designated recipient.

### SDK → EVM

In order to move an SDK based Coin to the EVM chain a user would need to sign and execute a `MsgLock` on the `ethbridge` module. This module will move the coins into escrow controlled by that module and event an event called `MsgLock`.

The relayer that is configured to listen for events on the SDK chain will hear the `MsgLock` event and create a `ProphecyClaim` containing the Lock event destined for the EVM chain. The event would contain the token address of a `BridgeToken` which was created with the same `symbol` as the Coin `denom`. Without a corresponding `BridgeToken` the `ProphecyClaim` will eventually fail.


If a `BridgeToken` exists with the correct `denom` / `symbol` and a threshold of validators relayers have made `ProphecyClaim`s then `processBridgeProphecy` can be executed on the `Oracle` contract, which in turn executes `completeProphecy` on the `CosmosBridge` contract. The `completeProphecy` may also be executed automatically within the `Oracle` contract upon a threshold complete `newOracleClaim` execution.

When `completeProphecy` is executed within `CosmosBridge`, the Claim type is deciphered as either a Lock or a Burn. Since this would be a Lock event, the corresponding denomination of  `BridgeToken` within the `BridgeBank` is requested to `mintBridgeTokens()` of the corresponding amount and to the corresponding recipient.

### SDK→EVM→SDK

Should an SDK native asset on the EVM chain need to be returned to the SDK chain, the asset is locked up similarly to the EVM native asset moving to the SDK. In the same way the relayer listens for this lock event and submits it on the SDK chain as an `MsgCreateEthBridgeClaim`. Here it is able to distinguish the denom and will not mint a new Coin but rather move the escrowed coin into the recipients account.
