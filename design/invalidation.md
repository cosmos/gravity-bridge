# Invalidation

`invalidation_id` and `invalidation_nonce` are used as replay protection in the Gravity arbitrary logic call functionality.

When a submitLogicCall transaction is submitted to the Ethereum contract, the contract checks uses `invalidation_id` to access a key in the invalidation mapping. The value at this key is checked against the supplied `invalidation_nonce`. The logic call is only allowed to go through if the supplied `invalidation_nonce` is higher.

This can be used to implement many different invalidation schemes:

## Easiest: timeout-only invalidation
If you don't know what this all means, when you send a logic call to the Gravity module from the Cosmos side, just set the `invalidation_id` to an incrementing integer that you keep track of in your module. Set the `invalidation_nonce` to zero each time. This will create a new entry in the invalidation mapping on Ethereum for each logic batch, providing replay protection, while allowing batches to be completely independent.

## Sequential invalidation
If you don't want it to be possible to submit an early logic call after a later logic call, you can instead set the `invalidation_id` to zero each time, and use an incrementing integer for the `invalidation_nonce`. This makes it so that any logic call that is successfully submitted will invalidate all previous logic calls.

## For example: Token based invalidation
In Gravity's core submitBatch functionality, we have batches of transactions for a given token invalidate earlier batches of that token, but not earlier batches of other tokens. To implement this on top of the submitLogicCall method, we would set the `invalidation_id` to the token address and keep an incrementing nonce for each token.
