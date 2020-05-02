pragma solidity ^0.6.4;
import "./SafeMath.sol";
import "./IERC20.sol";


// TODO gas optimization: break loops early
// TODO gas optimization: combine operations to avoid iterating over the same set
// multiple times
contract Peggy {
    using SafeMath for uint256;
    bytes32 public storedCheckpoint;
    uint256 public storedNonce;
    // TODO do we need a separate nonce or can it be stored in the checkpointed data?
    address public tokenContract;

    uint256 constant ENOUGH_POWER = 66666;

    event LogValsetUpdated(address[] _validators, uint256[] _powers);

    // - Check that the supplied current validator set matches the saved checkpoint
    function checkCheckpoint(
        address[] memory _currentValidators,
        uint256[] memory _currentPowers
    ) public view {
        bytes32 currentValidatorsHash = 0;
        for (uint256 i = 0; i < _currentValidators.length; i = i.add(1)) {
            currentValidatorsHash = keccak256(
                abi.encodePacked(
                    currentValidatorsHash,
                    _currentValidators[i],
                    _currentPowers[i]
                )
            );
        }

        require(
            currentValidatorsHash == storedCheckpoint,
            "Supplied validators and powers do not match checkpoint."
        );
    }

    function checkValidatorSignatures(
        // The current validator set and their powers
        address[] memory _currentValidators,
        uint256[] memory _currentPowers,
        // The current validator's signatures
        uint8[] memory _v,
        bytes32[] memory _r,
        bytes32[] memory _s,
        // This is what we are checking they have signed
        bytes32 theHash
    ) public pure {
        uint256 cumulativePower = 0;

        for (uint256 k = 0; k < _currentValidators.length; k = k.add(1)) {
            // Check that the current validator has signed off on the hash
            require(
                _currentValidators[k] ==
                    ecrecover(theHash, _v[k], _r[k], _s[k]),
                "Current validator signature does not match."
            );

            // Sum up cumulative power
            cumulativePower = cumulativePower + _currentPowers[k];

            // Break early to avoid wasting gas
            if (cumulativePower > ENOUGH_POWER) {
                break;
            }
        }

        // Check that there was enough power
        require(
            cumulativePower > ENOUGH_POWER,
            "Submitted validator set does not have enough power."
        );
    }

    function updateValset(
        // The new version of the validator set
        address[] memory _newValidators,
        uint256[] memory _newPowers,
        // The current validators that approve the change
        address[] memory _currentValidators,
        uint256[] memory _currentPowers,
        // These are arrays of the parts of the current validator's signatures
        uint8[] memory _v,
        bytes32[] memory _r,
        bytes32[] memory _s
    ) public {
        // CHECKS

        // Check that new validators and powers set is well-formed
        require(
            _newValidators.length == _newPowers.length,
            "Malformed new validator set"
        );

        // Check that current validators, powers, and signatures (v,r,s) set is well-formed
        require(
            _currentValidators.length == _currentPowers.length &&
                _currentValidators.length == _v.length &&
                _currentValidators.length == _r.length &&
                _currentValidators.length == _s.length,
            "Malformed current validator set"
        );

        // - Check that the supplied current validator set matches the saved checkpoint
        checkCheckpoint(_currentValidators, _currentPowers);

        // - Get hash of new validator set
        // - Check that validator powers are decreasing or equal (this prevents the
        // next caller from wasting gas)
        bytes32 newValidatorsHash = 0;
        {
            for (uint256 i = 0; i < _newValidators.length; i = i.add(1)) {
                if (i != 0) {
                    require(
                        !(_newPowers[i] > _newPowers[i - 1]),
                        "Validator power must not be higher than previous validator in batch"
                    );
                }
                newValidatorsHash = keccak256(
                    abi.encodePacked(
                        newValidatorsHash,
                        _newValidators[i],
                        _newPowers[i]
                    )
                );
            }
        }

        // - Check that enough current validators have signed off on the new validator
        // set hash
        checkValidatorSignatures(
            _currentValidators,
            _currentPowers,
            _v,
            _r,
            _s,
            newValidatorsHash
        );

        // ACTIONS

        storedCheckpoint = newValidatorsHash;

        // LOGS

        emit LogValsetUpdated(_newValidators, _newPowers);
    }

    function submitBatch(
        // The validators that approve the batch
        address[] memory _currentValidators,
        uint256[] memory _currentPowers,
        // These are arrays of the parts of the validators signatures
        uint8[] memory _v,
        bytes32[] memory _r,
        bytes32[] memory _s,
        // The batch of transactions
        uint256[] memory _amounts,
        address[] memory _destinations,
        uint256[] memory _fees,
        uint256[] memory _nonces // TODO: multi-erc20 support (input contract address). // Will be done once we get the basic version working
    ) public {
        // CHECKS

        // - Check that current validators, powers, and signatures (v,r,s) set is well-formed
        require(
            _currentValidators.length == _currentPowers.length &&
                _currentValidators.length == _v.length &&
                _currentValidators.length == _r.length &&
                _currentValidators.length == _s.length,
            "Malformed current validator set"
        );

        // - Check that the transaction batch is well-formed
        require(
            _amounts.length == _destinations.length &&
                _amounts.length == _fees.length &&
                _amounts.length == _nonces.length,
            "Malformed batch of transactions"
        );

        // - Check that the supplied current validator set matches the saved checkpoint
        checkCheckpoint(_currentValidators, _currentPowers);

        // - Get hash of the transaction batch
        // - Check that the tx nonces are higher than the stored nonce and are
        // strictly increasing (can have gaps)
        bytes32 transactionsHash = 0; // TODO: figure out if this is the best way to initialize the hash
        uint256 lastNonce = storedNonce;
        {
            for (uint256 i = 0; i < _amounts.length; i = i.add(1)) {
                require(
                    _nonces[i] > lastNonce,
                    "Transaction nonces in batch must be strictly increasing"
                );
                lastNonce = _nonces[i];

                transactionsHash = keccak256(
                    abi.encodePacked(
                        transactionsHash,
                        _amounts[i],
                        _destinations[i],
                        _fees[i],
                        _nonces[i]
                    )
                );
            }
        }

        // - Check that enough current validators have signed off on the transaction batch
        checkValidatorSignatures(
            _currentValidators,
            _currentPowers,
            _v,
            _r,
            _s,
            transactionsHash
        );

        // ACTIONS

        // Store nonce
        storedNonce = lastNonce;

        // - Send transaction amounts to destinations
        // - Send transaction fees to msg.sender
        {
            for (uint256 i = 0; i < _amounts.length; i = i.add(1)) {
                IERC20(tokenContract).transfer(_destinations[i], _amounts[i]);
                IERC20(tokenContract).transfer(msg.sender, _fees[i]);
            }
        }
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
