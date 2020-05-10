pragma solidity ^0.6.6;

import "@nomiclabs/buidler/console.sol";


contract SigningTest {
	function checkSignature(address _signer, bytes32 _theHash, uint8 _v, bytes32 _r, bytes32 _s)
		public
		pure
	{
		require(_signer == ecrecover(_theHash, _v, _r, _s), "Signature does not match.");
	}
}
