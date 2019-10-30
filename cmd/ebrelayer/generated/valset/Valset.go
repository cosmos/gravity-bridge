// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Valset

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

// ValsetABI is the input ABI used to generate the binding from.
const ValsetABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"seqCounter\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"validatorCount\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"powers\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"h\",\"type\":\"bytes32\"},{\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"recover\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"validators\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_validator\",\"type\":\"address\"}],\"name\":\"isActiveValidator\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_validator\",\"type\":\"address\"}],\"name\":\"getValidatorPower\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getTotalPower\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"operator\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"h\",\"type\":\"bytes32\"}],\"name\":\"toEthSignedMessageHash\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"activeValidators\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalPower\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_operator\",\"type\":\"address\"},{\"name\":\"_initValidatorAddresses\",\"type\":\"address[]\"},{\"name\":\"_initValidatorPowers\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// ValsetBin is the compiled bytecode used for deploying new contracts.
const ValsetBin = `6080604052600060035534801561001557600080fd5b50604051610b23380380610b238339810180604052606081101561003857600080fd5b8101908080519060200190929190805164010000000081111561005a57600080fd5b8281019050602081018481111561007057600080fd5b815185602082028301116401000000008211171561008d57600080fd5b505092919060200180516401000000008111156100a957600080fd5b828101905060208101848111156100bf57600080fd5b81518560208202830111640100000000821117156100dc57600080fd5b5050929190505050826000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506000600281905550815160018190555060008090505b6001548110156102585760016005600085848151811061015657fe5b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055508181815181106101bb57fe5b6020026020010151600660008584815181106101d357fe5b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208190555061024582828151811061022857fe5b602002602001015160025461026160201b6107771790919060201c565b600281905550808060010191505061013a565b505050506102e9565b6000808284019050838110156102df576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f536166654d6174683a206164646974696f6e206f766572666c6f77000000000081525060200191505060405180910390fd5b8091505092915050565b61082b806102f86000396000f3fe608060405234801561001057600080fd5b50600436106100b45760003560e01c8063473691a411610071578063473691a41461031c57806353976a2614610374578063570ca73514610392578063918a15cf146103dc578063ba26e6121461041e578063db3ad22c1461047a576100b4565b80630904baaa146100b95780630f43a677146100d75780630fd74ee0146100f557806319045a251461014d57806335aa2e441461025257806340550a1c146102c0575b600080fd5b6100c1610498565b6040518082815260200191505060405180910390f35b6100df61049e565b6040518082815260200191505060405180910390f35b6101376004803603602081101561010b57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506104a4565b6040518082815260200191505060405180910390f35b6102106004803603604081101561016357600080fd5b81019080803590602001909291908035906020019064010000000081111561018a57600080fd5b82018360208201111561019c57600080fd5b803590602001918460018302840111640100000000831117156101be57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505091929192905050506104bc565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b61027e6004803603602081101561026857600080fd5b81019080803590602001909291905050506104d9565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b610302600480360360208110156102d657600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610515565b604051808215151515815260200191505060405180910390f35b61035e6004803603602081101561033257600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061056b565b6040518082815260200191505060405180910390f35b61037c6105b4565b6040518082815260200191505060405180910390f35b61039a6105be565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b610408600480360360208110156103f257600080fd5b81019080803590602001909291905050506105e3565b6040518082815260200191505060405180910390f35b6104606004803603602081101561043457600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506105f5565b604051808215151515815260200191505060405180910390f35b610482610615565b6040518082815260200191505060405180910390f35b60035481565b60015481565b60066020528060005260406000206000915090505481565b60006104d1828461061b90919063ffffffff16565b905092915050565b600481815481106104e657fe5b906000526020600020016000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000600560008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff169050919050565b6000600660008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b6000600254905090565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60006105ee8261071f565b9050919050565b60056020528060005260406000206000915054906101000a900460ff1681565b60025481565b6000604182511461062f5760009050610719565b60008060006020850151925060408501519150606085015160001a90507f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a08260001c11156106835760009350505050610719565b601b8160ff161415801561069b5750601c8160ff1614155b156106ac5760009350505050610719565b60018682858560405160008152602001604052604051808581526020018460ff1660ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa158015610709573d6000803e3d6000fd5b5050506020604051035193505050505b92915050565b60008160405160200180807f19457468657265756d205369676e6564204d6573736167653a0a333200000000815250601c01828152602001915050604051602081830303815290604052805190602001209050919050565b6000808284019050838110156107f5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f536166654d6174683a206164646974696f6e206f766572666c6f77000000000081525060200191505060405180910390fd5b809150509291505056fea165627a7a7230582044e5de7c59e8fa037a06eebfabfd0b502b5b53081ea05bfb565a91c9c03a0e840029`

