pragma solidity ^0.6.6;
import "@openzeppelin/contracts/math/SafeMath.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@nomiclabs/buidler/console.sol";

contract Peggy {
	using SafeMath for uint256;

	// These are updated often
	bytes32 public state_lastCheckpoint;
	uint256 public state_lastTxNonce = 0;

	// These are set once at initialization
	address public state_tokenContract;
	bytes32 public state_peggyId;
	uint256 public state_powerThreshold;

	event ValsetUpdatedEvent(address[] _validators, uint256[] _powers);
	event TransferOutEvent(bytes32 _destination, uint256 _amount);

	// TEST FIXTURES
	// These are here to make it easier to measure gas usage. They should be removed before production
	function testMakeCheckpoint(
		address[] memory _validators,
		uint256[] memory _powers,
		uint256 _valsetNonce,
		bytes32 _peggyId
	) public {
		makeCheckpoint(_validators, _powers, _valsetNonce, _peggyId);
	}

	function testCheckValidatorSignatures(
		address[] memory _currentValidators,
		uint256[] memory _currentPowers,
		uint8[] memory _v,
		bytes32[] memory _r,
		bytes32[] memory _s,
		bytes32 _theHash,
		uint256 _powerThreshold
	) public {
		checkValidatorSignatures(
			_currentValidators,
			_currentPowers,
			_v,
			_r,
			_s,
			_theHash,
			_powerThreshold
		);
	}

	// END TEST FIXTURES

	// Utility function to verify geth style signatures
	function verifySig(
		address _signer,
		bytes32 _theHash,
		uint8 _v,
		bytes32 _r,
		bytes32 _s
	) private pure returns (bool) {
		bytes32 messageDigest = keccak256(
			abi.encodePacked("\x19Ethereum Signed Message:\n32", _theHash)
		);
		return _signer == ecrecover(messageDigest, _v, _r, _s);
	}

	// Make a new checkpoint from the supplied validator set
	// A checkpoint is a hash of all relevant information about the valset. This is stored by the contract,
	// instead of storing the information directly. This saves on storage and gas.
	// The format of the checkpoint is:
	// h(peggyId, "checkpoint", valsetNonce, validators[], powers[])
	// Where h is the keccak256 hash function.
	// The validator powers must be decreasing or equal. This is important for checking the signatures on the
	// next valset, since it allows the caller to stop verifying signatures once a quorum of signatures have been verified.
	function makeCheckpoint(
		address[] memory _validators,
		uint256[] memory _powers,
		uint256 _valsetNonce,
		bytes32 _peggyId
	) public pure returns (bytes32) {
		// bytes32 encoding of the string "checkpoint"
		bytes32 methodName = 0x636865636b706f696e7400000000000000000000000000000000000000000000;

		bytes32 checkpoint = keccak256(
			abi.encode(_peggyId, methodName, _valsetNonce, _validators, _powers)
		);

		return checkpoint;
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
		bytes32 _theHash,
		uint256 _powerThreshold
	) public pure {
		uint256 cumulativePower = 0;

		for (uint256 k = 0; k < _currentValidators.length; k = k.add(1)) {
			// Check that the current validator has signed off on the hash
			require(
				verifySig(_currentValidators[k], _theHash, _v[k], _r[k], _s[k]),
				"Validator signature does not match."
			);

			// Sum up cumulative power
			cumulativePower = cumulativePower + _currentPowers[k];

			// Break early to avoid wasting gas
			if (cumulativePower > _powerThreshold) {
				break;
			}
		}

		// Check that there was enough power
		require(
			cumulativePower > _powerThreshold,
			"Submitted validator set does not have enough power."
		);
	}

	// This updates the valset by checking that the validators in the current valset have signed off on the
	// new valset. The signatures supplied are the signatures of the current valset over the checkpoint hash
	// generated from the new valset.
	function updateValset(
		// The new version of the validator set
		address[] memory _newValidators,
		uint256[] memory _newPowers,
		uint256 _newValsetNonce,
		// The current validators that approve the change
		address[] memory _currentValidators,
		uint256[] memory _currentPowers,
		uint256 _currentValsetNonce,
		// These are arrays of the parts of the current validator's signatures
		uint8[] memory _v,
		bytes32[] memory _r,
		bytes32[] memory _s
	) public {
		// CHECKS

		// Check that new validators and powers set is well-formed
		require(_newValidators.length == _newPowers.length, "Malformed new validator set");

		// Check that current validators, powers, and signatures (v,r,s) set is well-formed
		require(
			_currentValidators.length == _currentPowers.length &&
				_currentValidators.length == _v.length &&
				_currentValidators.length == _r.length &&
				_currentValidators.length == _s.length,
			"Malformed current validator set"
		);

		// Check that the supplied current validator set matches the saved checkpoint
		require(
			makeCheckpoint(
				_currentValidators,
				_currentPowers,
				_currentValsetNonce,
				state_peggyId
			) == state_lastCheckpoint,
			"Supplied current validators and powers do not match checkpoint."
		);

		// Check that the valset nonce is incremented by one
		require(
			_newValsetNonce == _currentValsetNonce.add(1),
			"Valset nonce must be incremented by one"
		);

		// Check that enough current validators have signed off on the new validator set
		bytes32 newCheckpoint = makeCheckpoint(
			_newValidators,
			_newPowers,
			_newValsetNonce,
			state_peggyId
		);

		checkValidatorSignatures(
			_currentValidators,
			_currentPowers,
			_v,
			_r,
			_s,
			newCheckpoint,
			state_powerThreshold
		);

		// ACTIONS

		// Stored to be used next time to validate that the valset
		// supplied by the caller is correct.
		state_lastCheckpoint = newCheckpoint;

		// LOGS

		emit ValsetUpdatedEvent(_newValidators, _newPowers);
	}

	// This function submits a batch of transactions to be executed on Ethereum.
	// The caller must supply the current validator set, along with their signatures over the batch.
	// The contract checks that this validator set matches the saved checkpoint, then verifies their
	// signatures over a hash of the tx batch.
	function submitBatch(
		// The validators that approve the batch
		address[] memory _currentValidators,
		uint256[] memory _currentPowers,
		uint256 _currentValsetNonce,
		// These are arrays of the parts of the validators signatures
		uint8[] memory _v,
		bytes32[] memory _r,
		bytes32[] memory _s,
		// The batch of transactions
		uint256[] memory _amounts,
		address[] memory _destinations,
		uint256[] memory _fees,
		uint256[] memory _nonces // TODO: multi-erc20 support (input contract address).
	) public {
		// CHECKS

		// Check that current validators, powers, and signatures (v,r,s) set is well-formed
		require(
			_currentValidators.length == _currentPowers.length &&
				_currentValidators.length == _v.length &&
				_currentValidators.length == _r.length &&
				_currentValidators.length == _s.length,
			"Malformed current validator set"
		);

		// Check that the transaction batch is well-formed
		require(
			_amounts.length == _destinations.length &&
				_amounts.length == _fees.length &&
				_amounts.length == _nonces.length,
			"Malformed batch of transactions"
		);

		// Check that the supplied current validator set matches the saved checkpoint
		require(
			makeCheckpoint(
				_currentValidators,
				_currentPowers,
				_currentValsetNonce,
				state_peggyId
			) == state_lastCheckpoint,
			"Supplied current validators and powers do not match checkpoint."
		);

		// Check that the tx nonces are higher than the stored nonce and are
		// strictly increasing (can have gaps)
		uint256 lastTxNonceTemp = state_lastTxNonce;
		{
			for (uint256 i = 0; i < _nonces.length; i = i.add(1)) {
				require(
					_nonces[i] > lastTxNonceTemp,
					"Transaction nonces in batch must be strictly increasing"
				);

				lastTxNonceTemp = _nonces[i];
			}
		}

		// bytes32 encoding of "transactionBatch"
		bytes32 methodName = 0x7472616e73616374696f6e426174636800000000000000000000000000000000;
		bytes memory abiEncoded = abi.encode(
			state_peggyId,
			methodName,
			_amounts,
			_destinations,
			_fees,
			_nonces
		);

		// Get hash of the transaction batch
		bytes32 transactionsHash = keccak256(abiEncoded);

		// Check that enough current validators have signed off on the transaction batch
		checkValidatorSignatures(
			_currentValidators,
			_currentPowers,
			_v,
			_r,
			_s,
			transactionsHash,
			state_powerThreshold
		);

		// ACTIONS

		// Store nonce
		state_lastTxNonce = lastTxNonceTemp;

		// Send transaction amounts to destinations
		// Send transaction fees to msg.sender
		uint256 totalFee;
		{
			for (uint256 i = 0; i < _amounts.length; i = i.add(1)) {
				IERC20(state_tokenContract).transfer(_destinations[i], _amounts[i]);
				totalFee = totalFee.add(_fees[i]);
			}
			IERC20(state_tokenContract).transfer(msg.sender, totalFee);
		}
	}

	function transferOut(bytes32 _destination, uint256 _amount) public {
		IERC20(state_tokenContract).transferFrom(msg.sender, address(this), _amount);
		emit TransferOutEvent(_destination, _amount);
	}

	constructor(
		// The token that this bridge bridges
		address _tokenContract,
		// A unique identifier for this peggy instance to use in signatures
		bytes32 _peggyId,
		// How much voting power is needed to approve operations
		uint256 _powerThreshold,
		// The validator set
		address[] memory _validators,
		uint256[] memory _powers,
		// These are arrays of the parts of the validators signatures
		uint8[] memory _v,
		bytes32[] memory _r,
		bytes32[] memory _s
	) public {
		// CHECKS

		// Check that validators, powers, and signatures (v,r,s) set is well-formed
		require(
			_validators.length == _powers.length &&
				_validators.length == _v.length &&
				_validators.length == _r.length &&
				_validators.length == _s.length,
			"Malformed current validator set"
		);

		bytes32 newCheckpoint = makeCheckpoint(_validators, _powers, 0, _peggyId);

		checkValidatorSignatures(
			_validators,
			_powers,
			_v,
			_r,
			_s,
			keccak256(abi.encode(newCheckpoint, _tokenContract, _peggyId, _powerThreshold)),
			_powerThreshold
		);

		// ACTIONS

		state_tokenContract = _tokenContract;
		state_peggyId = _peggyId;
		state_powerThreshold = _powerThreshold;
		state_lastCheckpoint = newCheckpoint;
	}
}
