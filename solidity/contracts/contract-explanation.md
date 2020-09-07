# Functioning of the Peggy.sol contract

The Peggy contract locks assets on Ethereum to facilitate a Tendermint blockchain creating synthetic versions of those assets. It is designed to be used alongside software on the Tendermint blockchain, but this article focuses on the Ethereum side.

Usage example:

- You send 25 DAI to the Peggy contract, specifying which address on the Tendermint chain should recieve the syntehtic DAI.
- Validators on the Tendermint chain see that this has happened and mint 25 synthetic DAI for the address you specified on the Tendermint chain.
- You send the 25 synthetic DAI to Jim on the Tendermint chain.
- Jim sends the synthetic DAI to Peggy module on the Tendermint chain, specifying which Ethereum address should receive it.
- The Tendermint validators burn the synthetic DAI on the Tendermint chain and unlock 25 DAI for Jim on Ethereum

## Security model

The Peggy contract is basically a multisig with a few tweaks. Even though it is designed to be used with a consensus process on Tendermint, the Peggy contract itself encodes nothing about this consensus process. There are three main operations- UpdateValset, SubmitBatch, and TransferOut. UpdateValset updates the signers on the multisig, and their relative powers. This mirrors the validator set on the Tendermint chain, so that all the Tendermint validators are signers, in proportion to their staking power on the Tendermint chain. An UpdateValset transaction must be signed by 2/3's of the current valset to be accepted. SubmitBatch is used to submit a batch of transactions unlocking and transferring tokens to Ethereum addresses. It is used to send tokens from Cosmos to Ethereum. The batch must be signed by 2/3's of the current valset. TransferOut is used to send tokens onto the Tendermint chain. It simply locks the tokens in the contract and emits an event which is picked up by the Tendermint validators.

### UpdateValset

A valset consists of a list of validator's Ethereum addresses, their voting power, and a nonce for the entire valset. UpdateValset takes a new valset, the current valset, and the signatures of the current valset over the new valset. The valsets and the signatures are currently broken into separate arrays because it is not possible to pass arrays of structs into Solidity external functions. Because of this, UpdateValset first does a few checks to make sure that all the arrays that make up a valset are the same length.

Then, it checks the supplied current valset against the saved checkpoint. This requires some explanation. Because valsets contain over 100 validators, storing these all on the Ethereum blockchain each time would be quite expensive. Because of this, we only store a hash of the current valset, then let the caller supply the actual addresses, powers, and nonce of the valset. We call this hash the checkpoint. This is done with the function makeCheckpoint.

Once we are sure that the valset supplied by the caller is the correct one, we check that the new valset nonce is higher than current valset nonce (NOTE: It is still not entirely clear to us whether this rule should instead check that the nonce is incremented by exactly one). This ensures that old valsets cannot be submitted because their nonce is too low.

Now, we make a checkpoint from the submitted new valset, using makeCheckpoint again. In addition to be used as a checkpoint later on, we first use it as a digest to check the current valset's signature over the new valset. We use checkValidatorSignatures to do this.

CheckValidatorSignatures takes a valset, an array of signatures, a hash, and a power threshold. It checks that the powers of all the validators that have signed the hash add up to the threshold. This is how we know that the new valset has been approved by at least 2/3s of the current valset. We iterate over the current valset and the array of signatures, which should be the same length. For each validator, we first check if the signature is all zeros. This signifies that it was not possible to obtain the signature of a given validator. If this is the case, we just skip to the next validator in the list. Since we only need 2/3s of the signatures, it is not required that every validator sign every time, and skipping them stops any validator from being able to stop the bridge.

If we have a signature for a validator, we verify it, throwing an error if there is something wrong. We also increment a cumulativePower counter with the validator's power. Once this is over the threshold, we break out of the loop, and the signatures have been verified! If the loop ends without the threshold being met, we throw an error. Because of the way we break out of the loop once the threshold has been met, if the valset is sorted by descending power, we can usually skip evaluating the majority of signatures. To take advantage of this gas savings, it is important that valsets be produced by the validators in descending order of power.

At this point, all of the checks are complete, and it's time to update the valset! This is a bit anticlimactic, since all we do is save the new checkpoint over the old one. An event is also emitted.

### SubmitBatch

### TransferOut
