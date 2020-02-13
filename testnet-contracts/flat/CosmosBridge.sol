
// File: openzeppelin-solidity/contracts/math/SafeMath.sol

pragma solidity ^0.5.0;

/**
 * @dev Wrappers over Solidity's arithmetic operations with added overflow
 * checks.
 *
 * Arithmetic operations in Solidity wrap on overflow. This can easily result
 * in bugs, because programmers usually assume that an overflow raises an
 * error, which is the standard behavior in high level programming languages.
 * `SafeMath` restores this intuition by reverting the transaction when an
 * operation overflows.
 *
 * Using this library instead of the unchecked operations eliminates an entire
 * class of bugs, so it's recommended to use it always.
 */
library SafeMath {
    /**
     * @dev Returns the addition of two unsigned integers, reverting on
     * overflow.
     *
     * Counterpart to Solidity's `+` operator.
     *
     * Requirements:
     * - Addition cannot overflow.
     */
    function add(uint256 a, uint256 b) internal pure returns (uint256) {
        uint256 c = a + b;
        require(c >= a, "SafeMath: addition overflow");

        return c;
    }

    /**
     * @dev Returns the subtraction of two unsigned integers, reverting on
     * overflow (when the result is negative).
     *
     * Counterpart to Solidity's `-` operator.
     *
     * Requirements:
     * - Subtraction cannot overflow.
     */
    function sub(uint256 a, uint256 b) internal pure returns (uint256) {
        return sub(a, b, "SafeMath: subtraction overflow");
    }

    /**
     * @dev Returns the subtraction of two unsigned integers, reverting with custom message on
     * overflow (when the result is negative).
     *
     * Counterpart to Solidity's `-` operator.
     *
     * Requirements:
     * - Subtraction cannot overflow.
     *
     * _Available since v2.4.0._
     */
    function sub(uint256 a, uint256 b, string memory errorMessage) internal pure returns (uint256) {
        require(b <= a, errorMessage);
        uint256 c = a - b;

        return c;
    }

    /**
     * @dev Returns the multiplication of two unsigned integers, reverting on
     * overflow.
     *
     * Counterpart to Solidity's `*` operator.
     *
     * Requirements:
     * - Multiplication cannot overflow.
     */
    function mul(uint256 a, uint256 b) internal pure returns (uint256) {
        // Gas optimization: this is cheaper than requiring 'a' not being zero, but the
        // benefit is lost if 'b' is also tested.
        // See: https://github.com/OpenZeppelin/openzeppelin-contracts/pull/522
        if (a == 0) {
            return 0;
        }

        uint256 c = a * b;
        require(c / a == b, "SafeMath: multiplication overflow");

        return c;
    }

    /**
     * @dev Returns the integer division of two unsigned integers. Reverts on
     * division by zero. The result is rounded towards zero.
     *
     * Counterpart to Solidity's `/` operator. Note: this function uses a
     * `revert` opcode (which leaves remaining gas untouched) while Solidity
     * uses an invalid opcode to revert (consuming all remaining gas).
     *
     * Requirements:
     * - The divisor cannot be zero.
     */
    function div(uint256 a, uint256 b) internal pure returns (uint256) {
        return div(a, b, "SafeMath: division by zero");
    }

    /**
     * @dev Returns the integer division of two unsigned integers. Reverts with custom message on
     * division by zero. The result is rounded towards zero.
     *
     * Counterpart to Solidity's `/` operator. Note: this function uses a
     * `revert` opcode (which leaves remaining gas untouched) while Solidity
     * uses an invalid opcode to revert (consuming all remaining gas).
     *
     * Requirements:
     * - The divisor cannot be zero.
     *
     * _Available since v2.4.0._
     */
    function div(uint256 a, uint256 b, string memory errorMessage) internal pure returns (uint256) {
        // Solidity only automatically asserts when dividing by 0
        require(b > 0, errorMessage);
        uint256 c = a / b;
        // assert(a == b * c + a % b); // There is no case in which this doesn't hold

        return c;
    }

    /**
     * @dev Returns the remainder of dividing two unsigned integers. (unsigned integer modulo),
     * Reverts when dividing by zero.
     *
     * Counterpart to Solidity's `%` operator. This function uses a `revert`
     * opcode (which leaves remaining gas untouched) while Solidity uses an
     * invalid opcode to revert (consuming all remaining gas).
     *
     * Requirements:
     * - The divisor cannot be zero.
     */
    function mod(uint256 a, uint256 b) internal pure returns (uint256) {
        return mod(a, b, "SafeMath: modulo by zero");
    }

    /**
     * @dev Returns the remainder of dividing two unsigned integers. (unsigned integer modulo),
     * Reverts with custom message when dividing by zero.
     *
     * Counterpart to Solidity's `%` operator. This function uses a `revert`
     * opcode (which leaves remaining gas untouched) while Solidity uses an
     * invalid opcode to revert (consuming all remaining gas).
     *
     * Requirements:
     * - The divisor cannot be zero.
     *
     * _Available since v2.4.0._
     */
    function mod(uint256 a, uint256 b, string memory errorMessage) internal pure returns (uint256) {
        require(b != 0, errorMessage);
        return a % b;
    }
}

// File: openzeppelin-solidity/contracts/cryptography/ECDSA.sol

pragma solidity ^0.5.0;

/**
 * @dev Elliptic Curve Digital Signature Algorithm (ECDSA) operations.
 *
 * These functions can be used to verify that a message was signed by the holder
 * of the private keys of a given address.
 */
library ECDSA {
    /**
     * @dev Returns the address that signed a hashed message (`hash`) with
     * `signature`. This address can then be used for verification purposes.
     *
     * The `ecrecover` EVM opcode allows for malleable (non-unique) signatures:
     * this function rejects them by requiring the `s` value to be in the lower
     * half order, and the `v` value to be either 27 or 28.
     *
     * NOTE: This call _does not revert_ if the signature is invalid, or
     * if the signer is otherwise unable to be retrieved. In those scenarios,
     * the zero address is returned.
     *
     * IMPORTANT: `hash` _must_ be the result of a hash operation for the
     * verification to be secure: it is possible to craft signatures that
     * recover to arbitrary addresses for non-hashed data. A safe way to ensure
     * this is by receiving a hash of the original message (which may otherwise
     * be too long), and then calling {toEthSignedMessageHash} on it.
     */
    function recover(bytes32 hash, bytes memory signature) internal pure returns (address) {
        // Check the signature length
        if (signature.length != 65) {
            return (address(0));
        }

        // Divide the signature in r, s and v variables
        bytes32 r;
        bytes32 s;
        uint8 v;

        // ecrecover takes the signature parameters, and the only way to get them
        // currently is to use assembly.
        // solhint-disable-next-line no-inline-assembly
        assembly {
            r := mload(add(signature, 0x20))
            s := mload(add(signature, 0x40))
            v := byte(0, mload(add(signature, 0x60)))
        }

        // EIP-2 still allows signature malleability for ecrecover(). Remove this possibility and make the signature
        // unique. Appendix F in the Ethereum Yellow paper (https://ethereum.github.io/yellowpaper/paper.pdf), defines
        // the valid range for s in (281): 0 < s < secp256k1n ÷ 2 + 1, and for v in (282): v ∈ {27, 28}. Most
        // signatures from current libraries generate a unique signature with an s-value in the lower half order.
        //
        // If your library generates malleable signatures, such as s-values in the upper range, calculate a new s-value
        // with 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141 - s1 and flip v from 27 to 28 or
        // vice versa. If your library also generates signatures with 0/1 for v instead 27/28, add 27 to v to accept
        // these malleable signatures as well.
        if (uint256(s) > 0x7FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF5D576E7357A4501DDFE92F46681B20A0) {
            return address(0);
        }

        if (v != 27 && v != 28) {
            return address(0);
        }

        // If the signature is valid (and not malleable), return the signer address
        return ecrecover(hash, v, r, s);
    }

    /**
     * @dev Returns an Ethereum Signed Message, created from a `hash`. This
     * replicates the behavior of the
     * https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sign[`eth_sign`]
     * JSON-RPC method.
     *
     * See {recover}.
     */
    function toEthSignedMessageHash(bytes32 hash) internal pure returns (bytes32) {
        // 32 is the length in bytes of hash,
        // enforced by the type signature above
        return keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", hash));
    }
}

