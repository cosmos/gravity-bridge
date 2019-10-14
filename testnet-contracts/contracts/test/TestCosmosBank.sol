pragma solidity ^0.5.0;

import "../CosmosBank.sol";

contract TestCosmosBank is CosmosBank {

    function() external payable {}

    // //Wrapper function to test internal method
    // function callNewCosmosDeposit(
    //     bytes memory _cosmosSender,
    //     address payable _ethereumRecipient,
    //     address _token,
    //     uint256 _amount
    // )
    //     public
    //     returns(bytes32)
    // {
    //     return newCosmosDeposit(
    //         _cosmosSender,
    //         _ethereumRecipient,
    //         _token,
    //         _amount
    //     );
    // }

    // //Wrapper function to test internal method
    // function callMintCosmosToken(
    //     bytes memory _cosmosSender,
    //     address payable _ethereumRecipient,
    //     address _token,
    //     string memory _symbol,
    //     uint256 _amount
    // )
    //     public
    // {
    //     return mintCosmosToken(
    //         _cosmosSender,
    //         _ethereumRecipient,
    //         _token,
    //         _symbol,
    //         _amount
    //     );
    // }

    // //Wrapper function to test internal method
    // function callDeployNewCosmosToken(
    //     string memory _symbol
    // )
    //     public
    //     returns(address)
    // {
    //     return deployNewCosmosToken(
    //         _symbol
    //     );
    // }
}
