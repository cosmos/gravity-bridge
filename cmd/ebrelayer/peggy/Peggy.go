// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package peggy

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

// PeggyABI is the input ABI used to generate the binding from.
const PeggyABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"active\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"provider\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_id\",\"type\":\"bytes32\"}],\"name\":\"getStatus\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"activateLocking\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pauseLocking\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_id\",\"type\":\"bytes32\"}],\"name\":\"withdraw\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_recipient\",\"type\":\"bytes\"},{\"name\":\"_token\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"lock\",\"outputs\":[{\"name\":\"_id\",\"type\":\"bytes32\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"nonce\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_id\",\"type\":\"bytes32\"}],\"name\":\"viewItem\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"bytes\"},{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"ids\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_id\",\"type\":\"bytes32\"}],\"name\":\"unlock\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_id\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_to\",\"type\":\"bytes\"},{\"indexed\":false,\"name\":\"_token\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_symbol\",\"type\":\"string\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"_nonce\",\"type\":\"uint256\"}],\"name\":\"LogLock\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_id\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_token\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"_nonce\",\"type\":\"uint256\"}],\"name\":\"LogUnlock\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_id\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_token\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"_nonce\",\"type\":\"uint256\"}],\"name\":\"LogWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_time\",\"type\":\"uint256\"}],\"name\":\"LogLockingPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_time\",\"type\":\"uint256\"}],\"name\":\"LogLockingActivated\",\"type\":\"event\"}]"