// File: contracts/Valset.sol

pragma solidity ^0.5.0;



contract Valset {

    using SafeMath for uint256;
    using ECDSA for bytes32;

    /*
    * @dev: Variable declarations
    */
    address public operator;
    uint256 public totalPower;
    uint256 public currentValsetVersion;
    uint256 public validatorCount;
    mapping (bytes32 => bool) public validators;
    mapping(bytes32 => uint256) public powers;

    /*
    * @dev: Event declarations
    */
    event LogValidatorAdded(
        address _validator,
        uint256 _power,
        uint256 _currentValsetVersion,
        uint256 _validatorCount,
        uint256 _totalPower
    );

    event LogValidatorPowerUpdated(
        address _validator,
        uint256 _power,
        uint256 _currentValsetVersion,
        uint256 _validatorCount,
        uint256 _totalPower
    );

    event LogValidatorRemoved(
        address _validator,
        uint256 _power,
        uint256 _currentValsetVersion,
        uint256 _validatorCount,
        uint256 _totalPower
    );

    event LogValsetReset(
        uint256 _newValsetVersion,
        uint256 _validatorCount,
        uint256 _totalPower
    );

    event LogValsetUpdated(
        uint256 _newValsetVersion,
        uint256 _validatorCount,
        uint256 _totalPower
    );

    /*
    * @dev: Modifier which restricts access to the operator.
    */
    modifier onlyOperator()
    {
        require(
            msg.sender == operator,
            'Must be the operator.'
        );
        _;
    }

    /*
    * @dev: Constructor
    */
    constructor(
        address _operator,
        address[] memory _initValidators,
        uint256[] memory _initPowers
    )
        public
    {
        operator = _operator;
        currentValsetVersion = 0;

        updateValset(
            _initValidators,
            _initPowers
        );
    }

    function recover(
        string memory _message,
        bytes memory _signature
    )
        public
        pure
        returns (address)
    {
        bytes32 message = ethMessageHash(_message);
        return verify(message, _signature);
    }

    /*
    * @dev: addValidator
    */
    function addValidator(
        address _validatorAddress,
        uint256 _validatorPower
    )
        public
        onlyOperator
    {
        addValidatorInternal(
            _validatorAddress,
            _validatorPower
        );
    }

    /*
    * @dev: updateValidatorPower
    */
    function updateValidatorPower(
        address _validatorAddress,
        uint256 _newValidatorPower
    )
        public
        onlyOperator
    {
        // Create a unique key which for this validator's position in the current version of the mapping
        bytes32 key = keccak256(
            abi.encodePacked(
                currentValsetVersion,
                _validatorAddress
            )
        );

        require(
            validators[key],
            "Can only update the power of active valdiators"
        );


        // Adjust total power by new validator power
        uint256 priorPower = powers[key];
        totalPower = totalPower.sub(priorPower);
        totalPower = totalPower.add(_newValidatorPower);

        // Set validator's new power
        powers[key] = _newValidatorPower;

        emit LogValidatorPowerUpdated(
            _validatorAddress,
            _newValidatorPower,
            currentValsetVersion,
            validatorCount,
            totalPower
        );
    }

    /*
    * @dev: removeValidator
    */
    function removeValidator(
        address _validatorAddress
    )
        public
        onlyOperator
    {
        // Create a unique key which for this validator's position in the current version of the mapping
        bytes32 key = keccak256(
            abi.encodePacked(
                currentValsetVersion,
                _validatorAddress
            )
        );

        require(
            validators[key],
            "Can only remove active valdiators"
        );

        // Update validator count and total power
        validatorCount = validatorCount.sub(1);
        totalPower = totalPower.sub(powers[key]);

        // Delete validator and power
        delete validators[key];
        delete powers[key];

        emit LogValidatorRemoved(
            _validatorAddress,
            0,
            currentValsetVersion,
            validatorCount,
            totalPower
        );
    }

    /*
    * @dev: updateValset
    */
    function updateValset(
        address[] memory _validators,
        uint256[] memory _powers
    )
        public
        onlyOperator
    {
       require(
           _validators.length == _powers.length,
           "Every validator must have a corresponding power"
       );

       resetValset();

       for(uint256 i = 0; i < _validators.length; i = i.add(1)) {
           addValidatorInternal(_validators[i], _powers[i]);
       }

        emit LogValsetUpdated(
            currentValsetVersion,
            validatorCount,
            totalPower
        );
    }

    /*
    * @dev: isActiveValidator
    */
    function isActiveValidator(
        address _validatorAddress
    )
        public
        view
        returns(bool)
    {
        // Recreate the unique key for this address given the current mapping version
        bytes32 key = keccak256(
            abi.encodePacked(
                currentValsetVersion,
                _validatorAddress
            )
        );

        // Return bool indicating if this address is an active validator
        return validators[key];
    }

    /*
    * @dev: getValidatorPower
    */
    function getValidatorPower(
        address _validatorAddress
    )
        external
        view
        returns(uint256)
    {
        // Recreate the unique key for this address given the current mapping version
        bytes32 key = keccak256(
            abi.encodePacked(
                currentValsetVersion,
                _validatorAddress
            )
        );

        return powers[key];
    }

    /*
    * @dev: recoverGas
    */
    function recoverGas(
        uint256 _valsetVersion,
        address _validatorAddress
    )
        external
        onlyOperator
    {
        require(
            _valsetVersion < currentValsetVersion,
            "Gas recovery only allowed for previous validator sets"
        );

        // Recreate the unique key used to identify this validator in the given version
        bytes32 key = keccak256(
            abi.encodePacked(
                _valsetVersion,
                _validatorAddress
            )
        );

        // Delete from mappings and recover gas
        delete(validators[key]);
        delete(powers[key]);
    }

    /*
    * @dev: addValidatorInternal
    */
    function addValidatorInternal(
        address _validatorAddress,
        uint256 _validatorPower
    )
        internal
    {
        // Create a unique key which for this validator's position in the current version of the mapping
        bytes32 key = keccak256(
            abi.encodePacked(
                currentValsetVersion,
                _validatorAddress
            )
        );

        validatorCount = validatorCount.add(1);
        totalPower = totalPower.add(_validatorPower);

        // Set validator as active and set their power
        validators[key] = true;
        powers[key] = _validatorPower;

        emit LogValidatorAdded(
            _validatorAddress,
            _validatorPower,
            currentValsetVersion,
            validatorCount,
            totalPower
        );
    }

    /*
    * @dev: resetValset
    */
    function resetValset()
        internal
    {
        currentValsetVersion = currentValsetVersion.add(1);
        validatorCount = 0;
        totalPower = 0;

        emit LogValsetReset(
            currentValsetVersion,
            validatorCount,
            totalPower
        );
    }

  /*
    * @dev: Verify
    *
    */
    function verify(
        bytes32 h,
        bytes memory signature
    )
        internal
        pure
        returns (address)
    {
        bytes32 r;
        bytes32 s;
        uint8 v;

        // Check the signature length
        if (signature.length != 65) {
            return (address(0));
        }

        // Divide the signature in r, s and v variables
        // ecrecover takes the signature parameters, and the only way to get them
        // currently is to use assembly.
        // solium-disable-next-line security/no-inline-assembly
        assembly {
            r := mload(add(signature, 32))
            s := mload(add(signature, 64))
            v := byte(0, mload(add(signature, 96)))
        }

        // Version of signature should be 27 or 28, but 0 and 1 are also possible versions
        if (v < 27) {
            v += 27;
        }

        // If the version is correct return the signer address
        if (v != 27 && v != 28) {
            return (address(0));
        } else {
            // solium-disable-next-line arg-overflow
            return ecrecover(h, v, r, s);
        }
    }

    /**
    * @dev prefix a bytes32 value with "\x19Ethereum Signed Message:" and hash the result
    */
    function ethMessageHash(string memory message) internal pure returns (bytes32) {
        return keccak256(abi.encodePacked(
            "\x19Ethereum Signed Message:\n32", keccak256(abi.encodePacked(message)))
        );
    }
}

