// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Oracle

import (
  "math/big"
  "strings"

  ethereum "github.com/ethereum/go-ethereum"
  "github.com/ethereum/go-ethereum/accounts/abi"
  "github.com/ethereum/go-ethereum/accounts/abi/bind"
  "github.com/ethereum/go-ethereum/common"
  "github.com/ethereum/go-ethereum/core/types"
  "github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
  _ = big.NewInt
  _ = strings.NewReader
  _ = ethereum.NotFound
  _ = abi.U256
  _ = bind.Bind
  _ = common.Big1
  _ = types.BloomLookup
  _ = event.NewSubscription
)

// OracleABI is the input ABI used to generate the binding from.
const OracleABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"oracleClaimValidators\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"operator\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"valset\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_prophecyID\",\"type\":\"uint256\"}],\"name\":\"processBridgeProphecy\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"hasMadeClaim\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"cosmosBridge\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_prophecyID\",\"type\":\"uint256\"},{\"name\":\"_message\",\"type\":\"bytes32\"},{\"name\":\"_v\",\"type\":\"uint8\"},{\"name\":\"_r\",\"type\":\"bytes32\"},{\"name\":\"_s\",\"type\":\"bytes32\"}],\"name\":\"newOracleClaim\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_prophecyID\",\"type\":\"uint256\"}],\"name\":\"checkBridgeProphecy\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_operator\",\"type\":\"address\"},{\"name\":\"_valset\",\"type\":\"address\"},{\"name\":\"_cosmosBridge\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_prophecyID\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"_validatorAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_message\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_v\",\"type\":\"uint8\"},{\"indexed\":false,\"name\":\"_r\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_s\",\"type\":\"bytes32\"}],\"name\":\"LogNewOracleClaim\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_prophecyID\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"_weightedSignedPower\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"_weightedTotalPower\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"_submitter\",\"type\":\"address\"}],\"name\":\"LogProphecyProcessed\",\"type\":\"event\"}]"

