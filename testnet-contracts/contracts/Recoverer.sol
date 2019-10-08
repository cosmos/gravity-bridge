pragma solidity ^0.5.0;

// TODO: Refactor Recoverer to ECDSA based signature validation
contract Recoverer {
    bytes constant internal PREFIX = "\x19Ethereum Signed Message:\n32";

	function isValidSignature(
		address signer,
		bytes32 hash,
		uint8 v,
        bytes32 r,
        bytes32 s
	)
		public
		pure
		returns (bool valid)
	{
		bytes32 prefixedHash = keccak256(abi.encodePacked(PREFIX, hash));
		return ecrecover(prefixedHash, v, r, s) == signer;
	}
}
