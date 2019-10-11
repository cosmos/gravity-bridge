pragma solidity ^0.5.0;

import "../../node_modules/openzeppelin-solidity/contracts/math/SafeMath.sol";
import "../../node_modules/openzeppelin-solidity/contracts/token/ERC20/ERC20.sol";
import "./TokenERC20.sol";

  /*
   *  @title: EthereumBank
   *  @dev: EthereumBank requests for deposit locking and unlocking by
   *        storing an item's information then relaying the funds
   *        the original sender.
   */
contract EthereumBank {

    using SafeMath for uint256;

    /*
    * @dev: EthereumDeposit struct to store information.
    */
    struct EthereumDeposit {
        address payable sender;
        bytes recipient;
        address token;
        uint256 amount;
        uint256 nonce;
        bool locked;
    }

    uint256 public nonce;
    mapping(bytes32 => EthereumDeposit) private ethereumDeposits;

    /*
    * @dev: Constructor, initalizes item count.
    */
    constructor()
        public
    {
        nonce = 0;
    }

    modifier onlySender(bytes32 _id) {
        require(
            msg.sender == ethereumDeposits[_id].sender,
            'Must be the original sender.'
        );
        _;
    }

    modifier canDeliver(bytes32 _id) {
        if(ethereumDeposits[_id].token == address(0)) {
            require(
                address(this).balance >= ethereumDeposits[_id].amount,
                'Insufficient ethereum balance for delivery.'
            );
        } else {
            require(
                TokenERC20(ethereumDeposits[_id].token).balanceOf(address(this)) >= ethereumDeposits[_id].amount,
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
    function newEthereumDeposit(
        address payable _sender,
        bytes memory _recipient,
        address _token,
        uint256 _amount
    )
        internal
        returns(bytes32)
    {
        nonce++;

        bytes32 depositID = keccak256(
            abi.encodePacked(
                _sender,
                _recipient,
                _token,
                _amount,
                nonce
            )
        );

        ethereumDeposits[depositID] = EthereumDeposit(
            _sender,
            _recipient,
            _token,
            _amount,
            nonce,
            true
        );

        return depositID;
    }

    /*
    * @dev: Completes the deposit by sending the funds to the
    *       original sender and unlocking the item.
    *
    * @param _id: The item to be completed.
    */
    function unlockEthereumDeposit(
        bytes32 _id
    )
        internal
        canDeliver(_id)
        returns(address payable, address, uint256, uint256)
    {
        require(
            isLockedEthereumDeposit(_id),
            "The funds must currently be locked."
        );

        //Get locked deposit's attributes for return
        address payable sender = ethereumDeposits[_id].sender;
        address token = ethereumDeposits[_id].token;
        uint256 amount = ethereumDeposits[_id].amount;
        uint256 uniqueNonce = ethereumDeposits[_id].nonce;

        //Update lock status
        ethereumDeposits[_id].locked = false;

        //Transfers based on token address type
        if (token == address(0)) {
          sender.transfer(amount);
        } else {
          require(
              TokenERC20(token).transfer(sender, amount),
              "Token transfer failed, check contract token allowances and try again."
            );
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
    function isLockedEthereumDeposit(
        bytes32 _id
    )
        internal
        view
        returns(bool)
    {
        return(ethereumDeposits[_id].locked);
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
    function getEthereumDeposit(
        bytes32 _id
    )
        internal
        view
        returns(address payable, bytes memory, address, uint256, uint256)
    {
        EthereumDeposit memory deposit = ethereumDeposits[_id];

        return(
            deposit.sender,
            deposit.recipient,
            deposit.token,
            deposit.amount,
            deposit.nonce
        );
    }
}
