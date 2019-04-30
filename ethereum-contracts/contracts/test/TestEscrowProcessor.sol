pragma solidity ^0.5.0;

import "../EscrowProcessor.sol";

contract TestEscrowProcessor is EscrowProcessor {

    event LogEscrowCreated(bytes32 _id);
    
    function() external payable {}

    //Wrapper function to test internal method
    function callCreateEscrow(
        address payable _sender,
        bytes memory _recipient,
        address _token,
        uint256 _amount
    )
        public
        returns(bytes32)
    {
        bytes32 escrowId = createEscrow(_sender, _recipient, _token, _amount);
        emit LogEscrowCreated(escrowId);
        return escrowId;
    }

    //Wrapper function to test internal method
    function callCompleteEscrow(
        bytes32 _escrowId
    )
        public
    {
        completeEscrow(_escrowId);
    }

    //Wrapper function to test internal method
    function callIsEscrow(
        bytes32 _escrowId
    )
        public
        view
        returns(bool)
    {
        return isEscrow(_escrowId);
    }

    //Wrapper function to test internal method
    function callGetEscrow(
        bytes32 _escrowId
    )
        public 
        view
        returns(address, bytes memory, address, uint256, uint256)
    {
        return getEscrow(_escrowId);
    }

}
