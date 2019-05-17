# Unidirectional Peggy: Ethereum -> Cosmos

Unidirectional Peggy does not currently include functionality related to Cosmos -> Ethereum transfers, instead focusing on core features for unidirectional transfers. Building towards the eventual goal of bidirectional transfers between Ethereum and Cosmos, Unidirectional Peggy exclusively includes functionality to safely lock and unlock Ethereum and ERC20 tokens, providing Cosmos based Relayer and Oracle modules with the information required for unidirectional operatinos. Validators witness lock events and submit proof in the form of signed hashes to the Cosmos based modules, which are responsible for aggregating and tallying the validator's signatures and their respective signing power. The system is managed by the contract's deployer, known as the relayer, and follows a straightforward process:
1. Users lock Ethereum or ERC20 tokens on the contract, resulting in the emission of an event containing the created item's original sender's Ethereum address, the intended recipient's Cosmos address, the type of token, the amount locked, and the item's unique nonce.
2. Validators witness these lock events and sign a hash containing the unique item's information, which is submitted to a Cosmos Relayer module and communicated to the Oracle.
3. Once the Oracle module has verified that the validators' aggregated signing power is greater than the specified threshold, it mints the appropriate amount of tokens and forwards them to the intended recipient.

These contracts are for testing purposes only and are NOT intended for production. In order to prevent any loss of user funds, Ethereum and/or tokens locked in item can be withdrawn directly by the original sender at any time. Once the system's components are operational, these features will be removed (and others added) so that Unidirectional Peggy is production capable.

## Installation
Install Truffle: `$ npm install -g truffle`

Install dependencies: `$ npm install`


This project currently uses solc@0.5.0, make sure that this version of the Solidity compiler is being used to compile the contracts and does not conflict with other verions that may be installed on your machine.

## Testing
Run commands from the appropriate directory: `$ cd ethereum-contracts`

Start the truffle environment: `$ truffle develop`

In another tab, run tests: `$ truffle test`

Run individual tests: `$ truffle test test/<test_name.js>`


## Future Work
The related Cosmos modules are under active development. Once Ethereum -> Cosmos transfers have been successfully prototyped, Peggy functionality for bidirectional transfers (such as validator sets, signature validation, and secured token unlocking procedures) will be integrated into the contracts. Previous work in these areas is a valuable resource that will be leveraged once the complete system is ready for bidirectional transfers.