// OracleBin is the compiled bytecode used for deploying new contracts.
const OracleBin = `608060405234801561001057600080fd5b5060405160608061169f8339810180604052606081101561003057600080fd5b8101908080519060200190929190805190602001909291908051906020019092919050505082600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555081600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550505050611576806101296000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c8063a219763e1161005b578063a219763e146101c7578063b0e9ef711461022d578063c484a8a814610277578063e33a8b2a146102d057610088565b806336e413411461008d578063570ca735146101055780637f54af0c1461014f57806389ed70b714610199575b600080fd5b6100c3600480360360408110156100a357600080fd5b810190808035906020019092919080359060200190929190505050610324565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b61010d61036f565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b610157610395565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6101c5600480360360208110156101af57600080fd5b81019080803590602001909291905050506103bb565b005b610213600480360360408110156101dd57600080fd5b8101908080359060200190929190803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506105ba565b604051808215151515815260200191505060405180910390f35b6102356105e9565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6102ce600480360360a081101561028d57600080fd5b810190808035906020019092919080359060200190929190803560ff169060200190929190803590602001909291908035906020019092919050505061060e565b005b6102fc600480360360208110156102e657600080fd5b8101908080359060200190929190505050610c72565b6040518084151515158152602001838152602001828152602001935050505060405180910390f35b6003602052816000526040600020818154811061033d57fe5b906000526020600020016000915091509054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b80600115156000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663d8da69ea836040518263ffffffff1660e01b81526004018082815260200191505060206040518083038186803b15801561043257600080fd5b505afa158015610446573d6000803e3d6000fd5b505050506040513d602081101561045c57600080fd5b81019080805190602001909291905050501515146104c5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602f8152602001806114c1602f913960400191505060405180910390fd5b60008060006104d385610f81565b9250925092508261052f576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260488152602001806114796048913960600191505060405180910390fd5b610538856112db565b7f1d8e3fbd601d9d92db7022fb97f75e132841b94db732dcecb0c93cb31852fcbc85838333604051808581526020018481526020018381526020018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200194505050505060405180910390a15050505050565b60046020528160005260406000206020528060005260406000206000915091509054906101000a900460ff1681565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166340550a1c336040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b1580156106ad57600080fd5b505afa1580156106c1573d6000803e3d6000fd5b505050506040513d60208110156106d757600080fd5b810190808051906020019092919050505061075a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f4d75737420626520616e206163746976652076616c696461746f72000000000081525060200191505060405180910390fd5b84600115156000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663d8da69ea836040518263ffffffff1660e01b81526004018082815260200191505060206040518083038186803b1580156107d157600080fd5b505afa1580156107e5573d6000803e3d6000fd5b505050506040513d60208110156107fb57600080fd5b8101908080519060200190929190505050151514610864576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602f8152602001806114c1602f913960400191505060405180910390fd5b600033905060606040518060400160405280601c81526020017f19457468657265756d205369676e6564204d6573736167653a0a3332000000008152509050600081886040516020018080602001838152602001828103825284818151815260200191508051906020019080838360005b838110156108f05780820151818401526020810190506108d5565b50505050905090810190601f16801561091d5780820380516001836020036101000a031916815260200191505b50935050505060405160208183030381529060405280519060200120905060018188888860405160008152602001604052604051808581526020018460ff1660ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa158015610998573d6000803e3d6000fd5b5050506020604051035173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1614610a42576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601a8152602001807f496e76616c6964206d657373616765207369676e61747572652e00000000000081525060200191505060405180910390fd5b600460008a815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff1615610af6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252603a8152602001806114f0603a913960400191505060405180910390fd5b6001600460008b815260200190815260200160002060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff021916908315150217905550600360008a81526020019081526020016000208390806001815401808255809150509060018203906000526020600020016000909192909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550507f0e596df721fcd7ca12494ee97d05d86caaafbf3d92a47f3557947fbd0e4aba2789848a8a8a8a604051808781526020018673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018581526020018460ff1660ff168152602001838152602001828152602001965050505050505060405180910390a1505050505050505050565b6000806000600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610d3a576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260158152602001807f4d75737420626520746865206f70657261746f722e000000000000000000000081525060200191505060405180910390fd5b83600115156000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663d8da69ea836040518263ffffffff1660e01b81526004018082815260200191505060206040518083038186803b158015610db157600080fd5b505afa158015610dc5573d6000803e3d6000fd5b505050506040513d6020811015610ddb57600080fd5b8101908080519060200190929190505050151514610e44576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602f8152602001806114c1602f913960400191505060405180910390fd5b600115156000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663d8da69ea876040518263ffffffff1660e01b81526004018082815260200191505060206040518083038186803b158015610eba57600080fd5b505afa158015610ece573d6000803e3d6000fd5b505050506040513d6020811015610ee457600080fd5b8101908080519060200190929190505050151514610f6a576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f43616e206f6e6c7920636865636b206163746976652070726f7068656369657381525060200191505060405180910390fd5b610f7385610f81565b935093509350509193909250565b600080600080600090506000600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663db3ad22c6040518163ffffffff1660e01b815260040160206040518083038186803b158015610ff557600080fd5b505afa158015611009573d6000803e3d6000fd5b505050506040513d602081101561101f57600080fd5b8101908080519060200190929190505050905060008090505b600360008881526020019081526020016000208054905081101561128d57600060036000898152602001908152602001600020828154811061107657fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166340550a1c826040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b15801561114257600080fd5b505afa158015611156573d6000803e3d6000fd5b505050506040513d602081101561116c57600080fd5b8101908080519060200190929190505050156112715761126e600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663473691a4836040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b15801561122457600080fd5b505afa158015611238573d6000803e3d6000fd5b505050506040513d602081101561124e57600080fd5b81019080805190602001909291905050508561136a90919063ffffffff16565b93505b5061128660018261136a90919063ffffffff16565b9050611038565b5060006112a46003846113f290919063ffffffff16565b905060006112bc6002846113f290919063ffffffff16565b9050600081831015905080838397509750975050505050509193909250565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16636b3ce98c826040518263ffffffff1660e01b815260040180828152602001915050600060405180830381600087803b15801561134f57600080fd5b505af1158015611363573d6000803e3d6000fd5b5050505050565b6000808284019050838110156113e8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f536166654d6174683a206164646974696f6e206f766572666c6f77000000000081525060200191505060405180910390fd5b8091505092915050565b6000808314156114055760009050611472565b600082840290508284828161141657fe5b041461146d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602181526020018061152a6021913960400191505060405180910390fd5b809150505b9291505056fe5468652063756d756c617469766520706f776572206f66207369676e61746f72792076616c696461746f727320646f6573206e6f74206d65657420746865207468726573686f6c645468652070726f7068656379206d7573742062652070656e64696e6720666f722074686973206f7065726174696f6e43616e6e6f74206d616b65206475706c6963617465206f7261636c6520636c61696d732066726f6d207468652073616d6520616464726573732e536166654d6174683a206d756c7469706c69636174696f6e206f766572666c6f77a165627a7a72305820f5733700df1e7f626fc8ba6c44a824b0892ce097e2340b2413a7611fc9b1de120029`

