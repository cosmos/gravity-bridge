pragma solidity ^0.5.0;

import "../../../node_modules/openzeppelin-solidity/contracts/math/SafeMath.sol";
import "./BridgeToken.sol";

  /*
   *  @title: EthereumBank
   *  @dev: Ethereum bank which locks Ethereum/ERC20 token deposits, and unlocks
   *        Ethereum/ERC20 tokens once the prophecy has been successfully processed.
   */
contract EthereumBank {

    using SafeMath for uint256;

    uint256 public lockNonce;
    mapping(address => uint256) public lockedFunds;

    /*
    * @dev: Event declarations
    */
    event LogLock(
        address _from,
        bytes _to,
        address _token,
        string _symbol,
        uint256 _value,
        uint256 _nonce
    );

    event LogUnlock(
        address _to,
        address _token,
        string _symbol,
        uint256 _value
    );

    /*
    * @dev: Modifier declarations
    */

    modifier hasLockedFunds(
        address _token,
        uint256 _amount
    ) {
        require(
            lockedFunds[_token] >= _amount,
            "The Bank does not hold enough locked tokens to fulfill this request."
        );
        _;
    }

    modifier canDeliver(
        address _token,
        uint256 _amount
    )
    {
        if(_token == address(0)) {
            require(
                address(this).balance >= _amount,
                'Insufficient ethereum balance for delivery.'
            );
        } else {
            require(
                BridgeToken(_token).balanceOf(address(this)) >= _amount,
                'Insufficient ERC20 token balance for delivery.'
            );
        }
        _;
    }

    modifier availableNonce() {
        require(
            lockNonce + 1 > lockNonce,
            'No available nonces.'
        );
        _;
    }

    /*
    * @dev: Constructor which sets the lock nonce
    */
    constructor()
        public
    {
        lockNonce = 0;
    }

    /*
    * @dev: Creates a new Ethereum deposit with a unique id.
    *
    * @param _sender: The sender's ethereum address.
    * @param _recipient: The intended recipient's cosmos address.
    * @param _token: The currency type, either erc20 or ethereum.
    * @param _amount: The amount of erc20 tokens/ ethereum (in wei) to be itemized.
    */
    function lockFunds(
        address payable _sender,
        bytes memory _recipient,
        address _token,
        string memory _symbol,
        uint256 _amount
    )
        internal
    {
        // Incerment the lock nonce
        lockNonce = lockNonce.add(1);
        
        // Increment locked funds by the amount of tokens to be locked
        lockedFunds[_token] = lockedFunds[_token].add(_amount);

         emit LogLock(
            _sender,
            _recipient,
            _token,
            _symbol,
            _amount,
            lockNonce
        );
    }

    /*
    * @dev: Unlocks funds held on contract and sends them to the
    *       intended recipient
    *
    * @param _recipient: recipient's Ethereum address
    * @param _token: token contract address
    * @param _symbol: token symbol
    * @param _amount: wei amount or ERC20 token count
    */
    function unlockFunds(
        address payable _recipient,
        address _token,
        string memory _symbol,
        uint256 _amount
    )
        internal
    {
        // Decrement locked funds mapping by the amount of tokens to be unlocked
        lockedFunds[_token] = lockedFunds[_token].sub(_amount);

        // Transfer funds to intended recipient
        if (_token == address(0)) {
          _recipient.transfer(_amount);
        } else {
            require(
                BridgeToken(_token).transfer(_recipient, _amount),
                "Token transfer failed"
            );
        }

        emit LogUnlock(
            _recipient,
            _token,
            _symbol,
            _amount
        );
    }
}
