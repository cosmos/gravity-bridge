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
const OracleABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"oracleClaimValidators\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_prophecyID\",\"type\":\"uint256\"},{\"name\":\"_message\",\"type\":\"bytes32\"},{\"name\":\"_signature\",\"type\":\"bytes\"}],\"name\":\"newOracleClaim\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"operator\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"valset\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_prophecyID\",\"type\":\"uint256\"}],\"name\":\"processBridgeProphecy\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"hasMadeClaim\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"cosmosBridge\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_prophecyID\",\"type\":\"uint256\"}],\"name\":\"checkBridgeProphecy\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_operator\",\"type\":\"address\"},{\"name\":\"_valset\",\"type\":\"address\"},{\"name\":\"_cosmosBridge\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_prophecyID\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"_validatorAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_message\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_signature\",\"type\":\"bytes\"}],\"name\":\"LogNewOracleClaim\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_prophecyID\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"_weightedSignedPower\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"_weightedTotalPower\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"_submitter\",\"type\":\"address\"}],\"name\":\"LogProphecyProcessed\",\"type\":\"event\"}]"

// OracleBin is the compiled bytecode used for deploying new contracts.
const OracleBin = `608060405234801561001057600080fd5b506040516060806117428339810180604052606081101561003057600080fd5b8101908080519060200190929190805190602001909291908051906020019092919050505082600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555081600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550505050611619806101296000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c806389ed70b71161005b57806389ed70b714610268578063a219763e14610296578063b0e9ef71146102fc578063e33a8b2a1461034657610088565b806336e413411461008d578063568b3c4f14610105578063570ca735146101d45780637f54af0c1461021e575b600080fd5b6100c3600480360360408110156100a357600080fd5b81019080803590602001909291908035906020019092919050505061039a565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6101d26004803603606081101561011b57600080fd5b8101908080359060200190929190803590602001909291908035906020019064010000000081111561014c57600080fd5b82018360208201111561015e57600080fd5b8035906020019184600183028401116401000000008311171561018057600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505091929192905050506103e5565b005b6101dc610a76565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b610226610a9c565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6102946004803603602081101561027e57600080fd5b8101908080359060200190929190505050610ac2565b005b6102e2600480360360408110156102ac57600080fd5b8101908080359060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610cc1565b604051808215151515815260200191505060405180910390f35b610304610cf0565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6103726004803603602081101561035c57600080fd5b8101908080359060200190929190505050610d15565b6040518084151515158152602001838152602001828152602001935050505060405180910390f35b600360205281600052604060002081815481106103b357fe5b906000526020600020016000915091509054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166340550a1c336040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b15801561048457600080fd5b505afa158015610498573d6000803e3d6000fd5b505050506040513d60208110156104ae57600080fd5b8101908080519060200190929190505050610531576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f4d75737420626520616e206163746976652076616c696461746f72000000000081525060200191505060405180910390fd5b82600115156000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663d8da69ea836040518263ffffffff1660e01b81526004018082815260200191505060206040518083038186803b1580156105a857600080fd5b505afa1580156105bc573d6000803e3d6000fd5b505050506040513d60208110156105d257600080fd5b810190808051906020019092919050505015151461063b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602f815260200180611564602f913960400191505060405180910390fd5b6000339050600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166319045a2585856040518363ffffffff1660e01b81526004018083815260200180602001828103825283818151815260200191508051906020019080838360005b838110156106d55780820151818401526020810190506106ba565b50505050905090810190601f1680156107025780820380516001836020036101000a031916815260200191505b50935050505060206040518083038186803b15801561072057600080fd5b505afa158015610734573d6000803e3d6000fd5b505050506040513d602081101561074a57600080fd5b810190808051906020019092919050505073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16146107fb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601a8152602001807f496e76616c6964206d657373616765207369676e61747572652e00000000000081525060200191505060405180910390fd5b6004600086815260200190815260200160002060008273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16156108af576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252603a815260200180611593603a913960400191505060405180910390fd5b60016004600087815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff021916908315150217905550600360008681526020019081526020016000208190806001815401808255809150509060018203906000526020600020016000909192909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550507fdada3531177d355f0118a46f374e97a5147f06bcd8b130f884756cf193319a9185828686604051808581526020018473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200183815260200180602001828103825283818151815260200191508051906020019080838360005b83811015610a32578082015181840152602081019050610a17565b50505050905090810190601f168015610a5f5780820380516001836020036101000a031916815260200191505b509550505050505060405180910390a15050505050565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b80600115156000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663d8da69ea836040518263ffffffff1660e01b81526004018082815260200191505060206040518083038186803b158015610b3957600080fd5b505afa158015610b4d573d6000803e3d6000fd5b505050506040513d6020811015610b6357600080fd5b8101908080519060200190929190505050151514610bcc576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602f815260200180611564602f913960400191505060405180910390fd5b6000806000610bda85611024565b92509250925082610c36576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252604881526020018061151c6048913960600191505060405180910390fd5b610c3f8561137e565b7f1d8e3fbd601d9d92db7022fb97f75e132841b94db732dcecb0c93cb31852fcbc85838333604051808581526020018481526020018381526020018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200194505050505060405180910390a15050505050565b60046020528160005260406000206020528060005260406000206000915091509054906101000a900460ff1681565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000806000600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610ddd576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260158152602001807f4d75737420626520746865206f70657261746f722e000000000000000000000081525060200191505060405180910390fd5b83600115156000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663d8da69ea836040518263ffffffff1660e01b81526004018082815260200191505060206040518083038186803b158015610e5457600080fd5b505afa158015610e68573d6000803e3d6000fd5b505050506040513d6020811015610e7e57600080fd5b8101908080519060200190929190505050151514610ee7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602f815260200180611564602f913960400191505060405180910390fd5b600115156000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663d8da69ea876040518263ffffffff1660e01b81526004018082815260200191505060206040518083038186803b158015610f5d57600080fd5b505afa158015610f71573d6000803e3d6000fd5b505050506040513d6020811015610f8757600080fd5b810190808051906020019092919050505015151461100d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f43616e206f6e6c7920636865636b206163746976652070726f7068656369657381525060200191505060405180910390fd5b61101685611024565b935093509350509193909250565b600080600080600090506000600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663db3ad22c6040518163ffffffff1660e01b815260040160206040518083038186803b15801561109857600080fd5b505afa1580156110ac573d6000803e3d6000fd5b505050506040513d60208110156110c257600080fd5b8101908080519060200190929190505050905060008090505b600360008881526020019081526020016000208054905081101561133057600060036000898152602001908152602001600020828154811061111957fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166340550a1c826040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b1580156111e557600080fd5b505afa1580156111f9573d6000803e3d6000fd5b505050506040513d602081101561120f57600080fd5b81019080805190602001909291905050501561131457611311600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663473691a4836040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b1580156112c757600080fd5b505afa1580156112db573d6000803e3d6000fd5b505050506040513d60208110156112f157600080fd5b81019080805190602001909291905050508561140d90919063ffffffff16565b93505b5061132960018261140d90919063ffffffff16565b90506110db565b50600061134760038461149590919063ffffffff16565b9050600061135f60028461149590919063ffffffff16565b9050600081831015905080838397509750975050505050509193909250565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16636b3ce98c826040518263ffffffff1660e01b815260040180828152602001915050600060405180830381600087803b1580156113f257600080fd5b505af1158015611406573d6000803e3d6000fd5b5050505050565b60008082840190508381101561148b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f536166654d6174683a206164646974696f6e206f766572666c6f77000000000081525060200191505060405180910390fd5b8091505092915050565b6000808314156114a85760009050611515565b60008284029050828482816114b957fe5b0414611510576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260218152602001806115cd6021913960400191505060405180910390fd5b809150505b9291505056fe5468652063756d756c617469766520706f776572206f66207369676e61746f72792076616c696461746f727320646f6573206e6f74206d65657420746865207468726573686f6c645468652070726f7068656379206d7573742062652070656e64696e6720666f722074686973206f7065726174696f6e43616e6e6f74206d616b65206475706c6963617465206f7261636c6520636c61696d732066726f6d207468652073616d6520616464726573732e536166654d6174683a206d756c7469706c69636174696f6e206f766572666c6f77a165627a7a72305820bb0277add4f51007c60558223ec975ad8829d53ee1b2519221e05f6bfdc08bf80029`

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

