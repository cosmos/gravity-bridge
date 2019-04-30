pragma solidity ^0.5.0;

import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "openzeppelin-solidity/contracts/token/ERC20/ERC20.sol";

  /*
   *  @title: EscrowProcessor
   *  @dev: Processes requests for escrow creation and deletion by
   *        storing an escrow's information then relaying the funds
   *        the original sender.
   */
contract EscrowProcessor {

    using SafeMath for uint256;

    /*
    * @dev: Escrow struct to store information.
    */    
    struct Escrow {
        address payable sender;
        bytes recipient;
        address token;
        uint256 amount;
        uint256 nonce;
        bool isEscrow;
    }

    uint256 public nonce;
    mapping(bytes32 => Escrow) private escrows;

    /*
    * @dev: Constructor, initalizes escrow count.
    */
    constructor() 
        public
    {
        nonce = 0;
    }

    modifier onlySender(bytes32 _escrowId) {
        require(
            msg.sender == escrows[_escrowId].sender,
            'Must be the original sender of the escrow.'
        );
        _;
    }

    modifier canDeliver(bytes32 _escrowId) {
        if(escrows[_escrowId].token == address(0)) {
            require(
                address(this).balance >= escrows[_escrowId].amount,
                'Insufficient ethereum balance for delivery.'
            );
        } else {
            require(
                ERC20(escrows[_escrowId].token).balanceOf(address(this)) >= escrows[_escrowId].amount,
                'Insufficient ERC20 token balance for delivery.'
            );            
        }
        _;
    }
  
    modifier availableNonce() {
        require(
            nonce + 1 > nonce,
            'No available nonces.'
        );
        _;
    }

    /*
    * @dev: Creates an escrow with a unique id.
    *
    * @param _sender: The sender's ethereum address.
    * @param _recipient: The intended recipient's cosmos address.
    * @param _token: The currency type, either erc20 or ethereum.
    * @param _amount: The amount of erc20 tokens/ ethereum (in wei) to be escrowed.
    * @return: The newly created escrow's unique id.
    */
    function createEscrow(
        address payable _sender,
        bytes memory _recipient,
        address _token,
        uint256 _amount
    )
        internal
        returns(bytes32)
    {
        nonce++;

        bytes32 escrowKey = keccak256(
            abi.encodePacked(
                _sender,
                _recipient,
                _token,
                _amount,
                nonce
            )
        );
        
        escrows[escrowKey] = Escrow(
            _sender,
            _recipient,
            _token,
            _amount,
            nonce,
            true
        );

        return escrowKey;
    }

    /*
    * @dev: Completes the escrow by sending the funds to the
    *       original sender and deleting the escrow.
    *
    * @param _escrowId: The escrow to be completed.
    */
    function completeEscrow(
        bytes32 _escrowId
    )
        internal
        canDeliver(_escrowId)
        returns(address payable, address, uint256, uint256)
    {
        require(isEscrow(_escrowId));

        address payable sender = escrows[_escrowId].sender;
        address token = escrows[_escrowId].token;
        uint256 amount = escrows[_escrowId].amount;
        uint256 escrowNonce = escrows[_escrowId].nonce;
        
        //Delete escrow
        delete(escrows[_escrowId]);

        //Transfers based on token address type
        if (token == address(0)) {
          sender.transfer(amount);
        } else {
          require(ERC20(token).transfer(sender, amount));
        }       

        return(sender, token, amount, escrowNonce);
    }

    /*
    * @dev: Checks the current nonce.
    *
    * @return: The current nonce.
    */
    function getNonce()
        internal
        view
        returns(uint256)
    {
        return nonce;
    }

    /*
    * @dev: Checks if an individual escrow exists.
    *
    * @param _escrowId: The unique escrow's id.
    * @return: Boolean indicating if the escrow exists in memory.
    */
    function isEscrow(
        bytes32 _escrowId
    )
        internal 
        view
        returns(bool)
    {
        return(escrows[_escrowId].isEscrow);
    }

    /*
    * @dev: Gets an escrow's information
    *
    * @param _escrowId: The escrow containing the desired information.
    * @return: Sender's address.
    * @return: Recipient's address in bytes.
    * @return: Token address.
    * @return: Amount of ethereum/erc20 in the escrow.
    * @return: Unique nonce of the escrow.
    */
    function getEscrow(
        bytes32 _escrowId
    )
        internal 
        view
        returns(address payable, bytes memory, address, uint256, uint256)
    {
        Escrow memory escrow = escrows[_escrowId];

        return(
            escrow.sender,
            escrow.recipient,
            escrow.token,
            escrow.amount,
            escrow.nonce
        );
    }
}
