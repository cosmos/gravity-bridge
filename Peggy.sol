pragma solidity ^0.6.4;
import "./SafeMath.sol";

contract Peggy {
    using SafeMath for uint256;
    bytes32 public checkpoint;

    event LogValsetUpdated(address[] _validators, uint256[] _powers);
    event LogCheckpoint(bytes32 _hash);
    event LogIndex(uint256 _val);
    event LogAddress(address _val);

    constructor(bytes32 _checkpoint) public {
        checkpoint = _checkpoint;
    }

    function updateValset(
        // The new version of the validator set
        address[] memory _newValidators,
        uint256[] memory _newPowers,
        // The old validators that approve the change
        address[] memory _oldValidators,
        uint256[] memory _oldPowers,
        // These are arrays of the parts of the oldValidators signatures
        bytes32[] memory _r,
        uint8[] memory _v,
        bytes32[] memory _s
    ) public {
        // Check that oldValidators is correct:
        // Loop accumulates a hash of oldValidators and oldPowers
        // Checks it against the checkpoint
        bytes32 oldValidatorsHash = 0;
        for (uint256 i = 0; i < _oldValidators.length; i = i.add(1)) {
            oldValidatorsHash = keccak256(
                abi.encodePacked(
                    oldValidatorsHash,
                    _oldValidators[i],
                    _oldPowers[i]
                )
            );
        }

        //emit LogCheckpoint(oldValidatorsHash);
        require(oldValidatorsHash == checkpoint);

        // Generate hash of newValidators:
        // Loop accumulates a hash of newValidators and newPowers
        // assigns to var newValidatorsHash
        bytes32 newValidatorsHash = 0;
        for (uint256 j = 0; j < _newValidators.length; j = j.add(1)) {
            newValidatorsHash = keccak256(
                abi.encodePacked(
                    newValidatorsHash,
                    _newValidators[j],
                    _newPowers[j]
                )
            );
        }
        emit LogCheckpoint(newValidatorsHash);

        // Check that oldValidators signed off:
        // Loop checks _signatures against _oldValidators and newValidatorsHash and
        // sums _oldPowers until it reaches 2/3rds
        uint256 sumPower = 0;
        for (uint256 k = 0; k < _oldValidators.length; k = k.add(1)) {
            // Validate signature of each old validator over the newValidatorsHash
            address gotsig = ecrecover(newValidatorsHash, _v[k], _r[k], _s[k]);
            emit LogAddress(gotsig);

            require(
                _oldValidators[k] == gotsig
            );

            // Sum up cumulative power of all oldValidators that have signed
            sumPower = sumPower + _oldPowers[k];

            // If the cumulative power is greater than 66.666% of total
            // (we are arbitrarily choosing the maximum power to be 100,000, this should be enough for a PoC)
            if (sumPower > 66666) {
                break;
            }
        }

        checkpoint = newValidatorsHash;

        emit LogValsetUpdated(_newValidators, _newPowers);
    }
}
