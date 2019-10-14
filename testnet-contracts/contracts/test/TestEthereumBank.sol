pragma solidity ^0.5.0;

import "../EthereumBank.sol";

contract TestEthereumBank is EthereumBank {

    event LogNewEthereumDeposit(
        bytes32 _id
    );

    function() external payable {}

    // //Wrapper function to test internal method
    // function callNewEthereumDeposit(
    //     address payable _sender,
    //     bytes memory _recipient,
    //     address _token,
    //     uint256 _amount
    // )
    //     public
    //     returns(bytes32)
    // {
    //     bytes32 id = newEthereumDeposit(_sender, _recipient, _token, _amount);
    //     emit LogNewEthereumDeposit(id);
    //     return id;
    // }

    // //Wrapper function to test internal method
    // function callUnlockEthereumDeposit(
    //     bytes32 _id
    // )
    //     public
    // {
    //     unlockEthereumDeposit(_id);(_id);
    // }

    // //Wrapper function to test internal method
    // function callIsLockedEthereumDeposit(
    //     bytes32 _id
    // )
    //     public
    //     view
    //     returns(bool)
    // {
    //     return isLockedEthereumDeposit(_id);
    // }

    // //Wrapper function to test internal method
    // function callGetEthereumDeposit(
    //     bytes32 _id
    // )
    //     public
    //     view
    //     returns(address, bytes memory, address, uint256, uint256)
    // {
    //     return getEthereumDeposit(_id);
    // }

}
