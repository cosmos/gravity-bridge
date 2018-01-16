# Specifcation for 2-way peg between a Tendermint chain and an Ethereum chain

## Contents
* Overview
* Design A
* Design B

## Overview
The goal of the peg zone is to enable the move of assets between a Tendermint
chain and an Ethereum chain. It is designed to allow for secure and cheap
transfers of all Ethereum tokens (Ether and ERC20) as well as all Cosmos
tokens.

There are three major pieces:
* a Cosmos ABCI app
* a set of Ethereum smart contracts
* a relayer process to connect the above two

The below described points are applicable to both designs.

### Cosmos ABCI app
The Cosmos peg zone accepts and sends IBC packets. When the app receives an IBC
packet it processes it and then affects a change on the Ethereum state. When
the app is informed of a state change on Ethereum it generates and sends an IBC
packet.

### Ethereum Smart Contracts
The set of Ethereum smart contracts tracks the consensus (validator set) of the
Cosmos peg zone. It verifies updates to the state of the validator set and is
responsible for handling of the locked funds on the Ethereum side.

### Relayer Process
The relayer process can be someone manually copying information between the two
chains or an automatic process that is run alongside the peg zone. It is
responsible for taking data from the Cosmos peg zone and posting it to the 
Ethereum chain. Furthermore it takes data from Ethereum and posts it as a 
transaction to the Cosmos peg zone.
We do not yet define the specifics of this relayer process or the economic 
incentives to run it. Generally anyone that wants the peg zone to be successful
has an incentive to run the relayer process.
