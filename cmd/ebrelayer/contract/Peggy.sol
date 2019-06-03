pragma solidity ^0.5.0;

/**
 * @title SafeMath
 * @dev Unsigned math operations with safety checks that revert on error
 */
library SafeMath {
    /**
    * @dev Multiplies two unsigned integers, reverts on overflow.
    */
    function mul(uint256 a, uint256 b) internal pure returns (uint256) {
        // Gas optimization: this is cheaper than requiring 'a' not being zero, but the
        // benefit is lost if 'b' is also tested.
        // See: https://github.com/OpenZeppelin/openzeppelin-solidity/pull/522
        if (a == 0) {
            return 0;
        }

        uint256 c = a * b;
        require(c / a == b);

        return c;
    }

    /**
    * @dev Integer division of two unsigned integers truncating the quotient, reverts on division by zero.
    */
    function div(uint256 a, uint256 b) internal pure returns (uint256) {
        // Solidity only automatically asserts when dividing by 0
        require(b > 0);
        uint256 c = a / b;
        // assert(a == b * c + a % b); // There is no case in which this doesn't hold

        return c;
    }

    /**
    * @dev Subtracts two unsigned integers, reverts on overflow (i.e. if subtrahend is greater than minuend).
    */
    function sub(uint256 a, uint256 b) internal pure returns (uint256) {
        require(b <= a);
        uint256 c = a - b;

        return c;
    }

    /**
    * @dev Adds two unsigned integers, reverts on overflow.
    */
    function add(uint256 a, uint256 b) internal pure returns (uint256) {
        uint256 c = a + b;
        require(c >= a);

        return c;
    }

    /**
    * @dev Divides two unsigned integers and returns the remainder (unsigned integer modulo),
    * reverts when dividing by zero.
    */
    function mod(uint256 a, uint256 b) internal pure returns (uint256) {
        require(b != 0);
        return a % b;
    }
}


/**
 * @title ERC20 interface
 * @dev see https://github.com/ethereum/EIPs/issues/20
 */
interface IERC20 {
    function transfer(address to, uint256 value) external returns (bool);

    function approve(address spender, uint256 value) external returns (bool);

    function transferFrom(address from, address to, uint256 value) external returns (bool);

    function totalSupply() external view returns (uint256);

    function balanceOf(address who) external view returns (uint256);

    function allowance(address owner, address spender) external view returns (uint256);

    event Transfer(address indexed from, address indexed to, uint256 value);

    event Approval(address indexed owner, address indexed spender, uint256 value);
}

/**
 * @title Standard ERC20 token
 *
 * @dev Implementation of the basic standard token.
 * https://github.com/ethereum/EIPs/blob/master/EIPS/eip-20.md
 * Originally based on code by FirstBlood:
 * https://github.com/Firstbloodio/token/blob/master/smart_contract/FirstBloodToken.sol
 *
 * This implementation emits additional Approval events, allowing applications to reconstruct the allowance status for
 * all accounts just by listening to said events. Note that this isn't required by the specification, and other
 * compliant implementations may not do it.
 */
