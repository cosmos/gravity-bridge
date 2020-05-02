pragma solidity ^0.6.4;
import "./SafeMath.sol";


// TODO gas optimization: break loops early
// TODO gas optimization: combine operations to avoid iterating over the same set
// multiple times
contract Peggy {
    using SafeMath for uint256;
    bytes32 public checkpoint;
    uint256 public txNonce;
    // TODO do we need a separate nonce??????? <-- We think this is purely a gas optimization

    event LogValsetUpdated(address[] _validators, uint256[] _powers);
    event LogCheckpoint(bytes32 _hash);
    event LogIndex(uint256 _val);
    event LogAddress(address _val);

    // constructor(bytes32 _checkpoint) public {
    //     checkpoint = _checkpoint;
    // }

    // This function checks that the caller supplied data is consistent with the checkpoint
    // This lets the contract "store" data offchain.
    function checkCheckpoint(
        address[] memory _validators,
        uint256[] memory _powers,
        bytes32 _checkpoint
    ) private pure {
        // Loop accumulates a hash of oldValidators and oldPowers
        bytes32 validatorHash = 0;
        for (uint256 i = 0; i < _validators.length; i = i.add(1)) {
            validatorHash = keccak256(
                abi.encodePacked(validatorHash, _validators[i], _powers[i])
            );
        }

        // Checks it against the checkpoint
        require(
            validatorHash == _checkpoint,
            "Supplied validators and powers do not match checkpoint."
        );
    }

    // TODO: We need to make it so that you can't submit newValidators in a
    // non-descending order of power
    // (if this was possible, you could screw over the next guy)
    function hashValidators(
        address[] memory _validators,
        uint256[] memory _powers
    ) public pure returns (bytes32) {
        uint256 valsetLength = _validators.length;
        require(valsetLength == _powers.length, "Malformed validator set");

        // Generate hash of validators:
        // Loop accumulates a hash of validators and powers
        // assigns to validatorsHash
        bytes32 validatorsHash = 0;
        for (uint256 i = 0; i < valsetLength; i = i.add(1)) {
            validatorsHash = keccak256(
                abi.encodePacked(validatorsHash, _validators[i], _powers[i])
            );
        }

        return validatorsHash;
    }

    // Make sure that validator powers are equal or decreasing. This prevents someone
    // forcing the next caller to waste gas iterating all the validators
    function checkValidatorPowerOrder(uint256[] memory _powers) public pure {
        for (uint256 i = 0; i < _powers.length; i = i.add(1)) {
            if (i != 0) {
                require(
                    _powers[i] <= _powers[i - 1],
                    "Validator power must not be higher than previous validator in batch"
                );
            }
        }
    }

    function hashTransactions(
        uint256[] memory _amounts,
        address[] memory _destinations,
        uint256[] memory _fees,
        uint256[] memory _nonces
    ) public pure returns (bytes32) {
        // Check that all components of batch have same length
        uint256 batchLength = _amounts.length;
        require(
            batchLength == _destinations.length &&
                batchLength == _fees.length &&
                batchLength == _nonces.length,
            "Malformed batch of transactions"
        );

        // Loop accumulates a hash of _amounts, _destinations, _fees, and _nonces
        bytes32 batchHash = 0;
        for (uint256 i = 0; i < batchLength; i = i.add(1)) {
            batchHash = keccak256(
                abi.encodePacked(
                    batchHash,
                    _amounts[i],
                    _destinations[i],
                    _fees[i],
                    _nonces[i]
                )
            );
        }

        return batchHash;
    }

    function checkTxNonces(uint256[] memory _nonces) public view {
        uint256 lastNonce = txNonce;
        for (uint256 i = 0; i < _nonces.length; i = i.add(1)) {
            require(
                _nonces[i] > lastNonce,
                "Transaction nonces in batch must be strictly increasing"
            );
            lastNonce = _nonces[i];
        }
    }

    // This checks that enough old validators have signed the given hash
    // Will error if any signature is incorrect
    function checkSignatures(
        // These are the validators and their powers
        address[] memory _validators,
        // These are arrays of the parts of the validators signatures
        uint8[] memory _v,
        bytes32[] memory _r,
        bytes32[] memory _s,
        // This is the hash that we are checking signatures over
        bytes32 _hash
    ) public pure {
        // Loop checks signatures (v, r, s) against _validators and hash
        for (uint256 k = 0; k < _validators.length; k = k.add(1)) {
            // Validate signature of each old validator over the newValidatorsHash
            require(
                _validators[k] == ecrecover(_hash, _v[k], _r[k], _s[k]),
                "Old validator signature does not match."
            );
        }
    }

    // Checks if submitted powers are enough to approve this action
    function checkPowers(uint256[] memory _powers) public pure {
        uint256 sumPower = 0;
        for (uint256 k = 0; k < _powers.length; k = k.add(1)) {
            // Sum up cumulative power
            sumPower = sumPower + _powers[k];
        }

        // If the cumulative power is greater than 66.666% of total
        // (we are arbitrarily choosing the maximum power to be 100,000, this should be enough for a PoC)
        // TODO: make percentage configurable
        require(
            sumPower > 66666,
            "Submitted validator set does not have enough power."
        );
    }

    function updateValset(
        // The new version of the validator set
        address[] memory _newValidators,
        uint256[] memory _newPowers,
        // The old validators that approve the change
        address[] memory _oldValidators,
        uint256[] memory _oldPowers,
        // These are arrays of the parts of the oldValidators signatures
        uint8[] memory _v,
        bytes32[] memory _r,
        bytes32[] memory _s
    ) public {
        // Get hash of submitted new validator set
        bytes32 newValidatorsHash = hashValidators(_newValidators, _newPowers);

        // Check that old validator set matches the checkpoint
        checkCheckpoint(_oldValidators, _oldPowers, checkpoint);

        // Check that enough old validators have signed the new validator set
        checkSignatures(_oldValidators, _v, _r, _s, newValidatorsHash);

        // Save the new validator set
        checkpoint = newValidatorsHash;
        emit LogValsetUpdated(_newValidators, _newPowers);
    }

    function submitBatch(
        // The validators that approve the batch
        address[] memory _validators,
        uint256[] memory _powers,
        // These are arrays of the parts of the validators signatures
        uint8[] memory _v,
        bytes32[] memory _r,
        bytes32[] memory _s,
        // The batch of transactions
        uint256[] memory _amount,
        address[] memory _destination,
        uint256[] memory _fee,
        uint256[] memory _nonce // TODO: multi-erc20 support (input contract address). // Will be done once we get the basic version working
    ) public {
        // Iterate over batch and make sure than nonce is strictly increasing (can have gaps)
        // And that it is higher than the stored
        // As long as these are fulfilled, it's valid and you can do the stuff in the txs
    }

    constructor(
        // The token that this bridge bridges
        address _tokenContract,
        // The validator set
        address[] memory _validators,
        uint256[] memory _powers,
        // These are arrays of the parts of the validators signatures
        bytes32[] memory _r,
        uint8[] memory _v,
        bytes32[] memory _s
    ) public {}
}
