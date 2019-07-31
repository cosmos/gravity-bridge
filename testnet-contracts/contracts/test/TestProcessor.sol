pragma solidity ^0.5.0;

import "../Processor.sol";

contract TestProcessor is Processor {

    event LogItemCreated(bytes32 _id);
    
    function() external payable {}

    //Wrapper function to test internal method
    function callCreate(
        address payable _sender,
        bytes memory _recipient,
        address _token,
        uint256 _amount
    )
        public
        returns(bytes32)
    {
        bytes32 id = create(_sender, _recipient, _token, _amount);
        emit LogItemCreated(id);
        return id;
    }

    //Wrapper function to test internal method
    function callComplete(
        bytes32 _id
    )
        public
    {
        complete(_id);
    }

    //Wrapper function to test internal method
    function callIsLocked(
        bytes32 _id
    )
        public
        view
        returns(bool)
    {
        return isLocked(_id);
    }

    //Wrapper function to test internal method
    function callGetItem(
        bytes32 _id
    )
        public 
        view
        returns(address, bytes memory, address, uint256, uint256)
    {
        return getItem(_id);
    }

}