// File: openzeppelin-solidity/contracts/GSN/Context.sol

pragma solidity ^0.5.0;

/*
 * @dev Provides information about the current execution context, including the
 * sender of the transaction and its data. While these are generally available
 * via msg.sender and msg.data, they should not be accessed in such a direct
 * manner, since when dealing with GSN meta-transactions the account sending and
 * paying for execution may not be the actual sender (as far as an application
 * is concerned).
 *
 * This contract is only required for intermediate, library-like contracts.
 */
contract Context {
    // Empty internal constructor, to prevent people from mistakenly deploying
    // an instance of this contract, which should be used via inheritance.
    constructor () internal { }
    // solhint-disable-previous-line no-empty-blocks

    function _msgSender() internal view returns (address payable) {
        return msg.sender;
    }

    function _msgData() internal view returns (bytes memory) {
        this; // silence state mutability warning without generating bytecode - see https://github.com/ethereum/solidity/issues/2691
        return msg.data;
    }
}

// File: openzeppelin-solidity/contracts/token/ERC20/IERC20.sol

pragma solidity ^0.5.0;

/**
 * @dev Interface of the ERC20 standard as defined in the EIP. Does not include
 * the optional functions; to access them see {ERC20Detailed}.
 */
interface IERC20 {
    /**
     * @dev Returns the amount of tokens in existence.
     */
    function totalSupply() external view returns (uint256);

    /**
     * @dev Returns the amount of tokens owned by `account`.
     */
    function balanceOf(address account) external view returns (uint256);

    /**
     * @dev Moves `amount` tokens from the caller's account to `recipient`.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * Emits a {Transfer} event.
     */
    function transfer(address recipient, uint256 amount) external returns (bool);

    /**
     * @dev Returns the remaining number of tokens that `spender` will be
     * allowed to spend on behalf of `owner` through {transferFrom}. This is
     * zero by default.
     *
     * This value changes when {approve} or {transferFrom} are called.
     */
    function allowance(address owner, address spender) external view returns (uint256);

    /**
     * @dev Sets `amount` as the allowance of `spender` over the caller's tokens.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * IMPORTANT: Beware that changing an allowance with this method brings the risk
     * that someone may use both the old and the new allowance by unfortunate
     * transaction ordering. One possible solution to mitigate this race
     * condition is to first reduce the spender's allowance to 0 and set the
     * desired value afterwards:
     * https://github.com/ethereum/EIPs/issues/20#issuecomment-263524729
     *
     * Emits an {Approval} event.
     */
    function approve(address spender, uint256 amount) external returns (bool);

    /**
     * @dev Moves `amount` tokens from `sender` to `recipient` using the
     * allowance mechanism. `amount` is then deducted from the caller's
     * allowance.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * Emits a {Transfer} event.
     */
    function transferFrom(address sender, address recipient, uint256 amount) external returns (bool);

    /**
     * @dev Emitted when `value` tokens are moved from one account (`from`) to
     * another (`to`).
     *
     * Note that `value` may be zero.
     */
    event Transfer(address indexed from, address indexed to, uint256 value);

    /**
     * @dev Emitted when the allowance of a `spender` for an `owner` is set by
     * a call to {approve}. `value` is the new allowance.
     */
    event Approval(address indexed owner, address indexed spender, uint256 value);
}

// File: openzeppelin-solidity/contracts/token/ERC20/ERC20.sol

pragma solidity ^0.5.0;




/**
 * @dev Implementation of the {IERC20} interface.
 *
 * This implementation is agnostic to the way tokens are created. This means
 * that a supply mechanism has to be added in a derived contract using {_mint}.
 * For a generic mechanism see {ERC20Mintable}.
 *
 * TIP: For a detailed writeup see our guide
 * https://forum.zeppelin.solutions/t/how-to-implement-erc20-supply-mechanisms/226[How
 * to implement supply mechanisms].
 *
 * We have followed general OpenZeppelin guidelines: functions revert instead
 * of returning `false` on failure. This behavior is nonetheless conventional
 * and does not conflict with the expectations of ERC20 applications.
 *
 * Additionally, an {Approval} event is emitted on calls to {transferFrom}.
 * This allows applications to reconstruct the allowance for all accounts just
 * by listening to said events. Other implementations of the EIP may not emit
 * these events, as it isn't required by the specification.
 *
 * Finally, the non-standard {decreaseAllowance} and {increaseAllowance}
 * functions have been added to mitigate the well-known issues around setting
 * allowances. See {IERC20-approve}.
 */