contract ERC20 is IERC20 {
    using SafeMath for uint256;

    mapping (address => uint256) private _balances;

    mapping (address => mapping (address => uint256)) private _allowed;

    uint256 private _totalSupply;

    /**
    * @dev Total number of tokens in existence
    */
    function totalSupply() public view returns (uint256) {
        return _totalSupply;
    }

    /**
    * @dev Gets the balance of the specified address.
    * @param owner The address to query the balance of.
    * @return An uint256 representing the amount owned by the passed address.
    */
    function balanceOf(address owner) public view returns (uint256) {
        return _balances[owner];
    }

    /**
     * @dev Function to check the amount of tokens that an owner allowed to a spender.
     * @param owner address The address which owns the funds.
     * @param spender address The address which will spend the funds.
     * @return A uint256 specifying the amount of tokens still available for the spender.
     */
    function allowance(address owner, address spender) public view returns (uint256) {
        return _allowed[owner][spender];
    }

    /**
    * @dev Transfer token for a specified address
    * @param to The address to transfer to.
    * @param value The amount to be transferred.
    */
    function transfer(address to, uint256 value) public returns (bool) {
        _transfer(msg.sender, to, value);
        return true;
    }

    /**
     * @dev Approve the passed address to spend the specified amount of tokens on behalf of msg.sender.
     * Beware that changing an allowance with this method brings the risk that someone may use both the old
     * and the new allowance by unfortunate transaction ordering. One possible solution to mitigate this
     * race condition is to first reduce the spender's allowance to 0 and set the desired value afterwards:
     * https://github.com/ethereum/EIPs/issues/20#issuecomment-263524729
     * @param spender The address which will spend the funds.
     * @param value The amount of tokens to be spent.
     */
    function approve(address spender, uint256 value) public returns (bool) {
        require(spender != address(0));

        _allowed[msg.sender][spender] = value;
        emit Approval(msg.sender, spender, value);
        return true;
    }

    /**
     * @dev Transfer tokens from one address to another.
     * Note that while this function emits an Approval event, this is not required as per the specification,
     * and other compliant implementations may not emit the event.
     * @param from address The address which you want to send tokens from
     * @param to address The address which you want to transfer to
     * @param value uint256 the amount of tokens to be transferred
     */
    function transferFrom(address from, address to, uint256 value) public returns (bool) {
        _allowed[from][msg.sender] = _allowed[from][msg.sender].sub(value);
        _transfer(from, to, value);
        emit Approval(from, msg.sender, _allowed[from][msg.sender]);
        return true;
    }

    /**
     * @dev Increase the amount of tokens that an owner allowed to a spender.
     * approve should be called when allowed_[_spender] == 0. To increment
     * allowed value is better to use this function to avoid 2 calls (and wait until
     * the first transaction is mined)
     * From MonolithDAO Token.sol
     * Emits an Approval event.
     * @param spender The address which will spend the funds.
     * @param addedValue The amount of tokens to increase the allowance by.
     */
    function increaseAllowance(address spender, uint256 addedValue) public returns (bool) {
        require(spender != address(0));

        _allowed[msg.sender][spender] = _allowed[msg.sender][spender].add(addedValue);
        emit Approval(msg.sender, spender, _allowed[msg.sender][spender]);
        return true;
    }

    /**
     * @dev Decrease the amount of tokens that an owner allowed to a spender.
     * approve should be called when allowed_[_spender] == 0. To decrement
     * allowed value is better to use this function to avoid 2 calls (and wait until
     * the first transaction is mined)
     * From MonolithDAO Token.sol
     * Emits an Approval event.
     * @param spender The address which will spend the funds.
     * @param subtractedValue The amount of tokens to decrease the allowance by.
     */
    function decreaseAllowance(address spender, uint256 subtractedValue) public returns (bool) {
        require(spender != address(0));

        _allowed[msg.sender][spender] = _allowed[msg.sender][spender].sub(subtractedValue);
        emit Approval(msg.sender, spender, _allowed[msg.sender][spender]);
        return true;
    }

    /**
    * @dev Transfer token for a specified addresses
    * @param from The address to transfer from.
    * @param to The address to transfer to.
    * @param value The amount to be transferred.
    */
    function _transfer(address from, address to, uint256 value) internal {
        require(to != address(0));

        _balances[from] = _balances[from].sub(value);
        _balances[to] = _balances[to].add(value);
        emit Transfer(from, to, value);
    }

    /**
     * @dev Internal function that mints an amount of the token and assigns it to
     * an account. This encapsulates the modification of balances such that the
     * proper events are emitted.
     * @param account The account that will receive the created tokens.
     * @param value The amount that will be created.
     */
    function _mint(address account, uint256 value) internal {
        require(account != address(0));

        _totalSupply = _totalSupply.add(value);
        _balances[account] = _balances[account].add(value);
        emit Transfer(address(0), account, value);
    }

    /**
     * @dev Internal function that burns an amount of the token of a given
     * account.
     * @param account The account whose tokens will be burnt.
     * @param value The amount that will be burnt.
     */
    function _burn(address account, uint256 value) internal {
        require(account != address(0));

        _totalSupply = _totalSupply.sub(value);
        _balances[account] = _balances[account].sub(value);
        emit Transfer(account, address(0), value);
    }

    /**
     * @dev Internal function that burns an amount of the token of a given
     * account, deducting from the sender's allowance for said account. Uses the
     * internal burn function.
     * Emits an Approval event (reflecting the reduced allowance).
     * @param account The account whose tokens will be burnt.
     * @param value The amount that will be burnt.
     */
    function _burnFrom(address account, uint256 value) internal {
        _allowed[account][msg.sender] = _allowed[account][msg.sender].sub(value);
        _burn(account, value);
        emit Approval(account, msg.sender, _allowed[account][msg.sender]);
    }
}

  /*
   *  @title: Processor
   *  @dev: Processes requests for item locking and unlocking by
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
        bool locked;
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

    modifier onlySender(bytes32 _id) {
        require(
            msg.sender == items[_id].sender,
            'Must be the original sender.'
        );
        _;
    }

    modifier canDeliver(bytes32 _id) {
        if(items[_id].token == address(0)) {
            require(
                address(this).balance >= items[_id].amount,
                'Insufficient ethereum balance for delivery.'
            );
        } else {
            require(
                ERC20(items[_id].token).balanceOf(address(this)) >= items[_id].amount,
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
    *       original sender and unlocking the item.
    *
    * @param _id: The item to be completed.
    */
    function complete(
        bytes32 _id
    )
        internal
        canDeliver(_id)
        returns(address payable, address, uint256, uint256)
    {
        require(isLocked(_id));

        //Get locked item's attributes for return
        address payable sender = items[_id].sender;
        address token = items[_id].token;
        uint256 amount = items[_id].amount;
        uint256 uniqueNonce = items[_id].nonce;

        //Update lock status
        items[_id].locked = false;

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
    function isLocked(
        bytes32 _id
    )
        internal 
        view
        returns(bool)
    {
        return(items[_id].locked);
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
        bytes32 _id
    )
        internal 
        view
        returns(address payable, bytes memory, address, uint256, uint256)
    {
        Item memory item = items[_id];

        return(
            item.sender,
            item.recipient,
            item.token,
            item.amount,
            item.nonce
        );
    }
}

  /*
   *  @title: Peggy
   *  @dev: Peg zone contract for testing one-way transfers from Ethereum
   *        to Cosmos, facilitated by a trusted relayer. This contract is
   *        NOT intended to be used in production and users are empowered
   *        to withdraw their locked funds at any time.
   */