// DeployValset deploys a new Ethereum contract, binding an instance of Valset to it.
func DeployValset(auth *bind.TransactOpts, backend bind.ContractBackend, _operator common.Address, _initValidatorAddresses []common.Address, _initValidatorPowers []*big.Int) (common.Address, *types.Transaction, *Valset, error) {
	parsed, err := abi.JSON(strings.NewReader(ValsetABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ValsetBin), backend, _operator, _initValidatorAddresses, _initValidatorPowers)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Valset{ValsetCaller: ValsetCaller{contract: contract}, ValsetTransactor: ValsetTransactor{contract: contract}, ValsetFilterer: ValsetFilterer{contract: contract}}, nil
}

// Valset is an auto generated Go binding around an Ethereum contract.
type Valset struct {
	ValsetCaller     // Read-only binding to the contract
	ValsetTransactor // Write-only binding to the contract
	ValsetFilterer   // Log filterer for contract events
}

// ValsetCaller is an auto generated read-only Go binding around an Ethereum contract.
type ValsetCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ValsetTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ValsetTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ValsetFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ValsetFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ValsetSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ValsetSession struct {
	Contract     *Valset           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ValsetCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ValsetCallerSession struct {
	Contract *ValsetCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ValsetTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ValsetTransactorSession struct {
	Contract     *ValsetTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ValsetRaw is an auto generated low-level Go binding around an Ethereum contract.
type ValsetRaw struct {
	Contract *Valset // Generic contract binding to access the raw methods on
}

// ValsetCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ValsetCallerRaw struct {
	Contract *ValsetCaller // Generic read-only contract binding to access the raw methods on
}

// ValsetTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ValsetTransactorRaw struct {
	Contract *ValsetTransactor // Generic write-only contract binding to access the raw methods on
}

// NewValset creates a new instance of Valset, bound to a specific deployed contract.
func NewValset(address common.Address, backend bind.ContractBackend) (*Valset, error) {
	contract, err := bindValset(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Valset{ValsetCaller: ValsetCaller{contract: contract}, ValsetTransactor: ValsetTransactor{contract: contract}, ValsetFilterer: ValsetFilterer{contract: contract}}, nil
}

// NewValsetCaller creates a new read-only instance of Valset, bound to a specific deployed contract.
func NewValsetCaller(address common.Address, caller bind.ContractCaller) (*ValsetCaller, error) {
	contract, err := bindValset(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ValsetCaller{contract: contract}, nil
}

// NewValsetTransactor creates a new write-only instance of Valset, bound to a specific deployed contract.
func NewValsetTransactor(address common.Address, transactor bind.ContractTransactor) (*ValsetTransactor, error) {
	contract, err := bindValset(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ValsetTransactor{contract: contract}, nil
}

// NewValsetFilterer creates a new log filterer instance of Valset, bound to a specific deployed contract.
func NewValsetFilterer(address common.Address, filterer bind.ContractFilterer) (*ValsetFilterer, error) {
	contract, err := bindValset(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ValsetFilterer{contract: contract}, nil
}

// bindValset binds a generic wrapper to an already deployed contract.
func bindValset(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ValsetABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Valset *ValsetRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Valset.Contract.ValsetCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Valset *ValsetRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Valset.Contract.ValsetTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Valset *ValsetRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Valset.Contract.ValsetTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Valset *ValsetCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Valset.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Valset *ValsetTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Valset.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Valset *ValsetTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Valset.Contract.contract.Transact(opts, method, params...)
}

// ActiveValidators is a free data retrieval call binding the contract method 0xba26e612.
//
// Solidity: function activeValidators(address ) constant returns(bool)
func (_Valset *ValsetCaller) ActiveValidators(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Valset.contract.Call(opts, out, "activeValidators", arg0)
	return *ret0, err
}

// ActiveValidators is a free data retrieval call binding the contract method 0xba26e612.
//
// Solidity: function activeValidators(address ) constant returns(bool)
func (_Valset *ValsetSession) ActiveValidators(arg0 common.Address) (bool, error) {
	return _Valset.Contract.ActiveValidators(&_Valset.CallOpts, arg0)
}

// ActiveValidators is a free data retrieval call binding the contract method 0xba26e612.
//
// Solidity: function activeValidators(address ) constant returns(bool)
func (_Valset *ValsetCallerSession) ActiveValidators(arg0 common.Address) (bool, error) {
	return _Valset.Contract.ActiveValidators(&_Valset.CallOpts, arg0)
}

// GetTotalPower is a free data retrieval call binding the contract method 0x53976a26.
//
// Solidity: function getTotalPower() constant returns(uint256)
func (_Valset *ValsetCaller) GetTotalPower(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Valset.contract.Call(opts, out, "getTotalPower")
	return *ret0, err
}

// GetTotalPower is a free data retrieval call binding the contract method 0x53976a26.
//
// Solidity: function getTotalPower() constant returns(uint256)
func (_Valset *ValsetSession) GetTotalPower() (*big.Int, error) {
	return _Valset.Contract.GetTotalPower(&_Valset.CallOpts)
}

// GetTotalPower is a free data retrieval call binding the contract method 0x53976a26.
//
// Solidity: function getTotalPower() constant returns(uint256)
func (_Valset *ValsetCallerSession) GetTotalPower() (*big.Int, error) {
	return _Valset.Contract.GetTotalPower(&_Valset.CallOpts)
}

// GetValidatorPower is a free data retrieval call binding the contract method 0x473691a4.
//
// Solidity: function getValidatorPower(address _validator) constant returns(uint256)
func (_Valset *ValsetCaller) GetValidatorPower(opts *bind.CallOpts, _validator common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Valset.contract.Call(opts, out, "getValidatorPower", _validator)
	return *ret0, err
}

// GetValidatorPower is a free data retrieval call binding the contract method 0x473691a4.
//
// Solidity: function getValidatorPower(address _validator) constant returns(uint256)
func (_Valset *ValsetSession) GetValidatorPower(_validator common.Address) (*big.Int, error) {
	return _Valset.Contract.GetValidatorPower(&_Valset.CallOpts, _validator)
}

// GetValidatorPower is a free data retrieval call binding the contract method 0x473691a4.
//
// Solidity: function getValidatorPower(address _validator) constant returns(uint256)
func (_Valset *ValsetCallerSession) GetValidatorPower(_validator common.Address) (*big.Int, error) {
	return _Valset.Contract.GetValidatorPower(&_Valset.CallOpts, _validator)
}

// IsActiveValidator is a free data retrieval call binding the contract method 0x40550a1c.
//
// Solidity: function isActiveValidator(address _validator) constant returns(bool)
func (_Valset *ValsetCaller) IsActiveValidator(opts *bind.CallOpts, _validator common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Valset.contract.Call(opts, out, "isActiveValidator", _validator)
	return *ret0, err
}

// IsActiveValidator is a free data retrieval call binding the contract method 0x40550a1c.
//
// Solidity: function isActiveValidator(address _validator) constant returns(bool)
func (_Valset *ValsetSession) IsActiveValidator(_validator common.Address) (bool, error) {
	return _Valset.Contract.IsActiveValidator(&_Valset.CallOpts, _validator)
}

// IsActiveValidator is a free data retrieval call binding the contract method 0x40550a1c.
//
// Solidity: function isActiveValidator(address _validator) constant returns(bool)
func (_Valset *ValsetCallerSession) IsActiveValidator(_validator common.Address) (bool, error) {
	return _Valset.Contract.IsActiveValidator(&_Valset.CallOpts, _validator)
}

// Operator is a free data retrieval call binding the contract method 0x570ca735.
//
// Solidity: function operator() constant returns(address)
func (_Valset *ValsetCaller) Operator(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Valset.contract.Call(opts, out, "operator")
	return *ret0, err
}

// Operator is a free data retrieval call binding the contract method 0x570ca735.
//
// Solidity: function operator() constant returns(address)
func (_Valset *ValsetSession) Operator() (common.Address, error) {
	return _Valset.Contract.Operator(&_Valset.CallOpts)
}

// Operator is a free data retrieval call binding the contract method 0x570ca735.
//
// Solidity: function operator() constant returns(address)
func (_Valset *ValsetCallerSession) Operator() (common.Address, error) {
	return _Valset.Contract.Operator(&_Valset.CallOpts)
}

// Powers is a free data retrieval call binding the contract method 0x0fd74ee0.
//
// Solidity: function powers(address ) constant returns(uint256)
func (_Valset *ValsetCaller) Powers(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Valset.contract.Call(opts, out, "powers", arg0)
	return *ret0, err
}

// Powers is a free data retrieval call binding the contract method 0x0fd74ee0.
//
// Solidity: function powers(address ) constant returns(uint256)
func (_Valset *ValsetSession) Powers(arg0 common.Address) (*big.Int, error) {
	return _Valset.Contract.Powers(&_Valset.CallOpts, arg0)
}

// Powers is a free data retrieval call binding the contract method 0x0fd74ee0.
//
// Solidity: function powers(address ) constant returns(uint256)
func (_Valset *ValsetCallerSession) Powers(arg0 common.Address) (*big.Int, error) {
	return _Valset.Contract.Powers(&_Valset.CallOpts, arg0)
}

// Recover is a free data retrieval call binding the contract method 0x19045a25.
//
// Solidity: function recover(bytes32 h, bytes signature) constant returns(address)
func (_Valset *ValsetCaller) Recover(opts *bind.CallOpts, h [32]byte, signature []byte) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Valset.contract.Call(opts, out, "recover", h, signature)
	return *ret0, err
}

// Recover is a free data retrieval call binding the contract method 0x19045a25.
//
// Solidity: function recover(bytes32 h, bytes signature) constant returns(address)
func (_Valset *ValsetSession) Recover(h [32]byte, signature []byte) (common.Address, error) {
	return _Valset.Contract.Recover(&_Valset.CallOpts, h, signature)
}

// Recover is a free data retrieval call binding the contract method 0x19045a25.
//
// Solidity: function recover(bytes32 h, bytes signature) constant returns(address)
func (_Valset *ValsetCallerSession) Recover(h [32]byte, signature []byte) (common.Address, error) {
	return _Valset.Contract.Recover(&_Valset.CallOpts, h, signature)
}

// SeqCounter is a free data retrieval call binding the contract method 0x0904baaa.
//
// Solidity: function seqCounter() constant returns(uint256)
func (_Valset *ValsetCaller) SeqCounter(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Valset.contract.Call(opts, out, "seqCounter")
	return *ret0, err
}

// SeqCounter is a free data retrieval call binding the contract method 0x0904baaa.
//
// Solidity: function seqCounter() constant returns(uint256)
func (_Valset *ValsetSession) SeqCounter() (*big.Int, error) {
	return _Valset.Contract.SeqCounter(&_Valset.CallOpts)
}

// SeqCounter is a free data retrieval call binding the contract method 0x0904baaa.
//
// Solidity: function seqCounter() constant returns(uint256)
func (_Valset *ValsetCallerSession) SeqCounter() (*big.Int, error) {
	return _Valset.Contract.SeqCounter(&_Valset.CallOpts)
}

// ToEthSignedMessageHash is a free data retrieval call binding the contract method 0x918a15cf.
//
// Solidity: function toEthSignedMessageHash(bytes32 h) constant returns(bytes32)
func (_Valset *ValsetCaller) ToEthSignedMessageHash(opts *bind.CallOpts, h [32]byte) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Valset.contract.Call(opts, out, "toEthSignedMessageHash", h)
	return *ret0, err
}

// ToEthSignedMessageHash is a free data retrieval call binding the contract method 0x918a15cf.
//
// Solidity: function toEthSignedMessageHash(bytes32 h) constant returns(bytes32)
func (_Valset *ValsetSession) ToEthSignedMessageHash(h [32]byte) ([32]byte, error) {
	return _Valset.Contract.ToEthSignedMessageHash(&_Valset.CallOpts, h)
}

// ToEthSignedMessageHash is a free data retrieval call binding the contract method 0x918a15cf.
//
// Solidity: function toEthSignedMessageHash(bytes32 h) constant returns(bytes32)
func (_Valset *ValsetCallerSession) ToEthSignedMessageHash(h [32]byte) ([32]byte, error) {
	return _Valset.Contract.ToEthSignedMessageHash(&_Valset.CallOpts, h)
}

// TotalPower is a free data retrieval call binding the contract method 0xdb3ad22c.
//
// Solidity: function totalPower() constant returns(uint256)
func (_Valset *ValsetCaller) TotalPower(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Valset.contract.Call(opts, out, "totalPower")
	return *ret0, err
}

// TotalPower is a free data retrieval call binding the contract method 0xdb3ad22c.
//
// Solidity: function totalPower() constant returns(uint256)
func (_Valset *ValsetSession) TotalPower() (*big.Int, error) {
	return _Valset.Contract.TotalPower(&_Valset.CallOpts)
}

// TotalPower is a free data retrieval call binding the contract method 0xdb3ad22c.
//
// Solidity: function totalPower() constant returns(uint256)
func (_Valset *ValsetCallerSession) TotalPower() (*big.Int, error) {
	return _Valset.Contract.TotalPower(&_Valset.CallOpts)
}

// ValidatorCount is a free data retrieval call binding the contract method 0x0f43a677.
//
// Solidity: function validatorCount() constant returns(uint256)
func (_Valset *ValsetCaller) ValidatorCount(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Valset.contract.Call(opts, out, "validatorCount")
	return *ret0, err
}

// ValidatorCount is a free data retrieval call binding the contract method 0x0f43a677.
//
// Solidity: function validatorCount() constant returns(uint256)
func (_Valset *ValsetSession) ValidatorCount() (*big.Int, error) {
	return _Valset.Contract.ValidatorCount(&_Valset.CallOpts)
}

// ValidatorCount is a free data retrieval call binding the contract method 0x0f43a677.
//
// Solidity: function validatorCount() constant returns(uint256)
func (_Valset *ValsetCallerSession) ValidatorCount() (*big.Int, error) {
	return _Valset.Contract.ValidatorCount(&_Valset.CallOpts)
}

// Validators is a free data retrieval call binding the contract method 0x35aa2e44.
//
// Solidity: function validators(uint256 ) constant returns(address)
func (_Valset *ValsetCaller) Validators(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Valset.contract.Call(opts, out, "validators", arg0)
	return *ret0, err
}

// Validators is a free data retrieval call binding the contract method 0x35aa2e44.
//
// Solidity: function validators(uint256 ) constant returns(address)
func (_Valset *ValsetSession) Validators(arg0 *big.Int) (common.Address, error) {
	return _Valset.Contract.Validators(&_Valset.CallOpts, arg0)
}

// Validators is a free data retrieval call binding the contract method 0x35aa2e44.
//
// Solidity: function validators(uint256 ) constant returns(address)
func (_Valset *ValsetCallerSession) Validators(arg0 *big.Int) (common.Address, error) {
	return _Valset.Contract.Validators(&_Valset.CallOpts, arg0)
}
