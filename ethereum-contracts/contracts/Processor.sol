pragma solidity ^0.5.0;

import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "openzeppelin-solidity/contracts/token/ERC20/ERC20.sol";

  /*
   *  @title: Processor
   *  @dev: Processes requests for item creation and deletion by
   *        storing an item's information then relaying the funds
   *        the original sender.
   */
contract Processor {

    using SafeMath for uint256;

    /*
    * @dev: Item struct to store information.
    */    
    struct Item {
        address payable sender;
        bytes recipient;
        address token;
        uint256 amount;
        uint256 nonce;
        bool isItem;
    }

    uint256 public nonce;
    mapping(bytes32 => Item) private items;

    /*
    * @dev: Constructor, initalizes item count.
    */
    constructor() 
        public
    {
        nonce = 0;
    }

    modifier onlySender(bytes32 _itemId) {
        require(
            msg.sender == items[_itemId].sender,
            'Must be the original sender.'
        );
        _;
    }

    modifier canDeliver(bytes32 _itemId) {
        if(items[_itemId].token == address(0)) {
            require(
                address(this).balance >= items[_itemId].amount,
                'Insufficient ethereum balance for delivery.'
            );
        } else {
            require(
                ERC20(items[_itemId].token).balanceOf(address(this)) >= items[_itemId].amount,
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
    * @dev: Creates an item with a unique id.
    *
    * @param _sender: The sender's ethereum address.
    * @param _recipient: The intended recipient's cosmos address.
    * @param _token: The currency type, either erc20 or ethereum.
    * @param _amount: The amount of erc20 tokens/ ethereum (in wei) to be itemized.
    * @return: The newly created item's unique id.
    */
    function create(
        address payable _sender,
        bytes memory _recipient,
        address _token,
        uint256 _amount
    )
        internal
        returns(bytes32)
    {
        nonce++;

        bytes32 itemKey = keccak256(
            abi.encodePacked(
                _sender,
                _recipient,
                _token,
                _amount,
                nonce
            )
        );
        
        items[itemKey] = Item(
            _sender,
            _recipient,
            _token,
            _amount,
            nonce,
            true
        );

        return itemKey;
    }

    /*
    * @dev: Completes the item by sending the funds to the
    *       original sender and deleting the item.
    *
    * @param _id: The item to be completed.
    */
    function complete(
        bytes32 _itemId
    )
        internal
        canDeliver(_itemId)
        returns(address payable, address, uint256, uint256)
    {
        require(isItem(_itemId));

        address payable sender = items[_itemId].sender;
        address token = items[_itemId].token;
        uint256 amount = items[_itemId].amount;
        uint256 uniqueNonce = items[_itemId].nonce;
        
        //Delete item
        delete(items[_itemId]);

        //Transfers based on token address type
        if (token == address(0)) {
          sender.transfer(amount);
        } else {
          require(ERC20(token).transfer(sender, amount));
        }       

        return(sender, token, amount, uniqueNonce);
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
    * @dev: Checks if an individual item exists.
    *
    * @param _id: The unique item's id.
    * @return: Boolean indicating if the item exists in memory.
    */
    function isItem(
        bytes32 _itemId
    )
        internal 
        view
        returns(bool)
    {
        return(items[_itemId].isItem);
    }

    /*
    * @dev: Gets an item's information
    *
    * @param _Id: The item containing the desired information.
    * @return: Sender's address.
    * @return: Recipient's address in bytes.
    * @return: Token address.
    * @return: Amount of ethereum/erc20 in the item.
    * @return: Unique nonce of the item.
    */
    function getItem(
        bytes32 _itemId
    )
        internal 
        view
        returns(address payable, bytes memory, address, uint256, uint256)
    {
        Item memory item = items[_itemId];

        return(
            item.sender,
            item.recipient,
            item.token,
            item.amount,
            item.nonce
        );
    }
}
