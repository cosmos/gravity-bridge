# Arbitrary logic functionality

Gravity includes the functionality to make arbitrary calls out to other Ethereum contracts. This can be used to allow the Cosmos chain to take actions on Ethereum. This functionality is very general. It can even be used to implement the core token transferring functionality of the bridge. However, there is one important caveat: these arbitrary logic contracts can transact with ERC20 tokens, but not any other kind of asset, such as ERC721. Interacting with non-ERC20 assets would require modifications to the core Gravity contract.

# Architecture

`CreateContractCallTx`

Gravity offers a method which can be called by other modules to create an outgoing logic call. To use this method, a calling module must first assemble a logic call (more on this later). This is then submitted to the Gravity module with `CreateContractCallTx`. From here, it is signed by the validators. Once it has enough signatures, a Gravity relayer will pick it up and submit it to the Gravity contract on Ethereum.

`ContractCall`

`CreateContractCallTx` takes an `invalidationNonce`, `invalidationScope`, `payload`, `tokens`, `fees` 

Here is an explanation of its parameters:

- Tokens: These are tokens that are sent to the logic contract before it is executed. The contract can then take actions using the tokens. For example, Gravity could send the logic contract some Uniswap LP tokens that it would then use to redeem liquidity from Uniswap.
- Fees: These are tokens that will be paid by the core Gravity.sol contract to the Gravity relayer for executing the logic call. Fees are paid after the logic contract executes, so it is possible to pay the relayer with tokens that logic contract receives after executing, and then sends back to the core Gravity contract.
- Payload: This is the Ethereum abi encoded function call that will be executed on the logic contract. If you are using a batching middleware contract, then this abi encoded function call will itself contain an array of abi encoded function calls on the actual logic contract.
- Timeout: The logic call will not execute if the block timestamp on Ethereum is higher than the value of this timeout. 
- InvalidationScope and InvalidationNonce: More on these below:


## Invalidation

`invalidation_scope` and `invalidation_nonce` are used as replay protection in the Gravity arbitrary logic call functionality.

When a submitLogicCall transaction is submitted to the Ethereum contract, the contract checks uses `invalidation_scope` to access a key in the invalidation mapping. The value at this key is checked against the supplied `invalidation_nonce`. The logic call is only allowed to go through if the supplied `invalidation_nonce` is higher.

This can be used to implement different invalidation schemes:

### Easiest: timeout-only invalidation
If you don't know what this all means, when you send a logic call to the Gravity module from the Cosmos side, just set the `invalidation_id` to an incrementing integer that you keep track of in your module. Set the `invalidation_nonce` to zero each time. This will create a new entry in the invalidation mapping on Ethereum for each logic batch, providing replay protection, while allowing batches to be completely independent.

### Sequential invalidation
If you don't want it to be possible to submit an early logic call after a later logic call, you can instead set the `invalidation_id` to zero each time, and use an incrementing integer for the `invalidation_nonce`. This makes it so that any logic call that is successfully submitted will invalidate all previous logic calls.

### For example: Token based invalidation
In Gravity's core submitBatch functionality, we have batches of transactions for a given token invalidate earlier batches of that token, but not earlier batches of other tokens. To implement this on top of the submitLogicCall method, we would set the `invalidation_scope` to the token address and keep an incrementing nonce for each token.