// PeggyBin is the compiled bytecode used for deploying new contracts.
const PeggyBin = `608060405234801561001057600080fd5b506000808190555033600260016101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506001600260006101000a81548160ff0219169083151502179055507f9af033c3fdf318cb9968eac8a62b339bd18862abd1703fc74256e9d77cfc95df426040518082815260200191505060405180910390a16120d2806100ba6000396000f3fe60806040526004361061009c5760003560e01c80638e19899e116100645780638e19899e146101a85780639df2a385146101fb578063affed0e0146102f4578063c933dc5b1461031f578063cf7b4a0914610447578063ec9b5b3a1461049a5761009c565b806302fb0c5e146100a1578063085d4883146100d05780635de28ae01461012757806363faf36a1461017a5780638a5cd91e14610191575b600080fd5b3480156100ad57600080fd5b506100b66104ed565b604051808215151515815260200191505060405180910390f35b3480156100dc57600080fd5b506100e5610500565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561013357600080fd5b506101606004803603602081101561014a57600080fd5b8101908080359060200190929190505050610526565b604051808215151515815260200191505060405180910390f35b34801561018657600080fd5b5061018f610538565b005b34801561019d57600080fd5b506101a6610669565b005b3480156101b457600080fd5b506101e1600480360360208110156101cb57600080fd5b8101908080359060200190929190505050610799565b604051808215151515815260200191505060405180910390f35b6102de6004803603606081101561021157600080fd5b810190808035906020019064010000000081111561022e57600080fd5b82018360208201111561024057600080fd5b8035906020019184600183028401116401000000008311171561026257600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610bae565b6040518082815260200191505060405180910390f35b34801561030057600080fd5b506103096110a8565b6040518082815260200191505060405180910390f35b34801561032b57600080fd5b506103586004803603602081101561034257600080fd5b81019080803590602001909291905050506110ae565b604051808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001806020018573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001848152602001838152602001828103825286818151815260200191508051906020019080838360005b838110156104085780820151818401526020810190506103ed565b50505050905090810190601f1680156104355780820380516001836020036101000a031916815260200191505b50965050505050505060405180910390f35b34801561045357600080fd5b506104806004803603602081101561046a57600080fd5b81019080803590602001909291905050506110d6565b604051808215151515815260200191505060405180910390f35b3480156104a657600080fd5b506104d3600480360360208110156104bd57600080fd5b81019080803590602001909291905050506110f6565b604051808215151515815260200191505060405180910390f35b600260009054906101000a900460ff1681565b600260019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000610531826114f5565b9050919050565b600260019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146105fb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601f8152602001807f4d75737420626520746865207370656369666965642070726f76696465722e0081525060200191505060405180910390fd5b600260009054906101000a900460ff161561061557600080fd5b6001600260006101000a81548160ff0219169083151502179055507f9af033c3fdf318cb9968eac8a62b339bd18862abd1703fc74256e9d77cfc95df426040518082815260200191505060405180910390a1565b600260019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461072c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601f8152602001807f4d75737420626520746865207370656369666965642070726f76696465722e0081525060200191505060405180910390fd5b600260009054906101000a900460ff1661074557600080fd5b6000600260006101000a81548160ff0219169083151502179055507fbebc9a19c81e5697fda01edce5ac5aed2c5a0edb9a972fd5f58ac0419a405a82426040518082815260200191505060405180910390a1565b6000816001600082815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610873576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601c8152602001807f4d75737420626520746865206f726967696e616c2073656e6465722e0000000081525060200191505060405180910390fd5b82600073ffffffffffffffffffffffffffffffffffffffff166001600083815260200190815260200160002060020160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16141561096a5760016000828152602001908152602001600020600301543073ffffffffffffffffffffffffffffffffffffffff16311015610965576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602b815260200180611fbf602b913960400191505060405180910390fd5b610ac7565b60016000828152602001908152602001600020600301546001600083815260200190815260200160002060020160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b158015610a3457600080fd5b505afa158015610a48573d6000803e3d6000fd5b505050506040513d6020811015610a5e57600080fd5b81019080805190602001909291905050501015610ac6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602e815260200180612011602e913960400191505060405180910390fd5b5b610ad0846114f5565b610ad957600080fd5b600080600080610ae888611522565b93509350935093507f9cbca76b94cf51b34c3949f0c925da38fe8dbae8e6761e11389884a9c1354b2c8885858585604051808681526020018573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018381526020018281526020019550505050505060405180910390a160019650505050505050919050565b6000805460016000540111610c2b576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f4e6f20617661696c61626c65206e6f6e6365732e00000000000000000000000081525060200191505060405180910390fd5b60011515600260009054906101000a900460ff16151514610c97576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526027815260200180611fea6027913960400191505060405180910390fd5b606060003414610d2357600073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff1614610cda57600080fd5b823414610ce657600080fd5b6040518060400160405280600381526020017f45544800000000000000000000000000000000000000000000000000000000008152509050610efc565b8373ffffffffffffffffffffffffffffffffffffffff166323b872dd3330866040518463ffffffff1660e01b8152600401808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019350505050602060405180830381600087803b158015610dde57600080fd5b505af1158015610df2573d6000803e3d6000fd5b505050506040513d6020811015610e0857600080fd5b8101908080519060200190929190505050610e2257600080fd5b8373ffffffffffffffffffffffffffffffffffffffff166395d89b416040518163ffffffff1660e01b815260040160006040518083038186803b158015610e6857600080fd5b505afa158015610e7c573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f820116820180604052506020811015610ea657600080fd5b810190808051640100000000811115610ebe57600080fd5b82810190506020810184811115610ed457600080fd5b8151856001820283011164010000000082111715610ef157600080fd5b505092919050505090505b6000610f0a33878787611a66565b90507f3945646e76891f1dfa38b4aab98fac226e3f4ad3686493b722b62383358ba922813388888689610f3b611cce565b604051808881526020018773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001806020018673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200180602001858152602001848152602001838103835288818151815260200191508051906020019080838360005b83811015610ff5578082015181840152602081019050610fda565b50505050905090810190601f1680156110225780820380516001836020036101000a031916815260200191505b50838103825286818151815260200191508051906020019080838360005b8381101561105b578082015181840152602081019050611040565b50505050905090810190601f1680156110885780820380516001836020036101000a031916815260200191505b50995050505050505050505060405180910390a180925050509392505050565b60005481565b6000606060008060006110c086611cd7565b8494509450945094509450945091939590929450565b60036020528060005260406000206000915054906101000a900460ff1681565b6000600260019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146111bb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601f8152602001807f4d75737420626520746865207370656369666965642070726f76696465722e0081525060200191505060405180910390fd5b81600073ffffffffffffffffffffffffffffffffffffffff166001600083815260200190815260200160002060020160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614156112b25760016000828152602001908152602001600020600301543073ffffffffffffffffffffffffffffffffffffffff163110156112ad576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602b815260200180611fbf602b913960400191505060405180910390fd5b61140f565b60016000828152602001908152602001600020600301546001600083815260200190815260200160002060020160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b15801561137c57600080fd5b505afa158015611390573d6000803e3d6000fd5b505050506040513d60208110156113a657600080fd5b8101908080519060200190929190505050101561140e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602e815260200180612011602e913960400191505060405180910390fd5b5b611418836114f5565b61142157600080fd5b60008060008061143087611522565b93509350935093507fb3ceeb2ff57376fcabec63d51a010afad847c03e9365f20a168ca66db8b927408785858585604051808681526020018573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018381526020018281526020019550505050505060405180910390a1600195505050505050919050565b60006001600083815260200190815260200160002060050160009054906101000a900460ff169050919050565b60008060008084600073ffffffffffffffffffffffffffffffffffffffff166001600083815260200190815260200160002060020160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16141561161f5760016000828152602001908152602001600020600301543073ffffffffffffffffffffffffffffffffffffffff1631101561161a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602b815260200180611fbf602b913960400191505060405180910390fd5b61177c565b60016000828152602001908152602001600020600301546001600083815260200190815260200160002060020160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060206040518083038186803b1580156116e957600080fd5b505afa1580156116fd573d6000803e3d6000fd5b505050506040513d602081101561171357600080fd5b8101908080519060200190929190505050101561177b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602e815260200180612011602e913960400191505060405180910390fd5b5b611785866114f5565b6117da576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260238152602001806120846023913960400191505060405180910390fd5b60006001600088815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905060006001600089815260200190815260200160002060020160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690506000600160008a81526020019081526020016000206003015490506000600160008b81526020019081526020016000206004015490506000600160008c815260200190815260200160002060050160006101000a81548160ff021916908315150217905550600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415611936578373ffffffffffffffffffffffffffffffffffffffff166108fc839081150290604051600060405180830381858888f19350505050158015611930573d6000803e3d6000fd5b50611a4e565b8273ffffffffffffffffffffffffffffffffffffffff1663a9059cbb85846040518363ffffffff1660e01b8152600401808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b1580156119bd57600080fd5b505af11580156119d1573d6000803e3d6000fd5b505050506040513d60208110156119e757600080fd5b8101908080519060200190929190505050611a4d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252604581526020018061203f6045913960600191505060405180910390fd5b5b83838383985098509850985050505050509193509193565b60008060008154809291906001019190505550600085858585600054604051602001808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1660601b815260140185805190602001908083835b60208310611aed5780518252602082019150602081019050602083039250611aca565b6001836020036101000a0380198251168184511680821785525050505050509050018473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1660601b8152601401838152602001828152602001955050505050506040516020818303038152906040528051906020012090506040518060c001604052808773ffffffffffffffffffffffffffffffffffffffff1681526020018681526020018573ffffffffffffffffffffffffffffffffffffffff1681526020018481526020016000548152602001600115158152506001600083815260200190815260200160002060008201518160000160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506020820151816001019080519060200190611c43929190611eb5565b5060408201518160020160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550606082015181600301556080820151816004015560a08201518160050160006101000a81548160ff02191690831515021790555090505080915050949350505050565b60008054905090565b600060606000806000611ce8611f35565b600160008881526020019081526020016000206040518060c00160405290816000820160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001600182018054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015611df55780601f10611dca57610100808354040283529160200191611df5565b820191906000526020600020905b815481529060010190602001808311611dd857829003601f168201915b505050505081526020016002820160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200160038201548152602001600482015481526020016005820160009054906101000a900460ff161515151581525050905080600001518160200151826040015183606001518460800151839350955095509550955095505091939590929450565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10611ef657805160ff1916838001178555611f24565b82800160010185558215611f24579182015b82811115611f23578251825591602001919060010190611f08565b5b509050611f319190611f99565b5090565b6040518060c00160405280600073ffffffffffffffffffffffffffffffffffffffff16815260200160608152602001600073ffffffffffffffffffffffffffffffffffffffff16815260200160008152602001600081526020016000151581525090565b611fbb91905b80821115611fb7576000816000905550600101611f9f565b5090565b9056fe496e73756666696369656e7420657468657265756d2062616c616e636520666f722064656c69766572792e4c6f636b2066756e6374696f6e616c6974792069732063757272656e746c79207061757365642e496e73756666696369656e7420455243323020746f6b656e2062616c616e636520666f722064656c69766572792e546f6b656e207472616e73666572206661696c65642c20636865636b20636f6e747261637420746f6b656e20616c6c6f77616e63657320616e642074727920616761696e2e5468652066756e6473206d7573742063757272656e746c79206265206c6f636b65642ea165627a7a723058202f9ac9f9ba437d67811284991fcf36df5d1ac6a0d0e8c08675f5853ebc09b7c50029`