contract ERC20 is Context, IERC20 {
    using SafeMath for uint256;

    mapping (address => uint256) private _balances;

    mapping (address => mapping (address => uint256)) private _allowances;

    uint256 private _totalSupply;

    /**
     * @dev See {IERC20-totalSupply}.
     */
    function totalSupply() public view returns (uint256) {
        return _totalSupply;
    }

    /**
     * @dev See {IERC20-balanceOf}.
     */
    function balanceOf(address account) public view returns (uint256) {
        return _balances[account];
    }

    /**
     * @dev See {IERC20-transfer}.
     *
     * Requirements:
     *
     * - `recipient` cannot be the zero address.
     * - the caller must have a balance of at least `amount`.
     */
    function transfer(address recipient, uint256 amount) public returns (bool) {
        _transfer(_msgSender(), recipient, amount);
        return true;
    }

    /**
     * @dev See {IERC20-allowance}.
     */
    function allowance(address owner, address spender) public view returns (uint256) {
        return _allowances[owner][spender];
    }

    /**
     * @dev See {IERC20-approve}.
     *
     * Requirements:
     *
     * - `spender` cannot be the zero address.
     */
    function approve(address spender, uint256 amount) public returns (bool) {
        _approve(_msgSender(), spender, amount);
        return true;
    }

    /**
     * @dev See {IERC20-transferFrom}.
     *
     * Emits an {Approval} event indicating the updated allowance. This is not
     * required by the EIP. See the note at the beginning of {ERC20};
     *
     * Requirements:
     * - `sender` and `recipient` cannot be the zero address.
     * - `sender` must have a balance of at least `amount`.
     * - the caller must have allowance for `sender`'s tokens of at least
     * `amount`.
     */
    function transferFrom(address sender, address recipient, uint256 amount) public returns (bool) {
        _transfer(sender, recipient, amount);
        _approve(sender, _msgSender(), _allowances[sender][_msgSender()].sub(amount, "ERC20: transfer amount exceeds allowance"));
        return true;
    }

    /**
     * @dev Atomically increases the allowance granted to `spender` by the caller.
     *
     * This is an alternative to {approve} that can be used as a mitigation for
     * problems described in {IERC20-approve}.
     *
     * Emits an {Approval} event indicating the updated allowance.
     *
     * Requirements:
     *
     * - `spender` cannot be the zero address.
     */
    function increaseAllowance(address spender, uint256 addedValue) public returns (bool) {
        _approve(_msgSender(), spender, _allowances[_msgSender()][spender].add(addedValue));
        return true;
    }

    /**
     * @dev Atomically decreases the allowance granted to `spender` by the caller.
     *
     * This is an alternative to {approve} that can be used as a mitigation for
     * problems described in {IERC20-approve}.
     *
     * Emits an {Approval} event indicating the updated allowance.
     *
     * Requirements:
     *
     * - `spender` cannot be the zero address.
     * - `spender` must have allowance for the caller of at least
     * `subtractedValue`.
     */
    function decreaseAllowance(address spender, uint256 subtractedValue) public returns (bool) {
        _approve(_msgSender(), spender, _allowances[_msgSender()][spender].sub(subtractedValue, "ERC20: decreased allowance below zero"));
        return true;
    }

    /**
     * @dev Moves tokens `amount` from `sender` to `recipient`.
     *
     * This is internal function is equivalent to {transfer}, and can be used to
     * e.g. implement automatic token fees, slashing mechanisms, etc.
     *
     * Emits a {Transfer} event.
     *
     * Requirements:
     *
     * - `sender` cannot be the zero address.
     * - `recipient` cannot be the zero address.
     * - `sender` must have a balance of at least `amount`.
     */
    function _transfer(address sender, address recipient, uint256 amount) internal {
        require(sender != address(0), "ERC20: transfer from the zero address");
        require(recipient != address(0), "ERC20: transfer to the zero address");

        _balances[sender] = _balances[sender].sub(amount, "ERC20: transfer amount exceeds balance");
        _balances[recipient] = _balances[recipient].add(amount);
        emit Transfer(sender, recipient, amount);
    }

    /** @dev Creates `amount` tokens and assigns them to `account`, increasing
     * the total supply.
     *
     * Emits a {Transfer} event with `from` set to the zero address.
     *
     * Requirements
     *
     * - `to` cannot be the zero address.
     */
    function _mint(address account, uint256 amount) internal {
        require(account != address(0), "ERC20: mint to the zero address");

        _totalSupply = _totalSupply.add(amount);
        _balances[account] = _balances[account].add(amount);
        emit Transfer(address(0), account, amount);
    }

    /**
     * @dev Destroys `amount` tokens from `account`, reducing the
     * total supply.
     *
     * Emits a {Transfer} event with `to` set to the zero address.
     *
     * Requirements
     *
     * - `account` cannot be the zero address.
     * - `account` must have at least `amount` tokens.
     */
    function _burn(address account, uint256 amount) internal {
        require(account != address(0), "ERC20: burn from the zero address");

        _balances[account] = _balances[account].sub(amount, "ERC20: burn amount exceeds balance");
        _totalSupply = _totalSupply.sub(amount);
        emit Transfer(account, address(0), amount);
    }

    /**
     * @dev Sets `amount` as the allowance of `spender` over the `owner`s tokens.
     *
     * This is internal function is equivalent to `approve`, and can be used to
     * e.g. set automatic allowances for certain subsystems, etc.
     *
     * Emits an {Approval} event.
     *
     * Requirements:
     *
     * - `owner` cannot be the zero address.
     * - `spender` cannot be the zero address.
     */
    function _approve(address owner, address spender, uint256 amount) internal {
        require(owner != address(0), "ERC20: approve from the zero address");
        require(spender != address(0), "ERC20: approve to the zero address");

        _allowances[owner][spender] = amount;
        emit Approval(owner, spender, amount);
    }

    /**
     * @dev Destroys `amount` tokens from `account`.`amount` is then deducted
     * from the caller's allowance.
     *
     * See {_burn} and {_approve}.
     */
    function _burnFrom(address account, uint256 amount) internal {
        _burn(account, amount);
        _approve(account, _msgSender(), _allowances[account][_msgSender()].sub(amount, "ERC20: burn amount exceeds allowance"));
    }
}

// File: openzeppelin-solidity/contracts/access/Roles.sol

pragma solidity ^0.5.0;

/**
 * @title Roles
 * @dev Library for managing addresses assigned to a Role.
 */
library Roles {
    struct Role {
        mapping (address => bool) bearer;
    }

    /**
     * @dev Give an account access to this role.
     */
    function add(Role storage role, address account) internal {
        require(!has(role, account), "Roles: account already has role");
        role.bearer[account] = true;
    }

    /**
     * @dev Remove an account's access to this role.
     */
    function remove(Role storage role, address account) internal {
        require(has(role, account), "Roles: account does not have role");
        role.bearer[account] = false;
    }

    /**
     * @dev Check if an account has this role.
     * @return bool
     */
    function has(Role storage role, address account) internal view returns (bool) {
        require(account != address(0), "Roles: account is the zero address");
        return role.bearer[account];
    }
}