// DeployOracle deploys a new Ethereum contract, binding an instance of Oracle to it.
func DeployOracle(auth *bind.TransactOpts, backend bind.ContractBackend, _operator common.Address, _valset common.Address, _cosmosBridge common.Address) (common.Address, *types.Transaction, *Oracle, error) {
  parsed, err := abi.JSON(strings.NewReader(OracleABI))
  if err != nil {
    return common.Address{}, nil, nil, err
  }
  address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(OracleBin), backend, _operator, _valset, _cosmosBridge)
  if err != nil {
    return common.Address{}, nil, nil, err
  }
  return address, tx, &Oracle{OracleCaller: OracleCaller{contract: contract}, OracleTransactor: OracleTransactor{contract: contract}, OracleFilterer: OracleFilterer{contract: contract}}, nil
}

// Oracle is an auto generated Go binding around an Ethereum contract.
type Oracle struct {
  OracleCaller     // Read-only binding to the contract
  OracleTransactor // Write-only binding to the contract
  OracleFilterer   // Log filterer for contract events
}

// OracleCaller is an auto generated read-only Go binding around an Ethereum contract.
type OracleCaller struct {
  contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OracleTransactor struct {
  contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OracleFilterer struct {
  contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OracleSession struct {
  Contract     *Oracle           // Generic contract binding to set the session for
  CallOpts     bind.CallOpts     // Call options to use throughout this session
  TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OracleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OracleCallerSession struct {
  Contract *OracleCaller // Generic contract caller binding to set the session for
  CallOpts bind.CallOpts // Call options to use throughout this session
}

// OracleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OracleTransactorSession struct {
  Contract     *OracleTransactor // Generic contract transactor binding to set the session for
  TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OracleRaw is an auto generated low-level Go binding around an Ethereum contract.
type OracleRaw struct {
  Contract *Oracle // Generic contract binding to access the raw methods on
}

// OracleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OracleCallerRaw struct {
  Contract *OracleCaller // Generic read-only contract binding to access the raw methods on
}

// OracleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OracleTransactorRaw struct {
  Contract *OracleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOracle creates a new instance of Oracle, bound to a specific deployed contract.
func NewOracle(address common.Address, backend bind.ContractBackend) (*Oracle, error) {
  contract, err := bindOracle(address, backend, backend, backend)
  if err != nil {
    return nil, err
  }
  return &Oracle{OracleCaller: OracleCaller{contract: contract}, OracleTransactor: OracleTransactor{contract: contract}, OracleFilterer: OracleFilterer{contract: contract}}, nil
}

// NewOracleCaller creates a new read-only instance of Oracle, bound to a specific deployed contract.
func NewOracleCaller(address common.Address, caller bind.ContractCaller) (*OracleCaller, error) {
  contract, err := bindOracle(address, caller, nil, nil)
  if err != nil {
    return nil, err
  }
  return &OracleCaller{contract: contract}, nil
}

// NewOracleTransactor creates a new write-only instance of Oracle, bound to a specific deployed contract.
func NewOracleTransactor(address common.Address, transactor bind.ContractTransactor) (*OracleTransactor, error) {
  contract, err := bindOracle(address, nil, transactor, nil)
  if err != nil {
    return nil, err
  }
  return &OracleTransactor{contract: contract}, nil
}

// NewOracleFilterer creates a new log filterer instance of Oracle, bound to a specific deployed contract.
func NewOracleFilterer(address common.Address, filterer bind.ContractFilterer) (*OracleFilterer, error) {
  contract, err := bindOracle(address, nil, nil, filterer)
  if err != nil {
    return nil, err
  }
  return &OracleFilterer{contract: contract}, nil
}

// bindOracle binds a generic wrapper to an already deployed contract.
func bindOracle(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
  parsed, err := abi.JSON(strings.NewReader(OracleABI))
  if err != nil {
    return nil, err
  }
  return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Oracle *OracleRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
  return _Oracle.Contract.OracleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Oracle *OracleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
  return _Oracle.Contract.OracleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Oracle *OracleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
  return _Oracle.Contract.OracleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Oracle *OracleCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
  return _Oracle.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Oracle *OracleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
  return _Oracle.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Oracle *OracleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
  return _Oracle.Contract.contract.Transact(opts, method, params...)
}

// CheckBridgeProphecy is a free data retrieval call binding the contract method 0xe33a8b2a.
//
// Solidity: function checkBridgeProphecy(uint256 _prophecyID) constant returns(bool, uint256, uint256)
func (_Oracle *OracleCaller) CheckBridgeProphecy(opts *bind.CallOpts, _prophecyID *big.Int) (bool, *big.Int, *big.Int, error) {
  var (
    ret0 = new(bool)
    ret1 = new(*big.Int)
    ret2 = new(*big.Int)
  )
  out := &[]interface{}{
    ret0,
    ret1,
    ret2,
  }
  err := _Oracle.contract.Call(opts, out, "checkBridgeProphecy", _prophecyID)
  return *ret0, *ret1, *ret2, err
}

// CheckBridgeProphecy is a free data retrieval call binding the contract method 0xe33a8b2a.
//
// Solidity: function checkBridgeProphecy(uint256 _prophecyID) constant returns(bool, uint256, uint256)
func (_Oracle *OracleSession) CheckBridgeProphecy(_prophecyID *big.Int) (bool, *big.Int, *big.Int, error) {
  return _Oracle.Contract.CheckBridgeProphecy(&_Oracle.CallOpts, _prophecyID)
}

// CheckBridgeProphecy is a free data retrieval call binding the contract method 0xe33a8b2a.
//
// Solidity: function checkBridgeProphecy(uint256 _prophecyID) constant returns(bool, uint256, uint256)
func (_Oracle *OracleCallerSession) CheckBridgeProphecy(_prophecyID *big.Int) (bool, *big.Int, *big.Int, error) {
  return _Oracle.Contract.CheckBridgeProphecy(&_Oracle.CallOpts, _prophecyID)
}

// CosmosBridge is a free data retrieval call binding the contract method 0xb0e9ef71.
//
// Solidity: function cosmosBridge() constant returns(address)
func (_Oracle *OracleCaller) CosmosBridge(opts *bind.CallOpts) (common.Address, error) {
  var (
    ret0 = new(common.Address)
  )
  out := ret0
  err := _Oracle.contract.Call(opts, out, "cosmosBridge")
  return *ret0, err
}

// CosmosBridge is a free data retrieval call binding the contract method 0xb0e9ef71.
//
// Solidity: function cosmosBridge() constant returns(address)
func (_Oracle *OracleSession) CosmosBridge() (common.Address, error) {
  return _Oracle.Contract.CosmosBridge(&_Oracle.CallOpts)
}

// CosmosBridge is a free data retrieval call binding the contract method 0xb0e9ef71.
//
// Solidity: function cosmosBridge() constant returns(address)
func (_Oracle *OracleCallerSession) CosmosBridge() (common.Address, error) {
  return _Oracle.Contract.CosmosBridge(&_Oracle.CallOpts)
}

// HasMadeClaim is a free data retrieval call binding the contract method 0xa219763e.
//
// Solidity: function hasMadeClaim(uint256 , address ) constant returns(bool)
func (_Oracle *OracleCaller) HasMadeClaim(opts *bind.CallOpts, arg0 *big.Int, arg1 common.Address) (bool, error) {
  var (
    ret0 = new(bool)
  )
  out := ret0
  err := _Oracle.contract.Call(opts, out, "hasMadeClaim", arg0, arg1)
  return *ret0, err
}

// HasMadeClaim is a free data retrieval call binding the contract method 0xa219763e.
//
// Solidity: function hasMadeClaim(uint256 , address ) constant returns(bool)
func (_Oracle *OracleSession) HasMadeClaim(arg0 *big.Int, arg1 common.Address) (bool, error) {
  return _Oracle.Contract.HasMadeClaim(&_Oracle.CallOpts, arg0, arg1)
}

// HasMadeClaim is a free data retrieval call binding the contract method 0xa219763e.
//
// Solidity: function hasMadeClaim(uint256 , address ) constant returns(bool)
func (_Oracle *OracleCallerSession) HasMadeClaim(arg0 *big.Int, arg1 common.Address) (bool, error) {
  return _Oracle.Contract.HasMadeClaim(&_Oracle.CallOpts, arg0, arg1)
}

// Operator is a free data retrieval call binding the contract method 0x570ca735.
//
// Solidity: function operator() constant returns(address)
func (_Oracle *OracleCaller) Operator(opts *bind.CallOpts) (common.Address, error) {
  var (
    ret0 = new(common.Address)
  )
  out := ret0
  err := _Oracle.contract.Call(opts, out, "operator")
  return *ret0, err
}

// Operator is a free data retrieval call binding the contract method 0x570ca735.
//
// Solidity: function operator() constant returns(address)
func (_Oracle *OracleSession) Operator() (common.Address, error) {
  return _Oracle.Contract.Operator(&_Oracle.CallOpts)
}

// Operator is a free data retrieval call binding the contract method 0x570ca735.
//
// Solidity: function operator() constant returns(address)
func (_Oracle *OracleCallerSession) Operator() (common.Address, error) {
  return _Oracle.Contract.Operator(&_Oracle.CallOpts)
}

// OracleClaimValidators is a free data retrieval call binding the contract method 0x36e41341.
//
// Solidity: function oracleClaimValidators(uint256 , uint256 ) constant returns(address)
func (_Oracle *OracleCaller) OracleClaimValidators(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (common.Address, error) {
  var (
    ret0 = new(common.Address)
  )
  out := ret0
  err := _Oracle.contract.Call(opts, out, "oracleClaimValidators", arg0, arg1)
  return *ret0, err
}

// OracleClaimValidators is a free data retrieval call binding the contract method 0x36e41341.
//
// Solidity: function oracleClaimValidators(uint256 , uint256 ) constant returns(address)
func (_Oracle *OracleSession) OracleClaimValidators(arg0 *big.Int, arg1 *big.Int) (common.Address, error) {
  return _Oracle.Contract.OracleClaimValidators(&_Oracle.CallOpts, arg0, arg1)
}

// OracleClaimValidators is a free data retrieval call binding the contract method 0x36e41341.
//
// Solidity: function oracleClaimValidators(uint256 , uint256 ) constant returns(address)
func (_Oracle *OracleCallerSession) OracleClaimValidators(arg0 *big.Int, arg1 *big.Int) (common.Address, error) {
  return _Oracle.Contract.OracleClaimValidators(&_Oracle.CallOpts, arg0, arg1)
}

// Valset is a free data retrieval call binding the contract method 0x7f54af0c.
//
// Solidity: function valset() constant returns(address)
func (_Oracle *OracleCaller) Valset(opts *bind.CallOpts) (common.Address, error) {
  var (
    ret0 = new(common.Address)
  )
  out := ret0
  err := _Oracle.contract.Call(opts, out, "valset")
  return *ret0, err
}

// Valset is a free data retrieval call binding the contract method 0x7f54af0c.
//
// Solidity: function valset() constant returns(address)
func (_Oracle *OracleSession) Valset() (common.Address, error) {
  return _Oracle.Contract.Valset(&_Oracle.CallOpts)
}

// Valset is a free data retrieval call binding the contract method 0x7f54af0c.
//
// Solidity: function valset() constant returns(address)
func (_Oracle *OracleCallerSession) Valset() (common.Address, error) {
  return _Oracle.Contract.Valset(&_Oracle.CallOpts)
}

// NewOracleClaim is a paid mutator transaction binding the contract method 0xc484a8a8.
//
// Solidity: function newOracleClaim(uint256 _prophecyID, bytes32 _message, uint8 _v, bytes32 _r, bytes32 _s) returns()
func (_Oracle *OracleTransactor) NewOracleClaim(opts *bind.TransactOpts, _prophecyID *big.Int, _message [32]byte, _v uint8, _r [32]byte, _s [32]byte) (*types.Transaction, error) {
  return _Oracle.contract.Transact(opts, "newOracleClaim", _prophecyID, _message, _v, _r, _s)
}

// NewOracleClaim is a paid mutator transaction binding the contract method 0xc484a8a8.
//
// Solidity: function newOracleClaim(uint256 _prophecyID, bytes32 _message, uint8 _v, bytes32 _r, bytes32 _s) returns()
func (_Oracle *OracleSession) NewOracleClaim(_prophecyID *big.Int, _message [32]byte, _v uint8, _r [32]byte, _s [32]byte) (*types.Transaction, error) {
  return _Oracle.Contract.NewOracleClaim(&_Oracle.TransactOpts, _prophecyID, _message, _v, _r, _s)
}

// NewOracleClaim is a paid mutator transaction binding the contract method 0xc484a8a8.
//
// Solidity: function newOracleClaim(uint256 _prophecyID, bytes32 _message, uint8 _v, bytes32 _r, bytes32 _s) returns()
func (_Oracle *OracleTransactorSession) NewOracleClaim(_prophecyID *big.Int, _message [32]byte, _v uint8, _r [32]byte, _s [32]byte) (*types.Transaction, error) {
  return _Oracle.Contract.NewOracleClaim(&_Oracle.TransactOpts, _prophecyID, _message, _v, _r, _s)
}

// ProcessBridgeProphecy is a paid mutator transaction binding the contract method 0x89ed70b7.
//
// Solidity: function processBridgeProphecy(uint256 _prophecyID) returns()
func (_Oracle *OracleTransactor) ProcessBridgeProphecy(opts *bind.TransactOpts, _prophecyID *big.Int) (*types.Transaction, error) {
  return _Oracle.contract.Transact(opts, "processBridgeProphecy", _prophecyID)
}

// ProcessBridgeProphecy is a paid mutator transaction binding the contract method 0x89ed70b7.
//
// Solidity: function processBridgeProphecy(uint256 _prophecyID) returns()
func (_Oracle *OracleSession) ProcessBridgeProphecy(_prophecyID *big.Int) (*types.Transaction, error) {
  return _Oracle.Contract.ProcessBridgeProphecy(&_Oracle.TransactOpts, _prophecyID)
}

// ProcessBridgeProphecy is a paid mutator transaction binding the contract method 0x89ed70b7.
//
// Solidity: function processBridgeProphecy(uint256 _prophecyID) returns()
func (_Oracle *OracleTransactorSession) ProcessBridgeProphecy(_prophecyID *big.Int) (*types.Transaction, error) {
  return _Oracle.Contract.ProcessBridgeProphecy(&_Oracle.TransactOpts, _prophecyID)
}

// OracleLogNewOracleClaimIterator is returned from FilterLogNewOracleClaim and is used to iterate over the raw logs and unpacked data for LogNewOracleClaim events raised by the Oracle contract.
type OracleLogNewOracleClaimIterator struct {
  Event *OracleLogNewOracleClaim // Event containing the contract specifics and raw log

  contract *bind.BoundContract // Generic contract to use for unpacking event data
  event    string              // Event name to use for unpacking event data

  logs chan types.Log        // Log channel receiving the found contract events
  sub  ethereum.Subscription // Subscription for errors, completion and termination
  done bool                  // Whether the subscription completed delivering logs
  fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OracleLogNewOracleClaimIterator) Next() bool {
  // If the iterator failed, stop iterating
  if it.fail != nil {
    return false
  }
  // If the iterator completed, deliver directly whatever's available
  if it.done {
    select {
    case log := <-it.logs:
      it.Event = new(OracleLogNewOracleClaim)
      if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
        it.fail = err
        return false
      }
      it.Event.Raw = log
      return true

    default:
      return false
    }
  }
  // Iterator still in progress, wait for either a data or an error event
  select {
  case log := <-it.logs:
    it.Event = new(OracleLogNewOracleClaim)
    if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
      it.fail = err
      return false
    }
    it.Event.Raw = log
    return true

  case err := <-it.sub.Err():
    it.done = true
    it.fail = err
    return it.Next()
  }
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OracleLogNewOracleClaimIterator) Error() error {
  return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleLogNewOracleClaimIterator) Close() error {
  it.sub.Unsubscribe()
  return nil
}

// OracleLogNewOracleClaim represents a LogNewOracleClaim event raised by the Oracle contract.
type OracleLogNewOracleClaim struct {
  ProphecyID       *big.Int
  ValidatorAddress common.Address
  Message          [32]byte
  V                uint8
  R                [32]byte
  S                [32]byte
  Raw              types.Log // Blockchain specific contextual infos
}

// FilterLogNewOracleClaim is a free log retrieval operation binding the contract event 0x0e596df721fcd7ca12494ee97d05d86caaafbf3d92a47f3557947fbd0e4aba27.
//
// Solidity: event LogNewOracleClaim(uint256 _prophecyID, address _validatorAddress, bytes32 _message, uint8 _v, bytes32 _r, bytes32 _s)
func (_Oracle *OracleFilterer) FilterLogNewOracleClaim(opts *bind.FilterOpts) (*OracleLogNewOracleClaimIterator, error) {

  logs, sub, err := _Oracle.contract.FilterLogs(opts, "LogNewOracleClaim")
  if err != nil {
    return nil, err
  }
  return &OracleLogNewOracleClaimIterator{contract: _Oracle.contract, event: "LogNewOracleClaim", logs: logs, sub: sub}, nil
}

// WatchLogNewOracleClaim is a free log subscription operation binding the contract event 0x0e596df721fcd7ca12494ee97d05d86caaafbf3d92a47f3557947fbd0e4aba27.
//
// Solidity: event LogNewOracleClaim(uint256 _prophecyID, address _validatorAddress, bytes32 _message, uint8 _v, bytes32 _r, bytes32 _s)
func (_Oracle *OracleFilterer) WatchLogNewOracleClaim(opts *bind.WatchOpts, sink chan<- *OracleLogNewOracleClaim) (event.Subscription, error) {

  logs, sub, err := _Oracle.contract.WatchLogs(opts, "LogNewOracleClaim")
  if err != nil {
    return nil, err
  }
  return event.NewSubscription(func(quit <-chan struct{}) error {
    defer sub.Unsubscribe()
    for {
      select {
      case log := <-logs:
        // New log arrived, parse the event and forward to the user
        event := new(OracleLogNewOracleClaim)
        if err := _Oracle.contract.UnpackLog(event, "LogNewOracleClaim", log); err != nil {
          return err
        }
        event.Raw = log

        select {
        case sink <- event:
        case err := <-sub.Err():
          return err
        case <-quit:
          return nil
        }
      case err := <-sub.Err():
        return err
      case <-quit:
        return nil
      }
    }
  }), nil
}

// OracleLogProphecyProcessedIterator is returned from FilterLogProphecyProcessed and is used to iterate over the raw logs and unpacked data for LogProphecyProcessed events raised by the Oracle contract.
type OracleLogProphecyProcessedIterator struct {
  Event *OracleLogProphecyProcessed // Event containing the contract specifics and raw log

  contract *bind.BoundContract // Generic contract to use for unpacking event data
  event    string              // Event name to use for unpacking event data

  logs chan types.Log        // Log channel receiving the found contract events
  sub  ethereum.Subscription // Subscription for errors, completion and termination
  done bool                  // Whether the subscription completed delivering logs
  fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OracleLogProphecyProcessedIterator) Next() bool {
  // If the iterator failed, stop iterating
  if it.fail != nil {
    return false
  }
  // If the iterator completed, deliver directly whatever's available
  if it.done {
    select {
    case log := <-it.logs:
      it.Event = new(OracleLogProphecyProcessed)
      if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
        it.fail = err
        return false
      }
      it.Event.Raw = log
      return true

    default:
      return false
    }
  }
  // Iterator still in progress, wait for either a data or an error event
  select {
  case log := <-it.logs:
    it.Event = new(OracleLogProphecyProcessed)
    if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
      it.fail = err
      return false
    }
    it.Event.Raw = log
    return true

  case err := <-it.sub.Err():
    it.done = true
    it.fail = err
    return it.Next()
  }
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OracleLogProphecyProcessedIterator) Error() error {
  return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleLogProphecyProcessedIterator) Close() error {
  it.sub.Unsubscribe()
  return nil
}