// DeployPeggy deploys a new Ethereum contract, binding an instance of Peggy to it.
func DeployPeggy(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Peggy, error) {
	parsed, err := abi.JSON(strings.NewReader(PeggyABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(PeggyBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Peggy{PeggyCaller: PeggyCaller{contract: contract}, PeggyTransactor: PeggyTransactor{contract: contract}, PeggyFilterer: PeggyFilterer{contract: contract}}, nil
}

// Peggy is an auto generated Go binding around an Ethereum contract.
type Peggy struct {
	PeggyCaller     // Read-only binding to the contract
	PeggyTransactor // Write-only binding to the contract
	PeggyFilterer   // Log filterer for contract events
}

// PeggyCaller is an auto generated read-only Go binding around an Ethereum contract.
type PeggyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PeggyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PeggyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PeggyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PeggyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PeggySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PeggySession struct {
	Contract     *Peggy            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PeggyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PeggyCallerSession struct {
	Contract *PeggyCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// PeggyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PeggyTransactorSession struct {
	Contract     *PeggyTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PeggyRaw is an auto generated low-level Go binding around an Ethereum contract.
type PeggyRaw struct {
	Contract *Peggy // Generic contract binding to access the raw methods on
}

// PeggyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PeggyCallerRaw struct {
	Contract *PeggyCaller // Generic read-only contract binding to access the raw methods on
}

// PeggyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PeggyTransactorRaw struct {
	Contract *PeggyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPeggy creates a new instance of Peggy, bound to a specific deployed contract.
func NewPeggy(address common.Address, backend bind.ContractBackend) (*Peggy, error) {
	contract, err := bindPeggy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Peggy{PeggyCaller: PeggyCaller{contract: contract}, PeggyTransactor: PeggyTransactor{contract: contract}, PeggyFilterer: PeggyFilterer{contract: contract}}, nil
}

// NewPeggyCaller creates a new read-only instance of Peggy, bound to a specific deployed contract.
func NewPeggyCaller(address common.Address, caller bind.ContractCaller) (*PeggyCaller, error) {
	contract, err := bindPeggy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PeggyCaller{contract: contract}, nil
}

// NewPeggyTransactor creates a new write-only instance of Peggy, bound to a specific deployed contract.
func NewPeggyTransactor(address common.Address, transactor bind.ContractTransactor) (*PeggyTransactor, error) {
	contract, err := bindPeggy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PeggyTransactor{contract: contract}, nil
}

// NewPeggyFilterer creates a new log filterer instance of Peggy, bound to a specific deployed contract.
func NewPeggyFilterer(address common.Address, filterer bind.ContractFilterer) (*PeggyFilterer, error) {
	contract, err := bindPeggy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PeggyFilterer{contract: contract}, nil
}

// bindPeggy binds a generic wrapper to an already deployed contract.
func bindPeggy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(PeggyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Peggy *PeggyRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Peggy.Contract.PeggyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Peggy *PeggyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Peggy.Contract.PeggyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Peggy *PeggyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Peggy.Contract.PeggyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Peggy *PeggyCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Peggy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Peggy *PeggyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Peggy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Peggy *PeggyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Peggy.Contract.contract.Transact(opts, method, params...)
}

// Active is a free data retrieval call binding the contract method 0x02fb0c5e.
//
// Solidity: function active() constant returns(bool)
func (_Peggy *PeggyCaller) Active(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Peggy.contract.Call(opts, out, "active")
	return *ret0, err
}

// Active is a free data retrieval call binding the contract method 0x02fb0c5e.
//
// Solidity: function active() constant returns(bool)
func (_Peggy *PeggySession) Active() (bool, error) {
	return _Peggy.Contract.Active(&_Peggy.CallOpts)
}

// Active is a free data retrieval call binding the contract method 0x02fb0c5e.
//
// Solidity: function active() constant returns(bool)
func (_Peggy *PeggyCallerSession) Active() (bool, error) {
	return _Peggy.Contract.Active(&_Peggy.CallOpts)
}

// GetStatus is a free data retrieval call binding the contract method 0x5de28ae0.
//
// Solidity: function getStatus(bytes32 _id) constant returns(bool)
func (_Peggy *PeggyCaller) GetStatus(opts *bind.CallOpts, _id [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Peggy.contract.Call(opts, out, "getStatus", _id)
	return *ret0, err
}

// GetStatus is a free data retrieval call binding the contract method 0x5de28ae0.
//
// Solidity: function getStatus(bytes32 _id) constant returns(bool)
func (_Peggy *PeggySession) GetStatus(_id [32]byte) (bool, error) {
	return _Peggy.Contract.GetStatus(&_Peggy.CallOpts, _id)
}

// GetStatus is a free data retrieval call binding the contract method 0x5de28ae0.
//
// Solidity: function getStatus(bytes32 _id) constant returns(bool)
func (_Peggy *PeggyCallerSession) GetStatus(_id [32]byte) (bool, error) {
	return _Peggy.Contract.GetStatus(&_Peggy.CallOpts, _id)
}

// Ids is a free data retrieval call binding the contract method 0xcf7b4a09.
//
// Solidity: function ids(bytes32 ) constant returns(bool)
func (_Peggy *PeggyCaller) Ids(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Peggy.contract.Call(opts, out, "ids", arg0)
	return *ret0, err
}

// Ids is a free data retrieval call binding the contract method 0xcf7b4a09.
//
// Solidity: function ids(bytes32 ) constant returns(bool)
func (_Peggy *PeggySession) Ids(arg0 [32]byte) (bool, error) {
	return _Peggy.Contract.Ids(&_Peggy.CallOpts, arg0)
}

// Ids is a free data retrieval call binding the contract method 0xcf7b4a09.
//
// Solidity: function ids(bytes32 ) constant returns(bool)
func (_Peggy *PeggyCallerSession) Ids(arg0 [32]byte) (bool, error) {
	return _Peggy.Contract.Ids(&_Peggy.CallOpts, arg0)
}

// Nonce is a free data retrieval call binding the contract method 0xaffed0e0.
//
// Solidity: function nonce() constant returns(uint256)
func (_Peggy *PeggyCaller) Nonce(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Peggy.contract.Call(opts, out, "nonce")
	return *ret0, err
}

// Nonce is a free data retrieval call binding the contract method 0xaffed0e0.
//
// Solidity: function nonce() constant returns(uint256)
func (_Peggy *PeggySession) Nonce() (*big.Int, error) {
	return _Peggy.Contract.Nonce(&_Peggy.CallOpts)
}

// Nonce is a free data retrieval call binding the contract method 0xaffed0e0.
//
// Solidity: function nonce() constant returns(uint256)
func (_Peggy *PeggyCallerSession) Nonce() (*big.Int, error) {
	return _Peggy.Contract.Nonce(&_Peggy.CallOpts)
}

// Provider is a free data retrieval call binding the contract method 0x085d4883.
//
// Solidity: function provider() constant returns(address)
func (_Peggy *PeggyCaller) Provider(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Peggy.contract.Call(opts, out, "provider")
	return *ret0, err
}

// Provider is a free data retrieval call binding the contract method 0x085d4883.
//
// Solidity: function provider() constant returns(address)
func (_Peggy *PeggySession) Provider() (common.Address, error) {
	return _Peggy.Contract.Provider(&_Peggy.CallOpts)
}

// Provider is a free data retrieval call binding the contract method 0x085d4883.
//
// Solidity: function provider() constant returns(address)
func (_Peggy *PeggyCallerSession) Provider() (common.Address, error) {
	return _Peggy.Contract.Provider(&_Peggy.CallOpts)
}

// ViewItem is a free data retrieval call binding the contract method 0xc933dc5b.
//
// Solidity: function viewItem(bytes32 _id) constant returns(address, bytes, address, uint256, uint256)
func (_Peggy *PeggyCaller) ViewItem(opts *bind.CallOpts, _id [32]byte) (common.Address, []byte, common.Address, *big.Int, *big.Int, error) {
	var (
		ret0 = new(common.Address)
		ret1 = new([]byte)
		ret2 = new(common.Address)
		ret3 = new(*big.Int)
		ret4 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
		ret3,
		ret4,
	}
	err := _Peggy.contract.Call(opts, out, "viewItem", _id)
	return *ret0, *ret1, *ret2, *ret3, *ret4, err
}

// ViewItem is a free data retrieval call binding the contract method 0xc933dc5b.
//
// Solidity: function viewItem(bytes32 _id) constant returns(address, bytes, address, uint256, uint256)
func (_Peggy *PeggySession) ViewItem(_id [32]byte) (common.Address, []byte, common.Address, *big.Int, *big.Int, error) {
	return _Peggy.Contract.ViewItem(&_Peggy.CallOpts, _id)
}

// ViewItem is a free data retrieval call binding the contract method 0xc933dc5b.
//
// Solidity: function viewItem(bytes32 _id) constant returns(address, bytes, address, uint256, uint256)
func (_Peggy *PeggyCallerSession) ViewItem(_id [32]byte) (common.Address, []byte, common.Address, *big.Int, *big.Int, error) {
	return _Peggy.Contract.ViewItem(&_Peggy.CallOpts, _id)
}

// ActivateLocking is a paid mutator transaction binding the contract method 0x63faf36a.
//
// Solidity: function activateLocking() returns()
func (_Peggy *PeggyTransactor) ActivateLocking(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Peggy.contract.Transact(opts, "activateLocking")
}

// ActivateLocking is a paid mutator transaction binding the contract method 0x63faf36a.
//
// Solidity: function activateLocking() returns()
func (_Peggy *PeggySession) ActivateLocking() (*types.Transaction, error) {
	return _Peggy.Contract.ActivateLocking(&_Peggy.TransactOpts)
}

// ActivateLocking is a paid mutator transaction binding the contract method 0x63faf36a.
//
// Solidity: function activateLocking() returns()
func (_Peggy *PeggyTransactorSession) ActivateLocking() (*types.Transaction, error) {
	return _Peggy.Contract.ActivateLocking(&_Peggy.TransactOpts)
}

// Lock is a paid mutator transaction binding the contract method 0x9df2a385.
//
// Solidity: function lock(bytes _recipient, address _token, uint256 _amount) returns(bytes32 _id)
func (_Peggy *PeggyTransactor) Lock(opts *bind.TransactOpts, _recipient []byte, _token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Peggy.contract.Transact(opts, "lock", _recipient, _token, _amount)
}

// Lock is a paid mutator transaction binding the contract method 0x9df2a385.
//
// Solidity: function lock(bytes _recipient, address _token, uint256 _amount) returns(bytes32 _id)
func (_Peggy *PeggySession) Lock(_recipient []byte, _token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Peggy.Contract.Lock(&_Peggy.TransactOpts, _recipient, _token, _amount)
}

// Lock is a paid mutator transaction binding the contract method 0x9df2a385.
//
// Solidity: function lock(bytes _recipient, address _token, uint256 _amount) returns(bytes32 _id)
func (_Peggy *PeggyTransactorSession) Lock(_recipient []byte, _token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Peggy.Contract.Lock(&_Peggy.TransactOpts, _recipient, _token, _amount)
}

// PauseLocking is a paid mutator transaction binding the contract method 0x8a5cd91e.
//
// Solidity: function pauseLocking() returns()
func (_Peggy *PeggyTransactor) PauseLocking(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Peggy.contract.Transact(opts, "pauseLocking")
}

// PauseLocking is a paid mutator transaction binding the contract method 0x8a5cd91e.
//
// Solidity: function pauseLocking() returns()
func (_Peggy *PeggySession) PauseLocking() (*types.Transaction, error) {
	return _Peggy.Contract.PauseLocking(&_Peggy.TransactOpts)
}

// PauseLocking is a paid mutator transaction binding the contract method 0x8a5cd91e.
//
// Solidity: function pauseLocking() returns()
func (_Peggy *PeggyTransactorSession) PauseLocking() (*types.Transaction, error) {
	return _Peggy.Contract.PauseLocking(&_Peggy.TransactOpts)
}

// Unlock is a paid mutator transaction binding the contract method 0xec9b5b3a.
//
// Solidity: function unlock(bytes32 _id) returns(bool)
func (_Peggy *PeggyTransactor) Unlock(opts *bind.TransactOpts, _id [32]byte) (*types.Transaction, error) {
	return _Peggy.contract.Transact(opts, "unlock", _id)
}

// Unlock is a paid mutator transaction binding the contract method 0xec9b5b3a.
//
// Solidity: function unlock(bytes32 _id) returns(bool)
func (_Peggy *PeggySession) Unlock(_id [32]byte) (*types.Transaction, error) {
	return _Peggy.Contract.Unlock(&_Peggy.TransactOpts, _id)
}

// Unlock is a paid mutator transaction binding the contract method 0xec9b5b3a.
//
// Solidity: function unlock(bytes32 _id) returns(bool)
func (_Peggy *PeggyTransactorSession) Unlock(_id [32]byte) (*types.Transaction, error) {
	return _Peggy.Contract.Unlock(&_Peggy.TransactOpts, _id)
}

// Withdraw is a paid mutator transaction binding the contract method 0x8e19899e.
//
// Solidity: function withdraw(bytes32 _id) returns(bool)
func (_Peggy *PeggyTransactor) Withdraw(opts *bind.TransactOpts, _id [32]byte) (*types.Transaction, error) {
	return _Peggy.contract.Transact(opts, "withdraw", _id)
}

// Withdraw is a paid mutator transaction binding the contract method 0x8e19899e.
//
// Solidity: function withdraw(bytes32 _id) returns(bool)
func (_Peggy *PeggySession) Withdraw(_id [32]byte) (*types.Transaction, error) {
	return _Peggy.Contract.Withdraw(&_Peggy.TransactOpts, _id)
}

// Withdraw is a paid mutator transaction binding the contract method 0x8e19899e.
//
// Solidity: function withdraw(bytes32 _id) returns(bool)
func (_Peggy *PeggyTransactorSession) Withdraw(_id [32]byte) (*types.Transaction, error) {
	return _Peggy.Contract.Withdraw(&_Peggy.TransactOpts, _id)
}

// PeggyLogLockIterator is returned from FilterLogLock and is used to iterate over the raw logs and unpacked data for LogLock events raised by the Peggy contract.
type PeggyLogLockIterator struct {
	Event *PeggyLogLock // Event containing the contract specifics and raw log

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
func (it *PeggyLogLockIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PeggyLogLock)
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
		it.Event = new(PeggyLogLock)
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
func (it *PeggyLogLockIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PeggyLogLockIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PeggyLogLock represents a LogLock event raised by the Peggy contract.
type PeggyLogLock struct {
	Id     [32]byte
	From   common.Address
	To     []byte
	Token  common.Address
	Symbol string
	Value  *big.Int
	Nonce  *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterLogLock is a free log retrieval operation binding the contract event 0x3945646e76891f1dfa38b4aab98fac226e3f4ad3686493b722b62383358ba922.
//
// Solidity: event LogLock(bytes32 _id, address _from, bytes _to, address _token, string _symbol, uint256 _value, uint256 _nonce)
func (_Peggy *PeggyFilterer) FilterLogLock(opts *bind.FilterOpts) (*PeggyLogLockIterator, error) {

	logs, sub, err := _Peggy.contract.FilterLogs(opts, "LogLock")
	if err != nil {
		return nil, err
	}
	return &PeggyLogLockIterator{contract: _Peggy.contract, event: "LogLock", logs: logs, sub: sub}, nil
}

// WatchLogLock is a free log subscription operation binding the contract event 0x3945646e76891f1dfa38b4aab98fac226e3f4ad3686493b722b62383358ba922.
//
// Solidity: event LogLock(bytes32 _id, address _from, bytes _to, address _token, string _symbol, uint256 _value, uint256 _nonce)
func (_Peggy *PeggyFilterer) WatchLogLock(opts *bind.WatchOpts, sink chan<- *PeggyLogLock) (event.Subscription, error) {

	logs, sub, err := _Peggy.contract.WatchLogs(opts, "LogLock")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PeggyLogLock)
				if err := _Peggy.contract.UnpackLog(event, "LogLock", log); err != nil {
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

// PeggyLogLockingActivatedIterator is returned from FilterLogLockingActivated and is used to iterate over the raw logs and unpacked data for LogLockingActivated events raised by the Peggy contract.
type PeggyLogLockingActivatedIterator struct {
	Event *PeggyLogLockingActivated // Event containing the contract specifics and raw log

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
func (it *PeggyLogLockingActivatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PeggyLogLockingActivated)
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
		it.Event = new(PeggyLogLockingActivated)
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
func (it *PeggyLogLockingActivatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PeggyLogLockingActivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PeggyLogLockingActivated represents a LogLockingActivated event raised by the Peggy contract.
type PeggyLogLockingActivated struct {
	Time *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogLockingActivated is a free log retrieval operation binding the contract event 0x9af033c3fdf318cb9968eac8a62b339bd18862abd1703fc74256e9d77cfc95df.
//
// Solidity: event LogLockingActivated(uint256 _time)
func (_Peggy *PeggyFilterer) FilterLogLockingActivated(opts *bind.FilterOpts) (*PeggyLogLockingActivatedIterator, error) {

	logs, sub, err := _Peggy.contract.FilterLogs(opts, "LogLockingActivated")
	if err != nil {
		return nil, err
	}
	return &PeggyLogLockingActivatedIterator{contract: _Peggy.contract, event: "LogLockingActivated", logs: logs, sub: sub}, nil
}

// WatchLogLockingActivated is a free log subscription operation binding the contract event 0x9af033c3fdf318cb9968eac8a62b339bd18862abd1703fc74256e9d77cfc95df.
//
// Solidity: event LogLockingActivated(uint256 _time)
func (_Peggy *PeggyFilterer) WatchLogLockingActivated(opts *bind.WatchOpts, sink chan<- *PeggyLogLockingActivated) (event.Subscription, error) {

	logs, sub, err := _Peggy.contract.WatchLogs(opts, "LogLockingActivated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PeggyLogLockingActivated)
				if err := _Peggy.contract.UnpackLog(event, "LogLockingActivated", log); err != nil {
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

// PeggyLogLockingPausedIterator is returned from FilterLogLockingPaused and is used to iterate over the raw logs and unpacked data for LogLockingPaused events raised by the Peggy contract.
type PeggyLogLockingPausedIterator struct {
	Event *PeggyLogLockingPaused // Event containing the contract specifics and raw log

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
func (it *PeggyLogLockingPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PeggyLogLockingPaused)
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
		it.Event = new(PeggyLogLockingPaused)
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
func (it *PeggyLogLockingPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PeggyLogLockingPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PeggyLogLockingPaused represents a LogLockingPaused event raised by the Peggy contract.
type PeggyLogLockingPaused struct {
	Time *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterLogLockingPaused is a free log retrieval operation binding the contract event 0xbebc9a19c81e5697fda01edce5ac5aed2c5a0edb9a972fd5f58ac0419a405a82.
//
// Solidity: event LogLockingPaused(uint256 _time)
func (_Peggy *PeggyFilterer) FilterLogLockingPaused(opts *bind.FilterOpts) (*PeggyLogLockingPausedIterator, error) {

	logs, sub, err := _Peggy.contract.FilterLogs(opts, "LogLockingPaused")
	if err != nil {
		return nil, err
	}
	return &PeggyLogLockingPausedIterator{contract: _Peggy.contract, event: "LogLockingPaused", logs: logs, sub: sub}, nil
}

// WatchLogLockingPaused is a free log subscription operation binding the contract event 0xbebc9a19c81e5697fda01edce5ac5aed2c5a0edb9a972fd5f58ac0419a405a82.
//
// Solidity: event LogLockingPaused(uint256 _time)
func (_Peggy *PeggyFilterer) WatchLogLockingPaused(opts *bind.WatchOpts, sink chan<- *PeggyLogLockingPaused) (event.Subscription, error) {

	logs, sub, err := _Peggy.contract.WatchLogs(opts, "LogLockingPaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PeggyLogLockingPaused)
				if err := _Peggy.contract.UnpackLog(event, "LogLockingPaused", log); err != nil {
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

// PeggyLogUnlockIterator is returned from FilterLogUnlock and is used to iterate over the raw logs and unpacked data for LogUnlock events raised by the Peggy contract.
type PeggyLogUnlockIterator struct {
	Event *PeggyLogUnlock // Event containing the contract specifics and raw log

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
func (it *PeggyLogUnlockIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PeggyLogUnlock)
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
		it.Event = new(PeggyLogUnlock)
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
func (it *PeggyLogUnlockIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PeggyLogUnlockIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PeggyLogUnlock represents a LogUnlock event raised by the Peggy contract.
type PeggyLogUnlock struct {
	Id    [32]byte
	To    common.Address
	Token common.Address
	Value *big.Int
	Nonce *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterLogUnlock is a free log retrieval operation binding the contract event 0xb3ceeb2ff57376fcabec63d51a010afad847c03e9365f20a168ca66db8b92740.
//
// Solidity: event LogUnlock(bytes32 _id, address _to, address _token, uint256 _value, uint256 _nonce)
func (_Peggy *PeggyFilterer) FilterLogUnlock(opts *bind.FilterOpts) (*PeggyLogUnlockIterator, error) {

	logs, sub, err := _Peggy.contract.FilterLogs(opts, "LogUnlock")
	if err != nil {
		return nil, err
	}
	return &PeggyLogUnlockIterator{contract: _Peggy.contract, event: "LogUnlock", logs: logs, sub: sub}, nil
}

// WatchLogUnlock is a free log subscription operation binding the contract event 0xb3ceeb2ff57376fcabec63d51a010afad847c03e9365f20a168ca66db8b92740.
//
// Solidity: event LogUnlock(bytes32 _id, address _to, address _token, uint256 _value, uint256 _nonce)
func (_Peggy *PeggyFilterer) WatchLogUnlock(opts *bind.WatchOpts, sink chan<- *PeggyLogUnlock) (event.Subscription, error) {

	logs, sub, err := _Peggy.contract.WatchLogs(opts, "LogUnlock")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PeggyLogUnlock)
				if err := _Peggy.contract.UnpackLog(event, "LogUnlock", log); err != nil {
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

// PeggyLogWithdrawIterator is returned from FilterLogWithdraw and is used to iterate over the raw logs and unpacked data for LogWithdraw events raised by the Peggy contract.
type PeggyLogWithdrawIterator struct {
	Event *PeggyLogWithdraw // Event containing the contract specifics and raw log

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
func (it *PeggyLogWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PeggyLogWithdraw)
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
		it.Event = new(PeggyLogWithdraw)
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
func (it *PeggyLogWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PeggyLogWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PeggyLogWithdraw represents a LogWithdraw event raised by the Peggy contract.
type PeggyLogWithdraw struct {
	Id    [32]byte
	To    common.Address
	Token common.Address
	Value *big.Int
	Nonce *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterLogWithdraw is a free log retrieval operation binding the contract event 0x9cbca76b94cf51b34c3949f0c925da38fe8dbae8e6761e11389884a9c1354b2c.
//
// Solidity: event LogWithdraw(bytes32 _id, address _to, address _token, uint256 _value, uint256 _nonce)
func (_Peggy *PeggyFilterer) FilterLogWithdraw(opts *bind.FilterOpts) (*PeggyLogWithdrawIterator, error) {

	logs, sub, err := _Peggy.contract.FilterLogs(opts, "LogWithdraw")
	if err != nil {
		return nil, err
	}
	return &PeggyLogWithdrawIterator{contract: _Peggy.contract, event: "LogWithdraw", logs: logs, sub: sub}, nil
}

// WatchLogWithdraw is a free log subscription operation binding the contract event 0x9cbca76b94cf51b34c3949f0c925da38fe8dbae8e6761e11389884a9c1354b2c.
//
// Solidity: event LogWithdraw(bytes32 _id, address _to, address _token, uint256 _value, uint256 _nonce)
func (_Peggy *PeggyFilterer) WatchLogWithdraw(opts *bind.WatchOpts, sink chan<- *PeggyLogWithdraw) (event.Subscription, error) {

	logs, sub, err := _Peggy.contract.WatchLogs(opts, "LogWithdraw")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PeggyLogWithdraw)
				if err := _Peggy.contract.UnpackLog(event, "LogWithdraw", log); err != nil {
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