// File: openzeppelin-solidity/contracts/access/roles/MinterRole.sol

pragma solidity ^0.5.0;



contract MinterRole is Context {
    using Roles for Roles.Role;

    event MinterAdded(address indexed account);
    event MinterRemoved(address indexed account);

    Roles.Role private _minters;

    constructor () internal {
        _addMinter(_msgSender());
    }

    modifier onlyMinter() {
        require(isMinter(_msgSender()), "MinterRole: caller does not have the Minter role");
        _;
    }

    function isMinter(address account) public view returns (bool) {
        return _minters.has(account);
    }

    function addMinter(address account) public onlyMinter {
        _addMinter(account);
    }

    function renounceMinter() public {
        _removeMinter(_msgSender());
    }

    function _addMinter(address account) internal {
        _minters.add(account);
        emit MinterAdded(account);
    }

    function _removeMinter(address account) internal {
        _minters.remove(account);
        emit MinterRemoved(account);
    }
}

// File: openzeppelin-solidity/contracts/token/ERC20/ERC20Mintable.sol

pragma solidity ^0.5.0;



/**
 * @dev Extension of {ERC20} that adds a set of accounts with the {MinterRole},
 * which have permission to mint (create) new tokens as they see fit.
 *
 * At construction, the deployer of the contract is the only minter.
 */
contract ERC20Mintable is ERC20, MinterRole {
    /**
     * @dev See {ERC20-_mint}.
     *
     * Requirements:
     *
     * - the caller must have the {MinterRole}.
     */
    function mint(address account, uint256 amount) public onlyMinter returns (bool) {
        _mint(account, amount);
        return true;
    }
}

// File: openzeppelin-solidity/contracts/token/ERC20/ERC20Detailed.sol

pragma solidity ^0.5.0;


/**
 * @dev Optional functions from the ERC20 standard.
 */
contract ERC20Detailed is IERC20 {
    string private _name;
    string private _symbol;
    uint8 private _decimals;

    /**
     * @dev Sets the values for `name`, `symbol`, and `decimals`. All three of
     * these values are immutable: they can only be set once during
     * construction.
     */
    constructor (string memory name, string memory symbol, uint8 decimals) public {
        _name = name;
        _symbol = symbol;
        _decimals = decimals;
    }

    /**
     * @dev Returns the name of the token.
     */
    function name() public view returns (string memory) {
        return _name;
    }

    /**
     * @dev Returns the symbol of the token, usually a shorter version of the
     * name.
     */
    function symbol() public view returns (string memory) {
        return _symbol;
    }

    /**
     * @dev Returns the number of decimals used to get its user representation.
     * For example, if `decimals` equals `2`, a balance of `505` tokens should
     * be displayed to a user as `5,05` (`505 / 10 ** 2`).
     *
     * Tokens usually opt for a value of 18, imitating the relationship between
     * Ether and Wei.
     *
     * NOTE: This information is only used for _display_ purposes: it in
     * no way affects any of the arithmetic of the contract, including
     * {IERC20-balanceOf} and {IERC20-transfer}.
     */
    function decimals() public view returns (uint8) {
        return _decimals;
    }
}

// File: contracts/BridgeBank/BridgeToken.sol

pragma solidity ^0.5.0;



/**
 * @title BridgeToken
 * @dev Mintable, ERC20 compatible BankToken for use by BridgeBank
 **/

contract BridgeToken is ERC20Mintable, ERC20Detailed {

    constructor(
        string memory _symbol
    )
        public
        ERC20Detailed(
            _symbol,
            _symbol,
            18
        )
    {
        // Intentionally left blank
    }
}

// File: contracts/BridgeBank/CosmosBank.sol

pragma solidity ^0.5.0;



/**
 * @title CosmosBank
 * @dev Manages the deployment and minting of ERC20 compatible BridgeTokens
 *      which represent assets based on the Cosmos blockchain.
 **/

