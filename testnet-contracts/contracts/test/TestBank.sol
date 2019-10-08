pragma solidity ^0.5.0;

import "../Bank.sol";

contract TestBank is Bank {

    function() external payable {}

    //Wrapper function to test internal method
    function callDeliver(
        address _token,
        string memory _symbol,
        uint256 _amount,
        address _beneficiary
    )
        public
    {
        return deliver(_token, _symbol, _amount, _beneficiary);
    }

    //Wrapper function to test internal method
    function callDeployBankToken(
        string memory _symbol
    )
        public
        returns(address)
    {
        return deployBankToken(_symbol);
    }
}