contract Peggy is Processor {

    bool public active;
    address public relayer;
    mapping(bytes32 => bool) public ids;

    event LogLock(
        bytes32 _id,
        address _from,
        bytes _to,
        address _token,
        uint256 _value,
        uint256 _nonce
    );

    event LogUnlock(
        bytes32 _id,
        address _to,
        address _token,
        uint256 _value,
        uint256 _nonce
    );

    event LogWithdraw(
        bytes32 _id,
        address _to,
        address _token,
        uint256 _value,
        uint256 _nonce
    );

    event LogLockingPaused(
        uint256 _time
    );

    event LogLockingActivated(
        uint256 _time
    );

    /*
    * @dev: Modifier to restrict access to the relayer.
    *
    */
    modifier onlyRelayer()
    {
        require(
            msg.sender == relayer,
            'Must be the specified relayer.'
        );
        _;
    }

    /*
    * @dev: Modifier which restricts lock functionality when paused.
    *
    */
    modifier whileActive()
    {
        require(
            active == true,
            'Lock functionality is currently paused.'
        );
        _;
    }
    /*
    * @dev: Constructor, initalizes relayer and active status.
    *
    */
    constructor()
        public
    {
        relayer = msg.sender;
        active = true;
        emit LogLockingActivated(now);
    }

    /* 
     * @dev: Locks received funds and creates new items.
     *
     * @param _recipient: bytes representation of destination address.
     * @param _token: token address in origin chain (0x0 if ethereum)
     * @param _amount: value of item
     */
    function lock(
        bytes memory _recipient,
        address _token,
        uint256 _amount
    )
        public
        payable
        availableNonce()
        whileActive()
        returns(bytes32 _id)
    {
         //Actions based on token address type
        if (msg.value != 0) {
          require(_token == address(0));
          require(msg.value == _amount);
        } else {
          require(ERC20(_token).transferFrom(msg.sender, address(this), _amount));
        }

        //Create an item with a unique key.
        bytes32 id = create(
            msg.sender,
            _recipient,
            _token,
            _amount
        );

        emit LogLock(
            id,
            msg.sender,
            _recipient,
            _token,
            _amount,
            getNonce()
        );

        return id;
    }

    /*
     * @dev: Unlocks ethereum/erc20 tokens, called by relayer.
     *
     *       This is a shortcut utility method for testing purposes.
     *       In the future bidirectional system, unlocking functionality
     *       will be guarded by validator signatures.
     *
     * @param _id: Unique key of the item.
     */
    function unlock(
        bytes32 _id
    )
        onlyRelayer
        canDeliver(_id)
        external
        returns (bool)
    {
        require(isLocked(_id));

        // Transfer item's funds and unlock it
        (address payable sender,
            address token,
            uint256 amount,
            uint256 uniqueNonce) = complete(_id);

        //Emit unlock event
        emit LogUnlock(
            _id,
            sender,
            token,
            amount,
            uniqueNonce
        );
        return true;
    }

    /*
     * @dev: Withdraws ethereum/erc20 tokens, called original sender.
     *
     *       This is a backdoor utility method included for testing,
     *       purposes, allowing users to withdraw their funds. This
     *       functionality will be removed in production.
     *
     * @param _id: Unique key of the item.
     */
    function withdraw(
        bytes32 _id
    )
        onlySender(_id)
        canDeliver(_id)
        external
        returns (bool)
    {
        require(isLocked(_id));

        // Transfer item's funds and unlock it
        (address payable sender,
            address token,
            uint256 amount,
            uint256 uniqueNonce) = complete(_id);

        //Emit withdraw event
        emit LogWithdraw(
            _id,
            sender,
            token,
            amount,
            uniqueNonce
        );

        return true;
    }

    /*
    * @dev: Exposes an item's current status.
    *
    * @param _id: The item in question.
    * @return: Boolean indicating the lock status.
    */
    function getStatus(
        bytes32 _id
    )
        public 
        view
        returns(bool)
    {
        return isLocked(_id);
    }

    /*
    * @dev: Allows access to an item's information via its unique identifier.
    *
    * @param _id: The item to be viewed.
    * @return: Original sender's address.
    * @return: Intended receiver's address in bytes.
    * @return: The token's address.
    * @return: The amount locked in the item.
    * @return: The item's unique nonce.
    */
    function viewItem(
        bytes32 _id
    )
        public 
        view
        returns(address, bytes memory, address, uint256, uint256)
    {
        return getItem(_id);
    }

    /*
    * @dev: Relayer can pause fund locking without impacting other functionality.
    */
    function pauseLocking()
        public
        onlyRelayer
    {
        require(active);
        active = false;
        emit LogLockingPaused(now);
    }

    /*
    * @dev: Relayer can activate fund locking without impacting other functionality.
    */
    function activateLocking()
        public
        onlyRelayer
    {
        require(!active);
        active = true;
        emit LogLockingActivated(now);
    }
}