contract CosmosBank {

    using SafeMath for uint256;

    uint256 public bridgeTokenCount;
    mapping(address => bool) public bridgeTokenWhitelist;
    mapping(bytes32 => CosmosDeposit) cosmosDeposits;

    struct CosmosDeposit {
        bytes cosmosSender;
        address payable ethereumRecipient;
        address bridgeTokenAddress;
        uint256 amount;
        bool locked;
    }

    /*
    * @dev: Event declarations
    */
    event LogNewBridgeToken(
        address _token,
        string _symbol
    );

    event LogBridgeTokenMint(
        address _token,
        string _symbol,
        uint256 _amount,
        address _beneficiary
    );

    /*
    * @dev: Constructor, sets bridgeTokenCount
    */
    constructor () public {
        bridgeTokenCount = 0;
    }

    /*
    * @dev: Creates a new CosmosDeposit with a unique ID
    *
    * @param _cosmosSender: The sender's Cosmos address in bytes.
    * @param _ethereumRecipient: The intended recipient's Ethereum address.
    * @param _token: The currency type
    * @param _amount: The amount in the deposit.
    * @return: The newly created CosmosDeposit's unique id.
    */
    function newCosmosDeposit(
        bytes memory _cosmosSender,
        address payable _ethereumRecipient,
        address _token,
        uint256 _amount
    )
        internal
        returns(bytes32)
    {
        bytes32 depositID = keccak256(
            abi.encodePacked(
                _cosmosSender,
                _ethereumRecipient,
                _token,
                _amount
            )
        );

        cosmosDeposits[depositID] = CosmosDeposit(
            _cosmosSender,
            _ethereumRecipient,
            _token,
            _amount,
            true
        );

        return depositID;
    }

    /*
     * @dev: Deploys a new BridgeToken contract
     *
     * @param _symbol: The BridgeToken's symbol
     */
    function deployNewBridgeToken(
        string memory _symbol
    )
        internal
        returns(address)
    {
        // return address(0);
        bridgeTokenCount = bridgeTokenCount.add(1);

        // Deploy new bridge token contract
        BridgeToken newBridgeToken = (new BridgeToken)(_symbol);

        // Set address in tokens mapping
        address newBridgeTokenAddress = address(newBridgeToken);
        bridgeTokenWhitelist[newBridgeTokenAddress] = true;

        emit LogNewBridgeToken(
            newBridgeTokenAddress,
            _symbol
        );

        return newBridgeTokenAddress;
    }

    /*
     * @dev: Mints new cosmos tokens
     *
     * @param _cosmosSender: The sender's Cosmos address in bytes.
     * @param _ethereumRecipient: The intended recipient's Ethereum address.
     * @param _cosmosTokenAddress: The currency type
     * @param _symbol: comsos token symbol
     * @param _amount: number of comsos tokens to be minted
\    */
     function mintNewBridgeTokens(
        bytes memory _cosmosSender,
        address payable _intendedRecipient,
        address _bridgeTokenAddress,
        string memory _symbol,
        uint256 _amount
    )
        internal
    {
        // Must be whitelisted bridge token
        require(
            bridgeTokenWhitelist[_bridgeTokenAddress],
            "Token must be a whitelisted bridge token"
        );

        // Mint bridge tokens
        require(
            BridgeToken(_bridgeTokenAddress).mint(
                _intendedRecipient,
                _amount
            ),
            "Attempted mint of bridge tokens failed"
        );

        newCosmosDeposit(
            _cosmosSender,
            _intendedRecipient,
            _bridgeTokenAddress,
            _amount
        );

        emit LogBridgeTokenMint(
            _bridgeTokenAddress,
            _symbol,
            _amount,
            _intendedRecipient
        );
    }

    /*
    * @dev: Checks if an individual CosmosDeposit exists.
    *
    * @param _id: The unique CosmosDeposit's id.
    * @return: Boolean indicating if the CosmosDeposit exists in memory.
    */
    function isLockedCosmosDeposit(
        bytes32 _id
    )
        internal
        view
        returns(bool)
    {
        return(cosmosDeposits[_id].locked);
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
    function getCosmosDeposit(
        bytes32 _id
    )
        internal
        view
        returns(bytes memory, address payable, address, uint256)
    {
        CosmosDeposit memory deposit = cosmosDeposits[_id];

        return(
            deposit.cosmosSender,
            deposit.ethereumRecipient,
            deposit.bridgeTokenAddress,
            deposit.amount
        );
    }
}

// File: contracts/BridgeBank/EthereumBank.sol

pragma solidity ^0.5.0;



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

// File: contracts/Oracle.sol

pragma solidity ^0.5.0;

// import "../../node_modules/openzeppelin-solidity/contracts/math/SafeMath.sol";
// import "./Valset.sol";
// import "./CosmosBridge.sol";

contract Oracle {

    using SafeMath for uint256;

    /*
    * @dev: Public variable declarations
    */
    CosmosBridge public cosmosBridge;
    Valset public valset;
    address public operator;

    // Tracks the number of OracleClaims made on an individual BridgeClaim
    mapping(uint256 => address[]) public oracleClaimValidators;
    mapping(uint256 => mapping(address => bool)) public hasMadeClaim;

    /*
    * @dev: Event declarations
    */
    event LogNewOracleClaim(
        uint256 _prophecyID,
        string _message,
        address _validatorAddress,
        bytes _signature
    );

    event LogProphecyProcessed(
        uint256 _prophecyID,
        uint256 _weightedSignedPower,
        uint256 _weightedTotalPower,
        address _submitter
    );

    /*
    * @dev: Modifier to restrict access to the operator.
    */
    modifier onlyOperator()
    {
        require(
            msg.sender == operator,
            'Must be the operator.'
        );
        _;
    }

    /*
    * @dev: Modifier to restrict access to current ValSet validators
    */
    modifier onlyValidator()
    {
        require(
            valset.isActiveValidator(msg.sender),
            "Must be an active validator"
        );
        _;
    }

    /*
    * @dev: Modifier to restrict access to current ValSet validators
    */
    modifier isPending(
        uint256 _prophecyID
    )
    {
        require(
            cosmosBridge.isProphecyClaimActive(
                _prophecyID
            ) == true,
            "The prophecy must be pending for this operation"
        );
        _;
    }

    /*
    * @dev: Constructor
    */
    constructor(
        address _operator,
        address _valset,
        address _cosmosBridge
    )
        public
    {
        operator = _operator;
        cosmosBridge = CosmosBridge(_cosmosBridge);
        valset = Valset(_valset);
    }

    /*
    * @dev: newOracleClaim
    *       Allows validators to make new OracleClaims on an existing Prophecy
    */
    function newOracleClaim(
        uint256 _prophecyID,
        string memory _message,
        bytes memory _signature
    )
        public
        onlyValidator
        isPending(_prophecyID)
    {
        address validatorAddress = msg.sender;

        // Validate the msg.sender's signature
        require(
            validatorAddress == valset.recover(
                _message,
                _signature
            ),
            "Invalid message signature."
        );

        // Confirm that this address has not already made an oracle claim on this prophecy
        require(
            !hasMadeClaim[_prophecyID][validatorAddress],
            "Cannot make duplicate oracle claims from the same address."
        );

        hasMadeClaim[_prophecyID][validatorAddress] = true;
        oracleClaimValidators[_prophecyID].push(validatorAddress);

        emit LogNewOracleClaim(
            _prophecyID,
            _message,
            validatorAddress,
            _signature
        );

        // Process the prophecy
        (bool valid,
            uint256 weightedSignedPower,
            uint256 weightedTotalPower
        ) = getProphecyThreshold(_prophecyID);

        if (valid) {
            completeProphecy(
                _prophecyID
            );

            emit LogProphecyProcessed(
                _prophecyID,
                weightedSignedPower,
                weightedTotalPower,
                msg.sender
            );
        }

    }

    /*
    * @dev: processBridgeProphecy
    *       Pubically available method which attempts to process a bridge prophecy
    */
    function processBridgeProphecy(
        uint256 _prophecyID
    )
        public
        isPending(_prophecyID)
    {
        // Process the prophecy
        (bool valid,
            uint256 weightedSignedPower,
            uint256 weightedTotalPower
        ) = getProphecyThreshold(_prophecyID);

        require(
            valid,
            "The cumulative power of signatory validators does not meet the threshold"
        );

        // Update the BridgeClaim's status
        completeProphecy(
            _prophecyID
        );

        emit LogProphecyProcessed(
            _prophecyID,
            weightedSignedPower,
            weightedTotalPower,
            msg.sender
        );
    }

    /*
    * @dev: checkBridgeProphecy
    *       Operator accessor method which checks if a prophecy has passed
    *       the validity threshold, without actually completing the prophecy.
    */
    function checkBridgeProphecy(
        uint256 _prophecyID
    )
        public
        view
        onlyOperator
        isPending(_prophecyID)
        returns(bool, uint256, uint256)
    {
        require(
            cosmosBridge.isProphecyClaimActive(
                _prophecyID
            ) == true,
            "Can only check active prophecies"
        );
        return getProphecyThreshold(
            _prophecyID
        );
    }

    /*
    * @dev: processProphecy
    *       Calculates the status of a prophecy. The claim is considered valid if the
    *       combined active signatory validator powers pass the validation threshold.
    *       The hardcoded threshold is (Combined signed power * 2) >= (Total power * 3).
    */
    function getProphecyThreshold(
        uint256 _prophecyID
    )
        internal
        view
        returns(bool, uint256, uint256)
    {
        uint256 signedPower = 0;
        uint256 totalPower = valset.totalPower();

        // Iterate over the signatory addresses
        for (uint256 i = 0; i < oracleClaimValidators[_prophecyID].length; i = i.add(1)) {
            address signer = oracleClaimValidators[_prophecyID][i];

                // Only add the power of active validators
                if(valset.isActiveValidator(signer)) {
                    signedPower = signedPower.add(
                        valset.getValidatorPower(
                            signer
                        )
                    );
                }
        }

        // Calculate if weighted signed power has reached threshold of weighted total power
        uint256 weightedSignedPower = signedPower.mul(3);
        uint256 weightedTotalPower = totalPower.mul(2);
        bool hasReachedThreshold = weightedSignedPower >= weightedTotalPower;

        return(
            hasReachedThreshold,
            weightedSignedPower,
            weightedTotalPower
        );
    }

    /*
    * @dev: completeProphecy
    *       Completes a prophecy by completing the corresponding BridgeClaim
    *       on the CosmosBridge.
    */
    function completeProphecy(
        uint256 _prophecyID
    )
        internal
    {
        cosmosBridge.completeProphecyClaim(
            _prophecyID
        );
    }
}

// File: contracts/BridgeBank/BridgeBank.sol

pragma solidity ^0.5.0;



// import "../CosmosBridge.sol";

/**
 * @title BridgeBank
 * @dev Bank contract which coordinates asset-related functionality.
 *      CosmosBank manages the minting and burning of tokens which
 *      represent Cosmos based assets, while EthereumBank manages
 *      the locking and unlocking of Ethereum and ERC20 token assets
 *      based on Ethereum.
 **/

contract BridgeBank is CosmosBank, EthereumBank {

    using SafeMath for uint256;
    
    address public operator;
    Oracle public oracle;
    CosmosBridge public cosmosBridge;

    /*
    * @dev: Constructor, sets operator
    */
    constructor (
        address _operatorAddress,
        address _oracleAddress,
        address _cosmosBridgeAddress
    )
        public
    {
        operator = _operatorAddress;
        oracle = Oracle(_oracleAddress);
        cosmosBridge = CosmosBridge(_cosmosBridgeAddress);
    }

    /*
    * @dev: Modifier to restrict access to operator
    */
    modifier onlyOperator() {
        require(
            msg.sender == operator,
            'Must be BridgeBank operator.'
        );
        _;
    }

    /*
    * @dev: Modifier to restrict access to the oracle
    */
    modifier onlyOracle()
    {
        require(
            msg.sender == address(oracle),
            "Access restricted to the oracle"
        );
        _;
    }

    /*
    * @dev: Modifier to restrict access to the cosmos bridge
    */
    modifier onlyCosmosBridge()
    {
        require(
            msg.sender == address(cosmosBridge) || msg.sender == operator, // TODO: remove this after EthDenver
            "Access restricted to the cosmos bridge"
        );
        _;
    }

   /*
    * @dev: Fallback function allows operator to send funds to the bank directly
    *       This feature is used for testing and is available at the operator's own risk.
    */
    function() external payable onlyOperator {}

    /*
    * @dev: Creates a new BridgeToken
    *
    * @param _symbol: The new BridgeToken's symbol
    * @return: The new BridgeToken contract's address
    */
    function createNewBridgeToken(
        string memory _symbol
    )
        public
        onlyOperator
        returns(address)
    {
        return deployNewBridgeToken(_symbol);
    }

    /*
     * @dev: Mints new BankTokens
     *
     * @param _cosmosSender: The sender's Cosmos address in bytes.
     * @param _ethereumRecipient: The intended recipient's Ethereum address.
     * @param _cosmosTokenAddress: The currency type
     * @param _symbol: comsos token symbol
     * @param _amount: number of comsos tokens to be minted
     */
     function mintBridgeTokens(
        bytes memory _cosmosSender,
        address payable _intendedRecipient,
        address _bridgeTokenAddress,
        string memory _symbol,
        uint256 _amount
    )
        public
        onlyCosmosBridge
    {
        return mintNewBridgeTokens(
            _cosmosSender,
            _intendedRecipient,
            _bridgeTokenAddress,
            _symbol,
            _amount
        );
    }

    /*
    * @dev: Locks received Ethereum funds.
    *
    * @param _recipient: bytes representation of destination address.
    * @param _token: token address in origin chain (0x0 if ethereum)
    * @param _amount: value of deposit
    */
    function lock(
        bytes memory _recipient,
        address _token,
        uint256 _amount
    )
        public
        availableNonce()
        payable
    {
        string memory symbol;

        // Ethereum deposit
        if (msg.value > 0) {
          require(
              _token == address(0),
              "Ethereum deposits require the 'token' address to be the null address"
            );
          require(
              msg.value == _amount,
              "The transactions value must be equal the specified amount (in wei)"
            );

          // Set the the symbol to ETH
          symbol = "ETH";
          // ERC20 deposit
        } else {
          require(
              BridgeToken(_token).transferFrom(msg.sender, address(this), _amount),
              "Contract token allowances insufficient to complete this lock request"
          );
          // Set symbol to the ERC20 token's symbol
          symbol = BridgeToken(_token).symbol();
        }

        lockFunds(
            msg.sender,
            _recipient,
            _token,
            symbol,
            _amount
        );
    }

   /*
    * @dev: Unlocks Ethereum and ERC20 tokens held on the contract.
    *
    * @param _recipient: recipient's Ethereum address
    * @param _token: token contract address
    * @param _symbol: token symbol
    * @param _amount: wei amount or ERC20 token count
\   */
     function unlock(
        address payable _recipient,
        address _token,
        string memory _symbol,
        uint256 _amount
    )
        public
        onlyCosmosBridge
        hasLockedFunds(
            _token,
            _amount
        )
        canDeliver(
            _token,
            _amount
        )
    {
        unlockFunds(
            _recipient,
            _token,
            _symbol,
            _amount
        );
    }

    /*
    * @dev: Exposes an item's current status.
    *
    * @param _id: The item in question.
    * @return: Boolean indicating the lock status.
    */
    function getCosmosDepositStatus(
        bytes32 _id
    )
        public
        view
        returns(bool)
    {
        return isLockedCosmosDeposit(_id);
    }

    /*
    * @dev: Allows access to a Cosmos deposit's information via its unique identifier.
    *
    * @param _id: The deposit to be viewed.
    * @return: Original sender's Ethereum address.
    * @return: Intended Cosmos recipient's address in bytes.
    * @return: The lock deposit's currency, denoted by a token address.
    * @return: The amount locked in the deposit.
    * @return: The deposit's unique nonce.
    */
    function viewCosmosDeposit(
        bytes32 _id
    )
        public
        view
        returns(bytes memory, address payable, address, uint256)
    {
        return getCosmosDeposit(_id);
    }

}

// File: contracts/CosmosBridge.sol

pragma solidity ^0.5.0;




contract CosmosBridge {

    using SafeMath for uint256;

    /*
    * @dev: Public variable declarations
    */
    address public operator;
    Valset public valset;
    address public oracle;
    bool public hasOracle;
    BridgeBank public bridgeBank;
    bool public hasBridgeBank;

    uint256 public prophecyClaimCount;
    mapping(uint256 => ProphecyClaim) public prophecyClaims;

    enum Status {
        Null,
        Pending,
        Success,
        Failed
    }

    enum ClaimType {
        Unsupported,
        Burn,
        Lock
    }

    struct ProphecyClaim {
        ClaimType claimType;
        bytes cosmosSender;
        address payable ethereumReceiver;
        address originalValidator;
        address tokenAddress;
        string symbol;
        uint256 amount;
        Status status;
    }

    /*
    * @dev: Event declarations
    */
    event LogOracleSet(
        address _oracle
    );

    event LogBridgeBankSet(
        address _bridgeBank
    );

    event LogNewProphecyClaim(
        uint256 _prophecyID,
        ClaimType _claimType,
        bytes _cosmosSender,
        address payable _ethereumReceiver,
        address _validatorAddress,
        address _tokenAddress,
        string _symbol,
        uint256 _amount
    );

    event LogProphecyCompleted(
        uint256 _prophecyID,
        ClaimType _claimType
    );

    /*
    * @dev: Modifier which only allows access to currently pending prophecies
    */
    modifier isPending(
        uint256 _prophecyID
    )
    {
        require(
            isProphecyClaimActive(_prophecyID),
            "Prophecy claim is not active"
        );
        _;
    }

    /*
    * @dev: Modifier to restrict access to the operator.
    */
    modifier onlyOperator()
    {
        require(
            msg.sender == operator,
            'Must be the operator.'
        );
        _;
    }

      /*
    * @dev: The bridge is not active until oracle and bridge bank are set
    */
    modifier isActive()
    {
        require(
            hasOracle == true && hasBridgeBank == true,
            "The Operator must set the oracle and bridge bank for bridge activation"
        );
        _;
    }

    /*
    * @dev: Constructor
    */
    constructor(
        address _operator,
        address _valset
    )
        public
    {
        prophecyClaimCount = 0;
        operator = _operator;
        valset = Valset(_valset);
        hasOracle = false;
        hasBridgeBank = false;
    }

    /*
    * @dev: setOracle
    */
    function setOracle(
        address _oracle
    )
        public
        onlyOperator
    {
        require(
            !hasOracle,
            "The Oracle cannot be updated once it has been set"
        );

        hasOracle = true;
        oracle = _oracle;

        emit LogOracleSet(
            oracle
        );
    }

    /*
    * @dev: setBridgeBank
    */
    function setBridgeBank(
        address payable _bridgeBank
    )
        public
        onlyOperator
    {
        require(
            !hasBridgeBank,
            "The Bridge Bank cannot be updated once it has been set"
        );

        hasBridgeBank = true;
        bridgeBank = BridgeBank(_bridgeBank);

        emit LogBridgeBankSet(
            address(bridgeBank)
        );
    }

    /*
    * @dev: newProphecyClaim
    *       Creates a new burn or lock prophecy claim, adding it to the prophecyClaims mapping.
    *       Lock claims can only be created for BridgeTokens on BridgeBank's whitelist. The operator
    *       is responsible for adding them, and lock claims will fail until the operator has done so.
    */
    function newProphecyClaim(
        ClaimType _claimType,
        bytes memory _cosmosSender,
        address payable _ethereumReceiver,
        address _tokenAddress,
        string memory _symbol,
        uint256 _amount
    )
        public
        isActive
    {
        require(
            valset.isActiveValidator(msg.sender),
            "Must be an active validator"
        );

        // Increment the prophecy claim count
        prophecyClaimCount = prophecyClaimCount.add(1);

        address originalValidator = msg.sender;

        ClaimType claimType;
        if(_claimType == ClaimType.Burn){
            claimType = ClaimType.Burn;
        } else if(_claimType == ClaimType.Lock){
            claimType = ClaimType.Lock;
        }

        // Create the new ProphecyClaim
        ProphecyClaim memory prophecyClaim = ProphecyClaim(
            claimType,
            _cosmosSender,
            _ethereumReceiver,
            originalValidator,
            _tokenAddress,
            _symbol,
            _amount,
            Status.Pending
        );

        // Add the new ProphecyClaim to the mapping
        prophecyClaims[prophecyClaimCount] = prophecyClaim;

        emit LogNewProphecyClaim(
            prophecyClaimCount,
            claimType,
            _cosmosSender,
            _ethereumReceiver,
            originalValidator,
            _tokenAddress,
            _symbol,
            _amount
        );
    }

    /*
    * @dev: completeProphecyClaim
    *       Allows for the completion of ProphecyClaims once processed by the Oracle.
    *       Burn claims unlock tokens stored by BridgeBank.
    *       Lock claims mint BridgeTokens on BridgeBank's token whitelist.
    */
    function completeProphecyClaim(
        uint256 _prophecyID
    )
        public
        isPending(_prophecyID)
    {
        require(
            msg.sender == oracle,
            "Only the Oracle may complete prophecies"
        );

        prophecyClaims[_prophecyID].status = Status.Success;

        ClaimType claimType = prophecyClaims[_prophecyID].claimType;
        if(claimType == ClaimType.Burn) {
            unlockTokens(_prophecyID);
        } else {
            issueBridgeTokens(_prophecyID);
        }

        emit LogProphecyCompleted(
            _prophecyID,
            claimType
        );
    }

    /*
    * @dev: issueBridgeTokens
    *       Issues a request for the BridgeBank to mint new BridgeTokens
    */
    function issueBridgeTokens(
        uint256 _prophecyID
    )
        internal
    {
        ProphecyClaim memory prophecyClaim = prophecyClaims[_prophecyID];

        bridgeBank.mintBridgeTokens(
            prophecyClaim.cosmosSender,
            prophecyClaim.ethereumReceiver,
            prophecyClaim.tokenAddress,
            prophecyClaim.symbol,
            prophecyClaim.amount
        );
    }

    /*
    * @dev: unlockTokens
    *       Issues a request for the BridgeBank to unlock funds held on contract
    */
    function unlockTokens(
        uint256 _prophecyID
    )
        internal
    {
        ProphecyClaim memory prophecyClaim = prophecyClaims[_prophecyID];

        bridgeBank.unlock(
            prophecyClaim.ethereumReceiver,
            prophecyClaim.tokenAddress,
            prophecyClaim.symbol,
            prophecyClaim.amount
        );
    }

    /*
    * @dev: isProphecyClaimActive
    *       Returns boolean indicating if the ProphecyClaim is active
    */
    function isProphecyClaimActive(
        uint256 _prophecyID
    )
        public
        view
        returns(bool)
    {
        return prophecyClaims[_prophecyID].status == Status.Pending;
    }

    /*
    * @dev: isProphecyValidatorActive
    *       Returns boolean indicating if the validator that originally
    *       submitted the ProphecyClaim is still an active validator
    */
    function isProphecyClaimValidatorActive(
        uint256 _prophecyID
    )
        public
        view
        returns(bool)
    {
        return valset.isActiveValidator(
            prophecyClaims[_prophecyID].originalValidator
        );
    }
}
