pragma solidity ^0.5.0;

import "../../node_modules/openzeppelin-solidity/contracts/math/SafeMath.sol";
import "./BankToken.sol";

/**
 * @title Bank
 * @dev Manages the deployment and minting of ERC20 compatible BankTokens
 **/

contract Bank {

    using SafeMath for uint256;

    mapping(address => bool) public bankTokens;
    uint256 public numbTokens;

    event LogBankTokenDeploy(address _token);

    event LogBankTokenMint(
        address _token,
        string _symbol,
        uint256 _amount,
        address _beneficiary
    );

    constructor () public {
        numbTokens = 0;
    }

    /*
     * @dev: Delivers bank tokens
     *
     * @param _token: bank token contract address, address(0) indicates new token
     * @param _symbol: bank token symbol
     * @param _amount: number of bank tokens to be delivered
     * @param _beneficiary: recipient of the minted tokens
     */
     function deliver(
        address _token,
        string memory _symbol,
        uint256 _amount,
        address _beneficiary
    )
        internal
    {
        // If no token address, deploy a new bank token
        address bankToken = _token;
        if(address(_token) == address(0)) {
            bankToken = deployBankToken(_symbol);
        } else {
            bankToken = _token;
        }

        // Must be bank token controlled by the bank
        require(
            bankTokens[bankToken],
            "Invalid bank token address"
        );

        // Mint bank tokens
        require(
            ERC20Mintable(bankToken).mint(_beneficiary, _amount),
            "Failed to mint bank token"
        );

        emit LogBankTokenMint(bankToken, _symbol, _amount, _beneficiary);
    }

    /*
     * @dev: Deploys a new bank token contract
     *
     * @param _symbol: bank token symbol
     */
    function deployBankToken(
        string memory _symbol
    )
        internal
        returns(address)
    {
        numbTokens = numbTokens.add(1);

        // TODO: BankToken contract deployment puts Peggy over gas limit, causing deployment to fail
        // Deploy new token contract
        // BankToken newToken = (new BankToken)(_symbol);

        // Set address in tokens mapping
        // address newTokenAddress = address(newToken);
        address newTokenAddress = address(0);
        bankTokens[newTokenAddress] = true;

        emit LogBankTokenDeploy(newTokenAddress);

        return newTokenAddress;
    }
}