// OracleLogProphecyProcessed represents a LogProphecyProcessed event raised by the Oracle contract.
type OracleLogProphecyProcessed struct {
  ProphecyID          *big.Int
  WeightedSignedPower *big.Int
  WeightedTotalPower  *big.Int
  Submitter           common.Address
  Raw                 types.Log // Blockchain specific contextual infos
}

// FilterLogProphecyProcessed is a free log retrieval operation binding the contract event 0x1d8e3fbd601d9d92db7022fb97f75e132841b94db732dcecb0c93cb31852fcbc.
//
// Solidity: event LogProphecyProcessed(uint256 _prophecyID, uint256 _weightedSignedPower, uint256 _weightedTotalPower, address _submitter)
func (_Oracle *OracleFilterer) FilterLogProphecyProcessed(opts *bind.FilterOpts) (*OracleLogProphecyProcessedIterator, error) {

  logs, sub, err := _Oracle.contract.FilterLogs(opts, "LogProphecyProcessed")
  if err != nil {
    return nil, err
  }
  return &OracleLogProphecyProcessedIterator{contract: _Oracle.contract, event: "LogProphecyProcessed", logs: logs, sub: sub}, nil
}

// WatchLogProphecyProcessed is a free log subscription operation binding the contract event 0x1d8e3fbd601d9d92db7022fb97f75e132841b94db732dcecb0c93cb31852fcbc.
//
// Solidity: event LogProphecyProcessed(uint256 _prophecyID, uint256 _weightedSignedPower, uint256 _weightedTotalPower, address _submitter)
func (_Oracle *OracleFilterer) WatchLogProphecyProcessed(opts *bind.WatchOpts, sink chan<- *OracleLogProphecyProcessed) (event.Subscription, error) {

  logs, sub, err := _Oracle.contract.WatchLogs(opts, "LogProphecyProcessed")
  if err != nil {
    return nil, err
  }
  return event.NewSubscription(func(quit <-chan struct{}) error {
    defer sub.Unsubscribe()
    for {
      select {
      case log := <-logs:
        // New log arrived, parse the event and forward to the user
        event := new(OracleLogProphecyProcessed)
        if err := _Oracle.contract.UnpackLog(event, "LogProphecyProcessed", log); err != nil {
          return err
        }
        event.Raw = log

        select {
        case sink <- event:
        case err := <-sub.Err():
          return err
        case <-quit:
          return nil
        }
      case err := <-sub.Err():
        return err
      case <-quit:
        return nil
      }
    }
  }), nil
}
