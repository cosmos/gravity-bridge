pragma solidity ^0.5.0;

import "../CosmosBridge.sol";

contract TestCosmosBridge is CosmosBridge {

    constructor()
        public
        CosmosBridge()
    {
        // Intentionally left blank
    }

    //Wrapper function to test internal method
    function callNewCosmosBridgeClaim(
        uint256 _nonce,
        bytes memory _cosmosSender,
        address payable _ethereumReceiver,
        address _tokenAddress,
        string memory _symbol,
        uint256 _amount
    )
        public
        returns(bool)
    {
        return newCosmosBridgeClaim(
            _nonce,
            _cosmosSender,
            _ethereumReceiver,
            _tokenAddress,
            _symbol,
            _amount
        );
    }
}