// NewOracleClaim is a paid mutator transaction binding the contract method 0x568b3c4f.
//
// Solidity: function newOracleClaim(uint256 _prophecyID, bytes32 _message, bytes _signature) returns()
func (_Oracle *OracleTransactor) NewOracleClaim(opts *bind.TransactOpts, _prophecyID *big.Int, _message [32]byte, _signature []byte) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "newOracleClaim", _prophecyID, _message, _signature)
}

// NewOracleClaim is a paid mutator transaction binding the contract method 0x568b3c4f.
//
// Solidity: function newOracleClaim(uint256 _prophecyID, bytes32 _message, bytes _signature) returns()
func (_Oracle *OracleSession) NewOracleClaim(_prophecyID *big.Int, _message [32]byte, _signature []byte) (*types.Transaction, error) {
	return _Oracle.Contract.NewOracleClaim(&_Oracle.TransactOpts, _prophecyID, _message, _signature)
}

// NewOracleClaim is a paid mutator transaction binding the contract method 0x568b3c4f.
//
// Solidity: function newOracleClaim(uint256 _prophecyID, bytes32 _message, bytes _signature) returns()
func (_Oracle *OracleTransactorSession) NewOracleClaim(_prophecyID *big.Int, _message [32]byte, _signature []byte) (*types.Transaction, error) {
	return _Oracle.Contract.NewOracleClaim(&_Oracle.TransactOpts, _prophecyID, _message, _signature)
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
	Signature        []byte
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterLogNewOracleClaim is a free log retrieval operation binding the contract event 0xdada3531177d355f0118a46f374e97a5147f06bcd8b130f884756cf193319a91.
//
// Solidity: event LogNewOracleClaim(uint256 _prophecyID, address _validatorAddress, bytes32 _message, bytes _signature)
func (_Oracle *OracleFilterer) FilterLogNewOracleClaim(opts *bind.FilterOpts) (*OracleLogNewOracleClaimIterator, error) {

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "LogNewOracleClaim")
	if err != nil {
		return nil, err
	}
	return &OracleLogNewOracleClaimIterator{contract: _Oracle.contract, event: "LogNewOracleClaim", logs: logs, sub: sub}, nil
}

// WatchLogNewOracleClaim is a free log subscription operation binding the contract event 0xdada3531177d355f0118a46f374e97a5147f06bcd8b130f884756cf193319a91.
//
// Solidity: event LogNewOracleClaim(uint256 _prophecyID, address _validatorAddress, bytes32 _message, bytes _signature)
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
